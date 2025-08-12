# Logging System API Reference

## Overview

This document provides a technical reference for the spooky logging system APIs and implementation details. It covers core interfaces, type definitions, implementation patterns, and integration details for developers working with the logging system.

## Core Interfaces

### LogManager Interface

The `LogManager` interface provides the primary entry point for logging operations:

```go
type LogManager interface {
    // GetLogger returns a logger instance for the specified component
    GetLogger(component string) Logger
    
    // GetLoggerWithContext returns a logger with additional context
    GetLoggerWithContext(component string, ctx map[string]interface{}) Logger
    
    // ConfigureLogging configures the logging system
    ConfigureLogging(config *spookytypes.LoggingConfig) error
    
    // ValidateConfiguration validates logging configuration
    ValidateConfiguration(config *spookytypes.LoggingConfig) (*spookytypes.ValidationResult, error)
    
    // TestLogging tests the current logging configuration
    TestLogging(level string, component string) error
    
    // GetConfiguration returns the current logging configuration
    GetConfiguration() *spookytypes.LoggingConfig
    
    // SetLogLevel sets the log level for a component
    SetLogLevel(component string, level string) error
    
    // EnableComponent enables logging for a component
    EnableComponent(component string, enabled bool) error
}
```

### Logger Interface

The `Logger` interface provides the core logging functionality:

```go
type Logger interface {
    // Basic logging methods
    Trace(msg string, fields ...interface{})
    Debug(msg string, fields ...interface{})
    Info(msg string, fields ...interface{})
    Warn(msg string, fields ...interface{})
    Error(msg string, fields ...interface{})
    Fatal(msg string, fields ...interface{})
    
    // Context-aware logging
    With(fields ...interface{}) Logger
    
    // Structured logging with explicit fields
    Log(level string, msg string, fields map[string]interface{})
    
    // Performance logging
    TimeOperation(operation string, fn func() error) error
    TimeOperationWithResult(operation string, fn func() (interface{}, error)) (interface{}, error)
    
    // Audit logging
    Audit(event string, fields map[string]interface{})
    
    // Component information
    GetComponent() string
    GetLevel() string
    IsEnabled() bool
}
```

### LoggingIntegration Interface

The `LoggingIntegration` interface provides integration capabilities:

```go
type LoggingIntegration interface {
    // ConfigureLogging configures the logging system
    ConfigureLogging(ctx context.Context, config *spookytypes.LoggingConfig) error
    
    // ValidateConfiguration validates logging configuration
    ValidateConfiguration(ctx context.Context, config *spookytypes.LoggingConfig) (*spookytypes.ValidationResult, error)
    
    // TestLogging tests logging configuration
    TestLogging(ctx context.Context, level string, component string) error
    
    // GetLogger returns a logger for the specified component
    GetLogger(ctx context.Context, component string) (Logger, error)
    
    // SetGlobalLevel sets the global log level
    SetGlobalLevel(ctx context.Context, level string) error
    
    // SetComponentLevel sets the log level for a specific component
    SetComponentLevel(ctx context.Context, component string, level string) error
}
```

## Type Definitions

### LoggingConfig

The main configuration structure for the logging system:

```go
type LoggingConfig struct {
    // Global logging level
    Level string `hcl:"level,optional"`
    
    // Output format (json, structured, plain)
    Format string `hcl:"format,optional"`
    
    // Output destinations
    Outputs []LoggingOutput `hcl:"output,block"`
    
    // Component-specific configuration
    Components map[string]ComponentConfig `hcl:"components,block"`
    
    // Performance logging configuration
    Performance *PerformanceConfig `hcl:"performance,block"`
    
    // Audit logging configuration
    Audit *AuditConfig `hcl:"audit,block"`
    
    // Common entity fields
    spookytypescommon.TimestampedEntity
    spookytypescommon.NamedEntity
    spookytypescommon.MetadataEntity
}
```

### LoggingOutput

Configuration for output destinations:

```go
type LoggingOutput struct {
    // Output type (console, file)
    Type string `hcl:"type"`
    
    // Whether this output is enabled
    Enabled bool `hcl:"enabled,optional"`
    
    // Log level for this output
    Level string `hcl:"level,optional"`
    
    // Console-specific options
    Console *ConsoleOutput `hcl:"console,block"`
    
    // File-specific options
    File *FileOutput `hcl:"file,block"`
    
    // Common entity fields
    spookytypescommon.TimestampedEntity
}
```

### ConsoleOutput

Configuration for console output:

```go
type ConsoleOutput struct {
    // Enable colorization
    Colorize bool `hcl:"colorize,optional"`
    
    // Use stderr for errors
    UseStderr bool `hcl:"use_stderr,optional"`
    
    // Enable timestamps
    IncludeTimestamp bool `hcl:"include_timestamp,optional"`
    
    // Enable log levels
    IncludeLevel bool `hcl:"include_level,optional"`
    
    // Enable component names
    IncludeComponent bool `hcl:"include_component,optional"`
}
```

### FileOutput

Configuration for file output:

```go
type FileOutput struct {
    // File path
    Path string `hcl:"path"`
    
    // Maximum file size
    MaxSize string `hcl:"max_size,optional"`
    
    // Maximum age
    MaxAge string `hcl:"max_age,optional"`
    
    // Maximum number of backups
    MaxBackups int `hcl:"max_backups,optional"`
    
    // Enable compression
    Compress bool `hcl:"compress,optional"`
    
    // File permissions
    Permissions string `hcl:"permissions,optional"`
}
```

### ComponentConfig

Configuration for component-specific logging:

```go
type ComponentConfig struct {
    // Log level for this component
    Level string `hcl:"level,optional"`
    
    // Whether logging is enabled
    Enabled bool `hcl:"enabled,optional"`
    
    // Component-specific outputs
    Outputs []LoggingOutput `hcl:"output,block"`
    
    // Common entity fields
    spookytypescommon.TimestampedEntity
    spookytypescommon.NamedEntity
}
```

### PerformanceConfig

Configuration for performance logging:

```go
type PerformanceConfig struct {
    // Enable performance logging
    Enabled bool `hcl:"enabled,optional"`
    
    // Threshold for logging slow operations (in milliseconds)
    ThresholdMs int `hcl:"threshold_ms,optional"`
    
    // Include memory usage
    IncludeMemory bool `hcl:"include_memory,optional"`
    
    // Include CPU usage
    IncludeCPU bool `hcl:"include_cpu,optional"`
}
```

### AuditConfig

Configuration for audit logging:

```go
type AuditConfig struct {
    // Enable audit logging
    Enabled bool `hcl:"enabled,optional"`
    
    // Audit log level
    Level string `hcl:"level,optional"`
    
    // Include user information
    IncludeUser bool `hcl:"include_user,optional"`
    
    // Include IP address
    IncludeIP bool `hcl:"include_ip,optional"`
    
    // Include session information
    IncludeSession bool `hcl:"include_session,optional"`
}
```

## Implementation Details

### LogManager Implementation

The `LogManager` implementation provides the core logging functionality:

```go
type Manager struct {
    config     *spookytypes.LoggingConfig
    loggers    map[string]Logger
    outputs    []LoggingOutput
    mutex      sync.RWMutex
    logger     spookylogging.Logger
}

func NewManager(config *spookytypes.LoggingConfig) LogManager {
    return &Manager{
        config:  config,
        loggers: make(map[string]Logger),
        outputs: config.Outputs,
    }
}

func (m *Manager) GetLogger(component string) Logger {
    m.mutex.RLock()
    if logger, exists := m.loggers[component]; exists {
        m.mutex.RUnlock()
        return logger
    }
    m.mutex.RUnlock()
    
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    // Create new logger
    logger := NewLogger(component, m.config, m.outputs)
    m.loggers[component] = logger
    return logger
}
```

### Logger Implementation

The `Logger` implementation handles individual component logging:

```go
type logger struct {
    component string
    config    *spookytypes.LoggingConfig
    outputs   []LoggingOutput
    level     string
    enabled   bool
    fields    map[string]interface{}
    mutex     sync.RWMutex
}

func NewLogger(component string, config *spookytypes.LoggingConfig, outputs []LoggingOutput) Logger {
    level := config.Level
    enabled := true
    
    // Check component-specific configuration
    if compConfig, exists := config.Components[component]; exists {
        if compConfig.Level != "" {
            level = compConfig.Level
        }
        enabled = compConfig.Enabled
    }
    
    return &logger{
        component: component,
        config:    config,
        outputs:   outputs,
        level:     level,
        enabled:   enabled,
        fields:    make(map[string]interface{}),
    }
}

func (l *logger) Info(msg string, fields ...interface{}) {
    if !l.isLevelEnabled("info") {
        return
    }
    
    logEntry := l.createLogEntry("info", msg, fields...)
    l.writeToOutputs(logEntry)
}

func (l *logger) With(fields ...interface{}) Logger {
    newLogger := &logger{
        component: l.component,
        config:    l.config,
        outputs:   l.outputs,
        level:     l.level,
        enabled:   l.enabled,
        fields:    make(map[string]interface{}),
    }
    
    // Copy existing fields
    for k, v := range l.fields {
        newLogger.fields[k] = v
    }
    
    // Add new fields
    for i := 0; i < len(fields); i += 2 {
        if i+1 < len(fields) {
            key, ok := fields[i].(string)
            if ok {
                newLogger.fields[key] = fields[i+1]
            }
        }
    }
    
    return newLogger
}
```

### Output Handlers

#### Console Output Handler

```go
type ConsoleOutputHandler struct {
    config ConsoleOutput
    writer io.Writer
}

func NewConsoleOutputHandler(config ConsoleOutput) *ConsoleOutputHandler {
    writer := os.Stdout
    if config.UseStderr {
        writer = os.Stderr
    }
    
    return &ConsoleOutputHandler{
        config: config,
        writer: writer,
    }
}

func (h *ConsoleOutputHandler) Write(entry LogEntry) error {
    var output string
    
    switch h.config.Format {
    case "json":
        output = h.formatJSON(entry)
    case "structured":
        output = h.formatStructured(entry)
    default:
        output = h.formatPlain(entry)
    }
    
    _, err := fmt.Fprintln(h.writer, output)
    return err
}
```

#### File Output Handler

```go
type FileOutputHandler struct {
    config FileOutput
    file   *os.File
    writer *bufio.Writer
    mutex  sync.Mutex
}

func NewFileOutputHandler(config FileOutput) (*FileOutputHandler, error) {
    // Create directory if it doesn't exist
    dir := filepath.Dir(config.Path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create log directory: %w", err)
    }
    
    // Open or create log file
    file, err := os.OpenFile(config.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return nil, fmt.Errorf("failed to open log file: %w", err)
    }
    
    return &FileOutputHandler{
        config: config,
        file:   file,
        writer: bufio.NewWriter(file),
    }, nil
}

func (h *FileOutputHandler) Write(entry LogEntry) error {
    h.mutex.Lock()
    defer h.mutex.Unlock()
    
    output := h.formatEntry(entry)
    _, err := h.writer.WriteString(output + "\n")
    if err != nil {
        return fmt.Errorf("failed to write to log file: %w", err)
    }
    
    return h.writer.Flush()
}
```

## Error Handling

### LoggingError

Custom error type for logging-specific errors:

```go
type LoggingError struct {
    Component string
    Operation string
    Message   string
    Cause     error
    Fields    map[string]interface{}
}

func (e *LoggingError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("logging error in %s during %s: %s: %v", 
            e.Component, e.Operation, e.Message, e.Cause)
    }
    return fmt.Sprintf("logging error in %s during %s: %s", 
        e.Component, e.Operation, e.Message)
}

func (e *LoggingError) Unwrap() error {
    return e.Cause
}
```

### ValidationError

Error type for configuration validation:

```go
type ValidationError struct {
    Field   string
    Value   interface{}
    Message string
    Rule    string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error for field '%s' with value '%v': %s (rule: %s)", 
        e.Field, e.Value, e.Message, e.Rule)
}
```

## Validation Rules

### Configuration Validation

The logging system validates configuration using the following rules:

```go
func ValidateLoggingConfig(config *spookytypes.LoggingConfig) error {
    var errors []error
    
    // Validate log level
    if config.Level != "" {
        if !isValidLogLevel(config.Level) {
            errors = append(errors, &ValidationError{
                Field:   "level",
                Value:   config.Level,
                Message: "invalid log level",
                Rule:    "must be one of: trace, debug, info, warn, error, fatal",
            })
        }
    }
    
    // Validate format
    if config.Format != "" {
        if !isValidFormat(config.Format) {
            errors = append(errors, &ValidationError{
                Field:   "format",
                Value:   config.Format,
                Message: "invalid output format",
                Rule:    "must be one of: json, structured, plain",
            })
        }
    }
    
    // Validate outputs
    for i, output := range config.Outputs {
        if err := validateOutput(output, i); err != nil {
            errors = append(errors, err)
        }
    }
    
    // Validate components
    for component, compConfig := range config.Components {
        if err := validateComponentConfig(compConfig, component); err != nil {
            errors = append(errors, err)
        }
    }
    
    if len(errors) > 0 {
        return &AggregateError{
            Message: "logging configuration validation failed",
            Errors:  errors,
        }
    }
    
    return nil
}
```

### Output Validation

```go
func validateOutput(output LoggingOutput, index int) error {
    var errors []error
    
    // Validate output type
    if !isValidOutputType(output.Type) {
        errors = append(errors, &ValidationError{
            Field:   fmt.Sprintf("outputs[%d].type", index),
            Value:   output.Type,
            Message: "invalid output type",
            Rule:    "must be one of: console, file",
        })
    }
    
    // Validate file output
    if output.Type == "file" && output.File != nil {
        if output.File.Path == "" {
            errors = append(errors, &ValidationError{
                Field:   fmt.Sprintf("outputs[%d].file.path", index),
                Value:   output.File.Path,
                Message: "file path is required for file output",
                Rule:    "must be a valid file path",
            })
        }
    }
    
    if len(errors) > 0 {
        return &AggregateError{
            Message: fmt.Sprintf("output[%d] validation failed", index),
            Errors:  errors,
        }
    }
    
    return nil
}
```

## CLI Integration

### Logging Commands

The logging system integrates with the CLI through the following commands:

```go
var loggingCmd = &cobra.Command{
    Use:   "logging",
    Short: "Manage logging configuration",
    Long:  "Configure and manage logging for spooky operations",
}

var loggingShowCmd = &cobra.Command{
    Use:   "show",
    Short: "Show current logging configuration",
    RunE: func(cmd *cobra.Command, args []string) error {
        config, err := getLoggingConfig()
        if err != nil {
            return err
        }
        
        return displayLoggingConfig(config)
    },
}

var loggingConfigureCmd = &cobra.Command{
    Use:   "configure",
    Short: "Configure logging settings",
    RunE: func(cmd *cobra.Command, args []string) error {
        level, _ := cmd.Flags().GetString("level")
        format, _ := cmd.Flags().GetString("format")
        file, _ := cmd.Flags().GetString("file")
        
        config := &spookytypes.LoggingConfig{
            Level:  level,
            Format: format,
        }
        
        if file != "" {
            config.Outputs = append(config.Outputs, spookytypes.LoggingOutput{
                Type: "file",
                File: &spookytypes.FileOutput{
                    Path: file,
                },
            })
        }
        
        return configureLogging(config)
    },
}

var loggingValidateCmd = &cobra.Command{
    Use:   "validate",
    Short: "Validate logging configuration",
    RunE: func(cmd *cobra.Command, args []string) error {
        config, err := getLoggingConfig()
        if err != nil {
            return err
        }
        
        result, err := validateLoggingConfig(config)
        if err != nil {
            return err
        }
        
        return displayValidationResult(result)
    },
}

var loggingTestCmd = &cobra.Command{
    Use:   "test",
    Short: "Test logging configuration",
    RunE: func(cmd *cobra.Command, args []string) error {
        level, _ := cmd.Flags().GetString("level")
        component, _ := cmd.Flags().GetString("component")
        
        return testLogging(level, component)
    },
}
```

### Command Flags

```go
func init() {
    loggingConfigureCmd.Flags().String("level", "", "Set log level (trace, debug, info, warn, error, fatal)")
    loggingConfigureCmd.Flags().String("format", "", "Set output format (json, structured, plain)")
    loggingConfigureCmd.Flags().String("file", "", "Set log file path")
    loggingConfigureCmd.Flags().String("max-size", "", "Set maximum file size")
    loggingConfigureCmd.Flags().String("max-age", "", "Set maximum file age")
    loggingConfigureCmd.Flags().Int("max-backups", 0, "Set maximum number of backups")
    
    loggingTestCmd.Flags().String("level", "info", "Test log level")
    loggingTestCmd.Flags().String("component", "", "Test specific component")
}
```

## Performance Considerations

### Async Logging

For high-performance scenarios, the logging system supports async logging:

```go
type AsyncLogger struct {
    logger    Logger
    queue     chan LogEntry
    workers   int
    stopChan  chan struct{}
    wg        sync.WaitGroup
}

func NewAsyncLogger(logger Logger, workers int, queueSize int) *AsyncLogger {
    al := &AsyncLogger{
        logger:   logger,
        queue:    make(chan LogEntry, queueSize),
        workers:  workers,
        stopChan: make(chan struct{}),
    }
    
    // Start worker goroutines
    for i := 0; i < workers; i++ {
        al.wg.Add(1)
        go al.worker()
    }
    
    return al
}

func (al *AsyncLogger) worker() {
    defer al.wg.Done()
    
    for {
        select {
        case entry := <-al.queue:
            al.logger.Log(entry.Level, entry.Message, entry.Fields)
        case <-al.stopChan:
            return
        }
    }
}

func (al *AsyncLogger) Info(msg string, fields ...interface{}) {
    entry := LogEntry{
        Level:   "info",
        Message: msg,
        Fields:  al.parseFields(fields...),
    }
    
    select {
    case al.queue <- entry:
    default:
        // Queue full, log synchronously
        al.logger.Info(msg, fields...)
    }
}
```

### Log Rotation

The file output handler supports automatic log rotation:

```go
func (h *FileOutputHandler) checkRotation() error {
    info, err := h.file.Stat()
    if err != nil {
        return err
    }
    
    maxSize := h.parseSize(h.config.MaxSize)
    if maxSize > 0 && info.Size() >= maxSize {
        return h.rotate()
    }
    
    maxAge := h.parseDuration(h.config.MaxAge)
    if maxAge > 0 && time.Since(info.ModTime()) >= maxAge {
        return h.rotate()
    }
    
    return nil
}

func (h *FileOutputHandler) rotate() error {
    h.mutex.Lock()
    defer h.mutex.Unlock()
    
    // Close current file
    h.writer.Flush()
    h.file.Close()
    
    // Rotate existing backups
    h.rotateBackups()
    
    // Create new file
    file, err := os.OpenFile(h.config.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    
    h.file = file
    h.writer = bufio.NewWriter(file)
    
    return nil
}
```

## Integration Patterns

### Component Integration

Components integrate with the logging system through the `LogManager`:

```go
type FactManager struct {
    logger spookylogging.Logger
    // ... other fields
}

func NewFactManager(logManager spookylogging.LogManager) *FactManager {
    return &FactManager{
        logger: logManager.GetLogger("facts"),
        // ... initialize other fields
    }
}

func (m *FactManager) CollectFacts(server string) (*spookytypes.FactCollection, error) {
    m.logger.Info("Starting fact collection", "server", server)
    
    start := time.Now()
    defer func() {
        m.logger.Info("Fact collection completed", 
            "server", server,
            "duration_ms", time.Since(start).Milliseconds())
    }()
    
    // ... implementation
}
```

### Context-Aware Logging

```go
func (m *FactManager) CollectFactsWithContext(ctx context.Context, server string) (*spookytypes.FactCollection, error) {
    logger := m.logger.With(
        "request_id", ctx.Value("request_id"),
        "user", ctx.Value("user"),
        "session_id", ctx.Value("session_id"))
    
    logger.Info("Starting fact collection with context", "server", server)
    
    // ... implementation
}
```

## Testing

### Mock Logger

For testing, the logging system provides a mock logger:

```go
type MockLogger struct {
    entries []LogEntry
    mutex   sync.RWMutex
}

func NewMockLogger() *MockLogger {
    return &MockLogger{
        entries: make([]LogEntry, 0),
    }
}

func (m *MockLogger) Info(msg string, fields ...interface{}) {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    m.entries = append(m.entries, LogEntry{
        Level:   "info",
        Message: msg,
        Fields:  m.parseFields(fields...),
    })
}

func (m *MockLogger) GetEntries() []LogEntry {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    
    entries := make([]LogEntry, len(m.entries))
    copy(entries, m.entries)
    return entries
}

func (m *MockLogger) Clear() {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    m.entries = m.entries[:0]
}
```

### Test Utilities

```go
func TestLoggingConfiguration(t *testing.T) {
    config := &spookytypes.LoggingConfig{
        Level:  "debug",
        Format: "json",
        Outputs: []spookytypes.LoggingOutput{
            {
                Type: "console",
                Console: &spookytypes.ConsoleOutput{
                    Colorize: true,
                },
            },
        },
    }
    
    manager := NewManager(config)
    logger := manager.GetLogger("test")
    
    logger.Info("Test message", "key", "value")
    
    // Verify logging behavior
    // ... assertions
}
```

## Conclusion

The spooky logging system provides a comprehensive, interface-based solution for logging across all spooky components. The system is designed to be extensible, performant, and easy to integrate with existing code.

Key features include:
- Interface-based design for loose coupling
- Multiple output formats and destinations
- Component-specific configuration
- Performance optimization with async logging
- Comprehensive validation and error handling
- Easy integration with CLI and other components

For usage examples and best practices, refer to the [User Guide](LOGGING_USER_GUIDE.md) and [Troubleshooting Guide](LOGGING_TROUBLESHOOTING.md).
