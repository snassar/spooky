package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"spooky/internal/facts"
	"spooky/internal/schemas"
	"spooky/internal/ssh"

	"github.com/spf13/cobra"
)

var factsCmd = &cobra.Command{
	Use:   "facts",
	Short: "Gather facts from remote machines",
	Long: `Gather facts from remote machines via SSH.

This command collects three types of facts:
- Basic facts: System information gathered via SSH commands
- Enhanced facts: Detailed system information from spooky-facts tool (if available)
- Custom facts: Age-encrypted custom facts (if available)

The facts are collected in parallel from all machines defined in the project configuration.`,
}

var gatherFactsCmd = &cobra.Command{
	Use:   "gather",
	Short: "Gather facts from all machines",
	Run: func(cmd *cobra.Command, args []string) {
		// Load project configuration
		projectConfig, err := loadProjectConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading project config: %v\n", err)
			os.Exit(1)
		}

		// Load SSH configuration
		sshConfig, err := loadSSHConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading SSH config: %v\n", err)
			os.Exit(1)
		}

		// Create SSH manager
		sshManager := ssh.NewSimpleSSHManager(nil, sshConfig) // TODO: Add age encryption

		// Create facts gatherer
		gatherer := facts.NewGatherer(sshManager, projectConfig)

		// Get machines from configuration
		machines, err := getMachinesFromConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting machines: %v\n", err)
			os.Exit(1)
		}

		if len(machines) == 0 {
			fmt.Println("No machines configured for facts gathering")
			return
		}

		fmt.Printf("Gathering facts from %d machines...\n", len(machines))

		// Gather facts from all machines
		ctx := context.Background()
		machineFacts, err := gatherer.GatherFactsFromMachines(ctx, machines)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error gathering facts: %v\n", err)
			os.Exit(1)
		}

		// Process results
		successCount := 0
		errorCount := 0
		for _, machineFact := range machineFacts {
			if machineFact.Error != nil {
				fmt.Printf("❌ Failed to gather facts from %s: %v\n", machineFact.Machine.Hostname, machineFact.Error)
				errorCount++
			} else {
				fmt.Printf("✅ Successfully gathered facts from %s\n", machineFact.Machine.Hostname)
				successCount++
			}
		}

		fmt.Printf("\nFacts gathering completed: %d successful, %d failed\n", successCount, errorCount)

		// Export facts to HCL
		combinedFacts, err := gatherer.ExportFacts(machineFacts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error exporting facts: %v\n", err)
			os.Exit(1)
		}

		// Write facts to file
		outputPath := "facts.hcl"
		if len(args) > 0 {
			outputPath = args[0]
		}

		err = writeFactsToFile(combinedFacts, outputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing facts to file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Facts written to: %s\n", outputPath)
	},
}

var exportFactsCmd = &cobra.Command{
	Use:   "export [output-file]",
	Short: "Export gathered facts to HCL format",
	Run: func(cmd *cobra.Command, args []string) {
		// This would be used to export facts that were previously gathered
		// For now, it's a placeholder
		fmt.Println("Export command not yet implemented")
	},
}

func init() {
	factsCmd.AddCommand(gatherFactsCmd)
	factsCmd.AddCommand(exportFactsCmd)
	RootCmd.AddCommand(factsCmd)
}

// loadProjectConfig loads the project configuration
func loadProjectConfig() (*schemas.ProjectV1, error) {
	// TODO: Implement project config loading
	// For now, return a default configuration
	return &schemas.ProjectV1{
		FactsTimeout:            30,
		FactsParallelCollection: 10,
		FactsRetryAttempts:      3,
		FactsRetryDelay:         5,
	}, nil
}

// loadSSHConfig loads the SSH configuration
func loadSSHConfig() (*schemas.SpookySSHV1, error) {
	// TODO: Implement SSH config loading
	// For now, return a default configuration
	return &schemas.SpookySSHV1{
		Timeout:                   30,
		KeepaliveInterval:         60,
		KeepaliveCount:            3,
		KeyScanTimeout:            10,
		KnownHostsStrict:          true,
		Compression:               false,
		CompressionLevel:          6,
		TCPKeepAlive:              true,
		TCPKeepAliveCount:         3,
		TCPKeepAliveIdle:          60,
		TCPKeepAliveInterval:      10,
		TCPKeepAliveProbeInterval: 5,
	}, nil
}

// getMachinesFromConfig gets machines from the configuration
func getMachinesFromConfig() ([]*schemas.MachinesMachineV1, error) {
	// TODO: Implement machine loading from configuration
	// For now, return an empty list
	return []*schemas.MachinesMachineV1{}, nil
}

// writeFactsToFile writes facts to an HCL file
func writeFactsToFile(facts *schemas.FactsV1, outputPath string) error {
	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputPath)
	if outputDir != "." {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// TODO: Implement HCL writing
	// For now, just create an empty file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Write a placeholder comment
	_, err = file.WriteString("# Facts gathered by spooky\n# TODO: Implement HCL generation\n")
	if err != nil {
		return fmt.Errorf("failed to write to output file: %w", err)
	}

	return nil
}
