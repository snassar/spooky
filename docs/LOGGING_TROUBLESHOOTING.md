# Logging System Troubleshooting Guide

## Overview

This troubleshooting guide provides solutions for common issues encountered when working with the spooky logging system. It covers error messages, configuration problems, performance issues, and debugging techniques.

## Common Error Messages

### Configuration Errors

#### "Invalid log level: 'invalid_level'"

**Problem:** The specified log level is not valid.

**Solution:**
```bash
# Valid log levels are: trace, debug, info, warn, error, fatal
spooky logging configure --level info
```

**Check your configuration:**
```hcl
logging {
  level = "info"  # Must be one of the valid levels
}
```

#### "Invalid output format: 'invalid_format'"

**Problem:** The specified output format is not supported.

**Solution:**
```bash
# Valid formats are: json, structured, plain
spooky logging configure --format json
```

**Check your configuration:**
```hcl
logging {
  format = "json"  # Must be one of the valid formats
}
```

#### "File output path is required"

**Problem:** File output is configured but no path is specified.

**Solution:**
```hcl
logging {
  output {
    type = "file"
    path = "/path/to/logfile.log"  # Path is required for file output
  }
}
```

#### "Invalid file permissions: 'invalid_perms'"

**Problem:** File permissions are specified in an invalid format.

**Solution:**
```hcl
logging {
  output {
    type = "file"
    path = "/path/to/logfile.log"
    permissions = "0644"  # Use octal format (e.g., "0644", "0600")
  }
}
```

### File System Errors

#### "Failed to create log directory: permission denied"

**Problem:** Insufficient permissions to create the log directory.

**Solution:**
```bash
# Check directory permissions
ls -la /path/to/log/directory

# Create directory with proper permissions
mkdir -p /path/to/log/directory
chmod 755 /path/to/log/directory

# Or use a directory you have write access to
spooky logging configure --file ~/logs/spooky.log
```

#### "Failed to open log file: permission denied"

**Problem:** Cannot write to the specified log file.

**Solution:**
```bash
# Check file permissions
ls -la /path/to/logfile.log

# Fix permissions if needed
chmod 644 /path/to/logfile.log

# Or use a file you have write access to
spooky logging configure --file ~/logs/spooky.log
```

#### "Disk space full"

**Problem:** No space left on device for log files.

**Solution:**
```bash
# Check disk space
df -h

# Clean up old log files
find /var/log -name "*.log" -mtime +30 -delete

# Configure log rotation
spooky logging configure \
  --file /var/log/spooky.log \
  --max-size 100MB \
  --max-age 7d \
  --max-backups 5
```

### Component Errors

#### "Component 'invalid_component' not found"

**Problem:** Trying to configure logging for a non-existent component.

**Solution:**
```bash
# Available components: facts, machines, variables, actions, templates, ssh, config, cli
spooky logging configure --component facts --level debug
```

#### "Component logging is disabled"

**Problem:** Component logging is disabled in configuration.

**Solution:**
```hcl
logging {
  components {
    "facts" {
      enabled = true  # Enable component logging
      level = "debug"
    }
  }
}
```

## Configuration Problems

### Global Configuration Issues

#### Configuration Not Loading

**Problem:** Global configuration file is not being loaded.

**Diagnosis:**
```bash
# Check if configuration file exists
ls -la ~/.config/spooky/logging.hcl

# Check configuration file syntax
spooky logging validate
```

**Solution:**
```bash
# Create default configuration
spooky logging configure --level info --format json

# Or manually create the file
mkdir -p ~/.config/spooky
cat > ~/.config/spooky/logging.hcl << 'EOF'
logging {
  level = "info"
  format = "json"
  
  output {
    type = "console"
    enabled = true
  }
}
EOF
```

#### Configuration Override Issues

**Problem:** Project-specific configuration is not overriding global settings.

**Diagnosis:**
```bash
# Check project configuration
ls -la ./logging.hcl

# Validate project configuration
spooky logging validate --project ./
```

**Solution:**
```hcl
# In your project directory (./logging.hcl)
logging {
  level = "debug"  # This should override global level
  format = "structured"
  
  output {
    type = "file"
    enabled = true
    path = "./logs/project.log"
  }
}
```

### Environment Variable Issues

#### Environment Variables Not Applied

**Problem:** Environment variables are not overriding configuration.

**Diagnosis:**
```bash
# Check environment variables
env | grep SPOOKY_LOG

# Test with explicit environment variables
SPOOKY_LOG_LEVEL=debug spooky logging test
```

**Solution:**
```bash
# Set environment variables correctly
export SPOOKY_LOG_LEVEL=debug
export SPOOKY_LOG_FORMAT=json
export SPOOKY_LOG_FILE=./debug.log

# Verify they are applied
spooky logging show
```

### Output Configuration Issues

#### Console Output Not Working

**Problem:** Logs are not appearing in the console.

**Diagnosis:**
```bash
# Test console output
spooky logging test --level debug

# Check if console output is enabled
spooky logging show
```

**Solution:**
```hcl
logging {
  output {
    type = "console"
    enabled = true  # Ensure console output is enabled
    colorize = true
  }
}
```

#### File Output Not Working

**Problem:** Logs are not being written to files.

**Diagnosis:**
```bash
# Check if file exists and is writable
ls -la /path/to/logfile.log

# Test file output
spooky logging test --level debug

# Check file permissions
stat /path/to/logfile.log
```

**Solution:**
```hcl
logging {
  output {
    type = "file"
    enabled = true  # Ensure file output is enabled
    path = "/path/to/logfile.log"
    permissions = "0644"
  }
}
```

## Performance Issues

### High Logging Overhead

**Problem:** Logging is causing performance degradation.

**Diagnosis:**
```bash
# Check current log level
spooky logging show

# Monitor system performance
top -p $(pgrep spooky)
```

**Solution:**
```hcl
# Use higher log levels in production
logging {
  level = "warn"  # Only log warnings and errors
  
  # Disable verbose components
  components {
    "facts" {
      level = "error"  # Only log errors
    }
  }
}
```

### Large Log Files

**Problem:** Log files are growing too large.

**Diagnosis:**
```bash
# Check log file sizes
du -sh /var/log/spooky/*.log

# Check log file ages
ls -la /var/log/spooky/*.log
```

**Solution:**
```hcl
# Configure log rotation
logging {
  output {
    type = "file"
    path = "/var/log/spooky/app.log"
    max_size = "100MB"      # Rotate when file reaches 100MB
    max_age = "7d"          # Keep logs for 7 days
    max_backups = 10        # Keep 10 backup files
    compress = true         # Compress old logs
  }
}
```

### Memory Usage Issues

**Problem:** Logging is consuming too much memory.

**Diagnosis:**
```bash
# Check memory usage
ps aux | grep spooky

# Monitor memory over time
watch -n 1 'ps aux | grep spooky'
```

**Solution:**
```hcl
# Use async logging for high-volume operations
logging {
  async {
    enabled = true
    workers = 4
    queue_size = 1000
  }
}
```

## Output Format Issues

### JSON Format Problems

**Problem:** JSON output is malformed or unreadable.

**Diagnosis:**
```bash
# Test JSON output
spooky logging test --level debug | jq .

# Check JSON syntax
spooky logging test --level debug | python -m json.tool
```

**Solution:**
```hcl
logging {
  format = "json"
  
  json {
    pretty_print = true      # Make JSON readable
    include_timestamp = true
    include_level = true
    include_component = true
  }
}
```

### Structured Format Problems

**Problem:** Structured output is not formatted correctly.

**Diagnosis:**
```bash
# Test structured output
spooky logging test --level debug

# Check if colorization is causing issues
spooky logging configure --format structured
```

**Solution:**
```hcl
logging {
  format = "structured"
  
  structured {
    colorize = false         # Disable colors if causing issues
    include_timestamp = true
    include_level = true
    include_component = true
  }
}
```

### Plain Text Format Problems

**Problem:** Plain text output is missing information.

**Diagnosis:**
```bash
# Test plain text output
spooky logging test --level debug

# Compare with other formats
spooky logging configure --format json && spooky logging test
spooky logging configure --format structured && spooky logging test
```

**Solution:**
```hcl
logging {
  format = "plain"
  
  plain {
    include_timestamp = true
    include_level = true
    include_component = true
  }
}
```

## Component-Specific Issues

### Facts Component Logging

**Problem:** Facts component is not logging properly.

**Diagnosis:**
```bash
# Check facts component configuration
spooky logging show

# Test facts component logging
spooky logging test --component facts --level debug

# Run facts command to generate logs
spooky facts gather ./my-project
```

**Solution:**
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

### Machines Component Logging

**Problem:** Machines component is not logging properly.

**Diagnosis:**
```bash
# Check machines component configuration
spooky logging show

# Test machines component logging
spooky logging test --component machines --level debug

# Run machines command to generate logs
spooky machines ping ./my-project
```

**Solution:**
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

### Variables Component Logging

**Problem:** Variables component is not logging properly.

**Diagnosis:**
```bash
# Check variables component configuration
spooky logging show

# Test variables component logging
spooky logging test --component variables --level debug

# Run variables command to generate logs
spooky variables list ./my-project
```

**Solution:**
```hcl
logging {
  components {
    "variables" {
      level = "warn"
      enabled = true
      
      output {
        type = "file"
        path = "./logs/variables.log"
      }
    }
  }
}
```

## Debugging Techniques

### Verbose Logging

Enable verbose logging to get more detailed information:

```bash
# Enable debug logging
spooky logging configure --level debug

# Test with verbose output
spooky logging test --level debug

# Run commands with debug logging
SPOOKY_LOG_LEVEL=debug spooky facts gather ./my-project
```

### Configuration Validation

Validate your configuration to catch issues early:

```bash
# Validate global configuration
spooky logging validate

# Validate project configuration
spooky logging validate --project ./my-project

# Validate with strict checking
spooky logging validate --strict
```

### Log File Analysis

Analyze log files to understand issues:

```bash
# View recent logs
tail -f /var/log/spooky/app.log

# Search for errors
grep -i error /var/log/spooky/app.log

# Search for specific component
grep "component=facts" /var/log/spooky/app.log

# Analyze log patterns
awk '{print $4}' /var/log/spooky/app.log | sort | uniq -c
```

### Performance Monitoring

Monitor logging performance:

```bash
# Monitor log file growth
watch -n 5 'ls -lh /var/log/spooky/*.log'

# Monitor system resources
top -p $(pgrep spooky)

# Check disk I/O
iotop -p $(pgrep spooky)
```

## Best Practices for Troubleshooting

### 1. Start with Basic Configuration

Begin with a simple configuration and add complexity gradually:

```hcl
# Start with basic configuration
logging {
  level = "info"
  format = "structured"
  
  output {
    type = "console"
    enabled = true
  }
}
```

### 2. Use Component-Specific Logging

Configure logging per component to isolate issues:

```hcl
logging {
  level = "warn"  # Default level
  
  components {
    "facts" {
      level = "debug"  # More verbose for facts
    }
    
    "machines" {
      level = "info"   # Standard level for machines
    }
  }
}
```

### 3. Test Configuration Changes

Always test configuration changes:

```bash
# Test before applying
spooky logging validate

# Test after applying
spooky logging test

# Test with real commands
spooky facts gather ./my-project
```

### 4. Monitor Log Performance

Keep an eye on logging performance:

```bash
# Monitor log file sizes
du -sh /var/log/spooky/*.log

# Monitor system resources
ps aux | grep spooky

# Check for log rotation
ls -la /var/log/spooky/*.log*
```

### 5. Use Appropriate Log Levels

Use appropriate log levels for different environments:

```hcl
# Development
logging {
  level = "debug"
}

# Staging
logging {
  level = "info"
}

# Production
logging {
  level = "warn"
}
```

## Common Solutions

### Quick Fixes

#### Reset to Default Configuration

```bash
# Remove existing configuration
rm ~/.config/spooky/logging.hcl

# Create default configuration
spooky logging configure --level info --format structured
```

#### Enable Console Logging Only

```bash
# Configure console-only logging
spooky logging configure \
  --level debug \
  --format structured
```

#### Disable Problematic Components

```hcl
logging {
  components {
    "facts" {
      enabled = false  # Disable facts logging
    }
    
    "machines" {
      enabled = false  # Disable machines logging
    }
  }
}
```

### Advanced Solutions

#### Custom Log Format

```hcl
logging {
  format = "custom"
  
  custom {
    template = "{{.Timestamp}} [{{.Level}}] {{.Component}}: {{.Message}}"
    include_fields = ["server", "duration_ms"]
  }
}
```

#### Multiple Output Destinations

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
    path = "./logs/app.log"
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

#### Performance-Optimized Configuration

```hcl
logging {
  level = "warn"  # Higher level for performance
  
  # Async logging
  async {
    enabled = true
    workers = 4
    queue_size = 1000
  }
  
  # Aggressive log rotation
  output {
    type = "file"
    path = "/var/log/spooky/app.log"
    max_size = "50MB"
    max_age = "3d"
    max_backups = 5
    compress = true
  }
}
```

## Getting Help

### Documentation Resources

1. **User Guide** - For usage questions and best practices
2. **API Reference** - For technical implementation details
3. **Examples** - For configuration patterns and use cases

### Common Questions

#### "Why aren't my logs appearing?"

1. Check log level configuration
2. Verify output is enabled
3. Check file permissions
4. Validate configuration syntax

#### "How do I enable debug logging?"

```bash
spooky logging configure --level debug
```

#### "How do I rotate log files?"

```hcl
logging {
  output {
    type = "file"
    path = "/var/log/spooky/app.log"
    max_size = "100MB"
    max_age = "7d"
    max_backups = 10
  }
}
```

#### "How do I disable logging for a component?"

```hcl
logging {
  components {
    "facts" {
      enabled = false
    }
  }
}
```

### When to Seek Additional Help

- Configuration validation passes but logs still don't appear
- Performance issues persist after optimization
- Unusual error messages not covered in this guide
- Integration issues with other spooky components

For additional help, refer to the [User Guide](LOGGING_USER_GUIDE.md) and [API Reference](LOGGING_API_REFERENCE.md), or check the project documentation for more advanced troubleshooting techniques.
