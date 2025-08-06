package ssh

import (
	"time"

	"spooky/internal/ssh/types"
)

// SSHManager defines the main interface for SSH operations
type SSHManager interface {
	// Core SSH operations
	Connect(host string, config *types.SSHConfig) (*types.SSHConnection, error)
	ExecuteCommand(connection *types.SSHConnection, command string) (*types.CommandResult, error)
	ExecuteScript(connection *types.SSHConnection, script string) (*types.CommandResult, error)
	CloseConnection(connection *types.SSHConnection) error

	// Connection pool operations
	GetConnection(host string) (*types.SSHConnection, error)
	ReturnConnection(connection *types.SSHConnection) error
	CloseAllConnections() error

	// Authentication operations
	Authenticate(connection *types.SSHConnection, auth *types.AuthenticationConfig) error
	ValidateAuthentication(auth *types.AuthenticationConfig) error

	// Acting operations
	ExecuteAction(connection *types.SSHConnection, action *types.SSHAction) (*types.ActionResult, error)
	ExecuteTemplate(connection *types.SSHConnection, template *types.TemplateAction) (*types.ActionResult, error)

	// Configuration
	SetDefaultTimeout(timeout time.Duration) error
	SetMaxConnections(max int) error
	EnableConnectionPooling(enabled bool) error

	// Utility operations
	TestConnection(host string) error
	GetConnectionStats() *types.ConnectionStats
	Close() error
}

// SSHClient defines the interface for SSH client operations
type SSHClient interface {
	Connect(host string, config *types.SSHConfig) (*types.SSHConnection, error)
	ExecuteCommand(connection *types.SSHConnection, command string) (*types.CommandResult, error)
	ExecuteScript(connection *types.SSHConnection, script string) (*types.CommandResult, error)
	CloseConnection(connection *types.SSHConnection) error
}

// AuthenticationEngine defines the interface for authentication operations
type AuthenticationEngine interface {
	Authenticate(connection *types.SSHConnection, auth *types.AuthenticationConfig) error
	ValidateAuthentication(auth *types.AuthenticationConfig) error
	GetSupportedMethods() []string
}

// ConnectionPool defines the interface for connection pooling
type ConnectionPool interface {
	GetConnection(host string) (*types.SSHConnection, error)
	ReturnConnection(connection *types.SSHConnection) error
	CloseConnection(connection *types.SSHConnection) error
	CloseAllConnections() error
	GetStats() *types.PoolStats
}

// ActingEngine defines the interface for action execution
type ActingEngine interface {
	ExecuteAction(connection *types.SSHConnection, action *types.SSHAction) (*types.ActionResult, error)
	ExecuteTemplate(connection *types.SSHConnection, template *types.TemplateAction) (*types.ActionResult, error)
	ExecuteSequential(connection *types.SSHConnection, actions []*types.SSHAction) (*types.ActionResult, error)
	ExecuteParallel(connection *types.SSHConnection, actions []*types.SSHAction) (*types.ActionResult, error)
}

// SSHKeyManager defines the interface for SSH key management
type SSHKeyManager interface {
	LoadPrivateKey(path string) (*types.SSHKey, error)
	ValidateKeyFile(path string) error
}
