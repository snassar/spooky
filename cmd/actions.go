package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	spookyinterfaces "spooky/internal/interfaces"

	"github.com/spf13/cobra"
)

var (
	// actionsCmd represents the actions command
	actionsCmd = &cobra.Command{
		Use:   "actions",
		Short: "Manage and run actions",
		Long: `Manage and run actions on machines.

Actions are run in dependency order. If no action names are provided,
all actions in the project will be run.`,
	}

	// actionsListCmd represents the actions list command
	actionsListCmd = &cobra.Command{
		Use:   "list",
		Short: "List available actions",
		Long:  `List all available actions in the project.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("project directory is required")
			}

			projectPath := args[0]
			if err := validateProjectPath(projectPath); err != nil {
				return err
			}

			// Get actions integration
			actionsIntegration := getIntegrationManager().GetActionsIntegration()
			if actionsIntegration == nil {
				return fmt.Errorf("actions integration not available")
			}

			// Load actions from project
			ctx := context.Background()
			actions, err := actionsIntegration.LoadActions(ctx, projectPath)
			if err != nil {
				return fmt.Errorf("failed to load actions: %w", err)
			}

			if len(actions) == 0 {
				fmt.Println("No actions found in project")
				return nil
			}

			fmt.Printf("Available actions (%d found):\n", len(actions))
			for i, action := range actions {
				fmt.Printf("%d. %s", i+1, action.Name)
				if action.Description != "" {
					fmt.Printf(" - %s", action.Description)
				}
				fmt.Println()
			}

			return nil
		},
	}

	// actionsRunCmd represents the actions run command
	actionsRunCmd = &cobra.Command{
		Use:   "run",
		Short: "Run actions on machines",
		Long: `Run actions on target machines.

Actions are run in dependency order. If no action names are provided,
all actions in the project will be run.

Use --plan to see what would be run without actually running.
Use --dry-run to simulate running without making changes.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("project directory is required")
			}

			projectPath := args[0]
			if err := validateProjectPath(projectPath); err != nil {
				return err
			}

			// Get actions integration
			actionsIntegration := getIntegrationManager().GetActionsIntegration()
			if actionsIntegration == nil {
				return fmt.Errorf("actions integration not available")
			}

			ctx := context.Background()

			// Load actions from project
			actions, err := actionsIntegration.LoadActions(ctx, projectPath)
			if err != nil {
				return fmt.Errorf("failed to load actions: %w", err)
			}

			if len(actions) == 0 {
				return fmt.Errorf("no actions found in project")
			}

			// Check running mode
			if isFlagSet(cmd, "plan") {
				fmt.Printf("Action running plan (%d actions):\n", len(actions))
				for i, action := range actions {
					fmt.Printf("%d. %s", i+1, action.Name)
					if action.Description != "" {
						fmt.Printf(" - %s", action.Description)
					}
					fmt.Println()
				}
				return nil
			}

			if isFlagSet(cmd, "dry-run") {
				fmt.Println("Simulating running...")
				for _, action := range actions {
					fmt.Printf("[%s] Would run: %s\n", action.Name, action.Command)
				}
				return nil
			}

			// Normal running mode
			fmt.Printf("Running %d actions...\n", len(actions))

			// Load machines from project for action execution
			machinesIntegration := getIntegrationManager().GetMachinesIntegration()
			if machinesIntegration == nil {
				return fmt.Errorf("machines integration not available")
			}

			machines, err := machinesIntegration.LoadMachines(ctx, projectPath)
			if err != nil {
				return fmt.Errorf("failed to load machines: %w", err)
			}

			if len(machines) == 0 {
				return fmt.Errorf("no machines found in project for action execution")
			}

			fmt.Printf("Found %d machines for action execution\n", len(machines))

			// Run actions using the actions integration
			results, err := actionsIntegration.RunActions(ctx, actions, machines)
			if err != nil {
				return fmt.Errorf("failed to run actions: %w", err)
			}

			// Display results
			fmt.Printf("\n📊 Action Execution Results:\n")
			fmt.Printf("Total actions executed: %d\n", len(results))

			successCount := 0
			failureCount := 0

			for _, result := range results {
				if result.Status == "success" {
					successCount++
					fmt.Printf("✅ %s on %s: Success\n", result.ActionName, result.MachineName)
				} else {
					failureCount++
					fmt.Printf("❌ %s on %s: Failed - %s\n", result.ActionName, result.MachineName, result.Error)
				}
			}

			fmt.Printf("\n📈 Summary: %d successful, %d failed\n", successCount, failureCount)

			if failureCount > 0 {
				return fmt.Errorf("action execution completed with %d failures", failureCount)
			}

			fmt.Println("🎉 All actions completed successfully!")
			return nil
		},
	}

	// actionsValidateCmd represents the actions validate command
	actionsValidateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate action configurations",
		Long:  `Validate action configurations for syntax and dependencies.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("project directory is required")
			}

			projectPath := args[0]
			if err := validateProjectPath(projectPath); err != nil {
				return err
			}

			// Get actions integration
			actionsIntegration := getIntegrationManager().GetActionsIntegration()
			if actionsIntegration == nil {
				return fmt.Errorf("actions integration not available")
			}

			ctx := context.Background()

			// Load actions from project
			actions, err := actionsIntegration.LoadActions(ctx, projectPath)
			if err != nil {
				return fmt.Errorf("failed to load actions: %w", err)
			}

			if len(actions) == 0 {
				fmt.Println("No actions found to validate")
				return nil
			}

			// Validate actions
			result, err := actionsIntegration.ValidateActions(ctx, actions)
			if err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			if result.Valid {
				fmt.Printf("✅ All %d actions are valid!\n", len(actions))
			} else {
				fmt.Printf("❌ Validation failed for %d actions:\n", len(result.Errors))
				for _, err := range result.Errors {
					fmt.Printf("  - %s\n", err.Message)
				}
				return fmt.Errorf("action validation failed")
			}

			return nil
		},
	}
)

// Helper function to check if a flag is set
func isFlagSet(cmd *cobra.Command, flagName string) bool {
	return cmd.Flags().Lookup(flagName) != nil && cmd.Flags().Lookup(flagName).Changed
}

// Helper function to validate project path
func validateProjectPath(projectPath string) error {
	if projectPath == "" {
		return fmt.Errorf("project path cannot be empty")
	}

	// Check if project path exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return fmt.Errorf("project path does not exist: %s", projectPath)
	}

	// Check if it's a directory
	if info, err := os.Stat(projectPath); err == nil && !info.IsDir() {
		return fmt.Errorf("project path must be a directory: %s", projectPath)
	}

	// Check for project.hcl file
	projectFile := filepath.Join(projectPath, "project.hcl")
	if _, err := os.Stat(projectFile); os.IsNotExist(err) {
		return fmt.Errorf("project.hcl not found in: %s", projectPath)
	}

	return nil
}

// Helper function to get integration manager
func getIntegrationManager() spookyinterfaces.IntegrationManager {
	// Return the global integration manager from cmd package
	return integrationManager
}

func init() {
	RootCmd.AddCommand(actionsCmd)

	// Add subcommands
	actionsCmd.AddCommand(actionsListCmd)
	actionsCmd.AddCommand(actionsRunCmd)
	actionsCmd.AddCommand(actionsValidateCmd)

	// Add flags to run command
	actionsRunCmd.Flags().BoolP("plan", "p", false, "Show running plan without running")
	actionsRunCmd.Flags().BoolP("dry-run", "d", false, "Simulate running without making changes")
	actionsRunCmd.Flags().StringSliceP("machine", "m", nil, "Target specific machines")
	actionsRunCmd.Flags().StringSliceP("tags", "t", nil, "Target machines with specific tags")
	actionsRunCmd.Flags().StringP("filter", "f", "", "Complex filter expression")
	actionsRunCmd.Flags().IntP("parallel", "j", 1, "Number of parallel workers (minimum 2)")
}
