# Auto-Setup Configuration System

## Overview

The spooky auto-setup configuration system automatically creates and validates the necessary configuration files and directories when running any CLI command (excluding `--version` and `--help`). This ensures that spooky is properly configured before use and provides clear error messages for configuration issues.

## Current Implementation Status

### ✅ Fully Implemented

The auto-setup system is **fully implemented** and functional:

- **OS Detection**: Automatically detects the operating system (Windows, macOS, Linux, BSD)
- **XDG Directory Creation**: Creates OS-specific spooky config directories using XDG standards
- **Default Configuration**: Generates default `spooky.hcl` and `logging.hcl` files if missing
- **HCL Validation**: Validates existing configuration files for valid HCL syntax
- **Error Reporting**: Provides clear error messages for configuration validation failures
- **CLI Integration**: Integrated with all CLI commands (except `--version` and `--help`)

### What This Means for Users

- **Automatic Setup**: No manual configuration required - spooky sets up everything automatically
- **Cross-Platform Support**: Works on Windows, macOS, Linux, and BSD systems
- **XDG Compliance**: Follows XDG base directory specification for configuration storage
- **Validation**: Ensures configuration files are valid before proceeding
- **Clear Errors**: Provides actionable error messages for configuration issues

## Implementation Details

### OS Detection and Directory Creation

The system automatically detects the operating system and creates the appropriate configuration directory:

```go
// OS-specific configuration paths
func getConfigPath() string {
    switch runtime.GOOS {
    case "windows":
        return filepath.Join(os.Getenv("APPDATA"), "spooky")
    case "darwin":
        return filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "spooky")
    default: // Linux, BSD, etc.
        if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
            return filepath.Join(xdgConfig, "spooky")
        }
        return filepath.Join(os.Getenv("HOME"), ".config", "spooky")
    }
}
```

### Default Configuration Generation

When configuration files are missing, the system generates sensible defaults:

#### Default `spooky.hcl`
```hcl
# Global spooky configuration
# This file configures global settings for spooky

# SSH configuration
ssh {
  default_port = 22
  default_timeout = 30
  max_connections = 10
  max_retry_attempts = 3
  retry_delay = 5
  idle_timeout = 300
}

# Facts configuration
facts {
  storage_format = "badgerdb"
  storage_path = "~/.local/state/spooky/facts.db"
  collection_timeout = 60
  parallel_workers = 4
}

# Actions configuration
actions {
  default_timeout = 300
  max_parallel = 10
  retry_attempts = 3
  retry_delay = 5
}

# Variables configuration
variables {
  encryption_enabled = false
  default_scope = "project"
}

# Templates configuration
templates {
  default_engine = "go"
  cache_enabled = true
  cache_size = 100
}
```

#### Default `logging.hcl`
```hcl
# Global logging configuration for spooky
# This file configures logging behavior for all spooky operations

logging {
  # Log level (debug, info, warn, error, fatal)
  level = "info"
  
  # Output format (json, text, structured)
  format = "json"
  
  # Output destination (stdout, stderr, file, null)
  output = "stderr"
  
  # Component-specific filtering
  filtering {
    components = {
      # "ssh"     = "debug"
      # "facts"   = "info"
      # "actions" = "warn"
    }
  }
  
  # Performance optimization
  performance {
    buffer {
      enabled        = false
      size           = 4096
      flush_interval = "1s"
    }
    
    async {
      enabled       = false
      queue_size    = 1000
      workers       = 1
      drop_when_full = false
    }
  }
}
```

### HCL Validation

The system validates all configuration files for proper HCL syntax:

```go
func validateHCLFile(filePath string) error {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("failed to read config file: %w", err)
    }
    
    // Parse HCL to validate syntax
    var config interface{}
    if err := hcl.Unmarshal(data, &config); err != nil {
        return fmt.Errorf("invalid HCL syntax in %s: %w", filePath, err)
    }
    
    return nil
}
```

## Usage

### Automatic Setup

The auto-setup system runs automatically when you execute any spooky command:

```bash
# First run - auto-setup will create configuration
spooky facts gather ./my-project

# Subsequent runs - configuration is validated
spooky actions run ./my-project
```

### Manual Configuration

You can also manually create or modify configuration files:

```bash
# Create configuration directory
mkdir -p ~/.config/spooky

# Create custom spooky.hcl
cat > ~/.config/spooky/spooky.hcl << 'EOF'
ssh {
  default_timeout = 60
  max_connections = 20
}

facts {
  parallel_workers = 8
}
EOF

# Create custom logging.hcl
cat > ~/.config/spooky/logging.hcl << 'EOF'
logging {
  level = "debug"
  format = "structured"
  output = "file"
  
  file {
    path = "/var/log/spooky/spooky.log"
  }
}
EOF
```

## Configuration File Locations

### Linux/macOS
- **Global Config**: `~/.config/spooky/spooky.hcl`
- **Global Logging**: `~/.config/spooky/logging.hcl`
- **State Directory**: `~/.local/state/spooky/`

### Windows
- **Global Config**: `%APPDATA%\spooky\spooky.hcl`
- **Global Logging**: `%APPDATA%\spooky\logging.hcl`
- **State Directory**: `%APPDATA%\spooky\state\`

### XDG Override
If `XDG_CONFIG_HOME` is set, it takes precedence:
- **Global Config**: `$XDG_CONFIG_HOME/spooky/spooky.hcl`
- **Global Logging**: `$XDG_CONFIG_HOME/spooky/logging.hcl`

## Error Handling

### Configuration Validation Errors

When configuration files have invalid HCL syntax, spooky provides clear error messages:

```bash
$ spooky facts gather ./my-project
Error: invalid HCL syntax in ~/.config/spooky/spooky.hcl: unexpected token '}' at line 5
```

### Missing Directory Errors

If the configuration directory cannot be created:

```bash
$ spooky facts gather ./my-project
Error: failed to create config directory /path/to/config: permission denied
```

### File Permission Errors

If configuration files have incorrect permissions:

```bash
$ spooky facts gather ./my-project
Error: config file ~/.config/spooky/spooky.hcl has incorrect permissions (should be 600)
```

## Best Practices

### Configuration Management

1. **Use Version Control**: Keep your configuration files in version control for consistency
2. **Environment-Specific Configs**: Use different configurations for development, staging, and production
3. **Secure Permissions**: Ensure configuration files have appropriate permissions (600 for sensitive configs)
4. **Regular Validation**: Run `spooky config validate` to check configuration syntax

### Logging Configuration

1. **Development**: Use `level = "debug"` and `format = "structured"` for detailed output
2. **Production**: Use `level = "warn"` and `format = "json"` for performance and log aggregation
3. **File Rotation**: Configure log rotation for production environments
4. **Component Filtering**: Use component-specific log levels to reduce noise

### SSH Configuration

1. **Connection Limits**: Set appropriate `max_connections` based on your infrastructure
2. **Timeouts**: Configure timeouts based on network conditions and command complexity
3. **Retry Logic**: Use retry settings for unreliable network connections
4. **Key Management**: Ensure SSH keys are properly configured and secured

## Troubleshooting

### Common Issues

#### Configuration Not Found
```bash
Error: configuration file not found
```
**Solution**: The auto-setup system should create this automatically. Check if the process has write permissions to the config directory.

#### Invalid HCL Syntax
```bash
Error: invalid HCL syntax in ~/.config/spooky/spooky.hcl
```
**Solution**: Check the HCL syntax in your configuration file. Use `spooky config validate` to identify specific issues.

#### Permission Denied
```bash
Error: permission denied creating config directory
```
**Solution**: Ensure the user running spooky has write permissions to the parent directory.

#### XDG Configuration Issues
```bash
Error: XDG_CONFIG_HOME is set but directory is not writable
```
**Solution**: Check that `XDG_CONFIG_HOME` points to a writable directory, or unset it to use the default location.

### Debugging Configuration

Use the `--debug` flag to see detailed information about configuration loading:

```bash
spooky facts gather ./my-project --debug
```

This will show:
- Configuration file locations being checked
- Default values being applied
- Validation steps being performed
- Any warnings or errors encountered

## Integration with Other Systems

### Project System Integration
- **Project Overrides**: Project-specific configuration can override global settings
- **Validation**: Project configuration is validated against schemas
- **Isolation**: Each project can have its own configuration settings

### CLI System Integration
- **Command Integration**: Auto-setup runs before any command execution
- **Error Handling**: Configuration errors prevent command execution with clear messages
- **Help Integration**: Configuration help is available via `spooky config --help`

### Schema System Integration
- **Schema Validation**: Configuration files are validated against embedded schemas
- **Type Safety**: Configuration values are type-checked at runtime
- **Documentation**: Schema provides documentation for all configuration options

## Future Enhancements

### Planned Features

1. **Configuration Migration**: Automatic migration of configuration between versions
2. **Configuration Templates**: Pre-built configuration templates for common use cases
3. **Configuration Backup**: Automatic backup of configuration before changes
4. **Configuration Sync**: Synchronization of configuration across multiple machines
5. **Configuration Encryption**: Support for encrypted configuration files

### Integration Enhancements

1. **Environment Variables**: Enhanced support for environment variable overrides
2. **Configuration Inheritance**: Hierarchical configuration inheritance
3. **Configuration Validation**: Enhanced validation with custom rules
4. **Configuration Monitoring**: Monitoring of configuration changes and their impact

## Summary

The auto-setup configuration system ensures that spooky is properly configured before use, providing a seamless experience for users while maintaining security and validation. The system automatically creates necessary directories and files, validates configuration syntax, and provides clear error messages for any issues encountered.

**Status**: ✅ **Production Ready** - The auto-setup system is fully implemented and ready for production use.
