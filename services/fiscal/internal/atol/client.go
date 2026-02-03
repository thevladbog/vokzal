// Package atol — клиент АТОЛ ККТ через локальный агент.
package atol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// ATOLClient — клиент для работы с АТОЛ через локальный агент.
//
//nolint:revive // exported: имя ATOLClient намеренно (Client слишком общее в пакете atol).
type ATOLClient struct {
	client   *http.Client
	logger   *zap.Logger
	agentURL string
}

// ReceiptRequest — запрос на печать чека.
type ReceiptRequest struct {
	Company   Company       `json:"company"`
	Operation string        `json:"operation"`
	Items     []ReceiptItem `json:"items"`
	Payment   Payment       `json:"payment"`
}

// ReceiptItem — позиция чека.
type ReceiptItem struct {
	Name     string  `json:"name"`
	VAT      string  `json:"vat"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
}

// Payment — платёж в чеке.
type Payment struct {
	Type   string  `json:"type"` // cash, card, online
	Amount float64 `json:"amount"`
}

// Company — реквизиты организации для чека.
type Company struct {
	INN       string `json:"inn"`
	Name      string `json:"name"`
	TaxSystem string `json:"tax_system"`
}

// ReceiptResponse — ответ от ККТ.
type ReceiptResponse struct {
	OFDURL     string `json:"ofd_url"`
	FiscalSign string `json:"fiscal_sign"`
	KKTSerial  string `json:"kkt_serial"`
	ErrorMsg   string `json:"error_msg,omitempty"`
	Success    bool   `json:"success"`
}

// ZReportResponse — ответ на Z-отчёт.
type ZReportResponse struct {
	FiscalSign   string  `json:"fiscal_sign"`
	KKTSerial    string  `json:"kkt_serial,omitempty"`
	ErrorMsg     string  `json:"error_msg,omitempty"`
	ShiftNumber  int     `json:"shift_number"`
	TotalSales   float64 `json:"total_sales"`
	TotalRefunds float64 `json:"total_refunds"`
	SalesCount   int     `json:"sales_count"`
	RefundsCount int     `json:"refunds_count"`
	Success      bool    `json:"success"`
}

// NewATOLClient создаёт новый ATOLClient.
func NewATOLClient(agentURL string, logger *zap.Logger) *ATOLClient {
	return &ATOLClient{
		agentURL: agentURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// PrintReceipt отправляет запрос на печать чека через локальный агент.
func (c *ATOLClient) PrintReceipt(req *ReceiptRequest) (*ReceiptResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/kkt/receipt", c.agentURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	c.logger.Debug("Sending receipt to KKT", zap.String("url", url))

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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("KKT returned status %d: %s", resp.StatusCode, string(body))
	}

	var result ReceiptResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateZReport формирует Z-отчёт.
func (c *ATOLClient) CreateZReport() (*ZReportResponse, error) {
	url := fmt.Sprintf("%s/kkt/z-report", c.agentURL)
	httpReq, err := http.NewRequest("POST", url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.logger.Debug("Creating Z-report", zap.String("url", url))

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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("KKT returned status %d: %s", resp.StatusCode, string(body))
	}

	var result ZReportResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetKKTStatus получает статус ККТ.
func (c *ATOLClient) GetKKTStatus() (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/kkt/status", c.agentURL)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
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

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}
