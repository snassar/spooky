# test-readonly-files

**Purpose**: Test error handling with read-only configuration files.

project.hcl set to read-only (chmod 444) to test file permission handling.

**Expected Behavior**: CLI should fail gracefully when trying to write to read-only files.
