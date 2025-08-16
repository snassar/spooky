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
7. **Parallel Processing**: Support for parallel fact collection
8. **Filtering**: Support for machine, tag, and group filtering

### ⚠️ Known Issues

1. **SSH-Based Collection**: SSH-based fact collection has implementation issues
2. **Persistent Storage**: No persistent storage - facts are only stored in memory during export
3. **Fact History**: No historical fact tracking or comparison
4. **Import Functionality**: No fact import capabilities

### 🔄 In Progress

1. **SSH Collection Fixes**: Addressing SSH-based fact collection issues
2. **Storage Improvements**: Implementing persistent storage options
3. **Collection Enhancements**: Improving fact collection reliability

## Implementation Details

### Fact Collection System

The facts system uses a single `SystemFactCollector` that collects facts using SSH commands:

```go
type SystemFactCollector struct {
    name string
}

func (c *SystemFactCollector) Collect(ctx context.Context, machine *spookytypes.Machine) (*spookytypesfacts.FactCollection, error) {
    // Get machine ID from /etc/machine-id
    machineID, err := c.getMachineID(machine)
    if err != nil {
        return nil, fmt.Errorf("failed to get machine ID: %w", err)
    }
    
    // Collect system facts using SSH commands
    facts, err := c.collectSystemFacts(machine)
    if err != nil {
        return nil, fmt.Errorf("failed to collect system facts: %w", err)
    }
    
    return &spookytypesfacts.FactCollection{
        MachineID:   machineID,
        CollectedAt: time.Now(),
        Facts:       facts,
    }, nil
}
```

**Collected Facts:**
- **OS Facts**: Operating system information (name, version, architecture)
- **Hardware Facts**: CPU, memory, disk information
- **Network Facts**: Network interfaces and configuration
- **Load Average**: System load information
- **Process Facts**: Basic process information

### Memory Storage Implementation

Facts are gathered directly and exported without intermediate storage:

```go
type MemoryFactStorage struct {
    mutex sync.RWMutex
}

func (s *MemoryFactStorage) GetStats() (map[string]interface{}, error) {
    return map[string]interface{}{
        "storage_type": "in_memory_only",
        "description":  "Facts are stored in memory for the duration of operations only",
    }, nil
}

func (s *MemoryFactStorage) ExportToJSON(facts map[string]*spookytypesfacts.FactCollection, outputPath string) error {
    // Export facts directly to JSON format
    data, err := json.MarshalIndent(facts, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal facts to JSON: %w", err)
    }
    
    return os.WriteFile(outputPath, data, 0644)
}

func (s *MemoryFactStorage) ExportToHCL(facts map[string]*spookytypesfacts.FactCollection, outputPath string) error {
    // Export facts directly to HCL format
    // Implementation follows facts-structure.schema.hcl
    return nil
}
```

### CLI Integration

Facts commands integrate with the CLI system:

```go
// Facts export command implementation
func handleFactsExport(projectPath, format, outputPath string, machines, tags, groups []string) error {
    ctx := context.Background()
    
    // Initialize facts manager
    manager := spookyfacts.NewManager(collector, validator, logger)
    
    // Load project and get machines
    project, err := loadProject(projectPath)
    if err != nil {
        return fmt.Errorf("failed to load project: %w", err)
    }
    
    // Filter machines based on criteria
    targetMachines := filterMachines(project.Machines, machines, tags, groups)
    
    // Collect facts in parallel
    facts := make(map[string]*spookytypesfacts.FactCollection)
    var wg sync.WaitGroup
    results := make(chan error, len(targetMachines))
    
    for _, machine := range targetMachines {
        wg.Add(1)
        go func(m *spookytypes.Machine) {
            defer wg.Done()
            
            factCollection, err := manager.CollectFacts(ctx, m)
            if err != nil {
                results <- fmt.Errorf("failed to collect facts for %s: %w", m.Hostname, err)
                return
            }
            
            facts[m.Hostname] = factCollection
        }(machine)
    }
    
    wg.Wait()
    close(results)
    
    // Check for errors
    for err := range results {
        if err != nil {
            return err
        }
    }
    
    // Export facts
    return manager.ExportFacts(ctx, facts, format, outputPath)
}
```

## Type Definitions

### Fact Collection Types

```go
// FactCollection represents a collection of facts for a machine
type FactCollection struct {
    // Machine ID (32-character hex string from /etc/machine-id)
    MachineID string `json:"machine_id" hcl:"machine_id"`
    
    // Collection timestamp
    CollectedAt time.Time `json:"collected_at" hcl:"collected_at"`
    
    // Collection of facts for this machine
    Facts *Facts `json:"facts" hcl:"facts"`
    
    // Metadata about the collection
    Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// Facts represents the actual fact data
type Facts struct {
    // System facts
    System *SystemFacts `json:"system" hcl:"system"`
    
    // Hardware facts
    Hardware *HardwareFacts `json:"hardware" hcl:"hardware"`
    
    // Network facts
    Network *NetworkFacts `json:"network" hcl:"network"`
    
    // Custom facts
    Custom map[string]interface{} `json:"custom,omitempty" hcl:"custom,optional"`
}

// SystemFacts represents system information
type SystemFacts struct {
    // Operating system information
    OS *OSFacts `json:"os" hcl:"os"`
    
    // Load average
    LoadAverage *LoadAverageFacts `json:"load_average" hcl:"load_average"`
    
    // Process information
    Processes *ProcessFacts `json:"processes" hcl:"processes"`
}

// OSFacts represents operating system information
type OSFacts struct {
    Name         string `json:"name" hcl:"name"`
    Version      string `json:"version" hcl:"version"`
    Architecture string `json:"architecture" hcl:"architecture"`
    Kernel       string `json:"kernel" hcl:"kernel"`
    Distribution string `json:"distribution" hcl:"distribution"`
}
```

### Fact Storage Types

```go
// FactStorage provides storage operations for fact collections
type FactStorage interface {
    // Store stores facts for a machine
    Store(ctx context.Context, machineID string, facts *FactCollection) error
    
    // Get retrieves facts for a machine
    Get(ctx context.Context, machineID string) (*FactCollection, error)
    
    // List lists all machine IDs with stored facts
    List(ctx context.Context) ([]string, error)
    
    // Clear removes all facts from storage
    Clear(ctx context.Context) error
    
    // GetStats returns storage statistics for debugging
    GetStats() (map[string]interface{}, error)
}

// FactCollector collects facts from a machine
type FactCollector interface {
    // Collect collects facts from the given machine
    Collect(ctx context.Context, machine interface{}) (*FactCollection, error)
    
    // GetName returns the collector name
    GetName() string
}
```

## Error Handling

### Fact Collection Errors

```go
// FactCollectionError represents fact collection errors
type FactCollectionError struct {
    MachineID string `json:"machine_id" hcl:"machine_id"`
    Error     string `json:"error" hcl:"error"`
    Details   string `json:"details,omitempty" hcl:"details,optional"`
}

// FactValidationError represents fact validation errors
type FactValidationError struct {
    Field   string `json:"field" hcl:"field"`
    Message string `json:"message" hcl:"message"`
    Value   string `json:"value,omitempty" hcl:"value,optional"`
}
```

### Validation Implementation

```go
// ValidateFacts validates facts against schema
func (m *Manager) ValidateFacts(_ context.Context, facts *spookytypesfacts.FactCollection) (*spookytypes.ValidationResult, error) {
    if facts == nil {
        return &spookytypes.ValidationResult{
            Valid:    false,
            Errors:   []spookyschemas.SchemaError{{Message: "facts cannot be nil"}},
            Warnings: []spookyschemas.SchemaError{},
        }, nil
    }

    // Basic validation
    var errors []spookyschemas.SchemaError
    var warnings []spookyschemas.SchemaError

    // Validate machine ID
    if facts.MachineID == "" {
        errors = append(errors, spookyschemas.SchemaError{Message: "machine_id is required"})
    } else if !isValidMachineID(facts.MachineID) {
        errors = append(errors, spookyschemas.SchemaError{Message: "machine_id must be a 32-character hexadecimal string"})
    }

    // Validate collection timestamp
    if facts.CollectedAt.IsZero() {
        errors = append(errors, spookyschemas.SchemaError{Message: "collected_at is required"})
    }

    // Validate facts structure
    if facts.Facts == nil {
        errors = append(errors, spookyschemas.SchemaError{Message: "facts structure is required"})
    } else {
        // Validate system facts
        if facts.Facts.System == nil {
            errors = append(errors, spookyschemas.SchemaError{Message: "system facts are required"})
        } else {
            if facts.Facts.System.OS == nil {
                errors = append(errors, spookyschemas.SchemaError{Message: "system.os facts are required"})
            }
            if facts.Facts.System.Hardware == nil {
                errors = append(errors, spookyschemas.SchemaError{Message: "system.hardware facts are required"})
            }
            if facts.Facts.System.Network == nil {
                errors = append(errors, spookyschemas.SchemaError{Message: "system.network facts are required"})
            }
        }
    }

    valid := len(errors) == 0

    return &spookytypes.ValidationResult{
        Valid:    valid,
        Errors:   errors,
        Warnings: warnings,
    }, nil
}
```

## CLI Commands

### Facts Export Command

```bash
# Export facts to JSON format
spooky facts export ./my-project --format json --output facts.json

# Export facts to HCL format
spooky facts export ./my-project --format hcl --output facts.hcl

# Export facts for specific machines
spooky facts export ./my-project --machines web-server,app-server --format json --output selected-facts.json

# Export facts for machines with specific tags
spooky facts export ./my-project --tags production,web --format json --output production-facts.json

# Export facts for machines in specific groups
spooky facts export ./my-project --groups webservers,databases --format json --output group-facts.json
```

### Facts Validation Command

```bash
# Validate facts collection
spooky facts validate ./my-project

# Validate facts with comparison
spooky facts validate ./my-project --compare
```

## Integration Examples

### Basic Fact Collection

```go
// Basic fact collection example
func collectFactsForMachine(hostname string) error {
    ctx := context.Background()
    
    // Create fact collector
    collector := spookyfacts.NewSystemFactCollector()
    
    // Create machine object
    machine := &spookytypes.Machine{
        Hostname: hostname,
        Port:     22,
        User:     "admin",
    }
    
    // Collect facts
    facts, err := collector.Collect(ctx, machine)
    if err != nil {
        return fmt.Errorf("failed to collect facts: %w", err)
    }
    
    // Print facts
    fmt.Printf("Machine ID: %s\n", facts.MachineID)
    fmt.Printf("OS: %s %s\n", facts.Facts.System.OS.Name, facts.Facts.System.OS.Version)
    fmt.Printf("Architecture: %s\n", facts.Facts.System.OS.Architecture)
    
    return nil
}
```

### Parallel Fact Collection

```go
// Parallel fact collection example
func collectFactsForMultipleMachines(machines []*spookytypes.Machine) error {
    ctx := context.Background()
    
    // Create fact manager
    manager := spookyfacts.NewManager(collector, validator, logger)
    
    // Collect facts in parallel
    err := manager.CollectAndStoreFactsParallel(ctx, machines, 4)
    if err != nil {
        return fmt.Errorf("failed to collect facts: %w", err)
    }
    
    return nil
}
```

### Fact Export

```go
// Fact export example
func exportFactsToJSON(facts map[string]*spookytypesfacts.FactCollection, outputPath string) error {
    ctx := context.Background()
    
    // Create fact manager
    manager := spookyfacts.NewManager(collector, validator, logger)
    
    // Export facts
    err := manager.ExportFacts(ctx, facts, "json", outputPath)
    if err != nil {
        return fmt.Errorf("failed to export facts: %w", err)
    }
    
    return nil
}
```

## Current Limitations

### Storage Characteristics

1. **Ephemeral Storage**: Facts are only stored in memory during export (this is intentional)
2. **No Fact History**: No historical fact tracking or comparison (facts are gathered fresh each time)
3. **No Persistent Storage**: Facts are not saved to disk (memory-only during export)
4. **No Import Functionality**: Facts cannot be imported (only export is supported)

### Collection Limitations

1. **Single Collector**: Only system fact collector available
2. **No Custom Facts**: No support for custom fact collection
3. **No Fact Caching**: No caching of previously collected facts
4. **No Incremental Collection**: Always collects all facts

### Integration Limitations

1. **No Template Integration**: Facts not integrated with template system
2. **No Variable Integration**: Facts not used in variable resolution
3. **No Action Integration**: Facts not integrated with action system
4. **No Real-time Monitoring**: No real-time fact monitoring

## Future Enhancements

### Planned Features

1. **Persistent Storage**: Long-term fact storage in databases
2. **Fact History**: Historical fact tracking and comparison
3. **Advanced Collectors**: Custom fact collectors for specific data
4. **Fact Validation**: Enhanced validation and schema checking
5. **Fact Import**: Import facts from external sources
6. **Fact Comparison**: Compare facts across machines and time periods

### Integration Enhancements

1. **Template Integration**: Use facts in template rendering
2. **Variable Integration**: Use facts in variable resolution
3. **Action Integration**: Use facts in action conditions
4. **Real-time Monitoring**: Real-time fact collection and monitoring

## Summary

The facts system provides basic fact collection and export capabilities with some limitations. The system is functional for basic use cases but has known issues with SSH-based collection that need to be addressed.

**Status**: ⚠️ **Partially Implemented** - Basic functionality exists but SSH-based collection has issues that need to be resolved.
