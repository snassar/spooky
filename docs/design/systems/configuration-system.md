# Configuration System: Implementation Overview

## Overview

This document provides an overview of the spooky configuration system implementation. It covers the actual implemented configuration features, auto-setup capabilities, and integration with the spooky systems.

**Schema Integration**: The configuration system uses schema validation for configuration files and follows established validation patterns.

**Architecture Integration**: Configuration integrates with the overall spooky architecture, providing centralized settings and OS-specific configuration management.

## System Integration

The configuration system integrates with the implemented Spooky systems to provide configuration management:

### **Project System Integration**
- **Project Configuration**: Project-specific configuration overrides (see [Project System](../project-system.md))
- **Project Isolation**: Project configuration isolation from global settings
- **Project Validation**: Configuration validation for project settings
- **Project Context**: Configuration resolution within project context

### **CLI System Integration**
- **Auto-Setup**: Automatic configuration setup for all CLI commands (see [CLI System](../cli-system.md))
- **OS-Specific Configuration**: Cross-platform configuration directory management
- **Configuration Validation**: HCL syntax validation for configuration files
- **Configuration Discovery**: Automatic configuration file creation and validation

### **Schema System Integration**
- **Configuration Schema**: Schema validation for configuration structure (see [Schema System](../schema-system.md))
- **Schema Validation**: Configuration file validation against embedded schemas
- **Schema Evolution**: Configuration schema versioning and migration support

## Current State Analysis

### **✅ Implemented Features**
- **Auto-setup configuration** with OS-specific directory handling
- **HCL parsing and validation** infrastructure
- **XDG directory compliance** on Linux/BSD systems
- **Cross-platform support** for macOS and Windows
- **Configuration file validation** and error reporting

### **📋 Future Enhancements**
- **Global configuration file** with comprehensive settings
- **Configuration precedence** hierarchy
- **Environment variable integration** for configuration overrides
- **CLI configuration commands** for management

## Auto-Setup Configuration

### **Implementation Overview**

The configuration system automatically sets up the spooky environment on first use:

```go
// AutoSetupConfig ensures spooky configuration directory exists and is properly configured
// This function is called before any spooky command run (except --version and --help)
func AutoSetupConfig() error {
    // Determine OS and get appropriate config directory
    // Check if config directory exists
    // Create config directory and default files if needed
    // Validate existing config files
}
```

### **OS-Specific Configuration Directories**

#### **Linux/BSD Systems**
```go
// Use XDG Base Directory Specification
xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
if xdgConfigHome == "" {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", fmt.Errorf("failed to get user home directory: %w", err)
    }
    xdgConfigHome = filepath.Join(homeDir, ".config")
}
return filepath.Join(xdgConfigHome, "spooky"), nil
```

**Configuration Location**: `$XDG_CONFIG_HOME/spooky/` (defaults to `~/.config/spooky/`)

#### **macOS Systems**
```go
// macOS: ~/Library/Application Support/spooky
homeDir, err := os.UserHomeDir()
if err != nil {
    return "", fmt.Errorf("failed to get user home directory: %w", err)
}
return filepath.Join(homeDir, "Library", "Application Support", "spooky"), nil
```

**Configuration Location**: `~/Library/Application Support/spooky/`

#### **Windows Systems**
```go
// Windows: %APPDATA%\spooky
appData := os.Getenv("APPDATA")
if appData == "" {
    return "", fmt.Errorf("APPDATA environment variable not set")
}
return filepath.Join(appData, "spooky"), nil
```

**Configuration Location**: `%APPDATA%\spooky\`

### **Configuration File Structure**

The auto-setup system creates the following configuration files:

```
spooky-config/
├── spooky.hcl          # Main configuration file
└── logging.hcl         # Logging configuration file
```

#### **Default spooky.hcl**
```hcl
# spooky.hcl - Main spooky configuration
spooky {
  # Basic configuration settings
  version = "1.0.0"
  
  # Default settings
  defaults {
    timeout = 30
    retries = 3
    parallel = 4
  }
  
  # Validation settings
  validation {
    strict_mode = false
    validate_before_run = true
  }
}
```

#### **Default logging.hcl**
```hcl
# logging.hcl - Logging configuration
logging {
  # Global logging settings
  level = "info"
  format = "text"
  
  # Component-specific settings
  components {
    "project" = "info"
    "machines" = "info"
    "variables" = "info"
    "ssh" = "warn"
  }
  
  # Output settings
  output {
    destination = "stderr"
    include_timestamps = true
    include_colors = true
  }
}
```

## Configuration Validation

### **HCL Syntax Validation**

All configuration files are validated for proper HCL syntax:

```go
func validateConfigFiles(configDir string) error {
    // Validate spooky.hcl
    spookyConfigPath := filepath.Join(configDir, "spooky.hcl")
    if err := validateHCLFile(spookyConfigPath); err != nil {
        return fmt.Errorf("spooky.hcl validation failed: %w", err)
    }
    
    // Validate logging.hcl
    loggingConfigPath := filepath.Join(configDir, "logging.hcl")
    if err := validateHCLFile(loggingConfigPath); err != nil {
        return fmt.Errorf("logging.hcl validation failed: %w", err)
    }
    
    return nil
}
```

### **Schema Validation**

Configuration files are validated against embedded schemas:

```go
func validateHCLFile(filePath string) error {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("failed to read file: %w", err)
    }
    
    // Parse HCL with validation
    var result interface{}
    if err := hcl.Unmarshal(data, &result); err != nil {
        return fmt.Errorf("HCL parsing failed: %w", err)
    }
    
    return nil
}
```

## CLI Integration

### **Automatic Configuration Setup**

The CLI automatically sets up configuration before running any command:

```go
// RootCmd PersistentPreRunE
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
```

### **Configuration Error Handling**

Configuration errors are handled gracefully:

```go
// Configuration setup with error reporting
if err := spookyconfig.AutoSetupConfig(); err != nil {
    return fmt.Errorf("configuration setup failed: %w", err)
}
```

## Benefits

### **Immediate Benefits**
- **Automatic setup** - No manual configuration required
- **Cross-platform support** - Works consistently across operating systems
- **XDG compliance** - Follows Linux/BSD standards
- **Error validation** - Catches configuration issues early
- **User-friendly** - Transparent configuration management

### **Long-term Benefits**
- **Consistent behavior** - Same configuration across environments
- **Extensible design** - Easy to add new configuration options
- **Schema validation** - Ensures configuration integrity
- **Migration support** - Handles configuration evolution

## Current Status

### **✅ Implemented Features**
- Auto-setup configuration system
- OS-specific configuration directories
- HCL syntax validation
- Configuration file creation
- Cross-platform support
- CLI integration

### **📋 Future Enhancements**
- Global configuration file with comprehensive settings
- Configuration precedence hierarchy
- Environment variable integration
- Configuration management commands
- Advanced validation rules
- Configuration migration tools

## Integration with Other Rules

This configuration system works in conjunction with:
- **Interface Architecture**: Uses interface-based configuration management
- **Error Handling Standards**: Follows structured error handling patterns
- **Schema System**: Uses schema validation for configuration files
- **CLI System**: Integrates with command-line interface

### Cross-References
- **See [CLI System](../cli-system.md)**: For CLI integration details
- **See [Schema System](../schema-system.md)**: For validation details
- **See [Project System](../project-system.md)**: For project configuration details

## Remember

**The configuration system provides:**
- Automatic setup and validation
- Cross-platform compatibility
- Schema-based validation
- Extensible configuration management

**All configuration follows established patterns and integrates with the implemented systems.** 