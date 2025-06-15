package main

import (
	"log/slog"
	"os"

	"github.com/fornellas/slogxpert/log"
)

func main() {
	baseHandler := log.NewTerminalTreeHandler(os.Stdout, &log.TerminalHandlerOptions{})

	// Create a buffered handler that wraps the base handler
	bufferedHandler := log.NewBufferedHandler(baseHandler)

	logger := slog.New(bufferedHandler)

	// These logs will be stored in the buffer but not displayed yet
	logger.Info("Starting operation")
	logger.Debug("Processing item 1")
	logger.Debug("Processing item 2")

	// Use WithGroup/WithAttrs as normal - all derived loggers share the same buffer
	txLogger := logger.WithGroup("transaction")
	txLogger.Info("Transaction started", "tx_id", "abc123")

	// Some error occurred
	if true { // simulating an error condition
		txLogger.Error("Transaction failed", "reason", "timeout")
	}

	// Now flush all the buffered logs at once
	bufferedHandler.Flush()

	// New logs will again be buffered until flushed
	logger.Info("Operation complete")
	// This won't be displayed until Flush() is called again
}
