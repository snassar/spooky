# Facts System Troubleshooting Guide

## Overview

This troubleshooting guide provides solutions for common issues encountered when working with the spooky facts system. It covers error messages, configuration problems, performance issues, and debugging techniques.

## Common Error Messages

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
spooky facts gather ./my-project --machine problematic-machine
```

#### "Validation failed: system.hardware facts are required"

**Problem:** Required hardware facts are missing.

**Solution:**
```bash
# Check if gopsutil is working
python3 -c "import psutil; print(psutil.cpu_count())"

# Re-collect facts with verbose output
spooky facts gather ./my-project --machine problematic-machine --verbose

# Check for permission issues
sudo spooky facts gather ./my-project --machine problematic-machine
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
spooky facts gather ./my-project --machine problematic-machine
```

### Performance Issues

#### "Fact collection is very slow"

**Problem:** Fact collection is taking too long.

**Solution:**
```bash
# Use parallel collection
spooky facts gather ./my-project --parallel 8

# Increase timeout for slow machines
spooky facts gather ./my-project --timeout 180s

# Profile collection performance
spooky facts gather ./my-project --profile
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
spooky facts gather ./my-project --parallel 2

# Collect facts in smaller batches
spooky facts gather ./my-project --batch-size 5
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
sudo spooky facts gather ./my-project

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
sudo spooky facts gather ./my-project
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
sudo spooky facts gather ./my-project
```

### Process Collection Issues

#### "Failed to collect process information"

**Problem:** Cannot collect process facts.

**Solution:**
```bash
# Check process information manually
ps aux

# Check for permission issues
sudo spooky facts gather ./my-project

# Check if process collection is enabled
spooky facts gather ./my-project --include-processes
```

## Storage Issues

### Memory Management Problems

#### "Memory initialization failed"

**Problem:** Cannot initialize memory storage.

**Solution:**
```bash
# Check available memory
free -h

# Check system limits
ulimit -a

# Use memory-efficient mode
spooky facts gather ./my-project --memory-efficient

# Reduce memory usage
spooky facts gather ./my-project --parallel 1
```

#### "Memory corruption detected"

**Problem:** Memory corruption detected in fact storage.

**Solution:**
```bash
# Restart spooky process
pkill spooky

# Re-collect facts with fresh memory
spooky facts gather ./my-project

# Use memory validation
spooky facts gather ./my-project --validate-memory

# Check for memory leaks
spooky facts gather ./my-project --memory-profile
```

### Memory System Issues

#### "Failed to allocate facts memory"

**Problem:** Cannot allocate memory for fact gathering and export.

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

#### "Fact collection is slow for many machines"

**Problem:** Collection is slow with many machines.

**Solution:**
```bash
# Use parallel collection
spooky facts gather ./my-project --parallel 16

# Use connection pooling
spooky facts gather ./my-project --connection-pool 20

# Collect in batches
spooky facts gather ./my-project --batch-size 10
```

#### "High CPU usage during collection"

**Problem:** Fact collection is consuming too much CPU.

**Solution:**
```bash
# Reduce parallel workers
spooky facts gather ./my-project --parallel 2

# Use less intensive collection
spooky facts gather ./my-project --basic-facts-only

# Profile CPU usage
spooky facts gather ./my-project --profile
```

### Storage Performance

#### "Slow fact access"

**Problem:** Accessing facts from memory is slow.

**Solution:**
```bash
# Check memory performance
vmstat 1

# Optimize memory allocation
spooky facts gather ./my-project --memory-optimized

# Use memory pooling
spooky facts gather ./my-project --memory-pool
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
spooky facts gather ./my-project --verbose

# Enable SSH debug logging
export SPOOKY_SSH_DEBUG=true
spooky facts gather ./my-project

# Enable export debug logging
export SPOOKY_EXPORT_DEBUG=true
spooky facts export ./my-project
```

### Collect Debug Information

```bash
# Collect system information
spooky facts debug ./my-project

# Export facts for analysis
spooky facts export ./my-project --format json --output debug.json

# Validate facts with detailed output
spooky facts validate ./my-project --verbose
```

### Network Diagnostics

```bash
# Test SSH connectivity
spooky facts test-ssh ./my-project

# Check network latency
spooky facts ping ./my-project

# Test fact collection on single machine
spooky facts gather ./my-project --machine test-machine --verbose
```

## Recovery Procedures

### Memory Recovery

```bash
# Restart spooky process
pkill spooky

# Re-export facts with fresh memory
spooky facts export ./my-project --output fresh-facts.json

# Validate exported facts
spooky facts validate ./my-project

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
spooky facts validate ./my-project --machine problematic-machine
```

### Configuration Recovery

```bash
# Backup configuration
cp ./my-project/project.hcl ./my-project/project.hcl.backup

# Validate configuration
spooky facts validate-config ./my-project

# Restore from backup if needed
cp ./my-project/project.hcl.backup ./my-project/project.hcl
```

## Prevention Strategies

### Regular Maintenance

```bash
# Schedule regular fact collection
crontab -e
# Add: 0 2 * * * /usr/local/bin/spooky facts gather /path/to/project

# Regular validation
crontab -e
# Add: 0 3 * * * /usr/local/bin/spooky facts validate /path/to/project

# Regular memory monitoring
crontab -e
# Add: 0 4 * * * /usr/local/bin/spooky facts memory-check /path/to/project
```

### Monitoring

```bash
# Monitor fact collection health
spooky facts health-check ./my-project

# Monitor memory usage
spooky facts memory-stats ./my-project

# Monitor collection performance
spooky facts performance-stats ./my-project
```

### Memory Strategy

```bash
# Regular memory monitoring
crontab -e
# Add: 0 1 * * * /usr/local/bin/spooky facts memory-check /path/to/project

# Backup configuration
crontab -e
# Add: 0 1 * * * cp /path/to/project/project.hcl /backup/project.hcl.$(date +%Y%m%d)
```

This comprehensive troubleshooting guide provides solutions for the most common issues encountered with the spooky facts system, along with prevention strategies and recovery procedures.
