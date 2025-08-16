# SSH System API Reference

## Overview

This document provides a comprehensive API reference for the spooky SSH system. It covers all interfaces, types, methods, and implementation details for developers working with the SSH system.

**Status: Partially Implemented** - The SSH system has basic functionality but connection management and authentication have known issues that need to be addressed.

## Core Interfaces

### SSHManager Interface

The `SSHManager` interface provides the primary entry point for SSH operations:

```go
type SSHManager interface {
    // GetConnection gets an SSH connection for the given machine
    GetConnection(hostname string, port int, user string) (*ssh.Client, error)

    // ReturnConnection returns an SSH connection to the pool
    ReturnConnection(conn *ssh.Client)

    // TestConnection tests SSH connectivity to a machine
    TestConnection(hostname string, port int, user string) error

    // ExecuteCommand executes a command on a remote machine
    ExecuteCommand(conn *ssh.Client, command string) (string, error)

    // CopyFile copies a file to a remote machine
    CopyFile(conn *ssh.Client, localPath string, remotePath string) error

    // GetFile retrieves a file from a remote machine
    GetFile(conn *ssh.Client, remotePath string, localPath string) error
}
```

**Implementation Status**: ⚠️ **Partially Implemented** - Basic functionality exists but connection management has issues

### SSHClient Interface

The `SSHClient` interface provides SSH client operations:

```go
type SSHClient interface {
    // Connect establishes an SSH connection
    Connect() (*ssh.Client, error)

    // ExecuteCommand executes a command on the remote machine
    ExecuteCommand(command string) (string, error)

    // CopyFile copies a file to the remote machine
    CopyFile(localPath string, remotePath string) error

    // GetFile retrieves a file from the remote machine
    GetFile(remotePath string, localPath string) error

    // Close closes the SSH connection
    Close() error
}
```

**Implementation Status**: ⚠️ **Partially Implemented** - Basic connection exists but operations have issues

## Current Implementation Status

### ✅ Working Components

1. **SSH Connection**: Basic SSH connection establishment
2. **SSH Configuration**: SSH configuration management
3. **SSH Structure**: Proper SSH type definitions and structures
4. **CLI Integration**: SSH connectivity testing via `spooky machines ping`
5. **Project Integration**: SSH configuration loading from project configuration
6. **Basic Validation**: SSH configuration validation and error handling
7. **Authentication Support**: Support for SSH key authentication
8. **Connection Pooling**: Basic SSH connection pooling
9. **Host Key Management**: Basic host key validation
10. **File Transfer**: Basic file transfer operations

### ⚠️ Known Issues

1. **Connection Management**: SSH connection management has implementation issues
2. **Authentication Testing**: SSH authentication testing has problems
3. **Connection Pooling**: SSH connection pooling has issues
4. **Host Key Validation**: Host key validation has implementation problems
5. **File Transfer**: File transfer operations have issues
6. **Parallel Processing**: No parallel SSH operations support

### 🔄 In Progress

1. **Connection Fixes**: Addressing SSH connection management issues
2. **Authentication Improvements**: Implementing proper SSH authentication testing
3. **Pooling Fixes**: Fixing SSH connection pooling

## Implementation Details

### SSH Manager System

The SSH system manages SSH connections and operations:

```go
type Manager struct {
    logger          spookytypeslogging.Logger
    config          *spookytypes.SSHConfig
    connectionPool  *ConnectionPool
    hostKeyManager  *HostKeyManager
}

func NewManager(
    logger spookytypeslogging.Logger,
    config *spookytypes.SSHConfig,
) spookyinterfaces.SSHManager {
    return &Manager{
        logger:         logger,
        config:         config,
        connectionPool: NewConnectionPool(config),
        hostKeyManager: NewHostKeyManager(config),
    }
}
```

### SSH Connection Implementation

```go
// GetConnection gets an SSH connection for the given machine
func (m *Manager) GetConnection(hostname string, port int, user string) (*ssh.Client, error) {
    m.logger.Info("Getting SSH connection", map[string]interface{}{
        "hostname": hostname,
        "port":     port,
        "user":     user,
    })

    // Check connection pool first
    if conn := m.connectionPool.GetConnection(hostname); conn != nil {
        m.logger.Debug("Reusing connection from pool", map[string]interface{}{
            "hostname": hostname,
        })
        return conn, nil
    }

    // Create new connection
    conn, err := m.createConnection(hostname, port, user)
    if err != nil {
        m.logger.Error("Failed to create SSH connection", err, map[string]interface{}{
            "hostname": hostname,
            "port":     port,
            "user":     user,
        })
        return nil, fmt.Errorf("failed to create SSH connection: %w", err)
    }

    // Add to connection pool
    m.connectionPool.AddConnection(hostname, conn)

    m.logger.Info("SSH connection established", map[string]interface{}{
        "hostname": hostname,
    })

    return conn, nil
}

// ReturnConnection returns an SSH connection to the pool
func (m *Manager) ReturnConnection(conn *ssh.Client) {
    if conn == nil {
        return
    }

    // Return connection to pool
    m.connectionPool.ReturnConnection(conn)

    m.logger.Debug("SSH connection returned to pool")
}

func (m *Manager) createConnection(hostname string, port int, user string) (*ssh.Client, error) {
    // Create SSH client configuration
    config, err := m.createSSHConfig(user)
    if err != nil {
        return nil, fmt.Errorf("failed to create SSH config: %w", err)
    }

    // Establish connection
    conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port), config)
    if err != nil {
        return nil, fmt.Errorf("failed to establish SSH connection: %w", err)
    }

    return conn, nil
}

func (m *Manager) createSSHConfig(user string) (*ssh.ClientConfig, error) {
    // Load private key
    key, err := m.loadPrivateKey()
    if err != nil {
        return nil, fmt.Errorf("failed to load private key: %w", err)
    }

    // Create SSH client configuration
    config := &ssh.ClientConfig{
        User: user,
        Auth: []ssh.AuthMethod{
            ssh.PublicKeys(key),
        },
        HostKeyCallback: m.hostKeyManager.GetHostKeyCallback(),
        Timeout:         m.config.Timeout,
    }

    return config, nil
}

func (m *Manager) loadPrivateKey() (ssh.Signer, error) {
    // Read private key file
    keyData, err := os.ReadFile(m.config.KeyPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read private key file: %w", err)
    }

    // Parse private key
    var signer ssh.Signer
    if m.config.Passphrase != "" {
        signer, err = ssh.ParsePrivateKeyWithPassphrase(keyData, []byte(m.config.Passphrase))
    } else {
        signer, err = ssh.ParsePrivateKey(keyData)
    }
    if err != nil {
        return nil, fmt.Errorf("failed to parse private key: %w", err)
    }

    return signer, nil
}
```

### SSH Command Execution Implementation

```go
// ExecuteCommand executes a command on a remote machine
func (m *Manager) ExecuteCommand(conn *ssh.Client, command string) (string, error) {
    m.logger.Info("Executing SSH command", map[string]interface{}{
        "command": command,
    })

    // Create SSH session
    session, err := conn.NewSession()
    if err != nil {
        m.logger.Error("Failed to create SSH session", err, map[string]interface{}{
            "command": command,
        })
        return "", fmt.Errorf("failed to create SSH session: %w", err)
    }
    defer session.Close()

    // Execute command
    output, err := session.CombinedOutput(command)
    if err != nil {
        m.logger.Error("Failed to execute SSH command", err, map[string]interface{}{
            "command": command,
            "output":  string(output),
        })
        return string(output), fmt.Errorf("failed to execute SSH command: %w", err)
    }

    m.logger.Info("SSH command executed successfully", map[string]interface{}{
        "command": command,
        "output":  string(output),
    })

    return string(output), nil
}

// TestConnection tests SSH connectivity to a machine
func (m *Manager) TestConnection(hostname string, port int, user string) error {
    m.logger.Info("Testing SSH connection", map[string]interface{}{
        "hostname": hostname,
        "port":     port,
        "user":     user,
    })

    // Get connection
    conn, err := m.GetConnection(hostname, port, user)
    if err != nil {
        m.logger.Error("SSH connection test failed", err, map[string]interface{}{
            "hostname": hostname,
        })
        return fmt.Errorf("SSH connection test failed: %w", err)
    }
    defer m.ReturnConnection(conn)

    // Test with simple command
    _, err = m.ExecuteCommand(conn, "echo 'SSH connection test successful'")
    if err != nil {
        m.logger.Error("SSH command test failed", err, map[string]interface{}{
            "hostname": hostname,
        })
        return fmt.Errorf("SSH command test failed: %w", err)
    }

    m.logger.Info("SSH connection test successful", map[string]interface{}{
        "hostname": hostname,
    })

    return nil
}
```

### SSH File Transfer Implementation

```go
// CopyFile copies a file to a remote machine
func (m *Manager) CopyFile(conn *ssh.Client, localPath string, remotePath string) error {
    m.logger.Info("Copying file to remote machine", map[string]interface{}{
        "local":  localPath,
        "remote": remotePath,
    })

    // Open local file
    localFile, err := os.Open(localPath)
    if err != nil {
        m.logger.Error("Failed to open local file", err, map[string]interface{}{
            "local": localPath,
        })
        return fmt.Errorf("failed to open local file: %w", err)
    }
    defer localFile.Close()

    // Create SSH session
    session, err := conn.NewSession()
    if err != nil {
        m.logger.Error("Failed to create SSH session", err)
        return fmt.Errorf("failed to create SSH session: %w", err)
    }
    defer session.Close()

    // Create remote file
    remoteFile, err := session.StdinPipe()
    if err != nil {
        m.logger.Error("Failed to create stdin pipe", err)
        return fmt.Errorf("failed to create stdin pipe: %w", err)
    }

    // Start scp command
    if err := session.Start(fmt.Sprintf("scp -t %s", remotePath)); err != nil {
        m.logger.Error("Failed to start scp command", err)
        return fmt.Errorf("failed to start scp command: %w", err)
    }

    // Copy file data
    if _, err := io.Copy(remoteFile, localFile); err != nil {
        m.logger.Error("Failed to copy file data", err)
        return fmt.Errorf("failed to copy file data: %w", err)
    }

    // Close remote file
    if err := remoteFile.Close(); err != nil {
        m.logger.Error("Failed to close remote file", err)
        return fmt.Errorf("failed to close remote file: %w", err)
    }

    // Wait for session to complete
    if err := session.Wait(); err != nil {
        m.logger.Error("Failed to complete file copy", err)
        return fmt.Errorf("failed to complete file copy: %w", err)
    }

    m.logger.Info("File copied successfully", map[string]interface{}{
        "local":  localPath,
        "remote": remotePath,
    })

    return nil
}

// GetFile retrieves a file from a remote machine
func (m *Manager) GetFile(conn *ssh.Client, remotePath string, localPath string) error {
    m.logger.Info("Retrieving file from remote machine", map[string]interface{}{
        "remote": remotePath,
        "local":  localPath,
    })

    // Create SSH session
    session, err := conn.NewSession()
    if err != nil {
        m.logger.Error("Failed to create SSH session", err)
        return fmt.Errorf("failed to create SSH session: %w", err)
    }
    defer session.Close()

    // Create local file
    localFile, err := os.Create(localPath)
    if err != nil {
        m.logger.Error("Failed to create local file", err, map[string]interface{}{
            "local": localPath,
        })
        return fmt.Errorf("failed to create local file: %w", err)
    }
    defer localFile.Close()

    // Get remote file
    remoteFile, err := session.StdoutPipe()
    if err != nil {
        m.logger.Error("Failed to create stdout pipe", err)
        return fmt.Errorf("failed to create stdout pipe: %w", err)
    }

    // Start scp command
    if err := session.Start(fmt.Sprintf("scp -f %s", remotePath)); err != nil {
        m.logger.Error("Failed to start scp command", err)
        return fmt.Errorf("failed to start scp command: %w", err)
    }

    // Copy file data
    if _, err := io.Copy(localFile, remoteFile); err != nil {
        m.logger.Error("Failed to copy file data", err)
        return fmt.Errorf("failed to copy file data: %w", err)
    }

    // Wait for session to complete
    if err := session.Wait(); err != nil {
        m.logger.Error("Failed to complete file retrieval", err)
        return fmt.Errorf("failed to complete file retrieval: %w", err)
    }

    m.logger.Info("File retrieved successfully", map[string]interface{}{
        "remote": remotePath,
        "local":  localPath,
    })

    return nil
}
```

## Type Definitions

### SSH Types

```go
// SSHConfig represents SSH configuration
type SSHConfig struct {
    // SSH key path
    KeyPath string `json:"key_path" hcl:"key_path"`

    // SSH key passphrase (optional)
    Passphrase string `json:"passphrase,omitempty" hcl:"passphrase,optional"`

    // SSH timeout in seconds
    Timeout time.Duration `json:"timeout" hcl:"timeout"`

    // SSH connection retries
    Retries int `json:"retries,omitempty" hcl:"retries,optional"`

    // SSH connection retry delay in seconds
    RetryDelay time.Duration `json:"retry_delay,omitempty" hcl:"retry_delay,optional"`

    // SSH host key checking (default: true)
    HostKeyChecking bool `json:"host_key_checking" hcl:"host_key_checking"`

    // SSH known hosts file
    KnownHostsFile string `json:"known_hosts_file,omitempty" hcl:"known_hosts_file,optional"`

    // SSH connection pool size
    PoolSize int `json:"pool_size,omitempty" hcl:"pool_size,optional"`

    // SSH connection pool timeout
    PoolTimeout time.Duration `json:"pool_timeout,omitempty" hcl:"pool_timeout,optional"`
}

// SSHConnection represents an SSH connection
type SSHConnection struct {
    // Connection client
    Client *ssh.Client `json:"-" hcl:"-"`

    // Connection hostname
    Hostname string `json:"hostname" hcl:"hostname"`

    // Connection port
    Port int `json:"port" hcl:"port"`

    // Connection user
    User string `json:"user" hcl:"user"`

    // Connection creation time
    CreatedAt time.Time `json:"created_at" hcl:"created_at"`

    // Connection last used time
    LastUsed time.Time `json:"last_used" hcl:"last_used"`

    // Connection status
    Status string `json:"status" hcl:"status"`
}

// SSHCommand represents an SSH command
type SSHCommand struct {
    // Command to execute
    Command string `json:"command" hcl:"command"`

    // Command timeout
    Timeout time.Duration `json:"timeout,omitempty" hcl:"timeout,optional"`

    // Command working directory
    WorkingDir string `json:"working_dir,omitempty" hcl:"working_dir,optional"`

    // Command environment variables
    Environment map[string]string `json:"environment,omitempty" hcl:"environment,optional"`

    // Command user
    User string `json:"user,omitempty" hcl:"user,optional"`
}

// SSHCommandResult represents the result of an SSH command
type SSHCommandResult struct {
    // Command that was executed
    Command string `json:"command" hcl:"command"`

    // Command output
    Output string `json:"output" hcl:"output"`

    // Command error
    Error string `json:"error,omitempty" hcl:"error,optional"`

    // Command exit code
    ExitCode int `json:"exit_code" hcl:"exit_code"`

    // Command execution time
    ExecutionTime time.Duration `json:"execution_time" hcl:"execution_time"`

    // Command timestamp
    Timestamp time.Time `json:"timestamp" hcl:"timestamp"`
}
```

### SSH Configuration Types

```go
// SSHClientConfig represents SSH client configuration
type SSHClientConfig struct {
    // SSH server hostname
    Hostname string `json:"hostname" hcl:"hostname"`

    // SSH server port
    Port int `json:"port" hcl:"port"`

    // SSH user
    User string `json:"user" hcl:"user"`

    // SSH authentication method
    AuthMethod string `json:"auth_method" hcl:"auth_method"`

    // SSH key path (for key authentication)
    KeyPath string `json:"key_path,omitempty" hcl:"key_path,optional"`

    // SSH key passphrase (for key authentication)
    Passphrase string `json:"passphrase,omitempty" hcl:"passphrase,optional"`

    // SSH password (for password authentication)
    Password string `json:"password,omitempty" hcl:"password,optional"`

    // SSH timeout
    Timeout time.Duration `json:"timeout" hcl:"timeout"`

    // SSH host key checking
    HostKeyChecking bool `json:"host_key_checking" hcl:"host_key_checking"`

    // SSH known hosts file
    KnownHostsFile string `json:"known_hosts_file,omitempty" hcl:"known_hosts_file,optional"`
}

// SSHConnectionPool represents an SSH connection pool
type SSHConnectionPool struct {
    // Pool connections
    Connections map[string]*SSHConnection `json:"connections" hcl:"connections"`

    // Pool mutex
    Mutex sync.RWMutex `json:"-" hcl:"-"`

    // Pool configuration
    Config *SSHConfig `json:"config" hcl:"config"`

    // Pool size limit
    MaxSize int `json:"max_size" hcl:"max_size"`

    // Pool cleanup interval
    CleanupInterval time.Duration `json:"cleanup_interval" hcl:"cleanup_interval"`
}
```

## CLI Commands

### SSH Test Command

```bash
# Test SSH connectivity to a machine
spooky ssh test web-server

# Test SSH connectivity with custom port
spooky ssh test web-server --port 2222

# Test SSH connectivity with custom user
spooky ssh test web-server --user admin

# Test SSH connectivity with timeout
spooky ssh test web-server --timeout 30
```

### SSH Execute Command

```bash
# Execute command on remote machine
spooky ssh execute web-server "ls -la"

# Execute command with custom user
spooky ssh execute web-server "whoami" --user admin

# Execute command with timeout
spooky ssh execute web-server "sleep 10" --timeout 15

# Execute command with working directory
spooky ssh execute web-server "pwd" --working-dir /tmp
```

### SSH Copy Command

```bash
# Copy file to remote machine
spooky ssh copy web-server local-file.txt /tmp/remote-file.txt

# Copy file from remote machine
spooky ssh copy web-server /tmp/remote-file.txt local-file.txt

# Copy file with custom user
spooky ssh copy web-server local-file.txt /tmp/remote-file.txt --user admin

# Copy file with custom permissions
spooky ssh copy web-server local-file.txt /tmp/remote-file.txt --permissions 0644
```

## Integration Examples

### Basic SSH Connection

```go
// SSH connection example
func connectToMachine(hostname string, port int, user string) error {
    // Create SSH manager
    config := &spookytypes.SSHConfig{
        KeyPath:        "~/.ssh/id_rsa",
        Timeout:        30 * time.Second,
        HostKeyChecking: true,
    }
    
    manager := spookyssh.NewManager(logger, config)
    
    // Get SSH connection
    conn, err := manager.GetConnection(hostname, port, user)
    if err != nil {
        return fmt.Errorf("failed to get SSH connection: %w", err)
    }
    defer manager.ReturnConnection(conn)
    
    fmt.Printf("Connected to %s\n", hostname)
    return nil
}
```

### SSH Command Execution

```go
// SSH command execution example
func executeCommand(hostname string, command string) error {
    // Create SSH manager
    config := &spookytypes.SSHConfig{
        KeyPath:        "~/.ssh/id_rsa",
        Timeout:        30 * time.Second,
        HostKeyChecking: true,
    }
    
    manager := spookyssh.NewManager(logger, config)
    
    // Get SSH connection
    conn, err := manager.GetConnection(hostname, 22, "admin")
    if err != nil {
        return fmt.Errorf("failed to get SSH connection: %w", err)
    }
    defer manager.ReturnConnection(conn)
    
    // Execute command
    output, err := manager.ExecuteCommand(conn, command)
    if err != nil {
        return fmt.Errorf("failed to execute command: %w", err)
    }
    
    fmt.Printf("Command output: %s\n", output)
    return nil
}
```

### SSH File Transfer

```go
// SSH file transfer example
func copyFile(hostname string, localPath string, remotePath string) error {
    // Create SSH manager
    config := &spookytypes.SSHConfig{
        KeyPath:        "~/.ssh/id_rsa",
        Timeout:        30 * time.Second,
        HostKeyChecking: true,
    }
    
    manager := spookyssh.NewManager(logger, config)
    
    // Get SSH connection
    conn, err := manager.GetConnection(hostname, 22, "admin")
    if err != nil {
        return fmt.Errorf("failed to get SSH connection: %w", err)
    }
    defer manager.ReturnConnection(conn)
    
    // Copy file
    if err := manager.CopyFile(conn, localPath, remotePath); err != nil {
        return fmt.Errorf("failed to copy file: %w", err)
    }
    
    fmt.Printf("File copied successfully: %s -> %s\n", localPath, remotePath)
    return nil
}
```

## Error Handling

### SSH Errors

```go
// Error handling example
func handleSSHError(err error) {
    if err == nil {
        return
    }
    
    // Check for specific error types
    switch {
    case strings.Contains(err.Error(), "connection refused"):
        fmt.Println("SSH connection refused - check machine connectivity")
    case strings.Contains(err.Error(), "authentication failed"):
        fmt.Println("SSH authentication failed - check credentials")
    case strings.Contains(err.Error(), "host key verification failed"):
        fmt.Println("SSH host key verification failed - check known hosts")
    case strings.Contains(err.Error(), "timeout"):
        fmt.Println("SSH connection timeout - check network connectivity")
    case strings.Contains(err.Error(), "permission denied"):
        fmt.Println("SSH permission denied - check user privileges")
    default:
        fmt.Printf("SSH error: %v\n", err)
    }
}
```

### Connection Errors

```go
// Connection error handling
func handleConnectionError(err error) error {
    if err == nil {
        return nil
    }
    
    // Check for specific connection error types
    switch {
    case strings.Contains(err.Error(), "failed to create SSH connection"):
        return fmt.Errorf("SSH connection creation failed - check SSH configuration")
    case strings.Contains(err.Error(), "failed to load private key"):
        return fmt.Errorf("SSH private key loading failed - check key file and permissions")
    case strings.Contains(err.Error(), "failed to parse private key"):
        return fmt.Errorf("SSH private key parsing failed - check key format and passphrase")
    case strings.Contains(err.Error(), "failed to establish SSH connection"):
        return fmt.Errorf("SSH connection establishment failed - check network and SSH service")
    default:
        return fmt.Errorf("SSH connection error: %w", err)
    }
}
```

## Performance Considerations

### Connection Pooling

The SSH system supports connection pooling:

- SSH connections are pooled and reused
- Configurable pool size and timeout
- Thread-safe connection management

### Resource Management

The SSH system manages resources efficiently:

- SSH connections are properly closed
- Memory usage is optimized for large connection sets
- Timeouts prevent hanging connections

## Troubleshooting

### Common Issues

1. **Connection Refused**: Check machine connectivity and SSH service
2. **Authentication Failed**: Verify SSH key permissions and user access
3. **Host Key Verification Failed**: Check known hosts configuration
4. **Permission Denied**: Check SSH key permissions and user privileges
5. **Timeout Issues**: Adjust SSH timeouts for slow connections

### Debug Information

The SSH system provides comprehensive logging for debugging:

```go
// Enable debug logging
logger.SetLevel(spookytypes.LogLevelDebug)

// Check SSH configuration
fmt.Printf("SSH config: %+v\n", sshConfig)

// Validate SSH key
err := validateSSHKey(keyPath)
if err != nil {
    fmt.Printf("SSH key validation error: %v\n", err)
}
```

## Future Enhancements

### Planned Features

1. **Parallel Processing**: Implement parallel SSH operations
2. **Advanced Authentication**: Support for multiple authentication methods
3. **Connection Monitoring**: Add connection health monitoring
4. **Advanced File Transfer**: Improve file transfer reliability
5. **SSH Tunneling**: Add SSH tunneling support

### Integration Enhancements

1. **Actions Integration**: Use SSH in action execution
2. **Facts Integration**: Use SSH in fact collection
3. **Machines Integration**: Use SSH in machine connectivity testing
4. **Advanced Security**: Improve SSH security features
