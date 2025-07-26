# test-non-utf8-files

**Purpose**: Test handling of files with invalid UTF-8 encoding.

Contains invalid-encoding.hcl with non-UTF8 content to test encoding error handling.

**Expected Behavior**: Parser should fail with encoding error when reading non-UTF8 files.
