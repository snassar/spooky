# Spooky CLI Reference

## Overview

This document serves as the comprehensive reference for all spooky CLI commands, flags, and usage patterns. It is the single source of truth for command-line operations in the spooky automation and orchestration tool.

## Command Structure

All spooky commands follow a consistent **noun-verb** pattern:

```bash
spooky <noun> <verb> [arguments] [flags]
```

### Global Flags

These flags are available on all commands:

- `--version` - Show version information
- `--help` - Show help information for the command

## Command Categories

### 1. Project Management Commands

Commands for managing spooky projects, including initialization and validation.

#### `spooky project init`

Initialize a new spooky project with the required directory structure and configuration files.

**Syntax:**
```bash
spooky project init [project-name]
```

**Flags:**
- `--name string` - Project name (defaults to directory name)
- `--description string` - Project description
- `--version string` - Project version
- `--author string` - Project author
- `--email string` - Project email
- `--url string` - Project URL

**Examples:**
```bash
# Initialize a basic project
spooky project init my-automation-project

# Initialize with custom metadata
spooky project init my-project \
  --name "Production Automation" \
  --description "Automation for production environment" \
  --version "1.0.0" \
  --author "DevOps Team" \
  --email "devops@company.com" \
  --url "https://github.com/company/automation"
```

**What it does:**
- Creates project directory with required structure
- Generates default configuration files
- Validates project structure against schema
- Creates project.hcl with provided metadata

**Related Documentation:**
- [Project System User Guide](PROJECTS_DOCUMENTATION_SUMMARY.md)
- [Project System Design](../design/systems/project-system.md)

#### `spooky project validate`

Validate a spooky project structure and configuration.

**Syntax:**
```bash
spooky project validate [project-path]
```

**Examples:**
```bash
# Validate current directory as project
spooky project validate .

# Validate specific project
spooky project validate ./my-automation-project
```

**What it does:**
- Validates project directory structure
- Checks configuration file syntax
- Ensures compliance with project schema
- Reports validation errors and warnings

**Related Documentation:**
- [Project System User Guide](PROJECTS_DOCUMENTATION_SUMMARY.md)
- [Project System Design](../design/systems/project-system.md)

### 2. Machine Inventory Commands

Commands for managing machine inventory, including listing, validation, and connectivity testing.

#### `spooky machines list`

List all machines defined in the project's machine inventory.

**Syntax:**
```bash
spooky machines list [project-path]
```

**Examples:**
```bash
# List machines in current directory project
spooky machines list .

# List machines in specific project
spooky machines list ./my-automation-project
```

**What it does:**
- Reads machines.hcl files and machines/ directory
- Displays machine information (hostname, host, user)
- Shows connection status
- Groups machines by source file

**Related Documentation:**
- [Machines User Guide](MACHINES_USER_GUIDE.md)
- [Machines API Reference](MACHINES_API_REFERENCE.md)
- [Machines System Design](../design/systems/machines-system.md)

#### `spooky machines validate`

Validate machine inventory configuration and connectivity.

**Syntax:**
```bash
spooky machines validate [project-path]
```

**Examples:**
```bash
# Validate machines in current directory project
spooky machines validate .

# Validate machines in specific project
spooky machines validate ./my-automation-project
```

**What it does:**
- Validates machine configuration
- Checks required fields and authentication methods
- Validates SSH settings
- Reports configuration errors

**Related Documentation:**
- [Machines User Guide](MACHINES_USER_GUIDE.md)
- [Machines API Reference](MACHINES_API_REFERENCE.md)
- [Machines System Design](../design/systems/machines-system.md)

#### `spooky machines ping`

Test connectivity to machines in the inventory.

**Syntax:**
```bash
spooky machines ping [project-path]
```

**Flags:**
- `--format string` - Output format: text or json (default: text)
- `--verbose` - Show detailed output for all machines
- `--auth` - Test authentication in addition to connectivity

**Examples:**
```bash
# Basic connectivity test
spooky machines ping ./my-automation-project

# Verbose output with authentication testing
spooky machines ping ./my-automation-project --verbose --auth

# JSON output for scripting
spooky machines ping ./my-automation-project --format json
```

**What it does:**
- Tests network connectivity to all machines
- Checks SSH accessibility
- Reports response times and connection status
- Optionally tests authentication

**Related Documentation:**
- [Machines User Guide](MACHINES_USER_GUIDE.md)
- [Machines API Reference](MACHINES_API_REFERENCE.md)
- [Machines System Design](../design/systems/machines-system.md)

### 3. Variable Management Commands

Commands for managing project variables, including listing, validation, and resolution.

#### `spooky variables list`

List all variables defined in the project's variable files.

**Syntax:**
```bash
spooky variables list [project-path]
```

**Examples:**
```bash
# List variables in current directory project
spooky variables list .

# List variables in specific project
spooky variables list ./my-automation-project
```

**What it does:**
- Reads variables.hcl and variables/*.hcl files
- Displays variable information (name, type, description, scope)
- Shows variable metadata
- Groups variables by source file

**Related Documentation:**
- [Variables User Guide](VARIABLES_USER_GUIDE.md)
- [Variables API Reference](VARIABLES_API_REFERENCE.md)
- [Variables System Design](../design/systems/variables-system.md)

#### `spooky variables validate`

Validate variable definitions and dependencies.

**Syntax:**
```bash
spooky variables validate [project-path]
```

**Examples:**
```bash
# Validate variables in current directory project
spooky variables validate .

# Validate variables in specific project
spooky variables validate ./my-automation-project
```

**What it does:**
- Validates variable configuration
- Checks required fields and types
- Validates dependency relationships
- Reports validation errors

**Related Documentation:**
- [Variables User Guide](VARIABLES_USER_GUIDE.md)
- [Variables API Reference](VARIABLES_API_REFERENCE.md)
- [Variables System Design](../design/systems/variables-system.md)

#### `spooky variables resolve`

Resolve variables with the given context and display resolved values.

**Syntax:**
```bash
spooky variables resolve [project-path]
```

**Flags:**
- `--environment string` - Environment variables file (JSON)
- `--facts string` - Facts file (JSON)
- `--machines string` - Machine data file (JSON)
- `--user-data string` - User data file (JSON)
- `--json` - Output results in JSON format
- `--verbose` - Show detailed resolution information

**Examples:**
```bash
# Basic variable resolution
spooky variables resolve ./my-automation-project

# Resolution with external data sources
spooky variables resolve ./my-automation-project \
  --environment env.json \
  --facts facts.json \
  --machines machines.json

# JSON output for scripting
spooky variables resolve ./my-automation-project --json

# Verbose resolution with detailed information
spooky variables resolve ./my-automation-project --verbose
```

**What it does:**
- Loads variables from project
- Resolves using environment variables, facts, and machine data
- Displays resolved values
- Shows resolution context and errors

**Related Documentation:**
- [Variables User Guide](VARIABLES_USER_GUIDE.md)
- [Variables API Reference](VARIABLES_API_REFERENCE.md)
- [Variables System Design](../design/systems/variables-system.md)

### 4. Facts Management Commands

Commands for managing machine facts, including collection and export.

#### `spooky facts export`

Export facts from machines to files in various formats.

**Syntax:**
```bash
spooky facts export [project-path]
```

**Flags:**
- `--format string` - Export format (hcl, json) (default: hcl)
- `--output string` - Output file path (required)
- `--machine string` - Filter to specific machine
- `--tags stringSlice` - Filter by tags (supports key=value or key-only)
- `--groups stringSlice` - Filter by groups
- `--parallel int` - Number of parallel workers (default: 1)
- `--verbose` - Verbose output

**Examples:**
```bash
# Export all facts to HCL format
spooky facts export ./my-automation-project --output facts.hcl

# Export to JSON format
spooky facts export ./my-automation-project --format json --output facts.json

# Export with parallel processing
spooky facts export ./my-automation-project --parallel 4 --output facts.hcl

# Export from specific machine
spooky facts export ./my-automation-project --machine web-server --output web-server-facts.hcl

# Export with tag filtering
spooky facts export ./my-automation-project --tags environment=production --output prod-facts.hcl

# Export with group filtering
spooky facts export ./my-automation-project --groups webservers,database --output app-facts.hcl

# Export with multiple filters
spooky facts export ./my-automation-project \
  --tags role=web \
  --groups production \
  --output web-prod-facts.hcl

# Verbose export with detailed progress
spooky facts export ./my-automation-project --verbose --output facts.hcl
```

**What it does:**
- Automatically gathers facts from all machines in the project inventory
- Exports them to the specified format for backup, analysis, or transfer
- Supports filtering by machine, tags, and groups
- Provides parallel processing for large inventories

**Related Documentation:**
- [Facts User Guide](FACTS_USER_GUIDE.md)
- [Facts API Reference](FACTS_API_REFERENCE.md)
- [Facts System Design](../plans/facts-system-design.md)

## Common Usage Patterns

### Project Workflow

A typical spooky project workflow involves these commands:

```bash
# 1. Initialize a new project
spooky project init my-automation-project

# 2. Add machine inventory
# (Edit machines.hcl file)

# 3. Validate project structure
spooky project validate my-automation-project

# 4. Validate machine inventory
spooky machines validate my-automation-project

# 5. Test machine connectivity
spooky machines ping my-automation-project

# 6. Add variables
# (Edit variables.hcl file)

# 7. Validate variables
spooky variables validate my-automation-project

# 8. Resolve variables
spooky variables resolve my-automation-project

# 9. Export facts
spooky facts export my-automation-project --output facts.hcl
```

### Filtering and Selection

Many commands support filtering to work with subsets of resources:

```bash
# Filter machines by tags
spooky machines ping ./project --tags environment=production

# Filter facts export by machine and tags
spooky facts export ./project \
  --machine web-server \
  --tags role=web \
  --output web-facts.hcl

# Filter variables resolution with external data
spooky variables resolve ./project \
  --environment prod-env.json \
  --facts latest-facts.json
```

### Output Formats

Commands that support multiple output formats:

```bash
# Text output (default)
spooky machines ping ./project

# JSON output for scripting
spooky machines ping ./project --format json

# JSON output for variables
spooky variables resolve ./project --json
```

## Error Handling

All commands provide consistent error handling:

- **Validation errors** are clearly reported with file and line information
- **Connection errors** include detailed diagnostic information
- **Configuration errors** provide actionable suggestions
- **Permission errors** include guidance on required access levels

## Configuration

The spooky CLI automatically sets up configuration on first run:

- **Linux/BSD**: `$XDG_CONFIG_HOME/spooky/` (defaults to `~/.config/spooky/`)
- **macOS**: `~/Library/Application Support/spooky/`
- **Windows**: `%APPDATA%\spooky\`

Configuration files:
- `spooky.hcl` - Main configuration file
- `logging.hcl` - Logging configuration

## Shell Completion

Spooky supports shell completion for better user experience:

```bash
# Generate bash completion
spooky completion bash > ~/.local/share/bash-completion/completions/spooky

# Generate zsh completion
spooky completion zsh > ~/.zsh/completions/_spooky

# Generate fish completion
spooky completion fish > ~/.config/fish/completions/spooky.fish
```

## Troubleshooting

### Common Issues

1. **Project not found**: Ensure you're in the correct directory or provide the full path
2. **Permission denied**: Check file permissions and SSH key access
3. **Connection timeout**: Verify network connectivity and firewall settings
4. **Validation errors**: Check HCL syntax and required fields

### Debug Mode

Use the `--verbose` flag for detailed output:

```bash
spooky machines ping ./project --verbose
spooky facts export ./project --verbose --output facts.hcl
spooky variables resolve ./project --verbose
```

### Getting Help

```bash
# General help
spooky --help

# Command-specific help
spooky project --help
spooky machines --help
spooky variables --help
spooky facts --help

# Subcommand help
spooky project init --help
spooky machines ping --help
spooky variables resolve --help
spooky facts export --help
```

## Related Documentation

- [Development Guide](DEVELOPMENT.md) - Development setup and guidelines
- [Auto-Setup Configuration](AUTO_SETUP_CONFIGURATION.md) - Configuration management
- [CLI System Design](../design/systems/cli-system.md) - CLI architecture and design
- [Configuration System](../design/systems/configuration-system.md) - Configuration management
- [Schema System](../design/systems/schema-system.md) - Schema validation and management

## Version Information

```bash
spooky --version
```

This command displays the current version of spooky, including build information and ScalVer versioning details.
