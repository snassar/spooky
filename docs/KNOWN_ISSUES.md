# Known Issues

## Overview

This document consolidates all known issues, incomplete implementations, and current limitations in the spooky codebase with specific code paths, file references, and problematic code examples. This serves as a central reference for developers and users to understand the current state of the system and workarounds.

## Implementation Status Summary

### Production Ready Systems
- **SSH System**: Fully implemented with comprehensive functionality (`internal/ssh/`)
- **Machines System**: Complete inventory and connectivity management (`internal/machines/`)
- **Variables System**: Full variable management and resolution (`internal/variables/`)
- **Secrets System**: Complete encryption and key management (`internal/secrets/`)
- **Logging System**: Comprehensive logging and monitoring (`internal/logging/`)

### Partially Implemented Systems
- **Actions System**: Basic functionality with SSH orchestration issues (`internal/actions/`)
- **Facts System**: Basic collection with SSH-based gathering issues (`internal/facts/`)
- **Templates System**: Basic rendering with CLI command issues (`internal/templates/`)
- **Integrations System**: Basic functionality with advanced features in development (`internal/integration/`)



## SSH Integration Issues

### Primary Issue: SSH-Based Operations

**Status**: Critical - Affects multiple core systems

**Description**: While the SSH system itself is fully implemented and functional (`internal/ssh/`), several subsystems have significant issues with SSH-based operations. This creates a disconnect where SSH connectivity works but automated operations over SSH fail.

**Affected Systems**:
- Actions System (`internal/actions/`)
- Facts System (`internal/facts/`)
- Templates System (`internal/templates/`)
- Secrets System (`internal/secrets/`)

**Root Cause**: The integration layer between SSH connectivity and higher-level operations is incomplete or has implementation bugs.

**Impact**: Users cannot perform automated operations on remote machines despite having working SSH connectivity.

**Key Files Affected**:
- `internal/actions/manager.go` - SSH action orchestration issues
- `internal/facts/collector.go` - SSH fact collection issues
- `internal/templates/manager.go` - SSH template rendering issues
- `internal/ssh/manager.go` - SSH integration layer issues

### Actions System SSH Issues

**Issue**: SSH-based action orchestration has implementation issues

**Location**: `internal/actions/manager.go`

**Problematic Code Examples**:

```go
// Lines 420-558: runCommandAction - Creates new SSH connection for each action
func (m *Manager) runCommandAction(ctx context.Context, action *spookytypesactions.Action, machine *spookytypes.Machine, result *spookytypesactions.ActingResult) error {
    // Creates new SSH connection for each action - no connection reuse
    connectionRequest := &spookytypes.ConnectionRequest{
        Host:     machine.Host,
        Port:     machine.Port,
        User:     machine.User,
        Password: machine.Password,
        KeyPath:  machine.KeyFile,
        Timeout:  time.Duration(action.Timeout) * time.Second,
    }

    connectionResult, err := m.sshManager.Connect(ctx, connectionRequest)
    if err != nil {
        return fmt.Errorf("failed to connect to %s: %w", machine.Hostname, err)
    }

    // Creates new session for each action - no session reuse
    session, err := m.sshManager.CreateSession(ctx, connectionResult.Connection)
    if err != nil {
        return fmt.Errorf("failed to create session on %s: %w", machine.Hostname, err)
    }
    // ... rest of implementation
}
```

```go
// Lines 560-658: runScriptAction - Duplicates connection logic
func (m *Manager) runScriptAction(ctx context.Context, action *spookytypesactions.Action, machine *spookytypes.Machine, result *spookytypesactions.ActingResult) error {
    // Duplicates the same connection logic as runCommandAction
    connectionRequest := &spookytypes.ConnectionRequest{
        Host:     machine.Host,
        Port:     machine.Port,
        User:     machine.User,
        Password: machine.Password,
        KeyPath:  machine.KeyFile,
        Timeout:  time.Duration(action.Timeout) * time.Second,
    }
    // ... rest of implementation
}
```

**Symptoms**:
- Cannot properly run actions on remote machines
- No parallel action execution support
- Limited action planning capabilities
- No action result aggregation
- Inefficient connection management (new connection per action)
- Code duplication across action types

**Root Cause Analysis**:
1. **Connection Inefficiency**: Each action creates a new SSH connection instead of reusing connections
2. **No Connection Pooling**: No connection pooling mechanism for efficient resource usage
3. **Code Duplication**: Same connection logic repeated across different action types
4. **No Parallel Execution**: Actions run sequentially instead of in parallel
5. **Limited Error Handling**: Basic error handling without retry mechanisms

**Current Workarounds**:
- Use local action running for immediate needs
- Run actions manually on remote machines if needed
- Use machine filtering to limit running scope
- Monitor updates for improvements to SSH-based orchestration

**Technical Details**:
- SSH connectivity works (see SSH User Guide)
- Action loading and validation work
- CLI integration is functional
- The failure occurs in the action execution layer over SSH

### Facts System SSH Issues

**Issue**: SSH-based fact collection has implementation issues

**Location**: `internal/facts/collector.go`

**Problematic Code Examples**:

```go
// Lines 779-844: runSSHCommand - Creates new SSH connection for each command
func (c *SystemFactCollector) runSSHCommand(machine *spookytypes.Machine, command string) (string, error) {
    ctx := context.Background()

    // Creates new SSH connection for each fact collection command
    connectionRequest := &spookytypes.ConnectionRequest{
        Host:       machine.Hostname,
        Port:       machine.Port,
        User:       machine.User,
        Password:   machine.Password,
        KeyPath:    machine.KeyFile,
        Passphrase: machine.Passphrase,
        AuthMethod: spookytypesssh.AuthMethodPublicKey,
        Timeout:    30 * time.Second,
    }

    // Establishes new connection for each command
    connectionResult, err := c.sshManager.Connect(ctx, connectionRequest)
    if err != nil {
        return "", fmt.Errorf("failed to establish SSH connection to %s: %w", machine.Hostname, err)
    }

    // Creates new session for each command
    session, err := c.sshManager.CreateSession(ctx, connectionResult.Connection)
    if err != nil {
        return "", fmt.Errorf("failed to create SSH session on %s: %w", machine.Hostname, err)
    }
    // ... rest of implementation
}
```

```go
// Lines 212-779: collectSystemFactsViaSSH - Multiple SSH commands without connection reuse
func (c *SystemFactCollector) collectSystemFactsViaSSH(ctx context.Context, machine *spookytypes.Machine) (*spookytypesfacts.SystemFacts, error) {
    systemFacts := &spookytypesfacts.SystemFacts{}

    // Each of these calls creates a new SSH connection
    osFacts, err := c.collectOSFactsViaSSH(ctx, machine)
    if err != nil {
        return nil, fmt.Errorf("failed to collect OS facts: %w", err)
    }
    systemFacts.OS = osFacts

    hardwareFacts, err := c.collectHardwareFactsViaSSH(ctx, machine)
    if err != nil {
        return nil, fmt.Errorf("failed to collect hardware facts: %w", err)
    }
    systemFacts.Hardware = hardwareFacts

    networkFacts, err := c.collectNetworkFactsViaSSH(ctx, machine)
    if err != nil {
        return nil, fmt.Errorf("failed to collect network facts: %w", err)
    }
    systemFacts.Network = networkFacts

    return systemFacts, nil
}
```

**Symptoms**:
- Cannot reliably read `/etc/spooky/facts.*` files from remote machines
- Sequential collection only, no multi-machine parallel processing
- Cannot fully leverage existing SSH infrastructure and machine inventory
- Inefficient connection management (new connection per command)
- Multiple SSH connections for single fact collection operation

**Root Cause Analysis**:
1. **Connection Inefficiency**: Each SSH command creates a new connection instead of reusing connections
2. **No Connection Pooling**: No connection pooling mechanism for efficient resource usage
3. **Sequential Processing**: Facts are collected sequentially instead of in parallel
4. **Multiple Connections**: Single fact collection operation creates multiple SSH connections
5. **No Batch Commands**: Each fact type requires separate SSH command execution

**Current Workarounds**:
- Use local fact collection for immediate needs
- Export facts manually from remote machines if needed
- Use machine filtering to limit collection scope
- Monitor updates for improvements to SSH-based collection

**Technical Details**:
- SSH connectivity works (see SSH User Guide)
- Fact export functionality works
- CLI integration is functional
- The failure occurs in the automated fact collection layer over SSH

### Templates System SSH Issues

**Issue**: SSH-based template rendering has implementation issues

**Location**: `internal/templates/manager.go`

**Problematic Code Examples**:

```go
// Missing SSH-based template rendering implementation
// The templates manager lacks SSH integration for remote template rendering
type Manager struct {
    logger spookytypeslogging.Logger
    // Missing SSH manager integration
    // Missing remote rendering capabilities
}

// Local rendering works, but no SSH-based rendering
func (m *Manager) RenderTemplate(ctx context.Context, template *spookytypestemplates.Template, data map[string]interface{}) (string, error) {
    // Only supports local rendering
    // No SSH-based remote rendering implementation
    return m.renderLocally(template, data)
}
```

**Symptoms**:
- Cannot properly render templates on remote machines
- Limited template processing capabilities
- No template caching or optimization
- No SSH integration for remote template rendering
- Templates can only be rendered locally

**Root Cause Analysis**:
1. **Missing SSH Integration**: Templates manager lacks SSH manager integration
2. **No Remote Rendering**: No implementation for SSH-based template rendering
3. **Local-Only Rendering**: Templates can only be rendered locally, not on remote machines
4. **No Template Deployment**: No mechanism to deploy rendered templates to remote machines
5. **Limited Processing**: No advanced template processing capabilities

**Current Workarounds**:
- Use local template rendering for immediate needs
- Deploy templates through the Actions system (when working)
- Render templates manually on remote machines if needed
- Monitor updates for improvements to SSH-based template rendering

**Technical Details**:
- Template loading and validation work
- Local template rendering is functional
- The failure occurs in the remote template rendering layer over SSH

## CLI Command Issues

### Templates System CLI Issues

**Issue**: Template rendering CLI commands are partially implemented

**Location**: `cmd/templates.go`

**Problematic Code Examples**:

```go
// Lines 1-195: templatesRenderCmd - Limited functionality
var templatesRenderCmd = &cobra.Command{
    Use:   "render [project] [template]",
    Short: "Render a template",
    Long:  `Render a template with the given data and output to a file.`,
    Args:  cobra.ExactArgs(2),
    RunE: func(cmd *cobra.Command, args []string) error {
        projectPath := args[0]
        templatePath := args[1]

        // Get flags
        dataFile, _ := cmd.Flags().GetString("data")
        outputFile, _ := cmd.Flags().GetString("output")
        dryRun, _ := cmd.Flags().GetBool("dry-run")
        preview, _ := cmd.Flags().GetBool("preview")

        // Basic implementation - no SSH integration
        // No remote rendering capabilities
        // No advanced template processing
    },
}
```

```go
// Lines 195+: templatesValidateCmd - Basic validation only
var templatesValidateCmd = &cobra.Command{
    Use:   "validate [project]",
    Short: "Validate templates",
    Long:  `Validate templates in the project for syntax and security.`,
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        // Basic validation - no SSH-based validation
        // No remote template validation
        // Limited validation capabilities
    },
}
```

**Symptoms**:
- `spooky templates render` command has limited functionality
- Template management CLI commands are missing
- No CLI commands for template validation
- No SSH integration in CLI commands
- No remote template rendering via CLI

**Root Cause Analysis**:
1. **Limited CLI Implementation**: Basic CLI commands without advanced features
2. **No SSH Integration**: CLI commands lack SSH integration for remote operations
3. **Missing Commands**: Several template management commands are not implemented
4. **Basic Functionality**: Only basic template rendering and validation
5. **No Remote Support**: CLI commands only support local operations

**Current Workarounds**:
- Use templates through the Actions system for deployment
- Use local template rendering tools
- Monitor updates for CLI command improvements

**Technical Details**:
- Template engine works for local rendering
- CLI command infrastructure exists but is incomplete
- Integration with SSH-based rendering is the primary blocker

### Actions System CLI Issues

**Issue**: Limited CLI command functionality for advanced features

**Location**: `cmd/actions.go`

**Problematic Code Examples**:

```go
// Limited CLI implementation for actions
// Missing parallel execution support
// Missing action result aggregation
var actionsRunCmd = &cobra.Command{
    Use:   "run [project]",
    Short: "Run actions",
    Long:  `Run actions on target machines.`,
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        // Basic implementation - no parallel execution
        // No action result aggregation
        // Limited action planning capabilities
    },
}
```

**Symptoms**:
- No parallel execution support via CLI
- Limited action planning capabilities
- No action result aggregation
- Basic CLI functionality only
- No advanced action management features

**Root Cause Analysis**:
1. **Limited CLI Implementation**: Basic CLI commands without advanced features
2. **No Parallel Support**: CLI lacks parallel execution capabilities
3. **Missing Result Aggregation**: No mechanism to aggregate and display action results
4. **Basic Planning**: Limited action planning and dependency resolution
5. **No Advanced Features**: Missing advanced action management capabilities

**Current Workarounds**:
- Use basic action running commands
- Implement parallel execution manually if needed
- Monitor updates for CLI improvements

## CLI Logging Output Issue

**Issue**: Structured JSON logging was being output to stdout along with user-friendly messages, making the CLI output cluttered and difficult to read.

**Previous Behavior**:
```bash
$ spooky project init my-project
{"time":"2025-08-17T11:56:56.623783053+02:00","level":"INFO","msg":"IntegrationManager initialized",...}
{"time":"2025-08-17T11:56:56.623842825+02:00","level":"INFO","msg":"Initializing new project",...}
✅ Project initialized successfully: /path/to/my-project
📁 Project structure created according to project-directory.schema.hcl
```

**Current Behavior** (Fixed):
```bash
$ spooky project init my-project
✅ Project initialized successfully: /home/sn/Workshop/go/spooky/my-project
📁 Project structure created according to project-directory.schema.hcl
📄 Configuration files generated using project.schema.hcl
💡 Next steps:
   - Edit project.hcl to customize your project
   - Add machines.hcl for machine inventory
   - Add actions.hcl for automation tasks
   - Add variables.hcl for project variables
```

**Solution Implemented**:
- Updated default logging configuration to use `level = "error"` and `output = "null"`
- Modified auto-setup to create commented-out logging configuration files
- CLI now uses sensible defaults when no logging configuration exists
- User-friendly messages go to stdout, technical logs are suppressed by default

**Status**: ✅ **Fixed** - CLI output is now clean and user-friendly

## Performance Limitations

### Parallel Processing Issues

**Issue**: Limited parallel processing support across systems

**Affected Systems**:
- Actions System (`internal/actions/`)
- Facts System (`internal/facts/`)
- Templates System (`internal/templates/`)

**Problematic Code Examples**:

```go
// internal/actions/manager.go - Sequential action execution
func (m *Manager) runActionStep(ctx context.Context, session *spookytypesactions.ActingSession, actionNames []string, allActions []*spookytypesactions.Action, _ *spookytypesactions.ActionPlan) ([]spookytypesactions.ActingResult, error) {
    var results []spookytypesactions.ActingResult

    // Sequential processing - no parallel execution
    for _, actionName := range actionNames {
        var action *spookytypesactions.Action
        for _, a := range allActions {
            if a.Name == actionName {
                action = a
                break
            }
        }

        if action == nil {
            m.logger.Error("Action not found", fmt.Errorf("action %s not found", actionName), map[string]interface{}{
                "action": actionName,
            })
            continue
        }

        // Sequential action execution - no parallel processing
        actionResults, err := m.runAction(ctx, session, action)
        if err != nil {
            m.logger.Error("Failed to run action", err, map[string]interface{}{"action": actionName})
            continue
        }

        results = append(results, actionResults...)
    }

    return results, nil
}
```

```go
// cmd/facts.go - Sequential fact collection
func collectFactsParallel(ctx context.Context, machines []spookytypes.Machine, factsManager spookyinterfaces.FactsIntegration, parallel int, verbose bool, logger spookytypeslogging.Logger) (successCount, errorCount int) {
    // Limited parallel processing implementation
    // Uses semaphore but still has sequential bottlenecks
    semaphore := make(chan struct{}, parallel)

    for idx := range machines {
        machine := &machines[idx]
        wg.Add(1)
        go func(m *spookytypes.Machine) {
            defer wg.Done()
            semaphore <- struct{}{}        // Acquire
            defer func() { <-semaphore }() // Release

            // Sequential fact collection per machine
            facts, err := factsManager.CollectFacts(ctx, m)
            if err != nil {
                // Error handling
            }
        }(machine)
    }
}
```

**Symptoms**:
- Sequential processing only in most operations
- No multi-machine parallel processing
- Performance bottlenecks with large machine inventories
- Inefficient resource utilization
- No load balancing across machines

**Root Cause Analysis**:
1. **Sequential Execution**: Actions and facts are processed sequentially instead of in parallel
2. **No Connection Pooling**: Each operation creates new SSH connections
3. **Limited Concurrency**: Basic concurrency without proper load balancing
4. **No Resource Management**: No efficient resource allocation and management
5. **Bottleneck Operations**: Single-threaded operations create performance bottlenecks

**Current Workarounds**:
- Use machine filtering to limit scope
- Process machines in smaller batches
- Monitor updates for parallel processing improvements

### Caching Issues

**Issue**: No intelligent caching for frequently accessed data

**Affected Systems**:
- Facts System (`internal/facts/`)
- Templates System (`internal/templates/`)
- Variables System (`internal/variables/`)

**Problematic Code Examples**:

```go
// internal/facts/collector.go - No caching implementation
func (c *SystemFactCollector) Collect(ctx context.Context, machine *spookytypes.Machine) (*spookytypesfacts.FactCollection, error) {
    // No caching - collects facts every time
    // No cache validation or TTL
    // No cache invalidation mechanism
    
    machineID, err := c.getMachineID(machine)
    if err != nil {
        return nil, fmt.Errorf("failed to get machine ID: %w", err)
    }

    // Always collects fresh facts - no cache check
    facts := &spookytypesfacts.Facts{
        System: &spookytypesfacts.SystemFacts{},
    }

    systemFacts, err := c.collectSystemFacts(ctx, machine)
    if err != nil {
        return nil, fmt.Errorf("failed to collect system facts: %w", err)
    }
    facts.System = systemFacts

    // No caching of collected facts
    return &spookytypesfacts.FactCollection{
        MachineID:   machineID,
        CollectedAt: time.Now(),
        Facts:       facts,
    }, nil
}
```

```go
// internal/templates/manager.go - No template caching
func (m *Manager) RenderTemplate(ctx context.Context, template *spookytypestemplates.Template, data map[string]interface{}) (string, error) {
    // No template caching - renders every time
    // No compiled template caching
    // No result caching for repeated renders
    
    // Always renders from scratch
    return m.renderLocally(template, data)
}
```

**Symptoms**:
- Repeated data collection for the same information
- No template caching for repeated renders
- No variable resolution caching
- Inefficient resource usage
- No performance optimization for repeated operations

**Root Cause Analysis**:
1. **No Caching Layer**: Systems lack caching mechanisms for frequently accessed data
2. **No Cache Invalidation**: No mechanism to invalidate stale cached data
3. **No TTL Support**: No time-to-live support for cached data
4. **No Cache Persistence**: No persistent caching for data that could be reused
5. **No Cache Optimization**: No optimization for repeated operations

**Current Workarounds**:
- Cache data manually if needed
- Use local storage for frequently accessed data
- Monitor updates for caching improvements

### Connection Reuse Issues

**Issue**: Limited connection reuse optimization

**Affected Systems**:
- SSH System (`internal/ssh/`) when used by other systems
- All systems that depend on SSH

**Problematic Code Examples**:

```go
// internal/actions/manager.go - New connection per action
func (m *Manager) runCommandAction(ctx context.Context, action *spookytypesactions.Action, machine *spookytypes.Machine, result *spookytypesactions.ActingResult) error {
    // Creates new SSH connection for each action - no reuse
    connectionRequest := &spookytypes.ConnectionRequest{
        Host:     machine.Host,
        Port:     machine.Port,
        User:     machine.User,
        Password: machine.Password,
        KeyPath:  machine.KeyFile,
        Timeout:  time.Duration(action.Timeout) * time.Second,
    }

    // New connection every time - no connection pooling
    connectionResult, err := m.sshManager.Connect(ctx, connectionRequest)
    if err != nil {
        return fmt.Errorf("failed to connect to %s: %w", machine.Hostname, err)
    }

    // New session every time - no session reuse
    session, err := m.sshManager.CreateSession(ctx, connectionResult.Connection)
    if err != nil {
        return fmt.Errorf("failed to create session on %s: %w", machine.Hostname, err)
    }
}
```

```go
// internal/facts/collector.go - New connection per command
func (c *SystemFactCollector) runSSHCommand(machine *spookytypes.Machine, command string) (string, error) {
    ctx := context.Background()

    // New connection for each SSH command - no reuse
    connectionRequest := &spookytypes.ConnectionRequest{
        Host:       machine.Hostname,
        Port:       machine.Port,
        User:       machine.User,
        Password:   machine.Password,
        KeyPath:    machine.KeyFile,
        Passphrase: machine.Passphrase,
        AuthMethod: spookytypesssh.AuthMethodPublicKey,
        Timeout:    30 * time.Second,
    }

    // Establishes new connection for each command
    connectionResult, err := c.sshManager.Connect(ctx, connectionRequest)
    if err != nil {
        return "", fmt.Errorf("failed to establish SSH connection to %s: %w", machine.Hostname, err)
    }
}
```

**Symptoms**:
- New SSH connections for each operation
- Increased latency for multi-operation workflows
- Resource inefficiency
- High connection overhead
- No connection pooling

**Root Cause Analysis**:
1. **No Connection Pooling**: Each operation creates a new SSH connection
2. **No Session Reuse**: Sessions are not reused across operations
3. **No Connection Caching**: No caching mechanism for active connections
4. **High Overhead**: Connection establishment overhead for each operation
5. **Resource Waste**: Inefficient use of system resources

**Current Workarounds**:
- Batch operations when possible
- Use connection pooling manually if needed
- Monitor updates for connection optimization

## Integration Issues

### Secrets System Integration Issues

**Issue**: Limited integration with other systems

**Symptoms**:
- No SSH-based secret collection from remote machines
- Limited secret validation or lifecycle management
- No secret export/import functionality

**Current Workarounds**:
- Use local secret management
- Implement secret collection manually if needed
- Monitor updates for integration improvements

**Technical Details**:
- Age encryption/decryption works
- Basic CLI integration works
- Integration with other systems is incomplete

### Integrations System Issues

**Issue**: Many advanced features are still in development

**Symptoms**:
- Limited API integrations
- No webhook processing
- No database integrations
- No cloud provider integrations

**Current Workarounds**:
- Use basic integration functionality
- Implement custom integrations manually if needed
- Monitor updates for advanced feature development

## Validation and Error Handling Issues

### Validation Configuration Issues

**Issue**: Incomplete validation configuration in some systems

**Symptoms**:
- Limited validation rules
- Incomplete error reporting
- Missing validation for edge cases

**Current Workarounds**:
- Use basic validation where available
- Implement additional validation manually if needed
- Monitor updates for validation improvements

### Error Handling Issues

**Issue**: Inconsistent error handling across systems

**Symptoms**:
- Incomplete error messages
- Limited error recovery options
- Inconsistent error reporting formats

**Current Workarounds**:
- Use available error handling features
- Implement custom error handling if needed
- Monitor updates for error handling improvements

## Documentation Issues

### API Documentation Issues

**Issue**: Some API documentation has known issues

**Affected Documentation**:
- Logging API Reference
- Integrations API Reference
- Machines API Reference

**Symptoms**:
- Incomplete interface definitions
- Missing implementation details
- Outdated documentation

**Current Workarounds**:
- Use working API features
- Refer to source code for implementation details
- Monitor updates for documentation improvements

## Testing Issues

### Integration Testing Issues

**Issue**: Basic integration testing with known issues

**Symptoms**:
- Limited test coverage for SSH-based operations
- Incomplete test scenarios
- Missing edge case testing

**Current Workarounds**:
- Use available test coverage
- Implement additional tests manually if needed
- Monitor updates for testing improvements

## Development Workflow Issues

### Build and Development Issues

**Issue**: Some development tools have limitations

**Symptoms**:
- Limited tool integration
- Incomplete development workflows
- Missing automation features

**Current Workarounds**:
- Use available development tools
- Implement custom workflows if needed
- Monitor updates for tool improvements

## Reporting Issues

### Issue Reporting Process

**Issue**: Limited issue reporting and tracking

**Symptoms**:
- No centralized issue tracking
- Incomplete issue documentation
- Limited feedback mechanisms

**Current Workarounds**:
- Use available documentation
- Report issues through available channels
- Monitor updates for process improvements

## Workarounds and Mitigations

### General Workarounds

1. **Use Local Operations**: When SSH-based operations fail, use local alternatives
2. **Manual Execution**: Perform operations manually on remote machines when automation fails
3. **Filtered Scope**: Use machine filtering to limit the scope of operations
4. **Batch Processing**: Process operations in smaller batches to avoid performance issues
5. **Monitor Updates**: Watch for improvements and updates to affected systems

### System-Specific Workarounds

#### Actions System
- Use local action running for immediate needs
- Run actions manually on remote machines if needed
- Use machine filtering to limit running scope

#### Facts System
- Use local fact collection for immediate needs
- Export facts manually from remote machines if needed
- Use machine filtering to limit collection scope

#### Templates System
- Use local template rendering for immediate needs
- Deploy templates through the Actions system (when working)
- Render templates manually on remote machines if needed

#### Secrets System
- Use local secret management
- Implement secret collection manually if needed
- Use available encryption/decryption features

## Specific Files and Code Paths Requiring Fixes

### Critical Files for SSH Integration Fixes

#### 1. `internal/actions/manager.go`
**Lines 420-558**: `runCommandAction` function
- **Issue**: Creates new SSH connection for each action
- **Fix Required**: Implement connection pooling and reuse
- **Impact**: High - affects all action orchestration

**Lines 560-658**: `runScriptAction` function  
- **Issue**: Duplicates connection logic from `runCommandAction`
- **Fix Required**: Extract common connection logic to shared function
- **Impact**: Medium - code duplication and inefficiency

**Lines 660-720**: `runTemplateAction` and `runFileCopyAction` functions
- **Issue**: Similar connection duplication patterns
- **Fix Required**: Use shared connection management
- **Impact**: Medium - code duplication

#### 2. `internal/facts/collector.go`
**Lines 779-844**: `runSSHCommand` function
- **Issue**: Creates new SSH connection for each fact collection command
- **Fix Required**: Implement connection pooling for fact collection
- **Impact**: High - affects all SSH-based fact collection

**Lines 212-779**: `collectSystemFactsViaSSH` function
- **Issue**: Multiple SSH commands without connection reuse
- **Fix Required**: Batch commands or reuse connections
- **Impact**: High - inefficient fact collection

#### 3. `internal/templates/manager.go`
**Missing SSH Integration**: No SSH-based template rendering
- **Issue**: Templates can only be rendered locally
- **Fix Required**: Add SSH manager integration and remote rendering
- **Impact**: High - no remote template rendering capability

#### 4. `internal/ssh/manager.go`
**Lines 229-352**: `RunAction` function
- **Issue**: Basic SSH action orchestration without optimization
- **Fix Required**: Improve connection management and error handling
- **Impact**: Medium - affects SSH-based operations

### Critical Files for CLI Command Fixes

#### 1. `cmd/templates.go`
**Lines 1-195**: `templatesRenderCmd` implementation
- **Issue**: Limited functionality, no SSH integration
- **Fix Required**: Add SSH integration and remote rendering support
- **Impact**: High - affects template CLI usability

**Lines 195+**: `templatesValidateCmd` implementation
- **Issue**: Basic validation only, no remote validation
- **Fix Required**: Add remote template validation capabilities
- **Impact**: Medium - affects template validation

#### 2. `cmd/actions.go`
**Limited CLI Implementation**: Missing advanced features
- **Issue**: No parallel execution support, limited result aggregation
- **Fix Required**: Add parallel execution and result aggregation
- **Impact**: High - affects action CLI usability

### Critical Files for Performance Fixes

#### 1. `internal/actions/manager.go`
**Lines 280-320**: `runActionStep` function
- **Issue**: Sequential action execution
- **Fix Required**: Implement parallel execution with proper synchronization
- **Impact**: High - affects action performance

#### 2. `cmd/facts.go`
**Lines 251-332**: `collectFactsParallel` function
- **Issue**: Limited parallel processing implementation
- **Fix Required**: Improve parallel processing with proper resource management
- **Impact**: High - affects fact collection performance

### Critical Files for Caching Fixes

#### 1. `internal/facts/collector.go`
**Lines 40-100**: `Collect` function
- **Issue**: No caching implementation
- **Fix Required**: Add intelligent caching with TTL and invalidation
- **Impact**: High - affects fact collection efficiency

#### 2. `internal/templates/manager.go`
**Missing Caching**: No template caching
- **Issue**: Templates rendered from scratch every time
- **Fix Required**: Add template compilation and result caching
- **Impact**: Medium - affects template rendering performance

## Future Improvements

### Planned Fixes

1. **SSH Integration**: Complete SSH integration for all affected systems
2. **Parallel Processing**: Implement parallel processing across all systems
3. **Caching**: Add intelligent caching for frequently accessed data
4. **CLI Commands**: Complete CLI command implementation for all systems
5. **Error Handling**: Improve error handling and recovery across all systems

### Development Priorities

1. **Critical**: Fix SSH-based operations in Actions, Facts, and Templates systems
2. **High**: Implement parallel processing and caching
3. **Medium**: Complete CLI command implementation
4. **Low**: Improve documentation and testing

## Contributing to Issue Resolution

### How to Help

1. **Report Issues**: Document new issues with detailed reproduction steps
2. **Test Workarounds**: Verify that workarounds work in your environment
3. **Provide Feedback**: Share experiences with current limitations
4. **Monitor Updates**: Watch for improvements and updates

### Issue Reporting Guidelines

When reporting issues:

1. **Describe the Problem**: Clear description of what's not working
2. **Provide Steps**: Detailed steps to reproduce the issue
3. **Include Environment**: System details and configuration
4. **Show Error Messages**: Complete error messages and logs
5. **Suggest Workarounds**: Any workarounds you've found

## Related Documentation

- [SSH User Guide](SSH_USER_GUIDE.md) - SSH connectivity and authentication
- [Actions User Guide](ACTIONS_USER_GUIDE.md) - Action orchestration issues
- [Facts User Guide](FACTS_USER_GUIDE.md) - Fact collection issues
- [Templates User Guide](TEMPLATES_USER_GUIDE.md) - Template rendering issues
- [Secrets User Guide](SECRETS_USER_GUIDE.md) - Secret management issues
- [Integrations User Guide](INTEGRATIONS_USER_GUIDE.md) - Integration issues
- [User Guides Index](USER_GUIDES_INDEX.md) - Complete overview of all user guides

## Last Updated

This document was last updated on: **2024-12-15**

**Note**: This document is actively maintained and updated as issues are resolved and new issues are discovered.
