package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/fornellas/slogxt/log"
)

func main() {
	handler := log.NewTerminalLineHandler(os.Stderr, &log.TerminalHandlerOptions{
		HandlerOptions: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	})

	logger := slog.New(handler)

	ctx := context.Background()

	// Set given log to context
	ctx = log.WithLogger(ctx, logger)

	if !authenticateUser(ctx, "john", "bad password") {
		logger.Error("Authentication failure")
	}
}

func authenticateUser(ctx context.Context, username, password string) bool {
	ctx, logger := log.MustWithGroupAttrs(ctx, "login", "username", username)

	logger.Info("User authentication")
	return checkPassword(ctx, username, password)
}

func checkPassword(ctx context.Context, username, password string) bool {
	logger := log.MustLogger(ctx)
	logger.Debug("Checking user and password")
	return username == "john" && password == "secret"
}
