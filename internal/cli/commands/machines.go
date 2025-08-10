package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CreateMachinesCommand creates the main machines command
func CreateMachinesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "machines",
		Short: "Manage machines",
		Long:  "Manage and connect to remote machines",
	}

	// Add subcommands
	cmd.AddCommand(createMachinesListCommand())
	cmd.AddCommand(createMachinesPingCommand())

	return cmd
}

// createMachinesListCommand creates the machines list subcommand
func createMachinesListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List machines",
		Long:  "List all machines in a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement properly with correct types
			return fmt.Errorf("not implemented - interface mismatches need to be resolved")
		},
	}
	return cmd
}

// createMachinesPingCommand creates the machines ping subcommand
func createMachinesPingCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ping",
		Short: "Ping machines",
		Long:  "Ping machines to check connectivity",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement properly with correct types
			return fmt.Errorf("not implemented - interface mismatches need to be resolved")
		},
	}
	return cmd
}
