# Facts API Reference

## Overview

The facts system in spooky provides comprehensive functionality for collecting, storing, validating, and managing system facts from remote machines. This document provides the complete API reference for the facts system.

**Status: Implemented** - The facts system provides comprehensive functionality for fact collection, storage, and management.

## Core Components

### Facts Manager (`internal/facts/manager.go`)

The facts manager provides high-level fact operations:

```go
type Manager struct {
    logger   spookylogging.Logger
    storage  spookystorage.FactStorage
    config   *spookytypes.Config
    cache    spookyinterfaces.CacheManager
}
```

**Implemented Operations:**
- ✅ **CollectFacts**: Collect facts from remote machines
- ✅ **StoreFacts**: Store facts in persistent storage
- ✅ **ValidateFacts**: Validate fact data and structure
- ✅ **GetFacts**: Retrieve facts from storage
- ✅ **UpdateFacts**: Update existing facts
- ✅ **DeleteFacts**: Delete facts from storage

### Facts Collector (`internal/facts/collector.go`)

The facts collector handles fact collection from various sources:

```go
type Collector struct {
    logger   spookylogging.Logger
    sshClient spookyssh.SSHClient
    config   *spookytypes.Config
}
```

**Supported Collection Methods:**
- ✅ **SSH Collection**: Collect facts via SSH from remote machines
- ✅ **Local Collection**: Collect facts from local system
- ✅ **HCL Collection**: Collect facts from HCL configuration files
- ✅ **Custom Collection**: Support for custom fact collection scripts

### Facts Storage (`internal/facts/storage.go`)

The facts storage system provides persistent fact storage:

```go
type FactStorage struct {
    db       *badger.DB
    logger   spookylogging.Logger
    config   *spookytypes.Config
}
```

**Storage Features:**
- ✅ **BadgerDB Storage**: High-performance key-value storage
- ✅ **JSON Storage**: Human-readable JSON storage format
- ✅ **Encryption**: Optional fact encryption for sensitive data
- ✅ **Compression**: Data compression for large fact collections
- ✅ **Indexing**: Efficient fact indexing and querying

## API Methods

### Fact Collection

#### CollectFacts
```go
func (m *Manager) CollectFacts(server string) (*spookytypes.FactCollection, error)
```
Collects facts from a remote server via SSH.

**Parameters:**
- `server`: Target server identifier

**Returns:**
- `FactCollection`: Collected facts with metadata
- `error`: Error if collection fails

#### CollectFactsViaSSH
```go
func (c *Collector) CollectFactsViaSSH(ctx context.Context, machine *spookytypes.Machine) (*spookytypes.FactCollection, error)
```
Collects facts from a remote machine via SSH.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `machine`: Target machine configuration

**Returns:**
- `FactCollection`: Collected facts with metadata
- `error`: Error if collection fails

#### CollectFactsFromHCL
```go
func (c *Collector) CollectFactsFromHCL(filePath string) (*spookytypes.FactCollection, error)
```
Collects facts from HCL configuration files.

**Parameters:**
- `filePath`: Path to HCL fact file

**Returns:**
- `FactCollection`: Collected facts from HCL
- `error`: Error if collection fails

### Fact Storage

#### StoreFacts
```go
func (m *Manager) StoreFacts(collection *spookytypes.FactCollection) error
```
Stores facts in persistent storage.

**Parameters:**
- `collection`: Fact collection to store

**Returns:**
- `error`: Error if storage fails

#### GetFacts
```go
func (m *Manager) GetFacts(server string) (*spookytypes.FactCollection, error)
```
Retrieves facts for a specific server.

**Parameters:**
- `server`: Server identifier

**Returns:**
- `FactCollection`: Retrieved facts
- `error`: Error if retrieval fails

#### UpdateFacts
```go
func (m *Manager) UpdateFacts(collection *spookytypes.FactCollection) error
```
Updates existing facts in storage.

**Parameters:**
- `collection`: Updated fact collection

**Returns:**
- `error`: Error if update fails

#### DeleteFacts
```go
func (m *Manager) DeleteFacts(server string) error
```
Deletes facts for a specific server.

**Parameters:**
- `server`: Server identifier

**Returns:**
- `error`: Error if deletion fails

### Fact Validation

#### ValidateFacts
```go
func (m *Manager) ValidateFacts(collection *spookytypes.FactCollection) error
```
Validates fact data and structure.

**Parameters:**
- `collection`: Fact collection to validate

**Returns:**
- `error`: Error if validation fails

#### ValidateFactSchema
```go
func (m *Manager) ValidateFactSchema(facts map[string]interface{}) error
```
Validates facts against defined schemas.

**Parameters:**
- `facts`: Facts to validate

**Returns:**
- `error`: Error if validation fails

### Fact Processing

#### ProcessFacts
```go
func (m *Manager) ProcessFacts(collection *spookytypes.FactCollection) (*spookytypes.ProcessedFacts, error)
```
Processes and transforms collected facts.

**Parameters:**
- `collection`: Raw fact collection

**Returns:**
- `ProcessedFacts`: Processed and transformed facts
- `error`: Error if processing fails

#### FilterFacts
```go
func (m *Manager) FilterFacts(collection *spookytypes.FactCollection, filters map[string]string) (*spookytypes.FactCollection, error)
```
Filters facts based on specified criteria.

**Parameters:**
- `collection`: Fact collection to filter
- `filters`: Filter criteria

**Returns:**
- `FactCollection`: Filtered facts
- `error`: Error if filtering fails

## Configuration

### Facts Configuration
```go
type FactsConfig struct {
    CollectionTimeout time.Duration `json:"collection_timeout" hcl:"collection_timeout"`
    StoragePath       string        `json:"storage_path" hcl:"storage_path"`
    StorageFormat     string        `json:"storage_format" hcl:"storage_format"`
    EncryptionEnabled bool          `json:"encryption_enabled" hcl:"encryption_enabled"`
    CompressionEnabled bool         `json:"compression_enabled" hcl:"compression_enabled"`
    MaxFactSize       int64         `json:"max_fact_size" hcl:"max_fact_size"`
    CacheEnabled      bool          `json:"cache_enabled" hcl:"cache_enabled"`
    CacheTTL          time.Duration `json:"cache_ttl" hcl:"cache_ttl"`
}
```

### Fact Collection Configuration
```go
type CollectionConfig struct {
    SSHEnabled        bool          `json:"ssh_enabled" hcl:"ssh_enabled"`
    LocalEnabled      bool          `json:"local_enabled" hcl:"local_enabled"`
    HCLEnabled        bool          `json:"hcl_enabled" hcl:"hcl_enabled"`
    CustomEnabled     bool          `json:"custom_enabled" hcl:"custom_enabled"`
    ParallelWorkers   int           `json:"parallel_workers" hcl:"parallel_workers"`
    RetryAttempts     int           `json:"retry_attempts" hcl:"retry_attempts"`
    RetryDelay        time.Duration `json:"retry_delay" hcl:"retry_delay"`
}
```

### Fact Storage Configuration
```go
type StorageConfig struct {
    Backend           string        `json:"backend" hcl:"backend"`
    Path              string        `json:"path" hcl:"path"`
    Format            string        `json:"format" hcl:"format"`
    EncryptionKey     string        `json:"encryption_key" hcl:"encryption_key"`
    CompressionLevel  int           `json:"compression_level" hcl:"compression_level"`
    MaxFileSize       int64         `json:"max_file_size" hcl:"max_file_size"`
    BackupEnabled     bool          `json:"backup_enabled" hcl:"backup_enabled"`
    BackupInterval    time.Duration `json:"backup_interval" hcl:"backup_interval"`
}
```

## Data Types

### Fact Collection
```go
type FactCollection struct {
    ID          string                 `json:"id" hcl:"id"`
    Server      string                 `json:"server" hcl:"server"`
    CollectedAt time.Time              `json:"collected_at" hcl:"collected_at"`
    Source      string                 `json:"source" hcl:"source"`
    Facts       map[string]interface{} `json:"facts" hcl:"facts"`
    Metadata    *FactMetadata          `json:"metadata" hcl:"metadata"`
    Version     string                 `json:"version" hcl:"version"`
}
```

### Fact Metadata
```go
type FactMetadata struct {
    CollectionMethod string            `json:"collection_method" hcl:"collection_method"`
    CollectionTime   time.Duration     `json:"collection_time" hcl:"collection_time"`
    FactCount        int               `json:"fact_count" hcl:"fact_count"`
    TotalSize        int64             `json:"total_size" hcl:"total_size"`
    Errors           []string          `json:"errors" hcl:"errors"`
    Warnings         []string          `json:"warnings" hcl:"warnings"`
    Tags             map[string]string `json:"tags" hcl:"tags"`
}
```

### Processed Facts
```go
type ProcessedFacts struct {
    CollectionID string                 `json:"collection_id" hcl:"collection_id"`
    ProcessedAt  time.Time              `json:"processed_at" hcl:"processed_at"`
    Facts        map[string]interface{} `json:"facts" hcl:"facts"`
    Transformations []string            `json:"transformations" hcl:"transformations"`
    Statistics   *FactStatistics        `json:"statistics" hcl:"statistics"`
}
```

### Fact Statistics
```go
type FactStatistics struct {
    TotalFacts    int     `json:"total_facts" hcl:"total_facts"`
    UniqueValues  int     `json:"unique_values" hcl:"unique_values"`
    NullValues    int     `json:"null_values" hcl:"null_values"`
    AverageSize   float64 `json:"average_size" hcl:"average_size"`
    LargestFact   string  `json:"largest_fact" hcl:"largest_fact"`
    SmallestFact  string  `json:"smallest_fact" hcl:"smallest_fact"`
}
```

## Error Handling

### Fact Collection Errors
```go
type CollectionError struct {
    Server      string    `json:"server" hcl:"server"`
    Method      string    `json:"method" hcl:"method"`
    Error       string    `json:"error" hcl:"error"`
    Timestamp   time.Time `json:"timestamp" hcl:"timestamp"`
    RetryCount  int       `json:"retry_count" hcl:"retry_count"`
    Recoverable bool      `json:"recoverable" hcl:"recoverable"`
}
```

### Fact Validation Errors
```go
type ValidationError struct {
    Field       string `json:"field" hcl:"field"`
    Value       string `json:"value" hcl:"value"`
    Rule        string `json:"rule" hcl:"rule"`
    Message     string `json:"message" hcl:"message"`
    Severity    string `json:"severity" hcl:"severity"`
}
```

### Fact Storage Errors
```go
type StorageError struct {
    Operation   string    `json:"operation" hcl:"operation"`
    Path        string    `json:"path" hcl:"path"`
    Error       string    `json:"error" hcl:"error"`
    Timestamp   time.Time `json:"timestamp" hcl:"timestamp"`
    Recoverable bool      `json:"recoverable" hcl:"recoverable"`
}
```

## Integration Points

### SSH Integration
The facts system integrates with the SSH system for remote fact collection:

```go
// Collect facts via SSH
facts, err := sshCollector.CollectViaSSH(ctx, machine)
```

### Storage Integration
The facts system integrates with the storage system for persistent fact storage:

```go
// Store facts in persistent storage
err = factStorage.StoreFacts(collection)
```

### Actions Integration
The facts system integrates with the actions system for automated fact collection:

```go
// Collect facts as part of action execution
err = actionManager.CollectFactsForAction(ctx, action)
```

## Usage Examples

### Basic Fact Collection
```go
// Create facts manager
manager := NewFactsManager(config, logger, storage)

// Collect facts from server
facts, err := manager.CollectFacts("web-server")
if err != nil {
    log.Fatal(err)
}

// Store facts
err = manager.StoreFacts(facts)
if err != nil {
    log.Fatal(err)
}
```

### SSH Fact Collection
```go
// Create collector
collector := NewFactsCollector(config, logger, sshClient)

// Collect facts via SSH
machine := &spookytypes.Machine{
    Hostname: "web.example.com",
    User:     "admin",
    Port:     22,
}

facts, err := collector.CollectFactsViaSSH(ctx, machine)
if err != nil {
    log.Fatal(err)
}
```

### Fact Validation
```go
// Validate collected facts
err = manager.ValidateFacts(facts)
if err != nil {
    log.Printf("Fact validation failed: %v", err)
}

// Validate against schema
err = manager.ValidateFactSchema(facts.Facts)
if err != nil {
    log.Printf("Schema validation failed: %v", err)
}
```

### Fact Processing
```go
// Process and transform facts
processed, err := manager.ProcessFacts(facts)
if err != nil {
    log.Fatal(err)
}

// Filter facts
filters := map[string]string{
    "environment": "production",
    "role":        "web",
}

filtered, err := manager.FilterFacts(facts, filters)
if err != nil {
    log.Fatal(err)
}
```

### Fact Storage Operations
```go
// Store facts
err = manager.StoreFacts(facts)

// Retrieve facts
retrieved, err := manager.GetFacts("web-server")

// Update facts
err = manager.UpdateFacts(updatedFacts)

// Delete facts
err = manager.DeleteFacts("web-server")
```

## CLI Commands

### Available Commands
- ✅ `spooky facts gather` - Collect facts from machines
- ✅ `spooky facts list` - List stored facts
- ✅ `spooky facts validate` - Validate fact data
- ✅ `spooky facts export` - Export facts to various formats
- ✅ `spooky facts import` - Import facts from external sources

### Command Examples
```bash
# Collect facts from all machines
spooky facts gather my-project

# Collect facts from specific machines
spooky facts gather my-project --machine web-server

# Collect facts in parallel
spooky facts gather my-project --parallel 4

# Validate stored facts
spooky facts validate my-project

# Export facts to JSON
spooky facts export my-project --format json --output facts.json

# Export facts to HCL
spooky facts export my-project --format hcl --output facts.hcl
```

## Performance Features

### Parallel Collection
```go
// Configure parallel fact collection
config := &spookytypes.FactsConfig{
    ParallelWorkers: 4,
    CollectionTimeout: 30 * time.Second,
    RetryAttempts: 3,
    RetryDelay: 5 * time.Second,
}
```

### Caching
```go
// Enable fact caching
config := &spookytypes.FactsConfig{
    CacheEnabled: true,
    CacheTTL: 1 * time.Hour,
}
```

### Compression
```go
// Enable fact compression
config := &spookytypes.FactsConfig{
    CompressionEnabled: true,
    CompressionLevel: 6,
}
```

## Security Features

### Fact Encryption
```go
// Enable fact encryption
config := &spookytypes.FactsConfig{
    EncryptionEnabled: true,
    EncryptionKey: "your-encryption-key",
}
```

### Access Control
```go
// Configure fact access control
config := &spookytypes.FactsConfig{
    AccessControl: &FactAccessControl{
        ReadUsers:  []string{"admin", "operator"},
        WriteUsers: []string{"admin"},
        ReadGroups: []string{"facts-readers"},
        WriteGroups: []string{"facts-writers"},
    },
}
```

## Testing

### Unit Tests
- ✅ **Manager Tests**: Test facts manager functionality
- ✅ **Collector Tests**: Test fact collection methods
- ✅ **Storage Tests**: Test fact storage operations
- ✅ **Validation Tests**: Test fact validation logic

### Integration Tests
- ✅ **SSH Collection Tests**: Test SSH-based fact collection
- ✅ **Storage Integration Tests**: Test storage system integration
- ✅ **CLI Tests**: Test CLI command functionality
- ✅ **Performance Tests**: Test fact collection performance

### Test Coverage
- ✅ **Collection Methods**: Test all collection methods
- ✅ **Storage Operations**: Test all storage operations
- ✅ **Validation Logic**: Test validation scenarios
- ✅ **Error Handling**: Test error conditions and recovery

## Troubleshooting

### Common Issues

#### Collection Failures
- **SSH Connection Issues**: Check SSH connectivity and authentication
- **Timeout Issues**: Increase collection timeout for slow systems
- **Permission Issues**: Verify user permissions on target machines

#### Storage Issues
- **Disk Space**: Ensure sufficient disk space for fact storage
- **Permission Issues**: Verify storage directory permissions
- **Corruption Issues**: Check for storage corruption and repair

#### Validation Issues
- **Schema Mismatches**: Update fact schemas to match collected data
- **Data Type Issues**: Ensure fact data types match expected types
- **Required Fields**: Verify all required fields are present

### Debug Commands
```bash
# Enable debug logging
export SPOOKY_LOG_LEVEL=debug

# Test fact collection
spooky facts gather my-project --verbose

# Validate facts with detailed output
spooky facts validate my-project --verbose

# Check storage status
spooky facts list my-project --verbose
```

## Future Enhancements

### Planned Features
- **Real-time Fact Collection**: Continuous fact collection and monitoring
- **Fact Streaming**: Stream facts for real-time processing
- **Advanced Analytics**: Fact analytics and trend analysis
- **Fact Dependencies**: Handle fact dependencies and relationships

### Performance Improvements
- **Distributed Collection**: Distribute fact collection across multiple nodes
- **Incremental Collection**: Collect only changed facts
- **Smart Caching**: Intelligent fact caching strategies
- **Compression Optimization**: Optimize fact compression algorithms

## Summary

The facts system in spooky provides comprehensive functionality for:

1. **Fact Collection**: Multiple collection methods including SSH, local, and HCL
2. **Fact Storage**: Persistent storage with encryption and compression
3. **Fact Validation**: Comprehensive validation and schema checking
4. **Fact Processing**: Fact transformation and filtering
5. **Performance Optimization**: Parallel collection, caching, and compression
6. **Security Features**: Fact encryption and access control
7. **CLI Integration**: Complete CLI command support
8. **Error Handling**: Comprehensive error handling and recovery
9. **Testing**: Complete test coverage for all functionality
10. **Documentation**: Comprehensive API documentation and examples

The facts system is production-ready and provides all necessary functionality for reliable fact collection, storage, and management in the spooky automation platform.
