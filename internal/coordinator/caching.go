package coordinator

import (
	"sync"
	"time"
)

// CacheManager provides centralized cache management
type CacheManager struct {
	caches map[string]*Cache
	mu     sync.RWMutex
}

// Cache represents a simple in-memory cache
type Cache struct {
	data    map[string]*CacheEntry
	mu      sync.RWMutex
	ttl     time.Duration
	maxSize int
}

// CacheEntry represents a cached item
type CacheEntry struct {
	Value       interface{}
	ExpiresAt   time.Time
	CreatedAt   time.Time
	AccessCount int
}

// NewCacheManager creates a new cache manager
func NewCacheManager() *CacheManager {
	return &CacheManager{
		caches: make(map[string]*Cache),
	}
}

// GetCache gets or creates a cache with the specified name
func (cm *CacheManager) GetCache(name string, ttl time.Duration, maxSize int) *Cache {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cache, exists := cm.caches[name]; exists {
		return cache
	}

	cache := &Cache{
		data:    make(map[string]*CacheEntry),
		ttl:     ttl,
		maxSize: maxSize,
	}

	cm.caches[name] = cache
	return cache
}

// Set sets a value in the cache
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if we need to evict items
	if len(c.data) >= c.maxSize {
		c.evictOldest()
	}

	c.data[key] = &CacheEntry{
		Value:       value,
		ExpiresAt:   time.Now().Add(c.ttl),
		CreatedAt:   time.Now(),
		AccessCount: 0,
	}
}

// Get gets a value from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.data[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		delete(c.data, key)
		return nil, false
	}

	// Update access count
	entry.AccessCount++
	return entry.Value, true
}

// Delete removes a value from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// Clear clears all entries from the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]*CacheEntry)
}

// Size returns the number of entries in the cache
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data)
}

// evictOldest removes the oldest entry from the cache
func (c *Cache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range c.data {
		if oldestKey == "" || entry.CreatedAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.CreatedAt
		}
	}

	if oldestKey != "" {
		delete(c.data, oldestKey)
	}
}

// ClearAllCaches clears all caches in the manager
func (cm *CacheManager) ClearAllCaches() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for _, cache := range cm.caches {
		cache.Clear()
	}
}

// GetCacheStats returns statistics for all caches
func (cm *CacheManager) GetCacheStats() map[string]interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	stats := make(map[string]interface{})
	for name, cache := range cm.caches {
		stats[name] = map[string]interface{}{
			"size": cache.Size(),
			"ttl":  cache.ttl.String(),
		}
	}

	return stats
}
