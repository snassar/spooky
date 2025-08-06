package ssh

import (
	"time"

	spookysshtypes "spooky/internal/ssh/types"
)

// SSHManager defines the main interface for SSH operations
type SSHManager interface {
	// Core SSH operations
	Connect(host string, config *spookysshtypes.SSHConfig) (*spookysshtypes.SSHConnection, error)
	ExecuteCommand(connection *spookysshtypes.SSHConnection, command string) (*spookysshtypes.CommandResult, error)
	ExecuteScript(connection *spookysshtypes.SSHConnection, script string) (*spookysshtypes.CommandResult, error)
	CloseConnection(connection *spookysshtypes.SSHConnection) error

	// Connection pool operations
	GetConnection(host string) (*spookysshtypes.SSHConnection, error)
	ReturnConnection(connection *spookysshtypes.SSHConnection) error
	CloseAllConnections() error

	// Authentication operations
	Authenticate(connection *spookysshtypes.SSHConnection, auth *spookysshtypes.AuthenticationConfig) error
	ValidateAuthentication(auth *spookysshtypes.AuthenticationConfig) error

	// Acting operations
	ExecuteAction(connection *spookysshtypes.SSHConnection, action *spookysshtypes.SSHAction) (*spookysshtypes.ActionResult, error)
	ExecuteTemplate(connection *spookysshtypes.SSHConnection, template *spookysshtypes.TemplateAction) (*spookysshtypes.ActionResult, error)

	// Configuration
	SetDefaultTimeout(timeout time.Duration) error
	SetMaxConnections(max int) error
	EnableConnectionPooling(enabled bool) error

	// Utility operations
	TestConnection(host string) error
	GetConnectionStats() *spookysshtypes.ConnectionStats
	Close() error
}

// SSHClient defines the interface for SSH client operations
type SSHClient interface {
	Connect(host string, config *spookysshtypes.SSHConfig) (*spookysshtypes.SSHConnection, error)
	ExecuteCommand(connection *spookysshtypes.SSHConnection, command string) (*spookysshtypes.CommandResult, error)
	ExecuteScript(connection *spookysshtypes.SSHConnection, script string) (*spookysshtypes.CommandResult, error)
	CloseConnection(connection *spookysshtypes.SSHConnection) error
}

// AuthenticationEngine defines the interface for authentication operations
type AuthenticationEngine interface {
	Authenticate(connection *spookysshtypes.SSHConnection, auth *spookysshtypes.AuthenticationConfig) error
	ValidateAuthentication(auth *spookysshtypes.AuthenticationConfig) error
	GetSupportedMethods() []string
}

// ConnectionPool defines the interface for connection pooling
type ConnectionPool interface {
	GetConnection(host string) (*spookysshtypes.SSHConnection, error)
	ReturnConnection(connection *spookysshtypes.SSHConnection) error
	CloseConnection(connection *spookysshtypes.SSHConnection) error
	CloseAllConnections() error
	GetStats() *spookysshtypes.PoolStats
}

// ActingEngine defines the interface for action execution
type ActingEngine interface {
	ExecuteAction(connection *spookysshtypes.SSHConnection, action *spookysshtypes.SSHAction) (*spookysshtypes.ActionResult, error)
	ExecuteTemplate(connection *spookysshtypes.SSHConnection, template *spookysshtypes.TemplateAction) (*spookysshtypes.ActionResult, error)
	ExecuteSequential(connection *spookysshtypes.SSHConnection, actions []*spookysshtypes.SSHAction) (*spookysshtypes.ActionResult, error)
	ExecuteParallel(connection *spookysshtypes.SSHConnection, actions []*spookysshtypes.SSHAction) (*spookysshtypes.ActionResult, error)
}

// SSHKeyManager defines the interface for SSH key management
type SSHKeyManager interface {
	LoadPrivateKey(path string) (*spookysshtypes.SSHKey, error)
	ValidateKeyFile(path string) error
}
