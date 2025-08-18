# Performance Optimization Plan 2: Facts Caching

**Generated:** 2025-08-18  
**Recommendation:** Add Facts Caching  
**Priority:** Medium  
**Effort:** Low  
**Impact:** Medium performance improvement  
**Status:** Planning

## Overview

This plan addresses the performance bottleneck of repeated fact collection by implementing intelligent caching with TTL, memory management, and automatic cleanup mechanisms.

## Current State Analysis

### Problem Statement
- No caching mechanism for repeated fact collections
- Location: `internal/facts/manager.go:CollectFacts`
- Current Pattern: Always collect fresh facts, no caching mechanism
- Impact: Medium - Redundant fact collection without caching

### Current Implementation
```go
// CURRENT - No caching, always collect fresh facts
func (m *Manager) CollectFacts(ctx context.Context, machine interface{}) (*spookytypes.FactCollection, error) {
    // ... validation logic ...
    
    // Always collect fresh facts - NO CACHING
    facts, err := m.collector.Collect(ctx, machineObj)
    if err != nil {
        m.logger.Error("Failed to collect facts", err, map[string]interface{}{"machine": machineObj.Hostname})
        return nil, fmt.Errorf("failed to collect facts for %s: %w", machineObj.Hostname, err)
    }
    
    // ... validation logic ...
    
    return facts, nil
}
```

## Target State

### Desired Implementation
```go
// TARGET - Intelligent caching with TTL and memory management
type CachedFactManager struct {
    manager    *Manager
    cache      map[string]*CacheEntry
    mutex      sync.RWMutex
    ttl        time.Duration
    maxEntries int
    logger     spookytypeslogging.Logger
}

type CacheEntry struct {
    facts      *spookytypes.FactCollection
    expiresAt  time.Time
    lastAccess time.Time
    accessCount int64
}

func (c *CachedFactManager) CollectFacts(ctx context.Context, machine interface{}) (*spookytypes.FactCollection, error) {
    machineObj, ok := machine.(*spookytypes.Machine)
    if !ok {
        return nil, fmt.Errorf("machine must be of type *spookytypes.Machine")
    }
    
    cacheKey := fmt.Sprintf("%s:%s", machineObj.Hostname, machineObj.Host)
    
    // Check cache first
    c.mutex.RLock()
    if entry, exists := c.cache[cacheKey]; exists && time.Now().Before(entry.expiresAt) {
        entry.lastAccess = time.Now()
        entry.accessCount++
        c.mutex.RUnlock()
        c.logger.Debug("Returning cached facts", map[string]interface{}{
            "machine": machineObj.Hostname,
            "cache_hit": true,
        })
        return entry.facts, nil
    }
    c.mutex.RUnlock()
    
    // Collect fresh facts
    facts, err := c.manager.CollectFacts(ctx, machine)
    if err != nil {
        return nil, err
    }
    
    // Cache the results
    c.mutex.Lock()
    if len(c.cache) >= c.maxEntries {
        c.evictLeastRecentlyUsed()
    }
    c.cache[cacheKey] = &CacheEntry{
        facts:      facts,
        expiresAt:  time.Now().Add(c.ttl),
        lastAccess: time.Now(),
        accessCount: 1,
    }
    c.mutex.Unlock()
    
    return facts, nil
}
```

## Implementation Plan

### Phase 1: Cache Design and Architecture (Day 1-2)

#### 1.1 Cache Strategy Design
- **Task:** Design caching strategy and architecture
- **Deliverable:** Cache design document
- **Acceptance Criteria:** Clear caching strategy with TTL, eviction, and memory management
- **Effort:** 0.5 days

#### 1.2 Cache Key Design
- **Task:** Design cache key generation strategy
- **Deliverable:** Cache key design document
- **Acceptance Criteria:** Unique and consistent cache keys for machines
- **Effort:** 0.5 days

#### 1.3 Memory Management Strategy
- **Task:** Design memory management and eviction strategy
- **Deliverable:** Memory management design document
- **Acceptance Criteria:** LRU eviction with memory limits
- **Effort:** 0.5 days

#### 1.4 Configuration Design
- **Task:** Design cache configuration options
- **Deliverable:** Configuration design document
- **Acceptance Criteria:** Configurable TTL, max entries, and cleanup intervals
- **Effort:** 0.5 days

### Phase 2: Core Cache Implementation (Day 3-4)

#### 2.1 Implement Cache Structure
- **Task:** Implement cache data structures
- **File:** `internal/facts/cache.go`
- **Deliverable:** Cache entry and manager structures
- **Acceptance Criteria:** Thread-safe cache structures
- **Effort:** 0.5 days

#### 2.2 Implement Cache Operations
- **Task:** Implement cache get/set operations
- **File:** `internal/facts/cache.go`
- **Deliverable:** Cache get/set implementation
- **Acceptance Criteria:** Thread-safe cache operations
- **Effort:** 1 day

#### 2.3 Implement Cache Key Generation
- **Task:** Implement cache key generation
- **File:** `internal/facts/cache.go`
- **Deliverable:** Cache key generation implementation
- **Acceptance Criteria:** Unique and consistent cache keys
- **Effort:** 0.5 days

#### 2.4 Implement Cache Manager
- **Task:** Implement cached fact manager wrapper
- **File:** `internal/facts/cached_manager.go`
- **Deliverable:** Cached fact manager implementation
- **Acceptance Criteria:** Cached fact manager with proper interface compliance
- **Effort:** 1 day

### Phase 3: Memory Management and Cleanup (Day 5-6)

#### 3.1 Implement LRU Eviction
- **Task:** Implement least recently used eviction
- **File:** `internal/facts/cache.go`
- **Deliverable:** LRU eviction implementation
- **Acceptance Criteria:** Proper LRU eviction when cache is full
- **Effort:** 0.5 days

#### 3.2 Implement TTL Expiration
- **Task:** Implement time-based expiration
- **File:** `internal/facts/cache.go`
- **Deliverable:** TTL expiration implementation
- **Acceptance Criteria:** Automatic expiration of cached entries
- **Effort:** 0.5 days

#### 3.3 Implement Cleanup Routine
- **Task:** Implement background cleanup routine
- **File:** `internal/facts/cache.go`
- **Deliverable:** Background cleanup implementation
- **Acceptance Criteria:** Automatic cleanup of expired entries
- **Effort:** 0.5 days

#### 3.4 Implement Memory Monitoring
- **Task:** Implement memory usage monitoring
- **File:** `internal/facts/cache.go`
- **Deliverable:** Memory monitoring implementation
- **Acceptance Criteria:** Memory usage tracking and limits
- **Effort:** 0.5 days

### Phase 4: Integration and Testing (Day 7-8)

#### 4.1 Integration with Fact Manager
- **Task:** Integrate cache with existing fact manager
- **File:** `internal/facts/manager.go`
- **Deliverable:** Cache integration
- **Acceptance Criteria:** Seamless integration with existing fact collection
- **Effort:** 0.5 days

#### 4.2 Unit Testing
- **Task:** Implement unit tests for cache functionality
- **File:** `internal/facts/cache_test.go`
- **Deliverable:** Comprehensive unit tests
- **Acceptance Criteria:** 100% test coverage for cache functionality
- **Effort:** 1 day

#### 4.3 Integration Testing
- **Task:** Test cache with real fact collection
- **File:** `tests/integration/facts_cache_test.go`
- **Deliverable:** Integration test suite
- **Acceptance Criteria:** Cache works with real fact collection scenarios
- **Effort:** 0.5 days

#### 4.4 Performance Testing
- **Task:** Measure cache performance improvement
- **File:** `tests/performance/facts_cache_test.go`
- **Deliverable:** Performance benchmarks
- **Acceptance Criteria:** Measurable performance improvement
- **Effort:** 0.5 days

### Phase 5: Monitoring and Optimization (Day 9-10)

#### 5.1 Add Cache Metrics
- **Task:** Add cache performance metrics
- **File:** `internal/facts/metrics.go`
- **Deliverable:** Cache metrics implementation
- **Acceptance Criteria:** Comprehensive cache performance metrics
- **Effort:** 0.5 days

#### 5.2 Add Cache Logging
- **Task:** Add cache operation logging
- **File:** `internal/facts/cache.go`
- **Deliverable:** Cache logging implementation
- **Acceptance Criteria:** Detailed cache operation logging
- **Effort:** 0.5 days

#### 5.3 Cache Configuration
- **Task:** Add cache configuration options
- **File:** `internal/types/facts/facts.go`
- **Deliverable:** Cache configuration types
- **Acceptance Criteria:** Configurable cache parameters
- **Effort:** 0.5 days

#### 5.4 Documentation Update
- **Task:** Update documentation for cache functionality
- **File:** `docs/FACTS_SYSTEM.md`
- **Deliverable:** Updated documentation
- **Acceptance Criteria:** Clear documentation of cache behavior
- **Effort:** 0.5 days

## Technical Implementation Details

### Cache Structure
```go
// Cache entry structure
type CacheEntry struct {
    facts      *spookytypes.FactCollection
    expiresAt  time.Time
    lastAccess time.Time
    accessCount int64
    size       int64 // Memory size in bytes
}

// Cache manager structure
type CachedFactManager struct {
    manager    spookyinterfaces.FactsIntegration
    cache      map[string]*CacheEntry
    mutex      sync.RWMutex
    ttl        time.Duration
    maxEntries int
    maxMemory  int64 // Maximum memory usage in bytes
    currentMemory int64
    logger     spookytypeslogging.Logger
    metrics    *CacheMetrics
}

// Cache metrics
type CacheMetrics struct {
    hits       int64
    misses     int64
    evictions  int64
    expirations int64
    mutex      sync.RWMutex
}
```

### Cache Operations
```go
// Cache get operation
func (c *CachedFactManager) get(key string) (*spookytypes.FactCollection, bool) {
    c.mutex.RLock()
    entry, exists := c.cache[key]
    c.mutex.RUnlock()
    
    if !exists {
        c.metrics.recordMiss()
        return nil, false
    }
    
    // Check expiration
    if time.Now().After(entry.expiresAt) {
        c.mutex.Lock()
        delete(c.cache, key)
        c.currentMemory -= entry.size
        c.mutex.Unlock()
        c.metrics.recordExpiration()
        return nil, false
    }
    
    // Update access information
    c.mutex.Lock()
    entry.lastAccess = time.Now()
    entry.accessCount++
    c.mutex.Unlock()
    
    c.metrics.recordHit()
    return entry.facts, true
}

// Cache set operation
func (c *CachedFactManager) set(key string, facts *spookytypes.FactCollection) error {
    // Calculate entry size
    entrySize := c.calculateEntrySize(facts)
    
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    // Check if we need to evict entries
    for len(c.cache) >= c.maxEntries || (c.currentMemory+entrySize) > c.maxMemory {
        if !c.evictLeastRecentlyUsed() {
            return fmt.Errorf("unable to make space in cache")
        }
    }
    
    // Create cache entry
    entry := &CacheEntry{
        facts:      facts,
        expiresAt:  time.Now().Add(c.ttl),
        lastAccess: time.Now(),
        accessCount: 1,
        size:       entrySize,
    }
    
    // Remove existing entry if present
    if existing, exists := c.cache[key]; exists {
        c.currentMemory -= existing.size
    }
    
    // Add new entry
    c.cache[key] = entry
    c.currentMemory += entrySize
    
    return nil
}
```

### LRU Eviction
```go
// LRU eviction implementation
func (c *CachedFactManager) evictLeastRecentlyUsed() bool {
    if len(c.cache) == 0 {
        return false
    }
    
    var oldestKey string
    var oldestTime time.Time
    
    for key, entry := range c.cache {
        if oldestKey == "" || entry.lastAccess.Before(oldestTime) {
            oldestKey = key
            oldestTime = entry.lastAccess
        }
    }
    
    if oldestKey != "" {
        entry := c.cache[oldestKey]
        delete(c.cache, oldestKey)
        c.currentMemory -= entry.size
        c.metrics.recordEviction()
        return true
    }
    
    return false
}
```

### Cleanup Routine
```go
// Background cleanup routine
func (c *CachedFactManager) startCleanupRoutine() {
    go func() {
        ticker := time.NewTicker(5 * time.Minute)
        defer ticker.Stop()
        
        for range ticker.C {
            c.cleanupExpiredEntries()
        }
    }()
}

// Cleanup expired entries
func (c *CachedFactManager) cleanupExpiredEntries() {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    now := time.Now()
    expiredCount := 0
    
    for key, entry := range c.cache {
        if now.After(entry.expiresAt) {
            delete(c.cache, key)
            c.currentMemory -= entry.size
            expiredCount++
        }
    }
    
    if expiredCount > 0 {
        c.logger.Debug("Cleaned up expired cache entries", map[string]interface{}{
            "expired_count": expiredCount,
            "remaining_entries": len(c.cache),
        })
    }
}
```

### Cache Key Generation
```go
// Cache key generation
func (c *CachedFactManager) generateCacheKey(machine *spookytypes.Machine) string {
    // Include machine-specific information in cache key
    keyComponents := []string{
        machine.Hostname,
        machine.Host,
        machine.User,
    }
    
    // Add machine tags if present
    if len(machine.Tags) > 0 {
        sort.Strings(machine.Tags)
        keyComponents = append(keyComponents, strings.Join(machine.Tags, ","))
    }
    
    // Create hash of key components
    hash := sha256.Sum256([]byte(strings.Join(keyComponents, "|")))
    return fmt.Sprintf("facts:%x", hash[:8])
}
```

## Testing Strategy

### Unit Tests
```go
func TestCacheOperations(t *testing.T) {
    // Test cache get/set operations
    // Test cache expiration
    // Test LRU eviction
    // Test memory limits
}

func TestCacheKeyGeneration(t *testing.T) {
    // Test cache key uniqueness
    // Test cache key consistency
    // Test cache key with different machine configurations
}
```

### Integration Tests
```go
func TestCacheWithRealFactCollection(t *testing.T) {
    // Test cache with real fact collection
    // Test cache hit/miss scenarios
    // Test cache performance improvement
}
```

### Performance Tests
```go
func BenchmarkCachePerformance(b *testing.B) {
    // Benchmark cache hit performance
    // Benchmark cache miss performance
    // Benchmark cache memory usage
}
```

## Success Metrics

### Performance Metrics
- **Target:** 80% cache hit rate for repeated collections
- **Measurement:** Cache hit/miss ratio
- **Baseline:** Current fact collection performance
- **Monitoring:** Cache performance metrics

### Memory Metrics
- **Memory Usage:** Cache memory usage within configured limits
- **Eviction Rate:** LRU eviction rate < 10%
- **Expiration Rate:** TTL expiration rate < 20%
- **Memory Efficiency:** Memory usage per cached entry

### Quality Metrics
- **Cache Consistency:** 100% cache consistency with fresh data
- **Error Rate:** Cache error rate < 0.1%
- **Response Time:** Cache hit response time < 1ms
- **Reliability:** 99.9% cache operation success rate

## Configuration Options

### Cache Configuration
```go
type FactsCacheConfig struct {
    Enabled     bool          `hcl:"enabled,optional"`
    TTL         time.Duration `hcl:"ttl,optional"`
    MaxEntries  int           `hcl:"max_entries,optional"`
    MaxMemory   int64         `hcl:"max_memory,optional"`
    CleanupInterval time.Duration `hcl:"cleanup_interval,optional"`
}

func DefaultFactsCacheConfig() *FactsCacheConfig {
    return &FactsCacheConfig{
        Enabled:        true,
        TTL:            30 * time.Minute,
        MaxEntries:     1000,
        MaxMemory:      100 * 1024 * 1024, // 100MB
        CleanupInterval: 5 * time.Minute,
    }
}
```

## Risk Assessment

### Technical Risks
- **Memory Leaks:** Risk of memory leaks in cache
- **Mitigation:** Proper cleanup routines and memory monitoring
- **Cache Inconsistency:** Risk of stale cache data
- **Mitigation:** Configurable TTL and cache invalidation

### Performance Risks
- **Cache Thrashing:** Risk of excessive cache evictions
- **Mitigation:** Proper cache sizing and LRU eviction
- **Memory Pressure:** Risk of memory pressure from cache
- **Mitigation:** Memory limits and monitoring

### Functional Risks
- **Data Staleness:** Risk of serving stale fact data
- **Mitigation:** Configurable TTL and cache invalidation
- **Cache Corruption:** Risk of cache data corruption
- **Mitigation:** Proper error handling and cache validation

## Rollback Plan

### Rollback Triggers
- Cache hit rate < 50%
- Memory usage > 200MB
- Cache error rate > 1%
- Performance degradation > 10%

### Rollback Procedure
1. **Immediate:** Disable cache via configuration
2. **Short-term:** Revert to non-cached fact collection
3. **Long-term:** Fix issues and re-enable cache

## Dependencies

### Internal Dependencies
- `internal/facts/manager.go` - Core fact collection
- `internal/types/facts/` - Fact type definitions
- `internal/logging/` - Logging infrastructure

### External Dependencies
- Go sync package for thread safety
- Go crypto/sha256 for cache key generation
- No additional external dependencies required

## Timeline

### Week 1: Implementation
- Day 1-2: Cache design and architecture
- Day 3-4: Core cache implementation
- Day 5-6: Memory management and cleanup

### Week 2: Testing and Optimization
- Day 7-8: Integration and testing
- Day 9-10: Monitoring and optimization

## Conclusion

This improvement plan provides a systematic approach to implementing intelligent facts caching, addressing the performance bottleneck of repeated fact collection while maintaining data consistency and memory efficiency.

**Expected Outcomes:**
- 80% cache hit rate for repeated fact collections
- Reduced fact collection overhead
- Configurable memory management
- Enhanced user experience with faster fact access

The plan ensures a robust caching implementation that integrates seamlessly with the existing fact collection system while providing significant performance improvements.
