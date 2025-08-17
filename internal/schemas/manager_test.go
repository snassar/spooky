package schemas

import (
	"os"
	"path/filepath"
	"testing"

	spookytypeslogging "spooky/internal/types/logging"
)

// MockLogger is a mock implementation of Logger for testing
type MockLogger struct{}

func (m *MockLogger) Debug(msg string, fields ...map[string]interface{})                 {}
func (m *MockLogger) Info(msg string, fields ...map[string]interface{})                  {}
func (m *MockLogger) Warn(msg string, fields ...map[string]interface{})                  {}
func (m *MockLogger) Error(msg string, err error, fields ...map[string]interface{})      {}
func (m *MockLogger) Fatal(msg string, err error, fields ...map[string]interface{})      {}
func (m *MockLogger) WithFields(fields map[string]interface{}) spookytypeslogging.Logger { return m }
func (m *MockLogger) WithComponent(component string) spookytypeslogging.Logger           { return m }
func (m *MockLogger) WithOperation(operation string) spookytypeslogging.Logger           { return m }
func (m *MockLogger) SetLevel(level spookytypeslogging.LogLevel)                         {}
func (m *MockLogger) GetLevel() spookytypeslogging.LogLevel                              { return spookytypeslogging.LogLevelInfo }

func TestSchemaMetadataValidation(t *testing.T) {
	// Create a mock logger for testing
	logger := &MockLogger{}

	// Create a schema manager
	manager := NewManager(logger)

	// Test loading a valid schema with metadata
	validSchemaContent := `
# Test Schema
metadata {
  schema_version = "0.20250809.0"
  schema_type = "test-schema"
  schema_name = "Test Schema"
  last_updated = "2024-01-01"
  description = "A test schema for validation with proper metadata structure"
}

# Schema content - use a proper schema structure
test_schema {
  field = {
    type = "string"
    required = true
    description = "Test field"
  }
}
`

	schema, err := manager.LoadFromString(validSchemaContent)
	if err != nil {
		t.Fatalf("Failed to load valid schema: %v", err)
	}

	if schema == nil {
		t.Fatal("Expected schema to be loaded, got nil")
	}

	t.Logf("Successfully loaded schema: %s", schema.Name)
}

func TestSchemaMetadataValidationFailure(t *testing.T) {
	// Create a mock logger for testing
	logger := &MockLogger{}

	// Create a schema manager
	manager := NewManager(logger)

	// Test loading a schema without metadata (should fail)
	invalidSchemaContent := `
# Test Schema without metadata
test_schema {
  field = {
    type = "string"
    required = true
    description = "Test field"
  }
}
`

	_, err := manager.LoadFromString(invalidSchemaContent)
	if err == nil {
		t.Fatal("Expected error when loading schema without metadata, got nil")
	}

	expectedError := "schema metadata validation failed: schema must contain a metadata block"
	if err.Error() != expectedError {
		t.Fatalf("Expected error '%s', got '%s'", expectedError, err.Error())
	}

	t.Logf("Successfully caught missing metadata error: %v", err)
}

func TestSchemaMetadataValidationInvalidVersion(t *testing.T) {
	// Create a mock logger for testing
	logger := &MockLogger{}

	// Create a schema manager
	manager := NewManager(logger)

	// Test loading a schema with invalid version format (should fail)
	invalidSchemaContent := `
# Test Schema with invalid version
metadata {
  schema_version = "1.0.0"  # Invalid format, should be 0.YYYYMMDD.N
  schema_type = "test-schema"
  schema_name = "Test Schema"
  last_updated = "2024-01-01"
  description = "A test schema for validation"
}

# Schema content
test_schema {
  field = {
    type = "string"
    required = true
    description = "Test field"
  }
}
`

	_, err := manager.LoadFromString(invalidSchemaContent)
	if err == nil {
		t.Fatal("Expected error when loading schema with invalid version format, got nil")
	}

	t.Logf("Successfully caught invalid version format error: %v", err)
}

func TestLoadSchemasFromDirectoryWithMetadataValidation(t *testing.T) {
	// Create a debug logger that outputs to t.Log
	debugLogger := &DebugLogger{t: t}

	// Create a schema manager
	manager := NewManager(debugLogger)

	// Test loading schemas from the actual schemas directory
	// Use absolute path to ensure it works regardless of working directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// The test runs from internal/schemas, so we need to go up one level
	// and then into internal/schemas/schemas/structure
	schemasDir := filepath.Join(currentDir, "schemas", "structure")

	// Check if the structure schemas directory exists
	if _, err := os.Stat(schemasDir); os.IsNotExist(err) {
		t.Skipf("Structure schemas directory %s does not exist, skipping test", schemasDir)
	}

	schemas, err := manager.LoadSchemasFromDirectory(schemasDir)
	if err != nil {
		t.Fatalf("Failed to load schemas from directory: %v", err)
	}

	if len(schemas) == 0 {
		t.Fatal("Expected to load at least one schema, got none")
	}

	t.Logf("Successfully loaded %d schemas with metadata validation", len(schemas))

	// Verify that each schema has been validated
	for name, schema := range schemas {
		if schema == nil {
			t.Errorf("Schema %s is nil", name)
			continue
		}

		t.Logf("Schema %s loaded successfully: %s", name, schema.Name)
	}
}

// DebugLogger outputs log messages to the test log
type DebugLogger struct {
	t *testing.T
}

func (d *DebugLogger) Debug(msg string, fields ...map[string]interface{}) {
	d.t.Logf("DEBUG: %s %v", msg, fields)
}

func (d *DebugLogger) Info(msg string, fields ...map[string]interface{}) {
	d.t.Logf("INFO: %s %v", msg, fields)
}

func (d *DebugLogger) Warn(msg string, fields ...map[string]interface{}) {
	d.t.Logf("WARN: %s %v", msg, fields)
}

func (d *DebugLogger) Error(msg string, err error, fields ...map[string]interface{}) {
	if err != nil {
		d.t.Logf("ERROR: %s %v %v", msg, err, fields)
	} else {
		d.t.Logf("ERROR: %s %v", msg, fields)
	}
}

func (d *DebugLogger) Fatal(msg string, err error, fields ...map[string]interface{}) {
	d.t.Logf("FATAL: %s %v %v", msg, err, fields)
}

func (d *DebugLogger) WithFields(fields map[string]interface{}) spookytypeslogging.Logger {
	return d
}

func (d *DebugLogger) WithComponent(component string) spookytypeslogging.Logger {
	return d
}

func (d *DebugLogger) WithOperation(operation string) spookytypeslogging.Logger {
	return d
}

func (d *DebugLogger) SetLevel(level spookytypeslogging.LogLevel) {}

func (d *DebugLogger) GetLevel() spookytypeslogging.LogLevel {
	return spookytypeslogging.LogLevelInfo
}
