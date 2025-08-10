package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	spookycoordinator "spooky/internal/coordinator"
	spookylogging "spooky/internal/logging"
	spookyloggingtypes "spooky/internal/types/logging"

	"github.com/spf13/cobra"
)

// CreateMachinesCommand creates the machines command with all subcommands
func CreateMachinesCommand() *cobra.Command {
	machinesCmd := &cobra.Command{
		Use:   "machines",
		Short: "Manage machines",
		Long:  "Manage machine inventory and connectivity",
	}

	// Add subcommands
	machinesCmd.AddCommand(createMachinesListCommand())
	machinesCmd.AddCommand(createMachinesPingCommand())
	machinesCmd.AddCommand(createMachinesExportCommand())

	return machinesCmd
}

// createMachinesListCommand creates the machines list subcommand
func createMachinesListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List machines",
		Long:  "List all machines in the inventory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath := args[0]

			// Get flags
			_, _ = cmd.Flags().GetStringSlice("tags") // Keep flag but don't use for now
			_, _ = cmd.Flags().GetString("filter")    // Keep flag but don't use for now

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

			// Load machines context
			machinesContext, err := coord.Machines().LoadMachines(projectPath)
			if err != nil {
				return fmt.Errorf("failed to load machines: %w", err)
			}

			// List machines
			machines, err := coord.Machines().ListMachines(machinesContext)
			if err != nil {
				return fmt.Errorf("failed to list machines: %w", err)
			}

			// Display machines
			if len(machines) == 0 {
				fmt.Println("No machines found in inventory")
				return nil
			}

			fmt.Printf("Found %d machines in inventory:\n\n", len(machines))

			for _, machine := range machines {
				fmt.Printf("  %s\n", machine.Name)
				fmt.Printf("    Host: %s:%d\n", machine.Host, machine.Port)
				fmt.Printf("    Username: %s\n", machine.Username)
				if len(machine.Tags) > 0 {
					fmt.Printf("    Tags: %v\n", machine.Tags)
				}
				fmt.Println()
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringSliceP("tags", "t", []string{}, "Filter by tags")
	cmd.Flags().StringP("filter", "f", "", "Filter expression")

	return cmd
}

// createMachinesPingCommand creates the machines ping subcommand
func createMachinesPingCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ping",
		Short: "Ping machines",
		Long:  "Test connectivity to machines",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath := args[0]

			// Get flags
			machines, _ := cmd.Flags().GetStringSlice("machines")
			_, _ = cmd.Flags().GetStringSlice("tags") // Keep flag but don't use for now
			_, _ = cmd.Flags().GetInt("timeout")      // Keep flag but don't use for now

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

			// Load machines first to build indexes
			machinesContext, err := coord.Machines().LoadMachines(projectPath)
			if err != nil {
				return fmt.Errorf("failed to load machines: %w", err)
			}

			// Get machines to ping
			var machinesToPing []string
			if len(machines) > 0 {
				machinesToPing = machines
			} else {
				// Get all machines if none specified
				allMachines, err := coord.Machines().ListMachines(machinesContext)
				if err != nil {
					return fmt.Errorf("failed to list machines: %w", err)
				}
				for _, machine := range allMachines {
					machinesToPing = append(machinesToPing, machine.Name)
				}
			}

			if len(machinesToPing) == 0 {
				fmt.Println("No machines found to ping")
				return nil
			}

			// Ping machines
			fmt.Printf("Pinging %d machines...\n", len(machinesToPing))

			successCount := 0
			failureCount := 0

			for _, machineName := range machinesToPing {
				fmt.Printf("Pinging %s... ", machineName)

				if err := coord.Machines().PingMachine(machineName); err != nil {
					fmt.Printf("❌ FAILED: %v\n", err)
					failureCount++
				} else {
					fmt.Printf("✅ SUCCESS\n")
					successCount++
				}
			}

			fmt.Printf("\nPing results: %d successful, %d failed\n", successCount, failureCount)

			if failureCount > 0 {
				return fmt.Errorf("ping failed for %d machines", failureCount)
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringSliceP("machines", "m", []string{}, "Target machines")
	cmd.Flags().StringSliceP("tags", "t", []string{}, "Filter by tags")
	cmd.Flags().IntP("timeout", "T", 30, "Timeout in seconds")

	return cmd
}

// createMachinesExportCommand creates the machines export subcommand
func createMachinesExportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export machines",
		Long:  "Export machine inventory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath := args[0]

			// Get flags
			format, _ := cmd.Flags().GetString("format")
			output, _ := cmd.Flags().GetString("output")
			_, _ = cmd.Flags().GetStringSlice("tags") // Keep flag but don't use for now

			// Validate project path
			if _, err := os.Stat(projectPath); os.IsNotExist(err) {
				return fmt.Errorf("project directory does not exist: %s", projectPath)
			}

			// Validate format
			if format != "json" {
				return fmt.Errorf("unsupported format: %s (only json supported)", format)
			}

			// Set default output file if not specified
			if output == "" {
				output = filepath.Join(projectPath, "machines.json")
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

			// Export machines
			fmt.Printf("Exporting machines to %s...\n", output)

			// Load machines context
			machinesContext, err := coord.Machines().LoadMachines(projectPath)
			if err != nil {
				return fmt.Errorf("failed to load machines: %w", err)
			}

			// Get machines
			machines, err := coord.Machines().ListMachines(machinesContext)
			if err != nil {
				return fmt.Errorf("failed to list machines: %w", err)
			}

			// Create output file
			file, err := os.Create(output)
			if err != nil {
				return fmt.Errorf("failed to create output file: %w", err)
			}
			defer file.Close()

			// Export as JSON
			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(machines); err != nil {
				return fmt.Errorf("failed to encode JSON: %w", err)
			}

			fmt.Printf("✅ Machines exported successfully to %s\n", output)
			return nil
		},
	}

	// Add flags
	cmd.Flags().StringP("format", "f", "json", "Export format (json)")
	cmd.Flags().StringP("output", "o", "", "Output file path")
	cmd.Flags().StringSliceP("tags", "t", []string{}, "Filter by tags")

	return cmd
}
