// Package interfaces provides the core system interfaces for the spooky codebase.
// These interfaces define the contracts and architectural patterns for all system components.
package interfaces

import (
	"context"

	spookytypes "spooky/internal/types"
)

// IntegrationManager coordinates all system integrations
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
}

// ProjectManager manages project lifecycle and operations
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

// ProjectValidator validates project structure and configuration
type ProjectValidator interface {
	// ValidateProject validates a project structure
	ValidateProject(ctx context.Context, project *spookytypes.Project) (*spookytypes.ValidationResult, error)

	// ValidateProjectDirectory validates project directory structure
	ValidateProjectDirectory(ctx context.Context, projectPath string) (*spookytypes.ValidationResult, error)

	// ValidateProjectConfig validates project configuration
	ValidateProjectConfig(ctx context.Context, config *spookytypes.ProjectConfig) (*spookytypes.ValidationResult, error)
}

// ProjectLoader loads project data from various sources
type ProjectLoader interface {
	// LoadProject loads a project from disk
	LoadProject(ctx context.Context, projectPath string) (*spookytypes.Project, error)

	// LoadProjectConfig loads project configuration
	LoadProjectConfig(ctx context.Context, projectPath string) (*spookytypes.ProjectConfig, error)

	// LoadProjectMetadata loads project metadata
	LoadProjectMetadata(ctx context.Context, projectPath string) (*spookytypes.ProjectMetadata, error)
}

// ConfigManager manages configuration loading and validation
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

// LogManager manages logging operations and logger instances
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

// SchemaManager manages schema loading and validation
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

// CLIManager manages command-line interface operations
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

// FactsIntegration provides fact collection and memory storage
type FactsIntegration interface {
	// CollectFacts collects facts from the given source
	CollectFacts(ctx context.Context, source string) (interface{}, error)

	// StoreFacts stores facts in memory
	StoreFacts(ctx context.Context, facts interface{}) error

	// LoadFacts loads facts from memory
	LoadFacts(ctx context.Context) (interface{}, error)

	// ValidateFacts validates facts
	ValidateFacts(ctx context.Context, facts interface{}) (*spookytypes.ValidationResult, error)

	// GetManager returns the underlying fact manager
	GetManager() interface{}
}

// ActionsIntegration provides action management and orchestration
type ActionsIntegration interface {
	// LoadActions loads actions from the given source
	LoadActions(ctx context.Context, source string) ([]spookytypes.Action, error)

	// ValidateActions validates actions
	ValidateActions(ctx context.Context, actions []spookytypes.Action) (*spookytypes.ValidationResult, error)

	// RunActions runs actions on the given machines
	RunActions(ctx context.Context, actions []spookytypes.Action, machines []spookytypes.Machine) ([]spookytypes.ActingResult, error)
}

// VariableValidator provides variable validation
type VariableValidator interface {
	// ValidateVariables validates a collection of variables
	ValidateVariables(ctx context.Context, variables map[string]*spookytypes.Variable) (*spookytypes.ValidationResult, error)

	// ValidateVariable validates a single variable
	ValidateVariable(ctx context.Context, variable *spookytypes.Variable) (*spookytypes.ValidationResult, error)
}

// VariableLoader provides variable loading
type VariableLoader interface {
	// LoadVariablesFromFile loads variables from a file
	LoadVariablesFromFile(ctx context.Context, filePath string) (map[string]*spookytypes.Variable, error)

	// LoadVariablesFromDirectory loads variables from a directory
	LoadVariablesFromDirectory(ctx context.Context, dirPath string) (map[string]*spookytypes.Variable, error)
}

// VariablesIntegration provides variable management
type VariablesIntegration interface {
	// LoadVariables loads variables from the given source
	LoadVariables(ctx context.Context, source string) (map[string]*spookytypes.Variable, error)

	// ResolveVariables resolves variables with the given context
	ResolveVariables(ctx context.Context, variables map[string]*spookytypes.Variable, context *spookytypes.VariableContext) (*spookytypes.VariableResolutionResult, error)

	// ValidateVariables validates variables
	ValidateVariables(ctx context.Context, variables map[string]*spookytypes.Variable) (*spookytypes.ValidationResult, error)
}

// TemplatesIntegration provides template management
type TemplatesIntegration interface {
	// LoadTemplate loads a template from the given path
	LoadTemplate(ctx context.Context, templatePath string) (*spookytypes.Template, error)

	// RenderTemplate renders a template with the given data
	RenderTemplate(ctx context.Context, template *spookytypes.Template, data map[string]interface{}) (string, error)

	// ValidateTemplate validates a template
	ValidateTemplate(ctx context.Context, template *spookytypes.Template) (*spookytypes.ValidationResult, error)
}

// MachineValidator provides machine validation
type MachineValidator interface {
	// ValidateMachines validates machines
	ValidateMachines(ctx context.Context, machines []spookytypes.Machine) (*spookytypes.ValidationResult, error)

	// ValidateMachine validates a single machine
	ValidateMachine(ctx context.Context, machine spookytypes.Machine) (*spookytypes.ValidationResult, error)
}

// MachineLoader provides machine loading
type MachineLoader interface {
	// LoadMachinesFromFile loads machines from a file
	LoadMachinesFromFile(ctx context.Context, filePath string) ([]spookytypes.Machine, error)

	// LoadMachinesFromDirectory loads machines from a directory
	LoadMachinesFromDirectory(ctx context.Context, dirPath string) ([]spookytypes.Machine, error)
}

// MachinesIntegration provides machine management
type MachinesIntegration interface {
	// LoadMachines loads machines from the given source
	LoadMachines(ctx context.Context, source string) ([]spookytypes.Machine, error)

	// ValidateMachines validates machines
	ValidateMachines(ctx context.Context, machines []spookytypes.Machine) (*spookytypes.ValidationResult, error)

	// PingMachines pings machines to check connectivity
	PingMachines(ctx context.Context, machines []spookytypes.Machine) ([]spookytypes.MachineStatus, error)
}

// SecretsIntegration provides secrets management
type SecretsIntegration interface {
	// Encrypt encrypts data with the given key
	Encrypt(ctx context.Context, data []byte, key []byte) ([]byte, error)

	// Decrypt decrypts data with the given key
	Decrypt(ctx context.Context, data []byte, key []byte) ([]byte, error)

	// GenerateKey generates a new encryption key
	GenerateKey(ctx context.Context) ([]byte, error)

	// ValidateKey validates an encryption key
	ValidateKey(ctx context.Context, key []byte) error
}

// ActionValidator provides action validation
type ActionValidator interface {
	// ValidateActions validates a collection of actions
	ValidateActions(ctx context.Context, actions []spookytypes.Action) (*spookytypes.ValidationResult, error)

	// ValidateAction validates a single action
	ValidateAction(ctx context.Context, action spookytypes.Action) (*spookytypes.ValidationResult, error)
}

// ConfigIntegration provides configuration management
type ConfigIntegration interface {
	// LoadConfig loads configuration from the given source
	LoadConfig(ctx context.Context, source string) (*spookytypes.Config, error)

	// ValidateConfig validates configuration
	ValidateConfig(ctx context.Context, config *spookytypes.Config) (*spookytypes.ValidationResult, error)

	// SaveConfig saves configuration to the given destination
	SaveConfig(ctx context.Context, config *spookytypes.Config, destination string) error
}

// SSHManager manages SSH connections, authentication, and acting
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

// FactStorage provides minimal storage for debugging and statistics during export operations
type FactStorage interface {
	// GetStats returns memory usage statistics for debugging
	GetStats() (map[string]interface{}, error)
}
