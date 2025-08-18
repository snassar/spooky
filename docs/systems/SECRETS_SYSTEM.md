# Secrets System

## Overview

The Secrets System provides comprehensive secrets management and encryption capabilities for the spooky codebase. It enables secure storage, encryption, decryption, and management of sensitive data using the age encryption format and external key management.

**Status**: **Implemented** - Complete secrets system with encryption, decryption, key management, and CLI integration.

## Related Systems

This system integrates with and secures several other spooky systems:

- **[Variables System](VARIABLES_SYSTEM.md)** - Provides secrets for variable resolution and secure variable management
- **[Templates System](TEMPLATES_SYSTEM.md)** - Provides secrets for template rendering and secure template processing
- **[Actions System](ACTIONS_SYSTEM.md)** - Provides secrets for action running and secure action processing
- **[Logging System](LOGGING_SYSTEM.md)** - Provides secure logging with sensitive data protection
- **[Integrations System](INTEGRATIONS_SYSTEM.md)** - Provides secrets integration through the IntegrationManager
- **[Projects System](PROJECTS_SYSTEM.md)** - Provides project-level secrets management and encryption
- **[Machines System](MACHINES_SYSTEM.md)** - Provides secure machine authentication and credential management
- **[SSH System](SSH_SYSTEM.md)** - Provides secure SSH key and credential management

## Architecture

### Core Components

#### Secrets Manager
- **File**: `internal/secrets/manager.go`
- **Purpose**: Central secrets management with encryption and decryption
- **Features**:
  - Secrets encryption and decryption
  - Key management integration
  - Secrets storage and retrieval
  - Secrets validation and verification
  - Audit logging and compliance
  - Error handling and recovery

#### HCL Processor
- **File**: `internal/secrets/hcl_processor.go`
- **Purpose**: HCL configuration processing with secrets integration
- **Features**:
  - HCL configuration parsing
  - Secrets integration in HCL
  - Configuration validation
  - Secrets resolution
  - Error handling
  - Performance optimization

#### Secrets Integration
- **File**: `internal/secrets/integration.go`
- **Purpose**: Interface implementation for system integration
- **Features**:
  - EncryptSecrets - Encrypt sensitive data
  - DecryptSecrets - Decrypt sensitive data
  - ValidateSecrets - Validate secrets integrity
  - ManageKeys - Manage encryption keys

### Integration Points

#### Variables Integration
- Provides secrets for variable resolution
- Supports encrypted variable values
- Enables secure variable management

#### Templates Integration
- Provides secrets for template rendering
- Supports encrypted template variables
- Enables secure template processing

#### Actions Integration
- Provides secrets for action running
- Supports encrypted action parameters
- Enables secure action processing

#### Configuration Integration
- Provides secrets for configuration management
- Supports encrypted configuration values
- Enables secure configuration handling

## Secrets Types

### Secret Structure
```go
type Secret struct {
    ID              string                 // Secret identifier
    Name            string                 // Secret name
    Value           []byte                 // Encrypted secret value
    Type            string                 // Secret type
    Algorithm       string                 // Encryption algorithm
    KeyID           string                 // Encryption key ID
    Metadata        map[string]interface{} // Secret metadata
    ExpiresAt       *time.Time             // Expiration time
    CreatedAt       time.Time              // Creation timestamp
    UpdatedAt       time.Time              // Update timestamp
}
```

### Secret Collection
```go
type SecretCollection struct {
    ID              string                 // Collection identifier
    Name            string                 // Collection name
    Secrets         map[string]*Secret     // Secret collection
    Algorithm       string                 // Encryption algorithm
    KeyID           string                 // Encryption key ID
    Metadata        map[string]interface{} // Collection metadata
    CreatedAt       time.Time              // Creation timestamp
    UpdatedAt       time.Time              // Update timestamp
}
```

### Encryption Metadata
```go
type EncryptionMetadata struct {
    Algorithm       string                 // Encryption algorithm
    KeyID           string                 // Key identifier
    Version         string                 // Encryption version
    Timestamp       time.Time              // Encryption timestamp
    Recipients      []string               // Recipient identifiers
    Metadata        map[string]interface{} // Additional metadata
}
```

## Secret Categories

### Secret Types
- **Passwords**: User passwords and credentials
- **API Keys**: API keys and tokens
- **Certificates**: SSL/TLS certificates
- **Private Keys**: SSH and encryption keys
- **Database Credentials**: Database connection credentials
- **Service Accounts**: Service account credentials

### Security Levels
- **Low**: Non-critical secrets
- **Medium**: Important secrets
- **High**: Critical secrets
- **Maximum**: Ultra-sensitive secrets

### Access Levels
- **Read-only**: Read-only access to secrets
- **Read-write**: Read-write access to secrets
- **Admin**: Administrative access to secrets
- **Owner**: Full ownership of secrets

## Secrets Management

### Encryption Management
- **Age Encryption**: Use age encryption format
- **Key Management**: External key management
- **Algorithm Selection**: Choose encryption algorithms
- **Key Rotation**: Rotate encryption keys
- **Backup Encryption**: Encrypt secret backups

### Storage Management
- **Secure Storage**: Store secrets securely
- **Access Control**: Control access to secrets
- **Audit Logging**: Log secret access
- **Backup Management**: Manage secret backups
- **Recovery Procedures**: Secret recovery procedures

### Lifecycle Management
- **Secret Creation**: Create new secrets
- **Secret Updates**: Update existing secrets
- **Secret Rotation**: Rotate secrets regularly
- **Secret Expiration**: Handle secret expiration
- **Secret Deletion**: Secure secret deletion

## Secrets Operations

### Encryption Operations
- **Encrypt**: Encrypt sensitive data
- **Decrypt**: Decrypt encrypted data
- **Re-encrypt**: Re-encrypt with new keys
- **Validate**: Validate encryption integrity
- **Verify**: Verify encryption authenticity

### Management Operations
- **Create**: Create new secrets
- **Read**: Read secret values
- **Update**: Update secret values
- **Delete**: Delete secrets
- **List**: List available secrets

### Key Operations
- **Generate Keys**: Generate encryption keys
- **Import Keys**: Import external keys
- **Export Keys**: Export keys for backup
- **Rotate Keys**: Rotate encryption keys
- **Validate Keys**: Validate key integrity

## Security Features

### Encryption Security
- **Age Encryption**: Use age encryption format
- **Key Management**: External key management
- **Algorithm Security**: Secure encryption algorithms
- **Key Rotation**: Regular key rotation
- **Backup Security**: Secure backup encryption

### Access Security
- **Access Control**: Control access to secrets
- **Authentication**: Authenticate secret access
- **Authorization**: Authorize secret operations
- **Audit Logging**: Comprehensive audit logging
- **Compliance**: Compliance with security standards

### Data Protection
- **Data Encryption**: Encrypt all sensitive data
- **Data Validation**: Validate data integrity
- **Data Sanitization**: Sanitize input data
- **Secure Transmission**: Secure data transmission
- **Secure Storage**: Secure data storage

## Performance Features

### Encryption Performance
- **Fast Encryption**: Fast encryption algorithms
- **Parallel Processing**: Parallel encryption operations
- **Caching**: Cache encryption results
- **Optimization**: Optimize encryption performance
- **Monitoring**: Monitor encryption performance

### Storage Performance
- **Efficient Storage**: Efficient secret storage
- **Compression**: Compress secret data
- **Indexing**: Index secret metadata
- **Caching**: Cache frequently accessed secrets
- **Optimization**: Optimize storage performance

### Key Management Performance
- **Fast Key Operations**: Fast key operations
- **Key Caching**: Cache encryption keys
- **Parallel Key Operations**: Parallel key operations
- **Optimization**: Optimize key management
- **Monitoring**: Monitor key management performance

## CLI Commands

### Secrets Management
```bash
# Encrypt secrets
spooky secrets encrypt <file> --output <encrypted-file>

# Decrypt secrets
spooky secrets decrypt <encrypted-file> --output <file>

# List secrets
spooky secrets list

# Show secret information
spooky secrets info <secret-name>
```

### Key Management
```bash
# Generate encryption key
spooky secrets generate-key --output <key-file>

# Import external key
spooky secrets import-key <key-file>

# Export key for backup
spooky secrets export-key <key-id> --output <backup-file>

# List available keys
spooky secrets list-keys
```

### Secrets Validation
```bash
# Validate secrets
spooky secrets validate

# Check secret integrity
spooky secrets check-integrity <secret-name>

# Verify encryption
spooky secrets verify <encrypted-file>
```

### Secrets Operations
```bash
# Create new secret
spooky secrets create <secret-name> --value <secret-value>

# Update secret
spooky secrets update <secret-name> --value <new-value>

# Delete secret
spooky secrets delete <secret-name>

# Rotate secret
spooky secrets rotate <secret-name>
```

## Configuration

### Secrets Configuration
```hcl
# secrets/config.hcl
secrets_config {
  # Encryption settings
  encryption {
    algorithm = "age"
    key_path = "~/.config/spooky/keys"
    backup_keys = true
    rotation_interval = "90d"
  }
  
  # Storage settings
  storage {
    backend = "file"
    path = "~/.local/state/spooky/secrets"
    encryption = true
    compression = true
  }
  
  # Security settings
  security {
    access_control = true
    audit_logging = true
    compliance = true
    backup_encryption = true
  }
  
  # Performance settings
  performance {
    cache_enabled = true
    cache_ttl = "1h"
    parallel_operations = 4
    compression_level = 6
  }
}
```

### Secret Definition
```hcl
# secrets/database-credentials.hcl
secret "database-password" {
  name = "Database Password"
  description = "Main database password"
  
  # Encryption settings
  algorithm = "age"
  key_id = "database-key"
  
  # Security settings
  security_level = "high"
  access_level = "read-write"
  
  # Lifecycle settings
  expires_after = "365d"
  rotation_interval = "90d"
  
  # Metadata
  metadata {
    database = "postgresql"
    environment = "production"
    owner = "db-team"
  }
}

secret "api-key" {
  name = "API Key"
  description = "External API key"
  
  algorithm = "age"
  key_id = "api-key"
  
  security_level = "medium"
  access_level = "read-only"
  
  expires_after = "180d"
  rotation_interval = "60d"
  
  metadata {
    service = "external-api"
    environment = "production"
    owner = "api-team"
  }
}
```

## Examples

### Basic Secret Management
```hcl
# secrets/basic.hcl
secret "web-password" {
  name = "Web Server Password"
  description = "Web server admin password"
  
  algorithm = "age"
  key_id = "web-key"
  
  security_level = "medium"
  access_level = "read-write"
  
  expires_after = "180d"
  rotation_interval = "60d"
}

secret "ssh-key" {
  name = "SSH Private Key"
  description = "SSH private key for server access"
  
  algorithm = "age"
  key_id = "ssh-key"
  
  security_level = "high"
  access_level = "read-only"
  
  expires_after = "365d"
  rotation_interval = "90d"
}
```

### Advanced Secret Management
```hcl
# secrets/advanced.hcl
secret "database-master-password" {
  name = "Database Master Password"
  description = "Master password for database cluster"
  
  algorithm = "age"
  key_id = "database-master-key"
  
  security_level = "maximum"
  access_level = "admin"
  
  expires_after = "365d"
  rotation_interval = "90d"
  
  metadata {
    database = "postgresql"
    cluster = "production-cluster"
    environment = "production"
    owner = "db-admin"
    compliance = ["sox", "pci", "gdpr"]
  }
  
  # Advanced settings
  settings {
    backup_enabled = true
    audit_enabled = true
    compliance_logging = true
    access_notification = true
  }
}

secret "certificate-private-key" {
  name = "SSL Certificate Private Key"
  description = "Private key for SSL certificate"
  
  algorithm = "age"
  key_id = "certificate-key"
  
  security_level = "high"
  access_level = "read-only"
  
  expires_after = "730d"
  rotation_interval = "365d"
  
  metadata {
    certificate_type = "ssl"
    domain = "example.com"
    environment = "production"
    owner = "security-team"
  }
}
```

## Integration Examples

### Variables Integration
```go
// Use secrets in variables
secretsIntegration := manager.GetSecretsIntegration()
variablesIntegration := manager.GetVariablesIntegration()

// Decrypt secret for variable
secret, err := secretsIntegration.DecryptSecret("database-password")
if err != nil {
    return err
}

// Use decrypted secret in variable
variable := &spookytypes.Variable{
    Name: "db_password",
    Value: string(secret.Value),
    Encrypted: false,
}

err = variablesIntegration.SetVariable(variable)
if err != nil {
    return err
}
```

### Templates Integration
```go
// Use secrets in templates
secretsIntegration := manager.GetSecretsIntegration()
templatesIntegration := manager.GetTemplatesIntegration()

// Decrypt secrets for template
secrets := make(map[string]string)
secretNames := []string{"api-key", "database-password"}

for _, secretName := range secretNames {
    secret, err := secretsIntegration.DecryptSecret(secretName)
    if err != nil {
        log.Printf("Failed to decrypt secret %s: %v", secretName, err)
        continue
    }
    secrets[secretName] = string(secret.Value)
}

// Use secrets in template
template, err := templatesIntegration.LoadTemplate("templates/config.tmpl")
if err != nil {
    return err
}

result, err := templatesIntegration.RenderTemplate(template, map[string]interface{}{
    "secrets": secrets,
})
if err != nil {
    return err
}
```

### Actions Integration
```go
// Use secrets in actions
secretsIntegration := manager.GetSecretsIntegration()
actionsIntegration := manager.GetActionsIntegration()

// Decrypt secret for action
secret, err := secretsIntegration.DecryptSecret("ssh-key")
if err != nil {
    return err
}

// Use decrypted secret in action
action := &spookytypes.Action{
    Name: "deploy-application",
    Script: "deploy.sh",
    Secrets: map[string]string{
        "SSH_PRIVATE_KEY": string(secret.Value),
    },
}

result, err := actionsIntegration.RunAction(action)
if err != nil {
    return err
}
```

## Best Practices

### Security
- Use strong encryption algorithms
- Implement proper key management
- Regular key rotation
- Secure secret storage
- Comprehensive audit logging

### Management
- Use descriptive secret names
- Implement proper access controls
- Regular secret rotation
- Monitor secret usage
- Backup secrets securely

### Performance
- Use efficient encryption algorithms
- Implement caching for frequently accessed secrets
- Optimize secret operations
- Monitor performance
- Use appropriate timeouts

### Compliance
- Implement compliance logging
- Regular security audits
- Monitor access patterns
- Implement access controls
- Document security procedures

## Troubleshooting

### Common Issues

#### Encryption Issues
```bash
# Check encryption key
spooky secrets list-keys

# Validate encryption
spooky secrets verify <encrypted-file>

# Test encryption/decryption
spooky secrets encrypt test.txt --output test.enc
spooky secrets decrypt test.enc --output test-decrypted.txt
```

#### Key Management Issues
```bash
# Check key status
spooky secrets list-keys --verbose

# Validate key integrity
spooky secrets validate-key <key-id>

# Import missing key
spooky secrets import-key <key-file>
```

#### Access Issues
```bash
# Check secret access
spooky secrets info <secret-name>

# Validate permissions
spooky secrets validate --permissions

# Check audit logs
spooky secrets audit --secret <secret-name>
```

#### Performance Issues
```bash
# Check secret performance
spooky secrets performance

# Monitor cache usage
spooky secrets cache-status

# Optimize operations
spooky secrets optimize
```

## API Reference

### SecretsIntegration Interface
```go
type SecretsIntegration interface {
    EncryptSecret(ctx context.Context, name string, value []byte) (*spookytypes.Secret, error)
    DecryptSecret(ctx context.Context, name string) (*spookytypes.Secret, error)
    ValidateSecret(ctx context.Context, name string) (*spookytypes.ValidationResult, error)
    ManageKeys(ctx context.Context, operation string, keyData []byte) error
}
```

### Secrets Manager Methods
```go
// Secret management
EncryptSecret(ctx context.Context, name string, value []byte) (*spookytypes.Secret, error)
DecryptSecret(ctx context.Context, name string) (*spookytypes.Secret, error)
ValidateSecret(ctx context.Context, name string) (*spookytypes.ValidationResult, error)

// Key management
ManageKeys(ctx context.Context, operation string, keyData []byte) error
GenerateKey(ctx context.Context, algorithm string) (*spookytypes.Key, error)
ImportKey(ctx context.Context, keyData []byte) error
ExportKey(ctx context.Context, keyID string) ([]byte, error)

// Secret operations
CreateSecret(ctx context.Context, secret *spookytypes.Secret) error
GetSecret(ctx context.Context, name string) (*spookytypes.Secret, error)
UpdateSecret(ctx context.Context, secret *spookytypes.Secret) error
DeleteSecret(ctx context.Context, name string) error
ListSecrets(ctx context.Context) ([]*spookytypes.Secret, error)
```

## Related Documentation

- [Secrets API Reference](SECRETS_API_REFERENCE.md) - Complete API documentation
- [Secrets User Guide](SECRETS_USER_GUIDE.md) - User guide and examples
- [Secrets Troubleshooting](SECRETS_TROUBLESHOOTING.md) - Troubleshooting guide
- [Variables System](VARIABLES_SYSTEM.md) - Variables integration
- [Templates System](TEMPLATES_SYSTEM.md) - Templates integration
- [Actions System](ACTIONS_SYSTEM.md) - Actions integration
- [Configuration System](CONFIGURATION_SYSTEM.md) - Configuration integration
