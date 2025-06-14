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

Here's the output from [TerminalTreeHandler example](https://github.com/fornellas/slogxpert/blob/main/examples/TerminalTreeHandler/main.go):

![TerminalTreeHandler](https://raw.githubusercontent.com/fornellas/slogxpert/refs/heads/main/examples/TerminalTreeHandler/output.svg)

#### TerminalLineHandler

The `TerminalLineHandler` is similar to `TerminalTreeHandler`, but instead of nesting attributes, it always outputs one line per log message.

Here's the output from [TerminalLineHandler example](https://github.com/fornellas/slogxpert/blob/main/examples/TerminalLineHandler/main.go):

![TerminalLineHandler](https://raw.githubusercontent.com/fornellas/slogxpert/refs/heads/main/examples/TerminalLineHandler/output.svg)

#### Customizing

Handlers can be customized to your taste.

##### TerminalHandlerColorScheme

The full color scheme can be customized, by defining a color scheme, as in the [TerminalHandlerColorScheme example](https://github.com/fornellas/slogxpert/blob/main/examples/TerminalHandlerColorScheme/main.go), which tweaks the colors:

![TerminalHandlerColorScheme](https://raw.githubusercontent.com/fornellas/slogxpert/refs/heads/main/examples/TerminalHandlerColorScheme/output.svg)

##### TerminalHandlerOptions

`TerminalHandlerOptions` provides extensive customization options for both terminal handlers. In the [TerminalHandlerOptions example](https://github.com/fornellas/slogxpert/blob/main/examples/TerminalHandlerOptions/main.go), a custom log level, source code information, sensitive information masking and time are set:

![TerminalHandlerOptions](https://raw.githubusercontent.com/fornellas/slogxpert/refs/heads/main/examples/TerminalHandlerOptions/output.svg)

### BufferedHandler

The `BufferedHandler` allows you to buffer log records in memory until you explicitly flush them to the underlying handler. This is particularly useful when concurrent tasks generate logs, but you with that logging to be clustered per task, instead of interleaved.

In the [BufferedHandler example](https://github.com/fornellas/slogxpert/blob/main/examples/BufferedHandler/main.go), the output below only happens when `Flush()` is called:

![BufferedHandler](https://raw.githubusercontent.com/fornellas/slogxpert/refs/heads/main/examples/BufferedHandler/output.svg)

### MultiHandler

The `MultiHandler` dispatches log records to multiple handlers simultaneously. This is useful when you want to send logs to different destinations or format them differently for various purposes (e.g., console output, file logging, structured JSON for analysis).

In the [MultiHandler example](https://github.com/fornellas/slogxpert/blob/main/examples/MultiHandler/main.go), logs are dispatched to two handlers, one in the terminal:

![MultiHandler](https://raw.githubusercontent.com/fornellas/slogxpert/refs/heads/main/examples/MultiHandler/output.svg)

and another to a file:

```
2025-06-14 13:51:36 DEBUG Debug message
2025-06-14 13:51:36 INFO Info message
2025-06-14 13:51:36 WARN Warning message
2025-06-14 13:51:36 ERROR Error occurred [code: 500]
```

## Context

The package provides utilities for associating loggers with context objects, which is especially useful for structured logging in request-based applications. This helps track request-specific information throughout the call chain without manually passing loggers.

In the [Context example](https://github.com/fornellas/slogxpert/blob/main/examples/Context/main.go), loggers are set and retrieved from the context, resulting in the following output:

![Context](https://raw.githubusercontent.com/fornellas/slogxpert/refs/heads/main/examples/Context/output.svg)
