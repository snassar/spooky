package config

import (
	"errors"
	"fmt"
	"path/filepath"

	"spooky/internal/logging"

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
	logger := logging.GetLogger()

	logger.Info("Parsing configuration file",
		logging.String("config_file", filename),
	)

	parser := hclparse.NewParser()

	// Read the file
	file, diags := parser.ParseHCLFile(filename)
	if diags.HasErrors() {
		diagError := diags.Error()
		logger.Error("Failed to parse HCL file", errors.New(diagError),
			logging.String("config_file", filename),
		)
		return nil, errors.New("failed to parse HCL file: " + diagError)
	}

	// Decode the configuration
	var config Config
	diags = gohcl.DecodeBody(file.Body, nil, &config)
	if diags.HasErrors() {
		diagError := diags.Error()
		logger.Error("Failed to decode configuration", errors.New(diagError),
			logging.String("config_file", filename),
		)
		return nil, errors.New("failed to decode configuration: " + diagError)
	}

	logger.Debug("Configuration decoded successfully",
		logging.String("config_file", filename),
		logging.Int("server_count", len(config.Servers)),
		logging.Int("action_count", len(config.Actions)),
	)

	// Resolve relative paths in server configurations
	for i := range config.Servers {
		resolveServerPaths(filename, &config.Servers[i])
	}

	// Resolve relative paths in action configurations
	for i := range config.Actions {
		resolveActionPaths(filename, &config.Actions[i])
	}

	logger.Debug("Relative paths resolved",
		logging.String("config_file", filename),
	)

	// Set default values
	SetDefaults(&config)

	// Validate configuration
	if err := ValidateConfig(&config); err != nil {
		logger.Error("Configuration validation failed", err,
			logging.String("config_file", filename),
		)
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	logger.Info("Configuration parsed and validated successfully",
		logging.String("config_file", filename),
		logging.Int("server_count", len(config.Servers)),
		logging.Int("action_count", len(config.Actions)),
	)

	return &config, nil
}
