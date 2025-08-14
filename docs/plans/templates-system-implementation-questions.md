# Templates System Implementation Questions

## Overview

This document contains all the questions that need to be resolved before starting the templates system implementation. Each section contains specific questions with structured answer formats that will help guide the implementation.

**Purpose**: Resolve all implementation details and requirements before beginning development to ensure a complete and functional templates system.

**Format**: For each question, provide a clear, specific answer that an AI can use to implement the feature correctly.

---

## 1. Template Function Scope and Security

### 1.1 Template Function Whitelist
**Question**: What specific template functions should be allowed in the whitelist?

**Answer Format**:
```
ALLOWED_FUNCTIONS:
- function_name: "description of what it does"
- function_name: "description of what it does"

EXAMPLES:
- custom: "Access custom facts data with path-based lookup"
- system: "Access system facts data with path-based lookup"
- env: "Access environment variables"
- var: "Access project variables"
- varOrDefault: "Access project variables with default fallback"
- machines: "Get all machines from inventory"
- machine: "Get specific machine by name"
- data: "Access additional custom data sources"
```

### 1.2 Security Boundaries
**Question**: What are the security boundaries and restrictions for template functions?

**Answer Format**:
```
SECURITY_RESTRICTIONS:
- file_system_access: "yes/no - what paths allowed"
- network_access: "yes/no - what endpoints allowed"
- execution_timeout: "maximum seconds for template rendering"
- memory_limit: "maximum memory usage in MB"
- function_execution_limit: "maximum function calls per template"

SANDBOXING:
- execution_environment: "description of sandbox"
- resource_limits: "CPU, memory, disk limits"
- isolation_level: "process, thread, or function level"
```

### 1.3 Dangerous Pattern Blacklist
**Question**: What patterns should be blacklisted for security?

**Answer Format**:
```
BLACKLISTED_PATTERNS:
- pattern: "reason for blacklisting"
- pattern: "reason for blacklisting"

EXAMPLES:
- "{{.env.PASSWORD}}": "Prevents environment variable access to sensitive data"
- "{{.var.SECRET_KEY}}": "Prevents variable access to sensitive data"
- "{{.custom.private.*}}": "Prevents access to private custom data"
```

---

## 2. Template Context Data Structure

### 2.1 Context Structure Definition
**Question**: What is the exact structure of the TemplateContext and how does it integrate with existing systems?

**Answer Format**:
```
TEMPLATE_CONTEXT_STRUCTURE:
  Project:
    type: "spookytypes.ProjectConfig"
    description: "Project configuration and metadata"
    access_pattern: "{{.Project.Name}}"
  
  Facts:
    type: "map[string]interface{}"
    description: "Facts data from facts system"
    access_pattern: "{{.Facts.system.hostname}}"
  
  Machines:
    type: "[]*spookytypes.Machine"
    description: "Machine inventory data"
    access_pattern: "{{range .Machines}}{{.Name}}{{end}}"
  
  Variables:
    type: "spookytypes.VariableContext"
    description: "Project variables"
    access_pattern: "{{.Variables.project_name}}"
  
  Environment:
    type: "map[string]string"
    description: "Environment variables"
    access_pattern: "{{.Environment.HOME}}"
  
  CustomData:
    type: "map[string]interface{}"
    description: "Additional custom data sources"
    access_pattern: "{{.CustomData.user_data}}"
```

### 2.2 Facts Integration
**Question**: How should facts data be structured and accessed in templates?

**Answer Format**:
```
FACTS_INTEGRATION:
  data_structure:
    system_facts: "description of system facts structure"
    custom_facts: "description of custom facts structure"
    machine_facts: "description of machine-specific facts"
  
  access_patterns:
    system_facts: "{{.Facts.system.os}}"
    custom_facts: "{{.Facts.custom.application_version}}"
    machine_facts: "{{.Facts.machines.web_server.hostname}}"
  
  caching_strategy:
    cache_ttl: "time in seconds"
    cache_invalidation: "when to invalidate cache"
```

### 2.3 Variables Integration
**Question**: How should variables be resolved and prioritized in templates?

**Answer Format**:
```
VARIABLES_INTEGRATION:
  resolution_order:
    1: "highest priority source"
    2: "second priority source"
    3: "third priority source"
  
  access_patterns:
    simple_variable: "{{.Variables.project_name}}"
    nested_variable: "{{.Variables.config.database.host}}"
    default_value: "{{.Variables.port | default 8080}}"
  
  validation:
    required_variables: "how to handle missing required variables"
    type_validation: "how to validate variable types"
```

### 2.4 Machine Data Integration
**Question**: How should machine data be filtered and presented in templates?

**Answer Format**:
```
MACHINE_INTEGRATION:
  filtering_options:
    by_tags: "{{range (machinesByTag .Machines \"web\")}}{{.Name}}{{end}}"
    by_pattern: "{{range (machinesByPattern .Machines \"web-*\")}}{{.Name}}{{end}}"
    by_environment: "{{range (machinesByEnv .Machines \"production\")}}{{.Name}}{{end}}"
  
  machine_data_structure:
    basic_fields: "Name, Hostname, Port, User"
    authentication: "SSH key path, password, etc."
    metadata: "Tags, environment, etc."
  
  access_patterns:
    all_machines: "{{range .Machines}}{{.Name}}{{end}}"
    specific_machine: "{{(machine .Machines \"web-server\").Hostname}}"
```

---

## 3. CLI Command Scope and Flags

### 3.1 Command Structure
**Question**: What are the exact CLI commands and flags needed?

**Answer Format**:
```
CLI_COMMANDS:
  templates list:
    syntax: "spooky templates list <project-path> [flags]"
    flags:
      --verbose: "Show detailed template information"
      --pattern: "Filter templates by pattern"
      --format: "Output format (table, json, hcl)"
    description: "List all templates in project"
  
  templates validate:
    syntax: "spooky templates validate <project-path> [flags]"
    flags:
      --template: "Validate specific template file"
      --verbose: "Show detailed validation errors"
      --strict: "Use strict validation mode"
    description: "Validate template syntax and configuration"
  
  templates render:
    syntax: "spooky templates render <project-path> <template> [flags]"
    flags:
      --data: "Additional data file (JSON/HCL format)"
      --machine: "Target machine name or pattern"
      --output: "Output file path (default: stdout)"
      --preview: "Show preview without writing files"
      --dry-run: "Show what would be rendered"
    description: "Render template with project data"
  
  templates preview:
    syntax: "spooky templates preview <project-path> <template> [flags]"
    flags:
      --data: "Additional data file for preview"
      --mock: "Use mock data for preview"
    description: "Preview template rendering with analysis"
```

### 3.2 Flag Specifications
**Question**: What are the exact specifications for each flag?

**Answer Format**:
```
FLAG_SPECIFICATIONS:
  --data:
    format: "JSON or HCL file path"
    example: "--data config.json"
    description: "Additional data to merge with template context"
  
  --machine:
    format: "machine name, pattern, or tag"
    examples:
      - "--machine web-server"
      - "--machine web-*"
      - "--machine tag:production"
    description: "Target specific machine(s) for rendering"
  
  --output:
    format: "file path or '-' for stdout"
    example: "--output rendered.conf"
    description: "Output file path (default: stdout)"
  
  --preview:
    behavior: "Show rendered output without writing files"
    description: "Preview mode for template rendering"
  
  --dry-run:
    behavior: "Show what would be rendered without actual rendering"
    description: "Dry run mode for template rendering"
```

---

## 4. Template File Organization

### 4.1 Directory Structure
**Question**: How should templates be organized within projects?

**Answer Format**:
```
TEMPLATE_ORGANIZATION:
  directory_structure:
    templates/: "Root templates directory"
    templates/configs/: "Configuration templates"
    templates/scripts/: "Script templates"
    templates/documents/: "Document templates"
    templates/custom/: "Custom project templates"
  
  file_naming:
    pattern: "*.tmpl or *.template"
    examples:
      - "nginx.conf.tmpl"
      - "deploy.sh.template"
      - "README.md.tmpl"
  
  metadata:
    location: "Template metadata location"
    format: "Metadata format (comments, separate file, etc.)"
    fields: "Required metadata fields"
```

### 4.2 Template Metadata
**Question**: What metadata should templates include?

**Answer Format**:
```
TEMPLATE_METADATA:
  required_fields:
    name: "Template name"
    description: "Template description"
    version: "Template version"
  
  optional_fields:
    author: "Template author"
    tags: "Template tags for categorization"
    dependencies: "Required variables or facts"
    examples: "Usage examples"
  
  metadata_format:
    type: "HCL comments, YAML frontmatter, or separate file"
    example: "Show example metadata format"
```

---

## 5. Performance and Caching Strategy

### 5.1 Caching Requirements
**Question**: What are the specific performance requirements and caching strategies?

**Answer Format**:
```
CACHING_STRATEGY:
  template_cache:
    ttl: "Time in seconds"
    invalidation: "When to invalidate"
    storage: "Memory, disk, or both"
  
  context_cache:
    ttl: "Time in seconds"
    invalidation: "When to invalidate"
    scope: "Global or per-project"
  
  facts_cache:
    ttl: "Time in seconds"
    invalidation: "When facts are updated"
    scope: "Global or per-machine"
  
  variables_cache:
    ttl: "Time in seconds"
    invalidation: "When variables are updated"
    scope: "Per-project"
```

### 5.2 Performance Limits
**Question**: What are the performance limits and benchmarks?

**Answer Format**:
```
PERFORMANCE_LIMITS:
  execution_timeout:
    default: "Maximum seconds for template rendering"
    configurable: "Can be overridden in config"
  
  memory_limit:
    default: "Maximum memory usage in MB"
    configurable: "Can be overridden in config"
  
  template_size_limit:
    default: "Maximum template file size in KB"
    configurable: "Can be overridden in config"
  
  benchmarks:
    simple_template: "Expected rendering time"
    complex_template: "Expected rendering time"
    large_dataset: "Expected rendering time with large data"
```

---

## 6. Error Handling and Validation Scope

### 6.1 Validation Levels
**Question**: What level of validation and error reporting is required?

**Answer Format**:
```
VALIDATION_LEVELS:
  syntax_validation:
    scope: "What to validate"
    error_format: "Error message format"
    examples: "Example error messages"
  
  semantic_validation:
    scope: "What to validate"
    error_format: "Error message format"
    examples: "Example error messages"
  
  security_validation:
    scope: "What to validate"
    error_format: "Error message format"
    examples: "Example error messages"
  
  strict_mode:
    enabled_by: "How to enable strict mode"
    additional_checks: "What additional checks are performed"
```

### 6.2 Error Message Format
**Question**: What format should error messages follow?

**Answer Format**:
```
ERROR_FORMAT:
  structure:
    type: "Error type (syntax, semantic, security)"
    message: "Human-readable error message"
    line: "Line number in template"
    column: "Column number in template"
    context: "Additional context information"
  
  examples:
    syntax_error: "Show example syntax error"
    semantic_error: "Show example semantic error"
    security_error: "Show example security error"
  
  localization: "Support for multiple languages (yes/no)"
```

---

## 7. Integration Dependencies

### 7.1 Facts System Integration
**Question**: What are the exact dependencies on the facts system?

**Answer Format**:
```
FACTS_INTEGRATION:
  data_access:
    methods: "How to access facts data"
    caching: "How facts are cached"
    updates: "How facts updates are handled"
  
  data_formats:
    system_facts: "Format of system facts"
    custom_facts: "Format of custom facts"
    machine_facts: "Format of machine facts"
  
  error_handling:
    missing_facts: "How to handle missing facts"
    invalid_facts: "How to handle invalid facts"
```

### 7.2 Variables System Integration
**Question**: What are the exact dependencies on the variables system?

**Answer Format**:
```
VARIABLES_INTEGRATION:
  resolution:
    order: "Variable resolution order"
    precedence: "Variable precedence rules"
    scoping: "Variable scoping rules"
  
  types:
    supported_types: "What variable types are supported"
    type_conversion: "How type conversion works"
    validation: "How variables are validated"
  
  error_handling:
    missing_variables: "How to handle missing variables"
    invalid_variables: "How to handle invalid variables"
```

### 7.3 Machines System Integration
**Question**: What are the exact dependencies on the machines system?

**Answer Format**:
```
MACHINES_INTEGRATION:
  data_access:
    methods: "How to access machine data"
    filtering: "How machine filtering works"
    sorting: "How machine sorting works"
  
  authentication:
    ssh_keys: "How SSH keys are handled"
    credentials: "How credentials are managed"
    security: "Security considerations"
  
  error_handling:
    missing_machines: "How to handle missing machines"
    connection_errors: "How to handle connection errors"
```

---

## 8. Template Security Model

### 8.1 Execution Sandboxing
**Question**: What is the security model for template execution?

**Answer Format**:
```
SECURITY_MODEL:
  sandboxing:
    level: "Process, thread, or function level"
    isolation: "How isolation is implemented"
    resources: "What resources are available"
  
  access_control:
    file_system: "What file system access is allowed"
    network: "What network access is allowed"
    system: "What system calls are allowed"
  
  resource_limits:
    cpu: "CPU usage limits"
    memory: "Memory usage limits"
    disk: "Disk usage limits"
    time: "Execution time limits"
```

### 8.2 Function Security
**Question**: How are template functions secured?

**Answer Format**:
```
FUNCTION_SECURITY:
  whitelist_approach:
    allowed_functions: "List of allowed functions"
    denied_functions: "List of denied functions"
    custom_functions: "How custom functions are handled"
  
  execution_limits:
    max_calls: "Maximum function calls per template"
    max_depth: "Maximum call depth"
    timeout: "Function execution timeout"
  
  input_validation:
    parameter_validation: "How function parameters are validated"
    return_validation: "How function returns are validated"
```

---

## 9. Testing Strategy and Coverage

### 9.1 Unit Testing
**Question**: What are the specific testing requirements?

**Answer Format**:
```
UNIT_TESTING:
  coverage_requirements:
    overall: "Minimum overall coverage percentage"
    per_component: "Minimum coverage per component"
    critical_paths: "100% coverage for critical paths"
  
  test_categories:
    template_parsing: "Tests for template parsing"
    template_rendering: "Tests for template rendering"
    function_evaluation: "Tests for function evaluation"
    error_handling: "Tests for error handling"
    security: "Tests for security features"
  
  test_data:
    mock_data: "What mock data is needed"
    test_templates: "What test templates are needed"
    edge_cases: "What edge cases to test"
```

### 9.2 Integration Testing
**Question**: What integration testing is required?

**Answer Format**:
```
INTEGRATION_TESTING:
  system_integration:
    facts_integration: "Tests for facts system integration"
    variables_integration: "Tests for variables system integration"
    machines_integration: "Tests for machines system integration"
  
  cli_integration:
    command_execution: "Tests for CLI command execution"
    flag_handling: "Tests for flag handling"
    error_reporting: "Tests for error reporting"
  
  performance_testing:
    benchmarks: "Performance benchmarks to meet"
    load_testing: "Load testing requirements"
    stress_testing: "Stress testing requirements"
```

---

## 10. Backward Compatibility

### 10.1 Migration Strategy
**Question**: How should we handle existing template usage?

**Answer Format**:
```
MIGRATION_STRATEGY:
  current_usage:
    string_replacement: "How current string replacement works"
    limitations: "Current limitations"
    examples: "Examples of current usage"
  
  migration_path:
    automatic_migration: "Can migration be automatic (yes/no)"
    manual_migration: "What manual steps are required"
    tools: "What migration tools are needed"
  
  backward_compatibility:
    duration: "How long to maintain backward compatibility"
    deprecation_timeline: "When to deprecate old syntax"
    support: "What support is provided during transition"
```

### 10.2 Compatibility Features
**Question**: What compatibility features should be provided?

**Answer Format**:
```
COMPATIBILITY_FEATURES:
  legacy_support:
    old_syntax: "Support for old template syntax"
    old_functions: "Support for old template functions"
    old_data_access: "Support for old data access patterns"
  
  transition_features:
    warnings: "Warnings for deprecated features"
    suggestions: "Suggestions for new syntax"
    documentation: "Migration documentation"
  
  fallback_behavior:
    missing_functions: "How to handle missing functions"
    invalid_syntax: "How to handle invalid syntax"
    data_mismatches: "How to handle data mismatches"
```

---

## Implementation Priority

### Priority Order
**Question**: What is the implementation priority for these features?

**Answer Format**:
```
IMPLEMENTATION_PRIORITY:
  phase_0_critical:
    - "Feature 1: Description and reason"
    - "Feature 2: Description and reason"
  
  phase_1_core:
    - "Feature 1: Description and reason"
    - "Feature 2: Description and reason"
  
  phase_2_advanced:
    - "Feature 1: Description and reason"
    - "Feature 2: Description and reason"
  
  phase_3_optimization:
    - "Feature 1: Description and reason"
    - "Feature 2: Description and reason"
```

---

## Notes

- **Be specific**: Provide concrete examples and clear specifications
- **Consider constraints**: Think about performance, security, and usability constraints
- **Think ahead**: Consider future extensibility and maintenance
- **Follow patterns**: Align with existing spooky patterns and conventions
- **No placeholders**: All answers should be complete and actionable

**Next Steps**: Once all questions are answered, this document will serve as the complete specification for implementing the templates system.
