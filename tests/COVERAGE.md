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