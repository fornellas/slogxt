package main

import (
	"log/slog"
	"os"

	"github.com/fornellas/slogxpert"
)

func main() {
	handler := slogxpert.NewTerminalTreeHandler(os.Stderr, &slogxpert.TerminalHandlerOptions{
		HandlerOptions: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	})

	logger := slog.New(handler)

	logger.Info("Application started")

	logger.Info("User logged in", "user_id", 123, "username", "john_doe")

	userLogger := logger.WithGroup("user")
	userLogger.Info("Profile updated", "changes", 2)

	settingsLogger := userLogger.WithGroup("settings")
	settingsLogger.Info("Theme changed", "old", "light", "new", "dark")
}
