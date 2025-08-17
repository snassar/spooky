# Template Metadata Validation Rules
# Custom business logic validation rules for template metadata definitions
# These rules extend the basic structure validation with cross-field and business logic validation

# Template metadata validation rules
validation_rules {
  # Template metadata validation
  template_metadata_validation {
    rule {
      name = "template_name_format"
      description = "Template names must start with a letter and contain only alphanumeric characters, dots, underscores, and hyphens"
      condition = "name != null && !name.matches('^[a-zA-Z][a-zA-Z0-9._-]*$')"
      message = "Template names must start with a letter and contain only alphanumeric characters, dots, underscores, and hyphens"
      severity = "error"
    }
    
    rule {
      name = "template_name_unique"
      description = "Template names must be unique within the project"
      condition = "duplicateTemplateName(name)"
      message = "Template name must be unique within the project"
      severity = "error"
    }
  }
  
  # Cross-field validation
  cross_field_validation {
    rule {
      name = "category_subcategory_consistency"
      description = "If category is specified, subcategory should also be specified"
      condition = "category != null && subcategory == null"
      message = "If category is specified, subcategory should also be specified"
      severity = "warning"
    }
    
    rule {
      name = "dependency_validation"
      description = "Dependencies should reference valid template names"
      condition = "dependencies != null && !validateDependencies(dependencies)"
      message = "Dependencies should reference valid template names"
      severity = "error"
    }
    
    rule {
      name = "version_compatibility"
      description = "Version should be compatible with spooky version"
      condition = "version != null && !isCompatibleVersion(version)"
      message = "Version should be compatible with spooky version"
      severity = "warning"
    }
  }
  
  # Content validation
  content_validation {
    rule {
      name = "sensitive_data_check"
      description = "Description should not contain sensitive information"
      condition = "description != null && description.matches('(?i)(password|secret|key|token|credential)')"
      message = "Description may contain sensitive information"
      severity = "warning"
    }
    
    rule {
      name = "tag_meaningfulness"
      description = "Tags should be meaningful and at least 2 characters long"
      condition = "tags != null && tags.any(tag -> tag.length() < 2)"
      message = "Tags should be at least 2 characters long"
      severity = "warning"
    }
  }
  
  # Performance validation
  performance_validation {
    rule {
      name = "metadata_size"
      description = "Metadata size should be reasonable"
      condition = "getMetadataSize() > 1048576"  # 1MB
      message = "Metadata size exceeds maximum allowed size"
      severity = "error"
    }
    
    rule {
      name = "keyword_count"
      description = "Number of keywords should be reasonable"
      condition = "keywords != null && keywords.size() > 50"
      message = "Too many keywords may impact search performance"
      severity = "warning"
    }
  }
}
