package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"spooky/internal/facts"
	"spooky/internal/schemas"
	"spooky/internal/ssh"

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
		fmt.Println("DEBUG: gatherFactsCmd started")
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
		fmt.Println("DEBUG: About to call getMachinesFromConfig")
		machines, err := getMachinesFromConfig()
		fmt.Printf("DEBUG: getMachinesFromConfig returned %d machines, err: %v\n", len(machines), err)
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
	fmt.Println("DEBUG: loadProjectConfig called")
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

	// Use the validator to parse and validate the project configuration
	validator := schemas.NewValidator()
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
	fmt.Println("DEBUG: loadSSHConfig called")
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
	fmt.Println("DEBUG: getMachinesFromConfig called")

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
				Type:       "machines",
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

	fmt.Printf("DEBUG: Found %d machine blocks\n", len(machineContent.Blocks))

	// Create machines slice
	var machines []*schemas.MachinesMachineV1

	// Process each machine block
	for _, machineBlock := range machineContent.Blocks {
		machineName := machineBlock.Labels[0]
		fmt.Printf("DEBUG: Processing machine: %s\n", machineName)

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
			authMethod := authBlock.Labels[0]

			if authMethod == "password" {
				var passwordData struct {
					Password struct {
						Value     string `hcl:"value,attr"`
						Encrypted bool   `hcl:"encrypted,attr"`
					} `hcl:"password,block"`
				}

				if diags := gohcl.DecodeBody(authBlock.Body, nil, &passwordData); diags.HasErrors() {
					return nil, fmt.Errorf("failed to decode password auth for machine %s: %v", machineName, diags)
				}

				machine.Authentication.Password = schemas.MachinesMachineAuthenticationPasswordV1{
					Value:     passwordData.Password.Value,
					Encrypted: passwordData.Password.Encrypted,
				}
			} else if authMethod == "publickey" {
				var publicKeyData struct {
					PublicKeyPath string `hcl:"public_key_path,attr"`
					Passphrase    struct {
						Value     string `hcl:"value,attr"`
						Encrypted bool   `hcl:"encrypted,attr"`
					} `hcl:"passphrase,block"`
				}

				if diags := gohcl.DecodeBody(authBlock.Body, nil, &publicKeyData); diags.HasErrors() {
					return nil, fmt.Errorf("failed to decode publickey auth for machine %s: %v", machineName, diags)
				}

				machine.Authentication.PublicKeyPath = publicKeyData.PublicKeyPath
				machine.Authentication.Passphrase = schemas.MachinesMachineAuthenticationPassphraseV1{
					Value:     publicKeyData.Passphrase.Value,
					Encrypted: publicKeyData.Passphrase.Encrypted,
				}
			}
		}

		// Validate required fields
		if machine.Hostname == "" {
			return nil, fmt.Errorf("machine %s missing required hostname", machineName)
		}
		if machine.User == "" {
			return nil, fmt.Errorf("machine %s missing required user", machineName)
		}

		fmt.Printf("DEBUG: Successfully parsed machine: %s (%s@%s:%d)\n", machineName, machine.User, machine.Hostname, machine.Port)
		machines = append(machines, machine)
	}

	fmt.Printf("DEBUG: getMachinesFromConfig returning %d machines\n", len(machines))
	return machines, nil
}

// parseMachineBlock parses a single machine block into a MachinesMachineV1
func parseMachineBlock(block map[string]interface{}) (*schemas.MachinesMachineV1, error) {
	machine := &schemas.MachinesMachineV1{
		Port:              22, // Default SSH port
		ConnectionTimeout: 30, // Default timeout
		MaxRetries:        3,  // Default retries
		RetryDelay:        5,  // Default delay
	}

	// Extract name (from block label)
	if name, exists := block["name"]; exists {
		if _, ok := name.(string); ok {
			// The name is typically in the block label, not as a field
			// But we'll handle it if it's present
		}
	}

	// Extract hostname
	if hostname, exists := block["hostname"]; exists {
		if hostnameStr, ok := hostname.(string); ok {
			machine.Hostname = hostnameStr
		}
	}

	// Extract port
	if port, exists := block["port"]; exists {
		if portNum, ok := port.(int); ok {
			machine.Port = portNum
		}
	}

	// Extract user
	if user, exists := block["user"]; exists {
		if userStr, ok := user.(string); ok {
			machine.User = userStr
		}
	}

	// Extract description
	if desc, exists := block["description"]; exists {
		if descStr, ok := desc.(string); ok {
			machine.Description = descStr
		}
	}

	// Extract connection timeout
	if timeout, exists := block["connection_timeout"]; exists {
		if timeoutNum, ok := timeout.(int); ok {
			machine.ConnectionTimeout = timeoutNum
		}
	}

	// Extract max retries
	if retries, exists := block["max_retries"]; exists {
		if retriesNum, ok := retries.(int); ok {
			machine.MaxRetries = retriesNum
		}
	}

	// Extract retry delay
	if delay, exists := block["retry_delay"]; exists {
		if delayNum, ok := delay.(int); ok {
			machine.RetryDelay = delayNum
		}
	}

	// Parse authentication block
	if auth, exists := block["authentication"]; exists {
		if authMap, ok := auth.(map[string]interface{}); ok {
			authConfig, err := parseAuthenticationBlock(authMap)
			if err != nil {
				return nil, fmt.Errorf("failed to parse authentication: %w", err)
			}
			machine.Authentication = *authConfig
		}
	}

	// Validate required fields
	if machine.Hostname == "" {
		return nil, fmt.Errorf("machine hostname is required")
	}
	if machine.User == "" {
		return nil, fmt.Errorf("machine user is required")
	}

	return machine, nil
}

// parseAuthenticationBlock parses an authentication block
func parseAuthenticationBlock(block map[string]interface{}) (*schemas.MachinesMachineAuthenticationV1, error) {
	auth := &schemas.MachinesMachineAuthenticationV1{}

	// Look for password authentication
	if password, exists := block["password"]; exists {
		if passwordMap, ok := password.(map[string]interface{}); ok {
			auth.Password = schemas.MachinesMachineAuthenticationPasswordV1{}

			if value, exists := passwordMap["value"]; exists {
				if valueStr, ok := value.(string); ok {
					auth.Password.Value = valueStr
				}
			}

			if encrypted, exists := passwordMap["encrypted"]; exists {
				if encryptedBool, ok := encrypted.(bool); ok {
					auth.Password.Encrypted = encryptedBool
				}
			}
		}
	}

	// Look for public key path
	if publicKeyPath, exists := block["public_key_path"]; exists {
		if pathStr, ok := publicKeyPath.(string); ok {
			auth.PublicKeyPath = pathStr
		}
	}

	// Look for passphrase
	if passphrase, exists := block["passphrase"]; exists {
		if passphraseMap, ok := passphrase.(map[string]interface{}); ok {
			auth.Passphrase = schemas.MachinesMachineAuthenticationPassphraseV1{}

			if value, exists := passphraseMap["value"]; exists {
				if valueStr, ok := value.(string); ok {
					auth.Passphrase.Value = valueStr
				}
			}

			if encrypted, exists := passphraseMap["encrypted"]; exists {
				if encryptedBool, ok := encrypted.(bool); ok {
					auth.Passphrase.Encrypted = encryptedBool
				}
			}
		}
	}

	return auth, nil
}

// getMapKeys returns the keys of a map as a slice of strings
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
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
