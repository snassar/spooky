# Machines Inventory Troubleshooting Guide

## Overview

This troubleshooting guide provides solutions for common issues encountered when working with the spooky machines inventory system. It covers error messages, connectivity problems, validation issues, and performance problems.

## Table of Contents

1. [Common Error Messages](#common-error-messages)
2. [Connectivity Issues](#connectivity-issues)
3. [Validation Problems](#validation-problems)
4. [Performance Issues](#performance-issues)
5. [Configuration Problems](#configuration-problems)
6. [Debugging Techniques](#debugging-techniques)
7. [Best Practices for Troubleshooting](#best-practices-for-troubleshooting)

## Common Error Messages

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

#### "machine 'web-server' missing authentication method"

**Cause:** No SSH key file or authentication method specified.

**Solution:**
```hcl
# Add key_file attribute
machine "web-server" {
  host = "192.168.1.10"
  user = "admin"
  key_file = "~/.ssh/id_rsa"  # Add this line
}
```

#### "SSH key file '/path/to/key' not found"

**Cause:** SSH key file doesn't exist or path is incorrect.

**Solution:**
```bash
# Check if key file exists
ls -la ~/.ssh/id_rsa

# Generate a new SSH key if needed
ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa

# Set correct permissions
chmod 600 ~/.ssh/id_rsa

# Test SSH connection manually
ssh -i ~/.ssh/id_rsa admin@192.168.1.10
```

#### "port must be between 1 and 65535"

**Cause:** Invalid port number specified.

**Solution:**
```hcl
# Use valid port number
machine "web-server" {
  host = "192.168.1.10"
  user = "admin"
  port = 22  # Valid port number
  key_file = "~/.ssh/id_rsa"
}
```

### Connectivity Errors

#### "DNS resolution failed"

**Cause:** Hostname cannot be resolved to IP address.

**Solution:**
```bash
# Test DNS resolution manually
nslookup web-server.example.com
dig web-server.example.com

# Check if using IP address instead of hostname
# In machines.hcl:
machine "web-server" {
  host = "192.168.1.10"  # Use IP instead of hostname
  user = "admin"
  key_file = "~/.ssh/id_rsa"
}

# Or add hostname to /etc/hosts
echo "192.168.1.10 web-server.example.com" | sudo tee -a /etc/hosts
```

#### "connection refused"

**Cause:** SSH service not running or port blocked.

**Solution:**
```bash
# Check if SSH service is running on target machine
ssh admin@192.168.1.10 "systemctl status sshd"

# Check if port is open
telnet 192.168.1.10 22
nc -zv 192.168.1.10 22

# Check firewall settings
ssh admin@192.168.1.10 "sudo ufw status"
ssh admin@192.168.1.10 "sudo iptables -L"
```

#### "authentication failed"

**Cause:** SSH key authentication failed.

**Solution:**
```bash
# Check SSH key permissions
ls -la ~/.ssh/id_rsa
chmod 600 ~/.ssh/id_rsa

# Check if key is in authorized_keys on target machine
ssh admin@192.168.1.10 "cat ~/.ssh/authorized_keys"

# Add public key to target machine
ssh-copy-id -i ~/.ssh/id_rsa.pub admin@192.168.1.10

# Test SSH connection manually
ssh -i ~/.ssh/id_rsa admin@192.168.1.10
```

## Connectivity Issues

### Progressive Connectivity Testing

Spooky performs connectivity tests in stages. Understanding each stage helps with troubleshooting:

#### Stage 1: DNS Resolution
```bash
# Test DNS resolution manually
nslookup web-server.example.com
dig web-server.example.com

# Check if hostname is correct
cat ./my-project/machines.hcl | grep host
```

#### Stage 2: ICMP Ping (Simulated)
Currently simulated. In future versions, this will test actual ICMP connectivity.

#### Stage 3: TCP Port Scan
```bash
# Test SSH port manually
telnet web-server.example.com 22
nc -zv web-server.example.com 22

# Check if SSH service is running
ssh admin@web-server.example.com "systemctl status sshd"
```

#### Stage 4: SSH Authentication
```bash
# Test SSH connection manually
ssh -i ~/.ssh/id_rsa admin@web-server.example.com

# Check SSH configuration
ssh -v -i ~/.ssh/id_rsa admin@web-server.example.com
```

### Network Connectivity Troubleshooting

#### Check Network Connectivity
```bash
# Basic connectivity test
ping 192.168.1.10

# Test specific port
telnet 192.168.1.10 22

# Check routing
traceroute 192.168.1.10

# Check DNS resolution
nslookup web-server.example.com
```

#### Firewall Issues
```bash
# Check local firewall
sudo ufw status
sudo iptables -L

# Check remote firewall (if accessible)
ssh admin@192.168.1.10 "sudo ufw status"
ssh admin@192.168.1.10 "sudo iptables -L"

# Common firewall rules for SSH
sudo ufw allow 22/tcp
sudo ufw allow from 192.168.1.0/24 to any port 22
```

#### SSH Service Issues
```bash
# Check SSH service status
ssh admin@192.168.1.10 "sudo systemctl status sshd"

# Restart SSH service if needed
ssh admin@192.168.1.10 "sudo systemctl restart sshd"

# Check SSH configuration
ssh admin@192.168.1.10 "sudo cat /etc/ssh/sshd_config | grep -E 'Port|PermitRootLogin|PasswordAuthentication'"
```

## Validation Problems

### Environment-Specific Validation

#### Production Environment Warnings
```bash
# Common production warnings:
# - Missing resource specifications
# - Missing backup schedule
# - Missing cost center information

# Fix by adding required fields:
machine "prod-web-server" {
  host = "192.168.1.10"
  user = "admin"
  key_file = "~/.ssh/id_rsa"
  
  resources {
    cpu_cores = 8
    memory_gb = 32
    disk_gb = 500
  }
  
  metadata {
    environment = "production"
    backup_schedule = "daily"
    cost_center = "IT-001"
  }
}
```

#### Development Environment Warnings
```bash
# Common development warnings:
# - Missing owner information
# - Missing purpose documentation

# Fix by adding metadata:
machine "dev-web-server" {
  host = "192.168.1.20"
  user = "developer"
  key_file = "~/.ssh/dev_key"
  
  metadata {
    environment = "development"
    owner = "developer"
    purpose = "web development and testing"
  }
}
```

### Cross-File Validation Issues

#### Duplicate Host Addresses
```bash
# Warning: duplicate host address '192.168.1.10' used by multiple machines

# Solution: Use different IP addresses or document the reason
machine "web-server-1" {
  host = "192.168.1.10"
  user = "admin"
}

machine "web-server-2" {
  host = "192.168.1.11"  # Use different IP
  user = "admin"
}
```

#### Inconsistent Authentication Methods
```bash
# Warning: mixed authentication methods in production environment

# Solution: Use consistent authentication
# All production machines should use key-based authentication
machine "prod-web-server" {
  host = "192.168.1.10"
  user = "admin"
  key_file = "~/.ssh/prod_key"  # Consistent key file
}

machine "prod-db-server" {
  host = "192.168.1.20"
  user = "dbadmin"
  key_file = "~/.ssh/prod_key"  # Same key file
}
```

## Performance Issues

### Slow Connectivity Testing

#### Parallel Execution
```bash
# Use parallel execution for large inventories
spooky machines ping ./my-project --parallel 10

# Monitor system resources during testing
htop
iostat 1
```

#### Network Latency
```bash
# Check network latency to target machines
ping -c 10 192.168.1.10

# Use closer machines for testing
# Consider using local machines for development
```

#### SSH Connection Timeouts
```bash
# Check SSH connection timeouts
ssh -o ConnectTimeout=10 -i ~/.ssh/id_rsa admin@192.168.1.10

# Adjust timeout settings in SSH config
cat >> ~/.ssh/config << 'EOF'
Host *
    ConnectTimeout 10
    ServerAliveInterval 60
    ServerAliveCountMax 3
EOF
```

### Memory and CPU Usage

#### High Memory Usage
```bash
# Monitor memory usage during operations
free -h
ps aux | grep spooky

# Consider processing machines in batches
# Use streaming output for large inventories
spooky machines ping ./my-project --format json | head -100
```

#### High CPU Usage
```bash
# Monitor CPU usage
top
htop

# Consider reducing parallel connections
spooky machines ping ./my-project --parallel 5
```

## Configuration Problems

### HCL Syntax Errors

#### Common HCL Mistakes
```hcl
# WRONG: Missing quotes around strings
machine "web-server" {
  host = 192.168.1.10  # Should be "192.168.1.10"
  user = admin         # Should be "admin"
}

# CORRECT: Use quotes for strings
machine "web-server" {
  host = "192.168.1.10"
  user = "admin"
}

# WRONG: Incorrect array syntax
tags = [web, production]  # Should be ["web", "production"]

# CORRECT: Use quotes in arrays
tags = ["web", "production"]

# WRONG: Invalid block structure
machine "web-server" {
  host = "192.168.1.10"
  metadata environment = "production"  # Should be in metadata block
}

# CORRECT: Use proper block structure
machine "web-server" {
  host = "192.168.1.10"
  metadata {
    environment = "production"
  }
}
```

#### Schema Validation Errors
```bash
# Check schema compliance
spooky machines validate ./my-project --verbose

# Common schema violations:
# - Missing required fields
# - Invalid field types
# - Unsupported attributes
```

### File Organization Issues

#### Missing Source Files
```bash
# Check if all referenced files exist
find ./my-project -name "*.hcl" -exec echo "Found: {}" \;

# Check for broken symlinks
find ./my-project -type l -exec test ! -e {} \; -print

# Verify file permissions
ls -la ./my-project/machines/
ls -la ./my-project/machines.hcl
```

#### File Permission Issues
```bash
# Check file permissions
ls -la ./my-project/machines.hcl

# Fix permissions if needed
chmod 644 ./my-project/machines.hcl
chmod 755 ./my-project/machines/

# Check SSH key permissions
ls -la ~/.ssh/id_rsa
chmod 600 ~/.ssh/id_rsa
```

## Debugging Techniques

### Verbose Output

#### Enable Verbose Mode
```bash
# Get detailed output for all commands
spooky machines list ./my-project --verbose
spooky machines validate ./my-project --verbose
spooky machines ping ./my-project --verbose
```

#### JSON Output for Scripting
```bash
# Get machine status in JSON format
spooky machines ping ./my-project --format json

# Filter JSON output with jq
spooky machines ping ./my-project --format json | jq '.hostname, .status'

# Find offline machines
spooky machines ping ./my-project --format json | jq 'select(.status != "online")'
```

### Logging and Debugging

#### Enable Debug Logging
```bash
# Set debug log level
export SPOOKY_LOG_LEVEL=debug

# Run commands with debug output
spooky machines ping ./my-project
```

#### Manual Testing
```bash
# Test each component manually
# 1. Test HCL parsing
cat ./my-project/machines.hcl | hcl2json

# 2. Test DNS resolution
nslookup web-server.example.com

# 3. Test SSH connection
ssh -i ~/.ssh/id_rsa admin@192.168.1.10

# 4. Test port connectivity
telnet 192.168.1.10 22
```

### Step-by-Step Debugging

#### Debug Machine Loading
```bash
# Step 1: Check if project directory exists
ls -la ./my-project/

# Step 2: Check for machine files
ls -la ./my-project/machines.hcl
ls -la ./my-project/machines/

# Step 3: Validate HCL syntax
spooky machines validate ./my-project

# Step 4: Test individual machine loading
spooky machines list ./my-project
```

#### Debug Connectivity Issues
```bash
# Step 1: Test basic connectivity
ping 192.168.1.10

# Step 2: Test DNS resolution
nslookup web-server.example.com

# Step 3: Test SSH port
telnet 192.168.1.10 22

# Step 4: Test SSH authentication
ssh -i ~/.ssh/id_rsa admin@192.168.1.10

# Step 5: Test with spooky
spooky machines ping ./my-project --machine "web-server"
```

## Best Practices for Troubleshooting

### Systematic Approach

1. **Start with the Basics**
   - Check if files exist and are readable
   - Verify HCL syntax
   - Test basic connectivity

2. **Use Verbose Output**
   - Enable verbose mode for detailed information
   - Use JSON output for scripting and analysis

3. **Test Components Individually**
   - Test HCL parsing separately
   - Test network connectivity manually
   - Test SSH authentication manually

4. **Check Logs and Error Messages**
   - Read error messages carefully
   - Look for specific error codes
   - Check system logs if needed

### Common Patterns

#### Quick Health Check
```bash
#!/bin/bash
# Quick health check script

PROJECT_PATH="./my-project"

echo "=== Spooky Machines Health Check ==="

# Check project structure
echo "1. Checking project structure..."
if [ ! -d "$PROJECT_PATH" ]; then
    echo "ERROR: Project directory not found"
    exit 1
fi

# Check machine files
echo "2. Checking machine files..."
if [ ! -f "$PROJECT_PATH/machines.hcl" ] && [ ! -d "$PROJECT_PATH/machines" ]; then
    echo "ERROR: No machine inventory found"
    exit 1
fi

# Validate configuration
echo "3. Validating configuration..."
if ! spooky machines validate "$PROJECT_PATH"; then
    echo "ERROR: Configuration validation failed"
    exit 1
fi

# Test connectivity
echo "4. Testing connectivity..."
if ! spooky machines ping "$PROJECT_PATH" --format json | jq -e '.status == "online"' > /dev/null; then
    echo "WARNING: Some machines are offline"
fi

echo "=== Health check completed ==="
```

#### Automated Testing
```bash
#!/bin/bash
# Automated testing script

PROJECT_PATH="./my-project"
LOG_FILE="./machines-test.log"

echo "Starting automated testing at $(date)" | tee "$LOG_FILE"

# Test list command
echo "Testing list command..." | tee -a "$LOG_FILE"
if spooky machines list "$PROJECT_PATH" >> "$LOG_FILE" 2>&1; then
    echo "✓ List command successful"
else
    echo "✗ List command failed"
fi

# Test validate command
echo "Testing validate command..." | tee -a "$LOG_FILE"
if spooky machines validate "$PROJECT_PATH" >> "$LOG_FILE" 2>&1; then
    echo "✓ Validate command successful"
else
    echo "✗ Validate command failed"
fi

# Test ping command
echo "Testing ping command..." | tee -a "$LOG_FILE"
if spooky machines ping "$PROJECT_PATH" --format json >> "$LOG_FILE" 2>&1; then
    echo "✓ Ping command successful"
else
    echo "✗ Ping command failed"
fi

echo "Testing completed at $(date)" | tee -a "$LOG_FILE"
```

### Documentation

#### Keep Troubleshooting Notes
```markdown
# Troubleshooting Notes for Project: my-project

## Common Issues

### Issue 1: DNS Resolution Failed
- **Symptoms**: "DNS resolution failed" error
- **Cause**: Hostname not resolvable
- **Solution**: Use IP addresses or fix DNS
- **Date**: 2024-01-15

### Issue 2: SSH Authentication Failed
- **Symptoms**: "authentication failed" error
- **Cause**: SSH key not in authorized_keys
- **Solution**: Run ssh-copy-id
- **Date**: 2024-01-16

## Environment-Specific Notes

### Production Environment
- Uses dedicated SSH keys: ~/.ssh/prod_key
- All machines require resource specifications
- Backup schedule must be specified

### Development Environment
- Uses shared SSH key: ~/.ssh/dev_key
- More lenient validation rules
- Owner information required
```

This comprehensive troubleshooting guide provides solutions for the most common issues encountered when working with the spooky machines inventory system. Use it as a reference when encountering problems and follow the systematic approach for effective troubleshooting.
