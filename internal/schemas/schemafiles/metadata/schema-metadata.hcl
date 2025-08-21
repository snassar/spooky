# Schema Metadata Schema
# Meta-schema that defines the structure and validation rules for schema metadata blocks
# This schema validates the metadata blocks in all other schema files

# Schema metadata
metadata {
  schema_version = "0.20250809.0"
  schema_type = "schema-metadata"
  schema_name = "Schema Metadata Schema"
  last_updated = "2024-01-01"
  compatibility = ["0.20250809.0"]
  description = "Meta-schema that defines the structure and validation rules for schema metadata blocks used across all schema files"
  
  # ScalVer format: 0.YYYYMMDD.N
  # - 0: Development phase
  # - 20250809: Date (9 August 2025)
  # - 0: Patch version
  scalver_format = "0.20250809.0"
}

# Schema metadata structure
schema_metadata {
  # Required schema identification
  schema_version = {
    type = "string"
    required = true
    pattern = "^0\\.\\d{8}\\.\\d+$"
    description = "Schema version in ScalVer format (0.YYYYMMDD.N)"
    validation = {
      rule = "regex"
      pattern = "^0\\.\\d{8}\\.\\d+$"
      message = "Schema version must be in ScalVer format: 0.YYYYMMDD.N"
    }
  }
  
  schema_type = {
    type = "string"
    required = true
    pattern = "^[a-z][a-z0-9-]*$"
    description = "Schema type identifier (lowercase with hyphens)"
    validation = {
      rule = "regex"
      pattern = "^[a-z][a-z0-9-]*$"
      message = "Schema type must be lowercase with hyphens only"
    }
  }
  
  schema_name = {
    type = "string"
    required = true
    min_length = 1
    max_length = 128
    description = "Human-readable schema name"
    validation = {
      rule = "length"
      min = 1
      max = 128
      message = "Schema name must be between 1 and 128 characters"
    }
  }
  
  # Required schema information
  last_updated = {
    type = "string"
    required = true
    format = "date"
    pattern = "^\\d{4}-\\d{2}-\\d{2}$"
    description = "Last schema update date (YYYY-MM-DD format)"
    validation = {
      rule = "regex"
      pattern = "^\\d{4}-\\d{2}-\\d{2}$"
      message = "Last updated must be in YYYY-MM-DD format"
    }
  }
  
  description = {
    type = "string"
    required = true
    min_length = 10
    max_length = 500
    description = "Schema description explaining purpose and scope"
    validation = {
      rule = "length"
      min = 10
      max = 500
      message = "Description must be between 10 and 500 characters"
    }
  }
  
  # Optional schema information
  compatibility = {
    type = "array"
    required = false
    description = "List of compatible schema versions"
    items = {
      type = "string"
      pattern = "^0\\.\\d{8}\\.\\d+$"
    }
    validation = {
      rule = "array_items"
      pattern = "^0\\.\\d{8}\\.\\d+$"
      message = "Compatibility versions must be in ScalVer format"
    }
  }
  
  scalver_format = {
    type = "string"
    required = false
    pattern = "^0\\.\\d{8}\\.\\d+$"
    description = "ScalVer format specification for this schema"
    validation = {
      rule = "regex"
      pattern = "^0\\.\\d{8}\\.\\d+$"
      message = "ScalVer format must be in 0.YYYYMMDD.N format"
    }
  }
  
  # Schema lifecycle information
  lifecycle = {
    type = "object"
    required = false
    description = "Schema lifecycle information"
    
    properties = {
      status = {
        type = "string"
        required = false
        enum = ["draft", "active", "deprecated", "retired"]
        description = "Schema lifecycle status"
      }
      
      deprecation_date = {
        type = "string"
        required = false
        format = "date"
        description = "Date when schema will be deprecated"
      }
      
      retirement_date = {
        type = "string"
        required = false
        format = "date"
        description = "Date when schema will be retired"
      }
      
      migration_path = {
        type = "string"
        required = false
        description = "Migration path to newer schema version"
      }
    }
  }
  
  # Schema authorship and ownership
  authorship = {
    type = "object"
    required = false
    description = "Schema authorship and ownership information"
    
    properties = {
      author = {
        type = "string"
        required = false
        description = "Primary schema author"
      }
      
      maintainers = {
        type = "array"
        required = false
        description = "List of schema maintainers"
        items = {
          type = "string"
        }
      }
      
      contributors = {
        type = "array"
        required = false
        description = "List of schema contributors"
        items = {
          type = "string"
        }
      }
      
      license = {
        type = "string"
        required = false
        description = "Schema license information"
      }
      
      copyright = {
        type = "string"
        required = false
        description = "Schema copyright information"
      }
    }
  }
  
  # Schema dependencies and relationships
  dependencies = {
    type = "object"
    required = false
    description = "Schema dependencies and relationships"
    
    properties = {
      requires = {
        type = "array"
        required = false
        description = "Schemas that this schema requires"
        items = {
          type = "string"
          pattern = "^[a-z][a-z0-9-]*$"
        }
      }
      
      extends = {
        type = "array"
        required = false
        description = "Schemas that this schema extends"
        items = {
          type = "string"
          pattern = "^[a-z][a-z0-9-]*$"
        }
      }
      
      conflicts = {
        type = "array"
        required = false
        description = "Schemas that conflict with this schema"
        items = {
          type = "string"
          pattern = "^[a-z][a-z0-9-]*$"
        }
      }
      
      replaces = {
        type = "array"
        required = false
        description = "Schemas that this schema replaces"
        items = {
          type = "string"
          pattern = "^[a-z][a-z0-9-]*$"
        }
      }
    }
  }
  
  # Schema categorization
  categorization = {
    type = "object"
    required = false
    description = "Schema categorization and classification"
    
    properties = {
      category = {
        type = "string"
        required = false
        enum = ["core", "template", "facts", "machines", "actions", "variables", "configuration", "utility"]
        description = "Schema category"
      }
      
      tags = {
        type = "array"
        required = false
        description = "Schema tags for classification"
        items = {
          type = "string"
          pattern = "^[a-z][a-z0-9_-]*$"
        }
      }
      
      priority = {
        type = "string"
        required = false
        enum = ["critical", "high", "medium", "low"]
        description = "Schema priority level"
      }
      
      stability = {
        type = "string"
        required = false
        enum = ["experimental", "unstable", "stable", "frozen"]
        description = "Schema stability level"
      }
    }
  }
  
  # Schema documentation
  documentation = {
    type = "object"
    required = false
    description = "Schema documentation information"
    
    properties = {
      examples = {
        type = "array"
        required = false
        description = "Example usage files"
        items = {
          type = "string"
          pattern = "^examples/.*\\.hcl$"
        }
      }
      
      tutorials = {
        type = "array"
        required = false
        description = "Tutorial documentation files"
        items = {
          type = "string"
          pattern = "^docs/.*\\.md$"
        }
      }
      
      api_reference = {
        type = "string"
        required = false
        description = "API reference documentation URL"
      }
      
      changelog = {
        type = "string"
        required = false
        description = "Changelog file path"
      }
    }
  }
  
  # Schema validation rules
  validation = {
    # Required field validation
    required_fields = {
      rule = "required"
      fields = ["schema_version", "schema_type", "schema_name", "last_updated", "description"]
      message = "All required metadata fields must be present"
    }
    
    # Schema version format validation
    schema_version_format = {
      rule = "regex"
      field = "schema_version"
      pattern = "^0\\.\\d{8}\\.\\d+$"
      message = "Schema version must be in ScalVer format: 0.YYYYMMDD.N"
    }
    
    # Schema type format validation
    schema_type_format = {
      rule = "regex"
      field = "schema_type"
      pattern = "^[a-z][a-z0-9-]*$"
      message = "Schema type must be lowercase with hyphens only"
    }
    
    # Date format validation
    date_format = {
      rule = "regex"
      field = "last_updated"
      pattern = "^\\d{4}-\\d{2}-\\d{2}$"
      message = "Last updated must be in YYYY-MM-DD format"
    }
    
    # Description length validation
    description_length = {
      rule = "length"
      field = "description"
      min = 10
      max = 500
      message = "Description must be between 10 and 500 characters"
    }
    
    # Compatibility version validation
    compatibility_format = {
      rule = "conditional"
      condition = "compatibility != null"
      validation = {
        rule = "array_items"
        pattern = "^0\\.\\d{8}\\.\\d+$"
        message = "All compatibility versions must be in ScalVer format"
      }
    }
    
    # No circular dependencies
    no_circular_deps = {
      rule = "acyclic"
      field = "dependencies.requires"
      message = "Schema dependencies cannot be circular"
    }
    
    # Valid enum values
    valid_enums = {
      rule = "conditional"
      condition = "lifecycle.status != null"
      validation = {
        rule = "enum"
        enum = ["draft", "active", "deprecated", "retired"]
        message = "Lifecycle status must be one of: draft, active, deprecated, retired"
      }
    }
    
    # Tag format validation
    tag_format = {
      rule = "conditional"
      condition = "categorization.tags != null"
      validation = {
        rule = "array_items"
        pattern = "^[a-z][a-z0-9_-]*$"
        message = "Tags must be lowercase with underscores and hyphens only"
      }
    }
  }
}
