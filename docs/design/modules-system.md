# Spooky Modules System Design

## Overview

The spooky modules system enables extensible validation, custom schemas, and reusable components across projects. This document outlines the architectural design, use cases, and implementation approaches for a modular extension system.

## Problem Statement

### Current Limitations

1. **Hardcoded Validators**: All validation logic is built into spooky core
2. **Schema Duplication**: Teams reinvent validation logic across projects
3. **Limited Extensibility**: No way to add domain-specific validation rules
4. **No Reusability**: Custom validation patterns can't be shared between projects

### Desired Capabilities

1. **Custom Validators**: Project-specific validation rules
2. **Schema Sharing**: Reusable schema definitions
3. **Domain-Specific Logic**: Validation tailored to specific use cases
4. **Extensibility**: Easy addition of new validation capabilities

## Use Cases

### 1. Custom Action Validation

Teams building custom automation need validation specific to their domain:

```hcl
# schemas/custom-deployment.schema.hcl
schema "custom-deployment" {
  validation {
    fields {
      field "ssl_cert_path" {
        custom_validator = "ssl_cert_valid"
      }
      
      field "deployment_env" {
        custom_validator = "deployment_environment"
      }
    }
  }
}
```

### 2. Cross-Project Schema Reuse

Multiple projects with similar requirements:

```
project-a/
├── schemas/
│   └── deployment.schema.hcl  # Custom deployment validation

project-b/
├── schemas/
│   └── deployment.schema.hcl  # Same validation logic, duplicated
```

### 3. Domain-Specific Validation

Different teams have different validation needs:

- **Database Team**: Migration script validation
- **Security Team**: Compliance requirement validation  
- **DevOps Team**: Deployment configuration validation

## Architectural Approaches

### Approach 1: Full Module System

#### Module Structure
```
spooky-modules/
├── ssl-validation/
│   ├── spooky-module.hcl
│   ├── validators/
│   │   └── ssl_cert_valid.go
│   └── schemas/
│       └── ssl-cert.schema.hcl
├── deployment/
│   ├── spooky-module.hcl
│   ├── validators/
│   │   ├── deployment_environment.go
│   │   └── resource_limits.go
│   └── schemas/
│       └── deployment.schema.hcl
└── database/
    ├── spooky-module.hcl
    ├── validators/
    │   └── migration_safe.go
    └── schemas/
        └── migration.schema.hcl
```

#### Module Definition
```hcl
# ssl-validation/spooky-module.hcl
module "ssl-validation" {
  name = "ssl-validation"
  version = "1.0.0"
  description = "SSL certificate validation module"
  
  validators {
    validator "ssl_cert_valid" {
      description = "Validates SSL certificate format and expiration"
      function = "ssl_cert_valid"
    }
  }
  
  schemas {
    schema "ssl-cert" {
      file = "schemas/ssl-cert.schema.hcl"
    }
  }
}
```

#### Module Loading System
```go
type ModuleSystem struct {
    modules map[string]*Module
    validators map[string]ValidatorFunction
    schemas map[string]*Schema
}

type Module struct {
    Name        string
    Version     string
    Description string
    Validators  map[string]ValidatorFunction
    Schemas     map[string]*Schema
    Path        string
}

type ValidatorFunction func(data interface{}) error
```

#### Pros
- Full separation of concerns
- Version management
- Dependency resolution
- Reusable across projects

#### Cons
- High complexity
- Significant development effort
- Overhead for simple use cases

### Approach 2: Built-in Validator Registry

#### Built-in Validators
```go
var BuiltInValidators = map[string]ValidatorFunction{
    "required":        validateRequired,
    "email":          validateEmail,
    "url":            validateURL,
    "positive":       validatePositive,
    "file_exists":    validateFileExists,
    "ssl_cert_valid": validateSSLCert,
    "port_available": validatePortAvailable,
    "deployment_env": validateDeploymentEnvironment,
    "resource_limits": validateResourceLimits,
}
```

#### Project-Specific Validators
```go
// schemas/validators/custom.go
package validators

func ValidateDeploymentEnvironment(data interface{}) error {
    // Custom validation logic
    return nil
}

func ValidateResourceLimits(data interface{}) error {
    // Custom validation logic
    return nil
}
```

#### Schema Usage
```hcl
# schemas/custom-deployment.schema.hcl
schema "custom-deployment" {
  validation {
    fields {
      field "ssl_cert_path" {
        type = "string"
        custom_validator = "ssl_cert_valid"  # Built-in validator
      }
      
      field "deployment_env" {
        type = "string"
        custom_validator = "deployment_environment"  # Project-specific
      }
    }
  }
}
```

#### Pros
- Simpler implementation
- Faster development
- Lower complexity
- Immediate value

#### Cons
- Limited reusability
- No version management
- Harder to share across projects

### Approach 3: Hybrid Approach

#### Phase 1: Built-in Registry
Start with built-in validators and project-specific validators.

#### Phase 2: Simple Module System
Add basic module loading for sharing common validators.

#### Phase 3: Full Module System
Implement complete module system with dependency management.

## Recommended Implementation

### Phase 1: Built-in Validator Registry (Immediate)

1. **Extend Current Validator System**
   ```go
   func (v *Validator) executeCustomValidator(validator string, data interface{}) error {
       // Check built-in validators first
       if fn, exists := BuiltInValidators[validator]; exists {
           return fn(data)
       }
       
       // Check project-specific validators
       if fn, exists := v.projectValidators[validator]; exists {
           return fn(data)
       }
       
       return fmt.Errorf("unknown validator: %s", validator)
   }
   ```

2. **Add Common Built-in Validators**
   - File system validators (file_exists, dir_exists, permissions)
   - Network validators (url_accessible, port_available)
   - Security validators (ssl_cert_valid, key_permissions)
   - Business logic validators (deployment_environment, resource_limits)

3. **Project-Specific Validator Loading**
   ```go
   func (v *Validator) LoadProjectValidators(projectPath string) error {
       validatorsDir := filepath.Join(projectPath, "schemas", "validators")
       return v.loadValidatorsFromDirectory(validatorsDir)
   }
   ```

### Phase 2: Schema Loading Enhancement (Short-term)

1. **Dynamic Schema Discovery**
   ```go
   func (v *Validator) LoadProjectSchemas(projectPath string) error {
       schemasDir := filepath.Join(projectPath, "schemas")
       
       files, err := filepath.Glob(filepath.Join(schemasDir, "*.schema.hcl"))
       if err != nil {
           return err
       }
       
       for _, file := range files {
           schema, err := v.LoadSchema(file)
           if err != nil {
               return fmt.Errorf("failed to load schema %s: %w", file, err)
           }
           
           schemaName := filepath.Base(file)
           schemaName = strings.TrimSuffix(schemaName, ".schema.hcl")
           v.schemas[schemaName] = schema
       }
       
       return nil
   }
   ```

2. **Schema Reference in Actions**
   ```hcl
   # actions.hcl
   actions {
     action "deploy" {
       schema = "custom-deployment"  # References schemas/custom-deployment.schema.hcl
     }
   }
   ```

### Phase 3: Module System (Long-term)

1. **Module Definition Format**
2. **Module Discovery and Loading**
3. **Version Management**
4. **Dependency Resolution**

## Implementation Details

### Validator Function Signature
```go
type ValidatorFunction func(data interface{}) error

// Example implementations
func validateSSLCert(data interface{}) error {
    if certPath, ok := data.(string); ok {
        if !isValidSSLCertificate(certPath) {
            return fmt.Errorf("invalid SSL certificate: %s", certPath)
        }
    }
    return nil
}

func validateDeploymentEnvironment(data interface{}) error {
    if env, ok := data.(string); ok {
        if env == "production" && !hasApproval() {
            return fmt.Errorf("production deployment requires approval")
        }
    }
    return nil
}
```

### Schema Loading Integration
```go
func (v *Validator) ValidateAction(action *Action) error {
    // Check if action references a custom schema
    if action.Schema != "" {
        schema, exists := v.schemas[action.Schema]
        if !exists {
            return fmt.Errorf("schema not found: %s", action.Schema)
        }
        
        return v.validateWithSchema(action, schema)
    }
    
    // Use default validation
    return v.validateActionDefault(action)
}
```

### Error Handling
```go
type ValidationError struct {
    Field   string
    Message string
    Code    string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}
```

## Benefits

### 1. **Configuration Safety**
- Prevents invalid configurations from being deployed
- Catches common mistakes early in development
- Ensures consistent configuration patterns

### 2. **Developer Experience**
- Immediate feedback on configuration errors
- Self-documenting schemas
- Clear validation error messages

### 3. **Operational Reliability**
- Prevents runtime failures due to configuration issues
- Validates security-sensitive configurations
- Ensures compliance requirements

### 4. **Team Productivity**
- Reusable validation patterns
- Consistent standards across projects
- Reduced configuration errors

## Migration Path

### Current State
- Basic schema validation
- Hardcoded validators
- No extensibility

### Phase 1 (Immediate)
- Built-in validator registry
- Project-specific validators
- Enhanced error messages

### Phase 2 (Short-term)
- Dynamic schema loading
- Schema references in actions
- Improved validation coverage

### Phase 3 (Long-term)
- Full module system
- Version management
- Cross-project sharing

## Conclusion

The modules system provides a path for spooky to evolve from a simple configuration validator to a comprehensive, extensible validation platform. Starting with built-in validators and project-specific extensions provides immediate value while building toward a full module system.

The key insight is that validation needs vary significantly across different teams and use cases, and a one-size-fits-all approach limits spooky's effectiveness. A modular approach allows teams to define their own validation rules while maintaining the safety and consistency that comes with proper validation.

## Next Steps

1. **Implement Phase 1**: Built-in validator registry
2. **Add Common Validators**: File, network, security validators
3. **Enhance Schema Loading**: Dynamic schema discovery
4. **Document Examples**: Provide clear examples of custom validation
5. **Gather Feedback**: Understand real-world validation needs
6. **Plan Phase 2**: Design schema loading enhancements
7. **Evaluate Module System**: Assess need for full module system
