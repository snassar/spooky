# Scalable Fact Storage System Implementation Plan

## Overview

This document outlines the implementation plan for adding a scalable fact storage system to Spooky using BadgerDB and JSON support, as proposed in [issue #56](https://github.com/snassar/spooky/issues/56).

## Background

Spooky needs a scalable system for storing and querying facts about servers (machine-id, OS version, hardware info, etc.) that can handle deployments ranging from a few servers to thousands. The system should provide efficient storage, fast queries, and easy data portability.

## Current State Analysis

### ✅ Already in Place
1. **Go 1.24+** - Project uses Go 1.24.3
2. **Existing Fact System** - Well-structured with types, collectors, and manager
3. **JSON Support** - Facts already have JSON tags
4. **Project Structure** - Clean internal organization

### Current Fact System Structure
```
internal/facts/
├── manager.go      # Fact collection coordination
├── types.go        # Fact data structures
├── ssh_collector.go # SSH-based fact collection
├── local_collector.go # Local system fact collection
├── hcl_collector.go # HCL configuration facts
├── opentofu_collector.go # OpenTofu state facts
├── base_collector.go # Shared collector functionality
└── stub_collector.go # Stub implementations
```

## Prerequisites

### 1. Add BadgerDB Dependency
```bash
go get github.com/dgraph-io/badger/v4
```

### 2. Create Storage Package Structure
```
internal/
├── storage/
│   ├── badger/
│   │   ├── store.go      # BadgerDB implementation
│   │   ├── store_test.go
│   │   └── types.go      # Storage-specific types
│   ├── interface.go      # Storage interface
│   └── factory.go        # Storage factory
```

### 3. Update Fact Manager Integration
- Modify `internal/facts/manager.go` to use storage interface
- Add storage configuration options
- Implement fact persistence methods

### 4. Configuration Updates
- Add storage configuration to CLI flags
- Update config parsing for storage options
- Add storage directory configuration

### 5. Testing Infrastructure
- Add integration tests for BadgerDB storage
- Create test data generators
- Add performance benchmarks

## Implementation Plan

### Phase 1: Foundation (Week 1)
**Goal**: Add BadgerDB dependency and create storage interface

#### Tasks:
1. **Add BadgerDB dependency**
   ```bash
   go get github.com/dgraph-io/badger/v4
   ```

2. **Create storage interface** (`internal/storage/interface.go`)
   ```go
   type FactStorage interface {
       GetServerFacts(serverID string) (*ServerFacts, error)
       SetServerFacts(serverID string, facts *ServerFacts) error
       QueryFacts(query FactQuery) ([]*ServerFacts, error)
       ExportToJSON(w io.Writer) error
       ImportFromJSON(r io.Reader) error
       Close() error
   }
   ```

3. **Create storage factory** (`internal/storage/factory.go`)
   ```go
   type StorageOptions struct {
       Type StorageType
       Path string
   }
   
   func NewFactStorage(opts StorageOptions) (FactStorage, error)
   ```

4. **Add storage types** (`internal/storage/types.go`)
   ```go
   type StorageType string
   const (
       StorageTypeBadger StorageType = "badger"
       StorageTypeJSON   StorageType = "json"
   )
   ```

### Phase 2: BadgerDB Implementation (Week 2)
**Goal**: Implement BadgerDB storage backend

#### Tasks:
1. **Create BadgerDB store** (`internal/storage/badger/store.go`)
   - Implement `BadgerFactStorage` struct
   - Add database initialization and configuration
   - Implement all interface methods

2. **Add BadgerDB types** (`internal/storage/badger/types.go`)
   - Define BadgerDB-specific data structures
   - Add query optimization types

3. **Implement core methods**:
   - `GetServerFacts()` - Key-based retrieval
   - `SetServerFacts()` - ACID transactions
   - `QueryFacts()` - Iterator-based queries
   - `ExportToJSON()` - Full database export
   - `ImportFromJSON()` - Bulk data import

### Phase 3: JSON Storage Implementation (Week 2)
**Goal**: Implement JSON storage backend for small deployments

#### Tasks:
1. **Create JSON store** (`internal/storage/json/store.go`)
   - Implement `JSONFactStorage` struct
   - Add file-based persistence
   - Implement thread-safe operations

2. **Implement core methods**:
   - File-based read/write operations
   - Memory caching for performance
   - JSON marshaling/unmarshaling

### Phase 4: Fact Manager Integration (Week 3)
**Goal**: Integrate storage with existing fact manager

#### Tasks:
1. **Update fact manager** (`internal/facts/manager.go`)
   - Add storage dependency injection
   - Implement fact persistence methods
   - Add storage configuration options

2. **Add storage methods**:
   ```go
   func (m *Manager) PersistFacts(server string, collection *FactCollection) error
   func (m *Manager) LoadPersistedFacts(server string) (*FactCollection, error)
   func (m *Manager) QueryPersistedFacts(query FactQuery) ([]*FactCollection, error)
   ```

3. **Update fact collection flow**:
   - Collect facts from sources
   - Persist to storage backend
   - Cache in memory for performance

### Phase 5: CLI Integration (Week 3)
**Goal**: Add storage configuration to CLI

#### Tasks:
1. **Add CLI flags** (`internal/cli/commands.go`)
   ```go
   var (
       storageType = flag.String("storage", "badger", "Storage type: badger or json")
       storagePath = flag.String("storage-path", "facts.db", "Storage path")
       exportFacts = flag.String("export-facts", "", "Export facts to JSON file")
       importFacts = flag.String("import-facts", "", "Import facts from JSON file")
   )
   ```

2. **Add facts commands** (`internal/cli/facts.go`)
   ```go
   var FactsCmd = &cobra.Command{
       Use:   "facts",
       Short: "Manage server facts",
   }
   
   var QueryFactsCmd = &cobra.Command{
       Use:   "query",
       Short: "Query server facts",
   }
   ```

3. **Update execute command**:
   - Initialize storage backend
   - Handle export/import operations
   - Integrate with fact collection

### Phase 6: Testing and Validation (Week 4)
**Goal**: Comprehensive testing and performance validation

#### Tasks:
1. **Unit tests**:
   - Storage interface tests
   - BadgerDB implementation tests
   - JSON implementation tests
   - Fact manager integration tests

2. **Integration tests**:
   - End-to-end fact collection and storage
   - Export/import functionality
   - Query performance tests

3. **Performance benchmarks**:
   - Storage write performance
   - Query performance
   - Memory usage analysis
   - Large-scale deployment simulation

4. **Migration tests**:
   - JSON to BadgerDB migration
   - Data integrity validation

### Phase 7: Documentation and Examples (Week 4)
**Goal**: Complete documentation and usage examples

#### Tasks:
1. **Update documentation**:
   - Storage configuration guide
   - Performance characteristics
   - Migration guide
   - Troubleshooting guide

2. **Add examples**:
   - Basic usage examples
   - Advanced query examples
   - Migration examples
   - Performance tuning examples

3. **Update CLI help**:
   - Storage command documentation
   - Configuration examples
   - Best practices

## Technical Specifications

### Storage Interface
```go
type FactStorage interface {
    GetMachineFacts(machineID string) (*MachineFacts, error)
    SetMachineFacts(machineID string, facts *MachineFacts) error
    QueryFacts(query FactQuery) ([]*MachineFacts, error)
    DeleteFacts(query FactQuery) (int, error)  // Delete facts matching query, returns count
    DeleteMachineFacts(machineID string) error   // Delete specific machine facts
    ExportToJSON(w io.Writer) error
    ImportFromJSON(r io.Reader) error
    ExportToJSONWithEncryption(w io.Writer, opts ExportOptions) error
    ImportFromJSONWithDecryption(r io.Reader, identityFile string) error
    Close() error
}

type MachineFacts struct {
    MachineID     string            `json:"machine_id"`     // Machine ID as UUID
    MachineName   string            `json:"machine_name"`   // Human-readable name from HCL
    ActionFile    string            `json:"action_file"`    // Source action file path
    ProjectName   string            `json:"project_name"`   // Project name (portable across systems)
    ProjectPath   string            `json:"project_path"`   // Absolute path (for reference)
    Hostname      string            `json:"hostname"`
    IPAddress     string            `json:"ip_address"`
    OS            string            `json:"os"`
    OSVersion     string            `json:"os_version"`
    CPU           CPUInfo           `json:"cpu"`
    Memory        MemoryInfo        `json:"memory"`
    SystemID      string            `json:"system_id"`      // Actual system ID (/etc/machine-id)
    Tags          map[string]string `json:"tags"`           // Team, environment, etc.
    CreatedAt     time.Time         `json:"created_at"`
    UpdatedAt     time.Time         `json:"updated_at"`
}

type FactQuery struct {
    MachineName   string            // Query by human-readable machine name
    ActionFile    string            // Query by action file
    Tags          map[string]string // Query by tags
    OS            string            // Query by OS
    Environment   string            // Query by environment tag
    Limit         int               // Limit results
}
```

### Machine ID as UUID Strategy

The storage system uses the actual machine ID (from `/etc/machine-id`) as the primary key for facts storage. This approach provides:

1. **Global Uniqueness**: Machine IDs are globally unique across all systems
2. **Automatic Deduplication**: Same physical machine = same UUID regardless of action file
3. **No Manual Management**: No need to generate or track UUIDs manually
4. **Real-world Mapping**: Direct correlation to actual hardware

### Machine ID Collision Detection and Resolution

**Problem**: When machines are cloned (VMs, containers, etc.), `/etc/machine-id` is often not updated, leading to duplicate machine IDs across different systems.

**Detection Strategy**: Compare multiple identifying factors to detect potential collisions:

```go
type CollisionDetection struct {
    MachineID     string            // Primary identifier
    Hostname      string            // Should be unique per machine
    IPAddresses   []string          // Network identifiers
    MACAddresses  []string          // Hardware identifiers
    ActionFile    string            // Source configuration
    MachineName   string            // Human-readable name
    Timestamp     time.Time         // When facts were collected
}

type CollisionResolution struct {
    Type          CollisionType     // Type of collision detected
    Confidence    float64           // Confidence level (0.0-1.0)
    Resolution    ResolutionAction  // How to handle the collision
    Evidence      []string          // Supporting evidence
}

type CollisionType string
const (
    CollisionTypeNone       CollisionType = "none"
    CollisionTypeHostname   CollisionType = "hostname_mismatch"
    CollisionTypeNetwork    CollisionType = "network_mismatch"
    CollisionTypeHardware   CollisionType = "hardware_mismatch"
    CollisionTypeMultiple   CollisionType = "multiple_mismatches"
)

type ResolutionAction string
const (
    ResolutionActionUpdate     ResolutionAction = "update"      // Update existing record
    ResolutionActionCreate     ResolutionAction = "create"      // Create new record with suffix
    ResolutionActionMerge      ResolutionAction = "merge"       // Merge facts from both sources
    ResolutionActionWarn       ResolutionAction = "warn"        // Warn user and ask for action
    ResolutionActionSkip       ResolutionAction = "skip"        // Skip this collection
)
```

#### Collision Detection Logic

```go
func DetectCollision(storage FactStorage, newFacts *MachineFacts) (*CollisionResolution, error) {
    // Check if machine ID already exists
    existing, err := storage.GetMachineFacts(newFacts.SystemID)
    if err != nil {
        // No existing record - no collision
        return &CollisionResolution{Type: CollisionTypeNone}, nil
    }
    
    // Machine ID exists - check for collision indicators
    var mismatches []string
    confidence := 0.0
    
    // Hostname mismatch (high confidence)
    if existing.Hostname != newFacts.Hostname {
        mismatches = append(mismatches, fmt.Sprintf("hostname: %s vs %s", existing.Hostname, newFacts.Hostname))
        confidence += 0.4
    }
    
    // IP address mismatch (medium confidence)
    if existing.IPAddress != newFacts.IPAddress {
        mismatches = append(mismatches, fmt.Sprintf("ip: %s vs %s", existing.IPAddress, newFacts.IPAddress))
        confidence += 0.3
    }
    
    // Action file mismatch (low confidence - could be same machine in different configs)
    if existing.ActionFile != newFacts.ActionFile {
        mismatches = append(mismatches, fmt.Sprintf("action_file: %s vs %s", existing.ActionFile, newFacts.ActionFile))
        confidence += 0.1
    }
    
    // Machine name mismatch (low confidence - could be renamed)
    if existing.MachineName != newFacts.MachineName {
        mismatches = append(mismatches, fmt.Sprintf("machine_name: %s vs %s", existing.MachineName, newFacts.MachineName))
        confidence += 0.1
    }
    
    // Determine collision type and resolution
    if len(mismatches) == 0 {
        return &CollisionResolution{Type: CollisionTypeNone}, nil
    }
    
    var collisionType CollisionType
    var resolution ResolutionAction
    
    switch {
    case confidence >= 0.7:
        collisionType = CollisionTypeMultiple
        resolution = ResolutionActionWarn
    case confidence >= 0.5:
        collisionType = CollisionTypeHostname
        resolution = ResolutionActionUpdate
    case confidence >= 0.3:
        collisionType = CollisionTypeNetwork
        resolution = ResolutionActionMerge
    default:
        collisionType = CollisionTypeNone
        resolution = ResolutionActionUpdate
    }
    
    return &CollisionResolution{
        Type:       collisionType,
        Confidence: confidence,
        Resolution: resolution,
        Evidence:   mismatches,
    }, nil
}
```

#### Collision Resolution Strategies

```go
func ResolveCollision(storage FactStorage, newFacts *MachineFacts, resolution *CollisionResolution) error {
    switch resolution.Resolution {
    case ResolutionActionUpdate:
        // Update existing record with new facts
        return storage.SetMachineFacts(newFacts.SystemID, newFacts)
        
    case ResolutionActionCreate:
        // Create new record with machine ID suffix
        newMachineID := newFacts.SystemID + "-" + generateSuffix(newFacts)
        newFacts.MachineID = newMachineID
        newFacts.SystemID = newMachineID
        return storage.SetMachineFacts(newMachineID, newFacts)
        
    case ResolutionActionMerge:
        // Merge facts from both sources
        existing, _ := storage.GetMachineFacts(newFacts.SystemID)
        merged := mergeMachineFacts(existing, newFacts)
        return storage.SetMachineFacts(newFacts.SystemID, merged)
        
    case ResolutionActionWarn:
        // Log warning and ask for user input
        logCollisionWarning(newFacts, resolution)
        return fmt.Errorf("machine ID collision detected: %v", resolution.Evidence)
        
    case ResolutionActionSkip:
        // Skip this collection
        return fmt.Errorf("skipping fact collection due to collision")
        
    default:
        return fmt.Errorf("unknown resolution action: %s", resolution.Resolution)
    }
}

func mergeServerFacts(existing, new *ServerFacts) *ServerFacts {
    merged := &ServerFacts{
        ServerID:   existing.ServerID,
        MachineID:  existing.MachineID,
        CreatedAt:  existing.CreatedAt,
        UpdatedAt:  time.Now(),
        Tags:       make(map[string]string),
    }
    
    // Merge tags
    for k, v := range existing.Tags {
        merged.Tags[k] = v
    }
    for k, v := range new.Tags {
        merged.Tags[k] = v
    }
    
    // Prefer newer facts for most fields
    if new.UpdatedAt.After(existing.UpdatedAt) {
        merged.Hostname = new.Hostname
        merged.IPAddress = new.IPAddress
        merged.OS = new.OS
        merged.OSVersion = new.OSVersion
        merged.CPU = new.CPU
        merged.Memory = new.Memory
    } else {
        merged.Hostname = existing.Hostname
        merged.IPAddress = existing.IPAddress
        merged.OS = existing.OS
        merged.OSVersion = existing.OSVersion
        merged.CPU = existing.CPU
        merged.Memory = existing.Memory
    }
    
    // Keep both server names and action files
    if existing.ServerName != new.ServerName {
        merged.ServerName = existing.ServerName + "," + new.ServerName
    } else {
        merged.ServerName = existing.ServerName
    }
    
    if existing.ActionFile != new.ActionFile {
        merged.ActionFile = existing.ActionFile + "," + new.ActionFile
    } else {
        merged.ActionFile = existing.ActionFile
    }
    
    return merged
}
```

#### Implementation Details

```go
// Server ID generation strategy
func GenerateServerID(facts *FactCollection) string {
    // Use machine_id fact if available
    if machineID, exists := facts.Facts["machine_id"]; exists {
        if id, ok := machineID.Value.(string); ok && id != "" {
            return id
        }
    }
    
    // Fallback: generate UUID from hostname + IP + action file
    return generateUUIDFromFacts(facts)
}

// HCL configuration support
server "web-001" {
    host = "192.168.1.10"
    use_machine_id = true  // Use actual machine-id as UUID (default)
    tags = {
        team = "web-team"
        environment = "production"
    }
}
```

#### Data Segmentation and Storage Location

**XDG Base Directory Compliance**: Following the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/0.8/), facts are stored in `$XDG_DATA_HOME/spooky/` (defaults to `$HOME/.local/share/spooky/`).

**Default Storage Structure**:
```
$HOME/.local/share/spooky/
├── facts.db/                    # BadgerDB database directory
│   ├── 000001.vlog
│   ├── 000002.sst
│   └── MANIFEST
├── facts.json                   # JSON storage file (if using JSON backend)
├── config/                      # Spooky configuration
│   └── storage.conf
└── logs/                        # Storage operation logs
    └── facts.log
```

**BadgerDB Key-Value Structure**:
```
├── "6538b193d562410ab480b1ce22469fe6" → {
│     server_id: "6538b193d562410ab480b1ce22469fe6",
│     server_name: "web-001",
│     action_file: "actions/deploy-web.hcl",
│     project_name: "web-app",                    # Project name (portable)
│     project_path: "/home/user/projects/web-app", # Absolute path (for reference)
│     hostname: "web1.example.com",
│     machine_id: "6538b193d562410ab480b1ce22469fe6",
│     tags: {"team": "web-team", "environment": "production"}
│   }
├── "a1b2c3d4e5f6789012345678901234567" → {
│     server_id: "a1b2c3d4e5f6789012345678901234567", 
│     server_name: "db-001",
│     action_file: "actions/setup-db.hcl",
│     project_name: "database",                    # Project name (portable)
│     project_path: "/home/user/projects/database", # Absolute path (for reference)
│     hostname: "db1.example.com",
│     machine_id: "a1b2c3d4e5f6789012345678901234567",
│     tags: {"team": "db-team", "environment": "production"}
│   }
└── "f9e8d7c6b5a4321098765432109876543" → {
      server_id: "f9e8d7c6b5a4321098765432109876543",
      server_name: "web-001",  // Same name, different machine
      action_file: "actions/deploy-staging.hcl",
      project_name: "staging",                     # Project name (portable)
      project_path: "/home/user/my-little-projects/staging", # Absolute path (for reference)
      hostname: "staging-web1.example.com", 
      machine_id: "f9e8d7c6b5a4321098765432109876543",
      tags: {"team": "web-team", "environment": "staging"}
    }
```

#### Query Examples

```go
// Query by machine ID (exact match)
storage.GetServerFacts("6538b193d562410ab480b1ce22469fe6")

// Query by server name (across all action files)
query := FactQuery{ServerName: "web-001"}

// Query by action file
query := FactQuery{ActionFile: "actions/deploy-web.hcl"}

// Query by team
query := FactQuery{Tags: map[string]string{"team": "web-team"}}

// Query by environment
query := FactQuery{Tags: map[string]string{"environment": "production"}}
```

#### Collision Detection Examples

**Example 1: Cloned VM with same machine ID**
```
WARNING: Machine ID collision detected for 6538b193d562410ab480b1ce22469fe6 (confidence: 0.80)
  - hostname: web1.example.com vs web1-staging.example.com
  - ip: 192.168.1.10 vs 192.168.2.10
  - action_file: actions/deploy-web.hcl vs actions/deploy-staging.hcl
```

**Example 2: Same machine in different configurations**
```
WARNING: Machine ID collision detected for 6538b193d562410ab480b1ce22469fe6 (confidence: 0.20)
  - server_name: web-001 vs web-server
  - action_file: actions/deploy-web.hcl vs actions/setup-webserver.hcl
```

**Example 3: Network change on same machine**
```
WARNING: Machine ID collision detected for 6538b193d562410ab480b1ce22469fe6 (confidence: 0.30)
  - ip: 192.168.1.10 vs 192.168.1.15
```

#### Best Practices for Collision Handling

1. **Development/Testing Environments**: Use `--collision-policy=update` for rapid iteration
2. **Production Environments**: Use `--collision-policy=warn` to catch potential issues
3. **Strict Environments**: Use `--collision-policy=skip` to prevent data corruption
4. **Multi-Config Environments**: Use `--collision-policy=merge` to combine facts from different action files

#### HCL Configuration for Storage and Collision Handling

```hcl
# Global storage configuration
storage {
    type = "badger"                    # badger or json
    path = "/custom/path/facts.db"     # Override XDG default
    collision_policy = "warn"          # update, warn, skip, merge
}

# Global collision policy (if not in storage block)
collision_policy = "warn"

# Per-server override
server "web-001" {
    host = "192.168.1.10"
    collision_policy = "merge"  # Override global policy
    tags = {
        team = "web-team"
        environment = "production"
    }
}

# Force machine ID regeneration for cloned VMs
server "cloned-vm" {
    host = "192.168.1.20"
    regenerate_machine_id = true  # Force new machine ID generation
    tags = {
        team = "testing"
        environment = "staging"
    }
}
```

#### Portable Project Names

**Problem**: When sharing `facts.db` between users, absolute project paths become invalid:
- Your path: `/home/user/projects/web-app`
- Colleague's path: `/home/user/my-little-projects/web-app`

**Solution**: Store both project name and path:
- `project_name`: `"web-app"` (portable across systems)
- `project_path`: `"/home/user/projects/web-app"` (for reference)

**Usage**:
```bash
# Query by portable project name (works across systems)
spooky facts list --project-name "web-app"

# Query by absolute path (system-specific)
spooky facts list --project-path "/home/user/projects/web-app"

# Delete by portable project name
spooky facts delete --project-name "web-app"
```

#### Selective Fact Deletion Examples

**Delete facts from a specific action file:**
```bash
# Delete all facts collected from a specific action file
spooky facts delete --action-file "actions/deploy-web.hcl"

# Preview what will be deleted first
spooky facts list --action-file "actions/deploy-web.hcl"
```

**Delete facts from a specific project:**
```bash
# Delete facts by portable project name (recommended)
spooky facts delete --project-name "web-app"

# Delete facts by absolute project path (system-specific)
spooky facts delete --project-path "/home/user/projects/web-app"

# Delete facts from a specific environment
spooky facts delete --environment "staging"
```

**Delete facts for a specific team:**
```bash
# Delete all facts for a specific team
spooky facts delete --team "web-team"

# Delete facts for a specific server
spooky facts delete --server-name "web-001"
```

**Combined deletion criteria:**
```bash
# Delete facts from staging environment for web-team
spooky facts delete --environment "staging" --team "web-team"

# Delete facts from a specific action file in production
spooky facts delete --action-file "actions/deploy-db.hcl" --environment "production"
```

**Force deletion without confirmation:**
```bash
# Delete without prompting (useful for scripts)
spooky facts delete --action-file "actions/test-config.hcl" --confirm
```

**Configuration Precedence (highest to lowest):**
1. **CLI flags** (--storage-path, --storage) - Explicit user choice
2. **Project configuration** (project.conf.hcl in current directory) - Project-specific settings
3. **Environment variables** (SPOOKY_FACTS_PATH, SPOOKY_FACTS_FORMAT) - System-wide defaults
4. **XDG user config** ($HOME/.config/spooky/config.hcl) - User defaults
5. **XDG defaults** ($HOME/.local/share/spooky/) - Built-in fallback

**Storage location examples:**
```bash
# Use default XDG location ($HOME/.local/share/spooky/)
spooky facts list

# Use custom storage location via CLI
spooky facts list --storage-path "/custom/path/facts.db"

# Use environment variable
export SPOOKY_FACTS_PATH="/shared/facts.db"
spooky facts list

# Use JSON storage instead of BadgerDB
spooky facts list --storage json --storage-path "/tmp/facts.json"

# Use environment variable for storage format
export SPOOKY_FACTS_FORMAT="json"
spooky facts list --storage-path "/tmp/facts.json"

# Use badgerdb format (default)
export SPOOKY_FACTS_FORMAT="badgerdb"
spooky facts list
```

**Configuration hierarchy examples:**

**Project-specific configuration** (`project.conf.hcl` in current directory):
```hcl
# $HOME/Work/important-project/project.conf.hcl
storage {
    type = "json"
    path = ".facts.json"  # Project-specific storage
    collision_policy = "update"
}
```

**User configuration** (`$HOME/.config/spooky/config.hcl`):
```hcl
# $HOME/.config/spooky/config.hcl
storage {
    type = "badgerdb"
    path = "/var/lib/spooky/facts.db"  # User's preferred location
    collision_policy = "warn"
}
```

**Environment variables:**
```bash
export SPOOKY_FACTS_PATH="/shared/facts.db"
export SPOOKY_FACTS_FORMAT="badgerdb"
```

**Precedence example:**
```bash
# In $HOME/Work/important-project/
# 1. CLI flag wins
spooky facts list --storage-path "/tmp/override.db"

# 2. Project config wins over env/user config
spooky facts list  # Uses .facts.json from project.conf.hcl

# 3. Environment variable wins over user config
# (if no project config exists)

# 4. User config wins over XDG defaults
# (if no env vars set)

# 5. XDG defaults as fallback
# ($HOME/.local/share/spooky/facts.db)
```

#### Centralized Fact Gathering and Distribution

**Large Deployment Scenario**: Multiple people working on large fleets with centralized fact collection.

**Architecture:**
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Fleet Servers │    │ Central Spooky  │    │ Object Storage  │
│                 │    │   Container/VM  │    │                 │
│ • web-001       │◄──►│ • Runs every    │───►│ • facts-2024-   │
│ • web-002       │    │   15 minutes    │    │   01-15.json    │
│ • db-001        │    │ • Collects all  │    │ • facts-2024-   │
│ • db-002        │    │   fleet facts   │    │   01-15-12.json │
│ • ...           │    │ • Exports to    │    │ • ...           │
└─────────────────┘    │   JSON          │    └─────────────────┘
                       └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │ Team Members    │
                       │                 │
                       │ • Alice imports │
                       │   latest facts  │
                       │ • Bob starts    │
                       │   with fresh    │
                       │   data          │
                       └─────────────────┘
```

**Step 1: Centralized Fact Collection**
```bash
# Central container/VM runs scheduled fact collection
#!/bin/bash
# /opt/spooky/collect-facts.sh

# Collect facts from all servers in fleet
spooky facts gather --machine "web-001" --action-file "fleet/deploy-web.hcl"
spooky facts gather --machine "web-002" --action-file "fleet/deploy-web.hcl"
spooky facts gather --machine "db-001" --action-file "fleet/deploy-db.hcl"
# ... (automated for all servers)

# Export directly to S3/MinIO (no local file needed)
TIMESTAMP=$(date +"%Y-%m-%d-%H-%M")
spooky facts export --backend s3 --s3-bucket "spooky-facts" --s3-prefix "fleet" --output "facts-${TIMESTAMP}.json"

# Or export to HTTP endpoint
spooky facts export --backend http --http-url "https://facts-api.company.com/export" --output "facts-${TIMESTAMP}.json"
```

**Step 2: Team Member Bootstrap**
```bash
# New team member or returning colleague
# Import directly from S3/MinIO
spooky facts import --backend s3 --s3-bucket "spooky-facts" --s3-prefix "fleet" --input "facts-2024-01-15-14-30.json"

# Or import from HTTP endpoint
spooky facts import --backend http --http-url "https://facts-api.company.com/facts/latest" --input "latest-facts.json"

# Verify import
spooky facts list --project-name "fleet" | head -20
```

**Step 3: Automated Setup Script**
```bash
#!/bin/bash
# bootstrap-facts.sh

echo "Bootstrapping spooky facts from centralized collection..."

# Import latest facts directly from S3
spooky facts import --backend s3 --s3-bucket "spooky-facts" --s3-prefix "fleet" --input "latest"

# Or import from HTTP endpoint
# spooky facts import --backend http --http-url "https://facts-api.company.com/facts/latest" --input "latest"

echo "Import complete. Found $(spooky facts list --project-name fleet | wc -l) fleet servers"
```

**Step 4: Cron Job Setup**
```bash
# /etc/cron.d/spooky-facts
# Run fact collection every 15 minutes
*/15 * * * * root /opt/spooky/collect-facts.sh >> /var/log/spooky-facts.log 2>&1
```

#### Sharing Facts Databases

**Sharing with colleagues:**
```bash
# Export facts to JSON (portable format)
spooky facts export --output shared-facts.json

# Share the JSON file with colleague
# Colleague imports the facts
spooky facts import --input shared-facts.json
```

**Sharing BadgerDB directly:**
```bash
# Copy the entire facts.db directory
cp -r $HOME/.local/share/spooky/facts.db/ /path/to/share/

# Colleague uses the copied database
spooky facts list --storage-path "/path/to/share/facts.db"
```

**Best practices for sharing:**
1. **Use project names**: Query by `--project-name` instead of `--project-path`
2. **Export to JSON**: More portable than BadgerDB files
3. **Include tags**: Use team/environment tags for organization
4. **Document structure**: Note the project names and action files included
5. **Centralized collection**: Use scheduled jobs for large fleets
6. **Object storage**: Store facts in S3/GCS for team access
7. **Timestamped exports**: Keep historical fact snapshots

### Security Considerations

**Sensitive Information in Facts:**
```go
// High Risk - Always sanitize before export
type SensitiveFacts struct {
    Environment map[string]string `json:"environment"` // API keys, passwords, tokens
    NetworkIPs  []string          `json:"network_ips"` // Internal IP addresses
    NetworkMACs []string          `json:"network_macs"` // MAC addresses
    DNS         DNSInfo           `json:"dns"`         // Internal DNS servers
    MachineID   string            `json:"machine_id"`  // Unique system identifier
}

// Medium Risk - Consider sanitizing
type PotentiallySensitiveFacts struct {
    Hostname    string `json:"hostname"`    // Internal hostnames
    FQDN        string `json:"fqdn"`        // Internal domain names
    ProjectPath string `json:"project_path"` // Internal file paths
    ActionFile  string `json:"action_file"`  // Configuration file paths
}

// Low Risk - Generally safe
type SafeFacts struct {
    OS          OSInfo     `json:"os"`          // OS version, distribution
    CPU         CPUInfo    `json:"cpu"`         // CPU model, cores
    Memory      MemoryInfo `json:"memory"`      // Memory capacity
    Tags        map[string]string `json:"tags"` // Team, environment labels
}
```

**Export Encryption Options:**
```go
type ExportOptions struct {
    // Sanitization options
    SanitizeEnvironment bool     // Remove environment variables
    SanitizeNetwork     bool     // Remove IP/MAC addresses
    SanitizePaths       bool     // Remove file paths
    SanitizeHostnames   bool     // Remove hostnames
    SanitizeMachineID   bool     // Remove machine IDs
    RedactPatterns      []string // Custom redaction patterns
    
    // Encryption options
    EncryptFile         bool     // Encrypt entire JSON file with age
    AgePublicKey        string   // age public key for encryption
    EncryptFields       bool     // Encrypt sensitive fields individually
    EncryptedFields     []string // Fields to encrypt (environment, network_ips, etc.)
    AgeRecipient        string   // age recipient for field encryption
}

type EncryptedFact struct {
    MachineID     string            `json:"machine_id"`
    MachineName   string            `json:"machine_name"`
    ActionFile    string            `json:"action_file"`
    ProjectName   string            `json:"project_name"`
    ProjectPath   string            `json:"project_path"`
    Hostname      string            `json:"hostname,omitempty"`
    IPAddress     string            `json:"ip_address,omitempty"`
    OS            string            `json:"os"`
    OSVersion     string            `json:"os_version"`
    CPU           CPUInfo           `json:"cpu"`
    Memory        MemoryInfo        `json:"memory"`
    SystemID      string            `json:"system_id,omitempty"`
    Tags          map[string]string `json:"tags"`
    CreatedAt     time.Time         `json:"created_at"`
    UpdatedAt     time.Time         `json:"updated_at"`
    
    // Encrypted fields
    EncryptedEnvironment string `json:"encrypted_environment,omitempty"`
    EncryptedNetwork     string `json:"encrypted_network,omitempty"`
    EncryptedDNS         string `json:"encrypted_dns,omitempty"`
}

func (s *FactStorage) ExportToJSONWithEncryption(w io.Writer, opts ExportOptions) error {
    // Get all facts
    facts, err := s.QueryFacts(FactQuery{})
    if err != nil {
        return err
    }
    
    // Process facts for encryption
    var processedFacts []interface{}
    for _, fact := range facts {
        if opts.EncryptFields {
            processedFacts = append(processedFacts, s.encryptSensitiveFields(fact, opts))
        } else {
            processedFacts = append(processedFacts, s.sanitizeFactsForExport(fact, opts))
        }
    }
    
    // Marshal to JSON
    jsonData, err := json.MarshalIndent(processedFacts, "", "  ")
    if err != nil {
        return err
    }
    
    // Encrypt entire file if requested
    if opts.EncryptFile {
        return s.encryptFileWithAge(w, jsonData, opts.AgePublicKey)
    }
    
    // Write plain JSON
    _, err = w.Write(jsonData)
    return err
}

func (s *FactStorage) encryptSensitiveFields(facts *MachineFacts, opts ExportOptions) *EncryptedFact {
    encrypted := &EncryptedFact{
        MachineID:   facts.MachineID,
        MachineName: facts.MachineName,
        ActionFile:  facts.ActionFile,
        ProjectName: facts.ProjectName,
        ProjectPath: facts.ProjectPath,
        OS:          facts.OS,
        OSVersion:   facts.OSVersion,
        CPU:         facts.CPU,
        Memory:      facts.Memory,
        Tags:        facts.Tags,
        CreatedAt:   facts.CreatedAt,
        UpdatedAt:   facts.UpdatedAt,
    }
    
    // Conditionally include sensitive fields (encrypted)
    if !opts.SanitizeHostnames {
        encrypted.Hostname = facts.Hostname
    }
    if !opts.SanitizeMachineID {
        encrypted.SystemID = facts.SystemID
    }
    
    // Encrypt sensitive data
    if !opts.SanitizeEnvironment && facts.Environment != nil {
        if encryptedData, err := s.encryptFieldWithAge(facts.Environment, opts.AgeRecipient); err == nil {
            encrypted.EncryptedEnvironment = encryptedData
        }
    }
    
    if !opts.SanitizeNetwork && facts.NetworkIPs != nil {
        if encryptedData, err := s.encryptFieldWithAge(facts.NetworkIPs, opts.AgeRecipient); err == nil {
            encrypted.EncryptedNetwork = encryptedData
        }
    }
    
    return encrypted
}

func (s *FactStorage) encryptFieldWithAge(data interface{}, recipient string) (string, error) {
    // Marshal data to JSON
    jsonData, err := json.Marshal(data)
    if err != nil {
        return "", err
    }
    
    // Encrypt with age
    return s.encryptWithAge(jsonData, recipient)
}

func (s *FactStorage) encryptFileWithAge(w io.Writer, data []byte, publicKey string) error {
    // Encrypt entire file with age
    encryptedData, err := s.encryptWithAge(data, publicKey)
    if err != nil {
        return err
    }
    
    // Write encrypted data
    _, err = w.Write([]byte(encryptedData))
    return err
}

// age encryption implementation
func (s *FactStorage) encryptWithAge(data []byte, recipient string) (string, error) {
    // Parse recipient
    rec, err := age.ParseX25519Recipient(recipient)
    if err != nil {
        return "", fmt.Errorf("failed to parse age recipient: %w", err)
    }
    
    // Create encrypted writer
    var buf bytes.Buffer
    encryptedWriter, err := age.Encrypt(&buf, rec)
    if err != nil {
        return "", fmt.Errorf("failed to create age encryptor: %w", err)
    }
    
    // Write data
    if _, err := encryptedWriter.Write(data); err != nil {
        return "", fmt.Errorf("failed to write encrypted data: %w", err)
    }
    
    // Close writer
    if err := encryptedWriter.Close(); err != nil {
        return "", fmt.Errorf("failed to close age encryptor: %w", err)
    }
    
    return buf.String(), nil
}

func (s *FactStorage) decryptWithAge(encryptedData string, identityFile string) ([]byte, error) {
    // Read identity file
    identity, err := age.ParseIdentitiesFile(identityFile)
    if err != nil {
        return nil, fmt.Errorf("failed to parse age identity: %w", err)
    }
    
    // Create decrypted reader
    decryptedReader, err := age.Decrypt(strings.NewReader(encryptedData), identity...)
    if err != nil {
        return nil, fmt.Errorf("failed to create age decryptor: %w", err)
    }
    
    // Read decrypted data
    decryptedData, err := io.ReadAll(decryptedReader)
    if err != nil {
        return nil, fmt.Errorf("failed to read decrypted data: %w", err)
    }
    
    return decryptedData, nil
}

func (s *FactStorage) ImportFromJSONWithDecryption(r io.Reader, identityFile string) error {
    // Read all data from reader
    data, err := io.ReadAll(r)
    if err != nil {
        return fmt.Errorf("failed to read input data: %w", err)
    }
    
    var jsonData []byte
    
    // Check if data is age-encrypted (starts with age-1)
    if strings.HasPrefix(string(data), "age-1") {
        if identityFile == "" {
            return fmt.Errorf("age identity file required for decryption")
        }
        
        // Decrypt the data
        decryptedData, err := s.decryptWithAge(string(data), identityFile)
        if err != nil {
            return fmt.Errorf("failed to decrypt data: %w", err)
        }
        jsonData = decryptedData
    } else {
        // Data is not encrypted, use as-is
        jsonData = data
    }
    
    // Parse JSON data
    var facts []*MachineFacts
    if err := json.Unmarshal(jsonData, &facts); err != nil {
        return fmt.Errorf("failed to parse JSON: %w", err)
    }
    
    // Import each fact
    for _, fact := range facts {
        if err := s.SetMachineFacts(fact.MachineID, fact); err != nil {
            return fmt.Errorf("failed to import fact for %s: %w", fact.MachineID, err)
        }
    }
    
    return nil
}

func (s *FactStorage) decryptEncryptedFields(facts []*EncryptedFact, identityFile string) ([]*MachineFacts, error) {
    var decryptedFacts []*MachineFacts
    
    for _, encryptedFact := range facts {
        fact := &MachineFacts{
            MachineID:   encryptedFact.MachineID,
            MachineName: encryptedFact.MachineName,
            ActionFile:  encryptedFact.ActionFile,
            ProjectName: encryptedFact.ProjectName,
            ProjectPath: encryptedFact.ProjectPath,
            Hostname:    encryptedFact.Hostname,
            IPAddress:   encryptedFact.IPAddress,
            OS:          encryptedFact.OS,
            OSVersion:   encryptedFact.OSVersion,
            CPU:         encryptedFact.CPU,
            Memory:      encryptedFact.Memory,
            SystemID:    encryptedFact.SystemID,
            Tags:        encryptedFact.Tags,
            CreatedAt:   encryptedFact.CreatedAt,
            UpdatedAt:   encryptedFact.UpdatedAt,
        }
        
        // Decrypt environment variables if present
        if encryptedFact.EncryptedEnvironment != "" {
            decryptedEnv, err := s.decryptWithAge(encryptedFact.EncryptedEnvironment, identityFile)
            if err != nil {
                return nil, fmt.Errorf("failed to decrypt environment: %w", err)
            }
            
            var envVars []EnvironmentVariable
            if err := json.Unmarshal(decryptedEnv, &envVars); err != nil {
                return nil, fmt.Errorf("failed to parse decrypted environment: %w", err)
            }
            
            // Convert to map for storage
            fact.Environment = make(map[string]string)
            for _, env := range envVars {
                fact.Environment[env.Key] = env.Value
            }
        }
        
        // Decrypt network information if present
        if encryptedFact.EncryptedNetwork != "" {
            decryptedNetwork, err := s.decryptWithAge(encryptedFact.EncryptedNetwork, identityFile)
            if err != nil {
                return nil, fmt.Errorf("failed to decrypt network: %w", err)
            }
            
            var networkInfo NetworkInfo
            if err := json.Unmarshal(decryptedNetwork, &networkInfo); err != nil {
                return nil, fmt.Errorf("failed to parse decrypted network: %w", err)
            }
            
            fact.NetworkIPs = networkInfo.IPs
            fact.NetworkMACs = networkInfo.MACs
        }
        
        decryptedFacts = append(decryptedFacts, fact)
    }
    
    return decryptedFacts, nil
}
```

**CLI Export with Encryption:**
```bash
# Export with field-level encryption (sensitive fields encrypted)
spooky facts export --encrypt-fields --age-recipient "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p" --output "facts-encrypted.json"

# Export with file-level encryption (entire file encrypted)
spooky facts export --encrypt-file --age-public-key "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p" --output "facts.age"

# Export with both field and file encryption (maximum security)
spooky facts export --encrypt-fields --encrypt-file --age-recipient "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p" --age-public-key "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p" --output "facts-secure.age"

# Export with sanitization (no sensitive data)
spooky facts export --sanitize-env --sanitize-network --output "safe-facts.json"

# Export to stdout (default behavior)
spooky facts export --encrypt-fields --age-recipient "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"

# Decrypt and process encrypted facts
age -d -i ~/.age/identity.txt facts.age | jq '.[] | {machine_id, server_name}'

# Decrypt field-encrypted facts and extract environment variables
spooky facts export --encrypt-fields --age-recipient "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p" | \
  jq '.[] | .encrypted_environment' | \
  age -d -i ~/.age/identity.txt | \
  jq -r '.[] | select(.key | startswith("AWS_")) | "\(.key)=\(.value)"'

**CLI Import with Decryption:**
```bash
# Import encrypted facts file
spooky facts import --input "facts.age" --age-identity "~/.age/identity.txt"

# Import field-encrypted facts
spooky facts import --input "facts-encrypted.json" --age-identity "~/.age/identity.txt"

# Import from stdin (pipe from age decryption)
age -d -i ~/.age/identity.txt facts.age | spooky facts import

# Import with automatic decryption detection
spooky facts import --input "facts.age" --age-identity "~/.age/identity.txt"
```

**CLI Usage with Global Flags:**
```bash
# Use global configuration
spooky facts gather --config-dir "/etc/spooky" --ssh-key-path "~/.ssh/id_rsa" --log-level "debug"

# Enable verbose output with dry-run
spooky facts gather --verbose --dry-run --machine "web-001" --action-file "fleet/deploy-web.hcl"

# Suppress output for scripting
spooky facts export --quiet --output "facts.json"

# Use custom log file
spooky facts validate --log-file "/var/log/spooky-facts.log"
```

**CLI Validate Command Examples:**
```bash
# Validate all facts in database
spooky facts validate

# Validate with verbose output
spooky facts validate --verbose

# Validate specific facts by criteria
spooky facts validate --machine "web-001" --project "fleet"
```

**CLI Query Command Examples:**
```bash
# Query by OS
spooky facts query "os=ubuntu"

# Query by hostname pattern
spooky facts query "hostname=web-*"

# Query by project and limit results
spooky facts query "machine=web-001,limit=10"

# Query by tags
spooky facts query "tag=environment:production,tag=team:web-team"

# Complex query with multiple criteria
spooky facts query "os=ubuntu,machine=web-001,tag=environment:production,limit=5"
```

**Complete Encryption/Decryption Workflow:**
```bash
# 1. Export with field-level encryption
spooky facts export --encrypt-fields --age-recipient "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p" --output "facts-encrypted.json"

# 2. Import field-encrypted facts
spooky facts import --input "facts-encrypted.json" --age-identity "~/.age/identity.txt"

# 3. Export with file-level encryption
spooky facts export --encrypt-file --age-public-key "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p" --output "facts.age"

# 4. Import file-encrypted facts
spooky facts import --input "facts.age" --age-identity "~/.age/identity.txt"

# 5. Export with both field and file encryption (maximum security)
spooky facts export --encrypt-fields --encrypt-file --age-recipient "..." --age-public-key "..." --output "facts-secure.age"

# 6. Import with automatic detection (works with any encryption level)
spooky facts import --input "facts-secure.age" --age-identity "~/.age/identity.txt"
```
```

**Best Practices:**
1. **Always encrypt sensitive data** - Use age encryption for environment variables, network info, etc.
2. **Use field-level encryption** - Encrypt individual sensitive fields while keeping structure readable
3. **Use file-level encryption** - Encrypt entire exports for maximum security
4. **Never store unencrypted sensitive data** - Always use encryption or sanitization
5. **Use project names** instead of absolute paths for portability
6. **Consider machine ID sensitivity** in your environment
7. **Document what's being exported** for compliance
8. **Use HTTPS for HTTP storage** - PUT operations require HTTPS for security
9. **Validate TLS certificates** - Don't skip certificate verification in production
10. **Secure key management** - Store age keys securely and rotate regularly

### Storage Backends

**Storage Backend Types:**
```go
type StorageBackend string

const (
    BackendBadgerDB StorageBackend = "badgerdb"
    BackendJSON     StorageBackend = "json"
    BackendS3       StorageBackend = "s3"
    BackendMinIO    StorageBackend = "minio"
    BackendHTTP     StorageBackend = "http"
)
```

**Storage Configuration:**
```go
type StorageConfig struct {
    Backend     StorageBackend `json:"backend"`
    Path        string         `json:"path"`        // For local storage (BadgerDB, JSON)
    S3Config    *S3Config      `json:"s3_config"`   // For S3/MinIO
    HTTPConfig  *HTTPConfig    `json:"http_config"` // For HTTP endpoints
}

type S3Config struct {
    Region          string `json:"region"`
    Bucket          string `json:"bucket"`
    Prefix          string `json:"prefix"`
    AccessKeyID     string `json:"access_key_id"`
    SecretAccessKey string `json:"secret_access_key"`
    Endpoint        string `json:"endpoint"`        // For MinIO compatibility
    UseSSL          bool   `json:"use_ssl"`
}

type HTTPConfig struct {
    BaseURL     string            `json:"base_url"`
    Headers     map[string]string `json:"headers"`
    Timeout     time.Duration     `json:"timeout"`
    RetryCount  int               `json:"retry_count"`
}
```

**S3/MinIO Implementation:**
```go
type S3Storage struct {
    client *s3.Client
    bucket string
    prefix string
}

func NewS3Storage(config *S3Config) (*S3Storage, error) {
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        return nil, err
    }
    
    client := s3.NewFromConfig(cfg)
    return &S3Storage{
        client: client,
        bucket: config.Bucket,
        prefix: config.Prefix,
    }, nil
}

func (s *S3Storage) StoreFacts(serverID string, facts *ServerFacts) error {
    data, err := json.Marshal(facts)
    if err != nil {
        return err
    }
    
    key := fmt.Sprintf("%s/%s.json", s.prefix, serverID)
    _, err = s.client.PutObject(context.TODO(), &s3.PutObjectInput{
        Bucket: aws.String(s.bucket),
        Key:    aws.String(key),
        Body:   bytes.NewReader(data),
    })
    return err
}

func (s *S3Storage) GetFacts(serverID string) (*ServerFacts, error) {
    key := fmt.Sprintf("%s/%s.json", s.prefix, serverID)
    result, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
        Bucket: aws.String(s.bucket),
        Key:    aws.String(key),
    })
    if err != nil {
        return nil, err
    }
    defer result.Body.Close()
    
    var facts ServerFacts
    if err := json.NewDecoder(result.Body).Decode(&facts); err != nil {
        return nil, err
    }
    return &facts, nil
}

func (s *S3Storage) ExportToJSON(w io.Writer) error {
    // List all objects in bucket/prefix
    // Download and combine into single JSON
    return nil
}
```

**HTTP Storage Implementation:**
```go
type HTTPStorage struct {
    baseURL string
    client  *http.Client
    headers map[string]string
}

func NewHTTPStorage(config *HTTPConfig) (*HTTPStorage, error) {
    // Validate HTTPS for write operations
    if !strings.HasPrefix(config.BaseURL, "https://") {
        return nil, fmt.Errorf("HTTPS is required for HTTP storage backend: %s", config.BaseURL)
    }
    
    return &HTTPStorage{
        baseURL: config.BaseURL,
        client:  &http.Client{Timeout: config.Timeout},
        headers: config.Headers,
    }, nil
}

func (h *HTTPStorage) StoreFacts(serverID string, facts *ServerFacts) error {
    // Double-check HTTPS for write operations
    if !strings.HasPrefix(h.baseURL, "https://") {
        return fmt.Errorf("HTTPS is required for storing facts: %s", h.baseURL)
    }
    
    data, err := json.Marshal(facts)
    if err != nil {
        return err
    }
    
    url := fmt.Sprintf("%s/facts/%s", h.baseURL, serverID)
    req, err := http.NewRequest("PUT", url, bytes.NewReader(data))
    if err != nil {
        return err
    }
    
    req.Header.Set("Content-Type", "application/json")
    for k, v := range h.headers {
        req.Header.Set(k, v)
    }
    
    resp, err := h.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("HTTP error: %d", resp.StatusCode)
    }
    return nil
}

func (h *HTTPStorage) GetFacts(serverID string) (*ServerFacts, error) {
    // Warn about HTTP for read operations (but allow it)
    if !strings.HasPrefix(h.baseURL, "https://") {
        fmt.Fprintf(os.Stderr, "Warning: Using HTTP for fact retrieval. Consider using HTTPS for better security.\n")
    }
    
    url := fmt.Sprintf("%s/facts/%s", h.baseURL, serverID)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    for k, v := range h.headers {
        req.Header.Set(k, v)
    }
    
    resp, err := h.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
    }
    
    var facts ServerFacts
    if err := json.NewDecoder(resp.Body).Decode(&facts); err != nil {
        return nil, err
    }
    return &facts, nil
}
```

### BadgerDB Configuration
```go
type BadgerOptions struct {
    Path              string
    ValueLogFileSize  int64
    NumVersionsToKeep int
    MaxTableSize      int64
    LevelOneSize      int64
    ValueThreshold    int
}
```

### Performance Targets
- **Write Performance**: 10,000+ facts/second
- **Read Performance**: 100,000+ facts/second
- **Query Performance**: Sub-second response for 1,000+ servers
- **Memory Usage**: < 100MB for 1,000 servers
- **Storage Size**: < 1MB per server

## Success Criteria

1. **BadgerDB default**: Successfully use BadgerDB for all deployments
2. **JSON option**: Successfully use JSON for small deployments
3. **Performance**: Handle 10,000+ servers efficiently
4. **Portability**: Export/import facts between storage types
5. **CLI integration**: Seamless integration with existing spooky CLI
6. **Template integration**: Facts available in templates
7. **Migration**: Automatic migration between storage types
8. **Reliability**: ACID transactions and data integrity

## Risk Mitigation

### Technical Risks
1. **BadgerDB complexity**: Start with simple implementation, add features incrementally
2. **Performance issues**: Comprehensive benchmarking and optimization
3. **Data corruption**: Implement backup and recovery procedures
4. **Memory usage**: Monitor and optimize memory consumption

### Integration Risks
1. **Breaking changes**: Maintain backward compatibility
2. **Configuration complexity**: Provide sensible defaults
3. **Migration issues**: Thorough testing of migration paths

## Out of Scope

- **External databases**: No PostgreSQL, MySQL, etc.
- **Distributed storage**: No clustering or replication
- **Real-time sync**: No live synchronization between storage types
- **Custom protocols**: No custom fact collection protocols
- **External APIs**: No external fact collection services

## Timeline

- **Week 1**: Foundation and storage interface
- **Week 2**: BadgerDB and JSON implementations
- **Week 3**: Integration and CLI updates
- **Week 4**: Testing, documentation, and validation

## Development Instructions

### For AI Assistants and Developers

This section provides specific, actionable instructions for implementing the fact storage system.

#### Phase 1: Foundation Setup

**Step 1: Add Storage Dependencies**
```bash
# Add BadgerDB to go.mod
go get github.com/dgraph-io/badger/v4

# Add AWS SDK for S3/MinIO support
go get github.com/aws/aws-sdk-go-v2
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/service/s3

# Add age encryption support
go get filippo.io/age

# Verify dependencies were added
go mod tidy
go mod verify
```

**Step 2: Create Storage Package Structure**
```bash
# Create directory structure (integrated within internal/facts)
mkdir -p internal/facts

# Create initial files
touch internal/facts/storage.go
touch internal/facts/badger_storage.go
touch internal/facts/badger_storage_test.go
touch internal/facts/json_storage.go
touch internal/facts/json_storage_test.go
touch internal/facts/collision.go
touch internal/facts/collision_test.go
```

**Step 3: Implement Storage Interface**
```go
// File: internal/facts/storage.go
package facts

import (
    "io"
    "time"
    "path/filepath"
)

type FactStorage interface {
    GetMachineFacts(machineID string) (*MachineFacts, error)
    SetMachineFacts(machineID string, facts *MachineFacts) error
    QueryFacts(query FactQuery) ([]*MachineFacts, error)
    DeleteFacts(query FactQuery) (int, error)  // Delete facts matching query, returns count
    DeleteMachineFacts(machineID string) error   // Delete specific machine facts
    ExportToJSON(w io.Writer) error
    ImportFromJSON(r io.Reader) error
    ExportToJSONWithEncryption(w io.Writer, opts ExportOptions) error
    ImportFromJSONWithDecryption(r io.Reader, identityFile string) error
    Close() error
}

type MachineFacts struct {
    MachineID     string            `json:"machine_id"`     // Machine ID as UUID
    MachineName   string            `json:"machine_name"`   // Human-readable name from HCL
    ActionFile    string            `json:"action_file"`    // Source action file path
    ProjectName   string            `json:"project_name"`   // Project name (portable across systems)
    ProjectPath   string            `json:"project_path"`   // Absolute path (for reference)
    Hostname      string            `json:"hostname"`
    IPAddress     string            `json:"ip_address"`
    OS            string            `json:"os"`
    OSVersion     string            `json:"os_version"`
    CPU           CPUInfo           `json:"cpu"`
    Memory        MemoryInfo        `json:"memory"`
    SystemID      string            `json:"system_id"`      // Actual system ID (/etc/machine-id)
    Tags          map[string]string `json:"tags"`           // Team, environment, etc.
    CreatedAt     time.Time         `json:"created_at"`
    UpdatedAt     time.Time         `json:"updated_at"`
}

type FactQuery struct {
    MachineName     string            // Query by human-readable machine name
    ActionFile      string            // Query by action file
    ProjectName     string            // Query by project name (portable)
    ProjectPath     string            // Query by absolute project path
    Tags            map[string]string // Query by tags
    OS              string            // Query by OS
    Environment     string            // Query by environment tag
    Limit           int               // Limit results
    SearchQuery     string            // Text search query (supports regex)
    SearchField     string            // Field to search in
    UpdatedBefore   *time.Time        // Filter by update time
    UpdatedAfter    *time.Time        // Filter by update time
}

type BulkUpdateOperation struct {
    MachineID   string            `json:"machine_id"`
    Description string            `json:"description"`
    Updates     map[string]interface{} `json:"updates"`
    AddTags     map[string]string `json:"add_tags,omitempty"`
    RemoveTags  []string          `json:"remove_tags,omitempty"`
}

type DatabaseStats struct {
    TotalFacts           int       `json:"total_facts"`
    TotalServers         int       `json:"total_servers"`
    DatabaseSize         int64     `json:"database_size"`
    OldestFact           time.Time `json:"oldest_fact"`
    NewestFact           time.Time `json:"newest_fact"`
    AverageFactsPerServer float64  `json:"average_facts_per_server"`
    TopTags              []TagCount `json:"top_tags"`
    TopOS                []OSCount  `json:"top_os"`
}

type TagCount struct {
    Key   string `json:"key"`
    Count int    `json:"count"`
}

type OSCount struct {
    Name  string `json:"name"`
    Count int    `json:"count"`
}

type CPUInfo struct {
    Cores int    `json:"cores"`
    Model string `json:"model"`
    Arch  string `json:"arch"`
}

type MemoryInfo struct {
    Total     uint64 `json:"total"`
    Available uint64 `json:"available"`
    Used      uint64 `json:"used"`
}

// Machine ID generation strategy
func GenerateMachineID(facts *FactCollection) string {
    // Use machine_id fact if available
    if machineID, exists := facts.Facts["machine_id"]; exists {
        if id, ok := machineID.Value.(string); ok && id != "" {
            return id
        }
    }
    
    // Fallback: generate UUID from hostname + IP + action file
    return generateUUIDFromFacts(facts)
}

// Extract project name from path for portability
func ExtractProjectName(projectPath string) string {
    if projectPath == "" {
        return ""
    }
    
    // Get the last component of the path
    return filepath.Base(projectPath)
}
```

**Step 4: Implement Storage Factory with XDG Support**
```go
// File: internal/facts/storage.go (continued)
package facts

import (
    "fmt"
    "os"
    "path/filepath"
)

type StorageType string

const (
    StorageTypeBadger StorageType = "badger"
    StorageTypeJSON   StorageType = "json"
)

type StorageOptions struct {
    Type StorageType
    Path string
}

// GetXDGDataHome returns the XDG data home directory for spooky
func GetXDGDataHome() string {
    if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
        return filepath.Join(xdgDataHome, "spooky")
    }
    
    // Default to $HOME/.local/share/spooky
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return ".spooky" // Fallback to current directory
    }
    
    return filepath.Join(homeDir, ".local", "share", "spooky")
}

// EnsureSpookyDataDir creates the spooky data directory if it doesn't exist
func EnsureSpookyDataDir() error {
    dataDir := GetXDGDataHome()
    return os.MkdirAll(dataDir, 0755)
}

func NewFactStorage(opts StorageOptions) (FactStorage, error) {
    // If no path specified, use XDG data home
    if opts.Path == "" {
        if err := EnsureSpookyDataDir(); err != nil {
            return nil, fmt.Errorf("failed to create spooky data directory: %w", err)
        }
        
        switch opts.Type {
        case StorageTypeJSON:
            opts.Path = filepath.Join(GetXDGDataHome(), "facts.json")
        case StorageTypeBadger, "": // Default to BadgerDB
            opts.Path = filepath.Join(GetXDGDataHome(), "facts.db")
        }
    }
    
    switch opts.Type {
    case StorageTypeJSON:
        return NewJSONFactStorage(opts.Path)
    case StorageTypeBadger, "": // Default to BadgerDB
        return NewBadgerFactStorage(opts.Path)
    default:
        return nil, fmt.Errorf("unsupported storage type: %s", opts.Type)
    }
}

// Get storage configuration with proper precedence hierarchy
func getStorageConfig() StorageOptions {
    opts := StorageOptions{}
    
    // 1. CLI flags (highest priority)
    if *storageType != "" {
        opts.Type = StorageType(*storageType)
    }
    if *storagePath != "" {
        opts.Path = *storagePath
    }
    
    // 2. Project configuration (project.conf.hcl in current directory)
    if opts.Type == "" || opts.Path == "" {
        projectConfig := loadProjectConfig("project.conf.hcl")
        if opts.Type == "" && projectConfig.Storage.Type != "" {
            opts.Type = projectConfig.Storage.Type
        }
        if opts.Path == "" && projectConfig.Storage.Path != "" {
            opts.Path = projectConfig.Storage.Path
        }
    }
    
    // 3. Environment variables
    if opts.Type == "" {
        if envFormat := os.Getenv(EnvSpookyFactsFormat); envFormat != "" {
            switch envFormat {
            case "badgerdb":
                opts.Type = StorageTypeBadger
            case "json":
                opts.Type = StorageTypeJSON
            default:
                opts.Type = StorageTypeBadger
            }
        }
    }
    if opts.Path == "" {
        if envPath := os.Getenv(EnvSpookyFactsPath); envPath != "" {
            opts.Path = envPath
        }
    }
    
    // 4. XDG user config ($HOME/.config/spooky/config.hcl)
    if opts.Type == "" || opts.Path == "" {
        userConfig := loadUserConfig()
        if opts.Type == "" && userConfig.Storage.Type != "" {
            opts.Type = userConfig.Storage.Type
        }
        if opts.Path == "" && userConfig.Storage.Path != "" {
            opts.Path = userConfig.Storage.Path
        }
    }
    
    // 5. XDG defaults (lowest priority)
    if opts.Type == "" {
        opts.Type = StorageTypeBadger // Default to badgerdb
    }
    // Path defaults to XDG_DATA_HOME/spooky/ (handled in NewFactStorage)
    
    return opts
}
```
```

#### Phase 2: BadgerDB Implementation

**Step 1: Implement BadgerDB Store**
```go
// File: internal/facts/badger_storage.go
package facts

import (
    "encoding/json"
    "fmt"
    "io"
    "time"

    "github.com/dgraph-io/badger/v4"
)

type BadgerFactStorage struct {
    db *badger.DB
}

func NewBadgerFactStorage(dbPath string) (*BadgerFactStorage, error) {
    opts := badger.DefaultOptions(dbPath)
    opts.Logger = nil // Disable logging for cleaner output
    
    db, err := badger.Open(opts)
    if err != nil {
        return nil, fmt.Errorf("failed to open BadgerDB: %w", err)
    }
    
    return &BadgerFactStorage{db: db}, nil
}

func (b *BadgerFactStorage) GetMachineFacts(machineID string) (*MachineFacts, error) {
    var facts MachineFacts
    
    err := b.db.View(func(txn *badger.Txn) error {
        item, err := txn.Get([]byte(machineID))
        if err != nil {
            return err
        }
        
        return item.Value(func(val []byte) error {
            return json.Unmarshal(val, &facts)
        })
    })
    
    if err == badger.ErrKeyNotFound {
        return nil, fmt.Errorf("server facts not found: %s", serverID)
    }
    
    return &facts, err
}

func (b *BadgerFactStorage) SetServerFacts(serverID string, facts *ServerFacts) error {
    facts.UpdatedAt = time.Now()
    if facts.CreatedAt.IsZero() {
        facts.CreatedAt = facts.UpdatedAt
    }
    
    data, err := json.Marshal(facts)
    if err != nil {
        return fmt.Errorf("failed to marshal facts: %w", err)
    }
    
    return b.db.Update(func(txn *badger.Txn) error {
        return txn.Set([]byte(serverID), data)
    })
}

func (b *BadgerFactStorage) QueryFacts(query storage.FactQuery) ([]*storage.ServerFacts, error) {
    var results []*storage.ServerFacts
    
    err := b.db.View(func(txn *badger.Txn) error {
        opts := badger.DefaultIteratorOptions
        opts.PrefetchSize = 100
        
        it := txn.NewIterator(opts)
        defer it.Close()
        
        for it.Rewind(); it.Valid(); it.Next() {
            item := it.Item()
            var facts storage.ServerFacts
            
            err := item.Value(func(val []byte) error {
                return json.Unmarshal(val, &facts)
            })
            if err != nil {
                continue
            }
            
            if matchesQuery(&facts, query) {
                results = append(results, &facts)
                if query.Limit > 0 && len(results) >= query.Limit {
                    break
                }
            }
        }
        
        return nil
    })
    
    return results, err
}

func (b *BadgerFactStorage) DeleteFacts(query storage.FactQuery) (int, error) {
    var deletedCount int
    
    err := b.db.Update(func(txn *badger.Txn) error {
        opts := badger.DefaultIteratorOptions
        opts.PrefetchSize = 100
        
        it := txn.NewIterator(opts)
        defer it.Close()
        
        for it.Rewind(); it.Valid(); it.Next() {
            item := it.Item()
            var facts storage.ServerFacts
            
            err := item.Value(func(val []byte) error {
                return json.Unmarshal(val, &facts)
            })
            if err != nil {
                continue
            }
            
            if matchesQuery(&facts, query) {
                if err := txn.Delete(item.Key()); err != nil {
                    return fmt.Errorf("failed to delete facts for %s: %w", string(item.Key()), err)
                }
                deletedCount++
            }
        }
        
        return nil
    })
    
    return deletedCount, err
}

func (b *BadgerFactStorage) DeleteServerFacts(serverID string) error {
    return b.db.Update(func(txn *badger.Txn) error {
        return txn.Delete([]byte(serverID))
    })
}

func (b *BadgerFactStorage) ExportToJSON(w io.Writer) error {
    facts := make(map[string]*storage.ServerFacts)
    
    err := b.db.View(func(txn *badger.Txn) error {
        opts := badger.DefaultIteratorOptions
        it := txn.NewIterator(opts)
        defer it.Close()
        
        for it.Rewind(); it.Valid(); it.Next() {
            item := it.Item()
            var serverFacts storage.ServerFacts
            
            err := item.Value(func(val []byte) error {
                return json.Unmarshal(val, &serverFacts)
            })
            if err != nil {
                continue
            }
            
            facts[string(item.KeyCopy(nil))] = &serverFacts
        }
        
        return nil
    })
    
    if err != nil {
        return err
    }
    
    encoder := json.NewEncoder(w)
    encoder.SetIndent("", "  ")
    return encoder.Encode(facts)
}

func (b *BadgerFactStorage) ImportFromJSON(r io.Reader) error {
    var facts map[string]*storage.ServerFacts
    
    decoder := json.NewDecoder(r)
    if err := decoder.Decode(&facts); err != nil {
        return fmt.Errorf("failed to decode JSON: %w", err)
    }
    
    return b.db.Update(func(txn *badger.Txn) error {
        for serverID, serverFacts := range facts {
            data, err := json.Marshal(serverFacts)
            if err != nil {
                return fmt.Errorf("failed to marshal facts for %s: %w", serverID, err)
            }
            
            if err := txn.Set([]byte(serverID), data); err != nil {
                return fmt.Errorf("failed to set facts for %s: %w", serverID, err)
            }
        }
        return nil
    })
}

func (b *BadgerFactStorage) Close() error {
    return b.db.Close()
}

func matchesQuery(facts *storage.ServerFacts, query storage.FactQuery) bool {
    // Implement query matching logic
    if query.OS != "" && facts.OS != query.OS {
        return false
    }
    
    if query.Environment != "" {
        if env, exists := facts.Tags["environment"]; !exists || env != query.Environment {
            return false
        }
    }
    
    for key, value := range query.Tags {
        if tagValue, exists := facts.Tags[key]; !exists || tagValue != value {
            return false
        }
    }
    
    return true
}
```

#### Phase 3: Collision Detection Implementation

**Step 1: Implement Collision Detection**
```go
// File: internal/facts/collision.go
package facts

import (
    "fmt"
    "time"
)

// CollisionDetection represents collision detection logic
type CollisionDetection struct {
    MachineID     string            // Primary identifier
    Hostname      string            // Should be unique per machine
    IPAddresses   []string          // Network identifiers
    MACAddresses  []string          // Hardware identifiers
    ActionFile    string            // Source configuration
    ServerName    string            // Human-readable name
    Timestamp     time.Time         // When facts were collected
}

type CollisionResolution struct {
    Type          CollisionType     // Type of collision detected
    Confidence    float64           // Confidence level (0.0-1.0)
    Resolution    ResolutionAction  // How to handle the collision
    Evidence      []string          // Supporting evidence
}

type CollisionType string
const (
    CollisionTypeNone       CollisionType = "none"
    CollisionTypeHostname   CollisionType = "hostname_mismatch"
    CollisionTypeNetwork    CollisionType = "network_mismatch"
    CollisionTypeHardware   CollisionType = "hardware_mismatch"
    CollisionTypeMultiple   CollisionType = "multiple_mismatches"
)

type ResolutionAction string
const (
    ResolutionActionUpdate     ResolutionAction = "update"      // Update existing record
    ResolutionActionCreate     ResolutionAction = "create"      // Create new record with suffix
    ResolutionActionMerge      ResolutionAction = "merge"       // Merge facts from both sources
    ResolutionActionWarn       ResolutionAction = "warn"        // Warn user and ask for action
    ResolutionActionSkip       ResolutionAction = "skip"        // Skip this collection
)

// DetectCollision detects machine ID collisions
func DetectCollision(storage FactStorage, newFacts *ServerFacts) (*CollisionResolution, error) {
    // Implementation from the collision detection section above
    // ... (full implementation)
}

// ResolveCollision handles collision resolution
func ResolveCollision(storage FactStorage, newFacts *ServerFacts, resolution *CollisionResolution) error {
    // Implementation from the collision resolution section above
    // ... (full implementation)
}
```

#### Phase 4: JSON Implementation

**Step 1: Implement JSON Store**
```go
// File: internal/facts/json_storage.go
package facts

import (
    "encoding/json"
    "fmt"
    "io"
    "os"
    "sync"
    "time"

    "spooky/internal/storage"
)

type JSONFactStorage struct {
    filepath string
    facts    map[string]*storage.ServerFacts
    mu       sync.RWMutex
}

func NewJSONFactStorage(filepath string) (*JSONFactStorage, error) {
    storage := &JSONFactStorage{
        filepath: filepath,
        facts:    make(map[string]*storage.ServerFacts),
    }
    
    // Load existing data if file exists
    if err := storage.load(); err != nil && !os.IsNotExist(err) {
        return nil, fmt.Errorf("failed to load existing facts: %w", err)
    }
    
    return storage, nil
}

func (j *JSONFactStorage) GetServerFacts(serverID string) (*storage.ServerFacts, error) {
    j.mu.RLock()
    defer j.mu.RUnlock()
    
    if facts, exists := j.facts[serverID]; exists {
        return facts, nil
    }
    
    return nil, fmt.Errorf("server facts not found: %s", serverID)
}

func (j *JSONFactStorage) SetServerFacts(serverID string, facts *storage.ServerFacts) error {
    j.mu.Lock()
    defer j.mu.Unlock()
    
    facts.UpdatedAt = time.Now()
    if facts.CreatedAt.IsZero() {
        facts.CreatedAt = facts.UpdatedAt
    }
    
    j.facts[serverID] = facts
    
    return j.save()
}

func (j *JSONFactStorage) QueryFacts(query storage.FactQuery) ([]*storage.ServerFacts, error) {
    j.mu.RLock()
    defer j.mu.RUnlock()
    
    var results []*storage.ServerFacts
    
    for _, facts := range j.facts {
        if matchesQuery(facts, query) {
            results = append(results, facts)
            if query.Limit > 0 && len(results) >= query.Limit {
                break
            }
        }
    }
    
    return results, nil
}

func (j *JSONFactStorage) ExportToJSON(w io.Writer) error {
    j.mu.RLock()
    defer j.mu.RUnlock()
    
    encoder := json.NewEncoder(w)
    encoder.SetIndent("", "  ")
    return encoder.Encode(j.facts)
}

func (j *JSONFactStorage) ImportFromJSON(r io.Reader) error {
    j.mu.Lock()
    defer j.mu.Unlock()
    
    var facts map[string]*storage.ServerFacts
    
    decoder := json.NewDecoder(r)
    if err := decoder.Decode(&facts); err != nil {
        return fmt.Errorf("failed to decode JSON: %w", err)
    }
    
    j.facts = facts
    return j.save()
}

func (j *JSONFactStorage) Close() error {
    return j.save()
}

func (j *JSONFactStorage) load() error {
    file, err := os.Open(j.filepath)
    if err != nil {
        return err
    }
    defer file.Close()
    
    decoder := json.NewDecoder(file)
    return decoder.Decode(&j.facts)
}

func (j *JSONFactStorage) save() error {
    file, err := os.Create(j.filepath)
    if err != nil {
        return fmt.Errorf("failed to create file: %w", err)
    }
    defer file.Close()
    
    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    return encoder.Encode(j.facts)
}

func matchesQuery(facts *storage.ServerFacts, query storage.FactQuery) bool {
    // Same implementation as BadgerDB
    if query.OS != "" && facts.OS != query.OS {
        return false
    }
    
    if query.Environment != "" {
        if env, exists := facts.Tags["environment"]; !exists || env != query.Environment {
            return false
        }
    }
    
    for key, value := range query.Tags {
        if tagValue, exists := facts.Tags[key]; !exists || tagValue != value {
            return false
        }
    }
    
    return true
}
```

#### Phase 4: Integration with Fact Manager

**Step 1: Update Fact Manager**
```go
// Add to internal/facts/manager.go
type Manager struct {
    // ... existing fields ...
    storage storage.FactStorage
}

func NewManager(sshClient *ssh.SSHClient, storage storage.FactStorage) *Manager {
    return &Manager{
        // ... existing initialization ...
        storage: storage,
    }
}

func (m *Manager) PersistFacts(server string, collection *FactCollection) error {
    if m.storage == nil {
        return nil // No storage configured
    }
    
    // Generate server ID using machine ID strategy
    serverID := GenerateServerID(collection)
    
    // Convert FactCollection to ServerFacts
    serverFacts := &storage.ServerFacts{
        ServerID:   serverID,                    // Machine ID as UUID
        ServerName: server,                      // Human-readable name
        ActionFile: collection.ActionFile,       // Source action file
        Hostname:   collection.Server,
        CreatedAt:  collection.Timestamp,
        UpdatedAt:  time.Now(),
        Tags:       make(map[string]string),
    }
    
    // Extract facts and convert to storage format
    for key, fact := range collection.Facts {
        switch key {
        case "hostname":
            if str, ok := fact.Value.(string); ok {
                serverFacts.Hostname = str
            }
        case "machine_id":
            if str, ok := fact.Value.(string); ok {
                serverFacts.MachineID = str
            }
        case "os.name":
            if str, ok := fact.Value.(string); ok {
                serverFacts.OS = str
            }
        case "os.version":
            if str, ok := fact.Value.(string); ok {
                serverFacts.OSVersion = str
            }
        case "network.ips":
            if ips, ok := fact.Value.([]string); ok && len(ips) > 0 {
                serverFacts.IPAddress = ips[0]
            }
        // Add more fact mappings as needed
        }
    }
    
    // Detect and resolve collisions before storing
    resolution, err := DetectCollision(m.storage, serverFacts)
    if err != nil {
        return fmt.Errorf("failed to detect collision: %w", err)
    }
    
    if resolution.Type != CollisionTypeNone {
        // Log collision detection
        logCollisionDetection(serverFacts, resolution)
        
        // Handle collision based on configuration
        switch m.collisionPolicy {
        case CollisionPolicyWarn:
            if resolution.Confidence >= 0.7 {
                return fmt.Errorf("high-confidence machine ID collision detected: %v", resolution.Evidence)
            }
            // Fall through to update for lower confidence
        case CollisionPolicySkip:
            if resolution.Confidence >= 0.5 {
                return fmt.Errorf("skipping fact collection due to collision: %v", resolution.Evidence)
            }
            // Fall through to update for lower confidence
        case CollisionPolicyMerge:
            return ResolveCollision(m.storage, serverFacts, resolution)
        }
    }
    
    return m.storage.SetServerFacts(serverID, serverFacts)
}

// CollisionPolicy defines how to handle machine ID collisions
type CollisionPolicy string
const (
    CollisionPolicyUpdate CollisionPolicy = "update"  // Always update (default)
    CollisionPolicyWarn   CollisionPolicy = "warn"    // Warn on high confidence collisions
    CollisionPolicySkip   CollisionPolicy = "skip"    // Skip on medium+ confidence collisions
    CollisionPolicyMerge  CollisionPolicy = "merge"   // Always merge facts
)

func logCollisionDetection(facts *ServerFacts, resolution *CollisionResolution) {
    fmt.Printf("WARNING: Machine ID collision detected for %s (confidence: %.2f)\n", 
        facts.MachineID, resolution.Confidence)
    for _, evidence := range resolution.Evidence {
        fmt.Printf("  - %s\n", evidence)
    }
}
```

#### Phase 5: CLI Integration

**Step 1: Add CLI Flags and Commands**
```go
// Add to internal/cli/commands.go
var (
    storageType = flag.String("storage", "badger", "Storage type: badger or json")
    storagePath = flag.String("storage-path", "", "Storage path (defaults to XDG_DATA_HOME/spooky/)")
    exportFacts = flag.String("export-facts", "", "Export facts to JSON file")
    importFacts = flag.String("import-facts", "", "Import facts from JSON file")
    collisionPolicy = flag.String("collision-policy", "warn", "Machine ID collision policy: update, warn, skip, merge")
)

// Environment variable support
const (
    EnvSpookyFactsPath = "SPOOKY_FACTS_PATH"
    EnvSpookyFactsFormat = "SPOOKY_FACTS_FORMAT"
    EnvSpookyCollisionPolicy = "SPOOKY_COLLISION_POLICY"
)

// Environment variable support
const (
    EnvSpookyFactsPath = "SPOOKY_FACTS_PATH"
    EnvSpookyFactsFormat = "SPOOKY_FACTS_FORMAT"
    EnvSpookyCollisionPolicy = "SPOOKY_COLLISION_POLICY"
)

// Get storage configuration with proper precedence hierarchy
func getStorageConfig() StorageOptions {
    opts := StorageOptions{}
    
    // 1. CLI flags (highest priority)
    if *storageType != "" {
        opts.Type = StorageType(*storageType)
    }
    if *storagePath != "" {
        opts.Path = *storagePath
    }
    
    // 2. Project configuration (project.conf.hcl in current directory)
    if opts.Type == "" || opts.Path == "" {
        projectConfig := loadProjectConfig("project.conf.hcl")
        if opts.Type == "" && projectConfig.Storage.Type != "" {
            opts.Type = projectConfig.Storage.Type
        }
        if opts.Path == "" && projectConfig.Storage.Path != "" {
            opts.Path = projectConfig.Storage.Path
        }
    }
    
    // 3. Environment variables
    if opts.Type == "" {
        if envFormat := os.Getenv(EnvSpookyFactsFormat); envFormat != "" {
            switch envFormat {
            case "badgerdb":
                opts.Type = StorageTypeBadger
            case "json":
                opts.Type = StorageTypeJSON
            default:
                opts.Type = StorageTypeBadger
            }
        }
    }
    if opts.Path == "" {
        if envPath := os.Getenv(EnvSpookyFactsPath); envPath != "" {
            opts.Path = envPath
        }
    }
    
    // 4. XDG user config ($HOME/.config/spooky/config.hcl)
    if opts.Type == "" || opts.Path == "" {
        userConfig := loadUserConfig()
        if opts.Type == "" && userConfig.Storage.Type != "" {
            opts.Type = userConfig.Storage.Type
        }
        if opts.Path == "" && userConfig.Storage.Path != "" {
            opts.Path = userConfig.Storage.Path
        }
    }
    
    // 5. XDG defaults (lowest priority)
    if opts.Type == "" {
        opts.Type = StorageTypeBadger // Default to badgerdb
    }
    // Path defaults to XDG_DATA_HOME/spooky/ (handled in NewFactStorage)
    
    return opts
}

// Load project-specific configuration
func loadProjectConfig(filename string) *Config {
    // Look for project.conf.hcl in current directory
    configPath := filepath.Join(".", filename)
    if _, err := os.Stat(configPath); err == nil {
        return parseConfigFile(configPath)
    }
    return &Config{}
}

// Load user configuration from XDG config home
func loadUserConfig() *Config {
    configPath := filepath.Join(GetXDGConfigHome(), "config.hcl")
    if _, err := os.Stat(configPath); err == nil {
        return parseConfigFile(configPath)
    }
    return &Config{}
}

// Get XDG config home directory
func GetXDGConfigHome() string {
    if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
        return filepath.Join(xdgConfigHome, "spooky")
    }
    
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return ".spooky" // Fallback
    }
    
    return filepath.Join(homeDir, ".config", "spooky")
}
```

// Add facts management commands
var FactsCmd = &cobra.Command{
    Use:   "facts",
    Short: "Manage server facts",
}

var DeleteFactsCmd = &cobra.Command{
    Use:   "delete",
    Short: "Delete facts from storage",
    RunE:  runFactsDelete,
}

var ListFactsCmd = &cobra.Command{
    Use:   "list",
    Short: "List facts in storage",
    RunE:  runFactsList,
}

var ExportFactsCmd = &cobra.Command{
    Use:   "export",
    Short: "Export facts to storage backend",
    RunE:  runFactsExport,
}

var ImportFactsCmd = &cobra.Command{
    Use:   "import",
    Short: "Import facts from storage backend",
    RunE:  runFactsImport,
}

var CollectFactsCmd = &cobra.Command{
    Use:   "collect",
    Short: "Collect facts from servers",
    RunE:  runFactsCollect,
}

var PurgeFactsCmd = &cobra.Command{
    Use:   "purge",
    Short: "Purge facts from storage (permanent deletion)",
    RunE:  runFactsPurge,
}

var UpdateFactsCmd = &cobra.Command{
    Use:   "update",
    Short: "Update specific facts in storage",
    RunE:  runFactsUpdate,
}

var BulkUpdateFactsCmd = &cobra.Command{
    Use:   "bulk-update",
    Short: "Bulk update multiple facts in storage",
    RunE:  runFactsBulkUpdate,
}

var ShowFactsCmd = &cobra.Command{
    Use:   "show",
    Short: "Show detailed information about specific facts",
    RunE:  runFactsShow,
}

var SearchFactsCmd = &cobra.Command{
    Use:   "search",
    Short: "Search facts with advanced queries",
    RunE:  runFactsSearch,
}

var StatsFactsCmd = &cobra.Command{
    Use:   "stats",
    Short: "Show statistics about facts database",
    RunE:  runFactsStats,
}

// Flags for fact deletion
var (
    deleteActionFile = flag.String("action-file", "", "Delete facts from specific action file")
    deleteServerName = flag.String("server-name", "", "Delete facts for specific server name")
    deleteProjectName = flag.String("project-name", "", "Delete facts from specific project name (portable)")
    deleteProjectPath = flag.String("project-path", "", "Delete facts from specific project path (absolute)")
    deleteEnvironment = flag.String("environment", "", "Delete facts from specific environment")
    deleteTeam = flag.String("team", "", "Delete facts from specific team")
    deleteConfirm = flag.Bool("confirm", false, "Confirm deletion without prompting")
)

// Flags for storage backends
var (
    storageBackend   = flag.String("backend", "badgerdb", "Storage backend (badgerdb, json, s3, minio, http)")
    s3Bucket         = flag.String("s3-bucket", "", "S3 bucket name")
    s3Prefix         = flag.String("s3-prefix", "", "S3 key prefix")
    s3Region         = flag.String("s3-region", "", "S3 region")
    s3Endpoint       = flag.String("s3-endpoint", "", "S3 endpoint (for MinIO)")
    s3AccessKey      = flag.String("s3-access-key", "", "S3 access key")
    s3SecretKey      = flag.String("s3-secret-key", "", "S3 secret key")
    s3UseSSL         = flag.Bool("s3-use-ssl", true, "Use SSL for S3/MinIO")
    httpURL          = flag.String("http-url", "", "HTTP endpoint URL")
    httpHeaders      = flag.StringSlice("http-header", nil, "HTTP headers (key=value)")
    httpTimeout      = flag.Duration("http-timeout", 30*time.Second, "HTTP request timeout")
)

// Global flags
var (
    configDir   = flag.String("config-dir", "", "Configuration directory")
    sshKeyPath  = flag.String("ssh-key-path", "", "SSH key path")
    logLevel    = flag.String("log-level", "info", "Log level")
    logFile     = flag.String("log-file", "", "Log file path")
    dryRun      = flag.Bool("dry-run", false, "Show what would be done")
    verbose     = flag.Bool("verbose", false, "Enable verbose output")
    quiet       = flag.Bool("quiet", false, "Suppress output")
)

// Flags for encryption
var (
    encryptFile      = flag.Bool("encrypt-file", false, "Encrypt entire JSON file with age")
    agePublicKey     = flag.String("age-public-key", "", "age public key for file encryption")
    encryptFields    = flag.Bool("encrypt-fields", false, "Encrypt sensitive fields individually")
    ageRecipient     = flag.String("age-recipient", "", "age recipient for field encryption")
    ageIdentityFile  = flag.String("age-identity", "", "age identity file for decryption")
    outputFile       = flag.String("output", "", "Output file path (default: stdout)")
    inputFile        = flag.String("input", "", "Input file path (default: stdin)")
)

// Flags for fact operations
var (
    // Update flags
    updateMachineID  = flag.String("machine-id", "", "Machine ID to update")
    updateServerName = flag.String("server-name", "", "Server name to update")
    updateHostname   = flag.String("hostname", "", "Hostname to update")
    updateTags       = flag.StringSlice("tag", nil, "Tags to add/update (key=value)")
    updateRemoveTags = flag.StringSlice("remove-tag", nil, "Tags to remove")
    updateActionFile = flag.String("action-file", "", "Action file to update")
    
    // Search flags
    searchQuery      = flag.String("query", "", "Search query (supports regex)")
    searchField      = flag.String("field", "", "Field to search in (hostname, tags, etc.)")
    searchLimit      = flag.Int("limit", 100, "Limit search results")
    searchFormat     = flag.String("format", "json", "Output format (json only)")
    
    // Bulk update flags
    bulkUpdateFile   = flag.String("file", "", "JSON file with bulk update operations")
    bulkUpdateDryRun = flag.Bool("dry-run", false, "Show what would be updated without making changes")
    
    // Purge flags
    purgeOlderThan   = flag.Duration("older-than", 0, "Purge facts older than specified duration")
    purgeBefore      = flag.String("before", "", "Purge facts before specified date (YYYY-MM-DD)")
    purgeForce       = flag.Bool("force", false, "Force purge without confirmation")
)
```

// Update ExecuteCmd
var ExecuteCmd = &cobra.Command{
    Use:   "execute <source>",
    Short: "Execute configuration files or remote sources",
    RunE: func(_ *cobra.Command, args []string) error {
        // Initialize storage
        storage, err := storage.NewFactStorage(storage.StorageOptions{
            Type: storage.StorageType(*storageType),
            Path: *storagePath,
        })
        if err != nil {
            return fmt.Errorf("failed to initialize storage: %w", err)
        }
        defer storage.Close()
        
        // Handle export/import
        if *exportFacts != "" {
            file, err := os.Create(*exportFacts)
            if err != nil {
                return err
            }
            defer file.Close()
            return storage.ExportToJSON(file)
        }
        
        if *importFacts != "" {
            file, err := os.Open(*importFacts)
            if err != nil {
                return err
            }
            defer file.Close()
            return storage.ImportFromJSON(file)
        }
        
        // Normal execution with storage
        return executeWithStorage(args[0], storage)
    },
}

// Implementation of fact deletion command
func runFactsDelete(_ *cobra.Command, args []string) error {
    // Initialize storage
    storage, err := storage.NewFactStorage(storage.StorageOptions{
        Type: storage.StorageType(*storageType),
        Path: *storagePath,
    })
    if err != nil {
        return fmt.Errorf("failed to initialize storage: %w", err)
    }
    defer storage.Close()
    
    // Build deletion query
    query := storage.FactQuery{}
    if *deleteActionFile != "" {
        query.ActionFile = *deleteActionFile
    }
    if *deleteServerName != "" {
        query.ServerName = *deleteServerName
    }
    if *deleteEnvironment != "" {
        if query.Tags == nil {
            query.Tags = make(map[string]string)
        }
        query.Tags["environment"] = *deleteEnvironment
    }
    if *deleteTeam != "" {
        if query.Tags == nil {
            query.Tags = make(map[string]string)
        }
        query.Tags["team"] = *deleteTeam
    }
    
    // Show what will be deleted
    facts, err := storage.QueryFacts(query)
    if err != nil {
        return fmt.Errorf("failed to query facts: %w", err)
    }
    
    if len(facts) == 0 {
        fmt.Println("No facts found matching the specified criteria.")
        return nil
    }
    
    fmt.Printf("Found %d facts matching the criteria:\n", len(facts))
    for _, fact := range facts {
        fmt.Printf("  - %s (%s) from %s\n", fact.ServerName, fact.MachineID, fact.ActionFile)
    }
    
    // Confirm deletion
    if !*deleteConfirm {
        fmt.Print("\nAre you sure you want to delete these facts? (y/N): ")
        var response string
        fmt.Scanln(&response)
        if response != "y" && response != "Y" {
            fmt.Println("Deletion cancelled.")
            return nil
        }
    }
    
    // Perform deletion
    deletedCount, err := storage.DeleteFacts(query)
    if err != nil {
        return fmt.Errorf("failed to delete facts: %w", err)
    }
    
    fmt.Printf("Successfully deleted %d facts.\n", deletedCount)
    return nil
}

// Implementation of fact export command
func runFactsExport(_ *cobra.Command, args []string) error {
    if len(args) < 1 {
        return fmt.Errorf("export requires output file path")
    }
    outputFile := args[0]
    
    // Initialize storage based on backend
    var storage FactStorage
    var err error
    
    switch *storageBackend {
    case "s3", "minio":
        config := &S3Config{
            Region:          *s3Region,
            Bucket:          *s3Bucket,
            Prefix:          *s3Prefix,
            AccessKeyID:     *s3AccessKey,
            SecretAccessKey: *s3SecretKey,
            Endpoint:        *s3Endpoint,
            UseSSL:          *s3UseSSL,
        }
        storage, err = NewS3Storage(config)
        
    case "http":
        config := &HTTPConfig{
            BaseURL:    *httpURL,
            Timeout:    *httpTimeout,
            RetryCount: 3,
        }
        // Parse headers
        for _, header := range *httpHeaders {
            if parts := strings.SplitN(header, "=", 2); len(parts) == 2 {
                config.Headers[parts[0]] = parts[1]
            }
        }
        storage, err = NewHTTPStorage(config)
        
    default:
        // Local storage (BadgerDB, JSON)
        storage, err = storage.NewFactStorage(getStorageConfig())
    }
    
    if err != nil {
        return fmt.Errorf("failed to initialize storage: %w", err)
    }
    defer storage.Close()
    
    // Build export options
    exportOpts := ExportOptions{
        SanitizeEnvironment: false, // Keep sensitive data for encryption
        SanitizeNetwork:     false,
        SanitizePaths:       false,
        SanitizeHostnames:   false,
        SanitizeMachineID:   false,
        EncryptFile:         *encryptFile,
        AgePublicKey:        *agePublicKey,
        EncryptFields:       *encryptFields,
        AgeRecipient:        *ageRecipient,
    }
    
    // Determine output destination
    var output io.Writer = os.Stdout
    if *outputFile != "" {
        file, err := os.Create(*outputFile)
        if err != nil {
            return fmt.Errorf("failed to create output file: %w", err)
        }
        defer file.Close()
        output = file
    }
    
    // Export facts with encryption
    if err := storage.ExportToJSONWithEncryption(output, exportOpts); err != nil {
        return fmt.Errorf("failed to export facts: %w", err)
    }
    
    fmt.Printf("Successfully exported facts to %s\n", outputFile)
    return nil
}

// Implementation of fact import command
func runFactsImport(_ *cobra.Command, args []string) error {
    if len(args) < 1 {
        return fmt.Errorf("import requires input file path")
    }
    inputFile := args[0]
    
    // Initialize storage based on backend
    var storage FactStorage
    var err error
    
    switch *storageBackend {
    case "s3", "minio":
        config := &S3Config{
            Region:          *s3Region,
            Bucket:          *s3Bucket,
            Prefix:          *s3Prefix,
            AccessKeyID:     *s3AccessKey,
            SecretAccessKey: *s3SecretKey,
            Endpoint:        *s3Endpoint,
            UseSSL:          *s3UseSSL,
        }
        storage, err = NewS3Storage(config)
        
    case "http":
        config := &HTTPConfig{
            BaseURL:    *httpURL,
            Timeout:    *httpTimeout,
            RetryCount: 3,
        }
        // Parse headers
        for _, header := range *httpHeaders {
            if parts := strings.SplitN(header, "=", 2); len(parts) == 2 {
                config.Headers[parts[0]] = parts[1]
            }
        }
        storage, err = NewHTTPStorage(config)
        
    default:
        // Local storage (BadgerDB, JSON)
        storage, err = storage.NewFactStorage(getStorageConfig())
    }
    
    if err != nil {
        return fmt.Errorf("failed to initialize storage: %w", err)
    }
    defer storage.Close()
    
    // Determine input source
    var input io.Reader = os.Stdin
    if *inputFile != "" {
        file, err := os.Open(*inputFile)
        if err != nil {
            return fmt.Errorf("failed to open input file: %w", err)
        }
        defer file.Close()
        input = file
    }
    
    // Import facts with decryption
    if err := storage.ImportFromJSONWithDecryption(input, *ageIdentityFile); err != nil {
        return fmt.Errorf("failed to import facts: %w", err)
    }
    
    fmt.Printf("Successfully imported facts from %s\n", inputFile)
    return nil
}

// Implementation of fact gathering command
func runFactsGather(_ *cobra.Command, args []string) error {
    if len(args) < 2 {
        return fmt.Errorf("gather requires machine name and action file")
    }
    machineName := args[0]
    actionFile := args[1]
    
    // Initialize storage
    storage, err := storage.NewFactStorage(getStorageConfig())
    if err != nil {
        return fmt.Errorf("failed to initialize storage: %w", err)
    }
    defer storage.Close()
    
    // Create fact manager with storage
    manager := facts.NewManagerWithStorage(nil, storage)
    
    // Gather facts
    collection, err := manager.GatherAndPersistFacts(machineName)
    if err != nil {
        return fmt.Errorf("failed to gather facts: %w", err)
    }
    
    fmt.Printf("Successfully gathered facts for %s (%d facts)\n", machineName, len(collection.Facts))
    return nil
}

// Implementation of fact purge command
func runFactsPurge(_ *cobra.Command, args []string) error {
    storage, err := storage.NewFactStorage(getStorageConfig())
    if err != nil {
        return fmt.Errorf("failed to initialize storage: %w", err)
    }
    defer storage.Close()
    
    // Build purge query
    query := storage.FactQuery{}
    
    // Add time-based filters
    if *purgeOlderThan > 0 {
        cutoffTime := time.Now().Add(-*purgeOlderThan)
        query.UpdatedBefore = &cutoffTime
    }
    
    if *purgeBefore != "" {
        beforeTime, err := time.Parse("2006-01-02", *purgeBefore)
        if err != nil {
            return fmt.Errorf("invalid date format: %s (use YYYY-MM-DD)", *purgeBefore)
        }
        query.UpdatedBefore = &beforeTime
    }
    
    // Show what will be purged
    facts, err := storage.QueryFacts(query)
    if err != nil {
        return fmt.Errorf("failed to query facts: %w", err)
    }
    
    if len(facts) == 0 {
        fmt.Println("No facts found matching purge criteria.")
        return nil
    }
    
    fmt.Printf("Found %d facts to purge:\n", len(facts))
    for _, fact := range facts {
        fmt.Printf("  - %s (%s) updated %s\n", fact.ServerName, fact.MachineID, fact.UpdatedAt.Format("2006-01-02 15:04:05"))
    }
    
    // Confirm purge
    if !*purgeForce {
        fmt.Print("\nAre you sure you want to PERMANENTLY DELETE these facts? (type 'PURGE' to confirm): ")
        var response string
        fmt.Scanln(&response)
        if response != "PURGE" {
            fmt.Println("Purge cancelled.")
            return nil
        }
    }
    
    // Perform purge
    deletedCount, err := storage.DeleteFacts(query)
    if err != nil {
        return fmt.Errorf("failed to purge facts: %w", err)
    }
    
    fmt.Printf("Successfully purged %d facts.\n", deletedCount)
    return nil
}

// Implementation of fact update command
func runFactsUpdate(_ *cobra.Command, args []string) error {
    if len(args) < 1 {
        return fmt.Errorf("update requires machine ID")
    }
    machineID := args[0]
    
    storage, err := storage.NewFactStorage(getStorageConfig())
    if err != nil {
        return fmt.Errorf("failed to initialize storage: %w", err)
    }
    defer storage.Close()
    
    // Get existing facts
    facts, err := storage.GetServerFacts(machineID)
    if err != nil {
        return fmt.Errorf("failed to get facts for %s: %w", machineID, err)
    }
    
    // Apply updates
    if *updateServerName != "" {
        facts.ServerName = *updateServerName
    }
    if *updateHostname != "" {
        facts.Hostname = *updateHostname
    }
    if *updateActionFile != "" {
        facts.ActionFile = *updateActionFile
    }
    
    // Update tags
    if len(*updateTags) > 0 {
        if facts.Tags == nil {
            facts.Tags = make(map[string]string)
        }
        for _, tag := range *updateTags {
            if parts := strings.SplitN(tag, "=", 2); len(parts) == 2 {
                facts.Tags[parts[0]] = parts[1]
            }
        }
    }
    
    // Remove tags
    if len(*updateRemoveTags) > 0 {
        for _, tag := range *updateRemoveTags {
            delete(facts.Tags, tag)
        }
    }
    
    facts.UpdatedAt = time.Now()
    
    // Save updated facts
    if err := storage.SetServerFacts(machineID, facts); err != nil {
        return fmt.Errorf("failed to update facts: %w", err)
    }
    
    fmt.Printf("Successfully updated facts for %s\n", machineID)
    return nil
}

// Implementation of bulk update command
func runFactsBulkUpdate(_ *cobra.Command, args []string) error {
    if *bulkUpdateFile == "" {
        return fmt.Errorf("bulk-update requires --file with JSON operations")
    }
    
    storage, err := storage.NewFactStorage(getStorageConfig())
    if err != nil {
        return fmt.Errorf("failed to initialize storage: %w", err)
    }
    defer storage.Close()
    
    // Read bulk update file
    data, err := os.ReadFile(*bulkUpdateFile)
    if err != nil {
        return fmt.Errorf("failed to read bulk update file: %w", err)
    }
    
    var operations []BulkUpdateOperation
    if err := json.Unmarshal(data, &operations); err != nil {
        return fmt.Errorf("failed to parse bulk update file: %w", err)
    }
    
    var updatedCount int
    var errors []error
    
    for _, op := range operations {
        if *bulkUpdateDryRun {
            fmt.Printf("Would update %s: %s\n", op.MachineID, op.Description)
            continue
        }
        
        if err := storage.ApplyBulkUpdate(op); err != nil {
            errors = append(errors, fmt.Errorf("failed to update %s: %w", op.MachineID, err))
            continue
        }
        
        updatedCount++
        fmt.Printf("Updated %s: %s\n", op.MachineID, op.Description)
    }
    
    if *bulkUpdateDryRun {
        fmt.Printf("Dry run: would update %d facts\n", len(operations))
    } else {
        fmt.Printf("Successfully updated %d facts\n", updatedCount)
        if len(errors) > 0 {
            fmt.Printf("Errors: %d\n", len(errors))
            for _, err := range errors {
                fmt.Printf("  - %v\n", err)
            }
        }
    }
    
    return nil
}

// Implementation of fact show command
func runFactsShow(_ *cobra.Command, args []string) error {
    if len(args) < 1 {
        return fmt.Errorf("show requires machine ID")
    }
    machineID := args[0]
    
    storage, err := storage.NewFactStorage(getStorageConfig())
    if err != nil {
        return fmt.Errorf("failed to initialize storage: %w", err)
    }
    defer storage.Close()
    
    facts, err := storage.GetServerFacts(machineID)
    if err != nil {
        return fmt.Errorf("failed to get facts for %s: %w", machineID, err)
    }
    
    // Display detailed facts
    fmt.Printf("Facts for %s:\n", machineID)
    fmt.Printf("  Server Name: %s\n", facts.ServerName)
    fmt.Printf("  Hostname: %s\n", facts.Hostname)
    fmt.Printf("  Action File: %s\n", facts.ActionFile)
    fmt.Printf("  Project Name: %s\n", facts.ProjectName)
    fmt.Printf("  Project Path: %s\n", facts.ProjectPath)
    fmt.Printf("  OS: %s %s\n", facts.OS, facts.OSVersion)
    fmt.Printf("  CPU: %s (%d cores, %s)\n", facts.CPU.Model, facts.CPU.Cores, facts.CPU.Arch)
    fmt.Printf("  Memory: %d MB total, %d MB used\n", facts.Memory.Total/1024/1024, facts.Memory.Used/1024/1024)
    fmt.Printf("  IP Address: %s\n", facts.IPAddress)
    fmt.Printf("  Created: %s\n", facts.CreatedAt.Format("2006-01-02 15:04:05"))
    fmt.Printf("  Updated: %s\n", facts.UpdatedAt.Format("2006-01-02 15:04:05"))
    
    if len(facts.Tags) > 0 {
        fmt.Printf("  Tags:\n")
        for k, v := range facts.Tags {
            fmt.Printf("    %s: %s\n", k, v)
        }
    }
    
    return nil
}

// Implementation of fact search command
func runFactsSearch(_ *cobra.Command, args []string) error {
    if *searchQuery == "" {
        return fmt.Errorf("search requires --query")
    }
    
    storage, err := storage.NewFactStorage(getStorageConfig())
    if err != nil {
        return fmt.Errorf("failed to initialize storage: %w", err)
    }
    defer storage.Close()
    
    // Build search query
    query := storage.FactQuery{
        SearchQuery: *searchQuery,
        SearchField: *searchField,
        Limit:       *searchLimit,
    }
    
    facts, err := storage.SearchFacts(query)
    if err != nil {
        return fmt.Errorf("failed to search facts: %w", err)
    }
    
    if len(facts) == 0 {
        fmt.Println("No facts found matching search criteria.")
        return nil
    }
    
    // Output results as JSON (for use with jq and other tools)
    jsonData, err := json.MarshalIndent(facts, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal results: %w", err)
    }
    fmt.Println(string(jsonData))
    
    return nil
}

// Implementation of fact stats command
func runFactsStats(_ *cobra.Command, args []string) error {
    storage, err := storage.NewFactStorage(getStorageConfig())
    if err != nil {
        return fmt.Errorf("failed to initialize storage: %w", err)
    }
    defer storage.Close()
    
    stats, err := storage.GetStats()
    if err != nil {
        return fmt.Errorf("failed to get stats: %w", err)
    }
    
    fmt.Printf("Facts Database Statistics:\n")
    fmt.Printf("  Total Facts: %d\n", stats.TotalFacts)
    fmt.Printf("  Total Servers: %d\n", stats.TotalServers)
    fmt.Printf("  Database Size: %s\n", formatBytes(stats.DatabaseSize))
    fmt.Printf("  Oldest Fact: %s\n", stats.OldestFact.Format("2006-01-02 15:04:05"))
    fmt.Printf("  Newest Fact: %s\n", stats.NewestFact.Format("2006-01-02 15:04:05"))
    fmt.Printf("  Average Facts per Server: %.1f\n", stats.AverageFactsPerServer)
    
    if len(stats.TopTags) > 0 {
        fmt.Printf("  Top Tags:\n")
        for _, tag := range stats.TopTags {
            fmt.Printf("    %s: %d\n", tag.Key, tag.Count)
        }
    }
    
    if len(stats.TopOS) > 0 {
        fmt.Printf("  Top Operating Systems:\n")
        for _, os := range stats.TopOS {
            fmt.Printf("    %s: %d\n", os.Name, os.Count)
        }
    }
    
    return nil
}

// Implementation of fact validation command
func runFactsValidate(_ *cobra.Command, args []string) error {
    storage, err := storage.NewFactStorage(getStorageConfig())
    if err != nil {
        return fmt.Errorf("failed to initialize storage: %w", err)
    }
    defer storage.Close()
    
    // Get all facts for validation
    facts, err := storage.QueryFacts(FactQuery{})
    if err != nil {
        return fmt.Errorf("failed to query facts: %w", err)
    }
    
    var validationErrors []string
    
    for _, fact := range facts {
        // Validate required fields
        if fact.MachineID == "" {
            validationErrors = append(validationErrors, fmt.Sprintf("machine_id missing for %s", fact.ServerID))
        }
        if fact.Hostname == "" {
            validationErrors = append(validationErrors, fmt.Sprintf("hostname missing for %s", fact.ServerID))
        }
        if fact.OS == "" {
            validationErrors = append(validationErrors, fmt.Sprintf("OS missing for %s", fact.ServerID))
        }
        
        // Validate data types and ranges
        if fact.CPU.Cores <= 0 {
            validationErrors = append(validationErrors, fmt.Sprintf("invalid CPU cores for %s: %d", fact.ServerID, fact.CPU.Cores))
        }
        if fact.Memory.Total == 0 {
            validationErrors = append(validationErrors, fmt.Sprintf("memory total missing for %s", fact.ServerID))
        }
        
        // Validate timestamps
        if fact.CreatedAt.IsZero() {
            validationErrors = append(validationErrors, fmt.Sprintf("created_at missing for %s", fact.ServerID))
        }
        if fact.UpdatedAt.IsZero() {
            validationErrors = append(validationErrors, fmt.Sprintf("updated_at missing for %s", fact.ServerID))
        }
    }
    
    if len(validationErrors) > 0 {
        fmt.Fprintf(os.Stderr, "Validation failed with %d errors:\n", len(validationErrors))
        for _, err := range validationErrors {
            fmt.Fprintf(os.Stderr, "  - %s\n", err)
        }
        return fmt.Errorf("fact validation failed")
    }
    
    fmt.Printf("Validation passed: %d facts validated successfully\n", len(facts))
    return nil
}

// Implementation of fact query command
func runFactsQuery(_ *cobra.Command, args []string) error {
    if len(args) < 1 {
        return fmt.Errorf("query requires an expression")
    }
    expression := args[0]
    
    storage, err := storage.NewFactStorage(getStorageConfig())
    if err != nil {
        return fmt.Errorf("failed to initialize storage: %w", err)
    }
    defer storage.Close()
    
    // Parse query expression
    query, err := parseQueryExpression(expression)
    if err != nil {
        return fmt.Errorf("failed to parse query expression: %w", err)
    }
    
    // Execute query
    facts, err := storage.QueryFacts(query)
    if err != nil {
        return fmt.Errorf("failed to execute query: %w", err)
    }
    
    // Output results in JSON format
    if err := json.NewEncoder(os.Stdout).Encode(facts); err != nil {
        return fmt.Errorf("failed to encode results: %w", err)
    }
    
    return nil
}

func parseQueryExpression(expression string) (FactQuery, error) {
    query := FactQuery{}
    
    // Simple expression parser - can be enhanced for complex queries
    // Format: field=value,field2=value2
    parts := strings.Split(expression, ",")
    for _, part := range parts {
        if strings.Contains(part, "=") {
            kv := strings.SplitN(part, "=", 2)
            if len(kv) != 2 {
                return query, fmt.Errorf("invalid query format: %s", part)
            }
            
            key := strings.TrimSpace(kv[0])
            value := strings.TrimSpace(kv[1])
            
            switch key {
            case "os":
                query.OS = value
            case "hostname":
                query.Hostname = value
            case "project":
                query.ProjectName = value
            case "machine":
                query.MachineName = value
            case "tag":
                if query.Tags == nil {
                    query.Tags = make(map[string]string)
                }
                // Handle tag format: tag=key:value
                if strings.Contains(value, ":") {
                    tagParts := strings.SplitN(value, ":", 2)
                    query.Tags[tagParts[0]] = tagParts[1]
                } else {
                    query.Tags[value] = ""
                }
            case "limit":
                if limit, err := strconv.Atoi(value); err == nil {
                    query.Limit = limit
                }
            default:
                return query, fmt.Errorf("unknown query field: %s", key)
            }
        }
    }
    
    return query, nil
}

// Command registration
func registerFactsCommands(factsCmd *cobra.Command) {
    // Core commands from Issue #61
    factsCmd.AddCommand(&cobra.Command{
        Use:   "gather [hosts]",
        Short: "Gather facts from machines",
        Long:  "Gather system facts from remote or local machines using SSH or local execution",
        RunE:  runFactsGather,
    })
    
    factsCmd.AddCommand(&cobra.Command{
        Use:   "import <source> [flags]",
        Short: "Import facts from external sources",
        Long:  "Import facts from JSON files, S3, MinIO, or HTTP endpoints with optional decryption",
        RunE:  runFactsImport,
    })
    
    factsCmd.AddCommand(&cobra.Command{
        Use:   "export [flags]",
        Short: "Export facts to various formats",
        Long:  "Export facts to JSON with optional encryption and sanitization",
        RunE:  runFactsExport,
    })
    
    factsCmd.AddCommand(&cobra.Command{
        Use:   "validate [flags]",
        Short: "Validate facts against rules",
        Long:  "Validate facts for completeness, data types, and business rules",
        RunE:  runFactsValidate,
    })
    
    factsCmd.AddCommand(&cobra.Command{
        Use:   "query <expression> [flags]",
        Short: "Query facts with filters",
        Long:  "Query facts using expression syntax: field=value,field2=value2",
        RunE:  runFactsQuery,
    })
    
    // Additional useful commands
    factsCmd.AddCommand(&cobra.Command{
        Use:   "list [flags]",
        Short: "List all facts",
        Long:  "List all facts in the database with optional filtering",
        RunE:  runFactsList,
    })
    
    factsCmd.AddCommand(&cobra.Command{
        Use:   "show <machine-id> [flags]",
        Short: "Show detailed facts for a machine",
        Long:  "Show detailed facts for a specific machine by machine ID",
        RunE:  runFactsShow,
    })
    
    factsCmd.AddCommand(&cobra.Command{
        Use:   "delete [flags]",
        Short: "Delete facts by criteria",
        Long:  "Delete facts matching specified criteria",
        RunE:  runFactsDelete,
    })
    
    factsCmd.AddCommand(&cobra.Command{
        Use:   "purge [flags]",
        Short: "Purge old facts",
        Long:  "Purge facts older than specified time",
        RunE:  runFactsPurge,
    })
    
    factsCmd.AddCommand(&cobra.Command{
        Use:   "update <machine-id> [flags]",
        Short: "Update a single fact",
        Long:  "Update specific fields of a fact by machine ID",
        RunE:  runFactsUpdate,
    })
    
    factsCmd.AddCommand(&cobra.Command{
        Use:   "bulk-update [flags]",
        Short: "Bulk update facts from JSON",
        Long:  "Bulk update multiple facts from JSON file",
        RunE:  runFactsBulkUpdate,
    })
    
    factsCmd.AddCommand(&cobra.Command{
        Use:   "search [flags]",
        Short: "Search facts with text query",
        Long:  "Search facts using text queries with regex support",
        RunE:  runFactsSearch,
    })
    
    factsCmd.AddCommand(&cobra.Command{
        Use:   "stats [flags]",
        Short: "Show database statistics",
        Long:  "Show comprehensive database statistics and metrics",
        RunE:  runFactsStats,
    })
}

// Implementation of fact listing command
func runFactsList(_ *cobra.Command, args []string) error {
    // Initialize storage
    storage, err := storage.NewFactStorage(storage.StorageOptions{
        Type: storage.StorageType(*storageType),
        Path: *storagePath,
    })
    if err != nil {
        return fmt.Errorf("failed to initialize storage: %w", err)
    }
    defer storage.Close()
    
    // Build query (empty query lists all facts)
    query := storage.FactQuery{}
    if *deleteActionFile != "" {
        query.ActionFile = *deleteActionFile
    }
    if *deleteServerName != "" {
        query.ServerName = *deleteServerName
    }
    if *deleteEnvironment != "" {
        if query.Tags == nil {
            query.Tags = make(map[string]string)
        }
        query.Tags["environment"] = *deleteEnvironment
    }
    if *deleteTeam != "" {
        if query.Tags == nil {
            query.Tags = make(map[string]string)
        }
        query.Tags["team"] = *deleteTeam
    }
    
    facts, err := storage.QueryFacts(query)
    if err != nil {
        return fmt.Errorf("failed to query facts: %w", err)
    }
    
    if len(facts) == 0 {
        fmt.Println("No facts found.")
        return nil
    }
    
    fmt.Printf("Found %d facts:\n\n", len(facts))
    for _, fact := range facts {
        fmt.Printf("Machine: %s\n", fact.MachineName)
        fmt.Printf("  System ID: %s\n", fact.SystemID)
        fmt.Printf("  Action File: %s\n", fact.ActionFile)
        fmt.Printf("  Hostname: %s\n", fact.Hostname)
        if fact.IPAddress != "" {
            fmt.Printf("  IP Address: %s\n", fact.IPAddress)
        }
        if len(fact.Tags) > 0 {
            fmt.Printf("  Tags: %v\n", fact.Tags)
        }
        fmt.Printf("  Updated: %s\n", fact.UpdatedAt.Format("2006-01-02 15:04:05"))
        fmt.Println()
    }
    
    return nil
}
```

#### Testing Instructions

**Step 1: Run Unit Tests**
```bash
# Test fact storage implementations
go test ./internal/facts/...

# Test collision detection
go test ./internal/facts/ -run TestCollisionDetection

# Test configuration hierarchy
go test ./internal/facts/ -run TestConfigurationPrecedence

# Run all tests
go test ./...
```

**Step 2: Test Configuration Hierarchy**
```bash
# Test XDG default location
unset SPOOKY_FACTS_PATH
unset SPOOKY_FACTS_FORMAT
spooky facts list  # Should use $HOME/.local/share/spooky/facts.db

# Test environment variables
export SPOOKY_FACTS_PATH="/tmp/test-facts.db"
export SPOOKY_FACTS_FORMAT="json"
spooky facts list  # Should use /tmp/test-facts.json

# Test project configuration
echo 'storage { path = ".facts.db" type = "json" }' > project.conf.hcl
spooky facts list  # Should use .facts.json in current directory

# Test CLI override
spooky facts list --storage-path "/tmp/override.db"  # Should use /tmp/override.db
```

**Step 3: Test Collision Detection**
```bash
# Create facts with same machine ID but different hostnames
spooky facts gather --machine "server1" --hostname "host1.example.com"
spooky facts gather --machine "server2" --hostname "host2.example.com"  # Should detect collision

# Test different collision policies
spooky facts gather --collision-policy warn  # Should warn on collision
spooky facts gather --collision-policy merge  # Should merge facts
spooky facts gather --collision-policy skip   # Should skip collection
```

**Step 4: Test Centralized Workflow**
```bash
# Central server gathers facts from fleet
spooky facts gather "web-001" "fleet/deploy-web.hcl"
spooky facts gather "web-002" "fleet/deploy-web.hcl"
spooky facts gather "db-001" "fleet/deploy-db.hcl"

# Export facts to S3
spooky facts export --backend s3 --s3-bucket "spooky-facts" --s3-prefix "fleet" --output "facts-$(date +%Y-%m-%d-%H-%M).json"

# Export facts to MinIO
spooky facts export --backend minio --s3-bucket "spooky-facts" --s3-endpoint "http://minio:9000" --s3-prefix "fleet" --output "facts-$(date +%Y-%m-%d-%H-%M).json"

# Export facts to HTTP endpoint
spooky facts export --backend http --http-url "https://facts-api.company.com/export" --http-header "Authorization=Bearer token123" --output "facts-$(date +%Y-%m-%d-%H-%M).json"

# Team member imports facts from S3
spooky facts import --backend s3 --s3-bucket "spooky-facts" --s3-prefix "fleet" "facts-2024-01-15-14-30.json"

# Team member imports facts from HTTP
spooky facts import --backend http --http-url "https://facts-api.company.com/facts/latest" "latest-facts.json"

# Verify imported facts
spooky facts list --project-name "fleet"
```

**Step 5: Test Storage Backends**
```bash
# Test S3 storage
spooky facts export --backend s3 --s3-bucket "my-bucket" --s3-region "us-west-2" --s3-prefix "spooky" --output "facts.json"

# Test MinIO storage
spooky facts export --backend minio --s3-bucket "spooky" --s3-endpoint "http://localhost:9000" --s3-access-key "minioadmin" --s3-secret-key "minioadmin" --s3-use-ssl false --output "facts.json"

# Test HTTP storage (HTTPS required for writes)
spooky facts export --backend http --http-url "https://localhost:8443/api/facts" --http-header "Content-Type=application/json" --output "facts.json"

# Test local storage (default)
spooky facts export --output "facts.json"
```

**Step 6: Test Database Operations**
```bash
# Show database statistics
spooky facts stats

# Show detailed facts for a specific machine
spooky facts show "6538b193d562410ab480b1ce22469fe6"

# Search facts with text query (JSON output for use with jq)
spooky facts search --query "web-server" --field "hostname" | jq '.[] | .hostname'
spooky facts search --query "production" --field "tags" | jq '.[] | {machine_id, server_name, tags}'

# Update a single fact
spooky facts update "6538b193d562410ab480b1ce22469fe6" --server-name "web-001-prod" --tag "environment=production"

# Remove tags from a fact
spooky facts update "6538b193d562410ab480b1ce22469fe6" --remove-tag "staging" --remove-tag "test"

# Bulk update facts from JSON file
spooky facts bulk-update --file "bulk-updates.json" --dry-run
spooky facts bulk-update --file "bulk-updates.json"

# Purge old facts
spooky facts purge --older-than 30d
spooky facts purge --before "2024-01-01" --force

# Delete facts by criteria
spooky facts delete --action-file "staging/deploy.hcl" --confirm
spooky facts delete --environment "staging" --team "web-team"
```

**Step 7: Test Advanced Queries**
```bash
# List facts (JSON output for use with jq)
spooky facts list | jq '.[] | {machine_id, server_name, hostname, os}'

# Search with regex (JSON output for use with jq)
spooky facts search --query "web.*prod" --field "hostname" | jq '.[] | .hostname'
spooky facts search --query "ubuntu.*20" --field "os" | jq '.[] | {machine_id, os, os_version}'

# Query by time range (JSON output for use with jq)
spooky facts list --updated-after "2024-01-01" --updated-before "2024-01-31" | jq '.[] | {machine_id, server_name, updated_at}'

# Complex tag queries (JSON output for use with jq)
spooky facts list --tag "team=web-team" --tag "environment=production" | jq '.[] | {machine_id, server_name, tags}'
```

**Note**: All fact queries output JSON for easy integration with tools like `jq`. This approach:
- Avoids unwieldy table output with many facts
- Enables powerful data manipulation with `jq`
- Provides consistent, parseable output
- Allows easy filtering, sorting, and transformation

**Step 8: Bulk Update Examples**
```json
// bulk-updates.json
[
  {
    "machine_id": "6538b193d562410ab480b1ce22469fe6",
    "description": "Update web-001 to production environment",
    "updates": {
      "server_name": "web-001-prod",
      "action_file": "production/deploy-web.hcl"
    },
    "add_tags": {
      "environment": "production",
      "team": "web-team"
    },
    "remove_tags": ["staging", "test"]
  },
  {
    "machine_id": "a1b2c3d4e5f6789012345678901234567",
    "description": "Update db-001 tags",
    "add_tags": {
      "environment": "production",
      "team": "db-team",
      "backup": "enabled"
    }
  }
]
```

**Step 2: Run Integration Tests**
```bash
# Test BadgerDB storage
go test -v ./internal/storage/badger/...

# Test JSON storage
go test -v ./internal/storage/json/...

# Test CLI integration
go build -o build/spooky
./build/spooky execute examples/actions/example.hcl --storage badger
./build/spooky execute examples/actions/example.hcl --storage json --storage-path facts.json
```

**Step 3: Performance Testing**
```bash
# Run benchmarks
go test -bench=. ./internal/storage/...

# Test with large datasets
go test -bench=BenchmarkBadgerDB -benchmem ./internal/storage/badger/
```

### Implementation Checklist

- [ ] Add BadgerDB dependency to go.mod
- [ ] Create storage package structure
- [ ] Implement storage interface
- [ ] Implement BadgerDB store
- [ ] Implement JSON store
- [ ] Update fact manager integration
- [ ] Add CLI flags and commands
- [ ] Write unit tests
- [ ] Write integration tests
- [ ] Add performance benchmarks
- [ ] Update documentation
- [ ] Test with real configurations

## Complete CLI Command Structure

**Core Commands (from Issue #61):**
```bash
spooky facts gather [hosts]                  # Gather facts from machines
spooky facts import <source> [flags]         # Import facts from external sources
spooky facts export [flags]                  # Export facts to various formats
spooky facts validate [flags]                # Validate facts against rules
spooky facts query <expression> [flags]      # Query facts with filters
```

**Additional Commands:**
```bash
spooky facts list [flags]                    # List all facts
spooky facts show <machine-id> [flags]       # Show detailed facts for a machine
spooky facts delete [flags]                  # Delete facts by criteria
spooky facts purge [flags]                   # Purge old facts
spooky facts update <machine-id> [flags]     # Update a single fact
spooky facts bulk-update [flags]             # Bulk update facts from JSON
spooky facts search [flags]                  # Search facts with text query
spooky facts stats [flags]                   # Show database statistics
```

**Global Flags (from Issue #61):**
```bash
--config-dir string     # Configuration directory
--ssh-key-path string   # SSH key path
--log-level string      # Log level (debug, info, warn, error)
--log-file string       # Log file path
--dry-run              # Show what would be done
--verbose              # Enable verbose output
--quiet                # Suppress output
```

## Next Steps

1. **Review and approve** this implementation plan
2. **Set up development environment** with BadgerDB dependency
3. **Begin Phase 1** with storage interface creation
4. **Establish testing framework** for storage implementations
5. **Create integration test suite** for end-to-end validation

## References

- [Issue #56: Implement scalable fact storage system](https://github.com/snassar/spooky/issues/56)
- [BadgerDB Documentation](https://dgraph.io/docs/badger/)
- [Current Fact System Implementation](../internal/facts/)
- [Spooky CLI Specification](../cli-specification.md) 