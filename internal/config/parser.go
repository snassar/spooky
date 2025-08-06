package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	spookylogging "spooky/internal/logging"

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

// ParseProjectConfig parses a project configuration file
func ParseProjectConfig(filename string) (*ProjectConfig, error) {
	logger := spookylogging.GetLogger()

	logger.Info("Parsing project configuration",
		spookylogging.String("config_file", filename),
	)

	parser := hclparse.NewParser()

	// Read the file
	file, diags := parser.ParseHCLFile(filename)
	if diags.HasErrors() {
		diagError := diags.Error()
		logger.Error("Failed to parse project HCL file", errors.New(diagError),
			spookylogging.String("config_file", filename),
		)
		return nil, errors.New("failed to parse project HCL file: " + diagError)
	}

	// Decode the configuration using wrapper
	var wrapper ProjectConfigWrapper
	diags = gohcl.DecodeBody(file.Body, nil, &wrapper)
	if diags.HasErrors() {
		diagError := diags.Error()
		logger.Error("Failed to decode project configuration", errors.New(diagError),
			spookylogging.String("config_file", filename),
		)
		return nil, errors.New("failed to decode project configuration: " + diagError)
	}

	if wrapper.Project == nil {
		return nil, errors.New("no project block found in configuration")
	}

	config := wrapper.Project

	// Resolve relative paths
	resolveProjectPaths(filename, config)

	// Validate the project config
	validator := NewValidator()
	if err := validator.validate.Struct(config); err != nil {
		return nil, fmt.Errorf("project validation failed: %w", err)
	}

	logger.Info("Project configuration parsed successfully",
		spookylogging.String("config_file", filename),
		spookylogging.String("project_name", config.Name),
	)

	return config, nil
}

// ParseProjectConfigWithDebug parses a project configuration file with optional debug output
func ParseProjectConfigWithDebug(filename string, debug bool) (*ProjectConfig, error) {
	logger := spookylogging.GetLogger()

	logger.Info("Parsing project configuration file",
		spookylogging.String("config_file", filename),
	)

	parser := hclparse.NewParser()

	// Read the file
	file, diags := parser.ParseHCLFile(filename)
	if diags.HasErrors() {
		diagError := diags.Error()
		logger.Error("Failed to parse project HCL file", errors.New(diagError),
			spookylogging.String("config_file", filename),
		)
		return nil, errors.New("failed to parse project HCL file: " + diagError)
	}

	// Decode the configuration using wrapper
	var wrapper ProjectConfigWrapper
	diags = gohcl.DecodeBody(file.Body, nil, &wrapper)
	if diags.HasErrors() {
		diagError := diags.Error()
		logger.Error("Failed to decode project configuration", errors.New(diagError),
			spookylogging.String("config_file", filename),
		)
		return nil, errors.New("failed to decode project configuration: " + diagError)
	}

	if wrapper.Project == nil {
		return nil, errors.New("no project block found in configuration")
	}

	config := wrapper.Project

	logger.Debug("Project configuration decoded successfully",
		spookylogging.String("config_file", filename),
		spookylogging.String("project_name", config.Name),
	)

	// Resolve relative paths for file references
	if config.InventoryFile != "" {
		config.InventoryFile = resolvePath(filename, config.InventoryFile, debug)
	}
	if config.ActionsFile != "" {
		config.ActionsFile = resolvePath(filename, config.ActionsFile, debug)
	}

	logger.Info("Project configuration parsed successfully",
		spookylogging.String("config_file", filename),
		spookylogging.String("project_name", config.Name),
		spookylogging.String("inventory_file", config.InventoryFile),
		spookylogging.String("actions_file", config.ActionsFile),
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

// LoadActionsConfig loads actions from multiple sources and merges them
// 1. Load actions.hcl from project root (if exists)
// 2. Load all .hcl files from actions/ directory (if exists)
// 3. Merge all actions into a single ActionsConfig
func LoadActionsConfig(projectPath string) (*ActionsConfig, error) {
	logger := spookylogging.GetLogger()

	// Initialize merged config
	mergedConfig := &ActionsConfig{
		Actions: []Action{},
	}

	// 1. Try to load actions.hcl from project root
	rootActionsFile := filepath.Join(projectPath, "actions.hcl")
	if _, err := os.Stat(rootActionsFile); err == nil {
		logger.Info("Loading actions from root file", spookylogging.String("file", rootActionsFile))
		rootConfig, err := ParseActionsConfig(rootActionsFile)
		if err != nil {
			logger.Error("Failed to parse root actions file", err, spookylogging.String("file", rootActionsFile))
			return nil, fmt.Errorf("failed to parse root actions file: %w", err)
		}
		mergedConfig.Actions = append(mergedConfig.Actions, rootConfig.Actions...)
		logger.Info("Loaded actions from root file", spookylogging.Int("actions", len(rootConfig.Actions)))
	}

	// 2. Try to load all .hcl files from actions/ directory
	actionsDir := filepath.Join(projectPath, "actions")
	if _, err := os.Stat(actionsDir); err == nil {
		logger.Info("Loading actions from directory", spookylogging.String("dir", actionsDir))

		entries, err := os.ReadDir(actionsDir)
		if err != nil {
			logger.Error("Failed to read actions directory", err, spookylogging.String("dir", actionsDir))
			return nil, fmt.Errorf("failed to read actions directory: %w", err)
		}

		// Sort entries to ensure consistent loading order
		var actionFiles []string
		for _, entry := range entries {
			if !entry.IsDir() && filepath.Ext(entry.Name()) == ".hcl" {
				actionFiles = append(actionFiles, entry.Name())
			}
		}

		// Sort files to ensure consistent loading order (e.g., 01-dependencies.hcl comes before 02-database.hcl)
		sort.Strings(actionFiles)

		for _, fileName := range actionFiles {
			filePath := filepath.Join(actionsDir, fileName)
			logger.Info("Loading action file", spookylogging.String("file", filePath))

			fileConfig, err := ParseActionsConfig(filePath)
			if err != nil {
				logger.Error("Failed to parse action file", err, spookylogging.String("file", filePath))
				return nil, fmt.Errorf("failed to parse action file %s: %w", fileName, err)
			}

			mergedConfig.Actions = append(mergedConfig.Actions, fileConfig.Actions...)
			logger.Info("Loaded actions from file",
				spookylogging.String("file", fileName),
				spookylogging.Int("actions", len(fileConfig.Actions)))
		}
	}

	// Check if we loaded any actions
	if len(mergedConfig.Actions) == 0 {
		logger.Warn("No actions found in project", spookylogging.String("project_path", projectPath))
		return mergedConfig, nil
	}

	logger.Info("Successfully loaded all actions",
		spookylogging.String("project_path", projectPath),
		spookylogging.Int("total_actions", len(mergedConfig.Actions)))

	return mergedConfig, nil
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
	logger := spookylogging.GetLogger()

	logger.Info("Parsing "+configType+" configuration with wrapper block",
		spookylogging.String("config_file", filename),
	)

	parser := hclparse.NewParser()

	// Read the file
	file, diags := parser.ParseHCLFile(filename)
	if diags.HasErrors() {
		diagError := diags.Error()
		logger.Error("Failed to parse "+configType+" HCL file", errors.New(diagError),
			spookylogging.String("config_file", filename),
		)
		return nil, errors.New("failed to parse " + configType + " HCL file: " + diagError)
	}

	// Validate wrapper blocks
	if err := validateWrapperBlocks(file); err != nil {
		logger.Error("Failed to validate wrapper blocks", err,
			spookylogging.String("config_file", filename),
		)
		return nil, fmt.Errorf("wrapper block validation failed: %w", err)
	}

	// Decode the configuration using wrapper
	diags = gohcl.DecodeBody(file.Body, nil, wrapper)
	if diags.HasErrors() {
		diagError := diags.Error()
		logger.Error("Failed to decode "+configType+" configuration", errors.New(diagError),
			spookylogging.String("config_file", filename),
		)
		return nil, errors.New("failed to decode " + configType + " configuration: " + diagError)
	}

	// Extract configuration from wrapper
	config, err := extractConfig(wrapper)
	if err != nil {
		return nil, err
	}

	logger.Debug(configType+" configuration decoded successfully",
		spookylogging.String("config_file", filename),
	)

	// Resolve relative paths
	resolvePaths(config)

	logger.Info(configType+" configuration parsed successfully",
		spookylogging.String("config_file", filename),
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
