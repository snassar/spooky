# Spooky CLI Specification

## Overview
This document defines the command-line interface for spooky, a server configuration and automation tool. The CLI follows hierarchical command patterns using the Cobra library, providing a consistent and scalable interface.

## Command Structure

### Root Command
```bash
spooky [global-flags] <command> [subcommand] [flags] [arguments]
```

### Global Flags
```bash
--config-dir string     Directory containing configuration files (default: current directory)
--ssh-key-path string   Path to SSH private key or directory (default: ~/.ssh/)
--log-level string      Log level: debug, info, warn, error (default: info)
--log-file string       Log file path (default: $XDG_STATE_HOME/spooky/logs/spooky.log)
--dry-run              Show what would be done without making changes
--verbose              Enable verbose output
--quiet                Suppress all output except errors
```

## Core Commands

### `spooky execute`
Execute configuration files or remote sources.

```bash
spooky execute <source> [flags]
```

**Arguments:**
- `<source>`: Local file path, Git repository, S3 bucket, or HTTP URL

**Flags:**
```bash
--hosts string          Comma-separated list of target hosts
--parallel int          Number of parallel executions (default: 5)
--timeout duration      Execution timeout per host (default: 30s)
--retry int             Number of retry attempts (default: 3)
--tags string           Comma-separated list of tags to execute
--skip-tags string      Comma-separated list of tags to skip
```

**Examples:**
```bash
spooky execute config.hcl
spooky execute github.com/org/repo/configs/
spooky execute s3://my-bucket/configs/
spooky execute config.hcl --hosts web-001,web-002 --parallel 10
```

### `spooky validate`
Validate configuration files and templates.

```bash
spooky validate <source> [flags]
```

**Arguments:**
- `<source>`: Local file path, Git repository, S3 bucket, or HTTP URL

**Flags:**
```bash
--schema string         Path to schema file for validation
--strict               Enable strict validation mode
--format string        Output format: text, json, yaml (default: text)
```

**Examples:**
```bash
spooky validate config.hcl
spooky validate github.com/org/repo/configs/ --strict
spooky validate config.hcl --format json
```

### `spooky list`
List resources, facts, or configurations.

```bash
spooky list <resource> [flags]
```

**Resources:**
- `servers`: List known servers
- `facts`: List available facts
- `templates`: List available templates
- `configs`: List configuration files
- `actions`: List available actions

**Flags:**
```bash
--format string        Output format: table, json, yaml (default: table)
--filter string        Filter results by expression
--sort string          Sort field (default: name)
--reverse              Reverse sort order
```

**Examples:**
```bash
spooky list servers
spooky list facts --format json
spooky list servers --filter "os=ubuntu"
```

## Facts Management

### `spooky facts`
Manage server facts and fact collection.

```bash
spooky facts <subcommand> [flags]
```

#### `spooky facts gather`
Gather facts from target servers.

```bash
spooky facts gather [hosts] [flags]
```

**Arguments:**
- `[hosts]`: Comma-separated list of hosts (optional, uses inventory if not specified)

**Flags:**
```bash
--inventory string      Path to inventory file
--parallel int          Number of parallel fact gathering (default: 10)
--timeout duration      Timeout per host (default: 60s)
--facts string          Comma-separated list of fact types to gather
--update               Update existing facts instead of replacing
--cache-dir string      Directory for fact caching
```

**Examples:**
```bash
spooky facts gather
spooky facts gather web-001,web-002,db-001
spooky facts gather --inventory hosts.txt --parallel 20
spooky facts gather --facts hardware,network,os --update
```

#### `spooky facts import`
Import facts from external sources.

```bash
spooky facts import <source> [flags]
```

**Arguments:**
- `<source>`: Local JSON file, Git repository, S3 bucket, or HTTP URL

**Flags:**
```bash
--merge               Merge with existing facts instead of replacing
--validate            Validate facts before importing
--format string        Source format: json, yaml, csv (default: auto-detect)
--mapping string       Path to field mapping configuration
```

**Examples:**
```bash
spooky facts import facts.json
spooky facts import github.com/org/repo/facts/ --merge
spooky facts import s3://my-bucket/facts/ --validate
```

#### `spooky facts export`
Export facts to various formats.

```bash
spooky facts export [flags]
```

**Flags:**
```bash
--format string        Output format: json, yaml, csv, table (default: json)
--output string        Output file path (default: stdout)
--filter string        Filter facts by expression
--fields string        Comma-separated list of fields to include
--pretty              Pretty-print JSON output
```

**Examples:**
```bash
spooky facts export --format json --output facts.json
spooky facts export --filter "os=ubuntu" --format table
spooky facts export --fields hostname,ip,os --format csv
```

#### `spooky facts validate`
Validate facts against rules and schemas.

```bash
spooky facts validate [flags]
```

**Flags:**
```bash
--rules string         Path to validation rules file
--schema string        Path to schema file
--strict               Enable strict validation mode
--format string        Output format: text, json, html (default: text)
--output string        Output file path (default: stdout)
```

**Examples:**
```bash
spooky facts validate --rules validation-rules.yaml
spooky facts validate --schema facts-schema.json --strict
spooky facts validate --format html --output validation-report.html
```

#### `spooky facts query`
Query facts using expressions and filters.

```bash
spooky facts query <expression> [flags]
```

**Arguments:**
- `<expression>`: Query expression or filter

**Flags:**
```bash
--format string        Output format: table, json, yaml (default: table)
--output string        Output file path (default: stdout)
--limit int            Maximum number of results
--sort string          Sort field
--reverse              Reverse sort order
```

**Examples:**
```bash
spooky facts query "os=ubuntu and cpu_cores>=4"
spooky facts query "tags.environment=production"
spooky facts query "memory_gb>8" --limit 10 --sort memory_gb
```

## Template Management

### `spooky templates`
Manage configuration templates.

```bash
spooky templates <subcommand> [flags]
```

#### `spooky templates list`
List available templates.

```bash
spooky templates list [flags]
```

**Flags:**
```bash
--format string        Output format: table, json, yaml (default: table)
--filter string        Filter templates by expression
--sort string          Sort field (default: name)
```

**Examples:**
```bash
spooky templates list
spooky templates list --filter "type=nginx"
spooky templates list --format json
```

#### `spooky templates validate`
Validate template syntax and structure.

```bash
spooky templates validate <template> [flags]
```

**Arguments:**
- `<template>`: Path to template file

**Flags:**
```bash
--schema string        Path to template schema file
--strict               Enable strict validation mode
--format string        Output format: text, json (default: text)
```

**Examples:**
```bash
spooky templates validate nginx.conf.tmpl
spooky templates validate template.hcl --strict
```

#### `spooky templates render`
Render template locally with sample data.

```bash
spooky templates render <template> [flags]
```

**Arguments:**
- `<template>`: Path to template file

**Flags:**
```bash
--data string          Path to data file (JSON/YAML)
--output string        Output file path (default: stdout)
--facts string         Comma-separated list of fact types to include
--server string        Server ID to use for fact simulation
```

**Examples:**
```bash
spooky templates render nginx.conf.tmpl --data sample-data.json
spooky templates render config.hcl --facts os,network --server web-001
```

## Server Management

### `spooky servers`
Manage server information and connectivity.

```bash
spooky servers <subcommand> [flags]
```

#### `spooky servers list`
List known servers.

```bash
spooky servers list [flags]
```

**Flags:**
```bash
--format string        Output format: table, json, yaml (default: table)
--filter string        Filter servers by expression
--sort string          Sort field (default: hostname)
--fields string        Comma-separated list of fields to display
```

**Examples:**
```bash
spooky servers list
spooky servers list --filter "os=ubuntu"
spooky servers list --fields hostname,ip,os --format json
```

#### `spooky servers facts`
Get facts for specific server.

```bash
spooky servers facts <server> [flags]
```

**Arguments:**
- `<server>`: Server ID or hostname

**Flags:**
```bash
--format string        Output format: table, json, yaml (default: table)
--fields string        Comma-separated list of fields to display
--refresh              Refresh facts from server
```

**Examples:**
```bash
spooky servers facts web-001
spooky servers facts web-001 --format json
spooky servers facts web-001 --refresh
```

#### `spooky servers ping`
Test connectivity to server.

```bash
spooky servers ping <server> [flags]
```

**Arguments:**
- `<server>`: Server ID or hostname

**Flags:**
```bash
--timeout duration      Ping timeout (default: 10s)
--count int            Number of ping attempts (default: 3)
--ssh-only             Test only SSH connectivity
```

**Examples:**
```bash
spooky servers ping web-001
spooky servers ping web-001 --count 5
spooky servers ping web-001 --ssh-only
```

## Configuration Management

### `spooky config`
Manage configuration files and settings.

```bash
spooky config <subcommand> [flags]
```

#### `spooky config validate`
Validate configuration file syntax.

```bash
spooky config validate <file> [flags]
```

**Arguments:**
- `<file>`: Path to configuration file

**Flags:**
```bash
--schema string        Path to schema file
--strict               Enable strict validation mode
--format string        Output format: text, json (default: text)
```

**Examples:**
```bash
spooky config validate config.hcl
spooky config validate config.hcl --strict
```

#### `spooky config lint`
Lint configuration files for best practices.

```bash
spooky config lint <file> [flags]
```

**Arguments:**
- `<file>`: Path to configuration file

**Flags:**
```bash
--rules string         Path to linting rules file
--format string        Output format: text, json (default: text)
--fix                  Automatically fix issues where possible
```

**Examples:**
```bash
spooky config lint config.hcl
spooky config lint config.hcl --fix
```

#### `spooky config format`
Format configuration files.

```bash
spooky config format <file> [flags]
```

**Arguments:**
- `<file>`: Path to configuration file

**Flags:**
```bash
--output string        Output file path (default: overwrite input)
--check                Check if file is formatted without modifying
```

**Examples:**
```bash
spooky config format config.hcl
spooky config format config.hcl --check
```

## Help and Documentation

### `spooky help`
Show help information.

```bash
spooky help [command]
```

**Examples:**
```bash
spooky help
spooky help facts
spooky help facts gather
```

### `spooky version`
Show version information.

```bash
spooky version [flags]
```

**Flags:**
```bash
--format string        Output format: text, json (default: text)
--short                Show only version number
```

**Examples:**
```bash
spooky version
spooky version --format json
spooky version --short
```

## Environment Variables

The following environment variables can be used to set defaults:

```bash
SPOOKY_CONFIG_DIR      Default configuration directory
SPOOKY_SSH_KEY_PATH    Default SSH key path
SPOOKY_LOG_LEVEL       Default log level
SPOOKY_LOG_FILE        Default log file path
SPOOKY_CACHE_DIR       Default cache directory
SPOOKY_FACTS_DB_PATH   Default facts database path
```

## Exit Codes

- `0`: Success
- `1`: General error
- `2`: Configuration error
- `3`: Validation error
- `4`: Connection error
- `5`: Permission error
- `6`: Timeout error

## Examples

### Basic Usage
```bash
# Execute configuration
spooky execute config.hcl

# Gather facts from servers
spooky facts gather

# List servers
spooky list servers

# Validate configuration
spooky validate config.hcl
```

### Advanced Usage
```bash
# Execute from Git repository with specific hosts
spooky execute github.com/org/repo/configs/ --hosts web-001,web-002 --parallel 10

# Gather specific facts and update existing data
spooky facts gather --facts hardware,network --update --parallel 20

# Query facts with complex filter
spooky facts query "os=ubuntu and cpu_cores>=4 and memory_gb>=8" --limit 5

# Import facts from S3 and merge
spooky facts import s3://my-bucket/facts/ --merge --validate

# Render template with server-specific data
spooky templates render nginx.conf.tmpl --facts os,network --server web-001
```

### CI/CD Integration
```bash
# Validate configuration in CI
spooky config validate config.hcl --strict

# Execute with dry-run for safety
spooky execute config.hcl --dry-run --verbose

# Export facts for reporting
spooky facts export --format json --output facts-report.json
``` 