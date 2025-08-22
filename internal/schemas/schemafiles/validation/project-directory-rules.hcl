# Project Directory Validation Rules
# Validation rules for project-directory.hcl schema
# These rules validate schema compliance and data format correctness

# Project directory validation rules
validation_rules {
  # File or directory existence validation
  file_or_directory_validation {
    rule {
      name = "machines_file_or_directory_exists"
      description = "Either machines.hcl file or machines/ directory must exist"
      condition = "!fileExists('machines.hcl') && !directoryExists('machines')"
      message = "Either machines.hcl file or machines/ directory must exist"
      severity = "error"
    }
    
    rule {
      name = "actions_file_or_directory_exists"
      description = "Either actions.hcl file or actions/ directory must exist"
      condition = "!fileExists('actions.hcl') && !directoryExists('actions')"
      message = "Either actions.hcl file or actions/ directory must exist"
      severity = "error"
    }
    
    rule {
      name = "variables_file_or_directory_exists"
      description = "Either variables.hcl file or variables/ directory must exist"
      condition = "!fileExists('variables.hcl') && !directoryExists('variables')"
      message = "Either variables.hcl file or variables/ directory must exist"
      severity = "error"
    }
  }
  
  # Cross-file validation
  cross_file_validation {
    rule {
      name = "no_circular_references"
      description = "No circular references allowed between project files"
      condition = "hasCircularReferences(projectFiles)"
      message = "No circular references allowed between project files"
      severity = "error"
    }
    
    rule {
      name = "logging_file_output_requires_logs_directory"
      description = "Logging file output requires logs/ directory to exist"
      condition = "loggingOutput == 'file' && !directoryExists('logs')"
      message = "Logging file output requires logs/ directory to exist"
      severity = "error"
    }
    
    rule {
      name = "logging_file_path_validation"
      description = "Logging file path must be valid and writable"
      condition = "loggingOutput == 'file' && !isValidLoggingPath(loggingFilePath)"
      message = "Logging file path must be valid and writable"
      severity = "error"
    }
  }
  
  # File content validation
  file_content_validation {
    rule {
      name = "project_hcl_required_fields"
      description = "project.hcl must contain required fields"
      condition = "fileExists('project.hcl') && !hasRequiredProjectFields('project.hcl')"
      message = "project.hcl must contain required fields"
      severity = "error"
    }
    
    rule {
      name = "machines_hcl_format"
      description = "machines.hcl must contain valid machines block"
      condition = "fileExists('machines.hcl') && !hasValidMachinesBlock('machines.hcl')"
      message = "machines.hcl must contain valid machines block"
      severity = "error"
    }
    
    rule {
      name = "actions_hcl_format"
      description = "actions.hcl must contain valid actions block"
      condition = "fileExists('actions.hcl') && !hasValidActionsBlock('actions.hcl')"
      message = "actions.hcl must contain valid actions block"
      severity = "error"
    }
    
    rule {
      name = "variables_hcl_format"
      description = "variables.hcl must contain valid variables block"
      condition = "fileExists('variables.hcl') && !hasValidVariablesBlock('variables.hcl')"
      message = "variables.hcl must contain valid variables block"
      severity = "error"
    }
  }
  
  # Directory structure validation
  directory_structure_validation {
    rule {
      name = "templates_directory_optional"
      description = "templates/ directory is optional but must contain .tmpl files if present"
      condition = "directoryExists('templates') && !hasTemplateFiles('templates')"
      message = "templates/ directory must contain .tmpl files if present"
      severity = "warning"
    }
    
    rule {
      name = "files_directory_optional"
      description = "files/ directory is optional but must contain files if present"
      condition = "directoryExists('files') && !hasFiles('files')"
      message = "files/ directory must contain files if present"
      severity = "warning"
    }
    
    rule {
      name = "logs_directory_optional"
      description = "logs/ directory is optional but must be writable if present"
      condition = "directoryExists('logs') && !isWritableDirectory('logs')"
      message = "logs/ directory must be writable if present"
      severity = "error"
    }
  }
}
