package main

import (
	"log/slog"
	"os"

	"github.com/fornellas/slogxpert/log"
)

func main() {
	handler := log.NewTerminalTreeHandler(os.Stdout, &log.TerminalHandlerOptions{
		HandlerOptions: slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: true,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == "password" {
					return slog.String("password", "******")
				}
				return a
			},
		},
		TimeLayout: "2006-01-02 15:04:05",
		ForceColor: true,
	})

	logger := slog.New(handler)

	logger.Info("Application started")

	logger.Info("User logged in", "user_id", 123, "username", "john_doe", "password", "super_secret")

	userLogger := logger.WithGroup("user")
	userLogger.Info("Profile updated", "changes", 2)

	settingsLogger := userLogger.WithGroup("settings")
	settingsLogger.Info("Theme changed", "old", "light", "new", "dark")
}
