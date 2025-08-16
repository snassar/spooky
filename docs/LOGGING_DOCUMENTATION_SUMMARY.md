# Logging System Documentation Summary

## Overview

This document provides a comprehensive overview of the spooky logging system documentation. It serves as a guide to help you find the right documentation for your needs and understand how all the pieces fit together.

**Status: Implemented** - The logging system is fully implemented with comprehensive functionality for structured logging, configuration management, and integration with other systems.

## Documentation Structure

### 📚 Core Documentation

#### 1. [User Guide](LOGGING_USER_GUIDE.md)
**Audience:** End users, system administrators, DevOps engineers
**Purpose:** Complete guide to using the logging system

**What it covers:**
- Getting started with logging configuration
- Logging levels and formats
- Configuration management and validation
- Integration with other systems
- Real-world examples and use cases

**When to use:** Start here if you're new to spooky logging or need to understand how to use the system effectively.

#### 2. [API Reference](LOGGING_API_REFERENCE.md)
**Audience:** Developers, system integrators, contributors
**Purpose:** Technical reference for the logging system APIs and implementation

**What it covers:**
- Core interfaces and type definitions
- Implementation details and algorithms
- Error handling patterns
- Configuration rules and schemas
- CLI integration details
- Code examples and patterns

**When to use:** Use this when developing with the logging system, extending functionality, or debugging implementation issues.

#### 3. [Troubleshooting Guide](LOGGING_TROUBLESHOOTING.md)
**Audience:** System administrators, support engineers, users experiencing issues
**Purpose:** Solutions for common problems and debugging techniques

**What it covers:**
- Common error messages and solutions
- Configuration issues and debugging
- Performance problems and optimization
- Integration issues with other systems
- Best practices for troubleshooting

**When to use:** Use this when encountering problems or need to debug issues with the logging system.

### 📁 Examples Directory

#### [Examples Overview](examples/README.md)
**Audience:** All users
**Purpose:** Quick reference for available examples and use cases

**What it covers:**
- Available logging configuration examples
- Example configurations and scripts
- Common use case patterns
- Integration examples with other systems

**When to use:** Use this to quickly find relevant examples for your use case.

## Key Concepts

### Core Features

1. **Structured Logging** - JSON and text-based structured logging
2. **Configurable Levels** - Debug, info, warn, error logging levels
3. **Multiple Formats** - JSON, text, and custom logging formats
4. **Component-Based** - Component-specific logging configuration
5. **Performance Optimized** - Efficient logging with minimal overhead
6. **Integration Support** - Seamless integration with other systems
7. **Security Features** - Secure logging with sensitive data protection

### Architecture Principles

1. **Interface-First Design** - All functionality through well-defined interfaces
2. **Dependency Injection** - Loose coupling through interface-based dependencies
3. **Configurable Design** - Flexible configuration for different use cases
4. **Performance Optimized** - Efficient logging with minimal overhead
5. **Security by Default** - Secure handling of sensitive data

### Best Practices

1. **Use Appropriate Levels** - Choose appropriate logging levels for different scenarios
2. **Structure Log Messages** - Use structured logging for better analysis
3. **Protect Sensitive Data** - Never log sensitive information
4. **Configure Performance** - Optimize logging for production environments
5. **Monitor Log Volume** - Monitor and manage log volume appropriately
6. **Use Component Logging** - Use component-specific logging for better organization

## Logging System Overview

### Core Concepts

The logging system provides a comprehensive solution for structured logging in spooky projects. It supports:

- **Structured Logging** - JSON and text-based structured logging
- **Configurable Levels** - Debug, info, warn, error logging levels
- **Multiple Formats** - JSON, text, and custom logging formats
- **Component-Based** - Component-specific logging configuration
- **Performance Optimization** - Efficient logging with minimal overhead

### Logging Configuration

Logging is configured through HCL configuration files:

```hcl
logging {
  level = "info"
  format = "json"
  output = "stderr"
  
  components {
    ssh {
      level = "debug"
      format = "text"
    }
    
    facts {
      level = "info"
      format = "json"
    }
    
    actions {
      level = "warn"
      format = "json"
    }
  }
  
  security {
    redact_sensitive = true
    redact_patterns = ["password", "secret", "key", "token"]
  }
}
```

### CLI Commands

The logging system provides CLI commands for configuration management:

```bash
# Show logging configuration
spooky logging config show

# Validate logging configuration
spooky logging config validate

# Test logging configuration
spooky logging test

# Show logging status
spooky logging status
```

### Logging Levels

The logging system supports multiple logging levels:

#### Debug Level
```go
logger.Debug("Starting SSH connection", "host", hostname, "user", username)
```

#### Info Level
```go
logger.Info("SSH connection established", "host", hostname, "duration", duration)
```

#### Warn Level
```go
logger.Warn("SSH connection slow", "host", hostname, "duration", duration)
```

#### Error Level
```go
logger.Error("SSH connection failed", "host", hostname, "error", err)
```

### Logging Formats

The logging system supports multiple output formats:

#### JSON Format
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "info",
  "component": "ssh",
  "message": "SSH connection established",
  "fields": {
    "host": "web.example.com",
    "user": "admin",
    "duration": "1.2s"
  }
}
```

#### Text Format
```
2024-01-15T10:30:00Z [INFO] [ssh] SSH connection established host=web.example.com user=admin duration=1.2s
```

### Component-Based Logging

The logging system supports component-specific configuration:

```hcl
logging {
  level = "info"
  format = "json"
  
  components {
    ssh {
      level = "debug"
      format = "text"
    }
    
    facts {
      level = "info"
      format = "json"
    }
    
    actions {
      level = "warn"
      format = "json"
    }
    
    machines {
      level = "info"
      format = "text"
    }
  }
}
```

### Security Features

The logging system includes security features for sensitive data protection:

```hcl
logging {
  level = "info"
  format = "json"
  
  security {
    redact_sensitive = true
    redact_patterns = [
      "password",
      "secret", 
      "key",
      "token",
      "credential",
      "private_key"
    ]
  }
}
```

## Implementation Details

### Core Components

1. **Logger Factory** - Creates and manages logger instances
2. **Log Formatter** - Formats log messages in different formats
3. **Log Handler** - Handles log output and routing
4. **Configuration Manager** - Manages logging configuration
5. **Security Manager** - Handles sensitive data protection

### Integration Points

The logging system integrates with:

- **CLI System** - For logging configuration and management
- **Configuration System** - For logging configuration validation
- **All Other Systems** - For component-specific logging
- **Security System** - For sensitive data protection

### Error Handling

The logging system provides comprehensive error handling:

- **Configuration errors** - Invalid logging configuration
- **Format errors** - Invalid log format specifications
- **Output errors** - Log output and routing failures
- **Security errors** - Sensitive data protection issues
- **Performance errors** - Logging performance issues

## Best Practices

### Logging Configuration

1. **Use appropriate levels** for different environments
2. **Structure log messages** for better analysis
3. **Configure component-specific** logging when needed
4. **Protect sensitive data** with redaction patterns
5. **Optimize performance** for production environments

### Log Message Structure

1. **Use consistent field names** across components
2. **Include relevant context** in log messages
3. **Use appropriate log levels** for different types of messages
4. **Structure complex data** for better analysis
5. **Avoid logging sensitive** information

### Performance Optimization

1. **Use appropriate log levels** to control volume
2. **Structure log messages** for efficient processing
3. **Configure output destinations** appropriately
4. **Monitor log volume** and performance impact
5. **Use component-specific** configuration when needed

## Troubleshooting

### Common Issues

1. **Configuration errors** - Check logging configuration syntax
2. **Format errors** - Verify log format specifications
3. **Output errors** - Check log output destinations and permissions
4. **Performance issues** - Monitor log volume and performance impact
5. **Security issues** - Verify sensitive data protection configuration

### Debug Commands

```bash
# Enable debug logging
export SPOOKY_LOG_LEVEL=debug

# Show logging configuration
spooky logging config show --verbose

# Test logging configuration
spooky logging test --verbose

# Validate configuration
spooky logging config validate --verbose
```

### Common Patterns

1. **Environment-specific configuration** - Use different configurations for different environments
2. **Component-specific logging** - Configure logging per component
3. **Structured logging** - Use structured fields for better analysis
4. **Security configuration** - Always configure sensitive data protection
5. **Performance monitoring** - Monitor log volume and performance impact

## Related Documentation

- [Logging User Guide](LOGGING_USER_GUIDE.md) - Complete user guide
- [Logging API Reference](LOGGING_API_REFERENCE.md) - Technical reference
- [Logging Troubleshooting](LOGGING_TROUBLESHOOTING.md) - Troubleshooting guide
- [System Design](../design/systems/logging-system.md) - System design documentation
- [CLI Reference](CLI_REFERENCE.md) - CLI command reference
