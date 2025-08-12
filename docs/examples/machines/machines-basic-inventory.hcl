machines {
  # Basic web server
  machine "web-server-01" {
    host = "192.168.1.10"
    user = "admin"
    port = 22
    key_file = "~/.ssh/id_rsa"
    
    tags = ["web", "production"]
    groups = ["web-servers"]
    roles = ["web-server", "nginx"]
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      owner = "web-team"
      department = "Engineering"
      cost_center = "IT-001"
      maintenance_window = "Sunday 2-4 AM PST"
      backup_schedule = "daily"
    }
  }
  
  # Database server
  machine "db-server-01" {
    host = "192.168.1.20"
    user = "dbadmin"
    port = 22
    key_file = "~/.ssh/db_key"
    
    tags = ["database", "production"]
    groups = ["database-servers"]
    roles = ["database-server", "postgresql"]
    
    resources {
      cpu_cores = 16
      memory_gb = 64
      disk_gb = 2000
      network_mbps = 10000
    }
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
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
  
  # Load balancer
  machine "lb-primary" {
    host = "192.168.1.100"
    user = "admin"
    port = 22
    key_file = "~/.ssh/lb_key"
    
    tags = ["load-balancer", "production", "primary"]
    groups = ["load-balancers", "production-servers"]
    roles = ["load-balancer", "haproxy", "ssl-terminator"]
    
    resources {
      cpu_cores = 4
      memory_gb = 16
      disk_gb = 200
      network_mbps = 10000
    }
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      rack = "A-00"
      location = "San Francisco"
      owner = "infrastructure-team"
      department = "Engineering"
      cost_center = "IT-000"
      maintenance_window = "Sunday 2-4 AM PST"
      backup_schedule = "daily"
      monitoring = "prometheus"
      alerting = "pagerduty"
      sla = "99.99%"
    }
  }
}
