package sms

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

// SMSRuClient клиент для SMS.ru API
type SMSRuClient struct {
	apiID  string
	apiURL string
	client *http.Client
	logger *zap.Logger
}

type SMSResponse struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	SMS        map[string]SMSStatus `json:"sms"`
}

type SMSStatus struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	SMSID      string `json:"sms_id"`
}

func NewSMSRuClient(apiID, apiURL string, logger *zap.Logger) *SMSRuClient {
	return &SMSRuClient{
		apiID:  apiID,
		apiURL: apiURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// Send отправляет SMS
func (c *SMSRuClient) Send(phone, message string) error {
	params := url.Values{}
	params.Set("api_id", c.apiID)
	params.Set("to", phone)
	params.Set("msg", message)
	params.Set("json", "1")

	reqURL := fmt.Sprintf("%s?%s", c.apiURL, params.Encode())

	c.logger.Debug("Sending SMS", zap.String("phone", phone))

	resp, err := c.client.Get(reqURL)
	if err != nil {
		return fmt.Errorf("failed to send SMS request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var result SMSResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if result.StatusCode != 100 {
		return fmt.Errorf("SMS.ru error: status %d", result.StatusCode)
	}

	c.logger.Info("SMS sent successfully", zap.String("phone", phone))
	return nil
}
