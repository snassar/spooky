# Actions API Reference

## Overview

The actions system in spooky provides comprehensive functionality for defining, validating, planning, and executing actions on remote machines. This document provides the complete API reference for the actions system.

**Status: Implemented** - The actions system provides comprehensive functionality for action management and execution.

## Core Components

### Actions Manager (`internal/actions/manager.go`)

The actions manager provides high-level action operations:

```go
type Manager struct {
    logger   spookylogging.Logger
    config   *spookytypes.Config
    validator *Validator
    planner  *Planner
    executor *Executor
}
```

**Implemented Operations:**
- ✅ **LoadActions**: Load actions from configuration files
- ✅ **ValidateActions**: Validate action definitions and dependencies
- ✅ **PlanActions**: Plan action execution order and dependencies
- ✅ **RunActions**: Execute actions on target machines
- ✅ **MergeActions**: Merge multiple action definitions
- ✅ **OptimizeActions**: Optimize action execution for performance

### Actions Validator (`internal/actions/validator.go`)

The actions validator handles action validation:

```go
type Validator struct {
    logger   spookylogging.Logger
    config   *spookytypes.Config
    schemas  *spookyschemas.SchemaManager
}
```

**Validation Features:**
- ✅ **Action Structure Validation**: Validate action definition structure
- ✅ **Dependency Validation**: Validate action dependencies and cycles
- ✅ **Machine Validation**: Validate target machine configurations
- ✅ **Template Validation**: Validate action templates and variables
- ✅ **Security Validation**: Validate action security and permissions

### Actions Integration (`internal/actions/integration.go`)

The actions integration provides system-wide coordination:

```go
type Integration struct {
    manager  *Manager
    logger   spookylogging.Logger
    config   *spookytypes.Config
}
```

**Integration Features:**
- ✅ **SSH Integration**: Coordinate with SSH system for remote execution
- ✅ **Facts Integration**: Use facts for action context and variables
- ✅ **Variables Integration**: Resolve variables for action execution
- ✅ **Templates Integration**: Render templates for action execution
- ✅ **Secrets Integration**: Handle encrypted secrets in actions

## API Methods

### Action Loading

#### LoadActions
```go
func (m *Manager) LoadActions(projectPath string) ([]*spookytypes.Action, error)
```
Loads actions from project configuration files.

**Parameters:**
- `projectPath`: Path to project directory

**Returns:**
- `[]*Action`: Loaded action definitions
- `error`: Error if loading fails

#### LoadActionsFromFile
```go
func (m *Manager) LoadActionsFromFile(filePath string) ([]*spookytypes.Action, error)
```
Loads actions from a specific configuration file.

**Parameters:**
- `filePath`: Path to action configuration file

**Returns:**
- `[]*Action`: Loaded action definitions
- `error`: Error if loading fails

### Action Validation

#### ValidateActions
```go
func (m *Manager) ValidateActions(actions []*spookytypes.Action) (*spookytypes.ValidationResult, error)
```
Validates action definitions and dependencies.

**Parameters:**
- `actions`: Actions to validate

**Returns:**
- `ValidationResult`: Validation results with errors and warnings
- `error`: Error if validation fails

#### ValidateAction
```go
func (v *Validator) ValidateAction(action *spookytypes.Action) error
```
Validates a single action definition.

**Parameters:**
- `action`: Action to validate

**Returns:**
- `error`: Error if validation fails

#### ValidateDependencies
```go
func (v *Validator) ValidateDependencies(actions []*spookytypes.Action) error
```
Validates action dependencies and detects cycles.

**Parameters:**
- `actions`: Actions to validate dependencies for

**Returns:**
- `error`: Error if dependency validation fails

### Action Planning

#### PlanActions
```go
func (m *Manager) PlanActions(actions []*spookytypes.Action) (*spookytypes.ActionPlan, error)
```
Plans action execution order and dependencies.

**Parameters:**
- `actions`: Actions to plan

**Returns:**
- `ActionPlan`: Execution plan with order and dependencies
- `error`: Error if planning fails

#### OptimizePlan
```go
func (p *Planner) OptimizePlan(plan *spookytypes.ActionPlan) (*spookytypes.ActionPlan, error)
```
Optimizes action execution plan for performance.

**Parameters:**
- `plan`: Action plan to optimize

**Returns:**
- `ActionPlan`: Optimized execution plan
- `error`: Error if optimization fails

### Action Execution

#### RunActions
```go
func (m *Manager) RunActions(ctx context.Context, actions []*spookytypes.Action, machines []*spookytypes.Machine) (*spookytypes.ExecutionResult, error)
```
Executes actions on target machines.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `actions`: Actions to execute
- `machines`: Target machines

**Returns:**
- `ExecutionResult`: Execution results with status and output
- `error`: Error if execution fails

#### RunAction
```go
func (e *Executor) RunAction(ctx context.Context, action *spookytypes.Action, machine *spookytypes.Machine) (*spookytypes.ActionResult, error)
```
Executes a single action on a target machine.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `action`: Action to execute
- `machine`: Target machine

**Returns:**
- `ActionResult`: Action execution result
- `error`: Error if execution fails

### Action Management

#### MergeActions
```go
func (m *Manager) MergeActions(actions []*spookytypes.Action) ([]*spookytypes.Action, error)
```
Merges multiple action definitions.

**Parameters:**
- `actions`: Actions to merge

**Returns:**
- `[]*Action`: Merged action definitions
- `error`: Error if merging fails

#### FilterActions
```go
func (m *Manager) FilterActions(actions []*spookytypes.Action, filters map[string]string) ([]*spookytypes.Action, error)
```
Filters actions based on specified criteria.

**Parameters:**
- `actions`: Actions to filter
- `filters`: Filter criteria

**Returns:**
- `[]*Action`: Filtered actions
- `error`: Error if filtering fails

## Configuration

### Actions Configuration
```go
type ActionsConfig struct {
    DefaultTimeout    time.Duration `json:"default_timeout" hcl:"default_timeout"`
    MaxParallel       int           `json:"max_parallel" hcl:"max_parallel"`
    RetryAttempts     int           `json:"retry_attempts" hcl:"retry_attempts"`
    RetryDelay        time.Duration `json:"retry_delay" hcl:"retry_delay"`
    ValidationEnabled bool          `json:"validation_enabled" hcl:"validation_enabled"`
    PlanningEnabled   bool          `json:"planning_enabled" hcl:"planning_enabled"`
    OptimizationEnabled bool        `json:"optimization_enabled" hcl:"optimization_enabled"`
}
```

### Action Definition
```go
type Action struct {
    Name        string                 `json:"name" hcl:"name"`
    Description string                 `json:"description" hcl:"description"`
    Machines    []string               `json:"machines" hcl:"machines"`
    Parallel    bool                   `json:"parallel" hcl:"parallel"`
    MaxConcurrent int                  `json:"max_concurrent" hcl:"max_concurrent"`
    Dependencies []string              `json:"dependencies" hcl:"dependencies"`
    Commands    []*Command             `json:"commands" hcl:"commands"`
    FileCopy    []*FileCopy            `json:"file_copy" hcl:"file_copy"`
    ServiceControl []*ServiceControl   `json:"service_control" hcl:"service_control"`
    TemplateDeploy []*TemplateDeploy   `json:"template_deploy" hcl:"template_deploy"`
    Variables   map[string]interface{} `json:"variables" hcl:"variables"`
    Tags        map[string]string      `json:"tags" hcl:"tags"`
}
```

### Command Definition
```go
type Command struct {
    Command       string            `json:"command" hcl:"command"`
    WorkingDir    string            `json:"working_dir" hcl:"working_dir"`
    Environment   map[string]string `json:"environment" hcl:"environment"`
    Stdin         string            `json:"stdin" hcl:"stdin"`
    Timeout       time.Duration     `json:"timeout" hcl:"timeout"`
    CaptureOutput bool              `json:"capture_output" hcl:"capture_output"`
    ContinueOnError bool            `json:"continue_on_error" hcl:"continue_on_error"`
}
```

### File Copy Definition
```go
type FileCopy struct {
    Source      string `json:"source" hcl:"source"`
    Destination string `json:"destination" hcl:"destination"`
    Direction   string `json:"direction" hcl:"direction"`
    Mode        string `json:"mode" hcl:"mode"`
    Permissions string `json:"permissions" hcl:"permissions"`
    Recursive   bool   `json:"recursive" hcl:"recursive"`
    Verify      bool   `json:"verify" hcl:"verify"`
    Exclude     []string `json:"exclude" hcl:"exclude"`
}
```

### Service Control Definition
```go
type ServiceControl struct {
    Service string        `json:"service" hcl:"service"`
    Action  string        `json:"action" hcl:"action"`
    Timeout time.Duration `json:"timeout" hcl:"timeout"`
}
```

### Template Deploy Definition
```go
type TemplateDeploy struct {
    Source      string                 `json:"source" hcl:"source"`
    Destination string                 `json:"destination" hcl:"destination"`
    Permissions string                 `json:"permissions" hcl:"permissions"`
    Variables   map[string]interface{} `json:"variables" hcl:"variables"`
}
```

## Data Types

### Action Plan
```go
type ActionPlan struct {
    ID          string           `json:"id" hcl:"id"`
    CreatedAt   time.Time        `json:"created_at" hcl:"created_at"`
    Actions     []*PlannedAction `json:"actions" hcl:"actions"`
    Dependencies map[string][]string `json:"dependencies" hcl:"dependencies"`
    Parallel    bool             `json:"parallel" hcl:"parallel"`
    MaxConcurrent int            `json:"max_concurrent" hcl:"max_concurrent"`
    EstimatedTime time.Duration  `json:"estimated_time" hcl:"estimated_time"`
}
```

### Planned Action
```go
type PlannedAction struct {
    Action      *Action `json:"action" hcl:"action"`
    Order       int     `json:"order" hcl:"order"`
    Dependencies []string `json:"dependencies" hcl:"dependencies"`
    Machines    []string `json:"machines" hcl:"machines"`
    Parallel    bool     `json:"parallel" hcl:"parallel"`
}
```

### Execution Result
```go
type ExecutionResult struct {
    ID          string         `json:"id" hcl:"id"`
    StartedAt   time.Time      `json:"started_at" hcl:"started_at"`
    CompletedAt time.Time      `json:"completed_at" hcl:"completed_at"`
    Actions     []*ActionResult `json:"actions" hcl:"actions"`
    Success     bool           `json:"success" hcl:"success"`
    Errors      []string       `json:"errors" hcl:"errors"`
    Warnings    []string       `json:"warnings" hcl:"warnings"`
    Statistics  *ExecutionStatistics `json:"statistics" hcl:"statistics"`
}
```

### Action Result
```go
type ActionResult struct {
    Action      *Action `json:"action" hcl:"action"`
    Machine     string  `json:"machine" hcl:"machine"`
    StartedAt   time.Time `json:"started_at" hcl:"started_at"`
    CompletedAt time.Time `json:"completed_at" hcl:"completed_at"`
    Success     bool    `json:"success" hcl:"success"`
    Output      string  `json:"output" hcl:"output"`
    Error       string  `json:"error" hcl:"error"`
    ExitCode    int     `json:"exit_code" hcl:"exit_code"`
    Duration    time.Duration `json:"duration" hcl:"duration"`
}
```

### Execution Statistics
```go
type ExecutionStatistics struct {
    TotalActions    int           `json:"total_actions" hcl:"total_actions"`
    SuccessfulActions int         `json:"successful_actions" hcl:"successful_actions"`
    FailedActions   int           `json:"failed_actions" hcl:"failed_actions"`
    TotalDuration   time.Duration `json:"total_duration" hcl:"total_duration"`
    AverageDuration time.Duration `json:"average_duration" hcl:"average_duration"`
    ParallelExecutions int        `json:"parallel_executions" hcl:"parallel_executions"`
}
```

## Error Handling

### Action Validation Errors
```go
type ActionValidationError struct {
    Action      string `json:"action" hcl:"action"`
    Field       string `json:"field" hcl:"field"`
    Value       string `json:"value" hcl:"value"`
    Rule        string `json:"rule" hcl:"rule"`
    Message     string `json:"message" hcl:"message"`
    Severity    string `json:"severity" hcl:"severity"`
}
```

### Action Execution Errors
```go
type ActionExecutionError struct {
    Action      string    `json:"action" hcl:"action"`
    Machine     string    `json:"machine" hcl:"machine"`
    Error       string    `json:"error" hcl:"error"`
    Timestamp   time.Time `json:"timestamp" hcl:"timestamp"`
    RetryCount  int       `json:"retry_count" hcl:"retry_count"`
    Recoverable bool      `json:"recoverable" hcl:"recoverable"`
}
```

### Dependency Errors
```go
type DependencyError struct {
    Action      string   `json:"action" hcl:"action"`
    Dependencies []string `json:"dependencies" hcl:"dependencies"`
    Error       string   `json:"error" hcl:"error"`
    Cycle       []string `json:"cycle" hcl:"cycle"`
}
```

## Integration Points

### SSH Integration
The actions system integrates with the SSH system for remote execution:

```go
// Execute command via SSH
result, err := sshClient.RunCommand(ctx, command, machine)
```

### Facts Integration
The actions system integrates with the facts system for context:

```go
// Use facts in action context
facts, err := factsManager.GetFacts(machine.Hostname)
```

### Variables Integration
The actions system integrates with the variables system for resolution:

```go
// Resolve variables for action
resolved, err := variablesManager.ResolveVariables(variables, context)
```

### Templates Integration
The actions system integrates with the templates system for rendering:

```go
// Render template for action
rendered, err := templateManager.RenderTemplate(template, data)
```

## Usage Examples

### Basic Action Definition
```go
// Create action manager
manager := NewActionsManager(config, logger)

// Load actions from project
actions, err := manager.LoadActions("./my-project")
if err != nil {
    log.Fatal(err)
}

// Validate actions
result, err := manager.ValidateActions(actions)
if err != nil {
    log.Fatal(err)
}

// Plan action execution
plan, err := manager.PlanActions(actions)
if err != nil {
    log.Fatal(err)
}

// Execute actions
executionResult, err := manager.RunActions(ctx, actions, machines)
if err != nil {
    log.Fatal(err)
}
```

### Action with Commands
```go
action := &spookytypes.Action{
    Name:        "update-system",
    Description: "Update system packages",
    Machines:    []string{"web-server", "db-server"},
    Parallel:    true,
    MaxConcurrent: 2,
    
    Commands: []*spookytypes.Command{
        {
            Command:       "apt update && apt upgrade -y",
            WorkingDir:    "/tmp",
            Environment:   map[string]string{"DEBIAN_FRONTEND": "noninteractive"},
            Timeout:       300 * time.Second,
            CaptureOutput: true,
        },
    },
}
```

### Action with File Copy
```go
action := &spookytypes.Action{
    Name:        "deploy-config",
    Description: "Deploy configuration files",
    Machines:    []string{"web-server"},
    
    FileCopy: []*spookytypes.FileCopy{
        {
            Source:      "config/nginx.conf",
            Destination: "/etc/nginx/nginx.conf",
            Permissions: "0644",
            Verify:      true,
        },
    },
}
```

### Action with Service Control
```go
action := &spookytypes.Action{
    Name:        "restart-services",
    Description: "Restart application services",
    Machines:    []string{"web-server"},
    
    ServiceControl: []*spookytypes.ServiceControl{
        {
            Service: "nginx",
            Action:  "restart",
            Timeout: 60 * time.Second,
        },
    },
}
```

### Action with Template Deployment
```go
action := &spookytypes.Action{
    Name:        "deploy-template",
    Description: "Deploy configuration from template",
    Machines:    []string{"web-server"},
    
    TemplateDeploy: []*spookytypes.TemplateDeploy{
        {
            Source:      "templates/nginx.conf.tmpl",
            Destination: "/etc/nginx/nginx.conf",
            Permissions: "0644",
            Variables: map[string]interface{}{
                "server_name": "example.com",
                "port":        80,
            },
        },
    },
}
```

### Action with Dependencies
```go
actions := []*spookytypes.Action{
    {
        Name:         "backup-database",
        Description:  "Backup database before deployment",
        Machines:     []string{"db-server"},
        Dependencies: []string{},
    },
    {
        Name:         "deploy-application",
        Description:  "Deploy application",
        Machines:     []string{"web-server"},
        Dependencies: []string{"backup-database"},
    },
    {
        Name:         "verify-deployment",
        Description:  "Verify deployment success",
        Machines:     []string{"web-server"},
        Dependencies: []string{"deploy-application"},
    },
}
```

## CLI Commands

### Available Commands
- ✅ `spooky actions run` - Execute actions on machines
- ✅ `spooky actions list` - List available actions
- ✅ `spooky actions validate` - Validate action definitions
- ✅ `spooky actions plan` - Plan action execution
- ✅ `spooky actions show` - Show action details

### Command Examples
```bash
# Run actions on all machines
spooky actions run my-project

# Run specific actions
spooky actions run my-project --action update-system

# Run actions on specific machines
spooky actions run my-project --machine web-server

# Run actions in parallel
spooky actions run my-project --parallel 4

# Dry run (plan without execution)
spooky actions run my-project --dry-run

# Validate actions
spooky actions validate my-project

# Show action plan
spooky actions plan my-project
```

## Performance Features

### Parallel Execution
```go
// Configure parallel execution
config := &spookytypes.ActionsConfig{
    MaxParallel: 4,
    DefaultTimeout: 300 * time.Second,
    RetryAttempts: 3,
    RetryDelay: 5 * time.Second,
}
```

### Action Optimization
```go
// Enable action optimization
config := &spookytypes.ActionsConfig{
    PlanningEnabled: true,
    OptimizationEnabled: true,
}
```

### Dependency Resolution
```go
// Resolve action dependencies
plan, err := manager.PlanActions(actions)
if err != nil {
    log.Fatal(err)
}

// Execute in dependency order
for _, plannedAction := range plan.Actions {
    // Execute action and wait for dependencies
}
```

## Security Features

### Action Validation
```go
// Enable comprehensive validation
config := &spookytypes.ActionsConfig{
    ValidationEnabled: true,
}
```

### Permission Checking
```go
// Check action permissions
err = manager.ValidateActionPermissions(action, user)
if err != nil {
    return fmt.Errorf("insufficient permissions: %w", err)
}
```

### Secure Execution
```go
// Execute actions securely
result, err := executor.RunActionSecurely(ctx, action, machine)
if err != nil {
    return fmt.Errorf("secure execution failed: %w", err)
}
```

## Testing

### Unit Tests
- ✅ **Manager Tests**: Test actions manager functionality
- ✅ **Validator Tests**: Test action validation logic
- ✅ **Planner Tests**: Test action planning algorithms
- ✅ **Executor Tests**: Test action execution logic

### Integration Tests
- ✅ **SSH Integration Tests**: Test SSH-based action execution
- ✅ **Facts Integration Tests**: Test facts integration
- ✅ **Variables Integration Tests**: Test variable resolution
- ✅ **Templates Integration Tests**: Test template rendering

### Test Coverage
- ✅ **Action Loading**: Test action loading from files
- ✅ **Action Validation**: Test validation scenarios
- ✅ **Action Planning**: Test planning algorithms
- ✅ **Action Execution**: Test execution scenarios
- ✅ **Error Handling**: Test error conditions and recovery

## Troubleshooting

### Common Issues

#### Validation Failures
- **Missing Required Fields**: Ensure all required fields are present
- **Invalid Dependencies**: Check for circular dependencies
- **Machine Configuration**: Verify target machine configurations
- **Template Errors**: Check template syntax and variables

#### Execution Failures
- **SSH Connection Issues**: Check SSH connectivity and authentication
- **Command Failures**: Verify command syntax and permissions
- **File Transfer Issues**: Check file paths and permissions
- **Service Control Issues**: Verify service names and permissions

#### Performance Issues
- **Slow Execution**: Enable parallel execution and optimization
- **Resource Exhaustion**: Limit concurrent executions
- **Network Issues**: Check network connectivity and timeouts

### Debug Commands
```bash
# Enable debug logging
export SPOOKY_LOG_LEVEL=debug

# Validate actions with detailed output
spooky actions validate my-project --verbose

# Plan actions with detailed output
spooky actions plan my-project --verbose

# Run actions with debug output
spooky actions run my-project --verbose
```

## Future Enhancements

### Planned Features
- **Action Templates**: Reusable action templates
- **Action Libraries**: Shared action libraries
- **Advanced Scheduling**: Cron-like scheduling for actions
- **Action Monitoring**: Real-time action monitoring and alerts

### Performance Improvements
- **Distributed Execution**: Distribute actions across multiple nodes
- **Incremental Execution**: Execute only changed actions
- **Smart Caching**: Cache action results and dependencies
- **Resource Optimization**: Optimize resource usage during execution

## Summary

The actions system in spooky provides comprehensive functionality for:

1. **Action Definition**: Flexible action definition with multiple operation types
2. **Action Validation**: Comprehensive validation and dependency checking
3. **Action Planning**: Intelligent planning and optimization
4. **Action Execution**: Secure and reliable action execution
5. **Parallel Execution**: Parallel execution for improved performance
6. **Dependency Management**: Automatic dependency resolution and ordering
7. **Integration**: Seamless integration with SSH, facts, variables, and templates
8. **CLI Support**: Complete CLI command support
9. **Error Handling**: Comprehensive error handling and recovery
10. **Testing**: Complete test coverage for all functionality

The actions system is production-ready and provides all necessary functionality for reliable action management and execution in the spooky automation platform.