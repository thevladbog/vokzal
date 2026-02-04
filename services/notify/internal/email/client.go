// Package email предоставляет клиент для отправки email по SMTP.
package email

import (
	"bytes"
	"fmt"
	"html"
	"mime"
	"mime/quotedprintable"
	"net/mail"
	"net/smtp"
	"strings"

	"go.uber.org/zap"
)

// EmailClient — клиент для отправки email.
//
//nolint:revive // Имя сохраняем для ясности при использовании из других пакетов (email.Client).
type EmailClient struct {
	logger   *zap.Logger
	smtpHost string
	username string
	password string
	from     string
	smtpPort int
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

// sanitizeHeader удаляет CRLF символы из заголовков для предотвращения header injection.
func sanitizeHeader(input string) string {
	// Удаляем все символы перевода строки и возврата каретки
	sanitized := strings.ReplaceAll(input, "\r", "")
	sanitized = strings.ReplaceAll(sanitized, "\n", "")
	// Также удаляем нулевые байты
	sanitized = strings.ReplaceAll(sanitized, "\x00", "")
	return strings.TrimSpace(sanitized)
}

// validateEmail проверяет корректность email адреса.
func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

// Send отправляет email по указанному адресу с защитой от injection атак.
func (c *EmailClient) Send(to, subject, body string) error {
	// Валидация email адреса получателя
	if err := validateEmail(to); err != nil {
		return fmt.Errorf("invalid recipient email address: %w", err)
	}

	// Валидация email адреса отправителя
	if err := validateEmail(c.from); err != nil {
		return fmt.Errorf("invalid sender email address: %w", err)
	}

	// Санитизация заголовков для предотвращения header injection
	sanitizedTo := sanitizeHeader(to)
	sanitizedSubject := sanitizeHeader(subject)

	// Безопасное кодирование темы письма с помощью MIME Q-encoding
	encodedSubject := mime.QEncoding.Encode("UTF-8", sanitizedSubject)

	// HTML-экранирование тела письма для предотвращения XSS
	// Пользовательский ввод рассматривается как текст и безопасно встраивается в HTML-шаблон
	escapedBody := html.EscapeString(body)

	// Формируем простой HTML-шаблон, в который вставляем только экранированный текст
	htmlBody := "<html><body><pre style=\"white-space:pre-wrap;\">" + escapedBody + "</pre></body></html>"

	// Сборка MIME-сообщения с использованием стандартных средств кодирования
	var msgBuf bytes.Buffer

	// Заголовки письма
	headers := map[string]string{
		"From":         c.from,
		"To":           sanitizedTo,
		"Subject":      encodedSubject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=UTF-8",
	}

	for k, v := range headers {
		// sanitizeHeader уже удалил CRLF, что предотвращает header injection
		fmt.Fprintf(&msgBuf, "%s: %s\r\n", k, v)
	}

	// Пустая строка отделяет заголовки от тела
	msgBuf.WriteString("\r\n")

	// Кодируем тело в quoted-printable для безопасной передачи
	qpWriter := quotedprintable.NewWriter(&msgBuf)
	if _, err := qpWriter.Write([]byte(htmlBody)); err != nil {
		return fmt.Errorf("failed to encode email body: %w", err)
	}
	if err := qpWriter.Close(); err != nil {
		return fmt.Errorf("failed to finish encoding email body: %w", err)
	}

	auth := smtp.PlainAuth("", c.username, c.password, c.smtpHost)
	addr := fmt.Sprintf("%s:%d", c.smtpHost, c.smtpPort)

	c.logger.Debug("Sending email",
		zap.String("to", sanitizedTo),
		zap.String("subject", sanitizedSubject))

	if err := smtp.SendMail(addr, auth, c.from, []string{sanitizedTo}, msgBuf.Bytes()); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	c.logger.Info("Email sent successfully", zap.String("to", sanitizedTo))
	return nil
}
