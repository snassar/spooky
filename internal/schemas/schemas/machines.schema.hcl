# Spooky Machines Configuration Schema
# Comprehensive schema for machines.hcl files with enterprise-scale indexing and connectivity validation

# Schema metadata
metadata {
  schema_version = "0.20250809.0"
  schema_type = "machines"
  schema_name = "Spooky Machines Configuration Schema"
  last_updated = "2024-01-01"
  compatibility = ["0.20250809.0"]
  description = "Comprehensive schema for machines.hcl files with enterprise-scale indexing and connectivity validation"
  
  # ScalVer format: 0.YYYYMMDD.N
  # - 0: Development phase
  # - 20250809: Date (9 August 2025)
  # - 0: Patch version
  scalver_format = "0.20250809.0"
}

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
    
    passphrase = {
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
    
    roles = {
      type = "array"
      required = false
      description = "Functional roles this machine performs (e.g., web-server, database, load-balancer)"
      items = {
        type = "string"
        pattern = "^[a-zA-Z0-9._-]+$"
      }
    }
    
    classes = {
      type = "array"
      required = false
      description = "Configuration classes to apply to this machine"
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
    
    # Resource specifications (for capacity planning and validation)
    resources = {
      type = "object"
      required = false
      description = "Machine resource specifications"
      
      properties = {
        cpu_cores = {
          type = "integer"
          required = false
          min = 1
          max = 1024
          description = "Number of CPU cores"
        }
        
        memory_gb = {
          type = "integer"
          required = false
          min = 1
          max = 32768
          description = "Memory in GB"
        }
        
        disk_gb = {
          type = "integer"
          required = false
          min = 1
          max = 1048576
          description = "Disk space in GB"
        }
        
        network_speed = {
          type = "string"
          required = false
          pattern = "^[0-9]+(Gbps|Mbps)$"
          description = "Network speed (e.g., 10Gbps, 1Gbps)"
        }
      }
    }
    
    metadata = {
      type = "object"
      required = false
      description = "Additional machine metadata"
      additional_properties = "string"
    }
    
    environment = {
      type = "string"
      required = false
      pattern = "^[a-zA-Z0-9._-]+$"
      description = "Environment this machine belongs to (e.g., production, staging, development, testing)"
    }
    
    stage = {
      type = "string"
      required = false
      pattern = "^[a-zA-Z0-9._-]+$"
      description = "Deployment stage (e.g., blue, green, canary, main)"
    }
    
    # Lifecycle and maintenance information
    lifecycle = {
      type = "object"
      required = false
      description = "Machine lifecycle and maintenance information"
      
      properties = {
        status = {
          type = "string"
          required = false
          enum = ["active", "maintenance", "decommissioned", "standby", "testing"]
          description = "Current machine status - machines in 'maintenance' status are excluded from actions by default"
        }
        
        maintenance_window = {
          type = "object"
          required = false
          description = "Timezone-aware maintenance window"
          
          properties = {
            start_time = {
              type = "string"
              required = true
              pattern = "^[0-9]{2}:[0-9]{2}$"
              description = "Start time in HH:MM format"
            }
            
            end_time = {
              type = "string"
              required = true
              pattern = "^[0-9]{2}:[0-9]{2}$"
              description = "End time in HH:MM format"
            }
            
            timezone = {
              type = "string"
              required = true
              pattern = "^[A-Za-z_]+/[A-Za-z_]+$"
              description = "Timezone in IANA format (e.g., America/New_York, Europe/London, Asia/Tokyo)"
            }
            
            days_of_week = {
              type = "array"
              required = false
              default = ["sunday"]
              description = "Days of week when maintenance is allowed"
              items = {
                type = "string"
                enum = ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"]
              }
            }
            
            auto_exclude = {
              type = "boolean"
              required = false
              default = false
              description = "Automatically exclude machine from actions during maintenance window (overrides status-based exclusion)"
            }
            
            # Alternative: Simple string format for backward compatibility
            simple = {
              type = "string"
              required = false
              pattern = "^[0-9]{2}:[0-9]{2}-[0-9]{2}:[0-9]{2}$"
              description = "Simple maintenance window (e.g., 02:00-04:00) - local timezone only"
            }
          }
        }
        
        maintenance_team = {
          type = "string"
          required = false
          pattern = "^[a-zA-Z0-9._-]+$"
          description = "Team responsible for maintenance"
        }
        
        team_timezone = {
          type = "string"
          required = false
          pattern = "^[A-Za-z_]+/[A-Za-z_]+$"
          description = "Primary timezone of the maintenance team (IANA format)"
        }
        
        team_contact = {
          type = "object"
          required = false
          description = "Team contact information for maintenance coordination"
          
          properties = {
            primary = {
              type = "string"
              required = false
              description = "Primary contact (email or phone)"
            }
            
            secondary = {
              type = "string"
              required = false
              description = "Secondary contact (email or phone)"
            }
            
            slack_channel = {
              type = "string"
              required = false
              pattern = "^#[a-zA-Z0-9._-]+$"
              description = "Slack channel for maintenance coordination"
            }
            
            pagerduty_schedule = {
              type = "string"
              required = false
              description = "PagerDuty schedule ID for escalation"
            }
          }
        }
        
        warranty_expiry = {
          type = "string"
          required = false
          format = "date"
          description = "Hardware warranty expiry date (YYYY-MM-DD)"
        }
        
        retirement_date = {
          type = "string"
          required = false
          format = "date"
          description = "Planned retirement date (YYYY-MM-DD)"
        }
      }
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
    
    auth_mutual_exclusive = {
      rule = "conditional"
      condition = "!(password != null && key_file != null)"
      message = "Password and key_file authentication methods are mutually exclusive - specify only one"
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
    
    # Enhanced validation rules
    environment_validation = {
      rule = "conditional"
      condition = "environment != null && environment in ['production', 'staging', 'development', 'testing']"
      message = "Environment must be one of: production, staging, development, testing"
    }
    
    role_validation = {
      rule = "conditional"
      condition = "roles != null && len(roles) > 0"
      message = "At least one role should be specified for proper classification"
    }
    
    resource_validation = {
      rule = "conditional"
      condition = "resources != null && (resources.cpu_cores != null || resources.memory_gb != null)"
      message = "Resource specifications should include CPU cores and/or memory"
    }
    
    lifecycle_validation = {
      rule = "conditional"
      condition = "lifecycle != null && lifecycle.status != null"
      message = "Lifecycle status should be specified for proper management"
    }
    
    # Cross-field validation
    production_safety = {
      rule = "conditional"
      condition = "environment != 'production' || (lifecycle != null && lifecycle.status == 'active')"
      message = "Production machines must have active status"
    }
    
    maintenance_window_validation = {
      rule = "conditional"
      condition = "lifecycle == null || lifecycle.maintenance_window == null || (lifecycle.maintenance_window.simple != null || (lifecycle.maintenance_window.start_time != null && lifecycle.maintenance_window.end_time != null && lifecycle.maintenance_window.timezone != null))"
      message = "Maintenance window must be either simple format (HH:MM-HH:MM) or timezone-aware object with start_time, end_time, and timezone"
    }
    
    timezone_validation = {
      rule = "conditional"
      condition = "lifecycle == null || lifecycle.maintenance_window == null || lifecycle.maintenance_window.timezone == null || lifecycle.maintenance_window.timezone in ['America/New_York', 'America/Chicago', 'America/Denver', 'America/Los_Angeles', 'Europe/London', 'Europe/Paris', 'Europe/Berlin', 'Asia/Tokyo', 'Asia/Shanghai', 'Australia/Sydney', 'UTC']"
      message = "Timezone must be a valid IANA timezone identifier"
    }
  }
} 