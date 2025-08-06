package processors

import (
	"io"
	spookyfactstypes "spooky/internal/facts/types"
)

// Processor defines the interface for fact processing operations
type Processor interface {
	// Merging operations
	MergeCollections(collections ...*spookyfactstypes.FactCollection) (*spookyfactstypes.FactCollection, error)
	MergeWithPolicy(existing, new *spookyfactstypes.FactCollection, policy spookyfactstypes.MergePolicy) (*spookyfactstypes.FactCollection, error)

	// Validation operations
	ValidateCollection(collection *spookyfactstypes.FactCollection) error
	ValidateFact(fact *spookyfactstypes.Fact) error

	// Export operations
	ExportToJSON(collections []*spookyfactstypes.FactCollection, w io.Writer) error
	ExportToHCL(collections []*spookyfactstypes.FactCollection, w io.Writer) error

	// Import operations
	ImportFromJSON(r io.Reader) ([]*spookyfactstypes.FactCollection, error)
	ImportFromHCL(r io.Reader) ([]*spookyfactstypes.FactCollection, error)
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
