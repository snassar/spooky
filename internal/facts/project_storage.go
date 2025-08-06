package facts

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	spookyfactstypes "spooky/internal/facts/types"
	spookystorage "spooky/internal/storage"
	spookystoragebadger "spooky/internal/storage/badger"
)

// ProjectFactsStorage manages project-specific facts storage
type ProjectFactsStorage struct {
	storage    spookystorage.FactStorage
	projectDir string
	projectID  string
	schemaPath string
}

// NewProjectFactsStorage creates a new project facts storage instance
func NewProjectFactsStorage(projectDir string) (*ProjectFactsStorage, error) {
	// Validate project directory
	if projectDir == "" {
		return nil, fmt.Errorf("project directory cannot be empty")
	}

	// Check if project directory exists
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("project directory does not exist: %s", projectDir)
	}

	// Extract project ID from project directory name
	projectID := filepath.Base(projectDir)
	if projectID == "" {
		return nil, fmt.Errorf("invalid project directory name")
	}

	// Create facts database path within project
	factsDBPath := filepath.Join(projectDir, "facts.db")

	// Initialize BadgerDB storage for project facts
	storage, err := spookystoragebadger.NewBadgerFactStorage(factsDBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize project facts storage: %w", err)
	}

	return &ProjectFactsStorage{
		storage:    storage,
		projectDir: projectDir,
		projectID:  projectID,
		schemaPath: "facts-storage.schema.hcl",
	}, nil
}

// UpdateExpiredFacts updates facts that have expired based on TTL
func (pfs *ProjectFactsStorage) UpdateExpiredFacts() error {
	// Query all fact collections in this project
	query := &spookyfactstypes.FactQuery{
		Limit: 1000, // Reasonable limit for project facts
	}

	collections, err := pfs.storage.QueryFactCollections(query)
	if err != nil {
		return fmt.Errorf("failed to query fact collections: %w", err)
	}

	// Check for expired facts and mark them for update
	for _, collection := range collections {
		if pfs.belongsToProject(collection) && pfs.isExpired(collection) {
			// Mark facts as expired by updating timestamp
			collection.Timestamp = time.Now().Add(-2 * time.Hour) // Force expiration
			machineID := pfs.extractMachineID(collection)
			if machineID != "" {
				if err := pfs.storage.SetFactCollection(machineID, collection); err != nil {
					return fmt.Errorf("failed to mark expired facts: %w", err)
				}
			}
		}
	}

	return nil
}

// GetProjectID returns the project ID
func (pfs *ProjectFactsStorage) GetProjectID() string {
	return pfs.projectID
}

// GetProjectDir returns the project directory
func (pfs *ProjectFactsStorage) GetProjectDir() string {
	return pfs.projectDir
}

// GetStoragePath returns the storage path for debugging/info
func (pfs *ProjectFactsStorage) GetStoragePath() string {
	return filepath.Join(pfs.projectDir, "facts.db")
}

// Close closes the project facts storage
func (pfs *ProjectFactsStorage) Close() error {
	return pfs.storage.Close()
}

// validateMachineID validates machine ID format
func (pfs *ProjectFactsStorage) validateMachineID(machineID string) error {
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

// validateProjectFacts validates facts against project schema
func (pfs *ProjectFactsStorage) validateProjectFacts(facts *spookyfactstypes.FactCollection) error {
	utils := NewValidationUtils()
	return utils.ValidateFactsCollection(facts)
}

// StoreFactCollection stores a fact collection in project storage
func (pfs *ProjectFactsStorage) StoreFactCollection(machineID string, collection *spookyfactstypes.FactCollection) error {
	// Validate machine ID format (32-character hex)
	if err := pfs.validateMachineID(machineID); err != nil {
		return fmt.Errorf("invalid machine ID: %w", err)
	}

	// Validate facts against project schema
	if err := pfs.validateProjectFacts(collection); err != nil {
		return fmt.Errorf("facts validation failed: %w", err)
	}

	// Check for machine ID collisions within project
	if err := pfs.checkProjectCollision(machineID, collection); err != nil {
		return fmt.Errorf("machine ID collision detected in project: %w", err)
	}

	// Add project metadata to facts
	pfs.addProjectMetadata(collection)

	// Store facts with project TTL (1h default)
	collection.Timestamp = time.Now()
	for _, fact := range collection.Facts {
		if fact.TTL == 0 {
			fact.TTL = time.Hour // Project facts default TTL
		}
	}

	return pfs.storage.SetFactCollection(machineID, collection)
}

// GetFactCollection retrieves a fact collection from project storage
func (pfs *ProjectFactsStorage) GetFactCollection(machineID string) (*spookyfactstypes.FactCollection, error) {
	if err := pfs.validateMachineID(machineID); err != nil {
		return nil, fmt.Errorf("invalid machine ID: %w", err)
	}

	collection, err := pfs.storage.GetFactCollection(machineID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve fact collection: %w", err)
	}

	// Verify facts belong to this project
	if !pfs.belongsToProject(collection) {
		return nil, fmt.Errorf("fact collection does not belong to project %s", pfs.projectID)
	}

	// Check if facts are expired
	if pfs.isExpired(collection) {
		return nil, fmt.Errorf("fact collection has expired")
	}

	return collection, nil
}

// QueryFactCollections queries fact collections with project-specific filters
func (pfs *ProjectFactsStorage) QueryFactCollections(query *spookyfactstypes.FactQuery) ([]*spookyfactstypes.FactCollection, error) {
	// Validate query parameters
	if err := pfs.validateQuery(query); err != nil {
		return nil, fmt.Errorf("invalid query: %w", err)
	}

	collections, err := pfs.storage.QueryFactCollections(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query fact collections: %w", err)
	}

	// Filter to only include facts from this project
	var projectCollections []*spookyfactstypes.FactCollection
	for _, collection := range collections {
		if pfs.belongsToProject(collection) {
			projectCollections = append(projectCollections, collection)
		}
	}

	return projectCollections, nil
}

// DeleteFactCollection deletes a fact collection from project storage
func (pfs *ProjectFactsStorage) DeleteFactCollection(machineID string) error {
	if err := pfs.validateMachineID(machineID); err != nil {
		return fmt.Errorf("invalid machine ID: %w", err)
	}

	// Verify facts belong to this project before deletion
	collection, err := pfs.storage.GetFactCollection(machineID)
	if err == nil && !pfs.belongsToProject(collection) {
		return fmt.Errorf("fact collection does not belong to project %s", pfs.projectID)
	}

	return pfs.storage.DeleteFactCollection(machineID)
}

// DeleteFactCollections deletes fact collections matching query criteria
func (pfs *ProjectFactsStorage) DeleteFactCollections(query *spookyfactstypes.FactQuery) (int, error) {
	// Validate query parameters
	if err := pfs.validateQuery(query); err != nil {
		return 0, fmt.Errorf("invalid query: %w", err)
	}

	collections, err := pfs.storage.QueryFactCollections(query)
	if err != nil {
		return 0, fmt.Errorf("failed to query fact collections: %w", err)
	}

	// Only delete facts from this project
	var projectCollections []*spookyfactstypes.FactCollection
	for _, collection := range collections {
		if pfs.belongsToProject(collection) {
			projectCollections = append(projectCollections, collection)
		}
	}

	// Delete project-specific collections
	deletedCount := 0
	for _, collection := range projectCollections {
		// Extract machine ID from collection
		machineID := pfs.extractMachineID(collection)
		if machineID != "" {
			if err := pfs.storage.DeleteFactCollection(machineID); err != nil {
				return deletedCount, fmt.Errorf("failed to delete fact collection: %w", err)
			}
			deletedCount++
		}
	}

	return deletedCount, nil
}

// checkProjectCollision checks for machine ID collisions within the project
func (pfs *ProjectFactsStorage) checkProjectCollision(machineID string, newCollection *spookyfactstypes.FactCollection) error {
	existingCollection, err := pfs.storage.GetFactCollection(machineID)
	if err != nil {
		// No existing facts, no collision
		return nil
	}

	// Only check collision if existing facts belong to this project
	if !pfs.belongsToProject(existingCollection) {
		return nil
	}

	// Check if the existing facts are from a different machine
	if existingCollection.Server != newCollection.Server {
		return fmt.Errorf("machine ID collision in project: %s is already used by %s", machineID, existingCollection.Server)
	}

	return nil
}

// belongsToProject checks if a fact collection belongs to this project
func (pfs *ProjectFactsStorage) belongsToProject(collection *spookyfactstypes.FactCollection) bool {
	if collection == nil {
		return false
	}

	// Check if collection has project metadata
	if projectID, exists := collection.Facts["project.id"]; exists {
		if id, ok := projectID.Value.(string); ok && id == pfs.projectID {
			return true
		}
	}

	return false
}

// extractMachineID extracts machine ID from a fact collection
func (pfs *ProjectFactsStorage) extractMachineID(collection *spookyfactstypes.FactCollection) string {
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

// validateQuery validates query parameters
func (pfs *ProjectFactsStorage) validateQuery(query *spookyfactstypes.FactQuery) error {
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
func (pfs *ProjectFactsStorage) isExpired(collection *spookyfactstypes.FactCollection) bool {
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

// addProjectMetadata adds project-specific metadata to facts
func (pfs *ProjectFactsStorage) addProjectMetadata(facts *spookyfactstypes.FactCollection) {
	for _, fact := range facts.Facts {
		if fact.Metadata == nil {
			fact.Metadata = make(map[string]interface{})
		}
		fact.Metadata["project_id"] = pfs.projectID
		fact.Metadata["project_dir"] = pfs.projectDir
	}
}

// GetProjectFactsSummary returns a summary of facts in the project
func (pfs *ProjectFactsStorage) GetProjectFactsSummary() (*ProjectFactsSummary, error) {
	query := &spookyfactstypes.FactQuery{
		Limit: 10000,
	}

	collections, err := pfs.storage.QueryFactCollections(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query project facts: %w", err)
	}

	// Filter to only include facts from this project
	var projectCollections []*spookyfactstypes.FactCollection
	for _, collection := range collections {
		if pfs.belongsToProject(collection) {
			projectCollections = append(projectCollections, collection)
		}
	}

	summary := &ProjectFactsSummary{
		ProjectID:    pfs.projectID,
		ProjectDir:   pfs.projectDir,
		MachineCount: len(projectCollections),
		TotalFacts:   0,
		ExpiredFacts: 0,
		LastUpdated:  time.Time{},
	}

	for _, collection := range projectCollections {
		// Count actual facts
		summary.TotalFacts += len(collection.Facts)

		if collection.Timestamp.After(summary.LastUpdated) {
			summary.LastUpdated = collection.Timestamp
		}

		if pfs.isExpired(collection) {
			summary.ExpiredFacts++
		}
	}

	return summary, nil
}

// ProjectFactsSummary provides a summary of project facts
type ProjectFactsSummary struct {
	ProjectID    string    `json:"project_id"`
	ProjectDir   string    `json:"project_dir"`
	MachineCount int       `json:"machine_count"`
	TotalFacts   int       `json:"total_facts"`
	ExpiredFacts int       `json:"expired_facts"`
	LastUpdated  time.Time `json:"last_updated"`
}
