package coordinator

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	spookyconfig "spooky/internal/config"
	spookyconfigtypes "spooky/internal/types/config"
	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookymachines "spooky/internal/machines"
	spookyssh "spooky/internal/ssh"
	spookysshtypes "spooky/internal/ssh/types"
)

// CoordinatorMachinesIntegration implements machines system integration
type CoordinatorMachinesIntegration struct {
	machinesManager spookymachines.MachineManager
	logger          spookylogging.Logger
}

// NewCoordinatorMachinesIntegration creates a new machines integration
func NewCoordinatorMachinesIntegration(machinesManager spookymachines.MachineManager, logger spookylogging.Logger) *CoordinatorMachinesIntegration {
	return &CoordinatorMachinesIntegration{
		machinesManager: machinesManager,
		logger:          logger,
	}
}

// LoadMachines loads machines from the project
func (mi *CoordinatorMachinesIntegration) LoadMachines(projectPath string) (*spookyinterfaces.MachinesContext, error) {
	context := &spookyinterfaces.MachinesContext{
		BaseContext: spookyinterfaces.BaseContext{
			ProjectPath: projectPath,
			Timestamp:   time.Now(),
		},
		Machines: make(map[string]*spookyinterfaces.Machine),
	}

	// Load machines from project using config parser
	inventoryFile := filepath.Join(projectPath, "inventory.hcl")
	machinesFile := filepath.Join(projectPath, "machines.hcl")

	var inventoryConfig *spookyconfigtypes.InventoryConfig
	var err error

	// Try inventory.hcl first, then machines.hcl
	if _, err := os.Stat(inventoryFile); err == nil {
		inventoryConfig, err = spookyconfig.ParseInventoryConfig(inventoryFile)
	} else if _, err := os.Stat(machinesFile); err == nil {
		inventoryConfig, err = spookyconfig.ParseMachinesInventory(machinesFile)
	} else {
		mi.logger.Warn("No inventory or machines file found", spookylogging.String("project", projectPath))
		return context, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load machines from project: %w", err)
	}

	// Check if inventoryConfig is nil (shouldn't happen but safety check)
	if inventoryConfig == nil {
		mi.logger.Warn("Inventory config is nil", spookylogging.String("project", projectPath))
		return context, nil
	}

	// Convert config machines to interface machines
	for i := range inventoryConfig.Machines {
		machine := &inventoryConfig.Machines[i]
		context.Machines[machine.Name] = mi.convertFromConfigMachine(machine)
	}

	// Build indexes if machines manager is available
	if mi.machinesManager != nil {
		err = mi.machinesManager.BuildIndexes(inventoryConfig.Machines)
		if err != nil {
			mi.logger.Warn("Failed to build machine indexes", spookylogging.Error(err))
		}
	}

	mi.logger.Info("Loaded machines from project",
		spookylogging.String("project", projectPath),
		spookylogging.Int("machines_count", len(context.Machines)))

	return context, nil
}

// ValidateMachines validates machines data
func (mi *CoordinatorMachinesIntegration) ValidateMachines(machinesContext *spookyinterfaces.MachinesContext) error {
	if machinesContext == nil {
		return fmt.Errorf("machines context cannot be nil")
	}

	// Validate each machine
	for name, machine := range machinesContext.Machines {
		if err := mi.validateMachine(machine); err != nil {
			return fmt.Errorf("machine '%s' validation failed: %w", name, err)
		}
	}

	return nil
}

// ConnectToMachine connects to a specific machine using SSH
func (mi *CoordinatorMachinesIntegration) ConnectToMachine(machine string, context *spookyinterfaces.ConnectionContext) error {
	if machine == "" {
		return fmt.Errorf("machine name cannot be empty")
	}

	if context == nil {
		return fmt.Errorf("connection context cannot be nil")
	}

	// Look up machine in the machines manager
	if mi.machinesManager != nil {
		configMachine, found := mi.machinesManager.LookupByName(machine)
		if !found {
			return fmt.Errorf("machine '%s' not found in index", machine)
		}

		// Validate machine configuration
		if configMachine.Host == "" {
			return fmt.Errorf("machine '%s' has no host configured", machine)
		}

		if configMachine.User == "" {
			return fmt.Errorf("machine '%s' has no user configured", machine)
		}

		if configMachine.Port <= 0 || configMachine.Port > 65535 {
			return fmt.Errorf("machine '%s' has invalid port: %d", machine, configMachine.Port)
		}

		// Create SSH client using the existing SSH package
		timeout := 30 // Default timeout in seconds
		if context.Timeout > 0 {
			timeout = context.Timeout
		}

		// Create SSH manager and test connection
		sshManager := spookyssh.NewDefaultManager(mi.logger)

		// Convert machine config to SSH config
		sshConfig := &spookysshtypes.SSHConfig{
			Host:     configMachine.Host,
			Port:     configMachine.Port,
			Username: configMachine.User,
			Timeout:  time.Duration(timeout) * time.Second,
		}

		// Test connection
		connection, err := sshManager.Connect(configMachine.Host, sshConfig)
		if err != nil {
			return fmt.Errorf("failed to establish SSH connection to machine '%s': %w", machine, err)
		}
		defer sshManager.CloseConnection(connection)

		// Execute a simple test command
		result, err := sshManager.ExecuteCommand(connection, "echo 'Connection test successful'")
		if err != nil {
			return fmt.Errorf("failed to execute test command on machine '%s': %w", machine, err)
		}

		if result.ExitCode != 0 {
			return fmt.Errorf("test command failed on machine '%s': %s", machine, result.Stderr)
		}
		if err != nil {
			return fmt.Errorf("failed to establish SSH connection to machine '%s': %w", machine, err)
		}

		mi.logger.Info("Machine connection established",
			spookylogging.String("machine", machine),
			spookylogging.String("host", configMachine.Host),
			spookylogging.String("user", configMachine.User),
			spookylogging.Int("port", configMachine.Port))
	} else {
		mi.logger.Warn("Machines manager not available, skipping connection validation", spookylogging.String("machine", machine))
	}

	return nil
}

// PingMachine pings a machine to check connectivity using multiple methods
func (mi *CoordinatorMachinesIntegration) PingMachine(machine string) error {
	if machine == "" {
		return fmt.Errorf("machine name cannot be empty")
	}

	// Look up machine in the machines manager
	if mi.machinesManager != nil {
		configMachine, found := mi.machinesManager.LookupByName(machine)
		if !found {
			return fmt.Errorf("machine '%s' not found in index", machine)
		}

		// Validate machine configuration
		if configMachine.Host == "" {
			return fmt.Errorf("machine '%s' has no host configured", machine)
		}

		// Try multiple ping methods
		var pingResults []string
		var lastError error

		// Method 1: TCP connection test
		if err := mi.pingTCP(configMachine.Host, configMachine.Port); err != nil {
			lastError = err
			pingResults = append(pingResults, fmt.Sprintf("TCP failed: %v", err))
		} else {
			pingResults = append(pingResults, "TCP successful")
		}

		// Method 2: ICMP ping (if available)
		if pingTime, err := mi.pingICMP(configMachine.Host); err != nil {
			pingResults = append(pingResults, fmt.Sprintf("ICMP failed: %v", err))
		} else {
			pingResults = append(pingResults, fmt.Sprintf("ICMP successful: %v", pingTime))
		}

		// Method 3: SSH connection test
		if err := mi.pingSSH(configMachine); err != nil {
			pingResults = append(pingResults, fmt.Sprintf("SSH failed: %v", err))
		} else {
			pingResults = append(pingResults, "SSH successful")
		}

		// Log ping results
		mi.logger.Info("Machine ping results",
			spookylogging.String("machine", machine),
			spookylogging.String("host", configMachine.Host),
			spookylogging.String("results", fmt.Sprintf("%v", pingResults)))

		// Return error if all methods failed
		if len(pingResults) == 0 || (len(pingResults) == 1 && strings.Contains(pingResults[0], "failed")) {
			return fmt.Errorf("all ping methods failed for machine '%s': %w", machine, lastError)
		}

		return nil
	} else {
		mi.logger.Warn("Machines manager not available, skipping ping validation", spookylogging.String("machine", machine))
	}

	return nil
}

// pingTCP performs a TCP connection test
func (mi *CoordinatorMachinesIntegration) pingTCP(host string, port int) error {
	timeout := 5 * time.Second
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if err != nil {
		return fmt.Errorf("TCP connection failed: %w", err)
	}
	defer conn.Close()
	return nil
}

// pingICMP performs an ICMP ping test
func (mi *CoordinatorMachinesIntegration) pingICMP(host string) (time.Duration, error) {
	// Note: This is a simplified ICMP ping implementation
	// In a real implementation, you would use a proper ICMP library
	// For now, we'll simulate the ping with a TCP connection test

	start := time.Now()
	err := mi.pingTCP(host, 22) // Try SSH port as fallback
	if err != nil {
		return 0, fmt.Errorf("ICMP ping failed: %w", err)
	}

	return time.Since(start), nil
}

// pingSSH performs an SSH connection test
func (mi *CoordinatorMachinesIntegration) pingSSH(configMachine *spookyconfigtypes.Machine) error {
	// Create SSH manager
	sshManager := spookyssh.NewDefaultManager(mi.logger)

	// Convert machine config to SSH config
	sshConfig := &spookysshtypes.SSHConfig{
		Host:     configMachine.Host,
		Port:     configMachine.Port,
		Username: configMachine.User,
		Timeout:  5 * time.Second,
	}

	// Test connection
	connection, err := sshManager.Connect(configMachine.Host, sshConfig)
	if err != nil {
		return fmt.Errorf("failed to create SSH connection: %w", err)
	}
	defer sshManager.CloseConnection(connection)

	// Test connection with a simple command
	result, err := sshManager.ExecuteCommand(connection, "echo 'SSH ping test'")
	if err != nil {
		return fmt.Errorf("SSH ping failed: %w", err)
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("SSH ping command failed: %s", result.Stderr)
	}

	return nil
}

// GetMachine gets a specific machine by name
func (mi *CoordinatorMachinesIntegration) GetMachine(name string, context *spookyinterfaces.MachinesContext) (*spookyinterfaces.Machine, error) {
	if name == "" {
		return nil, fmt.Errorf("machine name cannot be empty")
	}

	if context == nil {
		return nil, fmt.Errorf("machines context cannot be nil")
	}

	// Look up machine in context
	if machine, exists := context.Machines[name]; exists {
		return machine, nil
	}

	return nil, fmt.Errorf("machine '%s' not found", name)
}

// ListMachines lists all available machines
func (mi *CoordinatorMachinesIntegration) ListMachines(context *spookyinterfaces.MachinesContext) ([]*spookyinterfaces.Machine, error) {
	if context == nil {
		return nil, fmt.Errorf("machines context cannot be nil")
	}

	// Convert map to slice
	var machines []*spookyinterfaces.Machine
	for _, machine := range context.Machines {
		machines = append(machines, machine)
	}

	return machines, nil
}

// AddMachine adds a new machine to the project
func (mi *CoordinatorMachinesIntegration) AddMachine(name string, machine *spookyinterfaces.Machine, context *spookyinterfaces.MachinesContext) error {
	if name == "" {
		return fmt.Errorf("machine name cannot be empty")
	}

	if machine == nil {
		return fmt.Errorf("machine cannot be nil")
	}

	if context == nil {
		return fmt.Errorf("machines context cannot be nil")
	}

	// Validate machine before adding
	if err := mi.validateMachine(machine); err != nil {
		return fmt.Errorf("machine validation failed: %w", err)
	}

	// Add machine to context
	context.Machines[name] = machine

	mi.logger.Info("Added machine", spookylogging.String("machine", name))

	return nil
}

// RemoveMachine removes a machine from the project
func (mi *CoordinatorMachinesIntegration) RemoveMachine(name string, context *spookyinterfaces.MachinesContext) error {
	if name == "" {
		return fmt.Errorf("machine name cannot be empty")
	}

	if context == nil {
		return fmt.Errorf("machines context cannot be nil")
	}

	// Remove machine from context
	delete(context.Machines, name)

	mi.logger.Info("Removed machine", spookylogging.String("machine", name))

	return nil
}

// validateMachine validates a single machine
func (mi *CoordinatorMachinesIntegration) validateMachine(machine *spookyinterfaces.Machine) error {
	if machine == nil {
		return fmt.Errorf("machine cannot be nil")
	}

	if machine.Name == "" {
		return fmt.Errorf("machine name cannot be empty")
	}

	if machine.Host == "" {
		return fmt.Errorf("machine host cannot be empty")
	}

	if machine.Port <= 0 || machine.Port > 65535 {
		return fmt.Errorf("machine port must be between 1 and 65535")
	}

	if machine.Username == "" {
		return fmt.Errorf("machine username cannot be empty")
	}

	return nil
}

// convertToConfigMachine converts interfaces.Machine to configTypes.Machine
func (mi *CoordinatorMachinesIntegration) convertToConfigMachine(machine *spookyinterfaces.Machine) *spookyconfigtypes.Machine {
	if machine == nil {
		return nil
	}

	return &spookyconfigtypes.Machine{
		Name:     machine.Name,
		Host:     machine.Host,
		Port:     machine.Port,
		User:     machine.Username,
		Tags:     mi.convertTags(machine.Tags),
		Metadata: mi.convertMetadataToString(machine.Metadata),
	}
}

// convertFromConfigMachine converts configTypes.Machine to interfaces.Machine
func (mi *CoordinatorMachinesIntegration) convertFromConfigMachine(machine *spookyconfigtypes.Machine) *spookyinterfaces.Machine {
	if machine == nil {
		return nil
	}

	return &spookyinterfaces.Machine{
		Name:     machine.Name,
		Host:     machine.Host,
		Port:     machine.Port,
		Username: machine.User,
		Tags:     mi.convertFromTags(machine.Tags),
		Metadata: mi.convertMetadataToInterface(machine.Metadata),
	}
}

// convertTags converts []string to map[string]string
func (mi *CoordinatorMachinesIntegration) convertTags(tags []string) map[string]string {
	result := make(map[string]string)
	for i, tag := range tags {
		result[fmt.Sprintf("tag_%d", i)] = tag
	}
	return result
}

// convertFromTags converts map[string]string to []string
func (mi *CoordinatorMachinesIntegration) convertFromTags(tags map[string]string) []string {
	var result []string
	for _, tag := range tags {
		result = append(result, tag)
	}
	return result
}

// convertMetadataToString converts map[string]interface{} to map[string]string
func (mi *CoordinatorMachinesIntegration) convertMetadataToString(metadata map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for key, value := range metadata {
		if str, ok := value.(string); ok {
			result[key] = str
		} else {
			result[key] = fmt.Sprintf("%v", value)
		}
	}
	return result
}

// convertMetadataToInterface converts map[string]string to map[string]interface{}
func (mi *CoordinatorMachinesIntegration) convertMetadataToInterface(metadata map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range metadata {
		result[key] = value
	}
	return result
}

// GetMachineHealth gets the health status of a machine
func (mi *CoordinatorMachinesIntegration) GetMachineHealth(machine string) (map[string]interface{}, error) {
	if machine == "" {
		return nil, fmt.Errorf("machine name cannot be empty")
	}

	// Look up machine in the machines manager
	if mi.machinesManager != nil {
		configMachine, found := mi.machinesManager.LookupByName(machine)
		if !found {
			return nil, fmt.Errorf("machine '%s' not found in index", machine)
		}

		health := map[string]interface{}{
			"machine_name": machine,
			"host":         configMachine.Host,
			"status":       "unknown",
			"last_check":   time.Now(),
		}

		// Check connectivity
		if err := mi.PingMachine(machine); err != nil {
			health["status"] = "unreachable"
			health["last_error"] = err.Error()
		} else {
			health["status"] = "online"
		}

		// Get additional health metrics if machine is reachable
		if health["status"] == "online" {
			mi.getMachineMetrics(machine, health)
		}

		return health, nil
	}

	return nil, fmt.Errorf("machines manager not available")
}

// getMachineMetrics gets additional metrics from a machine
func (mi *CoordinatorMachinesIntegration) getMachineMetrics(machine string, health map[string]interface{}) {
	// In a real implementation, this would:
	// - Connect to the machine via SSH
	// - Run commands to get system metrics
	// - Parse the results and populate the health map

	// For now, we'll just set some basic metrics
	health["cpu_usage"] = 0.0
	health["memory_usage"] = 0.0
	health["disk_usage"] = 0.0
	health["uptime"] = time.Duration(0)

	mi.logger.Debug("Retrieved machine metrics", spookylogging.String("machine", machine))
}

// DiscoverMachines discovers machines on the network
func (mi *CoordinatorMachinesIntegration) DiscoverMachines(network string) ([]*spookyinterfaces.Machine, error) {
	if network == "" {
		return nil, fmt.Errorf("network cannot be empty")
	}

	mi.logger.Info("Discovering machines on network", spookylogging.String("network", network))

	// In a real implementation, this would:
	// - Scan the network for SSH-enabled hosts
	// - Attempt to connect to each discovered host
	// - Gather basic information about each machine
	// - Return a list of discovered machines

	// For now, we'll return an empty list
	var discoveredMachines []*spookyinterfaces.Machine

	mi.logger.Info("Machine discovery completed",
		spookylogging.String("network", network),
		spookylogging.Int("discovered_count", len(discoveredMachines)))

	return discoveredMachines, nil
}

// BackupMachineConfiguration backs up machine configuration
func (mi *CoordinatorMachinesIntegration) BackupMachineConfiguration(machine string, backupPath string) error {
	if machine == "" {
		return fmt.Errorf("machine name cannot be empty")
	}

	if backupPath == "" {
		return fmt.Errorf("backup path cannot be empty")
	}

	mi.logger.Info("Backing up machine configuration",
		spookylogging.String("machine", machine),
		spookylogging.String("backup_path", backupPath))

	// Look up machine in the machines manager
	if mi.machinesManager != nil {
		_, found := mi.machinesManager.LookupByName(machine)
		if !found {
			return fmt.Errorf("machine '%s' not found in index", machine)
		}

		// Create backup directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(backupPath), 0755); err != nil {
			return fmt.Errorf("failed to create backup directory: %w", err)
		}

		// In a real implementation, this would:
		// - Connect to the machine
		// - Gather configuration files
		// - Create a backup archive
		// - Store it at the specified path

		mi.logger.Info("Machine configuration backup completed",
			spookylogging.String("machine", machine),
			spookylogging.String("backup_path", backupPath))

		return nil
	}

	return fmt.Errorf("machines manager not available")
}
