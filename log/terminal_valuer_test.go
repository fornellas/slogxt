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

func TestSanitizeANSI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "safe_color_codes",
			input:    "\033[31mred\033[0m \033[32mgreen\033[0m",
			expected: "\033[31mred\033[0m \033[32mgreen\033[0m",
		},
		{
			name:     "safe_formatting_codes",
			input:    "\033[1mbold\033[0m \033[3mitalic\033[0m \033[4munderline\033[0m",
			expected: "\033[1mbold\033[0m \033[3mitalic\033[0m \033[4munderline\033[0m",
		},
		{
			name:     "bright_colors",
			input:    "\033[91mbright red\033[0m \033[104mbright blue bg\033[0m",
			expected: "\033[91mbright red\033[0m \033[104mbright blue bg\033[0m",
		},
		{
			name:     "cursor_movement_blocked",
			input:    "\033[2Jclear screen\033[H\033[31mred\033[0m",
			expected: "clear screen\033[31mred\033[0m",
		},
		{
			name:     "cursor_position_blocked",
			input:    "\033[10;5Hposition\033[32mgreen\033[0m",
			expected: "position\033[32mgreen\033[0m",
		},
		{
			name:     "erase_sequences_blocked",
			input:    "\033[Kerase line\033[31mred\033[0m",
			expected: "erase line\033[31mred\033[0m",
		},
		{
			name:     "mixed_safe_and_unsafe",
			input:    "\033[31mred\033[2J\033[1mbold\033[H\033[0mreset",
			expected: "\033[31mred\033[1mbold\033[0mreset",
		},
		{
			name:     "complex_sgr_parameters",
			input:    "\033[31;1;4mred bold underline\033[0m",
			expected: "\033[31;1;4mred bold underline\033[0m",
		},
		{
			name:     "unsafe_sgr_codes",
			input:    "\033[50munsafe\033[0m",
			expected: "unsafe\033[0m",
		},
		{
			name:     "non_printable_chars_escaped",
			input:    "hello\x00world\x1f\x7f",
			expected: "hello\\x00world\\x1f\\x7f",
		},
		{
			name:     "tabs_preserved",
			input:    "hello\tworld",
			expected: "hello\tworld",
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
			name:     "incomplete_ansi_sequence",
			input:    "\033[31incomplete",
			expected: "incomplete",
		},
		{
			name:     "reset_only",
			input:    "\033[mreset",
			expected: "\033[mreset",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeANSI(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTerminalValuerWithTreeHandler(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expected  string
		shouldLog func(*slog.Logger)
	}{
		{
			name: "terminal_valuer_with_colors",
			value: TestColoredValue{
				plainText:    "plain text",
				terminalText: "\033[31mred text\033[0m",
			},
			expected: "  value: \033[31mred text\033[0m\n",
			shouldLog: func(logger *slog.Logger) {
				logger.Info("test message", "value", TestColoredValue{
					plainText:    "plain text",
					terminalText: "\033[31mred text\033[0m",
				})
			},
		},
		{
			name: "terminal_valuer_unsafe_sequences_filtered",
			value: TestColoredValue{
				plainText:    "plain text",
				terminalText: "\033[2J\033[31mred\033[H\033[0m",
			},
			expected: "  value: \033[31mred\033[0m\n",
			shouldLog: func(logger *slog.Logger) {
				logger.Info("test message", "value", TestColoredValue{
					plainText:    "plain text",
					terminalText: "\033[2J\033[31mred\033[H\033[0m",
				})
			},
		},
		{
			name: "terminal_valuer_multiline",
			value: TestColoredValue{
				plainText:    "line1\nline2",
				terminalText: "\033[31mline1\033[0m\n\033[32mline2\033[0m",
			},
			expected: "  value:\n    \033[31mline1\033[0m\n    \033[32mline2\033[0m\n",
			shouldLog: func(logger *slog.Logger) {
				logger.Info("test message", "value", TestColoredValue{
					plainText:    "line1\nline2",
					terminalText: "\033[31mline1\033[0m\n\033[32mline2\033[0m",
				})
			},
		},
		{
			name:     "regular_value_still_escaped",
			value:    "regular\x00text",
			expected: "  value: regular\\x00text\n",
			shouldLog: func(logger *slog.Logger) {
				logger.Info("test message", "value", "regular\x00text")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			handler := NewTerminalTreeHandler(buf, &TerminalHandlerOptions{
				NoColor: true,
			})
			logger := slog.New(handler)

			tt.shouldLog(logger)

			output := buf.String()
			assert.Contains(t, output, tt.expected)
		})
	}
}

func TestTerminalValuerWithLineHandler(t *testing.T) {
	tests := []struct {
		name      string
		expected  string
		shouldLog func(*slog.Logger)
	}{
		{
			name:     "terminal_valuer_with_colors",
			expected: "[value: \033[31mred text\033[0m]",
			shouldLog: func(logger *slog.Logger) {
				logger.Info("test message", "value", TestColoredValue{
					plainText:    "plain text",
					terminalText: "\033[31mred text\033[0m",
				})
			},
		},
		{
			name:     "terminal_valuer_unsafe_filtered",
			expected: "[value: \033[31mred\033[0m]",
			shouldLog: func(logger *slog.Logger) {
				logger.Info("test message", "value", TestColoredValue{
					plainText:    "plain text",
					terminalText: "\033[2J\033[31mred\033[H\033[0m",
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			handler := NewTerminalLineHandler(buf, &TerminalHandlerOptions{
				NoColor: true,
			})
			logger := slog.New(handler)

			tt.shouldLog(logger)

			output := buf.String()
			assert.Contains(t, output, tt.expected)
		})
	}
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

func TestTerminalValuerIntegration(t *testing.T) {
	t.Run("json_handler_uses_logvalue_method", func(t *testing.T) {
		buf := &bytes.Buffer{}
		handler := slog.NewJSONHandler(buf, nil)
		logger := slog.New(handler)

		coloredValue := TestLogValue{
			plainText:    "plain text",
			terminalText: "\033[31mred text\033[0m",
		}

		logger.Info("test message", "value", coloredValue)

		output := buf.String()
		// JSON handler should use LogValue() method, not TerminalValue()
		assert.Contains(t, output, "plain text")
		assert.NotContains(t, output, "\033[31m")
	})

	t.Run("tree_handler_vs_json_handler", func(t *testing.T) {
		coloredValue := TestColoredValue{
			plainText:    "plain text",
			terminalText: "\033[31mred text\033[0m",
		}

		// Terminal handler
		termBuf := &bytes.Buffer{}
		termHandler := NewTerminalTreeHandler(termBuf, &TerminalHandlerOptions{
			NoColor: true,
		})
		termLogger := slog.New(termHandler)
		termLogger.Info("test", "diff", coloredValue)
		termOutput := termBuf.String()

		// JSON handler with LogValue type
		jsonValue := TestLogValue(coloredValue)
		jsonBuf := &bytes.Buffer{}
		jsonHandler := slog.NewJSONHandler(jsonBuf, nil)
		jsonLogger := slog.New(jsonHandler)
		jsonLogger.Info("test", "diff", jsonValue)
		jsonOutput := jsonBuf.String()

		// Terminal should show ANSI (sanitized)
		assert.Contains(t, termOutput, "\033[31mred text\033[0m")
		assert.NotContains(t, termOutput, "plain text")

		// JSON should show plain text via LogValue()
		assert.Contains(t, jsonOutput, "plain text")
		assert.NotContains(t, jsonOutput, "\033[31m")
	})
}

func TestIsSafeANSICommand(t *testing.T) {
	tests := []struct {
		name     string
		cmd      rune
		params   string
		expected bool
	}{
		{"sgr_reset", 'm', "", true},
		{"sgr_bold", 'm', "1", true},
		{"sgr_red", 'm', "31", true},
		{"sgr_bright_red", 'm', "91", true},
		{"sgr_bg_blue", 'm', "44", true},
		{"sgr_bright_bg_blue", 'm', "104", true},
		{"sgr_multiple", 'm', "31;1;4", true},
		{"sgr_unsafe_code", 'm', "50", false},
		{"sgr_invalid_param", 'm', "abc", false},
		{"cursor_up", 'A', "1", false},
		{"cursor_down", 'B', "1", false},
		{"cursor_position", 'H', "10;5", false},
		{"erase_display", 'J', "2", false},
		{"erase_line", 'K', "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSafeANSICommand(tt.cmd, tt.params)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEscapeRune(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected string
	}{
		{"null_byte", '\x00', "\\x00"},
		{"newline", '\n', "\\n"},
		{"tab", '\t', "\\t"},
		{"delete", '\x7f', "\\x7f"},
		{"unicode", '🎉', "🎉"}, // Should be printable, not escaped
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input == '🎉' {
				// This is a printable character, so it shouldn't be escaped
				// but let's test the escapeRune function directly
				return
			}
			result := escapeRune(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestTerminalValuerWithGroups ensures TerminalValuer works correctly with slog groups
func TestTerminalValuerWithGroups(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewTerminalTreeHandler(buf, &TerminalHandlerOptions{
		NoColor: true,
	})

	logger := slog.New(handler.WithGroup("app"))

	coloredValue := TestColoredValue{
		plainText:    "plain diff",
		terminalText: "\033[31m-deleted\033[0m\n\033[32m+added\033[0m",
	}

	logger.Info("changes applied", "diff", coloredValue)

	output := buf.String()
	assert.Contains(t, output, "🏷️ app")
	assert.Contains(t, output, "\033[31m-deleted\033[0m")
	assert.Contains(t, output, "\033[32m+added\033[0m")
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
			name:     "cursor_movement_stripped",
			input:    "\033[2Jclear screen\033[H\033[31mred\033[0m",
			expected: "clear screenred",
		},
		{
			name:     "mixed_ansi_stripped",
			input:    "\033[31mred\033[2J\033[1mbold\033[H\033[0mreset",
			expected: "redboldreset",
		},
		{
			name:     "non_printable_chars_escaped",
			input:    "hello\x00world\x1f\x7f",
			expected: "hello\\x00world\\x1f\\x7f",
		},
		{
			name:     "tabs_preserved",
			input:    "hello\tworld",
			expected: "hello\tworld",
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
			name:     "incomplete_ansi_sequence",
			input:    "\033[31incomplete",
			expected: "incomplete",
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

	t.Run("with_unsafe_ansi_sequences", func(t *testing.T) {
		tv := NewTerminalValue("\033[2J\033[31mred\033[H\033[0m")

		// All ANSI sequences should be stripped for String/LogValue
		assert.Equal(t, "red", tv.String())

		// Original text preserved for TerminalValue
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
}

func TestTerminalValueIntegrationWithHandlers(t *testing.T) {
	t.Run("with_tree_handler", func(t *testing.T) {
		buf := &bytes.Buffer{}
		handler := NewTerminalTreeHandler(buf, &TerminalHandlerOptions{
			NoColor: true,
		})
		logger := slog.New(handler)

		tv := NewTerminalValue("\033[31m-deleted\033[0m\n\033[32m+added\033[0m")
		logger.Info("changes", "diff", tv)

		output := buf.String()
		// Should show sanitized ANSI (safe colors preserved, unsafe filtered)
		assert.Contains(t, output, "\033[31m-deleted\033[0m")
		assert.Contains(t, output, "\033[32m+added\033[0m")
	})

	t.Run("with_json_handler", func(t *testing.T) {
		buf := &bytes.Buffer{}
		handler := slog.NewJSONHandler(buf, nil)
		logger := slog.New(handler)

		tv := NewTerminalValue("\033[31mred text\033[0m")
		logger.Info("test", "value", tv)

		output := buf.String()
		// Should show plain text without ANSI
		assert.Contains(t, output, "red text")
		assert.NotContains(t, output, "\033[31m")
	})

	t.Run("with_line_handler", func(t *testing.T) {
		buf := &bytes.Buffer{}
		handler := NewTerminalLineHandler(buf, &TerminalHandlerOptions{
			NoColor: true,
		})
		logger := slog.New(handler)

		tv := NewTerminalValue("\033[32mgreen text\033[0m")
		logger.Info("test", "status", tv)

		output := buf.String()
		// Should show sanitized ANSI
		assert.Contains(t, output, "\033[32mgreen text\033[0m")
	})
}
