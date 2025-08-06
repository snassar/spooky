# Global Facts Enhanced Composition Schema

## Overview

This document describes the enhanced composition pattern implemented for global facts schemas, providing clear separation between base structure and format-specific features, with special considerations for global scope and cross-project access.

## Schema Architecture

### Base Schema: `global-facts-structure.hcl`
- **Purpose**: Common global fact structure definitions for all storage formats
- **Content**: 
  - Core fact structure (machine_id, collected_at, ttl, facts)
  - Comprehensive system fact definitions (OS, hardware, network, processes, enhanced)
  - Common validation rules (format-agnostic)
  - Common constraints (format-agnostic, global-specific)
  - Common metadata (format-agnostic)

### Format-Specific Schemas

#### 1. `global-facts-hcl.hcl` (Enhanced)
- **Purpose**: HCL-specific features and validation for global facts
- **Features**:
  - HCL interpolation support (`${var.name}`)
  - HCL expressions support (`var.count + 1`)
  - HCL heredoc support (`<<-EOT...EOT`)
  - HCL comments support (`#`, `//`, `/* */`)
  - HCL block syntax (`fact "name" { ... }`)
  - HCL attribute syntax (`attribute = value`)
  - Global scope variables and cross-project references

#### 2. `global-facts-json.hcl` (Enhanced)
- **Purpose**: JSON-specific features and validation for global facts
- **Features**:
  - JSON data types support
  - JSON export/import format structures
  - JSON Schema validation integration
  - JSON format limitations (no comments, no expressions, etc.)
  - JSON performance constraints (higher limits for global facts)
  - Global scope information in export/import

#### 3. `global-facts-badger.hcl` (Enhanced)
- **Purpose**: BadgerDB-specific features and validation for global facts
- **Features**:
  - BadgerDB compression (zstd)
  - BadgerDB encryption (age)
  - BadgerDB transactions (ACID)
  - BadgerDB indexing (prefix, value, global scope, cross-project)
  - BadgerDB garbage collection (longer intervals for global facts)
  - BadgerDB backup and recovery (global scope specific)
  - BadgerDB key/value structure (global scope specific)
  - BadgerDB query patterns (cross-project queries)

## Enhanced Composition Pattern

### Structure
```hcl
# Format-specific schema for global facts
global_facts {
  # Include base structure
  include = "global-facts-structure"
  
  # Format-specific metadata
  scope = "global"
  storage_location = "<format-specific-path>"
  description = "Global facts stored in <format> format"
  
  # Format-specific features
  <format>_features = {
    # Format-specific capabilities
  }
  
  # Format-specific validation
  <format>_validation = {
    # Format-specific validation rules
    # Global scope validation
  }
  
  # Format-specific constraints
  <format>_constraints = {
    # Format-specific limitations and constraints
    # Global facts specific constraints
  }
  
  # Format-specific structure extensions
  <format>_extensions = {
    # Format-specific structure enhancements
    # Global facts specific extensions
  }
}
```

### Benefits

1. **Clear Separation**: Base structure vs format-specific features
2. **Consistent Pattern**: All format schemas follow same composition pattern
3. **Extensible**: Easy to add new formats
4. **Maintainable**: Changes to base structure propagate to all formats
5. **Comprehensive**: Each format schema leverages format-specific capabilities
6. **Global Scope Aware**: Special handling for global facts and cross-project access

## Implementation Details

### Base Schema Enhancements
- Added common validation rules (machine_id_format, timestamp_format, etc.)
- Added format-agnostic constraints (scope, shared_access, default_ttl, etc.)
- Added common metadata (schema_version, schema_type, last_updated, etc.)
- Added global-specific constraints (persistence, isolation, etc.)

### HCL Schema Enhancements
- **Features**: Interpolation, expressions, heredoc, comments, block syntax
- **Validation**: HCL syntax, interpolation, expressions, block structure, global scope
- **Constraints**: Format limitations, performance constraints (higher limits for global facts)
- **Extensions**: Interpolated facts, block facts, template integration, global extensions

### JSON Schema Enhancements
- **Features**: Data types, export/import formats, schema validation
- **Validation**: JSON syntax, schema compliance, key/value validation, global scope
- **Constraints**: Format limitations, performance constraints (higher limits for global facts), export features
- **Extensions**: Export structure, import structure, schema integration, global extensions

### BadgerDB Schema Enhancements
- **Features**: Compression, encryption, transactions, indexing, garbage collection, backup/recovery
- **Validation**: Key format, value format, transaction, encryption, compression, global scope
- **Constraints**: Storage constraints, performance constraints (higher limits for global facts), encryption constraints
- **Extensions**: Key structure, value structure, query structure, backup structure, global extensions

## Global Facts Specific Features

### Global Scope Characteristics
- **Shared Access**: Global facts are shared across all projects
- **Persistent Storage**: Global facts persist across spooky sessions
- **Longer TTL**: Default TTL of 24h (vs 1h for project facts)
- **Cross-Project Access**: Facts accessible from any project
- **System-Wide Scope**: Facts apply to the entire system

### Enhanced Constraints for Global Facts
- **Larger File Sizes**: 2MB for HCL/JSON (vs 1MB for project facts)
- **More Facts per File**: 5000 for HCL/JSON (vs 1000 for project facts)
- **Deeper Nesting**: 15 levels (vs 10 for project facts)
- **Higher Performance Limits**: More concurrent operations, larger transactions
- **Enhanced Encryption**: More recipients, global scope keys

### Global Facts Extensions
- **Cross-Project References**: Facts that reference other projects
- **Global Variables**: Variables accessible to all projects
- **System-Wide Configuration**: Configuration that applies globally
- **Access Control**: Control who can access global facts
- **Backup and Recovery**: Global scope specific backup policies

## Usage Examples

### HCL Global Facts Example
```hcl
# global-facts.hcl
fact "system" {
  hostname = "${var.hostname}"
  os = "${var.os_name}"
  kernel = "${var.kernel_version}"
  
  system_info = <<-EOT
    Hostname: ${hostname}
    OS: ${os.name}
    Kernel: ${os.kernel}
    Architecture: ${os.arch}
  EOT
  
  # Global scope variables
  global_variables = {
    system_timezone = "${var.timezone}"
    system_locale = "${var.locale}"
  }
  
  # Cross-project references
  shared_projects = ["project-a", "project-b", "project-c"]
}
```

### JSON Global Facts Example
```json
{
  "metadata": {
    "version": "1.0.0",
    "exported_at": "2024-01-01T00:00:00Z",
    "format": "json",
    "global_scope": "true",
    "shared_projects": ["project-a", "project-b", "project-c"]
  },
  "facts": {
    "system": {
      "hostname": "server01.example.com",
      "os": {
        "name": "Ubuntu",
        "version": "22.04 LTS",
        "kernel": "5.15.0-91-generic"
      },
      "hardware": {
        "cpu": {
          "cores": 8,
          "model": "Intel(R) Core(TM) i7-9700K"
        },
        "memory": {
          "total": 17179869184
        }
      }
    }
  }
}
```

### BadgerDB Global Facts Example
```go
// Key: "global_facts:a1b2c3d4e5f67890123456789012345678"
// Value: JSON-encoded global fact collection with compression and encryption
{
  "machine_id": "a1b2c3d4e5f67890123456789012345678",
  "collected_at": "2024-01-01T00:00:00Z",
  "ttl": "24h",
  "facts": {
    "system": {
      "hostname": "server01.example.com",
      "os": {
        "name": "Ubuntu",
        "version": "22.04 LTS"
      }
    }
  },
  "encrypted": true,
  "compression": "zstd",
  "global_scope": "true",
  "shared_projects": ["project-a", "project-b", "project-c"]
}
```

## Migration Guide

### From Old Schema Structure
1. **Base Structure**: Enhanced with common validation, constraints, and metadata
2. **Format Schemas**: Completely rewritten to follow enhanced composition pattern
3. **Validation**: Moved from format-specific to base structure where appropriate
4. **Features**: Added comprehensive format-specific capabilities
5. **Global Scope**: Added global-specific features and constraints

### Schema Composition Logic
The schema composition logic should:
1. Load base structure schema
2. Load format-specific schema
3. Merge base structure with format-specific extensions
4. Validate composition results
5. Apply format-specific validation rules
6. Apply global scope validation rules

## Comparison with Project Facts

| Aspect | Project Facts | Global Facts |
|--------|---------------|--------------|
| **Scope** | Project-specific | System-wide |
| **Access** | Isolated per project | Shared across projects |
| **TTL** | 1h default | 24h default |
| **File Size** | 1MB limit | 2MB limit |
| **Facts per File** | 1000 | 5000 |
| **Nesting Depth** | 10 levels | 15 levels |
| **Concurrent Operations** | Lower limits | Higher limits |
| **Encryption Recipients** | 10 max | 20 max |
| **Backup** | Project scope | Global scope |
| **Cross-Project Access** | No | Yes |

## Future Enhancements

1. **Schema Composition Engine**: Implement automatic schema merging
2. **Format Detection**: Automatic format detection based on schema validation
3. **Schema Validation**: Comprehensive validation of composed schemas
4. **Schema Documentation**: Auto-generated documentation from schema definitions
5. **Schema Testing**: Comprehensive test suite for schema composition
6. **Global Scope Management**: Enhanced global scope management tools
7. **Cross-Project Integration**: Advanced cross-project fact sharing

## Conclusion

The enhanced composition pattern provides a robust, maintainable, and extensible approach to schema management for global facts. Each format schema now fully leverages its specific capabilities while maintaining consistency through the common base structure. The global facts schemas include special considerations for global scope, cross-project access, and system-wide configuration, making them suitable for managing facts that are shared across all projects in the spooky system. 