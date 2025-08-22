# Actions Validation Rules
# Validation rules for actions.hcl schema
# These rules validate schema compliance and data format correctness

# Actions validation rules
validation_rules {
  # Action name validation
  action_name_validation {
    rule {
      name = "action_name_format"
      description = "Action names must start with a letter and contain only alphanumeric characters, dots, underscores, and hyphens"
      condition = "name != null && !name.matches('^[a-zA-Z][a-zA-Z0-9._-]*$')"
      message = "Action names must start with a letter and contain only alphanumeric characters, dots, underscores, and hyphens"
      severity = "error"
    }
  }
  
  # Type-specific validation
  type_specific_validation {
    rule {
      name = "command_type_validation"
      description = "Command type actions must have a command specified"
      condition = "type == 'command' && command == null"
      message = "Command type actions must have a command specified"
      severity = "error"
    }
    
    rule {
      name = "script_type_validation"
      description = "Script type actions must have a script path specified"
      condition = "type == 'script' && script == null"
      message = "Script type actions must have a script path specified"
      severity = "error"
    }
    
    rule {
      name = "script_file_path_format"
      description = "Script must reference a file in files/ or templates/ directory"
      condition = "script != null && !script.matches('^(files|templates)/[a-zA-Z0-9/._-]+(\\.sh|\\.tmpl)?$')"
      message = "Script must reference a file in files/ or templates/ directory"
      severity = "error"
    }
    
    rule {
      name = "template_variables_required"
      description = "Templated scripts (.tmpl) should have variables defined"
      condition = "script != null && script.endsWith('.tmpl') && variables == null"
      message = "Templated scripts (.tmpl) should have variables defined"
      severity = "warning"
    }
    
    rule {
      name = "static_script_format"
      description = "Static scripts should be in files/ directory without .tmpl extension"
      condition = "script != null && script.startsWith('files/') && script.endsWith('.tmpl')"
      message = "Static scripts should be in files/ directory without .tmpl extension"
      severity = "warning"
    }
    
    rule {
      name = "template_deploy_validation"
      description = "Template deploy actions must have source and destination specified"
      condition = "type == 'template_deploy' && (source == null || destination == null)"
      message = "Template deploy actions must have source and destination specified"
      severity = "error"
    }
    
    rule {
      name = "template_deploy_source_format"
      description = "Template deploy source must be in templates/ directory with .tmpl extension"
      condition = "type == 'template_deploy' && source != null && !source.matches('^templates/[a-zA-Z0-9/._-]+\\.tmpl$')"
      message = "Template deploy source must be in templates/ directory with .tmpl extension"
      severity = "error"
    }
    
    rule {
      name = "template_deploy_destination_format"
      description = "Template deploy destination must be a valid file path"
      condition = "type == 'template_deploy' && destination != null && !destination.matches('^[a-zA-Z0-9/._-]+$')"
      message = "Template deploy destination must be a valid file path"
      severity = "error"
    }
    
    rule {
      name = "template_deploy_validate_boolean"
      description = "Template deploy validate must be a boolean value"
      condition = "type == 'template_deploy' && validate != null && typeof(validate) != 'boolean'"
      message = "Template deploy validate must be a boolean value"
      severity = "error"
    }
    
    rule {
      name = "template_deploy_backup_boolean"
      description = "Template deploy backup must be a boolean value"
      condition = "type == 'template_deploy' && backup != null && typeof(backup) != 'boolean'"
      message = "Template deploy backup must be a boolean value"
      severity = "error"
    }
    
    rule {
      name = "template_deploy_permissions_format"
      description = "Template deploy permissions must be in octal format"
      condition = "type == 'template_deploy' && permissions != null && !permissions.matches('^[0-7]{3,4}$')"
      message = "Template deploy permissions must be in octal format (e.g., 0644)"
      severity = "error"
    }
    
    rule {
      name = "template_deploy_owner_format"
      description = "Template deploy owner must be a valid username"
      condition = "type == 'template_deploy' && owner != null && !owner.matches('^[a-zA-Z0-9._-]+$')"
      message = "Template deploy owner must be a valid username"
      severity = "error"
    }
    
    rule {
      name = "template_deploy_group_format"
      description = "Template deploy group must be a valid group name"
      condition = "type == 'template_deploy' && group != null && !group.matches('^[a-zA-Z0-9._-]+$')"
      message = "Template deploy group must be a valid group name"
      severity = "error"
    }
    
    rule {
      name = "file_sync_validation"
      description = "File sync actions must have src and dest specified"
      condition = "type == 'file_sync' && (src == null || dest == null)"
      message = "File sync actions must have src and dest specified"
      severity = "error"
    }
    
    rule {
      name = "file_sync_source_exists"
      description = "Referenced file sync source file does not exist"
      condition = "type == 'file_sync' && src != null && !fileExists(src)"
      message = "Referenced file sync source file does not exist"
      severity = "error"
    }
    
    rule {
      name = "file_sync_mode_valid"
      description = "File sync mode must be one of the valid modes"
      condition = "type == 'file_sync' && mode != null && !['one-way-replica', 'one-way-safe', 'two-way-safe', 'two-way-resolved'].contains(mode)"
      message = "File sync mode must be one of: one-way-replica, one-way-safe, two-way-safe, two-way-resolved"
      severity = "error"
    }
    
    rule {
      name = "file_sync_block_size_valid"
      description = "File sync block size must be between 512 and 65536 bytes"
      condition = "type == 'file_sync' && block_size != null && (block_size < 512 || block_size > 65536)"
      message = "File sync block size must be between 512 and 65536 bytes"
      severity = "warning"
    }
    
    rule {
      name = "file_sync_preserve_attributes_boolean"
      description = "File sync preserve attributes must be a boolean value"
      condition = "type == 'file_sync' && preserve_attributes != null && typeof(preserve_attributes) != 'boolean'"
      message = "File sync preserve attributes must be a boolean value"
      severity = "error"
    }
    
    rule {
      name = "service_control_validation"
      description = "Service control actions must have service and action specified"
      condition = "type == 'service_control' && (service == null || action == null)"
      message = "Service control actions must have service and action specified"
      severity = "error"
    }
    
    rule {
      name = "service_control_service_format"
      description = "Service control service name must be a valid service name"
      condition = "type == 'service_control' && service != null && !service.matches('^[a-zA-Z0-9._-]+$')"
      message = "Service control service name must be a valid service name"
      severity = "error"
    }
    
    rule {
      name = "service_control_systemd_boolean"
      description = "Service control systemd must be a boolean value"
      condition = "type == 'service_control' && systemd != null && typeof(systemd) != 'boolean'"
      message = "Service control systemd must be a boolean value"
      severity = "error"
    }
  }
  
  # Security validation
  security_validation {
    rule {
      name = "command_security"
      description = "Commands cannot contain shell operators or special characters for security"
      condition = "command != null && command.matches('.*[;&|`$].*')"
      message = "Commands cannot contain shell operators or special characters for security"
      severity = "error"
    }
  }
  
  # Dependency validation
  dependency_validation {
    rule {
      name = "no_circular_dependencies"
      description = "Actions cannot have circular dependencies"
      condition = "dependencies != null && hasCircularDependencies(dependencies)"
      message = "Actions cannot have circular dependencies"
      severity = "error"
    }
    
    rule {
      name = "dependency_validation"
      description = "Action dependencies must reference valid action names"
      condition = "dependencies != null && !allDependenciesExist(dependencies)"
      message = "Action dependencies must reference valid action names"
      severity = "error"
    }
  }
  
  # Machine targeting validation
  machine_targeting_validation {
    rule {
      name = "machine_targeting_required"
      description = "Actions must target machines by name or tags"
      condition = "(machines == null || machines.isEmpty()) && (tags == null || tags.isEmpty())"
      message = "Actions must target machines by name or tags"
      severity = "error"
    }
  }
  
  # Service control specific validation
  service_control_validation {
    rule {
      name = "service_action_valid"
      description = "Service action must be one of: start, stop, restart, reload, enable, disable, status"
      condition = "type == 'service_control' && action != null && !['start', 'stop', 'restart', 'reload', 'enable', 'disable', 'status'].contains(action)"
      message = "Service action must be one of: start, stop, restart, reload, enable, disable, status"
      severity = "error"
    }
    
    rule {
      name = "service_timeout_reasonable"
      description = "Service timeout must be between 1 and 3600 seconds"
      condition = "type == 'service_control' && timeout != null && (timeout < 1 || timeout > 3600)"
      message = "Service timeout must be between 1 and 3600 seconds"
      severity = "warning"
    }
    
    rule {
      name = "wait_status_valid"
      description = "Wait status must be one of: active, inactive, failed, any"
      condition = "type == 'service_control' && wait_for_status != null && !['active', 'inactive', 'failed', 'any'].contains(wait_for_status)"
      message = "Wait status must be one of: active, inactive, failed, any"
      severity = "error"
    }
    
    rule {
      name = "wait_timeout_reasonable"
      description = "Wait timeout must be between 1 and 300 seconds"
      condition = "type == 'service_control' && wait_timeout != null && (wait_timeout < 1 || wait_timeout > 300)"
      message = "Wait timeout must be between 1 and 300 seconds"
      severity = "warning"
    }
  }
  
  # Execution configuration validation
  execution_validation {
    rule {
      name = "timeout_reasonable"
      description = "Action timeout must be between 1 and 3600 seconds"
      condition = "timeout != null && (timeout < 1 || timeout > 3600)"
      message = "Action timeout must be between 1 and 3600 seconds"
      severity = "warning"
    }
    
    rule {
      name = "retries_reasonable"
      description = "Retry attempts must be between 0 and 10"
      condition = "retries != null && (retries < 0 || retries > 10)"
      message = "Retry attempts must be between 0 and 10"
      severity = "warning"
    }
    
    rule {
      name = "retry_delay_reasonable"
      description = "Retry delay must be between 1 and 300 seconds"
      condition = "retry_delay != null && (retry_delay < 1 || retry_delay > 300)"
      message = "Retry delay must be between 1 and 300 seconds"
      severity = "warning"
    }
    
    rule {
      name = "critical_boolean"
      description = "Critical flag must be a boolean value"
      condition = "critical != null && typeof(critical) != 'boolean'"
      message = "Critical flag must be a boolean value"
      severity = "error"
    }
    
    rule {
      name = "allow_failure_boolean"
      description = "Allow failure flag must be a boolean value"
      condition = "allow_failure != null && typeof(allow_failure) != 'boolean'"
      message = "Allow failure flag must be a boolean value"
      severity = "error"
    }
    
    rule {
      name = "environment_variables_format"
      description = "Environment variables must be a valid object with string values"
      condition = "environment != null && !isValidEnvironmentObject(environment)"
      message = "Environment variables must be a valid object with string values"
      severity = "error"
    }
    
    rule {
      name = "working_directory_format"
      description = "Working directory must be a valid path"
      condition = "working_directory != null && !working_directory.matches('^[a-zA-Z0-9/._-]+$')"
      message = "Working directory must be a valid path"
      severity = "error"
    }
    
    rule {
      name = "user_format"
      description = "User must be a valid username"
      condition = "user != null && !user.matches('^[a-zA-Z0-9._-]+$')"
      message = "User must be a valid username"
      severity = "error"
    }
    
    rule {
      name = "sudo_boolean"
      description = "Sudo flag must be a boolean value"
      condition = "sudo != null && typeof(sudo) != 'boolean'"
      message = "Sudo flag must be a boolean value"
      severity = "error"
    }
  }
  
  # Tags validation
  tags_validation {
    rule {
      name = "tags_format"
      description = "Tags must contain only alphanumeric characters, underscores, and hyphens"
      condition = "tags != null && tags.any(tag -> !tag.matches('^[a-zA-Z0-9._-]+$'))"
      message = "Tags must contain only alphanumeric characters, underscores, and hyphens"
      severity = "error"
    }
    
    rule {
      name = "tags_count_limit"
      description = "Tags cannot exceed 20 items"
      condition = "tags != null && tags.size() > 20"
      message = "Tags cannot exceed 20 items"
      severity = "warning"
    }
  }
  
  # Cross-field validation
  cross_field_validation {
    rule {
      name = "critical_and_allow_failure_conflict"
      description = "Critical and allow failure flags cannot both be true"
      condition = "critical == true && allow_failure == true"
      message = "Critical and allow failure flags cannot both be true"
      severity = "error"
    }
    
    rule {
      name = "user_and_sudo_conflict"
      description = "User override and sudo flags should not be used together"
      condition = "user != null && sudo == true"
      message = "User override and sudo flags should not be used together"
      severity = "warning"
    }
  }
}
