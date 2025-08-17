# Facts Validation Rules
# Extracted from internal/schemas/schemas/structure/facts.hcl
# These rules validate schema compliance and data format correctness

# Facts validation rules
validation_rules {
  # Cross-field validation rules
  cross_field_validation {
    # Encryption metadata must be present when encrypted = true
    rule {
      name = "encryption_metadata_required"
      description = "Encryption metadata must be present when encrypted = true"
      condition = "encrypted == true && encryption_metadata == null"
      message = "Encryption metadata is required when encrypted = true"
      severity = "error"
    }
    
    # Encryption metadata must not be present when encrypted = false
    rule {
      name = "encryption_metadata_not_allowed"
      description = "Encryption metadata must not be present when encrypted = false"
      condition = "encrypted == false && encryption_metadata != null"
      message = "Encryption metadata is not allowed when encrypted = false"
      severity = "error"
    }
    
    # Encrypted facts must have type = "encrypted"
    rule {
      name = "encrypted_type_consistency"
      description = "Encrypted facts must have type = 'encrypted'"
      condition = "encrypted == true && type != 'encrypted'"
      message = "Encrypted facts must have type = 'encrypted'"
      severity = "error"
    }
    
    # Non-encrypted facts cannot have type = "encrypted"
    rule {
      name = "non_encrypted_type_consistency"
      description = "Non-encrypted facts cannot have type = 'encrypted'"
      condition = "encrypted == false && type == 'encrypted'"
      message = "Non-encrypted facts cannot have type = 'encrypted'"
      severity = "error"
    }
    
    # Age-encrypted values must start with 'age1'
    rule {
      name = "age_encrypted_format"
      description = "Age-encrypted values must start with 'age1'"
      condition = "encrypted == true && !value.startsWith('age1')"
      message = "Age-encrypted values must start with 'age1'"
      severity = "error"
    }
  }
  
  # Age encryption specific rules
  age_encryption_rules {
    # Age1 prefix detection for custom facts
    rule {
      name = "age1_prefix_detection"
      description = "Detect age1 prefix for custom facts"
      condition = "value.startsWith('age1') && encrypted == false"
      message = "Values starting with 'age1' should be marked as encrypted = true"
      severity = "warning"
    }
    
    # No regex validation for age-encrypted values
    rule {
      name = "no_regex_for_age_encrypted"
      description = "Do not apply regex validation to age-encrypted values"
      condition = "encrypted == true"
      message = "Age-encrypted values bypass regex validation"
      severity = "info"
    }
  }
}
