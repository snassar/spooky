# Spooky Interfaces API Reference

This document provides a comprehensive reference for all interfaces in the spooky codebase. It covers interface contracts, usage examples, implementation patterns, and best practices for developers working with the spooky system.

## Table of Contents

1. [Overview](#overview)
2. [Core Interfaces](#core-interfaces)
3. [Integration Interfaces](#integration-interfaces)
4. [Management Interfaces](#management-interfaces)
5. [Validation Interfaces](#validation-interfaces)
6. [Storage Interfaces](#storage-interfaces)
7. [Usage Examples](#usage-examples)
8. [Implementation Patterns](#implementation-patterns)
9. [Best Practices](#best-practices)
10. [Testing Interfaces](#testing-interfaces)

## Overview

The spooky codebase follows an interface-first architectural approach where all system components are defined through well-structured interfaces. This enables loose coupling, testability, and extensibility throughout the system.

### Interface Design Principles

- **Interface-First Design**: Define interfaces before implementations
- **Dependency Injection**: Use interfaces for all dependencies
- **Loose Coupling**: Components depend on interfaces, not concrete types
- **Comprehensive Coverage**: All public APIs are defined through interfaces
- **Clear Contracts**: Each interface has a single, well-defined responsibility

## Core Interfaces

### IntegrationManager

The `IntegrationManager` serves as the central coordinator for all system integrations, providing access to specialized integration components.

```go
type IntegrationManager interface {
    GetFactsIntegration() FactsIntegration
    GetActionsIntegration() ActionsIntegration
    GetVariablesIntegration() VariablesIntegration
    GetTemplatesIntegration() TemplatesIntegration
    GetMachinesIntegration() MachinesIntegration
    GetSecretsIntegration() SecretsIntegration
    GetConfigIntegration() ConfigIntegration
}
```

**Usage Example:**
```go
// Create integration manager
manager := NewIntegrationManager(config)

// Access facts integration
factsIntegration := manager.GetFactsIntegration()
facts, err := factsIntegration.CollectFacts(ctx, machine)

// Access actions integration
actionsIntegration := manager.GetActionsIntegration()
results, err := actionsIntegration.RunActions(ctx, actions, machines)

// Access machines integration
machinesIntegration := manager.GetMachinesIntegration()
inventory, err := machinesIntegration.GetFullInventory(ctx)
```

**Implementation Pattern:**
```go
type Manager struct {
    factsIntegration     FactsIntegration
    actionsIntegration   ActionsIntegration
    variablesIntegration VariablesIntegration
    templatesIntegration TemplatesIntegration
    machinesIntegration  MachinesIntegration
    secretsIntegration   SecretsIntegration
    configIntegration    ConfigIntegration
}

func (m *Manager) GetFactsIntegration() FactsIntegration {
    return m.factsIntegration
}

func (m *Manager) GetActionsIntegration() ActionsIntegration {
    return m.actionsIntegration
}

// ... other getter methods
```

### ProjectManager

The `ProjectManager` handles project lifecycle operations including initialization, loading, validation, and persistence.

```go
type ProjectManager interface {
    Initialize(ctx context.Context, projectPath string) (*spookytypes.Project, error)
    Load(ctx context.Context, projectPath string) (*spookytypes.Project, error)
    Validate(ctx context.Context, project *spookytypes.Project) (*spookytypes.ValidationResult, error)
    Save(ctx context.Context, project *spookytypes.Project) error
    Delete(ctx context.Context, projectPath string) error
}
```

**Usage Example:**
```go
// Initialize a new project
projectManager := NewProjectManager()
project, err := projectManager.Initialize(ctx, "./my-project")
if err != nil {
    return fmt.Errorf("failed to initialize project: %w", err)
}

// Load existing project
project, err = projectManager.Load(ctx, "./existing-project")
if err != nil {
    return fmt.Errorf("failed to load project: %w", err)
}

// Validate project structure
result, err := projectManager.Validate(ctx, project)
if err != nil {
    return fmt.Errorf("validation failed: %w", err)
}

if !result.IsValid() {
    return fmt.Errorf("project validation failed: %s", result.GetErrors())
}
```

### ConfigManager

The `ConfigManager` handles configuration loading, validation, and persistence operations.

```go
type ConfigManager interface {
    Load(ctx context.Context, configPath string) (*spookytypes.Config, error)
    Validate(ctx context.Context, config *spookytypes.Config) (*spookytypes.ValidationResult, error)
    Save(ctx context.Context, config *spookytypes.Config, configPath string) error
    GetDefaultConfig() *spookytypes.Config
}
```

**Usage Example:**
```go
// Load configuration
configManager := NewConfigManager()
config, err := configManager.Load(ctx, "./config.hcl")
if err != nil {
    return fmt.Errorf("failed to load config: %w", err)
}

// Validate configuration
result, err := configManager.Validate(ctx, config)
if err != nil {
    return fmt.Errorf("config validation failed: %w", err)
}

// Save configuration
err = configManager.Save(ctx, config, "./config.hcl")
if err != nil {
    return fmt.Errorf("failed to save config: %w", err)
}

// Get default configuration
defaultConfig := configManager.GetDefaultConfig()
```

### LogManager

The `LogManager` provides centralized logging operations and logger instance management.

```go
type LogManager interface {
    GetLogger(component string) spookytypes.Logger
    SetLevel(level spookytypes.LogLevel)
    GetLevel() spookytypes.LogLevel
    Configure(config *spookytypes.LogConfig) error
    Flush() error
    Close() error
}
```

**Usage Example:**
```go
// Create log manager
logManager := NewLogManager()

// Configure logging
config := &spookytypes.LogConfig{
    Level:  spookytypes.LogLevelInfo,
    Format: "json",
    Output: "stdout",
}
err := logManager.Configure(config)
if err != nil {
    return fmt.Errorf("failed to configure logging: %w", err)
}

// Get logger for component
logger := logManager.GetLogger("facts")
logger.Info("Starting fact collection", "server", "web-server")

// Set log level
logManager.SetLevel(spookytypes.LogLevelDebug)

// Flush and close
logManager.Flush()
logManager.Close()
```

## Integration Interfaces

### FactsIntegration

The `FactsIntegration` provides fact collection, storage, and validation capabilities.

```go
type FactsIntegration interface {
    CollectFacts(ctx context.Context, machine *spookytypes.Machine) (interface{}, error)
    StoreFacts(ctx context.Context, facts interface{}) error
    LoadFacts(ctx context.Context) (interface{}, error)
    ValidateFacts(ctx context.Context, facts interface{}) (*spookytypes.ValidationResult, error)
    GetManager() interface{}
}
```

**Usage Example:**
```go
// Get facts integration
factsIntegration := manager.GetFactsIntegration()

// Collect facts from machine
machine := &spookytypes.Machine{
    Hostname: "web-server",
    Port:     22,
    User:     "admin",
}
facts, err := factsIntegration.CollectFacts(ctx, machine)
if err != nil {
    return fmt.Errorf("failed to collect facts: %w", err)
}

// Store facts
err = factsIntegration.StoreFacts(ctx, facts)
if err != nil {
    return fmt.Errorf("failed to store facts: %w", err)
}

// Load facts
loadedFacts, err := factsIntegration.LoadFacts(ctx)
if err != nil {
    return fmt.Errorf("failed to load facts: %w", err)
}

// Validate facts
result, err := factsIntegration.ValidateFacts(ctx, loadedFacts)
if err != nil {
    return fmt.Errorf("fact validation failed: %w", err)
}
```

### ActionsIntegration

The `ActionsIntegration` provides action management and orchestration capabilities.

```go
type ActionsIntegration interface {
    LoadActions(ctx context.Context, source string) ([]spookytypes.Action, error)
    ValidateActions(ctx context.Context, actions []spookytypes.Action) (*spookytypes.ValidationResult, error)
    RunActions(ctx context.Context, actions []spookytypes.Action, machines []spookytypes.Machine) ([]spookytypes.ActingResult, error)
    GetSSHManager() SSHManager
}
```

**Usage Example:**
```go
// Get actions integration
actionsIntegration := manager.GetActionsIntegration()

// Load actions from file
actions, err := actionsIntegration.LoadActions(ctx, "./actions.hcl")
if err != nil {
    return fmt.Errorf("failed to load actions: %w", err)
}

// Validate actions
result, err := actionsIntegration.ValidateActions(ctx, actions)
if err != nil {
    return fmt.Errorf("action validation failed: %w", err)
}

// Run actions on machines
machines := []spookytypes.Machine{
    {Hostname: "web-server", Port: 22, User: "admin"},
    {Hostname: "db-server", Port: 22, User: "admin"},
}
results, err := actionsIntegration.RunActions(ctx, actions, machines)
if err != nil {
    return fmt.Errorf("failed to run actions: %w", err)
}

// Process results
for i, result := range results {
    if result.Error != nil {
        log.Printf("Action %d failed: %v", i, result.Error)
    } else {
        log.Printf("Action %d completed successfully", i)
    }
}
```

### VariablesIntegration

The `VariablesIntegration` provides variable loading, resolution, and validation capabilities.

```go
type VariablesIntegration interface {
    LoadVariables(ctx context.Context, source string) (map[string]*spookytypes.Variable, error)
    ResolveVariables(ctx context.Context, variables map[string]*spookytypes.Variable, context *spookytypes.VariableContext) (*spookytypes.VariableResolutionResult, error)
    ValidateVariables(ctx context.Context, variables map[string]*spookytypes.Variable) (*spookytypes.ValidationResult, error)
}
```

**Usage Example:**
```go
// Get variables integration
variablesIntegration := manager.GetVariablesIntegration()

// Load variables from directory
variables, err := variablesIntegration.LoadVariables(ctx, "./variables")
if err != nil {
    return fmt.Errorf("failed to load variables: %w", err)
}

// Create variable context
context := &spookytypes.VariableContext{
    Environment: "production",
    Machine:     "web-server",
    Project:     "my-project",
}

// Resolve variables
result, err := variablesIntegration.ResolveVariables(ctx, variables, context)
if err != nil {
    return fmt.Errorf("failed to resolve variables: %w", err)
}

// Use resolved variables
for name, value := range result.ResolvedVariables {
    log.Printf("Variable %s = %v", name, value)
}

// Validate variables
validationResult, err := variablesIntegration.ValidateVariables(ctx, variables)
if err != nil {
    return fmt.Errorf("variable validation failed: %w", err)
}
```

### TemplatesIntegration

The `TemplatesIntegration` provides template loading, rendering, and validation capabilities.

```go
type TemplatesIntegration interface {
    LoadTemplate(ctx context.Context, templatePath string) (*spookytypes.Template, error)
    RenderTemplate(ctx context.Context, template *spookytypes.Template, data map[string]interface{}) (string, error)
    ValidateTemplate(ctx context.Context, template *spookytypes.Template) (*spookytypes.ValidationResult, error)
}
```

**Usage Example:**
```go
// Get templates integration
templatesIntegration := manager.GetTemplatesIntegration()

// Load template
template, err := templatesIntegration.LoadTemplate(ctx, "./templates/deploy.sh.tmpl")
if err != nil {
    return fmt.Errorf("failed to load template: %w", err)
}

// Prepare template data
data := map[string]interface{}{
    "app_name":    "my-app",
    "version":     "1.0.0",
    "environment": "production",
    "servers":     []string{"web-1", "web-2"},
}

// Render template
rendered, err := templatesIntegration.RenderTemplate(ctx, template, data)
if err != nil {
    return fmt.Errorf("failed to render template: %w", err)
}

// Validate template
result, err := templatesIntegration.ValidateTemplate(ctx, template)
if err != nil {
    return fmt.Errorf("template validation failed: %w", err)
}

// Use rendered content
fmt.Println(rendered)
```

### MachinesIntegration

The `MachinesIntegration` provides machine inventory management and connectivity testing capabilities.

```go
type MachinesIntegration interface {
    LoadMachines(ctx context.Context, source string) ([]spookytypes.Machine, error)
    ValidateMachines(ctx context.Context, machines []spookytypes.Machine) (*spookytypes.ValidationResult, error)
    PingMachines(ctx context.Context, machines []spookytypes.Machine) ([]spookytypes.MachineStatus, error)
    ExportMachines(ctx context.Context, machines []spookytypes.Machine, outputPath string) error
    GetMachineByName(ctx context.Context, name string) (*spookytypes.Machine, error)
    GetMachinesByTags(ctx context.Context, tags []string) ([]spookytypes.Machine, error)
    GetFullInventory(ctx context.Context) ([]spookytypes.Machine, error)
    GetMachinesByFilter(ctx context.Context, filter interface{}) ([]spookytypes.Machine, error)
}
```

**Usage Example:**
```go
// Get machines integration
machinesIntegration := manager.GetMachinesIntegration()

// Load machines from file
machines, err := machinesIntegration.LoadMachines(ctx, "./machines.hcl")
if err != nil {
    return fmt.Errorf("failed to load machines: %w", err)
}

// Validate machines
result, err := machinesIntegration.ValidateMachines(ctx, machines)
if err != nil {
    return fmt.Errorf("machine validation failed: %w", err)
}

// Ping machines to check connectivity
statuses, err := machinesIntegration.PingMachines(ctx, machines)
if err != nil {
    return fmt.Errorf("failed to ping machines: %w", err)
}

// Check machine statuses
for i, status := range statuses {
    if status.Online {
        log.Printf("Machine %s is online", machines[i].Hostname)
    } else {
        log.Printf("Machine %s is offline: %v", machines[i].Hostname, status.Error)
    }
}

// Get machine by name
machine, err := machinesIntegration.GetMachineByName(ctx, "web-server")
if err != nil {
    return fmt.Errorf("failed to get machine: %w", err)
}

// Get machines by tags
webServers, err := machinesIntegration.GetMachinesByTags(ctx, []string{"web", "production"})
if err != nil {
    return fmt.Errorf("failed to get machines by tags: %w", err)
}

// Export machines
err = machinesIntegration.ExportMachines(ctx, machines, "./exported-machines.hcl")
if err != nil {
    return fmt.Errorf("failed to export machines: %w", err)
}
```

### SecretsIntegration

The `SecretsIntegration` provides encryption, decryption, and key validation capabilities.

```go
type SecretsIntegration interface {
    Encrypt(ctx context.Context, data []byte, key []byte) ([]byte, error)
    Decrypt(ctx context.Context, data []byte, key []byte) ([]byte, error)
    ValidateKey(ctx context.Context, key []byte) error
}
```

**Usage Example:**
```go
// Get secrets integration
secretsIntegration := manager.GetSecretsIntegration()

// Validate encryption key
key := []byte("my-secret-key-32-bytes-long")
err := secretsIntegration.ValidateKey(ctx, key)
if err != nil {
    return fmt.Errorf("invalid key: %w", err)
}

// Encrypt sensitive data
sensitiveData := []byte("password123")
encrypted, err := secretsIntegration.Encrypt(ctx, sensitiveData, key)
if err != nil {
    return fmt.Errorf("failed to encrypt data: %w", err)
}

// Decrypt data
decrypted, err := secretsIntegration.Decrypt(ctx, encrypted, key)
if err != nil {
    return fmt.Errorf("failed to decrypt data: %w", err)
}

// Verify decryption
if string(decrypted) != string(sensitiveData) {
    return fmt.Errorf("decryption verification failed")
}
```

## Management Interfaces

### SSHManager

The `SSHManager` provides comprehensive SSH connection, authentication, and acting capabilities.

```go
type SSHManager interface {
    CreateClient(ctx context.Context, config *spookytypes.ClientConfig) (*spookytypes.Client, error)
    Connect(ctx context.Context, request *spookytypes.ConnectionRequest) (*spookytypes.ConnectionResult, error)
    Authenticate(ctx context.Context, connection *spookytypes.Connection, auth *spookytypes.Authentication) (*spookytypes.AuthenticationResult, error)
    CreateSession(ctx context.Context, connection *spookytypes.Connection) (*spookytypes.Session, error)
    RunCommand(ctx context.Context, session *spookytypes.Session, command *spookytypes.SSHCommand) (*spookytypes.SSHCommandResult, error)
    CreateActingSession(ctx context.Context, connection *spookytypes.Connection) (*spookytypes.ActingSession, error)
    TransferFile(ctx context.Context, session *spookytypes.Session, transfer *spookytypes.FileTransfer) (*spookytypes.FileTransferResult, error)
    ValidateConnection(ctx context.Context, request *spookytypes.ConnectionRequest) (*spookytypes.ValidationResult, error)
    ValidateAuthentication(ctx context.Context, auth *spookytypes.Authentication) (*spookytypes.ValidationResult, error)
    GetConnectionPool() *spookytypes.ConnectionPool
    Close(ctx context.Context) error
}
```

**Usage Example:**
```go
// Create SSH manager
sshManager := NewSSHManager()

// Create client configuration
clientConfig := &spookytypes.ClientConfig{
    Timeout: 30 * time.Second,
    KeepAlive: 60 * time.Second,
}

// Create SSH client
client, err := sshManager.CreateClient(ctx, clientConfig)
if err != nil {
    return fmt.Errorf("failed to create SSH client: %w", err)
}

// Create connection request
request := &spookytypes.ConnectionRequest{
    Host: "web-server",
    Port: 22,
    User: "admin",
}

// Validate connection parameters
result, err := sshManager.ValidateConnection(ctx, request)
if err != nil {
    return fmt.Errorf("connection validation failed: %w", err)
}

// Connect to server
connection, err := sshManager.Connect(ctx, request)
if err != nil {
    return fmt.Errorf("failed to connect: %w", err)
}

// Create authentication
auth := &spookytypes.Authentication{
    Method: "ssh_key",
    KeyPath: "~/.ssh/id_rsa",
}

// Validate authentication
authResult, err := sshManager.ValidateAuthentication(ctx, auth)
if err != nil {
    return fmt.Errorf("authentication validation failed: %w", err)
}

// Authenticate
authResult, err = sshManager.Authenticate(ctx, connection, auth)
if err != nil {
    return fmt.Errorf("authentication failed: %w", err)
}

// Create session
session, err := sshManager.CreateSession(ctx, connection)
if err != nil {
    return fmt.Errorf("failed to create session: %w", err)
}

// Run command
command := &spookytypes.SSHCommand{
    Command: "ls -la",
    Timeout: 10 * time.Second,
}

cmdResult, err := sshManager.RunCommand(ctx, session, command)
if err != nil {
    return fmt.Errorf("command execution failed: %w", err)
}

log.Printf("Command output: %s", cmdResult.Output)

// Transfer file
transfer := &spookytypes.FileTransfer{
    LocalPath:  "./local-file.txt",
    RemotePath: "/tmp/remote-file.txt",
    Mode:       0644,
}

transferResult, err := sshManager.TransferFile(ctx, session, transfer)
if err != nil {
    return fmt.Errorf("file transfer failed: %w", err)
}

// Close connections
sshManager.Close(ctx)
```

### SchemaManager

The `SchemaManager` provides schema loading, validation, and registration capabilities.

```go
type SchemaManager interface {
    LoadSchema(ctx context.Context, schemaPath string) (*spookytypes.Schema, error)
    LoadEmbeddedSchema(ctx context.Context, schemaName string) (*spookytypes.Schema, error)
    Validate(ctx context.Context, schema *spookytypes.Schema, data interface{}) (*spookytypes.ValidationResult, error)
    Register(ctx context.Context, schema *spookytypes.Schema) error
}
```

**Usage Example:**
```go
// Create schema manager
schemaManager := NewSchemaManager()

// Load schema from file
schema, err := schemaManager.LoadSchema(ctx, "./schemas/project.schema.hcl")
if err != nil {
    return fmt.Errorf("failed to load schema: %w", err)
}

// Load embedded schema
embeddedSchema, err := schemaManager.LoadEmbeddedSchema(ctx, "project-directory")
if err != nil {
    return fmt.Errorf("failed to load embedded schema: %w", err)
}

// Validate data against schema
data := map[string]interface{}{
    "name": "my-project",
    "description": "A test project",
}

result, err := schemaManager.Validate(ctx, schema, data)
if err != nil {
    return fmt.Errorf("validation failed: %w", err)
}

if !result.IsValid() {
    return fmt.Errorf("data validation failed: %s", result.GetErrors())
}

// Register new schema
newSchema := &spookytypes.Schema{
    Name: "custom-schema",
    Version: "1.0.0",
    // ... schema definition
}

err = schemaManager.Register(ctx, newSchema)
if err != nil {
    return fmt.Errorf("failed to register schema: %w", err)
}
```

### CLIManager

The `CLIManager` provides command-line interface management and command registration capabilities.

```go
type CLIManager interface {
    RegisterCommand(command spookytypes.Command) error
    RunCommand(ctx context.Context, commandName string, args []string) error
    GetCommand(commandName string) (spookytypes.Command, bool)
    ListCommands() []spookytypes.Command
    ShowHelp(commandName string) error
    ShowVersion() error
}
```

**Usage Example:**
```go
// Create CLI manager
cliManager := NewCLIManager()

// Register commands
commands := []spookytypes.Command{
    NewFactsCommand(),
    NewActionsCommand(),
    NewMachinesCommand(),
    NewProjectCommand(),
}

for _, cmd := range commands {
    err := cliManager.RegisterCommand(cmd)
    if err != nil {
        return fmt.Errorf("failed to register command %s: %w", cmd.Name(), err)
    }
}

// List all commands
allCommands := cliManager.ListCommands()
for _, cmd := range allCommands {
    log.Printf("Registered command: %s", cmd.Name())
}

// Get specific command
if cmd, exists := cliManager.GetCommand("facts"); exists {
    log.Printf("Found facts command: %s", cmd.Description())
}

// Run command
err := cliManager.RunCommand(ctx, "facts", []string{"gather", "./project"})
if err != nil {
    return fmt.Errorf("failed to run facts command: %w", err)
}

// Show help
err = cliManager.ShowHelp("facts")
if err != nil {
    return fmt.Errorf("failed to show help: %w", err)
}

// Show version
err = cliManager.ShowVersion()
if err != nil {
    return fmt.Errorf("failed to show version: %w", err)
}
```

## Validation Interfaces

### ProjectValidator

The `ProjectValidator` provides project structure and configuration validation.

```go
type ProjectValidator interface {
    ValidateProject(ctx context.Context, project *spookytypes.Project) (*spookytypes.ValidationResult, error)
    ValidateProjectDirectory(ctx context.Context, projectPath string) (*spookytypes.ValidationResult, error)
    ValidateProjectConfig(ctx context.Context, config *spookytypes.ProjectConfig) (*spookytypes.ValidationResult, error)
}
```

**Usage Example:**
```go
// Create project validator
validator := NewProjectValidator()

// Validate project structure
project := &spookytypes.Project{
    Name: "my-project",
    Path: "./my-project",
}

result, err := validator.ValidateProject(ctx, project)
if err != nil {
    return fmt.Errorf("project validation failed: %w", err)
}

if !result.IsValid() {
    for _, error := range result.GetErrors() {
        log.Printf("Validation error: %s", error)
    }
}

// Validate project directory
result, err = validator.ValidateProjectDirectory(ctx, "./my-project")
if err != nil {
    return fmt.Errorf("directory validation failed: %w", err)
}

// Validate project configuration
config := &spookytypes.ProjectConfig{
    Name: "my-project",
    Description: "A test project",
}

result, err = validator.ValidateProjectConfig(ctx, config)
if err != nil {
    return fmt.Errorf("config validation failed: %w", err)
}
```

### ActionValidator

The `ActionValidator` provides action validation capabilities.

```go
type ActionValidator interface {
    ValidateActions(ctx context.Context, actions []spookytypes.Action) (*spookytypes.ValidationResult, error)
    ValidateAction(ctx context.Context, action spookytypes.Action) (*spookytypes.ValidationResult, error)
}
```

**Usage Example:**
```go
// Create action validator
validator := NewActionValidator()

// Validate single action
action := spookytypes.Action{
    Name: "deploy",
    Description: "Deploy application",
    Command: "deploy.sh",
}

result, err := validator.ValidateAction(ctx, action)
if err != nil {
    return fmt.Errorf("action validation failed: %w", err)
}

// Validate multiple actions
actions := []spookytypes.Action{
    {Name: "deploy", Command: "deploy.sh"},
    {Name: "restart", Command: "restart.sh"},
}

result, err = validator.ValidateActions(ctx, actions)
if err != nil {
    return fmt.Errorf("actions validation failed: %w", err)
}

if !result.IsValid() {
    for _, error := range result.GetErrors() {
        log.Printf("Action validation error: %s", error)
    }
}
```

### MachineValidator

The `MachineValidator` provides machine validation capabilities.

```go
type MachineValidator interface {
    ValidateMachines(ctx context.Context, machines []spookytypes.Machine) (*spookytypes.ValidationResult, error)
    ValidateMachine(ctx context.Context, machine spookytypes.Machine) (*spookytypes.ValidationResult, error)
}
```

**Usage Example:**
```go
// Create machine validator
validator := NewMachineValidator()

// Validate single machine
machine := spookytypes.Machine{
    Hostname: "web-server",
    Port:     22,
    User:     "admin",
}

result, err := validator.ValidateMachine(ctx, machine)
if err != nil {
    return fmt.Errorf("machine validation failed: %w", err)
}

// Validate multiple machines
machines := []spookytypes.Machine{
    {Hostname: "web-1", Port: 22, User: "admin"},
    {Hostname: "web-2", Port: 22, User: "admin"},
}

result, err = validator.ValidateMachines(ctx, machines)
if err != nil {
    return fmt.Errorf("machines validation failed: %w", err)
}
```

### VariableValidator

The `VariableValidator` provides variable validation capabilities.

```go
type VariableValidator interface {
    ValidateVariables(ctx context.Context, variables map[string]*spookytypes.Variable) (*spookytypes.ValidationResult, error)
    ValidateVariable(ctx context.Context, variable *spookytypes.Variable) (*spookytypes.ValidationResult, error)
}
```

**Usage Example:**
```go
// Create variable validator
validator := NewVariableValidator()

// Validate single variable
variable := &spookytypes.Variable{
    Name:  "app_version",
    Value: "1.0.0",
    Type:  "string",
}

result, err := validator.ValidateVariable(ctx, variable)
if err != nil {
    return fmt.Errorf("variable validation failed: %w", err)
}

// Validate multiple variables
variables := map[string]*spookytypes.Variable{
    "app_version": {Name: "app_version", Value: "1.0.0", Type: "string"},
    "debug_mode":  {Name: "debug_mode", Value: true, Type: "bool"},
}

result, err = validator.ValidateVariables(ctx, variables)
if err != nil {
    return fmt.Errorf("variables validation failed: %w", err)
}
```

## Storage Interfaces

### FactStorage

The `FactStorage` provides minimal storage capabilities for debugging and statistics.

```go
type FactStorage interface {
    GetStats() (map[string]interface{}, error)
}
```

**Usage Example:**
```go
// Create fact storage
storage := NewFactStorage()

// Get storage statistics
stats, err := storage.GetStats()
if err != nil {
    return fmt.Errorf("failed to get storage stats: %w", err)
}

// Use statistics
for key, value := range stats {
    log.Printf("Storage stat %s: %v", key, value)
}
```

## Usage Examples

### Complete Integration Example

Here's a complete example showing how to use multiple interfaces together:

```go
func RunProjectOrchestration(ctx context.Context, projectPath string) error {
    // Create integration manager
    manager := NewIntegrationManager(config)
    
    // Load project
    projectManager := NewProjectManager()
    project, err := projectManager.Load(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load project: %w", err)
    }
    
    // Validate project
    result, err := projectManager.Validate(ctx, project)
    if err != nil {
        return fmt.Errorf("project validation failed: %w", err)
    }
    
    if !result.IsValid() {
        return fmt.Errorf("project validation failed: %s", result.GetErrors())
    }
    
    // Get integrations
    factsIntegration := manager.GetFactsIntegration()
    actionsIntegration := manager.GetActionsIntegration()
    machinesIntegration := manager.GetMachinesIntegration()
    
    // Load machines
    machines, err := machinesIntegration.LoadMachines(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load machines: %w", err)
    }
    
    // Collect facts from all machines
    for _, machine := range machines {
        facts, err := factsIntegration.CollectFacts(ctx, &machine)
        if err != nil {
            log.Printf("Failed to collect facts from %s: %v", machine.Hostname, err)
            continue
        }
        
        err = factsIntegration.StoreFacts(ctx, facts)
        if err != nil {
            log.Printf("Failed to store facts for %s: %v", machine.Hostname, err)
        }
    }
    
    // Load and run actions
    actions, err := actionsIntegration.LoadActions(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load actions: %w", err)
    }
    
    results, err := actionsIntegration.RunActions(ctx, actions, machines)
    if err != nil {
        return fmt.Errorf("failed to run actions: %w", err)
    }
    
    // Process results
    for i, result := range results {
        if result.Error != nil {
            log.Printf("Action %d failed: %v", i, result.Error)
        } else {
            log.Printf("Action %d completed successfully", i)
        }
    }
    
    return nil
}
```

### Error Handling Example

Here's an example showing proper error handling with interfaces:

```go
func HandleIntegrationError(err error, operation string) error {
    // Check for specific error types
    switch {
    case errors.Is(err, ErrConnectionFailed):
        return fmt.Errorf("connection failed during %s: %w", operation, err)
        
    case errors.Is(err, ErrAuthenticationFailed):
        return fmt.Errorf("authentication failed during %s: %w", operation, err)
        
    case errors.Is(err, ErrValidationFailed):
        return fmt.Errorf("validation failed during %s: %w", operation, err)
        
    case errors.Is(err, ErrTimeout):
        return fmt.Errorf("operation %s timed out: %w", operation, err)
        
    default:
        return fmt.Errorf("unexpected error during %s: %w", operation, err)
    }
}

func SafeIntegrationOperation(ctx context.Context, operation func() error, operationName string) error {
    err := operation()
    if err != nil {
        return HandleIntegrationError(err, operationName)
    }
    return nil
}
```

## Implementation Patterns

### Dependency Injection Pattern

```go
// Good: Use dependency injection with interfaces
type Manager struct {
    factsIntegration     FactsIntegration
    actionsIntegration   ActionsIntegration
    machinesIntegration  MachinesIntegration
    logger              spookytypes.Logger
}

func NewManager(
    factsIntegration FactsIntegration,
    actionsIntegration ActionsIntegration,
    machinesIntegration MachinesIntegration,
    logger spookytypes.Logger,
) *Manager {
    return &Manager{
        factsIntegration:    factsIntegration,
        actionsIntegration:  actionsIntegration,
        machinesIntegration: machinesIntegration,
        logger:             logger,
    }
}

// Bad: Direct instantiation of concrete types
type BadManager struct {
    factsManager     *FactsManager
    actionsManager   *ActionsManager
    machinesManager  *MachinesManager
}
```

### Interface Composition Pattern

```go
// Good: Compose interfaces for complex operations
type OrchestrationManager interface {
    FactsIntegration
    ActionsIntegration
    MachinesIntegration
    ValidateOrchestration(ctx context.Context) (*spookytypes.ValidationResult, error)
}

type Manager struct {
    FactsIntegration
    ActionsIntegration
    MachinesIntegration
}

func (m *Manager) ValidateOrchestration(ctx context.Context) (*spookytypes.ValidationResult, error) {
    // Implementation
    return nil, nil
}
```

### Context Pattern

```go
// Good: Use context for state management
type OperationContext struct {
    ProjectPath string
    Environment string
    Variables   map[string]interface{}
    Metadata    map[string]interface{}
}

func (m *Manager) RunOperation(ctx context.Context, opCtx *OperationContext) error {
    // Use context for operation state
    return nil
}
```

## Best Practices

### 1. Interface Design

- **Single Responsibility**: Each interface should have a single, well-defined responsibility
- **Minimal Interface**: Keep interfaces small and focused
- **Clear Contracts**: Define clear contracts for all interface methods
- **Consistent Naming**: Use consistent naming conventions across interfaces

### 2. Implementation

- **Interface Compliance**: Ensure all implementations fully satisfy their interfaces
- **Error Handling**: Use structured error types and proper error wrapping
- **Context Usage**: Pass context through all interface methods
- **Resource Management**: Properly manage resources in implementations

### 3. Testing

- **Mock Interfaces**: Create mock implementations for testing
- **Interface Testing**: Test interface contracts and behavior
- **Integration Testing**: Test interface interactions
- **Error Testing**: Test error conditions and edge cases

### 4. Documentation

- **Clear Examples**: Provide clear, working examples for all interfaces
- **Error Scenarios**: Document error handling and recovery
- **Performance Notes**: Document performance characteristics
- **Migration Guides**: Document interface evolution and migration

## Testing Interfaces

### Mock Implementation Example

```go
// Create mock implementation for testing
type MockFactsIntegration struct {
    facts map[string]interface{}
    err   error
}

func (m *MockFactsIntegration) CollectFacts(ctx context.Context, machine *spookytypes.Machine) (interface{}, error) {
    if m.err != nil {
        return nil, m.err
    }
    return m.facts[machine.Hostname], nil
}

func (m *MockFactsIntegration) StoreFacts(ctx context.Context, facts interface{}) error {
    return m.err
}

func (m *MockFactsIntegration) LoadFacts(ctx context.Context) (interface{}, error) {
    if m.err != nil {
        return nil, m.err
    }
    return m.facts, nil
}

func (m *MockFactsIntegration) ValidateFacts(ctx context.Context, facts interface{}) (*spookytypes.ValidationResult, error) {
    if m.err != nil {
        return nil, m.err
    }
    return &spookytypes.ValidationResult{Valid: true}, nil
}

func (m *MockFactsIntegration) GetManager() interface{} {
    return nil
}

// Test with mock
func TestFactsIntegration(t *testing.T) {
    mock := &MockFactsIntegration{
        facts: map[string]interface{}{
            "test-server": map[string]interface{}{
                "os": "linux",
                "version": "ubuntu-20.04",
            },
        },
    }
    
    machine := &spookytypes.Machine{Hostname: "test-server"}
    facts, err := mock.CollectFacts(context.Background(), machine)
    
    assert.NoError(t, err)
    assert.NotNil(t, facts)
}
```

### Interface Contract Testing

```go
// Test interface contract compliance
func TestInterfaceCompliance(t *testing.T) {
    // Verify that implementations satisfy interfaces
    var _ FactsIntegration = &MockFactsIntegration{}
    var _ ActionsIntegration = &MockActionsIntegration{}
    var _ MachinesIntegration = &MockMachinesIntegration{}
}
```

This comprehensive interface documentation provides developers with clear guidance on how to use all interfaces in the spooky codebase, including practical examples, implementation patterns, and best practices.
