# test-broken-symlinks

**Purpose**: Test error handling for broken symbolic links.

Contains broken-link.hcl pointing to nonexistent file to test broken symlink handling.

**Expected Behavior**: CLI should fail gracefully when encountering broken symlinks.
