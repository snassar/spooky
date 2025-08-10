package json

import (
	spookytypes "spooky/internal/types"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	spookyfactscollectors "spooky/internal/facts/collectors"
	spookylogging "spooky/internal/logging"
	spookytypes "spooky/internal/types"
)

// Collector collects facts from JSON files
type Collector struct {
	spookyfactscollectors.BaseCollector
	sourcePath string
	logger     spookylogging.Logger
}

// NewCollector creates a new JSON fact collector
func NewCollector(sourcePath string, logger spookylogging.Logger) *Collector {
	return &Collector{
		BaseCollector: *spookyfactscollectors.NewBaseCollector(spookytypes.SourceJSON, spookytypes.MergePolicyReplace),
		sourcePath:    sourcePath,
		logger:        logger,
	}
}

// Collect gathers facts from JSON source
func (c *Collector) Collect(server string) (*spookytypes.FactCollection, error) {
	collection := &spookytypes.FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*spookytypes.Fact),
	}

	// Determine if source is file or directory
	info, err := os.Stat(c.sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat source path: %w", err)
	}

	if info.IsDir() {
		return c.collectFromDirectory(collection)
	}
	return c.collectFromFile(collection, c.sourcePath)
}

// CollectSpecific collects only the specified facts
func (c *Collector) CollectSpecific(server string, keys []string) (*spookytypes.FactCollection, error) {
	// Collect all facts first
	collection, err := c.Collect(server)
	if err != nil {
		return nil, err
	}

	// Filter to only requested keys
	filteredCollection := &spookytypes.FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*spookytypes.Fact),
	}

	for _, key := range keys {
		if fact, exists := collection.Facts[key]; exists {
			filteredCollection.Facts[key] = fact
		}
	}

	return filteredCollection, nil
}

// GetFact retrieves a single fact by key
func (c *Collector) GetFact(server, key string) (*spookytypes.Fact, error) {
	collection, err := c.Collect(server)
	if err != nil {
		return nil, err
	}

	if fact, exists := collection.Facts[key]; exists {
		return fact, nil
	}

	return nil, fmt.Errorf("fact '%s' not found", key)
}

// Validate validates the collector configuration
func (c *Collector) Validate() error {
	if c.sourcePath == "" {
		return fmt.Errorf("source path cannot be empty")
	}

	// Check if source exists
	if _, err := os.Stat(c.sourcePath); err != nil {
		return fmt.Errorf("source path does not exist: %w", err)
	}

	return nil
}

// collectFromFile processes a single JSON file
func (c *Collector) collectFromFile(collection *spookytypes.FactCollection, filePath string) (*spookytypes.FactCollection, error) {
	// Validate file before processing
	if err := c.validateJSONFile(filePath); err != nil {
		return nil, err
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Parse JSON
	var jsonData interface{}
	if err := json.Unmarshal(content, &jsonData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON in %s: %w", filePath, err)
	}

	// Convert JSON to facts
	facts := c.convertJSONToFacts(jsonData, filepath.Base(filePath))

	// Add facts to collection
	for key, fact := range facts {
		collection.Facts[key] = fact
	}

	return collection, nil
}

// collectFromDirectory processes all JSON files in a directory
func (c *Collector) collectFromDirectory(collection *spookytypes.FactCollection) (*spookytypes.FactCollection, error) {
	err := filepath.WalkDir(c.sourcePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Only process JSON files
		if !strings.HasSuffix(strings.ToLower(path), ".json") {
			return nil
		}

		// Process file
		_, err = c.collectFromFile(collection, path)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to process directory: %w", err)
	}

	return collection, nil
}

// convertJSONToFacts converts JSON data to facts
func (c *Collector) convertJSONToFacts(data interface{}, source string) map[string]*spookytypes.Fact {
	facts := make(map[string]*spookytypes.Fact)
	c.convertValueToFacts("", data, facts, source)
	return facts
}

// convertValueToFacts recursively converts JSON values to facts
func (c *Collector) convertValueToFacts(prefix string, value interface{}, facts map[string]*spookytypes.Fact, source string) {
	switch v := value.(type) {
	case map[string]interface{}:
		// Handle objects
		for key, val := range v {
			newPrefix := key
			if prefix != "" {
				newPrefix = prefix + "." + key
			}
			c.convertValueToFacts(newPrefix, val, facts, source)
		}
	case []interface{}:
		// Handle arrays
		if prefix == "" {
			// Root level array - store as single fact
			facts["array"] = &spookytypes.Fact{
				Key:       "array",
				Value:     v,
				Source:    string(c.GetSource()),
				Server:    "local", // JSON facts are always local
				Timestamp: time.Now(),
				TTL:       spookytypes.DefaultTTL,
				Metadata: map[string]interface{}{
					"source_machine": "local",
					"source_file":    source,
					"source_type":    "json",
					"json_type":      fmt.Sprintf("%T", v),
					"array_length":   len(v),
				},
			}
		} else {
			// Nested array - process individual elements
			for i, val := range v {
				newPrefix := fmt.Sprintf("%s[%d]", prefix, i)
				c.convertValueToFacts(newPrefix, val, facts, source)
			}
		}
	default:
		// Handle primitive values
		if prefix != "" {
			facts[prefix] = &spookytypes.Fact{
				Key:       prefix,
				Value:     v,
				Source:    string(c.GetSource()),
				Server:    "local", // JSON facts are always local
				Timestamp: time.Now(),
				TTL:       spookytypes.DefaultTTL,
				Metadata: map[string]interface{}{
					"source_machine": "local",
					"source_file":    source,
					"source_type":    "json",
					"json_type":      fmt.Sprintf("%T", v),
				},
			}
		}
	}
}

// validateJSONFile validates JSON file before processing
func (c *Collector) validateJSONFile(filePath string) error {
	// Check file extension
	if !strings.HasSuffix(strings.ToLower(filePath), ".json") {
		return fmt.Errorf("file %s is not a JSON file", filePath)
	}

	// Check file size (prevent processing huge files)
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	const maxFileSize = 10 * 1024 * 1024 // 10MB
	if info.Size() > maxFileSize {
		return fmt.Errorf("file %s is too large (%d bytes)", filePath, info.Size())
	}

	return nil
}
