# Global logging configuration for spooky
# This file should be placed at: ~/.config/spooky/logging.hcl
# It configures logging behavior for all spooky operations

logging {
  # Log level (debug, info, warn, error, fatal)
  level = "info"
  
  # Output format (json, text, structured)
  format = "json"
  
  # Output destination (stdout, stderr, file, null)
  output = "file"
  
  # File output configuration
  file {
    path        = "/var/log/spooky/spooky.log"
    permissions = "0644"
    append      = true
  }
  
  # Structured logging configuration
  structured {
    timestamp {
      enabled  = true
      format   = "RFC3339"
      timezone = "UTC"
    }
    
    component {
      key     = "component"
      enabled = true
    }
    
    operation {
      key     = "operation"
      enabled = true
    }
    
    error {
      key         = "error"
      include_stack = false
      include_type = true
    }
    
    fields {
      global = {
        service = "spooky"
        version = "0.20250812.0"
      }
      
      filter {
        sensitive = ["password", "token", "secret", "key"]
      }
    }
  }
  
  # Component-specific filtering
  filtering {
    components = {
      "ssh"     = "debug"   # Verbose SSH logging
      "facts"   = "info"    # Standard facts logging
      "actions" = "warn"    # Only warnings and errors for actions
    }
  }
  
  # Performance optimization
  performance {
    buffer {
      enabled        = true
      size           = 8192
      flush_interval = "5s"
    }
    
    async {
      enabled       = false  # Keep synchronous for now
      queue_size    = 1000
      workers       = 1
      drop_when_full = false
    }
  }
  
  # Log rotation
  rotation {
    enabled      = true
    max_size     = "100MB"
    max_age      = "30d"
    max_backups  = 5
    compress     = true
    local_time   = false
  }
}
