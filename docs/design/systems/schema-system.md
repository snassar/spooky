# Schema Framework Plan

## Overview

This document outlines the current state of spooky's configuration schema system and provides a roadmap for enhancing it with formal schema documentation, validation tools, and IDE integration. The schema system needs to cover not just configuration files, but also the facts system, project structure, and dynamic fact sources.

## Current Schema Status

### ✅ **What's Already Implemented**

#### 1. **Comprehensive Go Structs with HCL Tags** (`internal/config/types.go`)

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
- `Action` - Individual action configuration with execution options
- `TemplateConfig` - Template-specific configuration for file operations

#### 2. **Facts System Schema** (`internal/facts/types.go`)

**Core Facts Schema:**
- `Fact` - Individual fact with key, value, metadata, and TTL
- `FactCollection` - Collection of facts with machine ID and timestamp
- `SystemFacts` - System information (OS, hardware, network)
- `CustomFacts` - User-defined custom facts

**Fact Types Schema:**
- `OSInfo` - Operating system information
- `HardwareInfo` - Hardware details (CPU, memory, storage)
- `NetworkInfo` - Network interfaces and configuration
- `DNSInfo` - DNS configuration and resolution

**Fact Storage Schema:**
- `MachineFacts` - Machine-specific fact storage
- `FactQuery` - Query interface for fact retrieval
- `StorageOptions` - Storage backend configuration

#### 3. **Validation System** (`internal/config/validator.go`)

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

#### 4. **HCL Integration**

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
- Project metadata and versioning
- Project-specific configuration overrides
- Project dependencies and imports
- Project isolation rules

#### 2. **Dynamic Facts Schema**
- Dynamic fact source definitions
- Fact collision resolution rules
- Fact TTL configuration
- Fact change detection rules

#### 3. **Global Configuration Schema**
- XDG Base Directory configuration
- Global facts database configuration
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

## Inventory Configuration Schema (`inventory.hcl`)
### Root Block: `inventory`
### Machine Block
### Validation Rules

## Actions Configuration Schema (`actions.hcl`)
### Root Block: `actions`
### Action Block
### Template Block
### Validation Rules

## Facts System Schema
### Fact Structure
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
- `inventory-schema.json` - Inventory configuration schema
- `actions-schema.json` - Actions configuration schema
- `facts-schema.json` - Facts system schema
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

**Facts Schema:**
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
                "kernel": {"type": "string"}
              }
            },
            "hardware": {
              "type": "object",
              "properties": {
                "cpu": {
                  "type": "object",
                  "properties": {
                    "cores": {"type": "integer"},
                    "model": {"type": "string"},
                    "frequency": {"type": "number"}
                  }
                },
                "memory": {
                  "type": "object",
                  "properties": {
                    "total": {"type": "integer"},
                    "available": {"type": "integer"},
                    "used": {"type": "integer"}
                  }
                }
              }
            }
          }
        },
        "custom": {"type": "object"}
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

### 3. **Enhanced CLI Schema Commands**

#### **Priority: Medium**
**Goal:** Provide schema validation and generation tools for all spooky components

**New Commands:**
```bash
# Validate against schema
spooky validate --schema project-schema.json project.hcl
spooky validate --schema facts-schema.json facts.json

# Generate schema from existing config
spooky schema generate project.hcl
spooky schema generate facts.json

# Validate schema compatibility
spooky schema validate project-schema.json
spooky schema validate facts-schema.json

# Export current schema
spooky schema export --format json
spooky schema export --format markdown

# Facts schema commands
spooky facts schema export --format json
spooky facts schema validate facts.json
```

**Features:**
- Schema validation against JSON Schema
- Schema generation from Go structs
- Schema compatibility checking
- Multiple export formats
- Facts-specific schema validation

### 4. **IDE Integration**

#### **Priority: Medium**
**Goal:** Provide seamless development experience for all spooky schemas

**VSCode Extension:**
- HCL language server with spooky schema support
- Autocomplete for all configuration fields
- Real-time validation and error highlighting
- Snippets for common configurations
- Facts schema validation and autocomplete

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
   - Facts system schema reference
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
   - Generate facts system schemas

2. **Enhance validation system**
   - Add JSON Schema validation support
   - Implement schema validation CLI commands
   - Add schema export functionality
   - Facts schema validation

### Phase 3: Tooling Integration (Week 5-6)
1. **CLI schema commands**
   - Implement `spooky schema` subcommands
   - Add schema validation to existing `validate` command
   - Create schema generation utilities
   - Facts-specific schema commands

2. **Testing and validation**
   - Test schema validation with existing projects
   - Validate generated schemas against test cases
   - Performance testing for large configurations
   - Facts schema testing

### Phase 4: IDE Support (Week 7-8)
1. **VSCode extension development**
   - Basic HCL language server integration
   - Schema-aware autocomplete
   - Error highlighting and validation
   - Facts schema support

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

## Success Metrics

### Documentation Quality
- [ ] Complete schema reference documentation
- [ ] All field types and constraints documented
- [ ] Validation error messages documented
- [ ] Best practices and examples provided
- [ ] Facts system documentation complete

### Tooling Integration
- [ ] JSON Schema files generated and validated
- [ ] CLI schema commands implemented
- [ ] Schema validation working with existing projects
- [ ] Performance impact <100ms for typical projects
- [ ] Facts schema validation working

### Developer Experience
- [ ] IDE autocomplete working with schemas
- [ ] Real-time validation in supported editors
- [ ] Error highlighting and suggestions
- [ ] Configuration snippets available
- [ ] Facts schema IDE support

### Community Adoption
- [ ] Schema documentation referenced in tutorials
- [ ] JSON Schema files available for external tools
- [ ] IDE extensions published and maintained
- [ ] Community feedback and contributions
- [ ] Facts system widely adopted

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

### Adoption Risks
- **Tool Support:** Risk of limited IDE support for custom schemas
  - **Mitigation:** Focus on VSCode and popular editors first
- **Learning Curve:** Risk of schema complexity overwhelming users
  - **Mitigation:** Comprehensive documentation and examples
- **Maintenance Burden:** Risk of schema maintenance becoming too costly
  - **Mitigation:** Automated generation and validation
- **Facts Integration:** Risk of facts schema not integrating well
  - **Mitigation:** Thorough testing with existing facts data

## Conclusion

The current spooky schema system provides a solid foundation with comprehensive Go structs, validation, and HCL integration. However, it needs to be expanded to cover the facts system, project system, and dynamic fact sources.

The proposed enhancements will significantly improve the developer experience through better documentation, tooling integration, and IDE support across all spooky components.

The phased approach allows for incremental delivery of value while managing complexity and risk. Starting with documentation provides immediate benefits, while JSON Schema and IDE integration deliver long-term developer experience improvements.

**Next Steps:**
1. Begin Phase 1: Create comprehensive schema documentation including facts system
2. Validate approach with existing configuration files and facts data
3. Gather feedback from users and contributors
4. Proceed with JSON Schema implementation based on feedback
5. Ensure facts system schema integration is prioritized 