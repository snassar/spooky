# Configuration System: Comprehensive Implementation Plan

## Overview

This document is the authoritative source for all configuration system implementation details in spooky. It covers global configuration, environment variables, precedence hierarchy, and integration with all other spooky systems.

**Schema Integration**: This configuration system follows the schema validation patterns and global configuration definitions defined in [Schema System](../schema-system.md) for comprehensive configuration validation, schema-based defaults, and schema-enforced configuration integrity.

**Architecture Integration**: Configuration integrates with the overall spooky architecture as described in [Spooky Design](../spooky-design.md), providing centralized settings and precedence management for all system components.

## System Integration

This configuration system integrates with other core Spooky systems to provide centralized settings and precedence management:

### **Actions System Integration**
- **Action Execution Settings**: Timeouts, retries, parallel execution limits (see [Actions System](../actions-system.md))
- **SSH Configuration**: SSH settings for remote action execution
- **Logging Configuration**: Action execution logging and output settings
- **Performance Settings**: Action execution performance tuning
- **Security Settings**: Action execution security policies

### **Facts System Integration**
- **Facts Collection Settings**: Collection timeouts, caching TTL, parallel collection, retry attempts (see [Facts System](../facts-system.md))
- **Storage Configuration**: BadgerDB/JSON storage settings, compression, encryption
- **Performance Settings**: Concurrent fact collectors, memory limits, cache sizes
- **Configuration Precedence**: CLI flags → Project config → Global config → Environment variables → Defaults

### **Project System Integration**
- **Project Configuration**: Project-specific configuration overrides (see [Project System](../project-system.md))
- **Project Isolation**: Project configuration isolation from global settings
- **Project Defaults**: Default configuration values for new projects
- **Project Validation**: Configuration validation for project settings
- **Project Context**: Configuration resolution within project context

### **Variables System Integration**
- **Variable Resolution**: Configuration settings for variable resolution (see [Variables System](../variables-system.md))
- **Variable Storage**: Configuration for variable storage and caching
- **Variable Security**: Configuration for sensitive variable handling
- **Variable Performance**: Configuration for variable resolution performance
- **Variable Validation**: Configuration for variable validation settings

### **CLI System Integration**
- **Configuration Commands**: `spooky config` commands for configuration management (see [CLI System](../cli-system.md))
- **Configuration Display**: `spooky config show` for configuration information
- **Configuration Validation**: `spooky config validate` for configuration file validation
- **Configuration Discovery**: `spooky config list` for configuration section discovery
- **Global Configuration**: Integration with global configuration system

### **Schema System Integration**
- **Configuration Schema**: Schema validation for global configuration structure (see [Schema System](../schema-system.md))
- **Schema Validation**: Configuration file validation against embedded schemas
- **Schema Evolution**: Configuration schema versioning and migration
- **Schema Composition**: Runtime schema composition for configuration validation
- **Schema Integration**: Configuration integration with all system schemas

### **Machines System Integration**
- **Machine Configuration**: Machine inventory settings and SSH configuration (see [Machines System](../machines-system.md))
- **SSH Settings**: Global SSH timeout, keepalive, and security settings
- **Authentication Settings**: SSH key validation and connection policies
- **Performance Settings**: Connection pooling and parallel operation limits
- **Security Settings**: Machine access control and audit logging

### **Template System Integration**
- **Template Configuration**: Template rendering settings and security policies (see [Template System](../template-system.md))
- **Template Settings**: Template caching, size limits, and timeout configuration
- **Security Settings**: Template sandboxing and function restrictions
- **Performance Settings**: Template compilation and rendering optimization
- **Development Settings**: Template debugging and development tools

## Current State Analysis

### **What We Have**
- ✅ **Project-level configuration** with `project.hcl`
- ✅ **HCL parsing and validation** infrastructure
- ✅ **XDG directory handling** in facts system

### **What We Need**
- 🔄 **Global configuration file** with HCL schema
- 🔄 **Configuration precedence** hierarchy
- 🔄 **Schema validation** for global config
- 🔄 **Default values** and sensible defaults
- 🔄 **CLI integration** for config management

## CLI Integration

This configuration system integrates with the CLI system to provide comprehensive configuration management:

### **Configuration Commands**
- **`spooky config show`** - Display current configuration with precedence hierarchy
- **`spooky config validate`** - Validate configuration file syntax and schema
- **`spooky config list`** - List all configuration sections and values

### **Configuration Examples**
```bash
# Show current configuration
spooky config show

# Validate configuration file
spooky config validate $XDG_CONFIG_HOME/spooky/spooky.hcl

# List specific configuration section
spooky config list --section "facts"
```

### **Global Flags**
- **`--config <file>`** - Use alternative configuration file
- **`--verbose`** - Show detailed configuration information
- **`--quiet`** - Suppress configuration output

### **Environment Variable Overrides**
Configuration values can be overridden via environment variables:
- **`SPOOKY_FACTS_TIMEOUT`** → `spooky.facts.collection_timeout`
- **`SPOOKY_SSH_TIMEOUT`** → `spooky.ssh.timeout`
- **`SPOOKY_LOG_LEVEL`** → `spooky.logging.level`
- **`SPOOKY_STORAGE_TYPE`** → `spooky.storage.type`
- **`SPOOKY_STORAGE_PATH`** → `spooky.storage.path`

### **Configuration Precedence**
1. **CLI flags** - `--config`, `--verbose`, etc.
2. **Environment variables** - `SPOOKY_*` variables
3. **Global configuration** - `$XDG_CONFIG_HOME/spooky/spooky.hcl`
4. **Default values** - Built-in sensible defaults

**Reference**: See [CLI System](../cli-system.md) for detailed command patterns and usage examples.

## Configuration Schema Design

### **1. HCL Schema Structure**

```hcl
# $XDG_CONFIG_HOME/spooky/spooky.hcl
spooky {
  # Storage configuration
  storage {
    type = "badger"                    # badger, json
    path = "$XDG_STATE_HOME/spooky/global-facts.db"
    compression = true                 # Enable compression for BadgerDB
    encryption = false                 # Enable encryption (future)
  }
  
  # Facts collection settings
  facts {
    collection_timeout = "30s"         # SSH timeout for fact collection
    cache_ttl = "24h"                 # How long to cache facts
    auto_collect = true               # Auto-collect facts on first access
    parallel_collection = 4           # Number of parallel fact collectors
    retry_attempts = 3                # Retry attempts for failed collection
  }
  
  # SSH configuration
  ssh {
    timeout = "30s"                   # SSH connection timeout
    keepalive_interval = "60s"        # SSH keepalive interval
    keepalive_count = 3               # SSH keepalive count
    key_scan_timeout = "10s"          # SSH key scanning timeout
    known_hosts_strict = false        # Strict known_hosts checking
  }
  
  # Logging configuration
  logging {
    level = "info"                    # debug, info, warn, error
    format = "text"                   # text, json
    output = "stderr"                 # stdout, stderr, file
    file_path = ""                    # Path to log file (if output = "file")
    include_timestamps = true         # Include timestamps in logs
    include_colors = true             # Include colors in terminal output
  }
  
  # Template configuration
  templates {
    max_size = "10MB"                 # Maximum template file size
    allow_external_functions = false  # Allow external function calls
    timeout = "30s"                   # Template rendering timeout
    cache_compiled = true             # Cache compiled templates
  }
  
  # Security settings
  security {
    allow_unsafe_commands = false     # Allow potentially unsafe commands
    restrict_file_access = true       # Restrict file access to project dirs
    validate_ssh_keys = true          # Validate SSH key fingerprints
    audit_logging = false             # Enable audit logging (future)
  }
  
  # Performance settings
  performance {
    max_concurrent_actions = 10       # Maximum concurrent action execution
    memory_limit = "1GB"              # Memory limit for spooky process
    gc_interval = "5m"                # Garbage collection interval
    cache_size = "100MB"              # In-memory cache size
  }
  
  # Development settings
  development {
    debug_mode = false                # Enable debug mode
    profile_cpu = false               # Enable CPU profiling
    profile_memory = false            # Enable memory profiling
    trace_execution = false           # Trace action execution
  }
}
```

### **2. Go Schema Definition**

```go
type GlobalConfig struct {
    Spooky SpookyConfig `hcl:"spooky,block"`
}

type SpookyConfig struct {
    Storage     StorageConfig     `hcl:"storage,block"`
    Facts       FactsConfig       `hcl:"facts,block"`
    SSH         SSHConfig         `hcl:"ssh,block"`
    Logging     LoggingConfig     `hcl:"logging,block"`
    Templates   TemplatesConfig   `hcl:"templates,block"`
    Security    SecurityConfig    `hcl:"security,block"`
    Performance PerformanceConfig `hcl:"performance,block"`
    Development DevelopmentConfig `hcl:"development,block"`
}

type StorageConfig struct {
    Type        string `hcl:"type,optional"`
    Path        string `hcl:"path,optional"`
    Compression bool   `hcl:"compression,optional"`
    Encryption  bool   `hcl:"encryption,optional"`
}

type FactsConfig struct {
    CollectionTimeout   string `hcl:"collection_timeout,optional"`
    CacheTTL           string `hcl:"cache_ttl,optional"`
    AutoCollect        bool   `hcl:"auto_collect,optional"`
    ParallelCollection int    `hcl:"parallel_collection,optional"`
    RetryAttempts      int    `hcl:"retry_attempts,optional"`
}

type SSHConfig struct {
    Timeout           string `hcl:"timeout,optional"`
    KeepaliveInterval string `hcl:"keepalive_interval,optional"`
    KeepaliveCount    int    `hcl:"keepalive_count,optional"`
    KeyScanTimeout    string `hcl:"key_scan_timeout,optional"`
    KnownHostsStrict  bool   `hcl:"known_hosts_strict,optional"`
}

type LoggingConfig struct {
    Level            string `hcl:"level,optional"`
    Format           string `hcl:"format,optional"`
    Output           string `hcl:"output,optional"`
    FilePath         string `hcl:"file_path,optional"`
    IncludeTimestamps bool  `hcl:"include_timestamps,optional"`
    IncludeColors    bool  `hcl:"include_colors,optional"`
}

type TemplatesConfig struct {
    MaxSize              string `hcl:"max_size,optional"`
    AllowExternalFunctions bool  `hcl:"allow_external_functions,optional"`
    Timeout              string `hcl:"timeout,optional"`
    CacheCompiled        bool   `hcl:"cache_compiled,optional"`
}

type SecurityConfig struct {
    AllowUnsafeCommands bool `hcl:"allow_unsafe_commands,optional"`
    RestrictFileAccess  bool `hcl:"restrict_file_access,optional"`
    ValidateSSHKeys     bool `hcl:"validate_ssh_keys,optional"`
    AuditLogging        bool `hcl:"audit_logging,optional"`
}

type PerformanceConfig struct {
    MaxConcurrentActions int    `hcl:"max_concurrent_actions,optional"`
    MemoryLimit         string `hcl:"memory_limit,optional"`
    GCInterval          string `hcl:"gc_interval,optional"`
    CacheSize           string `hcl:"cache_size,optional"`
}

type DevelopmentConfig struct {
    DebugMode      bool `hcl:"debug_mode,optional"`
    ProfileCPU     bool `hcl:"profile_cpu,optional"`
    ProfileMemory  bool `hcl:"profile_memory,optional"`
    TraceExecution bool `hcl:"trace_execution,optional"`
}
```

## Configuration Precedence

### **1. Precedence Hierarchy**

```go
// Configuration precedence (highest to lowest priority)
func loadConfiguration() (*GlobalConfig, error) {
    config := &GlobalConfig{}
    
    // 1. Default values (lowest priority)
    setDefaults(config)
    
    // 2. Global config file ($XDG_CONFIG_HOME/spooky/spooky.hcl)
    if err := loadGlobalConfig(config); err != nil {
        return nil, fmt.Errorf("failed to load global config: %w", err)
    }
    
    // 3. Environment variables
    loadEnvironmentOverrides(config)
    
    // 4. CLI flags (highest priority)
    loadCLIOverrides(config)
    
    return config, nil
}
```

### **2. Environment Variable Mapping**

```go
const (
    // Storage
    EnvSpookyStorageType = "SPOOKY_STORAGE_TYPE"
    EnvSpookyStoragePath = "SPOOKY_STORAGE_PATH"
    
    // Facts
    EnvSpookyFactsTimeout = "SPOOKY_FACTS_TIMEOUT"
    EnvSpookyFactsCacheTTL = "SPOOKY_FACTS_CACHE_TTL"
    
    // SSH
    EnvSpookySSHTimeout = "SPOOKY_SSH_TIMEOUT"
    EnvSpookySSHKeepalive = "SPOOKY_SSH_KEEPALIVE"
    
    // Logging
    EnvSpookyLogLevel = "SPOOKY_LOG_LEVEL"
    EnvSpookyLogFormat = "SPOOKY_LOG_FORMAT"
    
    // Security
    EnvSpookyAllowUnsafe = "SPOOKY_ALLOW_UNSAFE_COMMANDS"
    EnvSpookyRestrictFiles = "SPOOKY_RESTRICT_FILE_ACCESS"
)
```

### **3. CLI Flag Integration**

```go
var (
    // Storage flags
    globalStorageType = flag.String("global-storage-type", "", "Global storage type")
    globalStoragePath = flag.String("global-storage-path", "", "Global storage path")
    
    // Facts flags
    globalFactsTimeout = flag.String("global-facts-timeout", "", "Global facts collection timeout")
    globalFactsCacheTTL = flag.String("global-facts-cache-ttl", "", "Global facts cache TTL")
    
    // SSH flags
    globalSSHTimeout = flag.String("global-ssh-timeout", "", "Global SSH timeout")
    
    // Logging flags
    globalLogLevel = flag.String("global-log-level", "", "Global log level")
    globalLogFormat = flag.String("global-log-format", "", "Global log format")
)
```

## Default Values

### **1. Sensible Defaults**

```go
func setDefaults(config *GlobalConfig) {
    // Storage defaults
    config.Spooky.Storage.Type = "badger"
    config.Spooky.Storage.Path = "$XDG_STATE_HOME/spooky/global-facts.db"
    config.Spooky.Storage.Compression = true
    config.Spooky.Storage.Encryption = false
    
    // Facts defaults
    config.Spooky.Facts.CollectionTimeout = "30s"
    config.Spooky.Facts.CacheTTL = "24h"
    config.Spooky.Facts.AutoCollect = true
    config.Spooky.Facts.ParallelCollection = 4
    config.Spooky.Facts.RetryAttempts = 3
    
    // SSH defaults
    config.Spooky.SSH.Timeout = "30s"
    config.Spooky.SSH.KeepaliveInterval = "60s"
    config.Spooky.SSH.KeepaliveCount = 3
    config.Spooky.SSH.KeyScanTimeout = "10s"
    config.Spooky.SSH.KnownHostsStrict = false
    
    // Logging defaults
    config.Spooky.Logging.Level = "info"
    config.Spooky.Logging.Format = "text"
    config.Spooky.Logging.Output = "stderr"
    config.Spooky.Logging.IncludeTimestamps = true
    config.Spooky.Logging.IncludeColors = true
    
    // Template defaults
    config.Spooky.Templates.MaxSize = "10MB"
    config.Spooky.Templates.AllowExternalFunctions = false
    config.Spooky.Templates.Timeout = "30s"
    config.Spooky.Templates.CacheCompiled = true
    
    // Security defaults
    config.Spooky.Security.AllowUnsafeCommands = false
    config.Spooky.Security.RestrictFileAccess = true
    config.Spooky.Security.ValidateSSHKeys = true
    config.Spooky.Security.AuditLogging = false
    
    // Performance defaults
    config.Spooky.Performance.MaxConcurrentActions = 10
    config.Spooky.Performance.MemoryLimit = "1GB"
    config.Spooky.Performance.GCInterval = "5m"
    config.Spooky.Performance.CacheSize = "100MB"
    
    // Development defaults
    config.Spooky.Development.DebugMode = false
    config.Spooky.Development.ProfileCPU = false
    config.Spooky.Development.ProfileMemory = false
    config.Spooky.Development.TraceExecution = false
}
```

## Implementation Phases

### **Phase 1: Schema Definition (Week 1)**
**Goal**: Define HCL schema and Go structs

#### Tasks:
1. **Create schema definitions** (`internal/config/schema.go`)
   - Define all configuration structs
   - Add HCL tags and validation rules
   - Create default value functions

2. **Add validation rules** (`internal/config/validation.go`)
   - Validate configuration values
   - Check for conflicting settings
   - Ensure required dependencies

### **Phase 2: Configuration Loading (Week 1)**
**Goal**: Implement configuration loading with precedence

#### Tasks:
1. **Configuration loader** (`internal/config/loader.go`)
   - Load from `$XDG_CONFIG_HOME/spooky/spooky.hcl`
   - Handle environment variable overrides
   - Apply CLI flag overrides

2. **XDG directory handling** (`internal/config/xdg.go`)
   - Implement XDG Base Directory compliance
   - Handle fallback paths
   - Create directories if needed

### **Phase 3: CLI Integration (Week 2)**
**Goal**: Add configuration management to CLI

#### Tasks:
1. **Configuration commands** (`internal/cli/config.go`)
   - `spooky config show` - Show current configuration
   - `spooky config validate` - Validate configuration file
   - `spooky config init` - Create default configuration

2. **Configuration flags** (`internal/cli/commands.go`)
   - Add global configuration flags
   - Integrate with existing commands

### **Phase 4: Integration (Week 2)**
**Goal**: Integrate global config with existing systems

#### Tasks:
1. **Facts system integration**
   - Use global storage settings
   - Apply facts collection settings
   - Use performance settings

2. **SSH client integration**
   - Apply SSH timeout settings
   - Use keepalive configuration
   - Apply security settings

3. **Logging integration**
   - Use global logging configuration
   - Apply format and level settings
   - Handle file output

### **Phase 5: Testing and Documentation (Week 3)**
**Goal**: Comprehensive testing and documentation

#### Tasks:
1. **Unit tests**
   - Configuration loading tests
   - Precedence hierarchy tests
   - Validation tests

2. **Integration tests**
   - End-to-end configuration flow
   - Environment variable handling
   - CLI flag integration

3. **Documentation**
   - Configuration file format
   - Environment variables
   - CLI usage examples

## Benefits

### **Immediate Benefits**
- **Centralized configuration** for all spooky projects
- **Consistent behavior** across different environments
- **User customization** without project modification
- **Environment-specific settings** via environment variables

### **Long-term Benefits**
- **Operational consistency** across teams
- **Security policy enforcement** via global settings
- **Performance tuning** for different environments
- **Debugging capabilities** via development settings

## Success Metrics

### **Functionality Metrics**
- [ ] Global configuration file loads correctly
- [ ] Configuration precedence works as expected
- [ ] Environment variables override defaults
- [ ] CLI flags override all other sources
- [ ] Configuration validation catches errors
- [ ] Default values are sensible and secure

### **User Experience Metrics**
- [ ] Configuration file format is intuitive
- [ ] Error messages are clear and actionable
- [ ] CLI commands provide helpful feedback
- [ ] Documentation is comprehensive
- [ ] Migration from existing settings is smooth

## Risk Assessment

### **Technical Risks**
- **Configuration conflicts** - Mitigation: Clear precedence rules
- **Schema evolution** - Mitigation: Backward compatibility
- **Performance impact** - Mitigation: Lazy loading

### **User Experience Risks**
- **Configuration complexity** - Mitigation: Sensible defaults
- **Migration confusion** - Mitigation: Clear migration guide
- **Debugging difficulty** - Mitigation: Comprehensive logging

### **Implementation Risks**
- **Scope creep** - Mitigation: Focus on core configuration first
- **Testing complexity** - Mitigation: Comprehensive test coverage
- **Integration issues** - Mitigation: Incremental integration 

### **Configuration System Integration**
- **Global Configuration**: Facts collection settings from `$XDG_CONFIG_HOME/spooky/spooky.hcl` (see [Configuration System](../configuration-system.md))
- **Facts Settings**: Collection timeouts, caching TTL, parallel collection, retry attempts
- **Storage Configuration**: BadgerDB/JSON storage settings, compression, encryption
- **Performance Settings**: Concurrent fact collectors, memory limits, cache sizes
- **Configuration Precedence**: CLI flags → Project config → Global config → Environment variables → Defaults

### **Actions System Integration**
- **Action Execution Settings**: Timeouts, retries, parallel execution limits (see [Actions System](../actions-system.md))
- **SSH Configuration**: SSH settings for remote action execution
- **Logging Configuration**: Action execution logging and output settings
- **Performance Settings**: Action execution performance tuning
- **Security Settings**: Action execution security policies

### **Project System Integration**
- **Project Configuration**: Project-specific configuration overrides (see [Project System](../project-system.md))
- **Project Isolation**: Project configuration isolation from global settings
- **Project Defaults**: Default configuration values for new projects
- **Project Validation**: Configuration validation for project settings
- **Project Context**: Configuration resolution within project context

### **Variables System Integration**
- **Variable Resolution**: Configuration settings for variable resolution (see [Variables System](../variables-system.md))
- **Variable Storage**: Configuration for variable storage and caching
- **Variable Security**: Configuration for sensitive variable handling
- **Variable Performance**: Configuration for variable resolution performance
- **Variable Validation**: Configuration for variable validation settings 