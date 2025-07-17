# Advanced Coverage Guide

[← Back to Coverage](coverage.md) | [← Back to Index](index.md)

---

This guide provides detailed information about test coverage configuration, thresholds, and advanced troubleshooting for developers working on spooky.

## Coverage Threshold Strategy

### Base Thresholds
- **File**: 50% - Minimum acceptable coverage for any file
- **Package**: 65% - Minimum acceptable coverage for packages
- **Total**: 60% - Minimum acceptable coverage for the entire project

### Override Strategy
- **Critical Core Files** (70-80%): `config.go`, `ssh.go`, `internal/cli/commands.go`
  - These contain the main business logic and should have high coverage
- **CLI Entry Point** (30%): `main.go`
  - Mostly initialization code, lower coverage acceptable
- **Test Infrastructure** (40-50%): Test helpers and infrastructure
  - Supporting code with moderate coverage requirements

## Test Structure

```
tests/
├── integration/           # Integration tests using gliderlabs/ssh
│   ├── config_integration_test.go
│   └── ssh_integration_test.go
├── fixtures/              # Test data and configuration files
│   ├── valid_config.hcl
│   └── invalid_config.hcl
├── helpers/               # Test utilities and helpers
│   └── test_helpers.go
└── infrastructure/        # Mock SSH servers for testing
    ├── simple_server/
    ├── public_key_server/
    ├── sftp_server/
    └── timeout_server/
```

## Test Types

### Unit Tests
- **Location**: Co-located with source files (e.g., `config_test.go` next to `config.go`)
- **Purpose**: Test individual functions and methods in isolation
- **Dependencies**: Minimal external dependencies
- **Speed**: Fast execution

### Integration Tests
- **Location**: `tests/integration/`
- **Purpose**: Test component interactions and end-to-end workflows
- **Dependencies**: Uses [gliderlabs/ssh](https://github.com/gliderlabs/ssh) for mock SSH servers
- **Speed**: Slower due to network operations and server setup

### Mock SSH Servers
The integration tests use [gliderlabs/ssh](https://github.com/gliderlabs/ssh) to create mock SSH servers that:
- Accept any password/key for testing
- Record executed commands
- Return predefined outputs
- Support custom command responses
- Run on dynamic ports to avoid conflicts

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
4. **Main**: Entry point typically has minimal logic


### Coverage Focus

The coverage metrics focus on:
- Core SSH client functionality (`spooky/ssh.go`)
- Configuration parsing (`spooky/config.go`)

- Command-line interface (`internal/cli/commands.go`)

This ensures coverage reflects the quality of the actual application code.

## Coverage-Ignore Comments

Some code sections are marked with `// coverage-ignore` comments for the following reasons:

### CLI Entry Points
- `main()` function - Entry point, tested via integration tests


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

## Coverage Thresholds and Requirements

### Threshold Configuration

Coverage thresholds are defined in `./tests/testcoverage.yml`:

```yaml
threshold:
  file: 50      # Individual files must have at least 50% coverage
  package: 65   # Packages must have at least 65% coverage
  total: 60     # Overall project must have at least 60% coverage
```

### Threshold Rationale

#### File-Level Threshold (50%)
- **Why 50%**: Some files contain error handling that's hard to test
- **Examples**: CLI validation, file system operations
- **Exclusions**: Files with `// coverage-ignore` comments

#### Package-Level Threshold (65%)
- **Why 65%**: Ensures core functionality is well-tested
- **Focus**: Business logic and critical paths
- **Current**: Project achieves ~61.7%, so 65% is achievable

#### Total Project Threshold (60%)
- **Why 60%**: Balances quality with practicality
- **Realistic**: Accounts for hard-to-test code
- **Maintainable**: Can be achieved and maintained

### Coverage Requirements by Code Type

#### High Priority (Must Test)
- SSH connection logic
- Configuration parsing
- Command execution
- Error handling paths

#### Medium Priority (Should Test)
- CLI argument parsing
- File operations
- Utility functions

#### Low Priority (Optional)
- Main entry points (tested via integration)
- Error exit paths (tested via integration)

### Exclusions and Ignored Code

#### Automatic Exclusions
- Test files (`_test.go`)
- Test infrastructure (`tests/`)
- Example files (`examples/`)
- Main entry point (`main.go`)

#### Manual Exclusions
- Code marked with `// coverage-ignore`
- CLI entry points
- File system error handling
- Network error handling

### Threshold Enforcement

#### CI/CD Pipeline
- Coverage checks run on every PR
- Workflow fails if thresholds not met
- Prevents merging code that reduces coverage

#### Local Development
- `make check-coverage` validates thresholds
- Immediate feedback on coverage issues
- HTML reports for detailed analysis

### Updating Thresholds

#### When to Increase
- Project has matured
- Test coverage has improved
- New testing patterns established

#### When to Decrease
- Thresholds are unrealistic
- Project scope has changed
- New code types added

#### Process
1. Update `./tests/testcoverage.yml`
2. Update documentation
3. Communicate changes to team
4. Monitor impact on development velocity

## Optimized Coverage Thresholds

### Threshold Strategy

Coverage thresholds are optimized for the spooky project structure:

#### Base Thresholds (Default)
- **File:** 50% - Allows for hard-to-test error paths
- **Package:** 65% - Ensures good overall package quality
- **Total:** 60% - Realistic target for development velocity

#### Package-Specific Overrides

##### Critical Code (High Thresholds)
- **`spooky/ssh.go`**: 75% - SSH connections are security-critical
- **`spooky/config.go`**: 70% - Configuration parsing affects all operations

##### Standard Code (Medium Thresholds)
- **`internal/cli/commands.go`**: 55% - CLI commands with error handling

##### Excluded Code
- **`main.go`**: Entry point, tested via integration
- **Test files**: Not production code
- **Examples**: Documentation only

### Rationale

#### Why Different Thresholds?
1. **Security**: SSH code has higher thresholds due to security implications
2. **Complexity**: Simple utilities have lower thresholds than complex logic
3. **Testability**: Error handling paths are harder to test than happy paths
4. **Development Speed**: Balanced thresholds maintain quality without slowing development

#### Threshold Categories
- **Critical (75%)**: Security-sensitive code, core functionality
- **Important (70%)**: Configuration, data parsing
- **Standard (60%)**: Business logic, utilities
- **Basic (55%)**: CLI, error handling, edge cases

### Monitoring and Adjustment

#### When to Adjust Thresholds
- **Increase**: When coverage consistently exceeds thresholds
- **Decrease**: When thresholds are unrealistic for certain code types
- **Review**: Quarterly review of threshold effectiveness

#### Development Velocity Impact
- Current thresholds balance quality with development speed
- Critical code maintains high standards
- Non-critical code allows for faster iteration

## Coverage Diff Tracking

### Overview
The project tracks coverage changes between branches to ensure code quality improvements.

### How It Works
1. **Base Coverage**: Coverage from the target branch (main) is stored as a baseline
2. **PR Coverage**: Coverage from the PR branch is calculated
3. **Diff Analysis**: Changes are computed and reported
4. **PR Comments**: Coverage changes are posted as PR comments

### Coverage Reports
- **Current Coverage**: Coverage percentage for the PR branch
- **Base Coverage**: Coverage percentage from the target branch
- **Change**: Difference between current and base coverage
- **Details**: File, package, and total coverage breakdown

### PR Comments
Coverage changes are automatically posted to PRs:
- 📈 **Green arrow**: Coverage increased
- 📉 **Red arrow**: Coverage decreased
- ⚠️ **Warning**: Coverage decreased (requires attention)

### Example PR Comment
```markdown
## Test Coverage Report 📈

**Current Coverage:** 62.3%
**Base Coverage:** 61.7%
**Change:** +0.6%

### Coverage Details
- **File Coverage:** 52%
- **Package Coverage:** 67%
- **Total Coverage:** 62.3%

📊 [View workflow run](https://github.com/owner/repo/actions/runs/123456789)
📄 Coverage artifacts: `coverage-reports` (HTML report: `tests/reports/coverage.html`)
```

### Local Diff Analysis
```bash
# Generate current coverage breakdown
go test ./... -coverprofile=./tests/coverage.out -covermode=atomic -coverpkg=./...
go run github.com/vladopajic/go-test-coverage/v2@latest \
  --config=./tests/testcoverage.yml \
  --breakdown-file-name=tests/coverage-breakdown.json

# Compare with base (if available)
go run github.com/vladopajic/go-test-coverage/v2@latest \
  --config=./tests/testcoverage.yml \
  --diff-base-breakdown-file-name=tests/base-coverage-breakdown.json
```

## Benefits

- **Visibility**: Clear coverage changes in PRs
- **Accountability**: Developers see impact of their changes
- **Quality**: Prevents coverage regression
- **Transparency**: Coverage data available to all contributors
- **Automation**: No manual coverage tracking required

This implementation will provide comprehensive coverage diff tracking that helps maintain and improve code quality through automated reporting and notifications.

---

## Navigation
- [← Back to Coverage](coverage.md)
- [← Back to Index](index.md)