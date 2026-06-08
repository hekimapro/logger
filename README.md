# logger

A structured logging facade for Go built on top of `log/slog`. Adds an interface abstraction for easier testing and dependency injection, with automatic context propagation for trace and request ID correlation.

## Installation

```bash
go get github.com/hekimapro/logger
```

## Quick Start

```go
import "github.com/hekimapro/logger"

log := logger.New(logger.DefaultConfiguration())

log.Info(ctx, "Server started", "port", 8080)
log.Warn(ctx, "Cache miss", "key", "user:42")
log.Error(ctx, "Database query failed", "error", err, "query", "SELECT ...")
```

## Configuration Presets

### Production

`DefaultConfiguration` emits JSON at `Info` level, suitable for log aggregators like Datadog, Loki, or Cloud Logging.

```go
log := logger.New(logger.DefaultConfiguration())
```

Output:
```json
{"time":"2026-06-08T10:00:00Z","level":"INFO","msg":"Server started","port":8080}
```

Default values:

| Setting | Value |
|---|---|
| `Level` | `Info` |
| `Output` | `os.Stdout` |
| `AddSource` | `false` |
| `UseJSON` | `true` |

### Development

`DevelopmentConfiguration` emits human-readable text at `Debug` level, including file and line numbers.

```go
log := logger.New(logger.DevelopmentConfiguration())
```

Output:
```
time=2026-06-08T10:00:00Z level=DEBUG source=main.go:42 msg="User fetched" id=99
```

## Full Configuration

```go
log := logger.New(logger.Configuration{
    Level:     slog.LevelWarn,  // Only Warn and Error
    Output:    os.Stderr,       // Write to stderr
    AddSource: true,            // Include file:line
    UseJSON:   true,            // JSON format
})
```

### `Configuration` Fields

| Field | Type | Description |
|---|---|---|
| `Level` | `slog.Level` | Minimum level to emit (`Debug`, `Info`, `Warn`, `Error`) |
| `Output` | `io.Writer` | Destination for log entries — defaults to `os.Stdout` |
| `AddSource` | `bool` | Append file name and line number to every entry |
| `UseJSON` | `bool` | `true` for JSON output, `false` for human-readable text |

## Logging with Key-Value Fields

All four log methods accept variadic key-value pairs after the message:

```go
log.Debug(ctx, "Cache lookup", "key", "user:42", "hit", false)
log.Info(ctx, "Request completed", "method", "GET", "path", "/users", "status", 200, "duration_ms", 12)
log.Warn(ctx, "Retry attempt", "attempt", 3, "max", 5)
log.Error(ctx, "Payment failed", "error", err, "user_id", userID, "amount", 9.99)
```

Keys and values alternate: every odd argument is a key (`string`), every even argument is its value (`any`).

## Child Loggers with Fixed Fields

Use `With` to create a child logger that always includes a set of fields — useful for adding request-scoped context once rather than repeating it on every call:

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    requestLogger := log.With(
        "request_id", r.Header.Get("X-Request-ID"),
        "method", r.Method,
        "path", r.URL.Path,
    )

    requestLogger.Info(ctx, "Request received")
    // ... handler logic ...
    requestLogger.Info(ctx, "Request completed", "status", 200)
}
```

Both log entries will automatically carry `request_id`, `method`, and `path`.

## Dependency Injection & Testing

`Logger` is an interface, making it straightforward to swap the real logger for a no-op or a mock in tests:

```go
// Inject into services
type UserService struct {
    log logger.Logger
    db  *sql.DB
}

func NewUserService(log logger.Logger, db *sql.DB) *UserService {
    return &UserService{log: log.With("service", "user"), db: db}
}
```

In tests, pass a no-op logger so log noise doesn't pollute test output:

```go
type noopLogger struct{}

func (n *noopLogger) Debug(ctx context.Context, msg string, args ...any) {}
func (n *noopLogger) Info(ctx context.Context, msg string, args ...any)  {}
func (n *noopLogger) Warn(ctx context.Context, msg string, args ...any)  {}
func (n *noopLogger) Error(ctx context.Context, msg string, args ...any) {}
func (n *noopLogger) With(args ...any) logger.Logger { return n }

svc := NewUserService(&noopLogger{}, testDB)
```

## Writing to a File

Pass any `io.Writer` as the output destination:

```go
file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
if err != nil {
    panic(err)
}

log := logger.New(logger.Configuration{
    Level:   slog.LevelInfo,
    Output:  file,
    UseJSON: true,
})
```

Or write to multiple destinations simultaneously:

```go
log := logger.New(logger.Configuration{
    Level:   slog.LevelInfo,
    Output:  io.MultiWriter(os.Stdout, file),
    UseJSON: true,
})
```

## Logger Interface

All implementations must satisfy:

```go
type Logger interface {
    Debug(ctx context.Context, message string, arguments ...any)
    Info(ctx context.Context, message string, arguments ...any)
    Warn(ctx context.Context, message string, arguments ...any)
    Error(ctx context.Context, message string, arguments ...any)
    With(arguments ...any) Logger
}
```

All implementations must be safe for concurrent use by multiple goroutines. `StandardLogger` (the built-in implementation) satisfies this requirement.

## License

See [LICENSE](LICENSE).
