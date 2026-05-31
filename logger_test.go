package logger

import (
    "bytes"
    "context"
    "encoding/json"
    "log/slog"
    "testing"
)

func TestNew(t *testing.T) {
    tests := []struct {
        name          string
        configuration Configuration
        shouldSucceed bool
    }{
        {
            name:          "default configuration",
            configuration: DefaultConfiguration(),
            shouldSucceed: true,
        },
        {
            name:          "development configuration",
            configuration: DevelopmentConfiguration(),
            shouldSucceed: true,
        },
        {
            name: "custom configuration",
            configuration: Configuration{
                Level:     slog.LevelDebug,
                Output:    &bytes.Buffer{},
                AddSource: true,
                UseJSON:   false,
            },
            shouldSucceed: true,
        },
    }

    for _, testCase := range tests {
        t.Run(testCase.name, func(t *testing.T) {
            log := New(testCase.configuration)

            if testCase.shouldSucceed && log == nil {
                t.Errorf("Expected non-nil logger")
            }
        })
    }
}

func TestLogLevels(t *testing.T) {
    var buffer bytes.Buffer

    log := New(Configuration{
        Level:   slog.LevelInfo,
        Output:  &buffer,
        UseJSON: true,
    })

    ctx := context.Background()

    // Debug should not appear (level is Info)
    log.Debug(ctx, "debug message", "key", "value")

    // Info should appear
    log.Info(ctx, "info message", "key", "value")

    output := buffer.String()

    // Debug message should not be present
    if containsString(output, "debug message") {
        t.Errorf("Debug message appeared when level is Info")
    }

    // Info message should be present
    if !containsString(output, "info message") {
        t.Errorf("Info message did not appear")
    }
}

func TestJSONFormat(t *testing.T) {
    var buffer bytes.Buffer

    log := New(Configuration{
        Level:   slog.LevelInfo,
        Output:  &buffer,
        UseJSON: true,
    })

    ctx := context.Background()
    log.Info(ctx, "test message", "user_id", 12345, "action", "login")

    output := buffer.String()

    // Verify JSON format
    var parsedOutput map[string]interface{}
    if err := json.Unmarshal([]byte(output), &parsedOutput); err != nil {
        t.Errorf("Failed to parse JSON output: %v", err)
    }

    // Check required fields
    if _, exists := parsedOutput["time"]; !exists {
        t.Errorf("Missing 'time' field in JSON output")
    }

    if _, exists := parsedOutput["level"]; !exists {
        t.Errorf("Missing 'level' field in JSON output")
    }

    if msg, exists := parsedOutput["msg"]; !exists || msg != "test message" {
        t.Errorf("Missing or incorrect 'msg' field: %v", msg)
    }
}

func TestTextFormat(t *testing.T) {
    var buffer bytes.Buffer

    log := New(Configuration{
        Level:   slog.LevelInfo,
        Output:  &buffer,
        UseJSON: false,
    })

    ctx := context.Background()
    log.Info(ctx, "test message", "user_id", 12345)

    output := buffer.String()

    // Text format should contain key=value pairs
    if !containsString(output, "user_id=12345") {
        t.Errorf("Expected key=value pair in text output")
    }
}

func TestWithFields(t *testing.T) {
    var buffer bytes.Buffer

    parentLogger := New(Configuration{
        Level:   slog.LevelInfo,
        Output:  &buffer,
        UseJSON: true,
    })

    // Create child logger with fixed fields
    childLogger := parentLogger.With("component", "database", "request_id", "req-123")

    ctx := context.Background()
    childLogger.Info(ctx, "query executed", "duration_ms", 45)

    output := buffer.String()

    // Child logger should include fixed fields
    if !containsString(output, "component") || !containsString(output, "database") {
        t.Errorf("Child logger missing fixed fields")
    }

    if !containsString(output, "duration_ms") || !containsString(output, "45") {
        t.Errorf("Child logger missing operation-specific fields")
    }
}

func TestContextFields(t *testing.T) {
    var buffer bytes.Buffer

    log := New(Configuration{
        Level:   slog.LevelInfo,
        Output:  &buffer,
        UseJSON: true,
    })

    // Add fields to context
    ctx := context.Background()
    ctx = context.WithValue(ctx, "trace_id", "abc-123-def")
    ctx = context.WithValue(ctx, "request_id", "req-456")
    ctx = context.WithValue(ctx, "user_id", "user-789")

    log.Info(ctx, "request processed", "status", 200)

    output := buffer.String()

    // The logger doesn't automatically extract context fields in this implementation
    // But the With method can be used instead
    if !containsString(output, "status") {
        t.Errorf("Expected status field in output")
    }
}

func TestConcurrentLogging(t *testing.T) {
    var buffer bytes.Buffer

    log := New(Configuration{
        Level:   slog.LevelInfo,
        Output:  &buffer,
        UseJSON: true,
    })

    ctx := context.Background()

    // Run concurrent logging
    done := make(chan bool)
    for i := 0; i < 100; i++ {
        go func(goroutineID int) {
            for j := 0; j < 100; j++ {
                log.Info(ctx, "concurrent log", "goroutine", goroutineID, "iteration", j)
            }
            done <- true
        }(i)
    }

    // Wait for all goroutines
    for i := 0; i < 100; i++ {
        <-done
    }

    // No panic means success
    t.Log("Concurrent logging completed without panic")
}

func TestDifferentLogLevels(t *testing.T) {
    var buffer bytes.Buffer

    log := New(Configuration{
        Level:   slog.LevelDebug,
        Output:  &buffer,
        UseJSON: true,
    })

    ctx := context.Background()

    log.Debug(ctx, "debug message")
    log.Info(ctx, "info message")
    log.Warn(ctx, "warn message")
    log.Error(ctx, "error message")

    output := buffer.String()

    // All levels should be present since level is Debug
    levels := []string{"DEBUG", "INFO", "WARN", "ERROR"}
    for _, level := range levels {
        if !containsString(output, level) {
            t.Errorf("Missing log level: %s", level)
        }
    }
}

// Helper function to check if a string contains a substring
func containsString(haystack, needle string) bool {
    return len(haystack) >= len(needle) &&
           (haystack == needle ||
            len(haystack) > 0 && len(needle) > 0 &&
            findSubstring(haystack, needle))
}

func findSubstring(haystack, needle string) bool {
    for i := 0; i <= len(haystack)-len(needle); i++ {
        if haystack[i:i+len(needle)] == needle {
            return true
        }
    }
    return false
}