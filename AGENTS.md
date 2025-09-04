# AGENTS.md

A dedicated guide for AI coding agents working on the **spooky** project.

## Project Overview

**spooky** is a Go-based configuration management and automation tool. The project focuses on creating self-contained binaries with embedded configuration schemas and validation rules.

## Development Environment

### Setup Commands

```bash
# Install Go dependencies
go mod tidy

# Run main spooky application
go run main.go

# Run tests
go test ./...

# Run tests
go test ./...

# Build the project
go build -o spooky main.go
```

### Project Structure

```
spooky/
├── commands/                 # Main spooky commands
├── internal/                 # Core packages and business logic
├── tools/                    # Development utilities (currently empty)
├── documentation/            # Project documentation
└── main.go                   # Entry point
```

## Code Style & Conventions

### Go Code

- **Package naming**: Use descriptive package names (`schemas`, `machines`, `actions`)
- **Error handling**: Always return and check errors, use `fmt.Errorf` for context
- **Embedding**: Use `//go:embed` directive for static file inclusion
- **Testing**: Write tests for all public functions, use `_test.go` suffix
- **Documentation**: Add comments for exported functions and types

### File Organization

- **Main Commands**: Place in `commands/` directory
- **Development Tools**: Place in `tools/` directory
- **Core Logic**: Place in `internal/` directory
- **Documentation**: Place in `documentation/` directory

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

### Test Structure

```go
func TestFunctionality(t *testing.T) {
    // Setup
    // ... setup code
    
    // Test valid cases
    // ... test logic
    
    // Test invalid cases
    // ... error checking
}
```

## Tool Architecture

### Main Application

The project focuses on a single main application:

**Main Application (`spooky`):**
- Focus: Automation and configuration management
- Commands: `commands/` directory
- Purpose: Core automation functionality (like Ansible)
- Structure: Self-contained binary with embedded configuration schemas

### Development Tools

**Note**: The `tools/` directory is currently empty. Development tools mentioned in previous documentation (`spooky-schemas-tool`, `spooky-hcl-tool`) are not currently implemented.

If development tools are needed in the future, they should:
- Be placed in the `tools/` directory with their own `main.go` files
- Share the same `go.mod` file for dependency consistency
- Import from `internal/` packages for shared functionality
- Use `spf13/cobra` for consistent CLI experience

## Development Workflow

### Adding New Features

1. **Plan**: Understand requirements and design approach
2. **Implement**: Write functional code with proper error handling
3. **Test**: Create comprehensive tests for new functionality
4. **Document**: Update documentation as needed
5. **Review**: Ensure code follows project conventions

### Code Quality

- Write clear, readable code
- Include proper error handling
- Add comprehensive tests
- Follow Go best practices
- Document public APIs

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

### Adding New Features

1. Understand the feature requirements
2. Implement the functionality with proper error handling
3. Add comprehensive tests for the new feature
4. Update documentation as needed
5. Ensure the feature follows project conventions

### Modifying Existing Code

1. Understand the current implementation
2. Make changes while maintaining backward compatibility
3. Update tests to reflect changes
4. Verify existing functionality still works
5. Update documentation if APIs change

### Main Application Command Development

1. Create command in `commands/` directory
2. Implement command logic in `internal/` packages
3. Add proper error handling and user feedback
4. Include help text and usage examples
5. Add tests for command functionality

### Development Tool Development (Future)

If development tools are needed:
1. Create tool in `tools/` directory with its own `main.go`
2. Import shared functionality from `internal/` packages
3. Use `spf13/cobra` for consistent CLI experience
4. Add proper error handling and user feedback
5. Include help text and usage examples
6. Add tests for tool functionality

## Error Handling

### Validation Errors

- Provide clear, actionable error messages
- Include relevant context and expected values
- Reference documentation when helpful
- Group related errors together

### File System Errors

- Handle missing files gracefully
- Provide helpful suggestions for common issues
- Log errors with appropriate detail levels
- Return structured error information

## Code Quality and Refactoring

### Refactoring Patterns

The project follows established refactoring patterns to maintain code quality:

**Function Decomposition:**
- Break large functions (>50 lines) into smaller, focused helper functions
- Each function should have a single responsibility
- Target cyclomatic complexity ≤ 10 (acceptable up to 15 for complex operations)

**Helper Function Naming:**
- `extract*` functions for data extraction logic
- `parse*` functions for parsing and validation logic  
- `process*` functions for data processing and transformation
- `validate*` functions for validation logic

**Common Refactoring Patterns:**
- Extract method for long functions with clear responsibilities
- Extract constant for magic numbers or strings used multiple times
- Simplify conditionals for complex nested if statements
- Consolidate duplicate code into reusable functions

### Testing Patterns

**Unit Testing:**
- Create comprehensive tests for all extracted helper functions
- Test both success and error cases
- Use table-driven tests for multiple scenarios
- Include benchmarks for performance-critical functions

**Test Structure:**
```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name        string
        input       InputType
        expected    ExpectedType
        expectError bool
        errorMsg    string
    }{
        // Test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

**Benchmark Testing:**
- Create benchmarks for performance-critical functions
- Use `-benchmem` flag to measure memory allocation
- Target sub-microsecond performance for simple operations
- Document performance characteristics

## Performance Considerations

### Resource Loading

- Use efficient loading mechanisms
- Cache frequently accessed data in memory
- Minimize file system operations
- Optimize parsing for large files

### Performance Optimization

- Profile code for bottlenecks
- Use efficient algorithms and data structures
- Cache results when appropriate
- Optimize for common use cases

### Performance Validation

- Run benchmarks after refactoring to ensure no regression
- Monitor memory allocation patterns
- Validate performance characteristics in CI/CD pipeline

## Security Guidelines

### Input Validation

- Validate all user input before processing
- Sanitize file paths and user inputs
- Prevent path traversal attacks
- Validate content for security

### Secret Handling

- Never log sensitive information
- Use secure methods for secret storage
- Validate secret formats and requirements
- Implement proper access controls

## Troubleshooting

### Common Issues

1. **File not found**: Check file paths and configuration
2. **Validation errors**: Verify data structure and test cases
3. **Build failures**: Ensure all dependencies are properly imported
4. **Test failures**: Check test data matches current implementations

### Debug Commands

```bash
# Run main application
go run main.go

# Run tests
go test ./...

# Run specific tests
go test -v ./internal/package/ -run TestSpecificFunction

# Check file structure
find . -name "*.go" -type f
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

- [Go Documentation](https://golang.org/doc/)
- [AGENTS.md Specification](https://github.com/openai/agents.md)
- [Project Documentation](./docs/)

---

*This AGENTS.md file provides comprehensive guidance for AI coding agents working on the spooky project. Update this file as the project evolves to ensure agents have the most current information.*
