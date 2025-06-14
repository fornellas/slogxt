package cobra

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	_ "github.com/spf13/pflag"

	"github.com/fornellas/slogxpert"
)

// LogHandlerValueOptions holds some options for [slogxpert.TerminalHandlerOptions].
type LogHandlerValueOptions struct {
	Level              slog.Level
	AddSource          bool
	TerminalTime       bool
	TerminalForceColor bool
}

var logHandlerNameFnMap = map[string]func(io.Writer, LogHandlerValueOptions) slog.Handler{
	"terminal-tree": func(writer io.Writer, options LogHandlerValueOptions) slog.Handler {
		var timeLayout string
		if options.TerminalTime {
			timeLayout = time.DateTime
		}
		return slogxpert.NewTerminalTreeHandler(writer, &slogxpert.TerminalHandlerOptions{
			HandlerOptions: slog.HandlerOptions{
				Level:     options.Level,
				AddSource: options.AddSource,
			},
			TimeLayout: timeLayout,
			ForceColor: options.TerminalForceColor,
		})
	},
	"terminal-line": func(writer io.Writer, options LogHandlerValueOptions) slog.Handler {
		var timeLayout string
		if options.TerminalTime {
			timeLayout = time.DateTime
		}
		return slogxpert.NewTerminalLineHandler(writer, &slogxpert.TerminalHandlerOptions{
			HandlerOptions: slog.HandlerOptions{
				Level:     options.Level,
				AddSource: options.AddSource,
			},
			TimeLayout: timeLayout,
			ForceColor: options.TerminalForceColor,
		})
	},
	"json": func(writer io.Writer, options LogHandlerValueOptions) slog.Handler {
		return slog.NewJSONHandler(writer, &slog.HandlerOptions{
			AddSource: options.AddSource,
			Level:     options.Level,
		})
	},
}

func logHandlerNames() (names []string) {
	for name := range logHandlerNameFnMap {
		names = append(names, name)
	}
	return names
}

var DefaultLogHandlerValue = "terminal-tree"

// LogHandlerValue implements [pflag.Value] interface for a [slog.Handler].
type LogHandlerValue struct {
	name string
}

func NewLogHandlerValue() *LogHandlerValue {
	return &LogHandlerValue{name: DefaultLogHandlerValue}
}

func (h *LogHandlerValue) String() string {
	return h.name
}

func (h *LogHandlerValue) Set(value string) error {
	if _, ok := logHandlerNameFnMap[value]; !ok {
		return fmt.Errorf("invalid log handler name '%s', valid options are %s", value, h.Type())
	}
	h.name = value
	return nil
}

func (h *LogHandlerValue) Reset() {
	if err := h.Set(DefaultLogHandlerValue); err != nil {
		panic(err)
	}
}

func (h *LogHandlerValue) Type() string {
	return fmt.Sprintf("[%s]", strings.Join(logHandlerNames(), "|"))
}

func (h *LogHandlerValue) GetHandler(
	writer io.Writer, options LogHandlerValueOptions,
) slog.Handler {
	fn, ok := logHandlerNameFnMap[h.name]
	if !ok {
		panic("bug detected: invalid handler name")
	}
	return fn(writer, options)
}
