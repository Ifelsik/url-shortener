package logger

import (
	"context"
	"errors"
)

type LoggerFields map[string]any

// Formatter types.
// TextFormatter is recommended for TTY output.
const (
	TextFormatter = 0
	JSONFormatter = 1
)

type LoggerConfig struct {
	Level      uint8 // TODO: add levels
	Formatter  uint8
	ShowCaller bool
}

type Logger interface {
	Infof(format string, args ...any)
	Debugf(format string, args ...any)
	Warningf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
	WithFields(LoggerFields) Logger
}

type loggerCtxKey string

const loggerKey loggerCtxKey = "logger"

var ErrNoLoggerInCtx = errors.New("no logger in context")

func ToContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) (Logger, error) {
	logger, ok := ctx.Value(loggerKey).(Logger)
	if logger == nil || !ok {
		return nil, ErrNoLoggerInCtx
	}
	return logger, nil
}
