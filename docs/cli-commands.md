# Spooky CLI Commands Reference

This document provides a complete reference for all spooky CLI commands and their available flags.

## Global Flags

All commands support these global flags:

- `--config` - Path to config directory
- `--verbose` - Enable verbose output  
- `--quiet` - Suppress output
- `--log-level` - Log level (debug, info, warn, error)
- `--log-file` - Log file path
- `--version` - Emit the version

## Command Structure

Spooky CLI follows a "spooky noun-verb" pattern for most commands, with the exception of `spooky project info` which uses a noun-noun pattern.

## Main Command Groups

### Actions Commands

Manage and execute actions on remote machines.

#### `spooky actions list <project>`

List actions in a project.

**Flags:**
- `--output` - Output format (hcl, json)
- `--tags` - Tag-based filter (simple key-value filtering)
- `--machines` - Machine-specific filter (target specific machines)

**Examples:**
```bash
# List all actions
spooky actions list ./my-project --output json

# List actions for specific machines
spooky actions list ./my-project --machines server1,server2

# List actions for machines with specific tags
spooky actions list ./my-project --tags environment=production,role=web
```

#### `spooky actions run <project>`

Execute actions on remote machines.

**Flags:**
- `--action` - Action name to execute
- `--machines` - Target machines
- `--dry-run` - Dry run mode
- `--decrypt` - Decrypt encrypted content
- `--parallel` - Parallel execution count (2 or larger)
- `--tags` - Tag-based targeting

**Examples:**
```bash
# Run all actions
spooky actions run ./my-project --machines server1,server2

# Run specific action
spooky actions run ./my-project --action update-system --dry-run

# Run with parallel execution
spooky actions run ./my-project --parallel 4 --tags production
```

#### `spooky actions validate <project>`

Validate actions configuration.

**Flags:**
- NONE. actions validate has no output options.

### Facts Commands

Manage machine facts collection and storage.

#### `spooky facts gather <project>`

Gather facts from machines.

**Flags:**
- `--parallel` - Parallel execution count (2 or larger)
- `--machines` - Target specific machine
- `--tags` - Tag-based targeting

**Examples:**
```bash
# Gather facts from all machines
spooky facts gather ./my-project

# Gather facts from specific machines
spooky facts gather ./my-project --machines server1,server2

# Gather facts from machines with specific tags
spooky facts gather ./my-project --tags environment=production,role=web

# Gather facts with parallel execution
spooky facts gather ./my-project --parallel 4 --tags web-servers
```

#### `spooky facts list <project>`

List stored facts.

**Flags:**
- `--output` - Output format (json, hcl)
- `--machines` - Machine-specific filter
- `--tags` - Tag-based filter

**Examples:**
```bash
# List facts for all machines
spooky facts list ./my-project --output json

# List facts for specific machines
spooky facts list ./my-project --output json --machines server1,server2

# List facts for machines with specific tags
spooky facts list ./my-project --output json --tags environment=production
```

#### `spooky facts validate <project>`

Validate stored facts.

**Flags:**
- NONE. validation has no output or compare flag

**Examples:**
```bash
# Validate stored facts
spooky facts validate ./my-project
```

### Machines Commands

Manage machine inventory and connectivity.

#### `spooky machines list <project>`

List machines in inventory.

**Flags:**
- `--tags` - Tag-based filter

**Examples:**
```bash
# List all machines
spooky machines list ./my-project --output json

# List machines with specific tags
spooky machines list ./my-project --tags environment=production

# List machines with multiple tag criteria
spooky machines list ./my-project --tags environment=production,role=web
```

#### `spooky machines ping <project>`

Test machine connectivity.

**Flags:**
- `--machines` - Target specific machine
- `--tags` - Tag-based targeting

**Examples:**
```bash
# Ping all machines
spooky machines ping ./my-project

# Ping specific machines
spooky machines ping ./my-project --machines server1,server2

# Ping machines with specific tags
spooky machines ping ./my-project --tags environment=production,role=web
```

#### `spooky machines export <project>`

Export machine inventory.

**Flags:**
- `--output` - Output file path
- `--format` - JSON or HCL

**Example:**
```bash
spooky machines export ./my-project --format json --output inventory.json
```

### Project Commands

Manage project configuration and structure.

#### `spooky project init <directory>`

Initialize a new project.

**Flags:**
- `--name` - Project name
- `--description` - Project description
- `--author` - Project author

**Example:**
```bash
spooky project init ./new-project --name "My Project" --description "Production deployment"
```

#### `spooky project info <project>`

Show project information.

**Flags:**
- NONE

**Example:**
```bash
spooky project info ./my-project
```

#### `spooky project validate <project>`

Validate project structure as well as all other fact, variable, machine inventory, and template validation.

**Flags:**
- NONE. validation does not have output

**Example:**
```bash
spooky project validate ./my-project
```

#### `spooky project encrypt <project>`

Encrypt project data.

**Flags:**
- `--variables` - Re-encrypt variables only
- `--facts` - Re-encrypt facts only

**Examples:**
```bash
# Encrypt all project data
spooky project encrypt ./my-project

# Encrypt only variables
spooky project encrypt ./my-project --variables
```

### Templates Commands

Manage configuration templates.

#### `spooky templates list <project>`

List templates.

**Flags:**
- NONE. Listing has no output

**Example:**
```bash
spooky templates list ./my-project
```

#### `spooky templates render <project> <template>`

Render a template.

**Flags:**
- `--output` - Output file path
- `--dry-run` - Dry run mode
- `--preview` - Preview mode

**Example:**
```bash
spooky templates render ./my-project templates/nginx.conf --output /tmp/nginx.conf
```

#### `spooky templates validate <project>`

Validate templates.

**Flags:**
- `--templates` - Single or Multiple template paths

**Examples:**
```bash
# Validate all templates
spooky templates validate ./my-project

# Validate specific template
spooky templates validate ./my-project --templates templates/nginx.conf
```

### Variables Commands

Manage project variables.

#### `spooky variables list <project>`

List variables.

**Example:**
```bash
spooky variables list ./my-project
```

#### `spooky variables validate <project>`

Validate variables.

**Flags:**
- NONE. validation doesn't have an output flag

**Example:**
```bash
spooky variables validate ./my-project
```

## Utility Commands

### Configuration Management

#### `spooky config show`

Show config path.

**Flags:**
- `--config` - Override config file path

**Example:**
```bash
spooky config show
```

#### `spooky config validate`

Validate config file.

**Flags:**
- `--config` - Config file to validate

**Example:**
```bash
spooky config validate --config /path/to/config/
```

#### `spooky config reencrypt`

Encrypt config data.

**Example:**
```bash
spooky config encrypt
```

### Shell Completion

#### `spooky completion generate <shell>`

Generate completion script.

**Flags:**
- `--output` - Output file path

**Supported shells:** bash, zsh, fish

**Example:**
```bash
spooky completion generate bash --output ~/.bash_completion.d/spooky
```

### Log Management

#### `spooky logs show <project>`

Show logs.

**Flags:**
- `--log-file` - Log file path

**Example:**
```bash
spooky logs show --log-file <project>/logs/logfile.log
```

## Alias Commands

These commands are aliases for other commands:

- **`spooky run <project>`** - Alias for `spooky actions run <project>`
- **`spooky init <directory>`** - Alias for `spooky project init <directory>`
- **`spooky ping <project>`** - Alias for `spooky machines ping <project>`

## Variable and Facts Re-encryption Aliases

- **`spooky variables reencrypt <project>`** - Alias for `spooky project reencrypt <project> --variables`
- **`spooky facts reencrypt <project>`** - Alias for `spooky project reencrypt <project> --facts`

## Version Information

#### `spooky --version`

Show version information.

**Note:** No `-v` short flag is available.

**Example:**
```bash
spooky --version
```

## Important Notes

### Command Conventions

1. **Export Formats**: Export format must always be explicitly specified (no default)
2. **Parallel Execution**: The `--parallel` flag must be 2 or larger (0 and 1 are invalid)
3. **Machine Export**: Machine inventory can only be exported to JSON format
4. **Decrypt Flag**: The `--decrypt` flag appears at the end of CLI commands
5. **Project Validation**: Uses both project-directory.schema.hcl (directory structure) and project.schema.hcl (configuration content)
6. **Facts Validation**: Defaults to verifying stored facts, with `--compare` to gather fresh facts

### Filtering Options

Most commands support three types of filtering, each with different use cases:

1. **`--machines`** - Simple machine name targeting (e.g., `server1,server2`)
   - Use when you know the exact machine names
   - Fastest and most direct targeting method
   - Example: `--machine web-server-1,db-server-2`

2. **`--tags`** - Tag-based filtering using key-value pairs
   - Use for logical grouping of machines by attributes
   - Supports multiple tag criteria (AND logic)
   - Example: `--tags environment=production,role=web`


### Output Formats

Supported output formats:
- `json` - JSON format
- `hcl` - HCL format

### Project Structure

Projects should follow the standard spooky project structure with:
- `project.hcl` - Project configuration
- `machines.hcl` and `machines/` - Machine inventories
- `actions.hcl` and `actions` - Action files
- `templates/` - Template files
- `variables.hcl` and `variables/` - Variable files
- `facts.db` - Facts database

### Environment Variables

- `SPOOKY_FACTS_PATH` - Facts database path
- `SPOOKY_FACTS_FORMAT` - Storage format (badgerdb or json)
