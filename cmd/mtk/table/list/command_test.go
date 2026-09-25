package list

import (
	"testing"
)

func TestShellEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple table name",
			input:    "users",
			expected: "users",
		},
		{
			name:     "table name with single quote",
			input:    "user's",
			expected: "user'\\''s",
		},
		{
			name:     "malicious table name with command injection",
			input:    "safe'; touch /tmp/marker; echo '",
			expected: "safe'\\\''; touch /tmp/marker; echo '\\''",
		},
		{
			name:     "table name with multiple single quotes",
			input:    "it's'a'test",
			expected: "it'\\''s'\\''a'\\''test",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "table name with spaces",
			input:    "my table",
			expected: "my table",
		},
		{
			name:     "table name with special chars but no quotes",
			input:    "table-name_123",
			expected: "table-name_123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shellEscape(tt.input)
			if result != tt.expected {
				t.Errorf("shellEscape(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
