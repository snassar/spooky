package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CreateActionsCommand creates the main actions command
func CreateActionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "actions",
		Short: "Manage actions",
		Long:  "Manage and execute actions on remote machines",
	}

	// Add subcommands
	cmd.AddCommand(createActionsListCommand())
	cmd.AddCommand(createActionsRunCommand())
	cmd.AddCommand(createActionsValidateCommand())

	return cmd
}

// createActionsListCommand creates the actions list subcommand
func createActionsListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List actions",
		Long:  "List all actions in a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement properly with correct types
			return fmt.Errorf("not implemented - interface mismatches need to be resolved")
		},
	}
	return cmd
}

// createActionsRunCommand creates the actions run subcommand
func createActionsRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run actions",
		Long:  "Execute actions on remote machines",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement properly with correct types
			return fmt.Errorf("not implemented - interface mismatches need to be resolved")
		},
	}
	return cmd
}

// createActionsValidateCommand creates the actions validate subcommand
func createActionsValidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate actions",
		Long:  "Validate action configurations",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement properly with correct types
			return fmt.Errorf("not implemented - interface mismatches need to be resolved")
		},
	}
	return cmd
}
