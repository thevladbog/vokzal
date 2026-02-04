package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeHeader(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "clean input",
			input:    "Normal Subject",
			expected: "Normal Subject",
		},
		{
			name:     "CRLF injection attempt",
			input:    "Subject\r\nBcc: attacker@evil.com",
			expected: "SubjectBcc: attacker@evil.com",
		},
		{
			name:     "newline injection",
			input:    "Subject\nBcc: attacker@evil.com",
			expected: "SubjectBcc: attacker@evil.com",
		},
		{
			name:     "null byte injection",
			input:    "Subject\x00Bcc: attacker@evil.com",
			expected: "SubjectBcc: attacker@evil.com",
		},
		{
			name:     "multiple CRLF",
			input:    "Subject\r\n\r\nMalicious content",
			expected: "SubjectMalicious content",
		},
		{
			name:     "leading/trailing whitespace",
			input:    "  Subject with spaces  ",
			expected: "Subject with spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeHeader(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "valid email",
			email:   "user@example.com",
			wantErr: false,
		},
		{
			name:    "valid email with name",
			email:   "John Doe <john@example.com>",
			wantErr: false,
		},
		{
			name:    "invalid email - no domain",
			email:   "userexample.com",
			wantErr: true,
		},
		{
			name:    "invalid email - no @",
			email:   "user.example.com",
			wantErr: true,
		},
		{
			name:    "empty email",
			email:   "",
			wantErr: true,
		},
		{
			name:    "email with injection attempt",
			email:   "user@example.com\r\nBcc: attacker@evil.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
