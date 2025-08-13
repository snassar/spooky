// Package facts provides fact collection, storage, and management functionality.
package facts

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	spookytypes "spooky/internal/types"
	spookytypesschemas "spooky/internal/types/schemas"
)

// Integration implements FactsIntegration interface
type Integration struct {
	manager FactManager
}

// NewIntegration creates a new facts integration
func NewIntegration(manager FactManager) *Integration {
	return &Integration{
		manager: manager,
	}
}

// GetManager returns the underlying fact manager
func (i *Integration) GetManager() interface{} {
	return i.manager
}

// CollectFacts collects facts from the given source
func (i *Integration) CollectFacts(ctx context.Context, source string) (interface{}, error) {
	// Create a placeholder machine for local fact collection
	machine := &spookytypes.Machine{
		Hostname: source,
		Host:     source,
		Port:     22,
		User:     "root",
	}

	// Collect facts using the manager
	facts, err := i.manager.CollectFacts(ctx, machine)
	if err != nil {
		return nil, fmt.Errorf("failed to collect facts from %s: %w", source, err)
	}

	return facts, nil
}

// StoreFacts stores facts in the given storage
func (i *Integration) StoreFacts(ctx context.Context, facts interface{}, storage spookytypes.FactStorage) error {
	if facts == nil {
		return fmt.Errorf("facts cannot be nil")
	}

	// Type assert to concrete type
	factCollection, ok := facts.(*FactCollection)
	if !ok {
		return fmt.Errorf("invalid facts type")
	}

	// Store facts using the manager
	if err := i.manager.StoreFacts(ctx, factCollection.MachineID, factCollection); err != nil {
		return fmt.Errorf("failed to store facts for %s: %w", factCollection.MachineID, err)
	}

	return nil
}

// LoadFacts loads facts from the given storage
func (i *Integration) LoadFacts(ctx context.Context, storage spookytypes.FactStorage) (interface{}, error) {
	// List all machine IDs with stored facts
	machineIDs, err := i.manager.ListFacts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list facts: %w", err)
	}

	if len(machineIDs) == 0 {
		return nil, nil
	}

	// For now, return the first machine's facts
	// In a real implementation, this would aggregate all facts
	if len(machineIDs) > 0 {
		facts, err := i.manager.GetFacts(ctx, machineIDs[0])
		if err != nil {
			return nil, fmt.Errorf("failed to get facts for %s: %w", machineIDs[0], err)
		}
		return facts, nil
	}

	return nil, nil
}

// ValidateFacts validates facts
func (i *Integration) ValidateFacts(ctx context.Context, facts interface{}) (*spookytypes.ValidationResult, error) {
	if facts == nil {
		return &spookytypes.ValidationResult{
			Valid:    false,
			Errors:   []spookytypesschemas.SchemaError{{Message: "facts cannot be nil"}},
			Warnings: []spookytypesschemas.SchemaError{},
		}, nil
	}

	// Type assert to concrete type
	factCollection, ok := facts.(*FactCollection)
	if !ok {
		return &spookytypes.ValidationResult{
			Valid:    false,
			Errors:   []spookytypesschemas.SchemaError{{Message: "invalid facts type"}},
			Warnings: []spookytypesschemas.SchemaError{},
		}, nil
	}

	// Validate facts using the manager
	result, err := i.manager.ValidateFacts(ctx, factCollection)
	if err != nil {
		return &spookytypes.ValidationResult{
			Valid:    false,
			Errors:   []spookytypesschemas.SchemaError{{Message: fmt.Sprintf("validation failed: %v", err)}},
			Warnings: []spookytypesschemas.SchemaError{},
		}, nil
	}

	return result, nil
}

// ExportFacts exports facts to the given format
func (i *Integration) ExportFacts(ctx context.Context, facts interface{}, format string, outputPath string) error {
	if facts == nil {
		return fmt.Errorf("facts cannot be nil")
	}

	// Type assert to concrete type
	factCollection, ok := facts.(*FactCollection)
	if !ok {
		return fmt.Errorf("invalid facts type")
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	switch format {
	case "json":
		return i.exportToJSON(factCollection, outputPath)
	case "hcl":
		return i.exportToHCL(factCollection, outputPath)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

// exportToJSON exports facts to JSON format
func (i *Integration) exportToJSON(facts *FactCollection, outputPath string) error {
	data, err := json.MarshalIndent(facts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal facts to JSON: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}

// exportToHCL exports facts to HCL format
func (i *Integration) exportToHCL(facts *FactCollection, outputPath string) error {
	// For now, export as JSON since HCL export is not implemented
	// In a real implementation, this would convert to HCL format
	return i.exportToJSON(facts, outputPath)
}
