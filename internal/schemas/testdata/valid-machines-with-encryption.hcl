# Valid Machines Configuration with Encrypted Authentication
# This file demonstrates how to use encrypted SSH credentials in machines.hcl

metadata {
  version = 1
  description = "A test machines configuration demonstrating encrypted SSH credentials"
}

machines {
  machine {
    name = "web-server-01"
    description = "Primary web server with encrypted SSH key passphrase"
    hostname = "web01.example.com"
    port = 22
    user = "admin"
    
    authentication {
      method = "ssh_key"
      ssh_key_path = "~/.ssh/id_rsa"
      
      # Encrypted SSH key passphrase
      ssh_key_passphrase = {
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
      timeout = 30
      retry_attempts = 3
      retry_delay = 5
    }
    
    tags = ["web", "production", "primary"]
    environment = "production"
    role = "web-server"
    location = "us-east-1"
  }
  
  machine {
    name = "db-server-01"
    description = "Database server with encrypted password"
    hostname = "192.168.1.100"
    port = 2222
    user = "dbadmin"
    
    authentication {
      method = "password"
      
      # Encrypted SSH password
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
  }
  
  machine {
    name = "cert-server-01"
    description = "Server using SSH certificate with encrypted passphrases"
    hostname = "cert01.example.com"
    port = 22
    user = "certuser"
    
    authentication {
      method = "certificate"
      certificate_path = "~/.ssh/cert.pub"
      certificate_key_path = "~/.ssh/cert_key"
      
      # Encrypted certificate passphrase (protects the certificate file)
      certificate_passphrase = {
        value = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
        encrypted = true
        encryption_metadata = {
          recipients = ["age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"]
          encrypted_at = "2024-01-01T12:00:00Z"
          method = "age"
        }
      }
      
      # Encrypted certificate private key passphrase (protects the private key)
      certificate_key_passphrase = {
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
      timeout = 45
      retry_attempts = 4
      retry_delay = 8
    }
    
    tags = ["certificate", "staging", "secondary"]
    environment = "staging"
    role = "application"
    location = "us-west-2"
  }
  
  machine {
    name = "dev-server-01"
    description = "Development server with plain text passphrase (will be encrypted by spooky project encrypt)"
    hostname = "dev01.example.com"
    port = 22
    user = "developer"
    
    authentication {
      method = "ssh_key"
      ssh_key_path = "~/.ssh/dev_key"
      
      # Plain text passphrase (will be encrypted by spooky project encrypt)
      ssh_key_passphrase = {
        value = "my-secret-passphrase"
        encrypted = false
      }
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
  }
}
