# Variables System

## Overview

The Variables System provides comprehensive variable management and resolution capabilities for the spooky codebase. It enables variable definition, resolution, validation, dependency management, and integration with all system components for dynamic configuration and template processing.

**Status**: **Implemented** - Complete variables system with loading, resolution, validation, dependency management, and CLI integration.

## Related Systems

This system integrates with and depends on several other spooky systems:

- **[Templates System](TEMPLATES_SYSTEM.md)** - Variables are used in templates for dynamic content rendering
- **[Actions System](ACTIONS_SYSTEM.md)** - Variables are used by actions for dynamic configuration
- **[Logging System](LOGGING_SYSTEM.md)** - Variable resolution generates comprehensive logs for monitoring and debugging
- **[Projects System](PROJECTS_SYSTEM.md)** - Variables are organized within projects
- **[Integrations System](INTEGRATIONS_SYSTEM.md)** - Variables integrate with other systems through the IntegrationManager
- **[Facts System](FACTS_SYSTEM.md)** - Variables can reference facts for dynamic values
- **[Machines System](MACHINES_SYSTEM.md)** - Variables can reference machine metadata
- **[SSH System](SSH_SYSTEM.md)** - Variables can be used in SSH configurations

## Architecture

### Core Components

#### Variables Manager
- **File**: `internal/variables/manager.go`
- **Purpose**: Central variables management with loading and resolution
- **Features**:
  - Variable loading and parsing
  - Variable resolution and interpolation
  - Dependency management
  - Variable validation
  - Error handling and recovery
  - Performance optimization

#### Variables Loader
- **File**: `internal/variables/loader.go`
- **Purpose**: Variable configuration loading and parsing
- **Features**:
  - HCL configuration parsing
  - Variable file loading
  - Configuration validation
  - Error handling
  - Performance optimization
  - File monitoring

#### Variables Validator
- **File**: `internal/variables/validator.go`
- **Purpose**: Variable configuration and dependency validation
- **Features**:
  - Variable validation
  - Dependency validation
  - Circular dependency detection
  - Type validation
  - Constraint validation
  - Error reporting

#### Variables Integration
- **File**: `internal/variables/integration.go`
- **Purpose**: Interface implementation for system integration
- **Features**:
  - LoadVariables - Load variable configurations
  - ResolveVariables - Resolve variable values
  - ValidateVariables - Validate variable configurations
  - GetVariable - Get specific variable value

### Integration Points

#### Templates Integration
- Provides variables for template rendering
- Supports variable interpolation in templates
- Enables dynamic template content generation

#### Actions Integration
- Provides variables for action running
- Supports variable substitution in actions
- Enables dynamic action configuration

#### Configuration Integration
- Provides variables for configuration management
- Supports variable-based configuration
- Enables dynamic configuration generation

#### Secrets Integration
- Provides encrypted variables for secure data
- Supports variable encryption and decryption
- Enables secure variable management

## Variables Types

### Variable Structure
```go
type Variable struct {
    ID              string                 // Variable identifier
    Name            string                 // Variable name
    Value           interface{}            // Variable value
    Type            string                 // Variable type
    Description     string                 // Variable description
    Default         interface{}            // Default value
    Required        bool                   // Required flag
    Sensitive       bool                   // Sensitive flag
    Encrypted       bool                   // Encrypted flag
    Dependencies    []string               // Variable dependencies
    Constraints     *VariableConstraints   // Variable constraints
    Metadata        map[string]interface{} // Variable metadata
    CreatedAt       time.Time              // Creation timestamp
    UpdatedAt       time.Time              // Update timestamp
}
```

### Variable Collection
```go
type VariableCollection struct {
    ID              string                 // Collection identifier
    Name            string                 // Collection name
    Variables       map[string]*Variable   // Variable collection
    Dependencies    map[string][]string    // Dependency graph
    Metadata        map[string]interface{} // Collection metadata
    CreatedAt       time.Time              // Creation timestamp
    UpdatedAt       time.Time              // Update timestamp
}
```

### Variable Constraints
```go
type VariableConstraints struct {
    MinValue        interface{}            // Minimum value
    MaxValue        interface{}            // Maximum value
    Pattern         string                 // Regex pattern
    AllowedValues   []interface{}          // Allowed values
    CustomValidator string                 // Custom validator
    Metadata        map[string]interface{} // Constraint metadata
}
```

## Variable Categories

### Variable Types
- **String**: String variables
- **Number**: Numeric variables
- **Boolean**: Boolean variables
- **List**: List variables
- **Map**: Map variables
- **Object**: Object variables

### Variable Sources
- **File**: Variables loaded from files
- **Environment**: Environment variables
- **Command Line**: Command line variables
- **Default**: Default variable values
- **Computed**: Computed variable values

### Variable Scopes
- **Global**: Global variables
- **Project**: Project-specific variables
- **Environment**: Environment-specific variables
- **Machine**: Machine-specific variables
- **Action**: Action-specific variables

## Variables Management

### Loading Management
- **File Loading**: Load variables from files
- **Configuration Parsing**: Parse variable configurations
- **Dependency Resolution**: Resolve variable dependencies
- **Validation**: Validate variable configurations
- **Error Handling**: Handle loading errors

### Resolution Management
- **Variable Resolution**: Resolve variable values
- **Interpolation**: Interpolate variable references
- **Dependency Resolution**: Resolve variable dependencies
- **Type Conversion**: Convert variable types
- **Error Handling**: Handle resolution errors

### Validation Management
- **Type Validation**: Validate variable types
- **Constraint Validation**: Validate variable constraints
- **Dependency Validation**: Validate variable dependencies
- **Circular Detection**: Detect circular dependencies
- **Error Reporting**: Report validation errors

## Variables Operations

### Loading Operations
- **Load Variables**: Load variables from files
- **Parse Configuration**: Parse variable configurations
- **Validate Configuration**: Validate variable configurations
- **Resolve Dependencies**: Resolve variable dependencies
- **Error Handling**: Handle loading errors

### Resolution Operations
- **Resolve Variable**: Resolve single variable
- **Resolve All**: Resolve all variables
- **Interpolate**: Interpolate variable references
- **Type Convert**: Convert variable types
- **Error Handling**: Handle resolution errors

### Validation Operations
- **Validate Variable**: Validate single variable
- **Validate All**: Validate all variables
- **Check Dependencies**: Check variable dependencies
- **Detect Circular**: Detect circular dependencies
- **Error Reporting**: Report validation errors

## Security Features

### Data Protection
- **Sensitive Variables**: Mark variables as sensitive
- **Encryption**: Encrypt sensitive variables
- **Access Control**: Control variable access
- **Audit Logging**: Log variable access
- **Data Sanitization**: Sanitize variable data

### Validation Security
- **Input Validation**: Validate variable input
- **Type Safety**: Ensure type safety
- **Constraint Validation**: Validate variable constraints
- **Security Scanning**: Scan for security issues
- **Compliance**: Ensure compliance requirements

### Access Security
- **Access Control**: Control variable access
- **Authentication**: Authenticate variable access
- **Authorization**: Authorize variable operations
- **Audit Logging**: Log variable operations
- **Monitoring**: Monitor variable access

## Performance Features

### Loading Performance
- **Lazy Loading**: Load variables on demand
- **Caching**: Cache variable configurations
- **Parallel Loading**: Load variables in parallel
- **Incremental Loading**: Load variables incrementally
- **Optimization**: Optimize loading performance

### Resolution Performance
- **Cached Resolution**: Cache resolved values
- **Parallel Resolution**: Resolve variables in parallel
- **Incremental Resolution**: Resolve variables incrementally
- **Optimization**: Optimize resolution performance
- **Monitoring**: Monitor resolution performance

### Validation Performance
- **Cached Validation**: Cache validation results
- **Parallel Validation**: Validate variables in parallel
- **Incremental Validation**: Validate variables incrementally
- **Optimization**: Optimize validation performance
- **Monitoring**: Monitor validation performance

## CLI Commands

### Variables Management
```bash
# List variables
spooky variables list <project>

# Show variable value
spooky variables get <project> <variable-name>

# Set variable value
spooky variables set <project> <variable-name> <value>

# Delete variable
spooky variables delete <project> <variable-name>
```

### Variables Validation
```bash
# Validate variables
spooky variables validate <project>

# Check dependencies
spooky variables check-dependencies <project>

# Detect circular dependencies
spooky variables detect-circular <project>

# Show variable dependencies
spooky variables dependencies <project> <variable-name>
```

### Variables Resolution
```bash
# Resolve variables
spooky variables resolve <project>

# Show resolved values
spooky variables show-resolved <project>

# Export variables
spooky variables export <project> --format json

# Import variables
spooky variables import <project> --file variables.json
```

### Variables Operations
```bash
# Create variable
spooky variables create <project> <variable-name> --type string --value "default"

# Update variable
spooky variables update <project> <variable-name> --value "new-value"

# Show variable info
spooky variables info <project> <variable-name>

# Search variables
spooky variables search <project> --pattern "db_*"
```

## Configuration

### Variables Configuration
```hcl
# variables/config.hcl
variables_config {
  # Loading settings
  loading {
    auto_load = true
    watch_files = true
    reload_interval = "30s"
    max_file_size = "1MB"
  }
  
  # Resolution settings
  resolution {
    strict_mode = true
    allow_undefined = false
    max_depth = 10
    timeout = "30s"
  }
  
  # Validation settings
  validation {
    validate_types = true
    validate_constraints = true
    check_dependencies = true
    detect_circular = true
  }
  
  # Security settings
  security {
    encrypt_sensitive = true
    audit_logging = true
    access_control = true
    data_sanitization = true
  }
  
  # Performance settings
  performance {
    cache_enabled = true
    cache_ttl = "1h"
    parallel_loading = true
    parallel_resolution = true
  }
}
```

### Variable Definition
```hcl
# variables/database.hcl
variable "database_host" {
  name = "Database Host"
  description = "Database server hostname"
  
  type = "string"
  default = "localhost"
  required = true
  
  constraints {
    pattern = "^[a-zA-Z0-9.-]+$"
    allowed_values = ["localhost", "127.0.0.1", "db.example.com"]
  }
  
  metadata {
    environment = "production"
    service = "database"
    owner = "db-team"
  }
}

variable "database_port" {
  name = "Database Port"
  description = "Database server port"
  
  type = "number"
  default = 5432
  required = false
  
  constraints {
    min_value = 1024
    max_value = 65535
  }
  
  metadata {
    environment = "production"
    service = "database"
  }
}

variable "database_password" {
  name = "Database Password"
  description = "Database password"
  
  type = "string"
  required = true
  sensitive = true
  encrypted = true
  
  constraints {
    min_length = 8
    pattern = "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d).*$"
  }
  
  metadata {
    environment = "production"
    service = "database"
    owner = "db-team"
  }
}
```

## Examples

### Basic Variables
```hcl
# variables/basic.hcl
variable "app_name" {
  name = "Application Name"
  description = "Name of the application"
  
  type = "string"
  default = "my-app"
  required = false
}

variable "app_version" {
  name = "Application Version"
  description = "Version of the application"
  
  type = "string"
  default = "1.0.0"
  required = false
}

variable "debug_mode" {
  name = "Debug Mode"
  description = "Enable debug mode"
  
  type = "boolean"
  default = false
  required = false
}
```

### Advanced Variables
```hcl
# variables/advanced.hcl
variable "database_config" {
  name = "Database Configuration"
  description = "Database configuration object"
  
  type = "object"
  default = {
    host = "localhost"
    port = 5432
    name = "myapp"
    ssl_mode = "require"
  }
  required = true
  
  constraints {
    custom_validator = "validate_database_config"
  }
  
  metadata {
    environment = "production"
    service = "database"
    owner = "db-team"
  }
}

variable "api_endpoints" {
  name = "API Endpoints"
  description = "List of API endpoints"
  
  type = "list"
  default = [
    "https://api.example.com/v1",
    "https://api.example.com/v2"
  ]
  required = false
  
  constraints {
    min_length = 1
    max_length = 10
  }
  
  metadata {
    environment = "production"
    service = "api"
    owner = "api-team"
  }
}

variable "feature_flags" {
  name = "Feature Flags"
  description = "Feature flags configuration"
  
  type = "map"
  default = {
    new_ui = true
    beta_features = false
    analytics = true
  }
  required = false
  
  metadata {
    environment = "production"
    service = "application"
    owner = "dev-team"
  }
}
```

## Integration Examples

### Templates Integration
```go
// Use variables in templates
variablesIntegration := manager.GetVariablesIntegration()
templatesIntegration := manager.GetTemplatesIntegration()

// Load variables
variables, err := variablesIntegration.LoadVariables("my-project")
if err != nil {
    return err
}

// Resolve variables
resolvedVars, err := variablesIntegration.ResolveVariables(variables)
if err != nil {
    return err
}

// Use variables in template
template, err := templatesIntegration.LoadTemplate("templates/config.tmpl")
if err != nil {
    return err
}

result, err := templatesIntegration.RenderTemplate(template, resolvedVars)
if err != nil {
    return err
}
```

### Actions Integration
```go
// Use variables in actions
variablesIntegration := manager.GetVariablesIntegration()
actionsIntegration := manager.GetActionsIntegration()

// Load and resolve variables
variables, err := variablesIntegration.LoadVariables("my-project")
if err != nil {
    return err
}

resolvedVars, err := variablesIntegration.ResolveVariables(variables)
if err != nil {
    return err
}

// Use variables in action
action := &spookytypes.Action{
    Name: "deploy-application",
    Script: "deploy.sh",
    Environment: map[string]string{
        "APP_NAME": resolvedVars["app_name"].(string),
        "APP_VERSION": resolvedVars["app_version"].(string),
        "DEBUG_MODE": fmt.Sprintf("%t", resolvedVars["debug_mode"].(bool)),
    },
}

result, err := actionsIntegration.RunAction(action)
if err != nil {
    return err
}
```

### Configuration Integration
```go
// Use variables in configuration
variablesIntegration := manager.GetVariablesIntegration()
configIntegration := manager.GetConfigIntegration()

// Load and resolve variables
variables, err := variablesIntegration.LoadVariables("my-project")
if err != nil {
    return err
}

resolvedVars, err := variablesIntegration.ResolveVariables(variables)
if err != nil {
    return err
}

// Use variables in configuration
config := &spookytypes.Config{
    Database: &spookytypes.DatabaseConfig{
        Host: resolvedVars["database_host"].(string),
        Port: resolvedVars["database_port"].(int),
        Name: resolvedVars["database_name"].(string),
    },
    Application: &spookytypes.ApplicationConfig{
        Name: resolvedVars["app_name"].(string),
        Version: resolvedVars["app_version"].(string),
        Debug: resolvedVars["debug_mode"].(bool),
    },
}

err = configIntegration.UpdateConfig(config)
if err != nil {
    return err
}
```

## Best Practices

### Variable Design
- Use descriptive variable names
- Provide clear descriptions
- Set appropriate default values
- Use proper variable types
- Implement proper constraints

### Security
- Mark sensitive variables appropriately
- Encrypt sensitive data
- Implement proper access controls
- Validate variable input
- Monitor variable access

### Performance
- Use lazy loading for large variable sets
- Implement caching for frequently accessed variables
- Optimize variable resolution
- Monitor variable performance
- Use appropriate timeouts

### Management
- Organize variables logically
- Use consistent naming conventions
- Document variable dependencies
- Regular variable validation
- Monitor variable usage

## Troubleshooting

### Common Issues

#### Loading Issues
```bash
# Check variable files
spooky variables list <project>

# Validate variable configuration
spooky variables validate <project>

# Check file syntax
spooky variables validate <project> --syntax-only

# Show loading errors
spooky variables list <project> --verbose
```

#### Resolution Issues
```bash
# Check variable resolution
spooky variables resolve <project>

# Show unresolved variables
spooky variables resolve <project> --show-unresolved

# Check dependencies
spooky variables check-dependencies <project>

# Detect circular dependencies
spooky variables detect-circular <project>
```

#### Validation Issues
```bash
# Validate variables
spooky variables validate <project>

# Check specific variable
spooky variables validate <project> --variable <name>

# Show validation errors
spooky variables validate <project> --verbose

# Fix validation issues
spooky variables validate <project> --fix
```

#### Performance Issues
```bash
# Check variable performance
spooky variables performance <project>

# Monitor variable loading
spooky variables list <project> --timing

# Check cache status
spooky variables cache-status <project>

# Optimize variables
spooky variables optimize <project>
```

## API Reference

### VariablesIntegration Interface
```go
type VariablesIntegration interface {
    LoadVariables(ctx context.Context, projectPath string) (*spookytypes.VariableCollection, error)
    ResolveVariables(ctx context.Context, variables *spookytypes.VariableCollection) (map[string]interface{}, error)
    ValidateVariables(ctx context.Context, variables *spookytypes.VariableCollection) (*spookytypes.ValidationResult, error)
    GetVariable(ctx context.Context, name string) (*spookytypes.Variable, error)
}
```

### Variables Manager Methods
```go
// Variable management
LoadVariables(ctx context.Context, projectPath string) (*spookytypes.VariableCollection, error)
ResolveVariables(ctx context.Context, variables *spookytypes.VariableCollection) (map[string]interface{}, error)
ValidateVariables(ctx context.Context, variables *spookytypes.VariableCollection) (*spookytypes.ValidationResult, error)

// Variable operations
GetVariable(ctx context.Context, name string) (*spookytypes.Variable, error)
SetVariable(ctx context.Context, variable *spookytypes.Variable) error
DeleteVariable(ctx context.Context, name string) error
ListVariables(ctx context.Context) ([]*spookytypes.Variable, error)

// Dependency management
CheckDependencies(ctx context.Context, variables *spookytypes.VariableCollection) (*spookytypes.ValidationResult, error)
DetectCircularDependencies(ctx context.Context, variables *spookytypes.VariableCollection) ([]string, error)
ResolveDependencies(ctx context.Context, variables *spookytypes.VariableCollection) ([]string, error)
```

## Related Documentation

- [Variables API Reference](VARIABLES_API_REFERENCE.md) - Complete API documentation
- [Variables User Guide](VARIABLES_USER_GUIDE.md) - User guide and examples
- [Variables Troubleshooting](VARIABLES_TROUBLESHOOTING.md) - Troubleshooting guide
- [Templates System](TEMPLATES_SYSTEM.md) - Templates integration
- [Actions System](ACTIONS_SYSTEM.md) - Actions integration
- [Configuration System](CONFIGURATION_SYSTEM.md) - Configuration integration
- [Secrets System](SECRETS_SYSTEM.md) - Secrets integration
