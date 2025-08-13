# Auto-Setup Configuration System

## Overview

The spooky CLI includes an automatic configuration setup system that ensures users have proper configuration from the start. This system automatically creates and validates the spooky configuration directory and files when running any spooky command (except `--version` and `--help`).

## How It Works

### 1. **OS Detection and Directory Selection**

The system automatically detects the operating system and creates the appropriate configuration directory:

- **Linux/BSD**: `~/.config/spooky/` (following XDG Base Directory Specification)
- **macOS**: `~/Library/Application Support/spooky/`
- **Windows**: `%APPDATA%\spooky\`

### 2. **Automatic Setup Process**

When running any spooky command (except `--version` and `--help`), the system:

1. **Detects the OS** and determines the appropriate config directory
2. **Checks if the config directory exists**
3. **If it doesn't exist**: Creates the directory and default configuration files
4. **If it exists**: Ensures required files exist and validates their HCL syntax
5. **If validation fails**: Emits an error and stops running

### 3. **Configuration Files Created**

The system creates two essential configuration files:

#### `spooky.hcl` - Main CLI Configuration
```hcl
# Spooky CLI Configuration
# This file contains global configuration for the spooky CLI tool

# Global CLI settings
cli {
  # Default timeout for operations (in seconds)
  default_timeout = 300
  
  # Maximum parallel operations
  max_parallel = 10
  
  # Default log level for CLI operations
  log_level = "info"
  
  # Enable colored output (if supported by terminal)
  colored_output = true
  
  # Show progress indicators for long-running operations
  show_progress = true
}

# SSH configuration
ssh {
  # Default SSH timeout (in seconds)
  timeout = 30
  
  # SSH connection retry attempts
  retry_attempts = 3
  
  # Delay between retry attempts (in seconds)
  retry_delay = 5
  
  # Enable SSH connection pooling
  connection_pooling = true
  
  # Maximum number of SSH connections to keep in pool
  max_connections = 20
}

# Facts collection configuration
facts {
  # Default facts collection timeout (in seconds)
  timeout = 60
  
  # Enable automatic facts collection
  auto_collect = false
  
  # Maximum parallel facts collection workers
  max_parallel = 5
  
  # Facts collection retry attempts
  retry_attempts = 3
  
  # Delay between facts collection retries (in seconds)
  retry_delay = 5
}

# Actions configuration
actions {
  # Default action timeout (in seconds)
  default_timeout = 300
  
  # Maximum parallel action runs
  max_parallel = 10
  
  # Enable dry-run mode by default
  dry_run_default = false
  
  # Validate actions before running
  validate_before_run = true
  
  # Create backups before making changes
  backup_before_changes = false
}

# Storage configuration
storage {
  # Default storage format for facts databases
  facts_format = "memory"
  
  # Enable compression for storage
  compression = true
  
  # Enable encryption for sensitive data
  encryption = false
}
```

#### `logging.hcl` - Global Logging Configuration
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
  
  # File output configuration (used when output = "file")
  # file {
  #   path        = "/var/log/spooky/spooky.log"
  #   permissions = "0644"
  #   append      = true
  # }
  
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
  
  # Log rotation (when using file output)
  # rotation {
  #   enabled      = true
  #   max_size     = "100MB"
  #   max_age      = "30d"
  #   max_backups  = 5
  #   compress     = true
  #   local_time   = false
  # }
}
```

## Validation System

### HCL Syntax Validation

The system performs comprehensive HCL syntax validation with enhanced error reporting:

1. **Balanced Braces**: Ensures all opening braces have corresponding closing braces
2. **Basic Structure**: Validates that HCL blocks are properly structured
3. **Assignment Syntax**: Checks that assignments follow `key = value` format
4. **Line-by-Line Validation**: Validates each non-comment line for proper syntax with line-specific error reporting
5. **Multi-Level Validation**: Performs file existence, read access, and syntax validation in sequence

### Validation Examples

#### Valid HCL
```hcl
cli {
  default_timeout = 300
  max_parallel = 10
}
```

#### Invalid HCL (Caught by Validation)
```hcl
cli {
  default_timeout = 300
  max_parallel = 10
  # Missing closing brace
```

**Error**: `unbalanced braces in HCL content`

#### Invalid HCL with Line-Specific Error Reporting
```hcl
cli {
  default_timeout = 300
  max_parallel = 10
  invalid_syntax = 
}
```

**Error**: `invalid HCL syntax at line 4: invalid_syntax =`

## Usage Examples

### First-Time User Experience

```bash
# User runs their first spooky command
$ spooky project init my-project

# System automatically:
# 1. Detects OS (Linux)
# 2. Creates ~/.config/spooky/ directory
# 3. Creates spooky.hcl and logging.hcl with defaults
# 4. Proceeds with project initialization

✅ Project initialized successfully: /path/to/my-project
```

### Existing User Experience

```bash
# User with existing configuration runs a command
$ spooky project validate my-project

# System automatically:
# 1. Detects existing ~/.config/spooky/ directory
# 2. Validates spooky.hcl and logging.hcl syntax
# 3. If valid, proceeds with validation
# 4. If invalid, shows error and stops
```

### Invalid Configuration Handling

```bash
# User has corrupted configuration file
$ spooky project init test-project

Error: configuration setup failed: config validation failed: invalid HCL syntax in spooky.hcl: unbalanced braces in HCL content

# User has syntax error with line-specific reporting
$ spooky project init test-project

Error: configuration setup failed: config validation failed: invalid HCL syntax in spooky.hcl: invalid HCL syntax at line 4: invalid_syntax =
```

### Version and Help Commands

```bash
# These commands don't trigger auto-setup
$ spooky --version
0.20250812.0-dev-5eb15a4

$ spooky --help
spooky is a powerful automation and orchestration tool...
```

## Implementation Details

### Core Functions

#### `AutoSetupConfig()`
Main entry point that orchestrates the entire auto-setup process.

#### `getConfigDirectory()`
Determines the appropriate configuration directory based on OS.

#### `configDirectoryExists()`
Checks if the spooky configuration directory exists.

#### `createConfigDirectory()`
Creates the configuration directory and default files.

#### `ensureConfigFiles()`
Ensures required configuration files exist, creating them if missing.

#### `validateConfigFiles()`
Validates that existing configuration files have valid HCL syntax.

#### `validateHCLSyntax()`
Performs comprehensive HCL syntax validation with line-specific error reporting and multi-level validation.

### Integration with CLI

The auto-setup is integrated into the root command using `PersistentPreRunE`:

```go
RootCmd = &cobra.Command{
    // ... other configuration ...
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        // Skip auto-setup for version and help commands
        if cmd.Name() == "version" || cmd.Name() == "help" {
            return nil
        }
        
        // Auto-setup configuration for all other commands
        if err := spookyconfig.AutoSetupConfig(); err != nil {
            return fmt.Errorf("configuration setup failed: %w", err)
        }
        
        return nil
    },
}
```

## Benefits

### 1. **Zero Configuration Setup**
- Users don't need to manually create configuration files
- Sensible defaults are provided automatically
- No setup documentation required

### 2. **Cross-Platform Compatibility**
- Automatically adapts to different operating systems
- Follows platform conventions (XDG, macOS, Windows)
- Consistent experience across platforms

### 3. **Error Prevention**
- Validates configuration before use
- Prevents runtime errors from invalid configuration
- Clear error messages guide users to fix issues
- Line-specific error reporting for precise problem identification
- Multi-level validation ensures comprehensive error detection

### 4. **Backward Compatibility**
- Existing configurations continue to work
- No breaking changes for current users
- Gradual migration path

### 5. **Maintainability**
- Centralized configuration management
- Consistent validation across all commands
- Easy to extend with new configuration options
- Enhanced error reporting for easier debugging
- Robust validation system for reliable operation

## Testing

### Test Utilities

Several test utilities have been created to validate the auto-setup system:

1. **`examples/test-hcl-validation-utility.go`**
   - Tests HCL syntax validation
   - Validates both valid and invalid HCL content
   - Tests actual configuration files

2. **`examples/setup-logging-config-utility.go`**
   - Creates default logging configuration
   - Demonstrates manual setup process
   - Provides configuration examples

### Test Scenarios

#### Scenario 1: First-Time Setup
```bash
# Remove existing configuration
rm -rf ~/.config/spooky

# Run spooky command
./build/spooky project init test-project

# Verify configuration was created
ls -la ~/.config/spooky/
cat ~/.config/spooky/spooky.hcl
cat ~/.config/spooky/logging.hcl
```

#### Scenario 2: Invalid Configuration
```bash
# Create invalid configuration
echo 'invalid hcl content { missing closing brace' > ~/.config/spooky/spooky.hcl

# Run spooky command
./build/spooky project init test-project

# Should show validation error
```

#### Scenario 3: Version/Help Commands
```bash
# These should not trigger auto-setup
./build/spooky --version
./build/spooky --help
```

## Migration Guide

### For Existing Users

1. **No Action Required**: Existing configurations continue to work
2. **Automatic Validation**: Configuration files are validated on each run
3. **Error Reporting**: Invalid configurations are caught and reported

### For New Users

1. **Automatic Setup**: Configuration is created automatically on first use
2. **Sensible Defaults**: Default configuration provides good starting point
3. **Easy Customization**: Edit configuration files to customize behavior

## Implementation Status

### Current Implementation Features

The configuration system is **fully implemented** with the following capabilities:

1. **Basic Auto-Setup**: Creates spooky.hcl and logging.hcl files automatically
2. **OS Detection**: Supports Linux/BSD, macOS, and Windows
3. **HCL Validation**: Basic syntax validation with error reporting
4. **Error Handling**: Clear error messages for configuration issues
5. **Cross-Platform**: Works consistently across operating systems

### Future Enhancements

1. **Configuration Migration**: Automatic migration of old configuration formats
2. **Configuration Backup**: Automatic backup before making changes
3. **Configuration Templates**: Multiple configuration templates for different use cases
4. **Interactive Setup**: Guided configuration setup for advanced users
5. **Configuration Validation**: Schema-based validation beyond syntax checking

### Integration Opportunities

1. **Project-Specific Overrides**: Allow projects to override global configuration
2. **Environment-Specific Configs**: Different configurations for development/production
3. **Configuration Sharing**: Share configurations across team members
4. **Configuration Versioning**: Track configuration changes over time

## Conclusion

The auto-setup configuration system provides a seamless user experience while ensuring robust configuration management. It eliminates the need for manual setup while providing comprehensive validation and error handling. The system is designed to be maintainable, extensible, and compatible with existing workflows.

The current implementation focuses on the essential configuration files (spooky.hcl and logging.hcl) with basic validation and error handling. This provides a solid foundation for future enhancements while meeting the immediate needs of users.
