# Test Coverage Configuration

## Coverage Exclusions

This project excludes certain files and directories from coverage calculations to focus on production code quality.

### Excluded Patterns

- `_test\.go$` - All test files (unit and integration)
- `^tests/` - Entire tests directory
- `^tests/infrastructure/` - Test infrastructure servers
- `^tests/helpers/` - Test helper utilities
- `^examples/` - Example configurations and code
- `main\.go$` - Main entry point (minimal logic)

### Rationale

1. **Test Files**: Test code itself doesn't need coverage - we care about production code
2. **Infrastructure**: Mock servers and test infrastructure are not production code
3. **Examples**: Example files are for documentation, not production use
4. **Utilities**: SSH key generation is a utility, not core functionality
5. **Main**: Entry point typically has minimal logic

### Coverage Focus

The coverage metrics focus on:
- Core SSH client functionality (`spooky/ssh.go`)
- Configuration parsing (`spooky/config.go`)
- Command-line interface (`commands.go`)

This ensures coverage reflects the quality of the actual application code.

## Coverage-Ignore Comments

Some code sections are marked with `// coverage-ignore` comments for the following reasons:

### CLI Entry Points
- `main()` function - Entry point, tested via integration tests
- `generateSSHKeys()` - CLI tool, tested via integration tests

### Error Handling
- File system errors - Hard to reliably test in unit tests
- CLI validation errors - Tested via integration tests
- Network errors - Tested via integration tests

### Rationale
These sections are excluded from coverage because:
1. They are tested via integration tests rather than unit tests
2. They involve external dependencies (file system, network) that are hard to mock
3. They are CLI entry points that are better tested end-to-end

### Guidelines for Adding Coverage-Ignore
- Only use for code that cannot be meaningfully unit tested
- Document the rationale in the comment
- Ensure the code is covered by integration tests
- Review periodically to see if unit testing has become feasible

## Verification commands:

```bash
# Run coverage to see the impact
go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
go run github.com/vladopajic/go-test-coverage/v2@latest --config=./tests/testcoverage.yml

# Check that coverage is more accurate (focused on testable code)
go run github.com/vladopajic/go-test-coverage/v2@latest --config=./tests/testcoverage.yml --debug
```

This approach will make coverage metrics more meaningful by focusing on code that can actually be unit tested, while ensuring that untestable code is still covered by integration tests.