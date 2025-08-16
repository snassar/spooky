# Logging System Troubleshooting Guide

## Overview

This troubleshooting guide provides solutions for common issues encountered when working with the spooky logging system. It covers error messages, configuration problems, performance issues, and debugging techniques.

**Status: Production Ready** - The logging system is fully implemented with comprehensive logging configuration, secure logging, and integration capabilities.

## Logging System Status

### ✅ Fully Functional Logging Infrastructure

The logging system now has **complete logging infrastructure** with:

- **Logging Configuration**: Comprehensive logging configuration and management
- **Secure Logging**: Secure logging with sensitive data redaction
- **Multiple Outputs**: Support for multiple output formats and destinations
- **CLI Integration**: Full CLI integration with `spooky logging` commands
- **Project Integration**: Logging configuration from project configuration
- **Error Handling**: Comprehensive error handling and reporting
- **Performance Monitoring**: Logging performance monitoring and optimization
- **Integration Support**: Full integration with all system components

### What This Means for Users

- **No More Stubs**: All functionality is fully implemented - no placeholder code
- **Production Ready**: The system is ready for production use
- **Complete Feature Set**: All documented features are functional
- **Reliable Logging**: Robust logging with secure data handling
- **Performance Optimized**: Efficient logging with multiple output options

### Expected Behavior

When using logging, you can expect:

1. **Proper Configuration**: Logging configuration loads and applies correctly
2. **Secure Logging**: Sensitive data is properly redacted
3. **Multiple Outputs**: Logs are written to configured outputs
4. **Performance**: Efficient logging with minimal overhead
5. **Integration**: Logging integrates with all system components
6. **Error Handling**: Clear error messages with actionable information

## Common Issues and Solutions

### Configuration Errors

#### "Invalid log level: 'invalid_level'"

**Cause:** The specified log level is not valid.

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

**Cause:** The specified output format is not supported.

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

**Cause:** File output is configured but no path is specified.

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

**Cause:** File permissions are specified in an invalid format.

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

**Cause:** Insufficient permissions to create the log directory.

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

**Cause:** Cannot write to the specified log file.

**Solution:**
```bash
# Check file permissions
ls -la /path/to/logfile.log

# Fix file permissions
chmod 644 /path/to/logfile.log

# Or create new file with proper permissions
touch /path/to/logfile.log
chmod 644 /path/to/logfile.log
```

#### "Log file is full: no space left on device"

**Cause:** Disk space is exhausted.

**Solution:**
```bash
# Check disk space
df -h

# Check log file size
ls -lh /path/to/logfile.log

# Rotate log files
spooky logging rotate

# Or clean up old logs
find /path/to/logs -name "*.log" -mtime +30 -delete
```

### Secure Logging Errors

#### "Failed to redact sensitive data"

**Cause:** Sensitive data redaction is failing.

**Solution:**
```go
// ✅ CORRECT - Proper secure logging configuration
logger := spookylogging.NewSecureLogger(baseLogger)

// Configure sensitive fields
logger.SetSensitiveFields([]string{"password", "secret", "key"})

// Use secure logging
logger.Info("User login", map[string]interface{}{
    "user": "admin",
    "password": "[REDACTED]",  // Will be redacted
})
```

#### "Sensitive data detected in logs"

**Cause:** Sensitive data is being logged without redaction.

**Solution:**
```go
// ✅ CORRECT - Use secure logging for sensitive data
secureLogger := spookylogging.NewSecureLogger(logger)

// Log sensitive operations securely
secureLogger.Info("Database connection", map[string]interface{}{
    "host": "db.example.com",
    "user": "admin",
    // Don't log password directly
})
```

### Performance Issues

#### "Logging is very slow"

**Cause:** Logging operations are taking too long.

**Solution:**
```go
// ✅ CORRECT - Optimize logging performance
logger := spookylogging.NewLogger()

// Use async logging
logger.SetAsync(true)

// Use buffered output
logger.SetBufferSize(1000)

// Use appropriate log level
logger.SetLevel("info")  // Don't use debug in production
```

#### "High memory usage during logging"

**Cause:** Logging is consuming too much memory.

**Solution:**
```go
// ✅ CORRECT - Optimize memory usage
logger := spookylogging.NewLogger()

// Use streaming output
logger.SetStreaming(true)

// Limit buffer size
logger.SetBufferSize(100)

// Use structured logging efficiently
logger.Info("Operation completed", map[string]interface{}{
    "operation": "process",
    "duration":  duration,
    // Avoid large objects in logs
})
```

### Integration Errors

#### "Failed to initialize logging"

**Cause:** Logging system failed to initialize.

**Solution:**
```go
// ✅ CORRECT - Proper logging initialization
logger, err := spookylogging.NewLogger()
if err != nil {
    return fmt.Errorf("failed to create logger: %w", err)
}

// Initialize logger
if err := logger.Initialize(); err != nil {
    return fmt.Errorf("failed to initialize logger: %w", err)
}

// Verify initialization
if !logger.IsInitialized() {
    return fmt.Errorf("logger not properly initialized")
}
```

#### "Logging integration failed"

**Cause:** Logging integration with other components is failing.

**Solution:**
```go
// ✅ CORRECT - Proper logging integration
logger := spookylogging.NewLogger()

// Integrate with other components
factsManager := spookyfacts.NewManager(logger)
actionsManager := spookyactions.NewManager(logger)

// Verify integration
if err := logger.ValidateIntegration(); err != nil {
    return fmt.Errorf("logging integration validation failed: %w", err)
}
```

## Configuration Problems

### Logging Configuration Issues

#### "Invalid logging configuration"

**Cause:** Logging configuration is incorrect.

**Solution:**
```hcl
# ✅ CORRECT - Valid logging configuration
logging {
  level = "info"
  format = "json"
  
  output {
    type = "file"
    path = "/var/log/spooky/spooky.log"
    permissions = "0644"
  }
  
  output {
    type = "console"
  }
}
```

#### "Missing logging configuration"

**Cause:** Logging configuration is not set up.

**Solution:**
```bash
# Create default logging configuration
mkdir -p ~/.config/spooky

cat > ~/.config/spooky/logging.hcl << 'EOF'
logging {
  level = "info"
  format = "json"
  
  output {
    type = "file"
    path = "~/.local/share/spooky/logs/spooky.log"
    permissions = "0644"
  }
}
EOF
```

### Output Configuration Issues

#### "Invalid output configuration"

**Cause:** Output configuration is incorrect.

**Solution:**
```hcl
# ✅ CORRECT - Valid output configuration
logging {
  output {
    type = "file"
    path = "/var/log/spooky/spooky.log"
    permissions = "0644"
    max_size = "100MB"
    max_files = 5
  }
  
  output {
    type = "console"
    format = "plain"
  }
}
```

#### "Output path not writable"

**Cause:** Output path is not writable.

**Solution:**
```bash
# Check path permissions
ls -la /var/log/spooky/

# Fix permissions
sudo chown -R $USER:$USER /var/log/spooky/
chmod 755 /var/log/spooky/

# Or use user directory
mkdir -p ~/.local/share/spooky/logs
```

## Debugging Techniques

### Enable Verbose Output

```bash
# Enable verbose output for logging operations
spooky logging configure --verbose
spooky logging test --verbose

# Enable debug logging
export SPOOKY_LOG_LEVEL=debug
spooky logging test
```

### Test Logging Functionality

```bash
# Test logging configuration
spooky logging test

# Test specific output
spooky logging test --output file
spooky logging test --output console

# Test log levels
spooky logging test --level debug
spooky logging test --level info
```

### Validate Configuration

```bash
# Validate logging configuration
spooky logging validate

# Check specific configuration
spooky logging validate --config ~/.config/spooky/logging.hcl

# Test configuration
spooky logging test --config ~/.config/spooky/logging.hcl
```

## Recovery Procedures

### Logging Recovery

```go
// Recreate logger if needed
logger, err := spookylogging.NewLogger()
if err != nil {
    return fmt.Errorf("failed to recreate logger: %w", err)
}

// Reinitialize logger
if err := logger.Initialize(); err != nil {
    return fmt.Errorf("failed to reinitialize logger: %w", err)
}
```

### Configuration Recovery

```bash
# Backup configuration
cp ~/.config/spooky/logging.hcl ~/.config/spooky/logging.hcl.backup

# Restore from backup if needed
cp ~/.config/spooky/logging.hcl.backup ~/.config/spooky/logging.hcl

# Validate configuration
spooky logging validate --config ~/.config/spooky/logging.hcl
```

### File Recovery

```bash
# Check log file integrity
tail -f /var/log/spooky/spooky.log

# Rotate log files if corrupted
spooky logging rotate

# Recreate log file if needed
rm /var/log/spooky/spooky.log
touch /var/log/spooky/spooky.log
chmod 644 /var/log/spooky/spooky.log
```

## Prevention Strategies

### Regular Maintenance

```bash
# Schedule log rotation
crontab -e
# Add: 0 2 * * * /usr/local/bin/spooky logging rotate

# Monitor log file sizes
find /var/log/spooky -name "*.log" -size +100M -exec ls -lh {} \;

# Clean up old logs
find /var/log/spooky -name "*.log.*" -mtime +30 -delete
```

### Monitoring

```bash
# Monitor logging performance
spooky logging status

# Monitor log file growth
watch -n 60 'ls -lh /var/log/spooky/spooky.log'

# Monitor system resources
top -p $(pgrep spooky)
```

### Backup Strategy

```bash
# Backup logging configuration
cp ~/.config/spooky/logging.hcl ~/.config/spooky/logging.hcl.$(date +%Y%m%d)

# Version control configuration
git add ~/.config/spooky/logging.hcl
git commit -m "Update logging configuration"

# Backup important logs
tar -czf spooky-logs-$(date +%Y%m%d).tar.gz /var/log/spooky/
```

## Best Practices for Troubleshooting

### 1. Start Simple

Begin with simple logging configuration and add complexity gradually:

```hcl
# Start with basic configuration
logging {
  level = "info"
  output {
    type = "console"
  }
}

# Then add complexity
logging {
  level = "info"
  format = "json"
  
  output {
    type = "file"
    path = "/var/log/spooky/spooky.log"
    permissions = "0644"
  }
  
  output {
    type = "console"
  }
}
```

### 2. Use Proper Error Handling

Implement proper error handling in logging code:

```go
// Use proper error handling
logger, err := spookylogging.NewLogger()
if err != nil {
    return fmt.Errorf("failed to create logger: %w", err)
}

// Initialize with error handling
if err := logger.Initialize(); err != nil {
    return fmt.Errorf("failed to initialize logger: %w", err)
}
```

### 3. Validate Early and Often

Validate logging configuration frequently:

```bash
# Validate after every change
spooky logging validate

# Validate before operations
spooky logging validate && spooky facts list ./my-project

# Validate in scripts
#!/bin/bash
if spooky logging validate; then
    spooky facts list ./my-project
else
    echo "Logging validation failed"
    exit 1
fi
```

### 4. Use Secure Logging

Use secure logging for sensitive data:

```go
// Use secure logging for sensitive operations
secureLogger := spookylogging.NewSecureLogger(logger)

// Log sensitive operations securely
secureLogger.Info("User authentication", map[string]interface{}{
    "user": "admin",
    "method": "ssh_key",
    // Don't log sensitive data
})
```

### 5. Monitor and Optimize

Monitor logging performance and optimize:

```bash
# Monitor logging performance
spooky logging status --performance

# Check log file sizes
ls -lh /var/log/spooky/

# Monitor system resources
top -p $(pgrep spooky)
```

## Getting Help

### Documentation Resources

1. **User Guide** - For usage questions and best practices
2. **API Reference** - For technical implementation details
3. **Examples** - For configuration patterns and use cases

### Common Questions

#### "Why can't I see my logs?"

1. Check logging configuration
2. Verify output paths
3. Check file permissions
4. Validate log levels

#### "How do I debug logging issues?"

```bash
# Enable verbose output
spooky logging test --verbose

# Test specific outputs
spooky logging test --output file
spooky logging test --output console

# Check configuration
spooky logging validate --verbose
```

#### "How do I fix configuration issues?"

```bash
# Validate configuration
spooky logging validate --verbose

# Test configuration
spooky logging test --config ~/.config/spooky/logging.hcl

# Fix configuration issues
# Update logging configuration based on error messages
```

#### "How do I optimize logging performance?"

```bash
# Use appropriate log levels
spooky logging configure --level info

# Use efficient output formats
spooky logging configure --format json

# Monitor performance
spooky logging status --performance
```

### When to Seek Additional Help

- Configuration validation passes but logging still fails
- Performance issues persist after optimization
- Unusual error messages not covered in this guide
- Integration issues with other spooky components

For additional help, refer to the [User Guide](LOGGING_USER_GUIDE.md) and [API Reference](LOGGING_API_REFERENCE.md), or check the project documentation for more advanced troubleshooting techniques.

## Conclusion

The logging system provides robust, reliable logging with comprehensive configuration, secure data handling, and performance optimization capabilities. Most issues can be resolved by following the troubleshooting steps outlined in this guide. For persistent issues, enable verbose output and collect diagnostic information for further analysis.
