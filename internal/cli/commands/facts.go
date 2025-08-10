package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CreateFactsCommand creates the main facts command
func CreateFactsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "facts",
		Short: "Manage facts",
		Long:  "Manage and collect facts from remote machines",
	}

	// Add subcommands
	cmd.AddCommand(createFactsGatherCommand())
	cmd.AddCommand(createFactsListCommand())
	cmd.AddCommand(createFactsValidateCommand())

	return cmd
}

// createFactsGatherCommand creates the facts gather subcommand
func createFactsGatherCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gather",
		Short: "Gather facts",
		Long:  "Collect facts from remote machines",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement properly with correct types
			return fmt.Errorf("not implemented - interface mismatches need to be resolved")
		},
	}
	return cmd
}

// createFactsListCommand creates the facts list subcommand
func createFactsListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List facts",
		Long:  "List collected facts",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement properly with correct types
			return fmt.Errorf("not implemented - interface mismatches need to be resolved")
		},
	}
	return cmd
}

// createFactsValidateCommand creates the facts validate subcommand
func createFactsValidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate facts",
		Long:  "Validate collected facts",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement properly with correct types
			return fmt.Errorf("not implemented - interface mismatches need to be resolved")
		},
	}
	return cmd
}
