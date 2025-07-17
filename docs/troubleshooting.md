# Troubleshooting

[← Back to Index](index.md) | [Next: Test Coverage](coverage.md)

---

## Test Coverage Issues

### Coverage Tool Problems

#### Tool Not Found
```bash
# Error: command not found: go-test-coverage
go install github.com/vladopajic/go-test-coverage/v2@latest
```

#### Permission Denied
```bash
# Error: permission denied
# On Windows, run PowerShell as Administrator
# On Linux/macOS, check file permissions
chmod +x $(go env GOPATH)/bin/go-test-coverage
```

#### Version Conflicts
```bash
# Check installed version
go-test-coverage --version

# Reinstall if needed
go install github.com/vladopajic/go-test-coverage/v2@latest
```

### Coverage Profile Issues

#### Profile Not Found
```bash
# Error: coverage profile not found
# Generate profile first
go test ./... -coverprofile=./tests/coverage.out -covermode=atomic -coverpkg=./...
```

#### Empty Profile
```bash
# Profile exists but is empty
# Check if tests are running
go test -v ./...

# Check if coverage mode is correct
go test ./... -coverprofile=./tests/coverage.out -covermode=atomic -coverpkg=./...
```

#### Profile Path Issues
```bash
# Error: cannot find coverage profile
# Use absolute path
go run github.com/vladopajic/go-test-coverage/v2@latest \
  --config=./tests/testcoverage.yml \
  --profile=$(pwd)/tests/coverage.out
```

### Threshold Failures

#### File Threshold Failed
```bash
# Error: file coverage below threshold
# Check specific file coverage
go tool cover -func=./tests/coverage.out | grep "filename.go"

# Add tests for uncovered functions
# Or add coverage-ignore comment if appropriate
```

#### Package Threshold Failed
```bash
# Error: package coverage below threshold
# Check package coverage
go tool cover -func=./tests/coverage.out

# Focus on high-priority functions
# Add integration tests if unit tests are insufficient
```

#### Total Threshold Failed
```bash
# Error: total coverage below threshold
# Review overall coverage
go tool cover -func=./tests/coverage.out

# Identify largest uncovered areas
# Prioritize testing based on business impact
```

### Configuration Issues

#### Config File Not Found
```bash
# Error: config file not found
# Check file exists
ls -la ./tests/testcoverage.yml

# Use absolute path
go run github.com/vladopajic/go-test-coverage/v2@latest \
  --config=$(pwd)/tests/testcoverage.yml
```

#### Invalid Config Syntax
```bash
# Error: invalid YAML
# Validate YAML syntax
yamllint ./tests/testcoverage.yml

# Check indentation and format
```

#### Exclusion Patterns Not Working
```bash
# Files still included despite exclusions
# Check regex patterns
# Test patterns manually
echo "filename.go" | grep -E "pattern"

# Use debug mode
go run github.com/vladopajic/go-test-coverage/v2@latest \
  --config=./tests/testcoverage.yml \
  --debug
```

### Performance Issues

#### Slow Coverage Generation
```bash
# Coverage generation is slow
# Exclude unnecessary files
# Use more specific package patterns
go test ./spooky/... -coverprofile=./tests/coverage.out -covermode=atomic -coverpkg=./spooky/...
```

#### Large Coverage Files
```bash
# Coverage file is very large
# Check what's being included
go tool cover -func=./tests/coverage.out | wc -l

# Review exclusions
# Remove unnecessary packages from coverage
```

## Test Environment Management

### Using spooky-test-env

The `spooky-test-env` tool can be used in two ways:

#### Method 1: Built Binary (Recommended)
```bash
# Build and install the tool
make build-test-env
make install-test-env

# Use the installed binary
./build/spooky-test-env preflight
./build/spooky-test-env start
./build/spooky-test-env status
./build/spooky-test-env stop
./build/spooky-test-env cleanup
```

#### Method 2: Go Run (Alternative)
```bash
# Run directly with go run (no build required)
go run tools/spooky-test-env/main.go preflight
go run tools/spooky-test-env/main.go start
go run tools/spooky-test-env/main.go status
go run tools/spooky-test-env/main.go stop
go run tools/spooky-test-env/main.go cleanup

# Or create an alias for convenience
alias spooky-test-env='go run tools/spooky-test-env/main.go'
spooky-test-env preflight
spooky-test-env start
```

### Integration Test Issues

#### Integration Tests Not Running
```bash
# Error: no tests found
# Check build tags
go test -tags=integration ./tests/integration/...

# Verify test files have correct tags
# //go:build integration
```

#### Mock Server Issues
```bash
# Integration tests failing
# Check mock server ports
# Verify SSH server is running
netstat -an | grep :3100

# Restart mock servers
make test-integration
```

### Port Conflicts

If you get port binding errors, ensure ports 2221-2223 are not in use:

```bash
sudo netstat -tlnp | grep :222
```

### Container Startup Issues

Check container logs:

```bash
podman logs spooky-server1
podman logs spooky-server2
podman logs spooky-server3
```

### Network Issues

If the network doesn't exist:

```bash
podman network create spooky-test
```

## Getting Help

If you're still experiencing issues:

1. Check the [Advanced Coverage Guide](coverage-advanced.md) for detailed coverage information
2. Review the [Test Environment](test-environment.md) documentation
3. Open an issue on GitHub with:
   - Error messages and logs
   - Steps to reproduce
   - Environment details (OS, Go version, etc.)

---

## Navigation
- [← Back to Index](index.md)
- [Previous: Test Environment](test-environment.md)
- [Next: Test Coverage](coverage.md) 