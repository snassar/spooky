package facts

import (
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	manager := NewManager(nil)
	if manager == nil {
		t.Fatal("NewManager returned nil")
	}

	if manager.defaultTTL != DefaultTTL {
		t.Errorf("Expected default TTL %v, got %v", DefaultTTL, manager.defaultTTL)
	}
}

func TestLocalFactCollection(t *testing.T) {
	manager := NewManager(nil)

	// Test collecting all facts from local machine
	collection, err := manager.CollectAllFacts("local")
	if err != nil {
		t.Fatalf("Failed to collect local facts: %v", err)
	}

	if collection == nil {
		t.Fatal("Collection is nil")
	}

	if collection.Server != "local" {
		t.Errorf("Expected server 'local', got '%s'", collection.Server)
	}

	// Check that we have some basic facts
	expectedFacts := []string{FactHostname, FactOSArch}
	for _, factKey := range expectedFacts {
		if fact, exists := collection.Facts[factKey]; !exists {
			t.Errorf("Expected fact '%s' not found", factKey)
		} else if fact == nil {
			t.Errorf("Fact '%s' is nil", factKey)
		} else if fact.Source != string(SourceLocal) {
			t.Errorf("Expected fact source '%s', got '%s'", SourceLocal, fact.Source)
		}
	}
}

func TestSpecificFactCollection(t *testing.T) {
	manager := NewManager(nil)

	// Test collecting specific facts
	specificFacts := []string{FactHostname, FactOSArch}
	collection, err := manager.CollectSpecificFacts("local", specificFacts)
	if err != nil {
		t.Fatalf("Failed to collect specific facts: %v", err)
	}

	if collection == nil {
		t.Fatal("Collection is nil")
	}

	// Check that we have at least the requested facts
	for _, expected := range specificFacts {
		if _, exists := collection.Facts[expected]; !exists {
			t.Errorf("Expected fact '%s' not found in collection", expected)
		}
	}
}

func TestSingleFactRetrieval(t *testing.T) {
	manager := NewManager(nil)

	// Test retrieving a single fact
	fact, err := manager.GetFact("local", FactHostname)
	if err != nil {
		t.Fatalf("Failed to get fact: %v", err)
	}

	if fact == nil {
		t.Fatal("Fact is nil")
	}

	if fact.Key != FactHostname {
		t.Errorf("Expected fact key '%s', got '%s'", FactHostname, fact.Key)
	}

	if fact.Source != string(SourceLocal) {
		t.Errorf("Expected fact source '%s', got '%s'", SourceLocal, fact.Source)
	}
}

func TestFactCaching(t *testing.T) {
	manager := NewManager(nil)

	// Collect facts twice
	collection1, err := manager.CollectAllFacts("local")
	if err != nil {
		t.Fatalf("Failed to collect facts first time: %v", err)
	}

	collection2, err := manager.CollectAllFacts("local")
	if err != nil {
		t.Fatalf("Failed to collect facts second time: %v", err)
	}

	// Both should be the same (cached)
	if collection1 != collection2 {
		t.Error("Cached collection should be the same instance")
	}
}

func TestCacheExpiration(t *testing.T) {
	manager := NewManager(nil)

	// Manually create expired facts in cache
	expiredFact := &Fact{
		Key:       "test-expired",
		Value:     "test-value",
		Source:    string(SourceLocal),
		Server:    "local",
		Timestamp: time.Now().Add(-2 * time.Hour), // 2 hours ago
		TTL:       1 * time.Hour,                  // 1 hour TTL
	}

	manager.cacheFact("local", expiredFact)

	// Verify fact is in cache (getCachedFacts returns nil for expired facts)
	if cached := manager.getCachedFacts("local"); cached != nil {
		t.Fatal("getCachedFacts should return nil for expired facts")
	}

	// Clear expired cache
	manager.ClearExpiredCache()

	// Check that the cache is empty after clearing expired facts
	manager.cacheMutex.RLock()
	collection, exists := manager.cache["local"]
	manager.cacheMutex.RUnlock()

	if exists && len(collection.Facts) > 0 {
		t.Error("Cache should be empty after clearing expired facts")
	}
}

func TestCacheClearing(t *testing.T) {
	manager := NewManager(nil)

	// Collect facts
	_, err := manager.CollectAllFacts("local")
	if err != nil {
		t.Fatalf("Failed to collect facts: %v", err)
	}

	// Verify cache has data
	if cached := manager.getCachedFacts("local"); cached == nil {
		t.Error("Cache should have data")
	}

	// Clear cache
	manager.ClearCache()

	// Verify cache is empty
	if cached := manager.getCachedFacts("local"); cached != nil {
		t.Error("Cache should be empty after clearing")
	}
}

func TestFactExpiration(t *testing.T) {
	manager := NewManager(nil)

	// Create a fact with short TTL
	fact := &Fact{
		Key:       "test",
		Value:     "test-value",
		Source:    string(SourceLocal),
		Server:    "local",
		Timestamp: time.Now().Add(-2 * time.Hour), // 2 hours ago
		TTL:       1 * time.Hour,                  // 1 hour TTL
	}

	// Should be expired
	if !manager.isExpired(fact) {
		t.Error("Fact should be expired")
	}

	// Create a fact with no TTL
	factNoTTL := &Fact{
		Key:       "test-no-ttl",
		Value:     "test-value",
		Source:    string(SourceLocal),
		Server:    "local",
		Timestamp: time.Now().Add(-2 * time.Hour),
		TTL:       0, // No expiration
	}

	// Should not be expired
	if manager.isExpired(factNoTTL) {
		t.Error("Fact with no TTL should not be expired")
	}
}
