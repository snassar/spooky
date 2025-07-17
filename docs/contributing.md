# Contributing

[← Back to Index](index.md) | [Next: Development Tools](tools.md)

---

## Development Setup

### Prerequisites
- Go 1.24 or later
- SSH access to target servers (for integration tests)
- Git

### Local Setup
```bash
git clone <repository-url>
cd spooky
go mod download
make build
```

## Test Coverage Requirements

- All PRs must maintain or improve test coverage
- Minimum thresholds: 60% total, 65% per package, 50% per file
- Run `make check-coverage` before submitting a PR
- Add tests for all new features and bug fixes
- Use `// coverage-ignore` only for code that cannot be unit tested

## Testing Requirements

### Test Coverage Standards

All contributions must maintain or improve test coverage:

- **Minimum Total Coverage**: 60%
- **Minimum Package Coverage**: 65%
- **Minimum File Coverage**: 50%

### Running Tests

#### Before Submitting
```bash
# Run all tests with coverage
make test

# Check coverage thresholds
make check-coverage

# Generate coverage report
make coverage-html
```

#### Test Types
- **Unit Tests**: Fast, isolated tests for individual functions
- **Integration Tests**: End-to-end tests with mock SSH servers
- **Coverage Tests**: Ensure coverage thresholds are met

### Adding New Features

#### Code Requirements
1. **Add tests** for all new functionality
2. **Maintain coverage** above thresholds
3. **Include integration tests** for SSH operations
4. **Document changes** in README and code comments

#### Test Guidelines
- Write unit tests for business logic
- Write integration tests for SSH operations
- Use mock servers for testing (see [Test Environment](test-environment.md))
- Test error conditions and edge cases

### Coverage Guidelines

#### What to Test
- **Business logic**: Core functionality and algorithms
- **Error handling**: Error paths and edge cases
- **Configuration**: Parsing and validation
- **SSH operations**: Connection and command execution

#### What Not to Test
- **CLI entry points**: Tested via integration tests
- **File system errors**: Hard to test reliably
- **Network errors**: Tested via integration tests
- **Main functions**: Entry points, minimal logic

#### Coverage-Ignore Comments
Use `// coverage-ignore` sparingly:
```go
// coverage-ignore: CLI entry point, tested via integration tests
func main() {
    // ...
}
```

## Pull Request Process

### Before Submitting
1. **Run tests**: `make test`
2. **Check coverage**: `make check-coverage`
3. **Generate report**: `make coverage-html`
4. **Review coverage**: Open `tests/reports/coverage.html` in browser
5. **Fix issues**: Address any coverage or test failures

### PR Requirements
- [ ] All tests pass
- [ ] Coverage thresholds met
- [ ] New code has appropriate tests
- [ ] Documentation updated
- [ ] Code follows project style

### Coverage Checks
- CI will automatically check coverage
- PR will be blocked if thresholds not met
- Review coverage report in CI artifacts

## Code Style

### Go Conventions
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Run `make fmt` before committing

### Project Conventions
- Use HCL2 for configuration files
- Include descriptive comments
- Add examples for new features
- Update documentation for changes
- Use long-form CLI flags for cross-platform compatibility
- Co-locate unit tests with source files

## Review Process

### What Reviewers Check
- Code quality and style
- Test coverage and quality
- Documentation updates
- Security considerations
- Performance impact
- Cross-platform compatibility

### Coverage Review
- Verify coverage thresholds are met
- Check that new code is tested
- Review coverage-ignore usage
- Ensure integration tests exist

## Troubleshooting

### Common Issues
- **Coverage too low**: Add tests for uncovered code
- **Tests failing**: Check mock server setup
- **Integration issues**: Verify SSH server configuration

### Getting Help
- Check [Troubleshooting](troubleshooting.md)
- Review [Advanced Coverage Guide](coverage-advanced.md)
- Open an issue for complex problems

## Questions?

If you have questions about contributing:
1. Check the documentation
2. Review existing issues and PRs
3. Open an issue for clarification
4. Ask in discussions

---

## Navigation
- [← Back to Index](index.md)
- [Previous: Test Coverage](coverage.md)
- [Next: Development Tools](tools.md) 