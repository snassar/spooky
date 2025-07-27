# Test: ParseProjectConfigWithDebug

## Priority: 1 (Critical) - Test #1

## Overview

**Function**: `ParseProjectConfigWithDebug`  
**Current Coverage**: 0%  
**Target Coverage**: 90%+  
**Coverage Impact**: High - Debug path resolution and logging

## Test Details

**Test Name**: `TestParseProjectConfigWithDebug`  
**Purpose**: Test debug-enabled project parsing  
**File**: `internal/config/parser_test.go`

## Test Cases

### 1. Debug enabled with valid config
**Description**: Test parsing with debug mode enabled and valid configuration  
**Expected**: Successfully parse and return valid ProjectConfig with debug logging  
**Coverage**: Debug path execution

**Test Data**:
```hcl
project {
  name = "test-project"
  description = "Test project with debug"
  version = "1.0.0"
  environment = "development"
  inventory_file = "inventory.hcl"
  actions_file = "actions.hcl"
}
```

### 2. Debug disabled with valid config
**Description**: Test parsing with debug mode disabled and valid configuration  
**Expected**: Successfully parse and return valid ProjectConfig without debug logging  
**Coverage**: Non-debug path execution

**Test Data**: Same as above

### 3. Debug enabled with invalid config
**Description**: Test parsing with debug mode enabled but invalid configuration  
**Expected**: Return error with debug information in logs  
**Coverage**: Debug error path

**Test Data**:
```hcl
project {
  name = "test-project"
  // Missing required fields
}
```

### 4. Path resolution with debug logging
**Description**: Test path resolution when debug mode is enabled  
**Expected**: Debug logs show path resolution steps  
**Coverage**: Debug path resolution

**Test Data**:
```hcl
project {
  name = "test-project"
  inventory_file = "../relative/inventory.hcl"
  actions_file = "./actions/actions.hcl"
}
```

## Implementation

```go
func TestParseProjectConfigWithDebug(t *testing.T) {
    t.Run("DebugEnabledValidConfig", func(t *testing.T) {
        // Create temporary valid project config file
        // Call ParseProjectConfigWithDebug with debug=true
        // Assert successful parsing
        // Verify debug logging occurred
        // Verify ProjectConfig struct fields
    })

    t.Run("DebugDisabledValidConfig", func(t *testing.T) {
        // Create temporary valid project config file
        // Call ParseProjectConfigWithDebug with debug=false
        // Assert successful parsing
        // Verify no debug logging occurred
        // Verify ProjectConfig struct fields
    })

    t.Run("DebugEnabledInvalidConfig", func(t *testing.T) {
        // Create temporary invalid project config file
        // Call ParseProjectConfigWithDebug with debug=true
        // Assert error returned
        // Verify debug logging occurred during error
        // Verify error message
    })

    t.Run("PathResolutionWithDebug", func(t *testing.T) {
        // Create temporary project config with relative paths
        // Call ParseProjectConfigWithDebug with debug=true
        // Assert successful parsing
        // Verify debug logs show path resolution
        // Verify paths are resolved correctly
    })
}
```

## Coverage Goals

- **Function Coverage**: 100% of `ParseProjectConfigWithDebug` function
- **Branch Coverage**: Both debug=true and debug=false paths
- **Line Coverage**: All lines in the function including debug logging

## Dependencies

- `testing` package
- `os` package for temporary file creation
- `path/filepath` for path handling
- `github.com/stretchr/testify/assert` for assertions
- Logging package for debug output verification

## Notes

- Use `testing.T.TempDir()` for temporary file creation
- Verify debug logging output using log capture or mock logger
- Test relative path resolution with debug enabled
- Ensure debug mode doesn't affect normal functionality
- Test both positive and negative debug scenarios
- **Updated**: Test cases now align with updated coverage plan (v2.0)
- **Updated**: Test file location changed from `config_test.go` to `parser_test.go`

## Success Criteria

- [x] All test cases pass
- [x] Function coverage reaches 100%
- [x] Debug logging works correctly when enabled
- [x] No debug logging when disabled
- [x] Path resolution works with debug mode
- [x] Test execution time < 1 second 