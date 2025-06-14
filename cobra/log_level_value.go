package cobra

import (
	"fmt"
	"log/slog"
	"strings"
)

var DefaultLevel = slog.LevelInfo

// LogLevelValue implements [pflag.Value] interface for [slog.Level].
type LogLevelValue slog.Level

func NewLogLevelValue() *LogLevelValue {
	logLevelValue := LogLevelValue(DefaultLevel)
	return &logLevelValue
}

func (l LogLevelValue) String() string {
	return strings.ToLower(slog.Level(l).String())
}

func (l *LogLevelValue) Set(value string) error {
	return (*slog.Level)(l).UnmarshalText([]byte(value))
}

func (l *LogLevelValue) Reset() {
	if err := l.Set(DefaultLevel.String()); err != nil {
		panic(err)
	}
}

func (l LogLevelValue) Type() string {
	return fmt.Sprintf("[%s]", strings.Join([]string{
		strings.ToLower(slog.LevelDebug.String()),
		strings.ToLower(slog.LevelInfo.String()),
		strings.ToLower(slog.LevelWarn.String()),
		strings.ToLower(slog.LevelError.String()),
	}, "|"))
}

func (l LogLevelValue) Level() slog.Level {
	return slog.Level(l)
}
