package config

import (
	"fmt"
)

// ValidateConfig validates the configuration structure
func ValidateConfig(config *Config) error {
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
			server.Port = DefaultSSHPort
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
