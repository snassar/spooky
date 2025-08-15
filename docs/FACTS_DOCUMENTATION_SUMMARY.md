# Facts System Documentation Summary

## Overview

This document provides a comprehensive overview of the spooky facts system documentation. It serves as a guide to help you find the right documentation for your needs and understand how all the pieces fit together.

**Status: Partially Implemented** - The facts system has basic functionality but SSH-based fact collection has known issues that need to be addressed.

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
- Current limitations and known issues
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

1. **Direct Export** - Facts are collected and exported immediately without intermediate storage
2. **Machine Integration** - Uses project machine inventory for target identification
3. **Multiple Export Formats** - JSON and HCL export formats
4. **Filtering Support** - Filter by machine, tags, and groups
5. **SSH-Based Collection** - Collect facts from remote machines via SSH (with known limitations)
6. **Local Collection** - Collect facts from local machine using system commands

### Architecture Principles

1. **Interface-First Design** - All functionality through well-defined interfaces
2. **Dependency Injection** - Loose coupling through interface-based dependencies
3. **Direct Export** - No intermediate storage, direct gather → export
4. **Extensible Design** - Easy to add new collection methods and export formats
5. **Performance Optimized** - Efficient collection and export with minimal overhead

### Best Practices

1. **Use Project Machine Inventory** - Leverage existing machine inventory for fact collection
2. **Filter Appropriately** - Use machine and tag filtering to limit collection scope
3. **Choose Export Format** - Use JSON for machine processing, HCL for human readability
4. **Handle SSH Issues** - Be aware of current SSH-based collection limitations
5. **Validate Before Export** - Ensure machine inventory is valid before collection
6. **Monitor Performance** - Use appropriate filtering to avoid overwhelming systems

## Implementation Status

### ✅ Completed Features

- **Core Fact Collection**
  - Local fact collection using system commands
  - Machine inventory integration for target identification
  - Direct export to JSON and HCL formats
  - Machine and tag-based filtering
  - Basic validation and error handling

- **CLI Integration**
  - `spooky facts export` command with full functionality
  - Project path validation and machine inventory loading
  - Export format selection and output file specification
  - Comprehensive error reporting and validation

- **Machine Integration**
  - Uses project machine inventory for target identification
  - Supports machine filtering by names and tags
  - Integrates with existing SSH infrastructure
  - Provides machine-specific collection context

### ⚠️ Known Issues and Limitations

- **SSH-Based Collection Issues**
  - SSH-based fact collection has known implementation issues
  - Remote `/etc/spooky/facts.*` reading is not fully functional
  - Parallel collection across multiple machines has limitations
  - SSH integration with machine inventory needs improvement

- **Current Workarounds**
  - Use local fact collection for immediate needs
  - Export facts manually from remote machines if needed
  - Monitor for updates to SSH-based collection functionality

### 🔄 Planned Improvements

- **SSH Collection Enhancement**
  - Fix SSH-based fact collection implementation
  - Improve remote facts file reading capabilities
  - Enhance parallel collection across multiple machines
  - Better integration with machine inventory authentication

- **Advanced Features**
  - Enhanced fact validation and processing
  - Custom fact collection scripts
  - Fact caching and performance optimization
  - Advanced filtering and querying capabilities

## Common Patterns

### Basic Fact Collection

```bash
# Export facts from a project
spooky facts export ./my-project --format json --output facts.json

# Export with machine filtering
spooky facts export ./my-project --machine web-server --format hcl --output web-facts.hcl

# Export with tag filtering
spooky facts export ./my-project --tags environment=production --format json --output prod-facts.json
```

### Advanced Configuration

```hcl
# Project configuration with facts integration
project {
  name = "my-project"
  description = "Project with fact collection"
  
  machines {
    machine "web-server" {
      hostname = "web.example.com"
      user = "admin"
      tags = {
        environment = "production"
        role = "web"
      }
    }
  }
}
```

## Integration with Other Systems

### CLI Integration
- [User Guide - Fact Collection](FACTS_USER_GUIDE.md#fact-collection)
- [API Reference - CLI Integration](FACTS_API_REFERENCE.md#cli-integration)
- [Examples - Export and Filtering](examples/README.md#export-and-filtering)

### Machine Integration
- [User Guide - Machine Integration](FACTS_USER_GUIDE.md#machine-integration)
- [API Reference - Machine Integration](FACTS_API_REFERENCE.md#machine-integration)
- [Examples - Machine Filtering](examples/README.md#machine-filtering)

## Recent Updates

### SSH Collection Issues (Latest)

The facts system has identified issues with SSH-based fact collection:

#### ⚠️ **Known Issues**
- **SSH-Based Collection**: SSH-based fact collection has implementation issues
- **Remote Facts Reading**: Cannot reliably read `/etc/spooky/facts.*` files from remote machines
- **Parallel Processing**: Sequential collection only, no multi-machine parallel processing
- **SSH Integration**: Cannot fully leverage existing SSH infrastructure and machine inventory
- **Authentication**: Limited support for machine inventory authentication methods

#### 🔧 **Current Workarounds**
- **Local Collection**: Use local fact collection for immediate needs
- **Manual Export**: Export facts manually from remote machines if needed
- **Filtered Collection**: Use machine filtering to limit collection scope
- **Monitor Updates**: Watch for improvements to SSH-based collection

#### 📚 **Updated Documentation**
- **User Guide**: Updated to reflect current SSH collection limitations
- **API Reference**: Added implementation status indicators for SSH-based features
- **Troubleshooting Guide**: Added SSH collection issues section
- **Integration Guides**: Updated to reflect current SSH integration status

## Remember

**Good facts system usage enables:**
- Efficient system information collection
- Machine inventory integration
- Flexible export and filtering
- Integration with other spooky systems
- Performance-optimized data gathering

**Always be aware of current SSH-based collection limitations and use appropriate workarounds until these issues are resolved.**