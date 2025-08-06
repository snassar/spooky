package facts

import (
	"testing"
	"time"

	spookyfactstypes "spooky/internal/facts/types"
	spookylogging "spooky/internal/logging"
	spookyloggingtypes "spooky/internal/logging/types"

	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	// Test creating a new manager without storage
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	manager := NewManager(nil, logger)

	assert.NotNil(t, manager)
}

func TestNewManagerWithStorage(t *testing.T) {
	// Test creating a new manager with storage
	storage := &MockFactStorage{}
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	manager := NewManagerWithStorage(nil, storage, logger)

	assert.NotNil(t, manager)
}

func TestManagerCollectAllFacts(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	manager := NewManager(nil, logger)

	// Test collecting facts for local server
	collection, err := manager.CollectAllFacts("local")

	assert.NoError(t, err)
	assert.NotNil(t, collection)
	assert.Equal(t, "local", collection.Server)
	assert.NotNil(t, collection.Facts)
}

func TestManagerCollectSpecificFacts(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	manager := NewManager(nil, logger)

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
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	manager := NewManager(nil, logger)

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
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	manager := NewManager(nil, logger)

	// Test cache operations
	manager.ClearCache()

	// Test setting default TTL
	newTTL := 30 * time.Minute
	manager.SetDefaultTTL(newTTL)
	assert.Equal(t, newTTL, manager.defaultTTL)
}

func TestManagerFactCollectionClone(t *testing.T) {
	// Test FactCollection cloning
	original := &spookyfactstypes.FactCollection{
		Server:    "test-server",
		Timestamp: time.Now(),
		Facts: map[string]*spookyfactstypes.Fact{
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
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	manager := NewManager(nil, logger)

	// Create a mock collector
	mockCollector := &MockFactCollector{}

	// Test registering custom collector
	manager.RegisterCustomCollector("mock", mockCollector)

	assert.Contains(t, manager.customCollectors, "mock")
	assert.Equal(t, mockCollector, manager.customCollectors["mock"])
}

func TestManagerClose(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	manager := NewManager(nil, logger)

	// Test closing manager
	err := manager.Close()
	assert.NoError(t, err)
}

// Mock implementations for testing

type MockFactStorage struct{}

func (m *MockFactStorage) GetFactCollection(machineID string) (*spookyfactstypes.FactCollection, error) {
	return &spookyfactstypes.FactCollection{
		Server:    machineID,
		Timestamp: time.Now(),
		Facts:     make(map[string]*spookyfactstypes.Fact),
	}, nil
}

func (m *MockFactStorage) SetFactCollection(_ string, _ *spookyfactstypes.FactCollection) error {
	return nil
}

func (m *MockFactStorage) QueryFactCollections(_ interface{}) ([]*spookyfactstypes.FactCollection, error) {
	return []*spookyfactstypes.FactCollection{}, nil
}

func (m *MockFactStorage) DeleteFactCollections(_ *spookyfactstypes.FactQuery) (int, error) {
	return 0, nil
}

func (m *MockFactStorage) DeleteFactCollection(_ string) error {
	return nil
}

func (m *MockFactStorage) ExportToJSON(_ interface{}) error {
	return nil
}

func (m *MockFactStorage) ImportFromJSON(_ interface{}) error {
	return nil
}

func (m *MockFactStorage) ExportToJSONWithEncryption(_ interface{}, _ interface{}) error {
	return nil
}

func (m *MockFactStorage) ImportFromJSONWithDecryption(_ interface{}, _ string) error {
	return nil
}

func (m *MockFactStorage) Close() error {
	return nil
}

func (m *MockFactStorage) ImportFromHCL(_ interface{}) error {
	return nil
}

type MockFactCollector struct{}

func (m *MockFactCollector) Collect(server string) (*spookyfactstypes.FactCollection, error) {
	return &spookyfactstypes.FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*spookyfactstypes.Fact),
	}, nil
}

func (m *MockFactCollector) CollectSpecific(server string, _ []string) (*spookyfactstypes.FactCollection, error) {
	return &spookyfactstypes.FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*spookyfactstypes.Fact),
	}, nil
}

func (m *MockFactCollector) GetFact(server, key string) (*spookyfactstypes.Fact, error) {
	return &spookyfactstypes.Fact{
		Key:       key,
		Value:     "mock-value",
		Source:    "mock",
		Server:    server,
		Timestamp: time.Now(),
	}, nil
}
