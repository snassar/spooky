# Machines Inventory API Reference

## Overview

This document provides a comprehensive API reference for the spooky machines inventory system. It covers all interfaces, types, methods, and implementation details for developers working with the machines system.

## Table of Contents

1. [Core Interfaces](#core-interfaces)
2. [Type Definitions](#type-definitions)
3. [Implementation Details](#implementation-details)
4. [Error Handling](#error-handling)
5. [Validation Rules](#validation-rules)
6. [CLI Integration](#cli-integration)
7. [Examples](#examples)

## Core Interfaces

### MachinesIntegration

The primary interface for machine management operations.

```go
type MachinesIntegration interface {
    // LoadMachines loads machines from the given source
    LoadMachines(ctx context.Context, source string) ([]spookytypes.Machine, error)
    
    // ValidateMachines validates machines
    ValidateMachines(ctx context.Context, machines []spookytypes.Machine) (*spookytypes.ValidationResult, error)
    
    // PingMachines pings machines to check connectivity
    PingMachines(ctx context.Context, machines []spookytypes.Machine) ([]spookytypes.MachineStatus, error)
}
```

**Methods:**

#### LoadMachines
```go
LoadMachines(ctx context.Context, source string) ([]spookytypes.Machine, error)
```
Loads machines from a project directory, supporting both `machines.hcl` files and `machines/` directories.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `source`: Project directory path

**Returns:**
- `[]spookytypes.Machine`: Array of loaded machines
- `error`: Error if loading fails

**Behavior:**
- Checks for `machines.hcl` file in the project root
- Checks for `machines/` directory and loads all `.hcl` files
- Combines machines from all sources
- Performs duplicate detection and cross-file validation
- Returns error if no machines found or validation fails

#### ValidateMachines
```go
ValidateMachines(ctx context.Context, machines []spookytypes.Machine) (*spookytypes.ValidationResult, error)
```
Validates machine configurations and returns detailed validation results.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `machines`: Array of machines to validate

**Returns:**
- `*spookytypes.ValidationResult`: Validation results with errors and warnings
- `error`: Error if validation process fails

**Validation Rules:**
- Required fields validation (hostname, host, user)
- SSH key file existence and permissions
- Port number range validation (1-65535)
- Hostname format validation
- Environment-specific validation rules

#### PingMachines
```go
PingMachines(ctx context.Context, machines []spookytypes.Machine) ([]spookytypes.MachineStatus, error)
```
Tests connectivity to machines using progressive connectivity checks.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `machines`: Array of machines to ping

**Returns:**
- `[]spookytypes.MachineStatus`: Array of machine status results
- `error`: Error if ping process fails

**Connectivity Checks:**
1. DNS resolution (for hostnames)
2. ICMP ping (simulated for now)
3. TCP port scan for SSH port
4. SSH connection and authentication (deferred)

### MachineValidator

Interface for machine-specific validation operations.

```go
type MachineValidator interface {
    // ValidateMachines validates machines
    ValidateMachines(ctx context.Context, machines []spookytypes.Machine) (*spookytypes.ValidationResult, error)
    
    // ValidateMachine validates a single machine
    ValidateMachine(ctx context.Context, machine spookytypes.Machine) (*spookytypes.ValidationResult, error)
}
```

**Methods:**

#### ValidateMachines
```go
ValidateMachines(ctx context.Context, machines []spookytypes.Machine) (*spookytypes.ValidationResult, error)
```
Validates multiple machines and aggregates results.

#### ValidateMachine
```go
ValidateMachine(ctx context.Context, machine spookytypes.Machine) (*spookytypes.ValidationResult, error)
```
Validates a single machine configuration.

### MachineLoader

Interface for loading machine configurations from files and directories.

```go
type MachineLoader interface {
    // LoadMachinesFromFile loads machines from a file
    LoadMachinesFromFile(ctx context.Context, filePath string) ([]spookytypes.Machine, error)
    
    // LoadMachinesFromDirectory loads machines from a directory
    LoadMachinesFromDirectory(ctx context.Context, dirPath string) ([]spookytypes.Machine, error)
}
```

**Methods:**

#### LoadMachinesFromFile
```go
LoadMachinesFromFile(ctx context.Context, filePath string) ([]spookytypes.Machine, error)
```
Loads machines from a single HCL file.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `filePath`: Path to the HCL file

**Returns:**
- `[]spookytypes.Machine`: Array of machines from the file
- `error`: Error if loading fails

**Behavior:**
- Parses HCL file using `machines.schema.hcl`
- Extracts machine blocks and attributes
- Adds source file information to machine metadata
- Returns error for invalid HCL syntax

#### LoadMachinesFromDirectory
```go
LoadMachinesFromDirectory(ctx context.Context, dirPath string) ([]spookytypes.Machine, error)
```
Loads machines from all `.hcl` files in a directory.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `dirPath`: Path to the directory

**Returns:**
- `[]spookytypes.Machine`: Array of machines from all files
- `error`: Error if loading fails

**Behavior:**
- Scans directory for `.hcl` files
- Loads machines from each file
- Combines results into single array
- Preserves source file information

## Type Definitions

### Machine

Core machine configuration structure.

```go
type Machine struct {
    // Basic connection information
    Hostname string `hcl:"hostname,label" json:"hostname"`
    Host     string `hcl:"host" json:"host"`
    Port     int    `hcl:"port,optional" json:"port,omitempty"`
    User     string `hcl:"user" json:"user"`
    
    // Authentication
    KeyFile   string `hcl:"key_file,optional" json:"key_file,omitempty"`
    Passphrase string `hcl:"passphrase,optional" json:"passphrase,omitempty"`
    
    // Organization
    Tags   []string `hcl:"tags,optional" json:"tags,omitempty"`
    Groups []string `hcl:"groups,optional" json:"groups,omitempty"`
    Roles  []string `hcl:"roles,optional" json:"roles,omitempty"`
    
    // Resources and metadata
    Resources      *MachineResources `hcl:"resources,block" json:"resources,omitempty"`
    MachineMetadata *MachineMetadata `hcl:"metadata,block" json:"metadata,omitempty"`
}
```

**Fields:**

- **Hostname**: Unique identifier for the machine (HCL label)
- **Host**: IP address or hostname for SSH connection
- **Port**: SSH port (default: 22)
- **User**: SSH username for authentication
- **KeyFile**: Path to SSH private key file
- **Passphrase**: Passphrase for encrypted SSH keys
- **Tags**: Array of tags for categorization
- **Groups**: Array of groups for organization
- **Roles**: Array of roles for automation
- **Resources**: Machine resource specifications
- **MachineMetadata**: Additional machine metadata

### MachineResources

Machine resource specifications for capacity planning.

```go
type MachineResources struct {
    CPUCores    int `hcl:"cpu_cores,optional" json:"cpu_cores,omitempty"`
    MemoryGB    int `hcl:"memory_gb,optional" json:"memory_gb,omitempty"`
    DiskGB      int `hcl:"disk_gb,optional" json:"disk_gb,omitempty"`
    NetworkMbps int `hcl:"network_mbps,optional" json:"network_mbps,omitempty"`
}
```

**Fields:**

- **CPUCores**: Number of CPU cores
- **MemoryGB**: Memory in gigabytes
- **DiskGB**: Disk space in gigabytes
- **NetworkMbps**: Network bandwidth in megabits per second

### MachineMetadata

Additional machine metadata for organization and management.

```go
type MachineMetadata struct {
    // Environment information
    Environment string `hcl:"environment,optional" json:"environment,omitempty"`
    Datacenter  string `hcl:"datacenter,optional" json:"datacenter,omitempty"`
    Rack        string `hcl:"rack,optional" json:"rack,omitempty"`
    Location    string `hcl:"location,optional" json:"location,omitempty"`
    
    // Ownership information
    Owner       string `hcl:"owner,optional" json:"owner,omitempty"`
    Department  string `hcl:"department,optional" json:"department,omitempty"`
    CostCenter  string `hcl:"cost_center,optional" json:"cost_center,omitempty"`
    
    // Operational information
    MaintenanceWindow string `hcl:"maintenance_window,optional" json:"maintenance_window,omitempty"`
    BackupSchedule    string `hcl:"backup_schedule,optional" json:"backup_schedule,omitempty"`
    Monitoring        string `hcl:"monitoring,optional" json:"monitoring,omitempty"`
    Alerting          string `hcl:"alerting,optional" json:"alerting,omitempty"`
    SLA               string `hcl:"sla,optional" json:"sla,omitempty"`
    
    // Custom fields
    CustomFields map[string]string `hcl:"custom_fields,optional" json:"custom_fields,omitempty"`
}
```

**Fields:**

- **Environment**: Environment name (production, staging, development)
- **Datacenter**: Datacenter identifier
- **Rack**: Rack location
- **Location**: Physical location
- **Owner**: Machine owner or team
- **Department**: Department responsible
- **CostCenter**: Cost center for billing
- **MaintenanceWindow**: Scheduled maintenance window
- **BackupSchedule**: Backup schedule
- **Monitoring**: Monitoring system
- **Alerting**: Alerting system
- **SLA**: Service level agreement
- **CustomFields**: Additional custom fields

### MachineStatus

Result of connectivity testing for a machine.

```go
type MachineStatus struct {
    Machine   spookytypes.Machine `json:"machine"`
    Status    string              `json:"status"`    // "online", "offline", "error"
    LastCheck time.Time           `json:"last_check"`
    Error     string              `json:"error,omitempty"`
    Latency   int64               `json:"latency_ms,omitempty"` // in milliseconds
    Details   map[string]interface{} `json:"details,omitempty"`
}
```

**Fields:**

- **Machine**: The machine configuration
- **Status**: Connectivity status ("online", "offline", "error")
- **LastCheck**: Timestamp of last connectivity check
- **Error**: Error message if connectivity failed
- **Latency**: Response time in milliseconds
- **Details**: Additional connectivity details

### MachineCollection

Collection of machines with metadata.

```go
type MachineCollection struct {
    Machines []spookytypes.Machine `json:"machines"`
    Metadata map[string]interface{} `json:"metadata,omitempty"`
    Source   string                 `json:"source,omitempty"`
}
```

**Fields:**

- **Machines**: Array of machines
- **Metadata**: Collection metadata
- **Source**: Source of the collection

## Implementation Details

### Manager Implementation

The `Manager` struct implements the `MachinesIntegration` interface.

```go
type Manager struct {
    logger    spookylogging.Logger
    loader    spookyinterfaces.MachineLoader
    validator spookyinterfaces.MachineValidator
}
```

**Constructor:**
```go
func NewManager(logger spookylogging.Logger, loader spookyinterfaces.MachineLoader, validator spookyinterfaces.MachineValidator) spookyinterfaces.MachinesIntegration
```

**Key Methods:**

#### LoadMachines Implementation
```go
func (m *Manager) LoadMachines(ctx context.Context, projectPath string) ([]spookytypes.Machine, error)
```

**Algorithm:**
1. Check for `machines.hcl` file in project root
2. If found, load machines using loader
3. Check for `machines/` directory
4. If found, load machines from all `.hcl` files
5. Combine machines from all sources
6. Perform duplicate detection
7. Perform cross-file validation
8. Return combined machines or error

#### PingMachines Implementation
```go
func (m *Manager) PingMachines(ctx context.Context, machines []spookytypes.Machine) ([]spookytypes.MachineStatus, error)
```

**Algorithm:**
1. Create status array for results
2. For each machine:
   - Perform DNS resolution (if hostname)
   - Perform ICMP ping (simulated)
   - Perform TCP port scan
   - Calculate latency
   - Create status object
3. Return status array

### Loader Implementation

The `Loader` struct implements the `MachineLoader` interface.

```go
type Loader struct {
    logger spookylogging.Logger
}
```

**Key Methods:**

#### LoadMachinesFromFile Implementation
```go
func (l *Loader) LoadMachinesFromFile(ctx context.Context, filePath string) ([]spookytypes.Machine, error)
```

**Algorithm:**
1. Read HCL file content
2. Parse HCL using `machines.schema.hcl`
3. Extract machine blocks
4. For each machine block:
   - Parse attributes
   - Parse nested blocks (resources, metadata)
   - Add source file information
   - Create machine object
5. Return machine array

#### LoadMachinesFromDirectory Implementation
```go
func (l *Loader) LoadMachinesFromDirectory(ctx context.Context, dirPath string) ([]spookytypes.Machine, error)
```

**Algorithm:**
1. Scan directory for `.hcl` files
2. For each `.hcl` file:
   - Load machines using `LoadMachinesFromFile`
   - Add to combined array
3. Return combined machine array

### Validator Implementation

The `Validator` struct implements the `MachineValidator` interface.

```go
type Validator struct {
    logger spookylogging.Logger
}
```

**Key Methods:**

#### ValidateMachines Implementation
```go
func (v *Validator) ValidateMachines(ctx context.Context, machines []spookytypes.Machine) (*spookytypes.ValidationResult, error)
```

**Algorithm:**
1. Create validation result
2. For each machine:
   - Validate individual machine
   - Collect errors and warnings
3. Perform cross-machine validation
4. Return aggregated results

#### ValidateMachine Implementation
```go
func (v *Validator) ValidateMachine(ctx context.Context, machine spookytypes.Machine) (*spookytypes.ValidationResult, error)
```

**Validation Rules:**
1. Required fields validation
2. Hostname format validation
3. Host validation (IP or hostname)
4. Port range validation
5. SSH key file validation
6. Environment-specific validation

## Error Handling

### Error Types

The machines system defines several error types for different scenarios.

#### MachineError
Base error type for machine operations.

```go
type MachineError struct {
    Message   string                 `json:"message"`
    Code      string                 `json:"code"`
    Details   map[string]interface{} `json:"details,omitempty"`
    Timestamp time.Time              `json:"timestamp"`
}
```

#### MachineValidationError
Error for validation failures.

```go
type MachineValidationError struct {
    MachineError
    Field   string `json:"field"`
    Value   string `json:"value"`
    Rule    string `json:"rule"`
}
```

#### MachineConnectionError
Error for connectivity failures.

```go
type MachineConnectionError struct {
    MachineError
    Hostname string `json:"hostname"`
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Stage    string `json:"stage"` // "dns", "icmp", "tcp", "ssh"
}
```

#### MachineLoadError
Error for loading failures.

```go
type MachineLoadError struct {
    MachineError
    FilePath string `json:"file_path"`
    Line     int    `json:"line,omitempty"`
    Column   int    `json:"column,omitempty"`
}
```

### Error Construction

Error constructors provide consistent error creation:

```go
// Create validation error
func NewMachineValidationError(message, field, value, rule string) *MachineValidationError

// Create connection error
func NewMachineConnectionError(message, hostname, host string, port int, stage string) *MachineConnectionError

// Create load error
func NewMachineLoadError(message, filePath string) *MachineLoadError
```

### Error Handling Patterns

**Validation Error Handling:**
```go
if err := validateMachine(machine); err != nil {
    return &spookytypesmachines.MachineValidationError{
        MachineError: spookytypesmachines.MachineError{
            Message:   err.Error(),
            Code:      "VALIDATION_FAILED",
            Timestamp: time.Now(),
        },
        Field: "hostname",
        Value: machine.Hostname,
        Rule:  "required",
    }
}
```

**Connection Error Handling:**
```go
if err := testConnectivity(machine); err != nil {
    return &spookytypesmachines.MachineConnectionError{
        MachineError: spookytypesmachines.MachineError{
            Message:   err.Error(),
            Code:      "CONNECTION_FAILED",
            Timestamp: time.Now(),
        },
        Hostname: machine.Hostname,
        Host:     machine.Host,
        Port:     machine.Port,
        Stage:    "dns",
    }
}
```

## Validation Rules

### Required Fields

Every machine must have these fields:

- **hostname**: Non-empty string, unique across all files
- **host**: Non-empty string, valid IP or hostname
- **user**: Non-empty string, valid username

### Field Validation

#### Hostname Validation
```go
func validateHostname(hostname string) error {
    if hostname == "" {
        return fmt.Errorf("hostname cannot be empty")
    }
    
    // Check for valid characters
    if !hostnameRegex.MatchString(hostname) {
        return fmt.Errorf("hostname contains invalid characters")
    }
    
    return nil
}
```

#### Host Validation
```go
func validateHost(host string) error {
    if host == "" {
        return fmt.Errorf("host cannot be empty")
    }
    
    // Check if it's an IP address
    if net.ParseIP(host) != nil {
        return nil
    }
    
    // Check if it's a valid hostname
    if !hostnameRegex.MatchString(host) {
        return fmt.Errorf("host is not a valid IP address or hostname")
    }
    
    return nil
}
```

#### Port Validation
```go
func validatePort(port int) error {
    if port < 1 || port > 65535 {
        return fmt.Errorf("port must be between 1 and 65535")
    }
    return nil
}
```

### Environment-Specific Validation

#### Production Environment Rules
```go
func validateProductionMachine(machine spookytypes.Machine) []string {
    var warnings []string
    
    // Require resource specifications
    if machine.Resources == nil {
        warnings = append(warnings, "production machines should specify resources")
    }
    
    // Require backup schedule
    if machine.MachineMetadata != nil && machine.MachineMetadata.BackupSchedule == "" {
        warnings = append(warnings, "production machines should specify backup schedule")
    }
    
    // Require cost center
    if machine.MachineMetadata != nil && machine.MachineMetadata.CostCenter == "" {
        warnings = append(warnings, "production machines should specify cost center")
    }
    
    return warnings
}
```

#### Development Environment Rules
```go
func validateDevelopmentMachine(machine spookytypes.Machine) []string {
    var warnings []string
    
    // More lenient validation for development
    if machine.MachineMetadata != nil && machine.MachineMetadata.Owner == "" {
        warnings = append(warnings, "development machines should specify owner")
    }
    
    return warnings
}
```

### Cross-File Validation

#### Duplicate Detection
```go
func validateDuplicateHostnames(machines []spookytypes.Machine) []string {
    var errors []string
    hostnameMap := make(map[string][]string)
    
    for _, machine := range machines {
        sourceFile := getSourceFile(machine)
        hostnameMap[machine.Hostname] = append(hostnameMap[machine.Hostname], sourceFile)
    }
    
    for hostname, sources := range hostnameMap {
        if len(sources) > 1 {
            errors = append(errors, fmt.Sprintf("duplicate hostname '%s' found in multiple files: %v", hostname, sources))
        }
    }
    
    return errors
}
```

#### Consistency Validation
```go
func validateEnvironmentConsistency(machines []spookytypes.Machine) []string {
    var warnings []string
    
    // Group machines by environment
    envGroups := make(map[string][]spookytypes.Machine)
    for _, machine := range machines {
        env := getEnvironment(machine)
        envGroups[env] = append(envGroups[env], machine)
    }
    
    // Check consistency within each environment
    for env, envMachines := range envGroups {
        if err := validateEnvironmentGroup(env, envMachines); err != nil {
            warnings = append(warnings, fmt.Sprintf("environment '%s': %v", env, err))
        }
    }
    
    return warnings
}
```

## CLI Integration

### Command Structure

The machines CLI follows the `spooky machines <verb>` pattern:

```go
var machinesCmd = &cobra.Command{
    Use:   "machines",
    Short: "Manage machine inventory",
    Long:  `Manage machine inventory including listing, validation, and connectivity testing.`,
}
```

### Available Commands

#### List Command
```go
var machinesListCmd = &cobra.Command{
    Use:   "list [project-path]",
    Short: "List machines in a project",
    Args:  cobra.ExactArgs(1),
    RunE:  func(cmd *cobra.Command, args []string) error {
        return handleMachinesList(args[0])
    },
}
```

#### Validate Command
```go
var machinesValidateCmd = &cobra.Command{
    Use:   "validate [project-path]",
    Short: "Validate machine inventory",
    Args:  cobra.ExactArgs(1),
    RunE:  func(cmd *cobra.Command, args []string) error {
        return handleMachinesValidate(args[0])
    },
}
```

#### Ping Command
```go
var machinesPingCmd = &cobra.Command{
    Use:   "ping [project-path]",
    Short: "Ping machines to test connectivity",
    Args:  cobra.ExactArgs(1),
    RunE:  func(cmd *cobra.Command, args []string) error {
        return handleMachinesPing(cmd, args[0])
    },
}
```

### Command Handlers

#### List Handler
```go
func handleMachinesList(projectPath string) error {
    ctx := context.Background()
    
    // Initialize dependencies
    if machinesManager == nil {
        if err := InitializeMachinesDependencies(); err != nil {
            return fmt.Errorf("failed to initialize machines dependencies: %w", err)
        }
    }
    
    // Load machines
    machines, err := machinesManager.LoadMachines(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load machines: %w", err)
    }
    
    // Display machines grouped by source
    displayMachinesBySource(machines)
    
    return nil
}
```

#### Validate Handler
```go
func handleMachinesValidate(projectPath string) error {
    ctx := context.Background()
    
    // Initialize dependencies
    if machinesManager == nil {
        if err := InitializeMachinesDependencies(); err != nil {
            return fmt.Errorf("failed to initialize machines dependencies: %w", err)
        }
    }
    
    // Load machines
    machines, err := machinesManager.LoadMachines(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load machines: %w", err)
    }
    
    // Validate machines
    result, err := machinesManager.ValidateMachines(ctx, machines)
    if err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    // Display validation results
    displayValidationResults(result)
    
    return nil
}
```

#### Ping Handler
```go
func handleMachinesPing(cmd *cobra.Command, projectPath string) error {
    ctx := context.Background()
    
    // Get command flags
    format, _ := cmd.Flags().GetString("format")
    verbose, _ := cmd.Flags().GetBool("verbose")
    
    // Initialize dependencies
    if machinesManager == nil {
        if err := InitializeMachinesDependencies(); err != nil {
            return fmt.Errorf("failed to initialize machines dependencies: %w", err)
        }
    }
    
    // Load machines
    machines, err := machinesManager.LoadMachines(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load machines: %w", err)
    }
    
    // Ping machines
    statuses, err := machinesManager.PingMachines(ctx, machines)
    if err != nil {
        return fmt.Errorf("ping failed: %w", err)
    }
    
    // Display results
    if format == "json" {
        return outputPingResultsJSON(statuses, verbose)
    } else {
        outputPingResultsText(statuses, verbose)
    }
    
    return nil
}
```

### Output Formats

#### Text Output
```go
func outputPingResultsText(statuses []spookytypes.MachineStatus, verbose bool) {
    fmt.Printf("🔍 Pinging machines in project\n")
    fmt.Printf("📊 Ping Results: Total machines: %d\n\n", len(statuses))
    
    for _, status := range statuses {
        if status.Status == "online" && !verbose {
            // Smart mode: minimal output for online machines
            fmt.Printf("✅ %s: online (%dms)\n", status.Machine.Hostname, status.Latency)
        } else {
            // Verbose mode or problematic machines
            fmt.Printf("%s %s (%s): %s", getStatusIcon(status.Status), status.Machine.Hostname, status.Machine.Host, status.Status)
            if status.Latency > 0 {
                fmt.Printf(" (%dms)", status.Latency)
            }
            if status.Error != "" {
                fmt.Printf(" (%s)", status.Error)
            }
            fmt.Printf("\n")
        }
    }
}
```

#### JSON Output
```go
func outputPingResultsJSON(statuses []spookytypes.MachineStatus, verbose bool) error {
    for _, status := range statuses {
        machineJSON := map[string]interface{}{
            "hostname": status.Machine.Hostname,
            "status":   status.Status,
        }
        
        // Add details only for problematic machines or verbose mode
        if status.Status != "online" || verbose {
            machineJSON["latency_ms"] = status.Latency
            if status.Status != "online" {
                machineJSON["error"] = getErrorMessage(status)
            }
        }
        
        jsonData, err := json.Marshal(machineJSON)
        if err != nil {
            return fmt.Errorf("failed to marshal machine JSON: %w", err)
        }
        
        fmt.Println(string(jsonData))
    }
    return nil
}
```

## Examples

### Basic Usage

**Loading Machines:**
```go
package main

import (
    "context"
    "fmt"
    "log"
    
    spookyinterfaces "spooky/internal/interfaces"
    spookymachines "spooky/internal/machines"
    spookylogging "spooky/internal/logging"
)

func main() {
    // Initialize dependencies
    logManager := spookylogging.NewLogManager()
    logger := logManager.GetLogger("machines")
    
    validator := spookymachines.NewValidator(logger)
    loader := spookymachines.NewLoader(logger)
    manager := spookymachines.NewManager(logger, loader, validator)
    
    // Load machines
    ctx := context.Background()
    machines, err := manager.LoadMachines(ctx, "./my-project")
    if err != nil {
        log.Fatalf("Failed to load machines: %v", err)
    }
    
    fmt.Printf("Loaded %d machines\n", len(machines))
    
    // Validate machines
    result, err := manager.ValidateMachines(ctx, machines)
    if err != nil {
        log.Fatalf("Validation failed: %v", err)
    }
    
    if len(result.Errors) > 0 {
        fmt.Printf("Validation errors: %v\n", result.Errors)
    }
    
    // Ping machines
    statuses, err := manager.PingMachines(ctx, machines)
    if err != nil {
        log.Fatalf("Ping failed: %v", err)
    }
    
    for _, status := range statuses {
        fmt.Printf("%s: %s\n", status.Machine.Hostname, status.Status)
    }
}
```

### Custom Validation

**Adding Custom Validation Rules:**
```go
func customValidateMachine(machine spookytypes.Machine) []string {
    var warnings []string
    
    // Custom validation: require tags for production machines
    if machine.MachineMetadata != nil && machine.MachineMetadata.Environment == "production" {
        if len(machine.Tags) == 0 {
            warnings = append(warnings, "production machines should have tags")
        }
    }
    
    // Custom validation: require resource specifications for high-performance machines
    if len(machine.Roles) > 0 {
        for _, role := range machine.Roles {
            if role == "high-performance" && machine.Resources == nil {
                warnings = append(warnings, "high-performance machines should specify resources")
                break
            }
        }
    }
    
    return warnings
}
```

### Error Handling

**Comprehensive Error Handling:**
```go
func handleMachineOperations(projectPath string) error {
    ctx := context.Background()
    
    // Initialize manager
    manager := createMachineManager()
    
    // Load machines with detailed error handling
    machines, err := manager.LoadMachines(ctx, projectPath)
    if err != nil {
        // Check for specific error types
        if loadErr, ok := err.(*spookytypesmachines.MachineLoadError); ok {
            return fmt.Errorf("failed to load machines from %s: %w", loadErr.FilePath, err)
        }
        return fmt.Errorf("failed to load machines: %w", err)
    }
    
    // Validate machines with detailed error reporting
    result, err := manager.ValidateMachines(ctx, machines)
    if err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    // Report validation issues
    if len(result.Errors) > 0 {
        fmt.Printf("Validation errors:\n")
        for _, err := range result.Errors {
            if valErr, ok := err.(*spookytypesmachines.MachineValidationError); ok {
                fmt.Printf("  - %s: %s (field: %s, value: %s)\n", 
                    valErr.Field, valErr.Message, valErr.Field, valErr.Value)
            } else {
                fmt.Printf("  - %s\n", err.Error())
            }
        }
        return fmt.Errorf("machine validation failed with %d errors", len(result.Errors))
    }
    
    if len(result.Warnings) > 0 {
        fmt.Printf("Validation warnings:\n")
        for _, warning := range result.Warnings {
            fmt.Printf("  - %s\n", warning)
        }
    }
    
    return nil
}
```

This comprehensive API reference provides all the details needed to work with the spooky machines inventory system, from basic usage to advanced customization and error handling.
