package interfaces

import (
	"time"

	spookytypesssh "spooky/internal/types/ssh"
)

// SSHManager defines the main interface for SSH operations
type SSHManager interface {
	// Core SSH operations
	Connect(host string, config *spookytypesssh.SSHConfig) (*spookytypesssh.SSHConnection, error)
	ExecuteCommand(connection *spookytypesssh.SSHConnection, command string) (*spookytypesssh.CommandResult, error)
	ExecuteScript(connection *spookytypesssh.SSHConnection, script string) (*spookytypesssh.CommandResult, error)
	CloseConnection(connection *spookytypesssh.SSHConnection) error

	// Connection pool operations
	GetConnection(host string) (*spookytypesssh.SSHConnection, error)
	ReturnConnection(connection *spookytypesssh.SSHConnection) error
	CloseAllConnections() error

	// Authentication operations
	Authenticate(connection *spookytypesssh.SSHConnection, auth *spookytypesssh.AuthenticationConfig) error
	ValidateAuthentication(auth *spookytypesssh.AuthenticationConfig) error

	// Acting operations
	ExecuteAction(connection *spookytypesssh.SSHConnection, action *spookytypesssh.SSHAction) (*spookytypesssh.ActionResult, error)
	ExecuteTemplate(connection *spookytypesssh.SSHConnection, template *spookytypesssh.TemplateAction) (*spookytypesssh.ActionResult, error)

	// Configuration
	SetDefaultTimeout(timeout time.Duration) error
	SetMaxConnections(max int) error
	EnableConnectionPooling(enabled bool) error

	// Utility operations
	TestConnection(host string) error
	GetConnectionStats() *spookytypesssh.ConnectionStats
	Close() error
}

// SSHClient defines the interface for SSH client operations
type SSHClient interface {
	Connect(host string, config *spookytypesssh.SSHConfig) (*spookytypesssh.SSHConnection, error)
	ExecuteCommand(connection *spookytypesssh.SSHConnection, command string) (*spookytypesssh.CommandResult, error)
	ExecuteScript(connection *spookytypesssh.SSHConnection, script string) (*spookytypesssh.CommandResult, error)
	CloseConnection(connection *spookytypesssh.SSHConnection) error
}

// AuthenticationEngine defines the interface for authentication operations
type AuthenticationEngine interface {
	Authenticate(connection *spookytypesssh.SSHConnection, auth *spookytypesssh.AuthenticationConfig) error
	ValidateAuthentication(auth *spookytypesssh.AuthenticationConfig) error
	GetSupportedMethods() []string
}

// ConnectionPool defines the interface for connection pooling
type ConnectionPool interface {
	GetConnection(host string) (*spookytypesssh.SSHConnection, error)
	ReturnConnection(connection *spookytypesssh.SSHConnection) error
	CloseConnection(connection *spookytypesssh.SSHConnection) error
	CloseAllConnections() error
	GetStats() *spookytypesssh.PoolStats
}

// ActingEngine defines the interface for action execution
type ActingEngine interface {
	ExecuteAction(connection *spookytypesssh.SSHConnection, action *spookytypesssh.SSHAction) (*spookytypesssh.ActionResult, error)
	ExecuteTemplate(connection *spookytypesssh.SSHConnection, template *spookytypesssh.TemplateAction) (*spookytypesssh.ActionResult, error)
	ExecuteSequential(connection *spookytypesssh.SSHConnection, actions []*spookytypesssh.SSHAction) (*spookytypesssh.ActionResult, error)
	ExecuteParallel(connection *spookytypesssh.SSHConnection, actions []*spookytypesssh.SSHAction) (*spookytypesssh.ActionResult, error)
}

// SSHKeyManager defines the interface for SSH key management
type SSHKeyManager interface {
	LoadPrivateKey(path string) (*spookytypesssh.SSHKey, error)
	ValidateKeyFile(path string) error
}

// ClientManager defines the interface for SSH client operations
type ClientManager interface {
	// Core client operations
	Connect(host string, config *spookytypesssh.SSHConfig) (*spookytypesssh.SSHConnection, error)
	ExecuteCommand(connection *spookytypesssh.SSHConnection, command string) (*spookytypesssh.CommandResult, error)
	ExecuteScript(connection *spookytypesssh.SSHConnection, script string) (*spookytypesssh.CommandResult, error)
	CloseConnection(connection *spookytypesssh.SSHConnection) error

	// Configuration
	SetDefaultTimeout(timeout time.Duration) error
	SetMaxRetries(max int) error
	EnableHostKeyChecking(enabled bool) error

	// Utility operations
	TestConnection(host string) error
	GetConnectionInfo(connection *spookytypesssh.SSHConnection) *spookytypesssh.ConnectionInfo
	Close() error
}

// ConnectionManager defines the interface for connection management
type ConnectionManager interface {
	Connect(host string, config *spookytypesssh.SSHConfig) (*spookytypesssh.SSHConnection, error)
	CloseConnection(connection *spookytypesssh.SSHConnection) error
	TestConnection(connection *spookytypesssh.SSHConnection) error
	GetConnectionState(connection *spookytypesssh.SSHConnection) *spookytypesssh.ConnectionState
}

// ExecutionManager defines the interface for command execution
type ExecutionManager interface {
	ExecuteCommand(connection *spookytypesssh.SSHConnection, command string) (*spookytypesssh.CommandResult, error)
	ExecuteScript(connection *spookytypesssh.SSHConnection, script string) (*spookytypesssh.CommandResult, error)
	ExecuteWithTimeout(connection *spookytypesssh.SSHConnection, command string, timeout time.Duration) (*spookytypesssh.CommandResult, error)
}

// HostKeyManager defines the interface for host key management
type HostKeyManager interface {
	GetHostKeyCallback() spookytypesssh.HostKeyCallback
	ValidateHostKey(host string, key []byte) error
	AddKnownHost(host string, key []byte) error
	RemoveKnownHost(host string) error
}
