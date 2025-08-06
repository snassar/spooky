package facts

import (
	"io"
	"time"

	spookyfactstypes "spooky/internal/facts/types"
	spookystorage "spooky/internal/storage"
)

// FactManager defines the main interface for fact collection and management
type FactManager interface {
	// Core collection operations
	CollectAllFacts(server string) (*spookyfactstypes.FactCollection, error)
	CollectSpecificFacts(server string, keys []string) (*spookyfactstypes.FactCollection, error)
	GetFact(server, key string) (*spookyfactstypes.Fact, error)

	// Storage operations
	PersistFacts(machineID string, collection *spookyfactstypes.FactCollection) error
	LoadPersistedFacts(server string) (*spookyfactstypes.FactCollection, error)
	QueryPersistedFacts(query *spookystorage.FactQuery) ([]*spookyfactstypes.FactCollection, error)
	DeletePersistedFacts(query *spookystorage.FactQuery) (int, error)

	// Export/Import operations
	ExportFacts(w io.Writer) error
	ImportFacts(r io.Reader) error
	ExportFactsWithEncryption(w io.Writer, opts spookystorage.ExportOptions) error
	ImportFactsWithDecryption(r io.Reader, identityFile string) error

	// Cache operations
	ClearCache()
	ClearExpiredCache()
	GetAllFacts() ([]*spookyfactstypes.Fact, error)

	// Configuration
	SetDefaultTTL(ttl time.Duration)
	RegisterCustomCollector(name string, collector spookyfactstypes.FactCollector)

	// Custom facts operations
	ImportCustomFacts(source, server string, mergePolicy spookyfactstypes.MergePolicy) (*spookyfactstypes.FactCollection, error)
	ImportCustomFactsWithOptions(source string, options *spookyfactstypes.ImportOptions) error
	GetCustomFacts(server string) (map[string]interface{}, error)

	// Utility operations
	GenerateMachineID(facts *spookyfactstypes.FactCollection) string
	Close() error

	// Coordinator integration methods
	GetFactCollection(machine string) (*spookyfactstypes.FactCollection, error)
	SetFactCollection(machine string, collection *spookyfactstypes.FactCollection) error

	// Storage access for coordinator integration
	GetStorage() spookystorage.FactStorage
}
