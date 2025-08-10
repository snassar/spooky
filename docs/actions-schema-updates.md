# Actions Schema Updates Required

## Overview

Based on the comprehensive script execution guide, the current `actions.schema.hcl` needs several updates to align with the new file-based script approach and remove support for inline scripts.

## Current Schema Issues

### 1. Script Field Definition

**Current (Problematic):**
```hcl
script = {
  type = "string"
  required = false
  pattern = "^[a-zA-Z0-9/._-]+$"
  description = "Script file path (for script type)"
}
```

**Issues:**
- Pattern is too restrictive for file paths
- Doesn't enforce file-based approach
- Allows any string content (could be inline scripts)

### 2. Missing Variables Support

**Missing:**
- No `variables` field for templated scripts
- No validation for script file existence
- No distinction between static and templated scripts

### 3. Validation Rules

**Current (Incomplete):**
```hcl
script_type = {
  rule = "conditional"
  condition = "type == 'script' && script != null"
  message = "Script type actions must have a script path specified"
}
```

**Issues:**
- Doesn't validate file path format
- Doesn't check if file exists
- Doesn't distinguish between static and templated scripts

## Required Schema Updates

### 1. Enhanced Script Field Definition

```hcl
script = {
  type = "string"
  required = false
  pattern = "^(files|templates)/[a-zA-Z0-9/._-]+(\\.sh|\\.tmpl)?$"
  description = "Script file path in files/ or templates/ directory (for script type)"
  validation = {
    pattern = "^(files|templates)/[a-zA-Z0-9/._-]+(\\.sh|\\.tmpl)?$"
    message = "Script must reference a file in files/ or templates/ directory"
  }
}
```

### 2. Add Variables Field

```hcl
variables = {
  type = "object"
  required = false
  description = "Variables for templated scripts (for script type with .tmpl files)"
  additional_properties = "string"
  validation = {
    conditional = "script ends with .tmpl"
    message = "Variables are only valid for templated scripts (.tmpl files)"
  }
}
```

### 3. Enhanced Validation Rules

```hcl
# Script type validation
script_type = {
  rule = "conditional"
  condition = "type == 'script' && script != null"
  message = "Script type actions must have a script path specified"
}

# Script file path validation
script_file_path = {
  rule = "regex"
  pattern = "^(files|templates)/[a-zA-Z0-9/._-]+(\\.sh|\\.tmpl)?$"
  message = "Script must reference a file in files/ or templates/ directory"
}

# Script file existence validation
script_file_exists = {
  rule = "file_exists"
  message = "Referenced script file does not exist"
}

# Template variables validation
template_variables = {
  rule = "conditional"
  condition = "script ends with .tmpl && variables != null"
  message = "Templated scripts (.tmpl) should have variables defined"
}

# Static script validation
static_script = {
  rule = "conditional"
  condition = "script starts with 'files/' && !script ends with .tmpl"
  message = "Static scripts should be in files/ directory without .tmpl extension"
}
```

### 4. Updated Type Enum

```hcl
type = {
  type = "string"
  required = true
  enum = ["command", "script", "template_deploy", "template_evaluate", "template_validate", "template_cleanup", "file_copy", "service_control"]
  description = "Action execution type"
}
```

## Complete Updated Schema Section

```hcl
# Action definitions
action "action_name" {
  description = {
    type = "string"
    required = true
    min_length = 1
    max_length = 500
    description = "Action description"
  }
  
  type = {
    type = "string"
    required = true
    enum = ["command", "script", "template_deploy", "template_evaluate", "template_validate", "template_cleanup", "file_copy", "service_control"]
    description = "Action execution type"
  }
  
  command = {
    type = "string"
    required = false
    min_length = 1
    max_length = 1000
    description = "Command to execute (for command type)"
    validation = {
      pattern = "^(?!.*[;&|`$]).*$"
      message = "Command cannot contain shell operators or special characters"
    }
  }
  
  script = {
    type = "string"
    required = false
    pattern = "^(files|templates)/[a-zA-Z0-9/._-]+(\\.sh|\\.tmpl)?$"
    description = "Script file path in files/ or templates/ directory (for script type)"
    validation = {
      pattern = "^(files|templates)/[a-zA-Z0-9/._-]+(\\.sh|\\.tmpl)?$"
      message = "Script must reference a file in files/ or templates/ directory"
    }
  }
  
  variables = {
    type = "object"
    required = false
    description = "Variables for templated scripts (for script type with .tmpl files)"
    additional_properties = "string"
    validation = {
      conditional = "script ends with .tmpl"
      message = "Variables are only valid for templated scripts (.tmpl files)"
    }
  }
  
  # ... rest of existing fields ...
}
```

## Validation Rules Updates

```hcl
validation = {
  # ... existing validation rules ...
  
  # Enhanced script validation
  script_type = {
    rule = "conditional"
    condition = "type == 'script' && script != null"
    message = "Script type actions must have a script path specified"
  }
  
  script_file_path = {
    rule = "regex"
    pattern = "^(files|templates)/[a-zA-Z0-9/._-]+(\\.sh|\\.tmpl)?$"
    message = "Script must reference a file in files/ or templates/ directory"
  }
  
  script_file_exists = {
    rule = "file_exists"
    message = "Referenced script file does not exist"
  }
  
  template_variables = {
    rule = "conditional"
    condition = "script ends with .tmpl && variables != null"
    message = "Templated scripts (.tmpl) should have variables defined"
  }
  
  static_script = {
    rule = "conditional"
    condition = "script starts with 'files/' && !script ends with .tmpl"
    message = "Static scripts should be in files/ directory without .tmpl extension"
  }
  
  # ... rest of existing validation rules ...
}
```

## Migration Impact

### Breaking Changes
1. **Inline scripts no longer supported**: All script actions must reference files
2. **Stricter file path validation**: Must be in `files/` or `templates/` directories
3. **File existence validation**: Referenced files must actually exist

### Required Actions
1. **Update existing actions**: Move inline scripts to files
2. **Update schema validation**: Implement new validation rules
3. **Update documentation**: Reflect new file-based approach
4. **Update tests**: Ensure new validation rules work correctly

## Benefits of Schema Updates

### 1. Enforced Best Practices
- Prevents inline script antipattern
- Ensures proper file organization
- Validates file existence

### 2. Better Validation
- Clear error messages for invalid configurations
- Prevents common mistakes
- Guides users toward correct patterns

### 3. Future-Proof Design
- Supports templated scripts with variables
- Clear separation between static and dynamic scripts
- Extensible for future enhancements

### 4. Improved Developer Experience
- Clear expectations about script organization
- Better error messages
- Consistent patterns across projects

These schema updates will ensure that the actions configuration aligns with the comprehensive script execution guide and enforces the file-based approach from the schema level.
