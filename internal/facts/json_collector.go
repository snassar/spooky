package facts

import (
	"fmt"
	"io"
	"os"

	"spooky/internal/logging"
)

// JSONCollector collects facts from local JSON files
type JSONCollector struct {
	filePath    string
	mergePolicy MergePolicy
	logger      logging.Logger
}

// NewJSONCollector creates a new JSON-based fact collector
func NewJSONCollector(filePath string, mergePolicy MergePolicy) *JSONCollector {
	return &JSONCollector{
		filePath:    filePath,
		mergePolicy: mergePolicy,
		logger:      logging.GetLogger(),
	}
}

// Validate validates the collector configuration
func (c *JSONCollector) Validate() error {
	if c.filePath == "" {
		return ErrInvalidSource("JSON file", "file path cannot be empty")
	}

	// Check if file exists and is readable
	if _, err := os.Stat(c.filePath); err != nil {
		return ErrInvalidSource(c.filePath, fmt.Sprintf("file not accessible: %v", err))
	}

	if err := ValidateMergePolicy(c.mergePolicy); err != nil {
		return ErrInvalidSource("merge policy", err.Error())
	}

	return nil
}

// Collect reads facts from a JSON file and converts them to FactCollection
func (c *JSONCollector) Collect(server string) (*FactCollection, error) {
	// Validate inputs
	if err := validateServer(server); err != nil {
		return nil, err
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}

	c.logger.Debug("Starting JSON fact collection",
		logging.Field{Key: "file", Value: c.filePath},
		logging.Field{Key: "server", Value: server},
		logging.Field{Key: "merge_policy", Value: c.mergePolicy})

	file, err := os.Open(c.filePath)
	if err != nil {
		c.logger.Error("Failed to open JSON file", err,
			logging.Field{Key: "file", Value: c.filePath})
		return nil, fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	collection, err := c.collectFromReader(file, server)
	if err != nil {
		c.logger.Error("Failed to collect facts from JSON file", err,
			logging.Field{Key: "file", Value: c.filePath},
			logging.Field{Key: "server", Value: server})
		return nil, err
	}

	c.logger.Info("Successfully collected facts from JSON file",
		logging.Field{Key: "file", Value: c.filePath},
		logging.Field{Key: "server", Value: server},
		logging.Field{Key: "fact_count", Value: len(collection.Facts)})

	return collection, nil
}

// CollectSpecific reads specific facts from a JSON file
func (c *JSONCollector) CollectSpecific(server string, keys []string) (*FactCollection, error) {
	return collectSpecificFacts(c, server, keys, c.logger, "JSON file")
}

// GetFact retrieves a single fact from the JSON file
func (c *JSONCollector) GetFact(server, key string) (*Fact, error) {
	return getSpecificFact(c, server, key, c.logger, "JSON file")
}

// collectFromReader reads facts from an io.Reader and converts them
func (c *JSONCollector) collectFromReader(r io.Reader, server string) (*FactCollection, error) {
	sourceInfo := buildStandardMetadata("json", c.filePath, "json")
	return parseJSONFromReader(r, server, sourceInfo)
}
