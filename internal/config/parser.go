package config

import (
	"errors"
	"fmt"
	"path/filepath"

	"spooky/internal/logging"

	"github.com/hashicorp/hcl/v2"
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

// resolveProjectPaths resolves relative paths in project configuration
func resolveProjectPaths(configFile string, project *ProjectConfig) {
	if project.InventoryFile != "" {
		project.InventoryFile = resolvePath(configFile, project.InventoryFile, false)
	}
	if project.ActionsFile != "" {
		project.ActionsFile = resolvePath(configFile, project.ActionsFile, false)
	}
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
	logger := logging.GetLogger()

	logger.Info("Parsing project configuration",
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

	// Resolve relative paths
	resolveProjectPaths(filename, config)

	logger.Info("Project configuration parsed successfully",
		logging.String("config_file", filename),
		logging.String("project_name", config.Name),
	)

	return config, nil
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
	return parseInventoryWithWrapper(filename)
}

// ParseActionsConfig parses an actions configuration file
func ParseActionsConfig(filename string) (*ActionsConfig, error) {
	return parseActionsWithWrapper(filename)
}

// parseInventoryWithWrapper parses an inventory configuration file with wrapper block
// nolint:dupl // Acceptable duplication - different types and purposes
func parseInventoryWithWrapper(filename string) (*InventoryConfig, error) {
	return parseConfigWithWrapper(filename, "inventory", &InventoryWrapper{},
		func(wrapper *InventoryWrapper) (*InventoryConfig, error) {
			if wrapper.Inventory == nil {
				return nil, errors.New("no inventory block found in configuration")
			}
			return wrapper.Inventory, nil
		},
		func(config *InventoryConfig) {
			for i := range config.Machines {
				resolveMachinePaths(filename, &config.Machines[i])
			}
		})
}

// parseActionsWithWrapper parses an actions configuration file with wrapper block
// nolint:dupl // Acceptable duplication - different types and purposes
func parseActionsWithWrapper(filename string) (*ActionsConfig, error) {
	return parseConfigWithWrapper(filename, "actions", &ActionsWrapper{},
		func(wrapper *ActionsWrapper) (*ActionsConfig, error) {
			if wrapper.Actions == nil {
				return nil, errors.New("no actions block found in configuration")
			}
			return wrapper.Actions, nil
		},
		func(config *ActionsConfig) {
			for i := range config.Actions {
				resolveActionPaths(filename, &config.Actions[i])
			}
		})
}

// parseConfigWithWrapper is a generic helper function to reduce code duplication
func parseConfigWithWrapper[T any, W any](
	filename, configType string,
	wrapper W,
	extractConfig func(W) (*T, error),
	resolvePaths func(*T),
) (*T, error) {
	logger := logging.GetLogger()

	logger.Info("Parsing "+configType+" configuration with wrapper block",
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
		return nil, errors.New("failed to parse " + configType + " HCL file: " + diagError)
	}

	// Validate wrapper blocks
	if err := validateWrapperBlocks(file); err != nil {
		logger.Error("Failed to validate wrapper blocks", err,
			logging.String("config_file", filename),
		)
		return nil, fmt.Errorf("wrapper block validation failed: %w", err)
	}

	// Decode the configuration using wrapper
	diags = gohcl.DecodeBody(file.Body, nil, wrapper)
	if diags.HasErrors() {
		diagError := diags.Error()
		logger.Error("Failed to decode "+configType+" configuration", errors.New(diagError),
			logging.String("config_file", filename),
		)
		return nil, errors.New("failed to decode " + configType + " configuration: " + diagError)
	}

	// Extract configuration from wrapper
	config, err := extractConfig(wrapper)
	if err != nil {
		return nil, err
	}

	logger.Debug(configType+" configuration decoded successfully",
		logging.String("config_file", filename),
	)

	// Resolve relative paths
	resolvePaths(config)

	logger.Info(configType+" configuration parsed successfully",
		logging.String("config_file", filename),
	)

	return config, nil
}

// validateWrapperBlocks ensures proper wrapper block usage
func validateWrapperBlocks(file *hcl.File) error {
	content, _, diags := file.Body.PartialContent(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "inventory"},
			{Type: "actions"},
		},
	})
	if diags.HasErrors() {
		return fmt.Errorf("failed to parse wrapper blocks: %s", diags.Error())
	}

	inventoryBlocks := 0
	actionsBlocks := 0

	for _, block := range content.Blocks {
		switch block.Type {
		case "inventory":
			inventoryBlocks++
		case "actions":
			actionsBlocks++
		}
	}

	if inventoryBlocks > 1 {
		return fmt.Errorf("multiple inventory blocks found (expected 1)")
	}
	if actionsBlocks > 1 {
		return fmt.Errorf("multiple actions blocks found (expected 1)")
	}

	return nil
}
