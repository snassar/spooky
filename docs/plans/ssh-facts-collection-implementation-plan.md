# SSH-Based Fact Collection Implementation Plan

## Overview

This document outlines the implementation plan for adding SSH-based fact collection functionality to spooky. Currently, the facts system only supports local fact collection, which severely limits its usefulness in multi-machine environments. This plan addresses the missing SSH-based fact collection and remote `/etc/spooky/facts.*` reading capabilities.

## Current State Analysis

### Existing Implementation
- **Location**: `internal/facts/`
- **Current Capabilities**: 
  - ✅ Local fact collection using `SystemFactCollector`
  - ✅ Memory-based fact storage and export
  - ✅ JSON and HCL export formats
  - ✅ Basic fact validation and management
- **Current Limitations**:
  - ❌ No SSH-based fact collection
  - ❌ No remote `/etc/spooky/facts.*` reading
  - ❌ No parallel collection across multiple machines
  - ❌ No integration with existing SSH infrastructure
  - ❌ No machine inventory authentication support

### Missing Functionality
1. **SSH Fact Collection**: Cannot collect facts from remote machines via SSH
2. **Remote Facts Reading**: Cannot read `/etc/spooky/facts.*` files from remote machines
3. **Parallel Processing**: Sequential collection only, no multi-machine parallel processing
4. **SSH Integration**: Cannot leverage existing SSH infrastructure and machine inventory
5. **Authentication**: No support for machine inventory authentication methods

## Implementation Objectives

1. **Implement SSH-based fact collection** for remote machines
2. **Support reading `/etc/spooky/facts.*` files** from remote machines
3. **Integrate with existing SSH infrastructure** and machine inventory
4. **Enable parallel fact collection** across multiple machines
5. **Maintain backward compatibility** with existing local collection
6. **Provide comprehensive error handling** and retry logic
7. **Enable fact collection filtering** and targeting
8. **Support custom fact collection scripts** on remote machines
9. **Provide detailed collection metrics** and progress reporting
10. **Support encrypted facts** using age encryption (future integration)
11. **Enable advanced fact collection** via separate advanced collector (see advanced-facts-collection-plan.md)

## Technical Architecture

### Core Components

#### 1. SSH Fact Collector
- **Location**: `internal/facts/collectors/ssh/`
- **Purpose**: Collect facts from remote machines via SSH
- **Integration**: Uses existing `SSHManager` interface
- **Authentication**: Leverages machine inventory authentication methods
- **Commands**: Executes system commands to gather facts

#### 2. Local Fact Collector (Enhanced)
- **Location**: `internal/facts/collectors/local/`
- **Purpose**: Collect facts from local machine using system commands
- **Replacement**: Replace gopsutil dependency with direct system commands
- **Consistency**: Same fact collection methods as SSH collector

#### 3. Facts Manager (Enhanced)
- **Location**: `internal/facts/manager.go`
- **Purpose**: Coordinate fact collection across multiple machines
- **Parallel Processing**: Support concurrent collection from multiple machines
- **Integration**: Integrate with machine inventory and SSH infrastructure

#### 4. CLI Integration
- **Location**: `cmd/facts.go`
- **Purpose**: Provide CLI interface for fact collection
- **Commands**: `spooky facts gather`, `spooky facts export`
- **Options**: Parallel processing, machine targeting, filtering

### Fact Collection Strategy

#### No-Sudo Policy
**Important**: SSH fact collection will NOT use sudo. This ensures:
- **Security**: No need to grant sudo access to spooky users
- **Simplicity**: No complex sudo configuration required
- **Admin Control**: Administrators explicitly control privileged information
- **Reliability**: No sudo timeout or permission issues

#### Fact Collection Sources
1. **System Commands**: Basic system information via SSH commands
2. **Custom Facts Files**: `/etc/spooky/facts.*` files on remote machines
3. **Advanced Collector Output**: Facts from advanced collector (separate plan)

### Admin Workflow for Privileged Facts
1. **Create custom facts files**: Place privileged information in `/etc/spooky/facts.hcl` or `/etc/spooky/facts.json`
2. **Set appropriate permissions**: Ensure spooky user can read the facts files
3. **Use custom collectors**: Implement collectors that run with appropriate permissions
4. **Configure system access**: Grant necessary file read permissions where possible
5. **Deploy advanced collectors**: Use spooky advanced collectors for comprehensive fact gathering (see advanced-facts-collection-plan.md)

## Implementation Details

### 1. SSH Fact Collector Implementation

#### Directory Structure
```
internal/facts/collectors/
├── local/
│   ├── collector.go      # Enhanced local collector (no gopsutil)
│   └── commands.go       # System command definitions
├── ssh/
│   ├── collector.go      # SSH-based fact collector
│   ├── commands.go       # SSH command execution
│   └── authentication.go # SSH authentication handling
└── common/
    ├── parser.go         # Shared fact parsing logic
    └── validator.go      # Shared fact validation
```

#### SSH Collector Interface
```go
// internal/facts/collectors/ssh/collector.go
type SSHFactCollector struct {
    sshManager spookyinterfaces.SSHManager
    logger     spookylogging.Logger
    config     *spookytypes.Config
}

func (c *SSHFactCollector) CollectFacts(ctx context.Context, machine *spookytypes.Machine) (*spookytypes.FactCollection, error) {
    // 1. Collect basic facts via SSH commands
    basicFacts, err := c.collectBasicFacts(ctx, machine)
    if err != nil {
        return nil, err
    }
    
    // 2. Check for custom facts files
    customFacts, err := c.collectCustomFacts(ctx, machine)
    if err == nil {
        return c.mergeFacts(basicFacts, customFacts), nil
    }
    
    return basicFacts, nil
}
```

#### SSH Command Execution
```go
// internal/facts/collectors/ssh/commands.go
func (c *SSHFactCollector) executeCommand(ctx context.Context, machine *spookytypes.Machine, command string) (string, error) {
    // Use existing SSH infrastructure
    session, err := c.sshManager.CreateSession(ctx, machine)
    if err != nil {
        return "", fmt.Errorf("failed to create SSH session: %w", err)
    }
    defer session.Close()
    
    // Execute command with timeout
    output, err := session.ExecuteCommand(command, 30*time.Second)
    if err != nil {
        return "", fmt.Errorf("failed to execute command: %w", err)
    }
    
    return output, nil
}
```

### 2. Enhanced Local Fact Collector

#### Replace gopsutil with System Commands
```go
// internal/facts/collectors/local/collector.go
type LocalFactCollector struct {
    logger spookylogging.Logger
}

func (c *LocalFactCollector) CollectFacts(ctx context.Context) (*spookytypes.FactCollection, error) {
    facts := &spookytypes.FactCollection{
        CollectedAt: time.Now(),
        Facts:       &spookytypesfacts.Facts{},
    }
    
    // Use same commands as SSH collector for consistency
    if err := c.collectSystemFacts(ctx, facts); err != nil {
        return nil, fmt.Errorf("failed to collect system facts: %w", err)
    }
    
    if err := c.collectNetworkFacts(ctx, facts); err != nil {
        return nil, fmt.Errorf("failed to collect network facts: %w", err)
    }
    
    if err := c.collectHardwareFacts(ctx, facts); err != nil {
        return nil, fmt.Errorf("failed to collect hardware facts: %w", err)
    }
    
    return facts, nil
}
```

#### System Command Definitions
```go
// internal/facts/collectors/common/commands.go
var SystemCommands = map[string]string{
    "os_release":     "cat /etc/os-release",
    "kernel_info":    "uname -a",
    "cpu_info":       "cat /proc/cpuinfo",
    "memory_info":    "free -b",
    "disk_usage":     "df -h",
    "network_interfaces": "ip addr show",
    "hostname":       "hostname",
    "uptime":         "uptime",
    "load_average":   "cat /proc/loadavg",
    "process_count":  "ps aux | wc -l",
}
```

### 3. Enhanced Facts Manager

#### Parallel Collection Support
```go
// internal/facts/manager.go
func (m *FactManager) CollectFactsFromMachines(ctx context.Context, machines []*spookytypes.Machine, parallel int) (map[string]*spookytypes.FactCollection, error) {
    results := make(map[string]*spookytypes.FactCollection)
    var mu sync.Mutex
    var wg sync.WaitGroup
    
    // Create worker pool
    semaphore := make(chan struct{}, parallel)
    
    for _, machine := range machines {
        wg.Add(1)
        go func(m *spookytypes.Machine) {
            defer wg.Done()
            semaphore <- struct{}{} // Acquire semaphore
            defer func() { <-semaphore }() // Release semaphore
            
            facts, err := m.collectFactsFromMachine(ctx, m)
            if err != nil {
                m.logger.Error("Failed to collect facts from machine", map[string]interface{}{
                    "machine": m.Hostname,
                    "error":   err.Error(),
                })
                return
            }
            
            mu.Lock()
            results[m.Hostname] = facts
            mu.Unlock()
        }(machine)
    }
    
    wg.Wait()
    return results, nil
}
```

### 4. CLI Integration

#### Enhanced Facts Commands
```go
// cmd/facts.go
func newFactsGatherCommand() *cobra.Command {
    var parallel int
    var machines []string
    var tags []string
    var filter string
    
    cmd := &cobra.Command{
        Use:   "gather [project]",
        Short: "Gather facts from machines",
        Long:  "Gather facts from local and remote machines in the project",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            projectPath := args[0]
            
            // Load project and machine inventory
            project, err := loadProject(projectPath)
            if err != nil {
                return err
            }
            
            // Filter machines based on flags
            targetMachines, err := filterMachines(project.Machines, machines, tags, filter)
            if err != nil {
                return err
            }
            
            // Collect facts in parallel
            facts, err := factManager.CollectFactsFromMachines(cmd.Context(), targetMachines, parallel)
            if err != nil {
                return err
            }
            
            // Store facts
            for hostname, factCollection := range facts {
                if err := factManager.StoreFacts(hostname, factCollection); err != nil {
                    return fmt.Errorf("failed to store facts for %s: %w", hostname, err)
                }
            }
            
            return nil
        },
    }
    
    cmd.Flags().IntVar(&parallel, "parallel", 5, "Number of parallel connections")
    cmd.Flags().StringSliceVar(&machines, "machine", nil, "Target specific machines")
    cmd.Flags().StringSliceVar(&tags, "tags", nil, "Target machines by tags")
    cmd.Flags().StringVar(&filter, "filter", "", "Complex machine filter")
    
    return cmd
}
```

## Facts Collection Commands

### System Facts (No Sudo Required)
- **OS Information**: `cat /etc/os-release` → Parse OS name, version, platform
- **Kernel Information**: `uname -a` → Kernel version, architecture
- **CPU Information**: `cat /proc/cpuinfo` → CPU cores, model, frequency
- **Memory Information**: `free -b` → Total, available, used memory
- **Disk Usage**: `df -h` → Filesystem usage and capacity
- **Network Interfaces**: `ip addr show` → Interface names, IP addresses, MAC addresses
- **Hostname**: `hostname` → System hostname
- **Uptime**: `uptime` → System uptime and load average
- **Load Average**: `cat /proc/loadavg` → 1, 5, 15 minute load averages
- **Process Count**: `ps aux | wc -l` → Number of running processes

### Custom Facts Files
- **HCL Files**: `/etc/spooky/facts.hcl` → Structured fact definitions
- **JSON Files**: `/etc/spooky/facts.json` → JSON fact data
- **Multiple Files**: Support for multiple fact files in `/etc/spooky/`

## Implementation Approach

The implementation should focus on:

1. **Core SSH Infrastructure**: Building the SSH fact collector and integration with existing SSH infrastructure
2. **Remote Facts Reading**: Implementing the ability to read `/etc/spooky/facts.*` files from remote machines
3. **Parallel Collection**: Enabling concurrent fact collection across multiple machines
4. **Integration**: Seamlessly integrating with the existing facts system and CLI
5. **Testing**: Comprehensive testing with both mock and real SSH environments

## Testing Strategy

### Unit Testing
- **Mock SSH Server**: Use mock SSH server for unit tests
- **Command Execution**: Test individual command execution and parsing
- **Error Handling**: Test SSH connection failures and timeouts
- **Fact Parsing**: Test fact parsing from command output

### Integration Testing
- **Real SSH Environment**: Test with actual SSH servers
- **Machine Inventory**: Test with real machine inventory and authentication
- **Parallel Collection**: Test concurrent fact collection
- **Error Recovery**: Test error handling and recovery mechanisms

### Performance Testing
- **Parallel Scaling**: Test performance with increasing parallel connections
- **Large Machine Sets**: Test with hundreds of machines
- **Network Latency**: Test with high-latency connections
- **Resource Usage**: Monitor memory and CPU usage during collection

## Success Criteria

1. **SSH Fact Collection**: Successfully collect facts from remote machines via SSH
2. **Remote Facts Reading**: Read `/etc/spooky/facts.*` files from remote machines
3. **Parallel Processing**: Support concurrent fact collection from multiple machines
4. **Integration**: Seamless integration with existing spooky architecture
5. **Performance**: Efficient fact collection with minimal resource usage
6. **Reliability**: Robust error handling and recovery mechanisms
7. **Backward Compatibility**: Maintain compatibility with existing local collection
8. **Security**: No sudo requirements, secure SSH authentication
9. **Usability**: Clear CLI interface and comprehensive documentation
10. **Testing**: Comprehensive test coverage and validation

## Next Steps

1. **Implementation**: Begin with core SSH fact collection infrastructure
2. **Testing Setup**: Establish testing environment with mock SSH servers
3. **Documentation**: Update user documentation and examples
4. **Integration**: Integrate with existing spooky architecture
5. **Validation**: Test with real SSH environments and user feedback
