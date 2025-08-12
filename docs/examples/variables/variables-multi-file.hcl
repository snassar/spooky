# This example demonstrates multi-file variable organization
# This would be the main variables.hcl file in the project root

variables {
  variable "project_name" {
    type = "string"
    description = "Project name"
    default = "my-project"
    scope = "project"
  }

  variable "project_version" {
    type = "string"
    description = "Project version"
    default = "1.0.0"
    scope = "project"
  }

  variable "environment" {
    type = "string"
    description = "Deployment environment"
    default = "development"
    scope = "project"
    
    validation {
      allowed_values = ["development", "staging", "production"]
    }
  }

  variable "region" {
    type = "string"
    description = "AWS region"
    default = "us-west-2"
    scope = "project"
  }
}
