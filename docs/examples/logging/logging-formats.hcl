# Logging Formats Example
# This file demonstrates different output format configurations
# Copy the relevant section to your logging.hcl file

# JSON Format Configuration
# Best for machine processing and log aggregation systems
logging {
  level = "info"
  format = "json"
  
  json {
    pretty_print = false        # Compact JSON for production
    include_timestamp = true
    include_level = true
    include_component = true
  }
  
  output {
    type = "console"
    enabled = true
  }
  
  output {
    type = "file"
    enabled = true
    path = "./logs/app.json"
  }
}

# Structured Format Configuration
# Best for human-readable logs with structured fields
logging {
  level = "info"
  format = "structured"
  
  structured {
    colorize = true             # Enable colors for console
    include_timestamp = true
    include_level = true
    include_component = true
    timestamp_format = "2006-01-02T15:04:05.000Z"
  }
  
  output {
    type = "console"
    enabled = true
  }
  
  output {
    type = "file"
    enabled = true
    path = "./logs/app.log"
  }
}

# Plain Text Format Configuration
# Simple text output for basic logging needs
logging {
  level = "info"
  format = "plain"
  
  plain {
    include_timestamp = true
    include_level = true
    include_component = true
    timestamp_format = "2006-01-02 15:04:05"
  }
  
  output {
    type = "console"
    enabled = true
  }
  
  output {
    type = "file"
    enabled = true
    path = "./logs/app.txt"
  }
}

# Development Configuration
# Verbose logging with multiple formats for debugging
logging {
  level = "debug"
  format = "structured"
  
  # Console with colors for development
  output {
    type = "console"
    enabled = true
    level = "debug"
    
    console {
      colorize = true
      include_timestamp = true
      include_level = true
      include_component = true
    }
  }
  
  # JSON file for log analysis
  output {
    type = "file"
    enabled = true
    path = "./logs/debug.json"
    level = "debug"
    
    file {
      max_size = "100MB"
      max_age = "1d"
      max_backups = 3
    }
  }
  
  # Plain text file for quick viewing
  output {
    type = "file"
    enabled = true
    path = "./logs/debug.log"
    level = "info"
  }
}

# Production Configuration
# Optimized for production environments
logging {
  level = "warn"  # Only warnings and errors
  
  # JSON format for log aggregation
  format = "json"
  
  json {
    pretty_print = false
    include_timestamp = true
    include_level = true
    include_component = true
  }
  
  # File output only (no console in production)
  output {
    type = "file"
    enabled = true
    path = "/var/log/spooky/app.log"
    
    file {
      max_size = "100MB"
      max_age = "30d"
      max_backups = 10
      compress = true
      permissions = "0644"
    }
  }
  
  # Separate error log
  output {
    type = "file"
    enabled = true
    path = "/var/log/spooky/errors.log"
    level = "error"
    
    file {
      max_size = "50MB"
      max_age = "90d"
      max_backups = 20
      compress = true
    }
  }
}

# Minimal Configuration
# Simple setup for basic logging needs
logging {
  level = "info"
  format = "structured"
  
  output {
    type = "console"
    enabled = true
  }
}

# Component-Specific Format Configuration
# Different formats for different components
logging {
  level = "info"
  format = "structured"  # Default format
  
  components {
    # Facts component - JSON for analysis
    "facts" {
      level = "debug"
      enabled = true
      
      output {
        type = "file"
        enabled = true
        path = "./logs/facts.json"
        
        # Override format for this component
        format = "json"
        json {
          pretty_print = true
          include_timestamp = true
          include_level = true
          include_component = true
        }
      }
    }
    
    # Machines component - Structured for readability
    "machines" {
      level = "info"
      enabled = true
      
      output {
        type = "file"
        enabled = true
        path = "./logs/machines.log"
        
        # Use structured format
        format = "structured"
        structured {
          colorize = false
          include_timestamp = true
          include_level = true
          include_component = true
        }
      }
    }
    
    # Variables component - Plain text for simplicity
    "variables" {
      level = "warn"
      enabled = true
      
      output {
        type = "file"
        enabled = true
        path = "./logs/variables.txt"
        
        # Use plain text format
        format = "plain"
        plain {
          include_timestamp = true
          include_level = true
          include_component = true
        }
      }
    }
  }
}
