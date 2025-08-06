package commands

import (
	"fmt"
	"os"

	spookyactionstypes "spooky/internal/actions/types"
	spookycoordinator "spooky/internal/coordinator"
	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookyloggingtypes "spooky/internal/logging/types"

	"github.com/spf13/cobra"
)

// CreateActionsCommand creates the actions command with all subcommands
func CreateActionsCommand() *cobra.Command {
	actionsCmd := &cobra.Command{
		Use:   "actions",
		Short: "Manage and execute actions",
		Long:  "Manage and execute actions on remote machines",
	}

	// Add subcommands
	actionsCmd.AddCommand(createActionsListCommand())
	actionsCmd.AddCommand(createActionsRunCommand())
	actionsCmd.AddCommand(createActionsValidateCommand())

	return actionsCmd
}

// createActionsListCommand creates the actions list subcommand
func createActionsListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available actions",
		Long:  "List all available actions in the project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath := args[0]

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

			// List actions from project
			actions, err := coord.Actions().ListActionsFromProject(projectPath)
			if err != nil {
				return fmt.Errorf("failed to list actions: %w", err)
			}

			// Display actions
			if len(actions) == 0 {
				fmt.Println("No actions found in project")
				return nil
			}

			fmt.Printf("Found %d actions in project:\n\n", len(actions))
			for _, action := range actions {
				fmt.Printf("  %s\n", action.Name)
				if action.Description != "" {
					fmt.Printf("    Description: %s\n", action.Description)
				}
				if action.Type != "" {
					fmt.Printf("    Type: %s\n", action.Type)
				}
				fmt.Println()
			}

			return nil
		},
	}
}

// createActionsRunCommand creates the actions run subcommand
func createActionsRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run actions",
		Long:  "Execute actions on remote machines",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath := args[0]

			// Get flags
			actionName, _ := cmd.Flags().GetString("action")
			machines, _ := cmd.Flags().GetStringSlice("machines")
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			decrypt, _ := cmd.Flags().GetBool("decrypt")
			_, _ = cmd.Flags().GetInt("parallel") // Keep flag but don't use for now

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

			// Load actions from project first
			actions, err := coord.Actions().ListActionsFromProject(projectPath)
			if err != nil {
				return fmt.Errorf("failed to load actions from project: %w", err)
			}

			// Execute action
			if actionName != "" {
				// Find the specific action
				var action *spookyactionstypes.Action
				for _, a := range actions {
					if a.Name == actionName {
						action = a
						break
					}
				}
				if action == nil {
					return fmt.Errorf("action '%s' not found in project", actionName)
				}

				// Create execution context
				execContext := &spookyinterfaces.ActionExecutionContext{
					BaseContext: spookyinterfaces.BaseContext{
						ProjectPath: projectPath,
					},
					MachineNames: machines,
					Action:       action,
					Decrypt:      decrypt,
				}

				if dryRun {
					fmt.Printf("DRY RUN: Would execute action '%s' on machines: %v\n", actionName, machines)
					return nil
				}

				fmt.Printf("Executing action '%s' on machines: %v\n", actionName, machines)
				return coord.Actions().ExecuteAction(action, execContext)
			} else {
				// Execute all actions
				if len(actions) == 0 {
					fmt.Println("No actions found to execute")
					return nil
				}

				if dryRun {
					fmt.Printf("DRY RUN: Would execute %d actions on machines: %v\n", len(actions), machines)
					return nil
				}

				fmt.Printf("Executing %d actions on machines: %v\n", len(actions), machines)
				for _, action := range actions {
					execContext := &spookyinterfaces.ActionExecutionContext{
						BaseContext: spookyinterfaces.BaseContext{
							ProjectPath: projectPath,
						},
						MachineNames: machines,
						Action:       action,
						Decrypt:      decrypt,
					}

					if err := coord.Actions().ExecuteAction(action, execContext); err != nil {
						fmt.Printf("Failed to execute action '%s': %v\n", action.Name, err)
					}
				}
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringP("action", "a", "", "Action name to execute")
	cmd.Flags().StringSliceP("machines", "m", []string{}, "Target machines")
	cmd.Flags().BoolP("dry-run", "d", false, "Dry run mode")
	cmd.Flags().BoolP("decrypt", "x", false, "Enable decryption of variables and facts for use in templates")
	cmd.Flags().IntP("parallel", "p", 1, "Number of parallel executions")

	return cmd
}

// createActionsValidateCommand creates the actions validate subcommand
func createActionsValidateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate actions",
		Long:  "Validate actions in the project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath := args[0]

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

			// Load actions from project
			actions, err := coord.Actions().ListActionsFromProject(projectPath)
			if err != nil {
				return fmt.Errorf("failed to load actions from project: %w", err)
			}

			// Validate each action
			if len(actions) == 0 {
				fmt.Println("No actions found in project")
				return nil
			}

			fmt.Printf("Validating %d actions in project:\n\n", len(actions))
			validCount := 0
			for _, action := range actions {
				// Create execution context for validation
				execContext := &spookyinterfaces.ActionExecutionContext{
					BaseContext: spookyinterfaces.BaseContext{
						ProjectPath: projectPath,
					},
					Action: action,
				}

				if err := coord.Actions().ValidateAction(action, execContext); err != nil {
					fmt.Printf("❌ %s: %v\n", action.Name, err)
				} else {
					fmt.Printf("✅ %s: Valid\n", action.Name)
					validCount++
				}
			}

			fmt.Printf("\nValidation complete: %d/%d actions valid\n", validCount, len(actions))
			return nil
		},
	}
}
