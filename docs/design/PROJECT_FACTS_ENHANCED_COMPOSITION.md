# Project Facts Enhanced Composition Schema

## Overview

This document describes the enhanced composition pattern implemented for project facts schemas, providing clear separation between base structure and format-specific features.

## Schema Architecture

### Base Schema: `project-facts-structure.hcl`
- **Purpose**: Common structure definitions for all storage formats
- **Content**: 
  - Core fact structure (machine_id, project_id, collected_at, ttl, facts)
  - Comprehensive fact definitions (applications, deployment, environment, monitoring, custom)
  - Common validation rules (format-agnostic)
  - Common constraints (format-agnostic)
  - Common metadata (format-agnostic)

### Format-Specific Schemas

#### 1. `project-facts-hcl.hcl` (Enhanced)
- **Purpose**: HCL-specific features and validation
- **Features**:
  - HCL interpolation support (`${var.name}`)
  - HCL expressions support (`var.count + 1`)
  - HCL heredoc support (`<<-EOT...EOT`)
  - HCL comments support (`#`, `//`, `/* */`)
  - HCL block syntax (`fact "name" { ... }`)
  - HCL attribute syntax (`attribute = value`)

#### 2. `project-facts-json.hcl` (Enhanced)
- **Purpose**: JSON-specific features and validation
- **Features**:
  - JSON data types support
  - JSON export/import format structures
  - JSON Schema validation integration
  - JSON format limitations (no comments, no expressions, etc.)
  - JSON performance constraints

#### 3. `project-facts-badger.hcl` (Enhanced)
- **Purpose**: BadgerDB-specific features and validation
- **Features**:
  - BadgerDB compression (zstd)
  - BadgerDB encryption (age)
  - BadgerDB transactions (ACID)
  - BadgerDB indexing (prefix, value)
  - BadgerDB garbage collection
  - BadgerDB key/value structure
  - BadgerDB query patterns
  - BadgerDB backup structure

## Enhanced Composition Pattern

### Structure
```hcl
# Format-specific schema
project_facts {
  # Include base structure
  include = "project-facts-structure"
  
  # Format-specific metadata
  scope = "project"
  storage_location = "<format-specific-path>"
  description = "Project-specific facts stored in <format> format"
  
  # Format-specific features
  <format>_features = {
    # Format-specific capabilities
  }
  
  # Format-specific validation
  <format>_validation = {
    # Format-specific validation rules
  }
  
  # Format-specific constraints
  <format>_constraints = {
    # Format-specific limitations and constraints
  }
  
  # Format-specific structure extensions
  <format>_extensions = {
    # Format-specific structure enhancements
  }
}
```

### Benefits

1. **Clear Separation**: Base structure vs format-specific features
2. **Consistent Pattern**: All format schemas follow same composition pattern
3. **Extensible**: Easy to add new formats
4. **Maintainable**: Changes to base structure propagate to all formats
5. **Comprehensive**: Each format schema leverages format-specific capabilities

## Implementation Details

### Base Schema Enhancements
- Added common validation rules (machine_id_format, project_id_required, etc.)
- Added format-agnostic constraints (scope, isolation, default_ttl, etc.)
- Added common metadata (schema_version, schema_type, last_updated, etc.)

### HCL Schema Enhancements
- **Features**: Interpolation, expressions, heredoc, comments, block syntax
- **Validation**: HCL syntax, interpolation, expressions, block structure
- **Constraints**: Format limitations, performance constraints
- **Extensions**: Interpolated facts, block facts, template integration

### JSON Schema Enhancements
- **Features**: Data types, export/import formats, schema validation
- **Validation**: JSON syntax, schema compliance, key/value validation
- **Constraints**: Format limitations, performance constraints, export features
- **Extensions**: Export structure, import structure, schema integration

### BadgerDB Schema Enhancements
- **Features**: Compression, encryption, transactions, indexing, garbage collection
- **Validation**: Key format, value format, transaction, encryption, compression
- **Constraints**: Storage constraints, performance constraints, encryption constraints
- **Extensions**: Key structure, value structure, query structure, backup structure

## Usage Examples

### HCL Facts Example
```hcl
# project-facts.hcl
fact "application" {
  version = "${var.app_version}"
  config_path = "${var.config_dir}/app.conf"
  port_count = var.base_port + 10
  
  config_content = <<-EOT
    server {
      port 8080;
      host "localhost";
    }
  EOT
}
```

### JSON Facts Example
```json
{
  "metadata": {
    "version": "1.0.0",
    "exported_at": "2024-01-01T00:00:00Z",
    "format": "json"
  },
  "facts": {
    "application": {
      "version": "1.0.0",
      "config_path": "/etc/app/config.conf",
      "port_count": 8080
    }
  }
}
```

### BadgerDB Facts Example
```go
// Key: "facts:a1b2c3d4e5f67890123456789012345678"
// Value: JSON-encoded fact collection with compression and encryption
{
  "machine_id": "a1b2c3d4e5f67890123456789012345678",
  "project_id": "my-project",
  "collected_at": "2024-01-01T00:00:00Z",
  "ttl": "1h",
  "facts": {
    "application": {
      "version": "1.0.0",
      "config_path": "/etc/app/config.conf"
    }
  },
  "encrypted": true,
  "compression": "zstd"
}
```

## Migration Guide

### From Old Schema Structure
1. **Base Structure**: Enhanced with common validation, constraints, and metadata
2. **Format Schemas**: Completely rewritten to follow enhanced composition pattern
3. **Validation**: Moved from format-specific to base structure where appropriate
4. **Features**: Added comprehensive format-specific capabilities

### Schema Composition Logic
The schema composition logic should:
1. Load base structure schema
2. Load format-specific schema
3. Merge base structure with format-specific extensions
4. Validate composition results
5. Apply format-specific validation rules

## Future Enhancements

1. **Schema Composition Engine**: Implement automatic schema merging
2. **Format Detection**: Automatic format detection based on schema validation
3. **Schema Validation**: Comprehensive validation of composed schemas
4. **Schema Documentation**: Auto-generated documentation from schema definitions
5. **Schema Testing**: Comprehensive test suite for schema composition

## Conclusion

The enhanced composition pattern provides a robust, maintainable, and extensible approach to schema management for project facts. Each format schema now fully leverages its specific capabilities while maintaining consistency through the common base structure. 