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
	return parseInventoryWithWrapper(filename)
}

// ParseActionsConfig parses an actions configuration file
func ParseActionsConfig(filename string) (*ActionsConfig, error) {
	return parseActionsWithWrapper(filename)
}

// parseInventoryWithWrapper parses inventory configuration with wrapper block
func parseInventoryWithWrapper(filename string) (*InventoryConfig, error) {
	logger := logging.GetLogger()

	logger.Info("Parsing inventory configuration with wrapper block",
		logging.String("config_file", filename),
	)

	parser := hclparse.NewParser()

	// Read the file
	file, diags := parser.ParseHCLFile(filename)
	if diags.HasErrors() {
		diagError := diags.Error()
		logger.Error("Failed to parse inventory HCL file", errors.New(diagError),
			logging.String("config_file", filename),
		)
		return nil, errors.New("failed to parse inventory HCL file: " + diagError)
	}

	// Validate wrapper blocks
	if err := validateWrapperBlocks(file); err != nil {
		logger.Error("Failed to validate wrapper blocks", err,
			logging.String("config_file", filename),
		)
		return nil, fmt.Errorf("wrapper block validation failed: %w", err)
	}

	// Decode the configuration using wrapper
	var wrapper InventoryWrapper
	diags = gohcl.DecodeBody(file.Body, nil, &wrapper)
	if diags.HasErrors() {
		diagError := diags.Error()
		logger.Error("Failed to decode inventory configuration", errors.New(diagError),
			logging.String("config_file", filename),
		)
		return nil, errors.New("failed to decode inventory configuration: " + diagError)
	}

	if wrapper.Inventory == nil {
		return nil, errors.New("no inventory block found in configuration")
	}

	config := wrapper.Inventory

	logger.Debug("Inventory configuration decoded successfully",
		logging.String("config_file", filename),
		logging.Int("machine_count", len(config.Machines)),
	)

	// Resolve relative paths in machine configurations
	for i := range config.Machines {
		resolveMachinePaths(filename, &config.Machines[i])
	}

	logger.Info("Inventory configuration parsed successfully",
		logging.String("config_file", filename),
		logging.Int("machine_count", len(config.Machines)),
	)

	return config, nil
}

// parseActionsWithWrapper parses actions configuration with wrapper block
func parseActionsWithWrapper(filename string) (*ActionsConfig, error) {
	logger := logging.GetLogger()

	logger.Info("Parsing actions configuration with wrapper block",
		logging.String("config_file", filename),
	)

	parser := hclparse.NewParser()

	// Read the file
	file, diags := parser.ParseHCLFile(filename)
	if diags.HasErrors() {
		diagError := diags.Error()
		logger.Error("Failed to parse actions HCL file", errors.New(diagError),
			logging.String("config_file", filename),
		)
		return nil, errors.New("failed to parse actions HCL file: " + diagError)
	}

	// Validate wrapper blocks
	if err := validateWrapperBlocks(file); err != nil {
		logger.Error("Failed to validate wrapper blocks", err,
			logging.String("config_file", filename),
		)
		return nil, fmt.Errorf("wrapper block validation failed: %w", err)
	}

	// Decode the configuration using wrapper
	var wrapper ActionsWrapper
	diags = gohcl.DecodeBody(file.Body, nil, &wrapper)
	if diags.HasErrors() {
		diagError := diags.Error()
		logger.Error("Failed to decode actions configuration", errors.New(diagError),
			logging.String("config_file", filename),
		)
		return nil, errors.New("failed to decode actions configuration: " + diagError)
	}

	if wrapper.Actions == nil {
		return nil, errors.New("no actions block found in configuration")
	}

	config := wrapper.Actions

	logger.Debug("Actions configuration decoded successfully",
		logging.String("config_file", filename),
		logging.Int("action_count", len(config.Actions)),
	)

	// Resolve relative paths in action configurations
	for i := range config.Actions {
		resolveActionPaths(filename, &config.Actions[i])
	}

	logger.Info("Actions configuration parsed successfully",
		logging.String("config_file", filename),
		logging.Int("action_count", len(config.Actions)),
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
