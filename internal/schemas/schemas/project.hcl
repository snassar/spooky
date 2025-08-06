# Spooky Project Configuration Schema
# Comprehensive schema for project.hcl files with project metadata and structure validation

# Project block structure
project {
  # Project metadata
  name = {
    type = "string"
    required = true
    pattern = "^[a-zA-Z][a-zA-Z0-9._-]*$"
    min_length = 1
    max_length = 100
    description = "Project name (used for identification and isolation)"
  }
  
  description = {
    type = "string"
    required = false
    max_length = 500
    description = "Project description"
  }
  
  version = {
    type = "string"
    required = false
    pattern = "^[0-9]+\\.[0-9]+\\.[0-9]+(-[a-zA-Z0-9._-]+)?$"
    description = "Project version (semantic versioning)"
  }
  
  author = {
    type = "string"
    required = false
    max_length = 100
    description = "Project author or maintainer"
  }
  
  contact = {
    type = "object"
    required = false
    description = "Project contact information"
    
    properties = {
      email = {
        type = "string"
        required = false
        format = "email"
        description = "Contact email address"
      }
      
      url = {
        type = "string"
        required = false
        format = "uri"
        description = "Project website or documentation URL"
      }
    }
  }
  
  # Project configuration
  environment = {
    type = "string"
    required = false
    enum = ["production", "staging", "development", "testing"]
    default = "development"
    description = "Project environment"
  }
  
  region = {
    type = "string"
    required = false
    description = "Primary geographic region for this project"
  }
  
  tags = {
    type = "array"
    required = false
    description = "Project tags for organization and filtering"
    items = {
      type = "string"
      pattern = "^[a-zA-Z0-9._-]+$"
    }
  }
  
  # Project structure
  structure = {
    type = "object"
    required = false
    description = "Project directory structure configuration"
    
    properties = {
      templates_dir = {
        type = "string"
        required = false
        default = "templates"
        pattern = "^[a-zA-Z0-9/._-]+$"
        description = "Directory containing template files"
      }
      
      data_dir = {
        type = "string"
        required = false
        default = "data"
        pattern = "^[a-zA-Z0-9/._-]+$"
        description = "Directory containing data files"
      }
      
      scripts_dir = {
        type = "string"
        required = false
        default = "scripts"
        pattern = "^[a-zA-Z0-9/._-]+$"
        description = "Directory containing script files"
      }
      
      logs_dir = {
        type = "string"
        required = false
        default = "logs"
        pattern = "^[a-zA-Z0-9/._-]+$"
        description = "Directory for log files"
      }
      
      backups_dir = {
        type = "string"
        required = false
        default = "backups"
        pattern = "^[a-zA-Z0-9/._-]+$"
        description = "Directory for backup files"
      }
    }
  }
  
  # Project isolation and security
  isolation = {
    type = "object"
    required = false
    description = "Project isolation and security settings"
    
    properties = {
      enabled = {
        type = "boolean"
        required = false
        default = true
        description = "Enable project isolation"
      }
      
      facts_scope = {
        type = "string"
        required = false
        enum = ["global", "project", "hybrid"]
        default = "project"
        description = "Scope for facts collection and storage"
      }
      
      variables_scope = {
        type = "string"
        required = false
        enum = ["project", "inherited"]
        default = "project"
        description = "Scope for variables (project-only or inherited from global)"
      }
      
      machine_access = {
        type = "string"
        required = false
        enum = ["all", "tagged", "explicit"]
        default = "all"
        description = "Machine access control level"
      }
      
      allowed_machines = {
        type = "array"
        required = false
        description = "List of machine names explicitly allowed for this project"
        items = {
          type = "string"
          pattern = "^[a-zA-Z0-9._-]+$"
        }
      }
      
      allowed_tags = {
        type = "array"
        required = false
        description = "List of tags for machines allowed in this project"
        items = {
          type = "string"
          pattern = "^[a-zA-Z0-9._-]+$"
        }
      }
    }
  }
  
  # Project dependencies
  dependencies = {
    type = "object"
    required = false
    description = "Project dependencies and imports"
    
    properties = {
      imports = {
        type = "array"
        required = false
        description = "List of other projects to import"
        items = {
          type = "string"
          pattern = "^[a-zA-Z0-9/._-]+$"
        }
      }
      
      shared_variables = {
        type = "array"
        required = false
        description = "List of variable names to share with imported projects"
        items = {
          type = "string"
          pattern = "^[a-zA-Z0-9._-]+$"
        }
      }
      
      shared_facts = {
        type = "array"
        required = false
        description = "List of fact names to share with imported projects"
        items = {
          type = "string"
          pattern = "^[a-zA-Z0-9._-]+$"
        }
      }
    }
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
    
    # Version validation
    version_format = {
      rule = "regex"
      pattern = "^[0-9]+\\.[0-9]+\\.[0-9]+(-[a-zA-Z0-9._-]+)?$"
      message = "Version must follow semantic versioning format (e.g., 1.0.0 or 1.0.0-beta)"
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
    
    # Directory validation
    directory_names = {
      rule = "regex"
      pattern = "^[a-zA-Z0-9/._-]+$"
      message = "Directory names must contain only alphanumeric characters, dots, underscores, hyphens, and forward slashes"
    }
    
    # Isolation validation
    isolation_config = {
      rule = "conditional"
      condition = "isolation.enabled == false || (isolation.machine_access != 'explicit' || isolation.allowed_machines != null || isolation.allowed_tags != null)"
      message = "Explicit machine access requires allowed_machines or allowed_tags to be specified"
    }
    
    # Dependency validation
    no_circular_deps = {
      rule = "acyclic"
      message = "Project dependencies cannot have circular references"
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
  }
} 