# Project Validation Rules
# Extracted from internal/schemas/schemas/structure/project.hcl
# These rules validate schema compliance and data format correctness

# Project validation rules
validation_rules {
  # Project name validation
  project_name_validation {
    rule {
      name = "project_name_format"
      description = "Project names must start with a letter and contain only alphanumeric characters, dots, underscores, and hyphens"
      condition = "name != null && !name.matches('^[a-zA-Z][a-zA-Z0-9._-]*$')"
      message = "Project names must start with a letter and contain only alphanumeric characters, dots, underscores, and hyphens"
      severity = "error"
    }
  }
  
  # Contact validation
  contact_validation {
    rule {
      name = "email_format"
      description = "Contact email must be a valid email address"
      condition = "contact.email != null && !isValidEmail(contact.email)"
      message = "Contact email must be a valid email address"
      severity = "error"
    }
    
    rule {
      name = "url_format"
      description = "Contact URL must be a valid URI"
      condition = "contact.url != null && !isValidURI(contact.url)"
      message = "Contact URL must be a valid URI"
      severity = "error"
    }
  }
  
  # Run validation
  run_validation {
    rule {
      name = "timeout_reasonable"
      description = "Timeouts must be between 1 and 3600 seconds"
      condition = "run.timeout != null && (run.timeout < 1 || run.timeout > 3600)"
      message = "Timeouts must be between 1 and 3600 seconds"
      severity = "warning"
    }
    
    rule {
      name = "parallel_limit"
      description = "Maximum parallel runs must be between 1 and 100"
      condition = "run.parallel != null && (run.parallel < 1 || run.parallel > 100)"
      message = "Maximum parallel runs must be between 1 and 100"
      severity = "warning"
    }
  }
  
  # Facts validation
  facts_validation {
    rule {
      name = "facts_timeout_reasonable"
      description = "Facts collection timeouts must be between 1 and 3600 seconds"
      condition = "facts.timeout != null && (facts.timeout < 1 || facts.timeout > 3600)"
      message = "Facts collection timeouts must be between 1 and 3600 seconds"
      severity = "warning"
    }
    
    rule {
      name = "facts_parallel_reasonable"
      description = "Facts parallel collection must be between 1 and 100 workers"
      condition = "facts.parallel != null && (facts.parallel < 1 || facts.parallel > 100)"
      message = "Facts parallel collection must be between 1 and 100 workers"
      severity = "warning"
    }
  }
  
  # Logging validation
  logging_validation {
    rule {
      name = "logging_level_valid"
      description = "Log level must be one of: debug, info, warn, error, fatal"
      condition = "logging.level != null && !['debug', 'info', 'warn', 'error', 'fatal'].contains(logging.level)"
      message = "Log level must be one of: debug, info, warn, error, fatal"
      severity = "error"
    }
    
    rule {
      name = "logging_format_valid"
      description = "Log format must be one of: json, text, structured"
      condition = "logging.format != null && !['json', 'text', 'structured'].contains(logging.format)"
      message = "Log format must be one of: json, text, structured"
      severity = "error"
    }
    
    rule {
      name = "logging_output_valid"
      description = "Log output must be one of: stdout, stderr, file, null"
      condition = "logging.output != null && !['stdout', 'stderr', 'file', 'null'].contains(logging.output)"
      message = "Log output must be one of: stdout, stderr, file, null"
      severity = "error"
    }
    
    rule {
      name = "logging_file_permissions_valid"
      description = "File permissions must be in octal format (e.g., 0644)"
      condition = "logging.file.permissions != null && !logging.file.permissions.matches('^[0-7]{3,4}$')"
      message = "File permissions must be in octal format (e.g., 0644)"
      severity = "error"
    }
    
    rule {
      name = "logging_rotation_backups_reasonable"
      description = "Maximum log backups must be between 1 and 100"
      condition = "logging.rotation.max_backups != null && (logging.rotation.max_backups < 1 || logging.rotation.max_backups > 100)"
      message = "Maximum log backups must be between 1 and 100"
      severity = "warning"
    }
    
    rule {
      name = "logging_file_path_required"
      description = "File path is required when output is set to 'file'"
      condition = "logging.output == 'file' && logging.file.path == null"
      message = "File path is required when output is set to 'file'"
      severity = "error"
    }
  }
}
