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
	storage   FactStorage
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
		storage:   NewMemoryFactStorage(),
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

// StoreFacts stores facts for a machine
func (m *Manager) StoreFacts(ctx context.Context, machineID string, facts *FactCollection) error {
	if machineID == "" {
		return fmt.Errorf("machine ID cannot be empty")
	}

	if facts == nil {
		return fmt.Errorf("facts cannot be nil")
	}

	m.logger.Info("Storing facts", map[string]interface{}{"machine_id": machineID})

	// Validate facts before storing
	validationResult, err := m.ValidateFacts(ctx, facts)
	if err != nil {
		m.logger.Error("Facts validation failed", err, map[string]interface{}{"machine_id": machineID})
		return fmt.Errorf("facts validation failed for %s: %w", machineID, err)
	}

	if !validationResult.Valid {
		m.logger.Error("Facts validation failed", fmt.Errorf("validation errors: %v", validationResult.Errors), map[string]interface{}{"machine_id": machineID})
		return fmt.Errorf("facts validation failed for %s: %w", machineID, fmt.Errorf("validation errors: %v", validationResult.Errors))
	}

	// Store facts
	if err := m.storage.Store(ctx, machineID, facts); err != nil {
		m.logger.Error("Failed to store facts", err, map[string]interface{}{"machine_id": machineID})
		return fmt.Errorf("failed to store facts for %s: %w", machineID, err)
	}

	m.logger.Info("Successfully stored facts", map[string]interface{}{"machine_id": machineID})

	return nil
}

// GetFacts retrieves facts for a machine
func (m *Manager) GetFacts(ctx context.Context, machineID string) (*FactCollection, error) {
	if machineID == "" {
		return nil, fmt.Errorf("machine ID cannot be empty")
	}

	m.logger.Debug("Retrieving facts", map[string]interface{}{"machine_id": machineID})

	facts, err := m.storage.Get(ctx, machineID)
	if err != nil {
		m.logger.Error("Failed to retrieve facts", err, map[string]interface{}{"machine_id": machineID})
		return nil, fmt.Errorf("failed to retrieve facts for %s: %w", machineID, err)
	}

	m.logger.Debug("Successfully retrieved facts", map[string]interface{}{"machine_id": machineID})

	return facts, nil
}

// ListFacts lists all machines with stored facts
func (m *Manager) ListFacts(ctx context.Context) ([]string, error) {
	m.logger.Debug("Listing facts")

	machineIDs, err := m.storage.List(ctx)
	if err != nil {
		m.logger.Error("Failed to list facts", err)
		return nil, fmt.Errorf("failed to list facts: %w", err)
	}

	m.logger.Debug("Successfully listed facts", map[string]interface{}{"count": len(machineIDs)})

	return machineIDs, nil
}

// DeleteFacts deletes facts for a machine
func (m *Manager) DeleteFacts(ctx context.Context, machineID string) error {
	if machineID == "" {
		return fmt.Errorf("machine ID cannot be empty")
	}

	m.logger.Info("Deleting facts", map[string]interface{}{"machine_id": machineID})

	if err := m.storage.Delete(ctx, machineID); err != nil {
		m.logger.Error("Failed to delete facts", err, map[string]interface{}{"machine_id": machineID})
		return fmt.Errorf("failed to delete facts for %s: %w", machineID, err)
	}

	m.logger.Info("Successfully deleted facts", map[string]interface{}{"machine_id": machineID})

	return nil
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
	if len(machineIDs) == 0 {
		return fmt.Errorf("no machine IDs provided for export")
	}

	if format == "" {
		format = "json"
	}

	if outputPath == "" {
		outputPath = fmt.Sprintf("facts-export-%s.%s", time.Now().Format("20060102-150405"), format)
	}

	m.logger.Info("Exporting facts", map[string]interface{}{"machines": len(machineIDs), "format": format, "output": outputPath})

	// Collect facts for all machines
	exportData := make(map[string]*FactCollection)

	for _, machineID := range machineIDs {
		facts, err := m.GetFacts(ctx, machineID)
		if err != nil {
			m.logger.Warn("Failed to get facts for export", map[string]interface{}{"machine_id": machineID, "error": err.Error()})
			continue
		}

		exportData[machineID] = facts
	}

	// Create export structure
	export := &FactExport{
		ExportedAt:   time.Now(),
		Format:       format,
		MachineCount: len(exportData),
		Facts:        exportData,
	}

	// Export based on format
	var data []byte
	var err error

	switch format {
	case "json":
		data, err = json.MarshalIndent(export, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal facts to JSON: %w", err)
		}

	case "hcl":
		data, err = m.exportToHCL(export)
		if err != nil {
			return fmt.Errorf("failed to marshal facts to HCL: %w", err)
		}

	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}

	// Write to file
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	m.logger.Info("Successfully exported facts", map[string]interface{}{"output": outputPath, "machines": len(exportData)})

	return nil
}

// ImportFacts imports facts from the given format
func (m *Manager) ImportFacts(ctx context.Context, format string, inputPath string) error {
	if format == "" {
		format = "json"
	}

	if inputPath == "" {
		return fmt.Errorf("input path is required")
	}

	m.logger.Info("Importing facts", map[string]interface{}{"format": format, "input": inputPath})

	// Read file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read import file: %w", err)
	}

	// Parse based on format
	var export *FactExport

	switch format {
	case "json":
		if err := json.Unmarshal(data, &export); err != nil {
			return fmt.Errorf("failed to unmarshal facts from JSON: %w", err)
		}

	case "hcl":
		// HCL import would be implemented here
		return fmt.Errorf("HCL import not implemented yet")

	default:
		return fmt.Errorf("unsupported import format: %s", format)
	}

	// Validate export structure
	if export == nil {
		return fmt.Errorf("invalid export structure")
	}

	// Import facts
	imported := 0
	failed := 0

	for machineID, facts := range export.Facts {
		// Validate facts before importing
		if validationResult, err := m.ValidateFacts(ctx, facts); err != nil {
			m.logger.Error("Failed to validate facts during import", err, map[string]interface{}{"machine_id": machineID})
			failed++
			continue
		} else if !validationResult.Valid {
			m.logger.Error("Facts validation failed during import", fmt.Errorf("validation errors: %v", validationResult.Errors), map[string]interface{}{"machine_id": machineID})
			failed++
			continue
		}

		// Store facts
		if err := m.StoreFacts(ctx, machineID, facts); err != nil {
			m.logger.Error("Failed to store facts during import", err, map[string]interface{}{"machine_id": machineID})
			failed++
			continue
		}

		imported++
	}

	m.logger.Info("Successfully imported facts", map[string]interface{}{"imported": imported, "failed": failed})

	return nil
}

// CollectAndStoreFacts collects and stores facts for a machine
func (m *Manager) CollectAndStoreFacts(ctx context.Context, machine *spookytypes.Machine) error {
	// Collect facts
	facts, err := m.CollectFacts(ctx, machine)
	if err != nil {
		return err
	}

	// Store facts
	return m.StoreFacts(ctx, facts.MachineID, facts)
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
	if memoryStorage, ok := m.storage.(*MemoryFactStorage); ok {
		return memoryStorage.GetStats()
	}

	return map[string]interface{}{
		"total_entries": 0,
		"total_size":    0,
		"storage_type":  "unknown",
	}, nil
}

// Close closes the manager and underlying storage
func (m *Manager) Close() error {
	if m.storage != nil {
		return m.storage.Close()
	}
	return nil
}

// FactExport represents exported facts data
type FactExport struct {
	ExportedAt   time.Time                  `json:"exported_at"`
	Format       string                     `json:"format"`
	MachineCount int                        `json:"machine_count"`
	Facts        map[string]*FactCollection `json:"facts"`
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

// exportToHCL exports facts to HCL format
func (m *Manager) exportToHCL(export *FactExport) ([]byte, error) {
	// Create HCL file
	f := hclwrite.NewFile()

	// Create root body
	rootBody := f.Body()

	// Add exported_at
	rootBody.SetAttributeValue("exported_at", cty.StringVal(export.ExportedAt.Format(time.RFC3339)))

	// Add format
	rootBody.SetAttributeValue("format", cty.StringVal(export.Format))

	// Add machine_count
	rootBody.SetAttributeValue("machine_count", cty.NumberIntVal(int64(export.MachineCount)))

	// Add facts block
	factsBlock := rootBody.AppendNewBlock("facts", nil)
	factsBody := factsBlock.Body()

	// Add each machine's facts
	for machineID, factCollection := range export.Facts {
		machineBlock := factsBody.AppendNewBlock("machine", []string{machineID})
		machineBody := machineBlock.Body()

		// Add machine_id
		machineBody.SetAttributeValue("machine_id", cty.StringVal(factCollection.MachineID))

		// Add collected_at
		machineBody.SetAttributeValue("collected_at", cty.StringVal(factCollection.CollectedAt.Format(time.RFC3339)))

		// Add facts data (simplified for now)
		if factCollection.Facts != nil {
			factsData := machineBody.AppendNewBlock("facts_data", nil)
			factsDataBody := factsData.Body()

			// Add system facts if available
			if factCollection.Facts.System != nil && factCollection.Facts.System.OS != nil {
				factsDataBody.SetAttributeValue("os_name", cty.StringVal(factCollection.Facts.System.OS.Name))
				if factCollection.Facts.System.OS.Version != "" {
					factsDataBody.SetAttributeValue("os_version", cty.StringVal(factCollection.Facts.System.OS.Version))
				}
			}
		}
	}

	return f.Bytes(), nil
}
