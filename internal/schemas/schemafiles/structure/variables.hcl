# Variables Structure Schema
# Common variable structure definitions for all storage formats
# This file defines the structure of variables that can be stored
# Used by HCL and JSON storage schemas for variables

# Schema metadata
metadata {
  schema_version = "0.20250809.2"
  schema_type = "variables-structure"
  schema_name = "Variables Structure Schema"
  last_updated = "2024-01-01"
  compatibility = ["0.20250809.0", "0.20250809.1", "0.20250809.2"]
  description = "Common variable structure definitions for all storage formats - defines the structure of variables that can be stored with age encryption support"
  
  # ScalVer format: 0.YYYYMMDD.N
  # - 0: Development phase
  # - 20250809: Date (9 August 2025)
  # - 2: Patch version (reorganized for validation system)
  scalver_format = "0.20250809.2"
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
  }
  
  # Variable description (optional)
  description {
    type = "string"
    required = false
    description = "Human-readable description of the variable"
  }
  
  # Variable value (required)
  # Can be string, number, boolean, or object
  value {
    type = "any"
    required = true
    description = "Variable value - can be string, number, boolean, or object"
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
  }
  
  # Encryption flag (optional)
  # Indicates if the variable value should be encrypted
  # This is the ONLY way to control encryption
  encrypted {
    type = "bool"
    required = false
    default = false
    description = "Whether the variable value should be encrypted using age encryption. This is the ONLY way to control encryption."
  }
  
  # Age encryption metadata (optional, only when encrypted = true)
  encryption_metadata {
    type = "object"
    required = false
    description = "Age encryption metadata - only present when encrypted = true"
    
    # Encryption metadata structure
    structure {
      # Recipients list (required for encryption)
      recipients {
        type = "list"
        required = true
        description = "List of age public keys that can decrypt this value"
        
        items {
          type = "string"
          description = "Age public key starting with 'age1'"
        }
      }
      
      # Encryption timestamp (optional)
      encrypted_at {
        type = "string"
        required = false
        description = "ISO 8601 timestamp when the value was encrypted"
      }
      
      # Encryption method (optional)
      method {
        type = "string"
        required = false
        default = "age"
        description = "Encryption method used"
      }
    }
  }
  
  # Variable tags (optional)
  tags {
    type = "list"
    required = false
    description = "Tags for categorizing and filtering variables"
    
    items {
      type = "string"
      description = "Tag value"
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
      }
      
      # Created timestamp (optional)
      created_at {
        type = "string"
        required = false
        description = "ISO 8601 timestamp when the variable was created"
      }
      
      # Last modified timestamp (optional)
      modified_at {
        type = "string"
        required = false
        description = "ISO 8601 timestamp when the variable was last modified"
      }
      
      # Version information (optional)
      version {
        type = "string"
        required = false
        description = "Version identifier for the variable"
      }
    }
  }
} 