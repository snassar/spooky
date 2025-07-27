project "test-type-mismatches" {
  description = "test-type-mismatches project"
  version = "1.0.0"
  environment = "development"
  
  # File references
  inventory_file = "inventory.hcl"
  actions_file = "actions.hcl"
  
  # Project settings with type mismatches
  default_timeout = "not_a_number"  # Should be int, but is string
  default_parallel = "not_a_boolean"  # Should be bool, but is string
  
  # Storage configuration
  storage {
    type = "badgerdb"
    path = "test.db"
  }
  
  # Logging configuration
  logging {
    level = "info"
    format = "json"
    output = "logs/spooky.log"
  }
  
  # SSH configuration with type mismatches
  ssh {
    default_user = "debian"
    default_port = "not_a_port"  # Should be int, but is string
    connection_timeout = "not_a_timeout"  # Should be int, but is string
    command_timeout = "not_a_timeout"  # Should be int, but is string
    retry_attempts = "not_a_number"  # Should be int, but is string
  }
  
  # Tags for project-wide targeting
  tags = {
    project = "test-type-mismatches"
  }
}