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
func resolvePath(configFile, path string, debug bool) string {
	if filepath.IsAbs(path) {
		return path
	}

	configDir := filepath.Dir(configFile)
	resolved := filepath.Join(configDir, path)
	if debug {
		fmt.Printf("[DEBUG] resolvePath: configFile=%q, path=%q, configDir=%q, resolved=%q\n", configFile, path, configDir, resolved)
	}
	return resolved
}

// resolveMachinePaths resolves relative paths in machine configuration
func resolveMachinePaths(configFile string, machine *Machine) {
	if machine.KeyFile != "" {
		machine.KeyFile = resolvePath(configFile, machine.KeyFile, false)
	}
}

// resolveActionPaths resolves relative paths in action configuration
func resolveActionPaths(configFile string, action *Action) {
	if action.Script != "" {
		action.Script = resolvePath(configFile, action.Script, false)
	}
}

// parseConfigFile is a common function for parsing HCL configuration files
func parseConfigFile(filename, configType string, target interface{}, resolvePaths func(string, interface{})) error {
	logger := logging.GetLogger()

	logger.Info("Parsing "+configType+" configuration file",
		logging.String("config_file", filename),
	)

	parser := hclparse.NewParser()

	// Read the file
	file, diags := parser.ParseHCLFile(filename)
	if diags.HasErrors() {
		diagError := diags.Error()
		logger.Error("Failed to parse "+configType+" HCL file", errors.New(diagError),
			logging.String("config_file", filename),
		)
		return errors.New("failed to parse " + configType + " HCL file: " + diagError)
	}

	// Decode the configuration
	diags = gohcl.DecodeBody(file.Body, nil, target)
	if diags.HasErrors() {
		diagError := diags.Error()
		logger.Error("Failed to decode "+configType+" configuration", errors.New(diagError),
			logging.String("config_file", filename),
		)
		return errors.New("failed to decode " + configType + " configuration: " + diagError)
	}

	// Resolve paths if resolver function is provided
	if resolvePaths != nil {
		resolvePaths(filename, target)
	}

	return nil
}

// parseConfigWithResolver is a generic function for parsing configs with path resolution
func parseConfigWithResolver[T any](filename, configType string, target *T, resolver func(string, *T)) (*T, error) {
	logger := logging.GetLogger()

	err := parseConfigFile(filename, configType, target, func(filename string, target interface{}) {
		if cfg, ok := target.(*T); ok {
			resolver(filename, cfg)
		}
	})
	if err != nil {
		return nil, err
	}

	logger.Info(configType+" configuration parsed successfully",
		logging.String("config_file", filename))
	return target, nil
}

// ParseConfig parses an HCL2 configuration file (legacy combined format)
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
		logging.Int("machine_count", len(config.Machines)),
		logging.Int("action_count", len(config.Actions)),
	)

	// Resolve relative paths in machine configurations
	for i := range config.Machines {
		resolveMachinePaths(filename, &config.Machines[i])
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
		logging.Int("machine_count", len(config.Machines)),
		logging.Int("action_count", len(config.Actions)),
	)

	return &config, nil
}

// ParseProjectConfig parses a project configuration file
func ParseProjectConfig(filename string) (*ProjectConfig, error) {
	return ParseProjectConfigWithDebug(filename, false)
}

// ParseProjectConfigWithDebug parses a project configuration file with optional debug output
func ParseProjectConfigWithDebug(filename string, debug bool) (*ProjectConfig, error) {
	logger := logging.GetLogger()

	logger.Info("Parsing project configuration file",
		logging.String("config_file", filename),
	)

	parser := hclparse.NewParser()

	// Read the file
	file, diags := parser.ParseHCLFile(filename)
	if diags.HasErrors() {
		diagError := diags.Error()
		logger.Error("Failed to parse project HCL file", errors.New(diagError),
			logging.String("config_file", filename),
		)
		return nil, errors.New("failed to parse project HCL file: " + diagError)
	}

	// Decode the configuration using wrapper
	var wrapper ProjectConfigWrapper
	diags = gohcl.DecodeBody(file.Body, nil, &wrapper)
	if diags.HasErrors() {
		diagError := diags.Error()
		logger.Error("Failed to decode project configuration", errors.New(diagError),
			logging.String("config_file", filename),
		)
		return nil, errors.New("failed to decode project configuration: " + diagError)
	}

	if wrapper.Project == nil {
		return nil, errors.New("no project block found in configuration")
	}

	config := wrapper.Project

	logger.Debug("Project configuration decoded successfully",
		logging.String("config_file", filename),
		logging.String("project_name", config.Name),
	)

	// Resolve relative paths for file references
	if config.InventoryFile != "" {
		config.InventoryFile = resolvePath(filename, config.InventoryFile, debug)
	}
	if config.ActionsFile != "" {
		config.ActionsFile = resolvePath(filename, config.ActionsFile, debug)
	}

	logger.Info("Project configuration parsed successfully",
		logging.String("config_file", filename),
		logging.String("project_name", config.Name),
		logging.String("inventory_file", config.InventoryFile),
		logging.String("actions_file", config.ActionsFile),
	)

	return config, nil
}

// ParseInventoryConfig parses an inventory configuration file
func ParseInventoryConfig(filename string) (*InventoryConfig, error) {
	var config InventoryConfig
	return parseConfigWithResolver(filename, "inventory", &config, func(filename string, cfg *InventoryConfig) {
		for i := range cfg.Machines {
			resolveMachinePaths(filename, &cfg.Machines[i])
		}
	})
}

// ParseActionsConfig parses an actions configuration file
func ParseActionsConfig(filename string) (*ActionsConfig, error) {
	var config ActionsConfig
	return parseConfigWithResolver(filename, "actions", &config, func(filename string, cfg *ActionsConfig) {
		for i := range cfg.Actions {
			resolveActionPaths(filename, &cfg.Actions[i])
		}
	})
}
