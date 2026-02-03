// Package sms предоставляет клиент для отправки SMS через SMS.ru API.
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

// SMSRuClient — клиент для SMS.ru API.
//
//nolint:revive // Имя RuClient сохраняем для однозначности (sms.RuClient).
type SMSRuClient struct {
	client *http.Client
	logger *zap.Logger
	apiID  string
	apiURL string
}

// SMSResponse — ответ API SMS.ru.
//
//nolint:revive // Имя сохраняем для ясности (sms.Response).
type SMSResponse struct {
	SMS        map[string]SMSStatus `json:"sms"`
	Status     string               `json:"status"`
	StatusCode int                  `json:"status_code"`
}

// SMSStatus — статус одного SMS в ответе.
//
//nolint:revive // Имя сохраняем для ясности (sms.Status).
type SMSStatus struct {
	Status     string `json:"status"`
	SMSID      string `json:"sms_id"`
	StatusCode int    `json:"status_code"`
}

// NewSMSRuClient создаёт клиент SMS.ru.
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

// Send отправляет SMS на указанный номер.
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
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Warn("failed to close response body", zap.Error(closeErr))
		}
	}()

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
