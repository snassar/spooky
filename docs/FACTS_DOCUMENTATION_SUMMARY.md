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
- **Storage issues?** → [Troubleshooting Guide](FACTS_TROUBLESHOOTING.md) - Memory Errors section
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
- **Troubleshooting** → [Troubleshooting Guide](FACTS_TROUBLESHOOTING.md) - Memory Issues
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
2. **FactStorage** - Provides minimal storage for debugging and statistics during export operations
3. **FactManager** - Orchestrates collection, validation, and export
4. **FactsIntegration** - Provides integration with other spooky components
5. **CLI Commands** - User interface for fact management

### Data Flow

1. **Machine Discovery** - Read machine inventory from project configuration
2. **SSH Connection** - Establish secure connection to target machine
3. **Fact Collection** - Use gopsutil to gather system information
4. **Data Processing** - Convert and validate collected data
5. **In-Memory Storage** - Provides minimal storage for debugging and statistics during export operations
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

# Export facts directly from machines
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
    storage_path = "memory"
    compression_enabled = true
  }
}
```