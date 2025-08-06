# Facts Storage Schema
# BadgerDB storage specifics for facts
# This schema defines storage constraints and features for facts in BadgerDB format

# Facts storage in BadgerDB format
facts_storage {
  # Include the base facts structure
  include = "facts-structure"
  
  # BadgerDB-specific metadata
  storage_type = "badgerdb"
  storage_location = "<project>/facts.db"
  description = "Facts stored in BadgerDB format"
  
  # BadgerDB-specific features
  badger_features = {
    # BadgerDB compression
    compression = {
      enabled = true
      algorithm = "zstd"
      level = 3
      description = "BadgerDB compression for fact storage"
    }
    
    # BadgerDB encryption
    encryption = {
      enabled = true
      algorithm = "age"
      description = "Age encryption for sensitive facts"
      
      properties = {
        recipients = {
          type = "array"
          description = "Age recipient keys for encryption"
        }
        
        identity_file = {
          type = "string"
          description = "Path to age identity file"
        }
      }
    }
    
    # BadgerDB transactions
    transactions = {
      enabled = true
      description = "ACID transactions for fact operations"
      
      properties = {
        read_only = {
          type = "boolean"
          description = "Read-only transaction mode"
        }
        
        serializable = {
          type = "boolean"
          description = "Serializable isolation level"
        }
      }
    }
    
    # BadgerDB indexing
    indexing = {
      enabled = true
      description = "Automatic indexing for fact queries"
      
      properties = {
        prefix_indexing = {
          type = "boolean"
          description = "Prefix-based indexing"
        }
        
        value_indexing = {
          type = "boolean"
          description = "Value-based indexing"
        }
      }
    }
    
    # BadgerDB garbage collection
    garbage_collection = {
      enabled = true
      description = "Automatic garbage collection"
      
      properties = {
        gc_interval = {
          type = "string"
          value = "1h"
          description = "Garbage collection interval"
        }
        
        gc_discard_ratio = {
          type = "number"
          value = 0.1
          description = "Garbage collection discard ratio"
        }
      }
    }
  }
  
  # BadgerDB-specific validation
  badger_validation = {
    # Valid BadgerDB key format
    valid_key_format = {
      rule = "regex"
      pattern = "^facts:[a-f0-9]{32}$"
      message = "BadgerDB keys must be in format: facts:<machine_id>"
    }
    
    # Valid BadgerDB value format
    valid_value_format = {
      rule = "json"
      message = "BadgerDB values must be valid JSON"
    }
    
    # Valid BadgerDB transaction
    valid_transaction = {
      rule = "badger_transaction"
      message = "BadgerDB operations must be within valid transactions"
    }
    
    # Valid BadgerDB encryption
    valid_encryption = {
      rule = "age_encryption"
      message = "Encrypted facts must use valid age encryption"
    }
    
    # Valid BadgerDB compression
    valid_compression = {
      rule = "zstd_compression"
      message = "Compressed facts must use valid zstd compression"
    }
    
    # Valid BadgerDB key size
    valid_key_size = {
      rule = "range"
      min = 1
      max = 65536
      message = "BadgerDB keys must be between 1 and 65536 bytes"
    }
    
    # Valid BadgerDB value size
    valid_value_size = {
      rule = "range"
      min = 1
      max = 1048576
      message = "BadgerDB values must be between 1 and 1MB"
    }
  }
  
  # BadgerDB-specific constraints
  badger_constraints = {
    # BadgerDB storage constraints
    storage_constraints = {
      type = "object"
      description = "BadgerDB storage constraints"
      
      properties = {
        max_db_size = {
          type = "integer"
          value = 1073741824
          description = "Maximum database size in bytes (1GB)"
        }
        
        max_value_size = {
          type = "integer"
          value = 1048576
          description = "Maximum value size in bytes (1MB)"
        }
        
        max_key_size = {
          type = "integer"
          value = 65536
          description = "Maximum key size in bytes (64KB)"
        }
        
        max_entries = {
          type = "integer"
          value = 1000000
          description = "Maximum number of entries"
        }
      }
    }
    
    # BadgerDB performance constraints
    performance_constraints = {
      type = "object"
      description = "BadgerDB performance constraints"
      
      properties = {
        max_concurrent_reads = {
          type = "integer"
          value = 100
          description = "Maximum concurrent read operations"
        }
        
        max_concurrent_writes = {
          type = "integer"
          value = 10
          description = "Maximum concurrent write operations"
        }
        
        max_transaction_size = {
          type = "integer"
          value = 10000
          description = "Maximum transaction size in entries"
        }
        
        max_memory_usage = {
          type = "integer"
          value = 1073741824
          description = "Maximum memory usage in bytes (1GB)"
        }
      }
    }
    
    # BadgerDB encryption constraints
    encryption_constraints = {
      type = "object"
      description = "BadgerDB encryption constraints"
      
      properties = {
        encryption_enabled = {
          type = "boolean"
          value = true
          description = "Whether encryption is enabled"
        }
        
        encryption_algorithm = {
          type = "string"
          value = "age"
          description = "Encryption algorithm (age)"
        }
        
        max_recipients = {
          type = "integer"
          value = 10
          description = "Maximum number of encryption recipients"
        }
        
        key_rotation = {
          type = "boolean"
          value = true
          description = "Whether key rotation is supported"
        }
      }
    }
  }
  
  # BadgerDB-specific structure extensions
  badger_extensions = {
    # BadgerDB key structure
    key_structure = {
      type = "object"
      description = "BadgerDB key structure for facts"
      
      properties = {
        prefix = {
          type = "string"
          value = "facts:"
          description = "Key prefix for facts"
        }
        
        machine_id = {
          type = "string"
          pattern = "^[a-f0-9]{32}$"
          description = "32-character machine ID"
        }
        
        key_format = {
          type = "string"
          value = "facts:<machine_id>"
          description = "Complete key format"
        }
      }
    }
    
    # BadgerDB value structure
    value_structure = {
      type = "object"
      description = "BadgerDB value structure for facts"
      
      properties = {
        fact_collection = {
          type = "object"
          description = "Serialized fact collection"
          
          properties = {
            machine_id = { type = "string" }
            collected_at = { type = "string" }
            facts = { type = "object" }
            encrypted = { type = "boolean" }
            compression = { type = "string" }
          }
        }
        
        metadata = {
          type = "object"
          description = "Value metadata"
          
          properties = {
            version = { type = "string" }
            created_at = { type = "string" }
            updated_at = { type = "string" }
            checksum = { type = "string" }
          }
        }
      }
    }
    
    # BadgerDB query structure
    query_structure = {
      type = "object"
      description = "BadgerDB query structure for facts"
      
      properties = {
        prefix_queries = {
          type = "object"
          description = "Prefix-based queries"
          
          properties = {
            all_facts = {
              type = "string"
              value = "facts:"
              description = "Query all facts"
            }
            
            machine_facts = {
              type = "string"
              value = "facts:<machine_id>"
              description = "Query facts for specific machine"
            }
          }
        }
        
        range_queries = {
          type = "object"
          description = "Range-based queries"
          
          properties = {
            time_range = {
              type = "object"
              description = "Time-based range queries"
            }
            
            value_range = {
              type = "object"
              description = "Value-based range queries"
            }
          }
        }
      }
    }
    
    # BadgerDB backup structure
    backup_structure = {
      type = "object"
      description = "BadgerDB backup structure"
      
      properties = {
        backup_format = {
          type = "string"
          value = "badger"
          description = "Backup format"
        }
        
        backup_compression = {
          type = "string"
          value = "gzip"
          description = "Backup compression"
        }
        
        backup_encryption = {
          type = "string"
          value = "age"
          description = "Backup encryption"
        }
        
        backup_metadata = {
          type = "object"
          description = "Backup metadata"
          
          properties = {
            created_at = { type = "string" }
            version = { type = "string" }
            checksum = { type = "string" }
            size = { type = "integer" }
          }
        }
      }
    }
  }
} 