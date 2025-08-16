# SSH System User Guide

## Overview

The spooky SSH system provides comprehensive SSH connectivity, authentication, and acting capabilities for remote machine operations. This guide covers everything from basic SSH configuration to advanced features like connection pooling, authentication methods, and acting operations.

**Status: Production Ready** - The SSH system is fully implemented with comprehensive connectivity, authentication, and acting capabilities.

## Getting Started

### Prerequisites

- spooky CLI installed and configured
- SSH access to target machines
- Basic understanding of SSH authentication methods
- Access to create and modify project files

### Quick Start

1. **Check Available SSH Commands**
   ```bash
   spooky ssh --help
   ```

2. **Test SSH Connectivity to a Machine**
   ```bash
   spooky machines ping ./my-project --machine web-server
   ```

3. **Connect to a Machine via SSH**
   ```bash
   spooky machines connect ./my-project --machine web-server
   ```

## Core Concepts

### SSH Authentication Methods

spooky supports multiple SSH authentication methods:

- **SSH Key Authentication** (recommended)
- **Password Authentication** (for testing only)
- **Certificate-based Authentication**
- **Multi-factor Authentication**

### Connection Management

The SSH system provides:

- **Connection Pooling** - Reuse connections for efficiency
- **Connection Limits** - Prevent resource exhaustion
- **Timeout Management** - Handle network issues gracefully
- **Health Monitoring** - Track connection status

### Acting Operations

SSH acting operations include:

- **Command Execution** - Run commands on remote machines
- **File Transfer** - Upload and download files
- **Script Execution** - Run scripts with proper environment
- **Service Control** - Start, stop, and manage services

## Configuration

### Machine SSH Configuration

Configure SSH settings in your `machines.hcl` file:

```hcl
machines {
  machine "web-server" {
    hostname = "web.example.com"
    port = 22
    user = "admin"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/id_rsa"
    }
    
    connection {
      timeout_seconds = 30
      retry_attempts = 3
      keepalive_seconds = 60
    }
    
    tags = ["web", "production"]
  }
}
```

### SSH Key Configuration

For SSH key authentication:

```hcl
authentication {
  method = "ssh_key"
  key_path = "~/.ssh/id_rsa"
  passphrase = "your-passphrase"  # Optional
}
```

### Password Authentication

For password authentication (testing only):

```hcl
authentication {
  method = "password"
  password = "your-password"
}
```

### Certificate Authentication

For certificate-based authentication:

```hcl
authentication {
  method = "certificate"
  certificate_path = "~/.ssh/user-cert.pub"
  key_path = "~/.ssh/id_rsa"
}
```

## CLI Commands

### SSH Connectivity Testing

Test SSH connectivity to machines:

```bash
# Test connectivity to all machines
spooky machines ping ./my-project

# Test specific machine
spooky machines ping ./my-project --machine web-server

# Test machines by tags
spooky machines ping ./my-project --tags production

# Test with custom timeout
spooky machines ping ./my-project --timeout 60
```

### SSH Connection Management

Connect to machines via SSH:

```bash
# Connect to specific machine
spooky machines connect ./my-project --machine web-server

# Connect with custom user
spooky machines connect ./my-project --machine web-server --user admin

# Connect with custom port
spooky machines connect ./my-project --machine web-server --port 2222
```

### SSH Acting Operations

Run actions that use SSH:

```bash
# Run actions on machines
spooky actions run ./my-project --machine web-server

# Run with parallel execution
spooky actions run ./my-project --parallel 4

# Run with dry-run mode
spooky actions run ./my-project --dry-run

# Run with decryption
spooky actions run ./my-project --decrypt
```

## Advanced Features

### Connection Pooling

The SSH system automatically manages connection pooling:

- **Reuse Connections** - Maintain connections for efficiency
- **Connection Limits** - Prevent resource exhaustion
- **Health Checks** - Monitor connection status
- **Automatic Cleanup** - Close idle connections

### Timeout Management

Configure timeouts for different operations:

```hcl
connection {
  timeout_seconds = 30        # Connection timeout
  command_timeout_seconds = 60 # Command execution timeout
  keepalive_seconds = 60      # Keepalive interval
}
```

### Retry Logic

Configure retry behavior for failed operations:

```hcl
connection {
  retry_attempts = 3
  retry_delay_seconds = 5
  backoff_multiplier = 2.0
}
```

### Health Monitoring

Monitor SSH connection health:

```bash
# Check connection status
spooky machines ping ./my-project --verbose

# Monitor connection metrics
spooky machines ping ./my-project --metrics
```

## Security Best Practices

### SSH Key Management

- Use Ed25519 or RSA 4096-bit keys
- Set proper file permissions (600)
- Use passphrases for additional security
- Rotate keys regularly

### Authentication Security

- Prefer SSH key authentication over passwords
- Use certificate-based authentication for large deployments
- Implement multi-factor authentication where possible
- Monitor authentication attempts

### Network Security

- Use SSH over secure networks
- Implement firewall rules for SSH access
- Use non-standard SSH ports when possible
- Monitor SSH access logs

## Troubleshooting

### Common SSH Issues

**Connection Refused**
```bash
# Check if SSH service is running
spooky machines ping ./my-project --machine web-server --verbose
```

**Authentication Failed**
```bash
# Verify SSH key permissions
ls -la ~/.ssh/id_rsa

# Test SSH key manually
ssh -i ~/.ssh/id_rsa user@hostname
```

**Timeout Issues**
```bash
# Increase timeout for slow connections
spooky machines ping ./my-project --timeout 120
```

### Debugging SSH Connections

Enable verbose output for debugging:

```bash
# Verbose SSH output
spooky machines ping ./my-project --verbose

# Debug connection issues
spooky machines connect ./my-project --machine web-server --debug
```

### Performance Optimization

Optimize SSH performance:

```bash
# Use connection pooling
spooky actions run ./my-project --parallel 4

# Monitor connection metrics
spooky machines ping ./my-project --metrics
```

## Integration with Other Systems

### Actions Integration

SSH integrates with the actions system for remote execution:

```hcl
actions {
  action "update-system" {
    description = "Update system packages"
    
    machines = ["web-server"]
    parallel = true
    
    command = "sudo apt update && sudo apt upgrade -y"
  }
}
```

### Facts Integration

SSH enables fact collection from remote machines:

```bash
# Collect facts via SSH
spooky facts gather ./my-project --parallel 4
```

### Variables Integration

SSH enables variable collection from remote machines:

```bash
# Collect variables via SSH
spooky variables gather ./my-project --parallel 4
```

## Examples

### Basic SSH Configuration

```hcl
machines {
  machine "web-server" {
    hostname = "web.example.com"
    port = 22
    user = "admin"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/id_rsa"
    }
    
    tags = ["web", "production"]
  }
}
```

### Advanced SSH Configuration

```hcl
machines {
  machine "database-server" {
    hostname = "db.example.com"
    port = 2222
    user = "dbadmin"
    
    authentication {
      method = "certificate"
      certificate_path = "~/.ssh/db-cert.pub"
      key_path = "~/.ssh/id_rsa"
    }
    
    connection {
      timeout_seconds = 60
      retry_attempts = 5
      keepalive_seconds = 30
    }
    
    tags = ["database", "production"]
  }
}
```

### SSH Acting Example

```hcl
actions {
  action "deploy-application" {
    description = "Deploy web application"
    
    machines = ["web-server"]
    parallel = true
    
    template {
      source = "templates/deploy.sh.tmpl"
      destination = "/tmp/deploy.sh"
      permissions = "0755"
    }
    
    command = "/tmp/deploy.sh"
  }
}
```

## Best Practices

### SSH Configuration

- Use descriptive machine names
- Group machines with tags
- Configure appropriate timeouts
- Use connection pooling for efficiency

### Security

- Use SSH key authentication
- Set proper file permissions
- Monitor authentication attempts
- Implement access controls

### Performance

- Use parallel execution for multiple machines
- Configure appropriate connection limits
- Monitor connection metrics
- Optimize timeout settings

### Monitoring

- Monitor SSH connection health
- Track authentication failures
- Monitor performance metrics
- Log SSH operations

## Next Steps

- Explore the [SSH API Reference](SSH_API_REFERENCE.md) for detailed technical information
- Check the [SSH Troubleshooting Guide](SSH_TROUBLESHOOTING.md) for common issues
- Review the [SSH Documentation Summary](SSH_DOCUMENTATION_SUMMARY.md) for implementation details
- Learn about [SSH Integration Patterns](INTEGRATIONS_USER_GUIDE.md) for advanced usage
