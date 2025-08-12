machines {
  # Production Environment - Web Servers
  machine "prod-web-01" {
    host = "10.0.1.10"
    user = "admin"
    port = 22
    key_file = "~/.ssh/prod_web_key"
    passphrase = "secure-production-passphrase"
    
    tags = ["web", "production", "load-balanced"]
    groups = ["web-servers", "production-servers"]
    roles = ["web-server", "nginx", "ssl-terminator"]
    
    resources {
      cpu_cores = 8
      memory_gb = 32
      disk_gb = 500
      network_mbps = 10000
    }
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      rack = "A-01"
      location = "San Francisco"
      owner = "web-team"
      department = "Engineering"
      cost_center = "IT-001"
      maintenance_window = "Sunday 2-4 AM PST"
      backup_schedule = "daily"
      monitoring = "prometheus"
      alerting = "pagerduty"
      sla = "99.9%"
    }
  }
  
  machine "prod-web-02" {
    host = "10.0.1.11"
    user = "admin"
    port = 22
    key_file = "~/.ssh/prod_web_key"
    passphrase = "secure-production-passphrase"
    
    tags = ["web", "production", "load-balanced"]
    groups = ["web-servers", "production-servers"]
    roles = ["web-server", "nginx", "ssl-terminator"]
    
    resources {
      cpu_cores = 8
      memory_gb = 32
      disk_gb = 500
      network_mbps = 10000
    }
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      rack = "A-02"
      location = "San Francisco"
      owner = "web-team"
      department = "Engineering"
      cost_center = "IT-001"
      maintenance_window = "Sunday 2-4 AM PST"
      backup_schedule = "daily"
      monitoring = "prometheus"
      alerting = "pagerduty"
      sla = "99.9%"
    }
  }
  
  # Production Environment - Database Servers
  machine "prod-db-primary" {
    host = "10.0.1.20"
    user = "dbadmin"
    port = 22
    key_file = "~/.ssh/prod_db_key"
    
    tags = ["database", "production", "primary"]
    groups = ["database-servers", "production-servers"]
    roles = ["database-server", "postgresql", "primary"]
    
    resources {
      cpu_cores = 16
      memory_gb = 64
      disk_gb = 2000
      network_mbps = 10000
    }
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      rack = "B-01"
      location = "San Francisco"
      owner = "db-team"
      department = "Engineering"
      cost_center = "IT-002"
      maintenance_window = "Sunday 2-4 AM PST"
      backup_schedule = "hourly"
      monitoring = "prometheus"
      alerting = "pagerduty"
      sla = "99.99%"
    }
  }
  
  machine "prod-db-replica" {
    host = "10.0.1.21"
    user = "dbadmin"
    port = 22
    key_file = "~/.ssh/prod_db_key"
    
    tags = ["database", "production", "replica"]
    groups = ["database-servers", "production-servers"]
    roles = ["database-server", "postgresql", "replica"]
    
    resources {
      cpu_cores = 16
      memory_gb = 64
      disk_gb = 2000
      network_mbps = 10000
    }
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      rack = "B-02"
      location = "San Francisco"
      owner = "db-team"
      department = "Engineering"
      cost_center = "IT-002"
      maintenance_window = "Sunday 2-4 AM PST"
      backup_schedule = "hourly"
      monitoring = "prometheus"
      alerting = "pagerduty"
      sla = "99.99%"
    }
  }
  
  # Staging Environment - Web Servers
  machine "staging-web-01" {
    host = "10.0.2.10"
    user = "admin"
    port = 22
    key_file = "~/.ssh/staging_key"
    
    tags = ["web", "staging"]
    groups = ["web-servers", "staging-servers"]
    roles = ["web-server", "nginx"]
    
    resources {
      cpu_cores = 4
      memory_gb = 16
      disk_gb = 200
      network_mbps = 1000
    }
    
    metadata {
      environment = "staging"
      datacenter = "us-west-1"
      owner = "web-team"
      department = "Engineering"
      cost_center = "IT-003"
      maintenance_window = "Saturday 2-4 AM PST"
      backup_schedule = "daily"
      monitoring = "prometheus"
      alerting = "slack"
      sla = "99.5%"
    }
  }
  
  # Staging Environment - Database Server
  machine "staging-db-01" {
    host = "10.0.2.20"
    user = "dbadmin"
    port = 22
    key_file = "~/.ssh/staging_key"
    
    tags = ["database", "staging"]
    groups = ["database-servers", "staging-servers"]
    roles = ["database-server", "postgresql"]
    
    resources {
      cpu_cores = 8
      memory_gb = 32
      disk_gb = 500
      network_mbps = 1000
    }
    
    metadata {
      environment = "staging"
      datacenter = "us-west-1"
      owner = "db-team"
      department = "Engineering"
      cost_center = "IT-003"
      maintenance_window = "Saturday 2-4 AM PST"
      backup_schedule = "daily"
      monitoring = "prometheus"
      alerting = "slack"
      sla = "99.5%"
    }
  }
  
  # Development Environment - All-in-One Server
  machine "dev-server-01" {
    host = "10.0.3.10"
    user = "developer"
    port = 22
    key_file = "~/.ssh/dev_key"
    
    tags = ["web", "database", "development"]
    groups = ["development-servers"]
    roles = ["web-server", "database-server", "nginx", "postgresql"]
    
    metadata {
      environment = "development"
      datacenter = "us-west-1"
      owner = "developer"
      department = "Engineering"
      purpose = "web and database development and testing"
      monitoring = "basic"
      alerting = "email"
    }
  }
}
