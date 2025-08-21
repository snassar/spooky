# Template Context Schema
# Enhanced template context schema with composition pattern
# Includes template-structure base and context-specific features

# Schema metadata
metadata {
  schema_version = "0.20250809.0"
  schema_type = "template-context"
  schema_name = "Template Context Schema"
  last_updated = "2024-01-01"
  compatibility = ["0.20250809.0"]
  description = "Enhanced template context schema with composition pattern - includes template-structure base and context-specific features"
  
  # ScalVer format: 0.YYYYMMDD.N
  # - 0: Development phase
  # - 20250809: Date (9 August 2025)
  # - 0: Patch version
  scalver_format = "0.20250809.0"
}

template_context {
  include = "template-structure"
  
  scope = "context"
  storage_location = "<project>/templates/context.hcl"
  description = "Template context data available during rendering"
  
  # Context-specific features
  context_features = {
    type = "object"
    description = "Context-specific features and capabilities"
    
    properties = {
      # Data binding capabilities
      data_binding = {
        type = "object"
        description = "Data binding features"
        properties = {
          facts_binding = {
            type = "boolean"
            value = true
            description = "Support for binding machine facts to templates"
          }
          variables_binding = {
            type = "boolean"
            value = true
            description = "Support for binding project variables to templates"
          }
          machines_binding = {
            type = "boolean"
            value = true
            description = "Support for binding machine inventory to templates"
          }
          environment_binding = {
            type = "boolean"
            value = true
            description = "Support for binding environment variables to templates"
          }
          project_binding = {
            type = "boolean"
            value = true
            description = "Support for binding project information to templates"
          }
        }
      }
      
      # Context resolution features
      context_resolution = {
        type = "object"
        description = "Context resolution features"
        properties = {
          lazy_loading = {
            type = "boolean"
            value = true
            description = "Support for lazy loading of context data"
          }
          caching = {
            type = "boolean"
            value = true
            description = "Support for caching resolved context data"
          }
          validation = {
            type = "boolean"
            value = true
            description = "Support for validating context data"
          }
          transformation = {
            type = "boolean"
            value = true
            description = "Support for transforming context data"
          }
        }
      }
      
      # Context scoping features
      context_scoping = {
        type = "object"
        description = "Context scoping features"
        properties = {
          global_scope = {
            type = "boolean"
            value = true
            description = "Support for global context scope"
          }
          project_scope = {
            type = "boolean"
            value = true
            description = "Support for project context scope"
          }
          machine_scope = {
            type = "boolean"
            value = true
            description = "Support for machine context scope"
          }
          template_scope = {
            type = "boolean"
            value = true
            description = "Support for template-specific context scope"
          }
        }
      }
    }
  }
  
  # Context-specific validation
  context_validation = {
    type = "object"
    description = "Context-specific validation rules"
    
    properties = {
      # Facts validation
      facts_validation = {
        type = "object"
        description = "Validation rules for facts data"
        properties = {
          required_facts = {
            type = "array"
            required = false
            items = {
              type = "string"
            }
            description = "Required facts that must be present"
          }
          facts_format = {
            type = "string"
            value = "json"
            enum = ["json", "hcl"]
            description = "Expected format for facts data"
          }
          facts_schema = {
            type = "string"
            required = false
            description = "Schema for validating facts data"
          }
        }
      }
      
      # Variables validation
      variables_validation = {
        type = "object"
        description = "Validation rules for variables data"
        properties = {
          required_variables = {
            type = "array"
            required = false
            items = {
              type = "string"
            }
            description = "Required variables that must be present"
          }
          variables_format = {
            type = "string"
            value = "hcl"
            enum = ["hcl", "json"]
            description = "Expected format for variables data"
          }
          variables_schema = {
            type = "string"
            required = false
            description = "Schema for validating variables data"
          }
        }
      }
      
      # Machines validation
      machines_validation = {
        type = "object"
        description = "Validation rules for machines data"
        properties = {
          required_machines = {
            type = "array"
            required = false
            items = {
              type = "string"
            }
            description = "Required machines that must be present"
          }
          machines_format = {
            type = "string"
            value = "hcl"
            enum = ["hcl", "json"]
            description = "Expected format for machines data"
          }
          machines_schema = {
            type = "string"
            required = false
            description = "Schema for validating machines data"
          }
        }
      }
      
      # Environment validation
      environment_validation = {
        type = "object"
        description = "Validation rules for environment data"
        properties = {
          required_environment = {
            type = "array"
            required = false
            items = {
              type = "string"
            }
            description = "Required environment variables that must be present"
          }
          environment_format = {
            type = "string"
            value = "key-value"
            enum = ["key-value", "json", "hcl"]
            description = "Expected format for environment data"
          }
        }
      }
      
      # Project validation
      project_validation = {
        type = "object"
        description = "Validation rules for project data"
        properties = {
          required_project = {
            type = "array"
            required = false
            items = {
              type = "string"
            }
            description = "Required project fields that must be present"
          }
          project_format = {
            type = "string"
            value = "hcl"
            enum = ["hcl", "json"]
            description = "Expected format for project data"
          }
          project_schema = {
            type = "string"
            required = false
            description = "Schema for validating project data"
          }
        }
      }
      
      # Context data validation
      context_data_validation = {
        type = "object"
        description = "General context data validation"
        properties = {
          data_integrity = {
            type = "boolean"
            value = true
            description = "Validate data integrity"
          }
          data_freshness = {
            type = "string"
            value = "5m"
            description = "Maximum age of context data"
          }
          data_size_limits = {
            type = "object"
            properties = {
              max_facts_size = {
                type = "integer"
                value = 10485760
                description = "Maximum size of facts data (10MB)"
              }
              max_variables_size = {
                type = "integer"
                value = 1048576
                description = "Maximum size of variables data (1MB)"
              }
              max_machines_size = {
                type = "integer"
                value = 5242880
                description = "Maximum size of machines data (5MB)"
              }
              max_environment_size = {
                type = "integer"
                value = 1048576
                description = "Maximum size of environment data (1MB)"
              }
              max_project_size = {
                type = "integer"
                value = 2097152
                description = "Maximum size of project data (2MB)"
              }
            }
          }
        }
      }
    }
  }
  
  # Context-specific constraints
  context_constraints = {
    type = "object"
    description = "Context-specific constraints and limits"
    
    properties = {
      # Performance constraints
      performance_constraints = {
        type = "object"
        description = "Performance-related constraints"
        properties = {
          max_context_size = {
            type = "integer"
            value = 52428800
            description = "Maximum total context size (50MB)"
          }
          max_context_resolution_time = {
            type = "integer"
            value = 5000
            description = "Maximum context resolution time (5s)"
          }
          max_context_cache_size = {
            type = "integer"
            value = 100
            description = "Maximum number of cached contexts"
          }
          max_context_cache_ttl = {
            type = "string"
            value = "10m"
            description = "Maximum context cache TTL"
          }
        }
      }
      
      # Security constraints
      security_constraints = {
        type = "object"
        description = "Security-related constraints"
        properties = {
          sensitive_data_filtering = {
            type = "boolean"
            value = true
            description = "Filter sensitive data from context"
          }
          data_encryption = {
            type = "boolean"
            value = false
            description = "Encrypt context data in memory"
          }
          access_control = {
            type = "boolean"
            value = true
            description = "Enforce access control on context data"
          }
          audit_logging = {
            type = "boolean"
            value = true
            description = "Log context access for auditing"
          }
        }
      }
      
      # Data constraints
      data_constraints = {
        type = "object"
        description = "Data-related constraints"
        properties = {
          max_facts_count = {
            type = "integer"
            value = 1000
            description = "Maximum number of facts per context"
          }
          max_variables_count = {
            type = "integer"
            value = 500
            description = "Maximum number of variables per context"
          }
          max_machines_count = {
            type = "integer"
            value = 100
            description = "Maximum number of machines per context"
          }
          max_environment_count = {
            type = "integer"
            value = 200
            description = "Maximum number of environment variables per context"
          }
          data_retention = {
            type = "string"
            value = "1h"
            description = "Context data retention period"
          }
        }
      }
      
      # Scope constraints
      scope_constraints = {
        type = "object"
        description = "Scope-related constraints"
        properties = {
          global_context_limit = {
            type = "integer"
            value = 10
            description = "Maximum number of global contexts"
          }
          project_context_limit = {
            type = "integer"
            value = 50
            description = "Maximum number of project contexts"
          }
          machine_context_limit = {
            type = "integer"
            value = 20
            description = "Maximum number of machine contexts"
          }
          template_context_limit = {
            type = "integer"
            value = 100
            description = "Maximum number of template contexts"
          }
        }
      }
    }
  }
  
  # Context-specific extensions
  context_extensions = {
    type = "object"
    description = "Context-specific extensions and enhancements"
    
    properties = {
      # Data transformation extensions
      data_transformations = {
        type = "object"
        description = "Data transformation capabilities"
        properties = {
          facts_transformation = {
            type = "object"
            description = "Facts data transformation"
            properties = {
              flatten_nested = {
                type = "boolean"
                value = true
                description = "Flatten nested facts structure"
              }
              normalize_keys = {
                type = "boolean"
                value = true
                description = "Normalize fact keys to consistent format"
              }
              filter_sensitive = {
                type = "boolean"
                value = true
                description = "Filter sensitive facts data"
              }
              add_metadata = {
                type = "boolean"
                value = true
                description = "Add metadata to facts"
              }
            }
          }
          
          variables_transformation = {
            type = "object"
            description = "Variables data transformation"
            properties = {
              resolve_references = {
                type = "boolean"
                value = true
                description = "Resolve variable references"
              }
              validate_types = {
                type = "boolean"
                value = true
                description = "Validate variable types"
              }
              apply_defaults = {
                type = "boolean"
                value = true
                description = "Apply default values"
              }
              expand_environment = {
                type = "boolean"
                value = true
                description = "Expand environment variables"
              }
            }
          }
          
          machines_transformation = {
            type = "object"
            description = "Machines data transformation"
            properties = {
              resolve_hostnames = {
                type = "boolean"
                value = true
                description = "Resolve machine hostnames"
              }
              validate_connectivity = {
                type = "boolean"
                value = false
                description = "Validate machine connectivity"
              }
              add_facts = {
                type = "boolean"
                value = true
                description = "Add facts to machines"
              }
              group_by_tags = {
                type = "boolean"
                value = true
                description = "Group machines by tags"
              }
            }
          }
          
          environment_transformation = {
            type = "object"
            description = "Environment data transformation"
            properties = {
              expand_variables = {
                type = "boolean"
                value = true
                description = "Expand environment variables"
              }
              filter_sensitive = {
                type = "boolean"
                value = true
                description = "Filter sensitive environment variables"
              }
              normalize_case = {
                type = "boolean"
                value = true
                description = "Normalize environment variable case"
              }
              add_prefix = {
                type = "string"
                value = "SPOOKY_"
                description = "Prefix for spooky environment variables"
              }
            }
          }
        }
      }
      
      # Context composition extensions
      context_composition = {
        type = "object"
        description = "Context composition capabilities"
        properties = {
          merge_strategies = {
            type = "object"
            description = "Context merge strategies"
            properties = {
              facts_merge = {
                type = "string"
                value = "deep"
                enum = ["shallow", "deep", "replace", "append"]
                description = "Strategy for merging facts"
              }
              variables_merge = {
                type = "string"
                value = "replace"
                enum = ["shallow", "deep", "replace", "append"]
                description = "Strategy for merging variables"
              }
              machines_merge = {
                type = "string"
                value = "append"
                enum = ["shallow", "deep", "replace", "append"]
                description = "Strategy for merging machines"
              }
              environment_merge = {
                type = "string"
                value = "replace"
                enum = ["shallow", "deep", "replace", "append"]
                description = "Strategy for merging environment"
              }
            }
          }
          
          inheritance = {
            type = "object"
            description = "Context inheritance capabilities"
            properties = {
              inherit_global = {
                type = "boolean"
                value = true
                description = "Inherit from global context"
              }
              inherit_project = {
                type = "boolean"
                value = true
                description = "Inherit from project context"
              }
              inherit_machine = {
                type = "boolean"
                value = true
                description = "Inherit from machine context"
              }
              override_rules = {
                type = "object"
                description = "Context override rules"
                properties = {
                  allow_facts_override = {
                    type = "boolean"
                    value = true
                    description = "Allow facts override"
                  }
                  allow_variables_override = {
                    type = "boolean"
                    value = true
                    description = "Allow variables override"
                  }
                  allow_machines_override = {
                    type = "boolean"
                    value = false
                    description = "Allow machines override"
                  }
                  allow_environment_override = {
                    type = "boolean"
                    value = true
                    description = "Allow environment override"
                  }
                }
              }
            }
          }
        }
      }
      
      # Context validation extensions
      context_validation_extensions = {
        type = "object"
        description = "Extended context validation capabilities"
        properties = {
          schema_validation = {
            type = "object"
            description = "Schema-based validation"
            properties = {
              validate_facts_schema = {
                type = "boolean"
                value = true
                description = "Validate facts against schema"
              }
              validate_variables_schema = {
                type = "boolean"
                value = true
                description = "Validate variables against schema"
              }
              validate_machines_schema = {
                type = "boolean"
                value = true
                description = "Validate machines against schema"
              }
              validate_project_schema = {
                type = "boolean"
                value = true
                description = "Validate project against schema"
              }
            }
          }
          
          data_validation = {
            type = "object"
            description = "Data-specific validation"
            properties = {
              validate_facts_integrity = {
                type = "boolean"
                value = true
                description = "Validate facts data integrity"
              }
              validate_variables_types = {
                type = "boolean"
                value = true
                description = "Validate variable types"
              }
              validate_machines_connectivity = {
                type = "boolean"
                value = false
                description = "Validate machine connectivity"
              }
              validate_environment_format = {
                type = "boolean"
                value = true
                description = "Validate environment variable format"
              }
            }
          }
          
          cross_reference_validation = {
            type = "object"
            description = "Cross-reference validation"
            properties = {
              validate_facts_references = {
                type = "boolean"
                value = true
                description = "Validate facts cross-references"
              }
              validate_variables_references = {
                type = "boolean"
                value = true
                description = "Validate variables cross-references"
              }
              validate_machines_references = {
                type = "boolean"
                value = true
                description = "Validate machines cross-references"
              }
              validate_project_references = {
                type = "boolean"
                value = true
                description = "Validate project cross-references"
              }
            }
          }
        }
      }
      
      # Context monitoring extensions
      context_monitoring = {
        type = "object"
        description = "Context monitoring and observability"
        properties = {
          metrics = {
            type = "object"
            description = "Context metrics collection"
            properties = {
              collect_resolution_time = {
                type = "boolean"
                value = true
                description = "Collect context resolution time"
              }
              collect_data_size = {
                type = "boolean"
                value = true
                description = "Collect context data size"
              }
              collect_cache_hits = {
                type = "boolean"
                value = true
                description = "Collect cache hit metrics"
              }
              collect_validation_errors = {
                type = "boolean"
                value = true
                description = "Collect validation error metrics"
              }
            }
          }
          
          logging = {
            type = "object"
            description = "Context logging capabilities"
            properties = {
              log_context_access = {
                type = "boolean"
                value = true
                description = "Log context access"
              }
              log_data_changes = {
                type = "boolean"
                value = true
                description = "Log context data changes"
              }
              log_validation_errors = {
                type = "boolean"
                value = true
                description = "Log validation errors"
              }
              log_performance_issues = {
                type = "boolean"
                value = true
                description = "Log performance issues"
              }
            }
          }
          
          alerting = {
            type = "object"
            description = "Context alerting capabilities"
            properties = {
              alert_on_validation_failure = {
                type = "boolean"
                value = true
                description = "Alert on validation failures"
              }
              alert_on_performance_degradation = {
                type = "boolean"
                value = true
                description = "Alert on performance degradation"
              }
              alert_on_data_corruption = {
                type = "boolean"
                value = true
                description = "Alert on data corruption"
              }
              alert_on_security_violation = {
                type = "boolean"
                value = true
                description = "Alert on security violations"
              }
            }
          }
        }
      }
    }
  }
  
  # Context-specific metadata
  context_metadata = {
    type = "object"
    description = "Context-specific metadata and information"
    
    properties = {
      # Context versioning
      versioning = {
        type = "object"
        description = "Context versioning information"
        properties = {
          context_version = {
            type = "string"
            value = "1.0.0"
            description = "Context schema version"
          }
          data_version = {
            type = "string"
            required = false
            description = "Data version for context"
          }
          schema_version = {
            type = "string"
            value = "1.0.0"
            description = "Schema version for context"
          }
        }
      }
      
      # Context lifecycle
      lifecycle = {
        type = "object"
        description = "Context lifecycle information"
        properties = {
          created_at = {
            type = "string"
            required = false
            format = "date-time"
            description = "Context creation timestamp"
          }
          updated_at = {
            type = "string"
            required = false
            format = "date-time"
            description = "Context last update timestamp"
          }
          expires_at = {
            type = "string"
            required = false
            format = "date-time"
            description = "Context expiration timestamp"
          }
          ttl = {
            type = "string"
            value = "1h"
            description = "Context time-to-live"
          }
        }
      }
      
      # Context ownership
      ownership = {
        type = "object"
        description = "Context ownership information"
        properties = {
          owner = {
            type = "string"
            required = false
            description = "Context owner"
          }
          project = {
            type = "string"
            required = false
            description = "Associated project"
          }
          machine = {
            type = "string"
            required = false
            description = "Associated machine"
          }
          template = {
            type = "string"
            required = false
            description = "Associated template"
          }
        }
      }
      
      # Context tags
      tags = {
        type = "array"
        required = false
        description = "Context tags"
        items = {
          type = "string"
        }
      }
      
      # Context description
      description = {
        type = "string"
        required = false
        description = "Context description"
      }
    }
  }
  
}