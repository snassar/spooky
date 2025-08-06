package client

import (
	"time"

	"spooky/internal/ssh/types"
)

// ClientManager defines the interface for SSH client operations
type ClientManager interface {
	// Core client operations
	Connect(host string, config *types.SSHConfig) (*types.SSHConnection, error)
	ExecuteCommand(connection *types.SSHConnection, command string) (*types.CommandResult, error)
	ExecuteScript(connection *types.SSHConnection, script string) (*types.CommandResult, error)
	CloseConnection(connection *types.SSHConnection) error

	// Configuration
	SetDefaultTimeout(timeout time.Duration) error
	SetMaxRetries(max int) error
	EnableHostKeyChecking(enabled bool) error

	// Utility operations
	TestConnection(host string) error
	GetConnectionInfo(connection *types.SSHConnection) *types.ConnectionInfo
	Close() error
}

// ConnectionManager defines the interface for connection management
type ConnectionManager interface {
	Connect(host string, config *types.SSHConfig) (*types.SSHConnection, error)
	CloseConnection(connection *types.SSHConnection) error
	TestConnection(connection *types.SSHConnection) error
	GetConnectionState(connection *types.SSHConnection) *types.ConnectionState
}

// ExecutionManager defines the interface for command execution
type ExecutionManager interface {
	ExecuteCommand(connection *types.SSHConnection, command string) (*types.CommandResult, error)
	ExecuteScript(connection *types.SSHConnection, script string) (*types.CommandResult, error)
	ExecuteWithTimeout(connection *types.SSHConnection, command string, timeout time.Duration) (*types.CommandResult, error)
}

// HostKeyManager defines the interface for host key management
type HostKeyManager interface {
	GetHostKeyCallback() types.HostKeyCallback
	ValidateHostKey(host string, key []byte) error
	AddKnownHost(host string, key []byte) error
	RemoveKnownHost(host string) error
}
