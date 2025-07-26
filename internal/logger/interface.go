package logger

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
	WithFields(LoggerFields) *Logger
}
