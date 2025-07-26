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

# Inventory with invalid port values

inventory {
  machine "port-zero" {
    host     = "192.168.1.200"
    port     = 0
    user     = "debian"
    password = "your-password"
    tags = {
      environment = "test"
      role = "web"
    }
  }
  machine "port-too-high" {
    host     = "192.168.1.201"
    port     = 70000
    user     = "debian"
    password = "your-password"
    tags = {
      environment = "test"
      role = "web"
    }
  }
}