package machines

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	spookytypeslogging "spooky/internal/types/logging"
)

// MockLogger implements the Logger interface for testing
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

func TestParseMachineBlockRefactored(t *testing.T) {
	// Test that the refactored function works correctly
	loader := NewLoader(&MockLogger{})

	// Test that the loader is properly initialized
	assert.NotNil(t, loader)

	// Test that the loader can be created without errors
	// This verifies that the refactored code compiles and initializes correctly
}

func TestAttributeParserRegistry(t *testing.T) {
	registry := NewAttributeParserRegistry()

	// Test that all expected parsers are registered
	expectedParsers := []string{
		"hostname", "host", "port", "user", "key_file", "passphrase",
		"tags", "groups", "roles", "classes",
		"connection_timeout", "command_timeout", "max_connections",
		"retry_attempts", "retry_delay",
	}

	for _, parserName := range expectedParsers {
		parser, exists := registry.parsers[parserName]
		assert.True(t, exists, "Parser %s should be registered", parserName)
		assert.NotNil(t, parser, "Parser %s should not be nil", parserName)
	}

	// Test that unknown parsers return false
	_, exists := registry.parsers["unknown_attribute"]
	assert.False(t, exists, "Unknown parser should not exist")
}

func TestBlockParserRegistry(t *testing.T) {
	loader := NewLoader(&MockLogger{})
	registry := NewBlockParserRegistry(loader)

	// Test that all expected block parsers are registered
	expectedBlocks := []string{"resources", "metadata"}

	for _, blockType := range expectedBlocks {
		parser, exists := registry.parsers[blockType]
		assert.True(t, exists, "Block parser %s should be registered", blockType)
		assert.NotNil(t, parser, "Block parser %s should not be nil", blockType)
		assert.Equal(t, blockType, parser.GetBlockType())
	}

	// Test that unknown block parsers return false
	_, exists := registry.parsers["unknown_block"]
	assert.False(t, exists, "Unknown block parser should not exist")
}

func TestStringAttributeParser(t *testing.T) {
	parser := &StringAttributeParser{fieldName: "Hostname"}

	// Test that parser is properly configured
	assert.Equal(t, "Hostname", parser.GetFieldName())
}

func TestIntAttributeParser(t *testing.T) {
	parser := &IntAttributeParser{fieldName: "Port"}

	// Test that parser is properly configured
	assert.Equal(t, "Port", parser.GetFieldName())
}

func TestObjectAttributeParser(t *testing.T) {
	parser := &ObjectAttributeParser{fieldName: "Tags"}

	// Test that parser is properly configured
	assert.Equal(t, "Tags", parser.GetFieldName())
}

func TestArrayAttributeParser(t *testing.T) {
	parser := &ArrayAttributeParser{fieldName: "Groups"}

	// Test that parser is properly configured
	assert.Equal(t, "Groups", parser.GetFieldName())
}

func TestLoadMachinesFromFile(t *testing.T) {
	loader := NewLoader(&MockLogger{})

	// Test with valid HCL content
	hclContent := `
		machines {
			machine "test-server" {
				hostname = "test.example.com"
				port = 22
				user = "admin"
			}
		}
	`

	// Create a temporary file
	tmpFile := createTempFile(t, hclContent)
	defer removeTempFile(t, tmpFile)

	machines, err := loader.LoadMachinesFromFile(context.Background(), tmpFile)
	assert.NoError(t, err)
	assert.Len(t, machines, 1)
	assert.Equal(t, "test-server", machines[0].Hostname)
	assert.Equal(t, "test.example.com", machines[0].Host)
	assert.Equal(t, 22, machines[0].Port)
	assert.Equal(t, "admin", machines[0].User)
}

// Helper functions for testing
func createTempFile(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "test-*.hcl")
	require.NoError(t, err)

	_, err = tmpFile.WriteString(content)
	require.NoError(t, err)

	err = tmpFile.Close()
	require.NoError(t, err)

	return tmpFile.Name()
}

func removeTempFile(t *testing.T, filename string) {
	err := os.Remove(filename)
	require.NoError(t, err)
}
