package spooky

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

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

// ParseConfig parses an HCL2 configuration file
func ParseConfig(filename string) (*Config, error) {
	parser := hclparse.NewParser()

	// Read the file
	file, diags := parser.ParseHCLFile(filename)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL file: %s", diags.Error())
	}

	// Decode the configuration
	var config Config
	diags = gohcl.DecodeBody(file.Body, nil, &config)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to decode configuration: %s", diags.Error())
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &config, nil
}

// validateConfig validates the configuration structure
func validateConfig(config *Config) error {
	if len(config.Servers) == 0 {
		return fmt.Errorf("at least one server must be defined")
	}

	// Actions are optional - you can have a config with just servers for listing purposes
	// if len(config.Actions) == 0 {
	// 	return fmt.Errorf("at least one action must be defined")
	// }

	// Validate servers
	serverNames := make(map[string]bool)
	for _, server := range config.Servers {
		if server.Name == "" {
			return fmt.Errorf("server name cannot be empty")
		}
		if server.Host == "" {
			return fmt.Errorf("server host cannot be empty for server %s", server.Name)
		}
		if server.User == "" {
			return fmt.Errorf("server user cannot be empty for server %s", server.Name)
		}
		if server.Port == 0 {
			server.Port = 22 // Default SSH port
		}
		if server.Password == "" && server.KeyFile == "" {
			return fmt.Errorf("either password or key_file must be specified for server %s", server.Name)
		}

		if serverNames[server.Name] {
			return fmt.Errorf("duplicate server name: %s", server.Name)
		}
		serverNames[server.Name] = true
	}

	// Validate actions
	actionNames := make(map[string]bool)
	for _, action := range config.Actions {
		if action.Name == "" {
			return fmt.Errorf("action name cannot be empty")
		}
		if action.Command == "" && action.Script == "" {
			return fmt.Errorf("either command or script must be specified for action %s", action.Name)
		}
		if action.Command != "" && action.Script != "" {
			return fmt.Errorf("cannot specify both command and script for action %s", action.Name)
		}

		if actionNames[action.Name] {
			return fmt.Errorf("duplicate action name: %s", action.Name)
		}
		actionNames[action.Name] = true
	}

	return nil
}

// GetServersForAction returns the list of servers that should execute an action
func GetServersForAction(action *Action, config *Config) ([]*Server, error) {
	var targetServers []*Server

	// If specific servers are specified, use those
	if len(action.Servers) > 0 {
		serverMap := make(map[string]*Server)
		for i := range config.Servers {
			serverMap[config.Servers[i].Name] = &config.Servers[i]
		}

		for _, serverName := range action.Servers {
			if server, exists := serverMap[serverName]; exists {
				targetServers = append(targetServers, server)
			} else {
				return nil, fmt.Errorf("server '%s' not found for action '%s'", serverName, action.Name)
			}
		}
		return targetServers, nil
	}

	// If tags are specified, find servers with matching tags
	if len(action.Tags) > 0 {
		for i := range config.Servers {
			server := &config.Servers[i]
			for _, tag := range action.Tags {
				if value, exists := server.Tags[tag]; exists && value != "" {
					targetServers = append(targetServers, server)
					break
				}
			}
		}
		return targetServers, nil
	}

	// If no servers or tags specified, use all servers
	for i := range config.Servers {
		targetServers = append(targetServers, &config.Servers[i])
	}

	return targetServers, nil
}
