package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	"spooky/internal/config"
	"spooky/internal/facts"
	"spooky/internal/logging"
)

var (
	// Project command flags
	listVerbose   bool
	validateDebug bool
)

func init() {
	// Add flags to ListCmd
	ListCmd.Flags().BoolVar(&listVerbose, "verbose", false, "Show detailed output including machine and action details")

	// Add flags to ValidateCmd
	ValidateCmd.Flags().BoolVar(&validateDebug, "debug", false, "Show debug output including path resolution details")

	// Add flags to RenderTemplateCmd
	RenderTemplateCmd.Flags().String("output", "", "Output file path (default: stdout)")
	RenderTemplateCmd.Flags().Bool("dry-run", false, "Show what would be rendered without writing output")
	RenderTemplateCmd.Flags().String("server", "", "Server name for fact integration (e.g., web-001)")
	RenderTemplateCmd.Flags().String("ssh-key-path", "~/.ssh/", "Path to SSH private key or directory")
	RenderTemplateCmd.Flags().String("templates-dir", "", "Override templates directory (default: templates/)")
	RenderTemplateCmd.Flags().String("data-dir", "", "Override data directory (default: data/)")

	// Add flags to ValidateTemplateCmd
	ValidateTemplateCmd.Flags().String("templates-dir", "", "Override templates directory (default: templates/)")
	ValidateTemplateCmd.Flags().String("data-dir", "", "Override data directory (default: data/)")

	// Add flags to GatherFactsCmd
	GatherFactsCmd.Flags().Bool("dry-run", false, "Show what facts would be collected without connecting to machines")
	GatherFactsCmd.Flags().String("ssh-key-path", "~/.ssh/", "Path to SSH private key or directory")
	GatherFactsCmd.Flags().String("facts-db-path", "", "Override facts database path (default: .facts.db)")

	// Add flags to ListFactsCmd
	ListFactsCmd.Flags().String("facts-db-path", "", "Override facts database path (default: .facts.db)")

	// Add flags to ListTemplatesCmd
	ListTemplatesCmd.Flags().String("templates-dir", "", "Override templates directory (default: templates/)")
}

var InitCmd = &cobra.Command{
	Use:   "init <PROJECT_NAME> [PATH]",
	Short: "Initialize a new spooky project",
	Long:  `Create a new spooky project with separated inventory and actions configuration`,
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := logging.GetLogger()

		// Get project name from positional args or flag
		var projectName string
		if len(args) > 0 {
			projectName = args[0]
		}
		if flagName, _ := cmd.Flags().GetString("project"); flagName != "" {
			projectName = flagName
		}

		// Validate that we have a project name
		if projectName == "" {
			return fmt.Errorf("project name is required (use positional argument or --project flag)")
		}

		// Get path from positional args only
		path := "."
		if len(args) > 1 {
			path = args[1]
		}

		return initProject(logger, projectName, path)
	},
}

var ValidateCmd = &cobra.Command{
	Use:   "validate [PROJECT_PATH]",
	Short: "Validate a spooky project",
	Long:  `Validate the project files (project.hcl, inventory.hcl, actions.hcl) in a spooky project`,
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

var ListCmd = &cobra.Command{
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

var ListMachinesCmd = &cobra.Command{
	Use:   "list-machines [PROJECT_PATH]",
	Short: "List machines in a spooky project",
	Long:  `List all machines from a spooky project's inventory configuration`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		logger := logging.GetLogger()
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		return listProjectMachines(logger, path)
	},
}

var ListActionsCmd = &cobra.Command{
	Use:   "list-actions [PROJECT_PATH]",
	Short: "List actions in a spooky project",
	Long:  `List all actions from a spooky project's actions configuration`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		logger := logging.GetLogger()
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		return listProjectActions(logger, path)
	},
}

var ListTemplatesCmd = &cobra.Command{
	Use:   "list-templates [PROJECT_PATH]",
	Short: "List templates in a spooky project",
	Long:  `List all template files available in a spooky project`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := logging.GetLogger()
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		templatesDir, _ := cmd.Flags().GetString("templates-dir")
		return listProjectTemplates(logger, path, templatesDir)
	},
}

var ListFactsCmd = &cobra.Command{
	Use:   "list-facts [PROJECT_PATH]",
	Short: "List facts in a spooky project",
	Long:  `List all facts available in a spooky project`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := logging.GetLogger()
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		factsDBPath, _ := cmd.Flags().GetString("facts-db-path")
		return listProjectFacts(logger, path, factsDBPath)
	},
}

var GatherFactsCmd = &cobra.Command{
	Use:   "gather-facts [PROJECT_PATH]",
	Short: "Gather facts for machines in a spooky project",
	Long:  `Gather facts from all machines defined in a spooky project's inventory`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := logging.GetLogger()
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		sshKeyPath, _ := cmd.Flags().GetString("ssh-key-path")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		factsDBPath, _ := cmd.Flags().GetString("facts-db-path")
		return gatherProjectFacts(logger, path, sshKeyPath, dryRun, factsDBPath)
	},
}

var RenderTemplateCmd = &cobra.Command{
	Use:   "render-template <TEMPLATE_FILE> [PROJECT_PATH]",
	Short: "Render a template file",
	Long:  `Render a template file in the context of a spooky project`,
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := logging.GetLogger()
		templateFile := args[0]
		path := "."
		if len(args) > 1 {
			path = args[1]
		}

		output, _ := cmd.Flags().GetString("output")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		server, _ := cmd.Flags().GetString("server")
		sshKeyPath, _ := cmd.Flags().GetString("ssh-key-path")
		templatesDir, _ := cmd.Flags().GetString("templates-dir")
		dataDir, _ := cmd.Flags().GetString("data-dir")

		return renderProjectTemplate(logger, templateFile, path, output, dryRun, server, sshKeyPath, templatesDir, dataDir)
	},
}

var ValidateTemplateCmd = &cobra.Command{
	Use:   "validate-template <TEMPLATE_FILE> [PROJECT_PATH]",
	Short: "Validate a template file",
	Long:  `Validate a template file in the context of a spooky project`,
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := logging.GetLogger()
		templateFile := args[0]
		path := "."
		if len(args) > 1 {
			path = args[1]
		}

		templatesDir, _ := cmd.Flags().GetString("templates-dir")
		dataDir, _ := cmd.Flags().GetString("data-dir")

		return validateProjectTemplate(logger, templateFile, path, templatesDir, dataDir)
	},
}

var commandsInitialized bool

// InitCommands initializes all CLI commands and their flags
func InitCommands() {
	if commandsInitialized {
		return // Already initialized
	}

	// Project command flags
	InitCmd.Flags().String("project", "", "Name of the project (alternative to positional argument)")

	// Initialize facts commands
	initFactsCommands()

	commandsInitialized = true
}

// Project functions

// initProject initializes a new spooky project
func initProject(logger logging.Logger, projectName, path string) error {
	// Validate project name
	if strings.TrimSpace(projectName) == "" {
		return fmt.Errorf("project name cannot be empty")
	}

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
	dirs := []string{"templates", "files", "logs", "actions", "data"}
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

	// Create actions.hcl (main actions file)
	actionsConfig := fmt.Sprintf(`# Main actions for %s project
# This file can contain core actions or reference actions in the actions/ directory

actions {
  action "check-status" {
    description = "Check server status"
    command     = "uptime && df -h"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
  }
}`, projectName)

	actionsFile := filepath.Join(projectDir, "actions.hcl")
	if err := os.WriteFile(actionsFile, []byte(actionsConfig), 0o600); err != nil {
		logger.Error("Failed to create actions.hcl", err,
			logging.String("file", actionsFile))
		return fmt.Errorf("failed to create actions.hcl: %w", err)
	}

	// Create example action files in actions/ directory
	exampleActions := []struct {
		filename string
		content  string
	}{
		{
			"01-dependencies.hcl",
			fmt.Sprintf(`# Dependencies setup for %s project

actions {
  action "install-dependencies" {
    description = "Install system dependencies"
    command     = "apt update && apt install -y curl wget git"
    tags        = ["role=web", "setup"]
    parallel    = true
    timeout     = 600
  }
}`, projectName),
		},
		{
			"02-system-update.hcl",
			fmt.Sprintf(`# System updates for %s project

actions {
  action "update-system" {
    description = "Update system packages"
    command     = "apt update && apt upgrade -y"
    tags        = ["environment=development", "maintenance"]
    parallel    = true
    timeout     = 600
  }
}`, projectName),
		},
		{
			"03-monitoring.hcl",
			fmt.Sprintf(`# Monitoring actions for %s project

actions {
  action "check-disk-space" {
    description = "Check disk space usage"
    command     = "df -h"
    tags        = ["monitoring"]
    parallel    = true
    timeout     = 60
  }

  action "check-memory" {
    description = "Check memory usage"
    command     = "free -h"
    tags        = ["monitoring"]
    parallel    = true
    timeout     = 60
  }
}`, projectName),
		},
	}

	for _, actionFile := range exampleActions {
		filePath := filepath.Join(projectDir, "actions", actionFile.filename)
		if err := os.WriteFile(filePath, []byte(actionFile.content), 0o600); err != nil {
			logger.Error("Failed to create action file", err,
				logging.String("file", filePath))
			return fmt.Errorf("failed to create action file %s: %w", actionFile.filename, err)
		}
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
├── actions.hcl          # Main action definitions (optional)
├── actions/             # Directory for organized action files
│   ├── 01-dependencies.hcl
│   ├── 02-system-update.hcl
│   └── 03-monitoring.hcl
├── templates/           # Template files for dynamic content
├── data/               # Data files for templates
├── files/              # Static files to be deployed
├── logs/               # Log files
├── .gitignore          # Git ignore rules
└── README.md           # This file

## Actions Organization

Spooky supports flexible action organization:

1. **actions.hcl** - Main actions file (optional)
2. **actions/** directory - Organized action files
   - Files are loaded in alphabetical order
   - Use numbered prefixes (01-, 02-, etc.) for ordering
   - Each file can contain multiple actions

## Usage

### List all actions:
spooky list-actions

### Execute an action:
spooky execute actions.hcl --inventory inventory.hcl --action check-status

### Execute with tag targeting:
spooky execute actions.hcl --inventory inventory.hcl --action update-system --tags "role=web"

## Configuration

- **project.hcl**: Project settings, storage, logging, and SSH configuration
- **inventory.hcl**: Machine definitions with tags for targeting
- **actions.hcl**: Main actions (optional)
- **actions/**: Organized action files for better maintainability

## Benefits

1. **Reusability**: Actions can be applied to different inventories
2. **Maintainability**: Clear separation of concerns
3. **Flexibility**: Mix and match actions with different machine groups
4. **Organization**: Group related actions in separate files
5. **Version Control**: Better tracking of changes to machines vs actions
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

// validateConfigFile validates a project file (inventory.hcl or actions.hcl) using the provided parser function
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

	// Validate actions from multiple sources
	logger.Info("Validating actions configuration")
	if _, err := config.LoadActionsConfig(path); err != nil {
		logger.Error("Failed to validate actions configuration", err)
		return fmt.Errorf("failed to validate actions configuration: %w", err)
	}
	logger.Info("Actions configuration validated successfully")

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

// listProjectMachines lists machines from a spooky project inventory
func listProjectMachines(logger logging.Logger, path string) error {
	logger.Info("Listing machines in spooky project",
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

	// List machines from inventory
	if projectConfig.InventoryFile == "" {
		fmt.Println("No inventory file configured")
		return nil
	}

	if _, err := os.Stat(projectConfig.InventoryFile); os.IsNotExist(err) {
		fmt.Printf("⚠️  Inventory file not found: %s\n", projectConfig.InventoryFile)
		return nil
	}

	inventoryConfig, err := config.ParseInventoryConfig(projectConfig.InventoryFile)
	if err != nil {
		logger.Error("Failed to parse inventory configuration", err,
			logging.String("file", projectConfig.InventoryFile))
		return fmt.Errorf("failed to parse inventory configuration: %w", err)
	}

	if len(inventoryConfig.Machines) == 0 {
		fmt.Println("No machines found in inventory")
		return nil
	}

	fmt.Printf("Machines (%d):\n", len(inventoryConfig.Machines))
	for _, machine := range inventoryConfig.Machines {
		fmt.Printf("  - %s (%s@%s:%d)\n", machine.Name, machine.User, machine.Host, machine.Port)
	}

	return nil
}

// listProjectActions lists actions in a spooky project
func listProjectActions(logger logging.Logger, path string) error {
	logger.Info("Listing project actions",
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

	// Load actions from multiple sources
	actionsConfig, err := config.LoadActionsConfig(path)
	if err != nil {
		logger.Error("Failed to load actions configuration", err)
		return fmt.Errorf("failed to load actions configuration: %w", err)
	}

	if len(actionsConfig.Actions) == 0 {
		fmt.Println("No actions found")
		return nil
	}

	fmt.Printf("Actions (%d):\n", len(actionsConfig.Actions))
	for i := range actionsConfig.Actions {
		action := &actionsConfig.Actions[i]
		desc := action.Description
		if desc == "" {
			desc = "No description"
		}
		fmt.Printf("  - %s: %s\n", action.Name, desc)
	}

	return nil
}

// listProjectTemplates lists templates in a spooky project
func listProjectTemplates(logger logging.Logger, path, templatesDir string) error {
	logger.Info("Listing project templates",
		logging.String("path", path))

	// Get templates directory
	var targetTemplatesDir string
	if templatesDir != "" {
		targetTemplatesDir = templatesDir
	} else {
		targetTemplatesDir = filepath.Join(path, "templates")
	}

	if _, err := os.Stat(targetTemplatesDir); os.IsNotExist(err) {
		fmt.Printf("No templates found in %s\n", targetTemplatesDir)
		return nil
	}

	// List template files
	files, err := os.ReadDir(targetTemplatesDir)
	if err != nil {
		logger.Error("Failed to read templates directory", err,
			logging.String("dir", targetTemplatesDir))
		return fmt.Errorf("failed to read templates directory: %w", err)
	}

	if len(files) == 0 {
		fmt.Printf("No templates found in %s\n", targetTemplatesDir)
		return nil
	}

	fmt.Printf("Templates found in %s:\n", targetTemplatesDir)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".tmpl") {
			filePath := filepath.Join(targetTemplatesDir, file.Name())
			info, err := os.Stat(filePath)
			if err == nil {
				fmt.Printf("  %s (%d bytes)\n", file.Name(), info.Size())
			} else {
				fmt.Printf("  %s\n", file.Name())
			}
		}
	}

	return nil
}

// listProjectFacts lists facts in a spooky project
func listProjectFacts(logger logging.Logger, path, _ string) error {
	logger.Info("Listing project facts",
		logging.String("path", path))

	// Create template context to get project info
	ctx, err := NewTemplateContext(logger, path)
	if err != nil {
		return fmt.Errorf("failed to create template context: %w", err)
	}

	// Create fact manager
	manager := facts.NewManager(nil)

	// Get all facts
	allFacts, err := manager.GetAllFacts()
	if err != nil {
		logger.Error("Failed to retrieve facts", err)
		return fmt.Errorf("failed to retrieve facts: %w", err)
	}

	if len(allFacts) == 0 {
		fmt.Printf("No facts found for project %s\n", ctx.Project.Name)
		fmt.Printf("Use 'spooky project gather-facts %s' to collect facts first.\n", path)
		return nil
	}

	// Group facts by server
	serverFacts := make(map[string][]*facts.Fact)
	for _, fact := range allFacts {
		serverFacts[fact.Server] = append(serverFacts[fact.Server], fact)
	}

	// Display facts in table format
	fmt.Printf("Facts for project %s:\n", ctx.Project.Name)
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

	logger.Info("Project facts listed successfully",
		logging.String("project", ctx.Project.Name),
		logging.Int("total_facts", len(allFacts)),
		logging.Int("servers", len(serverFacts)))
	return nil
}

// gatherProjectFacts gathers facts for machines in a spooky project
func gatherProjectFacts(logger logging.Logger, path, _ string, _ bool, _ string) error {
	logger.Info("Gathering project facts",
		logging.String("path", path))

	// Create template context to get project info and machines
	ctx, err := NewTemplateContext(logger, path)
	if err != nil {
		return fmt.Errorf("failed to create template context: %w", err)
	}

	// Create fact manager
	manager := facts.NewManager(nil)

	// Get all machines from the project
	machines := ctx.Machines
	if len(machines) == 0 {
		fmt.Printf("No machines found in project %s\n", ctx.Project.Name)
		return nil
	}

	fmt.Printf("Gathering facts for %d machines in project %s...\n", len(machines), ctx.Project.Name)

	// Collect facts from each machine
	var allCollections []*facts.FactCollection
	var errors []error

	for _, machine := range machines {
		logger.Info("Collecting facts from machine",
			logging.String("machine", machine.Name),
			logging.String("host", machine.Host))

		collection, err := manager.CollectAllFacts(machine.Host)
		if err != nil {
			logger.Error("Failed to collect facts from machine", err,
				logging.String("machine", machine.Name),
				logging.String("host", machine.Host))
			errors = append(errors, fmt.Errorf("failed to collect facts from %s (%s): %w", machine.Name, machine.Host, err))
			continue
		}

		allCollections = append(allCollections, collection)
		fmt.Printf("✓ Collected %d facts from %s (%s)\n", len(collection.Facts), machine.Name, machine.Host)
	}

	// Display summary
	fmt.Printf("\nFact Gathering Summary:\n")
	fmt.Printf("Project: %s\n", ctx.Project.Name)
	fmt.Printf("Machines processed: %d\n", len(allCollections))
	fmt.Printf("Machines failed: %d\n", len(errors))
	fmt.Printf("Total facts collected: %d\n", getTotalFactCount(allCollections))

	if len(errors) > 0 {
		fmt.Printf("\nErrors:\n")
		for _, err := range errors {
			fmt.Printf("  - %v\n", err)
		}
	}

	logger.Info("Project fact gathering completed",
		logging.String("project", ctx.Project.Name),
		logging.Int("machines_processed", len(allCollections)),
		logging.Int("machines_failed", len(errors)),
		logging.Int("total_facts", getTotalFactCount(allCollections)))

	return nil
}

// renderProjectTemplate renders a template in the context of a spooky project
func renderProjectTemplate(logger logging.Logger, templateFile, path, output string, dryRun bool, server, _, templatesDir, _ string) error {
	logger.Info("Rendering project template",
		logging.String("template", templateFile),
		logging.String("path", path),
		logging.String("output", output),
		logging.Bool("dry_run", dryRun),
		logging.String("server", server))

	// Create template context
	ctx, err := NewTemplateContext(logger, path)
	if err != nil {
		return fmt.Errorf("failed to create template context: %w", err)
	}

	// Load server-specific facts if server is specified
	if server != "" {
		if err := ctx.LoadServerFacts(logger, server); err != nil {
			logger.Warn("Failed to load server facts",
				logging.String("server", server),
				logging.String("error", err.Error()))
		}
	}

	// Check if template file exists
	var templatePath string
	if templatesDir != "" {
		templatePath = filepath.Join(templatesDir, templateFile)
	} else {
		templatePath = filepath.Join(path, templateFile)
	}

	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		logger.Error("Template file not found", err,
			logging.String("file", templatePath))
		return fmt.Errorf("template file not found: %s", templatePath)
	}

	// Read template content
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		logger.Error("Failed to read template file", err,
			logging.String("file", templatePath))
		return fmt.Errorf("failed to read template file: %w", err)
	}

	// Create template engine with context
	tmpl, err := template.New("project-template").Funcs(ctx.GetTemplateFunctions()).Parse(string(templateContent))
	if err != nil {
		logger.Error("Failed to parse template", err,
			logging.String("file", templatePath))
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		logger.Error("Failed to execute template", err,
			logging.String("file", templatePath))
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Handle output
	if dryRun {
		fmt.Printf("DRY RUN - Would render template '%s' to:\n", templateFile)
		fmt.Println("---")
		fmt.Print(buf.String())
		fmt.Println("---")
		return nil
	}

	if output == "" {
		// Output to stdout
		fmt.Print(buf.String())
	} else {
		// Output to file
		if err := os.WriteFile(output, buf.Bytes(), 0o600); err != nil {
			logger.Error("Failed to write output file", err,
				logging.String("output", output))
			return fmt.Errorf("failed to write output file: %w", err)
		}
		logger.Info("Template rendered successfully",
			logging.String("template", templateFile),
			logging.String("output", output))
	}

	return nil
}

// validateProjectTemplate validates a template in the context of a spooky project
func validateProjectTemplate(logger logging.Logger, templateFile, path, templatesDir, _ string) error {
	logger.Info("Validating project template",
		logging.String("template", templateFile),
		logging.String("path", path))

	// Create template context
	ctx, err := NewTemplateContext(logger, path)
	if err != nil {
		return fmt.Errorf("failed to create template context: %w", err)
	}

	// Check if template file exists
	var templatePath string
	if templatesDir != "" {
		templatePath = filepath.Join(templatesDir, templateFile)
	} else {
		templatePath = filepath.Join(path, templateFile)
	}

	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		logger.Error("Template file not found", err,
			logging.String("file", templatePath))
		return fmt.Errorf("template file not found: %s", templatePath)
	}

	// Read template content
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		logger.Error("Failed to read template file", err,
			logging.String("file", templatePath))
		return fmt.Errorf("failed to read template file: %w", err)
	}

	// Create template engine with context
	tmpl, err := template.New("project-template").Funcs(ctx.GetTemplateFunctions()).Parse(string(templateContent))
	if err != nil {
		logger.Error("Failed to parse template", err,
			logging.String("file", templatePath))
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Try to execute with empty data to validate syntax
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		logger.Error("Template validation failed", err,
			logging.String("file", templatePath))
		return fmt.Errorf("template validation failed: %w", err)
	}

	logger.Info("Template validation successful",
		logging.String("file", templatePath))
	fmt.Printf("✓ Template '%s' is valid\n", templateFile)

	return nil
}
