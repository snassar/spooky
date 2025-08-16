# Actions System API Reference

## Overview

This document provides a comprehensive API reference for the spooky actions system. It covers all interfaces, types, methods, and implementation details for developers working with the actions system.

**Status: Implemented** - The actions system provides comprehensive functionality for action management, orchestration, and execution.

## Core Interfaces

### ActionsIntegration Interface

The `ActionsIntegration` interface provides the primary entry point for actions operations:

```go
type ActionsIntegration interface {
    // LoadActions loads actions from the given project path
    LoadActions(ctx context.Context, projectPath string) (interface{}, error)
    
    // ValidateActions validates actions
    ValidateActions(ctx context.Context, actions interface{}) (*spookytypes.ValidationResult, error)
    
    // PlanActions plans actions
    PlanActions(ctx context.Context, actions interface{}, machines interface{}) (interface{}, error)
    
    // RunActions runs actions
    RunActions(ctx context.Context, actions interface{}, machines interface{}, secretsIntegration SecretsIntegration, identityPath string) error
    
    // MergeActions merges actions
    MergeActions(ctx context.Context, actions interface{}) (interface{}, error)
    
    // OptimizeActions optimizes actions
    OptimizeActions(ctx context.Context, actions interface{}) (interface{}, error)
    
    // GetManager returns the underlying action manager
    GetManager() interface{}
}
```

**Implementation Status**: ✅ **Implemented** - Complete functionality for action management and orchestration

### ActionManager Interface

The `ActionManager` interface provides action management and orchestration:

```go
type ActionManager interface {
    // LoadActions loads actions from the given project path
    LoadActions(ctx context.Context, projectPath string) ([]*spookytypes.Action, error)
    
    // ValidateActions validates actions
    ValidateActions(ctx context.Context, actions []*spookytypes.Action) (*spookytypes.ValidationResult, error)
    
    // PlanActions plans actions
    PlanActions(ctx context.Context, actions []*spookytypes.Action, machines []*spookytypes.Machine) (*spookytypes.ActionPlan, error)
    
    // RunActions runs actions
    RunActions(ctx context.Context, actions []*spookytypes.Action, machines []*spookytypes.Machine, secretsIntegration SecretsIntegration, identityPath string) error
    
    // MergeActions merges actions
    MergeActions(ctx context.Context, actions []*spookytypes.Action) ([]*spookytypes.Action, error)
    
    // OptimizeActions optimizes actions
    OptimizeActions(ctx context.Context, actions []*spookytypes.Action) ([]*spookytypes.Action, error)
}
```

**Implementation Status**: ✅ **Implemented** - Complete functionality for action management

## Core Types

### Action

```go
type Action struct {
    Name        string                 `hcl:"name" json:"name"`
    Description string                 `hcl:"description,optional" json:"description,omitempty"`
    Machines    []string               `hcl:"machines,optional" json:"machines,omitempty"`
    Tags        map[string]string      `hcl:"tags,optional" json:"tags,omitempty"`
    Parallel    bool                   `hcl:"parallel,optional" json:"parallel,omitempty"`
    Template    *ActionTemplate        `hcl:"template,block" json:"template,omitempty"`
    Command     string                 `hcl:"command,optional" json:"command,omitempty"`
    Script      string                 `hcl:"script,optional" json:"script,omitempty"`
    Variables   map[string]interface{} `hcl:"variables,optional" json:"variables,omitempty"`
    Metadata    *ActionMetadata        `hcl:"metadata,block" json:"metadata,omitempty"`
}

type ActionTemplate struct {
    Source      string `hcl:"source" json:"source"`
    Destination string `hcl:"destination" json:"destination"`
    Permissions string `hcl:"permissions,optional" json:"permissions,omitempty"`
}

type ActionMetadata struct {
    Version   string            `hcl:"version" json:"version"`
    Author    string            `hcl:"author,optional" json:"author,omitempty"`
    Tags      map[string]string `hcl:"tags,optional" json:"tags,omitempty"`
    CreatedAt time.Time         `hcl:"created_at" json:"created_at"`
    UpdatedAt time.Time         `hcl:"updated_at" json:"updated_at"`
}
```

### ActionPlan

```go
type ActionPlan struct {
    Actions    []*PlannedAction `hcl:"actions" json:"actions"`
    Machines   []string         `hcl:"machines" json:"machines"`
    Parallel   bool             `hcl:"parallel" json:"parallel"`
    Estimated  time.Duration    `hcl:"estimated" json:"estimated"`
    CreatedAt  time.Time        `hcl:"created_at" json:"created_at"`
}

type PlannedAction struct {
    Action     *spookytypes.Action `hcl:"action" json:"action"`
    Machine    string              `hcl:"machine" json:"machine"`
    Order      int                 `hcl:"order" json:"order"`
    Dependencies []string          `hcl:"dependencies,optional" json:"dependencies,omitempty"`
}
```

### Action Context

```go
type ActionRunContext struct {
    ProjectPath string                 `hcl:"project_path" json:"project_path"`
    Action      *spookytypes.Action    `hcl:"action" json:"action"`
    Machine     *spookytypes.Machine   `hcl:"machine" json:"machine"`
    Variables   map[string]interface{} `hcl:"variables" json:"variables"`
    Facts       *spookytypes.FactCollection `hcl:"facts,optional" json:"facts,omitempty"`
    DryRun      bool                   `hcl:"dry_run" json:"dry_run"`
    Parallel    bool                   `hcl:"parallel" json:"parallel"`
    Timeout     time.Duration          `hcl:"timeout" json:"timeout"`
}
```

## Current Implementation Status

### ✅ Working Components

1. **Action Loading**: Loading actions from HCL configuration files
2. **Action Validation**: Comprehensive validation of action structures
3. **Action Planning**: Planning actions with dependency resolution
4. **Action Orchestration**: Running actions on target machines
5. **CLI Integration**: `spooky actions run` command with all options
6. **Project Integration**: Actions loading from project configuration
7. **SSH Integration**: SSH-based action execution with authentication
8. **Template Support**: Template rendering for action scripts
9. **Variable Resolution**: Variable resolution in action contexts
10. **Age Encryption**: Support for age-encrypted action values
11. **Parallel Execution**: Parallel action execution support
12. **Dry Run Support**: Dry run mode for action planning
13. **Error Handling**: Comprehensive error handling and validation

### 🔧 Key Features

1. **SSH-Based Execution**: SSH-based action execution on remote machines
2. **Template Rendering**: Template rendering for dynamic action scripts
3. **Variable Resolution**: Variable resolution in action contexts
4. **Age Encryption**: Support for age-encrypted sensitive values
5. **Parallel Execution**: Parallel action execution for improved performance
6. **Dependency Resolution**: Action dependency resolution and ordering
7. **Dry Run Mode**: Dry run mode for action planning and validation

## Implementation Details

### Action Loading System

The actions system loads actions from HCL configuration files:

```go
// Load actions from project path
func (m *Manager) LoadActions(ctx context.Context, projectPath string) ([]*spookytypes.Action, error) {
    // Load actions.hcl file
    actionsPath := filepath.Join(projectPath, "actions.hcl")
    
    data, err := os.ReadFile(actionsPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read actions file: %w", err)
    }
    
    // Parse HCL
    var config struct {
        Actions []*spookytypes.Action `hcl:"action,block"`
    }
    
    if err := hcl.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse actions file: %w", err)
    }
    
    // Validate actions
    for _, action := range config.Actions {
        if err := m.validateAction(action); err != nil {
            return nil, fmt.Errorf("invalid action %s: %w", action.Name, err)
        }
    }
    
    return config.Actions, nil
}
```

### Action Validation System

```go
// Validate action
func (m *Manager) ValidateActions(ctx context.Context, actions []*spookytypes.Action) (*spookytypes.ValidationResult, error) {
    var errors []spookytypesschemas.SchemaError
    var warnings []spookytypesschemas.SchemaError
    
    for i, action := range actions {
        // Validate required fields
        if action.Name == "" {
            errors = append(errors, spookytypesschemas.SchemaError{
                Field:   fmt.Sprintf("actions[%d].name", i),
                Message: "action name is required",
            })
        }
        
        // Validate action has either command, script, or template
        if action.Command == "" && action.Script == "" && action.Template == nil {
            errors = append(errors, spookytypesschemas.SchemaError{
                Field:   fmt.Sprintf("actions[%d]", i),
                Message: "action must have command, script, or template",
            })
        }
        
        // Validate template if present
        if action.Template != nil {
            if action.Template.Source == "" {
                errors = append(errors, spookytypesschemas.SchemaError{
                    Field:   fmt.Sprintf("actions[%d].template.source", i),
                    Message: "template source is required",
                })
            }
            
            if action.Template.Destination == "" {
                errors = append(errors, spookytypesschemas.SchemaError{
                    Field:   fmt.Sprintf("actions[%d].template.destination", i),
                    Message: "template destination is required",
                })
            }
        }
    }
    
    valid := len(errors) == 0
    
    return &spookytypes.ValidationResult{
        Valid:    valid,
        Errors:   errors,
        Warnings: warnings,
    }, nil
}
```

### Action Planning System

```go
// Plan actions with dependency resolution
func (m *Manager) PlanActions(ctx context.Context, actions []*spookytypes.Action, machines []*spookytypes.Machine) (*spookytypes.ActionPlan, error) {
    plan := &spookytypes.ActionPlan{
        Actions:   make([]*spookytypes.PlannedAction, 0),
        Machines:  make([]string, 0),
        Parallel:  false,
        CreatedAt: time.Now(),
    }
    
    // Add machines to plan
    for _, machine := range machines {
        plan.Machines = append(plan.Machines, machine.Hostname)
    }
    
    // Plan each action for each machine
    order := 0
    for _, action := range actions {
        // Determine target machines
        targetMachines := m.getTargetMachines(action, machines)
        
        for _, machine := range targetMachines {
            plannedAction := &spookytypes.PlannedAction{
                Action:  action,
                Machine: machine.Hostname,
                Order:   order,
            }
            
            plan.Actions = append(plan.Actions, plannedAction)
            order++
        }
        
        // Check if any action is parallel
        if action.Parallel {
            plan.Parallel = true
        }
    }
    
    // Estimate execution time
    plan.Estimated = m.estimateExecutionTime(plan.Actions)
    
    return plan, nil
}
```

### Action Execution System

```go
// Run actions on target machines
func (m *Manager) RunActions(ctx context.Context, actions []*spookytypes.Action, machines []*spookytypes.Machine, secretsIntegration SecretsIntegration, identityPath string) error {
    // Create action plan
    plan, err := m.PlanActions(ctx, actions, machines)
    if err != nil {
        return fmt.Errorf("failed to plan actions: %w", err)
    }
    
    // Execute actions
    if plan.Parallel {
        return m.runActionsParallel(ctx, plan, secretsIntegration, identityPath)
    } else {
        return m.runActionsSequential(ctx, plan, secretsIntegration, identityPath)
    }
}

func (m *Manager) runActionsSequential(ctx context.Context, plan *spookytypes.ActionPlan, secretsIntegration SecretsIntegration, identityPath string) error {
    for _, plannedAction := range plan.Actions {
        // Create action context
        actionCtx := &spookytypes.ActionRunContext{
            Action:    plannedAction.Action,
            Machine:   m.findMachine(plannedAction.Machine, plan.Machines),
            Variables: plannedAction.Action.Variables,
            DryRun:    false,
            Timeout:   30 * time.Second,
        }
        
        // Execute action
        if err := m.executeAction(ctx, actionCtx, secretsIntegration, identityPath); err != nil {
            return fmt.Errorf("failed to execute action %s on %s: %w", plannedAction.Action.Name, plannedAction.Machine, err)
        }
    }
    
    return nil
}

func (m *Manager) runActionsParallel(ctx context.Context, plan *spookytypes.ActionPlan, secretsIntegration SecretsIntegration, identityPath string) error {
    var wg sync.WaitGroup
    errors := make(chan error, len(plan.Actions))
    
    for _, plannedAction := range plan.Actions {
        wg.Add(1)
        go func(action *spookytypes.PlannedAction) {
            defer wg.Done()
            
            // Create action context
            actionCtx := &spookytypes.ActionRunContext{
                Action:    action.Action,
                Machine:   m.findMachine(action.Machine, plan.Machines),
                Variables: action.Action.Variables,
                DryRun:    false,
                Timeout:   30 * time.Second,
            }
            
            // Execute action
            if err := m.executeAction(ctx, actionCtx, secretsIntegration, identityPath); err != nil {
                errors <- fmt.Errorf("failed to execute action %s on %s: %w", action.Action.Name, action.Machine, err)
            }
        }(plannedAction)
    }
    
    wg.Wait()
    close(errors)
    
    // Check for errors
    for err := range errors {
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

### Action Execution with SSH

```go
// Execute action on target machine
func (m *Manager) executeAction(ctx context.Context, actionCtx *spookytypes.ActionRunContext, secretsIntegration SecretsIntegration, identityPath string) error {
    // Create SSH connection
    connection, err := m.sshManager.Connect(ctx, &spookytypes.ConnectionRequest{
        Host:     actionCtx.Machine.Host,
        Port:     actionCtx.Machine.Port,
        User:     actionCtx.Machine.User,
        KeyPath:  actionCtx.Machine.KeyFile,
        Password: actionCtx.Machine.Password,
        Timeout:  actionCtx.Timeout,
    })
    if err != nil {
        return fmt.Errorf("failed to connect to %s: %w", actionCtx.Machine.Hostname, err)
    }
    
    // Handle different action types
    if actionCtx.Action.Template != nil {
        return m.executeTemplateAction(ctx, connection, actionCtx, secretsIntegration, identityPath)
    } else if actionCtx.Action.Script != "" {
        return m.executeScriptAction(ctx, connection, actionCtx, secretsIntegration, identityPath)
    } else if actionCtx.Action.Command != "" {
        return m.executeCommandAction(ctx, connection, actionCtx, secretsIntegration, identityPath)
    }
    
    return fmt.Errorf("no valid action type found")
}

func (m *Manager) executeTemplateAction(ctx context.Context, connection *spookytypes.Connection, actionCtx *spookytypes.ActionRunContext, secretsIntegration SecretsIntegration, identityPath string) error {
    // Render template
    templatePath := filepath.Join(actionCtx.ProjectPath, actionCtx.Action.Template.Source)
    
    rendered, err := m.renderTemplate(templatePath, actionCtx.Variables)
    if err != nil {
        return fmt.Errorf("failed to render template: %w", err)
    }
    
    // Decrypt if needed
    if secretsIntegration != nil {
        if err := m.decryptContent(rendered, secretsIntegration, identityPath); err != nil {
            return fmt.Errorf("failed to decrypt template content: %w", err)
        }
    }
    
    // Upload template to remote machine
    remotePath := actionCtx.Action.Template.Destination
    
    if err := m.sshManager.UploadFile(ctx, connection, []byte(rendered), remotePath); err != nil {
        return fmt.Errorf("failed to upload template: %w", err)
    }
    
    // Set permissions if specified
    if actionCtx.Action.Template.Permissions != "" {
        chmodCmd := &spookytypes.SSHCommand{
            Command: fmt.Sprintf("chmod %s %s", actionCtx.Action.Template.Permissions, remotePath),
        }
        
        result, err := m.sshManager.RunCommand(ctx, connection, chmodCmd)
        if err != nil {
            return fmt.Errorf("failed to set template permissions: %w", err)
        }
        
        if !result.Success {
            return fmt.Errorf("failed to set template permissions: %s", result.Error)
        }
    }
    
    return nil
}

func (m *Manager) executeScriptAction(ctx context.Context, connection *spookytypes.Connection, actionCtx *spookytypes.ActionRunContext, secretsIntegration SecretsIntegration, identityPath string) error {
    // Resolve script content
    script := actionCtx.Action.Script
    
    // Decrypt if needed
    if secretsIntegration != nil {
        if err := m.decryptContent(script, secretsIntegration, identityPath); err != nil {
            return fmt.Errorf("failed to decrypt script content: %w", err)
        }
    }
    
    // Execute script
    cmd := &spookytypes.SSHCommand{
        Command: script,
    }
    
    result, err := m.sshManager.RunCommand(ctx, connection, cmd)
    if err != nil {
        return fmt.Errorf("failed to execute script: %w", err)
    }
    
    if !result.Success {
        return fmt.Errorf("script execution failed: %s", result.Error)
    }
    
    return nil
}

func (m *Manager) executeCommandAction(ctx context.Context, connection *spookytypes.Connection, actionCtx *spookytypes.ActionRunContext, secretsIntegration SecretsIntegration, identityPath string) error {
    // Execute command
    cmd := &spookytypes.SSHCommand{
        Command: actionCtx.Action.Command,
    }
    
    result, err := m.sshManager.RunCommand(ctx, connection, cmd)
    if err != nil {
        return fmt.Errorf("failed to execute command: %w", err)
    }
    
    if !result.Success {
        return fmt.Errorf("command execution failed: %s", result.Error)
    }
    
    return nil
}
```

## Usage Examples

### Basic Action Loading

```go
// Create actions integration
actionsIntegration := NewActionsIntegration(manager)

// Load actions from project
actions, err := actionsIntegration.LoadActions(ctx, "./my-project")
if err != nil {
    return fmt.Errorf("failed to load actions: %w", err)
}

// Validate actions
result, err := actionsIntegration.ValidateActions(ctx, actions)
if err != nil {
    return fmt.Errorf("failed to validate actions: %w", err)
}

if !result.Valid {
    for _, error := range result.Errors {
        log.Printf("Validation error: %s", error.Message)
    }
    return fmt.Errorf("actions validation failed")
}
```

### Action Planning

```go
// Plan actions
plan, err := actionsIntegration.PlanActions(ctx, actions, machines)
if err != nil {
    return fmt.Errorf("failed to plan actions: %w", err)
}

log.Printf("Planned %d actions for %d machines", len(plan.Actions), len(plan.Machines))
log.Printf("Estimated execution time: %v", plan.Estimated)
log.Printf("Parallel execution: %v", plan.Parallel)
```

### Action Execution

```go
// Run actions
err = actionsIntegration.RunActions(ctx, actions, machines, secretsIntegration, "~/.age/identity.txt")
if err != nil {
    return fmt.Errorf("failed to run actions: %w", err)
}

log.Printf("Successfully executed all actions")
```

### CLI Usage

```bash
# Run all actions in a project
spooky actions run ./my-project

# Run actions on specific machines
spooky actions run ./my-project --machine web-server --machine db-server

# Run actions with specific tags
spooky actions run ./my-project --tags environment=production

# Run actions in parallel
spooky actions run ./my-project --parallel

# Dry run to see what would be executed
spooky actions run ./my-project --dry-run

# Run actions with age decryption
spooky actions run ./my-project --decrypt --identity ~/.age/identity.txt

# Run actions with custom command
spooky actions run ./my-project --command "systemctl restart nginx"
```

## Error Handling

### Action Loading Errors

```go
// Handle action loading errors
actions, err := actionsIntegration.LoadActions(ctx, projectPath)
if err != nil {
    if os.IsNotExist(err) {
        return fmt.Errorf("actions file not found: %s/actions.hcl", projectPath)
    }
    
    if strings.Contains(err.Error(), "failed to parse") {
        return fmt.Errorf("invalid actions file format: %w", err)
    }
    
    return fmt.Errorf("failed to load actions: %w", err)
}
```

### Action Validation Errors

```go
// Handle validation errors
result, err := actionsIntegration.ValidateActions(ctx, actions)
if err != nil {
    return fmt.Errorf("validation failed: %w", err)
}

if !result.Valid {
    // Log all validation errors
    for _, error := range result.Errors {
        log.Printf("Validation error in %s: %s", error.Field, error.Message)
    }
    
    // Log all validation warnings
    for _, warning := range result.Warnings {
        log.Printf("Validation warning in %s: %s", warning.Field, warning.Message)
    }
    
    return fmt.Errorf("actions validation failed with %d errors", len(result.Errors))
}
```

### Action Execution Errors

```go
// Handle execution errors
err = actionsIntegration.RunActions(ctx, actions, machines, secretsIntegration, identityPath)
if err != nil {
    if strings.Contains(err.Error(), "connection refused") {
        return fmt.Errorf("target machine is unreachable: %w", err)
    }
    
    if strings.Contains(err.Error(), "authentication failed") {
        return fmt.Errorf("SSH authentication failed: %w", err)
    }
    
    if strings.Contains(err.Error(), "permission denied") {
        return fmt.Errorf("insufficient permissions to execute action: %w", err)
    }
    
    return fmt.Errorf("action execution failed: %w", err)
}
```

## Testing

### Action Loading Testing

```go
func TestActionLoading(t *testing.T) {
    // Create actions manager
    manager := NewManager(logger, validator, sshManager, schemaValidator)
    
    // Test project path
    projectPath := "./testdata/test-project"
    
    // Load actions
    actions, err := manager.LoadActions(ctx, projectPath)
    if err != nil {
        t.Fatalf("Failed to load actions: %v", err)
    }
    
    // Validate actions
    if len(actions) == 0 {
        t.Error("Expected actions, got empty slice")
    }
    
    // Check first action
    action := actions[0]
    if action.Name == "" {
        t.Error("Expected action name, got empty string")
    }
}
```

### Action Validation Testing

```go
func TestActionValidation(t *testing.T) {
    // Create actions manager
    manager := NewManager(logger, validator, sshManager, schemaValidator)
    
    // Test actions
    actions := []*spookytypes.Action{
        {
            Name:    "test-action",
            Command: "echo 'hello world'",
        },
        {
            Name: "", // Invalid: missing name
        },
    }
    
    // Validate actions
    result, err := manager.ValidateActions(ctx, actions)
    if err != nil {
        t.Fatalf("Failed to validate actions: %v", err)
    }
    
    if result.Valid {
        t.Error("Expected validation to fail, got valid")
    }
    
    if len(result.Errors) == 0 {
        t.Error("Expected validation errors, got none")
    }
}
```

### Mock SSH Manager

```go
type MockSSHManager struct {
    // Mock implementation for testing
}

func (m *MockSSHManager) Connect(ctx context.Context, request *spookytypes.ConnectionRequest) (*spookytypes.ConnectionResult, error) {
    return &spookytypes.ConnectionResult{
        Success: true,
        Connection: &spookytypes.Connection{
            Host: request.Host,
            Port: request.Port,
            User: request.User,
        },
    }, nil
}

func (m *MockSSHManager) RunCommand(ctx context.Context, connection *spookytypes.Connection, command *spookytypes.SSHCommand) (*spookytypes.SSHCommandResult, error) {
    return &spookytypes.SSHCommandResult{
        Success: true,
        Output:  "command executed successfully",
        ExitCode: 0,
    }, nil
}

func (m *MockSSHManager) UploadFile(ctx context.Context, connection *spookytypes.Connection, data []byte, remotePath string) error {
    return nil
}
```

## Best Practices

### Action Design

1. **Use Templates**: Prefer templates over inline scripts for complex actions
2. **Validate Actions**: Always validate actions before execution
3. **Handle Errors**: Implement proper error handling in action scripts
4. **Use Variables**: Use variables for dynamic action configuration
5. **Encrypt Sensitive Data**: Use age encryption for sensitive action values

### Performance Optimization

```go
// Run actions in parallel when possible
func runActionsOptimized(actions []*spookytypes.Action, machines []*spookytypes.Machine, actionsIntegration ActionsIntegration) error {
    // Group actions by machine for parallel execution
    machineActions := make(map[string][]*spookytypes.Action)
    
    for _, action := range actions {
        for _, machine := range action.Machines {
            machineActions[machine] = append(machineActions[machine], action)
        }
    }
    
    // Execute actions per machine in parallel
    var wg sync.WaitGroup
    errors := make(chan error, len(machineActions))
    
    for machine, machineActions := range machineActions {
        wg.Add(1)
        go func(m string, actions []*spookytypes.Action) {
            defer wg.Done()
            
            // Find machine
            targetMachine := findMachine(m, machines)
            if targetMachine == nil {
                errors <- fmt.Errorf("machine %s not found", m)
                return
            }
            
            // Run actions for this machine
            if err := actionsIntegration.RunActions(ctx, actions, []*spookytypes.Machine{targetMachine}, nil, ""); err != nil {
                errors <- fmt.Errorf("failed to run actions on %s: %w", m, err)
            }
        }(machine, machineActions)
    }
    
    wg.Wait()
    close(errors)
    
    // Check for errors
    for err := range errors {
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

## Future Enhancements

### Planned Features

1. **Action Dependencies**: Support for action dependencies and ordering
2. **Action Rollback**: Automatic rollback for failed actions
3. **Action Templates**: Reusable action templates
4. **Action Scheduling**: Scheduled action execution
5. **Action Monitoring**: Real-time action monitoring and status
6. **Action History**: Action execution history and audit logs

### Architecture Improvements

1. **Distributed Execution**: Distributed action execution across multiple orchestrators
2. **Action Queuing**: Action queuing for high-load scenarios
3. **Action Caching**: Action result caching for improved performance
4. **Action Streaming**: Streaming action execution for real-time updates
5. **Action Analytics**: Analytics and reporting for action execution

## Related Documentation

- [Actions User Guide](ACTIONS_USER_GUIDE.md) - User guide for actions system
- [Actions Troubleshooting](ACTIONS_TROUBLESHOOTING.md) - Troubleshooting guide
- [SSH API Reference](SSH_API_REFERENCE.md) - SSH system API reference
- [Templates API Reference](TEMPLATES_API_REFERENCE.md) - Templates system API reference
- [Variables API Reference](VARIABLES_API_REFERENCE.md) - Variables system API reference
- [Secrets API Reference](SECRETS_API_REFERENCE.md) - Secrets management API reference