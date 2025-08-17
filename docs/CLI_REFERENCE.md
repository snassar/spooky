# CLI Reference

## Overview

This document provides a comprehensive reference for the spooky command-line interface (CLI). It covers all commands, options, and usage patterns for the spooky tool.

**Status: Implemented** - The CLI system has comprehensive functionality with all major commands implemented and working.

## Command Structure

The spooky CLI follows a consistent command structure:

```bash
spooky [global-options] <command> [command-options] <arguments>
```

### Global Options

```bash
--config <file>           # Specify configuration file path
--log-level <level>       # Set logging level (debug, info, warn, error)
--log-format <format>     # Set logging format (json, text)
--quiet                   # Suppress output
--verbose                 # Enable verbose output
--version                 # Show version information
--help                    # Show help information
```

## Basic Commands

### Version and Help

```bash
# Check version
$ spooky --version
0.20250817.0-dev-d3c8409

# Get help
$ spooky --help
spooky is a powerful automation and orchestration tool built with Go.

It provides declarative configuration, parallel run capabilities, 
and intelligent fact-driven decision making for heterogeneous environments.

Examples:
  spooky project init my-project
  spooky project validate my-project
  spooky --version
  spooky --help

Usage:
  spooky [command]

Available Commands:
  actions      Manage and run actions
  completion   Generate the autocompletion script for the specified shell
  facts        Export machine facts
  help         Help about any command
  integrations Manage system integrations
  machines     Manage machine inventory
  project      Manage spooky projects
  schemas      Manage and validate schemas
  secrets      Manage age encryption and secrets
  templates    Manage templates
  variables    Manage project variables

Flags:
  -h, --help      help for spooky
  -v, --version   version for spooky

Use "spooky [command] --help" for more information about a command.
```

### Command Help Examples

```bash
# Actions help
$ spooky actions --help
Manage and run actions on machines.

Actions are run in dependency order. If no action names are provided,
all actions in the project will be run.

Usage:
  spooky actions [command]

Available Commands:
  list        List available actions
  run         Run actions on machines
  validate    Validate action configurations

Flags:
  -h, --help   help for actions

Use "spooky actions [command] --help" for more information about a command.

# Facts help
$ spooky facts --help
Export machine facts to files in various formats.

Facts are system information collected from machines including OS details,
hardware information, network configuration, and custom data. Facts can be
exported to JSON or HCL format for backup, analysis, or transfer to other systems.

Usage:
  spooky facts [command]

Available Commands:
  export      Export facts to file

Flags:
  -h, --help   help for facts

Use "spooky facts [command] --help" for more information about a command.

# Machines help
$ spooky machines --help
Manage machine inventory including listing, validation, and connectivity testing.

Machine inventory is defined in machines.hcl files within spooky projects and contains
SSH connection details, authentication information, and machine metadata.

Usage:
  spooky machines [command]

Available Commands:
  encrypt     Encrypt machines in a project
  export      Export machines to HCL format
  list        List machines in a project
  ping        Ping machines to test connectivity
  validate    Validate machine inventory

Flags:
  -h, --help   help for machines

Use "spooky machines [command] --help" for more information about a command.
```

## Core Commands

### Project Commands

#### `spooky project init <project-name>`

Initialize a new spooky project with the specified name.

```bash
# Initialize a new project
$ spooky project init my-project

# Initialize with specific metadata
$ spooky project init my-project --name "Test Project" --description "A test project for documentation examples"

✅ Project initialized successfully: /home/sn/Workshop/go/spooky/my-project
📁 Project structure created according to project-directory.schema.hcl
📄 Configuration files generated using project.schema.hcl
💡 Next steps:
   - Edit project.hcl to customize your project
   - Add machines.hcl for machine inventory
   - Add actions.hcl for automation tasks
   - Add variables.hcl for project variables
```

**Options:**
- `--name <string>` - Project name (defaults to directory name)
- `--description <string>` - Project description
- `--version <string>` - Project version
- `--author <string>` - Project author
- `--email <string>` - Project email
- `--url <string>` - Project URL

#### `spooky project validate <project-path>`

Validate a spooky project structure and configuration.

```bash
# Validate project
$ spooky project validate ./my-project

🔍 Validating project: /path/to/my-project

❌ Validation issues found:
   - Project version must be in ScalVer format (MAJOR.DATE.PATCH)
⚠️  Warnings:
   - Optional file not found: machines.hcl
   - Optional file not found: actions.hcl
   - Optional file not found: variables.hcl
   - Optional directory not found: machines
   - Optional directory not found: actions
   - Optional directory not found: variables
   - Optional directory not found: templates
```

This command validates that the project follows the project-directory.schema.hcl schema and that all configuration files are properly formatted.

#### `spooky project encrypt <project-path>`

Encrypt all variables and machines in a project that have encrypted=true.

```bash
# Encrypt project
spooky project encrypt ./my-project

# Show what would be encrypted without making changes
spooky project encrypt ./my-project --dry-run
```

**Options:**
- `--dry-run` - Show what would be encrypted without making changes

### Actions Commands

#### `spooky actions list <project-path>`

List available actions in a project.

```bash
# List all actions
$ spooky actions list ./my-project

No actions found in project
```

#### `spooky actions run <project-path>`

Run actions in a spooky project.

```bash
# Run all actions
spooky actions run ./my-project

# Run with dry-run mode
spooky actions run ./my-project --dry-run

# Run with plan mode
spooky actions run ./my-project --plan

# Run with parallel execution
spooky actions run ./my-project --parallel 4

# Run with decryption
spooky actions run ./my-project --decrypt
```

**Options:**
- `--machine <list>` - Target specific machines
- `--tags <list>` - Target machines by tags
- `--filter <query>` - Use complex filter query
- `--parallel <number>` - Number of parallel workers (minimum 2)
- `--dry-run` - Simulate running without making changes
- `--plan` - Show running plan without running
- `--decrypt` - Decrypt encrypted variables and facts in-memory for debugging

#### `spooky actions validate <project-path>`

Validate actions in a project.

```bash
# Validate all actions
spooky actions validate ./my-project
```

### Facts Commands

#### `spooky facts export <project-path>`

Export facts to file in various formats.

```bash
# Export facts to HCL format
spooky facts export ./my-project --output facts.hcl

# Export facts to JSON format
spooky facts export ./my-project --format json --output facts.json

# Export facts for specific machines
spooky facts export ./my-project --machine web-server --output web-server-facts.hcl

# Export facts with parallel collection
spooky facts export ./my-project --parallel 4 --output facts.hcl

# Export facts with filtering
spooky facts export ./my-project --tags environment=production --output prod-facts.hcl
```

**Options:**
- `--format <string>` - Export format (hcl, json) [default: hcl]
- `--output <string>` - Output file path (required)
- `--machine <string>` - Filter to specific machine
- `--tags <list>` - Filter by tags (supports key=value or key-only)
- `--groups <list>` - Filter by groups
- `--parallel <number>` - Number of parallel workers [default: 1]
- `--verbose` - Verbose output

### Machines Commands

#### `spooky machines list <project-path>`

List machines in a project.

```bash
# List all machines
spooky machines list ./my-project
```

#### `spooky machines validate <project-path>`

Validate machine inventory configuration and connectivity.

```bash
# Validate machines
spooky machines validate ./my-project
```

#### `spooky machines ping <project-path>`

Ping machines to test network connectivity and SSH accessibility.

```bash
# Ping all machines
spooky machines ping ./my-project

# Ping with authentication testing
spooky machines ping ./my-project --auth

# Ping with verbose output
spooky machines ping ./my-project --verbose

# Ping with JSON output
spooky machines ping ./my-project --format json
```

**Options:**
- `--format <string>` - Output format: text or json [default: text]
- `--verbose` - Show detailed output for all machines
- `--auth` - Test authentication in addition to connectivity

#### `spooky machines export <project-path>`

Export machine inventory to HCL format.

```bash
# Export all machines
spooky machines export ./my-project --output machines.hcl

# Export specific machine
spooky machines export ./my-project --machine web-server --output web-server.hcl

# Export machines by tags
spooky machines export ./my-project --tags environment=production --output prod-machines.hcl
```

**Options:**
- `--output <string>` - Output file path (required)
- `--machine <string>` - Export specific machine by hostname
- `--tags <list>` - Filter machines by tags (key=value or key-only)

#### `spooky machines encrypt <project-path>`

Encrypt machines in a project using age encryption.

```bash
# Encrypt machines
spooky machines encrypt ./my-project

# Show what would be encrypted without making changes
spooky machines encrypt ./my-project --dry-run
```

**Options:**
- `--dry-run` - Show what would be encrypted without making changes

### Variables Commands

#### `spooky variables list <project-path>`

List variables in a project.

```bash
# List all variables
spooky variables list ./my-project
```

#### `spooky variables validate <project-path>`

Validate variable definitions and dependencies.

```bash
# Validate all variables
spooky variables validate ./my-project
```

#### `spooky variables resolve <project-path>`

Resolve variables with the given context and display resolved values.

```bash
# Resolve variables
spooky variables resolve ./my-project

# Resolve with environment variables
spooky variables resolve ./my-project --environment env.json

# Resolve with facts data
spooky variables resolve ./my-project --facts facts.json

# Resolve with JSON output
spooky variables resolve ./my-project --json
```

**Options:**
- `--environment <string>` - Environment variables file (JSON)
- `--facts <string>` - Facts file (JSON)
- `--machines <string>` - Machine data file (JSON)
- `--user-data <string>` - User data file (JSON)
- `--json` - Output results in JSON format
- `--verbose` - Show detailed resolution information

#### `spooky variables encrypt <project-path>`

Encrypt variables in a project.

```bash
# Encrypt variables
spooky variables encrypt ./my-project

# Show what would be encrypted without making changes
spooky variables encrypt ./my-project --dry-run
```

**Options:**
- `--dry-run` - Show what would be encrypted without making changes

#### `spooky variables armor <project-path>`

Armor (encrypt) variables in a project using age encryption.

```bash
# Armor variables
spooky variables armor ./my-project

# Show what would be encrypted without making changes
spooky variables armor ./my-project --dry-run
```

**Options:**
- `--dry-run` - Show what would be encrypted without making changes

### Templates Commands

#### `spooky templates render <project-path> <template-path>`

Render a template with the given data and output to a file.

```bash
# Basic template rendering
spooky templates render ./my-project templates/nginx.conf.tmpl

# Render with output file
spooky templates render ./my-project templates/nginx.conf.tmpl --output /etc/nginx/nginx.conf

# Preview mode
spooky templates render ./my-project templates/nginx.conf.tmpl --preview

# Dry run mode
spooky templates render ./my-project templates/nginx.conf.tmpl --dry-run
```

**Options:**
- `--data <string>` - Data file (JSON/HCL) for template variables
- `--output <string>` - Output file path
- `--dry-run` - Show what would be rendered without writing files
- `--preview` - Preview the rendered template

#### `spooky templates validate <project-path>`

Validate templates in the project for syntax and security.

```bash
# Validate all templates
spooky templates validate ./my-project

# Validate specific template
spooky templates validate ./my-project --template templates/nginx.conf.tmpl
```

**Options:**
- `--template <string>` - Specific template to validate

#### `spooky templates list <project-path>`

List all templates in the project.

```bash
# List all templates
spooky templates list ./my-project

# List with specific format
spooky templates list ./my-project --format json
spooky templates list ./my-project --format hcl
```

**Options:**
- `--format <string>` - Output format (table, json, hcl) [default: table]

#### `spooky templates search <project-path> <query>`

Search templates by query.

```bash
# Search templates by query
spooky templates search ./my-project "nginx"

# Search with tags
spooky templates search ./my-project "config" --tags web

# Search by category
spooky templates search ./my-project "deploy" --category deployment
```

**Options:**
- `--tags <list>` - Filter by tags
- `--category <string>` - Filter by category

### Schemas Commands

#### `spooky schemas validate <schema-file> <data-file>`

Validate HCL data files against their corresponding schemas.

```bash
# Validate data against schema
spooky schemas validate project.schema.hcl project.hcl
```

#### `spooky schemas list`

List all available schemas in the system.

```bash
# List schemas
spooky schemas list
```

### Secrets Commands

#### `spooky secrets validate <project-path>`

Validate age configuration and keys for a project.

```bash
# Validate secrets configuration
spooky secrets validate ./my-project
```

This command validates:
- Age configuration in spooky.hcl
- Identity files and permissions
- Recipients file format
- Encrypted values in project files

### Integrations Commands

#### `spooky integrations list`

List all available integrations and their current status.

```bash
# List integrations
spooky integrations list
```

#### `spooky integrations validate`

Validate that all integrations are working correctly.

```bash
# Validate integrations
spooky integrations validate
```

## Configuration

### Global Configuration

Spooky uses a global configuration file located at `$XDG_CONFIG_HOME/spooky/spooky.hcl` (defaulting to `$HOME/.config/spooky/spooky.hcl`).

### Project Configuration

Each spooky project contains configuration files that define the project structure and behavior:

- `project.hcl` - Project metadata and settings
- `machines.hcl` - Machine inventory and SSH configuration
- `actions.hcl` - Action definitions and templates
- `variables.hcl` - Project variables and configuration
- `variables/*.hcl` - Additional variable files

### Age Encryption Configuration

For encryption support, spooky uses age encryption with the following configuration:

- Identity files: `~/.config/spooky/identities/`
- Recipients file: `~/.config/spooky/recipients.txt`

## Examples

### Project Management

```bash
# Initialize a new project
spooky project init my-automation

# Validate project structure
spooky project validate my-automation

# Show project information
spooky project info my-automation

# Check project status
spooky project status my-automation
```

### Action Management

```bash
# List actions in project
spooky actions list my-automation

# Validate action definitions
spooky actions validate my-automation

# Run actions on all machines
spooky actions run my-automation

# Run actions on specific machines
spooky actions run my-automation --machine web-server

# Run actions with parallel execution
spooky actions run my-automation --parallel 4

# Run actions with dry-run mode
spooky actions run my-automation --dry-run

# Run actions with decryption
spooky actions run my-automation --decrypt
```

### Basic Project Workflow

```bash
# Initialize a new project
spooky project init my-automation

# Add machines to inventory
# Edit machines.hcl file

# Validate project structure
spooky project validate my-automation

# Test machine connectivity
spooky machines ping my-automation

# Add actions
# Edit actions.hcl file

# Validate actions
spooky actions validate my-automation

# Run actions
spooky actions run my-automation
```

### Fact Collection and Export

```bash
# Export facts from all machines
spooky facts export my-automation --output facts.hcl

# Export facts from specific machines
spooky facts export my-automation --machine web-server --output web-facts.hcl

# Export facts with parallel collection
spooky facts export my-automation --parallel 4 --output facts.hcl
```

### Variable Management

```bash
# List project variables
spooky variables list my-automation

# Validate variable definitions
spooky variables validate my-automation

# Resolve variables with context
spooky variables resolve my-automation

# Encrypt sensitive variables
spooky variables armor my-automation
```

### Machine Management

```bash
# List machines in project
spooky machines list my-automation

# Test machine connectivity
spooky machines ping my-automation

# Test authentication
spooky machines ping my-automation --auth

# Export machine inventory
spooky machines export my-automation --output inventory.hcl
```

### Template Management

```bash
# List templates in project
spooky templates list my-automation

# Validate templates
spooky templates validate my-automation

# Render template
spooky templates render my-automation templates/nginx.conf.tmpl --output nginx.conf

# Preview template rendering
spooky templates render my-automation templates/deploy.sh.tmpl --preview
```

### Schema Validation

```bash
# Validate project against schema
spooky schemas validate project.schema.hcl project.hcl

# List available schemas
spooky schemas list
```

### Secrets Management

```bash
# Validate age configuration
spooky secrets validate my-automation

# Encrypt project variables
spooky variables armor my-automation

# Run actions with decryption
spooky actions run my-automation --decrypt
```

### Integration Management

```bash
# List available integrations
spooky integrations list

# Validate all integrations
spooky integrations validate
```

## Error Handling

The spooky CLI provides comprehensive error reporting with:

- Clear error messages with context
- Detailed validation results
- Suggestions for fixing common issues
- Exit codes for programmatic use

## Exit Codes

- `0` - Success
- `1` - General error
- `2` - Configuration error
- `3` - Validation error
- `4` - Connection error
- `5` - Authentication error
