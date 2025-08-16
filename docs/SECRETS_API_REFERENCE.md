# Secrets System API Reference

## Overview

This document provides a comprehensive API reference for the spooky secrets system. It covers all interfaces, types, methods, and implementation details for developers working with the secrets system.

**Status: Fully Implemented** - The secrets system has complete age encryption functionality with HCL integration support, recipient management, and comprehensive validation.

## Core Interfaces

### SecretsIntegration Interface

The `SecretsIntegration` interface provides the primary entry point for secrets operations:

```go
type SecretsIntegration interface {
    // Age-specific methods
    EncryptWithAge(ctx context.Context, data []byte, recipients []string) ([]byte, error)
    DecryptWithAge(ctx context.Context, data []byte, identityPath string) ([]byte, error)
    EncryptWithPassphrase(ctx context.Context, data []byte, passphrase string) ([]byte, error)
    DecryptWithPassphrase(ctx context.Context, data []byte, passphrase string) ([]byte, error)

    // Key management
    ValidateAgeKey(ctx context.Context, keyPath string) error
    ListRecipients(ctx context.Context, encryptedData []byte) ([]string, error)

    // Application-level validation
    ValidateAgeEncryptedValue(ctx context.Context, value string) error
    LoadRecipients(ctx context.Context, recipientsPath string) ([]string, error)
    LoadIdentities(ctx context.Context, identitiesPath string) ([]string, error)

    // HCL encryption and decryption methods
    EncryptHCLValues(ctx context.Context, data interface{}, recipients []string, dryRun bool) error
    EncryptHCLValuesSensitive(ctx context.Context, data interface{}, recipients []string, dryRun bool, shouldEncrypt func(path []string, value interface{}) bool) error
    EncryptHCLValuesWithJSONSupport(ctx context.Context, data interface{}, recipients []string, dryRun bool) error
    DecryptHCLValues(ctx context.Context, data interface{}, identityPath string) error
    DecryptHCLValuesWithJSONSupport(ctx context.Context, data interface{}, identityPath string) error
}
```

**Implementation Status**: ✅ **Fully Implemented** - Complete age encryption/decryption functionality with HCL integration support

## Current Implementation Status

### ✅ Working Components

1. **Age Encryption**: Full age encryption/decryption functionality with X25519 and scrypt recipients
2. **Age Key Management**: Complete age key validation and management
3. **HCL Integration**: Full HCL value encryption/decryption with JSON support
4. **CLI Integration**: `spooky secrets validate` command
5. **Project Integration**: Secrets integration with variables, machines, and facts systems
6. **Validation**: Comprehensive age configuration and encrypted value validation
7. **Recipient Management**: Loading and managing age recipients from files
8. **Identity Management**: Loading and managing age identities from directories
9. **Passphrase Support**: Passphrase-based encryption/decryption using scrypt
10. **Armored Output**: Support for armored age output with configuration
11. **Integration Support**: Full integration with variables, machines, and facts systems
12. **Legacy Method Support**: Legacy methods return appropriate error messages directing users to age methods

### ✅ Fully Implemented Features

#### Age Encryption Methods
- `EncryptWithAge()` - Encrypt data with X25519 recipients
- `DecryptWithAge()` - Decrypt data with identity files
- `EncryptWithPassphrase()` - Encrypt data with scrypt passphrase
- `DecryptWithPassphrase()` - Decrypt data with scrypt passphrase

#### Key Management
- `ValidateAgeKey()` - Validate age identity files
- `ListRecipients()` - Extract recipient information from encrypted data
- `LoadRecipients()` - Load recipients from files
- `LoadIdentities()` - Load identities from directories

#### HCL Integration
- `EncryptHCLValues()` - Encrypt values in HCL structures
- `DecryptHCLValues()` - Decrypt values in HCL structures
- `EncryptHCLValuesSensitive()` - Encrypt with sensitivity awareness
- `EncryptHCLValuesWithJSONSupport()` - Encrypt with JSON serialization
- `DecryptHCLValuesWithJSONSupport()` - Decrypt with JSON deserialization

#### Validation
- `ValidateAgeEncryptedValue()` - Validate age-encrypted values
- Comprehensive error handling and logging
- Support for both armored and binary age formats

### ❌ Intentionally Not Implemented

The following components are intentionally not implemented as they are not needed for the secrets system design:

1. **Secret Loading**: No dedicated secret loading from HCL configuration files (secrets are handled through existing systems)
2. **Secret Storage**: No dedicated secret storage (secrets are part of machines inventories, custom facts, and variables)
3. **Secret Export/Import**: No export/import functionality (not needed for the system design)
4. **Secret Caching**: No secret caching (not needed for the system design)
5. **Secret Templates**: No secret template integration (not needed for the system design)

## Implementation Details

### Age Integration Implementation

The secrets system uses the `filippo.io/age` library for all encryption operations:

```go
// Integration implements the SecretsIntegration interface with age encryption
type Integration struct {
    logger       spookytypeslogging.Logger
    config       *spookytypes.AgeConfig
    hclProcessor *HCLProcessor
}

// NewIntegration creates a new age-focused secrets integration
func NewIntegration(logger spookytypeslogging.Logger, config *spookytypes.AgeConfig) spookyinterfaces.SecretsIntegration {
    return &Integration{
        logger:       logger,
        config:       config,
        hclProcessor: NewHCLProcessor(logger),
    }
}
```

### HCL Processor Implementation

The HCL processor handles encryption and decryption of values in HCL-compatible structures:

```go
// HCLProcessor handles encryption and decryption of HCL values
type HCLProcessor struct {
    logger spookytypeslogging.Logger
}

// NewHCLProcessor creates a new HCL processor
func NewHCLProcessor(logger spookytypeslogging.Logger) *HCLProcessor {
    return &HCLProcessor{
        logger: logger,
    }
}
```

### Configuration Support

The secrets system supports age configuration through the `AgeConfig` type:

```go
type AgeConfig struct {
    Encryption *AgeEncryptionConfig `hcl:"encryption,block"`
}

type AgeEncryptionConfig struct {
    Armor bool `hcl:"armor,optional"`
}
```

## Usage Examples

### Basic Age Encryption

```go
// Basic age encryption example
func encryptDataWithAge(data []byte, recipients []string) error {
    ctx := context.Background()
    
    // Create secrets integration
    secretsIntegration := spookysecrets.NewIntegration(logger, config)
    
    // Encrypt data
    encrypted, err := secretsIntegration.EncryptWithAge(ctx, data, recipients)
    if err != nil {
        return fmt.Errorf("failed to encrypt data: %w", err)
    }
    
    fmt.Printf("Encrypted data: %s\n", string(encrypted))
    return nil
}
```

### HCL Value Encryption

```go
// HCL value encryption example
func encryptHCLValues(data interface{}, recipients []string) error {
    ctx := context.Background()
    
    // Create secrets integration
    secretsIntegration := spookysecrets.NewIntegration(logger, config)
    
    // Encrypt HCL values
    err := secretsIntegration.EncryptHCLValues(ctx, data, recipients, false)
    if err != nil {
        return fmt.Errorf("failed to encrypt HCL values: %w", err)
    }
    
    return nil
}
```

### Age Key Validation

```go
// Age key validation example
func validateAgeKey(keyPath string) error {
    ctx := context.Background()
    
    // Create secrets integration
    secretsIntegration := spookysecrets.NewIntegration(logger, config)
    
    // Validate age key
    err := secretsIntegration.ValidateAgeKey(ctx, keyPath)
    if err != nil {
        return fmt.Errorf("invalid age key: %w", err)
    }
    
    return nil
}
```

## CLI Integration

### Secrets Validation Command

```bash
# Validate age configuration
spooky secrets validate

# Validate age keys
spooky secrets validate --key ~/.config/spooky/identities/age.key

# Validate encrypted values
spooky secrets validate --value "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
```

## Integration with Other Systems

### Variables System Integration

The secrets system integrates with the variables system to encrypt variable values:

```go
// Encrypt variables example
func encryptVariables(projectPath string, recipients []string) error {
    ctx := context.Background()
    
    // Get variables integration
    variablesIntegration := manager.GetVariablesIntegration()
    
    // Get secrets integration
    secretsIntegration := manager.GetSecretsIntegration()
    
    // Load variables
    variables, err := variablesIntegration.LoadVariables(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to load variables: %w", err)
    }
    
    // Encrypt variables
    err = variablesIntegration.EncryptVariables(ctx, projectPath, secretsIntegration, recipients, false)
    if err != nil {
        return fmt.Errorf("failed to encrypt variables: %w", err)
    }
    
    return nil
}
```

### Facts System Integration

The secrets system integrates with the facts system to decrypt encrypted fact values:

```go
// Decrypt facts example
func decryptFacts(facts interface{}, identityPath string) error {
    ctx := context.Background()
    
    // Get facts integration
    factsIntegration := manager.GetFactsIntegration()
    
    // Get secrets integration
    secretsIntegration := manager.GetSecretsIntegration()
    
    // Decrypt facts
    err := factsIntegration.DecryptFacts(ctx, facts, secretsIntegration, identityPath)
    if err != nil {
        return fmt.Errorf("failed to decrypt facts: %w", err)
    }
    
    return nil
}
```

### Machines System Integration

The secrets system integrates with the machines system to decrypt encrypted machine secrets:

```go
// Decrypt machines example
func decryptMachines(machines []spookytypes.Machine, identityPath string) error {
    ctx := context.Background()
    
    // Get machines integration
    machinesIntegration := manager.GetMachinesIntegration()
    
    // Get secrets integration
    secretsIntegration := manager.GetSecretsIntegration()
    
    // Decrypt machines
    err := machinesIntegration.DecryptMachines(ctx, machines, secretsIntegration, identityPath)
    if err != nil {
        return fmt.Errorf("failed to decrypt machines: %w", err)
    }
    
    return nil
}
```

## Error Handling

The secrets system provides comprehensive error handling:

```go
// Error handling example
func handleSecretsError(err error) {
    if err == nil {
        return
    }
    
    // Check for specific error types
    switch {
    case strings.Contains(err.Error(), "failed to parse identity file"):
        fmt.Println("Invalid age identity file format")
    case strings.Contains(err.Error(), "failed to create age encryptor"):
        fmt.Println("Invalid recipient format")
    case strings.Contains(err.Error(), "failed to decrypt data"):
        fmt.Println("Decryption failed - check identity file")
    default:
        fmt.Printf("Secrets error: %v\n", err)
    }
}
```

## Performance Considerations

### Memory Usage

The secrets system is designed for efficient memory usage:

- Encryption operations use streaming for large data
- HCL processing handles large structures efficiently
- No unnecessary data copying during operations

### Parallel Processing

The secrets system supports parallel operations:

- Multiple encryption operations can run concurrently
- HCL processing can handle multiple structures in parallel
- Integration with other systems supports parallel processing

## Security Best Practices

### Key Management

1. **Secure Key Storage**: Store age identity files with appropriate permissions (600)
2. **Key Rotation**: Regularly rotate age keys for enhanced security
3. **Recipient Management**: Carefully manage recipient lists to prevent unauthorized access
4. **Passphrase Security**: Use strong passphrases for scrypt-based encryption

### Configuration Security

1. **Configuration Validation**: Always validate age configuration before use
2. **Error Handling**: Handle encryption/decryption errors appropriately
3. **Logging**: Use secure logging practices to avoid exposing sensitive data

## Troubleshooting

### Common Issues

1. **Invalid Identity File**: Ensure age identity files are properly formatted
2. **Recipient Format**: Verify recipient strings are in correct age format
3. **Permission Errors**: Check file permissions for identity and recipient files
4. **HCL Structure**: Ensure HCL structures are properly formatted for encryption

### Debug Information

The secrets system provides comprehensive logging for debugging:

```go
// Enable debug logging
logger.SetLevel(spookytypes.LogLevelDebug)

// Check encryption configuration
fmt.Printf("Age config: %+v\n", config)

// Validate encrypted values
err := secretsIntegration.ValidateAgeEncryptedValue(ctx, encryptedValue)
if err != nil {
    fmt.Printf("Validation error: %v\n", err)
}
```
