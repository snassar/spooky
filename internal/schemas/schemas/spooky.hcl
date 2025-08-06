# Spooky Global Configuration Schema
# Schema for $XDG_CONFIG_HOME/spooky/spooky.hcl

# Global configuration block structure
spooky {
  # Default logging level
  log_level = {
    type = "string"
    required = false
    enum = ["debug", "info", "warn", "error"]
    default = "error"
    description = "Default logging level"
  }

  # Default SSH user for machines without explicit user
  default_ssh_user = {
    type = "string"
    required = false
    default = "debian"
    description = "Default SSH user for machines without explicit user"
  }

  # Default SSH port for machines without explicit port
  default_ssh_port = {
    type = "integer"
    required = false
    min = 1
    max = 65535
    default = 22
    description = "Default SSH port for machines without explicit port"
  }

  # SSH connection timeout in seconds
  ssh_timeout = {
    type = "integer"
    required = false
    min = 1
    max = 300
    default = 30
    description = "SSH connection timeout in seconds"
  }

  # Default facts storage backend
  facts_storage = {
    type = "string"
    required = false
    enum = ["badgerdb", "json"]
    default = "badgerdb"
    description = "Default facts storage backend"
  }

  # Default parallel execution for facts gathering
  default_parallel = {
    type = "integer"
    required = false
    min = 2
    max = 100
    default = 10
    description = "Default parallel execution limit"
  }

  # Enable colored output
  color_output = {
    type = "boolean"
    required = false
    default = true
    description = "Enable colored output"
  }

  # Enable progress bars
  progress_bars = {
    type = "boolean"
    required = false
    default = true
    description = "Show progress bars for long operations"
  }
} 