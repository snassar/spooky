# Actions System API Reference

## Overview

This document provides a technical reference for the spooky actions system APIs and implementation details. It covers core interfaces, type definitions, implementation patterns, and integration details for developers working with the actions system.

**Status: Production Ready** - The actions system is fully implemented with complete acting infrastructure, SSH integration, and comprehensive error handling.

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

**Implementation Status**: ✅ **Fully Implemented** - Complete HCL parsing, validation, and SSH-based execution

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

**Implementation Status**: ✅ **Fully Implemented** - Complete validation with type checking, dependency validation, and detailed error reporting

### SSHManager Interface

The `SSHManager` interface provides SSH connection and running capabilities:

```go
type SSHManager interface {
    // CreateSession creates an SSH session for the given machine
    CreateSession(ctx context.Context, machine spookytypes.Machine) (spookytypesssh.Session, error)
    
    // RunCommand runs a command on the given machine
    RunCommand(ctx context.Context, machine spookytypes.Machine, command string) (*spookytypesssh.CommandResult, error)
    
    // TransferFile transfers a file to the given machine
    TransferFile(ctx context.Context, machine spookytypes.Machine, source, destination string) error
    
    // Close closes all SSH connections
    Close() error
}
```

## Core Types

### Action

The `Action` type represents a single action to be run:

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
    
    // Run configuration
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

### ActionRunContext

The `ActionRunContext` type provides context for action running:

```go
type ActionRunContext struct {
    // Run context
    RunID string                          `json:"run_id"`
    ProjectPath string                    `json:"project_path"`
    StartTime   time.Time                 `json:"start_time"`
    
    // Action information
    Action      *Action                   `json:"action"`
    Machine     spookytypes.Machine       `json:"machine"`
    
    // Run state
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

The `ActingSession` type represents an action running session:

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
    
    // Run results
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

The `ActingResult` type represents the result of an action run:

```go
type ActingResult struct {
    // Result identification
    ResultID    string                    `json:"result_id"`
    ActionName  string                    `json:"action_name"`
    MachineName string                    `json:"machine_name"`
    
    // Run details
    StartTime   time.Time                 `json:"start_time"`
    EndTime     *time.Time                `json:"end_time,omitempty"`
    Duration    time.Duration             `json:"duration"`
    
    // Run status
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

The `ActionPlan` type represents a run plan for actions:

```go
type ActionPlan struct {
    // Plan information
    PlanID      string                    `json:"plan_id"`
    PlanName    string                    `json:"plan_name"`
    Description string                    `json:"description,omitempty"`
    
    // Actions and run order
    Actions     []*Action                 `json:"actions"`
    RunOrder [][]string                   `json:"run_order"`
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

### ActionRun

The `ActionRun` type represents action run metrics:

```go
type ActionRun struct {
    // Run information
    RunID string                          `json:"run_id"`
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
    
    // Run metrics
    TotalRuns int                         `json:"total_runs"`
    SuccessfulRuns int                    `json:"successful_runs"`
    FailedRuns int                        `json:"failed_runs"`
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