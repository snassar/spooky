package facts

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"spooky/internal/facts/types"
	"spooky/internal/schemas"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// CustomFactsCollector collects custom facts from /etc/spooky/facts.hcl file
type CustomFactsCollector struct {
	factsFile string
}

// NewCustomFactsCollector creates a new custom facts collector
func NewCustomFactsCollector(factsFile string) (*CustomFactsCollector, error) {
	if factsFile == "" {
		factsFile = "/etc/spooky/facts.hcl"
	}

	// Validate facts file
	if err := validateFactsFile(factsFile); err != nil {
		return nil, fmt.Errorf("invalid facts file: %w", err)
	}

	return &CustomFactsCollector{
		factsFile: factsFile,
	}, nil
}

// Collect collects custom facts from the facts.hcl file
func (cfc *CustomFactsCollector) Collect(server string) (*types.FactCollection, error) {
	collection := &types.FactCollection{
		Server:      server,
		Timestamp:   time.Now(),
		Facts:       make(map[string]*types.Fact),
		CustomFacts: make(map[string]map[string]interface{}),
	}

	// Check if facts file exists
	if _, err := os.Stat(cfc.factsFile); os.IsNotExist(err) {
		// File doesn't exist, return empty collection
		return collection, nil
	}

	// Read and parse custom facts file
	customFacts, err := cfc.readCustomFactsFile(server)
	if err != nil {
		return nil, fmt.Errorf("failed to read custom facts file: %w", err)
	}

	// Validate custom facts
	if err := cfc.validateCustomFacts(customFacts); err != nil {
		return nil, fmt.Errorf("custom facts validation failed: %w", err)
	}

	// Add custom facts to collection
	cfc.addCustomFactsToCollection(collection, customFacts)

	return collection, nil
}

// CollectSpecific collects specific custom facts by keys
func (cfc *CustomFactsCollector) CollectSpecific(server string, keys []string) (*types.FactCollection, error) {
	collection := &types.FactCollection{
		Server:      server,
		Timestamp:   time.Now(),
		Facts:       make(map[string]*types.Fact),
		CustomFacts: make(map[string]map[string]interface{}),
	}

	// Check if facts file exists
	if _, err := os.Stat(cfc.factsFile); os.IsNotExist(err) {
		return collection, nil
	}

	// Read all custom facts
	customFacts, err := cfc.readCustomFactsFile(server)
	if err != nil {
		return nil, fmt.Errorf("failed to read custom facts file: %w", err)
	}

	// Validate original facts first
	if err := cfc.validateCustomFacts(customFacts); err != nil {
		return nil, fmt.Errorf("custom facts validation failed: %w", err)
	}

	// Filter to only requested keys
	filteredFacts := make(map[string]interface{})
	requestedKeys := make(map[string]bool)
	for _, key := range keys {
		requestedKeys[key] = true
	}

	for key, value := range customFacts {
		// Skip metadata keys
		if strings.HasPrefix(key, "_") {
			continue
		}
		if requestedKeys[key] {
			filteredFacts[key] = value
		}
	}

	// Add filtered facts to collection
	cfc.addCustomFactsToCollection(collection, filteredFacts)

	return collection, nil
}

// GetFact retrieves a specific custom fact
func (cfc *CustomFactsCollector) GetFact(server, key string) (*types.Fact, error) {
	// Use new format (custom.key) only
	factKey := key
	if strings.HasPrefix(key, "custom.") {
		factKey = strings.TrimPrefix(key, "custom.")
	}

	// Check if facts file exists
	if _, err := os.Stat(cfc.factsFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("custom facts file does not exist: %s", cfc.factsFile)
	}

	// Read all custom facts
	customFacts, err := cfc.readCustomFactsFile(server)
	if err != nil {
		return nil, fmt.Errorf("failed to read custom facts file: %w", err)
	}

	// Check if the specific fact exists
	if _, exists := customFacts[factKey]; !exists {
		return nil, fmt.Errorf("custom fact %s not found", factKey)
	}

	// Create a fact collection with just this fact
	collection := &types.FactCollection{
		Server:      server,
		Timestamp:   time.Now(),
		Facts:       make(map[string]*types.Fact),
		CustomFacts: make(map[string]map[string]interface{}),
	}

	// Add the specific fact to collection
	collection.CustomFacts["facts"] = map[string]interface{}{
		factKey: customFacts[factKey],
	}

	// Convert to fact with proper key format
	fact := &types.Fact{
		Key:       "custom." + factKey,
		Value:     customFacts[factKey],
		Source:    "custom",
		Server:    server,
		Timestamp: time.Now(),
		TTL:       time.Hour, // Custom facts default TTL
		Metadata: map[string]interface{}{
			"file_path": cfc.factsFile,
			"fact_type": "custom",
			"filename":  "facts",
		},
	}

	return fact, nil
}

// readCustomFactsFile reads and parses the custom facts file
func (cfc *CustomFactsCollector) readCustomFactsFile(server string) (map[string]interface{}, error) {
	// Read file content
	content, err := os.ReadFile(cfc.factsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse HCL content
	customFacts, err := cfc.parseHCLContent(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HCL content: %w", err)
	}

	// Add file metadata
	if customFacts == nil {
		customFacts = make(map[string]interface{})
	}
	customFacts["_file_path"] = cfc.factsFile
	customFacts["_server"] = server
	customFacts["_collected_at"] = time.Now().Format(time.RFC3339)

	return customFacts, nil
}

// readCustomFactFile reads and parses a custom fact file
func (cfc *CustomFactsCollector) readCustomFactFile(filePath, server string) (map[string]interface{}, error) {
	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse HCL content
	customFact, err := cfc.parseHCLContent(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HCL content: %w", err)
	}

	// Add file metadata
	if customFact == nil {
		customFact = make(map[string]interface{})
	}
	customFact["_file_path"] = filePath
	customFact["_server"] = server
	customFact["_collected_at"] = time.Now().Format(time.RFC3339)

	return customFact, nil
}

// parseHCLContent parses HCL content from custom fact file using proper HCL library
func (cfc *CustomFactsCollector) parseHCLContent(content []byte) (map[string]interface{}, error) {
	// Use proper HCL parser
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(content, "custom-facts.hcl")
	if diags.HasErrors() {
		return nil, fmt.Errorf("HCL parsing failed: %v", diags)
	}

	// Extract facts from HCL body
	result := make(map[string]interface{})
	attrs, diags := file.Body.JustAttributes()
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL attributes: %v", diags)
	}

	for name, attr := range attrs {
		value, err := extractHCLValue(attr.Expr)
		if err != nil {
			return nil, fmt.Errorf("failed to extract value for %s: %w", name, err)
		}
		result[name] = value
	}

	return result, nil
}

// extractHCLValue extracts a value from an HCL expression
func extractHCLValue(expr hcl.Expression) (interface{}, error) {
	val, diags := expr.Value(nil)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to evaluate expression: %v", diags)
	}

	// Convert HCL value to Go interface{}
	if val.IsNull() {
		return nil, nil
	}

	switch {
	case val.Type().IsPrimitiveType() && val.Type().FriendlyName() == "string":
		return val.AsString(), nil
	case val.Type().IsPrimitiveType() && val.Type().FriendlyName() == "bool":
		return val.True(), nil
	case val.Type().IsPrimitiveType() && val.Type().FriendlyName() == "number":
		if val.CanIterateElements() {
			// Handle list length
			return val.LengthInt(), nil
		}
		// Try to get as float64
		if f := val.AsBigFloat(); f != nil {
			if f64, _ := f.Float64(); f64 != 0 {
				return f64, nil
			}
		}
		return val.AsBigFloat(), nil
	case val.Type().IsListType():
		var result []interface{}
		for _, item := range val.AsValueSlice() {
			if item.Type().IsPrimitiveType() && item.Type().FriendlyName() == "string" {
				result = append(result, item.AsString())
			} else if item.Type().IsPrimitiveType() && item.Type().FriendlyName() == "number" {
				if f := item.AsBigFloat(); f != nil {
					if f64, _ := f.Float64(); f64 != 0 {
						result = append(result, f64)
					}
				}
			} else {
				result = append(result, item)
			}
		}
		return result, nil
	default:
		return val, nil
	}
}

// validateCustomFacts validates custom facts using schema system
func (cfc *CustomFactsCollector) validateCustomFacts(customFacts map[string]interface{}) error {
	validator := schemas.NewSchemaValidator()
	if err := validator.LoadSchema(schemas.SchemaTypeCustomFactsHCL); err != nil {
		return fmt.Errorf("failed to load custom facts schema: %w", err)
	}

	if err := validator.ValidateData(customFacts, "custom-facts-hcl"); err != nil {
		return fmt.Errorf("custom facts validation failed: %w", err)
	}
	return nil
}

// validateCustomFact validates custom fact against schema
func (cfc *CustomFactsCollector) validateCustomFact(customFact map[string]interface{}) error {
	// Validate against custom-facts-hcl.hcl schema
	// For now, implement basic validation rules

	// Check required fields from schema
	requiredFields := []string{"app_name", "app_version"}
	for _, field := range requiredFields {
		if _, exists := customFact[field]; !exists {
			return fmt.Errorf("required field %s is missing", field)
		}
	}

	// Validate field types and constraints
	for key, value := range customFact {
		// Skip metadata fields
		if strings.HasPrefix(key, "_") {
			continue
		}

		// Validate key format
		keyPattern := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
		if !keyPattern.MatchString(key) {
			return fmt.Errorf("invalid key format: %s", key)
		}

		// Validate value is not nil
		if value == nil {
			return fmt.Errorf("value cannot be nil for key: %s", key)
		}

		// Validate specific fields based on schema
		switch key {
		case "app_name":
			if appName, ok := value.(string); ok {
				if appName == "" {
					return fmt.Errorf("app_name cannot be empty")
				}
			} else {
				return fmt.Errorf("app_name must be a string")
			}
		case "app_version":
			if appVersion, ok := value.(string); ok {
				if appVersion == "" {
					return fmt.Errorf("app_version cannot be empty")
				}
			} else {
				return fmt.Errorf("app_version must be a string")
			}
		case "config_path":
			if configPath, ok := value.(string); ok {
				if configPath != "" {
					// Validate path format
					if !strings.HasPrefix(configPath, "/") {
						return fmt.Errorf("config_path must be an absolute path")
					}
				}
			} else {
				return fmt.Errorf("config_path must be a string")
			}
		}
	}

	return nil
}

// addCustomFactToCollection adds custom fact to collection
func (cfc *CustomFactsCollector) addCustomFactToCollection(collection *types.FactCollection, customFact map[string]interface{}, filePath string) {
	// Extract filename from path
	fileName := filepath.Base(filePath)
	factKey := strings.TrimSuffix(fileName, ".fact")

	// Store custom facts in the hierarchical structure
	// This preserves the file-based organization: custom.<filename>.<key>
	if collection.CustomFacts == nil {
		collection.CustomFacts = make(map[string]map[string]interface{})
	}

	collection.CustomFacts[factKey] = customFact

	// Also store a reference in the flat Facts map for backward compatibility
	// This allows existing code to still access custom facts
	fact := &types.Fact{
		Key:       "custom." + factKey,
		Value:     customFact,
		Source:    "custom",
		Server:    collection.Server,
		Timestamp: collection.Timestamp,
		TTL:       time.Hour, // Custom facts default TTL
		Metadata: map[string]interface{}{
			"file_path": filePath,
			"fact_type": "custom",
			"filename":  factKey,
		},
	}

	collection.Facts["custom."+factKey] = fact
}

// addCustomFactsToCollection adds custom facts to collection from the single facts.hcl file
func (cfc *CustomFactsCollector) addCustomFactsToCollection(collection *types.FactCollection, customFacts map[string]interface{}) {
	// Store custom facts in the hierarchical structure
	// Use "facts" as the filename since it's from /etc/spooky/facts.hcl
	collection.CustomFacts["facts"] = customFacts

	// Also store a reference in the flat Facts map for backward compatibility
	fact := &types.Fact{
		Key:       "custom.facts",
		Value:     customFacts,
		Source:    "custom",
		Server:    collection.Server,
		Timestamp: collection.Timestamp,
		TTL:       time.Hour, // Custom facts default TTL
		Metadata: map[string]interface{}{
			"file_path": cfc.factsFile,
			"fact_type": "custom",
			"filename":  "facts",
		},
	}

	collection.Facts["custom.facts"] = fact
}

// validateFactsFile validates the facts file
func validateFactsFile(factsFile string) error {
	if factsFile == "" {
		return fmt.Errorf("facts file cannot be empty")
	}

	// Check if file exists
	if _, err := os.Stat(factsFile); os.IsNotExist(err) {
		// File doesn't exist, that's okay for custom facts
		return nil
	}

	// Check if it's a file
	info, err := os.Stat(factsFile)
	if err != nil {
		return fmt.Errorf("failed to stat facts file: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("facts file path is a directory: %s", factsFile)
	}

	// Check if file is readable
	if _, err := os.ReadFile(factsFile); err != nil {
		return fmt.Errorf("facts file is not readable: %w", err)
	}

	return nil
}

// GetFactsFile returns the facts file path
func (cfc *CustomFactsCollector) GetFactsFile() string {
	return cfc.factsFile
}

// ListAvailableFacts lists available custom facts
func (cfc *CustomFactsCollector) ListAvailableFacts() ([]string, error) {
	var facts []string

	// Check if facts file exists
	if _, err := os.Stat(cfc.factsFile); os.IsNotExist(err) {
		return facts, nil // Return empty list if file doesn't exist
	}

	// Read and parse the facts file to get available keys
	customFacts, err := cfc.readCustomFactsFile("")
	if err != nil {
		return nil, fmt.Errorf("failed to read facts file: %w", err)
	}

	// Extract keys from the custom facts
	for key := range customFacts {
		// Skip metadata keys
		if !strings.HasPrefix(key, "_") {
			facts = append(facts, key)
		}
	}

	return facts, nil
}

// ValidateCustomFactsFile validates the custom facts file
func (cfc *CustomFactsCollector) ValidateCustomFactsFile() error {
	// Check if file exists
	if _, err := os.Stat(cfc.factsFile); os.IsNotExist(err) {
		return fmt.Errorf("custom facts file does not exist: %s", cfc.factsFile)
	}

	// Read and parse file
	content, err := os.ReadFile(cfc.factsFile)
	if err != nil {
		return fmt.Errorf("failed to read custom facts file: %w", err)
	}

	// Parse HCL content
	customFacts, err := cfc.parseHCLContent(content)
	if err != nil {
		return fmt.Errorf("failed to parse custom facts file: %w", err)
	}

	// Validate custom facts
	if err := cfc.validateCustomFacts(customFacts); err != nil {
		return fmt.Errorf("custom facts validation failed: %w", err)
	}

	return nil
}
