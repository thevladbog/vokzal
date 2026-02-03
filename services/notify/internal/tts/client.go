// Package tts предоставляет клиент для голосовых оповещений (TTS) через локальный агент.
package tts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// TTSClient — клиент для отправки TTS-запросов на локальный агент.
//
//nolint:revive // Имя сохраняем для ясности при использовании (tts.Client).
type TTSClient struct {
	agentURL string
	client   *http.Client
	logger   *zap.Logger
}

// TTSRequest — тело запроса к TTS-агенту.
//
//nolint:revive // Имя сохраняем для ясности (tts.Request).
type TTSRequest struct {
	Text     string `json:"text"`
	Language string `json:"language"`
	Priority string `json:"priority"`
}

// NewTTSClient создаёт клиент TTS.
func NewTTSClient(agentURL string, logger *zap.Logger) *TTSClient {
	return &TTSClient{
		agentURL: agentURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// Announce отправляет голосовое оповещение (TTS).
func (c *TTSClient) Announce(text, language, priority string) error {
	req := &TTSRequest{
		Text:     text,
		Language: language,
		Priority: priority,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/tts/announce", c.agentURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	c.logger.Debug("Sending TTS announce", zap.String("text", text))

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("TTS returned status %d: %s", resp.StatusCode, string(body))
	}

	c.logger.Info("TTS announce sent", zap.String("language", language))
	return nil
}
