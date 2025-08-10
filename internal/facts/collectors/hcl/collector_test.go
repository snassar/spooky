package hcl

import (
	spookytypes "spooky/internal/types"
	"os"
	"path/filepath"
	"strings"
	"testing"

	spookytypes "spooky/internal/types"
)

func TestNewCollector(t *testing.T) {
	hclFiles := []string{"/test/path/config.hcl"}
	collector := NewCollector(hclFiles)

	if len(collector.hclFiles) != 1 || collector.hclFiles[0] != "/test/path/config.hcl" {
		t.Errorf("Expected hclFiles to contain '/test/path/config.hcl', got %v", collector.hclFiles)
	}

	if collector.GetSource() != spookytypes.SourceHCL {
		t.Errorf("Expected source to be SourceHCL, got '%s'", collector.GetSource())
	}

	if collector.GetMergePolicy() != spookytypes.MergePolicyMerge {
		t.Errorf("Expected merge policy to be MergePolicyMerge, got '%s'", collector.GetMergePolicy())
	}
}

func TestNewCollectorWithParser(t *testing.T) {
	hclFiles := []string{"/test/path/config.hcl"}
	mockParser := &MockParser{}
	collector := NewCollectorWithParser(hclFiles, mockParser)

	if collector.parser != mockParser {
		t.Errorf("Expected parser to be mockParser")
	}
}

func TestCollector_Validate(t *testing.T) {
	tests := []struct {
		name        string
		hclFiles    []string
		expectError bool
	}{
		{
			name:        "empty hcl files",
			hclFiles:    []string{},
			expectError: true,
		},
		{
			name:        "non-existent file",
			hclFiles:    []string{"/non/existent/file.hcl"},
			expectError: true,
		},
		{
			name:        "valid file",
			hclFiles:    []string{"/tmp"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := NewCollector(tt.hclFiles)
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
	tempDir := t.TempDir()

	// Create test HCL file
	testHCL := `
name = "test-server"
version = "1.0.0"
enabled = true
port = 8080
tags = ["web", "api", "production"]
`

	hclFile := filepath.Join(tempDir, "test.hcl")
	if err := os.WriteFile(hclFile, []byte(testHCL), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	collector := NewCollector([]string{hclFile})
	collection, err := collector.Collect("test-server")

	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	if collection.Server != "test-server" {
		t.Errorf("Expected server to be 'test-server', got '%s'", collection.Server)
	}

	// Check for expected facts
	expectedFacts := []string{"name", "version", "enabled", "port", "tags"}
	for _, factKey := range expectedFacts {
		if _, exists := collection.Facts[factKey]; !exists {
			t.Errorf("Expected fact '%s' to exist", factKey)
		}
	}

	// Check specific fact values
	if fact, exists := collection.Facts["name"]; exists {
		if fact.Value != "test-server" {
			t.Errorf("Expected name to be 'test-server', got '%v'", fact.Value)
		}
		if fact.Source != string(spookytypes.SourceHCL) {
			t.Errorf("Expected source to be SourceHCL, got '%s'", fact.Source)
		}
	}
}

func TestCollector_CollectFromDirectory(t *testing.T) {
	tempDir := t.TempDir()

	// Create multiple test HCL files
	testFiles := map[string]string{
		"config1.hcl": `name = "server1"`,
		"config2.hcl": `version = "1.0.0"`,
		"ignore.txt":  `not hcl content`,
	}

	for filename, content := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test file %s: %v", filename, err)
		}
	}

	collector := NewCollector([]string{tempDir})
	collection, err := collector.Collect("test-server")

	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	// Should have facts from both HCL files
	if len(collection.Facts) < 2 {
		t.Errorf("Expected at least 2 facts, got %d", len(collection.Facts))
	}

	// Check that non-HCL files were ignored
	for key := range collection.Facts {
		if strings.Contains(key, "ignore") {
			t.Errorf("Expected ignore.txt to be ignored, but found fact: %s", key)
		}
	}
}

func TestCollector_CollectSpecific(t *testing.T) {
	tempDir := t.TempDir()

	testHCL := `
name = "test-server"
version = "1.0.0"
enabled = true
port = 8080
`

	hclFile := filepath.Join(tempDir, "test.hcl")
	if err := os.WriteFile(hclFile, []byte(testHCL), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	collector := NewCollector([]string{hclFile})
	collection, err := collector.CollectSpecific("test-server", []string{"name", "version"})

	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	// Should only have requested facts
	if len(collection.Facts) != 2 {
		t.Errorf("Expected 2 facts, got %d", len(collection.Facts))
	}

	// Check specific facts exist
	if _, exists := collection.Facts["name"]; !exists {
		t.Errorf("Expected fact 'name' to exist")
	}
	if _, exists := collection.Facts["version"]; !exists {
		t.Errorf("Expected fact 'version' to exist")
	}
	if _, exists := collection.Facts["enabled"]; exists {
		t.Errorf("Expected fact 'enabled' to not exist")
	}
}

func TestCollector_GetFact(t *testing.T) {
	tempDir := t.TempDir()

	testHCL := `
name = "test-server"
version = "1.0.0"
`

	hclFile := filepath.Join(tempDir, "test.hcl")
	if err := os.WriteFile(hclFile, []byte(testHCL), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	collector := NewCollector([]string{hclFile})

	// Test existing fact
	fact, err := collector.GetFact("test-server", "name")
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}
	if fact.Value != "test-server" {
		t.Errorf("Expected name to be 'test-server', got '%v'", fact.Value)
	}

	// Test non-existing fact
	_, err = collector.GetFact("test-server", "non-existent")
	if err == nil {
		t.Errorf("Expected error for non-existent fact")
	}
}

func TestCollector_ValidateHCLFile(t *testing.T) {
	tempDir := t.TempDir()

	// Valid HCL file
	validHCL := `name = "test"`
	validFile := filepath.Join(tempDir, "valid.hcl")
	if err := os.WriteFile(validFile, []byte(validHCL), 0644); err != nil {
		t.Fatalf("Failed to write valid test file: %v", err)
	}

	// Invalid HCL file
	invalidHCL := `name = "test" = invalid`
	invalidFile := filepath.Join(tempDir, "invalid.hcl")
	if err := os.WriteFile(invalidFile, []byte(invalidHCL), 0644); err != nil {
		t.Fatalf("Failed to write invalid test file: %v", err)
	}

	collector := NewCollector([]string{})

	// Test valid file
	if err := collector.ValidateHCLFile(validFile); err != nil {
		t.Errorf("Expected no error for valid HCL file, got: %v", err)
	}

	// Test invalid file
	if err := collector.ValidateHCLFile(invalidFile); err == nil {
		t.Errorf("Expected error for invalid HCL file")
	}
}

func TestCollector_FindHCLFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Create test files
	testFiles := map[string]string{
		"config1.hcl":        `name = "server1"`,
		"config2.hcl":        `name = "server2"`,
		"ignore.txt":         `not hcl content`,
		"subdir/config3.hcl": `name = "server3"`,
	}

	// Create subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0o755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	for filename, content := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test file %s: %v", filename, err)
		}
	}

	collector := NewCollector([]string{})
	hclFiles, err := collector.FindHCLFiles(tempDir)

	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	// Should find 3 HCL files (including subdirectory)
	if len(hclFiles) != 3 {
		t.Errorf("Expected 3 HCL files, got %d", len(hclFiles))
	}

	// Check that only .hcl files are found
	for _, file := range hclFiles {
		if !strings.HasSuffix(file, ".hcl") {
			t.Errorf("Expected only .hcl files, got: %s", file)
		}
	}
}

func.*runsToHCL(t *testing.T) {
	collector := NewCollector([]string{})

	facts := map[string]interface{}{
		"name":    "test-server",
		"version": "1.0.0",
		"enabled": true,
		"port":    8080,
		"tags":    []string{"web", "api"},
		"config": map[string]interface{}{
			"host": "localhost",
			"ssl":  true,
		},
	}

	// Create temporary file for output
	tempFile := filepath.Join(t.TempDir(), "output.hcl")
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create output file: %v", err)
	}
	defer file.Close()

	err = collector.ExportFactsToHCL(facts, file)
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	// Read back the exported content
	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read exported file: %v", err)
	}

	exportedContent := string(content)

	// Check that exported content contains expected values
	expectedValues := []string{"test-server", "1.0.0", "true", "8080", "web", "api", "localhost"}
	for _, expected := range expectedValues {
		if !strings.Contains(exportedContent, expected) {
			t.Errorf("Expected exported content to contain '%s'", expected)
		}
	}
}

func TestDefaultParser_ParseContent(t *testing.T) {
	parser := &DefaultParser{}

	// Test basic HCL content
	content := `
machine_id = "1234567890abcdef1234567890abcdef"
collected_at = "2024-01-01T00:00:00Z"
version = "1.0.0"
`

	facts, err := parser.ParseContent([]byte(content))
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	expectedFacts := map[string]interface{}{
		"machine_id":   "1234567890abcdef1234567890abcdef",
		"collected_at": "2024-01-01T00:00:00Z",
		"version":      "1.0.0",
	}

	for key, expectedValue := range expectedFacts {
		if value, exists := facts[key]; !exists {
			t.Errorf("Expected fact '%s' to exist", key)
		} else if value != expectedValue {
			t.Errorf("Expected fact '%s' to be '%v', got '%v'", key, expectedValue, value)
		}
	}
}

func TestDefaultParser_ParseContentWithBlocks(t *testing.T) {
	parser := &DefaultParser{}

	// Test HCL content with blocks
	content := `
machine_id = "1234567890abcdef1234567890abcdef"
collected_at = "2024-01-01T00:00:00Z"

fact "system_info" {
  os = "linux"
  version = "20.04"
}

resource "aws_instance" "web" {
  instance_type = "t3.micro"
  ami = "ami-123456"
}
`

	facts, err := parser.ParseContent([]byte(content))
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	// Check that basic attributes are parsed
	if value, exists := facts["machine_id"]; !exists {
		t.Error("Expected fact 'machine_id' to exist")
	} else if value != "1234567890abcdef1234567890abcdef" {
		t.Errorf("Expected machine_id to be '1234567890abcdef1234567890abcdef', got '%v'", value)
	}

	// Check that blocks are parsed
	if factBlock, exists := facts["fact_system_info"]; !exists {
		t.Error("Expected fact block 'fact_system_info' to exist")
	} else {
		factMap, ok := factBlock.(map[string]interface{})
		if !ok {
			t.Error("Expected fact block to be a map")
		} else {
			if factMap["type"] != "fact" {
				t.Errorf("Expected block type to be 'fact', got '%v'", factMap["type"])
			}
			if factMap["os"] != "linux" {
				t.Errorf("Expected os to be 'linux', got '%v'", factMap["os"])
			}
			if factMap["version"] != "20.04" {
				t.Errorf("Expected version to be '20.04', got '%v'", factMap["version"])
			}
		}
	}

	if resourceBlock, exists := facts["resource_aws_instance_web"]; !exists {
		t.Error("Expected resource block 'resource_aws_instance_web' to exist")
	} else {
		resourceMap, ok := resourceBlock.(map[string]interface{})
		if !ok {
			t.Error("Expected resource block to be a map")
		} else {
			if resourceMap["type"] != "resource" {
				t.Errorf("Expected block type to be 'resource', got '%v'", resourceMap["type"])
			}
			if resourceMap["instance_type"] != "t3.micro" {
				t.Errorf("Expected instance_type to be 't3.micro', got '%v'", resourceMap["instance_type"])
			}
			if resourceMap["ami"] != "ami-123456" {
				t.Errorf("Expected ami to be 'ami-123456', got '%v'", resourceMap["ami"])
			}
		}
	}
}

func TestDefaultParser_ParseContentWithNestedBlocks(t *testing.T) {
	parser := &DefaultParser{}

	// Test HCL content with nested blocks
	content := `
machine_id = "1234567890abcdef1234567890abcdef"

resource "aws_instance" "web" {
  instance_type = "t3.micro"
  
  root_block_device {
    volume_size = 20
    volume_type = "gp3"
  }
}
`

	facts, err := parser.ParseContent([]byte(content))
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	// Check that nested blocks are parsed
	if resourceBlock, exists := facts["resource_aws_instance_web"]; !exists {
		t.Error("Expected resource block 'resource_aws_instance_web' to exist")
	} else {
		resourceMap, ok := resourceBlock.(map[string]interface{})
		if !ok {
			t.Error("Expected resource block to be a map")
		} else {
			// Check for nested blocks
			if nestedBlocks, exists := resourceMap["blocks"]; !exists {
				t.Error("Expected nested blocks to exist")
			} else {
				nestedMap, ok := nestedBlocks.(map[string]interface{})
				if !ok {
					t.Error("Expected nested blocks to be a map")
				} else {
					if rootBlock, exists := nestedMap["root_block_device_"]; !exists {
						t.Error("Expected nested root_block_device block to exist")
					} else {
						rootMap, ok := rootBlock.(map[string]interface{})
						if !ok {
							t.Error("Expected root_block_device block to be a map")
						} else {
							if rootMap["volume_size"] != float64(20) {
								t.Errorf("Expected volume_size to be 20, got '%v'", rootMap["volume_size"])
							}
							if rootMap["volume_type"] != "gp3" {
								t.Errorf("Expected volume_type to be 'gp3', got '%v'", rootMap["volume_type"])
							}
						}
					}
				}
			}
		}
	}
}

func TestDefaultParser_ParseFileWithDebug(t *testing.T) {
	parser := &DefaultParser{}

	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.hcl")
	testContent := `
machine_id = "1234567890abcdef1234567890abcdef"
collected_at = "2024-01-01T00:00:00Z"

fact "test" {
  value = "test-value"
}
`

	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Test ParseFileWithDebug
	facts, debugInfo, err := parser.ParseFileWithDebug(testFile)
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	// Check debug info
	if debugInfo.FilePath != testFile {
		t.Errorf("Expected FilePath to be '%s', got '%s'", testFile, debugInfo.FilePath)
	}

	if debugInfo.FileSize <= 0 {
		t.Error("Expected FileSize to be greater than 0")
	}

	if debugInfo.ParseTime <= 0 {
		t.Error("Expected ParseTime to be greater than 0")
	}

	if debugInfo.FactCount <= 0 {
		t.Error("Expected FactCount to be greater than 0")
	}

	// Check that facts were parsed
	if len(facts) == 0 {
		t.Error("Expected facts to be parsed")
	}
}

func TestDefaultParser_ParseContentWithTestFiles(t *testing.T) {
	parser := &DefaultParser{}

	// Test with the test files we created
	testFiles := []string{
		"testdata/simple_blocks.hcl",
		"testdata/nested_blocks.hcl",
		"testdata/facts_with_blocks.hcl",
	}

	for _, testFile := range testFiles {
		t.Run(testFile, func(t *testing.T) {
			content, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read test file %s: %v", testFile, err)
			}

			facts, err := parser.ParseContent(content)
			if err != nil {
				t.Fatalf("Failed to parse test file %s: %v", testFile, err)
			}

			// Basic validation - should have machine_id and collected_at
			if _, exists := facts["machine_id"]; !exists {
				t.Errorf("Expected machine_id to exist in %s", testFile)
			}

			if _, exists := facts["collected_at"]; !exists {
				t.Errorf("Expected collected_at to exist in %s", testFile)
			}

			// Check that blocks were parsed (should have at least one block)
			hasBlocks := false
			for key := range facts {
				if strings.Contains(key, "_") && key != "machine_id" && key != "collected_at" {
					hasBlocks = true
					break
				}
			}

			if !hasBlocks {
				t.Errorf("Expected blocks to be parsed in %s", testFile)
			}
		})
	}
}

// MockParser for testing
type MockParser struct{}

func (m *MockParser) ParseFile(_ string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"name":    "mock-server",
		"version": "1.0.0",
	}, nil
}

func (m *MockParser) ParseContent(_ []byte) (map[string]interface{}, error) {
	return map[string]interface{}{
		"name":    "mock-server",
		"version": "1.0.0",
	}, nil
}

func TestCollector_AddRemoveHCLFile(t *testing.T) {
	collector := NewCollector([]string{})

	// Test adding files
	collector.AddHCLFile("/path/to/file1.hcl")
	collector.AddHCLFile("/path/to/file2.hcl")

	if len(collector.hclFiles) != 2 {
		t.Errorf("Expected 2 files, got %d", len(collector.hclFiles))
	}

	// Test removing file
	collector.RemoveHCLFile("/path/to/file1.hcl")
	if len(collector.hclFiles) != 1 {
		t.Errorf("Expected 1 file after removal, got %d", len(collector.hclFiles))
	}

	if collector.hclFiles[0] != "/path/to/file2.hcl" {
		t.Errorf("Expected remaining file to be '/path/to/file2.hcl', got '%s'", collector.hclFiles[0])
	}

	// Test removing non-existent file
	collector.RemoveHCLFile("/non/existent/file.hcl")
	if len(collector.hclFiles) != 1 {
		t.Errorf("Expected 1 file after removing non-existent file, got %d", len(collector.hclFiles))
	}
}

func TestCollector_GetFactSources(t *testing.T) {
	hclFiles := []string{"/path/to/file1.hcl", "/path/to/file2.hcl"}
	collector := NewCollector(hclFiles)

	sources := collector.GetFactSources()
	if len(sources) != 2 {
		t.Errorf("Expected 2 sources, got %d", len(sources))
	}

	for i, source := range sources {
		if source != hclFiles[i] {
			t.Errorf("Expected source %d to be '%s', got '%s'", i, hclFiles[i], source)
		}
	}
}
