package cobra

import (
	"io"
	"log/slog"

	"github.com/spf13/cobra"
)

var logLevelValue = NewLogLevelValue()

var logHandlerValue = NewLogHandlerValue()

var defaultLogHandlerAddSource = false
var logHandlerAddSource = defaultLogHandlerAddSource

var defaultLogHandlerTerminlTime = false
var logHandlerTerminalTime = defaultLogHandlerTerminlTime

var defaultLogHandlerTerminalForceColor = false
var logHandlerTerminalForceColor = defaultLogHandlerTerminalForceColor

// AddLoggerFlags adds logger related flags to a Cobra command. A [slog.Logger] can then be retrieved
// with [GetLogger].
//
// These flags enable defining the log level, handler and some customization.
func AddLoggerFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().VarP(logLevelValue, "log-level", "l", "Logging level")

	cmd.PersistentFlags().VarP(logHandlerValue, "log-handler", "", "Logging handler")

	cmd.PersistentFlags().BoolVarP(
		&logHandlerAddSource, "log-handler-add-source", "", defaultLogHandlerAddSource,
		"Include source code position of the log statement when logging",
	)

	cmd.PersistentFlags().BoolVarP(
		&logHandlerTerminalTime, "log-handler-terminal-time", "", defaultLogHandlerTerminlTime,
		"Enable time for terminal handlers",
	)

	cmd.PersistentFlags().BoolVarP(
		&logHandlerTerminalForceColor, "log-handler-terminal-force-color", "", defaultLogHandlerTerminalForceColor,
		"Force ANSI colors even when terminal is not detected",
	)
}

// GetLogger returns a [slog.Logger] crafted as a function of the Cobra command flags from [AddLoggerFlags].
func GetLogger(writer io.Writer) *slog.Logger {
	handler := logHandlerValue.GetHandler(
		writer,
		LogHandlerValueOptions{
			Level:              logLevelValue.Level(),
			AddSource:          logHandlerAddSource,
			TerminalTime:       logHandlerTerminalTime,
			TerminalForceColor: logHandlerTerminalForceColor,
		},
	)
	return slog.New(handler)
}

// Reset the value of all flags from [AddLoggerFlags].
func Reset() {
	logLevelValue.Reset()
	logHandlerValue.Reset()
	logHandlerAddSource = defaultLogHandlerAddSource
	logHandlerTerminalTime = defaultLogHandlerTerminlTime
	logHandlerTerminalForceColor = defaultLogHandlerTerminalForceColor
}
