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
	RunE: func(cmd *cobra.Command, args []string) error {
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
	RunE: func(cmd *cobra.Command, args []string) error {
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

// handleMachinesList handles listing machines using the MachinesIntegration interface
func handleMachinesList(projectPath string) error {
	ctx := context.Background()

	// Initialize dependencies if not already done
	if machinesManager == nil {
		if err := InitializeMachinesDependencies(); err != nil {
			return fmt.Errorf("failed to initialize machines dependencies: %w", err)
		}
	}

	fmt.Printf("🔍 Loading machines from project: %s\n", projectPath)

	// Load machines using the enhanced manager (supports both machines.hcl and machines/ directory)
	machines, err := machinesManager.LoadMachines(ctx, projectPath)
	if err != nil {
		return fmt.Errorf("failed to load machines: %w", err)
	}

	fmt.Printf("📊 Found %d machines:\n\n", len(machines))

	if len(machines) == 0 {
		fmt.Printf("No machines found in inventory.\n")
		return nil
	}

	// Group machines by source file for better display
	machinesBySource := make(map[string][]spookytypes.Machine)
	for _, machine := range machines {
		sourceFile := "unknown"
		if machine.MachineMetadata != nil && machine.MachineMetadata.CustomFields != nil {
			if src, exists := machine.MachineMetadata.CustomFields["source_file"]; exists {
				sourceFile = src
			}
		}
		machinesBySource[sourceFile] = append(machinesBySource[sourceFile], machine)
	}

	// Display machines grouped by source
	for sourceFile, sourceMachines := range machinesBySource {
		fmt.Printf("📁 Source: %s (%d machines)\n", sourceFile, len(sourceMachines))
		fmt.Printf("%s\n", strings.Repeat("─", 50))

		for i, machine := range sourceMachines {
			fmt.Printf("%d. %s (%s)\n", i+1, machine.Hostname, machine.Host)
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
	}

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
		for i, err := range result.Errors {
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
		for i, warning := range result.Warnings {
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
	} else {
		fmt.Printf("✅ Validation passed with %d warnings\n", len(result.Warnings))
	}

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
			if statuses[i].Status == "online" {
				// Test authentication for this machine
				err := testMachineAuthentication(ctx, machines[i])
				if err != nil {
					statuses[i].Status = "auth_failed"
					statuses[i].Error = fmt.Sprintf("Authentication failed: %v", err)
				} else {
					statuses[i].Status = "authenticated"
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

// handleMachinesExport handles exporting machines using the MachinesIntegration interface
func handleMachinesExport(cmd *cobra.Command, projectPath string) error {
	ctx := context.Background()

	// Initialize dependencies if not already done
	if machinesManager == nil {
		if err := InitializeMachinesDependencies(); err != nil {
			return fmt.Errorf("failed to initialize machine dependencies: %w", err)
		}
	}

	// Get output path from flags
	outputPath, err := cmd.Flags().GetString("output")
	if err != nil {
		return fmt.Errorf("failed to get output flag: %w", err)
	}
	if outputPath == "" {
		return fmt.Errorf("--output flag is required")
	}

	// Load machines from project
	machines, err := machinesManager.LoadMachines(ctx, projectPath)
	if err != nil {
		return fmt.Errorf("failed to load machines from project: %w", err)
	}

	// Apply filters if specified
	filteredMachines := machines

	// Filter by machine name if specified
	if machineName, _ := cmd.Flags().GetString("machine"); machineName != "" {
		var filtered []spookytypes.Machine
		for _, machine := range machines {
			if machine.Hostname == machineName {
				filtered = append(filtered, machine)
				break
			}
		}
		if len(filtered) == 0 {
			return fmt.Errorf("machine '%s' not found in project", machineName)
		}
		filteredMachines = filtered
	}

	// Filter by tags if specified
	if tags, _ := cmd.Flags().GetStringArray("tags"); len(tags) > 0 {
		taggedMachines, err := machinesManager.GetMachinesByTags(ctx, tags)
		if err != nil {
			return fmt.Errorf("failed to filter machines by tags: %w", err)
		}
		// Intersect with already filtered machines
		var intersection []spookytypes.Machine
		for _, tagged := range taggedMachines {
			for _, filtered := range filteredMachines {
				if tagged.Hostname == filtered.Hostname {
					intersection = append(intersection, tagged)
					break
				}
			}
		}
		filteredMachines = intersection
	}

	// Export machines to HCL
	if err := machinesManager.ExportMachines(ctx, filteredMachines, outputPath); err != nil {
		return fmt.Errorf("failed to export machines: %w", err)
	}

	// Output success message
	fmt.Printf("✅ Successfully exported %d machines to: %s\n", len(filteredMachines), outputPath)

	return nil
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
		case "online":
			online++
			if verbose {
				fmt.Printf("%s: %s (latency: %dms)\n",
					status.Machine.Hostname, status.Status, status.Latency)
			} else {
				fmt.Printf("%s: %s\n", status.Machine.Hostname, status.Status)
			}
		case "offline":
			offline++
			if verbose {
				fmt.Printf("%s: %s (latency: %dms; Error: %s)\n",
					status.Machine.Hostname, status.Status, status.Latency,
					getErrorMessage(status))
			} else {
				fmt.Printf("%s: %s\n", status.Machine.Hostname, status.Status)
			}
		case "error":
			errors++
			if verbose {
				fmt.Printf("%s: %s (latency: %dms; Error: %s)\n",
					status.Machine.Hostname, status.Status, status.Latency,
					getErrorMessage(status))
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
		if status.Status != "online" || verbose {
			machineJSON["latency_ms"] = status.Latency
			if status.Status != "online" {
				machineJSON["error"] = getErrorMessage(status)
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
func testMachineAuthentication(ctx context.Context, machine spookytypes.Machine) error {
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
	if machine.KeyFile != "" {
		clientConfig.DefaultKeyPath = machine.KeyFile
		clientConfig.DefaultAuthMethod = "public_key"
	} else if machine.Password != "" {
		clientConfig.DefaultAuthMethod = "password"
	} else {
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
func getErrorMessage(status spookytypes.MachineStatus) string {
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

	// Add machines commands to root
	machinesCmd.AddCommand(machinesListCmd)
	machinesCmd.AddCommand(machinesValidateCmd)
	machinesCmd.AddCommand(machinesPingCmd)
	machinesCmd.AddCommand(machinesExportCmd)
	RootCmd.AddCommand(machinesCmd)
}
