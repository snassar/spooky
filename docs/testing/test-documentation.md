# Spooky Test Documentation

This document provides a comprehensive overview of all tests in the spooky project, organized by package and test type.

## Test Structure Overview

The spooky project follows a comprehensive testing strategy with multiple layers:

- **Unit Tests**: Co-located with source files (`*_test.go`)
- **Integration Tests**: Located in `tests/integration/`
- **Example-based Tests**: Located in `examples/testing/`
- **Coverage Requirements**: Defined in `tests/testcoverage.yml`

## Test Coverage Requirements

**File**: `tests/testcoverage.yml`

- **Individual Files**: Minimum 60% coverage
- **Packages**: Minimum 70% coverage  
- **Overall Project**: Minimum 75% coverage
- **Exclusions**: Test files, examples, main.go, generated files

## Internal Package Tests

### CLI Package (`internal/cli/`)

**Test Files**:
- `commands_test.go` (32KB, 1232 lines)
- `global_test.go` (11KB, 455 lines)
- `templates_test.go` (6.3KB, 305 lines)
- `facts_test.go` (32KB, 1218 lines)

**Test Categories**:

#### Command Tests (commands_test.go)
- `TestInitCmd` - Project initialization command
- `TestValidateCmd` - Configuration validation command
- `TestListCmd` - List command functionality
- `TestListMachinesCmd` - Machine listing
- `TestListActionsCmd` - Action listing
- `TestListTemplatesCmd` - Template listing
- `TestListFactsCmd` - Facts listing
- `TestGatherFactsCmd` - Facts gathering
- `TestRenderTemplateCmd` - Template rendering
- `TestValidateTemplateCmd` - Template validation

#### Project-specific Tests
- `TestInitProject` - Project initialization
- `TestValidateProject` - Project validation
- `TestListProject` - Project listing
- `TestListProjectMachines` - Project machines
- `TestListProjectActions` - Project actions
- `TestListProjectTemplates` - Project templates
- `TestListProjectFacts` - Project facts
- `TestGatherProjectFacts` - Project facts gathering
- `TestRenderProjectTemplate` - Project template rendering
- `TestValidateProjectTemplate` - Project template validation

#### Facts Command Tests (facts_test.go)
- `TestFactsCmd` - Facts command functionality
- `TestFactsGatherCmd` - Facts gathering
- `TestFactsImportCmd` - Facts import
- `TestFactsExportCmd` - Facts export
- `TestFactsValidateCmd` - Facts validation
- `TestFactsQueryCmd` - Facts querying
- `TestFactsCacheCmd` - Facts cache operations

#### Facts Execution Tests
- `TestRunFactsGather` - Facts gathering execution
- `TestRunFactsImport` - Facts import execution
- `TestRunFactsExport` - Facts export execution
- `TestRunFactsValidate` - Facts validation execution
- `TestRunFactsQuery` - Facts query execution
- `TestRunFactsCacheClear` - Cache clearing
- `TestRunFactsCacheExpired` - Cache expiration

#### Utility Tests
- `TestDetermineTargetHosts` - Host target determination
- `TestCollectFactsFromHosts` - Multi-host fact collection
- `TestCollectFactsFromHost` - Single host fact collection
- `TestDisplayFactGatheringResults` - Results display
- `TestDisplayHostFacts` - Host facts display
- `TestGetTotalFactCount` - Fact counting
- `TestIsCustomSource` - Custom source detection
- `TestParseQueryExpression` - Query parsing

#### Global Configuration Tests (global_test.go)
- `TestGetGlobalConfig` - Global configuration retrieval
- `TestAddGlobalFlags` - Global flag addition
- `TestGetEnvOrDefault` - Environment variable handling
- `TestGetDefaultLogFile` - Log file defaults
- `TestGetFactsDBPath` - Facts database path
- `TestGlobalFlagsIntegration` - Flag integration
- `TestGlobalConfigSerialization` - Config serialization
- `TestEnvironmentVariablePrecedence` - Environment precedence
- `TestFlagValidation` - Flag validation
- `TestPathResolution` - Path resolution
- `TestXDGStateHomeHandling` - XDG state home
- `TestUserHomeDirectoryFallback` - Home directory fallback
- `TestFlagConflictHandling` - Flag conflicts
- `TestLogFileDirectoryCreation` - Log directory creation
- `TestFactsDBPathValidation` - Facts DB path validation
- `TestTemplatePathResolution` - Template path resolution
- `TestDataPathResolution` - Data path resolution

#### Template Tests (templates_test.go)
- `TestTemplateValidation` - Template validation
- `TestTemplateRendering` - Template rendering
- `TestTemplateOutputHandling` - Template output handling

### Config Package (`internal/config/`)

**Test Files**:
- `config_test.go` (27KB, 899 lines)
- `machines_test.go` (16KB, 708 lines)
- `validator_test.go` (12KB, 501 lines)

**Test Categories**:

#### Configuration Parsing Tests (config_test.go)
- `TestParseInventoryConfig_ValidInventory` - Valid inventory parsing
- `TestParseActionsConfig_ValidActions` - Valid actions parsing

#### Machine Management Tests (machines_test.go)
- `TestBuildEnterpriseIndex` - Enterprise index building
- `TestGetMachinesForActionLarge` - Large-scale machine lookup
- `TestGetMachinesForActionLarge_InvalidMachines` - Invalid machine handling
- `TestIndexCache_GetIndex` - Index cache operations
- `TestIndexCache_GetIndexMetrics` - Index cache metrics
- `TestFindMachinesByName` - Machine name lookup
- `TestGetAllMachines` - All machines retrieval
- `TestGetMachinesForAction` - Action-specific machine lookup
- `TestGetMachinesForAction_WithTestDataFromExamples` - Example-based machine lookup
- `TestBuildEnterpriseIndex_Performance` - Index building performance
- `TestGetMachinesForActionLarge_Performance` - Large lookup performance
- `TestIndexCache_Concurrency` - Cache concurrency
- `TestSortTagsByPopularity` - Tag popularity sorting
- `TestUpdateLookupMetrics` - Lookup metrics updates
- `TestComputeConfigHash` - Configuration hashing
- `TestIndexCache_IsValid` - Cache validation
- `TestBuildEnterpriseIndex_EdgeCases` - Index edge cases
- `TestGetMachinesForActionLarge_EdgeCases` - Large lookup edge cases

#### Validation Tests (validator_test.go)
- `TestNewValidator` - Validator creation
- `TestValidateConfig_ValidConfig` - Valid configuration validation
- `TestValidateConfig_EmptyMachines` - Empty machines validation
- `TestValidateConfig_NilConfig` - Nil configuration handling
- `TestValidateMachine_ValidMachine` - Valid machine validation
- `TestValidateMachine_InvalidMachine` - Invalid machine validation
- `TestValidateAction_ValidAction` - Valid action validation
- `TestValidateAction_InvalidAction` - Invalid action validation
- `TestValidateSSHKeyFile_ValidFile` - Valid SSH key validation
- `TestValidateSSHKeyFile_NonExistentFile` - Non-existent SSH key handling
- `TestValidateScriptFile_ValidFile` - Valid script validation
- `TestValidateScriptFile_NonExecutableFile` - Non-executable script handling
- `TestValidateScriptFile_NonExistentFile` - Non-existent script handling
- `TestValidateConfig_WithTestDataFromExamples` - Example-based validation
- `TestValidateConfig_EdgeCases` - Validation edge cases
- `TestValidateConfig_Performance` - Validation performance

### Facts Package (`internal/facts/`)

**Test Files**:
- `manager_test.go` (12KB, 456 lines)
- `storage_test.go` (12KB, 446 lines)
- `validation_test.go` (7.8KB, 302 lines)
- `types_test.go` (12KB, 413 lines)

**Test Categories**:

#### Manager Tests (manager_test.go)
- `TestNewManager` - Manager creation
- `TestNewManagerWithStorage` - Manager with storage
- `TestManagerConfigureHCLCollector` - HCL collector configuration
- `TestManagerConfigureOpenTofuCollector` - OpenTofu collector configuration
- `TestManagerCollectAllFacts` - All facts collection
- `TestManagerCollectSpecificFacts` - Specific facts collection
- `TestManagerGetFact` - Fact retrieval
- `TestManagerCacheOperations` - Cache operations
- `TestManagerWithTestDataFromExamples` - Example-based manager tests
- `TestManagerImportCustomFacts` - Custom facts import
- `TestManagerImportCustomFactsWithOptions` - Custom facts import with options
- `TestManagerGetCustomFacts` - Custom facts retrieval
- `TestManagerFactQuery` - Fact querying
- `TestManagerExportImportFacts` - Facts export/import
- `TestManagerWithMissingRequiredFacts` - Missing required facts
- `TestManagerWithInvalidJSONFacts` - Invalid JSON facts
- `TestManagerFactCollectionClone` - Fact collection cloning
- `TestManagerRegisterCustomCollector` - Custom collector registration
- `TestManagerClose` - Manager cleanup

#### Storage Tests (storage_test.go)
- `TestNewFactStorage` - Storage creation
- `TestConvertFactCollectionToMachineFacts` - Fact collection conversion
- `TestConvertMachineFactsToFactCollection` - Machine facts conversion
- `TestStorageWithTestDataFromExamples` - Example-based storage tests
- `TestStorageQueryOperations` - Storage query operations
- `TestStorageExportImport` - Storage export/import
- `TestStorageDeleteOperations` - Storage delete operations

#### Type Tests (types_test.go)
- `TestFactStruct` - Fact structure
- `TestFactCollectionStruct` - Fact collection structure
- `TestFactCollectionClone` - Fact collection cloning
- `TestFactCollectionCloneNil` - Nil fact collection cloning
- `TestSystemFactsStruct` - System facts structure
- `TestOSInfoStruct` - OS info structure
- `TestHardwareInfoStruct` - Hardware info structure
- `TestNetworkInfoStruct` - Network info structure
- `TestFactSourceConstants` - Fact source constants
- `TestMergePolicyConstants` - Merge policy constants
- `TestCustomFactsStruct` - Custom facts structure
- `TestImportOptionsStruct` - Import options structure
- `TestMergeModeConstants` - Merge mode constants
- `TestFactWithTestDataFromExamples` - Example-based fact tests

#### Validation Tests (validation_test.go)
- `TestValidationError` - Validation error handling
- `TestValidationResult` - Validation result handling
- `TestValidateCustomFacts` - Custom facts validation
- `TestValidateCustomFactsWithOverrides` - Custom facts with overrides
- `TestValidateWithTestDataFromExamples` - Example-based validation

### Logging Package (`internal/logging/`)

**Test Files**:
- `logger_test.go` (10.0KB, 455 lines)
- `types_test.go` (7.4KB, 325 lines)
- `fields_test.go` (6.3KB, 240 lines)

**Test Categories**:

#### Field Tests (fields_test.go)
- `TestStringField` - String field creation
- `TestIntField` - Integer field creation
- `TestInt64Field` - Int64 field creation
- `TestFloat64Field` - Float64 field creation
- `TestBoolField` - Boolean field creation
- `TestErrorField` - Error field creation
- `TestDurationField` - Duration field creation
- `TestRequestIDField` - Request ID field creation
- `TestServerField` - Server field creation
- `TestActionField` - Action field creation
- `TestHostField` - Host field creation
- `TestPortField` - Port field creation
- `TestFieldHelpersWithTestDataFromExamples` - Example-based field tests
- `TestFieldHelpersEdgeCases` - Field edge cases
- `TestFieldHelpersWithSpecialCharacters` - Special character handling
- `TestFieldHelpersPerformance` - Field performance
- `TestFieldHelpersConcurrency` - Field concurrency
- `TestFieldHelpersCombinedUsage` - Combined field usage

#### Logger Tests (logger_test.go)
- `TestNewLogger` - Logger creation
- `TestConfigureLogger` - Logger configuration
- `TestLoggerWithFields` - Logger with fields
- `TestLoggerWithContext` - Logger with context
- `TestGlobalLoggerFunctions` - Global logger functions
- `TestContextLoggerFunctions` - Context logger functions
- `TestEnsureLogDirectory` - Log directory creation
- `TestLoggerWithTestDataFromExamples` - Example-based logger tests
- `TestLoggerLevels` - Logger levels
- `TestLoggerErrorHandling` - Logger error handling
- `TestLoggerFieldConversion` - Field conversion
- `TestLoggerSync` - Logger synchronization
- `TestLoggerInvalidConfigurations` - Invalid configurations
- `TestLoggerPerformance` - Logger performance
- `TestLoggerConcurrency` - Logger concurrency

#### Type Tests (types_test.go)
- `TestLogLevelConstants` - Log level constants
- `TestFieldStruct` - Field structure
- `TestConfigStruct` - Config structure
- `TestLoggerInterface` - Logger interface
- `TestLoggerInterfaceChaining` - Logger interface chaining
- `TestConfigDefaultValues` - Config default values
- `TestConfigWithCustomValues` - Config with custom values
- `TestFieldWithVariousTypes` - Various field types
- `TestLogLevelComparison` - Log level comparison
- `TestConfigEquality` - Config equality
- `TestFieldEquality` - Field equality
- `TestLoggerInterfaceNilHandling` - Nil handling
- `TestLoggerInterfaceWithTestDataFromExamples` - Example-based interface tests
- `TestConfigValidation` - Config validation

### SSH Package (`internal/ssh/`)

**Test Files**:
- `client_test.go` (14KB, 533 lines)
- `executor_test.go` (9.7KB, 447 lines)
- `template_executor_test.go` (15KB, 564 lines)
- `utils_test.go` (6.6KB, 241 lines)

**Test Categories**:

#### Client Tests (client_test.go)
- `TestNewSSHClient` - SSH client creation
- `TestNewSSHClientWithHostKeyCallback` - SSH client with host key callback
- `TestSSHClient_GetMachine` - Machine retrieval
- `TestSSHClient_ConnectAndClose` - Connection and cleanup
- `TestSSHClient_ExecuteCommand_NotConnected` - Command execution without connection
- `TestSSHClient_ExecuteScript_NotConnected` - Script execution without connection
- `TestGetHostKeyCallback` - Host key callback creation
- `TestGetHostKeyCallback_WithHomeDirectory` - Host key callback with home directory
- `TestSSHClient_WithTestDataFromExamples` - Example-based client tests
- `TestSSHClient_TimeoutConfiguration` - Timeout configuration
- `TestSSHClient_InvalidTimeout` - Invalid timeout handling
- `TestSSHClient_AuthenticationMethods` - Authentication methods

#### Executor Tests (executor_test.go)
- `TestExecuteConfig_NilConfig` - Nil configuration handling
- `TestExecuteConfig_EmptyActions` - Empty actions handling
- `TestExecuteConfig_EmptyMachines` - Empty machines handling
- `TestExecuteConfig_UnsupportedActionType` - Unsupported action types
- `TestExecuteConfig_CommandAction` - Command action execution
- `TestExecuteConfig_ScriptAction` - Script action execution
- `TestExecuteConfig_TemplateAction` - Template action execution
- `TestExecuteConfig_ParallelExecution` - Parallel execution
- `TestExecuteConfig_WithTestDataFromExamples` - Example-based execution tests
- `TestExecuteConfig_ActionTypes` - Different action types
- `TestIsTemplateAction` - Template action detection
- `TestExecuteConfig_MachineFiltering` - Machine filtering
- `TestExecuteConfig_TimeoutConfiguration` - Timeout configuration

#### Template Executor Tests (template_executor_test.go)
- `TestNewTemplateActionExecutor` - Template executor creation
- `TestTemplateActionExecutor_ExecuteAction_NilAction` - Nil action handling
- `TestTemplateActionExecutor_ExecuteAction_NilTemplate` - Nil template handling
- `TestTemplateActionExecutor_ExecuteAction_UnsupportedType` - Unsupported types
- `TestTemplateActionExecutor_ValidateTemplateSyntax` - Template syntax validation
- `TestTemplateActionExecutor_ValidateTemplateSyntax_WithTestExamples` - Example-based syntax validation
- `TestTemplateActionExecutor_CreateFuncMap` - Function map creation
- `TestTemplateActionExecutor_CreateFuncMapWithValues` - Function map with values
- `TestTemplateActionExecutor_ExecuteTemplateDeploy_FileNotExists` - Template deploy with missing file
- `TestTemplateActionExecutor_ExecuteTemplateDeploy_ValidTemplate` - Valid template deployment
- `TestTemplateActionExecutor_ExecuteTemplateDeploy_InvalidTemplate` - Invalid template deployment
- `TestTemplateActionExecutor_ExecuteTemplateOperation` - Template operations
- `TestTemplateActionExecutor_ExecuteTemplateOperation_WithFailingSSH` - Template operations with SSH failures
- `TestTemplateActionExecutor_WithTestDataFromExamples` - Example-based template tests
- `TestTemplateActionExecutor_TemplateActionTypes` - Template action types

#### Utility Tests (utils_test.go)
- `TestSSHUtils_BasicFunctionality` - Basic SSH utility functionality
- `TestSSHUtils_WithTestDataFromExamples` - Example-based utility tests
- `TestSSHUtils_ErrorHandling` - Error handling
- `TestSSHUtils_ConfigurationValidation` - Configuration validation
- `TestSSHUtils_NetworkAddressValidation` - Network address validation
- `TestSSHUtils_PortValidation` - Port validation
- `TestSSHUtils_TimeoutValidation` - Timeout validation
- `TestSSHUtils_AuthenticationValidation` - Authentication validation
- `TestSSHUtils_FilePathValidation` - File path validation

## Integration Tests (`tests/integration/`)

**Files**:
- `README.md` (5.5KB, 199 lines) - Integration test documentation
- `test-config.hcl` (661B, 30 lines) - Test configuration

**Planned Test Files** (referenced in README but not yet implemented):
- `podman_ci_test.go` - CI-specific tests for GitHub Actions
- `podman_integration_test.go` - Local integration tests for development
- `podman_basic_test.go` - Basic Podman environment tests

**Test Categories** (as documented in README):

#### Basic Environment Tests
- Podman availability and functionality
- Basic container operations

#### Local Integration Tests
- Full integration tests for local development
- Requires `-podman` flag to run

#### CI Integration Tests
- Tests designed for GitHub Actions CI environment
- Automatically skipped when not in CI

**Test Coverage** (as documented):
1. SSH Connection
2. Configuration Processing
3. Script Execution
4. Parallel Execution
5. Error Handling
6. CLI Commands

## Example-based Tests (`examples/testing/`)

The project includes extensive example-based testing scenarios in the `examples/testing/` directory. These are not traditional Go test files but rather configuration examples that test various edge cases and error conditions:

### Test Scenarios Include:
- Action command/script mutual exclusion
- Non-existent machine references
- Broken symlinks
- Circular references
- Conflicting permissions
- Corrupted facts database
- Missing data directories
- Duplicate actions/machines
- Invalid configurations
- Network timeouts
- Special characters
- Large files and strings
- Non-UTF8 files
- Template validation
- And many more...

## Test Execution

### Running Unit Tests
```bash
# Run all internal package tests
go test ./internal/... -v

# Run specific package tests
go test ./internal/cli -v
go test ./internal/config -v
go test ./internal/facts -v
go test ./internal/logging -v
go test ./internal/ssh -v
```

### Running Integration Tests
```bash
# Run basic Podman tests
go test ./tests/integration -v -podman-basic

# Run full integration tests (requires Podman)
go test ./tests/integration -v -podman

# Run all integration tests
go test ./tests/integration -v -podman -podman-basic
```

### Test Coverage
```bash
# Generate coverage report
go test ./internal/... -coverprofile=coverage.out

# View coverage in browser
go tool cover -html=coverage.out
```

## Test Quality Metrics

- **Total Test Files**: 15+ test files across internal packages
- **Total Test Functions**: 200+ individual test functions
- **Coverage Requirements**: 75% overall project coverage
- **Test Types**: Unit tests, integration tests, example-based tests
- **Test Categories**: Functionality, validation, error handling, performance, concurrency, edge cases

## Test Maintenance

### Adding New Tests
1. Follow existing naming conventions
2. Use appropriate test categories
3. Include example-based tests for edge cases
4. Ensure proper cleanup in test functions
5. Add documentation for new test scenarios

### Test Organization
- Unit tests are co-located with source files
- Integration tests are in dedicated `tests/` directory
- Example-based tests use real configuration scenarios
- Performance tests are separated from functional tests

### Continuous Integration
- Tests run automatically in GitHub Actions
- Coverage requirements are enforced
- Integration tests use Podman containers
- Example-based tests validate real-world scenarios 