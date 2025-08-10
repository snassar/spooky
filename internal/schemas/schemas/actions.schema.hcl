# Spooky Actions Configuration Schema
# Comprehensive schema for actions.hcl files with validation rules and security checks

# Schema metadata
metadata {
  schema_version = "0.20250809.0"
  schema_type = "actions"
  schema_name = "Spooky Actions Configuration Schema"
  last_updated = "2024-01-01"
  compatibility = ["0.20250809.0"]
  description = "Comprehensive schema for actions.hcl files with validation rules and security checks"
  
  # ScalVer format: 0.YYYYMMDD.N
  # - 0: Development phase
  # - 20250809: Date (9 August 2025)
  # - 0: Patch version
  scalver_format = "0.20250809.0"
}

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
      enum = ["command", "script", "template_deploy", "template_evaluate", "template_validate", "template_cleanup", "file_copy", "service_control"]
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
      pattern = "^(files|templates)/[a-zA-Z0-9/._-]+(\\.sh|\\.tmpl)?$"
      description = "Script file path in files/ or templates/ directory (for script type)"
      validation = {
        pattern = "^(files|templates)/[a-zA-Z0-9/._-]+(\\.sh|\\.tmpl)?$"
        message = "Script must reference a file in files/ or templates/ directory"
      }
    }
    
    variables = {
      type = "object"
      required = false
      description = "Variables for templated scripts (for script type with .tmpl files)"
      additional_properties = "string"
      validation = {
        conditional = "script ends with .tmpl"
        message = "Variables are only valid for templated scripts (.tmpl files)"
      }
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
    
    file_copy = {
      type = "object"
      required = false
      description = "File copy configuration (for file_copy type)"
      
      properties = {
        source = {
          type = "string"
          required = true
          pattern = "^files/[a-zA-Z0-9/._-]+$"
          description = "Source file path in files/ directory"
        }
        
        destination = {
          type = "string"
          required = true
          pattern = "^[a-zA-Z0-9/._-]+$"
          description = "Destination file path on target machine"
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
    
    service_control = {
      type = "object"
      required = false
      description = "Service control configuration (for service_control type)"
      
      properties = {
        service = {
          type = "string"
          required = true
          pattern = "^[a-zA-Z0-9._-]+$"
          description = "Name of the service to control"
        }
        
        action = {
          type = "string"
          required = true
          enum = ["start", "stop", "restart", "reload", "enable", "disable", "status"]
          description = "Action to perform on the service"
        }
        
        systemd = {
          type = "boolean"
          required = false
          default = true
          description = "Use systemd for service control (default: true)"
        }
        
        timeout = {
          type = "integer"
          required = false
          min = 1
          max = 300
          default = 30
          description = "Service operation timeout in seconds"
        }
        
        wait_for_status = {
          type = "string"
          required = false
          enum = ["active", "inactive", "failed", "any"]
          description = "Wait for specific service status after operation"
        }
        
        wait_timeout = {
          type = "integer"
          required = false
          min = 1
          max = 300
          default = 60
          description = "Timeout for waiting for service status"
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
    
    script_file_path = {
      rule = "regex"
      pattern = "^(files|templates)/[a-zA-Z0-9/._-]+(\\.sh|\\.tmpl)?$"
      message = "Script must reference a file in files/ or templates/ directory"
    }
    
    script_file_exists = {
      rule = "file_exists"
      message = "Referenced script file does not exist"
    }
    
    template_variables = {
      rule = "conditional"
      condition = "script ends with .tmpl && variables != null"
      message = "Templated scripts (.tmpl) should have variables defined"
    }
    
    static_script = {
      rule = "conditional"
      condition = "script starts with 'files/' && !script ends with .tmpl"
      message = "Static scripts should be in files/ directory without .tmpl extension"
    }
    
    template_type = {
      rule = "conditional"
      condition = "type == 'template_deploy' && template != null"
      message = "Template deploy actions must have template configuration"
    }
    
    file_copy_type = {
      rule = "conditional"
      condition = "type == 'file_copy' && file_copy != null"
      message = "File copy actions must have file_copy configuration"
    }
    
    file_copy_source_exists = {
      rule = "file_exists"
      message = "Referenced file copy source file does not exist"
    }
    
    service_control_type = {
      rule = "conditional"
      condition = "type == 'service_control' && service_control != null"
      message = "Service control actions must have service_control configuration"
    }
    
    service_control_required = {
      rule = "conditional"
      condition = "type == 'service_control' && service_control.service != null && service_control.action != null"
      message = "Service control actions must specify both service name and action"
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
    
    # Service control specific validation
    service_action_valid = {
      rule = "enum"
      enum = ["start", "stop", "restart", "reload", "enable", "disable", "status"]
      message = "Service action must be one of: start, stop, restart, reload, enable, disable, status"
    }
    
    service_timeout_reasonable = {
      rule = "range"
      min = 1
      max = 300
      message = "Service timeout must be between 1 and 300 seconds"
    }
    
    wait_status_valid = {
      rule = "enum"
      enum = ["active", "inactive", "failed", "any"]
      message = "Wait status must be one of: active, inactive, failed, any"
    }
  }
} 