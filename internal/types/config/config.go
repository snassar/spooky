// Package config provides types for configuration management in the spooky codebase.
// These types define the structure for spooky configuration and environment settings.
package config

import (
	"time"
)

// Config represents the main spooky configuration
type Config struct {
	// Configuration metadata
	Version   string    `json:"version" hcl:"version"`
	CreatedAt time.Time `json:"created_at" hcl:"created_at"`
	UpdatedAt time.Time `json:"updated_at" hcl:"updated_at"`

	// Global settings
	Global *GlobalConfig `json:"global" hcl:"global"`

	// CLI settings
	CLI *CLIConfig `json:"cli" hcl:"cli"`

	// Logging settings
	Logging *LoggingConfig `json:"logging" hcl:"logging"`

	// SSH settings
	SSH *SSHConfig `json:"ssh" hcl:"ssh"`

	// Storage settings
	Storage *StorageConfig `json:"storage" hcl:"storage"`

	// Security settings
	Security *SecurityConfig `json:"security" hcl:"security"`
}

// GlobalConfig provides global configuration settings
type GlobalConfig struct {
	// Default project path
	DefaultProjectPath string `json:"default_project_path,omitempty" hcl:"default_project_path,optional"`

	// Default parallel workers
	DefaultParallelWorkers int `json:"default_parallel_workers,omitempty" hcl:"default_parallel_workers,optional"`

	// Default timeout in seconds
	DefaultTimeout int `json:"default_timeout,omitempty" hcl:"default_timeout,optional"`

	// Default log level
	DefaultLogLevel string `json:"default_log_level,omitempty" hcl:"default_log_level,optional"`

	// Whether to enable dry-run by default
	DefaultDryRun bool `json:"default_dry_run,omitempty" hcl:"default_dry_run,optional"`

	// Whether to enable verbose output by default
	DefaultVerbose bool `json:"default_verbose,omitempty" hcl:"default_verbose,optional"`
}

// CLIConfig provides CLI-specific configuration
type CLIConfig struct {
	// CLI theme
	Theme string `json:"theme,omitempty" hcl:"theme,optional"`

	// CLI colors
	Colors bool `json:"colors,omitempty" hcl:"colors,optional"`

	// CLI progress bars
	ProgressBars bool `json:"progress_bars,omitempty" hcl:"progress_bars,optional"`

	// CLI confirmation prompts
	ConfirmPrompts bool `json:"confirm_prompts,omitempty" hcl:"confirm_prompts,optional"`

	// CLI output format
	OutputFormat string `json:"output_format,omitempty" hcl:"output_format,optional"`
}

// LoggingConfig provides logging configuration
type LoggingConfig struct {
	// Log level
	Level string `json:"level" hcl:"level"`

	// Log format
	Format string `json:"format" hcl:"format"`

	// Log output
	Output string `json:"output" hcl:"output"`

	// Log file path
	FilePath string `json:"file_path,omitempty" hcl:"file_path,optional"`

	// Log file permissions
	FilePermissions string `json:"file_permissions,omitempty" hcl:"file_permissions,optional"`

	// Log file max size in MB
	FileMaxSize int `json:"file_max_size,omitempty" hcl:"file_max_size,optional"`

	// Log file max age in days
	FileMaxAge int `json:"file_max_age,omitempty" hcl:"file_max_age,optional"`

	// Log file max backups
	FileMaxBackups int `json:"file_max_backups,omitempty" hcl:"file_max_backups,optional"`

	// Whether to compress log files
	FileCompress bool `json:"file_compress,omitempty" hcl:"file_compress,optional"`
}

// SSHConfig provides SSH configuration
type SSHConfig struct {
	// Default SSH port
	DefaultPort int `json:"default_port,omitempty" hcl:"default_port,optional"`

	// Default SSH user
	DefaultUser string `json:"default_user,omitempty" hcl:"default_user,optional"`

	// Default SSH key path
	DefaultKeyPath string `json:"default_key_path,omitempty" hcl:"default_key_path,optional"`

	// SSH connection timeout in seconds
	ConnectionTimeout int `json:"connection_timeout,omitempty" hcl:"connection_timeout,optional"`

	// SSH command timeout in seconds
	CommandTimeout int `json:"command_timeout,omitempty" hcl:"command_timeout,optional"`

	// SSH connection pool size
	ConnectionPoolSize int `json:"connection_pool_size,omitempty" hcl:"connection_pool_size,optional"`

	// SSH connection pool timeout in seconds
	ConnectionPoolTimeout int `json:"connection_pool_timeout,omitempty" hcl:"connection_pool_timeout,optional"`

	// Whether to enable SSH connection pooling
	EnableConnectionPool bool `json:"enable_connection_pool,omitempty" hcl:"enable_connection_pool,optional"`

	// Whether to enable SSH host key verification
	EnableHostKeyVerification bool `json:"enable_host_key_verification,omitempty" hcl:"enable_host_key_verification,optional"`

	// SSH known hosts file path
	KnownHostsPath string `json:"known_hosts_path,omitempty" hcl:"known_hosts_path,optional"`
}

// StorageConfig provides storage configuration
type StorageConfig struct {
	// Storage type
	Type string `json:"type" hcl:"type"`

	// Storage path
	Path string `json:"path" hcl:"path"`

	// Storage format
	Format string `json:"format" hcl:"format"`

	// Storage compression
	Compression bool `json:"compression,omitempty" hcl:"compression,optional"`

	// Storage encryption
	Encryption bool `json:"encryption,omitempty" hcl:"encryption,optional"`

	// Storage encryption key path
	EncryptionKeyPath string `json:"encryption_key_path,omitempty" hcl:"encryption_key_path,optional"`

	// Storage backup enabled
	BackupEnabled bool `json:"backup_enabled,omitempty" hcl:"backup_enabled,optional"`

	// Storage backup path
	BackupPath string `json:"backup_path,omitempty" hcl:"backup_path,optional"`

	// Storage backup retention in days
	BackupRetention int `json:"backup_retention,omitempty" hcl:"backup_retention,optional"`
}

// SecurityConfig provides security configuration
type SecurityConfig struct {
	// Whether to enable audit logging
	AuditLogging bool `json:"audit_logging,omitempty" hcl:"audit_logging,optional"`

	// Audit log path
	AuditLogPath string `json:"audit_log_path,omitempty" hcl:"audit_log_path,optional"`

	// Whether to enable sensitive data masking
	SensitiveDataMasking bool `json:"sensitive_data_masking,omitempty" hcl:"sensitive_data_masking,optional"`

	// Sensitive data patterns
	SensitiveDataPatterns []string `json:"sensitive_data_patterns,omitempty" hcl:"sensitive_data_patterns,optional"`

	// Whether to enable certificate verification
	CertificateVerification bool `json:"certificate_verification,omitempty" hcl:"certificate_verification,optional"`

	// Certificate authority path
	CAPath string `json:"ca_path,omitempty" hcl:"ca_path,optional"`
}

// Defaults provides default configuration values
type Defaults struct {
	// Default configuration values
	Defaults map[string]interface{} `json:"defaults" hcl:"defaults"`
}

// Supporting provides supporting configuration structures
type Supporting struct {
	// Supporting configuration
	Support map[string]interface{} `json:"support" hcl:"support"`
}

// Environment provides environment-specific configuration
type Environment struct {
	// Environment name
	Environment string `json:"environment" hcl:"environment"`

	// Environment-specific settings
	Settings map[string]interface{} `json:"settings" hcl:"settings"`
}

// Error represents a configuration-related error
type Error struct {
	// Error details
	Code        string                 `json:"code" hcl:"code"`
	Message     string                 `json:"message" hcl:"message"`
	Context     map[string]interface{} `json:"context,omitempty" hcl:"context,optional"`
	Stack       []string               `json:"stack,omitempty" hcl:"stack,optional"`
	Recoverable bool                   `json:"recoverable" hcl:"recoverable"`

	// Configuration path where the error occurred
	ConfigPath string `json:"config_path" hcl:"config_path"`

	// Configuration key where the error occurred
	ConfigKey string `json:"config_key,omitempty" hcl:"config_key,optional"`

	// Configuration value that caused the error
	ConfigValue interface{} `json:"config_value,omitempty" hcl:"config_value,optional"`
}

// NewError creates a new configuration error
func NewError(configPath, message string) *Error {
	return &Error{
		Code:        "config_error",
		Message:     message,
		Recoverable: true,
		ConfigPath:  configPath,
	}
}

// Error implements the error interface
func (e *Error) Error() string {
	return e.Message
}

// Unwrap returns the underlying error
func (e *Error) Unwrap() error {
	return nil
}
