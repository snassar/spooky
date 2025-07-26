# Inventory for test-duplicate-machines project
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
  
  # Duplicate machine name - should cause validation error
  machine "example-server" {
    host     = "192.168.1.101"
    port     = 22
    user     = "debian"
    password = "your-password"
    tags = {
      environment = "production"
      role = "db"
    }
  }
}