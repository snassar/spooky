inventory {
  machine "web-server" {
    host = "192.168.1.10"
    port = 22
    user = "admin"
    password = "secret123"
    tags = {
      environment = "production"
      role = "web"
      team = "frontend"
    }
  }

  machine "db-server" {
    host = "192.168.1.20"
    port = 22
    user = "dbadmin"
    key_file = "/home/user/.ssh/db_key"
    tags = {
      environment = "production"
      role = "database"
      team = "backend"
    }
  }

  machine "app-server" {
    host = "192.168.1.30"
    port = 22
    user = "appuser"
    password = "apppass"
    key_file = "/home/user/.ssh/app_key"
    tags = {
      environment = "staging"
      role = "application"
      team = "backend"
    }
  }

  machine "monitoring-server" {
    host = "192.168.1.40"
    port = 22
    user = "monitor"
    tags = {
      environment = "production"
      role = "monitoring"
      team = "ops"
    }
  }
}