# Logging System

## Overview

The Logging System provides comprehensive logging capabilities for the spooky codebase. It enables structured logging, configurable output formats, secure logging practices, and integration with all system components for observability and debugging.

**Status**: **Implemented** - Complete logging system with structured logging, configuration management, and CLI integration.

## Related Systems

This system integrates with and supports all other spooky systems:

- **[Actions System](ACTIONS_SYSTEM.md)** - Actions generate comprehensive logs for monitoring and debugging
- **[Facts System](FACTS_SYSTEM.md)** - Fact collection generates logs for tracking and troubleshooting
- **[Machines System](MACHINES_SYSTEM.md)** - Machine operations generate logs for connectivity and status monitoring
- **[SSH System](SSH_SYSTEM.md)** - SSH operations generate logs for connection and security monitoring
- **[Templates System](TEMPLATES_SYSTEM.md)** - Template rendering generates logs for debugging and validation
- **[Variables System](VARIABLES_SYSTEM.md)** - Variable resolution generates logs for troubleshooting
- **[Integrations System](INTEGRATIONS_SYSTEM.md)** - System integrations generate logs for coordination monitoring
- **[Projects System](PROJECTS_SYSTEM.md)** - Project operations generate logs for lifecycle management

## Architecture

### Core Components

#### Log Manager
- **File**: `internal/logging/logging.go`
- **Purpose**: Central logging management with configuration and output handling
- **Features**:
  - Logger instantiation and management
  - Log level configuration
  - Output format configuration
  - Log rotation and retention
  - Performance monitoring
  - Error handling and recovery

#### Secure Logger
- **File**: `internal/logging/secure_logger.go`
- **Purpose**: Secure logging with sensitive data protection
- **Features**:
  - Sensitive data redaction
  - Secure log output
  - Audit logging
  - Compliance logging
  - Data protection
  - Security validation

#### Logging Integration
- **File**: `internal/logging/integration.go`
- **Purpose**: Interface implementation for system integration
- **Features**:
  - GetLogger - Get configured logger instance
  - ConfigureLogging - Configure logging settings
  - SetLogLevel - Set log level dynamically
  - GetLogLevel - Get current log level

### Integration Points

#### All System Components
- Provides logging capabilities to all system components
- Supports consistent logging patterns
- Enables centralized log management

#### Configuration Integration
- Provides logging configuration management
- Supports dynamic configuration updates
- Enables environment-based configuration

#### CLI Integration
- Provides logging for CLI operations
- Supports verbose and debug modes
- Enables log output formatting

## Logging Types

### Log Structure
```go
type LogEntry struct {
    Timestamp   time.Time              // Log timestamp
    Level       string                 // Log level
    Message     string                 // Log message
    Fields      map[string]interface{} // Structured fields
    Component   string                 // Component name
    Operation   string                 // Operation name
    Error       error                  // Associated error
    StackTrace  string                 // Stack trace
    Metadata    map[string]interface{} // Additional metadata
}
```

### Log Levels
- **DEBUG**: Detailed debugging information
- **INFO**: General information messages
- **WARN**: Warning messages
- **ERROR**: Error messages
- **FATAL**: Fatal error messages

### Log Categories
- **System**: System-level operations
- **Security**: Security-related events
- **Performance**: Performance metrics
- **User**: User-initiated operations
- **Audit**: Audit trail events

## Logging Features

### Structured Logging
- **JSON Output**: Structured JSON log output
- **Field Support**: Custom fields and metadata
- **Context Support**: Request context and tracing
- **Correlation**: Request correlation IDs

### Configuration Management
- **Log Level Control**: Dynamic log level configuration
- **Output Format**: Configurable output formats
- **Output Destinations**: Multiple output destinations
- **Filtering**: Log filtering and routing

### Performance Features
- **Async Logging**: Asynchronous log processing
- **Buffering**: Log buffering for performance
- **Batching**: Log batching for efficiency
- **Compression**: Log compression for storage

### Security Features
- **Data Redaction**: Sensitive data redaction
- **Access Control**: Log access control
- **Audit Logging**: Comprehensive audit logging
- **Compliance**: Compliance logging support

## Logging Configuration

### Global Configuration
```hcl
# logging/global-config.hcl
logging_config {
  # General settings
  general {
    level = "info"
    format = "json"
    timestamp_format = "RFC3339"
    include_caller = true
    include_stack_trace = false
  }
  
  # Output settings
  output {
    destinations = ["stdout", "file"]
    file_path = "~/.local/state/spooky/logs/spooky.log"
    max_file_size = "100MB"
    max_files = 10
    compress_old_files = true
  }
  
  # Security settings
  security {
    redact_sensitive_fields = true
    sensitive_fields = [
      "password",
      "secret",
      "key",
      "token",
      "credential"
    ]
    audit_logging = true
    compliance_logging = true
  }
  
  # Performance settings
  performance {
    async_logging = true
    buffer_size = 1000
    flush_interval = "5s"
    batch_size = 100
  }
}
```

### Component Configuration
```hcl
# logging/component-config.hcl
component_logging {
  # Facts component
  facts {
    level = "info"
    include_debug = false
    sensitive_operations = ["collect_facts"]
  }
  
  # Actions component
  actions {
    level = "info"
    include_debug = true
    sensitive_operations = ["run_action"]
  }
  
  # SSH component
  ssh {
    level = "warn"
    include_debug = false
    sensitive_operations = ["connect", "authenticate"]
  }
  
  # Templates component
  templates {
    level = "info"
    include_debug = true
    sensitive_operations = ["render_template"]
  }
}
```

## CLI Commands

### Logging Configuration
```bash
# Show logging configuration
spooky logging config

# Set log level
spooky logging set-level debug

# Get current log level
spooky logging get-level

# Configure logging
spooky logging configure --level info --format json
```

### Log Management
```bash
# View logs
spooky logging view

# View logs with filtering
spooky logging view --level error --component facts

# View logs with time range
spooky logging view --since "1h ago" --until "now"

# View logs with search
spooky logging view --search "authentication failed"
```

### Log Analysis
```bash
# Analyze log patterns
spooky logging analyze

# Generate log report
spooky logging report --output report.json

# Check log health
spooky logging health

# Validate log configuration
spooky logging validate
```

## Examples

### Basic Logging
```go
// Get logger instance
logger := spookylogging.GetLogger()

// Basic logging
logger.Info("Application started")
logger.Debug("Debug information", "component", "facts")
logger.Warn("Warning message", "operation", "collect_facts")
logger.Error("Error occurred", "error", err, "component", "ssh")
```

### Structured Logging
```go
// Structured logging with fields
logger.Info("Fact collection completed",
    "machine", "web-server",
    "facts_collected", 25,
    "duration", "2.5s",
    "component", "facts",
    "operation", "collect_facts",
)

// Error logging with context
logger.Error("SSH connection failed",
    "machine", "web-server",
    "error", err,
    "component", "ssh",
    "operation", "connect",
    "retry_attempt", 3,
)
```

### Secure Logging
```go
// Secure logger for sensitive operations
secureLogger := spookylogging.NewSecureLogger(logger)

// Log with sensitive data redaction
secureLogger.Info("User authentication",
    "username", "admin",
    "password", "secret123", // Will be redacted
    "token", "abc123",       // Will be redacted
    "ip_address", "192.168.1.100",
)
```

### Component Logging
```go
// Component-specific logging
factsLogger := spookylogging.GetLogger("facts")
factsLogger.Info("Starting fact collection",
    "targets", []string{"web-server", "db-server"},
    "parallel", true,
)

actionsLogger := spookylogging.GetLogger("actions")
actionsLogger.Info("Action running started",
    "action", "deploy-web",
    "machines", []string{"web-server"},
    "parallel", false,
)
```

## Integration Examples

### Facts Integration
```go
// Logging in facts collection
func (m *FactManager) CollectFacts(machine string) (*spookytypes.FactCollection, error) {
    logger := spookylogging.GetLogger("facts")
    
    logger.Info("Starting fact collection",
        "machine", machine,
        "component", "facts",
        "operation", "collect_facts",
    )
    
    start := time.Now()
    facts, err := m.collector.Collect(machine)
    duration := time.Since(start)
    
    if err != nil {
        logger.Error("Fact collection failed",
            "machine", machine,
            "error", err,
            "duration", duration,
            "component", "facts",
            "operation", "collect_facts",
        )
        return nil, err
    }
    
    logger.Info("Fact collection completed",
        "machine", machine,
        "facts_collected", len(facts.Facts),
        "duration", duration,
        "component", "facts",
        "operation", "collect_facts",
    )
    
    return facts, nil
}
```

### Actions Integration
```go
// Logging in action running
func (m *ActionManager) RunAction(action *spookytypes.Action) error {
    logger := spookylogging.GetLogger("actions")
    
    logger.Info("Starting action running",
        "action", action.Name,
        "machine", action.Machine,
        "component", "actions",
        "operation", "run_action",
    )
    
    start := time.Now()
    err := m.orchestrator.Run(action)
    duration := time.Since(start)
    
    if err != nil {
        logger.Error("Action running failed",
            "action", action.Name,
            "machine", action.Machine,
            "error", err,
            "duration", duration,
            "component", "actions",
            "operation", "run_action",
        )
        return err
    }
    
    logger.Info("Action running completed",
        "action", action.Name,
        "machine", action.Machine,
        "duration", duration,
        "component", "actions",
        "operation", "run_action",
    )
    
    return nil
}
```

### SSH Integration
```go
// Logging in SSH operations
func (c *SSHClient) Connect(host string, config *spookytypes.SSHConfig) error {
    logger := spookylogging.GetLogger("ssh")
    
    logger.Info("Attempting SSH connection",
        "host", host,
        "user", config.User,
        "port", config.Port,
        "component", "ssh",
        "operation", "connect",
    )
    
    start := time.Now()
    err := c.connect(host, config)
    duration := time.Since(start)
    
    if err != nil {
        logger.Error("SSH connection failed",
            "host", host,
            "user", config.User,
            "error", err,
            "duration", duration,
            "component", "ssh",
            "operation", "connect",
        )
        return err
    }
    
    logger.Info("SSH connection established",
        "host", host,
        "user", config.User,
        "duration", duration,
        "component", "ssh",
        "operation", "connect",
    )
    
    return nil
}
```

## Best Practices

### Logging Design
- Use structured logging with consistent fields
- Include relevant context in log messages
- Use appropriate log levels
- Implement secure logging for sensitive data

### Performance
- Use async logging for high-performance operations
- Implement log buffering and batching
- Monitor logging performance impact
- Use log compression for storage efficiency

### Security
- Redact sensitive data in logs
- Implement audit logging for security events
- Control log access and permissions
- Use secure log storage

### Maintainability
- Use consistent logging patterns
- Document logging configuration
- Implement log rotation and retention
- Monitor log quality and completeness

## Troubleshooting

### Common Issues

#### Log Level Issues
```bash
# Check current log level
spooky logging get-level

# Set appropriate log level
spooky logging set-level debug

# Check component-specific levels
spooky logging config --component facts
```

#### Log Output Issues
```bash
# Check log file permissions
ls -la ~/.local/state/spooky/logs/

# Check log file size
du -h ~/.local/state/spooky/logs/spooky.log

# Check log rotation
spooky logging rotate
```

#### Performance Issues
```bash
# Check logging performance
spooky logging performance

# Check buffer usage
spooky logging buffer-status

# Flush log buffer
spooky logging flush
```

#### Security Issues
```bash
# Check sensitive data in logs
spooky logging audit --sensitive

# Validate log security
spooky logging validate --security

# Check audit logs
spooky logging audit --type security
```

## API Reference

### LoggingIntegration Interface
```go
type LoggingIntegration interface {
    GetLogger(component string) spookylogging.Logger
    ConfigureLogging(config *spookytypes.LoggingConfig) error
    SetLogLevel(level string) error
    GetLogLevel() string
}
```

### Log Manager Methods
```go
// Logger management
GetLogger(component string) spookylogging.Logger
GetSecureLogger(component string) spookylogging.SecureLogger

// Configuration management
ConfigureLogging(config *spookytypes.LoggingConfig) error
SetLogLevel(level string) error
GetLogLevel() string

// Log management
RotateLogs() error
FlushLogs() error
GetLogStats() (*spookytypes.LogStats, error)
```

### Logger Interface
```go
type Logger interface {
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    Fatal(msg string, fields ...Field)
    
    WithField(key string, value interface{}) Logger
    WithFields(fields map[string]interface{}) Logger
    WithError(err error) Logger
}
```

## Related Documentation

- [Logging API Reference](LOGGING_API_REFERENCE.md) - Complete API documentation
- [Logging User Guide](LOGGING_USER_GUIDE.md) - User guide and examples
- [Logging Troubleshooting](LOGGING_TROUBLESHOOTING.md) - Troubleshooting guide
- [Logging Framework](LOGGING_FRAMEWORK.md) - Framework documentation
- [Logging Configuration Strategy](LOGGING_CONFIGURATION_STRATEGY.md) - Configuration strategy
- [Actions System](ACTIONS_SYSTEM.md) - Actions logging integration
- [Facts System](FACTS_SYSTEM.md) - Facts logging integration
- [SSH System](SSH_SYSTEM.md) - SSH logging integration
