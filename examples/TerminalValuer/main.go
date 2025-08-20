package main

import (
	"bytes"
	"log/slog"
	"os"

	"github.com/fornellas/slogxt/log"
)

func main() {
	slog.Info("=== TerminalValuer Demo ===")

	// Create sample colored values using TerminalValue
	diff := log.NewTerminalValue("\033[31m-old line\033[0m\n\033[32m+new line\033[0m")
	status := log.NewTerminalValue("\033[32m✓\033[0m Operation completed successfully")

	// Demo 1: TerminalTreeHandler shows ANSI colors
	slog.Info("\n1. Terminal Tree Handler (with colors):")
	termHandler := log.NewTerminalTreeHandler(os.Stdout, &log.TerminalHandlerOptions{
		ForceColor: true, // Force colors even if not a TTY
	})
	termLogger := slog.New(termHandler)

	termLogger.Info("File changes applied", "diff", diff)
	termLogger.Info("Status update", "status", status)

	// Demo 2: JSON Handler shows plain text
	slog.Info("\n2. JSON Handler (plain text):")
	jsonBuf := &bytes.Buffer{}
	jsonHandler := slog.NewJSONHandler(jsonBuf, nil)
	jsonLogger := slog.New(jsonHandler)

	jsonLogger.Info("File changes applied", "diff", diff)
	jsonLogger.Info("Status update", "status", status)

	slog.Info("JSON Output:")
	os.Stdout.WriteString(jsonBuf.String())

	// Demo 3: TerminalLineHandler
	slog.Info("3. Terminal Line Handler (with colors):")
	lineHandler := log.NewTerminalLineHandler(os.Stdout, &log.TerminalHandlerOptions{
		ForceColor: true,
	})
	lineLogger := slog.New(lineHandler)

	lineLogger.Info("File changes applied", "diff", diff)
	lineLogger.Info("Status update", "status", status)

	slog.Info("\n=== End Demo ===")
}
