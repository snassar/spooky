# Actions System API Reference

## Overview

This document provides a comprehensive API reference for the spooky actions system. It covers all interfaces, types, methods, and implementation details for developers working with the actions system.

**Status: Partially Implemented** - The actions system has basic functionality but SSH-based action orchestration has known issues that need to be addressed.

## Core Interfaces

### ActionsIntegration Interface

The `ActionsIntegration` interface provides the primary entry point for actions operations:

```go
type ActionsIntegration interface {
    // LoadActions loads actions from the given source
    LoadActions(ctx context.Context, source string) ([]spookytypes.Action, error)

    // ValidateActions validates actions
    ValidateActions(ctx context.Context, actions []spookytypes.Action) (*spookytypes.ValidationResult, error)

    // RunActions runs actions on the given machines
    RunActions(ctx context.Context, actions []spookytypes.Action, machines []spookytypes.Machine) ([]spookytypes.ActingResult, error)

    // GetSSHManager returns the SSH manager for authentication testing
    GetSSHManager() SSHManager
}
```

**Implementation Status**: ⚠️ **Partially Implemented** - Basic functionality exists but SSH-based orchestration has issues

### ActionManager Interface

The `ActionManager` interface provides action management and orchestration:

```go
type ActionManager interface {
    // LoadActions loads actions from the given source
    LoadActions(ctx context.Context, source string) ([]spookytypes.Action, error)

    // ValidateActions validates actions
    ValidateActions(ctx context.Context, actions []spookytypes.Action) (*spookytypes.ValidationResult, error)

    // RunActions runs actions on the given machines
    RunActions(ctx context.Context, actions []spookytypes.Action, machines []spookytypes.Machine) ([]spookytypes.ActingResult, error)

    // GetSSHManager returns the SSH manager for authentication testing
    GetSSHManager() SSHManager
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
9. **Action Planning**: Basic action planning and dependency resolution
10. **SSH Manager Integration**: SSH manager for authentication testing

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
type Manager struct {
    logger          spookytypeslogging.Logger
    validator       spookyinterfaces.ActionValidator
    sshManager      spookyinterfaces.SSHManager
    schemaValidator *spookyschemas.Validator
}

func NewManager(
    logger spookytypeslogging.Logger,
    validator spookyinterfaces.ActionValidator,
    sshManager spookyinterfaces.SSHManager,
    schemaValidator *spookyschemas.Validator,
) spookyinterfaces.ActionsIntegration {
    return &Manager{
        logger:          logger,
        validator:       validator,
        sshManager:      sshManager,
        schemaValidator: schemaValidator,
    }
}
```

### Action Loading Implementation

```go
// LoadActions loads actions from the specified source
func (m *Manager) LoadActions(ctx context.Context, source string) ([]spookytypes.Action, error) {
    m.logger.Info("Loading actions", map[string]interface{}{
        "source": source,
    })

    // Check if source is a directory
    if info, err := os.Stat(source); err == nil && info.IsDir() {
        return m.loadActionsFromDirectory(ctx, source)
    }

    // Check if source is a file
    if _, err := os.Stat(source); err == nil {
        return m.loadActionsFromFile(ctx, source)
    }

    return nil, fmt.Errorf("source not found: %s", source)
}

// loadActionsFromFile loads actions from a single HCL file
func (m *Manager) loadActionsFromFile(_ context.Context, filePath string) ([]spookytypes.Action, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read actions file: %w", err)
    }

    var config struct {
        Actions []*spookytypesactions.Action `hcl:"action,block"`
    }

    if err := hclsimple.Decode(filePath, data, nil, &config); err != nil {
        return nil, fmt.Errorf("failed to parse actions file: %w", err)
    }

    // Convert to interface slice
    actions := make([]spookytypes.Action, len(config.Actions))
    for i, action := range config.Actions {
        actions[i] = action
    }

    return actions, nil
}
```

### Action Validation Implementation

```go
// ValidateActions validates a collection of actions
func (m *Manager) ValidateActions(ctx context.Context, actions []spookytypes.Action) (*spookytypes.ValidationResult, error) {
    m.logger.Info("Validating actions", map[string]interface{}{
        "actions": len(actions),
    })

    var errors []spookyschemas.SchemaError
    var warnings []spookyschemas.SchemaError

    for i, action := range actions {
        // Validate individual action
        if err := m.validateAction(action); err != nil {
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

func (m *Manager) validateAction(action spookytypes.Action) error {
    // Validate required fields
    if action.Name == "" {
        return fmt.Errorf("action name is required")
    }

    if action.Type == "" {
        return fmt.Errorf("action type is required")
    }

    // Validate action type
    validTypes := []string{"command", "script", "template_deploy", "file_copy", "service_control"}
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

    return nil
}
```

### Action Running Implementation

```go
// RunActions runs the specified actions on the target machines
func (m *Manager) RunActions(ctx context.Context, actions []spookytypes.Action, machines []spookytypes.Machine) ([]spookytypes.ActingResult, error) {
    m.logger.Info("Running actions", map[string]interface{}{
        "actions":  len(actions),
        "machines": len(machines),
    })

    // Create acting session
    session := &spookytypesactions.ActingSession{
        SessionID:    "session-" + fmt.Sprintf("%d", ctx.Value("session_id")),
        Status:       "active",
        TotalActions: len(actions),
    }

    // Plan action running
    plan, err := m.createActionPlan(ctx, actions, machines)
    if err != nil {
        m.logger.Error("Failed to create action plan", err, map[string]interface{}{
            "actions":  len(actions),
            "machines": len(machines),
        })
        return nil, fmt.Errorf("failed to create action plan: %w", err)
    }

    // Run action plan
    results, err := m.runActionPlan(ctx, session, plan)
    if err != nil {
        m.logger.Error("Failed to run action plan", err, map[string]interface{}{
            "plan": plan.Name,
        })
        return nil, fmt.Errorf("failed to run action plan: %w", err)
    }

    m.logger.Info("Action running completed", map[string]interface{}{
        "actions": len(actions),
        "results": len(results),
    })

    return results, nil
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

    // Action type (command, script, template_deploy, file_copy, service_control)
    Type string `json:"type" hcl:"type"`

    // Action configuration
    Config map[string]interface{} `json:"config,omitempty" hcl:"config,optional"`

    // Target machines
    Machines []string `json:"machines,omitempty" hcl:"machines,optional"`

    // Action dependencies
    Dependencies []string `json:"dependencies,omitempty" hcl:"dependencies,optional"`

    // Action tags
    Tags []string `json:"tags,omitempty" hcl:"tags,optional"`

    // Action metadata
    Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// ActingSession represents an action running session
type ActingSession struct {
    // Session ID
    SessionID string `json:"session_id" hcl:"session_id"`

    // Session status
    Status string `json:"status" hcl:"status"`

    // Total actions in session
    TotalActions int `json:"total_actions" hcl:"total_actions"`

    // Machine inventory
    MachineInventory []spookytypes.Machine `json:"machine_inventory,omitempty" hcl:"machine_inventory,optional"`
}

// ActingResult represents the result of an action
type ActingResult struct {
    // Action name
    ActionName string `json:"action_name" hcl:"action_name"`

    // Machine name
    MachineName string `json:"machine_name" hcl:"machine_name"`

    // Result status
    Status string `json:"status" hcl:"status"`

    // Start time
    StartTime time.Time `json:"start_time" hcl:"start_time"`

    // End time
    EndTime time.Time `json:"end_time" hcl:"end_time"`

    // Exit code
    ExitCode int `json:"exit_code" hcl:"exit_code"`

    // Error type
    ErrorType string `json:"error_type,omitempty" hcl:"error_type,optional"`

    // Error message
    Error string `json:"error,omitempty" hcl:"error,optional"`

    // Output
    Output string `json:"output,omitempty" hcl:"output,optional"`
}
```

### Action Configuration Types

```go
// CommandActionConfig represents command action configuration
type CommandActionConfig struct {
    // Command to execute
    Command string `json:"command" hcl:"command"`

    // Working directory
    WorkingDir string `json:"working_dir,omitempty" hcl:"working_dir,optional"`

    // Environment variables
    Environment map[string]string `json:"environment,omitempty" hcl:"environment,optional"`

    // Timeout in seconds
    Timeout int `json:"timeout,omitempty" hcl:"timeout,optional"`

    // User to run as
    User string `json:"user,omitempty" hcl:"user,optional"`
}

// ScriptActionConfig represents script action configuration
type ScriptActionConfig struct {
    // Script path
    Script string `json:"script" hcl:"script"`

    // Script arguments
    Args []string `json:"args,omitempty" hcl:"args,optional"`

    // Working directory
    WorkingDir string `json:"working_dir,omitempty" hcl:"working_dir,optional"`

    // Environment variables
    Environment map[string]string `json:"environment,omitempty" hcl:"environment,optional"`

    // Timeout in seconds
    Timeout int `json:"timeout,omitempty" hcl:"timeout,optional"`

    // User to run as
    User string `json:"user,omitempty" hcl:"user,optional"`
}
```

## CLI Commands

### Actions List Command

```bash
# List all actions in a project
spooky actions list ./my-project

# List actions with details
spooky actions list ./my-project --verbose

# List actions for specific machines
spooky actions list ./my-project --machine web-server

# List actions by tags
spooky actions list ./my-project --tags "environment=production"
```

### Actions Validate Command

```bash
# Validate action configuration
spooky actions validate ./my-project

# Validate with detailed output
spooky actions validate ./my-project --verbose

# Validate specific actions
spooky actions validate ./my-project --action restart-nginx
```

### Actions Run Command

```bash
# Run all actions
spooky actions run ./my-project

# Run specific actions
spooky actions run ./my-project --action restart-nginx

# Run actions with plan mode
spooky actions run ./my-project --plan

# Run actions in dry-run mode
spooky actions run ./my-project --dry-run

# Run actions in parallel
spooky actions run ./my-project --parallel 4

# Run actions with timeout
spooky actions run ./my-project --timeout 600
```

## Integration Examples

### Basic Action Definition

```hcl
# actions.hcl
actions {
  action "restart-nginx" {
    description = "Restart nginx service"
    type = "service_control"
    
    config {
      service = "nginx"
      action = "restart"
    }
    
    machines = ["web-server", "app-server"]
    tags = ["web", "production"]
  }
  
  action "update-packages" {
    description = "Update system packages"
    type = "command"
    
    config {
      command = "apt update && apt upgrade -y"
      timeout = 300
    }
    
    machines = ["web-server", "db-server"]
    tags = ["maintenance"]
  }
  
  action "deploy-config" {
    description = "Deploy configuration files"
    type = "template_deploy"
    
    config {
      template = "templates/nginx.conf.tmpl"
      destination = "/etc/nginx/nginx.conf"
      backup = true
    }
    
    machines = ["web-server"]
    dependencies = ["restart-nginx"]
  }
}
```

### Action Loading and Validation

```go
// Action loading and validation example
func loadAndValidateActions(projectPath string) error {
    ctx := context.Background()
    
    // Create action manager
    manager := spookyactions.NewManager(logger, validator, sshManager, schemaValidator)
    
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

### Action Running

```go
// Action running example
func runActions(projectPath string, machines []spookytypes.Machine) error {
    ctx := context.Background()
    
    // Create action manager
    manager := spookyactions.NewManager(logger, validator, sshManager, schemaValidator)
    
    // Load actions
    actions, err := manager.LoadActions(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load actions: %w", err)
    }
    
    // Run actions
    results, err := manager.RunActions(ctx, actions, machines)
    if err != nil {
        return fmt.Errorf("failed to run actions: %w", err)
    }
    
    // Process results
    for _, result := range results {
        if result.Error != "" {
            fmt.Printf("Action %s on %s failed: %s\n", result.ActionName, result.MachineName, result.Error)
        } else {
            fmt.Printf("Action %s on %s completed successfully\n", result.ActionName, result.MachineName)
        }
    }
    
    return nil
}
```

## Error Handling

### Action Errors

```go
// Error handling example
func handleActionError(err error) {
    if err == nil {
        return
    }
    
    // Check for specific error types
    switch {
    case strings.Contains(err.Error(), "failed to create action plan"):
        fmt.Println("Action planning failed - check action dependencies")
    case strings.Contains(err.Error(), "failed to run action plan"):
        fmt.Println("Action execution failed - check SSH connectivity")
    case strings.Contains(err.Error(), "invalid action type"):
        fmt.Println("Invalid action type - check action configuration")
    case strings.Contains(err.Error(), "action not found"):
        fmt.Println("Action not found - check action names")
    default:
        fmt.Printf("Action error: %v\n", err)
    }
}
```

### Validation Errors

```go
// Validation error handling
func handleValidationError(result *spookytypes.ValidationResult) error {
    if result.Valid {
        return nil
    }
    
    fmt.Println("Action validation failed:")
    for _, err := range result.Errors {
        fmt.Printf("  - %s\n", err.Message)
    }
    
    for _, warning := range result.Warnings {
        fmt.Printf("  Warning: %s\n", warning.Message)
    }
    
    return fmt.Errorf("action validation failed with %d errors", len(result.Errors))
}
```

## Performance Considerations

### Parallel Processing

The actions system supports parallel processing:

- Multiple actions can run concurrently on different machines
- Configurable parallel worker count
- Thread-safe action execution

### Resource Management

The actions system manages resources efficiently:

- SSH connections are pooled and reused
- Memory usage is optimized for large action sets
- Timeouts prevent hanging operations

## Troubleshooting

### Common Issues

1. **SSH Connection Failures**: Check machine connectivity and SSH configuration
2. **Authentication Errors**: Verify SSH key permissions and user access
3. **Action Dependencies**: Check action dependency configuration
4. **Permission Denied**: Check user privileges on target machines
5. **Timeout Issues**: Adjust action timeouts for slow operations

### Debug Information

The actions system provides comprehensive logging for debugging:

```go
// Enable debug logging
logger.SetLevel(spookytypes.LogLevelDebug)

// Check SSH configuration
fmt.Printf("SSH config: %+v\n", sshConfig)

// Validate action configuration
err := validateAction(action)
if err != nil {
    fmt.Printf("Action validation error: %v\n", err)
}
```

## Future Enhancements

### Planned Features

1. **Parallel Execution**: Implement parallel action execution
2. **Template Integration**: Add template rendering support
3. **Variable Integration**: Add variable resolution support
4. **Result Processing**: Improve action result handling
5. **Advanced Filtering**: Support complex action filtering

### Integration Enhancements

1. **Facts Integration**: Use facts in action execution
2. **Variables Integration**: Use variables in action configuration
3. **Templates Integration**: Use templates in action definitions
4. **Advanced Planning**: Improve action planning and dependency resolution