# Secrets System API Reference

## Overview

This document provides a comprehensive API reference for the spooky secrets system. It covers all interfaces, types, methods, and implementation details for developers working with the secrets system.

**Status: Partially Implemented** - The secrets system has basic functionality but SSH-based secret collection has known issues that need to be addressed.

## Core Interfaces

### SecretsIntegration Interface

The `SecretsIntegration` interface provides the primary entry point for secrets operations:

```go
type SecretsIntegration interface {
    // LoadSecrets loads secrets from the specified project path
    LoadSecrets(ctx context.Context, projectPath string) (map[string]interface{}, error)
    
    // ValidateSecrets validates secret definitions
    ValidateSecrets(ctx context.Context, secrets map[string]interface{}) (*ValidationResult, error)
    
    // EncryptSecrets encrypts secrets using age encryption
    EncryptSecrets(ctx context.Context, secrets map[string]interface{}, publicKey string) (map[string]string, error)
    
    // DecryptSecrets decrypts secrets using age decryption
    DecryptSecrets(ctx context.Context, encryptedSecrets map[string]string, privateKey string) (map[string]interface{}, error)
}
```

**Implementation Status**: ⚠️ **Partially Implemented** - Basic functionality exists but SSH-based collection has issues

### SecretsManager Interface

The `SecretsManager` interface provides secrets management and encryption:

```go
type SecretsManager interface {
    // LoadSecrets loads secrets from project configuration
    LoadSecrets(ctx context.Context, projectPath string) (map[string]*spookytypessecrets.Secret, error)
    
    // ValidateSecrets validates secret definitions
    ValidateSecrets(ctx context.Context, secrets map[string]*spookytypessecrets.Secret) (*ValidationResult, error)
    
    // EncryptSecrets encrypts secrets using age encryption
    EncryptSecrets(ctx context.Context, secrets map[string]*spookytypessecrets.Secret, publicKey string) (map[string]string, error)
    
    // DecryptSecrets decrypts secrets using age decryption
    DecryptSecrets(ctx context.Context, encryptedSecrets map[string]string, privateKey string) (map[string]*spookytypessecrets.Secret, error)
}
```

**Implementation Status**: ⚠️ **Partially Implemented** - Basic loading and validation exist but encryption has issues

## Current Implementation Status

### ✅ Working Components

1. **Secret Loading**: Loading secrets from HCL configuration files
2. **Secret Validation**: Basic validation of secret definitions
3. **Secret Structure**: Proper secret type definitions and structures
4. **CLI Integration**: `spooky secrets list` command with filtering options
5. **Project Integration**: Secrets loading from project configuration
6. **Basic Validation**: Secret definition validation and error handling
7. **Filtering Support**: Support for secret name and type filtering
8. **Export Support**: Secret export to JSON and HCL formats

### ⚠️ Known Issues

1. **SSH-Based Collection**: SSH-based secret collection has implementation issues
2. **Age Encryption**: Age encryption/decryption has implementation problems
3. **Remote Secret Reading**: Cannot properly read secrets from remote machines
4. **Parallel Processing**: No parallel secret collection support
5. **Import Functionality**: No secret import capabilities
6. **Template Integration**: No template secret integration

### 🔄 In Progress

1. **SSH Collection Fixes**: Addressing SSH-based secret collection issues
2. **Age Encryption Fixes**: Implementing proper age encryption/decryption
3. **Collection Enhancements**: Improving secret collection reliability

## Implementation Details

### Secret Loading System

The secrets system loads secrets from HCL configuration files:

```go
type SecretLoader struct {
    logger spookylogging.Logger
}

func (l *SecretLoader) LoadSecrets(ctx context.Context, projectPath string) (map[string]*spookytypessecrets.Secret, error) {
    secrets := make(map[string]*spookytypessecrets.Secret)
    
    // Load secrets.hcl file
    secretsPath := filepath.Join(projectPath, "secrets.hcl")
    if data, err := os.ReadFile(secretsPath); err == nil {
        if err := l.parseSecretsFile(data, secrets); err != nil {
            return nil, fmt.Errorf("failed to parse secrets.hcl: %w", err)
        }
    }
    
    // Load secrets from secrets/ directory
    secretsDir := filepath.Join(projectPath, "secrets")
    if entries, err := os.ReadDir(secretsDir); err == nil {
        for _, entry := range entries {
            if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".hcl") {
                filePath := filepath.Join(secretsDir, entry.Name())
                if data, err := os.ReadFile(filePath); err == nil {
                    if err := l.parseSecretsFile(data, secrets); err != nil {
                        return nil, fmt.Errorf("failed to parse %s: %w", entry.Name(), err)
                    }
                }
            }
        }
    }
    
    return secrets, nil
}

func (l *SecretLoader) parseSecretsFile(data []byte, secrets map[string]*spookytypessecrets.Secret) error {
    var config struct {
        Secrets []*spookytypessecrets.Secret `hcl:"secret,block"`
    }
    
    if err := hcl.Unmarshal(data, &config); err != nil {
        return fmt.Errorf("failed to parse HCL: %w", err)
    }
    
    for _, secret := range config.Secrets {
        if secret.Name == "" {
            return fmt.Errorf("secret name is required")
        }
        
        if _, exists := secrets[secret.Name]; exists {
            return fmt.Errorf("duplicate secret name: %s", secret.Name)
        }
        
        secrets[secret.Name] = secret
    }
    
    return nil
}
```

**Supported Secret Sources:**
- **Local Secrets**: Secrets defined in `secrets.hcl` and `secrets/*.hcl` files
- **Environment Secrets**: Secrets from environment variables
- **SSH Secrets**: Secrets collected from remote machines (has issues)
- **Computed Secrets**: Secrets computed from other secrets

### Secret Validation System

Secrets are validated against schemas and business rules:

```go
type SecretValidator struct {
    logger spookylogging.Logger
}

func (v *SecretValidator) ValidateSecrets(ctx context.Context, secrets map[string]*spookytypessecrets.Secret) (*spookytypes.ValidationResult, error) {
    var errors []spookyschemas.SchemaError
    var warnings []spookyschemas.SchemaError
    
    for name, secret := range secrets {
        // Validate secret name
        if name == "" {
            errors = append(errors, spookyschemas.SchemaError{
                Message: "secret name cannot be empty",
            })
            continue
        }
        
        // Validate secret structure
        if err := v.validateSecretStructure(secret); err != nil {
            errors = append(errors, spookyschemas.SchemaError{
                Message: fmt.Sprintf("secret %s: %s", name, err.Error()),
            })
        }
        
        // Validate secret type
        if err := v.validateSecretType(secret); err != nil {
            errors = append(errors, spookyschemas.SchemaError{
                Message: fmt.Sprintf("secret %s: %s", name, err.Error()),
            })
        }
        
        // Validate secret constraints
        if err := v.validateSecretConstraints(secret); err != nil {
            errors = append(errors, spookyschemas.SchemaError{
                Message: fmt.Sprintf("secret %s: %s", name, err.Error()),
            })
        }
    }
    
    return &spookytypes.ValidationResult{
        Valid:    len(errors) == 0,
        Errors:   errors,
        Warnings: warnings,
    }, nil
}
```

### Secret Encryption System

Secrets are encrypted using age encryption (currently has issues):

```go
type SecretEncryptor struct {
    logger spookylogging.Logger
}

func (e *SecretEncryptor) EncryptSecrets(ctx context.Context, secrets map[string]*spookytypessecrets.Secret, publicKey string) (map[string]string, error) {
    encrypted := make(map[string]string)
    
    for name, secret := range secrets {
        // Convert secret to string
        secretData, err := e.secretToString(secret)
        if err != nil {
            return nil, fmt.Errorf("failed to convert secret %s to string: %w", name, err)
        }
        
        // Encrypt secret using age
        encryptedData, err := e.encryptWithAge(secretData, publicKey)
        if err != nil {
            return nil, fmt.Errorf("failed to encrypt secret %s: %w", name, err)
        }
        
        encrypted[name] = encryptedData
    }
    
    return encrypted, nil
}

func (e *SecretEncryptor) DecryptSecrets(ctx context.Context, encryptedSecrets map[string]string, privateKey string) (map[string]*spookytypessecrets.Secret, error) {
    decrypted := make(map[string]*spookytypessecrets.Secret)
    
    for name, encryptedData := range encryptedSecrets {
        // Decrypt secret using age
        decryptedData, err := e.decryptWithAge(encryptedData, privateKey)
        if err != nil {
            return nil, fmt.Errorf("failed to decrypt secret %s: %w", name, err)
        }
        
        // Convert string back to secret
        secret, err := e.stringToSecret(decryptedData)
        if err != nil {
            return nil, fmt.Errorf("failed to convert decrypted data for secret %s: %w", name, err)
        }
        
        decrypted[name] = secret
    }
    
    return decrypted, nil
}
```

## Type Definitions

### Secret Types

```go
// Secret represents a secret definition
type Secret struct {
    // Secret name (required)
    Name string `json:"name" hcl:"name"`
    
    // Secret description (optional)
    Description string `json:"description,omitempty" hcl:"description,optional"`
    
    // Secret type (password, key, certificate, token)
    Type string `json:"type" hcl:"type"`
    
    // Secret value (for local secrets)
    Value interface{} `json:"value,omitempty" hcl:"value,optional"`
    
    // Secret source (for remote secrets)
    Source *SecretSource `json:"source,omitempty" hcl:"source,optional"`
    
    // Secret constraints
    Constraints *SecretConstraints `json:"constraints,omitempty" hcl:"constraints,optional"`
    
    // Secret metadata
    Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// SecretSource defines how to obtain secret values
type SecretSource struct {
    // Source type (environment, ssh, computed)
    Type string `json:"type" hcl:"type"`
    
    // Source configuration
    Config map[string]interface{} `json:"config,omitempty" hcl:"config,optional"`
    
    // SSH source configuration
    SSH *SSHSecretSource `json:"ssh,omitempty" hcl:"ssh,optional"`
    
    // Environment source configuration
    Environment *EnvironmentSecretSource `json:"environment,omitempty" hcl:"environment,optional"`
}

// SSHSecretSource defines SSH-based secret collection
type SSHSecretSource struct {
    // Target machine
    Machine string `json:"machine" hcl:"machine"`
    
    // SSH command to execute
    Command string `json:"command" hcl:"command"`
    
    // File path to read (alternative to command)
    File string `json:"file,omitempty" hcl:"file,optional"`
    
    // Parse configuration
    Parse *ParseConfig `json:"parse,omitempty" hcl:"parse,optional"`
}

// EnvironmentSecretSource defines environment secret collection
type EnvironmentSecretSource struct {
    // Environment variable name
    Name string `json:"name" hcl:"name"`
    
    // Default value if not set
    Default interface{} `json:"default,omitempty" hcl:"default,optional"`
}

// SecretConstraints defines secret validation constraints
type SecretConstraints struct {
    // Required constraint
    Required bool `json:"required,omitempty" hcl:"required,optional"`
    
    // String constraints
    MinLength int    `json:"min_length,omitempty" hcl:"min_length,optional"`
    MaxLength int    `json:"max_length,omitempty" hcl:"max_length,optional"`
    Pattern   string `json:"pattern,omitempty" hcl:"pattern,optional"`
    
    // Allowed values
    AllowedValues []interface{} `json:"allowed_values,omitempty" hcl:"allowed_values,optional"`
}

// ParseConfig defines how to parse secret values
type ParseConfig struct {
    // Parse format (json, yaml, hcl, text)
    Format string `json:"format" hcl:"format"`
    
    // Parse path (for nested values)
    Path string `json:"path,omitempty" hcl:"path,optional"`
    
    // Parse options
    Options map[string]interface{} `json:"options,omitempty" hcl:"options,optional"`
}
```

### Secret Context Types

```go
// SecretContext provides context for secret operations
type SecretContext struct {
    // Project path
    ProjectPath string `json:"project_path" hcl:"project_path"`
    
    // Secret being processed
    Secret *Secret `json:"secret" hcl:"secret"`
    
    // Operation timestamp
    Timestamp time.Time `json:"timestamp" hcl:"timestamp"`
    
    // Operation metadata
    Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// SecretResult represents the result of secret operations
type SecretResult struct {
    // Secret context
    Context *SecretContext `json:"context" hcl:"context"`
    
    // Operation success
    Success bool `json:"success" hcl:"success"`
    
    // Secret value
    Value interface{} `json:"value,omitempty" hcl:"value,optional"`
    
    // Operation error
    Error string `json:"error,omitempty" hcl:"error,optional"`
    
    // Operation duration
    Duration time.Duration `json:"duration" hcl:"duration"`
}
```

## Error Handling

### Secret Errors

```go
// SecretError represents secret operation errors
type SecretError struct {
    SecretName string `json:"secret_name" hcl:"secret_name"`
    Error      string `json:"error" hcl:"error"`
    Details    string `json:"details,omitempty" hcl:"details,optional"`
}

// SecretValidationError represents secret validation errors
type SecretValidationError struct {
    Field   string `json:"field" hcl:"field"`
    Message string `json:"message" hcl:"message"`
    Value   string `json:"value,omitempty" hcl:"value,optional"`
}
```

### Validation Implementation

```go
// ValidateSecret validates a single secret
func (v *SecretValidator) ValidateSecret(secret *spookytypessecrets.Secret) error {
    if secret == nil {
        return fmt.Errorf("secret cannot be nil")
    }
    
    // Validate required fields
    if secret.Name == "" {
        return fmt.Errorf("secret name is required")
    }
    
    if secret.Type == "" {
        return fmt.Errorf("secret type is required")
    }
    
    // Validate secret type
    validTypes := []string{"password", "key", "certificate", "token"}
    valid := false
    for _, t := range validTypes {
        if secret.Type == t {
            valid = true
            break
        }
    }
    if !valid {
        return fmt.Errorf("invalid secret type: %s (valid types: %v)", secret.Type, validTypes)
    }
    
    // Validate secret source
    if secret.Source != nil {
        if err := v.validateSecretSource(secret.Source); err != nil {
            return fmt.Errorf("invalid secret source: %w", err)
        }
    }
    
    // Validate secret constraints
    if secret.Constraints != nil {
        if err := v.validateSecretConstraints(secret.Constraints); err != nil {
            return fmt.Errorf("invalid secret constraints: %w", err)
        }
    }
    
    return nil
}
```

## CLI Commands

### Secrets List Command

```bash
# List all secrets in a project
spooky secrets list ./my-project

# List secrets with specific types
spooky secrets list ./my-project --type password,key

# List secrets with specific names
spooky secrets list ./my-project --names db_password,api_key

# List secrets with verbose output
spooky secrets list ./my-project --verbose
```

### Secrets Export Command

```bash
# Export secrets to JSON format
spooky secrets export ./my-project --format json --output secrets.json

# Export secrets to HCL format
spooky secrets export ./my-project --format hcl --output secrets.hcl

# Export secrets with specific types
spooky secrets export ./my-project --type password,key --format json --output password-secrets.json
```

### Secrets Validation Command

```bash
# Validate secrets in a project
spooky secrets validate ./my-project

# Validate secrets with verbose output
spooky secrets validate ./my-project --verbose
```

## Integration Examples

### Basic Secret Definition

```hcl
# secrets.hcl
secrets {
  secret "db_password" {
    description = "Database password"
    type = "password"
    value = "secret_password_123"
    
    constraints {
      required = true
      min_length = 8
    }
  }
  
  secret "api_key" {
    description = "API key for external service"
    type = "key"
    
    source {
      type = "ssh"
      ssh {
        machine = "app-server"
        command = "cat /etc/app/api.key"
      }
    }
    
    constraints {
      required = true
      min_length = 32
    }
  }
  
  secret "ssl_certificate" {
    description = "SSL certificate"
    type = "certificate"
    
    source {
      type = "ssh"
      ssh {
        machine = "web-server"
        file = "/etc/ssl/certs/website.crt"
      }
    }
    
    constraints {
      required = true
    }
  }
  
  secret "auth_token" {
    description = "Authentication token"
    type = "token"
    
    source {
      type = "environment"
      environment {
        name = "AUTH_TOKEN"
      }
    }
    
    constraints {
      required = true
      min_length = 16
    }
  }
}
```

### Secret Loading and Validation

```go
// Secret loading and validation example
func loadAndValidateSecrets(projectPath string) error {
    ctx := context.Background()
    
    // Create secret manager
    manager := spookysecrets.NewManager(loader, validator, logger)
    
    // Load secrets
    secrets, err := manager.LoadSecrets(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load secrets: %w", err)
    }
    
    // Validate secrets
    result, err := manager.ValidateSecrets(ctx, secrets)
    if err != nil {
        return fmt.Errorf("failed to validate secrets: %w", err)
    }
    
    if !result.Valid {
        fmt.Println("Secret validation failed:")
        for _, error := range result.Errors {
            fmt.Printf("  - %s\n", error.Message)
        }
        return fmt.Errorf("secret validation failed")
    }
    
    fmt.Printf("Loaded and validated %d secrets\n", len(secrets))
    return nil
}
```

### Secret Encryption

```go
// Secret encryption example
func encryptSecrets(projectPath string, publicKey string) error {
    ctx := context.Background()
    
    // Create secret manager
    manager := spookysecrets.NewManager(loader, validator, logger)
    
    // Load secrets
    secrets, err := manager.LoadSecrets(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load secrets: %w", err)
    }
    
    // Encrypt secrets
    encrypted, err := manager.EncryptSecrets(ctx, secrets, publicKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt secrets: %w", err)
    }
    
    // Print encrypted secrets
    for name, encryptedData := range encrypted {
        fmt.Printf("%s: %s\n", name, encryptedData)
    }
    
    return nil
}
```

## Current Limitations

### Collection Limitations

1. **SSH Integration Issues**: SSH-based secret collection has known problems
2. **No Parallel Collection**: Secrets are collected sequentially, not in parallel
3. **No Result Caching**: Secret values are not cached between operations
4. **No Incremental Collection**: Always collects all secrets
5. **Limited Source Types**: Only basic source types are supported

### Encryption Limitations

1. **Age Encryption Issues**: Age encryption/decryption has implementation problems
2. **No Key Management**: No proper key management system
3. **No Key Rotation**: No key rotation capabilities
4. **No Audit Logging**: No audit logging for secret operations

### Integration Limitations

1. **No Template Integration**: Secrets are not integrated with template system
2. **No Action Integration**: Secrets are not used in action system
3. **No Facts Integration**: Secrets are not integrated with facts system
4. **No Conditional Access**: No conditional secret access

## Future Enhancements

### Planned Features

1. **SSH Collection Fixes**: Resolve SSH-based secret collection issues
2. **Age Encryption Fixes**: Implement proper age encryption/decryption
3. **Parallel Collection**: Implement parallel secret collection
4. **Result Caching**: Add secret value caching
5. **Incremental Collection**: Support incremental secret collection
6. **Advanced Sources**: Support more secret source types

### Integration Enhancements

1. **Template Integration**: Integrate secrets with template system
2. **Action Integration**: Use secrets in action system
3. **Facts Integration**: Integrate secrets with facts system
4. **Advanced Encryption**: Support advanced encryption features

## Summary

The secrets system provides basic secret loading and validation capabilities but has significant limitations with SSH-based collection and age encryption that need to be addressed. The system is functional for basic use cases but requires improvements for production use.

**Status**: ⚠️ **Partially Implemented** - Basic functionality exists but SSH-based collection and age encryption have issues that need to be resolved.
