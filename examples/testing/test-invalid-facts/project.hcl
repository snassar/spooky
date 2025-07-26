project "test-invalid-facts" {
  description = "test-invalid-facts project"
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
    path = ".facts.db"
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
    project = "test-invalid-facts"
  }
}