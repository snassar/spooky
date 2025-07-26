# Test Scenarios Documentation

This document provides a comprehensive overview of all test scenarios created in `./examples/testing/` for unit and integration testing of the spooky CLI and backend logic.

## Overview

We have created **59 test scenarios** covering various edge cases, error conditions, and validation scenarios. These scenarios are designed to ensure robust error handling and full coverage of edge conditions across the spooky CLI and backend logic.

## Test Scenario Categories

### 1. Baseline Scenarios
| Scenario | Description | Purpose |
|----------|-------------|---------|
| `test-valid-project` | Fully valid project with all components | Positive test baseline |

### 2. Project Structure Edge Cases
| Scenario | Description | Purpose |
|----------|-------------|---------|
| `test-extra-unknown-files` | Unknown files in project directory | Test that CLI ignores/warns about unknown files without failing |
| `test-symlinks` | Symlinked config and template files | Test handling of symlinked configuration files |
| `test-readonly-files` | Read-only configuration files | Test error handling when trying to write to read-only files |
| `test-non-utf8-files` | Files with invalid UTF-8 encoding | Test that parser fails gracefully with encoding errors |
| `test-hidden-files` | Hidden files (.env, .DS_Store, etc.) | Test that CLI ignores hidden files appropriately |
| `test-external-symlinks` | Symlinks to external system files | Test handling of external symlinks |
| `test-broken-symlinks` | Broken symbolic links | Test error handling for broken symlinks |
| `test-duplicate-case-sensitive` | Files with same name, different case | Test case sensitivity handling |
| `test-conflicting-permissions` | Files with mixed permissions | Test permission handling |
| `test-mixed-line-endings` | Files with mixed line endings (CRLF/LF) | Test line ending handling |

### 3. Inventory/Machine Edge Cases
| Scenario | Description | Purpose |
|----------|-------------|---------|
| `test-duplicate-machines` | Duplicate machine names in inventory | Test validation error for duplicate machine names |
| `test-invalid-port-user` | Wrong data types and missing required fields | Test type validation and required field validation |
| `test-invalid-inventory` | Malformed inventory.hcl with syntax errors | Test HCL syntax validation |
| `test-missing-inventory` | Missing inventory.hcl file | Test handling of missing inventory file |
| `test-inventory-zero-machines` | Empty inventory block | Test minimum machine requirement validation |
| `test-invalid-tags` | Non-string tags and missing required tags | Test tag validation |
| `test-special-characters` | Special characters in names and values | Test parsing with special characters |
| `test-very-long-strings` | Very long strings in configuration | Test string length limits |

### 4. Actions Edge Cases
| Scenario | Description | Purpose |
|----------|-------------|---------|
| `test-duplicate-actions` | Duplicate action names | Test validation error for duplicate action names |
| `test-invalid-command-syntax` | Empty commands and invalid shell syntax | Test command validation |
| `test-unsupported-action-type` | Unknown action types | Test validation/execution error for unsupported types |
| `test-invalid-actions` | Malformed actions.hcl with syntax errors | Test HCL syntax validation |
| `test-missing-actions` | Missing actions.hcl and actions/ directory | Test handling of missing action files |
| `test-action-nonexistent-machine` | Action references nonexistent machine | Test valid_machines validation error |
| `test-action-command-script-mutual-excl` | Both/neither command and script | Test action_exec validation error |
| `test-port-timeout-range` | Invalid port (0, 70000) and timeout (0, 5000) | Test valid_port and valid_timeout validation |

### 5. Facts/Custom Facts Edge Cases
| Scenario | Description | Purpose |
|----------|-------------|---------|
| `test-invalid-json-facts` | Malformed JSON in facts files | Test JSON parsing error handling |
| `test-missing-required-facts` | Missing required fields in facts | Test schema validation for facts |
| `test-extra-facts-fields` | Unexpected fields in facts | Test that extra fields are ignored/warned about |
| `test-invalid-facts` | Invalid facts structure | Test facts validation error handling |
| `test-facts-missing-required-fields` | Missing MachineName, SystemID, OS | Test facts validation errors/warnings |
| `test-corrupted-facts-db` | Corrupted .facts.db file | Test facts database error handling |
| `test-deeply-nested-data` | Deeply nested JSON structures | Test parsing and memory handling |

### 6. Templates Edge Cases
| Scenario | Description | Purpose |
|----------|-------------|---------|
| `test-invalid-template-syntax` | Go template syntax errors | Test template validation |
| `test-template-missing-data` | Templates referencing missing variables | Test template rendering error handling |
| `test-template-undefined-function` | Templates using undefined functions | Test template function error handling |
| `test-template-unbalanced-braces` | Templates with unbalanced braces | Test template brace validation |

### 7. SSH/Connection Edge Cases
| Scenario | Description | Purpose |
|----------|-------------|---------|
| `test-unreachable-host` | Non-routable and invalid IP addresses | Test connection error handling |
| `test-invalid-ssh-key` | Non-existent SSH key files | Test SSH error handling |
| `test-password-no-user` | Password authentication without user field | Test validation error for missing user |
| `test-unreadable-sshkey-script` | Unreadable SSH keys and unexecutable scripts | Test file permission validation |
| `test-machine-no-auth` | Machine with neither password nor key_file | Test machine_auth validation error |
| `test-network-timeouts` | Very short SSH timeouts | Test timeout handling |

### 8. Logging/Output Edge Cases
| Scenario | Description | Purpose |
|----------|-------------|---------|
| `test-log-file-unwritable` | Unwritable log file paths | Test that logging errors are reported but don't crash the CLI |
| `test-logs-dir-not-writable` | Logs directory not writable | Test logging directory permission handling |

### 9. Data Directory Edge Cases
| Scenario | Description | Purpose |
|----------|-------------|---------|
| `test-data-dir-missing` | Missing data/ directory | Test that CLI creates it as needed or fails gracefully |
| `test-data-dir-not-directory` | data/ is a file, not a directory | Test error handling when data/ is not a directory |

### 10. Partial Configuration Scenarios
| Scenario | Description | Purpose |
|----------|-------------|---------|
| `test-only-project-hcl` | Only project.hcl present | Test validation error for incomplete configuration |
| `test-only-actions-hcl` | Only actions.hcl present | Test validation error for incomplete configuration |

### 11. Empty Project Scenarios
| Scenario | Description | Purpose |
|----------|-------------|---------|
| `test-empty-project` | No actions or inventory defined | Test handling of empty projects |

### 12. Performance/Large Scale Scenarios
| Scenario | Description | Purpose |
|----------|-------------|---------|
| `test-large-inventory` | 10+ machines in inventory | Test performance and memory handling with large inventories |
| `test-large-actions` | 10+ actions defined | Test performance and memory handling with large action sets |
| `test-extremely-large-files` | 10MB+ data files | Test large file handling |

### 13. Advanced Edge Cases
| Scenario | Description | Purpose |
|----------|-------------|---------|
| `test-circular-references` | Machines referencing each other in tags | Test circular reference detection |
| `test-project-dir-not-writable` | Project directory not writable | Test file creation error handling |
| `test-invalid-templates` | Invalid template configurations | Test template configuration validation |

## Usage in Tests

### Example Usage in Go Tests

```go
func TestInvalidProject(t *testing.T) {
    projectPath := "./examples/testing/test-invalid-project"
    
    // Test project validation
    err := validateProject(logger, projectPath)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "syntax error")
}

func TestMissingInventory(t *testing.T) {
    projectPath := "./examples/testing/test-missing-inventory"
    
    // Test inventory loading
    _, err := config.ParseInventoryConfig(filepath.Join(projectPath, "inventory.hcl"))
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "no such file")
}

func TestLargeInventoryPerformance(t *testing.T) {
    projectPath := "./examples/testing/test-large-inventory"
    
    start := time.Now()
    inventory, err := config.ParseInventoryConfig(filepath.Join(projectPath, "inventory.hcl"))
    duration := time.Since(start)
    
    assert.NoError(t, err)
    assert.Len(t, inventory.Machines, 11) // 10 + example-server
    assert.Less(t, duration, 100*time.Millisecond) // Performance assertion
}

func TestActionReferencesNonexistentMachine(t *testing.T) {
    projectPath := "./examples/testing/test-action-nonexistent-machine"
    
    // Test action validation
    _, err := config.LoadActionsConfig(projectPath)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "valid_machines")
}

func TestUnreadableSSHKey(t *testing.T) {
    projectPath := "./examples/testing/test-unreadable-sshkey-script"
    
    // Test SSH key validation
    inventory, err := config.ParseInventoryConfig(filepath.Join(projectPath, "inventory.hcl"))
    assert.NoError(t, err) // Should parse successfully
    // Validation would fail when trying to use the key
}

func TestPortTimeoutRangeValidation(t *testing.T) {
    projectPath := "./examples/testing/test-port-timeout-range"
    
    // Test port and timeout validation
    inventory, err := config.ParseInventoryConfig(filepath.Join(projectPath, "inventory.hcl"))
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "valid_port")
}
```

### CLI Testing Examples

```bash
# Test invalid project validation
./build/spooky validate ./examples/testing/test-invalid-project

# Test missing inventory
./build/spooky list-machines ./examples/testing/test-missing-inventory

# Test large inventory performance
./build/spooky list-machines ./examples/testing/test-large-inventory

# Test template validation
./build/spooky validate-template ./examples/testing/test-invalid-template-syntax/templates/invalid-syntax.tmpl

# Test action validation
./build/spooky validate ./examples/testing/test-action-nonexistent-machine

# Test SSH key validation
./build/spooky validate ./examples/testing/test-unreadable-sshkey-script

# Test network timeouts
./build/spooky gather-facts ./examples/testing/test-network-timeouts
```

## Test Coverage Goals

These scenarios are designed to achieve:

1. **100% Error Path Coverage**: Every error condition should be tested
2. **Edge Case Coverage**: All boundary conditions and edge cases
3. **Performance Testing**: Large scale scenarios for performance validation
4. **User Experience Testing**: Proper error messages and graceful failures
5. **Integration Testing**: End-to-end CLI command testing
6. **File System Testing**: Various file system edge cases
7. **Network Testing**: Connection and timeout scenarios
8. **Security Testing**: Permission and access control scenarios

## Maintenance

When adding new features or changing validation logic:

1. **Add new test scenarios** to cover new edge cases
2. **Update existing scenarios** if validation rules change
3. **Run all scenarios** to ensure they still produce expected results
4. **Document any changes** to scenario behavior

## Scenario Details

### File Structure Examples

Each test scenario follows this structure:
```
test-scenario-name/
├── project.hcl          # Project configuration (if applicable)
├── inventory.hcl        # Machine inventory (if applicable)
├── actions.hcl          # Action definitions (if applicable)
├── actions/             # Action subdirectory (if applicable)
│   ├── 01-dependencies.hcl
│   ├── 02-system-update.hcl
│   └── 03-monitoring.hcl
├── templates/           # Template files (if applicable)
├── data/               # Data files (if applicable)
├── files/              # File distribution (if applicable)
├── logs/               # Log files (if applicable)
└── README.md           # Scenario description
```

### Validation Expectations

Each scenario should produce predictable results:

- **Valid scenarios**: Should pass validation and execute successfully
- **Invalid scenarios**: Should fail with appropriate error messages
- **Performance scenarios**: Should complete within reasonable time limits
- **Error scenarios**: Should fail gracefully with clear error messages

## Conclusion

This comprehensive test suite provides a solid foundation for ensuring the reliability and robustness of the spooky CLI and backend logic. Regular testing against these scenarios will help catch regressions and ensure consistent behavior across different edge cases and error conditions.

The 59 test scenarios cover:
- **Syntax validation** (malformed HCL, JSON, templates)
- **Schema validation** (missing required fields, wrong types)
- **File system edge cases** (missing files, symlinks, permissions, hidden files)
- **Network/SSH edge cases** (unreachable hosts, invalid keys, timeouts)
- **Performance testing** (large inventories, actions, and files)
- **Error handling** (graceful failures, proper error messages)
- **Configuration edge cases** (partial configs, empty projects)
- **Advanced scenarios** (circular references, deeply nested data, special characters) 