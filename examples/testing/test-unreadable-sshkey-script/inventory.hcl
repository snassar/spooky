# Inventory with unreadable SSH key

inventory {
  machine "example-server" {
    host     = "192.168.1.100"
    port     = 22
    user     = "debian"
    ssh_key  = "unreadable.key"
    tags = {
      environment = "development"
      role = "web"
    }
  }
}