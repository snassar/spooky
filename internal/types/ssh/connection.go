// Package ssh provides SSH-related types for the spooky codebase.
// This package defines the data structures for SSH connections, authentication, and acting.
package ssh

import (
	"time"

	spookytypescommon "spooky/internal/types/common"
)

// Connection represents an SSH connection to a remote machine
type Connection struct {
	spookytypescommon.CompleteEntity

	// Connection details
	Host     string `json:"host" hcl:"host"`
	Port     int    `json:"port" hcl:"port"`
	User     string `json:"user" hcl:"user"`
	Protocol string `json:"protocol" hcl:"protocol" default:"ssh"`

	// Connection state
	Status      ConnectionStatus `json:"status" hcl:"status"`
	ConnectedAt *time.Time       `json:"connected_at,omitempty" hcl:"connected_at,optional"`
	LastUsed    *time.Time       `json:"last_used,omitempty" hcl:"last_used,optional"`

	// Connection metrics
	Latency      time.Duration `json:"latency,omitempty" hcl:"latency,optional"`
	ErrorCount   int           `json:"error_count,omitempty" hcl:"error_count,optional"`
	SuccessCount int           `json:"success_count,omitempty" hcl:"success_count,optional"`

	// Connection configuration
	Timeout           time.Duration `json:"timeout,omitempty" hcl:"timeout,optional"`
	KeepaliveInterval time.Duration `json:"keepalive_interval,omitempty" hcl:"keepalive_interval,optional"`
	KeepaliveCount    int           `json:"keepalive_count,omitempty" hcl:"keepalive_count,optional"`

	// Authentication information
	AuthMethod     AuthMethod `json:"auth_method" hcl:"auth_method"`
	KeyPath        string     `json:"key_path,omitempty" hcl:"key_path,optional"`
	KeyFingerprint string     `json:"key_fingerprint,omitempty" hcl:"key_fingerprint,optional"`

	// Host key verification
	HostKeyFingerprint string `json:"host_key_fingerprint,omitempty" hcl:"host_key_fingerprint,optional"`
	KnownHostsPath     string `json:"known_hosts_path,omitempty" hcl:"known_hosts_path,optional"`
	StrictHostKeyCheck bool   `json:"strict_host_key_check" hcl:"strict_host_key_check" default:"true"`

	// Connection metadata
	ClientVersion string `json:"client_version,omitempty" hcl:"client_version,optional"`
	ServerVersion string `json:"server_version,omitempty" hcl:"server_version,optional"`
	Compression   bool   `json:"compression" hcl:"compression" default:"false"`
}

// ConnectionStatus represents the status of an SSH connection
type ConnectionStatus string

const (
	ConnectionStatusDisconnected ConnectionStatus = "disconnected"
	ConnectionStatusConnecting   ConnectionStatus = "connecting"
	ConnectionStatusConnected    ConnectionStatus = "connected"
	ConnectionStatusFailed       ConnectionStatus = "failed"
	ConnectionStatusTimeout      ConnectionStatus = "timeout"
	ConnectionStatusClosed       ConnectionStatus = "closed"
)

// AuthMethod represents the authentication method used for SSH connections
type AuthMethod string

const (
	AuthMethodPassword  AuthMethod = "password"
	AuthMethodPublicKey AuthMethod = "public_key"
	AuthMethodKeyboard  AuthMethod = "keyboard_interactive"
	AuthMethodAgent     AuthMethod = "agent"
	AuthMethodNone      AuthMethod = "none"
)

// ConnectionPool represents a pool of SSH connections
type ConnectionPool struct {
	spookytypescommon.CompleteEntity

	// Pool configuration
	MaxConnections     int           `json:"max_connections" hcl:"max_connections" default:"10"`
	MaxIdleConnections int           `json:"max_idle_connections" hcl:"max_idle_connections" default:"5"`
	IdleTimeout        time.Duration `json:"idle_timeout" hcl:"idle_timeout" default:"300s"`
	ConnectionTimeout  time.Duration `json:"connection_timeout" hcl:"connection_timeout" default:"30s"`

	// Pool state
	ActiveConnections int `json:"active_connections" hcl:"active_connections"`
	IdleConnections   int `json:"idle_connections" hcl:"idle_connections"`
	TotalConnections  int `json:"total_connections" hcl:"total_connections"`

	// Pool metrics
	ConnectionAttempts int           `json:"connection_attempts" hcl:"connection_attempts"`
	ConnectionErrors   int           `json:"connection_errors" hcl:"connection_errors"`
	AverageLatency     time.Duration `json:"average_latency" hcl:"average_latency"`

	// Pool configuration
	EnableCompression  bool `json:"enable_compression" hcl:"enable_compression" default:"false"`
	EnableKeepalive    bool `json:"enable_keepalive" hcl:"enable_keepalive" default:"true"`
	EnableHostKeyCheck bool `json:"enable_host_key_check" hcl:"enable_host_key_check" default:"true"`

	// Connection factory
	ConnectionFactory ConnectionFactory `json:"connection_factory" hcl:"connection_factory"`
}

// ConnectionFactory creates new SSH connections
type ConnectionFactory struct {
	// Default authentication settings
	DefaultUser       string     `json:"default_user" hcl:"default_user"`
	DefaultKeyPath    string     `json:"default_key_path" hcl:"default_key_path"`
	DefaultAuthMethod AuthMethod `json:"default_auth_method" hcl:"default_auth_method" default:"public_key"`

	// Default connection settings
	DefaultTimeout           time.Duration `json:"default_timeout" hcl:"default_timeout" default:"30s"`
	DefaultKeepaliveInterval time.Duration `json:"default_keepalive_interval" hcl:"default_keepalive_interval" default:"60s"`
	DefaultKeepaliveCount    int           `json:"default_keepalive_count" hcl:"default_keepalive_count" default:"3"`

	// Security settings
	KnownHostsPath     string `json:"known_hosts_path" hcl:"known_hosts_path"`
	StrictHostKeyCheck bool   `json:"strict_host_key_check" hcl:"strict_host_key_check" default:"true"`
	AllowInsecureHosts bool   `json:"allow_insecure_hosts" hcl:"allow_insecure_hosts" default:"false"`

	// Performance settings
	EnableCompression bool `json:"enable_compression" hcl:"enable_compression" default:"false"`
	EnableKeepalive   bool `json:"enable_keepalive" hcl:"enable_keepalive" default:"true"`
}

// ConnectionRequest represents a request to establish an SSH connection
type ConnectionRequest struct {
	// Connection details
	Host string `json:"host" hcl:"host"`
	Port int    `json:"port" hcl:"port" default:"22"`
	User string `json:"user" hcl:"user"`

	// Authentication
	Password        string     `json:"password,omitempty" hcl:"password,optional" sensitive:"true"`
	KeyPath         string     `json:"key_path,omitempty" hcl:"key_path,optional"`
	CertificatePath string     `json:"certificate_path,omitempty" hcl:"certificate_path,optional"`
	Passphrase      string     `json:"passphrase,omitempty" hcl:"passphrase,optional" sensitive:"true"`
	AuthMethod      AuthMethod `json:"auth_method" hcl:"auth_method" default:"public_key"`

	// Connection settings
	Timeout           time.Duration `json:"timeout,omitempty" hcl:"timeout,optional"`
	KeepaliveInterval time.Duration `json:"keepalive_interval,omitempty" hcl:"keepalive_interval,optional"`
	KeepaliveCount    int           `json:"keepalive_count,omitempty" hcl:"keepalive_count,optional"`

	// Security settings
	KnownHostsPath     string `json:"known_hosts_path,omitempty" hcl:"known_hosts_path,optional"`
	StrictHostKeyCheck bool   `json:"strict_host_key_check" hcl:"strict_host_key_check" default:"true"`
	AllowInsecureHosts bool   `json:"allow_insecure_hosts" hcl:"allow_insecure_hosts" default:"false"`

	// Performance settings
	EnableCompression bool `json:"enable_compression" hcl:"enable_compression" default:"false"`
	EnableKeepalive   bool `json:"enable_keepalive" hcl:"enable_keepalive" default:"true"`

	// Request metadata
	RequestID   string    `json:"request_id" hcl:"request_id"`
	RequestedAt time.Time `json:"requested_at" hcl:"requested_at"`
	Priority    int       `json:"priority" hcl:"priority" default:"0"`
}

// ConnectionResult represents the result of an SSH connection attempt
type ConnectionResult struct {
	spookytypescommon.CompleteEntity

	// Connection details
	Connection *Connection        `json:"connection" hcl:"connection"`
	Request    *ConnectionRequest `json:"request" hcl:"request"`

	// Result status
	Success bool   `json:"success" hcl:"success"`
	Error   string `json:"error,omitempty" hcl:"error,optional"`

	// Connection metrics
	ConnectTime   time.Duration `json:"connect_time" hcl:"connect_time"`
	Latency       time.Duration `json:"latency" hcl:"latency"`
	RetryAttempts int           `json:"retry_attempts" hcl:"retry_attempts"`

	// Connection information
	ClientVersion      string `json:"client_version,omitempty" hcl:"client_version,optional"`
	ServerVersion      string `json:"server_version,omitempty" hcl:"server_version,optional"`
	HostKeyFingerprint string `json:"host_key_fingerprint,omitempty" hcl:"host_key_fingerprint,optional"`

	// Result metadata
	CompletedAt time.Time `json:"completed_at" hcl:"completed_at"`
}
