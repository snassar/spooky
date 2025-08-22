# Templates Validation Rules
# Validation rules for templates.hcl schema
# These rules validate schema compliance and data format correctness

# Templates validation rules
validation_rules {
  # Template name validation
  template_name_validation {
    rule {
      name = "template_name_format"
      description = "Template names must contain only alphanumeric characters, dots, underscores, and hyphens"
      condition = "name != null && !name.matches('^[a-zA-Z0-9._-]+$')"
      message = "Template names must contain only alphanumeric characters, dots, underscores, and hyphens"
      severity = "error"
    }
    
    rule {
      name = "template_name_length"
      description = "Template names must be between 1 and 128 characters"
      condition = "name != null && (name.length() < 1 || name.length() > 128)"
      message = "Template names must be between 1 and 128 characters"
      severity = "error"
    }
  }
  
  # Source path validation
  source_path_validation {
    rule {
      name = "source_path_format"
      description = "Source path must be in templates/ directory with .tmpl extension"
      condition = "source_path != null && !source_path.matches('^templates/.*\\.tmpl$')"
      message = "Source path must be in templates/ directory with .tmpl extension"
      severity = "error"
    }
  }
  
  # Destination path validation
  destination_path_validation {
    rule {
      name = "destination_path_format"
      description = "Destination path must be a valid file path"
      condition = "destination_path != null && !destination_path.matches('^[a-zA-Z0-9/._-]+$')"
      message = "Destination path must be a valid file path"
      severity = "error"
    }
  }
  
  # Description validation
  description_validation {
    rule {
      name = "description_length"
      description = "Template description must not exceed 500 characters"
      condition = "description != null && description.length() > 500"
      message = "Template description must not exceed 500 characters"
      severity = "error"
    }
  }
  
  # Tags validation
  tags_validation {
    rule {
      name = "tags_count_limit"
      description = "Templates cannot have more than 10 tags"
      condition = "tags != null && tags.size() > 10"
      message = "Templates cannot have more than 10 tags"
      severity = "warning"
    }
    
    rule {
      name = "tag_format"
      description = "Tags must contain only alphanumeric characters, underscores, and hyphens"
      condition = "tags != null && tags.any(tag -> !tag.matches('^[a-zA-Z0-9_-]+$'))"
      message = "Tags must contain only alphanumeric characters, underscores, and hyphens"
      severity = "error"
    }
    
    rule {
      name = "tag_length"
      description = "Tags must be between 1 and 32 characters"
      condition = "tags != null && tags.any(tag -> tag.length() < 1 || tag.length() > 32)"
      message = "Tags must be between 1 and 32 characters"
      severity = "error"
    }
  }
  
  # Enable flags validation
  enable_flags_validation {
    rule {
      name = "enable_machines_boolean"
      description = "Enable machines flag must be a boolean value"
      condition = "enable_machines != null && typeof(enable_machines) != 'boolean'"
      message = "Enable machines flag must be a boolean value"
      severity = "error"
    }
    
    rule {
      name = "enable_facts_boolean"
      description = "Enable facts flag must be a boolean value"
      condition = "enable_facts != null && typeof(enable_facts) != 'boolean'"
      message = "Enable facts flag must be a boolean value"
      severity = "error"
    }
    
    rule {
      name = "enable_variables_boolean"
      description = "Enable variables flag must be a boolean value"
      condition = "enable_variables != null && typeof(enable_variables) != 'boolean'"
      message = "Enable variables flag must be a boolean value"
      severity = "error"
    }
    
    rule {
      name = "enable_environment_boolean"
      description = "Enable environment flag must be a boolean value"
      condition = "enable_environment != null && typeof(enable_environment) != 'boolean'"
      message = "Enable environment flag must be a boolean value"
      severity = "error"
    }
  }
  
  # Constraints validation
  constraints_validation {
    rule {
      name = "max_size_reasonable"
      description = "Maximum template size must be between 1 and 10 MB"
      condition = "max_size_mb != null && (max_size_mb < 1 || max_size_mb > 10)"
      message = "Maximum template size must be between 1 and 10 MB"
      severity = "warning"
    }
    
    rule {
      name = "max_execution_time_reasonable"
      description = "Maximum execution time must be between 1 and 300 seconds"
      condition = "max_execution_time_seconds != null && (max_execution_time_seconds < 1 || max_execution_time_seconds > 300)"
      message = "Maximum execution time must be between 1 and 300 seconds"
      severity = "warning"
    }
    
    rule {
      name = "max_memory_reasonable"
      description = "Maximum memory usage must be between 1 and 100 MB"
      condition = "max_memory_mb != null && (max_memory_mb < 1 || max_memory_mb > 100)"
      message = "Maximum memory usage must be between 1 and 100 MB"
      severity = "warning"
    }
  }
  
  # Metadata validation
  metadata_validation {
    rule {
      name = "version_format"
      description = "Version must contain only alphanumeric characters, dots, underscores, and hyphens"
      condition = "version != null && !version.matches('^[a-zA-Z0-9._-]+$')"
      message = "Version must contain only alphanumeric characters, dots, underscores, and hyphens"
      severity = "error"
    }
    
    rule {
      name = "version_length"
      description = "Version must not exceed 32 characters"
      condition = "version != null && version.length() > 32"
      message = "Version must not exceed 32 characters"
      severity = "error"
    }
    
    rule {
      name = "author_length"
      description = "Author must not exceed 100 characters"
      condition = "author != null && author.length() > 100"
      message = "Author must not exceed 100 characters"
      severity = "warning"
    }
    
    rule {
      name = "timestamp_format"
      description = "Timestamps must be in ISO 8601 format"
      condition = "(created_at != null && !created_at.matches('^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$')) || (updated_at != null && !updated_at.matches('^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$'))"
      message = "Timestamps must be in ISO 8601 format"
      severity = "error"
    }
  }
  
  # Cross-field validation
  cross_field_validation {
    rule {
      name = "created_before_updated"
      description = "Created timestamp must be before or equal to updated timestamp"
      condition = "created_at != null && updated_at != null && created_at > updated_at"
      message = "Created timestamp must be before or equal to updated timestamp"
      severity = "error"
    }
    
    rule {
      name = "at_least_one_enable_flag"
      description = "At least one enable flag should be true for template context"
      condition = "enable_machines == false && enable_facts == false && enable_variables == false && enable_environment == false"
      message = "At least one enable flag should be true for template context"
      severity = "warning"
    }
  }
}
