# Facts System API Reference

## Overview

This document provides a comprehensive API reference for the spooky facts system. It covers all interfaces, types, methods, and implementation details for developers working with the facts system.

**Status: Implemented** - The facts system provides comprehensive functionality for fact collection, storage, and management.

## Core Interfaces

### FactsIntegration Interface

The `FactsIntegration` interface provides the primary entry point for facts operations:

```go
type FactsIntegration interface {
    // CollectFacts collects facts from the given machine
    CollectFacts(ctx context.Context, machine *spookytypes.Machine) (interface{}, error)
    
    // StoreFacts stores facts in memory
    StoreFacts(ctx context.Context, facts interface{}) error
    
    // LoadFacts loads facts from memory
    LoadFacts(ctx context.Context) (interface{}, error)
    
    // ValidateFacts validates facts
    ValidateFacts(ctx context.Context, facts interface{}) (*spookytypes.ValidationResult, error)
    
    // DecryptFacts decrypts age-encrypted values in facts collection
    DecryptFacts(ctx context.Context, facts interface{}, secretsIntegration SecretsIntegration, identityPath string) error
    
    // GetManager returns the underlying fact manager
    GetManager() interface{}
}
```

**Implementation Status**: ✅ **Implemented** - Complete functionality for fact collection and management

### FactManager Interface

The `FactManager` interface provides fact management and collection:

```go
type FactManager interface {
    // CollectFacts collects facts from the given machine
    CollectFacts(ctx context.Context, machine *spookytypes.Machine) (*spookytypes.FactCollection, error)
    
    // ValidateFacts validates facts
    ValidateFacts(ctx context.Context, facts *spookytypes.FactCollection) (*spookytypes.ValidationResult, error)
    
    // DecryptFacts decrypts age-encrypted values in facts collection
    DecryptFacts(ctx context.Context, facts *spookytypes.FactCollection, secretsIntegration SecretsIntegration, identityPath string) error
}
```

**Implementation Status**: ✅ **Implemented** - Complete functionality for fact management

## Core Types

### FactCollection

```go
type FactCollection struct {
    MachineID string                 `hcl:"machine_id" json:"machine_id"`
    Facts     map[string]interface{} `hcl:"facts" json:"facts"`
    Metadata  *FactMetadata          `hcl:"metadata,block" json:"metadata,omitempty"`
    CreatedAt time.Time              `hcl:"created_at" json:"created_at"`
    UpdatedAt time.Time              `hcl:"updated_at" json:"updated_at"`
}

type FactMetadata struct {
    Collector string            `hcl:"collector" json:"collector"`
    Version   string            `hcl:"version" json:"version"`
    Tags      map[string]string `hcl:"tags" json:"tags,omitempty"`
}
```

### Fact Types

```go
// Basic fact types
type SystemFacts struct {
    Hostname    string `hcl:"hostname" json:"hostname"`
    OS          string `hcl:"os" json:"os"`
    OSVersion   string `hcl:"os_version" json:"os_version"`
    Architecture string `hcl:"architecture" json:"architecture"`
    Kernel      string `hcl:"kernel" json:"kernel"`
    Uptime      int64  `hcl:"uptime" json:"uptime"`
}

type NetworkFacts struct {
    Hostname    string            `hcl:"hostname" json:"hostname"`
    IPAddresses map[string]string `hcl:"ip_addresses" json:"ip_addresses"`
    Interfaces  []string          `hcl:"interfaces" json:"interfaces"`
}

type HardwareFacts struct {
    CPUCount    int    `hcl:"cpu_count" json:"cpu_count"`
    MemoryTotal int64  `hcl:"memory_total" json:"memory_total"`
    DiskSpace   int64  `hcl:"disk_space" json:"disk_space"`
    Model       string `hcl:"model" json:"model"`
}
```

## Current Implementation Status

### ✅ Working Components

1. **Fact Collection**: SSH-based fact collection from remote machines
2. **Fact Storage**: In-memory fact storage during operations
3. **Fact Validation**: Comprehensive validation of fact structures
4. **CLI Integration**: `spooky facts gather` command with filtering options
5. **Project Integration**: Facts loading from project configuration
6. **SSH Integration**: SSH-based fact collection with authentication
7. **Age Encryption**: Support for age-encrypted fact values
8. **Filtering Support**: Support for machine and tag filtering
9. **Export Support**: Facts export to JSON format
10. **Error Handling**: Comprehensive error handling and validation

### 🔧 Key Features

1. **SSH-Based Collection**: SSH-based fact collection from remote machines
2. **Age Encryption**: Support for age-encrypted sensitive fact values
3. **Validation**: Comprehensive fact validation with detailed error reporting
4. **Memory Storage**: In-memory fact storage for efficient operations
5. **CLI Integration**: Full CLI integration with filtering and export options

## Implementation Details

### Fact Collection System

The facts system uses SSH to collect facts from remote machines:

```go
// Collect facts from a remote machine
func (m *Manager) CollectFacts(ctx context.Context, machine *spookytypes.Machine) (*spookytypes.FactCollection, error) {
    // Create SSH connection
    connection, err := m.sshManager.Connect(ctx, &spookytypes.ConnectionRequest{
        Host:     machine.Host,
        Port:     machine.Port,
        User:     machine.User,
        KeyPath:  machine.KeyFile,
        Password: machine.Password,
        Timeout:  30 * time.Second,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to connect to %s: %w", machine.Hostname, err)
    }
    
    // Collect system facts
    facts := &spookytypes.FactCollection{
        MachineID: machine.Hostname,
        Facts:     make(map[string]interface{}),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    // Collect basic system information
    if err := m.collectSystemFacts(ctx, connection, facts); err != nil {
        return nil, fmt.Errorf("failed to collect system facts: %w", err)
    }
    
    // Collect network information
    if err := m.collectNetworkFacts(ctx, connection, facts); err != nil {
        return nil, fmt.Errorf("failed to collect network facts: %w", err)
    }
    
    // Collect hardware information
    if err := m.collectHardwareFacts(ctx, connection, facts); err != nil {
        return nil, fmt.Errorf("failed to collect hardware facts: %w", err)
    }
    
    return facts, nil
}
```

### Fact Validation System

```go
// Validate fact collection
func (m *Manager) ValidateFacts(ctx context.Context, facts *spookytypes.FactCollection) (*spookytypes.ValidationResult, error) {
    var errors []spookytypesschemas.SchemaError
    var warnings []spookytypesschemas.SchemaError
    
    // Validate required fields
    if facts.MachineID == "" {
        errors = append(errors, spookytypesschemas.SchemaError{
            Field:   "machine_id",
            Message: "machine_id is required",
        })
    }
    
    if facts.Facts == nil {
        errors = append(errors, spookytypesschemas.SchemaError{
            Field:   "facts",
            Message: "facts cannot be nil",
        })
    }
    
    // Validate fact values
    for key, value := range facts.Facts {
        if value == nil {
            warnings = append(warnings, spookytypesschemas.SchemaError{
                Field:   fmt.Sprintf("facts.%s", key),
                Message: "fact value is nil",
            })
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

### Age Encryption Support

```go
// Decrypt age-encrypted fact values
func (m *Manager) DecryptFacts(ctx context.Context, facts *spookytypes.FactCollection, secretsIntegration SecretsIntegration, identityPath string) error {
    if facts == nil {
        return fmt.Errorf("facts cannot be nil")
    }
    
    if secretsIntegration == nil {
        return fmt.Errorf("secrets integration cannot be nil")
    }
    
    // Decrypt fact values recursively
    return m.decryptFactValues(facts.Facts, secretsIntegration, identityPath)
}

func (m *Manager) decryptFactValues(values map[string]interface{}, secretsIntegration SecretsIntegration, identityPath string) error {
    for key, value := range values {
        switch v := value.(type) {
        case string:
            // Check if value is age-encrypted
            if strings.HasPrefix(v, "age1") || strings.Contains(v, "-----BEGIN AGE ENCRYPTED FILE-----") {
                decrypted, err := secretsIntegration.DecryptWithAge(ctx, []byte(v), identityPath)
                if err != nil {
                    return fmt.Errorf("failed to decrypt fact value %s: %w", key, err)
                }
                values[key] = string(decrypted)
            }
        case map[string]interface{}:
            // Recursively decrypt nested values
            if err := m.decryptFactValues(v, secretsIntegration, identityPath); err != nil {
                return err
            }
        case []interface{}:
            // Decrypt array values
            for i, item := range v {
                if str, ok := item.(string); ok {
                    if strings.HasPrefix(str, "age1") || strings.Contains(str, "-----BEGIN AGE ENCRYPTED FILE-----") {
                        decrypted, err := secretsIntegration.DecryptWithAge(ctx, []byte(str), identityPath)
                        if err != nil {
                            return fmt.Errorf("failed to decrypt fact array value %s[%d]: %w", key, i, err)
                        }
                        v[i] = string(decrypted)
                    }
                }
            }
        }
    }
    
    return nil
}
```

## Usage Examples

### Basic Fact Collection

```go
// Create facts integration
factsIntegration := NewFactsIntegration(manager)

// Collect facts from a machine
machine := &spookytypes.Machine{
    Hostname: "web-server",
    Host:     "192.168.1.100",
    Port:     22,
    User:     "admin",
    KeyFile:  "~/.ssh/id_rsa",
}

facts, err := factsIntegration.CollectFacts(ctx, machine)
if err != nil {
    return fmt.Errorf("failed to collect facts: %w", err)
}

// Validate facts
result, err := factsIntegration.ValidateFacts(ctx, facts)
if err != nil {
    return fmt.Errorf("failed to validate facts: %w", err)
}

if !result.Valid {
    for _, error := range result.Errors {
        log.Printf("Validation error: %s", error.Message)
    }
    return fmt.Errorf("facts validation failed")
}
```

### Fact Collection with Age Decryption

```go
// Collect facts with age decryption
facts, err := factsIntegration.CollectFacts(ctx, machine)
if err != nil {
    return fmt.Errorf("failed to collect facts: %w", err)
}

// Decrypt age-encrypted values
err = factsIntegration.DecryptFacts(ctx, facts, secretsIntegration, "~/.age/identity.txt")
if err != nil {
    return fmt.Errorf("failed to decrypt facts: %w", err)
}

// Access decrypted facts
if hostname, ok := facts.Facts["hostname"].(string); ok {
    log.Printf("Hostname: %s", hostname)
}
```

### CLI Usage

```bash
# Collect facts from all machines in a project
spooky facts gather ./my-project

# Collect facts from specific machines
spooky facts gather ./my-project --machine web-server --machine db-server

# Collect facts from machines with specific tags
spooky facts gather ./my-project --tags environment=production

# Export facts to JSON
spooky facts gather ./my-project --export --output facts.json

# Collect facts with age decryption
spooky facts gather ./my-project --decrypt --identity ~/.age/identity.txt
```

## Error Handling

### Fact Collection Errors

```go
// Handle fact collection errors
facts, err := factsIntegration.CollectFacts(ctx, machine)
if err != nil {
    // Check for specific error types
    if strings.Contains(err.Error(), "connection refused") {
        return fmt.Errorf("machine %s is unreachable: %w", machine.Hostname, err)
    }
    
    if strings.Contains(err.Error(), "authentication failed") {
        return fmt.Errorf("authentication failed for %s: %w", machine.Hostname, err)
    }
    
    return fmt.Errorf("failed to collect facts from %s: %w", machine.Hostname, err)
}
```

### Validation Errors

```go
// Handle validation errors
result, err := factsIntegration.ValidateFacts(ctx, facts)
if err != nil {
    return fmt.Errorf("validation failed: %w", err)
}

if !result.Valid {
    // Log all validation errors
    for _, error := range result.Errors {
        log.Printf("Validation error in %s: %s", error.Field, error.Message)
    }
    
    // Log all validation warnings
    for _, warning := range result.Warnings {
        log.Printf("Validation warning in %s: %s", warning.Field, warning.Message)
    }
    
    return fmt.Errorf("facts validation failed with %d errors", len(result.Errors))
}
```

## Testing

### Fact Collection Testing

```go
func TestFactCollection(t *testing.T) {
    // Create mock SSH manager
    mockSSHManager := &MockSSHManager{}
    
    // Create facts manager
    manager := NewManager(logger, mockSSHManager)
    
    // Test machine
    machine := &spookytypes.Machine{
        Hostname: "test-server",
        Host:     "192.168.1.100",
        Port:     22,
        User:     "test",
    }
    
    // Collect facts
    facts, err := manager.CollectFacts(ctx, machine)
    if err != nil {
        t.Fatalf("Failed to collect facts: %v", err)
    }
    
    // Validate facts
    if facts.MachineID != "test-server" {
        t.Errorf("Expected machine ID 'test-server', got '%s'", facts.MachineID)
    }
    
    if facts.Facts == nil {
        t.Error("Expected facts map, got nil")
    }
}
```

### Mock SSH Manager

```go
type MockSSHManager struct {
    // Mock implementation for testing
}

func (m *MockSSHManager) Connect(ctx context.Context, request *spookytypes.ConnectionRequest) (*spookytypes.ConnectionResult, error) {
    return &spookytypes.ConnectionResult{
        Success: true,
        Connection: &spookytypes.Connection{
            Host: request.Host,
            Port: request.Port,
            User: request.User,
        },
    }, nil
}

func (m *MockSSHManager) RunCommand(ctx context.Context, connection *spookytypes.Connection, command *spookytypes.SSHCommand) (*spookytypes.SSHCommandResult, error) {
    // Return mock command results
    switch command.Command {
    case "hostname":
        return &spookytypes.SSHCommandResult{
            Success: true,
            Output:  "test-server",
            ExitCode: 0,
        }, nil
    case "uname -s":
        return &spookytypes.SSHCommandResult{
            Success: true,
            Output:  "Linux",
            ExitCode: 0,
        }, nil
    default:
        return &spookytypes.SSHCommandResult{
            Success:  false,
            Error:    "command not found",
            ExitCode: 127,
        }, nil
    }
}
```

## Best Practices

### Fact Collection

1. **Use SSH Keys**: Prefer SSH key authentication over passwords
2. **Validate Facts**: Always validate collected facts before use
3. **Handle Errors**: Implement proper error handling for collection failures
4. **Use Timeouts**: Set appropriate timeouts for SSH connections
5. **Encrypt Sensitive Data**: Use age encryption for sensitive fact values

### Performance Optimization

```go
// Collect facts in parallel
func collectFactsParallel(machines []*spookytypes.Machine, factsIntegration FactsIntegration) ([]*spookytypes.FactCollection, error) {
    var wg sync.WaitGroup
    results := make([]*spookytypes.FactCollection, len(machines))
    errors := make([]error, len(machines))
    
    for i, machine := range machines {
        wg.Add(1)
        go func(index int, m *spookytypes.Machine) {
            defer wg.Done()
            
            facts, err := factsIntegration.CollectFacts(ctx, m)
            if err != nil {
                errors[index] = err
                return
            }
            
            results[index] = facts
        }(i, machine)
    }
    
    wg.Wait()
    
    // Check for errors
    for _, err := range errors {
        if err != nil {
            return nil, fmt.Errorf("fact collection failed: %w", err)
        }
    }
    
    return results, nil
}
```

## Future Enhancements

### Planned Features

1. **Fact Caching**: Implement fact caching for improved performance
2. **Fact Templates**: Support for fact collection templates
3. **Custom Collectors**: Pluggable fact collectors
4. **Fact Dependencies**: Support for fact dependencies and relationships
5. **Fact Versioning**: Version control for fact collections
6. **Fact Analytics**: Analytics and reporting for fact data

### Architecture Improvements

1. **Distributed Collection**: Distributed fact collection across multiple collectors
2. **Streaming Collection**: Streaming fact collection for real-time updates
3. **Fact Compression**: Compression for large fact collections
4. **Fact Indexing**: Indexing for efficient fact queries
5. **Fact Replication**: Replication for high availability

## Related Documentation

- [Facts Collection System](FACTS_COLLECTION_SYSTEM.md) - Facts collection architecture
- [Facts User Guide](FACTS_USER_GUIDE.md) - User guide for facts system
- [Facts Troubleshooting](FACTS_TROUBLESHOOTING.md) - Troubleshooting guide
- [SSH API Reference](SSH_API_REFERENCE.md) - SSH system API reference
- [Secrets API Reference](SECRETS_API_REFERENCE.md) - Secrets management API reference
