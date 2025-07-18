package config

// Config represents the main configuration structure
type Config struct {
	Servers []Server `hcl:"server,block" validate:"required,min=1,dive"`
	Actions []Action `hcl:"action,block" validate:"dive"`
}

// Server represents a remote server configuration
type Server struct {
	Name     string            `hcl:"name,label" validate:"required"`
	Host     string            `hcl:"host" validate:"required"`
	Port     int               `hcl:"port,optional" validate:"omitempty,min=1,max=65535"`
	User     string            `hcl:"user" validate:"required"`
	Password string            `hcl:"password,optional"`
	KeyFile  string            `hcl:"key_file,optional"`
	Tags     map[string]string `hcl:"tags,optional" validate:"omitempty,dive,keys,required,endkeys,required"`
}

// Action represents an action to be executed on servers
type Action struct {
	Name        string   `hcl:"name,label" validate:"required"`
	Description string   `hcl:"description,optional"`
	Command     string   `hcl:"command,optional"`
	Script      string   `hcl:"script,optional"`
	Servers     []string `hcl:"servers,optional" validate:"omitempty,dive,required"`
	Tags        []string `hcl:"tags,optional" validate:"omitempty,dive,required"`
	Timeout     int      `hcl:"timeout,optional" validate:"omitempty,min=1,max=3600"`
	Parallel    bool     `hcl:"parallel,optional"`
}

// Custom validation tags for mutual exclusivity and authentication requirements
const (
	// Custom validation tags
	TagServerAuth   = "server_auth"   // Either password or key_file must be provided
	TagActionExec   = "action_exec"   // Either command or script must be provided, but not both
	TagUniqueServer = "unique_server" // Server names must be unique
	TagUniqueAction = "unique_action" // Action names must be unique
	TagValidPort    = "valid_port"    // Port must be valid (1-65535)
	TagValidTimeout = "valid_timeout" // Timeout must be reasonable (1-3600 seconds)
	TagValidTags    = "valid_tags"    // Tags must be non-empty strings
	TagValidServers = "valid_servers" // Server references must exist
)
