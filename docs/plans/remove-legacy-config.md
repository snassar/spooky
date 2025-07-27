# Remove Legacy Config Code

## Overview

The `ParseConfig` function and `Config` struct are legacy code that handles the old combined HCL format. The codebase has moved to a modular approach using separate files, making this legacy code obsolete.

## Current Status

### Functions to Remove
- `ParseConfig(filename string) (*Config, error)` - Legacy combined format parser
- `ValidateConfig(config *Config) error` - Only used by ParseConfig and tests

### Types to Remove
- `Config` struct - Legacy combined format structure

## Impact Analysis

### Files That Will Be Modified
1. **`internal/config/parser.go`**
   - Remove `ParseConfig` function (lines 55-120)
   - Remove `Config` struct usage

2. **`internal/config/validator.go`**
   - Remove `ValidateConfig` function
   - Remove `validateConfigStruct` function
   - Remove `validateConfig` method
   - Remove `Config` struct registration

3. **`internal/config/types.go`**
   - Remove `Config` struct definition

4. **`internal/config/config.go`**
   - Remove references to ParseConfig and ValidateConfig

5. **Test Files**
   - Update all tests to use new modular types
   - Remove tests that specifically test legacy format

### Files That Use Config Struct
- `internal/config/config_test.go` - 80+ test functions
- `internal/config/validator_test.go` - 10+ test functions  
- `internal/config/machines_test.go` - 15+ test functions
- `internal/ssh/executor_test.go` - 15+ test functions
- `internal/ssh/template_executor_test.go` - 8+ test functions
- `internal/facts/hcl_collector.go` - 2 usages

## Migration Strategy

### Phase 1: Update Tests
1. **Replace `Config` with `InventoryConfig` + `ActionsConfig`**
   - Update all test functions to use the new modular types
   - Create helper functions to build test data

2. **Update Validation Tests**
   - Replace `ValidateConfig` calls with individual validation functions
   - Use `ValidateMachine` and `ValidateAction` instead

### Phase 2: Update Internal Usage
1. **Update SSH Executor Tests**
   - Replace `config.Config` with separate inventory/actions configs
   - Update test setup functions

2. **Update Facts Collector**
   - Replace `config.Config` usage with appropriate new types

### Phase 3: Remove Legacy Code
1. **Remove Functions**
   - Delete `ParseConfig` function
   - Delete `ValidateConfig` function
   - Delete `validateConfigStruct` function

2. **Remove Types**
   - Delete `Config` struct
   - Remove struct registration in validator

3. **Clean Up Comments**
   - Remove references to legacy functions
   - Update documentation

## Benefits

1. **Reduced Code Complexity** - Remove unused legacy code
2. **Better Architecture** - Clear separation of concerns with modular files
3. **Improved Maintainability** - Single source of truth for each config type
4. **Better Testing** - Tests focus on actual functionality being used
5. **Reduced Confusion** - No more legacy vs new format confusion

## Risks

1. **Test Updates** - Large number of test files need updating
2. **Potential Bugs** - Risk of introducing bugs during migration
3. **Breaking Changes** - Any external code using these functions will break

## Implementation Plan

### Week 1: Test Migration
- [ ] Update `internal/config/config_test.go`
- [ ] Update `internal/config/validator_test.go`
- [ ] Update `internal/config/machines_test.go`

### Week 2: Internal Usage Migration
- [ ] Update `internal/ssh/executor_test.go`
- [ ] Update `internal/ssh/template_executor_test.go`
- [ ] Update `internal/facts/hcl_collector.go`

### Week 3: Legacy Code Removal
- [ ] Remove `ParseConfig` function
- [ ] Remove `ValidateConfig` function
- [ ] Remove `Config` struct
- [ ] Clean up comments and documentation

### Week 4: Verification
- [ ] Run all tests to ensure they pass
- [ ] Verify no regressions
- [ ] Update documentation

## Success Criteria

- [ ] All tests pass with new modular types
- [ ] No references to `ParseConfig` or `Config` remain
- [ ] Code coverage remains at or above current levels
- [ ] No functionality is lost
- [ ] Documentation is updated

## Notes

- This is a breaking change for any external code using these functions
- The migration should be done carefully to avoid regressions
- Consider adding deprecation warnings before removal if external usage is expected 