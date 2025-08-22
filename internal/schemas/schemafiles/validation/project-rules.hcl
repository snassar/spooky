# Project Validation Rules
# Validation rules for project.hcl schema
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
    
    rule {
      name = "project_name_length"
      description = "Project names must be between 1 and 128 characters"
      condition = "name != null && (name.length() < 1 || name.length() > 128)"
      message = "Project names must be between 1 and 128 characters"
      severity = "error"
    }
  }
  
  # Project description validation
  project_description_validation {
    rule {
      name = "project_description_length"
      description = "Project description must not exceed 1024 characters"
      condition = "description != null && description.length() > 1024"
      message = "Project description must not exceed 1024 characters"
      severity = "error"
    }
  }
  
  # Run validation
  run_validation {
    rule {
      name = "run_max_parallel_reasonable"
      description = "Maximum parallel runs must be between 1 and 100"
      condition = "run_max_parallel != null && (run_max_parallel < 1 || run_max_parallel > 100)"
      message = "Maximum parallel runs must be between 1 and 100"
      severity = "warning"
    }
    
    rule {
      name = "run_dry_run_default_boolean"
      description = "Dry run default must be a boolean value"
      condition = "run_dry_run_default != null && typeof(run_dry_run_default) != 'boolean'"
      message = "Dry run default must be a boolean value"
      severity = "error"
    }
    
    rule {
      name = "run_validate_before_run_boolean"
      description = "Validate before run must be a boolean value"
      condition = "run_validate_before_run != null && typeof(run_validate_before_run) != 'boolean'"
      message = "Validate before run must be a boolean value"
      severity = "error"
    }
    
    rule {
      name = "run_backup_before_changes_boolean"
      description = "Backup before changes must be a boolean value"
      condition = "run_backup_before_changes != null && typeof(run_backup_before_changes) != 'boolean'"
      message = "Backup before changes must be a boolean value"
      severity = "error"
    }
  }
  
  # Facts validation
  facts_validation {
    rule {
      name = "facts_timeout_reasonable"
      description = "Facts collection timeouts must be between 1 and 3600 seconds"
      condition = "facts_timeout != null && (facts_timeout < 1 || facts_timeout > 3600)"
      message = "Facts collection timeouts must be between 1 and 3600 seconds"
      severity = "warning"
    }
    
    rule {
      name = "facts_parallel_reasonable"
      description = "Facts parallel collection must be between 1 and 100 workers"
      condition = "facts_parallel_collection != null && (facts_parallel_collection < 1 || facts_parallel_collection > 100)"
      message = "Facts parallel collection must be between 1 and 100 workers"
      severity = "warning"
    }
    
    rule {
      name = "facts_retry_attempts_reasonable"
      description = "Facts retry attempts must be between 0 and 10"
      condition = "facts_retry_attempts != null && (facts_retry_attempts < 0 || facts_retry_attempts > 10)"
      message = "Facts retry attempts must be between 0 and 10"
      severity = "warning"
    }
    
    rule {
      name = "facts_retry_delay_reasonable"
      description = "Facts retry delay must be between 1 and 60 seconds"
      condition = "facts_retry_delay != null && (facts_retry_delay < 1 || facts_retry_delay > 60)"
      message = "Facts retry delay must be between 1 and 60 seconds"
      severity = "warning"
    }
  }
  
  # Cross-field validation
  cross_field_validation {
    rule {
      name = "facts_parallel_less_than_run_parallel"
      description = "Facts parallel collection should not exceed run max parallel"
      condition = "facts_parallel_collection != null && run_max_parallel != null && facts_parallel_collection > run_max_parallel"
      message = "Facts parallel collection should not exceed run max parallel"
      severity = "warning"
    }
    
    rule {
      name = "facts_timeout_less_than_overall_timeout"
      description = "Facts timeout should be less than overall timeout"
      condition = "facts_timeout != null && facts_timeout > 3600"
      message = "Facts timeout should be reasonable for collection operations"
      severity = "warning"
    }
  }
}
