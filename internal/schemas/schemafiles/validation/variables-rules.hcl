# Variables Validation Rules
# Validation rules for variables.hcl schema
# These rules validate schema compliance and data format correctness

# Variables validation rules
validation_rules {
  # Field-specific validation rules
  field_validation {
    # Variable name validation
    rule {
      name = "variable_name_format"
      description = "Variable names must be valid identifiers"
      condition = "name != null && !name.matches('^[a-zA-Z_][a-zA-Z0-9_]*$')"
      message = "Variable names must be valid identifiers"
      severity = "error"
    }
    
    # Variable description validation
    rule {
      name = "description_length"
      description = "Variable descriptions should be reasonable length"
      condition = "description != null && description.length() > 500"
      message = "Variable descriptions should not exceed 500 characters"
      severity = "warning"
    }
    
    # Tags validation
    rule {
      name = "tags_count_limit"
      description = "Variables cannot have more than 20 tags"
      condition = "tags != null && tags.size() > 20"
      message = "Variables cannot have more than 20 tags"
      severity = "warning"
    }
    
    rule {
      name = "tag_format"
      description = "Tags must contain only alphanumeric characters, underscores, and hyphens"
      condition = "tags != null && tags.any(tag -> !tag.matches('^[a-zA-Z0-9_-]+$'))"
      message = "Tags must contain only alphanumeric characters, underscores, and hyphens"
      severity = "error"
    }
    
    rule {
      name = "tag_length"
      description = "Tags must be between 1 and 32 characters"
      condition = "tags != null && tags.any(tag -> tag.length() < 1 || tag.length() > 32)"
      message = "Tags must be between 1 and 32 characters"
      severity = "error"
    }
    
    # Encryption metadata validation
    rule {
      name = "recipients_count"
      description = "Encryption must have between 1 and 10 recipients"
      condition = "encryption_metadata != null && encryption_metadata.recipients != null && (encryption_metadata.recipients.size() < 1 || encryption_metadata.recipients.size() > 10)"
      message = "Encryption must have between 1 and 10 recipients"
      severity = "error"
    }
    
    rule {
      name = "age_public_key_format"
      description = "Age public keys must start with 'age1' and be at least 50 characters"
      condition = "encryption_metadata != null && encryption_metadata.recipients != null && encryption_metadata.recipients.any(key -> !key.matches('^age1[a-z0-9]{50,}$'))"
      message = "Age public keys must start with 'age1' and be at least 50 characters"
      severity = "error"
    }
    
    rule {
      name = "encryption_method_valid"
      description = "Encryption method must be 'age'"
      condition = "encryption_metadata != null && encryption_metadata.method != null && encryption_metadata.method != 'age'"
      message = "Encryption method must be 'age'"
      severity = "error"
    }
    
    rule {
      name = "timestamp_format"
      description = "Timestamps must be in ISO 8601 format"
      condition = "(encryption_metadata != null && encryption_metadata.encrypted_at != null && !encryption_metadata.encrypted_at.matches('^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$')) || (metadata != null && metadata.created_at != null && !metadata.created_at.matches('^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$')) || (metadata != null && metadata.modified_at != null && !metadata.modified_at.matches('^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$'))"
      message = "Timestamps must be in ISO 8601 format"
      severity = "error"
    }
    
    rule {
      name = "version_format"
      description = "Version identifiers must contain only alphanumeric characters, dots, underscores, and hyphens"
      condition = "metadata != null && metadata.version != null && !metadata.version.matches('^[a-zA-Z0-9._-]+$')"
      message = "Version identifiers must contain only alphanumeric characters, dots, underscores, and hyphens"
      severity = "error"
    }
    
    rule {
      name = "version_length"
      description = "Version identifiers must not exceed 32 characters"
      condition = "metadata != null && metadata.version != null && metadata.version.length() > 32"
      message = "Version identifiers must not exceed 32 characters"
      severity = "error"
    }
  }
  
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
    
    # Metadata consistency validation
    rule {
      name = "created_before_modified"
      description = "Created timestamp must be before or equal to modified timestamp"
      condition = "metadata != null && metadata.created_at != null && metadata.modified_at != null && metadata.created_at > metadata.modified_at"
      message = "Created timestamp must be before or equal to modified timestamp"
      severity = "error"
    }
  }
}
