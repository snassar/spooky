# Implementation Plan: SSH Acting Engine

## Overview
Implement SSH-based action acting engine that can act actions, templates, and commands over SSH connections, support sequential and parallel acting, and provide comprehensive error handling and result collection.

## Task Details
- **Task ID**: 4.2
- **Priority**: High
- **File**: `internal/ssh/acting/manager.go`
- **Functions**: `ActAction`, `ActTemplate`, `ActSequential`, `ActParallel`


## Current State Analysis

### Existing Patterns
1. **SSH Manager**: Existing SSH manager in `internal/ssh/manager.go` provides connection and command acting
2. **Action Types**: Actions use `types.Action` with command, script, template configurations
3. **Context System**: `types.ActionContext` provides acting context with facts, variables, machines
4. **Result Structure**: `types.RunResult` captures acting results with output, error, exit code
5. **Error Handling**: Consistent error wrapping with context
6. **Connection Management**: SSH manager handles connection lifecycle

### Existing Implementation Examples
- **SSH Manager**: `internal/ssh/manager.go` provides connection and command acting
- **Action Types**: `internal/actions/types/action.go` defines action structure
- **Context System**: `internal/actions/types/context.go` provides acting context

## Implementation Requirements

### Interface Compliance
The SSH acting engine must:
1. **Act actions** over SSH connections
2. **Act templates** with rendering and deployment
3. **Support sequential acting** of multiple operations
4. **Support parallel acting** with concurrency control
5. **Handle connection management** and pooling
6. **Collect and aggregate** results from multiple machines
7. **Implement timeout** and error handling
8. **Support different action types** (command, script, template)

### Required Dependencies
- SSH manager for connection management
- Action context for acting context
- Template engine for template rendering
- Result aggregation system
- Connection pooling and management

## Detailed Implementation Plan

### Step 1: Create SSH Acting Engine Structure

**File**: `internal/ssh/acting/manager.go`

```go
package acting

import (
    "context"
    "fmt"
    "sync"
    "time"

    "spooky/internal/actions/types"
    "spooky/internal/config/types"
    "spooky/internal/logging"
    "spooky/internal/ssh"
    sshTypes "spooky/internal/ssh/types"
    "spooky/internal/templates"
)

// Manager manages SSH-based action acting
type Manager struct {
    sshManager ssh.Manager
    logger     logging.Logger
    pool       *connectionPool
}

// NewManager creates a new SSH acting manager
func NewManager(sshManager ssh.Manager, logger logging.Logger) *Manager {
    return &Manager{
        sshManager: sshManager,
        logger:     logger,
        pool:       newConnectionPool(logger),
    }
}
```

### Step 2: Implement Connection Pooling

#### 2.1 Connection Pool Structure
```go
// connectionPool manages SSH connections
type connectionPool struct {
    connections map[string]*pooledConnection
    mutex       sync.RWMutex
    logger      logging.Logger
}

// pooledConnection represents a pooled SSH connection
type pooledConnection struct {
    connection sshTypes.Connection
    machine    configtypes.Machine
    lastUsed   time.Time
    inUse      bool
}

// newConnectionPool creates a new connection pool
func newConnectionPool(logger logging.Logger) *connectionPool {
    return &connectionPool{
        connections: make(map[string]*pooledConnection),
        logger:      logger,
    }
}
```

#### 2.2 Connection Pool Operations
```go
// getConnection gets a connection from the pool or creates a new one
func (p *connectionPool) getConnection(machine configtypes.Machine) (sshTypes.Connection, error) {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    key := fmt.Sprintf("%s:%d", machine.Host, machine.Port)
    
    // Check if connection exists and is available
    if pooled, exists := p.connections[key]; exists && !pooled.inUse {
        pooled.inUse = true
        pooled.lastUsed = time.Now()
        return pooled.connection, nil
    }
    
    // Create new connection
    sshConfig := &sshTypes.SSHConfig{
        Host:     machine.Host,
        Port:     machine.Port,
        Username: machine.User,
        Timeout:  30 * time.Second,
    }
    
    connection, err := sshManager.Connect(machine.Host, sshConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to create SSH connection: %w", err)
    }
    
    // Add to pool
    p.connections[key] = &pooledConnection{
        connection: connection,
        machine:    machine,
        lastUsed:   time.Now(),
        inUse:      true,
    }
    
    return connection, nil
}

// returnConnection returns a connection to the pool
func (p *connectionPool) returnConnection(machine configtypes.Machine) {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    key := fmt.Sprintf("%s:%d", machine.Host, machine.Port)
    
    if pooled, exists := p.connections[key]; exists {
        pooled.inUse = false
        pooled.lastUsed = time.Now()
    }
}

// closeConnection closes and removes a connection from the pool
func (p *connectionPool) closeConnection(machine configtypes.Machine) error {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    key := fmt.Sprintf("%s:%d", machine.Host, machine.Port)
    
    if pooled, exists := p.connections[key]; exists {
        err := sshManager.CloseConnection(pooled.connection)
        delete(p.connections, key)
        return err
    }
    
    return nil
}
```

### Step 3: Implement Action Acting

#### 3.1 Main Action Acting Method
```go
// ActAction acts an action on a machine via SSH
func (m *Manager) ActAction(ctx context.Context, action *types.Action, machine configtypes.Machine, context *types.ActionContext) (*types.RunResult, error) {
    m.logger.Debug("Acting action via SSH",
        logging.String("action", action.Name),
        logging.String("machine", machine.Name),
        logging.String("type", action.Type))

    result := &types.RunResult{
        ActionName: action.Name,
        MachineID:  machine.Name,
        Status:     types.RunStatusRunning,
        StartTime:  time.Now(),
    }

    // Get connection from pool
    connection, err := m.pool.getConnection(machine)
    if err != nil {
        result.Status = types.RunStatusFailed
        result.Error = err.Error()
        result.EndTime = time.Now()
        result.Duration = result.EndTime.Sub(*result.StartTime)
        return result, fmt.Errorf("failed to get SSH connection: %w", err)
    }
    defer m.pool.returnConnection(machine)

    // Act based on action type
    switch action.Type {
    case "command":
        err = m.actCommandAction(ctx, connection, action, context, result)
    case "script":
        err = m.actScriptAction(ctx, connection, action, context, result)
    case "template_deploy":
        err = m.actTemplateAction(ctx, connection, action, context, result)
    default:
        err = fmt.Errorf("unsupported action type: %s", action.Type)
    }

    if err != nil {
        result.Status = types.RunStatusFailed
        result.Error = err.Error()
        result.EndTime = time.Now()
        result.Duration = result.EndTime.Sub(*result.StartTime)
        return result, err
    }

    result.Status = types.RunStatusCompleted
    result.EndTime = time.Now()
    result.Duration = result.EndTime.Sub(*result.StartTime)

    return result, nil
}
```

#### 3.2 Command Action Acting
```go
// actCommandAction acts a command action
func (m *Manager) actCommandAction(ctx context.Context, connection sshTypes.Connection, action *types.Action, context *types.ActionContext, result *types.RunResult) error {
    // Build acting command
    actingCmd := m.buildActingCommand(action.Command, context)
    
    // Act command
    commandResult, err := m.sshManager.ExecuteCommand(connection, actingCmd)
    if err != nil {
        return fmt.Errorf("command acting failed: %w", err)
    }

    result.Output = commandResult.Stdout
    result.Error = commandResult.Stderr
    result.ExitCode = commandResult.ExitCode

    return nil
}
```

#### 3.3 Script Action Acting
```go
// actScriptAction acts a script action
func (m *Manager) actScriptAction(ctx context.Context, connection sshTypes.Connection, action *types.Action, context *types.ActionContext, result *types.RunResult) error {
    // Upload script
    remoteScriptPath, err := m.uploadScript(connection, action.Script, context)
    if err != nil {
        return fmt.Errorf("failed to upload script: %w", err)
    }
    defer m.cleanupScript(connection, remoteScriptPath)

    // Build acting command
    actingCmd := m.buildScriptActingCommand(remoteScriptPath, action, context)
    
    // Act script
    scriptResult, err := m.sshManager.ExecuteCommand(connection, actingCmd)
    if err != nil {
        return fmt.Errorf("script acting failed: %w", err)
    }

    result.Output = scriptResult.Stdout
    result.Error = scriptResult.Stderr
    result.ExitCode = scriptResult.ExitCode

    return nil
}
```

#### 3.4 Template Action Acting
```go
// actTemplateAction acts a template action
func (m *Manager) actTemplateAction(ctx context.Context, connection sshTypes.Connection, action *types.Action, context *types.ActionContext, result *types.RunResult) error {
    // Render template
    renderedContent, err := m.renderTemplate(action.Template.Source, context)
    if err != nil {
        return fmt.Errorf("template rendering failed: %w", err)
    }

    // Upload rendered template
    err = m.uploadTemplate(connection, renderedContent, action.Template.Destination)
    if err != nil {
        return fmt.Errorf("template upload failed: %w", err)
    }

    // Set file attributes
    err = m.setFileAttributes(connection, action.Template)
    if err != nil {
        return fmt.Errorf("failed to set file attributes: %w", err)
    }

    result.Output = "Template deployed successfully"
    result.ExitCode = 0

    return nil
}
```

### Step 4: Implement Sequential Acting

#### 4.1 Sequential Acting Method
```go
// ActSequential acts actions sequentially on multiple machines
func (m *Manager) ActSequential(ctx context.Context, action *types.Action, machines []configtypes.Machine, context *types.ActionContext) ([]*types.RunResult, error) {
    var results []*types.RunResult
    var errors []error

    for _, machine := range machines {
        result, err := m.ActAction(ctx, action, machine, context)
        if err != nil {
            errors = append(errors, fmt.Errorf("machine %s: %w", machine.Name, err))
            if !action.AllowFailure {
                return results, fmt.Errorf("sequential acting failed: %v", errors)
            }
        }
        results = append(results, result)
    }

    return results, nil
}
```

### Step 5: Implement Parallel Acting

#### 5.1 Parallel Acting Method
```go
// ActParallel acts actions in parallel on multiple machines
func (m *Manager) ActParallel(ctx context.Context, action *types.Action, machines []configtypes.Machine, context *types.ActionContext) ([]*types.RunResult, error) {
    // Determine concurrency limit
    maxConcurrent := context.MaxConcurrent
    if maxConcurrent <= 0 {
        maxConcurrent = len(machines)
    }

    // Create worker pool
    semaphore := make(chan struct{}, maxConcurrent)
    results := make(chan *parallelExecutionResult, len(machines))
    errors := make(chan error, len(machines))

    // Start workers
    for _, machine := range machines {
        go func(m configtypes.Machine) {
            semaphore <- struct{}{} // Acquire semaphore
            defer func() { <-semaphore }() // Release semaphore

            result, err := m.ActAction(ctx, action, m, context)
            if err != nil {
                errors <- fmt.Errorf("machine %s: %w", m.Name, err)
            } else {
                results <- &parallelRunResult{
                    machine: m,
                    result:  result,
                }
            }
        }(machine)
    }

    // Collect results
    var allResults []*types.RunResult
    var allErrors []error

    for i := 0; i < len(machines); i++ {
        select {
        case executionResult := <-results:
            allResults = append(allResults, executionResult.result)
        case err := <-errors:
            allErrors = append(allErrors, err)
            if !action.AllowFailure {
                return allResults, fmt.Errorf("parallel acting failed: %v", allErrors)
            }
        case <-ctx.Done():
            return allResults, fmt.Errorf("parallel acting cancelled: %w", ctx.Err())
        }
    }

    return allResults, nil
}
```

### Step 6: Implement Template Acting

#### 6.1 Template Acting Method
```go
// ActTemplate acts a template on a machine
func (m *Manager) ActTemplate(ctx context.Context, templatePath string, machine configtypes.Machine, context *types.ActionContext) (*types.RunResult, error) {
    m.logger.Debug("Acting template via SSH",
        logging.String("template", templatePath),
        logging.String("machine", machine.Name))

    result := &types.RunResult{
        ActionName: "template_execution",
        MachineID:  machine.Name,
        Status:     types.RunStatusRunning,
        StartTime:  time.Now(),
    }

    // Get connection from pool
    connection, err := m.pool.getConnection(machine)
    if err != nil {
        result.Status = types.RunStatusFailed
        result.Error = err.Error()
        result.EndTime = time.Now()
        result.Duration = result.EndTime.Sub(*result.StartTime)
        return result, fmt.Errorf("failed to get SSH connection: %w", err)
    }
    defer m.pool.returnConnection(machine)

    // Render template
    renderedContent, err := m.renderTemplate(templatePath, context)
    if err != nil {
        result.Status = types.RunStatusFailed
        result.Error = err.Error()
        result.EndTime = time.Now()
        result.Duration = result.EndTime.Sub(*result.StartTime)
        return result, fmt.Errorf("template rendering failed: %w", err)
    }

    // Act rendered content as command
    commandResult, err := m.sshManager.ExecuteCommand(connection, string(renderedContent))
    if err != nil {
        result.Status = types.RunStatusFailed
        result.Error = err.Error()
        result.EndTime = time.Now()
        result.Duration = result.EndTime.Sub(*result.StartTime)
        return result, fmt.Errorf("template acting failed: %w", err)
    }

    result.Status = types.RunStatusCompleted
    result.Output = commandResult.Stdout
    result.Error = commandResult.Stderr
    result.ExitCode = commandResult.ExitCode
    result.EndTime = time.Now()
    result.Duration = result.EndTime.Sub(*result.StartTime)

    return result, nil
}
```

### Step 7: Implement Supporting Methods

#### 7.1 Command Building
```go
// buildActingCommand builds the command for acting
func (m *Manager) buildActingCommand(command string, context *types.ActionContext) string {
    var cmd strings.Builder
    
    // Set working directory
    if context.WorkingDirectory != "" {
        cmd.WriteString(fmt.Sprintf("cd %s && ", context.WorkingDirectory))
    }
    
    // Set environment variables
    for key, value := range context.Environment {
        cmd.WriteString(fmt.Sprintf("export %s='%s' && ", key, value))
    }
    
    // Add sudo if required
    if context.Sudo {
        cmd.WriteString("sudo ")
    }
    
    // Act command
    cmd.WriteString(command)
    
    return cmd.String()
}
```

#### 7.2 Template Rendering
```go
// renderTemplate renders a template with context data
func (m *Manager) renderTemplate(templatePath string, context *types.ActionContext) ([]byte, error) {
    // Read template content
    content, err := os.ReadFile(templatePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read template: %w", err)
    }
    
    // Create template context
    templateContext := &templates.TemplateContext{
        Facts:     context.Facts,
        Variables: context.Variables,
        Machines:  context.Machines,
        Metadata: map[string]interface{}{
            "project_path": context.ProjectPath,
            "project_name": context.ProjectName,
            "session_id":   context.SessionID,
        },
    }
    
    // Create template engine
    engine := templates.NewEngine(m.logger)
    
    // Render template
    rendered, err := engine.RenderTemplate(string(content), templateContext)
    if err != nil {
        return nil, fmt.Errorf("template rendering failed: %w", err)
    }
    
    return []byte(rendered), nil
}
```

#### 7.3 File Upload Operations
```go
// uploadScript uploads a script to the remote machine
func (m *Manager) uploadScript(connection sshTypes.Connection, scriptPath string, context *types.ActionContext) (string, error) {
    // Read script content
    content, err := os.ReadFile(scriptPath)
    if err != nil {
        return "", fmt.Errorf("failed to read script: %w", err)
    }

    // Generate remote path
    scriptName := filepath.Base(scriptPath)
    remotePath := fmt.Sprintf("/tmp/spooky_script_%d_%s", time.Now().Unix(), scriptName)

    // Upload file via SCP
    err = m.uploadFileViaSCP(connection, content, remotePath)
    if err != nil {
        return "", fmt.Errorf("failed to upload script: %w", err)
    }

    // Set execute permissions
    chmodCmd := fmt.Sprintf("chmod +x %s", remotePath)
    chmodResult, err := m.sshManager.ExecuteCommand(connection, chmodCmd)
    if err != nil {
        return "", fmt.Errorf("failed to set script permissions: %w", err)
    }

    if chmodResult.ExitCode != 0 {
        return "", fmt.Errorf("failed to set script permissions: %s", chmodResult.Stderr)
    }

    return remotePath, nil
}
```

### Step 8: Add Supporting Structures

#### 8.1 Parallel Execution Result
```go
// parallelExecutionResult holds result of parallel execution
type parallelExecutionResult struct {
    machine configtypes.Machine
    result  *types.RunResult
}
```

#### 8.2 Connection Pool Statistics
```go
// getPoolStats returns connection pool statistics
func (p *connectionPool) getPoolStats() map[string]interface{} {
    p.mutex.RLock()
    defer p.mutex.RUnlock()

    stats := map[string]interface{}{
        "total_connections": len(p.connections),
        "active_connections": 0,
        "idle_connections":   0,
    }

    for _, conn := range p.connections {
        if conn.inUse {
            stats["active_connections"] = stats["active_connections"].(int) + 1
        } else {
            stats["idle_connections"] = stats["idle_connections"].(int) + 1
        }
    }

    return stats
}
```







## Configuration Options

### Supported Options
- **MaxConcurrent**: Limit concurrent connections
- **ConnectionTimeout**: SSH connection timeout
- **CommandTimeout**: Command execution timeout
- **PoolSize**: Maximum connection pool size
- **RetryAttempts**: Connection retry attempts

## Dependencies

### Internal Dependencies
- `spooky/internal/actions/types`
- `spooky/internal/config/types`
- `spooky/internal/ssh`
- `spooky/internal/templates`
- `spooky/internal/logging`

### External Dependencies
- `context` (standard library)
- `time` (standard library)
- `strings` (standard library)
- `fmt` (standard library)
- `os` (standard library)
- `path/filepath` (standard library)
- `sync` (standard library)



## Implementation Order

1. Create SSH acting engine structure
2. Implement connection pooling
3. Implement action execution methods
4. Add sequential execution
5. Add parallel execution
6. Implement template execution
7. Add supporting methods
8. Add file upload operations
9. Add supporting structures
10. Write comprehensive tests
11. Performance optimization
12. Documentation and cleanup


