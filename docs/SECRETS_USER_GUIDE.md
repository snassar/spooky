# Secrets System User Guide

## Overview

The spooky secrets system provides comprehensive encryption, decryption, and key management capabilities for securing sensitive data. This guide covers everything from basic encryption operations to advanced features like key management, audit logging, and integration with other systems.

**Status: Production Ready** - The secrets system is fully implemented with comprehensive encryption, decryption, and key management capabilities.

## Related Documentation

- [Variables User Guide](VARIABLES_USER_GUIDE.md) - Encrypted variable management
- [Actions User Guide](ACTIONS_USER_GUIDE.md) - Decrypting variables during actions
- [Templates User Guide](TEMPLATES_USER_GUIDE.md) - Using encrypted variables in templates
- [Machines User Guide](MACHINES_USER_GUIDE.md) - Machine-specific secrets

> **See also**: [User Guides Index](USER_GUIDES_INDEX.md) - Complete overview of all user guides

## Getting Started

### Prerequisites

- spooky CLI installed and configured
- Basic understanding of encryption concepts
- Access to create and modify project files
- Understanding of key management best practices

### Quick Start

1. **Check Available Secrets Commands**
   ```bash
   spooky secrets --help
   ```

2. **Encrypt Sensitive Data**
   ```bash
   spooky variables armor ./my-project --variable database-password
   ```

3. **Decrypt Data During Actions**
   ```bash
   spooky actions run ./my-project --decrypt
   ```

## Core Concepts

### Encryption Methods

spooky supports multiple encryption methods:

- **Age Encryption** (recommended) - Modern, secure encryption
- **Symmetric Encryption** - For simple use cases
- **Asymmetric Encryption** - For secure key exchange

### Key Management

The secrets system provides:

- **Key Generation** - Create encryption keys
- **Key Storage** - Secure key storage and retrieval
- **Key Rotation** - Automated key rotation
- **Access Control** - Control who can access keys

### Security Features

Security features include:

- **Audit Logging** - Track all encryption/decryption operations
- **Access Control** - Control who can perform operations
- **Key Validation** - Validate key integrity and permissions
- **Secure Storage** - Store keys securely

### Integration with Other Systems

The secrets system integrates with other spooky systems:

- **Variables**: Encrypt sensitive [variables](VARIABLES_USER_GUIDE.md)
- **Actions**: Decrypt variables during [action execution](ACTIONS_USER_GUIDE.md)
- **Templates**: Use encrypted variables in [template rendering](TEMPLATES_USER_GUIDE.md)
- **Machines**: Store machine-specific secrets securely

## Configuration

### Secrets Configuration

Configure secrets settings in your `spooky.hcl` file:

```hcl
secrets {
  encryption {
    method = "age"
    key_path = "~/.config/spooky/keys/age.key"
  }
  
  audit {
    enabled = true
    log_path = "~/.local/state/spooky/audit.log"
  }
  
  access_control {
    allowed_users = ["admin", "deploy"]
    allowed_operations = ["encrypt", "decrypt"]
  }
}
```

### Age Encryption Configuration

For age encryption (recommended):

```hcl
encryption {
  method = "age"
  key_path = "~/.config/spooky/keys/age.key"
  recipients = ["age1...", "age1..."]
}
```

### Symmetric Encryption Configuration

For symmetric encryption:

```hcl
encryption {
  method = "symmetric"
  key_path = "~/.config/spooky/keys/symmetric.key"
  algorithm = "aes-256-gcm"
}
```

## CLI Commands

### Encryption Operations

Encrypt sensitive data:

```bash
# Encrypt a variable
spooky variables armor ./my-project --variable database-password

# Encrypt multiple variables
spooky variables armor ./my-project --variable db-password --variable api-key

# Encrypt with specific key
spooky variables armor ./my-project --variable secret --key ~/.config/spooky/keys/custom.key
```

### Decryption Operations

Decrypt data during operations:

```bash
# Decrypt during action execution
spooky actions run ./my-project --decrypt

# Decrypt during dry-run
spooky actions run ./my-project --dry-run --decrypt

# Decrypt during planning
spooky actions run ./my-project --plan --decrypt
```

### Key Management

Manage encryption keys:

```bash
# List available keys
spooky secrets keys list

# Validate key integrity
spooky secrets keys validate ~/.config/spooky/keys/age.key

# Rotate encryption keys
spooky secrets keys rotate --key ~/.config/spooky/keys/age.key
```

### Audit Operations

View audit logs:

```bash
# View recent audit entries
spooky secrets audit list

# View audit entries for specific operation
spooky secrets audit list --operation encrypt

# Export audit log
spooky secrets audit export --output audit.json
```

## Advanced Features

### Key Rotation

Automated key rotation capabilities:

```bash
# Rotate keys with re-encryption
spooky secrets keys rotate --re-encrypt

# Rotate specific key
spooky secrets keys rotate --key ~/.config/spooky/keys/age.key

# Validate rotation
spooky secrets keys validate --all
```

### Access Control

Configure access control for secrets:

```hcl
access_control {
  allowed_users = ["admin", "deploy", "ci"]
  allowed_operations = ["encrypt", "decrypt", "rotate"]
  allowed_resources = ["variables", "templates", "configs"]
}
```

### Audit Logging

Comprehensive audit logging:

```hcl
audit {
  enabled = true
  log_path = "~/.local/state/spooky/audit.log"
  log_level = "info"
  retention_days = 90
  
  events = [
    "encrypt",
    "decrypt", 
    "key_rotate",
    "access_denied"
  ]
}
```

## Security Best Practices

### Key Management

- Store keys securely with proper permissions (600)
- Use different keys for different environments
- Rotate keys regularly
- Monitor key access and usage

### Encryption Practices

- Use age encryption for sensitive data
- Encrypt all sensitive variables
- Use strong, unique keys
- Validate encryption before use

### Access Control

- Implement least privilege access
- Monitor access patterns
- Log all encryption/decryption operations
- Regular access reviews

### Audit and Monitoring

- Enable comprehensive audit logging
- Monitor for unusual access patterns
- Regular audit log reviews
- Alert on security events

## Troubleshooting

### Common Encryption Issues

**Key Not Found**
```bash
# Check key path and permissions
ls -la ~/.config/spooky/keys/

# Validate key integrity
spooky secrets keys validate ~/.config/spooky/keys/age.key
```

**Decryption Failed**
```bash
# Check if decryption flag is used
spooky actions run ./my-project --decrypt

# Verify key is available
spooky secrets keys list
```

**Access Denied**
```bash
# Check access control configuration
spooky secrets config show

# Verify user permissions
spooky secrets audit list --user $(whoami)
```

### Debugging Secrets Operations

Enable verbose output for debugging:

```bash
# Verbose encryption output
spooky variables armor ./my-project --variable secret --verbose

# Debug decryption issues
spooky actions run ./my-project --decrypt --debug
```

### Performance Optimization

Optimize secrets operations:

```bash
# Use key caching
spooky secrets config set --cache-enabled true

# Monitor performance
spooky secrets audit list --metrics
```

## Integration with Other Systems

### Variables Integration

Secrets integrate with the variables system for encrypted variables:

```hcl
variables {
  variable "database-password" {
    value = "encrypted:age1..."
    encrypted = true
  }
}
```

### Actions Integration

Secrets enable secure action execution:

```bash
# Run actions with decryption
spooky actions run ./my-project --decrypt

# Use encrypted variables in actions
spooky actions run ./my-project --variable encrypted-secret --decrypt
```

### Templates Integration

Secrets enable secure template rendering:

```hcl
templates {
  template "config.tmpl" {
    source = "templates/config.tmpl"
    destination = "/etc/app/config.conf"
    
    variables {
      db_password = "encrypted:age1..."
    }
  }
}
```

## Examples

### Basic Encryption Setup

```hcl
# spooky.hcl
secrets {
  encryption {
    method = "age"
    key_path = "~/.config/spooky/keys/age.key"
  }
  
  audit {
    enabled = true
    log_path = "~/.local/state/spooky/audit.log"
  }
}
```

### Encrypted Variables

```hcl
# variables.hcl
variables {
  variable "database-password" {
    value = "encrypted:age1..."
    encrypted = true
    description = "Database password"
  }
  
  variable "api-key" {
    value = "encrypted:age1..."
    encrypted = true
    description = "API key for external service"
  }
}
```

### Secure Action Configuration

```hcl
# actions.hcl
actions {
  action "deploy-application" {
    description = "Deploy application with encrypted secrets"
    
    machines = ["web-server"]
    parallel = true
    
    variables {
      db_password = "encrypted:age1..."
      api_key = "encrypted:age1..."
    }
    
    command = "deploy.sh"
  }
}
```

### Template with Encrypted Data

```bash
# templates/config.tmpl
DATABASE_URL=postgresql://user:{{.db_password}}@localhost/db
API_KEY={{.api_key}}
```

## Best Practices

### Key Management

- Use different keys for different environments
- Store keys securely with proper permissions
- Rotate keys regularly
- Monitor key usage and access

### Encryption

- Encrypt all sensitive data
- Use strong encryption algorithms
- Validate encryption before use
- Test decryption in safe environments

### Access Control

- Implement least privilege access
- Monitor access patterns
- Regular access reviews
- Log all security events

### Audit and Compliance

- Enable comprehensive audit logging
- Monitor for security events
- Regular compliance reviews
- Maintain audit trails

## Next Steps

- Explore the [Secrets API Reference](SECRETS_API_REFERENCE.md) for detailed technical information
- Check the [Secrets Troubleshooting Guide](SECRETS_TROUBLESHOOTING.md) for common issues
- Review the [Secrets Documentation Summary](SECRETS_DOCUMENTATION_SUMMARY.md) for implementation details
- Learn about [Secrets Integration Patterns](INTEGRATIONS_USER_GUIDE.md) for advanced usage
