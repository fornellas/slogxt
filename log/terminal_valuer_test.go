package log

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestColoredValue implements TerminalValuer for testing
type TestColoredValue struct {
	plainText    string
	terminalText string
}

func (v TestColoredValue) String() string {
	return v.plainText
}

func (v TestColoredValue) TerminalValue() slog.Value {
	return slog.StringValue(v.terminalText)
}

// TestLogValue implements both LogValuer and TerminalValuer for JSON integration tests
type TestLogValue struct {
	plainText    string
	terminalText string
}

func (v TestLogValue) String() string {
	return v.plainText
}

func (v TestLogValue) LogValue() slog.Value {
	return slog.StringValue(v.plainText)
}

func (v TestLogValue) TerminalValue() slog.Value {
	return slog.StringValue(v.terminalText)
}

func TestStripANSI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "color_codes_stripped",
			input:    "\033[31mred\033[0m \033[32mgreen\033[0m",
			expected: "red green",
		},
		{
			name:     "formatting_codes_stripped",
			input:    "\033[1mbold\033[0m \033[3mitalic\033[0m",
			expected: "bold italic",
		},
		{
			name:     "mixed_ansi_stripped",
			input:    "\033[31mred\033[1mbold\033[0mreset",
			expected: "redboldreset",
		},
		{
			name:     "empty_string",
			input:    "",
			expected: "",
		},
		{
			name:     "no_ansi_sequences",
			input:    "plain text",
			expected: "plain text",
		},
		{
			name:     "reset_only",
			input:    "\033[mreset",
			expected: "reset",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripANSI(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTerminalValue(t *testing.T) {
	t.Run("basic_functionality", func(t *testing.T) {
		tv := NewTerminalValue("\033[31mred text\033[0m")

		// String() should return stripped text
		assert.Equal(t, "red text", tv.String())

		// TerminalValue() should return original text with ANSI
		termVal := tv.TerminalValue()
		assert.Equal(t, "\033[31mred text\033[0m", termVal.String())
	})

	t.Run("with_multiline_content", func(t *testing.T) {
		tv := NewTerminalValue("\033[31m-deleted line\033[0m\n\033[32m+added line\033[0m")

		assert.Equal(t, "-deleted line\n+added line", tv.String())
		assert.Equal(t, "\033[31m-deleted line\033[0m\n\033[32m+added line\033[0m", tv.TerminalValue().String())
	})

	t.Run("with_non_sgr_ansi_sequences", func(t *testing.T) {
		tv := NewTerminalValue("\033[2J\033[31mred\033[H\033[0m")

		// All ANSI sequences should be stripped for String()
		assert.Equal(t, "red", tv.String())

		// Original text preserved for TerminalValue()
		assert.Equal(t, "\033[2J\033[31mred\033[H\033[0m", tv.TerminalValue().String())
	})

	t.Run("empty_text", func(t *testing.T) {
		tv := NewTerminalValue("")

		assert.Equal(t, "", tv.String())
		assert.Equal(t, "", tv.TerminalValue().String())
	})

	t.Run("plain_text_no_ansi", func(t *testing.T) {
		tv := NewTerminalValue("plain text")

		assert.Equal(t, "plain text", tv.String())
		assert.Equal(t, "plain text", tv.TerminalValue().String())
	})

	t.Run("json_handler_uses_logvalue_method", func(t *testing.T) {
		buf := &bytes.Buffer{}
		handler := slog.NewJSONHandler(buf, nil)
		logger := slog.New(handler)

		tv := NewTerminalValue("\033[31mred text\033[0m")

		logger.Info("test message", "value", tv)

		output := buf.String()

		assert.Contains(t, output, "red text")
		assert.NotContains(t, output, "\033[31m")
	})
}
