# Logging System User Guide

## Overview

The spooky logging system provides comprehensive logging capabilities for all spooky components. This guide covers everything from basic logging configuration to advanced features like structured logging, multiple output formats, and performance optimization.

## Getting Started

### Prerequisites

- spooky CLI installed and configured
- Basic understanding of HCL configuration syntax
- Access to create and modify configuration files

### Quick Start

1. **Check Current Logging Configuration**
   ```bash
   spooky logging show
   ```

2. **Configure Basic Logging**
   ```bash
   spooky logging configure --level info --format json
   ```

3. **Test Logging Output**
   ```bash
   spooky logging test
   ```

## Logging Configuration

### Global Logging Configuration

Global logging configuration is stored in `$XDG_CONFIG_HOME/spooky/logging.hcl` and applies to all spooky operations.

#### Basic Configuration

```hcl
logging {
  level = "info"
  format = "json"
  
  output {
    type = "console"
    enabled = true
  }
  
  output {
    type = "file"
    enabled = true
    path = "~/.local/state/spooky/logs/spooky.log"
  }
}
```

#### Advanced Configuration

```hcl
logging {
  level = "debug"
  format = "structured"
  
  # Console output
  output {
    type = "console"
    enabled = true
    colorize = true
  }
  
  # File output with rotation
  output {
    type = "file"
    enabled = true
    path = "~/.local/state/spooky/logs/spooky.log"
    max_size = "100MB"
    max_age = "30d"
    max_backups = 10
  }
  
  # Component-specific logging
  components {
    "facts" {
      level = "debug"
      enabled = true
    }
    
    "machines" {
      level = "info"
      enabled = true
    }
    
    "variables" {
      level = "warn"
      enabled = true
    }
  }
}
```

### Project-Specific Logging

Project-specific logging configuration overrides global settings for specific projects.

#### Project Configuration

```hcl
# In your project directory
logging {
  level = "debug"
  format = "json"
  
  output {
    type = "file"
    enabled = true
    path = "./logs/project.log"
  }
  
  # Override component logging for this project
  components {
    "actions" {
      level = "debug"
      enabled = true
    }
  }
}
```

### Environment Variables

You can override logging configuration using environment variables:

```bash
# Set log level
export SPOOKY_LOG_LEVEL=debug

# Set log format
export SPOOKY_LOG_FORMAT=json

# Set log file path
export SPOOKY_LOG_FILE=./debug.log

# Enable/disable console output
export SPOOKY_LOG_CONSOLE=true

# Enable/disable file output
export SPOOKY_LOG_FILE_OUTPUT=true
```

## Log Levels

### Available Levels

1. **trace** - Most detailed logging, includes function entry/exit
2. **debug** - Detailed debugging information
3. **info** - General information about operations
4. **warn** - Warning messages for potential issues
5. **error** - Error messages for failed operations
6. **fatal** - Critical errors that cause program termination

### Level Hierarchy

```
trace < debug < info < warn < error < fatal
```

Only messages at or above the configured level are logged.

### Component-Specific Levels

You can set different log levels for different components:

```hcl
logging {
  level = "info"  # Default level
  
  components {
    "facts" {
      level = "debug"  # More detailed for facts
    }
    
    "machines" {
      level = "warn"   # Less verbose for machines
    }
    
    "variables" {
      level = "error"  # Only errors for variables
    }
  }
}
```

## Output Formats

### JSON Format

Structured JSON output for machine processing:

```json
{
  "timestamp": "2024-01-15T10:30:45.123Z",
  "level": "info",
  "component": "facts",
  "message": "Collecting facts from server",
  "fields": {
    "server": "web.example.com",
    "method": "ssh",
    "duration_ms": 1250
  }
}
```

### Structured Format

Human-readable structured output:

```
2024-01-15T10:30:45.123Z [INFO] facts: Collecting facts from server server=web.example.com method=ssh duration_ms=1250
```

### Plain Text Format

Simple text output:

```
2024-01-15T10:30:45.123Z INFO facts: Collecting facts from server
```

### Configuration

```hcl
logging {
  format = "json"  # json, structured, or plain
  
  # Format-specific options
  json {
    pretty_print = false
    include_timestamp = true
    include_level = true
    include_component = true
  }
  
  structured {
    colorize = true
    include_timestamp = true
    include_level = true
    include_component = true
  }
  
  plain {
    include_timestamp = true
    include_level = true
    include_component = true
  }
}
```

## Output Destinations

### Console Output

Output to standard output/error:

```hcl
logging {
  output {
    type = "console"
    enabled = true
    colorize = true
    use_stderr = false  # Use stderr for errors
  }
}
```

### File Output

Output to log files with rotation:

```hcl
logging {
  output {
    type = "file"
    enabled = true
    path = "~/.local/state/spooky/logs/spooky.log"
    max_size = "100MB"
    max_age = "30d"
    max_backups = 10
    compress = true
  }
}
```

### Multiple Outputs

You can configure multiple outputs simultaneously:

```hcl
logging {
  # Console for immediate feedback
  output {
    type = "console"
    enabled = true
    level = "info"
  }
  
  # File for persistent storage
  output {
    type = "file"
    enabled = true
    path = "./logs/debug.log"
    level = "debug"
  }
  
  # Error file for errors only
  output {
    type = "file"
    enabled = true
    path = "./logs/errors.log"
    level = "error"
  }
}
```

## Component-Based Logging

### Available Components

- **facts** - Fact collection and processing
- **machines** - Machine inventory and connectivity
- **variables** - Variable management and resolution
- **actions** - Action orchestration and running
- **templates** - Template rendering and processing
- **ssh** - SSH connections and operations
- **config** - Configuration loading and validation
- **cli** - Command-line interface operations

### Component Configuration

```hcl
logging {
  components {
    "facts" {
      level = "debug"
      enabled = true
      output {
        type = "file"
        path = "./logs/facts.log"
      }
    }
    
    "machines" {
      level = "info"
      enabled = true
    }
    
    "variables" {
      level = "warn"
      enabled = false  # Disable logging for variables
    }
  }
}
```

## Advanced Features

### Structured Logging

Add structured fields to your log messages:

```go
// In your code
logger.Info("Processing machine", 
    "machine", "web.example.com",
    "method", "ssh",
    "duration_ms", 1250,
    "success", true)
```

### Log Context

Add context to log messages:

```go
// Create a logger with context
ctxLogger := logger.With(
    "project", "my-project",
    "user", "admin",
    "session_id", "abc123")

// Use context logger
ctxLogger.Info("Starting operation")
ctxLogger.Error("Operation failed", "error", err)
```

### Performance Logging

Log performance metrics:

```hcl
logging {
  performance {
    enabled = true
    threshold_ms = 1000  # Log operations taking > 1 second
    include_memory = true
    include_cpu = true
  }
}
```

### Audit Logging

Log security and audit events:

```hcl
logging {
  audit {
    enabled = true
    level = "info"
    include_user = true
    include_ip = true
    include_session = true
  }
}
```

## Logging Management

### CLI Commands

#### Show Current Configuration

```bash
# Show global configuration
spooky logging show

# Show project-specific configuration
spooky logging show --project ./my-project
```

#### Configure Logging

```bash
# Configure basic logging
spooky logging configure --level info --format json

# Configure with file output
spooky logging configure \
  --level debug \
  --format structured \
  --file ./logs/spooky.log \
  --max-size 100MB \
  --max-age 30d

# Configure component-specific logging
spooky logging configure \
  --component facts --level debug \
  --component machines --level warn
```

#### Validate Configuration

```bash
# Validate global configuration
spooky logging validate

# Validate project configuration
spooky logging validate --project ./my-project
```

#### Test Logging

```bash
# Test current configuration
spooky logging test

# Test with specific level
spooky logging test --level debug

# Test specific component
spooky logging test --component facts
```

### Configuration Files

#### Global Configuration Location

- **Linux/macOS**: `~/.config/spooky/logging.hcl`
- **Windows**: `%APPDATA%\spooky\logging.hcl`

#### Project Configuration

- **Location**: `./logging.hcl` in project directory
- **Overrides**: Global configuration for project-specific settings

## Best Practices

### Configuration Best Practices

1. **Start with Info Level** - Use `info` level for production
2. **Use Debug for Development** - Enable `debug` level during development
3. **Configure File Rotation** - Prevent log files from growing too large
4. **Separate Error Logs** - Keep error logs separate from general logs
5. **Use Structured Format** - Use structured or JSON format for better parsing

### Performance Best Practices

1. **Limit Debug Logging** - Debug logging can impact performance
2. **Use Appropriate Levels** - Don't log everything at trace level
3. **Configure Log Rotation** - Prevent disk space issues
4. **Monitor Log Performance** - Watch for logging overhead
5. **Use Async Logging** - Enable async logging for high-volume operations

### Security Best Practices

1. **Don't Log Sensitive Data** - Avoid logging passwords, keys, or tokens
2. **Use Audit Logging** - Enable audit logging for security events
3. **Secure Log Files** - Set appropriate permissions on log files
4. **Rotate Logs Regularly** - Prevent log file accumulation
5. **Monitor Log Access** - Track who accesses log files

### Development Best Practices

1. **Use Context Loggers** - Add context to log messages
2. **Include Relevant Fields** - Add structured fields for debugging
3. **Test Logging Configuration** - Validate logging setup
4. **Use Component Logging** - Configure per-component logging levels
5. **Document Logging Strategy** - Document your logging approach

## Examples

### Basic Project Setup

```hcl
# logging.hcl in your project directory
logging {
  level = "info"
  format = "structured"
  
  output {
    type = "console"
    enabled = true
    colorize = true
  }
  
  output {
    type = "file"
    enabled = true
    path = "./logs/project.log"
    max_size = "50MB"
    max_age = "7d"
    max_backups = 5
  }
}
```

### Development Configuration

```hcl
logging {
  level = "debug"
  format = "json"
  
  output {
    type = "console"
    enabled = true
  }
  
  output {
    type = "file"
    enabled = true
    path = "./logs/debug.log"
    level = "debug"
  }
  
  components {
    "facts" {
      level = "trace"
      enabled = true
    }
    
    "machines" {
      level = "debug"
      enabled = true
    }
  }
}
```

### Production Configuration

```hcl
logging {
  level = "warn"
  format = "json"
  
  output {
    type = "file"
    enabled = true
    path = "/var/log/spooky/app.log"
    max_size = "100MB"
    max_age = "30d"
    max_backups = 10
    compress = true
  }
  
  output {
    type = "file"
    enabled = true
    path = "/var/log/spooky/errors.log"
    level = "error"
    max_size = "50MB"
    max_age = "90d"
    max_backups = 20
  }
  
  audit {
    enabled = true
    level = "info"
  }
}
```

## Validation and Troubleshooting

### Configuration Validation

```bash
# Validate global configuration
spooky logging validate

# Validate project configuration
spooky logging validate --project ./my-project

# Check for common issues
spooky logging validate --strict
```

### Common Issues

#### Log Files Not Created

1. **Check Permissions** - Ensure write permissions to log directory
2. **Check Path** - Verify log file path is correct
3. **Check Configuration** - Validate logging configuration

#### Performance Issues

1. **Reduce Log Level** - Use higher log levels in production
2. **Enable Log Rotation** - Prevent large log files
3. **Use Async Logging** - Enable async logging for better performance

#### Missing Log Messages

1. **Check Log Level** - Ensure log level is appropriate
2. **Check Component Settings** - Verify component logging is enabled
3. **Check Output Configuration** - Ensure outputs are properly configured

### Debugging Logging Issues

```bash
# Test logging with verbose output
spooky logging test --verbose

# Check logging configuration
spooky logging show --verbose

# Validate configuration with details
spooky logging validate --verbose
```

## Integration with Other Systems

### Integration with Facts System

```hcl
logging {
  components {
    "facts" {
      level = "debug"
      enabled = true
      output {
        type = "file"
        path = "./logs/facts.log"
      }
    }
  }
}
```

### Integration with Machines System

```hcl
logging {
  components {
    "machines" {
      level = "info"
      enabled = true
      output {
        type = "file"
        path = "./logs/machines.log"
      }
    }
  }
}
```

### Integration with Variables System

```hcl
logging {
  components {
    "variables" {
      level = "warn"
      enabled = true
    }
  }
}
```

## Conclusion

The spooky logging system provides comprehensive logging capabilities that can be configured for various use cases, from simple development logging to complex production environments with multiple outputs and component-specific configurations.

Start with basic configuration and gradually add more advanced features as needed. Always validate your configuration and test logging output to ensure it meets your requirements.

For more advanced usage and troubleshooting, refer to the [API Reference](LOGGING_API_REFERENCE.md) and [Troubleshooting Guide](LOGGING_TROUBLESHOOTING.md).
