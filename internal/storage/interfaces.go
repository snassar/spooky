package storage

import (
	"io"
	"time"

	spookyfactstypes "spooky/internal/facts/types"
)

// StorageType defines the type of storage backend
type StorageType string

const (
	StorageTypeBadger StorageType = "badger"
	StorageTypeJSON   StorageType = "json"
	StorageTypeHCL    StorageType = "hcl"
)

// StorageBackend defines the basic storage operations
type StorageBackend interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	Close() error
}

// ExportableStorage defines operations for data export and import
type ExportableStorage interface {
	ExportToJSON(w io.Writer) error
	ExportToHCL(w io.Writer) error
	ImportFromJSON(r io.Reader) error
	ImportFromHCL(r io.Reader) error
	ExportToJSONWithEncryption(w io.Writer, options spookyfactstypes.ExportOptions) error
	ImportFromJSONWithDecryption(r io.Reader, identityFile string) error
}

// FactStorage defines the interface for fact-specific storage operations
type FactStorage interface {
	StorageBackend
	ExportableStorage

	// Fact-specific operations
	SetFactCollection(machineID string, collection *spookyfactstypes.FactCollection) error
	GetFactCollection(machineID string) (*spookyfactstypes.FactCollection, error)
	QueryFactCollections(query *spookyfactstypes.FactQuery) ([]*spookyfactstypes.FactCollection, error)
	DeleteFactCollection(machineID string) error
	DeleteFactCollections(query *spookyfactstypes.FactQuery) (int, error)
}

// StorageOptions defines configuration options for storage backends
type StorageOptions struct {
	Type           StorageType
	Path           string
	EncryptEnabled bool
	CryptoManager  interface{} // Will be properly typed when secrets package is available
	DefaultTTL     time.Duration
}

// StorageFactory creates storage instances
type StorageFactory interface {
	NewBadgerStorage(dbPath string) (FactStorage, error)
	NewBadgerStorageWithCrypto(dbPath string, cryptoManager interface{}, encryptEnabled bool) (FactStorage, error)
	NewBadgerStorageReadOnly(dbPath string) (FactStorage, error)
	NewJSONStorage(filePath string) (FactStorage, error)
	NewHCLStorage(filePath string) (FactStorage, error)
}
