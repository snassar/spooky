# Spooky Global Configuration Schema
# Schema for $XDG_CONFIG_HOME/spooky/spooky.hcl

# Schema metadata
metadata {
  schema_version = "0.20250809.0"
  schema_type = "spooky"
  schema_name = "Spooky Global Configuration Schema"
  last_updated = "2024-01-01"
  compatibility = ["0.20250809.0"]
  description = "Global configuration schema for spooky CLI - defines configuration structure for $XDG_CONFIG_HOME/spooky/spooky.hcl"
  
  # ScalVer format: 0.YYYYMMDD.N
  # - 0: Development phase
  # - 20250809: Date (9 August 2025)
  # - 0: Patch version
  scalver_format = "0.20250809.0"
}

# Global configuration block structure
spooky {
  # Storage configuration
  storage {
    type = {
      type = "string"
      required = false
      enum = ["memory", "json"]
      default = "memory"
      description = "Storage backend type (memory for in-memory storage)"
    }
    
    path = {
      type = "string"
      required = false
      default = ""
      description = "Storage path for facts export (not used for memory storage)"
    }
    
    compression = {
      type = "boolean"
      required = false
      default = true
      description = "Enable compression for storage"
    }
    
    encryption = {
      type = "boolean"
      required = false
      default = false
      description = "Enable age encryption for sensitive data"
    }
    
    backup_enabled = {
      type = "boolean"
      required = false
      default = true
      description = "Enable automatic backups"
    }
    
    backup_retention = {
      type = "integer"
      required = false
      min = 1
      max = 365
      default = 7
      description = "Number of backup files to retain"
    }
  }

  # Facts collection configuration
  facts {
    timeout = {
      type = "integer"
      required = false
      min = 1
      max = 3600
      default = 30
      description = "Timeout for facts collection in seconds"
    }
    
    cache_ttl = {
      type = "integer"
      required = false
      min = 0
      max = 86400
      default = 3600
      description = "Cache TTL for facts in seconds (0 = no cache)"
    }
    
    auto_collect = {
      type = "boolean"
      required = false
      default = false
      description = "Automatically collect facts on machine operations"
    }
    
    parallel_collection = {
      type = "integer"
      required = false
      min = 1
      max = 100
      default = 10
      description = "Number of parallel facts collection workers"
    }
    
    retry_attempts = {
      type = "integer"
      required = false
      min = 0
      max = 10
      default = 3
      description = "Number of retry attempts for failed facts collection"
    }
    
    retry_delay = {
      type = "integer"
      required = false
      min = 1
      max = 60
      default = 5
      description = "Delay between retry attempts in seconds"
    }
  }

  # SSH configuration
  ssh {
    timeout = {
      type = "integer"
      required = false
      min = 1
      max = 300
      default = 30
      description = "SSH connection timeout in seconds"
    }
    
    keepalive_interval = {
      type = "integer"
      required = false
      min = 1
      max = 300
      default = 60
      description = "SSH keepalive interval in seconds"
    }
    
    keepalive_count = {
      type = "integer"
      required = false
      min = 1
      max = 10
      default = 3
      description = "SSH keepalive count before considering connection dead"
    }
    
    key_scan_timeout = {
      type = "integer"
      required = false
      min = 1
      max = 60
      default = 10
      description = "SSH key scanning timeout in seconds"
    }
    
    known_hosts_strict = {
      type = "boolean"
      required = false
      default = true
      description = "Strict known_hosts checking"
    }
    
    connection_pool_size = {
      type = "integer"
      required = false
      min = 1
      max = 100
      default = 10
      description = "SSH connection pool size"
    }
  }

  # Template configuration
  templates {
    max_size = {
      type = "integer"
      required = false
      min = 1024
      max = 10485760
      default = 1048576
      description = "Maximum template file size in bytes"
    }
    
    allow_external_functions = {
      type = "boolean"
      required = false
      default = false
      description = "Allow external function calls in templates"
    }
    
    timeout = {
      type = "integer"
      required = false
      min = 1
      max = 300
      default = 30
      description = "Template rendering timeout in seconds"
    }
    
    cache_compiled = {
      type = "boolean"
      required = false
      default = true
      description = "Cache compiled templates"
    }
    
    sandbox_mode = {
      type = "boolean"
      required = false
      default = true
      description = "Enable sandbox mode for template running"
    }
    
    allowed_functions = {
      type = "list"
      required = false
      default = ["len", "join", "split", "replace", "trim", "upper", "lower"]
      description = "List of allowed functions in templates"
    }
  }

  # Security configuration
  security {
    allow_unsafe_commands = {
      type = "boolean"
      required = false
      default = false
      description = "Allow potentially unsafe commands"
    }
    
    restrict_file_access = {
      type = "boolean"
      required = false
      default = true
      description = "Restrict file access to project directory"
    }
    
    validate_ssh_keys = {
      type = "boolean"
      required = false
      default = true
      description = "Validate SSH keys before use"
    }
    
    audit_logging = {
      type = "boolean"
      required = false
      default = true
      description = "Enable audit logging for security events"
    }
    
    allowed_hosts = {
      type = "list"
      required = false
      default = []
      description = "List of allowed host patterns (CIDR or hostnames)"
    }
    
    blocked_hosts = {
      type = "list"
      required = false
      default = []
      description = "List of blocked host patterns (CIDR or hostnames)"
    }
  }

  # Age encryption configuration
  age {
    identities = {
      type = "string"
      required = false
      default = "~/.config/spooky/identities"
      description = "Path to directory containing age identity files (default: ~/.config/spooky/identities)"
    }
    
    recipients = {
      type = "string"
      required = false
      default = "~/.config/spooky/recipients.txt"
      description = "Path to file containing age recipients (public keys, one per line)"
    }
    
    passphrase = {
      type = "string"
      required = false
      description = "Passphrase for age encryption (if not using identity files)"
    }
    
    validation = {
      type = "object"
      required = false
      default = {
        strict_mode = true
        check_recipients = true
        validate_keys = true
      }
      description = "Age validation settings"
      
      properties = {
        strict_mode = {
          type = "boolean"
          required = false
          default = true
          description = "Enable strict validation mode"
        }
        
        check_recipients = {
          type = "boolean"
          required = false
          default = true
          description = "Validate recipient keys on encryption"
        }
        
        validate_keys = {
          type = "boolean"
          required = false
          default = true
          description = "Validate identity keys on decryption"
        }
      }
    }
    
    encryption = {
      type = "object"
      required = false
      default = {
        algorithm = "age"
        compression = true
        armor = false
      }
      description = "Age encryption settings"
      
      properties = {
        algorithm = {
          type = "string"
          required = false
          enum = ["age"]
          default = "age"
          description = "Encryption algorithm (only age supported)"
        }
        
        compression = {
          type = "boolean"
          required = false
          default = true
          description = "Enable compression before encryption"
        }
        
        armor = {
          type = "boolean"
          required = false
          default = false
          description = "Use armored output format"
        }
      }
    }
  }

  # Logging configuration
  logging {
    # Basic logging configuration
    level = {
      type = "string"
      required = false
      enum = ["debug", "info", "warn", "error", "fatal"]
      default = "info"
      description = "Minimum log level to output (debug, info, warn, error, fatal)"
    }
    
    format = {
      type = "string"
      required = false
      enum = ["json", "text", "structured"]
      default = "json"
      description = "Log output format (json for machine-readable, text for human-readable, structured for custom)"
    }
    
    output = {
      type = "string"
      required = false
      enum = ["stdout", "stderr", "file", "null"]
      default = "stderr"
      description = "Log output destination (stdout, stderr, file, null for no output)"
    }
    
    color_output = {
      type = "boolean"
      required = false
      default = true
      description = "Enable colored output"
    }
    
    progress_bars = {
      type = "boolean"
      required = false
      default = true
      description = "Show progress bars for long operations"
    }
    
    # File output configuration
    file = {
      type = "object"
      required = false
      description = "File output configuration (required when output is 'file')"
      
      properties = {
        path = {
          type = "string"
          required = true
          description = "Path to log file"
          pattern = "^[^<>:\"/\\|?*]+$"
        }
        
        permissions = {
          type = "string"
          required = false
          default = "0644"
          pattern = "^[0-7]{3,4}$"
          description = "File permissions in octal format (e.g., 0644)"
        }
        
        append = {
          type = "boolean"
          required = false
          default = true
          description = "Whether to append to existing file or truncate"
        }
      }
    }
    
    # Structured logging configuration
    structured = {
      type = "object"
      required = false
      description = "Structured logging configuration for custom formatting"
      
      properties = {
        # Timestamp configuration
        timestamp = {
          type = "object"
          required = false
          description = "Timestamp formatting configuration"
          
          properties = {
            enabled = {
              type = "boolean"
              required = false
              default = true
              description = "Whether to include timestamps in log entries"
            }
            
            format = {
              type = "string"
              required = false
              default = "RFC3339"
              enum = ["RFC3339", "RFC3339Nano", "Unix", "UnixNano", "ISO8601"]
              description = "Timestamp format (RFC3339, RFC3339Nano, Unix, UnixNano, ISO8601)"
            }
            
            timezone = {
              type = "string"
              required = false
              default = "UTC"
              description = "Timezone for timestamps (e.g., UTC, America/New_York)"
            }
          }
        }
        
        # Level configuration
        level = {
          type = "object"
          required = false
          description = "Log level configuration"
          
          properties = {
            key = {
              type = "string"
              required = false
              default = "level"
              description = "Field key for log level"
            }
            
            uppercase = {
              type = "boolean"
              required = false
              default = true
              description = "Whether to use uppercase level names"
            }
            
            color = {
              type = "boolean"
              required = false
              default = false
              description = "Whether to include color codes in level output"
            }
          }
        }
        
        # Message configuration
        message = {
          type = "object"
          required = false
          description = "Message formatting configuration"
          
          properties = {
            key = {
              type = "string"
              required = false
              default = "msg"
              description = "Field key for log message"
            }
            
            truncate = {
              type = "integer"
              required = false
              min = 0
              max = 10000
              description = "Maximum message length (0 for no truncation)"
            }
          }
        }
        
        # Component configuration
        component = {
          type = "object"
          required = false
          description = "Component identification configuration"
          
          properties = {
            key = {
              type = "string"
              required = false
              default = "component"
              description = "Field key for component name"
            }
            
            enabled = {
              type = "boolean"
              required = false
              default = true
              description = "Whether to include component information"
            }
            
            include_package = {
              type = "boolean"
              required = false
              default = false
              description = "Whether to include package information in component"
            }
          }
        }
        
        # Operation configuration
        operation = {
          type = "object"
          required = false
          description = "Operation tracking configuration"
          
          properties = {
            key = {
              type = "string"
              required = false
              default = "operation"
              description = "Field key for operation name"
            }
            
            enabled = {
              type = "boolean"
              required = false
              default = true
              description = "Whether to include operation information"
            }
            
            include_id = {
              type = "boolean"
              required = false
              default = false
              description = "Whether to include operation ID for correlation"
            }
          }
        }
        
        # Error configuration
        error = {
          type = "object"
          required = false
          description = "Error handling configuration"
          
          properties = {
            key = {
              type = "string"
              required = false
              default = "error"
              description = "Field key for error information"
            }
            
            include_stack = {
              type = "boolean"
              required = false
              default = false
              description = "Whether to include stack traces with errors"
            }
            
            include_type = {
              type = "boolean"
              required = false
              default = true
              description = "Whether to include error type information"
            }
            
            include_code = {
              type = "boolean"
              required = false
              default = false
              description = "Whether to include error codes"
            }
          }
        }
        
        # Caller information configuration
        caller = {
          type = "object"
          required = false
          description = "Caller information configuration"
          
          properties = {
            enabled = {
              type = "boolean"
              required = false
              default = false
              description = "Whether to include caller information (file, line, function)"
            }
            
            key = {
              type = "string"
              required = false
              default = "caller"
              description = "Field key for caller information"
            }
            
            include_package = {
              type = "boolean"
              required = false
              default = false
              description = "Whether to include package information in caller"
            }
            
            skip_frames = {
              type = "integer"
              required = false
              default = 0
              min = 0
              max = 10
              description = "Number of stack frames to skip when determining caller"
            }
          }
        }
        
        # Fields configuration
        fields = {
          type = "object"
          required = false
          description = "Additional fields configuration"
          
          properties = {
            # Global fields to include in all log entries
            global = {
              type = "object"
              required = false
              description = "Global fields to include in all log entries"
              additional_properties = true
            }
            
            # Field ordering for consistent output
            order = {
              type = "array"
              required = false
              description = "Order of fields in log output"
              items = {
                type = "string"
              }
            }
            
            # Field filtering
            filter = {
              type = "object"
              required = false
              description = "Field filtering configuration"
              
              properties = {
                exclude = {
                  type = "array"
                  required = false
                  description = "Field names to exclude from output"
                  items = {
                    type = "string"
                  }
                }
                
                include = {
                  type = "array"
                  required = false
                  description = "Field names to include in output (if empty, include all except excluded)"
                  items = {
                    type = "string"
                  }
                }
                
                sensitive = {
                  type = "array"
                  required = false
                  description = "Sensitive field names to mask or redact"
                  items = {
                    type = "string"
                  }
                }
              }
            }
          }
        }
      }
    }
    
    # Performance configuration
    performance = {
      type = "object"
      required = false
      description = "Performance optimization configuration"
      
      properties = {
        # Buffering configuration
        buffer = {
          type = "object"
          required = false
          description = "Log buffering configuration"
          
          properties = {
            enabled = {
              type = "boolean"
              required = false
              default = false
              description = "Whether to use buffered logging"
            }
            
            size = {
              type = "integer"
              required = false
              default = 4096
              min = 1024
              max = 1048576
              description = "Buffer size in bytes"
            }
            
            flush_interval = {
              type = "string"
              required = false
              default = "1s"
              description = "Flush interval (e.g., 1s, 100ms)"
            }
          }
        }
        
        # Async configuration
        async = {
          type = "object"
          required = false
          description = "Asynchronous logging configuration"
          
          properties = {
            enabled = {
              type = "boolean"
              required = false
              default = false
              description = "Whether to use asynchronous logging"
            }
            
            queue_size = {
              type = "integer"
              required = false
              default = 1000
              min = 100
              max = 100000
              description = "Queue size for async logging"
            }
            
            workers = {
              type = "integer"
              required = false
              default = 1
              min = 1
              max = 10
              description = "Number of worker goroutines for async logging"
            }
            
            drop_when_full = {
              type = "boolean"
              required = false
              default = false
              description = "Whether to drop logs when queue is full"
            }
          }
        }
      }
    }
    
    # Filtering configuration
    filtering = {
      type = "object"
      required = false
      description = "Log filtering configuration"
      
      properties = {
        # Component-based filtering
        components = {
          type = "object"
          required = false
          description = "Component-specific log level configuration"
          additional_properties = {
            type = "string"
            enum = ["debug", "info", "warn", "error", "fatal"]
          }
        }
        
        # Pattern-based filtering
        patterns = {
          type = "object"
          required = false
          description = "Pattern-based filtering configuration"
          
          properties = {
            include = {
              type = "array"
              required = false
              description = "Regex patterns to include (logs must match at least one)"
              items = {
                type = "string"
              }
            }
            
            exclude = {
              type = "array"
              required = false
              description = "Regex patterns to exclude (logs matching any are dropped)"
              items = {
                type = "string"
              }
            }
          }
        }
      }
    }
    
    # Rotation configuration
    rotation = {
      type = "object"
      required = false
      description = "Log file rotation configuration"
      
      properties = {
        enabled = {
          type = "boolean"
          required = false
          default = false
          description = "Whether to enable log file rotation"
        }
        
        max_size = {
          type = "string"
          required = false
          default = "100MB"
          description = "Maximum file size before rotation (e.g., 100MB, 1GB)"
        }
        
        max_age = {
          type = "string"
          required = false
          default = "30d"
          description = "Maximum age of rotated files (e.g., 30d, 7d, 24h)"
        }
        
        max_backups = {
          type = "integer"
          required = false
          default = 5
          min = 1
          max = 100
          description = "Maximum number of backup files to keep"
        }
        
        compress = {
          type = "boolean"
          required = false
          default = true
          description = "Whether to compress rotated log files"
        }
        
        local_time = {
          type = "boolean"
          required = false
          default = false
          description = "Whether to use local time for rotation timestamps"
        }
      }
    }
  }

  # Performance configuration
  performance {
    default_parallel = {
      type = "integer"
      required = false
      min = 2
      max = 100
      default = 10
      description = "Default parallel run limit"
    }
    
    max_memory = {
      type = "integer"
      required = false
      min = 64
      max = 8192
      default = 512
      description = "Maximum memory usage in MB"
    }
    
    gc_interval = {
      type = "integer"
      required = false
      min = 1
      max = 3600
      default = 300
      description = "Garbage collection interval in seconds"
    }
  }
  
} 