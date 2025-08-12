# Example project.hcl with logging override
# This project overrides the global logging configuration for development purposes

project {
  name = "development-project"
  description = "A development project with custom logging"
  version = "1.0.0"
  author = "developer"

  execution {
    default_timeout = 300
    max_parallel = 10
    dry_run_default = false
    validate_before_execute = true
    backup_before_changes = false
  }

  facts {
    timeout = 60
    max_parallel = 5
    retry_attempts = 3
    retry_delay = 5
    storage_format = "badgerdb"
    compression = true
    encryption = false
  }

  # Override global logging configuration for this project
  logging {
    # More verbose logging for development
    level = "debug"
    
    # Use project-specific log file
    output = "file"
    file {
      path        = "./logs/project.log"
      permissions = "0644"
      append      = false  # Start fresh for each run
    }
    
    # Component-specific filtering for development
    filtering {
      components = {
        "ssh"     = "debug"   # Very verbose SSH logging
        "facts"   = "debug"   # Debug facts collection
        "actions" = "debug"   # Debug action execution
        "project" = "debug"   # Debug project operations
      }
    }
    
    # Disable rotation for development (keep all logs)
    rotation {
      enabled = false
    }
  }
}
