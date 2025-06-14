[![ci](https://github.com/fornellas/slogxpert/actions/workflows/ci.yaml/badge.svg)](https://github.com/fornellas/slogxpert/actions/workflows/ci.yaml) [![update_deps](https://github.com/fornellas/slogxpert/actions/workflows/update_deps.yaml/badge.svg)](https://github.com/fornellas/slogxpert/actions/workflows/update_deps.yaml) [![Go Report Card](https://goreportcard.com/badge/github.com/fornellas/slogxpert)](https://goreportcard.com/report/github.com/fornellas/slogxpert) [![Coverage Status](https://coveralls.io/repos/github/fornellas/slogxpert/badge.svg?branch=master)](https://coveralls.io/github/fornellas/slogxpert?branch=master) [![Go Reference](https://pkg.go.dev/badge/github.com/fornellas/slogxpert.svg)](https://pkg.go.dev/github.com/fornellas/slogxpert) [![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0) [![Buy me a beer: donate](https://img.shields.io/badge/Donate-Buy%20me%20a%20beer-yellow)](https://www.paypal.com/donate?hosted_button_id=AX26JVRT2GS2Q)

# slogxpert

This package extends Go's `log/slog`, providing handlers for structured colored console output, buffering and other goodies.

The full detailed documentation can be found [here](https://pkg.go.dev/github.com/fornellas/slogxpert), here are some highlights.

Contributions are accepted, check [README.development.md](README.development.md) for instructions on how to run the build, and cut a PR with your changes.

## Handlers

### Terminal

#### TerminalTreeHandler

TODO add explanation about this handler, add short code example of how to use it, and what the output looks like.

#### TerminalLineHandler

TODO add explanation about this handler, add short code example of how to use it, and what the output looks like.

#### Customizing

##### TerminalHandlerColorScheme

TODO add an example of how to customize the color scheme for either TerminalTreeHandler or TerminalLineHandler

##### TerminalHandlerOptions

TODO add an example of how to customize either TerminalTreeHandler or TerminalLineHandler using TerminalHandlerOptions.

### BufferedHandler

TODO add explanation about this handler, add short code example of how to use it, and what the output looks like.

### MultiHandler

TODO add explanation about this handler, add short code example of how to use it, and what the output looks like.

## Context

TODO add explanation and examples on how to use the functions from context.go to set / retrieve a logger from the context
