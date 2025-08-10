package processors

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"

	spookyfactstypes "spooky/internal/types/facts"
	spookylogging "spooky/internal/logging"
)

// ProcessorManager implements Processor interface
type ProcessorManager struct {
	logger spookylogging.Logger
}

// NewProcessorManager creates a new processor manager
func NewProcessorManager(logger spookylogging.Logger) *ProcessorManager {
	return &ProcessorManager{
		logger: logger,
	}
}

// ExportToJSON exports fact collections to JSON format
func (m *ProcessorManager) ExportToJSON(collections []*spookyfactstypes.FactCollection, w io.Writer) error {
	// Convert to export format
	exportData := make(map[string]interface{})
	for _, collection := range collections {
		machineID := ""
		if machineIDFact, exists := collection.Facts["machine_id"]; exists {
			if id, ok := machineIDFact.Value.(string); ok {
				machineID = id
			}
		}
		if machineID == "" {
			machineID = collection.Server
		}

		exportData[machineID] = map[string]interface{}{
			"server":    collection.Server,
			"timestamp": collection.Timestamp.Format(time.RFC3339),
			"facts":     collection.Facts,
		}
	}

	// Export as JSON
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(exportData)
}

// ExportToHCL exports fact collections to HCL format
func (m *ProcessorManager) ExportToHCL(collections []*spookyfactstypes.FactCollection, w io.Writer) error {
	// Generate HCL content according to schema
	hclContent := "facts = [\n"

	for _, collection := range collections {
		machineID := ""
		if machineIDFact, exists := collection.Facts["machine_id"]; exists {
			if id, ok := machineIDFact.Value.(string); ok {
				machineID = id
			}
		}
		if machineID == "" {
			machineID = collection.Server
		}

		hclContent += fmt.Sprintf("  {\n")
		hclContent += fmt.Sprintf("    machine_id = \"%s\"\n", machineID)
		hclContent += fmt.Sprintf("    server = \"%s\"\n", collection.Server)
		hclContent += fmt.Sprintf("    timestamp = \"%s\"\n", collection.Timestamp.Format(time.RFC3339))

		hclContent += "    facts = {\n"
		for key, fact := range collection.Facts {
			hclContent += fmt.Sprintf("      \"%s\" = %s\n", key, formatHCLValue(fact.Value))
		}
		hclContent += "    }\n"

		hclContent += "  },\n"
	}

	hclContent += "]\n"

	// Write HCL content
	_, err := w.Write([]byte(hclContent))
	return err
}

// ImportFromJSON imports fact collections from JSON format
func (m *ProcessorManager) ImportFromJSON(r io.Reader) ([]*spookyfactstypes.FactCollection, error) {
	// Read and parse JSON
	var data map[string]interface{}
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Convert to fact collections
	collections := make([]*spookyfactstypes.FactCollection, 0, len(data))
	for machineID, machineData := range data {
		collection, err := convertToFactCollection(machineID, machineData)
		if err != nil {
			return nil, fmt.Errorf("failed to convert fact collection: %w", err)
		}
		collections = append(collections, collection)
	}

	return collections, nil
}

// ImportFromHCL imports fact collections from HCL format
func (m *ProcessorManager) ImportFromHCL(r io.Reader) ([]*spookyfactstypes.FactCollection, error) {
	// Read HCL content
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read HCL content: %w", err)
	}

	// Parse HCL
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(content, "facts.hcl")
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL: %w", diags)
	}

	// Parse HCL to fact collections
	collections, err := parseHCLToFactCollections(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HCL to fact collections: %w", err)
	}

	return collections, nil
}

// Helper functions
func formatHCLValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", v)
	case int, int64:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case []interface{}:
		items := make([]string, len(v))
		for i, item := range v {
			items[i] = formatHCLValue(item)
		}
		return fmt.Sprintf("[%s]", strings.Join(items, ", "))
	case map[string]interface{}:
		pairs := make([]string, 0, len(v))
		for key, val := range v {
			pairs = append(pairs, fmt.Sprintf("\"%s\" = %s", key, formatHCLValue(val)))
		}
		return fmt.Sprintf("{\n      %s\n    }", strings.Join(pairs, "\n      "))
	default:
		return fmt.Sprintf("%v", v)
	}
}

func convertToFactCollection(machineID string, data interface{}) (*spookyfactstypes.FactCollection, error) {
	// Implementation would convert map to FactCollection struct
	// This is a placeholder - actual implementation would parse the data structure
	return nil, fmt.Errorf("conversion not yet implemented")
}

func parseHCLToFactCollections(file *hcl.File) ([]*spookyfactstypes.FactCollection, error) {
	// Implementation would traverse the HCL AST and convert to FactCollection structs
	// This is a placeholder - actual implementation would parse the HCL structure
	return nil, fmt.Errorf("HCL parsing not yet implemented")
}
