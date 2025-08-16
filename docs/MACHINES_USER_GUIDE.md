# Machines System User Guide

## Overview

The spooky machines system provides comprehensive machine inventory management, connectivity testing, and validation capabilities. This guide covers everything from basic machine configuration to advanced features like SSH authentication, connectivity testing, and inventory validation.

**Status: Production Ready** - The machines system is fully implemented with comprehensive inventory management, connectivity testing, and validation capabilities.

## Getting Started

### Prerequisites

- spooky CLI installed and configured
- SSH access to target machines
- Basic understanding of HCL configuration syntax
- Access to create and modify project files

### Quick Start

1. **Check Available Machines Commands**
   ```bash
   spooky machines --help
   ```

2. **List Machines in a Project**
   ```bash
   spooky machines list ./my-project
   ```

3. **Test Machine Connectivity**
   ```bash
   spooky machines ping ./my-project
   ```

4. **Validate Machine Configuration**
   ```bash
   spooky machines validate ./my-project
   ```

## Machines System Concepts

### What are Machines?

Machines represent the servers, workstations, and other systems that spooky can manage. Each machine has:

- **Connection Information**: Hostname, IP address, SSH port
- **Authentication Details**: SSH keys, usernames, passwords
- **Metadata**: Tags, descriptions, environment information
- **Configuration**: SSH settings, timeouts, retry policies

### Machine Inventory

The machine inventory is stored in `machines.hcl` files and provides:

- **Centralized Management**: All machine definitions in one place
- **Tag-Based Organization**: Group machines by environment, role, or purpose
- **SSH Configuration**: SSH connection settings and authentication
- **Validation**: Automatic validation of machine configurations

### Machine Connectivity

The machines system provides comprehensive connectivity testing:

- **DNS Resolution**: Verify hostname resolution
- **Network Connectivity**: Test basic network reachability
- **SSH Connectivity**: Test SSH connection and authentication
- **Authentication Testing**: Verify SSH key and password authentication

## Current Implementation Status

### ✅ Fully Implemented Features

- **Complete Machine Inventory Management**: Full machine inventory loading and validation
- **SSH Connectivity Testing**: Comprehensive SSH connection testing with authentication
- **Machine Validation**: Complete machine configuration validation
- **CLI Integration**: All `spooky machines` commands fully functional
- **Project Integration**: Seamless integration with project configuration
- **SSH Authentication**: Support for SSH keys, passwords, and certificates
- **Tag-Based Filtering**: Machine filtering by tags and metadata
- **Export Functionality**: Machine inventory export to JSON format
- **Error Handling**: Comprehensive error handling and reporting

### 🎯 Production Ready

The machines system is **production-ready** with:
- **100% Functional Infrastructure**: No stubs or placeholders
- **Complete SSH Integration**: Full SSH connectivity testing
- **Robust Error Handling**: Comprehensive error recovery and reporting
- **Performance Optimized**: Efficient connectivity testing with proper timeouts
- **Type Safe**: All interface contracts satisfied with proper validation

## Basic Usage

### Listing Machines

List all machines in a project:

```bash
# List all machines
spooky machines list ./my-project
```

**Example Output**:
```
Machines in project (3 found):
1. web-server (web.example.com:22) - admin@web.example.com
2. db-server (db.example.com:22) - admin@db.example.com
3. cache-server (cache.example.com:22) - admin@cache.example.com
```

### Testing Connectivity

Test connectivity to machines:

```bash
# Test connectivity to all machines
spooky machines ping ./my-project

# Test connectivity with authentication
spooky machines ping ./my-project --auth
```

**Example Output**:
```
Testing connectivity to 3 machines...

web-server (web.example.com:22):
  ✅ DNS Resolution: OK
  ✅ Network Connectivity: OK (5ms)
  ✅ SSH Connection: OK
  ✅ Authentication: OK

db-server (db.example.com:22):
  ✅ DNS Resolution: OK
  ✅ Network Connectivity: OK (3ms)
  ✅ SSH Connection: OK
  ✅ Authentication: OK

cache-server (cache.example.com:22):
  ❌ DNS Resolution: FAILED
  ❌ Network Connectivity: FAILED
  ❌ SSH Connection: FAILED
  ❌ Authentication: FAILED

Summary: 2/3 machines reachable
```

### Validating Configuration

Validate machine configurations:

```bash
# Validate all machines
spooky machines validate ./my-project
```

**Example Output**:
```
✅ Machine configuration validation passed

Validated 3 machines:
- web-server: ✅ Valid
- db-server: ✅ Valid
- cache-server: ✅ Valid

All machines have valid SSH configuration and authentication methods.
```

### Exporting Inventory

Export machine inventory to JSON:

```bash
# Export machine inventory
spooky machines export ./my-project --output inventory.json
```

## Project Configuration

### Basic Machine Configuration

Create machine inventory in your project:

```hcl
# machines.hcl
machines {
  machine "web-server" {
    hostname = "web.example.com"
    host = "192.168.1.10"
    user = "admin"
    port = 22
    
    tags = {
      environment = "production"
      role = "web"
      datacenter = "us-east-1"
    }
  }
  
  machine "db-server" {
    hostname = "db.example.com"
    host = "192.168.1.11"
    user = "admin"
    port = 22
    
    tags = {
      environment = "production"
      role = "database"
      datacenter = "us-east-1"
    }
  }
  
  machine "cache-server" {
    hostname = "cache.example.com"
    host = "192.168.1.12"
    user = "admin"
    port = 22
    
    tags = {
      environment = "production"
      role = "cache"
      datacenter = "us-east-1"
    }
  }
}
```

### SSH Authentication Configuration

Configure SSH authentication for machines:

```hcl
# machines.hcl
machines {
  machine "web-server" {
    hostname = "web.example.com"
    user = "admin"
    
    # SSH key authentication
    key_file = "~/.ssh/id_ed25519"
    
    # Optional passphrase for encrypted keys
    passphrase = "your-passphrase"
    
    tags = {
      environment = "production"
    }
  }
  
  machine "db-server" {
    hostname = "db.example.com"
    user = "admin"
    
    # Password authentication
    password = "your-password"
    
    tags = {
      environment = "production"
    }
  }
}
```

### Advanced SSH Configuration

Configure advanced SSH settings:

```hcl
# machines.hcl
machines {
  machine "web-server" {
    hostname = "web.example.com"
    user = "admin"
    port = 2222  # Custom SSH port
    
    # SSH key with custom path
    key_file = "~/.ssh/custom_key"
    
    # SSH connection settings
    ssh_config {
      connect_timeout = 30
      command_timeout = 60
      retry_attempts = 3
      retry_delay = 5
    }
    
    tags = {
      environment = "production"
    }
  }
}
```

### Machine Tags and Metadata

Use tags to organize and filter machines:

```hcl
# machines.hcl
machines {
  machine "web-prod-1" {
    hostname = "web-prod-1.example.com"
    user = "admin"
    key_file = "~/.ssh/id_ed25519"
    
    tags = {
      environment = "production"
      role = "web"
      datacenter = "us-east-1"
      tier = "frontend"
    }
  }
  
  machine "web-prod-2" {
    hostname = "web-prod-2.example.com"
    user = "admin"
    key_file = "~/.ssh/id_ed25519"
    
    tags = {
      environment = "production"
      role = "web"
      datacenter = "us-east-1"
      tier = "frontend"
    }
  }
  
  machine "db-prod-1" {
    hostname = "db-prod-1.example.com"
    user = "admin"
    key_file = "~/.ssh/id_ed25519"
    
    tags = {
      environment = "production"
      role = "database"
      datacenter = "us-east-1"
      tier = "backend"
    }
  }
}
```

## Advanced Usage

### Environment-Specific Configuration

Use different configurations for different environments:

```hcl
# machines.hcl
machines {
  # Production machines
  machine "web-prod" {
    hostname = "web-prod.example.com"
    user = "admin"
    key_file = "~/.ssh/prod_key"
    
    tags = {
      environment = "production"
      role = "web"
    }
  }
  
  machine "db-prod" {
    hostname = "db-prod.example.com"
    user = "admin"
    key_file = "~/.ssh/prod_key"
    
    tags = {
      environment = "production"
      role = "database"
    }
  }
  
  # Staging machines
  machine "web-staging" {
    hostname = "web-staging.example.com"
    user = "admin"
    key_file = "~/.ssh/staging_key"
    
    tags = {
      environment = "staging"
      role = "web"
    }
  }
  
  machine "db-staging" {
    hostname = "db-staging.example.com"
    user = "admin"
    key_file = "~/.ssh/staging_key"
    
    tags = {
      environment = "staging"
      role = "database"
    }
  }
}
```

### SSH Certificate Authentication

Use SSH certificates for authentication:

```hcl
# machines.hcl
machines {
  machine "web-server" {
    hostname = "web.example.com"
    user = "admin"
    
    # SSH certificate authentication
    certificate_file = "~/.ssh/cert.pub"
    key_file = "~/.ssh/id_ed25519"
    
    tags = {
      environment = "production"
    }
  }
}
```

### Machine Groups and Patterns

Use naming patterns for machine organization:

```hcl
# machines.hcl
machines {
  # Web servers
  machine "web-01" {
    hostname = "web-01.example.com"
    user = "admin"
    key_file = "~/.ssh/id_ed25519"
    
    tags = {
      environment = "production"
      role = "web"
      instance = "01"
    }
  }
  
  machine "web-02" {
    hostname = "web-02.example.com"
    user = "admin"
    key_file = "~/.ssh/id_ed25519"
    
    tags = {
      environment = "production"
      role = "web"
      instance = "02"
    }
  }
  
  # Database servers
  machine "db-primary" {
    hostname = "db-primary.example.com"
    user = "admin"
    key_file = "~/.ssh/id_ed25519"
    
    tags = {
      environment = "production"
      role = "database"
      type = "primary"
    }
  }
  
  machine "db-replica" {
    hostname = "db-replica.example.com"
    user = "admin"
    key_file = "~/.ssh/id_ed25519"
    
    tags = {
      environment = "production"
      role = "database"
      type = "replica"
    }
  }
}
```

## Connectivity Testing

### Basic Connectivity Testing

Test basic connectivity to machines:

```bash
# Test all machines
spooky machines ping ./my-project

# Test specific machines
spooky machines ping ./my-project --machine web-server

# Test machines by tags
spooky machines ping ./my-project --tags environment=production
```

### Authentication Testing

Test SSH authentication:

```bash
# Test authentication for all machines
spooky machines ping ./my-project --auth

# Test authentication for specific machines
spooky machines ping ./my-project --machine web-server --auth
```

### Detailed Connectivity Information

Get detailed connectivity information:

```bash
# Test with verbose output
spooky machines ping ./my-project --verbose
```

**Example Output**:
```
Testing connectivity to 3 machines...

web-server (web.example.com:22):
  ✅ DNS Resolution: OK (web.example.com -> 192.168.1.10)
  ✅ Network Connectivity: OK (5ms, 64 bytes)
  ✅ SSH Connection: OK (SSH-2.0-OpenSSH_8.2p1)
  ✅ Authentication: OK (publickey)
  ✅ User Access: OK (admin)

db-server (db.example.com:22):
  ✅ DNS Resolution: OK (db.example.com -> 192.168.1.11)
  ✅ Network Connectivity: OK (3ms, 64 bytes)
  ✅ SSH Connection: OK (SSH-2.0-OpenSSH_8.2p1)
  ✅ Authentication: OK (publickey)
  ✅ User Access: OK (admin)

cache-server (cache.example.com:22):
  ❌ DNS Resolution: FAILED (NXDOMAIN)
  ❌ Network Connectivity: FAILED (no route to host)
  ❌ SSH Connection: FAILED
  ❌ Authentication: FAILED
  ❌ User Access: FAILED

Summary: 2/3 machines reachable
Total time: 15.2s
```

## Troubleshooting

### Common Issues

#### SSH Connection Failures

**Problem**: SSH connections fail during connectivity testing

**Solutions**:
1. Verify SSH connectivity manually:
   ```bash
   ssh admin@web.example.com
   ```

2. Check SSH key permissions:
   ```bash
   chmod 600 ~/.ssh/id_ed25519
   ```

3. Verify machine inventory configuration:
   ```bash
   spooky machines validate ./my-project
   ```

4. Check SSH service status on target machine:
   ```bash
   sudo systemctl status ssh
   ```

#### DNS Resolution Issues

**Problem**: Hostnames cannot be resolved

**Solutions**:
1. Check DNS configuration:
   ```bash
   nslookup web.example.com
   ```

2. Use IP addresses instead of hostnames:
   ```hcl
   machine "web-server" {
     hostname = "192.168.1.10"  # Use IP instead of hostname
     user = "admin"
   }
   ```

3. Check `/etc/hosts` file for local resolution

#### Authentication Failures

**Problem**: SSH authentication fails

**Solutions**:
1. Verify SSH key is in authorized_keys:
   ```bash
   ssh-copy-id admin@web.example.com
   ```

2. Check SSH key format:
   ```bash
   ssh-keygen -l -f ~/.ssh/id_ed25519
   ```

3. Test SSH key manually:
   ```bash
   ssh -i ~/.ssh/id_ed25519 admin@web.example.com
   ```

#### Configuration Validation Errors

**Problem**: Machine configuration validation fails

**Solutions**:
1. Check required fields:
   ```hcl
   machine "web-server" {
     hostname = "web.example.com"  # Required
     user = "admin"                # Required
     # Authentication method required (key_file, password, or certificate_file)
     key_file = "~/.ssh/id_ed25519"
   }
   ```

2. Verify HCL syntax:
   ```bash
   spooky machines validate ./my-project --verbose
   ```

3. Check file permissions:
   ```bash
   chmod 600 ~/.ssh/id_ed25519
   ```

### Debugging

Enable verbose output for debugging:

```bash
# Enable debug logging
export SPOOKY_LOG_LEVEL=debug

# Test with verbose output
spooky machines ping ./my-project --verbose
```

## Best Practices

### Machine Organization

1. **Use Descriptive Names**: Choose clear, descriptive machine names
2. **Organize with Tags**: Use tags to group machines by environment, role, or purpose
3. **Consistent Naming**: Use consistent naming patterns for similar machines
4. **Documentation**: Include descriptions for complex configurations

### Security Considerations

1. **SSH Key Management**: Use dedicated SSH keys for different environments
2. **Key Permissions**: Ensure SSH keys have correct permissions (600)
3. **Network Security**: Use VPN or secure networks for machine access
4. **Access Control**: Limit SSH access to necessary users and IPs

### Performance Optimization

1. **Efficient Connectivity Testing**: Use appropriate timeouts for connectivity testing
2. **Parallel Testing**: Test multiple machines in parallel when possible
3. **Caching**: Cache connectivity results for frequently accessed machines
4. **Monitoring**: Monitor connectivity patterns and performance

## Integration with Other Systems

### Actions Integration

Use machine inventory in actions:

```hcl
# actions.hcl
actions {
  action "update-web-servers" {
    type = "command"
    command = "apt update && apt upgrade -y"
    
    machines = ["web-01", "web-02"]  # Use machine names from inventory
  }
  
  action "update-production" {
    type = "command"
    command = "systemctl restart nginx"
    
    # Use tag-based targeting
    tags = ["environment=production", "role=web"]
  }
}
```

### Facts Integration

Use machine inventory for fact collection:

```bash
# Export facts from specific machines
spooky facts export ./my-project --machine web-server --output web-facts.hcl

# Export facts from machines with specific tags
spooky facts export ./my-project --tags environment=production --output prod-facts.hcl
```

### Variables Integration

Use machine information in variables:

```hcl
# variables.hcl
variables {
  web_servers = ["web-01", "web-02"]
  db_servers = ["db-primary", "db-replica"]
  
  # Use machine tags for dynamic configuration
  production_machines = "${machines.tags.environment=production}"
}
```

## CLI Reference

### `spooky machines list`

List machines in a project.

**Syntax**:
```bash
spooky machines list <project-path>
```

**Examples**:
```bash
# List all machines
spooky machines list ./my-project
```

### `spooky machines validate`

Validate machine configurations.

**Syntax**:
```bash
spooky machines validate <project-path>
```

**Examples**:
```bash
# Validate all machines
spooky machines validate ./my-project
```

### `spooky machines ping`

Test connectivity to machines.

**Syntax**:
```bash
spooky machines ping <project-path> [options]
```

**Options**:
- `--machine <list>` - Test specific machines
- `--tags <list>` - Test machines by tags
- `--auth` - Test SSH authentication
- `--verbose` - Show detailed output

**Examples**:
```bash
# Test all machines
spooky machines ping ./my-project

# Test specific machines
spooky machines ping ./my-project --machine web-server

# Test with authentication
spooky machines ping ./my-project --auth

# Test with verbose output
spooky machines ping ./my-project --verbose
```

### `spooky machines export`

Export machine inventory to JSON.

**Syntax**:
```bash
spooky machines export <project-path> --output <file>
```

**Examples**:
```bash
# Export machine inventory
spooky machines export ./my-project --output inventory.json
```

## Remember

**Good machines system usage enables:**
- Efficient machine inventory management
- Reliable connectivity testing
- Secure SSH authentication
- Integration with other spooky systems
- Scalable infrastructure management

**The machines system is production-ready and provides comprehensive machine management capabilities.**
