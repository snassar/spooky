# Facts System Documentation Summary

## Overview

This document provides a comprehensive overview of the spooky facts system documentation. It serves as a guide to help you find the right documentation for your needs and understand how all the pieces fit together.

## Documentation Structure

### 📚 Core Documentation

#### 1. [User Guide](FACTS_USER_GUIDE.md)
**Audience:** End users, system administrators, DevOps engineers
**Purpose:** Complete guide to using the facts system

**What it covers:**
- Getting started with fact collection
- Basic and advanced usage patterns
- Project and machine configuration
- Integration with other spooky components
- Monitoring and maintenance
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
- Configuration problems and fixes
- Performance issues and optimization
- Network and connectivity issues
- Storage and database problems
- Recovery procedures and prevention strategies

**When to use:** Use this when encountering problems or need to debug issues with the facts system.

### 📁 Examples Directory

#### [Examples Overview](examples/README.md)
**Audience:** All users
**Purpose:** Practical examples and configuration patterns

**What it covers:**
- Basic fact collection setup
- Project configuration examples
- Machine-specific configurations
- Integration examples with other components
- Best practices and patterns

**Example Files:**
- [`basic-facts-project.hcl`](examples/basic-facts-project.hcl) - Basic fact collection project
- [`advanced-facts-config.hcl`](examples/advanced-facts-config.hcl) - Advanced configuration
- [`facts-integration.hcl`](examples/facts-integration.hcl) - Integration with other components

**When to use:** Use these as starting points for your own configurations or to learn best practices.

## Quick Start Guide

### For New Users

1. **Read the User Guide** - Start with [FACTS_USER_GUIDE.md](FACTS_USER_GUIDE.md) to understand the basics
2. **Try the Examples** - Copy and customize examples from the [examples/](examples/) directory
3. **Test Your Setup** - Use `spooky facts gather` and `spooky facts validate` to test
4. **Check Troubleshooting** - If you encounter issues, refer to [FACTS_TROUBLESHOOTING.md](FACTS_TROUBLESHOOTING.md)

### For Developers

1. **Review the API Reference** - Understand the interfaces and implementation in [FACTS_API_REFERENCE.md](FACTS_API_REFERENCE.md)
2. **Study the Examples** - See how the APIs are used in practice
3. **Check the Code** - Review the actual implementation in `internal/facts/`
4. **Test Your Changes** - Use the examples to test your modifications

### For System Administrators

1. **Start with User Guide** - Understand the system capabilities
2. **Review Examples** - See real-world configuration patterns
3. **Plan Your Fact Collection** - Design your fact collection strategy
4. **Implement Gradually** - Start with basic collection and expand
5. **Monitor and Validate** - Use validation and testing regularly

## Documentation Navigation

### By Use Case

#### Getting Started
- **New to spooky facts?** → [User Guide](FACTS_USER_GUIDE.md) - Getting Started section
- **Setting up your first project?** → [User Guide](FACTS_USER_GUIDE.md) - Basic Usage section
- **Need examples?** → [Examples Directory](examples/) - Basic examples

#### Configuration
- **Project configuration?** → [User Guide](FACTS_USER_GUIDE.md) - Project Configuration section
- **Machine-specific settings?** → [User Guide](FACTS_USER_GUIDE.md) - Machine-Specific Configuration section
- **Advanced configuration?** → [User Guide](FACTS_USER_GUIDE.md) - Advanced Usage section

#### Integration
- **Using facts in variables?** → [User Guide](FACTS_USER_GUIDE.md) - Variables Integration section
- **Using facts in templates?** → [User Guide](FACTS_USER_GUIDE.md) - Templates Integration section
- **Using facts in actions?** → [User Guide](FACTS_USER_GUIDE.md) - Actions Integration section

#### Troubleshooting
- **Collection failures?** → [Troubleshooting Guide](FACTS_TROUBLESHOOTING.md) - Collection Errors section
- **Storage issues?** → [Troubleshooting Guide](FACTS_TROUBLESHOOTING.md) - Storage Errors section
- **Performance problems?** → [Troubleshooting Guide](FACTS_TROUBLESHOOTING.md) - Performance Issues section

#### Development
- **Extending the system?** → [API Reference](FACTS_API_REFERENCE.md) - Core Interfaces section
- **Adding new collectors?** → [API Reference](FACTS_API_REFERENCE.md) - FactCollector Interface section
- **Custom storage backends?** → [API Reference](FACTS_API_REFERENCE.md) - FactStorage Interface section

### By Component

#### Fact Collection
- **Overview** → [User Guide](FACTS_USER_GUIDE.md) - Facts System Concepts
- **Configuration** → [User Guide](FACTS_USER_GUIDE.md) - Gathering Facts
- **Troubleshooting** → [Troubleshooting Guide](FACTS_TROUBLESHOOTING.md) - Collection Errors
- **API** → [API Reference](FACTS_API_REFERENCE.md) - FactCollector Interface

#### Fact Storage
- **Overview** → [User Guide](FACTS_USER_GUIDE.md) - Fact Storage
- **Configuration** → [User Guide](FACTS_USER_GUIDE.md) - Storage Management
- **Troubleshooting** → [Troubleshooting Guide](FACTS_TROUBLESHOOTING.md) - Storage Issues
- **API** → [API Reference](FACTS_API_REFERENCE.md) - FactStorage Interface

#### Fact Validation
- **Overview** → [User Guide](FACTS_USER_GUIDE.md) - Validating Facts
- **Configuration** → [User Guide](FACTS_USER_GUIDE.md) - Validation Settings
- **Troubleshooting** → [Troubleshooting Guide](FACTS_TROUBLESHOOTING.md) - Validation Errors
- **API** → [API Reference](FACTS_API_REFERENCE.md) - Validation Rules

#### CLI Commands
- **Overview** → [User Guide](FACTS_USER_GUIDE.md) - Basic Usage
- **Commands** → [User Guide](FACTS_USER_GUIDE.md) - Command Options
- **Troubleshooting** → [Troubleshooting Guide](FACTS_TROUBLESHOOTING.md) - CLI Issues
- **API** → [API Reference](FACTS_API_REFERENCE.md) - CLI Integration

## Key Concepts

### Facts System Architecture

The facts system consists of several key components:

1. **FactCollector** - Collects system information from machines
2. **FactStorage** - Stores and retrieves facts from persistent storage
3. **FactManager** - Orchestrates collection, storage, and validation
4. **FactsIntegration** - Provides integration with other spooky components
5. **CLI Commands** - User interface for fact management

### Data Flow

1. **Machine Discovery** - Read machine inventory from project configuration
2. **SSH Connection** - Establish secure connection to target machine
3. **Fact Collection** - Use gopsutil to gather system information
4. **Data Processing** - Convert and validate collected data
5. **In-Memory Storage** - Store facts in memory for the duration of the action run
6. **Integration** - Make facts available to other spooky components

### Fact Types

The system collects various types of facts:

- **System Facts** - OS, hardware, network information
- **Enhanced Facts** - Additional system details
- **Application Facts** - Application-specific information
- **Custom Facts** - User-defined facts and metadata

## Common Patterns

### Basic Fact Collection

```bash
# Collect facts from all machines
spooky facts gather ./my-project

# List stored facts
spooky facts list ./my-project

# Validate facts
spooky facts validate ./my-project
```

### Advanced Configuration

```hcl
project {
  facts {
    parallel_workers = 8
    timeout_seconds = 120
    storage_path = "facts.db"
    compression_enabled = true
  }
}
```

### Integration with Other Components

```hcl
# Use facts in variables
variables {
  cpu_cores = "{{ facts.web-server.system.hardware.cpu.cores }}"
}

# Use facts in templates
template {
  source = "deploy.sh.tmpl"
  data = {
    machine = "web-server"
    facts = "{{ facts.web-server }}"
  }
}
```

## Best Practices

### Fact Collection

- **Regular Collection** - Collect facts regularly (daily or weekly)
- **Parallel Processing** - Use parallel collection for multiple machines
- **Error Handling** - Monitor collection failures and retry
- **Validation** - Always validate collected facts

### Memory Management

- **Memory Monitoring** - Monitor memory usage during fact collection
- **Garbage Collection** - Proper cleanup of temporary objects
- **Memory Pooling** - Reuse memory structures for better performance
- **Memory Limits** - Set appropriate memory limits for large collections

### Security

- **Memory Protection** - Facts stored in memory are not persisted to disk
- **Access Control** - Facts are only accessible during the action run
- **Audit Logging** - Log fact collection and usage
- **Data Sanitization** - Remove sensitive information before collection

### Performance

- **Parallel Collection** - Use parallel workers for faster collection
- **Timeout Tuning** - Adjust timeouts based on network conditions
- **Memory Tuning** - Optimize memory allocation and usage
- **Memory Management** - Monitor memory usage during collection

## Troubleshooting Quick Reference

### Common Issues

| Issue | Quick Fix | Full Solution |
|-------|-----------|---------------|
| SSH connection failed | Check SSH key permissions | [Troubleshooting Guide](FACTS_TROUBLESHOOTING.md) - Collection Errors |
| Fact collection timeout | Increase timeout | [Troubleshooting Guide](FACTS_TROUBLESHOOTING.md) - Performance Issues |
| Memory allocation failed | Reduce parallel workers | [Troubleshooting Guide](FACTS_TROUBLESHOOTING.md) - Memory Issues |
| Validation errors | Check fact format | [Troubleshooting Guide](FACTS_TROUBLESHOOTING.md) - Validation Errors |
| Slow collection | Use parallel workers | [Troubleshooting Guide](FACTS_TROUBLESHOOTING.md) - Performance Optimization |

### Debug Commands

```bash
# Enable debug logging
export SPOOKY_LOG_LEVEL=debug
spooky facts gather ./my-project --verbose

# Test SSH connectivity
spooky facts test-ssh ./my-project

# Export facts for inspection
spooky facts export ./my-project --format json --output debug.json

# Validate with detailed output
spooky facts validate ./my-project --verbose
```

## Integration Points

### With Other Spooky Components

1. **Variables System** - Facts can be used in variable resolution
2. **Templates System** - Facts can be used in template rendering
3. **Actions System** - Facts can be used in action conditions
4. **Machines System** - Facts are collected from machine inventory
5. **SSH System** - Facts collection uses SSH connections

### With External Systems

1. **Monitoring Systems** - Export facts for monitoring dashboards
2. **Configuration Management** - Use facts for dynamic configuration
3. **Automation Tools** - Integrate facts with other automation tools
4. **Reporting Systems** - Generate reports from collected facts

## Performance Considerations

### Collection Performance

- **Parallel Workers** - Use 4-8 parallel workers for optimal performance
- **Timeout Settings** - Set appropriate timeouts based on network conditions
- **Batch Processing** - Process machines in batches for large environments
- **Connection Pooling** - Use connection pooling for SSH connections

### Memory Performance

- **Memory Optimization** - Optimize memory usage for your workload
- **Efficient Allocation** - Optimize memory allocation for fact storage
- **Garbage Collection** - Proper cleanup of temporary objects
- **Memory Pooling** - Reuse memory structures for better performance

### Memory Usage

- **Streaming** - Stream large fact collections to avoid memory issues
- **Chunking** - Process facts in chunks for large datasets
- **Garbage Collection** - Proper cleanup of temporary objects
- **Monitoring** - Monitor memory usage during collection

## Security Considerations

### Data Protection

- **Memory Protection** - Facts stored in memory are not persisted to disk
- **Access Control** - Facts are only accessible during the action run
- **Audit Logging** - Logging of fact collection and usage
- **Data Sanitization** - Removal of sensitive information before collection

### Network Security

- **SSH Authentication** - Secure SSH authentication for fact collection
- **Connection Encryption** - Encrypted connections to target machines
- **Certificate Validation** - Validation of SSH certificates
- **Timeout Protection** - Protection against hanging connections

## Monitoring and Maintenance

### Health Checks

```bash
# Regular health checks
spooky facts health-check ./my-project

# Monitor collection status
spooky facts list ./my-project --check-freshness

# Validate fact integrity
spooky facts validate ./my-project
```

### Maintenance Tasks

```bash
# Regular cleanup
spooky facts cleanup ./my-project --older-than 30d

# Database optimization
spooky facts optimize ./my-project

# Backup creation
cp ./my-project/facts.db ./backup/facts.db.$(date +%Y%m%d)
```

### Monitoring Metrics

- **Collection Success Rate** - Percentage of successful fact collections
- **Collection Duration** - Time taken for fact collection
- **Storage Usage** - Disk space used by facts database
- **Validation Errors** - Number of validation failures
- **Performance Metrics** - Memory and CPU usage during collection

## Getting Help

### Documentation Resources

- **User Guide** - [FACTS_USER_GUIDE.md](FACTS_USER_GUIDE.md) for usage instructions
- **API Reference** - [FACTS_API_REFERENCE.md](FACTS_API_REFERENCE.md) for technical details
- **Troubleshooting** - [FACTS_TROUBLESHOOTING.md](FACTS_TROUBLESHOOTING.md) for problem solving
- **Examples** - [Examples Directory](examples/) for practical examples

### Support Channels

- **GitHub Issues** - Report bugs and request features
- **Documentation** - Check this documentation for solutions
- **Community** - Ask questions in the community forums
- **Code** - Review the source code for implementation details

### Contributing

- **Bug Reports** - Report issues with detailed information
- **Feature Requests** - Suggest new features and improvements
- **Code Contributions** - Submit pull requests for fixes and enhancements
- **Documentation** - Help improve the documentation

This documentation summary provides a comprehensive guide to the spooky facts system documentation, helping you find the right information for your needs and understand how all the components work together.
