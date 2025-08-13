# Spooky Project Configuration Schema
# Comprehensive schema for project.hcl files with project metadata and structure validation

# Schema metadata
metadata {
  schema_version = "0.20250809.0"
  schema_type = "project"
  schema_name = "Spooky Project Configuration Schema"
  last_updated = "2024-01-01"
  compatibility = ["0.20250809.0"]
  description = "Comprehensive schema for project.hcl files with project metadata and structure validation"
  
  # ScalVer format: 0.YYYYMMDD.N
  # - 0: Development phase
  # - 20250809: Date (9 August 2025)
  # - 0: Patch version
  scalver_format = "0.20250809.0"
}

# Project block structure
project {
  # Project metadata
  name = {
    type = "string"
    required = true
    pattern = "^[a-zA-Z][a-zA-Z0-9._-]*$"
    min_length = 1
    max_length = 128
    description = "Project name (used for identification and isolation)"
  }
  
  description = {
    type = "string"
    required = false
    max_length = 1024
    description = "Project description"
  }
  
  version = {
    type = "string"
    required = false
    max_length = 64
    description = "Project version"
  }
  
  author = {
    type = "string"
    required = false
    max_length = 128
    description = "Project author or maintainer"
  }
  
  email = {
    type = "string"
    required = false
    format = "email"
    description = "Project contact email address"
  }
  
  url = {
    type = "string"
    required = false
    format = "uri"
    description = "Project website or documentation URL"
  }
  
  # Project run settings
run = {
    type = "object"
    required = false
          description = "Project run configuration"
    
    properties = {
      default_timeout = {
        type = "integer"
        required = false
        min = 1
        max = 3600
        default = 300
        description = "Default timeout for actions in seconds"
      }
      
      max_parallel = {
        type = "integer"
        required = false
        min = 1
        max = 100
        default = 10
        description = "Maximum parallel runs"
      }
      
      dry_run_default = {
        type = "boolean"
        required = false
        default = false
        description = "Default dry-run mode for actions"
      }
      
      validate_before_run = {
        type = "boolean"
        required = false
        default = true
        description = "Validate project configuration before running"
      }
      
      backup_before_changes = {
        type = "boolean"
        required = false
        default = false
        description = "Create backups before making changes"
      }
    }
  }
  
  # Facts collection configuration
  facts = {
    type = "object"
    required = false
    description = "Project-specific facts collection configuration"
    
    properties = {
      timeout = {
        type = "integer"
        required = false
        min = 1
        max = 3600
        default = 30
        description = "Timeout for facts collection in seconds"
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
      
      storage_format = {
        type = "string"
        required = false
        enum = ["memory", "json"]
        default = "memory"
        description = "Storage format for project facts (memory for in-memory storage)"
      }
      
      compression = {
        type = "boolean"
        required = false
        default = true
        description = "Enable compression for facts storage"
      }
      
      encryption = {
        type = "boolean"
        required = false
        default = false
        description = "Enable age encryption for sensitive facts data"
      }
    }
  }
  
  # Project metadata
  metadata = {
    type = "object"
    required = false
    description = "Additional project metadata"
    additional_properties = "string"
  }
  
  # Project-specific logging configuration (overrides global logging settings)
  logging = {
    type = "object"
    required = false
    description = "Project-specific logging configuration that overrides global settings"
    
    properties = {
      level = {
        type = "string"
        required = false
        enum = ["debug", "info", "warn", "error", "fatal"]
        description = "Log level for this project (overrides global setting)"
      }
      
      format = {
        type = "string"
        required = false
        enum = ["json", "text", "structured"]
        description = "Log format for this project (overrides global setting)"
      }
      
      output = {
        type = "string"
        required = false
        enum = ["stdout", "stderr", "file", "null"]
        description = "Log output destination for this project (overrides global setting)"
      }
      
      file = {
        type = "object"
        required = false
        description = "File output configuration for this project"
        
        properties = {
          path = {
            type = "string"
            required = false
            description = "Path to log file (relative to project directory or absolute)"
          }
          
          permissions = {
            type = "string"
            required = false
            pattern = "^[0-7]{3,4}$"
            default = "0644"
            description = "File permissions in octal format (e.g., 0644)"
          }
          
          append = {
            type = "boolean"
            required = false
            default = true
            description = "Append to existing log file instead of overwriting"
          }
        }
      }
      
      filtering = {
        type = "object"
        required = false
        description = "Component-specific filtering for this project"
        
        properties = {
          components = {
            type = "object"
            required = false
            additional_properties = "string"
            description = "Component-specific log levels (e.g., 'ssh' = 'debug')"
          }
          
          patterns = {
            type = "array"
            required = false
            items = "string"
            description = "Pattern-based filtering rules"
          }
        }
      }
      
      rotation = {
        type = "object"
        required = false
        description = "Log rotation configuration for this project"
        
        properties = {
          enabled = {
            type = "boolean"
            required = false
            default = false
            description = "Enable log rotation for this project"
          }
          
          max_size = {
            type = "string"
            required = false
            description = "Maximum log file size (e.g., '100MB', '1GB')"
          }
          
          max_age = {
            type = "string"
            required = false
            description = "Maximum log file age (e.g., '24h', '7d', '30d')"
          }
          
          max_backups = {
            type = "integer"
            required = false
            min = 1
            max = 100
            description = "Maximum number of backup files to keep"
          }
          
          compress = {
            type = "boolean"
            required = false
            default = true
            description = "Compress rotated log files"
          }
        }
      }
    }
  }
  
  # Validation rules
  validation = {
    # Project name validation
    project_name = {
      rule = "regex"
      pattern = "^[a-zA-Z][a-zA-Z0-9._-]*$"
      message = "Project names must start with a letter and contain only alphanumeric characters, dots, underscores, and hyphens"
    }
    
    # Email validation
    email_format = {
      rule = "format"
      format = "email"
      message = "Contact email must be a valid email address"
    }
    
    # URL validation
    url_format = {
      rule = "format"
      format = "uri"
      message = "Contact URL must be a valid URI"
    }
    
    # Run validation
    timeout_reasonable = {
      rule = "range"
      min = 1
      max = 3600
      message = "Timeouts must be between 1 and 3600 seconds"
    }
    
    parallel_limit = {
      rule = "range"
      min = 1
      max = 100
              message = "Maximum parallel runs must be between 1 and 100"
    }
    
    # Facts validation
    facts_timeout_reasonable = {
      rule = "range"
      min = 1
      max = 3600
      message = "Facts collection timeouts must be between 1 and 3600 seconds"
    }
    
    facts_parallel_reasonable = {
      rule = "range"
      min = 1
      max = 100
      message = "Facts parallel collection must be between 1 and 100 workers"
    }
    
    # Logging validation
    logging_level_valid = {
      rule = "enum"
      values = ["debug", "info", "warn", "error", "fatal"]
      message = "Log level must be one of: debug, info, warn, error, fatal"
    }
    
    logging_format_valid = {
      rule = "enum"
      values = ["json", "text", "structured"]
      message = "Log format must be one of: json, text, structured"
    }
    
    logging_output_valid = {
      rule = "enum"
      values = ["stdout", "stderr", "file", "null"]
      message = "Log output must be one of: stdout, stderr, file, null"
    }
    
    logging_file_permissions_valid = {
      rule = "regex"
      pattern = "^[0-7]{3,4}$"
      message = "File permissions must be in octal format (e.g., 0644)"
    }
    
    logging_rotation_backups_reasonable = {
      rule = "range"
      min = 1
      max = 100
      message = "Maximum log backups must be between 1 and 100"
    }
    
    logging_file_path_required = {
      rule = "required_if"
      condition = "output == 'file'"
      field = "logging.file.path"
      message = "File path is required when output is set to 'file'"
    }
  }
  
} 