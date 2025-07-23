package facts

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPCollector_FlatJSON(t *testing.T) {
	// Create test server
	testData := map[string]interface{}{
		"hostname":     "test-server",
		"os.name":      "Linux",
		"os.version":   "Ubuntu 22.04",
		"cpu.cores":    4,
		"memory.total": 8589934592,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(testData); err != nil {
			t.Errorf("Failed to encode test data: %v", err)
		}
	}))
	defer server.Close()

	// Create HTTP collector
	collector := NewHTTPCollector(server.URL, nil, 30*time.Second, MergePolicyReplace)

	// Collect facts
	collection, err := collector.Collect("test-server")
	if err != nil {
		t.Fatalf("Failed to collect facts: %v", err)
	}

	// Verify facts
	if len(collection.Facts) != 5 {
		t.Errorf("Expected 5 facts, got %d", len(collection.Facts))
	}

	// Check specific facts
	if fact, exists := collection.Facts["hostname"]; !exists {
		t.Error("hostname fact not found")
	} else if fact.Value != "test-server" {
		t.Errorf("Expected hostname 'test-server', got %v", fact.Value)
	}

	if fact, exists := collection.Facts["cpu.cores"]; !exists {
		t.Error("cpu.cores fact not found")
	} else if fact.Value != float64(4) {
		t.Errorf("Expected cpu.cores 4, got %v", fact.Value)
	}
}

func TestHTTPCollector_ArrayJSON(t *testing.T) {
	// Create test server
	testData := []map[string]interface{}{
		{"key": "hostname", "value": "test-server"},
		{"key": "os.name", "value": "Linux"},
		{"key": "os.version", "value": "Ubuntu 22.04"},
		{"key": "cpu.cores", "value": 4},
		{"key": "memory.total", "value": 8589934592},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(testData); err != nil {
			t.Errorf("Failed to encode test data: %v", err)
		}
	}))
	defer server.Close()

	// Create HTTP collector
	collector := NewHTTPCollector(server.URL, nil, 30*time.Second, MergePolicyReplace)

	// Collect facts
	collection, err := collector.Collect("test-server")
	if err != nil {
		t.Fatalf("Failed to collect facts: %v", err)
	}

	// Verify facts
	if len(collection.Facts) != 5 {
		t.Errorf("Expected 5 facts, got %d", len(collection.Facts))
	}

	// Check specific facts
	if fact, exists := collection.Facts["hostname"]; !exists {
		t.Error("hostname fact not found")
	} else if fact.Value != "test-server" {
		t.Errorf("Expected hostname 'test-server', got %v", fact.Value)
	}
}

func TestHTTPCollector_CollectSpecific(t *testing.T) {
	// Create test server
	testData := map[string]interface{}{
		"hostname":     "test-server",
		"os.name":      "Linux",
		"os.version":   "Ubuntu 22.04",
		"cpu.cores":    4,
		"memory.total": 8589934592,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(testData); err != nil {
			t.Errorf("Failed to encode test data: %v", err)
		}
	}))
	defer server.Close()

	// Create HTTP collector
	collector := NewHTTPCollector(server.URL, nil, 30*time.Second, MergePolicyReplace)

	// Collect specific facts
	collection, err := collector.CollectSpecific("test-server", []string{"hostname", "cpu.cores"})
	if err != nil {
		t.Fatalf("Failed to collect specific facts: %v", err)
	}

	// Verify only requested facts are present
	if len(collection.Facts) != 2 {
		t.Errorf("Expected 2 facts, got %d", len(collection.Facts))
	}

	if _, exists := collection.Facts["hostname"]; !exists {
		t.Error("hostname fact not found")
	}

	if _, exists := collection.Facts["cpu.cores"]; !exists {
		t.Error("cpu.cores fact not found")
	}

	if _, exists := collection.Facts["os.name"]; exists {
		t.Error("os.name fact should not be present")
	}
}

func TestHTTPCollector_GetFact(t *testing.T) {
	// Create test server
	testData := map[string]interface{}{
		"hostname": "test-server",
		"os.name":  "Linux",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(testData); err != nil {
			t.Errorf("Failed to encode test data: %v", err)
		}
	}))
	defer server.Close()

	// Create HTTP collector
	collector := NewHTTPCollector(server.URL, nil, 30*time.Second, MergePolicyReplace)

	// Get specific fact
	fact, err := collector.GetFact("test-server", "hostname")
	if err != nil {
		t.Fatalf("Failed to get fact: %v", err)
	}

	if fact.Value != "test-server" {
		t.Errorf("Expected hostname 'test-server', got %v", fact.Value)
	}

	// Try to get non-existent fact
	_, err = collector.GetFact("test-server", "non-existent")
	if err == nil {
		t.Error("Expected error for non-existent fact")
	}
}
