# Comprehensive Script Execution Guide for Spooky

## Overview

This guide covers the current state of script execution in Spooky, proposed improvements, and how different action types function. It addresses the evolution from inline script content to a more practical file-based approach with excellent debugging support.

## Current State vs Proposed Design

### Current Implementation (Inline Scripts - Deprecated)

**How it works now:**
- Script actions contain script content as strings in the action definition
- Temporary files are created from inline content for execution
- Only local execution is currently implemented
- **This pattern is deprecated and will not be supported in future versions**

**Example (Deprecated - Will Be Removed):**
```hcl
action "setup_nginx" {
  type = "script"
  script = """
    #!/bin/bash
    apt-get update
    apt-get install -y nginx
    systemctl enable nginx
    systemctl start nginx
  """
}
```

**Current limitations:**
- Script content embedded in action definitions (deprecated antipattern)
- No script reuse across actions
- Limited script size due to HCL constraints
- Temporary file overhead for every execution
- No remote execution support yet
- Poor maintainability and version control
- Difficult to debug and test scripts independently

### Proposed Design (File-Based Scripts)

**How it should work:**
- Scripts stored in project directories (`files/` and `templates/`)
- Actions reference script files instead of containing content
- Universal execution (local and remote)
- Debugging-friendly temporary file naming

**Example (Recommended Pattern):**
```hcl
action "install_nginx" {
  type = "script"
  script = "files/install-nginx.sh"
  machines = ["web-server-1"]
}
```

## Action Type Functional Differences

### Command Actions (`type = "command"`)

**Purpose**: Execute a single shell command directly

**Execution**: Runs the command string directly on target machines

**Use Case**: Simple, one-off commands like `systemctl restart nginx`

**Configuration**: Requires `command` field with the shell command

**Security**: Validates against shell injection (blocks `;&|`$` characters)

**Example:**
```hcl
action "restart_nginx" {
  type = "command"
  command = "systemctl restart nginx"
  machines = ["web-server-1"]
}
```

### Script Actions (`type = "script"`)

**Purpose**: Execute multi-line script content

**Current Execution Process**:
1. Creates a temporary script file from the script content
2. Executes the script file on target machines
3. Automatically cleans up the temporary file

**Proposed Execution Process**:
1. **File resolution**: Read script from `files/` or render from `templates/`
2. **Temporary file creation**: Copy to debug-friendly temp file
3. **Execution**: Run locally or remotely via SSH
4. **Cleanup**: Remove temporary file

**Use Case**: Complex multi-step operations, conditional logic, loops

**Configuration**: Requires `script` field with file path only

**Advantage**: Can contain complex logic that would be unsafe as a single command

### Template Actions (`type = "template_deploy"`)

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

**Example:**
```hcl
action "deploy_nginx_config" {
  type = "template_deploy"
  template {
    source = "templates/nginx.conf.tmpl"
    destination = "/etc/nginx/nginx.conf"
    backup = true
    permissions = "0644"
  }
  machines = ["web-server-1"]
}
```

## Script Organization Strategy

### Directory Structure

```
project/
├── files/                    # Static scripts (no variables)
│   ├── install-nginx.sh
│   ├── setup-database.sh
│   └── backup-data.sh
├── templates/                # Scripts with variables
│   ├── install-service.sh.tmpl
│   ├── configure-app.sh.tmpl
│   └── deploy-config.sh.tmpl
└── actions.hcl              # Action definitions
```

### File Types

#### Static Scripts (`files/`)
- **Purpose**: Scripts without variables or templating
- **Content**: Fixed, reusable across multiple actions
- **Example**: `files/install-nginx.sh`

```bash
#!/bin/bash
apt-get update
apt-get install -y nginx
systemctl enable nginx
systemctl start nginx
```

#### Templated Scripts (`templates/`)
- **Purpose**: Scripts with variables and dynamic content
- **Content**: Contains template variables like `{{.service_name}}`
- **Example**: `templates/install-service.sh.tmpl`

```bash
#!/bin/bash
SERVICE_NAME="{{.service_name}}"
VERSION="{{.version}}"

echo "Installing {{.service_name}} version {{.version}}"
apt-get update
apt-get install -y {{.service_name}}=${VERSION}
systemctl enable {{.service_name}}
systemctl start {{.service_name}}
```

## Why Temporary Files Are Needed

### The Core Issue: Script Content vs File Path

**Current implementation**: Script actions contain script content as strings, not file paths

```go
type Action struct {
    Script string `hcl:"script,optional"`  // Contains script content, not file path
}
```

**Execution requirements**:
- Executable programs require file paths, not content strings
- Scripts need executable permissions (`chmod 0755`)
- Proper file extensions (e.g., `.sh` for shell scripts)

**The solution**: Create temporary files to bridge the gap between script content and executable file requirements.

## Debugging-Friendly Temporary File Naming

### Naming Convention
```
/tmp/spooky-<action_name>-<script_name>-<timestamp>.sh
```

### Examples:
- **Static script**: `/tmp/spooky-install_nginx-install-nginx-20241201-143022.sh`
- **Templated script**: `/tmp/spooky-install_service-install-service-20241201-143022.sh`
- **Inline script**: `/tmp/spooky-custom_action-inline-20241201-143022.sh`

### Benefits for Debugging:
1. **Action identification**: Shows which action created the file
2. **Script identification**: Shows the original script name
3. **Timestamp**: Provides chronological ordering
4. **Consistent prefix**: Makes it easy to find all Spooky temp files
5. **Unique naming**: Ensures no conflicts between parallel executions

### Debugging Commands:
```bash
# Find all Spooky temp files
find /tmp -name "spooky-*" -type f

# Find temp files for a specific action
find /tmp -name "spooky-install_nginx-*" -type f

# Find temp files created in the last hour
find /tmp -name "spooky-*" -type f -mmin -60

# View the content of a temp file (if it still exists)
cat /tmp/spooky-install_nginx-install-nginx-20241201-143022.sh
```

## Implementation Strategy

### Phase 1: File-Only Script Resolution

```go
func (a *actorImpl) resolveScript() (string, error) {
    // Script field must contain file path
    if a.action.Script == "" {
        return "", fmt.Errorf("script field is required and must reference a file")
    }
    
    // Validate file path format
    if !strings.HasPrefix(a.action.Script, "files/") && 
       !strings.HasPrefix(a.action.Script, "templates/") {
        return "", fmt.Errorf("script must reference a file in files/ or templates/ directory")
    }
    
    // Check if it's a template
    if strings.HasSuffix(a.action.Script, ".tmpl") {
        return a.renderTemplateScript(a.action.Script, a.action.Variables)
    } else {
        // Static script - read file content
        return a.readScriptFile(a.action.Script)
    }
}
```

### Phase 2: Universal Execution Interface

```go
func (a *actorImpl) executeScript(ctx context.Context, scriptContent string, machine string) (string, int, error) {
    if machine == "local" || machine == "localhost" {
        return a.executeScriptLocally(ctx, scriptContent)
    } else {
        return a.executeScriptRemotely(ctx, scriptContent, machine)
    }
}
```

### Phase 3: Enhanced Action Type

```go
type Action struct {
    // ... existing fields ...
    
    // Enhanced script handling
    Script     string            `hcl:"script,optional"`      // File path or content
    Variables  map[string]string `hcl:"variables,optional"`   // For templated scripts
}
```

## Key Functional Differences Summary

| Aspect | Command | Script (Current) | Script (Proposed) | Template |
|--------|---------|------------------|-------------------|----------|
| **Content** | Single command string | Inline script content | File path or content | Template file + data |
| **Execution** | Direct shell execution | Temporary file execution | Universal temp file execution | Render + deploy |
| **Complexity** | Simple, linear | Complex logic possible | Complex logic possible | Dynamic content |
| **File Management** | None | Temporary file cleanup | Debug-friendly temp files | Backup, permissions, ownership |
| **Context Usage** | Limited | Limited | Limited | Full access to facts/variables |
| **Reusability** | None | None | High (file-based) | High (template-based) |
| **Use Case** | Quick commands | Complex operations | Complex operations | Configuration management |

## Migration Path

### Migration Required
- Existing inline script actions must be migrated to file-based scripts
- No backward compatibility for inline scripts in future versions
- File-based scripts are the only supported pattern

### Recommended Migration Steps:
1. **Extract common scripts** from actions to `files/` directory
2. **Create templates** for parameterized scripts in `templates/` directory
3. **Update actions** to reference files instead of inline content
4. **Test** both local and remote execution
5. **Clean up** old inline script actions

## Example Migration

### Before (Antipattern - Avoid This):
```hcl
action "install_nginx" {
  type = "script"
  script = """
    #!/bin/bash
    apt-get update
    apt-get install -y nginx
    systemctl enable nginx
    systemctl start nginx
  """
  machines = ["web-server-1"]
}
```

### After (Recommended Pattern):
```hcl
# files/install-nginx.sh
#!/bin/bash
apt-get update
apt-get install -y nginx
systemctl enable nginx
systemctl start nginx
```

```hcl
# actions.hcl
action "install_nginx" {
  type = "script"
  script = "files/install-nginx.sh"
  machines = ["web-server-1"]
}
```

## Benefits of the Proposed Design

### 1. Consistency
- Same action definition works for local and remote execution
- No need to duplicate scripts for different environments

### 2. Organization
- Clear separation between static and templated scripts
- Easy to find and manage scripts by purpose

### 3. Reusability
- Static scripts can be reused across multiple actions
- Templates can be parameterized for different use cases

### 4. Maintainability
- Scripts are version-controlled with the project
- Changes to scripts are tracked alongside action changes

### 5. Flexibility
- Supports file references only
- Clean, consistent approach without legacy baggage

### 6. Debugging
- Clear, predictable temp file names
- Easy to correlate logs with temp files
- Simple cleanup and troubleshooting

This comprehensive approach provides a clear path from the current deprecated inline script pattern to a more practical, maintainable, and debuggable file-based system. Inline scripts will not be supported in future versions.
