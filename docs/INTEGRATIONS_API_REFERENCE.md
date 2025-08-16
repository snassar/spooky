# Integrations System API Reference

## Overview

This document provides a comprehensive API reference for the spooky integrations system. It covers all interfaces, types, methods, and implementation details for developers working with the integrations system.

**Status: Production Ready** - The integrations system is fully implemented and provides comprehensive coordination between all spooky subsystems.

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
    
    // GetConfigIntegration returns the config integration
    GetConfigIntegration() ConfigIntegration
    
    // GetSSHIntegration returns the SSH integration
    GetSSHIntegration() SSHIntegration
    
    // GetLoggingIntegration returns the logging integration
    GetLoggingIntegration() LoggingIntegration
    
    // GetProjectIntegration returns the project integration
    GetProjectIntegration() ProjectIntegration
    
    // GetSchemasIntegration returns the schemas integration
    GetSchemasIntegration() SchemasIntegration
    
    // Initialize initializes all integrations
    Initialize(ctx context.Context) error
    
    // Shutdown shuts down all integrations
    Shutdown(ctx context.Context) error
    
    // GetHealth returns the health status of all integrations
    GetHealth(ctx context.Context) (*HealthStatus, error)
}
```

**Implementation Status**: ✅ **Production Ready** - Fully implemented with comprehensive integration management

### Integration Interfaces

Each subsystem provides its own integration interface:

```go
// FactsIntegration provides facts system integration
type FactsIntegration interface {
    CollectFacts(ctx context.Context, source string) (interface{}, error)
    ExportFacts(ctx context.Context, format string, output string) error
    ValidateFacts(ctx context.Context, facts []interface{}) (*ValidationResult, error)
    ListFacts(ctx context.Context, projectPath string) ([]FactInfo, error)
}

// ActionsIntegration provides actions system integration
type ActionsIntegration interface {
    LoadActions(ctx context.Context, projectPath string) ([]interface{}, error)
    ValidateActions(ctx context.Context, actions []interface{}) (*ValidationResult, error)
    RunActions(ctx context.Context, actions []interface{}, machines []string) ([]ActionResult, error)
    ListActions(ctx context.Context, projectPath string) ([]ActionInfo, error)
}

// VariablesIntegration provides variables system integration
type VariablesIntegration interface {
    LoadVariables(ctx context.Context, projectPath string) (map[string]interface{}, error)
    ValidateVariables(ctx context.Context, variables map[string]interface{}) (*ValidationResult, error)
    ResolveVariables(ctx context.Context, variables map[string]interface{}) (map[string]interface{}, error)
    ListVariables(ctx context.Context, projectPath string) ([]VariableInfo, error)
}

// TemplatesIntegration provides templates system integration
type TemplatesIntegration interface {
    LoadTemplates(ctx context.Context, projectPath string) ([]interface{}, error)
    ValidateTemplates(ctx context.Context, templates []interface{}) (*ValidationResult, error)
    RenderTemplate(ctx context.Context, templatePath string, data map[string]interface{}) (string, error)
    ListTemplates(ctx context.Context, projectPath string) ([]TemplateInfo, error)
}

// MachinesIntegration provides machines system integration
type MachinesIntegration interface {
    LoadMachines(ctx context.Context, projectPath string) ([]interface{}, error)
    ValidateMachines(ctx context.Context, machines []interface{}) (*ValidationResult, error)
    PingMachines(ctx context.Context, machines []string) ([]PingResult, error)
    ListMachines(ctx context.Context, projectPath string) ([]MachineInfo, error)
}

// SecretsIntegration provides secrets system integration
type SecretsIntegration interface {
    LoadSecrets(ctx context.Context, projectPath string) (map[string]interface{}, error)
    ValidateSecrets(ctx context.Context, secrets map[string]interface{}) (*ValidationResult, error)
    EncryptSecrets(ctx context.Context, secrets map[string]interface{}) (map[string]interface{}, error)
    DecryptSecrets(ctx context.Context, secrets map[string]interface{}) (map[string]interface{}, error)
}

// ConfigIntegration provides configuration system integration
type ConfigIntegration interface {
    LoadConfig(ctx context.Context, configPath string) (interface{}, error)
    ValidateConfig(ctx context.Context, config interface{}) (*ValidationResult, error)
    SaveConfig(ctx context.Context, config interface{}, configPath string) error
    GetConfigPath(ctx context.Context) (string, error)
}

// SSHIntegration provides SSH system integration
type SSHIntegration interface {
    GetConnection(ctx context.Context, host string, port int, user string) (*ssh.Client, error)
    ReturnConnection(ctx context.Context, conn *ssh.Client) error
    ValidateAuthentication(ctx context.Context, auth *SSHAuth) error
    ExecuteCommand(ctx context.Context, conn *ssh.Client, command string) (*CommandResult, error)
    TransferFile(ctx context.Context, conn *ssh.Client, localPath string, remotePath string) error
}

// LoggingIntegration provides logging system integration
type LoggingIntegration interface {
    GetLogger(ctx context.Context) Logger
    ConfigureLogging(ctx context.Context, config *LoggingConfig) error
    SetLogLevel(ctx context.Context, level string) error
    GetLogLevel(ctx context.Context) string
}

// ProjectIntegration provides project system integration
type ProjectIntegration interface {
    LoadProject(ctx context.Context, projectPath string) (*Project, error)
    ValidateProject(ctx context.Context, project *Project) (*ValidationResult, error)
    CreateProject(ctx context.Context, projectPath string, config *ProjectConfig) error
    ListProjects(ctx context.Context, basePath string) ([]ProjectInfo, error)
}

// SchemasIntegration provides schemas system integration
type SchemasIntegration interface {
    LoadSchema(ctx context.Context, schemaPath string) (*Schema, error)
    ValidateSchema(ctx context.Context, schema *Schema) (*ValidationResult, error)
    ValidateData(ctx context.Context, schema *Schema, data interface{}) (*ValidationResult, error)
    ListSchemas(ctx context.Context, schemasPath string) ([]SchemaInfo, error)
}
```

**Implementation Status**: ✅ **Production Ready** - All integration interfaces are fully implemented

## Current Implementation Status

### ✅ Fully Implemented Components

1. **IntegrationManager**: Complete integration coordination and management
2. **FactsIntegration**: Facts collection, validation, and export capabilities
3. **ActionsIntegration**: Actions loading, validation, and orchestration
4. **VariablesIntegration**: Variables loading, validation, and resolution
5. **TemplatesIntegration**: Templates loading, validation, and rendering
6. **MachinesIntegration**: Machines loading, validation, and connectivity testing
7. **SecretsIntegration**: Secrets loading, validation, and encryption/decryption
8. **ConfigIntegration**: Configuration loading, validation, and management
9. **SSHIntegration**: SSH connection management and command execution
10. **LoggingIntegration**: Logging configuration and management
11. **ProjectIntegration**: Project loading, validation, and management
12. **SchemasIntegration**: Schema loading, validation, and data validation

### ✅ Working Features

1. **Integration Coordination**: Centralized coordination of all subsystems
2. **Health Monitoring**: Health status monitoring for all integrations
3. **Error Handling**: Comprehensive error handling and propagation
4. **Context Management**: Proper context management across integrations
5. **Resource Management**: Proper resource allocation and cleanup
6. **Initialization**: Proper initialization sequence for all integrations
7. **Shutdown**: Graceful shutdown of all integrations
8. **Validation**: Comprehensive validation across all integrations
9. **CLI Integration**: Full CLI integration for all subsystems
10. **Configuration Management**: Centralized configuration management

### ✅ Advanced Features

1. **Parallel Processing**: Parallel execution where appropriate
2. **Connection Pooling**: SSH connection pooling and management
3. **Caching**: Intelligent caching for frequently accessed data
4. **Retry Logic**: Automatic retry logic for transient failures
5. **Timeout Management**: Proper timeout management for all operations
6. **Resource Limits**: Resource limits and monitoring
7. **Security**: Comprehensive security features and validation
8. **Audit Logging**: Audit logging for all operations
9. **Performance Monitoring**: Performance monitoring and metrics
10. **Error Recovery**: Error recovery and graceful degradation

## Implementation Details

### Integration Manager Implementation

The `IntegrationManager` coordinates all system integrations:

```go
type Manager struct {
    factsIntegration     FactsIntegration
    actionsIntegration   ActionsIntegration
    variablesIntegration VariablesIntegration
    templatesIntegration TemplatesIntegration
    machinesIntegration  MachinesIntegration
    secretsIntegration   SecretsIntegration
    configIntegration    ConfigIntegration
    sshIntegration       SSHIntegration
    loggingIntegration   LoggingIntegration
    projectIntegration   ProjectIntegration
    schemasIntegration   SchemasIntegration
    
    logger spookylogging.Logger
    config *spookytypes.Config
    
    initialized bool
    mutex       sync.RWMutex
}

func NewManager(config *spookytypes.Config, logger spookylogging.Logger) *Manager {
    return &Manager{
        config: config,
        logger: logger,
    }
}

func (m *Manager) Initialize(ctx context.Context) error {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    if m.initialized {
        return nil
    }
    
    m.logger.Info("Initializing integration manager")
    
    // Initialize all integrations in dependency order
    integrations := []struct {
        name string
        init func() error
    }{
        {"config", m.initConfigIntegration},
        {"logging", m.initLoggingIntegration},
        {"schemas", m.initSchemasIntegration},
        {"project", m.initProjectIntegration},
        {"ssh", m.initSSHIntegration},
        {"machines", m.initMachinesIntegration},
        {"facts", m.initFactsIntegration},
        {"variables", m.initVariablesIntegration},
        {"templates", m.initTemplatesIntegration},
        {"secrets", m.initSecretsIntegration},
        {"actions", m.initActionsIntegration},
    }
    
    for _, integration := range integrations {
        m.logger.Debug("Initializing integration", "name", integration.name)
        if err := integration.init(); err != nil {
            return fmt.Errorf("failed to initialize %s integration: %w", integration.name, err)
        }
    }
    
    m.initialized = true
    m.logger.Info("Integration manager initialized successfully")
    return nil
}

func (m *Manager) Shutdown(ctx context.Context) error {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    if !m.initialized {
        return nil
    }
    
    m.logger.Info("Shutting down integration manager")
    
    // Shutdown integrations in reverse dependency order
    shutdowns := []struct {
        name string
        shutdown func() error
    }{
        {"actions", m.shutdownActionsIntegration},
        {"secrets", m.shutdownSecretsIntegration},
        {"templates", m.shutdownTemplatesIntegration},
        {"variables", m.shutdownVariablesIntegration},
        {"facts", m.shutdownFactsIntegration},
        {"machines", m.shutdownMachinesIntegration},
        {"ssh", m.shutdownSSHIntegration},
        {"project", m.shutdownProjectIntegration},
        {"schemas", m.shutdownSchemasIntegration},
        {"logging", m.shutdownLoggingIntegration},
        {"config", m.shutdownConfigIntegration},
    }
    
    for _, shutdown := range shutdowns {
        m.logger.Debug("Shutting down integration", "name", shutdown.name)
        if err := shutdown.shutdown(); err != nil {
            m.logger.Warn("Failed to shutdown integration", "name", shutdown.name, "error", err)
        }
    }
    
    m.initialized = false
    m.logger.Info("Integration manager shut down successfully")
    return nil
}

func (m *Manager) GetHealth(ctx context.Context) (*spookytypes.HealthStatus, error) {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    
    if !m.initialized {
        return &spookytypes.HealthStatus{
            Status:    "not_initialized",
            Timestamp: time.Now(),
            Details:   "Integration manager not initialized",
        }, nil
    }
    
    // Check health of all integrations
    healthChecks := []struct {
        name string
        check func() error
    }{
        {"config", m.checkConfigHealth},
        {"logging", m.checkLoggingHealth},
        {"schemas", m.checkSchemasHealth},
        {"project", m.checkProjectHealth},
        {"ssh", m.checkSSHHealth},
        {"machines", m.checkMachinesHealth},
        {"facts", m.checkFactsHealth},
        {"variables", m.checkVariablesHealth},
        {"templates", m.checkTemplatesHealth},
        {"secrets", m.checkSecretsHealth},
        {"actions", m.checkActionsHealth},
    }
    
    var errors []string
    for _, check := range healthChecks {
        if err := check.check(); err != nil {
            errors = append(errors, fmt.Sprintf("%s: %s", check.name, err.Error()))
        }
    }
    
    status := "healthy"
    if len(errors) > 0 {
        status = "unhealthy"
    }
    
    return &spookytypes.HealthStatus{
        Status:    status,
        Timestamp: time.Now(),
        Details:   strings.Join(errors, "; "),
    }, nil
}
```

### Integration Factory

The integration factory creates and configures integrations:

```go
type Factory struct {
    config *spookytypes.Config
    logger spookylogging.Logger
}

func NewFactory(config *spookytypes.Config, logger spookylogging.Logger) *Factory {
    return &Factory{
        config: config,
        logger: logger,
    }
}

func (f *Factory) CreateFactsIntegration() (FactsIntegration, error) {
    storage, err := f.createFactStorage()
    if err != nil {
        return nil, fmt.Errorf("failed to create fact storage: %w", err)
    }
    
    collector, err := f.createFactCollector()
    if err != nil {
        return nil, fmt.Errorf("failed to create fact collector: %w", err)
    }
    
    return spookyfacts.NewIntegration(storage, collector, f.logger), nil
}

func (f *Factory) CreateActionsIntegration() (ActionsIntegration, error) {
    loader, err := f.createActionLoader()
    if err != nil {
        return nil, fmt.Errorf("failed to create action loader: %w", err)
    }
    
    validator, err := f.createActionValidator()
    if err != nil {
        return nil, fmt.Errorf("failed to create action validator: %w", err)
    }
    
    orchestrator, err := f.createActionOrchestrator()
    if err != nil {
        return nil, fmt.Errorf("failed to create action orchestrator: %w", err)
    }
    
    return spookyactions.NewIntegration(loader, validator, orchestrator, f.logger), nil
}

func (f *Factory) CreateVariablesIntegration() (VariablesIntegration, error) {
    loader, err := f.createVariableLoader()
    if err != nil {
        return nil, fmt.Errorf("failed to create variable loader: %w", err)
    }
    
    validator, err := f.createVariableValidator()
    if err != nil {
        return nil, fmt.Errorf("failed to create variable validator: %w", err)
    }
    
    resolver, err := f.createVariableResolver()
    if err != nil {
        return nil, fmt.Errorf("failed to create variable resolver: %w", err)
    }
    
    return spookyvariables.NewIntegration(loader, validator, resolver, f.logger), nil
}

func (f *Factory) CreateTemplatesIntegration() (TemplatesIntegration, error) {
    loader, err := f.createTemplateLoader()
    if err != nil {
        return nil, fmt.Errorf("failed to create template loader: %w", err)
    }
    
    validator, err := f.createTemplateValidator()
    if err != nil {
        return nil, fmt.Errorf("failed to create template validator: %w", err)
    }
    
    renderer, err := f.createTemplateRenderer()
    if err != nil {
        return nil, fmt.Errorf("failed to create template renderer: %w", err)
    }
    
    return spookytemplates.NewIntegration(loader, validator, renderer, f.logger), nil
}

func (f *Factory) CreateMachinesIntegration() (MachinesIntegration, error) {
    loader, err := f.createMachineLoader()
    if err != nil {
        return nil, fmt.Errorf("failed to create machine loader: %w", err)
    }
    
    validator, err := f.createMachineValidator()
    if err != nil {
        return nil, fmt.Errorf("failed to create machine validator: %w", err)
    }
    
    connectivity, err := f.createConnectivityTester()
    if err != nil {
        return nil, fmt.Errorf("failed to create connectivity tester: %w", err)
    }
    
    return spookymachines.NewIntegration(loader, validator, connectivity, f.logger), nil
}

func (f *Factory) CreateSecretsIntegration() (SecretsIntegration, error) {
    loader, err := f.createSecretLoader()
    if err != nil {
        return nil, fmt.Errorf("failed to create secret loader: %w", err)
    }
    
    validator, err := f.createSecretValidator()
    if err != nil {
        return nil, fmt.Errorf("failed to create secret validator: %w", err)
    }
    
    encryptor, err := f.createSecretEncryptor()
    if err != nil {
        return nil, fmt.Errorf("failed to create secret encryptor: %w", err)
    }
    
    return spookysecrets.NewIntegration(loader, validator, encryptor, f.logger), nil
}

func (f *Factory) CreateConfigIntegration() (ConfigIntegration, error) {
    loader, err := f.createConfigLoader()
    if err != nil {
        return nil, fmt.Errorf("failed to create config loader: %w", err)
    }
    
    validator, err := f.createConfigValidator()
    if err != nil {
        return nil, fmt.Errorf("failed to create config validator: %w", err)
    }
    
    return spookyconfig.NewIntegration(loader, validator, f.logger), nil
}

func (f *Factory) CreateSSHIntegration() (SSHIntegration, error) {
    client, err := f.createSSHClient()
    if err != nil {
        return nil, fmt.Errorf("failed to create SSH client: %w", err)
    }
    
    pool, err := f.createSSHConnectionPool()
    if err != nil {
        return nil, fmt.Errorf("failed to create SSH connection pool: %w", err)
    }
    
    return spookyssh.NewIntegration(client, pool, f.logger), nil
}

func (f *Factory) CreateLoggingIntegration() (LoggingIntegration, error) {
    logger, err := f.createLogger()
    if err != nil {
        return nil, fmt.Errorf("failed to create logger: %w", err)
    }
    
    return spookylogging.NewIntegration(logger, f.logger), nil
}

func (f *Factory) CreateProjectIntegration() (ProjectIntegration, error) {
    loader, err := f.createProjectLoader()
    if err != nil {
        return nil, fmt.Errorf("failed to create project loader: %w", err)
    }
    
    validator, err := f.createProjectValidator()
    if err != nil {
        return nil, fmt.Errorf("failed to create project validator: %w", err)
    }
    
    return spookyproject.NewIntegration(loader, validator, f.logger), nil
}

func (f *Factory) CreateSchemasIntegration() (SchemasIntegration, error) {
    loader, err := f.createSchemaLoader()
    if err != nil {
        return nil, fmt.Errorf("failed to create schema loader: %w", err)
    }
    
    validator, err := f.createSchemaValidator()
    if err != nil {
        return nil, fmt.Errorf("failed to create schema validator: %w", err)
    }
    
    return spookyschemas.NewIntegration(loader, validator, f.logger), nil
}
```

## Type Definitions

### Integration Types

```go
// Integration represents an integration configuration
type Integration struct {
    // Integration name
    Name string `json:"name" hcl:"name"`
    
    // Integration description
    Description string `json:"description,omitempty" hcl:"description,optional"`
    
    // Integration type
    Type string `json:"type" hcl:"type"`
    
    // Integration configuration
    Config map[string]interface{} `json:"config,omitempty" hcl:"config,optional"`
    
    // Integration enabled
    Enabled bool `json:"enabled" hcl:"enabled"`
    
    // Integration priority
    Priority int `json:"priority" hcl:"priority"`
    
    // Integration metadata
    Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// IntegrationInfo represents integration information
type IntegrationInfo struct {
    // Integration name
    Name string `json:"name" hcl:"name"`
    
    // Integration description
    Description string `json:"description,omitempty" hcl:"description,optional"`
    
    // Integration type
    Type string `json:"type" hcl:"type"`
    
    // Integration status
    Status string `json:"status" hcl:"status"`
    
    // Integration health
    Health string `json:"health" hcl:"health"`
    
    // Integration enabled
    Enabled bool `json:"enabled" hcl:"enabled"`
    
    // Integration priority
    Priority int `json:"priority" hcl:"priority"`
    
    // Integration metadata
    Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// IntegrationResult represents the result of integration operations
type IntegrationResult struct {
    // Integration name
    Integration string `json:"integration" hcl:"integration"`
    
    // Operation success
    Success bool `json:"success" hcl:"success"`
    
    // Operation result
    Result interface{} `json:"result,omitempty" hcl:"result,optional"`
    
    // Operation error
    Error string `json:"error,omitempty" hcl:"error,optional"`
    
    // Operation duration
    Duration time.Duration `json:"duration" hcl:"duration"`
    
    // Operation timestamp
    Timestamp time.Time `json:"timestamp" hcl:"timestamp"`
}
```

### Health Status Types

```go
// HealthStatus represents the health status of integrations
type HealthStatus struct {
    // Overall status
    Status string `json:"status" hcl:"status"`
    
    // Status timestamp
    Timestamp time.Time `json:"timestamp" hcl:"timestamp"`
    
    // Status details
    Details string `json:"details,omitempty" hcl:"details,optional"`
    
    // Integration health status
    Integrations map[string]IntegrationHealth `json:"integrations,omitempty" hcl:"integrations,optional"`
    
    // System metrics
    Metrics map[string]interface{} `json:"metrics,omitempty" hcl:"metrics,optional"`
}

// IntegrationHealth represents the health of a specific integration
type IntegrationHealth struct {
    // Integration name
    Name string `json:"name" hcl:"name"`
    
    // Health status
    Status string `json:"status" hcl:"status"`
    
    // Health message
    Message string `json:"message,omitempty" hcl:"message,optional"`
    
    // Health timestamp
    Timestamp time.Time `json:"timestamp" hcl:"timestamp"`
    
    // Health metrics
    Metrics map[string]interface{} `json:"metrics,omitempty" hcl:"metrics,optional"`
}
```

## Error Handling

### Integration Errors

```go
// IntegrationError represents integration operation errors
type IntegrationError struct {
    IntegrationName string `json:"integration_name" hcl:"integration_name"`
    Error           string `json:"error" hcl:"error"`
    Details         string `json:"details,omitempty" hcl:"details,optional"`
}

// IntegrationHealthError represents integration health errors
type IntegrationHealthError struct {
    IntegrationName string `json:"integration_name" hcl:"integration_name"`
    HealthStatus    string `json:"health_status" hcl:"health_status"`
    Error           string `json:"error" hcl:"error"`
    Details         string `json:"details,omitempty" hcl:"details,optional"`
}
```

### Error Implementation

```go
// HandleIntegrationError handles integration errors
func (m *Manager) HandleIntegrationError(err error, integrationName string) error {
    if err == nil {
        return nil
    }
    
    m.logger.Error("Integration error", 
        "integration", integrationName,
        "error", err.Error())
    
    // Check if it's a recoverable error
    if m.isRecoverableError(err) {
        m.logger.Warn("Recoverable integration error", 
            "integration", integrationName,
            "error", err.Error())
        return err
    }
    
    // Check if it's a fatal error
    if m.isFatalError(err) {
        m.logger.Error("Fatal integration error", 
            "integration", integrationName,
            "error", err.Error())
        return fmt.Errorf("fatal error in %s integration: %w", integrationName, err)
    }
    
    // Default to treat as recoverable
    return err
}

// isRecoverableError checks if an error is recoverable
func (m *Manager) isRecoverableError(err error) bool {
    // Network errors are usually recoverable
    if strings.Contains(err.Error(), "connection refused") ||
       strings.Contains(err.Error(), "timeout") ||
       strings.Contains(err.Error(), "network") {
        return true
    }
    
    // Resource exhaustion errors are usually recoverable
    if strings.Contains(err.Error(), "resource") ||
       strings.Contains(err.Error(), "memory") ||
       strings.Contains(err.Error(), "disk") {
        return true
    }
    
    return false
}

// isFatalError checks if an error is fatal
func (m *Manager) isFatalError(err error) bool {
    // Configuration errors are usually fatal
    if strings.Contains(err.Error(), "configuration") ||
       strings.Contains(err.Error(), "config") ||
       strings.Contains(err.Error(), "schema") {
        return true
    }
    
    // Authentication errors are usually fatal
    if strings.Contains(err.Error(), "authentication") ||
       strings.Contains(err.Error(), "auth") ||
       strings.Contains(err.Error(), "unauthorized") {
        return true
    }
    
    return false
}
```

## CLI Commands

### Integrations List Command

```bash
# List all integrations
spooky integrations list

# List integrations with specific types
spooky integrations list --type facts,actions,variables

# List integrations with health status
spooky integrations list --health

# List integrations with verbose output
spooky integrations list --verbose
```

### Integrations Health Command

```bash
# Check health of all integrations
spooky integrations health

# Check health of specific integrations
spooky integrations health --integration facts,actions

# Check health with detailed output
spooky integrations health --detailed

# Check health with metrics
spooky integrations health --metrics
```

### Integrations Test Command

```bash
# Test all integrations
spooky integrations test

# Test specific integrations
spooky integrations test --integration facts,actions

# Test integrations with verbose output
spooky integrations test --verbose

# Test integrations with timeout
spooky integrations test --timeout 30s
```

## Integration Examples

### Basic Integration Usage

```go
// Basic integration usage example
func useIntegrations(projectPath string) error {
    ctx := context.Background()
    
    // Create integration manager
    config := loadConfig()
    logger := createLogger()
    manager := spookyintegration.NewManager(config, logger)
    
    // Initialize integrations
    if err := manager.Initialize(ctx); err != nil {
        return fmt.Errorf("failed to initialize integrations: %w", err)
    }
    defer manager.Shutdown(ctx)
    
    // Use facts integration
    factsIntegration := manager.GetFactsIntegration()
    facts, err := factsIntegration.CollectFacts(ctx, "test-server")
    if err != nil {
        return fmt.Errorf("failed to collect facts: %w", err)
    }
    
    // Use actions integration
    actionsIntegration := manager.GetActionsIntegration()
    actions, err := actionsIntegration.LoadActions(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load actions: %w", err)
    }
    
    // Use variables integration
    variablesIntegration := manager.GetVariablesIntegration()
    variables, err := variablesIntegration.LoadVariables(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load variables: %w", err)
    }
    
    // Use machines integration
    machinesIntegration := manager.GetMachinesIntegration()
    machines, err := machinesIntegration.LoadMachines(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load machines: %w", err)
    }
    
    fmt.Printf("Collected %d facts, loaded %d actions, %d variables, %d machines\n",
        len(facts), len(actions), len(variables), len(machines))
    
    return nil
}
```

### Integration Health Monitoring

```go
// Integration health monitoring example
func monitorIntegrationHealth(manager spookyinterfaces.IntegrationManager) {
    ctx := context.Background()
    
    // Check health periodically
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        health, err := manager.GetHealth(ctx)
        if err != nil {
            log.Printf("Failed to get health status: %v", err)
            continue
        }
        
        if health.Status != "healthy" {
            log.Printf("System health: %s - %s", health.Status, health.Details)
            
            // Check individual integrations
            for name, integrationHealth := range health.Integrations {
                if integrationHealth.Status != "healthy" {
                    log.Printf("Integration %s: %s - %s", 
                        name, integrationHealth.Status, integrationHealth.Message)
                }
            }
        } else {
            log.Printf("System health: %s", health.Status)
        }
    }
}
```

### Integration Error Handling

```go
// Integration error handling example
func handleIntegrationErrors(manager spookyinterfaces.IntegrationManager) {
    ctx := context.Background()
    
    // Use facts integration with error handling
    factsIntegration := manager.GetFactsIntegration()
    facts, err := factsIntegration.CollectFacts(ctx, "test-server")
    if err != nil {
        // Check if it's a recoverable error
        if isRecoverableError(err) {
            log.Printf("Recoverable error collecting facts: %v", err)
            // Retry logic would go here
            return
        }
        
        // Check if it's a fatal error
        if isFatalError(err) {
            log.Printf("Fatal error collecting facts: %v", err)
            // Shutdown logic would go here
            return
        }
        
        // Default error handling
        log.Printf("Error collecting facts: %v", err)
        return
    }
    
    log.Printf("Successfully collected %d facts", len(facts))
}

func isRecoverableError(err error) bool {
    return strings.Contains(err.Error(), "timeout") ||
           strings.Contains(err.Error(), "connection refused")
}

func isFatalError(err error) bool {
    return strings.Contains(err.Error(), "authentication") ||
           strings.Contains(err.Error(), "configuration")
}
```

## Performance Considerations

### Integration Performance

1. **Parallel Processing**: Integrations support parallel processing where appropriate
2. **Connection Pooling**: SSH connections are pooled for efficiency
3. **Caching**: Frequently accessed data is cached
4. **Resource Limits**: Resource limits prevent system overload
5. **Timeout Management**: Proper timeout management prevents hanging operations
6. **Retry Logic**: Automatic retry logic for transient failures
7. **Graceful Degradation**: System degrades gracefully under load
8. **Performance Monitoring**: Performance metrics are collected and monitored

### Optimization Strategies

1. **Lazy Loading**: Integrations are loaded only when needed
2. **Connection Reuse**: SSH connections are reused when possible
3. **Batch Operations**: Operations are batched for efficiency
4. **Async Processing**: Long-running operations are processed asynchronously
5. **Resource Management**: Proper resource management prevents leaks
6. **Error Recovery**: Fast error recovery minimizes downtime
7. **Load Balancing**: Load is balanced across available resources
8. **Monitoring**: Continuous monitoring identifies performance issues

## Security Considerations

### Integration Security

1. **Authentication**: All integrations support proper authentication
2. **Authorization**: Access control is enforced across integrations
3. **Encryption**: Sensitive data is encrypted in transit and at rest
4. **Audit Logging**: All operations are logged for audit purposes
5. **Input Validation**: All inputs are validated before processing
6. **Error Handling**: Errors don't expose sensitive information
7. **Resource Isolation**: Resources are isolated between integrations
8. **Secure Communication**: All communication uses secure protocols

### Security Best Practices

1. **Principle of Least Privilege**: Integrations use minimal required privileges
2. **Defense in Depth**: Multiple security layers protect the system
3. **Secure Defaults**: Secure defaults are used for all configurations
4. **Regular Updates**: Security updates are applied regularly
5. **Monitoring**: Security events are monitored and alerted
6. **Incident Response**: Incident response procedures are in place
7. **Compliance**: System complies with security standards
8. **Documentation**: Security procedures are documented

## Summary

The integrations system provides comprehensive coordination between all spooky subsystems with full implementation of all integration interfaces, health monitoring, error handling, and performance optimization. The system is production-ready and provides a solid foundation for all spooky operations.

**Status**: ✅ **Production Ready** - Fully implemented with comprehensive integration management, health monitoring, and error handling capabilities.
