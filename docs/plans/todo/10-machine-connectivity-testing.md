# Implementation Plan: Machine Connectivity Testing

## Overview
Implement machine connectivity testing functionality that can test SSH connectivity, validate machine configurations, perform health checks, and provide detailed connectivity reports.

## Task Details
- **Task ID**: 6.4
- **Priority**: Medium
- **File**: `internal/cli/commands/machines.go`
- **Function**: `testConnectivity`


## Current State Analysis

### Existing Patterns
1. **CLI Structure**: CLI commands use `spf13/cobra` framework with consistent patterns
2. **SSH Manager**: Existing SSH manager provides connection and command execution capabilities
3. **Machine Types**: Machine configurations defined in `internal/machines/types/`
4. **Error Handling**: Consistent error wrapping with context
5. **Output Formatting**: Structured output with JSON and text formats
6. **Parallel Operations**: Support for parallel operations across machines

### Existing Implementation Examples
- **CLI Commands**: `internal/cli/commands/machines.go` provides machine-related commands
- **SSH Manager**: `internal/ssh/manager.go` provides SSH operations
- **Machine Types**: `internal/machines/types/structures.go` defines machine structures
- **Output Formatting**: `internal/cli/help/renderer.go` provides output formatting

## Implementation Requirements

### Interface Compliance
The machine connectivity testing must:
1. **Test SSH connectivity** to target machines
2. **Validate machine configurations** and parameters
3. **Perform health checks** and basic system tests
4. **Support parallel testing** across multiple machines
5. **Provide detailed connectivity reports** and diagnostics
6. **Handle different authentication methods** (key-based, password)
7. **Support timeout and retry logic** for reliable testing
8. **Generate comprehensive test reports** with status and metrics

### Required Dependencies
- SSH manager for connectivity testing
- CLI framework for command handling
- Output formatting system
- Parallel execution system

## Detailed Implementation Plan

### Step 1: Create Machine Connectivity Testing Structure

**File**: `internal/cli/commands/connectivity.go`

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

// ConnectivityTester handles machine connectivity testing
type ConnectivityTester struct {
    sshManager ssh.Manager
    logger     logging.Logger
}

// ConnectivityResult represents the result of connectivity testing
type ConnectivityResult struct {
    MachineID       string                 `json:"machine_id"`
    Host            string                 `json:"host"`
    Port            int                    `json:"port"`
    Username        string                 `json:"username"`
    Status          ConnectivityStatus     `json:"status"`
    ConnectionTime  time.Duration          `json:"connection_time"`
    ResponseTime    time.Duration          `json:"response_time"`
    SSHVersion      string                 `json:"ssh_version"`
    OSInfo          *OSInfo                `json:"os_info"`
    HealthChecks    *HealthCheckResult     `json:"health_checks"`
    Errors          []string               `json:"errors"`
    Warnings        []string               `json:"warnings"`
    Metadata        map[string]interface{} `json:"metadata"`
}

// ConnectivityStatus represents the connectivity status
type ConnectivityStatus string

const (
    ConnectivityStatusUnknown   ConnectivityStatus = "unknown"
    ConnectivityStatusReachable ConnectivityStatus = "reachable"
    ConnectivityStatusUnreachable ConnectivityStatus = "unreachable"
    ConnectivityStatusTimeout   ConnectivityStatus = "timeout"
    ConnectivityStatusAuthFailed ConnectivityStatus = "auth_failed"
    ConnectivityStatusError     ConnectivityStatus = "error"
)

// OSInfo represents operating system information
type OSInfo struct {
    Name        string `json:"name"`
    Version     string `json:"version"`
    Architecture string `json:"architecture"`
    Kernel      string `json:"kernel"`
    Hostname    string `json:"hostname"`
    Uptime      string `json:"uptime"`
}

// HealthCheckResult represents health check results
type HealthCheckResult struct {
    DiskSpace   *DiskSpaceCheck   `json:"disk_space"`
    MemoryUsage *MemoryUsageCheck `json:"memory_usage"`
    LoadAverage *LoadAverageCheck `json:"load_average"`
    Network     *NetworkCheck     `json:"network"`
    Services    *ServicesCheck    `json:"services"`
}

// DiskSpaceCheck represents disk space check results
type DiskSpaceCheck struct {
    Available   int64  `json:"available"`
    Total       int64  `json:"total"`
    Used        int64  `json:"used"`
    UsagePercent float64 `json:"usage_percent"`
    Status      string `json:"status"`
}

// MemoryUsageCheck represents memory usage check results
type MemoryUsageCheck struct {
    Available   int64  `json:"available"`
    Total       int64  `json:"total"`
    Used        int64  `json:"used"`
    UsagePercent float64 `json:"usage_percent"`
    Status      string `json:"status"`
}

// LoadAverageCheck represents load average check results
type LoadAverageCheck struct {
    OneMinute   float64 `json:"one_minute"`
    FiveMinutes float64 `json:"five_minutes"`
    FifteenMinutes float64 `json:"fifteen_minutes"`
    Status      string `json:"status"`
}

// NetworkCheck represents network connectivity check results
type NetworkCheck struct {
    Interfaces  []NetworkInterface `json:"interfaces"`
    DNS         *DNSCheck          `json:"dns"`
    Gateway     *GatewayCheck      `json:"gateway"`
    Status      string             `json:"status"`
}

// NetworkInterface represents a network interface
type NetworkInterface struct {
    Name        string `json:"name"`
    IPAddress   string `json:"ip_address"`
    MACAddress  string `json:"mac_address"`
    Status      string `json:"status"`
}

// DNSCheck represents DNS check results
type DNSCheck struct {
    Resolvable  bool   `json:"resolvable"`
    ResponseTime time.Duration `json:"response_time"`
    Status      string `json:"status"`
}

// GatewayCheck represents gateway check results
type GatewayCheck struct {
    Reachable   bool   `json:"reachable"`
    ResponseTime time.Duration `json:"response_time"`
    Status      string `json:"status"`
}

// ServicesCheck represents services check results
type ServicesCheck struct {
    SSH         *ServiceStatus `json:"ssh"`
    Systemd     *ServiceStatus `json:"systemd"`
    Status      string         `json:"status"`
}

// ServiceStatus represents a service status
type ServiceStatus struct {
    Running     bool   `json:"running"`
    Enabled     bool   `json:"enabled"`
    Status      string `json:"status"`
}

// TestingOptions represents options for connectivity testing
type TestingOptions struct {
    Parallel      bool                   `json:"parallel"`
    Timeout       time.Duration          `json:"timeout"`
    RetryCount    int                    `json:"retry_count"`
    RetryDelay    time.Duration          `json:"retry_delay"`
    HealthChecks  bool                   `json:"health_checks"`
    OutputFormat  string                 `json:"output_format"`
    ShowProgress  bool                   `json:"show_progress"`
    Detailed      bool                   `json:"detailed"`
    StopOnError   bool                   `json:"stop_on_error"`
}

// NewConnectivityTester creates a new connectivity tester
func NewConnectivityTester(sshManager ssh.Manager, logger logging.Logger) *ConnectivityTester {
    return &ConnectivityTester{
        sshManager: sshManager,
        logger:     logger,
    }
}
```

### Step 2: Implement Connectivity Testing CLI Command

#### 2.1 Main CLI Command
```go
// NewTestConnectivityCommand creates the connectivity testing CLI command
func NewTestConnectivityCommand(tester *ConnectivityTester) *cobra.Command {
    var (
        parallel     bool
        timeout      time.Duration
        retryCount   int
        retryDelay   time.Duration
        healthChecks bool
        outputFormat string
        showProgress bool
        detailed     bool
        stopOnError  bool
    )

    cmd := &cobra.Command{
        Use:   "test-connectivity [project]",
        Short: "Test connectivity to target machines",
        Long: `Test connectivity to target machines in a project.

Examples:
  spooky test-connectivity myproject
  spooky test-connectivity myproject --parallel --health-checks
  spooky test-connectivity myproject --timeout 30s --detailed
  spooky test-connectivity myproject --output json`,
        Args: cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            projectPath := args[0]

            options := &TestingOptions{
                Parallel:     parallel,
                Timeout:      timeout,
                RetryCount:   retryCount,
                RetryDelay:   retryDelay,
                HealthChecks: healthChecks,
                OutputFormat: outputFormat,
                ShowProgress: showProgress,
                Detailed:     detailed,
                StopOnError:  stopOnError,
            }

            return tester.TestConnectivity(cmd.Context(), projectPath, options)
        },
    }

    // Add flags
    cmd.Flags().BoolVarP(&parallel, "parallel", "p", false, "Test machines in parallel")
    cmd.Flags().DurationVarP(&timeout, "timeout", "t", 30*time.Second, "Connection timeout")
    cmd.Flags().IntVarP(&retryCount, "retry", "r", 1, "Number of retry attempts")
    cmd.Flags().DurationVarP(&retryDelay, "retry-delay", "", 5*time.Second, "Delay between retries")
    cmd.Flags().BoolVarP(&healthChecks, "health-checks", "h", false, "Perform health checks")
    cmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text, json, table)")
    cmd.Flags().BoolVarP(&showProgress, "progress", "", true, "Show testing progress")
    cmd.Flags().BoolVarP(&detailed, "detailed", "d", false, "Show detailed information")
    cmd.Flags().BoolVarP(&stopOnError, "stop-on-error", "", false, "Stop testing on first error")

    return cmd
}
```

#### 2.2 Connectivity Testing Method
```go
// TestConnectivity tests connectivity to target machines
func (t *ConnectivityTester) TestConnectivity(ctx context.Context, projectPath string, options *TestingOptions) error {
    t.logger.Info("Testing machine connectivity",
        logging.String("project", projectPath),
        logging.Bool("parallel", options.Parallel),
        logging.Bool("health_checks", options.HealthChecks))

    // Load project configuration
    project, err := t.loadProject(projectPath)
    if err != nil {
        return fmt.Errorf("failed to load project: %w", err)
    }

    // Load machine inventory
    machines, err := t.loadMachines(projectPath)
    if err != nil {
        return fmt.Errorf("failed to load machines: %w", err)
    }

    if len(machines) == 0 {
        return fmt.Errorf("no machines found in project")
    }

    // Test connectivity based on parallel flag
    var results []*ConnectivityResult
    if options.Parallel {
        results, err = t.testParallel(ctx, machines, options)
    } else {
        results, err = t.testSequential(ctx, machines, options)
    }

    if err != nil {
        return fmt.Errorf("connectivity testing failed: %w", err)
    }

    // Format and display results
    return t.displayResults(results, options)
}
```

### Step 3: Implement Parallel and Sequential Testing

#### 3.1 Parallel Testing Method
```go
// testParallel tests connectivity in parallel across machines
func (t *ConnectivityTester) testParallel(ctx context.Context, machines []*types.Machine, options *TestingOptions) ([]*ConnectivityResult, error) {
    t.logger.Debug("Testing connectivity in parallel",
        logging.Int("machine_count", len(machines)))

    results := make([]*ConnectivityResult, 0, len(machines))
    resultChan := make(chan *ConnectivityResult, len(machines))
    errorChan := make(chan error, len(machines))

    // Start testing for each machine
    for _, machine := range machines {
        go func(m *types.Machine) {
            result, err := t.testMachine(ctx, m, options)
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
                t.showProgress(result, completed, len(machines))
            }

            // Check if we should stop on error
            if options.StopOnError && result.Status == ConnectivityStatusUnreachable {
                t.logger.Warn("Stopping testing due to unreachable machine",
                    logging.String("machine", result.MachineID))
                return results, fmt.Errorf("testing stopped due to unreachable machine %s", result.MachineID)
            }

        case err := <-errorChan:
            t.logger.Error("Testing error",
                logging.Error(err))
            return results, err

        case <-ctx.Done():
            t.logger.Warn("Testing cancelled")
            return results, ctx.Err()
        }
    }

    return results, nil
}
```

#### 3.2 Sequential Testing Method
```go
// testSequential tests connectivity sequentially across machines
func (t *ConnectivityTester) testSequential(ctx context.Context, machines []*types.Machine, options *TestingOptions) ([]*ConnectivityResult, error) {
    t.logger.Debug("Testing connectivity sequentially",
        logging.Int("machine_count", len(machines)))

    results := make([]*ConnectivityResult, 0, len(machines))

    for i, machine := range machines {
        if options.ShowProgress {
            t.showProgress(nil, i, len(machines))
        }

        result, err := t.testMachine(ctx, machine, options)
        if err != nil {
            t.logger.Error("Testing failed",
                logging.String("machine", machine.Name),
                logging.Error(err))
            return results, fmt.Errorf("machine %s: %w", machine.Name, err)
        }

        results = append(results, result)

        // Check if we should stop on error
        if options.StopOnError && result.Status == ConnectivityStatusUnreachable {
            t.logger.Warn("Stopping testing due to unreachable machine",
                logging.String("machine", result.MachineID))
            return results, fmt.Errorf("testing stopped due to unreachable machine %s", result.MachineID)
        }

        // Check for cancellation
        select {
        case <-ctx.Done():
            t.logger.Warn("Testing cancelled")
            return results, ctx.Err()
        default:
        }
    }

    return results, nil
}
```

### Step 4: Implement Machine Testing

#### 4.1 Machine Testing Method
```go
// testMachine tests connectivity to a specific machine
func (t *ConnectivityTester) testMachine(ctx context.Context, machine *types.Machine, options *TestingOptions) (*ConnectivityResult, error) {
    startTime := time.Now()

    result := &ConnectivityResult{
        MachineID: machine.Name,
        Host:      machine.Host,
        Port:      machine.Port,
        Username:  machine.User,
        Status:    ConnectivityStatusUnknown,
        Metadata:  make(map[string]interface{}),
    }

    // Create SSH configuration
    sshConfig := &sshTypes.SSHConfig{
        Host:     machine.Host,
        Port:     machine.Port,
        Username: machine.User,
        Timeout:  options.Timeout,
    }

    // Test connectivity with retry logic
    var lastErr error
    for attempt := 0; attempt <= options.RetryCount; attempt++ {
        if attempt > 0 {
            t.logger.Debug("Retrying connectivity test",
                logging.String("machine", machine.Name),
                logging.Int("attempt", attempt+1))
            time.Sleep(options.RetryDelay)
        }

        testResult, err := t.testConnectivity(ctx, sshConfig, options)
        if err != nil {
            lastErr = err
            continue
        }

        // Test successful
        result.Status = testResult.Status
        result.ConnectionTime = testResult.ConnectionTime
        result.ResponseTime = testResult.ResponseTime
        result.SSHVersion = testResult.SSHVersion
        result.OSInfo = testResult.OSInfo
        result.HealthChecks = testResult.HealthChecks
        break
    }

    if result.Status == ConnectivityStatusUnknown {
        result.Status = ConnectivityStatusUnreachable
        result.Errors = append(result.Errors, lastErr.Error())
    }

    result.Metadata["attempts"] = options.RetryCount + 1

    return result, nil
}
```

#### 4.2 Connectivity Testing Method
```go
// testConnectivity performs the actual connectivity test
func (t *ConnectivityTester) testConnectivity(ctx context.Context, sshConfig *sshTypes.SSHConfig, options *TestingOptions) (*ConnectivityResult, error) {
    connectionStart := time.Now()

    // Test SSH connection
    connection, err := t.sshManager.Connect(sshConfig.Host, sshConfig)
    if err != nil {
        return nil, fmt.Errorf("SSH connection failed: %w", err)
    }
    defer t.sshManager.CloseConnection(connection)

    connectionTime := time.Since(connectionStart)

    // Test basic command execution
    responseStart := time.Now()
    result, err := t.sshManager.ExecuteCommand(connection, "echo 'connectivity test successful'")
    if err != nil {
        return nil, fmt.Errorf("command execution failed: %w", err)
    }

    responseTime := time.Since(responseStart)

    // Gather system information
    osInfo, err := t.gatherOSInfo(connection)
    if err != nil {
        t.logger.Warn("Failed to gather OS info",
            logging.String("machine", sshConfig.Host),
            logging.Error(err))
    }

    // Perform health checks if requested
    var healthChecks *HealthCheckResult
    if options.HealthChecks {
        healthChecks, err = t.performHealthChecks(connection)
        if err != nil {
            t.logger.Warn("Failed to perform health checks",
                logging.String("machine", sshConfig.Host),
                logging.Error(err))
        }
    }

    return &ConnectivityResult{
        Status:         ConnectivityStatusReachable,
        ConnectionTime: connectionTime,
        ResponseTime:   responseTime,
        SSHVersion:     "OpenSSH", // This would be extracted from SSH handshake
        OSInfo:         osInfo,
        HealthChecks:   healthChecks,
    }, nil
}
```

### Step 5: Implement System Information Gathering

#### 5.1 OS Information Gathering
```go
// gatherOSInfo gathers operating system information
func (t *ConnectivityTester) gatherOSInfo(connection sshTypes.Connection) (*OSInfo, error) {
    osInfo := &OSInfo{}

    // Get OS name and version
    result, err := t.sshManager.ExecuteCommand(connection, "cat /etc/os-release")
    if err == nil && result.ExitCode == 0 {
        lines := strings.Split(result.Stdout, "\n")
        for _, line := range lines {
            if strings.HasPrefix(line, "PRETTY_NAME=") {
                osInfo.Name = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
            } else if strings.HasPrefix(line, "VERSION=") {
                osInfo.Version = strings.Trim(strings.TrimPrefix(line, "VERSION="), "\"")
            }
        }
    }

    // Get architecture
    result, err = t.sshManager.ExecuteCommand(connection, "uname -m")
    if err == nil && result.ExitCode == 0 {
        osInfo.Architecture = strings.TrimSpace(result.Stdout)
    }

    // Get kernel version
    result, err = t.sshManager.ExecuteCommand(connection, "uname -r")
    if err == nil && result.ExitCode == 0 {
        osInfo.Kernel = strings.TrimSpace(result.Stdout)
    }

    // Get hostname
    result, err = t.sshManager.ExecuteCommand(connection, "hostname")
    if err == nil && result.ExitCode == 0 {
        osInfo.Hostname = strings.TrimSpace(result.Stdout)
    }

    // Get uptime
    result, err = t.sshManager.ExecuteCommand(connection, "uptime -p")
    if err == nil && result.ExitCode == 0 {
        osInfo.Uptime = strings.TrimSpace(result.Stdout)
    }

    return osInfo, nil
}
```

#### 5.2 Health Checks
```go
// performHealthChecks performs health checks on the machine
func (t *ConnectivityTester) performHealthChecks(connection sshTypes.Connection) (*HealthCheckResult, error) {
    healthChecks := &HealthCheckResult{}

    // Disk space check
    diskSpace, err := t.checkDiskSpace(connection)
    if err == nil {
        healthChecks.DiskSpace = diskSpace
    }

    // Memory usage check
    memoryUsage, err := t.checkMemoryUsage(connection)
    if err == nil {
        healthChecks.MemoryUsage = memoryUsage
    }

    // Load average check
    loadAverage, err := t.checkLoadAverage(connection)
    if err == nil {
        healthChecks.LoadAverage = loadAverage
    }

    // Network check
    network, err := t.checkNetwork(connection)
    if err == nil {
        healthChecks.Network = network
    }

    // Services check
    services, err := t.checkServices(connection)
    if err == nil {
        healthChecks.Services = services
    }

    return healthChecks, nil
}
```

### Step 6: Implement Health Check Methods

#### 6.1 Disk Space Check
```go
// checkDiskSpace checks disk space usage
func (t *ConnectivityTester) checkDiskSpace(connection sshTypes.Connection) (*DiskSpaceCheck, error) {
    result, err := t.sshManager.ExecuteCommand(connection, "df -h / | tail -1")
    if err != nil || result.ExitCode != 0 {
        return nil, fmt.Errorf("disk space check failed")
    }

    // Parse df output
    fields := strings.Fields(result.Stdout)
    if len(fields) < 5 {
        return nil, fmt.Errorf("invalid df output format")
    }

    // Extract values (this is a simplified parser)
    total := int64(0) // Would parse from fields[1]
    used := int64(0)  // Would parse from fields[2]
    available := int64(0) // Would parse from fields[3]

    usagePercent := float64(used) / float64(total) * 100

    status := "OK"
    if usagePercent > 90 {
        status = "CRITICAL"
    } else if usagePercent > 80 {
        status = "WARNING"
    }

    return &DiskSpaceCheck{
        Available:    available,
        Total:        total,
        Used:         used,
        UsagePercent: usagePercent,
        Status:       status,
    }, nil
}
```

#### 6.2 Memory Usage Check
```go
// checkMemoryUsage checks memory usage
func (t *ConnectivityTester) checkMemoryUsage(connection sshTypes.Connection) (*MemoryUsageCheck, error) {
    result, err := t.sshManager.ExecuteCommand(connection, "free -m | grep Mem")
    if err != nil || result.ExitCode != 0 {
        return nil, fmt.Errorf("memory usage check failed")
    }

    // Parse free output
    fields := strings.Fields(result.Stdout)
    if len(fields) < 3 {
        return nil, fmt.Errorf("invalid free output format")
    }

    total := int64(0)    // Would parse from fields[1]
    used := int64(0)     // Would parse from fields[2]
    available := int64(0) // Would parse from fields[6]

    usagePercent := float64(used) / float64(total) * 100

    status := "OK"
    if usagePercent > 90 {
        status = "CRITICAL"
    } else if usagePercent > 80 {
        status = "WARNING"
    }

    return &MemoryUsageCheck{
        Available:    available,
        Total:        total,
        Used:         used,
        UsagePercent: usagePercent,
        Status:       status,
    }, nil
}
```

### Step 7: Implement Result Display

#### 7.1 Result Display Method
```go
// displayResults displays connectivity test results
func (t *ConnectivityTester) displayResults(results []*ConnectivityResult, options *TestingOptions) error {
    switch options.OutputFormat {
    case "json":
        return t.displayJSONResults(results)
    case "table":
        return t.displayTableResults(results)
    case "text":
    default:
        return t.displayTextResults(results, options)
    }
    return nil
}
```

#### 7.2 Text Output Display
```go
// displayTextResults displays results in text format
func (t *ConnectivityTester) displayTextResults(results []*ConnectivityResult, options *TestingOptions) error {
    fmt.Println("Machine Connectivity Test Results:")
    fmt.Println("==================================")

    reachableCount := 0
    unreachableCount := 0

    for _, result := range results {
        fmt.Printf("\nMachine: %s (%s:%d)\n", result.MachineID, result.Host, result.Port)
        fmt.Printf("Status: %s\n", result.Status)
        fmt.Printf("Connection Time: %v\n", result.ConnectionTime)
        fmt.Printf("Response Time: %v\n", result.ResponseTime)

        if result.OSInfo != nil {
            fmt.Printf("OS: %s\n", result.OSInfo.Name)
            fmt.Printf("Hostname: %s\n", result.OSInfo.Hostname)
            fmt.Printf("Kernel: %s\n", result.OSInfo.Kernel)
        }

        if options.HealthChecks && result.HealthChecks != nil {
            fmt.Printf("Health Checks:\n")
            if result.HealthChecks.DiskSpace != nil {
                fmt.Printf("  Disk Space: %s (%.1f%% used)\n", 
                    result.HealthChecks.DiskSpace.Status, 
                    result.HealthChecks.DiskSpace.UsagePercent)
            }
            if result.HealthChecks.MemoryUsage != nil {
                fmt.Printf("  Memory: %s (%.1f%% used)\n", 
                    result.HealthChecks.MemoryUsage.Status, 
                    result.HealthChecks.MemoryUsage.UsagePercent)
            }
        }

        if len(result.Errors) > 0 {
            fmt.Printf("Errors:\n")
            for _, err := range result.Errors {
                fmt.Printf("  - %s\n", err)
            }
        }

        if len(result.Warnings) > 0 {
            fmt.Printf("Warnings:\n")
            for _, warning := range result.Warnings {
                fmt.Printf("  - %s\n", warning)
            }
        }

        if result.Status == ConnectivityStatusReachable {
            reachableCount++
        } else {
            unreachableCount++
        }
    }

    fmt.Printf("\nSummary: %d reachable, %d unreachable\n", reachableCount, unreachableCount)
    return nil
}
```







## Configuration Options

### Supported Options
- **Parallel**: Enable/disable parallel testing
- **Timeout**: Connection timeout
- **RetryCount**: Number of retry attempts
- **RetryDelay**: Delay between retries
- **HealthChecks**: Enable/disable health checks
- **OutputFormat**: Text, JSON, or table output
- **ShowProgress**: Enable/disable progress reporting
- **Detailed**: Show detailed information
- **StopOnError**: Stop testing on first error

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



## Implementation Order

1. Create connectivity testing structure
2. Implement CLI command definition
3. Add connectivity testing method
4. Implement parallel testing
5. Add sequential testing
6. Implement machine testing
7. Add system information gathering
8. Implement health checks
9. Add result display
10. Write comprehensive tests
11. Performance optimization
12. Documentation and cleanup


