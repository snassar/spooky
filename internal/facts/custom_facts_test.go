package facts

import (
	"os"
	"path/filepath"
	spookyfactstypes "spooky/internal/facts/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomFactsCollectorWithHierarchicalStructure(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	factsFile := filepath.Join(tempDir, "facts.hcl")

	// Create a test facts.hcl file
	testFactsContent := `
app_name = "test-app"
app_version = "1.2.3"
environment = "testing"
config_path = "/etc/test-app/config.hcl"
log_path = "/var/log/test-app"
deployment_state = "active"
ssl_enabled = true
prometheus_port = 9090
`
	err := os.WriteFile(factsFile, []byte(testFactsContent), 0o600)
	require.NoError(t, err)

	// Create custom facts collector
	collector, err := NewCustomFactsCollector(factsFile)
	require.NoError(t, err)

	// Collect custom facts
	collection, err := collector.Collect("test-server")
	require.NoError(t, err)

	// Verify the collection has custom facts
	assert.True(t, collection.HasCustomFacts())
	assert.Equal(t, 1, len(collection.CustomFacts))

	// Verify custom facts are stored in hierarchical structure
	facts, exists := collection.GetCustomFactsByFile("facts")
	assert.True(t, exists)
	assert.Equal(t, "test-app", facts["app_name"])
	assert.Equal(t, "1.2.3", facts["app_version"])
	assert.Equal(t, "testing", facts["environment"])
	assert.Equal(t, "/etc/test-app/config.hcl", facts["config_path"])
	assert.Equal(t, "/var/log/test-app", facts["log_path"])
	assert.Equal(t, "active", facts["deployment_state"])
	assert.Equal(t, true, facts["ssl_enabled"])
	assert.Equal(t, 9090, facts["prometheus_port"])

	// Verify individual fact access
	value, exists := collection.GetCustomFact("facts", "app_name")
	assert.True(t, exists)
	assert.Equal(t, "test-app", value)

	// Verify backward compatibility with flat Facts map
	assert.Contains(t, collection.Facts, "custom.facts")
	fact := collection.Facts["custom.facts"]
	assert.Equal(t, "custom.facts", fact.Key)
	assert.Equal(t, "custom", fact.Source)
	assert.Equal(t, "test-server", fact.Server)
	assert.Equal(t, "facts", fact.Metadata["filename"])
}

func TestCustomFactsCollectorWithMissingFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	factsFile := filepath.Join(tempDir, "nonexistent.hcl")

	// Create custom facts collector
	collector, err := NewCustomFactsCollector(factsFile)
	require.NoError(t, err)

	// Collect custom facts (should return empty collection)
	collection, err := collector.Collect("test-server")
	require.NoError(t, err)

	// Verify the collection is empty
	assert.False(t, collection.HasCustomFacts())
	assert.Equal(t, 0, len(collection.CustomFacts))
}

func TestCustomFactsCollectorSpecificKeys(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	factsFile := filepath.Join(tempDir, "facts.hcl")

	// Create a test facts.hcl file
	testFactsContent := `
app_name = "test-app"
app_version = "1.2.3"
environment = "testing"
config_path = "/etc/test-app/config.hcl"
log_path = "/var/log/test-app"
`
	err := os.WriteFile(factsFile, []byte(testFactsContent), 0o600)
	require.NoError(t, err)

	// Create custom facts collector
	collector, err := NewCustomFactsCollector(factsFile)
	require.NoError(t, err)

	// Collect specific custom facts
	collection, err := collector.CollectSpecific("test-server", []string{"app_name", "environment"})
	require.NoError(t, err)

	// Verify only requested facts are collected
	assert.True(t, collection.HasCustomFacts())
	facts, exists := collection.GetCustomFactsByFile("facts")
	assert.True(t, exists)
	assert.Equal(t, 2, len(facts))
	assert.Equal(t, "test-app", facts["app_name"])
	assert.Equal(t, "testing", facts["environment"])
	assert.NotContains(t, facts, "app_version")
	assert.NotContains(t, facts, "config_path")
}

func TestCustomFactsCollectorGetFact(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	factsFile := filepath.Join(tempDir, "facts.hcl")

	// Create a test facts.hcl file
	testFactsContent := `
app_name = "test-app"
app_version = "1.2.3"
environment = "testing"
`
	err := os.WriteFile(factsFile, []byte(testFactsContent), 0o600)
	require.NoError(t, err)

	// Create custom facts collector
	collector, err := NewCustomFactsCollector(factsFile)
	require.NoError(t, err)

	// Get specific fact
	fact, err := collector.GetFact("test-server", "app_name")
	require.NoError(t, err)
	assert.Equal(t, "custom.app_name", fact.Key)
	assert.Equal(t, "test-app", fact.Value)
	assert.Equal(t, "custom", fact.Source)
	assert.Equal(t, "test-server", fact.Server)

	// Get fact with custom. prefix
	fact, err = collector.GetFact("test-server", "custom.app_version")
	require.NoError(t, err)
	assert.Equal(t, "custom.app_version", fact.Key)
	assert.Equal(t, "1.2.3", fact.Value)

	// Try to get non-existent fact
	_, err = collector.GetFact("test-server", "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "custom fact nonexistent not found")
}

func TestFactCollectionCustomFactsMethods(t *testing.T) {
	collection := &spookyfactstypes.FactCollection{
		Server:      "test-server",
		Timestamp:   time.Now(),
		Facts:       make(map[string]*spookyfactstypes.Fact),
		CustomFacts: make(map[string]map[string]interface{}),
	}

	// Test SetCustomFact
	collection.SetCustomFact("test-file", "key1", "value1")
	collection.SetCustomFact("test-file", "key2", "value2")
	collection.SetCustomFact("another-file", "key3", "value3")

	// Test GetCustomFact
	value, exists := collection.GetCustomFact("test-file", "key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", value)

	value, exists = collection.GetCustomFact("test-file", "key2")
	assert.True(t, exists)
	assert.Equal(t, "value2", value)

	value, exists = collection.GetCustomFact("another-file", "key3")
	assert.True(t, exists)
	assert.Equal(t, "value3", value)

	value, exists = collection.GetCustomFact("test-file", "nonexistent")
	assert.False(t, exists)
	assert.Nil(t, value)

	// Test GetCustomFactsByFile
	facts, exists := collection.GetCustomFactsByFile("test-file")
	assert.True(t, exists)
	assert.Equal(t, 2, len(facts))
	assert.Equal(t, "value1", facts["key1"])
	assert.Equal(t, "value2", facts["key2"])

	facts, exists = collection.GetCustomFactsByFile("nonexistent")
	assert.False(t, exists)
	assert.Nil(t, facts)

	// Test GetAllCustomFacts
	allFacts := collection.GetAllCustomFacts()
	assert.Equal(t, 2, len(allFacts))
	assert.Contains(t, allFacts, "test-file")
	assert.Contains(t, allFacts, "another-file")

	// Test HasCustomFacts
	assert.True(t, collection.HasCustomFacts())

	// Test with empty collection
	emptyCollection := &spookyfactstypes.FactCollection{
		Server:    "test-server",
		Timestamp: time.Now(),
		Facts:     make(map[string]*spookyfactstypes.Fact),
	}
	assert.False(t, emptyCollection.HasCustomFacts())
}
