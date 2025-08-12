# Project-Specific Logging Configuration Example
# This file should be placed in your project directory as ./logging.hcl
# It overrides global logging settings for this specific project

logging {
  # Override global log level for this project
  level = "debug"
  
  # Use JSON format for this project (overrides global structured format)
  format = "json"
  
  # Console output with project-specific settings
  output {
    type = "console"
    enabled = true
    colorize = true
    
    console {
      include_timestamp = true
      include_level = true
      include_component = true
    }
  }
  
  # Project-specific log file
  output {
    type = "file"
    enabled = true
    path = "./logs/project.log"
    
    file {
      max_size = "50MB"
      max_age = "7d"
      max_backups = 5
      compress = true
      permissions = "0644"
    }
  }
  
  # Error-only log file for this project
  output {
    type = "file"
    enabled = true
    path = "./logs/errors.log"
    level = "error"  # Only log errors to this file
    
    file {
      max_size = "25MB"
      max_age = "30d"
      max_backups = 10
      compress = true
    }
  }
  
  # Component-specific overrides for this project
  components {
    # Facts component - very verbose for debugging
    "facts" {
      level = "trace"  # Most detailed logging
      enabled = true
      
      output {
        type = "file"
        enabled = true
        path = "./logs/facts.log"
        
        file {
          max_size = "100MB"
          max_age = "3d"
          max_backups = 3
        }
      }
    }
    
    # Machines component - debug level for connectivity issues
    "machines" {
      level = "debug"
      enabled = true
      
      output {
        type = "file"
        enabled = true
        path = "./logs/machines.log"
      }
    }
    
    # Variables component - info level for this project
    "variables" {
      level = "info"  # Override global warn level
      enabled = true
    }
    
    # Actions component - debug level for development
    "actions" {
      level = "debug"
      enabled = true
      
      output {
        type = "file"
        enabled = true
        path = "./logs/actions.log"
      }
    }
    
    # Templates component - debug level for template issues
    "templates" {
      level = "debug"
      enabled = true
    }
    
    # SSH component - info level for connection debugging
    "ssh" {
      level = "info"  # Override global warn level
      enabled = true
      
      output {
        type = "file"
        enabled = true
        path = "./logs/ssh.log"
      }
    }
    
    # Config component - warn level (less verbose)
    "config" {
      level = "warn"
      enabled = true
    }
    
    # CLI component - info level
    "cli" {
      level = "info"
      enabled = true
    }
  }
  
  # Project-specific performance logging
  performance {
    enabled = true
    threshold_ms = 500  # Log operations taking > 500ms (lower than global)
    include_memory = true
    include_cpu = true
  }
  
  # Project-specific audit logging
  audit {
    enabled = true
    level = "info"
    include_user = true
    include_ip = true
    include_session = true
  }
}
