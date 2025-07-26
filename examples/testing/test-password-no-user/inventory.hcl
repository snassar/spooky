# Inventory for test-password-no-user project
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
  
  # Password but no user - should cause validation error
  machine "no-user-server" {
    host     = "192.168.1.101"
    port     = 22
    # user field missing
    password = "your-password"
    tags = {
      environment = "development"
      role = "web"
    }
  }
}