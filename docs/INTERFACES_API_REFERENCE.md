# Interfaces System API Reference

## Overview

This document provides a comprehensive API reference for the spooky interfaces system. It covers all interfaces, types, methods, and implementation details for developers working with the interfaces system.

**Status: Partially Implemented** - The interfaces system has basic functionality but the interface definitions and implementation details have known issues that need to be addressed.

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

### ValidationResult

```go
type ValidationResult struct {
    Valid    bool                   `json:"valid"`
    Errors   []ValidationError      `json:"errors,omitempty"`
    Warnings []ValidationWarning    `json:"warnings,omitempty"`
    Details  map[string]interface{} `json:"details,omitempty"`
}

type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
    Value   string `json:"value,omitempty"`
}

type ValidationWarning struct {
    Field   string `json:"field"`
    Message string `json:"message"`
    Value   string `json:"value,omitempty"`
}
```

### Context Interfaces

```go
// BaseContext provides base context for all integrations
type BaseContext interface {
    GetProjectPath() string
    GetTimestamp() time.Time
    GetMetadata() map[string]interface{}
}

// ProjectContext contains all project-related contexts
type ProjectContext interface {
    BaseContext
    GetFactsContext() FactsContext
    GetVariablesContext() VariablesContext
    GetTemplatesContext() TemplatesContext
    GetMachinesContext() MachinesContext
    GetActionsContext() ActionsContext
}

// FactsContext provides facts data for integrations
type FactsContext interface {
    BaseContext
    GetFacts() map[string]*spookytypes.FactCollection
    GetFactCollection(machine string) (*spookytypes.FactCollection, error)
}

// VariablesContext provides variables data for integrations
type VariablesContext interface {
    BaseContext
    GetVariables() map[string]*spookytypes.Variable
    GetVariable(name string) (*spookytypes.Variable, error)
}

// TemplatesContext provides templates data for integrations
type TemplatesContext interface {
    BaseContext
    GetTemplates() map[string]*spookytypes.Template
    GetTemplate(name string) (*spookytypes.Template, error)
}

// MachinesContext provides machines data for integrations
type MachinesContext interface {
    BaseContext
    GetMachines() []*spookytypes.Machine
    GetMachine(hostname string) (*spookytypes.Machine, error)
}

// ActionsContext provides actions data for integrations
type ActionsContext interface {
    BaseContext
    GetActions() []*spookytypes.Action
    GetAction(name string) (*spookytypes.Action, error)
}

// ActionRunContext provides acting context for actions
type ActionRunContext interface {
    BaseContext
    GetAction() *spookytypes.Action
    GetMachine() *spookytypes.Machine
    GetVariables() map[string]interface{}
    GetFacts() map[string]interface{}
}
```

## Implementation Details

### Current Implementation Status

The interfaces system currently has:

1. **Basic Interface Definitions**: Core interface definitions are in place
2. **Context Interfaces**: Context interfaces for state management
3. **Validation Interfaces**: Basic validation result structures
4. **Health Status**: Health status definitions

### Missing Features

1. **Complete Interface Implementations**: Many interfaces are not fully implemented
2. **Error Handling**: Limited error handling and recovery
3. **Performance Optimization**: No performance optimization
4. **Monitoring**: Limited monitoring and metrics
5. **Configuration Management**: Limited configuration management

## Usage Examples

### Basic Interface Usage

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

### Context Usage

```go
// Create project context
projectCtx := NewProjectContext("./my-project")

// Get facts context
factsCtx := projectCtx.GetFactsContext()

// Use facts context
facts, err := factsCtx.GetFactCollection("web-server")
if err != nil {
    return fmt.Errorf("failed to get facts: %w", err)
}
```

### Validation Usage

```go
// Validate facts
result, err := factsIntegration.ValidateFacts(ctx, facts)
if err != nil {
    return fmt.Errorf("validation failed: %w", err)
}

if !result.Valid {
    for _, error := range result.Errors {
        log.Printf("Validation error: %s - %s", error.Field, error.Message)
    }
    return fmt.Errorf("facts validation failed")
}
```

## Error Handling

### Interface Errors

```go
// Handle interface errors gracefully
if err := integration.Operation(ctx); err != nil {
    // Check if error is recoverable
    if isRecoverableError(err) {
        // Retry with exponential backoff
        return retryWithBackoff(func() error {
            return integration.Operation(ctx)
        })
    }
    return fmt.Errorf("unrecoverable error: %w", err)
}
```

### Validation Errors

```go
// Handle validation errors
result, err := validator.Validate(data)
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
    
    return fmt.Errorf("validation failed with %d errors", len(result.Errors))
}
```

## Testing

### Interface Testing

```go
func TestFactsIntegration(t *testing.T) {
    // Create mock facts integration
    mockIntegration := &MockFactsIntegration{
        facts: make(map[string]*spookytypes.FactCollection),
    }
    
    // Add test data
    testFacts := &spookytypes.FactCollection{
        Machine: "test-server",
        Facts:   map[string]interface{}{"os": "linux"},
    }
    mockIntegration.facts["test-server"] = testFacts
    
    // Test facts collection
    facts, err := mockIntegration.CollectFacts(ctx, "test-server")
    if err != nil {
        t.Fatalf("Failed to collect facts: %v", err)
    }
    
    if facts.Machine != "test-server" {
        t.Errorf("Expected machine 'test-server', got '%s'", facts.Machine)
    }
}
```

### Mock Implementation

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

// Implement other interface methods...
func (m *MockFactsIntegration) GetFacts(ctx context.Context, machine string) (*spookytypes.FactCollection, error) {
    return m.CollectFacts(ctx, machine)
}

func (m *MockFactsIntegration) ListFacts(ctx context.Context) ([]string, error) {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    
    machines := make([]string, 0, len(m.facts))
    for machine := range m.facts {
        machines = append(machines, machine)
    }
    return machines, nil
}

func (m *MockFactsIntegration) ValidateFacts(ctx context.Context, facts *spookytypes.FactCollection) error {
    if facts == nil {
        return fmt.Errorf("facts cannot be nil")
    }
    if facts.Machine == "" {
        return fmt.Errorf("machine name is required")
    }
    return nil
}

func (m *MockFactsIntegration) ExportFacts(ctx context.Context, format string, output string) error {
    // Mock implementation
    return nil
}
```

## Best Practices

### Interface Design

1. **Single Responsibility**: Each interface should have a single responsibility
2. **Consistent Naming**: Use consistent naming conventions across interfaces
3. **Error Handling**: Define clear error handling patterns
4. **Context Usage**: Use context for state management
5. **Validation**: Include validation in interface contracts

### Implementation Guidelines

```go
// Use dependency injection
type Manager struct {
    factsIntegration     FactsIntegration
    actionsIntegration   ActionsIntegration
    variablesIntegration VariablesIntegration
    // ... other integrations
}

// Implement proper error handling
func (m *Manager) Operation(ctx context.Context) error {
    if err := m.validateContext(ctx); err != nil {
        return fmt.Errorf("invalid context: %w", err)
    }
    
    if err := m.performOperation(ctx); err != nil {
        return fmt.Errorf("operation failed: %w", err)
    }
    
    return nil
}

// Use context for state management
func (m *Manager) performOperation(ctx context.Context) error {
    // Extract context information
    projectPath := ctx.Value("project_path").(string)
    
    // Use context for operation
    return m.factsIntegration.CollectFacts(ctx, "machine")
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

1. **Plugin System**: Pluggable interface system
2. **Service Discovery**: Automatic service discovery
3. **Load Balancing**: Load balancing for interfaces
4. **Circuit Breaker**: Circuit breaker pattern implementation
5. **Metrics Integration**: Integration with metrics systems
6. **Distributed Tracing**: Integration with distributed tracing systems

### Architecture Improvements

1. **Microservices**: Move to microservices architecture
2. **Event-Driven**: Event-driven interface patterns
3. **API Gateway**: API gateway for interfaces
4. **Service Mesh**: Service mesh integration
5. **Containerization**: Containerized interface deployment

## Related Documentation

- [Interface Architecture](INTERFACE_ARCHITECTURE.md) - Interface-based architecture
- [Interface Definitions](INTERFACE_DEFINITIONS.md) - Complete interface definitions
- [Interface Dos and Don'ts](INTERFACE_DOS_DONTS.md) - Interface best practices
- [System Design](../design/systems/) - System design documentation
