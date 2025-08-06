# Template Metadata Schema
# Enhanced template metadata schema with composition pattern
# Includes template-structure base and metadata-specific features
template_metadata {
  include = "template-structure"
  
  scope = "metadata"
  storage_location = "<project>/templates/metadata.hcl"
  description = "Template metadata and configuration information"
  
  # Metadata-specific features
  metadata_features = {
    type = "object"
    description = "Metadata-specific features and capabilities"
    
    properties = {
      # Metadata categories
      metadata_categories = {
        type = "object"
        description = "Metadata categories and capabilities"
        properties = {
          basic_metadata = {
            type = "object"
            description = "Basic template metadata"
            properties = {
              name = {
                type = "boolean"
                value = true
                description = "Template name support"
              }
              description = {
                type = "boolean"
                value = true
                description = "Template description support"
              }
              version = {
                type = "boolean"
                value = true
                description = "Template version support"
              }
              author = {
                type = "boolean"
                value = true
                description = "Template author support"
              }
            }
          }
          
          classification_metadata = {
            type = "object"
            description = "Template classification metadata"
            properties = {
              tags = {
                type = "boolean"
                value = true
                description = "Template tags support"
              }
              categories = {
                type = "boolean"
                value = true
                description = "Template categories support"
              }
              output_format = {
                type = "boolean"
                value = true
                description = "Output format classification"
              }
              template_type = {
                type = "boolean"
                value = true
                description = "Template type classification"
              }
            }
          }
          
          lifecycle_metadata = {
            type = "object"
            description = "Template lifecycle metadata"
            properties = {
              created_at = {
                type = "boolean"
                value = true
                description = "Creation timestamp support"
              }
              updated_at = {
                type = "boolean"
                value = true
                description = "Update timestamp support"
              }
              deprecated_at = {
                type = "boolean"
                value = true
                description = "Deprecation timestamp support"
              }
              removed_at = {
                type = "boolean"
                value = true
                description = "Removal timestamp support"
              }
            }
          }
          
          dependency_metadata = {
            type = "object"
            description = "Template dependency metadata"
            properties = {
              dependencies = {
                type = "boolean"
                value = true
                description = "Template dependencies support"
              }
              requirements = {
                type = "boolean"
                value = true
                description = "Template requirements support"
              }
              conflicts = {
                type = "boolean"
                value = true
                description = "Template conflicts support"
              }
              provides = {
                type = "boolean"
                value = true
                description = "Template provides support"
              }
            }
          }
          
          documentation_metadata = {
            type = "object"
            description = "Template documentation metadata"
            properties = {
              usage = {
                type = "boolean"
                value = true
                description = "Usage documentation support"
              }
              examples = {
                type = "boolean"
                value = true
                description = "Example documentation support"
              }
              api_reference = {
                type = "boolean"
                value = true
                description = "API reference documentation support"
              }
              changelog = {
                type = "boolean"
                value = true
                description = "Changelog documentation support"
              }
            }
          }
        }
      }
      
      # Metadata management features
      metadata_management = {
        type = "object"
        description = "Metadata management capabilities"
        properties = {
          metadata_validation = {
            type = "boolean"
            value = true
            description = "Support for metadata validation"
          }
          metadata_inheritance = {
            type = "boolean"
            value = true
            description = "Support for metadata inheritance"
          }
          metadata_merging = {
            type = "boolean"
            value = true
            description = "Support for metadata merging"
          }
          metadata_versioning = {
            type = "boolean"
            value = true
            description = "Support for metadata versioning"
          }
          metadata_migration = {
            type = "boolean"
            value = true
            description = "Support for metadata migration"
          }
        }
      }
      
      # Metadata discovery features
      metadata_discovery = {
        type = "object"
        description = "Metadata discovery capabilities"
        properties = {
          metadata_indexing = {
            type = "boolean"
            value = true
            description = "Support for metadata indexing"
          }
          metadata_search = {
            type = "boolean"
            value = true
            description = "Support for metadata search"
          }
          metadata_filtering = {
            type = "boolean"
            value = true
            description = "Support for metadata filtering"
          }
          metadata_sorting = {
            type = "boolean"
            value = true
            description = "Support for metadata sorting"
          }
          metadata_grouping = {
            type = "boolean"
            value = true
            description = "Support for metadata grouping"
          }
        }
      }
      
      # Metadata export features
      metadata_export = {
        type = "object"
        description = "Metadata export capabilities"
        properties = {
          export_formats = {
            type = "array"
            items = {
              type = "string"
            }
            value = ["json", "hcl"]
            description = "Supported export formats"
          }
          export_validation = {
            type = "boolean"
            value = true
            description = "Validate exported metadata"
          }
          export_compression = {
            type = "boolean"
            value = true
            description = "Support for export compression"
          }
          export_encryption = {
            type = "boolean"
            value = false
            description = "Support for export encryption"
          }
        }
      }
    }
  }
  
  # Metadata-specific validation
  metadata_validation = {
    type = "object"
    description = "Metadata-specific validation rules"
    
    properties = {
      # Basic metadata validation
      basic_metadata_validation = {
        type = "object"
        description = "Validation rules for basic metadata"
        properties = {
          name_validation = {
            type = "object"
            description = "Name validation rules"
            properties = {
              name_pattern = {
                type = "string"
                value = "^[a-zA-Z0-9._-]+$"
                description = "Regex pattern for template names"
              }
              name_length = {
                type = "object"
                properties = {
                  min = {
                    type = "integer"
                    value = 1
                    description = "Minimum name length"
                  }
                  max = {
                    type = "integer"
                    value = 100
                    description = "Maximum name length"
                  }
                }
              }
              name_uniqueness = {
                type = "boolean"
                value = true
                description = "Require unique template names"
              }
            }
          }
          
          description_validation = {
            type = "object"
            description = "Description validation rules"
            properties = {
              description_length = {
                type = "object"
                properties = {
                  min = {
                    type = "integer"
                    value = 0
                    description = "Minimum description length"
                  }
                  max = {
                    type = "integer"
                    value = 1000
                    description = "Maximum description length"
                  }
                }
              }
              description_format = {
                type = "string"
                value = "text"
                enum = ["text", "markdown", "html"]
                description = "Description format"
              }
            }
          }
          
          version_validation = {
            type = "object"
            description = "Version validation rules"
            properties = {
              version_pattern = {
                type = "string"
                value = "^\\d+\\.\\d+\\.\\d+(-[a-zA-Z0-9._-]+)?(\\+[a-zA-Z0-9._-]+)?$"
                description = "Semantic versioning pattern"
              }
              version_compatibility = {
                type = "boolean"
                value = true
                description = "Check version compatibility"
              }
            }
          }
          
          author_validation = {
            type = "object"
            description = "Author validation rules"
            properties = {
              author_pattern = {
                type = "string"
                value = "^[a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+\\.[a-zA-Z]{2,}$"
                description = "Author email pattern"
              }
              author_required = {
                type = "boolean"
                value = false
                description = "Require author information"
              }
            }
          }
        }
      }
      
      # Classification metadata validation
      classification_metadata_validation = {
        type = "object"
        description = "Validation rules for classification metadata"
        properties = {
          tags_validation = {
            type = "object"
            description = "Tags validation rules"
            properties = {
              tag_pattern = {
                type = "string"
                value = "^[a-zA-Z0-9._-]+$"
                description = "Tag pattern"
              }
              max_tags = {
                type = "integer"
                value = 20
                description = "Maximum number of tags"
              }
              tag_length = {
                type = "object"
                properties = {
                  min = {
                    type = "integer"
                    value = 1
                    description = "Minimum tag length"
                  }
                  max = {
                    type = "integer"
                    value = 50
                    description = "Maximum tag length"
                  }
                }
              }
            }
          }
          
          categories_validation = {
            type = "object"
            description = "Categories validation rules"
            properties = {
              category_pattern = {
                type = "string"
                value = "^[a-zA-Z0-9._-]+$"
                description = "Category pattern"
              }
              max_categories = {
                type = "integer"
                value = 10
                description = "Maximum number of categories"
              }
              allowed_categories = {
                type = "array"
                items = {
                  type = "string"
                }
                value = ["config", "script", "documentation", "data", "template", "other"]
                description = "Allowed category values"
              }
            }
          }
          
          output_format_validation = {
            type = "object"
            description = "Output format validation rules"
            properties = {
              allowed_formats = {
                type = "array"
                items = {
                  type = "string"
                }
                value = ["config", "script", "documentation", "data", "template", "other"]
                description = "Allowed output formats"
              }
              format_required = {
                type = "boolean"
                value = false
                description = "Require output format"
              }
            }
          }
        }
      }
      
      # Lifecycle metadata validation
      lifecycle_metadata_validation = {
        type = "object"
        description = "Validation rules for lifecycle metadata"
        properties = {
          timestamp_validation = {
            type = "object"
            description = "Timestamp validation rules"
            properties = {
              timestamp_format = {
                type = "string"
                value = "RFC3339"
                description = "Timestamp format"
              }
              timestamp_order = {
                type = "boolean"
                value = true
                description = "Validate timestamp order"
              }
              future_timestamps = {
                type = "boolean"
                value = false
                description = "Allow future timestamps"
              }
            }
          }
          
          lifecycle_state_validation = {
            type = "object"
            description = "Lifecycle state validation rules"
            properties = {
              valid_states = {
                type = "array"
                items = {
                  type = "string"
                }
                value = ["active", "deprecated", "removed"]
                description = "Valid lifecycle states"
              }
              state_transitions = {
                type = "object"
                description = "Valid state transitions"
                properties = {
                  active_to_deprecated = {
                    type = "boolean"
                    value = true
                    description = "Allow active to deprecated transition"
                  }
                  deprecated_to_removed = {
                    type = "boolean"
                    value = true
                    description = "Allow deprecated to removed transition"
                  }
                  removed_to_active = {
                    type = "boolean"
                    value = false
                    description = "Allow removed to active transition"
                  }
                }
              }
            }
          }
        }
      }
      
      # Dependency metadata validation
      dependency_metadata_validation = {
        type = "object"
        description = "Validation rules for dependency metadata"
        properties = {
          dependency_validation = {
            type = "object"
            description = "Dependency validation rules"
            properties = {
              dependency_pattern = {
                type = "string"
                value = "^[a-zA-Z0-9._-]+(>=|<=|==|!=|>|<)\\d+\\.\\d+\\.\\d+$"
                description = "Dependency version pattern"
              }
              max_dependencies = {
                type = "integer"
                value = 50
                description = "Maximum number of dependencies"
              }
              circular_dependency_check = {
                type = "boolean"
                value = true
                description = "Check for circular dependencies"
              }
            }
          }
          
          requirement_validation = {
            type = "object"
            description = "Requirement validation rules"
            properties = {
              requirement_pattern = {
                type = "string"
                value = "^[a-zA-Z0-9._-]+$"
                description = "Requirement pattern"
              }
              max_requirements = {
                type = "integer"
                value = 20
                description = "Maximum number of requirements"
              }
            }
          }
          
          conflict_validation = {
            type = "object"
            description = "Conflict validation rules"
            properties = {
              conflict_pattern = {
                type = "string"
                value = "^[a-zA-Z0-9._-]+$"
                description = "Conflict pattern"
              }
              max_conflicts = {
                type = "integer"
                value = 10
                description = "Maximum number of conflicts"
              }
            }
          }
        }
      }
      
      # Documentation metadata validation
      documentation_metadata_validation = {
        type = "object"
        description = "Validation rules for documentation metadata"
        properties = {
          usage_validation = {
            type = "object"
            description = "Usage documentation validation"
            properties = {
              usage_required = {
                type = "boolean"
                value = false
                description = "Require usage documentation"
              }
              usage_format = {
                type = "string"
                value = "markdown"
                enum = ["text", "markdown", "html"]
                description = "Usage documentation format"
              }
              usage_length = {
                type = "object"
                properties = {
                  min = {
                    type = "integer"
                    value = 0
                    description = "Minimum usage length"
                  }
                  max = {
                    type = "integer"
                    value = 10000
                    description = "Maximum usage length"
                  }
                }
              }
            }
          }
          
          example_validation = {
            type = "object"
            description = "Example documentation validation"
            properties = {
              example_required = {
                type = "boolean"
                value = false
                description = "Require example documentation"
              }
              max_examples = {
                type = "integer"
                value = 10
                description = "Maximum number of examples"
              }
              example_format = {
                type = "string"
                value = "markdown"
                enum = ["text", "markdown", "html", "json"]
                description = "Example documentation format"
              }
            }
          }
          
          api_reference_validation = {
            type = "object"
            description = "API reference validation"
            properties = {
              api_reference_required = {
                type = "boolean"
                value = false
                description = "Require API reference"
              }
              api_reference_format = {
                type = "string"
                value = "json"
                enum = ["json", "xml", "openapi"]
                description = "API reference format"
              }
            }
          }
        }
      }
    }
  }
  
  # Metadata-specific constraints
  metadata_constraints = {
    type = "object"
    description = "Metadata-specific constraints and limits"
    
    properties = {
      # Size constraints
      size_constraints = {
        type = "object"
        description = "Size-related constraints"
        properties = {
          max_metadata_size = {
            type = "integer"
            value = 1048576
            description = "Maximum metadata size (1MB)"
          }
          max_name_length = {
            type = "integer"
            value = 100
            description = "Maximum name length"
          }
          max_description_length = {
            type = "integer"
            value = 1000
            description = "Maximum description length"
          }
          max_tag_length = {
            type = "integer"
            value = 50
            description = "Maximum tag length"
          }
          max_tags_count = {
            type = "integer"
            value = 20
            description = "Maximum number of tags"
          }
          max_categories_count = {
            type = "integer"
            value = 10
            description = "Maximum number of categories"
          }
          max_dependencies_count = {
            type = "integer"
            value = 50
            description = "Maximum number of dependencies"
          }
          max_requirements_count = {
            type = "integer"
            value = 20
            description = "Maximum number of requirements"
          }
          max_conflicts_count = {
            type = "integer"
            value = 10
            description = "Maximum number of conflicts"
          }
          max_examples_count = {
            type = "integer"
            value = 10
            description = "Maximum number of examples"
          }
        }
      }
      
      # Performance constraints
      performance_constraints = {
        type = "object"
        description = "Performance-related constraints"
        properties = {
          metadata_validation_timeout = {
            type = "integer"
            value = 5000
            description = "Metadata validation timeout (5s)"
          }
          metadata_processing_timeout = {
            type = "integer"
            value = 10000
            description = "Metadata processing timeout (10s)"
          }
          max_metadata_cache_size = {
            type = "integer"
            value = 1000
            description = "Maximum metadata cache size"
          }
          metadata_cache_ttl = {
            type = "string"
            value = "1h"
            description = "Metadata cache TTL"
          }
        }
      }
      
      # Security constraints
      security_constraints = {
        type = "object"
        description = "Security-related constraints"
        properties = {
          metadata_encryption = {
            type = "boolean"
            value = false
            description = "Encrypt metadata at rest"
          }
          metadata_access_control = {
            type = "boolean"
            value = true
            description = "Enforce metadata access control"
          }
          metadata_audit_logging = {
            type = "boolean"
            value = true
            description = "Log metadata access for auditing"
          }
          sensitive_data_filtering = {
            type = "boolean"
            value = true
            description = "Filter sensitive data from metadata"
          }
        }
      }
      
      # Content constraints
      content_constraints = {
        type = "object"
        description = "Content-related constraints"
        properties = {
          allowed_characters = {
            type = "string"
            value = "a-zA-Z0-9._-"
            description = "Allowed characters in metadata"
          }
          forbidden_patterns = {
            type = "array"
            items = {
              type = "string"
            }
            value = [
              "<script>",
              "javascript:",
              "data:text/html",
              "vbscript:",
              "onload=",
              "onerror="
            ]
            description = "Forbidden patterns in metadata"
          }
          content_sanitization = {
            type = "boolean"
            value = true
            description = "Sanitize metadata content"
          }
          html_escaping = {
            type = "boolean"
            value = true
            description = "Escape HTML in metadata"
          }
        }
      }
      
      # Versioning constraints
      versioning_constraints = {
        type = "object"
        description = "Versioning-related constraints"
        properties = {
          version_scheme = {
            type = "string"
            value = "semver"
            enum = ["semver", "calver", "custom"]
            description = "Version scheme"
          }
          version_compatibility = {
            type = "boolean"
            value = true
            description = "Check version compatibility"
          }
          version_migration = {
            type = "boolean"
            value = true
            description = "Support version migration"
          }
          version_deprecation = {
            type = "boolean"
            value = true
            description = "Support version deprecation"
          }
        }
      }
    }
  }
  
  # Metadata-specific extensions
  metadata_extensions = {
    type = "object"
    description = "Metadata-specific extensions and enhancements"
    
    properties = {
      # Metadata transformation extensions
      metadata_transformation_extensions = {
        type = "object"
        description = "Metadata transformation capabilities"
        properties = {
          format_transformation = {
            type = "object"
            description = "Format transformation capabilities"
            properties = {
              
              hcl_to_json = {
                type = "boolean"
                value = true
                description = "Convert HCL to JSON"
              }
              json_to_hcl = {
                type = "boolean"
                value = true
                description = "Convert JSON to HCL"
              }
            }
          }
          
          content_transformation = {
            type = "object"
            description = "Content transformation capabilities"
            properties = {
              markdown_to_html = {
                type = "boolean"
                value = true
                description = "Convert Markdown to HTML"
              }
              html_to_markdown = {
                type = "boolean"
                value = true
                description = "Convert HTML to Markdown"
              }
              text_normalization = {
                type = "boolean"
                value = true
                description = "Normalize text content"
              }
              character_encoding = {
                type = "boolean"
                value = true
                description = "Handle character encoding"
              }
            }
          }
          
          validation_transformation = {
            type = "object"
            description = "Validation transformation capabilities"
            properties = {
              schema_validation = {
                type = "boolean"
                value = true
                description = "Validate against schema"
              }
              content_validation = {
                type = "boolean"
                value = true
                description = "Validate content integrity"
              }
              format_validation = {
                type = "boolean"
                value = true
                description = "Validate format compliance"
              }
              cross_reference_validation = {
                type = "boolean"
                value = true
                description = "Validate cross-references"
              }
            }
          }
        }
      }
      
      # Metadata indexing extensions
      metadata_indexing_extensions = {
        type = "object"
        description = "Metadata indexing capabilities"
        properties = {
          full_text_indexing = {
            type = "object"
            description = "Full-text indexing capabilities"
            properties = {
              enable_indexing = {
                type = "boolean"
                value = true
                description = "Enable full-text indexing"
              }
              index_fields = {
                type = "array"
                items = {
                  type = "string"
                }
                value = ["name", "description", "tags", "author"]
                description = "Fields to index"
              }
              search_engine = {
                type = "string"
                value = "bleve"
                enum = ["bleve", "elasticsearch", "lucene"]
                description = "Search engine"
              }
              index_refresh_rate = {
                type = "string"
                value = "1s"
                description = "Index refresh rate"
              }
            }
          }
          
          metadata_categorization = {
            type = "object"
            description = "Metadata categorization capabilities"
            properties = {
              auto_categorization = {
                type = "boolean"
                value = true
                description = "Automatic categorization"
              }
              category_suggestions = {
                type = "boolean"
                value = true
                description = "Category suggestions"
              }
              tag_suggestions = {
                type = "boolean"
                value = true
                description = "Tag suggestions"
              }
              similarity_matching = {
                type = "boolean"
                value = true
                description = "Similarity matching"
              }
            }
          }
          
          metadata_analytics = {
            type = "object"
            description = "Metadata analytics capabilities"
            properties = {
              usage_analytics = {
                type = "boolean"
                value = true
                description = "Usage analytics"
              }
              popularity_metrics = {
                type = "boolean"
                value = true
                description = "Popularity metrics"
              }
              quality_metrics = {
                type = "boolean"
                value = true
                description = "Quality metrics"
              }
              trend_analysis = {
                type = "boolean"
                value = true
                description = "Trend analysis"
              }
            }
          }
        }
      }
      
      # Metadata discovery extensions
      metadata_discovery_extensions = {
        type = "object"
        description = "Metadata discovery capabilities"
        properties = {
          search_capabilities = {
            type = "object"
            description = "Search capabilities"
            properties = {
              full_text_search = {
                type = "boolean"
                value = true
                description = "Full-text search"
              }
              fuzzy_search = {
                type = "boolean"
                value = true
                description = "Fuzzy search"
              }
              regex_search = {
                type = "boolean"
                value = true
                description = "Regex search"
              }
              faceted_search = {
                type = "boolean"
                value = true
                description = "Faceted search"
              }
            }
          }
          
          filtering_capabilities = {
            type = "object"
            description = "Filtering capabilities"
            properties = {
              tag_filtering = {
                type = "boolean"
                value = true
                description = "Tag-based filtering"
              }
              category_filtering = {
                type = "boolean"
                value = true
                description = "Category-based filtering"
              }
              author_filtering = {
                type = "boolean"
                value = true
                description = "Author-based filtering"
              }
              version_filtering = {
                type = "boolean"
                value = true
                description = "Version-based filtering"
              }
              date_filtering = {
                type = "boolean"
                value = true
                description = "Date-based filtering"
              }
            }
          }
          
          sorting_capabilities = {
            type = "object"
            description = "Sorting capabilities"
            properties = {
              name_sorting = {
                type = "boolean"
                value = true
                description = "Sort by name"
              }
              date_sorting = {
                type = "boolean"
                value = true
                description = "Sort by date"
              }
              popularity_sorting = {
                type = "boolean"
                value = true
                description = "Sort by popularity"
              }
              quality_sorting = {
                type = "boolean"
                value = true
                description = "Sort by quality"
              }
              relevance_sorting = {
                type = "boolean"
                value = true
                description = "Sort by relevance"
              }
            }
          }
        }
      }
      
      # Metadata export extensions
      metadata_export_extensions = {
        type = "object"
        description = "Metadata export capabilities"
        properties = {
          format_export = {
            type = "object"
            description = "Format export capabilities"
            properties = {
              json_export = {
                type = "boolean"
                value = true
                description = "Export to JSON"
              }
              


            }
          }
          
          content_export = {
            type = "object"
            description = "Content export capabilities"
            properties = {
              full_export = {
                type = "boolean"
                value = true
                description = "Export full metadata"
              }
              partial_export = {
                type = "boolean"
                value = true
                description = "Export partial metadata"
              }
              filtered_export = {
                type = "boolean"
                value = true
                description = "Export filtered metadata"
              }
              templated_export = {
                type = "boolean"
                value = true
                description = "Export using templates"
              }
            }
          }
          
          delivery_export = {
            type = "object"
            description = "Delivery export capabilities"
            properties = {
              file_export = {
                type = "boolean"
                value = true
                description = "Export to file"
              }
              stream_export = {
                type = "boolean"
                value = true
                description = "Export to stream"
              }
              api_export = {
                type = "boolean"
                value = true
                description = "Export via API"
              }
              webhook_export = {
                type = "boolean"
                value = true
                description = "Export via webhook"
              }
            }
          }
        }
      }
    }
  }
  
  # Metadata-specific metadata
  metadata_metadata = {
    type = "object"
    description = "Metadata-specific metadata and information"
    
    properties = {
      # Metadata versioning
      versioning = {
        type = "object"
        description = "Metadata versioning information"
        properties = {
          metadata_version = {
            type = "string"
            value = "1.0.0"
            description = "Metadata schema version"
          }
          schema_version = {
            type = "string"
            value = "1.0.0"
            description = "Schema version"
          }
          api_version = {
            type = "string"
            value = "1.0.0"
            description = "API version"
          }
        }
      }
      
      # Metadata lifecycle
      lifecycle = {
        type = "object"
        description = "Metadata lifecycle information"
        properties = {
          created_at = {
            type = "string"
            required = false
            format = "date-time"
            description = "Metadata creation timestamp"
          }
          updated_at = {
            type = "string"
            required = false
            format = "date-time"
            description = "Metadata last update timestamp"
          }
          expires_at = {
            type = "string"
            required = false
            format = "date-time"
            description = "Metadata expiration timestamp"
          }
          ttl = {
            type = "string"
            value = "1h"
            description = "Metadata time-to-live"
          }
        }
      }
      
      # Metadata ownership
      ownership = {
        type = "object"
        description = "Metadata ownership information"
        properties = {
          owner = {
            type = "string"
            required = false
            description = "Metadata owner"
          }
          maintainer = {
            type = "string"
            required = false
            description = "Metadata maintainer"
          }
          contributors = {
            type = "array"
            items = {
              type = "string"
            }
            required = false
            description = "Metadata contributors"
          }
          license = {
            type = "string"
            required = false
            description = "Metadata license"
          }
        }
      }
      
      # Metadata tags
      tags = {
        type = "array"
        required = false
        description = "Metadata tags"
        items = {
          type = "string"
        }
      }
      
      # Metadata description
      description = {
        type = "string"
        required = false
        description = "Metadata description"
      }
    }
  }
}