# Convert createMachinesHCL to schema-driven generation

## Problem Statement
The `createMachinesHCL` function in `internal/project/manager.go` uses hardcoded strings instead of schema-driven generation, which creates inconsistency and maintenance issues.

### Current Implementation
```go
func (m *Manager) createMachinesHCL(filePath string) error {
	content := `# Machine inventory for spooky project
# Define your machines here

machines {
  # Example machine definition
  # machine "web-server" {
  #   hostname = "web.example.com"
  #   port = 22
  #   user = "admin"
  #   
  #   authentication {
  #     method = "ssh_key"
  #     key_path = "~/.ssh/id_rsa"
  #   }
  #   
  #   tags = ["web", "production"]
  # }
}
`
	return os.WriteFile(filePath, []byte(content), 0o600)
}
```

### Issues
1. **Hardcoded content** - No schema validation or consistency
2. **Manual maintenance** - Changes to machines schema require manual updates
3. **Potential validation errors** - Generated content may not match current schema
4. **Inconsistent architecture** - Doesn't follow schema-driven patterns used elsewhere

### Proposed Solution
Convert to schema-driven generation similar to the fixed `createProjectHCL`:

1. **Load machines schema** from `internal/schemas/schemas/structure/machines.hcl`
2. **Generate content** based on schema structure and validation rules
3. **Include proper examples** that match the schema requirements
4. **Ensure validation compliance** with machines schema

### Benefits
- **Consistency** with schema definitions
- **Automatic updates** when schema changes
- **Validation compliance** guaranteed
- **Architectural alignment** with schema-driven approach

## Priority
**High Priority** - Core functionality for project initialization

## Related Files
- `internal/project/manager.go` - Current implementation
- `internal/schemas/schemas/structure/machines.hcl` - Schema definition
- `internal/machines/` - Machines system implementation

## Acceptance Criteria
- [ ] Replace hardcoded template with schema-driven generation
- [ ] Load machines schema from file
- [ ] Generate valid machines.hcl content based on schema
- [ ] Include proper examples that match schema requirements
- [ ] Ensure generated file passes validation
- [ ] Add comprehensive tests for schema-driven generation
- [ ] Update documentation to reflect schema-driven approach

## Implementation Notes
- Follow the pattern established by the fixed `createProjectHCL` function
- Use the schema manager to load and validate the machines schema
- Generate examples that demonstrate proper machine configuration
- Ensure the generated content is immediately valid and usable
