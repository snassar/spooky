# Spooky Development Guide

## Overview

This guide covers everything needed to develop and contribute to the spooky automation and orchestration tool.

## Prerequisites

### Required Tools

- **Go 1.24+**: [Download from golang.org](https://golang.org/dl/)
- **Just**: Command runner for build automation - [Installation instructions](https://github.com/casey/just)
- **Git**: Version control
- **golangci-lint**: Code linting

### Optional Tools

- **BadgerDB**: For local development and testing
- **SSH tools**: For testing SSH functionality
- **Podman/Docker**: For integration testing

## Quick Start

### 1. Clone and Setup

```bash
git clone https://github.com/your-org/spooky.git
cd spooky
```

### 2. Build and Test

```bash
just build           # Build development version
just test            # Run all tests
just check           # Run all checks (lint, format, deps, tests)
```

### 3. Development Workflow

```bash
just dev             # Build, test, and show version
just release         # Clean, check, build, test (release workflow)
```

## Development Environment

### Project Structure

```
spooky/
├── cmd/                 # CLI command implementations
├── internal/            # Internal packages
│   ├── cli/            # CLI management
│   ├── interfaces/     # System interfaces
│   ├── logging/        # Logging system
│   ├── schemas/        # HCL schema definitions
│   └── types/          # Type definitions
├── docs/               # Documentation
├── build/              # Build artifacts
├── justfile            # Build automation
├── go.mod              # Go module definition
└── README.md           # Project overview
```

### Key Development Files

- **`justfile`**: Build automation and development commands
- **`internal/interfaces/interfaces.go`**: Core system interfaces
- **`internal/types/types.go`**: Unified type definitions
- **`internal/schemas/schemas/`**: HCL schema definitions
- **`cmd/root.go`**: CLI root command
- **`cmd/project.go`**: Project management commands

## Build System

### Just Commands

The project uses [just](https://github.com/casey/just) for build automation. Key commands:

```bash
# Core development
just build              # Build development version
just test               # Run all tests
just lint               # Run linter
just fmt                # Format code
just check              # Run all checks

# Version management
just version            # Show version information
just build-release      # Build from git tag
just build-scalver      # Build with specific version

# Development workflow
just dev                # Build, test, show version
just release            # Full release workflow
just clean              # Clean build artifacts

# Git integration
just git-status         # Show git status
just git-info           # Show branch and commit info
```

### Version Management

Spooky uses ScalVer versioning format: `MAJOR.DATE.PATCH`

- **Development builds**: `0.20250101.0-dev-abc123`
- **Official releases**: `0.20250101.0` (from git tags)
- **Build time**: Automatically generated from current time
- **Git commit**: Short commit hash for development builds

## Development Workflow

### 1. Code Quality Standards

#### Interface-First Development
- Define interfaces before implementations
- Use interfaces for all public APIs
- Follow established interface patterns

#### Type Safety
- Use unified types from `internal/types/`
- Follow domain-specific type organization
- Maintain type composition patterns

#### Error Handling
- Use structured error types with context
- Implement proper error wrapping
- Provide actionable error messages

### 2. Testing Strategy

#### Test Organization
```bash
just test-unit          # Unit tests only
just test-integration   # Integration tests only
just test-coverage      # Tests with coverage report
```

#### Test Requirements
- **Minimum coverage**: 75% overall, 70% per package
- **Interface testing**: 100% interface method coverage
- **Error path testing**: 100% error path coverage
- **Integration testing**: End-to-end workflow testing

### 3. Code Review Process

#### Pre-commit Checks
```bash
just check              # Run all checks before committing
```

#### Review Checklist
- [ ] All interfaces are properly implemented
- [ ] Tests pass with required coverage
- [ ] Code follows established patterns
- [ ] Documentation is updated
- [ ] No compilation errors
- [ ] Linting passes

## Architecture Guidelines

### Interface Architecture

#### Core Principles
1. **Interface-first design**: Define contracts before implementations
2. **Dependency injection**: Use interfaces for all dependencies
3. **Loose coupling**: Components depend on interfaces, not concrete types
4. **Separation of concerns**: Clear boundaries between components

#### Key Interfaces
- **IntegrationManager**: Coordinates all system integrations
- **ProjectManager**: Manages project lifecycle
- **ConfigManager**: Handles configuration loading and validation
- **LogManager**: Manages logging operations
- **SchemaManager**: Handles schema loading and validation

### Type System

#### Type Organization
- **`internal/types/`**: Unified access to all types
- **Domain subpackages**: Organized by functional domain
- **Common types**: Shared types used across domains
- **Import aliases**: Consistent naming patterns

#### Type Patterns
- **Entity composition**: Use embedded structs for common functionality
- **Validation types**: Comprehensive validation with context
- **Error types**: Structured errors with operation context
- **Configuration types**: Hierarchical and composable

### Schema System

#### Schema Organization
- **`internal/schemas/schemas/`**: All HCL schema definitions
- **Domain-specific schemas**: Separate schemas for each domain
- **Validation rules**: Comprehensive validation patterns
- **Version management**: ScalVer versioning for schemas

#### Schema Development
- **Schema-first**: Define schemas before implementation
- **Validation**: Include comprehensive validation rules
- **Documentation**: Clear field descriptions and examples
- **Compatibility**: Maintain backward compatibility

## CLI Development

### Command Structure

#### Command Patterns
- **Noun-verb**: `spooky project init`, `spooky facts gather`
- **Long-form flags**: Use `--flag` instead of `-f`
- **Consistent help**: Clear descriptions and examples
- **Error handling**: User-friendly error messages

#### Command Implementation
```go
// Example command structure
var projectInitCmd = &cobra.Command{
    Use:   "init [project-name]",
    Short: "Initialize a new spooky project",
    Long:  `Detailed description...`,
    Args:  cobra.ExactArgs(1),
    RunE:  handleProjectInit,
}
```

### CLI Guidelines

#### Flag Conventions
- **Long-form only**: No short flags (`-f`)
- **Descriptive names**: Clear, self-documenting flag names
- **Consistent types**: Use appropriate data types
- **Default values**: Provide sensible defaults

#### Error Handling
- **User-friendly messages**: Clear, actionable error messages
- **Context information**: Include relevant context in errors
- **Exit codes**: Use appropriate exit codes for different errors
- **Help integration**: Suggest solutions in error messages

## Configuration Management

### HCL Configuration

#### Configuration Files
- **`project.hcl`**: Main project configuration
- **`machines.hcl`**: Machine inventory
- **`actions.hcl`**: Action definitions
- **`variables.hcl`**: Project variables

#### Configuration Validation
- **Schema validation**: Validate against HCL schemas
- **Cross-file validation**: Validate relationships between files
- **Environment integration**: Support environment variable overrides
- **Default values**: Provide sensible defaults

### Environment Variables

#### Supported Variables
- **`XDG_CONFIG_HOME`**: Configuration directory
- **`SPOOKY_CONFIG_PATH`**: Custom config file path
- **`SPOOKY_LOG_LEVEL`**: Logging level
- **`SPOOKY_FACTS_PATH`**: Facts database path
- **`SPOOKY_FACTS_FORMAT`**: Storage format (badgerdb/json)

## Testing Guidelines

### Test Organization

#### Test Structure
```
tests/
├── unit/              # Unit tests
├── integration/       # Integration tests
├── fixtures/          # Test data and fixtures
└── helpers/           # Test helper functions
```

#### Test Patterns
- **Table-driven tests**: For comprehensive validation
- **Mock implementations**: For interface testing
- **Test fixtures**: Reusable test data
- **Integration tests**: End-to-end workflow testing

### Test Requirements

#### Coverage Requirements
- **Overall coverage**: Minimum 75%
- **Package coverage**: Minimum 70% per package
- **Interface coverage**: 100% interface method coverage
- **Error path coverage**: 100% error path coverage

#### Test Quality
- **Deterministic**: Tests must be repeatable
- **Fast**: Unit tests should run quickly
- **Isolated**: Tests should not depend on each other
- **Comprehensive**: Test all code paths

## Performance Guidelines

### Performance Requirements

#### Benchmarks
- **Build time**: Fast incremental builds
- **Test time**: Quick test running
- **Memory usage**: Efficient memory allocation
- **Startup time**: Fast CLI startup

#### Optimization
- **Interface performance**: Efficient interface implementations
- **Memory management**: Proper resource cleanup
- **Concurrency**: Parallel processing where appropriate
- **Caching**: Intelligent caching strategies

## Security Guidelines

### Security Requirements

#### Data Protection
- **Encryption**: Support for sensitive data encryption
- **Key management**: Secure key storage and rotation
- **Access control**: Proper access control mechanisms
- **Audit logging**: Comprehensive audit trails

#### SSH Security
- **Authentication**: Secure SSH authentication
- **Key validation**: Proper SSH key validation
- **Connection security**: Secure connection handling
- **Host verification**: SSH host key verification

## Documentation Guidelines

### Documentation Requirements

#### Code Documentation
- **Interface documentation**: Complete interface documentation
- **Function documentation**: Clear function descriptions
- **Type documentation**: Comprehensive type documentation
- **Example usage**: Working code examples

#### User Documentation
- **CLI documentation**: Complete command documentation
- **Configuration guides**: Configuration examples and guides
- **Tutorials**: Step-by-step tutorials
- **Troubleshooting**: Common issues and solutions

## Release Process

### Release Workflow

#### Pre-release Checklist
- [ ] All tests pass
- [ ] Code coverage meets requirements
- [ ] Documentation is updated
- [ ] Version is updated
- [ ] Changelog is updated

#### Release Steps
1. **Create release branch**: `git checkout -b release/v1.0.0`
2. **Update version**: Update version in code
3. **Run full test suite**: `just release`
4. **Create git tag**: `git tag v1.0.0`
5. **Build release**: `just build-release`
6. **Test release binary**: Verify functionality
7. **Merge and push**: Merge to main and push tag

### Version Management

#### ScalVer Format
- **Format**: `MAJOR.DATE.PATCH`
- **Example**: `0.20250101.0`
- **Development**: `0.20250101.0-dev-abc123`

#### Version Increment Rules
- **Major**: Breaking changes
- **Date**: Build date (YYYYMMDD)
- **Patch**: Bug fixes and minor changes

## Troubleshooting

### Common Issues

#### Build Issues
```bash
# Clean and rebuild
just clean
just build

# Check dependencies
just check-deps

# Verify Go version
go version  # Should be 1.24+
```

#### Test Issues
```bash
# Run specific test
go test ./internal/types -v

# Run with race detection
go test -race ./...

# Check test coverage
just test-coverage
```

#### Linting Issues
```bash
# Run linter
just lint

# Auto-fix issues (if supported)
golangci-lint run --fix

# Check specific file
golangci-lint run internal/types/types.go
```

### Development Tips

#### Efficient Development
- **Use just commands**: Leverage the justfile for common tasks
- **Incremental builds**: Use `just build` for quick builds
- **Test frequently**: Run tests often during development
- **Check early**: Run `just check` before committing

#### Debugging
- **Use logging**: Add appropriate logging statements
- **Test isolation**: Create focused unit tests
- **Error context**: Include context in error messages
- **Documentation**: Keep documentation updated

## Contributing

### Contribution Guidelines

#### Code Standards
- **Follow established patterns**: Use existing code patterns
- **Interface compliance**: Implement interfaces completely
- **Test coverage**: Maintain high test coverage
- **Documentation**: Update documentation with changes

#### Pull Request Process
1. **Create feature branch**: `git checkout -b feature/description`
2. **Make changes**: Implement feature or fix
3. **Add tests**: Include appropriate tests
4. **Update documentation**: Update relevant documentation
5. **Run checks**: `just check`
6. **Submit PR**: Create pull request with clear description

#### Review Process
- **Code review**: All changes require review
- **Test verification**: Ensure tests pass
- **Documentation review**: Verify documentation updates
- **Integration testing**: Test with existing functionality

## Resources

### Documentation
- **Project README**: `README.md`
- **API Documentation**: `docs/` directory
- **Schema Documentation**: `internal/schemas/schemas/`
- **Interface Documentation**: `internal/interfaces/`

### Tools
- **Just**: [just.systems](https://just.systems/)
- **Cobra**: [github.com/spf13/cobra](https://github.com/spf13/cobra)
- **HCL**: [github.com/hashicorp/hcl](https://github.com/hashicorp/hcl)
- **golangci-lint**: [golangci-lint.run](https://golangci-lint.run/)

### Development Tools
- **Go**: [golang.org](https://golang.org/)
- **Git**: [git-scm.com](https://git-scm.com/)
- **BadgerDB**: [github.com/dgraph-io/badger](https://github.com/dgraph-io/badger)

---

This development guide should be updated as the project evolves. For questions or clarifications, please refer to the project documentation or create an issue.
