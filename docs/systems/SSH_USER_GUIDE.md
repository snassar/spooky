# SSH User Guide

## Overview

The SSH system in spooky provides comprehensive functionality for SSH connections, authentication, command execution, and file transfer operations. This guide covers how to configure and use SSH features effectively.

**Status: Fully Implemented** - All SSH functionality is implemented and ready for production use.

## Related Documentation

- [Machines User Guide](MACHINES_USER_GUIDE.md) - Machine inventory and connectivity testing
- [Actions User Guide](ACTIONS_USER_GUIDE.md) - SSH-based action orchestration
- [Facts User Guide](FACTS_USER_GUIDE.md) - SSH-based fact collection
- [Templates User Guide](TEMPLATES_USER_GUIDE.md) - SSH-based template rendering
- [Variables User Guide](VARIABLES_USER_GUIDE.md) - Variable usage in SSH operations

> **See also**: [User Guides Index](USER_GUIDES_INDEX.md) - Complete overview of all user guides

## Quick Start

### Basic SSH Configuration

1. **Configure SSH keys**:
   ```bash
   # Generate Ed25519 key (recommended)
   ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -C "your-email@example.com"
   
   # Or generate RSA 4096-bit key
   ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa -C "your-email@example.com"
   ```

2. **Set proper permissions**:
   ```bash
   chmod 600 ~/.ssh/id_ed25519
   chmod 644 ~/.ssh/id_ed25519.pub
   ```

3. **Add public key to remote servers**:
   ```bash
   ssh-copy-id -i ~/.ssh/id_ed25519.pub user@hostname
   ```

### Test SSH Connectivity

```bash
# Test SSH connectivity to machines
spooky machines ping my-project

# Test specific machine
spooky machines ping my-project --machine web-server
```

> **Note**: SSH connectivity testing is integrated with the [Machines User Guide](MACHINES_USER_GUIDE.md) system. The `spooky machines ping` command uses SSH to test connectivity to target machines.

## Configuration

### Machine Configuration

Configure SSH settings in your machine inventory:

```hcl
machines {
  machine "web-server" {
    hostname = "web.example.com"
    port = 22
    user = "admin"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/id_ed25519"
      passphrase = "optional_passphrase"
    }
    
    ssh {
      timeout = "30s"
      keepalive = "60s"
      host_key_validation = true
    }
    
    tags = {
      environment = "production"
      role = "web"
    }
  }
  
  machine "db-server" {
    hostname = "db.example.com"
    port = 22
    user = "dbadmin"
    
    authentication {
      method = "password"
      password = "secure_password"
    }
    
    ssh {
      timeout = "60s"
      keepalive = "120s"
      host_key_validation = false  # Development only
    }
    
    tags = {
      environment = "production"
      role = "database"
    }
  }
}
```

### SSH Client Configuration

Configure global SSH settings:

```hcl
ssh {
  default_port = 22
  default_timeout = "30s"
  max_connections = 10
  max_retry_attempts = 3
  retry_delay = "5s"
  idle_timeout = "300s"
  known_hosts_path = "~/.ssh/known_hosts"
  strict_host_key_check = true
  allow_insecure_hosts = false
}
```

## Authentication Methods

### SSH Key Authentication (Recommended)

SSH key authentication is the most secure and recommended method:

```hcl
authentication {
  method = "ssh_key"
  key_path = "~/.ssh/id_ed25519"
  passphrase = "optional_passphrase"
}
```

**Supported Key Types**:
- **Ed25519**: Recommended for most use cases
- **RSA 4096-bit**: For compatibility with older systems

**Key Requirements**:
- Private key permissions: 600 (`-rw-------`)
- Public key permissions: 644 (`-rw-r--r--`)
- Key size: RSA keys must be at least 4096 bits

### Password Authentication

Password authentication is supported but less secure:

```hcl
authentication {
  method = "password"
  password = "user_password"
}
```

**Security Considerations**:
- Passwords are stored in plain text in configuration
- Use SSH keys when possible
- Consider using environment variables for passwords

### SSH Agent Authentication

Use SSH agent for key management:

```hcl
authentication {
  method = "agent"
}
```

**Setup**:
```bash
# Start SSH agent
eval $(ssh-agent)

# Add key to agent
ssh-add ~/.ssh/id_ed25519
```

## Connection Management

### Connection Pooling

The SSH system uses connection pooling for improved performance:

```hcl
ssh {
  max_connections = 10      # Maximum concurrent connections
  idle_timeout = "300s"     # Close idle connections after 5 minutes
  keepalive = "60s"         # Send keepalive every 60 seconds
}
```

**Benefits**:
- Reuse connections for multiple operations
- Reduce connection establishment overhead
- Improve performance for parallel operations

### Connection Health Monitoring

The system automatically monitors connection health:

- **Health Checks**: Regular health checks for active connections
- **Automatic Cleanup**: Clean up stale or unhealthy connections
- **Retry Logic**: Automatic retry for failed connections

### Timeout Configuration

Configure appropriate timeouts for your environment:

```hcl
ssh {
  default_timeout = "30s"   # Connection timeout
  command_timeout = "60s"   # Command execution timeout
  transfer_timeout = "300s" # File transfer timeout
}
```

## Command Execution

### Basic Command Execution

Execute commands on remote machines through actions:

```hcl
actions {
  action "check-system-status" {
    description = "Check system status on all machines"
    
    machines = ["web-server", "db-server"]
    parallel = true
    
    command {
      command = "systemctl status"
      timeout = 30
    }
  }
  
  action "update-packages" {
    description = "Update system packages"
    
    machines = ["web-server"]
    
    command {
      command = "apt update && apt upgrade -y"
      timeout = 300
      working_dir = "/tmp"
      environment = {
        DEBIAN_FRONTEND = "noninteractive"
      }
    }
  }
}
```

### Command with Stdin Input

Execute commands with standard input:

```hcl
actions {
  action "create-user" {
    description = "Create new user account"
    
    machines = ["web-server"]
    
    command {
      command = "adduser --gecos '' newuser"
      stdin = "password\npassword\n"  # Provide password twice
      timeout = 60
    }
  }
}
```

### Environment Variables

Set environment variables for commands:

```hcl
actions {
  action "deploy-application" {
    description = "Deploy application with environment variables"
    
    machines = ["web-server"]
    
    command {
      command = "./deploy.sh"
      working_dir = "/opt/app"
      environment = {
        NODE_ENV = "production"
        DATABASE_URL = "postgresql://user:pass@localhost/db"
        API_KEY = "secret-key"
      }
      timeout = 120
    }
  }
}
```

## File Transfer

### SFTP File Transfer

Transfer files using SFTP (recommended):

```hcl
actions {
  action "deploy-config" {
    description = "Deploy configuration files"
    
    machines = ["web-server"]
    
    file_copy {
      source = "config/nginx.conf"
      destination = "/etc/nginx/nginx.conf"
      permissions = "0644"
      verify = true
    }
  }
  
  action "backup-logs" {
    description = "Backup log files from remote server"
    
    machines = ["web-server"]
    
    file_copy {
      source = "/var/log/nginx/access.log"
      destination = "backups/nginx-access.log"
      direction = "download"
      verify = true
    }
  }
}
```

### SCP File Transfer

Transfer files using SCP:

```hcl
actions {
  action "upload-script" {
    description = "Upload script file using SCP"
    
    machines = ["web-server"]
    
    file_copy {
      source = "scripts/deploy.sh"
      destination = "/tmp/deploy.sh"
      mode = "scp"
      permissions = "0755"
    }
  }
}
```

### Directory Transfer

Transfer entire directories:

```hcl
actions {
  action "deploy-application" {
    description = "Deploy entire application directory"
    
    machines = ["web-server"]
    
    file_copy {
      source = "app/"
      destination = "/opt/app/"
      recursive = true
      permissions = "0755"
      exclude = ["*.tmp", "*.log"]
    }
  }
}
```

## Service Management

### Service Control Actions

Control system services:

```hcl
actions {
  action "restart-nginx" {
    description = "Restart nginx service"
    
    machines = ["web-server"]
    
    service_control {
      service = "nginx"
      action = "restart"
      timeout = 60
    }
  }
  
  action "stop-services" {
    description = "Stop all application services"
    
    machines = ["web-server", "db-server"]
    parallel = true
    
    service_control {
      service = "app"
      action = "stop"
      timeout = 30
    }
  }
}
```

**Supported Service Actions**:
- `start`: Start a service
- `stop`: Stop a service
- `restart`: Restart a service
- `reload`: Reload service configuration
- `status`: Check service status
- `enable`: Enable service at boot
- `disable`: Disable service at boot

## Template Deployment

### Template-Based File Deployment

Deploy files with template rendering:

```hcl
actions {
  action "deploy-config-template" {
    description = "Deploy configuration from template"
    
    machines = ["web-server"]
    
    template_deploy {
      source = "templates/nginx.conf.tmpl"
      destination = "/etc/nginx/nginx.conf"
      permissions = "0644"
      variables = {
        server_name = "example.com"
        port = 80
        ssl_enabled = true
      }
    }
  }
}
```

### Template with Machine Facts

Use machine facts in templates:

```hcl
actions {
  action "deploy-host-config" {
    description = "Deploy host-specific configuration"
    
    machines = ["web-server", "db-server"]
    parallel = true
    
    template_deploy {
      source = "templates/hosts.tmpl"
      destination = "/etc/hosts"
      permissions = "0644"
      variables = {
        hostname = "{{.Machine.Hostname}}"
        ip_address = "{{.Machine.Host}}"
        environment = "{{.Machine.Tags.environment}}"
      }
    }
  }
}
```

## Parallel Operations

### Parallel Command Execution

Execute commands in parallel across multiple machines:

```hcl
actions {
  action "parallel-update" {
    description = "Update packages on all machines in parallel"
    
    machines = ["web-server", "db-server", "app-server"]
    parallel = true
    max_concurrent = 3
    
    command {
      command = "apt update && apt upgrade -y"
      timeout = 300
    }
  }
}
```

### Parallel File Transfer

Transfer files in parallel:

```hcl
actions {
  action "parallel-deploy" {
    description = "Deploy application to all machines in parallel"
    
    machines = ["web-server", "db-server", "app-server"]
    parallel = true
    max_concurrent = 2
    
    file_copy {
      source = "app/"
      destination = "/opt/app/"
      recursive = true
      permissions = "0755"
    }
  }
}
```

## Error Handling

### Command Error Handling

Handle command execution errors:

```hcl
actions {
  action "safe-update" {
    description = "Safe system update with error handling"
    
    machines = ["web-server"]
    
    command {
      command = "apt update && apt upgrade -y"
      timeout = 300
      continue_on_error = true  # Continue even if this command fails
    }
    
    command {
      command = "systemctl restart nginx"
      timeout = 60
      depends_on = ["safe-update"]  # Only run if previous command succeeds
    }
  }
}
```

### File Transfer Error Handling

Handle file transfer errors:

```hcl
actions {
  action "backup-with-retry" {
    description = "Backup files with retry logic"
    
    machines = ["web-server"]
    
    file_copy {
      source = "/var/log/nginx/access.log"
      destination = "backups/nginx-access.log"
      direction = "download"
      verify = true
      retry_attempts = 3
      retry_delay = "5s"
    }
  }
}
```

## Security Best Practices

### SSH Key Management

1. **Use strong keys**:
   ```bash
   # Generate Ed25519 key (recommended)
   ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -C "your-email@example.com"
   
   # Or RSA 4096-bit key
   ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa -C "your-email@example.com"
   ```

2. **Set proper permissions**:
   ```bash
   chmod 600 ~/.ssh/id_ed25519
   chmod 644 ~/.ssh/id_ed25519.pub
   ```

3. **Use passphrases**:
   ```bash
   ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -C "your-email@example.com"
   # Enter a strong passphrase when prompted
   ```

### Host Key Validation

Enable host key validation for security:

```hcl
ssh {
  strict_host_key_check = true
  known_hosts_path = "~/.ssh/known_hosts"
}
```

**Add host keys**:
```bash
# Add host key to known_hosts
ssh-keyscan -H web.example.com >> ~/.ssh/known_hosts
```

### Network Security

1. **Use non-standard SSH ports** (if configured):
   ```hcl
   machine "web-server" {
     hostname = "web.example.com"
     port = 2222  # Non-standard SSH port
     user = "admin"
   }
   ```

2. **Use VPN or private networks** for sensitive operations
3. **Limit SSH access** to specific IP addresses
4. **Monitor SSH logs** for suspicious activity

## Performance Optimization

### Connection Pooling

Optimize connection performance:

```hcl
ssh {
  max_connections = 20      # Increase for high-concurrency operations
  idle_timeout = "600s"     # Keep connections alive longer
  keepalive = "30s"         # More frequent keepalive
}
```

### Parallel Operations

Use parallel operations for better performance:

```hcl
actions {
  action "mass-update" {
    description = "Update all machines in parallel"
    
    machines = ["web-server", "db-server", "app-server", "cache-server"]
    parallel = true
    max_concurrent = 4      # Limit concurrent operations
    
    command {
      command = "apt update && apt upgrade -y"
      timeout = 300
    }
  }
}
```

### File Transfer Optimization

Optimize file transfers:

```hcl
actions {
  action "fast-deploy" {
    description = "Fast application deployment"
    
    machines = ["web-server"]
    
    file_copy {
      source = "app/"
      destination = "/opt/app/"
      recursive = true
      mode = "sftp"         # Use SFTP for better performance
      compression = true    # Enable compression
      buffer_size = 65536   # Larger buffer for faster transfers
    }
  }
}
```

## Monitoring and Logging

### Enable Debug Logging

Enable detailed SSH logging:

```bash
# Set debug log level
export SPOOKY_LOG_LEVEL=debug

# Run commands with verbose output
spooky actions run my-project --verbose
```

### Monitor Connection Performance

Monitor SSH connection performance:

```bash
# Test connection performance
time spooky machines ping my-project

# Monitor connection pool usage
spooky machines ping my-project --verbose
```

### Log Analysis

Analyze SSH logs for troubleshooting:

```bash
# Check SSH connection logs
grep "ssh" /var/log/auth.log

# Check spooky SSH logs
grep "SSH" spooky.log
```

## Troubleshooting

### Common Issues

1. **Authentication failures**:
   ```bash
   # Test SSH manually
   ssh -i ~/.ssh/id_ed25519 user@hostname
   
   # Check key permissions
   ls -la ~/.ssh/id_ed25519
   ```

2. **Connection timeouts**:
   ```bash
   # Test network connectivity
   ping hostname
   telnet hostname 22
   ```

3. **Host key issues**:
   ```bash
   # Add host key
   ssh-keyscan -H hostname >> ~/.ssh/known_hosts
   ```

### Debug Commands

```bash
# Enable verbose SSH output
ssh -v user@hostname

# Test specific machine
spooky machines ping my-project --machine web-server

# Check machine configuration
spooky machines list my-project --verbose
```

## Examples

### Complete Deployment Example

```hcl
# Project configuration
project {
  name = "web-application"
  description = "Web application deployment"
}

# Machine inventory
machines {
  machine "web-server" {
    hostname = "web.example.com"
    port = 22
    user = "admin"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/id_ed25519"
    }
    
    tags = {
      environment = "production"
      role = "web"
    }
  }
  
  machine "db-server" {
    hostname = "db.example.com"
    port = 22
    user = "dbadmin"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/id_ed25519"
    }
    
    tags = {
      environment = "production"
      role = "database"
    }
  }
}

# Actions for deployment
actions {
  action "deploy-application" {
    description = "Deploy web application"
    
    machines = ["web-server"]
    
    # Stop application
    service_control {
      service = "nginx"
      action = "stop"
      timeout = 30
    }
    
    # Deploy files
    file_copy {
      source = "app/"
      destination = "/opt/app/"
      recursive = true
      permissions = "0755"
    }
    
    # Deploy configuration
    template_deploy {
      source = "templates/nginx.conf.tmpl"
      destination = "/etc/nginx/nginx.conf"
      permissions = "0644"
      variables = {
        server_name = "example.com"
        port = 80
      }
    }
    
    # Start application
    service_control {
      service = "nginx"
      action = "start"
      timeout = 30
    }
    
    # Verify deployment
    command {
      command = "curl -f http://localhost/health"
      timeout = 30
    }
  }
  
  action "backup-database" {
    description = "Backup database"
    
    machines = ["db-server"]
    
    command {
      command = "pg_dump myapp > /backup/myapp_$(date +%Y%m%d_%H%M%S).sql"
      timeout = 300
    }
    
    file_copy {
      source = "/backup/"
      destination = "backups/"
      direction = "download"
      recursive = true
    }
  }
}
```

### Parallel Operations Example

```hcl
actions {
  action "parallel-maintenance" {
    description = "Perform maintenance on all servers"
    
    machines = ["web-server", "db-server", "app-server"]
    parallel = true
    max_concurrent = 3
    
    # Update packages
    command {
      command = "apt update && apt upgrade -y"
      timeout = 300
    }
    
    # Restart services
    service_control {
      service = "nginx"
      action = "restart"
      timeout = 60
    }
    
    # Check system status
    command {
      command = "systemctl status"
      timeout = 30
    }
  }
}
```

## Summary

The SSH system in spooky provides comprehensive functionality for secure, efficient, and reliable SSH operations. Key features include:

1. **Multiple Authentication Methods**: SSH keys, passwords, and agent authentication
2. **Connection Pooling**: Efficient connection management and reuse
3. **File Transfer**: SFTP and SCP file transfer with progress tracking
4. **Command Execution**: Secure command execution with output capture
5. **Service Management**: System service control and management
6. **Template Deployment**: Template-based file deployment with variable substitution
7. **Parallel Operations**: Parallel execution across multiple machines
8. **Error Handling**: Comprehensive error handling and retry logic
9. **Security Features**: Host key validation, key management, and security best practices
10. **Performance Optimization**: Connection pooling, parallel operations, and transfer optimization

The SSH system is production-ready and provides all necessary functionality for secure, efficient, and reliable SSH operations in the spooky automation platform.
