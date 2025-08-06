package facts

import (
	"fmt"
	"io"
)

// JSONFactStorage implements FactStorage using JSON files
type JSONFactStorage struct {
	*BaseFileStorage
}

// NewJSONFactStorage creates a new JSON-based fact storage
func NewJSONFactStorage(filepath string) (*JSONFactStorage, error) {
	base, err := NewBaseFileStorage(filepath, FileStorageFormatJSON)
	if err != nil {
		return nil, err
	}

	return &JSONFactStorage{
		BaseFileStorage: base,
	}, nil
}

// ImportFromHCL imports fact collections from HCL format
func (j *JSONFactStorage) ImportFromHCL(r io.Reader) error {
	// For JSON storage, we need to convert HCL to JSON format first
	// This is a placeholder implementation
	return fmt.Errorf("HCL import not yet implemented for JSON storage")
}
