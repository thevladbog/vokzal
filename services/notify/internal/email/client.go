// Package email provides a client for sending emails via SMTP.
package email

import (
	"bytes"
	"fmt"
	"mime"
	"mime/quotedprintable"
	"net/mail"
	"net/smtp"
	"strings"

	"go.uber.org/zap"
)

// EmailClient is a client for sending emails.
//
//nolint:revive // Name preserved for clarity when used from other packages (email.Client).
type EmailClient struct {
	logger   *zap.Logger
	smtpHost string
	username string
	password string
	from     string
	smtpPort int
}

// NewEmailClient creates a client for sending emails.
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

// sanitizeHeader removes CRLF characters from headers to prevent header injection.
func sanitizeHeader(input string) string {
	// Remove all newline and carriage return characters
	sanitized := strings.ReplaceAll(input, "\r", "")
	sanitized = strings.ReplaceAll(sanitized, "\n", "")
	// Also remove null bytes
	sanitized = strings.ReplaceAll(sanitized, "\x00", "")
	return strings.TrimSpace(sanitized)
}

// sanitizeBody removes dangerous control characters from the email body while
// preserving legitimate line breaks. This allows multi-line email bodies.
func sanitizeBody(input string) string {
	// Map over runes and drop dangerous ASCII control characters (0x00–0x1F)
	// but explicitly allow newline (\n), carriage return (\r), and tab (\t).
	sanitized := strings.Map(func(r rune) rune {
		// Allow newline, carriage return, and tab
		if r == '\n' || r == '\r' || r == '\t' {
			return r
		}
		// Drop other control characters (especially null bytes and other dangerous chars)
		if r < 0x20 {
			return -1
		}
		return r
	}, input)
	// Don't use TrimSpace as it would remove legitimate trailing newlines
	return sanitized
}

// validateEmail validates email address format.
func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

// Send sends an email to the specified address with protection against injection attacks.
func (c *EmailClient) Send(to, subject, body string) error {
	// Validate recipient email address
	toAddr, err := mail.ParseAddress(to)
	if err != nil {
		return fmt.Errorf("invalid recipient email address: %w", err)
	}

	// Validate sender email address
	fromAddr, err := mail.ParseAddress(c.from)
	if err != nil {
		return fmt.Errorf("invalid sender email address: %w", err)
	}

	// Extract bare mailbox addresses for SMTP envelope
	// smtp.SendMail requires bare addresses (e.g., "user@example.com")
	// not display-name format (e.g., "John Doe <user@example.com>")
	bareToAddr := toAddr.Address
	bareFromAddr := fromAddr.Address

	// Sanitize headers to prevent header injection
	sanitizedTo := sanitizeHeader(to)
	sanitizedSubject := sanitizeHeader(subject)

	// Safely encode subject using MIME Q-encoding
	encodedSubject := mime.QEncoding.Encode("UTF-8", sanitizedSubject)

	// Sanitize body content by removing any potential control characters
	// Use text/plain instead of HTML to eliminate XSS risks entirely
	sanitizedBody := sanitizeBody(body)

	// Build MIME message using standard encoding methods
	var msgBuf bytes.Buffer

	// Email headers
	headers := map[string]string{
		"From":                      c.from,
		"To":                        sanitizedTo,
		"Subject":                   encodedSubject,
		"MIME-Version":              "1.0",
		"Content-Type":              "text/plain; charset=UTF-8",
		"Content-Transfer-Encoding": "quoted-printable",
	}

	for k, v := range headers {
		// sanitizeHeader already removed CRLF, preventing header injection
		fmt.Fprintf(&msgBuf, "%s: %s\r\n", k, v)
	}

	// Empty line separates headers from body
	msgBuf.WriteString("\r\n")

	// Encode body in quoted-printable for safe transmission
	// Use only sanitized text without HTML tags
	qpWriter := quotedprintable.NewWriter(&msgBuf)
	if _, err := qpWriter.Write([]byte(sanitizedBody)); err != nil {
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

	// Use bare addresses for SMTP envelope (required by smtp.SendMail)
	if err := smtp.SendMail(addr, auth, bareFromAddr, []string{bareToAddr}, msgBuf.Bytes()); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	c.logger.Info("Email sent successfully", zap.String("to", sanitizedTo))
	return nil
}
