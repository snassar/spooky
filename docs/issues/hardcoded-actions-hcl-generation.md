# Convert createActionsHCL to schema-driven generation

## Problem Statement
The `createActionsHCL` function in `internal/project/manager.go` uses hardcoded strings instead of schema-driven generation, which creates inconsistency and maintenance issues.



### Current Implementation
```go
func (m *Manager) createActionsHCL(filePath string) error {
	content := `# Actions for spooky project
# Define your actions here

actions {
  # Example action definition
  # action "deploy-web" {
  #   description = "Deploy web application"
  #   
  #   machines = ["web-server"]
  #   parallel = true
  #   
  #   template {
  #     source = "templates/deploy.sh.tmpl"
  #     destination = "/tmp/deploy.sh"
  #     permissions = "0755"
  #   }
  #   
  #   command = "/tmp/deploy.sh"
  # }
}
`
	return os.WriteFile(filePath, []byte(content), 0o600)
}
```

### Issues
1. **Hardcoded content** - No schema validation or consistency
2. **Manual maintenance** - Changes to actions schema require manual updates
3. **Potential validation errors** - Generated content may not match current schema
4. **Inconsistent architecture** - Doesn't follow schema-driven patterns used elsewhere
5. **Outdated examples** - May not reflect current action schema structure

### Proposed Solution
Convert to schema-driven generation similar to the fixed `createProjectHCL`:

1. **Load actions schema** from `internal/schemas/schemas/structure/actions.hcl`
2. **Generate content** based on schema structure and validation rules
3. **Include proper examples** that match the schema requirements
4. **Ensure validation compliance** with actions schema

### Benefits
- **Consistency** with schema definitions
- **Automatic updates** when schema changes
- **Validation compliance** guaranteed
- **Architectural alignment** with schema-driven approach
- **Up-to-date examples** that reflect current schema

## Priority
**High Priority** - Core functionality for project initialization

## Related Files
- `internal/project/manager.go` - Current implementation
- `internal/schemas/schemas/structure/actions.hcl` - Schema definition
- `internal/actions/` - Actions system implementation

## Acceptance Criteria
- [ ] Replace hardcoded template with schema-driven generation
- [ ] Load actions schema from file
- [ ] Generate valid actions.hcl content based on schema
- [ ] Include proper examples that match schema requirements
- [ ] Ensure generated file passes validation
- [ ] Add comprehensive tests for schema-driven generation
- [ ] Update documentation to reflect schema-driven approach
- [ ] Verify examples reflect current action schema structure

## Implementation Notes
- Follow the pattern established by the fixed `createProjectHCL` function
- Use the schema manager to load and validate the actions schema
- Generate examples that demonstrate proper action configuration
- Ensure the generated content is immediately valid and usable
- Include examples for different action types (template-based, command-based, etc.)
