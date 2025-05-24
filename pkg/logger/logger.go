//go:generate go tool enumer -type=Level -output level_enumer___generated.go -transform=lower
package logger

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

// Level represents the severity of a log message.
// It can be one of DEBUG, INFO, WARN, ERROR, ALWAYS, or SILENT.
type Level int

const (
	SILENT Level = iota - 1 // no logging
	DEBUG                   // detailed debug information
	INFO                    // normal operational messages
	WARN                    // potentially harmful situations
	ERROR                   // error events
	ALWAYS                  // always shown regardless of current log level
)

// Logger holds the configuration for logging.
// It includes the current logging level, the output writer, and color mappings for each log level.
type Logger struct {
	output io.Writer
	colors map[Level]*color.Color
	level  Level
}

// ErrInvalidLogLevel is returned when an invalid log level is provided.
var ErrInvalidLogLevel = errors.New("invalid log level")

// NewCustom creates a new Logger instance with the specified log level and output writer.
// If an invalid log level is provided, it defaults to INFO.
func NewCustom(level Level, output io.Writer) (*Logger, error) {
	if !level.IsALevel() {
		return nil, fmt.Errorf("%w: %q, setting to %q\n", ErrInvalidLogLevel, level, INFO)
	}

	return &Logger{
		level:  level,
		output: output,
		colors: map[Level]*color.Color{
			DEBUG:  color.New(color.FgBlue),
			INFO:   color.New(color.FgGreen),
			WARN:   color.New(color.FgYellow),
			ERROR:  color.New(color.FgRed),
			ALWAYS: color.New(color.FgGreen),
		},
	}, nil
}

// New creates a new Logger instance with the specified log level and writes to stdout.
func New(level Level) (*Logger, error) {
	return NewCustom(level, os.Stdout)
}

// log prints a log message with the specified level, respecting the configured minimum level.
func (l *Logger) log(level Level, format string, args ...any) {
	if level < l.level && level != ALWAYS {
		return
	}

	message := fmt.Sprintf(format, args...)
	if c, ok := l.colors[level]; ok {
		_, _ = c.Fprintln(l.output, message)
	} else {
		_ = writeSilently(l.output, message)
	}
}

func writeSilently(w io.Writer, msg string) error {
	_, err := fmt.Fprintln(w, msg)

	return err
}

func (l *Logger) logPlain(level Level, message string) {
	l.log(level, "%s", message)
}

// Always logs a message at the ALWAYS level, regardless of the current log level.
func (l *Logger) Always(msg string) { l.logPlain(ALWAYS, msg) }

// Alwaysf logs a formatted message at the ALWAYS level.
func (l *Logger) Alwaysf(f string, a ...any) { l.log(ALWAYS, f, a...) }

// Debug logs a message at the DEBUG level, if enabled.
func (l *Logger) Debug(msg string) { l.logPlain(DEBUG, msg) }

// Debugf logs a formatted message at the DEBUG level.
func (l *Logger) Debugf(f string, a ...any) { l.log(DEBUG, f, a...) }

// Info logs a message at the INFO level, if enabled.
func (l *Logger) Info(msg string) { l.logPlain(INFO, msg) }

// Infof logs a formatted message at the INFO level.
func (l *Logger) Infof(f string, a ...any) { l.log(INFO, f, a...) }

// Warn logs a message at the WARN level, if enabled.
func (l *Logger) Warn(msg string) { l.logPlain(WARN, msg) }

// Warnf logs a formatted message at the WARN level.
func (l *Logger) Warnf(f string, a ...any) { l.log(WARN, f, a...) }

// Error logs a message at the ERROR level, if enabled.
func (l *Logger) Error(msg string) { l.logPlain(ERROR, msg) }

// Errorf logs a formatted message at the ERROR level.
func (l *Logger) Errorf(f string, a ...any) { l.log(ERROR, f, a...) }

// SetLevel updates the logger's level at runtime.
func (l *Logger) SetLevel(level Level) {
	if level.IsALevel() {
		l.level = level
	}
}

// Level returns the current log level.
func (l *Logger) Level() Level {
	return l.level
}
