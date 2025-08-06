package facts

import (
	"io"
	"time"

	"spooky/internal/facts/types"
)

// FactManager defines the main interface for fact collection and management
type FactManager interface {
	// Core collection operations
	CollectAllFacts(server string) (*types.FactCollection, error)
	CollectSpecificFacts(server string, keys []string) (*types.FactCollection, error)
	GetFact(server, key string) (*types.Fact, error)

	// Storage operations
	PersistFacts(machineID string, collection *types.FactCollection) error
	LoadPersistedFacts(server string) (*types.FactCollection, error)
	QueryPersistedFacts(query *FactQuery) ([]*types.FactCollection, error)
	DeletePersistedFacts(query *FactQuery) (int, error)

	// Export/Import operations
	ExportFacts(w io.Writer) error
	ImportFacts(r io.Reader) error
	ExportFactsWithEncryption(w io.Writer, opts ExportOptions) error
	ImportFactsWithDecryption(r io.Reader, identityFile string) error

	// Cache operations
	ClearCache()
	ClearExpiredCache()
	GetAllFacts() ([]*types.Fact, error)

	// Configuration
	SetDefaultTTL(ttl time.Duration)
	RegisterCustomCollector(name string, collector types.FactCollector)

	// Custom facts operations
	ImportCustomFacts(source, server string, mergePolicy types.MergePolicy) (*types.FactCollection, error)
	ImportCustomFactsWithOptions(source string, options *types.ImportOptions) error
	GetCustomFacts(server string) (map[string]interface{}, error)

	// Utility operations
	GenerateMachineID(facts *types.FactCollection) string
	Close() error

	// Coordinator integration methods
	GetFactCollection(machine string) (*types.FactCollection, error)
	SetFactCollection(machine string, collection *types.FactCollection) error

	// Storage access for coordinator integration
	GetStorage() FactStorage
}
