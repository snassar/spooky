# Spooky Global Configuration Schema
# Schema for $XDG_CONFIG_HOME/spooky/spooky.hcl

# Global configuration block structure
spooky {
  # Storage configuration
  storage {
    type = {
      type = "string"
      required = false
      enum = ["badgerdb", "json"]
      default = "badgerdb"
      description = "Storage backend type"
    }
    
    path = {
      type = "string"
      required = false
      default = "$XDG_DATA_HOME/spooky/facts.db"
      description = "Storage path for facts database"
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
    enabled = {
      type = "boolean"
      required = false
      default = false
      description = "Enable age encryption for sensitive values"
    }
    
    public_key = {
      type = "string"
      required = false
      description = "Age public key for encryption"
    }
    
    private_key_path = {
      type = "string"
      required = false
      description = "Path to age private key file"
    }
    
    recipients = {
      type = "list"
      required = false
      default = []
      description = "List of age recipients (public keys)"
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