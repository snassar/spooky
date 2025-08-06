# Variables HCL Schema
# Schema for variables stored in HCL format
# This schema composes the base variables structure with HCL-specific storage rules

# Variables in HCL format
variables {
  # Include the base variables structure
  include = "variables-structure"
  
  # Variables specific metadata
  scope = "project"
  storage_location = "<project>/variables.hcl or <project>/variables/*.hcl"
  description = "Project variables stored in HCL format"
  
  # Variables specific validation
  variables_validation = {
    # Variables must be in a variables block
    variables_block_required = {
      rule = "required"
      field = "variables"
      message = "Variables must be defined within a 'variables' block"
    }
    
    # Variable names must be unique within the file
    unique_names = {
      rule = "unique"
      field = "variable.name"
      message = "Variable names must be unique within the file"
    }
    
    # Variable types must be valid
    valid_types = {
      rule = "enum"
      field = "variable.type"
      values = ["string", "number", "float", "bool", "list", "map", "object", "duration", "ip", "cidr", "path", "file", "secret"]
      message = "Variable type must be one of: string, number, float, bool, list, map, object, duration, ip, cidr, path, file, secret"
    }
    
    # Required variables must have defaults or be provided
    required_validation = {
      rule = "conditional"
      condition = "required_variables_have_defaults_or_environment"
      message = "Required variables must have default values or be provided via environment variables"
    }
    
    # No circular dependencies
    no_circular_deps = {
      rule = "acyclic"
      field = "variable.dependencies"
      message = "Variables cannot have circular dependencies"
    }
  }
  
  # Variables specific constraints
  variables_constraints = {
    # Variables are project-scoped by default
    project_scope = {
      type = "string"
      value = "project"
      description = "Variables are project-scoped by default"
    }
    
    # Variables support interpolation
    interpolation = {
      type = "boolean"
      value = true
      description = "Variables support HCL interpolation"
    }
    
    # Variables support references
    references = {
      type = "boolean"
      value = true
      description = "Variables support references to other variables"
    }
    
    # Variables support expressions
    expressions = {
      type = "boolean"
      value = true
      description = "Variables support HCL expressions in default values"
    }
    
    # Variables support multi-line strings
    multiline_strings = {
      type = "boolean"
      value = true
      description = "Variables support multi-line string values using heredoc syntax"
    }
    
    # Variables support comments
    comments = {
      type = "boolean"
      value = true
      description = "Variables support HCL comments"
    }
  }
  
  # HCL-specific features
  hcl_features = {
    # HCL block syntax
    block_syntax = {
      type = "string"
      value = "variable \"name\" { ... }"
      description = "Variables use HCL block syntax"
    }
    
    # HCL attribute syntax
    attribute_syntax = {
      type = "string"
      value = "attribute = value"
      description = "Variable properties use HCL attribute syntax"
    }
    
    # HCL interpolation
    interpolation_syntax = {
      type = "string"
      value = "${var.other_variable}"
      description = "Variables support HCL interpolation syntax"
    }
    
    # HCL expressions
    expression_syntax = {
      type = "string"
      value = "default = var.count + 1"
      description = "Variables support HCL expressions in default values"
    }
    
    # HCL heredoc
    heredoc_syntax = {
      type = "string"
      value = "default = <<-EOT\nmulti-line\ncontent\nEOT"
      description = "Variables support HCL heredoc syntax for multi-line strings"
    }
  }
} 