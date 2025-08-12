# SSH System Troubleshooting Guide

## Overview

This troubleshooting guide provides solutions for common issues encountered when working with the spooky SSH system. It covers error messages, key validation problems, certificate issues, connection problems, and performance issues.

## Table of Contents

1. [Common Error Messages](#common-error-messages)
2. [Key Validation Issues](#key-validation-issues)
3. [Certificate Problems](#certificate-problems)
4. [Connection Issues](#connection-issues)
5. [Authentication Failures](#authentication-failures)
6. [Performance Issues](#performance-issues)
7. [Configuration Problems](#configuration-problems)
8. [Debugging Techniques](#debugging-techniques)
9. [Best Practices for Troubleshooting](#best-practices-for-troubleshooting)

## Common Error Messages

### Key Validation Errors

#### "key validation failed for rsa: RSA key size 2048 bits is less than minimum required 4096 bits"

**Cause:** RSA key is smaller than the required 4096-bit minimum.

**Solution:**
```bash
# Generate a new 4096-bit RSA key
ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa_4096 -C "spooky-rsa-key"

# Update machine configuration
# In machines.hcl:
machine "server" {
  host = "192.168.1.10"
  user = "admin"
  key_file = "~/.ssh/id_rsa_4096"  # Use the new 4096-bit key
}
```

**Prevention:**
- Always use 4096-bit RSA keys for new deployments
- Consider using ED25519 keys instead (more secure and efficient)

#### "key validation failed for dsa: unsupported key type: ssh-dss. Supported types: ed25519, ed25519-sk, rsa-4096"

**Cause:** Using an unsupported key type (DSA keys are not supported).

**Solution:**
```bash
# Generate a supported key type (ED25519 recommended)
ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -C "spooky-ed25519-key"

# Or generate 4096-bit RSA key
ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa_4096 -C "spooky-rsa-key"

# Update machine configuration
machine "server" {
  host = "192.168.1.10"
  user = "admin"
  key_file = "~/.ssh/id_ed25519"  # Use supported key type
}
```

**Supported Key Types:**
- ED25519 (recommended)
- ED25519-SK (hardware security key, planned)
- RSA 4096-bit (minimum)

#### "failed to read private key file: open ~/.ssh/id_rsa: no such file or directory"

**Cause:** Private key file doesn't exist or path is incorrect.

**Solution:**
```bash
# Check if key file exists
ls -la ~/.ssh/id_rsa

# Generate key if it doesn't exist
ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519

# Verify the path in configuration
# In machines.hcl, ensure the path is correct:
machine "server" {
  host = "192.168.1.10"
  user = "admin"
  key_file = "~/.ssh/id_ed25519"  # Verify this path
}
```

**Path Resolution:**
- `~` expands to user's home directory
- Use absolute paths if needed: `/home/user/.ssh/id_ed25519`
- Check file permissions: `chmod 600 ~/.ssh/id_ed25519`

### Connection Errors

#### "connection failed: dial tcp 192.168.1.10:22: connect: connection refused"

**Cause:** SSH service not running or port not accessible.

**Solution:**
```bash
# Check if SSH service is running on target machine
ssh user@192.168.1.10 "sudo systemctl status sshd"

# Start SSH service if not running
ssh user@192.168.1.10 "sudo systemctl start sshd"

# Check if port 22 is open
nmap -p 22 192.168.1.10

# Check firewall settings
ssh user@192.168.1.10 "sudo ufw status"
```

**Troubleshooting Steps:**
1. Verify SSH service is running
2. Check firewall settings
3. Verify port configuration
4. Test basic connectivity

#### "connection failed: dial tcp 192.168.1.10:22: i/o timeout"

**Cause:** Network connectivity issues or firewall blocking connection.

**Solution:**
```bash
# Test basic connectivity
ping 192.168.1.10

# Test port connectivity
telnet 192.168.1.10 22

# Check network routing
traceroute 192.168.1.10

# Verify firewall rules
sudo iptables -L | grep 22
```

**Network Troubleshooting:**
1. Check physical network connectivity
2. Verify routing configuration
3. Check firewall rules
4. Test with different network paths

#### "connection failed: ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey]"

**Cause:** Authentication failed - no valid authentication method available.

**Solution:**
```bash
# Check if public key is installed on target machine
ssh user@192.168.1.10 "cat ~/.ssh/authorized_keys"

# Install public key if missing
ssh-copy-id -i ~/.ssh/id_ed25519.pub user@192.168.1.10

# Or manually copy the key
cat ~/.ssh/id_ed25519.pub | ssh user@192.168.1.10 "mkdir -p ~/.ssh && cat >> ~/.ssh/authorized_keys"

# Set proper permissions
ssh user@192.168.1.10 "chmod 700 ~/.ssh && chmod 600 ~/.ssh/authorized_keys"
```

**Authentication Checklist:**
1. Private key exists and is readable
2. Public key is in `~/.ssh/authorized_keys` on target
3. File permissions are correct (700 for .ssh, 600 for authorized_keys)
4. SSH service allows public key authentication

### Certificate Errors

#### "failed to load SSH certificate: failed to parse SSH certificate: ssh: no key found"

**Cause:** Certificate file is not in valid SSH certificate format.

**Solution:**
```bash
# Check certificate format
cat ~/.ssh/id_ed25519-cert.pub

# Certificate should start with "ssh-ed25519-cert-v01@openssh.com"
# If not, regenerate the certificate
ssh-keygen -s ~/.ssh/ca_key -I "spooky-cert" -n "admin" -V +52w ~/.ssh/id_ed25519.pub
```

**Certificate Format Validation:**
- Certificate must be in OpenSSH certificate format
- Must include certificate header and signature
- Must be accompanied by private key

#### "failed to create certificate signer: ssh: certificate and private key do not match"

**Cause:** Certificate and private key are not a matching pair.

**Solution:**
```bash
# Verify certificate and key match
ssh-keygen -L -f ~/.ssh/id_ed25519-cert.pub

# Check the public key fingerprint
ssh-keygen -lf ~/.ssh/id_ed25519.pub

# Regenerate certificate with correct key
ssh-keygen -s ~/.ssh/ca_key -I "spooky-cert" -n "admin" -V +52w ~/.ssh/id_ed25519.pub
```

**Certificate-Key Matching:**
1. Certificate must be signed for the specific public key
2. Private key must correspond to the public key in the certificate
3. Certificate and key must be from the same key pair

## Key Validation Issues

### RSA Key Size Problems

#### Problem: RSA key too small
```bash
Error: key validation failed for rsa: RSA key size 2048 bits is less than minimum required 4096 bits
```

**Diagnosis:**
```bash
# Check current key size
ssh-keygen -lf ~/.ssh/id_rsa

# Look for key size in output (e.g., "2048 SHA256:...")
```

**Solution:**
```bash
# Generate new 4096-bit RSA key
ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa_4096 -C "spooky-4096-key"

# Update configuration
# In machines.hcl:
machine "server" {
  host = "192.168.1.10"
  user = "admin"
  key_file = "~/.ssh/id_rsa_4096"
}
```

### Unsupported Key Types

#### Problem: DSA, ECDSA, or other unsupported key types
```bash
Error: key validation failed for dsa: unsupported key type: ssh-dss
```

**Solution:**
```bash
# Generate supported key type
ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -C "spooky-ed25519-key"

# Update all machine configurations to use supported keys
```

**Migration Strategy:**
1. Generate new supported keys
2. Install public keys on target machines
3. Update machine configurations
4. Test connections
5. Remove old keys

### Key File Access Issues

#### Problem: Key file not readable
```bash
Error: failed to read private key file: permission denied
```

**Solution:**
```bash
# Check file permissions
ls -la ~/.ssh/id_ed25519

# Set correct permissions
chmod 600 ~/.ssh/id_ed25519
chmod 644 ~/.ssh/id_ed25519.pub

# Check directory permissions
chmod 700 ~/.ssh
```

**Permission Requirements:**
- Private key: 600 (owner read/write only)
- Public key: 644 (owner read/write, others read)
- .ssh directory: 700 (owner read/write/run only)

## Certificate Problems

### Certificate Format Issues

#### Problem: Invalid certificate format
```bash
Error: failed to parse SSH certificate: ssh: no key found
```

**Diagnosis:**
```bash
# Check certificate format
head -1 ~/.ssh/id_ed25519-cert.pub

# Should show: ssh-ed25519-cert-v01@openssh.com AAAAC3...
```

**Solution:**
```bash
# Regenerate certificate with proper format
ssh-keygen -s ~/.ssh/ca_key -I "spooky-cert" -n "admin" -V +52w ~/.ssh/id_ed25519.pub

# Verify certificate
ssh-keygen -L -f ~/.ssh/id_ed25519-cert.pub
```

### Certificate Expiration

#### Problem: Certificate expired
```bash
Error: certificate expired at 2024-01-01 12:00:00 +0000 UTC
```

**Solution:**
```bash
# Check certificate expiration
ssh-keygen -L -f ~/.ssh/id_ed25519-cert.pub

# Regenerate certificate with new expiration
ssh-keygen -s ~/.ssh/ca_key -I "spooky-cert" -n "admin" -V +52w ~/.ssh/id_ed25519.pub
```

**Certificate Management:**
1. Monitor certificate expiration dates
2. Set up automated renewal processes
3. Use appropriate validity periods
4. Keep CA key secure

### Certificate-Private Key Mismatch

#### Problem: Certificate doesn't match private key
```bash
Error: certificate and private key do not match
```

**Diagnosis:**
```bash
# Check certificate public key
ssh-keygen -L -f ~/.ssh/id_ed25519-cert.pub | grep "Public key"

# Check private key public key
ssh-keygen -y -f ~/.ssh/id_ed25519

# Compare the two outputs
```

**Solution:**
```bash
# Regenerate certificate for the correct key
ssh-keygen -s ~/.ssh/ca_key -I "spooky-cert" -n "admin" -V +52w ~/.ssh/id_ed25519.pub
```

## Connection Issues

### DNS Resolution Problems

#### Problem: Hostname cannot be resolved
```bash
Error: connection failed: dial tcp: lookup example.com: no such host
```

**Solution:**
```bash
# Test DNS resolution
nslookup example.com
dig example.com

# Check /etc/hosts file
cat /etc/hosts | grep example.com

# Use IP address instead of hostname
# In machines.hcl:
machine "server" {
  host = "192.168.1.10"  # Use IP instead of hostname
  user = "admin"
}
```

### Port Configuration Issues

#### Problem: Wrong SSH port
```bash
Error: connection failed: dial tcp 192.168.1.10:22: connection refused
```

**Solution:**
```bash
# Check if SSH is running on different port
nmap -p 22,2222,2022 192.168.1.10

# Update configuration with correct port
# In machines.hcl:
machine "server" {
  host = "192.168.1.10"
  port = 2222  # Use correct port
  user = "admin"
}
```

### Firewall Issues

#### Problem: Firewall blocking SSH
```bash
Error: connection failed: dial tcp 192.168.1.10:22: i/o timeout
```

**Solution:**
```bash
# Check local firewall
sudo ufw status
sudo iptables -L

# Check remote firewall
ssh user@192.168.1.10 "sudo ufw status"
ssh user@192.168.1.10 "sudo iptables -L"

# Allow SSH through firewall
sudo ufw allow ssh
ssh user@192.168.1.10 "sudo ufw allow ssh"
```

## Authentication Failures

### Public Key Authentication Issues

#### Problem: Public key not in authorized_keys
```bash
Error: ssh: handshake failed: ssh: unable to authenticate
```

**Solution:**
```bash
# Check authorized_keys on target machine
ssh user@192.168.1.10 "cat ~/.ssh/authorized_keys"

# Install public key
ssh-copy-id -i ~/.ssh/id_ed25519.pub user@192.168.1.10

# Or manually copy
cat ~/.ssh/id_ed25519.pub | ssh user@192.168.1.10 "mkdir -p ~/.ssh && cat >> ~/.ssh/authorized_keys"
```

### Passphrase Issues

#### Problem: Incorrect passphrase
```bash
Error: failed to parse private key: ssh: this private key is passphrase protected
```

**Solution:**
```bash
# Test passphrase manually
ssh-keygen -y -f ~/.ssh/id_ed25519

# Update configuration with correct passphrase
# In machines.hcl:
machine "server" {
  host = "192.168.1.10"
  user = "admin"
  key_file = "~/.ssh/id_ed25519"
  passphrase = "correct-passphrase"  # Ensure this is correct
}
```

### User Permission Issues

#### Problem: User doesn't have SSH access
```bash
Error: ssh: handshake failed: ssh: unable to authenticate
```

**Solution:**
```bash
# Check user permissions on target machine
ssh user@192.168.1.10 "ls -la ~/.ssh/"

# Fix permissions
ssh user@192.168.1.10 "chmod 700 ~/.ssh && chmod 600 ~/.ssh/authorized_keys"

# Check SSH configuration
ssh user@192.168.1.10 "sudo grep -i pubkey /etc/ssh/sshd_config"
```

## Performance Issues

### Connection Pooling Problems

#### Problem: Too many connections
```bash
Error: connection pool at capacity (10)
```

**Solution:**
```bash
# Increase connection pool size
# In configuration:
config {
  max_connections = 20  # Increase from default 10
}
```

**Optimization:**
1. Increase pool size for high-concurrency workloads
2. Implement connection cleanup
3. Monitor pool utilization
4. Use connection timeouts

### Slow Connection Performance

#### Problem: Slow SSH connections
```bash
# Connection takes >30 seconds
```

**Solution:**
```bash
# Enable compression for slow connections
# In machines.hcl:
machine "slow-server" {
  host = "remote.example.com"
  user = "admin"
  key_file = "~/.ssh/id_ed25519"
  compression = true  # Enable compression
}

# Optimize connection settings
machine "optimized-server" {
  host = "192.168.1.10"
  user = "admin"
  key_file = "~/.ssh/id_ed25519"
  timeout = "60s"  # Increase timeout
  keepalive_interval = "30s"
  keepalive_count = 3
}
```

### Command Execution Timeouts

#### Problem: Commands timing out
```bash
Error: command run timeout after 30s
```

**Solution:**
```bash
# Increase command timeout
# In actions.hcl:
action "long-running-command" {
  machines = ["server"]
  command = "long-running-script.sh"
  timeout = "300s"  # 5 minutes
}
```

## Configuration Problems

### HCL Syntax Errors

#### Problem: Invalid HCL syntax
```bash
Error: failed to parse machines.hcl: unexpected block type
```

**Solution:**
```bash
# Check HCL syntax
spooky validate --project ./my-project

# Common syntax fixes:
machines {
  machine "server" {
    host = "192.168.1.10"  # Use quotes for strings
    user = "admin"
    key_file = "~/.ssh/id_ed25519"
    
    tags = ["web", "production"]  # Use array syntax
    
    metadata {
      environment = "production"
      owner = "web-team"
    }
  }
}
```

### Missing Required Fields

#### Problem: Missing required configuration
```bash
Error: validation failed: host is required
```

**Solution:**
```bash
# Ensure all required fields are present
machine "server" {
  host = "192.168.1.10"  # Required
  user = "admin"         # Required
  key_file = "~/.ssh/id_ed25519"  # Required for key auth
}
```

**Required Fields:**
- `host`: Target hostname or IP
- `user`: SSH username
- `key_file` or `password`: Authentication method

## Debugging Techniques

### Enable Verbose Logging

```bash
# Enable verbose output for SSH operations
spooky machines ping ./my-project --verbose

# Check logs for detailed error information
spooky machines ping ./my-project --log-level debug
```

### Test Individual Components

```bash
# Test key validation
ssh-keygen -lf ~/.ssh/id_ed25519

# Test SSH connection manually
ssh -i ~/.ssh/id_ed25519 user@192.168.1.10

# Test certificate
ssh-keygen -L -f ~/.ssh/id_ed25519-cert.pub
```

### Network Diagnostics

```bash
# Test basic connectivity
ping 192.168.1.10

# Test port connectivity
telnet 192.168.1.10 22

# Check routing
traceroute 192.168.1.10

# Check DNS resolution
nslookup example.com
```

### SSH Debug Mode

```bash
# Enable SSH debug mode for detailed connection information
ssh -v -i ~/.ssh/id_ed25519 user@192.168.1.10

# Very verbose mode
ssh -vv -i ~/.ssh/id_ed25519 user@192.168.1.10

# Maximum verbosity
ssh -vvv -i ~/.ssh/id_ed25519 user@192.168.1.10
```

## Best Practices for Troubleshooting

### Systematic Approach

1. **Identify the Problem**: Understand the exact error message
2. **Check Configuration**: Verify all configuration files
3. **Test Connectivity**: Ensure basic network connectivity
4. **Validate Authentication**: Test authentication methods
5. **Check Permissions**: Verify file and directory permissions
6. **Review Logs**: Check system and application logs
7. **Test Incrementally**: Test each component separately

### Common Checklist

#### Before Troubleshooting
- [ ] Verify network connectivity
- [ ] Check SSH service status
- [ ] Validate key file existence and permissions
- [ ] Confirm user permissions
- [ ] Check firewall settings

#### During Troubleshooting
- [ ] Use verbose logging
- [ ] Test with manual SSH commands
- [ ] Validate configuration syntax
- [ ] Check error messages carefully
- [ ] Test with different authentication methods

#### After Resolution
- [ ] Document the solution
- [ ] Update configuration if needed
- [ ] Test the fix thoroughly
- [ ] Monitor for recurrence
- [ ] Update documentation

### Prevention Strategies

1. **Key Management**:
   - Use supported key types only
   - Implement key rotation policies
   - Monitor key expiration
   - Use certificates for enhanced security

2. **Configuration Management**:
   - Validate configurations before deployment
   - Use consistent naming conventions
   - Document configuration changes
   - Implement configuration testing

3. **Monitoring and Alerting**:
   - Monitor SSH connection success rates
   - Track authentication failures
   - Alert on unusual patterns
   - Regular connectivity testing

4. **Security Best Practices**:
   - Use strong key types (ED25519 preferred)
   - Implement proper file permissions
   - Use certificates for enhanced security
   - Regular security audits

This comprehensive troubleshooting guide provides solutions for the most common SSH system issues and helps maintain reliable SSH connectivity in spooky deployments.
