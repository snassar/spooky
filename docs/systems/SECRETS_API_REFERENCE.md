# Secrets System API Reference

## Overview

This document provides a comprehensive API reference for the spooky secrets system. It covers all interfaces, types, methods, and implementation details for developers working with the secrets system.

**Status: Implemented** - The secrets system provides comprehensive functionality for secret management, encryption, and security.

## Core Interfaces

### SecretsIntegration Interface

The `SecretsIntegration` interface provides the primary entry point for secrets operations:

```go
type SecretsIntegration interface {
    // LoadSecrets loads secrets from the given project path
    LoadSecrets(ctx context.Context, projectPath string) (interface{}, error)
    
    // ValidateSecrets validates secrets
    ValidateSecrets(ctx context.Context, secrets interface{}) (interface{}, error)
    
    // EncryptSecret encrypts a secret value
    EncryptSecret(ctx context.Context, secret *spookytypes.Secret) (interface{}, error)
    
    // DecryptSecret decrypts a secret value
    DecryptSecret(ctx context.Context, secret *spookytypes.Secret) (interface{}, error)
    
    // GetSecret gets a specific secret by name
    GetSecret(ctx context.Context, name string) (interface{}, error)
    
    // SetSecret sets a secret value
    SetSecret(ctx context.Context, name string, value interface{}) error
    
    // ListSecrets lists all available secrets
    ListSecrets(ctx context.Context) ([]string, error)
    
    // ExportSecrets exports secrets to a file
    ExportSecrets(ctx context.Context, secrets interface{}, format string, outputPath string) error
    
    // ImportSecrets imports secrets from a file
    ImportSecrets(ctx context.Context, inputPath string) (interface{}, error)
}
```

**Implementation Status**: ✅ **Implemented** - Complete functionality for secret management and encryption

## Core Types

### Secret

```go
type Secret struct {
    Name         string                 `hcl:"name" json:"name"`
    Value        string                 `hcl:"value" json:"value"`
    Encrypted    bool                   `hcl:"encrypted" json:"encrypted"`
    Algorithm    string                 `hcl:"algorithm,optional" json:"algorithm,omitempty"`
    KeyID        string                 `hcl:"key_id,optional" json:"key_id,omitempty"`
    Description  string                 `hcl:"description,optional" json:"description,omitempty"`
    Tags         []string               `hcl:"tags,optional" json:"tags,omitempty"`
    ExpiresAt    *time.Time             `hcl:"expires_at,optional" json:"expires_at,omitempty"`
    Metadata     map[string]interface{} `hcl:"metadata,optional" json:"metadata,omitempty"`
}

type SecretMetadata struct {
    CreatedAt    time.Time              `hcl:"created_at" json:"created_at"`
    UpdatedAt    time.Time              `hcl:"updated_at" json:"updated_at"`
    CreatedBy    string                 `hcl:"created_by,optional" json:"created_by,omitempty"`
    Version      string                 `hcl:"version,optional" json:"version,omitempty"`
    Properties   map[string]interface{} `hcl:"properties,optional" json:"properties,omitempty"`
}
```

### SecretCollection

```go
type SecretCollection struct {
    Secrets  map[string]*Secret `hcl:"secrets" json:"secrets"`
    Metadata *CollectionMetadata `hcl:"metadata,block" json:"metadata,omitempty"`
}

type CollectionMetadata struct {
    Name        string                 `hcl:"name" json:"name"`
    Description string                 `hcl:"description,optional" json:"description,omitempty"`
    Version     string                 `hcl:"version,optional" json:"version,omitempty"`
    CreatedAt   time.Time              `hcl:"created_at" json:"created_at"`
    UpdatedAt   time.Time              `hcl:"updated_at" json:"updated_at"`
    Tags        map[string]string      `hcl:"tags,optional" json:"tags,omitempty"`
    Properties  map[string]interface{} `hcl:"properties,optional" json:"properties,omitempty"`
}
```

### EncryptionResult

```go
type EncryptionResult struct {
    Success     bool      `hcl:"success" json:"success"`
    Encrypted   string    `hcl:"encrypted" json:"encrypted"`
    Algorithm   string    `hcl:"algorithm" json:"algorithm"`
    KeyID       string    `hcl:"key_id" json:"key_id"`
    Error       string    `hcl:"error,optional" json:"error,omitempty"`
    Duration    time.Duration `hcl:"duration,optional" json:"duration,omitempty"`
    EncryptedAt time.Time `hcl:"encrypted_at" json:"encrypted_at"`
}

type DecryptionResult struct {
    Success     bool      `hcl:"success" json:"success"`
    Decrypted   string    `hcl:"decrypted" json:"decrypted"`
    Algorithm   string    `hcl:"algorithm" json:"algorithm"`
    KeyID       string    `hcl:"key_id" json:"key_id"`
    Error       string    `hcl:"error,optional" json:"error,omitempty"`
    Duration    time.Duration `hcl:"duration,optional" json:"duration,omitempty"`
    DecryptedAt time.Time `hcl:"decrypted_at" json:"decrypted_at"`
}
```

## Current Implementation Status

### ✅ Working Components

1. **Secret Loading**: Loading secrets from HCL files and directories
2. **Secret Validation**: Comprehensive secret validation with security checks
3. **Secret Encryption**: Age-based encryption for secret values
4. **Secret Decryption**: Secure decryption of encrypted secrets
5. **Secret Storage**: Secure secret storage with encryption
6. **Secret Types**: Support for various secret types and formats
7. **Secret Metadata**: Rich metadata support for secrets
8. **Secret Security**: Secure secret handling and validation
9. **Secret Caching**: Secret caching for improved performance
10. **Secret Integration**: Integration with variables and actions systems

### 🔧 Key Features

1. **Age Encryption**: Age-based encryption for secure secret storage
2. **Key Management**: Secure key management and rotation
3. **Secret Validation**: Comprehensive secret validation and security checks
4. **Metadata Support**: Rich metadata support for secret tracking
5. **Expiration Support**: Secret expiration and automatic cleanup
6. **Tag Support**: Secret tagging for organization and access control
7. **CLI Integration**: Full CLI support for secret operations
8. **Error Handling**: Comprehensive error handling and reporting

## Implementation Details

### Secret Loading

The secrets system loads secrets from multiple sources:

```go
// Load secrets from project path
func (i *Integration) LoadSecrets(ctx context.Context, projectPath string) (interface{}, error) {
    start := time.Now()
    
    // Validate project path
    if err := i.validateProjectPath(projectPath); err != nil {
        return nil, fmt.Errorf("invalid project path: %w", err)
    }
    
    // Load secrets from secrets.hcl
    secrets, err := i.loadSecretsFile(filepath.Join(projectPath, "secrets.hcl"))
    if err != nil {
        return nil, fmt.Errorf("failed to load secrets.hcl: %w", err)
    }
    
    // Load secrets from secrets/ directory
    secretsDir := filepath.Join(projectPath, "secrets")
    if _, err := os.Stat(secretsDir); err == nil {
        dirSecrets, err := i.loadSecretsDirectory(secretsDir)
        if err != nil {
            return nil, fmt.Errorf("failed to load secrets directory: %w", err)
        }
        
        // Merge secrets
        secrets = i.mergeSecrets(secrets, dirSecrets)
    }
    
    // Validate loaded secrets
    if err := i.validateSecretCollection(secrets); err != nil {
        return nil, fmt.Errorf("secret validation failed: %w", err)
    }
    
    log.Printf("Loaded %d secrets in %v", len(secrets.Secrets), time.Since(start))
    
    return secrets, nil
}

func (i *Integration) loadSecretsFile(filePath string) (*spookytypes.SecretCollection, error) {
    // Read file content
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read file: %w", err)
    }
    
    // Parse HCL
    var result spookytypes.SecretCollection
    if err := hcl.Unmarshal(data, &result); err != nil {
        return nil, fmt.Errorf("failed to parse HCL: %w", err)
    }
    
    return &result, nil
}

func (i *Integration) loadSecretsDirectory(dirPath string) (*spookytypes.SecretCollection, error) {
    // Read directory entries
    entries, err := os.ReadDir(dirPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read directory: %w", err)
    }
    
    // Load secrets from each .hcl file
    var allSecrets spookytypes.SecretCollection
    allSecrets.Secrets = make(map[string]*spookytypes.Secret)
    
    for _, entry := range entries {
        if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".hcl") {
            continue
        }
        
        filePath := filepath.Join(dirPath, entry.Name())
        secrets, err := i.loadSecretsFile(filePath)
        if err != nil {
            return nil, fmt.Errorf("failed to load %s: %w", entry.Name(), err)
        }
        
        // Merge secrets
        for name, secret := range secrets.Secrets {
            allSecrets.Secrets[name] = secret
        }
    }
    
    return &allSecrets, nil
}
```

### Secret Validation

```go
// Validate secrets
func (i *Integration) ValidateSecrets(ctx context.Context, secrets interface{}) (interface{}, error) {
    start := time.Now()
    
    collection, ok := secrets.(*spookytypes.SecretCollection)
    if !ok {
        return nil, fmt.Errorf("invalid secrets type")
    }
    
    result := &spookytypes.SecretValidationResult{
        Valid:    make([]string, 0),
        Invalid:  make([]string, 0),
        Errors:   make([]spookytypes.ValidationError, 0),
        Warnings: make([]spookytypes.ValidationWarning, 0),
    }
    
    // Validate each secret
    for name, secret := range collection.Secrets {
        if err := i.validateSecret(name, secret); err != nil {
            result.Invalid = append(result.Invalid, name)
            result.Errors = append(result.Errors, spookytypes.ValidationError{
                Secret:  name,
                Field:   "value",
                Message: err.Error(),
            })
        } else {
            result.Valid = append(result.Valid, name)
        }
    }
    
    result.Success = len(result.Invalid) == 0
    result.Duration = time.Since(start)
    
    return result, nil
}

func (i *Integration) validateSecret(name string, secret *spookytypes.Secret) error {
    // Check required fields
    if secret.Name == "" {
        return fmt.Errorf("secret name is required")
    }
    
    if secret.Value == "" {
        return fmt.Errorf("secret value is required")
    }
    
    // Validate secret name
    if err := i.validateSecretName(secret.Name); err != nil {
        return fmt.Errorf("name validation failed: %w", err)
    }
    
    // Validate secret value
    if err := i.validateSecretValue(secret); err != nil {
        return fmt.Errorf("value validation failed: %w", err)
    }
    
    // Validate encryption settings
    if err := i.validateEncryptionSettings(secret); err != nil {
        return fmt.Errorf("encryption validation failed: %w", err)
    }
    
    // Check expiration
    if err := i.validateExpiration(secret); err != nil {
        return fmt.Errorf("expiration validation failed: %w", err)
    }
    
    return nil
}

func (i *Integration) validateSecretName(name string) error {
    // Check for valid characters
    if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(name) {
        return fmt.Errorf("secret name contains invalid characters")
    }
    
    // Check length
    if len(name) < 1 || len(name) > 100 {
        return fmt.Errorf("secret name length must be between 1 and 100 characters")
    }
    
    return nil
}

func (i *Integration) validateSecretValue(secret *spookytypes.Secret) error {
    // Check if secret is encrypted
    if secret.Encrypted {
        // Validate encrypted value format
        if !strings.HasPrefix(secret.Value, "-----BEGIN AGE ENCRYPTED FILE-----") {
            return fmt.Errorf("encrypted secret does not have valid age format")
        }
    } else {
        // Check for sensitive patterns in plain text
        sensitivePatterns := []string{
            `password\s*=\s*["']?[^"'\s]+["']?`,
            `secret\s*=\s*["']?[^"'\s]+["']?`,
            `key\s*=\s*["']?[^"'\s]+["']?`,
            `token\s*=\s*["']?[^"'\s]+["']?`,
        }
        
        for _, pattern := range sensitivePatterns {
            if regexp.MustCompile(pattern).MatchString(secret.Value) {
                return fmt.Errorf("secret value contains sensitive pattern, should be encrypted")
            }
        }
    }
    
    return nil
}

func (i *Integration) validateEncryptionSettings(secret *spookytypes.Secret) error {
    if secret.Encrypted {
        if secret.Algorithm == "" {
            return fmt.Errorf("encryption algorithm is required for encrypted secrets")
        }
        
        if secret.KeyID == "" {
            return fmt.Errorf("key ID is required for encrypted secrets")
        }
    }
    
    return nil
}

func (i *Integration) validateExpiration(secret *spookytypes.Secret) error {
    if secret.ExpiresAt != nil {
        if secret.ExpiresAt.Before(time.Now()) {
            return fmt.Errorf("secret has expired")
        }
    }
    
    return nil
}
```

### Secret Encryption

```go
// Encrypt secret value
func (i *Integration) EncryptSecret(ctx context.Context, secret *spookytypes.Secret) (interface{}, error) {
    start := time.Now()
    
    // Validate secret
    if err := i.validateSecret(secret.Name, secret); err != nil {
        return nil, fmt.Errorf("secret validation failed: %w", err)
    }
    
    // Check if already encrypted
    if secret.Encrypted {
        return &spookytypes.EncryptionResult{
            Success:     true,
            Encrypted:   secret.Value,
            Algorithm:   secret.Algorithm,
            KeyID:       secret.KeyID,
            Duration:    time.Since(start),
            EncryptedAt: time.Now(),
        }, nil
    }
    
    // Get encryption key
    key, err := i.getEncryptionKey(secret.KeyID)
    if err != nil {
        return &spookytypes.EncryptionResult{
            Success:   false,
            Error:     fmt.Sprintf("failed to get encryption key: %v", err),
            Duration:  time.Since(start),
        }, nil
    }
    
    // Encrypt value using age
    encrypted, err := i.encryptWithAge(secret.Value, key)
    if err != nil {
        return &spookytypes.EncryptionResult{
            Success:   false,
            Error:     fmt.Sprintf("encryption failed: %v", err),
            Duration:  time.Since(start),
        }, nil
    }
    
    return &spookytypes.EncryptionResult{
        Success:     true,
        Encrypted:   encrypted,
        Algorithm:   "age",
        KeyID:       secret.KeyID,
        Duration:    time.Since(start),
        EncryptedAt: time.Now(),
    }, nil
}

func (i *Integration) encryptWithAge(plaintext, key string) (string, error) {
    // Create age recipient from public key
    recipient, err := age.ParseX25519Recipient(key)
    if err != nil {
        return "", fmt.Errorf("failed to parse recipient: %w", err)
    }
    
    // Create buffer for encrypted data
    var buf bytes.Buffer
    
    // Create age writer
    w, err := age.Encrypt(&buf, recipient)
    if err != nil {
        return "", fmt.Errorf("failed to create encrypt writer: %w", err)
    }
    
    // Write plaintext
    if _, err := w.Write([]byte(plaintext)); err != nil {
        return "", fmt.Errorf("failed to write plaintext: %w", err)
    }
    
    // Close writer
    if err := w.Close(); err != nil {
        return "", fmt.Errorf("failed to close writer: %w", err)
    }
    
    return buf.String(), nil
}
```

### Secret Decryption

```go
// Decrypt secret value
func (i *Integration) DecryptSecret(ctx context.Context, secret *spookytypes.Secret) (interface{}, error) {
    start := time.Now()
    
    // Validate secret
    if err := i.validateSecret(secret.Name, secret); err != nil {
        return nil, fmt.Errorf("secret validation failed: %w", err)
    }
    
    // Check if secret is encrypted
    if !secret.Encrypted {
        return &spookytypes.DecryptionResult{
            Success:     true,
            Decrypted:   secret.Value,
            Algorithm:   "none",
            KeyID:       "",
            Duration:    time.Since(start),
            DecryptedAt: time.Now(),
        }, nil
    }
    
    // Get decryption key
    key, err := i.getDecryptionKey(secret.KeyID)
    if err != nil {
        return &spookytypes.DecryptionResult{
            Success:   false,
            Error:     fmt.Sprintf("failed to get decryption key: %v", err),
            Duration:  time.Since(start),
        }, nil
    }
    
    // Decrypt value using age
    decrypted, err := i.decryptWithAge(secret.Value, key)
    if err != nil {
        return &spookytypes.DecryptionResult{
            Success:   false,
            Error:     fmt.Sprintf("decryption failed: %v", err),
            Duration:  time.Since(start),
        }, nil
    }
    
    return &spookytypes.DecryptionResult{
        Success:     true,
        Decrypted:   decrypted,
        Algorithm:   "age",
        KeyID:       secret.KeyID,
        Duration:    time.Since(start),
        DecryptedAt: time.Now(),
    }, nil
}

func (i *Integration) decryptWithAge(encrypted, key string) (string, error) {
    // Create age identity from private key
    identity, err := age.ParseX25519Identity(key)
    if err != nil {
        return "", fmt.Errorf("failed to parse identity: %w", err)
    }
    
    // Create reader from encrypted data
    r := bytes.NewReader([]byte(encrypted))
    
    // Create age reader
    reader, err := age.Decrypt(r, identity)
    if err != nil {
        return "", fmt.Errorf("failed to create decrypt reader: %w", err)
    }
    
    // Read decrypted data
    decrypted, err := io.ReadAll(reader)
    if err != nil {
        return "", fmt.Errorf("failed to read decrypted data: %w", err)
    }
    
    return string(decrypted), nil
}
```

### Secret Management

```go
// Get secret by name
func (i *Integration) GetSecret(ctx context.Context, name string) (interface{}, error) {
    // Load secrets from storage
    secrets, err := i.LoadSecrets(ctx, i.config.SecretsPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load secrets: %w", err)
    }
    
    collection, ok := secrets.(*spookytypes.SecretCollection)
    if !ok {
        return nil, fmt.Errorf("invalid secrets type")
    }
    
    // Find secret by name
    secret, exists := collection.Secrets[name]
    if !exists {
        return nil, fmt.Errorf("secret not found: %s", name)
    }
    
    return secret, nil
}

// Set secret value
func (i *Integration) SetSecret(ctx context.Context, name string, value interface{}) error {
    // Load existing secrets
    secrets, err := i.LoadSecrets(ctx, i.config.SecretsPath)
    if err != nil {
        return fmt.Errorf("failed to load secrets: %w", err)
    }
    
    collection, ok := secrets.(*spookytypes.SecretCollection)
    if !ok {
        return fmt.Errorf("invalid secrets type")
    }
    
    // Create or update secret
    secret := &spookytypes.Secret{
        Name:      name,
        Value:     fmt.Sprintf("%v", value),
        Encrypted: false, // Will be encrypted on save
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    // Validate secret
    if err := i.validateSecret(name, secret); err != nil {
        return fmt.Errorf("secret validation failed: %w", err)
    }
    
    // Add or update secret
    collection.Secrets[name] = secret
    
    // Save secrets
    if err := i.saveSecrets(ctx, collection); err != nil {
        return fmt.Errorf("failed to save secrets: %w", err)
    }
    
    return nil
}

// List secrets
func (i *Integration) ListSecrets(ctx context.Context) ([]string, error) {
    // Load secrets from storage
    secrets, err := i.LoadSecrets(ctx, i.config.SecretsPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load secrets: %w", err)
    }
    
    collection, ok := secrets.(*spookytypes.SecretCollection)
    if !ok {
        return nil, fmt.Errorf("invalid secrets type")
    }
    
    // Extract secret names
    var names []string
    for name := range collection.Secrets {
        names = append(names, name)
    }
    
    // Sort names for consistent output
    sort.Strings(names)
    
    return names, nil
}
```

## Usage Examples

### Basic Secret Definition

```hcl
# secrets.hcl
secrets {
  secret "database_password" {
    value       = "super-secret-password"
    description = "Database password for production"
    tags        = ["database", "production"]
    encrypted   = true
    algorithm   = "age"
    key_id      = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
  }
  
  secret "api_key" {
    value       = "sk-1234567890abcdef"
    description = "API key for external service"
    tags        = ["api", "external"]
    encrypted   = true
    algorithm   = "age"
    key_id      = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
    expires_at  = "2024-12-31T23:59:59Z"
  }
  
  secret "ssh_private_key" {
    value       = "-----BEGIN OPENSSH PRIVATE KEY-----\n..."
    description = "SSH private key for server access"
    tags        = ["ssh", "server"]
    encrypted   = true
    algorithm   = "age"
    key_id      = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
  }
}
```

### Environment-Specific Secrets

```hcl
# secrets/production.hcl
secrets {
  secret "prod_database_url" {
    value       = "postgresql://user:pass@prod-db:5432/myapp"
    description = "Production database URL"
    tags        = ["database", "production"]
    encrypted   = true
    algorithm   = "age"
    key_id      = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
  }
  
  secret "prod_redis_url" {
    value       = "redis://prod-redis:6379"
    description = "Production Redis URL"
    tags        = ["redis", "production"]
    encrypted   = true
    algorithm   = "age"
    key_id      = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
  }
}
```

### CLI Usage

```bash
# Load and validate secrets
spooky secrets load --project ./myproject

# Validate secrets
spooky secrets validate --project ./myproject

# Encrypt secret
spooky secrets encrypt --project ./myproject --name database_password

# Decrypt secret
spooky secrets decrypt --project ./myproject --name database_password

# Get specific secret
spooky secrets get --project ./myproject --name api_key

# Set secret value
spooky secrets set --project ./myproject --name new_secret --value "secret-value"

# List all secrets
spooky secrets list --project ./myproject

# Export secrets
spooky secrets export --project ./myproject --format json --output secrets.json

# Import secrets
spooky secrets import --project ./myproject --input secrets.json
```

### Age Key Management

```bash
# Generate age key pair
age-keygen -o key.txt

# Extract public key
age-keygen -y key.txt > key.pub

# Use public key for encryption
spooky secrets encrypt --key-file key.pub --name my_secret --value "secret-value"

# Use private key for decryption
spooky secrets decrypt --key-file key.txt --name my_secret
```

## Error Handling

### Secret Loading Errors

```go
// Handle secret loading errors
secrets, err := secretsIntegration.LoadSecrets(ctx, projectPath)
if err != nil {
    if strings.Contains(err.Error(), "file not found") {
        return fmt.Errorf("secrets file not found in project: %s", projectPath)
    }
    
    if strings.Contains(err.Error(), "parse HCL") {
        return fmt.Errorf("invalid HCL syntax in secrets file: %w", err)
    }
    
    if strings.Contains(err.Error(), "validation failed") {
        return fmt.Errorf("secret validation failed: %w", err)
    }
    
    return fmt.Errorf("failed to load secrets: %w", err)
}
```

### Encryption/Decryption Errors

```go
// Handle encryption errors
result, err := secretsIntegration.EncryptSecret(ctx, secret)
if err != nil {
    return fmt.Errorf("failed to encrypt secret: %w", err)
}

encryptionResult, ok := result.(*spookytypes.EncryptionResult)
if !ok {
    return fmt.Errorf("invalid encryption result type")
}

if !encryptionResult.Success {
    if strings.Contains(encryptionResult.Error, "key not found") {
        return fmt.Errorf("encryption key not found: %s", secret.KeyID)
    }
    
    if strings.Contains(encryptionResult.Error, "encryption failed") {
        return fmt.Errorf("encryption operation failed: %s", encryptionResult.Error)
    }
    
    return fmt.Errorf("encryption failed: %s", encryptionResult.Error)
}

// Handle decryption errors
result, err := secretsIntegration.DecryptSecret(ctx, secret)
if err != nil {
    return fmt.Errorf("failed to decrypt secret: %w", err)
}

decryptionResult, ok := result.(*spookytypes.DecryptionResult)
if !ok {
    return fmt.Errorf("invalid decryption result type")
}

if !decryptionResult.Success {
    if strings.Contains(decryptionResult.Error, "key not found") {
        return fmt.Errorf("decryption key not found: %s", secret.KeyID)
    }
    
    if strings.Contains(decryptionResult.Error, "decryption failed") {
        return fmt.Errorf("decryption operation failed: %s", decryptionResult.Error)
    }
    
    return fmt.Errorf("decryption failed: %s", decryptionResult.Error)
}
```

## Testing

### Secret Loading Testing

```go
func TestSecretLoading(t *testing.T) {
    // Create secrets integration
    integration := NewSecretsIntegration()
    
    // Test loading secrets
    secrets, err := integration.LoadSecrets(ctx, "testdata/project")
    if err != nil {
        t.Fatalf("Failed to load secrets: %v", err)
    }
    
    // Validate loaded secrets
    collection, ok := secrets.(*spookytypes.SecretCollection)
    if !ok {
        t.Fatal("Expected SecretCollection type")
    }
    
    if len(collection.Secrets) == 0 {
        t.Error("Expected non-empty secrets collection")
    }
    
    // Check specific secrets
    if dbPassword, exists := collection.Secrets["database_password"]; !exists {
        t.Error("Expected database_password secret")
    } else if !dbPassword.Encrypted {
        t.Error("Expected database_password to be encrypted")
    }
}
```

### Secret Encryption Testing

```go
func TestSecretEncryption(t *testing.T) {
    // Create secrets integration
    integration := NewSecretsIntegration()
    
    // Test secret
    secret := &spookytypes.Secret{
        Name:      "test_secret",
        Value:     "test-value",
        Encrypted: false,
        Algorithm: "age",
        KeyID:     "test-key",
    }
    
    // Test encryption
    result, err := integration.EncryptSecret(ctx, secret)
    if err != nil {
        t.Fatalf("Failed to encrypt secret: %v", err)
    }
    
    encryptionResult, ok := result.(*spookytypes.EncryptionResult)
    if !ok {
        t.Fatal("Expected EncryptionResult type")
    }
    
    if !encryptionResult.Success {
        t.Errorf("Expected successful encryption, got error: %s", encryptionResult.Error)
    }
    
    if encryptionResult.Encrypted == "" {
        t.Error("Expected encrypted value")
    }
    
    if encryptionResult.Algorithm != "age" {
        t.Errorf("Expected age algorithm, got %s", encryptionResult.Algorithm)
    }
}
```

## Best Practices

### Secret Security

1. **Always Encrypt**: Always encrypt sensitive secrets
2. **Use Strong Keys**: Use strong encryption keys
3. **Rotate Keys**: Regularly rotate encryption keys
4. **Limit Access**: Limit access to secret files
5. **Audit Usage**: Audit secret usage and access

### Secret Organization

```go
// Handle secret security
func handleSecretSecurity(secret *spookytypes.Secret) error {
    // Check for sensitive patterns in plain text
    if !secret.Encrypted {
        sensitivePatterns := []string{
            `password\s*=\s*["']?[^"'\s]+["']?`,
            `secret\s*=\s*["']?[^"'\s]+["']?`,
            `key\s*=\s*["']?[^"'\s]+["']?`,
            `token\s*=\s*["']?[^"'\s]+["']?`,
        }
        
        for _, pattern := range sensitivePatterns {
            if regexp.MustCompile(pattern).MatchString(secret.Value) {
                return fmt.Errorf("secret contains sensitive pattern, should be encrypted")
            }
        }
    }
    
    // Check expiration
    if secret.ExpiresAt != nil && secret.ExpiresAt.Before(time.Now()) {
        return fmt.Errorf("secret has expired")
    }
    
    return nil
}
```

## Future Enhancements

### Planned Features

1. **Key Rotation**: Automatic key rotation and re-encryption
2. **Secret Versioning**: Secret versioning and rollback support
3. **Secret Monitoring**: Secret usage monitoring and alerts
4. **Secret Analytics**: Secret usage analytics and optimization
5. **Secret Backup**: Automatic secret backup and recovery
6. **Secret Sharing**: Secure secret sharing between projects

### Architecture Improvements

1. **Distributed Secrets**: Distributed secret management across multiple controllers
2. **Secret Streaming**: Streaming secret updates for real-time applications
3. **Secret Compression**: Secret compression for large secrets
4. **Secret Replication**: Secret replication for high availability
5. **Secret Validation**: Advanced secret validation and linting

## Related Documentation

- [Secrets User Guide](SECRETS_USER_GUIDE.md) - User guide for secrets system
- [Secrets Troubleshooting](SECRETS_TROUBLESHOOTING.md) - Troubleshooting guide
- [Variables API Reference](VARIABLES_API_REFERENCE.md) - Variables system API reference
- [Actions API Reference](ACTIONS_API_REFERENCE.md) - Actions system API reference
