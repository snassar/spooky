# Facts Export Schema Compatibility Update

## Overview

This document describes the compatibility update made to `facts-export.hcl` to ensure it works seamlessly with the enhanced `global-facts-structure.hcl` and `project-facts-structure.hcl` schemas.

## Problem Statement

The original `facts-export.hcl` schema was **incompatible** with the enhanced structure schemas due to:

1. **Schema Structure Mismatch**: Used simple HCL syntax vs detailed type definitions
2. **Missing Fields**: Lacked fields present in enhanced schemas
3. **Field Type Mismatches**: Used generic types vs specific field definitions
4. **Missing Validation**: No validation rules vs comprehensive validation

## Solution Implemented

### **Complete Schema Restructure**

The `facts-export.hcl` schema has been completely restructured to match the enhanced schemas:

#### **Before (Incompatible)**
```hcl
facts_export = {
  metadata = {
    exported_at = string
    project_id = string
    export_format = string
    version = string
  }
  
  global_facts = [
    {
      machine_id = string
      collected_at = string
      ttl = string
      facts = {
        system = {
          os = {
            name = string
            version = string
            # ... generic structure
          }
        }
      }
    }
  ]
}
```

#### **After (Compatible)**
```hcl
facts_export = {
  metadata = {
    exported_at = {
      type = "string"
      required = true
      format = "date-time"
      description = "Timestamp when facts were exported (ISO 8601 format)"
    }
    project_id = {
      type = "string"
      required = true
      description = "Project identifier for the exported facts"
    }
    # ... detailed structure matching enhanced schemas
  }
  
  global_facts = [
    {
      machine_id = {
        type = "string"
        required = true
        pattern = "^[a-f0-9]{32}$"
        description = "Unique machine identifier from /etc/machine-id"
      }
      # ... complete structure matching global-facts-structure.hcl
    }
  ]
}
```

## Key Changes Made

### **1. Metadata Structure Enhancement**
- **Added**: Detailed type definitions with validation rules
- **Added**: Source files tracking
- **Added**: Fact count information
- **Added**: Export format validation (json/hcl only)

### **2. Global Facts Compatibility**
- **Updated**: Complete structure to match `global-facts-structure.hcl`
- **Added**: All missing fields (custom facts, detailed hardware info)
- **Added**: Proper validation rules and constraints
- **Added**: Enhanced system information support

### **3. Project Facts Compatibility**
- **Updated**: Complete structure to match `project-facts-structure.hcl`
- **Added**: Detailed application, deployment, environment, monitoring structures
- **Added**: Proper validation rules and constraints
- **Added**: Custom facts support

### **4. Field Type Standardization**
- **Replaced**: Generic `string = number` patterns
- **Added**: Specific field names and types
- **Added**: Proper validation patterns (regex, enums, etc.)
- **Added**: Required/optional field specifications

## Compatibility Benefits

### **✅ Full Import/Export Compatibility**
- Exported HCL/JSON files can be imported using enhanced schemas
- No data loss or structure mismatches
- Proper validation during import/export operations

### **✅ Enhanced Validation**
- Export data validated against enhanced schemas
- Proper error reporting for invalid data
- Consistent validation rules across all operations

### **✅ Future-Proof Design**
- Ready for schema composition features
- Compatible with enhanced composition pattern
- Extensible for new fact types

### **✅ Consistent Architecture**
- Same structure patterns across all schemas
- Unified validation and constraint handling
- Consistent metadata and documentation

## Usage Examples

### **Export Operation**
```bash
# Export facts using enhanced schema validation
spooky facts export --format json --output facts.json
```

### **Import Operation**
```bash
# Import facts using enhanced schema validation
spooky facts import facts.json
```

### **Generated Export File Structure**
```json
{
  "metadata": {
    "exported_at": "2024-01-01T00:00:00Z",
    "project_id": "my-project",
    "export_format": "json",
    "version": "1.0.0",
    "source_files": ["global-facts.db", "project-facts.db"],
    "fact_count": {
      "global_facts": 5,
      "project_facts": 3
    }
  },
  "global_facts": [
    {
      "machine_id": "a1b2c3d4e5f67890123456789012345678",
      "collected_at": "2024-01-01T00:00:00Z",
      "ttl": "24h",
      "facts": {
        "system": {
          "os": {
            "name": "Ubuntu",
            "version": "22.04 LTS",
            "arch": "x86_64",
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
        },
        "custom": {
          "environment": "production"
        }
      }
    }
  ],
  "project_facts": [
    {
      "machine_id": "a1b2c3d4e5f67890123456789012345678",
      "project_id": "my-project",
      "collected_at": "2024-01-01T00:00:00Z",
      "ttl": "1h",
      "facts": {
        "applications": {
          "versions": {
            "nginx": "1.18.0",
            "postgresql": "14.5"
          }
        },
        "deployment": {
          "state": "deployed",
          "info": {
            "version": "1.0.0",
            "deployed_at": "2024-01-01T00:00:00Z"
          }
        }
      }
    }
  ]
}
```

## Migration Guide

### **For Existing Export Files**
1. **Backward Compatibility**: Existing export files will continue to work
2. **Enhanced Validation**: New validation rules will apply to future exports
3. **Gradual Migration**: No immediate action required for existing files

### **For New Exports**
1. **Automatic Enhancement**: New exports automatically use enhanced schema
2. **Better Validation**: Improved error detection and reporting
3. **Rich Metadata**: Enhanced metadata for better tracking

### **For Import Operations**
1. **Enhanced Validation**: All imports validated against enhanced schemas
2. **Better Error Messages**: More specific error reporting
3. **Data Integrity**: Improved data integrity checks

## Testing

### **Build Verification**
```bash
# Verify schema compilation
go build ./internal/schemas
```

### **Schema Validation**
```bash
# Test export/import operations
spooky facts export --format json --output test.json
spooky facts import test.json
```

## Conclusion

The `facts-export.hcl` schema is now **fully compatible** with the enhanced `global-facts-structure.hcl` and `project-facts-structure.hcl` schemas. This ensures:

- **Seamless Import/Export**: Exported files can be imported without issues
- **Enhanced Validation**: Better data validation and error reporting
- **Consistent Architecture**: Unified schema patterns across the system
- **Future-Proof Design**: Ready for advanced schema features

The compatibility update maintains backward compatibility while providing enhanced functionality for new operations. 