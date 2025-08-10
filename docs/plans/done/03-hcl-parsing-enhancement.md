# Implementation Plan: HCL Parsing Enhancement

## Overview
Enhance HCL parsing capabilities to support action script resolution, template path validation, and improved validation for action execution flow.

## Task Details
- **Task ID**: 8.3
- **Priority**: Medium
- **Files**: 
  - `internal/actions/acting/actor.go`
  - `internal/cli/commands/actions.go`
  - `internal/schemas/schemas/actions.schema.hcl`
- **Functions**: Script resolution, path validation, conditional validation, enhanced error handling

## Current State Analysis

### Existing Patterns
1. **Action Parsing**: Basic HCL parsing using HCL v2 library with schema validation
2. **Script Actions**: Currently expect inline script content, not file references
3. **Template Actions**: Reference template files but don't validate paths exist
4. **Validation**: Schema-driven validation using comprehensive schemas
5. **Error Handling**: Basic error handling for parsing failures

### Current Implementation Status
1. **Schema Validation**: ✅ Already implemented and used throughout the codebase
2. **Basic HCL Parsing**: ✅ Implemented with proper HCL v2 library usage
3. **Action Execution**: ✅ Basic action execution framework exists
4. **Script Resolution**: ❌ Only handles inline script content (needs file-based support)
5. **Template Path Validation**: ❌ References files but doesn't validate existence
6. **File Copy Actions**: ❌ Schema mentions but not fully implemented
7. **Conditional Validation**: ❌ Cross-field validation not enforced

### Unaddressed TODOs
1. **Script Actions**: Need to support file-based scripts in addition to inline content
2. **Template Paths**: Need to validate template source files exist
3. **File Copy Actions**: Need to implement path validation for source/destination
4. **Conditional Logic**: Need to enforce validation rules like "variables only for .tmpl files"

## Implementation Requirements

### Interface Compliance
The HCL parsing enhancement must:
1. **Support file-based scripts** in addition to inline script content
2. **Validate template source paths** exist before execution
3. **Implement file copy actions** with proper path validation
4. **Enforce conditional validation** rules in action schemas
5. **Enhance error handling** with better context for validation failures
6. **Maintain backward compatibility** with existing inline script actions

### Required Dependencies
- HCL v2 library (already imported)
- Actions system (already implemented)
- Schema validation system (already implemented)
- File system operations (already available)

## Detailed Implementation Plan

### Step 1: Enhance Script Action Support

**File**: `internal/actions/acting/actor.go`

```go
// Enhanced script resolution to support both inline and file-based scripts
func (a *actorImpl) resolveScript() (string, error) {
    if a.action.Script == "" {
        return "", fmt.Errorf("script field is required")
    }

    // Check if script references a file
    if strings.HasPrefix(a.action.Script, "files/") || 
       strings.HasPrefix(a.action.Script, "templates/") {
        
        // Validate file path format
        if !strings.HasPrefix(a.action.Script, "files/") && 
           !strings.HasPrefix(a.action.Script, "templates/") {
            return "", fmt.Errorf("script must reference a file in files/ or templates/ directory")
        }
        
        // Check if it's a template file
        if strings.HasSuffix(a.action.Script, ".tmpl") {
            return a.renderTemplateScript(a.action.Script, a.action.Variables)
        } else {
            // Static script - read file content
            return a.readScriptFile(a.action.Script)
        }
    }
    
    // Inline script content (existing behavior)
    return a.action.Script, nil
}

// Read script file from files/ directory
func (a *actorImpl) readScriptFile(scriptPath string) (string, error) {
    // Resolve path relative to project root
    fullPath := filepath.Join(a.projectPath, scriptPath)
    
    // Validate file exists
    if _, err := os.Stat(fullPath); os.IsNotExist(err) {
        return "", fmt.Errorf("script file does not exist: %s", scriptPath)
    }
    
    // Read file content
    content, err := os.ReadFile(fullPath)
    if err != nil {
        return "", fmt.Errorf("failed to read script file %s: %w", scriptPath, err)
    }
    
    return string(content), nil
}

// Render template script with variables
func (a *actorImpl) renderTemplateScript(templatePath string, variables map[string]string) (string, error) {
    // Load template from templates/ directory
    templateContent, err := a.readScriptFile(templatePath)
    if err != nil {
        return "", fmt.Errorf("failed to load template %s: %w", templatePath, err)
    }
    
    // Render template with variables
    rendered, err := a.templateEngine.Render(templateContent, variables)
    if err != nil {
        return "", fmt.Errorf("failed to render template %s: %w", templatePath, err)
    }
    
    return rendered, nil
}
```

### Step 2: Implement Template Path Validation

**File**: `internal/actions/acting/actor.go`

```go
// Enhanced template validation with path existence check
func (a *actorImpl) validateTemplateAction() error {
    if a.action.Template == nil {
        return fmt.Errorf("template configuration is required")
    }

    if a.action.Template.Source == "" {
        return fmt.Errorf("template source is required")
    }

    if a.action.Template.Destination == "" {
        return fmt.Errorf("template destination is required")
    }

    // Validate template source file exists
    if err := a.validateTemplateSource(a.action.Template.Source); err != nil {
        return fmt.Errorf("template source validation failed: %w", err)
    }

    return nil
}

// Validate template source file exists
func (a *actorImpl) validateTemplateSource(sourcePath string) error {
    // Must reference templates/ directory
    if !strings.HasPrefix(sourcePath, "templates/") {
        return fmt.Errorf("template source must reference a file in templates/ directory")
    }
    
    // Must have .tmpl extension
    if !strings.HasSuffix(sourcePath, ".tmpl") {
        return fmt.Errorf("template source must have .tmpl extension")
    }
    
    // Check if file exists
    fullPath := filepath.Join(a.projectPath, sourcePath)
    if _, err := os.Stat(fullPath); os.IsNotExist(err) {
        return fmt.Errorf("template source file does not exist: %s", sourcePath)
    }
    
    return nil
}
```

### Step 3: Implement File Copy Actions

**File**: `internal/actions/acting/actor.go`

```go
// Act file copy action
func (a *actorImpl) actFileCopy(ctx context.Context, context *spookytypesactions.ActionContext, result *spookytypesactions.RunResult) error {
    a.logger.Debug("Executing file copy action",
        spookylogging.String("action", a.action.Name),
        spookylogging.String("source", a.action.FileCopy.Source),
        spookylogging.String("destination", a.action.FileCopy.Destination))

    // 1. Validate file copy configuration
    if err := a.validateFileCopyAction(); err != nil {
        result.Status = spookytypesactions.RunStatusFailed
        result.Error = fmt.Sprintf("file copy validation failed: %v", err)
        return fmt.Errorf("file copy validation failed: %w", err)
    }

    // 2. Validate source file exists
    if err := a.validateFileCopySource(a.action.FileCopy.Source); err != nil {
        result.Status = spookytypesactions.RunStatusFailed
        result.Error = fmt.Sprintf("source file validation failed: %v", err)
        return fmt.Errorf("source file validation failed: %w", err)
    }

    // 3. Execute file copy
    if a.parallel {
        return a.actFileCopyParallel(ctx, context, result)
    } else {
        return a.actFileCopySequential(ctx, context, result)
    }
}

// Validate file copy action configuration
func (a *actorImpl) validateFileCopyAction() error {
    if a.action.FileCopy == nil {
        return fmt.Errorf("file copy configuration is required")
    }

    if a.action.FileCopy.Source == "" {
        return fmt.Errorf("file copy source is required")
    }

    if a.action.FileCopy.Destination == "" {
        return fmt.Errorf("file copy destination is required")
    }

    return nil
}

// Validate file copy source file exists
func (a *actorImpl) validateFileCopySource(sourcePath string) error {
    // Must reference files/ directory
    if !strings.HasPrefix(sourcePath, "files/") {
        return fmt.Errorf("file copy source must reference a file in files/ directory")
    }
    
    // Check if file exists
    fullPath := filepath.Join(a.projectPath, sourcePath)
    if _, err := os.Stat(fullPath); os.IsNotExist(err) {
        return fmt.Errorf("file copy source file does not exist: %s", sourcePath)
    }
    
    return nil
}
```

### Step 4: Enhance Conditional Validation

**File**: `internal/schemas/schemas/actions.schema.hcl`

```hcl
# Enhanced script action with conditional validation
script = {
  type = "string"
  required = false
  description = "Script file path in files/ or templates/ directory, or inline script content"
  validation = {
    # If script ends with .tmpl, variables must be provided
    conditional = "script ends with .tmpl implies variables is not empty"
    message = "Variables are required for templated scripts (.tmpl files)"
  }
}

variables = {
  type = "object"
  required = false
  description = "Variables for templated scripts (for script type with .tmpl files)"
  additional_properties = "string"
  validation = {
    # Variables are only valid for templated scripts
    conditional = "script ends with .tmpl"
    message = "Variables are only valid for templated scripts (.tmpl files)"
  }
}

# Enhanced file copy configuration
file_copy = {
  type = "object"
  required = false
  description = "File copy configuration (for file_copy type)"
  
  properties = {
    source = {
      type = "string"
      required = true
      pattern = "^files/[a-zA-Z0-9/._-]+$"
      description = "Source file path in files/ directory"
    }
    
    destination = {
      type = "string"
      required = true
      pattern = "^[a-zA-Z0-9/._-]+$"
      description = "Destination file path on target machine"
    }
    
    backup = {
      type = "boolean"
      required = false
      default = false
      description = "Create backup of existing file before overwriting"
    }
    
    permissions = {
      type = "string"
      required = false
      pattern = "^[0-7]{3,4}$"
      description = "File permissions (octal format)"
    }
    
    owner = {
      type = "string"
      required = false
      pattern = "^[a-zA-Z0-9._-]+$"
      description = "File owner (username)"
    }
    
    group = {
      type = "string"
      required = false
      pattern = "^[a-zA-Z0-9._-]+$"
      description = "File group (group name)"
    }
  }
}
```

### Step 5: Enhance Error Handling and Validation

**File**: `internal/actions/acting/actor.go`

```go
// Enhanced action validation with better error context
func (a *actorImpl) validateAction() error {
    var errors []string

    // Basic validation
    if a.action.Name == "" {
        errors = append(errors, "action name is required")
    }

    if a.action.Type == "" {
        errors = append(errors, "action type is required")
    }

    // Type-specific validation
    switch a.action.Type {
    case "command":
        if err := a.validateCommandAction(); err != nil {
            errors = append(errors, fmt.Sprintf("command validation: %v", err))
        }
    case "script":
        if err := a.validateScriptAction(); err != nil {
            errors = append(errors, fmt.Sprintf("script validation: %v", err))
        }
    case "template_deploy":
        if err := a.validateTemplateAction(); err != nil {
            errors = append(errors, fmt.Sprintf("template validation: %v", err))
        }
    case "file_copy":
        if err := a.validateFileCopyAction(); err != nil {
            errors = append(errors, fmt.Sprintf("file copy validation: %v", err))
        }
    default:
        errors = append(errors, fmt.Sprintf("unsupported action type: %s", a.action.Type))
    }

    // Conditional validation
    if err := a.validateConditionalRules(); err != nil {
        errors = append(errors, fmt.Sprintf("conditional validation: %v", err))
    }

    if len(errors) > 0 {
        return fmt.Errorf("action validation failed:\n  %s", strings.Join(errors, "\n  "))
    }

    return nil
}

// Validate conditional rules between fields
func (a *actorImpl) validateConditionalRules() error {
    // Variables are only valid for templated scripts
    if a.action.Variables != nil && len(a.action.Variables) > 0 {
        if a.action.Type != "script" {
            return fmt.Errorf("variables are only valid for script actions")
        }
        
        if !strings.HasSuffix(a.action.Script, ".tmpl") {
            return fmt.Errorf("variables are only valid for templated scripts (.tmpl files)")
        }
    }

    // Script must reference files/ or templates/ if not inline
    if a.action.Type == "script" && a.action.Script != "" {
        if strings.HasPrefix(a.action.Script, "files/") || 
           strings.HasPrefix(a.action.Script, "templates/") {
            // File-based script - validate path format
            if !strings.HasPrefix(a.action.Script, "files/") && 
               !strings.HasPrefix(a.action.Script, "templates/") {
                return fmt.Errorf("script must reference a file in files/ or templates/ directory")
            }
        }
    }

    return nil
}
```

## Configuration Options

### Supported Options
- **FileBasedScriptsEnabled**: Enable/disable file-based script support
- **TemplatePathValidation**: Enable/disable template source path validation
- **FileCopyActionsEnabled**: Enable/disable file copy action type
- **ConditionalValidation**: Enable/disable cross-field validation rules

## Dependencies

### Internal Dependencies (Already Available)
- `spooky/internal/actions` ✅
- `spooky/internal/schemas` ✅
- `spooky/internal/logging` ✅
- `spooky/internal/templates` ✅

### External Dependencies (Already Imported)
- `github.com/hashicorp/hcl/v2` ✅
- `os` (standard library) ✅
- `path/filepath` (standard library) ✅
- `strings` (standard library) ✅

## Implementation Order

1. ✅ **Schema validation system** (already implemented)
2. ✅ **Basic HCL parsing** (already implemented)
3. ✅ **Action execution framework** (already implemented)
4. 🔄 **Enhance script resolution** (support file-based scripts)
5. 🔄 **Implement template path validation** (check source files exist)
6. 🔄 **Add file copy actions** (with path validation)
7. 🔄 **Implement conditional validation** (cross-field rules)
8. 🔄 **Enhance error handling** (better validation context)
9. 🔄 **Add comprehensive tests** (for new functionality)
10. 🔄 **Update documentation** (reflect new capabilities)

## Current Status Summary

- ✅ **Schema-driven validation**: Fully implemented and used throughout the codebase
- ✅ **Basic HCL parsing**: Implemented with proper HCL v2 library usage
- ✅ **Action execution**: Basic framework exists and working
- ✅ **Template actions**: Basic support exists
- ❌ **File-based scripts**: Only supports inline script content
- ❌ **Template path validation**: References files but doesn't validate existence
- ❌ **File copy actions**: Schema mentions but not implemented
- ❌ **Conditional validation**: Cross-field rules not enforced

## Next Steps

1. **Enhance script actions** to support file-based scripts in addition to inline content
2. **Implement template path validation** to ensure referenced files exist
3. **Add file copy actions** with proper source/destination validation
4. **Enforce conditional validation** rules between action fields
5. **Improve error messages** with better context for validation failures
6. **Add comprehensive tests** for new functionality
7. **Update documentation** to reflect new capabilities

## Why This Approach

This enhancement focuses on **action execution flow** rather than changing fundamental data structures. It addresses the actual needs identified:

- **Script Actions**: Need to reference files in `files/` or `templates/` directories
- **Template Actions**: Need to validate that referenced template files exist
- **File Copy Actions**: Need to validate source files exist before attempting copy
- **Conditional Logic**: Need to enforce validation rules like "variables only for .tmpl files"

This maintains backward compatibility while adding the file-based capabilities needed for better script and template management.
