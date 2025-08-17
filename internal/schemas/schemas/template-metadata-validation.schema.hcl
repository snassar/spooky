# Template Metadata Validation Schema
# This schema defines the actual structure and validation rules for template metadata
# Used for validating template metadata against the enhanced structure

# Schema metadata
metadata {
  schema_version = "0.20250115.0"
  schema_type = "template-metadata-validation"
  schema_name = "Template Metadata Validation Schema"
  last_updated = "2025-01-15"
  description = "Validation schema for enhanced template metadata structure with categorization, search optimization, and analytics fields"
  compatibility = ["0.20250115.0"]
}

# Template metadata structure definition
template_metadata {
  # Basic metadata fields
  name = {
    type = "string"
    required = true
    pattern = "^[a-zA-Z0-9._-]+$"
    min_length = 1
    max_length = 100
    description = "Template name (required, alphanumeric with dots, underscores, and hyphens)"
  }
  
  description = {
    type = "string"
    required = false
    max_length = 1000
    description = "Template description"
  }
  
  author = {
    type = "string"
    required = false
    pattern = "^[a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+\\.[a-zA-Z]{2,}$"
    description = "Template author (email format)"
  }
  
  version = {
    type = "string"
    required = false
    pattern = "^\\d+\\.\\d+\\.\\d+(-[a-zA-Z0-9._-]+)?(\\+[a-zA-Z0-9._-]+)?$"
    description = "Template version (semantic versioning format)"
  }
  
  tags = {
    type = "array"
    required = false
    max_items = 20
    items = {
      type = "string"
      pattern = "^[a-zA-Z0-9._-]+$"
      min_length = 1
      max_length = 50
    }
    description = "Template tags (max 20 tags, alphanumeric with dots, underscores, and hyphens)"
  }
  
  license = {
    type = "string"
    required = false
    enum = ["MIT", "Apache-2.0", "GPL-3.0", "BSD-3-Clause", "CC0", "Custom"]
    description = "Template license"
  }
  
  # Enhanced categorization fields
  category = {
    type = "string"
    required = false
    pattern = "^[a-zA-Z0-9._-]+$"
    enum = ["deployment", "configuration", "script", "documentation", "data", "template", "other"]
    description = "Template category"
  }
  
  subcategory = {
    type = "string"
    required = false
    pattern = "^[a-zA-Z0-9._-]+$"
    description = "Template subcategory"
  }
  
  priority = {
    type = "integer"
    required = false
    minimum = 1
    maximum = 10
    default = 5
    description = "Template priority (1-10, higher is more important)"
  }
  
  # Search optimization fields
  keywords = {
    type = "array"
    required = false
    max_items = 50
    items = {
      type = "string"
      pattern = "^[a-zA-Z0-9._-]+$"
      min_length = 1
      max_length = 50
    }
    description = "Search keywords (max 50 keywords)"
  }
  
  # Dependency fields
  dependencies = {
    type = "array"
    required = false
    max_items = 50
    items = {
      type = "string"
      pattern = "^[a-zA-Z0-9._-]+$"
    }
    description = "Template dependencies (max 50 dependencies)"
  }
  
  compatibility = {
    type = "object"
    required = false
    additional_properties = {
      type = "string"
      pattern = "^[><=!]*\\d+\\.\\d+\\.\\d+$"
    }
    description = "Compatibility requirements (e.g., 'spooky': '>=1.0.0')"
  }
  
  # Analytics fields
  usage_count = {
    type = "integer"
    required = false
    minimum = 0
    description = "Number of times template has been used"
  }
  
  last_used = {
    type = "string"
    required = false
    format = "date-time"
    description = "Last usage timestamp (RFC3339 format)"
  }
  
  # Lifecycle fields
  created_at = {
    type = "string"
    required = false
    format = "date-time"
    description = "Creation timestamp (RFC3339 format)"
  }
  
  updated_at = {
    type = "string"
    required = false
    format = "date-time"
    description = "Last update timestamp (RFC3339 format)"
  }
  
  # Status fields
  status = {
    type = "string"
    required = false
    enum = ["active", "deprecated", "removed"]
    default = "active"
    description = "Template status"
  }
  
  # Quality fields
  quality_score = {
    type = "number"
    required = false
    minimum = 0.0
    maximum = 1.0
    description = "Template quality score (0.0-1.0)"
  }
  
  # Security fields
  security_level = {
    type = "string"
    required = false
    enum = ["restricted", "standard", "elevated", "trusted"]
    default = "standard"
    description = "Template security level"
  }
  
  # Documentation fields
  documentation_url = {
    type = "string"
    required = false
    format = "uri"
    description = "URL to template documentation"
  }
  
  examples = {
    type = "array"
    required = false
    max_items = 10
    items = {
      type = "object"
      properties = {
        name = {
          type = "string"
          required = true
          description = "Example name"
        }
        description = {
          type = "string"
          required = false
          description = "Example description"
        }
        input = {
          type = "object"
          required = false
          description = "Example input data"
        }
        output = {
          type = "string"
          required = false
          description = "Example output"
        }
      }
    }
    description = "Template usage examples (max 10 examples)"
  }
}

# Validation rules
validation_rules {
  # Cross-field validation
  cross_field_validation = {
    # If category is specified, subcategory should also be specified
    category_subcategory_consistency = {
      condition = "if category is present, subcategory should also be present"
      severity = "warning"
    }
    
    # If dependencies are specified, they should be valid template names
    dependency_validation = {
      condition = "dependencies should reference valid template names"
      severity = "error"
    }
    
    # Version should be compatible with spooky version if specified
    version_compatibility = {
      condition = "version should be compatible with spooky version"
      severity = "warning"
    }
  }
  
  # Content validation
  content_validation = {
    # Description should not contain sensitive information
    sensitive_data_check = {
      pattern = "(?i)(password|secret|key|token|credential)"
      severity = "warning"
      message = "Description may contain sensitive information"
    }
    
    # Tags should be meaningful
    tag_meaningfulness = {
      min_tag_length = 2
      severity = "warning"
      message = "Tags should be at least 2 characters long"
    }
  }
  
  # Performance validation
  performance_validation = {
    # Metadata size should be reasonable
    metadata_size = {
      max_size_bytes = 1048576  # 1MB
      severity = "error"
      message = "Metadata size exceeds maximum allowed size"
    }
    
    # Number of keywords should be reasonable
    keyword_count = {
      max_keywords = 50
      severity = "warning"
      message = "Too many keywords may impact search performance"
    }
  }
}
