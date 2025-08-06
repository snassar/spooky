# Variables Structure Schema
# Common variable structure definitions for all storage formats
# This file defines the structure of variables that can be stored
# Used by HCL and JSON storage schemas for variables

# Common variables structure
variables_structure {
  # Main variables file (variables.hcl in project root)
  main_file = {
    type = "string"
    required = false
    pattern = "^variables\\.hcl$"
    description = "Main variables file in project root"
  }
  
  # Variables directory (variables/ in project root)
  variables_directory = {
    type = "string"
    required = false
    pattern = "^variables/$"
    description = "Variables directory containing organized variable files"
  }
  
  # Variable definition structure
  variable_blocks = {
    type = "list(object)"
    required = true
    min_items = 0
    
    properties = {
      name = {
        type = "string"
        required = true
        description = "Variable name"
        pattern = "^[a-z][a-z0-9_]*$"
        min_length = 1
        max_length = 100
      }
      
      type = {
        type = "string"
        required = true
        description = "Variable type"
        enum = ["string", "number", "float", "bool", "list", "map", "object", "duration", "ip", "cidr", "path", "file", "secret"]
      }
      
      description = {
        type = "string"
        required = false
        description = "Variable description"
        max_length = 500
      }
      
      default = {
        type = "any"
        required = false
        description = "Default value for the variable"
      }
      
      required = {
        type = "bool"
        required = false
        default = false
        description = "Whether the variable is required"
      }
      
      sensitive = {
        type = "bool"
        required = false
        default = false
        description = "Whether the variable contains sensitive data"
      }
      
      encrypted = {
        type = "bool"
        required = false
        default = false
        description = "Whether the variable should be encrypted at rest"
      }
      
      scope = {
        type = "string"
        required = false
        enum = ["project", "global", "inherited"]
        default = "project"
        description = "Variable scope (project-only, global, or inherited)"
      }
      
      dependencies = {
        type = "array"
        required = false
        description = "List of other variables this variable depends on"
        items = {
          type = "string"
          pattern = "^[a-z][a-z0-9_]*$"
        }
      }
      
      validation = {
        type = "object"
        required = false
        max_occurrences = 1
        
        properties = {
          condition = {
            type = "string"
            required = true
            description = "Validation condition using HCL syntax"
            max_length = 1000
          }
          
          error_message = {
            type = "string"
            required = true
            description = "Error message for validation failure"
            max_length = 500
          }
          
          warning_message = {
            type = "string"
            required = false
            description = "Warning message for validation warning"
            max_length = 500
          }
        }
      }
      
      constraints = {
        type = "object"
        required = false
        description = "Type-specific constraints"
        
        properties = {
          # String constraints
          min_length = {
            type = "integer"
            required = false
            min = 0
            description = "Minimum string length"
          }
          
          max_length = {
            type = "integer"
            required = false
            min = 1
            description = "Maximum string length"
          }
          
          pattern = {
            type = "string"
            required = false
            description = "Regular expression pattern for string validation"
          }
          
          # Numeric constraints
          min_value = {
            type = "number"
            required = false
            description = "Minimum numeric value"
          }
          
          max_value = {
            type = "number"
            required = false
            description = "Maximum numeric value"
          }
          
          # List constraints
          min_items = {
            type = "integer"
            required = false
            min = 0
            description = "Minimum number of list items"
          }
          
          max_items = {
            type = "integer"
            required = false
            min = 1
            description = "Maximum number of list items"
          }
          
          # File constraints
          file_exists = {
            type = "bool"
            required = false
            default = false
            description = "Whether the file must exist"
          }
          
          file_readable = {
            type = "bool"
            required = false
            default = false
            description = "Whether the file must be readable"
          }
          
          file_size_max = {
            type = "string"
            required = false
            pattern = "^\\d+[KMGT]?B$"
            description = "Maximum file size (e.g., '10MB', '1GB')"
          }
          
          # Path constraints
          path_exists = {
            type = "bool"
            required = false
            default = false
            description = "Whether the path must exist"
          }
          
          path_absolute = {
            type = "bool"
            required = false
            default = false
            description = "Whether the path must be absolute"
          }
          
          path_relative = {
            type = "bool"
            required = false
            default = false
            description = "Whether the path must be relative"
          }
        }
      }
      
      metadata = {
        type = "object"
        required = false
        description = "Additional variable metadata"
        additional_properties = "string"
      }
    }
  }
  
  # Variable resolution and interpolation
  resolution = {
    type = "object"
    required = false
    description = "Variable resolution configuration"
    
    properties = {
      allow_self_reference = {
        type = "bool"
        required = false
        default = false
        description = "Allow variables to reference themselves"
      }
      
      allow_circular_deps = {
        type = "bool"
        required = false
        default = false
        description = "Allow circular dependencies between variables"
      }
      
      max_resolution_depth = {
        type = "integer"
        required = false
        min = 1
        max = 100
        default = 10
        description = "Maximum depth for variable resolution"
      }
      
      fail_on_missing = {
        type = "bool"
        required = false
        default = true
        description = "Fail if required variables are missing"
      }
      
      use_environment = {
        type = "bool"
        required = false
        default = true
        description = "Use environment variables as fallback"
      }
      
      environment_prefix = {
        type = "string"
        required = false
        default = "SPOOKY_"
        pattern = "^[A-Z_]*$"
        description = "Prefix for environment variable names"
      }
    }
  }
  
  # Security and encryption
  security = {
    type = "object"
    required = false
    description = "Security and encryption settings"
    
    properties = {
      encryption_enabled = {
        type = "bool"
        required = false
        default = false
        description = "Enable encryption for sensitive variables"
      }
      
      encryption_method = {
        type = "string"
        required = false
        enum = ["age", "gpg", "none"]
        default = "age"
        description = "Encryption method for sensitive variables"
      }
      
      key_file = {
        type = "string"
        required = false
        pattern = "^[a-zA-Z0-9/._-]+$"
        description = "Path to encryption key file"
      }
      
      key_id = {
        type = "string"
        required = false
        description = "Encryption key identifier"
      }
      
      mask_sensitive = {
        type = "bool"
        required = false
        default = true
        description = "Mask sensitive variable values in logs"
      }
      
      audit_trail = {
        type = "bool"
        required = false
        default = false
        description = "Enable audit trail for variable access"
      }
    }
  }
  
  # Validation rules
  validation = {
    # File location validation
    file_location = {
      rule = "path"
      pattern = "^(variables\\.hcl|variables/[^/]+\\.hcl)$"
      message = "Variable files must be variables.hcl in project root or .hcl files in variables/ directory"
    }
    
    # Variable name validation
    variable_name = {
      rule = "regex"
      pattern = "^[a-z][a-z0-9_]*$"
      message = "Variable names must be lowercase with underscores, starting with a letter"
    }
    
    # Type validation
    variable_type = {
      rule = "enum"
      allowed_values = ["string", "number", "float", "bool", "list", "map", "object", "duration", "ip", "cidr", "path", "file", "secret"]
      message = "Variable type must be one of: string, number, float, bool, list, map, object, duration, ip, cidr, path, file, secret"
    }
    
    # Required variable validation
    required_variable = {
      rule = "required"
      condition = "variable_is_required"
      message = "Required variables must have a default value or be provided via environment variable"
    }
    
    # Validation condition syntax
    validation_condition = {
      rule = "hcl"
      message = "Validation condition must be valid HCL syntax"
    }
    
    # No circular dependencies
    no_circular_deps = {
      rule = "acyclic"
      message = "Variables cannot have circular dependencies"
    }
    
    # Constraint validation
    constraint_consistency = {
      rule = "conditional"
      condition = "constraints_are_consistent_with_type"
      message = "Constraints must be consistent with variable type"
    }
    
    # Security validation
    sensitive_encryption = {
      rule = "conditional"
      condition = "sensitive_variables_are_encrypted"
      message = "Sensitive variables should be encrypted"
    }
    
    # Resolution validation
    resolution_depth = {
      rule = "range"
      min = 1
      max = 100
      message = "Resolution depth must be between 1 and 100"
    }
  }
} 