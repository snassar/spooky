# Logging System API Reference

## Overview

This document provides a comprehensive API reference for the spooky logging system. It covers all interfaces, types, methods, and implementation details for developers working with the logging system.

**Status: Production Ready** - The logging system is fully implemented with comprehensive configuration, formatters, and secure logging capabilities.

## Core Interfaces

### LoggingManager Interface

The `LoggingManager` interface provides the primary entry point for logging operations:

```go
type LoggingManager interface {
    // GetLogger returns a logger instance
    GetLogger() Logger
    
    // GetLoggerWithContext returns a logger with context
    GetLoggerWithContext(ctx context.Context) Logger
    
    // ConfigureLogging configures logging from configuration
    ConfigureLogging(ctx context.Context, config *spookytypeslogging.Config) error
    
    // SetLogLevel sets the global log level
    SetLogLevel(level string) error
    
    // GetLogLevel returns the current log level
    GetLogLevel() string
}
```

**Implementation Status**: ✅ **Fully Implemented** - Complete logging management functionality

### Logger Interface

The `Logger` interface provides logging operations:

```go
type Logger interface {
    // Basic logging methods
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    Fatal(msg string, fields ...Field)
    
    // Context-aware logging
    WithContext(ctx context.Context) Logger
    WithFields(fields ...Field) Logger
    
    // Structured logging
    WithField(key string, value interface{}) Logger
    WithError(err error) Logger
    
    // Sync flushes any buffered log entries
    Sync() error
}
```

**Implementation Status**: ✅ **Fully Implemented** - Complete logging functionality with structured logging

## Current Implementation Status

### ✅ Working Components

1. **Logging Configuration**: Comprehensive logging configuration system
2. **Multiple Formatters**: JSON, text, and custom formatters
3. **Multiple Outputs**: Console, file, and custom outputs
4. **Log Levels**: Full log level support with filtering
5. **Structured Logging**: Field-based structured logging
6. **Context Support**: Context-aware logging
7. **Secure Logging**: Sensitive data filtering and redaction
8. **CLI Integration**: `spooky logging` commands with configuration
9. **Project Integration**: Project-specific logging configuration
10. **Global Integration**: Global logging configuration
11. **Auto-Setup**: Automatic logging configuration setup
12. **Validation**: Logging configuration validation

### ✅ Advanced Features

1. **Formatter System**: Pluggable log formatters
2. **Output System**: Pluggable log outputs
3. **Filter System**: Log filtering capabilities
4. **Performance**: Optimized logging performance
5. **Thread Safety**: Thread-safe logging operations
6. **Error Handling**: Comprehensive error handling
7. **Configuration Validation**: Logging configuration validation
8. **Default Configuration**: Sensible default logging configuration

### ✅ Integration Features

1. **CLI Commands**: Complete CLI integration
2. **Project Support**: Project-specific logging
3. **Global Support**: Global logging configuration
4. **Auto-Setup**: Automatic configuration setup
5. **Schema Validation**: HCL schema validation
6. **Configuration Loading**: Configuration file loading

## Implementation Details

### Logging Configuration System

The logging system supports comprehensive configuration:

```go
type LoggingConfig struct {
    // Global logging settings
    Level   string `json:"level" hcl:"level"`
    Format  string `json:"format" hcl:"format"`
    Output  string `json:"output" hcl:"output"`
    
    // File output settings
    File *FileOutputConfig `json:"file,omitempty" hcl:"file,optional"`
    
    // Console output settings
    Console *ConsoleOutputConfig `json:"console,omitempty" hcl:"console,optional"`
    
    // Formatter settings
    Formatters map[string]*FormatterConfig `json:"formatters,omitempty" hcl:"formatters,optional"`
    
    // Output settings
    Outputs map[string]*OutputConfig `json:"outputs,omitempty" hcl:"outputs,optional"`
    
    // Security settings
    Security *SecurityConfig `json:"security,omitempty" hcl:"security,optional"`
    
    // Performance settings
    Performance *PerformanceConfig `json:"performance,omitempty" hcl:"performance,optional"`
}

type FileOutputConfig struct {
    Path       string `json:"path" hcl:"path"`
    MaxSize    int    `json:"max_size" hcl:"max_size"`
    MaxBackups int    `json:"max_backups" hcl:"max_backups"`
    MaxAge     int    `json:"max_age" hcl:"max_age"`
    Compress   bool   `json:"compress" hcl:"compress"`
}

type ConsoleOutputConfig struct {
    Colorize bool `json:"colorize" hcl:"colorize"`
    Timestamp bool `json:"timestamp" hcl:"timestamp"`
}

type FormatterConfig struct {
    Type   string                 `json:"type" hcl:"type"`
    Config map[string]interface{} `json:"config,omitempty" hcl:"config,optional"`
}

type OutputConfig struct {
    Type   string                 `json:"type" hcl:"type"`
    Config map[string]interface{} `json:"config,omitempty" hcl:"config,optional"`
}

type SecurityConfig struct {
    RedactSensitive bool     `json:"redact_sensitive" hcl:"redact_sensitive"`
    SensitiveFields []string `json:"sensitive_fields" hcl:"sensitive_fields"`
    MaskPatterns    []string `json:"mask_patterns" hcl:"mask_patterns"`
}

type PerformanceConfig struct {
    BufferSize int `json:"buffer_size" hcl:"buffer_size"`
    Async      bool `json:"async" hcl:"async"`
}
```

### Logging Manager Implementation

The logging manager provides centralized logging control:

```go
type LoggingManager struct {
    logger  Logger
    config  *spookytypeslogging.Config
    mutex   sync.RWMutex
}

func NewLoggingManager() *LoggingManager {
    return &LoggingManager{
        logger: NewDefaultLogger(),
        config: &spookytypeslogging.Config{},
    }
}

func (m *LoggingManager) ConfigureLogging(ctx context.Context, config *spookytypeslogging.Config) error {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    // Validate configuration
    if err := m.validateConfig(config); err != nil {
        return fmt.Errorf("invalid logging configuration: %w", err)
    }
    
    // Create new logger with configuration
    logger, err := m.createLogger(config)
    if err != nil {
        return fmt.Errorf("failed to create logger: %w", err)
    }
    
    m.logger = logger
    m.config = config
    
    return nil
}

func (m *LoggingManager) GetLogger() Logger {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    return m.logger
}

func (m *LoggingManager) GetLoggerWithContext(ctx context.Context) Logger {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    return m.logger.WithContext(ctx)
}
```

### Secure Logging Implementation

The secure logging system filters sensitive information:

```go
type SecureLogger struct {
    logger          Logger
    sensitiveFields []string
    maskPatterns    []*regexp.Regexp
}

func NewSecureLogger(logger Logger, config *SecurityConfig) *SecureLogger {
    var patterns []*regexp.Regexp
    for _, pattern := range config.MaskPatterns {
        if re, err := regexp.Compile(pattern); err == nil {
            patterns = append(patterns, re)
        }
    }
    
    return &SecureLogger{
        logger:          logger,
        sensitiveFields: config.SensitiveFields,
        maskPatterns:    patterns,
    }
}

func (l *SecureLogger) Info(msg string, fields ...Field) {
    l.logger.Info(msg, l.sanitizeFields(fields)...)
}

func (l *SecureLogger) Error(msg string, fields ...Field) {
    l.logger.Error(msg, l.sanitizeFields(fields)...)
}

func (l *SecureLogger) sanitizeFields(fields []Field) []Field {
    sanitized := make([]Field, len(fields))
    for i, field := range fields {
        if l.isSensitiveField(field.Key) {
            sanitized[i] = String(field.Key, "[REDACTED]")
        } else {
            sanitized[i] = l.sanitizeValue(field)
        }
    }
    return sanitized
}

func (l *SecureLogger) isSensitiveField(key string) bool {
    keyLower := strings.ToLower(key)
    for _, sensitive := range l.sensitiveFields {
        if strings.Contains(keyLower, sensitive) {
            return true
        }
    }
    return false
}

func (l *SecureLogger) sanitizeValue(field Field) Field {
    if str, ok := field.Value.(string); ok {
        for _, pattern := range l.maskPatterns {
            str = pattern.ReplaceAllString(str, "[MASKED]")
        }
        return String(field.Key, str)
    }
    return field
}
```

## Type Definitions

### Logging Types

```go
// LogLevel represents logging levels
type LogLevel string

const (
    LogLevelDebug LogLevel = "debug"
    LogLevelInfo  LogLevel = "info"
    LogLevelWarn  LogLevel = "warn"
    LogLevelError LogLevel = "error"
    LogLevelFatal LogLevel = "fatal"
)

// Field represents a structured logging field
type Field struct {
    Key   string
    Value interface{}
}

// LogEntry represents a log entry
type LogEntry struct {
    Timestamp time.Time              `json:"timestamp"`
    Level     LogLevel               `json:"level"`
    Message   string                 `json:"message"`
    Fields    map[string]interface{} `json:"fields,omitempty"`
    Context   map[string]interface{} `json:"context,omitempty"`
    Error     string                 `json:"error,omitempty"`
}

// LoggingContext provides context for logging operations
type LoggingContext struct {
    // Project path
    ProjectPath string `json:"project_path" hcl:"project_path"`
    
    // Logging configuration
    Config *LoggingConfig `json:"config" hcl:"config"`
    
    // Operation timestamp
    Timestamp time.Time `json:"timestamp" hcl:"timestamp"`
    
    // Operation metadata
    Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// LoggingResult represents the result of logging operations
type LoggingResult struct {
    // Logging context
    Context *LoggingContext `json:"context" hcl:"context"`
    
    // Operation success
    Success bool `json:"success" hcl:"success"`
    
    // Log entries written
    EntriesWritten int `json:"entries_written" hcl:"entries_written"`
    
    // Operation error
    Error string `json:"error,omitempty" hcl:"error,optional"`
    
    // Operation duration
    Duration time.Duration `json:"duration" hcl:"duration"`
}
```

### Configuration Types

```go
// LoggingConfig represents logging configuration
type LoggingConfig struct {
    // Global settings
    Level   string `json:"level" hcl:"level"`
    Format  string `json:"format" hcl:"format"`
    Output  string `json:"output" hcl:"output"`
    
    // File output configuration
    File *FileOutputConfig `json:"file,omitempty" hcl:"file,optional"`
    
    // Console output configuration
    Console *ConsoleOutputConfig `json:"console,omitempty" hcl:"console,optional"`
    
    // Formatter configurations
    Formatters map[string]*FormatterConfig `json:"formatters,omitempty" hcl:"formatters,optional"`
    
    // Output configurations
    Outputs map[string]*OutputConfig `json:"outputs,omitempty" hcl:"outputs,optional"`
    
    // Security configuration
    Security *SecurityConfig `json:"security,omitempty" hcl:"security,optional"`
    
    // Performance configuration
    Performance *PerformanceConfig `json:"performance,omitempty" hcl:"performance,optional"`
}

// FileOutputConfig represents file output configuration
type FileOutputConfig struct {
    Path       string `json:"path" hcl:"path"`
    MaxSize    int    `json:"max_size" hcl:"max_size"`
    MaxBackups int    `json:"max_backups" hcl:"max_backups"`
    MaxAge     int    `json:"max_age" hcl:"max_age"`
    Compress   bool   `json:"compress" hcl:"compress"`
}

// ConsoleOutputConfig represents console output configuration
type ConsoleOutputConfig struct {
    Colorize  bool `json:"colorize" hcl:"colorize"`
    Timestamp bool `json:"timestamp" hcl:"timestamp"`
}

// FormatterConfig represents formatter configuration
type FormatterConfig struct {
    Type   string                 `json:"type" hcl:"type"`
    Config map[string]interface{} `json:"config,omitempty" hcl:"config,optional"`
}

// OutputConfig represents output configuration
type OutputConfig struct {
    Type   string                 `json:"type" hcl:"type"`
    Config map[string]interface{} `json:"config,omitempty" hcl:"config,optional"`
}

// SecurityConfig represents security configuration
type SecurityConfig struct {
    RedactSensitive bool     `json:"redact_sensitive" hcl:"redact_sensitive"`
    SensitiveFields []string `json:"sensitive_fields" hcl:"sensitive_fields"`
    MaskPatterns    []string `json:"mask_patterns" hcl:"mask_patterns"`
}

// PerformanceConfig represents performance configuration
type PerformanceConfig struct {
    BufferSize int  `json:"buffer_size" hcl:"buffer_size"`
    Async      bool `json:"async" hcl:"async"`
}
```

## Error Handling

### Logging Errors

```go
// LoggingError represents logging operation errors
type LoggingError struct {
    Operation string `json:"operation" hcl:"operation"`
    Error     string `json:"error" hcl:"error"`
    Details   string `json:"details,omitempty" hcl:"details,optional"`
}

// LoggingValidationError represents logging validation errors
type LoggingValidationError struct {
    Field   string `json:"field" hcl:"field"`
    Message string `json:"message" hcl:"message"`
    Value   string `json:"value,omitempty" hcl:"value,optional"`
}
```

### Validation Implementation

```go
// ValidateLoggingConfig validates logging configuration
func ValidateLoggingConfig(config *LoggingConfig) error {
    if config == nil {
        return fmt.Errorf("logging configuration cannot be nil")
    }
    
    // Validate log level
    validLevels := []string{"debug", "info", "warn", "error", "fatal"}
    valid := false
    for _, level := range validLevels {
        if config.Level == level {
            valid = true
            break
        }
    }
    if !valid {
        return fmt.Errorf("invalid log level: %s (valid levels: %v)", config.Level, validLevels)
    }
    
    // Validate format
    validFormats := []string{"json", "text", "custom"}
    valid = false
    for _, format := range validFormats {
        if config.Format == format {
            valid = true
            break
        }
    }
    if !valid {
        return fmt.Errorf("invalid log format: %s (valid formats: %v)", config.Format, validFormats)
    }
    
    // Validate output
    validOutputs := []string{"console", "file", "custom"}
    valid = false
    for _, output := range validOutputs {
        if config.Output == output {
            valid = true
            break
        }
    }
    if !valid {
        return fmt.Errorf("invalid log output: %s (valid outputs: %v)", config.Output, validOutputs)
    }
    
    // Validate file output configuration
    if config.File != nil {
        if err := validateFileOutputConfig(config.File); err != nil {
            return fmt.Errorf("invalid file output configuration: %w", err)
        }
    }
    
    // Validate security configuration
    if config.Security != nil {
        if err := validateSecurityConfig(config.Security); err != nil {
            return fmt.Errorf("invalid security configuration: %w", err)
        }
    }
    
    return nil
}
```

## CLI Commands

### Logging Configuration Command

```bash
# Show current logging configuration
spooky logging config

# Show logging configuration for a project
spooky logging config ./my-project

# Show global logging configuration
spooky logging config --global
```

### Logging Validation Command

```bash
# Validate logging configuration
spooky logging validate

# Validate logging configuration for a project
spooky logging validate ./my-project

# Validate global logging configuration
spooky logging validate --global
```

### Logging Test Command

```bash
# Test logging configuration
spooky logging test

# Test logging configuration for a project
spooky logging test ./my-project

# Test logging with specific level
spooky logging test --level debug
```

## Integration Examples

### Basic Logging Configuration

```hcl
# logging.hcl
logging {
  level = "info"
  format = "json"
  output = "console"
  
  console {
    colorize = true
    timestamp = true
  }
  
  security {
    redact_sensitive = true
    sensitive_fields = ["password", "secret", "key", "token"]
    mask_patterns = [
      "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b",
      "\\b\\d{4}[- ]?\\d{4}[- ]?\\d{4}[- ]?\\d{4}\\b"
    ]
  }
  
  performance {
    buffer_size = 1024
    async = true
  }
}
```

### Project-Specific Logging

```hcl
# project/logging.hcl
logging {
  level = "debug"
  format = "text"
  output = "file"
  
  file {
    path = "logs/project.log"
    max_size = 100
    max_backups = 3
    max_age = 30
    compress = true
  }
  
  security {
    redact_sensitive = true
    sensitive_fields = ["api_key", "database_password"]
  }
}
```

### Logging Integration

```go
// Logging integration example
func setupLogging(projectPath string) error {
    ctx := context.Background()
    
    // Create logging manager
    manager := spookylogging.NewLoggingManager()
    
    // Load logging configuration
    config, err := loadLoggingConfig(projectPath)
    if err != nil {
        return fmt.Errorf("failed to load logging configuration: %w", err)
    }
    
    // Configure logging
    if err := manager.ConfigureLogging(ctx, config); err != nil {
        return fmt.Errorf("failed to configure logging: %w", err)
    }
    
    // Get logger
    logger := manager.GetLogger()
    
    // Test logging
    logger.Info("Logging system initialized", 
        String("project", projectPath),
        String("level", config.Level),
        String("format", config.Format),
    )
    
    return nil
}

// Structured logging example
func logOperation(ctx context.Context, operation string, details map[string]interface{}) {
    logger := spookylogging.GetLogger().WithContext(ctx)
    
    fields := []Field{
        String("operation", operation),
        String("timestamp", time.Now().Format(time.RFC3339)),
    }
    
    for key, value := range details {
        fields = append(fields, Any(key, value))
    }
    
    logger.Info("Operation completed", fields...)
}
```

### Secure Logging Example

```go
// Secure logging example
func logSensitiveData(ctx context.Context, userData map[string]interface{}) {
    logger := spookylogging.GetLogger().WithContext(ctx)
    
    // Log with sensitive data (will be redacted)
    logger.Info("User data processed",
        String("user_id", userData["id"].(string)),
        String("email", userData["email"].(string)),
        String("password", userData["password"].(string)), // Will be redacted
    )
    
    // Log without sensitive data
    logger.Info("User authentication successful",
        String("user_id", userData["id"].(string)),
        String("timestamp", time.Now().Format(time.RFC3339)),
    )
}
```

## Current Capabilities

### Configuration Management

1. **Global Configuration**: System-wide logging configuration
2. **Project Configuration**: Project-specific logging configuration
3. **Auto-Setup**: Automatic configuration file creation
4. **Validation**: Configuration validation and error reporting
5. **Schema Support**: HCL schema validation
6. **Default Values**: Sensible default configuration

### Output Management

1. **Console Output**: Colored console output with timestamps
2. **File Output**: Rotating file output with compression
3. **Custom Outputs**: Pluggable output system
4. **Multiple Outputs**: Support for multiple simultaneous outputs
5. **Performance**: Optimized output performance

### Formatting System

1. **JSON Format**: Structured JSON logging
2. **Text Format**: Human-readable text logging
3. **Custom Formats**: Pluggable formatter system
4. **Field Support**: Structured field logging
5. **Context Support**: Context-aware formatting

### Security Features

1. **Sensitive Data Redaction**: Automatic redaction of sensitive fields
2. **Pattern Masking**: Regex-based pattern masking
3. **Field Filtering**: Configurable sensitive field filtering
4. **Secure Defaults**: Secure default configuration
5. **Audit Support**: Audit logging capabilities

### Performance Features

1. **Buffered Output**: Configurable output buffering
2. **Async Logging**: Asynchronous logging support
3. **Level Filtering**: Efficient level-based filtering
4. **Memory Management**: Optimized memory usage
5. **Thread Safety**: Thread-safe logging operations

## Summary

The logging system is fully implemented and production-ready with comprehensive configuration, multiple output formats, secure logging capabilities, and excellent performance characteristics. It provides all the features needed for production logging in the spooky system.

**Status**: ✅ **Production Ready** - Complete logging system with advanced features and excellent performance.
