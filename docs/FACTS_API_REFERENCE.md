# Facts System API Reference

## Overview

This document provides a comprehensive API reference for the spooky facts system. It covers all interfaces, types, methods, and implementation details for developers working with the facts system.

**Status: Partially Implemented** - The facts system has basic functionality but SSH-based fact collection has known issues that need to be addressed.

## Core Interfaces

### FactsIntegration Interface

The `FactsIntegration` interface provides the primary entry point for facts operations:

```go
type FactsIntegration interface {
    // CollectFacts collects facts from the given source
    CollectFacts(ctx context.Context, source string) (interface{}, error)
    
    // ExportFacts exports facts to the specified format and output
    ExportFacts(ctx context.Context, projectPath string, format string, outputPath string) error
    
    // ValidateFacts validates fact collection and storage
    ValidateFacts(ctx context.Context, projectPath string) (*ValidationResult, error)
}
```

**Implementation Status**: ⚠️ **Partially Implemented** - Basic functionality exists but SSH-based collection has issues

### FactsManager Interface

The `FactsManager` interface provides fact collection and management:

```go
type FactsManager interface {
    // CollectFacts collects facts from the given machine
    CollectFacts(ctx context.Context, machine *spookytypes.Machine) (*spookytypesfacts.FactCollection, error)
    
    // ExportFacts exports facts to the specified format
    ExportFacts(ctx context.Context, facts *spookytypesfacts.FactCollection, format string, outputPath string) error
    
    // ValidateFacts validates fact collection
    ValidateFacts(ctx context.Context, facts *spookytypesfacts.FactCollection) (*ValidationResult, error)
}
```

**Implementation Status**: ⚠️ **Partially Implemented** - Basic collection exists but SSH integration has issues

## Current Implementation Status

### ✅ Working Components

1. **Basic Fact Collection**: System fact collection using SSH commands
2. **Memory Storage**: In-memory fact storage during export operations
3. **Export Functionality**: Facts export to JSON and HCL formats
4. **CLI Integration**: `spooky facts export` command with filtering options
5. **Machine Integration**: Facts collection from project machine inventory
6. **Basic Validation**: Fact collection validation and error handling

### ❌ Known Issues

1. **SSH-based Fact Collection**: Has implementation issues that prevent proper remote collection
2. **Remote Facts Reading**: Cannot properly read `/etc/spooky/facts.*` files from remote machines
3. **Parallel Processing**: Sequential collection only, no multi-machine parallel processing
4. **SSH Integration**: SSH integration has bugs that prevent proper remote machine access

## Type Definitions

### Fact Collection Types

```go
// FactCollection represents a complete fact collection for a machine
type FactCollection struct {
    MachineID   string                 `json:"machine_id" hcl:"machine_id"`
    CollectedAt time.Time              `json:"collected_at" hcl:"collected_at"`
    Facts       *Facts                 `json:"facts" hcl:"facts"`
    Metadata    map[string]interface{} `json:"metadata" hcl:"metadata"`
}

// Facts represents the complete facts structure for a machine
type Facts struct {
    // System facts (from SSH commands - user level)
    System *SystemFacts `json:"system" hcl:"system"`
    
    // Collector facts (from spooky-collector binary - comprehensive gopsutil coverage)
    Collector *CollectorFacts `json:"collector,omitempty" hcl:"collector,optional"`
    
    // Custom facts (user-defined from /etc/spooky/custom.hcl)
    Custom map[string]interface{} `json:"custom,omitempty" hcl:"custom,optional"`
}
```

### System Facts Types

```go
// SystemFacts represents system-level facts from SSH commands
type SystemFacts struct {
    // Operating system facts
    OS *OSFacts `json:"os" hcl:"os"`
    
    // Hardware facts
    Hardware *HardwareFacts `json:"hardware" hcl:"hardware"`
    
    // Network facts
    Network *NetworkFacts `json:"network" hcl:"network"`
    
    // Load average facts
    LoadAverage *LoadAverageFacts `json:"load_average,omitempty" hcl:"load_average,optional"`
    
    // Process information
    Processes *ProcessFacts `json:"processes,omitempty" hcl:"processes,optional"`
}

// OSFacts represents operating system facts
type OSFacts struct {
    Name        string `json:"name" hcl:"name"`
    Version     string `json:"version" hcl:"version"`
    Architecture string `json:"architecture" hcl:"architecture"`
    Kernel      string `json:"kernel" hcl:"kernel"`
    Hostname    string `json:"hostname" hcl:"hostname"`
    Uptime      int64  `json:"uptime" hcl:"uptime"`
}

// HardwareFacts represents hardware facts
type HardwareFacts struct {
    CPUCount    int   `json:"cpu_count" hcl:"cpu_count"`
    MemoryTotal int64 `json:"memory_total" hcl:"memory_total"`
    DiskTotal   int64 `json:"disk_total" hcl:"disk_total"`
}

// NetworkFacts represents network facts
type NetworkFacts struct {
    PrimaryIP   string   `json:"primary_ip" hcl:"primary_ip"`
    Interfaces  []string `json:"interfaces" hcl:"interfaces"`
    Hostname    string   `json:"hostname" hcl:"hostname"`
}
```

## Implementation Details

### SystemFactCollector

The `SystemFactCollector` provides SSH-based fact collection:

```go
// SystemFactCollector collects system facts using SSH commands
type SystemFactCollector struct {
    name       string
    sshManager spookyinterfaces.SSHManager
}

// Collect collects facts from the given machine
func (c *SystemFactCollector) Collect(ctx context.Context, machine *spookytypes.Machine) (*spookytypesfacts.FactCollection, error) {
    // Get machine ID
    machineID, err := c.getMachineID(machine)
    if err != nil {
        return nil, fmt.Errorf("failed to get machine ID: %w", err)
    }
    
    // Collect system facts
    facts := &spookytypesfacts.Facts{
        System: &spookytypesfacts.SystemFacts{},
    }
    
    // Collect system facts via SSH
    systemFacts, err := c.collectSystemFacts(ctx, machine)
    if err != nil {
        return nil, fmt.Errorf("failed to collect system facts: %w", err)
    }
    facts.System = systemFacts
    
    // Create fact collection
    collection := &spookytypesfacts.FactCollection{
        MachineID:   machineID,
        CollectedAt: time.Now(),
        Facts:       facts,
        Metadata: map[string]interface{}{
            "collector": c.name,
            "machine":   machine.Hostname,
        },
    }
    
    return collection, nil
}
```

### SSH Command Execution

The collector uses SSH commands to gather system information:

```go
// runSSHCommand runs a command via SSH on the target machine
func (c *SystemFactCollector) runSSHCommand(machine *spookytypes.Machine, command string) (string, error) {
    ctx := context.Background()
    
    // Create connection request with actual machine configuration
    connectionRequest := &spookytypes.ConnectionRequest{
        Host:       machine.Hostname,
        Port:       machine.Port,
        User:       machine.User,
        Password:   machine.Password,
        KeyFile:    machine.KeyFile,
        Passphrase: machine.Passphrase,
        AuthMethod: spookytypesssh.AuthMethodPublicKey,
        Timeout:    30 * time.Second,
    }
    
    // Establish connection
    connectionResult, err := c.sshManager.Connect(ctx, connectionRequest)
    if err != nil {
        return "", fmt.Errorf("failed to establish SSH connection to %s: %w", machine.Hostname, err)
    }
    
    // Create session and run command
    session, err := c.sshManager.CreateSession(ctx, connectionResult.Connection)
    if err != nil {
        return "", fmt.Errorf("failed to create SSH session on %s: %w", machine.Hostname, err)
    }
    
    sshCommand := &spookytypes.SSHCommand{
        Command: command,
        Timeout: 30 * time.Second,
    }
    
    commandResult, err := c.sshManager.RunCommand(ctx, session, sshCommand)
    if err != nil {
        return "", fmt.Errorf("failed to run SSH command on %s: %w", machine.Hostname, err)
    }
    
    if !commandResult.Success {
        return "", fmt.Errorf("SSH command failed on %s with exit code %d: %s",
            machine.Hostname, commandResult.ExitCode, commandResult.Stderr)
    }
    
    return commandResult.Stdout, nil
}
```

## Fact Collection Commands

### System Information Commands

The collector uses standard SSH commands to gather system information:

```go
// collectOSFacts collects operating system facts via SSH
func (c *SystemFactCollector) collectOSFacts(machine *spookytypes.Machine) (*spookytypesfacts.OSFacts, error) {
    facts := &spookytypesfacts.OSFacts{}
    
    // Get OS name and version
    osRelease, err := c.runSSHCommand(machine, "cat /etc/os-release")
    if err == nil {
        osInfo := c.parseOSRelease(osRelease)
        facts.Name = osInfo["NAME"]
        facts.Version = osInfo["VERSION"]
    }
    
    // Get kernel version
    kernel, err := c.runSSHCommand(machine, "uname -r")
    if err == nil {
        facts.Kernel = strings.TrimSpace(kernel)
    }
    
    // Get hostname
    hostname, err := c.runSSHCommand(machine, "hostname")
    if err == nil {
        facts.Hostname = strings.TrimSpace(hostname)
    }
    
    // Get uptime
    uptime, err := c.runSSHCommand(machine, "cat /proc/uptime")
    if err == nil {
        if uptimeSeconds, err := strconv.ParseFloat(strings.Fields(uptime)[0], 64); err == nil {
            facts.Uptime = int64(uptimeSeconds)
        }
    }
    
    return facts, nil
}
```

### Hardware Information Commands

```go
// collectHardwareFacts collects hardware facts via SSH
func (c *SystemFactCollector) collectHardwareFacts(machine *spookytypes.Machine) (*spookytypesfacts.HardwareFacts, error) {
    facts := &spookytypesfacts.HardwareFacts{}
    
    // Get CPU count
    cpuCount, err := c.runSSHCommand(machine, "nproc")
    if err == nil {
        if count, err := strconv.Atoi(strings.TrimSpace(cpuCount)); err == nil {
            facts.CPUCount = count
        }
    }
    
    // Get memory information
    memInfo, err := c.runSSHCommand(machine, "cat /proc/meminfo")
    if err == nil {
        facts.MemoryTotal = c.parseMemoryTotal(memInfo)
    }
    
    // Get disk information
    diskInfo, err := c.runSSHCommand(machine, "df -B1 /")
    if err == nil {
        facts.DiskTotal = c.parseDiskTotal(diskInfo)
    }
    
    return facts, nil
}
```

## Error Handling

### Fact Collection Errors

```go
// FactCollectionError represents fact collection errors
type FactCollectionError struct {
    Machine     string
    ErrorType   string
    Message     string
    Recoverable bool
}

func (e *FactCollectionError) Error() string {
    return fmt.Sprintf("fact collection error on %s (%s): %s", e.Machine, e.ErrorType, e.Message)
}
```

### Common Error Scenarios

1. **SSH Connection Failure**: Cannot establish SSH connection to remote machine
2. **Command Execution Failure**: SSH command fails on remote machine
3. **Permission Denied**: Command requires elevated privileges
4. **Machine ID Issues**: Cannot read or validate machine ID
5. **Parse Errors**: Cannot parse command output

## Integration Issues

### Known SSH Integration Problems

The facts system has known issues with SSH-based fact collection:

1. **Connection Problems**: SSH connections may fail or timeout
2. **Authentication Issues**: SSH authentication may not work properly
3. **Command Execution**: Remote commands may fail or return unexpected results
4. **Machine Configuration**: Machine inventory configuration may not be properly used

### Workarounds

Until SSH integration is fixed:

1. **Test SSH Connectivity**: Use `spooky machines ping` to test SSH connectivity first
2. **Verify Machine Configuration**: Ensure machine inventory is properly configured
3. **Check SSH Keys**: Verify SSH keys are properly configured and accessible
4. **Use Verbose Output**: Enable verbose output to see detailed error information

## CLI Integration

### Facts Export Command

```go
// factsExportCmd represents the facts export command
var factsExportCmd = &cobra.Command{
    Use:   "export [project-path]",
    Short: "Export machine facts",
    Long: `Export machine facts to files in various formats.
    
Facts are system information collected from machines including OS details,
hardware information, network configuration, and custom data.`,
    Args: cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        projectPath := args[0]
        
        // Get export format
        format, _ := cmd.Flags().GetString("format")
        if format == "" {
            format = "json"
        }
        
        // Get output path
        outputPath, _ := cmd.Flags().GetString("output")
        if outputPath == "" {
            outputPath = fmt.Sprintf("facts.%s", format)
        }
        
        // Export facts
        return factsManager.ExportFacts(cmd.Context(), projectPath, format, outputPath)
    },
}
```

### Command Options

```go
// Add command flags
factsExportCmd.Flags().String("format", "json", "Export format (json, hcl)")
factsExportCmd.Flags().String("output", "", "Output file path")
factsExportCmd.Flags().StringSlice("machines", nil, "Target specific machines")
factsExportCmd.Flags().StringSlice("tags", nil, "Target machines by tags")
factsExportCmd.Flags().String("filter", "", "Complex machine filter query")
factsExportCmd.Flags().Int("parallel", 1, "Number of parallel connections")
```

## Configuration

### Machine Inventory Integration

Facts collection uses machine inventory for target identification:

```hcl
# machines.hcl
machines {
  machine "web-server" {
    hostname = "web.example.com"
    host     = "192.168.1.100"
    port     = 22
    user     = "admin"
    key_file = "~/.ssh/id_ed25519"
    
    tags = ["web", "production"]
  }
}
```

### Project Configuration

Facts export uses project configuration:

```hcl
# project.hcl
project {
  name = "my-project"
  description = "Example project for facts collection"
  
  metadata {
    version = "1.0.0"
    author = "admin"
  }
}
```

## Future Enhancements

### Planned Improvements

1. **Fix SSH Integration**: Resolve SSH-based fact collection issues
2. **Parallel Processing**: Implement parallel fact collection across multiple machines
3. **Custom Facts**: Support for custom fact collection from remote files
4. **Fact Validation**: Enhanced validation and schema checking
5. **Performance Optimization**: Improve collection performance and reliability

### Advanced Features

1. **Fact History**: Historical fact tracking and comparison
2. **Real-time Collection**: Real-time fact monitoring
3. **Fact Comparison**: Compare facts across machines and time periods
4. **Template Integration**: Use facts in template rendering
5. **Variable Integration**: Use facts in variable resolution

## Conclusion

The facts system provides basic fact collection and export capabilities but has known issues with SSH-based collection that need to be addressed. The system is functional for basic use cases but requires improvements for production use with remote machines.

### Current Limitations

- SSH-based fact collection has implementation issues
- No parallel processing for multiple machines
- Limited error handling and recovery
- No custom fact collection from remote files
- Basic validation and schema checking

### Recommendations

1. **Test Thoroughly**: Test SSH connectivity before running facts collection
2. **Monitor Errors**: Pay attention to error messages and SSH connection issues
3. **Use Verbose Output**: Enable verbose output for debugging
4. **Plan for Improvements**: Consider the planned enhancements for future use
