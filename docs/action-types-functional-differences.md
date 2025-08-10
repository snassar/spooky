# Action Type Functional Differences in Spooky

Based on analysis of the codebase, here's how the different action types function in Spooky:

## Command Actions (`type = "command"`)

**Purpose**: Execute a single shell command directly

**Execution**: Runs the command string directly on target machines

**Use Case**: Simple, one-off commands like `systemctl restart nginx`

**Configuration**: Requires `command` field with the shell command

**Security**: Validates against shell injection (blocks `;&|`$` characters)

## Script Actions (`type = "script"`)

**Purpose**: Execute multi-line script content

**Execution Process**:
1. Creates a temporary script file from the script content
2. Executes the script file on target machines
3. Automatically cleans up the temporary file

**Use Case**: Complex multi-step operations, conditional logic, loops

**Configuration**: Requires `script` field with script content

**Advantage**: Can contain complex logic that would be unsafe as a single command

## Template Actions (`type = "template_deploy"`)

**Purpose**: Render and deploy configuration files with dynamic content

**Execution Process**:
1. **Loads** template from `templates/` directory
2. **Renders** template with context data (facts, variables, etc.)
3. **Creates backup** of existing file (if `backup = true`)
4. **Deploys** rendered content to destination
5. **Sets permissions** and ownership (if specified)

**Use Case**: Configuration management, dynamic file generation

**Configuration**: Requires `template` block with `source` and `destination`

**Features**:
- Template validation before deployment
- Backup creation
- Permission/ownership management
- Context-aware rendering

## Additional Template Types

- **`template_evaluate`**: Renders template but doesn't deploy (for testing/preview)
- **`template_validate`**: Validates template syntax without rendering
- **`template_cleanup`**: Removes deployed files and backups

## Key Functional Differences

| Aspect | Command | Script | Template |
|--------|---------|--------|----------|
| **Content** | Single command string | Multi-line script | Template file + data |
| **Execution** | Direct shell execution | Temporary file execution | Render + deploy |
| **Complexity** | Simple, linear | Complex logic possible | Dynamic content |
| **File Management** | None | Temporary file cleanup | Backup, permissions, ownership |
| **Context Usage** | Limited | Limited | Full access to facts/variables |
| **Use Case** | Quick commands | Complex operations | Configuration management |

## Summary

The choice between these types depends on the complexity of the operation and whether you need dynamic content generation or just command execution:

- **Commands**: For simple, direct shell operations
- **Scripts**: For complex multi-step operations with logic
- **Templates**: For configuration management with dynamic content
