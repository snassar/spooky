# Machines Schema
# Schema for machine inventory and connectivity
# Defines the structure and validation rules for machines

# Schema metadata
metadata {
  schema_version = "0.20250809.1"
  schema_type = "machines"
  schema_name = "Machines Schema"
  last_updated = "2024-01-01"
  compatibility = ["0.20250809.0", "0.20250809.1"]
  description = "Schema for machine inventory and connectivity with age encryption support"
  
  # ScalVer format: 0.YYYYMMDD.N
  # - 0: Development phase
  # - 20250809: Date (9 August 2025)
  # - 1: Patch version (added age encryption support)
  scalver_format = "0.20250809.1"
}

# Machine inventory structure
machines_inventory {
  # Machine list
  machines {
    type = "list"
    required = true
    description = "List of machines in the inventory"
    
    validation {
      min_items = 0
      max_items = 10000
    }
    
    # Individual machine structure
    items {
      type = "object"
      required = true
      description = "Individual machine with connectivity and authentication"
      
      structure {
        # Machine name
        name {
          type = "string"
          required = true
          description = "Name of the machine"
          
          validation {
            pattern = "^[a-zA-Z0-9_.-]+$"
            min_length = 1
            max_length = 64
          }
        }
        
        # Machine description
        description {
          type = "string"
          required = false
          description = "Description of the machine"
          
          validation {
            max_length = 256
          }
        }
        
        # Hostname or IP address
        hostname {
          type = "string"
          required = true
          description = "Hostname or IP address of the machine"
          
          validation {
            # Allow hostnames and IP addresses
            # Application-level validation for format
            min_length = 1
            max_length = 253
          }
        }
        
        # SSH port
        port {
          type = "integer"
          required = false
          default = 22
          description = "SSH port number"
          
          validation {
            min = 1
            max = 65535
          }
        }
        
        # SSH user
        user {
          type = "string"
          required = true
          description = "SSH username"
          
          validation {
            pattern = "^[a-zA-Z0-9_.-]+$"
            min_length = 1
            max_length = 32
          }
        }
        
        # Authentication configuration
        authentication {
          type = "object"
          required = true
          description = "SSH authentication configuration with age encryption support"
          
          structure {
            # Authentication method
            method {
              type = "string"
              required = true
              description = "Authentication method"
              
              validation {
                allowed_values = ["ssh_key", "password", "passphrase", "certificate", "agent"]
              }
            }
            
            # SSH key path (for ssh_key method)
            key_path {
              type = "string"
              required = false
              description = "Path to SSH private key file"
              
              validation {
                # Path validation handled at application level
                min_length = 1
                max_length = 512
              }
            }
            
            # SSH key passphrase (optional, can be age-encrypted)
            passphrase {
              type = "object"
              required = false
              description = "SSH key passphrase with age encryption support"
              
              structure {
                # Passphrase value (can be plaintext or age-encrypted)
                value {
                  type = "string"
                  required = true
                  description = "Passphrase value - can be plaintext or age-encrypted string"
                  
                  validation {
                    # Application-level validation for age encryption
                    min_length = 1
                    max_length = 1024
                  }
                }
                
                # Encryption flag
                encrypted {
                  type = "bool"
                  required = false
                  default = false
                  description = "Whether the passphrase is age-encrypted"
                  
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
              }
            }
            
            # Password (for password method, can be age-encrypted)
            password {
              type = "object"
              required = false
              description = "SSH password with age encryption support"
              
              structure {
                # Password value (can be plaintext or age-encrypted)
                value {
                  type = "string"
                  required = true
                  description = "Password value - can be plaintext or age-encrypted string"
                  
                  validation {
                    # Application-level validation for age encryption
                    min_length = 1
                    max_length = 1024
                  }
                }
                
                # Encryption flag
                encrypted {
                  type = "bool"
                  required = false
                  default = false
                  description = "Whether the password is age-encrypted"
                  
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
              }
            }
            
            # Certificate path (for certificate method)
            certificate_path {
              type = "string"
              required = false
              description = "Path to SSH certificate file"
              
              validation {
                # Path validation handled at application level
                min_length = 1
                max_length = 512
              }
            }
            
            # Certificate key path (for certificate method)
            certificate_key_path {
              type = "string"
              required = false
              description = "Path to SSH certificate private key file"
              
              validation {
                # Path validation handled at application level
                min_length = 1
                max_length = 512
              }
            }
            
            # Certificate key passphrase (optional, can be age-encrypted)
            certificate_key_passphrase {
              type = "object"
              required = false
              description = "SSH certificate key passphrase with age encryption support"
              
              structure {
                # Passphrase value (can be plaintext or age-encrypted)
                value {
                  type = "string"
                  required = true
                  description = "Passphrase value - can be plaintext or age-encrypted string"
                  
                  validation {
                    # Application-level validation for age encryption
                    min_length = 1
                    max_length = 1024
                  }
                }
                
                # Encryption flag
                encrypted {
                  type = "bool"
                  required = false
                  default = false
                  description = "Whether the passphrase is age-encrypted"
                  
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
              }
            }
          }
        }
        
        # Connection settings
        connection {
          type = "object"
          required = false
          description = "SSH connection settings"
          
          structure {
            # Connection timeout
            timeout {
              type = "string"
              required = false
              default = "30s"
              description = "SSH connection timeout (e.g., '30s', '2m')"
              
              validation {
                pattern = "^\\d+[smhd]$"
              }
            }
            
            # Keep-alive interval
            keepalive_interval {
              type = "string"
              required = false
              default = "30s"
              description = "SSH keep-alive interval (e.g., '30s', '2m')"
              
              validation {
                pattern = "^\\d+[smhd]$"
              }
            }
            
            # Max retries
            max_retries {
              type = "integer"
              required = false
              default = 3
              description = "Maximum number of connection retries"
              
              validation {
                min = 0
                max = 10
              }
            }
            
            # Retry delay
            retry_delay {
              type = "string"
              required = false
              default = "5s"
              description = "Delay between retries (e.g., '5s', '1m')"
              
              validation {
                pattern = "^\\d+[smhd]$"
              }
            }
            
            # Host key verification
            host_key_verification {
              type = "bool"
              required = false
              default = true
              description = "Whether to verify host keys"
            }
            
            # Known hosts file
            known_hosts_file {
              type = "string"
              required = false
              description = "Path to known hosts file"
              
              validation {
                # Path validation handled at application level
                min_length = 1
                max_length = 512
              }
            }
          }
        }
        
        # Machine tags
        tags {
          type = "list"
          required = false
          description = "Tags for categorizing machines"
          
          validation {
            max_items = 20
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
        
        # Machine metadata
        metadata {
          type = "object"
          required = false
          description = "Additional metadata for the machine"
          
          structure {
            # Environment information
            environment {
              type = "string"
              required = false
              description = "Environment name (e.g., 'development', 'staging', 'production')"
              
              validation {
                allowed_values = ["development", "staging", "production", "testing", "qa"]
              }
            }
            
            # Role information
            role {
              type = "string"
              required = false
              description = "Machine role (e.g., 'web', 'database', 'load-balancer')"
              
              validation {
                max_length = 64
              }
            }
            
            # Location information
            location {
              type = "string"
              required = false
              description = "Machine location (e.g., 'us-east-1', 'eu-west-1')"
              
              validation {
                max_length = 64
              }
            }
            
            # Created timestamp
            created_at {
              type = "string"
              required = false
              description = "ISO 8601 timestamp when the machine was added"
              
              validation {
                pattern = "^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$"
              }
            }
            
            # Last accessed timestamp
            last_accessed {
              type = "string"
              required = false
              description = "ISO 8601 timestamp when the machine was last accessed"
              
              validation {
                pattern = "^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$"
              }
            }
            
            # Version information
            version {
              type = "string"
              required = false
              description = "Version of the machine configuration"
              
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
  
  # Inventory metadata
  metadata {
    type = "object"
    required = false
    description = "Metadata about the machine inventory"
    
    structure {
      # Inventory name
      name {
        type = "string"
        required = false
        description = "Name of the machine inventory"
        
        validation {
          pattern = "^[a-zA-Z0-9_.-]+$"
          min_length = 1
          max_length = 64
        }
      }
      
      # Inventory description
      description {
        type = "string"
        required = false
        description = "Description of the machine inventory"
        
        validation {
          max_length = 256
        }
      }
      
      # Created timestamp
      created_at {
        type = "string"
        required = false
        description = "ISO 8601 timestamp when the inventory was created"
        
        validation {
          pattern = "^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$"
        }
      }
      
      # Last updated timestamp
      updated_at {
        type = "string"
        required = false
        description = "ISO 8601 timestamp when the inventory was last updated"
        
        validation {
          pattern = "^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$"
        }
      }
      
      # Version information
      version {
        type = "string"
        required = false
        description = "Version of the inventory"
        
        validation {
          pattern = "^[a-zA-Z0-9._-]+$"
          max_length = 32
        }
      }
    }
  }
}

# Validation rules for machines
validation_rules {
  # Cross-field validation rules
  cross_field_validation {
    # Authentication method validation
    rule {
      name = "ssh_key_method_validation"
      description = "SSH key method requires key_path"
      condition = "authentication.method == 'ssh_key' && authentication.key_path == null"
      message = "SSH key authentication requires key_path"
    }
    
    rule {
      name = "password_method_validation"
      description = "Password method requires password"
      condition = "authentication.method == 'password' && authentication.password == null"
      message = "Password authentication requires password"
    }
    
    rule {
      name = "certificate_method_validation"
      description = "Certificate method requires certificate_path and certificate_key_path"
      condition = "authentication.method == 'certificate' && (authentication.certificate_path == null || authentication.certificate_key_path == null)"
      message = "Certificate authentication requires certificate_path and certificate_key_path"
    }
    
    # Encryption metadata validation
    rule {
      name = "passphrase_encryption_metadata_required"
      description = "Encryption metadata must be present when passphrase is encrypted"
      condition = "authentication.passphrase.encrypted == true && authentication.passphrase.encryption_metadata == null"
      message = "Encryption metadata is required when passphrase is encrypted"
    }
    
    rule {
      name = "password_encryption_metadata_required"
      description = "Encryption metadata must be present when password is encrypted"
      condition = "authentication.password.encrypted == true && authentication.password.encryption_metadata == null"
      message = "Encryption metadata is required when password is encrypted"
    }
    
    rule {
      name = "certificate_passphrase_encryption_metadata_required"
      description = "Encryption metadata must be present when certificate passphrase is encrypted"
      condition = "authentication.certificate_key_passphrase.encrypted == true && authentication.certificate_key_passphrase.encryption_metadata == null"
      message = "Encryption metadata is required when certificate passphrase is encrypted"
    }
    
    # Age-encrypted values must start with 'age1'
    rule {
      name = "passphrase_age_encrypted_format"
      description = "Age-encrypted passphrases must start with 'age1'"
      condition = "authentication.passphrase.encrypted == true && !authentication.passphrase.value.startsWith('age1')"
      message = "Age-encrypted passphrases must start with 'age1'"
    }
    
    rule {
      name = "password_age_encrypted_format"
      description = "Age-encrypted passwords must start with 'age1'"
      condition = "authentication.password.encrypted == true && !authentication.password.value.startsWith('age1')"
      message = "Age-encrypted passwords must start with 'age1'"
    }
    
    rule {
      name = "certificate_passphrase_age_encrypted_format"
      description = "Age-encrypted certificate passphrases must start with 'age1'"
      condition = "authentication.certificate_key_passphrase.encrypted == true && !authentication.certificate_key_passphrase.value.startsWith('age1')"
      message = "Age-encrypted certificate passphrases must start with 'age1'"
    }
  }
  
  # Application-level validation notes
  application_validation {
    # These validations require application-level logic
    # Schema validation cannot enforce all type relationships
    
    note = "Application must validate that encrypted = true is only used with age-encrypted strings"
    note = "Application must validate age public key format and validity"
    note = "Application must validate that encrypted values are valid age-encrypted strings"
    note = "Application must handle decryption during authentication"
    note = "Application must validate SSH key file permissions and format"
    note = "Application must validate certificate file permissions and format"
    note = "Application must validate hostname/IP address format"
    note = "Application must validate file paths exist and are accessible"
  }
  
  # Age encryption specific rules
  age_encryption_rules {
    # Age1 prefix detection for authentication values
    rule {
      name = "age1_prefix_detection_passphrase"
      description = "Detect age1 prefix for passphrases"
      condition = "authentication.passphrase.value.startsWith('age1') && authentication.passphrase.encrypted == false"
      message = "Values starting with 'age1' should be marked as encrypted = true"
    }
    
    rule {
      name = "age1_prefix_detection_password"
      description = "Detect age1 prefix for passwords"
      condition = "authentication.password.value.startsWith('age1') && authentication.password.encrypted == false"
      message = "Values starting with 'age1' should be marked as encrypted = true"
    }
    
    rule {
      name = "age1_prefix_detection_certificate_passphrase"
      description = "Detect age1 prefix for certificate passphrases"
      condition = "authentication.certificate_key_passphrase.value.startsWith('age1') && authentication.certificate_key_passphrase.encrypted == false"
      message = "Values starting with 'age1' should be marked as encrypted = true"
    }
  }
} 