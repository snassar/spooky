package hcl

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"spooky/internal/facts/collectors"
	"spooky/internal/facts/types"

	hcl2 "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// Collector collects facts from HCL configuration files
type Collector struct {
	collectors.BaseCollector
	hclFiles []string
	parser   HCLParser
}

// HCLParser defines the interface for parsing HCL content
type HCLParser interface {
	ParseFile(filePath string) (map[string]interface{}, error)
	ParseContent(content []byte) (map[string]interface{}, error)
}

// DefaultParser is a simple HCL parser implementation
type DefaultParser struct{}

// NewCollector creates a new HCL fact collector
func NewCollector(hclFiles []string) *Collector {
	return &Collector{
		BaseCollector: *collectors.NewBaseCollector(types.SourceHCL, types.MergePolicyMerge),
		hclFiles:      hclFiles,
		parser:        &DefaultParser{},
	}
}

// NewCollectorWithParser creates a new HCL fact collector with a custom parser
func NewCollectorWithParser(hclFiles []string, parser HCLParser) *Collector {
	return &Collector{
		BaseCollector: *collectors.NewBaseCollector(types.SourceHCL, types.MergePolicyMerge),
		hclFiles:      hclFiles,
		parser:        parser,
	}
}

// Collect gathers facts from all configured HCL files
func (c *Collector) Collect(server string) (*types.FactCollection, error) {
	collection := &types.FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*types.Fact),
	}

	// Collect facts from each HCL file
	for _, hclFile := range c.hclFiles {
		if err := c.collectFromFile(collection, hclFile); err != nil {
			return nil, fmt.Errorf("failed to collect from %s: %w", hclFile, err)
		}
	}

	return collection, nil
}

// CollectSpecific collects only the specified facts from HCL files
func (c *Collector) CollectSpecific(server string, keys []string) (*types.FactCollection, error) {
	collection := &types.FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*types.Fact),
	}

	// Collect all facts first, then filter
	allFacts, err := c.Collect(server)
	if err != nil {
		return nil, err
	}

	// Copy facts to our collection
	for key, fact := range allFacts.Facts {
		collection.Facts[key] = fact
	}

	// Filter to only requested keys
	filteredCollection := &types.FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*types.Fact),
	}

	for _, key := range keys {
		if fact, exists := collection.Facts[key]; exists {
			filteredCollection.Facts[key] = fact
		}
	}

	return filteredCollection, nil
}

// GetFact retrieves a single fact from HCL files
func (c *Collector) GetFact(server, key string) (*types.Fact, error) {
	collection := &types.FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*types.Fact),
	}

	// Search through all HCL files for the specific fact
	for _, hclFile := range c.hclFiles {
		if err := c.collectFromFile(collection, hclFile); err != nil {
			continue // Skip files that can't be parsed
		}

		if fact, exists := collection.Facts[key]; exists {
			return fact, nil
		}
	}

	return nil, fmt.Errorf("fact %s not found in any HCL file", key)
}

// Validate validates the collector configuration
func (c *Collector) Validate() error {
	if len(c.hclFiles) == 0 {
		return fmt.Errorf("no HCL files configured")
	}

	// Check if all HCL files exist
	for _, hclFile := range c.hclFiles {
		if _, err := os.Stat(hclFile); os.IsNotExist(err) {
			return fmt.Errorf("HCL file does not exist: %s", hclFile)
		}
	}

	return nil
}

// AddHCLFile adds an HCL file to the collection list
func (c *Collector) AddHCLFile(filePath string) {
	c.hclFiles = append(c.hclFiles, filePath)
}

// RemoveHCLFile removes an HCL file from the collection list
func (c *Collector) RemoveHCLFile(filePath string) {
	for i, file := range c.hclFiles {
		if file == filePath {
			c.hclFiles = append(c.hclFiles[:i], c.hclFiles[i+1:]...)
			break
		}
	}
}

// GetHCLFiles returns the list of configured HCL files
func (c *Collector) GetHCLFiles() []string {
	result := make([]string, len(c.hclFiles))
	copy(result, c.hclFiles)
	return result
}

// collectFromFile collects facts from a single HCL file
func (c *Collector) collectFromFile(collection *types.FactCollection, filePath string) error {
	// Parse the HCL file
	facts, err := c.parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse HCL file: %w", err)
	}

	// Add facts to collection
	for key, value := range facts {
		c.createFact(collection, key, value)
	}

	return nil
}

// createFact creates a fact in the collection
func (c *Collector) createFact(collection *types.FactCollection, key string, value interface{}) {
	collection.Facts[key] = &types.Fact{
		Key:       key,
		Value:     value,
		Source:    string(c.GetSource()),
		Timestamp: time.Now(),
	}
}

// ParseHCLFile parses an HCL file and returns the facts
func (c *Collector) ParseHCLFile(filePath string) (map[string]interface{}, error) {
	return c.parser.ParseFile(filePath)
}

// ParseHCLContent parses HCL content and returns the facts
func (c *Collector) ParseHCLContent(content []byte) (map[string]interface{}, error) {
	return c.parser.ParseContent(content)
}

// ExportFactsToHCL exports facts to HCL format
func (c *Collector) ExportFactsToHCL(facts map[string]interface{}, w io.Writer) error {
	// Write HCL header
	if _, err := fmt.Fprintln(w, "# Facts exported by spooky HCL collector"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "# Generated at:", time.Now().Format(time.RFC3339)); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, ""); err != nil {
		return err
	}

	// Write facts
	for key, value := range facts {
		if err := c.writeHCLValue(w, key, value); err != nil {
			return err
		}
	}

	return nil
}

// writeHCLValue writes a single value to HCL format
func (c *Collector) writeHCLValue(w io.Writer, key string, value interface{}) error {
	switch v := value.(type) {
	case string:
		if _, err := fmt.Fprintf(w, "%s = %q\n", key, v); err != nil {
			return err
		}
	case int, int32, int64, uint, uint32, uint64:
		if _, err := fmt.Fprintf(w, "%s = %v\n", key, v); err != nil {
			return err
		}
	case float32, float64:
		if _, err := fmt.Fprintf(w, "%s = %v\n", key, v); err != nil {
			return err
		}
	case bool:
		if _, err := fmt.Fprintf(w, "%s = %t\n", key, v); err != nil {
			return err
		}
	case []interface{}:
		if _, err := fmt.Fprintf(w, "%s = [\n", key); err != nil {
			return err
		}
		for _, item := range v {
			if _, err := fmt.Fprintf(w, "  %q,\n", item); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w, "]"); err != nil {
			return err
		}
	case map[string]interface{}:
		if _, err := fmt.Fprintf(w, "%s = {\n", key); err != nil {
			return err
		}
		for k, v := range v {
			if _, err := fmt.Fprintf(w, "  %s = %q\n", k, v); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w, "}"); err != nil {
			return err
		}
	default:
		if _, err := fmt.Fprintf(w, "%s = %q\n", key, fmt.Sprintf("%v", v)); err != nil {
			return err
		}
	}

	return nil
}

// DefaultParser implementation

// ParseFile parses an HCL file and returns the facts
func (p *DefaultParser) ParseFile(filePath string) (map[string]interface{}, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return p.ParseContent(content)
}

// ParseContent parses HCL content and returns the facts using proper HCL library
func (p *DefaultParser) ParseContent(content []byte) (map[string]interface{}, error) {
	// Use proper HCL parser
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(content, "facts.hcl")
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
		value, err := p.extractHCLValue(attr.Expr)
		if err != nil {
			return nil, fmt.Errorf("failed to extract value for %s: %w", name, err)
		}
		facts[name] = value
	}

	return facts, nil
}

// extractHCLValue extracts a value from an HCL expression
func (p *DefaultParser) extractHCLValue(expr hcl2.Expression) (interface{}, error) {
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

// FindHCLFiles finds all HCL files in a directory
func (c *Collector) FindHCLFiles(directory string) ([]string, error) {
	var hclFiles []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".hcl") {
			hclFiles = append(hclFiles, path)
		}

		return nil
	})

	return hclFiles, err
}

// LoadHCLFilesFromDirectory loads all HCL files from a directory
func (c *Collector) LoadHCLFilesFromDirectory(directory string) error {
	hclFiles, err := c.FindHCLFiles(directory)
	if err != nil {
		return fmt.Errorf("failed to find HCL files: %w", err)
	}

	c.hclFiles = append(c.hclFiles, hclFiles...)
	return nil
}

// ValidateHCLFile validates that a file is valid HCL
func (c *Collector) ValidateHCLFile(filePath string) error {
	_, err := c.parser.ParseFile(filePath)
	return err
}

// GetFactSources returns information about where facts were collected from
func (c *Collector) GetFactSources() []string {
	return c.hclFiles
}
