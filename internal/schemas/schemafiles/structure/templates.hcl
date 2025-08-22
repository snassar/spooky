# Templates Schema
# Schema for Go template configuration and enrichment
# Templates are enriched with machines, facts, and variables data

# Schema metadata
metadata {
  version = "1"
  description = "Template configuration schema for spooky Go templates"
}

# Template structure
templates {
  # Template name
  name = {
    type = "string"
    required = true
    pattern = "^[a-zA-Z0-9._-]+$"
    min_length = 1
    max_length = 128
    description = "Template name"
  }
  
  # Template source file
  source_path = {
    type = "string"
    required = true
    pattern = "^templates/.*\\.tmpl$"
    description = "Path to Go template file (.tmpl extension)"
  }
  
  # Template destination (optional)
  destination_path = {
    type = "string"
    required = false
    description = "Default destination path for rendered output"
  }
  
  # Template description
  description = {
    type = "string"
    required = false
    max_length = 500
    description = "Template description"
  }
  
  # Template tags
  tags = {
    type = "array"
    required = false
    max_items = 10
    description = "Template tags for categorization"
    items = {
      type = "string"
      pattern = "^[a-zA-Z0-9_-]+$"
      min_length = 1
      max_length = 32
    }
  }
  
  # Enable machines data
  enable_machines = {
    type = "boolean"
    required = false
    default = true
    description = "Enable machines data in template context"
  }
  
  # Enable facts data
  enable_facts = {
    type = "boolean"
    required = false
    default = true
    description = "Enable facts data in template context"
  }
  
  # Enable variables data
  enable_variables = {
    type = "boolean"
    required = false
    default = true
    description = "Enable variables data in template context"
  }
  
  # Enable environment variables
  enable_environment = {
    type = "boolean"
    required = false
    default = true
    description = "Enable environment variables in template context"
  }
  
  # Maximum template size
  max_size_mb = {
    type = "integer"
    required = false
    min = 1
    max = 10
    default = 1
    description = "Maximum template file size in MB"
  }
  
  # Maximum execution time
  max_execution_time_seconds = {
    type = "integer"
    required = false
    min = 1
    max = 300
    default = 30
    description = "Maximum template execution time in seconds"
  }
  
  # Maximum memory usage
  max_memory_mb = {
    type = "integer"
    required = false
    min = 1
    max = 100
    default = 10
    description = "Maximum memory usage in MB"
  }
  
  # Template version
  version = {
    type = "string"
    required = false
    pattern = "^[a-zA-Z0-9._-]+$"
    max_length = 32
    description = "Template version"
  }
  
  # Template author
  author = {
    type = "string"
    required = false
    max_length = 100
    description = "Template author"
  }
  
  # Created timestamp
  created_at = {
    type = "string"
    required = false
    pattern = "^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$"
    description = "Template creation timestamp (ISO 8601)"
  }
  
  # Updated timestamp
  updated_at = {
    type = "string"
    required = false
    pattern = "^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$"
    description = "Template last update timestamp (ISO 8601)"
  }
} 