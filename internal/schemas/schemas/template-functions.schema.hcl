# Template Functions Schema
# Enhanced template functions schema with composition pattern
# Includes template-structure base and function-specific features

# Schema metadata
metadata {
  schema_version = "0.20250809.0"
  schema_type = "template-functions"
  schema_name = "Template Functions Schema"
  last_updated = "2024-01-01"
  compatibility = ["0.20250809.0"]
  description = "Enhanced template functions schema with composition pattern - includes template-structure base and function-specific features"
  
  # ScalVer format: 0.YYYYMMDD.N
  # - 0: Development phase
  # - 20250809: Date (9 August 2025)
  # - 0: Patch version
  scalver_format = "0.20250809.0"
}

template_functions {
  include = "template-structure"
  
  scope = "functions"
  storage_location = "<project>/templates/functions.hcl"
  description = "Template functions and security restrictions"
  
  # Function-specific features
  function_features = {
    type = "object"
    description = "Function-specific features and capabilities"
    
    properties = {
      # Function categories
      function_categories = {
        type = "object"
        description = "Function categories and capabilities"
        properties = {
          data_functions = {
            type = "object"
            description = "Data access functions"
            properties = {
              var = {
                type = "boolean"
                value = true
                description = "Access template variables"
              }
              varOrDefault = {
                type = "boolean"
                value = true
                description = "Access template variables with defaults"
              }
              facts = {
                type = "boolean"
                value = true
                description = "Access machine facts"
              }
              fact = {
                type = "boolean"
                value = true
                description = "Access specific machine fact"
              }
              env = {
                type = "boolean"
                value = true
                description = "Access environment variables"
              }
              data = {
                type = "boolean"
                value = true
                description = "Access external data sources"
              }
            }
          }
          
          project_functions = {
            type = "object"
            description = "Project-related functions"
            properties = {
              project = {
                type = "boolean"
                value = true
                description = "Access project information"
              }
              projectName = {
                type = "boolean"
                value = true
                description = "Get project name"
              }
              machines = {
                type = "boolean"
                value = true
                description = "Access machine inventory"
              }
              machine = {
                type = "boolean"
                value = true
                description = "Access specific machine"
              }
            }
          }
          
          string_functions = {
            type = "object"
            description = "String manipulation functions"
            properties = {
              upper = {
                type = "boolean"
                value = true
                description = "Convert string to uppercase"
              }
              lower = {
                type = "boolean"
                value = true
                description = "Convert string to lowercase"
              }
              trim = {
                type = "boolean"
                value = true
                description = "Trim whitespace from string"
              }
              join = {
                type = "boolean"
                value = true
                description = "Join array elements into string"
              }
              split = {
                type = "boolean"
                value = true
                description = "Split string into array"
              }
              replace = {
                type = "boolean"
                value = true
                description = "Replace string patterns"
              }
              contains = {
                type = "boolean"
                value = true
                description = "Check if string contains substring"
              }
            }
          }
          
          math_functions = {
            type = "object"
            description = "Mathematical functions"
            properties = {
              add = {
                type = "boolean"
                value = true
                description = "Add numbers"
              }
              sub = {
                type = "boolean"
                value = true
                description = "Subtract numbers"
              }
              mul = {
                type = "boolean"
                value = true
                description = "Multiply numbers"
              }
              div = {
                type = "boolean"
                value = true
                description = "Divide numbers"
              }
            }
          }
          
          array_functions = {
            type = "object"
            description = "Array manipulation functions"
            properties = {
              length = {
                type = "boolean"
                value = true
                description = "Get array length"
              }
              index = {
                type = "boolean"
                value = true
                description = "Get array element by index"
              }
              first = {
                type = "boolean"
                value = true
                description = "Get first array element"
              }
              last = {
                type = "boolean"
                value = true
                description = "Get last array element"
              }
              sort = {
                type = "boolean"
                value = true
                description = "Sort array elements"
              }
              reverse = {
                type = "boolean"
                value = true
                description = "Reverse array order"
              }
              unique = {
                type = "boolean"
                value = true
                description = "Remove duplicate array elements"
              }
            }
          }
          
          system_functions = {
            type = "object"
            description = "System-related functions"
            properties = {
              system = {
                type = "boolean"
                value = false
                description = "Execute system commands (restricted)"
              }
              custom = {
                type = "boolean"
                value = true
                description = "Execute custom functions"
              }
              default = {
                type = "boolean"
                value = true
                description = "Provide default values"
              }
            }
          }
        }
      }
      
      # Function execution features
      execution_features = {
        type = "object"
        description = "Function execution capabilities"
        properties = {
          lazy_evaluation = {
            type = "boolean"
            value = true
            description = "Support for lazy function evaluation"
          }
          caching = {
            type = "boolean"
            value = true
            description = "Support for function result caching"
          }
          parallel_execution = {
            type = "boolean"
            value = false
            description = "Support for parallel function execution"
          }
          error_handling = {
            type = "boolean"
            value = true
            description = "Support for function error handling"
          }
          timeout_handling = {
            type = "boolean"
            value = true
            description = "Support for function timeout handling"
          }
        }
      }
      
      # Function security features
      security_features = {
        type = "object"
        description = "Function security capabilities"
        properties = {
          sandboxing = {
            type = "boolean"
            value = true
            description = "Support for function sandboxing"
          }
          access_control = {
            type = "boolean"
            value = true
            description = "Support for function access control"
          }
          audit_logging = {
            type = "boolean"
            value = true
            description = "Support for function audit logging"
          }
          pattern_filtering = {
            type = "boolean"
            value = true
            description = "Support for pattern-based filtering"
          }
          resource_limiting = {
            type = "boolean"
            value = true
            description = "Support for resource limiting"
          }
        }
      }
    }
  }
  
  # Function-specific validation
  function_validation = {
    type = "object"
    description = "Function-specific validation rules"
    
    properties = {
      # Function name validation
      function_name_validation = {
        type = "object"
        description = "Validation rules for function names"
        properties = {
          name_pattern = {
            type = "string"
            value = "^[a-zA-Z_][a-zA-Z0-9_]*$"
            description = "Regex pattern for function names"
          }
          reserved_names = {
            type = "array"
            items = {
              type = "string"
            }
            value = ["exec", "system", "eval", "import", "reflect", "os"]
            description = "Reserved function names"
          }
          max_name_length = {
            type = "integer"
            value = 50
            description = "Maximum function name length"
          }
        }
      }
      
      # Function argument validation
      function_argument_validation = {
        type = "object"
        description = "Validation rules for function arguments"
        properties = {
          max_arguments = {
            type = "integer"
            value = 10
            description = "Maximum number of function arguments"
          }
          argument_types = {
            type = "array"
            items = {
              type = "string"
            }
            value = ["string", "number", "boolean", "array", "object"]
            description = "Allowed argument types"
          }
          required_arguments = {
            type = "object"
            description = "Required arguments for specific functions"
            additional_properties = {
              type = "array"
              items = {
                type = "string"
              }
            }
          }
        }
      }
      
      # Function return validation
      function_return_validation = {
        type = "object"
        description = "Validation rules for function return values"
        properties = {
          return_types = {
            type = "array"
            items = {
              type = "string"
            }
            value = ["string", "number", "boolean", "array", "object", "null"]
            description = "Allowed return types"
          }
          max_return_size = {
            type = "integer"
            value = 1048576
            description = "Maximum return value size (1MB)"
          }
          return_value_validation = {
            type = "boolean"
            value = true
            description = "Validate function return values"
          }
        }
      }
      
      # Function pattern validation
      function_pattern_validation = {
        type = "object"
        description = "Validation rules for function patterns"
        properties = {
          pattern_syntax = {
            type = "string"
            value = "regex"
            enum = ["regex", "glob", "wildcard"]
            description = "Pattern syntax type"
          }
          pattern_validation = {
            type = "boolean"
            value = true
            description = "Validate pattern syntax"
          }
          dangerous_patterns = {
            type = "array"
            items = {
              type = "string"
            }
            value = [
              "{{.*os\\.Exec.*}}",
              "{{.*exec.*}}",
              "{{.*system.*}}",
              "{{.*eval.*}}",
              "{{.*import.*}}",
              "{{.*reflect.*}}",
              "{{.*\\.\\./.*}}",
              "{{.*/etc/.*}}",
              "{{.*/proc/.*}}",
              "{{.*/sys/.*}}"
            ]
            description = "Dangerous patterns to block"
          }
        }
      }
      
      # Function performance validation
      function_performance_validation = {
        type = "object"
        description = "Validation rules for function performance"
        properties = {
          max_execution_time = {
            type = "integer"
            value = 5000
            description = "Maximum function execution time (ms)"
          }
          max_memory_usage = {
            type = "integer"
            value = 10485760
            description = "Maximum function memory usage (10MB)"
          }
          max_cpu_usage = {
            type = "number"
            value = 0.5
            description = "Maximum function CPU usage (50%)"
          }
          max_io_operations = {
            type = "integer"
            value = 1000
            description = "Maximum I/O operations per function"
          }
        }
      }
      
      # Function security validation
      function_security_validation = {
        type = "object"
        description = "Validation rules for function security"
        properties = {
          file_access_validation = {
            type = "boolean"
            value = true
            description = "Validate file access patterns"
          }
          network_access_validation = {
            type = "boolean"
            value = true
            description = "Validate network access patterns"
          }
          process_access_validation = {
            type = "boolean"
            value = true
            description = "Validate process access patterns"
          }
          environment_access_validation = {
            type = "boolean"
            value = true
            description = "Validate environment access patterns"
          }
        }
      }
    }
  }
  
  # Function-specific constraints
  function_constraints = {
    type = "object"
    description = "Function-specific constraints and limits"
    
    properties = {
      # Performance constraints
      performance_constraints = {
        type = "object"
        description = "Performance-related constraints"
        properties = {
          max_concurrent_functions = {
            type = "integer"
            value = 20
            description = "Maximum concurrent function executions"
          }
          max_total_execution_time = {
            type = "integer"
            value = 30000
            description = "Maximum total function execution time (30s)"
          }
          max_total_memory = {
            type = "integer"
            value = 52428800
            description = "Maximum total function memory usage (50MB)"
          }
          function_cache_size = {
            type = "integer"
            value = 200
            description = "Maximum number of cached function results"
          }
          function_cache_ttl = {
            type = "string"
            value = "5m"
            description = "Function cache time-to-live"
          }
        }
      }
      
      # Security constraints
      security_constraints = {
        type = "object"
        description = "Security-related constraints"
        properties = {
          function_sandboxing = {
            type = "object"
            description = "Function sandboxing constraints"
            properties = {
              enable_sandboxing = {
                type = "boolean"
                value = true
                description = "Enable function sandboxing"
              }
              sandbox_timeout = {
                type = "integer"
                value = 10000
                description = "Sandbox timeout (10s)"
              }
              sandbox_memory_limit = {
                type = "integer"
                value = 20971520
                description = "Sandbox memory limit (20MB)"
              }
              sandbox_cpu_limit = {
                type = "number"
                value = 0.25
                description = "Sandbox CPU limit (25%)"
              }
            }
          }
          
          access_control_constraints = {
            type = "object"
            description = "Access control constraints"
            properties = {
              function_permissions = {
                type = "object"
                description = "Function permission levels"
                properties = {
                  read_only = {
                    type = "array"
                    value = ["var", "varOrDefault", "facts", "fact", "env", "data"]
                    description = "Read-only functions"
                  }
                  read_write = {
                    type = "array"
                    value = ["project", "projectName", "machines", "machine"]
                    description = "Read-write functions"
                  }
                  restricted = {
                    type = "array"
                    value = ["system", "custom"]
                    description = "Restricted functions"
                  }
                }
              }
              
              resource_access = {
                type = "object"
                description = "Resource access constraints"
                properties = {
                  file_access = {
                    type = "boolean"
                    value = false
                    description = "Allow file system access"
                  }
                  network_access = {
                    type = "boolean"
                    value = false
                    description = "Allow network access"
                  }
                  process_access = {
                    type = "boolean"
                    value = false
                    description = "Allow process access"
                  }
                  environment_access = {
                    type = "boolean"
                    value = true
                    description = "Allow environment access"
                  }
                }
              }
            }
          }
          
          audit_constraints = {
            type = "object"
            description = "Audit and logging constraints"
            properties = {
              log_function_calls = {
                type = "boolean"
                value = true
                description = "Log all function calls"
              }
              log_function_arguments = {
                type = "boolean"
                value = false
                description = "Log function arguments"
              }
              log_function_results = {
                type = "boolean"
                value = false
                description = "Log function results"
              }
              log_security_violations = {
                type = "boolean"
                value = true
                description = "Log security violations"
              }
            }
          }
        }
      }
      
      # Function category constraints
      function_category_constraints = {
        type = "object"
        description = "Constraints by function category"
        properties = {
          data_functions_constraints = {
            type = "object"
            description = "Data function constraints"
            properties = {
              max_data_size = {
                type = "integer"
                value = 10485760
                description = "Maximum data size (10MB)"
              }
              max_data_operations = {
                type = "integer"
                value = 100
                description = "Maximum data operations per template"
              }
              data_cache_ttl = {
                type = "string"
                value = "10m"
                description = "Data cache TTL"
              }
            }
          }
          
          string_functions_constraints = {
            type = "object"
            description = "String function constraints"
            properties = {
              max_string_length = {
                type = "integer"
                value = 1048576
                description = "Maximum string length (1MB)"
              }
              max_string_operations = {
                type = "integer"
                value = 1000
                description = "Maximum string operations per template"
              }
              string_cache_size = {
                type = "integer"
                value = 100
                description = "String operation cache size"
              }
            }
          }
          
          math_functions_constraints = {
            type = "object"
            description = "Math function constraints"
            properties = {
              max_number_precision = {
                type = "integer"
                value = 15
                description = "Maximum number precision"
              }
              max_math_operations = {
                type = "integer"
                value = 1000
                description = "Maximum math operations per template"
              }
              math_cache_size = {
                type = "integer"
                value = 50
                description = "Math operation cache size"
              }
            }
          }
          
          array_functions_constraints = {
            type = "object"
            description = "Array function constraints"
            properties = {
              max_array_size = {
                type = "integer"
                value = 10000
                description = "Maximum array size"
              }
              max_array_operations = {
                type = "integer"
                value = 500
                description = "Maximum array operations per template"
              }
              array_cache_size = {
                type = "integer"
                value = 50
                description = "Array operation cache size"
              }
            }
          }
          
          system_functions_constraints = {
            type = "object"
            description = "System function constraints"
            properties = {
              system_function_timeout = {
                type = "integer"
                value = 5000
                description = "System function timeout (5s)"
              }
              system_function_memory = {
                type = "integer"
                value = 10485760
                description = "System function memory limit (10MB)"
              }
              system_function_cpu = {
                type = "number"
                value = 0.1
                description = "System function CPU limit (10%)"
              }
            }
          }
        }
      }
      
      # Engine-specific constraints
      engine_function_constraints = {
        type = "object"
        description = "Engine-specific function constraints"
        properties = {
          go_template_constraints = {
            type = "object"
            description = "Go template function constraints"
            properties = {
              max_functions = {
                type = "integer"
                value = 50
                description = "Maximum functions for Go templates"
              }
              max_nesting = {
                type = "integer"
                value = 10
                description = "Maximum function nesting for Go templates"
              }
              builtin_functions = {
                type = "array"
                value = ["len", "index", "printf", "html", "js", "call"]
                description = "Built-in Go template functions"
              }
            }
          }
          
          jinja2_constraints = {
            type = "object"
            description = "Jinja2 function constraints"
            properties = {
              max_functions = {
                type = "integer"
                value = 75
                description = "Maximum functions for Jinja2 templates"
              }
              max_nesting = {
                type = "integer"
                value = 15
                description = "Maximum function nesting for Jinja2 templates"
              }
              builtin_functions = {
                type = "array"
                value = ["range", "dict", "cycler", "joiner", "namespace"]
                description = "Built-in Jinja2 functions"
              }
            }
          }
          
          handlebars_constraints = {
            type = "object"
            description = "Handlebars function constraints"
            properties = {
              max_functions = {
                type = "integer"
                value = 30
                description = "Maximum functions for Handlebars templates"
              }
              max_nesting = {
                type = "integer"
                value = 8
                description = "Maximum function nesting for Handlebars templates"
              }
              builtin_functions = {
                type = "array"
                value = ["if", "unless", "each", "with", "lookup"]
                description = "Built-in Handlebars functions"
              }
            }
          }
        }
      }
    }
  }
  
  # Function-specific extensions
  function_extensions = {
    type = "object"
    description = "Function-specific extensions and enhancements"
    
    properties = {
      # Custom function extensions
      custom_function_extensions = {
        type = "object"
        description = "Custom function capabilities"
        properties = {
          function_registration = {
            type = "object"
            description = "Custom function registration"
            properties = {
              allow_custom_functions = {
                type = "boolean"
                value = true
                description = "Allow registration of custom functions"
              }
              custom_function_validation = {
                type = "boolean"
                value = true
                description = "Validate custom function signatures"
              }
              custom_function_sandboxing = {
                type = "boolean"
                value = true
                description = "Sandbox custom function execution"
              }
              custom_function_caching = {
                type = "boolean"
                value = true
                description = "Cache custom function results"
              }
            }
          }
          
          function_plugins = {
            type = "object"
            description = "Function plugin system"
            properties = {
              plugin_support = {
                type = "boolean"
                value = true
                description = "Support for function plugins"
              }
              plugin_validation = {
                type = "boolean"
                value = true
                description = "Validate function plugins"
              }
              plugin_isolation = {
                type = "boolean"
                value = true
                description = "Isolate function plugins"
              }
              plugin_versioning = {
                type = "boolean"
                value = true
                description = "Support plugin versioning"
              }
            }
          }
        }
      }
      
      # Function optimization extensions
      function_optimization_extensions = {
        type = "object"
        description = "Function optimization capabilities"
        properties = {
          function_inlining = {
            type = "object"
            description = "Function inlining optimization"
            properties = {
              enable_inlining = {
                type = "boolean"
                value = true
                description = "Enable function inlining"
              }
              inline_threshold = {
                type = "integer"
                value = 100
                description = "Function size threshold for inlining"
              }
              inline_depth_limit = {
                type = "integer"
                value = 3
                description = "Maximum inlining depth"
              }
            }
          }
          
          function_memoization = {
            type = "object"
            description = "Function memoization optimization"
            properties = {
              enable_memoization = {
                type = "boolean"
                value = true
                description = "Enable function memoization"
              }
              memoization_cache_size = {
                type = "integer"
                value = 1000
                description = "Memoization cache size"
              }
              memoization_ttl = {
                type = "string"
                value = "1h"
                description = "Memoization cache TTL"
              }
            }
          }
          
          function_parallelization = {
            type = "object"
            description = "Function parallelization optimization"
            properties = {
              enable_parallelization = {
                type = "boolean"
                value = false
                description = "Enable function parallelization"
              }
              max_parallel_functions = {
                type = "integer"
                value = 4
                description = "Maximum parallel functions"
              }
              parallel_timeout = {
                type = "integer"
                value = 10000
                description = "Parallel execution timeout"
              }
            }
          }
        }
      }
      
      # Function monitoring extensions
      function_monitoring_extensions = {
        type = "object"
        description = "Function monitoring and observability"
        properties = {
          function_metrics = {
            type = "object"
            description = "Function metrics collection"
            properties = {
              collect_execution_time = {
                type = "boolean"
                value = true
                description = "Collect function execution time"
              }
              collect_memory_usage = {
                type = "boolean"
                value = true
                description = "Collect function memory usage"
              }
              collect_call_count = {
                type = "boolean"
                value = true
                description = "Collect function call count"
              }
              collect_error_count = {
                type = "boolean"
                value = true
                description = "Collect function error count"
              }
              collect_cache_hits = {
                type = "boolean"
                value = true
                description = "Collect function cache hits"
              }
            }
          }
          
          function_profiling = {
            type = "object"
            description = "Function profiling capabilities"
            properties = {
              enable_profiling = {
                type = "boolean"
                value = false
                description = "Enable function profiling"
              }
              profile_sampling_rate = {
                type = "number"
                value = 0.1
                description = "Function profiling sampling rate"
              }
              profile_output_format = {
                type = "string"
                value = "json"
                enum = ["json", "text", "pprof"]
                description = "Profile output format"
              }
            }
          }
          
          function_tracing = {
            type = "object"
            description = "Function tracing capabilities"
            properties = {
              enable_tracing = {
                type = "boolean"
                value = false
                description = "Enable function tracing"
              }
              trace_sampling_rate = {
                type = "number"
                value = 0.01
                description = "Function tracing sampling rate"
              }
              trace_output_format = {
                type = "string"
                value = "jaeger"
                enum = ["jaeger", "zipkin", "otlp"]
                description = "Trace output format"
              }
            }
          }
        }
      }
      
      # Function debugging extensions
      function_debugging_extensions = {
        type = "object"
        description = "Function debugging capabilities"
        properties = {
          function_debugging = {
            type = "object"
            description = "Function debugging features"
            properties = {
              enable_debugging = {
                type = "boolean"
                value = false
                description = "Enable function debugging"
              }
              debug_log_level = {
                type = "string"
                value = "info"
                enum = ["debug", "info", "warn", "error"]
                description = "Function debug log level"
              }
              debug_output_format = {
                type = "string"
                value = "text"
                enum = ["text", "json", "structured"]
                description = "Debug output format"
              }
            }
          }
          
          function_testing = {
            type = "object"
            description = "Function testing capabilities"
            properties = {
              enable_testing = {
                type = "boolean"
                value = true
                description = "Enable function testing"
              }
              test_coverage = {
                type = "boolean"
                value = true
                description = "Collect function test coverage"
              }
              test_timeout = {
                type = "integer"
                value = 5000
                description = "Function test timeout"
              }
            }
          }
        }
      }
    }
  }
  
  # Function-specific metadata
  function_metadata = {
    type = "object"
    description = "Function-specific metadata and information"
    
    properties = {
      # Function versioning
      versioning = {
        type = "object"
        description = "Function versioning information"
        properties = {
          function_version = {
            type = "string"
            value = "1.0.0"
            description = "Function schema version"
          }
          api_version = {
            type = "string"
            value = "1.0.0"
            description = "Function API version"
          }
          compatibility_version = {
            type = "string"
            value = "1.0.0"
            description = "Function compatibility version"
          }
        }
      }
      
      # Function lifecycle
      lifecycle = {
        type = "object"
        description = "Function lifecycle information"
        properties = {
          created_at = {
            type = "string"
            required = false
            format = "date-time"
            description = "Function creation timestamp"
          }
          updated_at = {
            type = "string"
            required = false
            format = "date-time"
            description = "Function last update timestamp"
          }
          deprecated_at = {
            type = "string"
            required = false
            format = "date-time"
            description = "Function deprecation timestamp"
          }
          removed_at = {
            type = "string"
            required = false
            format = "date-time"
            description = "Function removal timestamp"
          }
        }
      }
      
      # Function documentation
      documentation = {
        type = "object"
        description = "Function documentation"
        properties = {
          description = {
            type = "string"
            required = false
            description = "Function description"
          }
          usage_examples = {
            type = "array"
            required = false
            description = "Function usage examples"
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
                  type = "string"
                  required = false
                  description = "Example input"
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
            description = "Function API reference"
            additional_properties = true
          }
        }
      }
      
      # Function tags
      tags = {
        type = "array"
        required = false
        description = "Function tags"
        items = {
          type = "string"
        }
      }
    }
  }
}