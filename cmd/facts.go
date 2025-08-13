// Package cmd provides command implementations for spooky CLI.
package cmd

import (
	"context"
	"fmt"
	"strings"

	spookyfacts "spooky/internal/facts"
	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"

	"github.com/spf13/cobra"
)

// Global instances for facts dependency injection
var (
	factsManager spookyinterfaces.FactsIntegration
	factsLogger  spookytypeslogging.Logger
)

// InitializeFactsDependencies initializes facts-related dependencies
func InitializeFactsDependencies() error {
	// Create log manager for facts component
	logManager := spookylogging.NewLogManager()
	factsLogger = logManager.GetLogger("facts")

	// Initialize facts components
	collector := spookyfacts.NewSystemFactCollector()
	manager := spookyfacts.NewManager(collector, nil, factsLogger)

	// Create facts integration
	factsManager = spookyfacts.NewIntegration(manager)

	return nil
}

// factsCmd represents the facts command
var factsCmd = &cobra.Command{
	Use:   "facts",
	Short: "Manage machine facts",
	Long: `Manage machine facts including collection, storage, and validation.

Facts are system information collected from machines including OS details,
hardware information, network configuration, and custom data. Facts are
stored in memory and can be used for decision making in actions.`,
}

// factsExportCmd represents the facts export command
var factsExportCmd = &cobra.Command{
	Use:   "export [project-path]",
	Short: "Export facts to file",
	Long: `Export facts to a file in JSON or HCL format.

This command automatically gathers facts from all machines in the project
inventory and exports them to the specified format for backup, analysis,
or transfer to other systems.

Examples:
  spooky facts export ./my-project --output facts.hcl
  spooky facts export ./my-project --format json --output facts.json
  spooky facts export ./my-project --machine web-server --output web-server-facts.hcl`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleFactsExport(cmd, args[0])
	},
}

// handleFactsExport handles exporting facts using the FactsManager interface
func handleFactsExport(cmd *cobra.Command, projectPath string) error {
	ctx := context.Background()

	// Initialize dependencies if not already done
	if factsManager == nil {
		if err := InitializeFactsDependencies(); err != nil {
			return fmt.Errorf("failed to initialize facts dependencies: %w", err)
		}
	}

	// Get flags
	format, _ := cmd.Flags().GetString("format")
	outputPath, _ := cmd.Flags().GetString("output")
	machineFilter, _ := cmd.Flags().GetString("machine")
	parallel, _ := cmd.Flags().GetInt("parallel")
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Load machines from project
	machines, err := loadMachinesFromProject(projectPath)
	if err != nil {
		return fmt.Errorf("failed to load machines: %w", err)
	}

	// Filter machines if specified
	if machineFilter != "" {
		machines = filterMachines(machines, machineFilter)
		if len(machines) == 0 {
			return fmt.Errorf("no machines found matching filter: %s", machineFilter)
		}
	}

	if verbose {
		factsLogger.Info("Starting fact collection for export", map[string]interface{}{
			"project_path":  projectPath,
			"machine_count": len(machines),
			"parallel":      parallel,
		})
	}

	// Collect facts from each machine
	successCount := 0
	errorCount := 0

	for _, machine := range machines {
		if verbose {
			factsLogger.Info("Collecting facts for export", map[string]interface{}{
				"machine": machine.Hostname,
			})
		}

		// Collect facts using the integration
		facts, err := factsManager.CollectFacts(ctx, machine.Hostname)
		if err != nil {
			errorCount++
			factsLogger.Error("Failed to collect facts for export", err, map[string]interface{}{
				"machine": machine.Hostname,
			})
			continue
		}

		// Store facts
		if err := factsManager.StoreFacts(ctx, facts, nil); err != nil {
			errorCount++
			factsLogger.Error("Failed to store facts for export", err, map[string]interface{}{
				"machine": machine.Hostname,
			})
			continue
		}

		successCount++
		if verbose {
			factsLogger.Info("Successfully collected facts for export", map[string]interface{}{
				"machine": machine.Hostname,
			})
		}
	}

	if verbose {
		fmt.Printf("Fact collection completed:\n")
		fmt.Printf("  Success: %d machines\n", successCount)
		fmt.Printf("  Errors: %d machines\n", errorCount)
		fmt.Printf("  Total: %d machines\n", len(machines))
	}

	// Get the underlying fact manager
	managerInterface := factsManager.GetManager()
	manager, ok := managerInterface.(spookyfacts.FactManager)
	if !ok {
		return fmt.Errorf("failed to get underlying fact manager")
	}

	// Get all machine IDs with stored facts
	machineIDs, err := manager.ListFacts(ctx)
	if err != nil {
		return fmt.Errorf("failed to list facts: %w", err)
	}

	if len(machineIDs) == 0 {
		return fmt.Errorf("no facts found to export")
	}

	// Export facts using the manager's export functionality
	if err := manager.ExportFacts(ctx, machineIDs, format, outputPath); err != nil {
		return fmt.Errorf("failed to export facts: %w", err)
	}

	fmt.Printf("Successfully exported facts to: %s\n", outputPath)
	return nil
}

// Helper functions

func loadMachinesFromProject(projectPath string) ([]spookytypes.Machine, error) {
	// This would load machines from the project's machines.hcl file
	// For now, return a placeholder implementation
	return []spookytypes.Machine{}, nil
}

func filterMachines(machines []spookytypes.Machine, filter string) []spookytypes.Machine {
	var filtered []spookytypes.Machine
	for _, machine := range machines {
		if strings.Contains(machine.Hostname, filter) {
			filtered = append(filtered, machine)
		}
	}
	return filtered
}

func init() {

	// Add flags to facts export command
	factsExportCmd.Flags().String("format", "hcl", "Export format (hcl, json)")
	factsExportCmd.Flags().String("output", "", "Output file path (required)")
	factsExportCmd.Flags().String("machine", "", "Filter to specific machine")
	factsExportCmd.Flags().Int("parallel", 1, "Number of parallel workers")
	factsExportCmd.Flags().Bool("verbose", false, "Verbose output")
	factsExportCmd.MarkFlagRequired("output")

	// Add commands to facts command
	factsCmd.AddCommand(factsExportCmd)

	// Add facts command to root
	RootCmd.AddCommand(factsCmd)
}
