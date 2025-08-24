# Variables Structure Schema
# Common variable structure definitions for all storage formats
# This file defines the structure of variables that can be stored
# Used by HCL and JSON storage schemas for variables

# Schema metadata
metadata {
  version = "1"
  description = "Variables configuration schema for spooky variable definitions"
}

# Variables block structure (like machines and actions)
variables {
  # Individual variable definitions
  variable {
    # Variable name (required)
    name = "variable_name"
    
    # Variable type (optional)
    type = "string"
    
    # Variable value (required)
    value = "variable_value"
    
    # Variable description (optional)
    description = "Variable description"
    
    # Sensitive flag (optional)
    sensitive = false
    
    # Encryption flag (optional)
    encrypted = false
  }
}
