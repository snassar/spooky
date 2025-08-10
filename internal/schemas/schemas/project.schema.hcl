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
  
  # Project execution settings
  execution = {
    type = "object"
    required = false
    description = "Project execution configuration"
    
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
        description = "Maximum parallel executions"
      }
      
      dry_run_default = {
        type = "boolean"
        required = false
        default = false
        description = "Default dry-run mode for actions"
      }
      
      validate_before_execute = {
        type = "boolean"
        required = false
        default = true
        description = "Validate project configuration before execution"
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
        enum = ["badgerdb", "json"]
        default = "badgerdb"
        description = "Storage format for project facts database"
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
    
    # Execution validation
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
      message = "Maximum parallel executions must be between 1 and 100"
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
  }
} 