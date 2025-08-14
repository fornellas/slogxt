package log

import (
	"log/slog"
	"strconv"
	"strings"
)

// TerminalValuer is implemented by any value that wants to provide
// a terminal representation for slogxt terminal handlers.
// This allows values to include ANSI escape sequences that will be
// sanitized and displayed in terminal output while still working
// correctly with other handlers (like JSON) which will use the
// standard String() method or LogValue() method.
type TerminalValuer interface {
	// TerminalValue returns a slog.Value that may contain any
	// ANSI escape sequences. Terminal handlers will sanitize
	// these sequences to keep only safe formatting codes.
	TerminalValue() slog.Value
}

// isCSIStart checks if the current position starts a CSI sequence
func isCSIStart(runes []rune, i int) bool {
	return i < len(runes)-1 && runes[i] == '\033' && runes[i+1] == '['
}

// parseANSIParameters parses the numeric parameters of an ANSI sequence
func parseANSIParameters(runes []rune, i *int) []rune {
	var params []rune
	for *i < len(runes) && (runes[*i] >= '0' && runes[*i] <= '9' || runes[*i] == ';') {
		params = append(params, runes[*i])
		*i++
	}
	return params
}

// escapeRune returns the escaped representation of a non-printable rune.
func escapeRune(r rune) string {
	e := strconv.QuoteRune(r)
	return e[1 : len(e)-1] // Remove the surrounding quotes
}

// isValidANSICommand checks if a character is a valid ANSI command character.
func isValidANSICommand(cmd rune) bool {
	// Only allow actual ANSI CSI command characters
	switch cmd {
	case 'A', 'B', 'C', 'D': // Cursor movement
		return true
	case 'E', 'F': // Cursor next/previous line
		return true
	case 'G': // Cursor horizontal absolute
		return true
	case 'H', 'f': // Cursor position
		return true
	case 'J': // Erase in display
		return true
	case 'K': // Erase in line
		return true
	case 'S', 'T': // Scroll up/down
		return true
	case 'm': // SGR (Select Graphic Rendition) - colors and formatting
		return true
	case 'n': // Device status report
		return true
	case 's', 'u': // Save/restore cursor position
		return true
	default:
		return false
	}
}

// isSafeSGRCode checks if a single SGR code is safe
func isSafeSGRCode(param string) bool {
	code, err := strconv.Atoi(param)
	if err != nil {
		return false
	}

	// Allow safe SGR codes:
	// 0: reset, 1-9: formatting, 30-37: fg colors, 40-47: bg colors
	// 90-97: bright fg colors, 100-107: bright bg colors
	return code == 0 ||
		(code >= 1 && code <= 9) ||
		(code >= 30 && code <= 37) ||
		(code >= 40 && code <= 47) ||
		(code >= 90 && code <= 97) ||
		(code >= 100 && code <= 107)
}

// isSafeSGRParameters checks if SGR parameters are safe
func isSafeSGRParameters(params string) bool {
	if params == "" {
		return true // Reset
	}

	for _, param := range strings.Split(params, ";") {
		if param == "" {
			continue
		}
		if !isSafeSGRCode(param) {
			return false
		}
	}
	return true
}

// isSafeANSICommand determines if an ANSI command is safe for terminal output.
func isSafeANSICommand(cmd rune, params string) bool {
	if cmd != 'm' {
		// Block all non-SGR CSI sequences (cursor movement, screen clearing, etc.)
		return false
	}
	return isSafeSGRParameters(params)
}

// processRegularCharacter processes a non-ANSI character and returns the next position
func processRegularCharacter(runes []rune, i int, result *[]rune) int {
	r := runes[i]
	if r == '\t' || strconv.IsPrint(r) {
		*result = append(*result, r)
	} else {
		escaped := escapeRune(r)
		*result = append(*result, []rune(escaped)...)
	}
	return i + 1
}

// processRegularCharacterForStripping processes a non-ANSI character for stripping
func processRegularCharacterForStripping(runes []rune, i int, result *[]rune) int {
	r := runes[i]
	if r == '\t' || r == '\n' || r == '\r' || strconv.IsPrint(r) {
		*result = append(*result, r)
	} else {
		escaped := escapeRune(r)
		*result = append(*result, []rune(escaped)...)
	}
	return i + 1
}

// processANSISequence processes an ANSI escape sequence and returns the next position
func processANSISequence(runes []rune, i int, result *[]rune) int {
	start := i
	i += 2 // Skip '\033['

	// Parse parameters
	params := parseANSIParameters(runes, &i)

	if i >= len(runes) {
		// Incomplete sequence, skip it
		return len(runes)
	}

	cmd := runes[i]
	if isValidANSICommand(cmd) {
		if isSafeANSICommand(cmd, string(params)) {
			// Include the safe ANSI sequence
			for j := start; j <= i; j++ {
				*result = append(*result, runes[j])
			}
		}
		return i + 1
	}

	// Invalid command, process as regular character
	return processRegularCharacter(runes, i, result)
}

// stripANSISequence processes and strips an ANSI escape sequence
func stripANSISequence(runes []rune, i int, result *[]rune) int {
	i += 2 // Skip '\033['

	// Skip parameters
	for i < len(runes) && (runes[i] >= '0' && runes[i] <= '9' || runes[i] == ';') {
		i++
	}

	if i >= len(runes) {
		// Incomplete sequence, skip it
		return len(runes)
	}

	cmd := runes[i]
	if isValidANSICommand(cmd) {
		// Valid ANSI sequence, skip it entirely
		return i + 1
	}

	// Invalid command, process as regular character
	return processRegularCharacterForStripping(runes, i, result)
}

// sanitizeANSI removes dangerous ANSI escape sequences while preserving
// safe color and formatting sequences. It allows:
// - Color changes (30-37, 90-97, 40-47, 100-107)
// - Text formatting (bold, dim, italic, underline, etc.)
// - Reset sequences
// But blocks:
// - Cursor movement
// - Screen clearing
// - Other potentially disruptive sequences
func sanitizeANSI(s string) string {
	result := make([]rune, 0, len(s))
	runes := []rune(s)

	for i := 0; i < len(runes); {
		if isCSIStart(runes, i) {
			i = processANSISequence(runes, i, &result)
		} else {
			i = processRegularCharacter(runes, i, &result)
		}
	}

	return string(result)
}

// stripANSI removes all ANSI escape sequences from a string
func stripANSI(s string) string {
	result := make([]rune, 0, len(s))
	runes := []rune(s)

	for i := 0; i < len(runes); {
		if isCSIStart(runes, i) {
			i = stripANSISequence(runes, i, &result)
		} else {
			i = processRegularCharacterForStripping(runes, i, &result)
		}
	}

	return string(result)
}

// TerminalValue represents a string value that may contain ANSI escape sequences.
// It automatically provides different representations for terminal vs. other handlers:
// - String() and MarshalText() return the text with all ANSI sequences removed
// - TerminalValue() returns the original text with ANSI sequences for terminal handlers
type TerminalValue struct {
	text string
}

// NewTerminalValue creates a new TerminalValue from a string that may contain ANSI sequences
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
