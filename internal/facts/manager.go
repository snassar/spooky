// Package facts provides fact collection, storage, and management functionality.
package facts

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	spookytypes "spooky/internal/types"
	spookytypesschemas "spooky/internal/types/schemas"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// Manager implements FactManager
type Manager struct {
	collector FactCollector
	validator spookytypesschemas.SchemaValidator
	logger    spookytypes.Logger
}

// NewManager creates a new fact manager
func NewManager(
	collector FactCollector,
	validator spookytypesschemas.SchemaValidator,
	logger spookytypes.Logger,
) *Manager {
	return &Manager{
		collector: collector,
		validator: validator,
		logger:    logger,
	}
}

// CollectFacts collects facts from the given machine
func (m *Manager) CollectFacts(ctx context.Context, machine *spookytypes.Machine) (*FactCollection, error) {
	if machine == nil {
		return nil, fmt.Errorf("machine cannot be nil")
	}

	m.logger.Info("Collecting facts", map[string]interface{}{"machine": machine.Hostname})

	// Collect facts using the collector
	facts, err := m.collector.Collect(ctx, machine)
	if err != nil {
		m.logger.Error("Failed to collect facts", err, map[string]interface{}{"machine": machine.Hostname})
		return nil, fmt.Errorf("failed to collect facts for %s: %w", machine.Hostname, err)
	}

	m.logger.Info("Successfully collected facts", map[string]interface{}{"machine": machine.Hostname, "collector": m.collector.GetName()})

	return facts, nil
}

// ValidateFacts validates facts against schema
func (m *Manager) ValidateFacts(ctx context.Context, facts *FactCollection) (*spookytypes.ValidationResult, error) {
	if facts == nil {
		return &spookytypes.ValidationResult{
			Valid:    false,
			Errors:   []spookytypesschemas.SchemaError{{Message: "facts cannot be nil"}},
			Warnings: []spookytypesschemas.SchemaError{},
		}, nil
	}

	// Basic validation
	var errors []spookytypesschemas.SchemaError
	var warnings []spookytypesschemas.SchemaError

	// Validate machine ID
	if facts.MachineID == "" {
		errors = append(errors, spookytypesschemas.SchemaError{Message: "machine_id is required"})
	} else if !isValidMachineID(facts.MachineID) {
		errors = append(errors, spookytypesschemas.SchemaError{Message: "machine_id must be a 32-character hexadecimal string"})
	}

	// Validate collection timestamp
	if facts.CollectedAt.IsZero() {
		errors = append(errors, spookytypesschemas.SchemaError{Message: "collected_at is required"})
	}

	// Validate facts structure
	if facts.Facts == nil {
		errors = append(errors, spookytypesschemas.SchemaError{Message: "facts structure is required"})
	} else {
		// Validate system facts
		if facts.Facts.System == nil {
			errors = append(errors, spookytypesschemas.SchemaError{Message: "system facts are required"})
		} else {
			if facts.Facts.System.OS == nil {
				errors = append(errors, spookytypesschemas.SchemaError{Message: "system.os facts are required"})
			}
			if facts.Facts.System.Hardware == nil {
				errors = append(errors, spookytypesschemas.SchemaError{Message: "system.hardware facts are required"})
			}
			if facts.Facts.System.Network == nil {
				errors = append(errors, spookytypesschemas.SchemaError{Message: "system.network facts are required"})
			}
		}
	}

	// Schema validation would be implemented here using the validator
	// For now, we'll do basic validation

	valid := len(errors) == 0

	return &spookytypes.ValidationResult{
		Valid:    valid,
		Errors:   errors,
		Warnings: warnings,
	}, nil
}

// ExportFacts exports facts to the given format
func (m *Manager) ExportFacts(ctx context.Context, machineIDs []string, format string, outputPath string) error {
	m.logger.Info("Exporting facts", map[string]interface{}{
		"machines": len(machineIDs),
		"format":   format,
		"output":   outputPath,
	})

	// Collect facts for the specified machines
	var allFacts []*FactCollection
	for _, machineID := range machineIDs {
		facts, err := m.CollectFacts(ctx, &spookytypes.Machine{Hostname: machineID})
		if err != nil {
			m.logger.Error("Failed to collect facts for machine", err, map[string]interface{}{
				"machine": machineID,
			})
			return fmt.Errorf("failed to collect facts for machine %s: %w", machineID, err)
		}
		allFacts = append(allFacts, facts)
	}

	// Export based on format
	switch format {
	case "json":
		return m.exportToJSON(allFacts, outputPath)
	case "hcl":
		return m.exportToHCL(allFacts, outputPath)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

// exportToJSON exports facts to JSON format following facts-structure.schema.hcl
func (m *Manager) exportToJSON(facts []*FactCollection, outputPath string) error {
	// Export each fact collection individually following the schema structure
	var exportData []map[string]interface{}

	for _, fact := range facts {
		if fact == nil {
			continue
		}

		// Create the facts_structure format
		factStructure := map[string]interface{}{
			"machine_id":   fact.MachineID,
			"collected_at": fact.CollectedAt.Format(time.RFC3339),
			"facts":        fact.Facts,
		}

		exportData = append(exportData, factStructure)
	}

	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal facts to JSON: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON export file: %w", err)
	}

	m.logger.Info("Successfully exported facts to JSON", map[string]interface{}{"output": outputPath, "machines": len(facts)})
	return nil
}

// exportToHCL exports facts to HCL format following facts-structure.schema.hcl
func (m *Manager) exportToHCL(facts []*FactCollection, outputPath string) error {
	file := hclwrite.NewEmptyFile()
	rootBody := file.Body()

	// Export each fact collection following the facts_structure schema
	for _, fact := range facts {
		if fact == nil {
			continue
		}

		// Create facts_structure block
		factsStructureBlock := rootBody.AppendNewBlock("facts_structure", nil)
		factsStructureBody := factsStructureBlock.Body()

		// Add machine_id
		factsStructureBody.SetAttributeValue("machine_id", cty.StringVal(fact.MachineID))

		// Add collected_at
		factsStructureBody.SetAttributeValue("collected_at", cty.StringVal(fact.CollectedAt.Format(time.RFC3339)))

		// Add facts object
		if fact.Facts != nil {
			factsBlock := factsStructureBody.AppendNewBlock("facts", nil)
			factsBody := factsBlock.Body()

			// Add system facts if available
			if fact.Facts.System != nil {
				systemBlock := factsBody.AppendNewBlock("system", nil)
				systemBody := systemBlock.Body()

				// OS facts
				if fact.Facts.System.OS != nil {
					osBlock := systemBody.AppendNewBlock("os", nil)
					osBody := osBlock.Body()
					osBody.SetAttributeValue("name", cty.StringVal(fact.Facts.System.OS.Name))
					osBody.SetAttributeValue("version", cty.StringVal(fact.Facts.System.OS.Version))
					osBody.SetAttributeValue("arch", cty.StringVal(fact.Facts.System.OS.Arch))
					if fact.Facts.System.OS.Kernel != "" {
						osBody.SetAttributeValue("kernel", cty.StringVal(fact.Facts.System.OS.Kernel))
					}
					if fact.Facts.System.OS.Platform != "" {
						osBody.SetAttributeValue("platform", cty.StringVal(fact.Facts.System.OS.Platform))
					}
					if fact.Facts.System.OS.Family != "" {
						osBody.SetAttributeValue("family", cty.StringVal(fact.Facts.System.OS.Family))
					}
				}

				// Hardware facts
				if fact.Facts.System.Hardware != nil {
					hardwareBlock := systemBody.AppendNewBlock("hardware", nil)
					hardwareBody := hardwareBlock.Body()

					// CPU facts
					if fact.Facts.System.Hardware.CPU != nil {
						cpuBlock := hardwareBody.AppendNewBlock("cpu", nil)
						cpuBody := cpuBlock.Body()
						cpuBody.SetAttributeValue("cores", cty.NumberIntVal(int64(fact.Facts.System.Hardware.CPU.Cores)))
						cpuBody.SetAttributeValue("model", cty.StringVal(fact.Facts.System.Hardware.CPU.Model))
						if fact.Facts.System.Hardware.CPU.Frequency > 0 {
							cpuBody.SetAttributeValue("frequency", cty.NumberFloatVal(fact.Facts.System.Hardware.CPU.Frequency))
						}
						if fact.Facts.System.Hardware.CPU.Architecture != "" {
							cpuBody.SetAttributeValue("architecture", cty.StringVal(fact.Facts.System.Hardware.CPU.Architecture))
						}
						if fact.Facts.System.Hardware.CPU.Vendor != "" {
							cpuBody.SetAttributeValue("vendor", cty.StringVal(fact.Facts.System.Hardware.CPU.Vendor))
						}
					}

					// Memory facts
					if fact.Facts.System.Hardware.Memory != nil {
						memoryBlock := hardwareBody.AppendNewBlock("memory", nil)
						memoryBody := memoryBlock.Body()
						memoryBody.SetAttributeValue("total", cty.NumberIntVal(fact.Facts.System.Hardware.Memory.Total))
						if fact.Facts.System.Hardware.Memory.Available > 0 {
							memoryBody.SetAttributeValue("available", cty.NumberIntVal(fact.Facts.System.Hardware.Memory.Available))
						}
						if fact.Facts.System.Hardware.Memory.Used > 0 {
							memoryBody.SetAttributeValue("used", cty.NumberIntVal(fact.Facts.System.Hardware.Memory.Used))
						}
						if fact.Facts.System.Hardware.Memory.Free > 0 {
							memoryBody.SetAttributeValue("free", cty.NumberIntVal(fact.Facts.System.Hardware.Memory.Free))
						}
						if fact.Facts.System.Hardware.Memory.Buffers > 0 {
							memoryBody.SetAttributeValue("buffers", cty.NumberIntVal(fact.Facts.System.Hardware.Memory.Buffers))
						}
						if fact.Facts.System.Hardware.Memory.Cached > 0 {
							memoryBody.SetAttributeValue("cached", cty.NumberIntVal(fact.Facts.System.Hardware.Memory.Cached))
						}
					}
				}

				// Network facts
				if fact.Facts.System.Network != nil {
					networkBlock := systemBody.AppendNewBlock("network", nil)
					networkBody := networkBlock.Body()

					if fact.Facts.System.Network.Hostname != "" {
						networkBody.SetAttributeValue("hostname", cty.StringVal(fact.Facts.System.Network.Hostname))
					}

					if len(fact.Facts.System.Network.Interfaces) > 0 {
						interfacesBlock := networkBody.AppendNewBlock("interfaces", nil)
						interfacesBody := interfacesBlock.Body()

						for _, iface := range fact.Facts.System.Network.Interfaces {
							ifaceBlock := interfacesBody.AppendNewBlock("interface", []string{iface.Name})
							ifaceBody := ifaceBlock.Body()
							ifaceBody.SetAttributeValue("name", cty.StringVal(iface.Name))
							if iface.MACAddress != "" {
								ifaceBody.SetAttributeValue("mac_address", cty.StringVal(iface.MACAddress))
							}
							if len(iface.IPAddresses) > 0 {
								// Convert IP addresses to HCL list
								ipList := make([]cty.Value, len(iface.IPAddresses))
								for i, ip := range iface.IPAddresses {
									ipList[i] = cty.StringVal(ip)
								}
								ifaceBody.SetAttributeValue("ip_addresses", cty.ListVal(ipList))
							}
							if iface.MTU > 0 {
								ifaceBody.SetAttributeValue("mtu", cty.NumberIntVal(int64(iface.MTU)))
							}
						}
					}

					if len(fact.Facts.System.Network.IPAddresses) > 0 {
						// Convert IP addresses to HCL list
						ipList := make([]cty.Value, len(fact.Facts.System.Network.IPAddresses))
						for i, ip := range fact.Facts.System.Network.IPAddresses {
							ipList[i] = cty.StringVal(ip)
						}
						networkBody.SetAttributeValue("ip_addresses", cty.ListVal(ipList))
					}

					if fact.Facts.System.Network.PrimaryIP != "" {
						networkBody.SetAttributeValue("primary_ip", cty.StringVal(fact.Facts.System.Network.PrimaryIP))
					}
				}
			}
		}
	}

	// Write to file
	if err := os.WriteFile(outputPath, file.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write HCL export file: %w", err)
	}

	m.logger.Info("Successfully exported facts to HCL", map[string]interface{}{"output": outputPath, "machines": len(facts)})
	return nil
}

// CollectAndStoreFacts collects and stores facts for a machine
func (m *Manager) CollectAndStoreFacts(ctx context.Context, machine *spookytypes.Machine) error {
	// Collect facts
	facts, err := m.CollectFacts(ctx, machine)
	if err != nil {
		return err
	}

	// Validate facts
	validationResult, err := m.ValidateFacts(ctx, facts)
	if err != nil {
		return fmt.Errorf("facts validation failed: %w", err)
	}

	if !validationResult.Valid {
		return fmt.Errorf("facts validation failed: %v", validationResult.Errors)
	}

	return nil
}

// CollectAndStoreFactsParallel collects and stores facts for multiple machines in parallel
func (m *Manager) CollectAndStoreFactsParallel(ctx context.Context, machines []*spookytypes.Machine, maxWorkers int) error {
	if len(machines) == 0 {
		return fmt.Errorf("no machines provided")
	}

	if maxWorkers <= 0 {
		maxWorkers = 4
	}

	m.logger.Info("Starting parallel fact collection", map[string]interface{}{"machines": len(machines), "workers": maxWorkers})

	// Create worker pool
	semaphore := make(chan struct{}, maxWorkers)
	results := make(chan error, len(machines))

	// Start workers
	for _, machine := range machines {
		go func(machine *spookytypes.Machine) {
			semaphore <- struct{}{} // Acquire semaphore
			defer func() {
				<-semaphore // Release semaphore
			}()

			err := m.CollectAndStoreFacts(ctx, machine)
			results <- err
		}(machine)
	}

	// Collect results
	var errors []error
	for i := 0; i < len(machines); i++ {
		if err := <-results; err != nil {
			errors = append(errors, err)
		}
	}

	// Report results
	successful := len(machines) - len(errors)
	m.logger.Info("Completed parallel fact collection", map[string]interface{}{"successful": successful, "failed": len(errors)})

	if len(errors) > 0 {
		return fmt.Errorf("fact collection failed for %d machines", len(errors))
	}

	return nil
}

// GetStorageStats returns storage statistics
func (m *Manager) GetStorageStats() (map[string]interface{}, error) {
	return map[string]interface{}{
		"storage_type": "memory_only",
		"description":  "Direct export without intermediate storage",
	}, nil
}

// GetFacts retrieves facts for a specific machine
func (m *Manager) GetFacts(ctx context.Context, machineID string) (*FactCollection, error) {
	m.logger.Info("Getting facts for machine", map[string]interface{}{
		"machine": machineID,
	})

	// For memory storage, we need to collect facts on demand
	return m.CollectFacts(ctx, &spookytypes.Machine{Hostname: machineID})
}

// FactExport represents exported facts data
type FactExport struct {
	ExportedAt   time.Time                  `json:"exported_at"`
	Format       string                     `json:"format"`
	MachineCount int                        `json:"machine_count"`
	Facts        map[string]*FactCollection `json:"facts"`
}

// ClearFacts removes all facts from memory
func (m *Manager) ClearFacts(ctx context.Context) error {
	m.logger.Info("Clearing all facts from memory")

	// Memory-only storage - no cleanup needed
	m.logger.Info("Successfully cleared all facts from memory")

	return nil
}

// isValidMachineID validates machine ID format
func isValidMachineID(machineID string) bool {
	if len(machineID) != 32 {
		return false
	}

	for _, char := range machineID {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
			return false
		}
	}

	return true
}
