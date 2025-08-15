# Configuration Encryption for Sensitive Data

## Overview

The spooky codebase currently stores sensitive configuration data in plain text, including SSH keys, passwords, and other credentials. This document outlines recommendations for implementing configuration encryption to improve security and protect sensitive data at rest.

## Current State

- Sensitive configuration data stored in plain text
- SSH private keys and passphrases unencrypted
- No encryption for configuration files
- Limited protection for credentials
- No key management for encrypted configurations

## Recommendations

### 1. Configuration Encryption Interface

Define a comprehensive configuration encryption interface:

```go
type ConfigurationEncryption interface {
    EncryptConfig(config *spookytypes.Config, key []byte) (*EncryptedConfig, error)
    DecryptConfig(encrypted *EncryptedConfig, key []byte) (*spookytypes.Config, error)
    EncryptSensitiveField(value string, key []byte) (string, error)
    DecryptSensitiveField(encryptedValue string, key []byte) (string, error)
    GenerateEncryptionKey() ([]byte, error)
    ValidateEncryptionKey(key []byte) error
    RotateEncryptionKey(oldKey []byte, newKey []byte) error
}

type EncryptedConfig struct {
    Version       string            `json:"version"`
    Algorithm     string            `json:"algorithm"`
    Salt          []byte            `json:"salt"`
    Nonce         []byte            `json:"nonce"`
    EncryptedData []byte            `json:"encrypted_data"`
    Metadata      map[string]string `json:"metadata"`
    CreatedAt     time.Time         `json:"created_at"`
    LastModified  time.Time         `json:"last_modified"`
}

type EncryptionConfig struct {
    Algorithm     string            `hcl:"algorithm"`
    KeySize       int               `hcl:"key_size"`
    SaltSize      int               `hcl:"salt_size"`
    NonceSize     int               `hcl:"nonce_size"`
    Iterations    int               `hcl:"iterations"`
    KeyDerivation string            `hcl:"key_derivation"`
    SensitiveFields []string        `hcl:"sensitive_fields"`
}
```

### 2. AES-GCM Encryption Implementation

Implement AES-GCM encryption for configuration data:

```go
type AESGCMEncryption struct {
    config EncryptionConfig
    logger spookylogging.Logger
}

func NewAESGCMEncryption(config EncryptionConfig) *AESGCMEncryption {
    return &AESGCMEncryption{
        config: config,
        logger: spookylogging.GetLogger(),
    }
}

func (ae *AESGCMEncryption) EncryptConfig(config *spookytypes.Config, key []byte) (*EncryptedConfig, error) {
    // Validate encryption key
    if err := ae.ValidateEncryptionKey(key); err != nil {
        return nil, fmt.Errorf("invalid encryption key: %w", err)
    }
    
    // Serialize configuration to JSON
    configData, err := json.Marshal(config)
    if err != nil {
        return nil, fmt.Errorf("failed to serialize config: %w", err)
    }
    
    // Generate salt and nonce
    salt := make([]byte, ae.config.SaltSize)
    if _, err := rand.Read(salt); err != nil {
        return nil, fmt.Errorf("failed to generate salt: %w", err)
    }
    
    nonce := make([]byte, ae.config.NonceSize)
    if _, err := rand.Read(nonce); err != nil {
        return nil, fmt.Errorf("failed to generate nonce: %w", err)
    }
    
    // Derive encryption key
    derivedKey, err := ae.deriveKey(key, salt)
    if err != nil {
        return nil, fmt.Errorf("failed to derive key: %w", err)
    }
    
    // Create AES-GCM cipher
    block, err := aes.NewCipher(derivedKey)
    if err != nil {
        return nil, fmt.Errorf("failed to create cipher: %w", err)
    }
    
    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM: %w", err)
    }
    
    // Encrypt configuration data
    encryptedData := aesgcm.Seal(nil, nonce, configData, nil)
    
    return &EncryptedConfig{
        Version:       "1.0",
        Algorithm:     "AES-GCM",
        Salt:          salt,
        Nonce:         nonce,
        EncryptedData: encryptedData,
        Metadata: map[string]string{
            "key_size":    strconv.Itoa(ae.config.KeySize),
            "salt_size":   strconv.Itoa(ae.config.SaltSize),
            "nonce_size":  strconv.Itoa(ae.config.NonceSize),
            "iterations":  strconv.Itoa(ae.config.Iterations),
        },
        CreatedAt:    time.Now(),
        LastModified: time.Now(),
    }, nil
}

func (ae *AESGCMEncryption) DecryptConfig(encrypted *EncryptedConfig, key []byte) (*spookytypes.Config, error) {
    // Validate encryption key
    if err := ae.ValidateEncryptionKey(key); err != nil {
        return nil, fmt.Errorf("invalid encryption key: %w", err)
    }
    
    // Derive encryption key
    derivedKey, err := ae.deriveKey(key, encrypted.Salt)
    if err != nil {
        return nil, fmt.Errorf("failed to derive key: %w", err)
    }
    
    // Create AES-GCM cipher
    block, err := aes.NewCipher(derivedKey)
    if err != nil {
        return nil, fmt.Errorf("failed to create cipher: %w", err)
    }
    
    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM: %w", err)
    }
    
    // Decrypt configuration data
    decryptedData, err := aesgcm.Open(nil, encrypted.Nonce, encrypted.EncryptedData, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt data: %w", err)
    }
    
    // Deserialize configuration
    var config spookytypes.Config
    if err := json.Unmarshal(decryptedData, &config); err != nil {
        return nil, fmt.Errorf("failed to deserialize config: %w", err)
    }
    
    return &config, nil
}

func (ae *AESGCMEncryption) deriveKey(key []byte, salt []byte) ([]byte, error) {
    switch ae.config.KeyDerivation {
    case "pbkdf2":
        return pbkdf2.Key(key, salt, ae.config.Iterations, ae.config.KeySize/8, sha256.New)
    case "scrypt":
        return scrypt.Key(key, salt, 16384, 8, 1, ae.config.KeySize/8)
    default:
        return nil, fmt.Errorf("unsupported key derivation: %s", ae.config.KeyDerivation)
    }
}
```

### 3. Selective Field Encryption

Implement selective encryption for sensitive fields:

```go
type SelectiveEncryption struct {
    encryption ConfigurationEncryption
    config     EncryptionConfig
    logger     spookylogging.Logger
}

func (se *SelectiveEncryption) EncryptSensitiveFields(config *spookytypes.Config, key []byte) error {
    // Encrypt SSH authentication data
    if config.Machines != nil {
        for i, machine := range config.Machines {
            if machine.Authentication != nil {
                if err := se.encryptSSHAuthentication(machine.Authentication, key); err != nil {
                    return fmt.Errorf("failed to encrypt SSH authentication for machine %d: %w", i, err)
                }
            }
        }
    }
    
    // Encrypt secrets configuration
    if config.Secrets != nil {
        if err := se.encryptSecretsConfig(config.Secrets, key); err != nil {
            return fmt.Errorf("failed to encrypt secrets config: %w", err)
        }
    }
    
    // Encrypt other sensitive fields
    if err := se.encryptOtherSensitiveFields(config, key); err != nil {
        return fmt.Errorf("failed to encrypt other sensitive fields: %w", err)
    }
    
    return nil
}

func (se *SelectiveEncryption) encryptSSHAuthentication(auth *spookytypes.SSHAuthentication, key []byte) error {
    if auth.PrivateKeyPath != "" {
        encryptedPath, err := se.encryption.EncryptSensitiveField(auth.PrivateKeyPath, key)
        if err != nil {
            return fmt.Errorf("failed to encrypt private key path: %w", err)
        }
        auth.PrivateKeyPath = encryptedPath
    }
    
    if auth.Passphrase != "" {
        encryptedPassphrase, err := se.encryption.EncryptSensitiveField(auth.Passphrase, key)
        if err != nil {
            return fmt.Errorf("failed to encrypt passphrase: %w", err)
        }
        auth.Passphrase = encryptedPassphrase
    }
    
    if auth.Password != "" {
        encryptedPassword, err := se.encryption.EncryptSensitiveField(auth.Password, key)
        if err != nil {
            return fmt.Errorf("failed to encrypt password: %w", err)
        }
        auth.Password = encryptedPassword
    }
    
    return nil
}

func (se *SelectiveEncryption) decryptSSHAuthentication(auth *spookytypes.SSHAuthentication, key []byte) error {
    if auth.PrivateKeyPath != "" {
        decryptedPath, err := se.encryption.DecryptSensitiveField(auth.PrivateKeyPath, key)
        if err != nil {
            return fmt.Errorf("failed to decrypt private key path: %w", err)
        }
        auth.PrivateKeyPath = decryptedPath
    }
    
    if auth.Passphrase != "" {
        decryptedPassphrase, err := se.encryption.DecryptSensitiveField(auth.Passphrase, key)
        if err != nil {
            return fmt.Errorf("failed to decrypt passphrase: %w", err)
        }
        auth.Passphrase = decryptedPassphrase
    }
    
    if auth.Password != "" {
        decryptedPassword, err := se.encryption.DecryptSensitiveField(auth.Password, key)
        if err != nil {
            return fmt.Errorf("failed to decrypt password: %w", err)
        }
        auth.Password = decryptedPassword
    }
    
    return nil
}
```

### 4. Key Management System

Implement a secure key management system:

```go
type KeyManager interface {
    GenerateKey() ([]byte, error)
    StoreKey(key []byte, identifier string) error
    RetrieveKey(identifier string) ([]byte, error)
    RotateKey(oldIdentifier string, newIdentifier string) error
    DeleteKey(identifier string) error
    ListKeys() ([]string, error)
    ValidateKey(key []byte) error
}

type SecureKeyManager struct {
    storage KeyStorage
    config  KeyManagerConfig
    logger  spookylogging.Logger
}

type KeyManagerConfig struct {
    KeySize        int           `hcl:"key_size"`
    KeyAlgorithm   string        `hcl:"key_algorithm"`
    KeyExpiration  time.Duration `hcl:"key_expiration"`
    AutoRotation   bool          `hcl:"auto_rotation"`
    RotationPeriod time.Duration `hcl:"rotation_period"`
}

type KeyStorage interface {
    Store(key []byte, identifier string, metadata map[string]string) error
    Retrieve(identifier string) ([]byte, map[string]string, error)
    Delete(identifier string) error
    List() ([]string, error)
    Exists(identifier string) (bool, error)
}

func (skm *SecureKeyManager) GenerateKey() ([]byte, error) {
    key := make([]byte, skm.config.KeySize/8)
    
    switch skm.config.KeyAlgorithm {
    case "AES":
        if _, err := rand.Read(key); err != nil {
            return nil, fmt.Errorf("failed to generate AES key: %w", err)
        }
    case "ChaCha20":
        if _, err := rand.Read(key); err != nil {
            return nil, fmt.Errorf("failed to generate ChaCha20 key: %w", err)
        }
    default:
        return nil, fmt.Errorf("unsupported key algorithm: %s", skm.config.KeyAlgorithm)
    }
    
    return key, nil
}

func (skm *SecureKeyManager) StoreKey(key []byte, identifier string) error {
    metadata := map[string]string{
        "algorithm":  skm.config.KeyAlgorithm,
        "key_size":   strconv.Itoa(skm.config.KeySize),
        "created_at": time.Now().Format(time.RFC3339),
        "expires_at": time.Now().Add(skm.config.KeyExpiration).Format(time.RFC3339),
    }
    
    return skm.storage.Store(key, identifier, metadata)
}

func (skm *SecureKeyManager) RetrieveKey(identifier string) ([]byte, error) {
    key, metadata, err := skm.storage.Retrieve(identifier)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve key: %w", err)
    }
    
    // Check if key has expired
    if expiresAt, ok := metadata["expires_at"]; ok {
        if expiration, err := time.Parse(time.RFC3339, expiresAt); err == nil {
            if time.Now().After(expiration) {
                return nil, fmt.Errorf("key %s has expired", identifier)
            }
        }
    }
    
    return key, nil
}
```

### 5. Configuration File Encryption

Implement encrypted configuration file handling:

```go
type EncryptedConfigManager struct {
    encryption ConfigurationEncryption
    keyManager KeyManager
    logger     spookylogging.Logger
    config     EncryptedConfigManagerConfig
}

type EncryptedConfigManagerConfig struct {
    DefaultKeyIdentifier string            `hcl:"default_key_identifier"`
    AutoEncrypt          bool              `hcl:"auto_encrypt"`
    EncryptionRequired   bool              `hcl:"encryption_required"`
    BackupUnencrypted    bool              `hcl:"backup_unencrypted"`
    SensitiveFields      []string          `hcl:"sensitive_fields"`
}

func (ecm *EncryptedConfigManager) LoadEncryptedConfig(filePath string) (*spookytypes.Config, error) {
    // Read encrypted configuration file
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }
    
    // Parse encrypted configuration
    var encryptedConfig EncryptedConfig
    if err := json.Unmarshal(data, &encryptedConfig); err != nil {
        return nil, fmt.Errorf("failed to parse encrypted config: %w", err)
    }
    
    // Retrieve encryption key
    key, err := ecm.keyManager.RetrieveKey(ecm.config.DefaultKeyIdentifier)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve encryption key: %w", err)
    }
    
    // Decrypt configuration
    config, err := ecm.encryption.DecryptConfig(&encryptedConfig, key)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt config: %w", err)
    }
    
    // Decrypt sensitive fields
    if err := ecm.decryptSensitiveFields(config, key); err != nil {
        return nil, fmt.Errorf("failed to decrypt sensitive fields: %w", err)
    }
    
    return config, nil
}

func (ecm *EncryptedConfigManager) SaveEncryptedConfig(config *spookytypes.Config, filePath string) error {
    // Retrieve encryption key
    key, err := ecm.keyManager.RetrieveKey(ecm.config.DefaultKeyIdentifier)
    if err != nil {
        return fmt.Errorf("failed to retrieve encryption key: %w", err)
    }
    
    // Encrypt sensitive fields
    if err := ecm.encryptSensitiveFields(config, key); err != nil {
        return fmt.Errorf("failed to encrypt sensitive fields: %w", err)
    }
    
    // Encrypt entire configuration
    encryptedConfig, err := ecm.encryption.EncryptConfig(config, key)
    if err != nil {
        return fmt.Errorf("failed to encrypt config: %w", err)
    }
    
    // Serialize encrypted configuration
    data, err := json.MarshalIndent(encryptedConfig, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to serialize encrypted config: %w", err)
    }
    
    // Write to file
    if err := os.WriteFile(filePath, data, 0600); err != nil {
        return fmt.Errorf("failed to write encrypted config: %w", err)
    }
    
    ecm.logger.Info("Configuration encrypted and saved", "file", filePath)
    return nil
}
```

### 6. CLI Integration

Integrate configuration encryption into CLI commands:

```go
type EncryptionCLI struct {
    configManager *EncryptedConfigManager
    keyManager    KeyManager
    logger        spookylogging.Logger
}

func (ecli *EncryptionCLI) EncryptConfigFile(inputPath, outputPath string) error {
    // Load unencrypted configuration
    config, err := spookyconfig.LoadConfig(inputPath)
    if err != nil {
        return fmt.Errorf("failed to load config: %w", err)
    }
    
    // Save as encrypted configuration
    if err := ecli.configManager.SaveEncryptedConfig(config, outputPath); err != nil {
        return fmt.Errorf("failed to save encrypted config: %w", err)
    }
    
    ecli.logger.Info("Configuration encrypted successfully",
        "input", inputPath,
        "output", outputPath)
    
    return nil
}

func (ecli *EncryptionCLI) DecryptConfigFile(inputPath, outputPath string) error {
    // Load encrypted configuration
    config, err := ecli.configManager.LoadEncryptedConfig(inputPath)
    if err != nil {
        return fmt.Errorf("failed to load encrypted config: %w", err)
    }
    
    // Save as unencrypted configuration
    if err := spookyconfig.SaveConfig(config, outputPath); err != nil {
        return fmt.Errorf("failed to save decrypted config: %w", err)
    }
    
    ecli.logger.Info("Configuration decrypted successfully",
        "input", inputPath,
        "output", outputPath)
    
    return nil
}

func (ecli *EncryptionCLI) GenerateEncryptionKey(identifier string) error {
    // Generate new encryption key
    key, err := ecli.keyManager.GenerateKey()
    if err != nil {
        return fmt.Errorf("failed to generate key: %w", err)
    }
    
    // Store key with identifier
    if err := ecli.keyManager.StoreKey(key, identifier); err != nil {
        return fmt.Errorf("failed to store key: %w", err)
    }
    
    ecli.logger.Info("Encryption key generated and stored",
        "identifier", identifier,
        "key_size", len(key)*8)
    
    return nil
}
```

## Implementation Plan

### Phase 1: Core Encryption
1. Implement AES-GCM encryption
2. Create key management system
3. Add selective field encryption
4. Implement configuration encryption interface

### Phase 2: Configuration Integration
1. Integrate encryption into configuration loading
2. Add encrypted configuration file support
3. Implement CLI commands for encryption
4. Add key rotation capabilities

### Phase 3: Advanced Features
1. Implement hardware security module (HSM) support
2. Add key escrow and recovery
3. Create encryption monitoring and auditing
4. Implement encryption policy management

## Benefits

- **Data Protection**: Sensitive configuration data encrypted at rest
- **Compliance**: Meets security requirements for sensitive data
- **Key Management**: Secure key generation, storage, and rotation
- **Selective Encryption**: Only encrypt sensitive fields for performance
- **Audit Trail**: Complete encryption/decryption logging

## Risks and Mitigation

### Risks
- Key loss leading to data inaccessibility
- Performance impact from encryption/decryption
- Complexity in key management
- Potential for encryption bypass

### Mitigation
- Robust key backup and recovery procedures
- Efficient encryption algorithms and selective encryption
- Clear key management documentation and procedures
- Regular security audits and penetration testing

## Success Metrics

- All sensitive configuration data encrypted
- Successful key rotation procedures
- Reduced security incidents related to configuration exposure
- Compliance with security requirements
- Minimal performance impact from encryption

## Related Documentation

- [Configuration Management](mdc:configuration-management)
- [Security Audit Logging](mdc:security-audit-logging)
- [SSH Implementation](mdc:ssh-implementation)
