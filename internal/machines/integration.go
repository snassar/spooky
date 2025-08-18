// Package machines provides machine management functionality for the spooky codebase.
package machines

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
)

// Machine connectivity status constants
const (
	MachineConnectivityReachable      = "reachable"
	MachineConnectivityUnreachable    = "unreachable"
	MachineConnectivitySSHUnreachable = "ssh_unreachable"
)

// Integration implements the MachinesIntegration interface
type Integration struct {
	manager    spookyinterfaces.MachinesIntegration
	loader     spookyinterfaces.MachineLoader
	validator  spookyinterfaces.MachineValidator
	sshManager spookyinterfaces.SSHManager
	logger     spookytypeslogging.Logger
}

// NewIntegration creates a new machines integration
func NewIntegration(
	manager spookyinterfaces.MachinesIntegration,
	loader spookyinterfaces.MachineLoader,
	validator spookyinterfaces.MachineValidator,
	sshManager spookyinterfaces.SSHManager,
	logger spookytypeslogging.Logger,
) spookyinterfaces.MachinesIntegration {
	return &Integration{
		manager:    manager,
		loader:     loader,
		validator:  validator,
		sshManager: sshManager,
		logger:     logger,
	}
}

// LoadMachines loads machines from the given source
func (i *Integration) LoadMachines(ctx context.Context, source string) ([]spookytypes.Machine, error) {
	i.logger.Debug("Loading machines from source", map[string]interface{}{
		"source": source,
	})

	// Check if source is a file or directory
	if strings.HasSuffix(source, ".hcl") {
		// Load from single file
		machines, err := i.loader.LoadMachinesFromFile(ctx, source)
		if err != nil {
			return nil, fmt.Errorf("failed to load machines from file %s: %w", source, err)
		}

		i.logger.Info("Loaded machines from file", map[string]interface{}{
			"file":  source,
			"count": len(machines),
		})

		return machines, nil
	}

	// Load from directory
	machines, err := i.loader.LoadMachinesFromDirectory(ctx, source)
	if err != nil {
		return nil, fmt.Errorf("failed to load machines from directory %s: %w", source, err)
	}

	i.logger.Info("Loaded machines from directory", map[string]interface{}{
		"directory": source,
		"count":     len(machines),
	})

	return machines, nil
}

// ValidateMachines validates machines
func (i *Integration) ValidateMachines(ctx context.Context, machines []spookytypes.Machine) (*spookytypes.ValidationResult, error) {
	i.logger.Debug("Validating machines", map[string]interface{}{
		"count": len(machines),
	})

	// Validate machines using the validator
	result, err := i.validator.ValidateMachines(ctx, machines)
	if err != nil {
		return nil, fmt.Errorf("failed to validate machines: %w", err)
	}

	i.logger.Info("Machines validation completed", map[string]interface{}{
		"count":  len(machines),
		"valid":  result.Valid,
		"errors": len(result.Errors),
	})

	return result, nil
}

// PingMachines pings machines to check connectivity
func (i *Integration) PingMachines(ctx context.Context, machines []spookytypes.Machine) ([]spookytypes.MachineStatus, error) {
	i.logger.Debug("Pinging machines", map[string]interface{}{
		"count": len(machines),
	})

	var statuses []spookytypes.MachineStatus

	for idx := range machines {
		machine := &machines[idx]
		status := spookytypes.MachineStatus{
			Machine:   machine,
			Status:    "unknown",
			LastCheck: time.Now(),
		}

		// Create connection request
		request := &spookytypes.ConnectionRequest{
			Host: machine.Host,
			Port: machine.Port,
			User: machine.User,
		}

		// Validate connection parameters
		validationResult, err := i.sshManager.ValidateConnection(ctx, request)
		if err != nil {
			status.Status = "invalid"
			status.Error = fmt.Sprintf("connection validation failed: %v", err)
			statuses = append(statuses, status)
			continue
		}

		if !validationResult.Valid {
			status.Status = "invalid"
			status.Error = "connection parameters invalid"
			statuses = append(statuses, status)
			continue
		}

		// Attempt SSH connection
		connectionResult, err := i.sshManager.Connect(ctx, request)
		if err != nil {
			status.Status = MachineConnectivityUnreachable
			status.Error = fmt.Sprintf("SSH connection failed: %v", err)
			statuses = append(statuses, status)
			continue
		}

		if !connectionResult.Success {
			status.Status = MachineConnectivityUnreachable
			status.Error = connectionResult.Error
			statuses = append(statuses, status)
			continue
		}

		// Connection successful
		status.Status = MachineConnectivityReachable
		status.Latency = int(connectionResult.ConnectTime.Milliseconds())
		statuses = append(statuses, status)
	}

	i.logger.Info("Machine ping completed", map[string]interface{}{
		"total":       len(machines),
		"reachable":   countReachableMachines(statuses),
		"unreachable": countUnreachableMachines(statuses),
	})

	return statuses, nil
}

// ExportMachines exports machines to HCL format according to machines schema
func (i *Integration) ExportMachines(ctx context.Context, machines []spookytypes.Machine, outputPath string) error {
	i.logger.Debug("Exporting machines to HCL", map[string]interface{}{
		"count":       len(machines),
		"output_path": outputPath,
	})

	// Validate machines before export
	validationResult, err := i.ValidateMachines(ctx, machines)
	if err != nil {
		return fmt.Errorf("failed to validate machines before export: %w", err)
	}

	if !validationResult.Valid {
		return fmt.Errorf("cannot export invalid machines: %d validation errors", len(validationResult.Errors))
	}

	// Create HCL content
	hclContent, err := i.generateHCLContent(machines)
	if err != nil {
		return fmt.Errorf("failed to generate HCL content: %w", err)
	}

	// Write to file
	if err := os.WriteFile(outputPath, []byte(hclContent), 0o600); err != nil {
		return fmt.Errorf("failed to write machines to file %s: %w", outputPath, err)
	}

	i.logger.Info("Machines exported successfully", map[string]interface{}{
		"count":       len(machines),
		"output_path": outputPath,
	})

	return nil
}

// GetMachineByName looks up a machine by hostname
func (i *Integration) GetMachineByName(ctx context.Context, name string) (*spookytypes.Machine, error) {
	i.logger.Debug("Looking up machine by name", map[string]interface{}{
		"name": name,
	})

	// Load all machines from the current project
	projectPath := i.getCurrentProjectPath()
	machines, err := i.LoadMachines(ctx, projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load machines: %w", err)
	}

	// Find machine by name
	for idx := range machines {
		machine := &machines[idx]
		if machine.Hostname == name {
			i.logger.Debug("Found machine by name", map[string]interface{}{
				"name":     name,
				"hostname": machine.Hostname,
			})
			return machine, nil
		}
	}

	i.logger.Debug("Machine not found by name", map[string]interface{}{
		"name": name,
	})

	return nil, fmt.Errorf("machine not found: %s", name)
}

// GetMachinesByTags filters machines by tags (supports key=value and key-only matching)
func (i *Integration) GetMachinesByTags(ctx context.Context, tags []string) ([]spookytypes.Machine, error) {
	i.logger.Debug("Filtering machines by tags", map[string]interface{}{
		"tags": tags,
	})

	// Load all machines from the current project
	projectPath := i.getCurrentProjectPath()
	machines, err := i.LoadMachines(ctx, projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load machines: %w", err)
	}

	var filteredMachines []spookytypes.Machine

	for idx := range machines {
		if i.machineMatchesTags(&machines[idx], tags) {
			filteredMachines = append(filteredMachines, machines[idx])
		}
	}

	i.logger.Info("Machines filtered by tags", map[string]interface{}{
		"tags":           tags,
		"total_machines": len(machines),
		"filtered_count": len(filteredMachines),
	})

	return filteredMachines, nil
}

// GetFullInventory returns the complete machine inventory
func (i *Integration) GetFullInventory(ctx context.Context) ([]spookytypes.Machine, error) {
	i.logger.Debug("Getting full machine inventory")

	// Load all machines from the current project
	projectPath := i.getCurrentProjectPath()
	machines, err := i.LoadMachines(ctx, projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load full inventory: %w", err)
	}

	i.logger.Info("Full inventory retrieved", map[string]interface{}{
		"count": len(machines),
	})

	return machines, nil
}

// GetMachinesByFilter applies complex filtering criteria to machines
func (i *Integration) GetMachinesByFilter(ctx context.Context, filter interface{}) ([]spookytypes.Machine, error) {
	i.logger.Debug("Filtering machines with complex filter", map[string]interface{}{
		"filter": filter,
	})

	// Load all machines from the current project
	projectPath := i.getCurrentProjectPath()
	machines, err := i.LoadMachines(ctx, projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load machines for filtering: %w", err)
	}

	// Apply complex filter
	filteredMachines, err := i.applyComplexFilter(machines, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to apply complex filter: %w", err)
	}

	i.logger.Info("Machines filtered with complex filter", map[string]interface{}{
		"total_machines": len(machines),
		"filtered_count": len(filteredMachines),
	})

	return filteredMachines, nil
}

// SaveMachines saves machines to the given destination
func (i *Integration) SaveMachines(ctx context.Context, machines []spookytypes.Machine, destination string) error {
	return i.manager.SaveMachines(ctx, machines, destination)
}

// EncryptMachines encrypts all machine secrets that have encrypted=true
func (i *Integration) EncryptMachines(ctx context.Context, projectPath string, secretsIntegration spookyinterfaces.SecretsIntegration, recipients []string, dryRun bool) error {
	return i.manager.EncryptMachines(ctx, projectPath, secretsIntegration, recipients, dryRun)
}

// DecryptMachines decrypts age-encrypted values in machines for debugging
func (i *Integration) DecryptMachines(ctx context.Context, machines []spookytypes.Machine, secretsIntegration spookyinterfaces.SecretsIntegration, identityPath string) error {
	return i.manager.DecryptMachines(ctx, machines, secretsIntegration, identityPath)
}

// Helper methods

// countReachableMachines counts reachable machines
func countReachableMachines(statuses []spookytypes.MachineStatus) int {
	count := 0
	for _, status := range statuses {
		if status.Status == MachineConnectivityReachable {
			count++
		}
	}
	return count
}

// countUnreachableMachines counts unreachable machines
func countUnreachableMachines(statuses []spookytypes.MachineStatus) int {
	count := 0
	for _, status := range statuses {
		if status.Status == MachineConnectivityUnreachable {
			count++
		}
	}
	return count
}

// generateHCLContent generates HCL content for machines
func (i *Integration) generateHCLContent(machines []spookytypes.Machine) (string, error) {
	var content strings.Builder

	content.WriteString("machines {\n")

	for idx := range machines {
		machine := &machines[idx]
		content.WriteString(fmt.Sprintf("  machine %q {\n", machine.Host))
		content.WriteString(fmt.Sprintf("    host = %q\n", machine.Host))
		content.WriteString(fmt.Sprintf("    port = %d\n", machine.Port))
		content.WriteString(fmt.Sprintf("    user = %q\n", machine.User))

		// Add authentication fields if present
		if machine.KeyFile != "" {
			content.WriteString(fmt.Sprintf("    key_file = %q\n", machine.KeyFile))
		}
		if machine.Password != "" {
			content.WriteString(fmt.Sprintf("    password = %q\n", machine.Password))
		}
		if machine.Passphrase != "" {
			content.WriteString(fmt.Sprintf("    passphrase = %q\n", machine.Passphrase))
		}

		// Add tags if present
		if len(machine.Tags) > 0 {
			content.WriteString("    tags = [\n")
			for _, tag := range machine.Tags {
				content.WriteString(fmt.Sprintf("      %q,\n", tag))
			}
			content.WriteString("    ]\n")
		}

		content.WriteString("  }\n")
	}

	content.WriteString("}\n")

	return content.String(), nil
}

// getCurrentProjectPath gets the current project path
func (i *Integration) getCurrentProjectPath() string {
	// This would typically be determined from context or configuration
	// For now, return a default path
	return "."
}

// machineMatchesTags checks if a machine matches the given tags
func (i *Integration) machineMatchesTags(machine *spookytypes.Machine, tags []string) bool {
	machineTags := make(map[string]string)
	for _, tag := range machine.Tags {
		if strings.Contains(tag, "=") {
			parts := strings.SplitN(tag, "=", 2)
			if len(parts) == 2 {
				machineTags[parts[0]] = parts[1]
			}
		} else {
			machineTags[tag] = ""
		}
	}

	for _, filterTag := range tags {
		if strings.Contains(filterTag, "=") {
			// Key=value matching
			parts := strings.SplitN(filterTag, "=", 2)
			if len(parts) == 2 {
				if value, exists := machineTags[parts[0]]; !exists || value != parts[1] {
					return false
				}
			}
		} else {
			// Key-only matching
			if _, exists := machineTags[filterTag]; !exists {
				return false
			}
		}
	}

	return true
}

// applyComplexFilter applies complex filtering criteria to machines
func (i *Integration) applyComplexFilter(machines []spookytypes.Machine, _ interface{}) ([]spookytypes.Machine, error) {
	// This would implement complex filtering logic based on the filter interface
	// For now, return all machines
	return machines, nil
}
