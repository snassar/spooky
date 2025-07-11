# Test Coverage Configuration

## Coverage Threshold Strategy

### Base Thresholds
- **File**: 50% - Minimum acceptable coverage for any file
- **Package**: 65% - Minimum acceptable coverage for packages
- **Total**: 60% - Minimum acceptable coverage for the entire project

### Override Strategy
- **Critical Core Files** (70-80%): `config.go`, `ssh.go`, `commands.go`
  - These contain the main business logic and should have high coverage
- **CLI Entry Point** (30%): `main.go`
  - Mostly initialization code, lower coverage acceptable
- **Test Infrastructure** (40-50%): Test helpers and infrastructure
  - Supporting code with moderate coverage requirements

## Pre-commit Hooks

### Setup Options

#### Option 1: Go Script (Cross-Platform - Recommended)
```bash
# Build the pre-commit script
go build -o scripts/pre-commit scripts/pre-commit.go

# Install on Windows
copy scripts\pre-commit.exe .git\hooks\pre-commit.exe

# Install on Unix-like systems
cp scripts/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

#### Option 2: pre-commit Framework
```bash
# Install pre-commit framework
pip install pre-commit

# Install hooks
pre-commit install
```

#### Option 3: PowerShell (Windows Only)
```powershell
# Set execution policy (run as administrator)
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# Copy PowerShell script
Copy-Item .git\hooks\pre-commit.ps1 .git\hooks\pre-commit
```

### What the Pre-commit Hook Does
1. Checks if Go files are staged for commit
2. Runs tests with coverage profiling
3. Validates coverage against configured thresholds
4. Prevents commit if coverage thresholds are not met
5. Provides helpful error messages and suggestions

### Bypassing Pre-commit Hooks
If you need to bypass coverage checks (emergency fixes, etc.):
```bash
git commit --no-verify -m "Emergency fix - bypassing coverage checks"
```

### Local Coverage Check

```bash
# Install the coverage tool (if not already installed)
go install github.com/vladopajic/go-test-coverage/v2@latest

# Run tests and generate coverage profile
go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...

# Check coverage against thresholds
go-test-coverage --config=./tests/testcoverage.yml

# Generate HTML report
go tool cover -html=./cover.out -o coverage.html
```

---

## 1. Add/Update Coverage Section in `README.md`

**Ensure the following is present and up-to-date:**
- A dedicated "Test Coverage" section
- Coverage requirements (thresholds)
- How to run coverage checks locally (with Go 1.24+)
- How to generate and view HTML reports
- How CI/CD enforces coverage

**Example:**
```markdown
## Test Coverage

This project enforces test coverage to ensure code quality.

### Coverage Requirements

- **Total Coverage:** Minimum 60%
- **Package Coverage:** Minimum 65%
- **File Coverage:** Minimum 50%

### Running Coverage Checks Locally

```bash
# Run all tests and check coverage thresholds
make check-coverage

# Generate HTML coverage report
make coverage-html

# Run coverage tool manually
go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
go run github.com/vladopajic/go-test-coverage/v2@latest --config=./tests/testcoverage.yml
```

### Viewing Coverage Reports

- Open `coverage.html` in your browser for a detailed report.
- Download coverage artifacts from GitHub Actions workflow runs.

### CI/CD

- Coverage is checked on every PR and push to main.
- PRs that drop coverage below thresholds will fail.
```

---

## 2. Document How to Run Coverage Checks Locally

**Ensure this is present in `README.md` and/or `docs/COVERAGE.md`:**
- Step-by-step instructions for running coverage checks
- How to install the coverage tool (with Go 1.24+)
- How to interpret results

**Example:**
```markdown
### Local Coverage Check

```bash
# Install the coverage tool (if not already installed)
go install github.com/vladopajic/go-test-coverage/v2@latest

# Run tests and generate coverage profile
go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...

# Check coverage against thresholds
go-test-coverage --config=./tests/testcoverage.yml

# Generate HTML report
go tool cover -html=./cover.out -o coverage.html
```
```

---

## 3. Explain Coverage Thresholds and Requirements

**Ensure this is present in `tests/COVERAGE.md` and referenced in the README:**
- What the thresholds are and why
- What code is excluded (automatic/manual)
- How thresholds are enforced (locally and in CI)

**Example:**
```markdown
## Coverage Thresholds

Thresholds are defined in `./tests/testcoverage.yml`:

```yaml
threshold:
  file: 50
  package: 65
  total: 60
```

- **File:** Each file must have at least 50% coverage
- **Package:** Each package must have at least 65% coverage
- **Total:** The project must have at least 60% overall coverage

See [tests/COVERAGE.md](tests/COVERAGE.md) for details and rationale.
```

---

## 4. Add Troubleshooting Section for Coverage Issues

**Ensure this is present in `docs/TROUBLESHOOTING.md` and referenced in the README:**
- Common problems (tool not found, profile not found, threshold failures, etc.)
- Solutions and debug tips

**Example:**
```markdown
## Troubleshooting Coverage

- **Tool not found:**  
  `go install github.com/vladopajic/go-test-coverage/v2@latest`

- **Profile not found:**  
  Run `go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...` first.

- **Threshold failures:**  
  Add or improve tests for uncovered code.

- **Debug:**  
  Run with `--debug` for more info:
  ```bash
  go-test-coverage --config=./tests/testcoverage.yml --debug
  ```
```

---

## 5. Update Contributing Guidelines with Coverage Requirements

**Ensure this is present in `CONTRIBUTING.md`:**
- State the minimum coverage requirements for PRs
- Require contributors to run coverage checks before submitting
- Mention that PRs will be blocked if coverage drops below thresholds

**Example:**
```markdown
## Test Coverage Requirements

- All PRs must maintain or improve test coverage.
- Minimum thresholds: 60% total, 65% per package, 50% per file.
- Run `make check-coverage` before submitting a PR.
- Add tests for all new features and bug fixes.
- Use `// coverage-ignore` only for code that cannot be unit tested (see [tests/COVERAGE.md](tests/COVERAGE.md)).
```

---

## 6. Cross-Check for Consistency

- All references to Go version should be **Go 1.24 or later**.
- All documentation should reference the correct coverage tool and config file locations.
- README, CONTRIBUTING.md, tests/COVERAGE.md, and docs/TROUBLESHOOTING.md should be in sync.

---

## Summary Table

| File                        | Section/Change                                                                 |
|-----------------------------|-------------------------------------------------------------------------------|
| README.md                   | Add/Update "Test Coverage" section, local/CI instructions, Go 1.24+           |
| tests/COVERAGE.md           | Thresholds, rationale, exclusions, enforcement, Go 1.24+                      |
| docs/TROUBLESHOOTING.md     | Coverage troubleshooting, Go 1.24+                                            |
| CONTRIBUTING.md             | Coverage requirements for PRs, Go 1.24+                                       |
| .github/workflows/*.yml     | Ensure Go version is 1.24, coverage steps are present                         |

---

**If you make these changes, you will fully satisfy [issue #27](https://github.com/snassar/spooky/issues/27) and have robust, clear, and up-to-date documentation for test coverage in your project.**

## Coverage Exclusions

This project excludes certain files and directories from coverage calculations to focus on production code quality.

### Excluded Patterns

- `_test\.go$` - All test files (unit and integration)
- `^tests/` - Entire tests directory
- `^tests/infrastructure/` - Test infrastructure servers
- `^tests/helpers/` - Test helper utilities
- `^examples/` - Example configurations and code
- `main\.go$` - Main entry point (minimal logic)
- `ssh_keygen\.go$` - SSH key generation utilities

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
- **`ssh_keygen.go`**: 60% - SSH key generation utilities
- **`commands.go`**: 55% - CLI commands with error handling

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
📁 Coverage artifacts: `coverage-reports` (HTML report: `coverage.html`)
```

### Local Diff Analysis
```bash
# Generate current coverage breakdown
go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
go run github.com/vladopajic/go-test-coverage/v2@latest \
  --config=./tests/testcoverage.yml \
  --breakdown-file-name=coverage-breakdown.json

# Compare with base (if available)
go run github.com/vladopajic/go-test-coverage/v2@latest \
  --config=./tests/testcoverage.yml \
  --diff-base-breakdown-file-name=base-coverage-breakdown.json
```

## Task 6: Update README.md

**File to modify:** `README.md`

**Add to Test Coverage section:**
```markdown
### Coverage Diff Tracking

- Coverage changes are automatically tracked in pull requests
- PR comments show coverage increases/decreases
- Coverage reports are available as workflow artifacts
- Decreases in coverage trigger warnings

For detailed coverage analysis, see [tests/COVERAGE.md](tests/COVERAGE.md).
```

## Task 7: Test diff tracking functionality

**Commands to test locally:**
```bash
# Generate current coverage breakdown
go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
go run github.com/vladopajic/go-test-coverage/v2@latest \
  --config=./tests/testcoverage.yml \
  --breakdown-file-name=coverage-breakdown.json

# Verify breakdown file was created
ls -la coverage-breakdown.json

# Check breakdown file content
cat coverage-breakdown.json
```

## Expected outcomes:

1. ✅ **Breakdown files generated** - JSON files with detailed coverage data
2. ✅ **Diff analysis** - Comparison between PR and base branch coverage
3. ✅ **PR comments** - Automatic coverage change notifications
4. ✅ **Documentation** - Clear explanation of diff tracking functionality
5. ✅ **Local testing** - Ability to test diff tracking locally

## Benefits:

- **Visibility**: Clear coverage changes in PRs
- **Accountability**: Developers see impact of their changes
- **Quality**: Prevents coverage regression
- **Transparency**: Coverage data available to all contributors
- **Automation**: No manual coverage tracking required

This implementation will provide comprehensive coverage diff tracking that helps maintain and improve code quality through automated reporting and notifications.