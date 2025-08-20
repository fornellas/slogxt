package log

import (
	"log/slog"
	"regexp"
)

// TerminalValuer is implemented by any value that wants to provide
// a terminal representation for slogxt terminal handlers.
//
// IMPORTANT: Only use Select Graphic Rendition (SGR) escape sequences
// for colors and text formatting. Other ANSI control sequences (cursor
// movement, screen clearing, etc.) will disrupt terminal output.
type TerminalValuer interface {
	// TerminalValue returns a slog.Value that should only contain
	// Select Graphic Rendition (SGR) ANSI escape sequences.
	// Other control sequences will disrupt terminal display.
	TerminalValue() slog.Value
}

var ansiRegex = regexp.MustCompile(`\033\[[0-9;]*[A-Za-z]`)

// stripANSI removes all ANSI escape sequences from a string
func stripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

// TerminalValue represents a string value that may contain ANSI escape sequences.
// It automatically provides different representations for terminal vs. other handlers:
// - String() and MarshalText() return the text with all ANSI sequences removed
// - TerminalValue() returns the original text with ANSI sequences for terminal handlers
type TerminalValue struct {
	text string
}

// NewTerminalValue creates a new TerminalValue from a string that may contain
// Select Graphic Rendition (SGR) ANSI sequences.
//
// IMPORTANT: Only use Select Graphic Rendition (SGR) escape sequences
// for colors and text formatting. Other ANSI control sequences (cursor
// movement, screen clearing, etc.) will disrupt terminal output.
func NewTerminalValue(text string) TerminalValue {
	return TerminalValue{text: text}
}

// String returns the text with all ANSI escape sequences removed
func (tv TerminalValue) String() string {
	return stripANSI(tv.text)
}

// MarshalText implements encoding.TextMarshaler for JSON compatibility
func (tv TerminalValue) MarshalText() ([]byte, error) {
	return []byte(tv.String()), nil
}

// TerminalValue implements TerminalValuer, returning the original text with ANSI sequences
func (tv TerminalValue) TerminalValue() slog.Value {
	return slog.StringValue(tv.text)
}
