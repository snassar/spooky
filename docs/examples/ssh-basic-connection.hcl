# Basic SSH Connection Example
# This file demonstrates basic SSH connection configuration for spooky

machines {
  # Basic ED25519 key connection
  machine "web-server-01" {
    host = "192.168.1.10"
    user = "admin"
    port = 22
    
    key_file = "~/.ssh/id_ed25519"
    passphrase = "my-secure-passphrase"  # Optional
    
    tags = ["web", "production"]
    groups = ["web-servers"]
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      owner = "web-team"
      key_type = "ed25519"
    }
  }
  
  # Basic RSA 4096-bit key connection
  machine "db-server-01" {
    host = "192.168.1.20"
    user = "dbadmin"
    port = 22
    
    key_file = "~/.ssh/id_rsa_4096"
    # No passphrase = unencrypted key
    
    tags = ["database", "production"]
    groups = ["database-servers"]
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      owner = "db-team"
      key_type = "rsa-4096"
    }
  }
  
  # Password authentication (less secure, not recommended for production)
  machine "legacy-server" {
    host = "192.168.1.30"
    user = "admin"
    port = 22
    
    password = "my-secure-password"
    
    tags = ["legacy", "development"]
    groups = ["legacy-servers"]
    
    metadata {
      environment = "development"
      datacenter = "us-west-1"
      owner = "legacy-team"
      auth_method = "password"
    }
  }
  
  # Custom SSH port
  machine "custom-port-server" {
    host = "192.168.1.40"
    user = "admin"
    port = 2222  # Custom SSH port
    
    key_file = "~/.ssh/id_ed25519"
    
    tags = ["custom", "production"]
    groups = ["custom-servers"]
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      owner = "custom-team"
      key_type = "ed25519"
    }
  }
}
