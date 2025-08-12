# Spooky Project Directory Structure Schema
# Schema for validating project directory structure and content
# This schema defines the expected structure of a spooky project directory

# Schema metadata
metadata {
  schema_version = "0.20250809.0"
  schema_type = "project-directory"
  schema_name = "Spooky Project Directory Structure Schema"
  last_updated = "2024-01-01"
  compatibility = ["0.20250809.0"]
  description = "Schema for validating project directory structure and content - defines the expected structure of a spooky project directory"
  
  # ScalVer format: 0.YYYYMMDD.N
  # - 0: Development phase
  # - 20250809: Date (9 August 2025)
  # - 0: Patch version
  scalver_format = "0.20250809.0"
}

# Project directory structure validation
project_directory "project_root" {
  # Required files
  file "project.hcl" {
    type = "file"
    required = true
    description = "Main project configuration file"
    validate = "hcl_project_config"
    pattern = "project \"[a-zA-Z0-9_-]+\" {"
  }
  
  # Optional machine inventory - either file or directory
  file "machines.hcl" {
    type = "file"
    required = false
    description = "Machine inventory definitions (optional if machines/ directory exists)"
    validate = "hcl_machines_config"
    pattern = "machines {"
  }
  
  directory "machines" {
    type = "directory"
    required = false
    description = "Machine inventory files directory (optional if machines.hcl exists)"
    validate = "hcl_machines_files"
    pattern = ".*\\.hcl$"
  }
  
  # Optional actions - either file or directory
  file "actions.hcl" {
    type = "file"
    required = false
    description = "Main actions file (optional if actions/ directory exists)"
    validate = "hcl_actions_config"
    pattern = "actions {"
  }
  
  directory "actions" {
    type = "directory"
    required = false
    description = "Organized action files (optional if actions.hcl exists)"
    validate = "hcl_actions_files"
    pattern = ".*\\.hcl$"
  }
  
  # Optional variables - either file or directory
  file "variables.hcl" {
    type = "file"
    required = false
    description = "Main variables file (optional if variables/ directory exists)"
    validate = "hcl_variables_config"
    pattern = "variables {"
  }
  
  directory "variables" {
    type = "directory"
    required = false
    description = "Variables files directory (optional if variables.hcl exists)"
    validate = "hcl_variables_files"
    pattern = ".*\\.hcl$"
  }
  
  file "README.md" {
    type = "file"
    required = false
    description = "Project documentation"
    pattern = "# .*"
  }
  
  # Optional directories
  directory "templates" {
    type = "directory"
    required = false
    description = "Template files for dynamic content"
    validate = "directory_exists"
  }
  
  directory "files" {
    type = "directory"
    required = false
    description = "Static files to be deployed"
    validate = "directory_exists"
  }
  
  directory "logs" {
    type = "directory"
    required = false
    description = "Log files directory"
    validate = "directory_exists"
  }
  
  # Facts database directory
  directory "facts.db" {
    type = "directory"
    required = true
    description = "Facts database directory (initialized during project creation)"
    validate = "badgerdb_initialized"
  }
  
  # Cross-file validation rules
  validation_rules = [
    "machines_file_or_directory_exists",
    "actions_file_or_directory_exists", 
    "variables_file_or_directory_exists",
    "facts_database_initialized",
    "no_circular_references",
    "logging_file_output_requires_logs_directory",
    "logging_file_path_validation"
  ]
}