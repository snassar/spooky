# Troubleshooting Guide

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
go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
```

#### Empty Profile
```bash
# Profile exists but is empty
# Check if tests are running
go test -v ./...

# Check if coverage mode is correct
go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
```

#### Profile Path Issues
```bash
# Error: cannot find coverage profile
# Use absolute path
go run github.com/vladopajic/go-test-coverage/v2@latest \
  --config=./tests/testcoverage.yml \
  --profile=$(pwd)/cover.out
```

### Threshold Failures

#### File Threshold Failed
```bash
# Error: file coverage below threshold
# Check specific file coverage
go tool cover -func=./cover.out | grep "filename.go"

# Add tests for uncovered functions
# Or add coverage-ignore comment if appropriate
```

#### Package Threshold Failed
```bash
# Error: package coverage below threshold
# Check package coverage
go tool cover -func=./cover.out

# Focus on high-priority functions
# Add integration tests if unit tests are insufficient
```

#### Total Threshold Failed
```bash
# Error: total coverage below threshold
# Review overall coverage
go tool cover -func=./cover.out

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
go test ./spooky/... -coverprofile=./cover.out -covermode=atomic -coverpkg=./spooky/...
```

#### Large Coverage Files
```bash
# Coverage file is very large
# Check what's being included
go tool cover -func=./cover.out | wc -l

# Review exclusions
# Remove unnecessary packages from coverage
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

## Getting Help

### Debug Mode
Always run with debug mode when troubleshooting:
```bash
go run github.com/vladopajic/go-test-coverage/v2@latest \
  --config=./tests/testcoverage.yml \
  --debug
```

### Logs and Output
- Check console output for error messages
- Review coverage HTML report for details
- Check GitHub Actions logs for CI issues

### Common Solutions
1. **Reinstall tools**: `go install github.com/vladopajic/go-test-coverage/v2@latest`
2. **Regenerate profiles**: `go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...`
3. **Check paths**: Use absolute paths for config and profile files
4. **Validate config**: Check YAML syntax and regex patterns
5. **Update exclusions**: Add appropriate exclusions for untestable code

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
