// Package cmd provides command implementations for spooky CLI.
package cmd

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesvariables "spooky/internal/types/variables"
	spookyvariables "spooky/internal/variables"

	"github.com/spf13/cobra"
)

// Global instances for variable dependency injection
var (
	variablesManager spookyinterfaces.VariablesIntegration
	variablesLogger  spookytypeslogging.Logger
)

// InitializeVariablesDependencies initializes variable-related dependencies
func InitializeVariablesDependencies() error {
	// Create log manager for variables component
	logManager := spookylogging.NewLogManager()
	variablesLogger = logManager.GetLogger("variables")

	// Initialize variable components
	validator := spookyvariables.NewValidator(variablesLogger)
	loader := spookyvariables.NewLoader(variablesLogger)
	variablesManager = spookyvariables.NewManager(variablesLogger, loader, validator)

	return nil
}

// variablesCmd represents the variables command
var variablesCmd = &cobra.Command{
	Use:   "variables",
	Short: "Manage project variables",
	Long: `Manage project variables including listing, validation, and resolution.

Variables are defined in variables.hcl files or variables/*.hcl files within spooky projects
and provide dynamic configuration values for templates and actions.`,
}

// variablesListCmd represents the variables list command
var variablesListCmd = &cobra.Command{
	Use:   "list [project-path]",
	Short: "List variables in a project",
	Long: `List all variables defined in the project's variable files.

This command reads variables.hcl files and variables/*.hcl files and displays information
about all configured variables including name, type, description, and scope.`,
	Args: cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		return handleVariablesList(args[0])
	},
}

// variablesValidateCmd represents the variables validate command
var variablesValidateCmd = &cobra.Command{
	Use:   "validate [project-path]",
	Short: "Validate variable definitions",
	Long: `Validate variable definitions and dependencies.

This command validates that all variables in the project have proper configuration
including required fields, valid types, and dependency relationships.`,
	Args: cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		return handleVariablesValidate(args[0])
	},
}

// variablesResolveCmd represents the variables resolve command
var variablesResolveCmd = &cobra.Command{
	Use:   "resolve [project-path]",
	Short: "Resolve variables with context",
	Long: `Resolve variables with the given context and display resolved values.

This command loads variables from the project and resolves them using the provided
context including environment variables, facts, and machine data.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleVariablesResolve(cmd, args[0])
	},
}

// variablesEncryptCmd represents the variables encrypt command
var variablesEncryptCmd = &cobra.Command{
	Use:   "encrypt [project-path]",
	Short: "Encrypt variables in a project",
	Long: `Encrypt all variables in a project that have encrypted=true.

This command processes variables files and encrypts any variables that have
encrypted=true set. It will re-encrypt if identities/recipients have changed.

Examples:
  spooky variables encrypt ./my-project
  spooky variables encrypt ./my-project --dry-run`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath := args[0]
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		return handleVariablesEncrypt(projectPath, dryRun)
	},
}

// variablesArmorCmd represents the variables armor command
var variablesArmorCmd = &cobra.Command{
	Use:   "armor [project-path]",
	Short: "Armor (encrypt) variables in a project",
	Long: `Armor (encrypt) variables in a project using age encryption.

This command processes variables files and encrypts any variables that have
encrypted=true set. It will re-encrypt if identities/recipients have changed.

Examples:
  spooky variables armor ./my-project
  spooky variables armor ./my-project --dry-run`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath := args[0]
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		return handleVariablesArmor(projectPath, dryRun)
	},
}

// handleVariablesList handles listing variables using the VariablesIntegration interface
func handleVariablesList(projectPath string) error {
	ctx := context.Background()

	// Initialize dependencies if not already done
	if variablesManager == nil {
		if err := InitializeVariablesDependencies(); err != nil {
			return fmt.Errorf("failed to initialize variables dependencies: %w", err)
		}
	}

	// Load variables from project
	variables, err := variablesManager.LoadVariables(ctx, projectPath)
	if err != nil {
		return fmt.Errorf("failed to load variables: %w", err)
	}

	// Display variables in a structured format
	fmt.Printf("Variables in project: %s\n", projectPath)
	fmt.Printf("Total variables: %d\n\n", len(variables))

	if len(variables) == 0 {
		fmt.Println("No variables found.")
		return nil
	}

	// Group variables by source file
	sourceGroups := make(map[string][]*spookytypesvariables.Variable)
	for _, variable := range variables {
		source := variable.SourceFile
		if source == "" {
			source = "unknown"
		}
		sourceGroups[source] = append(sourceGroups[source], variable)
	}

	// Display variables grouped by source
	for source, vars := range sourceGroups {
		fmt.Printf("Source: %s (%d variables)\n", source, len(vars))
		fmt.Println(strings.Repeat("-", 80))

		for _, variable := range vars {
			fmt.Printf("  %s (%s)", variable.Name, variable.Type)
			if variable.Description != "" {
				fmt.Printf(" - %s", variable.Description)
			}
			if variable.Required {
				fmt.Printf(" [required]")
			}
			if variable.Sensitive {
				fmt.Printf(" [sensitive]")
			}
			if variable.Encrypted {
				fmt.Printf(" [encrypted]")
			}
			if len(variable.Dependencies) > 0 {
				fmt.Printf(" [deps: %s]", strings.Join(variable.Dependencies, ", "))
			}
			fmt.Println()
		}
		fmt.Println()
	}

	return nil
}

// handleVariablesValidate handles validating variables using the VariablesIntegration interface
func handleVariablesValidate(projectPath string) error {
	ctx := context.Background()

	// Initialize dependencies if not already done
	if variablesManager == nil {
		if err := InitializeVariablesDependencies(); err != nil {
			return fmt.Errorf("failed to initialize variables dependencies: %w", err)
		}
	}

	// Load variables from project
	variables, err := variablesManager.LoadVariables(ctx, projectPath)
	if err != nil {
		return fmt.Errorf("failed to load variables: %w", err)
	}

	// Validate variables
	result, err := variablesManager.ValidateVariables(ctx, variables)
	if err != nil {
		return fmt.Errorf("failed to validate variables: %w", err)
	}

	// Display validation results
	fmt.Printf("Variable validation for project: %s\n", projectPath)
	fmt.Printf("Total variables: %d\n", len(variables))

	if result.Valid {
		fmt.Println("✅ All variables are valid!")
	} else {
		fmt.Printf("❌ Validation failed with %d errors\n", len(result.Errors))
	}

	// Display errors
	if len(result.Errors) > 0 {
		fmt.Println("\nErrors:")
		for i := range result.Errors {
			err := &result.Errors[i]
			fmt.Printf("  %d. %s\n", i+1, err.Message)
		}
	}

	// Display warnings
	if len(result.Warnings) > 0 {
		fmt.Println("\nWarnings:")
		for i := range result.Warnings {
			warning := &result.Warnings[i]
			fmt.Printf("  %d. %s\n", i+1, warning.Message)
		}
	}

	// Return error if validation failed
	if !result.Valid {
		return fmt.Errorf("variable validation failed")
	}

	return nil
}

// Helper function for dependency initialization
func initializeVariablesDependenciesIfNeeded() error {
	if variablesManager == nil {
		if err := InitializeVariablesDependencies(); err != nil {
			return fmt.Errorf("failed to initialize variables dependencies: %w", err)
		}
	}
	return nil
}

// Helper function for loading variables
func loadProjectVariables(ctx context.Context, projectPath string) (map[string]*spookytypesvariables.Variable, error) {
	variables, err := variablesManager.LoadVariables(ctx, projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load variables: %w", err)
	}
	return variables, nil
}

// Helper function for creating variable context
func createVariableContext(projectPath string) *spookytypesvariables.VariableContext {
	return &spookytypesvariables.VariableContext{
		ProjectPath: projectPath,
		Environment: make(map[string]interface{}),
		Facts:       make(map[string]interface{}),
		Machines:    make(map[string]interface{}),
		UserData:    make(map[string]interface{}),
		Timestamp:   time.Now(),
	}
}

// Helper function for resolving variables
func resolveProjectVariables(ctx context.Context, variables map[string]*spookytypesvariables.Variable, variableContext *spookytypesvariables.VariableContext) (*spookytypesvariables.VariableResolutionResult, error) {
	resolutionResult, err := variablesManager.ResolveVariables(ctx, variables, variableContext)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve variables: %w", err)
	}
	return resolutionResult, nil
}

// Helper function for displaying resolution summary
func displayResolutionSummary(projectPath string, variables map[string]*spookytypesvariables.Variable, resolutionResult *spookytypesvariables.VariableResolutionResult) {
	fmt.Printf("Variable resolution for project: %s\n", projectPath)
	fmt.Printf("Total variables: %d\n", len(variables))
	fmt.Printf("Resolved variables: %d\n", len(resolutionResult.Resolved))
	fmt.Printf("Resolution time: %v\n", resolutionResult.Duration)
}

// Helper function for displaying errors
func displayResolutionErrors(resolutionResult *spookytypesvariables.VariableResolutionResult) {
	if len(resolutionResult.Errors) > 0 {
		fmt.Printf("\n❌ Resolution errors (%d):\n", len(resolutionResult.Errors))
		for i := range resolutionResult.Errors {
			err := &resolutionResult.Errors[i]
			fmt.Printf("  %d. %s: %s\n", i+1, err.VariableName, err.ErrorDetails.Message)
		}
	}
}

// Helper function for displaying warnings
func displayResolutionWarnings(resolutionResult *spookytypesvariables.VariableResolutionResult) {
	if len(resolutionResult.Warnings) > 0 {
		fmt.Printf("\n⚠️  Resolution warnings (%d):\n", len(resolutionResult.Warnings))
		for i := range resolutionResult.Warnings {
			warning := &resolutionResult.Warnings[i]
			fmt.Printf("  %d. %s: %s\n", i+1, warning.VariableName, warning.ErrorDetails.Message)
		}
	}
}

// Helper function for displaying resolved values
func displayResolvedValues(resolutionResult *spookytypesvariables.VariableResolutionResult) {
	if len(resolutionResult.Resolved) > 0 {
		fmt.Println("\nResolved values:")
		fmt.Println(strings.Repeat("-", 80))

		// Sort variables by name for consistent output
		var names []string
		for name := range resolutionResult.Resolved {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			displaySingleVariable(name, resolutionResult.Resolved[name], resolutionResult.Variables[name])
		}
	}
}

// Helper function for displaying a single variable
func displaySingleVariable(name string, value interface{}, variable *spookytypesvariables.Variable) {
	fmt.Printf("  %s = ", name)
	if variable.Sensitive {
		fmt.Print("[SENSITIVE]")
	} else {
		displayFormattedValue(value)
	}
	fmt.Printf(" (%s)", variable.Type)
	if variable.Description != "" {
		fmt.Printf(" - %s", variable.Description)
	}
	fmt.Println()
}

// Helper function for formatting values
func displayFormattedValue(value interface{}) {
	switch v := value.(type) {
	case string:
		fmt.Printf("%q", v)
	case nil:
		fmt.Print("null")
	default:
		fmt.Printf("%v", v)
	}
}

// Helper function for checking resolution errors
func checkResolutionErrors(resolutionResult *spookytypesvariables.VariableResolutionResult) error {
	if len(resolutionResult.Errors) > 0 {
		return fmt.Errorf("variable resolution completed with errors")
	}
	return nil
}

// handleVariablesResolve handles resolving variables using the VariablesIntegration interface
func handleVariablesResolve(_ *cobra.Command, projectPath string) error {
	ctx := context.Background()

	// Initialize dependencies
	if err := initializeVariablesDependenciesIfNeeded(); err != nil {
		return err
	}

	// Load variables from project
	variables, err := loadProjectVariables(ctx, projectPath)
	if err != nil {
		return err
	}

	// Create variable context
	variableContext := createVariableContext(projectPath)

	// Resolve variables
	resolutionResult, err := resolveProjectVariables(ctx, variables, variableContext)
	if err != nil {
		return err
	}

	// Display results
	displayResolutionSummary(projectPath, variables, resolutionResult)
	displayResolutionErrors(resolutionResult)
	displayResolutionWarnings(resolutionResult)
	displayResolvedValues(resolutionResult)

	// Check for errors
	return checkResolutionErrors(resolutionResult)
}

// handleVariablesEncrypt handles variables encryption
func handleVariablesEncrypt(projectPath string, dryRun bool) error {
	manager := GetIntegrationManager()
	variablesIntegration := manager.GetVariablesIntegration()
	if variablesIntegration == nil {
		return fmt.Errorf("variables integration not available")
	}

	return handleEncryptionOperation(projectPath, dryRun, "Variables", variablesIntegration.EncryptVariables)
}

// handleVariablesArmor handles variables armoring (encryption)
func handleVariablesArmor(projectPath string, dryRun bool) error {
	manager := GetIntegrationManager()
	variablesIntegration := manager.GetVariablesIntegration()
	if variablesIntegration == nil {
		return fmt.Errorf("variables integration not available")
	}

	return handleEncryptionOperation(projectPath, dryRun, "Variables", variablesIntegration.EncryptVariables)
}

func init() {
	// Add variables command to root
	RootCmd.AddCommand(variablesCmd)

	// Add subcommands to variables command
	variablesCmd.AddCommand(variablesListCmd)
	variablesCmd.AddCommand(variablesValidateCmd)
	variablesCmd.AddCommand(variablesResolveCmd)
	variablesCmd.AddCommand(variablesEncryptCmd)
	variablesCmd.AddCommand(variablesArmorCmd)

	// Add flags to variables resolve command
	variablesResolveCmd.Flags().String("environment", "", "Environment variables file (JSON)")
	variablesResolveCmd.Flags().String("facts", "", "Facts file (JSON)")
	variablesResolveCmd.Flags().String("machines", "", "Machine data file (JSON)")
	variablesResolveCmd.Flags().String("user-data", "", "User data file (JSON)")
	variablesResolveCmd.Flags().Bool("json", false, "Output results in JSON format")
	variablesResolveCmd.Flags().Bool("verbose", false, "Show detailed resolution information")

	// Add flags to variables encrypt command
	variablesEncryptCmd.Flags().Bool("dry-run", false, "Show what would be encrypted without making changes")

	// Add flags to variables armor command
	variablesArmorCmd.Flags().Bool("dry-run", false, "Show what would be encrypted without making changes")
}
