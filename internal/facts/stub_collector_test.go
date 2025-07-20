package facts

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewStubCollector(t *testing.T) {
	collector := NewStubCollector("test-collector")

	assert.NotNil(t, collector)
	assert.Equal(t, "test-collector", collector.name)
}

func TestStubCollectorCollect(t *testing.T) {
	collector := NewStubCollector("test-collector")

	// Test Collect method
	collection, err := collector.Collect("test-server")

	assert.NoError(t, err)
	assert.NotNil(t, collection)
	assert.Equal(t, "test-server", collection.Server)
	assert.NotNil(t, collection.Facts)
	assert.Empty(t, collection.Facts) // Should be empty map
	assert.WithinDuration(t, time.Now(), collection.Timestamp, time.Second)
}

func TestStubCollectorCollectSpecific(t *testing.T) {
	collector := NewStubCollector("test-collector")

	// Test CollectSpecific method
	keys := []string{"hostname", "os.name", "cpu.cores"}
	collection, err := collector.CollectSpecific("test-server", keys)

	assert.NoError(t, err)
	assert.NotNil(t, collection)
	assert.Equal(t, "test-server", collection.Server)
	assert.NotNil(t, collection.Facts)
	assert.Empty(t, collection.Facts) // Should be empty map regardless of keys
	assert.WithinDuration(t, time.Now(), collection.Timestamp, time.Second)
}

func TestStubCollectorGetFact(t *testing.T) {
	collector := NewStubCollector("test-collector")

	// Test GetFact method
	fact, err := collector.GetFact("test-server", "hostname")

	assert.Error(t, err)
	assert.Nil(t, fact)
	assert.Contains(t, err.Error(), "test-collector fact collection not yet implemented")

	// Test with different fact key
	fact, err = collector.GetFact("test-server", "os.name")

	assert.Error(t, err)
	assert.Nil(t, fact)
	assert.Contains(t, err.Error(), "test-collector fact collection not yet implemented")
}

func TestStubCollectorWithDifferentNames(t *testing.T) {
	// Test with different collector names
	collectors := []string{"hcl-collector", "opentofu-collector", "custom-collector"}

	for _, name := range collectors {
		t.Run(name, func(t *testing.T) {
			collector := NewStubCollector(name)

			// Test Collect
			collection, err := collector.Collect("test-server")
			assert.NoError(t, err)
			assert.NotNil(t, collection)

			// Test GetFact error message
			_, err = collector.GetFact("test-server", "hostname")
			assert.Error(t, err)
			assert.Contains(t, err.Error(), name+" fact collection not yet implemented")
		})
	}
}

func TestStubCollectorEmptyKeys(t *testing.T) {
	collector := NewStubCollector("test-collector")

	// Test CollectSpecific with empty keys
	collection, err := collector.CollectSpecific("test-server", []string{})

	assert.NoError(t, err)
	assert.NotNil(t, collection)
	assert.Equal(t, "test-server", collection.Server)
	assert.Empty(t, collection.Facts)
}

func TestStubCollectorNilKeys(t *testing.T) {
	collector := NewStubCollector("test-collector")

	// Test CollectSpecific with nil keys
	collection, err := collector.CollectSpecific("test-server", nil)

	assert.NoError(t, err)
	assert.NotNil(t, collection)
	assert.Equal(t, "test-server", collection.Server)
	assert.Empty(t, collection.Facts)
}

func TestStubCollectorEmptyServer(t *testing.T) {
	collector := NewStubCollector("test-collector")

	// Test with empty server name
	collection, err := collector.Collect("")

	assert.NoError(t, err)
	assert.NotNil(t, collection)
	assert.Equal(t, "", collection.Server)
	assert.Empty(t, collection.Facts)

	// Test GetFact with empty server
	fact, err := collector.GetFact("", "hostname")
	assert.Error(t, err)
	assert.Nil(t, fact)
}
