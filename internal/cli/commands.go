package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"spooky/internal/config"
	"spooky/internal/facts"
	"spooky/internal/logging"
)

var (
	// Validate command flags
	schema string
	strict bool
	format string

	// List command flags
	filter  string
	sort    string
	reverse bool

	// Project command flags
	projectPath   string
	projectName   string
	listVerbose   bool
	validateDebug bool
)

var ProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage spooky projects",
	Long:  `Manage spooky projects with separated inventory and actions configuration`,
}

func init() {
	// Add subcommands to ProjectCmd
	ProjectCmd.AddCommand(ProjectInitCmd)
	ProjectCmd.AddCommand(ProjectValidateCmd)
	ProjectCmd.AddCommand(ProjectListCmd)

	// Add flags to ProjectListCmd
	ProjectListCmd.Flags().BoolVar(&listVerbose, "verbose", false, "Show detailed output including machine and action details")

	// Add flags to ProjectValidateCmd
	ProjectValidateCmd.Flags().BoolVar(&validateDebug, "debug", false, "Show debug output including path resolution details")
}

var ProjectInitCmd = &cobra.Command{
	Use:   "init <PROJECT_NAME> [PATH]",
	Short: "Initialize a new spooky project",
	Long:  `Create a new spooky project with separated inventory and actions configuration`,
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(_ *cobra.Command, args []string) error {
		logger := logging.GetLogger()
		projectName := args[0]
		path := "."
		if len(args) > 1 {
			path = args[1]
		}

		return initProject(logger, projectName, path)
	},
}

var ProjectValidateCmd = &cobra.Command{
	Use:   "validate [PROJECT_PATH]",
	Short: "Validate a spooky project",
	Long:  `Validate the configuration files in a spooky project`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		logger := logging.GetLogger()
		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		return validateProject(logger, path)
	},
}

var ProjectListCmd = &cobra.Command{
	Use:   "list [PROJECT_PATH]",
	Short: "List project resources",
	Long:  `List machines and actions from a spooky project`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		logger := logging.GetLogger()
		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		return listProject(logger, path)
	},
}

var ValidateCmd = &cobra.Command{
	Use:   "validate <source>",
	Short: "Validate configuration files and templates",
	Long:  `Validate the syntax and structure of configuration files and templates`,
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		logger := logging.GetLogger()
		source := args[0]

		// TODO: Support remote sources (Git, S3, HTTP)
		// For now, only support local files
		if !isLocalFile(source) {
			return fmt.Errorf("remote sources not yet supported: %s", source)
		}

		// Validate config file exists
		if _, err := os.Stat(source); os.IsNotExist(err) {
			logger.Error("Config file not found", err, logging.String("config_file", source))
			return fmt.Errorf("config file %s does not exist", source)
		}

		// Parse configuration
		config, err := config.ParseConfig(source)
		if err != nil {
			logger.Error("Configuration validation failed", err, logging.String("config_file", source))
			return fmt.Errorf("validation failed: %w", err)
		}

		logger.Info("Configuration file validated successfully",
			logging.String("config_file", source),
			logging.Int("machine_count", len(config.Machines)),
			logging.Int("action_count", len(config.Actions)),
		)

		return nil
	},
}

var ListCmd = &cobra.Command{
	Use:   "list <resource|config-file>",
	Short: "List resources or configurations",
	Long:  `Display resources or configurations based on the specified resource type or from a configuration file. Use 'spooky facts list' for facts.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		logger := logging.GetLogger()

		// If no arguments provided, show help
		if len(args) == 0 {
			return fmt.Errorf("resource type or config file is required. Use 'spooky list <resource|config-file>' or 'spooky facts list' for facts")
		}

		resource := args[0]

		// Check if it's a config file path
		if strings.HasSuffix(resource, ".hcl") || strings.HasSuffix(resource, ".yaml") || strings.HasSuffix(resource, ".yml") {
			return listFromConfigFile(logger, resource)
		}

		// Otherwise treat as resource type
		switch resource {
		case "machines":
			return listMachines(logger)
		case "templates":
			return listTemplates(logger)
		case "configs":
			return listConfigs(logger)
		case "actions":
			return listActions(logger)
		default:
			return fmt.Errorf("unknown resource type: %s. Supported types: machines, templates, configs, actions. Use 'spooky facts list' for facts", resource)
		}
	},
}

// InitCommands initializes all CLI commands and their flags
func InitCommands() {
	// Check if flags are already initialized to prevent redefinition
	if ValidateCmd.Flags().Lookup("schema") != nil {
		return // Already initialized
	}

	// Validate command flags
	ValidateCmd.Flags().StringVar(&schema, "schema", "", "Path to schema file for validation")
	ValidateCmd.Flags().BoolVar(&strict, "strict", false, "Enable strict validation mode")
	ValidateCmd.Flags().StringVar(&format, "format", "text", "Output format: text, json, yaml")

	// List command flags
	ListCmd.Flags().StringVar(&format, "format", "table", "Output format: table, json, yaml")
	ListCmd.Flags().StringVar(&filter, "filter", "", "Filter results by expression")
	ListCmd.Flags().StringVar(&sort, "sort", "name", "Sort field")
	ListCmd.Flags().BoolVar(&reverse, "reverse", false, "Reverse sort order")

	// Project command flags
	ProjectCmd.Flags().StringVar(&projectPath, "path", ".", "Path to the project directory")
	ProjectCmd.Flags().StringVar(&projectName, "name", "", "Name of the project (required for init)")

	// Initialize facts commands
	initFactsCommands()
}

// isLocalFile checks if the source is a local file
func isLocalFile(source string) bool {
	// Check if it's a remote source (Git, S3, HTTP)
	if strings.HasPrefix(source, "github.com/") ||
		strings.HasPrefix(source, "git://") ||
		strings.HasPrefix(source, "s3://") ||
		strings.HasPrefix(source, "http://") ||
		strings.HasPrefix(source, "https://") {
		return false
	}
	return true
}

// listFromConfigFile lists resources from a configuration file
func listFromConfigFile(logger logging.Logger, configFile string) error {
	// Validate config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		logger.Error("Config file not found", err, logging.String("config_file", configFile))
		return fmt.Errorf("config file %s does not exist", configFile)
	}

	// Parse configuration
	config, err := config.ParseConfig(configFile)
	if err != nil {
		logger.Error("Failed to parse configuration", err, logging.String("config_file", configFile))
		return fmt.Errorf("failed to parse config: %w", err)
	}

	logger.Info("Listing resources from configuration file",
		logging.String("config_file", configFile),
		logging.Int("machine_count", len(config.Machines)),
		logging.Int("action_count", len(config.Actions)))

	// Display machines
	if len(config.Machines) > 0 {
		fmt.Printf("Machines (%d):\n", len(config.Machines))
		for _, machine := range config.Machines {
			fmt.Printf("  - %s (%s@%s:%d)\n", machine.Name, machine.User, machine.Host, machine.Port)
		}
		fmt.Println()
	}

	// Display actions
	if len(config.Actions) > 0 {
		fmt.Printf("Actions (%d):\n", len(config.Actions))
		for i := range config.Actions {
			action := &config.Actions[i]
			desc := action.Description
			if desc == "" {
				desc = "No description"
			}
			fmt.Printf("  - %s: %s\n", action.Name, desc)
		}
		fmt.Println()
	}

	// Display summary
	fmt.Printf("Configuration Summary:\n")
	fmt.Printf("  File: %s\n", configFile)
	fmt.Printf("  Machines: %d\n", len(config.Machines))
	fmt.Printf("  Actions: %d\n", len(config.Actions))

	return nil
}

// List functions for different resource types
func listMachines(logger logging.Logger) error {
	// TODO: Implement machine listing from configuration or inventory
	logger.Info("Listing machines (not yet implemented)")
	return fmt.Errorf("machine listing not yet implemented")
}

func listFacts(logger logging.Logger) error {
	// Create fact manager
	manager := facts.NewManager(nil)

	// Try to get cached facts first
	allFacts, err := manager.GetAllFacts()
	if err != nil {
		logger.Error("Failed to retrieve cached facts", err)
		return fmt.Errorf("failed to retrieve cached facts: %w", err)
	}

	// If no cached facts, collect from local server
	if len(allFacts) == 0 {
		logger.Info("No cached facts found, collecting from local server")
		collection, err := manager.CollectAllFacts("local")
		if err != nil {
			logger.Error("Failed to collect facts from local server", err)
			return fmt.Errorf("failed to collect facts from local server: %w", err)
		}

		// Convert collection to slice of facts
		for _, fact := range collection.Facts {
			allFacts = append(allFacts, fact)
		}
	}

	if len(allFacts) == 0 {
		fmt.Println("No facts found. Use 'spooky facts collect <server>' to gather facts first.")
		return nil
	}

	// Check output format
	if format == "json" {
		// Create a structured response for JSON output
		response := map[string]interface{}{
			"summary": map[string]interface{}{
				"total_facts": len(allFacts),
				"servers":     len(getUniqueServers(allFacts)),
			},
			"facts": allFacts,
		}

		jsonData, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal facts to JSON: %w", err)
		}
		fmt.Println(string(jsonData))
	} else {
		// Group facts by server
		serverFacts := make(map[string][]*facts.Fact)
		for _, fact := range allFacts {
			serverFacts[fact.Server] = append(serverFacts[fact.Server], fact)
		}

		// Display facts in table format
		fmt.Printf("Facts Summary:\n")
		fmt.Printf("Total facts: %d\n", len(allFacts))
		fmt.Printf("Servers: %d\n\n", len(serverFacts))

		for server, facts := range serverFacts {
			fmt.Printf("Server: %s (%d facts)\n", server, len(facts))
			fmt.Printf("%-30s %-15s %-20s %s\n", "KEY", "SOURCE", "TTL", "VALUE")
			fmt.Printf("%s\n", strings.Repeat("-", 80))

			for _, fact := range facts {
				// Truncate value if too long
				valueStr := fmt.Sprintf("%v", fact.Value)
				if len(valueStr) > 40 {
					valueStr = valueStr[:37] + "..."
				}

				// Format TTL
				ttlStr := "expired"
				if fact.TTL > 0 {
					ttlStr = fact.TTL.String()
				}

				fmt.Printf("%-30s %-15s %-20s %s\n",
					fact.Key,
					fact.Source,
					ttlStr,
					valueStr)
			}
			fmt.Println()
		}
	}

	logger.Info("Facts listed successfully", logging.Int("total_facts", len(allFacts)), logging.Int("servers", len(getUniqueServers(allFacts))))
	return nil
}

// Helper function to get unique servers from facts
func getUniqueServers(facts []*facts.Fact) []string {
	servers := make(map[string]bool)
	for _, fact := range facts {
		servers[fact.Server] = true
	}

	result := make([]string, 0, len(servers))
	for server := range servers {
		result = append(result, server)
	}
	return result
}

func listTemplates(logger logging.Logger) error {
	// TODO: Implement template listing
	logger.Info("Listing templates (not yet implemented)")
	return fmt.Errorf("template listing not yet implemented")
}

func listConfigs(logger logging.Logger) error {
	// TODO: Implement config listing
	logger.Info("Listing configs (not yet implemented)")
	return fmt.Errorf("config listing not yet implemented")
}

func listActions(logger logging.Logger) error {
	// TODO: Implement actions listing from configuration
	logger.Info("Listing actions (not yet implemented)")
	return fmt.Errorf("actions listing not yet implemented")
}

// Project functions

// initProject initializes a new spooky project
func initProject(logger logging.Logger, projectName, path string) error {
	logger.Info("Initializing new spooky project",
		logging.String("project_name", projectName),
		logging.String("path", path))

	// Create project directory
	projectDir := filepath.Join(path, projectName)
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		logger.Error("Failed to create project directory", err,
			logging.String("project_dir", projectDir))
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create subdirectories
	dirs := []string{"templates", "files", "logs"}
	for _, dir := range dirs {
		dirPath := filepath.Join(projectDir, dir)
		if err := os.MkdirAll(dirPath, 0o755); err != nil {
			logger.Error("Failed to create subdirectory", err,
				logging.String("dir", dirPath))
			return fmt.Errorf("failed to create subdirectory %s: %w", dir, err)
		}
	}

	// Create project.hcl
	projectConfig := fmt.Sprintf(`project "%s" {
  description = "%s project"
  version = "1.0.0"
  environment = "development"
  
  # File references
  inventory_file = "inventory.hcl"
  actions_file = "actions.hcl"
  
  # Project settings
  default_timeout = 300
  default_parallel = true
  
  # Storage configuration
  storage {
    type = "badgerdb"
    path = ".facts.db"
  }
  
  # Logging configuration
  logging {
    level = "info"
    format = "json"
    output = "logs/spooky.log"
  }
  
  # SSH configuration
  ssh {
    default_user = "debian"
    default_port = 22
    connection_timeout = 30
    command_timeout = 300
    retry_attempts = 3
  }
  
  # Tags for project-wide targeting
  tags = {
    project = "%s"
  }
}`, projectName, projectName, projectName)

	projectFile := filepath.Join(projectDir, "project.hcl")
	if err := os.WriteFile(projectFile, []byte(projectConfig), 0o600); err != nil {
		logger.Error("Failed to create project.hcl", err,
			logging.String("file", projectFile))
		return fmt.Errorf("failed to create project.hcl: %w", err)
	}

	// Create inventory.hcl
	inventoryConfig := fmt.Sprintf(`# Inventory for %s project
# Add your machine definitions here

inventory {
  machine "example-server" {
    host     = "192.168.1.100"
    port     = 22
    user     = "debian"
    password = "your-password"
    tags = {
      environment = "development"
      role = "web"
    }
  }
}`, projectName)

	inventoryFile := filepath.Join(projectDir, "inventory.hcl")
	if err := os.WriteFile(inventoryFile, []byte(inventoryConfig), 0o600); err != nil {
		logger.Error("Failed to create inventory.hcl", err,
			logging.String("file", inventoryFile))
		return fmt.Errorf("failed to create inventory.hcl: %w", err)
	}

	// Create actions.hcl
	actionsConfig := fmt.Sprintf(`# Actions for %s project
# Add your action definitions here

actions {
  action "check-status" {
    description = "Check server status"
    command     = "uptime && df -h"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
  }

  action "update-system" {
    description = "Update system packages"
    command     = "apt update && apt upgrade -y"
    tags        = ["environment=development"]
    parallel    = true
    timeout     = 600
  }
}`, projectName)

	actionsFile := filepath.Join(projectDir, "actions.hcl")
	if err := os.WriteFile(actionsFile, []byte(actionsConfig), 0o600); err != nil {
		logger.Error("Failed to create actions.hcl", err,
			logging.String("file", actionsFile))
		return fmt.Errorf("failed to create actions.hcl: %w", err)
	}

	// Create .gitignore
	gitignore := `# Spooky project gitignore

# Facts database
.facts.db/
*.db

# Logs
logs/
*.log

# Temporary files
*.tmp
*.temp

# SSH keys and sensitive files
*.pem
*.key
id_rsa*
*.pub

# Environment files
.env
.env.local

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# IDE files
.vscode/
.idea/
*.swp
*.swo
*~

# Backup files
*.bak
*.backup`

	gitignoreFile := filepath.Join(projectDir, ".gitignore")
	if err := os.WriteFile(gitignoreFile, []byte(gitignore), 0o600); err != nil {
		logger.Error("Failed to create .gitignore", err,
			logging.String("file", gitignoreFile))
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	// Create README.md
	readme := fmt.Sprintf(`# %s Project

This is a spooky project with separated inventory and actions configuration.

## Project Structure

%s/
├── project.hcl          # Project configuration and settings
├── inventory.hcl        # Machine definitions
├── actions.hcl          # Action definitions for automation
├── .gitignore          # Git ignore rules
├── templates/           # Template files for dynamic content
├── files/              # Static files to be deployed
└── README.md           # This file

## Usage

### Execute an action:
spooky execute actions.hcl --inventory inventory.hcl --action check-status

### Execute with tag targeting:
spooky execute actions.hcl --inventory inventory.hcl --action update-system --tags "role=web"

## Configuration

- **project.hcl**: Project settings, storage, logging, and SSH configuration
- **inventory.hcl**: Machine definitions with tags for targeting
- **actions.hcl**: Automation actions with tag-based targeting

## Benefits

1. **Reusability**: Actions can be applied to different inventories
2. **Maintainability**: Clear separation of concerns
3. **Flexibility**: Mix and match actions with different machine groups
4. **Version Control**: Better tracking of changes to machines vs actions
`, projectName, projectName)

	readmeFile := filepath.Join(projectDir, "README.md")
	if err := os.WriteFile(readmeFile, []byte(readme), 0o600); err != nil {
		logger.Error("Failed to create README.md", err,
			logging.String("file", readmeFile))
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	logger.Info("Project initialized successfully",
		logging.String("project_name", projectName),
		logging.String("project_dir", projectDir))

	fmt.Printf("✅ Project '%s' initialized successfully at %s\n", projectName, projectDir)
	fmt.Printf("📁 Project structure created with separated inventory and actions\n")
	fmt.Printf("📝 Edit inventory.hcl to add your machines\n")
	fmt.Printf("⚡ Edit actions.hcl to add your automation actions\n")
	fmt.Printf("🔧 Edit project.hcl to configure project settings\n")

	return nil
}

// validateConfigFile validates a configuration file using the provided parser function
func validateConfigFile(logger logging.Logger, filePath, fileType string, parser func(string) error) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logger.Warn(fileType+" file not found",
			logging.String("file", filePath))
		return nil
	}

	if err := parser(filePath); err != nil {
		logger.Error("Failed to parse "+fileType+" configuration", err,
			logging.String("file", filePath))
		return fmt.Errorf("failed to parse "+fileType+" configuration: %w", err)
	}

	logger.Info(fileType+" configuration validated",
		logging.String("file", filePath))
	return nil
}

// validateProject validates a spooky project
func validateProject(logger logging.Logger, path string) error {
	logger.Info("Validating spooky project",
		logging.String("path", path))

	// Check if project.hcl exists
	projectFile := filepath.Join(path, "project.hcl")
	if _, err := os.Stat(projectFile); os.IsNotExist(err) {
		logger.Error("Project file not found", err,
			logging.String("file", projectFile))
		return fmt.Errorf("project.hcl not found in %s", path)
	}

	// Parse project configuration with debug flag
	projectConfig, err := config.ParseProjectConfigWithDebug(projectFile, validateDebug)
	if err != nil {
		logger.Error("Failed to parse project configuration", err,
			logging.String("file", projectFile))
		return fmt.Errorf("failed to parse project configuration: %w", err)
	}

	logger.Info("Project configuration validated",
		logging.String("project_name", projectConfig.Name),
		logging.String("inventory_file", projectConfig.InventoryFile),
		logging.String("actions_file", projectConfig.ActionsFile))

	// Validate inventory file if it exists
	if projectConfig.InventoryFile != "" {
		if err := validateConfigFile(logger, projectConfig.InventoryFile, "Inventory", func(file string) error {
			_, err := config.ParseInventoryConfig(file)
			return err
		}); err != nil {
			return err
		}
	}

	// Validate actions file if it exists
	if projectConfig.ActionsFile != "" {
		if err := validateConfigFile(logger, projectConfig.ActionsFile, "Actions", func(file string) error {
			_, err := config.ParseActionsConfig(file)
			return err
		}); err != nil {
			return err
		}
	}

	fmt.Printf("✅ Project validation successful\n")
	fmt.Printf("📋 Project: %s\n", projectConfig.Name)
	fmt.Printf("📁 Path: %s\n", path)
	if projectConfig.Description != "" {
		fmt.Printf("📝 Description: %s\n", projectConfig.Description)
	}

	return nil
}

// listProject lists resources from a spooky project
func listProject(logger logging.Logger, path string) error {
	logger.Info("Listing spooky project resources",
		logging.String("path", path))

	// Check if project.hcl exists
	projectFile := filepath.Join(path, "project.hcl")
	if _, err := os.Stat(projectFile); os.IsNotExist(err) {
		logger.Error("Project file not found", err,
			logging.String("file", projectFile))
		return fmt.Errorf("project.hcl not found in %s", path)
	}

	// Parse project configuration
	projectConfig, err := config.ParseProjectConfig(projectFile)
	if err != nil {
		logger.Error("Failed to parse project configuration", err,
			logging.String("file", projectFile))
		return fmt.Errorf("failed to parse project configuration: %w", err)
	}

	fmt.Printf("Project: %s\n", projectConfig.Name)
	if projectConfig.Description != "" {
		fmt.Printf("Description: %s\n", projectConfig.Description)
	}
	fmt.Printf("Path: %s\n\n", path)

	machineCount := 0
	actionCount := 0

	// Count machines from inventory
	if projectConfig.InventoryFile != "" {
		if _, err := os.Stat(projectConfig.InventoryFile); os.IsNotExist(err) {
			fmt.Printf("⚠️  Inventory file not found: %s\n", projectConfig.InventoryFile)
		} else {
			inventoryConfig, err := config.ParseInventoryConfig(projectConfig.InventoryFile)
			if err != nil {
				logger.Error("Failed to parse inventory configuration", err,
					logging.String("file", projectConfig.InventoryFile))
				return fmt.Errorf("failed to parse inventory configuration: %w", err)
			}

			machineCount = len(inventoryConfig.Machines)

			if listVerbose {
				fmt.Printf("Machines (%d):\n", machineCount)
				for _, machine := range inventoryConfig.Machines {
					fmt.Printf("  - %s (%s@%s:%d)\n", machine.Name, machine.User, machine.Host, machine.Port)
				}
				fmt.Println()
			}
		}
	}

	// Count actions
	if projectConfig.ActionsFile != "" {
		if _, err := os.Stat(projectConfig.ActionsFile); os.IsNotExist(err) {
			fmt.Printf("⚠️  Actions file not found: %s\n", projectConfig.ActionsFile)
		} else {
			actionsConfig, err := config.ParseActionsConfig(projectConfig.ActionsFile)
			if err != nil {
				logger.Error("Failed to parse actions configuration", err,
					logging.String("file", projectConfig.ActionsFile))
				return fmt.Errorf("failed to parse actions configuration: %w", err)
			}

			actionCount = len(actionsConfig.Actions)

			if listVerbose {
				fmt.Printf("Actions (%d):\n", actionCount)
				for i := range actionsConfig.Actions {
					action := &actionsConfig.Actions[i]
					desc := action.Description
					if desc == "" {
						desc = "No description"
					}
					fmt.Printf("  - %s: %s\n", action.Name, desc)
				}
				fmt.Println()
			}
		}
	}

	// Show summary
	fmt.Printf("Summary: %d machines, %d actions\n", machineCount, actionCount)
	if !listVerbose {
		fmt.Printf("Use --verbose for detailed output\n")
	}

	return nil
}
