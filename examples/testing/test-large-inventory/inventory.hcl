# Large inventory for test-large-inventory project
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
  
  # Generate many machines for performance testing
  machine "web-server-1" {
    host     = "192.168.1.101"
    port     = 22
    user     = "debian"
    password = "your-password"
    tags = {
      environment = "production"
      role = "web"
      region = "us-east"
    }
  }
  
  machine "web-server-2" {
    host     = "192.168.1.102"
    port     = 22
    user     = "debian"
    password = "your-password"
    tags = {
      environment = "production"
      role = "web"
      region = "us-east"
    }
  }
  
  machine "db-server-1" {
    host     = "192.168.1.103"
    port     = 22
    user     = "debian"
    password = "your-password"
    tags = {
      environment = "production"
      role = "database"
      region = "us-east"
    }
  }
  
  machine "db-server-2" {
    host     = "192.168.1.104"
    port     = 22
    user     = "debian"
    password = "your-password"
    tags = {
      environment = "production"
      role = "database"
      region = "us-east"
    }
  }
  
  machine "cache-server-1" {
    host     = "192.168.1.105"
    port     = 22
    user     = "debian"
    password = "your-password"
    tags = {
      environment = "production"
      role = "cache"
      region = "us-east"
    }
  }
  
  machine "cache-server-2" {
    host     = "192.168.1.106"
    port     = 22
    user     = "debian"
    password = "your-password"
    tags = {
      environment = "production"
      role = "cache"
      region = "us-east"
    }
  }
  
  machine "load-balancer-1" {
    host     = "192.168.1.107"
    port     = 22
    user     = "debian"
    password = "your-password"
    tags = {
      environment = "production"
      role = "load-balancer"
      region = "us-east"
    }
  }
  
  machine "load-balancer-2" {
    host     = "192.168.1.108"
    port     = 22
    user     = "debian"
    password = "your-password"
    tags = {
      environment = "production"
      role = "load-balancer"
      region = "us-east"
    }
  }
  
  machine "monitoring-1" {
    host     = "192.168.1.109"
    port     = 22
    user     = "debian"
    password = "your-password"
    tags = {
      environment = "production"
      role = "monitoring"
      region = "us-east"
    }
  }
  
  machine "monitoring-2" {
    host     = "192.168.1.110"
    port     = 22
    user     = "debian"
    password = "your-password"
    tags = {
      environment = "production"
      role = "monitoring"
      region = "us-east"
    }
  }
}