# Machines Inventory Troubleshooting Guide

## Overview

This troubleshooting guide provides solutions for common issues encountered when working with the spooky machines inventory system. It covers error messages, connectivity problems, validation issues, and performance problems.

**Status: Production Ready** - The machines system is fully implemented with comprehensive inventory management, connectivity testing, and validation capabilities.

## Machines System Status

### ✅ Fully Functional Machines Infrastructure

The machines system now has **complete machines infrastructure** with:

- **Inventory Management**: Comprehensive machine inventory loading and validation
- **Connectivity Testing**: SSH-based connectivity testing with DNS resolution and ICMP checks
- **Machine Validation**: Complete machine configuration validation
- **CLI Integration**: Full CLI integration with `spooky machines` commands
- **Project Integration**: Machine inventory loading from project configuration
- **Export Functionality**: Machine inventory export to JSON format
- **Error Handling**: Comprehensive error handling and reporting

### What This Means for Users

- **No More Stubs**: All functionality is fully implemented - no placeholder code
- **Production Ready**: The system is ready for production use
- **Complete Feature Set**: All documented features are functional
- **Reliable Connectivity**: Robust connectivity testing and validation
- **Performance Optimized**: Efficient inventory management and connectivity testing

### Expected Behavior

When using machines, you can expect:

1. **Proper Inventory Loading**: Machine inventory loads correctly from HCL files
2. **Connectivity Testing**: SSH connectivity testing with proper error reporting
3. **Validation**: Comprehensive machine configuration validation
4. **Export Functionality**: Machine inventory export to JSON format
5. **Error Handling**: Clear error messages with actionable information

## Common Issues and Solutions

### Loading Errors

#### "Failed to load machines: no machines found in project"

**Cause:** No machine inventory files found in the project directory.

**Solution:**
```bash
# Check if machines.hcl exists
ls -la ./my-project/machines.hcl

# Check if machines/ directory exists
ls -la ./my-project/machines/

# Create a basic machines.hcl file
cat > ./my-project/machines.hcl << 'EOF'
machines {
  machine "test-server" {
    host = "192.168.1.10"
    user = "admin"
    key_file = "~/.ssh/id_rsa"
  }
}
EOF
```

#### "Failed to parse machine block: Unexpected block type"

**Cause:** Invalid HCL syntax in machine configuration.

**Solution:**
```bash
# Check HCL syntax
spooky machines validate ./my-project

# Common syntax errors to fix:
# 1. Missing quotes around strings
# 2. Incorrect block structure
# 3. Invalid attribute names

# Example of correct syntax:
machines {
  machine "web-server" {
    host = "192.168.1.10"  # Use quotes for strings
    user = "admin"
    key_file = "~/.ssh/id_rsa"
    
    tags = ["web", "production"]  # Use array syntax
    
    metadata {
      environment = "production"
      owner = "web-team"
    }
  }
}
```

#### "duplicate hostname 'web-server' found in multiple files"

**Cause:** Same hostname defined in multiple inventory files.

**Solution:**
```bash
# Find all occurrences of the duplicate hostname
grep -r "web-server" ./my-project/machines/

# Rename one of the machines to make it unique
# In machines/production.hcl:
machine "prod-web-server" {  # Changed from "web-server"
  host = "192.168.1.10"
  user = "admin"
}

# In machines/staging.hcl:
machine "staging-web-server" {  # Changed from "web-server"
  host = "192.168.1.20"
  user = "admin"
}
```

### Validation Errors

#### "Validation failed: missing required field 'host'"

**Cause:** Machine configuration is missing required fields.

**Solution:**
```hcl
# ✅ CORRECT - All required fields present
machines {
  machine "web-server" {
    host = "192.168.1.10"    # Required
    user = "admin"           # Required
    key_file = "~/.ssh/id_rsa"  # Required for SSH key authentication
  }
}
```

#### "Validation failed: invalid port number"

**Cause:** Port number is outside valid range (1-65535).

**Solution:**
```hcl
# ✅ CORRECT - Valid port number
machines {
  machine "web-server" {
    host = "192.168.1.10"
    port = 22                # Valid port (1-65535)
    user = "admin"
    key_file = "~/.ssh/id_rsa"
  }
}
```

#### "Validation failed: SSH key file not found"

**Cause:** SSH key file path is invalid or file doesn't exist.

**Solution:**
```bash
# Check if key file exists
ls -la ~/.ssh/id_rsa

# Check key file permissions
chmod 600 ~/.ssh/id_rsa

# Generate new key if needed
ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa
```

### Connectivity Issues

#### "DNS resolution failed for hostname"

**Cause:** Cannot resolve hostname to IP address.

**Solution:**
```bash
# Test DNS resolution manually
nslookup example.com

# Check DNS configuration
cat /etc/resolv.conf

# Use IP address instead of hostname
machines {
  machine "web-server" {
    host = "192.168.1.10"  # Use IP instead of hostname
    user = "admin"
    key_file = "~/.ssh/id_rsa"
  }
}
```

#### "SSH port not reachable"

**Cause:** SSH port is not accessible (firewall, service not running).

**Solution:**
```bash
# Test SSH port connectivity
telnet 192.168.1.10 22

# Check if SSH service is running on target
ssh admin@192.168.1.10 "systemctl status sshd"

# Check firewall settings
ssh admin@192.168.1.10 "sudo ufw status"
```

#### "SSH authentication failed"

**Cause:** SSH authentication is failing (wrong key, user, or permissions).

**Solution:**
```bash
# Test SSH connection manually
ssh -i ~/.ssh/id_rsa admin@192.168.1.10

# Check key permissions
chmod 600 ~/.ssh/id_rsa

# Verify key is authorized on server
ssh admin@192.168.1.10 "cat ~/.ssh/authorized_keys | grep $(cat ~/.ssh/id_rsa.pub)"
```

### Performance Issues

#### "Machine ping is very slow"

**Cause:** Network latency or slow DNS resolution.

**Solution:**
```bash
# Test network latency
ping -c 5 192.168.1.10

# Use IP addresses instead of hostnames
machines {
  machine "web-server" {
    host = "192.168.1.10"  # Use IP instead of hostname
    user = "admin"
    key_file = "~/.ssh/id_rsa"
  }
}

# Increase timeout for slow connections
machines {
  machine "web-server" {
    host = "192.168.1.10"
    user = "admin"
    key_file = "~/.ssh/id_rsa"
    connection_timeout = 60  # Increase timeout
  }
}
```

#### "High memory usage during machine operations"

**Cause:** Too many machines or inefficient operations.

**Solution:**
```bash
# Test fewer machines at once
spooky machines ping ./my-project --machines machine1,machine2

# Use parallel operations with limits
spooky machines ping ./my-project --parallel 4

# Monitor memory usage
top -p $(pgrep spooky)
```

## Configuration Problems

### Project Configuration Issues

#### "Invalid machine configuration: unknown field"

**Cause:** Unknown fields in machine configuration.

**Solution:**
```hcl
# ✅ CORRECT - Valid machine configuration
machines {
  machine "web-server" {
    host = "192.168.1.10"
    port = 22
    user = "admin"
    key_file = "~/.ssh/id_rsa"
    
    tags = ["web", "production"]
    
    metadata {
      environment = "production"
      owner = "web-team"
    }
  }
}
```

#### "Invalid authentication configuration: missing method"

**Cause:** Authentication configuration is incomplete.

**Solution:**
```hcl
# ✅ CORRECT - Complete authentication configuration
machines {
  machine "web-server" {
    host = "192.168.1.10"
    user = "admin"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/id_rsa"
    }
  }
}
```

### SSH Configuration Issues

#### "SSH connection failed: no such identity"

**Cause:** SSH key file not found.

**Solution:**
```bash
# Check if key file exists
ls -la ~/.ssh/id_rsa

# Generate new SSH key if needed
ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa

# Copy key to target machine
ssh-copy-id -i ~/.ssh/id_rsa admin@192.168.1.10
```

#### "SSH connection failed: host key verification failed"

**Cause:** SSH host key verification is failing.

**Solution:**
```bash
# Add host key to known_hosts
ssh-keyscan -H 192.168.1.10 >> ~/.ssh/known_hosts

# Or disable host key checking (not recommended for production)
ssh -o StrictHostKeyChecking=no admin@192.168.1.10
```

## Network Issues

### Connectivity Problems

#### "Network unreachable"

**Cause:** Cannot reach the target machine.

**Solution:**
```bash
# Check network connectivity
ping 192.168.1.10

# Check routing
traceroute 192.168.1.10

# Check firewall rules
sudo iptables -L
```

#### "Connection timed out"

**Cause:** Network connection is timing out.

**Solution:**
```bash
# Check if machine is reachable
telnet 192.168.1.10 22

# Test SSH connection with timeout
ssh -o ConnectTimeout=10 admin@192.168.1.10

# Check network latency
ping -c 5 192.168.1.10
```

### Firewall Issues

#### "Connection refused by firewall"

**Cause:** Firewall is blocking SSH connections.

**Solution:**
```bash
# Check local firewall
sudo ufw status

# Check remote firewall
ssh admin@192.168.1.10 "sudo iptables -L"

# Open SSH port
sudo ufw allow 22/tcp
```

## Debugging Techniques

### Enable Verbose Output

```bash
# Enable verbose output for machine operations
spooky machines list ./my-project --verbose
spooky machines ping ./my-project --verbose

# Enable debug logging
export SPOOKY_LOG_LEVEL=debug
spooky machines ping ./my-project
```

### Test Connectivity Manually

```bash
# Test DNS resolution
nslookup example.com

# Test network connectivity
ping 192.168.1.10

# Test SSH port
telnet 192.168.1.10 22

# Test SSH connection
ssh -i ~/.ssh/id_rsa admin@192.168.1.10 "echo 'SSH working'"
```

### Validate Configuration

```bash
# Validate machine configuration
spooky machines validate ./my-project

# Check specific machine
spooky machines validate ./my-project --machine web-server

# Export machine inventory for inspection
spooky machines export ./my-project --output machines.json
```

### Network Diagnostics

```bash
# Test network connectivity
ping -c 5 192.168.1.10

# Check routing
traceroute 192.168.1.10

# Check DNS resolution
nslookup example.com

# Check SSH connectivity
ssh -o ConnectTimeout=10 admin@192.168.1.10
```

## Recovery Procedures

### Machine Configuration Recovery

```bash
# Backup configuration
cp machines.hcl machines.hcl.backup

# Validate configuration
spooky machines validate ./my-project

# Restore from backup if needed
cp machines.hcl.backup machines.hcl
```

### SSH Connection Recovery

```bash
# Test SSH connectivity
spooky machines ping ./my-project

# Check SSH keys
ls -la ~/.ssh/
chmod 600 ~/.ssh/id_rsa

# Regenerate SSH key if needed
ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa

# Copy key to machines
ssh-copy-id -i ~/.ssh/id_rsa admin@192.168.1.10
```

### Network Recovery

```bash
# Check network connectivity
ping 192.168.1.10

# Check firewall settings
sudo ufw status

# Check SSH service
ssh admin@192.168.1.10 "systemctl status sshd"
```

## Prevention Strategies

### Regular Validation

```bash
# Schedule regular validation
crontab -e
# Add: 0 2 * * * /usr/local/bin/spooky machines validate /path/to/project

# Validate before operations
spooky machines validate ./my-project

# Validate in CI/CD pipeline
spooky machines validate ./my-project --strict
```

### Monitoring

```bash
# Monitor machine connectivity
spooky machines ping ./my-project

# Monitor SSH connections
netstat -an | grep :22

# Monitor system resources
top -p $(pgrep spooky)
```

### Backup Strategy

```bash
# Backup machine configurations
cp machines.hcl machines.hcl.$(date +%Y%m%d)

# Version control configurations
git add machines.hcl
git commit -m "Update machine configuration"

# Backup project structure
tar -czf project-backup-$(date +%Y%m%d).tar.gz ./
```

## Best Practices for Troubleshooting

### 1. Start Simple

Begin with simple machine configurations and add complexity gradually:

```hcl
# Start with basic configuration
machines {
  machine "test-server" {
    host = "192.168.1.10"
    user = "admin"
    key_file = "~/.ssh/id_rsa"
  }
}

# Then add complexity
machines {
  machine "web-server" {
    host = "192.168.1.10"
    user = "admin"
    key_file = "~/.ssh/id_rsa"
    
    tags = ["web", "production"]
    
    metadata {
      environment = "production"
      owner = "web-team"
    }
  }
}
```

### 2. Use Descriptive Names

Use clear, descriptive machine names:

```hcl
# ✅ GOOD - Descriptive names
machines {
  machine "prod-web-server-01" {
    host = "192.168.1.10"
    user = "admin"
    key_file = "~/.ssh/id_rsa"
  }
}

# ❌ BAD - Unclear names
machines {
  machine "server1" {
    host = "192.168.1.10"
    user = "admin"
    key_file = "~/.ssh/id_rsa"
  }
}
```

### 3. Validate Early and Often

Validate configurations frequently:

```bash
# Validate after every change
spooky machines validate ./my-project

# Validate before operations
spooky machines validate ./my-project && spooky machines ping ./my-project

# Validate in scripts
#!/bin/bash
if spooky machines validate ./my-project; then
    spooky machines ping ./my-project
else
    echo "Validation failed"
    exit 1
fi
```

### 4. Use Proper Error Handling

Implement proper error handling in configurations:

```hcl
# Use proper SSH key configuration
machines {
  machine "web-server" {
    host = "192.168.1.10"
    user = "admin"
    key_file = "~/.ssh/id_rsa"  # Ensure key exists and has correct permissions
    
    # Add connection options for reliability
    connection_timeout = 30
    retry_attempts = 3
  }
}
```

### 5. Monitor and Log

Monitor machine operations and maintain logs:

```bash
# Enable verbose logging
spooky machines ping ./my-project --verbose

# Monitor operations
watch -n 1 'ps aux | grep spooky'

# Check logs
tail -f /var/log/spooky/machines.log
```

## Getting Help

### Documentation Resources

1. **User Guide** - For usage questions and best practices
2. **API Reference** - For technical implementation details
3. **Examples** - For configuration patterns and use cases

### Common Questions

#### "Why can't I connect to my machines?"

1. Check machine configuration
2. Verify SSH connectivity
3. Check network connectivity
4. Validate SSH keys

#### "How do I debug connectivity issues?"

```bash
# Enable verbose output
spooky machines ping ./my-project --verbose

# Test connectivity manually
ssh -i ~/.ssh/id_rsa admin@192.168.1.10

# Check network connectivity
ping 192.168.1.10
```

#### "How do I fix validation issues?"

```bash
# Validate configuration
spooky machines validate ./my-project --verbose

# Check specific errors
spooky machines validate ./my-project --machine problematic-machine

# Fix configuration issues
# Update machine configuration based on error messages
```

#### "How do I optimize machine operations?"

```bash
# Use parallel operations
spooky machines ping ./my-project --parallel 4

# Use machine filtering
spooky machines ping ./my-project --machines machine1,machine2

# Monitor resource usage
top -p $(pgrep spooky)
```

### When to Seek Additional Help

- Configuration validation passes but connectivity still fails
- Performance issues persist after optimization
- Unusual error messages not covered in this guide
- Integration issues with other spooky components

For additional help, refer to the [User Guide](MACHINES_USER_GUIDE.md) and [API Reference](MACHINES_API_REFERENCE.md), or check the project documentation for more advanced troubleshooting techniques.

## Conclusion

The machines system provides robust, reliable machine inventory management with comprehensive connectivity testing and validation capabilities. Most issues can be resolved by following the troubleshooting steps outlined in this guide. For persistent issues, enable verbose output and collect diagnostic information for further analysis.
