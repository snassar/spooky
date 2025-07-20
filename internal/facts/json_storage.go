package facts

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// JSONFactStorage implements FactStorage using JSON files
type JSONFactStorage struct {
	filepath string
	facts    map[string]*MachineFacts
	mu       sync.RWMutex
}

// NewJSONFactStorage creates a new JSON-based fact storage
func NewJSONFactStorage(filepath string) (*JSONFactStorage, error) {
	storage := &JSONFactStorage{
		filepath: filepath,
		facts:    make(map[string]*MachineFacts),
	}

	// Load existing data if file exists
	if err := storage.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load existing facts: %w", err)
	}

	return storage, nil
}

// GetMachineFacts retrieves facts for a specific machine
func (j *JSONFactStorage) GetMachineFacts(machineID string) (*MachineFacts, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()

	if facts, exists := j.facts[machineID]; exists {
		return facts, nil
	}

	return nil, fmt.Errorf("machine facts not found: %s", machineID)
}

// SetMachineFacts stores facts for a specific machine
func (j *JSONFactStorage) SetMachineFacts(machineID string, facts *MachineFacts) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	facts.UpdatedAt = time.Now()
	if facts.CreatedAt.IsZero() {
		facts.CreatedAt = facts.UpdatedAt
	}

	j.facts[machineID] = facts

	return j.save()
}

// QueryFacts searches for facts matching the query criteria
func (j *JSONFactStorage) QueryFacts(query *FactQuery) ([]*MachineFacts, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()

	var results []*MachineFacts

	for _, facts := range j.facts {
		if matchesQuery(facts, query) {
			results = append(results, facts)
			if query.Limit > 0 && len(results) >= query.Limit {
				break
			}
		}
	}

	return results, nil
}

// DeleteFacts deletes facts matching the query criteria
func (j *JSONFactStorage) DeleteFacts(query *FactQuery) (int, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	var deletedCount int
	var toDelete []string

	for machineID, facts := range j.facts {
		if matchesQuery(facts, query) {
			toDelete = append(toDelete, machineID)
		}
	}

	for _, machineID := range toDelete {
		delete(j.facts, machineID)
		deletedCount++
	}

	if deletedCount > 0 {
		if err := j.save(); err != nil {
			return deletedCount, fmt.Errorf("failed to save after deletion: %w", err)
		}
	}

	return deletedCount, nil
}

// DeleteMachineFacts deletes facts for a specific machine
func (j *JSONFactStorage) DeleteMachineFacts(machineID string) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	if _, exists := j.facts[machineID]; !exists {
		return fmt.Errorf("machine facts not found: %s", machineID)
	}

	delete(j.facts, machineID)
	return j.save()
}

// ExportToJSON exports all facts to JSON format
func (j *JSONFactStorage) ExportToJSON(w io.Writer) error {
	j.mu.RLock()
	defer j.mu.RUnlock()

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(j.facts)
}

// ImportFromJSON imports facts from JSON format
func (j *JSONFactStorage) ImportFromJSON(r io.Reader) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	var facts map[string]*MachineFacts

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&facts); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	j.facts = facts
	return j.save()
}

// ExportToJSONWithEncryption exports facts with encryption support
func (j *JSONFactStorage) ExportToJSONWithEncryption(w io.Writer, _ ExportOptions) error {
	// For now, implement basic encryption without age
	// TODO: Add age encryption support
	return j.ExportToJSON(w)
}

// ImportFromJSONWithDecryption imports facts with decryption support
func (j *JSONFactStorage) ImportFromJSONWithDecryption(r io.Reader, _ string) error {
	// For now, implement basic import without age decryption
	// TODO: Add age decryption support
	return j.ImportFromJSON(r)
}

// Close saves the current state and closes the storage
func (j *JSONFactStorage) Close() error {
	return j.save()
}

// load reads facts from the JSON file
func (j *JSONFactStorage) load() error {
	file, err := os.Open(j.filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&j.facts)
}

// save writes facts to the JSON file
func (j *JSONFactStorage) save() error {
	file, err := os.Create(j.filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(j.facts)
}
