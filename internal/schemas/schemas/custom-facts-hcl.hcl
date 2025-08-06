# Spooky Custom Facts HCL Schema
# Schema for custom facts stored in HCL format in /etc/spooky/facts/ on target machines
#
# NOTE: Any string field may be stored as an age-encrypted value (pattern ^age1[0-9a-z]+$) if 'encrypted = true'.
# If a field is marked as 'encrypted = true', its value must match the age-encrypted pattern.

# Custom facts HCL structure
facts {
  # File naming convention
  file_pattern = {
    type = "string"
    required = true
    pattern = "^facts\\.hcl$"
    description = "Custom facts must be in facts.hcl file"
  }
  
  # File location
  file_location = {
    type = "string"
    required = true
    pattern = "^/etc/spooky/facts\\.hcl$"
    description = "Custom facts must be in /etc/spooky/facts.hcl file"
  }
  
  # HCL content structure
  hcl_content = {
    type = "object"
    required = true
    description = "HCL content for custom facts"
    
    properties = {
      # Application facts
      app_name = {
        type = "string"
        required = false
        description = "Application name"
        pattern = "^[a-zA-Z0-9_-]+$"
      }
      
      app_version = {
        type = "string"
        required = false
        description = "Application version"
        pattern = "^[0-9]+\\.[0-9]+\\.[0-9]+(-[a-zA-Z0-9]+)?$"
      }
      
      config_path = {
        type = "string"
        required = false
        description = "Configuration file path"
        pattern = "^/[a-zA-Z0-9/_-]+$"
      }
      
      log_path = {
        type = "string"
        required = false
        description = "Log file path"
        pattern = "^/[a-zA-Z0-9/_-]+$"
      }
      
      # Environment facts
      environment = {
        type = "string"
        required = false
        description = "Environment name"
        enum = ["development", "staging", "production", "testing"]
      }
      
      datacenter = {
        type = "string"
        required = false
        description = "Datacenter identifier"
        pattern = "^[a-zA-Z0-9_-]+$"
      }
      
      rack = {
        type = "string"
        required = false
        description = "Rack identifier"
        pattern = "^[A-Z][0-9]{2}$"
      }
      
      power_zone = {
        type = "string"
        required = false
        description = "Power zone identifier"
        pattern = "^PZ-[0-9]+$"
      }
      
      region = {
        type = "string"
        required = false
        description = "Cloud region"
        pattern = "^[a-z]+-[a-z]+-[0-9]+$"
      }
      
      # Monitoring facts
      prometheus_port = {
        type = "integer"
        required = false
        description = "Prometheus metrics port"
        min = 1024
        max = 65535
      }
      
      alert_manager = {
        type = "string"
        required = false
        description = "Alert manager URL"
        pattern = "^[a-zA-Z0-9.-]+(:[0-9]+)?$"
      }
      
      log_level = {
        type = "string"
        required = false
        description = "Logging level"
        enum = ["debug", "info", "warn", "error", "fatal"]
      }
      
      metrics_enabled = {
        type = "bool"
        required = false
        description = "Whether metrics collection is enabled"
      }
      
      # Deployment facts (from dynamic scripts)
      deployment_state = {
        type = "string"
        required = false
        description = "Deployment state"
        enum = ["active", "inactive", "failed", "starting", "stopping"]
      }
      
      last_deploy = {
        type = "integer"
        required = false
        description = "Last deployment timestamp (Unix epoch)"
        min = 0
      }
      
      uptime = {
        type = "string"
        required = false
        description = "System uptime"
        pattern = "^[0-9]+ [a-z]+(, [0-9]+ [a-z]+)*$"
      }
      
      # SSL facts
      ssl_enabled = {
        type = "bool"
        required = false
        description = "Whether SSL is enabled"
      }
      
      ssl_cert = {
        type = "string"
        required = false
        description = "SSL certificate path"
        pattern = "^/[a-zA-Z0-9/_-]+\\.(crt|pem)$"
      }
      
      ssl_key = {
        type = "string"
        required = false
        description = "SSL private key path"
        pattern = "^/[a-zA-Z0-9/_-]+\\.(key|pem)$"
      }
      
      # Disk usage facts (from dynamic scripts)
      web_root_usage = {
        type = "integer"
        required = false
        description = "Web root disk usage percentage"
        min = 0
        max = 100
      }
      
      web_root_path = {
        type = "string"
        required = false
        description = "Web root directory path"
        pattern = "^/[a-zA-Z0-9/_-]+$"
      }
      
      # Network facts
      external_ip = {
        type = "string"
        required = false
        description = "External IP address"
        pattern = "^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$"
      }
      
      domain = {
        type = "string"
        required = false
        description = "Domain name"
        pattern = "^[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
      }
      
      # Service facts
      service_status = {
        type = "string"
        required = false
        description = "Service status"
        enum = ["running", "stopped", "failed", "starting", "stopping"]
      }
      
      service_enabled = {
        type = "bool"
        required = false
        description = "Whether service is enabled at boot"
      }

      # Example for a sensitive fact:
      database_password = {
        type = "string"
        required = false
        description = "Database password (plain or age-encrypted)"
        pattern = "^(age1[0-9a-z]+|.+)$"
        encrypted = {
          type = "bool"
          required = false
          default = false
          description = "Whether this fact is encrypted with age"
        }
      }
      # Repeat the above pattern for any other sensitive string fields (e.g., api_key, ssl_certificate, etc.)
    }
  }
  
  # File permissions
  file_permissions = {
    type = "string"
    required = true
    pattern = "^[0-7]{4}$"
    description = "File permissions (e.g., 0644 for static files, 0755 for executable scripts)"
  }
  
  # File size limits
  max_file_size = {
    type = "integer"
    required = true
    default = 1048576  # 1MB
    description = "Maximum file size in bytes"
    min = 1024
    max = 10485760  # 10MB
  }
  
  # Validation rules
  validation = {
    # File naming validation
    file_name_validation = {
      rule = "regex"
      pattern = "^facts\\.hcl$"
      message = "Custom facts must be in facts.hcl file"
    }
    
    # File location validation
    file_location_validation = {
      rule = "path"
      pattern = "^/etc/spooky/facts\\.hcl$"
      message = "Custom facts must be in /etc/spooky/facts.hcl file"
    }
    
    # HCL syntax validation
    hcl_syntax_validation = {
      rule = "hcl"
      message = "Custom fact files must contain valid HCL syntax"
    }
    
    # Executable script validation
    executable_validation = {
      rule = "executable"
      condition = "file_is_executable"
      message = "Executable fact scripts must be valid shell scripts that output HCL"
    }
    
    # Port number validation
    port_validation = {
      rule = "range"
      field = "prometheus_port"
      min = 1024
      max = 65535
      message = "Prometheus port must be between 1024 and 65535"
    }
    
    # IP address validation
    ip_validation = {
      rule = "regex"
      field = "external_ip"
      pattern = "^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$"
      message = "External IP must be a valid IPv4 address"
    }
    
    # Domain name validation
    domain_validation = {
      rule = "regex"
      field = "domain"
      pattern = "^[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
      message = "Domain must be a valid domain name"
    }

    encrypted_value_format = {
      rule = "conditional_regex"
      condition = "encrypted == true"
      pattern = "^age1[0-9a-z]+$"
      message = "Encrypted values must be valid age-encrypted strings"
    }
  }
  
  # Template integration
  template_context = {
    description = "How custom facts are made available in templates"
    context = {
      custom_facts = {
        description = "Custom facts from /etc/spooky/facts.hcl"
        access = "{{.facts.custom.facts.<key>}}"
      }
    }
    precedence = [
      "custom_facts",  # Custom facts (highest priority)
      "system_facts"   # System facts from gopsutil (lower priority)
    ]
  }
} 