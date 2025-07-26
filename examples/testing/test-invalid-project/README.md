# test-invalid-project

**Purpose**: Test HCL syntax validation with malformed project.hcl.

Contains a project.hcl file with invalid HCL syntax (missing closing brace) to test:
- HCL parser error handling
- Validation error reporting
- Graceful failure when project configuration is malformed

**Expected Behavior**: `spooky validate` should fail with syntax error.
