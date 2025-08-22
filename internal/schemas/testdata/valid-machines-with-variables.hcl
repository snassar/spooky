# Valid Machines Configuration with Variables
# This file demonstrates per-machine and group variables with encryption support

metadata {
  version = 1
  description = "A test machines configuration demonstrating per-machine and group variables"
}

machines {
  machine {
    name = "web-server-01"
    description = "Primary web server with machine-specific variables"
    hostname = "web01.example.com"
    port = 22
    user = "admin"
    
    authentication {
      method = "ssh_key"
      ssh_key_path = "~/.ssh/id_rsa"
    }
    
    connection {
      timeout = 30
      retry_attempts = 3
      retry_delay = 5
    }
    
    tags = ["web", "production", "primary"]
    environment = "production"
    role = "web-server"
    location = "us-east-1"
    
    # Machine-specific variables
    variables {
      # Database connection string (encrypted)
      db_connection_string = {
        value = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
        encrypted = true
        encryption_metadata = {
          recipients = ["age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"]
          encrypted_at = "2024-01-01T12:00:00Z"
          method = "age"
        }
        description = "Database connection string for this specific web server"
        sensitive = true
        tags = ["database", "connection"]
      }
      
      # Application port (plain text)
      app_port = {
        value = 8080
        encrypted = false
        description = "Application port for this web server"
        tags = ["application", "port"]
      }
      
      # Feature flags (object)
      feature_flags = {
        value = {
          new_ui = true
          beta_features = false
          debug_mode = false
        }
        encrypted = false
        description = "Feature flags for this web server"
        tags = ["features", "configuration"]
      }
    }
  }
  
  machine {
    name = "db-server-01"
    description = "Database server with encrypted variables"
    hostname = "192.168.1.100"
    port = 2222
    user = "dbadmin"
    
    authentication {
      method = "password"
      password = {
        value = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
        encrypted = true
        encryption_metadata = {
          recipients = ["age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"]
          encrypted_at = "2024-01-01T12:00:00Z"
          method = "age"
        }
      }
    }
    
    connection {
      timeout = 60
      retry_attempts = 5
      retry_delay = 10
    }
    
    tags = ["database", "production", "primary"]
    environment = "production"
    role = "database"
    location = "us-east-1"
    
    # Machine-specific variables
    variables {
      # Database password (encrypted)
      db_password = {
        value = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
        encrypted = true
        encryption_metadata = {
          recipients = ["age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"]
          encrypted_at = "2024-01-01T12:00:00Z"
          method = "age"
        }
        description = "Database password for this server"
        sensitive = true
        tags = ["database", "password"]
      }
      
      # Database port (plain text)
      db_port = {
        value = 5432
        encrypted = false
        description = "Database port"
        tags = ["database", "port"]
      }
    }
  }
  
  machine {
    name = "dev-server-01"
    description = "Development server with plain text variables"
    hostname = "dev01.example.com"
    port = 22
    user = "developer"
    
    authentication {
      method = "ssh_key"
      ssh_key_path = "~/.ssh/dev_key"
    }
    
    connection {
      timeout = 20
      retry_attempts = 2
      retry_delay = 3
    }
    
    tags = ["development", "dev"]
    environment = "development"
    role = "development"
    location = "local"
    
    # Machine-specific variables (plain text for development)
    variables {
      # Development database URL (plain text)
      dev_db_url = {
        value = "postgresql://dev:devpass@localhost:5432/devdb"
        encrypted = false
        description = "Development database URL"
        tags = ["database", "development"]
      }
      
      # Debug mode (boolean)
      debug_mode = {
        value = true
        encrypted = false
        description = "Enable debug mode for development"
        tags = ["development", "debug"]
      }
    }
  }
  
  # Group-level variables (applied to all machines in specific groups)
  group_variables {
    # Production group variables
    production = {
      # Production environment variables
      env = {
        value = "production"
        encrypted = false
        description = "Environment for production machines"
        tags = ["environment", "production"]
      }
      
      # Production API key (encrypted)
      api_key = {
        value = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
        encrypted = true
        encryption_metadata = {
          recipients = ["age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"]
          encrypted_at = "2024-01-01T12:00:00Z"
          method = "age"
        }
        description = "API key for production services"
        sensitive = true
        tags = ["api", "production"]
      }
      
      # Production monitoring settings
      monitoring = {
        value = {
          enabled = true
          interval = 60
          alerting = true
        }
        encrypted = false
        description = "Monitoring configuration for production"
        tags = ["monitoring", "production"]
      }
    }
    
    # Database group variables
    database = {
      # Database backup settings
      backup_retention = {
        value = 30
        encrypted = false
        description = "Number of days to retain database backups"
        tags = ["database", "backup"]
      }
      
      # Database encryption key (encrypted)
      encryption_key = {
        value = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
        encrypted = true
        encryption_metadata = {
          recipients = ["age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"]
          encrypted_at = "2024-01-01T12:00:00Z"
          method = "age"
        }
        description = "Database encryption key"
        sensitive = true
        tags = ["database", "encryption"]
      }
    }
    
    # Development group variables
    development = {
      # Development environment variables
      env = {
        value = "development"
        encrypted = false
        description = "Environment for development machines"
        tags = ["environment", "development"]
      }
      
      # Development logging level
      log_level = {
        value = "debug"
        encrypted = false
        description = "Logging level for development"
        tags = ["logging", "development"]
      }
      
      # Development feature flags
      features = {
        value = {
          experimental = true
          beta_features = true
          debug_mode = true
        }
        encrypted = false
        description = "Feature flags for development"
        tags = ["features", "development"]
      }
    }
  }
}
