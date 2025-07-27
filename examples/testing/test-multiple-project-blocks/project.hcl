project "first" {
  description = "first project"
  version = "1.0.0"
  environment = "development"
  
  # File references
  inventory_file = "inventory.hcl"
  actions_file = "actions.hcl"
  
  # Project settings
  default_timeout = 300
  default_parallel = true
  
  # Storage configuration
  storage {
    type = "badgerdb"
    path = "test1.db"
  }
  
  # Logging configuration
  logging {
    level = "info"
    format = "json"
    output = "logs/spooky.log"
  }
  
  # SSH configuration
  ssh {
    default_user = "debian"
    default_port = 22
    connection_timeout = 30
    command_timeout = 300
    retry_attempts = 3
  }
  
  # Tags for project-wide targeting
  tags = {
    project = "first"
  }
}

project "second" {
  description = "second project"
  version = "2.0.0"
  environment = "production"
  
  # File references
  inventory_file = "inventory.hcl"
  actions_file = "actions.hcl"
  
  # Project settings
  default_timeout = 600
  default_parallel = false
  
  # Storage configuration
  storage {
    type = "badgerdb"
    path = "test2.db"
  }
  
  # Logging configuration
  logging {
    level = "warn"
    format = "text"
    output = "logs/spooky.log"
  }
  
  # SSH configuration
  ssh {
    default_user = "admin"
    default_port = 2222
    connection_timeout = 60
    command_timeout = 600
    retry_attempts = 5
  }
  
  # Tags for project-wide targeting
  tags = {
    project = "second"
  }
}