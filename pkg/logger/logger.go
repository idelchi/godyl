package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

// Level represents the severity of a log message
type Level string

const (
	DEBUG Level = "debug"
	INFO  Level = "info"
	WARN  Level = "warn"
	ERROR Level = "error"
)

func (l Level) Int() int {
	switch l {
	case DEBUG:
		return 0
	case INFO:
		return 1
	case WARN:
		return 2
	case ERROR:
		return 3
	default:
		return 0
	}
}

func (l Level) IsAllowed() bool {
	switch l {
	case DEBUG, INFO, WARN, ERROR:
		return true
	default:
		return false
	}
}

// Logger struct to hold the log level, output writer, and color functions
type Logger struct {
	level  Level
	output io.Writer
	colors map[Level]*color.Color
}

// New creates a new Logger instance
func NewCustom(level Level, output io.Writer) *Logger {
	if !level.IsAllowed() {
		fmt.Fprintf(os.Stderr, "Invalid log level: %q, setting to %q\n", level, INFO)

		level = INFO
	}

	return &Logger{
		level:  level,
		output: output,
		colors: map[Level]*color.Color{
			DEBUG: color.New(color.FgBlue),
			INFO:  color.New(color.FgGreen),
			WARN:  color.New(color.FgYellow),
			ERROR: color.New(color.FgRed),
		},
	}
}

func New(level Level) *Logger {
	if !level.IsAllowed() {
		fmt.Fprintf(os.Stderr, "Invalid log level: %q, setting to %q\n", level, INFO)

		level = INFO
	}

	return &Logger{
		level:  level,
		output: os.Stdout,
		colors: map[Level]*color.Color{
			DEBUG: color.New(color.FgBlue),
			INFO:  color.New(color.FgGreen),
			WARN:  color.New(color.FgYellow),
			ERROR: color.New(color.FgRed),
		},
	}
}

// log prints a colored message if the log level is sufficient
func (l *Logger) log(level Level, format string, args ...any) {
	if level.Int() >= l.level.Int() {
		message := fmt.Sprintf(format, args...)
		if c, ok := l.colors[level]; ok {
			c.Fprintln(l.output, message)
		} else {
			fmt.Fprintln(l.output, message)
		}
	}
}

func (l *Logger) Always(format string, args ...any) {
	l.log(INFO, format, args...)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...any) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...any) {
	l.log(INFO, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...any) {
	l.log(WARN, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...any) {
	l.log(ERROR, format, args...)
}
