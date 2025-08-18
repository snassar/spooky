# SSH System

## Overview

The SSH System provides comprehensive SSH connectivity and acting capabilities for the spooky codebase. It enables secure SSH connections, authentication management, command running, file transfer, and acting operations across multiple machines.

**Status**: **Implemented** - Complete SSH system with connection management, authentication, acting, and CLI integration.

## Related Systems

This system integrates with and depends on several other spooky systems:

- **[Machines System](MACHINES_SYSTEM.md)** - SSH connects to machines defined in the inventory
- **[Actions System](ACTIONS_SYSTEM.md)** - SSH runs actions on remote machines
- **[Facts System](FACTS_SYSTEM.md)** - SSH collects facts from remote machines
- **[Logging System](LOGGING_SYSTEM.md)** - SSH operations generate comprehensive logs for monitoring and debugging
- **[Integrations System](INTEGRATIONS_SYSTEM.md)** - SSH integrates with other systems through the IntegrationManager
- **[Templates System](TEMPLATES_SYSTEM.md)** - SSH deploys rendered templates to remote machines
- **[Variables System](VARIABLES_SYSTEM.md)** - SSH configurations can use variables
- **[Projects System](PROJECTS_SYSTEM.md)** - SSH configurations are organized within projects

## Architecture

### Core Components

#### SSH Manager
- **File**: `internal/ssh/manager.go`
- **Purpose**: Central SSH management with connection and acting coordination
- **Features**:
  - SSH connection management
  - Authentication handling
  - Acting session management
  - Connection pooling
  - Error handling and recovery
  - Performance monitoring

#### SSH Client
- **File**: `internal/ssh/client.go`
- **Purpose**: SSH client implementation with connection and authentication
- **Features**:
  - SSH connection establishment
  - Authentication methods
  - Connection validation
  - Error handling
  - Performance optimization
  - Security validation

#### Connection Pool
- **File**: `internal/ssh/connection_pool.go`
- **Purpose**: SSH connection pooling for performance and efficiency
- **Features**:
  - Connection pooling
  - Connection reuse
  - Load balancing
  - Health monitoring
  - Resource management
  - Performance optimization

#### File Transfer
- **File**: `internal/ssh/file_transfer.go`
- **Purpose**: Secure file transfer operations via SSH
- **Features**:
  - File upload and download
  - Directory synchronization
  - Progress monitoring
  - Error handling
  - Security validation
  - Performance optimization

#### Advanced Authentication
- **File**: `internal/ssh/advanced_auth.go`
- **Purpose**: Advanced SSH authentication methods and security
- **Features**:
  - Multiple authentication methods
  - Key management
  - Certificate validation
  - Security hardening
  - Audit logging
  - Compliance support

#### SSH Integration
- **File**: `internal/ssh/integration.go`
- **Purpose**: Interface implementation for system integration
- **Features**:
  - Connect - Establish SSH connections
  - RunCommand - Run commands via SSH
  - TransferFile - Transfer files via SSH
  - ValidateConnection - Validate SSH connections

### Integration Points

#### Machines Integration
- Provides SSH connectivity for machine operations
- Supports machine-specific SSH configurations
- Enables machine-based SSH operations

#### Actions Integration
- Provides SSH connectivity for action running
- Supports SSH-based action acting
- Enables remote command running

#### Facts Integration
- Provides SSH connectivity for fact collection
- Supports SSH-based fact gathering
- Enables remote system information collection

#### Templates Integration
- Provides SSH connectivity for template deployment
- Supports SSH-based file transfer
- Enables remote configuration deployment

## SSH Types

### SSH Connection
```go
type SSHConnection struct {
    ID              string                 // Connection identifier
    Host            string                 // Target host
    Port            int                    // SSH port
    User            string                 // SSH user
    Authentication  *SSHAuthentication    // Authentication configuration
    Client          *ssh.Client            // SSH client
    Session         *ssh.Session           // SSH session
    Status          string                 // Connection status
    LastUsed        time.Time              // Last used timestamp
    CreatedAt       time.Time              // Creation timestamp
}
```

### SSH Authentication
```go
type SSHAuthentication struct {
    Method          string                 // Authentication method
    KeyPath         string                 // SSH key path
    KeyData         []byte                 // SSH key data
    Password        string                 // SSH password
    Passphrase      string                 // Key passphrase
    Certificate     string                 // SSH certificate
    Options         map[string]interface{} // Additional options
}
```

### SSH Acting Session
```go
type SSHActingSession struct {
    ID              string                 // Session identifier
    Connection      *SSHConnection         // SSH connection
    Commands        []string               // Commands to run
    Results         []*ActingResult        // Acting results
    Status          string                 // Session status
    StartedAt       time.Time              // Session start time
    CompletedAt     *time.Time             // Session completion time
}
```

### Acting Result
```go
type ActingResult struct {
    Command         string                 // Run command
    Output          string                 // Command output
    Error           string                 // Command error
    ExitCode        int                    // Command exit code
    Duration        time.Duration          // Running duration
Timestamp       time.Time              // Running timestamp
}
```

## SSH Categories

### Connection Types
- **Direct**: Direct SSH connections
- **Proxy**: SSH connections through proxy
- **Jump Host**: SSH connections through jump host
- **Tunnel**: SSH tunnel connections
- **Multi-hop**: Multi-hop SSH connections

### Authentication Methods
- **SSH Key**: SSH key-based authentication
- **Password**: Password-based authentication
- **Certificate**: SSH certificate authentication
- **Multi-factor**: Multi-factor authentication
- **Agent**: SSH agent authentication

### Acting Types
- **Command Running**: Run single commands
- **Script Running**: Run shell scripts
- **File Transfer**: Transfer files and directories
- **Configuration**: Deploy configurations
- **Monitoring**: Monitor system status

## SSH Management

### Connection Management
- **Connection Establishment**: Establish SSH connections
- **Connection Pooling**: Pool SSH connections
- **Connection Validation**: Validate connection health
- **Connection Cleanup**: Clean up connections
- **Connection Monitoring**: Monitor connection status

### Authentication Management
- **Key Management**: Manage SSH keys
- **Certificate Management**: Manage SSH certificates
- **Password Management**: Manage passwords securely
- **Multi-factor Setup**: Setup multi-factor authentication
- **Security Validation**: Validate authentication security

### Acting Management
- **Session Management**: Manage acting sessions
- **Command Running**: Run commands securely
- **Result Collection**: Collect acting results
- **Error Handling**: Handle acting errors
- **Performance Monitoring**: Monitor acting performance

## SSH Operations

### Connection Operations
- **Connect**: Establish SSH connection
- **Disconnect**: Close SSH connection
- **Reconnect**: Reconnect SSH connection
- **Validate**: Validate connection health
- **Monitor**: Monitor connection status

### Authentication Operations
- **Authenticate**: Authenticate SSH connection
- **Validate Key**: Validate SSH key
- **Validate Certificate**: Validate SSH certificate
- **Test Authentication**: Test authentication methods
- **Update Credentials**: Update authentication credentials

### Acting Operations
- **Run Command**: Run single command
- **Run Script**: Run shell script
- **Transfer File**: Transfer file via SCP/SFTP
- **Sync Directory**: Synchronize directory
- **Monitor System**: Monitor system status

## Security Features

### Connection Security
- **Encryption**: Use strong encryption algorithms
- **Key Exchange**: Secure key exchange protocols
- **Host Validation**: Validate host keys
- **Connection Limits**: Limit concurrent connections
- **Timeout Management**: Manage connection timeouts

### Authentication Security
- **Key Validation**: Validate SSH keys
- **Certificate Validation**: Validate SSH certificates
- **Permission Validation**: Validate file permissions
- **Access Control**: Control SSH access
- **Audit Logging**: Log SSH access

### Acting Security
- **Command Validation**: Validate commands
- **Output Sanitization**: Sanitize command output
- **Error Handling**: Handle security errors
- **Session Isolation**: Isolate acting sessions
- **Resource Limits**: Limit resource usage

## Performance Features

### Connection Pooling
- **Connection Reuse**: Reuse SSH connections
- **Connection Limits**: Limit connection pool size
- **Load Balancing**: Balance load across connections
- **Health Monitoring**: Monitor connection health
- **Resource Optimization**: Optimize resource usage

### Parallel Acting
- **Parallel Running**: Run commands in parallel
- **Concurrent Sessions**: Manage concurrent sessions
- **Resource Management**: Manage acting resources
- **Performance Monitoring**: Monitor acting performance
- **Load Distribution**: Distribute load across machines

### Optimization
- **Connection Optimization**: Optimize connection performance
- **Authentication Optimization**: Optimize authentication
- **Acting Optimization**: Optimize acting performance
- **Resource Optimization**: Optimize resource usage
- **Monitoring**: Monitor performance metrics

## CLI Commands

### SSH Connection
```bash
# Test SSH connection
spooky ssh test <host>

# Connect to host
spooky ssh connect <host>

# List connections
spooky ssh list-connections

# Close connection
spooky ssh disconnect <host>
```

### SSH Authentication
```bash
# Test authentication
spooky ssh test-auth <host>

# Validate SSH key
spooky ssh validate-key <key-file>

# Import SSH key
spooky ssh import-key <key-file>

# List SSH keys
spooky ssh list-keys
```

### SSH Acting
```bash
# Run command
spooky ssh run <host> --command "uname -a"

# Run script
spooky ssh run <host> --script deploy.sh

# Transfer file
spooky ssh transfer <host> --upload local.txt --remote /tmp/remote.txt

# Monitor system
spooky ssh monitor <host> --metrics cpu,memory,disk
```

### SSH Management
```bash
# Show SSH status
spooky ssh status

# Show connection pool
spooky ssh pool-status

# Show SSH configuration
spooky ssh config

# Validate SSH setup
spooky ssh validate
```

## Configuration

### SSH Configuration
```hcl
# ssh/config.hcl
ssh_config {
  # Connection settings
  connection {
    default_timeout = 30  # seconds
    max_connections = 10
    connection_pool_size = 5
    retry_attempts = 3
  }
  
  # Authentication settings
  authentication {
    default_method = "ssh_key"
    key_path = "~/.ssh/id_rsa"
    key_passphrase = ""
    certificate_path = ""
  }
  
  # Acting settings
  acting {
    default_timeout = 300  # seconds
    max_concurrent = 5
    output_buffer_size = 1024
    error_buffer_size = 1024
  }
  
  # Security settings
  security {
    validate_host_keys = true
    validate_key_permissions = true
    audit_logging = true
    access_control = true
  }
  
  # Performance settings
  performance {
    connection_pooling = true
    parallel_running = true
    compression = true
    keepalive = true
  }
}
```

### SSH Host Configuration
```hcl
# ssh/hosts.hcl
ssh_hosts {
  host "web-server" {
    hostname = "web.example.com"
    port = 22
    user = "admin"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/web_keys/id_rsa"
    }
    
    connection {
      timeout = 30
      retries = 3
    }
    
    acting {
      timeout = 300
      max_concurrent = 2
    }
    
    security {
      validate_host_key = true
      strict_host_checking = true
    }
  }
  
  host "db-server" {
    hostname = "db.example.com"
    port = 22
    user = "dbadmin"
    
    authentication {
      method = "certificate"
      certificate = "~/.ssh/db_cert.pem"
      key_path = "~/.ssh/db_key"
    }
    
    connection {
      timeout = 60
      retries = 5
    }
    
    acting {
      timeout = 600
      max_concurrent = 1
    }
    
    security {
      validate_host_key = true
      strict_host_checking = true
      audit_logging = true
    }
  }
}
```

## Examples

### Basic SSH Connection
```hcl
# ssh/basic.hcl
ssh_hosts {
  host "server" {
    hostname = "192.168.1.100"
    port = 22
    user = "admin"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/id_rsa"
    }
  }
}
```

### Advanced SSH Configuration
```hcl
# ssh/advanced.hcl
ssh_hosts {
  host "production-server" {
    hostname = "prod.example.com"
    port = 22
    user = "prodadmin"
    
    authentication {
      method = "certificate"
      certificate = "~/.ssh/prod_cert.pem"
      key_path = "~/.ssh/prod_key"
      passphrase = "prod-key-passphrase"
    }
    
    connection {
      timeout = 60
      retries = 5
      keepalive = true
      keepalive_interval = 30
    }
    
    acting {
      timeout = 600
      max_concurrent = 3
      output_buffer_size = 2048
      error_buffer_size = 2048
    }
    
    security {
      validate_host_key = true
      strict_host_checking = true
      audit_logging = true
      access_control = true
    }
    
    performance {
      compression = true
      connection_pooling = true
      parallel_running = true
    }
  }
}
```

## Integration Examples

### Machines Integration
```go
// Use SSH with machines
sshIntegration := manager.GetSSHIntegration()
machinesIntegration := manager.GetMachinesIntegration()

// Get machine inventory
machines, err := machinesIntegration.GetMachines("my-project")
if err != nil {
    return err
}

// Connect to machines via SSH
for _, machine := range machines {
    connection, err := sshIntegration.Connect(machine.Hostname, machine.Authentication)
    if err != nil {
        log.Printf("Failed to connect to %s: %v", machine.Hostname, err)
        continue
    }
    defer connection.Close()
    
    // Run command on machine
result, err := sshIntegration.RunCommand(connection, "uname -a")
if err != nil {
    log.Printf("Failed to run command on %s: %v", machine.Hostname, err)
    continue
}
    
    log.Printf("Machine %s: %s", machine.Hostname, result.Output)
}
```

### Actions Integration
```go
// Use SSH with actions
sshIntegration := manager.GetSSHIntegration()
actionsIntegration := manager.GetActionsIntegration()

// Run action via SSH
action := &spookytypes.Action{
    Name: "deploy-application",
    Script: "deploy.sh",
    Machine: "web-server",
}

// Get SSH connection for machine
connection, err := sshIntegration.Connect("web-server", authConfig)
if err != nil {
    return err
}
defer connection.Close()

// Run action script
result, err := sshIntegration.RunScript(connection, action.Script)
if err != nil {
    return err
}

log.Printf("Action completed: %s", result.Output)
```

### Facts Integration
```go
// Use SSH for fact collection
sshIntegration := manager.GetSSHIntegration()
factsIntegration := manager.GetFactsIntegration()

// Collect facts via SSH
machines := []string{"web-server", "db-server"}

for _, machine := range machines {
    connection, err := sshIntegration.Connect(machine, authConfig)
    if err != nil {
        log.Printf("Failed to connect to %s: %v", machine, err)
        continue
    }
    defer connection.Close()
    
    // Collect system facts
    facts := make(map[string]interface{})
    
    // CPU info
cpuResult, err := sshIntegration.RunCommand(connection, "lscpu")
    if err == nil {
        facts["cpu"] = parseCPUInfo(cpuResult.Output)
    }
    
    // Memory info
memResult, err := sshIntegration.RunCommand(connection, "free -h")
    if err == nil {
        facts["memory"] = parseMemoryInfo(memResult.Output)
    }
    
    // Store facts
    err = factsIntegration.StoreFacts(machine, facts)
    if err != nil {
        log.Printf("Failed to store facts for %s: %v", machine, err)
    }
}
```

## Best Practices

### Security
- Use SSH keys instead of passwords
- Implement proper key management
- Validate host keys
- Use strong encryption
- Regular security audits

### Performance
- Use connection pooling
- Implement parallel running
- Monitor connection health
- Optimize resource usage
- Use appropriate timeouts

### Reliability
- Implement retry mechanisms
- Handle connection failures
- Monitor connection status
- Implement health checks
- Use connection validation

### Management
- Regular key rotation
- Monitor SSH access
- Implement access controls
- Audit SSH operations
- Document SSH procedures

## Troubleshooting

### Common Issues

#### Connection Issues
```bash
# Test basic connectivity
spooky ssh test <host>

# Check SSH configuration
spooky ssh config

# Validate SSH key
spooky ssh validate-key <key-file>

# Test authentication
spooky ssh test-auth <host>
```

#### Authentication Issues
```bash
# Check key permissions
ls -la ~/.ssh/id_rsa

# Test SSH key
ssh-keygen -l -f ~/.ssh/id_rsa

# Check SSH agent
ssh-add -l

# Test manual SSH connection
ssh -i ~/.ssh/id_rsa user@host
```

#### Acting Issues
```bash
# Test command running
spooky ssh run <host> --command "echo 'test'"

# Check acting timeout
spooky ssh run <host> --command "sleep 10" --timeout 15

# Monitor acting performance
spooky ssh run <host> --command "uname -a" --verbose
```

#### Performance Issues
```bash
# Check connection pool
spooky ssh pool-status

# Monitor SSH performance
spooky ssh performance

# Check connection limits
spooky ssh config --show-limits
```

## API Reference

### SSHIntegration Interface
```go
type SSHIntegration interface {
    Connect(ctx context.Context, host string, auth *spookytypes.SSHAuthentication) (*spookyssh.Connection, error)
    RunCommand(ctx context.Context, connection *spookyssh.Connection, command string) (*spookytypes.ActingResult, error)
    TransferFile(ctx context.Context, connection *spookyssh.Connection, localPath, remotePath string) error
    ValidateConnection(ctx context.Context, host string) (*spookytypes.ValidationResult, error)
}
```

### SSH Manager Methods
```go
// Connection management
Connect(ctx context.Context, host string, auth *spookytypes.SSHAuthentication) (*spookyssh.Connection, error)
Disconnect(ctx context.Context, connection *spookyssh.Connection) error
ValidateConnection(ctx context.Context, host string) (*spookytypes.ValidationResult, error)

// Acting operations
RunCommand(ctx context.Context, connection *spookyssh.Connection, command string) (*spookytypes.ActingResult, error)
RunScript(ctx context.Context, connection *spookyssh.Connection, script string) (*spookytypes.ActingResult, error)
TransferFile(ctx context.Context, connection *spookyssh.Connection, localPath, remotePath string) error

// Session management
CreateActingSession(ctx context.Context, connection *spookyssh.Connection) (*spookyssh.ActingSession, error)
RunSession(ctx context.Context, session *spookyssh.ActingSession, commands []string) ([]*spookytypes.ActingResult, error)
CloseSession(ctx context.Context, session *spookyssh.ActingSession) error
```

## Related Documentation

- [SSH API Reference](SSH_API_REFERENCE.md) - Complete API documentation
- [SSH User Guide](SSH_USER_GUIDE.md) - User guide and examples
- [SSH Troubleshooting](SSH_TROUBLESHOOTING.md) - Troubleshooting guide
- [Machines System](MACHINES_SYSTEM.md) - Machines integration
- [Actions System](ACTIONS_SYSTEM.md) - Actions integration
- [Facts System](FACTS_SYSTEM.md) - Facts integration
- [Templates System](TEMPLATES_SYSTEM.md) - Templates integration
