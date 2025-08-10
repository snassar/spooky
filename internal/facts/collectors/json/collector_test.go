package json

import (
	spookytypes "spooky/internal/types"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	spookylogging "spooky/internal/logging"
	spookytypes "spooky/internal/types"
)

func TestNewCollector(t *testing.T) {
	logger := spookylogging.GetLogger()
	collector := NewCollector("/test/path", logger)

	if collector.sourcePath != "/test/path" {
		t.Errorf("Expected sourcePath to be '/test/path', got '%s'", collector.sourcePath)
	}

	if collector.GetSource() != spookytypes.SourceJSON {
		t.Errorf("Expected source to be SourceJSON, got '%s'", collector.GetSource())
	}

	if collector.GetMergePolicy() != spookytypes.MergePolicyReplace {
		t.Errorf("Expected merge policy to be MergePolicyReplace, got '%s'", collector.GetMergePolicy())
	}
}

func TestCollector_Validate(t *testing.T) {
	logger := spookylogging.GetLogger()

	tests := []struct {
		name        string
		sourcePath  string
		expectError bool
	}{
		{
			name:        "empty source path",
			sourcePath:  "",
			expectError: true,
		},
		{
			name:        "non-existent path",
			sourcePath:  "/non/existent/path",
			expectError: true,
		},
		{
			name:        "valid path",
			sourcePath:  "/tmp",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := NewCollector(tt.sourcePath, logger)
			err := collector.Validate()

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestCollector_CollectFromFile(t *testing.T) {
	logger := spookylogging.GetLogger()
	tempDir := t.TempDir()

	// Create test JSON file
	testData := map[string]interface{}{
		"name":    "test-server",
		"version": "1.0.0",
		"config": map[string]interface{}{
			"port":    8080,
			"enabled": true,
		},
		"tags": []string{"web", "api"},
	}

	jsonFile := filepath.Join(tempDir, "test.json")
	jsonData, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}
	if err := os.WriteFile(jsonFile, jsonData, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	collector := NewCollector(jsonFile, logger)
	collection, err := collector.Collect("test-server")

	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	if collection.Server != "test-server" {
		t.Errorf("Expected server to be 'test-server', got '%s'", collection.Server)
	}

	// Check that facts were created
	expectedFacts := []string{"name", "version", "config.port", "config.enabled", "tags[0]", "tags[1]"}
	for _, factKey := range expectedFacts {
		if fact, exists := collection.Facts[factKey]; !exists {
			t.Errorf("Expected fact '%s' to exist", factKey)
		} else {
			if fact.Source != string(spookytypes.SourceJSON) {
				t.Errorf("Expected fact source to be 'json', got '%s'", fact.Source)
			}
			if fact.Server != "local" {
				t.Errorf("Expected fact server to be 'local', got '%s'", fact.Server)
			}
		}
	}

	// Check specific fact values
	if fact, exists := collection.Facts["name"]; exists {
		if fact.Value != "test-server" {
			t.Errorf("Expected name fact value to be 'test-server', got '%v'", fact.Value)
		}
	}

	if fact, exists := collection.Facts["config.port"]; exists {
		if fact.Value != float64(8080) { // JSON numbers are unmarshaled as float64
			t.Errorf("Expected config.port fact value to be 8080, got '%v'", fact.Value)
		}
	}
}

func TestCollector_CollectFromDirectory(t *testing.T) {
	logger := spookylogging.GetLogger()
	tempDir := t.TempDir()

	// Create multiple JSON files
	files := []struct {
		name string
		data map[string]interface{}
	}{
		{
			name: "file1.json",
			data: map[string]interface{}{
				"name": "server1",
				"port": 8080,
			},
		},
		{
			name: "file2.json",
			data: map[string]interface{}{
				"name": "server2",
				"port": 8081,
			},
		},
		{
			name: "ignore.txt", // Should be ignored
			data: map[string]interface{}{
				"name": "ignored",
			},
		},
	}

	for _, file := range files {
		jsonData, err := json.Marshal(file.data)
		if err != nil {
			t.Fatalf("Failed to marshal test data: %v", err)
		}
		filePath := filepath.Join(tempDir, file.name)
		if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
	}

	collector := NewCollector(tempDir, logger)
	collection, err := collector.Collect("test-server")

	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	// Should have facts from both JSON files
	expectedFacts := []string{"name", "port"}
	for _, factKey := range expectedFacts {
		if fact, exists := collection.Facts[factKey]; !exists {
			t.Errorf("Expected fact '%s' to exist", factKey)
		} else if fact.Source != string(spookytypes.SourceJSON) {
			t.Errorf("Expected fact source to be 'json', got '%s'", fact.Source)
		}
	}

	// Should not have facts from the .txt file
	if _, exists := collection.Facts["ignored"]; exists {
		t.Errorf("Expected fact 'ignored' to not exist")
	}
}

func TestCollector_CollectSpecific(t *testing.T) {
	logger := spookylogging.GetLogger()
	tempDir := t.TempDir()

	testData := map[string]interface{}{
		"name":    "test-server",
		"version": "1.0.0",
		"port":    8080,
		"enabled": true,
	}

	jsonFile := filepath.Join(tempDir, "test.json")
	jsonData, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}
	if err := os.WriteFile(jsonFile, jsonData, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	collector := NewCollector(jsonFile, logger)
	collection, err := collector.CollectSpecific("test-server", []string{"name", "port"})

	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	// Should only have the requested facts
	expectedFacts := []string{"name", "port"}
	unexpectedFacts := []string{"version", "enabled"}

	for _, factKey := range expectedFacts {
		if _, exists := collection.Facts[factKey]; !exists {
			t.Errorf("Expected fact '%s' to exist", factKey)
		}
	}

	for _, factKey := range unexpectedFacts {
		if _, exists := collection.Facts[factKey]; exists {
			t.Errorf("Expected fact '%s' to not exist", factKey)
		}
	}
}

func TestCollector_GetFact(t *testing.T) {
	logger := spookylogging.GetLogger()
	tempDir := t.TempDir()

	testData := map[string]interface{}{
		"name": "test-server",
		"port": 8080,
	}

	jsonFile := filepath.Join(tempDir, "test.json")
	jsonData, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}
	if err := os.WriteFile(jsonFile, jsonData, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	collector := NewCollector(jsonFile, logger)

	// Test getting existing fact
	fact, err := collector.GetFact("test-server", "name")
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	if fact.Value != "test-server" {
		t.Errorf("Expected fact value to be 'test-server', got '%v'", fact.Value)
	}

	// Test getting non-existent fact
	_, err = collector.GetFact("test-server", "non-existent")
	if err == nil {
		t.Errorf("Expected error but got none")
	}
}

func TestCollector_ValidateJSONFile(t *testing.T) {
	logger := spookylogging.GetLogger()
	tempDir := t.TempDir()

	collector := NewCollector(tempDir, logger)

	// Test non-JSON file
	txtFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(txtFile, []byte("not json"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	err := collector.validateJSONFile(txtFile)
	if err == nil {
		t.Errorf("Expected error for non-JSON file but got none")
	}

	// Test valid JSON file
	jsonFile := filepath.Join(tempDir, "test.json")
	if err := os.WriteFile(jsonFile, []byte(`{"name": "test"}`), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	err = collector.validateJSONFile(jsonFile)
	if err != nil {
		t.Errorf("Expected no error for valid JSON file but got: %v", err)
	}
}

func.*runs(t *testing.T) {
	logger := spookylogging.GetLogger()
	collector := NewCollector("/test/path", logger)

	facts := make(map[string]*spookytypes.Fact)

	// Test nested object
	testData := map[string]interface{}{
		"server": map[string]interface{}{
			"name": "test-server",
			"config": map[string]interface{}{
				"port": 8080,
			},
		},
		"tags": []interface{}{"web", "api"},
	}

	collector.convertValueToFacts("", testData, facts, "test.json")

	expectedFacts := []string{
		"server.name",
		"server.config.port",
		"tags[0]",
		"tags[1]",
	}

	for _, factKey := range expectedFacts {
		if fact, exists := facts[factKey]; !exists {
			t.Errorf("Expected fact '%s' to exist", factKey)
		} else {
			if fact.Source != string(spookytypes.SourceJSON) {
				t.Errorf("Expected fact source to be 'json', got '%s'", fact.Source)
			}
			if fact.Server != "local" {
				t.Errorf("Expected fact server to be 'local', got '%s'", fact.Server)
			}
			if fact.TTL != spookytypes.DefaultTTL {
				t.Errorf("Expected fact TTL to be DefaultTTL, got '%v'", fact.TTL)
			}
		}
	}

	// Check metadata
	if fact, exists := facts["server.name"]; exists {
		if fact.Metadata["source_file"] != "test.json" {
			t.Errorf("Expected source_file metadata to be 'test.json', got '%v'", fact.Metadata["source_file"])
		}
		if fact.Metadata["source_type"] != "json" {
			t.Errorf("Expected source_type metadata to be 'json', got '%v'", fact.Metadata["source_type"])
		}
	}
}

func TestCollector_FileSizeLimit(t *testing.T) {
	logger := spookylogging.GetLogger()
	tempDir := t.TempDir()

	collector := NewCollector(tempDir, logger)

	// Create a large JSON file (over 10MB)
	largeData := make(map[string]interface{})
	for i := 0; i < 100000; i++ {
		largeData[fmt.Sprintf("key%d", i)] = strings.Repeat("value", 100)
	}

	jsonFile := filepath.Join(tempDir, "large.json")
	jsonData, err := json.Marshal(largeData)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}
	if err := os.WriteFile(jsonFile, jsonData, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Check file size
	info, _ := os.Stat(jsonFile)
	if info.Size() < 10*1024*1024 {
		t.Skip("File is not large enough to test size limit")
	}

	err = collector.validateJSONFile(jsonFile)
	if err == nil {
		t.Errorf("Expected error for large file but got none")
	}
}
