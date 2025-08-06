package facts

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	spookyfactstypes "spooky/internal/facts/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewFactStorage tests the creation of new fact storage instances
func TestNewFactStorage(t *testing.T) {
	// Test with temporary directory
	tempDir := t.TempDir()
	storagePath := filepath.Join(tempDir, "test-facts.db")

	storage, err := NewBadgerFactStorage(storagePath)
	require.NoError(t, err)
	require.NotNil(t, storage)

	// Test that the storage can be closed
	err = storage.Close()
	assert.NoError(t, err)
}

// TestFactCollectionSerialization tests the serialization and deserialization of fact collections
func TestFactCollectionSerialization(t *testing.T) {
	// Test FactCollection serialization and deserialization
	now := time.Now()
	collection := &spookyfactstypes.FactCollection{
		Server:    "test-server",
		Timestamp: now,
		Facts: map[string]*spookyfactstypes.Fact{
			"hostname": {
				Key:       "hostname",
				Value:     "test-host",
				Source:    "ssh",
				Server:    "test-server",
				Timestamp: now,
			},
			"os.name": {
				Key:       "os.name",
				Value:     "linux",
				Source:    "ssh",
				Server:    "test-server",
				Timestamp: now,
			},
			"os.version": {
				Key:       "os.version",
				Value:     "20.04",
				Source:    "ssh",
				Server:    "test-server",
				Timestamp: now,
			},
			"hardware.cpu.cores": {
				Key:       "hardware.cpu.cores",
				Value:     4,
				Source:    "ssh",
				Server:    "test-server",
				Timestamp: now,
			},
			"hardware.cpu.model": {
				Key:       "hardware.cpu.model",
				Value:     "Intel i7",
				Source:    "ssh",
				Server:    "test-server",
				Timestamp: now,
			},
			"hardware.memory.total": {
				Key:       "hardware.memory.total",
				Value:     uint64(16 * 1024 * 1024 * 1024), // 16GB
				Source:    "ssh",
				Server:    "test-server",
				Timestamp: now,
			},
			"network.interfaces.eth0.addresses": {
				Key:       "network.interfaces.eth0.addresses",
				Value:     []string{"192.168.1.100"},
				Source:    "ssh",
				Server:    "test-server",
				Timestamp: now,
			},
		},
	}

	// Test that the collection can be serialized and deserialized
	assert.NotNil(t, collection)
	assert.Equal(t, "test-server", collection.Server)
	assert.Equal(t, "test-host", getFactValue(collection, "hostname"))
	assert.Equal(t, "linux", getFactValue(collection, "os.name"))
	assert.Equal(t, "20.04", getFactValue(collection, "os.version"))
}

func getFactValue(collection *spookyfactstypes.FactCollection, key string) string {
	if fact, exists := collection.Facts[key]; exists {
		if value, ok := fact.Value.(string); ok {
			return value
		}
	}
	return ""
}

func getFactIntValue(collection *spookyfactstypes.FactCollection, key string) int {
	if fact, exists := collection.Facts[key]; exists {
		if value, ok := fact.Value.(int); ok {
			return value
		}
	}
	return 0
}

func getFactUint64Value(collection *spookyfactstypes.FactCollection, key string) uint64 {
	if fact, exists := collection.Facts[key]; exists {
		if value, ok := fact.Value.(uint64); ok {
			return value
		}
	}
	return 0
}

func TestFactCollectionCloning(t *testing.T) {
	// Test FactCollection cloning functionality
	now := time.Now()
	original := &spookyfactstypes.FactCollection{
		Server:    "test-server",
		Timestamp: now,
		Facts: map[string]*spookyfactstypes.Fact{
			"hostname": {
				Key:       "hostname",
				Value:     "test-host",
				Source:    "ssh",
				Server:    "test-server",
				Timestamp: now,
			},
			"os.name": {
				Key:       "os.name",
				Value:     "linux",
				Source:    "ssh",
				Server:    "test-server",
				Timestamp: now,
			},
		},
	}

	// Test cloning
	clone := original.Clone()
	assert.NotNil(t, clone)
	assert.Equal(t, original.Server, clone.Server)
	assert.Equal(t, original.Timestamp, clone.Timestamp)
	assert.Equal(t, len(original.Facts), len(clone.Facts))

	// Test that modifying the clone doesn't affect the original
	clone.Facts["hostname"].Value = "modified-host"
	assert.Equal(t, "test-host", getFactValue(original, "hostname"))
	assert.Equal(t, "modified-host", getFactValue(clone, "hostname"))
}

// TestStorageWithTestDataFromExamples tests storage with data from the examples
func TestStorageWithTestDataFromExamples(t *testing.T) {
	// Create temporary storage
	tempDir := t.TempDir()
	storagePath := filepath.Join(tempDir, "test-facts.db")

	storage, err := NewBadgerFactStorage(storagePath)
	require.NoError(t, err)
	defer storage.Close()

	// Create test data similar to what would be collected from real systems
	now := time.Now()
	testCollections := []*spookyfactstypes.FactCollection{
		{
			Server:    "web-server-01",
			Timestamp: now,
			Facts: map[string]*spookyfactstypes.Fact{
				"hostname": {
					Key:       "hostname",
					Value:     "web-server-01",
					Source:    "ssh",
					Server:    "web-server-01",
					Timestamp: now,
				},
				"os.name": {
					Key:       "os.name",
					Value:     "ubuntu",
					Source:    "ssh",
					Server:    "web-server-01",
					Timestamp: now,
				},
				"os.version": {
					Key:       "os.version",
					Value:     "20.04",
					Source:    "ssh",
					Server:    "web-server-01",
					Timestamp: now,
				},
				"hardware.cpu.cores": {
					Key:       "hardware.cpu.cores",
					Value:     8,
					Source:    "ssh",
					Server:    "web-server-01",
					Timestamp: now,
				},
				"hardware.memory.total": {
					Key:       "hardware.memory.total",
					Value:     uint64(32 * 1024 * 1024 * 1024), // 32GB
					Source:    "ssh",
					Server:    "web-server-01",
					Timestamp: now,
				},
			},
		},
		{
			Server:    "db-server-01",
			Timestamp: now,
			Facts: map[string]*spookyfactstypes.Fact{
				"hostname": {
					Key:       "hostname",
					Value:     "db-server-01",
					Source:    "ssh",
					Server:    "db-server-01",
					Timestamp: now,
				},
				"os.name": {
					Key:       "os.name",
					Value:     "centos",
					Source:    "ssh",
					Server:    "db-server-01",
					Timestamp: now,
				},
				"os.version": {
					Key:       "os.version",
					Value:     "8",
					Source:    "ssh",
					Server:    "db-server-01",
					Timestamp: now,
				},
				"hardware.cpu.cores": {
					Key:       "hardware.cpu.cores",
					Value:     16,
					Source:    "ssh",
					Server:    "db-server-01",
					Timestamp: now,
				},
				"hardware.memory.total": {
					Key:       "hardware.memory.total",
					Value:     uint64(64 * 1024 * 1024 * 1024), // 64GB
					Source:    "ssh",
					Server:    "db-server-01",
					Timestamp: now,
				},
			},
		},
	}

	// Store all collections
	for _, collection := range testCollections {
		err := storage.SetFactCollection(collection.Server, collection)
		require.NoError(t, err)
	}

	// Test retrieval
	for _, expectedCollection := range testCollections {
		retrievedCollection, err := storage.GetFactCollection(expectedCollection.Server)
		require.NoError(t, err)
		require.NotNil(t, retrievedCollection)

		assert.Equal(t, expectedCollection.Server, retrievedCollection.Server)
		assert.Equal(t, expectedCollection.Timestamp.Unix(), retrievedCollection.Timestamp.Unix())
		assert.Equal(t, len(expectedCollection.Facts), len(retrievedCollection.Facts))

		// Test specific facts
		assert.Equal(t, getFactValue(expectedCollection, "hostname"), getFactValue(retrievedCollection, "hostname"))
		assert.Equal(t, getFactValue(expectedCollection, "os.name"), getFactValue(retrievedCollection, "os.name"))
		assert.Equal(t, getFactIntValue(expectedCollection, "hardware.cpu.cores"), getFactIntValue(retrievedCollection, "hardware.cpu.cores"))
	}
}

// TestStorageQueryOperations tests various query operations
func TestStorageQueryOperations(t *testing.T) {
	// Create temporary storage
	tempDir := t.TempDir()
	storagePath := filepath.Join(tempDir, "test-facts.db")

	storage, err := NewBadgerFactStorage(storagePath)
	require.NoError(t, err)
	defer storage.Close()

	// Create test data
	now := time.Now()
	testCollections := []*spookyfactstypes.FactCollection{
		{
			Server:    "server-01",
			Timestamp: now,
			Facts: map[string]*spookyfactstypes.Fact{
				"hostname": {
					Key:       "hostname",
					Value:     "server-01",
					Source:    "ssh",
					Server:    "server-01",
					Timestamp: now,
				},
				"os.name": {
					Key:       "os.name",
					Value:     "ubuntu",
					Source:    "ssh",
					Server:    "server-01",
					Timestamp: now,
				},
				"environment": {
					Key:       "environment",
					Value:     "production",
					Source:    "ssh",
					Server:    "server-01",
					Timestamp: now,
				},
			},
		},
		{
			Server:    "server-02",
			Timestamp: now,
			Facts: map[string]*spookyfactstypes.Fact{
				"hostname": {
					Key:       "hostname",
					Value:     "server-02",
					Source:    "ssh",
					Server:    "server-02",
					Timestamp: now,
				},
				"os.name": {
					Key:       "os.name",
					Value:     "centos",
					Source:    "ssh",
					Server:    "server-02",
					Timestamp: now,
				},
				"environment": {
					Key:       "environment",
					Value:     "staging",
					Source:    "ssh",
					Server:    "server-02",
					Timestamp: now,
				},
			},
		},
		{
			Server:    "server-03",
			Timestamp: now,
			Facts: map[string]*spookyfactstypes.Fact{
				"hostname": {
					Key:       "hostname",
					Value:     "server-03",
					Source:    "ssh",
					Server:    "server-03",
					Timestamp: now,
				},
				"os.name": {
					Key:       "os.name",
					Value:     "ubuntu",
					Source:    "ssh",
					Server:    "server-03",
					Timestamp: now,
				},
				"environment": {
					Key:       "environment",
					Value:     "production",
					Source:    "ssh",
					Server:    "server-03",
					Timestamp: now,
				},
			},
		},
	}

	// Store all collections
	for _, collection := range testCollections {
		err := storage.SetFactCollection(collection.Server, collection)
		require.NoError(t, err)
	}

	// Test querying by fact value using the query interface
	query := &spookyfactstypes.FactQuery{
		OS: "ubuntu",
	}
	ubuntuResults, err := storage.QueryFactCollections(query)
	require.NoError(t, err)
	assert.Len(t, ubuntuResults, 2)

	// Test querying by fact value
	query = &spookyfactstypes.FactQuery{
		Environment: "production",
	}
	productionResults, err := storage.QueryFactCollections(query)
	require.NoError(t, err)
	assert.Len(t, productionResults, 2)
}

// TestStorageExportImport tests export and import functionality
func TestStorageExportImport(t *testing.T) {
	// Create temporary storage
	tempDir := t.TempDir()
	storagePath := filepath.Join(tempDir, "test-facts.db")

	storage, err := NewBadgerFactStorage(storagePath)
	require.NoError(t, err)
	defer storage.Close()

	// Create test data
	now := time.Now()
	testCollection := &spookyfactstypes.FactCollection{
		Server:    "test-server",
		Timestamp: now,
		Facts: map[string]*spookyfactstypes.Fact{
			"hostname": {
				Key:       "hostname",
				Value:     "test-host",
				Source:    "ssh",
				Server:    "test-server",
				Timestamp: now,
			},
			"os.name": {
				Key:       "os.name",
				Value:     "linux",
				Source:    "ssh",
				Server:    "test-server",
				Timestamp: now,
			},
		},
	}

	// Store the collection
	err = storage.StoreFacts(testCollection.Server, testCollection)
	require.NoError(t, err)

	// Test export to JSON
	exportPath := filepath.Join(tempDir, "export.json")
	err = storage.ExportToJSON(exportPath)
	require.NoError(t, err)

	// Verify export file exists
	_, err = os.Stat(exportPath)
	assert.NoError(t, err)

	// Test import from JSON
	newStoragePath := filepath.Join(tempDir, "import-facts.db")
	newStorage, err := NewBadgerFactStorage(newStoragePath)
	require.NoError(t, err)
	defer newStorage.Close()

	err = newStorage.ImportFromJSON(exportPath)
	require.NoError(t, err)

	// Verify imported data
	importedCollection, err := newStorage.GetFacts(testCollection.Server)
	require.NoError(t, err)
	require.NotNil(t, importedCollection)

	assert.Equal(t, testCollection.Server, importedCollection.Server)
	assert.Equal(t, testCollection.Timestamp.Unix(), importedCollection.Timestamp.Unix())
	assert.Equal(t, len(testCollection.Facts), len(importedCollection.Facts))
}

// TestStorageDeleteOperations tests delete operations
func TestStorageDeleteOperations(t *testing.T) {
	// Create temporary storage
	tempDir := t.TempDir()
	storagePath := filepath.Join(tempDir, "test-facts.db")

	storage, err := NewBadgerFactStorage(storagePath)
	require.NoError(t, err)
	defer storage.Close()

	// Create test data
	now := time.Now()
	testCollection := &spookyfactstypes.FactCollection{
		Server:    "test-server",
		Timestamp: now,
		Facts: map[string]*spookyfactstypes.Fact{
			"hostname": {
				Key:       "hostname",
				Value:     "test-host",
				Source:    "ssh",
				Server:    "test-server",
				Timestamp: now,
			},
			"os.name": {
				Key:       "os.name",
				Value:     "linux",
				Source:    "ssh",
				Server:    "test-server",
				Timestamp: now,
			},
		},
	}

	// Store the collection
	err = storage.StoreFacts(testCollection.Server, testCollection)
	require.NoError(t, err)

	// Verify it exists
	retrievedCollection, err := storage.GetFacts(testCollection.Server)
	require.NoError(t, err)
	assert.NotNil(t, retrievedCollection)

	// Test deleting specific facts
	err = storage.DeleteFact(testCollection.Server, "hostname")
	require.NoError(t, err)

	// Verify the fact was deleted
	retrievedCollection, err = storage.GetFacts(testCollection.Server)
	require.NoError(t, err)
	assert.NotNil(t, retrievedCollection)
	assert.NotContains(t, retrievedCollection.Facts, "hostname")
	assert.Contains(t, retrievedCollection.Facts, "os.name")

	// Test deleting all facts for a server
	err = storage.DeleteFacts(testCollection.Server)
	require.NoError(t, err)

	// Verify all facts were deleted
	retrievedCollection, err = storage.GetFacts(testCollection.Server)
	require.NoError(t, err)
	assert.Nil(t, retrievedCollection)
}
