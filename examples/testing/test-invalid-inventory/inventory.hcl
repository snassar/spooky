# Inventory for test-invalid-inventory project
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
  
  # Invalid machine definition - missing required fields
  machine "invalid-server" {
    # Missing host field
    port     = 22
    user     = "debian"
    tags = {
      environment = "development"
    }
  }
  
  # Invalid syntax - missing closing brace
  machine "broken-server" {
    host = "192.168.1.101"
    port = 22
    user = "debian"
    # Missing closing brace