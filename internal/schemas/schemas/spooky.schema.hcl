# Spooky Global Configuration Schema
# Schema for $XDG_CONFIG_HOME/spooky/spooky.hcl

# Schema metadata
metadata {
  schema_version = "0.20250809.0"
  schema_type = "spooky"
  schema_name = "Spooky Global Configuration Schema"
  last_updated = "2024-01-01"
  compatibility = ["0.20250809.0"]
  description = "Global configuration schema for spooky CLI - defines configuration structure for $XDG_CONFIG_HOME/spooky/spooky.hcl"
  
  # ScalVer format: 0.YYYYMMDD.N
  # - 0: Development phase
  # - 20250809: Date (9 August 2025)
  # - 0: Patch version
  scalver_format = "0.20250809.0"
}

# Global configuration block structure
spooky {
  # Storage configuration
  storage {
    type = {
      type = "string"
      required = false
      enum = ["memory", "json"]
      default = "memory"
      description = "Storage backend type (memory for in-memory storage)"
    }
    
    path = {
      type = "string"
      required = false
      default = ""
      description = "Storage path for facts export (not used for memory storage)"
    }
    
    compression = {
      type = "boolean"
      required = false
      default = true
      description = "Enable compression for storage"
    }
    
    encryption = {
      type = "boolean"
      required = false
      default = false
      description = "Enable age encryption for sensitive data"
    }
    
    backup_enabled = {
      type = "boolean"
      required = false
      default = true
      description = "Enable automatic backups"
    }
    
    backup_retention = {
      type = "integer"
      required = false
      min = 1
      max = 365
      default = 7
      description = "Number of backup files to retain"
    }
  }

  # Facts collection configuration
  facts {
    timeout = {
      type = "integer"
      required = false
      min = 1
      max = 3600
      default = 30
      description = "Timeout for facts collection in seconds"
    }
    
    cache_ttl = {
      type = "integer"
      required = false
      min = 0
      max = 86400
      default = 3600
      description = "Cache TTL for facts in seconds (0 = no cache)"
    }
    
    auto_collect = {
      type = "boolean"
      required = false
      default = false
      description = "Automatically collect facts on machine operations"
    }
    
    parallel_collection = {
      type = "integer"
      required = false
      min = 1
      max = 100
      default = 10
      description = "Number of parallel facts collection workers"
    }
    
    retry_attempts = {
      type = "integer"
      required = false
      min = 0
      max = 10
      default = 3
      description = "Number of retry attempts for failed facts collection"
    }
    
    retry_delay = {
      type = "integer"
      required = false
      min = 1
      max = 60
      default = 5
      description = "Delay between retry attempts in seconds"
    }
  }

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

  # Template configuration
  templates {
    max_size = {
      type = "integer"
      required = false
      min = 1024
      max = 10485760
      default = 1048576
      description = "Maximum template file size in bytes"
    }
    
    allow_external_functions = {
      type = "boolean"
      required = false
      default = false
      description = "Allow external function calls in templates"
    }
    
    timeout = {
      type = "integer"
      required = false
      min = 1
      max = 300
      default = 30
      description = "Template rendering timeout in seconds"
    }
    
    cache_compiled = {
      type = "boolean"
      required = false
      default = true
      description = "Cache compiled templates"
    }
    
    sandbox_mode = {
      type = "boolean"
      required = false
      default = true
      description = "Enable sandbox mode for template running"
    }
    
    allowed_functions = {
      type = "list"
      required = false
      default = ["len", "join", "split", "replace", "trim", "upper", "lower"]
      description = "List of allowed functions in templates"
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
    
    restrict_file_access = {
      type = "boolean"
      required = false
      default = true
      description = "Restrict file access to project directory"
    }
    
    validate_ssh_keys = {
      type = "boolean"
      required = false
      default = true
      description = "Validate SSH keys before use"
    }
    
    audit_logging = {
      type = "boolean"
      required = false
      default = true
      description = "Enable audit logging for security events"
    }
    
    allowed_hosts = {
      type = "list"
      required = false
      default = []
      description = "List of allowed host patterns (CIDR or hostnames)"
    }
    
    blocked_hosts = {
      type = "list"
      required = false
      default = []
      description = "List of blocked host patterns (CIDR or hostnames)"
    }
  }

  # Age encryption configuration
  age {
    identities = {
      type = "string"
      required = false
      default = "~/.config/spooky/identities"
      description = "Path to directory containing age identity files (default: ~/.config/spooky/identities)"
    }
    
    recipients = {
      type = "string"
      required = false
      default = "~/.config/spooky/recipients.txt"
      description = "Path to file containing age recipients (public keys, one per line)"
    }
    
    passphrase = {
      type = "string"
      required = false
      description = "Passphrase for age encryption (if not using identity files)"
    }
    
    validation = {
      type = "object"
      required = false
      default = {
        strict_mode = true
        check_recipients = true
        validate_keys = true
      }
      description = "Age validation settings"
      
      properties = {
        strict_mode = {
          type = "boolean"
          required = false
          default = true
          description = "Enable strict validation mode"
        }
        
        check_recipients = {
          type = "boolean"
          required = false
          default = true
          description = "Validate recipient keys on encryption"
        }
        
        validate_keys = {
          type = "boolean"
          required = false
          default = true
          description = "Validate identity keys on decryption"
        }
      }
    }
    
    encryption = {
      type = "object"
      required = false
      default = {
        algorithm = "age"
        compression = true
        armor = false
      }
      description = "Age encryption settings"
      
      properties = {
        algorithm = {
          type = "string"
          required = false
          enum = ["age"]
          default = "age"
          description = "Encryption algorithm (only age supported)"
        }
        
        compression = {
          type = "boolean"
          required = false
          default = true
          description = "Enable compression before encryption"
        }
        
        armor = {
          type = "boolean"
          required = false
          default = false
          description = "Use armored output format"
        }
      }
    }
  }

  # Logging configuration
  logging {
    level = {
      type = "string"
      required = false
      enum = ["debug", "info", "warn", "error"]
      default = "info"
      description = "Default logging level"
    }
    
    format = {
      type = "string"
      required = false
      enum = ["json", "text"]
      default = "text"
      description = "Log format"
    }
    
    output = {
      type = "string"
      required = false
      default = "stderr"
      description = "Log output destination"
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
  }

  # Performance configuration
  performance {
    default_parallel = {
      type = "integer"
      required = false
      min = 2
      max = 100
      default = 10
      description = "Default parallel run limit"
    }
    
    max_memory = {
      type = "integer"
      required = false
      min = 64
      max = 8192
      default = 512
      description = "Maximum memory usage in MB"
    }
    
    gc_interval = {
      type = "integer"
      required = false
      min = 1
      max = 3600
      default = 300
      description = "Garbage collection interval in seconds"
    }
  }
  
} 