package main

import (
	"log/slog"
	"os"

	"github.com/fornellas/slogxpert"
)

func main() {
	// Create a file for logging
	logFile, err := os.Create("application.log")
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	// Create different handlers for different outputs
	// 1. Terminal handler for human-readable console output
	consoleHandler := slogxpert.NewTerminalTreeHandler(os.Stdout, &slogxpert.TerminalHandlerOptions{
		HandlerOptions: slog.HandlerOptions{
			Level: slog.LevelInfo, // Only INFO and above for console
		},
		TimeLayout: "15:04:05",
	})

	// 2. Line handler for file output
	fileHandler := slogxpert.NewTerminalLineHandler(logFile, &slogxpert.TerminalHandlerOptions{
		HandlerOptions: slog.HandlerOptions{
			Level: slog.LevelDebug, // All logs for file
		},
		TimeLayout: "2006-01-02 15:04:05",
		NoColor:    true, // No color codes in file
	})

	// Combine all handlers into a MultiHandler
	multiHandler := slogxpert.NewMultiHandler(
		consoleHandler,
		fileHandler,
	)

	// Create a logger with the multi handler
	logger := slog.New(multiHandler)

	// All log entries will go to all three handlers according to their level settings
	logger.Debug("Debug message")               // Only goes to file (level too low for console)
	logger.Info("Info message")                 // Goes to console and file
	logger.Warn("Warning message")              // Goes to all handlers
	logger.Error("Error occurred", "code", 500) // Goes to all handlers
}
