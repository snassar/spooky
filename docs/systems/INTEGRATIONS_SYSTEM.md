# Integrations System

## Overview

The Integrations System provides comprehensive system coordination and integration capabilities for the spooky codebase. It serves as the central orchestrator that coordinates all system components, manages dependencies, and provides unified access to all system integrations.

**Status**: **Implemented** - Complete integration system with coordination, dependency management, and unified access patterns.

## Related Systems

This system coordinates and integrates all other spooky systems:

- **[Actions System](ACTIONS_SYSTEM.md)** - Provides action orchestration and running through ActionsIntegration
- **[Facts System](FACTS_SYSTEM.md)** - Provides fact collection and storage through FactsIntegration
- **[Variables System](VARIABLES_SYSTEM.md)** - Provides variable resolution and management through VariablesIntegration
- **[Templates System](TEMPLATES_SYSTEM.md)** - Provides template rendering and validation through TemplatesIntegration
- **[Machines System](MACHINES_SYSTEM.md)** - Provides machine inventory and connectivity through MachinesIntegration
- **[SSH System](SSH_SYSTEM.md)** - Provides SSH connectivity and acting through SSH integration
- **[Logging System](LOGGING_SYSTEM.md)** - Provides logging capabilities to all integrated systems
- **[Projects System](PROJECTS_SYSTEM.md)** - Provides project management and organization through ConfigIntegration

## Architecture

### Core Components

#### Integration Manager
- **File**: `internal/integration/manager.go`
- **Purpose**: Central integration coordination with dependency management and unified access
- **Features**:
  - System component coordination
  - Dependency injection management
  - Integration lifecycle management
  - Error handling and recovery
  - Performance monitoring
  - Health validation

#### Integration Factory
- **File**: `internal/integration/factory.go`
- **Purpose**: Integration instantiation and configuration management
- **Features**:
  - Integration component creation
  - Configuration management
  - Dependency resolution
  - Component initialization
  - Error handling
  - Resource management

#### Integration Interfaces
- **File**: `internal/interfaces/interfaces.go`
- **Purpose**: Core interface definitions for all system integrations
- **Features**:
  - IntegrationManager interface
  - Component integration interfaces
  - Context management interfaces
  - Validation interfaces
  - Error handling interfaces

### Integration Points

#### Actions Integration
- Provides action orchestration and running
- Supports action planning and validation
- Enables action lifecycle management

#### Facts Integration
- Provides fact collection and storage
- Supports fact processing and validation
- Enables fact-based operations

#### Variables Integration
- Provides variable resolution and management
- Supports variable validation and processing
- Enables dynamic variable access

#### Templates Integration
- Provides template rendering and validation
- Supports template processing and caching
- Enables dynamic content generation

#### Machines Integration
- Provides machine inventory and connectivity
- Supports machine targeting and validation
- Enables machine-based operations

#### Secrets Integration
- Provides secrets management and encryption
- Supports secure data handling
- Enables encrypted operations

#### Config Integration
- Provides configuration management
- Supports configuration validation
- Enables dynamic configuration access

## Integration Manager

### Manager Structure
```go
type IntegrationManager struct {
    factsIntegration      spookyinterfaces.FactsIntegration
    actionsIntegration    spookyinterfaces.ActionsIntegration
    variablesIntegration  spookyinterfaces.VariablesIntegration
    templatesIntegration  spookyinterfaces.TemplatesIntegration
    machinesIntegration   spookyinterfaces.MachinesIntegration
    secretsIntegration    spookyinterfaces.SecretsIntegration
    configIntegration     spookyinterfaces.ConfigIntegration
    logger               spookylogging.Logger
    healthStatus         map[string]HealthStatus
    mutex                sync.RWMutex
}
```

### Manager Interface
```go
type IntegrationManager interface {
    // Component access
    GetFactsIntegration() spookyinterfaces.FactsIntegration
    GetActionsIntegration() spookyinterfaces.ActionsIntegration
    GetVariablesIntegration() spookyinterfaces.VariablesIntegration
    GetTemplatesIntegration() spookyinterfaces.TemplatesIntegration
    GetMachinesIntegration() spookyinterfaces.MachinesIntegration
    GetSecretsIntegration() spookyinterfaces.SecretsIntegration
    GetConfigIntegration() spookyinterfaces.ConfigIntegration
    
    // Health and status
    GetHealthStatus() map[string]HealthStatus
    ValidateHealth() error
    
    // Lifecycle management
    Initialize() error
    Shutdown() error
}
```

### Health Status
```go
type HealthStatus struct {
    Component     string                 // Component name
    Status        string                 // Health status
    LastCheck     time.Time              // Last health check
    Error         error                  // Last error
    Metrics       map[string]interface{} // Health metrics
    Dependencies  []string               // Component dependencies
}
```

## Integration Patterns

### Dependency Injection
- **Purpose**: Provide loose coupling between components
- **Features**: Interface-based dependencies, configuration-driven injection
- **Benefits**: Testability, flexibility, maintainability

### Factory Pattern
- **Purpose**: Centralized component creation and configuration
- **Features**: Configuration management, dependency resolution
- **Benefits**: Consistent initialization, error handling

### Coordinator Pattern
- **Purpose**: Centralized system coordination
- **Features**: Component lifecycle management, error handling
- **Benefits**: Unified access, consistent behavior

### Health Monitoring
- **Purpose**: Monitor component health and status
- **Features**: Health checks, dependency validation, error reporting
- **Benefits**: Reliability, observability, troubleshooting

## Integration Lifecycle

### Initialization Process
1. **Configuration Loading**: Load integration configuration
2. **Dependency Resolution**: Resolve component dependencies
3. **Component Creation**: Create integration components
4. **Component Initialization**: Initialize all components
5. **Health Validation**: Validate component health
6. **Service Registration**: Register integration services

### Shutdown Process
1. **Service Deregistration**: Deregister integration services
2. **Component Shutdown**: Shutdown all components
3. **Resource Cleanup**: Clean up allocated resources
4. **Health Finalization**: Finalize health monitoring
5. **Error Handling**: Handle shutdown errors

### Health Monitoring
1. **Component Health Checks**: Check individual component health
2. **Dependency Validation**: Validate component dependencies
3. **Error Aggregation**: Aggregate health errors
4. **Status Reporting**: Report health status
5. **Recovery Actions**: Run recovery actions

## Integration Features

### Unified Access
- **Single Entry Point**: Access all integrations through IntegrationManager
- **Consistent Interface**: Consistent interface patterns across components
- **Error Handling**: Unified error handling and recovery
- **Context Management**: Consistent context management

### Dependency Management
- **Dependency Injection**: Interface-based dependency injection
- **Dependency Resolution**: Automatic dependency resolution
- **Circular Dependency Detection**: Detect and prevent circular dependencies
- **Dependency Validation**: Validate dependency requirements

### Configuration Management
- **Centralized Configuration**: Centralized configuration management
- **Configuration Validation**: Validate configuration integrity
- **Dynamic Configuration**: Support dynamic configuration updates
- **Configuration Persistence**: Persist configuration changes

### Error Handling
- **Error Aggregation**: Aggregate errors from multiple components
- **Error Recovery**: Automatic error recovery mechanisms
- **Error Reporting**: Comprehensive error reporting
- **Error Context**: Maintain error context and stack traces

## Performance Features

### Connection Pooling
- **Resource Reuse**: Reuse connections and resources
- **Connection Limits**: Limit concurrent connections
- **Connection Timeouts**: Implement connection timeouts
- **Connection Health**: Monitor connection health

### Caching
- **Result Caching**: Cache integration results
- **Configuration Caching**: Cache configuration data
- **Health Caching**: Cache health check results
- **Metadata Caching**: Cache component metadata

### Optimization
- **Lazy Loading**: Load components on demand
- **Resource Optimization**: Optimize resource usage
- **Performance Monitoring**: Monitor integration performance
- **Load Balancing**: Balance load across components

## Security Features

### Access Control
- **Component Access Control**: Control access to components
- **Method Access Control**: Control access to methods
- **Data Access Control**: Control access to data
- **Audit Logging**: Log access and operations

### Data Protection
- **Data Encryption**: Encrypt sensitive data
- **Data Validation**: Validate data integrity
- **Data Sanitization**: Sanitize input data
- **Secure Communication**: Secure inter-component communication

### Authentication
- **Component Authentication**: Authenticate component access
- **Method Authentication**: Authenticate method calls
- **Session Management**: Manage authentication sessions
- **Token Management**: Manage authentication tokens

## CLI Integration

### Integration Status
```bash
# Check integration status
spooky integrations status

# Check component health
spooky integrations health

# Check component dependencies
spooky integrations dependencies
```

### Integration Configuration
```bash
# Show integration configuration
spooky integrations config

# Validate integration configuration
spooky integrations validate

# Test integration connectivity
spooky integrations test
```

### Integration Management
```bash
# Initialize integrations
spooky integrations init

# Restart integrations
spooky integrations restart

# Shutdown integrations
spooky integrations shutdown
```

## Configuration

### Integration Configuration
```hcl
# integrations/config.hcl
integration_config {
  # Component settings
  components {
    facts {
      enabled = true
      timeout = 300  # seconds
      retries = 3
    }
    
    actions {
      enabled = true
      timeout = 600  # seconds
      retries = 3
    }
    
    variables {
      enabled = true
      timeout = 60  # seconds
      retries = 2
    }
    
    templates {
      enabled = true
      timeout = 120  # seconds
      retries = 2
    }
    
    machines {
      enabled = true
      timeout = 180  # seconds
      retries = 3
    }
    
    secrets {
      enabled = true
      timeout = 60  # seconds
      retries = 2
    }
    
    config {
      enabled = true
      timeout = 30  # seconds
      retries = 1
    }
  }
  
  # Health monitoring
  health {
    check_interval = 30  # seconds
    timeout = 10  # seconds
    max_failures = 3
  }
  
  # Performance settings
  performance {
    connection_pool_size = 10
    cache_ttl = 300  # seconds
    max_concurrent_operations = 100
  }
  
  # Security settings
  security {
    enable_encryption = true
    enable_audit_logging = true
    access_control_enabled = true
  }
}
```

### Component Configuration
```hcl
# integrations/facts-config.hcl
facts_integration {
  # Collection settings
  collection {
    timeout = 300  # seconds
    max_parallel = 10
    retry_attempts = 3
  }
  
  # Storage settings
  storage {
    backend = "badgerdb"
    location = "~/.local/state/spooky/facts.db"
    encryption = true
  }
  
  # Processing settings
  processing {
    enable_validation = true
    enable_normalization = true
  }
}
```

## Examples

### Basic Integration Usage
```go
// Create integration manager
manager := spookyintegration.NewManager(config)

// Initialize integrations
if err := manager.Initialize(); err != nil {
    log.Fatalf("Failed to initialize integrations: %v", err)
}

// Access facts integration
factsIntegration := manager.GetFactsIntegration()
facts, err := factsIntegration.CollectFacts("web-server")
if err != nil {
    log.Printf("Failed to collect facts: %v", err)
}

// Access actions integration
actionsIntegration := manager.GetActionsIntegration()
actions, err := actionsIntegration.LoadActions("my-project")
if err != nil {
    log.Printf("Failed to load actions: %v", err)
}
```

### Health Monitoring
```go
// Check integration health
healthStatus := manager.GetHealthStatus()
for component, status := range healthStatus {
    if status.Status != "healthy" {
        log.Printf("Component %s is unhealthy: %v", component, status.Error)
    }
}

// Validate overall health
if err := manager.ValidateHealth(); err != nil {
    log.Printf("Integration health validation failed: %v", err)
}
```

### Error Handling
```go
// Handle integration errors
factsIntegration := manager.GetFactsIntegration()
facts, err := factsIntegration.CollectFacts("web-server")
if err != nil {
    // Check if it's a connectivity error
    if spookyinterfaces.IsConnectivityError(err) {
        log.Printf("Connectivity error: %v", err)
        // Implement retry logic
    } else {
        log.Printf("Fact collection error: %v", err)
    }
}
```

## Integration Examples

### Cross-Component Operations
```go
// Use facts in actions
factsIntegration := manager.GetFactsIntegration()
actionsIntegration := manager.GetActionsIntegration()

// Collect facts
facts, err := factsIntegration.CollectFacts("web-server")
if err != nil {
    return err
}

// Use facts in action context
action := &spookytypes.Action{
    Name: "system-info",
    Script: "echo 'CPU: {{.facts.cpu.cores}}'",
    Context: map[string]interface{}{
        "facts": facts,
    },
}

// Run action
result, err := actionsIntegration.RunAction(action)
if err != nil {
    return err
}
```

### Template Rendering with Variables
```go
// Use variables in templates
variablesIntegration := manager.GetVariablesIntegration()
templatesIntegration := manager.GetTemplatesIntegration()

// Load variables
variables, err := variablesIntegration.LoadVariables("my-project")
if err != nil {
    return err
}

// Render template with variables
template, err := templatesIntegration.LoadTemplate("templates/config.tmpl")
if err != nil {
    return err
}

result, err := templatesIntegration.RenderTemplate(template, map[string]interface{}{
    "variables": variables,
})
if err != nil {
    return err
}
```

### Machine Operations
```go
// Use machines with actions
machinesIntegration := manager.GetMachinesIntegration()
actionsIntegration := manager.GetActionsIntegration()

// Get machine inventory
machines, err := machinesIntegration.GetMachines("my-project")
if err != nil {
    return err
}

// Run action on machines
for _, machine := range machines {
    action := &spookytypes.Action{
        Name: "system-update",
        Script: "apt update && apt upgrade -y",
        Machine: machine.Hostname,
    }
    
    result, err := actionsIntegration.RunAction(action)
    if err != nil {
        log.Printf("Failed to update %s: %v", machine.Hostname, err)
    }
}
```

## Best Practices

### Integration Design
- Use interface-based design for loose coupling
- Implement proper error handling and recovery
- Use dependency injection for component management
- Implement health monitoring for reliability

### Performance
- Use connection pooling for resource efficiency
- Implement caching for frequently accessed data
- Monitor integration performance
- Optimize resource usage

### Security
- Implement proper access controls
- Use encryption for sensitive data
- Validate all inputs and outputs
- Audit integration operations

### Maintainability
- Use consistent interface patterns
- Implement comprehensive error handling
- Document integration contracts
- Use proper logging and monitoring

## Troubleshooting

### Common Issues

#### Component Initialization Failures
```bash
# Check component configuration
spooky integrations config

# Check component dependencies
spooky integrations dependencies

# Check component health
spooky integrations health
```

#### Integration Errors
```bash
# Check integration status
spooky integrations status

# Check integration logs
spooky integrations logs

# Test integration connectivity
spooky integrations test
```

#### Performance Issues
```bash
# Check integration performance
spooky integrations performance

# Check resource usage
spooky integrations resources

# Check connection pools
spooky integrations connections
```

#### Health Issues
```bash
# Check component health
spooky integrations health

# Check health history
spooky integrations health-history

# Restart unhealthy components
spooky integrations restart --component <component>
```

## API Reference

### IntegrationManager Interface
```go
type IntegrationManager interface {
    // Component access
    GetFactsIntegration() spookyinterfaces.FactsIntegration
    GetActionsIntegration() spookyinterfaces.ActionsIntegration
    GetVariablesIntegration() spookyinterfaces.VariablesIntegration
    GetTemplatesIntegration() spookyinterfaces.TemplatesIntegration
    GetMachinesIntegration() spookyinterfaces.MachinesIntegration
    GetSecretsIntegration() spookyinterfaces.SecretsIntegration
    GetConfigIntegration() spookyinterfaces.ConfigIntegration
    
    // Health and status
    GetHealthStatus() map[string]HealthStatus
    ValidateHealth() error
    
    // Lifecycle management
    Initialize() error
    Shutdown() error
}
```

### Integration Manager Methods
```go
// Component access
GetFactsIntegration() spookyinterfaces.FactsIntegration
GetActionsIntegration() spookyinterfaces.ActionsIntegration
GetVariablesIntegration() spookyinterfaces.VariablesIntegration
GetTemplatesIntegration() spookyinterfaces.TemplatesIntegration
GetMachinesIntegration() spookyinterfaces.MachinesIntegration
GetSecretsIntegration() spookyinterfaces.SecretsIntegration
GetConfigIntegration() spookyinterfaces.ConfigIntegration

// Health monitoring
GetHealthStatus() map[string]HealthStatus
ValidateHealth() error
CheckComponentHealth(component string) (*HealthStatus, error)

// Lifecycle management
Initialize() error
Shutdown() error
RestartComponent(component string) error
```

## Related Documentation

- [Integrations API Reference](INTEGRATIONS_API_REFERENCE.md) - Complete API documentation
- [Integrations User Guide](INTEGRATIONS_USER_GUIDE.md) - User guide and examples
- [Integrations Troubleshooting](INTEGRATIONS_TROUBLESHOOTING.md) - Troubleshooting guide
- [Actions System](ACTIONS_SYSTEM.md) - Actions integration
- [Facts System](FACTS_SYSTEM.md) - Facts integration
- [Variables System](VARIABLES_SYSTEM.md) - Variables integration
- [Templates System](TEMPLATES_SYSTEM.md) - Templates integration
- [Machines System](MACHINES_SYSTEM.md) - Machines integration
- [Secrets System](SECRETS_SYSTEM.md) - Secrets integration
- [Config System](CONFIG_SYSTEM.md) - Config integration
