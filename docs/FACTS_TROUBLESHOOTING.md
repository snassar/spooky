# Facts System Troubleshooting Guide

## Overview

This troubleshooting guide provides solutions for common issues encountered when working with the spooky facts system. It covers error messages, configuration problems, performance issues, and debugging techniques.

**Status: Partially Implemented** - The facts system has basic functionality but SSH-based fact collection has known issues that need to be addressed.

## Facts System Status

### ⚠️ Partially Functional Facts Infrastructure

The facts system currently has **basic facts infrastructure** with:

- **Local Fact Collection**: Basic system fact collection using gopsutil
- **Memory Storage**: In-memory fact storage during export operations
- **Export Functionality**: Facts export to JSON and HCL formats
- **CLI Integration**: `spooky facts export` command with filtering options
- **Project Integration**: Facts collection from project machine inventory
- **Basic Validation**: Fact collection validation and error handling

### Known Limitations

- **SSH-Based Collection**: SSH-based fact collection has implementation issues
- **Remote Facts Reading**: Cannot read `/etc/spooky/facts.*` files from remote machines
- **Parallel Processing**: Limited parallel collection support
- **Fact Caching**: No persistent fact storage or caching
- **Advanced Collectors**: Only basic system fact collector available

### Expected Behavior

When using facts, you can expect:

1. **Local Collection**: Basic system facts can be collected locally
2. **Export Functionality**: Facts can be exported to JSON and HCL formats
3. **Project Integration**: Facts operations work with project configuration
4. **Basic Validation**: Fact validation with error reporting
5. **Memory Management**: Facts stored in memory during operations only

## Common Issues and Solutions

### Collection Errors

#### "Failed to connect to machine: connection refused"

**Problem:** Cannot establish SSH connection to the target machine.

**Solution:**
```bash
# Check if SSH service is running on target machine
ssh -i ~/.ssh/id_rsa user@machine.example.com

# Verify SSH configuration
ssh -v -i ~/.ssh/id_rsa user@machine.example.com

# Check firewall settings
sudo ufw status
```

**Check your machine configuration:**
```hcl
machines {
  machine "web-server" {
    hostname = "web.example.com"  # Verify hostname is correct
    port = 22                     # Verify SSH port
    user = "admin"                # Verify username
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/id_rsa"  # Verify key path
    }
  }
}
```

#### "Authentication failed: permission denied"

**Problem:** SSH authentication is failing.

**Solution:**
```bash
# Check SSH key permissions
chmod 600 ~/.ssh/id_rsa
chmod 700 ~/.ssh

# Test SSH connection manually
ssh -i ~/.ssh/id_rsa user@machine.example.com

# Check if key is in authorized_keys
ssh user@machine.example.com "cat ~/.ssh/authorized_keys"
```

**Check your authentication configuration:**
```hcl
authentication {
  method = "ssh_key"
  key_path = "~/.ssh/id_rsa"  # Must be readable by spooky
}
```

#### "Failed to gather facts: memory limit exceeded"

**Problem:** Memory limit exceeded during fact gathering and export.

**Solution:**
```bash
# Check available memory
free -h

# Reduce parallel workers
spooky facts export ./my-project --parallel 2

# Monitor memory usage during export
spooky facts export ./my-project --verbose

# Export fewer machines at once
spooky facts export ./my-project --machine specific-machine --output single-machine-facts.json
```

#### "Failed to read facts: memory corruption detected"

**Problem:** Memory corruption detected during fact gathering and export.

**Solution:**
```bash
# Restart spooky process
pkill spooky

# Re-export facts with fresh memory
spooky facts export ./my-project --output fresh-facts.json

# Use memory validation
spooky facts export ./my-project --verbose

# Check for memory leaks
spooky facts export ./my-project --parallel 1
```

### Validation Errors

#### "Validation failed: machine_id must be a 32-character hexadecimal string"

**Problem:** The machine ID format is invalid.

**Solution:**
```bash
# Check machine ID format
cat /etc/machine-id

# Regenerate machine ID if needed
sudo rm /etc/machine-id
sudo systemd-machine-id-setup

# Re-collect facts
spooky facts export ./my-project --machine problematic-machine
```

#### "Validation failed: system.hardware facts are required"

**Problem:** Required hardware facts are missing.

**Solution:**
```bash
# Check if gopsutil is working
python3 -c "import psutil; print(psutil.cpu_count())"

# Re-collect facts with verbose output
spooky facts export ./my-project --machine problematic-machine --verbose

# Check for permission issues
sudo spooky facts export ./my-project --machine problematic-machine
```

#### "Validation failed: invalid data type for field"

**Problem:** Collected data doesn't match expected types.

**Solution:**
```bash
# Export facts for inspection
spooky facts export ./my-project --format json --output debug.json

# Check specific field values
cat debug.json | jq '.machines[].facts.system.hardware.cpu.cores'

# Re-collect facts
spooky facts export ./my-project --machine problematic-machine
```

### Performance Issues

#### "Fact collection is very slow"

**Problem:** Fact collection is taking too long.

**Solution:**
```bash
# Use parallel collection
spooky facts export ./my-project --parallel 8

# Increase timeout for slow machines
spooky facts export ./my-project --timeout 180s

# Profile collection performance
spooky facts export ./my-project --profile
```

**Optimize your configuration:**
```hcl
project {
  facts {
    parallel_workers = 8        # Increase parallel workers
    timeout_seconds = 180       # Increase timeout
    retry_attempts = 2          # Reduce retries for faster failure
  }
}
```

#### "High memory usage during fact collection"

**Problem:** Fact collection is consuming too much memory.

**Solution:**
```bash
# Monitor memory usage
top -p $(pgrep spooky)

# Use fewer parallel workers
spooky facts export ./my-project --parallel 2

# Collect facts in smaller batches
spooky facts export ./my-project --batch-size 5
```

## Configuration Problems

### Project Configuration Issues

#### "Invalid facts configuration: unknown field"

**Problem:** Unknown fields in facts configuration.

**Solution:**
```hcl
project {
  facts {
    # Valid fields only
    parallel_workers = 4
    timeout_seconds = 60
    retry_attempts = 3
    compression_enabled = true
  }
}
```

#### "Invalid machine configuration: missing authentication"

**Problem:** Machine configuration is missing authentication.

**Solution:**
```hcl
machines {
  machine "web-server" {
    hostname = "web.example.com"
    port = 22
    user = "admin"
    
    # Authentication is required
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/id_rsa"
    }
  }
}
```

### SSH Configuration Issues

#### "SSH connection failed: no such identity"

**Problem:** SSH key file not found.

**Solution:**
```bash
# Check if key file exists
ls -la ~/.ssh/id_rsa

# Note: spooky does not generate SSH keys - use openssh to generate new SSH key if needed
ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa

# Copy key to target machine
ssh-copy-id -i ~/.ssh/id_rsa user@machine.example.com
```

#### "SSH connection failed: host key verification failed"

**Problem:** SSH host key verification is failing.

**Solution:**
```bash
# Add host key to known_hosts
ssh-keyscan -H machine.example.com >> ~/.ssh/known_hosts

# Or disable host key checking (not recommended for production)
ssh -o StrictHostKeyChecking=no user@machine.example.com
```

## Network Issues

### Connectivity Problems

#### "Network unreachable"

**Problem:** Cannot reach the target machine.

**Solution:**
```bash
# Check network connectivity
ping machine.example.com

# Check DNS resolution
nslookup machine.example.com

# Check routing
traceroute machine.example.com

# Verify firewall rules
sudo iptables -L
```

#### "Connection timed out"

**Problem:** Network connection is timing out.

**Solution:**
```bash
# Check if machine is reachable
telnet machine.example.com 22

# Test SSH connection
ssh -o ConnectTimeout=10 user@machine.example.com

# Check network latency
ping -c 5 machine.example.com
```

### Firewall Issues

#### "Connection refused by firewall"

**Problem:** Firewall is blocking SSH connections.

**Solution:**
```bash
# Check local firewall
sudo ufw status

# Check remote firewall
ssh user@machine.example.com "sudo iptables -L"

# Open SSH port
sudo ufw allow 22/tcp
```

## Data Collection Issues

### gopsutil Problems

#### "Failed to collect CPU information"

**Problem:** Cannot collect CPU facts.

**Solution:**
```bash
# Check if gopsutil is working
python3 -c "import psutil; print(psutil.cpu_count())"

# Check system permissions
sudo spooky facts export ./my-project

# Check for system-specific issues
cat /proc/cpuinfo
```

#### "Failed to collect memory information"

**Problem:** Cannot collect memory facts.

**Solution:**
```bash
# Check memory information manually
free -h

# Check /proc/meminfo
cat /proc/meminfo

# Check for permission issues
sudo spooky facts export ./my-project
```

#### "Failed to collect disk information"

**Problem:** Cannot collect disk facts.

**Solution:**
```bash
# Check disk information manually
df -h

# Check disk partitions
lsblk

# Check for permission issues
sudo spooky facts export ./my-project
```

### Process Collection Issues

#### "Failed to collect process information"

**Problem:** Cannot collect process facts.

**Solution:**
```bash
# Check process information manually
ps aux

# Check for permission issues
sudo spooky facts export ./my-project

# Check if process collection is enabled
spooky facts export ./my-project --include-processes
```

## Export Issues

### Memory Management Problems

#### "Memory allocation failed during export"

**Problem:** Cannot allocate memory for fact export.

**Solution:**
```bash
# Check available memory
free -h

# Check system limits
ulimit -a

# Reduce memory usage
spooky facts export ./my-project --parallel 1

# Export fewer machines at once
spooky facts export ./my-project --machine specific-machine --output single-machine-facts.json
```

#### "Memory corruption detected during export"

**Problem:** Memory corruption detected during fact export.

**Solution:**
```bash
# Restart spooky process
pkill spooky

# Re-export facts with fresh memory
spooky facts export ./my-project

# Use memory validation
spooky facts export ./my-project --validate-memory

# Check for memory leaks
spooky facts export ./my-project --memory-profile
```

### Export System Issues

#### "Failed to allocate facts memory"

**Problem:** Cannot allocate memory for fact export.

**Solution:**
```bash
# Check available memory
free -h

# Check system limits
ulimit -a

# Reduce memory usage
spooky facts export ./my-project --parallel 1

# Export fewer machines at once
spooky facts export ./my-project --machine specific-machine --output single-machine-facts.json
```

## Performance Optimization

### Collection Performance

#### "Fact export is slow for many machines"

**Problem:** Export is slow with many machines.

**Solution:**
```bash
# Use parallel collection
spooky facts export ./my-project --parallel 16

# Use connection pooling
spooky facts export ./my-project --connection-pool 20

# Export in batches
spooky facts export ./my-project --batch-size 10
```

#### "High CPU usage during export"

**Problem:** Fact export is consuming too much CPU.

**Solution:**
```bash
# Reduce parallel workers
spooky facts export ./my-project --parallel 2

# Use less intensive collection
spooky facts export ./my-project --basic-facts-only

# Profile CPU usage
spooky facts export ./my-project --profile
```

### Export Performance

#### "Slow fact access"

**Problem:** Accessing facts from memory is slow.

**Solution:**
```bash
# Check memory performance
vmstat 1

# Optimize memory allocation
spooky facts export ./my-project --memory-optimized

# Use memory pooling
spooky facts export ./my-project --memory-pool
```

#### "High memory usage during fact export"

**Problem:** Exporting facts is causing high memory usage.

**Solution:**
```bash
# Reduce parallel workers
spooky facts export ./my-project --parallel 2

# Export fewer machines at once
spooky facts export ./my-project --machine specific-machine --output single-machine-facts.json

# Monitor memory usage
spooky facts export ./my-project --verbose

# Use smaller batch sizes
spooky facts export ./my-project --machine machine1,machine2 --output batch-facts.json
```

## Debugging Techniques

### Enable Debug Logging

```bash
# Enable debug logging
export SPOOKY_LOG_LEVEL=debug
spooky facts export ./my-project --verbose

# Enable SSH debug logging
export SPOOKY_SSH_DEBUG=true
spooky facts export ./my-project

# Enable export debug logging
export SPOOKY_EXPORT_DEBUG=true
spooky facts export ./my-project
```

### Collect Debug Information

```bash
# Collect system information
spooky facts export ./my-project --verbose

# Export facts for analysis
spooky facts export ./my-project --format json --output debug.json

# Validate facts with detailed output
spooky facts export ./my-project --verbose
```

### Network Diagnostics

```bash
# Test SSH connectivity
spooky machines ping ./my-project

# Check network latency
ping machine.example.com

# Test fact collection on single machine
spooky facts export ./my-project --machine test-machine --verbose
```

## Recovery Procedures

### Memory Recovery

```bash
# Restart spooky process
pkill spooky

# Re-export facts with fresh memory
spooky facts export ./my-project --output fresh-facts.json

# Validate exported facts
spooky facts export ./my-project

# Check for memory leaks
spooky facts export ./my-project --parallel 1
```

### Fact Recollection

```bash
# Export facts for specific machine
spooky facts export ./my-project --machine problematic-machine --output problematic-machine-facts.json

# Re-export facts for specific machine
spooky facts export ./my-project --machine problematic-machine --output recollected-facts.json

# Validate exported facts
spooky facts export ./my-project --machine problematic-machine
```

### Configuration Recovery

```bash
# Backup configuration
cp ./my-project/project.hcl ./my-project/project.hcl.backup

# Validate configuration
spooky project validate ./my-project

# Restore from backup if needed
cp ./my-project/project.hcl.backup ./my-project/project.hcl
```

## Prevention Strategies

### Regular Maintenance

```bash
# Schedule regular fact collection
crontab -e
# Add: 0 2 * * * /usr/local/bin/spooky facts export /path/to/project

# Regular validation
crontab -e
# Add: 0 3 * * * /usr/local/bin/spooky project validate /path/to/project

# Regular memory monitoring
crontab -e
# Add: 0 4 * * * /usr/local/bin/spooky facts export /path/to/project --verbose
```

### Monitoring

```bash
# Monitor fact collection health
spooky facts export ./my-project --verbose

# Monitor memory usage
free -h

# Monitor collection performance
spooky facts export ./my-project --profile
```

### Memory Strategy

```bash
# Regular memory monitoring
crontab -e
# Add: 0 1 * * * free -h >> /var/log/spooky/memory.log

# Backup configuration
crontab -e
# Add: 0 1 * * * cp /path/to/project/project.hcl /backup/project.hcl.$(date +%Y%m%d)
```

## Known Issues and Workarounds

> **See also**: [Known Issues](KNOWN_ISSUES.md#facts-system-ssh-issues) - Comprehensive documentation of all known issues and workarounds

### SSH-Based Collection Issues

**Issue:** SSH-based fact collection has implementation problems.

**Workaround:**
```bash
# Use local fact collection for immediate needs
spooky facts export ./my-project --local-only

# Collect facts manually and import
ssh user@machine.example.com "cat /etc/machine-id" > machine-id.txt
```

### Parallel Collection Issues

**Issue:** Parallel fact collection may fail with many machines.

**Workaround:**
```bash
# Use sequential collection
spooky facts export ./my-project --parallel 1

# Collect facts in smaller batches
spooky facts export ./my-project --machine machine1,machine2 --output batch1.json
spooky facts export ./my-project --machine machine3,machine4 --output batch2.json
```

### Memory Issues

**Issue:** High memory usage during fact export.

**Workaround:**
```bash
# Export facts in smaller batches
spooky facts export ./my-project --machine single-machine --output single.json

# Use lower parallel settings
spooky facts export ./my-project --parallel 2

# Monitor memory usage
watch -n 1 'free -h'
```

## Getting Help

### Enable Debug Logging

```bash
# Enable debug logging
export SPOOKY_LOG_LEVEL=debug
spooky facts export ./my-project --verbose
```

### Collect Diagnostic Information

```bash
# Collect system information
spooky facts export ./my-project --verbose > facts-debug.log 2>&1

# Collect memory information
free -h >> facts-debug.log 2>&1

# Collect SSH information
spooky machines ping ./my-project --verbose >> facts-debug.log 2>&1
```

### Report Issues

When reporting facts issues, include:

1. **Error Message**: Complete error message
2. **Configuration**: Relevant project and machine configuration
3. **System Information**: OS, memory, CPU information
4. **Network**: Network connectivity information
5. **Logs**: Debug logs and verbose output
6. **Steps**: Steps to reproduce the issue

## Conclusion

The facts system provides basic fact collection and export capabilities with some limitations. Most issues can be resolved by following the troubleshooting steps outlined in this guide. For persistent issues, enable verbose output and collect diagnostic information for further analysis. SSH-based fact collection improvements are planned for future releases.
