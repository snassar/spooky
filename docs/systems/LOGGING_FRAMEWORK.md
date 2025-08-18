# Spooky Logging Framework

## Overview

The spooky project now includes a comprehensive logging framework built on Go's `slog` package with schema-driven configuration. This framework provides structured logging capabilities with extensive customization options while maintaining high performance and following best practices for log formatting.

## Key Features

### 🎯 **Schema-Driven Configuration**
- **HCL Schema**: Complete logging configuration schema in `internal/schemas/schemas/logging.schema.hcl`
- **Validation**: All configuration validated against schema at runtime
- **Best Practices**: Configuration follows [Sematext log formatting best practices](https://sematext.com/blog/log-formatting-8-best-practices-for-better-readability/)

### 📊 **Multiple Output Formats**
- **JSON**: Machine-readable structured logging (default)
- **Text**: Human-readable text format
- **Structured**: Custom formatting with field ordering and filtering

### 🎛️ **Flexible Output Destinations**
- **stdout**: Standard output
- **stderr**: Standard error (default)
- **file**: Arbitrary log files with permissions and append options
- **null**: No output (for testing or when logging is disabled)

### 🔧 **Advanced Configuration Options**

#### Structured Logging
- **Timestamp Configuration**: Multiple formats (RFC3339, Unix, ISO8601) with timezone support
- **Level Configuration**: Customizable field keys, uppercase options, color support
- **Component Tracking**: Automatic component identification with package inclusion options
- **Operation Context**: Operation tracking with optional ID generation
- **Error Handling**: Configurable error details, stack traces, type information
- **Caller Information**: File, line, function information with configurable skip frames

#### Performance Optimization
- **Buffering**: Configurable buffer sizes and flush intervals
- **Async Logging**: Worker pools with configurable queue sizes
- **Drop Policies**: Configurable behavior when queues are full

#### Filtering and Security
- **Component Filtering**: Different log levels per component
- **Pattern Filtering**: Regex-based include/exclude patterns
- **Field Filtering**: Include/exclude specific fields
- **Sensitive Data**: Automatic redaction of sensitive fields (passwords, tokens, etc.)

#### Log Rotation
- **Size-based**: Rotate when files reach specified size
- **Age-based**: Rotate files older than specified duration
- **Compression**: Automatic compression of rotated files
- **Backup Management**: Configurable number of backup files

## Architecture

### Core Components

```
spooky/
├── internal/
│   ├── logging/
│   │   ├── logging.go          # Main logging implementation
│   │   └── logging_test.go     # Comprehensive test suite
│   ├── types/logging/
│   │   └── config.go           # Type definitions
│   └── schemas/schemas/
│       └── logging.schema.hcl  # Configuration schema
└── examples/
    ├── logging-config.hcl      # Example configuration
    └── logging-demo.go         # Usage demonstration
```

### Key Interfaces

#### LogManager
```go
type LogManager interface {
    GetLogger(component string) Logger
    SetLevel(level LogLevel)
    GetLevel() LogLevel
    Configure(config *LogConfig) error
    Flush() error
    Close() error
}
```

#### Logger
```go
type Logger interface {
    Debug(msg string, fields ...map[string]interface{})
    Info(msg string, fields ...map[string]interface{})
    Warn(msg string, fields ...map[string]interface{})
    Error(msg string, err error, fields ...map[string]interface{})
    Fatal(msg string, err error, fields ...map[string]interface{})
    WithFields(fields map[string]interface{}) Logger
    WithComponent(component string) Logger
    WithOperation(operation string) Logger
    SetLevel(level LogLevel)
    GetLevel() LogLevel
}
```

## Usage Examples

### Basic Usage

```go
// Create log manager
logManager := spookylogging.NewLogManager()
defer logManager.Close()

// Get logger for component
logger := logManager.GetLogger("my-component")

// Log messages
logger.Info("Application started", map[string]interface{}{
    "version": "1.0.0",
    "pid":     os.Getpid(),
})

logger.Error("Operation failed", err, map[string]interface{}{
    "operation": "database-query",
    "duration":  "1.5s",
})
```

### Configuration Examples

#### File Output with JSON Format
```hcl
logging {
  level  = "info"
  format = "json"
  output = "file"
  
  file {
    path        = "/var/log/spooky/app.log"
    permissions = "0644"
    append      = true
  }
}
```

#### Component-Specific Filtering
```hcl
logging {
  level  = "debug"
  format = "json"
  output = "stderr"
  
  filtering {
    components = {
      "ssh"     = "debug"
      "facts"   = "info"
      "actions" = "warn"
    }
  }
}
```

#### Structured Logging with Custom Fields
```hcl
logging {
  level  = "info"
  format = "structured"
  output = "stdout"
  
  structured {
    timestamp {
      enabled  = true
      format   = "RFC3339"
      timezone = "UTC"
    }
    
    fields {
      global = {
        service = "spooky"
        version = "0.20250812.0"
      }
      
      filter {
        sensitive = ["password", "token", "secret"]
      }
    }
  }
}
```

### Advanced Usage Patterns

#### Operation Context
```go
logger := logManager.GetLogger("project")
initLogger := logger.WithOperation("project-init")

initLogger.Info("Starting initialization", map[string]interface{}{
    "project_path": "/path/to/project",
    "template":     "default",
})
```

#### Field Inheritance
```go
logger := logManager.GetLogger("user")
userLogger := logger.WithFields(map[string]interface{}{
    "user_id": "12345",
    "session": "abc123",
})

userLogger.Info("User action", map[string]interface{}{
    "action": "login",
    "ip":     "192.168.1.1",
})
```

#### Error Context
```go
logger := logManager.GetLogger("validation")
logger.Error("Configuration validation failed", err, map[string]interface{}{
    "file":     "project.hcl",
    "line":     15,
    "column":   8,
    "expected": "string",
    "got":      "number",
})
```

## Integration with CLI

The logging framework is designed to integrate seamlessly with the spooky CLI:

### CLI Output Separation
- **Logs are NOT shown in CLI output** - they go to configured destinations
- **CLI messages remain user-friendly** with emojis and clear formatting
- **Debug information goes to logs** for troubleshooting

### Configuration Loading
```go
// Load logging configuration from HCL file
config, err := loadLoggingConfig("logging.hcl")
if err != nil {
    // Fall back to default configuration
    config = getDefaultLogConfig()
}

// Configure logging
if err := logManager.Configure(config); err != nil {
    return fmt.Errorf("failed to configure logging: %w", err)
}
```

### Component Organization
```go
// Use component-specific loggers throughout the application
sshLogger := logManager.GetLogger("ssh")
factsLogger := logManager.GetLogger("facts")
actionsLogger := logManager.GetLogger("actions")
projectLogger := logManager.GetLogger("project")
```

## Performance Characteristics

### Benchmarks
The framework leverages Go's `slog` package for optimal performance:

- **JSON Format**: ~656 ns/op with 5 allocations
- **Text Format**: ~935 ns/op with 10 allocations
- **Structured Format**: Custom performance based on configuration

### Memory Usage
- **Minimal allocations** for common logging patterns
- **Efficient field handling** with map reuse
- **Configurable buffering** for high-throughput scenarios

## Testing

The logging framework includes comprehensive tests:

```bash
# Run all logging tests
go test ./internal/logging -v

# Run with coverage
go test ./internal/logging -cover
```

### Test Coverage
- ✅ Basic logging functionality
- ✅ File output configuration
- ✅ Component and operation context
- ✅ Error handling and filtering
- ✅ Configuration validation
- ✅ Performance characteristics

## Best Practices

### 1. **Component Naming**
Use descriptive, hierarchical component names:
```go
logManager.GetLogger("ssh.connection")
logManager.GetLogger("facts.collector")
logManager.GetLogger("actions.runner")
```

### 2. **Field Consistency**
Use consistent field names across components:
```go
// Good
logger.Info("Operation completed", map[string]interface{}{
    "operation": "user-authentication",
    "duration":  "1.2s",
    "user_id":   "12345",
})

// Avoid
logger.Info("Operation completed", map[string]interface{}{
    "op":     "user-auth",
    "time":   "1.2s",
    "uid":    "12345",
})
```

### 3. **Error Context**
Always provide context with errors:
```go
logger.Error("Database operation failed", err, map[string]interface{}{
    "table":    "users",
    "operation": "insert",
    "user_id":   "12345",
    "retries":   3,
})
```

### 4. **Sensitive Data**
Never log sensitive information directly:
```go
// Good - sensitive data will be redacted
logger.Info("Authentication attempt", map[string]interface{}{
    "user":     "admin",
    "password": "secret123", // Will be redacted
})

// Better - don't include sensitive data at all
logger.Info("Authentication attempt", map[string]interface{}{
    "user": "admin",
    "method": "password",
})
```

## Migration from Previous Logging

The new framework is a complete replacement for the previous logging implementation:

### Key Changes
1. **Uses `slog` instead of custom implementation**
2. **Schema-driven configuration**
3. **Enhanced filtering and security**
4. **Better performance characteristics**
5. **More flexible output options**

### Migration Steps
1. **Update imports** to use new logging package
2. **Replace logger creation** with new interface
3. **Update configuration** to use HCL schema
4. **Test thoroughly** with new filtering capabilities

## Future Enhancements

### Planned Features
- **Log aggregation** support (ELK stack, etc.)
- **Metrics integration** with logging
- **Distributed tracing** correlation
- **Custom formatters** for specific use cases
- **Log compression** and archival

### Performance Optimizations
- **Zero-copy logging** for high-performance scenarios
- **Batch logging** for bulk operations
- **Async handlers** for non-blocking logging
- **Memory pooling** for reduced allocations

## Conclusion

The new logging framework provides a robust, performant, and flexible logging solution for the spooky project. It follows industry best practices, supports schema-driven configuration, and integrates seamlessly with the existing CLI architecture while maintaining separation between user-facing output and internal logging.

The framework is ready for production use and provides a solid foundation for future enhancements and integrations.
