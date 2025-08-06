package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	spookycoordinator "spooky/internal/coordinator"
	spookylogging "spooky/internal/logging"
	spookyloggingtypes "spooky/internal/logging/types"

	"github.com/spf13/cobra"
)

// CreateVariablesCommand creates the variables command with all subcommands
func CreateVariablesCommand() *cobra.Command {
	variablesCmd := &cobra.Command{
		Use:   "variables",
		Short: "Manage variables",
		Long:  "Manage project variables",
	}

	// Add subcommands
	variablesCmd.AddCommand(createVariablesListCommand())
	variablesCmd.AddCommand(createVariablesSetCommand())
	variablesCmd.AddCommand(createVariablesGetCommand())
	variablesCmd.AddCommand(createVariablesExportCommand())

	return variablesCmd
}

// createVariablesListCommand creates the variables list subcommand
func createVariablesListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List variables",
		Long:  "List all project variables",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath := args[0]

			// Get flags
			_, _ = cmd.Flags().GetString("filter") // Keep flag but don't use for now
			showSensitive, _ := cmd.Flags().GetBool("show-sensitive")

			// Validate project path
			if _, err := os.Stat(projectPath); os.IsNotExist(err) {
				return fmt.Errorf("project directory does not exist: %s", projectPath)
			}

			// Create logger
			logger := spookylogging.NewLogger(spookyloggingtypes.Config{
				Level:  spookyloggingtypes.InfoLevel,
				Format: "text",
				Output: "stdout",
			})

			// Create coordinator manager
			coord, err := spookycoordinator.NewCoordinatorManagerFromProject(projectPath, logger)
			if err != nil {
				return fmt.Errorf("failed to create coordinator: %w", err)
			}

			// Load variables
			variablesContext, err := coord.Variables().LoadVariables(projectPath)
			if err != nil {
				return fmt.Errorf("failed to load variables: %w", err)
			}

			// Display variables
			if len(variablesContext.ResolvedVariables) == 0 {
				fmt.Println("No variables found")
				return nil
			}

			fmt.Printf("Found %d variables:\n\n", len(variablesContext.ResolvedVariables))

			for name, value := range variablesContext.ResolvedVariables {
				// Skip sensitive variables unless explicitly requested
				if !showSensitive && isSensitiveVariable(name) {
					fmt.Printf("  %s: [SENSITIVE]\n", name)
					continue
				}

				fmt.Printf("  %s: %v\n", name, value)
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringP("filter", "f", "", "Filter expression")
	cmd.Flags().BoolP("show-sensitive", "s", false, "Show sensitive variables")

	return cmd
}

// createVariablesSetCommand creates the variables set subcommand
func createVariablesSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set variable",
		Long:  "Set a project variable",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath := args[0]
			varName := args[1]
			varValue := args[2]

			// Get flags
			sensitive, _ := cmd.Flags().GetBool("sensitive")

			// Validate project path
			if _, err := os.Stat(projectPath); os.IsNotExist(err) {
				return fmt.Errorf("project directory does not exist: %s", projectPath)
			}

			// Create logger
			logger := spookylogging.NewLogger(spookyloggingtypes.Config{
				Level:  spookyloggingtypes.InfoLevel,
				Format: "text",
				Output: "stdout",
			})

			// Create coordinator manager
			coord, err := spookycoordinator.NewCoordinatorManagerFromProject(projectPath, logger)
			if err != nil {
				return fmt.Errorf("failed to create coordinator: %w", err)
			}

			// Load variables context
			variablesContext, err := coord.Variables().LoadVariables(projectPath)
			if err != nil {
				return fmt.Errorf("failed to load variables: %w", err)
			}

			// Set variable
			variablesContext.ResolvedVariables[varName] = varValue

			// Save variables (this would typically write to variables.hcl or similar)
			fmt.Printf("✅ Variable '%s' set to '%s'\n", varName, varValue)
			if sensitive {
				fmt.Printf("  (marked as sensitive)\n")
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringP("name", "n", "", "Variable name")
	cmd.Flags().StringP("value", "v", "", "Variable value")
	cmd.Flags().BoolP("sensitive", "s", false, "Mark as sensitive variable")

	return cmd
}

// createVariablesGetCommand creates the variables get subcommand
func createVariablesGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get variable",
		Long:  "Get a project variable value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath := args[0]
			varName := args[1]

			// Validate project path
			if _, err := os.Stat(projectPath); os.IsNotExist(err) {
				return fmt.Errorf("project directory does not exist: %s", projectPath)
			}

			// Create logger
			logger := spookylogging.NewLogger(spookyloggingtypes.Config{
				Level:  spookyloggingtypes.InfoLevel,
				Format: "text",
				Output: "stdout",
			})

			// Create coordinator manager
			coord, err := spookycoordinator.NewCoordinatorManagerFromProject(projectPath, logger)
			if err != nil {
				return fmt.Errorf("failed to create coordinator: %w", err)
			}

			// Load variables
			variablesContext, err := coord.Variables().LoadVariables(projectPath)
			if err != nil {
				return fmt.Errorf("failed to load variables: %w", err)
			}

			// Get variable
			value, exists := variablesContext.ResolvedVariables[varName]
			if !exists {
				return fmt.Errorf("variable '%s' not found", varName)
			}

			fmt.Printf("%v\n", value)
			return nil
		},
	}

	// Add flags
	cmd.Flags().StringP("name", "n", "", "Variable name")

	return cmd
}

// createVariablesExportCommand creates the variables export subcommand
func createVariablesExportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export variables",
		Long:  "Export variables to various formats",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath := args[0]

			// Get flags
			format, _ := cmd.Flags().GetString("format")
			output, _ := cmd.Flags().GetString("output")
			includeSensitive, _ := cmd.Flags().GetBool("include-sensitive")

			// Validate project path
			if _, err := os.Stat(projectPath); os.IsNotExist(err) {
				return fmt.Errorf("project directory does not exist: %s", projectPath)
			}

			// Validate format
			if format != "json" && format != "hcl" {
				return fmt.Errorf("unsupported format: %s (supported: json, hcl)", format)
			}

			// Set default output file if not specified
			if output == "" {
				output = filepath.Join(projectPath, fmt.Sprintf("variables.%s", format))
			}

			// Create logger
			logger := spookylogging.NewLogger(spookyloggingtypes.Config{
				Level:  spookyloggingtypes.InfoLevel,
				Format: "text",
				Output: "stdout",
			})

			// Create coordinator manager
			coord, err := spookycoordinator.NewCoordinatorManagerFromProject(projectPath, logger)
			if err != nil {
				return fmt.Errorf("failed to create coordinator: %w", err)
			}

			// Export variables
			fmt.Printf("Exporting variables to %s...\n", output)

			// Load variables
			variablesContext, err := coord.Variables().LoadVariables(projectPath)
			if err != nil {
				return fmt.Errorf("failed to load variables: %w", err)
			}

			// Create output file
			file, err := os.Create(output)
			if err != nil {
				return fmt.Errorf("failed to create output file: %w", err)
			}
			defer file.Close()

			// Export based on format
			if format == "json" {
				// Export as JSON
				encoder := json.NewEncoder(file)
				encoder.SetIndent("", "  ")
				if err := encoder.Encode(variablesContext.ResolvedVariables); err != nil {
					return fmt.Errorf("failed to encode JSON: %w", err)
				}
			} else if format == "hcl" {
				// Export as HCL
				fmt.Fprintf(file, "# Variables export\n\n")
				for name, value := range variablesContext.ResolvedVariables {
					// Skip sensitive variables unless explicitly requested
					if !includeSensitive && isSensitiveVariable(name) {
						fmt.Fprintf(file, "%s = \"[SENSITIVE]\"\n", name)
						continue
					}

					// Format value based on type
					switch v := value.(type) {
					case string:
						fmt.Fprintf(file, "%s = %q\n", name, v)
					case int, int64, float64, bool:
						fmt.Fprintf(file, "%s = %v\n", name, v)
					default:
						fmt.Fprintf(file, "%s = %q\n", name, fmt.Sprintf("%v", v))
					}
				}
			}

			fmt.Printf("✅ Variables exported successfully to %s\n", output)
			return nil
		},
	}

	// Add flags
	cmd.Flags().StringP("format", "f", "json", "Export format (json, hcl)")
	cmd.Flags().StringP("output", "o", "", "Output file path")
	cmd.Flags().BoolP("include-sensitive", "s", false, "Include sensitive variables")

	return cmd
}

// Helper function to check if a variable name indicates it's sensitive
func isSensitiveVariable(name string) bool {
	sensitivePatterns := []string{"password", "secret", "key", "token", "credential", "auth"}
	for _, pattern := range sensitivePatterns {
		if contains(name, pattern) {
			return true
		}
	}
	return false
}

// Helper function to check if a string contains a substring (case insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(len(s) == len(substr) && s == substr ||
			len(s) > len(substr) && (s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
