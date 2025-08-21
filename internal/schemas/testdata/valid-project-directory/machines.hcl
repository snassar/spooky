# Valid Machines Configuration
# This file should pass all schema validations

metadata {
  schema_version = "0.20250809.1"
  schema_type = "machines"
  schema_name = "Test Machines Configuration"
  last_updated = "2024-01-01"
  compatibility = ["0.20250809.0", "0.20250809.1"]
  description = "A test machines configuration that should pass all validations"
}

machines_inventory {
  machines = [
    {
      name = "web-server-01"
      description = "Primary web server"
      hostname = "web01.example.com"
      port = 22
      user = "admin"
      
      authentication {
        method = "ssh_key"
        ssh_key_path = "~/.ssh/id_rsa"
        ssh_key_passphrase = "encrypted:age1..."
      }
      
      connection {
        timeout = 30
        retry_attempts = 3
        retry_delay = 5
      }
      
      tags = ["web", "production", "primary"]
    },
    {
      name = "db-server-01"
      description = "Database server"
      hostname = "192.168.1.100"
      port = 2222
      user = "dbadmin"
      
      authentication {
        method = "password"
        password = "encrypted:age1..."
      }
      
      connection {
        timeout = 60
        retry_attempts = 5
        retry_delay = 10
      }
      
      tags = ["database", "production", "primary"]
    }
  ]
}
