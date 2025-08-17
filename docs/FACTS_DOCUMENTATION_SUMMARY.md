# Facts System Documentation Summary

## Overview

This document provides a comprehensive overview of the spooky facts system documentation. It serves as a guide to help you find the right documentation for your needs and understand how all the pieces fit together.

**Status: Implemented** - The facts system is fully implemented with comprehensive functionality for fact collection and export.

## Documentation Structure

### 📚 Core Documentation

#### 1. [User Guide](FACTS_USER_GUIDE.md)
**Audience:** End users, system administrators, DevOps engineers
**Purpose:** Complete guide to using the facts system

**What it covers:**
- Getting started with fact collection
- Fact collection and export functionality
- Machine inventory integration
- Export formats and options
- Real-world examples and use cases

**When to use:** Start here if you're new to spooky facts or need to understand how to use the system effectively.

#### 2. [API Reference](FACTS_API_REFERENCE.md)
**Audience:** Developers, system integrators, contributors
**Purpose:** Technical reference for the facts system APIs and implementation

**What it covers:**
- Core interfaces and type definitions
- Implementation details and algorithms
- Error handling patterns
- Configuration rules and schemas
- CLI integration details
- Code examples and patterns

**When to use:** Use this when developing with the facts system, extending functionality, or debugging implementation issues.

#### 3. [Troubleshooting Guide](FACTS_TROUBLESHOOTING.md)
**Audience:** System administrators, support engineers, users experiencing issues
**Purpose:** Solutions for common problems and debugging techniques

**What it covers:**
- Common error messages and solutions
- SSH-based collection issues and workarounds

> **See also**: [Known Issues](KNOWN_ISSUES.md#facts-system-ssh-issues) - Comprehensive documentation of all known issues and workarounds
- Performance issues and optimization
- Export format issues
- Configuration problems and debugging
- Best practices for troubleshooting

**When to use:** Use this when encountering problems or need to debug issues with the facts system.

### 📁 Examples Directory

#### [Examples Overview](examples/README.md)
**Audience:** All users
**Purpose:** Quick reference for available examples and use cases

**What it covers:**
- Available fact collection examples
- Example configurations and scripts
- Common use case patterns
- Integration examples with other systems

**When to use:** Use this to quickly find relevant examples for your use case.

## Key Concepts

### Core Features

1. **Fact Collection** - Collect system information from machines via SSH
2. **Fact Export** - Export facts to various formats (JSON, HCL)
3. **Machine Integration** - Uses project machine inventory for target identification
4. **Multiple Export Formats** - JSON and HCL export formats
5. **Filtering Support** - Filter by machine, tags, and groups
6. **Parallel Collection** - Collect facts from multiple machines concurrently

### Architecture Principles

1. **Interface-First Design** - All functionality through well-defined interfaces
2. **Dependency Injection** - Loose coupling through interface-based dependencies
3. **Direct Export** - Facts are collected and exported directly without intermediate storage
4. **Extensible Design** - Easy to add new collection methods and export formats
5. **Performance Optimized** - Efficient collection and export with parallel processing

### Best Practices

1. **Use Project Machine Inventory** - Leverage existing machine inventory for fact collection
2. **Filter Appropriately** - Use machine and tag filtering to limit collection scope
3. **Choose Export Format** - Use JSON for machine processing, HCL for human readability
4. **Use Parallel Collection** - Enable parallel collection for better performance
5. **Validate Before Collection** - Ensure machine inventory is valid before collection
6. **Monitor Performance** - Use appropriate filtering to avoid overwhelming systems

## Facts System Overview

### Core Concepts

The facts system provides a comprehensive solution for collecting, storing, and exporting system information from machines. Facts are collected via SSH and can include:

- **System information** - OS details, kernel version, architecture
- **Hardware information** - CPU, memory, disk, network interfaces
- **Network configuration** - IP addresses, hostnames, routing
- **Custom data** - User-defined facts and metadata

### Fact Collection Process

1. **Load Machine Inventory** - Read machines from project inventory
2. **Filter Machines** - Apply machine, tag, and group filters
3. **Collect Facts** - Connect to machines via SSH and collect system information
4. **Export Facts** - Export facts directly to requested format (JSON or HCL)

### CLI Commands

The facts system provides comprehensive CLI commands:

```bash
# Export facts from all machines
spooky facts export ./my-project --output facts.hcl

# Export facts to JSON format
spooky facts export ./my-project --format json --output facts.json

# Export facts from specific machines
spooky facts export ./my-project --machine web-server --output web-facts.hcl

# Export facts with parallel collection
spooky facts export ./my-project --parallel 4 --output facts.hcl

# Export facts with filtering
spooky facts export ./my-project --tags environment=production --output prod-facts.hcl
```

### Export Formats

The facts system supports multiple export formats:

#### HCL Format
```hcl
facts {
  fact_collection "web-server" {
    hostname = "web-server.example.com"
    system {
      os_name = "Ubuntu"
      os_version = "22.04"
      kernel_version = "5.15.0-88-generic"
      architecture = "x86_64"
    }
    hardware {
      cpu_cores = 4
      memory_total = "8192MB"
      disk_total = "100GB"
    }
    network {
      interfaces = ["eth0", "eth1"]
      ip_addresses = ["192.168.1.100", "10.0.0.100"]
    }
  }
}
```

#### JSON Format
```json
{
  "facts": {
    "web-server": {
      "hostname": "web-server.example.com",
      "system": {
        "os_name": "Ubuntu",
        "os_version": "22.04",
        "kernel_version": "5.15.0-88-generic",
        "architecture": "x86_64"
      },
      "hardware": {
        "cpu_cores": 4,
        "memory_total": "8192MB",
        "disk_total": "100GB"
      },
      "network": {
        "interfaces": ["eth0", "eth1"],
        "ip_addresses": ["192.168.1.100", "10.0.0.100"]
      }
    }
  }
}
```

### Machine Filtering

Facts can be filtered using various criteria:

```bash
# Filter by machine name
spooky facts export ./my-project --machine web-server

# Filter by tags
spooky facts export ./my-project --tags environment=production

# Filter by groups
spooky facts export ./my-project --groups webservers

# Combine filters
spooky facts export ./my-project --tags environment=production --groups webservers
```

## Implementation Details

### Core Components

1. **Fact Collector** - Collects facts from machines via SSH
2. **Fact Manager** - Manages fact collection and export lifecycle
3. **Fact Integration** - Provides integration with other system components
4. **SSH Manager** - Handles SSH connections for fact collection

### Integration Points

The facts system integrates with:

- **Machines System** - For machine inventory and connectivity
- **SSH System** - For secure connections to remote machines
- **CLI System** - For user interface and command execution

### Error Handling

The facts system provides comprehensive error handling:

- **Connection errors** - SSH connectivity issues
- **Collection errors** - Fact collection failures
- **Export errors** - Format conversion and file I/O issues
- **Validation errors** - Configuration and data validation issues

## Best Practices

### Fact Collection

1. **Use appropriate filtering** to limit collection scope
2. **Enable parallel collection** for better performance
3. **Validate machine inventory** before collection
4. **Monitor collection progress** and handle errors gracefully
5. **Use appropriate timeouts** for SSH connections

### Fact Export

1. **Regularly export facts** to backup important data
2. **Use appropriate export formats** for your use case
3. **Validate exported data** for completeness and accuracy

### Performance Optimization

1. **Use parallel collection** for multiple machines
2. **Filter machines appropriately** to reduce collection scope
3. **Optimize SSH connections** with connection pooling
4. **Monitor resource usage** during collection
5. **Use appropriate timeouts** to prevent hanging connections

## Troubleshooting

### Common Issues

1. **SSH connection failures** - Check machine connectivity and authentication
2. **Fact collection timeouts** - Adjust timeout settings and check network
3. **Export format errors** - Verify format compatibility and file permissions
5. **Filtering issues** - Verify machine names, tags, and group configurations

### Debug Commands

```bash
# Enable verbose logging
export SPOOKY_LOG_LEVEL=debug

# Test machine connectivity
spooky machines ping ./my-project --auth

# Export with verbose output
spooky facts export ./my-project --verbose --output facts.hcl

# Check exported facts files
ls -la *.hcl *.json
```

## Related Documentation

- [Facts User Guide](FACTS_USER_GUIDE.md) - Complete user guide
- [Facts API Reference](FACTS_API_REFERENCE.md) - Technical reference
- [Facts Troubleshooting](FACTS_TROUBLESHOOTING.md) - Troubleshooting guide
- [System Design](../design/systems/facts-system.md) - System design documentation
- [CLI Reference](CLI_REFERENCE.md) - CLI command reference