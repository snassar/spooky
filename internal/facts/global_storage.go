package facts

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	spookybadger "spooky/internal/facts/storage/badger"
	spookyfactstypes "spooky/internal/facts/types"
)

// GlobalFactsStorage manages global facts storage with XDG compliance
type GlobalFactsStorage struct {
	storage    FactStorage
	basePath   string
	schemaPath string
}

// NewGlobalFactsStorage creates a new global facts storage instance
func NewGlobalFactsStorage() (*GlobalFactsStorage, error) {
	// Get XDG state directory
	stateDir := os.Getenv("XDG_STATE_HOME")
	if stateDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		stateDir = filepath.Join(homeDir, ".local", "state")
	}

	// Create spooky directory structure
	spookyStateDir := filepath.Join(stateDir, "spooky")
	if err := os.MkdirAll(spookyStateDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create spooky state directory: %w", err)
	}

	// Initialize BadgerDB storage for global facts
	storage, err := spookybadger.NewBadgerFactStorage(filepath.Join(spookyStateDir, "global-facts.db"))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize global facts storage: %w", err)
	}

	return &GlobalFactsStorage{
		storage:    storage,
		basePath:   spookyStateDir,
		schemaPath: "global-facts-badger.hcl",
	}, nil
}

// StoreFactCollection stores a fact collection in global storage
func (gfs *GlobalFactsStorage) StoreFactCollection(machineID string, collection *spookyfactstypes.FactCollection) error {
	// Validate machine ID format (32-character hex)
	if err := gfs.validateMachineID(machineID); err != nil {
		return fmt.Errorf("invalid machine ID: %w", err)
	}

	// Validate facts against global schema
	if err := gfs.validateGlobalFacts(collection); err != nil {
		return fmt.Errorf("facts validation failed: %w", err)
	}

	// Check for machine ID collisions
	if err := gfs.checkCollision(machineID, collection); err != nil {
		return fmt.Errorf("machine ID collision detected: %w", err)
	}

	// Store facts with global TTL (24h default)
	collection.Timestamp = time.Now()
	for _, fact := range collection.Facts {
		if fact.TTL == 0 {
			fact.TTL = 24 * time.Hour // Global facts default TTL
		}
	}

	return gfs.storage.SetFactCollection(machineID, collection)
}

// GetFactCollection retrieves a fact collection from global storage
func (gfs *GlobalFactsStorage) GetFactCollection(machineID string) (*spookyfactstypes.FactCollection, error) {
	if err := gfs.validateMachineID(machineID); err != nil {
		return nil, fmt.Errorf("invalid machine ID: %w", err)
	}

	collection, err := gfs.storage.GetFactCollection(machineID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve fact collection: %w", err)
	}

	// Check if facts are expired (basic check)
	if gfs.isExpired(collection) {
		return nil, fmt.Errorf("fact collection has expired")
	}

	return collection, nil
}

// QueryFactCollections queries fact collections with filters
func (gfs *GlobalFactsStorage) QueryFactCollections(query *spookyfactstypes.FactQuery) ([]*spookyfactstypes.FactCollection, error) {
	// Validate query parameters
	if err := gfs.validateQuery(query); err != nil {
		return nil, fmt.Errorf("invalid query: %w", err)
	}

	return gfs.storage.QueryFactCollections(query)
}

// DeleteFactCollection deletes a fact collection from global storage
func (gfs *GlobalFactsStorage) DeleteFactCollection(machineID string) error {
	if err := gfs.validateMachineID(machineID); err != nil {
		return fmt.Errorf("invalid machine ID: %w", err)
	}

	return gfs.storage.DeleteFactCollection(machineID)
}

// DeleteFactCollections deletes fact collections matching query criteria
func (gfs *GlobalFactsStorage) DeleteFactCollections(query *spookyfactstypes.FactQuery) (int, error) {
	// Validate query parameters
	if err := gfs.validateQuery(query); err != nil {
		return 0, fmt.Errorf("invalid query: %w", err)
	}

	return gfs.storage.DeleteFactCollections(query)
}

// UpdateExpiredFacts updates facts that have expired based on TTL
func (gfs *GlobalFactsStorage) UpdateExpiredFacts() error {
	// Query all fact collections
	query := &spookyfactstypes.FactQuery{
		Limit: 1000, // Reasonable limit for global facts
	}

	collections, err := gfs.storage.QueryFactCollections(query)
	if err != nil {
		return fmt.Errorf("failed to query fact collections: %w", err)
	}

	// Check for expired facts and mark them for update
	for _, collection := range collections {
		if gfs.isExpired(collection) {
			// Mark facts as expired by updating timestamp
			collection.Timestamp = time.Now().Add(-25 * time.Hour) // Force expiration
			machineID := gfs.extractMachineID(collection)
			if machineID != "" {
				if err := gfs.storage.SetFactCollection(machineID, collection); err != nil {
					return fmt.Errorf("failed to mark expired facts: %w", err)
				}
			}
		}
	}

	return nil
}

// GetStoragePath returns the storage path for debugging/info
func (gfs *GlobalFactsStorage) GetStoragePath() string {
	return gfs.basePath
}

// Close closes the global facts storage
func (gfs *GlobalFactsStorage) Close() error {
	return gfs.storage.Close()
}

// validateMachineID validates machine ID format
func (gfs *GlobalFactsStorage) validateMachineID(machineID string) error {
	if machineID == "" {
		return fmt.Errorf("machine ID cannot be empty")
	}

	// Check 32-character hex pattern from schema
	if len(machineID) != 32 {
		return fmt.Errorf("machine ID must be 32 characters long")
	}

	for _, char := range machineID {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
			return fmt.Errorf("machine ID must contain only hexadecimal characters")
		}
	}

	return nil
}

// validateGlobalFacts validates facts against global schema
func (gfs *GlobalFactsStorage) validateGlobalFacts(facts *spookyfactstypes.FactCollection) error {
	utils := NewValidationUtils()
	return utils.ValidateFactsCollection(facts)
}

// checkCollision checks for machine ID collisions
func (gfs *GlobalFactsStorage) checkCollision(machineID string, newCollection *spookyfactstypes.FactCollection) error {
	existingCollection, err := gfs.storage.GetFactCollection(machineID)
	if err != nil {
		// No existing facts, no collision
		return nil
	}

	// Check if the existing facts are from a different machine
	if existingCollection.Server != newCollection.Server {
		return fmt.Errorf("machine ID collision: %s is already used by %s", machineID, existingCollection.Server)
	}

	return nil
}

// validateQuery validates query parameters
func (gfs *GlobalFactsStorage) validateQuery(query *spookyfactstypes.FactQuery) error {
	if query == nil {
		return fmt.Errorf("query cannot be nil")
	}

	if query.Limit < 0 {
		return fmt.Errorf("query limit cannot be negative")
	}

	if query.Limit > 10000 {
		return fmt.Errorf("query limit cannot exceed 10000")
	}

	return nil
}

// isExpired checks if a fact collection has expired
func (gfs *GlobalFactsStorage) isExpired(collection *spookyfactstypes.FactCollection) bool {
	if collection == nil {
		return true
	}

	// Check if any facts have expired
	for _, fact := range collection.Facts {
		if fact.TTL > 0 && time.Since(fact.Timestamp) > fact.TTL {
			return true
		}
	}

	return false
}

// extractMachineID extracts machine ID from a fact collection
func (gfs *GlobalFactsStorage) extractMachineID(collection *spookyfactstypes.FactCollection) string {
	if collection == nil {
		return ""
	}

	if machineID, exists := collection.Facts["machine_id"]; exists {
		if id, ok := machineID.Value.(string); ok {
			return id
		}
	}

	return ""
}
