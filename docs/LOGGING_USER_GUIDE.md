# Logging System User Guide

## Overview

The spooky logging system provides comprehensive logging capabilities for monitoring and debugging operations across all system components. This guide covers everything from basic logging configuration to advanced features like structured logging, log aggregation, and performance monitoring.

**Status: Production Ready** - The logging system is fully implemented with comprehensive logging, formatting, and output capabilities.

## Getting Started

### Prerequisites

- spooky CLI installed and configured
- Basic understanding of logging concepts
- Access to create and modify configuration files
- Understanding of log levels and formatting

### Quick Start

1. **Check Available Logging Commands**
   ```bash
   spooky logging --help
   ```

2. **View Current Logging Configuration**
   ```bash
   spooky logging config show
   ```

3. **Set Log Level**
   ```bash
   spooky logging config set --level debug
   ```

## Core Concepts

### Log Levels

spooky supports standard log levels:

- **DEBUG** - Detailed debugging information
- **INFO** - General operational information
- **WARN** - Warning messages for potential issues
- **ERROR** - Error messages for failed operations
- **FATAL** - Critical errors that cause system failure

### Logging Components

The logging system provides:

- **Structured Logging** - JSON and text formatting
- **Multiple Outputs** - Console, file, and network outputs
- **Performance Monitoring** - Operation timing and metrics
- **Context Tracking** - Request and operation context

### Logging Features

Key features include:

- **Configurable Levels** - Set different levels for different components
- **Structured Output** - JSON and text formatting options
- **Performance Metrics** - Operation timing and resource usage
- **Context Propagation** - Track operations across components

## Configuration

### Global Logging Configuration

Configure logging in your `spooky.hcl` file:

```hcl
logging {
  level = "info"
  format = "json"
  
  outputs {
    console {
      enabled = true
      level = "info"
    }
    
    file {
      enabled = true
      path = "~/.local/state/spooky/logs/spooky.log"
      level = "debug"
      max_size_mb = 100
      max_files = 5
    }
  }
  
  performance {
    enabled = true
    threshold_ms = 1000
  }
}
```

### Console Output Configuration

Configure console logging:

```hcl
outputs {
  console {
    enabled = true
    level = "info"
    format = "text"
    colors = true
  }
}
```

### File Output Configuration

Configure file logging:

```hcl
outputs {
  file {
    enabled = true
    path = "~/.local/state/spooky/logs/spooky.log"
    level = "debug"
    format = "json"
    max_size_mb = 100
    max_files = 5
    compress = true
  }
}
```

### Network Output Configuration

Configure network logging:

```hcl
outputs {
  syslog {
    enabled = true
    address = "localhost:514"
    facility = "local0"
    level = "warn"
  }
  
  http {
    enabled = true
    url = "https://logs.example.com/ingest"
    level = "error"
    headers = {
      "Authorization" = "Bearer token"
    }
  }
}
```

## CLI Commands

### Logging Configuration

Manage logging configuration:

```bash
# Show current configuration
spooky logging config show

# Set log level
spooky logging config set --level debug

# Set output format
spooky logging config set --format json

# Enable performance logging
spooky logging config set --performance-enabled true
```

### Log Management

Manage log files:

```bash
# List log files
spooky logging files list

# View recent logs
spooky logging files tail

# Export logs
spooky logging files export --output logs.json

# Rotate logs
spooky logging files rotate
```

### Performance Monitoring

Monitor performance:

```bash
# View performance metrics
spooky logging performance show

# Export performance data
spooky logging performance export --output metrics.json

# Set performance threshold
spooky logging performance set --threshold 500
```

## Advanced Features

### Structured Logging

Enable structured logging for better analysis:

```hcl
logging {
  format = "json"
  
  fields {
    service = "spooky"
    version = "1.0.0"
    environment = "production"
  }
}
```

### Context Tracking

Track operations across components:

```hcl
logging {
  context {
    enabled = true
    fields = ["request_id", "user_id", "operation"]
  }
}
```

### Performance Monitoring

Monitor operation performance:

```hcl
logging {
  performance {
    enabled = true
    threshold_ms = 1000
    metrics = ["duration", "memory", "cpu"]
  }
}
```

### Log Aggregation

Configure log aggregation:

```hcl
outputs {
  elasticsearch {
    enabled = true
    url = "http://localhost:9200"
    index = "spooky-logs"
    level = "info"
  }
  
  fluentd {
    enabled = true
    address = "localhost:24224"
    tag = "spooky"
    level = "debug"
  }
}
```

## Security Best Practices

### Log Security

- Never log sensitive information (passwords, keys, tokens)
- Use appropriate log levels for different environments
- Implement log retention policies
- Monitor log access and tampering

### Access Control

- Restrict access to log files
- Use secure transport for network logging
- Implement log encryption for sensitive data
- Monitor log access patterns

### Compliance

- Implement audit logging for security events
- Maintain log retention for compliance
- Use structured logging for analysis
- Regular log reviews and monitoring

## Troubleshooting

### Common Logging Issues

**Logs Not Appearing**
```bash
# Check log level configuration
spooky logging config show

# Verify output configuration
spooky logging config show --outputs
```

**Performance Issues**
```bash
# Check performance metrics
spooky logging performance show

# Adjust performance threshold
spooky logging performance set --threshold 2000
```

**File Permission Issues**
```bash
# Check log directory permissions
ls -la ~/.local/state/spooky/logs/

# Fix permissions if needed
chmod 755 ~/.local/state/spooky/logs/
```

### Debugging Logging

Enable debug logging for troubleshooting:

```bash
# Set debug level
spooky logging config set --level debug

# View debug logs
spooky logging files tail --level debug
```

### Performance Optimization

Optimize logging performance:

```bash
# Disable performance logging in production
spooky logging config set --performance-enabled false

# Use async logging
spooky logging config set --async true

# Adjust buffer size
spooky logging config set --buffer-size 1024
```

## Integration with Other Systems

### Actions Integration

Logging integrates with the actions system:

```bash
# Run actions with detailed logging
spooky actions run ./my-project --log-level debug

# View action logs
spooky logging files tail --filter "action"
```

### Facts Integration

Logging tracks fact collection operations:

```bash
# Collect facts with logging
spooky facts gather ./my-project --log-level info

# View fact collection logs
spooky logging files tail --filter "facts"
```

### SSH Integration

Logging tracks SSH operations:

```bash
# Test SSH with logging
spooky machines ping ./my-project --log-level debug

# View SSH logs
spooky logging files tail --filter "ssh"
```

## Examples

### Basic Logging Configuration

```hcl
# spooky.hcl
logging {
  level = "info"
  format = "text"
  
  outputs {
    console {
      enabled = true
      level = "info"
    }
    
    file {
      enabled = true
      path = "~/.local/state/spooky/logs/spooky.log"
      level = "debug"
    }
  }
}
```

### Advanced Logging Configuration

```hcl
# spooky.hcl
logging {
  level = "debug"
  format = "json"
  
  outputs {
    console {
      enabled = true
      level = "info"
      format = "text"
      colors = true
    }
    
    file {
      enabled = true
      path = "~/.local/state/spooky/logs/spooky.log"
      level = "debug"
      format = "json"
      max_size_mb = 100
      max_files = 5
    }
    
    syslog {
      enabled = true
      address = "localhost:514"
      facility = "local0"
      level = "warn"
    }
  }
  
  performance {
    enabled = true
    threshold_ms = 1000
  }
  
  context {
    enabled = true
    fields = ["request_id", "operation"]
  }
}
```

### Project-Specific Logging

```hcl
# project.hcl
project {
  name = "my-project"
  
  logging {
    level = "debug"
    file {
      enabled = true
      path = "./logs/project.log"
    }
  }
}
```

## Best Practices

### Log Level Management

- Use DEBUG for development and troubleshooting
- Use INFO for normal operations
- Use WARN for potential issues
- Use ERROR for failed operations
- Use FATAL sparingly for critical failures

### Log Formatting

- Use structured logging (JSON) for production
- Use text formatting for development
- Include relevant context in log messages
- Use consistent field names across components

### Performance

- Monitor logging performance impact
- Use appropriate log levels for different environments
- Implement log rotation and retention
- Use async logging for high-throughput operations

### Security

- Never log sensitive information
- Implement log access controls
- Monitor log access patterns
- Use secure transport for network logging

## Next Steps

- Explore the [Logging API Reference](LOGGING_API_REFERENCE.md) for detailed technical information
- Check the [Logging Troubleshooting Guide](LOGGING_TROUBLESHOOTING.md) for common issues
- Review the [Logging Documentation Summary](LOGGING_DOCUMENTATION_SUMMARY.md) for implementation details
- Learn about [Logging Integration Patterns](INTEGRATIONS_USER_GUIDE.md) for advanced usage
