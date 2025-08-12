# SSH Certificates Example
# This file demonstrates SSH certificate authentication for spooky

machines {
  # ED25519 Certificate Authentication
  machine "ed25519-cert-server" {
    host = "192.168.1.10"
    user = "admin"
    port = 22
    
    # Private key (required for certificate authentication)
    key_file = "~/.ssh/id_ed25519"
    passphrase = "my-secure-passphrase"
    
    # Certificate (signed by CA)
    certificate_file = "~/.ssh/id_ed25519-cert.pub"
    
    tags = ["ed25519", "certificate", "production"]
    groups = ["certificate-servers"]
    
    metadata {
      environment = "production"
      key_type = "ed25519"
      auth_method = "certificate"
      certificate_type = "ed25519"
      ca_authority = "spooky-ca"
      certificate_expiry = "2024-12-31"
    }
  }
  
  # RSA 4096-bit Certificate Authentication
  machine "rsa-cert-server" {
    host = "192.168.1.20"
    user = "admin"
    port = 22
    
    # Private key
    key_file = "~/.ssh/id_rsa_4096"
    passphrase = "my-secure-passphrase"
    
    # Certificate
    certificate_file = "~/.ssh/id_rsa_4096-cert.pub"
    
    tags = ["rsa-4096", "certificate", "production"]
    groups = ["certificate-servers"]
    
    metadata {
      environment = "production"
      key_type = "rsa-4096"
      auth_method = "certificate"
      certificate_type = "rsa-4096"
      ca_authority = "spooky-ca"
      certificate_expiry = "2024-12-31"
    }
  }
  
  # Certificate with specific principals
  machine "principal-cert-server" {
    host = "192.168.1.30"
    user = "admin"
    port = 22
    
    key_file = "~/.ssh/id_ed25519"
    passphrase = "my-secure-passphrase"
    certificate_file = "~/.ssh/id_ed25519-admin-cert.pub"
    
    tags = ["certificate", "principal", "production"]
    groups = ["admin-servers"]
    
    metadata {
      environment = "production"
      key_type = "ed25519"
      auth_method = "certificate"
      certificate_principals = ["admin", "root"]
      ca_authority = "spooky-ca"
      certificate_expiry = "2024-12-31"
    }
  }
  
  # Certificate with extensions
  machine "extended-cert-server" {
    host = "192.168.1.40"
    user = "admin"
    port = 22
    
    key_file = "~/.ssh/id_ed25519"
    passphrase = "my-secure-passphrase"
    certificate_file = "~/.ssh/id_ed25519-extended-cert.pub"
    
    tags = ["certificate", "extended", "production"]
    groups = ["extended-servers"]
    
    metadata {
      environment = "production"
      key_type = "ed25519"
      auth_method = "certificate"
      certificate_extensions = ["permit-pty", "permit-port-forwarding"]
      ca_authority = "spooky-ca"
      certificate_expiry = "2024-12-31"
    }
  }
  
  # Short-lived certificate (for temporary access)
  machine "temp-cert-server" {
    host = "192.168.1.50"
    user = "temp-user"
    port = 22
    
    key_file = "~/.ssh/id_ed25519_temp"
    passphrase = "temp-passphrase"
    certificate_file = "~/.ssh/id_ed25519_temp-cert.pub"
    
    tags = ["certificate", "temporary", "development"]
    groups = ["temp-servers"]
    
    metadata {
      environment = "development"
      key_type = "ed25519"
      auth_method = "certificate"
      certificate_type = "temporary"
      ca_authority = "spooky-ca"
      certificate_expiry = "2024-01-31"  # Short expiry
      access_type = "temporary"
    }
  }
}

# Certificate Authority Setup Commands
#
# 1. Generate CA key:
#    ssh-keygen -t ed25519 -f ~/.ssh/ca_key -C "spooky-ca"
#
# 2. Generate user key:
#    ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -C "spooky-user"
#
# 3. Sign certificate for user:
#    ssh-keygen -s ~/.ssh/ca_key -I "spooky-cert" -n "admin" -V +52w ~/.ssh/id_ed25519.pub
#
# 4. Sign certificate with specific principals:
#    ssh-keygen -s ~/.ssh/ca_key -I "spooky-admin-cert" -n "admin,root" -V +52w ~/.ssh/id_ed25519.pub
#
# 5. Sign certificate with extensions:
#    ssh-keygen -s ~/.ssh/ca_key -I "spooky-extended-cert" -n "admin" -O permit-pty -O permit-port-forwarding -V +52w ~/.ssh/id_ed25519.pub
#
# 6. Sign short-lived certificate:
#    ssh-keygen -s ~/.ssh/ca_key -I "spooky-temp-cert" -n "temp-user" -V +7d ~/.ssh/id_ed25519.pub
#
# Certificate Validation Commands
#
# 1. View certificate details:
#    ssh-keygen -L -f ~/.ssh/id_ed25519-cert.pub
#
# 2. Check certificate expiration:
#    ssh-keygen -L -f ~/.ssh/id_ed25519-cert.pub | grep "Valid:"
#
# 3. Verify certificate signature:
#    ssh-keygen -l -f ~/.ssh/id_ed25519-cert.pub
#
# 4. Test certificate authentication:
#    ssh -i ~/.ssh/id_ed25519 -i ~/.ssh/id_ed25519-cert.pub user@host
#
# Certificate Management Best Practices
#
# 1. Keep CA key secure and offline
# 2. Use appropriate certificate validity periods
# 3. Implement certificate revocation procedures
# 4. Monitor certificate expiration
# 5. Use different CAs for different environments
# 6. Implement certificate renewal automation
