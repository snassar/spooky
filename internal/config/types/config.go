package types

// No imports needed

// Config represents the main configuration structure
type Config struct {
	GlobalConfig  *GlobalConfig          `hcl:"global,optional"`
	ProjectConfig *ProjectConfig         `hcl:"project,optional"`
	Environment   map[string]interface{} `hcl:"environment,optional"`
	Source        ConfigSource           `hcl:"source,optional"`
	CLIFlags      map[string]interface{} `hcl:"cli_flags,optional"`
}

// GlobalConfig represents global configuration
type GlobalConfig struct {
	LogLevel   string `hcl:"log_level,optional"`
	LogFile    string `hcl:"log_file,optional"`
	Quiet      bool   `hcl:"quiet,optional"`
	Verbose    bool   `hcl:"verbose,optional"`
	ConfigPath string `hcl:"config_path,optional"`

	// Configuration blocks
	Storage     *StorageConfig     `hcl:"storage,block"`
	Facts       *FactsConfig       `hcl:"facts,block"`
	SSH         *SSHConfig         `hcl:"ssh,block"`
	Templates   *TemplatesConfig   `hcl:"templates,block"`
	Security    *SecurityConfig    `hcl:"security,block"`
	Age         *AgeConfig         `hcl:"age,block"`
	Logging     *LoggingConfig     `hcl:"logging,block"`
	Performance *PerformanceConfig `hcl:"performance,block"`
}

// ProjectConfig represents project configuration
type ProjectConfig struct {
	Name        string `hcl:"name,label" validate:"required"`
	Description string `hcl:"description,optional"`
	Version     string `hcl:"version,optional"`
	Environment string `hcl:"environment,optional"`

	// File references
	InventoryFile string `hcl:"inventory_file,optional"`
	ActionsFile   string `hcl:"actions_file,optional"`

	// Project settings
	DefaultTimeout        int  `hcl:"default_timeout,optional" validate:"omitempty,min=1,max=3600"`
	DefaultParallel       bool `hcl:"default_parallel,optional"`
	DryRunDefault         bool `hcl:"dry_run_default,optional"`
	ValidateBeforeExecute bool `hcl:"validate_before_execute,optional"`
	BackupBeforeChanges   bool `hcl:"backup_before_changes,optional"`

	// Configuration blocks
	Storage   *StorageConfig   `hcl:"storage,block"`
	Logging   *LoggingConfig   `hcl:"logging,block"`
	SSH       *SSHConfig       `hcl:"ssh,block"`
	Isolation *IsolationConfig `hcl:"isolation,block"`

	// Project-wide tags
	Tags map[string]string `hcl:"tags,optional"`
}

// Configuration types
type LoadingConfig struct {
	ConfigPath    string        `hcl:"config_path,optional"`
	DefaultConfig *GlobalConfig `hcl:"default_config,optional"`
	AutoReload    bool          `hcl:"auto_reload,optional"`
}

type ValidationConfig struct {
	ValidationRules  *ValidationRules `hcl:"validation_rules,optional"`
	StrictValidation bool             `hcl:"strict_validation,optional"`
}

type EnvironmentConfig struct {
	EnvironmentFile   string `hcl:"environment_file,optional"`
	ValidateVariables bool   `hcl:"validate_variables,optional"`
}

type ValidationRules struct {
	RequiredFields  []string `hcl:"required_fields,optional"`
	ForbiddenFields []string `hcl:"forbidden_fields,optional"`
}

// ConfigSource represents the source of configuration
type ConfigSource string

const (
	SourceDefault     ConfigSource = "default"
	SourceGlobal      ConfigSource = "global"
	SourceProject     ConfigSource = "project"
	SourceEnvironment ConfigSource = "environment"
	SourceCLI         ConfigSource = "cli"
)

// Error types
type ValidationError struct {
	Field   string `hcl:"field"`
	Message string `hcl:"message"`
}

type ConfigError struct {
	Operation string `hcl:"operation"`
	Error     string `hcl:"error"`
}
