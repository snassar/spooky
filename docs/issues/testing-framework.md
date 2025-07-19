# Implement comprehensive testing framework for new features

## Background
As spooky grows with new features like facts management, template processing, and remote configuration sources, we need a robust testing framework that can validate all functionality across different environments and scales.

## Requirements

### Integration Testing Framework

#### Fact Collection Testing
```go
// Test fact gathering from various sources
func TestFactGathering(t *testing.T) {
    // Test SSH-based fact collection
    // Test local fact collection
    // Test OpenTofu integration
    // Test HCL config parsing
}
```

#### Template Processing Testing
```go
// Test template evaluation and rendering
func TestTemplateProcessing(t *testing.T) {
    // Test server-side template evaluation
    // Test fact injection into templates
    // Test template validation
    // Test template rendering with various data types
}

#### Remote Source Testing
```go
// Test remote configuration sources
func TestRemoteSources(t *testing.T) {
    // Test Git repository access
    // Test S3 bucket access
    // Test HTTP endpoint access
    // Test authentication and authorization
}
```

### Performance Testing

#### Scalability Tests
```bash
# Test with different server counts
spooky facts gather --hosts 10,100,1000
spooky execute config.hcl --parallel 1,10,50,100
```

#### Memory Usage Tests
```go
// Monitor memory usage during operations
func TestMemoryUsage(t *testing.T) {
    // Test fact storage memory usage
    // Test template processing memory usage
    // Test large configuration file processing
}
```

#### Database Performance Tests
```go
// Test BadgerDB performance
func TestDatabasePerformance(t *testing.T) {
    // Test fact storage performance
    // Test query performance
    // Test concurrent access
    // Test data import/export performance
}
```

### Security Testing

#### SSH Security Tests
```go
// Test SSH connection security
func TestSSHSecurity(t *testing.T) {
    // Test key authentication
    // Test connection encryption
    // Test host key verification
    // Test connection timeouts
}
```

#### Input Validation Tests
```go
// Test input validation and sanitization
func TestInputValidation(t *testing.T) {
    // Test HCL file validation
    // Test template input validation
    // Test fact data validation
    // Test remote URL validation
}
```

### Test Environment Setup

#### Podman-Based Test Environment
```bash
# Extend existing Podman test environment
# Add containers for:
# - Multiple OS types (Ubuntu, CentOS, Alpine)
# - Different SSH configurations
# - Mock S3-compatible storage
# - Git server for testing
```

#### Mock Services
```go
// Mock external services for testing
type MockS3Service struct {
    // Mock S3 operations
}

type MockGitService struct {
    // Mock Git operations
}

type MockSSHServer struct {
    // Mock SSH server for testing
}
```

## Implementation

### Phase 1: Core Testing Infrastructure
1. **Extend existing Podman test environment**
2. **Create mock services for external dependencies**
3. **Implement basic integration test framework**
4. **Add performance testing utilities**

### Phase 2: Feature-Specific Testing
1. **Fact collection integration tests**
2. **Template processing integration tests**
3. **Remote source integration tests**
4. **Database performance tests**

### Phase 3: Security and Compliance Testing
1. **SSH security validation tests**
2. **Input validation and sanitization tests**
3. **Access control and permission tests**
4. **Audit logging tests**

### Phase 4: Scalability and Performance Testing
1. **Large-scale deployment tests**
2. **Memory usage optimization tests**
3. **Concurrent operation tests**
4. **Stress testing**

## Test Categories

### Unit Tests
- **Individual component testing**
- **Mock-based testing**
- **Fast execution (< 1s per test)**
- **High coverage (> 90%)**

### Integration Tests
- **End-to-end workflow testing**
- **Real service integration**
- **Medium execution time (< 30s per test)**
- **Realistic scenarios**

### Performance Tests
- **Scalability testing**
- **Memory usage monitoring**
- **Long execution time (minutes)**
- **Resource utilization analysis**

### Security Tests
- **Input validation testing**
- **Authentication testing**
- **Authorization testing**
- **Encryption testing**

## Test Data Management

### Test Data Sets
```yaml
# test-data/servers.yaml
servers:
  - name: web-001
    host: 192.168.1.10
    os: ubuntu-22.04
    cpu_cores: 4
    memory_gb: 8
  - name: db-001
    host: 192.168.1.20
    os: centos-8
    cpu_cores: 8
    memory_gb: 16
```

### Test Templates
```hcl
# test-data/templates/nginx.conf.tmpl
server {
    listen {{ .port }};
    server_name {{ .hostname }};
    root {{ .web_root }};
}
```

### Test Configurations
```hcl
# test-data/configs/test-config.hcl
servers {
    web-001 {
        host = "192.168.1.10"
        user = "ubuntu"
    }
}

actions {
    deploy-nginx {
        description = "Deploy nginx configuration"
        command = "sudo systemctl restart nginx"
    }
}
```

## CI/CD Integration

### GitHub Actions Workflow
```yaml
# .github/workflows/test.yml
name: Comprehensive Testing

on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run unit tests
        run: go test ./... -v -race -cover

  integration-tests:
    runs-on: ubuntu-latest
    services:
      podman:
        image: docker://docker.io/library/hello-world
    steps:
      - uses: actions/checkout@v3
      - name: Run integration tests
        run: go test ./tests/integration/... -v

  performance-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run performance tests
        run: go test ./tests/performance/... -v -timeout 30m
```

## Success Criteria

### Coverage Requirements
- [ ] Unit test coverage > 90%
- [ ] Integration test coverage > 80%
- [ ] All critical paths tested
- [ ] Error conditions tested

### Performance Requirements
- [ ] Unit tests complete in < 30s
- [ ] Integration tests complete in < 5m
- [ ] Performance tests complete in < 30m
- [ ] Memory usage within acceptable limits

### Quality Requirements
- [ ] All tests pass consistently
- [ ] No flaky tests
- [ ] Clear test failure messages
- [ ] Comprehensive test documentation

## Dependencies

### Testing Libraries
- **github.com/stretchr/testify**: Assertions and mocking
- **github.com/ory/dockertest**: Docker-based testing
- **github.com/hashicorp/go-multierror**: Error handling
- **github.com/sirupsen/logrus**: Test logging

### Test Utilities
- **github.com/onsi/ginkgo**: BDD testing framework
- **github.com/onsi/gomega**: Matcher library
- **github.com/golang/mock**: Mock generation

## Implementation Notes

### Test Organization
```
tests/
├── unit/              # Unit tests
├── integration/       # Integration tests
├── performance/       # Performance tests
├── security/          # Security tests
├── data/              # Test data
└── utils/             # Test utilities
```

### Mock Services
- **Mock SSH server**: For testing SSH connections
- **Mock S3 service**: For testing S3 operations
- **Mock Git server**: For testing Git operations
- **Mock OpenTofu**: For testing OpenTofu integration

### Test Environment
- **Podman containers**: Multiple OS types and configurations
- **Network isolation**: Separate networks for different test scenarios
- **Resource limits**: Controlled resource usage for performance testing
- **Cleanup procedures**: Automatic cleanup after tests 