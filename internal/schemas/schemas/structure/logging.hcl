# Spooky Logging Configuration Schema
# Comprehensive schema for logging configuration with best practices for readability
# Based on log formatting best practices from Sematext

# Schema metadata
metadata {
  schema_version = "0.20250812.0"
  schema_type = "logging"
  schema_name = "Spooky Logging Configuration Schema"
  last_updated = "2024-01-01"
  compatibility = ["0.20250812.0"]
  description = "Comprehensive schema for logging configuration with structured output and best practices"
  
  # ScalVer format: 0.YYYYMMDD.N
  # - 0: Development phase
  # - 20250812: Date (12 August 2025)
  # - 0: Patch version
  scalver_format = "0.20250812.0"
}

# Logging configuration block
logging {
  # Log level configuration
  level = {
    type = "string"
    required = false
    default = "info"
    enum = ["debug", "info", "warn", "error", "fatal"]
    description = "Minimum log level to output (debug, info, warn, error, fatal)"
  }
  
  # Output format configuration
  format = {
    type = "string"
    required = false
    default = "json"
    enum = ["json", "text", "structured"]
    description = "Log output format (json for machine-readable, text for human-readable, structured for custom)"
  }
  
  # Output destination configuration
  output = {
    type = "string"
    required = false
    default = "stderr"
    enum = ["stdout", "stderr", "file", "null"]
    description = "Log output destination (stdout, stderr, file, null for no output)"
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
