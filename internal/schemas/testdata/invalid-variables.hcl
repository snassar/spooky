# Invalid Variables Configuration
# This file should fail schema validations

metadata {
  version = 1
  description = "A test variables configuration that should fail validations"
}

variables {
  # Invalid environment variable
  environment = {
    type = "string"
    value = "invalid-env"  # Invalid: not in enum
    description = "Deployment environment"
    sensitive = false
    validation {
      enum = ["development", "staging", "production"]
      required = true
    }
  }
  
  # Invalid app version
  app_version = {
    type = "string"
    value = "invalid-version"  # Invalid: doesn't match pattern
    description = "Application version to deploy"
    sensitive = false
    validation {
      pattern = "^\\d+\\.\\d+\\.\\d+$"
      required = true
    }
  }
  
  # Invalid database port
  db_port = {
    type = "integer"
    value = 99999  # Invalid: too high
    description = "Database port"
    sensitive = false
    validation {
      min = 1
      max = 65535
      required = true
    }
  }
  
  # Invalid password (too short)
  db_password = {
    type = "string"
    value = "123"  # Invalid: too short
    description = "Database password"
    sensitive = true
    validation {
      min_length = 8
      required = true
    }
  }
  
  # Invalid max connections
  max_connections = {
    type = "integer"
    value = 9999  # Invalid: too high
    description = "Maximum database connections"
    sensitive = false
    validation {
      min = 1
      max = 1000
      required = true
    }
  }
  
  # Invalid list (too many items)
  allowed_ips = {
    type = "list"
    value = [
      "192.168.1.0/24", 
      "10.0.0.0/8", 
      "172.16.0.0/12",
      "8.8.8.8/32",
      "1.1.1.1/32",
      "9.9.9.9/32",
      "208.67.222.222/32",
      "208.67.220.220/32",
      "8.8.4.4/32",
      "1.0.0.1/32",
      "extra-ip"  # This makes it 11 items, exceeding max of 10
    ]
    description = "Allowed IP ranges"
    sensitive = false
    validation {
      min_items = 1
      max_items = 10
      required = true
    }
  }
  
  # Invalid file path
  log_file = {
    type = "string"
    value = "relative/path.log"  # Invalid: not absolute path
    description = "Application log file path"
    sensitive = false
    validation {
      pattern = "^/[a-zA-Z0-9/._-]+$"
      required = true
    }
  }
}
