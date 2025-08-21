# Valid Project Configuration
# This file should pass all schema validations

metadata {
  schema_version = "0.20250809.0"
  schema_type = "project"
  schema_name = "Test Project Configuration"
  last_updated = "2024-01-01"
  compatibility = ["0.20250809.0"]
  description = "A test project configuration that should pass all validations"
}

project {
  name = "test-automation-project"
  description = "A test automation project for schema validation"
  version = "1.0.0"
  author = "Test User"
  email = "test@example.com"
  url = "https://example.com/test-project"
  
  run {
    default_timeout = 300
    max_parallel = 10
    dry_run_default = false
    validate_before_run = true
  }
  
  facts {
    enabled = true
    timeout = 60
    parallel = 5
    storage = "local"
  }
  
  logging {
    level = "info"
    format = "text"
    output = "stdout"
    file_path = "logs/project.log"
  }
}
