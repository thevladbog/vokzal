package kkt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// ATOLClient клиент для работы с ККТ через АТОЛ драйвер
type ATOLClient struct {
	driverURL string
	inn       string
	ofdURL    string
	enabled   bool
	logger    *zap.Logger
	client    *http.Client
}

type ReceiptRequest struct {
	Operation string         `json:"operation"` // sell, refund
	Items     []ReceiptItem  `json:"items"`
	Payment   PaymentInfo    `json:"payment"`
	INN       string         `json:"inn"`
}

type ReceiptItem struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
	VAT      string  `json:"vat"` // vat20, vat10, vat0, none
}

type PaymentInfo struct {
	Type   string  `json:"type"`   // cash, card, sbp
	Amount float64 `json:"amount"`
}

type ReceiptResponse struct {
	Success      bool   `json:"success"`
	FiscalSign   string `json:"fiscal_sign"`
	OFDURL       string `json:"ofd_url"`
	ReceiptNum   int    `json:"receipt_num"`
	ShiftNum     int    `json:"shift_num"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type ZReportResponse struct {
	Success    bool   `json:"success"`
	ShiftNum   int    `json:"shift_num"`
	TotalSales float64 `json:"total_sales"`
	DocumentURL string `json:"document_url"`
}

func NewATOLClient(driverURL, inn, ofdURL string, enabled bool, logger *zap.Logger) *ATOLClient {
	return &ATOLClient{
		driverURL: driverURL,
		inn:       inn,
		ofdURL:    ofdURL,
		enabled:   enabled,
		logger:    logger,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// PrintReceipt печатает чек на ККТ
func (c *ATOLClient) PrintReceipt(req *ReceiptRequest) (*ReceiptResponse, error) {
	if !c.enabled {
		c.logger.Info("KKT disabled, simulating receipt")
		return &ReceiptResponse{
			Success:    true,
			FiscalSign: "SIMULATED_" + time.Now().Format("20060102150405"),
			OFDURL:     "http://check.ofd.ru/simulated",
			ReceiptNum: int(time.Now().Unix() % 10000),
			ShiftNum:   1,
		}, nil
	}

	req.INN = c.inn

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/receipt", c.driverURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	c.logger.Info("Printing receipt on KKT", zap.String("operation", req.Operation))

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to KKT: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result ReceiptResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("KKT error: %s", result.ErrorMessage)
	}

	c.logger.Info("Receipt printed successfully", zap.String("fiscal_sign", result.FiscalSign))
	return &result, nil
}

// CreateZReport создаёт Z-отчёт
func (c *ATOLClient) CreateZReport() (*ZReportResponse, error) {
	if !c.enabled {
		c.logger.Info("KKT disabled, simulating Z-report")
		return &ZReportResponse{
			Success:     true,
			ShiftNum:    1,
			TotalSales:  15000.00,
			DocumentURL: "http://check.ofd.ru/simulated/z-report",
		}, nil
	}

	url := fmt.Sprintf("%s/z-report", c.driverURL)
	httpReq, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.logger.Info("Creating Z-report")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to KKT: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result ZReportResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("Z-report creation failed")
	}

	c.logger.Info("Z-report created", zap.Int("shift", result.ShiftNum))
	return &result, nil
}

// GetStatus получает статус ККТ
func (c *ATOLClient) GetStatus() (map[string]interface{}, error) {
	if !c.enabled {
		return map[string]interface{}{
			"status":     "simulated",
			"shift_open": true,
			"fn_status":  "ok",
		}, nil
	}

	url := fmt.Sprintf("%s/status", c.driverURL)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get KKT status: %w", err)
	}
	defer resp.Body.Close()

	var status map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode status: %w", err)
	}

	return status, nil
}
