package facts

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestJSONCollector_FlatJSON(t *testing.T) {
	// Create test JSON file
	testData := map[string]interface{}{
		"hostname":     "test-server",
		"os.name":      "Linux",
		"os.version":   "Ubuntu 22.04",
		"cpu.cores":    4,
		"memory.total": 8589934592,
	}

	tempFile, err := os.CreateTemp("", "test-facts-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data
	encoder := json.NewEncoder(tempFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tempFile.Close()

	// Create collector
	collector := NewJSONCollector(tempFile.Name(), MergePolicyReplace)

	// Collect facts
	collection, err := collector.Collect("test-server")
	if err != nil {
		t.Fatalf("Failed to collect facts: %v", err)
	}

	// Verify results
	if collection.Server != "test-server" {
		t.Errorf("Expected server 'test-server', got '%s'", collection.Server)
	}

	if len(collection.Facts) != 5 {
		t.Errorf("Expected 5 facts, got %d", len(collection.Facts))
	}

	// Check specific facts
	if fact, exists := collection.Facts["hostname"]; exists {
		if fact.Value != "test-server" {
			t.Errorf("Expected hostname 'test-server', got '%v'", fact.Value)
		}
		if fact.Source != string(SourceCustom) {
			t.Errorf("Expected source 'custom', got '%s'", fact.Source)
		}
	} else {
		t.Error("hostname fact not found")
	}

	if fact, exists := collection.Facts["cpu.cores"]; exists {
		if fact.Value != float64(4) { // JSON numbers are float64
			t.Errorf("Expected cpu.cores 4, got %v", fact.Value)
		}
	} else {
		t.Error("cpu.cores fact not found")
	}
}

func TestJSONCollector_ArrayJSON(t *testing.T) {
	// Create test JSON file with array format
	testData := []map[string]interface{}{
		{
			"key":    "hostname",
			"value":  "test-server",
			"source": "custom",
			"ttl":    3600,
		},
		{
			"key":    "os.name",
			"value":  "Linux",
			"source": "custom",
		},
		{
			"key":   "cpu.cores",
			"value": 4,
			"metadata": map[string]interface{}{
				"unit": "cores",
			},
		},
	}

	tempFile, err := os.CreateTemp("", "test-facts-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data
	encoder := json.NewEncoder(tempFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tempFile.Close()

	// Create collector
	collector := NewJSONCollector(tempFile.Name(), MergePolicyReplace)

	// Collect facts
	collection, err := collector.Collect("test-server")
	if err != nil {
		t.Fatalf("Failed to collect facts: %v", err)
	}

	// Verify results
	if len(collection.Facts) != 3 {
		t.Errorf("Expected 3 facts, got %d", len(collection.Facts))
	}

	// Check specific facts
	if fact, exists := collection.Facts["hostname"]; exists {
		if fact.Value != "test-server" {
			t.Errorf("Expected hostname 'test-server', got '%v'", fact.Value)
		}
		if fact.Source != "custom" {
			t.Errorf("Expected source 'custom', got '%s'", fact.Source)
		}
		if fact.TTL != 3600*time.Second {
			t.Errorf("Expected TTL 3600s, got %v", fact.TTL)
		}
	} else {
		t.Error("hostname fact not found")
	}

	if fact, exists := collection.Facts["cpu.cores"]; exists {
		if fact.Value != float64(4) {
			t.Errorf("Expected cpu.cores 4, got %v", fact.Value)
		}
		if metadata, ok := fact.Metadata["unit"]; !ok || metadata != "cores" {
			t.Errorf("Expected metadata unit=cores, got %v", fact.Metadata)
		}
	} else {
		t.Error("cpu.cores fact not found")
	}
}

func TestJSONCollector_CollectSpecific(t *testing.T) {
	// Create test JSON file
	testData := map[string]interface{}{
		"hostname":     "test-server",
		"os.name":      "Linux",
		"cpu.cores":    4,
		"memory.total": 8589934592,
	}

	tempFile, err := os.CreateTemp("", "test-facts-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data
	encoder := json.NewEncoder(tempFile)
	if err := encoder.Encode(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tempFile.Close()

	// Create collector
	collector := NewJSONCollector(tempFile.Name(), MergePolicyReplace)

	// Collect specific facts
	collection, err := collector.CollectSpecific("test-server", []string{"hostname", "os.name"})
	if err != nil {
		t.Fatalf("Failed to collect specific facts: %v", err)
	}

	// Verify results
	if len(collection.Facts) != 2 {
		t.Errorf("Expected 2 facts, got %d", len(collection.Facts))
	}

	if _, exists := collection.Facts["hostname"]; !exists {
		t.Error("hostname fact not found")
	}

	if _, exists := collection.Facts["os.name"]; !exists {
		t.Error("os.name fact not found")
	}

	if _, exists := collection.Facts["cpu.cores"]; exists {
		t.Error("cpu.cores fact should not be included")
	}
}

func TestJSONCollector_GetFact(t *testing.T) {
	// Create test JSON file
	testData := map[string]interface{}{
		"hostname": "test-server",
		"os.name":  "Linux",
	}

	tempFile, err := os.CreateTemp("", "test-facts-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data
	encoder := json.NewEncoder(tempFile)
	if err := encoder.Encode(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tempFile.Close()

	// Create collector
	collector := NewJSONCollector(tempFile.Name(), MergePolicyReplace)

	// Get specific fact
	fact, err := collector.GetFact("test-server", "hostname")
	if err != nil {
		t.Fatalf("Failed to get fact: %v", err)
	}

	if fact.Value != "test-server" {
		t.Errorf("Expected hostname 'test-server', got '%v'", fact.Value)
	}

	// Test non-existent fact
	_, err = collector.GetFact("test-server", "non-existent")
	if err == nil {
		t.Error("Expected error for non-existent fact")
	}
}
