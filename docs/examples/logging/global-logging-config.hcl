# Global Logging Configuration Example
# This file should be placed at $XDG_CONFIG_HOME/spooky/logging.hcl
# (typically ~/.config/spooky/logging.hcl on Linux/macOS)

logging {
  # Global log level - applies to all components unless overridden
  level = "info"
  
  # Output format - json, structured, or plain
  format = "structured"
  
  # Console output for immediate feedback
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
  
  # File output for persistent storage
  output {
    type = "file"
    enabled = true
    path = "~/.local/state/spooky/logs/spooky.log"
    
    file {
      max_size = "100MB"
      max_age = "30d"
      max_backups = 10
      compress = true
      permissions = "0644"
    }
  }
  
  # Component-specific logging configuration
  components {
    # Facts component - more verbose for debugging
    "facts" {
      level = "debug"
      enabled = true
      
      output {
        type = "file"
        enabled = true
        path = "~/.local/state/spooky/logs/facts.log"
        
        file {
          max_size = "50MB"
          max_age = "7d"
          max_backups = 5
        }
      }
    }
    
    # Machines component - standard level
    "machines" {
      level = "info"
      enabled = true
    }
    
    # Variables component - warnings and errors only
    "variables" {
      level = "warn"
      enabled = true
    }
    
    # Actions component - debug level for development
    "actions" {
      level = "debug"
      enabled = true
    }
    
    # Templates component - info level
    "templates" {
      level = "info"
      enabled = true
    }
    
    # SSH component - warnings and errors only
    "ssh" {
      level = "warn"
      enabled = true
    }
    
    # Config component - info level
    "config" {
      level = "info"
      enabled = true
    }
    
    # CLI component - info level
    "cli" {
      level = "info"
      enabled = true
    }
  }
  
  # Performance logging configuration
  performance {
    enabled = true
    threshold_ms = 1000  # Log operations taking > 1 second
    include_memory = true
    include_cpu = true
  }
  
  # Audit logging configuration
  audit {
    enabled = true
    level = "info"
    include_user = true
    include_ip = true
    include_session = true
  }
}
