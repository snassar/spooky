# Age Encryption Troubleshooting Guide

## Overview

This guide helps you diagnose and resolve common issues with the spooky age encryption system. It covers problems from basic setup to advanced encryption/decryption scenarios.

## Table of Contents

1. [Quick Diagnosis](#quick-diagnosis)
2. [Setup and Configuration Issues](#setup-and-configuration-issues)
3. [Key Management Problems](#key-management-problems)
4. [Encryption Issues](#encryption-issues)
5. [Decryption Issues](#decryption-issues)
6. [CLI Command Problems](#cli-command-problems)
7. [Performance Issues](#performance-issues)
8. [Integration Problems](#integration-problems)
9. [Debugging Tools](#debugging-tools)
10. [Getting Help](#getting-help)

## Quick Diagnosis

### Diagnostic Commands

Run these commands to quickly identify common issues:

```bash
# 1. Check spooky version and basic functionality
spooky --version

# 2. Validate age configuration
spooky secrets validate ./my-project

# 3. Test age CLI tools
age-keygen --version

# 4. Check identity file permissions
ls -la ~/.config/spooky/identities/

# 5. Verify recipients file
cat ~/.config/spooky/recipients.txt

# 6. Test basic encryption/decryption
echo "test" | age -r $(head -1 ~/.config/spooky/recipients.txt) | age -d -i ~/.config/spooky/identities/identity.txt
```

### Common Error Patterns

| Error Pattern | Likely Cause | Quick Fix |
|---------------|--------------|-----------|
| `identity file not found` | Missing identity file | Generate with `age-keygen` |
| `invalid recipient` | Malformed public key | Check recipients.txt format |
| `decryption failed` | Wrong identity or corrupted data | Verify identity file |
| `permission denied` | Incorrect file permissions | Set 600 for identity files |
| `no recipients found` | Empty recipients file | Add recipients to file |

## Setup and Configuration Issues

### Problem: Missing Age CLI Tools

**Symptoms:**
```bash
Error: age-keygen: command not found
```

**Solution:**
```bash
# Ubuntu/Debian
sudo apt update && sudo apt install age

# macOS
brew install age

# Other systems
# Download from https://github.com/FiloSottile/age/releases
```

### Problem: Missing Configuration Directory

**Symptoms:**
```bash
Error: configuration directory not found: ~/.config/spooky
```

**Solution:**
```bash
# Create configuration directory
mkdir -p ~/.config/spooky/identities

# Set proper permissions
chmod 700 ~/.config/spooky
chmod 700 ~/.config/spooky/identities
```

### Problem: Invalid Configuration File

**Symptoms:**
```bash
Error: failed to parse spooky.hcl: invalid HCL syntax
```

**Solution:**
```bash
# Validate HCL syntax
hclfmt ~/.config/spooky/spooky.hcl

# Check for common syntax errors:
# - Missing quotes around strings
# - Incorrect block syntax
# - Invalid attribute names
```

**Example of correct configuration:**
```hcl
age {
  identities = "~/.config/spooky/identities"
  recipients = "~/.config/spooky/recipients.txt"
  
  validation {
    strict_mode = true
    check_recipients = true
    validate_keys = true
  }
  
  encryption {
    algorithm = "age"
    armor = true
  }
}
```

### Problem: Configuration Not Found

**Symptoms:**
```bash
Error: configuration file not found: ~/.config/spooky/spooky.hcl
```

**Solution:**
```bash
# Create default configuration
cat > ~/.config/spooky/spooky.hcl << 'EOF'
age {
  identities = "~/.config/spooky/identities"
  recipients = "~/.config/spooky/recipients.txt"
  
  validation {
    strict_mode = true
    check_recipients = true
    validate_keys = true
  }
  
  encryption {
    algorithm = "age"
    armor = true
  }
}
EOF
```

## Key Management Problems

### Problem: Missing Identity File

**Symptoms:**
```bash
Error: identity file not found: ~/.config/spooky/identities/identity.txt
```

**Solution:**
```bash
# Generate new identity file
age-keygen -o ~/.config/spooky/identities/identity.txt

# Set proper permissions
chmod 600 ~/.config/spooky/identities/identity.txt

# Extract public key for recipients file
age-keygen -y ~/.config/spooky/identities/identity.txt > ~/.config/spooky/recipients.txt
```

### Problem: Incorrect File Permissions

**Symptoms:**
```bash
Error: identity file has incorrect permissions: ~/.config/spooky/identities/identity.txt
```

**Solution:**
```bash
# Fix identity file permissions
chmod 600 ~/.config/spooky/identities/identity.txt

# Fix directory permissions
chmod 700 ~/.config/spooky/identities

# Verify permissions
ls -la ~/.config/spooky/identities/
```

### Problem: Invalid Recipients File

**Symptoms:**
```bash
Error: invalid recipient in recipients.txt: age1invalidkey
```

**Solution:**
```bash
# Check recipients file format
cat ~/.config/spooky/recipients.txt

# Each line should contain a valid age public key
# Format: age1[base32-encoded-public-key]

# Regenerate recipients file from identity
age-keygen -y ~/.config/spooky/identities/identity.txt > ~/.config/spooky/recipients.txt

# Or manually add valid recipients
echo "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p" >> ~/.config/spooky/recipients.txt
```

### Problem: Empty Recipients File

**Symptoms:**
```bash
Error: no recipients found in recipients file
```

**Solution:**
```bash
# Add recipients to file
echo "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p" > ~/.config/spooky/recipients.txt

# Or regenerate from identity
age-keygen -y ~/.config/spooky/identities/identity.txt > ~/.config/spooky/recipients.txt
```

## Encryption Issues

### Problem: Encryption Fails with "Invalid Recipient"

**Symptoms:**
```bash
Error: failed to encrypt variable 'database_password': invalid recipient
```

**Diagnosis:**
```bash
# Check recipient format
cat ~/.config/spooky/recipients.txt

# Test recipient with age CLI
echo "test" | age -r age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p
```

**Solution:**
```bash
# Regenerate valid recipients
age-keygen -y ~/.config/spooky/identities/identity.txt > ~/.config/spooky/recipients.txt

# Or manually fix recipient format
# Remove any comments, extra spaces, or invalid characters
```

### Problem: Encryption Fails with "No Recipients"

**Symptoms:**
```bash
Error: failed to encrypt variable 'api_key': no recipients provided
```

**Solution:**
```bash
# Check if recipients file exists and has content
ls -la ~/.config/spooky/recipients.txt
cat ~/.config/spooky/recipients.txt

# If empty or missing, regenerate
age-keygen -y ~/.config/spooky/identities/identity.txt > ~/.config/spooky/recipients.txt
```

### Problem: Object Encryption Fails

**Symptoms:**
```bash
Error: failed to encrypt object 'config': JSON serialization failed
```

**Diagnosis:**
```bash
# Check if object contains unsupported types
# Age encryption only supports basic JSON types
```

**Solution:**
```hcl
# Ensure object contains only supported types:
# - strings
# - numbers
# - booleans
# - arrays
# - objects (nested)

# Example of supported object:
variable "config" {
  type = "object"
  default = {
    host = "db.example.com"
    port = 5432
    ssl = true
    options = ["require", "verify-ca"]
  }
  encrypted = true
}
```

### Problem: Dry Run Shows No Changes

**Symptoms:**
```bash
$ spooky variables encrypt ./my-project --dry-run
No variables to encrypt
```

**Diagnosis:**
```bash
# Check if variables have encrypted=true
grep -r "encrypted = true" ./my-project/variables.hcl
```

**Solution:**
```hcl
# Add encrypted=true to variables that should be encrypted
variable "database_password" {
  type = "string"
  default = "secret-password"
  encrypted = true  # This was missing
}
```

## Decryption Issues

### Problem: Decryption Fails with "Wrong Identity"

**Symptoms:**
```bash
Error: failed to decrypt variable 'database_password': decryption failed
```

**Diagnosis:**
```bash
# Check if identity file matches the one used for encryption
age-keygen -y ~/.config/spooky/identities/identity.txt

# Compare with the public key that was used for encryption
```

**Solution:**
```bash
# Use the correct identity file
# If you have multiple identity files, specify the correct one:

# Option 1: Update configuration
# Edit ~/.config/spooky/spooky.hcl to point to correct identity

# Option 2: Use different identity file
spooky actions run ./my-project --decrypt --identity ~/.ssh/age-identity.txt
```

### Problem: Decryption Fails with "Invalid Format"

**Symptoms:**
```bash
Error: failed to decrypt variable 'api_key': invalid age format
```

**Diagnosis:**
```bash
# Check if the encrypted value is actually age-encrypted
# Age-encrypted values start with "age1" or "-----BEGIN AGE ENCRYPTED FILE-----"
```

**Solution:**
```bash
# The value might not be age-encrypted
# Check the original value in variables.hcl

# If it should be encrypted, run encryption first:
spooky variables encrypt ./my-project
```

### Problem: Object Decryption Fails

**Symptoms:**
```bash
Error: failed to decrypt object 'config': JSON deserialization failed
```

**Diagnosis:**
```bash
# The encrypted object might be corrupted or not properly encrypted
```

**Solution:**
```bash
# Re-encrypt the object
spooky variables encrypt ./my-project

# Or check if the object was encrypted as a whole (not individual fields)
```

### Problem: Facts Decryption Fails

**Symptoms:**
```bash
Error: failed to decrypt facts: no identity file found
```

**Solution:**
```bash
# Ensure identity file exists and is accessible
ls -la ~/.config/spooky/identities/identity.txt

# Check permissions
chmod 600 ~/.config/spooky/identities/identity.txt

# Verify identity file is valid
age-keygen -y ~/.config/spooky/identities/identity.txt
```

## CLI Command Problems

### Problem: "spooky project encrypt" Fails

**Symptoms:**
```bash
Error: failed to encrypt project: no variables or machines to encrypt
```

**Diagnosis:**
```bash
# Check if any variables or machines have encrypted=true
grep -r "encrypted = true" ./my-project/
```

**Solution:**
```hcl
# Add encrypted=true to variables or machines that should be encrypted

# In variables.hcl:
variable "secret" {
  type = "string"
  default = "value"
  encrypted = true
}

# In machines.hcl:
machine "server" {
  hostname = "example.com"
  authentication {
    method = "password"
    password = {
      value = "secret"
      encrypted = true
    }
  }
}
```

### Problem: "spooky secrets validate" Fails

**Symptoms:**
```bash
Error: age configuration validation failed
```

**Diagnosis:**
```bash
# Run validation with verbose output
spooky secrets validate ./my-project --verbose
```

**Solution:**
```bash
# Fix common validation issues:

# 1. Check identity file exists and is valid
ls -la ~/.config/spooky/identities/identity.txt
age-keygen -y ~/.config/spooky/identities/identity.txt

# 2. Check recipients file exists and has valid keys
cat ~/.config/spooky/recipients.txt

# 3. Verify configuration syntax
hclfmt ~/.config/spooky/spooky.hcl
```

### Problem: "--decrypt" Flag Not Working

**Symptoms:**
```bash
Error: decryption failed: identity file not found
```

**Solution:**
```bash
# Ensure identity file exists
ls -la ~/.config/spooky/identities/identity.txt

# Check configuration points to correct identity file
cat ~/.config/spooky/spooky.hcl

# Verify identity file permissions
chmod 600 ~/.config/spooky/identities/identity.txt
```

## Performance Issues

### Problem: Encryption is Slow

**Symptoms:**
```bash
# Encryption takes a long time for large objects
```

**Diagnosis:**
```bash
# Check object size and complexity
# Large objects or deeply nested structures can be slow
```

**Solution:**
```bash
# Optimize object structure
# Break large objects into smaller ones
# Avoid deeply nested structures

# Use batch operations
spooky project encrypt ./my-project  # Instead of individual commands
```

### Problem: Memory Usage is High

**Symptoms:**
```bash
# High memory usage during encryption/decryption
```

**Solution:**
```bash
# Process smaller batches
# Use dry-run to check what will be encrypted first
spooky variables encrypt ./my-project --dry-run
```

## Integration Problems

### Problem: Variables Integration Fails

**Symptoms:**
```bash
Error: variables integration failed: encryption error
```

**Diagnosis:**
```bash
# Check if secrets integration is properly configured
spooky secrets validate ./my-project
```

**Solution:**
```bash
# Ensure secrets integration is working
# Check age configuration
cat ~/.config/spooky/spooky.hcl

# Test basic encryption
echo "test" | age -r $(head -1 ~/.config/spooky/recipients.txt)
```

### Problem: Machines Integration Fails

**Symptoms:**
```bash
Error: machines integration failed: decryption error
```

**Solution:**
```bash
# Check machine authentication encryption
# Ensure identity file can decrypt machine passwords/passphrases

# Test decryption manually
echo "encrypted-value" | age -d -i ~/.config/spooky/identities/identity.txt
```

### Problem: Facts Integration Fails

**Symptoms:**
```bash
Error: facts integration failed: decryption error
```

**Solution:**
```bash
# Check if facts on target machines are properly encrypted
# Verify identity file can decrypt the facts

# Test fact decryption
spooky facts gather ./my-project --decrypt --dry-run
```

## Debugging Tools

### Manual Age Testing

```bash
# Test basic encryption/decryption
echo "secret data" | age -r age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p | age -d -i ~/.config/spooky/identities/identity.txt

# Test with armored output
echo "secret data" | age -r age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p -a | age -d -i ~/.config/spooky/identities/identity.txt

# List recipients from encrypted data
age -d -i ~/.config/spooky/identities/identity.txt encrypted_file.age
```

### Configuration Validation

```bash
# Validate HCL syntax
hclfmt ~/.config/spooky/spooky.hcl

# Check configuration structure
spooky config validate ~/.config/spooky/spooky.hcl

# Validate project configuration
spooky validate --config ~/.config/spooky/spooky.hcl ./my-project
```

### Logging and Debugging

```bash
# Enable debug logging
export SPOOKY_LOG_LEVEL=debug

# Run commands with verbose output
spooky variables encrypt ./my-project --verbose

# Check spooky logs
tail -f ~/.local/share/spooky/logs/spooky.log
```

### Key Management Tools

```bash
# Generate new identity
age-keygen -o ~/.config/spooky/identities/new-identity.txt

# Extract public key
age-keygen -y ~/.config/spooky/identities/identity.txt

# Validate identity file
age-keygen -y ~/.config/spooky/identities/identity.txt > /dev/null && echo "Valid" || echo "Invalid"

# List all identities in directory
for f in ~/.config/spooky/identities/*; do
  echo "File: $f"
  age-keygen -y "$f"
done
```

## Getting Help

### Information to Collect

When seeking help, collect this information:

1. **Error messages**: Complete error output
2. **Configuration**: Contents of `~/.config/spooky/spooky.hcl`
3. **Identity file**: Output of `age-keygen -y ~/.config/spooky/identities/identity.txt`
4. **Recipients file**: Contents of `~/.config/spooky/recipients.txt`
5. **File permissions**: Output of `ls -la ~/.config/spooky/identities/`
6. **Age version**: Output of `age --version`
7. **Spooky version**: Output of `spooky --version`
8. **Operating system**: `uname -a`

### Debugging Checklist

Before seeking help, verify:

- [ ] Age CLI tools are installed and working
- [ ] Identity file exists and has correct permissions (600)
- [ ] Recipients file exists and contains valid public keys
- [ ] Configuration file syntax is valid
- [ ] Variables/machines have `encrypted = true`
- [ ] Identity file can decrypt test data
- [ ] No typos in file paths or configuration

### Common Solutions Summary

| Problem | Quick Fix |
|---------|-----------|
| Missing identity file | `age-keygen -o ~/.config/spooky/identities/identity.txt` |
| Wrong permissions | `chmod 600 ~/.config/spooky/identities/identity.txt` |
| Invalid recipients | `age-keygen -y ~/.config/spooky/identities/identity.txt > ~/.config/spooky/recipients.txt` |
| Configuration errors | `hclfmt ~/.config/spooky/spooky.hcl` |
| No encrypted variables | Add `encrypted = true` to variables |
| Decryption failures | Check identity file matches encryption key |

### Additional Resources

- [Age Documentation](https://github.com/FiloSottile/age)
- [Spooky User Guide](SECRETS_USER_GUIDE.md)
- [Spooky API Reference](SECRETS_API_REFERENCE.md)
- [Age Key Management](https://github.com/FiloSottile/age#usage)

## Conclusion

Most age encryption issues can be resolved by following this troubleshooting guide. The key is to systematically check each component: CLI tools, configuration, keys, and permissions.

If you continue to experience issues after following this guide, collect the debugging information listed above and seek help from the spooky community or maintainers.
