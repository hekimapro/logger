package main

import (
    "context"
    "log/slog"
    "time"

    "github.com/hekimapro/logger"
)

func main() {
    println("=== Logger Package Examples ===\n")

    // Example 1: Basic usage with production config
    exampleBasicUsage()

    // Example 2: Development configuration
    exampleDevelopmentConfig()

    // Example 3: Child loggers with fixed fields
    exampleChildLoggers()

    // Example 4: Different log levels
    exampleLogLevels()

    // Example 5: Structured logging
    exampleStructuredLogging()
}

func exampleBasicUsage() {
    println("1. Basic Usage (Production Configuration)")
    println("   --------------------------------------")

    // Create production logger (JSON output, Info level)
    appLogger := logger.New(logger.DefaultConfiguration())
    ctx := context.Background()

    appLogger.Info(ctx, "Application started", "version", "1.0.0", "environment", "production")
    appLogger.Info(ctx, "Server listening", "port", 8080, "host", "0.0.0.0")

    println()
}

func exampleDevelopmentConfig() {
    println("2. Development Configuration")
    println("   ------------------------")

    // Create development logger (text output, Debug level)
    devLogger := logger.New(logger.DevelopmentConfiguration())
    ctx := context.Background()

    devLogger.Debug(ctx, "This debug message only appears in development")
    devLogger.Info(ctx, "User logged in", "user_id", 12345, "ip", "192.168.1.100")
    devLogger.Warn(ctx, "Rate limit approaching", "current", 95, "limit", 100)
    devLogger.Error(ctx, "Failed to send email", "error", "connection timeout", "retry_count", 3)

    println()
}

func exampleChildLoggers() {
    println("3. Child Loggers with Fixed Fields")
    println("   -------------------------------")

    baseLogger := logger.New(logger.DefaultConfiguration())
    ctx := context.Background()

    // Create component-specific loggers
    databaseLogger := baseLogger.With("component", "database", "service", "postgres")
    httpLogger := baseLogger.With("component", "http", "service", "api")
    workerLogger := baseLogger.With("component", "worker", "service", "email")

    databaseLogger.Info(ctx, "Connection pool initialized", "max_connections", 25, "idle_connections", 10)
    httpLogger.Info(ctx, "Request received", "method", "POST", "path", "/api/users", "duration_ms", 45)
    workerLogger.Info(ctx, "Processing job", "job_id", "job-12345", "queue", "emails")

    println()
}

func exampleLogLevels() {
    println("4. Different Log Levels")
    println("   -------------------")

    appLogger := logger.New(logger.DefaultConfiguration())
    ctx := context.Background()

    // Simulate different severity levels
    appLogger.Debug(ctx, "Cache hit", "key", "user:123", "ttl_remaining", 300)
    appLogger.Info(ctx, "Order placed", "order_id", "ORD-789", "total", 99.95)
    appLogger.Warn(ctx, "High memory usage", "memory_percent", 85, "threshold", 80)
    appLogger.Error(ctx, "Payment gateway timeout", "gateway", "stripe", "attempt", 3)

    println()
}

func exampleStructuredLogging() {
    println("5. Structured Logging with Complex Data")
    println("   -----------------------------------")

    appLogger := logger.New(logger.DefaultConfiguration())
    ctx := context.Background()

    // Log with multiple fields
    requestStartTime := time.Now()

    // Simulate some work
    time.Sleep(10 * time.Millisecond)

    requestDuration := time.Since(requestStartTime)

    appLogger.Info(ctx, "API request completed",
        slog.String("method", "POST"),
        slog.String("path", "/api/v1/users"),
        slog.Int("status_code", 201),
        slog.Duration("duration", requestDuration),
        slog.String("client_ip", "192.168.1.100"),
        slog.String("user_agent", "Mozilla/5.0"),
        slog.Int64("response_size", 1024),
    )

    // Log error with context
    err := simulateDatabaseError()
    appLogger.Error(ctx, "Database operation failed",
        slog.String("operation", "INSERT"),
        slog.String("table", "users"),
        slog.Any("error", err),
        slog.Int("retry_count", 3),
        slog.Bool("will_retry", true),
    )

    println()
}

func simulateDatabaseError() error {
    return &databaseError{message: "connection pool exhausted", code: "POOL_FULL"}
}

type databaseError struct {
    message string
    code    string
}

func (err *databaseError) Error() string {
    return err.message
}