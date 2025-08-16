# SSH System API Reference

## Overview

This document provides a comprehensive API reference for the spooky SSH system. It covers all interfaces, types, methods, and implementation details for developers working with the SSH system.

**Status: Implemented** - The SSH system provides comprehensive functionality for SSH connections, authentication, and command execution.

## Core Interfaces

### SSHManager Interface

The `SSHManager` interface provides the primary entry point for SSH operations:

```go
type SSHManager interface {
    // Connect establishes an SSH connection to the target host
    Connect(ctx context.Context, request *spookytypes.ConnectionRequest) (*spookytypes.ConnectionResult, error)
    
    // RunCommand executes a command on the remote host
    RunCommand(ctx context.Context, connection *spookytypes.Connection, command *spookytypes.SSHCommand) (*spookytypes.SSHCommandResult, error)
    
    // UploadFile uploads a file to the remote host
    UploadFile(ctx context.Context, connection *spookytypes.Connection, data []byte, remotePath string) error
    
    // DownloadFile downloads a file from the remote host
    DownloadFile(ctx context.Context, connection *spookytypes.Connection, remotePath string) ([]byte, error)
    
    // CloseConnection closes an SSH connection
    CloseConnection(ctx context.Context, connection *spookytypes.Connection) error
    
    // TestAuthentication tests SSH authentication without executing commands
    TestAuthentication(ctx context.Context, request *spookytypes.ConnectionRequest) (*spookytypes.AuthenticationResult, error)
}
```

**Implementation Status**: ✅ **Implemented** - Complete functionality for SSH connections and operations

## Core Types

### ConnectionRequest

```go
type ConnectionRequest struct {
    Host     string        `hcl:"host" json:"host"`
    Port     int           `hcl:"port" json:"port"`
    User     string        `hcl:"user" json:"user"`
    KeyPath  string        `hcl:"key_path,optional" json:"key_path,omitempty"`
    Password string        `hcl:"password,optional" json:"password,omitempty"`
    Timeout  time.Duration `hcl:"timeout,optional" json:"timeout,omitempty"`
    Config   *SSHConfig    `hcl:"config,block" json:"config,omitempty"`
}

type SSHConfig struct {
    HostKeyCallback   string            `hcl:"host_key_callback,optional" json:"host_key_callback,omitempty"`
    KnownHostsFile    string            `hcl:"known_hosts_file,optional" json:"known_hosts_file,omitempty"`
    StrictHostKeyChecking bool           `hcl:"strict_host_key_checking,optional" json:"strict_host_key_checking,omitempty"`
    Compression        bool              `hcl:"compression,optional" json:"compression,omitempty"`
    KeepAliveInterval time.Duration     `hcl:"keep_alive_interval,optional" json:"keep_alive_interval,omitempty"`
    KeepAliveCountMax int               `hcl:"keep_alive_count_max,optional" json:"keep_alive_count_max,omitempty"`
    Environment       map[string]string `hcl:"environment,optional" json:"environment,omitempty"`
}
```

### ConnectionResult

```go
type ConnectionResult struct {
    Success     bool                `hcl:"success" json:"success"`
    Connection  *spookytypes.Connection `hcl:"connection,block" json:"connection,omitempty"`
    Error       string              `hcl:"error,optional" json:"error,omitempty"`
    Duration    time.Duration       `hcl:"duration,optional" json:"duration,omitempty"`
    Metadata    *ConnectionMetadata `hcl:"metadata,block" json:"metadata,omitempty"`
}

type Connection struct {
    Host        string    `hcl:"host" json:"host"`
    Port        int       `hcl:"port" json:"port"`
    User        string    `hcl:"user" json:"user"`
    ConnectedAt time.Time `hcl:"connected_at" json:"connected_at"`
    SessionID   string    `hcl:"session_id" json:"session_id"`
}

type ConnectionMetadata struct {
    ServerVersion    string            `hcl:"server_version" json:"server_version"`
    ClientVersion    string            `hcl:"client_version" json:"client_version"`
    CipherSuite      string            `hcl:"cipher_suite" json:"cipher_suite"`
    KeyExchange      string            `hcl:"key_exchange" json:"key_exchange"`
    MAC              string            `hcl:"mac" json:"mac"`
    Compression      string            `hcl:"compression" json:"compression"`
    RemoteAddress    string            `hcl:"remote_address" json:"remote_address"`
    LocalAddress     string            `hcl:"local_address" json:"local_address"`
}
```

### SSHCommand

```go
type SSHCommand struct {
    Command     string            `hcl:"command" json:"command"`
    WorkingDir  string            `hcl:"working_dir,optional" json:"working_dir,omitempty"`
    Environment map[string]string `hcl:"environment,optional" json:"environment,omitempty"`
    Timeout     time.Duration     `hcl:"timeout,optional" json:"timeout,omitempty"`
    User        string            `hcl:"user,optional" json:"user,omitempty"`
    Stdin       string            `hcl:"stdin,optional" json:"stdin,omitempty"`
}

type SSHCommandResult struct {
    Success   bool          `hcl:"success" json:"success"`
    Output    string        `hcl:"output" json:"output"`
    Error     string        `hcl:"error,optional" json:"error,omitempty"`
    ExitCode  int           `hcl:"exit_code" json:"exit_code"`
    Duration  time.Duration `hcl:"duration,optional" json:"duration,omitempty"`
    Stderr    string        `hcl:"stderr,optional" json:"stderr,omitempty"`
}
```

### AuthenticationResult

```go
type AuthenticationResult struct {
    Success     bool      `hcl:"success" json:"success"`
    Method      string    `hcl:"method" json:"method"`
    Error       string    `hcl:"error,optional" json:"error,omitempty"`
    Duration    time.Duration `hcl:"duration,optional" json:"duration,omitempty"`
    TestedAt    time.Time `hcl:"tested_at" json:"tested_at"`
}
```

## Current Implementation Status

### ✅ Working Components

1. **SSH Connections**: SSH connection establishment and management
2. **Authentication**: SSH key and password authentication
3. **Command Execution**: Remote command execution with full output capture
4. **File Transfer**: File upload and download capabilities
5. **Connection Pooling**: Efficient connection pooling and reuse
6. **Host Key Validation**: Host key validation and management
7. **Error Handling**: Comprehensive error handling and reporting
8. **Timeout Management**: Configurable timeouts for all operations
9. **Environment Support**: Environment variable support for commands
10. **Session Management**: SSH session management and cleanup

### 🔧 Key Features

1. **Multiple Authentication Methods**: SSH key and password authentication
2. **Connection Pooling**: Efficient connection pooling for performance
3. **Host Key Validation**: Configurable host key validation
4. **File Transfer**: Secure file upload and download
5. **Command Execution**: Full command execution with output capture
6. **Timeout Management**: Configurable timeouts for all operations
7. **Error Recovery**: Automatic error recovery and retry mechanisms

## Implementation Details

### SSH Connection Management

The SSH system manages connections efficiently:

```go
// Connect to SSH host
func (m *Manager) Connect(ctx context.Context, request *spookytypes.ConnectionRequest) (*spookytypes.ConnectionResult, error) {
    start := time.Now()
    
    // Validate request
    if err := m.validateConnectionRequest(request); err != nil {
        return &spookytypes.ConnectionResult{
            Success: false,
            Error:   err.Error(),
        }, nil
    }
    
    // Create SSH client config
    config, err := m.createSSHConfig(request)
    if err != nil {
        return &spookytypes.ConnectionResult{
            Success: false,
            Error:   fmt.Sprintf("failed to create SSH config: %v", err),
        }, nil
    }
    
    // Establish connection
    client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", request.Host, request.Port), config)
    if err != nil {
        return &spookytypes.ConnectionResult{
            Success: false,
            Error:   fmt.Sprintf("failed to connect to %s:%d: %v", request.Host, request.Port, err),
        }, nil
    }
    
    // Create session
    session, err := client.NewSession()
    if err != nil {
        client.Close()
        return &spookytypes.ConnectionResult{
            Success: false,
            Error:   fmt.Sprintf("failed to create SSH session: %v", err),
        }, nil
    }
    
    // Create connection object
    connection := &spookytypes.Connection{
        Host:        request.Host,
        Port:        request.Port,
        User:        request.User,
        ConnectedAt: time.Now(),
        SessionID:   generateSessionID(),
    }
    
    // Store connection in pool
    m.connectionPool[connection.SessionID] = &SSHConnection{
        Client:  client,
        Session: session,
    }
    
    return &spookytypes.ConnectionResult{
        Success:    true,
        Connection: connection,
        Duration:   time.Since(start),
    }, nil
}
```

### SSH Command Execution

```go
// Execute command on remote host
func (m *Manager) RunCommand(ctx context.Context, connection *spookytypes.Connection, command *spookytypes.SSHCommand) (*spookytypes.SSHCommandResult, error) {
    start := time.Now()
    
    // Get connection from pool
    sshConn, exists := m.connectionPool[connection.SessionID]
    if !exists {
        return &spookytypes.SSHCommandResult{
            Success:  false,
            Error:    "connection not found in pool",
            ExitCode: -1,
        }, nil
    }
    
    // Create new session for command execution
    session, err := sshConn.Client.NewSession()
    if err != nil {
        return &spookytypes.SSHCommandResult{
            Success:  false,
            Error:    fmt.Sprintf("failed to create session: %v", err),
            ExitCode: -1,
        }, nil
    }
    defer session.Close()
    
    // Set up command environment
    if command.WorkingDir != "" {
        session.RequestPty("xterm", 80, 40, ssh.TerminalModes{})
        session.Stdin = strings.NewReader(fmt.Sprintf("cd %s && %s", command.WorkingDir, command.Command))
    } else {
        session.Stdin = strings.NewReader(command.Command)
    }
    
    // Set environment variables
    if len(command.Environment) > 0 {
        for key, value := range command.Environment {
            session.Setenv(key, value)
        }
    }
    
    // Capture output
    var stdout, stderr bytes.Buffer
    session.Stdout = &stdout
    session.Stderr = &stderr
    
    // Execute command with timeout
    if command.Timeout > 0 {
        ctx, cancel := context.WithTimeout(ctx, command.Timeout)
        defer cancel()
        
        done := make(chan error, 1)
        go func() {
            done <- session.Run(command.Command)
        }()
        
        select {
        case err := <-done:
            if err != nil {
                return &spookytypes.SSHCommandResult{
                    Success:  false,
                    Error:    err.Error(),
                    Output:   stdout.String(),
                    Stderr:   stderr.String(),
                    ExitCode: -1,
                    Duration: time.Since(start),
                }, nil
            }
        case <-ctx.Done():
            session.Signal(ssh.SIGKILL)
            return &spookytypes.SSHCommandResult{
                Success:  false,
                Error:    "command execution timed out",
                Output:   stdout.String(),
                Stderr:   stderr.String(),
                ExitCode: -1,
                Duration: time.Since(start),
            }, nil
        }
    } else {
        if err := session.Run(command.Command); err != nil {
            return &spookytypes.SSHCommandResult{
                Success:  false,
                Error:    err.Error(),
                Output:   stdout.String(),
                Stderr:   stderr.String(),
                ExitCode: -1,
                Duration: time.Since(start),
            }, nil
        }
    }
    
    return &spookytypes.SSHCommandResult{
        Success:  true,
        Output:   stdout.String(),
        Stderr:   stderr.String(),
        ExitCode: 0,
        Duration: time.Since(start),
    }, nil
}
```

### File Transfer Operations

```go
// Upload file to remote host
func (m *Manager) UploadFile(ctx context.Context, connection *spookytypes.Connection, data []byte, remotePath string) error {
    // Get connection from pool
    sshConn, exists := m.connectionPool[connection.SessionID]
    if !exists {
        return fmt.Errorf("connection not found in pool")
    }
    
    // Create SFTP client
    sftpClient, err := sftp.NewClient(sshConn.Client)
    if err != nil {
        return fmt.Errorf("failed to create SFTP client: %w", err)
    }
    defer sftpClient.Close()
    
    // Create remote file
    remoteFile, err := sftpClient.Create(remotePath)
    if err != nil {
        return fmt.Errorf("failed to create remote file %s: %w", remotePath, err)
    }
    defer remoteFile.Close()
    
    // Write data to remote file
    if _, err := remoteFile.Write(data); err != nil {
        return fmt.Errorf("failed to write to remote file %s: %w", remotePath, err)
    }
    
    return nil
}

// Download file from remote host
func (m *Manager) DownloadFile(ctx context.Context, connection *spookytypes.Connection, remotePath string) ([]byte, error) {
    // Get connection from pool
    sshConn, exists := m.connectionPool[connection.SessionID]
    if !exists {
        return nil, fmt.Errorf("connection not found in pool")
    }
    
    // Create SFTP client
    sftpClient, err := sftp.NewClient(sshConn.Client)
    if err != nil {
        return nil, fmt.Errorf("failed to create SFTP client: %w", err)
    }
    defer sftpClient.Close()
    
    // Open remote file
    remoteFile, err := sftpClient.Open(remotePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open remote file %s: %w", remotePath, err)
    }
    defer remoteFile.Close()
    
    // Read file data
    data, err := io.ReadAll(remoteFile)
    if err != nil {
        return nil, fmt.Errorf("failed to read remote file %s: %w", remotePath, err)
    }
    
    return data, nil
}
```

### Authentication Testing

```go
// Test SSH authentication
func (m *Manager) TestAuthentication(ctx context.Context, request *spookytypes.ConnectionRequest) (*spookytypes.AuthenticationResult, error) {
    start := time.Now()
    
    // Create SSH client config
    config, err := m.createSSHConfig(request)
    if err != nil {
        return &spookytypes.AuthenticationResult{
            Success: false,
            Error:   fmt.Sprintf("failed to create SSH config: %v", err),
        }, nil
    }
    
    // Attempt connection
    client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", request.Host, request.Port), config)
    if err != nil {
        return &spookytypes.AuthenticationResult{
            Success: false,
            Error:   fmt.Sprintf("authentication failed: %v", err),
        }, nil
    }
    defer client.Close()
    
    // Determine authentication method
    authMethod := "unknown"
    if request.KeyPath != "" {
        authMethod = "public_key"
    } else if request.Password != "" {
        authMethod = "password"
    }
    
    return &spookytypes.AuthenticationResult{
        Success:  true,
        Method:   authMethod,
        Duration: time.Since(start),
        TestedAt: time.Now(),
    }, nil
}
```

## Usage Examples

### Basic SSH Connection

```go
// Create SSH manager
sshManager := NewSSHManager()

// Create connection request
request := &spookytypes.ConnectionRequest{
    Host:    "192.168.1.100",
    Port:    22,
    User:    "admin",
    KeyPath: "~/.ssh/id_rsa",
    Timeout: 30 * time.Second,
}

// Connect to remote host
result, err := sshManager.Connect(ctx, request)
if err != nil {
    return fmt.Errorf("failed to connect: %w", err)
}

if !result.Success {
    return fmt.Errorf("connection failed: %s", result.Error)
}

log.Printf("Successfully connected to %s", result.Connection.Host)
```

### Command Execution

```go
// Execute command
command := &spookytypes.SSHCommand{
    Command: "ls -la /var/log",
    Timeout: 10 * time.Second,
}

cmdResult, err := sshManager.RunCommand(ctx, result.Connection, command)
if err != nil {
    return fmt.Errorf("failed to execute command: %w", err)
}

if cmdResult.Success {
    log.Printf("Command output: %s", cmdResult.Output)
} else {
    log.Printf("Command failed: %s", cmdResult.Error)
    log.Printf("Stderr: %s", cmdResult.Stderr)
}
```

### File Transfer

```go
// Upload file
data := []byte("Hello, World!")
err = sshManager.UploadFile(ctx, result.Connection, data, "/tmp/test.txt")
if err != nil {
    return fmt.Errorf("failed to upload file: %w", err)
}

// Download file
downloadedData, err := sshManager.DownloadFile(ctx, result.Connection, "/tmp/test.txt")
if err != nil {
    return fmt.Errorf("failed to download file: %w", err)
}

log.Printf("Downloaded data: %s", string(downloadedData))
```

### Authentication Testing

```go
// Test authentication
authResult, err := sshManager.TestAuthentication(ctx, request)
if err != nil {
    return fmt.Errorf("failed to test authentication: %w", err)
}

if authResult.Success {
    log.Printf("Authentication successful using %s method", authResult.Method)
} else {
    log.Printf("Authentication failed: %s", authResult.Error)
}
```

### CLI Usage

```bash
# Test SSH connectivity
spooky ssh test --host 192.168.1.100 --user admin --key-file ~/.ssh/id_rsa

# Execute remote command
spooky ssh exec --host 192.168.1.100 --user admin --key-file ~/.ssh/id_rsa --command "ls -la"

# Upload file
spooky ssh upload --host 192.168.1.100 --user admin --key-file ~/.ssh/id_rsa --local-file ./config.txt --remote-file /tmp/config.txt

# Download file
spooky ssh download --host 192.168.1.100 --user admin --key-file ~/.ssh/id_rsa --remote-file /var/log/messages --local-file ./messages.log
```

## Error Handling

### Connection Errors

```go
// Handle connection errors
result, err := sshManager.Connect(ctx, request)
if err != nil {
    return fmt.Errorf("connection error: %w", err)
}

if !result.Success {
    // Provide specific guidance based on error type
    if strings.Contains(result.Error, "connection refused") {
        return fmt.Errorf("SSH service not running on %s:%d", request.Host, request.Port)
    }
    
    if strings.Contains(result.Error, "authentication failed") {
        return fmt.Errorf("authentication failed - check credentials")
    }
    
    if strings.Contains(result.Error, "host key verification failed") {
        return fmt.Errorf("host key verification failed - check known_hosts")
    }
    
    return fmt.Errorf("connection failed: %s", result.Error)
}
```

### Command Execution Errors

```go
// Handle command execution errors
cmdResult, err := sshManager.RunCommand(ctx, connection, command)
if err != nil {
    return fmt.Errorf("command execution error: %w", err)
}

if !cmdResult.Success {
    // Check exit code for specific error types
    switch cmdResult.ExitCode {
    case 127:
        return fmt.Errorf("command not found: %s", command.Command)
    case 126:
        return fmt.Errorf("command not executable: %s", command.Command)
    case 13:
        return fmt.Errorf("permission denied: %s", command.Command)
    default:
        return fmt.Errorf("command failed with exit code %d: %s", cmdResult.ExitCode, cmdResult.Error)
    }
}
```

## Testing

### SSH Connection Testing

```go
func TestSSHConnection(t *testing.T) {
    // Create SSH manager
    manager := NewSSHManager()
    
    // Test connection request
    request := &spookytypes.ConnectionRequest{
        Host:    "test-host",
        Port:    22,
        User:    "test-user",
        KeyPath: "test-key",
        Timeout: 10 * time.Second,
    }
    
    // Test connection
    result, err := manager.Connect(ctx, request)
    if err != nil {
        t.Fatalf("Failed to connect: %v", err)
    }
    
    // Validate result
    if !result.Success {
        t.Errorf("Expected successful connection, got error: %s", result.Error)
    }
    
    if result.Connection == nil {
        t.Error("Expected connection object, got nil")
    }
    
    if result.Connection.Host != request.Host {
        t.Errorf("Expected host %s, got %s", request.Host, result.Connection.Host)
    }
}
```

### Command Execution Testing

```go
func TestSSHCommandExecution(t *testing.T) {
    // Create SSH manager
    manager := NewSSHManager()
    
    // Create test connection
    connection := &spookytypes.Connection{
        Host:        "test-host",
        Port:        22,
        User:        "test-user",
        ConnectedAt: time.Now(),
        SessionID:   "test-session",
    }
    
    // Test command
    command := &spookytypes.SSHCommand{
        Command: "echo 'hello world'",
        Timeout: 5 * time.Second,
    }
    
    // Execute command
    result, err := manager.RunCommand(ctx, connection, command)
    if err != nil {
        t.Fatalf("Failed to execute command: %v", err)
    }
    
    // Validate result
    if !result.Success {
        t.Errorf("Expected successful command execution, got error: %s", result.Error)
    }
    
    if result.ExitCode != 0 {
        t.Errorf("Expected exit code 0, got %d", result.ExitCode)
    }
    
    expectedOutput := "hello world\n"
    if result.Output != expectedOutput {
        t.Errorf("Expected output '%s', got '%s'", expectedOutput, result.Output)
    }
}
```

## Best Practices

### SSH Configuration

1. **Use SSH Keys**: Prefer SSH key authentication over passwords
2. **Validate Host Keys**: Always validate host keys for security
3. **Use Timeouts**: Set appropriate timeouts for all operations
4. **Handle Errors**: Implement proper error handling for all SSH operations
5. **Clean Up Connections**: Always close SSH connections when done

### Performance Optimization

```go
// Use connection pooling for multiple operations
func executeMultipleCommands(sshManager SSHManager, connection *spookytypes.Connection, commands []string) error {
    for _, cmd := range commands {
        command := &spookytypes.SSHCommand{
            Command: cmd,
            Timeout: 10 * time.Second,
        }
        
        result, err := sshManager.RunCommand(ctx, connection, command)
        if err != nil {
            return fmt.Errorf("failed to execute command '%s': %w", cmd, err)
        }
        
        if !result.Success {
            return fmt.Errorf("command '%s' failed: %s", cmd, result.Error)
        }
        
        log.Printf("Command '%s' completed successfully", cmd)
    }
    
    return nil
}
```

## Future Enhancements

### Planned Features

1. **Connection Multiplexing**: SSH connection multiplexing for improved performance
2. **Advanced Authentication**: Support for additional authentication methods
3. **SSH Tunneling**: SSH tunneling and port forwarding
4. **Connection Monitoring**: Real-time connection monitoring and health checks
5. **SSH Agent Support**: SSH agent integration for key management
6. **Connection Encryption**: Additional encryption options and algorithms

### Architecture Improvements

1. **Distributed SSH**: Distributed SSH management across multiple controllers
2. **SSH Caching**: SSH connection caching for improved performance
3. **SSH Streaming**: Streaming SSH operations for real-time data
4. **SSH Compression**: SSH compression for improved performance
5. **SSH Replication**: SSH connection replication for high availability

## Related Documentation

- [SSH User Guide](SSH_USER_GUIDE.md) - User guide for SSH system
- [SSH Troubleshooting](SSH_TROUBLESHOOTING.md) - Troubleshooting guide
- [Machines API Reference](MACHINES_API_REFERENCE.md) - Machines system API reference
- [Actions API Reference](ACTIONS_API_REFERENCE.md) - Actions system API reference
