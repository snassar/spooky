package facts

import (
	"fmt"
	"io"
)

// HCLFactStorage implements FactStorage using HCL files
type HCLFactStorage struct {
	*BaseFileStorage
}

// NewHCLFactStorage creates a new HCL-based fact storage
func NewHCLFactStorage(filepath string) (*HCLFactStorage, error) {
	base, err := NewBaseFileStorage(filepath, FileStorageFormatHCL)
	if err != nil {
		return nil, err
	}

	return &HCLFactStorage{
		BaseFileStorage: base,
	}, nil
}

// ImportFromHCL imports fact collections from HCL format
func (h *HCLFactStorage) ImportFromHCL(r io.Reader) error {
	// For HCL storage, we can directly import HCL content
	// This is a placeholder implementation
	return fmt.Errorf("HCL import not yet implemented for HCL storage")
}
