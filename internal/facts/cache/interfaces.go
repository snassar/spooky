package cache

import (
	"sync"
	"time"

	"spooky/internal/facts/types"
)

// Cache defines the interface for fact caching operations
type Cache interface {
	// Core cache operations
	Get(server string) (*types.FactCollection, error)
	Set(server string, collection *types.FactCollection) error
	Delete(server string) error
	Clear() error

	// Cache management
	ClearExpired() error
	GetAll() ([]*types.FactCollection, error)
	GetFiltered(server string, keys []string) (*types.FactCollection, error)

	// Configuration
	SetTTL(ttl time.Duration)
	GetTTL() time.Duration
}

// Manager provides cache management functionality
type Manager struct {
	cache Cache
	ttl   time.Duration
	mutex sync.RWMutex
}

// NewManager creates a new cache manager
func NewManager(cache Cache) *Manager {
	return &Manager{
		cache: cache,
		ttl:   30 * time.Minute, // Default TTL
	}
}

// GetCache returns the underlying cache
func (m *Manager) GetCache() Cache {
	return m.cache
}

// SetCache sets the underlying cache
func (m *Manager) SetCache(cache Cache) {
	m.cache = cache
}

// SetTTL sets the cache TTL
func (m *Manager) SetTTL(ttl time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.ttl = ttl
	if m.cache != nil {
		m.cache.SetTTL(ttl)
	}
}

// GetTTL gets the cache TTL
func (m *Manager) GetTTL() time.Duration {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.ttl
}
