project "test-malformed-hcl" {
  description = "test-malformed-hcl project"
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
    # Missing closing brace for storage block
  
  # Logging configuration
  logging {
    level = "info"
    format = "json"
    output = "logs/spooky.log"
    # Missing closing brace for logging block
  
  # SSH configuration
  ssh {
    default_user = "debian"
    default_port = 22
    connection_timeout = 30
    command_timeout = 300
    retry_attempts = 3
    # Missing closing brace for ssh block
  
  # Tags for project-wide targeting
  tags = {
    project = "test-malformed-hcl"
    # Missing closing brace for tags block
  
  # Missing closing brace for project block