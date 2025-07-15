package config

// Config represents the main configuration structure
type Config struct {
	Servers []Server `hcl:"server,block"`
	Actions []Action `hcl:"action,block"`
}

// Server represents a remote server configuration
type Server struct {
	Name     string            `hcl:"name,label"`
	Host     string            `hcl:"host"`
	Port     int               `hcl:"port,optional"`
	User     string            `hcl:"user"`
	Password string            `hcl:"password,optional"`
	KeyFile  string            `hcl:"key_file,optional"`
	Tags     map[string]string `hcl:"tags,optional"`
}

// Action represents an action to be executed on servers
type Action struct {
	Name        string   `hcl:"name,label"`
	Description string   `hcl:"description,optional"`
	Command     string   `hcl:"command,optional"`
	Script      string   `hcl:"script,optional"`
	Servers     []string `hcl:"servers,optional"`
	Tags        []string `hcl:"tags,optional"`
	Timeout     int      `hcl:"timeout,optional"`
	Parallel    bool     `hcl:"parallel,optional"`
}
