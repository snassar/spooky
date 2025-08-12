// Package ssh provides SSH client types for the spooky codebase.
// This package defines the data structures for SSH client operations and session management.
package ssh

import (
	"os"
	"time"

	spookytypescommon "spooky/internal/types/common"
)

// Client represents an SSH client for remote operations
type Client struct {
	spookytypescommon.CompleteEntity

	// Client configuration
	Config *ClientConfig `json:"config" hcl:"config"`

	// Client state
	Status    ClientStatus `json:"status" hcl:"status"`
	Connected bool         `json:"connected" hcl:"connected"`

	// Client metrics
	TotalConnections int           `json:"total_connections" hcl:"total_connections"`
	ActiveSessions   int           `json:"active_sessions" hcl:"active_sessions"`
	AverageLatency   time.Duration `json:"average_latency" hcl:"average_latency"`

	// Connection pool
	ConnectionPool *ConnectionPool `json:"connection_pool" hcl:"connection_pool"`

	// Client metadata
	Version      string   `json:"version" hcl:"version"`
	UserAgent    string   `json:"user_agent" hcl:"user_agent"`
	Capabilities []string `json:"capabilities" hcl:"capabilities"`
}

// ClientConfig represents SSH client configuration
type ClientConfig struct {
	// Connection settings
	DefaultHost    string        `json:"default_host" hcl:"default_host"`
	DefaultPort    int           `json:"default_port" hcl:"default_port" default:"22"`
	DefaultUser    string        `json:"default_user" hcl:"default_user"`
	DefaultTimeout time.Duration `json:"default_timeout" hcl:"default_timeout" default:"30s"`

	// Authentication settings
	DefaultKeyPath    string     `json:"default_key_path" hcl:"default_key_path"`
	DefaultAuthMethod AuthMethod `json:"default_auth_method" hcl:"default_auth_method" default:"public_key"`

	// Security settings
	KnownHostsPath     string `json:"known_hosts_path" hcl:"known_hosts_path"`
	StrictHostKeyCheck bool   `json:"strict_host_key_check" hcl:"strict_host_key_check" default:"true"`
	AllowInsecureHosts bool   `json:"allow_insecure_hosts" hcl:"allow_insecure_hosts" default:"false"`

	// Performance settings
	EnableCompression bool `json:"enable_compression" hcl:"enable_compression" default:"false"`
	EnableKeepalive   bool `json:"enable_keepalive" hcl:"enable_keepalive" default:"true"`

	// Connection pool settings
	MaxConnections     int           `json:"max_connections" hcl:"max_connections" default:"10"`
	MaxIdleConnections int           `json:"max_idle_connections" hcl:"max_idle_connections" default:"5"`
	IdleTimeout        time.Duration `json:"idle_timeout" hcl:"idle_timeout" default:"300s"`

	// Retry settings
	MaxRetryAttempts int           `json:"max_retry_attempts" hcl:"max_retry_attempts" default:"3"`
	RetryDelay       time.Duration `json:"retry_delay" hcl:"retry_delay" default:"5s"`

	// Keepalive settings
	KeepaliveInterval time.Duration `json:"keepalive_interval" hcl:"keepalive_interval" default:"60s"`
	KeepaliveCount    int           `json:"keepalive_count" hcl:"keepalive_count" default:"3"`
}

// ClientStatus represents the status of an SSH client
type ClientStatus string

const (
	ClientStatusInitialized  ClientStatus = "initialized"
	ClientStatusConnecting   ClientStatus = "connecting"
	ClientStatusConnected    ClientStatus = "connected"
	ClientStatusDisconnected ClientStatus = "disconnected"
	ClientStatusError        ClientStatus = "error"
	ClientStatusClosed       ClientStatus = "closed"
)

// Session represents an SSH session for command execution
type Session struct {
	spookytypescommon.CompleteEntity

	// Session details
	SessionID  string      `json:"session_id" hcl:"session_id"`
	Connection *Connection `json:"connection" hcl:"connection"`
	Client     *Client     `json:"client" hcl:"client"`

	// Session state
	Status    SessionStatus `json:"status" hcl:"status"`
	StartedAt time.Time     `json:"started_at" hcl:"started_at"`
	EndedAt   *time.Time    `json:"ended_at,omitempty" hcl:"ended_at,optional"`

	// Session configuration
	Environment map[string]string `json:"environment,omitempty" hcl:"environment,optional"`
	WorkingDir  string            `json:"working_dir,omitempty" hcl:"working_dir,optional"`
	Pty         *PtyConfig        `json:"pty,omitempty" hcl:"pty,optional"`

	// Session metrics
	CommandsExecuted int           `json:"commands_executed" hcl:"commands_executed"`
	TotalExecTime    time.Duration `json:"total_exec_time" hcl:"total_exec_time"`
	AverageExecTime  time.Duration `json:"average_exec_time" hcl:"average_exec_time"`

	// Session metadata
	UserAgent    string `json:"user_agent,omitempty" hcl:"user_agent,optional"`
	TerminalType string `json:"terminal_type,omitempty" hcl:"terminal_type,optional"`
}

// SessionStatus represents the status of an SSH session
type SessionStatus string

const (
	SessionStatusCreated   SessionStatus = "created"
	SessionStatusStarting  SessionStatus = "starting"
	SessionStatusActive    SessionStatus = "active"
	SessionStatusExecuting SessionStatus = "executing"
	SessionStatusCompleted SessionStatus = "completed"
	SessionStatusFailed    SessionStatus = "failed"
	SessionStatusClosed    SessionStatus = "closed"
)

// PtyConfig represents pseudo-terminal configuration
type PtyConfig struct {
	// Terminal type
	Term string `json:"term" hcl:"term" default:"xterm"`

	// Terminal dimensions
	Width  int `json:"width" hcl:"width" default:"80"`
	Height int `json:"height" hcl:"height" default:"24"`

	// Terminal modes
	Modes map[string]uint32 `json:"modes,omitempty" hcl:"modes,optional"`

	// Terminal features
	EnableEcho      bool `json:"enable_echo" hcl:"enable_echo" default:"true"`
	EnableCanonical bool `json:"enable_canonical" hcl:"enable_canonical" default:"true"`
	EnableSignals   bool `json:"enable_signals" hcl:"enable_signals" default:"true"`
}

// Command represents a command to be executed via SSH
type Command struct {
	spookytypescommon.CompleteEntity

	// Command details
	Command     string            `json:"command" hcl:"command"`
	Args        []string          `json:"args,omitempty" hcl:"args,optional"`
	WorkingDir  string            `json:"working_dir,omitempty" hcl:"working_dir,optional"`
	Environment map[string]string `json:"environment,omitempty" hcl:"environment,optional"`

	// Command configuration
	Timeout       time.Duration `json:"timeout,omitempty" hcl:"timeout,optional"`
	Pty           *PtyConfig    `json:"pty,omitempty" hcl:"pty,optional"`
	Stdin         string        `json:"stdin,omitempty" hcl:"stdin,optional"`
	CaptureOutput bool          `json:"capture_output" hcl:"capture_output" default:"true"`

	// Command metadata
	Priority    int       `json:"priority" hcl:"priority" default:"0"`
	ScheduledAt time.Time `json:"scheduled_at" hcl:"scheduled_at"`
	Tags        []string  `json:"tags,omitempty" hcl:"tags,optional"`

	// Security settings
	AllowUnsafe bool `json:"allow_unsafe" hcl:"allow_unsafe" default:"false"`
	RequireSudo bool `json:"require_sudo" hcl:"require_sudo" default:"false"`
}

// CommandResult represents the result of a command execution
type CommandResult struct {
	spookytypescommon.CompleteEntity

	// Command details
	Command *Command `json:"command" hcl:"command"`
	Session *Session `json:"session" hcl:"session"`

	// Execution results
	Success  bool   `json:"success" hcl:"success"`
	ExitCode int    `json:"exit_code" hcl:"exit_code"`
	Stdout   string `json:"stdout,omitempty" hcl:"stdout,optional"`
	Stderr   string `json:"stderr,omitempty" hcl:"stderr,optional"`
	Error    string `json:"error,omitempty" hcl:"error,optional"`

	// Execution metrics
	StartTime     time.Time     `json:"start_time" hcl:"start_time"`
	EndTime       time.Time     `json:"end_time" hcl:"end_time"`
	Duration      time.Duration `json:"duration" hcl:"duration"`
	RetryAttempts int           `json:"retry_attempts" hcl:"retry_attempts"`

	// Execution metadata
	Hostname    string            `json:"hostname,omitempty" hcl:"hostname,optional"`
	Username    string            `json:"username,omitempty" hcl:"username,optional"`
	WorkingDir  string            `json:"working_dir,omitempty" hcl:"working_dir,optional"`
	Environment map[string]string `json:"environment,omitempty" hcl:"environment,optional"`

	// Security information
	CommandHash string `json:"command_hash,omitempty" hcl:"command_hash,optional"`
	AuditTrail  string `json:"audit_trail,omitempty" hcl:"audit_trail,optional"`
}

// FileTransfer represents a file transfer operation via SSH
type FileTransfer struct {
	spookytypescommon.CompleteEntity

	// Transfer details
	LocalPath  string            `json:"local_path" hcl:"local_path"`
	RemotePath string            `json:"remote_path" hcl:"remote_path"`
	Direction  TransferDirection `json:"direction" hcl:"direction"`

	// Transfer configuration
	Mode        TransferMode `json:"mode" hcl:"mode" default:"scp"`
	Permissions os.FileMode  `json:"permissions,omitempty" hcl:"permissions,optional"`
	Owner       string       `json:"owner,omitempty" hcl:"owner,optional"`
	Group       string       `json:"group,omitempty" hcl:"group,optional"`

	// Transfer settings
	Resume       bool `json:"resume" hcl:"resume" default:"false"`
	Verify       bool `json:"verify" hcl:"verify" default:"true"`
	Compress     bool `json:"compress" hcl:"compress" default:"false"`
	PreserveTime bool `json:"preserve_time" hcl:"preserve_time" default:"true"`

	// Transfer metadata
	Size        int64     `json:"size,omitempty" hcl:"size,optional"`
	Checksum    string    `json:"checksum,omitempty" hcl:"checksum,optional"`
	ScheduledAt time.Time `json:"scheduled_at" hcl:"scheduled_at"`
}

// TransferDirection represents the direction of a file transfer
type TransferDirection string

const (
	TransferDirectionUpload   TransferDirection = "upload"
	TransferDirectionDownload TransferDirection = "download"
)

// TransferMode represents the mode of file transfer
type TransferMode string

const (
	TransferModeSCP  TransferMode = "scp"
	TransferModeSFTP TransferMode = "sftp"
)

// FileTransferResult represents the result of a file transfer operation
type FileTransferResult struct {
	spookytypescommon.CompleteEntity

	// Transfer details
	Transfer *FileTransfer `json:"transfer" hcl:"transfer"`
	Session  *Session      `json:"session" hcl:"session"`

	// Transfer results
	Success          bool   `json:"success" hcl:"success"`
	Error            string `json:"error,omitempty" hcl:"error,optional"`
	BytesTransferred int64  `json:"bytes_transferred" hcl:"bytes_transferred"`

	// Transfer metrics
	StartTime    time.Time     `json:"start_time" hcl:"start_time"`
	EndTime      time.Time     `json:"end_time" hcl:"end_time"`
	Duration     time.Duration `json:"duration" hcl:"duration"`
	TransferRate float64       `json:"transfer_rate" hcl:"transfer_rate"`

	// Transfer metadata
	LocalPath   string      `json:"local_path" hcl:"local_path"`
	RemotePath  string      `json:"remote_path" hcl:"remote_path"`
	Checksum    string      `json:"checksum,omitempty" hcl:"checksum,optional"`
	Permissions os.FileMode `json:"permissions,omitempty" hcl:"permissions,optional"`
}
