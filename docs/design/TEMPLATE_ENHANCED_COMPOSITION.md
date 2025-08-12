# Template Enhanced Composition Pattern

## Overview

This document describes the implementation of the **Enhanced Composition Pattern** for template schemas in the Spooky project. This pattern provides a robust, consistent, and extensible approach to template schema management, following the same architectural principles established for facts schemas.

## Architecture

### **Base Structure Schema**
- **`template-structure.hcl`**: Common template structure definitions for all template-related schemas
- Provides shared validation rules, constraints, and metadata
- Defines the foundation for all template schemas

### **Enhanced Schema Components**
- **`template-context.hcl`**: Enhanced context schema with composition pattern
- **`template-functions.hcl`**: Enhanced functions schema with composition pattern  
- **`template-metadata.hcl`**: Enhanced metadata schema with composition pattern

### **Composition Pattern**
Each enhanced schema includes:
- **`include = "template-structure"`**: References the base structure
- **Schema-specific features**: Unique capabilities and functionality
- **Enhanced validation**: Comprehensive validation rules
- **Extended constraints**: Performance, security, and operational limits
- **Rich extensions**: Advanced capabilities and integrations
- **Detailed metadata**: Versioning, lifecycle, and ownership information

## Benefits

### **✅ Clear Separation**
- Base structure vs schema-specific features
- Common functionality vs specialized capabilities
- Validation rules vs implementation details

### **✅ Consistent Pattern**
- All template schemas follow the same composition pattern
- Unified structure across context, functions, and metadata
- Consistent validation and constraint handling

### **✅ Extensible Design**
- Easy to add new template schema types
- Support for custom extensions and features
- Plugin-based architecture for advanced capabilities

### **✅ Maintainable Code**
- Changes to base structure propagate to all schemas
- Centralized validation and constraint management
- Clear documentation and versioning

## Implementation Details

### **Phase 1: Base Structure Schema**

#### **`template-structure.hcl`**
The foundation schema that provides:

**Core Template Fields**:
- `template_id`: Unique template identifier
- `source_path`: Template file location
- `destination_path`: Output destination
- `template_type`: Schema classification
- `scope`: Usage scope (project, global, machine, system)
- `security_level`: Security classification
- `engine`: Template rendering engine

**Variable Definitions**:
- Required and optional variables
- Type validation and constraints
- Default values and descriptions
- Validation rules and patterns

**Context Data Structure**:
- Facts, variables, machines, environment, project data
- Flexible data binding capabilities
- Cross-reference support

**Function Definitions**:
- Allowed function categories
- Security restrictions and patterns
- Performance limits and constraints
- Engine-specific capabilities

**Metadata Structure**:
- Basic metadata (name, description, version, author)
- Classification metadata (tags, categories, output format)
- Lifecycle metadata (timestamps, states)
- Documentation metadata (usage, examples, API reference)

**Common Validation Rules**:
- Template ID format validation
- Source path validation
- Variable name validation
- Function name validation
- Pattern validation
- Size and performance limits
- Security pattern blocking

**Common Constraints**:
- Scope-specific constraints
- Security constraints by level
- Performance constraints
- Engine-specific constraints

**Common Metadata**:
- Schema versioning
- Author and license information
- Dependencies and tags
- Description and documentation

### **Phase 2: Enhanced Context Schema**

#### **`template-context.hcl`**
Enhanced context schema with:

**Context-Specific Features**:
- Data binding capabilities (facts, variables, machines, environment, project)
- Context resolution features (lazy loading, caching, validation, transformation)
- Context scoping features (global, project, machine, template)

**Context-Specific Validation**:
- Facts validation (required facts, format, schema)
- Variables validation (required variables, format, schema)
- Machines validation (required machines, format, schema)
- Environment validation (required environment, format)
- Project validation (required project, format, schema)
- Context data validation (integrity, freshness, size limits)

**Context-Specific Constraints**:
- Performance constraints (size, resolution time, cache)
- Security constraints (filtering, encryption, access control, audit)
- Data constraints (count limits, retention)
- Scope constraints (limits per scope)

**Context-Specific Extensions**:
- Data transformation extensions (facts, variables, machines, environment)
- Context composition extensions (merge strategies, inheritance)
- Context validation extensions (schema, data, cross-reference)
- Context monitoring extensions (metrics, logging, alerting)

**Context-Specific Metadata**:
- Context versioning
- Context lifecycle
- Context ownership
- Context tags and description

### **Phase 3: Enhanced Functions Schema**

#### **`template-functions.hcl`**
Enhanced functions schema with:

**Function-Specific Features**:
- Function categories (data, project, string, math, array, system)
- Run features (lazy evaluation, caching, parallel running, error handling)
- Security features (sandboxing, access control, audit logging, pattern filtering)

**Function-Specific Validation**:
- Function name validation (pattern, reserved names, length)
- Function argument validation (count, types, required arguments)
- Function return validation (types, size, value validation)
- Function pattern validation (syntax, dangerous patterns)
- Function performance validation (run time, memory, CPU, I/O)
- Function security validation (file, network, process, environment access)

**Function-Specific Constraints**:
- Performance constraints (concurrent functions, total time, memory, cache)
- Security constraints (sandboxing, access control, audit)
- Function category constraints (data, string, math, array, system)
- Engine-specific constraints (Go template, Jinja2, Handlebars)

**Function-Specific Extensions**:
- Custom function extensions (registration, plugins)
- Function optimization extensions (inlining, memoization, parallelization)
- Function monitoring extensions (metrics, profiling, tracing)
- Function debugging extensions (debugging, testing)

**Function-Specific Metadata**:
- Function versioning
- Function lifecycle
- Function documentation
- Function tags

### **Phase 4: Enhanced Metadata Schema**

#### **`template-metadata.hcl`**
Enhanced metadata schema with:

**Metadata-Specific Features**:
- Metadata categories (basic, classification, lifecycle, dependency, documentation)
- Metadata management features (validation, inheritance, merging, versioning)
- Metadata discovery features (indexing, search, filtering, sorting)
- Metadata export features (formats, validation, compression, encryption)

**Metadata-Specific Validation**:
- Basic metadata validation (name, description, version, author)
- Classification metadata validation (tags, categories, output format)
- Lifecycle metadata validation (timestamps, states)
- Dependency metadata validation (dependencies, requirements, conflicts)
- Documentation metadata validation (usage, examples, API reference)

**Metadata-Specific Constraints**:
- Size constraints (metadata size, field lengths, counts)
- Performance constraints (validation timeout, processing timeout, cache)
- Security constraints (encryption, access control, audit logging, filtering)
- Content constraints (allowed characters, forbidden patterns, sanitization)
- Versioning constraints (scheme, compatibility, migration, deprecation)

**Metadata-Specific Extensions**:
- Metadata transformation extensions (format, content, validation)
- Metadata indexing extensions (full-text, categorization, analytics)
- Metadata discovery extensions (search, filtering, sorting)
- Metadata export extensions (format, content, delivery)

**Metadata-Specific Metadata**:
- Metadata versioning
- Metadata lifecycle
- Metadata ownership
- Metadata tags and description

## Template-Specific Features

### **Context Schema Features**

**Data Binding Capabilities**:
- Facts binding: Support for binding machine facts to templates
- Variables binding: Support for binding project variables to templates
- Machines binding: Support for binding machine inventory to templates
- Environment binding: Support for binding environment variables to templates
- Project binding: Support for binding project information to templates

**Context Resolution Features**:
- Lazy loading: Support for lazy loading of context data
- Caching: Support for caching resolved context data
- Validation: Support for validating context data
- Transformation: Support for transforming context data

**Context Scoping Features**:
- Global scope: Support for global context scope
- Project scope: Support for project context scope
- Machine scope: Support for machine context scope
- Template scope: Support for template-specific context scope

### **Functions Schema Features**

**Function Categories**:
- Data functions: Access template variables, facts, environment, external data
- Project functions: Access project information, machine inventory
- String functions: String manipulation (upper, lower, trim, join, split, replace)
- Math functions: Mathematical operations (add, sub, mul, div)
- Array functions: Array manipulation (length, index, first, last, sort, reverse)
- System functions: System-related operations (restricted by security level)

**Security Features**:
- Sandboxing: Function run sandboxing
- Access control: Function permission levels
- Audit logging: Function call logging
- Pattern filtering: Dangerous pattern blocking

**Performance Features**:
- Caching: Function result caching
- Optimization: Function inlining and memoization
- Monitoring: Function metrics and profiling
- Debugging: Function debugging and testing

### **Metadata Schema Features**

**Metadata Categories**:
- Basic metadata: Name, description, version, author
- Classification metadata: Tags, categories, output format, template type
- Lifecycle metadata: Creation, update, deprecation, removal timestamps
- Dependency metadata: Dependencies, requirements, conflicts, provides
- Documentation metadata: Usage, examples, API reference, changelog

**Management Features**:
- Validation: Metadata validation against schemas
- Inheritance: Metadata inheritance from parent schemas
- Merging: Metadata merging strategies
- Versioning: Metadata versioning and compatibility

**Discovery Features**:
- Indexing: Full-text indexing of metadata
- Search: Advanced search capabilities
- Filtering: Multi-criteria filtering
- Sorting: Flexible sorting options

## Usage Examples

### **Basic Template Context Usage**
```hcl
template_context {
  include = "template-structure"
  
  scope = "context"
  template_type = "context"
  security_level = "standard"
  engine = "go-template"
  
  context_features = {
    data_binding = {
      facts_binding = true
      variables_binding = true
      machines_binding = true
    }
  }
  
  context_validation = {
    facts_validation = {
      required_facts = ["os", "hardware", "network"]
      facts_format = "json"
    }
  }
}
```

### **Advanced Template Functions Usage**
```hcl
template_functions {
  include = "template-structure"
  
  scope = "functions"
  template_type = "functions"
  security_level = "elevated"
  engine = "jinja2"
  
  function_features = {
    function_categories = {
      data_functions = {
        var = true
        facts = true
        env = true
      }
      string_functions = {
        upper = true
        lower = true
        trim = true
      }
    }
  }
  
  function_validation = {
    function_name_validation = {
      name_pattern = "^[a-zA-Z_][a-zA-Z0-9_]*$"
      max_name_length = 50
    }
  }
  
  function_constraints = {
    performance_constraints = {
      max_concurrent_functions = 20
      max_total_run_time = 30000
    }
  }
}
```

### **Comprehensive Template Metadata Usage**
```hcl
template_metadata {
  include = "template-structure"
  
  scope = "metadata"
  template_type = "metadata"
  security_level = "standard"
  engine = "go-template"
  
  metadata_features = {
    metadata_categories = {
      basic_metadata = {
        name = true
        description = true
        version = true
        author = true
      }
      classification_metadata = {
        tags = true
        categories = true
        output_format = true
      }
    }
  }
  
  metadata_validation = {
    basic_metadata_validation = {
      name_validation = {
        name_pattern = "^[a-zA-Z0-9._-]+$"
        name_length = {
          min = 1
          max = 100
        }
      }
    }
  }
  
  metadata_extensions = {
    metadata_transformation_extensions = {
      format_transformation = {
        json_to_yaml = true
        yaml_to_json = true
      }
    }
  }
}
```

## Migration Guide

### **From Simple Schemas to Enhanced Schemas**

#### **Step 1: Update Schema Structure**
```hcl
# Before (Simple)
template_context {
  facts = {
    type = "object"
    description = "Machine facts available to templates"
    additional_properties = true
  }
}

# After (Enhanced)
template_context {
  include = "template-structure"
  
  scope = "context"
  template_type = "context"
  security_level = "standard"
  engine = "go-template"
  
  context_features = {
    data_binding = {
      facts_binding = true
    }
  }
  
  context_validation = {
    facts_validation = {
      facts_format = "json"
    }
  }
}
```

#### **Step 2: Add Enhanced Features**
- Include base structure reference
- Add schema-specific features
- Implement comprehensive validation
- Define performance and security constraints
- Add monitoring and debugging extensions

#### **Step 3: Update Validation**
- Replace simple validation with comprehensive rules
- Add cross-reference validation
- Implement security pattern validation
- Add performance and resource validation

#### **Step 4: Add Extensions**
- Implement data transformation capabilities
- Add monitoring and observability features
- Include debugging and testing support
- Add export and integration capabilities

### **Backward Compatibility**
- Existing template schemas continue to work
- Gradual migration to enhanced schemas
- Automatic validation and constraint application
- Enhanced functionality for new templates

## Future Enhancements

### **Planned Features**
1. **Schema Composition Tools**: Automated schema composition and validation
2. **Template Schema Registry**: Centralized template schema management
3. **Advanced Validation**: Machine learning-based validation and suggestions
4. **Performance Optimization**: Advanced caching and optimization strategies
5. **Security Enhancements**: Advanced security patterns and threat detection

### **Integration Opportunities**
1. **CI/CD Integration**: Automated schema validation in pipelines
2. **IDE Support**: Enhanced editor support for template schemas
3. **API Integration**: RESTful API for schema management
4. **Plugin System**: Extensible plugin architecture for custom features

## Conclusion

The Enhanced Composition Pattern for template schemas provides:

- **Robust Architecture**: Solid foundation for template schema management
- **Consistent Patterns**: Unified approach across all template schemas
- **Extensible Design**: Easy to add new features and capabilities
- **Maintainable Code**: Clear separation and centralized management
- **Future-Proof**: Ready for advanced features and integrations

This implementation establishes a comprehensive template schema system that matches the quality and sophistication of the enhanced facts schema system, providing a consistent and powerful foundation for template management in the Spooky project. 