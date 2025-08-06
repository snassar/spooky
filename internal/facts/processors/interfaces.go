package processors

import (
	"io"
	"spooky/internal/facts/types"
)

// Processor defines the interface for fact processing operations
type Processor interface {
	// Merging operations
	MergeCollections(collections ...*types.FactCollection) (*types.FactCollection, error)
	MergeWithPolicy(existing, new *types.FactCollection, policy types.MergePolicy) (*types.FactCollection, error)

	// Validation operations
	ValidateCollection(collection *types.FactCollection) error
	ValidateFact(fact *types.Fact) error

	// Export operations
	ExportToJSON(collections []*types.FactCollection, w io.Writer) error
	ExportToHCL(collections []*types.FactCollection, w io.Writer) error

	// Import operations
	ImportFromJSON(r io.Reader) ([]*types.FactCollection, error)
	ImportFromHCL(r io.Reader) ([]*types.FactCollection, error)
}

// Manager provides processor management functionality
type Manager struct {
	processor Processor
}

// NewManager creates a new processor manager
func NewManager(processor Processor) *Manager {
	return &Manager{
		processor: processor,
	}
}

// GetProcessor returns the underlying processor
func (m *Manager) GetProcessor() Processor {
	return m.processor
}

// SetProcessor sets the underlying processor
func (m *Manager) SetProcessor(processor Processor) {
	m.processor = processor
}
