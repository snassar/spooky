package config

// Config represents the main configuration structure (legacy combined format)
type Config struct {
	Machines []Machine `hcl:"machine,block" validate:"required,min=1,dive"`
	Actions  []Action  `hcl:"action,block" validate:"dive"`
}

// ProjectConfig represents a project configuration
type ProjectConfig struct {
	Name        string `hcl:"name,label" validate:"required"`
	Description string `hcl:"description,optional"`
	Version     string `hcl:"version,optional"`
	Environment string `hcl:"environment,optional"`

	// File references
	InventoryFile string `hcl:"inventory_file,optional"`
	ActionsFile   string `hcl:"actions_file,optional"`

	// Project settings
	DefaultTimeout  int  `hcl:"default_timeout,optional" validate:"omitempty,min=1,max=3600"`
	DefaultParallel bool `hcl:"default_parallel,optional"`

	// Configuration blocks
	Storage *StorageConfig `hcl:"storage,block"`
	Logging *LoggingConfig `hcl:"logging,block"`
	SSH     *SSHConfig     `hcl:"ssh,block"`

	// Project-wide tags
	Tags map[string]string `hcl:"tags,optional"`
}

// ProjectConfigWrapper wraps ProjectConfig for HCL parsing
type ProjectConfigWrapper struct {
	Project *ProjectConfig `hcl:"project,block"`
}

// InventoryWrapper wraps InventoryConfig for HCL parsing
type InventoryWrapper struct {
	Inventory *InventoryConfig `hcl:"inventory,block"`
}

// ActionsWrapper wraps ActionsConfig for HCL parsing
type ActionsWrapper struct {
	Actions *ActionsConfig `hcl:"actions,block"`
}

// StorageConfig represents storage configuration
type StorageConfig struct {
	Type string `hcl:"type" validate:"required,oneof=badgerdb json"`
	Path string `hcl:"path" validate:"required"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level  string `hcl:"level,optional" validate:"omitempty,oneof=debug info warn error"`
	Format string `hcl:"format,optional" validate:"omitempty,oneof=json text"`
	Output string `hcl:"output,optional"`
}

// SSHConfig represents SSH configuration
type SSHConfig struct {
	DefaultUser       string `hcl:"default_user,optional"`
	DefaultPort       int    `hcl:"default_port,optional" validate:"omitempty,min=1,max=65535"`
	ConnectionTimeout int    `hcl:"connection_timeout,optional" validate:"omitempty,min=1,max=300"`
	CommandTimeout    int    `hcl:"command_timeout,optional" validate:"omitempty,min=1,max=3600"`
	RetryAttempts     int    `hcl:"retry_attempts,optional" validate:"omitempty,min=0,max=10"`
}

// InventoryConfig represents an inventory configuration (machines only)
type InventoryConfig struct {
	Machines []Machine `hcl:"machine,block" validate:"required,min=1,dive"`
}

// ActionsConfig represents an actions configuration (actions only)
type ActionsConfig struct {
	Actions []Action `hcl:"action,block" validate:"dive"`
}

// Machine represents a remote machine configuration
type Machine struct {
	Name     string            `hcl:"name,label" validate:"required"`
	Host     string            `hcl:"host" validate:"required"`
	Port     int               `hcl:"port,optional" validate:"omitempty,min=1,max=65535"`
	User     string            `hcl:"user" validate:"required"`
	Password string            `hcl:"password,optional"`
	KeyFile  string            `hcl:"key_file,optional"`
	Tags     map[string]string `hcl:"tags,optional" validate:"omitempty,dive,keys,required,endkeys,required"`
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

// Custom validation tags for mutual exclusivity and authentication requirements
const (
	// Custom validation tags
	TagMachineAuth   = "machine_auth"   // Either password or key_file must be provided
	TagActionExec    = "action_exec"    // Either command or script must be provided, but not both
	TagUniqueMachine = "unique_machine" // Machine names must be unique
	TagUniqueAction  = "unique_action"  // Action names must be unique
	TagValidPort     = "valid_port"     // Port must be valid (1-65535)
	TagValidTimeout  = "valid_timeout"  // Timeout must be reasonable (1-3600 seconds)
	TagValidTags     = "valid_tags"     // Tags must be non-empty strings
	TagValidMachines = "valid_machines" // Machine references must exist
)
