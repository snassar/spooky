# Test Coverage

[← Back to Index](index.md) | [Next: Contributing](contributing.md)

---

## Coverage Policy

- **Total Coverage:** Minimum 60%
- **Package Coverage:** Minimum 65%
- **File Coverage:** Minimum 50%
- **Critical Code:** Higher thresholds for SSH and configuration code

## Checking Coverage Locally

```bash
# Run all tests and check coverage thresholds
make check-coverage

# Generate HTML coverage report
make coverage-html

# Run coverage tool manually
make install-development-tools
go test ./... -coverprofile=./tests/coverage.out -covermode=atomic -coverpkg=./...
go run github.com/vladopajic/go-test-coverage/v2@latest --config=./tests/testcoverage.yml
```

## Viewing Coverage Reports
- Open `tests/reports/coverage.html` in your browser for a detailed report
- Download coverage artifacts from CI workflow runs

## CI/CD
- Coverage is checked on every PR and push to main
- PRs that drop coverage below thresholds will fail

## Troubleshooting
See [Troubleshooting](troubleshooting.md) for common coverage issues.

## Advanced Coverage Documentation
For detailed coverage configuration, threshold rationale, and advanced troubleshooting, see [Advanced Coverage Guide](coverage-advanced.md).

---

## Navigation
- [← Back to Index](index.md)
- [Previous: Troubleshooting](troubleshooting.md)
- [Next: Contributing](contributing.md) 