package custom

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	spookyfactscollectors "spooky/internal/facts/collectors"
	spookyfactstypes "spooky/internal/types/facts"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// Collector collects custom facts from various sources
type Collector struct {
	spookyfactscollectors.BaseCollector
	customFactsPath string
	overrides       map[string]interface{}
}

// NewCollector creates a new custom fact collector
func NewCollector(customFactsPath string) *Collector {
	return &Collector{
		BaseCollector:   *spookyfactscollectors.NewBaseCollector(spookyfactstypes.SourceCustom, spookyfactstypes.MergePolicyMerge),
		customFactsPath: customFactsPath,
		overrides:       make(map[string]interface{}),
	}
}

// Collect gathers custom facts from the configured sources
func (c *Collector) Collect(server string) (*spookyfactstypes.FactCollection, error) {
	collection := &spookyfactstypes.FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*spookyfactstypes.Fact),
	}

	// Load custom facts from file
	if err := c.loadCustomFactsFromFile(collection); err != nil {
		return nil, fmt.Errorf("failed to load custom facts from file: %w", err)
	}

	// Apply overrides
	c.applyOverrides(collection)

	return collection, nil
}

// CollectSpecific collects only the specified custom facts
func (c *Collector) CollectSpecific(server string, keys []string) (*spookyfactstypes.FactCollection, error) {
	collection := &spookyfactstypes.FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*spookyfactstypes.Fact),
	}

	// Load all custom facts first
	if err := c.loadCustomFactsFromFile(collection); err != nil {
		return nil, fmt.Errorf("failed to load custom facts from file: %w", err)
	}

	// Filter to only requested keys
	filteredCollection := &spookyfactstypes.FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*spookyfactstypes.Fact),
	}

	for _, key := range keys {
		if fact, exists := collection.Facts[key]; exists {
			filteredCollection.Facts[key] = fact
		}
	}

	// Apply overrides
	c.applyOverrides(filteredCollection)

	return filteredCollection, nil
}

// GetFact retrieves a single custom fact
func (c *Collector) GetFact(server, key string) (*spookyfactstypes.Fact, error) {
	collection := &spookyfactstypes.FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*spookyfactstypes.Fact),
	}

	// Load custom facts from file
	if err := c.loadCustomFactsFromFile(collection); err != nil {
		return nil, fmt.Errorf("failed to load custom facts from file: %w", err)
	}

	// Check if fact exists
	fact, exists := collection.Facts[key]
	if !exists {
		return nil, fmt.Errorf("custom fact %s not found", key)
	}

	// Apply overrides if they exist for this key
	if override, hasOverride := c.overrides[key]; hasOverride {
		fact.Value = override
	}

	return fact, nil
}

// Validate validates the collector configuration
func (c *Collector) Validate() error {
	if c.customFactsPath == "" {
		return fmt.Errorf("custom facts path is not configured")
	}

	// Check if the custom facts file exists
	if _, err := os.Stat(c.customFactsPath); os.IsNotExist(err) {
		return fmt.Errorf("custom facts file does not exist: %s", c.customFactsPath)
	}

	return nil
}

// SetOverride sets an override value for a fact
func (c *Collector) SetOverride(key string, value interface{}) {
	c.overrides[key] = value
}

// ClearOverride clears an override for a fact
func (c *Collector) ClearOverride(key string) {
	delete(c.overrides, key)
}

// ClearAllOverrides clears all overrides
func (c *Collector) ClearAllOverrides() {
	c.overrides = make(map[string]interface{})
}

// GetOverrides returns all current overrides
func (c *Collector) GetOverrides() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range c.overrides {
		result[k] = v
	}
	return result
}

// loadCustomFactsFromFile loads custom facts from the configured file
func (c *Collector) loadCustomFactsFromFile(collection *spookyfactstypes.FactCollection) error {
	if c.customFactsPath == "" {
		return fmt.Errorf("custom facts path is not configured")
	}

	// Determine file type and load accordingly
	ext := strings.ToLower(filepath.Ext(c.customFactsPath))

	switch ext {
	case ".json":
		return c.loadFromJSON(collection)
	case ".hcl":
		return c.loadFromHCL(collection)
	default:
		return fmt.Errorf("unsupported file format: %s", ext)
	}
}

// loadFromJSON loads custom facts from a JSON file
func (c *Collector) loadFromJSON(collection *spookyfactstypes.FactCollection) error {
	file, err := os.Open(c.customFactsPath)
	if err != nil {
		return fmt.Errorf("failed to open custom facts file: %w", err)
	}
	defer file.Close()

	var customFacts spookyfactstypes.CustomFacts
	if err := json.NewDecoder(file).Decode(&customFacts); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Add custom facts to collection
	for key, value := range customFacts.Custom {
		c.createFact(collection, key, value)
	}

	// Add overrides to collection
	for key, value := range customFacts.Overrides {
		c.createFact(collection, key, value)
	}

	return nil
}

// loadFromHCL loads custom facts from an HCL file
func (c *Collector) loadFromHCL(collection *spookyfactstypes.FactCollection) error {
	// For now, we'll use a simple approach to parse HCL
	// In a real implementation, you would use the HCL parser
	content, err := os.ReadFile(c.customFactsPath)
	if err != nil {
		return fmt.Errorf("failed to read HCL file: %w", err)
	}

	// Parse HCL content (simplified - in real implementation use HCL parser)
	facts, err := c.parseHCLContent(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse HCL content: %w", err)
	}

	// Add facts to collection
	for key, value := range facts {
		c.createFact(collection, key, value)
	}

	return nil
}

// parseHCLContent parses HCL content and extracts facts using proper HCL library
func (c *Collector) parseHCLContent(content string) (map[string]interface{}, error) {
	// Use proper HCL parser
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL([]byte(content), "custom-facts.hcl")
	if diags.HasErrors() {
		return nil, fmt.Errorf("HCL parsing failed: %v", diags)
	}

	// Extract facts from HCL body
	facts := make(map[string]interface{})
	attrs, diags := file.Body.JustAttributes()
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL attributes: %v", diags)
	}

	for name, attr := range attrs {
		value, err := extractHCLValue(attr.Expr)
		if err != nil {
			return nil, fmt.Errorf("failed to extract value for %s: %w", name, err)
		}
		facts[name] = value
	}

	return facts, nil
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

// applyOverrides applies any configured overrides to the collection
func (c *Collector) applyOverrides(collection *spookyfactstypes.FactCollection) {
	for key, value := range c.overrides {
		if fact, exists := collection.Facts[key]; exists {
			// Update existing fact
			fact.Value = value
			fact.Timestamp = time.Now()
		} else {
			// Create new fact
			c.createFact(collection, key, value)
		}
	}
}

// createFact creates a fact in the collection
func (c *Collector) createFact(collection *spookyfactstypes.FactCollection, key string, value interface{}) {
	collection.Facts[key] = &spookyfactstypes.Fact{
		Key:       key,
		Value:     value,
		Source:    string(c.GetSource()),
		Timestamp: time.Now(),
	}
}

// ImportCustomFacts imports custom facts from a source
func (c *Collector) ImportCustomFacts(source string, mergePolicy spookyfactstypes.MergePolicy) error {
	// Read source file
	file, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer file.Close()

	// Determine file type
	ext := strings.ToLower(filepath.Ext(source))

	var customFacts spookyfactstypes.CustomFacts
	switch ext {
	case ".json":
		if err := json.NewDecoder(file).Decode(&customFacts); err != nil {
			return fmt.Errorf("failed to decode JSON: %w", err)
		}
	case ".hcl":
		// For HCL, we would need to parse it properly
		return fmt.Errorf("HCL import not yet implemented")
	default:
		return fmt.Errorf("unsupported file format: %s", ext)
	}

	// Apply merge policy
	switch mergePolicy {
	case spookyfactstypes.MergePolicyReplace:
		c.overrides = make(map[string]interface{})
		for k, v := range customFacts.Custom {
			c.overrides[k] = v
		}
	case spookyfactstypes.MergePolicyMerge:
		for k, v := range customFacts.Custom {
			c.overrides[k] = v
		}
	case spookyfactstypes.MergePolicyAppend:
		// Append logic would depend on the data type
		for k, v := range customFacts.Custom {
			c.overrides[k] = v
		}
	case spookyfactstypes.MergePolicySkip:
		// Only add facts that don't already exist
		for k, v := range customFacts.Custom {
			if _, exists := c.overrides[k]; !exists {
				c.overrides[k] = v
			}
		}
	}

	return nil
}

// ExportCustomFacts exports custom facts to a file
func (c *Collector) ExportCustomFacts(w io.Writer, format string) error {
	customFacts := spookyfactstypes.CustomFacts{
		Custom:    make(map[string]interface{}),
		Overrides: c.overrides,
		Source:    c.customFactsPath,
	}

	switch format {
	case "json":
		return json.NewEncoder(w).Encode(customFacts)
	case "hcl":
		return c.exportToHCL(w, customFacts)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

// exportToHCL exports custom facts to HCL format
func (c *Collector) exportToHCL(w io.Writer, customFacts spookyfactstypes.CustomFacts) error {
	// Write HCL header
	if _, err := fmt.Fprintln(w, "# Custom Facts"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "# Generated by spooky custom collector"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, ""); err != nil {
		return err
	}

	// Write custom facts
	if _, err := fmt.Fprintln(w, "# Custom facts"); err != nil {
		return err
	}
	for key, value := range customFacts.Custom {
		if _, err := fmt.Fprintf(w, "%s = %q\n", key, value); err != nil {
			return err
		}
	}

	// Write overrides
	if len(customFacts.Overrides) > 0 {
		if _, err := fmt.Fprintln(w, ""); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, "# Overrides"); err != nil {
			return err
		}
		for key, value := range customFacts.Overrides {
			if _, err := fmt.Fprintf(w, "%s = %q\n", key, value); err != nil {
				return err
			}
		}
	}

	return nil
}
