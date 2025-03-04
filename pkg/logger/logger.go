package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type LogLevel int

// Log Levels
const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

// Colors for terminal output (optional)
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Green  = "\033[32m"
	Blue   = "\033[34m"
)

// Map log levels to colors
var levelColors = map[LogLevel]string{
	DEBUG:   Blue,
	INFO:    Green,
	WARNING: Yellow,
	ERROR:   Red,
}

// Logger struct with level filtering
type Logger struct {
	level  LogLevel
	logger *log.Logger // Embed the standard logger
}

// Global package-level logger
var logInstance *Logger

// init initializes the default logger at DEBUG level
func init() {
	logInstance = NewLogger(DEBUG)
}

// NewLogger initializes a logger with a specific log level
func NewLogger(level LogLevel) *Logger {
	stdLogger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	return &Logger{level: level, logger: stdLogger}
}

// SetLevel dynamically sets the logging level
func SetLevel(level LogLevel) {
	logInstance.level = level
}

func SetLogLevel(level *string) {
	var logLevel LogLevel

	switch strings.ToLower(*level) {
	case "info", "i":
		logLevel = INFO
	case "debug", "d":
		logLevel = DEBUG
	case "warning", "warn", "w":
		logLevel = WARNING
	default:
		logLevel = ERROR
	}

	SetLevel(logLevel)
}

// logMessage formats and logs messages with levels and variadic arguments
func (l *Logger) logMessage(level LogLevel, label string, color string, args ...interface{}) {
	if level < l.level {
		return // Do not log messages below the current level
	}

	var message string

	if len(args) == 0 {
		message = ""
	} else if len(args) == 1 {
		switch v := args[0].(type) {
		case error:
			message = v.Error()
		case string:
			message = v
		default:
			message = fmt.Sprintf("%+v", v)
		}
	} else {
		// Assume the first argument is a format string
		format, ok := args[0].(string)
		if !ok {
			// If the first argument isn't a string, format the entire arguments
			message = fmt.Sprintf("%+v", args)
		} else {
			message = fmt.Sprintf(format, args[1:]...)
		}
	}

	// Use log.Output with calldepth = 3 to go up to the caller of Debug/Info/Warning/Error
	calldepth := 3
	_ = l.logger.Output(calldepth, fmt.Sprintf("%s[%s] %s%s", color, label, message, Reset))
}

// Debug logs a debug message
func Debug(args ...interface{}) {
	logInstance.logMessage(DEBUG, "DEBUG", levelColors[DEBUG], args...)
}

// Info logs an info message
func Info(args ...interface{}) {
	logInstance.logMessage(INFO, "INFO", levelColors[INFO], args...)
}

// Warning logs a warning message
func Warning(args ...interface{}) {
	logInstance.logMessage(WARNING, "WARNING", levelColors[WARNING], args...)
}

// Error logs an error message
func Error(args ...interface{}) {
	logInstance.logMessage(ERROR, "ERROR", levelColors[ERROR], args...)
}

// ErrorF logs an error message and exits the app
func ErrorF(args ...interface{}) {
	Error(args...)
	os.Exit(1)
}
