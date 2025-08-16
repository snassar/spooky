# Logging System API Reference

## Overview

This document provides a comprehensive API reference for the spooky logging system. It covers all interfaces, types, methods, and implementation details for developers working with the logging system.

**Status: Partially Implemented** - The logging system has basic functionality but the interface definitions and implementation details have known issues that need to be addressed.

## Core Interfaces

### LogManager Interface

The `LogManager` interface provides the primary entry point for logging operations:

```go
type LogManager interface {
    // GetLogger gets a logger instance for the given name
    GetLogger(name string) Logger
    
    // GetLoggerWithContext gets a logger with context
    GetLoggerWithContext(name string, ctx LogContext) Logger
    
    // SetLevel sets the logging level for all loggers
    SetLevel(level LogLevel) error
    
    // GetLevel gets the current logging level
    GetLevel() LogLevel
    
    // AddHandler adds a log handler
    AddHandler(handler LogHandler) error
    
    // RemoveHandler removes a log handler
    RemoveHandler(handler LogHandler) error
    
    // Flush flushes all pending log entries
    Flush() error
    
    // Close closes the log manager
    Close() error
}
```

### Logger Interface

The `Logger` interface provides logging methods:

```go
type Logger interface {
    // Debug logs a debug message
    Debug(msg string, fields ...LogField)
    
    // Info logs an info message
    Info(msg string, fields ...LogField)
    
    // Warn logs a warning message
    Warn(msg string, fields ...LogField)
    
    // Error logs an error message
    Error(msg string, fields ...LogField)
    
    // Fatal logs a fatal message and exits
    Fatal(msg string, fields ...LogField)
    
    // WithField adds a field to the logger
    WithField(key string, value interface{}) Logger
    
    // WithFields adds multiple fields to the logger
    WithFields(fields map[string]interface{}) Logger
    
    // WithContext adds context to the logger
    WithContext(ctx LogContext) Logger
}
```

### LogHandler Interface

The `LogHandler` interface defines log output handlers:

```go
type LogHandler interface {
    // Handle handles a log entry
    Handle(entry LogEntry) error
    
    // Flush flushes pending log entries
    Flush() error
    
    // Close closes the handler
    Close() error
}
```

## Core Types

### LogLevel

```go
type LogLevel int

const (
    LogLevelDebug LogLevel = iota
    LogLevelInfo
    LogLevelWarn
    LogLevelError
    LogLevelFatal
)

func (l LogLevel) String() string {
    switch l {
    case LogLevelDebug:
        return "debug"
    case LogLevelInfo:
        return "info"
    case LogLevelWarn:
        return "warn"
    case LogLevelError:
        return "error"
    case LogLevelFatal:
        return "fatal"
    default:
        return "unknown"
    }
}
```

### LogEntry

```go
type LogEntry struct {
    Timestamp time.Time
    Level     LogLevel
    Message   string
    Fields    map[string]interface{}
    Context   LogContext
    Logger    string
}
```

### LogContext

```go
type LogContext struct {
    RequestID    string
    UserID       string
    SessionID    string
    CorrelationID string
    Metadata     map[string]interface{}
}
```

### LogField

```go
type LogField struct {
    Key   string
    Value interface{}
}

func String(key, value string) LogField {
    return LogField{Key: key, Value: value}
}

func Int(key string, value int) LogField {
    return LogField{Key: key, Value: value}
}

func Float(key string, value float64) LogField {
    return LogField{Key: key, Value: value}
}

func Bool(key string, value bool) LogField {
    return LogField{Key: key, Value: value}
}

func Error(key string, err error) LogField {
    return LogField{Key: key, Value: err}
}
```

## Implementation Details

### Current Implementation Status

The logging system currently has:

1. **Basic Logger Implementation**: Simple logging with levels
2. **Secure Logger**: Sensitive data filtering
3. **Configuration Support**: HCL-based configuration
4. **Multiple Output Formats**: JSON and text formats

### Missing Features

1. **Structured Logging**: Limited structured logging support
2. **Log Handlers**: No pluggable log handlers
3. **Context Support**: Limited context propagation
4. **Performance Optimization**: No async logging
5. **Log Rotation**: No log rotation support

## Configuration

### Logging Configuration Structure

```hcl
# Global logging configuration
logging {
  level = "info"
  format = "text"
  
  # File output configuration
  file {
    path = "/var/log/spooky.log"
    max_size = "100MB"
    max_age = "30d"
    max_backups = 10
  }
  
  # Console output configuration
  console {
    enabled = true
    color = true
  }
  
  # Security configuration
  security {
    redact_sensitive = true
    sensitive_fields = ["password", "secret", "key", "token"]
  }
}
```

### Project-Level Logging

```hcl
# Project-specific logging configuration
project {
  name = "my-project"
  
  logging {
    level = "debug"
    format = "json"
    
    # Project-specific file output
    file {
      path = "./logs/project.log"
      max_size = "50MB"
    }
  }
}
```

## Usage Examples

### Basic Logging

```go
// Get a logger
logger := logManager.GetLogger("my-component")

// Log messages
logger.Info("Application started")
logger.Debug("Processing request", String("request_id", "123"))
logger.Error("Failed to connect", Error("error", err))
```

### Structured Logging

```go
// Create logger with fields
logger := logManager.GetLogger("database").
    WithField("component", "database").
    WithField("version", "1.0.0")

// Log with additional fields
logger.Info("Query executed",
    String("query", "SELECT * FROM users"),
    Int("duration_ms", 150),
    Int("rows", 1000))
```

### Context-Aware Logging

```go
// Create context
ctx := LogContext{
    RequestID: "req-123",
    UserID:    "user-456",
}

// Get logger with context
logger := logManager.GetLoggerWithContext("api", ctx)

// Log with context
logger.Info("Request processed",
    String("endpoint", "/api/users"),
    Int("status_code", 200))
```

### Error Logging

```go
// Log errors with context
if err != nil {
    logger.Error("Database operation failed",
        String("operation", "insert"),
        String("table", "users"),
        Error("error", err),
        String("sql", query))
}
```

## Security Features

### Sensitive Data Filtering

The logging system includes automatic filtering of sensitive data:

```go
// Sensitive data is automatically redacted
logger.Info("User login",
    String("username", "john.doe"),
    String("password", "secret123")) // Password will be redacted

// Output: User login username=john.doe password=[REDACTED]
```

### Secure Logger Implementation

```go
type SecureLogger struct {
    logger          Logger
    sensitiveFields []string
}

func (l *SecureLogger) Info(msg string, fields ...LogField) {
    l.logger.Info(msg, l.sanitizeFields(fields)...)
}

func (l *SecureLogger) sanitizeFields(fields []LogField) []LogField {
    sanitized := make([]LogField, len(fields))
    for i, field := range fields {
        if l.isSensitiveField(field.Key) {
            sanitized[i] = String(field.Key, "[REDACTED]")
        } else {
            sanitized[i] = field
        }
    }
    return sanitized
}
```

## Performance Considerations

### Current Limitations

1. **Synchronous Logging**: All logging operations are synchronous
2. **No Buffering**: No log entry buffering for performance
3. **Limited Async Support**: No asynchronous log processing
4. **Memory Usage**: No memory usage optimization

### Recommended Improvements

1. **Async Logging**: Implement asynchronous log processing
2. **Log Buffering**: Add log entry buffering
3. **Performance Monitoring**: Add logging performance metrics
4. **Memory Optimization**: Optimize memory usage for high-volume logging

## Integration with Other Systems

### SSH Integration

```go
// SSH operations use structured logging
sshLogger := logManager.GetLogger("ssh").
    WithField("component", "ssh")

sshLogger.Info("SSH connection established",
    String("host", hostname),
    String("user", username),
    Int("port", port))
```

### Facts Integration

```go
// Facts collection uses context-aware logging
factsLogger := logManager.GetLogger("facts").
    WithField("component", "facts")

factsLogger.Info("Facts collected",
    String("machine", machine),
    Int("fact_count", len(facts)),
    String("collection_time", duration.String()))
```

### Actions Integration

```go
// Action execution uses structured logging
actionLogger := logManager.GetLogger("actions").
    WithField("component", "actions")

actionLogger.Info("Action executed",
    String("action", actionName),
    String("machine", machine),
    Int("exit_code", exitCode),
    String("duration", duration.String()))
```

## Error Handling

### Logging Errors

```go
// Handle logging errors gracefully
if err := logger.Error("Failed to process request", Error("error", err)); err != nil {
    // Fallback to stderr if logging fails
    fmt.Fprintf(os.Stderr, "Logging failed: %v\n", err)
}
```

### Configuration Errors

```go
// Validate logging configuration
if err := validateLogConfig(config); err != nil {
    return fmt.Errorf("invalid logging configuration: %w", err)
}
```

## Testing

### Logger Testing

```go
func TestLogger(t *testing.T) {
    // Create test logger
    logger := NewTestLogger()
    
    // Test logging
    logger.Info("Test message")
    
    // Verify log entries
    entries := logger.GetEntries()
    assert.Len(t, entries, 1)
    assert.Equal(t, "Test message", entries[0].Message)
}
```

### Mock Logger

```go
type MockLogger struct {
    entries []LogEntry
    mutex   sync.RWMutex
}

func (l *MockLogger) Info(msg string, fields ...LogField) {
    l.mutex.Lock()
    defer l.mutex.Unlock()
    
    l.entries = append(l.entries, LogEntry{
        Timestamp: time.Now(),
        Level:     LogLevelInfo,
        Message:   msg,
        Fields:    fieldsToMap(fields),
    })
}

func (l *MockLogger) GetEntries() []LogEntry {
    l.mutex.RLock()
    defer l.mutex.RUnlock()
    
    result := make([]LogEntry, len(l.entries))
    copy(result, l.entries)
    return result
}
```

## Best Practices

### Log Level Usage

1. **Debug**: Detailed information for debugging
2. **Info**: General operational information
3. **Warn**: Warning conditions that don't stop operation
4. **Error**: Error conditions that affect operation
5. **Fatal**: Critical errors that require immediate attention

### Structured Logging

```go
// Good: Structured logging with fields
logger.Info("User action completed",
    String("user_id", userID),
    String("action", action),
    Int("duration_ms", duration),
    String("status", status))

// Bad: Unstructured logging
logger.Info(fmt.Sprintf("User %s completed action %s in %dms with status %s",
    userID, action, duration, status))
```

### Error Context

```go
// Good: Include error context
logger.Error("Database query failed",
    String("query", query),
    Error("error", err),
    String("table", tableName))

// Bad: Minimal error information
logger.Error("Database error: " + err.Error())
```

### Performance Logging

```go
// Good: Include performance metrics
start := time.Now()
// ... operation ...
duration := time.Since(start)

logger.Info("Operation completed",
    String("operation", operationName),
    Duration("duration", duration),
    Int("result_count", len(results)))
```

## Future Enhancements

### Planned Features

1. **Log Aggregation**: Support for log aggregation systems
2. **Metrics Integration**: Integration with metrics systems
3. **Log Sampling**: Configurable log sampling for high-volume scenarios
4. **Custom Formatters**: Pluggable log formatters
5. **Log Compression**: Automatic log compression
6. **Distributed Tracing**: Integration with distributed tracing systems

### Architecture Improvements

1. **Plugin System**: Pluggable log handlers and formatters
2. **Performance Optimization**: Async logging and buffering
3. **Configuration Validation**: Enhanced configuration validation
4. **Monitoring Integration**: Integration with monitoring systems
5. **Security Enhancements**: Enhanced security features

## Related Documentation

- [Logging User Guide](LOGGING_USER_GUIDE.md) - User guide for logging features
- [Logging Configuration Strategy](LOGGING_CONFIGURATION_STRATEGY.md) - Configuration strategies
- [Logging Framework](LOGGING_FRAMEWORK.md) - Framework architecture
- [Logging Troubleshooting](LOGGING_TROUBLESHOOTING.md) - Troubleshooting guide
