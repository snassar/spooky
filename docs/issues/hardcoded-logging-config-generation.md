# 3
# Convert createDefaultLoggingConfig to schema-driven generation

## Problem Statement
The `createDefaultLoggingConfig` function in `internal/config/auto_setup.go` uses hardcoded strings instead of schema-driven generation, which creates inconsistency and maintenance issues.

### Current Implementation
```go
func createDefaultLoggingConfig(configDir string) error {
	content := `# Global logging configuration for spooky
# This file configures logging behavior for all spooky operations
# 
# Default behavior (when this file doesn't exist or is empty):
# - Log level: error (only errors are shown)
# - Format: json (structured logging)
# - Output: null (no logging output to terminal)

# Uncomment and modify the settings below to customize logging behavior:

# logging {
#   # Log level (debug, info, warn, error, fatal)
#   # Default: error (only show errors)
#   level = "error"
#   
#   # Output format (json, text, structured)
#   # Default: json (structured logging)
#   format = "json"
#   
#   # Output destination (stdout, stderr, file, null)
#   # Default: null (no output to terminal)
#   # Options: stdout, stderr, file, null
#   output = "null"
#   
#   # ... (more hardcoded configuration)
# }
`
	configPath := filepath.Join(configDir, "logging.hcl")
	return os.WriteFile(configPath, []byte(content), 0o600)
}
```

### Issues
1. **Hardcoded content** - No schema validation or consistency
2. **Manual maintenance** - Changes to logging schema require manual updates
3. **Potential validation errors** - Generated content may not match current schema
4. **Inconsistent architecture** - Doesn't follow schema-driven patterns used elsewhere
5. **Commented-out configuration** - Generated file is mostly comments, not functional configuration

### Proposed Solution
Convert to schema-driven generation similar to the fixed `createProjectHCL`:

1. **Load logging schema** from `internal/schemas/schemas/structure/logging.hcl`
2. **Generate content** based on schema structure and validation rules
3. **Include proper examples** that match the schema requirements
4. **Ensure validation compliance** with logging schema

### Benefits
- **Consistency** with schema definitions
- **Automatic updates** when schema changes
- **Validation compliance** guaranteed
- **Architectural alignment** with schema-driven approach
- **Better user experience** with functional default configuration

## Priority
**Medium Priority** - Global configuration but not core functionality

## Related Files
- `internal/config/auto_setup.go` - Current implementation
- `internal/schemas/schemas/structure/logging.hcl` - Schema definition
- `internal/logging/` - Logging system implementation

## Acceptance Criteria
- [ ] Replace hardcoded template with schema-driven generation
- [ ] Load logging schema from file
- [ ] Generate valid logging.hcl content based on schema
- [ ] Include proper examples that match schema requirements
- [ ] Ensure generated file passes validation
- [ ] Add comprehensive tests for schema-driven generation
- [ ] Update documentation to reflect schema-driven approach
- [ ] Generate functional default configuration instead of commented examples
- [ ] Maintain backward compatibility with existing configuration files

## Implementation Notes
- Follow the pattern established by the fixed `createProjectHCL` function
- Use the schema manager to load and validate the logging schema
- Generate examples that demonstrate proper logging configuration
- Ensure the generated content is immediately valid and usable
- Consider generating a functional default configuration rather than commented examples
- Include examples for different logging scenarios (development, production, debugging, etc.)
