# Schema System: Comprehensive Implementation Plan

## Overview

This document is the authoritative source for all schema system implementation details in spooky. It covers schema validation, embedded composition, runtime validation, and integration with all other spooky systems.

**Schema Integration**: This schema system provides the foundation for all other Spooky systems through embedded validation and schema composition.

**Architecture Integration**: Schema system integrates with the overall spooky architecture as described in [Spooky Design](../spooky-design.md), providing validation and type safety for all system components.

## System Integration

This schema system provides the foundation for all other Spooky systems through embedded validation and schema composition:

### **Facts System Integration**
- **Facts Schema**: Comprehensive schema validation for facts storage and structure (see [Facts System](../facts-system.md))
- **Storage Validation**: BadgerDB, JSON, and HCL storage format validation
- **Fact Collection**: Schema validation for fact collection and processing
- **Schema Composition**: Runtime schema composition for facts validation

### **Project System Integration**
- **Project Schema**: Project configuration and structure validation (see [Project System](../project-system.md))
- **Directory Validation**: Project directory structure schema enforcement
- **Metadata Schema**: Project metadata and configuration validation
- **Schema Evolution**: Project schemas evolve with system changes

### **Variables System Integration**
- **Variable Schema**: Variable definition and type validation (see [Variables System](../variables-system.md))
- **File Merging**: Schema validation for variable file merging and conflict resolution
- **Template Integration**: Schema validation for variable usage in templates
- **Export Schema**: Variable export format and metadata validation

### **Actions System Integration**
- **Action Schema**: Action definition and dependency validation (see [Actions System](../actions-system.md))
- **Dependency Schema**: Action dependency graph validation and circular reference detection
- **Run Schema**: Action run parameters and validation
- **File Merging**: Schema validation for action file merging and conflict resolution

### **CLI System Integration**
- **Command Schema**: CLI command structure and parameter validation (see [CLI System](../cli-system.md))
- **Validation Schema**: CLI validation command schema enforcement
- **Export Schema**: CLI export command format validation
- **Configuration Schema**: CLI configuration integration validation

### **Configuration System Integration**
- **Configuration Schema**: Schema validation for global configuration structure (see [Configuration System](../configuration-system.md))
- **Configuration Validation**: Configuration file validation against embedded schemas
- **Configuration Evolution**: Configuration schema versioning and migration
- **Configuration Composition**: Runtime schema composition for configuration validation
- **Configuration Integration**: Configuration integration with all system schemas

### **Machines System Integration**
- **Machine Schema**: Machine inventory validation against embedded schemas (see [Machines System](../machines-system.md))
- **Machine Validation**: Machine configuration validation using schema system
- **Machine Evolution**: Machine schema versioning and migration support
- **Machine Composition**: Runtime schema composition for machine validation
- **Machine Integration**: Machine system integration with all system schemas

### **Template System Integration**
- **Template Schema**: Template validation against embedded schemas (see [Template System](../template-system.md))
- **Template Validation**: Template configuration validation using schema system
- **Template Evolution**: Template schema versioning and migration support
- **Template Composition**: Runtime schema composition for template validation
- **Template Integration**: Template system integration with all system schemas

## Current Schema Status

### ✅ **What's Already Implemented**

#### 1. **Embedded Schema Composition System** (`internal/schemas/embedded.go`)

**Runtime Schema Composition:**
- Uses Go's `embed` directive to embed source schema files
- Composes schemas at runtime using template processing
- No separate build tools or build steps required
- Self-contained binary distribution

**Composed Schema Types:**
- `facts-hcl-composed.hcl` - HCL format facts schema
- `facts-json-composed.hcl` - JSON format facts schema
- `facts-badger-composed.hcl` - BadgerDB format facts schema

**API Functions:**
- `ComposeSchemas()` - Compose all schema combinations
- `GetComposedSchema(name)` - Get a specific composed schema
- `ListComposedSchemas()` - List available composed schemas
- `GetSchema(schemaType)` - Get embedded schema by type

## Build Process and Schema Composition

### **Embedded Runtime Composition**
- **Location**: `internal/schemas/embedded.go`
- **Process**: Runtime composition using Go's `embed` directive
- **Output**: Composed schemas generated on-demand
- **Status**: Current implementation

### Current Build Process

#### **Schema File Structure**
```
internal/schemas/schemas/
├── facts-structure.hcl          # Base fact structure (all gopsutil data)
├── storage/
│   ├── hcl.hcl                  # HCL-specific validation rules
│   ├── json.hcl                 # JSON-specific validation rules
│   └── badger.hcl               # BadgerDB-specific validation rules
├── facts-hcl.hcl         # Composed schema (includes base + format)
├── facts-json.hcl        # Composed schema (includes base + format)
└── facts-badger.hcl      # Composed schema (includes base + format)
```

#### **Build Integration**
```bash
# Build spooky with embedded schema composition
make build

# No additional steps required - schemas are composed at runtime
```

#### **Runtime Composition Process**
1. **Source Embedding**: Source schema files embedded using `go:embed`
2. **Template Processing**: Schemas composed using `text/template`
3. **On-Demand Generation**: Composed schemas generated when requested
4. **Caching**: Composed schemas cached for performance

### Benefits of Embedded Approach

- **DRY Principle**: Single source of truth for fact structure
- **Consistency**: All storage formats use identical fact definitions
- **Maintainability**: Changes to fact structure update all formats automatically
- **Extensibility**: Easy to add new storage formats
- **Performance**: No build-time overhead, runtime composition only when needed
- **Simplicity**: No separate build tools or build steps required
- **Reliability**: Self-contained binary with embedded schemas

#### 2. **Facts System Schema** (`internal/schemas/schemas/facts-structure.hcl`)

Validates facts system data structures and storage formats (see [Facts System](../facts-system.md) for implementation details):

**Comprehensive Fact Structure:**
- Complete gopsutil data coverage (OS, hardware, network, processes)
- Single source of truth for all fact definitions
- Extensible structure for custom facts
- TTL and metadata support

**Fact Categories:**
- **System Facts**: OS information, kernel version, architecture
- **Hardware Facts**: CPU, memory, storage, disk I/O
- **Network Facts**: Interfaces, IP addresses, connections
- **Process Facts**: Process counts, top processes by resource usage
- **Enhanced Facts**: Virtualization, package managers, services

**Storage Format Schemas:**
- `storage/hcl.hcl` - HCL-specific validation rules
- `storage/json.hcl` - JSON-specific validation rules
- `storage/badger.hcl` - BadgerDB-specific validation rules

#### 3. **Comprehensive Go Structs with HCL Tags** (`internal/config/types.go`)

**Project Configuration Schema:**
- `ProjectConfig` - Complete project configuration with all required and optional fields
- `StorageConfig`, `LoggingConfig`, `SSHConfig` - Nested configuration blocks
- Proper HCL tags for field mapping and validation

**Inventory Configuration Schema:**
- `InventoryConfig` - Machine inventory wrapper
- `Machine` - Individual machine configuration with authentication options
- Support for tags, ports, and connection parameters

**Actions Configuration Schema:**
- `ActionsConfig` - Action definitions wrapper
- `Action` - Individual action configuration with run options
- `TemplateConfig` - Template-specific configuration for file operations

#### 4. **Validation System** (`internal/config/validator.go`)

**Core Validation Features:**
- Uses `go-playground/validator` for struct validation
- Custom validation functions for SSH keys and script files
- Cross-field validation (e.g., machine authentication requirements)
- Comprehensive error formatting and user-friendly messages

**Validation Rules Implemented:**
- Required field validation
- Range validation (ports, timeouts, retry attempts)
- Enum validation (log levels, storage types, action types)
- File existence and permissions validation
- Mutual exclusivity validation (password vs key_file, command vs script)

#### 5. **HCL Integration**

**Wrapper Block System:**
- `ProjectConfigWrapper` - Explicit project block declaration
- `InventoryWrapper` - Explicit inventory block declaration
- `ActionsWrapper` - Explicit actions block declaration
- Proper HCL tags for field mapping, blocks, and labels

**File Structure Support:**
- Support for optional fields and blocks
- Label-based entity naming (machine names, action names)
- Nested block configuration (storage, logging, ssh, template)

### 🔄 **What Needs Schema Definition**

#### 1. **Project System Schema**

Validates project configuration and directory structure (see [Project System](../project-system.md) for implementation details):

**Project Metadata:**
- Project name, description, version
- Environment and cost center
- SLA tier and maintenance window
- Contact information

**Project Configuration:**
- Storage, logging, SSH, and template configurations
- Default timeouts and parallel running
- Inventory file and actions file paths

**Project Dependencies:**
- External dependencies and their versions
- Import paths for external libraries

**Project Isolation:**
- Namespace isolation for different projects
- Resource limits and quotas

#### 2. **Dynamic Facts Schema**
- Dynamic fact source definitions
- Fact collision resolution rules
- Fact TTL configuration
- Fact change detection rules

#### 3. **Variables System Schema** (`internal/schemas/schemas/variables-structure.hcl`)

Validates variable definitions and file merging (see [Variables System](../variables-system.md) for implementation details):

**Variable Definition Schema:**
- Variable name and type validation
- Default value and description fields
- Variable scope and visibility rules
- Variable dependency resolution

**File Merging Schema:**
- Multiple file merging rules
- Conflict detection and resolution
- Variable precedence hierarchy
- File loading order validation

**Template Integration Schema:**
- Variable usage in templates
- Variable interpolation rules
- Variable context validation
- Template function integration

#### 4. **Global Configuration Schema**
- XDG Base Directory configuration
- Facts database configuration
- Global SSH defaults
- Environment variable overrides

## Recommended Schema Enhancements

### 1. **Comprehensive Schema Documentation**

#### **Priority: High**
**Goal:** Create comprehensive, user-friendly schema documentation for all spooky components

**Deliverables:**
- Complete schema reference documentation in markdown
- Field-by-field explanations with examples
- Validation rules and error messages
- Best practices and common patterns

**Content Structure:**
```markdown
# Spooky Configuration Schema Reference

## Project Configuration Schema (`project.hcl`)
### Root Block: `project`
### File References
### Project Settings
### Storage Block
### Logging Block
### SSH Block
### Tags Block

## Machines Configuration Schema (`machines.hcl`)
### Root Block: `machines`
### Machine Block
### Validation Rules

## Actions Configuration Schema (`actions.hcl`)
### Root Block: `actions`
### Action Block
### Template Block
### Validation Rules

## Facts System Schema
### Fact Structure (facts-structure.hcl)
### Storage Formats (storage/*.hcl)
### Composed Schemas (embedded composition)
### Dynamic Facts
### Fact Storage
### Fact Validation

## Project System Schema
### Project Metadata
### Project Configuration
### Project Dependencies
### Project Isolation
```

**Benefits:**
- Improved user experience and onboarding
- Reduced configuration errors
- Better developer documentation
- Foundation for tooling integration

### 2. **JSON Schema Files**

#### **Priority: High**
**Goal:** Enable IDE integration and external tooling support for all spooky schemas

**Deliverables:**
- `project-schema.json` - Project configuration schema
- `machines-schema.json` - Machines configuration schema
- `actions-schema.json` - Actions configuration schema
- `facts-schema.json` - Facts system schema (composed)
- `global-config-schema.json` - Global configuration schema
- Schema validation utilities

**Schema Structure Examples:**

**Project Schema:**
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Spooky Project Configuration",
  "type": "object",
  "properties": {
    "project": {
      "type": "object",
      "properties": {
        "name": {"type": "string", "pattern": "^[a-zA-Z0-9_-]+$"},
        "description": {"type": "string"},
        "version": {"type": "string", "pattern": "^\\d+\\.\\d+\\.\\d+$"},
        "environment": {"type": "string"},
        "inventory_file": {"type": "string"},
        "actions_file": {"type": "string"},
        "default_timeout": {"type": "integer", "minimum": 1, "maximum": 3600},
        "default_parallel": {"type": "boolean"},
        "metadata": {
          "type": "object",
          "properties": {
            "team": {"type": "string"},
            "environment": {"type": "string"},
            "cost_center": {"type": "string"},
            "sla_tier": {"type": "string"},
            "maintenance_window": {"type": "string"},
            "contact": {"type": "string"}
          }
        },
        "config": {
          "type": "object",
          "properties": {
            "facts": {
              "type": "object",
              "properties": {
                "ttl": {"type": "string", "pattern": "^\\d+[hms]$"},
                "dynamic_fact_ttl": {"type": "string", "pattern": "^\\d+[hms]$"},
                "auto_collect": {"type": "boolean"},
                "collision_policy": {"type": "string", "enum": ["warn", "error", "merge", "highest_priority"]},
                "ttl_overrides": {"type": "object"},
                "collision_rules": {"type": "object"}
              }
            },
            "ssh": {
              "type": "object",
              "properties": {
                "timeout": {"type": "string", "pattern": "^\\d+[hms]$"},
                "keepalive_interval": {"type": "string", "pattern": "^\\d+[hms]$"},
                "retry_attempts": {"type": "integer", "minimum": 0, "maximum": 10},
                "parallel_connections": {"type": "integer", "minimum": 1, "maximum": 10}
              }
            }
          }
        },
        "dependencies": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "name": {"type": "string"},
              "version": {"type": "string"},
              "source": {"type": "string"}
            },
            "required": ["name"]
          }
        }
      },
      "required": ["name"]
    }
  },
  "required": ["project"]
}
```

**Facts Schema (Composed):**
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Spooky Facts Schema",
  "type": "object",
  "properties": {
    "machine_id": {"type": "string", "pattern": "^[a-f0-9]{32}$"},
    "collected_at": {"type": "string", "format": "date-time"},
    "ttl": {"type": "string", "pattern": "^\\d+[hms]$"},
    "facts": {
      "type": "object",
      "properties": {
        "system": {
          "type": "object",
          "properties": {
            "os": {
              "type": "object",
              "properties": {
                "name": {"type": "string"},
                "version": {"type": "string"},
                "arch": {"type": "string"},
                "kernel": {"type": "string"},
                "platform": {"type": "string"},
                "family": {"type": "string"}
              }
            },
            "hardware": {
              "type": "object",
              "properties": {
                "cpu": {
                  "type": "object",
                  "properties": {
                    "cores": {"type": "integer", "minimum": 1},
                    "model": {"type": "string"},
                    "frequency": {"type": "number"},
                    "architecture": {"type": "string"},
                    "vendor": {"type": "string"},
                    "percent": {"type": "number"}
                  }
                },
                "memory": {
                  "type": "object",
                  "properties": {
                    "total": {"type": "integer", "minimum": 1},
                    "available": {"type": "integer"},
                    "used": {"type": "integer"},
                    "free": {"type": "integer"},
                    "percent": {"type": "number"}
                  }
                },
                "disks": {
                  "type": "array",
                  "items": {
                    "type": "object",
                    "properties": {
                      "device": {"type": "string"},
                      "mount_point": {"type": "string"},
                      "total": {"type": "integer"},
                      "used": {"type": "integer"},
                      "free": {"type": "integer"},
                      "filesystem": {"type": "string"}
                    }
                  }
                }
              }
            },
            "network": {
              "type": "object",
              "properties": {
                "hostname": {"type": "string"},
                "interfaces": {
                  "type": "array",
                  "items": {
                    "type": "object",
                    "properties": {
                      "name": {"type": "string"},
                      "mac_address": {"type": "string"},
                      "ip_addresses": {"type": "array", "items": {"type": "string"}},
                      "mtu": {"type": "integer"}
                    }
                  }
                },
                "ip_addresses": {"type": "array", "items": {"type": "string"}},
                "primary_ip": {"type": "string"}
              }
            }
          }
        },
        "enhanced": {
          "type": "object",
          "properties": {
            "virtualization": {
              "type": "object",
              "properties": {
                "system": {"type": "string"},
                "role": {"type": "string"}
              }
            },
            "package_manager": {
              "type": "object",
              "properties": {
                "type": {"type": "string", "enum": ["apt", "yum", "dnf", "zypper", "pacman", "apk", "unknown"]},
                "version": {"type": "string"}
              }
            },
            "service_manager": {
              "type": "object",
              "properties": {
                "type": {"type": "string", "enum": ["systemd", "upstart", "init", "runit", "openrc", "unknown"]},
                "version": {"type": "string"}
              }
            }
          }
        }
      }
    }
  },
  "required": ["machine_id", "collected_at", "facts"]
}
```

**Benefits:**
- IDE autocomplete and validation
- External tool integration
- Schema-aware editors
- Automated documentation generation

### 3. **Internal Schema Management**

#### **Priority: Medium**
**Goal:** Provide robust internal schema validation and management for all spooky components

**Internal Schema Usage:**
- **Configuration Validation**: All HCL files validated against embedded schemas
- **Facts Validation**: Facts data validated against composed schemas
- **Project Validation**: Project structure validated against project schemas
- **Template Validation**: Templates validated against variable and facts schemas

**Schema Integration Points:**
- **`spooky project validate`**: Uses embedded schemas for all validation operations
- **`spooky facts export`**: Uses composed facts schemas for data validation
- **`spooky project validate`**: Uses project schemas for structure validation
- **`spooky variables validate`**: Uses variable schemas for configuration validation
- **`spooky templates validate`**: Uses template schemas for syntax validation

**Features:**
- Automatic schema loading from embedded sources
- Runtime schema composition for facts system
- Schema validation for all configuration files
- Facts data validation against composed schemas
- Template syntax validation against variable schemas

### 4. **IDE Integration**

#### **Priority: Medium**
**Goal:** Provide seamless development experience for all spooky schemas

**VSCode Extension:**
- HCL language server with spooky schema support
- Autocomplete for all configuration fields
- Real-time validation and error highlighting
- Snippets for common configurations
- Facts schema validation and autocomplete
- Embedded schema support

**JetBrains Plugin:**
- IntelliJ IDEA/GoLand plugin for spooky configs
- Integration with existing HCL support
- Custom validation rules
- Refactoring support
- Facts system integration

**Vim/Neovim Support:**
- Tree-sitter grammar for spooky HCL
- Language server integration
- Syntax highlighting and folding
- Facts schema support

### 5. **Schema Versioning and Migration**

#### **Priority: Low**
**Goal:** Handle schema evolution gracefully across all components

**Implementation:**
```hcl
project "my-project" {
  schema_version = "1.0"  # New field for versioning
  # ... rest of config
}

# Facts schema versioning
facts {
  schema_version = "1.0"
  # ... facts config
}
```

**Features:**
- Backward compatibility support
- Migration tools for schema updates
- Deprecation warnings
- Version-specific validation rules
- Facts schema migration support

## Implementation Roadmap

### Phase 1: Documentation Foundation (Week 1-2)
1. **Create comprehensive schema documentation**
   - Project configuration schema reference
   - Inventory configuration schema reference
   - Actions configuration schema reference
   - Facts system schema reference (including embedded composition)
   - Global configuration schema reference
   - Validation rules and error messages
   - Best practices and examples

2. **Update existing documentation**
   - Integrate schema reference into main docs
   - Add schema examples to configuration guide
   - Update CLI specification with schema commands
   - Add facts system documentation

### Phase 2: JSON Schema Implementation (Week 3-4)
1. **Generate JSON Schema files**
   - Create schema generation tool
   - Export current Go structs to JSON Schema
   - Validate schema accuracy against existing configs
   - Generate facts system schemas (including composed schemas)

2. **Enhance validation system**
   - Add JSON Schema validation support
   - Implement schema validation CLI commands
   - Add schema export functionality
   - Facts schema validation

### Phase 3: Internal Schema Integration (Week 5-6)
1. **Schema validation integration**
   - Integrate embedded schemas into existing `spooky project validate` command
   - Add facts schema validation to `spooky facts export`
   - Add project schema validation to `spooky project validate`
   - Add variable schema validation to `spooky variables validate`
   - Add template schema validation to `spooky templates validate`

2. **Testing and validation**
   - Test embedded schema validation with existing projects
   - Validate composed schemas against test cases
   - Performance testing for large configurations
   - Facts schema testing with real data

### Phase 4: IDE Support (Week 7-8)
1. **VSCode extension development**
   - Basic HCL language server integration
   - Schema-aware autocomplete
   - Error highlighting and validation
   - Facts schema support
   - Embedded schema support

2. **Documentation and examples**
   - IDE setup guides
   - Configuration examples
   - Troubleshooting documentation
   - Facts system examples

## Technical Considerations

### Schema Generation Strategy

**Option 1: Go Struct Reflection**
- Use Go reflection to generate JSON Schema from structs
- Maintains single source of truth
- Automatic updates when structs change
- More complex implementation
- Works for facts system and configuration

**Option 2: Manual Schema Definition**
- Manually define JSON Schema files
- Easier to customize and optimize
- Requires manual maintenance
- Risk of drift between Go structs and schema

**Recommendation:** Start with Option 2 for Phase 2, then implement Option 1 for Phase 3

### Validation Performance

**Current Performance:**
- Go struct validation: ~1ms per configuration file
- HCL parsing: ~5-10ms per file
- Embedded schema composition: ~1-2ms per schema
- Total validation time: <50ms for typical projects

**Expected Performance with JSON Schema:**
- JSON Schema validation: ~2-5ms per file
- Schema loading: ~1-2ms per schema
- Total validation time: <100ms for typical projects
- Facts validation: ~1-3ms per machine

### Backward Compatibility

**Schema Evolution Rules:**
1. New fields must be optional
2. Existing fields cannot be removed without deprecation period
3. Field type changes require major version bump
4. Validation rule changes must be backward compatible
5. Facts schema changes must maintain data compatibility
6. Embedded schema composition must remain backward compatible

## Success Metrics

### Documentation Quality
- [ ] Complete schema reference documentation
- [ ] All field types and constraints documented
- [ ] Validation error messages documented
- [ ] Best practices and examples provided
- [ ] Facts system documentation complete
- [ ] Embedded schema composition documented

### Tooling Integration
- [ ] JSON Schema files generated and validated
- [ ] Embedded schemas integrated into existing CLI commands
- [ ] Schema validation working with existing projects
- [ ] Performance impact <100ms for typical projects
- [ ] Facts schema validation working
- [ ] Internal schema validation working across all commands

### Developer Experience
- [ ] IDE autocomplete working with schemas
- [ ] Real-time validation in supported editors
- [ ] Error highlighting and suggestions
- [ ] Configuration snippets available
- [ ] Facts schema IDE support
- [ ] Embedded schema IDE support

### Community Adoption
- [ ] Schema documentation referenced in tutorials
- [ ] JSON Schema files available for external tools
- [ ] IDE extensions published and maintained
- [ ] Community feedback and contributions
- [ ] Facts system widely adopted
- [ ] Internal schema validation approach adopted

## Risk Assessment

### Technical Risks
- **Schema Drift:** Risk of JSON Schema becoming out of sync with Go structs
  - **Mitigation:** Automated schema generation and validation tests
- **Performance Impact:** Risk of validation becoming too slow
  - **Mitigation:** Performance testing and optimization
- **Complexity:** Risk of over-engineering the schema system
  - **Mitigation:** Start simple, iterate based on user feedback
- **Facts Schema Complexity:** Risk of facts schema being too complex
  - **Mitigation:** Modular schema design with clear separation
- **Embedded Schema Maintenance:** Risk of embedded schemas becoming difficult to maintain
  - **Mitigation:** Automated composition and validation

### Adoption Risks
- **Tool Support:** Risk of limited IDE support for custom schemas
  - **Mitigation:** Focus on VSCode and popular editors first
- **Learning Curve:** Risk of schema complexity overwhelming users
  - **Mitigation:** Comprehensive documentation and examples
- **Maintenance Burden:** Risk of schema maintenance becoming too costly
  - **Mitigation:** Automated generation and validation
- **Facts Integration:** Risk of facts schema not integrating well
  - **Mitigation:** Thorough testing with existing facts data
- **Embedded Approach:** Risk of embedded approach being too complex
  - **Mitigation:** Clear documentation and examples

## Conclusion

The current spooky schema system provides a solid foundation with comprehensive Go structs, validation, HCL integration, and embedded schema composition. The facts system is well-integrated with runtime schema composition that eliminates build complexity.

The proposed enhancements will significantly improve the developer experience through better documentation, tooling integration, and IDE support across all spooky components, including the embedded schema composition system.

The phased approach allows for incremental delivery of value while managing complexity and risk. Starting with documentation provides immediate benefits, while JSON Schema and IDE integration deliver long-term developer experience improvements.

**Next Steps:**
1. Begin Phase 1: Create comprehensive schema documentation including facts system and embedded composition
2. Validate approach with existing configuration files and facts data
3. Gather feedback from users and contributors
4. Proceed with JSON Schema implementation based on feedback
5. Ensure facts system schema integration and embedded composition are prioritized 