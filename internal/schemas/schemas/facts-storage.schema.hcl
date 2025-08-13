# Facts Storage Schema
# Memory storage specifics for facts
# This schema defines storage constraints and features for facts in memory format

# Facts storage in memory format
facts_storage {
  # Include the base facts structure
  include = "facts-structure"
  
  # Memory-specific metadata
  storage_type = "memory"
  storage_location = "in-memory"
  description = "Facts stored in memory format"
  
  # Memory-specific features
  memory_features = {
    # Memory allocation
    allocation = {
      enabled = true
      strategy = "dynamic"
      description = "Dynamic memory allocation for fact storage"
    }
    
    # Memory pooling
    pooling = {
      enabled = false
      description = "Memory pooling for efficient allocation"
      
      properties = {
        pool_size = {
          type = "integer"
          description = "Size of memory pool in bytes"
        }
        
        pre_allocate = {
          type = "boolean"
          description = "Pre-allocate memory pool"
        }
      }
    }
    
    # Memory limits
    limits = {
      enabled = true
      description = "Memory usage limits for fact storage"
      
      properties = {
        max_memory = {
          type = "string"
          value = "1GB"
          description = "Maximum memory usage for facts"
        }
        
        max_entries = {
          type = "integer"
          value = 10000
          description = "Maximum number of fact entries"
        }
      }
    }
    
    # Memory cleanup
    cleanup = {
      enabled = true
      description = "Automatic memory cleanup"
      
      properties = {
        gc_interval = {
          type = "string"
          value = "5m"
          description = "Garbage collection interval"
        }
        
        cleanup_on_exit = {
          type = "boolean"
          value = true
          description = "Clean up memory on application exit"
        }
      }
    }
  }
  
  # Memory-specific validation
  memory_validation = {
    # Valid memory key format
    key_format = {
      rule = "memory_key_format"
      message = "Memory keys must be valid machine IDs"
    }
    
    # Valid memory value format
    value_format = {
      rule = "memory_value_format"
      message = "Memory values must be valid fact collections"
    }
    
    # Valid memory size
    size_limits = {
      rule = "memory_size_limits"
      message = "Memory usage must be within configured limits"
    }
    
    # Valid memory cleanup
    cleanup_validation = {
      rule = "memory_cleanup"
      message = "Memory cleanup must be performed regularly"
    }
  }
  
  # Memory-specific constraints
  memory_constraints = {
    # Memory storage constraints
    storage_constraints = {
      type = "object"
      description = "Memory storage constraints"
      
      properties = {
        max_concurrent_access = {
          type = "integer"
          value = 100
          description = "Maximum concurrent access to memory storage"
        }
        
        thread_safety = {
          type = "boolean"
          value = true
          description = "Thread-safe memory operations"
        }
        
        atomic_operations = {
          type = "boolean"
          value = true
          description = "Atomic memory operations"
        }
      }
    }
    
    # Memory performance constraints
    performance_constraints = {
      type = "object"
      description = "Memory performance constraints"
      
      properties = {
        read_latency = {
          type = "string"
          value = "1ms"
          description = "Maximum read latency"
        }
        
        write_latency = {
          type = "string"
          value = "1ms"
          description = "Maximum write latency"
        }
        
        memory_efficiency = {
          type = "number"
          value = 0.9
          description = "Minimum memory efficiency ratio"
        }
      }
    }
    
    # Memory security constraints
    security_constraints = {
      type = "object"
      description = "Memory security constraints"
      
      properties = {
        data_isolation = {
          type = "boolean"
          value = true
          description = "Data isolation between different fact collections"
        }
        
        access_control = {
          type = "boolean"
          value = true
          description = "Access control for memory operations"
        }
        
        memory_protection = {
          type = "boolean"
          value = true
          description = "Memory protection against corruption"
        }
      }
    }
  }
  
  # Memory-specific structure extensions
  memory_extensions = {
    # Memory key structure
    key_structure = {
      type = "object"
      description = "Memory key structure for facts"
      
      properties = {
        machine_id = {
          type = "string"
          pattern = "^[a-f0-9]{32}$"
          description = "32-character hexadecimal machine ID"
        }
        
        collection_timestamp = {
          type = "string"
          format = "datetime"
          description = "Timestamp when facts were collected"
        }
        
        version = {
          type = "string"
          pattern = "^[0-9]+\\.[0-9]+\\.[0-9]+$"
          description = "Fact collection version"
        }
      }
    }
    
    # Memory value structure
    value_structure = {
      type = "object"
      description = "Memory value structure for facts"
      
      properties = {
        fact_collection = {
          type = "object"
          description = "Complete fact collection data"
          include = "facts-structure"
        }
        
        metadata = {
          type = "object"
          description = "Additional metadata for the fact collection"
          
          properties = {
            collection_method = {
              type = "string"
              description = "Method used to collect facts"
            }
            
            collector_version = {
              type = "string"
              description = "Version of the fact collector"
            }
            
            validation_status = {
              type = "string"
              enum = ["valid", "invalid", "unknown"]
              description = "Validation status of the facts"
            }
          }
        }
      }
    }
    
    # Memory query structure
    query_structure = {
      type = "object"
      description = "Memory query structure for facts"
      
      properties = {
        query_type = {
          type = "string"
          enum = ["exact", "prefix", "range", "filter"]
          description = "Type of memory query"
        }
        
        query_parameters = {
          type = "object"
          description = "Parameters for the memory query"
          
          properties = {
            machine_id = {
              type = "string"
              description = "Machine ID to query"
            }
            
            time_range = {
              type = "object"
              description = "Time range for query"
              
              properties = {
                start = {
                  type = "string"
                  format = "datetime"
                  description = "Start time for query"
                }
                
                end = {
                  type = "string"
                  format = "datetime"
                  description = "End time for query"
                }
              }
            }
            
            filters = {
              type = "object"
              description = "Additional filters for query"
            }
          }
        }
      }
    }
    
    # Memory export structure
    export_structure = {
      type = "object"
      description = "Memory export structure"
      
      properties = {
        export_format = {
          type = "string"
          enum = ["json", "hcl", "yaml"]
          description = "Format for exported facts"
        }
        
        export_data = {
          type = "object"
          description = "Exported fact data"
          
          properties = {
            exported_at = {
              type = "string"
              format = "datetime"
              description = "Timestamp when facts were exported"
            }
            
            facts = {
              type = "object"
              description = "Exported facts data"
            }
          }
        }
      }
    }
  }
} 