# Facts System API Reference

## Overview

This document provides a technical reference for the spooky facts system APIs and implementation details. It covers core interfaces, type definitions, implementation patterns, and integration details for developers working with the facts system.

## Core Interfaces

### FactManager Interface

The `FactManager` interface provides the primary entry point for fact management operations:

```go
type FactManager interface {
    // CollectFacts collects facts from the given machine
    CollectFacts(ctx context.Context, machine *spookytypes.Machine) (*FactCollection, error)
    
    // StoreFacts stores facts for a machine
    StoreFacts(ctx context.Context, machineID string, facts *FactCollection) error
    
    // GetFacts retrieves facts for a machine
    GetFacts(ctx context.Context, machineID string) (*FactCollection, error)
    
    // ListFacts lists all machines with stored facts
    ListFacts(ctx context.Context) ([]string, error)
    
    // DeleteFacts deletes facts for a machine
    DeleteFacts(ctx context.Context, machineID string) error
    
    // ValidateFacts validates facts against schema
    ValidateFacts(ctx context.Context, facts *FactCollection) (*spookytypes.ValidationResult, error)
    
    // ExportFacts exports facts to the given format
    ExportFacts(ctx context.Context, machineIDs []string, format string, outputPath string) error
    
    // ImportFacts imports facts from the given format
    ImportFacts(ctx context.Context, format string, inputPath string) error
}
```

### FactCollector Interface

The `FactCollector` interface provides fact collection capabilities:

```go
type FactCollector interface {
    // Collect collects facts from the given machine
    Collect(ctx context.Context, machine *spookytypes.Machine) (*FactCollection, error)
    
    // GetName returns the collector name
    GetName() string
}
```

### FactCollector Interface

The `FactCollector` interface provides fact collection capabilities:

```go
type FactCollector interface {
    // Collect collects facts from the given machine
    Collect(ctx context.Context, machine *spookytypes.Machine) (*FactCollection, error)
    
    // GetName returns the collector name
    GetName() string
}
```

### FactsIntegration Interface

The `FactsIntegration` interface provides integration capabilities:

```go
type FactsIntegration interface {
    // CollectFacts collects facts from the given machine
    CollectFacts(ctx context.Context, machine *spookytypes.Machine) (interface{}, error)
    
    // ValidateFacts validates facts against schema
    ValidateFacts(ctx context.Context, facts interface{}) (*spookytypes.ValidationResult, error)
    
    // ExportFacts exports facts to the given format
    ExportFacts(ctx context.Context, machineIDs []string, format string, outputPath string) error
}
```

## Core Types

### FactCollection

The `FactCollection` type represents a collection of facts for a machine:

```go
type FactCollection struct {
    // Machine ID (32-character hex string from /etc/machine-id)
    MachineID string `json:"machine_id" hcl:"machine_id"`
    
    // Collection timestamp
    CollectedAt time.Time `json:"collected_at" hcl:"collected_at"`
    
    // Collection of facts for this machine
    Facts *spookytypesfacts.Facts `json:"facts" hcl:"facts"`
    
    // Metadata about the collection
    Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}
```

### Facts Structure

The `Facts` type represents the complete facts structure for a machine:

```go
type Facts struct {
    // System-level facts (from gopsutil)
    System *SystemFacts `json:"system" hcl:"system"`

    // Enhanced system information
    Enhanced *EnhancedFacts `json:"enhanced,omitempty" hcl:"enhanced,optional"`

    // Application facts
    Applications *ApplicationFacts `json:"applications,omitempty" hcl:"applications,optional"`

    // Deployment facts
    Deployment *DeploymentFacts `json:"deployment,omitempty" hcl:"deployment,optional"`

    // Environment facts
    Environment *EnvironmentFacts `json:"environment,omitempty" hcl:"environment,optional"`

    // Monitoring facts
    Monitoring *MonitoringFacts `json:"monitoring,omitempty" hcl:"monitoring,optional"`

    // Custom facts (user-defined)
    Custom map[string]interface{} `json:"custom,omitempty" hcl:"custom,optional"`
}
```

### SystemFacts

The `SystemFacts` type represents system-level facts from gopsutil:

```go
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
```

### HardwareFacts

The `HardwareFacts` type represents hardware information:

```go
type HardwareFacts struct {
    CPU    *CPUFacts    `json:"cpu" hcl:"cpu"`
    Memory *MemoryFacts `json:"memory" hcl:"memory"`
    Disks  []*DiskFacts `json:"disks" hcl:"disks"`
    DiskIO *DiskIOFacts `json:"disk_io,omitempty" hcl:"disk_io,optional"`
}
```

### CPUFacts

The `CPUFacts` type represents CPU information:

```go
type CPUFacts struct {
    Cores        int              `json:"cores" hcl:"cores"`
    Model        string           `json:"model" hcl:"model"`
    Frequency    float64          `json:"frequency,omitempty" hcl:"frequency,optional"`
    Architecture string           `json:"architecture,omitempty" hcl:"architecture,optional"`
    Vendor       string           `json:"vendor,omitempty" hcl:"vendor,optional"`
    Times        *CPUTimes        `json:"times,omitempty" hcl:"times,optional"`
    Percent      float64          `json:"percent,omitempty" hcl:"percent,optional"`
    CoresDetail  []*CPUCoreDetail `json:"cores_detail,omitempty" hcl:"cores_detail,optional"`
}
```

### MemoryFacts

The `MemoryFacts` type represents memory information:

```go
type MemoryFacts struct {
    Total     int64   `json:"total" hcl:"total"`
    Available int64   `json:"available" hcl:"available"`
    Used      int64   `json:"used" hcl:"used"`
    Free      int64   `json:"free" hcl:"free"`
    Percent   float64 `json:"percent,omitempty" hcl:"percent,optional"`
}
```

### NetworkFacts

The `NetworkFacts` type represents network information:

```go
type NetworkFacts struct {
    Interfaces []*NetworkInterface `json:"interfaces" hcl:"interfaces"`
    Connections []*NetworkConnection `json:"connections,omitempty" hcl:"connections,optional"`
}
```

## Implementation Details

### Fact Collection Process

The fact collection process follows these steps:

1. **Machine Validation**: Validate the machine configuration and connectivity
2. **SSH Connection**: Establish SSH connection to the target machine
3. **Fact Collection**: Use gopsutil to collect system information
4. **Data Processing**: Convert and validate collected data
5. **In-Memory Storage**: Store facts in memory for the duration of the action run

### In-Memory Implementation

The in-memory facts implementation provides:

- **Fast Access**: Facts available immediately after collection
- **No I/O Overhead**: No disk operations during fact access
- **Session-Based**: Facts persist for the duration of the action run
- **Efficient Memory Usage**: Optimized memory allocation for fact storage
- **Parallel Collection**: Efficient parallel fact collection across machines

### Validation Rules

The facts validation system enforces these rules:

- **Machine ID**: Must be a 32-character hexadecimal string
- **Required Facts**: System, hardware, and network facts are required
- **Data Types**: All numeric fields must be valid numbers
- **String Lengths**: String fields must not exceed maximum lengths
- **Custom Facts**: Custom facts must be valid JSON

### Error Handling

The facts system provides comprehensive error handling:

- **Collection Errors**: Network timeouts, SSH failures, permission issues
- **Memory Errors**: Memory allocation failures, out of memory conditions
- **Validation Errors**: Schema violations, data type mismatches
- **Integration Errors**: Interface contract violations

## CLI Integration

### Command Structure

The facts CLI provides these commands:

```bash
spooky facts gather [project-path]     # Collect facts from machines
spooky facts export [project-path]     # Export facts to file (includes internal validation)
```

### Command Options

Common options across facts commands:

- `--machine`: Target specific machine
- `--parallel`: Number of parallel workers
- `--timeout`: Collection timeout
- `--format`: Export format (json, hcl)
- `--output`: Output file path

### Error Reporting

CLI commands provide detailed error reporting:

- **Validation Errors**: Schema validation failures with field details
- **Collection Errors**: Machine-specific collection failures
- **Memory Errors**: Memory allocation and management failures
- **Network Errors**: Connectivity and timeout issues

## Performance Considerations

### Collection Performance

- **Parallel Collection**: Support for parallel fact collection across machines
- **Timeout Management**: Configurable timeouts to prevent hanging
- **Retry Logic**: Automatic retry for transient failures
- **Caching**: Optional caching of frequently accessed facts

### Memory Performance

- **Efficient Allocation**: Optimized memory allocation for fact storage
- **Garbage Collection**: Proper cleanup of temporary objects
- **Memory Pooling**: Reuse of memory structures for better performance
- **Memory Monitoring**: Track memory usage during fact collection

### Memory Usage

- **Streaming**: Stream large fact collections to avoid memory issues
- **Chunking**: Process facts in chunks for large datasets
- **Garbage Collection**: Proper cleanup of temporary objects

## Security Considerations

### Data Protection

- **Memory Protection**: Facts stored in memory are not persisted to disk
- **Access Control**: Facts are only accessible during the action run
- **Audit Logging**: Logging of fact collection and usage
- **Data Sanitization**: Removal of sensitive information before collection

### Network Security

- **SSH Authentication**: Secure SSH authentication for fact collection
- **Connection Encryption**: Encrypted connections to target machines
- **Certificate Validation**: Validation of SSH certificates
- **Timeout Protection**: Protection against hanging connections

## Testing

### Unit Testing

The facts system includes comprehensive unit tests:

- **Interface Testing**: Test all interface implementations
- **Mock Testing**: Mock dependencies for isolated testing
- **Error Testing**: Test error conditions and edge cases
- **Validation Testing**: Test validation rules and schemas

### Integration Testing

Integration tests cover:

- **End-to-End Workflows**: Complete fact collection workflows
- **Memory Integration**: In-memory fact storage testing
- **CLI Integration**: Command-line interface testing
- **Error Scenarios**: Real-world error scenarios

### Performance Testing

Performance tests validate:

- **Collection Performance**: Fact collection speed and efficiency
- **Memory Performance**: Memory allocation and access performance
- **Memory Usage**: Memory consumption patterns
- **Scalability**: Performance with large numbers of machines

## Best Practices

### Fact Collection

- **Regular Collection**: Collect facts regularly for up-to-date information
- **Parallel Processing**: Use parallel collection for multiple machines
- **Error Handling**: Implement proper error handling and retry logic
- **Validation**: Always validate collected facts before storage

### Memory Management

- **Memory Monitoring**: Monitor memory usage during fact collection
- **Garbage Collection**: Proper cleanup of temporary objects
- **Memory Pooling**: Reuse memory structures for better performance
- **Memory Limits**: Set appropriate memory limits for large collections

### Integration

- **Interface Compliance**: Ensure implementations comply with interfaces
- **Error Propagation**: Properly propagate errors through the system
- **Logging**: Implement comprehensive logging for debugging
- **Configuration**: Use configuration for customization

## Troubleshooting

### Common Issues

- **Collection Failures**: Network issues, SSH problems, permission errors
- **Memory Issues**: Memory allocation failures, out of memory conditions
- **Validation Errors**: Schema violations, data type issues
- **Performance Issues**: Slow collection, high memory usage

### Debugging Techniques

- **Logging**: Enable debug logging for detailed information
- **Validation**: Use validation commands to check fact integrity
- **Export/Import**: Use export/import for data analysis
- **Storage Inspection**: Direct inspection of BadgerDB contents

### Performance Optimization

- **Parallel Collection**: Increase parallel workers for faster collection
- **Timeout Tuning**: Adjust timeouts based on network conditions
- **Storage Tuning**: Optimize BadgerDB configuration
- **Memory Management**: Monitor and optimize memory usage
