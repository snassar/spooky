# Inventory with potential circular references in tags

inventory {
  machine "server-a" {
    host     = "192.168.1.100"
    port     = 22
    user     = "debian"
    password = "your-password"
    tags = {
      environment = "development"
      role = "web"
      depends_on = "server-b"
    }
  }
  
  machine "server-b" {
    host     = "192.168.1.101"
    port     = 22
    user     = "debian"
    password = "your-password"
    tags = {
      environment = "development"
      role = "database"
      depends_on = "server-c"
    }
  }
  
  machine "server-c" {
    host     = "192.168.1.102"
    port     = 22
    user     = "debian"
    password = "your-password"
    tags = {
      environment = "development"
      role = "cache"
      depends_on = "server-a"  # Circular reference back to server-a
    }
  }
}