package config

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// resolvePath resolves a path relative to the config file's directory
func resolvePath(configFile, path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	configDir := filepath.Dir(configFile)
	return filepath.Join(configDir, path)
}

// resolveServerPaths resolves relative paths in server configuration
func resolveServerPaths(configFile string, server *Server) {
	if server.KeyFile != "" {
		server.KeyFile = resolvePath(configFile, server.KeyFile)
	}
}

// resolveActionPaths resolves relative paths in action configuration
func resolveActionPaths(configFile string, action *Action) {
	if action.Script != "" {
		action.Script = resolvePath(configFile, action.Script)
	}
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

	// Resolve relative paths in server configurations
	for i := range config.Servers {
		resolveServerPaths(filename, &config.Servers[i])
	}

	// Resolve relative paths in action configurations
	for i := range config.Actions {
		resolveActionPaths(filename, &config.Actions[i])
	}

	// Set default values
	SetDefaults(&config)

	// Validate configuration
	if err := ValidateConfig(&config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &config, nil
}
