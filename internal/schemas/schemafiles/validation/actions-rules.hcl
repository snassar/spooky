# Actions Validation Rules
# Extracted from internal/schemas/schemas/structure/actions.hcl
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
      name = "template_type_validation"
      description = "Template deploy actions must have template configuration"
      condition = "type == 'template_deploy' && template == null"
      message = "Template deploy actions must have template configuration"
      severity = "error"
    }
    
    rule {
      name = "file_copy_type_validation"
      description = "File copy actions must have file_copy configuration"
      condition = "type == 'file_copy' && file_copy == null"
      message = "File copy actions must have file_copy configuration"
      severity = "error"
    }
    
    rule {
      name = "service_control_type_validation"
      description = "Service control actions must have service_control configuration"
      condition = "type == 'service_control' && service_control == null"
      message = "Service control actions must have service_control configuration"
      severity = "error"
    }
    
    rule {
      name = "service_control_required_fields"
      description = "Service control actions must specify both service name and action"
      condition = "type == 'service_control' && (service_control.service == null || service_control.action == null)"
      message = "Service control actions must specify both service name and action"
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
      condition = "type == 'service_control' && service_control.action != null && !['start', 'stop', 'restart', 'reload', 'enable', 'disable', 'status'].contains(service_control.action)"
      message = "Service action must be one of: start, stop, restart, reload, enable, disable, status"
      severity = "error"
    }
    
    rule {
      name = "service_timeout_reasonable"
      description = "Service timeout must be between 1 and 300 seconds"
      condition = "type == 'service_control' && service_control.timeout != null && (service_control.timeout < 1 || service_control.timeout > 300)"
      message = "Service timeout must be between 1 and 300 seconds"
      severity = "warning"
    }
    
    rule {
      name = "wait_status_valid"
      description = "Wait status must be one of: active, inactive, failed, any"
      condition = "type == 'service_control' && service_control.wait_status != null && !['active', 'inactive', 'failed', 'any'].contains(service_control.wait_status)"
      message = "Wait status must be one of: active, inactive, failed, any"
      severity = "error"
    }
  }
}
