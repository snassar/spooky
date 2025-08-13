package facts

import (
	"context"
	"encoding/json"
	"io"
	"sync"
	"time"
)

// MemoryFactStorage provides minimal storage for debugging and statistics
type MemoryFactStorage struct {
	mutex sync.RWMutex
}

// NewMemoryFactStorage creates a new minimal fact storage
func NewMemoryFactStorage() *MemoryFactStorage {
	return &MemoryFactStorage{}
}

// ExportToJSON exports facts to JSON format (for direct export)
func (s *MemoryFactStorage) ExportToJSON(w io.Writer, facts []*FactCollection) error {
	export := struct {
		ExportedAt time.Time         `json:"exported_at"`
		Facts      []*FactCollection `json:"facts"`
	}{
		ExportedAt: time.Now(),
		Facts:      facts,
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(export)
}

// GetStats returns storage statistics for debugging
func (s *MemoryFactStorage) GetStats() (map[string]interface{}, error) {
	return map[string]interface{}{
		"storage_type": "minimal",
		"description":  "Direct export without intermediate storage",
	}, nil
}

// Clear removes all facts from memory
func (s *MemoryFactStorage) Clear(ctx context.Context) error {
	// Memory storage is ephemeral - no cleanup needed
	return nil
}
