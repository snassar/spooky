package config

// Config represents the main configuration structure
type Config struct {
	Machines []Machine `hcl:"machine,block" validate:"required,min=1,dive"`
	Actions  []Action  `hcl:"action,block" validate:"dive"`
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
	Name        string   `hcl:"name,label" validate:"required"`
	Description string   `hcl:"description,optional"`
	Command     string   `hcl:"command,optional"`
	Script      string   `hcl:"script,optional"`
	Machines    []string `hcl:"machines,optional" validate:"omitempty,dive,required"`
	Tags        []string `hcl:"tags,optional" validate:"omitempty,dive,required"`
	Timeout     int      `hcl:"timeout,optional" validate:"omitempty,min=1,max=3600"`
	Parallel    bool     `hcl:"parallel,optional"`
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
