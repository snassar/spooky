# Invalid Project Configuration
# This file should fail schema validations

metadata {
  schema_version = "invalid-version"  # Invalid ScalVer format
  schema_type = "project"
  schema_name = "Test Project Configuration"
  last_updated = "2024-01-01"
  compatibility = ["0.20250809.0"]
  description = "A test project configuration that should fail validations"
}

project {
  name = "123-invalid-name"  # Invalid: starts with number
  description = "A test automation project for schema validation"
  version = "1.0.0"
  author = "Test User"
  email = "invalid-email"  # Invalid email format
  url = "not-a-valid-url"  # Invalid URL format
  
  run {
    default_timeout = 9999  # Invalid: too high
    max_parallel = 999  # Invalid: too high
    dry_run_default = false
    validate_before_run = true
  }
  
  facts {
    enabled = true
    timeout = 9999  # Invalid: too high
    parallel = 999  # Invalid: too high
    storage = "invalid-storage"  # Invalid: not in enum
  }
  
  logging {
    level = "invalid-level"  # Invalid: not in enum
    format = "invalid-format"  # Invalid: not in enum
    output = "invalid-output"  # Invalid: not in enum
    file_path = "logs/project.log"
  }
}
