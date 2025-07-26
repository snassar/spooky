# Inventory with invalid and missing tags

inventory {
  machine "non-string-tag-server" {
    host     = "192.168.1.251"
    port     = 22
    user     = "debian"
    password = "your-password"
    tags = {
      environment = 12345  # Non-string value
      role = "web"
    }
  }
  machine "missing-tag-server" {
    host     = "192.168.1.252"
    port     = 22
    user     = "debian"
    password = "your-password"
    tags = {
      # environment tag missing
      role = "web"
    }
  }
}