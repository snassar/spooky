# Variables JSON Schema
# Schema for variables stored in JSON format
# This schema composes the base variables structure with JSON-specific storage rules

# Schema metadata
metadata {
  schema_version = "0.20250809.0"
  schema_type = "variables-json"
  schema_name = "Variables JSON Schema"
  last_updated = "2024-01-01"
  compatibility = ["0.20250809.0"]
  description = "Schema for variables stored in JSON format - composes the base variables structure with JSON-specific storage rules"
  
  # ScalVer format: 0.YYYYMMDD.N
  # - 0: Development phase
  # - 20250809: Date (9 August 2025)
  # - 0: Patch version
  scalver_format = "0.20250809.0"
}

# Variables in JSON format
variables {
  # Include the base variables structure
  include = "variables-structure"
  
  # Variables specific metadata
  scope = "export/import"
  storage_location = "exported variables in JSON format"
  description = "Variables exported/imported in JSON format"
  
  # Variables specific validation
  variables_validation = {
    # JSON must be valid JSON syntax
    valid_json = {
      rule = "json"
      message = "Variables must be valid JSON"
    }
    
    # Variables must have a metadata section
    metadata_required = {
      rule = "required"
      field = "metadata"
      message = "JSON export must include metadata section"
    }
    
    # Variables must have a variables section
    variables_required = {
      rule = "required"
      field = "variables"
      message = "JSON export must include variables section"
    }
    
    # Variable names must be valid JSON keys
    valid_json_keys = {
      rule = "pattern"
      field = "variables.*"
      pattern = "^[a-z][a-z0-9_]*$"
      message = "Variable names must be valid JSON keys and follow the same naming convention as HCL variables"
    }
    
    # JSON values must be valid JSON types
    valid_json_types = {
      rule = "json_types"
      message = "Variable values must be valid JSON types (string, number, boolean, null, array, object)"
    }
  }
  
  # Variables specific constraints
  variables_constraints = {
    # JSON format limitations
    json_limitations = {
      type = "object"
      description = "JSON format limitations for variables"
      
      properties = {
        no_comments = {
          type = "boolean"
          value = true
          description = "JSON does not support comments"
        }
        
        no_trailing_commas = {
          type = "boolean"
          value = true
          description = "JSON does not support trailing commas"
        }
        
        no_unquoted_keys = {
          type = "boolean"
          value = true
          description = "JSON requires quoted keys"
        }
        
        no_expressions = {
          type = "boolean"
          value = true
          description = "JSON does not support expressions or interpolation"
        }
        
        no_heredoc = {
          type = "boolean"
          value = true
          description = "JSON does not support heredoc syntax"
        }
      }
    }
    
    # JSON export features
    export_features = {
      type = "object"
      description = "JSON export-specific features"
      
      properties = {
        include_metadata = {
          type = "boolean"
          value = true
          description = "JSON export includes metadata section"
        }
        
        include_sources = {
          type = "boolean"
          value = true
          description = "JSON export can include source information"
        }
        
        include_definitions = {
          type = "boolean"
          value = true
          description = "JSON export can include full variable definitions"
        }
        
        resolved_values = {
          type = "boolean"
          value = true
          description = "JSON export can include resolved variable values"
        }
      }
    }
  }
  
  # JSON-specific structure
  json_structure = {
    # JSON export structure
    export_format = {
      type = "object"
      description = "JSON export structure"
      
      properties = {
        metadata = {
          type = "object"
          description = "Export metadata"
          
          properties = {
            version = {
              type = "string"
              description = "Export format version"
            }
            
            exported_at = {
              type = "string"
              format = "date-time"
              description = "Export timestamp"
            }
            
            project_path = {
              type = "string"
              description = "Source project path"
            }
            
            format = {
              type = "string"
              description = "Export format (json)"
            }
            
            source_files = {
              type = "array"
              items = {
                type = "string"
              }
              description = "Source files used for export"
            }
            
            variable_count = {
              type = "integer"
              description = "Number of exported variables"
            }
            
            options = {
              type = "object"
              description = "Export options used"
            }
          }
        }
        
        variables = {
          type = "object"
          description = "Exported variables - key-value pairs where keys are variable names and values are variable definitions or resolved values"
          additional_properties = {
            type = "object"
            description = "Variable definition or resolved value"
            properties = {
              name = {
                type = "string"
                description = "Variable name"
              }
              type = {
                type = "string"
                description = "Variable type"
              }
              value = {
                type = "any"
                description = "Variable value (resolved or default)"
              }
              description = {
                type = "string"
                description = "Variable description"
              }
              required = {
                type = "boolean"
                description = "Whether the variable is required"
              }
              sensitive = {
                type = "boolean"
                description = "Whether the variable contains sensitive data"
              }
              encrypted = {
                type = "boolean"
                description = "Whether the variable is encrypted"
              }
              scope = {
                type = "string"
                description = "Variable scope"
              }
              dependencies = {
                type = "array"
                items = {
                  type = "string"
                }
                description = "Variable dependencies"
              }
              metadata = {
                type = "object"
                description = "Additional variable metadata"
              }
            }
          }
        }
        
        sources = {
          type = "object"
          description = "Variable source information"
          additional_properties = {
            type = "object"
            properties = {
              type = {
                type = "string"
                description = "Source type"
              }
              
              file = {
                type = "string"
                description = "Source file"
              }
              
              line = {
                type = "integer"
                description = "Source line number"
              }
              
              priority = {
                type = "integer"
                description = "Source priority"
              }
            }
          }
        }
        
        definitions = {
          type = "object"
          description = "Full variable definitions"
          additional_properties = {
            type = "object"
            description = "Variable definition"
          }
        }
      }
    }
  }
} 