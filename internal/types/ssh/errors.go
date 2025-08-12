// Package ssh provides SSH error types for the spooky codebase.
// This package defines the data structures for SSH error handling and validation.
package ssh

import (
	"time"

	spookytypescommon "spooky/internal/types/common"
)

// SSHError represents a generic SSH error
type SSHError struct {
	spookytypescommon.ErrorDetails

	// Error details
	ErrorType    SSHErrorType `json:"error_type" hcl:"error_type"`
	ErrorCode    int          `json:"error_code" hcl:"error_code"`
	ErrorMessage string       `json:"error_message" hcl:"error_message"`

	// Error context
	Hostname  string `json:"hostname,omitempty" hcl:"hostname,optional"`
	Port      int    `json:"port,omitempty" hcl:"port,optional"`
	Username  string `json:"username,omitempty" hcl:"username,optional"`
	Operation string `json:"operation,omitempty" hcl:"operation,optional"`

	// Error metadata
	Timestamp   time.Time `json:"timestamp" hcl:"timestamp"`
	Retryable   bool      `json:"retryable" hcl:"retryable"`
	Recoverable bool      `json:"recoverable" hcl:"recoverable"`

	// Error details
	Details map[string]interface{} `json:"details,omitempty" hcl:"details,optional"`
}

// SSHErrorType represents the type of SSH error
type SSHErrorType string

const (
	SSHErrorTypeConnection     SSHErrorType = "connection"
	SSHErrorTypeAuthentication SSHErrorType = "authentication"
	SSHErrorTypeAuthorization  SSHErrorType = "authorization"
	SSHErrorTypeTimeout        SSHErrorType = "timeout"
	SSHErrorTypeProtocol       SSHErrorType = "protocol"
	SSHErrorTypeHostKey        SSHErrorType = "host_key"
	SSHErrorTypeCommand        SSHErrorType = "command"
	SSHErrorTypeSession        SSHErrorType = "session"
	SSHErrorTypeFileTransfer   SSHErrorType = "file_transfer"
	SSHErrorTypeValidation     SSHErrorType = "validation"
	SSHErrorTypeConfiguration  SSHErrorType = "configuration"
	SSHErrorTypeUnknown        SSHErrorType = "unknown"
)

// ConnectionError represents an SSH connection error
type ConnectionError struct {
	SSHError

	// Connection details
	ConnectionAttempts int           `json:"connection_attempts" hcl:"connection_attempts"`
	ConnectionTimeout  time.Duration `json:"connection_timeout" hcl:"connection_timeout"`
	LastAttemptTime    time.Time     `json:"last_attempt_time" hcl:"last_attempt_time"`

	// Network details
	NetworkError    string `json:"network_error,omitempty" hcl:"network_error,optional"`
	DNSResolution   bool   `json:"dns_resolution" hcl:"dns_resolution"`
	PortReachable   bool   `json:"port_reachable" hcl:"port_reachable"`
	FirewallBlocked bool   `json:"firewall_blocked" hcl:"firewall_blocked"`

	// Connection context
	LocalAddress  string `json:"local_address,omitempty" hcl:"local_address,optional"`
	RemoteAddress string `json:"remote_address,omitempty" hcl:"remote_address,optional"`
	ProxyUsed     bool   `json:"proxy_used" hcl:"proxy_used"`
}

// AuthenticationError represents an SSH authentication error
type AuthenticationError struct {
	SSHError

	// Authentication details
	AuthMethod      AuthMethod `json:"auth_method" hcl:"auth_method"`
	AuthAttempts    int        `json:"auth_attempts" hcl:"auth_attempts"`
	MaxAuthAttempts int        `json:"max_auth_attempts" hcl:"max_auth_attempts"`

	// Key details (for public key auth)
	KeyPath        string  `json:"key_path,omitempty" hcl:"key_path,optional"`
	KeyType        KeyType `json:"key_type,omitempty" hcl:"key_type,optional"`
	KeyFingerprint string  `json:"key_fingerprint,omitempty" hcl:"key_fingerprint,optional"`
	KeyValid       bool    `json:"key_valid" hcl:"key_valid"`

	// Password details (for password auth)
	PasswordProvided bool `json:"password_provided" hcl:"password_provided"`
	PasswordValid    bool `json:"password_valid" hcl:"password_valid"`

	// Authentication context
	ServerAuthMethods []string `json:"server_auth_methods,omitempty" hcl:"server_auth_methods,optional"`
	ClientAuthMethods []string `json:"client_auth_methods,omitempty" hcl:"client_auth_methods,optional"`
}

// HostKeyError represents an SSH host key error
type HostKeyError struct {
	SSHError

	// Host key details
	HostKeyType         KeyType `json:"host_key_type" hcl:"host_key_type"`
	HostKeyFingerprint  string  `json:"host_key_fingerprint" hcl:"host_key_fingerprint"`
	ExpectedFingerprint string  `json:"expected_fingerprint,omitempty" hcl:"expected_fingerprint,optional"`

	// Host key validation
	KnownHostsPath     string `json:"known_hosts_path,omitempty" hcl:"known_hosts_path,optional"`
	StrictHostKeyCheck bool   `json:"strict_host_key_check" hcl:"strict_host_key_check"`
	HostKeyFound       bool   `json:"host_key_found" hcl:"host_key_found"`
	HostKeyTrusted     bool   `json:"host_key_trusted" hcl:"host_key_trusted"`

	// Host key context
	HostKeyAlgorithm string `json:"host_key_algorithm,omitempty" hcl:"host_key_algorithm,optional"`
	HostKeySize      int    `json:"host_key_size,omitempty" hcl:"host_key_size,optional"`
	HostKeyComment   string `json:"host_key_comment,omitempty" hcl:"host_key_comment,optional"`
}

// CommandError represents an SSH command execution error
type CommandError struct {
	SSHError

	// Command details
	Command    string   `json:"command" hcl:"command"`
	Args       []string `json:"args,omitempty" hcl:"args,optional"`
	WorkingDir string   `json:"working_dir,omitempty" hcl:"working_dir,optional"`
	ExitCode   int      `json:"exit_code" hcl:"exit_code"`

	// Command execution
	CommandTimeout time.Duration `json:"command_timeout" hcl:"command_timeout"`
	ExecutionTime  time.Duration `json:"execution_time" hcl:"execution_time"`
	Killed         bool          `json:"killed" hcl:"killed"`

	// Command output
	Stdout string `json:"stdout,omitempty" hcl:"stdout,optional"`
	Stderr string `json:"stderr,omitempty" hcl:"stderr,optional"`

	// Command context
	Environment map[string]string `json:"environment,omitempty" hcl:"environment,optional"`
	User        string            `json:"user,omitempty" hcl:"user,optional"`
	Group       string            `json:"group,omitempty" hcl:"group,optional"`
}

// SessionError represents an SSH session error
type SessionError struct {
	SSHError

	// Session details
	SessionID      string        `json:"session_id" hcl:"session_id"`
	SessionType    string        `json:"session_type" hcl:"session_type"`
	SessionTimeout time.Duration `json:"session_timeout" hcl:"session_timeout"`

	// Session state
	SessionActive  bool       `json:"session_active" hcl:"session_active"`
	SessionCreated time.Time  `json:"session_created" hcl:"session_created"`
	SessionClosed  *time.Time `json:"session_closed,omitempty" hcl:"session_closed,optional"`

	// Session metrics
	CommandsExecuted int           `json:"commands_executed" hcl:"commands_executed"`
	TotalExecTime    time.Duration `json:"total_exec_time" hcl:"total_exec_time"`
	LastCommandTime  *time.Time    `json:"last_command_time,omitempty" hcl:"last_command_time,optional"`

	// Session context
	PtyRequested bool   `json:"pty_requested" hcl:"pty_requested"`
	PtyAllocated bool   `json:"pty_allocated" hcl:"pty_allocated"`
	TerminalType string `json:"terminal_type,omitempty" hcl:"terminal_type,optional"`
	TerminalSize string `json:"terminal_size,omitempty" hcl:"terminal_size,optional"`
}

// FileTransferError represents an SSH file transfer error
type FileTransferError struct {
	SSHError

	// Transfer details
	TransferMode      TransferMode      `json:"transfer_mode" hcl:"transfer_mode"`
	TransferDirection TransferDirection `json:"transfer_direction" hcl:"transfer_direction"`
	LocalPath         string            `json:"local_path" hcl:"local_path"`
	RemotePath        string            `json:"remote_path" hcl:"remote_path"`

	// Transfer state
	BytesTransferred int64   `json:"bytes_transferred" hcl:"bytes_transferred"`
	TotalBytes       int64   `json:"total_bytes" hcl:"total_bytes"`
	TransferProgress float64 `json:"transfer_progress" hcl:"transfer_progress"`

	// Transfer metrics
	TransferStartTime time.Time     `json:"transfer_start_time" hcl:"transfer_start_time"`
	TransferEndTime   time.Time     `json:"transfer_end_time" hcl:"transfer_end_time"`
	TransferDuration  time.Duration `json:"transfer_duration" hcl:"transfer_duration"`
	TransferRate      float64       `json:"transfer_rate" hcl:"transfer_rate"`

	// Transfer context
	FilePermissions string `json:"file_permissions,omitempty" hcl:"file_permissions,optional"`
	FileOwner       string `json:"file_owner,omitempty" hcl:"file_owner,optional"`
	FileGroup       string `json:"file_group,omitempty" hcl:"file_group,optional"`
	FileSize        int64  `json:"file_size,omitempty" hcl:"file_size,optional"`
	FileChecksum    string `json:"file_checksum,omitempty" hcl:"file_checksum,optional"`
}

// ValidationError represents an SSH validation error
type ValidationError struct {
	SSHError

	// Validation details
	ValidationType ValidationType `json:"validation_type" hcl:"validation_type"`
	FieldName      string         `json:"field_name" hcl:"field_name"`
	FieldValue     interface{}    `json:"field_value,omitempty" hcl:"field_value,optional"`
	ExpectedValue  interface{}    `json:"expected_value,omitempty" hcl:"expected_value,optional"`

	// Validation context
	ValidationRule    string          `json:"validation_rule,omitempty" hcl:"validation_rule,optional"`
	ValidationMessage string          `json:"validation_message" hcl:"validation_message"`
	ValidationLevel   ValidationLevel `json:"validation_level" hcl:"validation_level"`

	// Validation metadata
	ValidationSource string    `json:"validation_source,omitempty" hcl:"validation_source,optional"`
	ValidationTime   time.Time `json:"validation_time" hcl:"validation_time"`
}

// ValidationType represents the type of validation error
type ValidationType string

const (
	ValidationTypeRequired    ValidationType = "required"
	ValidationTypeFormat      ValidationType = "format"
	ValidationTypeRange       ValidationType = "range"
	ValidationTypeLength      ValidationType = "length"
	ValidationTypePattern     ValidationType = "pattern"
	ValidationTypeEnum        ValidationType = "enum"
	ValidationTypeCustom      ValidationType = "custom"
	ValidationTypeDependency  ValidationType = "dependency"
	ValidationTypeConsistency ValidationType = "consistency"
)

// ValidationLevel represents the level of validation error
type ValidationLevel string

const (
	ValidationLevelError   ValidationLevel = "error"
	ValidationLevelWarning ValidationLevel = "warning"
	ValidationLevelInfo    ValidationLevel = "info"
)

// ConfigurationError represents an SSH configuration error
type ConfigurationError struct {
	SSHError

	// Configuration details
	ConfigPath    string `json:"config_path" hcl:"config_path"`
	ConfigSection string `json:"config_section" hcl:"config_section"`
	ConfigKey     string `json:"config_key" hcl:"config_key"`

	// Configuration context
	ConfigValue  interface{} `json:"config_value,omitempty" hcl:"config_value,optional"`
	ExpectedType string      `json:"expected_type,omitempty" hcl:"expected_type,optional"`
	DefaultValue interface{} `json:"default_value,omitempty" hcl:"default_value,optional"`

	// Configuration validation
	ConfigValid  bool   `json:"config_valid" hcl:"config_valid"`
	ConfigError  string `json:"config_error,omitempty" hcl:"config_error,optional"`
	ConfigSource string `json:"config_source,omitempty" hcl:"config_source,optional"`

	// Configuration metadata
	ConfigLine     int       `json:"config_line,omitempty" hcl:"config_line,optional"`
	ConfigColumn   int       `json:"config_column,omitempty" hcl:"config_column,optional"`
	ConfigModified time.Time `json:"config_modified" hcl:"config_modified"`
}
