# Performance Optimization Plan 4: Template Caching

**Generated:** 2025-08-18  
**Recommendation:** Template Caching  
**Priority:** Low  
**Effort:** Low  
**Impact:** Low (performance improvement for repeated renders)  
**Status:** Planning

## Overview

This plan addresses the performance bottleneck of repeated template parsing and compilation by implementing intelligent template caching with compilation caching, AST caching, and memory management.

## Current State Analysis

### Problem Statement
- No caching mechanism for template parsing and compilation
- Location: `internal/templates/manager.go`
- Current Pattern: Parse and compile templates on every render
- Impact: Low - Redundant template parsing without caching

### Current Implementation
```go
// CURRENT - No caching, parse and compile on every render
func (m *Manager) RenderTemplate(ctx context.Context, templatePath string, data map[string]interface{}) (string, error) {
    // ... validation logic ...
    
    // Always parse and compile template - NO CACHING
    template, err := m.parseTemplate(templatePath)
    if err != nil {
        return "", fmt.Errorf("failed to parse template %s: %w", templatePath, err)
    }
    
    // Always compile template - NO CACHING
    compiled, err := m.compileTemplate(template)
    if err != nil {
        return "", fmt.Errorf("failed to compile template %s: %w", templatePath, err)
    }
    
    // Render template
    result, err := m.renderTemplate(compiled, data)
    if err != nil {
        return "", fmt.Errorf("failed to render template %s: %w", templatePath, err)
    }
    
    return result, nil
}
```

## Target State

### Desired Implementation
```go
// TARGET - Intelligent template caching with compilation caching
type CachedTemplateManager struct {
    manager    *Manager
    cache      map[string]*TemplateCacheEntry
    mutex      sync.RWMutex
    maxEntries int
    logger     spookytypeslogging.Logger
}

type TemplateCacheEntry struct {
    template     *template.Template
    compiled     interface{}
    lastModified time.Time
    lastAccess   time.Time
    accessCount  int64
    size         int64
}

func (c *CachedTemplateManager) RenderTemplate(ctx context.Context, templatePath string, data map[string]interface{}) (string, error) {
    cacheKey := c.generateCacheKey(templatePath)
    
    // Check cache first
    c.mutex.RLock()
    if entry, exists := c.cache[cacheKey]; exists {
        // Check if template file has been modified
        if c.isTemplateModified(templatePath, entry.lastModified) {
            c.mutex.RUnlock()
            // Template modified, remove from cache
            c.mutex.Lock()
            delete(c.cache, cacheKey)
            c.mutex.Unlock()
        } else {
            // Use cached template
            entry.lastAccess = time.Now()
            entry.accessCount++
            c.mutex.RUnlock()
            
            c.logger.Debug("Using cached template", map[string]interface{}{
                "template": templatePath,
                "cache_hit": true,
            })
            
            return c.renderCachedTemplate(entry, data)
        }
    } else {
        c.mutex.RUnlock()
    }
    
    // Parse and compile fresh template
    template, err := c.manager.parseTemplate(templatePath)
    if err != nil {
        return "", err
    }
    
    compiled, err := c.manager.compileTemplate(template)
    if err != nil {
        return "", err
    }
    
    // Cache the compiled template
    c.mutex.Lock()
    if len(c.cache) >= c.maxEntries {
        c.evictLeastRecentlyUsed()
    }
    
    fileInfo, _ := os.Stat(templatePath)
    c.cache[cacheKey] = &TemplateCacheEntry{
        template:     template,
        compiled:     compiled,
        lastModified: fileInfo.ModTime(),
        lastAccess:   time.Now(),
        accessCount:  1,
        size:         c.calculateEntrySize(template, compiled),
    }
    c.mutex.Unlock()
    
    return c.renderCachedTemplate(c.cache[cacheKey], data)
}
```

## Implementation Plan

### Phase 1: Cache Design and Architecture (Day 1)

#### 1.1 Template Cache Strategy Design
- **Task:** Design template caching strategy and architecture
- **Deliverable:** Cache design document
- **Acceptance Criteria:** Clear caching strategy with file modification detection
- **Effort:** 0.5 days

#### 1.2 Compilation Cache Design
- **Task:** Design compilation caching strategy
- **Deliverable:** Compilation cache design document
- **Acceptance Criteria:** Strategy for caching compiled templates
- **Effort:** 0.5 days

### Phase 2: Core Cache Implementation (Day 2-3)

#### 2.1 Implement Cache Structure
- **Task:** Implement template cache data structures
- **File:** `internal/templates/cache.go`
- **Deliverable:** Cache entry and manager structures
- **Acceptance Criteria:** Thread-safe cache structures
- **Effort:** 0.5 days

#### 2.2 Implement Template Cache Operations
- **Task:** Implement template cache get/set operations
- **File:** `internal/templates/cache.go`
- **Deliverable:** Template cache operations implementation
- **Acceptance Criteria:** Thread-safe template cache operations
- **Effort:** 1 day

#### 2.3 Implement File Modification Detection
- **Task:** Implement file modification detection
- **File:** `internal/templates/cache.go`
- **Deliverable:** File modification detection implementation
- **Acceptance Criteria:** Accurate file modification detection
- **Effort:** 0.5 days

#### 2.4 Implement Cache Manager
- **Task:** Implement cached template manager wrapper
- **File:** `internal/templates/cached_manager.go`
- **Deliverable:** Cached template manager implementation
- **Acceptance Criteria:** Cached template manager with proper interface compliance
- **Effort:** 1 day

### Phase 3: Compilation Caching (Day 4-5)

#### 3.1 Implement Compilation Cache
- **Task:** Implement compilation result caching
- **File:** `internal/templates/compilation.go`
- **Deliverable:** Compilation cache implementation
- **Acceptance Criteria:** Cached compilation results
- **Effort:** 0.5 days

#### 3.2 Implement AST Caching
- **Task:** Implement AST caching for templates
- **File:** `internal/templates/ast_cache.go`
- **Deliverable:** AST cache implementation
- **Acceptance Criteria:** Cached AST structures
- **Effort:** 0.5 days

#### 3.3 Implement Memory Management
- **Task:** Implement memory management for cached templates
- **File:** `internal/templates/cache.go`
- **Deliverable:** Memory management implementation
- **Acceptance Criteria:** Proper memory management for cached templates
- **Effort:** 0.5 days

#### 3.4 Implement Cache Cleanup
- **Task:** Implement cache cleanup routine
- **File:** `internal/templates/cache.go`
- **Deliverable:** Cache cleanup implementation
- **Acceptance Criteria:** Automatic cleanup of unused templates
- **Effort:** 0.5 days

### Phase 4: Integration and Testing (Day 6-7)

#### 4.1 Integration with Template Manager
- **Task:** Integrate cache with existing template manager
- **File:** `internal/templates/manager.go`
- **Deliverable:** Cache integration
- **Acceptance Criteria:** Seamless integration with existing template rendering
- **Effort:** 0.5 days

#### 4.2 Unit Testing
- **Task:** Implement unit tests for template cache
- **File:** `internal/templates/cache_test.go`
- **Deliverable:** Comprehensive unit tests
- **Acceptance Criteria:** 100% test coverage for cache functionality
- **Effort:** 1 day

#### 4.3 Integration Testing
- **Task:** Test cache with real template scenarios
- **File:** `tests/integration/template_cache_test.go`
- **Deliverable:** Integration test suite
- **Acceptance Criteria:** Cache works with real template scenarios
- **Effort:** 0.5 days

#### 4.4 Performance Testing
- **Task:** Measure cache performance improvement
- **File:** `tests/performance/template_cache_test.go`
- **Deliverable:** Performance benchmarks
- **Acceptance Criteria:** Measurable performance improvement
- **Effort:** 0.5 days

### Phase 5: Monitoring and Optimization (Day 8-10)

#### 5.1 Add Cache Metrics
- **Task:** Add template cache performance metrics
- **File:** `internal/templates/metrics.go`
- **Deliverable:** Cache metrics implementation
- **Acceptance Criteria:** Comprehensive cache performance metrics
- **Effort:** 0.5 days

#### 5.2 Add Cache Logging
- **Task:** Add template cache operation logging
- **File:** `internal/templates/cache.go`
- **Deliverable:** Cache logging implementation
- **Acceptance Criteria:** Detailed cache operation logging
- **Effort:** 0.5 days

#### 5.3 Template Cache Configuration
- **Task:** Add template cache configuration options
- **File:** `internal/types/templates/template.go`
- **Deliverable:** Cache configuration types
- **Acceptance Criteria:** Configurable cache parameters
- **Effort:** 0.5 days

#### 5.4 Documentation Update
- **Task:** Update documentation for template cache
- **File:** `docs/TEMPLATES_SYSTEM.md`
- **Deliverable:** Updated documentation
- **Acceptance Criteria:** Clear documentation of cache behavior
- **Effort:** 0.5 days

## Technical Implementation Details

### Cache Structure
```go
// Template cache entry structure
type TemplateCacheEntry struct {
    template     *template.Template
    compiled     interface{}
    ast          interface{} // Cached AST
    lastModified time.Time
    lastAccess   time.Time
    accessCount  int64
    size         int64
    checksum     string // File checksum for modification detection
}

// Template cache manager structure
type CachedTemplateManager struct {
    manager    spookyinterfaces.TemplatesIntegration
    cache      map[string]*TemplateCacheEntry
    mutex      sync.RWMutex
    maxEntries int
    maxMemory  int64
    currentMemory int64
    logger     spookytypeslogging.Logger
    metrics    *TemplateCacheMetrics
}

// Template cache metrics
type TemplateCacheMetrics struct {
    hits       int64
    misses     int64
    evictions  int64
    compilations int64
    mutex      sync.RWMutex
}
```

### Cache Operations
```go
// Template cache get operation
func (c *CachedTemplateManager) get(key string) (*TemplateCacheEntry, bool) {
    c.mutex.RLock()
    entry, exists := c.cache[key]
    c.mutex.RUnlock()
    
    if !exists {
        c.metrics.recordMiss()
        return nil, false
    }
    
    // Check file modification
    if c.isTemplateModified(entry.templatePath, entry.lastModified) {
        c.mutex.Lock()
        delete(c.cache, key)
        c.currentMemory -= entry.size
        c.mutex.Unlock()
        c.metrics.recordEviction()
        return nil, false
    }
    
    // Update access information
    c.mutex.Lock()
    entry.lastAccess = time.Now()
    entry.accessCount++
    c.mutex.Unlock()
    
    c.metrics.recordHit()
    return entry, true
}

// Template cache set operation
func (c *CachedTemplateManager) set(key string, template *template.Template, compiled interface{}, ast interface{}) error {
    // Calculate entry size
    entrySize := c.calculateEntrySize(template, compiled, ast)
    
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    // Check if we need to evict entries
    for len(c.cache) >= c.maxEntries || (c.currentMemory+entrySize) > c.maxMemory {
        if !c.evictLeastRecentlyUsed() {
            return fmt.Errorf("unable to make space in cache")
        }
    }
    
    // Get file information
    fileInfo, err := os.Stat(templatePath)
    if err != nil {
        return fmt.Errorf("failed to stat template file: %w", err)
    }
    
    // Calculate checksum
    checksum, err := c.calculateFileChecksum(templatePath)
    if err != nil {
        return fmt.Errorf("failed to calculate file checksum: %w", err)
    }
    
    // Create cache entry
    entry := &TemplateCacheEntry{
        template:     template,
        compiled:     compiled,
        ast:          ast,
        lastModified: fileInfo.ModTime(),
        lastAccess:   time.Now(),
        accessCount:  1,
        size:         entrySize,
        checksum:     checksum,
    }
    
    // Remove existing entry if present
    if existing, exists := c.cache[key]; exists {
        c.currentMemory -= existing.size
    }
    
    // Add new entry
    c.cache[key] = entry
    c.currentMemory += entrySize
    c.metrics.recordCompilation()
    
    return nil
}
```

### File Modification Detection
```go
// File modification detection
func (c *CachedTemplateManager) isTemplateModified(templatePath string, lastModified time.Time) bool {
    fileInfo, err := os.Stat(templatePath)
    if err != nil {
        // File doesn't exist or can't be accessed, consider it modified
        return true
    }
    
    return fileInfo.ModTime().After(lastModified)
}

// File checksum calculation
func (c *CachedTemplateManager) calculateFileChecksum(templatePath string) (string, error) {
    data, err := os.ReadFile(templatePath)
    if err != nil {
        return "", err
    }
    
    hash := sha256.Sum256(data)
    return fmt.Sprintf("%x", hash), nil
}

// Checksum-based modification detection
func (c *CachedTemplateManager) isTemplateModifiedByChecksum(templatePath string, cachedChecksum string) bool {
    currentChecksum, err := c.calculateFileChecksum(templatePath)
    if err != nil {
        return true
    }
    
    return currentChecksum != cachedChecksum
}
```

### Compilation Caching
```go
// Compilation cache implementation
type CompilationCache struct {
    cache map[string]*CompiledTemplate
    mutex sync.RWMutex
}

type CompiledTemplate struct {
    template *template.Template
    ast      interface{}
    compiled interface{}
    size     int64
}

func (cc *CompilationCache) Get(key string) (*CompiledTemplate, bool) {
    cc.mutex.RLock()
    defer cc.mutex.RUnlock()
    
    compiled, exists := cc.cache[key]
    return compiled, exists
}

func (cc *CompilationCache) Set(key string, template *template.Template, ast interface{}, compiled interface{}) {
    cc.mutex.Lock()
    defer cc.mutex.Unlock()
    
    size := cc.calculateSize(template, ast, compiled)
    cc.cache[key] = &CompiledTemplate{
        template: template,
        ast:      ast,
        compiled: compiled,
        size:     size,
    }
}

func (cc *CompilationCache) calculateSize(template *template.Template, ast interface{}, compiled interface{}) int64 {
    // Calculate approximate memory size
    var size int64
    
    // Template size
    if template != nil {
        size += int64(len(template.DefinedTemplates()) * 1024) // Approximate
    }
    
    // AST size
    if ast != nil {
        astData, _ := json.Marshal(ast)
        size += int64(len(astData))
    }
    
    // Compiled size
    if compiled != nil {
        compiledData, _ := json.Marshal(compiled)
        size += int64(len(compiledData))
    }
    
    return size
}
```

### LRU Eviction
```go
// LRU eviction for template cache
func (c *CachedTemplateManager) evictLeastRecentlyUsed() bool {
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

### Cache Cleanup
```go
// Background cleanup routine
func (c *CachedTemplateManager) startCleanupRoutine() {
    go func() {
        ticker := time.NewTicker(10 * time.Minute)
        defer ticker.Stop()
        
        for range ticker.C {
            c.cleanupUnusedTemplates()
        }
    }()
}

// Cleanup unused templates
func (c *CachedTemplateManager) cleanupUnusedTemplates() {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    cutoff := time.Now().Add(-30 * time.Minute) // Templates unused for 30 minutes
    removedCount := 0
    
    for key, entry := range c.cache {
        if entry.lastAccess.Before(cutoff) {
            delete(c.cache, key)
            c.currentMemory -= entry.size
            removedCount++
        }
    }
    
    if removedCount > 0 {
        c.logger.Debug("Cleaned up unused template cache entries", map[string]interface{}{
            "removed_count": removedCount,
            "remaining_entries": len(c.cache),
        })
    }
}
```

### Cache Key Generation
```go
// Template cache key generation
func (c *CachedTemplateManager) generateCacheKey(templatePath string) string {
    // Include template path and modification time in cache key
    fileInfo, err := os.Stat(templatePath)
    if err != nil {
        // Fallback to path-based key
        return fmt.Sprintf("template:%s", templatePath)
    }
    
    // Create hash of path and modification time
    keyData := fmt.Sprintf("%s:%d", templatePath, fileInfo.ModTime().Unix())
    hash := sha256.Sum256([]byte(keyData))
    return fmt.Sprintf("template:%x", hash[:8])
}
```

## Testing Strategy

### Unit Tests
```go
func TestTemplateCacheOperations(t *testing.T) {
    // Test template cache get/set operations
    // Test file modification detection
    // Test LRU eviction
    // Test memory limits
}

func TestFileModificationDetection(t *testing.T) {
    // Test file modification detection
    // Test checksum-based detection
    // Test cache invalidation on file changes
}
```

### Integration Tests
```go
func TestTemplateCacheWithRealTemplates(t *testing.T) {
    // Test cache with real template files
    // Test cache hit/miss scenarios
    // Test cache performance improvement
}
```

### Performance Tests
```go
func BenchmarkTemplateCachePerformance(b *testing.B) {
    // Benchmark cache hit performance
    // Benchmark cache miss performance
    // Benchmark template compilation caching
}
```

## Success Metrics

### Performance Metrics
- **Target:** 90% cache hit rate for repeated template renders
- **Measurement:** Cache hit/miss ratio
- **Baseline:** Current template parsing performance
- **Monitoring:** Template cache performance metrics

### Memory Metrics
- **Memory Usage:** Cache memory usage within configured limits
- **Eviction Rate:** LRU eviction rate < 5%
- **Compilation Cache Hit Rate:** 95% compilation cache hit rate
- **Memory Efficiency:** Memory usage per cached template

### Quality Metrics
- **Cache Consistency:** 100% cache consistency with file modifications
- **Error Rate:** Cache error rate < 0.1%
- **Response Time:** Cache hit response time < 1ms
- **Reliability:** 99.9% cache operation success rate

## Configuration Options

### Template Cache Configuration
```go
type TemplateCacheConfig struct {
    Enabled     bool          `hcl:"enabled,optional"`
    MaxEntries  int           `hcl:"max_entries,optional"`
    MaxMemory   int64         `hcl:"max_memory,optional"`
    CleanupInterval time.Duration `hcl:"cleanup_interval,optional"`
    CompilationCache bool     `hcl:"compilation_cache,optional"`
    ASTCache      bool        `hcl:"ast_cache,optional"`
}

func DefaultTemplateCacheConfig() *TemplateCacheConfig {
    return &TemplateCacheConfig{
        Enabled:        true,
        MaxEntries:     100,
        MaxMemory:      50 * 1024 * 1024, // 50MB
        CleanupInterval: 10 * time.Minute,
        CompilationCache: true,
        ASTCache:       true,
    }
}
```

## Risk Assessment

### Technical Risks
- **Memory Leaks:** Risk of memory leaks in template cache
- **Mitigation:** Proper cleanup routines and memory monitoring
- **Cache Inconsistency:** Risk of stale cached templates
- **Mitigation:** File modification detection and cache invalidation

### Performance Risks
- **Cache Thrashing:** Risk of excessive cache evictions
- **Mitigation:** Proper cache sizing and LRU eviction
- **Memory Pressure:** Risk of memory pressure from cache
- **Mitigation:** Memory limits and monitoring

### Functional Risks
- **Template Staleness:** Risk of serving stale template data
- **Mitigation:** File modification detection and cache invalidation
- **Cache Corruption:** Risk of cache data corruption
- **Mitigation:** Proper error handling and cache validation

## Rollback Plan

### Rollback Triggers
- Cache hit rate < 70%
- Memory usage > 100MB
- Cache error rate > 1%
- Performance degradation > 10%

### Rollback Procedure
1. **Immediate:** Disable cache via configuration
2. **Short-term:** Revert to non-cached template rendering
3. **Long-term:** Fix issues and re-enable cache

## Dependencies

### Internal Dependencies
- `internal/templates/manager.go` - Core template management
- `internal/types/templates/` - Template type definitions
- `internal/logging/` - Logging infrastructure

### External Dependencies
- Go sync package for thread safety
- Go crypto/sha256 for file checksums
- Go encoding/json for size calculation
- No additional external dependencies required

## Timeline

### Week 1: Implementation
- Day 1: Cache design and architecture
- Day 2-3: Core cache implementation
- Day 4-5: Compilation caching

### Week 2: Testing and Optimization
- Day 6-7: Integration and testing
- Day 8-10: Monitoring and optimization

## Conclusion

This improvement plan provides a systematic approach to implementing intelligent template caching, addressing the performance bottleneck of repeated template parsing and compilation while maintaining data consistency and memory efficiency.

**Expected Outcomes:**
- 90% cache hit rate for repeated template renders
- Reduced template parsing overhead
- Configurable memory management
- Enhanced user experience with faster template rendering

The plan ensures a robust caching implementation that integrates seamlessly with the existing template system while providing significant performance improvements.

**Key Benefits:**
- **Performance Improvement:** Faster template rendering through caching
- **Resource Efficiency:** Reduced CPU usage for template parsing
- **Scalability:** Better performance under high template usage
- **User Experience:** Faster response times for template operations
