# Secrets Management Troubleshooting Guide

## Overview

This troubleshooting guide provides solutions for common issues encountered when working with the spooky secrets management system. It covers error messages, encryption problems, key management issues, and performance problems.

**Status: Production Ready** - The secrets system is fully implemented with comprehensive encryption, key management, and integration capabilities.

## Secrets System Status

### ✅ Fully Functional Secrets Infrastructure

The secrets system now has **complete secrets infrastructure** with:

- **Age Encryption**: Full age encryption and decryption support
- **Key Management**: Comprehensive key generation and management
- **HCL Processing**: Complete HCL file encryption and decryption
- **CLI Integration**: Full CLI integration with `spooky secrets` commands
- **Project Integration**: Secrets management from project configuration
- **Variable Integration**: Variable encryption and decryption capabilities
- **Export Functionality**: Secrets export and management
- **Error Handling**: Comprehensive error handling and reporting

### What This Means for Users

- **No More Stubs**: All functionality is fully implemented - no placeholder code
- **Production Ready**: The system is ready for production use
- **Complete Feature Set**: All documented features are functional
- **Reliable Encryption**: Robust age-based encryption and decryption
- **Performance Optimized**: Efficient secrets management and processing

### Expected Behavior

When using secrets, you can expect:

1. **Proper Key Management**: Age key generation and management
2. **Encryption**: Reliable encryption of sensitive data
3. **Decryption**: Secure decryption with proper key validation
4. **HCL Processing**: Encryption and decryption of HCL files
5. **Variable Integration**: Variable encryption and decryption
6. **Error Handling**: Clear error messages with actionable information

## Common Issues and Solutions

### Key Management Errors

#### "Failed to generate age key: age-keygen not found"

**Cause:** Age CLI tools are not installed on the system.

**Solution:**
```bash
# Install age CLI tools
# Ubuntu/Debian
sudo apt update && sudo apt install age

# macOS
brew install age

# CentOS/RHEL
sudo yum install age

# Verify installation
age-keygen --version
```

#### "Failed to read identity file: permission denied"

**Cause:** Identity file has incorrect permissions.

**Solution:**
```bash
# Check current permissions
ls -la ~/.config/spooky/identities/

# Fix permissions
chmod 600 ~/.config/spooky/identities/identity.txt
chmod 700 ~/.config/spooky/identities/

# Verify permissions
ls -la ~/.config/spooky/identities/
```

#### "Identity file not found"

**Cause:** Age identity file doesn't exist.

**Solution:**
```bash
# Generate new identity
age-keygen -o ~/.config/spooky/identities/identity.txt

# Set correct permissions
chmod 600 ~/.config/spooky/identities/identity.txt

# Verify identity file
ls -la ~/.config/spooky/identities/identity.txt
```

### Encryption Errors

#### "Failed to encrypt data: no recipients specified"

**Cause:** No encryption recipients are configured.

**Solution:**
```bash
# Add recipients to configuration
echo "age1..." >> ~/.config/spooky/recipients.txt

# Or use environment variable
export SPOOKY_RECIPIENTS="age1..."

# Or specify recipients in command
spooky secrets encrypt --recipients age1... data.txt
```

#### "Failed to encrypt: invalid recipient format"

**Cause:** Recipient key format is invalid.

**Solution:**
```bash
# Check recipient format
cat ~/.config/spooky/recipients.txt

# Valid age public key format: age1...
# Example: age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p

# Fix invalid recipients
# Remove invalid lines and add correct ones
echo "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p" > ~/.config/spooky/recipients.txt
```

#### "Encryption failed: file not found"

**Cause:** File to encrypt doesn't exist.

**Solution:**
```bash
# Check if file exists
ls -la data.txt

# Create test file if needed
echo "sensitive data" > data.txt

# Try encryption again
spooky secrets encrypt data.txt
```

### Decryption Errors

#### "Failed to decrypt: no identity found"

**Cause:** No age identity file is available for decryption.

**Solution:**
```bash
# Check if identity exists
ls -la ~/.config/spooky/identities/

# Generate identity if missing
age-keygen -o ~/.config/spooky/identities/identity.txt

# Set permissions
chmod 600 ~/.config/spooky/identities/identity.txt
```

#### "Decryption failed: wrong identity"

**Cause:** Using wrong identity file for decryption.

**Solution:**
```bash
# Check which identity was used for encryption
# Age files contain recipient information

# Use correct identity
spooky secrets decrypt --identity ~/.config/spooky/identities/correct_identity.txt data.txt.age

# Or set identity in environment
export SPOOKY_IDENTITY_PATH=~/.config/spooky/identities/correct_identity.txt
spooky secrets decrypt data.txt.age
```

#### "Decryption failed: corrupted data"

**Cause:** Encrypted data is corrupted or incomplete.

**Solution:**
```bash
# Check if file is complete
ls -la data.txt.age

# Verify file integrity
file data.txt.age

# Try manual decryption to test
age -d -i ~/.config/spooky/identities/identity.txt data.txt.age
```

### HCL Processing Errors

#### "Failed to process HCL file: invalid syntax"

**Cause:** HCL file has syntax errors.

**Solution:**
```bash
# Validate HCL syntax
spooky secrets validate-hcl config.hcl

# Check for common HCL errors:
# - Missing quotes around strings
# - Incorrect block structure
# - Invalid attribute names

# Fix HCL syntax
# Example of correct syntax:
variables {
  variable "secret_key" {
    type = "string"
    default = "encrypted_value"
  }
}
```

#### "Failed to encrypt HCL file: no sensitive values found"

**Cause:** HCL file doesn't contain values marked for encryption.

**Solution:**
```hcl
# Mark values for encryption
variables {
  variable "database_password" {
    type = "string"
    default = "secret123"  # This will be encrypted
    sensitive = true       # Mark as sensitive
  }
}
```

### CLI Command Errors

#### "Command not found: spooky secrets"

**Cause:** Secrets command is not available.

**Solution:**
```bash
# Check if spooky is installed
which spooky

# Check available commands
spooky --help

# Verify secrets command
spooky secrets --help
```

#### "Invalid command syntax"

**Cause:** Command syntax is incorrect.

**Solution:**
```bash
# Correct command syntax
spooky secrets encrypt input.txt --output encrypted.txt.age

# Or use default output
spooky secrets encrypt input.txt

# For decryption
spooky secrets decrypt input.txt.age --output decrypted.txt

# For HCL processing
spooky secrets encrypt-hcl config.hcl
spooky secrets decrypt-hcl config.hcl.age
```

### Performance Issues

#### "Encryption is very slow"

**Cause:** Large files or inefficient processing.

**Solution:**
```bash
# Use streaming for large files
spooky secrets encrypt --stream large_file.txt

# Process files in chunks
split -b 1M large_file.txt chunk_
for chunk in chunk_*; do
    spooky secrets encrypt "$chunk"
done

# Monitor system resources
top -p $(pgrep spooky)
```

#### "High memory usage during encryption"

**Cause:** Loading entire files into memory.

**Solution:**
```bash
# Use streaming mode
spooky secrets encrypt --stream input.txt

# Process smaller files
split -b 100K large_file.txt small_chunk_
for chunk in small_chunk_*; do
    spooky secrets encrypt "$chunk"
done
```

## Configuration Problems

### Age Configuration Issues

#### "Invalid age configuration"

**Cause:** Age configuration is incorrect.

**Solution:**
```bash
# Check age configuration
cat ~/.config/spooky/spooky.hcl

# Valid age configuration:
age {
  identity_path = "~/.config/spooky/identities/identity.txt"
  recipients_file = "~/.config/spooky/recipients.txt"
}

# Fix configuration
mkdir -p ~/.config/spooky
cat > ~/.config/spooky/spooky.hcl << 'EOF'
age {
  identity_path = "~/.config/spooky/identities/identity.txt"
  recipients_file = "~/.config/spooky/recipients.txt"
}
EOF
```

#### "Missing age configuration"

**Cause:** Age configuration is not set up.

**Solution:**
```bash
# Create age configuration
mkdir -p ~/.config/spooky/identities

# Generate identity
age-keygen -o ~/.config/spooky/identities/identity.txt

# Create recipients file
echo "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p" > ~/.config/spooky/recipients.txt

# Set permissions
chmod 600 ~/.config/spooky/identities/identity.txt
chmod 600 ~/.config/spooky/recipients.txt
```

### File Permission Issues

#### "Permission denied accessing identity file"

**Cause:** Identity file has wrong permissions.

**Solution:**
```bash
# Check current permissions
ls -la ~/.config/spooky/identities/identity.txt

# Fix permissions
chmod 600 ~/.config/spooky/identities/identity.txt
chmod 700 ~/.config/spooky/identities/

# Verify permissions
ls -la ~/.config/spooky/identities/identity.txt
```

#### "Permission denied creating output file"

**Cause:** Output directory has insufficient permissions.

**Solution:**
```bash
# Check directory permissions
ls -la /path/to/output/directory/

# Fix directory permissions
chmod 755 /path/to/output/directory/

# Or use different output location
spooky secrets encrypt input.txt --output ~/encrypted.txt.age
```

## Debugging Techniques

### Enable Verbose Output

```bash
# Enable verbose output for secrets operations
spooky secrets encrypt --verbose input.txt
spooky secrets decrypt --verbose input.txt.age

# Enable debug logging
export SPOOKY_LOG_LEVEL=debug
spooky secrets encrypt input.txt
```

### Test Age Functionality

```bash
# Test age CLI directly
echo "test data" | age -r $(head -1 ~/.config/spooky/recipients.txt) | age -d -i ~/.config/spooky/identities/identity.txt

# Test identity file
age-keygen -y ~/.config/spooky/identities/identity.txt

# Test recipients
cat ~/.config/spooky/recipients.txt | while read recipient; do
    echo "Testing recipient: $recipient"
    echo "test" | age -r "$recipient" > /dev/null && echo "✓ Valid" || echo "✗ Invalid"
done
```

### Validate Configuration

```bash
# Validate secrets configuration
spooky secrets validate ./my-project

# Check specific configuration
spooky secrets validate --config ~/.config/spooky/spooky.hcl

# Test HCL processing
spooky secrets validate-hcl config.hcl
```

## Recovery Procedures

### Key Recovery

```bash
# Backup identity file
cp ~/.config/spooky/identities/identity.txt ~/.config/spooky/identities/identity.txt.backup

# Generate new identity if needed
age-keygen -o ~/.config/spooky/identities/identity.txt

# Set permissions
chmod 600 ~/.config/spooky/identities/identity.txt
```

### Configuration Recovery

```bash
# Backup configuration
cp ~/.config/spooky/spooky.hcl ~/.config/spooky/spooky.hcl.backup

# Restore from backup if needed
cp ~/.config/spooky/spooky.hcl.backup ~/.config/spooky/spooky.hcl

# Validate configuration
spooky secrets validate --config ~/.config/spooky/spooky.hcl
```

### Data Recovery

```bash
# Test decryption with different identities
for identity in ~/.config/spooky/identities/*.txt; do
    echo "Testing identity: $identity"
    spooky secrets decrypt --identity "$identity" data.txt.age
done

# Try manual age decryption
age -d -i ~/.config/spooky/identities/identity.txt data.txt.age
```

## Prevention Strategies

### Regular Key Rotation

```bash
# Schedule key rotation
crontab -e
# Add: 0 2 1 * * /usr/local/bin/spooky secrets rotate-keys

# Manual key rotation
spooky secrets rotate-keys

# Backup old keys
cp ~/.config/spooky/identities/identity.txt ~/.config/spooky/identities/identity.txt.old
```

### Monitoring

```bash
# Monitor encryption operations
spooky secrets encrypt --verbose input.txt

# Monitor system resources
top -p $(pgrep spooky)

# Check log files
tail -f /var/log/spooky/secrets.log
```

### Backup Strategy

```bash
# Backup identity files
cp ~/.config/spooky/identities/identity.txt ~/.config/spooky/identities/identity.txt.$(date +%Y%m%d)

# Backup configuration
cp ~/.config/spooky/spooky.hcl ~/.config/spooky/spooky.hcl.$(date +%Y%m%d)

# Version control configuration (not identity files!)
git add ~/.config/spooky/spooky.hcl
git commit -m "Update secrets configuration"
```

## Best Practices for Troubleshooting

### 1. Start Simple

Begin with simple encryption and add complexity gradually:

```bash
# Start with basic encryption
echo "test" > test.txt
spooky secrets encrypt test.txt

# Then add complexity
spooky secrets encrypt --recipients age1... test.txt
spooky secrets encrypt-hcl config.hcl
```

### 2. Use Proper Key Management

Implement proper key management practices:

```bash
# Generate keys securely
age-keygen -o ~/.config/spooky/identities/identity.txt

# Set proper permissions
chmod 600 ~/.config/spooky/identities/identity.txt

# Backup keys securely
cp ~/.config/spooky/identities/identity.txt /secure/backup/location/
```

### 3. Validate Early and Often

Validate configurations frequently:

```bash
# Validate after every change
spooky secrets validate ./my-project

# Validate before operations
spooky secrets validate ./my-project && spooky secrets encrypt data.txt

# Validate in scripts
#!/bin/bash
if spooky secrets validate ./my-project; then
    spooky secrets encrypt data.txt
else
    echo "Validation failed"
    exit 1
fi
```

### 4. Use Proper Error Handling

Implement proper error handling in scripts:

```bash
#!/bin/bash
# Encrypt with error handling
if spooky secrets encrypt input.txt; then
    echo "Encryption successful"
else
    echo "Encryption failed"
    exit 1
fi
```

### 5. Monitor and Log

Monitor secrets operations and maintain logs:

```bash
# Enable verbose logging
spooky secrets encrypt --verbose input.txt

# Monitor operations
watch -n 1 'ps aux | grep spooky'

# Check logs
tail -f /var/log/spooky/secrets.log
```

## Getting Help

### Documentation Resources

1. **User Guide** - For usage questions and best practices
2. **API Reference** - For technical implementation details
3. **Examples** - For configuration patterns and use cases

### Common Questions

#### "Why can't I encrypt my data?"

1. Check age CLI tools installation
2. Verify key configuration
3. Check file permissions
4. Validate recipients

#### "How do I debug encryption issues?"

```bash
# Enable verbose output
spooky secrets encrypt --verbose input.txt

# Test age CLI directly
echo "test" | age -r $(head -1 ~/.config/spooky/recipients.txt)

# Check configuration
spooky secrets validate ./my-project
```

#### "How do I fix decryption issues?"

```bash
# Check identity file
ls -la ~/.config/spooky/identities/identity.txt

# Test decryption manually
age -d -i ~/.config/spooky/identities/identity.txt data.txt.age

# Verify file integrity
file data.txt.age
```

#### "How do I optimize secrets operations?"

```bash
# Use streaming for large files
spooky secrets encrypt --stream large_file.txt

# Process files in chunks
split -b 1M large_file.txt chunk_
for chunk in chunk_*; do
    spooky secrets encrypt "$chunk"
done

# Monitor resource usage
top -p $(pgrep spooky)
```

### When to Seek Additional Help

- Configuration validation passes but encryption still fails
- Performance issues persist after optimization
- Unusual error messages not covered in this guide
- Integration issues with other spooky components

For additional help, refer to the [User Guide](SECRETS_USER_GUIDE.md) and [API Reference](SECRETS_API_REFERENCE.md), or check the project documentation for more advanced troubleshooting techniques.

## Conclusion

The secrets system provides robust, reliable encryption and key management with comprehensive age-based security capabilities. Most issues can be resolved by following the troubleshooting steps outlined in this guide. For persistent issues, enable verbose output and collect diagnostic information for further analysis.
