# Facts System: Comprehensive Implementation Plan

## Overview

This document is the authoritative source for all facts system implementation details in spooky. It covers facts collection, storage, validation, caching, and integration with the actions and template systems.

**Schema Integration**: This facts system implements the schema validation patterns and data structures defined in [Schema System](../schema-system.md) for comprehensive facts validation, storage format consistency, and schema-based data integrity.

**Architecture Integration**: Facts integrate with the overall spooky architecture as described in [Spooky Design](../spooky-design.md), providing machine data to the actions system for dynamic running and template rendering.

## System Integration

This facts system integrates with other core Spooky systems to provide comprehensive fact management and data sharing:

### **Project System Integration**
- **Project Initialization**: Facts database creation during `spooky project init` (see [Project System](../project-system.md))
- **Project Isolation**: Project-specific facts storage with global facts sharing
- **Project Configuration**: Facts storage settings in `project.hcl` configuration
- **Project Context**: Facts available in project run context

### **Configuration System Integration**
- **Global Configuration**: Facts collection settings from `$XDG_CONFIG_HOME/spooky/spooky.hcl` (see [Configuration System](../configuration-system.md))
- **Facts Settings**: Collection timeouts, caching TTL, parallel collection, retry attempts
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

## Storage Architecture

### **1. Hybrid Approach: Global + Project Facts**

Based on our analysis and the architectural principles outlined in [Spooky Design](../spooky-design.md), we've decided on a hybrid approach that combines the benefits of global fact sharing with project isolation:

```
$XDG_STATE_HOME/spooky/global-facts.db    # Global facts (gopsutil, hardware, system)
project/
├── facts.db/                                # Project-specific facts (optional)
├── project.hcl                              # Project configuration
├── machines.hcl                             # Machine inventory
├── actions.hcl                              # Action definitions
├── variables.hcl                            # Variable definitions
└── ...
```

### **2. Fact Classification**

**Global Facts** (stored in `$XDG_STATE_HOME/spooky/global-facts.db`):
- ✅ **gopsutil facts** (hardware, system, network)
- ✅ **BIOS information** (vendor, version, date)
- ✅ **SSH host keys** (public keys)
- ✅ **Package manager detection** (apt, yum, dnf, etc.)
- ✅ **Service manager detection** (systemd, upstart, etc.)
- ✅ **SELinux status** (enabled, mode, type)
- ✅ **System identification** (machine-id, hostname, IPs)

**Project Facts** (stored in `project/facts.db/` - optional):
- ❌ **Application versions** (nginx, postgres, etc.)
- ❌ **Deployment states** (installed, configured, running)
- ❌ **Configuration values** (custom app configs)
- ❌ **Project-specific variables** (environment-specific data)

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
# In any project, global facts are always available
action "use-hardware" {
  command = "setup.sh {{.facts.server1.cpu.cores}} {{.facts.server1.memory.total}}"
}

# Project-specific facts are isolated
action "deploy-app" {
  command = "deploy.sh {{.facts.server1.projects.myapp.app_version}}"
}

# Facts merge automatically in templates
template "config.tmpl" {
  # Global facts available everywhere
  # {{.facts.server1.os}} - always available
  # {{.facts.server1.cpu.cores}} - always available
  
  # Project facts only in project context
  # {{.facts.server1.projects.myapp.app_version}} - project-specific
}
```

### **5. Benefits of This Approach**

1. **No Duplication**: gopsutil facts stored once, referenced everywhere
2. **Clear Boundaries**: Global vs project facts are explicit
3. **Performance**: No re-gathering of expensive system information
4. **Simplicity**: One global database, clear access patterns
5. **Project Isolation**: Project facts don't leak between projects
6. **Cross-Project Intelligence**: Hardware facts shared across projects

### **6. Storage Types**

#### **BadgerDB (Default)**
- **Format**: Embedded key-value database
- **Performance**: High performance for large datasets
- **Use Case**: Production deployments, large server fleets
- **Location**: `<project>/facts.db/` directory

#### **JSON Storage**
- **Format**: Human-readable JSON files
- **Performance**: Good for small to medium deployments
- **Use Case**: Development, testing, small deployments
- **Location**: `<project>/facts.json` file

#### **HCL Storage**
- **Format**: HashiCorp Configuration Language
- **Performance**: Good for configuration-heavy deployments
- **Use Case**: Configuration management, infrastructure as code
- **Location**: `<project>/facts.hcl` file

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

### **2. Dynamic Facts**

You can also use runnable scripts in `/etc/spooky/facts/` to provide facts at runtime:

```bash
#!/bin/bash
# /etc/spooky/facts/deployment.fact (runnable)

# Get current deployment status
DEPLOYMENT_STATE=$(systemctl is-active nginx)
LAST_DEPLOY=$(stat -c %Y /var/www/html/index.html)

# Output as HCL
cat << EOF
deployment_state = "${DEPLOYMENT_STATE}"
last_deploy = ${LAST_DEPLOY}
uptime = $(uptime -p | sed 's/up //')
EOF
```

**Example dynamic fact usage:**
```hcl
# In templates
{{if eq .facts.custom.deployment.deployment_state "active"}}
  # Server is running
  status = "healthy"
{{else}}
  # Server is down
  status = "unhealthy"
{{end}}

# Show last deployment time
last_updated = "{{.facts.custom.deployment.last_deploy}}"
```

### **3. Summary**

**Custom facts are files on the remote host, discovered at runtime, and merged into the fact namespace.**

### **4. Implementation Details**

#### **Fact Discovery Process**

```go
// internal/facts/local_collector.go
func (c *LocalCollector) CollectCustomFacts() (map[string]interface{}, error) {
    customFacts := make(map[string]interface{})
    factsDir := "/etc/spooky/facts"
    
    // Read all .fact files in /etc/spooky/facts/
    entries, err := os.ReadDir(factsDir)
    if err != nil {
        if os.IsNotExist(err) {
            return customFacts, nil // Directory doesn't exist, no custom facts
        }
        return nil, err
    }
    
    for _, entry := range entries {
        if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".fact") {
            continue
        }
        
        filePath := filepath.Join(factsDir, entry.Name())
        factName := strings.TrimSuffix(entry.Name(), ".fact")
        
        // Check if file is runnable
        if info, err := os.Stat(filePath); err == nil && info.Mode()&0111 != 0 {
            // Run dynamic fact script
            customFacts[factName] = c.runDynamicFact(filePath)
        } else {
            // Parse static HCL fact file
            customFacts[factName] = c.parseStaticFact(filePath)
        }
    }
    
    return customFacts, nil
}
```

#### **Template Context Integration**

```go
// internal/cli/template_context.go
func (ctx *TemplateContext) buildFactsContext() map[string]interface{} {
    facts := make(map[string]interface{})
    
    // System facts from gopsutil
    facts["hostname"] = ctx.SystemFacts.Hostname
    facts["os"] = ctx.SystemFacts.OS
    facts["cpu"] = ctx.SystemFacts.CPU
    facts["memory"] = ctx.SystemFacts.Memory
    facts["network"] = ctx.SystemFacts.Network
    
    // Custom facts from /etc/spooky/facts/
    facts["custom"] = ctx.CustomFacts
    
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
5. **✅ Dynamic Support** - Executable scripts for runtime facts
6. **✅ Clear Namespace** - Custom facts under `.facts.custom.*`
7. **✅ No Import Step** - Facts are automatically discovered and available

## Dynamic Facts and Collision Handling

### **1. Dynamic Facts Overview**

**Dynamic facts** (similar to Ansible's dynamic facts) allow importing custom fact data from external sources:

```go
// Dynamic fact sources
type DynamicFactSource struct {
    Type     string `json:"type"`      // "json", "script", "yaml", "hcl"
    Path     string `json:"path"`      // File path or script location
    Timeout  int    `json:"timeout"`   // Execution timeout in seconds
    Priority int    `json:"priority"`  // Fact precedence (higher = more important)
}

// Example dynamic fact sources
dynamic_facts = [
    {
        type = "json"
        path = "facts.d/app-versions.json"
        priority = 100
    },
    {
        type = "script" 
        path = "facts.d/deployment-state.sh"
        timeout = 30
        priority = 200
    },
    {
        type = "yaml"
        path = "facts.d/custom-facts.yaml"
        priority = 50
    }
]
```

### **2. Dynamic Facts Collision Resolution**

**Problem**: When multiple sources provide the same fact key, we need clear precedence rules.

**Solution**: Priority-based collision resolution with explicit merging strategies:

```go
type FactCollision struct {
    Key        string                 // The conflicting fact key
    Sources    []FactSource           // All sources providing this fact
    Values     []interface{}          // The conflicting values
    Priorities []int                  // Priority of each source
    Resolution CollisionResolution    // How to resolve the conflict
}

type CollisionResolution string
const (
    ResolutionHighestPriority CollisionResolution = "highest_priority"  // Use highest priority
    ResolutionMerge           CollisionResolution = "merge"            // Merge all values
    ResolutionConcatenate     CollisionResolution = "concatenate"      // Join as strings
    ResolutionAppend          CollisionResolution = "append"           // Add to array
    ResolutionWarn            CollisionResolution = "warn"             // Warn and use first
    ResolutionError           CollisionResolution = "error"            // Fail on collision
)

// Collision resolution strategy
func resolveFactCollision(collision *FactCollision) (interface{}, error) {
    switch collision.Resolution {
    case ResolutionHighestPriority:
        return selectHighestPriority(collision)
    case ResolutionMerge:
        return mergeFactValues(collision)
    case ResolutionConcatenate:
        return concatenateFactValues(collision)
    case ResolutionAppend:
        return appendFactValues(collision)
    case ResolutionWarn:
        logCollisionWarning(collision)
        return collision.Values[0], nil
    case ResolutionError:
        return nil, fmt.Errorf("fact collision for key '%s': %v", collision.Key, collision.Values)
    default:
        return collision.Values[0], nil // Default to first value
    }
}
```

### **3. Machine ID Collision Detection**

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

### **5. Fact TTL and Change Detection**

**TTL (Time To Live) Configuration**:

```hcl
# project.hcl - Project-specific TTL settings
fact_ttl = "6h"                    # Global facts expire after 6 hours
dynamic_fact_ttl = "1h"            # Dynamic facts expire after 1 hour
system_fact_ttl = "24h"            # System facts expire after 24 hours

# Per-fact-type TTL overrides
fact_ttl_overrides = {
    "cpu" = "1h"                   # CPU facts expire quickly (frequent changes)
    "memory" = "30m"               # Memory facts expire very quickly
    "disk" = "6h"                  # Disk facts expire moderately
    "network" = "2h"               # Network facts expire quickly
    "app_version" = "24h"          # App versions change slowly
    "deployment_state" = "5m"      # Deployment state changes frequently
}

# ~/.config/spooky/config.hcl - User-wide TTL settings
fact_ttl = "12h"                   # Default TTL for all projects
dynamic_fact_ttl = "2h"            # Default dynamic fact TTL
system_fact_ttl = "48h"            # Default system fact TTL
```

**What happens when facts change**:

```go
// Enhanced fact collection with change detection
func (e *EnhancedFactManager) CollectAndStoreFacts(machine string) error {
    // 1. Get existing facts for comparison
    existingFacts, _ := e.globalStorage.GetMachineFacts(machine)
    
    // 2. Collect fresh facts
    freshFacts, err := e.collectFreshFacts(machine)
    if err != nil {
        return err
    }
    
    // 3. Detect changes and handle accordingly
    changes := e.detectFactChanges(existingFacts, freshFacts)
    
    if len(changes) > 0 {
        // Log what changed
        e.logFactChanges(machine, changes)
        
        // Check if changes are significant enough to update
        if e.shouldUpdateFacts(changes) {
            // Update facts in global storage
            return e.globalStorage.SetMachineFacts(machine, freshFacts)
        } else {
            // Minor changes - extend TTL but don't update
            e.extendFactTTL(machine, existingFacts)
            return nil
        }
    }
    
    // No changes - extend TTL
    e.extendFactTTL(machine, existingFacts)
    return nil
}
```

**Example TTL scenarios**:

```bash
# Scenario 1: Facts haven't changed, TTL extended
$ spooky act
Using cached facts for server1 (TTL: 6h remaining)

# Scenario 2: Facts expired, re-collecting
$ spooky act
Facts for server1 expired (TTL: 6h), re-collecting...
Facts changed for server1:
  memory.used: 8GB → 9GB (low impact)
  uptime: 5d → 6d (low impact)
Using updated facts

# Scenario 3: Critical change detected
$ spooky act
Facts changed for server1:
  cpu.cores: 4 → 8 (high impact)
  memory.total: 16GB → 32GB (high impact)
Updating facts due to high-impact changes

# Scenario 4: Dynamic facts changed
$ spooky act
Facts changed for server1:
  projects.myapp.app_version: 1.2.0 → 1.2.1 (medium impact)
  projects.myapp.deployment_state: running → updating (medium impact)
Updating facts due to medium-impact changes
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

# Global facts operations
spooky facts gather --global                # Collect facts to global storage
spooky facts list --global                  # List facts from global storage
spooky facts export --global --format json  # Export global facts
```

### Project Integration

#### **Project Initialization**
Following the project system architecture described in [Spooky Design](../spooky-design.md) and implemented in [Project System](../project-system.md), project initialization creates project-local facts storage:

```bash
# Creates project-local facts.db during init
spooky project init my-project

# Results in:
my-project/
├── project.hcl
├── machines.hcl
├── actions.hcl
├── facts.db/           # BadgerDB database
├── templates/
├── data/
└── logs/
```

#### **Facts Database Creation**
- **Automatic Creation**: `spooky project init` creates a functioning BadgerDB
- **Schema Validation**: Database is validated against facts schema
- **Self-Check**: Project initialization includes database validation
- **Error Handling**: Clear error messages if database creation fails

### Storage Configuration

#### **Project Configuration** (`project.hcl`)
```hcl
project {
  name = "my-web-app"
  version = "1.0.0"
  environment = "production"
  
  facts_storage {
    type = "badgerdb"  # badgerdb, json, hcl
    path = "facts.db"  # relative to project root
  }
  
  parallel = 10  # default parallelism for fact gathering
}
```

#### **Global Configuration** (`$XDG_CONFIG_HOME/spooky/spooky.hcl`)
Following the configuration system architecture described in [Spooky Design](../spooky-design.md), global configuration provides system-wide facts storage settings:

```hcl
spooky {
  facts_storage {
    type = "badgerdb"  # badgerdb, json, hcl
    path = "/var/lib/spooky/facts.db"  # absolute path
  }
  
  collision_policy = "warn"  # update, warn, skip, merge
}
```

#### **Environment Variables**
```bash
# Override storage type
export SPOOKY_FACTS_FORMAT="json"

# Override storage path
export SPOOKY_FACTS_PATH="/custom/path/facts.db"

# Override collision policy
export SPOOKY_COLLISION_POLICY="warn"
```

### Enhanced CLI Commands

#### **TTL and Fact Update Commands**
```bash
# Force refresh facts (ignore TTL)
$ spooky facts refresh server1
Forcing fact refresh for server1...
Facts changed for server1:
  memory.used: 8GB → 9GB (low impact)
  uptime: 5d → 6d (low impact)
Updated facts for server1

# Check fact TTL status
$ spooky facts ttl server1
Facts for server1:
  Last updated: 2024-01-15 10:30:00
  TTL: 6h (expires: 2024-01-15 16:30:00)
  Status: Valid (2h remaining)

# Extend TTL for specific facts
$ spooky facts extend-ttl server1 --fact cpu --duration 12h
Extended TTL for server1.cpu to 12h

# Clear facts cache (force re-collection)
$ spooky facts clear server1
Cleared cached facts for server1
Next spooky act will re-collect all facts

# Show fact change history
$ spooky facts history server1 --limit 10
Fact change history for server1:
  2024-01-15 10:30:00: memory.used 8GB → 9GB (low impact)
  2024-01-15 09:15:00: app_version 1.2.0 → 1.2.1 (medium impact)
  2024-01-15 08:00:00: cpu.cores 4 → 8 (high impact)
```

#### **Configuration Commands**
```bash
# Set project TTL
$ spooky facts config --project --ttl 4h
Set project fact TTL to 4h

# Set per-fact TTL override
$ spooky facts config --project --fact cpu --ttl 30m
Set cpu fact TTL to 30m in project

# Set change impact thresholds
$ spooky facts config --project --update-on-medium
Enable auto-update for medium impact changes

# Show current configuration
$ spooky facts config --show
Fact configuration:
  Project TTL: 6h
  Dynamic fact TTL: 1h
  System fact TTL: 24h
  Update on medium changes: true
  Update on low changes: false
  TTL overrides:
    cpu: 1h
    memory: 30m
    disk: 6h
```

#### **Advanced Fact Management Commands**
```bash
# Delete facts with filtering
$ spooky facts delete --machine-name "web-001" --confirm
Found 1 facts matching the criteria:
  - web-001 (abc123) from actions/deploy-web.hcl
Successfully deleted 1 facts.

# Export facts with filtering
$ spooky facts export --machine "web-001" --format json --output web-001-facts.json
Successfully exported facts to web-001-facts.json

# Import facts from file
$ spooky facts import facts-backup.json
Successfully imported facts from facts-backup.json

# Collect facts from specific machine
$ spooky facts collect web-001 actions/deploy-web.hcl
Successfully gathered facts for web-001 (15 facts)

# Show fact statistics
$ spooky facts stats
Fact database statistics:
  Total machines: 25
  Total facts: 375
  Last updated: 2024-01-15 10:30:00
  Storage size: 2.3MB
  Average facts per machine: 15
```

#### **Dynamic Facts Management**
```bash
# List dynamic fact sources
$ spooky facts dynamic list
Dynamic fact sources:
  - app-versions.json (priority: 100, type: json)
  - deployment-state.sh (priority: 200, type: script)
  - custom-facts.yaml (priority: 50, type: yaml)

# Add new dynamic fact source
$ spooky facts dynamic add --type script --path facts.d/monitoring.sh --priority 150
Added dynamic fact source: facts.d/monitoring.sh

# Remove dynamic fact source
$ spooky facts dynamic remove facts.d/custom-facts.yaml
Removed dynamic fact source: facts.d/custom-facts.yaml

# Test dynamic fact source
$ spooky facts dynamic test facts.d/deployment-state.sh
Testing dynamic fact source: facts.d/deployment-state.sh
Output: {"deployment_state": "running", "last_deploy": "2024-01-15T10:00:00Z"}
```

#### **Collision Management**
```bash
# List machine ID collisions
$ spooky facts collisions list
Machine ID collisions:
  - abc123: web-001 (confidence: 0.8)
    Evidence: hostname mismatch, IP mismatch
    Resolution: warn

# Resolve collision manually
$ spooky facts collisions resolve abc123 --action create
Created new machine record: abc123-web-001

# Set collision policy
$ spooky facts collisions policy --set warn
Set collision policy to: warn

# Show collision history
$ spooky facts collisions history --limit 5
Collision history:
  2024-01-15 10:30:00: abc123 (web-001) - resolved by create
  2024-01-15 09:15:00: def456 (db-001) - resolved by merge
  2024-01-15 08:00:00: ghi789 (app-001) - resolved by update
```

## Implementation Phases

### **Phase 1: Global Facts Database Foundation (Week 1)**
**Goal**: Create global facts database with automatic initialization

#### Tasks:
1. **Add BadgerDB dependency**
   ```bash
   go get github.com/dgraph-io/badger/v4
   ```

2. **Create global facts storage** (`internal/facts/global_storage.go`)
   - Implement `ensureGlobalFactsDB()` function
   - Add XDG Base Directory compliance
   - Create BadgerDB initialization

3. **Update fact manager** (`internal/facts/manager.go`)
   - Add global storage integration
   - Implement fact persistence to global database
   - Add fact retrieval from global database

4. **Schema validation integration**
   - Integrate facts schema validation with global storage
   - Add self-checking database creation
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

3. **Integrate with global storage**
   - Store all enhanced facts in global database
   - Make facts available across all projects

4. **Enhanced fact collection**
   - Implement dynamic fact sources (JSON, scripts, YAML, HCL)
   - Add priority-based collision resolution
   - Create fact merging strategies

### **Phase 3: Project Facts and Hybrid Storage (Week 3)**
**Goal**: Add optional project-specific fact storage with hybrid approach

#### Tasks:
1. **Project facts storage** (`internal/facts/project_storage.go`)
   - Implement project-specific fact storage
   - Add project isolation
   - Create project fact access patterns

2. **Fact merging** (`internal/facts/merger.go`)
   - Merge global and project facts
   - Implement clear access patterns
   - Add fact precedence rules

3. **Collision detection and resolution**
   - Implement machine ID collision detection
   - Add confidence-based collision handling
   - Create collision resolution strategies

4. **TTL and change detection**
   - Implement fact TTL configuration
   - Add change impact assessment
   - Create automatic update policies

### **Phase 4: CLI Integration and Management (Week 3)**
**Goal**: Add comprehensive facts management to CLI

#### Tasks:
1. **Add facts commands** (`internal/cli/facts.go`)
   - `spooky facts list` - List facts in global database
   - `spooky facts show <machine>` - Show facts for specific machine
   - `spooky facts export` - Export facts to JSON
   - `spooky facts import` - Import facts from JSON
   - `spooky facts delete` - Delete facts with filtering
   - `spooky facts refresh` - Force refresh facts
   - `spooky facts ttl` - TTL management
   - `spooky facts config` - Configuration management
   - `spooky facts dynamic` - Dynamic facts management
   - `spooky facts collisions` - Collision management

2. **Configuration integration**
   - Use global configuration for storage settings
   - Apply facts collection settings from config
   - Implement environment variable support

3. **Enhanced CLI examples**
   - Add comprehensive usage examples
   - Include TTL and collision scenarios
   - Provide configuration examples

### **Phase 5: Testing and Documentation (Week 4)**
**Goal**: Comprehensive testing and documentation

#### Tasks:
1. **Unit tests**
   - Global storage tests
   - Enhanced collector tests
   - Fact merging tests
   - Collision detection tests
   - TTL and change detection tests

2. **Integration tests**
   - End-to-end fact collection and storage
   - Cross-project fact sharing
   - Performance benchmarks
   - CLI command testing

3. **Documentation**
   - Global facts database usage
   - Enhanced fact collection
   - CLI usage examples
   - Configuration examples
   - Troubleshooting guide

## Current Implementation Status

### ✅ **Completed Features**

#### **1. Storage Interface**
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

#### **2. Storage Factory**
```go
// Clear separation of concerns
func NewFactStorage(opts StorageOptions) (FactStorage, error)  // For creating new storage
func OpenFactStorage(opts StorageOptions) (FactStorage, error) // For reading existing storage
```

#### **3. BadgerDB Implementation**
- ✅ Full BadgerDB storage backend
- ✅ ACID transactions
- ✅ Efficient querying
- ✅ Export/import functionality
- ✅ Read-only operations for export

#### **4. JSON Storage Implementation**
- ✅ File-based JSON storage
- ✅ Thread-safe operations
- ✅ Export/import functionality
- ✅ Human-readable format

#### **5. HCL Storage Implementation**
- ✅ HCL-based storage backend
- ✅ Configuration-friendly format
- ✅ Export/import functionality

#### **6. CLI Integration**
- ✅ `spooky facts gather` - Collect facts from machines
- ✅ `spooky facts list` - List available facts
- ✅ `spooky facts validate` - Validate facts data
- ✅ `spooky facts export` - Export facts to JSON/HCL

#### **7. Project Integration**
- ✅ `spooky project init` creates facts.db
- ✅ Schema validation during initialization
- ✅ Self-checking database creation
- ✅ Clear error handling

#### **8. Export Safety**
- ✅ Export only works with existing databases
- ✅ No database creation during export
- ✅ Clear error messages for missing databases

### 🔄 **In Progress Features**

#### **1. Global Facts Database**
- 🔄 Global facts storage initialization
- 🔄 Cross-project fact sharing
- 🔄 XDG Base Directory compliance

#### **2. Enhanced Fact Collection**
- 🔄 Package manager detection
- 🔄 Service manager detection
- 🔄 SELinux status collection
- 🔄 SSH host key collection
- 🔄 BIOS information collection

#### **3. Dynamic Facts System**
- 🔄 Dynamic fact sources (JSON, scripts, YAML, HCL)
- 🔄 Priority-based collision resolution
- 🔄 Fact merging strategies

### ❌ **Planned Features**

#### **1. Collision Detection and Resolution**
- ❌ Machine ID collision detection logic
- ❌ Collision resolution strategies
- ❌ Confidence-based collision handling

#### **2. TTL and Change Detection**
- ❌ Fact TTL configuration
- ❌ Change impact assessment
- ❌ Automatic update policies

#### **3. Advanced CLI Commands**
- ❌ `spooky facts delete` - Delete specific facts
- ❌ `spooky facts refresh` - Force refresh facts
- ❌ `spooky facts ttl` - TTL management
- ❌ `spooky facts config` - Configuration management
- ❌ `spooky facts dynamic` - Dynamic facts management
- ❌ `spooky facts collisions` - Collision management

#### **4. Encryption Support**
- ❌ Age encryption for sensitive data
- ❌ Field-level encryption
- ❌ File-level encryption

#### **5. Cloud Storage Backends**
- ❌ S3/MinIO storage support
- ❌ HTTP storage endpoints
- ❌ Multi-backend synchronization

#### **6. Statistics and Monitoring**
- ❌ Database statistics
- ❌ Performance metrics
- ❌ Health monitoring

## Technical Specifications

### Schema Files

The fact storage system uses three schema files to validate data across different storage formats:

#### **BadgerDB Schema** (`internal/schemas/schemas/facts-badger.hcl`)
- **Purpose**: Validates facts stored in BadgerDB format
- **Location**: `<project>/facts.db/` directory
- **Validation**: Checks for BadgerDB files (`MANIFEST`, `KEYREGISTRY`, `.vlog`)

#### **JSON Schema** (`internal/schemas/schemas/facts-json.hcl`)
- **Purpose**: Validates facts stored in JSON format
- **Location**: `<project>/facts.json` file
- **Format**: Human-readable JSON with structured facts

#### **HCL Schema** (`internal/schemas/schemas/facts-hcl.hcl`)
- **Purpose**: Validates facts stored in HCL format
- **Location**: `<project>/facts.hcl` file
- **Format**: HashiCorp Configuration Language with facts blocks

### Schema-Based Data Structure

All three schema files define the same core structure for facts data:

```hcl
# Core facts structure (from facts-badger.hcl, facts-json.hcl, facts-hcl.hcl)
facts {
  # Machine ID (32-character hex string from /etc/machine-id)
  machine_id = "a1b2c3d4e5f678901234567890123456"
  
  # Collection timestamp (ISO 8601 format)
  collected_at = "2024-07-28T12:00:00Z"
  
  # Time-to-live (optional, e.g., "24h", "30m")
  ttl = "24h"
  
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
    TTL         string                 `json:"ttl,omitempty"` // Optional TTL
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

### Data Segmentation and Storage Location

#### **Project-Local Storage Structure**
```
my-project/
├── project.hcl
├── machines.hcl
├── actions.hcl
├── facts.db/                    # BadgerDB database directory
│   ├── 000001.vlog
│   ├── 000002.sst
│   └── MANIFEST
├── facts.json                   # JSON storage file (if using JSON backend)
├── facts.hcl                    # HCL storage file (if using HCL backend)
├── templates/
├── data/
└── logs/
```

#### **Global Storage Structure** (Optional)
```
$HOME/.local/share/spooky/
├── facts.db/                    # Global BadgerDB database
│   ├── 000001.vlog
│   ├── 000002.sst
│   └── MANIFEST
├── facts.json                   # Global JSON storage file
├── facts.hcl                    # Global HCL storage file
└── config/                      # Spooky configuration
    └── spooky.hcl
```

#### **BadgerDB Key-Value Structure**
```
├── "6538b193d562410ab480b1ce22469fe6" → {
│     machine_id: "6538b193d562410ab480b1ce22469fe6",
│     collected_at: "2024-07-28T12:00:00Z",
│     ttl: "24h",
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
│     ttl: "24h",
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

### **Project-Local Facts Operations**
```bash
# Initialize project with facts database
spooky project init my-project

# Gather facts from all machines in project
spooky facts gather ./my-project

# Gather facts from specific machines
spooky facts gather ./my-project --machine "web-001"

# Gather facts with specific tags
spooky facts gather ./my-project --tags "production"

# List all facts in project
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

### **Global Facts Operations**
```bash
# Gather facts to global storage
spooky facts gather --global

# List facts from global storage
spooky facts list --global

# Export global facts
spooky facts export --global --format json --output global-facts.json

# Validate global facts
spooky facts validate --global
```

### **Storage Configuration Examples**
```bash
# Use JSON storage for project
spooky project init my-project
# Edit project.hcl to set facts_storage.type = "json"

# Use custom storage path
export SPOOKY_FACTS_PATH="/custom/path/facts.db"
spooky facts gather ./my-project

# Use HCL storage format
export SPOOKY_FACTS_FORMAT="hcl"
spooky facts gather ./my-project
```

## Benefits

### **Immediate Benefits**
- **Persistent fact storage** across spooky sessions
- **Automatic deduplication** using machine IDs
- **Cross-project fact sharing** with portable project names
- **Efficient querying** with BadgerDB backend
- **Data portability** with JSON export/import
- **Project isolation** for application-specific facts
- **Global fact sharing** for system-level information
- **Enhanced fact collection** with package managers, service managers, etc.

### **Long-term Benefits**
- **Scalable fact management** for large infrastructures
- **Historical fact tracking** with timestamps
- **Team collaboration** with shared fact databases
- **Automated collision detection** for cloned systems
- **Integration with CI/CD** for fact validation
- **Dynamic fact sources** for custom data collection
- **TTL-based caching** for performance optimization
- **Change impact assessment** for intelligent updates

### **Architectural Benefits**
- **Hybrid approach** solves the "Ansible variable passing problem"
- **Clear boundaries** between global and project facts
- **No duplication** of expensive system information
- **Performance optimization** through intelligent caching
- **Flexible storage** with multiple backend options
- **Schema validation** ensures data integrity
- **Collision resolution** handles complex scenarios

## Success Metrics

### **Functionality Metrics**
- [ ] Facts can be stored and retrieved from BadgerDB
- [ ] Facts can be stored and retrieved from JSON files
- [ ] Facts can be stored and retrieved from HCL files
- [ ] Machine ID collision detection works correctly
- [ ] Collision resolution strategies function properly
- [ ] Fact querying with filters works as expected
- [ ] Export/import functionality preserves data integrity
- [ ] XDG Base Directory compliance is maintained
- [ ] CLI commands provide clear feedback and error messages
- [ ] Global facts are shared across projects
- [ ] Project facts are isolated per project
- [ ] Dynamic fact sources work correctly
- [ ] TTL-based caching functions properly
- [ ] Change detection and impact assessment work

### **Performance Metrics**
- [ ] BadgerDB storage handles 10,000+ machine facts efficiently
- [ ] JSON storage performs well for small to medium datasets
- [ ] HCL storage performs well for configuration-heavy deployments
- [ ] Query performance scales linearly with dataset size
- [ ] Memory usage stays reasonable for large fact collections
- [ ] Export/import operations complete within acceptable timeframes
- [ ] Fact collection with TTL caching reduces collection time by 80%
- [ ] Global fact sharing reduces duplicate collection by 90%

### **User Experience Metrics**
- [ ] Storage configuration is intuitive and follows XDG standards
- [ ] Collision warnings are clear and actionable
- [ ] CLI commands provide helpful usage information
- [ ] Error messages are descriptive and actionable
- [ ] Documentation is comprehensive and up-to-date
- [ ] TTL configuration is flexible and user-friendly
- [ ] Dynamic fact sources are easy to configure
- [ ] Fact access patterns are intuitive in templates

### **Integration Metrics**
- [ ] Facts storage integrates seamlessly with project initialization
- [ ] Schema validation works during project creation
- [ ] Export safety prevents accidental database creation
- [ ] Global facts are available in all project contexts
- [ ] Project facts are properly isolated
- [ ] CLI integration follows noun-verb patterns
- [ ] Configuration precedence works correctly
- [ ] Environment variable overrides function properly

## Risk Assessment

### **Technical Risks**
- **BadgerDB corruption** - Mitigation: Regular backups and validation
- **Machine ID collisions** - Mitigation: Robust detection and resolution
- **Performance degradation** - Mitigation: Efficient algorithms and indexing
- **Data loss** - Mitigation: Atomic operations and error handling
- **Memory usage** - Mitigation: Efficient BadgerDB implementation
- **Schema validation complexity** - Mitigation: Clear validation rules

### **User Experience Risks**
- **Configuration complexity** - Mitigation: Sensible defaults and clear documentation
- **Collision confusion** - Mitigation: Clear warnings and resolution options
- **Storage location confusion** - Mitigation: XDG compliance and clear paths
- **TTL configuration complexity** - Mitigation: Hierarchical configuration with defaults
- **Dynamic facts complexity** - Mitigation: Clear examples and documentation

### **Implementation Risks**
- **Scope creep** - Mitigation: Focus on core functionality first
- **Testing complexity** - Mitigation: Comprehensive test coverage
- **Integration issues** - Mitigation: Incremental implementation and testing
- **Performance bottlenecks** - Mitigation: Benchmarking and optimization
- **Schema evolution** - Mitigation: Versioned schemas and migration tools

## Out of Scope

- **External databases**: No PostgreSQL, MySQL, etc.
- **Distributed storage**: No clustering or replication
- **Real-time sync**: No live synchronization between storage types
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

#### **Adding TTL Support**

1. **Add TTL configuration**:
```go
type TTLConfig struct {
    DefaultTTL    time.Duration            `json:"default_ttl"`
    FactTypeTTLs  map[string]time.Duration `json:"fact_type_ttls"`
    UpdatePolicy  UpdatePolicy             `json:"update_policy"`
}

type UpdatePolicy struct {
    UpdateOnLow    bool `json:"update_on_low"`
    UpdateOnMedium bool `json:"update_on_medium"`
    UpdateOnHigh   bool `json:"update_on_high"`
}
```

2. **Implement TTL checking**:
```go
func (s *FactStorage) IsFactsValid(machineID string) (bool, time.Duration, error) {
    facts, err := s.GetMachineFacts(machineID)
    if err != nil {
        return false, 0, err
    }
    
    ttl := s.getTTLForFacts(facts)
    remaining := ttl - time.Since(facts.UpdatedAt)
    
    return remaining > 0, remaining, nil
}
```

#### **Adding Dynamic Facts**

1. **Create dynamic fact collector**:
```go
type DynamicFactCollector struct {
    sources []DynamicFactSource
    timeout time.Duration
}

func (d *DynamicFactCollector) Collect() (map[string]*Fact, error) {
    // Implementation from the dynamic facts section
}
```

2. **Integrate with fact manager**:
```go
func (m *Manager) CollectAllFacts(machine string) (*FactCollection, error) {
    // Collect system facts
    systemFacts, err := m.collectSystemFacts(machine)
    if err != nil {
        return nil, err
    }
    
    // Collect dynamic facts
    dynamicFacts, err := m.dynamicCollector.Collect()
    if err != nil {
        return nil, err
    }
    
    // Merge with collision resolution
    return m.mergeFactsWithCollisionResolution(systemFacts, dynamicFacts), nil
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
    
    // 2. Load global facts (system facts from global storage)
    if globalFacts, err := manager.GetGlobalFacts(server); err == nil {
        context["facts"] = globalFacts
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