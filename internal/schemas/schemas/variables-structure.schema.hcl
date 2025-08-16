# Variables Structure Schema
# Common variable structure definitions for all storage formats
# This file defines the structure of variables that can be stored
# Used by HCL and JSON storage schemas for variables

# Schema metadata
metadata {
  schema_version = "0.20250809.0"
  schema_type = "variables-structure"
  schema_name = "Variables Structure Schema"
  last_updated = "2024-01-01"
  compatibility = ["0.20250809.0", "0.20250809.1"]
  description = "Common variable structure definitions for all storage formats - defines the structure of variables that can be stored with age encryption support"
  
  # ScalVer format: 0.YYYYMMDD.N
  # - 0: Development phase
  # - 20250809: Date (9 August 2025)
  # - 1: Patch version (added age encryption support)
  scalver_format = "0.20250809.1"
}

# Variable structure definition
# This defines the common structure for all variable types
variable_structure {
  # Variable name (required)
  # Must be a valid identifier
  name {
    type = "string"
    required = true
    description = "Variable name - must be a valid identifier"
    
    # Validation rules
    validation {
      # Must be a valid identifier
      pattern = "^[a-zA-Z_][a-zA-Z0-9_]*$"
      min_length = 1
      max_length = 64
    }
  }
  
  # Variable description (optional)
  description {
    type = "string"
    required = false
    description = "Human-readable description of the variable"
    
    validation {
      max_length = 256
    }
  }
  
  # Variable value (required)
  # Can be string, number, boolean, or object
  value {
    type = "any"
    required = true
    description = "Variable value - can be string, number, boolean, or object"
    
    # Value type validation
    validation {
      # Allow any type for flexibility
      # Type-specific validation handled by application
    }
  }
  
  # Sensitive flag (optional)
  # Indicates if the variable contains sensitive data
  # Used for logging, display, and documentation purposes
  # Does NOT control encryption - use encrypted flag for that
  sensitive {
    type = "bool"
    required = false
    default = false
    description = "Whether the variable contains sensitive data - used for logging, display, and documentation. Does NOT control encryption."
    
    # Sensitive variables are:
    # - Masked in logs and output
    # - Not displayed in plain text
    # - Marked in documentation as sensitive
    # - Still need encrypted = true to be encrypted
    validation {
      # No specific validation rules
      # Application-level validation ensures proper handling
    }
  }
  
  # Encryption flag (optional)
  # Indicates if the variable value should be encrypted
  # This is the ONLY way to control encryption
  encrypted {
    type = "bool"
    required = false
    default = false
    description = "Whether the variable value should be encrypted using age encryption. This is the ONLY way to control encryption."
    
    # Only applicable to string, number, and object values
    # Booleans cannot be encrypted (they're simple true/false values)
    validation {
      # Application-level validation ensures this is only set for encryptable types
      # Schema validation cannot enforce this cross-field relationship
    }
  }
  
  # Age encryption metadata (optional, only when encrypted = true)
  encryption_metadata {
    type = "object"
    required = false
    description = "Age encryption metadata - only present when encrypted = true"
    
    # Only validate if encrypted = true
    validation {
      # Cross-field validation: only present when encrypted = true
      # Application-level validation required
    }
    
    # Encryption metadata structure
    structure {
      # Recipients list (required for encryption)
      recipients {
        type = "list"
        required = true
        description = "List of age public keys that can decrypt this value"
        
        validation {
          min_items = 1
          max_items = 10
        }
        
        # Each recipient must be a valid age public key
        items {
          type = "string"
          validation {
            # Age public key format: age1...
            pattern = "^age1[a-z0-9]{50,}$"
            description = "Must be a valid age public key starting with 'age1'"
          }
        }
      }
      
      # Encryption timestamp (optional)
      encrypted_at {
        type = "string"
        required = false
        description = "ISO 8601 timestamp when the value was encrypted"
        
        validation {
          # ISO 8601 format
          pattern = "^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$"
        }
      }
      
      # Encryption method (optional)
      method {
        type = "string"
        required = false
        default = "age"
        description = "Encryption method used"
        
        validation {
          # Currently only age encryption is supported
          allowed_values = ["age"]
        }
      }
    }
  }
  
  # Variable tags (optional)
  tags {
    type = "list"
    required = false
    description = "Tags for categorizing and filtering variables"
    
    validation {
      max_items = 10
    }
    
    items {
      type = "string"
      validation {
        pattern = "^[a-zA-Z0-9_-]+$"
        min_length = 1
        max_length = 32
      }
    }
  }
  
  # Variable metadata (optional)
  metadata {
    type = "object"
    required = false
    description = "Additional metadata for the variable"
    
    # Flexible metadata structure
    structure {
      # Source information (optional)
      source {
        type = "string"
        required = false
        description = "Source of the variable (e.g., 'environment', 'file', 'manual')"
        
        validation {
          allowed_values = ["environment", "file", "manual", "computed", "imported"]
        }
      }
      
      # Created timestamp (optional)
      created_at {
        type = "string"
        required = false
        description = "ISO 8601 timestamp when the variable was created"
        
        validation {
          pattern = "^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$"
        }
      }
      
      # Last modified timestamp (optional)
      modified_at {
        type = "string"
        required = false
        description = "ISO 8601 timestamp when the variable was last modified"
        
        validation {
          pattern = "^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$"
        }
      }
      
      # Version information (optional)
      version {
        type = "string"
        required = false
        description = "Version identifier for the variable"
        
        validation {
          pattern = "^[a-zA-Z0-9._-]+$"
          max_length = 32
        }
      }
    }
  }
}

# Validation rules for variable structure
validation_rules {
  # Cross-field validation rules
  cross_field_validation {
    # Encryption metadata must be present when encrypted = true
    rule {
      name = "encryption_metadata_required"
      description = "Encryption metadata must be present when encrypted = true"
      condition = "encrypted == true && encryption_metadata == null"
      message = "Encryption metadata is required when encrypted = true"
    }
    
    # Encryption metadata must not be present when encrypted = false
    rule {
      name = "encryption_metadata_not_allowed"
      description = "Encryption metadata must not be present when encrypted = false"
      condition = "encrypted == false && encryption_metadata != null"
      message = "Encryption metadata is not allowed when encrypted = false"
    }
    
    # Value type validation for encryption
    rule {
      name = "encryptable_value_types"
      description = "Only string, number, and object values can be encrypted"
      condition = "encrypted == true && (value_type == 'bool')"
      message = "Boolean values cannot be encrypted"
    }
  }
  
  # Application-level validation notes
  application_validation {
    # These validations require application-level logic
    # Schema validation cannot enforce all type relationships
    
    note = "Application must validate that encrypted = true is only used with string, number, or object values"
    note = "Application must validate age public key format and validity"
    note = "Application must validate that encrypted values are valid age-encrypted strings"
    note = "Application must handle decryption and re-encryption during variable operations"
    note = "Sensitive variables should be masked in logs and output, regardless of encryption status"
    note = "Only encrypted = true controls encryption - sensitive = true is for display/logging only"
  }
} 