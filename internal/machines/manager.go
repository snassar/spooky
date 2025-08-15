// Package machines provides machine inventory management for the spooky codebase.
package machines

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"

	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookyssh "spooky/internal/ssh"
	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesmachines "spooky/internal/types/machines"
)

// Manager provides machine management functionality
type Manager struct {
	logger    spookytypeslogging.Logger
	loader    spookyinterfaces.MachineLoader
	validator spookyinterfaces.MachineValidator
}

// NewManager creates a new machine manager
func NewManager(
	logger spookytypeslogging.Logger,
	loader spookyinterfaces.MachineLoader,
	validator spookyinterfaces.MachineValidator,
) spookyinterfaces.MachinesIntegration {
	return &Manager{
		logger:    logger,
		loader:    loader,
		validator: validator,
	}
}

// LoadMachines loads machines from the given project path
// Supports both machines.hcl file and machines/ directory
func (m *Manager) LoadMachines(ctx context.Context, projectPath string) ([]spookytypes.Machine, error) {
	m.logger.Debug("Loading machines from project", map[string]interface{}{
		"project_path": projectPath,
	})

	var allMachines []spookytypes.Machine
	var loadErrors []string

	// Check for machines.hcl file
	machinesFile := filepath.Join(projectPath, "machines.hcl")
	if _, err := os.Stat(machinesFile); err == nil {
		m.logger.Debug("Found machines.hcl file", map[string]interface{}{
			"file_path": machinesFile,
		})

		machines, err := m.loader.LoadMachinesFromFile(ctx, machinesFile)
		if err != nil {
			loadErrors = append(loadErrors, fmt.Sprintf("machines.hcl: %v", err))
		} else {
			allMachines = append(allMachines, machines...)
			m.logger.Info("Loaded machines from file", map[string]interface{}{
				"file_path": machinesFile,
				"count":     len(machines),
			})
		}
	}

	// Check for machines/ directory
	machinesDir := filepath.Join(projectPath, "machines")
	if _, err := os.Stat(machinesDir); err == nil {
		m.logger.Debug("Found machines directory", map[string]interface{}{
			"dir_path": machinesDir,
		})

		machines, err := m.loader.LoadMachinesFromDirectory(ctx, machinesDir)
		if err != nil {
			loadErrors = append(loadErrors, fmt.Sprintf("machines/ directory: %v", err))
		} else {
			allMachines = append(allMachines, machines...)
			m.logger.Info("Loaded machines from directory", map[string]interface{}{
				"dir_path": machinesDir,
				"count":    len(machines),
			})
		}
	}

	// Check if no machines were found
	if len(allMachines) == 0 {
		if len(loadErrors) > 0 {
			return nil, fmt.Errorf("failed to load machines: %s", strings.Join(loadErrors, "; "))
		}
		return nil, fmt.Errorf("no machines found in project: %s (neither machines.hcl nor machines/ directory found)", projectPath)
	}

	// Validate for duplicates and consistency
	if err := m.validateMachineCollection(allMachines); err != nil {
		return nil, fmt.Errorf("machine validation failed: %w", err)
	}

	m.logger.Info("Machines loaded successfully", map[string]interface{}{
		"project_path": projectPath,
		"total_count":  len(allMachines),
		"load_errors":  len(loadErrors),
	})

	return allMachines, nil
}

// ExportMachines exports machines to HCL format according to machines schema
func (m *Manager) ExportMachines(ctx context.Context, machines []spookytypes.Machine, outputPath string) error {
	m.logger.Info("Exporting machines to HCL", map[string]interface{}{
		"count":       len(machines),
		"output_path": outputPath,
	})

	// Create HCL file
	file := hclwrite.NewEmptyFile()
	rootBody := file.Body()

	// Create machines block
	machinesBlock := rootBody.AppendNewBlock("machines", nil)
	machinesBody := machinesBlock.Body()

	// Add each machine
	for _, machine := range machines {
		machineBlock := machinesBody.AppendNewBlock("machine", []string{machine.Hostname})
		machineBody := machineBlock.Body()

		// Add basic attributes
		if machine.Host != "" && machine.Host != machine.Hostname {
			machineBody.SetAttributeValue("host", cty.StringVal(machine.Host))
		}
		if machine.Port != 22 {
			machineBody.SetAttributeValue("port", cty.NumberIntVal(int64(machine.Port)))
		}
		if machine.User != "" {
			machineBody.SetAttributeValue("user", cty.StringVal(machine.User))
		}
		if machine.KeyFile != "" {
			machineBody.SetAttributeValue("key_file", cty.StringVal(machine.KeyFile))
		}
		if machine.Password != "" {
			machineBody.SetAttributeValue("password", cty.StringVal(machine.Password))
		}

		// Add tags if present
		if len(machine.Tags) > 0 {
			tagsMap := make(map[string]cty.Value)
			for k, v := range machine.Tags {
				tagsMap[k] = cty.StringVal(v)
			}
			machineBody.SetAttributeValue("tags", cty.MapVal(tagsMap))
		}

		// Add groups if present
		if len(machine.Groups) > 0 {
			groupsList := make([]cty.Value, len(machine.Groups))
			for i, group := range machine.Groups {
				groupsList[i] = cty.StringVal(group)
			}
			machineBody.SetAttributeValue("groups", cty.ListVal(groupsList))
		}

		// Add roles if present
		if len(machine.Roles) > 0 {
			rolesList := make([]cty.Value, len(machine.Roles))
			for i, role := range machine.Roles {
				rolesList[i] = cty.StringVal(role)
			}
			machineBody.SetAttributeValue("roles", cty.ListVal(rolesList))
		}

		// Add classes if present
		if len(machine.Classes) > 0 {
			classesList := make([]cty.Value, len(machine.Classes))
			for i, class := range machine.Classes {
				classesList[i] = cty.StringVal(class)
			}
			machineBody.SetAttributeValue("classes", cty.ListVal(classesList))
		}

		// Add connection settings if different from defaults
		if machine.ConnectionTimeout != 0 {
			machineBody.SetAttributeValue("connection_timeout", cty.NumberIntVal(int64(machine.ConnectionTimeout)))
		}
		if machine.CommandTimeout != 0 {
			machineBody.SetAttributeValue("command_timeout", cty.NumberIntVal(int64(machine.CommandTimeout)))
		}
		if machine.MaxConnections != 0 {
			machineBody.SetAttributeValue("max_connections", cty.NumberIntVal(int64(machine.MaxConnections)))
		}
		if machine.RetryAttempts != 0 {
			machineBody.SetAttributeValue("retry_attempts", cty.NumberIntVal(int64(machine.RetryAttempts)))
		}
		if machine.RetryDelay != 0 {
			machineBody.SetAttributeValue("retry_delay", cty.NumberIntVal(int64(machine.RetryDelay)))
		}

		// Add resources block if present
		if machine.Resources != nil {
			resourcesBlock := machineBody.AppendNewBlock("resources", nil)
			resourcesBody := resourcesBlock.Body()

			if machine.Resources.CPUCores != 0 {
				resourcesBody.SetAttributeValue("cpu_cores", cty.NumberIntVal(int64(machine.Resources.CPUCores)))
			}
			if machine.Resources.MemoryGB != 0 {
				resourcesBody.SetAttributeValue("memory_gb", cty.NumberIntVal(int64(machine.Resources.MemoryGB)))
			}
			if machine.Resources.DiskGB != 0 {
				resourcesBody.SetAttributeValue("disk_gb", cty.NumberIntVal(int64(machine.Resources.DiskGB)))
			}
			if machine.Resources.NetworkSpeed != "" {
				resourcesBody.SetAttributeValue("network_speed", cty.StringVal(machine.Resources.NetworkSpeed))
			}
		}

		// Add metadata block if present
		if machine.MachineMetadata != nil {
			metadataBlock := machineBody.AppendNewBlock("metadata", nil)
			metadataBody := metadataBlock.Body()

			if machine.MachineMetadata.Environment != "" {
				metadataBody.SetAttributeValue("environment", cty.StringVal(machine.MachineMetadata.Environment))
			}
			if machine.MachineMetadata.Location != "" {
				metadataBody.SetAttributeValue("location", cty.StringVal(machine.MachineMetadata.Location))
			}
			if machine.MachineMetadata.Owner != "" {
				metadataBody.SetAttributeValue("owner", cty.StringVal(machine.MachineMetadata.Owner))
			}
			if len(machine.MachineMetadata.CustomFields) > 0 {
				customFieldsMap := make(map[string]cty.Value)
				for k, v := range machine.MachineMetadata.CustomFields {
					customFieldsMap[k] = cty.StringVal(v)
				}
				metadataBody.SetAttributeValue("custom_fields", cty.MapVal(customFieldsMap))
			}
		}
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write HCL file
	if err := os.WriteFile(outputPath, file.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write HCL file: %w", err)
	}

	m.logger.Info("Successfully exported machines to HCL", map[string]interface{}{
		"output_path": outputPath,
		"count":       len(machines),
	})

	return nil
}

// GetMachineByName looks up a machine by hostname
func (m *Manager) GetMachineByName(ctx context.Context, name string) (*spookytypes.Machine, error) {
	if name == "" {
		return nil, fmt.Errorf("machine name cannot be empty")
	}

	// Load machines from the current project context
	// Note: This assumes we have access to the project path from context
	// In a real implementation, this would be passed as a parameter or stored in the manager
	projectPath := m.getProjectPathFromContext(ctx)
	if projectPath == "" {
		return nil, fmt.Errorf("project path not available in context")
	}

	machines, err := m.LoadMachines(ctx, projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load machines: %w", err)
	}

	// Find machine by hostname
	for i := range machines {
		if machines[i].Hostname == name {
			return &machines[i], nil
		}
	}

	return nil, fmt.Errorf("machine '%s' not found in inventory", name)
}

// GetMachinesByTags filters machines by tags (supports key=value and key-only matching)
func (m *Manager) GetMachinesByTags(ctx context.Context, tags []string) ([]spookytypes.Machine, error) {
	if len(tags) == 0 {
		return nil, fmt.Errorf("at least one tag must be specified")
	}

	// Load machines from the current project context
	projectPath := m.getProjectPathFromContext(ctx)
	if projectPath == "" {
		return nil, fmt.Errorf("project path not available in context")
	}

	machines, err := m.LoadMachines(ctx, projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load machines: %w", err)
	}

	var filteredMachines []spookytypes.Machine

	// Filter machines by tags
	for i := range machines {
		if m.machineHasTags(&machines[i], tags) {
			filteredMachines = append(filteredMachines, machines[i])
		}
	}

	m.logger.Debug("Filtered machines by tags", map[string]interface{}{
		"tags":              tags,
		"total_machines":    len(machines),
		"filtered_machines": len(filteredMachines),
	})

	return filteredMachines, nil
}

// GetFullInventory returns the complete machine inventory
func (m *Manager) GetFullInventory(ctx context.Context) ([]spookytypes.Machine, error) {
	projectPath := m.getProjectPathFromContext(ctx)
	if projectPath == "" {
		return nil, fmt.Errorf("project path not available in context")
	}

	return m.LoadMachines(ctx, projectPath)
}

// GetMachinesByFilter applies complex filtering criteria to machines
func (m *Manager) GetMachinesByFilter(ctx context.Context, filter interface{}) ([]spookytypes.Machine, error) {
	projectPath := m.getProjectPathFromContext(ctx)
	if projectPath == "" {
		return nil, fmt.Errorf("project path not available in context")
	}

	machines, err := m.LoadMachines(ctx, projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load machines: %w", err)
	}

	// Apply filter based on type
	switch f := filter.(type) {
	case *spookytypesmachines.MachineFilter:
		return m.applyMachineFilter(machines, f)
	case map[string]interface{}:
		return m.applyMapFilter(machines, f)
	default:
		return nil, fmt.Errorf("unsupported filter type: %T", filter)
	}
}

// ValidateMachines validates machines
func (m *Manager) ValidateMachines(ctx context.Context, machines []spookytypes.Machine) (*spookytypes.ValidationResult, error) {
	m.logger.Debug("Validating machines", map[string]interface{}{
		"count": len(machines),
	})

	result, err := m.validator.ValidateMachines(ctx, machines)
	if err != nil {
		return nil, fmt.Errorf("failed to validate machines: %w", err)
	}

	m.logger.Info("Machine validation completed", map[string]interface{}{
		"count": len(machines),
		"valid": len(result.Errors) == 0,
	})

	return result, nil
}

// PingMachines pings machines to check connectivity
func (m *Manager) PingMachines(ctx context.Context, machines []spookytypes.Machine) ([]spookytypes.MachineStatus, error) {
	m.logger.Info("Pinging machines", map[string]interface{}{
		"count": len(machines),
	})

	var statuses []spookytypes.MachineStatus

	for i := range machines {
		status, err := m.pingMachine(ctx, &machines[i])
		if err != nil {
			m.logger.Warn("Failed to ping machine", map[string]interface{}{
				"hostname": machines[i].Hostname,
				"error":    err.Error(),
			})
			// Create error status
			status = &spookytypes.MachineStatus{
				Machine:   &machines[i],
				Status:    "error",
				LastCheck: time.Now(),
				Error:     err.Error(),
			}
		}
		statuses = append(statuses, *status)
	}

	// Count statuses
	online := 0
	offline := 0
	errors := 0
	for _, status := range statuses {
		switch status.Status {
		case "online":
			online++
		case "offline":
			offline++
		case "error":
			errors++
		}
	}

	m.logger.Info("Machine ping completed", map[string]interface{}{
		"count":    len(machines),
		"statuses": len(statuses),
		"online":   online,
		"offline":  offline,
		"errors":   errors,
	})

	return statuses, nil
}

// pingMachine pings a single machine to check connectivity
func (m *Manager) pingMachine(ctx context.Context, machine *spookytypes.Machine) (*spookytypes.MachineStatus, error) {
	startTime := time.Now()

	// Step 1: DNS Resolution (if hostname is not an IP)
	host := machine.Host
	if !isIPAddress(host) {
		ips, err := net.LookupHost(host)
		if err != nil {
			return &spookytypes.MachineStatus{
				Machine:   machine,
				Status:    "offline",
				LastCheck: time.Now(),
				Error:     fmt.Sprintf("DNS resolution failed: %v", err),
				Latency:   0,
			}, nil
		}
		if len(ips) == 0 {
			return &spookytypes.MachineStatus{
				Machine:   machine,
				Status:    "offline",
				LastCheck: time.Now(),
				Error:     "DNS resolution returned no IP addresses",
				Latency:   0,
			}, nil
		}
		// Use the first IP for connectivity checks
		host = ips[0]
	}

	// Step 2: TCP Port Scan for SSH
	sshPort := strconv.Itoa(machine.Port)
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, sshPort), 5*time.Second)
	if err != nil {
		return &spookytypes.MachineStatus{
			Machine:   machine,
			Status:    "offline",
			LastCheck: time.Now(),
			Error:     fmt.Sprintf("SSH port %s not reachable: %v", sshPort, err),
			Latency:   0,
		}, nil
	}
	conn.Close()

	// Step 3: SSH Authentication Test (if SSH manager is available)
	var sshStatus string
	var sshError string

	// Create SSH manager for authentication testing
	logManager := spookylogging.NewLogManager()
	sshLogger := logManager.GetLogger("ssh")
	sshManager := spookyssh.NewManager(sshLogger)

	// Test SSH connectivity with authentication
	connectionRequest := &spookytypes.ConnectionRequest{
		Host:     machine.Host,
		Port:     machine.Port,
		User:     machine.User,
		Password: machine.Password,
		KeyPath:  machine.KeyFile,
		Timeout:  10 * time.Second, // Short timeout for ping
	}

	_, err = sshManager.Connect(ctx, connectionRequest)
	if err != nil {
		sshStatus = "ssh_unreachable"
		sshError = err.Error()
	} else {
		sshStatus = "reachable"
	}

	// Calculate latency
	latency := int(time.Since(startTime).Milliseconds())

	// Determine final status
	var status string
	var errorMsg string

	if sshStatus == "reachable" {
		status = "online"
	} else if sshStatus == "ssh_unreachable" {
		status = "ssh_unreachable"
		errorMsg = sshError
	} else {
		status = "offline"
		errorMsg = "connectivity check failed"
	}

	return &spookytypes.MachineStatus{
		Machine:   machine,
		Status:    status,
		LastCheck: time.Now(),
		Error:     errorMsg,
		Latency:   latency,
	}, nil
}

// isIPAddress checks if a string is a valid IP address
func isIPAddress(host string) bool {
	// Check for IPv4
	if net.ParseIP(host) != nil {
		return true
	}

	// Check for IPv6 (enclosed in brackets)
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		ip := strings.Trim(host, "[]")
		return net.ParseIP(ip) != nil
	}

	return false
}

// getMachineHostnames returns a list of machine hostnames
func getMachineHostnames(machines []spookytypes.Machine) []string {
	var hostnames []string
	for _, machine := range machines {
		hostnames = append(hostnames, machine.Hostname)
	}
	return hostnames
}

// validateMachineCollection validates the entire machine collection
func (m *Manager) validateMachineCollection(machines []spookytypes.Machine) error {
	var errors []string
	var warnings []string

	// Check for duplicate hostnames with file source information
	hostnameMap := make(map[string][]string) // hostname -> []file_sources
	for _, machine := range machines {
		sourceFile := "unknown"
		if machine.MachineMetadata != nil && machine.MachineMetadata.CustomFields != nil {
			if src, exists := machine.MachineMetadata.CustomFields["source_file"]; exists {
				sourceFile = src
			}
		}
		hostnameMap[machine.Hostname] = append(hostnameMap[machine.Hostname], sourceFile)
	}

	for hostname, sources := range hostnameMap {
		if len(sources) > 1 {
			// Remove duplicates from sources
			uniqueSources := make([]string, 0)
			sourceMap := make(map[string]bool)
			for _, source := range sources {
				if !sourceMap[source] {
					sourceMap[source] = true
					uniqueSources = append(uniqueSources, source)
				}
			}
			errors = append(errors, fmt.Sprintf("duplicate hostname '%s' found in multiple files: %v", hostname, uniqueSources))
		}
	}

	// Check for duplicate host addresses with file source information
	hostMap := make(map[string][]string)       // host -> []hostnames
	hostSourceMap := make(map[string][]string) // host -> []file_sources
	for _, machine := range machines {
		sourceFile := "unknown"
		if machine.MachineMetadata != nil && machine.MachineMetadata.CustomFields != nil {
			if src, exists := machine.MachineMetadata.CustomFields["source_file"]; exists {
				sourceFile = src
			}
		}
		hostMap[machine.Host] = append(hostMap[machine.Host], machine.Hostname)
		hostSourceMap[machine.Host] = append(hostSourceMap[machine.Host], sourceFile)
	}

	for host, hostnames := range hostMap {
		if len(hostnames) > 1 {
			sources := hostSourceMap[host]
			uniqueSources := make([]string, 0)
			sourceMap := make(map[string]bool)
			for _, source := range sources {
				if !sourceMap[source] {
					sourceMap[source] = true
					uniqueSources = append(uniqueSources, source)
				}
			}
			warnings = append(warnings, fmt.Sprintf("duplicate host address '%s' used by multiple machines (%v) in files: %v", host, hostnames, uniqueSources))
		}
	}

	// Validate environment consistency
	if err := m.validateEnvironmentConsistency(machines); err != nil {
		warnings = append(warnings, fmt.Sprintf("environment consistency: %v", err))
	}

	// Validate authentication consistency
	if err := m.validateAuthenticationConsistency(machines); err != nil {
		warnings = append(warnings, fmt.Sprintf("authentication consistency: %v", err))
	}

	// Validate cross-file consistency
	if err := m.validateCrossFileConsistency(machines); err != nil {
		warnings = append(warnings, fmt.Sprintf("cross-file consistency: %v", err))
	}

	// Report errors and warnings
	if len(errors) > 0 {
		errorMsg := strings.Join(errors, "; ")
		m.logger.Error("Machine validation errors", fmt.Errorf("%s", errorMsg), map[string]interface{}{
			"errors": errors,
		})
		return fmt.Errorf("validation errors: %s", errorMsg)
	}

	if len(warnings) > 0 {
		m.logger.Warn("Machine validation warnings", map[string]interface{}{
			"warnings": warnings,
		})
	}

	return nil
}

// validateEnvironmentConsistency validates environment-specific rules
func (m *Manager) validateEnvironmentConsistency(machines []spookytypes.Machine) error {
	productionMachines := make([]spookytypes.Machine, 0)
	stagingMachines := make([]spookytypes.Machine, 0)

	// Group machines by environment
	for _, machine := range machines {
		if machine.MachineMetadata != nil {
			switch machine.MachineMetadata.Environment {
			case "production":
				productionMachines = append(productionMachines, machine)
			case "staging":
				stagingMachines = append(stagingMachines, machine)
			}
		}
	}

	// Validate production machines
	for _, machine := range productionMachines {
		// Production machines should have proper authentication
		if machine.KeyFile == "" {
			return fmt.Errorf("production machine '%s' must use key-based authentication", machine.Hostname)
		}

		// Production machines should have reasonable timeouts
		if machine.ConnectionTimeout > 60 {
			return fmt.Errorf("production machine '%s' has excessive connection timeout (%ds)", machine.Hostname, machine.ConnectionTimeout)
		}

		// Production machines should have proper resource specifications
		if machine.Resources == nil {
			return fmt.Errorf("production machine '%s' should have resource specifications", machine.Hostname)
		}
	}

	// Validate staging machines
	for _, machine := range stagingMachines {
		// Staging machines should have reasonable timeouts
		if machine.ConnectionTimeout > 120 {
			return fmt.Errorf("staging machine '%s' has excessive connection timeout (%ds)", machine.Hostname, machine.ConnectionTimeout)
		}
	}

	return nil
}

// validateAuthenticationConsistency validates authentication method consistency
func (m *Manager) validateAuthenticationConsistency(machines []spookytypes.Machine) error {
	keyBasedCount := 0

	for _, machine := range machines {
		if machine.KeyFile != "" {
			keyBasedCount++
		}
		// Note: Both password and key-based authentication are supported
		// Key-based authentication is recommended for security
	}

	// Recommend key-based authentication for all machines
	if keyBasedCount < len(machines) {
		return fmt.Errorf("recommend using key-based authentication for all machines (currently %d/%d use keys)", keyBasedCount, len(machines))
	}

	return nil
}

// validateCrossFileConsistency validates consistency across multiple files
func (m *Manager) validateCrossFileConsistency(machines []spookytypes.Machine) error {
	// Group machines by environment and check for consistency
	envGroups := make(map[string][]spookytypes.Machine)

	for _, machine := range machines {
		env := "unknown"
		if machine.MachineMetadata != nil && machine.MachineMetadata.Environment != "" {
			env = machine.MachineMetadata.Environment
		}
		envGroups[env] = append(envGroups[env], machine)
	}

	// Check for consistent authentication methods within environments
	for env, envMachines := range envGroups {
		if len(envMachines) < 2 {
			continue // Need at least 2 machines to check consistency
		}

		keyBasedCount := 0
		for _, machine := range envMachines {
			if machine.KeyFile != "" {
				keyBasedCount++
			}
		}

		// Check if authentication methods are consistent within environment
		if keyBasedCount > 0 && keyBasedCount < len(envMachines) {
			return fmt.Errorf("inconsistent authentication methods in %s environment (%d/%d use keys)", env, keyBasedCount, len(envMachines))
		}

		// Check for consistent timeout settings within environment
		timeoutValues := make(map[int]int)
		for _, machine := range envMachines {
			timeoutValues[machine.ConnectionTimeout]++
		}

		if len(timeoutValues) > 1 {
			return fmt.Errorf("inconsistent connection timeouts in %s environment: %v", env, timeoutValues)
		}
	}

	return nil
}

// machineHasTags checks if a machine has the specified tags
// Supports both key=value and key-only tag matching
func (m *Manager) machineHasTags(machine *spookytypes.Machine, tags []string) bool {
	if len(tags) == 0 {
		return true
	}

	// Handle case where machine has no tags
	if len(machine.Tags) == 0 {
		return false
	}

	// Check each required tag
	for _, requiredTag := range tags {
		tagMatched := false

		// Check if tag is in key=value format
		if strings.Contains(requiredTag, "=") {
			parts := strings.SplitN(requiredTag, "=", 2)
			if len(parts) == 2 {
				key, value := parts[0], parts[1]
				if machineValue, exists := machine.Tags[key]; exists && machineValue == value {
					tagMatched = true
				}
			}
		} else {
			// Key-only format - check if the key exists in machine tags
			if _, exists := machine.Tags[requiredTag]; exists {
				tagMatched = true
			}
		}

		// If any required tag doesn't match, machine doesn't have all required tags
		if !tagMatched {
			return false
		}
	}

	return true
}

// applyMachineFilter applies a MachineFilter to the machine list
func (m *Manager) applyMachineFilter(machines []spookytypes.Machine, filter *spookytypesmachines.MachineFilter) ([]spookytypes.Machine, error) {
	var filteredMachines []spookytypes.Machine

	for i := range machines {
		machine := &machines[i]

		// Check hostname filter
		if len(filter.Hostnames) > 0 {
			hostnameMatch := false
			for _, hostname := range filter.Hostnames {
				if machine.Hostname == hostname {
					hostnameMatch = true
					break
				}
			}
			if !hostnameMatch {
				continue
			}
		}

		// Check groups filter
		if len(filter.Groups) > 0 {
			groupMatch := false
			for _, group := range filter.Groups {
				for _, machineGroup := range machine.Groups {
					if machineGroup == group {
						groupMatch = true
						break
					}
				}
				if groupMatch {
					break
				}
			}
			if !groupMatch {
				continue
			}
		}

		// Check roles filter
		if len(filter.Roles) > 0 {
			roleMatch := false
			for _, role := range filter.Roles {
				for _, machineRole := range machine.Roles {
					if machineRole == role {
						roleMatch = true
						break
					}
				}
				if roleMatch {
					break
				}
			}
			if !roleMatch {
				continue
			}
		}

		// Check tags filter
		if len(filter.Tags) > 0 {
			tagMatch := false
			for key, value := range filter.Tags {
				if machineValue, exists := machine.Tags[key]; exists && machineValue == value {
					tagMatch = true
					break
				}
			}
			if !tagMatch {
				continue
			}
		}

		// Check patterns filter
		if len(filter.Patterns) > 0 {
			patternMatch := false
			for _, pattern := range filter.Patterns {
				if strings.Contains(machine.Hostname, pattern) {
					patternMatch = true
					break
				}
			}
			if !patternMatch {
				continue
			}
		}

		filteredMachines = append(filteredMachines, *machine)
	}

	return filteredMachines, nil
}

// applyMapFilter applies a map-based filter to the machine list
func (m *Manager) applyMapFilter(machines []spookytypes.Machine, filter map[string]interface{}) ([]spookytypes.Machine, error) {
	var filteredMachines []spookytypes.Machine

	for i := range machines {
		machine := &machines[i]
		match := true

		for key, value := range filter {
			switch key {
			case "hostname":
				if expected, ok := value.(string); ok && machine.Hostname != expected {
					match = false
				}
			case "host":
				if expected, ok := value.(string); ok && machine.Host != expected {
					match = false
				}
			case "user":
				if expected, ok := value.(string); ok && machine.User != expected {
					match = false
				}
			case "environment":
				if machine.MachineMetadata != nil {
					if expected, ok := value.(string); ok && machine.MachineMetadata.Environment != expected {
						match = false
					}
				} else {
					match = false
				}
			}
		}

		if match {
			filteredMachines = append(filteredMachines, *machine)
		}
	}

	return filteredMachines, nil
}

// getProjectPathFromContext extracts project path from context
func (m *Manager) getProjectPathFromContext(ctx context.Context) string {
	// Check if project path is stored in context
	if projectPath, ok := ctx.Value("project_path").(string); ok && projectPath != "" {
		return projectPath
	}

	// Check if project is stored in context
	if project, ok := ctx.Value("project").(*spookytypes.Project); ok && project != nil {
		return project.Path
	}

	// Check if project path is stored in context with different key
	if projectPath, ok := ctx.Value("projectPath").(string); ok && projectPath != "" {
		return projectPath
	}

	// If no project path found in context, return empty string
	// This allows callers to handle the case appropriately
	return ""
}
