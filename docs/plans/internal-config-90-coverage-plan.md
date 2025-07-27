# Internal Config 90% Coverage Plan

## Overview

This document outlines the comprehensive test plan to increase coverage of the `internal/config` package from the current **69.1%** to **90%+**. The plan is based on detailed analysis of uncovered functions and code paths.

## Current Coverage Status

| Package | Current Coverage | Target | Gap |
|---------|------------------|--------|-----|
| **internal/config** | **54.6%** | **90%+** | **35.4%+** |

### Coverage Analysis by Function

Based on `go tool cover -func` output:

#### Zero Coverage Functions (Critical Priority)
- `ParseProjectConfigWithDebug` - 0% (Debug-enabled parsing)
- `LoadActionsConfig` - 0% (Multi-source action loading)
- `validateSSHKeyFile` - 0% (SSH key file validation)
- `validateScriptFile` - 0% (Script file validation)
- `SetMachineDefaults` - 0% (Machine default value application)
- `SetActionDefaults` - 0% (Action default value application)
- `SetInventoryDefaults` - 0% (Inventory default value application)
- `SetActionsDefaults` - 0% (Actions default value application)
- `resolveMachinePaths` - 0% (Machine path resolution)
- `resolveActionPaths` - 0% (Action path resolution)
- `validateWrapperBlocks` - 0% (Wrapper block validation)
- `formatMinValidation` - 0% (Min validation error formatting)

#### Low Coverage Functions
- `resolvePath` - 71.4% (Path resolution logic)
- `registerCustomValidations` - 66.7% (Validation registration)
- `parseInventoryWithWrapper` - 16.7% (Inventory wrapper parsing)
- `parseActionsWithWrapper` - 16.7% (Actions wrapper parsing)
- `parseConfigWithWrapper` - 34.8% (Generic wrapper parsing)

## Test Implementation Plan

### **Priority 1: Zero Coverage Functions (Critical)**

#### 1. ParseProjectConfigWithDebug Function (0% coverage)
**Test**: `TestParseProjectConfigWithDebug`
**Purpose**: Test debug-enabled project parsing
**Coverage Impact**: High - Debug path resolution and logging

**Test Cases**:
- Debug enabled with valid config
- Debug disabled with valid config
- Debug enabled with invalid config
- Path resolution with debug logging

#### 2. LoadActionsConfig Function (0% coverage)
**Test**: `TestLoadActionsConfig`
**Purpose**: Test loading actions from multiple sources
**Coverage Impact**: High - Complex file loading and merging logic

**Test Cases**:
- Load from root actions.hcl only
- Load from actions/ directory only
- Load from both sources
- Merge conflicts resolution
- Invalid files in actions/ directory
- No actions files found

#### 3. validateSSHKeyFile Function (0% coverage)
**Test**: `TestValidateSSHKeyFile`
**Purpose**: Test SSH key file validation
**Coverage Impact**: Medium - File existence and readability checks

**Test Cases**:
- Valid SSH key file
- Non-existent file
- Unreadable file
- Empty key file (should pass)
- Directory instead of file

#### 4. validateScriptFile Function (0% coverage)
**Test**: `TestValidateScriptFile`
**Purpose**: Test script file validation
**Coverage Impact**: Medium - File existence and executable checks

**Test Cases**:
- Valid executable script
- Non-existent file
- Non-executable file
- Empty script file (should pass)
- Directory instead of file

#### 5. SetMachineDefaults Function (0% coverage)
**Test**: `TestSetMachineDefaults`
**Purpose**: Test machine default value application
**Coverage Impact**: Medium - Default value logic

**Test Cases**:
- Machine with no port (should set default)
- Machine with existing port (should not change)
- Machine with other defaults
- Nil machine handling

#### 6. SetActionDefaults Function (0% coverage)
**Test**: `TestSetActionDefaults`
**Purpose**: Test action default value application
**Coverage Impact**: Medium - Default value logic

**Test Cases**:
- Action with no timeout (should set default)
- Action with existing timeout (should not change)
- Action with other defaults
- Nil action handling

#### 7. SetInventoryDefaults Function (0% coverage)
**Test**: `TestSetInventoryDefaults`
**Purpose**: Test inventory default value application
**Coverage Impact**: Medium - Default value logic

**Test Cases**:
- Inventory with multiple machines
- Empty inventory
- Nil inventory handling

#### 8. SetActionsDefaults Function (0% coverage)
**Test**: `TestSetActionsDefaults`
**Purpose**: Test actions default value application
**Coverage Impact**: Medium - Default value logic

**Test Cases**:
- Actions with multiple actions
- Empty actions
- Nil actions handling

#### 9. resolveMachinePaths Function (0% coverage)
**Test**: `TestResolveMachinePaths`
**Purpose**: Test machine path resolution
**Coverage Impact**: Medium - Path resolution logic

**Test Cases**:
- Key file path resolution
- Script path resolution
- Both paths present
- No paths present
- Invalid paths

#### 10. resolveActionPaths Function (0% coverage)
**Test**: `TestResolveActionPaths`
**Purpose**: Test action path resolution
**Coverage Impact**: Medium - Path resolution logic

**Test Cases**:
- Script path resolution
- Template path resolution
- Both paths present
- No paths present
- Invalid paths

#### 11. validateWrapperBlocks Function (0% coverage)
**Test**: `TestValidateWrapperBlocks`
**Purpose**: Test wrapper block validation
**Coverage Impact**: Medium - Block validation logic

**Test Cases**:
- Single valid block
- Multiple blocks (should fail)
- No blocks (should fail)
- Invalid block structure

#### 12. formatMinValidation Function (0% coverage)
**Test**: `TestFormatMinValidation`
**Purpose**: Test min validation error formatting
**Coverage Impact**: Low - Error message formatting

**Test Cases**:
- Min validation errors
- Max validation errors
- Required validation errors
- Custom validation errors

### **Priority 2: Low Coverage Functions**

#### 13. resolvePath Function (71.4% coverage)
**Test**: `TestResolvePath`
**Purpose**: Test path resolution logic
**Coverage Impact**: Medium - Relative to absolute path conversion

**Test Cases**:
- Absolute paths (no change)
- Relative paths (resolve correctly)
- Empty paths
- Paths with special characters
- Debug mode path resolution

#### 14. registerCustomValidations Function (66.7% coverage)
**Test**: `TestRegisterCustomValidations`
**Purpose**: Test custom validation registration
**Coverage Impact**: Low - Validation tag registration

**Test Cases**:
- SSH key file validation registration
- Script file validation registration
- Machine auth validation registration
- Action exec validation registration

#### 15. parseInventoryWithWrapper Function (16.7% coverage)
**Test**: `TestParseInventoryWithWrapper`
**Purpose**: Test inventory wrapper parsing
**Coverage Impact**: Medium - Inventory parsing logic

**Test Cases**:
- Valid inventory wrapper
- Invalid inventory wrapper
- Missing inventory block
- Multiple inventory blocks

#### 16. parseActionsWithWrapper Function (16.7% coverage)
**Test**: `TestParseActionsWithWrapper`
**Purpose**: Test actions wrapper parsing
**Coverage Impact**: Medium - Actions parsing logic

**Test Cases**:
- Valid actions wrapper
- Invalid actions wrapper
- Missing actions block
- Multiple actions blocks

#### 17. parseConfigWithWrapper Function (34.8% coverage)
**Test**: `TestParseConfigWithWrapper`
**Purpose**: Test generic wrapper parsing
**Coverage Impact**: Medium - Generic parsing logic

**Test Cases**:
- Valid generic wrapper
- Invalid generic wrapper
- Error handling scenarios
- Different wrapper types

### **Priority 3: Edge Cases and Error Handling**

#### 18. Error Handling Tests
**Test**: `TestParserErrorHandling`
**Purpose**: Test comprehensive error scenarios
**Coverage Impact**: High - All error paths in parser functions

**Test Cases**:
- HCL parsing errors
- File read errors
- Decoding errors
- Validation errors
- Path resolution errors

#### 19. Path Resolution Edge Cases
**Test**: `TestPathResolutionEdgeCases`
**Purpose**: Test complex path resolution scenarios
**Coverage Impact**: Medium - Edge cases in path handling

**Test Cases**:
- Paths with `../` components
- Paths with `./` components
- Absolute paths on different OS
- Paths with spaces and special characters
- Symlink resolution

#### 20. File System Integration Tests
**Test**: `TestFileSystemIntegration`
**Purpose**: Test actual file system operations
**Coverage Impact**: High - Real file operations for validation

**Test Cases**:
- Create temporary SSH key files
- Create temporary script files
- Test file permissions
- Test file readability
- Test file executability

### **Priority 4: Performance and Concurrency**

#### 21. Large Configuration Tests
**Test**: `TestLargeConfigurationParsing`
**Purpose**: Test performance with large configs
**Coverage Impact**: Medium - Memory usage and parsing speed

**Test Cases**:
- 1000+ machines
- 100+ actions
- Large tag sets
- Deep nested structures

#### 22. Concurrent Parsing Tests
**Test**: `TestConcurrentParsing`
**Purpose**: Test thread safety
**Coverage Impact**: Low - Concurrent access to parser functions

**Test Cases**:
- Multiple goroutines parsing same file
- Multiple goroutines parsing different files
- Race condition detection

### **Priority 5: Integration Tests**

#### 23. End-to-End Configuration Tests
**Test**: `TestEndToEndConfiguration`
**Purpose**: Test complete configuration workflow
**Coverage Impact**: High - Full configuration loading pipeline

**Test Cases**:
- Load project → inventory → actions
- Validate complete configuration
- Resolve all paths
- Apply defaults
- Final validation

#### 24. Configuration Migration Tests
**Test**: `TestConfigurationMigration`
**Purpose**: Test legacy to new format migration
**Coverage Impact**: Medium - Backward compatibility

**Test Cases**:
- Legacy format parsing
- Format conversion
- Validation after conversion

## Coverage Impact Estimation

| Function Category | Current Coverage | Target Coverage | Tests Needed | Coverage Impact |
|------------------|------------------|-----------------|--------------|-----------------|
| Zero Coverage Functions | 0% | 90%+ | 12 major tests | **+25-30%** |
| Low Coverage Functions | 16-71% | 90%+ | 5 tests | **+8-12%** |
| Edge Cases | N/A | 90%+ | 3 tests | **+3-5%** |
| Performance | N/A | 90%+ | 2 tests | **+1-2%** |
| Integration | N/A | 90%+ | 2 tests | **+2-3%** |

**Total Tests Needed**: **24 comprehensive tests**
**Expected Coverage Increase**: From 54.6% to **90%+**

## Implementation Strategy

### Phase 1: Critical Functions (Week 1-2)
1. Implement tests for zero coverage functions (12 tests)
2. Focus on parser functions, validation functions, and default functions
3. **Expected Coverage**: 70-75%

### Phase 2: Low Coverage Functions (Week 3)
1. Implement tests for low coverage functions (5 tests)
2. Focus on wrapper parsing and path resolution functions
3. **Expected Coverage**: 80-85%

### Phase 3: Edge Cases (Week 4)
1. Implement error handling and edge case tests (3 tests)
2. Focus on file system integration and path resolution edge cases
3. **Expected Coverage**: 85-90%

### Phase 4: Performance & Integration (Week 5)
1. Implement performance and integration tests (4 tests)
2. Focus on large configs and end-to-end workflows
3. **Expected Coverage**: 90%+

## Test File Organization

### Recommended Structure

```
internal/config/
├── config_test.go          # Parser and main configuration tests
├── validator_test.go       # Validation function tests  
├── machines_test.go        # Machine-related tests (already well covered)
└── integration_test.go     # End-to-end and performance tests
```

### Test Naming Convention

- `TestParseProjectConfigWithDebug_ValidConfig`
- `TestLoadActionsConfig_MultipleSources`
- `TestValidateSSHKeyFile_ValidFile`
- `TestResolvePath_RelativePaths`
- `TestErrorHandling_HCLParsingErrors`

## Success Metrics

### Coverage Targets
- **Individual Functions**: 90%+ coverage
- **Overall Package**: 90%+ coverage
- **Critical Paths**: 100% coverage

### Quality Metrics
- **Test Execution Time**: < 30 seconds
- **Test Reliability**: 100% pass rate
- **Error Scenarios**: All error paths covered
- **Edge Cases**: Comprehensive edge case coverage

## Risk Mitigation

### Potential Challenges
1. **File System Dependencies**: Use temporary files and cleanup
2. **OS-Specific Paths**: Test on multiple operating systems
3. **Concurrent Access**: Use proper synchronization
4. **Large Configs**: Monitor memory usage and performance

### Mitigation Strategies
1. **Temporary File Management**: Use `testing.T.TempDir()`
2. **Cross-Platform Testing**: Use path/filepath for OS compatibility
3. **Race Detection**: Use `go test -race`
4. **Memory Profiling**: Use `go test -memprofile`

## Conclusion

This comprehensive test plan will systematically increase the `internal/config` package coverage from 54.6% to 90%+ through:

1. **24 targeted tests** covering all uncovered functions
2. **Phased implementation** ensuring steady progress
3. **Quality-focused approach** with comprehensive error handling
4. **Performance considerations** for large configurations

The plan prioritizes critical functions first, ensuring maximum coverage impact with each test implementation. By following this structured approach, we can achieve the 90% coverage target while maintaining high test quality and reliability.

---

**Document Version**: 2.0  
**Created**: 2025-07-27  
**Updated**: 2025-01-27  
**Target Completion**: 5 weeks  
**Expected Coverage**: 90%+ 