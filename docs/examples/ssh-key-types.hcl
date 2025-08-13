# SSH Key Types Example
# This file demonstrates different supported key types for spooky SSH connections

machines {
  # ED25519 Key (Recommended - Modern, secure, efficient)
  machine "ed25519-server" {
    host = "192.168.1.10"
    user = "admin"
    port = 22
    
    key_file = "~/.ssh/id_ed25519"
    passphrase = "my-secure-passphrase"
    
    tags = ["ed25519", "modern", "production"]
    groups = ["modern-servers"]
    
    metadata {
      environment = "production"
      key_type = "ed25519"
      key_size = "256"  # Fixed size
      security_level = "high"
      performance = "excellent"
    }
  }
  
  # ED25519-SK Key (Hardware Security Key - Planned feature)
  machine "ed25519-sk-server" {
    host = "192.168.1.20"
    user = "admin"
    port = 22
    
    key_file = "~/.ssh/id_ed25519_sk"
    passphrase = "my-secure-passphrase"
    
    tags = ["ed25519-sk", "hardware", "production"]
    groups = ["hardware-key-servers"]
    
    metadata {
      environment = "production"
      key_type = "ed25519-sk"
      key_size = "256"  # Fixed size
      security_level = "very_high"
      hardware_backed = true
      performance = "excellent"
    }
  }
  
  # RSA 4096-bit Key (Traditional, high security)
  machine "rsa-4096-server" {
    host = "192.168.1.30"
    user = "admin"
    port = 22
    
    key_file = "~/.ssh/id_rsa_4096"
    passphrase = "my-secure-passphrase"
    
    tags = ["rsa-4096", "traditional", "production"]
    groups = ["traditional-servers"]
    
    metadata {
      environment = "production"
      key_type = "rsa-4096"
      key_size = "4096"
      security_level = "high"
      performance = "good"
      legacy_compatible = true
    }
  }
  
  # RSA 8192-bit Key (Maximum security)
  machine "rsa-8192-server" {
    host = "192.168.1.40"
    user = "admin"
    port = 22
    
    key_file = "~/.ssh/id_rsa_8192"
    passphrase = "my-secure-passphrase"
    
    tags = ["rsa-8192", "maximum-security", "production"]
    groups = ["maximum-security-servers"]
    
    metadata {
      environment = "production"
      key_type = "rsa-8192"
      key_size = "8192"
      security_level = "maximum"
      performance = "acceptable"
      legacy_compatible = true
    }
  }
  
  # Multiple key types for different environments
  machine "multi-key-server" {
    host = "192.168.1.50"
    user = "admin"
    port = 22
    
    # Primary key (ED25519)
    key_file = "~/.ssh/id_ed25519"
    passphrase = "my-secure-passphrase"
    
    # Fallback key (RSA 4096) - for legacy compatibility
    # Note: This would require additional configuration in the SSH system
    
    tags = ["multi-key", "production", "legacy-compatible"]
    groups = ["multi-key-servers"]
    
    metadata {
      environment = "production"
      primary_key_type = "ed25519"
      fallback_key_type = "rsa-4096"
      security_level = "high"
      legacy_compatible = true
    }
  }
}

# Key Generation Commands Reference
# 
# Note: spooky does not generate SSH keys
# Use openssh to generate keys:
#
# Generate ED25519 key:
#   ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -C "spooky-ed25519-key"
#   ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -C "spooky-ed25519-key" -N "my-passphrase"
#
# Generate ED25519-SK key (hardware security key):
#   ssh-keygen -t ed25519-sk -f ~/.ssh/id_ed25519_sk -C "spooky-ed25519-sk-key"
#   # Note: Requires hardware security key (planned feature)
#
# Generate RSA 4096-bit key:
#   ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa_4096 -C "spooky-rsa-4096-key"
#   ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa_4096 -C "spooky-rsa-4096-key" -N "my-passphrase"
#
# Generate RSA 8192-bit key:
#   ssh-keygen -t rsa -b 8192 -f ~/.ssh/id_rsa_8192 -C "spooky-rsa-8192-key"
#   ssh-keygen -t rsa -b 8192 -f ~/.ssh/id_rsa_8192 -C "spooky-rsa-8192-key" -N "my-passphrase"
#
# Set proper permissions:
#   chmod 600 ~/.ssh/id_ed25519
#   chmod 644 ~/.ssh/id_ed25519.pub
#   chmod 700 ~/.ssh
#
# Validate key:
#   ssh-keygen -lf ~/.ssh/id_ed25519
