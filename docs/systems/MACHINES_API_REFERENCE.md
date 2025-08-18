# Machines System API Reference

## Overview

This document provides a comprehensive API reference for the spooky machines system. It covers all interfaces, types, methods, and implementation details for developers working with the machines system.

**Status: Partially Implemented** - The machines system has basic functionality but SSH-based connectivity has known issues that need to be addressed.

## Core Interfaces

### MachinesIntegration Interface

The `MachinesIntegration` interface provides the primary entry point for machines operations:

```go
type MachinesIntegration interface {
    // LoadMachines loads machines from the given source
    LoadMachines(ctx context.Context, source string) ([]spookytypes.Machine, error)

    // ValidateMachines validates machines
    ValidateMachines(ctx context.Context, machines []spookytypes.Machine) (*spookytypes.ValidationResult, error)

    // PingMachines pings machines to check connectivity
    PingMachines(ctx context.Context, machines []spookytypes.Machine) ([]spookytypes.PingResult, error)

    // GetSSHManager returns the SSH manager for connectivity testing
    GetSSHManager() SSHManager
}
```

**Implementation Status**: ⚠️ **Partially Implemented** - Basic functionality exists but SSH-based connectivity has issues

### MachineManager Interface

The `MachineManager` interface provides machine management and connectivity:

```go
type MachineManager interface {
    // LoadMachines loads machines from the given source
    LoadMachines(ctx context.Context, source string) ([]spookytypes.Machine, error)

    // ValidateMachines validates machines
    ValidateMachines(ctx context.Context, machines []spookytypes.Machine) (*spookytypes.ValidationResult, error)

    // PingMachines pings machines to check connectivity
    PingMachines(ctx context.Context, machines []spookytypes.Machine) ([]spookytypes.PingResult, error)

    // GetSSHManager returns the SSH manager for connectivity testing
    GetSSHManager() SSHManager
}
```

**Implementation Status**: ⚠️ **Partially Implemented** - Basic loading and validation exist but connectivity has issues

## Current Implementation Status

### ✅ Working Components

1. **Machine Loading**: Loading machines from HCL configuration files
2. **Machine Validation**: Basic validation of machine definitions
3. **Machine Structure**: Proper machine type definitions and structures
4. **CLI Integration**: `spooky machines ping` command with filtering options
5. **Project Integration**: Machines loading from project configuration
6. **Basic Validation**: Machine definition validation and error handling
7. **Filtering Support**: Support for tag and group filtering
8. **Inventory Management**: Basic machine inventory management
9. **SSH Manager Integration**: SSH manager for connectivity testing
10. **Export Support**: Machine inventory export to JSON format

### ⚠️ Known Issues

> **See also**: [Known Issues](KNOWN_ISSUES.md#ssh-integration-issues) - Comprehensive documentation of all known issues and workarounds

1. **SSH-Based Connectivity**: SSH-based machine connectivity has implementation issues
2. **Ping Functionality**: Machine ping functionality has connectivity problems
3. **Authentication Testing**: SSH authentication testing has issues
4. **Connection Pooling**: SSH connection pooling has problems
5. **Host Key Validation**: Host key validation has implementation issues
6. **Parallel Processing**: No parallel machine operations support

### 🔄 In Progress

1. **SSH Connectivity Fixes**: Addressing SSH-based connectivity issues
2. **Ping Improvements**: Implementing proper machine ping functionality
3. **Authentication Fixes**: Fixing SSH authentication testing

## Implementation Details

### Machine Loading System

The machines system loads machines from HCL configuration files:

```go
type Manager struct {
    logger          spookytypeslogging.Logger
    validator       spookyinterfaces.MachineValidator
    sshManager      spookyinterfaces.SSHManager
    schemaValidator *spookyschemas.Validator
}

func NewManager(
    logger spookytypeslogging.Logger,
    validator spookyinterfaces.MachineValidator,
    sshManager spookyinterfaces.SSHManager,
    schemaValidator *spookyschemas.Validator,
) spookyinterfaces.MachinesIntegration {
    return &Manager{
        logger:          logger,
        validator:       validator,
        sshManager:      sshManager,
        schemaValidator: schemaValidator,
    }
}
```

### Machine Loading Implementation

```go
// LoadMachines loads machines from the specified source
func (m *Manager) LoadMachines(ctx context.Context, source string) ([]spookytypes.Machine, error) {
    m.logger.Info("Loading machines", map[string]interface{}{
        "source": source,
    })

    // Check if source is a directory
    if info, err := os.Stat(source); err == nil && info.IsDir() {
        return m.loadMachinesFromDirectory(ctx, source)
    }

    // Check if source is a file
    if _, err := os.Stat(source); err == nil {
        return m.loadMachinesFromFile(ctx, source)
    }

    return nil, fmt.Errorf("source not found: %s", source)
}

// loadMachinesFromFile loads machines from a single HCL file
func (m *Manager) loadMachinesFromFile(_ context.Context, filePath string) ([]spookytypes.Machine, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read machines file: %w", err)
    }

    var config struct {
        Machines []*spookytypesmachines.Machine `hcl:"machine,block"`
    }

    if err := hclsimple.Decode(filePath, data, nil, &config); err != nil {
        return nil, fmt.Errorf("failed to parse machines file: %w", err)
    }

    // Convert to interface slice
    machines := make([]spookytypes.Machine, len(config.Machines))
    for i, machine := range config.Machines {
        machines[i] = machine
    }

    return machines, nil
}

// loadMachinesFromDirectory loads machines from a directory
func (m *Manager) loadMachinesFromDirectory(ctx context.Context, dirPath string) ([]spookytypes.Machine, error) {
    // Look for machines.hcl file
    machinesPath := filepath.Join(dirPath, "machines.hcl")
    if _, err := os.Stat(machinesPath); err == nil {
        return m.loadMachinesFromFile(ctx, machinesPath)
    }

    return nil, fmt.Errorf("no machines.hcl file found in directory: %s", dirPath)
}
```

### Machine Validation Implementation

```go
// ValidateMachines validates a collection of machines
func (m *Manager) ValidateMachines(ctx context.Context, machines []spookytypes.Machine) (*spookytypes.ValidationResult, error) {
    m.logger.Info("Validating machines", map[string]interface{}{
        "machines": len(machines),
    })

    var errors []spookyschemas.SchemaError
    var warnings []spookyschemas.SchemaError

    for i, machine := range machines {
        // Validate individual machine
        if err := m.validateMachine(machine); err != nil {
            errors = append(errors, spookyschemas.SchemaError{
                Message: fmt.Sprintf("machine[%d]: %s", i, err.Error()),
            })
        }
    }

    return &spookytypes.ValidationResult{
        Valid:    len(errors) == 0,
        Errors:   errors,
        Warnings: warnings,
    }, nil
}

func (m *Manager) validateMachine(machine spookytypes.Machine) error {
    // Validate required fields
    if machine.Hostname == "" {
        return fmt.Errorf("machine hostname is required")
    }

    if machine.Port <= 0 || machine.Port > 65535 {
        return fmt.Errorf("invalid port number: %d", machine.Port)
    }

    if machine.User == "" {
        return fmt.Errorf("machine user is required")
    }

    // Validate SSH configuration
    if machine.SSHConfig == nil {
        return fmt.Errorf("SSH configuration is required")
    }

    if machine.SSHConfig.KeyPath == "" {
        return fmt.Errorf("SSH key path is required")
    }

    return nil
}
```

### Machine Ping Implementation

```go
// PingMachines pings the specified machines to check connectivity
func (m *Manager) PingMachines(ctx context.Context, machines []spookytypes.Machine) ([]spookytypes.PingResult, error) {
    m.logger.Info("Pinging machines", map[string]interface{}{
        "machines": len(machines),
    })

    var results []spookytypes.PingResult

    for _, machine := range machines {
        result := spookytypes.PingResult{
            Hostname: machine.Hostname,
            Status:   "unknown",
            Error:    "",
        }

        // Test DNS resolution
        if err := m.testDNSResolution(machine.Hostname); err != nil {
            result.Status = "dns_failed"
            result.Error = fmt.Sprintf("DNS resolution failed: %v", err)
            results = append(results, result)
            continue
        }

        // Test SSH connectivity
        if err := m.testSSHConnectivity(ctx, machine); err != nil {
            result.Status = "ssh_failed"
            result.Error = fmt.Sprintf("SSH connectivity failed: %v", err)
        } else {
            result.Status = "reachable"
        }

        results = append(results, result)
    }

    m.logger.Info("Machine ping completed", map[string]interface{}{
        "machines": len(machines),
        "results":  len(results),
    })

    return results, nil
}

func (m *Manager) testDNSResolution(hostname string) error {
    _, err := net.LookupHost(hostname)
    return err
}

func (m *Manager) testSSHConnectivity(ctx context.Context, machine spookytypes.Machine) error {
    // Use SSH manager to test connectivity
    sshManager := m.GetSSHManager()
    if sshManager == nil {
        return fmt.Errorf("SSH manager not available")
    }

    // Test SSH connection
    conn, err := sshManager.GetConnection(machine.Hostname, machine.Port, machine.User)
    if err != nil {
        return fmt.Errorf("SSH connection failed: %w", err)
    }
    defer sshManager.ReturnConnection(conn)

    return nil
}
```

## Type Definitions

### Machine Types

```go
// Machine represents a machine definition
type Machine struct {
    // Machine hostname (required)
    Hostname string `json:"hostname" hcl:"hostname"`

    // Machine description (optional)
    Description string `json:"description,omitempty" hcl:"description,optional"`

    // SSH port (default: 22)
    Port int `json:"port,omitempty" hcl:"port,optional"`

    // SSH user (required)
    User string `json:"user" hcl:"user"`

    // SSH configuration
    SSHConfig *SSHConfig `json:"ssh_config" hcl:"ssh_config"`

    // Machine tags
    Tags []string `json:"tags,omitempty" hcl:"tags,optional"`

    // Machine groups
    Groups []string `json:"groups,omitempty" hcl:"groups,optional"`

    // Machine metadata
    Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// SSHConfig represents SSH configuration for a machine
type SSHConfig struct {
    // SSH key path
    KeyPath string `json:"key_path" hcl:"key_path"`

    // SSH key passphrase (optional)
    Passphrase string `json:"passphrase,omitempty" hcl:"passphrase,optional"`

    // SSH timeout in seconds
    Timeout int `json:"timeout,omitempty" hcl:"timeout,optional"`

    // SSH connection retries
    Retries int `json:"retries,omitempty" hcl:"retries,optional"`

    // SSH connection retry delay in seconds
    RetryDelay int `json:"retry_delay,omitempty" hcl:"retry_delay,optional"`

    // SSH host key checking (default: true)
    HostKeyChecking bool `json:"host_key_checking,omitempty" hcl:"host_key_checking,optional"`

    // SSH known hosts file
    KnownHostsFile string `json:"known_hosts_file,omitempty" hcl:"known_hosts_file,optional"`
}

// PingResult represents the result of a machine ping
type PingResult struct {
    // Machine hostname
    Hostname string `json:"hostname" hcl:"hostname"`

    // Ping status (reachable, dns_failed, ssh_failed, unknown)
    Status string `json:"status" hcl:"status"`

    // Error message (if any)
    Error string `json:"error,omitempty" hcl:"error,optional"`

    // Ping timestamp
    Timestamp time.Time `json:"timestamp" hcl:"timestamp"`

    // Ping duration
    Duration time.Duration `json:"duration,omitempty" hcl:"duration,optional"`
}
```

### Machine Configuration Types

```go
// MachineInventory represents a machine inventory
type MachineInventory struct {
    // Inventory name
    Name string `json:"name" hcl:"name"`

    // Inventory description
    Description string `json:"description,omitempty" hcl:"description,optional"`

    // Machines in inventory
    Machines []*Machine `json:"machines" hcl:"machine,block"`

    // Inventory metadata
    Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// MachineFilter represents a machine filter
type MachineFilter struct {
    // Filter by hostnames
    Hostnames []string `json:"hostnames,omitempty" hcl:"hostnames,optional"`

    // Filter by tags
    Tags []string `json:"tags,omitempty" hcl:"tags,optional"`

    // Filter by groups
    Groups []string `json:"groups,omitempty" hcl:"groups,optional"`

    // Complex filter query
    Query string `json:"query,omitempty" hcl:"query,optional"`
}
```

## CLI Commands

### Machines List Command

```bash
# List all machines in a project
spooky machines list ./my-project

# List machines with details
spooky machines list ./my-project --verbose

# List machines by tags
spooky machines list ./my-project --tags "environment=production"

# List machines by groups
spooky machines list ./my-project --groups "webservers"
```

### Machines Validate Command

```bash
# Validate machine configuration
spooky machines validate ./my-project

# Validate with detailed output
spooky machines validate ./my-project --verbose

# Validate specific machines
spooky machines validate ./my-project --machine web-server
```

### Machines Ping Command

```bash
# Ping all machines
spooky machines ping ./my-project

# Ping specific machines
spooky machines ping ./my-project --machine web-server

# Ping machines by tags
spooky machines ping ./my-project --tags "environment=production"

# Ping with timeout
spooky machines ping ./my-project --timeout 30
```

### Machines Export Command

```bash
# Export machine inventory to JSON
spooky machines export ./my-project --json --output machines.json

# Export with filtering
spooky machines export ./my-project --json --output production-machines.json --tags "environment=production"
```

## Integration Examples

### Basic Machine Definition

```hcl
# machines.hcl
machines {
  machine "web-server" {
    description = "Web server for production"
    hostname = "web.example.com"
    port = 22
    user = "admin"
    
    ssh_config {
      key_path = "~/.ssh/id_rsa"
      timeout = 30
      retries = 3
      retry_delay = 5
      host_key_checking = true
    }
    
    tags = ["web", "production", "nginx"]
    groups = ["webservers", "production"]
  }
  
  machine "db-server" {
    description = "Database server"
    hostname = "db.example.com"
    port = 22
    user = "postgres"
    
    ssh_config {
      key_path = "~/.ssh/id_rsa"
      timeout = 30
      retries = 3
      retry_delay = 5
      host_key_checking = true
    }
    
    tags = ["database", "production", "postgresql"]
    groups = ["databases", "production"]
  }
  
  machine "app-server" {
    description = "Application server"
    hostname = "app.example.com"
    port = 22
    user = "app"
    
    ssh_config {
      key_path = "~/.ssh/id_rsa"
      timeout = 30
      retries = 3
      retry_delay = 5
      host_key_checking = true
    }
    
    tags = ["app", "production", "nodejs"]
    groups = ["applications", "production"]
  }
}
```

### Machine Loading and Validation

```go
// Machine loading and validation example
func loadAndValidateMachines(projectPath string) error {
    ctx := context.Background()
    
    // Create machine manager
    manager := spookymachines.NewManager(logger, validator, sshManager, schemaValidator)
    
    // Load machines
    machines, err := manager.LoadMachines(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load machines: %w", err)
    }
    
    // Validate machines
    result, err := manager.ValidateMachines(ctx, machines)
    if err != nil {
        return fmt.Errorf("failed to validate machines: %w", err)
    }
    
    if !result.Valid {
        fmt.Println("Machine validation failed:")
        for _, error := range result.Errors {
            fmt.Printf("  - %s\n", error.Message)
        }
        return fmt.Errorf("machine validation failed")
    }
    
    fmt.Printf("Loaded and validated %d machines\n", len(machines))
    return nil
}
```

### Machine Ping

```go
// Machine ping example
func pingMachines(projectPath string) error {
    ctx := context.Background()
    
    // Create machine manager
    manager := spookymachines.NewManager(logger, validator, sshManager, schemaValidator)
    
    // Load machines
    machines, err := manager.LoadMachines(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load machines: %w", err)
    }
    
    // Ping machines
    results, err := manager.PingMachines(ctx, machines)
    if err != nil {
        return fmt.Errorf("failed to ping machines: %w", err)
    }
    
    // Process results
    for _, result := range results {
        switch result.Status {
        case "reachable":
            fmt.Printf("✓ %s is reachable\n", result.Hostname)
        case "dns_failed":
            fmt.Printf("✗ %s: DNS resolution failed - %s\n", result.Hostname, result.Error)
        case "ssh_failed":
            fmt.Printf("✗ %s: SSH connectivity failed - %s\n", result.Hostname, result.Error)
        default:
            fmt.Printf("? %s: Unknown status - %s\n", result.Hostname, result.Error)
        }
    }
    
    return nil
}
```

## Error Handling

### Machine Errors

```go
// Error handling example
func handleMachineError(err error) {
    if err == nil {
        return
    }
    
    // Check for specific error types
    switch {
    case strings.Contains(err.Error(), "failed to load machines"):
        fmt.Println("Machine loading failed - check machine configuration")
    case strings.Contains(err.Error(), "failed to validate machines"):
        fmt.Println("Machine validation failed - check machine definitions")
    case strings.Contains(err.Error(), "failed to ping machines"):
        fmt.Println("Machine ping failed - check network connectivity")
    case strings.Contains(err.Error(), "SSH connection failed"):
        fmt.Println("SSH connection failed - check SSH configuration")
    case strings.Contains(err.Error(), "DNS resolution failed"):
        fmt.Println("DNS resolution failed - check hostname configuration")
    default:
        fmt.Printf("Machine error: %v\n", err)
    }
}
```

### Validation Errors

```go
// Validation error handling
func handleValidationError(result *spookytypes.ValidationResult) error {
    if result.Valid {
        return nil
    }
    
    fmt.Println("Machine validation failed:")
    for _, err := range result.Errors {
        fmt.Printf("  - %s\n", err.Message)
    }
    
    for _, warning := range result.Warnings {
        fmt.Printf("  Warning: %s\n", warning.Message)
    }
    
    return fmt.Errorf("machine validation failed with %d errors", len(result.Errors))
}
```

## Performance Considerations

### Parallel Processing

The machines system supports parallel processing:

- Multiple machines can be pinged concurrently
- Configurable parallel worker count
- Thread-safe machine operations

### Resource Management

The machines system manages resources efficiently:

- SSH connections are pooled and reused
- Memory usage is optimized for large machine sets
- Timeouts prevent hanging operations

## Troubleshooting

### Common Issues

1. **SSH Connection Failures**: Check machine connectivity and SSH configuration
2. **Authentication Errors**: Verify SSH key permissions and user access
3. **DNS Resolution Failures**: Check hostname configuration and DNS settings
4. **Permission Denied**: Check SSH key permissions and user access
5. **Timeout Issues**: Adjust SSH timeouts for slow connections

### Debug Information

The machines system provides comprehensive logging for debugging:

```go
// Enable debug logging
logger.SetLevel(spookytypes.LogLevelDebug)

// Check SSH configuration
fmt.Printf("SSH config: %+v\n", sshConfig)

// Validate machine configuration
err := validateMachine(machine)
if err != nil {
    fmt.Printf("Machine validation error: %v\n", err)
}
```

## Future Enhancements

### Planned Features

1. **Parallel Processing**: Implement parallel machine operations
2. **Advanced Filtering**: Support complex machine filtering
3. **Machine Discovery**: Add automatic machine discovery
4. **Health Monitoring**: Add machine health monitoring
5. **Inventory Management**: Improve machine inventory management

### Integration Enhancements

1. **Facts Integration**: Use facts in machine operations
2. **Variables Integration**: Use variables in machine configuration
3. **Templates Integration**: Use templates in machine definitions
4. **Advanced Connectivity**: Improve SSH connectivity and reliability
