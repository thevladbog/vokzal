// Package tinkoff provides a client for the Tinkoff Acquiring API.
package tinkoff

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap"
)

// TinkoffClient is a client for interacting with the Tinkoff Acquiring API.
//
//nolint:revive // Name preserved for clarity (tinkoff.Client).
type TinkoffClient struct {
	client      *http.Client
	logger      *zap.Logger
	terminalKey string
	apiSecret   string // API secret (not a user password) used to sign requests
	apiURL      string // API base URL
}

// InitRequest is a payment initialization request.
type InitRequest struct {
	TerminalKey     string `json:"TerminalKey"`
	OrderID         string `json:"OrderId"`
	Description     string `json:"Description"`
	Token           string `json:"Token"`
	NotificationURL string `json:"NotificationURL,omitempty"`
	SuccessURL      string `json:"SuccessURL,omitempty"`
	FailURL         string `json:"FailURL,omitempty"`
	Amount          int64  `json:"Amount"`
}

// InitResponse is a payment initialization response.
type InitResponse struct {
	ErrorCode  string `json:"ErrorCode"`
	Message    string `json:"Message"`
	PaymentID  string `json:"PaymentId"`
	PaymentURL string `json:"PaymentURL"`
	Success    bool   `json:"Success"`
}

// GetStateRequest is a payment status request.
type GetStateRequest struct {
	TerminalKey string `json:"TerminalKey"`
	PaymentID   string `json:"PaymentId"`
	Token       string `json:"Token"`
}

// GetStateResponse is a payment status response.
type GetStateResponse struct {
	ErrorCode string `json:"ErrorCode"`
	Message   string `json:"Message"`
	Status    string `json:"Status"`
	PaymentID string `json:"PaymentId"`
	OrderID   string `json:"OrderId"`
	Amount    int64  `json:"Amount"`
	Success   bool   `json:"Success"`
}

// NewTinkoffClient creates a Tinkoff Acquiring client.
// apiSecret is the terminal API secret (not a user password) used to sign API requests.
func NewTinkoffClient(terminalKey, apiSecret, apiURL string, logger *zap.Logger) *TinkoffClient {
	return &TinkoffClient{
		terminalKey: terminalKey,
		apiSecret:   apiSecret,
		apiURL:      apiURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// Init initializes a payment.
func (c *TinkoffClient) Init(orderID string, amount float64, description string) (*InitResponse, error) {
	amountKopecks := int64(amount * 100)

	req := &InitRequest{
		TerminalKey: c.terminalKey,
		Amount:      amountKopecks,
		OrderID:     orderID,
		Description: description,
	}

	// Generate token
	req.Token = c.generateToken(map[string]interface{}{
		"TerminalKey": req.TerminalKey,
		"Amount":      req.Amount,
		"OrderId":     req.OrderID,
	})

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/Init", c.apiURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	c.logger.Debug("Tinkoff Init request", zap.String("url", url), zap.String("order_id", orderID))

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Warn("failed to close response body", zap.Error(closeErr))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result InitResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("tinkoff error: %s - %s", result.ErrorCode, result.Message)
	}

	return &result, nil
}

// GetState retrieves the payment status.
func (c *TinkoffClient) GetState(paymentID string) (*GetStateResponse, error) {
	req := &GetStateRequest{
		TerminalKey: c.terminalKey,
		PaymentID:   paymentID,
	}

	req.Token = c.generateToken(map[string]interface{}{
		"TerminalKey": req.TerminalKey,
		"PaymentId":   req.PaymentID,
	})

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/GetState", c.apiURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Warn("failed to close response body", zap.Error(closeErr))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result GetStateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// generateToken generates a token for signing requests according to the Tinkoff Acquiring API specification.
//
// SECURITY NOTE: SHA-256 is used here for generating API request signatures (HMAC-like mechanism),
// NOT for hashing user passwords. This complies with the official Tinkoff documentation.
// SHA-256 is a cryptographically strong hash function for this use case.
func (c *TinkoffClient) generateToken(params map[string]interface{}) string {
	// Add API secret to parameters (Tinkoff API requirement)
	params["Password"] = c.apiSecret

	// Sort keys
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build string for hashing
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%v", params[k]))
	}
	str := strings.Join(parts, "")

	// SHA-256 hashing (Tinkoff Acquiring API requirement)
	hash := sha256.Sum256([]byte(str))
	return fmt.Sprintf("%x", hash)
}
