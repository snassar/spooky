# SSH System User Guide

## Overview

The spooky SSH system provides secure, efficient, and validated SSH connectivity for remote machine management. This guide covers everything from basic SSH connections to advanced features like key type validation, certificate authentication, and connection pooling.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Key Types and Authentication](#key-types-and-authentication)
3. [SSH Certificates](#ssh-certificates)
4. [Connection Management](#connection-management)
5. [Command Execution](#command-execution)
6. [Advanced Features](#advanced-features)
7. [Best Practices](#best-practices)
8. [Examples](#examples)

## Getting Started

### Basic SSH Connection

The SSH system in spooky is designed to work seamlessly with the machines inventory system. SSH connections are automatically established when needed for machine operations.

**Basic Machine Configuration with SSH:**
```hcl
machines {
  machine "web-server-01" {
    host = "192.168.1.10"
    user = "admin"
    port = 22
    
    key_file = "~/.ssh/id_ed25519"
    passphrase = "my-secure-passphrase"
    
    tags = ["web", "production"]
    groups = ["web-servers"]
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      owner = "web-team"
    }
  }
}
```

### Testing SSH Connectivity

Use the machines ping command to test SSH connectivity:

```bash
# Test connectivity to all machines
spooky machines ping ./my-project

# Test specific machine
spooky machines ping ./my-project --machine web-server-01

# Test with verbose output
spooky machines ping ./my-project --verbose
```

## Key Types and Authentication

### Supported Key Types

The SSH system supports three key types with strict validation:

#### 1. ED25519 Keys
Modern, secure elliptic curve keys with fixed 256-bit size.

**Generation:**
```bash
# Generate ED25519 key
ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -C "spooky-ssh-key"

# Generate with passphrase
ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -C "spooky-ssh-key" -N "my-passphrase"
```

**Configuration:**
```hcl
machines {
  machine "modern-server" {
    host = "192.168.1.20"
    user = "admin"
    key_file = "~/.ssh/id_ed25519"
    passphrase = "my-passphrase"  # Optional
  }
}
```

#### 2. ED25519-SK Keys
Hardware security key-based ED25519 keys (planned feature).

**Note:** This key type is marked as supported but implementation is pending.

#### 3. RSA Keys (4096-bit minimum)
Traditional RSA keys with enhanced security requirements.

**Generation:**
```bash
# Generate 4096-bit RSA key
ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa_4096 -C "spooky-rsa-key"

# Generate with passphrase
ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa_4096 -C "spooky-rsa-key" -N "my-passphrase"
```

**Configuration:**
```hcl
machines {
  machine "legacy-server" {
    host = "192.168.1.30"
    user = "admin"
    key_file = "~/.ssh/id_rsa_4096"
    passphrase = "my-passphrase"  # Optional
  }
}
```

### Key Validation

The SSH system automatically validates all keys:

- **Type Validation**: Only supported key types are accepted
- **Size Validation**: RSA keys must be 4096-bit minimum
- **Format Validation**: Keys must be in valid SSH format
- **Access Validation**: Key files must be readable

**Validation Errors:**
```bash
# Example: RSA key too small
Error: key validation failed for rsa: RSA key size 2048 bits is less than minimum required 4096 bits

# Example: Unsupported key type
Error: key validation failed for dsa: unsupported key type: ssh-dss. Supported types: ed25519, ed25519-sk, rsa-4096
```

### Authentication Methods

#### Public Key Authentication
Standard SSH key-based authentication (recommended).

```hcl
machines {
  machine "key-auth-server" {
    host = "192.168.1.40"
    user = "admin"
    key_file = "~/.ssh/id_ed25519"
    # No passphrase = unencrypted key
  }
}
```

#### Password Authentication
Traditional password-based authentication (less secure).

```hcl
machines {
  machine "password-auth-server" {
    host = "192.168.1.50"
    user = "admin"
    password = "my-secure-password"  # Not recommended for production
  }
}
```

## SSH Certificates

### Certificate Authentication

SSH certificates provide enhanced security and key management capabilities.

#### Certificate Setup

1. **Generate Certificate Authority (CA):**
```bash
# Generate CA key
ssh-keygen -t ed25519 -f ~/.ssh/ca_key -C "spooky-ca"

# Generate CA certificate
ssh-keygen -s ~/.ssh/ca_key -I "spooky-ca" -n "admin" -V +52w ~/.ssh/id_ed25519.pub
```

2. **Configure Certificate Authentication:**
```hcl
machines {
  machine "cert-server" {
    host = "192.168.1.60"
    user = "admin"
    key_file = "~/.ssh/id_ed25519"           # Private key
    certificate_file = "~/.ssh/id_ed25519-cert.pub"  # Certificate
    passphrase = "my-passphrase"             # Optional
  }
}
```

#### Certificate Benefits

- **Enhanced Security**: Certificates include identity and authorization information
- **Key Management**: Centralized key management through CA
- **Access Control**: Fine-grained access control through certificate principals
- **Audit Trail**: Certificate usage can be tracked and audited

### Certificate Validation

The SSH system validates certificates:

- **Format Validation**: Certificate must be in valid SSH format
- **Private Key Requirement**: Certificate must be accompanied by private key
- **Expiration Checking**: Certificate expiration is validated (planned)
- **Principal Validation**: Certificate principals are checked (planned)

## Connection Management

### Connection Pooling

The SSH system uses connection pooling for efficiency:

```hcl
# Connection pooling is automatic
# Multiple operations reuse connections when possible
```

**Pool Configuration:**
- **Max Connections**: Default 10 concurrent connections
- **Idle Timeout**: Default 300 seconds
- **Connection Timeout**: Default 30 seconds
- **Retry Attempts**: Default 3 attempts

### Connection Health

Monitor connection health through the machines ping command:

```bash
# Check connection health
spooky machines ping ./my-project

# Output shows connection status:
# ✓ web-server-01: Connected (latency: 15ms)
# ✗ db-server-01: Connection failed (timeout)
```

### Connection Retry Logic

The SSH system implements robust retry logic:

1. **Initial Connection**: Attempt to establish connection
2. **Retry on Failure**: Retry up to 3 times with exponential backoff
3. **Timeout Handling**: Configurable timeouts for different operations
4. **Error Reporting**: Detailed error messages for troubleshooting

## Command Execution

### Basic Command Execution

Commands are executed through the actions system:

```hcl
actions {
  action "check-system-info" {
    description = "Check system information on remote machines"
    
    machines = ["web-server-01", "db-server-01"]
    parallel = true
    
    command = "uname -a && df -h && free -h"
  }
}
```

### Command Configuration

Commands support various configuration options:

```hcl
actions {
  action "complex-command" {
    description = "Complex command with environment and working directory"
    
    machines = ["web-server-01"]
    
    command = "python3 /opt/scripts/health_check.py"
    working_dir = "/opt/scripts"
    environment = {
      "PYTHONPATH" = "/opt/lib"
      "LOG_LEVEL" = "DEBUG"
    }
    timeout = "60s"
  }
}
```

### Command Output Handling

Command output is automatically captured and processed:

- **Standard Output**: Captured and available for processing
- **Standard Error**: Captured and reported separately
- **Exit Codes**: Available for error handling
- **Execution Time**: Tracked for performance monitoring

## Advanced Features

### Environment Variables

Set environment variables for command execution:

```hcl
actions {
  action "environment-test" {
    machines = ["web-server-01"]
    command = "echo $CUSTOM_VAR && env | grep CUSTOM"
    
    environment = {
      "CUSTOM_VAR" = "custom-value"
      "DEBUG" = "true"
      "PATH" = "/usr/local/bin:/usr/bin:/bin"
    }
  }
}
```

### Working Directory

Set working directory for command execution:

```hcl
actions {
  action "working-dir-test" {
    machines = ["web-server-01"]
    command = "pwd && ls -la"
    working_dir = "/opt/app"
  }
}
```

### Timeout Configuration

Configure timeouts for different operations:

```hcl
machines {
  machine "slow-server" {
    host = "192.168.1.70"
    user = "admin"
    key_file = "~/.ssh/id_ed25519"
    
    # Connection timeout
    timeout = "60s"
    
    # Keepalive settings
    keepalive_interval = "30s"
    keepalive_count = 3
  }
}
```

### Compression

Enable compression for slow connections:

```hcl
machines {
  machine "slow-connection" {
    host = "remote.example.com"
    user = "admin"
    key_file = "~/.ssh/id_ed25519"
    
    # Enable compression
    compression = true
  }
}
```

## Best Practices

### Key Management

1. **Use ED25519 keys** for new deployments:
   - Modern, secure, and efficient
   - Fixed size (256 bits)
   - Resistant to timing attacks

2. **Use 4096-bit RSA keys** if RSA is required:
   - Minimum security requirement
   - Compatible with legacy systems
   - Larger key size for enhanced security

3. **Secure key storage**:
   ```bash
   # Set proper permissions
   chmod 600 ~/.ssh/id_ed25519
   chmod 644 ~/.ssh/id_ed25519.pub
   
   # Store keys in secure location
   # Use passphrases for additional security
   ```

4. **Key rotation**:
   - Rotate keys regularly (every 90-180 days)
   - Use certificates for easier key management
   - Monitor key usage and access

### Connection Management

1. **Use connection pooling**:
   - Multiple operations reuse connections
   - Reduces connection overhead
   - Improves performance

2. **Set appropriate timeouts**:
   - Connection timeout: 30-60 seconds
   - Command timeout: Based on command complexity
   - Idle timeout: 300-600 seconds

3. **Implement retry logic**:
   - Handle transient network issues
   - Exponential backoff for retries
   - Maximum retry attempts (3-5)

4. **Monitor connection health**:
   - Regular connectivity testing
   - Performance monitoring
   - Error tracking and alerting

### Security Considerations

1. **Key validation**:
   - Always validate keys before use
   - Use only supported key types
   - Verify key permissions and ownership

2. **Certificate usage**:
   - Use certificates for enhanced security
   - Implement proper CA management
   - Monitor certificate expiration

3. **Access control**:
   - Follow least privilege principles
   - Use dedicated service accounts
   - Implement proper user management

4. **Audit and monitoring**:
   - Log all SSH connections
   - Monitor for suspicious activity
   - Track key and certificate usage

### Performance Optimization

1. **Connection pooling**:
   - Reuse connections when possible
   - Configure appropriate pool sizes
   - Monitor pool utilization

2. **Parallel execution**:
   - Use parallel execution for multiple machines
   - Configure appropriate concurrency limits
   - Monitor resource usage

3. **Compression**:
   - Enable compression for slow connections
   - Monitor compression effectiveness
   - Balance compression vs. CPU usage

4. **Caching**:
   - Cache connection information
   - Cache key validation results
   - Implement appropriate cache invalidation

## Examples

### Basic SSH Setup

**Project Structure:**
```
my-ssh-project/
├── project.hcl
├── machines.hcl
└── actions.hcl
```

**Project Configuration (`project.hcl`):**
```hcl
project {
  name = "ssh-example-project"
  description = "Example SSH project with different key types"
  
  metadata {
    version = "1.0.0"
    author = "admin"
    tags = ["ssh", "example"]
  }
}
```

**Machine Inventory (`machines.hcl`):**
```hcl
machines {
  # ED25519 key example
  machine "modern-server" {
    host = "192.168.1.10"
    user = "admin"
    key_file = "~/.ssh/id_ed25519"
    passphrase = "my-passphrase"
    
    tags = ["modern", "production"]
    groups = ["web-servers"]
    
    metadata {
      environment = "production"
      key_type = "ed25519"
    }
  }
  
  # RSA 4096-bit key example
  machine "legacy-server" {
    host = "192.168.1.20"
    user = "admin"
    key_file = "~/.ssh/id_rsa_4096"
    
    tags = ["legacy", "production"]
    groups = ["database-servers"]
    
    metadata {
      environment = "production"
      key_type = "rsa-4096"
    }
  }
  
  # Certificate authentication example
  machine "cert-server" {
    host = "192.168.1.30"
    user = "admin"
    key_file = "~/.ssh/id_ed25519"
    certificate_file = "~/.ssh/id_ed25519-cert.pub"
    passphrase = "my-passphrase"
    
    tags = ["certificate", "production"]
    groups = ["app-servers"]
    
    metadata {
      environment = "production"
      auth_method = "certificate"
    }
  }
}
```

**Actions Configuration (`actions.hcl`):**
```hcl
actions {
  # Basic system information
  action "system-info" {
    description = "Get system information from all machines"
    
    machines = ["modern-server", "legacy-server", "cert-server"]
    parallel = true
    
    command = "uname -a && hostname && date"
  }
  
  # Disk usage check
  action "disk-usage" {
    description = "Check disk usage on all machines"
    
    machines = ["modern-server", "legacy-server", "cert-server"]
    parallel = true
    
    command = "df -h"
  }
  
  # Process monitoring
  action "process-check" {
    description = "Check running processes"
    
    machines = ["modern-server", "legacy-server", "cert-server"]
    parallel = true
    
    command = "ps aux | head -20"
  }
}
```

### Testing the Setup

```bash
# Validate the project
spooky validate --project ./my-ssh-project

# Test SSH connectivity
spooky machines ping ./my-ssh-project

# Run actions
spooky actions run ./my-ssh-project --action system-info
spooky actions run ./my-ssh-project --action disk-usage
spooky actions run ./my-ssh-project --action process-check
```

### Advanced Configuration

**Complex Machine Configuration:**
```hcl
machines {
  machine "production-web" {
    host = "web.prod.example.com"
    user = "webadmin"
    port = 2222  # Custom SSH port
    
    key_file = "~/.ssh/prod_web_key"
    passphrase = "production-passphrase"
    
    # Connection settings
    timeout = "60s"
    keepalive_interval = "30s"
    keepalive_count = 3
    compression = true
    
    # Security settings
    strict_host_key_check = true
    known_hosts_path = "~/.ssh/known_hosts"
    
    tags = ["web", "production", "load-balanced"]
    groups = ["web-servers", "production-servers"]
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      availability_zone = "us-west-1a"
      instance_type = "t3.large"
      cost_center = "IT-001"
      owner = "web-team"
      key_type = "ed25519"
      last_updated = "2024-01-15"
    }
  }
}
```

**Complex Action Configuration:**
```hcl
actions {
  action "deploy-web-application" {
    description = "Deploy web application to production servers"
    
    machines = ["production-web"]
    parallel = false  # Sequential deployment
    
    # Environment setup
    environment = {
      "DEPLOY_ENV" = "production"
      "APP_VERSION" = "1.2.3"
      "LOG_LEVEL" = "INFO"
      "NODE_ENV" = "production"
    }
    
    working_dir = "/opt/webapp"
    
    # Multi-step deployment
    command = """
    echo "Starting deployment of version $APP_VERSION"
    git pull origin main
    npm install --production
    npm run build
    sudo systemctl restart webapp
    echo "Deployment completed successfully"
    """
    
    timeout = "300s"  # 5 minutes
    
    tags = ["deployment", "web", "production"]
  }
}
```

This comprehensive user guide provides everything needed to effectively use the SSH system in spooky, from basic setup to advanced configurations and best practices.
