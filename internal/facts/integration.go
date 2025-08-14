// Package facts provides fact collection, storage, and management functionality.
package facts

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

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
	// Create a machine representation for local fact collection
	// The source parameter represents the local machine identifier
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

// StoreFacts stores facts in memory
func (i *Integration) StoreFacts(ctx context.Context, facts interface{}) error {
	if facts == nil {
		return fmt.Errorf("facts cannot be nil")
	}

	// Type assert to concrete type
	_, ok := facts.(*FactCollection)
	if !ok {
		return fmt.Errorf("invalid facts type")
	}

	// Facts are collected and exported directly - no storage needed
	return nil
}

// LoadFacts loads facts from memory
func (i *Integration) LoadFacts(ctx context.Context) (interface{}, error) {
	// Memory-only storage - no facts to load
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
	// Convert facts to HCL format
	hclContent := "fact_collection {\n"

	// Add fact collection metadata
	if facts.MachineID != "" {
		hclContent += fmt.Sprintf("  machine_id = \"%s\"\n", facts.MachineID)
	}
	hclContent += fmt.Sprintf("  collected_at = \"%s\"\n", facts.CollectedAt.Format(time.RFC3339))

	// Add facts block
	if facts.Facts != nil {
		hclContent += "  facts {\n"

		// Add system facts if available
		if facts.Facts.System != nil {
			hclContent += "    system {\n"
			if facts.Facts.System.OS != nil {
				hclContent += "      os {\n"
				hclContent += fmt.Sprintf("        name = \"%s\"\n", facts.Facts.System.OS.Name)
				hclContent += fmt.Sprintf("        version = \"%s\"\n", facts.Facts.System.OS.Version)
				hclContent += fmt.Sprintf("        arch = \"%s\"\n", facts.Facts.System.OS.Arch)
				hclContent += fmt.Sprintf("        kernel = \"%s\"\n", facts.Facts.System.OS.Kernel)
				if facts.Facts.System.OS.Platform != "" {
					hclContent += fmt.Sprintf("        platform = \"%s\"\n", facts.Facts.System.OS.Platform)
				}
				if facts.Facts.System.OS.Family != "" {
					hclContent += fmt.Sprintf("        family = \"%s\"\n", facts.Facts.System.OS.Family)
				}
				hclContent += "      }\n"
			}
			hclContent += "    }\n"
		}

		// Add custom facts if available
		if len(facts.Facts.Custom) > 0 {
			hclContent += "    custom {\n"
			for key, value := range facts.Facts.Custom {
				switch v := value.(type) {
				case string:
					hclContent += fmt.Sprintf("      %s = \"%s\"\n", key, v)
				case int, int64:
					hclContent += fmt.Sprintf("      %s = %v\n", key, v)
				case float64:
					hclContent += fmt.Sprintf("      %s = %f\n", key, v)
				case bool:
					hclContent += fmt.Sprintf("      %s = %t\n", key, v)
				default:
					hclContent += fmt.Sprintf("      %s = \"%v\"\n", key, v)
				}
			}
			hclContent += "    }\n"
		}

		hclContent += "  }\n"
	}

	// Add metadata if available
	if len(facts.Metadata) > 0 {
		hclContent += "  metadata {\n"
		for key, value := range facts.Metadata {
			switch v := value.(type) {
			case string:
				hclContent += fmt.Sprintf("    %s = \"%s\"\n", key, v)
			case int, int64:
				hclContent += fmt.Sprintf("    %s = %v\n", key, v)
			case float64:
				hclContent += fmt.Sprintf("    %s = %f\n", key, v)
			case bool:
				hclContent += fmt.Sprintf("    %s = %t\n", key, v)
			default:
				hclContent += fmt.Sprintf("    %s = \"%v\"\n", key, v)
			}
		}
		hclContent += "  }\n"
	}

	hclContent += "}\n"

	// Write HCL content to file
	if err := os.WriteFile(outputPath, []byte(hclContent), 0644); err != nil {
		return fmt.Errorf("failed to write HCL file: %w", err)
	}

	return nil
}
