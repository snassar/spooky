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

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
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

	for _, machine := range machines {
		status, err := m.pingMachine(ctx, machine)
		if err != nil {
			m.logger.Warn("Failed to ping machine", map[string]interface{}{
				"hostname": machine.Hostname,
				"error":    err.Error(),
			})
			// Create error status
			status = &spookytypes.MachineStatus{
				Machine:   &machine,
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
func (m *Manager) pingMachine(ctx context.Context, machine spookytypes.Machine) (*spookytypes.MachineStatus, error) {
	startTime := time.Now()

	// Step 1: DNS Resolution (if hostname is not an IP)
	host := machine.Host
	if !isIPAddress(host) {
		ips, err := net.LookupHost(host)
		if err != nil {
			return &spookytypes.MachineStatus{
				Machine:   &machine,
				Status:    "offline",
				LastCheck: time.Now(),
				Error:     fmt.Sprintf("DNS resolution failed: %v", err),
				Latency:   0,
			}, nil
		}
		if len(ips) == 0 {
			return &spookytypes.MachineStatus{
				Machine:   &machine,
				Status:    "offline",
				LastCheck: time.Now(),
				Error:     "DNS resolution returned no IP addresses",
				Latency:   0,
			}, nil
		}
		// Use the first IP for connectivity checks
		host = ips[0]
	}

	// Step 2: ICMP Ping (simulated for now - would need root privileges)
	// In a real implementation, this would use raw sockets or external ping command
	icmpReachable := true // Simulated for now

	// Step 3: TCP Port Scan for SSH
	sshPort := strconv.Itoa(machine.Port)
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, sshPort), 5*time.Second)
	if err != nil {
		return &spookytypes.MachineStatus{
			Machine:   &machine,
			Status:    "offline",
			LastCheck: time.Now(),
			Error:     fmt.Sprintf("SSH port %s not reachable: %v", sshPort, err),
			Latency:   0,
		}, nil
	}
	conn.Close()

	// Calculate latency
	latency := int(time.Since(startTime).Milliseconds())

	// Step 4: SSH Authentication (skipped for now as requested)
	// This would be implemented later when SSH functionality is added

	// Determine final status
	var status string
	var errorMsg string

	if icmpReachable && latency > 0 {
		status = "online"
	} else {
		status = "offline"
		errorMsg = "connectivity check failed"
	}

	return &spookytypes.MachineStatus{
		Machine:   &machine,
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
		m.logger.Error("Machine validation errors", fmt.Errorf(errorMsg), map[string]interface{}{
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
		// Note: password authentication is not supported in current implementation
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
