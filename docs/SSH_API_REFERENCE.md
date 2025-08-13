# SSH System API Reference

## Overview

This document provides a comprehensive API reference for the spooky SSH system. It covers all interfaces, types, methods, and implementation details for developers working with the SSH system.

## Table of Contents

1. [Core Interfaces](#core-interfaces)
2. [Advanced SSH Capabilities](#advanced-ssh-capabilities)
3. [Type Definitions](#type-definitions)
4. [Implementation Details](#implementation-details)
5. [Error Handling](#error-handling)
6. [Key Validation Rules](#key-validation-rules)
7. [Certificate Handling](#certificate-handling)
8. [File Transfer](#file-transfer)
9. [Advanced Authentication](#advanced-authentication)
11. [Examples](#examples)

## Advanced SSH Capabilities

The spooky SSH system now includes advanced capabilities for file transfer (as a library component for actions) and multi-factor authentication.

### FileTransferManager

Manages SFTP and SCP file transfers with progress tracking and verification.

```go
type FileTransferManager struct {
	client *Client
    logger spookytypeslogging.Logger
}

func NewFileTransferManager(client *Client, logger spookytypeslogging.Logger) *FileTransferManager
```

**Methods:**

#### TransferFile
```go
TransferFile(ctx context.Context, connection *spookytypes.Connection, transfer *spookytypesssh.FileTransfer) (*spookytypesssh.FileTransferResult, error)
```
Transfers a file using SFTP or SCP protocol.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `connection`: SSH connection details
- `transfer`: File transfer configuration

**Returns:**
- `*spookytypesssh.FileTransferResult`: Transfer result with metrics
- `error`: Error if transfer fails

**Features:**
- Supports SFTP and SCP protocols
- Progress tracking and logging
- Post-transfer verification
- Automatic retry on failure

#### BatchTransfer
```go
BatchTransfer(ctx context.Context, connection *spookytypes.Connection, transfers []*spookytypesssh.FileTransfer) ([]*spookytypesssh.FileTransferResult, error)
```
Performs multiple file transfers concurrently.



### AdvancedAuthManager

Manages multi-factor authentication and certificate-based authentication.

```go
type AdvancedAuthManager struct {
    logger spookytypeslogging.Logger
    agent  *Agent
}

func NewAdvancedAuthManager(logger spookytypeslogging.Logger) *AdvancedAuthManager
```

**Methods:**

#### GetAuthMethods
```go
GetAuthMethods(config *MultiFactorAuthConfig) ([]ssh.AuthMethod, error)
```
Creates authentication methods for multi-factor authentication.

#### GenerateCertificate
```go
GenerateCertificate(config *CertificateConfig) (*ssh.Certificate, error)
```
Generates SSH certificates for certificate-based authentication.

#### SaveCertificate
```go
SaveCertificate(cert *ssh.Certificate, path string) error
```
Saves an SSH certificate to file.



## Core Interfaces

### SSHManager

The primary interface for SSH management operations.

```go
type SSHManager interface {
    // CreateClient creates a new SSH client with the given configuration
    CreateClient(ctx context.Context, config *spookytypes.ClientConfig) (*spookytypes.Client, error)
    
    // Connect establishes an SSH connection to the given host
    Connect(ctx context.Context, request *spookytypes.ConnectionRequest) (*spookytypes.ConnectionResult, error)
    
    // Authenticate authenticates with the given credentials
    Authenticate(ctx context.Context, connection *spookytypes.Connection, auth *spookytypes.Authentication) (*spookytypes.AuthenticationResult, error)
    
    // CreateSession creates a new SSH session
    CreateSession(ctx context.Context, connection *spookytypes.Connection) (*spookytypes.Session, error)
    
    // RunCommand runs a command via SSH
    RunCommand(ctx context.Context, session *spookytypes.Session, command *spookytypes.SSHCommand) (*spookytypes.SSHCommandResult, error)
    
    // CreateActingSession creates a new SSH acting session
    CreateActingSession(ctx context.Context, connection *spookytypes.Connection) (*spookytypes.ActingSession, error)
    
    // TransferFile transfers a file via SSH
    TransferFile(ctx context.Context, session *spookytypes.Session, transfer *spookytypes.FileTransfer) (*spookytypes.FileTransferResult, error)
    
    // ValidateConnection validates SSH connection parameters
    ValidateConnection(ctx context.Context, request *spookytypes.ConnectionRequest) (*spookytypes.ValidationResult, error)
    
    // ValidateAuthentication validates SSH authentication parameters
    ValidateAuthentication(ctx context.Context, auth *spookytypes.Authentication) (*spookytypes.ValidationResult, error)
    
    // GetConnectionPool returns the connection pool
    GetConnectionPool() *spookytypes.ConnectionPool
    
    // Close closes all SSH connections
    Close(ctx context.Context) error
}
```

**Methods:**

#### CreateClient
```go
CreateClient(ctx context.Context, config *spookytypes.ClientConfig) (*spookytypes.Client, error)
```
Creates a new SSH client with the specified configuration.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `config`: SSH client configuration

**Returns:**
- `*spookytypes.Client`: SSH client instance
- `error`: Error if client creation fails

**Behavior:**
- Validates client configuration
- Initializes connection pooling
- Sets up logging and error handling
- Returns configured client instance

#### Connect
```go
Connect(ctx context.Context, request *spookytypes.ConnectionRequest) (*spookytypes.ConnectionResult, error)
```
Establishes an SSH connection to the specified host.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `request`: Connection request with host details and authentication

**Returns:**
- `*spookytypes.ConnectionResult`: Connection result with status and metrics
- `error`: Error if connection fails

**Behavior:**
- Validates connection request parameters
- Loads and validates authentication keys/certificates
- Establishes TCP connection to SSH port
- Performs SSH handshake and authentication
- Returns connection result with metrics

#### RunCommand
```go
RunCommand(ctx context.Context, session *spookytypes.Session, command *spookytypes.SSHCommand) (*spookytypes.SSHCommandResult, error)
```
Runs a command via SSH session.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `session`: SSH session for command running
- `command`: Command to run with configuration

**Returns:**
- `*spookytypes.SSHCommandResult`: Command run result
- `error`: Error if command run fails

**Behavior:**
- Sets up command environment and working directory
- Runs command with specified timeout
- Captures standard output and error streams
- Returns run result with metrics

### Client

The main SSH client implementation.

```go
type Client struct {
    config      *spookytypes.ClientConfig
    logger      spookytypeslogging.Logger
    connections map[string]*ssh.Client
    mu          sync.RWMutex
    closed      bool
}
```

**Key Methods:**

#### NewClient
```go
func NewClient(config *spookytypes.ClientConfig, logger spookytypeslogging.Logger) *Client
```
Creates a new SSH client with default configuration if none provided.

**Default Configuration:**
- Default Port: 22
- Default Timeout: 30 seconds
- Max Connections: 10
- Max Retry Attempts: 3
- Retry Delay: 5 seconds
- Idle Timeout: 300 seconds

#### loadPrivateKey
```go
func (c *Client) loadPrivateKey(keyPath, passphrase string) (ssh.Signer, error)
```
Loads and validates a private key from file.

**Key Validation:**
- Type validation (ED25519, ED25519-SK, RSA 4096-bit)
- Size validation for RSA keys (minimum 4096 bits)
- Format validation (valid SSH key format)
- Access validation (readable file)

#### validateKeyType
```go
func (c *Client) validateKeyType(signer ssh.Signer) error
```
Validates that the key is of a supported type.

**Supported Types:**
- `ssh.KeyAlgoED25519`: ED25519 keys (always valid)
- `ssh.KeyAlgoRSA`: RSA keys (4096-bit minimum)
- Other types: Rejected with detailed error message

## Type Definitions

### Connection Types

#### Connection
```go
type Connection struct {
    spookytypescommon.CompleteEntity

    // Connection details
    Host     string `json:"host" hcl:"host"`
    Port     int    `json:"port" hcl:"port"`
    User     string `json:"user" hcl:"user"`
    Protocol string `json:"protocol" hcl:"protocol" default:"ssh"`

    // Connection state
    Status      ConnectionStatus `json:"status" hcl:"status"`
    ConnectedAt *time.Time       `json:"connected_at,omitempty" hcl:"connected_at,optional"`
    LastUsed    *time.Time       `json:"last_used,omitempty" hcl:"last_used,optional"`

    // Connection metrics
    Latency      time.Duration `json:"latency,omitempty" hcl:"latency,optional"`
    ErrorCount   int           `json:"error_count,omitempty" hcl:"error_count,optional"`
    SuccessCount int           `json:"success_count,omitempty" hcl:"success_count,optional"`

    // Connection configuration
    Timeout           time.Duration `json:"timeout,omitempty" hcl:"timeout,optional"`
    KeepaliveInterval time.Duration `json:"keepalive_interval,omitempty" hcl:"keepalive_interval,optional"`
    KeepaliveCount    int           `json:"keepalive_count,omitempty" hcl:"keepalive_count,optional"`

    // Authentication information
    AuthMethod     AuthMethod `json:"auth_method" hcl:"auth_method"`
    KeyPath        string     `json:"key_path,omitempty" hcl:"key_path,optional"`
    KeyFingerprint string     `json:"key_fingerprint,omitempty" hcl:"key_fingerprint,optional"`

    // Host key verification
    HostKeyFingerprint string `json:"host_key_fingerprint,omitempty" hcl:"host_key_fingerprint,optional"`
    KnownHostsPath     string `json:"known_hosts_path,omitempty" hcl:"known_hosts_path,optional"`
    StrictHostKeyCheck bool   `json:"strict_host_key_check" hcl:"strict_host_key_check" default:"true"`

    // Connection metadata
    ClientVersion string `json:"client_version,omitempty" hcl:"client_version,optional"`
    ServerVersion string `json:"server_version,omitempty" hcl:"server_version,optional"`
    Compression   bool   `json:"compression" hcl:"compression" default:"false"`
}
```

#### ConnectionRequest
```go
type ConnectionRequest struct {
    // Connection details
    Host string `json:"host" hcl:"host"`
    Port int    `json:"port" hcl:"port" default:"22"`
    User string `json:"user" hcl:"user"`

    // Authentication
    Password        string     `json:"password,omitempty" hcl:"password,optional" sensitive:"true"`
    KeyPath         string     `json:"key_path,omitempty" hcl:"key_path,optional"`
    CertificatePath string     `json:"certificate_path,omitempty" hcl:"certificate_path,optional"`
    Passphrase      string     `json:"passphrase,omitempty" hcl:"passphrase,optional" sensitive:"true"`
    AuthMethod      AuthMethod `json:"auth_method" hcl:"auth_method" default:"public_key"`

    // Connection settings
    Timeout           time.Duration `json:"timeout,omitempty" hcl:"timeout,optional"`
    KeepaliveInterval time.Duration `json:"keepalive_interval,omitempty" hcl:"keepalive_interval,optional"`
    KeepaliveCount    int           `json:"keepalive_count,omitempty" hcl:"keepalive_count,optional"`

    // Security settings
    KnownHostsPath     string `json:"known_hosts_path,omitempty" hcl:"known_hosts_path,optional"`
    StrictHostKeyCheck bool   `json:"strict_host_key_check" hcl:"strict_host_key_check" default:"true"`
    AllowInsecureHosts bool   `json:"allow_insecure_hosts" hcl:"allow_insecure_hosts" default:"false"`

    // Performance settings
    EnableCompression bool `json:"enable_compression" hcl:"enable_compression" default:"false"`
    EnableKeepalive   bool `json:"enable_keepalive" hcl:"enable_keepalive" default:"true"`

    // Request metadata
    RequestID   string    `json:"request_id" hcl:"request_id"`
    RequestedAt time.Time `json:"requested_at" hcl:"requested_at"`
    Priority    int       `json:"priority" hcl:"priority" default:"0"`
}
```

#### ConnectionResult
```go
type ConnectionResult struct {
    spookytypescommon.CompleteEntity

    // Connection details
    Connection *Connection        `json:"connection" hcl:"connection"`
    Request    *ConnectionRequest `json:"request" hcl:"request"`

    // Result status
    Success bool   `json:"success" hcl:"success"`
    Error   string `json:"error,omitempty" hcl:"error,optional"`

    // Connection metrics
    ConnectTime   time.Duration `json:"connect_time" hcl:"connect_time"`
    Latency       time.Duration `json:"latency" hcl:"latency"`
    RetryAttempts int           `json:"retry_attempts" hcl:"retry_attempts"`

    // Connection information
    ClientVersion      string `json:"client_version,omitempty" hcl:"client_version,optional"`
    ServerVersion      string `json:"server_version,omitempty" hcl:"server_version,optional"`
    HostKeyFingerprint string `json:"host_key_fingerprint,omitempty" hcl:"host_key_fingerprint,optional"`

    // Result metadata
    CompletedAt time.Time `json:"completed_at" hcl:"completed_at"`
}
```

### Client Types

#### Client
```go
type Client struct {
    spookytypescommon.CompleteEntity

    // Client configuration
    Config *ClientConfig `json:"config" hcl:"config"`

    // Client state
    Status    ClientStatus `json:"status" hcl:"status"`
    CreatedAt time.Time    `json:"created_at" hcl:"created_at"`
    LastUsed  *time.Time   `json:"last_used,omitempty" hcl:"last_used,optional"`

    // Client metrics
    TotalConnections int           `json:"total_connections" hcl:"total_connections"`
    ActiveConnections int          `json:"active_connections" hcl:"active_connections"`
    TotalCommands    int           `json:"total_commands" hcl:"total_commands"`
    AverageLatency   time.Duration `json:"average_latency" hcl:"average_latency"`

    // Connection pool
    ConnectionPool *ConnectionPool `json:"connection_pool" hcl:"connection_pool"`
}
```

#### ClientConfig
```go
type ClientConfig struct {
    // Default connection settings
    DefaultPort        int           `json:"default_port" hcl:"default_port" default:"22"`
    DefaultTimeout     time.Duration `json:"default_timeout" hcl:"default_timeout" default:"30s"`
    DefaultKeyPath     string        `json:"default_key_path" hcl:"default_key_path"`
    DefaultAuthMethod  AuthMethod    `json:"default_auth_method" hcl:"default_auth_method" default:"public_key"`

    // Connection pool settings
    MaxConnections     int           `json:"max_connections" hcl:"max_connections" default:"10"`
    MaxRetryAttempts   int           `json:"max_retry_attempts" hcl:"max_retry_attempts" default:"3"`
    RetryDelay         time.Duration `json:"retry_delay" hcl:"retry_delay" default:"5s"`
    IdleTimeout        time.Duration `json:"idle_timeout" hcl:"idle_timeout" default:"300s"`

    // Security settings
    KnownHostsPath     string `json:"known_hosts_path" hcl:"known_hosts_path"`
    StrictHostKeyCheck bool   `json:"strict_host_key_check" hcl:"strict_host_key_check" default:"true"`
    AllowInsecureHosts bool   `json:"allow_insecure_hosts" hcl:"allow_insecure_hosts" default:"false"`

    // Performance settings
    EnableCompression bool `json:"enable_compression" hcl:"enable_compression" default:"false"`
    EnableKeepalive   bool `json:"enable_keepalive" hcl:"enable_keepalive" default:"true"`
}
```

### Session Types

#### Session
```go
type Session struct {
    spookytypescommon.CompleteEntity

    // Session identification
    SessionID  string `json:"session_id" hcl:"session_id"`
    Connection *Connection `json:"connection" hcl:"connection"`

    // Session state
    Status    SessionStatus `json:"status" hcl:"status"`
    StartedAt time.Time     `json:"started_at" hcl:"started_at"`
    EndedAt   *time.Time    `json:"ended_at,omitempty" hcl:"ended_at,optional"`

    // Session configuration
    PtyConfig *PtyConfig `json:"pty_config,omitempty" hcl:"pty_config,optional"`

    // Session metrics
    ActionsRun int           `json:"actions_run" hcl:"actions_run"`
    TotalRunTime time.Duration `json:"total_run_time" hcl:"total_run_time"`
    AverageCommandTime time.Duration `json:"average_command_time" hcl:"average_command_time"`
}
```

#### SSHCommand
```go
type SSHCommand struct {
    spookytypescommon.CompleteEntity

    // Command details
    Command      string   `json:"command" hcl:"command"`
    Args         []string `json:"args,omitempty" hcl:"args,optional"`
    WorkingDir   string   `json:"working_dir,omitempty" hcl:"working_dir,optional"`
    Environment  map[string]string `json:"environment,omitempty" hcl:"environment,optional"`

    // Input/Output configuration
    Stdin        string `json:"stdin,omitempty" hcl:"stdin,optional"`
    CaptureOutput bool  `json:"capture_output" hcl:"capture_output" default:"true"`

    // Run configuration
    Timeout      time.Duration `json:"timeout,omitempty" hcl:"timeout,optional"`
    Priority     int           `json:"priority" hcl:"priority" default:"0"`
    ScheduledAt  time.Time     `json:"scheduled_at" hcl:"scheduled_at"`

    // Command metadata
    Tags         []string `json:"tags,omitempty" hcl:"tags,optional"`
    Description  string   `json:"description,omitempty" hcl:"description,optional"`
}
```

#### SSHCommandResult
```go
type SSHCommandResult struct {
    spookytypescommon.CompleteEntity

    // Command details
    Command *SSHCommand `json:"command" hcl:"command"`
    Session *Session    `json:"session" hcl:"session"`

    // Run results
    Success   bool   `json:"success" hcl:"success"`
    ExitCode  int    `json:"exit_code" hcl:"exit_code"`
    Stdout    string `json:"stdout,omitempty" hcl:"stdout,optional"`
    Stderr    string `json:"stderr,omitempty" hcl:"stderr,optional"`
    Error     string `json:"error,omitempty" hcl:"error,optional"`

    // Run metrics
    StartTime time.Time     `json:"start_time" hcl:"start_time"`
    EndTime   time.Time     `json:"end_time" hcl:"end_time"`
    Duration  time.Duration `json:"duration" hcl:"duration"`
}
```

### Authentication Types

#### Authentication
```go
type Authentication struct {
    spookytypescommon.CompleteEntity

    // Authentication method
    Method AuthMethod `json:"method" hcl:"method"`

    // Key-based authentication
    KeyPath    string `json:"key_path,omitempty" hcl:"key_path,optional"`
    Passphrase string `json:"passphrase,omitempty" hcl:"passphrase,optional" sensitive:"true"`

    // Certificate authentication
    CertificatePath string `json:"certificate_path,omitempty" hcl:"certificate_path,optional"`

    // Password authentication
    Password string `json:"password,omitempty" hcl:"password,optional" sensitive:"true"`

    // Authentication metadata
    KeyType     KeyType `json:"key_type,omitempty" hcl:"key_type,optional"`
    KeyFingerprint string `json:"key_fingerprint,omitempty" hcl:"key_fingerprint,optional"`
}
```

#### AuthMethod
```go
type AuthMethod string

const (
    AuthMethodPassword  AuthMethod = "password"
    AuthMethodPublicKey AuthMethod = "public_key"
    AuthMethodKeyboard  AuthMethod = "keyboard_interactive"
    AuthMethodAgent     AuthMethod = "agent"
    AuthMethodNone      AuthMethod = "none"
)
```

#### KeyType
```go
type KeyType string

const (
    KeyTypeRSA     KeyType = "rsa"
    KeyTypeDSA     KeyType = "dsa"
    KeyTypeECDSA   KeyType = "ecdsa"
    KeyTypeED25519 KeyType = "ed25519"
)
```

## Implementation Details

### Key Validation Algorithm

The SSH system implements comprehensive key validation:

```go
func (c *SimpleClient) validateKeyType(signer ssh.Signer) error {
    pubKey := signer.PublicKey()
    keyType := pubKey.Type()

    switch keyType {
    case ssh.KeyAlgoED25519:
        return c.validateED25519Key(pubKey)
    case ssh.KeyAlgoRSA:
        return c.validateRSAKey(pubKey)
    default:
        return &KeyValidationError{
            KeyType: keyType,
            Reason: fmt.Sprintf("unsupported key type: %s. Supported types: %s, %s, %s", 
                keyType, KeyTypeED25519, KeyTypeED25519SK, KeyTypeRSA4096),
        }
    }
}
```

### RSA Key Size Validation

RSA keys must meet minimum size requirements:

```go
func (c *SimpleClient) validateRSAKey(pubKey ssh.PublicKey) error {
    if pubKey.Type() != ssh.KeyAlgoRSA {
        return fmt.Errorf("expected RSA key, got %s", pubKey.Type())
    }

    rsaPubKey, ok := pubKey.(ssh.CryptoPublicKey)
    if !ok {
        return fmt.Errorf("failed to extract RSA public key")
    }

    cryptoPubKey := rsaPubKey.CryptoPublicKey()
    rsaKey, ok := cryptoPubKey.(*rsa.PublicKey)
    if !ok {
        return fmt.Errorf("failed to cast to RSA public key")
    }

    if rsaKey.Size()*8 < MinRSAKeySize {
        return fmt.Errorf("RSA key size %d bits is less than minimum required %d bits", 
            rsaKey.Size()*8, MinRSAKeySize)
    }

    return nil
}
```

### Certificate Loading

SSH certificates require both certificate and private key:

```go
func (c *SimpleClient) loadSSHCertificate(certPath, keyPath, passphrase string) (ssh.Signer, error) {
    // Load certificate
    certData, err := os.ReadFile(certPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read certificate file: %w", err)
    }

    // Load private key
    keyData, err := os.ReadFile(keyPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read private key file: %w", err)
    }

    // Parse private key
    var signer ssh.Signer
    if passphrase != "" {
        signer, err = ssh.ParsePrivateKeyWithPassphrase(keyData, []byte(passphrase))
    } else {
        signer, err = ssh.ParsePrivateKey(keyData)
    }
    if err != nil {
        return nil, fmt.Errorf("failed to parse private key: %w", err)
    }

    // Parse certificate
    cert, err := ssh.ParsePublicKey(certData)
    if err != nil {
        return nil, fmt.Errorf("failed to parse SSH certificate: %w", err)
    }

    // Create certificate signer
    certSigner, err := ssh.NewCertSigner(cert.(*ssh.Certificate), signer)
    if err != nil {
        return nil, fmt.Errorf("failed to create certificate signer: %w", err)
    }

    return certSigner, nil
}
```

### Connection Pooling

The SSH system implements efficient connection pooling:

```go
type SimpleClient struct {
    config      *spookytypes.ClientConfig
    logger      spookytypeslogging.Logger
    connections map[string]*ssh.Client
    mu          sync.RWMutex
    closed      bool
}
```

**Pool Management:**
- **Connection Reuse**: Existing connections are reused when possible
- **Health Checks**: Connections are tested before reuse
- **Cleanup**: Dead connections are automatically removed
- **Thread Safety**: Pool operations are thread-safe

## Error Handling

### Error Types

#### SSHError
```go
type SSHError struct {
    spookytypescommon.CompleteEntity

    // Error details
    Type        SSHErrorType `json:"type" hcl:"type"`
    Message     string       `json:"message" hcl:"message"`
    Operation   string       `json:"operation" hcl:"operation"`
    Host        string       `json:"host,omitempty" hcl:"host,optional"`
    User        string       `json:"user,omitempty" hcl:"user,optional"`

    // Error context
    Timestamp   time.Time    `json:"timestamp" hcl:"timestamp"`
    Details     string       `json:"details,omitempty" hcl:"details,optional"`
    Recoverable bool         `json:"recoverable" hcl:"recoverable"`
}
```

#### SSHErrorType
```go
type SSHErrorType string

const (
    SSHErrorTypeConnection     SSHErrorType = "connection"
    SSHErrorTypeAuthentication SSHErrorType = "authentication"
    SSHErrorTypeAuthorization  SSHErrorType = "authorization"
    SSHErrorTypeTimeout        SSHErrorType = "timeout"
    SSHErrorTypeProtocol       SSHErrorType = "protocol"
    SSHErrorTypeHostKey        SSHErrorType = "host_key"
    SSHErrorTypeCommand        SSHErrorType = "command"
    SSHErrorTypeSession        SSHErrorType = "session"
    SSHErrorTypeFileTransfer   SSHErrorType = "file_transfer"
    SSHErrorTypeValidation     SSHErrorType = "validation"
    SSHErrorTypeConfiguration  SSHErrorType = "configuration"
    SSHErrorTypeUnknown        SSHErrorType = "unknown"
)
```

#### KeyValidationError
```go
type KeyValidationError struct {
    KeyType string
    Reason  string
}

func (e *KeyValidationError) Error() string {
    return fmt.Sprintf("key validation failed for %s: %s", e.KeyType, e.Reason)
}
```

### Error Handling Patterns

#### Connection Errors
```go
// Handle connection failures
if !result.Success {
    switch {
    case strings.Contains(result.Error, "timeout"):
        // Handle timeout errors
        return fmt.Errorf("connection timeout to %s: %w", request.Host, err)
    case strings.Contains(result.Error, "authentication"):
        // Handle authentication errors
        return fmt.Errorf("authentication failed for %s@%s: %w", request.User, request.Host, err)
    default:
        // Handle other connection errors
        return fmt.Errorf("connection failed to %s: %w", request.Host, err)
    }
}
```

#### Key Validation Errors
```go
// Handle key validation errors
if err := client.validateKeyType(signer); err != nil {
    if keyErr, ok := err.(*KeyValidationError); ok {
        switch keyErr.KeyType {
        case "rsa":
            return fmt.Errorf("RSA key validation failed: %s", keyErr.Reason)
        case "ed25519":
            return fmt.Errorf("ED25519 key validation failed: %s", keyErr.Reason)
        default:
            return fmt.Errorf("unsupported key type: %s", keyErr.KeyType)
        }
    }
    return fmt.Errorf("key validation failed: %w", err)
}
```

## Key Validation Rules

### Supported Key Types

1. **ED25519 Keys**
   - **Type**: `ssh.KeyAlgoED25519`
   - **Size**: Fixed 256 bits
   - **Validation**: Always valid (fixed size)
   - **Security**: High (modern elliptic curve)

2. **ED25519-SK Keys**
   - **Type**: `ssh.KeyAlgoED25519` with security key
   - **Size**: Fixed 256 bits
   - **Validation**: Hardware key validation (planned)
   - **Security**: Very high (hardware-backed)

3. **RSA Keys**
   - **Type**: `ssh.KeyAlgoRSA`
   - **Size**: Minimum 4096 bits
   - **Validation**: Size validation required
   - **Security**: High (with sufficient key size)

### Validation Rules

#### Type Validation
```go
// Only supported key types are accepted
supportedTypes := map[string]bool{
    ssh.KeyAlgoED25519: true,
    ssh.KeyAlgoRSA:     true,
}

if !supportedTypes[keyType] {
    return &KeyValidationError{
        KeyType: keyType,
        Reason:  "unsupported key type",
    }
}
```

#### Size Validation
```go
// RSA keys must be 4096-bit minimum
if keyType == ssh.KeyAlgoRSA {
    if rsaKey.Size()*8 < MinRSAKeySize {
        return &KeyValidationError{
            KeyType: KeyTypeRSA4096,
            Reason:  fmt.Sprintf("key size %d bits < minimum %d bits", 
                rsaKey.Size()*8, MinRSAKeySize),
        }
    }
}
```

#### Format Validation
```go
// Keys must be in valid SSH format
if err := ssh.ParsePrivateKey(keyData); err != nil {
    return &KeyValidationError{
        KeyType: "unknown",
        Reason:  fmt.Sprintf("invalid key format: %v", err),
    }
}
```

## Certificate Handling

### Certificate Types

#### SSH Certificate
```go
type Certificate struct {
    spookytypescommon.CompleteEntity

    // Certificate details
    Type        string    `json:"type" hcl:"type"`
    Serial      uint64    `json:"serial" hcl:"serial"`
    KeyID       string    `json:"key_id" hcl:"key_id"`
    Principals  []string  `json:"principals" hcl:"principals"`
    ValidAfter  time.Time `json:"valid_after" hcl:"valid_after"`
    ValidBefore time.Time `json:"valid_before" hcl:"valid_before"`

    // Certificate metadata
    SignatureKey string `json:"signature_key" hcl:"signature_key"`
    Extensions   map[string]string `json:"extensions,omitempty" hcl:"extensions,optional"`
}
```

### Certificate Validation

#### Format Validation
```go
// Certificate must be in valid SSH format
cert, err := ssh.ParsePublicKey(certData)
if err != nil {
    return fmt.Errorf("failed to parse SSH certificate: %w", err)
}
```

#### Private Key Requirement
```go
// Certificate must be accompanied by private key
if keyPath == "" {
    return fmt.Errorf("certificate requires private key")
}
```

#### Expiration Validation (Planned)
```go
// Check certificate expiration
if time.Now().After(cert.ValidBefore) {
    return fmt.Errorf("certificate expired at %s", cert.ValidBefore)
}
```

## File Transfer

### SFTP Transfer

SFTP (SSH File Transfer Protocol) provides secure file transfer capabilities with progress tracking and verification.

```go
// Create file transfer manager
ftm := ssh.NewFileTransferManager(client, logger)

// Configure file transfer
transfer := &spookytypesssh.FileTransfer{
    LocalPath:   "/local/file.txt",
    RemotePath:  "/remote/file.txt",
    Direction:   spookytypesssh.TransferDirectionUpload,
    Mode:        spookytypesssh.TransferModeSFTP,
    Verify:      true,
    Permissions: 0o644,
}

// Run transfer
result, err := ftm.TransferFile(ctx, connection, transfer)
if err != nil {
    log.Printf("Transfer failed: %v", err)
    return
}

log.Printf("Transfer successful: %d bytes in %v", 
    result.BytesTransferred, result.Duration)
```

### SCP Transfer

SCP (Secure Copy Protocol) provides efficient file transfer using SSH.

```go
// Configure SCP transfer
transfer := &spookytypesssh.FileTransfer{
    LocalPath:   "/local/file.txt",
    RemotePath:  "/remote/file.txt",
    Direction:   spookytypesssh.TransferDirectionUpload,
    Mode:        spookytypesssh.TransferModeSCP,
    Verify:      true,
}

// Run transfer
result, err := ftm.TransferFile(ctx, connection, transfer)
```

### Batch Transfer

Perform multiple file transfers concurrently.

```go
transfers := []*spookytypesssh.FileTransfer{
    {
        LocalPath:  "/local/file1.txt",
        RemotePath: "/remote/file1.txt",
        Direction:  spookytypesssh.TransferDirectionUpload,
        Mode:       spookytypesssh.TransferModeSFTP,
    },
    {
        LocalPath:  "/local/file2.txt",
        RemotePath: "/remote/file2.txt",
        Direction:  spookytypesssh.TransferDirectionUpload,
        Mode:       spookytypesssh.TransferModeSFTP,
    },
}

results, err := ftm.BatchTransfer(ctx, connection, transfers)
```



## Advanced Authentication

### Multi-Factor Authentication

Configure multiple authentication methods with fallback support.

```go
// Create advanced auth manager
aam := ssh.NewAdvancedAuthManager(logger)

// Configure multi-factor authentication
authConfig := &ssh.MultiFactorAuthConfig{
    PrimaryMethod:   spookytypesssh.AuthMethodPublicKey,
    PrimaryKey:      "~/.ssh/id_rsa",
    SecondaryMethod: spookytypesssh.AuthMethodPassword,
    SecondaryPass:   "password",
    TOTPSecret:      "your_totp_secret",
    TOTPAlgorithm:   "sha1",
    TOTPDigits:      6,
    TOTPPeriod:      30,
    AuthOrder:       []string{"primary", "secondary", "totp"},
    MaxRetries:      3,
    RetryDelay:      5 * time.Second,
}

// Get authentication methods
authMethods, err := aam.GetAuthMethods(authConfig)
if err != nil {
    log.Printf("Failed to create auth methods: %v", err)
    return
}
```

### Certificate-Based Authentication

Generate and use SSH certificates for authentication.

```go
// Generate certificate configuration
certConfig := &ssh.CertificateConfig{
    KeyType:         "rsa",
    KeySize:         4096,
    Serial:          1,
    CertType:        1, // User certificate
    KeyID:           "user-cert",
    Principals:      []string{"user"},
    ValidAfter:      time.Now(),
    ValidBefore:     time.Now().Add(24 * time.Hour),
    CriticalOptions: map[string]string{},
    Extensions:      map[string]string{},
}

// Generate certificate
cert, err := aam.GenerateCertificate(certConfig)
if err != nil {
    log.Printf("Failed to generate certificate: %v", err)
    return
}

// Save certificate
err = aam.SaveCertificate(cert, "/path/to/cert.pub")
if err != nil {
    log.Printf("Failed to save certificate: %v", err)
    return
}
```

### SSH Agent Integration

Connect to and use the local SSH agent.

```go
// Connect to SSH agent
err := aam.agent.Connect()
if err != nil {
    log.Printf("Failed to connect to SSH agent: %v", err)
    return
}

// List available keys
keys, err := aam.agent.ListKeys()
if err != nil {
    log.Printf("Failed to list keys: %v", err)
    return
}

log.Printf("SSH agent has %d keys", len(keys))

// Add key to agent
err = aam.agent.AddKey(signer)
if err != nil {
    log.Printf("Failed to add key: %v", err)
    return
}
```

## Examples

### Basic SSH Client Usage

```go
package main

import (
    "context"
    "fmt"
    "time"

    spookylogging "spooky/internal/logging"
    spookytypes "spooky/internal/types"
    "spooky/internal/ssh"
)

func main() {
    // Create logger
    logManager := spookylogging.NewLogManager()
    logger := logManager.GetLogger("ssh-example")

    // Create SSH client
    config := &spookytypes.ClientConfig{
        DefaultPort:        22,
        DefaultTimeout:     30 * time.Second,
        MaxConnections:     5,
        MaxRetryAttempts:   3,
        RetryDelay:         5 * time.Second,
        IdleTimeout:        300 * time.Second,
    }

    client := ssh.NewClient(config, logger)

    // Create connection request
    request := &spookytypes.ConnectionRequest{
        Host:     "example.com",
        Port:     22,
        User:     "admin",
        KeyPath:  "~/.ssh/id_ed25519",
        Timeout:  30 * time.Second,
        RequestID: "example-connection",
        RequestedAt: time.Now(),
    }

    // Connect
    ctx := context.Background()
    result, err := client.Connect(ctx, request)
    if err != nil {
        log.Fatal(err)
    }

    if !result.Success {
        log.Fatalf("Connection failed: %s", result.Error)
    }

    fmt.Printf("Connected to %s:%d as %s\n", 
        result.Connection.Host, 
        result.Connection.Port, 
        result.Connection.User)

    // Run command
    command := &spookytypes.SSHCommand{
        Command:      "uname",
        Args:         []string{"-a"},
        CaptureOutput: true,
        Timeout:      10 * time.Second,
    }

    cmdResult, err := client.RunCommand(ctx, result.Connection, command)
    if err != nil {
        log.Fatal(err)
    }

    if cmdResult.Success {
        fmt.Printf("Command output: %s\n", cmdResult.Stdout)
    } else {
        fmt.Printf("Command failed: %s\n", cmdResult.Stderr)
    }

    // Close client
    client.Close(ctx)
}
```

### Key Validation Example

```go
package main

import (
    "fmt"
    "spooky/internal/ssh"
    "spooky/internal/logging"
)

func main() {
    // Create client for key validation
    logManager := logging.NewLogManager()
    logger := logManager.GetLogger("key-validation")
    client := ssh.NewClient(nil, logger)

    // Test key generation and validation
    keyTypes := []string{ssh.KeyTypeED25519, ssh.KeyTypeRSA4096}

    for _, keyType := range keyTypes {
        fmt.Printf("Testing %s key generation and validation...\n", keyType)
        
        signer, err := client.generateSupportedKey(keyType)
        if err != nil {
            fmt.Printf("Failed to generate %s key: %v\n", keyType, err)
            continue
        }

        // Validate the generated key
        if err := client.validateKeyType(signer); err != nil {
            fmt.Printf("Key validation failed for %s: %v\n", keyType, err)
        } else {
            fmt.Printf("✓ %s key validation passed\n", keyType)
        }
    }
}
```

### Certificate Authentication Example

```go
package main

import (
    "context"
    "fmt"
    "time"

    spookylogging "spooky/internal/logging"
    spookytypes "spooky/internal/types"
    "spooky/internal/ssh"
)

func main() {
    // Create logger and client
    logManager := spookylogging.NewLogManager()
    logger := logManager.GetLogger("cert-example")
    client := ssh.NewClient(nil, logger)

    // Create connection request with certificate
    request := &spookytypes.ConnectionRequest{
        Host:            "cert-server.example.com",
        Port:            22,
        User:            "admin",
        KeyPath:         "~/.ssh/id_ed25519",        // Private key
        CertificatePath: "~/.ssh/id_ed25519-cert.pub", // Certificate
        Passphrase:      "my-passphrase",            // Optional
        Timeout:         30 * time.Second,
        RequestID:       "certificate-connection",
        RequestedAt:     time.Now(),
    }

    // Connect using certificate authentication
    ctx := context.Background()
    result, err := client.Connect(ctx, request)
    if err != nil {
        log.Fatal(err)
    }

    if !result.Success {
        log.Fatalf("Certificate authentication failed: %s", result.Error)
    }

    fmt.Printf("Connected using certificate authentication\n")
    fmt.Printf("Certificate: %s\n", request.CertificatePath)
    fmt.Printf("Private key: %s\n", request.KeyPath)

    // Run command
    command := &spookytypes.SSHCommand{
        Command:      "whoami",
        CaptureOutput: true,
        Timeout:      10 * time.Second,
    }

    cmdResult, err := client.RunCommand(ctx, result.Connection, command)
    if err != nil {
        log.Fatal(err)
    }

    if cmdResult.Success {
        fmt.Printf("User: %s\n", cmdResult.Stdout)
    }

    client.Close(ctx)
}
```

This comprehensive API reference provides all the technical details needed to work with the SSH system, from basic usage to advanced features and error handling.
