# Secrets Management Implementation Plan

## Overview

This document outlines the implementation plan for integrating [age encryption](https://github.com/FiloSottile/age) into spooky for comprehensive secrets management. The plan covers migrating from the current AES-GCM implementation to age-based encryption for variables, facts, and machine inventories.

## Current State Analysis

### Existing Implementation
- **Location**: `internal/secrets/integration.go`
- **Current Method**: AES-GCM encryption with 32-byte keys
- **Interface**: `SecretsIntegration` with `Encrypt`, `Decrypt`, `ValidateKey` methods
- **Configuration**: Basic age configuration exists in `spooky.schema.hcl`

### Current Limitations
1. Uses AES-GCM instead of age encryption
2. Limited to symmetric encryption only
3. No support for multiple recipients
4. No integration with age key management
5. No support for SSH key-based encryption

## Implementation Goals

### Primary Objectives
1. **Replace AES-GCM with age encryption** for all secrets management
2. **Support multiple encryption scenarios**:
   - Encrypted variables in `variables.hcl` (encrypt/decrypt)
   - Encrypted facts in `/etc/spooky/facts.*` on target machines (decrypt only)
   - Encrypted machine inventory secrets (passphrases, keys)
3. **Maintain backward compatibility** during transition
4. **Provide comprehensive CLI support** for age key management
5. **Integrate with existing spooky architecture**

### Secondary Objectives
1. **Support multiple recipients** for encrypted data
2. **Enable SSH key-based encryption** for convenience
3. **Provide audit logging** for encryption/decryption operations
4. **Support passphrase-based encryption** as fallback
5. **Integrate with existing validation systems**

## Technical Architecture

### Age Integration Design

#### Core Components
```go
// Enhanced SecretsIntegration interface
type SecretsIntegration interface {
    // Age-specific methods
    EncryptWithAge(ctx context.Context, data []byte, recipients []string) ([]byte, error)
    DecryptWithAge(ctx context.Context, data []byte, identityPath string) ([]byte, error)
    EncryptWithPassphrase(ctx context.Context, data []byte, passphrase string) ([]byte, error)
    DecryptWithPassphrase(ctx context.Context, data []byte, passphrase string) ([]byte, error)
    
    // Key management
    ValidateAgeKey(ctx context.Context, keyPath string) error
    ListRecipients(ctx context.Context, encryptedData []byte) ([]string, error)
    
    // Legacy support (deprecated)
    Encrypt(ctx context.Context, data []byte, key []byte) ([]byte, error)
    Decrypt(ctx context.Context, data []byte, key []byte) ([]byte, error)
    ValidateKey(ctx context.Context, key []byte) error
}
```

#### Configuration Structure
```hcl
# Enhanced age configuration in spooky.hcl
age {
  enabled = true
  
  # Primary age key configuration
  public_key = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
  private_key_path = "~/.config/spooky/age.key"
  
  # Multiple recipients support
  recipients = [
    "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p",
    "age1lggyhqrw2nlhcxprm67z43rta597azn8gknawjehu9d9dl0jq3yqqvfafg"
  ]
  
  # SSH key support
  ssh_keys = [
    "~/.ssh/id_ed25519.pub",
    "~/.ssh/id_rsa.pub"
  ]
  
  # Fallback passphrase configuration
  passphrase_fallback = true
  passphrase_prompt = "Enter encryption passphrase:"
  
  # Security settings
  audit_logging = true
  key_validation = true
  recipient_validation = true
}
```

### Use Case Implementations

#### 1. Encrypted Variables (Encrypt/Decrypt)
```hcl
# variables.hcl with encrypted values
# These can be encrypted by spooky or pre-encrypted by users
variables {
  database_password = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
  api_key = "age1lggyhqrw2nlhcxprm67z43rta597azn8gknawjehu9d9dl0jq3yqqvfafg"
  
  # Encrypted with multiple recipients
  shared_secret = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
  
  # Passphrase encrypted
  backup_key = "age1passphrase1..."
}
```

#### 2. Encrypted Facts (Decrypt Only)
```hcl
# /etc/spooky/facts/custom.hcl on target machines
# Facts are pre-encrypted by machine administrators
# Spooky only decrypts them during fact collection (facts are not stored on disk)
facts {
  # Encrypted sensitive system information (pre-encrypted by machine admin)
  database_credentials = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
  api_tokens = "age1lggyhqrw2nlhcxprm67z43rta597azn8gknawjehu9d9dl0jq3yqqvfafg"
  
  # Encrypted with SSH key for convenience (pre-encrypted by machine admin)
  local_secrets = "age1ssh-ed25519..."
}
```

#### 3. Actions Run with Decryption
```bash
# Run actions with decryption enabled for secret material
spooky actions run my-project --decrypt

# This allows the orchestrator to use encrypted variables and facts
# but ensures secrets are never logged or displayed in plaintext
```

#### 3. Encrypted Machine Inventory
```hcl
# machines.hcl with encrypted authentication
machines {
  machine "web-server" {
    hostname = "web.example.com"
    port = 22
    user = "admin"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/id_rsa"
      
      # Encrypted passphrase for SSH key
      key_passphrase = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
    }
    
    # Encrypted sudo password
    sudo_password = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
  }
}
```

## Implementation Phases

### Phase 1: Core Age Integration (Week 1-2)

#### 1.1 Dependencies and Setup
- [ ] Add `filippo.io/age` dependency to `go.mod`
- [ ] Create age-specific types in `internal/types/secrets/`
- [ ] Update `internal/types/types.go` to re-export new types
- [ ] Create age configuration types

#### 1.2 Core Implementation
- [ ] Implement age encryption/decryption in `internal/secrets/integration.go`
- [ ] Add age key validation and management
- [ ] Implement recipient list extraction
- [ ] Add passphrase encryption support
- [ ] Create comprehensive tests

#### 1.3 Configuration Integration
- [ ] Update `internal/schemas/schemas/spooky.schema.hcl` with enhanced age config
- [ ] Create age configuration validation
- [ ] Add age configuration loading in `internal/config/`
- [ ] Update configuration examples

### Phase 2: Variable Encryption/Decryption (Week 3)

#### 2.1 Variable Schema Updates
- [ ] Update `internal/schemas/schemas/variables-structure.schema.hcl`
- [ ] Add age encryption/decryption support to variable types
- [ ] Create variable encryption/decryption logic
- [ ] Update variable validation to handle encrypted values

#### 2.2 Variable Integration
- [ ] Modify `internal/variables/` to support age encryption/decryption
- [ ] Add automatic encryption/decryption during variable resolution
- [ ] Create variable encryption/decryption CLI commands
- [ ] Add variable encryption/decryption examples

### Phase 3: Facts Decryption (Week 4)

#### 3.1 Facts Schema Updates
- [ ] Update `internal/schemas/schemas/custom-facts-hcl.schema.hcl`
- [ ] Add age decryption support to fact types
- [ ] Create fact decryption logic (read-only)
- [ ] Update fact validation to handle encrypted values

#### 3.2 Facts Integration
- [ ] Modify `internal/facts/` to support age decryption
- [ ] Add automatic decryption during fact collection from `/etc/spooky/facts.*`
- [ ] Facts are not stored on disk - only decrypted in memory during collection
- [ ] Add fact decryption examples and documentation

### Phase 4: Actions Run with Decryption (Week 5)

#### 4.1 Actions Run Integration
- [ ] Add `--decrypt` flag to `spooky actions run` command
- [ ] Modify action orchestration to use decrypted secret material
- [ ] Implement secure secret handling during action execution
- [ ] Ensure secrets are never logged or displayed in plaintext

#### 4.2 Security and Logging
- [ ] Implement secret masking in all logging output
- [ ] Add audit logging for decryption operations
- [ ] Create secure secret material handling patterns
- [ ] Add validation to prevent secret exposure in output

### Phase 5: CLI and User Experience (Week 6)

#### 5.1 CLI Commands
- [ ] Add `spooky variables encrypt <variable>` command for encrypting variable values
- [ ] Ensure secrets are never emitted in plaintext in logs or console output
- [ ] Follow spooky's existing CLI patterns and command structure

#### 5.2 Documentation and Examples
- [ ] Create comprehensive secrets management documentation
- [ ] Add age encryption examples to `docs/examples/`
- [ ] Create troubleshooting guide for secrets issues
- [ ] Add migration guide from AES to age

### Phase 6: Testing and Validation (Week 7)

#### 6.1 Comprehensive Testing
- [ ] Unit tests for all age encryption functions
- [ ] Integration tests with real age keys
- [ ] Performance tests for encryption/decryption
- [ ] Security tests for key management

#### 6.2 Validation and Audit
- [ ] Security review of age integration
- [ ] Performance validation
- [ ] User acceptance testing
- [ ] Documentation review

## Detailed Implementation

### 1. Age Types Definition

```go
// internal/types/secrets/types.go
package spookytypessecrets

import (
    "time"
    spookytypescommon "spooky/internal/types/common"
)

// AgeConfig represents age encryption configuration
type AgeConfig struct {
    Enabled              bool     `json:"enabled" hcl:"enabled"`
    PublicKey            string   `json:"public_key,omitempty" hcl:"public_key,optional"`
    PrivateKeyPath       string   `json:"private_key_path,omitempty" hcl:"private_key_path,optional"`
    Recipients           []string `json:"recipients,omitempty" hcl:"recipients,optional"`
    SSHKeys              []string `json:"ssh_keys,omitempty" hcl:"ssh_keys,optional"`
    PassphraseFallback   bool     `json:"passphrase_fallback,omitempty" hcl:"passphrase_fallback,optional"`
    PassphrasePrompt     string   `json:"passphrase_prompt,omitempty" hcl:"passphrase_prompt,optional"`
    AuditLogging         bool     `json:"audit_logging,omitempty" hcl:"audit_logging,optional"`
    KeyValidation        bool     `json:"key_validation,omitempty" hcl:"key_validation,optional"`
    RecipientValidation  bool     `json:"recipient_validation,omitempty" hcl:"recipient_validation,optional"`
}

// AgeKey represents an age key pair
type AgeKey struct {
    spookytypescommon.CompleteEntity
    
    PublicKey  string    `json:"public_key" hcl:"public_key"`
    PrivateKey string    `json:"private_key,omitempty" hcl:"private_key,optional" sensitive:"true"`
    KeyPath    string    `json:"key_path,omitempty" hcl:"key_path,optional"`
    CreatedAt  time.Time `json:"created_at" hcl:"created_at"`
    ExpiresAt  *time.Time `json:"expires_at,omitempty" hcl:"expires_at,optional"`
    IsValid    bool      `json:"is_valid" hcl:"is_valid"`
}

// AgeRecipient represents an age recipient
type AgeRecipient struct {
    spookytypescommon.NamedEntity
    
    PublicKey string `json:"public_key" hcl:"public_key"`
    Type      string `json:"type" hcl:"type"` // "age", "ssh"
    Source    string `json:"source,omitempty" hcl:"source,optional"`
    IsValid   bool   `json:"is_valid" hcl:"is_valid"`
}

// AgeEncryptedData represents encrypted data with metadata
type AgeEncryptedData struct {
    Data       []byte                 `json:"data" hcl:"data"`
    Recipients []string               `json:"recipients" hcl:"recipients"`
    Metadata   map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
    CreatedAt  time.Time              `json:"created_at" hcl:"created_at"`
    Version    string                 `json:"version" hcl:"version"`
}
```

### 2. Enhanced Secrets Integration

```go
// internal/secrets/integration.go (enhanced)
package secrets

import (
    "context"
    "fmt"
    "os"
    "strings"
    
    "filippo.io/age"
    "filippo.io/age/agessh"
    
    spookyinterfaces "spooky/internal/interfaces"
    spookytypessecrets "spooky/internal/types/secrets"
    spookytypeslogging "spooky/internal/types/logging"
)

// Integration implements the enhanced SecretsIntegration interface
type Integration struct {
    logger spookytypeslogging.Logger
    config *spookytypessecrets.AgeConfig
}

// NewIntegration creates a new secrets integration
func NewIntegration(logger spookytypeslogging.Logger, config *spookytypessecrets.AgeConfig) spookyinterfaces.SecretsIntegration {
    return &Integration{
        logger: logger,
        config: config,
    }
}

// EncryptWithAge encrypts data with age encryption
func (i *Integration) EncryptWithAge(ctx context.Context, data []byte, recipients []string) ([]byte, error) {
    if len(data) == 0 {
        return nil, fmt.Errorf("data cannot be empty")
    }
    
    if len(recipients) == 0 {
        return nil, fmt.Errorf("at least one recipient is required")
    }
    
    // Parse recipients
    var ageRecipients []age.Recipient
    for _, recipient := range recipients {
        if strings.HasPrefix(recipient, "age1") {
            // Age public key
            r, err := age.ParseX25519Recipient(recipient)
            if err != nil {
                return nil, fmt.Errorf("invalid age recipient %s: %w", recipient, err)
            }
            ageRecipients = append(ageRecipients, r)
        } else if strings.HasPrefix(recipient, "ssh-") {
            // SSH public key
            r, err := agessh.ParseRecipient(recipient)
            if err != nil {
                return nil, fmt.Errorf("invalid SSH recipient %s: %w", recipient, err)
            }
            ageRecipients = append(ageRecipients, r)
        } else {
            return nil, fmt.Errorf("unsupported recipient format: %s", recipient)
        }
    }
    
    // Encrypt data
    encrypted, err := age.Encrypt(ageRecipients, strings.NewReader(string(data)))
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt data: %w", err)
    }
    
    // Read encrypted data
    encryptedData, err := io.ReadAll(encrypted)
    if err != nil {
        return nil, fmt.Errorf("failed to read encrypted data: %w", err)
    }
    
    i.logger.Info("Data encrypted with age successfully", map[string]interface{}{
        "data_size":       len(data),
        "ciphertext_size": len(encryptedData),
        "recipients":      recipients,
    })
    
    return encryptedData, nil
}

// DecryptWithAge decrypts data with age encryption
func (i *Integration) DecryptWithAge(ctx context.Context, data []byte, identityPath string) ([]byte, error) {
    if len(data) == 0 {
        return nil, fmt.Errorf("data cannot be empty")
    }
    
    // Load identity
    identity, err := i.loadIdentity(identityPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load identity: %w", err)
    }
    
    // Decrypt data
    decrypted, err := age.Decrypt(identity, bytes.NewReader(data))
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt data: %w", err)
    }
    
    // Read decrypted data
    decryptedData, err := io.ReadAll(decrypted)
    if err != nil {
        return nil, fmt.Errorf("failed to read decrypted data: %w", err)
    }
    
    i.logger.Info("Data decrypted with age successfully", map[string]interface{}{
        "ciphertext_size": len(data),
        "plaintext_size":  len(decryptedData),
        "identity_path":   identityPath,
    })
    
    return decryptedData, nil
}

// loadIdentity loads an age identity from file or passphrase
func (i *Integration) loadIdentity(identityPath string) (age.Identity, error) {
    if strings.HasPrefix(identityPath, "age1passphrase1") {
        // Passphrase-based identity
        passphrase, err := i.promptPassphrase()
        if err != nil {
            return nil, fmt.Errorf("failed to get passphrase: %w", err)
        }
        return age.NewScryptIdentity(passphrase)
    }
    
    // File-based identity
    f, err := os.Open(identityPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open identity file: %w", err)
    }
    defer f.Close()
    
    identities, err := age.ParseIdentities(f)
    if err != nil {
        return nil, fmt.Errorf("failed to parse identity file: %w", err)
    }
    
    if len(identities) == 0 {
        return nil, fmt.Errorf("no identities found in file")
    }
    
    return identities[0], nil
}

// promptPassphrase prompts for a passphrase
func (i *Integration) promptPassphrase() (string, error) {
    // Implementation would use terminal input
    // For now, return error indicating passphrase input not implemented
    return "", fmt.Errorf("passphrase input not implemented")
}

// ListRecipients lists recipients from encrypted data
func (i *Integration) ListRecipients(ctx context.Context, encryptedData []byte) ([]string, error) {
    // Parse encrypted data to extract recipient information
    // This is a simplified implementation
    return []string{}, fmt.Errorf("recipient listing not implemented")
}

// ValidateAgeKey validates an age key
func (i *Integration) ValidateAgeKey(ctx context.Context, keyPath string) error {
    f, err := os.Open(keyPath)
    if err != nil {
        return fmt.Errorf("failed to open key file: %w", err)
    }
    defer f.Close()
    
    identities, err := age.ParseIdentities(f)
    if err != nil {
        return fmt.Errorf("failed to parse key file: %w", err)
    }
    
    if len(identities) == 0 {
        return fmt.Errorf("no valid identities found in key file")
    }
    
    i.logger.Info("Age key validated successfully", map[string]interface{}{
        "key_path": keyPath,
    })
    
    return nil
}

// Legacy methods (deprecated)
func (i *Integration) Encrypt(ctx context.Context, data []byte, key []byte) ([]byte, error) {
    i.logger.Warn("Using deprecated AES encryption, migrate to age encryption", map[string]interface{}{})
    // Implementation remains for backward compatibility
    return i.encryptAES(data, key)
}

func (i *Integration) Decrypt(ctx context.Context, data []byte, key []byte) ([]byte, error) {
    i.logger.Warn("Using deprecated AES decryption, migrate to age decryption", map[string]interface{}{})
    // Implementation remains for backward compatibility
    return i.decryptAES(data, key)
}

func (i *Integration) ValidateKey(ctx context.Context, key []byte) error {
    i.logger.Warn("Using deprecated AES key validation, migrate to age key validation", map[string]interface{}{})
    // Implementation remains for backward compatibility
    return i.validateAESKey(key)
}
```

### 3. CLI Commands

**NOTE: This section contains fabricated code that doesn't match spooky's actual CLI structure. CLI commands need to be designed based on actual spooky patterns.**

The actual spooky CLI commands are:
- `spooky project` - Project management
- `spooky actions` - Action management  
- `spooky variables` - Variable management
- `spooky machines` - Machine management
- `spooky facts` - Facts management
- `spooky schemas` - Schema management
- `spooky integrations` - Integration management

Any secrets CLI commands would need to follow these existing patterns and integrate with the actual spooky CLI architecture.

## Migration Strategy

### Phase 1: Parallel Support (Weeks 1-4)
- Maintain both AES and age encryption during transition
- Add deprecation warnings for AES methods
- Provide migration tools and documentation

### Phase 2: Gradual Migration (Weeks 5-6)
- Encourage users to migrate to age encryption
- Provide examples and migration guides
- Update documentation to prioritize age

### Phase 3: Deprecation (Week 7+)
- Mark AES methods as deprecated
- Remove AES support in future major version
- Complete migration to age-only encryption

## Testing Strategy

### Unit Tests
- Age encryption/decryption with various key types
- Key validation and parsing
- Error handling and edge cases
- Performance benchmarks

### Integration Tests
- End-to-end encryption workflows
- CLI command testing
- Configuration integration
- Cross-platform compatibility

### Security Tests
- Key management security
- Encryption strength validation
- Audit logging verification
- Access control testing

## Documentation Requirements

### User Documentation
- Age encryption setup guide
- Key management best practices
- Migration guide from AES to age
- Troubleshooting common issues

### Developer Documentation
- API reference for age integration
- Integration examples
- Security considerations
- Performance guidelines

### Operational Documentation
- Key rotation procedures
- Backup and recovery
- Monitoring and alerting
- Incident response

## Success Criteria

### Functional Requirements
- [ ] Variables can be encrypted and decrypted with age
- [ ] Facts can be decrypted from `/etc/spooky/facts.*` files (read-only, not stored on disk)
- [ ] Machine inventory secrets can be encrypted/decrypted
- [ ] Support for multiple recipients
- [ ] SSH key integration works
- [ ] Passphrase fallback works
- [ ] CLI commands function correctly

### Performance Requirements
- [ ] Age encryption/decryption performance meets benchmarks
- [ ] No significant performance regression from AES
- [ ] Memory usage remains reasonable
- [ ] Startup time not significantly impacted

### Security Requirements
- [ ] Age encryption provides equivalent or better security than AES
- [ ] Key management follows security best practices
- [ ] Audit logging captures all operations
- [ ] No sensitive data exposed in logs

### Usability Requirements
- [ ] Migration from AES is straightforward
- [ ] CLI commands are intuitive
- [ ] Error messages are helpful
- [ ] Documentation is comprehensive

## Risk Mitigation

### Technical Risks
- **Age library compatibility**: Test with multiple age versions
- **Performance impact**: Benchmark and optimize
- **Key management complexity**: Provide clear documentation and tools

### Operational Risks
- **Migration complexity**: Provide automated migration tools
- **User adoption**: Create comprehensive examples and guides
- **Backward compatibility**: Maintain parallel support during transition

### Security Risks
- **Key exposure**: Follow security best practices
- **Audit trail gaps**: Comprehensive logging
- **Access control**: Proper file permissions and validation

## Conclusion

This implementation plan provides a comprehensive approach to integrating age encryption into spooky. The phased approach ensures minimal disruption while providing robust secrets management capabilities. The plan addresses all major use cases (variables, facts, machine inventories) while maintaining security and usability standards.

The implementation will significantly improve spooky's security posture by leveraging age's modern encryption standards and providing flexible key management options. The migration strategy ensures existing users can transition smoothly while new users benefit from the enhanced security features.
