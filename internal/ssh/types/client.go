package types

import (
	"net"
	"time"
)

// SSHConfig represents SSH connection configuration
type SSHConfig struct {
	Host     string        `hcl:"host"`
	Port     int           `hcl:"port,optional"`
	Username string        `hcl:"username,optional"`
	Password string        `hcl:"password,optional"`
	KeyFile  string        `hcl:"key_file,optional"`
	Timeout  time.Duration `hcl:"timeout,optional"`
}

// SSHConnection represents an SSH connection
type SSHConnection struct {
	Host      string    `hcl:"host"`
	Port      int       `hcl:"port"`
	Username  string    `hcl:"username"`
	Connected bool      `hcl:"connected"`
	CreatedAt time.Time `hcl:"created_at"`
	LastUsed  time.Time `hcl:"last_used"`
}

// CommandResult represents the result of a command execution
type CommandResult struct {
	Command  string        `hcl:"command"`
	ExitCode int           `hcl:"exit_code"`
	Stdout   string        `hcl:"stdout"`
	Stderr   string        `hcl:"stderr"`
	Duration time.Duration `hcl:"duration"`
	Error    string        `hcl:"error,optional"`
}

// ConnectionInfo represents connection information
type ConnectionInfo struct {
	Host      string    `hcl:"host"`
	Port      int       `hcl:"port"`
	Username  string    `hcl:"username"`
	Connected bool      `hcl:"connected"`
	CreatedAt time.Time `hcl:"created_at"`
	LastUsed  time.Time `hcl:"last_used"`
}

// ConnectionState represents the state of a connection
type ConnectionState struct {
	Connected bool      `hcl:"connected"`
	CreatedAt time.Time `hcl:"created_at"`
	LastUsed  time.Time `hcl:"last_used"`
	Error     string    `hcl:"error,optional"`
}

// HostKeyCallback represents a host key callback function
type HostKeyCallback func(hostname string, remote net.Addr, key []byte) error

// ClientConfig represents client configuration
type ClientConfig struct {
	DefaultTimeout  time.Duration `hcl:"default_timeout,optional"`
	MaxRetries      int           `hcl:"max_retries,optional"`
	HostKeyChecking bool          `hcl:"host_key_checking,optional"`
	KnownHostsFile  string        `hcl:"known_hosts_file,optional"`
}

// Config represents the main SSH configuration
type Config struct {
	ClientConfig            *ClientConfig         `hcl:"client,optional"`
	AuthenticationConfig    *AuthenticationConfig `hcl:"authentication,optional"`
	PoolConfig              *PoolConfig           `hcl:"pool,optional"`
	ActingConfig            *ActingConfig         `hcl:"acting,optional"`
	KeysConfig              *KeysConfig           `hcl:"keys,optional"`
	DefaultTimeout          time.Duration         `hcl:"default_timeout,optional"`
	MaxConnections          int                   `hcl:"max_connections,optional"`
	EnableConnectionPooling bool                  `hcl:"enable_connection_pooling,optional"`
}

// ConnectionStats represents connection statistics
type ConnectionStats struct {
	PoolStats *PoolStats `hcl:"pool_stats,optional"`
}

// Machine represents a machine for SSH operations
type Machine struct {
	Host     string `hcl:"host"`
	Port     int    `hcl:"port,optional"`
	Username string `hcl:"username,optional"`
	Password string `hcl:"password,optional"`
	KeyFile  string `hcl:"key_file,optional"`
}
