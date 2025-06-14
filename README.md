[![ci](https://github.com/fornellas/slogxpert/actions/workflows/ci.yaml/badge.svg)](https://github.com/fornellas/slogxpert/actions/workflows/ci.yaml) [![update_deps](https://github.com/fornellas/slogxpert/actions/workflows/update_deps.yaml/badge.svg)](https://github.com/fornellas/slogxpert/actions/workflows/update_deps.yaml) [![Go Report Card](https://goreportcard.com/badge/github.com/fornellas/slogxpert)](https://goreportcard.com/report/github.com/fornellas/slogxpert) [![Coverage Status](https://coveralls.io/repos/github/fornellas/slogxpert/badge.svg?branch=master)](https://coveralls.io/github/fornellas/slogxpert?branch=master) [![Go Reference](https://pkg.go.dev/badge/github.com/fornellas/slogxpert.svg)](https://pkg.go.dev/github.com/fornellas/slogxpert) [![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0) [![Buy me a beer: donate](https://img.shields.io/badge/Donate-Buy%20me%20a%20beer-yellow)](https://www.paypal.com/donate?hosted_button_id=AX26JVRT2GS2Q)

# slogxpert

This package extends Go's `log/slog`, providing handlers for structured colored console output, buffering and other goodies.

The full detailed documentation can be found [here](https://pkg.go.dev/github.com/fornellas/slogxpert), here are some highlights.

Full executable examples can be found [here](https://github.com/fornellas/slogxpert/tree/main/examples).

Contributions are accepted, check [README.development.md](README.development.md) for instructions on how to run the build, and cut a PR with your changes.

## Handlers

### Terminal

#### TerminalTreeHandler

The `TerminalTreeHandler` handler bridges the gap between structured logging and human friendly console logging. It organizes log entries with proper indentation for groups, colorized level indicators, and formatted attributes.

This is particularly useful in somme scenarios:

- CLI applications that with to express hyerarchical context to users.
- Local development of server applications that use structured logging (eg: JSON), which are easier to read for humans with `TerminalTreeHandler`.

```go
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
```

Output will look like:

```
INFO Application started
INFO User logged in
  user_id: 123
  username: john_doe
🏷 user
  INFO Profile updated
    changes: 2
  🏷 settings
    INFO Theme changed
      old: light
      new: dark
```

#### TerminalLineHandler

The `TerminalLineHandler` is similar to `TerminalTreeHandler`, but instead of nesting attributes, it always outputs one line per log message.

```go
handler := slogxpert.NewTerminalLineHandler(os.Stderr, &slogxpert.TerminalHandlerOptions{
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
```

Output will look like:

```
INFO Application started
INFO User logged in [user_id: 123, username: john_doe]
INFO 🏷 user: Profile updated [changes: 2]
INFO 🏷 user > 🏷 settings: Theme changed [old: light, new: dark]
```

#### Customizing

Handlers can be customized to your taste.

##### TerminalHandlerColorScheme

The full color scheme can be customized:

```go
customColors := &slogxpert.TerminalHandlerColorScheme{
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

handler := slogxpert.NewTerminalTreeHandler(os.Stdout, &slogxpert.TerminalHandlerOptions{
	ColorScheme: customColors,
})

logger := slog.New(handler)

logger.Info("Application started")

logger.Info("User logged in", "user_id", 123, "username", "john_doe")

userLogger := logger.WithGroup("user")
userLogger.Info("Profile updated", "changes", 2)

settingsLogger := userLogger.WithGroup("settings")
settingsLogger.Info("Theme changed", "old", "light", "new", "dark")
```

##### TerminalHandlerOptions

`TerminalHandlerOptions` provides extensive customization options for both terminal handlers:

```go
handler := slogxpert.NewTerminalTreeHandler(os.Stdout, &slogxpert.TerminalHandlerOptions{
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
```

Output will look like:

```
INFO Application started
  2025-06-14 13:45:32
  /home/fornellas/src/slogxpert/examples/TerminalHandlerOptions/main.go:28 (main.main)
INFO User logged in
  2025-06-14 13:45:32
  /home/fornellas/src/slogxpert/examples/TerminalHandlerOptions/main.go:30 (main.main)
  user_id: 123
  username: john_doe
  password: ******
🏷 user
  INFO Profile updated
    2025-06-14 13:45:32
    /home/fornellas/src/slogxpert/examples/TerminalHandlerOptions/main.go:33 (main.main)
    changes: 2
  🏷 settings
    INFO Theme changed
      2025-06-14 13:45:32
      /home/fornellas/src/slogxpert/examples/TerminalHandlerOptions/main.go:36 (main.main)
      old: light
      new: dark
```

### BufferedHandler

The `BufferedHandler` allows you to buffer log records in memory until you explicitly flush them to the underlying handler. This is particularly useful when concurrent tasks generate logs, but you with that logging to be clustered per task, instead of interleaved.

```go
baseHandler := slogxpert.NewTerminalTreeHandler(os.Stdout, &slogxpert.TerminalHandlerOptions{})

// Create a buffered handler that wraps the base handler
bufferedHandler := slogxpert.NewBufferedHandler(baseHandler)

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
```

When `Flush()` is called, the output will show all the buffered logs at once:

```
INFO Starting operation
🏷 transaction
  INFO Transaction started
    tx_id: abc123
  ERROR Transaction failed
    reason: timeout
```

The final "Operation complete" log won't appear until `Flush()` is called again.

### MultiHandler

The `MultiHandler` dispatches log records to multiple handlers simultaneously. This is useful when you want to send logs to different destinations or format them differently for various purposes (e.g., console output, file logging, structured JSON for analysis).

```go
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
```

The console output will be:

```
INFO Info message
  13:51:36
WARN Warning message
  13:51:36
ERROR Error occurred
  13:51:36
  code: 500
```

and `application.log` contents:

```
2025-06-14 13:51:36 DEBUG Debug message
2025-06-14 13:51:36 INFO Info message
2025-06-14 13:51:36 WARN Warning message
2025-06-14 13:51:36 ERROR Error occurred [code: 500]
```

## Context

The package provides utilities for associating loggers with context objects, which is especially useful for structured logging in request-based applications. This helps track request-specific information throughout the call chain without manually passing loggers.

```go
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
```

Output will look like:

```go
INFO 🏷 login [username: john]: User authentication
DEBUG 🏷 login [username: john]: Checking user and password
ERROR Authentication failure
```
