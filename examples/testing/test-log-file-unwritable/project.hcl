project "test-log-file-unwritable" {
  description = "test-log-file-unwritable project"
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
  
  # Logging configuration - unwritable path
  logging {
    level = "info"
    format = "json"
    output = "/root/unwritable.log"
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
    project = "test-log-file-unwritable"
  }
}