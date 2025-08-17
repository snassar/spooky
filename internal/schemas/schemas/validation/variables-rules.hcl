# Variables Validation Rules
# Extracted from internal/schemas/schemas/structure/variables.hcl
# These rules validate schema compliance and data format correctness

# Variables validation rules
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
    
    # Value type validation for encryption
    rule {
      name = "encryptable_value_types"
      description = "Only string, number, and object values can be encrypted"
      condition = "encrypted == true && (value_type == 'bool')"
      message = "Boolean values cannot be encrypted"
      severity = "error"
    }
  }
}
