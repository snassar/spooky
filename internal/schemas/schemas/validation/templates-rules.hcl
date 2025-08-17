# Templates Validation Rules
# Extracted from internal/schemas/schemas/structure/templates.hcl
# These rules validate schema compliance and data format correctness

# Templates validation rules
validation_rules {
  # Template ID validation
  template_id_validation {
    rule {
      name = "template_id_format"
      description = "Template ID must contain only alphanumeric characters, dots, underscores, and hyphens"
      condition = "template_id != null && !template_id.matches('^[a-zA-Z0-9._-]+$')"
      message = "Template ID must contain only alphanumeric characters, dots, underscores, and hyphens"
      severity = "error"
    }
  }
  
  # Source path validation
  source_path_validation {
    rule {
      name = "source_path_format"
      description = "Source path must be in templates/ directory with .tmpl extension"
      condition = "source_path != null && !source_path.matches('^templates/.*\\.tmpl$')"
      message = "Source path must be in templates/ directory with .tmpl extension"
      severity = "error"
    }
  }
  
  # Required fields validation
  required_fields_validation {
    rule {
      name = "template_type_required"
      description = "Template type is required"
      condition = "template_type == null"
      message = "Template type is required"
      severity = "error"
    }
    
    rule {
      name = "scope_required"
      description = "Template scope is required"
      condition = "scope == null"
      message = "Template scope is required"
      severity = "error"
    }
    
    rule {
      name = "security_level_required"
      description = "Security level is required"
      condition = "security_level == null"
      message = "Security level is required"
      severity = "error"
    }
    
    rule {
      name = "engine_required"
      description = "Template engine is required"
      condition = "engine == null"
      message = "Template engine is required"
      severity = "error"
    }
  }
  
  # Variable validation
  variable_validation {
    rule {
      name = "variable_name_format"
      description = "Variable names must start with letter or underscore and contain only alphanumeric characters and underscores"
      condition = "variables != null && !allVariableNamesValid(variables)"
      message = "Variable names must start with letter or underscore and contain only alphanumeric characters and underscores"
      severity = "error"
    }
  }
  
  # Function validation
  function_validation {
    rule {
      name = "function_name_format"
      description = "Function names must start with letter or underscore and contain only alphanumeric characters and underscores"
      condition = "functions != null && !allFunctionNamesValid(functions)"
      message = "Function names must start with letter or underscore and contain only alphanumeric characters and underscores"
      severity = "error"
    }
  }
  
  # Pattern validation
  pattern_validation {
    rule {
      name = "pattern_format"
      description = "Restricted patterns must be valid regex patterns"
      condition = "restricted_patterns != null && !allPatternsValid(restricted_patterns)"
      message = "Restricted patterns must be valid regex patterns"
      severity = "error"
    }
  }
  
  # Size validation
  size_validation {
    rule {
      name = "template_size_limits"
      description = "Template size must be between 1KB and 10MB"
      condition = "size != null && (size < 1024 || size > 10485760)"
      message = "Template size must be between 1KB and 10MB"
      severity = "warning"
    }
  }
  
  # Nesting validation
  nesting_validation {
    rule {
      name = "nesting_depth_limits"
      description = "Nesting depth must be between 1 and 50"
      condition = "nesting_depth != null && (nesting_depth < 1 || nesting_depth > 50)"
      message = "Nesting depth must be between 1 and 50"
      severity = "warning"
    }
  }
  
  # Run time validation
  runtime_validation {
    rule {
      name = "run_time_limits"
      description = "Run time must be between 100ms and 30s"
      condition = "run_time != null && (run_time < 100 || run_time > 30000)"
      message = "Run time must be between 100ms and 30s"
      severity = "warning"
    }
  }
  
  # Memory usage validation
  memory_validation {
    rule {
      name = "memory_usage_limits"
      description = "Memory usage must be between 1MB and 100MB"
      condition = "memory_usage != null && (memory_usage < 1048576 || memory_usage > 104857600)"
      message = "Memory usage must be between 1MB and 100MB"
      severity = "warning"
    }
  }
  
  # Security validation
  security_validation {
    rule {
      name = "no_circular_refs"
      description = "No circular references allowed in template definitions"
      condition = "hasCircularReferences(template_content)"
      message = "No circular references allowed in template definitions"
      severity = "error"
    }
    
    rule {
      name = "no_dangerous_patterns"
      description = "Dangerous patterns are not allowed in templates"
      condition = "hasDangerousPatterns(template_content)"
      message = "Dangerous patterns are not allowed in templates"
      severity = "error"
    }
  }
}
