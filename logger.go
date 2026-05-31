// Package logger provides a structured logging facade with context propagation.
// Built on Go's log/slog package but adds interface abstraction for easier testing
// and dependency injection.
//
// The logger automatically includes trace and request IDs when provided through
// context.Context, enabling end-to-end request correlation.
package logger

import (
    "context"
    "io"
    "log/slog"
    "os"
)

// Logger defines the interface for structured logging.
// Implementations must be safe for concurrent use by multiple goroutines.
type Logger interface {
    // Debug logs at debug level, typically disabled in production
    Debug(ctx context.Context, message string, arguments ...any)

    // Info logs at info level for normal operations
    Info(ctx context.Context, message string, arguments ...any)

    // Warn logs at warning level for unexpected but recoverable issues
    Warn(ctx context.Context, message string, arguments ...any)

    // Error logs at error level for operations that failed
    Error(ctx context.Context, message string, arguments ...any)

    // With creates a child logger with additional fixed fields
    With(arguments ...any) Logger
}

// StandardLogger implements the Logger interface using Go's slog package.
type StandardLogger struct {
    internal *slog.Logger
}

// Configuration holds all options for creating a new logger.
type Configuration struct {
    // Level determines which log messages are emitted (Debug, Info, Warn, Error)
    Level slog.Level

    // Output where log entries are written (defaults to os.Stdout)
    Output io.Writer

    // AddSource includes file and line number information
    AddSource bool

    // UseJSON determines output format (true for JSON, false for human-readable text)
    UseJSON bool
}

// DefaultConfiguration returns a production-ready configuration.
func DefaultConfiguration() Configuration {
    return Configuration{
        Level:     slog.LevelInfo,
        Output:    os.Stdout,
        AddSource: false,
        UseJSON:   true,
    }
}

// DevelopmentConfiguration returns a developer-friendly configuration.
func DevelopmentConfiguration() Configuration {
    return Configuration{
        Level:     slog.LevelDebug,
        Output:    os.Stdout,
        AddSource: true,
        UseJSON:   false,
    }
}

// New creates a new logger with the specified configuration.
func New(configuration Configuration) Logger {
    handlerOptions := &slog.HandlerOptions{
        Level:     configuration.Level,
        AddSource: configuration.AddSource,
    }

    var handler slog.Handler
    if configuration.UseJSON {
        handler = slog.NewJSONHandler(configuration.Output, handlerOptions)
    } else {
        handler = slog.NewTextHandler(configuration.Output, handlerOptions)
    }

    return &StandardLogger{
        internal: slog.New(handler),
    }
}

// Debug implements Logger.Debug
func (standardLogger *StandardLogger) Debug(ctx context.Context, message string, arguments ...any) {
    standardLogger.internal.Log(ctx, slog.LevelDebug, message, arguments...)
}

// Info implements Logger.Info
func (standardLogger *StandardLogger) Info(ctx context.Context, message string, arguments ...any) {
    standardLogger.internal.Log(ctx, slog.LevelInfo, message, arguments...)
}

// Warn implements Logger.Warn
func (standardLogger *StandardLogger) Warn(ctx context.Context, message string, arguments ...any) {
    standardLogger.internal.Log(ctx, slog.LevelWarn, message, arguments...)
}

// Error implements Logger.Error
func (standardLogger *StandardLogger) Error(ctx context.Context, message string, arguments ...any) {
    standardLogger.internal.Log(ctx, slog.LevelError, message, arguments...)
}

// With creates a child logger with additional fields
func (standardLogger *StandardLogger) With(arguments ...any) Logger {
    return &StandardLogger{
        internal: standardLogger.internal.With(arguments...),
    }
}