// Package cmd provides command implementations for spooky CLI.
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookymachines "spooky/internal/machines"
	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"

	"github.com/spf13/cobra"
)

// Machine status constants
const (
	MachineStatusOnline        = "online"
	MachineStatusOffline       = "offline"
	MachineStatusError         = "error"
	MachineStatusAuthenticated = "authenticated"
	MachineStatusAuthFailed    = "auth_failed"
)

// Global instances for machine dependency injection
var (
	machinesManager spookyinterfaces.MachinesIntegration
	machinesLogger  spookytypeslogging.Logger
)

// InitializeMachinesDependencies initializes machine-related dependencies
func InitializeMachinesDependencies() error {
	// Create log manager for machines component
	logManager := spookylogging.NewLogManager()
	machinesLogger = logManager.GetLogger("machines")

	// Initialize machine components
	validator := spookymachines.NewValidator(machinesLogger)
	loader := spookymachines.NewLoader(machinesLogger)
	machinesManager = spookymachines.NewManager(machinesLogger, loader, validator)

	return nil
}

// machinesCmd represents the machines command
var machinesCmd = &cobra.Command{
	Use:   "machines",
	Short: "Manage machine inventory",
	Long: `Manage machine inventory including listing, validation, and connectivity testing.

Machine inventory is defined in machines.hcl files within spooky projects and contains
SSH connection details, authentication information, and machine metadata.`,
}

// machinesListCmd represents the machines list command
var machinesListCmd = &cobra.Command{
	Use:   "list [project-path]",
	Short: "List machines in a project",
	Long: `List all machines defined in the project's machine inventory.

This command reads machines.hcl files and displays information about all configured
machines including hostname, host, user, and connection status.`,
	Args: cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		return handleMachinesList(args[0])
	},
}

// machinesValidateCmd represents the machines validate command
var machinesValidateCmd = &cobra.Command{
	Use:   "validate [project-path]",
	Short: "Validate machine inventory",
	Long: `Validate machine inventory configuration and connectivity.

This command validates that all machines in the inventory have proper configuration
including required fields, valid authentication methods, and SSH settings.`,
	Args: cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		return handleMachinesValidate(args[0])
	},
}

// machinesPingCmd represents the machines ping command
var machinesPingCmd = &cobra.Command{
	Use:   "ping [project-path]",
	Short: "Ping machines to test connectivity",
	Long: `Ping machines to test network connectivity and SSH accessibility.

This command tests connectivity to all machines in the inventory and reports
their status including response times and connection success/failure.

By default, shows minimal output for working machines and detailed output for
problematic machines. Use --verbose for detailed output on all machines.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleMachinesPing(cmd, args[0])
	},
}

// machinesExportCmd represents the machines export command
var machinesExportCmd = &cobra.Command{
	Use:   "export [project-path]",
	Short: "Export machines to HCL format",
	Long: `Export machine inventory to HCL format according to machines schema.

This command exports all machines from the project's inventory to a single HCL file
that follows the machines schema specification. The exported file can be used for
backup, analysis, or transfer to other systems.

Examples:
  spooky machines export ./my-project --output machines.hcl
  spooky machines export ./my-project --machine web-server --output web-server.hcl
  spooky machines export ./my-project --tags environment=production --output prod-machines.hcl`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleMachinesExport(cmd, args[0])
	},
}

// machinesEncryptCmd represents the machines encrypt command
var machinesEncryptCmd = &cobra.Command{
	Use:   "encrypt [project-path]",
	Short: "Encrypt machines in a project",
	Long: `Encrypt machines in a project using age encryption.

This command processes machines.hcl files and encrypts any authentication
credentials that have encrypted=true set. It will re-encrypt if identities/recipients have changed.

Examples:
  spooky machines encrypt ./my-project
  spooky machines encrypt ./my-project --dry-run`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath := args[0]
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		return handleMachinesEncrypt(projectPath, dryRun)
	},
}

// handleMachinesList handles listing machines using the MachinesIntegration interface

// Helper function for extracting source file information
func extractSourceFile(machine *spookytypes.Machine) string {
	sourceFile := "unknown"
	if machine.MachineMetadata != nil && machine.MachineMetadata.CustomFields != nil {
		if src, exists := machine.MachineMetadata.CustomFields["source_file"]; exists {
			sourceFile = src
		}
	}
	return sourceFile
}

// Helper function for grouping machines by source
func groupMachinesBySource(machines []spookytypes.Machine) map[string][]spookytypes.Machine {
	machinesBySource := make(map[string][]spookytypes.Machine)

	for idx := range machines {
		machine := &machines[idx]
		sourceFile := extractSourceFile(machine)
		machinesBySource[sourceFile] = append(machinesBySource[sourceFile], *machine)
	}

	return machinesBySource
}

// Helper function for displaying empty state
func displayEmptyState() {
	fmt.Printf("No machines found in inventory.\n")
}

// Helper function for displaying a single machine
func displayMachine(machine *spookytypes.Machine, index int) {
	fmt.Printf("%d. %s (%s)\n", index+1, machine.Hostname, machine.Host)
	fmt.Printf("   User: %s\n", machine.User)
	fmt.Printf("   Port: %d\n", machine.Port)

	// Show environment if available
	if machine.MachineMetadata != nil && machine.MachineMetadata.Environment != "" {
		fmt.Printf("   Environment: %s\n", machine.MachineMetadata.Environment)
	}

	if len(machine.Groups) > 0 {
		fmt.Printf("   Groups: %v\n", machine.Groups)
	}

	if len(machine.Roles) > 0 {
		fmt.Printf("   Roles: %v\n", machine.Roles)
	}

	if len(machine.Tags) > 0 {
		fmt.Printf("   Tags: %v\n", machine.Tags)
	}

	fmt.Printf("\n")
}

// Helper function for displaying machines grouped by source
func displayMachinesBySource(machinesBySource map[string][]spookytypes.Machine) {
	for sourceFile, sourceMachines := range machinesBySource {
		fmt.Printf("📁 Source: %s (%d machines)\n", sourceFile, len(sourceMachines))
		fmt.Printf("%s\n", strings.Repeat("─", 50))

		for i := range sourceMachines {
			machine := &sourceMachines[i]
			displayMachine(machine, i)
		}
	}
}

// Helper function for displaying summary information
func displayMachinesSummary(projectPath string, machineCount int) {
	fmt.Printf("🔍 Loading machines from project: %s\n", projectPath)
	fmt.Printf("📊 Found %d machines:\n\n", machineCount)
}

// handleMachinesList handles listing machines using the MachinesIntegration interface
func handleMachinesList(projectPath string) error {
	ctx := context.Background()

	// Initialize dependencies
	if err := initializeMachinesDependenciesIfNeeded(); err != nil {
		return err
	}

	// Load machines
	machines, err := loadProjectMachines(ctx, projectPath)
	if err != nil {
		return err
	}

	// Display summary
	displayMachinesSummary(projectPath, len(machines))

	// Handle empty state
	if len(machines) == 0 {
		displayEmptyState()
		return nil
	}

	// Group and display machines
	machinesBySource := groupMachinesBySource(machines)
	displayMachinesBySource(machinesBySource)

	return nil
}

// handleMachinesValidate handles machine validation using the MachinesIntegration interface
func handleMachinesValidate(projectPath string) error {
	ctx := context.Background()

	// Initialize dependencies if not already done
	if machinesManager == nil {
		if err := InitializeMachinesDependencies(); err != nil {
			return fmt.Errorf("failed to initialize machines dependencies: %w", err)
		}
	}

	fmt.Printf("🔍 Validating machines in project: %s\n", projectPath)

	// Load machines using the enhanced manager
	machines, err := machinesManager.LoadMachines(ctx, projectPath)
	if err != nil {
		return fmt.Errorf("failed to load machines: %w", err)
	}

	if len(machines) == 0 {
		fmt.Printf("No machines found in inventory.\n")
		return nil
	}

	fmt.Printf("📊 Validating %d machines...\n", len(machines))

	// Validate machines using the manager
	result, err := machinesManager.ValidateMachines(ctx, machines)
	if err != nil {
		return fmt.Errorf("failed to validate machines: %w", err)
	}

	// Display validation results
	fmt.Printf("\n✅ Validation Results:\n")
	fmt.Printf("%s\n", strings.Repeat("─", 50))

	if len(result.Errors) == 0 && len(result.Warnings) == 0 {
		fmt.Printf("🎉 All machines are valid!\n")
		return nil
	}

	// Display errors
	if len(result.Errors) > 0 {
		fmt.Printf("❌ Errors (%d):\n", len(result.Errors))
		for i := range result.Errors {
			err := &result.Errors[i]
			fmt.Printf("  %d. %s\n", i+1, err.Message)
			if err.Context != nil {
				fmt.Printf("     Context: %v\n", err.Context)
			}
		}
		fmt.Printf("\n")
	}

	// Display warnings
	if len(result.Warnings) > 0 {
		fmt.Printf("⚠️  Warnings (%d):\n", len(result.Warnings))
		for i := range result.Warnings {
			warning := &result.Warnings[i]
			fmt.Printf("  %d. %s\n", i+1, warning.Message)
			if warning.Context != nil {
				fmt.Printf("     Context: %v\n", warning.Context)
			}
		}
		fmt.Printf("\n")
	}

	// Summary
	if len(result.Errors) > 0 {
		fmt.Printf("❌ Validation failed with %d errors and %d warnings\n", len(result.Errors), len(result.Warnings))
		return fmt.Errorf("validation failed")
	}
	fmt.Printf("✅ Validation passed with %d warnings\n", len(result.Warnings))

	return nil
}

// handleMachinesPing handles pinging machines using the MachinesIntegration interface
func handleMachinesPing(cmd *cobra.Command, projectPath string) error {
	ctx := context.Background()

	// Initialize dependencies if not already done
	if machinesManager == nil {
		if err := InitializeMachinesDependencies(); err != nil {
			return fmt.Errorf("failed to initialize machines dependencies: %w", err)
		}
	}

	fmt.Printf("🔍 Loading machines from project: %s\n", projectPath)

	// Load machines using the enhanced manager
	machines, err := machinesManager.LoadMachines(ctx, projectPath)
	if err != nil {
		return fmt.Errorf("failed to load machines: %w", err)
	}

	if len(machines) == 0 {
		fmt.Printf("No machines found in inventory.\n")
		return nil
	}

	// Get auth flag
	auth, _ := cmd.Flags().GetBool("auth")

	if auth {
		fmt.Printf("🔐 Testing connectivity and authentication for %d machines...\n", len(machines))
	} else {
		fmt.Printf("📊 Pinging %d machines...\n", len(machines))
	}

	// Ping machines using the manager
	statuses, err := machinesManager.PingMachines(ctx, machines)
	if err != nil {
		return fmt.Errorf("failed to ping machines: %w", err)
	}

	// If auth flag is set, test authentication for each machine
	if auth {
		fmt.Printf("🔑 Testing authentication...\n")
		for i := range statuses {
			if statuses[i].Status == MachineStatusOnline {
				// Test authentication for this machine
				err := testMachineAuthentication(ctx, &machines[i])
				if err != nil {
					statuses[i].Status = MachineStatusAuthFailed
					statuses[i].Error = fmt.Sprintf("Authentication failed: %v", err)
				} else {
					statuses[i].Status = MachineStatusAuthenticated
					statuses[i].Error = "" // Clear any previous errors
				}
			}
		}
	}

	// Get output format and verbose flags
	format, _ := cmd.Flags().GetString("format")
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Output results based on format
	switch format {
	case "json":
		return outputPingResultsJSON(statuses, verbose)
	default:
		return outputPingResultsText(statuses, verbose)
	}
}

// Helper function for dependency initialization
func initializeMachinesDependenciesIfNeeded() error {
	if machinesManager == nil {
		if err := InitializeMachinesDependencies(); err != nil {
			return fmt.Errorf("failed to initialize machine dependencies: %w", err)
		}
	}
	return nil
}

// Helper function for output path validation
func validateOutputPath(cmd *cobra.Command) (string, error) {
	outputPath, err := cmd.Flags().GetString("output")
	if err != nil {
		return "", fmt.Errorf("failed to get output flag: %w", err)
	}
	if outputPath == "" {
		return "", fmt.Errorf("--output flag is required")
	}
	return outputPath, nil
}

// Helper function for loading machines
func loadProjectMachines(ctx context.Context, projectPath string) ([]spookytypes.Machine, error) {
	machines, err := machinesManager.LoadMachines(ctx, projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load machines from project: %w", err)
	}
	return machines, nil
}

// Helper function for filtering machines by name
func filterMachinesByName(machines []spookytypes.Machine, machineName string) ([]spookytypes.Machine, error) {
	if machineName == "" {
		return machines, nil
	}

	var filtered []spookytypes.Machine
	for idx := range machines {
		machine := &machines[idx]
		if machine.Hostname == machineName {
			filtered = append(filtered, *machine)
			break
		}
	}

	if len(filtered) == 0 {
		return nil, fmt.Errorf("machine '%s' not found in project", machineName)
	}

	return filtered, nil
}

// Helper function for filtering machines by tags
func filterMachinesByTags(ctx context.Context, machines []spookytypes.Machine, tags []string) ([]spookytypes.Machine, error) {
	if len(tags) == 0 {
		return machines, nil
	}

	taggedMachines, err := machinesManager.GetMachinesByTags(ctx, tags)
	if err != nil {
		return nil, fmt.Errorf("failed to filter machines by tags: %w", err)
	}

	return intersectMachines(machines, taggedMachines), nil
}

// Helper function for intersecting machine lists
func intersectMachines(list1, list2 []spookytypes.Machine) []spookytypes.Machine {
	var intersection []spookytypes.Machine

	for idx := range list2 {
		tagged := &list2[idx]
		for idx := range list1 {
			filtered := &list1[idx]
			if tagged.Hostname == filtered.Hostname {
				intersection = append(intersection, *tagged)
				break
			}
		}
	}

	return intersection
}

// Helper function for exporting machines
func exportMachinesToFile(ctx context.Context, machines []spookytypes.Machine, outputPath string) error {
	if err := machinesManager.ExportMachines(ctx, machines, outputPath); err != nil {
		return fmt.Errorf("failed to export machines: %w", err)
	}
	return nil
}

// Helper function for displaying success output
func displayExportSuccess(machineCount int, outputPath string) {
	fmt.Printf("✅ Successfully exported %d machines to: %s\n", machineCount, outputPath)
}

// Helper function for applying all filters
func applyMachineFilters(ctx context.Context, machines []spookytypes.Machine, cmd *cobra.Command) ([]spookytypes.Machine, error) {
	filteredMachines := machines

	// Apply name filter
	if machineName, _ := cmd.Flags().GetString("machine"); machineName != "" {
		nameFiltered, err := filterMachinesByName(filteredMachines, machineName)
		if err != nil {
			return nil, err
		}
		filteredMachines = nameFiltered
	}

	// Apply tags filter
	if tags, _ := cmd.Flags().GetStringArray("tags"); len(tags) > 0 {
		tagsFiltered, err := filterMachinesByTags(ctx, filteredMachines, tags)
		if err != nil {
			return nil, err
		}
		filteredMachines = tagsFiltered
	}

	return filteredMachines, nil
}

// handleMachinesExport handles exporting machines using the MachinesIntegration interface
func handleMachinesExport(cmd *cobra.Command, projectPath string) error {
	ctx := context.Background()

	// Initialize dependencies
	if err := initializeMachinesDependenciesIfNeeded(); err != nil {
		return err
	}

	// Validate output path
	outputPath, err := validateOutputPath(cmd)
	if err != nil {
		return err
	}

	// Load machines from project
	machines, err := loadProjectMachines(ctx, projectPath)
	if err != nil {
		return err
	}

	// Apply filters
	filteredMachines, err := applyMachineFilters(ctx, machines, cmd)
	if err != nil {
		return err
	}

	// Export machines
	if err := exportMachinesToFile(ctx, filteredMachines, outputPath); err != nil {
		return err
	}

	// Display success
	displayExportSuccess(len(filteredMachines), outputPath)

	return nil
}

// handleMachinesEncrypt handles machines encryption
func handleMachinesEncrypt(projectPath string, dryRun bool) error {
	manager := GetIntegrationManager()
	machinesIntegration := manager.GetMachinesIntegration()
	if machinesIntegration == nil {
		return fmt.Errorf("machines integration not available")
	}

	return handleEncryptionOperation(projectPath, dryRun, "Machines", machinesIntegration.EncryptMachines)
}

// outputPingResultsText outputs ping results in text format
func outputPingResultsText(statuses []spookytypes.MachineStatus, verbose bool) error {
	fmt.Printf("\n📊 Ping Results:\n")
	fmt.Printf("Total machines: %d\n\n", len(statuses))

	online := 0
	offline := 0
	errors := 0

	for _, status := range statuses {
		switch status.Status {
		case MachineStatusOnline:
			online++
			if verbose {
				fmt.Printf("%s: %s (latency: %dms)\n",
					status.Machine.Hostname, status.Status, status.Latency)
			} else {
				fmt.Printf("%s: %s\n", status.Machine.Hostname, status.Status)
			}
		case MachineStatusOffline:
			offline++
			if verbose {
				fmt.Printf("%s: %s (latency: %dms; Error: %s)\n",
					status.Machine.Hostname, status.Status, status.Latency,
					getErrorMessage(&status))
			} else {
				fmt.Printf("%s: %s\n", status.Machine.Hostname, status.Status)
			}
		case MachineStatusError:
			errors++
			if verbose {
				fmt.Printf("%s: %s (latency: %dms; Error: %s)\n",
					status.Machine.Hostname, status.Status, status.Latency,
					getErrorMessage(&status))
			} else {
				fmt.Printf("%s: %s\n", status.Machine.Hostname, status.Status)
			}
		default:
			if verbose {
				fmt.Printf("%s: %s (latency: %dms)\n",
					status.Machine.Hostname, status.Status, status.Latency)
			} else {
				fmt.Printf("%s: %s\n", status.Machine.Hostname, status.Status)
			}
		}
	}

	fmt.Printf("\n📈 Summary: %d online, %d offline, %d errors\n", online, offline, errors)
	return nil
}

// outputPingResultsJSON outputs ping results in JSON format (always streaming)
func outputPingResultsJSON(statuses []spookytypes.MachineStatus, verbose bool) error {
	for _, status := range statuses {
		machineJSON := map[string]interface{}{
			"hostname": status.Machine.Hostname,
			"status":   status.Status,
		}

		// Add details only for problematic machines or verbose mode
		if status.Status != MachineStatusOnline || verbose {
			machineJSON["latency_ms"] = status.Latency
			if status.Status != MachineStatusOnline {
				machineJSON["error"] = getErrorMessage(&status)
			}
		}

		jsonData, err := json.Marshal(machineJSON)
		if err != nil {
			return fmt.Errorf("failed to marshal machine JSON: %w", err)
		}

		fmt.Println(string(jsonData))
	}
	return nil
}

// testMachineAuthentication tests authentication for a single machine
func testMachineAuthentication(ctx context.Context, machine *spookytypes.Machine) error {
	// Get actions integration to access SSH manager
	actionsIntegration := getIntegrationManager().GetActionsIntegration()
	if actionsIntegration == nil {
		return fmt.Errorf("actions integration not available")
	}

	// Get SSH manager from actions integration
	sshManager := actionsIntegration.GetSSHManager()
	if sshManager == nil {
		return fmt.Errorf("SSH manager not available")
	}

	// Create SSH client configuration from machine
	clientConfig := &spookytypes.ClientConfig{
		DefaultHost:      machine.Hostname,
		DefaultPort:      machine.Port,
		DefaultUser:      machine.User,
		DefaultTimeout:   time.Duration(machine.ConnectionTimeout) * time.Second,
		MaxRetryAttempts: machine.RetryAttempts,
		RetryDelay:       time.Duration(machine.RetryDelay) * time.Second,
	}

	// Set authentication method based on available credentials
	switch {
	case machine.KeyFile != "":
		clientConfig.DefaultKeyPath = machine.KeyFile
		clientConfig.DefaultAuthMethod = "public_key"
	case machine.Password != "":
		clientConfig.DefaultAuthMethod = "password"
	default:
		return fmt.Errorf("no authentication credentials provided for %s", machine.Hostname)
	}

	// Create connection request with authentication
	connectionRequest := &spookytypes.ConnectionRequest{
		Host:    machine.Hostname,
		Port:    machine.Port,
		User:    machine.User,
		Timeout: time.Duration(machine.ConnectionTimeout) * time.Second,
	}

	// Set authentication method and credentials
	if machine.KeyFile != "" {
		connectionRequest.AuthMethod = "public_key"
		connectionRequest.KeyPath = machine.KeyFile
		if machine.Passphrase != "" {
			connectionRequest.Passphrase = machine.Passphrase
		}
	} else if machine.Password != "" {
		connectionRequest.AuthMethod = "password"
		connectionRequest.Password = machine.Password
	}

	// Attempt to connect and authenticate
	connectionResult, err := sshManager.Connect(ctx, connectionRequest)
	if err != nil {
		return fmt.Errorf("authentication failed for %s: %w", machine.Hostname, err)
	}

	if !connectionResult.Success {
		return fmt.Errorf("authentication failed for %s: %s", machine.Hostname, connectionResult.Error)
	}

	fmt.Printf("✅ Authentication successful for %s\n", machine.Hostname)
	return nil
}

// getErrorMessage extracts error message from machine status
func getErrorMessage(status *spookytypes.MachineStatus) string {
	if status.Error != "" {
		return status.Error
	}
	return "unknown error"
}

func init() {
	// Add format and verbose flags to ping command
	machinesPingCmd.Flags().String("format", "text", "Output format: text or json")
	machinesPingCmd.Flags().Bool("verbose", false, "Show detailed output for all machines")
	machinesPingCmd.Flags().Bool("auth", false, "Test authentication in addition to connectivity")

	// Add flags to export command
	machinesExportCmd.Flags().String("output", "", "Output file path (required)")
	machinesExportCmd.Flags().String("machine", "", "Export specific machine by hostname")
	machinesExportCmd.Flags().StringArray("tags", []string{}, "Filter machines by tags (key=value or key-only)")

	// Add flags to machines encrypt command
	machinesEncryptCmd.Flags().Bool("dry-run", false, "Show what would be encrypted without making changes")

	// Add machines commands to root
	machinesCmd.AddCommand(machinesListCmd)
	machinesCmd.AddCommand(machinesValidateCmd)
	machinesCmd.AddCommand(machinesPingCmd)
	machinesCmd.AddCommand(machinesExportCmd)
	machinesCmd.AddCommand(machinesEncryptCmd)
	RootCmd.AddCommand(machinesCmd)
}
