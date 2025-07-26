# test-missing-actions

**Purpose**: Test handling of missing actions.hcl and actions/ directory.

Project without any action definitions to test CLI behavior when no actions are available.

**Expected Behavior**: spooky list-actions should fail with 'no actions found' error.
