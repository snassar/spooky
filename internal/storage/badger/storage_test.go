package badger

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	spookyfactstypes "spooky/internal/types/facts"
)

func TestNewBadgerFactStorage(t *testing.T) {
	// Test creating new storage
	storage, err := NewBadgerFactStorage(":memory:")
	if err != nil {
		t.Fatalf("Failed to create BadgerFactStorage: %v", err)
	}
	defer storage.Close()

	if storage.db == nil {
		t.Fatal("Database should not be nil")
	}
}

func TestSetAndGetFactCollection(t *testing.T) {
	storage, err := NewBadgerFactStorage(":memory:")
	if err != nil {
		t.Fatalf("Failed to create BadgerFactStorage: %v", err)
	}
	defer storage.Close()

	// Create a test fact collection
	collection := &spookyfactstypes.FactCollection{
		Server:    "test-server",
		Timestamp: time.Now(),
		Facts: map[string]*spookyfactstypes.Fact{
			"os.name": {
				Key:   "os.name",
				Value: "Ubuntu",
			},
			"os.version": {
				Key:   "os.version",
				Value: "22.04",
			},
		},
	}

	// Store the collection
	err = storage.SetFactCollection("test-machine", collection)
	if err != nil {
		t.Fatalf("Failed to set fact collection: %v", err)
	}

	// Retrieve the collection
	retrieved, err := storage.GetFactCollection("test-machine")
	if err != nil {
		t.Fatalf("Failed to get fact collection: %v", err)
	}

	// Verify the data
	if retrieved.Server != collection.Server {
		t.Errorf("Expected server %s, got %s", collection.Server, retrieved.Server)
	}

	if len(retrieved.Facts) != len(collection.Facts) {
		t.Errorf("Expected %d facts, got %d", len(collection.Facts), len(retrieved.Facts))
	}

	// Check specific facts
	if osName, exists := retrieved.Facts["os.name"]; !exists {
		t.Error("os.name fact not found")
	} else if osName.Value != "Ubuntu" {
		t.Errorf("Expected os.name 'Ubuntu', got %v", osName.Value)
	}
}

func TestQueryFactCollections(t *testing.T) {
	storage, err := NewBadgerFactStorage(":memory:")
	if err != nil {
		t.Fatalf("Failed to create BadgerFactStorage: %v", err)
	}
	defer storage.Close()

	// Create test collections
	collection1 := &spookyfactstypes.FactCollection{
		Server:    "server1",
		Timestamp: time.Now(),
		Facts: map[string]*spookyfactstypes.Fact{
			"os.name": {
				Key:   "os.name",
				Value: "Ubuntu",
			},
		},
	}

	collection2 := &spookyfactstypes.FactCollection{
		Server:    "server2",
		Timestamp: time.Now(),
		Facts: map[string]*spookyfactstypes.Fact{
			"os.name": {
				Key:   "os.name",
				Value: "CentOS",
			},
		},
	}

	// Store collections
	err = storage.SetFactCollection("machine1", collection1)
	if err != nil {
		t.Fatalf("Failed to set fact collection 1: %v", err)
	}

	err = storage.SetFactCollection("machine2", collection2)
	if err != nil {
		t.Fatalf("Failed to set fact collection 2: %v", err)
	}

	// Query all collections
	query := &spookyfactstypes.FactQuery{}
	results, err := storage.QueryFactCollections(query)
	if err != nil {
		t.Fatalf("Failed to query fact collections: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// Query by OS
	query = &spookyfactstypes.FactQuery{OS: "Ubuntu"}
	results, err = storage.QueryFactCollections(query)
	if err != nil {
		t.Fatalf("Failed to query fact collections by OS: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result for Ubuntu, got %d", len(results))
	}

	if results[0].Server != "server1" {
		t.Errorf("Expected server1, got %s", results[0].Server)
	}
}

func TestExportToJSON(t *testing.T) {
	storage, err := NewBadgerFactStorage(":memory:")
	if err != nil {
		t.Fatalf("Failed to create BadgerFactStorage: %v", err)
	}
	defer storage.Close()

	// Create a test collection
	collection := &spookyfactstypes.FactCollection{
		Server:    "test-server",
		Timestamp: time.Now(),
		Facts: map[string]*spookyfactstypes.Fact{
			"os.name": {
				Key:   "os.name",
				Value: "Ubuntu",
			},
		},
	}

	// Store the collection
	err = storage.SetFactCollection("test-machine", collection)
	if err != nil {
		t.Fatalf("Failed to set fact collection: %v", err)
	}

	// Export to JSON
	var buf bytes.Buffer
	err = storage.ExportToJSON(&buf)
	if err != nil {
		t.Fatalf("Failed to export to JSON: %v", err)
	}

	// Verify the JSON output
	var exported map[string]*spookyfactstypes.FactCollection
	err = json.Unmarshal(buf.Bytes(), &exported)
	if err != nil {
		t.Fatalf("Failed to unmarshal exported JSON: %v", err)
	}

	if len(exported) != 1 {
		t.Errorf("Expected 1 exported collection, got %d", len(exported))
	}

	exportedCollection, exists := exported["test-server"]
	if !exists {
		t.Fatal("Expected collection 'test-server' not found in export")
	}

	if exportedCollection.Server != "test-server" {
		t.Errorf("Expected server 'test-server', got %s", exportedCollection.Server)
	}
}

func TestImportFromJSON(t *testing.T) {
	storage, err := NewBadgerFactStorage(":memory:")
	if err != nil {
		t.Fatalf("Failed to create BadgerFactStorage: %v", err)
	}
	defer storage.Close()

	// Create test data
	testData := map[string]*spookyfactstypes.FactCollection{
		"imported-server": {
			Server:    "imported-server",
			Timestamp: time.Now(),
			Facts: map[string]*spookyfactstypes.Fact{
				"os.name": {
					Key:   "os.name",
					Value: "Debian",
				},
			},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	// Import from JSON
	reader := bytes.NewReader(jsonData)
	err = storage.ImportFromJSON(reader)
	if err != nil {
		t.Fatalf("Failed to import from JSON: %v", err)
	}

	// Verify the imported data
	imported, err := storage.GetFactCollection("imported-server")
	if err != nil {
		t.Fatalf("Failed to get imported collection: %v", err)
	}

	if imported.Server != "imported-server" {
		t.Errorf("Expected server 'imported-server', got %s", imported.Server)
	}

	if osName, exists := imported.Facts["os.name"]; !exists {
		t.Error("os.name fact not found in imported collection")
	} else if osName.Value != "Debian" {
		t.Errorf("Expected os.name 'Debian', got %v", osName.Value)
	}
}

func TestDeleteFactCollection(t *testing.T) {
	storage, err := NewBadgerFactStorage(":memory:")
	if err != nil {
		t.Fatalf("Failed to create BadgerFactStorage: %v", err)
	}
	defer storage.Close()

	// Create and store a collection
	collection := &spookyfactstypes.FactCollection{
		Server:    "test-server",
		Timestamp: time.Now(),
		Facts: map[string]*spookyfactstypes.Fact{
			"os.name": {
				Key:   "os.name",
				Value: "Ubuntu",
			},
		},
	}

	err = storage.SetFactCollection("test-machine", collection)
	if err != nil {
		t.Fatalf("Failed to set fact collection: %v", err)
	}

	// Verify it exists
	_, err = storage.GetFactCollection("test-machine")
	if err != nil {
		t.Fatalf("Failed to get fact collection before deletion: %v", err)
	}

	// Delete the collection
	err = storage.DeleteFactCollection("test-machine")
	if err != nil {
		t.Fatalf("Failed to delete fact collection: %v", err)
	}

	// Verify it's gone
	_, err = storage.GetFactCollection("test-machine")
	if err == nil {
		t.Fatal("Expected error when getting deleted collection")
	}
}

func TestDeleteFactCollections(t *testing.T) {
	storage, err := NewBadgerFactStorage(":memory:")
	if err != nil {
		t.Fatalf("Failed to create BadgerFactStorage: %v", err)
	}
	defer storage.Close()

	// Create test collections
	collection1 := &spookyfactstypes.FactCollection{
		Server:    "server1",
		Timestamp: time.Now(),
		Facts: map[string]*spookyfactstypes.Fact{
			"os.name": {
				Key:   "os.name",
				Value: "Ubuntu",
			},
		},
	}

	collection2 := &spookyfactstypes.FactCollection{
		Server:    "server2",
		Timestamp: time.Now(),
		Facts: map[string]*spookyfactstypes.Fact{
			"os.name": {
				Key:   "os.name",
				Value: "CentOS",
			},
		},
	}

	// Store collections
	err = storage.SetFactCollection("machine1", collection1)
	if err != nil {
		t.Fatalf("Failed to set fact collection 1: %v", err)
	}

	err = storage.SetFactCollection("machine2", collection2)
	if err != nil {
		t.Fatalf("Failed to set fact collection 2: %v", err)
	}

	// Delete collections by OS
	query := &spookyfactstypes.FactQuery{OS: "Ubuntu"}
	deletedCount, err := storage.DeleteFactCollections(query)
	if err != nil {
		t.Fatalf("Failed to delete fact collections: %v", err)
	}

	if deletedCount != 1 {
		t.Errorf("Expected 1 deleted collection, got %d", deletedCount)
	}

	// Verify only Ubuntu collection was deleted
	_, err = storage.GetFactCollection("machine1")
	if err == nil {
		t.Fatal("Expected Ubuntu collection to be deleted")
	}

	_, err = storage.GetFactCollection("machine2")
	if err != nil {
		t.Fatalf("Expected CentOS collection to still exist: %v", err)
	}
}
