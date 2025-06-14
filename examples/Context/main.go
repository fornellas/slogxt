package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/fornellas/slogxpert"
)

func main() {
	handler := slogxpert.NewTerminalLineHandler(os.Stderr, &slogxpert.TerminalHandlerOptions{
		HandlerOptions: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	})

	logger := slog.New(handler)

	ctx := context.Background()

	// Set given log to context
	ctx = slogxpert.WithLogger(ctx, logger)

	if !authenticateUser(ctx, "john", "bad password") {
		logger.Error("Authentication failure")
	}
}

func authenticateUser(ctx context.Context, username, password string) bool {
	logger := slogxpert.MustLogger(ctx)

	ctx, logger = slogxpert.MustWithGroupAttrs(ctx, "login", "username", username)

	logger.Info("User authentication")
	return checkPassword(ctx, username, password)
}

func checkPassword(ctx context.Context, username, password string) bool {
	logger := slogxpert.MustLogger(ctx)
	logger.Debug("Checking user and password")
	return username == "john" && "secret" == password
}
