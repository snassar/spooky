# Convert createDefaultSpookyConfig to schema-driven generation

## Problem Statement
The `createDefaultSpookyConfig` function in `internal/config/auto_setup.go` uses hardcoded strings instead of schema-driven generation, which creates inconsistency and maintenance issues.

### Current Implementation
```go
func createDefaultSpookyConfig(configDir string) error {
	content := `# Spooky CLI Configuration
# This file contains global configuration for the spooky CLI tool

# Global CLI settings
cli {
  # Default timeout for operations (in seconds)
  default_timeout = 300
  
  # Maximum parallel operations
  max_parallel = 10
  
  # Default log level for CLI operations
  log_level = "info"
  
  # Enable colored output (if supported by terminal)
  colored_output = true
  
  # Show progress indicators for long-running operations
  show_progress = true
}

# SSH configuration
ssh {
  # Default SSH timeout (in seconds)
  timeout = 30
  
  # SSH connection retry attempts
  retry_attempts = 3
  
  # Delay between retry attempts (in seconds)
  retry_delay = 5
  
  # Enable SSH connection pooling
  connection_pooling = true
  
  # Maximum number of SSH connections to keep in pool
  max_connections = 20
}

# ... (more hardcoded configuration)
`
	configPath := filepath.Join(configDir, "spooky.hcl")
	return os.WriteFile(configPath, []byte(content), 0o600)
}
```

### Issues
1. **Hardcoded content** - No schema validation or consistency
2. **Manual maintenance** - Changes to spooky config schema require manual updates
3. **Potential validation errors** - Generated content may not match current schema
4. **Inconsistent architecture** - Doesn't follow schema-driven patterns used elsewhere
5. **Large hardcoded template** - Difficult to maintain and update

### Proposed Solution
Convert to schema-driven generation similar to the fixed `createProjectHCL`:

1. **Load spooky config schema** from `internal/schemas/schemas/structure/spooky.hcl`
2. **Generate content** based on schema structure and validation rules
3. **Include proper examples** that match the schema requirements
4. **Ensure validation compliance** with spooky config schema

### Benefits
- **Consistency** with schema definitions
- **Automatic updates** when schema changes
- **Validation compliance** guaranteed
- **Architectural alignment** with schema-driven approach
- **Easier maintenance** of configuration templates

## Priority
**Medium Priority** - Global configuration but not core functionality

## Related Files
- `internal/config/auto_setup.go` - Current implementation
- `internal/schemas/schemas/structure/spooky.hcl` - Schema definition
- `internal/config/` - Config system implementation

## Acceptance Criteria
- [ ] Replace hardcoded template with schema-driven generation
- [ ] Load spooky config schema from file
- [ ] Generate valid spooky.hcl content based on schema
- [ ] Include proper examples that match schema requirements
- [ ] Ensure generated file passes validation
- [ ] Add comprehensive tests for schema-driven generation
- [ ] Update documentation to reflect schema-driven approach
- [ ] Maintain all current configuration options in generated content

### Implementation Notes
- Follow the pattern established by the fixed `createProjectHCL` function
- Use the schema manager to load and validate the spooky config schema
- Generate examples that demonstrate proper configuration
- Ensure the generated content is immediately valid and usable
- Consider the complexity of the spooky config schema and ensure all sections are properly generated
- Maintain backward compatibility with existing configuration files
