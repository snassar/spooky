# Spooky User Guides Index

## Overview

This index provides navigation to all spooky user guides. Each guide focuses on a specific system component while cross-referencing related functionality.

## Core System Guides

### [SSH User Guide](SSH_USER_GUIDE.md)
**Status: Fully Implemented** - SSH connectivity, authentication, and command execution
- SSH key management and authentication
- Connection testing and troubleshooting
- SSH-based operations for other systems

### [Machines User Guide](MACHINES_USER_GUIDE.md)
**Status: Production Ready** - Machine inventory and connectivity management
- Machine inventory configuration
- Connectivity testing and validation
- Machine targeting and filtering

### [Variables User Guide](VARIABLES_USER_GUIDE.md)
**Status: Production Ready** - Variable management and resolution
- Variable definition and resolution
- Environment integration
- Encrypted variable support

### [Secrets User Guide](SECRETS_USER_GUIDE.md)
**Status: Production Ready** - Encryption and key management
- Age encryption for sensitive data
- Key management and rotation
- Audit logging and access control

### [Logging User Guide](LOGGING_USER_GUIDE.md)
**Status: Production Ready** - Logging and monitoring
- Structured logging configuration
- Performance monitoring
- Log aggregation and analysis

## Operational System Guides

### [Actions User Guide](ACTIONS_USER_GUIDE.md)
**Status: Partially Implemented** - Action orchestration and execution
- Action definition and configuration
- Parallel execution and dependency management
- SSH-based action running

### [Facts User Guide](FACTS_USER_GUIDE.md)
**Status: Partially Implemented** - System fact collection and export
- Machine fact collection via SSH
- Fact export and integration
- Custom fact definition

### [Templates User Guide](TEMPLATES_USER_GUIDE.md)
**Status: Partially Implemented** - Template rendering and deployment
- Go template syntax and functions
- Variable integration in templates
- Template validation and deployment

### [Integrations User Guide](INTEGRATIONS_USER_GUIDE.md)
**Status: Partially Implemented** - External system integrations
- API integrations and webhooks
- Authentication and rate limiting
- Custom integration development

## Implementation Status Summary

### Production Ready Systems
- **SSH**: Fully implemented with comprehensive functionality
- **Machines**: Complete inventory and connectivity management
- **Variables**: Full variable management and resolution
- **Secrets**: Complete encryption and key management
- **Logging**: Comprehensive logging and monitoring

### Partially Implemented Systems
- **Actions**: Basic functionality with SSH orchestration issues
- **Facts**: Basic collection with SSH-based gathering issues
- **Templates**: Basic rendering with CLI command issues
- **Integrations**: Basic functionality with advanced features in development

> **See also**: [Known Issues](KNOWN_ISSUES.md) - Comprehensive documentation of all known issues and workarounds

## Common Workflows

### Basic Machine Management
1. [Machines User Guide](MACHINES_USER_GUIDE.md) - Configure machine inventory
2. [SSH User Guide](SSH_USER_GUIDE.md) - Set up SSH authentication
3. [Machines User Guide](MACHINES_USER_GUIDE.md) - Test connectivity

### Action Orchestration
1. [Actions User Guide](ACTIONS_USER_GUIDE.md) - Define actions
2. [Variables User Guide](VARIABLES_USER_GUIDE.md) - Configure variables
3. [Templates User Guide](TEMPLATES_USER_GUIDE.md) - Create templates
4. [Actions User Guide](ACTIONS_USER_GUIDE.md) - Run actions

### Secure Operations
1. [Secrets User Guide](SECRETS_USER_GUIDE.md) - Set up encryption
2. [Variables User Guide](VARIABLES_USER_GUIDE.md) - Encrypt sensitive variables
3. [Actions User Guide](ACTIONS_USER_GUIDE.md) - Run with decryption

### Fact Collection
1. [Machines User Guide](MACHINES_USER_GUIDE.md) - Configure target machines
2. [SSH User Guide](SSH_USER_GUIDE.md) - Ensure SSH connectivity
3. [Facts User Guide](FACTS_USER_GUIDE.md) - Collect and export facts

## CLI Command Reference

### Core Commands
- `spooky machines` - Machine inventory management
- `spooky ssh` - SSH connectivity testing
- `spooky variables` - Variable management
- `spooky secrets` - Encryption and key management

### Operational Commands
- `spooky actions` - Action orchestration
- `spooky facts` - Fact collection and export
- `spooky templates` - Template rendering
- `spooky integrations` - External integrations

### Common Flags
- `--machine <list>` - Target specific machines
- `--tags <list>` - Filter by machine tags
- `--filter <query>` - Complex filtering queries
- `--parallel <number>` - Parallel execution (minimum 2)
- `--dry-run` - Simulate without changes
- `--plan` - Show execution plan
- `--decrypt` - Decrypt variables during execution

## Troubleshooting

For troubleshooting specific systems, see the individual user guides:

- [SSH Troubleshooting](SSH_TROUBLESHOOTING.md)
- [Actions Troubleshooting](ACTIONS_TROUBLESHOOTING.md)
- [Facts Troubleshooting](FACTS_TROUBLESHOOTING.md)
- [Templates Troubleshooting](TEMPLATES_TROUBLESHOOTING.md)
- [Variables Troubleshooting](VARIABLES_TROUBLESHOOTING.md)
- [Secrets Troubleshooting](SECRETS_TROUBLESHOOTING.md)
- [Machines Troubleshooting](MACHINES_TROUBLESHOOTING.md)
- [Integrations Troubleshooting](INTEGRATIONS_TROUBLESHOOTING.md)
- [Logging Troubleshooting](LOGGING_TROUBLESHOOTING.md)

## API Reference

For detailed API documentation, see:

- [CLI Reference](CLI_REFERENCE.md)
- [Interfaces API Reference](INTERFACES_API_REFERENCE.md)
- [Actions API Reference](ACTIONS_API_REFERENCE.md)
- [Facts API Reference](FACTS_API_REFERENCE.md)
- [Machines API Reference](MACHINES_API_REFERENCE.md)
- [Variables API Reference](VARIABLES_API_REFERENCE.md)
- [Secrets API Reference](SECRETS_API_REFERENCE.md)
- [Templates API Reference](TEMPLATES_API_REFERENCE.md)
- [Integrations API Reference](INTEGRATIONS_API_REFERENCE.md)
- [Logging API Reference](LOGGING_API_REFERENCE.md)
