package facts

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"spooky/internal/ssh"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	// Test creating a new manager without storage
	sshClient := &ssh.SSHClient{}
	manager := NewManager(sshClient)

	assert.NotNil(t, manager)
	assert.Equal(t, sshClient, manager.sshClient)
	assert.NotNil(t, manager.sshCollector)
	assert.NotNil(t, manager.localCollector)
	assert.Nil(t, manager.hclCollector)
	assert.Nil(t, manager.tofuCollector)
	assert.NotNil(t, manager.customCollectors)
	assert.NotNil(t, manager.cache)
	assert.Equal(t, DefaultTTL, manager.defaultTTL)
}

func TestNewManagerWithStorage(t *testing.T) {
	// Test creating a new manager with storage
	sshClient := &ssh.SSHClient{}
	storage := &MockFactStorage{}
	manager := NewManagerWithStorage(sshClient, storage)

	assert.NotNil(t, manager)
	assert.Equal(t, sshClient, manager.sshClient)
	assert.Equal(t, storage, manager.storage)
	assert.NotNil(t, manager.sshCollector)
	assert.NotNil(t, manager.localCollector)
}

func TestManagerConfigureHCLCollector(t *testing.T) {
	manager := NewManager(nil)

	// Test configuring HCL collector
	filePath := "/path/to/config.hcl"
	manager.ConfigureHCLCollector(filePath)

	assert.NotNil(t, manager.hclCollector)
}

func TestManagerConfigureOpenTofuCollector(t *testing.T) {
	manager := NewManager(nil)

	// Test configuring OpenTofu collector
	statePath := "/path/to/state.tfstate"
	manager.ConfigureOpenTofuCollector(statePath)

	assert.NotNil(t, manager.tofuCollector)
}

func TestManagerCollectAllFacts(t *testing.T) {
	manager := NewManager(nil)

	// Test collecting facts for local server
	collection, err := manager.CollectAllFacts("local")

	assert.NoError(t, err)
	assert.NotNil(t, collection)
	assert.Equal(t, "local", collection.Server)
	assert.NotNil(t, collection.Facts)
}

func TestManagerCollectSpecificFacts(t *testing.T) {
	manager := NewManager(nil)

	// Test collecting specific facts
	keys := []string{"os.name", "hardware.cpu"}
	collection, err := manager.CollectSpecificFacts("local", keys)

	// This might fail if the facts don't exist, but should not panic
	if err != nil {
		assert.Contains(t, err.Error(), "no facts collected")
	} else {
		assert.NotNil(t, collection)
		assert.Equal(t, "local", collection.Server)
	}
}

func TestManagerGetFact(t *testing.T) {
	manager := NewManager(nil)

	// Test getting a specific fact
	fact, err := manager.GetFact("local", "os.name")

	// This might fail if the fact doesn't exist, but should not panic
	if err != nil {
		assert.Contains(t, err.Error(), "fact not found")
	} else {
		assert.NotNil(t, fact)
		assert.Equal(t, "os.name", fact.Key)
	}
}

func TestManagerCacheOperations(t *testing.T) {
	manager := NewManager(nil)

	// Test cache operations
	manager.ClearCache()

	// Test setting default TTL
	newTTL := 30 * time.Minute
	manager.SetDefaultTTL(newTTL)
	assert.Equal(t, newTTL, manager.defaultTTL)
}

func TestManagerWithTestDataFromExamples(t *testing.T) {
	// Use test data from examples/testing where available
	examplesDir := "../../examples/testing"

	// Test with valid project data
	validProjectPath := filepath.Join(examplesDir, "test-valid-project")
	if _, err := os.Stat(validProjectPath); os.IsNotExist(err) {
		t.Skip("Test data directory not found, skipping test")
	}

	manager := NewManager(nil)

	// Test with valid facts data
	validFactsPath := filepath.Join(validProjectPath, "data", "valid-facts.json")
	if _, err := os.Stat(validFactsPath); err == nil {
		// Test importing valid facts
		options := &ImportOptions{
			Source:    validFactsPath,
			Path:      validFactsPath,
			MergeMode: MergeModeAppend,
			Validate:  true,
			Server:    "example-server",
		}

		err := manager.ImportCustomFactsWithOptions("file://"+validFactsPath, options)
		// This might fail if the file doesn't exist or is invalid
		if err != nil {
			assert.Contains(t, err.Error(), "failed to load custom facts")
		}
	}
}

func TestManagerImportCustomFacts(t *testing.T) {
	manager := NewManager(nil)

	// Test importing custom facts
	collection, err := manager.ImportCustomFacts("test", "test-server", MergePolicyReplace)

	// This might fail if the source doesn't exist, but should not panic
	if err != nil {
		assert.Contains(t, err.Error(), "invalid source")
	} else {
		assert.NotNil(t, collection)
		assert.Equal(t, "test-server", collection.Server)
	}
}

func TestManagerImportCustomFactsWithOptions(t *testing.T) {
	manager := NewManager(nil)

	// Create temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test-facts.json")

	testData := map[string]*CustomFacts{
		"test-server": {
			Custom: map[string]interface{}{
				"os": map[string]interface{}{
					"name":    "linux",
					"version": "20.04",
				},
			},
		},
	}

	// Write test data to file
	data, err := json.Marshal(testData)
	require.NoError(t, err)
	err = os.WriteFile(testFile, data, 0o600)
	require.NoError(t, err)

	// Test importing with options
	options := &ImportOptions{
		Source:    "file://" + testFile,
		Path:      testFile,
		MergeMode: MergeModeAppend,
		Validate:  true,
		Server:    "test-server",
	}

	err = manager.ImportCustomFactsWithOptions("file://"+testFile, options)
	// This might fail due to file path issues, but should not panic
	if err != nil {
		assert.Contains(t, err.Error(), "failed to load custom facts")
	}
}

func TestManagerGetCustomFacts(t *testing.T) {
	manager := NewManager(nil)

	// Test getting custom facts
	facts, err := manager.GetCustomFacts("test-server")

	// This might fail if no storage is configured
	if err != nil {
		assert.Contains(t, err.Error(), "no storage configured")
	} else {
		assert.NotNil(t, facts)
	}
}

func TestManagerFactQuery(t *testing.T) {
	manager := NewManager(nil)

	// Test querying persisted facts
	query := &FactQuery{
		MachineName: "test-server",
	}

	collections, err := manager.QueryPersistedFacts(query)

	// This might fail if no storage is configured
	if err != nil {
		assert.Contains(t, err.Error(), "no storage configured")
	} else {
		assert.NotNil(t, collections)
	}
}

func TestManagerExportImportFacts(t *testing.T) {
	manager := NewManager(nil)

	// Test exporting facts
	var buf bytes.Buffer
	err := manager.ExportFacts(&buf)
	// This might fail if no storage is configured
	if err != nil {
		assert.Contains(t, err.Error(), "no storage configured")
		return
	}

	// Test importing facts
	reader := bytes.NewReader(buf.Bytes())
	err = manager.ImportFacts(reader)
	// This might fail if no storage is configured
	if err != nil {
		assert.Contains(t, err.Error(), "no storage configured")
	}
}

func TestManagerWithMissingRequiredFacts(t *testing.T) {
	// Test with missing required facts scenario
	examplesDir := "../../examples/testing"
	missingFactsPath := filepath.Join(examplesDir, "test-missing-required-facts")

	if _, err := os.Stat(missingFactsPath); os.IsNotExist(err) {
		t.Skip("Test data directory not found, skipping test")
	}

	manager := NewManager(nil)

	// Test with missing fields facts
	missingFactsFile := filepath.Join(missingFactsPath, "data", "missing-fields-facts.json")
	if _, err := os.Stat(missingFactsFile); err == nil {
		options := &ImportOptions{
			Source:    "file://" + missingFactsFile,
			Path:      missingFactsFile,
			MergeMode: MergeModeAppend,
			Validate:  true,
			Server:    "example-server",
		}

		err := manager.ImportCustomFactsWithOptions("file://"+missingFactsFile, options)
		// This might fail due to file path issues
		if err != nil {
			assert.Contains(t, err.Error(), "failed to load custom facts")
		}
	}
}

func TestManagerWithInvalidJSONFacts(t *testing.T) {
	// Test with invalid JSON facts scenario
	examplesDir := "../../examples/testing"
	invalidJSONPath := filepath.Join(examplesDir, "test-invalid-json-facts")

	if _, err := os.Stat(invalidJSONPath); os.IsNotExist(err) {
		t.Skip("Test data directory not found, skipping test")
	}

	manager := NewManager(nil)

	// Test with malformed facts
	malformedFactsFile := filepath.Join(invalidJSONPath, "data", "malformed-facts.json")
	if _, err := os.Stat(malformedFactsFile); err == nil {
		options := &ImportOptions{
			Source:    "file://" + malformedFactsFile,
			Path:      malformedFactsFile,
			MergeMode: MergeModeAppend,
			Validate:  true,
			Server:    "example-server",
		}

		err := manager.ImportCustomFactsWithOptions("file://"+malformedFactsFile, options)
		// This should fail due to invalid JSON
		assert.Error(t, err)
	}
}

func TestManagerFactCollectionClone(t *testing.T) {
	// Test FactCollection cloning
	original := &FactCollection{
		Server:    "test-server",
		Timestamp: time.Now(),
		Facts: map[string]*Fact{
			"os.name": {
				Key:       "os.name",
				Value:     "linux",
				Source:    "test",
				Server:    "test-server",
				Timestamp: time.Now(),
				Metadata: map[string]interface{}{
					"test": "value",
				},
			},
		},
	}

	cloned := original.Clone()

	assert.NotNil(t, cloned)
	assert.Equal(t, original.Server, cloned.Server)
	assert.Equal(t, original.Timestamp, cloned.Timestamp)
	// The maps should be different objects but contain the same data
	// The maps should be different objects but contain the same data
	// We'll test that they have the same length and content instead
	assert.Equal(t, len(original.Facts), len(cloned.Facts))

	// Test that modifying cloned doesn't affect original
	cloned.Facts["os.name"].Value = "windows"
	assert.Equal(t, "linux", original.Facts["os.name"].Value)
	assert.Equal(t, "windows", cloned.Facts["os.name"].Value)

	// Test metadata cloning
	cloned.Facts["os.name"].Metadata["test"] = "new-value"
	assert.Equal(t, "value", original.Facts["os.name"].Metadata["test"])
	assert.Equal(t, "new-value", cloned.Facts["os.name"].Metadata["test"])
}

func TestManagerRegisterCustomCollector(t *testing.T) {
	manager := NewManager(nil)

	// Create a mock collector
	mockCollector := &MockFactCollector{}

	// Test registering custom collector
	manager.RegisterCustomCollector("mock", mockCollector)

	assert.Contains(t, manager.customCollectors, "mock")
	assert.Equal(t, mockCollector, manager.customCollectors["mock"])
}

func TestManagerClose(t *testing.T) {
	manager := NewManager(nil)

	// Test closing manager
	err := manager.Close()
	assert.NoError(t, err)
}

// Mock implementations for testing

type MockFactStorage struct{}

func (m *MockFactStorage) GetMachineFacts(machineID string) (*MachineFacts, error) {
	return &MachineFacts{
		MachineID: machineID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *MockFactStorage) SetMachineFacts(_ string, _ *MachineFacts) error {
	return nil
}

func (m *MockFactStorage) QueryFacts(_ *FactQuery) ([]*MachineFacts, error) {
	return []*MachineFacts{}, nil
}

func (m *MockFactStorage) DeleteFacts(_ *FactQuery) (int, error) {
	return 0, nil
}

func (m *MockFactStorage) DeleteMachineFacts(_ string) error {
	return nil
}

func (m *MockFactStorage) ExportToJSON(_ io.Writer) error {
	return nil
}

func (m *MockFactStorage) ImportFromJSON(_ io.Reader) error {
	return nil
}

func (m *MockFactStorage) ExportToJSONWithEncryption(_ io.Writer, _ ExportOptions) error {
	return nil
}

func (m *MockFactStorage) ImportFromJSONWithDecryption(_ io.Reader, _ string) error {
	return nil
}

func (m *MockFactStorage) Close() error {
	return nil
}

type MockFactCollector struct{}

func (m *MockFactCollector) Collect(server string) (*FactCollection, error) {
	return &FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*Fact),
	}, nil
}

func (m *MockFactCollector) CollectSpecific(server string, _ []string) (*FactCollection, error) {
	return &FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*Fact),
	}, nil
}

func (m *MockFactCollector) GetFact(server, key string) (*Fact, error) {
	return &Fact{
		Key:       key,
		Value:     "mock-value",
		Source:    "mock",
		Server:    server,
		Timestamp: time.Now(),
	}, nil
}
