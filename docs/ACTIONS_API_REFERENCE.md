# Actions System API Reference

## Overview

This document provides a technical reference for the spooky actions system APIs and implementation details. It covers core interfaces, type definitions, implementation patterns, and integration details for developers working with the actions system.

## Core Interfaces

### ActionsIntegration Interface

The `ActionsIntegration` interface provides the primary entry point for action management operations:

```go
type ActionsIntegration interface {
    // LoadActions loads actions from the specified source
    LoadActions(ctx context.Context, source string) ([]spookytypes.Action, error)
    
    // ValidateActions validates a collection of actions
    ValidateActions(ctx context.Context, actions []spookytypes.Action) (*spookytypes.ValidationResult, error)
    
    // RunActions runs the specified actions on the given machines
    RunActions(ctx context.Context, actions []spookytypes.Action, machines []spookytypes.Machine) ([]spookytypes.ActingResult, error)
}
```

### ActionValidator Interface

The `ActionValidator` interface provides action validation capabilities:

```go
type ActionValidator interface {
    // ValidateActions validates a collection of actions
    ValidateActions(ctx context.Context, actions []spookytypes.Action) (*spookytypes.ValidationResult, error)
    
    // ValidateAction validates a single action
    ValidateAction(ctx context.Context, action spookytypes.Action) (*spookytypes.ValidationResult, error)
}
```

### SSHManager Interface

The `SSHManager` interface provides SSH connection and execution capabilities:

```go
type SSHManager interface {
    // CreateSession creates an SSH session for the given machine
    CreateSession(ctx context.Context, machine spookytypes.Machine) (spookytypesssh.Session, error)
    
    // RunCommand executes a command on the given machine
    RunCommand(ctx context.Context, machine spookytypes.Machine, command string) (*spookytypesssh.CommandResult, error)
    
    // TransferFile transfers a file to the given machine
    TransferFile(ctx context.Context, machine spookytypes.Machine, source, destination string) error
    
    // Close closes all SSH connections
    Close() error
}
```

## Core Types

### Action

The `Action` type represents a single action to be executed:

```go
type Action struct {
    // Basic action information
    Name        string `hcl:"name" json:"name"`
    Type        string `hcl:"type" json:"type"`
    Description string `hcl:"description,optional" json:"description,omitempty"`
    
    // Action-specific configurations
    Command     string                    `hcl:"command,optional" json:"command,omitempty"`
    Script      string                    `hcl:"script,optional" json:"script,omitempty"`
    Template    *TemplateConfig           `hcl:"template,block" json:"template,omitempty"`
    FileCopy    *FileCopyConfig           `hcl:"file_copy,block" json:"file_copy,omitempty"`
    ServiceControl *ServiceConfig         `hcl:"service_control,block" json:"service_control,omitempty"`
    
    // Execution configuration
    Timeout     int                       `hcl:"timeout,optional" json:"timeout,omitempty"`
    Parallel    bool                      `hcl:"parallel,optional" json:"parallel,omitempty"`
    Retries     int                       `hcl:"retries,optional" json:"retries,omitempty"`
    RetryDelay  int                       `hcl:"retry_delay,optional" json:"retry_delay,omitempty"`
    Dependencies []string                 `hcl:"dependencies,optional" json:"dependencies,omitempty"`
    Environment map[string]string         `hcl:"environment,optional" json:"environment,omitempty"`
    WorkingDirectory string               `hcl:"working_directory,optional" json:"working_directory,omitempty"`
    User        string                    `hcl:"user,optional" json:"user,omitempty"`
    Sudo        bool                      `hcl:"sudo,optional" json:"sudo,omitempty"`
    DryRun      bool                      `hcl:"dry_run,optional" json:"dry_run,omitempty"`
    
    // Metadata and organization
    Category    string                    `hcl:"category,optional" json:"category,omitempty"`
    Priority    int                       `hcl:"priority,optional" json:"priority,omitempty"`
    Critical    bool                      `hcl:"critical,optional" json:"critical,omitempty"`
    Metadata    map[string]string         `hcl:"metadata,optional" json:"metadata,omitempty"`
    
    // Security and validation
    ValidateBeforeRun bool                `hcl:"validate_before_run,optional" json:"validate_before_run,omitempty"`
    AllowFailure      bool                `hcl:"allow_failure,optional" json:"allow_failure,omitempty"`
    
    // Performance and resource limits
    MaxConcurrent  int                    `hcl:"max_concurrent,optional" json:"max_concurrent,omitempty"`
    ResourceLimits *ResourceLimits        `hcl:"resource_limits,block" json:"resource_limits,omitempty"`
    
    // Targeting
    Machines     []string                 `hcl:"machines,optional" json:"machines,omitempty"`
    Tags         []string                 `hcl:"tags,optional" json:"tags,omitempty"`
}
```

### ActionExecutionContext

The `ActionExecutionContext` type provides context for action execution:

```go
type ActionExecutionContext struct {
    // Execution context
    ExecutionID string                    `json:"execution_id"`
    ProjectPath string                    `json:"project_path"`
    StartTime   time.Time                 `json:"start_time"`
    
    // Action information
    Action      *Action                   `json:"action"`
    Machine     spookytypes.Machine       `json:"machine"`
    
    // Execution state
    Status      string                    `json:"status"`
    Attempt     int                       `json:"attempt"`
    MaxAttempts int                       `json:"max_attempts"`
    
    // Environment and variables
    Environment map[string]string         `json:"environment"`
    Variables   map[string]interface{}    `json:"variables"`
    
    // SSH session
    SSHSession  spookytypesssh.Session    `json:"ssh_session,omitempty"`
    
    // Common entity fields
    spookytypescommon.TimestampedEntity
}
```

### ActingSession

The `ActingSession` type represents an action execution session:

```go
type ActingSession struct {
    // Session information
    SessionID   string                    `json:"session_id"`
    ProjectPath string                    `json:"project_path"`
    StartTime   time.Time                 `json:"start_time"`
    EndTime     *time.Time                `json:"end_time,omitempty"`
    
    // Session state
    Status      string                    `json:"status"`
    Actions     []*Action                 `json:"actions"`
    Machines    []spookytypes.Machine     `json:"machines"`
    
    // Execution results
    Results     []ActingResult            `json:"results"`
    Errors      []error                   `json:"errors,omitempty"`
    
    // Configuration
    Parallel    bool                      `json:"parallel"`
    MaxConcurrent int                     `json:"max_concurrent"`
    Timeout     time.Duration             `json:"timeout"`
    
    // Common entity fields
    spookytypescommon.TimestampedEntity
}
```

### ActingResult

The `ActingResult` type represents the result of an action execution:

```go
type ActingResult struct {
    // Result identification
    ResultID    string                    `json:"result_id"`
    ActionName  string                    `json:"action_name"`
    MachineName string                    `json:"machine_name"`
    
    // Execution details
    StartTime   time.Time                 `json:"start_time"`
    EndTime     *time.Time                `json:"end_time,omitempty"`
    Duration    time.Duration             `json:"duration"`
    
    // Execution status
    Status      string                    `json:"status"`
    ExitCode    int                       `json:"exit_code"`
    Error       error                     `json:"error,omitempty"`
    ErrorType   string                    `json:"error_type,omitempty"`
    ErrorMessage string                   `json:"error_message,omitempty"`
    
    // Output and results
    Stdout      string                    `json:"stdout,omitempty"`
    Stderr      string                    `json:"stderr,omitempty"`
    Output      map[string]interface{}    `json:"output,omitempty"`
    
    // Metadata
    Attempt     int                       `json:"attempt"`
    MaxAttempts int                       `json:"max_attempts"`
    Retries     int                       `json:"retries"`
    
    // Common entity fields
    spookytypescommon.TimestampedEntity
}
```

### ActionCollection

The `ActionCollection` type represents a collection of actions:

```go
type ActionCollection struct {
    // Collection information
    Name        string                    `json:"name"`
    Description string                    `json:"description,omitempty"`
    
    // Actions
    Actions     []*Action                 `json:"actions"`
    
    // Metadata
    Metadata    map[string]interface{}    `json:"metadata,omitempty"`
    
    // Common entity fields
    spookytypescommon.TimestampedEntity
}
```

### ActionPlan

The `ActionPlan` type represents an execution plan for actions:

```go
type ActionPlan struct {
    // Plan information
    PlanID      string                    `json:"plan_id"`
    PlanName    string                    `json:"plan_name"`
    Description string                    `json:"description,omitempty"`
    
    // Actions and execution order
    Actions     []*Action                 `json:"actions"`
    ExecutionOrder [][]string             `json:"execution_order"`
    Dependencies map[string][]string      `json:"dependencies"`
    
    // Plan configuration
    Parallel    bool                      `json:"parallel"`
    MaxConcurrent int                     `json:"max_concurrent"`
    Timeout     time.Duration             `json:"timeout"`
    
    // Validation
    Validated   bool                      `json:"validated"`
    ValidationErrors []error              `json:"validation_errors,omitempty"`
    
    // Common entity fields
    spookytypescommon.TimestampedEntity
}
```

### ActionDependency

The `ActionDependency` type represents dependencies between actions:

```go
type ActionDependency struct {
    // Dependency information
    FromAction  string                    `json:"from_action"`
    ToAction    string                    `json:"to_action"`
    Type        string                    `json:"type"` // "requires", "triggers", "conflicts"
    
    // Dependency metadata
    Description string                    `json:"description,omitempty"`
    Metadata    map[string]interface{}    `json:"metadata,omitempty"`
    
    // Common entity fields
    spookytypescommon.TimestampedEntity
}
```

### ActionExecution

The `ActionExecution` type represents action execution metrics:

```go
type ActionExecution struct {
    // Execution information
    ExecutionID string                    `json:"execution_id"`
    ActionName  string                    `json:"action_name"`
    MachineName string                    `json:"machine_name"`
    
    // Timing
    StartTime   time.Time                 `json:"start_time"`
    EndTime     *time.Time                `json:"end_time,omitempty"`
    Duration    time.Duration             `json:"duration"`
    
    // Performance metrics
    CPUUsage    float64                   `json:"cpu_usage,omitempty"`
    MemoryUsage int64                     `json:"memory_usage,omitempty"`
    DiskUsage   int64                     `json:"disk_usage,omitempty"`
    
    // Network metrics
    BytesSent   int64                     `json:"bytes_sent,omitempty"`
    BytesReceived int64                   `json:"bytes_received,omitempty"`
    
    // Common entity fields
    spookytypescommon.TimestampedEntity
}
```

### ActionValidation

The `ActionValidation` type represents action validation results:

```go
type ActionValidation struct {
    // Validation information
    ActionName  string                    `json:"action_name"`
    Valid       bool                      `json:"valid"`
    
    // Validation details
    Errors      []ValidationError         `json:"errors,omitempty"`
    Warnings    []ValidationWarning       `json:"warnings,omitempty"`
    
    // Validation metadata
    ValidatedAt time.Time                 `json:"validated_at"`
    Validator   string                    `json:"validator"`
    
    // Common entity fields
    spookytypescommon.TimestampedEntity
}
```

### ActionMetrics

The `ActionMetrics` type represents action performance metrics:

```go
type ActionMetrics struct {
    // Metrics identification
    ActionName  string                    `json:"action_name"`
    TimeRange   spookytypescommon.TimeRange `json:"time_range"`
    
    // Execution metrics
    TotalExecutions int                   `json:"total_executions"`
    SuccessfulExecutions int              `json:"successful_executions"`
    FailedExecutions int                  `json:"failed_executions"`
    SuccessRate     float64               `json:"success_rate"`
    
    // Performance metrics
    AverageDuration time.Duration         `json:"average_duration"`
    MinDuration     time.Duration         `json:"min_duration"`
    MaxDuration     time.Duration         `json:"max_duration"`
    
    // Resource metrics
    AverageCPUUsage float64               `json:"average_cpu_usage"`
    AverageMemoryUsage int64              `json:"average_memory_usage"`
    AverageDiskUsage int64                `json:"average_disk_usage"`
    
    // Common entity fields
    spookytypescommon.TimestampedEntity
}
```

## Configuration Types

### TemplateConfig

Configuration for template deployment actions:

```go
type TemplateConfig struct {
    Source      string `hcl:"source" json:"source"`
    Destination string `hcl:"destination" json:"destination"`
    Validate    bool   `hcl:"validate,optional" json:"validate,omitempty"`
    Backup      bool   `hcl:"backup,optional" json:"backup,omitempty"`
    Permissions string `hcl:"permissions,optional" json:"permissions,omitempty"`
    Owner       string `hcl:"owner,optional" json:"owner,omitempty"`
    Group       string `hcl:"group,optional" json:"group,omitempty"`
}
```

### FileCopyConfig

Configuration for file copy actions:

```go
type FileCopyConfig struct {
    Source      string `hcl:"source" json:"source"`
    Destination string `hcl:"destination" json:"destination"`
    Backup      bool   `hcl:"backup,optional" json:"backup,omitempty"`
    Permissions string `hcl:"permissions,optional" json:"permissions,omitempty"`
    Owner       string `hcl:"owner,optional" json:"owner,omitempty"`
    Group       string `hcl:"group,optional" json:"group,omitempty"`
}
```

### ServiceConfig

Configuration for service control actions:

```go
type ServiceConfig struct {
    Service       string `hcl:"service" json:"service"`
    Action        string `hcl:"action" json:"action"`
    Systemd       bool   `hcl:"systemd,optional" json:"systemd,omitempty"`
    Timeout       int    `hcl:"timeout,optional" json:"timeout,omitempty"`
    WaitForStatus string `hcl:"wait_for_status,optional" json:"wait_for_status,omitempty"`
    WaitTimeout   int    `hcl:"wait_timeout,optional" json:"wait_timeout,omitempty"`
}
```

### ResourceLimits

Configuration for resource limits:

```go
type ResourceLimits struct {
    MemoryMB   int `hcl:"memory_mb,optional" json:"memory_mb,omitempty"`
    CPUPercent int `hcl:"cpu_percent,optional" json:"cpu_percent,omitempty"`
    DiskMB     int `hcl:"disk_mb,optional" json:"disk_mb,omitempty"`
}
```

## Error Types

### ActionError

Error type for action-specific errors:

```go
type ActionError struct {
    ActionName  string                    `json:"action_name"`
    MachineName string                    `json:"machine_name"`
    Operation   string                    `json:"operation"`
    Message     string                    `json:"message"`
    Cause       error                     `json:"cause,omitempty"`
    
    spookytypescommon.ErrorDetails
}

func (e *ActionError) Error() string {
    return fmt.Sprintf("action error in %s on %s during %s: %s", 
        e.ActionName, e.MachineName, e.Operation, e.Message)
}

func (e *ActionError) Unwrap() error {
    return nil
}
```

### ActingError

Error type for action execution errors:

```go
type ActingError struct {
    SessionID   string                    `json:"session_id"`
    ActionName  string                    `json:"action_name"`
    MachineName string                    `json:"machine_name"`
    Stage       string                    `json:"stage"`
    Message     string                    `json:"message"`
    Cause       error                     `json:"cause,omitempty"`
    
    spookytypescommon.ErrorDetails
}

func (e *ActingError) Error() string {
    return fmt.Sprintf("acting error in session %s for action %s on %s during %s: %s", 
        e.SessionID, e.ActionName, e.MachineName, e.Stage, e.Message)
}

func (e *ActingError) Unwrap() error {
    return nil
}
```

### PlanningError

Error type for action planning errors:

```go
type PlanningError struct {
    PlanID      string                    `json:"plan_id"`
    ActionName  string                    `json:"action_name"`
    Stage       string                    `json:"stage"`
    Message     string                    `json:"message"`
    Cause       error                     `json:"cause,omitempty"`
    
    spookytypescommon.ErrorDetails
}

func (e *PlanningError) Error() string {
    return fmt.Sprintf("planning error in plan %s for action %s during %s: %s", 
        e.PlanID, e.ActionName, e.Stage, e.Message)
}

func (e *PlanningError) Unwrap() error {
    return nil
}
```

### ValidationError

Error type for action validation errors:

```go
type ValidationError struct {
    ActionName  string                    `json:"action_name"`
    Field       string                    `json:"field"`
    Value       interface{}               `json:"value"`
    Rule        string                    `json:"rule"`
    Message     string                    `json:"message"`
    
    spookytypescommon.ErrorDetails
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error for action %s field %s with value %v: %s (rule: %s)", 
        e.ActionName, e.Field, e.Value, e.Message, e.Rule)
}

func (e *ValidationError) Unwrap() error {
    return nil
}
```

### DependencyError

Error type for dependency-related errors:

```go
type DependencyError struct {
    ActionName  string                    `json:"action_name"`
    Dependency  string                    `json:"dependency"`
    Type        string                    `json:"type"`
    Message     string                    `json:"message"`
    Cause       error                     `json:"cause,omitempty"`
    
    spookytypescommon.ErrorDetails
}

func (e *DependencyError) Error() string {
    return fmt.Sprintf("dependency error for action %s with dependency %s (%s): %s", 
        e.ActionName, e.Dependency, e.Type, e.Message)
}

func (e *DependencyError) Unwrap() error {
    return nil
}
```

## Implementation Details

### Action Loading Process

The action loading process follows these steps:

1. **File Discovery**: Scan for `actions.hcl` and files in `actions/` directory
2. **HCL Parsing**: Parse HCL files using schema validation
3. **Action Creation**: Create Action structs from parsed data
4. **Validation**: Validate action configurations
5. **Dependency Resolution**: Resolve action dependencies
6. **Collection**: Return collection of validated actions

### Action Validation Process

The action validation process validates:

- **Required Fields**: All required fields are present
- **Action Type**: Action type is valid and supported
- **Configuration**: Action-specific configuration is valid
- **Dependencies**: Dependencies reference valid actions
- **Circular Dependencies**: No circular dependency chains
- **Machine Targeting**: Target machines exist in inventory
- **Resource Limits**: Resource limits are reasonable

### Action Execution Process

The action execution process follows these steps:

1. **Session Creation**: Create acting session for the execution
2. **Plan Creation**: Create execution plan with dependency resolution
3. **Machine Targeting**: Determine target machines for each action
4. **SSH Connection**: Establish SSH connections to target machines
5. **Action Execution**: Execute actions according to the plan
6. **Result Collection**: Collect and aggregate execution results
7. **Session Cleanup**: Clean up SSH connections and resources

### Dependency Resolution

The dependency resolution process:

1. **Graph Construction**: Build dependency graph from action dependencies
2. **Cycle Detection**: Detect and report circular dependencies
3. **Topological Sort**: Create execution order using topological sort
4. **Parallel Grouping**: Group actions that can run in parallel
5. **Plan Creation**: Create execution plan with proper ordering

## CLI Integration

### Action Commands

The actions system integrates with the CLI through the following commands:

```go
var actionsCmd = &cobra.Command{
    Use:   "actions",
    Short: "Manage and run actions",
    Long:  "Manage and run actions on remote machines",
}

var actionsListCmd = &cobra.Command{
    Use:   "list [project-path]",
    Short: "List actions in a project",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation for listing actions
        return nil
    },
}

var actionsValidateCmd = &cobra.Command{
    Use:   "validate [project-path]",
    Short: "Validate action configuration",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation for validating actions
        return nil
    },
}

var actionsRunCmd = &cobra.Command{
    Use:   "run [project-path]",
    Short: "Run actions on machines",
    Long:  "Run actions on machines with support for --plan and --dry-run modes",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation for running actions
        return nil
    },
}
```

### Command Flags

```go
func init() {
    // Common flags
    actionsListCmd.Flags().String("machine", "", "Target specific machine")
    actionsListCmd.Flags().String("tags", "", "Filter by tags")
    actionsListCmd.Flags().Bool("verbose", false, "Verbose output")
    
    actionsValidateCmd.Flags().String("action", "", "Validate specific action")
    actionsValidateCmd.Flags().Bool("verbose", false, "Verbose output")
    
    actionsRunCmd.Flags().String("action", "", "Run specific action")
    actionsRunCmd.Flags().String("machine", "", "Target specific machine")
    actionsRunCmd.Flags().String("tags", "", "Filter by tags")
    actionsRunCmd.Flags().Int("parallel", 1, "Number of parallel workers")
    actionsRunCmd.Flags().Int("timeout", 300, "Action timeout in seconds")
    actionsRunCmd.Flags().Bool("plan", false, "Show execution plan without running")
    actionsRunCmd.Flags().Bool("dry-run", false, "Simulate execution without running")
    actionsRunCmd.Flags().Bool("verbose", false, "Verbose output")
}
```

## Performance Considerations

### Parallel Execution

The actions system supports parallel execution:

- **Action-Level Parallelism**: Actions can be configured to run in parallel
- **Machine-Level Parallelism**: Actions can run on multiple machines simultaneously
- **Concurrency Control**: Configurable limits on concurrent executions
- **Resource Management**: Proper resource cleanup for parallel executions

### Resource Management

The actions system manages resources efficiently:

- **SSH Connection Pooling**: Reuse SSH connections when possible
- **Memory Management**: Efficient memory allocation and cleanup
- **File Handling**: Proper file handling and cleanup
- **Timeout Management**: Configurable timeouts for all operations

### Caching

The actions system implements caching for:

- **Action Validation**: Cache validation results
- **Dependency Resolution**: Cache dependency graphs
- **Machine Targeting**: Cache machine targeting results
- **SSH Connections**: Cache SSH connection information

## Security Considerations

### SSH Security

The actions system ensures SSH security:

- **Authentication**: Secure SSH authentication methods
- **Key Management**: Proper SSH key management
- **Connection Encryption**: Encrypted SSH connections
- **Host Verification**: SSH host key verification

### Action Security

The actions system implements action security:

- **Command Validation**: Validate commands for security
- **File Permissions**: Proper file permission handling
- **User Context**: Execute actions in appropriate user context
- **Sudo Usage**: Secure sudo usage when required

### Data Protection

The actions system protects sensitive data:

- **Environment Variables**: Secure handling of environment variables
- **File Content**: Secure handling of file content
- **Logging**: Secure logging without sensitive data exposure
- **Error Messages**: Secure error messages without sensitive data

## Testing

### Unit Testing

The actions system includes comprehensive unit tests:

- **Interface Testing**: Test all interface implementations
- **Mock Testing**: Mock dependencies for isolated testing
- **Error Testing**: Test error conditions and edge cases
- **Validation Testing**: Test validation rules and schemas

### Integration Testing

Integration tests cover:

- **End-to-End Workflows**: Complete action execution workflows
- **SSH Integration**: SSH connection and execution testing
- **CLI Integration**: Command-line interface testing
- **Error Scenarios**: Real-world error scenarios

### Performance Testing

Performance tests validate:

- **Execution Performance**: Action execution speed and efficiency
- **Parallel Performance**: Parallel execution performance
- **Resource Usage**: Resource consumption patterns
- **Scalability**: Performance with large numbers of actions

## Best Practices

### Action Design

- **Single Responsibility**: Each action should have a single, clear purpose
- **Descriptive Names**: Use clear, descriptive action names
- **Proper Documentation**: Include descriptions for all actions
- **Error Handling**: Implement proper error handling in actions
- **Resource Cleanup**: Ensure proper resource cleanup

### Configuration Management

- **Schema Validation**: Always validate action configurations
- **Template Usage**: Use templates for dynamic configuration
- **Variable Substitution**: Use variables for configuration values
- **Environment Separation**: Separate configuration by environment
- **Version Control**: Version control all action configurations

### Error Handling

- **Graceful Degradation**: Handle errors gracefully
- **Retry Logic**: Implement appropriate retry logic
- **Error Reporting**: Provide clear error reporting
- **Recovery Procedures**: Implement recovery procedures
- **Monitoring**: Monitor action execution for errors

### Performance Optimization

- **Parallel Execution**: Use parallel execution when appropriate
- **Resource Limits**: Set appropriate resource limits
- **Caching**: Use caching for performance optimization
- **Connection Pooling**: Use connection pooling for SSH
- **Timeout Management**: Set appropriate timeouts

## Conclusion

The spooky actions system provides a comprehensive, interface-based solution for action orchestration across all spooky components. The system is designed to be extensible, performant, and easy to integrate with existing code.

Key features include:
- Interface-based design for loose coupling
- Multiple action types and configurations
- Comprehensive dependency management
- Parallel execution capabilities
- Extensive validation and error handling
- Easy integration with CLI and other components

For usage examples and best practices, refer to the [User Guide](ACTIONS_USER_GUIDE.md) and [Troubleshooting Guide](ACTIONS_TROUBLESHOOTING.md).
