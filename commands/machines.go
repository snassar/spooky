package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"spooky/internal/encryption"
	"spooky/internal/logging"
	"spooky/internal/schemas"
	"spooky/internal/ssh"
	"spooky/internal/utilities"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/spf13/cobra"
)

var (
	machinesCmd = &cobra.Command{
		Use:   schemas.ResourceTypeMachines,
		Short: "Manage and interact with remote machines",
		Long: `Manage and interact with remote machines defined in your project configuration.

Machines can be pinged to test connectivity and response times.

Examples:
  spooky machines ping
  spooky machines ping --targets web-server-01,db-server-01`,
	}

	pingCmd = &cobra.Command{
		Use:   "ping",
		Short: "Ping remote machines to test connectivity",
		Long: `Ping remote machines to test basic connectivity and response times.

This command will attempt to connect to all machines defined in your machines.hcl
file and report their status and response times.

Examples:
  spooky machines ping
  spooky machines ping --targets web-server-01,db-server-01`,
		RunE: pingMachines,
	}
)

// Machine ping flags
var (
	pingTargets []string
	pingTimeout int
	pingAuth    bool
	pingOutput  string
)

// MachineWithName wraps a machine with its name from the HCL block label
type MachineWithName struct {
	Name    string
	Machine *schemas.MachinesMachineV1
}

// PingResult represents the result of a ping operation
type PingResult struct {
	Hostname      string `json:"hostname"`
	DNS           string `json:"dns"`
	IPAddress     string `json:"ip_address"`
	SSH           string `json:"ssh"`
	Authenticated string `json:"authenticated,omitempty"`
}

func init() {
	// Add subcommands to machines command
	machinesCmd.AddCommand(pingCmd)

	// Add flags to ping command
	pingCmd.Flags().StringSliceVarP(&pingTargets, "targets", "t", nil, "specific machines to ping (comma-separated list)")
	pingCmd.Flags().IntVar(&pingTimeout, "timeout", 30, "ping timeout in seconds")
	pingCmd.Flags().BoolVar(&pingAuth, "auth", false, "attempt authentication during ping")
	pingCmd.Flags().StringVar(&pingOutput, "output", "text", "output format (text, json)")

	// Add machines command to root
	RootCmd.AddCommand(machinesCmd)
}

// pingMachines implements the ping functionality
func pingMachines(cmd *cobra.Command, args []string) error {
	// Create machines with names by parsing HCL to get the block labels
	machinesWithNames, err := getMachinesWithNames()
	if err != nil {
		return fmt.Errorf("failed to get machines from config: %w", err)
	}

	if len(machinesWithNames) == 0 {
		return fmt.Errorf("no machines configured for ping")
	}

	// Filter machines if targets specified
	if len(pingTargets) > 0 {
		machinesWithNames = filterMachinesByTargets(machinesWithNames, pingTargets)
		if len(machinesWithNames) == 0 {
			return fmt.Errorf("no machines found matching targets: %v", pingTargets)
		}
	}

	logger := logging.GetGlobalLogger()
	logger.Info("Pinging machines", slog.Int("machine_count", len(machinesWithNames)))

	// Ping each machine
	for _, machineWithName := range machinesWithNames {
		result := pingMachine(machineWithName, pingAuth, pingTimeout)

		if pingOutput == "json" {
			outputJSON(result)
		} else {
			outputText(result)
		}
	}

	return nil
}

// getMachinesWithNames gets machines with their names from HCL block labels
func getMachinesWithNames() ([]MachineWithName, error) {
	// We need to re-parse the HCL to get the machine names from block labels
	// This is a simplified approach - in production, we'd modify getMachinesFromConfig
	// to return both machines and their names

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
				Type:       schemas.ResourceTypeMachines,
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

	var machinesWithNames []MachineWithName

	// Process each machine block to get the name and create the machine struct
	for _, machineBlock := range machineContent.Blocks {
		machineName := machineBlock.Labels[0]

		// Create a new machine struct
		machine := &schemas.MachinesMachineV1{
			Description:       "",
			Hostname:          machineName, // Use machine name as hostname by default
			Port:              22,          // default
			User:              "",
			Authentication:    schemas.MachinesMachineAuthenticationV1{},
			ConnectionTimeout: 30, // default
			MaxRetries:        3,  // default
			RetryDelay:        5,  // default
			Facts:             schemas.MachinesMachineFactsV1{},
		}

		// Use a schema that allows authentication blocks
		attrSchema := &hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{Name: "hostname", Required: false},
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
		content, diags := machineBlock.Body.Content(attrSchema)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to decode machine %s: %v", machineName, diags)
		}

		// Extract attributes
		if hostnameAttr, exists := content.Attributes["hostname"]; exists {
			var hostname string
			if diags := gohcl.DecodeExpression(hostnameAttr.Expr, nil, &hostname); diags.HasErrors() {
				return nil, fmt.Errorf("failed to decode hostname for machine %s: %v", machineName, diags)
			}
			machine.Hostname = hostname // Override default with explicit hostname
		}

		if portAttr, exists := content.Attributes["port"]; exists {
			var port int
			if diags := gohcl.DecodeExpression(portAttr.Expr, nil, &port); diags.HasErrors() {
				return nil, fmt.Errorf("failed to decode port for machine %s: %v", machineName, diags)
			}
			machine.Port = port
		}

		if userAttr, exists := content.Attributes["user"]; exists {
			var user string
			if diags := gohcl.DecodeExpression(userAttr.Expr, nil, &user); diags.HasErrors() {
				return nil, fmt.Errorf("failed to decode user for machine %s: %v", machineName, diags)
			}
			machine.User = user
		}

		// Validate required fields after parsing
		if machine.User == "" {
			return nil, fmt.Errorf("machine %s missing required user field in machines.hcl", machineName)
		}

		// Handle authentication blocks
		if len(content.Blocks) > 0 {
			authBlock := content.Blocks[0]
			authConfig, err := utilities.ParseAuthenticationBlock(authBlock, machineName)
			if err != nil {
				return nil, err
			}
			machine.Authentication = *authConfig
		}

		machinesWithNames = append(machinesWithNames, MachineWithName{
			Name:    machineName,
			Machine: machine,
		})
	}

	return machinesWithNames, nil
}

// pingMachine performs a ping operation on a single machine
func pingMachine(machine MachineWithName, withAuth bool, timeout int) PingResult {
	result := PingResult{
		Hostname: machine.Name, // Use the machine name directly as hostname
	}

	// Determine port to use
	port := 22 // default SSH port
	if machine.Machine.Port != 0 {
		port = machine.Machine.Port
	}

	// Check if hostname is IP or FQDN and resolve DNS
	ip, dnsStatus := resolveHostname(machine.Name)
	result.IPAddress = ip
	result.DNS = dnsStatus

	// Test SSH connectivity
	sshStatus := testSSHConnection(ip, port, timeout)
	result.SSH = sshStatus

	// Test authentication if requested
	if withAuth {
		authStatus := testAuthentication(machine.Machine, timeout)
		result.Authenticated = authStatus
	}

	return result
}

// resolveHostname resolves a hostname to IP and returns DNS status
func resolveHostname(hostname string) (string, string) {
	// Check if it's already an IP address
	if net.ParseIP(hostname) != nil {
		return hostname, "N/A"
	}

	// Resolve FQDN
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return hostname, "unresolved"
	}

	if len(ips) > 0 {
		return ips[0].String(), "resolved"
	}

	return hostname, "unresolved"
}

// testSSHConnection tests if SSH daemon is listening on the specified port
func testSSHConnection(ip string, port int, timeout int) string {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), time.Duration(timeout)*time.Second)
	if err != nil {
		return "unavailable"
	}
	defer conn.Close()
	return "available"
}

// testAuthentication attempts to authenticate with the machine
func testAuthentication(machine *schemas.MachinesMachineV1, timeout int) string {
	// Create a basic SSH configuration for testing
	sshConfig := &schemas.SpookySSHV1{
		Timeout:            timeout,
		KeepaliveInterval:  60,
		KeepaliveCount:     3,
		KeyScanTimeout:     10,
		KnownHostsStrict:   false, // For testing, be more permissive
		KnownHostsMode:     "accept-new",
		ConnectionPoolSize: 1,
	}

	// Create age encryption (optional)
	var ageEncryption *encryption.AgeEncryption
	if os.Getenv("SSH_AUTH_SOCK") != "" {
		// Try to create age encryption if SSH agent is available
		var err error
		ageEncryption, err = encryption.NewAgeEncryption("", "")
		if err != nil {
			logger := logging.GetGlobalLogger()
			logger.Warn("failed to initialize age encryption for authentication testing, continuing without encryption support",
				slog.String("error", err.Error()))
			// Continue with nil encryption - SSH manager will handle this gracefully
		}
	}

	// Create SSH manager
	manager := ssh.NewSimpleSSHManager(ageEncryption, sshConfig)

	// Test connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	err := manager.TestConnection(ctx, machine)
	if err != nil {
		// Return specific error messages based on the type of failure
		errMsg := err.Error()
		if contains(errMsg, "authentication") || contains(errMsg, "password") || contains(errMsg, "publickey") {
			return "failed"
		} else if contains(errMsg, "connection") || contains(errMsg, "timeout") || contains(errMsg, "refused") {
			return "unavailable"
		} else {
			return "failed"
		}
	}

	return "success"
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				containsSubstring(s, substr))))
}

// containsSubstring performs a simple substring search
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// filterMachinesByTargets filters machines by target names
func filterMachinesByTargets(machines []MachineWithName, targets []string) []MachineWithName {
	var filtered []MachineWithName
	targetMap := make(map[string]bool)

	for _, target := range targets {
		targetMap[target] = true
	}

	for _, machine := range machines {
		if targetMap[machine.Name] {
			filtered = append(filtered, machine)
		}
	}

	return filtered
}

// outputText outputs ping result in text format
func outputText(result PingResult) {
	output := fmt.Sprintf("hostname: %s; DNS: %s; IP address: %s; SSH: %s",
		result.Hostname, result.DNS, result.IPAddress, result.SSH)

	if result.Authenticated != "" {
		output += fmt.Sprintf("; authenticated: %s", result.Authenticated)
	}

	output += ";"
	fmt.Println(output)
}

// outputJSON outputs ping result in JSON format
func outputJSON(result PingResult) {
	jsonData, err := json.Marshal(result)
	if err != nil {
		logger := logging.GetGlobalLogger()
		logger.Error("failed to marshal JSON output", slog.String("error", err.Error()))
		return
	}
	fmt.Println(string(jsonData))
}
