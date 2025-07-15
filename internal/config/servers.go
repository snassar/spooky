package config

import (
	"fmt"
)

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
