project "test-network-timeouts" {
  description = "test-network-timeouts project"
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
  
  # SSH configuration with very short timeouts
  ssh {
    default_user = "debian"
    default_port = 22
    connection_timeout = 1  # Very short connection timeout
    command_timeout = 2     # Very short command timeout
    retry_attempts = 1      # Minimal retries
  }
  
  # Tags for project-wide targeting
  tags = {
    project = "test-network-timeouts"
  }
}