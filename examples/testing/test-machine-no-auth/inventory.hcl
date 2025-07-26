# Inventory for test-valid-project project
# Add your machine definitions here

inventory {
  machine "example-server" {
    host     = "192.168.1.100"
    port     = 22
    user     = "debian"
    password = "your-password"
    tags = {
      environment = "development"
      role = "web"
    }
  }
}

# Inventory with machine missing both password and key_file

inventory {
  machine "no-auth-server" {
    host     = "192.168.1.250"
    port     = 22
    user     = "debian"
    # No password or key_file
    tags = {
      environment = "test"
      role = "web"
    }
  }
}