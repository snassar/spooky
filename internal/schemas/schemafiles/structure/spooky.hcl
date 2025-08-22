# Spooky Global Configuration Schema
# Schema for $XDG_CONFIG_HOME/spooky/spooky.hcl

# Schema metadata
metadata {
  version = "1"
  description = "Global configuration schema for spooky CLI"
}

# Global configuration block structure
spooky {
  # SSH configuration
  ssh {
    timeout = {
      type = "integer"
      required = false
      min = 1
      max = 300
      default = 30
      description = "SSH connection timeout in seconds"
    }
    
    keepalive_interval = {
      type = "integer"
      required = false
      min = 1
      max = 300
      default = 60
      description = "SSH keepalive interval in seconds"
    }
    
    keepalive_count = {
      type = "integer"
      required = false
      min = 1
      max = 10
      default = 3
      description = "SSH keepalive count before considering connection dead"
    }
    
    key_scan_timeout = {
      type = "integer"
      required = false
      min = 1
      max = 60
      default = 10
      description = "SSH key scanning timeout in seconds"
    }
    
    known_hosts_strict = {
      type = "boolean"
      required = false
      default = true
      description = "Strict known_hosts checking"
    }
    
    connection_pool_size = {
      type = "integer"
      required = false
      min = 1
      max = 100
      default = 10
      description = "SSH connection pool size"
    }
  }



  # Security configuration
  security {
    allow_unsafe_commands = {
      type = "boolean"
      required = false
      default = false
      description = "Allow potentially unsafe commands"
    }
    

    
    audit_logging = {
      type = "boolean"
      required = false
      default = true
      description = "Enable audit logging for security events"
    }
  }

  # Age encryption configuration
  age {
    identities = {
      type = "string"
      required = false
      description = "Path to directory containing age identity files"
    }
    
    recipients = {
      type = "string"
      required = false
      description = "Path to file containing age recipients (public keys, one per line)"
    }
    

  }

  # Logging configuration
  logging {
    # Basic logging configuration
    level = {
      type = "string"
      required = false
      enum = ["debug", "info", "warn", "error", "fatal"]
      default = "info"
      description = "Minimum log level to output (debug, info, warn, error, fatal)"
    }
    
    format = {
      type = "string"
      required = false
      enum = ["json", "text", "structured"]
      default = "json"
      description = "Log output format (json for machine-readable, text for human-readable, structured for custom)"
    }
    
    output = {
      type = "string"
      required = false
      enum = ["stdout", "stderr", "file", "null"]
      default = "stderr"
      description = "Log output destination (stdout, stderr, file, null for no output)"
    }
    
    color_output = {
      type = "boolean"
      required = false
      default = true
      description = "Enable colored output"
    }
    
    progress_bars = {
      type = "boolean"
      required = false
      default = true
      description = "Show progress bars for long operations"
    }
    
    # File output configuration
    file_path = {
      type = "string"
      required = false
      description = "Path to log file (required when output is 'file')"
      pattern = "^[^<>:\"/\\|?*]+$"
    }
    
    file_permissions = {
      type = "string"
      required = false
      default = "0644"
      pattern = "^[0-7]{3,4}$"
      description = "File permissions in octal format (e.g., 0644)"
    }
    
    file_append = {
      type = "boolean"
      required = false
      default = true
      description = "Whether to append to existing file or truncate"
    }
    
    # Component-specific logging configuration
    filtering_components = {
      type = "object"
      required = false
      additional_properties = {
        type = "string"
        enum = ["debug", "info", "warn", "error", "fatal"]
      }
      description = "Component-specific log level configuration"
    }
  }

  # Performance configuration
  performance {
    max_memory_mb = {
      type = "integer"
      required = false
      min = 64
      max = 8192
      default = 512
      description = "Maximum memory usage in MB for facts, variables, and template processing"
    }
    
    overall_timeout = {
      type = "integer"
      required = false
      min = 1
      max = 7200
      default = 3600
      description = "Global maximum timeout for all operations in seconds"
    }
  }
} 