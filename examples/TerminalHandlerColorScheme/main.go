package main

import (
	"log/slog"
	"os"

	"github.com/fornellas/slogxpert/ansi"
	"github.com/fornellas/slogxpert/log"
)

func main() {
	customColors := &log.TerminalHandlerColorScheme{
		GroupName:    ansi.SGRs{ansi.FgMagenta, ansi.Bold},
		AttrKey:      ansi.SGRs{ansi.FgBlue},
		AttrValue:    ansi.SGRs{ansi.FgWhite},
		Time:         ansi.SGRs{ansi.FgYellow},
		LevelDebug:   ansi.SGRs{ansi.FgCyan},
		MessageDebug: ansi.SGRs{ansi.FgCyan},
		LevelInfo:    ansi.SGRs{ansi.FgGreen, ansi.Bold},
		MessageInfo:  ansi.SGRs{ansi.FgWhite, ansi.Bold},
		LevelWarn:    ansi.SGRs{ansi.FgYellow, ansi.Bold},
		MessageWarn:  ansi.SGRs{ansi.FgYellow},
		LevelError:   ansi.SGRs{ansi.FgRed, ansi.Bold, ansi.Blink},
		MessageError: ansi.SGRs{ansi.FgRed, ansi.Bold},
		File:         ansi.SGRs{ansi.FgBlue, ansi.Italic},
		Line:         ansi.SGRs{ansi.FgBlue, ansi.Bold},
		Function:     ansi.SGRs{ansi.FgBlue, ansi.Dim},
	}

	handler := log.NewTerminalTreeHandler(os.Stdout, &log.TerminalHandlerOptions{
		ColorScheme: customColors,
	})

	logger := slog.New(handler)

	logger.Info("Application started")

	logger.Info("User logged in", "user_id", 123, "username", "john_doe")

	userLogger := logger.WithGroup("user")
	userLogger.Info("Profile updated", "changes", 2)

	settingsLogger := userLogger.WithGroup("settings")
	settingsLogger.Info("Theme changed", "old", "light", "new", "dark")
}
