# Inventory for parser_wrapper project
# Add your machine definitions here

inventory {
  machine "web-server-1" {
    host = "192.168.1.100"
    port = 22
    user = "debian"
    password = "test-password"
    tags = {
      environment = "testing"
      role = "web"
      region = "us-west"
    }
  }
  
  machine "web-server-2" {
    host = "192.168.1.101"
    port = 22
    user = "debian"
    key_file = "/path/to/key.pem"
    tags = {
      environment = "testing"
      role = "web"
      region = "us-west"
    }
  }
  
  machine "db-server" {
    host = "192.168.1.102"
    port = 22
    user = "postgres"
    password = "db-password"
    tags = {
      environment = "testing"
      role = "database"
      region = "us-west"
    }
  }
  
  machine "app-server" {
    host = "192.168.1.103"
    port = 22
    user = "appuser"
    password = "app-password"
    tags = {
      environment = "testing"
      role = "application"
      region = "us-east"
    }
  }
}