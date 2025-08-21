# Logging Validation Rules
# Created for internal/schemas/schemas/structure/logging.hcl
# These rules validate schema compliance and data format correctness

# Logging validation rules
validation_rules {
  # Log level validation
  log_level_validation {
    rule {
      name = "log_level_valid"
      description = "Log level must be one of: debug, info, warn, error, fatal"
      condition = "level != null && !['debug', 'info', 'warn', 'error', 'fatal'].contains(level)"
      message = "Log level must be one of: debug, info, warn, error, fatal"
      severity = "error"
    }
  }
  
  # Output format validation
  output_format_validation {
    rule {
      name = "output_format_valid"
      description = "Output format must be one of: json, text, structured"
      condition = "format != null && !['json', 'text', 'structured'].contains(format)"
      message = "Output format must be one of: json, text, structured"
      severity = "error"
    }
  }
  
  # Output destination validation
  output_destination_validation {
    rule {
      name = "output_destination_valid"
      description = "Output destination must be one of: stdout, stderr, file, null"
      condition = "output != null && !['stdout', 'stderr', 'file', 'null'].contains(output)"
      message = "Output destination must be one of: stdout, stderr, file, null"
      severity = "error"
    }
  }
  
  # File output validation
  file_output_validation {
    rule {
      name = "file_path_required_when_output_file"
      description = "File path is required when output is set to 'file'"
      condition = "output == 'file' && (file == null || file.path == null)"
      message = "File path is required when output is set to 'file'"
      severity = "error"
    }
    
    rule {
      name = "file_path_format"
      description = "File path must not contain invalid characters"
      condition = "file != null && file.path != null && file.path.matches('.*[<>:\"/\\\\|?*].*')"
      message = "File path must not contain invalid characters"
      severity = "error"
    }
    
    rule {
      name = "file_permissions_format"
      description = "File permissions must be in octal format (e.g., 0644)"
      condition = "file != null && file.permissions != null && !file.permissions.matches('^[0-7]{3,4}$')"
      message = "File permissions must be in octal format (e.g., 0644)"
      severity = "error"
    }
  }
  
  # Rotation validation
  rotation_validation {
    rule {
      name = "max_backups_reasonable"
      description = "Maximum log backups must be between 1 and 100"
      condition = "rotation != null && rotation.max_backups != null && (rotation.max_backups < 1 || rotation.max_backups > 100)"
      message = "Maximum log backups must be between 1 and 100"
      severity = "warning"
    }
    
    rule {
      name = "max_size_format"
      description = "Maximum file size must be in valid format (e.g., 100MB, 1GB)"
      condition = "rotation != null && rotation.max_size != null && !rotation.max_size.matches('^\\d+[KMGT]?B$')"
      message = "Maximum file size must be in valid format (e.g., 100MB, 1GB)"
      severity = "error"
    }
    
    rule {
      name = "max_age_format"
      description = "Maximum age must be in valid format (e.g., 30d, 7d, 24h)"
      condition = "rotation != null && rotation.max_age != null && !rotation.max_age.matches('^\\d+[dhms]$')"
      message = "Maximum age must be in valid format (e.g., 30d, 7d, 24h)"
      severity = "error"
    }
  }
  
  # Component filtering validation
  component_filtering_validation {
    rule {
      name = "component_log_levels_valid"
      description = "Component log levels must be valid log levels"
      condition = "filtering != null && filtering.components != null && !allComponentLevelsValid(filtering.components)"
      message = "Component log levels must be valid log levels"
      severity = "error"
    }
  }
  
  # Pattern filtering validation
  pattern_filtering_validation {
    rule {
      name = "include_patterns_valid_regex"
      description = "Include patterns must be valid regex patterns"
      condition = "filtering != null && filtering.patterns != null && filtering.patterns.include != null && !allPatternsValid(filtering.patterns.include)"
      message = "Include patterns must be valid regex patterns"
      severity = "error"
    }
    
    rule {
      name = "exclude_patterns_valid_regex"
      description = "Exclude patterns must be valid regex patterns"
      condition = "filtering != null && filtering.patterns != null && filtering.patterns.exclude != null && !allPatternsValid(filtering.patterns.exclude)"
      message = "Exclude patterns must be valid regex patterns"
      severity = "error"
    }
  }
}
