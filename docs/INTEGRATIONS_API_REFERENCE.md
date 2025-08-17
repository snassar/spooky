# Integrations System API Reference

## Overview

This document provides a comprehensive API reference for the spooky integrations system. It covers all interfaces, types, methods, and implementation details for developers working with the integrations system.

**Status: Partially Implemented** - The integrations system has basic functionality but the interface definitions and implementation details have known issues that need to be addressed.

> **See also**: [Known Issues](KNOWN_ISSUES.md#integrations-system-issues) - Comprehensive documentation of all known issues and workarounds

## Core Interfaces

### IntegrationManager Interface

The `IntegrationManager` interface provides the primary entry point for all system integrations:

```go
type IntegrationManager interface {
    // GetFactsIntegration returns the facts integration
    GetFactsIntegration() FactsIntegration
    
    // GetActionsIntegration returns the actions integration
    GetActionsIntegration() ActionsIntegration
    
    // GetVariablesIntegration returns the variables integration
    GetVariablesIntegration() VariablesIntegration
    
    // GetTemplatesIntegration returns the templates integration
    GetTemplatesIntegration() TemplatesIntegration
    
    // GetMachinesIntegration returns the machines integration
    GetMachinesIntegration() MachinesIntegration
    
    // GetSecretsIntegration returns the secrets integration
    GetSecretsIntegration() SecretsIntegration
    
    // GetConfigIntegration returns the configuration integration
    GetConfigIntegration() ConfigIntegration
    
    // Initialize initializes all integrations
    Initialize(ctx context.Context) error
    
    // Shutdown shuts down all integrations
    Shutdown(ctx context.Context) error
    
    // HealthCheck performs health checks on all integrations
    HealthCheck(ctx context.Context) (map[string]HealthStatus, error)
}
```

### FactsIntegration Interface

The `FactsIntegration` interface provides facts collection and management:

```go
type FactsIntegration interface {
    // CollectFacts collects facts from the given machine
    CollectFacts(ctx context.Context, machine string) (*spookytypes.FactCollection, error)
    
    // GetFacts gets facts for the given machine
    GetFacts(ctx context.Context, machine string) (*spookytypes.FactCollection, error)
    
    // StoreFacts stores facts for the given machine
    StoreFacts(ctx context.Context, machine string, facts *spookytypes.FactCollection) error
    
    // ListFacts lists all available facts
    ListFacts(ctx context.Context) ([]string, error)
    
    // ValidateFacts validates facts
    ValidateFacts(ctx context.Context, facts *spookytypes.FactCollection) error
    
    // ExportFacts exports facts to the given format
    ExportFacts(ctx context.Context, format string, output string) error
}
```

### ActionsIntegration Interface

The `ActionsIntegration` interface provides action management and execution:

```go
type ActionsIntegration interface {
    // LoadActions loads actions from the given source
    LoadActions(ctx context.Context, source string) ([]*spookytypes.Action, error)
    
    // RunAction runs the given action on the given machine
    RunAction(ctx context.Context, action *spookytypes.Action, machine string) error
    
    // RunActions runs multiple actions
    RunActions(ctx context.Context, actions []*spookytypes.Action, machines []string) error
    
    // ListActions lists all available actions
    ListActions(ctx context.Context) ([]*spookytypes.Action, error)
    
    // ValidateAction validates the given action
    ValidateAction(ctx context.Context, action *spookytypes.Action) error
    
    // PlanAction plans the execution of the given action
    PlanAction(ctx context.Context, action *spookytypes.Action, machine string) (*spookytypes.ActionPlan, error)
}
```

### VariablesIntegration Interface

The `VariablesIntegration` interface provides variable management:

```go
type VariablesIntegration interface {
    // LoadVariables loads variables from the given source
    LoadVariables(ctx context.Context, source string) (map[string]*spookytypes.Variable, error)
    
    // GetVariable gets the given variable
    GetVariable(ctx context.Context, name string) (*spookytypes.Variable, error)
    
    // SetVariable sets the given variable
    SetVariable(ctx context.Context, name string, value *spookytypes.Variable) error
    
    // ListVariables lists all available variables
    ListVariables(ctx context.Context) ([]string, error)
    
    // ResolveVariables resolves variables with the given context
    ResolveVariables(ctx context.Context, variables map[string]*spookytypes.Variable, ctx *spookytypes.VariablesContext) (map[string]interface{}, error)
    
    // ValidateVariable validates the given variable
    ValidateVariable(ctx context.Context, variable *spookytypes.Variable) error
}
```

### TemplatesIntegration Interface

The `TemplatesIntegration` interface provides template management:

```go
type TemplatesIntegration interface {
    // LoadTemplates loads templates from the given source
    LoadTemplates(ctx context.Context, source string) (map[string]*spookytypes.Template, error)
    
    // GetTemplate gets the given template
    GetTemplate(ctx context.Context, name string) (*spookytypes.Template, error)
    
    // RenderTemplate renders the given template with the given data
    RenderTemplate(ctx context.Context, template *spookytypes.Template, data map[string]interface{}) (string, error)
    
    // ListTemplates lists all available templates
    ListTemplates(ctx context.Context) ([]string, error)
    
    // ValidateTemplate validates the given template
    ValidateTemplate(ctx context.Context, template *spookytypes.Template) error
}
```

### MachinesIntegration Interface

The `MachinesIntegration` interface provides machine management:

```go
type MachinesIntegration interface {
    // LoadMachines loads machines from the given source
    LoadMachines(ctx context.Context, source string) ([]*spookytypes.Machine, error)
    
    // GetMachine gets the given machine
    GetMachine(ctx context.Context, hostname string) (*spookytypes.Machine, error)
    
    // ListMachines lists all available machines
    ListMachines(ctx context.Context) ([]*spookytypes.Machine, error)
    
    // PingMachine pings the given machine
    PingMachine(ctx context.Context, hostname string) error
    
    // ValidateMachine validates the given machine
    ValidateMachine(ctx context.Context, machine *spookytypes.Machine) error
    
    // ExportMachines exports machines to the given format
    ExportMachines(ctx context.Context, format string, output string) error
}
```

### SecretsIntegration Interface

The `SecretsIntegration` interface provides secrets management:

```go
type SecretsIntegration interface {
    // EncryptWithAge encrypts data with age
    EncryptWithAge(ctx context.Context, data []byte, recipients []string) ([]byte, error)
    
    // DecryptWithAge decrypts data with age
    DecryptWithAge(ctx context.Context, data []byte, passphrase string) ([]byte, error)
    
    // ProcessHCLWithSecrets processes HCL with secrets
    ProcessHCLWithSecrets(ctx context.Context, hclData []byte, secrets map[string]string) ([]byte, error)
    
    // ValidateSecrets validates secrets
    ValidateSecrets(ctx context.Context, secrets map[string]string) error
}
```

### ConfigIntegration Interface

The `ConfigIntegration` interface provides configuration management:

```go
type ConfigIntegration interface {
    // LoadConfig loads configuration from the given source
    LoadConfig(ctx context.Context, source string) (*spookytypes.Config, error)
    
    // GetConfig gets the current configuration
    GetConfig(ctx context.Context) (*spookytypes.Config, error)
    
    // ValidateConfig validates the given configuration
    ValidateConfig(ctx context.Context, config *spookytypes.Config) error
    
    // SaveConfig saves the given configuration
    SaveConfig(ctx context.Context, config *spookytypes.Config) error
}
```

## Core Types

### HealthStatus

```go
type HealthStatus struct {
    Status    string                 `json:"status"`
    Message   string                 `json:"message"`
    Timestamp time.Time              `json:"timestamp"`
    Details   map[string]interface{} `json:"details,omitempty"`
}

const (
    HealthStatusHealthy   = "healthy"
    HealthStatusDegraded  = "degraded"
    HealthStatusUnhealthy = "unhealthy"
)
```

### IntegrationContext

```go
type IntegrationContext struct {
    ProjectPath string                 `json:"project_path"`
    Config      *spookytypes.Config    `json:"config"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
    Timestamp   time.Time              `json:"timestamp"`
}
```

### IntegrationResult

```go
type IntegrationResult struct {
    Success    bool                   `json:"success"`
    Message    string                 `json:"message"`
    Data       interface{}            `json:"data,omitempty"`
    Error      string                 `json:"error,omitempty"`
    Timestamp  time.Time              `json:"timestamp"`
    Duration   time.Duration          `json:"duration"`
    Metadata   map[string]interface{} `json:"metadata,omitempty"`
}
```

## Implementation Details

### Current Implementation Status

The integrations system currently has:

1. **Basic Integration Manager**: Simple integration coordination
2. **Integration Factory**: Factory pattern for creating integrations
3. **Health Checking**: Basic health check functionality
4. **Context Management**: Integration context handling

### Missing Features

1. **Complete Integration Implementations**: Many integrations are not fully implemented
2. **Error Handling**: Limited error handling and recovery
3. **Performance Optimization**: No performance optimization
4. **Monitoring**: Limited monitoring and metrics
5. **Configuration Management**: Limited configuration management

## Configuration

### Integration Configuration Structure

```hcl
# Integration configuration
integrations {
  # Facts integration configuration
  facts {
    enabled = true
    storage_path = "./facts.db"
    collection_timeout = "30s"
  }
  
  # Actions integration configuration
  actions {
    enabled = true
    parallel_workers = 4
    execution_timeout = "5m"
  }
  
  # Variables integration configuration
  variables {
    enabled = true
    encryption_enabled = true
    storage_path = "./variables.db"
  }
  
  # Templates integration configuration
  templates {
    enabled = true
    cache_enabled = true
    cache_size = 100
  }
  
  # Machines integration configuration
  machines {
    enabled = true
    ping_timeout = "10s"
    connection_timeout = "30s"
  }
  
  # Secrets integration configuration
  secrets {
    enabled = true
    age_key_path = "~/.config/spooky/age.key"
    encryption_enabled = true
  }
  
  # Configuration integration
  config {
    enabled = true
    auto_reload = true
    reload_interval = "30s"
  }
}
```

## Usage Examples

### Basic Integration Usage

```go
// Create integration manager
manager := NewIntegrationManager()

// Initialize integrations
if err := manager.Initialize(ctx); err != nil {
    return fmt.Errorf("failed to initialize integrations: %w", err)
}
defer manager.Shutdown(ctx)

// Get facts integration
factsIntegration := manager.GetFactsIntegration()

// Collect facts
facts, err := factsIntegration.CollectFacts(ctx, "web-server")
if err != nil {
    return fmt.Errorf("failed to collect facts: %w", err)
}
```

### Health Checking

```go
// Perform health checks
healthStatus, err := manager.HealthCheck(ctx)
if err != nil {
    return fmt.Errorf("health check failed: %w", err)
}

// Check individual integration health
for integration, status := range healthStatus {
    switch status.Status {
    case HealthStatusHealthy:
        log.Printf("Integration %s is healthy", integration)
    case HealthStatusDegraded:
        log.Printf("Integration %s is degraded: %s", integration, status.Message)
    case HealthStatusUnhealthy:
        log.Printf("Integration %s is unhealthy: %s", integration, status.Message)
    }
}
```

### Error Handling

```go
// Handle integration errors
result, err := factsIntegration.CollectFacts(ctx, "web-server")
if err != nil {
    // Check if it's a recoverable error
    if isRecoverableError(err) {
        // Retry with exponential backoff
        return retryWithBackoff(func() error {
            return factsIntegration.CollectFacts(ctx, "web-server")
        })
    }
    return fmt.Errorf("unrecoverable error collecting facts: %w", err)
}
```

## Performance Considerations

### Current Limitations

1. **Synchronous Operations**: Most integration operations are synchronous
2. **No Caching**: Limited caching support
3. **No Connection Pooling**: No connection pooling for external services
4. **Memory Usage**: No memory usage optimization

### Recommended Improvements

1. **Async Operations**: Implement asynchronous integration operations
2. **Caching**: Add caching for frequently accessed data
3. **Connection Pooling**: Implement connection pooling for external services
4. **Performance Monitoring**: Add performance metrics and monitoring

## Integration with Other Systems

### SSH Integration

```go
// SSH operations through machines integration
machinesIntegration := manager.GetMachinesIntegration()

// Ping machine
if err := machinesIntegration.PingMachine(ctx, "web-server"); err != nil {
    log.Printf("Machine web-server is not reachable: %v", err)
}
```

### Facts Integration

```go
// Facts collection through facts integration
factsIntegration := manager.GetFactsIntegration()

// Collect facts
facts, err := factsIntegration.CollectFacts(ctx, "web-server")
if err != nil {
    return fmt.Errorf("failed to collect facts: %w", err)
}

// Store facts
if err := factsIntegration.StoreFacts(ctx, "web-server", facts); err != nil {
    return fmt.Errorf("failed to store facts: %w", err)
}
```

### Actions Integration

```go
// Action execution through actions integration
actionsIntegration := manager.GetActionsIntegration()

// Load actions
actions, err := actionsIntegration.LoadActions(ctx, "./actions.hcl")
if err != nil {
    return fmt.Errorf("failed to load actions: %w", err)
}

// Run action
if err := actionsIntegration.RunAction(ctx, actions[0], "web-server"); err != nil {
    return fmt.Errorf("failed to run action: %w", err)
}
```

## Error Handling

### Integration Errors

```go
// Handle integration errors gracefully
if err := integration.Operation(ctx); err != nil {
    // Log error with context
    log.Printf("Integration operation failed: %v", err)
    
    // Check if error is recoverable
    if isRecoverableError(err) {
        // Implement retry logic
        return retryOperation(integration, ctx)
    }
    
    // Return error for non-recoverable errors
    return fmt.Errorf("integration operation failed: %w", err)
}
```

### Health Check Errors

```go
// Handle health check errors
healthStatus, err := manager.HealthCheck(ctx)
if err != nil {
    log.Printf("Health check failed: %v", err)
    
    // Continue with degraded functionality
    return nil
}

// Handle unhealthy integrations
for integration, status := range healthStatus {
    if status.Status == HealthStatusUnhealthy {
        log.Printf("Integration %s is unhealthy: %s", integration, status.Message)
        
        // Implement fallback logic
        if err := handleUnhealthyIntegration(integration); err != nil {
            log.Printf("Failed to handle unhealthy integration %s: %v", integration, err)
        }
    }
}
```

## Testing

### Integration Testing

```go
func TestIntegrationManager(t *testing.T) {
    // Create test integration manager
    manager := NewTestIntegrationManager()
    
    // Initialize integrations
    if err := manager.Initialize(ctx); err != nil {
        t.Fatalf("Failed to initialize integrations: %v", err)
    }
    
    // Test health checks
    healthStatus, err := manager.HealthCheck(ctx)
    if err != nil {
        t.Fatalf("Health check failed: %v", err)
    }
    
    // Verify all integrations are healthy
    for integration, status := range healthStatus {
        if status.Status != HealthStatusHealthy {
            t.Errorf("Integration %s is not healthy: %s", integration, status.Status)
        }
    }
}
```

### Mock Integration

```go
type MockFactsIntegration struct {
    facts map[string]*spookytypes.FactCollection
    mutex sync.RWMutex
}

func (m *MockFactsIntegration) CollectFacts(ctx context.Context, machine string) (*spookytypes.FactCollection, error) {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    
    if facts, exists := m.facts[machine]; exists {
        return facts, nil
    }
    
    return nil, fmt.Errorf("facts not found for machine: %s", machine)
}

func (m *MockFactsIntegration) StoreFacts(ctx context.Context, machine string, facts *spookytypes.FactCollection) error {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    m.facts[machine] = facts
    return nil
}
```

## Best Practices

### Integration Management

1. **Initialization**: Always initialize integrations before use
2. **Shutdown**: Always shutdown integrations when done
3. **Health Checks**: Regular health checks for all integrations
4. **Error Handling**: Proper error handling and recovery
5. **Monitoring**: Monitor integration performance and health

### Performance Optimization

```go
// Use connection pooling for external services
type ConnectionPool struct {
    connections map[string]interface{}
    mutex       sync.RWMutex
    maxConnections int
}

// Implement caching for frequently accessed data
type Cache struct {
    data    map[string]interface{}
    mutex   sync.RWMutex
    ttl     time.Duration
}

// Use async operations for long-running tasks
func (i *Integration) AsyncOperation(ctx context.Context) <-chan error {
    result := make(chan error, 1)
    
    go func() {
        defer close(result)
        result <- i.performOperation(ctx)
    }()
    
    return result
}
```

### Error Recovery

```go
// Implement retry logic with exponential backoff
func retryWithBackoff(operation func() error) error {
    maxRetries := 3
    backoff := time.Second
    
    for attempt := 0; attempt < maxRetries; attempt++ {
        if err := operation(); err == nil {
            return nil
        }
        
        if attempt < maxRetries-1 {
            time.Sleep(backoff)
            backoff *= 2
        }
    }
    
    return fmt.Errorf("operation failed after %d attempts", maxRetries)
}
```

## Future Enhancements

### Planned Features

1. **Plugin System**: Pluggable integration system
2. **Service Discovery**: Automatic service discovery
3. **Load Balancing**: Load balancing for integrations
4. **Circuit Breaker**: Circuit breaker pattern implementation
5. **Metrics Integration**: Integration with metrics systems
6. **Distributed Tracing**: Integration with distributed tracing systems

### Architecture Improvements

1. **Microservices**: Move to microservices architecture
2. **Event-Driven**: Event-driven integration patterns
3. **API Gateway**: API gateway for integrations
4. **Service Mesh**: Service mesh integration
5. **Containerization**: Containerized integration deployment

## Related Documentation

- [Integrations User Guide](INTEGRATIONS_USER_GUIDE.md) - User guide for integration features
- [Integrations Troubleshooting](INTEGRATIONS_TROUBLESHOOTING.md) - Troubleshooting guide
- [Interface Architecture](INTERFACE_ARCHITECTURE.md) - Interface-based architecture
- [System Design](../design/systems/) - System design documentation
