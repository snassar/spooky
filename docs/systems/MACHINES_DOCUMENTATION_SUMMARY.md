# Machines System Documentation Summary

## Overview

This document provides a comprehensive overview of the spooky machines system documentation. It serves as a guide to help you find the right documentation for your needs and understand how all the pieces fit together.

**Status: Implemented** - The machines system is fully implemented with comprehensive functionality for machine inventory management, connectivity testing, and export capabilities.

## Documentation Structure

### 📚 Core Documentation

#### 1. [User Guide](MACHINES_USER_GUIDE.md)
**Audience:** End users, system administrators, DevOps engineers
**Purpose:** Complete guide to using the machines system

**What it covers:**
- Getting started with machine inventory
- Machine configuration and management
- Connectivity testing and validation
- Export and filtering capabilities
- Authentication and security
- Real-world examples and use cases

**When to use:** Start here if you're new to spooky machines or need to understand how to use the system effectively.

#### 2. [API Reference](MACHINES_API_REFERENCE.md)
**Audience:** Developers, system integrators, contributors
**Purpose:** Technical reference for the machines system APIs and implementation

**What it covers:**
- Core interfaces and type definitions
- Implementation details and algorithms
- Error handling patterns
- Configuration rules and schemas
- CLI integration details
- Code examples and patterns

**When to use:** Use this when developing with the machines system, extending functionality, or debugging implementation issues.

#### 3. [Troubleshooting Guide](MACHINES_TROUBLESHOOTING.md)
**Audience:** System administrators, support engineers, users experiencing issues
**Purpose:** Solutions for common problems and debugging techniques

**What it covers:**
- Common error messages and solutions
- SSH connectivity issues and workarounds

> **See also**: [Known Issues](KNOWN_ISSUES.md#ssh-integration-issues) - Comprehensive documentation of all known issues and workarounds
- Authentication problems and debugging
- Export and filtering issues
- Configuration problems and debugging
- Best practices for troubleshooting

**When to use:** Use this when encountering problems or need to debug issues with the machines system.

### 📁 Examples Directory

#### [Examples Overview](examples/README.md)
**Audience:** All users
**Purpose:** Quick reference for available examples and use cases

**What it covers:**
- Available machine inventory examples
- Example configurations and scripts
- Common use case patterns
- Integration examples with other systems

**When to use:** Use this to quickly find relevant examples for your use case.

## Key Concepts

### Core Features

1. **Machine Inventory Management** - Define and manage machine inventories in HCL format
2. **SSH Connectivity Testing** - Test SSH connectivity to machines (not ICMP ping)
3. **Machine Export** - Export machine inventory to JSON format
4. **Authentication Support** - Multiple authentication methods (SSH keys, passwords)
5. **Filtering and Grouping** - Filter machines by tags, groups, and custom criteria
6. **Validation** - Comprehensive validation of machine configurations
7. **Integration** - Seamless integration with other spooky systems

### Architecture Principles

1. **Interface-First Design** - All functionality through well-defined interfaces
2. **Dependency Injection** - Loose coupling through interface-based dependencies
3. **Inventory-Based** - Machine inventory as the source of truth
4. **Extensible Design** - Easy to add new authentication methods and features
5. **Performance Optimized** - Efficient connectivity testing and export

### Best Practices

1. **Use Descriptive Names** - Use clear, descriptive machine names
2. **Organize with Tags** - Use tags to organize and filter machines
3. **Secure Authentication** - Use SSH keys for authentication when possible
4. **Validate Configurations** - Always validate machine configurations
5. **Test Connectivity** - Regularly test SSH connectivity to machines
6. **Backup Inventory** - Export and backup machine inventory regularly

## Machines System Overview

### Core Concepts

The machines system provides a comprehensive solution for managing machine inventories and testing connectivity. Machines are defined in HCL format and can include:

- **Basic Information** - Hostname, IP address, port, user
- **Authentication** - SSH keys, passwords, authentication methods
- **Metadata** - Tags, groups, descriptions, custom attributes
- **Connectivity** - SSH configuration, timeouts, retry settings

### Machine Inventory Structure

Machine inventories are defined in `machines.hcl` files:

```hcl
machines {
  machine "web-server" {
    hostname = "web.example.com"
    port = 22
    user = "admin"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/id_rsa"
    }
    
    tags = {
      environment = "production"
      role = "web"
      datacenter = "us-west"
    }
    
    groups = ["webservers", "production"]
    
    metadata {
      description = "Primary web server"
      owner = "web-team"
      maintenance_window = "Sunday 2-4 AM UTC"
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
    
    tags = {
      environment = "production"
      role = "database"
      datacenter = "us-west"
    }
    
    groups = ["databases", "production"]
  }
}
```

### CLI Commands

The machines system provides comprehensive CLI commands:

```bash
# Test SSH connectivity to all machines
spooky machines ping ./my-project

# Test connectivity with authentication
spooky machines ping ./my-project --auth

# Test connectivity to specific machines
spooky machines ping ./my-project --machine web-server

# Test connectivity with filtering
spooky machines ping ./my-project --tags environment=production

# Export machine inventory to JSON
spooky machines export ./my-project --output machines.json

# Export with filtering
spooky machines export ./my-project --tags role=web --output web-servers.json

# Validate machine inventory
spooky machines validate ./my-project
```

### Authentication Methods

The machines system supports multiple authentication methods:

#### SSH Key Authentication
```hcl
authentication {
  method = "ssh_key"
  key_path = "~/.ssh/id_rsa"
  passphrase = "optional_passphrase"
}
```

#### Password Authentication
```hcl
authentication {
  method = "password"
  password = "user_password"
}
```

#### Agent Authentication
```hcl
authentication {
  method = "agent"
}
```

### Machine Filtering

Machines can be filtered using various criteria:

```bash
# Filter by machine name
spooky machines ping ./my-project --machine web-server

# Filter by tags
spooky machines ping ./my-project --tags environment=production

# Filter by groups
spooky machines ping ./my-project --groups webservers

# Combine filters
spooky machines ping ./my-project --tags environment=production --groups webservers
```

### Export Formats

The machines system exports to JSON format:

```json
{
  "machines": {
    "web-server": {
      "hostname": "web.example.com",
      "port": 22,
      "user": "admin",
      "authentication": {
        "method": "ssh_key",
        "key_path": "~/.ssh/id_rsa"
      },
      "tags": {
        "environment": "production",
        "role": "web",
        "datacenter": "us-west"
      },
      "groups": ["webservers", "production"],
      "metadata": {
        "description": "Primary web server",
        "owner": "web-team",
        "maintenance_window": "Sunday 2-4 AM UTC"
      }
    }
  }
}
```

## Implementation Details

### Core Components

1. **Machine Loader** - Loads machine inventory from HCL files
2. **Machine Validator** - Validates machine configurations
3. **SSH Manager** - Handles SSH connections and authentication
4. **Machine Integration** - Provides integration with other system components
5. **Export Manager** - Handles machine inventory export

### Integration Points

The machines system integrates with:

- **SSH System** - For connectivity testing and authentication
- **Facts System** - For machine-specific fact collection
- **Actions System** - For machine-specific action execution
- **CLI System** - For user interface and command execution

### Error Handling

The machines system provides comprehensive error handling:

- **Configuration errors** - Invalid machine configurations
- **Authentication errors** - SSH authentication failures
- **Connectivity errors** - Network and SSH connection issues
- **Validation errors** - Schema and data validation issues
- **Export errors** - File I/O and format conversion issues

## Best Practices

### Machine Configuration

1. **Use descriptive names** for easy identification
2. **Organize with tags** for flexible filtering
3. **Use SSH keys** for secure authentication
4. **Set appropriate timeouts** for connectivity testing
5. **Include metadata** for better organization
6. **Validate configurations** before use

### Connectivity Testing

1. **Test regularly** to ensure machines are accessible
2. **Use authentication testing** to verify credentials
3. **Filter appropriately** to test specific machine groups
4. **Monitor results** and address connectivity issues
5. **Use appropriate timeouts** to prevent hanging tests

### Security

1. **Use SSH keys** instead of passwords when possible
2. **Secure key storage** with appropriate permissions
3. **Use passphrases** for additional security
4. **Limit access** to machine inventory files
5. **Regularly rotate** authentication credentials

## Troubleshooting

### Common Issues

1. **SSH connection failures** - Check network connectivity and SSH configuration
2. **Authentication errors** - Verify SSH keys and passwords
3. **Configuration errors** - Validate machine inventory syntax
4. **Export errors** - Check file permissions and disk space
5. **Filtering issues** - Verify machine names, tags, and group configurations

### Debug Commands

```bash
# Enable verbose logging
export SPOOKY_LOG_LEVEL=debug

# Test SSH connectivity manually
ssh -i ~/.ssh/id_rsa user@hostname

# Validate machine inventory
spooky machines validate ./my-project --verbose

# Test with authentication
spooky machines ping ./my-project --auth --verbose

# Check SSH agent
ssh-add -l
```

### SSH Troubleshooting

1. **Check SSH configuration** on target machines
2. **Verify SSH key permissions** (should be 600)
3. **Test SSH connectivity manually** before using spooky
4. **Check SSH agent** for loaded keys
5. **Verify user permissions** on target machines

## Related Documentation

- [Machines User Guide](MACHINES_USER_GUIDE.md) - Complete user guide
- [Machines API Reference](MACHINES_API_REFERENCE.md) - Technical reference
- [Machines Troubleshooting](MACHINES_TROUBLESHOOTING.md) - Troubleshooting guide
- [System Design](../design/systems/machines-system.md) - System design documentation
- [CLI Reference](CLI_REFERENCE.md) - CLI command reference
