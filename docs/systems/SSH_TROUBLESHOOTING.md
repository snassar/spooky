# SSH Troubleshooting Guide

## Overview

This guide provides comprehensive troubleshooting information for SSH-related issues in the spooky system. It covers common problems, diagnostic techniques, and solutions for SSH connections, authentication, and file transfer operations.

**Status: Implemented** - The SSH system provides comprehensive functionality with robust error handling and diagnostic capabilities.

## Common Issues and Solutions

### Authentication Failures

#### SSH Key Authentication Issues

**Problem**: SSH key authentication fails with "Permission denied" or "Authentication failed"

**Diagnostic Steps**:
1. **Check SSH key permissions**:
   ```bash
   ls -la ~/.ssh/id_rsa
   # Should show: -rw------- (600 permissions)
   ```

2. **Verify SSH key format**:
   ```bash
   ssh-keygen -l -f ~/.ssh/id_rsa
   # Should show key fingerprint
   ```

3. **Test SSH key manually**:
   ```bash
   ssh -i ~/.ssh/id_rsa user@hostname
   ```

**Solutions**:
- **Fix permissions**: `chmod 600 ~/.ssh/id_rsa`
- **Regenerate key**: Use Ed25519 or RSA 4096-bit keys
- **Check passphrase**: Verify passphrase for encrypted keys
- **Update authorized_keys**: Ensure public key is in remote `~/.ssh/authorized_keys`

**Code Example**:
```go
// Validate SSH key before use
func validateSSHKey(keyPath string) error {
    info, err := os.Stat(keyPath)
    if err != nil {
        return fmt.Errorf("key file not found: %w", err)
    }
    
    mode := info.Mode()
    if mode&0077 != 0 {
        return fmt.Errorf("key file has too permissive permissions: %v", mode)
    }
    
    return nil
}
```

#### Password Authentication Issues

**Problem**: Password authentication fails

**Diagnostic Steps**:
1. **Check user credentials**:
   ```bash
   ssh user@hostname
   # Test manual password authentication
   ```

2. **Verify SSH configuration**:
   ```bash
   # Check if password authentication is enabled
   ssh -o PreferredAuthentications=password user@hostname
   ```

**Solutions**:
- **Verify username**: Ensure correct username is specified
- **Check password**: Verify password is correct
- **Enable password auth**: Ensure `PasswordAuthentication yes` in SSH config
- **Use SSH keys**: Prefer SSH key authentication over passwords

### Connection Issues

#### Connection Timeouts

**Problem**: SSH connections timeout or hang

**Diagnostic Steps**:
1. **Test network connectivity**:
   ```bash
   ping hostname
   telnet hostname 22
   ```

2. **Check SSH service**:
   ```bash
   ssh -v user@hostname
   # Verbose output shows connection details
   ```

3. **Test with different timeout**:
   ```bash
   ssh -o ConnectTimeout=10 user@hostname
   ```

**Solutions**:
- **Increase timeout**: Set longer connection timeout in configuration
- **Check firewall**: Verify port 22 is open
- **Check SSH service**: Ensure SSH service is running on target
- **Use different port**: If SSH is on non-standard port

**Code Example**:
```go
// Configure connection timeout
request := &spookytypes.ConnectionRequest{
    Host:     "example.com",
    Port:     22,
    User:     "user",
    KeyPath:  "~/.ssh/id_rsa",
    Timeout:  60 * time.Second, // Increase timeout
}
```

#### Host Key Validation Issues

**Problem**: Host key verification fails

**Diagnostic Steps**:
1. **Check known_hosts**:
   ```bash
   ssh-keygen -F hostname
   # Check if host key is in known_hosts
   ```

2. **Test host key manually**:
   ```bash
   ssh -o StrictHostKeyChecking=no user@hostname
   # Temporarily disable strict checking
   ```

**Solutions**:
- **Add host key**: `ssh-keyscan -H hostname >> ~/.ssh/known_hosts`
- **Update host key**: Remove old key and add new one
- **Disable strict checking**: Set `StrictHostKeyCheck: false` (development only)
- **Use custom known_hosts**: Specify custom known_hosts file

### File Transfer Issues

#### SFTP Transfer Failures

**Problem**: SFTP file transfers fail

**Diagnostic Steps**:
1. **Check file permissions**:
   ```bash
   ls -la /path/to/file
   # Verify source file exists and is readable
   ```

2. **Test SFTP manually**:
   ```bash
   sftp user@hostname
   # Test manual SFTP connection
   ```

3. **Check disk space**:
   ```bash
   df -h
   # Check available disk space
   ```

**Solutions**:
- **Fix permissions**: Ensure proper file permissions
- **Check disk space**: Ensure sufficient disk space
- **Create directories**: Ensure remote directories exist
- **Use SCP**: Try SCP instead of SFTP for problematic transfers

**Code Example**:
```go
// Validate file transfer request
func validateTransferRequest(transfer *spookytypesssh.FileTransfer) error {
    if transfer.LocalPath == "" {
        return fmt.Errorf("local path is required")
    }
    
    if transfer.RemotePath == "" {
        return fmt.Errorf("remote path is required")
    }
    
    // Check local file exists for upload
    if transfer.Direction == spookytypesssh.TransferDirectionUpload {
        if _, err := os.Stat(transfer.LocalPath); os.IsNotExist(err) {
            return fmt.Errorf("local file does not exist: %s", transfer.LocalPath)
        }
    }
    
    return nil
}
```

#### SCP Transfer Failures

**Problem**: SCP file transfers fail

**Diagnostic Steps**:
1. **Test SCP manually**:
   ```bash
   scp /local/file user@hostname:/remote/path
   # Test manual SCP transfer
   ```

2. **Check remote permissions**:
   ```bash
   ssh user@hostname "ls -la /remote/path"
   # Check remote directory permissions
   ```

**Solutions**:
- **Fix remote permissions**: Ensure write permissions on remote directory
- **Create remote directories**: Ensure remote directories exist
- **Use SFTP**: Try SFTP instead of SCP for problematic transfers
- **Check file size**: Ensure file size is within limits

### Command Execution Issues

#### Command Not Found

**Problem**: Remote commands fail with "command not found"

**Diagnostic Steps**:
1. **Check command path**:
   ```bash
   ssh user@hostname "which command"
   # Check if command exists
   ```

2. **Check PATH**:
   ```bash
   ssh user@hostname "echo \$PATH"
   # Check PATH environment variable
   ```

**Solutions**:
- **Use full path**: Specify full path to command
- **Set PATH**: Set PATH environment variable in command
- **Install command**: Ensure command is installed on remote system
- **Use different shell**: Specify different shell if needed

**Code Example**:
```go
// Execute command with full path
command := &spookytypes.SSHCommand{
    Command:       "/usr/bin/ls -la",
    WorkingDir:    "/home/user",
    Environment:   map[string]string{"PATH": "/usr/bin:/usr/local/bin"},
    Timeout:       30 * time.Second,
    CaptureOutput: true,
}
```

#### Permission Denied

**Problem**: Commands fail with "Permission denied"

**Diagnostic Steps**:
1. **Check user permissions**:
   ```bash
   ssh user@hostname "id"
   # Check user ID and groups
   ```

2. **Test command manually**:
   ```bash
   ssh user@hostname "sudo -l"
   # Check sudo permissions
   ```

**Solutions**:
- **Use sudo**: Prefix commands with `sudo` if needed
- **Fix permissions**: Ensure proper file/directory permissions
- **Change user**: Use different user with appropriate permissions
- **Configure sudo**: Configure sudo access for commands

### Performance Issues

#### Slow Connections

**Problem**: SSH connections are slow

**Diagnostic Steps**:
1. **Test connection speed**:
   ```bash
   time ssh user@hostname "echo hello"
   # Measure connection time
   ```

2. **Check network latency**:
   ```bash
   ping -c 10 hostname
   # Check network latency
   ```

**Solutions**:
- **Use connection pooling**: Enable connection pooling for reuse
- **Optimize network**: Check network configuration
- **Use compression**: Enable SSH compression
- **Use faster algorithms**: Configure faster SSH algorithms

**Code Example**:
```go
// Configure connection pooling
config := &spookytypes.ClientConfig{
    MaxConnections:   10,
    IdleTimeout:      300 * time.Second,
    DefaultTimeout:   30 * time.Second,
    MaxRetryAttempts: 3,
    RetryDelay:       5 * time.Second,
}
```

#### Slow File Transfers

**Problem**: File transfers are slow

**Diagnostic Steps**:
1. **Check transfer speed**:
   ```bash
   scp -v /large/file user@hostname:/remote/path
   # Verbose transfer shows progress
   ```

2. **Check network bandwidth**:
   ```bash
   iperf3 -c hostname
   # Test network bandwidth
   ```

**Solutions**:
- **Use compression**: Enable transfer compression
- **Use parallel transfers**: Transfer multiple files in parallel
- **Optimize buffer size**: Adjust transfer buffer size
- **Use faster protocol**: Try different transfer protocols

## Debugging Techniques

### Enable Debug Logging

**Enable verbose SSH logging**:
```bash
export SPOOKY_LOG_LEVEL=debug
spooky machines ping my-project
```

**Enable SSH debug output**:
```bash
ssh -v user@hostname
# Verbose SSH output
```

### Connection Diagnostics

**Test SSH connectivity step by step**:
```bash
# 1. Test DNS resolution
nslookup hostname

# 2. Test network connectivity
ping hostname

# 3. Test port connectivity
telnet hostname 22

# 4. Test SSH connection
ssh -v user@hostname

# 5. Test authentication
ssh -i ~/.ssh/id_rsa user@hostname
```

### File Transfer Diagnostics

**Test file transfer step by step**:
```bash
# 1. Test SFTP connection
sftp user@hostname

# 2. Test SCP connection
scp /test/file user@hostname:/tmp/

# 3. Test file permissions
ssh user@hostname "ls -la /tmp/test/file"

# 4. Test directory creation
ssh user@hostname "mkdir -p /remote/path"
```

## Error Messages and Solutions

### Common Error Messages

#### "Permission denied (publickey,password)"
**Cause**: Authentication failed
**Solution**: Check SSH key permissions, verify credentials, ensure public key is in authorized_keys

#### "Connection timed out"
**Cause**: Network connectivity issue
**Solution**: Check network connectivity, firewall settings, SSH service status

#### "Host key verification failed"
**Cause**: Host key mismatch
**Solution**: Update known_hosts file, verify host key, or temporarily disable strict checking

#### "No such file or directory"
**Cause**: File or directory not found
**Solution**: Check file paths, create missing directories, verify file existence

#### "Disk quota exceeded"
**Cause**: Insufficient disk space
**Solution**: Free up disk space, check disk quotas, use different location

### Error Code Reference

| Exit Code | Meaning | Solution |
|-----------|---------|----------|
| 0 | Success | No action needed |
| 1 | General error | Check command syntax and permissions |
| 2 | Misuse of shell builtins | Check command syntax |
| 126 | Command not executable | Check file permissions |
| 127 | Command not found | Install command or use full path |
| 128+n | Signal n | Process terminated by signal |

## Best Practices

### Security Best Practices

1. **Use SSH keys**: Prefer SSH key authentication over passwords
2. **Validate host keys**: Always validate host keys for security
3. **Secure key storage**: Store SSH keys with 600 permissions
4. **Use strong keys**: Use Ed25519 or RSA 4096-bit keys
5. **Rotate keys**: Regularly rotate SSH keys

### Performance Best Practices

1. **Use connection pooling**: Enable connection pooling for efficiency
2. **Optimize timeouts**: Set appropriate timeouts for operations
3. **Use compression**: Enable compression for slow connections
4. **Parallel operations**: Use parallel operations when possible
5. **Monitor performance**: Monitor connection and transfer performance

### Troubleshooting Best Practices

1. **Test manually first**: Always test SSH operations manually before automation
2. **Use verbose logging**: Enable verbose logging for debugging
3. **Check permissions**: Verify file and directory permissions
4. **Test step by step**: Test each step of the process individually
5. **Document solutions**: Document solutions for future reference

## Configuration Examples

### SSH Client Configuration

```hcl
# SSH client configuration
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

### Machine Configuration

```hcl
# Machine configuration with SSH
machines {
  machine "web-server" {
    hostname = "web.example.com"
    port = 22
    user = "admin"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/id_rsa"
      passphrase = "optional_passphrase"
    }
    
    ssh {
      timeout = "30s"
      keepalive = "60s"
      host_key_validation = true
    }
  }
}
```

### File Transfer Configuration

```hcl
# File transfer configuration
file_transfer {
  mode = "sftp"
  verify = true
  permissions = 0644
  buffer_size = 32768
  compression = true
}
```

## Getting Help

### Debug Information

When reporting SSH issues, include:

1. **Error messages**: Complete error messages and stack traces
2. **Configuration**: Relevant configuration files
3. **Environment**: Operating system, spooky version, SSH version
4. **Steps to reproduce**: Detailed steps to reproduce the issue
5. **Manual test results**: Results of manual SSH tests

### Log Files

Check these log locations for SSH-related information:

- **spooky logs**: Application logs with SSH operations
- **SSH logs**: System SSH logs (`/var/log/auth.log` on Linux)
- **Debug logs**: Debug logs when `SPOOKY_LOG_LEVEL=debug` is set

### Support Resources

- **Documentation**: Check SSH API reference and user guide
- **Examples**: Review SSH configuration examples
- **Community**: Check community forums and issue trackers
- **Debugging**: Use debugging techniques described in this guide

## Summary

The SSH system in spooky provides comprehensive functionality with robust error handling and diagnostic capabilities. This troubleshooting guide covers the most common issues and provides practical solutions for resolving SSH-related problems.

Key points for effective troubleshooting:

1. **Test manually first**: Always test SSH operations manually before automation
2. **Check permissions**: Verify file, directory, and SSH key permissions
3. **Use verbose logging**: Enable debug logging for detailed diagnostics
4. **Follow security best practices**: Use SSH keys, validate host keys, secure key storage
5. **Monitor performance**: Monitor connection and transfer performance
6. **Document solutions**: Document solutions for future reference

The SSH system is production-ready and provides all necessary functionality for secure, efficient, and reliable SSH operations in the spooky automation platform.
