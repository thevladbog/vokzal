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

func TestSanitizeBody(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "clean single line",
			input:    "Hello, world!",
			expected: "Hello, world!",
		},
		{
			name:     "multi-line with newlines",
			input:    "Line 1\nLine 2\nLine 3",
			expected: "Line 1\nLine 2\nLine 3",
		},
		{
			name:     "multi-line with CRLF",
			input:    "Line 1\r\nLine 2\r\nLine 3",
			expected: "Line 1\nLine 2\nLine 3", // CRLF normalized to LF
		},
		{
			name:     "text with tabs",
			input:    "Column1\tColumn2\tColumn3",
			expected: "Column1\tColumn2\tColumn3",
		},
		{
			name:     "null byte removal",
			input:    "Text\x00WithNull",
			expected: "TextWithNull",
		},
		{
			name:     "bell character removal",
			input:    "Text\x07WithBell",
			expected: "TextWithBell",
		},
		{
			name:     "multiple control chars removed but newlines preserved",
			input:    "Line 1\x00\x01\x02\nLine 2\x07\x08\r\nLine 3",
			expected: "Line 1\nLine 2\nLine 3", // Control chars removed, CRLF normalized
		},
		{
			name:     "mixed whitespace preserved",
			input:    "  Line with spaces  \n\tTabbed line\n",
			expected: "  Line with spaces  \n\tTabbed line\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeBody(tt.input)
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
