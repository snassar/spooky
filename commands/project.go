package commands

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"spooky/internal/encryption"
	"spooky/internal/logging"
	"spooky/internal/schemas"
	"spooky/internal/utilities"

	"github.com/spf13/cobra"
)

var (
	projectName        string
	projectDescription string
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage spooky projects",
	Long: `Manage spooky projects including initialization, configuration, and management.

A spooky project is a directory containing configuration files that define
automation tasks, machine inventory, and deployment configurations.`,
}

var projectInitCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initialize a new spooky project",
	Long: `Initialize a new spooky project in the specified directory.

This command creates a new project directory with the necessary configuration
files including project.hcl, machines.hcl, actions.hcl, and variables.hcl.

The directory will be created if it doesn't exist.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runProjectInit,
}

var projectValidateCmd = &cobra.Command{
	Use:   "validate [directory]",
	Short: "Validate a spooky project",
	Long: `Validate a spooky project directory structure and configuration files.

This command checks:
- Project directory structure compliance
- HCL file syntax validation
- Schema compliance for machines, actions, and variables
- File and directory existence rules

The directory defaults to the current directory if not specified.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runProjectValidate,
}

var projectEncryptCmd = &cobra.Command{
	Use:   "encrypt [directory]",
	Short: "Encrypt sensitive values in a spooky project",
	Long: `Encrypt sensitive values in a spooky project using age encryption.

This command:
- Finds variables marked with encrypted = true
- Encrypts their plaintext values using age encryption
- Updates the HCL files with encrypted values
- Requires age identities and recipients to be configured

The directory defaults to the current directory if not specified.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runProjectEncrypt,
}

var projectConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current configuration information",
	Long: `Show information about the current spooky configuration.

This command displays:
- Which configuration file is being used
- Configuration source (custom, user, or embedded default)
- Age encryption settings
- Configuration file details`,
	RunE: runProjectConfig,
}

func init() {
	// Add project command to root
	RootCmd.AddCommand(projectCmd)

	// Add init subcommand to project
	projectCmd.AddCommand(projectInitCmd)

	// Add validate subcommand to project
	projectCmd.AddCommand(projectValidateCmd)

	// Add encrypt subcommand to project
	projectCmd.AddCommand(projectEncryptCmd)

	// Add config subcommand to project
	projectCmd.AddCommand(projectConfigCmd)

	// Add flags for project init
	projectInitCmd.Flags().StringVar(&projectName, "name", "", "Project name (required)")
	projectInitCmd.Flags().StringVar(&projectDescription, "description", "", "Project description")

	// Mark name as required
	if err := projectInitCmd.MarkFlagRequired("name"); err != nil {
		// Log error but continue - this is a configuration issue that shouldn't prevent the program from running
		slog.Warn("Failed to mark name flag as required", "error", err)
	}
}

func runProjectInit(cmd *cobra.Command, args []string) error {
	// Determine target directory
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	// Validate project name
	if err := validateProjectName(projectName); err != nil {
		return fmt.Errorf("invalid project name: %w", err)
	}

	// Create target directory if it doesn't exist
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
	}

	// Create project.hcl
	if err := createProjectHCL(targetDir, projectName, projectDescription); err != nil {
		return fmt.Errorf("failed to create project.hcl: %w", err)
	}

	// Create machines.hcl
	if err := createMachinesHCL(targetDir); err != nil {
		return fmt.Errorf("failed to create machines.hcl: %w", err)
	}

	// Create actions.hcl
	if err := createActionsHCL(targetDir); err != nil {
		return fmt.Errorf("failed to create actions.hcl: %w", err)
	}

	// Create variables.hcl
	if err := createVariablesHCL(targetDir); err != nil {
		return fmt.Errorf("failed to create variables.hcl: %w", err)
	}

	// Create README.md
	if err := createREADME(targetDir, projectName, projectDescription); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	logger := logging.GetGlobalLogger()
	logger.Info("✅ Successfully initialized spooky project",
		slog.String("project_name", projectName),
		slog.String("target_directory", targetDir))

	logger.Info("📁 Project files created",
		slog.String("project_hcl", "project.hcl"),
		slog.String("machines_hcl", "machines.hcl"),
		slog.String("actions_hcl", "actions.hcl"),
		slog.String("variables_hcl", "variables.hcl"),
		slog.String("readme", "README.md"))

	logger.Info("🚀 Next steps",
		slog.String("step1", fmt.Sprintf("cd %s", targetDir)),
		slog.String("step2", "Edit the configuration files as needed"),
		slog.String("step3", "Run 'spooky project validate' to check configuration"))

	return nil
}

func validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	if len(name) > 128 {
		return fmt.Errorf("project name too long (max 128 characters)")
	}

	// Check if name matches the pattern from schema: ^[a-zA-Z][a-zA-Z0-9._-]*$
	if !regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9._-]*$").MatchString(name) {
		return fmt.Errorf("project name must start with a letter and contain only letters, numbers, dots, underscores, and hyphens")
	}

	return nil
}

// createProjectConfigFile creates a project configuration file with the specified content and filename
func createProjectConfigFile(targetDir, filename, content string) error {
	return os.WriteFile(filepath.Join(targetDir, filename), []byte(content), 0o644)
}

func createProjectHCL(targetDir, name, description string) error {
	// Generate project configuration directly from Go structs
	content, err := schemas.GenerateProjectConfigFromStructs(name, description)
	if err != nil {
		return fmt.Errorf("failed to generate project config: %w", err)
	}
	return createProjectConfigFile(targetDir, "project.hcl", content)
}

func createMachinesHCL(targetDir string) error {
	// Generate machines configuration directly from Go structs
	content := schemas.GenerateMachinesConfigFromStructs()
	return createProjectConfigFile(targetDir, "machines.hcl", content)
}

func createActionsHCL(targetDir string) error {
	// Generate actions configuration directly from Go structs
	content := schemas.GenerateActionsConfigFromStructs()
	return createProjectConfigFile(targetDir, "actions.hcl", content)
}

func createVariablesHCL(targetDir string) error {
	// Generate variables configuration directly from Go structs
	content := schemas.GenerateVariablesConfigFromStructs()
	return createProjectConfigFile(targetDir, "variables.hcl", content)
}

func createREADME(targetDir, name, description string) error {
	content := fmt.Sprintf("# %s\n\n%s\n\n## Overview\n\nThis is a Spooky automation project that defines configuration management, \ndeployment automation, and infrastructure management tasks.\n\n## Project Structure\n\n- project.hcl - Project configuration and metadata\n- machines.hcl - Machine inventory and connectivity settings\n- actions.hcl - Automation tasks and deployment actions\n- variables.hcl - Project-wide variables and configuration values\n- templates/ - Template files for deployment (create as needed)\n- files/ - Static files for deployment (create as needed)\n\n## Getting Started\n\n1. **Configure Machines**: Edit machines.hcl to define your target machines\n2. **Define Actions**: Edit actions.hcl to create automation tasks\n3. **Set Variables**: Edit variables.hcl to configure project variables\n4. **Validate**: Run 'spooky project validate' to check configuration\n5. **Execute**: Run 'spooky run <action-name>' to execute actions\n\n## Examples\n\n### Running Actions\n```bash\n# Run a specific action\nspooky run deploy-application\n\n# Run actions with specific tags\nspooky run --tags deployment\n\n# Dry run to see what would happen\nspooky run --dry-run deploy-application\n```\n\n### Managing Machines\n```bash\n# List all machines\nspooky machines list\n\n# Test connectivity\nspooky machines test-connection\n\n# Collect facts\nspooky machines collect-facts\n```\n\n## Documentation\n\nFor more information about Spooky, visit the project documentation or run:\n```bash\nspooky --help\nspooky project --help\nspooky machines --help\nspooky actions --help\n```\n\n## Support\n\nIf you encounter issues or have questions, please refer to the Spooky documentation\nor create an issue in the project repository.", name, description)

	return os.WriteFile(filepath.Join(targetDir, "README.md"), []byte(content), 0o644)
}

func runProjectValidate(cmd *cobra.Command, args []string) error {
	// Determine target directory
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	// Check if directory exists
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", targetDir)
	}

	// Use the enhanced project validator with schema-driven validation
	validator := utilities.NewProjectValidator()
	result := validator.ValidateProject(targetDir)

	// Display validation result
	if result.IsValid {
		// Try to get project name for display, but don't fail if it's missing
		projectName, _ := getProjectName(targetDir)
		if projectName == "" {
			projectName = "unknown"
		}
		logger := logging.GetGlobalLogger()
		logger.Info("Project validation successful", slog.String("project_name", projectName))
	} else {
		logger := logging.GetGlobalLogger()
		logger.Error("project validation failed",
			slog.Int("error_count", len(result.Errors)))
		for _, err := range result.Errors {
			logger.Error("validation error",
				slog.String("message", err.Message))
		}
		return fmt.Errorf("project validation failed with %d errors", len(result.Errors))
	}

	return nil
}

// getProjectName extracts the project name from project.hcl file
func getProjectName(targetDir string) (string, error) {
	projectHCLPath := filepath.Join(targetDir, "project.hcl")
	content, err := os.ReadFile(projectHCLPath)
	if err != nil {
		return "", nil // Return empty string instead of error for missing files
	}

	// Regex to extract project name from "project \"name\" {"
	re := regexp.MustCompile(`project\s+["']([^"']+)["']\s*{`)
	matches := re.FindStringSubmatch(string(content))
	if len(matches) < 2 {
		return "unknown", nil // Return default name instead of error
	}

	return matches[1], nil
}

func runProjectEncrypt(cmd *cobra.Command, args []string) error {
	// Determine target directory
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	// Validate project directory
	validator := utilities.NewProjectValidator()
	result := validator.ValidateProject(targetDir)
	if !result.IsValid {
		return fmt.Errorf("project directory validation failed: %s", result.Errors[0].Message)
	}

	// Get age configuration from spooky config
	configManager, err := utilities.NewConfigManager()
	if err != nil {
		return fmt.Errorf("failed to create config manager: %w", err)
	}

	// Get effective config to find age settings (respect --config flag)
	customConfigFile := GetConfigFile()
	_, err = configManager.GetEffectiveConfig(customConfigFile)
	if err != nil {
		return fmt.Errorf("failed to get effective config: %w", err)
	}

	// Parse config to find age settings
	// For now, we'll use default paths - in a real implementation,
	// you'd parse the HCL config to get the actual paths
	identitiesPath := "~/.config/spooky/age/identities"
	recipientsPath := "~/.config/spooky/age/recipients"

	// Expand home directory
	if strings.HasPrefix(identitiesPath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		identitiesPath = filepath.Join(homeDir, identitiesPath[2:])
	}

	if strings.HasPrefix(recipientsPath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		recipientsPath = filepath.Join(homeDir, recipientsPath[2:])
	}

	// Create age encryption instance
	ageEncryption, err := encryption.NewAgeEncryption(identitiesPath, recipientsPath)
	if err != nil {
		return fmt.Errorf("failed to create age encryption: %w", err)
	}

	// Validate age configuration
	if err := ageEncryption.ValidateConfiguration(); err != nil {
		return fmt.Errorf("age encryption configuration invalid: %w", err)
	}

	logger := logging.GetGlobalLogger()
	logger.Info("Age encryption configured",
		slog.Int("identities_count", ageEncryption.GetIdentitiesCount()),
		slog.Int("recipients_count", ageEncryption.GetRecipientsCount()))

	// Create HCL updater
	updater := encryption.NewHCLUpdater(ageEncryption)

	// Process the project directory
	logger.Info("Processing project directory", slog.String("target_directory", targetDir))
	if err := updater.UpdateDirectory(targetDir); err != nil {
		return fmt.Errorf("failed to process project directory: %w", err)
	}

	logger.Info("Project encryption completed successfully!")
	return nil
}

func runProjectConfig(cmd *cobra.Command, args []string) error {
	// Get custom config file from --config flag
	customConfigFile := GetConfigFile()

	logger := logging.GetGlobalLogger()
	logger.Info("=== Spooky Configuration Information ===")

	// Create config manager
	configManager, err := utilities.NewConfigManager()
	if err != nil {
		return fmt.Errorf("failed to create config manager: %w", err)
	}

	// Get effective config info
	effectiveInfo, err := configManager.GetEffectiveConfigInfo(customConfigFile)
	if err != nil {
		return fmt.Errorf("failed to get effective config info: %w", err)
	}

	// Display configuration information
	logger.Info("Configuration information",
		slog.String("source", effectiveInfo.Source),
		slog.String("config_file", effectiveInfo.ConfigFile))

	if effectiveInfo.Exists {
		logger.Info("Configuration file details",
			slog.Int64("size_bytes", effectiveInfo.Size))
		if effectiveInfo.ModTime != "" {
			logger.Info("Configuration file modified", slog.String("modified_time", effectiveInfo.ModTime))
		}
	} else {
		logger.Info("Configuration file status", slog.String("status", "Does not exist"))
	}

	// Show config priority
	logger.Info("Configuration priority information",
		slog.String("priority1", "Custom config file (--config flag)"),
		slog.String("priority2", "User config (~/.config/spooky/spooky.hcl)"),
		slog.String("priority3", "Embedded default"))

	// Show current --config flag status
	logger.Info("Configuration flag status")
	if customConfigFile != "" {
		logger.Info("Custom config file specified", slog.String("config_file", customConfigFile))
	} else {
		logger.Info("Custom config file not specified")
	}

	return nil
}
