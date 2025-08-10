package facts

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookytypesfacts "spooky/internal/types/facts"
)

// parseJSONFromReader handles JSON parsing with shared logic
func parseJSONFromReader(r io.Reader, server string, sourceInfo map[string]interface{}) (*spookytypesfacts.FactCollection, error) {
	// Try to parse as different JSON formats
	var data interface{}
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	collection := &spookytypesfacts.FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*spookytypesfacts.Fact),
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
func parseArrayJSON(data []interface{}, collection *spookytypesfacts.FactCollection, sourceInfo map[string]interface{}) (*spookytypesfacts.FactCollection, error) {
	for _, item := range data {
		if factObj, ok := item.(map[string]interface{}); ok {
			key, keyOk := factObj["key"].(string)
			value, valueOk := factObj["value"]

			if keyOk && valueOk {
				fact := &spookytypesfacts.Fact{
					Key:       key,
					Value:     value,
					Source:    string(spookytypesfacts.SourceCustom),
					Server:    collection.Server,
					Timestamp: collection.Timestamp,
					TTL:       spookytypesfacts.DefaultTTL,
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
func parseFlatJSON(data map[string]interface{}, collection *spookytypesfacts.FactCollection, sourceInfo map[string]interface{}) (*spookytypesfacts.FactCollection, error) {
	for key, value := range data {
		fact := &spookytypesfacts.Fact{
			Key:       key,
			Value:     value,
			Source:    string(spookytypesfacts.SourceCustom),
			Server:    collection.Server,
			Timestamp: collection.Timestamp,
			TTL:       spookytypesfacts.DefaultTTL,
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
func filterFactCollection(collection *spookytypesfacts.FactCollection, keys []string, logger spookyinterfaces.Logger, source string) *spookytypesfacts.FactCollection {
	filtered := &spookytypesfacts.FactCollection{
		Server:    collection.Server,
		Timestamp: collection.Timestamp,
		Facts:     make(map[string]*spookytypesfacts.Fact),
	}

	for _, key := range keys {
		if fact, exists := collection.Facts[key]; exists {
			filtered.Facts[key] = fact
		} else {
			logger.Warn("Requested fact not found",
				spookylogging.String("key", key),
				spookylogging.String("source", source))
		}
	}

	logger.Debug("Completed fact filtering",
		spookylogging.String("source", source),
		spookylogging.String("server", collection.Server),
		spookylogging.Int("requested_count", len(keys)),
		spookylogging.Int("found_count", len(filtered.Facts)))

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
func collectSpecificFacts(collector spookytypesfacts.FactCollector, server string, keys []string, logger spookyinterfaces.Logger, source string) (*spookytypesfacts.FactCollection, error) {
	logger.Debug("Starting specific fact collection",
		spookylogging.String("server", server),
		spookylogging.String("requested_keys", fmt.Sprintf("%v", keys)),
		spookylogging.String("source", source))

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
func getSpecificFact(collector spookytypesfacts.FactCollector, server, key string, logger spookyinterfaces.Logger, source string) (*spookytypesfacts.Fact, error) {
	logger.Debug("Getting specific fact",
		spookylogging.String("server", server),
		spookylogging.String("key", key),
		spookylogging.String("source", source))

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
			spookylogging.String("server", server),
			spookylogging.String("key", key),
			spookylogging.String("source", source))
		return fact, nil
	}

	logger.Warn("Fact not found",
		spookylogging.String("server", server),
		spookylogging.String("key", key),
		spookylogging.String("source", source))

	return nil, ErrFactNotFoundInSource(key, server, source)
}

// buildStandardMetadata creates standardized metadata for a collector
func buildStandardMetadata(collectorType, source, format string) map[string]interface{} {
	return map[string]interface{}{
		"source":       collectorType,
		"source_path":  source,
		"format":       format,
		"collected_at": time.Now().UTC().Format(time.RFC3339),
	}
}

// collectFromFile is a common implementation for file-based collectors
func collectFromFile(filePath, server, sourceType string, logger spookyinterfaces.Logger,
	parseFunc func() (interface{}, error),
	extractFunc func(interface{}, string) map[string]*spookytypesfacts.Fact) (*spookytypesfacts.FactCollection, error) {

	logger.Debug("Starting file-based fact collection",
		spookylogging.String("file", filePath),
		spookylogging.String("server", server),
		spookylogging.String("source_type", sourceType))

	// Validate inputs
	if err := validateServer(server); err != nil {
		return nil, err
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logger.Error("File not found",
			spookylogging.String("file", filePath),
			spookylogging.Error(err))
		return nil, fmt.Errorf("invalid source %s: file not found", filePath)
	}

	// Parse file
	data, err := parseFunc()
	if err != nil {
		logger.Error("Failed to parse file",
			spookylogging.String("file", filePath),
			spookylogging.Error(err))
		return nil, fmt.Errorf("failed to parse %s file: %w", sourceType, err)
	}

	// Extract facts from data
	facts := extractFunc(data, server)

	collection := &spookytypesfacts.FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     facts,
	}

	logger.Info("Successfully collected facts from file",
		spookylogging.String("file", filePath),
		spookylogging.String("server", server),
		spookylogging.Int("fact_count", len(facts)),
		spookylogging.String("source_type", sourceType))

	return collection, nil
}
