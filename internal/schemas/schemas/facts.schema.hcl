# Facts Schema
# Schema for fact collection and storage
# Defines the structure and validation rules for facts

# Schema metadata
metadata {
  schema_version = "0.20250809.1"
  schema_type = "facts"
  schema_name = "Facts Schema"
  last_updated = "2024-01-01"
  compatibility = ["0.20250809.0", "0.20250809.1"]
  description = "Schema for fact collection and storage with age encryption support"
  
  # ScalVer format: 0.YYYYMMDD.N
  # - 0: Development phase
  # - 20250809: Date (9 August 2025)
  # - 1: Patch version (added age encryption support)
  scalver_format = "0.20250809.1"
}

# Facts collection structure
facts_collection {
  # Collection metadata
  metadata {
    type = "object"
    required = true
    description = "Metadata about the fact collection"
    
    structure {
      # Collection name
      name {
        type = "string"
        required = true
        description = "Name of the fact collection"
        
        validation {
          pattern = "^[a-zA-Z0-9_-]+$"
          min_length = 1
          max_length = 64
        }
      }
      
      # Collection description
      description {
        type = "string"
        required = false
        description = "Description of the fact collection"
        
        validation {
          max_length = 256
        }
      }
      
      # Collection timestamp
      timestamp {
        type = "string"
        required = true
        description = "ISO 8601 timestamp when facts were collected"
        
        validation {
          pattern = "^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$"
        }
      }
      
      # Source information
      source {
        type = "string"
        required = true
        description = "Source of the facts (e.g., 'ssh', 'local', 'file')"
        
        validation {
          allowed_values = ["ssh", "local", "file", "custom", "imported"]
        }
      }
      
      # Collection method
      method {
        type = "string"
        required = false
        description = "Method used to collect facts"
        
        validation {
          max_length = 64
        }
      }
      
      # Version information
      version {
        type = "string"
        required = false
        description = "Version of the fact collection"
        
        validation {
          pattern = "^[a-zA-Z0-9._-]+$"
          max_length = 32
        }
      }
    }
  }
  
  # Facts list
  facts {
    type = "list"
    required = true
    description = "List of collected facts"
    
    validation {
      min_items = 0
      max_items = 10000
    }
    
    # Individual fact structure
    items {
      type = "object"
      required = true
      description = "Individual fact with name, value, and metadata"
      
      structure {
        # Fact name
        name {
          type = "string"
          required = true
          description = "Name of the fact"
          
          validation {
            pattern = "^[a-zA-Z0-9_.-]+$"
            min_length = 1
            max_length = 128
          }
        }
        
        # Fact value
        value {
          type = "any"
          required = true
          description = "Value of the fact - can be string, number, boolean, object, or age-encrypted string"
          
          # Value validation handled at application level
          # Age-encrypted values (age1...) are valid strings
          validation {
            # Application-level validation for age encryption
            # Schema cannot enforce age format validation
          }
        }
        
        # Fact type
        type {
          type = "string"
          required = true
          description = "Type of the fact value"
          
          validation {
            allowed_values = ["string", "number", "boolean", "object", "array", "encrypted"]
          }
        }
        
        # Encryption flag
        encrypted {
          type = "bool"
          required = false
          default = false
          description = "Whether the fact value is age-encrypted"
          
          # Only applicable to string values that start with 'age1'
          validation {
            # Application-level validation ensures this is only true for age-encrypted strings
          }
        }
        
        # Age encryption metadata (optional, only when encrypted = true)
        encryption_metadata {
          type = "object"
          required = false
          description = "Age encryption metadata - only present when encrypted = true"
          
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
        
        # Fact description
        description {
          type = "string"
          required = false
          description = "Description of the fact"
          
          validation {
            max_length = 256
          }
        }
        
        # Fact tags
        tags {
          type = "list"
          required = false
          description = "Tags for categorizing facts"
          
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
        
        # Fact metadata
        metadata {
          type = "object"
          required = false
          description = "Additional metadata for the fact"
          
          structure {
            # Source information
            source {
              type = "string"
              required = false
              description = "Source of the fact (e.g., 'system', 'file', 'command')"
              
              validation {
                allowed_values = ["system", "file", "command", "custom", "imported"]
              }
            }
            
            # Collection timestamp
            collected_at {
              type = "string"
              required = false
              description = "ISO 8601 timestamp when the fact was collected"
              
              validation {
                pattern = "^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$"
              }
            }
            
            # Expiration timestamp
            expires_at {
              type = "string"
              required = false
              description = "ISO 8601 timestamp when the fact expires"
              
              validation {
                pattern = "^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$"
              }
            }
            
            # Priority level
            priority {
              type = "string"
              required = false
              description = "Priority level of the fact"
              
              validation {
                allowed_values = ["low", "normal", "high", "critical"]
              }
            }
            
            # Version information
            version {
              type = "string"
              required = false
              description = "Version of the fact"
              
              validation {
                pattern = "^[a-zA-Z0-9._-]+$"
                max_length = 32
              }
            }
          }
        }
      }
    }
  }
  
  # Collection statistics
  statistics {
    type = "object"
    required = false
    description = "Statistics about the fact collection"
    
    structure {
      # Total facts count
      total_facts {
        type = "integer"
        required = true
        description = "Total number of facts in the collection"
        
        validation {
          min = 0
        }
      }
      
      # Encrypted facts count
      encrypted_facts {
        type = "integer"
        required = false
        description = "Number of age-encrypted facts in the collection"
        
        validation {
          min = 0
        }
      }
      
      # Collection duration
      collection_duration {
        type = "string"
        required = false
        description = "Duration of fact collection (e.g., '1.5s', '2m30s')"
        
        validation {
          pattern = "^\\d+(\\.\\d+)?[smhd]$"
        }
      }
      
      # Error count
      error_count {
        type = "integer"
        required = false
        description = "Number of errors during fact collection"
        
        validation {
          min = 0
        }
      }
      
      # Warning count
      warning_count {
        type = "integer"
        required = false
        description = "Number of warnings during fact collection"
        
        validation {
          min = 0
        }
      }
    }
  }
}

# Validation rules for facts
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
    
    # Encrypted facts must have type = "encrypted"
    rule {
      name = "encrypted_type_consistency"
      description = "Encrypted facts must have type = 'encrypted'"
      condition = "encrypted == true && type != 'encrypted'"
      message = "Encrypted facts must have type = 'encrypted'"
    }
    
    # Non-encrypted facts cannot have type = "encrypted"
    rule {
      name = "non_encrypted_type_consistency"
      description = "Non-encrypted facts cannot have type = 'encrypted'"
      condition = "encrypted == false && type == 'encrypted'"
      message = "Non-encrypted facts cannot have type = 'encrypted'"
    }
    
    # Age-encrypted values must start with 'age1'
    rule {
      name = "age_encrypted_format"
      description = "Age-encrypted values must start with 'age1'"
      condition = "encrypted == true && !value.startsWith('age1')"
      message = "Age-encrypted values must start with 'age1'"
    }
  }
  
  # Application-level validation notes
  application_validation {
    # These validations require application-level logic
    # Schema validation cannot enforce all type relationships
    
    note = "Application must validate that encrypted = true is only used with age-encrypted strings"
    note = "Application must validate age public key format and validity"
    note = "Application must validate that encrypted values are valid age-encrypted strings"
    note = "Application must handle decryption during fact access"
    note = "Application must validate age1 prefix detection for custom facts"
    note = "Application must remove regex validation for age-encrypted values"
    note = "Application must integrate age encryption validation with fact collection"
  }
  
  # Age encryption specific rules
  age_encryption_rules {
    # Age1 prefix detection for custom facts
    rule {
      name = "age1_prefix_detection"
      description = "Detect age1 prefix for custom facts"
      condition = "value.startsWith('age1') && encrypted == false"
      message = "Values starting with 'age1' should be marked as encrypted = true"
    }
    
    # No regex validation for age-encrypted values
    rule {
      name = "no_regex_for_age_encrypted"
      description = "Do not apply regex validation to age-encrypted values"
      condition = "encrypted == true"
      message = "Age-encrypted values bypass regex validation"
    }
  }
}
