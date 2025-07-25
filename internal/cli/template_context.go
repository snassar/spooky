package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"spooky/internal/config"
	"spooky/internal/facts"
	"spooky/internal/logging"
	"spooky/internal/ssh"
)

// TemplateContext holds all data available to templates
type TemplateContext struct {
	// Project configuration
	Project *config.ProjectConfig

	// Machine facts (from facts.db or JSON)
	Facts map[string]interface{}

	// Inventory information
	Machines []*config.Machine

	// Actions configuration
	Actions []*config.Action

	// Server-specific data (when --server is specified)
	ServerFacts map[string]interface{}

	// Environment variables
	Environment map[string]string

	// Custom data files
	CustomData map[string]interface{}
}

// NewTemplateContext creates a new template context for a project
func NewTemplateContext(logger logging.Logger, projectPath string) (*TemplateContext, error) {
	ctx := &TemplateContext{
		Facts:       make(map[string]interface{}),
		Environment: make(map[string]string),
		CustomData:  make(map[string]interface{}),
	}

	logger.Info("Creating template context",
		logging.String("project_path", projectPath))

	// Load project configuration
	if err := ctx.loadProjectConfig(logger, projectPath); err != nil {
		return nil, fmt.Errorf("failed to load project config: %w", err)
	}

	// Load facts
	if err := ctx.loadFacts(logger, projectPath); err != nil {
		logger.Warn("Failed to load facts", logging.String("error", err.Error()))
	}

	// Load inventory
	if err := ctx.loadInventory(logger, projectPath); err != nil {
		logger.Warn("Failed to load inventory", logging.String("error", err.Error()))
	} else {
		logger.Info("Inventory loaded successfully", logging.Int("machines", len(ctx.Machines)))
	}

	// Load actions
	if err := ctx.loadActions(logger, projectPath); err != nil {
		logger.Warn("Failed to load actions", logging.String("error", err.Error()))
	} else {
		logger.Info("Actions loaded successfully", logging.Int("actions", len(ctx.Actions)))
	}

	// Load environment variables
	ctx.loadEnvironment()

	// Load custom data files
	if err := ctx.loadCustomData(logger, projectPath); err != nil {
		logger.Warn("Failed to load custom data", logging.String("error", err.Error()))
	}

	logger.Info("Template context created successfully",
		logging.Int("machines", len(ctx.Machines)),
		logging.Int("actions", len(ctx.Actions)),
		logging.Int("facts", len(ctx.Facts)))

	// Direct output for debugging
	fmt.Printf("DEBUG: Template context created with %d machines, %d actions, %d facts\n",
		len(ctx.Machines), len(ctx.Actions), len(ctx.Facts))
	fmt.Printf("DEBUG: Project config - Name: %s, InventoryFile: %s, ActionsFile: %s\n",
		ctx.Project.Name, ctx.Project.InventoryFile, ctx.Project.ActionsFile)

	// Test if files exist
	inventoryPath := ctx.Project.InventoryFile
	actionsPath := ctx.Project.ActionsFile

	fmt.Printf("DEBUG: Resolved paths - Inventory: %s (exists: %t), Actions: %s (exists: %t)\n",
		inventoryPath, fileExists(inventoryPath), actionsPath, fileExists(actionsPath))

	return ctx, nil
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// loadProjectConfig loads the project configuration
func (ctx *TemplateContext) loadProjectConfig(logger logging.Logger, projectPath string) error {
	projectFile := filepath.Join(projectPath, "project.hcl")
	logger.Info("Loading project config",
		logging.String("project_file", projectFile))

	projectConfig, err := config.ParseProjectConfig(projectFile)
	if err != nil {
		logger.Error("Failed to parse project config", err,
			logging.String("project_file", projectFile))
		return fmt.Errorf("failed to parse project config: %w", err)
	}

	ctx.Project = projectConfig
	logger.Info("Project config loaded successfully",
		logging.String("project_name", projectConfig.Name),
		logging.String("inventory_file", projectConfig.InventoryFile),
		logging.String("actions_file", projectConfig.ActionsFile))

	return nil
}

// loadFacts loads facts from facts.db or JSON files
func (ctx *TemplateContext) loadFacts(logger logging.Logger, projectPath string) error {
	// Try to load from facts.db first
	factsDBPath := filepath.Join(projectPath, ".facts.db")
	if _, err := os.Stat(factsDBPath); err == nil {
		_ = facts.NewManager(nil) // TODO: Implement facts loading from badgerdb
		logger.Info("Found facts.db, but loading not yet implemented")
		return nil
	}

	// Try to load from JSON files
	factsDir := filepath.Join(projectPath, "facts")
	if _, err := os.Stat(factsDir); err == nil {
		return ctx.loadFactsFromJSON(logger, factsDir)
	}

	// Try to load from custom-facts directory
	customFactsDir := filepath.Join(projectPath, "custom-facts")
	if _, err := os.Stat(customFactsDir); err == nil {
		return ctx.loadFactsFromJSON(logger, customFactsDir)
	}

	return nil
}

// loadFactsFromJSON loads facts from JSON files in a directory
func (ctx *TemplateContext) loadFactsFromJSON(logger logging.Logger, factsDir string) error {
	entries, err := os.ReadDir(factsDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(factsDir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			logger.Warn("Failed to read facts file",
				logging.String("file", filePath),
				logging.String("error", err.Error()))
			continue
		}

		var facts map[string]interface{}
		if err := json.Unmarshal(data, &facts); err != nil {
			logger.Warn("Failed to parse facts file",
				logging.String("file", filePath),
				logging.String("error", err.Error()))
			continue
		}

		// Use filename (without extension) as the key
		key := filepath.Base(entry.Name())
		key = key[:len(key)-len(filepath.Ext(key))]
		ctx.Facts[key] = facts
	}

	return nil
}

// loadInventory loads inventory configuration
func (ctx *TemplateContext) loadInventory(logger logging.Logger, projectPath string) error {
	logger.Info("Loading inventory",
		logging.String("project_path", projectPath),
		logging.String("inventory_file", ctx.Project.InventoryFile))

	// Use the resolved path directly since it's already absolute
	inventoryPath := ctx.Project.InventoryFile

	logger.Info("Resolved inventory path", logging.String("inventory_path", inventoryPath))

	return loadConfigFileWithProcessor(logger, projectPath, inventoryPath, "inventory",
		func(inventoryConfig *config.InventoryConfig) {
			ctx.Machines = make([]*config.Machine, len(inventoryConfig.Machines))
			for i := range inventoryConfig.Machines {
				ctx.Machines[i] = &inventoryConfig.Machines[i]
			}
			logger.Info("Inventory config processed", logging.Int("machines", len(ctx.Machines)))
		})
}

// loadActions loads actions configuration
func (ctx *TemplateContext) loadActions(logger logging.Logger, projectPath string) error {
	logger.Info("Loading actions",
		logging.String("project_path", projectPath),
		logging.String("actions_file", ctx.Project.ActionsFile))

	// Use the resolved path directly since it's already absolute
	actionsPath := ctx.Project.ActionsFile

	logger.Info("Resolved actions path", logging.String("actions_path", actionsPath))

	return loadConfigFileWithProcessor(logger, projectPath, actionsPath, "actions",
		func(actionsConfig *config.ActionsConfig) {
			ctx.Actions = make([]*config.Action, len(actionsConfig.Actions))
			for i := range actionsConfig.Actions {
				ctx.Actions[i] = &actionsConfig.Actions[i]
			}
			logger.Info("Actions config processed", logging.Int("actions", len(ctx.Actions)))
		})
}

// loadConfigFileWithProcessor is a generic helper to load configuration files
func loadConfigFileWithProcessor[T any](logger logging.Logger, projectPath, fileName, configType string, processor func(*T)) error {
	if fileName == "" {
		logger.Info("No file name provided for config type", logging.String("config_type", configType))
		return nil
	}

	// Use the fileName directly since it's already resolved
	configPath := fileName

	logger.Info("Loading config file",
		logging.String("config_type", configType),
		logging.String("config_path", configPath))

	// Use type-specific parser based on configType
	var parsedConfig interface{}
	var err error

	switch configType {
	case "inventory":
		parsedConfig, err = config.ParseInventoryConfig(configPath)
	case "actions":
		parsedConfig, err = config.ParseActionsConfig(configPath)
	default:
		return fmt.Errorf("unknown config type: %s", configType)
	}

	if err != nil {
		logger.Error("Failed to parse config file", err,
			logging.String("config_type", configType),
			logging.String("config_path", configPath))
		return fmt.Errorf("failed to parse %s: %w", configType, err)
	}

	// Type assertion and processing
	if typedConfig, ok := parsedConfig.(*T); ok {
		processor(typedConfig)
		logger.Info("Config file processed successfully",
			logging.String("config_type", configType),
			logging.String("config_path", configPath))
	} else {
		logger.Error("Type assertion failed", nil,
			logging.String("config_type", configType),
			logging.String("config_path", configPath))
		return fmt.Errorf("type assertion failed for %s config", configType)
	}

	return nil
}

// loadEnvironment loads environment variables
func (ctx *TemplateContext) loadEnvironment() {
	for _, env := range os.Environ() {
		key, value, found := cut(env, "=")
		if found {
			ctx.Environment[key] = value
		}
	}
}

// loadCustomData loads custom data files
func (ctx *TemplateContext) loadCustomData(logger logging.Logger, projectPath string) error {
	dataDir := filepath.Join(projectPath, "data")
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		return nil
	}

	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(dataDir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			logger.Warn("Failed to read data file",
				logging.String("file", filePath),
				logging.String("error", err.Error()))
			continue
		}

		var customData interface{}
		if err := json.Unmarshal(data, &customData); err != nil {
			logger.Warn("Failed to parse data file",
				logging.String("file", filePath),
				logging.String("error", err.Error()))
			continue
		}

		// Use filename (without extension) as the key
		key := filepath.Base(entry.Name())
		key = key[:len(key)-len(filepath.Ext(key))]
		ctx.CustomData[key] = customData
	}

	return nil
}

// LoadServerFacts loads facts for a specific server
func (ctx *TemplateContext) LoadServerFacts(logger logging.Logger, serverName string) error {
	// Find the machine in inventory
	var targetMachine *config.Machine
	for _, machine := range ctx.Machines {
		if machine.Name == serverName {
			targetMachine = machine
			break
		}
	}

	if targetMachine == nil {
		return fmt.Errorf("server '%s' not found in inventory", serverName)
	}

	// Create SSH client with 30 second timeout
	sshClient, err := ssh.NewSSHClient(targetMachine, 30)
	if err != nil {
		return fmt.Errorf("failed to create SSH client for %s: %w", serverName, err)
	}
	defer sshClient.Close()

	// Initialize server facts map
	ctx.ServerFacts = make(map[string]interface{})

	// 1. Add all machine fields and tags to ServerFacts
	ctx.ServerFacts["name"] = targetMachine.Name
	ctx.ServerFacts["host"] = targetMachine.Host
	ctx.ServerFacts["port"] = targetMachine.Port
	ctx.ServerFacts["user"] = targetMachine.User
	ctx.ServerFacts["tags"] = targetMachine.Tags
	ctx.ServerFacts["keyfile"] = targetMachine.KeyFile
	ctx.ServerFacts["password"] = targetMachine.Password

	// 2. Collect system facts (extensible map)
	systemFacts := map[string]string{
		"machine_id":   "cat /etc/machine-id 2>/dev/null || echo 'unknown'",
		"os_version":   "uname -r",
		"hostname":     "hostname",
		"ip_address":   "hostname -I | awk '{print $1}'",
		"disk_space":   "df -h / | tail -1 | awk '{print $4}'",
		"memory_info":  "free -h | grep '^Mem:' | awk '{print $2}'",
		"cpu_info":     "nproc",
		"uptime":       "uptime -p",
		"kernel":       "uname -s",
		"architecture": "uname -m",
	}

	for factName, command := range systemFacts {
		result, err := sshClient.ExecuteCommand(command)
		if err != nil {
			logger.Warn("Failed to collect fact",
				logging.String("fact", factName),
				logging.String("server", serverName),
				logging.String("error", err.Error()))
			ctx.ServerFacts[factName] = "unknown"
			continue
		}
		ctx.ServerFacts[factName] = strings.TrimSpace(result)
	}

	// 3. Collect file system facts
	fileFacts := map[string]string{
		"etc_exists":   "test -f /etc/passwd && echo 'true' || echo 'false'",
		"var_exists":   "test -d /var && echo 'true' || echo 'false'",
		"tmp_writable": "test -w /tmp && echo 'true' || echo 'false'",
		"home_exists":  "test -d /home && echo 'true' || echo 'false'",
	}

	for factName, command := range fileFacts {
		result, err := sshClient.ExecuteCommand(command)
		if err != nil {
			logger.Warn("Failed to collect file fact",
				logging.String("fact", factName),
				logging.String("server", serverName),
				logging.String("error", err.Error()))
			ctx.ServerFacts[factName] = "false"
			continue
		}
		ctx.ServerFacts[factName] = strings.TrimSpace(result) == "true"
	}

	// 4. Collect service facts
	serviceFacts := map[string]string{
		"nginx_running":    "systemctl is-active nginx 2>/dev/null || echo 'inactive'",
		"apache_running":   "systemctl is-active apache2 2>/dev/null || echo 'inactive'",
		"mysql_running":    "systemctl is-active mysql 2>/dev/null || echo 'inactive'",
		"postgres_running": "systemctl is-active postgresql 2>/dev/null || echo 'inactive'",
		"docker_running":   "systemctl is-active docker 2>/dev/null || echo 'inactive'",
	}

	for factName, command := range serviceFacts {
		result, err := sshClient.ExecuteCommand(command)
		if err != nil {
			logger.Warn("Failed to collect service fact",
				logging.String("fact", factName),
				logging.String("server", serverName),
				logging.String("error", err.Error()))
			ctx.ServerFacts[factName] = "unknown"
			continue
		}
		ctx.ServerFacts[factName] = strings.TrimSpace(result)
	}

	// 5. (Future) Custom fact commands from inventory or project config
	// Example: ctx.ServerFacts["custom_fact"] = ...

	logger.Info("Server facts collected successfully",
		logging.String("server", serverName),
		logging.Int("fact_count", len(ctx.ServerFacts)))

	return nil
}

// GetTemplateFunctions returns template functions for the context
func (ctx *TemplateContext) GetTemplateFunctions() map[string]interface{} {
	return map[string]interface{}{
		// Project functions
		"project":            func() *config.ProjectConfig { return ctx.Project },
		"projectName":        func() string { return ctx.Project.Name },
		"projectDescription": func() string { return ctx.Project.Description },

		// Facts functions
		"facts": func() map[string]interface{} { return ctx.Facts },
		"fact":  func(key string) interface{} { return ctx.Facts[key] },

		// Machine functions
		"machines": func() []*config.Machine { return ctx.Machines },
		"machine": func(name string) *config.Machine {
			for _, m := range ctx.Machines {
				if m.Name == name {
					return m
				}
			}
			return nil
		},

		// Action functions
		"actions": func() []*config.Action { return ctx.Actions },
		"action": func(name string) *config.Action {
			for _, a := range ctx.Actions {
				if a.Name == name {
					return a
				}
			}
			return nil
		},

		// Environment functions
		"env": func(key string) string { return ctx.Environment[key] },
		"envOrDefault": func(key, defaultValue string) string {
			if value, exists := ctx.Environment[key]; exists {
				return value
			}
			return defaultValue
		},

		// Custom data functions
		"data":    func() map[string]interface{} { return ctx.CustomData },
		"dataKey": func(key string) interface{} { return ctx.CustomData[key] },

		// Server facts functions
		"serverFacts": func() map[string]interface{} { return ctx.ServerFacts },
		"serverFact":  func(key string) interface{} { return ctx.ServerFacts[key] },
	}
}

// cut is a helper function similar to strings.Cut (Go 1.18+)
func cut(s, sep string) (before, after string, found bool) {
	if i := index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}

// index is a helper function similar to strings.Index
func index(s, substr string) int {
	n := len(substr)
	switch {
	case n == 0:
		return 0
	case n == 1:
		return indexByte(s, substr[0])
	case n == len(s):
		if substr == s {
			return 0
		}
		return -1
	case n > len(s):
		return -1
	}

	for i := 0; i <= len(s)-n; i++ {
		if s[i:i+n] == substr {
			return i
		}
	}
	return -1
}

// indexByte is a helper function similar to strings.IndexByte
func indexByte(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}
