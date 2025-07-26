# Inventory for test-invalid-port-user project
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
  
  # Invalid port type - should cause validation error
  machine "invalid-port-server" {
    host     = "192.168.1.101"
    port     = "22"  # String instead of number
    user     = "debian"
    password = "your-password"
    tags = {
      environment = "development"
      role = "web"
    }
  }
  
  # Missing user - should cause validation error
  machine "missing-user-server" {
    host     = "192.168.1.102"
    port     = 22
    # user field missing
    password = "your-password"
    tags = {
      environment = "development"
      role = "web"
    }
  }
}