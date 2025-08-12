# Logging Examples

This directory contains examples of logging configurations for the spooky logging system. These examples demonstrate various patterns and best practices for configuring logging in different environments and use cases.

## Example Files

### [`global-logging-config.hcl`](global-logging-config.hcl)
**Purpose:** Global logging configuration for system-wide settings

**Use Case:** Place this file at `$XDG_CONFIG_HOME/spooky/logging.hcl` (typically `~/.config/spooky/logging.hcl`) to configure logging for all spooky operations.

**Features:**
- System-wide logging defaults
- Component-specific configuration
- Performance and audit logging
- File rotation and compression
- Multiple output destinations

**Best For:** Setting up logging for your entire spooky installation.

### [`project-logging-config.hcl`](project-logging-config.hcl)
**Purpose:** Project-specific logging configuration that overrides global settings

**Use Case:** Place this file as `./logging.hcl` in your project directory to customize logging for that specific project.

**Features:**
- Overrides global logging settings
- Project-specific log files
- Component-level customization
- Error-only logging
- Development-friendly configuration

**Best For:** Customizing logging for individual projects while maintaining global defaults.

### [`logging-formats.hcl`](logging-formats.hcl)
**Purpose:** Examples of different output format configurations

**Use Case:** Reference this file to understand how to configure different output formats (JSON, structured, plain text) for various scenarios.

**Features:**
- JSON format for machine processing
- Structured format for human readability
- Plain text format for simplicity
- Development vs production configurations
- Component-specific format overrides

**Best For:** Understanding output format options and choosing the right format for your needs.

## Using the Examples

### Getting Started

1. **Copy the appropriate example** to your desired location:
   ```bash
   # For global configuration
   cp global-logging-config.hcl ~/.config/spooky/logging.hcl
   
   # For project configuration
   cp project-logging-config.hcl ./logging.hcl
   ```

2. **Customize the configuration** for your specific needs:
   - Adjust log levels
   - Change file paths
   - Modify component settings
   - Configure output formats

3. **Validate your configuration**:
   ```bash
   spooky logging validate
   ```

4. **Test the configuration**:
   ```bash
   spooky logging test
   ```

### Configuration Patterns

#### Basic Setup
```hcl
logging {
  level = "info"
  format = "structured"
  
  output {
    type = "console"
    enabled = true
  }
}
```

#### Development Configuration
```hcl
logging {
  level = "debug"
  format = "structured"
  
  output {
    type = "console"
    enabled = true
    colorize = true
  }
  
  output {
    type = "file"
    enabled = true
    path = "./logs/debug.log"
  }
}
```

#### Production Configuration
```hcl
logging {
  level = "warn"
  format = "json"
  
  output {
    type = "file"
    enabled = true
    path = "/var/log/spooky/app.log"
    
    file {
      max_size = "100MB"
      max_age = "30d"
      max_backups = 10
      compress = true
    }
  }
}
```

## Best Practices

### 1. Start Simple
Begin with a basic configuration and add complexity as needed:
```hcl
logging {
  level = "info"
  format = "structured"
  
  output {
    type = "console"
    enabled = true
  }
}
```

### 2. Use Appropriate Log Levels
- **Development**: `debug` or `trace` for detailed information
- **Staging**: `info` for general information
- **Production**: `warn` or `error` for important events only

### 3. Configure Log Rotation
Always configure log rotation to prevent disk space issues:
```hcl
file {
  max_size = "100MB"
  max_age = "30d"
  max_backups = 10
  compress = true
}
```

### 4. Use Component-Specific Configuration
Configure logging per component to isolate issues:
```hcl
components {
  "facts" {
    level = "debug"
    enabled = true
  }
  
  "machines" {
    level = "info"
    enabled = true
  }
}
```

### 5. Choose the Right Format
- **JSON**: For log aggregation and machine processing
- **Structured**: For human-readable logs with structured fields
- **Plain**: For simple text output

## Environment-Specific Configurations

### Development Environment
- Use `debug` or `trace` log levels
- Enable console output with colors
- Use structured format for readability
- Configure component-specific logging

### Staging Environment
- Use `info` log level
- Enable both console and file output
- Use JSON format for analysis
- Configure performance logging

### Production Environment
- Use `warn` or `error` log levels
- Disable console output
- Use JSON format for aggregation
- Configure aggressive log rotation
- Enable audit logging

## Testing and Validation

### Validate Configuration
```bash
# Validate global configuration
spooky logging validate

# Validate project configuration
spooky logging validate --project ./
```

### Test Logging
```bash
# Test current configuration
spooky logging test

# Test with specific level
spooky logging test --level debug

# Test specific component
spooky logging test --component facts
```

### Monitor Logs
```bash
# View recent logs
tail -f ./logs/project.log

# Search for errors
grep -i error ./logs/project.log

# Monitor log file sizes
du -sh ./logs/*.log
```

## Troubleshooting

### Common Issues

1. **Logs not appearing**: Check log level and output configuration
2. **Permission errors**: Verify file and directory permissions
3. **Large log files**: Configure log rotation
4. **Performance issues**: Use higher log levels in production

### Debugging Steps

1. **Check configuration**:
   ```bash
   spooky logging show
   ```

2. **Validate configuration**:
   ```bash
   spooky logging validate
   ```

3. **Test logging**:
   ```bash
   spooky logging test
   ```

4. **Check file permissions**:
   ```bash
   ls -la ./logs/
   ```

## Integration with Other Systems

### Facts System
```hcl
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
```

### Machines System
```hcl
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
```

### Variables System
```hcl
components {
  "variables" {
    level = "warn"
    enabled = true
  }
}
```

## Advanced Features

### Performance Logging
```hcl
performance {
  enabled = true
  threshold_ms = 1000
  include_memory = true
  include_cpu = true
}
```

### Audit Logging
```hcl
audit {
  enabled = true
  level = "info"
  include_user = true
  include_ip = true
  include_session = true
}
```

### Async Logging
```hcl
async {
  enabled = true
  workers = 4
  queue_size = 1000
}
```

## Conclusion

These examples provide a solid foundation for configuring logging in the spooky system. Start with a simple configuration and gradually add features as needed. Always validate and test your configuration to ensure it works correctly in your environment.

For more detailed information, refer to the [Logging User Guide](../../LOGGING_USER_GUIDE.md), [API Reference](../../LOGGING_API_REFERENCE.md), and [Troubleshooting Guide](../../LOGGING_TROUBLESHOOTING.md).
