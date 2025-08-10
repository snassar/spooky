package interfaces

import (
	"io"
	"time"

	spookytypesfacts "spooky/internal/types/facts"
	spookystorage "spooky/internal/storage"
)

// FactManager defines the main interface for fact collection and management
type FactManager interface {
	// Core collection operations
	CollectAllFacts(server string) (*spookytypesfacts.FactCollection, error)
	CollectSpecificFacts(server string, keys []string) (*spookytypesfacts.FactCollection, error)
	GetFact(server, key string) (*spookytypesfacts.Fact, error)

	// Storage operations
	PersistFacts(machineID string, collection *spookytypesfacts.FactCollection) error
	LoadPersistedFacts(server string) (*spookytypesfacts.FactCollection, error)
	QueryPersistedFacts(query *spookystorage.FactQuery) ([]*spookytypesfacts.FactCollection, error)
	DeletePersistedFacts(query *spookystorage.FactQuery) (int, error)

	// Export/Import operations
	ExportFacts(w io.Writer) error
	ImportFacts(r io.Reader) error
	ExportFactsWithEncryption(w io.Writer, opts spookystorage.ExportOptions) error
	ImportFactsWithDecryption(r io.Reader, identityFile string) error

	// Cache operations
	ClearCache()
	ClearExpiredCache()
	GetAllFacts() ([]*spookytypesfacts.Fact, error)

	// Configuration
	SetDefaultTTL(ttl time.Duration)
	RegisterCustomCollector(name string, collector spookytypesfacts.FactCollector)

	// Custom facts operations
	ImportCustomFacts(source, server string, mergePolicy spookytypesfacts.MergePolicy) (*spookytypesfacts.FactCollection, error)
	ImportCustomFactsWithOptions(source string, options *spookytypesfacts.ImportOptions) error
	GetCustomFacts(server string) (map[string]interface{}, error)

	// Utility operations
	GenerateMachineID(facts *spookytypesfacts.FactCollection) string
	Close() error

	// Coordinator integration methods
	GetFactCollection(machine string) (*spookytypesfacts.FactCollection, error)
	SetFactCollection(machine string, collection *spookytypesfacts.FactCollection) error

	// Storage access for coordinator integration
	GetStorage() spookystorage.FactStorage
}

// FactCollector defines the interface for fact collection
type FactCollector interface {
	CollectFacts(server string) (*spookytypesfacts.FactCollection, error)
	CollectSpecificFacts(server string, keys []string) (*spookytypesfacts.FactCollection, error)
	ValidateFacts(collection *spookytypesfacts.FactCollection) error
}

// FactProcessor defines the interface for fact processing
type FactProcessor interface {
	ProcessFacts(collection *spookytypesfacts.FactCollection) error
	TransformFacts(collection *spookytypesfacts.FactCollection, transform string) error
	FilterFacts(collection *spookytypesfacts.FactCollection, filter string) error
}

// FactValidator defines the interface for fact validation
type FactValidator interface {
	ValidateFact(fact *spookytypesfacts.Fact) error
	ValidateCollection(collection *spookytypesfacts.FactCollection) error
	ValidateSchema(collection *spookytypesfacts.FactCollection) error
}

// FactCache defines the interface for fact caching
type FactCache interface {
	GetFact(server, key string) (*spookytypesfacts.Fact, error)
	SetFact(server, key string, fact *spookytypesfacts.Fact) error
	DeleteFact(server, key string) error
	ClearCache() error
	ClearExpiredCache() error
}
