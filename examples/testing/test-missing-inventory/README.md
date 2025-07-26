# test-missing-inventory

**Purpose**: Test handling of missing inventory.hcl file.

Project without inventory.hcl to test CLI behavior when machine definitions are missing.

**Expected Behavior**: spooky list-machines should fail with 'inventory.hcl not found' error.
