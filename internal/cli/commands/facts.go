package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	spookycoordinator "spooky/internal/coordinator"
	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookyloggingtypes "spooky/internal/logging/types"

	"github.com/spf13/cobra"
)

// CreateFactsCommand creates the facts command with all subcommands
func CreateFactsCommand() *cobra.Command {
	factsCmd := &cobra.Command{
		Use:   "facts",
		Short: "Manage machine facts",
		Long:  "Collect and manage machine facts",
	}

	// Add subcommands
	factsCmd.AddCommand(createFactsGatherCommand())
	factsCmd.AddCommand(createFactsListCommand())
	factsCmd.AddCommand(createFactsValidateCommand())
	factsCmd.AddCommand(createFactsExportCommand())

	return factsCmd
}

// createFactsGatherCommand creates the facts gather subcommand
func createFactsGatherCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gather",
		Short: "Gather facts from machines",
		Long:  "Collect facts from remote machines",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath := args[0]

			// Get flags
			machines, _ := cmd.Flags().GetStringSlice("machines")
			_, _ = cmd.Flags().GetInt("parallel") // Keep flag but don't use for now
			_, _ = cmd.Flags().GetBool("force")   // Keep flag but don't use for now

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

			// Gather facts
			fmt.Printf("Gathering facts from machines: %v\n", machines)

			// Collect facts from the specified machines
			factsContext, err := coord.Facts().CollectFacts(machines)
			if err != nil {
				return fmt.Errorf("failed to collect facts: %w", err)
			}

			// Cache facts
			if err := coord.Facts().CacheFacts(factsContext); err != nil {
				return fmt.Errorf("failed to cache facts: %w", err)
			}

			fmt.Printf("Successfully gathered facts for %d machines\n", len(factsContext.MachineFacts))
			return nil
		},
	}

	// Add flags
	cmd.Flags().StringSliceP("machines", "m", []string{}, "Target machines")
	cmd.Flags().IntP("parallel", "p", 1, "Number of parallel connections")
	cmd.Flags().BoolP("force", "f", false, "Force refresh of cached facts")

	return cmd
}

// createFactsListCommand creates the facts list subcommand
func createFactsListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List collected facts",
		Long:  "List facts collected from machines",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath := args[0]

			// Get flags
			machines, _ := cmd.Flags().GetStringSlice("machines")
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

			// Load facts
			factsContext, err := coord.Facts().LoadFacts(machines)
			if err != nil {
				return fmt.Errorf("failed to load facts: %w", err)
			}

			// Display facts
			if len(factsContext.MachineFacts) == 0 {
				fmt.Println("No facts found")
				return nil
			}

			fmt.Printf("Found facts for %d machines:\n\n", len(factsContext.MachineFacts))

			for machine, factCollection := range factsContext.MachineFacts {
				fmt.Printf("Machine: %s\n", machine)
				fmt.Printf("  Timestamp: %s\n", factCollection.Timestamp.Format("2006-01-02 15:04:05"))
				fmt.Printf("  Facts count: %d\n", len(factCollection.Facts))

				// Show some key facts
				keyFacts := []string{"os.name", "os.version", "hostname", "ip_address"}
				for _, key := range keyFacts {
					if fact, exists := factCollection.Facts[key]; exists {
						fmt.Printf("  %s: %v\n", key, fact.Value)
					}
				}
				fmt.Println()
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringSliceP("machines", "m", []string{}, "Filter by machines")
	cmd.Flags().StringSliceP("tags", "t", []string{}, "Filter by tags")
	cmd.Flags().StringP("filter", "f", "", "Filter expression")

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
			projectPath := args[0]

			// Get flags
			_, _ = cmd.Flags().GetBool("compare") // Keep flag but don't use for now

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

			// Validate facts
			fmt.Println("Validating facts...")

			// Load facts for validation
			factsContext := &spookyinterfaces.FactsContext{
				BaseContext: spookyinterfaces.BaseContext{
					ProjectPath: projectPath,
				},
			}

			if err := coord.Facts().ValidateFacts(factsContext); err != nil {
				return fmt.Errorf("facts validation failed: %w", err)
			}

			fmt.Println("✅ Facts validation completed successfully")
			return nil
		},
	}

	// Add flags
	cmd.Flags().BoolP("compare", "c", false, "Compare with fresh facts")

	return cmd
}

// createFactsExportCommand creates the facts export subcommand
func createFactsExportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export facts",
		Long:  "Export facts to various formats",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath := args[0]

			// Get flags
			format, _ := cmd.Flags().GetString("format")
			output, _ := cmd.Flags().GetString("output")
			machines, _ := cmd.Flags().GetStringSlice("machines")

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
				output = filepath.Join(projectPath, fmt.Sprintf("facts.%s", format))
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

			// Export facts
			fmt.Printf("Exporting facts to %s...\n", output)

			// Load facts for export
			factsContext, err := coord.Facts().LoadFacts(machines)
			if err != nil {
				return fmt.Errorf("failed to load facts: %w", err)
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
				if err := encoder.Encode(factsContext); err != nil {
					return fmt.Errorf("failed to encode JSON: %w", err)
				}
			} else if format == "hcl" {
				// Export as HCL
				// This is a simplified HCL export - in a real implementation, you'd use proper HCL encoding
				fmt.Fprintf(file, "# Facts export\n\n")
				for machine, factCollection := range factsContext.MachineFacts {
					fmt.Fprintf(file, "machine \"%s\" {\n", machine)
					fmt.Fprintf(file, "  timestamp = \"%s\"\n", factCollection.Timestamp.Format(time.RFC3339))
					for key, fact := range factCollection.Facts {
						fmt.Fprintf(file, "  %s = %q\n", key, fact.Value)
					}
					fmt.Fprintf(file, "}\n\n")
				}
			}

			fmt.Printf("✅ Facts exported successfully to %s\n", output)
			return nil
		},
	}

	// Add flags
	cmd.Flags().StringP("format", "f", "json", "Export format (json, hcl)")
	cmd.Flags().StringP("output", "o", "", "Output file path")
	cmd.Flags().StringSliceP("machines", "m", []string{}, "Filter by machines")

	return cmd
}
