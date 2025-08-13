package facts

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"
)

// MemoryFactStorage provides in-memory storage for facts
type MemoryFactStorage struct {
	facts map[string]*FactCollection
	mutex sync.RWMutex
}

// NewMemoryFactStorage creates a new in-memory fact storage
func NewMemoryFactStorage() *MemoryFactStorage {
	return &MemoryFactStorage{
		facts: make(map[string]*FactCollection),
	}
}

// Store stores facts for a machine
func (s *MemoryFactStorage) Store(_ context.Context, machineID string, facts *FactCollection) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if facts == nil {
		return fmt.Errorf("facts cannot be nil")
	}

	s.facts[machineID] = facts
	return nil
}

// Get retrieves facts for a machine
func (s *MemoryFactStorage) Get(_ context.Context, machineID string) (*FactCollection, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	facts, exists := s.facts[machineID]
	if !exists {
		return nil, fmt.Errorf("facts not found for machine %s", machineID)
	}

	return facts, nil
}

// List lists all machine IDs with stored facts
func (s *MemoryFactStorage) List(_ context.Context) ([]string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	machineIDs := make([]string, 0, len(s.facts))
	for machineID := range s.facts {
		machineIDs = append(machineIDs, machineID)
	}

	return machineIDs, nil
}

// Delete deletes facts for a machine
func (s *MemoryFactStorage) Delete(_ context.Context, machineID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.facts, machineID)
	return nil
}

// Close closes the storage
func (s *MemoryFactStorage) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Clear all facts from memory
	s.facts = make(map[string]*FactCollection)
	return nil
}

// ExportToJSON exports facts to JSON format
func (s *MemoryFactStorage) ExportToJSON(w io.Writer) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	export := struct {
		ExportedAt time.Time                  `json:"exported_at"`
		Facts      map[string]*FactCollection `json:"facts"`
	}{
		ExportedAt: time.Now(),
		Facts:      s.facts,
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(export)
}

// ImportFromJSON imports facts from JSON format
func (s *MemoryFactStorage) ImportFromJSON(r io.Reader) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var importData struct {
		ExportedAt time.Time                  `json:"exported_at"`
		Facts      map[string]*FactCollection `json:"facts"`
	}

	if err := json.NewDecoder(r).Decode(&importData); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Import facts
	for machineID, facts := range importData.Facts {
		s.facts[machineID] = facts
	}

	return nil
}

// GetStats returns storage statistics
func (s *MemoryFactStorage) GetStats() (map[string]interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	totalEntries := len(s.facts)
	totalSize := 0

	// Calculate approximate memory usage
	for _, facts := range s.facts {
		if data, err := json.Marshal(facts); err == nil {
			totalSize += len(data)
		}
	}

	return map[string]interface{}{
		"total_entries": totalEntries,
		"total_size":    totalSize,
		"storage_type":  "memory",
	}, nil
}
