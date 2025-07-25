project "template-testing" {
  description = "Test project for template functionality"
  version = "1.0.0"
  environment = "testing"
  
  inventory_file = "inventory.hcl"
  actions_file = "actions.hcl"
  
  default_timeout = 300
  default_parallel = true
  
  storage {
    type = "json"
    path = ".facts.db"
  }
  
  logging {
    level = "info"
    format = "text"
  }
  
  ssh {
    default_user = "testuser"
    default_port = 22
    connection_timeout = 30
    command_timeout = 300
    retry_attempts = 3
  }
  
  tags = {
    environment = "testing"
    purpose = "unit-testing"
  }
} 