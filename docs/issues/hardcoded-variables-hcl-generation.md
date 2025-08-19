# Convert createVariablesHCL to schema-driven generation

## Problem Statement
The `createVariablesHCL` function in `internal/project/manager.go` uses hardcoded strings instead of schema-driven generation, which creates inconsistency and maintenance issues.

### Current Implementation
```go
func (m *Manager) createVariablesHCL(filePath string) error {
	content := `# Variables for spooky project
# Define your variables here

variables {
  # Example variable definitions
  # app_version = "1.0.0"
  # environment = "production"
  # backup_retention_days = 30
}
`
	return os.WriteFile(filePath, []byte(content), 0o600)
}
```

### Issues
1. **Hardcoded content** - No schema validation or consistency
2. **Manual maintenance** - Changes to variables schema require manual updates
3. **Potential validation errors** - Generated content may not match current schema
4. **Inconsistent architecture** - Doesn't follow schema-driven patterns used elsewhere
5. **Limited examples** - Basic examples don't demonstrate full variable capabilities

### Proposed Solution
Convert to schema-driven generation similar to the fixed `createProjectHCL`:

1. **Load variables schema** from `internal/schemas/schemas/structure/variables.hcl`
2. **Generate content** based on schema structure and validation rules
3. **Include comprehensive examples** that match the schema requirements
4. **Ensure validation compliance** with variables schema

### Benefits
- **Consistency** with schema definitions
- **Automatic updates** when schema changes
- **Validation compliance** guaranteed
- **Architectural alignment** with schema-driven approach
- **Better examples** that demonstrate variable capabilities

## Priority
**Medium Priority** - Important for project initialization but not core functionality

## Related Files
- `internal/project/manager.go` - Current implementation
- `internal/schemas/schemas/structure/variables.hcl` - Schema definition
- `internal/variables/` - Variables system implementation

## Acceptance Criteria
- [ ] Replace hardcoded template with schema-driven generation
- [ ] Load variables schema from file
- [ ] Generate valid variables.hcl content based on schema
- [ ] Include comprehensive examples that match schema requirements
- [ ] Ensure generated file passes validation
- [ ] Add comprehensive tests for schema-driven generation
- [ ] Update documentation to reflect schema-driven approach
- [ ] Include examples for different variable types and use cases

### Implementation Notes
- Follow the pattern established by the fixed `createProjectHCL` function
- Use the schema manager to load and validate the variables schema
- Generate examples that demonstrate proper variable configuration
- Ensure the generated content is immediately valid and usable
- Include examples for different variable types (strings, numbers, booleans, etc.)
- Consider including examples for variable dependencies and constraints
