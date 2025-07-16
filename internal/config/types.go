package config

import "fmt"

// Config represents the main configuration structure
type Config struct {
	Servers []Server `hcl:"server,block" validate:"required,min=1,dive"`
	Actions []Action `hcl:"action,block" validate:"dive"`
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

// Validate ensures either password or key_file is provided for servers
func (s *Server) Validate() error {
	if s.Password == "" && s.KeyFile == "" {
		return fmt.Errorf("either password or key_file must be specified for server %s", s.Name)
	}
	return nil
}

// Validate ensures either command or script is provided for actions
func (a *Action) Validate() error {
	if a.Command == "" && a.Script == "" {
		return fmt.Errorf("either command or script must be specified for action %s", a.Name)
	}
	if a.Command != "" && a.Script != "" {
		return fmt.Errorf("cannot specify both command and script for action %s", a.Name)
	}
	return nil
}
