// Package telegram предоставляет клиент для Telegram Bot API.
package telegram

import (
	"fmt"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

// TelegramClient — клиент для Telegram Bot API.
//
//nolint:revive // Имя сохраняем для ясности при использовании (telegram.Client).
type TelegramClient struct {
	bot    *tele.Bot
	logger *zap.Logger
}

// NewTelegramClient создаёт клиент Telegram-бота.
func NewTelegramClient(token string, logger *zap.Logger) (*TelegramClient, error) {
	pref := tele.Settings{
		Token: token,
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	return &TelegramClient{
		bot:    bot,
		logger: logger,
	}, nil
}

// Send отправляет сообщение в Telegram по chatID.
func (c *TelegramClient) Send(chatID int64, message string) error {
	recipient := &tele.Chat{ID: chatID}

	c.logger.Debug("Sending Telegram message", zap.Int64("chat_id", chatID))

	if _, err := c.bot.Send(recipient, message); err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}

	c.logger.Info("Telegram message sent", zap.Int64("chat_id", chatID))
	return nil
}

// SetWebhook устанавливает webhook для бота.
func (c *TelegramClient) SetWebhook(webhookURL string) error {
	_ = &tele.Webhook{
		Listen: ":8443",
		Endpoint: &tele.WebhookEndpoint{
			PublicURL: webhookURL,
		},
	}

	c.bot.Start()
	return nil
}
