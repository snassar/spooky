// Package cmd provides command implementations for spooky CLI.
package cmd

import (
	"context"
	"fmt"

	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookyproject "spooky/internal/project"
	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"

	"github.com/spf13/cobra"
)

// Global instances for dependency injection
var (
	projectManager   spookyinterfaces.ProjectManager
	projectValidator spookyinterfaces.ProjectValidator
	projectLoader    spookyinterfaces.ProjectLoader
	logger           spookytypeslogging.Logger
)

// InitializeProjectDependencies initializes the project dependencies
func InitializeProjectDependencies() error {
	// Initialize logger
	logManager := spookylogging.NewLogManager()
	logger = logManager.GetLogger("project")

	// Initialize project components
	projectValidator = spookyproject.NewValidator(logger)
	projectLoader = spookyproject.NewLoader(logger)
	projectManager = spookyproject.NewManager(logger, projectValidator, projectLoader)

	return nil
}

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage spooky projects",
	Long: `Manage spooky projects including initialization and validation.

A spooky project contains configuration files, templates, and other resources
needed for automation and orchestration tasks.`,
}

// projectInitCmd represents the project init command
var projectInitCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new spooky project",
	Long: `Initialize a new spooky project with the specified name.

This command creates a new spooky project directory with the required structure
and configuration files according to the project-directory.schema.hcl schema.

Flags allow customization of project metadata during initialization.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		version, _ := cmd.Flags().GetString("version")
		author, _ := cmd.Flags().GetString("author")
		email, _ := cmd.Flags().GetString("email")
		url, _ := cmd.Flags().GetString("url")

		return handleProjectInit(args[0], name, description, version, author, email, url)
	},
}

// projectValidateCmd represents the project validate command
var projectValidateCmd = &cobra.Command{
	Use:   "validate [project-path]",
	Short: "Validate a spooky project",
	Long: `Validate a spooky project structure and configuration.

This command validates that the project follows the project-directory.schema.hcl
schema and that all configuration files are properly formatted.`,
	Args: cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		return handleProjectValidate(args[0])
	},
}

// projectEncryptCmd represents the project encrypt command
var projectEncryptCmd = &cobra.Command{
	Use:   "encrypt [project-path]",
	Short: "Encrypt all variables and machines in a project",
	Long: `Encrypt all variables and machines in a project that have encrypted=true.

This command processes both variables and machines, encrypting any that have
encrypted=true set. It will re-encrypt if identities/recipients have changed.

Examples:
  spooky project encrypt ./my-project
  spooky project encrypt ./my-project --dry-run`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath := args[0]
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		return handleProjectEncrypt(projectPath, dryRun)
	},
}

// Helper function for dependency initialization
func initializeProjectDependenciesIfNeeded() error {
	if projectManager == nil {
		if err := InitializeProjectDependencies(); err != nil {
			return fmt.Errorf("failed to initialize project dependencies: %w", err)
		}
	}
	return nil
}

// Helper function for project initialization
func createProject(ctx context.Context, projectPath string) (*spookytypes.Project, error) {
	project, err := projectManager.Initialize(ctx, projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize project: %w", err)
	}
	return project, nil
}

// Helper function to check if configuration updates are needed
func hasConfigurationUpdates(name, description, version, author, url string) bool {
	return name != "" || description != "" || version != "" || author != "" || url != ""
}

// Helper function for updating project configuration
func updateProjectConfiguration(project *spookytypes.Project, name, description string) {
	if name != "" {
		project.Config.Name = name
	}
	if description != "" {
		project.Config.Description = description
	}
}

// Helper function for updating project metadata
func updateProjectMetadata(project *spookytypes.Project, version, author, url string) {
	if project.Config.Metadata == nil {
		project.Config.Metadata = &spookytypes.ProjectMetadata{}
	}

	if version != "" {
		project.Config.Metadata.Version = version
	}
	if author != "" {
		project.Config.Metadata.Author = author
	}
	if url != "" {
		project.Config.Metadata.URL = url
	}
}

// Helper function for saving project configuration
func saveProjectConfiguration(ctx context.Context, project *spookytypes.Project) error {
	if err := projectManager.Save(ctx, project); err != nil {
		return fmt.Errorf("failed to save updated project configuration: %w", err)
	}
	return nil
}

// Helper function for displaying success output
func displayProjectInitSuccess(projectPath string) {
	fmt.Printf("✅ Project initialized successfully: %s\n", projectPath)
	fmt.Printf("📁 Project structure created according to project-directory.schema.hcl\n")
	fmt.Printf("📄 Configuration files generated using project.schema.hcl\n")
	fmt.Printf("💡 Next steps:\n")
	fmt.Printf("   - Edit project.hcl to customize your project\n")
	fmt.Printf("   - Add machines.hcl for machine inventory\n")
	fmt.Printf("   - Add actions.hcl for automation tasks\n")
	fmt.Printf("   - Add variables.hcl for project variables\n")
}

// handleProjectInit handles project initialization using the ProjectManager interface
func handleProjectInit(projectPath, name, description, version, author, _, url string) error {
	ctx := context.Background()

	// Initialize dependencies
	if err := initializeProjectDependenciesIfNeeded(); err != nil {
		return err
	}

	// Create project
	project, err := createProject(ctx, projectPath)
	if err != nil {
		return err
	}

	// Update configuration if needed
	if hasConfigurationUpdates(name, description, version, author, url) {
		updateProjectConfiguration(project, name, description)
		updateProjectMetadata(project, version, author, url)

		if err := saveProjectConfiguration(ctx, project); err != nil {
			return err
		}
	}

	// Display success output
	displayProjectInitSuccess(projectPath)

	return nil
}

// handleProjectValidate handles project validation using the ProjectManager interface
func handleProjectValidate(projectPath string) error {
	ctx := context.Background()

	// Initialize dependencies if not already done
	if projectManager == nil {
		if err := InitializeProjectDependencies(); err != nil {
			return fmt.Errorf("failed to initialize project dependencies: %w", err)
		}
	}

	// Load the project using the manager
	project, err := projectManager.Load(ctx, projectPath)
	if err != nil {
		return fmt.Errorf("failed to load project: %w", err)
	}

	fmt.Printf("🔍 Validating project: %s\n", project.Path)

	// Validate the project using the manager
	result, err := projectManager.Validate(ctx, project)
	if err != nil {
		return fmt.Errorf("failed to validate project: %w", err)
	}

	// Report validation results
	if result.Valid {
		fmt.Printf("✅ Project validation passed - all required components present\n")
		fmt.Printf("📋 Schema compliance: project-directory.schema.hcl ✅\n")
		fmt.Printf("📋 Schema compliance: project.schema.hcl ✅\n")
	} else {
		if len(result.Errors) > 0 {
			fmt.Printf("❌ Validation issues found:\n")
			for i := range result.Errors {
				err := &result.Errors[i]
				fmt.Printf("   - %s\n", err.Message)
			}
		}
		if len(result.Warnings) > 0 {
			fmt.Printf("⚠️  Warnings:\n")
			for i := range result.Warnings {
				warning := &result.Warnings[i]
				fmt.Printf("   - %s\n", warning.Message)
			}
		}
	}

	return nil
}

// handleProjectEncrypt handles project-wide encryption
func handleProjectEncrypt(projectPath string, dryRun bool) error {
	secretsIntegration, err := setupEncryption()
	if err != nil {
		return err
	}

	manager := GetIntegrationManager()
	variablesIntegration := manager.GetVariablesIntegration()
	machinesIntegration := manager.GetMachinesIntegration()

	if variablesIntegration == nil || machinesIntegration == nil {
		return fmt.Errorf("required integrations not available")
	}

	fmt.Printf("Project encryption for %s (dry-run: %t)\n", projectPath, dryRun)

	recipients, err := loadRecipients(secretsIntegration)
	if err != nil {
		return err
	}

	fmt.Printf("Loaded %d recipients for encryption\n", len(recipients))

	// Process variables
	if err := variablesIntegration.EncryptVariables(context.Background(), projectPath, secretsIntegration, recipients, dryRun); err != nil {
		return fmt.Errorf("failed to encrypt variables: %w", err)
	}

	// Process machines
	if err := machinesIntegration.EncryptMachines(context.Background(), projectPath, secretsIntegration, recipients, dryRun); err != nil {
		return fmt.Errorf("failed to encrypt machines: %w", err)
	}

	if dryRun {
		fmt.Println("Dry run completed - no changes made")
	} else {
		fmt.Println("Project encryption completed successfully")
	}

	return nil
}

func init() {
	// Add flags to project init command
	projectInitCmd.Flags().String("name", "", "Project name (defaults to directory name)")
	projectInitCmd.Flags().String("description", "", "Project description")
	projectInitCmd.Flags().String("version", "", "Project version")
	projectInitCmd.Flags().String("author", "", "Project author")
	projectInitCmd.Flags().String("email", "", "Project email")
	projectInitCmd.Flags().String("url", "", "Project URL")

	// Add flags to project encrypt command
	projectEncryptCmd.Flags().Bool("dry-run", false, "Show what would be encrypted without making changes")

	projectCmd.AddCommand(projectInitCmd)
	projectCmd.AddCommand(projectValidateCmd)
	projectCmd.AddCommand(projectEncryptCmd)
	RootCmd.AddCommand(projectCmd)
}
