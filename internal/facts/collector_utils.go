package facts

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"spooky/internal/logging"
)

// parseJSONFromReader handles JSON parsing with shared logic
func parseJSONFromReader(r io.Reader, server string, sourceInfo map[string]interface{}) (*FactCollection, error) {
	// Try to parse as different JSON formats
	var data interface{}
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	collection := &FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*Fact),
	}

	// Handle different JSON formats
	switch v := data.(type) {
	case map[string]interface{}:
		// Format: {"fact_key": "value", "another_fact": 123}
		return parseFlatJSON(v, collection, sourceInfo)
	case []interface{}:
		// Format: [{"key": "fact1", "value": "value1"}, ...]
		return parseArrayJSON(v, collection, sourceInfo)
	default:
		return nil, fmt.Errorf("unsupported JSON format")
	}
}

// parseArrayJSON handles array-based JSON format with shared logic
func parseArrayJSON(data []interface{}, collection *FactCollection, sourceInfo map[string]interface{}) (*FactCollection, error) {
	for _, item := range data {
		if factObj, ok := item.(map[string]interface{}); ok {
			key, keyOk := factObj["key"].(string)
			value, valueOk := factObj["value"]

			if keyOk && valueOk {
				fact := &Fact{
					Key:       key,
					Value:     value,
					Source:    string(SourceCustom),
					Server:    collection.Server,
					Timestamp: collection.Timestamp,
					TTL:       DefaultTTL,
					Metadata:  make(map[string]interface{}),
				}

				// Copy source info to metadata
				for k, v := range sourceInfo {
					fact.Metadata[k] = v
				}

				// Copy additional metadata if present
				if source, ok := factObj["source"].(string); ok {
					fact.Source = source
				}
				if ttl, ok := factObj["ttl"].(float64); ok {
					fact.TTL = time.Duration(ttl) * time.Second
				}
				if metadata, ok := factObj["metadata"].(map[string]interface{}); ok {
					for k, v := range metadata {
						fact.Metadata[k] = v
					}
				}

				collection.Facts[key] = fact
			}
		}
	}
	return collection, nil
}

// parseFlatJSON handles flat key-value JSON format with shared logic
func parseFlatJSON(data map[string]interface{}, collection *FactCollection, sourceInfo map[string]interface{}) (*FactCollection, error) {
	for key, value := range data {
		fact := &Fact{
			Key:       key,
			Value:     value,
			Source:    string(SourceCustom),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       DefaultTTL,
			Metadata:  make(map[string]interface{}),
		}

		// Copy source info to metadata
		for k, v := range sourceInfo {
			fact.Metadata[k] = v
		}

		collection.Facts[key] = fact
	}
	return collection, nil
}

// filterFactCollection filters a fact collection to only include specified keys
func filterFactCollection(collection *FactCollection, keys []string, logger logging.Logger, source string) *FactCollection {
	filtered := &FactCollection{
		Server:    collection.Server,
		Timestamp: collection.Timestamp,
		Facts:     make(map[string]*Fact),
	}

	for _, key := range keys {
		if fact, exists := collection.Facts[key]; exists {
			filtered.Facts[key] = fact
		} else {
			logger.Warn("Requested fact not found",
				logging.Field{Key: "key", Value: key},
				logging.Field{Key: "source", Value: source})
		}
	}

	logger.Debug("Completed fact filtering",
		logging.Field{Key: "source", Value: source},
		logging.Field{Key: "server", Value: collection.Server},
		logging.Field{Key: "requested_count", Value: len(keys)},
		logging.Field{Key: "found_count", Value: len(filtered.Facts)})

	return filtered
}

// validateFactKey validates that a fact key is not empty
func validateFactKey(key string) error {
	if key == "" {
		return fmt.Errorf("fact key cannot be empty")
	}
	return nil
}

// validateServer validates that a server name is not empty
func validateServer(server string) error {
	if server == "" {
		return fmt.Errorf("server name cannot be empty")
	}
	return nil
}

// validateKeys validates that the keys slice is not empty
func validateKeys(keys []string) error {
	if len(keys) == 0 {
		return fmt.Errorf("keys slice cannot be empty")
	}
	for i, key := range keys {
		if err := validateFactKey(key); err != nil {
			return fmt.Errorf("invalid key at index %d: %w", i, err)
		}
	}
	return nil
}

// collectSpecificFacts is a common implementation for CollectSpecific
func collectSpecificFacts(collector FactCollector, server string, keys []string, logger logging.Logger, source string) (*FactCollection, error) {
	logger.Debug("Starting specific fact collection",
		logging.Field{Key: "server", Value: server},
		logging.Field{Key: "requested_keys", Value: keys},
		logging.Field{Key: "source", Value: source})

	// Validate inputs
	if err := validateServer(server); err != nil {
		return nil, err
	}
	if err := validateKeys(keys); err != nil {
		return nil, err
	}

	collection, err := collector.Collect(server)
	if err != nil {
		return nil, err
	}

	return filterFactCollection(collection, keys, logger, source), nil
}

// getSpecificFact is a common implementation for GetFact
func getSpecificFact(collector FactCollector, server, key string, logger logging.Logger, source string) (*Fact, error) {
	logger.Debug("Getting specific fact",
		logging.Field{Key: "server", Value: server},
		logging.Field{Key: "key", Value: key},
		logging.Field{Key: "source", Value: source})

	// Validate inputs
	if err := validateServer(server); err != nil {
		return nil, err
	}
	if err := validateFactKey(key); err != nil {
		return nil, err
	}

	collection, err := collector.Collect(server)
	if err != nil {
		return nil, err
	}

	if fact, exists := collection.Facts[key]; exists {
		logger.Debug("Successfully retrieved fact",
			logging.Field{Key: "server", Value: server},
			logging.Field{Key: "key", Value: key},
			logging.Field{Key: "source", Value: source})
		return fact, nil
	}

	logger.Warn("Fact not found",
		logging.Field{Key: "server", Value: server},
		logging.Field{Key: "key", Value: key},
		logging.Field{Key: "source", Value: source})

	return nil, ErrFactNotFoundInSource(key, server, source)
}

// buildStandardMetadata creates standardized metadata for a collector
func buildStandardMetadata(collectorType, source, format string) map[string]interface{} {
	return NewMetadataBuilder().
		WithCollectorType(collectorType).
		WithSourcePath(source).
		WithFormat(format).
		Build()
}

// collectFromFile is a common implementation for file-based collectors
func collectFromFile(filePath, server, sourceType string, logger logging.Logger,
	parseFunc func() (interface{}, error),
	extractFunc func(interface{}, string) map[string]*Fact) (*FactCollection, error) {

	logger.Debug("Starting file-based fact collection",
		logging.Field{Key: "file", Value: filePath},
		logging.Field{Key: "server", Value: server},
		logging.Field{Key: "source_type", Value: sourceType})

	// Validate inputs
	if err := validateServer(server); err != nil {
		return nil, err
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logger.Error("File not found", err,
			logging.Field{Key: "file", Value: filePath})
		return nil, ErrInvalidSource(filePath, "file not found")
	}

	// Parse file
	data, err := parseFunc()
	if err != nil {
		logger.Error("Failed to parse file", err,
			logging.Field{Key: "file", Value: filePath})
		return nil, fmt.Errorf("failed to parse %s file: %w", sourceType, err)
	}

	// Extract facts from data
	facts := extractFunc(data, server)

	collection := &FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     facts,
	}

	logger.Info("Successfully collected facts from file",
		logging.Field{Key: "file", Value: filePath},
		logging.Field{Key: "server", Value: server},
		logging.Field{Key: "fact_count", Value: len(facts)},
		logging.Field{Key: "source_type", Value: sourceType})

	return collection, nil
}
