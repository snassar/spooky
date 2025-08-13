// Package cmd provides command implementations for spooky CLI.
package cmd

import (
	"context"
	"encoding/json"
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

// factsGatherCmd represents the facts gather command
var factsGatherCmd = &cobra.Command{
	Use:   "gather [project-path]",
	Short: "Gather facts from machines",
	Long: `Gather facts from all machines in the project inventory.

This command connects to each machine in the project's machine inventory
and collects system facts including OS information, hardware details,
network configuration, and other system data. Facts are stored in memory
for the duration of the session.

Examples:
  spooky facts gather ./my-project
  spooky facts gather ./my-project --parallel 4
  spooky facts gather ./my-project --machine web-server`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleFactsGather(cmd, args[0])
	},
}

// factsListCmd represents the facts list command
var factsListCmd = &cobra.Command{
	Use:   "list [project-path]",
	Short: "List stored facts",
	Long: `List all machines with stored facts in the project.

This command displays information about all machines that have facts
stored in memory, including collection timestamps and fact summaries.

Examples:
  spooky facts list ./my-project
  spooky facts list ./my-project --machine web-server`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleFactsList(cmd, args[0])
	},
}

// factsValidateCmd represents the facts validate command
var factsValidateCmd = &cobra.Command{
	Use:   "validate [project-path]",
	Short: "Validate stored facts",
	Long: `Validate facts stored in memory.

This command validates that all stored facts are properly formatted and
contain required data. It checks fact structure, timestamps, and data
integrity.

Examples:
  spooky facts validate ./my-project
  spooky facts validate ./my-project --compare`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleFactsValidate(cmd, args[0])
	},
}

// factsExportCmd represents the facts export command
var factsExportCmd = &cobra.Command{
	Use:   "export [project-path]",
	Short: "Export facts to file",
	Long: `Export facts from memory to a file.

This command exports facts to JSON or HCL format for backup, analysis,
or transfer to other systems.

Examples:
  spooky facts export ./my-project --format json --output facts.json
  spooky facts export ./my-project --format hcl --output facts.hcl`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleFactsExport(cmd, args[0])
	},
}

// handleFactsGather handles gathering facts using the FactsIntegration interface
func handleFactsGather(cmd *cobra.Command, projectPath string) error {
	ctx := context.Background()

	// Initialize dependencies if not already done
	if factsManager == nil {
		if err := InitializeFactsDependencies(); err != nil {
			return fmt.Errorf("failed to initialize facts dependencies: %w", err)
		}
	}

	// Get flags
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
		factsLogger.Info("Starting fact collection", map[string]interface{}{
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
			factsLogger.Info("Collecting facts", map[string]interface{}{
				"machine": machine.Hostname,
			})
		}

		// Collect facts using the integration
		facts, err := factsManager.CollectFacts(ctx, machine.Hostname)
		if err != nil {
			errorCount++
			factsLogger.Error("Failed to collect facts", err, map[string]interface{}{
				"machine": machine.Hostname,
			})
			continue
		}

		// Store facts
		if err := factsManager.StoreFacts(ctx, facts, nil); err != nil {
			errorCount++
			factsLogger.Error("Failed to store facts", err, map[string]interface{}{
				"machine": machine.Hostname,
			})
			continue
		}

		successCount++
		if verbose {
			factsLogger.Info("Successfully collected facts", map[string]interface{}{
				"machine": machine.Hostname,
			})
		}
	}

	// Print summary
	fmt.Printf("Fact collection completed:\n")
	fmt.Printf("  Success: %d machines\n", successCount)
	fmt.Printf("  Errors: %d machines\n", errorCount)
	fmt.Printf("  Total: %d machines\n", len(machines))

	if errorCount > 0 {
		return fmt.Errorf("fact collection completed with %d errors", errorCount)
	}

	return nil
}

// handleFactsList handles listing facts using the FactsIntegration interface
func handleFactsList(cmd *cobra.Command, projectPath string) error {
	ctx := context.Background()

	// Initialize dependencies if not already done
	if factsManager == nil {
		if err := InitializeFactsDependencies(); err != nil {
			return fmt.Errorf("failed to initialize facts dependencies: %w", err)
		}
	}

	// Get flags
	format, _ := cmd.Flags().GetString("format")

	// Load facts from storage
	facts, err := factsManager.LoadFacts(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to load facts: %w", err)
	}

	if facts == nil {
		fmt.Println("No facts found in storage")
		return nil
	}

	// Display facts based on format
	switch format {
	case "json":
		data, err := json.MarshalIndent(facts, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal facts to JSON: %w", err)
		}
		fmt.Println(string(data))

	default:
		fmt.Printf("Facts loaded: %+v\n", facts)
	}

	return nil
}

// handleFactsValidate handles validating facts using the FactsIntegration interface
func handleFactsValidate(cmd *cobra.Command, projectPath string) error {
	ctx := context.Background()

	// Initialize dependencies if not already done
	if factsManager == nil {
		if err := InitializeFactsDependencies(); err != nil {
			return fmt.Errorf("failed to initialize facts dependencies: %w", err)
		}
	}

	// Get flags
	compare, _ := cmd.Flags().GetBool("compare")

	// Load facts from storage
	facts, err := factsManager.LoadFacts(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to load facts: %w", err)
	}

	if facts == nil {
		fmt.Println("No facts found to validate")
		return nil
	}

	// Validate facts
	validationResult, err := factsManager.ValidateFacts(ctx, facts)
	if err != nil {
		return fmt.Errorf("failed to validate facts: %w", err)
	}

	// Display validation results
	fmt.Printf("Facts validation results:\n")
	fmt.Printf("  Valid: %t\n", validationResult.Valid)
	fmt.Printf("  Errors: %d\n", len(validationResult.Errors))
	fmt.Printf("  Warnings: %d\n", len(validationResult.Warnings))

	if len(validationResult.Errors) > 0 {
		fmt.Printf("\nErrors:\n")
		for _, err := range validationResult.Errors {
			fmt.Printf("  - %s\n", err.Message)
		}
	}

	if len(validationResult.Warnings) > 0 {
		fmt.Printf("\nWarnings:\n")
		for _, warning := range validationResult.Warnings {
			fmt.Printf("  - %s\n", warning.Message)
		}
	}

	if !validationResult.Valid {
		return fmt.Errorf("facts validation failed with %d errors", len(validationResult.Errors))
	}

	// Compare with fresh facts if requested
	if compare {
		fmt.Printf("\nComparing with fresh facts...\n")
		if err := compareWithFreshFacts(ctx, projectPath, facts); err != nil {
			return fmt.Errorf("failed to compare facts: %w", err)
		}
	}

	return nil
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

func compareWithFreshFacts(ctx context.Context, projectPath string, storedFacts interface{}) error {
	// This would compare stored facts with fresh facts from machines
	// For now, just indicate that comparison is not implemented
	fmt.Println("Fact comparison not yet implemented")
	return nil
}

func init() {
	// Add flags to facts gather command
	factsGatherCmd.Flags().String("machine", "", "Filter to specific machine")
	factsGatherCmd.Flags().Int("parallel", 1, "Number of parallel workers")
	factsGatherCmd.Flags().Bool("verbose", false, "Verbose output")

	// Add flags to facts list command
	factsListCmd.Flags().String("format", "summary", "Output format (summary, table, json)")

	// Add flags to facts validate command
	factsValidateCmd.Flags().Bool("compare", false, "Compare with fresh facts")

	// Add flags to facts export command
	factsExportCmd.Flags().String("format", "json", "Export format (json, hcl)")
	factsExportCmd.Flags().String("output", "", "Output file path (required)")
	factsExportCmd.MarkFlagRequired("output")

	// Add commands to facts command
	factsCmd.AddCommand(factsGatherCmd)
	factsCmd.AddCommand(factsListCmd)
	factsCmd.AddCommand(factsValidateCmd)
	factsCmd.AddCommand(factsExportCmd)

	// Add facts command to root
	RootCmd.AddCommand(factsCmd)
}
