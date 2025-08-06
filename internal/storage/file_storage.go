package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	spookyfactstypes "spooky/internal/facts/types"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// FileStorageFormat defines the format for file-based storage
type FileStorageFormat string

const (
	FileStorageFormatJSON FileStorageFormat = "json"
	FileStorageFormatHCL  FileStorageFormat = "hcl"
)

// BaseFileStorage provides common functionality for file-based fact storage
type BaseFileStorage struct {
	filepath string
	format   FileStorageFormat
	facts    map[string]*spookyfactstypes.FactCollection
	mu       sync.RWMutex
}

// NewBaseFileStorage creates a new base file storage
func NewBaseFileStorage(filepath string, format FileStorageFormat) (*BaseFileStorage, error) {
	storage := &BaseFileStorage{
		filepath: filepath,
		format:   format,
		facts:    make(map[string]*spookyfactstypes.FactCollection),
	}

	// Load existing data if file exists
	if err := storage.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load existing facts: %w", err)
	}

	return storage, nil
}

// SetFactCollection stores a fact collection for a specific machine
func (b *BaseFileStorage) SetFactCollection(machineID string, collection *spookyfactstypes.FactCollection) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Update timestamp
	collection.Timestamp = time.Now()
	b.facts[machineID] = collection

	return b.save()
}

// GetFactCollection retrieves a fact collection for a specific machine
func (b *BaseFileStorage) GetFactCollection(machineID string) (*spookyfactstypes.FactCollection, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if collection, exists := b.facts[machineID]; exists {
		return collection, nil
	}

	return nil, fmt.Errorf("fact collection not found: %s", machineID)
}

// QueryFactCollections searches for fact collections matching the query criteria
func (b *BaseFileStorage) QueryFactCollections(query *spookyfactstypes.FactQuery) ([]*spookyfactstypes.FactCollection, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var results []*spookyfactstypes.FactCollection

	for _, collection := range b.facts {
		if matchesQuery(collection, query) {
			results = append(results, collection)
			if query.Limit > 0 && len(results) >= query.Limit {
				break
			}
		}
	}

	return results, nil
}

// DeleteFactCollection deletes a fact collection for a specific machine
func (b *BaseFileStorage) DeleteFactCollection(machineID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.facts, machineID)
	return b.save()
}

// DeleteFactCollections deletes fact collections matching the query criteria
func (b *BaseFileStorage) DeleteFactCollections(query *spookyfactstypes.FactQuery) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	var deletedCount int
	for machineID, collection := range b.facts {
		if matchesQuery(collection, query) {
			delete(b.facts, machineID)
			deletedCount++
		}
	}

	if deletedCount > 0 {
		return deletedCount, b.save()
	}

	return deletedCount, nil
}

// ExportToJSON exports all fact collections to JSON
func (b *BaseFileStorage) ExportToJSON(w io.Writer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(b.facts)
}

// ImportFromJSON imports fact collections from JSON
func (b *BaseFileStorage) ImportFromJSON(r io.Reader) error {
	var facts map[string]*spookyfactstypes.FactCollection
	if err := json.NewDecoder(r).Decode(&facts); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.facts = facts
	return b.save()
}

// ExportToJSONWithEncryption exports fact collections with encryption
func (b *BaseFileStorage) ExportToJSONWithEncryption(w io.Writer, _ ExportOptions) error {
	// TODO: Implement encryption support
	return b.ExportToJSON(w)
}

// ImportFromJSONWithDecryption imports fact collections with decryption
func (b *BaseFileStorage) ImportFromJSONWithDecryption(r io.Reader, _ string) error {
	// TODO: Implement decryption support
	return b.ImportFromJSON(r)
}

// Close closes the file storage (no-op for file-based storage)
func (b *BaseFileStorage) Close() error {
	return nil
}

// load loads facts from the file
func (b *BaseFileStorage) load() error {
	switch b.format {
	case FileStorageFormatJSON:
		return b.loadJSON()
	case FileStorageFormatHCL:
		return b.loadHCL()
	default:
		return fmt.Errorf("unsupported storage format: %s", b.format)
	}
}

// save saves facts to the file
func (b *BaseFileStorage) save() error {
	switch b.format {
	case FileStorageFormatJSON:
		return b.saveJSON()
	case FileStorageFormatHCL:
		return b.saveHCL()
	default:
		return fmt.Errorf("unsupported storage format: %s", b.format)
	}
}

// loadJSON loads facts from JSON file
func (b *BaseFileStorage) loadJSON() error {
	file, err := os.Open(b.filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&b.facts)
}

// saveJSON saves facts to JSON file
func (b *BaseFileStorage) saveJSON() error {
	file, err := os.Create(b.filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(b.facts)
}

// HCLFactsWrapper wraps facts for HCL serialization
type HCLFactsWrapper struct {
	Facts map[string]*spookyfactstypes.FactCollection `hcl:"facts,block"`
}

// loadHCL loads facts from HCL file
func (b *BaseFileStorage) loadHCL() error {
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(b.filepath)
	if diags.HasErrors() {
		// If HCL parsing fails, try JSON as fallback for backward compatibility
		return b.loadJSON()
	}

	var wrapper HCLFactsWrapper
	diags = gohcl.DecodeBody(file.Body, nil, &wrapper)
	if diags.HasErrors() {
		// If HCL decoding fails, try JSON as fallback for backward compatibility
		return b.loadJSON()
	}

	b.facts = wrapper.Facts
	return nil
}

// saveHCL saves facts to HCL file
func (b *BaseFileStorage) saveHCL() error {
	// Generate HCL content as a string, following the existing codebase patterns
	content := b.generateHCLContent()

	// Write to file
	return os.WriteFile(b.filepath, []byte(content), 0o600)
}

// generateHCLContent generates HCL content for facts storage
func (b *BaseFileStorage) generateHCLContent() string {
	var lines []string

	// Add header comment
	lines = append(lines, "# Spooky Facts Storage")
	lines = append(lines, "# This file contains collected facts from machines")
	lines = append(lines, "")

	// Add facts wrapper block
	lines = append(lines, "facts {")

	// Add each machine's facts
	for machineID, collection := range b.facts {
		lines = append(lines, fmt.Sprintf("  machine \"%s\" {", machineID))
		lines = append(lines, fmt.Sprintf("    server = \"%s\"", collection.Server))
		lines = append(lines, fmt.Sprintf("    timestamp = \"%s\"", collection.Timestamp.Format(time.RFC3339)))

		// Add facts
		if len(collection.Facts) > 0 {
			lines = append(lines, "    facts = {")
			for key, fact := range collection.Facts {
				lines = append(lines, fmt.Sprintf("      %s = %s", key, b.formatHCLValue(fact.Value)))
			}
			lines = append(lines, "    }")
		}

		// Add custom facts if present
		if len(collection.CustomFacts) > 0 {
			lines = append(lines, "    custom_facts = {")
			for filename, facts := range collection.CustomFacts {
				lines = append(lines, fmt.Sprintf("      %s = {", filename))
				for key, value := range facts {
					lines = append(lines, fmt.Sprintf("        %s = %s", key, b.formatHCLValue(value)))
				}
				lines = append(lines, "      }")
			}
			lines = append(lines, "    }")
		}

		lines = append(lines, "  }")
		lines = append(lines, "")
	}

	lines = append(lines, "}")

	return strings.Join(lines, "\n")
}

// formatHCLValue formats a value for HCL output
func (b *BaseFileStorage) formatHCLValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", v)
	case int, int32, int64:
		return fmt.Sprintf("%v", v)
	case float32, float64:
		return fmt.Sprintf("%v", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case nil:
		return "null"
	default:
		// For complex types, convert to JSON string
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("\"%v\"", v)
		}
		return fmt.Sprintf("\"%s\"", string(jsonBytes))
	}
}
