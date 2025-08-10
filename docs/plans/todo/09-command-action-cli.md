# Implementation Plan: Command Acting CLI

## Overview
Implement command acting CLI functionality that can act commands on target machines, handle parallel acting, provide real-time output, and manage command acting lifecycle.

## Task Details
- **Task ID**: 5.2
- **Priority**: Medium
- **File**: `internal/cli/commands/executor.go`
- **Function**: `actCommand`


## Current State Analysis

### Existing Patterns
1. **CLI Structure**: CLI commands use `spf13/cobra` framework with consistent patterns
2. **Command Context**: Commands receive context with project, machines, and configuration
3. **SSH Manager**: Existing SSH manager provides command acting capabilities
4. **Error Handling**: Consistent error wrapping with context
5. **Output Formatting**: Structured output with JSON and text formats
6. **Parallel Acting**: Support for parallel operations across machines

### Existing Implementation Examples
- **CLI Commands**: `internal/cli/commands/` provides command implementations
- **SSH Manager**: `internal/ssh/manager.go` provides SSH operations
- **Command Types**: `internal/cli/types/command.go` defines command structure
- **Output Formatting**: `internal/cli/help/renderer.go` provides output formatting

## Implementation Requirements

### Interface Compliance
The command execution CLI must:
1. **Act commands** on target machines via SSH using `--machines` flag (plural)
2. **Handle parallel acting** across multiple machines with `--parallel` flag
3. **Provide real-time output** and progress reporting
4. **Support different output formats** (text, JSON, table) where applicable
5. **Handle command timeouts** and cancellation
6. **Support command chaining** and conditional acting
7. **Provide detailed acting reporting** and status
8. **Support dry-run mode** for safe acting
9. **Support tag-based filtering** using `--tags` flag
10. **Remove support for deprecated `--filter` flag**
11. **Remove short flags** (`-a`, `-m`, `-d`) from actions commands

### Required Dependencies
- SSH manager for command execution
- CLI framework for command handling
- Output formatting system
- Parallel execution system

## Detailed Implementation Plan

### Step 1: Create Command Acting CLI Structure

**File**: `internal/cli/commands/executor.go`

```go
package commands

import (
    "context"
    "fmt"
    "strings"
    "time"

    "github.com/spf13/cobra"
    "spooky/internal/cli/types"
    "spooky/internal/logging"
    "spooky/internal/ssh"
    sshTypes "spooky/internal/ssh/types"
)

// CommandExecutor handles command acting via CLI
type CommandExecutor struct {
    sshManager ssh.Manager
    logger     logging.Logger
}

// RunResult represents the result of command acting
type RunResult struct {
    MachineID     string                 `json:"machine_id"`
    Command       string                 `json:"command"`
    ExitCode      int                    `json:"exit_code"`
    Stdout        string                 `json:"stdout"`
    Stderr        string                 `json:"stderr"`
    ExecutionTime time.Duration          `json:"execution_time"`
    Status        ExecutionStatus        `json:"status"`
    Error         string                 `json:"error,omitempty"`
    Metadata      map[string]interface{} `json:"metadata"`
}

// RunStatus represents the status of command acting
type RunStatus string

const (
    RunStatusPending   RunStatus = "pending"
    RunStatusRunning   RunStatus = "running"
    RunStatusCompleted RunStatus = "completed"
    RunStatusFailed    RunStatus = "failed"
    RunStatusTimeout   RunStatus = "timeout"
    RunStatusCancelled RunStatus = "cancelled"
)

// ActingOptions represents options for command acting
type ActingOptions struct {
    Parallel      bool                   `json:"parallel"`
    Timeout       time.Duration          `json:"timeout"`
    DryRun        bool                   `json:"dry_run"`
    OutputFormat  string                 `json:"output_format"`
    ShowProgress  bool                   `json:"show_progress"`
    StopOnError   bool                   `json:"stop_on_error"`
    Environment   map[string]string      `json:"environment"`
    WorkingDir    string                 `json:"working_dir"`
    User          string                 `json:"user"`
    Sudo          bool                   `json:"sudo"`
    RetryCount    int                    `json:"retry_count"`
    RetryDelay    time.Duration          `json:"retry_delay"`
}

// NewCommandExecutor creates a new command executor
func NewCommandExecutor(sshManager ssh.Manager, logger logging.Logger) *CommandExecutor {
    return &CommandExecutor{
        sshManager: sshManager,
        logger:     logger,
    }
}
```

### Step 2: Implement Command Execution CLI Command

#### 2.1 Main CLI Command
```go
// NewExecuteCommand creates the command execution CLI command
func NewExecuteCommand(executor *CommandExecutor) *cobra.Command {
    var (
        parallel     bool
        timeout      time.Duration
        dryRun       bool
        outputFormat string
        showProgress bool
        stopOnError  bool
        workingDir   string
        user         string
        sudo         bool
        retryCount   int
        retryDelay   time.Duration
        environment  []string
    )

    cmd := &cobra.Command{
        Use:   "execute [project] [command]",
        Short: "Execute commands on target machines",
        Long: `Execute commands on target machines in a project.

Examples:
  spooky execute myproject "ls -la"
  spooky execute myproject "systemctl status nginx" --parallel
  spooky execute myproject "apt update" --sudo --timeout 5m
  spooky execute myproject "echo $PATH" --env PATH=/usr/local/bin:/usr/bin`,
        Args: cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            projectPath := args[0]
            command := args[1]

            // Parse environment variables
            envMap := make(map[string]string)
            for _, env := range environment {
                if parts := strings.SplitN(env, "=", 2); len(parts) == 2 {
                    envMap[parts[0]] = parts[1]
                }
            }

            options := &ExecutionOptions{
                Parallel:     parallel,
                Timeout:      timeout,
                DryRun:       dryRun,
                OutputFormat: outputFormat,
                ShowProgress: showProgress,
                StopOnError:  stopOnError,
                Environment:  envMap,
                WorkingDir:   workingDir,
                User:         user,
                Sudo:         sudo,
                RetryCount:   retryCount,
                RetryDelay:   retryDelay,
            }

            return executor.ExecuteCommand(cmd.Context(), projectPath, command, options)
        },
    }

    // Add flags
    cmd.Flags().BoolVarP(&parallel, "parallel", "p", false, "Execute commands in parallel")
    cmd.Flags().DurationVarP(&timeout, "timeout", "t", 5*time.Minute, "Command execution timeout")
    cmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Show what would be executed without running")
    cmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text, json, table)")
    cmd.Flags().BoolVarP(&showProgress, "progress", "", true, "Show execution progress")
    cmd.Flags().BoolVarP(&stopOnError, "stop-on-error", "", false, "Stop execution on first error")
    cmd.Flags().StringVarP(&workingDir, "working-dir", "w", "", "Working directory for command execution")
    cmd.Flags().StringVarP(&user, "user", "u", "", "User to execute command as")
    cmd.Flags().BoolVarP(&sudo, "sudo", "s", false, "Execute command with sudo")
    cmd.Flags().IntVarP(&retryCount, "retry", "r", 0, "Number of retries on failure")
    cmd.Flags().DurationVarP(&retryDelay, "retry-delay", "", 5*time.Second, "Delay between retries")
    cmd.Flags().StringArrayVarP(&environment, "env", "e", nil, "Environment variables (KEY=VALUE)")

    return cmd
}
```

#### 2.2 Command Execution Method
```go
// ExecuteCommand executes a command on target machines
func (e *CommandExecutor) ExecuteCommand(ctx context.Context, projectPath, command string, options *ExecutionOptions) error {
    e.logger.Info("Executing command",
        logging.String("project", projectPath),
        logging.String("command", command),
        logging.Bool("parallel", options.Parallel),
        logging.Bool("dry_run", options.DryRun))

    // Load project configuration
    project, err := e.loadProject(projectPath)
    if err != nil {
        return fmt.Errorf("failed to load project: %w", err)
    }

    // Load machine inventory
    machines, err := e.loadMachines(projectPath)
    if err != nil {
        return fmt.Errorf("failed to load machines: %w", err)
    }

    if len(machines) == 0 {
        return fmt.Errorf("no machines found in project")
    }

    // Prepare command for execution
    preparedCommand := e.prepareCommand(command, options)

    // Execute command based on parallel flag
    var results []*ExecutionResult
    if options.Parallel {
        results, err = e.executeParallel(ctx, machines, preparedCommand, options)
    } else {
        results, err = e.executeSequential(ctx, machines, preparedCommand, options)
    }

    if err != nil {
        return fmt.Errorf("command execution failed: %w", err)
    }

    // Format and display results
    return e.displayResults(results, options)
}
```

### Step 3: Implement Command Preparation

#### 3.1 Command Preparation
```go
// prepareCommand prepares the command for execution
func (e *CommandExecutor) prepareCommand(command string, options *ExecutionOptions) string {
    preparedCommand := command

    // Add working directory if specified
    if options.WorkingDir != "" {
        preparedCommand = fmt.Sprintf("cd %s && %s", options.WorkingDir, preparedCommand)
    }

    // Add environment variables if specified
    if len(options.Environment) > 0 {
        envVars := make([]string, 0, len(options.Environment))
        for key, value := range options.Environment {
            envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
        }
        preparedCommand = fmt.Sprintf("export %s && %s", strings.Join(envVars, " "), preparedCommand)
    }

    // Add sudo if requested
    if options.Sudo {
        preparedCommand = fmt.Sprintf("sudo %s", preparedCommand)
    }

    // Add user execution if specified
    if options.User != "" {
        preparedCommand = fmt.Sprintf("su - %s -c '%s'", options.User, preparedCommand)
    }

    return preparedCommand
}
```

#### 3.2 Project and Machine Loading
```go
// loadProject loads project configuration
func (e *CommandExecutor) loadProject(projectPath string) (*types.Project, error) {
    // This would integrate with the existing project loading system
    // For now, we'll use a placeholder implementation
    return &types.Project{
        Path: projectPath,
        Name: "project",
    }, nil
}

// loadMachines loads machine inventory
func (e *CommandExecutor) loadMachines(projectPath string) ([]*types.Machine, error) {
    // This would integrate with the existing machine loading system
    // For now, we'll use a placeholder implementation
    return []*types.Machine{
        {
            Name: "localhost",
            Host: "127.0.0.1",
            Port: 22,
            User: "root",
        },
    }, nil
}
```

### Step 4: Implement Parallel Execution

#### 4.1 Parallel Execution Method
```go
// executeParallel executes commands in parallel across machines
func (e *CommandExecutor) executeParallel(ctx context.Context, machines []*types.Machine, command string, options *ExecutionOptions) ([]*ExecutionResult, error) {
    e.logger.Debug("Executing commands in parallel",
        logging.Int("machine_count", len(machines)))

    results := make([]*ExecutionResult, 0, len(machines))
    resultChan := make(chan *ExecutionResult, len(machines))
    errorChan := make(chan error, len(machines))

    // Start execution for each machine
    for _, machine := range machines {
        go func(m *types.Machine) {
            result, err := e.executeOnMachine(ctx, m, command, options)
            if err != nil {
                errorChan <- fmt.Errorf("machine %s: %w", m.Name, err)
                return
            }
            resultChan <- result
        }(machine)
    }

    // Collect results
    completed := 0
    for completed < len(machines) {
        select {
        case result := <-resultChan:
            results = append(results, result)
            completed++
            
            if options.ShowProgress {
                e.showProgress(result, completed, len(machines))
            }

            // Check if we should stop on error
            if options.StopOnError && result.Status == ExecutionStatusFailed {
                e.logger.Warn("Stopping execution due to error",
                    logging.String("machine", result.MachineID))
                return results, fmt.Errorf("execution stopped due to error on machine %s", result.MachineID)
            }

        case err := <-errorChan:
            e.logger.Error("Execution error",
                logging.Error(err))
            return results, err

        case <-ctx.Done():
            e.logger.Warn("Execution cancelled")
            return results, ctx.Err()
        }
    }

    return results, nil
}
```

#### 4.2 Sequential Execution Method
```go
// executeSequential executes commands sequentially across machines
func (e *CommandExecutor) executeSequential(ctx context.Context, machines []*types.Machine, command string, options *ExecutionOptions) ([]*ExecutionResult, error) {
    e.logger.Debug("Executing commands sequentially",
        logging.Int("machine_count", len(machines)))

    results := make([]*ExecutionResult, 0, len(machines))

    for i, machine := range machines {
        if options.ShowProgress {
            e.showProgress(nil, i, len(machines))
        }

        result, err := e.executeOnMachine(ctx, machine, command, options)
        if err != nil {
            e.logger.Error("Execution failed",
                logging.String("machine", machine.Name),
                logging.Error(err))
            return results, fmt.Errorf("machine %s: %w", machine.Name, err)
        }

        results = append(results, result)

        // Check if we should stop on error
        if options.StopOnError && result.Status == ExecutionStatusFailed {
            e.logger.Warn("Stopping execution due to error",
                logging.String("machine", result.MachineID))
            return results, fmt.Errorf("execution stopped due to error on machine %s", result.MachineID)
        }

        // Check for cancellation
        select {
        case <-ctx.Done():
            e.logger.Warn("Execution cancelled")
            return results, ctx.Err()
        default:
        }
    }

    return results, nil
}
```

### Step 5: Implement Machine Execution

#### 5.1 Machine Execution Method
```go
// executeOnMachine executes a command on a specific machine
func (e *CommandExecutor) executeOnMachine(ctx context.Context, machine *types.Machine, command string, options *ExecutionOptions) (*ExecutionResult, error) {
    startTime := time.Now()

    result := &ExecutionResult{
        MachineID: machine.Name,
        Command:   command,
        Status:    ExecutionStatusRunning,
        Metadata:  make(map[string]interface{}),
    }

    // Create SSH configuration
    sshConfig := &sshTypes.SSHConfig{
        Host:     machine.Host,
        Port:     machine.Port,
        Username: machine.User,
        Timeout:  options.Timeout,
    }

    // Execute with retry logic
    var lastErr error
    for attempt := 0; attempt <= options.RetryCount; attempt++ {
        if attempt > 0 {
            e.logger.Debug("Retrying command execution",
                logging.String("machine", machine.Name),
                logging.Int("attempt", attempt+1))
            time.Sleep(options.RetryDelay)
        }

        execResult, err := e.executeWithSSH(ctx, sshConfig, command)
        if err != nil {
            lastErr = err
            continue
        }

        // Command executed successfully
        result.ExitCode = execResult.ExitCode
        result.Stdout = execResult.Stdout
        result.Stderr = execResult.Stderr
        result.Status = ExecutionStatusCompleted
        break
    }

    if result.Status != ExecutionStatusCompleted {
        result.Status = ExecutionStatusFailed
        result.Error = lastErr.Error()
    }

    result.ExecutionTime = time.Since(startTime)
    result.Metadata["attempts"] = options.RetryCount + 1

    return result, nil
}
```

#### 5.2 SSH Execution Method
```go
// executeWithSSH executes a command via SSH
func (e *CommandExecutor) executeWithSSH(ctx context.Context, sshConfig *sshTypes.SSHConfig, command string) (*sshTypes.CommandResult, error) {
    // Create SSH connection
    connection, err := e.sshManager.Connect(sshConfig.Host, sshConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to connect: %w", err)
    }
    defer e.sshManager.CloseConnection(connection)

    // Execute command
    result, err := e.sshManager.ExecuteCommand(connection, command)
    if err != nil {
        return nil, fmt.Errorf("command execution failed: %w", err)
    }

    return result, nil
}
```

### Step 6: Implement Progress Reporting

#### 6.1 Progress Display
```go
// showProgress displays execution progress
func (e *CommandExecutor) showProgress(result *ExecutionResult, completed, total int) {
    if result != nil {
        status := "✓"
        if result.Status == ExecutionStatusFailed {
            status = "✗"
        }
        fmt.Printf("[%s] %s: %s (%v)\n", status, result.MachineID, result.Status, result.ExecutionTime)
    } else {
        fmt.Printf("Progress: %d/%d machines\n", completed, total)
    }
}
```

### Step 7: Implement Result Display

#### 7.1 Result Display Method
```go
// displayResults displays execution results
func (e *CommandExecutor) displayResults(results []*ExecutionResult, options *ExecutionOptions) error {
    switch options.OutputFormat {
    case "json":
        return e.displayJSONResults(results)
    case "table":
        return e.displayTableResults(results)
    case "text":
    default:
        return e.displayTextResults(results)
    }
    return nil
}
```

#### 7.2 Text Output Display
```go
// displayTextResults displays results in text format
func (e *CommandExecutor) displayTextResults(results []*ExecutionResult) error {
    fmt.Println("Command Execution Results:")
    fmt.Println("==========================")

    successCount := 0
    failureCount := 0

    for _, result := range results {
        fmt.Printf("\nMachine: %s\n", result.MachineID)
        fmt.Printf("Command: %s\n", result.Command)
        fmt.Printf("Status: %s\n", result.Status)
        fmt.Printf("Exit Code: %d\n", result.ExitCode)
        fmt.Printf("Execution Time: %v\n", result.ExecutionTime)

        if result.Stdout != "" {
            fmt.Printf("Stdout:\n%s\n", result.Stdout)
        }

        if result.Stderr != "" {
            fmt.Printf("Stderr:\n%s\n", result.Stderr)
        }

        if result.Error != "" {
            fmt.Printf("Error: %s\n", result.Error)
        }

        if result.Status == ExecutionStatusCompleted && result.ExitCode == 0 {
            successCount++
        } else {
            failureCount++
        }
    }

    fmt.Printf("\nSummary: %d successful, %d failed\n", successCount, failureCount)
    return nil
}
```

#### 7.3 JSON Output Display
```go
// displayJSONResults displays results in JSON format
func (e *CommandExecutor) displayJSONResults(results []*ExecutionResult) error {
    // This would use the existing JSON output formatting
    // For now, we'll use a simple implementation
    output := map[string]interface{}{
        "results": results,
        "summary": map[string]interface{}{
            "total":     len(results),
            "successful": 0,
            "failed":     0,
        },
    }

    // Calculate summary
    for _, result := range results {
        if result.Status == ExecutionStatusCompleted && result.ExitCode == 0 {
            output["summary"].(map[string]interface{})["successful"] = output["summary"].(map[string]interface{})["successful"].(int) + 1
        } else {
            output["summary"].(map[string]interface{})["failed"] = output["summary"].(map[string]interface{})["failed"].(int) + 1
        }
    }

    // Use existing JSON output formatter
    return nil
}
```

#### 7.4 Table Output Display
```go
// displayTableResults displays results in table format
func (e *CommandExecutor) displayTableResults(results []*ExecutionResult) error {
    fmt.Println("Command Execution Results:")
    fmt.Println("==========================")
    fmt.Printf("%-20s %-10s %-8s %-12s %s\n", "Machine", "Status", "Exit", "Time", "Error")
    fmt.Println(strings.Repeat("-", 80))

    for _, result := range results {
        status := string(result.Status)
        if result.Status == ExecutionStatusCompleted && result.ExitCode == 0 {
            status = "SUCCESS"
        } else if result.Status == ExecutionStatusFailed {
            status = "FAILED"
        }

        exitCode := fmt.Sprintf("%d", result.ExitCode)
        executionTime := result.ExecutionTime.String()
        errorMsg := ""
        if result.Error != "" {
            errorMsg = result.Error
        }

        fmt.Printf("%-20s %-10s %-8s %-12s %s\n", result.MachineID, status, exitCode, executionTime, errorMsg)
    }

    return nil
}
```







## Configuration Options

### Supported Options
- **Parallel**: Enable/disable parallel execution
- **Timeout**: Command execution timeout
- **DryRun**: Show commands without execution
- **OutputFormat**: Text, JSON, or table output
- **ShowProgress**: Enable/disable progress reporting
- **StopOnError**: Stop execution on first error
- **RetryCount**: Number of retry attempts
- **RetryDelay**: Delay between retries

## Dependencies

### Internal Dependencies
- `spooky/internal/cli/types`
- `spooky/internal/ssh`
- `spooky/internal/ssh/types`
- `spooky/internal/logging`

### External Dependencies
- `github.com/spf13/cobra`
- `context` (standard library)
- `strings` (standard library)
- `time` (standard library)


7. **Testing**: Comprehensive test coverage
8. **Documentation**: Clear code documentation

## Implementation Order

1. Create command execution CLI structure
2. Implement CLI command definition
3. Add command execution method
4. Implement command preparation
5. Add parallel execution
6. Implement sequential execution
7. Add machine execution
8. Implement progress reporting
9. Add result display
10. Write comprehensive tests
11. Performance optimization
12. Documentation and cleanup


