// Package sbp предоставляет клиент для СБП (Система быстрых платежей).
package sbp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// SBPClient — клиент для работы с СБП.
//
//nolint:revive // Имя сохраняем для ясности (sbp.Client).
type SBPClient struct {
	client     *http.Client
	logger     *zap.Logger
	merchantID string
	apiURL     string
	apiKey     string
}

// QRRequest — запрос на генерацию QR-кода.
type QRRequest struct {
	MerchantID string  `json:"merchantId"`
	Currency   string  `json:"currency"`
	Purpose    string  `json:"purpose"`
	QRType     string  `json:"qrType"`
	Amount     float64 `json:"amount"`
}

// QRResponse — ответ с QR-кодом.
type QRResponse struct {
	QRCode    string `json:"qrCode"`
	QRString  string `json:"qrString"`
	PaymentID string `json:"paymentId"`
	ErrorMsg  string `json:"errorMsg,omitempty"`
	Success   bool   `json:"success"`
}

// StatusRequest — запрос статуса платежа.
type StatusRequest struct {
	MerchantID string `json:"merchantId"`
	PaymentID  string `json:"paymentId"`
}

// StatusResponse — ответ со статусом платежа.
type StatusResponse struct {
	PaidAt    *time.Time `json:"paidAt,omitempty"`
	Status    string     `json:"status"`
	PaymentID string     `json:"paymentId"`
	ErrorMsg  string     `json:"errorMsg,omitempty"`
	Amount    float64    `json:"amount"`
	Success   bool       `json:"success"`
}

// NewSBPClient создаёт клиент СБП.
func NewSBPClient(merchantID, apiURL, apiKey string, logger *zap.Logger) *SBPClient {
	return &SBPClient{
		merchantID: merchantID,
		apiURL:     apiURL,
		apiKey:     apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// GenerateQR генерирует QR-код для оплаты.
func (c *SBPClient) GenerateQR(amount float64, purpose string) (*QRResponse, error) {
	req := &QRRequest{
		MerchantID: c.merchantID,
		Amount:     amount,
		Currency:   "RUB",
		Purpose:    purpose,
		QRType:     "dynamic",
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/qr/generate", c.apiURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	c.logger.Debug("SBP GenerateQR request", zap.String("url", url), zap.Float64("amount", amount))

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

	var result QRResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("SBP error: %s", result.ErrorMsg)
	}

	return &result, nil
}

// GetStatus получает статус платежа.
func (c *SBPClient) GetStatus(paymentID string) (*StatusResponse, error) {
	req := &StatusRequest{
		MerchantID: c.merchantID,
		PaymentID:  paymentID,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/payment/status", c.apiURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

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

	var result StatusResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}
