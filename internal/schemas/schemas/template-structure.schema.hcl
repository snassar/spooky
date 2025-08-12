# Template Structure Schema
# Common template structure definitions for all template-related schemas
# This file defines the base structure for templates that can be used
# Used by context, functions, and metadata schemas for templates

# Schema metadata
metadata {
  schema_version = "0.20250809.0"
  schema_type = "template-structure"
  schema_name = "Template Structure Schema"
  last_updated = "2024-01-01"
  compatibility = ["0.20250809.0"]
  description = "Common template structure definitions for all template-related schemas - defines the base structure for templates that can be used"
  
  # ScalVer format: 0.YYYYMMDD.N
  # - 0: Development phase
  # - 20250809: Date (9 August 2025)
  # - 0: Patch version
  scalver_format = "0.20250809.0"
}

# Common template structure
template_structure {
  # Template identification
  template_id = {
    type = "string"
    required = true
    pattern = "^[a-zA-Z0-9._-]+$"
    description = "Unique template identifier"
  }

  # Template source location
  source_path = {
    type = "string"
    required = true
    pattern = "^templates/.*\\.tmpl$"
    description = "Path to template file relative to project root"
  }

  # Template destination
  destination_path = {
    type = "string"
    required = false
    description = "Default destination path for rendered output"
  }

  # Template type classification
  template_type = {
    type = "string"
    required = true
    enum = ["context", "functions", "metadata", "combined"]
    description = "Type of template schema"
  }

  # Template scope
  scope = {
    type = "string"
    required = true
    enum = ["project", "global", "machine", "system"]
    description = "Scope of template usage"
  }

  # Template security level
  security_level = {
    type = "string"
    required = true
    enum = ["restricted", "standard", "elevated", "trusted"]
          description = "Security level for template running"
  }

  # Template rendering engine
  engine = {
    type = "string"
    required = true
    enum = ["go-template", "jinja2", "handlebars", "custom"]
    description = "Template rendering engine"
  }

  # Template variables
  variables = {
    type = "object"
    required = false
    description = "Template variable definitions"
    
    properties = {
      required_variables = {
        type = "array"
        required = false
        description = "List of required variables"
        items = {
          type = "object"
          properties = {
            name = {
              type = "string"
              required = true
              description = "Variable name"
            }
            type = {
              type = "string"
              required = true
              enum = ["string", "number", "boolean", "list", "object", "any"]
              description = "Variable data type"
            }
            required = {
              type = "boolean"
              required = true
              default = true
              description = "Whether variable is required"
            }
            default = {
              type = "any"
              required = false
              description = "Default value for variable"
            }
            description = {
              type = "string"
              required = false
              description = "Variable description"
            }
            validation = {
              type = "object"
              required = false
              description = "Variable validation rules"
              properties = {
                pattern = {
                  type = "string"
                  required = false
                  description = "Regex pattern for validation"
                }
                min = {
                  type = "number"
                  required = false
                  description = "Minimum value"
                }
                max = {
                  type = "number"
                  required = false
                  description = "Maximum value"
                }
                enum = {
                  type = "array"
                  required = false
                  items = {
                    type = "any"
                  }
                  description = "Allowed values"
                }
              }
            }
          }
        }
      }
      
      optional_variables = {
        type = "array"
        required = false
        description = "List of optional variables"
        items = {
          type = "object"
          properties = {
            name = {
              type = "string"
              required = true
              description = "Variable name"
            }
            type = {
              type = "string"
              required = true
              enum = ["string", "number", "boolean", "list", "object", "any"]
              description = "Variable data type"
            }
            default = {
              type = "any"
              required = false
              description = "Default value for variable"
            }
            description = {
              type = "string"
              required = false
              description = "Variable description"
            }
            validation = {
              type = "object"
              required = false
              description = "Variable validation rules"
              properties = {
                pattern = {
                  type = "string"
                  required = false
                  description = "Regex pattern for validation"
                }
                min = {
                  type = "number"
                  required = false
                  description = "Minimum value"
                }
                max = {
                  type = "number"
                  required = false
                  description = "Maximum value"
                }
                enum = {
                  type = "array"
                  required = false
                  items = {
                    type = "any"
                  }
                  description = "Allowed values"
                }
              }
            }
          }
        }
      }
    }
  }

  # Template context data
  context_data = {
    type = "object"
    required = false
    description = "Context data available to templates"
    
    properties = {
      facts = {
        type = "object"
        required = false
        description = "Machine facts available to templates"
        additional_properties = true
      }
      
      variables = {
        type = "object"
        required = false
        description = "Project variables available to templates"
        additional_properties = true
      }
      
      machines = {
        type = "array"
        required = false
        description = "List of machines in inventory"
        items = {
          type = "object"
          properties = {
            name = {
              type = "string"
              required = true
              description = "Machine name"
            }
            host = {
              type = "string"
              required = true
              description = "Machine hostname or IP"
            }
            port = {
              type = "integer"
              required = false
              description = "SSH port"
              min = 1
              max = 65535
            }
            user = {
              type = "string"
              required = false
              description = "SSH username"
            }
            tags = {
              type = "object"
              required = false
              description = "Machine tags"
              additional_properties = {
                type = "string"
              }
            }
            metadata = {
              type = "object"
              required = false
              description = "Machine metadata for templates"
              additional_properties = {
                type = "string"
              }
            }
          }
        }
      }
      
      environment = {
        type = "object"
        required = false
        description = "Environment variables available to templates"
        additional_properties = {
          type = "string"
        }
      }
      
      project = {
        type = "object"
        required = false
        description = "Project information available to templates"
        properties = {
          name = {
            type = "string"
            required = false
            description = "Project name"
          }
          path = {
            type = "string"
            required = false
            description = "Project path"
          }
          config = {
            type = "object"
            required = false
            description = "Project configuration"
            additional_properties = true
          }
        }
      }
    }
  }

  # Template functions and restrictions
  functions = {
    type = "object"
    required = false
    description = "Template functions and restrictions"
    
    properties = {
      allowed_functions = {
        type = "array"
        required = false
        description = "List of allowed template functions"
        items = {
          type = "string"
          enum = [
            "custom", "system", "env", "data", "var", "varOrDefault",
            "project", "projectName", "machines", "machine", "facts", "fact",
            "upper", "lower", "trim", "add", "sub", "mul", "div",
            "join", "split", "replace", "contains", "length", "index",
            "first", "last", "sort", "reverse", "unique", "default"
          ]
        }
      }
      
      restricted_patterns = {
        type = "array"
        required = false
        description = "List of forbidden patterns for security"
        items = {
          type = "string"
          description = "Regex patterns that are not allowed in templates"
        }
        default = [
          "{{.*os\\.Exec.*}}",
          "{{.*exec.*}}",
          "{{.*system.*}}",
          "{{.*eval.*}}",
          "{{.*import.*}}",
          "{{.*reflect.*}}"
        ]
      }
      
      max_template_size = {
        type = "integer"
        required = false
        description = "Maximum template file size in bytes"
        min = 1024
        max = 10485760  # 10MB
        default = 1048576  # 1MB
      }
      
      max_nesting_depth = {
        type = "integer"
        required = false
        description = "Maximum nesting depth for template functions"
        min = 1
        max = 50
        default = 10
      }
      
      max_run_time = {
        type = "integer"
        required = false
        description = "Maximum template run time in milliseconds"
        min = 100
        max = 30000
        default = 5000
      }
      
      max_memory_usage = {
        type = "integer"
        required = false
        description = "Maximum memory usage in bytes"
        min = 1048576  # 1MB
        max = 104857600  # 100MB
        default = 10485760  # 10MB
      }
    }
  }

  # Template metadata
  metadata = {
    type = "object"
    required = false
    description = "Template metadata and information"
    
    properties = {
      name = {
        type = "string"
        required = false
        description = "Template display name"
      }
      
      description = {
        type = "string"
        required = false
        description = "Template description"
      }
      
      author = {
        type = "string"
        required = false
        description = "Template author"
      }
      
      version = {
        type = "string"
        required = false
        description = "Template version"
        pattern = "^\\d+\\.\\d+\\.\\d+$"
      }
      
      tags = {
        type = "array"
        required = false
        description = "Template tags"
        items = {
          type = "string"
        }
      }
      
      output_format = {
        type = "string"
        required = false
        enum = ["config", "script", "documentation", "data", "template", "other"]
        description = "Output format type"
      }
      
      dependencies = {
        type = "array"
        required = false
        description = "Template dependencies"
        items = {
          type = "string"
        }
      }
      
      created_at = {
        type = "string"
        required = false
        format = "date-time"
        description = "Template creation timestamp"
      }
      
      updated_at = {
        type = "string"
        required = false
        format = "date-time"
        description = "Template last update timestamp"
      }
      
      license = {
        type = "string"
        required = false
        description = "Template license"
      }
      
      documentation = {
        type = "object"
        required = false
        description = "Template documentation"
        properties = {
          usage = {
            type = "string"
            required = false
            description = "Usage instructions"
          }
          
          examples = {
            type = "array"
            required = false
            description = "Usage examples"
            items = {
              type = "object"
              properties = {
                name = {
                  type = "string"
                  required = true
                  description = "Example name"
                }
                description = {
                  type = "string"
                  required = false
                  description = "Example description"
                }
                input = {
                  type = "object"
                  required = false
                  description = "Example input data"
                  additional_properties = true
                }
                output = {
                  type = "string"
                  required = false
                  description = "Example output"
                }
              }
            }
          }
          
          api_reference = {
            type = "object"
            required = false
            description = "API reference documentation"
            additional_properties = true
          }
        }
      }
    }
  }

  # Common validation rules (template-agnostic)
  validation = {
    # Template ID validation
    template_id_format = {
      rule = "regex"
      pattern = "^[a-zA-Z0-9._-]+$"
      message = "Template ID must contain only alphanumeric characters, dots, underscores, and hyphens"
    }
    
    # Source path validation
    source_path_format = {
      rule = "regex"
      pattern = "^templates/.*\\.tmpl$"
      message = "Source path must be in templates/ directory with .tmpl extension"
    }
    
    # Template type validation
    template_type_required = {
      rule = "required"
      field = "template_type"
      message = "Template type is required"
    }
    
    # Scope validation
    scope_required = {
      rule = "required"
      field = "scope"
      message = "Template scope is required"
    }
    
    # Security level validation
    security_level_required = {
      rule = "required"
      field = "security_level"
      message = "Security level is required"
    }
    
    # Engine validation
    engine_required = {
      rule = "required"
      field = "engine"
      message = "Template engine is required"
    }
    
    # Variable validation
    variable_name_format = {
      rule = "regex"
      pattern = "^[a-zA-Z_][a-zA-Z0-9_]*$"
      message = "Variable names must start with letter or underscore and contain only alphanumeric characters and underscores"
    }
    
    # Function validation
    function_name_format = {
      rule = "regex"
      pattern = "^[a-zA-Z_][a-zA-Z0-9_]*$"
      message = "Function names must start with letter or underscore and contain only alphanumeric characters and underscores"
    }
    
    # Pattern validation
    pattern_format = {
      rule = "regex"
      pattern = "^.*$"
      message = "Restricted patterns must be valid regex patterns"
    }
    
    # Size validation
    template_size_limits = {
      rule = "range"
      min = 1024
      max = 10485760
      message = "Template size must be between 1KB and 10MB"
    }
    
    # Nesting validation
    nesting_depth_limits = {
      rule = "range"
      min = 1
      max = 50
      message = "Nesting depth must be between 1 and 50"
    }
    
    # Run time validation
run_time_limits = {
      rule = "range"
      min = 100
      max = 30000
              message = "Run time must be between 100ms and 30s"
    }
    
    # Memory usage validation
    memory_usage_limits = {
      rule = "range"
      min = 1048576
      max = 104857600
      message = "Memory usage must be between 1MB and 100MB"
    }
    
    # No circular references
    no_circular_refs = {
      rule = "acyclic"
      message = "No circular references allowed in template definitions"
    }
    
    # No dangerous patterns
    no_dangerous_patterns = {
      rule = "forbidden_patterns"
      patterns = [
        "{{.*os\\.Exec.*}}",
        "{{.*exec.*}}",
        "{{.*system.*}}",
        "{{.*eval.*}}",
        "{{.*import.*}}",
        "{{.*reflect.*}}"
      ]
      message = "Dangerous patterns are not allowed in templates"
    }
  }

  # Common constraints (template-agnostic)
  constraints = {
    # Template scope constraints
    scope_constraints = {
      type = "object"
      description = "Scope-specific constraints"
      
      properties = {
        project_scope = {
          type = "object"
          description = "Project scope constraints"
          properties = {
            max_templates = {
              type = "integer"
              value = 100
              description = "Maximum templates per project"
            }
            max_size_per_template = {
              type = "integer"
              value = 1048576
              description = "Maximum size per template in project scope"
            }
          }
        }
        
        global_scope = {
          type = "object"
          description = "Global scope constraints"
          properties = {
            max_templates = {
              type = "integer"
              value = 1000
              description = "Maximum templates in global scope"
            }
            max_size_per_template = {
              type = "integer"
              value = 5242880
              description = "Maximum size per template in global scope"
            }
          }
        }
        
        machine_scope = {
          type = "object"
          description = "Machine scope constraints"
          properties = {
            max_templates = {
              type = "integer"
              value = 50
              description = "Maximum templates per machine"
            }
            max_size_per_template = {
              type = "integer"
              value = 2097152
              description = "Maximum size per template in machine scope"
            }
          }
        }
      }
    }
    
    # Security constraints
    security_constraints = {
      type = "object"
      description = "Security-related constraints"
      
      properties = {
        restricted_level = {
          type = "object"
          description = "Restricted security level constraints"
          properties = {
            max_functions = {
              type = "integer"
              value = 10
              description = "Maximum functions allowed in restricted mode"
            }
            allowed_functions = {
              type = "array"
              value = ["var", "varOrDefault", "upper", "lower", "trim"]
              description = "Functions allowed in restricted mode"
            }
            max_run_time = {
              type = "integer"
              value = 1000
              description = "Maximum run time in restricted mode (ms)"
            }
          }
        }
        
        standard_level = {
          type = "object"
          description = "Standard security level constraints"
          properties = {
            max_functions = {
              type = "integer"
              value = 25
              description = "Maximum functions allowed in standard mode"
            }
            max_run_time = {
              type = "integer"
              value = 5000
              description = "Maximum run time in standard mode (ms)"
            }
          }
        }
        
        elevated_level = {
          type = "object"
          description = "Elevated security level constraints"
          properties = {
            max_functions = {
              type = "integer"
              value = 50
              description = "Maximum functions allowed in elevated mode"
            }
            max_run_time = {
              type = "integer"
              value = 15000
              description = "Maximum run time in elevated mode (ms)"
            }
          }
        }
        
        trusted_level = {
          type = "object"
          description = "Trusted security level constraints"
          properties = {
            max_functions = {
              type = "integer"
              value = 100
              description = "Maximum functions allowed in trusted mode"
            }
            max_run_time = {
              type = "integer"
              value = 30000
              description = "Maximum run time in trusted mode (ms)"
            }
          }
        }
      }
    }
    
    # Performance constraints
    performance_constraints = {
      type = "object"
      description = "Performance-related constraints"
      
      properties = {
        max_concurrent_templates = {
          type = "integer"
          value = 10
          description = "Maximum concurrent template runs"
        }
        
        max_total_memory = {
          type = "integer"
          value = 104857600
          description = "Maximum total memory usage for all templates (100MB)"
        }
        
        max_total_run_time = {
          type = "integer"
          value = 60000
          description = "Maximum total run time for all templates (60s)"
        }
        
        cache_size = {
          type = "integer"
          value = 100
          description = "Maximum number of cached template results"
        }
        
        cache_ttl = {
          type = "string"
          value = "1h"
          description = "Cache time-to-live for template results"
        }
      }
    }
    
    # Engine constraints
    engine_constraints = {
      type = "object"
      description = "Engine-specific constraints"
      
      properties = {
        go_template = {
          type = "object"
          description = "Go template engine constraints"
          properties = {
            max_nesting = {
              type = "integer"
              value = 10
              description = "Maximum nesting depth for Go templates"
            }
            max_functions = {
              type = "integer"
              value = 50
              description = "Maximum functions for Go templates"
            }
          }
        }
        
        jinja2 = {
          type = "object"
          description = "Jinja2 template engine constraints"
          properties = {
            max_nesting = {
              type = "integer"
              value = 15
              description = "Maximum nesting depth for Jinja2 templates"
            }
            max_functions = {
              type = "integer"
              value = 75
              description = "Maximum functions for Jinja2 templates"
            }
          }
        }
        
        handlebars = {
          type = "object"
          description = "Handlebars template engine constraints"
          properties = {
            max_nesting = {
              type = "integer"
              value = 8
              description = "Maximum nesting depth for Handlebars templates"
            }
            max_functions = {
              type = "integer"
              value = 30
              description = "Maximum functions for Handlebars templates"
            }
          }
        }
        
        custom = {
          type = "object"
          description = "Custom template engine constraints"
          properties = {
            max_nesting = {
              type = "integer"
              value = 20
              description = "Maximum nesting depth for custom templates"
            }
            max_functions = {
              type = "integer"
              value = 100
              description = "Maximum functions for custom templates"
            }
          }
        }
      }
    }
  }

  # Common metadata (template-agnostic)
  metadata = {
    # Schema version
    schema_version = {
      type = "string"
      value = "1.0.0"
      description = "Template structure schema version"
    }
    
    # Schema type
    schema_type = {
      type = "string"
      value = "template-structure"
      description = "Schema type identifier"
    }
    
    # Last updated
    last_updated = {
      type = "string"
      value = "2024-01-01"
      description = "Last schema update date"
    }
    
    # Description
    description = {
      type = "string"
      value = "Common template structure definitions for all template-related schemas"
      description = "Schema description"
    }
    
    # Author
    author = {
      type = "string"
      value = "Spooky Template System"
      description = "Schema author"
    }
    
    # License
    license = {
      type = "string"
      value = "MIT"
      description = "Schema license"
    }
    
    # Dependencies
    dependencies = {
      type = "array"
      value = []
      description = "Schema dependencies"
    }
    
    # Tags
    tags = {
      type = "array"
      value = ["template", "structure", "base", "common"]
      description = "Schema tags"
    }
  }
} 