# Machine Validation Rules
# Validation rules for machines.hcl schema
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
      severity = "error"
    }
    
    rule {
      name = "password_method_validation"
      description = "Password method requires password"
      condition = "authentication.method == 'password' && authentication.password == null"
      message = "Password authentication requires password"
      severity = "error"
    }
    
    rule {
      name = "certificate_method_validation"
      description = "Certificate method requires certificate_path and certificate_key_path"
      condition = "authentication.method == 'certificate' && (authentication.certificate_path == null || authentication.certificate_key_path == null)"
      message = "Certificate authentication requires certificate_path and certificate_key_path"
      severity = "error"
    }
    
    # Encryption metadata validation
    rule {
      name = "passphrase_encryption_metadata_required"
      description = "Encryption metadata must be present when passphrase is encrypted"
      condition = "authentication.passphrase.encrypted == true && authentication.passphrase.encryption_metadata == null"
      message = "Encryption metadata is required when passphrase is encrypted"
      severity = "error"
    }
    
    rule {
      name = "password_encryption_metadata_required"
      description = "Encryption metadata must be present when password is encrypted"
      condition = "authentication.password.encrypted == true && authentication.password.encryption_metadata == null"
      message = "Encryption metadata is required when password is encrypted"
      severity = "error"
    }
    
    rule {
      name = "certificate_passphrase_encryption_metadata_required"
      description = "Encryption metadata must be present when certificate passphrase is encrypted"
      condition = "authentication.certificate_passphrase.encrypted == true && authentication.certificate_passphrase.encryption_metadata == null"
      message = "Encryption metadata is required when certificate passphrase is encrypted"
      severity = "error"
    }
    
    rule {
      name = "certificate_key_passphrase_encryption_metadata_required"
      description = "Encryption metadata must be present when certificate key passphrase is encrypted"
      condition = "authentication.certificate_key_passphrase.encrypted == true && authentication.certificate_key_passphrase.encryption_metadata == null"
      message = "Encryption metadata is required when certificate key passphrase is encrypted"
      severity = "error"
    }
    
    # Age-encrypted values must start with 'age1'
    rule {
      name = "passphrase_age_encrypted_format"
      description = "Age-encrypted passphrases must start with 'age1'"
      condition = "authentication.passphrase.encrypted == true && !authentication.passphrase.value.startsWith('age1')"
      message = "Age-encrypted passphrases must start with 'age1'"
      severity = "error"
    }
    
    rule {
      name = "password_age_encrypted_format"
      description = "Age-encrypted passwords must start with 'age1'"
      condition = "authentication.password.encrypted == true && !authentication.password.value.startsWith('age1')"
      message = "Age-encrypted passwords must start with 'age1'"
      severity = "error"
    }
    
    rule {
      name = "certificate_passphrase_age_encrypted_format"
      description = "Age-encrypted certificate passphrases must start with 'age1'"
      condition = "authentication.certificate_passphrase.encrypted == true && !authentication.certificate_passphrase.value.startsWith('age1')"
      message = "Age-encrypted certificate passphrases must start with 'age1'"
      severity = "error"
    }
    
    rule {
      name = "certificate_key_passphrase_age_encrypted_format"
      description = "Age-encrypted certificate key passphrases must start with 'age1'"
      condition = "authentication.certificate_key_passphrase.encrypted == true && !authentication.certificate_key_passphrase.value.startsWith('age1')"
      message = "Age-encrypted certificate key passphrases must start with 'age1'"
      severity = "error"
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
      severity = "warning"
    }
    
    rule {
      name = "age1_prefix_detection_password"
      description = "Detect age1 prefix for passwords"
      condition = "authentication.password.value.startsWith('age1') && authentication.password.encrypted == false"
      message = "Values starting with 'age1' should be marked as encrypted = true"
      severity = "warning"
    }
    
    rule {
      name = "age1_prefix_detection_certificate_passphrase"
      description = "Detect age1 prefix for certificate passphrases"
      condition = "authentication.certificate_passphrase.value.startsWith('age1') && authentication.certificate_passphrase.encrypted == false"
      message = "Values starting with 'age1' should be marked as encrypted = true"
      severity = "warning"
    }
    
    rule {
      name = "age1_prefix_detection_certificate_key_passphrase"
      description = "Detect age1 prefix for certificate key passphrases"
      condition = "authentication.certificate_key_passphrase.value.startsWith('age1') && authentication.certificate_key_passphrase.encrypted == false"
      message = "Values starting with 'age1' should be marked as encrypted = true"
      severity = "warning"
    }
  }
  
  # Variable validation rules
  variable_validation {
    # Machine-specific variables validation
    rule {
      name = "machine_variable_encryption_metadata_required"
      description = "Machine variable encryption metadata must be present when encrypted = true"
      condition = "variables != null && variables.any(var -> var.encrypted == true && var.encryption_metadata == null)"
      message = "Machine variable encryption metadata is required when encrypted = true"
      severity = "error"
    }
    
    rule {
      name = "machine_variable_age_encrypted_format"
      description = "Age-encrypted machine variables must start with 'age1'"
      condition = "variables != null && variables.any(var -> var.encrypted == true && !var.value.startsWith('age1'))"
      message = "Age-encrypted machine variables must start with 'age1'"
      severity = "error"
    }
    
    rule {
      name = "machine_variable_age1_prefix_detection"
      description = "Detect age1 prefix for machine variables"
      condition = "variables != null && variables.any(var -> var.value.startsWith('age1') && var.encrypted == false)"
      message = "Machine variable values starting with 'age1' should be marked as encrypted = true"
      severity = "warning"
    }
    
    # Group variables validation
    rule {
      name = "group_variable_encryption_metadata_required"
      description = "Group variable encryption metadata must be present when encrypted = true"
      condition = "group_variables != null && group_variables.any(group -> group.variables.any(var -> var.encrypted == true && var.encryption_metadata == null))"
      message = "Group variable encryption metadata is required when encrypted = true"
      severity = "error"
    }
    
    rule {
      name = "group_variable_age_encrypted_format"
      description = "Age-encrypted group variables must start with 'age1'"
      condition = "group_variables != null && group_variables.any(group -> group.variables.any(var -> var.encrypted == true && !var.value.startsWith('age1')))"
      message = "Age-encrypted group variables must start with 'age1'"
      severity = "error"
    }
    
    rule {
      name = "group_variable_age1_prefix_detection"
      description = "Detect age1 prefix for group variables"
      condition = "group_variables != null && group_variables.any(group -> group.variables.any(var -> var.value.startsWith('age1') && var.encrypted == false))"
      message = "Group variable values starting with 'age1' should be marked as encrypted = true"
      severity = "warning"
    }
    
    # Variable name validation
    rule {
      name = "variable_name_format"
      description = "Variable names must be valid identifiers"
      condition = "(variables != null && variables.any(var -> var.name != null && !var.name.matches('^[a-zA-Z_][a-zA-Z0-9_]*$'))) || (group_variables != null && group_variables.any(group -> group.variables.any(var -> var.name != null && !var.name.matches('^[a-zA-Z_][a-zA-Z0-9_]*$'))))"
      message = "Variable names must be valid identifiers"
      severity = "error"
    }
    
    # Variable tags validation
    rule {
      name = "variable_tags_count_limit"
      description = "Variables cannot have more than 10 tags"
      condition = "(variables != null && variables.any(var -> var.tags != null && var.tags.size() > 10)) || (group_variables != null && group_variables.any(group -> group.variables.any(var -> var.tags != null && var.tags.size() > 10)))"
      message = "Variables cannot have more than 10 tags"
      severity = "warning"
    }
    
    rule {
      name = "variable_tag_format"
      description = "Variable tags must contain only alphanumeric characters, underscores, and hyphens"
      condition = "(variables != null && variables.any(var -> var.tags != null && var.tags.any(tag -> !tag.matches('^[a-zA-Z0-9_-]+$')))) || (group_variables != null && group_variables.any(group -> group.variables.any(var -> var.tags != null && var.tags.any(tag -> !tag.matches('^[a-zA-Z0-9_-]+$')))))"
      message = "Variable tags must contain only alphanumeric characters, underscores, and hyphens"
      severity = "error"
    }
  }
}
