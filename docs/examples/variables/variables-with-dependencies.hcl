variables {
  # Base configuration variables
  variable "base_url" {
    type = "string"
    description = "Base URL for the application"
    default = "https://api.example.com"
    scope = "project"
  }

  variable "api_version" {
    type = "string"
    description = "API version"
    default = "v1"
    scope = "project"
  }

  variable "timeout" {
    type = "number"
    description = "Request timeout in seconds"
    default = 30
    scope = "project"
    
    constraints {
      min_value = 1
      max_value = 300
    }
  }

  variable "retry_count" {
    type = "number"
    description = "Number of retries"
    default = 3
    scope = "project"
    
    constraints {
      min_value = 0
      max_value = 10
    }
  }

  # Dependent variables
  variable "api_url" {
    type = "string"
    description = "Full API URL"
    scope = "project"
    dependencies = ["base_url", "api_version"]
  }

  variable "api_config" {
    type = "map"
    description = "API configuration"
    scope = "project"
    dependencies = ["api_url", "timeout", "retry_count"]
  }

  variable "app_name" {
    type = "string"
    description = "Application name"
    default = "my-app"
    scope = "project"
  }

  variable "app_title" {
    type = "string"
    description = "Application title"
    scope = "project"
    dependencies = ["app_name"]
  }

  variable "app_url" {
    type = "string"
    description = "Application URL"
    scope = "project"
    dependencies = ["base_url", "app_name"]
  }

  # Environment-specific variables
  variable "db_host" {
    type = "string"
    description = "Database host"
    default = "localhost"
    scope = "project"
  }

  variable "db_port" {
    type = "number"
    description = "Database port"
    default = 5432
    scope = "project"
  }

  variable "db_name" {
    type = "string"
    description = "Database name"
    scope = "project"
    dependencies = ["app_name", "environment"]
  }

  variable "db_url" {
    type = "string"
    description = "Database connection URL"
    scope = "project"
    dependencies = ["db_host", "db_port", "db_name"]
  }

  # Complex configuration
  variable "logging_config" {
    type = "map"
    description = "Logging configuration"
    scope = "project"
    dependencies = ["environment", "debug"]
  }

  variable "monitoring_config" {
    type = "map"
    description = "Monitoring configuration"
    scope = "project"
    dependencies = ["environment", "app_name"]
  }

  # Sensitive variables
  variable "api_key" {
    type = "string"
    description = "API authentication key"
    scope = "project"
    sensitive = true
  }

  variable "db_password" {
    type = "string"
    description = "Database password"
    scope = "project"
    sensitive = true
    encrypted = true
  }

  # Validation examples
  variable "email" {
    type = "string"
    description = "Contact email"
    scope = "project"
    
    validation {
      pattern = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
    }
  }

  variable "port" {
    type = "number"
    description = "Application port"
    default = 8080
    scope = "project"
    
    validation {
      min_value = 1024
      max_value = 65535
    }
  }
}
