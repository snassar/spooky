# Actions System API Reference

## Overview

This document provides a comprehensive API reference for the spooky actions system. It covers all interfaces, types, methods, and implementation details for developers working with the actions system.

**Status: Partially Implemented** - The actions system has basic functionality but SSH-based action orchestration has known issues that need to be addressed.

## Core Interfaces

### ActionsIntegration Interface

The `ActionsIntegration` interface provides the primary entry point for actions operations:

```go
type ActionsIntegration interface {
    // LoadActions loads actions from the specified project path
    LoadActions(ctx context.Context, projectPath string) ([]interface{}, error)
    
    // ValidateActions validates action definitions
    ValidateActions(ctx context.Context, actions []interface{}) (*ValidationResult, error)
    
    // RunActions runs actions on target machines
    RunActions(ctx context.Context, projectPath string, actions []interface{}, machines []string) error
}
```

**Implementation Status**: ⚠️ **Partially Implemented** - Basic functionality exists but SSH-based orchestration has issues

### ActionManager Interface

The `ActionManager` interface provides action management and orchestration:

```go
type ActionManager interface {
    // LoadActions loads actions from project configuration
    LoadActions(ctx context.Context, projectPath string) ([]*spookytypesactions.Action, error)
    
    // ValidateActions validates action definitions
    ValidateActions(ctx context.Context, actions []*spookytypesactions.Action) (*ValidationResult, error)
    
    // RunActions runs actions on target machines
    RunActions(ctx context.Context, actions []*spookytypesactions.Action, machines []*spookytypes.Machine) error
}
```

**Implementation Status**: ⚠️ **Partially Implemented** - Basic loading and validation exist but orchestration has issues

## Current Implementation Status

### ✅ Working Components

1. **Action Loading**: Loading actions from HCL configuration files
2. **Action Validation**: Basic validation of action definitions
3. **Action Structure**: Proper action type definitions and structures
4. **CLI Integration**: `spooky actions run` command with filtering options
5. **Project Integration**: Actions loading from project configuration
6. **Basic Validation**: Action definition validation and error handling
7. **Filtering Support**: Support for machine, tag, and group filtering
8. **Dry Run Support**: Dry run mode for action validation

### ⚠️ Known Issues

1. **SSH-Based Orchestration**: SSH-based action orchestration has implementation issues
2. **Action Execution**: Actions cannot be properly executed on remote machines
3. **Parallel Processing**: No parallel action execution support
4. **Result Handling**: No action result collection or processing
5. **Template Integration**: No template rendering in actions
6. **Variable Integration**: No variable resolution in actions

### 🔄 In Progress

1. **SSH Orchestration Fixes**: Addressing SSH-based action orchestration issues
2. **Execution Improvements**: Implementing proper action execution
3. **Parallel Support**: Adding parallel action execution support

## Implementation Details

### Action Loading System

The actions system loads actions from HCL configuration files:

```go
type ActionLoader struct {
    logger spookylogging.Logger
}

func (l *ActionLoader) LoadActions(ctx context.Context, projectPath string) ([]*spookytypesactions.Action, error) {
    // Load actions.hcl file
    actionsPath := filepath.Join(projectPath, "actions.hcl")
    
    data, err := os.ReadFile(actionsPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read actions file: %w", err)
    }
    
    // Parse HCL configuration
    var config struct {
        Actions []*spookytypesactions.Action `hcl:"action,block"`
    }
    
    if err := hcl.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse actions configuration: %w", err)
    }
    
    return config.Actions, nil
}
```

**Supported Action Types:**
- **Command Actions**: Execute commands on remote machines
- **Script Actions**: Run scripts on remote machines
- **File Actions**: Copy files to/from remote machines
- **Template Actions**: Render templates and deploy to remote machines

### Action Validation System

Actions are validated against schemas and business rules:

```go
type ActionValidator struct {
    logger spookylogging.Logger
}

func (v *ActionValidator) ValidateActions(ctx context.Context, actions []*spookytypesactions.Action) (*spookytypes.ValidationResult, error) {
    var errors []spookyschemas.SchemaError
    var warnings []spookyschemas.SchemaError
    
    for i, action := range actions {
        // Validate action name
        if action.Name == "" {
            errors = append(errors, spookyschemas.SchemaError{
                Message: fmt.Sprintf("action[%d]: name is required", i),
            })
        }
        
        // Validate action description
        if action.Description == "" {
            warnings = append(warnings, spookyschemas.SchemaError{
                Message: fmt.Sprintf("action[%d]: description is recommended", i),
            })
        }
        
        // Validate action type
        if err := v.validateActionType(action); err != nil {
            errors = append(errors, spookyschemas.SchemaError{
                Message: fmt.Sprintf("action[%d]: %s", i, err.Error()),
            })
        }
        
        // Validate action targets
        if err := v.validateActionTargets(action); err != nil {
            errors = append(errors, spookyschemas.SchemaError{
                Message: fmt.Sprintf("action[%d]: %s", i, err.Error()),
            })
        }
    }
    
    return &spookytypes.ValidationResult{
        Valid:    len(errors) == 0,
        Errors:   errors,
        Warnings: warnings,
    }, nil
}
```

### Action Orchestration System

Actions are orchestrated through the SSH system (currently has issues):

```go
type ActionOrchestrator struct {
    sshManager spookyinterfaces.SSHManager
    logger     spookylogging.Logger
}

func (o *ActionOrchestrator) RunActions(ctx context.Context, actions []*spookytypesactions.Action, machines []*spookytypes.Machine) error {
    for _, action := range actions {
        // Filter machines based on action targets
        targetMachines := o.filterMachinesForAction(action, machines)
        
        // Run action on target machines
        for _, machine := range targetMachines {
            if err := o.runActionOnMachine(ctx, action, machine); err != nil {
                o.logger.Error("Failed to run action on machine", 
                    "action", action.Name, 
                    "machine", machine.Hostname, 
                    "error", err)
                return fmt.Errorf("failed to run action %s on %s: %w", action.Name, machine.Hostname, err)
            }
        }
    }
    
    return nil
}

func (o *ActionOrchestrator) runActionOnMachine(ctx context.Context, action *spookytypesactions.Action, machine *spookytypes.Machine) error {
    // Create SSH connection
    conn, err := o.sshManager.GetConnection(machine.Hostname, machine.Port, machine.User)
    if err != nil {
        return fmt.Errorf("failed to establish SSH connection: %w", err)
    }
    defer o.sshManager.ReturnConnection(conn)
    
    // Execute action based on type
    switch action.Type {
    case "command":
        return o.executeCommandAction(ctx, conn, action)
    case "script":
        return o.executeScriptAction(ctx, conn, action)
    case "file":
        return o.executeFileAction(ctx, conn, action)
    case "template":
        return o.executeTemplateAction(ctx, conn, action)
    default:
        return fmt.Errorf("unsupported action type: %s", action.Type)
    }
}
```

## Type Definitions

### Action Types

```go
// Action represents an action definition
type Action struct {
    // Action name (required)
    Name string `json:"name" hcl:"name"`
    
    // Action description (optional)
    Description string `json:"description,omitempty" hcl:"description,optional"`
    
    // Action type (command, script, file, template)
    Type string `json:"type" hcl:"type"`
    
    // Action targets (machines, tags, groups)
    Targets *ActionTargets `json:"targets,omitempty" hcl:"targets,optional"`
    
    // Action configuration
    Config *ActionConfig `json:"config,omitempty" hcl:"config,optional"`
    
    // Action metadata
    Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// ActionTargets defines action execution targets
type ActionTargets struct {
    // Target machines by name
    Machines []string `json:"machines,omitempty" hcl:"machines,optional"`
    
    // Target machines by tags
    Tags []string `json:"tags,omitempty" hcl:"tags,optional"`
    
    // Target machines by groups
    Groups []string `json:"groups,omitempty" hcl:"groups,optional"`
    
    // Complex filter query
    Filter string `json:"filter,omitempty" hcl:"filter,optional"`
}

// ActionConfig provides action-specific configuration
type ActionConfig struct {
    // Command to execute (for command actions)
    Command string `json:"command,omitempty" hcl:"command,optional"`
    
    // Script path (for script actions)
    Script string `json:"script,omitempty" hcl:"script,optional"`
    
    // File operations (for file actions)
    Files []*FileOperation `json:"files,omitempty" hcl:"files,optional"`
    
    // Template operations (for template actions)
    Templates []*TemplateOperation `json:"templates,omitempty" hcl:"templates,optional"`
    
    // Execution options
    Options *ExecutionOptions `json:"options,omitempty" hcl:"options,optional"`
}

// FileOperation represents a file operation
type FileOperation struct {
    // Source file path
    Source string `json:"source" hcl:"source"`
    
    // Destination file path
    Destination string `json:"destination" hcl:"destination"`
    
    // Operation type (copy, move, delete)
    Operation string `json:"operation" hcl:"operation"`
    
    // File permissions
    Permissions string `json:"permissions,omitempty" hcl:"permissions,optional"`
}

// TemplateOperation represents a template operation
type TemplateOperation struct {
    // Template source path
    Source string `json:"source" hcl:"source"`
    
    // Template destination path
    Destination string `json:"destination" hcl:"destination"`
    
    // Template data
    Data map[string]interface{} `json:"data,omitempty" hcl:"data,optional"`
    
    // Template permissions
    Permissions string `json:"permissions,omitempty" hcl:"permissions,optional"`
}

// ExecutionOptions provides execution configuration
type ExecutionOptions struct {
    // Parallel execution
    Parallel bool `json:"parallel,omitempty" hcl:"parallel,optional"`
    
    // Execution timeout
    Timeout time.Duration `json:"timeout,omitempty" hcl:"timeout,optional"`
    
    // Working directory
    WorkingDir string `json:"working_dir,omitempty" hcl:"working_dir,optional"`
    
    // Environment variables
    Environment map[string]string `json:"environment,omitempty" hcl:"environment,optional"`
}
```

### Action Context Types

```go
// ActionContext provides context for action execution
type ActionContext struct {
    // Project path
    ProjectPath string `json:"project_path" hcl:"project_path"`
    
    // Action being executed
    Action *Action `json:"action" hcl:"action"`
    
    // Target machine
    Machine *spookytypes.Machine `json:"machine" hcl:"machine"`
    
    // Execution timestamp
    Timestamp time.Time `json:"timestamp" hcl:"timestamp"`
    
    // Execution metadata
    Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// ActionResult represents the result of action execution
type ActionResult struct {
    // Action context
    Context *ActionContext `json:"context" hcl:"context"`
    
    // Execution success
    Success bool `json:"success" hcl:"success"`
    
    // Execution output
    Output string `json:"output,omitempty" hcl:"output,optional"`
    
    // Execution error
    Error string `json:"error,omitempty" hcl:"error,optional"`
    
    // Execution duration
    Duration time.Duration `json:"duration" hcl:"duration"`
    
    // Exit code
    ExitCode int `json:"exit_code,omitempty" hcl:"exit_code,optional"`
}
```

## Error Handling

### Action Errors

```go
// ActionError represents action execution errors
type ActionError struct {
    ActionName string `json:"action_name" hcl:"action_name"`
    MachineID  string `json:"machine_id" hcl:"machine_id"`
    Error      string `json:"error" hcl:"error"`
    Details    string `json:"details,omitempty" hcl:"details,optional"`
}

// ActionValidationError represents action validation errors
type ActionValidationError struct {
    Field   string `json:"field" hcl:"field"`
    Message string `json:"message" hcl:"message"`
    Value   string `json:"value,omitempty" hcl:"value,optional"`
}
```

### Validation Implementation

```go
// ValidateAction validates a single action
func (v *ActionValidator) ValidateAction(action *spookytypesactions.Action) error {
    if action == nil {
        return fmt.Errorf("action cannot be nil")
    }
    
    // Validate required fields
    if action.Name == "" {
        return fmt.Errorf("action name is required")
    }
    
    if action.Type == "" {
        return fmt.Errorf("action type is required")
    }
    
    // Validate action type
    validTypes := []string{"command", "script", "file", "template"}
    valid := false
    for _, t := range validTypes {
        if action.Type == t {
            valid = true
            break
        }
    }
    if !valid {
        return fmt.Errorf("invalid action type: %s (valid types: %v)", action.Type, validTypes)
    }
    
    // Validate action configuration based on type
    switch action.Type {
    case "command":
        if action.Config == nil || action.Config.Command == "" {
            return fmt.Errorf("command action requires command configuration")
        }
    case "script":
        if action.Config == nil || action.Config.Script == "" {
            return fmt.Errorf("script action requires script configuration")
        }
    case "file":
        if action.Config == nil || len(action.Config.Files) == 0 {
            return fmt.Errorf("file action requires file operations")
        }
    case "template":
        if action.Config == nil || len(action.Config.Templates) == 0 {
            return fmt.Errorf("template action requires template operations")
        }
    }
    
    return nil
}
```

## CLI Commands

### Actions Run Command

```bash
# Run all actions in a project
spooky actions run ./my-project

# Run actions with dry run mode
spooky actions run ./my-project --dry-run

# Run actions on specific machines
spooky actions run ./my-project --machines web-server,app-server

# Run actions on machines with specific tags
spooky actions run ./my-project --tags production,web

# Run actions on machines in specific groups
spooky actions run ./my-project --groups webservers,databases

# Run actions with parallel execution
spooky actions run ./my-project --parallel 4
```

### Actions Validation Command

```bash
# Validate actions in a project
spooky actions validate ./my-project

# Validate actions with verbose output
spooky actions validate ./my-project --verbose
```

## Integration Examples

### Basic Action Definition

```hcl
# actions.hcl
actions {
  action "update-system" {
    description = "Update system packages"
    type = "command"
    
    targets {
      tags = ["production", "web"]
    }
    
    config {
      command = "apt update && apt upgrade -y"
      
      options {
        timeout = "300s"
        parallel = true
      }
    }
  }
  
  action "deploy-config" {
    description = "Deploy configuration files"
    type = "template"
    
    targets {
      machines = ["web-server", "app-server"]
    }
    
    config {
      templates {
        source = "templates/nginx.conf.tmpl"
        destination = "/etc/nginx/nginx.conf"
        permissions = "0644"
      }
      
      options {
        timeout = "60s"
      }
    }
  }
}
```

### Action Loading and Validation

```go
// Action loading and validation example
func loadAndValidateActions(projectPath string) error {
    ctx := context.Background()
    
    // Create action manager
    manager := spookyactions.NewManager(loader, validator, logger)
    
    // Load actions
    actions, err := manager.LoadActions(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load actions: %w", err)
    }
    
    // Validate actions
    result, err := manager.ValidateActions(ctx, actions)
    if err != nil {
        return fmt.Errorf("failed to validate actions: %w", err)
    }
    
    if !result.Valid {
        fmt.Println("Action validation failed:")
        for _, error := range result.Errors {
            fmt.Printf("  - %s\n", error.Message)
        }
        return fmt.Errorf("action validation failed")
    }
    
    fmt.Printf("Loaded and validated %d actions\n", len(actions))
    return nil
}
```

### Action Execution

```go
// Action execution example
func executeActions(projectPath string, machines []*spookytypes.Machine) error {
    ctx := context.Background()
    
    // Create action manager
    manager := spookyactions.NewManager(loader, validator, logger)
    
    // Load actions
    actions, err := manager.LoadActions(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load actions: %w", err)
    }
    
    // Run actions
    err = manager.RunActions(ctx, actions, machines)
    if err != nil {
        return fmt.Errorf("failed to run actions: %w", err)
    }
    
    return nil
}
```

## Current Limitations

### Orchestration Limitations

1. **SSH Integration Issues**: SSH-based action orchestration has known problems
2. **No Parallel Execution**: Actions are executed sequentially, not in parallel
3. **No Result Collection**: Action results are not collected or processed
4. **No Error Recovery**: No error recovery or retry mechanisms
5. **No Progress Tracking**: No progress tracking or status reporting

### Integration Limitations

1. **No Template Integration**: Templates are not rendered in actions
2. **No Variable Integration**: Variables are not resolved in actions
3. **No Facts Integration**: Facts are not used in action conditions
4. **No Conditional Execution**: No conditional action execution based on facts or variables

### Execution Limitations

1. **Basic Command Execution**: Only basic command execution is supported
2. **No Script Support**: Script execution is not properly implemented
3. **No File Operations**: File copy/move operations are not implemented
4. **No Template Rendering**: Template rendering is not implemented

## Future Enhancements

### Planned Features

1. **SSH Orchestration Fixes**: Resolve SSH-based action orchestration issues
2. **Parallel Execution**: Implement parallel action execution
3. **Result Collection**: Collect and process action results
4. **Error Recovery**: Implement error recovery and retry mechanisms
5. **Progress Tracking**: Add progress tracking and status reporting
6. **Conditional Execution**: Support conditional action execution

### Integration Enhancements

1. **Template Integration**: Integrate template rendering in actions
2. **Variable Integration**: Integrate variable resolution in actions
3. **Facts Integration**: Use facts in action conditions and execution
4. **Advanced Execution**: Support advanced execution options and features

## Summary

The actions system provides basic action loading and validation capabilities but has significant limitations with SSH-based orchestration that need to be addressed. The system is functional for basic use cases but requires improvements for production use.

**Status**: ⚠️ **Partially Implemented** - Basic functionality exists but SSH-based orchestration has issues that need to be resolved.