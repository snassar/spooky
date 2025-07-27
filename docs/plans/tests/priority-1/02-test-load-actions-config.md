# Test: LoadActionsConfig

## Priority: 1 (Critical) - Test #3

## Overview

**Function**: `LoadActionsConfig`  
**Current Coverage**: 0%  
**Target Coverage**: 90%+  
**Coverage Impact**: High - Complex file loading and merging logic

## Test Details

**Test Name**: `TestLoadActionsConfig`  
**Purpose**: Test loading actions from multiple sources  
**File**: `internal/config/parser_test.go`

## Test Cases

### 1. Load from Root actions.hcl Only
**Description**: Test loading actions from root actions.hcl file only  
**Expected**: Successfully load actions from root file  
**Coverage**: Single source loading path

**Test Project**: `examples/testing/test-only-actions-hcl`
- Has `actions.hcl` in root
- Has `actions/` directory with additional files
- Tests single source loading by using only the root file

### 2. Load from actions/ Directory Only
**Description**: Test loading actions from actions/ directory only  
**Expected**: Successfully load actions from directory files  
**Coverage**: Directory loading path

**Test Project**: `examples/testing/test-missing-actions`
- Has no `actions.hcl` in root
- Has `actions/` directory with multiple files
- Tests directory-only loading scenario

### 3. Load from Both Sources
**Description**: Test loading actions from both root file and actions/ directory  
**Expected**: Successfully merge actions from both sources  
**Coverage**: Multi-source merging path

**Test Project**: `examples/testing/test-valid-project`
- Has `actions.hcl` in root
- Has `actions/` directory with multiple files
- Tests merging from both sources

### 4. Merge Conflicts Resolution
**Description**: Test handling of action name conflicts between sources  
**Expected**: Handle conflicts appropriately (last wins or error)  
**Coverage**: Conflict resolution path

**Test Project**: `examples/testing/test-duplicate-actions`
- Has duplicate action names in root `actions.hcl`
- Tests conflict resolution behavior
- Contains "check-status" action defined twice

### 5. Invalid Files in actions/ Directory
**Description**: Test handling of invalid HCL files in actions/ directory  
**Expected**: Skip invalid files and continue with valid ones  
**Coverage**: Error handling path

**Test Project**: `examples/testing/test-invalid-actions`
- Contains invalid action definitions
- Tests error handling for malformed action files
- Verifies graceful handling of invalid HCL

### 6. No Actions Files Found
**Description**: Test when no actions files exist  
**Expected**: Return empty ActionsConfig  
**Coverage**: Empty result path

**Test Project**: `examples/testing/test-missing-project-file`
- Has no `actions.hcl` in root
- Has no `actions/` directory
- Tests empty result scenario

## Implementation

```go
func TestLoadActionsConfig(t *testing.T) {
    t.Run("RootActionsFileOnly", func(t *testing.T) {
        // Use test-only-actions-hcl project
        projectPath := filepath.Join("..", "..", "examples", "testing", "test-only-actions-hcl")
        // Call LoadActionsConfig
        // Assert successful loading
        // Verify actions from root file only
    })

    t.Run("ActionsDirectoryOnly", func(t *testing.T) {
        // Use test-missing-actions project
        projectPath := filepath.Join("..", "..", "examples", "testing", "test-missing-actions")
        // Call LoadActionsConfig
        // Assert successful loading
        // Verify actions from directory files only
    })

    t.Run("BothSources", func(t *testing.T) {
        // Use test-valid-project
        projectPath := filepath.Join("..", "..", "examples", "testing", "test-valid-project")
        // Call LoadActionsConfig
        // Assert successful loading
        // Verify actions merged from both sources
    })

    t.Run("MergeConflicts", func(t *testing.T) {
        // Use test-duplicate-actions project
        projectPath := filepath.Join("..", "..", "examples", "testing", "test-duplicate-actions")
        // Call LoadActionsConfig
        // Assert conflict resolution behavior
        // Verify final action set
    })

    t.Run("InvalidFilesInDirectory", func(t *testing.T) {
        // Use test-invalid-actions project
        projectPath := filepath.Join("..", "..", "examples", "testing", "test-invalid-actions")
        // Call LoadActionsConfig
        // Assert successful loading despite invalid files
        // Verify only valid actions loaded
    })

    t.Run("NoActionsFiles", func(t *testing.T) {
        // Use test-missing-project-file project
        projectPath := filepath.Join("..", "..", "examples", "testing", "test-missing-project-file")
        // Call LoadActionsConfig
        // Assert successful loading
        // Verify empty ActionsConfig returned
    })
}
```

## Coverage Goals

- **Function Coverage**: 100% of `LoadActionsConfig` function
- **Branch Coverage**: All file existence checks and error paths
- **Line Coverage**: All lines including file system operations

## Dependencies

- `testing` package
- `os` package for temporary file and directory creation
- `path/filepath` for path handling
- `github.com/stretchr/testify/assert` for assertions
- `io/ioutil` for file operations

## Notes

- Use existing test projects from `examples/testing/` directory
- No temporary file creation needed
- Test file permissions and access issues using existing projects
- Verify action merging order and precedence
- Test with various file sizes and complexities using real projects
- No cleanup needed since using existing files

## Success Criteria

- [ ] All test cases pass
- [ ] Function coverage reaches 100%
- [ ] Actions are correctly merged from multiple sources
- [ ] Invalid files are handled gracefully
- [ ] No temporary files left behind
- [ ] Test execution time < 2 seconds 