# Spooky Machines Configuration Schema
# Comprehensive schema for machines.hcl files with enterprise-scale indexing and connectivity validation

# Machines block structure
machines {
  # Machine definitions
  machine "machine_name" {
    hostname = {
      type = "string"
      required = true
      pattern = "^[a-zA-Z0-9.-]+$"
      description = "Machine hostname"
    }
    
    host = {
      type = "string"
      required = true
      format = "ipv4|ipv6|hostname"
      description = "Machine hostname or IP address"
    }
    
    port = {
      type = "integer"
      required = false
      min = 1
      max = 65535
      default = 22
      description = "SSH port number"
    }
    
    user = {
      type = "string"
      required = true
      pattern = "^[a-zA-Z0-9._-]+$"
      description = "SSH username"
    }
    
    password = {
      type = "string"
      required = false
      sensitive = true
      description = "SSH password (mutually exclusive with key_file)"
    }
    
    key_file = {
      type = "string"
      required = false
      pattern = "^[^/].*"
      description = "Path to SSH private key file (mutually exclusive with password)"
    }
    
    key_passphrase = {
      type = "string"
      required = false
      sensitive = true
      description = "Passphrase for SSH private key"
    }
    
    tags = {
      type = "object"
      required = false
      description = "Machine-specific tags for targeting and organization"
      additional_properties = "string"
    }
    
    groups = {
      type = "array"
      required = false
      description = "Machine groups for organization and targeting"
      items = {
        type = "string"
        pattern = "^[a-zA-Z0-9._-]+$"
      }
    }
    
    # SSH connection configuration (overrides global defaults)
    connection_timeout = {
      type = "integer"
      required = false
      min = 1
      max = 300
      default = 30
      description = "SSH connection timeout in seconds (overrides global default)"
    }
    
    command_timeout = {
      type = "integer"
      required = false
      min = 1
      max = 3600
      default = 300
      description = "Command execution timeout in seconds (overrides global default)"
    }
    
    max_connections = {
      type = "integer"
      required = false
      min = 1
      max = 100
      default = 10
      description = "Maximum concurrent SSH connections for this machine (overrides global default)"
    }
    
    retry_attempts = {
      type = "integer"
      required = false
      min = 0
      max = 10
      default = 3
      description = "Number of connection retry attempts (overrides global default)"
    }
    
    retry_delay = {
      type = "integer"
      required = false
      min = 1
      max = 60
      default = 5
      description = "Delay between retry attempts in seconds (overrides global default)"
    }
    
    metadata = {
      type = "object"
      required = false
      description = "Additional machine metadata"
      additional_properties = "string"
    }
  }
  
  # Validation rules
  validation = {
    # Machine name validation
    machine_name = {
      rule = "regex"
      pattern = "^[a-zA-Z][a-zA-Z0-9._-]*$"
      message = "Machine names must start with a letter and contain only alphanumeric characters, dots, underscores, and hyphens"
    }
    
    # Authentication validation
    auth_method = {
      rule = "conditional"
      condition = "password != null || key_file != null"
      message = "Machine must have either password or key_file authentication method"
    }
    
    # Host validation
    host_format = {
      rule = "format"
      format = "ipv4|ipv6|hostname"
      message = "Host must be a valid IPv4, IPv6, or hostname"
    }
    
    # Port validation
    port_range = {
      rule = "range"
      min = 1
      max = 65535
      message = "Port must be between 1 and 65535"
    }
    
    # Tag validation
    tag_names = {
      rule = "regex"
      pattern = "^[a-zA-Z][a-zA-Z0-9._-]*$"
      message = "Tag names must start with a letter and contain only alphanumeric characters, dots, underscores, and hyphens"
    }
    
    # Group validation
    group_names = {
      rule = "regex"
      pattern = "^[a-zA-Z][a-zA-Z0-9._-]*$"
      message = "Group names must start with a letter and contain only alphanumeric characters, dots, underscores, and hyphens"
    }
    
    # SSH configuration validation
    ssh_timeout_reasonable = {
      rule = "range"
      min = 1
      max = 3600
      message = "SSH timeouts must be between 1 and 3600 seconds"
    }
    
    ssh_connection_limit_reasonable = {
      rule = "range"
      min = 1
      max = 100
      message = "SSH connection limits must be between 1 and 100"
    }
    
    ssh_retry_reasonable = {
      rule = "range"
      min = 0
      max = 10
      message = "SSH retry attempts must be between 0 and 10"
    }
  }
} 