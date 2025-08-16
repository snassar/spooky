# Secrets Management Implementation Plan

## Overview

This document outlines the implementation plan for integrating [age encryption](https://github.com/FiloSottile/age) into spooky for comprehensive secrets management. The plan covers implementing age-based encryption for variables, facts, and machine inventories with explicit decryption control.

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

### Objectives
1. **Replace AES-GCM with age encryption** for all secrets management
2. **Support multiple encryption scenarios**:
   - Encrypted variables in `variables.hcl` (encrypt/decrypt)
   - Encrypted facts in `/etc/spooky/custom.hcl` on target machines (decrypt only)
   - Encrypted machine inventory secrets (passphrases, keys)
3. **Maintain explicit decryption control** - no automatic decryption
4. **Provide comprehensive CLI support** for age key management
5. **Integrate with existing spooky architecture**
6. **Support multiple recipients** for encrypted data
7. **Provide audit logging** for encryption/decryption operations
8. **Integrate with existing validation systems**

## Technical Architecture

### Age Integration Design

#### Core Components
```go
// Age-focused SecretsIntegration interface (breaking change)
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
}
```

#### Configuration Structure
```hcl
# Enhanced age configuration in spooky.hcl
age {
  # Primary age key configuration
  identities = "~/.config/spooky/identities"
  recipients = "~/.config/spooky/recipients.txt"
  
  # Security settings
  key_validation = true
  recipient_validation = true
}
```

### Use Case Implementations

#### 1. Encrypted Variables (Encrypt/Decrypt)
```hcl
# variables.hcl with encrypted values
variables {
  variable "database_password" {
    value = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
    encrypted = true
    description = "Database password for production"
  }
  
  variable "api_key" {
    value = "sk-1234567890abcdef"
    description = "API key for external service" # encrypted = false is omitted (defaults to false)
  }
  
  variable "secrets" {
    type = "object"
    value = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
    encrypted = true
    description = "Encrypted configuration object"
  }
}
```

#### 2. Encrypted Facts (Decrypt Only)
```hcl
# /etc/spooky/custom.hcl on target machines
custom {
  # Plaintext environment
  environment = "production"
  
  # Encrypted database connection (automatically detected by age1 prefix)
  database_connection_string = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
  
  # Application info with encrypted secrets (automatically detected by age1 prefix)
  application = {
    name = "web-app"
    version = "1.2.3"
    api_key = "age1abc123def456ghi789jkl012mno345pqr678stu901vwx234yz"
  }
}
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
      passphrase = {
        value = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
        encrypted = true
      }
    }
  }
  
  machine "db-server" {
    hostname = "db.example.com"
    port = 22
    user = "postgres"
    
    authentication {
      method = "password"
      password = {
        value = "age1abc123def456ghi789jkl012mno345pqr678stu901vwx234yz"
        encrypted = true
      }
    }
  }
}
```

#### 4. Actions Run with Decryption
```bash
# Run actions with decryption enabled for secret material
spooky actions run my-project --decrypt

# This allows the orchestrator to use encrypted variables and facts
# but ensures secrets are never logged or displayed in plaintext
```

## Key Management Strategy

### Key Generation
- **Spooky generates**: No - users should use age CLI tools
- **User provides**: Yes - via age CLI: `age-keygen -o ~/.config/spooky/identities/identity.txt`
- **Key format**: age1... public keys in recipients.txt files

### Key Storage
- **Private keys**: Age identity files can contain multiple keys, one per line, with comments starting with #
- **Public keys**: recipients.txt files with one public key per line, no comments
- **Permissions**: 600 for identity files, 644 for recipients.txt
- **Backup**: Users responsible for backing up ~/.config/spooky/identities/ directory

### Key Rotation
- **Supported**: Yes. Age encryption is additive - when you encrypt to multiple recipients, the data can be decrypted by ANY of those recipients
- **Automatic**: No - manual process via spooky variables encrypt command
- **Manual process**: 1. Add new recipient to recipients.txt, 2. Run spooky variables encrypt, 3. Remove old recipient from recipients.txt
- **Migration**: Single re-encryption with spooky variables encrypt after updating recipients.txt

### Key Management Implementation Details

#### Identity File Format
- **Format**: Follow age documentation format exactly
- **Example**:
  ```
  # This is a comment
  AGE-SECRET-KEY-1GQZ8TGD35TCWUPJQWF2E9Y62WR73QSH2SJ7K3KM53G3Q0MFRQCGS6T6PSG
  # Another comment
  AGE-SECRET-KEY-1GQZ8TGD35TCWUPJQWF2E9Y62WR73QSH2SJ7K3KM53G3Q0MFRQCGS6T6PSG
  ```
- **Spooky's role**: Delegate entirely to age library - do NOT parse these files
- **Implementation**: Use `filippo.io/age` library's built-in identity file parsing

#### Recipients File Handling
- **Project-level recipients.txt**: Optional - if it doesn't exist, spooky continues without it
- **Global recipients.txt**: Required for encryption operations - spooky should error if missing
- **Malformed files**: Delegate to age library - let age library handle validation
- **Error handling**: If age library rejects the file, spooky should report the age library's error

#### Recipient Key Validation
- **Delegate validation to age library** - that's why we're using `filippo.io/age`
- **No custom validation** - spooky should not try to validate age keys itself
- **Pass-through approach**: Read recipients.txt, pass to age library, let age library validate

## CLI Command Design

### Core Commands

#### spooky project encrypt
- **Syntax**: `spooky project encrypt <project> [flags]`
- **Flags**: `--dry-run` - Show what would be encrypted without making changes
- **Description**: Encrypt all variables and machines in project that have encrypted=true, and re-encrypt if identities/recipients changed (on-disk changes). Objects and maps are serialized before encryption.

#### spooky variables encrypt
- **Syntax**: `spooky variables encrypt <project> [flags]`
- **Flags**: `--dry-run` - Show what would be encrypted without making changes
- **Description**: Encrypt all variables in project that have encrypted=true, and re-encrypt if identities/recipients changed (on-disk changes). Objects and maps are serialized before encryption.

#### spooky machines encrypt
- **Syntax**: `spooky machines encrypt <project> [flags]`
- **Flags**: `--dry-run` - Show what would be encrypted without making changes
- **Description**: Encrypt all machines in project that have encrypted=true, and re-encrypt if identities/recipients changed (on-disk changes). Objects and maps are serialized before encryption.

#### --decrypt flag
- **Syntax**: `--decrypt` (flag on spooky actions run)
- **Description**: Decrypt variables in-memory for debugging (no on-disk changes)
- **Usage**: `spooky actions run . --decrypt --dry-run`

#### spooky secrets validate
- **Syntax**: `spooky secrets validate <project>`
- **Description**: Validate age configuration and keys for project

### Backend Architecture
- **Shared backend**: Single encrypt backend function that handles encryption for all scopes
- **Scope parameter**: Backend takes scope parameter: 'variables', 'machines', or 'all'
- **Implementation**: One implementation, multiple CLI entry points
- **Code reuse**: No code duplication between encrypt commands

### CLI Command Behavior

#### Encrypt Commands Behavior
- **When no `encrypted = true` fields are found**: Command should succeed and report "nothing to do"
- **Not silent** - provide clear feedback about what was checked
- **Example output**: `"Checked 15 variables, 0 require encryption. Nothing to do."`

#### Encryption Logic
- **If value is already encrypted** (age1... format): Re-encrypt with current recipients
- **If value is not encrypted** but `encrypted = true`: Encrypt and replace
- **If value is not encrypted** and `encrypted = false` or omitted: Skip

#### Re-encryption Detection
- **Recipients are always determined by**: `$XDG_CONFIG_HOME/spooky/spooky.hcl` + project-level `recipients.txt` (if exists)
- **Detection method**: Always re-encrypt - don't try to detect changes
- **Rationale**: Simpler and more reliable than change detection
- **Behavior**: Every `spooky variables encrypt` re-encrypts all `encrypted = true` values with current recipients

#### --dry-run Behavior
- **Show exactly what would be encrypted**
- **Format**: List each field that would be encrypted/re-encrypted
- **Example output**:
  ```
  Would encrypt/re-encrypt:
  - variables.hcl: database_password (string)
  - variables.hcl: api_secrets (object) 
  - machines.hcl: web-server.passphrase (string)
  - machines.hcl: db-server.password (string)
  
  Total: 4 values would be processed
  ```

#### Decryption Flag Behavior
- **When no encrypted values are found**: Don't attempt decryption - just continue normally
- **No error** - this is a valid scenario
- **Example output**: `"No encrypted values found. Continuing without decryption."`

#### Error Handling Strategy
- **Continue and report all errors** - don't fail fast
- **Collect all decryption errors** and report them together
- **Example output**:
  ```
  Decryption errors:
  - variables.hcl: database_password - invalid age format
  - machines.hcl: web-server.passphrase - missing identity file
  - machines.hcl: db-server.password - decryption failed
  
  Total: 3 decryption errors
  ```

## Schema Updates

### Variables Schema Updates
- **Encrypted field**: Variables use encrypted=true field to indicate age-encrypted values
- **Validation**: Use application-level validation with age library for encrypted values
- **Serialization**: Objects and maps are serialized to HCL/JSON before encryption, deserialized after decryption
- **Mixed content**: Support for mixing encrypted and plain values (encrypted=true for sensitive, omitted for plaintext)
- **Object mixing**: Objects and maps cannot mix encrypted/plaintext - they are either completely encrypted or completely plaintext

### Variable Types Support
- **String encryption**: Yes - primary use case
- **Number encryption**: Yes - useful for sensitive numbers
- **Boolean encryption**: No - booleans are too simple to encrypt
- **Object encryption**: Yes - entire object encrypted as one blob, then deserialized in memory
- **Map encryption**: Yes - entire map encrypted as one blob, then deserialized in memory

### Facts Schema Updates
- **Encrypted facts**: Support age1... strings as values for any custom fact
- **Validation**: Use application-level validation with age library (no regex validation)
- **Sources**: Support for encrypted facts in /etc/spooky/custom.hcl
- **Types**: Only custom facts can be encrypted, not system or machine facts
- **Detection method**: age1 prefix is the ONLY detection method for facts
- **Mixed encrypted/plaintext objects**: Objects and maps can ONLY be encrypted as a whole - no mixing within the same object

### Machines Schema Updates
- **Authentication encryption**: SSH key passphrases and passwords can be encrypted
- **Secret fields**: Use nested value and encrypted fields for password and passphrase
- **Mixed authentication**: Support for mixing encrypted and plain authentication
- **Validation**: Use application-level validation with age library for encrypted values

### Schema Validation Strategy

#### Application-Level Validation
```go
// Application-level validation using age library
func validateAgeEncryptedValue(value string) error {
    // Let the age library handle validation
    // This will catch malformed age strings, invalid recipients, etc.
    _, err := age.ParseRecipients(strings.NewReader(value))
    if err != nil {
        return fmt.Errorf("invalid age-encrypted value: %w", err)
    }
    return nil
}
```

#### Schema-Level Validation (Limited)
- **Basic structure validation**: Only validate field exists, is string, etc.
- **No regex validation**: Remove age1... pattern validation from schemas
- **No cryptographic validation**: Let age library handle all cryptographic validation

#### Validation Flow
1. **Schema validation**: Only validate basic structure (field exists, is string, etc.)
2. **Application validation**: Use age library to validate encrypted values during processing
3. **Error handling**: Let age library errors bubble up with clear context

#### Schema Updates Needed
```hcl
# Remove this from facts.schema.hcl and other schemas
# age_encrypted_values = {
#   rule = "regex"
#   pattern = "^age1[a-zA-Z0-9]+"
#   message = "Age-encrypted values must start with 'age1'"
#   apply_to = ["custom.*"]
# }
```

## Security and Logging Protection

### Critical Security Requirements
- **No automatic decryption**: Decryption only happens with explicit `--decrypt` flag
- **Logging protection**: Prevent decrypted values in logs (critical security requirement)
- **Age string safety**: Age1... strings are safe to log (they are encrypted and designed for exposure)

### Redaction Patterns
- **Age strings**: Age1... strings are safe to log (they are encrypted and designed for exposure)
- **Decrypted values**: Redact all decrypted variable values from logs (replace with [REDACTED_VALUE])
- **Object values**: Redact decrypted object/map values from logs (replace with [REDACTED_OBJECT])
- **Sensitive fields**: Redact field names that contain 'password', 'secret', 'key', 'token' (replace with [REDACTED_FIELD])

### Protection Methods
- **Pre-logging**: Scan all log messages before output for decrypted values (age1... strings are safe)
- **Post-logging**: Scan log files after writing for any leaked decrypted secrets
- **Field filtering**: Filter sensitive field names from structured logging
- **Value sanitization**: Sanitize all decrypted variable values before logging

### Log Levels
- **Debug**: No decrypted values in debug logs (only redacted placeholders)
- **Info**: No decrypted values in info logs (only redacted placeholders)
- **Error**: No decrypted values in error logs (only redacted placeholders)
- **Trace**: No decrypted values in trace logs (only redacted placeholders)

### Logging Integration Implementation

#### Redaction Patterns Implementation
```go
// Redaction patterns for different data types
type RedactionPatterns struct {
    // Age strings are safe to log (they're encrypted)
    AgeStringPattern *regexp.Regexp
    
    // Decrypted values to redact
    DecryptedValuePattern *regexp.Regexp
    
    // Sensitive field names to redact
    SensitiveFieldPattern *regexp.Regexp
    
    // Object/map values to redact
    ObjectValuePattern *regexp.Regexp
}

func NewRedactionPatterns() *RedactionPatterns {
    return &RedactionPatterns{
        AgeStringPattern:     regexp.MustCompile(`^age1[a-zA-Z0-9]+`),
        DecryptedValuePattern: regexp.MustCompile(`(?i)(password|secret|key|token|credential)`),
        SensitiveFieldPattern: regexp.MustCompile(`(?i)(password|secret|key|token|credential|private_key|ssh_key|auth_key)`),
        ObjectValuePattern:    regexp.MustCompile(`\{.*\}`), // Simple object detection
    }
}
```

#### Pre-Logging Scanning Implementation
```go
// Secure logger that scans before output
type SecureLogger struct {
    logger     spookylogging.Logger
    patterns   *RedactionPatterns
    redactMode bool
}

func (l *SecureLogger) Info(msg string, fields ...spookylogging.Field) {
    l.logger.Info(msg, l.sanitizeFields(fields)...)
}

func (l *SecureLogger) sanitizeFields(fields []spookylogging.Field) []spookylogging.Field {
    sanitized := make([]spookylogging.Field, len(fields))
    
    for i, field := range fields {
        if l.shouldRedactField(field) {
            sanitized[i] = spookylogging.String(field.Key, "[REDACTED]")
        } else {
            sanitized[i] = l.sanitizeValue(field)
        }
    }
    
    return sanitized
}

func (l *SecureLogger) shouldRedactField(field spookylogging.Field) bool {
    // Redact sensitive field names
    if l.patterns.SensitiveFieldPattern.MatchString(field.Key) {
        return true
    }
    
    // Redact if in decryption mode and value is not age-encrypted
    if l.redactMode {
        switch v := field.Value.(type) {
        case string:
            if !l.patterns.AgeStringPattern.MatchString(v) {
                return true
            }
        case map[string]interface{}, []interface{}:
            return true // Always redact objects/maps when decrypted
        }
    }
    
    return false
}

func (l *SecureLogger) sanitizeValue(field spookylogging.Field) spookylogging.Field {
    if !l.redactMode {
        return field // No redaction needed
    }
    
    switch v := field.Value.(type) {
    case string:
        if l.patterns.AgeStringPattern.MatchString(v) {
            return field // Age strings are safe
        }
        return spookylogging.String(field.Key, "[REDACTED_VALUE]")
        
    case map[string]interface{}:
        return spookylogging.String(field.Key, "[REDACTED_OBJECT]")
        
    case []interface{}:
        return spookylogging.String(field.Key, "[REDACTED_ARRAY]")
        
    default:
        return field
    }
}
```

#### Post-Logging Scanning Implementation
```go
// Post-logging scanner for leaked secrets
type PostLogScanner struct {
    patterns *RedactionPatterns
    logFile  string
}

func (s *PostLogScanner) ScanForLeaks() error {
    data, err := os.ReadFile(s.logFile)
    if err != nil {
        return fmt.Errorf("failed to read log file: %w", err)
    }
    
    var leaks []string
    
    // Scan for decrypted values that shouldn't be in logs
    lines := strings.Split(string(data), "\n")
    for i, line := range lines {
        if s.containsLeakedSecret(line) {
            leaks = append(leaks, fmt.Sprintf("line %d: %s", i+1, line[:50]+"..."))
        }
    }
    
    if len(leaks) > 0 {
        return fmt.Errorf("potential secret leaks detected in logs:\n%s", strings.Join(leaks, "\n"))
    }
    
    return nil
}

func (s *PostLogScanner) containsLeakedSecret(line string) bool {
    // Look for patterns that suggest decrypted secrets
    if s.patterns.DecryptedValuePattern.MatchString(line) {
        // Check if it's not an age string
        if !s.patterns.AgeStringPattern.MatchString(line) {
            return true
        }
    }
    
    return false
}
```

#### Integration with Existing Logging System
```go
// Integration with spooky's logging system
func SetupSecureLogging(ctx context.Context, logger spookylogging.Logger, redactMode bool) spookylogging.Logger {
    patterns := NewRedactionPatterns()
    
    secureLogger := &SecureLogger{
        logger:     logger,
        patterns:   patterns,
        redactMode: redactMode,
    }
    
    // Set up post-logging scanner
    if redactMode {
        scanner := &PostLogScanner{
            patterns: patterns,
            logFile:  getLogFilePath(), // Get from logging config
        }
        
        // Schedule post-scan
        go func() {
            defer scanner.ScanForLeaks()
            <-ctx.Done()
        }()
    }
    
    return secureLogger
}

// Usage in CLI commands
func runWithSecureLogging(ctx context.Context, project string, decrypt bool) error {
    logger := spookylogging.GetLogger()
    
    if decrypt {
        logger = SetupSecureLogging(ctx, logger, true)
    }
    
    // Use logger for all operations
    logger.Info("Processing project", "project", project, "decrypt", decrypt)
    
    return nil
}
```

## Error Handling and Validation

### Application-Level Validation Strategy
- **Age Library Authority**: Use `filippo.io/age` library for all cryptographic validation
- **No Schema Regex Validation**: Remove age1... pattern validation from schemas
- **Comprehensive Validation**: Let age library handle format, recipient, and cryptographic validation
- **Clear Error Messages**: Use age library error messages with spooky context

### Validation Implementation
```go
// Application-level validation using age library
func (s *SecretsIntegration) ValidateAgeEncryptedValue(ctx context.Context, value string) error {
    // Let the age library handle all validation
    // This will catch malformed age strings, invalid recipients, etc.
    _, err := age.ParseRecipients(strings.NewReader(value))
    if err != nil {
        return fmt.Errorf("invalid age-encrypted value: %w", err)
    }
    return nil
}

// Validation during processing
func (s *SecretsIntegration) ProcessEncryptedVariable(value string) error {
    // Validate using age library
    if err := s.ValidateAgeEncryptedValue(context.Background(), value); err != nil {
        return fmt.Errorf("failed to validate encrypted variable: %w", err)
    }
    
    // Process the validated value
    return nil
}
```

### Error Scenarios
- **Malformed Age Strings**: Error out, log and display - don't continue with invalid data
- **Invalid Recipients**: Error out, log and display - fail fast on invalid recipients
- **Decryption Failures**: Error out, log and display - don't continue with failed decryption
- **Missing Identity Files**: Error out, log and display - critical for decryption operations
- **Incorrect Permissions**: Error out, log and display - security-critical issue
- **Individual Decryption Failures**: Error out, log and display - don't continue processing other variables/facts

### Error Handling Implementation
```go
func (s *SecretsIntegration) decryptValue(value string, identityPath string) ([]byte, error) {
    // Check identity file exists
    if _, err := os.Stat(identityPath); os.IsNotExist(err) {
        return nil, fmt.Errorf("identity file not found: %s", identityPath)
    }
    
    // Check identity file permissions
    info, err := os.Stat(identityPath)
    if err != nil {
        return nil, fmt.Errorf("failed to stat identity file: %w", err)
    }
    
    mode := info.Mode()
    if mode&0077 != 0 {
        return nil, fmt.Errorf("identity file has incorrect permissions: %s (expected 600, got %v)", 
            identityPath, mode)
    }
    
    // Validate using age library first
    if err := s.ValidateAgeEncryptedValue(context.Background(), value); err != nil {
        return nil, fmt.Errorf("invalid age-encrypted value: %w", err)
    }
    
    // Attempt decryption with age library
    decrypted, err := age.Decrypt(value, identityPath)
    if err != nil {
        return nil, fmt.Errorf("decryption failed: %w", err)
    }
    
    return decrypted, nil
}
```

### CLI Error Output Examples

**Missing Identity File:**
```bash
$ spooky actions run my-project --decrypt
Error: identity file not found: ~/.config/spooky/identities/identity.txt
Please run: age-keygen -o ~/.config/spooky/identities/identity.txt
```

**Incorrect Permissions:**
```bash
$ spooky actions run my-project --decrypt
Error: identity file has incorrect permissions: ~/.config/spooky/identities/identity.txt (expected 600, got 644)
Please run: chmod 600 ~/.config/spooky/identities/identity.txt
```

**Decryption Failure:**
```bash
$ spooky actions run my-project --decrypt
Error: failed to decrypt variable 'database_password': invalid age format
```

**Invalid Recipients:**
```bash
$ spooky variables encrypt my-project
Error: invalid recipient in recipients.txt: age1xyz... (age library error: unsupported key type)
```

## Backward Compatibility

### Clean Break to Age Encryption
- **AES-GCM Removal**: Complete removal of AES-GCM encryption - no deprecation timeline
- **No migration path**: Not necessary since we're removing AES-GCM entirely
- **Breaking change**: This is a major version change that replaces AES-GCM with age encryption only
- **Interface replacement**: Complete replacement of SecretsIntegration interface with age-focused methods

### Implementation Approach
```go
// Complete replacement of SecretsIntegration interface
type SecretsIntegration interface {
    // Age-specific methods only
    EncryptWithAge(ctx context.Context, data []byte, recipients []string) ([]byte, error)
    DecryptWithAge(ctx context.Context, data []byte, identityPath string) ([]byte, error)
    EncryptWithPassphrase(ctx context.Context, data []byte, passphrase string) ([]byte, error)
    DecryptWithPassphrase(ctx context.Context, data []byte, passphrase string) ([]byte, error)
    
    // Key management
    ValidateAgeKey(ctx context.Context, keyPath string) error
    ListRecipients(ctx context.Context, encryptedData []byte) ([]string, error)
    
    // Application-level validation
    ValidateAgeEncryptedValue(ctx context.Context, value string) error
}
```

### Migration Strategy
- **New installations**: Use age encryption from the start
- **Existing users**: Need to re-encrypt any existing AES-GCM data with age
- **Documentation**: Provide clear migration guide for users with existing encrypted data
- **Version bump**: This change requires a major version bump (e.g., 1.0.0 to 2.0.0)

## Security Implementation Details

### Memory Management
- **Decrypted Value Lifetime**: Only as long as `spooky actions run` takes - clear immediately after use
- **Per-variable basis**: Decrypt, use, clear, move to next variable
- **No persistent storage**: No persistent storage of decrypted values in memory

### Modern Memory Clearing in Go
```go
// Secure memory clearing utilities
type SecureMemory struct{}

// Clear string from memory
func (sm *SecureMemory) ClearString(s *string) {
    if s == nil {
        return
    }
    
    // Convert to bytes for clearing
    bytes := []byte(*s)
    sm.ClearBytes(bytes)
    
    // Clear the string reference
    *s = ""
}

// Clear bytes from memory
func (sm *SecureMemory) ClearBytes(b []byte) {
    if b == nil {
        return
    }
    
    // Zero out the memory
    for i := range b {
        b[i] = 0
    }
    
    // Use runtime.KeepAlive to prevent optimization
    runtime.KeepAlive(b)
}

// Clear sensitive struct fields
func (sm *SecureMemory) ClearVariable(v *spookytypes.Variable) {
    if v == nil {
        return
    }
    
    // Clear resolved value if it's a string
    if str, ok := v.ResolvedValue.(string); ok {
        sm.ClearString(&str)
        v.ResolvedValue = nil
    }
    
    // Clear other sensitive fields
    sm.ClearString(&v.Description)
    if v.Metadata != nil {
        for key, value := range v.Metadata {
            if str, ok := value.(string); ok {
                sm.ClearString(&str)
            }
            delete(v.Metadata, key)
        }
    }
}
```

### Integration with Variable Processing
```go
// Secure variable processor
type SecureVariableProcessor struct {
    memory *SecureMemory
}

func (p *SecureVariableProcessor) ProcessVariable(ctx context.Context, variable *spookytypes.Variable, decrypt bool) error {
    defer func() {
        // Always clear sensitive data when done
        p.memory.ClearVariable(variable)
    }()
    
    if !decrypt || !variable.Encrypted {
        return nil // No decryption needed
    }
    
    // Decrypt the value
    decrypted, err := p.decryptValue(variable.Default.(string))
    if err != nil {
        return fmt.Errorf("failed to decrypt variable %s: %w", variable.Name, err)
    }
    
    // Use the decrypted value
    variable.ResolvedValue = decrypted
    
    // Process the variable...
    // (variable is used here)
    
    return nil
}

func (p *SecureVariableProcessor) decryptValue(encryptedValue string) (string, error) {
    // Decrypt using age
    decrypted, err := age.Decrypt(encryptedValue, identityPath)
    if err != nil {
        return "", err
    }
    
    // Convert to string
    result := string(decrypted)
    
    // Clear the decrypted bytes immediately
    p.memory.ClearBytes(decrypted)
    
    return result, nil
}
```

### Memory Clearing Best Practices
```go
// Memory clearing patterns for different data types
type MemoryClearingPatterns struct {
    memory *SecureMemory
}

// Clear sensitive strings
func (m *MemoryClearingPatterns) ClearSensitiveString(s *string) {
    defer m.memory.ClearString(s)
    // Use the string here
}

// Clear sensitive byte slices
func (m *MemoryClearingPatterns) ClearSensitiveBytes(b []byte) {
    defer m.memory.ClearBytes(b)
    // Use the bytes here
}

// Clear sensitive maps
func (m *MemoryClearingPatterns) ClearSensitiveMap(data map[string]interface{}) {
    defer func() {
        for key, value := range data {
            if str, ok := value.(string); ok {
                m.memory.ClearString(&str)
            }
            delete(data, key)
        }
    }()
    // Use the map here
}

// Clear sensitive structs
func (m *MemoryClearingPatterns) ClearSensitiveStruct(v interface{}) {
    defer func() {
        // Use reflection to clear all string fields
        m.clearStructFields(v)
    }()
    // Use the struct here
}

func (m *MemoryClearingPatterns) clearStructFields(v interface{}) {
    val := reflect.ValueOf(v)
    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }
    
    typ := val.Type()
    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)
        if field.Kind() == reflect.String {
            if field.CanSet() {
                field.SetString("")
            }
        }
    }
}
```

### Integration with CLI Commands
```go
// Secure CLI command execution
func runActionsWithSecureMemory(ctx context.Context, project string, decrypt bool) error {
    memory := &SecureMemory{}
    processor := &SecureVariableProcessor{memory: memory}
    
    // Ensure cleanup on exit
    defer func() {
        // Final memory cleanup
        runtime.GC()
        runtime.KeepAlive(memory)
    }()
    
    // Process variables with secure memory management
    variables, err := loadVariables(project)
    if err != nil {
        return err
    }
    
    for _, variable := range variables {
        if err := processor.ProcessVariable(ctx, variable, decrypt); err != nil {
            return err
        }
        // Variable is automatically cleared after processing
    }
    
    return nil
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
- [ ] Replace AES-GCM interface with age-focused interface (breaking change)
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
- [ ] Update variable validation to use age library (no regex validation)

#### 2.2 Variable Integration
- [ ] Modify `internal/variables/` to support age encryption/decryption
- [ ] Add explicit decryption control (no automatic decryption)
- [ ] Create variable encryption/decryption CLI commands
- [ ] Add variable encryption/decryption examples

### Phase 3: Facts Decryption (Week 4)

#### 3.1 Facts Schema Updates
- [ ] Update `internal/schemas/schemas/facts.schema.hcl`
- [ ] Remove regex validation for age-encrypted values
- [ ] Add age decryption support to fact types
- [ ] Create fact decryption logic (read-only)
- [ ] Update fact validation to use age library

#### 3.2 Facts Integration
- [ ] Modify `internal/facts/` to support age decryption
- [ ] Add automatic age1 prefix detection for custom facts
- [ ] Create fact decryption logic (read-only)
- [ ] Add fact decryption examples

### Phase 4: Machine Inventory Encryption (Week 5)

#### 4.1 Machine Schema Updates
- [ ] Update `internal/schemas/schemas/machines.schema.hcl`
- [ ] Add age encryption support to machine authentication types
- [ ] Create machine encryption/decryption logic
- [ ] Update machine validation to use age library

#### 4.2 Machine Integration
- [ ] Modify `internal/machines/` to support age encryption/decryption
- [ ] Add encryption support for SSH passphrases and passwords
- [ ] Create machine encryption/decryption CLI commands
- [ ] Add machine encryption/decryption examples

### Phase 5: CLI Commands and Integration (Week 6)

#### 5.1 CLI Command Implementation
- [ ] Implement `spooky project encrypt` command
- [ ] Implement `spooky variables encrypt` command
- [ ] Implement `spooky machines encrypt` command
- [ ] Implement `spooky secrets validate` command
- [ ] Add `--decrypt` flag to `spooky actions run`

#### 5.2 Integration Testing
- [ ] Create comprehensive integration tests
- [ ] Test all CLI commands with various scenarios
- [ ] Test error handling and edge cases
- [ ] Test memory management and security

### Phase 6: Logging and Security (Week 7)

#### 6.1 Logging Integration
- [ ] Implement secure logging with redaction patterns
- [ ] Add pre-logging scanning for decrypted values
- [ ] Add post-logging scanning for leaked secrets
- [ ] Integrate with existing logging system

#### 6.2 Security Implementation
- [ ] Implement secure memory management
- [ ] Add memory clearing utilities
- [ ] Test memory security and cleanup
- [ ] Add security audit logging

### Phase 7: Documentation and Examples (Week 8)

#### 7.1 Documentation Updates
- [ ] Update API documentation for secrets management
- [ ] Create user guides for encryption/decryption
- [ ] Add troubleshooting guides
- [ ] Update configuration examples

#### 7.2 Example Implementation
- [ ] Create comprehensive examples for all use cases
- [ ] Add example projects with encrypted variables
- [ ] Add example machine inventories with encrypted authentication
- [ ] Add example custom facts with encrypted values

## Security Best Practices

### Memory Management
- **Zero out memory** - don't just set to empty string
- **Use defer for cleanup** - ensure cleanup happens
- **Clear at multiple levels** - strings, bytes, structs, maps
- **Prevent optimization** - use `runtime.KeepAlive` to prevent compiler optimization
- **Force garbage collection** - call `runtime.GC()` after sensitive operations

### Error Handling
- **Error out, log and display** for all validation/decryption failures
- **Fail fast** - stop processing on first error
- **Clear error messages** with specific details and suggested fixes
- **Use age library error messages** when available

### Logging Security
- **Pre-logging scanning** - sanitize before output
- **Post-logging scanning** - detect leaks after writing
- **Configurable redaction** - patterns and modes
- **Integration with existing logger** - wrap current logging system
- **Context-aware redaction** - only redact when decryption is active

## Conclusion

This implementation plan provides a comprehensive approach to integrating age encryption into spooky for secure secrets management. The plan addresses all major concerns including schema clarity, key management, CLI behavior, error handling, integration points, backward compatibility, and security implementation details.

The phased approach ensures systematic implementation while maintaining security and usability throughout the development process.