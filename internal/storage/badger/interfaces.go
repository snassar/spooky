package badger

import (
	"io"

	spookyfactstypes "spooky/internal/facts/types"
)

// FactStorage defines the interface for persistent fact storage
type FactStorage interface {
	// Core fact collection operations
	SetFactCollection(machineID string, collection *spookyfactstypes.FactCollection) error
	GetFactCollection(machineID string) (*spookyfactstypes.FactCollection, error)
	QueryFactCollections(query *spookyfactstypes.FactQuery) ([]*spookyfactstypes.FactCollection, error)
	DeleteFactCollection(machineID string) error
	DeleteFactCollections(query *spookyfactstypes.FactQuery) (int, error)

	// Import/Export operations
	ExportToJSON(w io.Writer) error
	ImportFromJSON(r io.Reader) error
	ImportFromHCL(r io.Reader) error
	ExportToJSONWithEncryption(w io.Writer, opts spookyfactstypes.ExportOptions) error
	ImportFromJSONWithDecryption(r io.Reader, identityFile string) error

	// Close the storage connection
	Close() error
}

// ExportOptions defines options for encrypted export
type ExportOptions = spookyfactstypes.ExportOptions

// FactQuery defines query parameters for fact collections
type FactQuery = spookyfactstypes.FactQuery
