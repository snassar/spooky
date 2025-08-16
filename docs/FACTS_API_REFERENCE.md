# Facts System API Reference

## Overview

This document provides a comprehensive API reference for the spooky facts system. It covers all interfaces, types, methods, and implementation details for developers working with the facts system.

**Status: Partially Implemented** - The facts system has basic functionality but SSH-based fact collection has known issues that need to be addressed.

## Core Interfaces

### FactsIntegration Interface

The `FactsIntegration` interface provides the primary entry point for facts operations:

```go
type FactsIntegration interface {
    // CollectFacts collects facts from the given machine
    CollectFacts(ctx context.Context, machine spookytypes.Machine) (*spookytypes.FactCollection, error)

    // StoreFacts stores facts in the facts database
    StoreFacts(ctx context.Context, facts *spookytypes.FactCollection) error

    // GetFacts retrieves facts from the facts database
    GetFacts(ctx context.Context, machine string) (*spookytypes.FactCollection, error)

    // ListFacts lists all facts in the facts database
    ListFacts(ctx context.Context) ([]*spookytypes.FactCollection, error)

    // ValidateFacts validates facts
    ValidateFacts(ctx context.Context, facts *spookytypes.FactCollection) (*spookytypes.ValidationResult, error)

    // GetSSHManager returns the SSH manager for fact collection
    GetSSHManager() SSHManager
}
```

**Implementation Status**: ⚠️ **Partially Implemented** - Basic functionality exists but SSH-based fact collection has issues

### FactManager Interface

The `FactManager` interface provides fact management and collection:

```go
type FactManager interface {
    // CollectFacts collects facts from the given machine
    CollectFacts(ctx context.Context, machine spookytypes.Machine) (*spookytypes.FactCollection, error)

    // StoreFacts stores facts in the facts database
    StoreFacts(ctx context.Context, facts *spookytypes.FactCollection) error

    // GetFacts retrieves facts from the facts database
    GetFacts(ctx context.Context, machine string) (*spookytypes.FactCollection, error)

    // ListFacts lists all facts in the facts database
    ListFacts(ctx context.Context) ([]*spookytypes.FactCollection, error)

    // ValidateFacts validates facts
    ValidateFacts(ctx context.Context, facts *spookytypes.FactCollection) (*spookytypes.ValidationResult, error)

    // GetSSHManager returns the SSH manager for fact collection
    GetSSHManager() SSHManager
}
```

**Implementation Status**: ⚠️ **Partially Implemented** - Basic loading and validation exist but collection has issues

## Current Implementation Status

### ✅ Working Components

1. **Fact Loading**: Loading facts from HCL configuration files
2. **Fact Validation**: Basic validation of fact definitions
3. **Fact Structure**: Proper fact type definitions and structures
4. **CLI Integration**: `spooky facts gather` command with filtering options
5. **Project Integration**: Facts loading from project configuration
6. **Basic Validation**: Fact definition validation and error handling
7. **Filtering Support**: Support for machine and tag filtering
8. **Export Management**: Basic fact export management
9. **SSH Manager Integration**: SSH manager for fact collection
10. **Export Support**: Facts export to JSON format

### ⚠️ Known Issues

1. **SSH-Based Collection**: SSH-based fact collection has implementation issues
2. **Fact Collection**: Facts cannot be properly collected from remote machines
3. **Authentication Testing**: SSH authentication testing has issues
4. **Connection Pooling**: SSH connection pooling has problems
5. **Host Key Validation**: Host key validation has implementation issues
6. **Parallel Processing**: No parallel fact collection support

### 🔄 In Progress

1. **SSH Collection Fixes**: Addressing SSH-based fact collection issues
2. **Collection Improvements**: Implementing proper fact collection functionality
3. **Authentication Fixes**: Fixing SSH authentication testing

## Implementation Details

### Fact Loading System

The facts system loads facts from HCL configuration files:

```go
type Manager struct {
    logger          spookytypeslogging.Logger
    validator       spookyinterfaces.FactValidator
    sshManager      spookyinterfaces.SSHManager
    schemaValidator *spookyschemas.Validator
    
}

func NewManager(
    logger spookytypeslogging.Logger,
    validator spookyinterfaces.FactValidator,
    sshManager spookyinterfaces.SSHManager,
    schemaValidator *spookyschemas.Validator,
    
) spookyinterfaces.FactsIntegration {
    return &Manager{
        logger:          logger,
        validator:       validator,
        sshManager:      sshManager,
        schemaValidator: schemaValidator,
        
    }
}
```

### Fact Collection Implementation

```go
// CollectFacts collects facts from the specified machine
func (m *Manager) CollectFacts(ctx context.Context, machine spookytypes.Machine) (*spookytypes.FactCollection, error) {
    m.logger.Info("Collecting facts", map[string]interface{}{
        "machine": machine.Hostname,
    })

    // Create fact collection
    collection := &spookytypes.FactCollection{
        Machine:   machine.Hostname,
        Timestamp: time.Now(),
        Facts:     make(map[string]interface{}),
    }

    // Collect system facts
    if err := m.collectSystemFacts(ctx, machine, collection); err != nil {
        m.logger.Error("Failed to collect system facts", err, map[string]interface{}{
            "machine": machine.Hostname,
        })
        return nil, fmt.Errorf("failed to collect system facts: %w", err)
    }

    // Collect network facts
    if err := m.collectNetworkFacts(ctx, machine, collection); err != nil {
        m.logger.Error("Failed to collect network facts", err, map[string]interface{}{
            "machine": machine.Hostname,
        })
        return nil, fmt.Errorf("failed to collect network facts: %w", err)
    }

    // Collect application facts
    if err := m.collectApplicationFacts(ctx, machine, collection); err != nil {
        m.logger.Error("Failed to collect application facts", err, map[string]interface{}{
            "machine": machine.Hostname,
        })
        return nil, fmt.Errorf("failed to collect application facts: %w", err)
    }

    m.logger.Info("Fact collection completed", map[string]interface{}{
        "machine": machine.Hostname,
        "facts":   len(collection.Facts),
    })

    return collection, nil
}

func (m *Manager) collectSystemFacts(ctx context.Context, machine spookytypes.Machine, collection *spookytypes.FactCollection) error {
    // Use SSH manager to collect system facts
    sshManager := m.GetSSHManager()
    if sshManager == nil {
        return fmt.Errorf("SSH manager not available")
    }

    // Get SSH connection
    conn, err := sshManager.GetConnection(machine.Hostname, machine.Port, machine.User)
    if err != nil {
        return fmt.Errorf("failed to get SSH connection: %w", err)
    }
    defer sshManager.ReturnConnection(conn)

    // Collect basic system facts
    facts := map[string]interface{}{
        "os":        m.getOSInfo(conn),
        "hostname":  m.getHostname(conn),
        "uptime":    m.getUptime(conn),
        "memory":    m.getMemoryInfo(conn),
        "disk":      m.getDiskInfo(conn),
        "cpu":       m.getCPUInfo(conn),
        "kernel":    m.getKernelInfo(conn),
    }

    // Add facts to collection
    for key, value := range facts {
        collection.Facts[key] = value
    }

    return nil
}

func (m *Manager) collectNetworkFacts(ctx context.Context, machine spookytypes.Machine, collection *spookytypes.FactCollection) error {
    // Use SSH manager to collect network facts
    sshManager := m.GetSSHManager()
    if sshManager == nil {
        return fmt.Errorf("SSH manager not available")
    }

    // Get SSH connection
    conn, err := sshManager.GetConnection(machine.Hostname, machine.Port, machine.User)
    if err != nil {
        return fmt.Errorf("failed to get SSH connection: %w", err)
    }
    defer sshManager.ReturnConnection(conn)

    // Collect network facts
    facts := map[string]interface{}{
        "interfaces": m.getNetworkInterfaces(conn),
        "routes":     m.getNetworkRoutes(conn),
        "dns":        m.getDNSInfo(conn),
    }

    // Add facts to collection
    for key, value := range facts {
        collection.Facts[key] = value
    }

    return nil
}

func (m *Manager) collectApplicationFacts(ctx context.Context, machine spookytypes.Machine, collection *spookytypes.FactCollection) error {
    // Use SSH manager to collect application facts
    sshManager := m.GetSSHManager()
    if sshManager == nil {
        return fmt.Errorf("SSH manager not available")
    }

    // Get SSH connection
    conn, err := sshManager.GetConnection(machine.Hostname, machine.Port, machine.User)
    if err != nil {
        return fmt.Errorf("failed to get SSH connection: %w", err)
    }
    defer sshManager.ReturnConnection(conn)

    // Collect application facts
    facts := map[string]interface{}{
        "services": m.getServiceStatus(conn),
        "processes": m.getProcessInfo(conn),
        "packages": m.getPackageInfo(conn),
    }

    // Add facts to collection
    for key, value := range facts {
        collection.Facts[key] = value
    }

    return nil
}
```

### Fact Export Implementation

```go
// ExportFacts exports facts directly to file
func (m *Manager) ExportFacts(ctx context.Context, machineIDs []string, format, outputPath string) error {
    m.logger.Info("Exporting facts", map[string]interface{}{
        "machines": len(machineIDs),
        "format":   format,
        "output":   outputPath,
    })

    // Collect facts for the specified machines
    var allFacts []*spookytypesfacts.FactCollection
    for _, machineID := range machineIDs {
        facts, err := m.CollectFacts(ctx, &spookytypes.Machine{Hostname: machineID})
        if err != nil {
            m.logger.Error("Failed to collect facts for machine", err, map[string]interface{}{
                "machine": machineID,
            })
            return fmt.Errorf("failed to collect facts for machine %s: %w", machineID, err)
        }
        allFacts = append(allFacts, facts)
    }

    // Export based on format
    switch format {
    case "json":
        return m.exportToJSON(allFacts, outputPath)
    case "hcl":
        return m.exportToHCL(allFacts, outputPath)
    default:
        return fmt.Errorf("unsupported export format: %s", format)
    }
}

// GetFacts retrieves facts for a specific machine (collects on demand)
func (m *Manager) GetFacts(ctx context.Context, machineID string) (*spookytypesfacts.FactCollection, error) {
    m.logger.Info("Getting facts for machine", map[string]interface{}{
        "machine": machineID,
    })

    // Facts are collected on demand since there's no persistent storage
    return m.CollectFacts(ctx, &spookytypes.Machine{Hostname: machineID})
}
```

### Fact Validation Implementation

```go
// ValidateFacts validates a fact collection
func (m *Manager) ValidateFacts(ctx context.Context, facts *spookytypes.FactCollection) (*spookytypes.ValidationResult, error) {
    m.logger.Info("Validating facts", map[string]interface{}{
        "machine": facts.Machine,
        "facts":   len(facts.Facts),
    })

    var errors []spookyschemas.SchemaError
    var warnings []spookyschemas.SchemaError

    // Validate fact collection structure
    if err := m.validateFactCollection(facts); err != nil {
        errors = append(errors, spookyschemas.SchemaError{
            Message: err.Error(),
        })
    }

    // Validate individual facts
    for key, value := range facts.Facts {
        if err := m.validateFact(key, value); err != nil {
            errors = append(errors, spookyschemas.SchemaError{
                Message: fmt.Sprintf("fact[%s]: %s", key, err.Error()),
            })
        }
    }

    return &spookytypes.ValidationResult{
        Valid:    len(errors) == 0,
        Errors:   errors,
        Warnings: warnings,
    }, nil
}

func (m *Manager) validateFactCollection(facts *spookytypes.FactCollection) error {
    if facts == nil {
        return fmt.Errorf("fact collection cannot be nil")
    }

    if facts.Machine == "" {
        return fmt.Errorf("machine name is required")
    }

    if facts.Facts == nil {
        return fmt.Errorf("facts map cannot be nil")
    }

    return nil
}

func (m *Manager) validateFact(key string, value interface{}) error {
    if key == "" {
        return fmt.Errorf("fact key cannot be empty")
    }

    if value == nil {
        return fmt.Errorf("fact value cannot be nil")
    }

    // Validate fact key format
    if !regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`).MatchString(key) {
        return fmt.Errorf("invalid fact key format: %s", key)
    }

    return nil
}
```

## Type Definitions

### Fact Types

```go
// FactCollection represents a collection of facts for a machine
type FactCollection struct {
    // Machine hostname
    Machine string `json:"machine" hcl:"machine"`

    // Collection timestamp
    Timestamp time.Time `json:"timestamp" hcl:"timestamp"`

    // Facts map
    Facts map[string]interface{} `json:"facts" hcl:"facts"`

    // Collection metadata
    Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// Fact represents a single fact
type Fact struct {
    // Fact key
    Key string `json:"key" hcl:"key"`

    // Fact value
    Value interface{} `json:"value" hcl:"value"`

    // Fact type
    Type string `json:"type" hcl:"type"`

    // Fact description
    Description string `json:"description,omitempty" hcl:"description,optional"`

    // Fact metadata
    Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// SystemFacts represents system-related facts
type SystemFacts struct {
    // Operating system information
    OS *OSInfo `json:"os" hcl:"os"`

    // Hostname
    Hostname string `json:"hostname" hcl:"hostname"`

    // System uptime
    Uptime time.Duration `json:"uptime" hcl:"uptime"`

    // Memory information
    Memory *MemoryInfo `json:"memory" hcl:"memory"`

    // Disk information
    Disk *DiskInfo `json:"disk" hcl:"disk"`

    // CPU information
    CPU *CPUInfo `json:"cpu" hcl:"cpu"`

    // Kernel information
    Kernel *KernelInfo `json:"kernel" hcl:"kernel"`
}

// NetworkFacts represents network-related facts
type NetworkFacts struct {
    // Network interfaces
    Interfaces []*NetworkInterface `json:"interfaces" hcl:"interfaces"`

    // Network routes
    Routes []*NetworkRoute `json:"routes" hcl:"routes"`

    // DNS information
    DNS *DNSInfo `json:"dns" hcl:"dns"`
}

// ApplicationFacts represents application-related facts
type ApplicationFacts struct {
    // Service status
    Services []*ServiceStatus `json:"services" hcl:"services"`

    // Process information
    Processes []*ProcessInfo `json:"processes" hcl:"processes"`

    // Package information
    Packages []*PackageInfo `json:"packages" hcl:"packages"`
}
```

### Fact Configuration Types

```go
// FactCollector represents a fact collector
type FactCollector struct {
    // Collector name
    Name string `json:"name" hcl:"name"`

    // Collector description
    Description string `json:"description,omitempty" hcl:"description,optional"`

    // Collector type
    Type string `json:"type" hcl:"type"`

    // Collector configuration
    Config map[string]interface{} `json:"config,omitempty" hcl:"config,optional"`

    // Collector enabled
    Enabled bool `json:"enabled" hcl:"enabled"`

    // Collector metadata
    Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// FactFilter represents a fact filter
type FactFilter struct {
    // Filter by machines
    Machines []string `json:"machines,omitempty" hcl:"machines,optional"`

    // Filter by fact keys
    Keys []string `json:"keys,omitempty" hcl:"keys,optional"`

    // Filter by fact types
    Types []string `json:"types,omitempty" hcl:"types,optional"`

    // Complex filter query
    Query string `json:"query,omitempty" hcl:"query,optional"`
}
```

## CLI Commands

### Facts Gather Command

```bash
# Gather facts from all machines
spooky facts gather ./my-project

# Gather facts from specific machines
spooky facts gather ./my-project --machine web-server

# Gather facts by tags
spooky facts gather ./my-project --tags "environment=production"

# Gather facts in parallel
spooky facts gather ./my-project --parallel 4

# Gather facts with timeout
spooky facts gather ./my-project --timeout 300
```

### Facts List Command

```bash
# List all facts
spooky facts list ./my-project

# List facts with details
spooky facts list ./my-project --verbose

# List facts for specific machines
spooky facts list ./my-project --machine web-server

# List facts by keys
spooky facts list ./my-project --key "os,hostname,uptime"
```

### Facts Validate Command

```bash
# Validate facts
spooky facts validate ./my-project

# Validate with detailed output
spooky facts validate ./my-project --verbose

# Validate specific facts
spooky facts validate ./my-project --machine web-server
```

### Facts Export Command

```bash
# Export facts to JSON
spooky facts export ./my-project --json --output facts.json

# Export with filtering
spooky facts export ./my-project --json --output production-facts.json --tags "environment=production"
```

## Integration Examples

### Basic Fact Collection

```go
// Fact collection example
func collectFacts(projectPath string, machines []spookytypes.Machine) error {
    ctx := context.Background()
    
    // Create fact manager
    manager := spookyfacts.NewManager(logger, validator, sshManager, schemaValidator)
    
    // Collect facts from each machine
    for _, machine := range machines {
        // Collect facts
        facts, err := manager.CollectFacts(ctx, machine)
        if err != nil {
            fmt.Printf("Failed to collect facts from %s: %v\n", machine.Hostname, err)
            continue
        }
        
        fmt.Printf("Collected %d facts from %s\n", len(facts.Facts), machine.Hostname)
    }
    
    return nil
}
```

### Fact Retrieval and Validation

```go
// Fact retrieval and validation example
func retrieveAndValidateFacts(projectPath string, machine string) error {
    ctx := context.Background()
    
    // Create fact manager
    manager := spookyfacts.NewManager(logger, validator, sshManager, schemaValidator)
    
    // Get facts
    facts, err := manager.GetFacts(ctx, machine)
    if err != nil {
        return fmt.Errorf("failed to get facts: %w", err)
    }
    
    // Validate facts
    result, err := manager.ValidateFacts(ctx, facts)
    if err != nil {
        return fmt.Errorf("failed to validate facts: %w", err)
    }
    
    if !result.Valid {
        fmt.Println("Fact validation failed:")
        for _, error := range result.Errors {
            fmt.Printf("  - %s\n", error.Message)
        }
        return fmt.Errorf("fact validation failed")
    }
    
    fmt.Printf("Retrieved and validated %d facts from %s\n", len(facts.Facts), machine)
    return nil
}
```

## Error Handling

### Fact Errors

```go
// Error handling example
func handleFactError(err error) {
    if err == nil {
        return
    }
    
    // Check for specific error types
    switch {
    case strings.Contains(err.Error(), "failed to collect facts"):
        fmt.Println("Fact collection failed - check SSH connectivity")
    case strings.Contains(err.Error(), "failed to export facts"):
        fmt.Println("Fact export failed - check file permissions and format")
    case strings.Contains(err.Error(), "failed to collect facts"):
        fmt.Println("Fact collection failed - check SSH connectivity")
    case strings.Contains(err.Error(), "SSH connection failed"):
        fmt.Println("SSH connection failed - check SSH configuration")
    case strings.Contains(err.Error(), "fact validation failed"):
        fmt.Println("Fact validation failed - check fact format")
    default:
        fmt.Printf("Fact error: %v\n", err)
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
    
    fmt.Println("Fact validation failed:")
    for _, err := range result.Errors {
        fmt.Printf("  - %s\n", err.Message)
    }
    
    for _, warning := range result.Warnings {
        fmt.Printf("  Warning: %s\n", warning.Message)
    }
    
    return fmt.Errorf("fact validation failed with %d errors", len(result.Errors))
}
```

## Performance Considerations

### Parallel Processing

The facts system supports parallel processing:

- Multiple machines can be processed concurrently
- Configurable parallel worker count
- Thread-safe fact operations

### Resource Management

The facts system manages resources efficiently:

- SSH connections are pooled and reused
- Memory usage is optimized for large fact sets
- Timeouts prevent hanging operations

## Troubleshooting

### Common Issues

1. **SSH Connection Failures**: Check machine connectivity and SSH configuration
2. **Authentication Errors**: Verify SSH key permissions and user access
3. **Fact Collection Failures**: Check command execution permissions
4. **Export Errors**: Check file permissions and export format
5. **Timeout Issues**: Adjust collection timeouts for slow machines

### Debug Information

The facts system provides comprehensive logging for debugging:

```go
// Enable debug logging
logger.SetLevel(spookytypes.LogLevelDebug)

// Check SSH configuration
fmt.Printf("SSH config: %+v\n", sshConfig)

// Validate fact configuration
err := validateFact(fact)
if err != nil {
    fmt.Printf("Fact validation error: %v\n", err)
}
```

## Future Enhancements

### Planned Features

1. **Parallel Processing**: Implement parallel fact collection
2. **Advanced Filtering**: Support complex fact filtering
3. **Fact Discovery**: Add automatic fact discovery
4. **Fact Caching**: Add fact caching for performance
5. **Fact Versioning**: Add fact versioning support

### Integration Enhancements

1. **Actions Integration**: Use facts in action execution
2. **Variables Integration**: Use facts in variable resolution
3. **Templates Integration**: Use facts in template rendering
4. **Advanced Collection**: Improve fact collection reliability
