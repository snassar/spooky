package coordinator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCacheManager(t *testing.T) {
	manager := NewCacheManager()
	assert.NotNil(t, manager)
}

func TestCacheManager_GetCache(t *testing.T) {
	manager := NewCacheManager()

	// Test getting a new cache
	cache := manager.GetCache("test-cache", 1*time.Hour, 100)
	assert.NotNil(t, cache)

	// Test getting the same cache again
	cache2 := manager.GetCache("test-cache", 1*time.Hour, 100)
	assert.Equal(t, cache, cache2)
}

func TestCache_SetAndGet(t *testing.T) {
	manager := NewCacheManager()
	cache := manager.GetCache("test-cache", 1*time.Hour, 100)

	// Test setting and getting a value
	key := "test-key"
	value := "test-value"

	cache.Set(key, value)

	retrieved, found := cache.Get(key)
	assert.True(t, found)
	assert.Equal(t, value, retrieved)
}

func TestCache_Expiration(t *testing.T) {
	manager := NewCacheManager()
	cache := manager.GetCache("expire-cache", 1*time.Millisecond, 100)

	// Test expiration
	key := "expire-key"
	value := "expire-value"

	cache.Set(key, value)

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	_, found := cache.Get(key)
	assert.False(t, found)
}

func TestCache_Delete(t *testing.T) {
	manager := NewCacheManager()
	cache := manager.GetCache("delete-cache", 1*time.Hour, 100)

	key := "delete-key"
	value := "delete-value"

	cache.Set(key, value)

	// Verify it exists
	_, found := cache.Get(key)
	assert.True(t, found)

	// Delete it
	cache.Delete(key)

	// Verify it's gone
	_, found = cache.Get(key)
	assert.False(t, found)
}

func TestCache_Clear(t *testing.T) {
	manager := NewCacheManager()
	cache := manager.GetCache("clear-cache", 1*time.Hour, 100)

	// Add multiple items
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	// Verify they exist
	_, found1 := cache.Get("key1")
	_, found2 := cache.Get("key2")
	_, found3 := cache.Get("key3")
	assert.True(t, found1)
	assert.True(t, found2)
	assert.True(t, found3)

	// Clear cache
	cache.Clear()

	// Verify they're gone
	_, found1 = cache.Get("key1")
	_, found2 = cache.Get("key2")
	_, found3 = cache.Get("key3")
	assert.False(t, found1)
	assert.False(t, found2)
	assert.False(t, found3)
}

func TestCache_Size(t *testing.T) {
	manager := NewCacheManager()
	cache := manager.GetCache("size-cache", 1*time.Hour, 100)

	// Initially empty
	assert.Equal(t, 0, cache.Size())

	// Add items
	cache.Set("key1", "value1")
	assert.Equal(t, 1, cache.Size())

	cache.Set("key2", "value2")
	assert.Equal(t, 2, cache.Size())

	// Delete item
	cache.Delete("key1")
	assert.Equal(t, 1, cache.Size())
}

func TestCacheManager_ClearAllCaches(t *testing.T) {
	manager := NewCacheManager()

	// Create multiple caches
	cache1 := manager.GetCache("cache1", 1*time.Hour, 100)
	cache2 := manager.GetCache("cache2", 1*time.Hour, 100)

	// Add items to caches
	cache1.Set("key1", "value1")
	cache2.Set("key2", "value2")

	// Verify items exist
	_, found1 := cache1.Get("key1")
	_, found2 := cache2.Get("key2")
	assert.True(t, found1)
	assert.True(t, found2)

	// Clear all caches
	manager.ClearAllCaches()

	// Verify items are gone
	_, found1 = cache1.Get("key1")
	_, found2 = cache2.Get("key2")
	assert.False(t, found1)
	assert.False(t, found2)
}

func TestCacheManager_GetCacheStats(t *testing.T) {
	manager := NewCacheManager()

	// Create a cache and add some items
	cache := manager.GetCache("stats-cache", 1*time.Hour, 100)
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	stats := manager.GetCacheStats()
	assert.NotNil(t, stats)
	assert.IsType(t, map[string]interface{}{}, stats)
}
