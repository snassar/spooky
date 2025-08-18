# Interfaces System API Reference

## Overview

This document provides a comprehensive API reference for the spooky interfaces system. It covers all interfaces, types, methods, and implementation details for developers working with the interfaces system.

**Status: Implemented** - The interfaces system provides comprehensive functionality for system coordination and integration.

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
}
```

### ProjectManager Interface

The `ProjectManager` interface manages project lifecycle and operations:

```go
type ProjectManager interface {
    // Initialize initializes a new project
    Initialize(ctx context.Context, projectPath string) (*spookytypes.Project, error)
    
    // Load loads a project from the given path
    Load(ctx context.Context, projectPath string) (*spookytypes.Project, error)
    
    // Validate validates a project
    Validate(ctx context.Context, project *spookytypes.Project) (*spookytypes.ValidationResult, error)
    
    // Save saves a project to disk
    Save(ctx context.Context, project *spookytypes.Project) error
    
    // Delete deletes a project
    Delete(ctx context.Context, projectPath string) error
}
```

### ProjectValidator Interface

The `ProjectValidator` interface validates project structure and configuration:

```go
type ProjectValidator interface {
    // ValidateProject validates a project structure
    ValidateProject(ctx context.Context, project *spookytypes.Project) (*spookytypes.ValidationResult, error)
    
    // ValidateProjectDirectory validates project directory structure
    ValidateProjectDirectory(ctx context.Context, projectPath string) (*spookytypes.ValidationResult, error)
    
    // ValidateProjectConfig validates project configuration
    ValidateProjectConfig(ctx context.Context, config *spookytypes.ProjectConfig) (*spookytypes.ValidationResult, error)
}
```

### ProjectLoader Interface

The `ProjectLoader` interface loads project data from various sources:

```go
type ProjectLoader interface {
    // LoadProject loads a project from disk
    LoadProject(ctx context.Context, projectPath string) (*spookytypes.Project, error)
    
    // LoadProjectConfig loads project configuration
    LoadProjectConfig(ctx context.Context, projectPath string) (*spookytypes.ProjectConfig, error)
    
    // LoadProjectMetadata loads project metadata
    LoadProjectMetadata(ctx context.Context, projectPath string) (*spookytypes.ProjectMetadata, error)
}
```

### ConfigManager Interface

The `ConfigManager` interface manages configuration loading and validation:

```go
type ConfigManager interface {
    // Load loads configuration from the given path
    Load(ctx context.Context, configPath string) (*spookytypes.Config, error)
    
    // Validate validates configuration
    Validate(ctx context.Context, config *spookytypes.Config) (*spookytypes.ValidationResult, error)
    
    // Save saves configuration to disk
    Save(ctx context.Context, config *spookytypes.Config, configPath string) error
    
    // GetDefaultConfig returns default configuration
    GetDefaultConfig() *spookytypes.Config
}
```

### LogManager Interface

The `LogManager` interface manages logging operations and logger instances:

```go
type LogManager interface {
    // GetLogger returns a logger for the given component
    GetLogger(component string) spookytypes.Logger
    
    // SetLevel sets the log level for all loggers
    SetLevel(level spookytypes.LogLevel)
    
    // GetLevel returns the current log level
    GetLevel() spookytypes.LogLevel
    
    // Configure configures logging with the given configuration
    Configure(config *spookytypes.LogConfig) error
    
    // Flush flushes all pending log entries
    Flush() error
    
    // Close closes the log manager
    Close() error
}
```

### SchemaManager Interface

The `SchemaManager` interface manages schema loading and validation:

```go
type SchemaManager interface {
    // LoadSchema loads a schema from the given path
    LoadSchema(ctx context.Context, schemaPath string) (*spookytypes.Schema, error)
    
    // LoadEmbeddedSchema loads an embedded schema
    LoadEmbeddedSchema(ctx context.Context, schemaName string) (*spookytypes.Schema, error)
    
    // Validate validates data against a schema
    Validate(ctx context.Context, schema *spookytypes.Schema, data interface{}) (*spookytypes.ValidationResult, error)
    
    // Register registers a new schema
    Register(ctx context.Context, schema *spookytypes.Schema) error
}
```

### CLIManager Interface

The `CLIManager` interface manages command-line interface operations:

```go
type CLIManager interface {
    // RegisterCommand registers a new command
    RegisterCommand(command spookytypes.Command) error
    
    // RunCommand runs a command
    RunCommand(ctx context.Context, commandName string, args []string) error
    
    // GetCommand returns a command by name
    GetCommand(commandName string) (spookytypes.Command, bool)
    
    // ListCommands returns all registered commands
    ListCommands() []spookytypes.Command
    
    // ShowHelp shows help for a command
    ShowHelp(commandName string) error
    
    // ShowVersion shows version information
    ShowVersion() error
}
```

### FactsIntegration Interface

The `FactsIntegration` interface provides fact collection and memory storage:

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

### ActionsIntegration Interface

The `ActionsIntegration` interface provides action management and orchestration:

```go
type ActionsIntegration interface {
    // LoadActions loads actions from the given source
    LoadActions(ctx context.Context, source string) ([]spookytypes.Action, error)
    
    // ValidateActions validates actions
    ValidateActions(ctx context.Context, actions []spookytypes.Action) (*spookytypes.ValidationResult, error)
    
    // RunActions runs actions on the given machines
    RunActions(ctx context.Context, actions []spookytypes.Action, machines []spookytypes.Machine) ([]spookytypes.ActingResult, error)
    
    // GetSSHManager returns the SSH manager for authentication testing
    GetSSHManager() SSHManager
}
```

### VariablesIntegration Interface

The `VariablesIntegration` interface provides variable management:

```go
type VariablesIntegration interface {
    // LoadVariables loads variables from the given source
    LoadVariables(ctx context.Context, source string) (map[string]*spookytypes.Variable, error)
    
    // ResolveVariables resolves variables with the given context
    ResolveVariables(ctx context.Context, variables map[string]*spookytypes.Variable, context *spookytypes.VariableContext) (*spookytypes.VariableResolutionResult, error)
    
    // ValidateVariables validates variables
    ValidateVariables(ctx context.Context, variables map[string]*spookytypes.Variable) (*spookytypes.ValidationResult, error)
    
    // SaveVariables saves variables to the given destination
    SaveVariables(ctx context.Context, variables map[string]*spookytypes.Variable, destination string) error
    
    // EncryptVariables encrypts all variables that have encrypted=true
    EncryptVariables(ctx context.Context, projectPath string, secretsIntegration SecretsIntegration, recipients []string, dryRun bool) error
    
    // DecryptVariables decrypts age-encrypted values in variables for debugging
    DecryptVariables(ctx context.Context, variables map[string]*spookytypes.Variable, secretsIntegration SecretsIntegration, identityPath string) error
}
```

### TemplatesIntegration Interface

The `TemplatesIntegration` interface provides template management:

```go
type TemplatesIntegration interface {
    // LoadTemplate loads a template from the given path
    LoadTemplate(ctx context.Context, templatePath string) (*spookytypes.Template, error)
    
    // RenderTemplate renders a template with the given data
    RenderTemplate(ctx context.Context, template *spookytypes.Template, data map[string]interface{}) (string, error)
    
    // ValidateTemplate validates a template
    ValidateTemplate(ctx context.Context, template *spookytypes.Template) (*spookytypes.ValidationResult, error)
}
```

### MachinesIntegration Interface

The `MachinesIntegration` interface provides machine management:

```go
type MachinesIntegration interface {
    // LoadMachines loads machines from the given source
    LoadMachines(ctx context.Context, source string) ([]spookytypes.Machine, error)
    
    // SaveMachines saves machines to the given destination
    SaveMachines(ctx context.Context, machines []spookytypes.Machine, destination string) error
    
    // EncryptMachines encrypts all machine secrets that have encrypted=true
    EncryptMachines(ctx context.Context, projectPath string, secretsIntegration SecretsIntegration, recipients []string, dryRun bool) error
    
    // ValidateMachines validates machines
    ValidateMachines(ctx context.Context, machines []spookytypes.Machine) (*spookytypes.ValidationResult, error)
    
    // PingMachines pings machines to check connectivity
    PingMachines(ctx context.Context, machines []spookytypes.Machine) ([]spookytypes.MachineStatus, error)
    
    // ExportMachines exports machines to HCL format according to machines schema
    ExportMachines(ctx context.Context, machines []spookytypes.Machine, outputPath string) error
    
    // GetMachineByName looks up a machine by hostname
    GetMachineByName(ctx context.Context, name string) (*spookytypes.Machine, error)
    
    // GetMachinesByTags filters machines by tags (supports key=value and key-only matching)
    GetMachinesByTags(ctx context.Context, tags []string) ([]spookytypes.Machine, error)
    
    // GetFullInventory returns the complete machine inventory
    GetFullInventory(ctx context.Context) ([]spookytypes.Machine, error)
    
    // GetMachinesByFilter applies complex filtering criteria to machines
    GetMachinesByFilter(ctx context.Context, filter interface{}) ([]spookytypes.Machine, error)
    
    // DecryptMachines decrypts age-encrypted values in machines for debugging
    DecryptMachines(ctx context.Context, machines []spookytypes.Machine, secretsIntegration SecretsIntegration, identityPath string) error
}
```

### SecretsIntegration Interface

The `SecretsIntegration` interface provides secrets management with age encryption:

```go
type SecretsIntegration interface {
    // Age-specific methods
    EncryptWithAge(ctx context.Context, data []byte, recipients []string) ([]byte, error)
    DecryptWithAge(ctx context.Context, data []byte, identityPath string) ([]byte, error)
    EncryptWithPassphrase(ctx context.Context, data []byte, passphrase string) ([]byte, error)
    DecryptWithPassphrase(ctx context.Context, data []byte, passphrase string) ([]byte, error)
    
    // Key management
    ValidateAgeKey(ctx context.Context, keyPath string) error
    ListRecipients(ctx context.Context, encryptedData []byte) ([]string, error)
    
    // Application-level validation
    ValidateAgeEncryptedValue(ctx context.Context, value string) error
    LoadRecipients(ctx context.Context, recipientsPath string) ([]string, error)
    LoadIdentities(ctx context.Context, identitiesPath string) ([]string, error)
    
    // HCL encryption and decryption methods
    EncryptHCLValues(ctx context.Context, data interface{}, recipients []string, dryRun bool) error
    EncryptHCLValuesSensitive(ctx context.Context, data interface{}, recipients []string, dryRun bool, shouldEncrypt func(path []string, value interface{}) bool) error
    EncryptHCLValuesWithJSONSupport(ctx context.Context, data interface{}, recipients []string, dryRun bool) error
    DecryptHCLValues(ctx context.Context, data interface{}, identityPath string) error
    DecryptHCLValuesWithJSONSupport(ctx context.Context, data interface{}, identityPath string) error
}
```

### SSHManager Interface

The `SSHManager` interface manages SSH connections, authentication, and acting:

```go
type SSHManager interface {
    // CreateClient creates a new SSH client with the given configuration
    CreateClient(ctx context.Context, config *spookytypes.ClientConfig) (*spookytypes.Client, error)
    
    // Connect establishes an SSH connection to the given host
    Connect(ctx context.Context, request *spookytypes.ConnectionRequest) (*spookytypes.ConnectionResult, error)
    
    // Authenticate authenticates with the given credentials
    Authenticate(ctx context.Context, connection *spookytypes.Connection, auth *spookytypes.Authentication) (*spookytypes.AuthenticationResult, error)
    
    // CreateSession creates a new SSH session
    CreateSession(ctx context.Context, connection *spookytypes.Connection) (*spookytypes.Session, error)
    
    // RunCommand runs a command via SSH
    RunCommand(ctx context.Context, session *spookytypes.Session, command *spookytypes.SSHCommand) (*spookytypes.SSHCommandResult, error)
    
    // CreateActingSession creates a new SSH acting session
    CreateActingSession(ctx context.Context, connection *spookytypes.Connection) (*spookytypes.ActingSession, error)
    
    // TransferFile transfers a file via SSH
    TransferFile(ctx context.Context, session *spookytypes.Session, transfer *spookytypes.FileTransfer) (*spookytypes.FileTransferResult, error)
    
    // ValidateConnection validates SSH connection parameters
    ValidateConnection(ctx context.Context, request *spookytypes.ConnectionRequest) (*spookytypes.ValidationResult, error)
    
    // ValidateAuthentication validates SSH authentication parameters
    ValidateAuthentication(ctx context.Context, auth *spookytypes.Authentication) (*spookytypes.ValidationResult, error)
    
    // GetConnectionPool returns the connection pool
    GetConnectionPool() *spookytypes.ConnectionPool
    
    // Close closes all SSH connections
    Close(ctx context.Context) error
}
```

### Validation Interfaces

#### ActionValidator Interface

```go
type ActionValidator interface {
    // ValidateActions validates a collection of actions
    ValidateActions(ctx context.Context, actions []spookytypes.Action) (*spookytypes.ValidationResult, error)
    
    // ValidateAction validates a single action
    ValidateAction(ctx context.Context, action *spookytypes.Action) (*spookytypes.ValidationResult, error)
}
```

#### VariableValidator Interface

```go
type VariableValidator interface {
    // ValidateVariables validates a collection of variables
    ValidateVariables(ctx context.Context, variables map[string]*spookytypes.Variable) (*spookytypes.ValidationResult, error)
    
    // ValidateVariable validates a single variable
    ValidateVariable(ctx context.Context, variable *spookytypes.Variable) (*spookytypes.ValidationResult, error)
}
```

#### VariableLoader Interface

```go
type VariableLoader interface {
    // LoadVariablesFromFile loads variables from a file
    LoadVariablesFromFile(ctx context.Context, filePath string) (map[string]*spookytypes.Variable, error)
    
    // LoadVariablesFromDirectory loads variables from a directory
    LoadVariablesFromDirectory(ctx context.Context, dirPath string) (map[string]*spookytypes.Variable, error)
}
```

#### MachineValidator Interface

```go
type MachineValidator interface {
    // ValidateMachines validates machines
    ValidateMachines(ctx context.Context, machines []spookytypes.Machine) (*spookytypes.ValidationResult, error)
    
    // ValidateMachine validates a single machine
    ValidateMachine(ctx context.Context, machine *spookytypes.Machine) (*spookytypes.ValidationResult, error)
}
```

#### MachineLoader Interface

```go
type MachineLoader interface {
    // LoadMachinesFromFile loads machines from a file
    LoadMachinesFromFile(ctx context.Context, filePath string) ([]spookytypes.Machine, error)
    
    // LoadMachinesFromDirectory loads machines from a directory
    LoadMachinesFromDirectory(ctx context.Context, dirPath string) ([]spookytypes.Machine, error)
}
```

### ConfigIntegration Interface

```go
type ConfigIntegration interface {
    // LoadConfig loads configuration from the given source
    LoadConfig(ctx context.Context, source string) (*spookytypes.Config, error)
    
    // ValidateConfig validates configuration
    ValidateConfig(ctx context.Context, config *spookytypes.Config) (*spookytypes.ValidationResult, error)
    
    // SaveConfig saves configuration to the given destination
    SaveConfig(ctx context.Context, config *spookytypes.Config, destination string) error
}
```

### FactStorage Interface

```go
type FactStorage interface {
    // GetStats returns memory usage statistics for debugging
    GetStats() (map[string]interface{}, error)
}
```

## Core Types

### ValidationResult

```go
type ValidationResult struct {
    Valid    bool                   `json:"valid"`
    Errors   []SchemaError          `json:"errors,omitempty"`
    Warnings []SchemaError          `json:"warnings,omitempty"`
    Details  map[string]interface{} `json:"details,omitempty"`
}

type SchemaError struct {
    Field   string `json:"field,omitempty"`
    Message string `json:"message"`
    Value   string `json:"value,omitempty"`
}
```

## Implementation Details

### Current Implementation Status

The interfaces system currently has:

1. **Complete Interface Definitions**: All core interfaces are defined and implemented
2. **Integration Implementations**: All integration interfaces have concrete implementations
3. **Validation Support**: Comprehensive validation interfaces and implementations
4. **Error Handling**: Structured error handling through validation results
5. **Context Support**: Context-aware operations throughout the system

### Key Features

1. **Interface-Based Architecture**: All system components use interfaces for loose coupling
2. **Comprehensive Validation**: All components support validation with detailed error reporting
3. **Age Encryption Integration**: Built-in support for age encryption in secrets management
4. **SSH Management**: Complete SSH connection and acting management
5. **Project Management**: Full project lifecycle management
6. **Configuration Management**: Flexible configuration loading and validation

## Usage Examples

### Basic Integration Usage

```go
// Create integration manager
manager := NewIntegrationManager()

// Get facts integration
factsIntegration := manager.GetFactsIntegration()

// Collect facts
machine := &spookytypes.Machine{
    Hostname: "web-server",
    Host:     "192.168.1.100",
    Port:     22,
    User:     "admin",
}
facts, err := factsIntegration.CollectFacts(ctx, machine)
if err != nil {
    return fmt.Errorf("failed to collect facts: %w", err)
}
```

### Project Management Usage

```go
// Create project manager
projectManager := NewProjectManager()

// Initialize new project
project, err := projectManager.Initialize(ctx, "./my-project")
if err != nil {
    return fmt.Errorf("failed to initialize project: %w", err)
}

// Validate project
result, err := projectManager.Validate(ctx, project)
if err != nil {
    return fmt.Errorf("failed to validate project: %w", err)
}

if !result.Valid {
    for _, error := range result.Errors {
        log.Printf("Validation error: %s", error.Message)
    }
    return fmt.Errorf("project validation failed")
}
```

### Secrets Management Usage

```go
// Create secrets integration
secretsIntegration := NewSecretsIntegration(logger, config)

// Encrypt data with age
data := []byte("sensitive data")
recipients := []string{"age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"}
encrypted, err := secretsIntegration.EncryptWithAge(ctx, data, recipients)
if err != nil {
    return fmt.Errorf("failed to encrypt data: %w", err)
}

// Decrypt data
decrypted, err := secretsIntegration.DecryptWithAge(ctx, encrypted, "~/.age/identity.txt")
if err != nil {
    return fmt.Errorf("failed to decrypt data: %w", err)
}
```

### SSH Management Usage

```go
// Create SSH manager
sshManager := NewSSHManager(logger)

// Create connection request
request := &spookytypes.ConnectionRequest{
    Host:     "192.168.1.100",
    Port:     22,
    User:     "admin",
    KeyPath:  "~/.ssh/id_rsa",
    Timeout:  30 * time.Second,
}

// Establish connection
connectionResult, err := sshManager.Connect(ctx, request)
if err != nil {
    return fmt.Errorf("failed to connect: %w", err)
}

// Create session
session, err := sshManager.CreateSession(ctx, connectionResult.Connection)
if err != nil {
    return fmt.Errorf("failed to create session: %w", err)
}

// Run command
command := &spookytypes.SSHCommand{
    Command: "ls -la",
    Timeout: 30 * time.Second,
}
result, err := sshManager.RunCommand(ctx, session, command)
if err != nil {
    return fmt.Errorf("failed to run command: %w", err)
}
```

## Error Handling

### Validation Error Handling

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

### Integration Error Handling

```go
// Handle integration errors gracefully
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
        MachineID: "test-server",
        Facts:     map[string]interface{}{"os": "linux"},
    }
    mockIntegration.facts["test-server"] = testFacts
    
    // Test facts collection
    machine := &spookytypes.Machine{Hostname: "test-server"}
    facts, err := mockIntegration.CollectFacts(ctx, machine)
    if err != nil {
        t.Fatalf("Failed to collect facts: %v", err)
    }
    
    if facts == nil {
        t.Error("Expected facts, got nil")
    }
}
```

### Mock Implementation

```go
type MockFactsIntegration struct {
    facts map[string]*spookytypes.FactCollection
    mutex sync.RWMutex
}

func (m *MockFactsIntegration) CollectFacts(ctx context.Context, machine *spookytypes.Machine) (interface{}, error) {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    
    if facts, exists := m.facts[machine.Hostname]; exists {
        return facts, nil
    }
    
    return nil, fmt.Errorf("facts not found for machine: %s", machine.Hostname)
}

func (m *MockFactsIntegration) StoreFacts(ctx context.Context, facts interface{}) error {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    if factCollection, ok := facts.(*spookytypes.FactCollection); ok {
        m.facts[factCollection.MachineID] = factCollection
        return nil
    }
    
    return fmt.Errorf("invalid facts type")
}

// Implement other interface methods...
func (m *MockFactsIntegration) LoadFacts(ctx context.Context) (interface{}, error) {
    return nil, nil
}

func (m *MockFactsIntegration) ValidateFacts(ctx context.Context, facts interface{}) (*spookytypes.ValidationResult, error) {
    if facts == nil {
        return &spookytypes.ValidationResult{
            Valid:    false,
            Errors:   []spookytypesschemas.SchemaError{{Message: "facts cannot be nil"}},
            Warnings: []spookytypesschemas.SchemaError{},
        }, nil
    }
    return &spookytypes.ValidationResult{Valid: true}, nil
}

func (m *MockFactsIntegration) DecryptFacts(ctx context.Context, facts interface{}, secretsIntegration SecretsIntegration, identityPath string) error {
    return nil
}

func (m *MockFactsIntegration) GetManager() interface{} {
    return nil
}
```

## Best Practices

### Interface Design

1. **Single Responsibility**: Each interface should have a single responsibility
2. **Consistent Naming**: Use consistent naming conventions across interfaces
3. **Error Handling**: Define clear error handling patterns through validation results
4. **Context Usage**: Use context for cancellation and timeouts
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
    return m.factsIntegration.CollectFacts(ctx, machine)
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

1. **Plugin System**: Pluggable interface system for extensibility
2. **Service Discovery**: Automatic service discovery for integrations
3. **Load Balancing**: Load balancing for high-availability integrations
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
