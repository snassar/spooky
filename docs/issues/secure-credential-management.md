# Secure Credential Management

## Overview

The spooky codebase currently lacks a comprehensive secure credential management system for handling SSH keys, passwords, tokens, and other sensitive credentials. This document outlines recommendations for implementing secure credential management to improve security, reduce credential exposure, and provide better credential lifecycle management.

## Current State

- Credentials stored in plain text configuration files
- No secure credential storage mechanism
- Limited credential rotation capabilities
- No credential validation or health checking
- Insufficient credential access controls
- No audit trail for credential usage

## Recommendations

### 1. Credential Management Interface

Define a comprehensive credential management interface:

```go
type CredentialManager interface {
    StoreCredential(credential *Credential) error
    RetrieveCredential(identifier string) (*Credential, error)
    UpdateCredential(identifier string, credential *Credential) error
    DeleteCredential(identifier string) error
    ListCredentials(filters CredentialFilters) ([]*Credential, error)
    ValidateCredential(credential *Credential) error
    RotateCredential(identifier string) error
    GetCredentialHealth(identifier string) (*CredentialHealth, error)
    AuditCredentialAccess(identifier string) ([]CredentialAccess, error)
}

type Credential struct {
    ID          string                 `json:"id"`
    Type        CredentialType         `json:"type"`
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Data        map[string]interface{} `json:"data"`
    Metadata    map[string]string      `json:"metadata"`
    Tags        []string               `json:"tags"`
    ExpiresAt   *time.Time             `json:"expires_at,omitempty"`
    CreatedAt   time.Time              `json:"created_at"`
    UpdatedAt   time.Time              `json:"updated_at"`
    CreatedBy   string                 `json:"created_by"`
    LastUsed    *time.Time             `json:"last_used,omitempty"`
    UsageCount  int64                  `json:"usage_count"`
    Status      CredentialStatus       `json:"status"`
}

type CredentialType string

const (
    CredentialTypeSSHKey     CredentialType = "ssh_key"
    CredentialTypePassword   CredentialType = "password"
    CredentialTypeToken      CredentialType = "token"
    CredentialTypeCertificate CredentialType = "certificate"
    CredentialTypeAPIKey     CredentialType = "api_key"
    CredentialTypeSecret     CredentialType = "secret"
)

type CredentialStatus string

const (
    CredentialStatusActive   CredentialStatus = "active"
    CredentialStatusInactive CredentialStatus = "inactive"
    CredentialStatusExpired  CredentialStatus = "expired"
    CredentialStatusRevoked  CredentialStatus = "revoked"
)

type CredentialHealth struct {
    ID              string            `json:"id"`
    Status          CredentialStatus  `json:"status"`
    LastValidation  time.Time         `json:"last_validation"`
    ValidationResult ValidationResult `json:"validation_result"`
    ExpirationDate  *time.Time        `json:"expiration_date"`
    DaysUntilExpiry int               `json:"days_until_expiry"`
    UsagePattern    UsagePattern      `json:"usage_pattern"`
    SecurityScore   int               `json:"security_score"`
}

type ValidationResult struct {
    Valid    bool     `json:"valid"`
    Errors   []string `json:"errors"`
    Warnings []string `json:"warnings"`
}
```

### 2. Secure Credential Storage

Implement secure credential storage with encryption:

```go
type SecureCredentialStorage struct {
    encryption CredentialEncryption
    storage    CredentialStorage
    config     CredentialStorageConfig
    logger     spookylogging.Logger
}

type CredentialStorageConfig struct {
    EncryptionEnabled bool              `hcl:"encryption_enabled"`
    EncryptionKey     string            `hcl:"encryption_key"`
    StorageBackend    string            `hcl:"storage_backend"`
    BackupEnabled     bool              `hcl:"backup_enabled"`
    BackupLocation    string            `hcl:"backup_location"`
    RetentionPolicy   RetentionPolicy   `hcl:"retention_policy"`
}

type RetentionPolicy struct {
    MaxAge           time.Duration `hcl:"max_age"`
    MaxUsageCount    int64         `hcl:"max_usage_count"`
    AutoExpiration   bool          `hcl:"auto_expiration"`
    ArchiveExpired   bool          `hcl:"archive_expired"`
}

func (scs *SecureCredentialStorage) StoreCredential(credential *Credential) error {
    // Validate credential before storage
    if err := scs.validateCredential(credential); err != nil {
        return fmt.Errorf("credential validation failed: %w", err)
    }
    
    // Encrypt sensitive credential data
    encryptedCredential, err := scs.encryptCredential(credential)
    if err != nil {
        return fmt.Errorf("failed to encrypt credential: %w", err)
    }
    
    // Store encrypted credential
    if err := scs.storage.Store(encryptedCredential); err != nil {
        return fmt.Errorf("failed to store credential: %w", err)
    }
    
    // Create backup if enabled
    if scs.config.BackupEnabled {
        if err := scs.createBackup(encryptedCredential); err != nil {
            scs.logger.Warn("Failed to create credential backup", "error", err)
        }
    }
    
    scs.logger.Info("Credential stored successfully", "id", credential.ID, "type", credential.Type)
    return nil
}

func (scs *SecureCredentialStorage) RetrieveCredential(identifier string) (*Credential, error) {
    // Retrieve encrypted credential
    encryptedCredential, err := scs.storage.Retrieve(identifier)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve credential: %w", err)
    }
    
    // Decrypt credential data
    credential, err := scs.decryptCredential(encryptedCredential)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt credential: %w", err)
    }
    
    // Update usage statistics
    credential.LastUsed = timePtr(time.Now())
    credential.UsageCount++
    
    // Store updated credential
    if err := scs.StoreCredential(credential); err != nil {
        scs.logger.Warn("Failed to update credential usage statistics", "error", err)
    }
    
    // Audit credential access
    scs.auditCredentialAccess(credential.ID, "retrieve")
    
    return credential, nil
}

func (scs *SecureCredentialStorage) encryptCredential(credential *Credential) (*EncryptedCredential, error) {
    // Serialize credential data
    data, err := json.Marshal(credential)
    if err != nil {
        return nil, fmt.Errorf("failed to serialize credential: %w", err)
    }
    
    // Encrypt credential data
    encryptedData, err := scs.encryption.Encrypt(data, scs.config.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt credential data: %w", err)
    }
    
    return &EncryptedCredential{
        ID:            credential.ID,
        EncryptedData: encryptedData,
        CreatedAt:     time.Now(),
        Metadata:      credential.Metadata,
    }, nil
}
```

### 3. SSH Key Management

Implement specialized SSH key management:

```go
type SSHKeyManager struct {
    credentialManager CredentialManager
    keyValidator     SSHKeyValidator
    keyGenerator     SSHKeyGenerator
    logger           spookylogging.Logger
}

type SSHKeyValidator interface {
    ValidatePrivateKey(keyData []byte, passphrase string) error
    ValidatePublicKey(keyData []byte) error
    CheckKeyPermissions(keyPath string) error
    ValidateKeyPair(privateKey []byte, publicKey []byte) error
}

type SSHKeyGenerator interface {
    GenerateKeyPair(keyType string, bits int, passphrase string) (*SSHKeyPair, error)
    GenerateKeyPairWithComment(keyType string, bits int, passphrase string, comment string) (*SSHKeyPair, error)
}

type SSHKeyPair struct {
    PrivateKey []byte `json:"private_key"`
    PublicKey  []byte `json:"public_key"`
    KeyType    string `json:"key_type"`
    Bits       int    `json:"bits"`
    Comment    string `json:"comment"`
}

func (skm *SSHKeyManager) StoreSSHKey(keyPath string, passphrase string, metadata map[string]string) error {
    // Read private key file
    privateKeyData, err := os.ReadFile(keyPath)
    if err != nil {
        return fmt.Errorf("failed to read private key file: %w", err)
    }
    
    // Validate private key
    if err := skm.keyValidator.ValidatePrivateKey(privateKeyData, passphrase); err != nil {
        return fmt.Errorf("invalid private key: %w", err)
    }
    
    // Read public key file
    publicKeyPath := keyPath + ".pub"
    publicKeyData, err := os.ReadFile(publicKeyPath)
    if err != nil {
        return fmt.Errorf("failed to read public key file: %w", err)
    }
    
    // Validate public key
    if err := skm.keyValidator.ValidatePublicKey(publicKeyData); err != nil {
        return fmt.Errorf("invalid public key: %w", err)
    }
    
    // Validate key pair
    if err := skm.keyValidator.ValidateKeyPair(privateKeyData, publicKeyData); err != nil {
        return fmt.Errorf("key pair validation failed: %w", err)
    }
    
    // Create credential
    credential := &Credential{
        ID:          generateCredentialID(),
        Type:        CredentialTypeSSHKey,
        Name:        filepath.Base(keyPath),
        Description: "SSH key pair",
        Data: map[string]interface{}{
            "private_key": string(privateKeyData),
            "public_key":  string(publicKeyData),
            "key_path":    keyPath,
        },
        Metadata: metadata,
        Tags:      []string{"ssh", "key"},
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        CreatedBy: getCurrentUser(),
        Status:    CredentialStatusActive,
    }
    
    // Store credential
    return skm.credentialManager.StoreCredential(credential)
}

func (skm *SSHKeyManager) GenerateAndStoreSSHKey(keyType string, bits int, passphrase string, name string) error {
    // Generate new SSH key pair
    keyPair, err := skm.keyGenerator.GenerateKeyPair(keyType, bits, passphrase)
    if err != nil {
        return fmt.Errorf("failed to generate SSH key pair: %w", err)
    }
    
    // Create credential
    credential := &Credential{
        ID:          generateCredentialID(),
        Type:        CredentialTypeSSHKey,
        Name:        name,
        Description: fmt.Sprintf("Generated SSH %s key (%d bits)", keyType, bits),
        Data: map[string]interface{}{
            "private_key": string(keyPair.PrivateKey),
            "public_key":  string(keyPair.PublicKey),
            "key_type":    keyPair.KeyType,
            "bits":        keyPair.Bits,
            "comment":     keyPair.Comment,
        },
        Tags:      []string{"ssh", "key", "generated"},
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        CreatedBy: getCurrentUser(),
        Status:    CredentialStatusActive,
    }
    
    // Store credential
    return skm.credentialManager.StoreCredential(credential)
}
```

### 4. Credential Validation and Health Checking

Implement credential validation and health monitoring:

```go
type CredentialValidator struct {
    validators map[CredentialType]CredentialTypeValidator
    logger     spookylogging.Logger
}

type CredentialTypeValidator interface {
    Validate(credential *Credential) (ValidationResult, error)
    CheckHealth(credential *Credential) (*CredentialHealth, error)
    TestCredential(credential *Credential) error
}

type SSHKeyValidator struct {
    logger spookylogging.Logger
}

func (skv *SSHKeyValidator) Validate(credential *Credential) (ValidationResult, error) {
    result := ValidationResult{Valid: true}
    
    // Check required fields
    if credential.Data["private_key"] == nil {
        result.Valid = false
        result.Errors = append(result.Errors, "private key is required")
    }
    
    if credential.Data["public_key"] == nil {
        result.Valid = false
        result.Errors = append(result.Errors, "public key is required")
    }
    
    // Validate private key format
    if privateKey, ok := credential.Data["private_key"].(string); ok {
        if err := skv.validatePrivateKeyFormat([]byte(privateKey)); err != nil {
            result.Valid = false
            result.Errors = append(result.Errors, fmt.Sprintf("invalid private key format: %v", err))
        }
    }
    
    // Validate public key format
    if publicKey, ok := credential.Data["public_key"].(string); ok {
        if err := skv.validatePublicKeyFormat([]byte(publicKey)); err != nil {
            result.Valid = false
            result.Errors = append(result.Errors, fmt.Sprintf("invalid public key format: %v", err))
        }
    }
    
    // Check expiration
    if credential.ExpiresAt != nil && time.Now().After(*credential.ExpiresAt) {
        result.Valid = false
        result.Errors = append(result.Errors, "credential has expired")
    }
    
    return result, nil
}

func (skv *SSHKeyValidator) CheckHealth(credential *Credential) (*CredentialHealth, error) {
    health := &CredentialHealth{
        ID:             credential.ID,
        Status:         credential.Status,
        LastValidation: time.Now(),
    }
    
    // Validate credential
    validationResult, err := skv.Validate(credential)
    if err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    health.ValidationResult = validationResult
    
    // Calculate days until expiry
    if credential.ExpiresAt != nil {
        health.ExpirationDate = credential.ExpiresAt
        health.DaysUntilExpiry = int(credential.ExpiresAt.Sub(time.Now()).Hours() / 24)
        
        if health.DaysUntilExpiry < 0 {
            health.Status = CredentialStatusExpired
        } else if health.DaysUntilExpiry < 30 {
            health.ValidationResult.Warnings = append(health.ValidationResult.Warnings, 
                fmt.Sprintf("credential expires in %d days", health.DaysUntilExpiry))
        }
    }
    
    // Calculate security score
    health.SecurityScore = skv.calculateSecurityScore(credential)
    
    return health, nil
}

func (skv *SSHKeyValidator) calculateSecurityScore(credential *Credential) int {
    score := 100
    
    // Deduct points for various security issues
    if credential.UsageCount > 1000 {
        score -= 20 // High usage might indicate overuse
    }
    
    if credential.LastUsed != nil && time.Since(*credential.LastUsed) > 365*24*time.Hour {
        score -= 30 // Unused for over a year
    }
    
    if credential.ExpiresAt != nil && time.Until(*credential.ExpiresAt) < 30*24*time.Hour {
        score -= 40 // Expires soon
    }
    
    // Check key strength
    if bits, ok := credential.Data["bits"].(int); ok {
        if bits < 2048 {
            score -= 30 // Weak key
        } else if bits >= 4096 {
            score += 10 // Strong key
        }
    }
    
    return max(0, score)
}
```

### 5. Credential Rotation and Lifecycle Management

Implement credential rotation and lifecycle management:

```go
type CredentialLifecycleManager struct {
    credentialManager CredentialManager
    rotationPolicy    RotationPolicy
    logger            spookylogging.Logger
}

type RotationPolicy struct {
    AutoRotation     bool          `hcl:"auto_rotation"`
    RotationInterval time.Duration `hcl:"rotation_interval"`
    WarningPeriod    time.Duration `hcl:"warning_period"`
    GracePeriod      time.Duration `hcl:"grace_period"`
    MaxAge           time.Duration `hcl:"max_age"`
    RotationStrategy string        `hcl:"rotation_strategy"`
}

func (clm *CredentialLifecycleManager) StartRotationMonitoring() {
    ticker := time.NewTicker(24 * time.Hour) // Check daily
    defer ticker.Stop()
    
    for range ticker.C {
        if err := clm.checkAndRotateCredentials(); err != nil {
            clm.logger.Error("Failed to check and rotate credentials", "error", err)
        }
    }
}

func (clm *CredentialLifecycleManager) checkAndRotateCredentials() error {
    // Get all active credentials
    filters := CredentialFilters{
        Status: []CredentialStatus{CredentialStatusActive},
    }
    
    credentials, err := clm.credentialManager.ListCredentials(filters)
    if err != nil {
        return fmt.Errorf("failed to list credentials: %w", err)
    }
    
    for _, credential := range credentials {
        if err := clm.checkCredentialRotation(credential); err != nil {
            clm.logger.Error("Failed to check credential rotation", 
                "credential_id", credential.ID, "error", err)
        }
    }
    
    return nil
}

func (clm *CredentialLifecycleManager) checkCredentialRotation(credential *Credential) error {
    // Check if credential needs rotation
    if clm.needsRotation(credential) {
        clm.logger.Info("Credential needs rotation", "credential_id", credential.ID)
        
        if clm.rotationPolicy.AutoRotation {
            if err := clm.rotateCredential(credential); err != nil {
                return fmt.Errorf("failed to rotate credential: %w", err)
            }
        } else {
            // Send warning notification
            clm.sendRotationWarning(credential)
        }
    }
    
    return nil
}

func (clm *CredentialLifecycleManager) needsRotation(credential *Credential) bool {
    // Check age-based rotation
    if clm.rotationPolicy.MaxAge > 0 {
        age := time.Since(credential.CreatedAt)
        if age > clm.rotationPolicy.MaxAge {
            return true
        }
    }
    
    // Check usage-based rotation
    if credential.UsageCount > 10000 {
        return true
    }
    
    // Check expiration-based rotation
    if credential.ExpiresAt != nil {
        timeUntilExpiry := time.Until(*credential.ExpiresAt)
        if timeUntilExpiry < clm.rotationPolicy.WarningPeriod {
            return true
        }
    }
    
    return false
}

func (clm *CredentialLifecycleManager) rotateCredential(credential *Credential) error {
    clm.logger.Info("Starting credential rotation", "credential_id", credential.ID)
    
    // Create new credential
    newCredential, err := clm.createNewCredential(credential)
    if err != nil {
        return fmt.Errorf("failed to create new credential: %w", err)
    }
    
    // Store new credential
    if err := clm.credentialManager.StoreCredential(newCredential); err != nil {
        return fmt.Errorf("failed to store new credential: %w", err)
    }
    
    // Update old credential status
    credential.Status = CredentialStatusInactive
    credential.UpdatedAt = time.Now()
    if err := clm.credentialManager.UpdateCredential(credential.ID, credential); err != nil {
        return fmt.Errorf("failed to update old credential: %w", err)
    }
    
    // Notify about rotation
    clm.notifyCredentialRotation(credential, newCredential)
    
    clm.logger.Info("Credential rotation completed", 
        "old_credential_id", credential.ID,
        "new_credential_id", newCredential.ID)
    
    return nil
}
```

### 6. Credential Access Control and Auditing

Implement credential access control and auditing:

```go
type CredentialAccessControl struct {
    credentialManager CredentialManager
    accessPolicy      AccessPolicy
    auditLogger       AuditLogger
    logger            spookylogging.Logger
}

type AccessPolicy struct {
    RequireAuthentication bool     `hcl:"require_authentication"`
    AllowedUsers         []string `hcl:"allowed_users"`
    AllowedRoles         []string `hcl:"allowed_roles"`
    DeniedUsers          []string `hcl:"denied_users"`
    DeniedRoles          []string `hcl:"denied_roles"`
    MaxConcurrentAccess  int      `hcl:"max_concurrent_access"`
    AccessTimeout        time.Duration `hcl:"access_timeout"`
}

type CredentialAccess struct {
    ID            string    `json:"id"`
    CredentialID  string    `json:"credential_id"`
    User          string    `json:"user"`
    Action        string    `json:"action"`
    Timestamp     time.Time `json:"timestamp"`
    IPAddress     string    `json:"ip_address"`
    UserAgent     string    `json:"user_agent"`
    Success       bool      `json:"success"`
    ErrorMessage  string    `json:"error_message,omitempty"`
    SessionID     string    `json:"session_id,omitempty"`
}

func (cac *CredentialAccessControl) CheckAccess(credentialID string, user string, action string) (bool, error) {
    // Check if user is explicitly denied
    if cac.isUserDenied(user) {
        cac.auditAccess(credentialID, user, action, false, "user explicitly denied")
        return false, fmt.Errorf("access denied for user: %s", user)
    }
    
    // Check if user is explicitly allowed
    if !cac.isUserAllowed(user) {
        cac.auditAccess(credentialID, user, action, false, "user not in allowed list")
        return false, fmt.Errorf("access denied for user: %s", user)
    }
    
    // Check concurrent access limits
    if err := cac.checkConcurrentAccess(credentialID); err != nil {
        cac.auditAccess(credentialID, user, action, false, err.Error())
        return false, fmt.Errorf("concurrent access limit exceeded: %w", err)
    }
    
    // Audit successful access
    cac.auditAccess(credentialID, user, action, true, "")
    
    return true, nil
}

func (cac *CredentialAccessControl) auditAccess(credentialID, user, action string, success bool, errorMessage string) {
    access := CredentialAccess{
        ID:           generateAccessID(),
        CredentialID: credentialID,
        User:         user,
        Action:       action,
        Timestamp:    time.Now(),
        IPAddress:    getClientIP(),
        UserAgent:    getUserAgent(),
        Success:      success,
        ErrorMessage: errorMessage,
        SessionID:    getSessionID(),
    }
    
    if err := cac.auditLogger.LogCredentialAccess(access); err != nil {
        cac.logger.Error("Failed to log credential access", "error", err)
    }
}
```

## Implementation Plan

### Phase 1: Core Credential Management
1. Implement credential storage and encryption
2. Create credential validation system
3. Add SSH key management
4. Implement basic access controls

### Phase 2: Advanced Features
1. Add credential rotation and lifecycle management
2. Implement credential health monitoring
3. Create comprehensive auditing
4. Add credential backup and recovery

### Phase 3: Integration and Automation
1. Integrate with SSH connection management
2. Add automated credential rotation
3. Implement credential monitoring dashboards
4. Create credential policy management

## Benefits

- **Enhanced Security**: Secure storage and handling of sensitive credentials
- **Credential Lifecycle Management**: Automated rotation and expiration handling
- **Access Control**: Granular control over credential access
- **Audit Trail**: Complete tracking of credential usage
- **Compliance**: Meets security and compliance requirements

## Risks and Mitigation

### Risks
- Credential loss leading to service disruption
- Performance impact from encryption/decryption
- Complexity in credential management
- Potential for credential misuse

### Mitigation
- Robust backup and recovery procedures
- Efficient credential caching and storage
- Clear credential management policies
- Regular security audits and monitoring

## Success Metrics

- All sensitive credentials encrypted and securely stored
- Successful automated credential rotation
- Reduced credential-related security incidents
- Complete audit trail for credential access
- Improved compliance with security requirements

## Related Documentation

- [Configuration Encryption](mdc:configuration-encryption)
- [Security Audit Logging](mdc:security-audit-logging)
- [SSH Implementation](mdc:ssh-implementation)
