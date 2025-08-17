# Code Walkthroughs

This directory contains detailed function-by-function walkthroughs of all spooky validate commands. Each walkthrough provides a comprehensive understanding of the execution flow, from entry point to completion.

## Available Walkthroughs

### Project Validation
- **[Project Validate](project-validate.md)** - `spooky project validate <project-path>`
  - Validates project structure and configuration
  - Checks project.hcl syntax and content
  - Validates project directory structure

### Actions Validation
- **[Actions Validate](actions-validate.md)** - `spooky actions validate <project-path>`
  - Validates action configurations for syntax and dependencies
  - Checks action definitions and templates
  - Validates action relationships and references

### Machines Validation
- **[Machines Validate](machines-validate.md)** - `spooky machines validate <project-path>`
  - Validates machine inventory configurations
  - Checks SSH connectivity and authentication
  - Validates machine definitions and settings

### Templates Validation
- **[Templates Validate](templates-validate.md)** - `spooky templates validate <project-path>`
  - Validates template syntax and variables
  - Checks template rendering and data integration
  - Validates template functions and helpers

### Variables Validation
- **[Variables Validate](variables-validate.md)** - `spooky variables validate <project-path>`
  - Validates variable definitions and dependencies
  - Checks variable types and constraints
  - Validates variable resolution and relationships

### Secrets Validation
- **[Secrets Validate](secrets-validate.md)** - `spooky secrets validate <project-path>`
  - Validates age configuration and keys
  - Checks identity files and permissions
  - Validates encrypted values and recipients

### Configuration Validation
- **[Config Validate](config-validate.md)** - `spooky config validate [--config <file>]`
  - Validates global spooky configuration
  - Checks logging, SSH, and facts settings
  - Validates age encryption configuration

## Common Patterns

All validate commands follow these common architectural patterns:

### 1. Entry Point Flow
```
main() → cmd.Execute() → RootCmd.PersistentPreRunE → Command Routing → Handler
```

### 2. Dependency Injection
- **Integration Layer**: Provides CLI access to domain functionality
- **Manager Layer**: Coordinates operations and business logic
- **Validator Layer**: Performs validation and error checking
- **Loader Layer**: Handles file loading and parsing

### 3. Validation Process
1. **Initialization**: Set up dependencies and logging
2. **Loading**: Load configuration files and data
3. **Validation**: Perform comprehensive validation
4. **Reporting**: Display detailed results to user

### 4. Error Handling
- **Structured Errors**: Use consistent error types and codes
- **Context Information**: Provide detailed error context
- **Warning System**: Distinguish between errors and warnings
- **Aggregation**: Collect all issues before reporting

### 5. Output Format
- **Consistent UI**: Use emojis and formatting for readability
- **Detailed Results**: Show specific errors and warnings
- **Summary**: Provide clear success/failure indication
- **Context**: Include relevant file paths and line numbers

## Architecture Components

### Integration Layer
- **Purpose**: Bridge between CLI and domain logic
- **Responsibilities**: Command handling, dependency management
- **Pattern**: Interface-based delegation to managers

### Manager Layer
- **Purpose**: Coordinate domain operations
- **Responsibilities**: Business logic, workflow coordination
- **Pattern**: Dependency injection with specialized components

### Validator Layer
- **Purpose**: Validate data and configurations
- **Responsibilities**: Syntax checking, constraint validation
- **Pattern**: Comprehensive validation with detailed reporting

### Loader Layer
- **Purpose**: Load data from various sources
- **Responsibilities**: File parsing, HCL processing, data extraction
- **Pattern**: Multiple source support with error handling

## Key Design Principles

### 1. Interface-Based Design
- All components use interfaces for loose coupling
- Dependency injection enables testing and flexibility
- Clear separation of concerns across layers

### 2. Comprehensive Validation
- Validate all aspects of configuration and data
- Provide detailed error messages with context
- Support both strict and lenient validation modes

### 3. User-Friendly Output
- Clear, readable error messages
- Consistent formatting and emoji usage
- Detailed context for troubleshooting

### 4. Error Aggregation
- Collect all validation issues before reporting
- Distinguish between errors and warnings
- Provide actionable feedback for fixes

### 5. Extensible Architecture
- Easy to add new validation rules
- Support for custom validation logic
- Pluggable validation components

## Usage Examples

### Basic Validation
```bash
# Validate project structure
spooky project validate ./my-project

# Validate actions configuration
spooky actions validate ./my-project

# Validate machine inventory
spooky machines validate ./my-project
```

### Advanced Validation
```bash
# Validate specific templates
spooky templates validate ./my-project --template templates/deploy.sh.tmpl

# Validate with custom config
spooky config validate --config /path/to/spooky.hcl

# Validate secrets with specific identity
spooky secrets validate ./my-project
```

## Troubleshooting

### Common Issues
1. **Missing Dependencies**: Ensure all required components are initialized
2. **File Permissions**: Check file and directory permissions
3. **HCL Syntax**: Validate HCL syntax in configuration files
4. **Path Resolution**: Verify file paths and environment variables

### Debugging Tips
1. **Enable Debug Logging**: Set log level to debug for detailed output
2. **Check File Paths**: Verify all referenced files exist
3. **Validate HCL**: Use HCL syntax checkers for configuration files
4. **Review Error Context**: Pay attention to error context information

## Contributing

When adding new validate commands or modifying existing ones:

1. **Follow Patterns**: Use established architectural patterns
2. **Add Walkthrough**: Create comprehensive walkthrough documentation
3. **Update Index**: Add new walkthrough to this README
4. **Test Thoroughly**: Ensure all validation scenarios are covered
5. **Document Changes**: Update relevant documentation

## Related Documentation

- **[Interface Architecture](../rules/interface-architecture.mdc)** - Interface-based design patterns
- **[Code Quality Standards](../rules/code-quality-standards.mdc)** - Quality requirements
- **[Testing Standards](../rules/testing.mdc)** - Testing strategies
- **[Error Handling Standards](../rules/error-handling-standards.mdc)** - Error patterns
