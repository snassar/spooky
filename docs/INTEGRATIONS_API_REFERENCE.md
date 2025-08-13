# Integrations API Reference

## Overview

The IntegrationManager provides centralized coordination for all system integrations in spooky. It serves as the primary interface for accessing and managing facts, actions, variables, templates, machines, secrets, and configuration integrations.

## Core Interfaces

### IntegrationManager

The central coordinator for all system integrations.

```go
type IntegrationManager interface {
    // Get individual integrations
    GetFactsIntegration() FactsIntegration
    GetActionsIntegration() ActionsIntegration
    GetVariablesIntegration() VariablesIntegration
    GetTemplatesIntegration() TemplatesIntegration
    GetMachinesIntegration() MachinesIntegration
    GetSecretsIntegration() SecretsIntegration
    GetConfigIntegration() ConfigIntegration
}
```

### Individual Integration Interfaces

#### FactsIntegration
```go
type FactsIntegration interface {
    CollectFacts(ctx context.Context, server string) (*FactCollection, error)
    StoreFacts(ctx context.Context, facts *FactCollection) error
    GetFacts(ctx context.Context, server string) (*FactCollection, error)
    ValidateFacts(ctx context.Context, facts *FactCollection) (*ValidationResult, error)
}
```

#### ActionsIntegration
```go
type ActionsIntegration interface {
    LoadActions(ctx context.Context, path string) ([]Action, error)
    ValidateActions(ctx context.Context, actions []Action) (*ValidationResult, error)
    RunActions(ctx context.Context, actions []Action, ctx ActionRunContext) error
    PlanActions(ctx context.Context, actions []Action) (*ActionPlan, error)
}
```

#### VariablesIntegration
```go
type VariablesIntegration interface {
    LoadVariables(ctx context.Context, path string) ([]Variable, error)
    ResolveVariables(ctx context.Context, variables []Variable, ctx VariableContext) (*VariableResolutionResult, error)
    ValidateVariables(ctx context.Context, variables []Variable) (*ValidationResult, error)
}
```

#### TemplatesIntegration
```go
type TemplatesIntegration interface {
    LoadTemplate(ctx context.Context, templatePath string) (*Template, error)
    RenderTemplate(ctx context.Context, template *Template, data map[string]interface{}) (string, error)
    ValidateTemplate(ctx context.Context, template *Template) (*ValidationResult, error)
}
```

#### MachinesIntegration
```go
type MachinesIntegration interface {
    LoadMachines(ctx context.Context, path string) ([]Machine, error)
    ValidateMachines(ctx context.Context, machines []Machine) (*ValidationResult, error)
    ConnectToMachine(ctx context.Context, machine *Machine) (SSHConnection, error)
    PingMachine(ctx context.Context, machine *Machine) error
}
```

#### SecretsIntegration
```go
type SecretsIntegration interface {
    Encrypt(ctx context.Context, data []byte, key []byte) ([]byte, error)
    Decrypt(ctx context.Context, data []byte, key []byte) ([]byte, error)
    ValidateKey(ctx context.Context, key []byte) error
}
```

#### ConfigIntegration
```go
type ConfigIntegration interface {
    LoadConfig(ctx context.Context, source string) (*Config, error)
    ValidateConfig(ctx context.Context, config *Config) (*ValidationResult, error)
    SaveConfig(ctx context.Context, config *Config, destination string) error
}
```

## Implementation Details

### Manager Implementation

The `Manager` struct implements the `IntegrationManager` interface:

```go
type Manager struct {
    logger                Logger
    factsIntegration      FactsIntegration
    actionsIntegration    ActionsIntegration
    variablesIntegration  VariablesIntegration
    templatesIntegration  TemplatesIntegration
    machinesIntegration   MachinesIntegration
    secretsIntegration    SecretsIntegration
    configIntegration     ConfigIntegration
    healthStatus          map[string]bool
    healthMu              sync.RWMutex
}
```

### Health Monitoring

The IntegrationManager includes built-in health monitoring:

```go
// Internal methods for health management
func (m *Manager) ValidateSystemHealth(ctx context.Context) (*ValidationResult, error)
func (m *Manager) GetHealthStatus() map[string]bool
func (m *Manager) UpdateHealthStatus(integration string, healthy bool)
func (m *Manager) CoordinatedOperation(ctx context.Context, operation func() error) error
```

### Factory Pattern

Use the `Factory` to create IntegrationManager instances:

```go
type Factory struct {
    logger Logger
}

func (f *Factory) CreateIntegrationManager() IntegrationManager
```

## Usage Examples

### Basic Usage

```go
// Create integration manager
factory := integration.NewFactory(logger)
manager := factory.CreateIntegrationManager()

// Access integrations
facts := manager.GetFactsIntegration()
actions := manager.GetActionsIntegration()
templates := manager.GetTemplatesIntegration()

// Use integrations
factsCollection, err := facts.CollectFacts(ctx, "server.example.com")
if err != nil {
    return err
}
```

### Health Validation

```go
// Validate system health
result, err := manager.ValidateSystemHealth(ctx)
if err != nil {
    return err
}

if !result.Valid {
    for _, error := range result.Errors {
        log.Printf("Integration error: %s", error.Message)
    }
}
```

### Coordinated Operations

```go
// Perform coordinated operations
err := manager.CoordinatedOperation(ctx, func() error {
    // All integration operations within this function
    // will be validated for health before execution
    return nil
})
```

## CLI Integration

The IntegrationManager is accessible via CLI commands:

```bash
# List available integrations and their status
spooky integrations list

# Validate all integrations are working
spooky integrations validate
```

### CLI Output Examples

```bash
$ spooky integrations list
Available Integrations:
=======================
facts        ✅ available
actions      ✅ available
variables    ❌ unavailable
templates    ✅ available
machines     ✅ available
secrets      ✅ available
config       ✅ available

$ spooky integrations validate
✅ All integrations are working correctly
```

## Error Handling

All integration methods return structured errors:

```go
// Validation errors include context and field information
type ValidationError struct {
    Message string
    Field   string
    Type    string
    Rule    string
}

// Validation results aggregate multiple errors
type ValidationResult struct {
    Valid    bool
    Errors   []ValidationError
    Warnings []ValidationError
}
```

## Best Practices

1. **Always validate health** before performing critical operations
2. **Use coordinated operations** for multi-integration workflows
3. **Handle validation errors** with proper context and user feedback
4. **Monitor integration status** regularly in production environments
5. **Use the factory pattern** for consistent IntegrationManager creation

## Integration Status

Integration status is tracked internally:

- **Available**: Integration is properly initialized and functional
- **Unavailable**: Integration is nil or not properly configured
- **Health monitoring**: Real-time status tracking with thread-safe access

## Thread Safety

The IntegrationManager is thread-safe:

- Health status updates use read-write mutex protection
- Integration access is concurrent-safe
- Coordinated operations provide atomic execution guarantees
