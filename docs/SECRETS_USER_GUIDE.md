# Age Encryption User Guide

## Overview

The spooky age encryption system provides comprehensive secrets management using [age encryption](https://github.com/FiloSottile/age) for secure storage and transmission of sensitive data. This guide covers everything from basic encryption setup to advanced workflows for variables, facts, and machine authentication.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Basic Concepts](#basic-concepts)
3. [Configuration Setup](#configuration-setup)
4. [Variable Encryption](#variable-encryption)
5. [Machine Authentication Encryption](#machine-authentication-encryption)
6. [Facts Decryption](#facts-decryption)
7. [CLI Commands](#cli-commands)
8. [Advanced Usage](#advanced-usage)
9. [Best Practices](#best-practices)
10. [Troubleshooting](#troubleshooting)

## Getting Started

### Prerequisites

1. **Install age CLI tools** (spooky does not generate keys):
   ```bash
   # On Ubuntu/Debian
   sudo apt install age
   
   # On macOS
   brew install age
   
   # On other systems, download from https://github.com/FiloSottile/age/releases
   ```

2. **Generate age keys**:
   ```bash
   # Create spooky config directory
   mkdir -p ~/.config/spooky/identities
   
   # Generate identity file
   age-keygen -o ~/.config/spooky/identities/identity.txt
   
   # Extract public key for recipients file
   age-keygen -y ~/.config/spooky/identities/identity.txt > ~/.config/spooky/recipients.txt
   ```

3. **Set proper permissions**:
   ```bash
   chmod 600 ~/.config/spooky/identities/identity.txt
   chmod 644 ~/.config/spooky/recipients.txt
   ```

### Quick Start

1. **Initialize a project**:
   ```bash
   spooky project init ./my-project
   ```

2. **Add encrypted variables**:
   ```hcl
   # variables.hcl
   variables {
     variable "database_password" {
       type = "string"
       description = "Database password"
       default = "my-secret-password"
       encrypted = true
     }
   }
   ```

3. **Encrypt the variables**:
   ```bash
   spooky variables encrypt ./my-project
   ```

4. **Run actions with decryption**:
   ```bash
   spooky actions run ./my-project --decrypt
   ```

## Basic Concepts

### Age Encryption

Age is a modern encryption tool that provides:
- **Asymmetric encryption** using X25519 keys
- **Multiple recipients** support
- **Armored output** for easy transmission
- **Passphrase encryption** as an alternative

### Encryption Scenarios

1. **Variables**: Encrypt sensitive configuration values
2. **Machine Authentication**: Encrypt SSH passphrases and passwords
3. **Facts**: Decrypt age-encrypted facts from target machines

### Key Management

- **Identity files**: Private keys for decryption (stored securely)
- **Recipients files**: Public keys for encryption (can be shared)
- **No key generation**: spooky does not generate keys - use age CLI tools

## Configuration Setup

### Global Configuration

Configure age encryption in `~/.config/spooky/spooky.hcl`:

```hcl
age {
  # Identity file path (for decryption)
  identities = "~/.config/spooky/identities"
  
  # Recipients file path (for encryption)
  recipients = "~/.config/spooky/recipients.txt"
  
  # Security settings
  validation {
    strict_mode = true
    check_recipients = true
    validate_keys = true
  }
  
  # Encryption settings
  encryption {
    algorithm = "age"
    compression = false
    armor = true  # Use armored output for easy transmission
  }
}
```

### Project-Level Recipients

Optionally add project-specific recipients in `./my-project/recipients.txt`:

```txt
# Project-specific recipients
age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p
age1abc123def456ghi789jkl012mno345pqr678stu901vwx234yz
```

## Variable Encryption

### Basic Variable Encryption

Encrypt individual variables:

```hcl
# variables.hcl
variables {
  variable "api_key" {
    type = "string"
    description = "API key for external service"
    default = "sk-1234567890abcdef"
    encrypted = true
  }
  
  variable "webhook_secret" {
    type = "string"
    description = "Webhook secret for verification"
    default = "whsec_abcdef123456"
    encrypted = true
  }
}
```

### Object Encryption

Encrypt entire objects or maps:

```hcl
# variables.hcl
variables {
  variable "database_config" {
    type = "object"
    description = "Database configuration with credentials"
    default = {
      host = "db.example.com"
      port = 5432
      username = "app_user"
      password = "secret_password"
      ssl_mode = "require"
    }
    encrypted = true
  }
  
  variable "api_config" {
    type = "map"
    description = "API configuration with secrets"
    default = {
      base_url = "https://api.example.com"
      api_key = "sk-1234567890abcdef"
      webhook_secret = "whsec_abcdef123456"
      timeout = 30
    }
    encrypted = true
  }
}
```

### Mixed Content

Mix encrypted and plaintext variables:

```hcl
# variables.hcl
variables {
  # Plaintext variables (not encrypted)
  variable "app_name" {
    type = "string"
    description = "Application name"
    default = "my-app"
  }
  
  variable "environment" {
    type = "string"
    description = "Environment name"
    default = "production"
  }
  
  # Encrypted variables
  variable "database_password" {
    type = "string"
    description = "Database password"
    default = "secret_password"
    encrypted = true
  }
  
  variable "secrets" {
    type = "object"
    description = "Application secrets"
    default = {
      jwt_secret = "super-secret-jwt-key"
      api_key = "sk-1234567890abcdef"
    }
    encrypted = true
  }
}
```

### Encryption Commands

Encrypt variables in your project:

```bash
# Show what would be encrypted (dry run)
spooky variables encrypt ./my-project --dry-run

# Actually encrypt the variables
spooky variables encrypt ./my-project
```

## Machine Authentication Encryption

### SSH Key Passphrase Encryption

Encrypt SSH key passphrases:

```hcl
# machines.hcl
machines {
  machine "web-server" {
    hostname = "web.example.com"
    port = 22
    user = "admin"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/id_rsa"
      passphrase = {
        value = "my-ssh-passphrase"
        encrypted = true
      }
    }
  }
}
```

### Password Encryption

Encrypt SSH passwords:

```hcl
# machines.hcl
machines {
  machine "db-server" {
    hostname = "db.example.com"
    port = 22
    user = "postgres"
    
    authentication {
      method = "password"
      password = {
        value = "database-password"
        encrypted = true
      }
    }
  }
}
```

### Mixed Authentication

Mix encrypted and plaintext authentication:

```hcl
# machines.hcl
machines {
  # Machine with encrypted authentication
  machine "secure-server" {
    hostname = "secure.example.com"
    port = 22
    user = "admin"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/id_rsa"
      passphrase = {
        value = "encrypted-passphrase"
        encrypted = true
      }
    }
  }
  
  # Machine with plaintext authentication
  machine "dev-server" {
    hostname = "dev.example.com"
    port = 22
    user = "developer"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/id_ed25519"
      # No passphrase - uses SSH agent
    }
  }
}
```

### Encryption Commands

Encrypt machine authentication:

```bash
# Show what would be encrypted (dry run)
spooky machines encrypt ./my-project --dry-run

# Actually encrypt the machines
spooky machines encrypt ./my-project
```

## Facts Decryption

### Encrypted Facts on Target Machines

Target machines can have age-encrypted facts in `/etc/spooky/custom.hcl`:

```hcl
# /etc/spooky/custom.hcl on target machine
custom {
  # Plaintext facts
  environment = "production"
  datacenter = "us-west-1"
  
  # Encrypted facts (automatically detected by age1 prefix)
  database_connection_string = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
  
  # Encrypted objects
  application_secrets = {
    api_key = "age1abc123def456ghi789jkl012mno345pqr678stu901vwx234yz"
    webhook_secret = "age1xyz789abc123def456ghi789jkl012mno345pqr678stu901vwx"
  }
}
```

### Automatic Decryption

Facts are automatically decrypted when using the `--decrypt` flag:

```bash
# Decrypt facts during fact collection
spooky facts gather ./my-project --decrypt

# Decrypt facts during actions
spooky actions run ./my-project --decrypt
```

## CLI Commands

### Project-Wide Encryption

Encrypt all variables and machines in a project:

```bash
# Show what would be encrypted (dry run)
spooky project encrypt ./my-project --dry-run

# Actually encrypt everything
spooky project encrypt ./my-project
```

### Variable Encryption

Encrypt only variables:

```bash
# Show what would be encrypted (dry run)
spooky variables encrypt ./my-project --dry-run

# Actually encrypt variables
spooky variables encrypt ./my-project
```

### Machine Encryption

Encrypt only machines:

```bash
# Show what would be encrypted (dry run)
spooky machines encrypt ./my-project --dry-run

# Actually encrypt machines
spooky machines encrypt ./my-project
```

### Decryption During Actions

Decrypt encrypted data during action execution:

```bash
# Run actions with decryption enabled
spooky actions run ./my-project --decrypt

# Combine with other flags
spooky actions run ./my-project --decrypt --dry-run --plan
```

### Secrets Validation

Validate age configuration and keys:

```bash
# Validate age configuration for a project
spooky secrets validate ./my-project
```

## Advanced Usage

### Multiple Recipients

Encrypt to multiple recipients for team access:

```bash
# Add multiple recipients to recipients.txt
cat >> ~/.config/spooky/recipients.txt << EOF
# Team member 1
age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p

# Team member 2
age1abc123def456ghi789jkl012mno345pqr678stu901vwx234yz

# CI/CD system
age1xyz789abc123def456ghi789jkl012mno345pqr678stu901vwx
EOF
```

### Key Rotation

Rotate encryption keys:

```bash
# 1. Add new recipient to recipients.txt
echo "age1newkey1234567890abcdefghijklmnopqrstuvwxyz" >> ~/.config/spooky/recipients.txt

# 2. Re-encrypt all data with new recipients
spooky project encrypt ./my-project

# 3. Remove old recipient from recipients.txt
# (Edit recipients.txt to remove old keys)

# 4. Re-encrypt again to remove old recipient access
spooky project encrypt ./my-project
```

### Passphrase Encryption

Use passphrase encryption as an alternative:

```hcl
# variables.hcl with passphrase encryption
variables {
  variable "backup_key" {
    type = "string"
    description = "Backup encryption key"
    default = "my-backup-key"
    encrypted = true
    # Uses passphrase from spooky.hcl age.passphrase
  }
}
```

### Conditional Encryption

Use environment-specific encryption:

```hcl
# variables.hcl with conditional encryption
variables {
  variable "api_key" {
    type = "string"
    description = "API key"
    default = "sk-1234567890abcdef"
    encrypted = true  # Only encrypt in production
  }
}

# In development, you might not encrypt
# In production, always encrypt sensitive data
```

## Best Practices

### Key Management

1. **Secure storage**: Store identity files with 600 permissions
2. **Backup keys**: Regularly backup identity files
3. **Key rotation**: Rotate keys periodically
4. **Access control**: Limit access to identity files

### Encryption Strategy

1. **Encrypt sensitive data**: Always encrypt passwords, keys, tokens
2. **Object encryption**: Encrypt entire objects rather than individual fields
3. **Mixed content**: Use plaintext for non-sensitive data
4. **Environment separation**: Use different keys for different environments

### Security Considerations

1. **No automatic decryption**: Decryption only happens with explicit `--decrypt` flag
2. **Memory management**: Decrypted values are cleared from memory after use
3. **Logging protection**: Decrypted values are never logged
4. **Audit trail**: All encryption/decryption operations are logged

### Performance Optimization

1. **Batch operations**: Use project-wide encryption for efficiency
2. **Dry runs**: Always use `--dry-run` first to see what will be encrypted
3. **Selective encryption**: Only encrypt what's necessary
4. **Caching**: Recipients are cached for performance

## Troubleshooting

### Common Issues

1. **Missing identity file**:
   ```bash
   Error: identity file not found: ~/.config/spooky/identities/identity.txt
   Solution: Generate identity file with age-keygen
   ```

2. **Incorrect permissions**:
   ```bash
   Error: identity file has incorrect permissions: ~/.config/spooky/identities/identity.txt
   Solution: Set permissions to 600
   ```

3. **Invalid recipients**:
   ```bash
   Error: invalid recipient in recipients.txt
   Solution: Check recipient format and regenerate if needed
   ```

4. **Decryption failures**:
   ```bash
   Error: failed to decrypt variable 'database_password'
   Solution: Check identity file and encrypted value format
   ```

### Debugging Commands

```bash
# Validate age configuration
spooky secrets validate ./my-project

# Check identity file
age-keygen -y ~/.config/spooky/identities/identity.txt

# Test encryption/decryption
echo "test" | age -r age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p | age -d -i ~/.config/spooky/identities/identity.txt
```

### Getting Help

1. **Check logs**: Look for detailed error messages
2. **Validate configuration**: Use `spooky secrets validate`
3. **Test manually**: Use age CLI tools to test encryption/decryption
4. **Review documentation**: Check this guide and API reference

## Conclusion

The spooky age encryption system provides secure, flexible secrets management for your automation projects. By following this guide, you can effectively encrypt sensitive data while maintaining security and usability.

For more advanced usage and troubleshooting, refer to the [API Reference](SECRETS_API_REFERENCE.md) and [Troubleshooting Guide](SECRETS_TROUBLESHOOTING.md).
