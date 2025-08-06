package facts

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	spookyfactstypes "spooky/internal/facts/types"
	spookyschemas "spooky/internal/schemas"
	spookystorage "spooky/internal/storage"
)

// Exporter handles exporting facts to various formats
type Exporter struct {
	storage spookystorage.FactStorage
}

// NewExporter creates a new facts exporter
func NewExporter(storage spookystorage.FactStorage) *Exporter {
	return &Exporter{
		storage: storage,
	}
}

// ExportToJSON exports facts to JSON format
func (e *Exporter) ExportToJSON(w io.Writer, query *spookystorage.FactQuery) error {
	// Query facts from storage
	collections, err := e.storage.QueryFactCollections(query)
	if err != nil {
		return fmt.Errorf("failed to query facts: %w", err)
	}

	// Validate against schema
	validator := spookyschemas.NewSchemaValidator()
	if err := validator.LoadSchema(spookyschemas.SchemaTypeFactsStructure); err != nil {
		return fmt.Errorf("failed to load facts schema: %w", err)
	}

	// Create unified export structure
	exportData := map[string]interface{}{
		"metadata": map[string]interface{}{
			"exported_at":   time.Now().Format(time.RFC3339),
			"project_id":    "current-project", // TODO: Get from context
			"export_format": "json",
			"version":       "1.0",
		},
		"global_facts":  []interface{}{},
		"project_facts": []interface{}{},
	}

	// Process collections and separate global vs project facts
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

		// Determine if this is a global fact or project fact based on content
		isGlobalFact := e.isGlobalFact(collection)

		// Filter out custom facts from export
		filteredFacts := e.filterCustomFacts(collection.Facts)

		factEntry := map[string]interface{}{
			"machine_id":   machineID,
			"collected_at": collection.Timestamp.Format(time.RFC3339),
			"ttl":          "24h",
			"facts":        filteredFacts,
		}

		if isGlobalFact {
			exportData["global_facts"] = append(exportData["global_facts"].([]interface{}), factEntry)
		} else {
			// Add project_id for project facts
			factEntry["project_id"] = "current-project" // TODO: Get from context
			exportData["project_facts"] = append(exportData["project_facts"].([]interface{}), factEntry)
		}
	}

	// Validate against schema
	if err := validator.ValidateData(exportData, "facts_export"); err != nil {
		return fmt.Errorf("facts export validation failed: %w", err)
	}

	// Write JSON
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(exportData)
}

// ExportToHCL exports facts to HCL format
func (e *Exporter) ExportToHCL(w io.Writer, query *spookystorage.FactQuery) error {
	// Query facts from storage
	collections, err := e.storage.QueryFactCollections(query)
	if err != nil {
		return fmt.Errorf("failed to query facts: %w", err)
	}

	// Validate against schema
	validator := spookyschemas.NewSchemaValidator()
	if err := validator.LoadSchema(spookyschemas.SchemaTypeFactsStructure); err != nil {
		return fmt.Errorf("failed to load facts schema: %w", err)
	}

	// Generate HCL content according to unified schema
	hclContent := "facts_export = {\n"

	// Add metadata
	hclContent += "  metadata = {\n"
	hclContent += fmt.Sprintf("    exported_at = \"%s\"\n", time.Now().Format(time.RFC3339))
	hclContent += "    project_id = \"current-project\"\n" // TODO: Get from context
	hclContent += "    export_format = \"hcl\"\n"
	hclContent += "    version = \"1.0\"\n"
	hclContent += "  }\n"

	// Separate global and project facts
	var globalFacts []*spookyfactstypes.FactCollection
	var projectFacts []*spookyfactstypes.FactCollection

	for _, collection := range collections {
		if e.isGlobalFact(collection) {
			globalFacts = append(globalFacts, collection)
		} else {
			projectFacts = append(projectFacts, collection)
		}
	}

	// Export global facts
	hclContent += "  global_facts = [\n"
	for _, collection := range globalFacts {
		machineID := ""
		if machineIDFact, exists := collection.Facts["machine_id"]; exists {
			if id, ok := machineIDFact.Value.(string); ok {
				machineID = id
			}
		}
		if machineID == "" {
			machineID = collection.Server
		}

		hclContent += fmt.Sprintf("    {\n")
		hclContent += fmt.Sprintf("      machine_id = \"%s\"\n", machineID)
		hclContent += fmt.Sprintf("      collected_at = \"%s\"\n", collection.Timestamp.Format(time.RFC3339))
		hclContent += fmt.Sprintf("      ttl = \"24h\"\n")
		hclContent += "      facts = {\n"

		// Filter out custom facts and export actual facts
		filteredFacts := e.filterCustomFacts(collection.Facts)
		for key, fact := range filteredFacts {
			hclContent += fmt.Sprintf("        \"%s\" = %s\n", key, formatHCLValue(fact.Value))
		}

		hclContent += "      }\n"
		hclContent += "    },\n"
	}
	hclContent += "  ]\n"

	// Export project facts
	hclContent += "  project_facts = [\n"
	for _, collection := range projectFacts {
		machineID := ""
		if machineIDFact, exists := collection.Facts["machine_id"]; exists {
			if id, ok := machineIDFact.Value.(string); ok {
				machineID = id
			}
		}
		if machineID == "" {
			machineID = collection.Server
		}

		hclContent += fmt.Sprintf("    {\n")
		hclContent += fmt.Sprintf("      machine_id = \"%s\"\n", machineID)
		hclContent += fmt.Sprintf("      project_id = \"current-project\"\n") // TODO: Get from context
		hclContent += fmt.Sprintf("      collected_at = \"%s\"\n", collection.Timestamp.Format(time.RFC3339))
		hclContent += fmt.Sprintf("      ttl = \"24h\"\n")
		hclContent += "      facts = {\n"

		// Filter out custom facts and export actual facts
		filteredFacts := e.filterCustomFacts(collection.Facts)
		for key, fact := range filteredFacts {
			hclContent += fmt.Sprintf("        \"%s\" = %s\n", key, formatHCLValue(fact.Value))
		}

		hclContent += "      }\n"
		hclContent += "    },\n"
	}
	hclContent += "  ]\n"

	hclContent += "}\n"

	// Write HCL content
	_, err = w.Write([]byte(hclContent))
	return err
}

// formatHCLValue formats a value for HCL output
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

// ExportToFile exports facts to a file in the specified format
func (e *Exporter) ExportToFile(outputPath, format string, query *spookystorage.FactQuery) error {
	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Export based on format
	switch format {
	case "json":
		return e.ExportToJSON(file, query)
	case "hcl":
		return e.ExportToHCL(file, query)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

// ExportProjectFacts exports facts for a specific project
func (e *Exporter) ExportProjectFacts(projectPath, outputPath, format string) error {
	// Determine facts storage path for the project
	factsPath := filepath.Join(projectPath, "facts.db")

	// Check if facts.db directory exists
	if _, err := os.Stat(factsPath); os.IsNotExist(err) {
		return fmt.Errorf("facts database not found at %s", factsPath)
	}

	// Check if the directory has BadgerDB files BEFORE trying to open the database
	entries, err := os.ReadDir(factsPath)
	if err != nil {
		return fmt.Errorf("failed to read facts directory: %w", err)
	}

	// Check if directory has BadgerDB files
	hasBadgerFiles := false
	for _, entry := range entries {
		if !entry.IsDir() && (entry.Name() == "MANIFEST" || entry.Name() == "KEYREGISTRY" || strings.HasSuffix(entry.Name(), ".vlog")) {
			hasBadgerFiles = true
			break
		}
	}

	if !hasBadgerFiles {
		return fmt.Errorf("facts database exists but contains no data at %s", factsPath)
	}

	// Open storage for the project (read-only)
	storage, err := spookystorage.OpenFactStorage(spookystorage.StorageOptions{
		Type: spookystorage.StorageTypeBadger,
		Path: factsPath,
	})
	if err != nil {
		return fmt.Errorf("failed to open facts storage: %w", err)
	}
	defer storage.Close()

	// Create exporter
	exporter := NewExporter(storage)

	// Export all facts
	query := &spookystorage.FactQuery{}
	return exporter.ExportToFile(outputPath, format, query)
}

// ExportMachineFacts exports facts for a specific machine
func (e *Exporter) ExportMachineFacts(projectPath, machineID, outputPath, format string) error {
	// Determine facts storage path for the project
	factsPath := filepath.Join(projectPath, "facts.db")

	// Check if facts.db directory exists
	if _, err := os.Stat(factsPath); os.IsNotExist(err) {
		return fmt.Errorf("facts database not found at %s", factsPath)
	}

	// Check if the directory has BadgerDB files BEFORE trying to open the database
	entries, err := os.ReadDir(factsPath)
	if err != nil {
		return fmt.Errorf("failed to read facts directory: %w", err)
	}

	// Check if directory has BadgerDB files
	hasBadgerFiles := false
	for _, entry := range entries {
		if !entry.IsDir() && (entry.Name() == "MANIFEST" || entry.Name() == "KEYREGISTRY" || strings.HasSuffix(entry.Name(), ".vlog")) {
			hasBadgerFiles = true
			break
		}
	}

	if !hasBadgerFiles {
		return fmt.Errorf("facts database exists but contains no data at %s", factsPath)
	}

	// Open storage for the project (read-only)
	storage, err := spookystorage.OpenFactStorage(spookystorage.StorageOptions{
		Type: spookystorage.StorageTypeBadger,
		Path: factsPath,
	})
	if err != nil {
		return fmt.Errorf("failed to open facts storage: %w", err)
	}
	defer storage.Close()

	// Create exporter
	exporter := NewExporter(storage)

	// Export facts for specific machine
	query := &spookystorage.FactQuery{
		MachineName: machineID, // Use machine name for querying
	}
	return exporter.ExportToFile(outputPath, format, query)
}

// ExportFilteredFacts exports facts based on a filter
func (e *Exporter) ExportFilteredFacts(projectPath, filter, outputPath, format string) error {
	// Determine facts storage path for the project
	factsPath := filepath.Join(projectPath, "facts.db")

	// Check if facts.db directory exists
	if _, err := os.Stat(factsPath); os.IsNotExist(err) {
		return fmt.Errorf("facts database not found at %s", factsPath)
	}

	// Check if the directory has BadgerDB files BEFORE trying to open the database
	entries, err := os.ReadDir(factsPath)
	if err != nil {
		return fmt.Errorf("failed to read facts directory: %w", err)
	}

	// Check if directory has BadgerDB files
	hasBadgerFiles := false
	for _, entry := range entries {
		if !entry.IsDir() && (entry.Name() == "MANIFEST" || entry.Name() == "KEYREGISTRY" || strings.HasSuffix(entry.Name(), ".vlog")) {
			hasBadgerFiles = true
			break
		}
	}

	if !hasBadgerFiles {
		return fmt.Errorf("facts database exists but contains no data at %s", factsPath)
	}

	// Open storage for the project (read-only)
	storage, err := spookystorage.OpenFactStorage(spookystorage.StorageOptions{
		Type: spookystorage.StorageTypeBadger,
		Path: factsPath,
	})
	if err != nil {
		return fmt.Errorf("failed to open facts storage: %w", err)
	}
	defer storage.Close()

	// Create exporter
	exporter := NewExporter(storage)

	// Parse filter and create query
	query, err := ParseFilter(filter)
	if err != nil {
		return fmt.Errorf("failed to parse filter: %w", err)
	}

	return exporter.ExportToFile(outputPath, format, query)
}

// ParseFilter parses a filter string into a FactQuery
func ParseFilter(filter string) (*spookystorage.FactQuery, error) {
	// This is a simple filter parser - can be enhanced for more complex queries
	query := &spookystorage.FactQuery{}

	// For now, just support basic machine name filtering
	// TODO: Implement more sophisticated filter parsing
	if filter != "" {
		query.MachineName = filter
	}

	return query, nil
}

// isGlobalFact determines if a fact collection contains global facts
func (e *Exporter) isGlobalFact(collection *spookyfactstypes.FactCollection) bool {
	// Check for system-related facts that indicate global facts
	systemKeys := []string{"os", "hardware", "network", "cpu", "memory", "disk", "load", "processes"}

	for _, key := range systemKeys {
		if _, exists := collection.Facts[key]; exists {
			return true
		}
	}

	// Check for enhanced facts that indicate global facts
	enhancedKeys := []string{"virtualization", "package_manager", "service_manager", "selinux", "ssh_keys", "bios", "sensors", "docker"}

	for _, key := range enhancedKeys {
		if _, exists := collection.Facts[key]; exists {
			return true
		}
	}

	// If no system or enhanced facts found, assume it's a project fact
	return false
}

// filterCustomFacts removes custom facts from the facts map
// Custom facts are those that come from /etc/spooky/facts.hcl and should not be exported
func (e *Exporter) filterCustomFacts(facts map[string]*spookyfactstypes.Fact) map[string]*spookyfactstypes.Fact {
	filtered := make(map[string]*spookyfactstypes.Fact)

	for key, fact := range facts {
		// Skip custom facts that come from /etc/spooky/facts.hcl
		if fact.Source == string(spookyfactstypes.SourceCustom) {
			continue
		}

		// Skip facts that are explicitly marked as custom
		if strings.HasPrefix(key, "custom_") {
			continue
		}

		filtered[key] = fact
	}

	return filtered
}
