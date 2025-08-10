# Proposed Script Execution Design

## Overview

The user proposes a unified script execution system that works consistently across local and remote machines, with clear separation between static scripts and templated scripts.

## Core Requirements

### 1. Universal Execution
- **Local execution**: Scripts run on the machine where Spooky is executed
- **Remote execution**: Scripts run on target machines via SSH
- **Same interface**: Commands and scripts work identically in both environments

### 2. Script Organization
- **`files/` directory**: Static scripts without variables
  - Example: `files/install-nginx.sh`
  - Content: Fixed, no templating needed
- **`templates/` directory**: Scripts with variables/templating
  - Example: `templates/install-service.sh.tmpl`
  - Content: Contains variables like `{{.service_name}}`, `{{.version}}`

## Proposed Execution Flow

### For Static Scripts (`files/`)

```hcl
action "install_nginx" {
  type = "script"
  script = "files/install-nginx.sh"
  machines = ["web-server-1", "web-server-2"]
}
```

**Execution Process**:
1. **Local execution**:
   - Copy `files/install-nginx.sh` to local temp file (`/tmp/spooky-<action_name>-<script_name>-<timestamp>.sh`)
   - Set executable permissions (`chmod 0755`)
   - Execute locally using `exec.CommandContext()`
   - Clean up temp file

2. **Remote execution**:
   - Copy `files/install-nginx.sh` to remote temp file (`/tmp/spooky-<action_name>-<script_name>-<timestamp>.sh`)
   - Upload via SCP to target machine
   - Set executable permissions remotely (`chmod +x`)
   - Execute via SSH
   - Clean up remote temp file

### For Templated Scripts (`templates/`)

```hcl
action "install_service" {
  type = "script"
  script = "templates/install-service.sh.tmpl"
  variables = {
    service_name = "nginx"
    version = "1.18.0"
    config_path = "/etc/nginx/nginx.conf"
  }
  machines = ["web-server-1"]
}
```

**Execution Process**:
1. **Template rendering**:
   - Load template from `templates/install-service.sh.tmpl`
   - Render with variables (replace `{{.service_name}}` with `nginx`)
   - Generate final script content

2. **Local execution**:
   - Write rendered content to local temp file (`/tmp/spooky-<action_name>-<template_name>-<timestamp>.sh`)
   - Set executable permissions
   - Execute locally
   - Clean up temp file

3. **Remote execution**:
   - Write rendered content to remote temp file (`/tmp/spooky-<action_name>-<template_name>-<timestamp>.sh`)
   - Upload via SCP
   - Set executable permissions remotely
   - Execute via SSH
   - Clean up remote temp file

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
1. **Action identification**: `spooky-install_nginx-` shows which action created the file
2. **Script identification**: `install-nginx-` shows the original script name
3. **Timestamp**: `20241201-143022` provides chronological ordering
4. **Consistent prefix**: `spooky-` makes it easy to find all Spooky temp files
5. **Unique naming**: Timestamp ensures no conflicts between parallel executions

### Logging Integration:
```go
func (a *actorImpl) createTemporaryScript(scriptContent, actionName, scriptPath string) (string, error) {
    // Generate debug-friendly temp file name
    timestamp := time.Now().Format("20060102-150405")
    scriptName := filepath.Base(scriptPath)
    if scriptName == "" {
        scriptName = "inline"
    }
    
    tempFileName := fmt.Sprintf("spooky-%s-%s-%s.sh", actionName, scriptName, timestamp)
    tempFilePath := filepath.Join(os.TempDir(), tempFileName)
    
    // Log the temp file creation
    a.logger.Debug("Creating temporary script file",
        spookylogging.String("action", actionName),
        spookylogging.String("script", scriptPath),
        spookylogging.String("temp_file", tempFilePath))
    
    // Create and write the file
    if err := os.WriteFile(tempFilePath, []byte(scriptContent), 0755); err != nil {
        return "", fmt.Errorf("failed to create temporary script: %w", err)
    }
    
    return tempFilePath, nil
}
```

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

## Implementation Changes Needed

### 1. Action Type Enhancement

```go
type Action struct {
    // ... existing fields ...
    
    // Enhanced script handling
    Script     string            `hcl:"script,optional"`      // File path or content
    ScriptFile string            `hcl:"script_file,optional"` // Explicit file reference
    Variables  map[string]string `hcl:"variables,optional"`   // For templated scripts
}
```

### 2. Script Resolution Logic

```go
func (a *actorImpl) resolveScript() (string, error) {
    // If script field contains file path
    if strings.HasPrefix(a.action.Script, "files/") || 
       strings.HasPrefix(a.action.Script, "templates/") {
        
        // Check if it's a template
        if strings.HasSuffix(a.action.Script, ".tmpl") {
            return a.renderTemplateScript(a.action.Script, a.action.Variables)
        } else {
            // Static script - read file content
            return a.readScriptFile(a.action.Script)
        }
    }
    
    // Inline script content (existing behavior)
    return a.action.Script, nil
}
```

### 3. Template Rendering

```go
func (a *actorImpl) renderTemplateScript(templatePath string, variables map[string]string) (string, error) {
    // Load template from templates/ directory
    templateContent, err := a.loadTemplate(templatePath)
    if err != nil {
        return "", fmt.Errorf("failed to load template: %w", err)
    }
    
    // Render template with variables
    rendered, err := a.templateEngine.Render(templateContent, variables)
    if err != nil {
        return "", fmt.Errorf("failed to render template: %w", err)
    }
    
    return rendered, nil
}
```

### 4. File Reading

```go
func (a *actorImpl) readScriptFile(scriptPath string) (string, error) {
    // Resolve path relative to project root
    fullPath := filepath.Join(a.projectPath, scriptPath)
    
    // Read file content
    content, err := os.ReadFile(fullPath)
    if err != nil {
        return "", fmt.Errorf("failed to read script file: %w", err)
    }
    
    return string(content), nil
}
```

### 5. Universal Execution Interface

```go
func (a *actorImpl) executeScript(ctx context.Context, scriptContent string, machine string) (string, int, error) {
    if machine == "local" || machine == "localhost" {
        return a.executeScriptLocally(ctx, scriptContent)
    } else {
        return a.executeScriptRemotely(ctx, scriptContent, machine)
    }
}
```

## Benefits of This Design

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
- Supports both inline content and file references
- Backward compatible with existing script actions

### 6. Debugging
- Clear, predictable temp file names
- Easy to correlate logs with temp files
- Simple cleanup and troubleshooting

## Example Usage

### Static Script
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

### Templated Script
```hcl
# templates/install-service.sh.tmpl
#!/bin/bash
SERVICE_NAME="{{.service_name}}"
VERSION="{{.version}}"
CONFIG_PATH="{{.config_path}}"

echo "Installing {{.service_name}} version {{.version}}"
apt-get update
apt-get install -y {{.service_name}}=${VERSION}
systemctl enable {{.service_name}}
systemctl start {{.service_name}}
```

```hcl
# actions.hcl
action "install_nginx" {
  type = "script"
  script = "templates/install-service.sh.tmpl"
  variables = {
    service_name = "nginx"
    version = "1.18.0"
    config_path = "/etc/nginx/nginx.conf"
  }
  machines = ["web-server-1"]
}
```

This design provides a clean, consistent, and flexible approach to script execution that works seamlessly across local and remote environments, with excellent debugging support through predictable temporary file naming.
