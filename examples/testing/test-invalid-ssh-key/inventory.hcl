# Inventory for test-invalid-ssh-key project
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
  
  # Invalid SSH key path - non-existent file
  machine "invalid-key-server" {
    host     = "192.168.1.101"
    port     = 22
    user     = "debian"
    ssh_key  = "/path/to/non/existent/key"
    tags = {
      environment = "development"
      role = "web"
    }
  }
  
  # Invalid SSH key path - directory instead of file
  machine "directory-key-server" {
    host     = "192.168.1.102"
    port     = 22
    user     = "debian"
    ssh_key  = "/tmp"
    tags = {
      environment = "development"
      role = "web"
    }
  }
}