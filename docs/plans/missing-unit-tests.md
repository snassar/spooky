# Missing Unit Tests

This document outlines the unit tests that are missing from the spooky project. These tests are essential for ensuring code quality, reliability, and maintainability.

## Overview

Based on analysis of the codebase, there are **significant gaps** in unit test coverage across multiple packages. Many public functions lack corresponding test functions, which could lead to undetected bugs and regressions.

## Missing Tests by Package

### 1. internal/cli

#### commands.go
The following functions lack unit tests:

- `listFromConfigFile()` - No test for listing from config file
- `listMachines()` - Currently returns "not yet implemented" error
- `listFacts()` - No test for listing facts
- `getUniqueServers()` - No test for extracting unique servers from facts
- `listTemplates()` - Currently returns "not yet implemented" error  
- `listConfigs()` - Currently returns "not yet implemented" error
- `listActions()` - Currently returns "not yet implemented" error
- `initProject()` - No test for project initialization
- `validateConfigFile()` - No test for config file validation
- `exportFacts()` - No test for facts export functionality
- `importFacts()` - No test for facts import functionality

#### config.go
- `loadConfig()` - No test for configuration loading
- `getConfigPath()` - No test for config path resolution

#### facts.go
- `collectFacts()` - No test for fact collection process
- `getFactsPath()` - No test for facts path resolution

#### global.go
- `setupLogging()` - No test for logging setup
- `setupFactsManager()` - No test for facts manager initialization

#### machines.go
- `connectToMachine()` - No test for machine connection
- `executeOnMachine()` - No test for remote execution
- `getMachineFacts()` - No test for machine fact retrieval

#### template_context.go
- `buildTemplateContext()` - No test for template context building
- `getMachineFactsForTemplate()` - No test for template-specific fact retrieval

#### template_engine.go
- `renderTemplate()` - No test for template rendering
- `validateTemplate()` - No test for template validation

#### templates.go
- `listProjectTemplates()` - No test for project template listing
- `renderProjectTemplate()` - No test for project template rendering

### 2. internal/ssh

#### client.go
- `NewClient()` - No test for SSH client creation
- `Connect()` - No test for SSH connection establishment
- `Close()` - No test for SSH connection cleanup
- `ExecuteCommand()` - No test for command execution

#### executor.go
- `NewExecutor()` - No test for executor creation
- `Execute()` - No test for command execution
- `ExecuteWithTimeout()` - No test for timeout-based execution

#### template_executor.go
- `NewTemplateExecutor()` - No test for template executor creation
- `ExecuteTemplate()` - No test for template-based execution
- `RenderTemplate()` - No test for template rendering

#### utils.go
- `parseSSHConfig()` - No test for SSH config parsing
- `getSSHKeyPath()` - No test for SSH key path resolution

### 3. internal/config

#### config.go
- `LoadConfig()` - No test for configuration loading
- `ValidateConfig()` - No test for configuration validation
- `GetDefaultConfig()` - No test for default configuration

#### machines.go
- `LoadMachines()` - No test for machine configuration loading
- `ValidateMachines()` - No test for machine validation
- `GetMachineByHost()` - No test for machine lookup

#### parser.go
- `ParseConfig()` - No test for configuration parsing
- `ParseMachines()` - No test for machine parsing
- `ParseActions()` - No test for action parsing

#### validator.go
- `ValidateProjectConfig()` - No test for project validation
- `ValidateMachineConfig()` - No test for machine validation
- `ValidateActionConfig()` - No test for action validation

### 4. internal/facts

#### manager.go
- `NewManager()` - No test for manager creation
- `AddFact()` - No test for fact addition
- `GetFact()` - No test for fact retrieval
- `DeleteFact()` - No test for fact deletion
- `ListFacts()` - No test for fact listing
- `CollectFacts()` - No test for fact collection
- `ExportFacts()` - No test for fact export
- `ImportFacts()` - No test for fact import

#### storage.go
- `NewStorage()` - No test for storage creation
- `Open()` - No test for storage opening
- `Close()` - No test for storage closing
- `Get()` - No test for storage retrieval
- `Set()` - No test for storage setting
- `Delete()` - No test for storage deletion

#### collector_utils.go
- `ValidateFactKey()` - No test for fact key validation
- `NormalizeFactKey()` - No test for fact key normalization
- `MergeFacts()` - No test for fact merging

### 5. internal/logging

#### logger.go
- `NewLogger()` - No test for logger creation
- `SetLevel()` - No test for log level setting
- `Log()` - No test for logging functionality

## Priority Levels

### High Priority
These tests are critical for core functionality:
- All `internal/facts` tests (data integrity)
- All `internal/config` tests (configuration reliability)
- All `internal/ssh` tests (remote execution safety)

### Medium Priority
These tests are important for CLI functionality:
- All `internal/cli` tests (user interface reliability)
- All `internal/logging` tests (debugging support)

### Low Priority
These tests are nice to have:
- Integration tests for end-to-end workflows
- Performance tests for large datasets
- Stress tests for concurrent operations

## Implementation Guidelines

### Test Structure
Each test should follow this pattern:
```go
func TestFunctionName(t *testing.T) {
    // Setup
    // Execute
    // Assert
}
```

### Test Coverage Goals
- **Line Coverage**: Aim for 90%+ line coverage
- **Branch Coverage**: Aim for 80%+ branch coverage
- **Function Coverage**: 100% of public functions should have tests

### Test Categories
1. **Happy Path Tests**: Normal operation scenarios
2. **Error Path Tests**: Error conditions and edge cases
3. **Boundary Tests**: Input validation and limits
4. **Integration Tests**: Multi-component interactions

### Mocking Strategy
- Use interfaces for external dependencies
- Mock SSH connections for testing
- Use temporary files for storage tests
- Use in-memory storage for fast tests

## Implementation Plan

### Phase 1: Core Infrastructure (Week 1)
- Implement `internal/facts` tests
- Implement `internal/config` tests
- Set up test utilities and mocks

### Phase 2: SSH and Execution (Week 2)
- Implement `internal/ssh` tests
- Add integration tests for remote execution
- Test timeout and error handling

### Phase 3: CLI Interface (Week 3)
- Implement `internal/cli` tests
- Test command-line argument parsing
- Test user interaction flows

### Phase 4: Logging and Utilities (Week 4)
- Implement `internal/logging` tests
- Add performance benchmarks
- Complete integration test suite

## Success Metrics

- [ ] 90%+ line coverage across all packages
- [ ] 100% of public functions have corresponding tests
- [ ] All tests pass consistently
- [ ] No flaky tests in CI/CD pipeline
- [ ] Test execution time under 30 seconds

## Notes

- Some functions currently return "not yet implemented" errors - these should be implemented before testing
- Consider using table-driven tests for functions with multiple input scenarios
- Ensure tests are deterministic and don't depend on external state
- Add benchmarks for performance-critical functions 