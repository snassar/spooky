# Machines Schema
# Schema for machine inventory and connectivity
# Defines the structure and validation rules for machines

# Schema metadata
metadata {
  version = "1"
  description = "Machine inventory schema for spooky host definitions"
}

# Machine inventory structure (single machine)
machines {
  # Individual machine blocks
  machine {
    # Machine name
    name = {
      type = "string"
      required = true
      pattern = "^[a-zA-Z0-9_.-]+$"
      min_length = 1
      max_length = 64
      description = "Name of the machine"
    }
    
    # Machine description
    description = {
      type = "string"
      required = false
      max_length = 256
      description = "Description of the machine"
    }
    
    # Hostname or IP address
    hostname = {
      type = "string"
      required = true
      min_length = 1
      max_length = 253
      description = "Hostname or IP address of the machine"
    }
    
    # SSH port
    port = {
      type = "integer"
      required = false
      min = 1
      max = 65535
      default = 22
      description = "SSH port number"
    }
    
    # SSH user
    user = {
      type = "string"
      required = true
      pattern = "^[a-zA-Z0-9_.-]+$"
      min_length = 1
      max_length = 32
      description = "SSH username"
    }
    
    # SSH Authentication configuration
    authentication = {
      type = "object"
      required = true
      description = "SSH authentication configuration"
      
      # Authentication method
      method = {
        type = "string"
        required = true
        enum = ["ssh_key", "password", "certificate"]
        description = "SSH authentication method"
      }
      
      # SSH Key authentication fields
      ssh_key_path = {
        type = "string"
        required = false
        description = "Path to SSH private key file (required when method = 'ssh_key')"
      }
      
      # SSH Key passphrase (can be encrypted)
      ssh_key_passphrase = {
        type = "object"
        required = false
        description = "SSH key passphrase - can be encrypted using age encryption"
        
        # Passphrase value
        value = {
          type = "string"
          required = true
          description = "SSH key passphrase value - can be plain text or age-encrypted"
        }
        
        # Encryption flag
        encrypted = {
          type = "bool"
          required = false
          default = false
          description = "Whether the passphrase is encrypted using age encryption"
        }
        
        # Age encryption metadata (optional, only when encrypted = true)
        encryption_metadata = {
          type = "object"
          required = false
          description = "Age encryption metadata - only present when encrypted = true"
          
          # Recipients list (required for encryption)
          recipients = {
            type = "list"
            required = true
            description = "List of age public keys that can decrypt this value"
            
            items = {
              type = "string"
              description = "Age public key starting with 'age1'"
            }
          }
          
          # Encryption timestamp (optional)
          encrypted_at = {
            type = "string"
            required = false
            description = "ISO 8601 timestamp when the value was encrypted"
          }
          
          # Encryption method (optional)
          method = {
            type = "string"
            required = false
            default = "age"
            description = "Encryption method used"
          }
        }
      }
      
      # Password authentication (can be encrypted)
      password = {
        type = "object"
        required = false
        description = "SSH password - can be encrypted using age encryption"
        
        # Password value
        value = {
          type = "string"
          required = true
          description = "SSH password value - can be plain text or age-encrypted"
        }
        
        # Encryption flag
        encrypted = {
          type = "bool"
          required = false
          default = false
          description = "Whether the password is encrypted using age encryption"
        }
        
        # Age encryption metadata (optional, only when encrypted = true)
        encryption_metadata = {
          type = "object"
          required = false
          description = "Age encryption metadata - only present when encrypted = true"
          
          # Recipients list (required for encryption)
          recipients = {
            type = "list"
            required = true
            description = "List of age public keys that can decrypt this value"
            
            items = {
              type = "string"
              description = "Age public key starting with 'age1'"
            }
          }
          
          # Encryption timestamp (optional)
          encrypted_at = {
            type = "string"
            required = false
            description = "ISO 8601 timestamp when the value was encrypted"
          }
          
          # Encryption method (optional)
          method = {
            type = "string"
            required = false
            default = "age"
            description = "Encryption method used"
          }
        }
      }
      
      # Certificate authentication fields
      certificate_path = {
        type = "string"
        required = false
        description = "Path to SSH certificate file (required when method = 'certificate')"
      }
      
      certificate_key_path = {
        type = "string"
        required = false
        description = "Path to SSH certificate private key file (required when method = 'certificate')"
      }
      
      # Certificate passphrase (can be encrypted) - protects the certificate file itself
      certificate_passphrase = {
        type = "object"
        required = false
        description = "SSH certificate passphrase - can be encrypted using age encryption"
        
        # Passphrase value
        value = {
          type = "string"
          required = true
          description = "SSH certificate passphrase value - can be plain text or age-encrypted"
        }
        
        # Encryption flag
        encrypted = {
          type = "bool"
          required = false
          default = false
          description = "Whether the certificate passphrase is encrypted using age encryption"
        }
        
        # Age encryption metadata (optional, only when encrypted = true)
        encryption_metadata = {
          type = "object"
          required = false
          description = "Age encryption metadata - only present when encrypted = true"
          
          # Recipients list (required for encryption)
          recipients = {
            type = "list"
            required = true
            description = "List of age public keys that can decrypt this value"
            
            items = {
              type = "string"
              description = "Age public key starting with 'age1'"
            }
          }
          
          # Encryption timestamp (optional)
          encrypted_at = {
            type = "string"
            required = false
            description = "ISO 8601 timestamp when the value was encrypted"
          }
          
          # Encryption method (optional)
          method = {
            type = "string"
            required = false
            default = "age"
            description = "Encryption method used"
          }
        }
      }
      
      # Certificate private key passphrase (can be encrypted) - protects the private key used with the certificate
      certificate_key_passphrase = {
        type = "object"
        required = false
        description = "SSH certificate private key passphrase - can be encrypted using age encryption"
        
        # Passphrase value
        value = {
          type = "string"
          required = true
          description = "SSH certificate private key passphrase value - can be plain text or age-encrypted"
        }
        
        # Encryption flag
        encrypted = {
          type = "bool"
          required = false
          default = false
          description = "Whether the certificate private key passphrase is encrypted using age encryption"
        }
        
        # Age encryption metadata (optional, only when encrypted = true)
        encryption_metadata = {
          type = "object"
          required = false
          description = "Age encryption metadata - only present when encrypted = true"
          
          # Recipients list (required for encryption)
          recipients = {
            type = "list"
            required = true
            description = "List of age public keys that can decrypt this value"
            
            items = {
              type = "string"
              description = "Age public key starting with 'age1'"
            }
          }
          
          # Encryption timestamp (optional)
          encrypted_at = {
            type = "string"
            required = false
            description = "ISO 8601 timestamp when the value was encrypted"
          }
          
          # Encryption method (optional)
          method = {
            type = "string"
            required = false
            default = "age"
            description = "Encryption method used"
          }
        }
      }
    }
    
    # Connection configuration
    connection = {
      type = "object"
      required = false
      description = "SSH connection configuration"
      
      # Connection timeout
      timeout = {
        type = "integer"
        required = false
        min = 1
        max = 300
        default = 30
        description = "SSH connection timeout in seconds"
      }
      
      # Retry attempts
      retry_attempts = {
        type = "integer"
        required = false
        min = 0
        max = 10
        default = 3
        description = "Number of connection retry attempts"
      }
      
      # Retry delay
      retry_delay = {
        type = "integer"
        required = false
        min = 1
        max = 60
        default = 5
        description = "Delay between retry attempts in seconds"
      }
    }
    
    # Machine tags
    tags = {
      type = "array"
      required = false
      max_items = 20
      description = "Tags for categorizing machines"
      items = {
        type = "string"
        pattern = "^[a-zA-Z0-9_-]+$"
        min_length = 1
        max_length = 32
      }
    }
    
    # Machine groups
    groups = {
      type = "array"
      required = false
      max_items = 10
      description = "Groups for targeting machines (Spooky creates groups dynamically)"
      items = {
        type = "string"
        pattern = "^[a-zA-Z0-9_-]+$"
        min_length = 1
        max_length = 32
      }
    }
    
    # Environment information
    environment = {
      type = "string"
      required = false
      enum = ["development", "staging", "production", "testing", "qa"]
      description = "Environment name (e.g., 'development', 'staging', 'production')"
    }
    
    # Role information
    role = {
      type = "string"
      required = false
      max_length = 64
      description = "Machine role (e.g., 'web', 'database', 'load-balancer')"
    }
    
    # Location information
    location = {
      type = "string"
      required = false
      max_length = 64
      description = "Machine location (e.g., 'us-east-1', 'eu-west-1')"
    }
    
    # Last accessed timestamp
    last_accessed = {
      type = "string"
      required = false
      pattern = "^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$"
      description = "ISO 8601 timestamp when the machine was last accessed"
    }
    
    # Machine-specific variables (behaves like variables.hcl)
    variables = {
      type = "object"
      required = false
      description = "Machine-specific variables that override project variables"
      
      # Variable structure follows the same pattern as variables.hcl
      structure = {
        # Variable name
        name = {
          type = "string"
          required = true
          description = "Variable name - must be a valid identifier"
        }
        
        # Variable description
        description = {
          type = "string"
          required = false
          description = "Human-readable description of the variable"
        }
        
        # Variable value
        value = {
          type = "any"
          required = true
          description = "Variable value - can be string, number, boolean, or object"
        }
        
        # Sensitive flag
        sensitive = {
          type = "bool"
          required = false
          default = false
          description = "Whether the variable contains sensitive data - used for logging, display, and documentation. Does NOT control encryption."
        }
        
        # Encryption flag
        encrypted = {
          type = "bool"
          required = false
          default = false
          description = "Whether the variable value should be encrypted using age encryption. This is the ONLY way to control encryption."
        }
        
        # Age encryption metadata (optional, only when encrypted = true)
        encryption_metadata = {
          type = "object"
          required = false
          description = "Age encryption metadata - only present when encrypted = true"
          
          # Recipients list (required for encryption)
          recipients = {
            type = "list"
            required = true
            description = "List of age public keys that can decrypt this value"
            
            items = {
              type = "string"
              description = "Age public key starting with 'age1'"
            }
          }
          
          # Encryption timestamp (optional)
          encrypted_at = {
            type = "string"
            required = false
            description = "ISO 8601 timestamp when the value was encrypted"
          }
          
          # Encryption method (optional)
          method = {
            type = "string"
            required = false
            default = "age"
            description = "Encryption method used"
          }
        }
        
        # Variable tags
        tags = {
          type = "array"
          required = false
          max_items = 10
          description = "Tags for categorizing and filtering variables"
          items = {
            type = "string"
            pattern = "^[a-zA-Z0-9_-]+$"
            min_length = 1
            max_length = 32
            description = "Tag value"
          }
        }
        
        # Variable metadata
        metadata = {
          type = "object"
          required = false
          description = "Additional metadata for the variable"
          
          # Flexible metadata structure
          structure = {
            # Source information (optional)
            source = {
              type = "string"
              required = false
              description = "Source of the variable (e.g., 'environment', 'file', 'manual')"
            }
            
            # Created timestamp (optional)
            created_at = {
              type = "string"
              required = false
              description = "ISO 8601 timestamp when the variable was created"
            }
            
            # Last modified timestamp (optional)
            modified_at = {
              type = "string"
              required = false
              description = "ISO 8601 timestamp when the variable was last modified"
            }
            
            # Version information (optional)
            version = {
              type = "string"
              required = false
              description = "Version identifier for the variable"
            }
          }
        }
      }
    }
  }
  
  # Group-level variables (applied to all machines in specific groups)
  group_variables = {
    type = "object"
    required = false
    description = "Group-level variables that apply to machines in specific groups"
    
    # Group variable structure
    structure = {
      # Group name
      group_name = {
        type = "string"
        required = true
        pattern = "^[a-zA-Z0-9_-]+$"
        min_length = 1
        max_length = 32
        description = "Name of the group these variables apply to"
      }
      
      # Variables for this group (follows same structure as machine variables)
      variables = {
        type = "object"
        required = true
        description = "Variables that apply to all machines in this group"
        
        # Variable structure follows the same pattern as variables.hcl
        structure = {
          # Variable name
          name = {
            type = "string"
            required = true
            description = "Variable name - must be a valid identifier"
          }
          
          # Variable description
          description = {
            type = "string"
            required = false
            description = "Human-readable description of the variable"
          }
          
          # Variable value
          value = {
            type = "any"
            required = true
            description = "Variable value - can be string, number, boolean, or object"
          }
          
          # Sensitive flag
          sensitive = {
            type = "bool"
            required = false
            default = false
            description = "Whether the variable contains sensitive data - used for logging, display, and documentation. Does NOT control encryption."
          }
          
          # Encryption flag
          encrypted = {
            type = "bool"
            required = false
            default = false
            description = "Whether the variable value should be encrypted using age encryption. This is the ONLY way to control encryption."
          }
          
          # Age encryption metadata (optional, only when encrypted = true)
          encryption_metadata = {
            type = "object"
            required = false
            description = "Age encryption metadata - only present when encrypted = true"
            
            # Recipients list (required for encryption)
            recipients = {
              type = "list"
              required = true
              description = "List of age public keys that can decrypt this value"
              
              items = {
                type = "string"
                description = "Age public key starting with 'age1'"
              }
            }
            
            # Encryption timestamp (optional)
            encrypted_at = {
              type = "string"
              required = false
              description = "ISO 8601 timestamp when the value was encrypted"
            }
            
            # Encryption method (optional)
            method = {
              type = "string"
              required = false
              default = "age"
              description = "Encryption method used"
            }
          }
          
          # Variable tags
          tags = {
            type = "array"
            required = false
            max_items = 10
            description = "Tags for categorizing and filtering variables"
            items = {
              type = "string"
              pattern = "^[a-zA-Z0-9_-]+$"
              min_length = 1
              max_length = 32
              description = "Tag value"
            }
          }
          
          # Variable metadata
          metadata = {
            type = "object"
            required = false
            description = "Additional metadata for the variable"
            
                      # Flexible metadata structure
          structure = {
              # Source information (optional)
              source = {
                type = "string"
                required = false
                description = "Source of the variable (e.g., 'environment', 'file', 'manual')"
              }
              
              # Created timestamp (optional)
              created_at = {
                type = "string"
                required = false
                description = "ISO 8601 timestamp when the variable was created"
              }
              
              # Last modified timestamp (optional)
              modified_at = {
                type = "string"
                required = false
                description = "ISO 8601 timestamp when the variable was last modified"
              }
              
              # Version information (optional)
              version = {
                type = "string"
                required = false
                description = "Version identifier for the variable"
              }
            }
          }
        }
      }
    }
  }
}

 