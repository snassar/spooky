# SSH System User Guide

## Overview

The spooky SSH system provides comprehensive SSH connectivity and management capabilities with enhanced key support, SSH certificate support, and robust connection management. This guide covers everything from basic SSH configuration to advanced features like certificate authentication and connection pooling.

## Getting Started

### Prerequisites

- spooky CLI installed and configured
- SSH keys (ED25519, ED25519-SK, or RSA 4096-bit minimum)
- SSH access to target machines
- Basic understanding of SSH authentication

### Quick Start

1. **Check SSH System Status**
   ```bash
   spooky ssh --help
   ```

2. **Test SSH Connection**
   ```bash
   spooky ssh connect example.com --user admin --key ~/.ssh/id_ed25519
   ```

3. **Validate SSH Key**
   ```bash
   spooky ssh validate-key ~/.ssh/id_ed25519
   ```

## SSH System Concepts

### What is the SSH System?

The SSH system provides:

- **Secure Connections**: Encrypted SSH connections to remote machines
- **Key Management**: Support for modern SSH key types (ED25519, RSA 4096-bit)
- **Certificate Authentication**: SSH certificate support for enhanced security
- **Connection Pooling**: Efficient connection management and reuse
- **Command Execution**: Remote command execution via SSH
- **File Transfer**: Secure file transfer capabilities (SFTP/SCP)

### Supported Key Types

The SSH system supports the following key types:

1. **ED25519 Keys** - Modern, secure elliptic curve keys
   - Fixed 256-bit size
   - Always valid (no size validation needed)
   - Recommended for new deployments

2. **ED25519-SK Keys** - Hardware security key support
   - Security key-based ED25519 keys
   - Hardware-backed security
   - Implementation pending

3. **RSA Keys** - Traditional RSA keys with enhanced security
   - Minimum 4096-bit key size enforced
   - Compatible with legacy systems
   - Must meet minimum size requirements

### SSH Certificate Support

The system includes comprehensive SSH certificate support:

- **Certificate Authentication**: Use SSH certificates for authentication
- **Private Key Requirement**: Certificates must be accompanied by private keys
- **Passphrase Support**: Encrypted private keys supported
- **Certificate Validation**: Automatic certificate format validation

## Configuration

### SSH Key Configuration

#### ED25519 Key Setup
```bash
# Generate ED25519 key
ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -C "your-email@example.com"

# Set proper permissions
chmod 600 ~/.ssh/id_ed25519
chmod 644 ~/.ssh/id_ed25519.pub

# Add to SSH agent
ssh-add ~/.ssh/id_ed25519
```

#### RSA 4096-bit Key Setup
```bash
# Generate RSA 4096-bit key
ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa_4096 -C "your-email@example.com"

# Set proper permissions
chmod 600 ~/.ssh/id_rsa_4096
chmod 644 ~/.ssh/id_rsa_4096.pub

# Add to SSH agent
ssh-add ~/.ssh/id_rsa_4096
```

#### SSH Certificate Setup
```bash
# Generate certificate (requires CA)
ssh-keygen -s /path/to/ca_key -I "user-cert" -n "user" -V +1d ~/.ssh/id_ed25519.pub

# Certificate will be saved as ~/.ssh/id_ed25519-cert.pub
# Use with private key for authentication
```

### Machine Configuration

SSH settings are configured in machine inventory files:

```hcl
# machines.hcl
machines {
  machine "web-server" {
    hostname = "web.example.com"
    host     = "192.168.1.100"
    port     = 22
    user     = "admin"
    
    # SSH key authentication
    key_file = "~/.ssh/id_ed25519"
    
    # SSH certificate authentication (optional)
    certificate_path = "~/.ssh/id_ed25519-cert.pub"
    passphrase       = "your-passphrase"  # Optional
    
    # Connection settings
    connection_timeout = 30
    command_timeout    = 300
    max_connections    = 10
    retry_attempts     = 3
    retry_delay        = 5
    
    tags = ["web", "production"]
  }
}
```

## CLI Commands

### SSH Connection Commands

#### Test SSH Connection
```bash
# Basic connection test
spooky ssh connect example.com --user admin --key ~/.ssh/id_ed25519

# Connection with certificate
spooky ssh connect example.com \
  --user admin \
  --key ~/.ssh/id_ed25519 \
  --certificate ~/.ssh/id_ed25519-cert.pub

# Connection with custom timeout
spooky ssh connect example.com \
  --user admin \
  --key ~/.ssh/id_ed25519 \
  --timeout 60
```

#### Run SSH Command
```bash
# Run single command
spooky ssh run example.com \
  --user admin \
  --key ~/.ssh/id_ed25519 \
  --command "uname -a"

# Run command with arguments
spooky ssh run example.com \
  --user admin \
  --key ~/.ssh/id_ed25519 \
  --command "ls" \
  --args "-la" \
  --args "/etc"

# Run command with environment variables
spooky ssh run example.com \
  --user admin \
  --key ~/.ssh/id_ed25519 \
  --command "echo \$ENV_VAR" \
  --env "ENV_VAR=value"
```

### SSH Key Management Commands

#### Validate SSH Key
```bash
# Validate key type and format
spooky ssh validate-key ~/.ssh/id_ed25519

# Validate RSA key size
spooky ssh validate-key ~/.ssh/id_rsa_4096

# Validate with verbose output
spooky ssh validate-key ~/.ssh/id_ed25519 --verbose
```

#### Generate Key Fingerprint
```bash
# Generate SHA256 fingerprint
spooky ssh fingerprint ~/.ssh/id_ed25519

# Generate fingerprint for certificate
spooky ssh fingerprint ~/.ssh/id_ed25519-cert.pub
```

### SSH Configuration Commands

#### List SSH Configuration
```bash
# List SSH configuration
spooky ssh config list

# Show specific configuration
spooky ssh config show --key-path ~/.ssh/id_ed25519
```

#### Test SSH Configuration
```bash
# Test SSH configuration
spooky ssh config test

# Test with specific machine
spooky ssh config test --machine web-server
```

## Advanced Features

### Connection Pooling

The SSH system implements efficient connection pooling:

- **Connection Reuse**: Existing connections are reused when possible
- **Health Checks**: Connections are tested before reuse
- **Automatic Cleanup**: Dead connections are automatically removed
- **Thread Safety**: Pool operations are thread-safe

### Certificate Authentication

SSH certificates provide enhanced security:

```bash
# Connect using certificate authentication
spooky ssh connect example.com \
  --user admin \
  --key ~/.ssh/id_ed25519 \
  --certificate ~/.ssh/id_ed25519-cert.pub \
  --passphrase "your-passphrase"
```

### Multi-Machine Operations

Execute commands across multiple machines:

```bash
# Run command on multiple machines
spooky ssh run-multiple \
  --machines "web-server,app-server,db-server" \
  --user admin \
  --key ~/.ssh/id_ed25519 \
  --command "uptime"

# Run with parallel execution
spooky ssh run-multiple \
  --machines "web-server,app-server,db-server" \
  --user admin \
  --key ~/.ssh/id_ed25519 \
  --command "systemctl status nginx" \
  --parallel
```

## Integration with Other Systems

### Facts System Integration

The SSH system integrates with the facts system for remote data collection:

```bash
# Collect facts using SSH
spooky facts export ./my-project \
  --machines "web-server,app-server" \
  --format json \
  --output facts.json
```

### Actions System Integration

The SSH system powers the actions system for remote command execution:

```bash
# Run actions using SSH
spooky actions run ./my-project \
  --machines "web-server" \
  --action "deploy-application"
```

### Machines System Integration

The SSH system validates machine connectivity:

```bash
# Test machine connectivity
spooky machines ping ./my-project \
  --machines "web-server,app-server"
```

## Troubleshooting

### Common Issues

#### Key Validation Errors
```bash
# Error: Unsupported key type
# Solution: Use ED25519, ED25519-SK, or RSA 4096-bit keys

# Error: RSA key size too small
# Solution: Generate new RSA key with 4096-bit minimum
ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa_4096
```

#### Connection Errors
```bash
# Error: Connection timeout
# Solution: Check network connectivity and firewall settings

# Error: Authentication failed
# Solution: Verify key permissions and server configuration
chmod 600 ~/.ssh/id_ed25519
```

#### Certificate Errors
```bash
# Error: Certificate requires private key
# Solution: Ensure private key is provided with certificate

# Error: Invalid certificate format
# Solution: Verify certificate is in valid SSH format
```

### Debugging Commands

#### Verbose SSH Output
```bash
# Enable verbose output
spooky ssh connect example.com \
  --user admin \
  --key ~/.ssh/id_ed25519 \
  --verbose
```

#### SSH Configuration Test
```bash
# Test SSH configuration
spooky ssh config test --verbose

# Test specific connection
spooky ssh config test \
  --host example.com \
  --user admin \
  --key ~/.ssh/id_ed25519 \
  --verbose
```

## Security Best Practices

### Key Management

1. **Use Strong Keys**: Prefer ED25519 keys over RSA
2. **Enforce Key Size**: Use 4096-bit minimum for RSA keys
3. **Secure Permissions**: Set 600 permissions on private keys
4. **Key Rotation**: Regularly rotate SSH keys
5. **Key Storage**: Store keys securely, not in version control

### Certificate Management

1. **Certificate Expiration**: Monitor certificate expiration dates
2. **Private Key Security**: Keep private keys secure and separate
3. **Certificate Validation**: Validate certificates before use
4. **CA Management**: Secure certificate authority keys

### Connection Security

1. **Host Key Verification**: Verify host keys (TODO: implement)
2. **Connection Timeouts**: Use appropriate timeouts
3. **Retry Limits**: Limit retry attempts to prevent brute force
4. **Logging**: Monitor SSH connection logs

## Performance Optimization

### Connection Pooling

- **Pool Size**: Configure appropriate pool size for your workload
- **Idle Timeout**: Set idle timeout to free unused connections
- **Health Checks**: Enable health checks for connection validation

### Parallel Operations

- **Concurrent Connections**: Use parallel execution for multiple machines
- **Connection Limits**: Respect server connection limits
- **Resource Management**: Monitor connection resource usage

## Future Enhancements

### Planned Features

1. **ED25519-SK Support**: Hardware security key integration
2. **Enhanced Certificate Support**: Certificate chain validation
3. **Host Key Verification**: Known hosts file integration
4. **Performance Optimizations**: Connection pooling improvements

### Roadmap

- **Q1 2024**: ED25519-SK hardware key support
- **Q2 2024**: Enhanced certificate validation
- **Q3 2024**: Host key verification implementation
- **Q4 2024**: Performance optimization and caching

## Conclusion

The SSH system provides robust, secure connectivity with support for modern key types and SSH certificates. The system is production-ready and integrates seamlessly with other spooky systems. Follow security best practices and monitor system performance for optimal operation.
