package commands

import (
	"testing"
)

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "no truncation needed",
			input:    "Short title",
			maxLen:   50,
			expected: "Short title",
		},
		{
			name:     "truncate long string",
			input:    "This is a very long title that needs to be truncated",
			maxLen:   30,
			expected: "This is a very long title t...",
		},
		{
			name:     "exact length",
			input:    "Exactly 20 chars!!!",
			maxLen:   19,
			expected: "Exactly 20 chars!!!",
		},
		{
			name:     "very short maxLen",
			input:    "Hello World",
			maxLen:   3,
			expected: "Hel",
		},
		{
			name:     "maxLen of 1",
			input:    "Hello",
			maxLen:   1,
			expected: "H",
		},
		{
			name:     "empty string",
			input:    "",
			maxLen:   10,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateString(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("truncateString(%q, %d) = %q, expected %q",
					tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		seconds  int
		expected string
	}{
		{
			name:     "zero seconds",
			seconds:  0,
			expected: "00:00:00",
		},
		{
			name:     "under one minute",
			seconds:  45,
			expected: "00:00:45",
		},
		{
			name:     "exactly one hour",
			seconds:  3600,
			expected: "01:00:00",
		},
		{
			name:     "complex time",
			seconds:  3665, // 1h 1m 5s
			expected: "01:01:05",
		},
		{
			name:     "long match",
			seconds:  5432, // 1h 30m 32s
			expected: "01:30:32",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.seconds)
			if result != tt.expected {
				t.Errorf("formatDuration(%d) = %q, expected %q",
					tt.seconds, result, tt.expected)
			}
		})
	}
}
