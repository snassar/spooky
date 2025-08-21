# AGENTS.md

A dedicated guide for AI coding agents working on the **spooky** project.

## Project Overview

**spooky** is a Go-based configuration management and automation tool that uses HCL (HashiCorp Configuration Language) schemas for validation. The project focuses on embedding schemas, validation rules, and test data into Go binaries for self-contained operation.

### Key Components

- **Schema System**: HCL-based schemas for project, machines, actions, variables, templates, logging, and secrets
- **Embedding**: Go `embed` directive to include schemas and test data in binaries
- **Validation**: Multi-layer validation (structure, rules, metadata)
- **Test Data**: Comprehensive test cases for schema validation

## Development Environment

### Setup Commands

```bash
# Install Go dependencies
go mod tidy

# Run the schema embedder demo
go run main.go

# Run tests
go test ./...

# Build the project
go build -o spooky main.go
```

### Project Structure

```
spooky/
├── internal/schemas/          # Schema definitions and embedding
│   ├── schemafiles/          # HCL schema files
│   │   ├── structure/        # Data structure definitions
│   │   ├── validation/       # Business logic validation rules
│   │   └── metadata/         # Schema metadata definitions
│   ├── testdata/             # Test data for validation
│   └── schema_embedder.go    # Schema embedding logic
├── cmd/                      # CLI commands
├── internal/                 # Core packages
└── main.go                   # Entry point
```

## Code Style & Conventions

### Go Code

- **Package naming**: Use descriptive package names (`schemas`, `machines`, `actions`)
- **Error handling**: Always return and check errors, use `fmt.Errorf` for context
- **Embedding**: Use `//go:embed` directive for static file inclusion
- **Testing**: Write tests for all public functions, use `_test.go` suffix
- **Documentation**: Add comments for exported functions and types

### HCL Schemas

- **Naming**: Use kebab-case for schema files (`project.hcl`, `machines.hcl`)
- **Structure**: Group related schemas in subdirectories
- **Validation**: Include both structure and business rule validation
- **Metadata**: Always include `metadata` blocks with ScalVer versioning

### File Organization

- **Schemas**: Place in `internal/schemas/schemafiles/`
- **Test Data**: Place in `internal/schemas/testdata/`
- **CLI Commands**: Place in `cmd/` directory
- **Core Logic**: Place in `internal/` directory

## Testing Instructions

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/schemas/...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...
```

### Test Data Validation

- **Valid cases**: Test data should pass all schema validations
- **Invalid cases**: Test data should trigger appropriate validation errors
- **Edge cases**: Include boundary conditions and error scenarios
- **Completeness**: Test data should cover all schema fields and rules

### Test Structure

```go
func TestSchemaValidation(t *testing.T) {
    // Setup
    embedder, err := schemas.NewSchemaEmbedder()
    require.NoError(t, err)
    
    // Test valid data
    validData := embedder.GetTestData("valid-project")
    // ... validation logic
    
    // Test invalid data
    invalidData := embedder.GetTestData("invalid-project")
    // ... error checking
}
```

## Schema Development

### Creating New Schemas

1. **Structure Schema**: Define data structure in `schemafiles/structure/`
2. **Validation Rules**: Add business logic in `schemafiles/validation/`
3. **Test Data**: Create test cases in `testdata/`
4. **Update Embedder**: Ensure new files are included in embedding

### Schema Validation Layers

1. **Structure Validation**: Field types, requirements, patterns
2. **Business Rules**: Complex validation logic and relationships
3. **Metadata Validation**: Version, type, and description consistency

### ScalVer Versioning

Use ScalVer format: `0.YYYYMMDD.N`
- Example: `0.20250809.0`
- Update version when schema changes
- Maintain backward compatibility when possible

## Git Workflow

### Branch Strategy

- **main**: Clean, stable branch
- **development**: Active development work
- **feature branches**: For specific features or fixes

### Commit Guidelines

**See `.cursor/rules/no-placeholder-commits.mdc` for detailed commit guidelines.**

Key principles:
- Only commit functional, working code
- Never commit placeholder or stub implementations
- Ensure code compiles and passes tests
- Include comprehensive test coverage

### Commit Message Format

```
type: brief description

- Detailed bullet points
- Additional context
- Impact and benefits
```

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `chore`

## Common Tasks

### Adding New Schema Fields

1. Update structure schema in `schemafiles/structure/`
2. Add validation rules in `schemafiles/validation/`
3. Update test data to include new fields
4. Test with both valid and invalid data
5. Update documentation if needed

### Embedding New Files

1. Place files in appropriate `schemafiles/` subdirectory
2. Update `schema_embedder.go` if new file types are needed
3. Add test data for new schemas
4. Verify embedding works correctly

### CLI Command Development

1. Create command in `cmd/` directory
2. Implement command logic in `internal/` packages
3. Add proper error handling and user feedback
4. Include help text and usage examples
5. Add tests for command functionality

## Error Handling

### Schema Validation Errors

- Provide clear, actionable error messages
- Include field names and expected values
- Reference schema documentation when helpful
- Group related validation errors together

### File System Errors

- Handle missing files gracefully
- Provide helpful suggestions for common issues
- Log errors with appropriate detail levels
- Return structured error information

## Performance Considerations

### Schema Loading

- Use Go `embed` directive for compile-time inclusion
- Cache loaded schemas in memory
- Minimize file system operations
- Optimize schema parsing for large files

### Validation Performance

- Validate only necessary fields
- Use efficient validation algorithms
- Cache validation results when appropriate
- Profile validation performance for large datasets

## Security Guidelines

### Input Validation

- Validate all HCL input before processing
- Sanitize file paths and user inputs
- Prevent path traversal attacks
- Validate template content for security

### Secret Handling

- Never log sensitive information
- Use secure methods for secret storage
- Validate secret formats and requirements
- Implement proper access controls

## Troubleshooting

### Common Issues

1. **Schema not found**: Check file paths and embedding configuration
2. **Validation errors**: Verify schema structure and test data
3. **Build failures**: Ensure all dependencies are properly imported
4. **Test failures**: Check test data matches current schema versions

### Debug Commands

```bash
# Check schema embedding
go run main.go

# Validate specific schema
go test -v ./internal/schemas/ -run TestSpecificSchema

# Check file structure
find internal/schemas/ -name "*.hcl" -type f
```

## Project Rules

### Cursor Rules

The project uses Cursor rules to maintain code quality and consistency:

- **`.cursor/rules/no-placeholder-commits.mdc`** - Prevents committing placeholder or stub code
- *Additional rules will be added here as the project evolves*

### Rule Management

- Rules are stored in `.cursor/rules/` directory
- Each rule file uses `.mdc` format with frontmatter
- Rules can specify file patterns and application scope
- Update this section when adding new rules

## Resources

- [HCL Documentation](https://github.com/hashicorp/hcl)
- [Go Embed Directive](https://pkg.go.dev/embed)
- [AGENTS.md Specification](https://github.com/openai/agents.md)
- [Project Documentation](./docs/)

---

*This AGENTS.md file provides comprehensive guidance for AI coding agents working on the spooky project. Update this file as the project evolves to ensure agents have the most current information.*
