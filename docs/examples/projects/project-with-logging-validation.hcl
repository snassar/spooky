# Example project.hcl with logging configuration
# This example demonstrates how to configure project-specific logging
# that overrides global logging settings

project "example-project" {
  name = "example-project"
  description = "Example project demonstrating logging configuration"
  version = "1.0.0"
  author = "spooky-user"
  email = "user@example.com"
  url = "https://github.com/example/spooky-project"

  # Project run settings
  run {
    default_timeout = 300
    max_parallel = 10
    dry_run_default = false
    validate_before_run = true
    backup_before_changes = false
  }

  # Facts collection configuration
  facts {
    timeout = 60
    auto_collect = true
    parallel_collection = 5
    retry_attempts = 3
    retry_delay = 5
    storage_format = "badgerdb"
    compression = true
    encryption = false
  }

  # Project-specific logging configuration
  # This overrides global logging settings from ~/.config/spooky/logging.hcl
  logging {
    # Override global log level for this project
    level = "debug"
    
    # Override global format for this project
    format = "text"
    
    # Override global output destination for this project
    output = "file"
    
    # File output configuration for this project
    file {
      # Relative path - will be created in project directory
      path = "./logs/project.log"
      
      # File permissions in octal format
      permissions = "0644"
      
      # Start fresh each run (don't append to existing file)
      append = false
    }
    
    # Component-specific filtering for this project
    filtering {
      components = {
        # Very verbose SSH logging for debugging
        "ssh"     = "debug"
        
        # Standard facts logging
        "facts"   = "info"
        
        # Only warnings and errors for actions
        "actions" = "warn"
        
        # Debug project operations
        "project" = "debug"
        
        # Only errors for machine operations
        "machines" = "error"
      }
      
      # Pattern-based filtering (optional)
      patterns = [
        "password=***",
        "token=***",
        "secret=***"
      ]
    }
    
    # Log rotation configuration for this project
    rotation {
      # Enable rotation for this project
      enabled = true
      
      # Maximum file size before rotation
      max_size = "50MB"
      
      # Maximum age of log files
      max_age = "7d"
      
      # Keep up to 5 backup files
      max_backups = 5
      
      # Compress rotated files
      compress = true
    }
  }

  # Additional project metadata
  metadata = {
    environment = "development"
    team = "platform"
    priority = "high"
  }
}
