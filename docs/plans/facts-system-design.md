# Facts System: Comprehensive Implementation Plan

## Overview

This document is the authoritative source for all facts system implementation details in spooky. It covers facts collection, storage, validation, caching, and integration with the actions and template systems.

**Schema Integration**: This facts system implements the schema validation patterns and data structures defined in [Schema System](../schema-system.md) for comprehensive facts validation, storage format consistency, and schema-based data integrity.

**Architecture Integration**: Facts integrate with the overall spooky architecture as described in [Spooky Design](../spooky-design.md), providing machine data to the actions system for dynamic running and template rendering.

## System Integration

This facts system integrates with other core Spooky systems to provide comprehensive fact management and data sharing:

### **Project System Integration**
- **Project Initialization**: Facts database creation during `spooky project init` (see [Project System](../project-system.md))
- **Project Isolation**: Project-specific facts storage
- **Project Configuration**: Facts storage settings in `project.hcl` configuration
- **Project Context**: Facts available in project run context

### **Configuration System Integration**
- **Global Configuration**: Facts collection settings from `$XDG_CONFIG_HOME/spooky/spooky.hcl` (see [Configuration System](../configuration-system.md))
- **Facts Settings**: Collection timeouts, parallel collection, retry attempts
- **Storage Configuration**: BadgerDB/JSON storage settings, compression, encryption
- **Performance Settings**: Concurrent fact collectors, memory limits, cache sizes
- **Configuration Precedence**: CLI flags → Project config → Global config → Environment variables → Defaults

### **Variables System Integration**
- **Template Context**: Facts available as variables in templates (see [Variables System](../variables-system.md))
- **Variable Resolution**: Facts can be referenced in variable definitions
- **Data Merging**: Facts and variables merge in template rendering context
- **Export Integration**: Facts export includes variable context information

### **Schema System Integration**
- **Schema Validation**: All facts data validated against embedded schemas (see [Schema System](../schema-system.md))
- **Schema Composition**: Runtime schema composition for facts validation
- **Schema Evolution**: Facts schemas evolve with system changes
- **Schema Documentation**: Facts schemas documented and versioned

### **CLI System Integration**
- **Facts Commands**: Facts management through `spooky facts` commands (see [CLI System](../cli-system.md))
- **Facts Collection**: `spooky facts gather` for machine facts collection
- **Facts Validation**: `spooky facts validate` for facts database validation
- **Facts Export**: `spooky facts export` for facts data export
- **Facts Discovery**: `spooky facts list` for facts discovery and filtering

### **Actions System Integration**
- **Facts in Actions**: Actions can access machine facts for dynamic running (see [Actions System](../actions-system.md))
- **Facts-Based Commands**: Actions can use facts data in command running
- **Template Context**: Actions have access to facts through template context
- **Facts Resolution**: Actions resolve facts at run time for each machine

### **Machines System Integration**
- **Machine Facts**: Facts collection uses machine inventory for target identification (see [Machines System](../machines-system.md))
- **Machine Metadata**: Machine inventory provides metadata for facts collection
- **Facts Storage**: Machine facts are stored using machine IDs and names from inventory
- **Facts Association**: Collected facts are associated with specific machines from inventory
- **Enterprise Scale**: Facts collection optimized for large machine inventories

### **Template System Integration**
- **Facts in Templates**: Templates can access machine facts for dynamic content (see [Template System](../template-system.md))
- **System Facts**: Templates have access to system facts (OS, hardware, network)
- **Custom Facts**: Templates can use custom facts for project-specific data
- **Facts Functions**: Template functions for accessing facts data (`system()`, `custom()`)
- **Template Context**: Facts data available in template rendering context

## Current State Analysis

### ✅ Already Implemented
1. **Go 1.24+** - Project uses Go 1.24.3
2. **Existing Fact System** - Well-structured with types, collectors, and manager
3. **JSON/HCL Support** - Facts already have JSON tags and HCL support
4. **Project Structure** - Clean internal organization
5. **Noun-Verb CLI** - Integrated with `spooky facts` commands
6. **Project Integration** - Facts storage integrated with `spooky project init`
7. **Schema Validation** - Facts database validation during project initialization
8. **Storage Interface** - Clean separation between creation and reading operations

### Current Fact System Structure
```
internal/facts/
├── manager.go         # Fact collection coordination
├── types.go           # Fact data structures
├── storage.go         # Storage interface and factory
├── badger_storage.go  # BadgerDB implementation
├── json_storage.go    # JSON storage implementation
├── hcl_storage.go     # HCL storage implementation
├── exporter.go        # Export/import functionality
├── ssh_collector.go   # SSH-based fact collection
├── local_collector.go # Local system fact collection
├── hcl_collector.go   # HCL configuration facts
├── opentofu_collector.go # OpenTofu state facts
├── base_collector.go  # Shared collector functionality
└── stub_collector.go  # Stub implementations
```

## Memory Architecture

### **1. In-Memory Facts Storage**

Based on our analysis and the architectural principles outlined in [Spooky Design](../spooky-design.md), we've decided on an in-memory approach for facts storage:

```
project/
├── project.hcl                              # Project configuration
├── machines.hcl                             # Machine inventory
├── actions.hcl                              # Action definitions
├── variables.hcl                            # Variable definitions
└── ...
# Facts are stored in memory during action runs
```

### **2. Fact Classification**

**In-Memory Facts** (stored in memory during action runs):
- ✅ **gopsutil facts** (hardware, system, network)
- ✅ **BIOS information** (vendor, version, date)
- ✅ **SSH host keys** (public keys)
- ✅ **Package manager detection** (apt, yum, dnf, etc.)
- ✅ **Service manager detection** (systemd, upstart, etc.)
- ✅ **SELinux status** (enabled, mode, type)
- ✅ **System identification** (machine-id, hostname, IPs)
- ✅ **Application versions** (nginx, postgres, etc.)
- ✅ **Deployment states** (installed, configured, running)
- ✅ **Configuration values** (custom app configs)
- ✅ **Project-specific variables** (environment-specific data)

### **3. Enhanced Fact Collection**

Our current implementation uses `github.com/shirou/gopsutil/v4` which provides comprehensive system information. We can enhance this with additional fact collectors:

#### **Package Manager Detection**
```go
// Package manager detection via system commands
func detectPackageManager() string {
    commands := map[string]string{
        "apt":     "which apt-get >/dev/null 2>&1 && echo 'apt'",
        "yum":     "which yum >/dev/null 2>&1 && echo 'yum'", 
        "dnf":     "which dnf >/dev/null 2>&1 && echo 'dnf'",
        "zypper":  "which zypper >/dev/null 2>&1 && echo 'zypper'",
        "pacman":  "which pacman >/dev/null 2>&1 && echo 'pacman'",
        "apk":     "which apk >/dev/null 2>&1 && echo 'apk'",
    }
    
    for pkgMgr, cmd := range commands {
        if output, err := exec.Command("sh", "-c", cmd).Output(); err == nil {
            return strings.TrimSpace(string(output))
        }
    }
    return "unknown"
}
```

#### **Service Manager Detection**
```go
// Service manager detection
func detectServiceManager() string {
    managers := []string{"systemd", "upstart", "init", "runit", "openrc"}
    
    for _, manager := range managers {
        if _, err := os.Stat(fmt.Sprintf("/proc/1/comm")); err == nil {
            if data, err := os.ReadFile("/proc/1/comm"); err == nil {
                if strings.Contains(string(data), manager) {
                    return manager
                }
            }
        }
    }
    return "unknown"
}
```

#### **SELinux Status**
```go
// SELinux status detection
func getSELinuxStatus() map[string]interface{} {
    status := map[string]interface{}{
        "enabled": false,
        "mode":    "disabled",
        "type":    "",
    }
    
    // Check if SELinux is enabled
    if _, err := os.Stat("/sys/fs/selinux"); err == nil {
        status["enabled"] = true
        
        // Get SELinux mode
        if data, err := os.ReadFile("/sys/fs/selinux/enforce"); err == nil {
            if strings.TrimSpace(string(data)) == "1" {
                status["mode"] = "enforcing"
            } else {
                status["mode"] = "permissive"
            }
        }
        
        // Get SELinux type
        if data, err := os.ReadFile("/sys/fs/selinux/policy/version"); err == nil {
            status["type"] = strings.TrimSpace(string(data))
        }
    }
    
    return status
}
```

#### **SSH Host Keys**
```go
// SSH host key collection
func getSSHHostKeys() map[string]string {
    keys := make(map[string]string)
    keyFiles := []string{
        "/etc/ssh/ssh_host_rsa_key.pub",
        "/etc/ssh/ssh_host_ecdsa_key.pub", 
        "/etc/ssh/ssh_host_ed25519_key.pub",
    }
    
    for _, keyFile := range keyFiles {
        if data, err := os.ReadFile(keyFile); err == nil {
            keyType := strings.TrimSuffix(filepath.Base(keyFile), ".pub")
            keys[keyType] = strings.TrimSpace(string(data))
        }
    }
    
    return keys
}
```

#### **BIOS Information**
```go
// BIOS information collection
func getBIOSInfo() map[string]string {
    bios := make(map[string]string)
    
    // Read from /sys/class/dmi/id/
    dmiFiles := map[string]string{
        "bios_vendor":   "/sys/class/dmi/id/bios_vendor",
        "bios_version":  "/sys/class/dmi/id/bios_version", 
        "bios_date":     "/sys/class/dmi/id/bios_date",
        "bios_release":  "/sys/class/dmi/id/bios_release",
        "board_vendor":  "/sys/class/dmi/id/board_vendor",
        "board_name":    "/sys/class/dmi/id/board_name",
        "board_version": "/sys/class/dmi/id/board_version",
    }
    
    for key, file := range dmiFiles {
        if data, err := os.ReadFile(file); err == nil {
            bios[key] = strings.TrimSpace(string(data))
        }
    }
    
    return bios
}
```

### **4. Fact Access Patterns**

```hcl
# Project facts are available in the project context
action "use-hardware" {
  command = "setup.sh {{.facts.server1.cpu.cores}} {{.facts.server1.memory.total}}"
}

# Project-specific facts are available
action "deploy-app" {
  command = "deploy.sh {{.facts.server1.app_version}}"
}

# Facts are available in templates
template "config.tmpl" {
  # Project facts available in project context
  # {{.facts.server1.os}} - available in project
  # {{.facts.server1.cpu.cores}} - available in project
  
  # Project facts in project context
  # {{.facts.server1.app_version}} - project-specific
}
```

### **5. Benefits of This Approach**

1. **Session-Based**: Facts stored in memory for the duration of action runs
2. **Clear Boundaries**: Facts are ephemeral and isolated to each run
3. **Performance**: No disk I/O overhead, fast fact access
4. **Simplicity**: No database management, clear access patterns
5. **Fresh Data**: Facts are always fresh since they're collected before each run
6. **Memory Efficiency**: Optimized memory usage for fact storage

### **6. Memory Types**

#### **In-Memory Storage (Default)**
- **Format**: Go structs in memory
- **Performance**: High performance for all deployments
- **Use Case**: All deployments, from small to large server fleets
- **Location**: Memory during action runs

#### **Memory Pooling**
- **Format**: Reused memory structures
- **Performance**: Optimized memory allocation
- **Use Case**: Large deployments with many machines
- **Location**: Memory pool during action runs

#### **Memory-Efficient Storage**
- **Format**: Compressed memory structures
- **Performance**: Reduced memory usage
- **Use Case**: Memory-constrained environments
- **Location**: Memory during action runs

## Custom Facts Integration

### **1. Ansible-Inspired Custom Facts Approach**

Following the established pattern from Ansible, Spooky implements custom facts using a node-local approach:

**Location:**  
Place custom fact files (HCL) in `/etc/spooky/facts/` on the target machine.

**Format:**  
```hcl
# /etc/spooky/facts/app.fact
app_name = "nginx"
app_version = "1.18.0"
config_path = "/etc/nginx/nginx.conf"
log_path = "/var/log/nginx"

# /etc/spooky/facts/environment.fact
environment = "production"
datacenter = "fra00"
rack = "A01"
power_zone = "PZ-1"

# /etc/spooky/facts/monitoring.fact
prometheus_port = 9100
alert_manager = "alert.example.com"
log_level = "info"
```

**Discovery:**  
When you run `spooky facts gather`, Spooky collects facts using `github.com/shirou/gopsutil/v4` and custom facts from `/etc/spooky/facts/`.

**Access in actions and templates:**  
Custom facts are available as `{{.facts.custom.<filename>.<key>}}`:

```hcl
# In templates
server {
  listen 80;
  server_name {{.facts.hostname}};
  
  # Access custom facts
  root {{.facts.custom.app.config_path}};
  access_log {{.facts.custom.app.log_path}}/access.log;
  
  # Environment-specific config
  location /health {
    return 200 "healthy from {{.facts.custom.environment.datacenter}}-{{.facts.custom.environment.rack}}";
  }
}

# In actions
action "deploy-app" {
  command = "deploy.sh {{.facts.custom.app.app_name}} {{.facts.custom.app.app_version}}"
}

action "configure-monitoring" {
  command = "setup-monitoring.sh --port {{.facts.custom.monitoring.prometheus_port}} --alert-manager {{.facts.custom.monitoring.alert_manager}}"
}
```


    
    return facts
}
```

### **5. CLI Examples**

```bash
# Gather facts including custom facts
$ spooky facts gather ./my-project
Gathering facts from 3 machines...
✓ web-001: 15 system facts + 3 custom facts
✓ web-002: 15 system facts + 3 custom facts  
✓ db-001: 15 system facts + 2 custom facts

# List facts including custom facts
$ spooky facts list ./my-project
Machine: web-001
  System Facts:
    Hostname: web-001.example.com
    OS: Ubuntu 22.04 LTS
    CPU: 4 cores, Intel(R) Core(TM) i7-8550U
    Memory: 16384 MB total, 8192 MB available
  
  Custom Facts:
    app: nginx, 1.18.0, /etc/nginx/nginx.conf
    environment: production, fra00, A01
    monitoring: prometheus:9100, alert.example.com

# Export facts including custom facts
$ spooky facts export ./my-project --format json --output facts.json
Exported 3 machines with system and custom facts to facts.json
```

### **6. Benefits of This Approach**

1. **✅ Follows Ansible Pattern** - Familiar to users coming from Ansible
2. **✅ Node-Local Facts** - Facts live on the target system where they belong
3. **✅ Runtime Discovery** - Facts are discovered when `spooky facts gather` runs
4. **✅ Simple Structure** - Just HCL files in a directory

6. **✅ Clear Namespace** - Custom facts under `.facts.custom.*`
7. **✅ No Import Step** - Facts are automatically discovered and available

## Collision Handling

### **1. Machine ID Collision Detection**

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

### **4. Collision Detection Algorithm**

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



## CLI Integration

### Noun-Verb Command Structure

Following the `cli-system.md` design and the architectural patterns established in [Spooky Design](../spooky-design.md), facts management uses the `spooky facts` noun:

```bash
# Core facts operations
spooky facts gather <project directory>     # Collect facts from machines
spooky facts list <project directory>       # List available facts
spooky facts validate <project directory>   # Validate facts data
spooky facts export <project directory>     # Export facts to various formats


```

### Project Integration

#### **Project Initialization**
Following the project system architecture described in [Spooky Design](../spooky-design.md) and implemented in [Project System](../project-system.md), project initialization creates project structure:

```bash
# Creates project structure during init
spooky project init my-project

# Results in:
my-project/
├── project.hcl
├── machines.hcl
├── actions.hcl
├── templates/
├── data/
└── logs/
# Facts are collected and stored in memory during action runs
```

#### **Facts Memory Management**
- **Automatic Management**: `spooky facts gather` manages memory allocation
- **Schema Validation**: Facts are validated against facts schema
- **Self-Check**: Fact collection includes validation
- **Error Handling**: Clear error messages if memory allocation fails

### Memory Configuration

#### **Project Configuration** (`project.hcl`)
```hcl
project {
  name = "my-web-app"
  version = "1.0.0"
  environment = "production"
  
  facts {
    memory_limit = "1GB"  # memory limit for facts
    memory_efficient = true  # use memory-efficient storage
    parallel_workers = 10  # default parallelism for fact gathering
  }
}
```

#### **Global Configuration** (`$XDG_CONFIG_HOME/spooky/spooky.hcl`)
Following the configuration system architecture described in [Spooky Design](../spooky-design.md), global configuration provides system-wide facts memory settings:

```hcl
spooky {
  facts {
    memory_limit = "2GB"  # global memory limit for facts
    memory_efficient = true  # use memory-efficient storage
    collision_policy = "warn"  # update, warn, skip, merge
  }
}
```

#### **Environment Variables**
```bash
# Override memory limit
export SPOOKY_FACTS_MEMORY_LIMIT="1GB"

# Override memory efficiency
export SPOOKY_FACTS_MEMORY_EFFICIENT="true"

# Override collision policy
export SPOOKY_COLLISION_POLICY="warn"
```



## Implementation Phases

### **Phase 1: In-Memory Facts Foundation (Week 1)**
**Goal**: Create in-memory facts system with automatic management

#### Tasks:
1. **Remove BadgerDB dependency**
   ```bash
   go mod tidy  # Remove unused dependencies
   ```

2. **Create in-memory facts storage** (`internal/facts/memory_storage.go`)
   - Implement `NewMemoryFactStorage()` function
   - Add memory allocation management
   - Create efficient memory structures

3. **Update fact manager** (`internal/facts/manager.go`)
   - Add in-memory storage integration
   - Implement fact storage in memory
   - Add fact retrieval from memory

4. **Schema validation integration**
   - Integrate facts schema validation with memory storage
   - Add self-checking fact validation
   - Implement clear error handling

### **Phase 2: Enhanced Fact Collection (Week 2)**
**Goal**: Upgrade to gopsutil v4 and add additional fact collectors

#### Tasks:
1. **Upgrade gopsutil**
   ```bash
   go get github.com/shirou/gopsutil/v4
   ```

2. **Add enhanced collectors** (`internal/facts/enhanced_collectors.go`)
   - Package manager detection
   - Service manager detection
   - SELinux status
   - SSH host keys
   - BIOS information

3. **Integrate with project storage**
   - Store all enhanced facts in project database
   - Make facts available within project context

4. **Enhanced fact collection**
   - Add priority-based collision resolution
   - Create fact merging strategies

### **Phase 3: Enhanced In-Memory Facts Storage (Week 3)**
**Goal**: Add enhanced in-memory fact storage

#### Tasks:
1. **Enhanced in-memory facts storage** (`internal/facts/memory_storage.go`)
   - Implement enhanced in-memory fact storage
   - Add memory optimization
   - Create efficient fact access patterns

2. **Fact organization** (`internal/facts/organizer.go`)
   - Organize facts efficiently in memory
   - Implement clear access patterns
   - Add fact categorization rules

3. **Collision detection and resolution**
   - Implement machine ID collision detection
   - Add confidence-based collision handling
   - Create collision resolution strategies



### **Phase 4: CLI Integration and Management (Week 3)**
**Goal**: Add comprehensive facts management to CLI

#### Tasks:
1. **Add facts commands** (`internal/cli/facts.go`)
   - `spooky facts list` - List facts in project database
   - `spooky facts show <machine>` - Show facts for specific machine
   - `spooky facts export` - Export facts to HCL/JSON
   - `spooky facts import` - Import facts from HCL/JSON
   - `spooky facts gather` - gather facts from the machines and save to facts.db


2. **Configuration integration**
   - Use global configuration for storage settings
   - Apply facts collection settings from config
   - Implement environment variable support

3. **Enhanced CLI examples**
   - Add comprehensive usage examples
   - Include collision scenarios
   - Provide configuration examples

### **Phase 5: Testing and Documentation (Week 4)**
**Goal**: Comprehensive testing and documentation

#### Tasks:
1. **Unit tests**
   - Global storage tests
   - Enhanced collector tests
   - Fact merging tests
   - Collision detection tests

2. **Integration tests**
   - End-to-end fact collection and storage
   - Cross-project fact sharing
   - Performance benchmarks
   - CLI command testing

3. **Documentation**
   - Project facts database usage
   - Enhanced fact collection
   - CLI usage examples
   - Configuration examples
   - Troubleshooting guide

## Current Implementation Status

### ✅ **Completed Features**

#### **1. Memory Storage Interface**
```go
type FactStorage interface {
    GetMachineFacts(machineID string) (*MachineFacts, error)
    SetMachineFacts(machineID string, facts *MachineFacts) error
    QueryFacts(query *FactQuery) ([]*MachineFacts, error)
    DeleteFacts(query *FactQuery) (int, error)
    DeleteMachineFacts(machineID string) error
    ExportToJSON(w io.Writer) error
    ImportFromJSON(r io.Reader) error
    Close() error
}
```

#### **2. Memory Storage Factory**
```go
// Clear separation of concerns
func NewFactStorage(opts StorageOptions) (FactStorage, error)  // For creating new memory storage
func OpenFactStorage(opts StorageOptions) (FactStorage, error) // For reading existing memory storage
```

#### **3. In-Memory Implementation**
- ✅ Full in-memory storage backend
- ✅ Memory-efficient operations
- ✅ Fast querying
- ✅ Export/import functionality
- ✅ Memory-optimized operations





#### **6. CLI Integration**
- ✅ `spooky facts gather` - Collect facts from machines
- ✅ `spooky facts list` - List available facts
- ✅ `spooky facts validate` - Validate facts data
- ✅ `spooky facts export` - Export facts to JSON/HCL

#### **7. Project Integration**
- ✅ `spooky project init` creates project structure
- ✅ Schema validation during fact collection
- ✅ Self-checking fact validation
- ✅ Clear error handling

#### **8. Export Safety**
- ✅ Export only works with existing facts in memory
- ✅ No database creation during export
- ✅ Clear error messages for missing facts

### 🔄 **In Progress Features**

#### **1. In-Memory Facts Database**
- 🔄 In-memory facts storage initialization
- 🔄 In-memory fact storage
- 🔄 Memory allocation compliance

#### **2. Enhanced Fact Collection**
- 🔄 Package manager detection
- 🔄 Service manager detection
- 🔄 SELinux status collection
- 🔄 SSH host key collection
- 🔄 BIOS information collection

#### **3. Dynamic Facts System**
- 🔄 Fact merging strategies

### ❌ **Planned Features**

#### **1. Collision Detection and Resolution**
- ❌ Machine ID collision detection logic
- ❌ Collision resolution strategies
- ❌ Confidence-based collision handling



#### **3. Advanced CLI Commands**
- ❌ `spooky facts delete` - Delete specific facts
- ❌ `spooky facts collisions` - Collision management

#### **4. Encryption Support**
- ❌ Age encryption for sensitive data
- ❌ Field-level encryption
- ❌ File-level encryption

#### **5. Advanced Memory Management**
- ❌ Memory pooling for large deployments
- ❌ Memory compression for efficiency
- ❌ Memory monitoring and profiling

#### **6. Statistics and Monitoring**
- ❌ Memory usage statistics
- ❌ Performance metrics
- ❌ Health monitoring

## Technical Specifications

### Schema Files

The fact storage system uses schema files to validate data across different formats:

#### **Facts Schema** (`internal/schemas/schemas/facts-structure.hcl`)
- **Purpose**: Validates facts structure and format
- **Location**: Embedded in binary
- **Validation**: Checks for facts structure compliance



### Schema-Based Data Structure

The schema file defines the core structure for facts data:

```hcl
# Core facts structure (from facts-badger.hcl)
facts {
  # Machine ID (32-character hex string from /etc/machine-id)
  machine_id = "a1b2c3d4e5f678901234567890123456"
  
  # Collection timestamp (ISO 8601 format)
  collected_at = "2024-07-28T12:00:00Z"
  

  
  # Collection of facts
  facts {
    # System-level facts
    system {
      # Operating system facts
      os {
        name = "Ubuntu"
        version = "22.04 LTS"
        arch = "x86_64"
        kernel = "5.15.0-107-generic"
      }
      
      # Hardware facts
      hardware {
        # CPU information
        cpu {
          cores = 4
          model = "Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz"
          frequency = 1800
        }
        
        # Memory information
        memory {
          total = 16777216      # Total memory in bytes
          available = 8388608   # Available memory in bytes
        }
        
        # Disk information
        disks = [
          {
            device = "/dev/sda1"
            mount_point = "/"
            total = 256000000000
            used = 128000000000
          }
        ]
      }
      
      # Network facts
      network {
        hostname = "web-server-01"
        interfaces = [
          {
            name = "eth0"
            mac_address = "00:11:22:33:44:55"
            ip_addresses = ["192.168.1.100"]
          }
        ]
        ip_addresses = ["192.168.1.100", "10.0.0.100"]
      }
    }
    
    # Custom facts (user-defined)
    custom {
      app_version = "1.2.3"
      deployment_status = "completed"
      team = "web-team"
      environment = "production"
    }
    
    # Environment facts
    environment {
      SPOOKY_PROJECT = "web-app"
      SPOOKY_ENVIRONMENT = "production"
    }
  }
}
```

### Storage Interface

The Go implementation provides a clean interface that maps to the schema structure:

```go
type MachineFacts struct {
    MachineID   string                 `json:"machine_id"`   // 32-char hex from /etc/machine-id
    CollectedAt time.Time              `json:"collected_at"` // ISO 8601 timestamp

    Facts       map[string]interface{} `json:"facts"`        // Schema-compliant facts structure
}

type FactQuery struct {
    MachineID     string            // Query by machine ID
    Tags          map[string]string // Query by custom facts tags
    OS            string            // Query by system.os.name
    Environment   string            // Query by custom.environment
    Limit         int               // Limit results
    SearchQuery   string            // Text search in facts
    SearchField   string            // Specific field to search
    UpdatedBefore *time.Time        // Filter by collection time
    UpdatedAfter  *time.Time        // Filter by collection time
}
```

### Schema Validation

The system validates facts data against the appropriate schema:

```go
// Validate facts against schema
func (sv *SchemaValidator) ValidateFacts(factsPath string, format string) *ValidationResult {
    switch format {
    case "badgerdb":
        return sv.validateBadgerDB(factsPath)
    case "json":
        return sv.validateJSON(factsPath)
    case "hcl":
        return sv.validateHCL(factsPath)
    default:
        return &ValidationResult{Valid: false, Errors: []ValidationError{
            {Field: "format", Message: fmt.Sprintf("unsupported format: %s", format)},
        }}
    }
}
```

### Machine ID Strategy

The storage system uses the actual machine ID (from `/etc/machine-id`) as the primary key for facts storage. This approach provides:

1. **Global Uniqueness**: Machine IDs are globally unique across all systems
2. **Automatic Deduplication**: Same physical machine = same UUID regardless of action file
3. **No Manual Management**: No need to generate or track UUIDs manually
4. **Real-world Mapping**: Direct correlation to actual hardware
5. **Schema Compliance**: Matches the `machine_id` field pattern `^[a-f0-9]{32}$`

### Data Segmentation and Memory Location

#### **In-Memory Storage Structure**
```
my-project/
├── project.hcl
├── machines.hcl
├── actions.hcl
├── templates/
├── data/
└── logs/
# Facts stored in memory during action runs
```

#### **Global Configuration Structure**
```
$HOME/.local/share/spooky/
└── config/                      # Spooky configuration
    └── spooky.hcl
# Facts stored in memory during action runs
```

#### **In-Memory Key-Value Structure**
```
├── "6538b193d562410ab480b1ce22469fe6" → {
│     machine_id: "6538b193d562410ab480b1ce22469fe6",
│     collected_at: "2024-07-28T12:00:00Z",
│     facts: {
│       system: {
│         os: { name: "Ubuntu", version: "22.04 LTS" },
│         hardware: { cpu: { cores: 4 }, memory: { total: 16777216 } },
│         network: { hostname: "web1.example.com", ip_addresses: ["192.168.1.10"] }
│       },
│       custom: { team: "web-team", environment: "production" }
│     }
│   }
├── "a1b2c3d4e5f6789012345678901234567" → {
│     machine_id: "a1b2c3d4e5f6789012345678901234567", 
│     collected_at: "2024-07-28T12:00:00Z",
│     facts: {
│       system: {
│         os: { name: "CentOS", version: "8" },
│         hardware: { cpu: { cores: 8 }, memory: { total: 33554432 } },
│         network: { hostname: "db1.example.com", ip_addresses: ["192.168.1.20"] }
│       },
│       custom: { team: "db-team", environment: "production" }
│     }
│   }
```

## CLI Examples

### **In-Memory Facts Operations**
```bash
# Initialize project structure
spooky project init my-project

# Gather facts from all machines in project
spooky facts gather ./my-project

# Gather facts from specific machines
spooky facts gather ./my-project --machine "web-001"

# Gather facts with specific tags
spooky facts gather ./my-project --tags "production"

# List all facts in memory
spooky facts list ./my-project

# List facts for specific machine
spooky facts list ./my-project --machine "web-001"

# Validate facts data
spooky facts validate ./my-project

# Export facts to JSON
spooky facts export ./my-project --format json --output facts.json

# Export facts to HCL
spooky facts export ./my-project --format hcl --output facts.hcl

# Export specific facts
spooky facts export ./my-project --machine "web-001" --format json --output web-001-facts.json
```



### **Memory Configuration Examples**
```bash
# Use custom memory limit
export SPOOKY_FACTS_MEMORY_LIMIT="1GB"
spooky facts gather ./my-project
```

## Benefits

### **Immediate Benefits**
- **Session-based fact storage** during spooky action runs
- **Automatic deduplication** using machine IDs
- **In-memory fact storage** with fast access
- **Efficient querying** with memory backend
- **Data portability** with JSON export/import
- **Session isolation** for application-specific facts
- **Fresh fact collection** for system-level information
- **Enhanced fact collection** with package managers, service managers, etc.

### **Long-term Benefits**
- **Scalable fact management** for large infrastructures
- **Real-time fact tracking** with timestamps
- **Team collaboration** with shared fact collection
- **Automated collision detection** for cloned systems
- **Integration with CI/CD** for fact validation




### **Architectural Benefits**
- **Session-focused approach** solves the "Ansible variable passing problem"
- **Clear boundaries** for session facts
- **Efficient memory usage** of system information
- **Performance optimization** through intelligent memory management
- **Flexible memory** with multiple optimization options
- **Schema validation** ensures data integrity
- **Collision resolution** handles complex scenarios

## Success Metrics

### **Functionality Metrics**
- [ ] Facts can be stored and retrieved from memory
- [ ] Facts can be exported to JSON files
- [ ] Facts can be exported to HCL files
- [ ] Machine ID collision detection works correctly
- [ ] Collision resolution strategies function properly
- [ ] Fact querying with filters works as expected
- [ ] Export/import functionality preserves data integrity
- [ ] XDG Base Directory compliance is maintained
- [ ] CLI commands provide clear feedback and error messages
- [ ] Facts are stored in memory per session
- [ ] Facts are isolated per session




### **Performance Metrics**
- [ ] In-memory storage handles 10,000+ machine facts efficiently

- [ ] Query performance scales linearly with dataset size
- [ ] Memory usage stays reasonable for large fact collections
- [ ] Export/import operations complete within acceptable timeframes

- [ ] In-memory fact storage reduces duplicate collection by 90%

### **User Experience Metrics**
- [ ] Memory configuration is intuitive and follows XDG standards
- [ ] Collision warnings are clear and actionable
- [ ] CLI commands provide helpful usage information
- [ ] Error messages are descriptive and actionable
- [ ] Documentation is comprehensive and up-to-date


- [ ] Fact access patterns are intuitive in templates

### **Integration Metrics**
- [ ] Facts memory management integrates seamlessly with project initialization
- [ ] Schema validation works during fact collection
- [ ] Export safety prevents accidental memory allocation
- [ ] Facts are available in session context
- [ ] Facts are properly isolated per session
- [ ] CLI integration follows noun-verb patterns
- [ ] Configuration precedence works correctly
- [ ] Environment variable overrides function properly

## Risk Assessment

### **Technical Risks**
- **Memory corruption** - Mitigation: Memory validation and error handling
- **Machine ID collisions** - Mitigation: Robust detection and resolution
- **Performance degradation** - Mitigation: Efficient algorithms and memory management
- **Data loss** - Mitigation: Memory protection and error handling
- **Memory usage** - Mitigation: Efficient memory implementation
- **Schema validation complexity** - Mitigation: Clear validation rules

### **User Experience Risks**
- **Configuration complexity** - Mitigation: Sensible defaults and clear documentation
- **Collision confusion** - Mitigation: Clear warnings and resolution options
- **Memory usage confusion** - Mitigation: Clear memory monitoring and limits



### **Implementation Risks**
- **Scope creep** - Mitigation: Focus on core functionality first
- **Testing complexity** - Mitigation: Comprehensive test coverage
- **Integration issues** - Mitigation: Incremental implementation and testing
- **Performance bottlenecks** - Mitigation: Benchmarking and optimization
- **Schema evolution** - Mitigation: Versioned schemas and migration tools

## Out of Scope

- **External databases**: No PostgreSQL, MySQL, etc.
- **Persistent storage**: No disk-based fact storage
- **Real-time sync**: No live synchronization between sessions
- **Custom protocols**: No custom fact collection protocols
- **External APIs**: No external fact collection services

## Timeline

- **✅ Week 1**: Foundation and storage interface (COMPLETED)
- **✅ Week 2**: BadgerDB, JSON, and HCL implementations (COMPLETED)
- **✅ Week 3**: Integration and CLI updates (COMPLETED)
- **🔄 Week 4**: Advanced features (collision detection, encryption)

## Development Instructions

### For AI Assistants and Developers

This section provides specific, actionable instructions for extending the fact storage system, following the implementation patterns and architectural principles established in [Spooky Design](../spooky-design.md).

#### **Adding New Storage Backends**

To add a new storage backend (e.g., SQLite):

1. **Implement the interface**:
```go
type SQLiteFactStorage struct {
    db *sql.DB
}

func (s *SQLiteFactStorage) GetMachineFacts(machineID string) (*MachineFacts, error) {
    // Implementation
}

// Implement all other interface methods...
```

2. **Add to storage factory**:
```go
func NewFactStorage(opts StorageOptions) (FactStorage, error) {
    switch opts.Type {
    case StorageTypeSQLite:
        return NewSQLiteFactStorage(opts.Path)
    // ... existing cases
    }
}
```

3. **Add storage type constant**:
```go
const (
    StorageTypeBadger StorageType = "badger"
    StorageTypeJSON   StorageType = "json"
    StorageTypeHCL    StorageType = "hcl"
    StorageTypeSQLite StorageType = "sqlite"  // New
)
```

#### **Adding New CLI Commands**

To add a new facts command (e.g., `spooky facts delete`):

1. **Add command to facts.go**:
```go
var factsDeleteCmd = &cobra.Command{
    Use:   "delete",
    Short: "Delete facts from storage",
    RunE:  runFactsDelete,
}

func initFactsCommands() {
    FactsCmd.AddCommand(factsDeleteCmd)
    // ... existing commands
}
```

2. **Implement the command**:
```go
func runFactsDelete(cmd *cobra.Command, args []string) error {
    // Implementation using OpenFactStorage for read-only operations
    // or NewFactStorage for write operations
}
```

#### **Adding Collision Detection**

1. **Create collision detection logic**:
```go
func DetectCollision(storage FactStorage, newFacts *MachineFacts) (*CollisionResolution, error) {
    // Implementation from the collision detection section
}
```

2. **Integrate with fact gathering**:
```go
func (m *Manager) PersistFacts(server string, collection *FactCollection) error {
    // Detect collision before storing
    resolution, err := DetectCollision(m.storage, serverFacts)
    if err != nil {
        return err
    }
    
    // Handle collision based on configuration
    // ... implementation
}
```





### **5. Template Integration with Data Variables**

#### **Enhanced Template Context**

Custom facts are integrated into the existing template system alongside `data/` variables, following the template system architecture described in [Spooky Design](../spooky-design.md) and the variable management patterns in [Variables System](../variables-system.md). This provides a unified context for all project-specific data:

```go
// Enhanced template context that includes both data/ variables and custom facts
func buildTemplateContext(server string, manager *facts.Manager, projectPath string) (map[string]interface{}, error) {
    context := make(map[string]interface{})
    
    // 1. Load data/ variables (existing functionality)
    if dataVars, err := loadDataVariables(projectPath); err == nil {
        context["data"] = dataVars
    }
    
    // 2. Load project facts (system facts from project storage)
    if projectFacts, err := manager.GetProjectFacts(server); err == nil {
        context["facts"] = projectFacts
    }
    
    // 3. Load custom facts (project-specific facts)
    if customFacts, err := manager.GetCustomFacts(server); err == nil {
        context["custom"] = customFacts
    }
    
    // 4. Load fact overrides (project-specific overrides)
    if overrides, err := manager.GetFactOverrides(server); err == nil {
        context["overrides"] = overrides
    }
    
    return context, nil
}

// GetCustomFacts retrieves custom facts for template usage
func (m *Manager) GetCustomFacts(server string) (map[string]interface{}, error) {
    if m.storage == nil {
        return nil, fmt.Errorf("no storage configured")
    }
    
    // Load persisted facts from project storage
    collection, err := m.LoadProjectFacts(server)
    if err != nil {
        return nil, err
    }
    
    // Extract custom facts
    customFacts := make(map[string]interface{})
    for key, fact := range collection.Facts {
        if strings.HasPrefix(key, "custom.") {
            parts := strings.Split(key, ".")
            if len(parts) >= 3 {
                category := parts[1]
                factKey := parts[2]
                
                if customFacts[category] == nil {
                    customFacts[category] = make(map[string]interface{})
                }
                
                if categoryMap, ok := customFacts[category].(map[string]interface{}); ok {
                    categoryMap[factKey] = fact.Value
                }
            }
        }
    }
    
    return customFacts, nil
}

// GetFactOverrides retrieves fact overrides for template usage
func (m *Manager) GetFactOverrides(server string) (map[string]interface{}, error) {
    if m.storage == nil {
        return nil, fmt.Errorf("no storage configured")
    }
    
    // Load persisted facts from project storage
    collection, err := m.LoadProjectFacts(server)
    if err != nil {
        return nil, err
    }
    
    // Extract overrides
    overrides := make(map[string]interface{})
    for key, fact := range collection.Facts {
        if strings.HasPrefix(key, "override.") {
            parts := strings.Split(key, ".")
            if len(parts) >= 3 {
                category := parts[1]
                factKey := parts[2]
                
                if overrides[category] == nil {
                    overrides[category] = make(map[string]interface{})
                }
                
                if categoryMap, ok := overrides[category].(map[string]interface{}); ok {
                    categoryMap[factKey] = fact.Value
                }
            }
        }
    }
    
    return overrides, nil
}
```