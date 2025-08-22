# Spooky Validation Rules
# Validation rules for spooky.hcl schema
# These rules validate schema compliance and data format correctness

# Spooky validation rules
validation_rules {
  # SSH validation
  ssh_validation {
    rule {
      name = "ssh_timeout_reasonable"
      description = "SSH timeout must be between 1 and 300 seconds"
      condition = "ssh.timeout != null && (ssh.timeout < 1 || ssh.timeout > 300)"
      message = "SSH timeout must be between 1 and 300 seconds"
      severity = "warning"
    }
    
    rule {
      name = "ssh_keepalive_interval_reasonable"
      description = "SSH keepalive interval must be between 1 and 300 seconds"
      condition = "ssh.keepalive_interval != null && (ssh.keepalive_interval < 1 || ssh.keepalive_interval > 300)"
      message = "SSH keepalive interval must be between 1 and 300 seconds"
      severity = "warning"
    }
    
    rule {
      name = "ssh_keepalive_count_reasonable"
      description = "SSH keepalive count must be between 1 and 10"
      condition = "ssh.keepalive_count != null && (ssh.keepalive_count < 1 || ssh.keepalive_count > 10)"
      message = "SSH keepalive count must be between 1 and 10"
      severity = "warning"
    }
    
    rule {
      name = "ssh_key_scan_timeout_reasonable"
      description = "SSH key scan timeout must be between 1 and 60 seconds"
      condition = "ssh.key_scan_timeout != null && (ssh.key_scan_timeout < 1 || ssh.key_scan_timeout > 60)"
      message = "SSH key scan timeout must be between 1 and 60 seconds"
      severity = "warning"
    }
    
    rule {
      name = "ssh_known_hosts_strict_boolean"
      description = "SSH known hosts strict must be a boolean value"
      condition = "ssh.known_hosts_strict != null && typeof(ssh.known_hosts_strict) != 'boolean'"
      message = "SSH known hosts strict must be a boolean value"
      severity = "error"
    }
    
    rule {
      name = "ssh_connection_pool_size_reasonable"
      description = "SSH connection pool size must be between 1 and 100"
      condition = "ssh.connection_pool_size != null && (ssh.connection_pool_size < 1 || ssh.connection_pool_size > 100)"
      message = "SSH connection pool size must be between 1 and 100"
      severity = "warning"
    }
  }
  
  # Security validation
  security_validation {
    rule {
      name = "allow_unsafe_commands_boolean"
      description = "Allow unsafe commands must be a boolean value"
      condition = "security.allow_unsafe_commands != null && typeof(security.allow_unsafe_commands) != 'boolean'"
      message = "Allow unsafe commands must be a boolean value"
      severity = "error"
    }
    
    rule {
      name = "audit_logging_boolean"
      description = "Audit logging must be a boolean value"
      condition = "security.audit_logging != null && typeof(security.audit_logging) != 'boolean'"
      message = "Audit logging must be a boolean value"
      severity = "error"
    }
  }
  
  # Age encryption validation
  age_validation {
    rule {
      name = "age_identities_path_valid"
      description = "Age identities path must be a valid directory path"
      condition = "age.identities != null && !isValidDirectoryPath(age.identities)"
      message = "Age identities path must be a valid directory path"
      severity = "error"
    }
    
    rule {
      name = "age_recipients_path_valid"
      description = "Age recipients path must be a valid file path"
      condition = "age.recipients != null && !isValidFilePath(age.recipients)"
      message = "Age recipients path must be a valid file path"
      severity = "error"
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
      name = "logging_color_output_boolean"
      description = "Color output must be a boolean value"
      condition = "logging.color_output != null && typeof(logging.color_output) != 'boolean'"
      message = "Color output must be a boolean value"
      severity = "error"
    }
    
    rule {
      name = "logging_progress_bars_boolean"
      description = "Progress bars must be a boolean value"
      condition = "logging.progress_bars != null && typeof(logging.progress_bars) != 'boolean'"
      message = "Progress bars must be a boolean value"
      severity = "error"
    }
    
    rule {
      name = "logging_file_path_required"
      description = "File path is required when output is set to 'file'"
      condition = "logging.output == 'file' && logging.file_path == null"
      message = "File path is required when output is set to 'file'"
      severity = "error"
    }
    
    rule {
      name = "logging_file_path_format"
      description = "Log file path must not contain invalid characters"
      condition = "logging.file_path != null && logging.file_path.matches('.*[<>:\"/\\\\|?*].*')"
      message = "Log file path must not contain invalid characters"
      severity = "error"
    }
    
    rule {
      name = "logging_file_permissions_valid"
      description = "File permissions must be in octal format (e.g., 0644)"
      condition = "logging.file_permissions != null && !logging.file_permissions.matches('^[0-7]{3,4}$')"
      message = "File permissions must be in octal format (e.g., 0644)"
      severity = "error"
    }
    
    rule {
      name = "logging_file_append_boolean"
      description = "File append must be a boolean value"
      condition = "logging.file_append != null && typeof(logging.file_append) != 'boolean'"
      message = "File append must be a boolean value"
      severity = "error"
    }
    
    rule {
      name = "logging_component_levels_valid"
      description = "Component log levels must be valid log levels"
      condition = "logging.filtering_components != null && !allComponentLevelsValid(logging.filtering_components)"
      message = "Component log levels must be valid log levels"
      severity = "error"
    }
  }
  
  # Performance validation
  performance_validation {
    rule {
      name = "max_memory_reasonable"
      description = "Maximum memory must be between 64MB and 8GB"
      condition = "performance.max_memory_mb != null && (performance.max_memory_mb < 64 || performance.max_memory_mb > 8192)"
      message = "Maximum memory must be between 64MB and 8GB"
      severity = "warning"
    }
    
    rule {
      name = "overall_timeout_reasonable"
      description = "Overall timeout must be between 1 and 7200 seconds"
      condition = "performance.overall_timeout != null && (performance.overall_timeout < 1 || performance.overall_timeout > 7200)"
      message = "Overall timeout must be between 1 and 7200 seconds"
      severity = "warning"
    }
  }
  
  # Cross-field validation
  cross_field_validation {
    rule {
      name = "ssh_keepalive_interval_greater_than_timeout"
      description = "SSH keepalive interval should be greater than timeout"
      condition = "ssh.keepalive_interval != null && ssh.timeout != null && ssh.keepalive_interval <= ssh.timeout"
      message = "SSH keepalive interval should be greater than timeout"
      severity = "warning"
    }
    
    rule {
      name = "overall_timeout_greater_than_ssh_timeout"
      description = "Overall timeout should be greater than SSH timeout"
      condition = "performance.overall_timeout != null && ssh.timeout != null && performance.overall_timeout <= ssh.timeout"
      message = "Overall timeout should be greater than SSH timeout"
      severity = "warning"
    }
  }
}
