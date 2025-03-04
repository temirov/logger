package logger

import (
	"bytes"
	"errors"
	"regexp"
	"testing"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// Helper function to capture log output and strip ANSI color codes.
func captureLogOutput(f func()) string {
	var buf bytes.Buffer
	originalOutput := logInstance.logger.Writer()
	defer logInstance.logger.SetOutput(originalOutput)

	logInstance.logger.SetOutput(&buf)
	f()
	return ansiRegex.ReplaceAllString(buf.String(), "")
}

// TestLoggerLevels tests logging behavior at different levels with string inputs.
func TestLoggerLevels(t *testing.T) {
	testCases := []struct {
		name           string
		level          LogLevel
		logFunc        func(args ...interface{})
		message        string
		expectedOutput string
	}{
		{
			name:           "DEBUG level - Debug message is logged",
			level:          DEBUG,
			logFunc:        Debug,
			message:        "Debugging info",
			expectedOutput: "[DEBUG] Debugging info",
		},
		{
			name:           "DEBUG level - Info message is logged",
			level:          DEBUG,
			logFunc:        Info,
			message:        "Info message",
			expectedOutput: "[INFO] Info message",
		},
		{
			name:           "INFO level - Debug message is ignored",
			level:          INFO,
			logFunc:        Debug,
			message:        "Debugging info",
			expectedOutput: "",
		},
		{
			name:           "INFO level - Info message is logged",
			level:          INFO,
			logFunc:        Info,
			message:        "Info message",
			expectedOutput: "[INFO] Info message",
		},
		{
			name:           "WARNING level - Warning message is logged",
			level:          WARNING,
			logFunc:        Warning,
			message:        "Warning message",
			expectedOutput: "[WARNING] Warning message",
		},
		{
			name:           "ERROR level - Error message is logged",
			level:          ERROR,
			logFunc:        Error,
			message:        "Error message",
			expectedOutput: "[ERROR] Error message",
		},
		{
			name:           "WARNING level - Info message is ignored",
			level:          WARNING,
			logFunc:        Info,
			message:        "Info message",
			expectedOutput: "",
		},
		{
			name:           "ERROR level - Warning message is ignored",
			level:          ERROR,
			logFunc:        Warning,
			message:        "Warning message",
			expectedOutput: "",
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			SetLevel(tc.level)
			output := captureLogOutput(func() {
				tc.logFunc(tc.message)
			})

			if tc.expectedOutput != "" {
				prefix := `^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} logger_test.go:\d+: `
				suffix := `\n$`
				re := regexp.MustCompile(prefix + regexp.QuoteMeta(tc.expectedOutput) + suffix)
				if !re.MatchString(output) {
					t.Errorf("Expected output to match %q, but got %q", re.String(), output)
				}
			} else {
				if output != "" {
					t.Errorf("Expected no output, but got %q", output)
				}
			}
		})
	}
}

// TestLoggerErrors tests logging behavior with error inputs.
func TestLoggerErrors(t *testing.T) {
	testCases := []struct {
		name           string
		level          LogLevel
		logFunc        func(args ...interface{})
		err            error
		expectedOutput string
	}{
		{
			name:           "DEBUG level - Error message is logged",
			level:          DEBUG,
			logFunc:        Error,
			err:            errors.New("debug error"),
			expectedOutput: "[ERROR] debug error",
		},
		{
			name:           "INFO level - Error message is logged",
			level:          INFO,
			logFunc:        Error,
			err:            errors.New("info error"),
			expectedOutput: "[ERROR] info error",
		},
		{
			name:           "WARNING level - Error message is logged",
			level:          WARNING,
			logFunc:        Error,
			err:            errors.New("warning error"),
			expectedOutput: "[ERROR] warning error",
		},
		{
			name:           "ERROR level - Error message is logged",
			level:          ERROR,
			logFunc:        Error,
			err:            errors.New("error level error"),
			expectedOutput: "[ERROR] error level error",
		},
		{
			name:           "INFO level - Debug error is ignored",
			level:          INFO,
			logFunc:        Debug,
			err:            errors.New("debug error"),
			expectedOutput: "",
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			SetLevel(tc.level)
			output := captureLogOutput(func() {
				tc.logFunc(tc.err)
			})

			if tc.expectedOutput != "" {
				prefix := `^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} logger_test.go:\d+: `
				suffix := `\n$`
				re := regexp.MustCompile(prefix + regexp.QuoteMeta(tc.expectedOutput) + suffix)
				if !re.MatchString(output) {
					t.Errorf("Expected output to match %q, but got %q", re.String(), output)
				}
			} else {
				if output != "" {
					t.Errorf("Expected no output, but got %q", output)
				}
			}
		})
	}
}

// TestSetLevel verifies that the log level can be dynamically changed.
func TestSetLevel(t *testing.T) {
	SetLevel(INFO)
	if logInstance.level != INFO {
		t.Errorf("Expected log level to be INFO, but got %d", logInstance.level)
	}

	SetLevel(DEBUG)
	if logInstance.level != DEBUG {
		t.Errorf("Expected log level to be DEBUG, but got %d", logInstance.level)
	}
}

// TestColorCodes verifies that the color codes exist for all levels.
func TestColorCodes(t *testing.T) {
	levels := []LogLevel{DEBUG, INFO, WARNING, ERROR}
	for _, level := range levels {
		if _, exists := levelColors[level]; !exists {
			t.Errorf("Expected color code for level %d, but not found", level)
		}
	}
}

// TestVariadicArguments tests logging functions with multiple variadic arguments.
func TestVariadicArguments(t *testing.T) {
	testCases := []struct {
		name           string
		level          LogLevel
		logFunc        func(args ...interface{})
		format         string
		args           []interface{}
		expectedOutput string
	}{
		{
			name:           "DEBUG message with multiple arguments",
			level:          DEBUG,
			logFunc:        Debug,
			format:         "Debugging %s at %d%%",
			args:           []interface{}{"progress", 50},
			expectedOutput: "[DEBUG] Debugging progress at 50%",
		},
		{
			name:           "INFO message with numbers",
			level:          INFO,
			logFunc:        Info,
			format:         "Processed %d items successfully",
			args:           []interface{}{42},
			expectedOutput: "[INFO] Processed 42 items successfully",
		},
		{
			name:           "ERROR message with struct",
			level:          ERROR,
			logFunc:        Error,
			format:         "Error processing user %+v",
			args:           []interface{}{struct{ Name string }{"Alice"}},
			expectedOutput: "[ERROR] Error processing user {Name:Alice}",
		},
		{
			name:           "WARNING message with no formatting",
			level:          WARNING,
			logFunc:        Warning,
			format:         "Simple warning message",
			args:           nil,
			expectedOutput: "[WARNING] Simple warning message",
		},
		{
			name:           "INFO message with string only",
			level:          INFO,
			logFunc:        Info,
			format:         "Info message without formatting",
			args:           []interface{}{},
			expectedOutput: "[INFO] Info message without formatting",
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			SetLevel(tc.level)
			output := captureLogOutput(func() {
				if tc.args != nil && len(tc.args) > 0 {
					args := make([]interface{}, 0, len(tc.args)+1)
					args = append(args, tc.format)
					args = append(args, tc.args...)
					tc.logFunc(args...)
				} else {
					tc.logFunc(tc.format)
				}
			})
			if tc.expectedOutput != "" {
				prefix := `^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} logger_test.go:\d+: `
				suffix := `\n$`
				re := regexp.MustCompile(prefix + regexp.QuoteMeta(tc.expectedOutput) + suffix)
				if !re.MatchString(output) {
					t.Errorf("Expected output to match %q, but got %q", re.String(), output)
				}
			} else {
				if output != "" {
					t.Errorf("Expected no output, but got %q", output)
				}
			}
		})
	}
}

// TestErrorF tests the ErrorF function which logs an error and exits the application.
// Note: Testing functions that call os.Exit is non-trivial because os.Exit terminates the test process.
// One common approach is to refactor the logger to allow injecting a custom exit function, which can be mocked during tests.
// For simplicity, this test is skipped.
func TestErrorF(t *testing.T) {
	t.Skip("Skipping TestErrorF because it calls os.Exit(1)")
}
