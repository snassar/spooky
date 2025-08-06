package types

// StorageConfig represents storage configuration
type StorageConfig struct {
	Type            string `hcl:"type,optional" validate:"omitempty,oneof=badgerdb json"`
	Path            string `hcl:"path,optional"`
	Compression     bool   `hcl:"compression,optional"`
	Encryption      bool   `hcl:"encryption,optional"`
	BackupEnabled   bool   `hcl:"backup_enabled,optional"`
	BackupRetention int    `hcl:"backup_retention,optional" validate:"omitempty,min=1,max=365"`
}

// FactsConfig represents facts collection configuration
type FactsConfig struct {
	Timeout            int  `hcl:"timeout,optional" validate:"omitempty,min=1,max=3600"`
	CacheTTL           int  `hcl:"cache_ttl,optional" validate:"omitempty,min=0,max=86400"`
	AutoCollect        bool `hcl:"auto_collect,optional"`
	ParallelCollection int  `hcl:"parallel_collection,optional" validate:"omitempty,min=1,max=100"`
	RetryAttempts      int  `hcl:"retry_attempts,optional" validate:"omitempty,min=0,max=10"`
	RetryDelay         int  `hcl:"retry_delay,optional" validate:"omitempty,min=1,max=60"`
}

// SSHConfig represents SSH configuration
type SSHConfig struct {
	Timeout            int  `hcl:"timeout,optional" validate:"omitempty,min=1,max=300"`
	KeepaliveInterval  int  `hcl:"keepalive_interval,optional" validate:"omitempty,min=1,max=300"`
	KeepaliveCount     int  `hcl:"keepalive_count,optional" validate:"omitempty,min=1,max=10"`
	KeyScanTimeout     int  `hcl:"key_scan_timeout,optional" validate:"omitempty,min=1,max=60"`
	KnownHostsStrict   bool `hcl:"known_hosts_strict,optional"`
	ConnectionPoolSize int  `hcl:"connection_pool_size,optional" validate:"omitempty,min=1,max=100"`
}

// TemplatesConfig represents template configuration
type TemplatesConfig struct {
	MaxSize                int      `hcl:"max_size,optional" validate:"omitempty,min=1024,max=10485760"`
	AllowExternalFunctions bool     `hcl:"allow_external_functions,optional"`
	Timeout                int      `hcl:"timeout,optional" validate:"omitempty,min=1,max=300"`
	CacheCompiled          bool     `hcl:"cache_compiled,optional"`
	SandboxMode            bool     `hcl:"sandbox_mode,optional"`
	AllowedFunctions       []string `hcl:"allowed_functions,optional"`
}

// SecurityConfig represents security configuration
type SecurityConfig struct {
	AllowUnsafeCommands bool     `hcl:"allow_unsafe_commands,optional"`
	RestrictFileAccess  bool     `hcl:"restrict_file_access,optional"`
	ValidateSSHKeys     bool     `hcl:"validate_ssh_keys,optional"`
	AuditLogging        bool     `hcl:"audit_logging,optional"`
	AllowedHosts        []string `hcl:"allowed_hosts,optional"`
	BlockedHosts        []string `hcl:"blocked_hosts,optional"`
}

// AgeConfig represents age encryption configuration
type AgeConfig struct {
	Enabled        bool     `hcl:"enabled,optional"`
	PublicKey      string   `hcl:"public_key,optional"`
	PrivateKeyPath string   `hcl:"private_key_path,optional"`
	Recipients     []string `hcl:"recipients,optional"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level        string `hcl:"level,optional" validate:"omitempty,oneof=debug info warn error"`
	Format       string `hcl:"format,optional" validate:"omitempty,oneof=json text"`
	Output       string `hcl:"output,optional"`
	ColorOutput  bool   `hcl:"color_output,optional"`
	ProgressBars bool   `hcl:"progress_bars,optional"`
}

// PerformanceConfig represents performance configuration
type PerformanceConfig struct {
	DefaultParallel int `hcl:"default_parallel,optional" validate:"omitempty,min=2,max=100"`
	MaxMemory       int `hcl:"max_memory,optional" validate:"omitempty,min=64,max=8192"`
	GCInterval      int `hcl:"gc_interval,optional" validate:"omitempty,min=1,max=3600"`
}

// IsolationConfig represents project isolation configuration
type IsolationConfig struct {
	Enabled         bool     `hcl:"enabled,optional"`
	FactsScope      string   `hcl:"facts_scope,optional"`
	VariablesScope  string   `hcl:"variables_scope,optional"`
	MachineAccess   string   `hcl:"machine_access,optional"`
	AllowedMachines []string `hcl:"allowed_machines,optional"`
	AllowedTags     []string `hcl:"allowed_tags,optional"`
}

// Machine represents a remote machine configuration
type Machine struct {
	Name          string            `hcl:"name,label" validate:"required"`
	Host          string            `hcl:"host" validate:"required"`
	Port          int               `hcl:"port,optional" validate:"omitempty,min=1,max=65535"`
	User          string            `hcl:"user" validate:"required"`
	Password      string            `hcl:"password,optional"`
	KeyFile       string            `hcl:"key_file,optional"`
	KeyPassphrase string            `hcl:"key_passphrase,optional"`
	Tags          map[string]string `hcl:"tags,optional" validate:"omitempty,dive,keys,required,endkeys,required"`
	Groups        []string          `hcl:"groups,optional" validate:"omitempty,dive,required"`
	Metadata      map[string]string `hcl:"metadata,optional" validate:"omitempty,dive,keys,required,endkeys,required"`

	// SSH configuration overrides
	ConnectionTimeout int `hcl:"connection_timeout,optional" validate:"omitempty,min=1,max=300"`
	CommandTimeout    int `hcl:"command_timeout,optional" validate:"omitempty,min=1,max=3600"`
	MaxConnections    int `hcl:"max_connections,optional" validate:"omitempty,min=1,max=100"`
	RetryAttempts     int `hcl:"retry_attempts,optional" validate:"omitempty,min=0,max=10"`
	RetryDelay        int `hcl:"retry_delay,optional" validate:"omitempty,min=1,max=60"`
}

// Action represents an action to be executed on machines
type Action struct {
	Name        string          `hcl:"name,label" validate:"required"`
	Description string          `hcl:"description,optional"`
	Type        string          `hcl:"type,optional" validate:"omitempty,oneof=command script template_deploy template_evaluate template_validate template_cleanup"`
	Command     string          `hcl:"command,optional"`
	Script      string          `hcl:"script,optional"`
	Template    *TemplateConfig `hcl:"template,block"`
	Machines    []string        `hcl:"machines,optional" validate:"omitempty,dive,required"`
	Tags        []string        `hcl:"tags,optional" validate:"omitempty,dive,required"`
	Timeout     int             `hcl:"timeout,optional" validate:"omitempty,min=1,max=3600"`
	Parallel    bool            `hcl:"parallel,optional"`

	// Extended properties for retries and error handling
	Retries    int `hcl:"retries,optional" validate:"omitempty,min=0,max=10"`
	RetryDelay int `hcl:"retry_delay,optional" validate:"omitempty,min=1,max=300"`

	// Dependencies and execution control
	Dependencies []string `hcl:"dependencies,optional" validate:"omitempty,dive,required"`

	// Environment and execution context
	Environment      map[string]string `hcl:"environment,optional" validate:"omitempty,dive,keys,required,endkeys,required"`
	WorkingDirectory string            `hcl:"working_directory,optional"`
	User             string            `hcl:"user,optional"`
	Sudo             bool              `hcl:"sudo,optional"`
	DryRun           bool              `hcl:"dry_run,optional"`

	// Action metadata and organization
	Category string            `hcl:"category,optional"`
	Priority int               `hcl:"priority,optional" validate:"omitempty,min=1,max=10"`
	Critical bool              `hcl:"critical,optional"`
	Metadata map[string]string `hcl:"metadata,optional" validate:"omitempty,dive,keys,required,endkeys,required"`

	// Security and validation
	ValidateBeforeRun bool `hcl:"validate_before_run,optional"`
	AllowFailure      bool `hcl:"allow_failure,optional"`

	// Performance and resource management
	MaxConcurrent  int                   `hcl:"max_concurrent,optional" validate:"omitempty,min=1,max=100"`
	ResourceLimits *ActionResourceLimits `hcl:"resource_limits,block"`
}

// ActionResourceLimits represents resource limits for action execution
type ActionResourceLimits struct {
	MemoryMB   int `hcl:"memory_mb,optional" validate:"omitempty,min=1,max=32768"`
	CPUPercent int `hcl:"cpu_percent,optional" validate:"omitempty,min=1,max=100"`
	DiskMB     int `hcl:"disk_mb,optional" validate:"omitempty,min=1,max=1048576"`
}

// TemplateConfig represents template-specific configuration
type TemplateConfig struct {
	Source      string `hcl:"source" validate:"required"`
	Destination string `hcl:"destination" validate:"required"`
	Validate    bool   `hcl:"validate,optional"`
	Backup      bool   `hcl:"backup,optional"`
	Permissions string `hcl:"permissions,optional" validate:"omitempty,regexp=^[0-7]{3,4}$"`
	Owner       string `hcl:"owner,optional"`
	Group       string `hcl:"group,optional"`
}

// InventoryConfig represents an inventory configuration (machines only)
type InventoryConfig struct {
	Machines []Machine `hcl:"machine,block" validate:"required,min=1,dive"`
}

// ActionsConfig represents an actions configuration (actions only)
type ActionsConfig struct {
	Actions []Action `hcl:"action,block" validate:"dive"`
}

// Wrapper types for HCL parsing
type ProjectConfigWrapper struct {
	Project *ProjectConfig `hcl:"project,block"`
}

type InventoryWrapper struct {
	Inventory *InventoryConfig `hcl:"inventory,block"`
}

type MachinesWrapper struct {
	Machines *InventoryConfig `hcl:"machines,block"`
}

type ActionsWrapper struct {
	Actions *ActionsConfig `hcl:"actions,block"`
}

type GlobalConfigWrapper struct {
	Spooky *GlobalConfig `hcl:"spooky,block"`
}
