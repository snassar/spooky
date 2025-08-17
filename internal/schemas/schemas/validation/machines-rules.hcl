# Machine Validation Rules
# Extracted from internal/schemas/schemas/structure/machines.hcl
# These rules validate schema compliance and data format correctness

# Machine validation rules
validation_rules {
  # Cross-field validation rules
  cross_field_validation {
    # Authentication method validation
    rule {
      name = "ssh_key_method_validation"
      description = "SSH key method requires key_path"
      condition = "authentication.method == 'ssh_key' && authentication.key_path == null"
      message = "SSH key authentication requires key_path"
    }
    
    rule {
      name = "password_method_validation"
      description = "Password method requires password"
      condition = "authentication.method == 'password' && authentication.password == null"
      message = "Password authentication requires password"
    }
    
    rule {
      name = "certificate_method_validation"
      description = "Certificate method requires certificate_path and certificate_key_path"
      condition = "authentication.method == 'certificate' && (authentication.certificate_path == null || authentication.certificate_key_path == null)"
      message = "Certificate authentication requires certificate_path and certificate_key_path"
    }
    
    # Encryption metadata validation
    rule {
      name = "passphrase_encryption_metadata_required"
      description = "Encryption metadata must be present when passphrase is encrypted"
      condition = "authentication.passphrase.encrypted == true && authentication.passphrase.encryption_metadata == null"
      message = "Encryption metadata is required when passphrase is encrypted"
    }
    
    rule {
      name = "password_encryption_metadata_required"
      description = "Encryption metadata must be present when password is encrypted"
      condition = "authentication.password.encrypted == true && authentication.password.encryption_metadata == null"
      message = "Encryption metadata is required when password is encrypted"
    }
    
    rule {
      name = "certificate_passphrase_encryption_metadata_required"
      description = "Encryption metadata must be present when certificate passphrase is encrypted"
      condition = "authentication.certificate_key_passphrase.encrypted == true && authentication.certificate_key_passphrase.encryption_metadata == null"
      message = "Encryption metadata is required when certificate passphrase is encrypted"
    }
    
    # Age-encrypted values must start with 'age1'
    rule {
      name = "passphrase_age_encrypted_format"
      description = "Age-encrypted passphrases must start with 'age1'"
      condition = "authentication.passphrase.encrypted == true && !authentication.passphrase.value.startsWith('age1')"
      message = "Age-encrypted passphrases must start with 'age1'"
    }
    
    rule {
      name = "password_age_encrypted_format"
      description = "Age-encrypted passwords must start with 'age1'"
      condition = "authentication.password.encrypted == true && !authentication.password.value.startsWith('age1')"
      message = "Age-encrypted passwords must start with 'age1'"
    }
    
    rule {
      name = "certificate_passphrase_age_encrypted_format"
      description = "Age-encrypted certificate passphrases must start with 'age1'"
      condition = "authentication.certificate_key_passphrase.encrypted == true && !authentication.certificate_key_passphrase.value.startsWith('age1')"
      message = "Age-encrypted certificate passphrases must start with 'age1'"
    }
  }
  
  # Age encryption specific rules
  age_encryption_rules {
    # Age1 prefix detection for authentication values
    rule {
      name = "age1_prefix_detection_passphrase"
      description = "Detect age1 prefix for passphrases"
      condition = "authentication.passphrase.value.startsWith('age1') && authentication.passphrase.encrypted == false"
      message = "Values starting with 'age1' should be marked as encrypted = true"
    }
    
    rule {
      name = "age1_prefix_detection_password"
      description = "Detect age1 prefix for passwords"
      condition = "authentication.password.value.startsWith('age1') && authentication.password.encrypted == false"
      message = "Values starting with 'age1' should be marked as encrypted = true"
    }
    
    rule {
      name = "age1_prefix_detection_certificate_passphrase"
      description = "Detect age1 prefix for certificate passphrases"
      condition = "authentication.certificate_key_passphrase.value.startsWith('age1') && authentication.certificate_key_passphrase.encrypted == false"
      message = "Values starting with 'age1' should be marked as encrypted = true"
    }
  }
}
