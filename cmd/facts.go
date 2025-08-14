// Package cmd provides command implementations for spooky CLI.
package cmd

import (
	"context"
	"fmt"
	"strings"
	"sync"

	spookyfacts "spooky/internal/facts"
	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookymachines "spooky/internal/machines"
	spookyssh "spooky/internal/ssh"
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

	// Create SSH manager for facts collection
	sshManager := spookyssh.NewManager(factsLogger)

	// Initialize facts components
	collector := spookyfacts.NewSystemFactCollector(sshManager)
	manager := spookyfacts.NewManager(collector, nil, factsLogger)

	// Create facts integration
	factsManager = spookyfacts.NewIntegration(manager)

	return nil
}

// factsCmd represents the facts command
var factsCmd = &cobra.Command{
	Use:   "facts",
	Short: "Export machine facts",
	Long: `Export machine facts to files in various formats.

Facts are system information collected from machines including OS details,
hardware information, network configuration, and custom data. Facts can be
exported to JSON or HCL format for backup, analysis, or transfer to other systems.`,
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
  spooky facts export ./my-project --machine web-server --output web-server-facts.hcl
  spooky facts export ./my-project --tags environment=production --output prod-facts.hcl
  spooky facts export ./my-project --groups webservers,database --output app-facts.hcl
  spooky facts export ./my-project --tags role=web --groups production --output web-prod-facts.hcl`,
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
	tagsFilter, _ := cmd.Flags().GetStringSlice("tags")
	groupsFilter, _ := cmd.Flags().GetStringSlice("groups")
	parallel, _ := cmd.Flags().GetInt("parallel")
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Load machines from project
	machines, err := loadMachinesFromProject(projectPath)
	if err != nil {
		return fmt.Errorf("failed to load machines: %w", err)
	}

	// Filter machines if specified
	if machineFilter != "" || len(tagsFilter) > 0 || len(groupsFilter) > 0 {
		machines = filterMachinesAdvanced(machines, machineFilter, tagsFilter, groupsFilter)
		if len(machines) == 0 {
			return fmt.Errorf("no machines found matching filters")
		}
	}

	if verbose {
		factsLogger.Info("Starting fact collection for export", map[string]interface{}{
			"project_path":  projectPath,
			"machine_count": len(machines),
			"parallel":      parallel,
		})
	}

	// Collect facts from machines in parallel
	successCount, errorCount := collectFactsParallel(ctx, machines, factsManager, parallel, verbose, factsLogger)

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

	// Extract machine hostnames for export
	var machineIDs []string
	for _, machine := range machines {
		machineIDs = append(machineIDs, machine.Hostname)
	}

	if len(machineIDs) == 0 {
		return fmt.Errorf("no machines found to export")
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
	// Create a log manager and get a logger for the machines manager
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("machines")

	// Create a machines loader
	loader := spookymachines.NewLoader(logger)

	// Create a machines validator
	validator := spookymachines.NewValidator(logger)

	// Create a machines manager with the loader and validator
	manager := spookymachines.NewManager(logger, loader, validator)

	// Load machines from the project using the manager
	// This supports both machines.hcl file and machines/ directory
	machines, err := manager.LoadMachines(context.Background(), projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load machines from project %s: %w", projectPath, err)
	}

	return machines, nil
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

func filterMachinesAdvanced(machines []spookytypes.Machine, machineFilter string, tagsFilter []string, groupsFilter []string) []spookytypes.Machine {
	var filtered []spookytypes.Machine

	for _, machine := range machines {
		// Check if machine matches all filters
		if matchesMachineFilter(machine, machineFilter) &&
			matchesTagsFilter(machine, tagsFilter) &&
			matchesGroupsFilter(machine, groupsFilter) {
			filtered = append(filtered, machine)
		}
	}

	return filtered
}

func matchesMachineFilter(machine spookytypes.Machine, filter string) bool {
	if filter == "" {
		return true
	}
	return strings.Contains(machine.Hostname, filter)
}

func matchesTagsFilter(machine spookytypes.Machine, tagsFilter []string) bool {
	if len(tagsFilter) == 0 {
		return true
	}

	// Check if machine has any of the specified tags
	for _, tag := range tagsFilter {
		if machine.Tags != nil {
			for machineTagKey, machineTagValue := range machine.Tags {
				// Support both key=value and key-only filtering
				if strings.Contains(tag, "=") {
					// Key=value format
					parts := strings.SplitN(tag, "=", 2)
					if len(parts) == 2 && machineTagKey == parts[0] && machineTagValue == parts[1] {
						return true
					}
				} else {
					// Key-only format
					if machineTagKey == tag {
						return true
					}
				}
			}
		}
	}

	return false
}

func matchesGroupsFilter(machine spookytypes.Machine, groupsFilter []string) bool {
	if len(groupsFilter) == 0 {
		return true
	}

	// Check if machine belongs to any of the specified groups
	for _, group := range groupsFilter {
		if machine.Groups != nil {
			for _, machineGroup := range machine.Groups {
				if machineGroup == group {
					return true
				}
			}
		}
	}

	return false
}

// collectFactsParallel collects facts from multiple machines in parallel
func collectFactsParallel(ctx context.Context, machines []spookytypes.Machine, factsManager spookyinterfaces.FactsIntegration, parallel int, verbose bool, logger spookytypeslogging.Logger) (int, int) {
	successCount := 0
	errorCount := 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Create semaphore for parallel execution
	semaphore := make(chan struct{}, parallel)

	for _, machine := range machines {
		wg.Add(1)
		go func(m spookytypes.Machine) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			if verbose {
				logger.Info("Collecting facts for export", map[string]interface{}{
					"machine": m.Hostname,
				})
			}

			// Collect facts using the integration with actual machine object
			facts, err := factsManager.CollectFacts(ctx, &m)
			if err != nil {
				mu.Lock()
				errorCount++
				mu.Unlock()
				logger.Error("Failed to collect facts for export", err, map[string]interface{}{
					"machine": m.Hostname,
				})
				return
			}

			// Store facts
			if err := factsManager.StoreFacts(ctx, facts); err != nil {
				mu.Lock()
				errorCount++
				mu.Unlock()
				logger.Error("Failed to store facts for export", err, map[string]interface{}{
					"machine": m.Hostname,
				})
				return
			}

			mu.Lock()
			successCount++
			mu.Unlock()

			if verbose {
				logger.Info("Successfully collected facts for export", map[string]interface{}{
					"machine": m.Hostname,
				})
			}
		}(machine)
	}

	wg.Wait()
	return successCount, errorCount
}

func init() {

	// Add flags to facts export command
	factsExportCmd.Flags().String("format", "hcl", "Export format (hcl, json)")
	factsExportCmd.Flags().String("output", "", "Output file path (required)")
	factsExportCmd.Flags().String("machine", "", "Filter to specific machine")
	factsExportCmd.Flags().StringSlice("tags", []string{}, "Filter by tags (supports key=value or key-only)")
	factsExportCmd.Flags().StringSlice("groups", []string{}, "Filter by groups")
	factsExportCmd.Flags().Int("parallel", 1, "Number of parallel workers")
	factsExportCmd.Flags().Bool("verbose", false, "Verbose output")
	if err := factsExportCmd.MarkFlagRequired("output"); err != nil {
		panic(fmt.Sprintf("failed to mark output flag as required: %v", err))
	}

	// Add commands to facts command
	factsCmd.AddCommand(factsExportCmd)

	// Add facts command to root
	RootCmd.AddCommand(factsCmd)
}
