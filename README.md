# Logger

A lightweight, color-coded logging library for Go with support for log levels, dynamic log filtering, and variadic
arguments. Designed to enhance terminal output with clear, readable, and formatted logs.

---

## Features

- ✅ Log levels: `DEBUG`, `INFO`, `WARNING`, `ERROR`
- ✅ Colorized terminal output
- ✅ Global, package-level logger
- ✅ Dynamic log level configuration
- ✅ Supports `string`, `error`, and formatted messages with variadic arguments
- ✅ Safe to use concurrently
- ✅ Well-tested with comprehensive test coverage

---

## Installation

```
go get github.com/temirov/logger
```

---

## Usage

### Import

```
import "github.com/temirov/logger"
```

### Basic Logging

```
logger.Debug("Debugging info")
logger.Info("Application started")
logger.Warning("Warning: %s", "disk space low")
logger.Error(errors.New("something went wrong"))
```

### Fatal Error Logging (exits the app)

```
logger.ErrorF("Fatal error, exiting")
```

---

## Log Levels

| Level   | Description           | Usage                        |
|---------|-----------------------|------------------------------|
| DEBUG   | Verbose debug details | `logger.Debug()`             |
| INFO    | General info logs     | `logger.Info()`              |
| WARNING | Warnings              | `logger.Warning()`           |
| ERROR   | Errors                | `logger.Error()`, `ErrorF()` |

Default log level is `DEBUG`. To change the log level at runtime:

```
logger.SetLevel(logger.INFO)
```

Or from a string value (useful for configs/env vars):

```
level := "warning"
logger.SetLogLevel(&level)
```

---

## Output Example

```
2025/03/03 14:23:01 logger.go:123: [DEBUG] Starting process
2025/03/03 14:23:01 logger.go:124: [INFO] Process completed successfully
2025/03/03 14:23:01 logger.go:125: [WARNING] Low memory detected
2025/03/03 14:23:01 logger.go:126: [ERROR] failed to connect to database
```

Colors are applied automatically when running in a terminal that supports ANSI codes.

---

## Testing

Run tests with:

```
go test ./...
```

Tests cover:

- Log level filtering
- Error handling
- Variadic formatting
- Dynamic level setting
- Color code correctness

---

## License

MIT License

---

## Author

[temirov](https://github.com/temirov)