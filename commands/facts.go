package commands

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"spooky/internal/facts"
	internalhcl "spooky/internal/hcl"
	"spooky/internal/logging"
	"spooky/internal/schemas"
	"spooky/internal/ssh"
	"spooky/internal/utilities"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
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
		logger := logging.GetGlobalLogger()
		logger.Debug("gatherFactsCmd started")
		// Load project configuration
		projectConfig, err := loadProjectConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading project config: %v\n", err)
			os.Exit(1)
		}

		// Load SSH configuration
		sshConfig, err := loadSSHConfig()
		if err != nil {
			logger.Error("failed to load SSH configuration", slog.String("error", err.Error()))
			os.Exit(1)
		}

		// Create SSH manager
		sshManager := ssh.NewSimpleSSHManager(nil, sshConfig) // Note: Age encryption not yet implemented for facts gathering

		// Create facts gatherer
		gatherer := facts.NewGatherer(sshManager, projectConfig)

		// Get machines from configuration
		logger.Debug("about to call getMachinesFromConfig")
		machines, err := getMachinesFromConfig()
		logger.Debug("getMachinesFromConfig completed",
			slog.Int("machine_count", len(machines)),
			slog.String("error", fmt.Sprintf("%v", err)))
		if err != nil {
			logger.Error("failed to get machines from configuration", slog.String("error", err.Error()))
			os.Exit(1)
		}

		if len(machines) == 0 {
			logger.Info("No machines configured for facts gathering")
			return
		}

		logger.Info("Gathering facts from machines", slog.Int("machine_count", len(machines)))

		// Gather facts from all machines
		ctx := context.Background()
		machineFacts, err := gatherer.GatherFactsFromMachines(ctx, machines)
		if err != nil {
			logger.Error("failed to gather facts from machines", slog.String("error", err.Error()))
			os.Exit(1)
		}

		// Process results
		successCount := 0
		errorCount := 0
		for _, machineFact := range machineFacts {
			if machineFact.Error != nil {
				logger.Error("❌ Failed to gather facts from machine",
					slog.String("hostname", machineFact.Machine.Hostname),
					slog.String("error", machineFact.Error.Error()))
				errorCount++
			} else {
				logger.Info("✅ Successfully gathered facts from machine",
					slog.String("hostname", machineFact.Machine.Hostname))
				successCount++
			}
		}

		logger.Info("Facts gathering completed",
			slog.Int("successful", successCount),
			slog.Int("failed", errorCount))

		// Export facts to HCL
		combinedFacts, err := gatherer.ExportFacts(machineFacts)
		if err != nil {
			logger.Error("failed to export facts", slog.String("error", err.Error()))
			os.Exit(1)
		}

		// Write facts to file
		outputPath := "facts.hcl"
		if len(args) > 0 {
			outputPath = args[0]
		}

		err = writeFactsToFile(combinedFacts, outputPath)
		if err != nil {
			logger.Error("failed to write facts to file",
				slog.String("file", outputPath),
				slog.String("error", err.Error()))
			os.Exit(1)
		}

		logger.Info("Facts written to file", slog.String("file", outputPath))
	},
}

var exportFactsCmd = &cobra.Command{
	Use:   "export [output-file]",
	Short: "Export gathered facts to HCL format",
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetGlobalLogger()
		// This would be used to export facts that were previously gathered
		// For now, it's a placeholder
		logger.Info("Export command not yet implemented")
	},
}

func init() {
	factsCmd.AddCommand(gatherFactsCmd)
	factsCmd.AddCommand(exportFactsCmd)
	RootCmd.AddCommand(factsCmd)
}

// loadProjectConfig loads the project configuration
func loadProjectConfig() (*schemas.ProjectV1, error) {
	logger := logging.GetGlobalLogger()
	logger.Debug("loading project configuration")
	// Look for project.hcl in current directory
	projectHCLPath := "project.hcl"
	if _, err := os.Stat(projectHCLPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("project.hcl not found in current directory")
	}

	// Read and parse project.hcl
	content, err := os.ReadFile(projectHCLPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read project.hcl: %w", err)
	}

	// Use the simplified validator to parse and validate the project configuration
	validator := schemas.NewSimpleValidator()
	result, err := validator.ValidateHCLContent("project", string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to validate project.hcl: %w", err)
	}

	if !result.IsValid {
		return nil, fmt.Errorf("project.hcl validation failed: %s", result.Errors[0].Message)
	}

	// Parse the HCL content to extract project configuration
	parsedData, err := validator.ParseHCLContent(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse project.hcl: %w", err)
	}

	// Extract project configuration from parsed data
	projectBlock, exists := parsedData["project"]
	if !exists {
		return nil, fmt.Errorf("no project block found in project.hcl")
	}

	projectMap, ok := projectBlock.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid project block structure")
	}

	// Create a default project configuration
	projectConfig := &schemas.ProjectV1{
		FactsTimeout:            30,
		FactsParallelCollection: 10,
		FactsRetryAttempts:      3,
		FactsRetryDelay:         5,
	}

	// Extract project name from block label or name field
	if name, exists := projectMap["name"]; exists {
		if nameStr, ok := name.(string); ok {
			projectConfig.Name = nameStr
		}
	}

	// Extract description
	if desc, exists := projectMap["description"]; exists {
		if descStr, ok := desc.(string); ok {
			projectConfig.Description = descStr
		}
	}

	return projectConfig, nil
}

// loadSSHConfig loads the SSH configuration
func loadSSHConfig() (*schemas.SpookySSHV1, error) {
	logger := logging.GetGlobalLogger()
	logger.Debug("loading SSH configuration")
	// For now, return a default SSH configuration
	// In a full implementation, this would load from spooky.hcl or environment
	return &schemas.SpookySSHV1{
		Timeout:                   30,
		KeepaliveInterval:         60,
		KeepaliveCount:            3,
		KeyScanTimeout:            10,
		KnownHostsStrict:          false,        // Deprecated, use KnownHostsMode instead
		KnownHostsMode:            "accept-new", // Accept new hosts silently, warn on changed keys
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
	logger := logging.GetGlobalLogger()
	logger.Debug("loading machines configuration")

	// Look for machines.hcl in current directory
	machinesHCLPath := "machines.hcl"
	if _, err := os.Stat(machinesHCLPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("machines.hcl not found in current directory")
	}

	// Read and parse machines.hcl using the HCL library properly
	content, err := os.ReadFile(machinesHCLPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read machines.hcl: %w", err)
	}

	// Parse HCL using the library
	file, diags := hclsyntax.ParseConfig(content, machinesHCLPath, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse machines.hcl: %v", diags)
	}

	// Define the schema for machines block
	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       ResourceTypeMachines,
				LabelNames: []string{},
			},
		},
	}

	// Extract the machines block
	bodyContent, diags := file.Body.Content(schema)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to decode machines block: %v", diags)
	}

	if len(bodyContent.Blocks) == 0 {
		return nil, fmt.Errorf("no machines block found in machines.hcl")
	}

	if len(bodyContent.Blocks) > 1 {
		return nil, fmt.Errorf("multiple machines blocks found in machines.hcl")
	}

	machinesBlock := bodyContent.Blocks[0]

	// Define schema for machine blocks inside machines
	machineSchema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "machine",
				LabelNames: []string{"name"},
			},
		},
	}

	// Extract machine blocks
	machineContent, diags := machinesBlock.Body.Content(machineSchema)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to decode machine blocks: %v", diags)
	}

	logger.Debug("found machine blocks", slog.Int("count", len(machineContent.Blocks)))

	// Create machines slice
	var machines []*schemas.MachinesMachineV1

	// Process each machine block
	for _, machineBlock := range machineContent.Blocks {
		machineName := machineBlock.Labels[0]
		logger.Debug("processing machine", slog.String("machine_name", machineName))

		// Create a new machine struct
		machine := &schemas.MachinesMachineV1{
			Description:       "",
			Hostname:          "",
			Port:              22, // default
			User:              "",
			Authentication:    schemas.MachinesMachineAuthenticationV1{},
			ConnectionTimeout: 30, // default
			MaxRetries:        3,  // default
			RetryDelay:        5,  // default
			Facts:             schemas.MachinesMachineFactsV1{},
		}

		// Use a schema that allows authentication blocks
		machineSchema := &hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{Name: "hostname", Required: true},
				{Name: "port", Required: false},
				{Name: "user", Required: true},
			},
			Blocks: []hcl.BlockHeaderSchema{
				{
					Type:       "authentication",
					LabelNames: []string{"method"},
				},
			},
		}

		// Extract the machine content
		machineContent, diags := machineBlock.Body.Content(machineSchema)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to decode machine %s: %v", machineName, diags)
		}

		// Extract attributes
		if hostnameAttr, exists := machineContent.Attributes["hostname"]; exists {
			var hostname string
			if diags := gohcl.DecodeExpression(hostnameAttr.Expr, nil, &hostname); diags.HasErrors() {
				return nil, fmt.Errorf("failed to decode hostname for machine %s: %v", machineName, diags)
			}
			machine.Hostname = hostname
		}

		if portAttr, exists := machineContent.Attributes["port"]; exists {
			var port int
			if diags := gohcl.DecodeExpression(portAttr.Expr, nil, &port); diags.HasErrors() {
				return nil, fmt.Errorf("failed to decode port for machine %s: %v", machineName, diags)
			}
			machine.Port = port
		}

		if userAttr, exists := machineContent.Attributes["user"]; exists {
			var user string
			if diags := gohcl.DecodeExpression(userAttr.Expr, nil, &user); diags.HasErrors() {
				return nil, fmt.Errorf("failed to decode user for machine %s: %v", machineName, diags)
			}
			machine.User = user
		}

		// Handle authentication block
		if len(machineContent.Blocks) > 0 {
			authBlock := machineContent.Blocks[0]
			authConfig, err := utilities.ParseAuthenticationBlock(authBlock, machineName)
			if err != nil {
				return nil, err
			}
			machine.Authentication = *authConfig
		}

		// Validate required fields
		if machine.Hostname == "" {
			return nil, fmt.Errorf("machine %s missing required hostname", machineName)
		}
		if machine.User == "" {
			return nil, fmt.Errorf("machine %s missing required user", machineName)
		}

		logger.Debug("successfully parsed machine",
			slog.String("machine_name", machineName),
			slog.String("user", machine.User),
			slog.String("hostname", machine.Hostname),
			slog.Int("port", machine.Port))
		machines = append(machines, machine)
	}

	logger.Debug("getMachinesFromConfig completed", slog.Int("machine_count", len(machines)))
	return machines, nil
}

// writeFactsToFile writes facts to an HCL file
func writeFactsToFile(facts *schemas.FactsV1, outputPath string) error {
	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputPath)
	if outputDir != "." {
		if err := os.MkdirAll(outputDir, 0o755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Generate HCL from facts using the generic HCL generator
	hclContent, err := internalhcl.GenerateHCL(facts, "facts")
	if err != nil {
		return fmt.Errorf("failed to generate HCL from facts: %w", err)
	}

	// Add header comments
	hclContent = "# Facts gathered by spooky\n" +
		fmt.Sprintf("# Generated at: %s\n\n", time.Now().Format(time.RFC3339)) +
		hclContent

	// Write HCL content to file
	err = utilities.WriteFile(outputPath, hclContent)
	if err != nil {
		return fmt.Errorf("failed to write HCL to file: %w", err)
	}

	return nil
}
