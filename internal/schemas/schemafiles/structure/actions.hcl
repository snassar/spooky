# Spooky Actions Configuration Schema
# Comprehensive schema for actions.hcl files with validation rules and security checks

# Schema metadata
metadata {
  version = "1"
  description = "Actions configuration schema for spooky automation tasks"
}

# Actions block structure
actions {
  # Action definitions
  action {
    name = "action_name"
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
      enum = ["command", "script", "template_deploy", "file_sync", "service_control"]
      description = "Action run type"
    }
    
    command = {
      type = "string"
      required = false
      min_length = 1
      max_length = 1000
      description = "Command to run (for command type)"
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
    
    source = {
      type = "string"
      required = false
      pattern = "^templates/[a-zA-Z0-9/._-]+\\.tmpl$"
      description = "Source template file path (for template_deploy type)"
    }
    
    destination = {
      type = "string"
      required = false
      pattern = "^[a-zA-Z0-9/._-]+$"
      description = "Destination file path (local or remote) (for template_deploy type)"
    }
    
    validate = {
      type = "boolean"
      required = false
      default = false
      description = "Validate template syntax before deployment (for template_deploy type)"
    }
    
    backup = {
      type = "boolean"
      required = false
      default = false
      description = "Create backup of existing file before overwriting (for template_deploy type)"
    }
    
    permissions = {
      type = "string"
      required = false
      pattern = "^[0-7]{3,4}$"
      description = "File permissions in octal format (for template_deploy type)"
    }
    
    owner = {
      type = "string"
      required = false
      pattern = "^[a-zA-Z0-9._-]+$"
      description = "File owner username (for template_deploy type)"
    }
    
    group = {
      type = "string"
      required = false
      pattern = "^[a-zA-Z0-9._-]+$"
      description = "File group name (for template_deploy type)"
    }
    
    mode = {
      type = "string"
      required = false
      enum = ["one-way-replica", "one-way-safe", "two-way-safe", "two-way-resolved"]
      default = "one-way-replica"
      description = "Synchronization mode and directionality (for file_sync type)"
    }
    
    block_size = {
      type = "integer"
      required = false
      min = 512
      max = 65536
      default = 2048
      description = "Block size for rsync algorithm in bytes (for file_sync type)"
    }
    
    preserve_attributes = {
      type = "boolean"
      required = false
      default = true
      description = "Preserve file attributes (permissions, owner, group) (for file_sync type)"
    }
    
    service = {
      type = "string"
      required = false
      pattern = "^[a-zA-Z0-9._-]+$"
      description = "Name of the service to control (for service_control type)"
    }
    
    action = {
      type = "string"
      required = false
      enum = ["start", "stop", "restart", "reload", "enable", "disable", "status"]
      description = "Action to perform on the service (for service_control type)"
    }
    
    systemd = {
      type = "boolean"
      required = false
      default = true
      description = "Use systemd for service control (for service_control type)"
    }
    
    timeout = {
      type = "integer"
      required = false
      min = 1
      max = 3600
      default = 300
      description = "Action timeout in seconds (overrides project default)"
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
    
    critical = {
      type = "boolean"
      required = false
      default = false
      description = "Whether action failure should stop execution"
    }
    
    allow_failure = {
      type = "boolean"
      required = false
      default = false
      description = "Allow action to fail without stopping"
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
      description = "Working directory for command running"
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
    
    wait_for_status = {
      type = "string"
      required = false
      enum = ["active", "inactive", "failed", "any"]
      description = "Wait for specific service status after operation (for service_control type)"
    }
    
    wait_timeout = {
      type = "integer"
      required = false
      min = 1
      max = 300
      default = 60
      description = "Timeout for waiting for service status (for service_control type)"
    }
  }
} 