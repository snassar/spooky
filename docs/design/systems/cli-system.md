# CLI System: Implementation Overview

## Overview

This document provides an overview of the spooky CLI system implementation. It covers the actual implemented command patterns, user interface design, and integration with the spooky systems.

**Schema Integration**: The CLI system uses schema validation for configuration files and project structure validation.

**Architecture Integration**: CLI integrates with the overall spooky architecture, providing the primary user interface for system operations and management.

## System Integration

The CLI system integrates with the implemented Spooky systems to provide command-line management:

### **Project System Integration**
- **Project Commands**: `spooky project` commands for project management (see [Project System](../project-system.md))
- **Project Initialization**: `spooky project init` for project creation
- **Project Validation**: `spooky project validate` for project structure validation

### **Machines System Integration**
- **Machine Commands**: Machine management through `spooky machines` commands (see [Machines System](../machines-system.md))
- **Command Patterns**: Machine commands follow the established `spooky noun verb` CLI pattern
- **Validation Commands**: `spooky machines validate` for inventory validation
- **Management Commands**: `spooky machines list`, `spooky machines ping`

### **Variables System Integration**
- **Variables Commands**: `spooky variables` commands for variable management (see [Variables System](../variables-system.md))
- **Variables Validation**: `spooky variables validate` for variable file validation
- **Variables Discovery**: `spooky variables list` for variable discovery and filtering
- **Variables Resolution**: `spooky variables resolve` for variable resolution with context

### **Configuration System Integration**
- **Auto-Setup**: Automatic configuration directory creation and validation
- **OS-Specific Configuration**: XDG compliance on Linux/BSD, macOS, and Windows support
- **Configuration Validation**: HCL syntax validation for configuration files

## Implemented CLI Design

```
# Core project and system management
spooky project {init, validate} <project directory>
spooky machines {list, validate, ping} <project directory>
spooky variables {list, validate, resolve} <project directory>

# Global flags
spooky --version                   # show version information
spooky --help                      # show help
```

## Implementation Details

### **Command Structure**

The CLI follows a consistent noun-verb pattern:

```bash
spooky <noun> <verb> [arguments]
```

**Available Nouns:**
- `project` - Project management and validation
- `machines` - Machine inventory operations
- `variables` - Variable management and resolution

**Available Verbs:**
- `init` - Initialize new resources (project)
- `validate` - Validate configuration and structure
- `list` - List resources (machines, variables)
- `ping` - Test connectivity (machines)
- `resolve` - Resolve with context (variables)

### **Project Commands**

#### `spooky project init`
```bash
spooky project init [project-name]
```
Creates a new spooky project with the required directory structure and configuration files.

**Features:**
- Creates project directory with required structure
- Generates default configuration files
- Supports customization via flags (name, description, version, author, email, url)
- Validates project structure against schema

#### `spooky project validate`
```bash
spooky project validate [project-path]
```
Validates a spooky project structure and configuration.

**Features:**
- Validates project directory structure
- Checks configuration file syntax
- Ensures compliance with project schema
- Reports validation errors and warnings

### **Machines Commands**

#### `spooky machines list`
```bash
spooky machines list [project-path]
```
Lists all machines defined in the project's machine inventory.

**Features:**
- Reads machines.hcl files
- Displays machine information (hostname, host, user)
- Shows connection status
- Supports verbose output

#### `spooky machines validate`
```bash
spooky machines validate [project-path]
```
Validates machine inventory configuration and connectivity.

**Features:**
- Validates machine configuration
- Checks required fields and authentication methods
- Validates SSH settings
- Reports configuration errors

#### `spooky machines ping`
```bash
spooky machines ping [project-path]
```
Tests connectivity to machines in the inventory.

**Features:**
- Tests network connectivity
- Checks SSH accessibility
- Reports response times and connection status
- Supports verbose output for detailed results

### **Variables Commands**

#### `spooky variables list`
```bash
spooky variables list [project-path]
```
Lists all variables defined in the project's variable files.

**Features:**
- Reads variables.hcl and variables/*.hcl files
- Displays variable information (name, type, description, scope)
- Shows variable metadata
- Supports filtering and sorting

#### `spooky variables validate`
```bash
spooky variables validate [project-path]
```
Validates variable definitions and dependencies.

**Features:**
- Validates variable configuration
- Checks required fields and types
- Validates dependency relationships
- Reports validation errors

#### `spooky variables resolve`
```bash
spooky variables resolve [project-path]
```
Resolves variables with the given context and displays resolved values.

**Features:**
- Loads variables from project
- Resolves using environment variables, facts, and machine data
- Displays resolved values
- Shows resolution context

## Configuration Integration

### **Auto-Setup Configuration**

The CLI automatically sets up configuration on first run:

```go
// Auto-setup runs before any command (except --version and --help)
func AutoSetupConfig() error {
    // Determine OS-specific config directory
    // Create config directory if needed
    // Generate default configuration files
    // Validate existing configuration
}
```

### **OS-Specific Configuration**

- **Linux/BSD**: Follows XDG Base Directory Specification (`$XDG_CONFIG_HOME/spooky/`)
- **macOS**: Uses `~/Library/Application Support/spooky/`
- **Windows**: Uses `%APPDATA%\spooky\`

### **Configuration Files**

- `spooky.hcl` - Main configuration file
- `logging.hcl` - Logging configuration

## Error Handling

### **Command Error Handling**
```go
// All commands use consistent error handling
func handleCommand(cmd *cobra.Command, args []string) error {
    // Initialize dependencies
    // Execute command logic
    // Return structured errors
}
```

### **Validation Error Reporting**
- Clear error messages with file and line information
- Structured validation results
- Actionable error suggestions

## Benefits of the Current Design

1. **Consistent Patterns**: All commands follow the same noun-verb structure
2. **Clear Separation**: Each noun represents a distinct domain
3. **Extensible**: New verbs can be added to existing nouns
4. **User-Friendly**: Intuitive command structure
5. **Validation-First**: All commands include validation capabilities
6. **Cross-Platform**: Works consistently across operating systems

## Current Status

### **✅ Implemented Features**
- Project initialization and validation
- Machine inventory management and connectivity testing
- Variable management and resolution
- Auto-setup configuration system
- Cross-platform configuration support
- Schema-based validation

### **📋 Future Enhancements**
- Additional command nouns (facts, actions, templates)
- Enhanced filtering and search capabilities
- Shell completion support
- Configuration management commands
- Export and import functionality

## Integration with Other Rules

This CLI system works in conjunction with:
- **Interface Architecture**: Uses interface-based dependency injection
- **Error Handling Standards**: Follows structured error handling patterns
- **Configuration Management**: Integrates with auto-setup configuration
- **Schema System**: Uses schema validation for all operations

### Cross-References
- **See [Project System](../project-system.md)**: For project management details
- **See [Machines System](../machines-system.md)**: For machine inventory details
- **See [Variables System](../variables-system.md)**: For variable management details
- **See [Schema System](../schema-system.md)**: For validation details

## Remember

**The CLI system provides:**
- Consistent and intuitive command structure
- Comprehensive validation and error reporting
- Cross-platform configuration management
- Extensible architecture for future features

**All commands follow established patterns and integrate with the implemented systems.**