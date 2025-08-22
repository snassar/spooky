# Logging Validation Rules
# Validation rules for logging.hcl schema
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
      condition = "output == 'file' && file_path == null"
      message = "File path is required when output is set to 'file'"
      severity = "error"
    }
    
    rule {
      name = "file_path_format"
      description = "File path must not contain invalid characters"
      condition = "file_path != null && file_path.matches('.*[<>:\"/\\\\|?*].*')"
      message = "File path must not contain invalid characters"
      severity = "error"
    }
    
    rule {
      name = "file_permissions_format"
      description = "File permissions must be in octal format (e.g., 0644)"
      condition = "file_permissions != null && !file_permissions.matches('^[0-7]{3,4}$')"
      message = "File permissions must be in octal format (e.g., 0644)"
      severity = "error"
    }
    
    rule {
      name = "file_append_boolean"
      description = "File append must be a boolean value"
      condition = "file_append != null && typeof(file_append) != 'boolean'"
      message = "File append must be a boolean value"
      severity = "error"
    }
  }
  
  # Structured logging validation
  structured_logging_validation {
    rule {
      name = "structured_timestamp_enabled_boolean"
      description = "Structured timestamp enabled must be a boolean value"
      condition = "structured_timestamp_enabled != null && typeof(structured_timestamp_enabled) != 'boolean'"
      message = "Structured timestamp enabled must be a boolean value"
      severity = "error"
    }
    
    rule {
      name = "structured_timestamp_format_valid"
      description = "Structured timestamp format must be one of: RFC3339, RFC3339Nano, Unix, UnixNano, ISO8601"
      condition = "structured_timestamp_format != null && !['RFC3339', 'RFC3339Nano', 'Unix', 'UnixNano', 'ISO8601'].contains(structured_timestamp_format)"
      message = "Structured timestamp format must be one of: RFC3339, RFC3339Nano, Unix, UnixNano, ISO8601"
      severity = "error"
    }
    
    rule {
      name = "structured_timestamp_timezone_format"
      description = "Structured timestamp timezone must be a valid timezone identifier"
      condition = "structured_timestamp_timezone != null && !structured_timestamp_timezone.matches('^[A-Za-z_]+/[A-Za-z_]+$')"
      message = "Structured timestamp timezone must be a valid timezone identifier (e.g., UTC, America/New_York)"
      severity = "error"
    }
    
    rule {
      name = "structured_level_key_format"
      description = "Structured level key must be a valid identifier"
      condition = "structured_level_key != null && !structured_level_key.matches('^[a-zA-Z_][a-zA-Z0-9_]*$')"
      message = "Structured level key must be a valid identifier"
      severity = "error"
    }
    
    rule {
      name = "structured_message_key_format"
      description = "Structured message key must be a valid identifier"
      condition = "structured_message_key != null && !structured_message_key.matches('^[a-zA-Z_][a-zA-Z0-9_]*$')"
      message = "Structured message key must be a valid identifier"
      severity = "error"
    }
    
    rule {
      name = "structured_error_key_format"
      description = "Structured error key must be a valid identifier"
      condition = "structured_error_key != null && !structured_error_key.matches('^[a-zA-Z_][a-zA-Z0-9_]*$')"
      message = "Structured error key must be a valid identifier"
      severity = "error"
    }
    
    rule {
      name = "structured_fields_include_array"
      description = "Structured fields include must be an array of strings"
      condition = "structured_fields_include != null && !isArrayOfStrings(structured_fields_include)"
      message = "Structured fields include must be an array of strings"
      severity = "error"
    }
    
    rule {
      name = "structured_fields_exclude_array"
      description = "Structured fields exclude must be an array of strings"
      condition = "structured_fields_exclude != null && !isArrayOfStrings(structured_fields_exclude)"
      message = "Structured fields exclude must be an array of strings"
      severity = "error"
    }
    
    rule {
      name = "structured_fields_filter_sensitive_array"
      description = "Structured fields filter sensitive must be an array of strings"
      condition = "structured_fields_filter_sensitive != null && !isArrayOfStrings(structured_fields_filter_sensitive)"
      message = "Structured fields filter sensitive must be an array of strings"
      severity = "error"
    }
  }
  
  # Performance configuration validation
  performance_validation {
    rule {
      name = "performance_buffer_enabled_boolean"
      description = "Performance buffer enabled must be a boolean value"
      condition = "performance_buffer_enabled != null && typeof(performance_buffer_enabled) != 'boolean'"
      message = "Performance buffer enabled must be a boolean value"
      severity = "error"
    }
    
    rule {
      name = "performance_buffer_size_range"
      description = "Performance buffer size must be between 1024 and 1048576 bytes"
      condition = "performance_buffer_size != null && (performance_buffer_size < 1024 || performance_buffer_size > 1048576)"
      message = "Performance buffer size must be between 1024 and 1048576 bytes"
      severity = "error"
    }
    
    rule {
      name = "performance_buffer_flush_interval_format"
      description = "Performance buffer flush interval must be in valid duration format"
      condition = "performance_buffer_flush_interval != null && !performance_buffer_flush_interval.matches('^\\d+[nsuµmh]s?$')"
      message = "Performance buffer flush interval must be in valid duration format (e.g., 1s, 100ms)"
      severity = "error"
    }
    
    rule {
      name = "performance_async_enabled_boolean"
      description = "Performance async enabled must be a boolean value"
      condition = "performance_async_enabled != null && typeof(performance_async_enabled) != 'boolean'"
      message = "Performance async enabled must be a boolean value"
      severity = "error"
    }
    
    rule {
      name = "performance_async_queue_size_range"
      description = "Performance async queue size must be between 100 and 100000"
      condition = "performance_async_queue_size != null && (performance_async_queue_size < 100 || performance_async_queue_size > 100000)"
      message = "Performance async queue size must be between 100 and 100000"
      severity = "error"
    }
    
    rule {
      name = "performance_async_workers_range"
      description = "Performance async workers must be between 1 and 10"
      condition = "performance_async_workers != null && (performance_async_workers < 1 || performance_async_workers > 10)"
      message = "Performance async workers must be between 1 and 10"
      severity = "error"
    }
    
    rule {
      name = "performance_async_drop_when_full_boolean"
      description = "Performance async drop when full must be a boolean value"
      condition = "performance_async_drop_when_full != null && typeof(performance_async_drop_when_full) != 'boolean'"
      message = "Performance async drop when full must be a boolean value"
      severity = "error"
    }
  }
  
  # Filtering configuration validation
  filtering_validation {
    rule {
      name = "filtering_components_valid"
      description = "Component log levels must be valid log levels"
      condition = "filtering_components != null && !allComponentLevelsValid(filtering_components)"
      message = "Component log levels must be valid log levels"
      severity = "error"
    }
    
    rule {
      name = "filtering_patterns_include_valid_regex"
      description = "Include patterns must be valid regex patterns"
      condition = "filtering_patterns_include != null && !allPatternsValid(filtering_patterns_include)"
      message = "Include patterns must be valid regex patterns"
      severity = "error"
    }
    
    rule {
      name = "filtering_patterns_exclude_valid_regex"
      description = "Exclude patterns must be valid regex patterns"
      condition = "filtering_patterns_exclude != null && !allPatternsValid(filtering_patterns_exclude)"
      message = "Exclude patterns must be valid regex patterns"
      severity = "error"
    }
  }
  
  # Rotation configuration validation
  rotation_validation {
    rule {
      name = "rotation_enabled_boolean"
      description = "Rotation enabled must be a boolean value"
      condition = "rotation_enabled != null && typeof(rotation_enabled) != 'boolean'"
      message = "Rotation enabled must be a boolean value"
      severity = "error"
    }
    
    rule {
      name = "rotation_max_size_format"
      description = "Maximum file size must be in valid format (e.g., 100MB, 1GB)"
      condition = "rotation_max_size != null && !rotation_max_size.matches('^\\d+[KMGT]?B$')"
      message = "Maximum file size must be in valid format (e.g., 100MB, 1GB)"
      severity = "error"
    }
    
    rule {
      name = "rotation_max_age_format"
      description = "Maximum age must be in valid format (e.g., 30d, 7d, 24h)"
      condition = "rotation_max_age != null && !rotation_max_age.matches('^\\d+[dhms]$')"
      message = "Maximum age must be in valid format (e.g., 30d, 7d, 24h)"
      severity = "error"
    }
    
    rule {
      name = "rotation_max_backups_reasonable"
      description = "Maximum log backups must be between 1 and 100"
      condition = "rotation_max_backups != null && (rotation_max_backups < 1 || rotation_max_backups > 100)"
      message = "Maximum log backups must be between 1 and 100"
      severity = "warning"
    }
    
    rule {
      name = "rotation_compress_boolean"
      description = "Rotation compress must be a boolean value"
      condition = "rotation_compress != null && typeof(rotation_compress) != 'boolean'"
      message = "Rotation compress must be a boolean value"
      severity = "error"
    }
    
    rule {
      name = "rotation_local_time_boolean"
      description = "Rotation local time must be a boolean value"
      condition = "rotation_local_time != null && typeof(rotation_local_time) != 'boolean'"
      message = "Rotation local time must be a boolean value"
      severity = "error"
    }
  }
  
  # Cross-field validation
  cross_field_validation {
    rule {
      name = "performance_buffer_requires_buffer_enabled"
      description = "Performance buffer size requires buffer enabled to be true"
      condition = "performance_buffer_size != null && performance_buffer_enabled == false"
      message = "Performance buffer size requires buffer enabled to be true"
      severity = "warning"
    }
    
    rule {
      name = "performance_async_requires_async_enabled"
      description = "Performance async settings require async enabled to be true"
      condition = "(performance_async_queue_size != null || performance_async_workers != null || performance_async_drop_when_full != null) && performance_async_enabled == false"
      message = "Performance async settings require async enabled to be true"
      severity = "warning"
    }
    
    rule {
      name = "rotation_requires_rotation_enabled"
      description = "Rotation settings require rotation enabled to be true"
      condition = "(rotation_max_size != null || rotation_max_age != null || rotation_max_backups != null || rotation_compress != null || rotation_local_time != null) && rotation_enabled == false"
      message = "Rotation settings require rotation enabled to be true"
      severity = "warning"
    }
  }
}
