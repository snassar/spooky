# SSH API Reference

## Overview

The SSH system provides comprehensive functionality for SSH connections, authentication, and command execution in the spooky codebase.

**Status: Implemented** - The SSH system provides comprehensive functionality for SSH connections, authentication, and command execution.

## Core Components

### SSH Client (`internal/ssh/client.go`)

The SSH client provides the foundation for all SSH operations:

```go
type Client struct {
    config              *spookytypes.ClientConfig
    logger              spookytypeslogging.Logger
    hostKeyManager      *HostKeyManager
    connectionPool      *AdvancedConnectionPool
    fileTransferManager *FileTransferManager
    advancedAuthManager *AdvancedAuthManager
    closed              bool
}
```

**Key Features:**
- ✅ **Connection Management**: Advanced connection pooling with health monitoring
- ✅ **Authentication**: Support for SSH keys, passwords, and agent authentication
- ✅ **Host Key Validation**: Configurable host key checking with known hosts support
- ✅ **File Transfer**: SFTP and SCP file transfer capabilities
- ✅ **Advanced Authentication**: Multi-factor authentication support

### SSH Manager (`internal/ssh/manager.go`)

The SSH manager provides high-level SSH operations:

```go
type Manager struct {
    client *Client
    logger spookytypeslogging.Logger
}
```

**Implemented Operations:**
- ✅ **Connect**: Establish SSH connections with authentication
- ✅ **CreateSession**: Create SSH sessions for command execution
- ✅ **RunCommand**: Execute commands on remote machines
- ✅ **TransferFile**: Transfer files using SFTP or SCP
- ✅ **Ping**: Test SSH connectivity with lightweight checks
- ✅ **ValidateConnection**: Validate connection parameters

### File Transfer Manager (`internal/ssh/file_transfer.go`)

Comprehensive file transfer capabilities:

```go
type FileTransferManager struct {
    client *Client
    logger spookytypeslogging.Logger
}
```

**Supported Operations:**
- ✅ **SFTP Transfer**: Upload and download files via SFTP
- ✅ **SCP Transfer**: Upload and download files via SCP
- ✅ **Progress Tracking**: Real-time transfer progress monitoring
- ✅ **File Verification**: Checksum verification for transferred files
- ✅ **Permission Management**: Set file permissions during transfer
- ✅ **Directory Creation**: Automatic remote directory creation

### Advanced Connection Pool (`internal/ssh/connection_pool.go`)

Efficient connection management:

```go
type AdvancedConnectionPool struct {
    connections map[string]*PooledConnection
    mutex       sync.RWMutex
    config      *spookytypes.ClientConfig
    logger      spookytypeslogging.Logger
}
```

**Features:**
- ✅ **Connection Reuse**: Reuse connections for improved performance
- ✅ **Health Monitoring**: Monitor connection health and cleanup stale connections
- ✅ **Load Balancing**: Distribute connections across multiple targets
- ✅ **Automatic Cleanup**: Clean up idle connections automatically
- ✅ **Metrics Tracking**: Track connection usage and performance

## API Methods

### Connection Management

#### Connect
```go
func (c *Client) Connect(ctx context.Context, request *spookytypes.ConnectionRequest) (*spookytypes.ConnectionResult, error)
```
Establishes an SSH connection with proper authentication.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `request`: Connection request with host, port, user, and authentication details

**Returns:**
- `ConnectionResult`: Connection result with success status and connection details
- `error`: Error if connection fails

#### CreateSession
```go
func (m *Manager) CreateSession(ctx context.Context, connection *spookytypes.Connection) (*spookytypes.Session, error)
```
Creates an SSH session for command execution.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `connection`: SSH connection to create session on

**Returns:**
- `Session`: SSH session for command execution
- `error`: Error if session creation fails

### Command Execution

#### RunCommand
```go
func (m *Manager) RunCommand(ctx context.Context, session *spookytypes.Session, command *spookytypes.SSHCommand) (*spookytypes.CommandResult, error)
```
Executes a command on a remote machine via SSH.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `session`: SSH session for command execution
- `command`: Command to execute with arguments, working directory, and environment

**Returns:**
- `CommandResult`: Command execution result with output and exit code
- `error`: Error if command execution fails

#### RunCommandWithStdin
```go
func (c *ReusableSSHClient) RunCommandWithStdin(ctx context.Context, machine *spookytypes.Machine, command string, stdin string) (string, error)
```
Executes a command with stdin input.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `machine`: Target machine configuration
- `command`: Command to execute
- `stdin`: Standard input data

**Returns:**
- `string`: Command output
- `error`: Error if command execution fails

### File Transfer

#### TransferFile
```go
func (ft *FileTransferManager) TransferFile(ctx context.Context, connection *spookytypes.Connection, transfer *spookytypesssh.FileTransfer) (*spookytypesssh.FileTransferResult, error)
```
Transfers files using SFTP or SCP.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `connection`: SSH connection for file transfer
- `transfer`: File transfer configuration with source, destination, and mode

**Returns:**
- `FileTransferResult`: Transfer result with bytes transferred and duration
- `error`: Error if transfer fails

**Supported Modes:**
- `TransferModeSFTP`: Use SFTP for file transfer
- `TransferModeSCP`: Use SCP for file transfer

**Supported Directions:**
- `TransferDirectionUpload`: Upload file from local to remote
- `TransferDirectionDownload`: Download file from remote to local

### Connectivity Testing

#### Ping
```go
func (m *Manager) Ping(ctx context.Context, request *spookytypes.PingRequest) (*spookytypes.PingResult, error)
```
Tests SSH connectivity to a remote machine.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `request`: Ping request with host and authentication details

**Returns:**
- `PingResult`: Ping result with response time and status
- `error`: Error if ping fails

**Features:**
- DNS resolution testing
- ICMP reachability testing (lightweight TCP-based)
- SSH connection testing
- Authentication validation

### Connection Validation

#### ValidateConnection
```go
func (m *Manager) ValidateConnection(ctx context.Context, request *spookytypes.ConnectionRequest) (*spookytypes.ValidationResult, error)
```
Validates SSH connection parameters without establishing a connection.

**Parameters:**
- `ctx`: Context for cancellation and timeouts
- `request`: Connection request to validate

**Returns:**
- `ValidationResult`: Validation result with errors and warnings
- `error`: Error if validation fails

## Configuration

### Client Configuration
```go
type ClientConfig struct {
    DefaultPort       int           `json:"default_port" hcl:"default_port"`
    DefaultTimeout    time.Duration `json:"default_timeout" hcl:"default_timeout"`
    MaxConnections    int           `json:"max_connections" hcl:"max_connections"`
    MaxRetryAttempts  int           `json:"max_retry_attempts" hcl:"max_retry_attempts"`
    RetryDelay        time.Duration `json:"retry_delay" hcl:"retry_delay"`
    IdleTimeout       time.Duration `json:"idle_timeout" hcl:"idle_timeout"`
    KnownHostsPath    string        `json:"known_hosts_path" hcl:"known_hosts_path"`
    StrictHostKeyCheck bool         `json:"strict_host_key_check" hcl:"strict_host_key_check"`
    AllowInsecureHosts bool         `json:"allow_insecure_hosts" hcl:"allow_insecure_hosts"`
}
```

### Connection Request
```go
type ConnectionRequest struct {
    Host       string        `json:"host" hcl:"host"`
    Port       int           `json:"port" hcl:"port"`
    User       string        `json:"user" hcl:"user"`
    Password   string        `json:"password,omitempty" hcl:"password,optional"`
    KeyPath    string        `json:"key_path,omitempty" hcl:"key_path,optional"`
    Passphrase string        `json:"passphrase,omitempty" hcl:"passphrase,optional"`
    Timeout    time.Duration `json:"timeout" hcl:"timeout"`
}
```

### SSH Command
```go
type SSHCommand struct {
    Command       string            `json:"command" hcl:"command"`
    Args          []string          `json:"args,omitempty" hcl:"args,optional"`
    WorkingDir    string            `json:"working_dir,omitempty" hcl:"working_dir,optional"`
    Environment   map[string]string `json:"environment,omitempty" hcl:"environment,optional"`
    Timeout       time.Duration     `json:"timeout" hcl:"timeout"`
    CaptureOutput bool              `json:"capture_output" hcl:"capture_output"`
}
```

### File Transfer
```go
type FileTransfer struct {
    LocalPath   string `json:"local_path" hcl:"local_path"`
    RemotePath  string `json:"remote_path" hcl:"remote_path"`
    Direction   string `json:"direction" hcl:"direction"`
    Mode        string `json:"mode" hcl:"mode"`
    Permissions uint32 `json:"permissions,omitempty" hcl:"permissions,optional"`
    Verify      bool   `json:"verify" hcl:"verify"`
}
```

## Error Handling

The SSH system provides comprehensive error handling:

### Connection Errors
- Authentication failures
- Host key validation errors
- Network connectivity issues
- Timeout errors

### Command Execution Errors
- Command not found
- Permission denied
- Timeout errors
- Exit code errors

### File Transfer Errors
- File not found
- Permission denied
- Disk space issues
- Network errors

## Security Features

### Authentication Methods
- ✅ **SSH Key Authentication**: Support for Ed25519 and RSA keys
- ✅ **Password Authentication**: Secure password-based authentication
- ✅ **SSH Agent**: Integration with SSH agent for key management
- ✅ **Multi-factor Authentication**: Support for TOTP and other MFA methods

### Host Key Validation
- ✅ **Known Hosts**: Support for known_hosts file validation
- ✅ **Strict Mode**: Configurable strict host key checking
- ✅ **Insecure Mode**: Option to allow insecure hosts (development only)

### Security Best Practices
- ✅ **Connection Encryption**: All connections use SSH encryption
- ✅ **Key Validation**: SSH key format and permission validation
- ✅ **Timeout Management**: Configurable timeouts for all operations
- ✅ **Resource Cleanup**: Automatic cleanup of connections and sessions

## Performance Features

### Connection Pooling
- ✅ **Connection Reuse**: Reuse connections for improved performance
- ✅ **Health Monitoring**: Monitor connection health automatically
- ✅ **Load Balancing**: Distribute connections across targets
- ✅ **Automatic Cleanup**: Clean up idle connections

### File Transfer Optimization
- ✅ **Progress Tracking**: Real-time transfer progress monitoring
- ✅ **Parallel Transfers**: Support for parallel file transfers
- ✅ **Resume Capability**: Resume interrupted transfers
- ✅ **Compression**: Optional compression for large files

### Metrics and Monitoring
- ✅ **Connection Metrics**: Track connection usage and performance
- ✅ **Transfer Metrics**: Monitor file transfer performance
- ✅ **Error Tracking**: Track and report SSH errors
- ✅ **Performance Logging**: Log performance metrics for analysis

## Integration with Other Systems

### Actions Integration
The SSH system integrates with the actions system for command execution:

```go
// Run command action via SSH
err = m.runCommandAction(ctx, action, machine, result)
```

### Facts Integration
The SSH system integrates with the facts system for remote fact collection:

```go
// Collect facts via SSH
facts, err := sshCollector.CollectViaSSH(ctx, machine)
```

### File Transfer Integration
The SSH system provides file transfer capabilities for template deployment and file operations:

```go
// Transfer files via SSH
result, err := m.TransferFile(ctx, session, transfer)
```

## Usage Examples

### Basic Connection
```go
client := NewClient(config, logger)
defer client.Close(context.Background())

request := &spookytypes.ConnectionRequest{
    Host:     "example.com",
    Port:     22,
    User:     "user",
    KeyPath:  "~/.ssh/id_rsa",
    Timeout:  30 * time.Second,
}

result, err := client.Connect(context.Background(), request)
if err != nil {
    log.Fatal(err)
}
```

### Command Execution
```go
session, err := manager.CreateSession(ctx, connection)
if err != nil {
    return err
}

command := &spookytypes.SSHCommand{
    Command:       "ls -la",
    WorkingDir:    "/home/user",
    Timeout:       30 * time.Second,
    CaptureOutput: true,
}

result, err := manager.RunCommand(ctx, session, command)
if err != nil {
    return err
}

fmt.Printf("Output: %s\n", result.Stdout)
```

### File Transfer
```go
transfer := &spookytypesssh.FileTransfer{
    LocalPath:  "/local/file.txt",
    RemotePath: "/remote/file.txt",
    Direction:  spookytypesssh.TransferDirectionUpload,
    Mode:       spookytypesssh.TransferModeSFTP,
    Verify:     true,
}

result, err := fileTransferManager.TransferFile(ctx, connection, transfer)
if err != nil {
    return err
}

fmt.Printf("Transferred %d bytes in %v\n", result.BytesTransferred, result.Duration)
```

### Connectivity Testing
```go
request := &spookytypes.PingRequest{
    Host:     "example.com",
    Port:     22,
    User:     "user",
    KeyPath:  "~/.ssh/id_rsa",
    Timeout:  10 * time.Second,
}

result, err := manager.Ping(ctx, request)
if err != nil {
    return err
}

fmt.Printf("Ping successful: %v\n", result.ResponseTime)
```

## Testing

The SSH system includes comprehensive testing:

### Unit Tests
- ✅ **Client Tests**: Test SSH client functionality
- ✅ **Manager Tests**: Test SSH manager operations
- ✅ **File Transfer Tests**: Test file transfer capabilities
- ✅ **Connection Pool Tests**: Test connection pooling

### Integration Tests
- ✅ **Mock SSH Server**: Test with mock SSH server
- ✅ **Real SSH Server**: Test with real SSH server
- ✅ **Authentication Tests**: Test various authentication methods
- ✅ **File Transfer Tests**: Test file transfer with real server

### Performance Tests
- ✅ **Connection Pooling**: Test connection pool performance
- ✅ **File Transfer**: Test file transfer performance
- ✅ **Concurrent Operations**: Test concurrent SSH operations
- ✅ **Memory Usage**: Test memory usage under load

## Troubleshooting

### Common Issues

#### Authentication Failures
- Verify SSH key permissions (should be 600)
- Check SSH key format (Ed25519 or RSA 4096-bit)
- Verify passphrase for encrypted keys
- Check user permissions on remote machine

#### Connection Timeouts
- Verify network connectivity
- Check firewall settings
- Verify SSH service is running
- Check SSH configuration on remote machine

#### File Transfer Failures
- Verify file permissions
- Check disk space on remote machine
- Verify directory exists on remote machine
- Check SSH configuration for file transfer

### Debugging

#### Enable Debug Logging
```go
logger.SetLevel(spookytypeslogging.DebugLevel)
```

#### Connection Diagnostics
```go
result, err := manager.ValidateConnection(ctx, request)
if err != nil {
    fmt.Printf("Validation errors: %v\n", result.Errors)
}
```

#### Performance Monitoring
```go
stats, err := connectionPool.GetStats()
if err != nil {
    fmt.Printf("Connection pool stats: %v\n", stats)
}
```

## Future Enhancements

### Planned Features
- **SSH Key Management**: Automated SSH key generation and management
- **Advanced Authentication**: Support for additional authentication methods
- **Connection Multiplexing**: SSH connection multiplexing for improved performance
- **Advanced File Transfer**: Support for rsync and other transfer protocols

### Performance Improvements
- **Connection Optimization**: Further optimize connection pooling
- **Transfer Optimization**: Improve file transfer performance
- **Memory Optimization**: Reduce memory usage for large operations
- **Concurrency Optimization**: Improve concurrent operation handling
