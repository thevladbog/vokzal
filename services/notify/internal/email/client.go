// Package email предоставляет клиент для отправки email по SMTP.
package email

import (
	"fmt"
	"net/smtp"

	"go.uber.org/zap"
)

// EmailClient — клиент для отправки email.
//
//nolint:revive // Имя сохраняем для ясности при использовании из других пакетов (email.Client).
type EmailClient struct {
	smtpHost string
	smtpPort int
	username string
	password string
	from     string
	logger   *zap.Logger
}

// NewEmailClient создаёт клиент для отправки email.
func NewEmailClient(smtpHost string, smtpPort int, username, password, from string, logger *zap.Logger) *EmailClient {
	return &EmailClient{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		username: username,
		password: password,
		from:     from,
		logger:   logger,
	}
}

// Send отправляет email по указанному адресу.
func (c *EmailClient) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", c.username, c.password, c.smtpHost)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", c.from, to, subject, body)

	addr := fmt.Sprintf("%s:%d", c.smtpHost, c.smtpPort)

	c.logger.Debug("Sending email", zap.String("to", to), zap.String("subject", subject))

	if err := smtp.SendMail(addr, auth, c.from, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	c.logger.Info("Email sent successfully", zap.String("to", to))
	return nil
}
