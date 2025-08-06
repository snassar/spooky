# Spooky Actions Configuration Schema
# Comprehensive schema for actions.hcl files with validation rules and security checks

# Actions block structure
actions {
  # Action definitions
  action "action_name" {
    description = {
      type = "string"
      required = true
      min_length = 1
      max_length = 500
      description = "Action description"
    }
    
    type = {
      type = "string"
      required = true
      enum = ["command", "script", "template_deploy", "file_copy", "service_control"]
      description = "Action execution type"
    }
    
    command = {
      type = "string"
      required = false
      min_length = 1
      max_length = 1000
      description = "Command to execute (for command type)"
      validation = {
        pattern = "^(?!.*[;&|`$]).*$"
        message = "Command cannot contain shell operators or special characters"
      }
    }
    
    script = {
      type = "string"
      required = false
      pattern = "^[a-zA-Z0-9/._-]+$"
      description = "Script file path (for script type)"
    }
    
    template = {
      type = "object"
      required = false
      description = "Template configuration (for template_deploy type)"
      
      properties = {
        source = {
          type = "string"
          required = true
          pattern = "^templates/[a-zA-Z0-9/._-]+\\.tmpl$"
          description = "Source template file path"
        }
        
        destination = {
          type = "string"
          required = true
          pattern = "^[a-zA-Z0-9/._-]+$"
          description = "Destination file path on target machine"
        }
        
        validate = {
          type = "boolean"
          required = false
          default = false
          description = "Validate template syntax before deployment"
        }
        
        backup = {
          type = "boolean"
          required = false
          default = false
          description = "Create backup of existing file before overwriting"
        }
        
        permissions = {
          type = "string"
          required = false
          pattern = "^[0-7]{3,4}$"
          description = "File permissions (octal format)"
        }
        
        owner = {
          type = "string"
          required = false
          pattern = "^[a-zA-Z0-9._-]+$"
          description = "File owner (username)"
        }
        
        group = {
          type = "string"
          required = false
          pattern = "^[a-zA-Z0-9._-]+$"
          description = "File group (group name)"
        }
      }
    }
    
    machines = {
      type = "array"
      required = false
      description = "List of target machine names"
      items = {
        type = "string"
        pattern = "^[a-zA-Z0-9._-]+$"
      }
    }
    
    tags = {
      type = "array"
      required = false
      description = "List of tags for targeting machines"
      items = {
        type = "string"
        pattern = "^[a-zA-Z0-9._-]+$"
      }
    }
    
    timeout = {
      type = "integer"
      required = false
      min = 1
      max = 3600
      default = 300
      description = "Action timeout in seconds"
    }
    
    parallel = {
      type = "boolean"
      required = false
      default = false
      description = "Execute action in parallel across machines"
    }
    
    retries = {
      type = "integer"
      required = false
      min = 0
      max = 10
      default = 0
      description = "Number of retry attempts on failure"
    }
    
    retry_delay = {
      type = "integer"
      required = false
      min = 1
      max = 300
      default = 5
      description = "Delay between retries in seconds"
    }
    
    dependencies = {
      type = "array"
      required = false
      description = "List of action names that must complete before this action"
      items = {
        type = "string"
        pattern = "^[a-zA-Z0-9._-]+$"
      }
    }
    
    environment = {
      type = "object"
      required = false
      description = "Environment variables for the action"
      additional_properties = "string"
    }
    
    working_directory = {
      type = "string"
      required = false
      pattern = "^[a-zA-Z0-9/._-]+$"
      description = "Working directory for command execution"
    }
    
    user = {
      type = "string"
      required = false
      pattern = "^[a-zA-Z0-9._-]+$"
      description = "User to run the action as (overrides machine user)"
    }
    
    sudo = {
      type = "boolean"
      required = false
      default = false
      description = "Run action with sudo privileges"
    }
    
    dry_run = {
      type = "boolean"
      required = false
      default = false
      description = "Show what would be executed without actually running"
    }
    
    # Extended properties for action metadata and organization
    category = {
      type = "string"
      required = false
      pattern = "^[a-zA-Z0-9._-]+$"
      description = "Action category for organization"
    }
    
    priority = {
      type = "integer"
      required = false
      min = 1
      max = 10
      default = 5
      description = "Action priority (1=lowest, 10=highest)"
    }
    
    critical = {
      type = "boolean"
      required = false
      default = false
      description = "Whether action failure should stop execution"
    }
    
    metadata = {
      type = "object"
      required = false
      description = "Additional metadata for the action"
      additional_properties = "string"
    }
    
    # Security and validation properties
    validate_before_run = {
      type = "boolean"
      required = false
      default = true
      description = "Validate action before execution"
    }
    
    allow_failure = {
      type = "boolean"
      required = false
      default = false
      description = "Allow action to fail without stopping execution"
    }
    
    # Performance and resource properties
    max_concurrent = {
      type = "integer"
      required = false
      min = 1
      max = 100
      default = 1
      description = "Maximum concurrent executions of this action"
    }
    
    resource_limits = {
      type = "object"
      required = false
      description = "Resource limits for action execution"
      
      properties = {
        memory_mb = {
          type = "integer"
          required = false
          min = 1
          max = 32768
          description = "Memory limit in MB"
        }
        
        cpu_percent = {
          type = "integer"
          required = false
          min = 1
          max = 100
          description = "CPU usage limit as percentage"
        }
        
        disk_mb = {
          type = "integer"
          required = false
          min = 1
          max = 1048576
          description = "Disk usage limit in MB"
        }
      }
    }
  }
  
  # Validation rules
  validation = {
    # Action name validation
    action_name = {
      rule = "regex"
      pattern = "^[a-zA-Z][a-zA-Z0-9._-]*$"
      message = "Action names must start with a letter and contain only alphanumeric characters, dots, underscores, and hyphens"
    }
    
    # Type-specific validation
    command_type = {
      rule = "conditional"
      condition = "type == 'command' && command != null"
      message = "Command type actions must have a command specified"
    }
    
    script_type = {
      rule = "conditional"
      condition = "type == 'script' && script != null"
      message = "Script type actions must have a script path specified"
    }
    
    template_type = {
      rule = "conditional"
      condition = "type == 'template_deploy' && template != null"
      message = "Template deploy actions must have template configuration"
    }
    
    # Security validation
    command_security = {
      rule = "regex"
      pattern = "^(?!.*[;&|`$]).*$"
      message = "Commands cannot contain shell operators or special characters for security"
    }
    
    # Dependency validation
    no_circular_deps = {
      rule = "acyclic"
      message = "Actions cannot have circular dependencies"
    }
    
    # Machine targeting validation
    machine_targeting = {
      rule = "conditional"
      condition = "machines != null || tags != null"
      message = "Actions must target machines by name or tags"
    }
    
    # Extended validation rules
    dependency_validation = {
      rule = "dependency_check"
      message = "Action dependencies must reference valid action names"
    }
    
    resource_validation = {
      rule = "resource_check"
      message = "Resource limits must be reasonable for the target environment"
    }
    
    priority_validation = {
      rule = "priority_check"
      message = "Action priority must be between 1 and 10"
    }
  }
} 