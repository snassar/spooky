# Valid Variables Configuration
# This file should pass all schema validations

metadata {
  version = 1
  description = "A test variables configuration that should pass all validations"
}

variables {
  # Environment variables
  environment = {
    type = "string"
    value = "production"
    description = "Deployment environment"
    sensitive = false
    validation {
      enum = ["development", "staging", "production"]
      required = true
    }
  }
  
  app_version = {
    type = "string"
    value = "1.2.3"
    description = "Application version to deploy"
    sensitive = false
    validation {
      pattern = "^\\d+\\.\\d+\\.\\d+$"
      required = true
    }
  }
  
  # Database configuration
  db_host = {
    type = "string"
    value = "db.example.com"
    description = "Database hostname"
    sensitive = false
    validation {
      format = "hostname"
      required = true
    }
  }
  
  db_port = {
    type = "integer"
    value = 5432
    description = "Database port"
    sensitive = false
    validation {
      min = 1
      max = 65535
      required = true
    }
  }
  
  db_password = {
    type = "string"
    value = "encrypted:age1..."
    description = "Database password"
    sensitive = true
    validation {
      min_length = 8
      required = true
    }
  }
  
  # Application configuration
  max_connections = {
    type = "integer"
    value = 100
    description = "Maximum database connections"
    sensitive = false
    validation {
      min = 1
      max = 1000
      required = true
    }
  }
  
  debug_mode = {
    type = "boolean"
    value = false
    description = "Enable debug mode"
    sensitive = false
    validation {
      required = false
    }
  }
  
  # File paths
  log_file = {
    type = "string"
    value = "/var/log/app.log"
    description = "Application log file path"
    sensitive = false
    validation {
      pattern = "^/[a-zA-Z0-9/._-]+$"
      required = true
    }
  }
  
  # Lists and maps
  allowed_ips = {
    type = "list"
    value = ["192.168.1.0/24", "10.0.0.0/8"]
    description = "Allowed IP ranges"
    sensitive = false
    validation {
      min_items = 1
      max_items = 10
      required = true
    }
  }
  
  feature_flags = {
    type = "map"
    value = {
      "new_ui" = true
      "beta_features" = false
      "maintenance_mode" = false
    }
    description = "Feature flags configuration"
    sensitive = false
    validation {
      required = false
    }
  }
}
