# Logging Configuration Strategy

## Overview

The spooky logging framework uses a **hierarchical configuration strategy** that allows for both global and project-specific logging configuration. This provides flexibility while maintaining consistency across the system.

## Configuration Hierarchy

### 1. **Global Configuration** (Primary)
**Location**: `~/.config/spooky/logging.hcl` (or `$XDG_CONFIG_HOME/spooky/logging.hcl`)

This is the **primary** logging configuration that applies to all spooky operations unless overridden by project-specific settings.

### 2. **Project-Specific Overrides** (Optional)
**Location**: `project.hcl` within individual project directories

Projects can optionally override global logging settings for development, debugging, or specific requirements.

## Configuration Loading Order

The logging system loads configuration in the following order:

1. **Default Configuration** (built-in)
2. **Global Configuration** (`~/.config/spooky/logging.hcl`)
3. **Project Configuration** (`project.hcl` - if exists)

Each level can override settings from the previous level.

## File Locations

### Global Configuration
```
~/.config/spooky/
├── spooky.hcl          # Main CLI configuration
└── logging.hcl         # Global logging configuration
```

### Project Configuration
```
my-project/
├── project.hcl         # Project configuration (includes optional logging section)
├── machines.hcl        # Machine inventory
├── actions.hcl         # Action definitions
└── logs/               # Project-specific log directory (if configured)
    └── project.log
```

## Setup Instructions

### 1. Create Global Configuration

Use the setup utility to create the default global configuration:

```bash
go run examples/setup-logging-config-utility.go
```

This creates `~/.config/spooky/logging.hcl` with sensible defaults.

### 2. Customize Global Configuration

Edit the global configuration file:

```bash
# Edit global logging configuration
nano ~/.config/spooky/logging.hcl
```

### 3. Optional: Add Project-Specific Configuration

Add a `logging` section to your `project.hcl`:

```hcl
project {
  name = "my-project"
  # ... other project settings

  # Override global logging for this project
  logging {
    level = "debug"  # More verbose for development
    output = "file"
    file {
      path = "./logs/project.log"
    }
  }
}
```

## Configuration Examples

### Global Configuration (Production)

```hcl
# ~/.config/spooky/logging.hcl
logging {
  level = "info"
  format = "json"
  output = "file"
  
  file {
    path        = "/var/log/spooky/spooky.log"
    permissions = "0644"
    append      = true
  }
  
  filtering {
    components = {
      "ssh"     = "warn"    # Reduce SSH verbosity
      "facts"   = "info"    # Standard facts logging
      "actions" = "info"    # Standard action logging
    }
  }
  
  rotation {
    enabled      = true
    max_size     = "100MB"
    max_age      = "30d"
    max_backups  = 5
    compress     = true
  }
}
```

### Global Configuration (Development)

```hcl
# ~/.config/spooky/logging.hcl
logging {
  level = "debug"
  format = "text"
  output = "stderr"
  
  filtering {
    components = {
      "ssh"     = "debug"   # Verbose SSH logging
      "facts"   = "debug"   # Verbose facts logging
      "actions" = "debug"   # Verbose action logging
    }
  }
}
```

### Project-Specific Override (Development Project)

```hcl
# my-development-project/project.hcl
project {
  name = "development-project"
  
  # Override global logging for this project
  logging {
    level = "debug"
    output = "file"
    file {
      path   = "./logs/dev.log"
      append = false  # Start fresh each run
    }
    
    filtering {
      components = {
        "project" = "debug"   # Very verbose project logging
        "ssh"     = "debug"   # Very verbose SSH logging
      }
    }
    
    rotation {
      enabled = false  # Keep all logs for debugging
    }
  }
}
```

### Project-Specific Override (Production Project)

```hcl
# my-production-project/project.hcl
project {
  name = "production-project"
  
  # Override global logging for this project
  logging {
    level = "warn"  # Only warnings and errors
    output = "file"
    file {
      path        = "/var/log/spooky/production.log"
      permissions = "0640"
      append      = true
    }
    
    filtering {
      components = {
        "ssh"     = "error"   # Only SSH errors
        "facts"   = "warn"    # Only fact warnings/errors
        "actions" = "warn"    # Only action warnings/errors
      }
    }
  }
}
```

## Configuration Merging Strategy

When both global and project configurations exist, the system merges them with the following rules:

### **Global Configuration** (Base)
- Provides default values for all settings
- Sets up production-ready defaults
- Configures system-wide behavior

### **Project Configuration** (Override)
- Can override any setting from global configuration
- Only specified settings are overridden
- Unspecified settings inherit from global configuration

### **Merging Examples**

#### Example 1: Project Overrides Level Only
```hcl
# Global: level = "info", output = "file", format = "json"
# Project: level = "debug"

# Result: level = "debug", output = "file", format = "json"
```

#### Example 2: Project Overrides Multiple Settings
```hcl
# Global: level = "info", output = "file", format = "json"
# Project: level = "debug", output = "stderr"

# Result: level = "debug", output = "stderr", format = "json"
```

#### Example 3: Project Uses Global Settings
```hcl
# Global: level = "info", output = "file", format = "json"
# Project: (no logging section)

# Result: level = "info", output = "file", format = "json"
```

## Best Practices

### 1. **Global Configuration**
- Set **production-ready defaults** in global configuration
- Use **file output** with rotation for production
- Configure **component filtering** for noise reduction
- Set **appropriate log levels** for different components

### 2. **Project Configuration**
- Use **debug level** for development projects
- Use **project-specific log files** for isolation
- **Disable rotation** for development debugging
- **Override only necessary settings**

### 3. **Security Considerations**
- Use **appropriate file permissions** (0644 for logs, 0640 for sensitive)
- Configure **sensitive field redaction** in global config
- **Separate log files** for different environments
- **Rotate logs regularly** in production

### 4. **Performance Considerations**
- Use **buffering** for high-throughput scenarios
- Configure **async logging** for non-blocking operations
- **Filter unnecessary components** to reduce log volume
- Use **appropriate log levels** to control verbosity

## Environment Variables

The logging system respects the following environment variables:

- `XDG_CONFIG_HOME`: Override the global configuration directory
- `SPOOKY_LOG_LEVEL`: Override the log level (for debugging)
- `SPOOKY_LOG_OUTPUT`: Override the output destination (for debugging)

## Troubleshooting

### Common Issues

#### 1. **No Logs Appearing**
- Check if global configuration exists: `ls ~/.config/spooky/logging.hcl`
- Verify log level is appropriate for your components
- Check if output destination is accessible

#### 2. **Permission Denied**
- Ensure log directory exists and is writable
- Check file permissions on log files
- Verify user has write access to log directory

#### 3. **Too Much Logging**
- Increase log level (info → warn → error)
- Use component filtering to reduce specific component verbosity
- Configure pattern-based filtering

#### 4. **Too Little Logging**
- Decrease log level (error → warn → info → debug)
- Remove component filtering restrictions
- Check if project configuration is overriding global settings

### Debugging Commands

```bash
# Check if global configuration exists
ls -la ~/.config/spooky/logging.hcl

# View current global configuration
cat ~/.config/spooky/logging.hcl

# Test logging configuration
go run examples/logging-demo.go

# Check log file permissions
ls -la /var/log/spooky/

# Monitor logs in real-time
tail -f /var/log/spooky/spooky.log
```

## Migration Guide

### From Previous Logging System

1. **Create global configuration** using the setup utility
2. **Migrate existing log settings** to the new HCL format
3. **Test configuration** with the logging demo
4. **Update project configurations** to use the new format
5. **Verify log output** in new locations

### Configuration Mapping

| Old Setting | New Location | Example |
|-------------|--------------|---------|
| Log level | `logging.level` | `level = "info"` |
| Output file | `logging.file.path` | `file { path = "/var/log/spooky.log" }` |
| Format | `logging.format` | `format = "json"` |
| Component filtering | `logging.filtering.components` | `components = { "ssh" = "debug" }` |

## Conclusion

The hierarchical logging configuration strategy provides:

- **Consistency**: Global defaults ensure consistent behavior
- **Flexibility**: Project-specific overrides allow customization
- **Maintainability**: Centralized configuration management
- **Security**: Proper file permissions and sensitive data handling
- **Performance**: Configurable buffering and filtering options

This approach follows the principle of "sensible defaults with optional customization" and integrates seamlessly with the existing spooky project structure.
