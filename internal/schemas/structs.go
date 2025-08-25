package schemas

import "time"

// SchemaVersion represents the current supported schema version
const SchemaVersion = "1"

// SupportedVersions lists all supported schema versions (current + 4 previous)
var SupportedVersions = []string{"1"}

// ============================================================================
// PROJECT SCHEMA STRUCTS
// ============================================================================

// ProjectV1 represents a project configuration (version 1)
type ProjectV1 struct {
	// Project metadata
	Name        string `json:"name" required:"true" pattern:"^[a-zA-Z][a-zA-Z0-9._-]*$" min_length:"1" max_length:"128" description:"Project name (used for identification and isolation)"`
	Description string `json:"description" max_length:"1024" description:"Project description"`

	// Execution behavior configuration
	RunMaxParallel         int  `json:"run_max_parallel" min:"1" max:"100" default:"10" description:"Maximum parallel action executions for this project"`
	RunDryRunDefault       bool `json:"run_dry_run_default" default:"false" description:"Default dry-run mode for actions"`
	RunValidateBeforeRun   bool `json:"run_validate_before_run" default:"true" description:"Validate project configuration before running"`
	RunBackupBeforeChanges bool `json:"run_backup_before_changes" default:"false" description:"Create backups before making changes"`

	// Facts collection configuration
	FactsTimeout            int `json:"facts_timeout" min:"1" max:"3600" default:"30" description:"Timeout for facts collection in seconds"`
	FactsParallelCollection int `json:"facts_parallel_collection" min:"1" max:"100" default:"10" description:"Number of parallel facts collection workers for this project"`
	FactsRetryAttempts      int `json:"facts_retry_attempts" min:"0" max:"10" default:"3" description:"Number of retry attempts for failed facts collection"`
	FactsRetryDelay         int `json:"facts_retry_delay" min:"1" max:"60" default:"5" description:"Delay between retry attempts in seconds"`
}

// ============================================================================
// MACHINES SCHEMA STRUCTS
// ============================================================================

// MachinesV1 represents a machines configuration (version 1)
// In HCL2, this becomes: machines { machine "name" { ... } }
type MachinesV1 struct {
	Machine []MachinesMachineV1 `json:"machine" description:"Individual machine configurations"`
	Group   []MachinesGroupV1   `json:"group" description:"Machine group configurations"`
}

// MachinesMachineV1 represents an individual machine configuration
// In HCL2, this becomes: machine "name" { ... } where "name" is the block label
type MachinesMachineV1 struct {
	// Machine identification (name comes from block label)
	Description string `json:"description" max_length:"256" description:"Description of the machine"`

	// Connection details
	Hostname string `json:"hostname" required:"true" min_length:"1" max_length:"253" description:"Hostname or IP address of the machine"`
	Port     int    `json:"port" min:"1" max:"65535" default:"22" description:"SSH port number"`
	User     string `json:"user" required:"true" pattern:"^[a-zA-Z0-9_.-]+$" min_length:"1" max_length:"32" description:"SSH username"`

	// Authentication configuration
	Authentication MachinesMachineAuthenticationV1 `json:"authentication" required:"true" description:"SSH authentication configuration"`

	// Connection settings
	ConnectionTimeout int `json:"connection_timeout" min:"1" max:"3600" default:"30" description:"Connection timeout in seconds"`
	MaxRetries        int `json:"max_retries" min:"0" max:"10" default:"3" description:"Maximum connection retry attempts"`
	RetryDelay        int `json:"retry_delay" min:"1" max:"60" default:"5" description:"Delay between retry attempts in seconds"`

	// Machine facts (optional)
	Facts MachinesMachineFactsV1 `json:"facts" description:"Machine facts and metadata"`
}

// MachinesMachineAuthenticationV1 represents SSH authentication configuration
// In HCL2, this becomes: authentication "method" { ... } where "method" is the block label
// Supported methods: publickey, password, certificate
type MachinesMachineAuthenticationV1 struct {
	// Authentication method comes from block label: authentication "publickey" { ... }
	// SSH Key authentication fields (for publickey method)
	PublicKeyPath string                                    `json:"public_key_path" description:"Path to SSH private key file (required for publickey method)"`
	Passphrase    MachinesMachineAuthenticationPassphraseV1 `json:"passphrase" description:"SSH key passphrase - can be encrypted using age encryption"`

	// Password authentication fields (for password method)
	Password MachinesMachineAuthenticationPasswordV1 `json:"password" description:"SSH password - can be encrypted using age encryption"`

	// Certificate authentication fields (for certificate method)
	// Certificate-based auth requires both private key and certificate
	PrivateKeyPath        string                                    `json:"private_key_path" description:"Path to SSH private key file (required for certificate method)"`
	CertificatePath       string                                    `json:"certificate_path" description:"Path to SSH certificate file (required for certificate method)"`
	CertificatePassphrase MachinesMachineAuthenticationPassphraseV1 `json:"certificate_passphrase" description:"SSH certificate passphrase if the certificate is encrypted"`
}

// MachinesMachineAuthenticationPassphraseV1 represents SSH key passphrase configuration
type MachinesMachineAuthenticationPassphraseV1 struct {
	Value          string            `json:"value" description:"SSH key passphrase value (mutually exclusive with encrypted_value)"`
	Encrypted      bool              `json:"encrypted" default:"false" description:"Whether the passphrase should be encrypted (triggers transformation)"`
	EncryptedValue *EncryptedValueV1 `json:"encrypted_value,omitempty" description:"Structured encrypted value (mutually exclusive with value)"`
}

// MachinesMachineAuthenticationPasswordV1 represents SSH password configuration
type MachinesMachineAuthenticationPasswordV1 struct {
	Value          string            `json:"value" description:"SSH password value (mutually exclusive with encrypted_value)"`
	Encrypted      bool              `json:"encrypted" default:"false" description:"Whether the password should be encrypted (triggers transformation)"`
	EncryptedValue *EncryptedValueV1 `json:"encrypted_value,omitempty" description:"Structured encrypted value (mutually exclusive with value)"`
}

// MachinesMachineFactsV1 represents machine facts and metadata
type MachinesMachineFactsV1 struct {
	OSFamily     string `json:"os_family" description:"Operating system family (e.g., linux, windows, darwin)"`
	OSVersion    string `json:"os_version" description:"Operating system version"`
	Architecture string `json:"architecture" description:"System architecture (e.g., x86_64, arm64)"`
	Hostname     string `json:"hostname" description:"System hostname"`
	Domain       string `json:"domain" description:"System domain"`
	Location     string `json:"location" description:"Physical or logical location"`
	Environment  string `json:"environment" description:"Environment (e.g., dev, staging, prod)"`
	Role         string `json:"role" description:"Machine role (e.g., web, database, load_balancer)"`
}

// MachinesGroupV1 represents a machine group configuration
// In HCL2, this becomes: group "name" { ... } where "name" is the block label
type MachinesGroupV1 struct {
	// Group identification (name comes from block label)
	Description string   `json:"description" max_length:"256" description:"Description of the group"`
	Machines    []string `json:"machines" required:"true" min_items:"1" description:"List of machine names in this group"`

	// Group-level settings (inherited by machines in the group)
	User           string                          `json:"user" description:"Default SSH username for machines in this group"`
	Port           int                             `json:"port" min:"1" max:"65535" default:"22" description:"Default SSH port for machines in this group"`
	Authentication MachinesMachineAuthenticationV1 `json:"authentication" description:"Default authentication configuration for machines in this group"`
}

// ============================================================================
// ACTIONS SCHEMA STRUCTS
// ============================================================================

// ActionsV1 represents an actions configuration (version 1)
// In HCL2, this becomes: actions { action "name" { ... } }
type ActionsV1 struct {
	Action []ActionsActionV1 `json:"action" description:"Individual action configurations"`
}

// ActionsActionV1 represents an individual action configuration
// In HCL2, this becomes: action "name" { ... } where "name" is the block label
type ActionsActionV1 struct {
	// Action identification (name comes from block label)
	Description string   `json:"description" required:"true" min_length:"1" max_length:"500" description:"Action description"`
	Type        string   `json:"type" required:"true" enum:"command,script,template_deploy,file_sync,service_control" description:"Action run type"`
	Tags        []string `json:"tags" description:"Tags for categorizing and filtering actions"`

	// Command type fields
	Command string `json:"command" min_length:"1" max_length:"1000" description:"Command to run (for command type)"`

	// Script type fields
	Script    string            `json:"script" pattern:"^(files|templates)/[a-zA-Z0-9/._-]+(\\.sh|\\.tmpl)?$" description:"Script file path in files/ or templates/ directory (for script type)"`
	Variables map[string]string `json:"variables" description:"Variables for templated scripts (for script type with .tmpl files)"`

	// Template deploy type fields
	Source      string `json:"source" pattern:"^templates/[a-zA-Z0-9/._-]+\\.tmpl$" description:"Source template file path (for template_deploy type)"`
	Destination string `json:"destination" pattern:"^[a-zA-Z0-9/._-]+$" description:"Destination file path (local or remote) (for template_deploy type)"`
	Validate    bool   `json:"validate" default:"false" description:"Validate template syntax before deployment (for template_deploy type)"`
	Backup      bool   `json:"backup" default:"false" description:"Create backup of existing file before overwriting (for template_deploy type)"`
	Permissions string `json:"permissions" pattern:"^[0-7]{3,4}$" description:"File permissions in octal format (for template_deploy type)"`
	Owner       string `json:"owner" pattern:"^[a-zA-Z0-9_.-]+$" description:"File owner (for template_deploy type)"`
	Group       string `json:"group" pattern:"^[a-zA-Z0-9_.-]+$" description:"File group (for template_deploy type)"`

	// File sync type fields
	SyncSource      string `json:"sync_source" required:"true" description:"Source directory or file path (for file_sync type)"`
	SyncDestination string `json:"sync_destination" required:"true" description:"Destination directory or file path (for file_sync type)"`
	SyncDelete      bool   `json:"sync_delete" default:"false" description:"Delete files in destination that don't exist in source (for file_sync type)"`
	SyncPreserve    bool   `json:"sync_preserve" default:"true" description:"Preserve file attributes (for file_sync type)"`

	// Service control type fields
	ServiceName   string `json:"service_name" required:"true" description:"Name of the service to control (for service_control type)"`
	ServiceAction string `json:"service_action" required:"true" enum:"start,stop,restart,reload,enable,disable,status" description:"Action to perform on the service (for service_control type)"`

	// Execution configuration
	Targets    []string `json:"targets" description:"Target machines or groups for this action"`
	RunAs      string   `json:"run_as" description:"User to run the action as"`
	Sudo       bool     `json:"sudo" default:"false" description:"Whether to use sudo for execution"`
	Timeout    int      `json:"timeout" min:"1" max:"3600" default:"300" description:"Action timeout in seconds"`
	Retries    int      `json:"retries" min:"0" max:"10" default:"0" description:"Number of retry attempts on failure"`
	RetryDelay int      `json:"retry_delay" min:"1" max:"60" default:"5" description:"Delay between retry attempts in seconds"`

	// Conditional execution
	When   string `json:"when" description:"Condition for when to execute this action"`
	OnlyIf string `json:"only_if" description:"Command to check before execution (action only runs if command succeeds)"`
	Unless string `json:"unless" description:"Command to check before execution (action only runs if command fails)"`

	// Output handling
	ExpectExitCode int    `json:"expect_exit_code" description:"Expected exit code for successful execution"`
	OutputFile     string `json:"output_file" description:"File to capture command output"`
	ErrorFile      string `json:"error_file" description:"File to capture error output"`
	Quiet          bool   `json:"quiet" default:"false" description:"Suppress output unless there's an error"`
}

// ============================================================================
// VARIABLES SCHEMA STRUCTS
// ============================================================================

// VariablesV1 represents a variables configuration (version 1)
// In HCL2, this becomes: variables { variable "name" { ... } }
type VariablesV1 struct {
	Variable []VariablesVariableV1 `json:"variable" description:"Individual variable configurations"`
}

// VariablesVariableV1 represents an individual variable configuration
// In HCL2, this becomes: variable "name" { ... } where "name" is the block label
type VariablesVariableV1 struct {
	// Variable identification (name comes from block label)
	Type           string            `json:"type" enum:"string,number,boolean,list,map" description:"Variable type"`
	Value          string            `json:"value" description:"Variable value (mutually exclusive with encrypted_value)"`
	Description    string            `json:"description" max_length:"256" description:"Variable description"`
	Sensitive      bool              `json:"sensitive" default:"false" description:"Whether this variable contains sensitive information (auto-true if encrypted_value exists)"`
	Encrypted      bool              `json:"encrypted" default:"false" description:"Whether this variable should be encrypted (triggers transformation)"`
	EncryptedValue *EncryptedValueV1 `json:"encrypted_value,omitempty" description:"Structured encrypted value (mutually exclusive with value)"`

	// Validation rules
	Required  bool   `json:"required" default:"false" description:"Whether this variable is required"`
	Default   string `json:"default" description:"Default value for this variable"`
	Pattern   string `json:"pattern" description:"Regex pattern for value validation"`
	MinLength int    `json:"min_length" min:"0" description:"Minimum length for string values"`
	MaxLength int    `json:"max_length" min:"1" description:"Maximum length for string values"`
	Min       int    `json:"min" description:"Minimum value for numeric values"`
	Max       int    `json:"max" description:"Maximum value for numeric values"`
	Enum      string `json:"enum" description:"Comma-separated list of allowed values"`
}

// ============================================================================
// COMMON SCHEMA STRUCTS
// ============================================================================

// EncryptedValueV1 represents a structured encrypted value
type EncryptedValueV1 struct {
	Data        string `json:"data" required:"true" description:"The encrypted data (base64 content only, no headers/footers)"`
	Format      string `json:"format" enum:"base64,armored,compact" default:"base64" description:"Format of the encrypted data"`
	Algorithm   string `json:"algorithm" default:"age" description:"Encryption algorithm used"`
	Version     string `json:"version" default:"v1" description:"Encryption version"`
	EncryptedAt string `json:"encrypted_at" description:"ISO 8601 timestamp when the value was encrypted"`
}

// MetadataSchemaV1 represents schema metadata
type MetadataSchemaV1 struct {
	Version     string `json:"version" required:"true" description:"Schema version"`
	Description string `json:"description" required:"true" description:"Schema description"`
}

// ValidationRuleV1 represents a validation rule
type ValidationRuleV1 struct {
	Pattern   string `json:"pattern" description:"Regex pattern for validation"`
	Message   string `json:"message" description:"Validation error message"`
	Min       int    `json:"min" description:"Minimum value"`
	Max       int    `json:"max" description:"Maximum value"`
	MinLength int    `json:"min_length" description:"Minimum length"`
	MaxLength int    `json:"max_length" description:"Maximum length"`
	Required  bool   `json:"required" description:"Whether field is required"`
	Default   string `json:"default" description:"Default value"`
	Enum      string `json:"enum" description:"Comma-separated allowed values"`
}

// ============================================================================
// PROJECT DIRECTORY SCHEMA STRUCTS
// ============================================================================

// ProjectDirectoryV1 represents project directory structure validation (version 1)
type ProjectDirectoryV1 struct {
	Name       string                   `json:"name" required:"true" description:"Project root directory name"`
	File       []ProjectDirectoryFileV1 `json:"file" description:"Required and optional files"`
	Directory  []ProjectDirectoryDirV1  `json:"directory" description:"Required and optional directories"`
	Validation []ProjectDirectoryRuleV1 `json:"validation" description:"Directory validation rules"`
}

// ProjectDirectoryFileV1 represents a file requirement in project directory
type ProjectDirectoryFileV1 struct {
	Name        string `json:"name" required:"true" description:"File name"`
	Type        string `json:"type" required:"true" enum:"file" description:"Type (always 'file')"`
	Required    bool   `json:"required" default:"false" description:"Whether this file is required"`
	Description string `json:"description" description:"Description of the file's purpose"`
	Validate    string `json:"validate" description:"Validation rule to apply to this file"`
	Pattern     string `json:"pattern" description:"Regex pattern for file content validation"`
}

// ProjectDirectoryDirV1 represents a directory requirement in project directory
type ProjectDirectoryDirV1 struct {
	Name        string `json:"name" required:"true" description:"Directory name"`
	Type        string `json:"type" required:"true" enum:"directory" description:"Type (always 'directory')"`
	Required    bool   `json:"required" default:"false" description:"Whether this directory is required"`
	Description string `json:"description" description:"Description of the directory's purpose"`
	Validate    string `json:"validate" description:"Validation rule to apply to this directory"`
	Pattern     string `json:"pattern" description:"Regex pattern for directory content validation"`
}

// ProjectDirectoryRuleV1 represents a directory validation rule
type ProjectDirectoryRuleV1 struct {
	Name        string `json:"name" required:"true" description:"Rule name"`
	Description string `json:"description" description:"Rule description"`
	Type        string `json:"type" required:"true" enum:"file_exists,directory_exists,hcl_config,pattern_match" description:"Validation rule type"`
	Pattern     string `json:"pattern" description:"Regex pattern for validation"`
	Message     string `json:"message" description:"Validation error message"`
}

// ============================================================================
// SPOOKY GLOBAL CONFIGURATION SCHEMA STRUCTS
// ============================================================================

// SpookyV1 represents global spooky CLI configuration (version 1)
type SpookyV1 struct {
	SSH      SpookySSHV1      `json:"ssh" description:"SSH configuration"`
	Security SpookySecurityV1 `json:"security" description:"Security configuration"`
	Age      SpookyAgeV1      `json:"age" description:"Age encryption configuration"`
	Logging  SpookyLoggingV1  `json:"logging" description:"Logging configuration"`
}

// SpookySSHV1 represents SSH configuration in global spooky config
type SpookySSHV1 struct {
	Timeout            int  `json:"timeout" min:"1" max:"300" default:"30" description:"SSH connection timeout in seconds"`
	KeepaliveInterval  int  `json:"keepalive_interval" min:"1" max:"300" default:"60" description:"SSH keepalive interval in seconds"`
	KeepaliveCount     int  `json:"keepalive_count" min:"1" max:"10" default:"3" description:"SSH keepalive count before considering connection dead"`
	KeyScanTimeout     int  `json:"key_scan_timeout" min:"1" max:"60" default:"10" description:"SSH key scanning timeout in seconds"`
	KnownHostsStrict   bool `json:"known_hosts_strict" default:"true" description:"Strict known_hosts checking"`
	ConnectionPoolSize int  `json:"connection_pool_size" min:"1" max:"100" default:"10" description:"SSH connection pool size"`

	// Proxy configuration
	ProxyCommand string `json:"proxy_command" description:"SSH proxy command (e.g., 'ssh -W %h:%p bastion.example.com')"`
	ProxyJump    string `json:"proxy_jump" description:"SSH proxy jump host (e.g., 'bastion.example.com')"`

	// Compression configuration
	Compression      bool `json:"compression" default:"false" description:"Enable SSH compression"`
	CompressionLevel int  `json:"compression_level" min:"1" max:"9" default:"6" description:"SSH compression level (1-9)"`

	// TCP keepalive configuration
	TCPKeepAlive              bool          `json:"tcp_keepalive" default:"true" description:"Enable TCP keepalive"`
	TCPKeepAliveCount         int           `json:"tcp_keepalive_count" min:"1" max:"10" default:"3" description:"TCP keepalive count"`
	TCPKeepAliveIdle          time.Duration `json:"tcp_keepalive_idle" default:"60s" description:"TCP keepalive idle time"`
	TCPKeepAliveInterval      time.Duration `json:"tcp_keepalive_interval" default:"10s" description:"TCP keepalive interval"`
	TCPKeepAliveProbeInterval time.Duration `json:"tcp_keepalive_probe_interval" default:"5s" description:"TCP keepalive probe interval"`
}

// SpookySecurityV1 represents security configuration in global spooky config
type SpookySecurityV1 struct {
	AllowUnsafeCommands bool `json:"allow_unsafe_commands" default:"false" description:"Allow potentially unsafe commands"`
	AuditLogging        bool `json:"audit_logging" default:"true" description:"Enable audit logging for security events"`
}

// SpookyAgeV1 represents age encryption configuration in global spooky config
type SpookyAgeV1 struct {
	Identities string `json:"identities" description:"Path to directory containing age identity files"`
	Recipients string `json:"recipients" description:"Path to file containing age recipients (public keys, one per line)"`
}

// SpookyLoggingV1 represents logging configuration in global spooky config
type SpookyLoggingV1 struct {
	Level      string `json:"level" default:"info" enum:"debug,info,warn,error,fatal" description:"Log level"`
	Format     string `json:"format" default:"json" enum:"json,text,structured" description:"Log format"`
	Output     string `json:"output" default:"stderr" enum:"stdout,stderr,file,null" description:"Log output destination"`
	FilePath   string `json:"file_path" description:"Path to log file (required when output is 'file')"`
	FilePerms  string `json:"file_permissions" default:"0644" pattern:"^[0-7]{3,4}$" description:"File permissions in octal format"`
	FileAppend bool   `json:"file_append" default:"true" description:"Whether to append to existing file or truncate"`
}

// ============================================================================
// FACTS SCHEMA STRUCTS
// ============================================================================

// FactsV1 represents facts configuration (version 1)
// Based on documentation/schemafiles/structure/facts.hcl
type FactsV1 struct {
	Name               string                     `json:"name" required:"true" pattern:"^[a-zA-Z0-9_.-]+$" min_length:"1" max_length:"128" description:"Name of the fact"`
	Value              interface{}                `json:"value" required:"true" description:"Value of the fact - can be string, number, boolean, object, or age-encrypted string"`
	Type               string                     `json:"type" required:"true" enum:"string,number,boolean,object,array,encrypted" description:"Type of the fact value"`
	Encrypted          bool                       `json:"encrypted" default:"false" description:"Whether the fact value is age-encrypted"`
	EncryptionMetadata *FactsEncryptionMetadataV1 `json:"encryption_metadata,omitempty" description:"Age encryption metadata - only present when encrypted = true"`
	Description        string                     `json:"description" max_length:"256" description:"Description of the fact"`
	Tags               []string                   `json:"tags" max_items:"10" description:"Tags for categorizing facts"`
	Metadata           *FactsMetadataV1           `json:"metadata,omitempty" description:"Additional metadata for the fact"`
}

// FactsEncryptionMetadataV1 represents age encryption metadata for facts
type FactsEncryptionMetadataV1 struct {
	Recipients  []string `json:"recipients" required:"true" description:"List of age public keys that can decrypt this value"`
	EncryptedAt string   `json:"encrypted_at" description:"ISO 8601 timestamp when the value was encrypted"`
	Method      string   `json:"method" default:"age" description:"Encryption method used"`
}

// FactsMetadataV1 represents additional metadata for facts
type FactsMetadataV1 struct {
	Source     string `json:"source" description:"Source of the fact (e.g., 'environment', 'file', 'manual')"`
	CreatedAt  string `json:"created_at" description:"ISO 8601 timestamp when the fact was created"`
	ModifiedAt string `json:"modified_at" description:"ISO 8601 timestamp when the fact was last modified"`
	Version    string `json:"version" description:"Version identifier for the fact"`
}

// ============================================================================
// LOGGING SCHEMA STRUCTS
// ============================================================================

// LoggingV1 represents logging configuration (version 1)
// Based on documentation/schemafiles/structure/logging.hcl
type LoggingV1 struct {
	// Basic configuration
	Level  string `json:"level" default:"info" enum:"debug,info,warn,error,fatal" description:"Minimum log level to output"`
	Format string `json:"format" default:"json" enum:"json,text,structured" description:"Log output format"`
	Output string `json:"output" default:"stderr" enum:"stdout,stderr,file,null" description:"Log output destination"`

	// File output configuration
	FilePath        string `json:"file_path" description:"Path to log file (required when output is 'file')"`
	FilePermissions string `json:"file_permissions" default:"0644" pattern:"^[0-7]{3,4}$" description:"File permissions in octal format"`
	FileAppend      bool   `json:"file_append" default:"true" description:"Whether to append to existing file or truncate"`

	// Structured logging configuration
	StructuredTimestampEnabled  bool   `json:"structured_timestamp_enabled" default:"true" description:"Whether to include timestamps in log entries"`
	StructuredTimestampFormat   string `json:"structured_timestamp_format" default:"RFC3339" enum:"RFC3339,RFC3339Nano,Unix,UnixNano,ISO8601" description:"Timestamp format"`
	StructuredTimestampTimezone string `json:"structured_timestamp_timezone" default:"UTC" description:"Timezone for timestamps"`
	StructuredLevelKey          string `json:"structured_level_key" default:"level" description:"Field key for log level"`
	StructuredMessageKey        string `json:"structured_message_key" default:"message" description:"Field key for log message"`
	StructuredErrorKey          string `json:"structured_error_key" default:"error" description:"Field key for error information"`

	// Performance configuration
	PerformanceBufferEnabled       bool   `json:"performance_buffer_enabled" default:"false" description:"Whether to use buffered logging"`
	PerformanceBufferSize          int    `json:"performance_buffer_size" default:"4096" min:"1024" max:"1048576" description:"Buffer size in bytes"`
	PerformanceBufferFlushInterval string `json:"performance_buffer_flush_interval" default:"1s" description:"Flush interval"`
	PerformanceAsyncEnabled        bool   `json:"performance_async_enabled" default:"false" description:"Whether to use asynchronous logging"`
	PerformanceAsyncQueueSize      int    `json:"performance_async_queue_size" default:"1000" min:"100" max:"100000" description:"Queue size for async logging"`
	PerformanceAsyncWorkers        int    `json:"performance_async_workers" default:"1" min:"1" max:"10" description:"Number of worker goroutines for async logging"`
	PerformanceAsyncDropWhenFull   bool   `json:"performance_async_drop_when_full" default:"false" description:"Whether to drop logs when queue is full"`

	// Rotation configuration
	RotationEnabled    bool   `json:"rotation_enabled" default:"false" description:"Whether to enable log file rotation"`
	RotationMaxSize    string `json:"rotation_max_size" default:"100MB" description:"Maximum file size before rotation"`
	RotationMaxAge     string `json:"rotation_max_age" default:"30d" description:"Maximum age of rotated files"`
	RotationMaxBackups int    `json:"rotation_max_backups" default:"5" min:"1" max:"100" description:"Maximum number of backup files to keep"`
	RotationCompress   bool   `json:"rotation_compress" default:"true" description:"Whether to compress rotated log files"`
	RotationLocalTime  bool   `json:"rotation_local_time" default:"false" description:"Whether to use local time for rotation timestamps"`
}
