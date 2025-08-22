# Spooky Project Configuration Schema
# Comprehensive schema for project.hcl files with project metadata and structure validation

# Schema metadata
metadata {
  version = "1"
  description = "Project configuration schema for spooky projects"
}

# Project block structure
project {
  # Project metadata
  name = {
    type = "string"
    required = true
    pattern = "^[a-zA-Z][a-zA-Z0-9._-]*$"
    min_length = 1
    max_length = 128
    description = "Project name (used for identification and isolation)"
  }
  
  description = {
    type = "string"
    required = false
    max_length = 1024
    description = "Project description"
  }
  

  
  run_max_parallel = {
    type = "integer"
    required = false
    min = 1
    max = 100
    default = 10
    description = "Maximum parallel action executions for this project"
  }
  
  run_dry_run_default = {
    type = "boolean"
    required = false
    default = false
    description = "Default dry-run mode for actions"
  }
  
  run_validate_before_run = {
    type = "boolean"
    required = false
    default = true
    description = "Validate project configuration before running"
  }
  
  run_backup_before_changes = {
    type = "boolean"
    required = false
    default = false
    description = "Create backups before making changes"
  }
  
  # Facts collection configuration
  facts_timeout = {
    type = "integer"
    required = false
    min = 1
    max = 3600
    default = 30
    description = "Timeout for facts collection in seconds"
  }
  
  facts_parallel_collection = {
    type = "integer"
    required = false
    min = 1
    max = 100
    default = 10
    description = "Number of parallel facts collection workers for this project"
  }
  
  facts_retry_attempts = {
    type = "integer"
    required = false
    min = 0
    max = 10
    default = 3
    description = "Number of retry attempts for failed facts collection"
  }
  
  facts_retry_delay = {
    type = "integer"
    required = false
    min = 1
    max = 60
    default = 5
    description = "Delay between retry attempts in seconds"
  }
} 