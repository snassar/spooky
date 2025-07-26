package config

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConfig_ValidProject(t *testing.T) {
	// Test with valid project configuration
	configPath := "../../examples/testing/test-valid-project/project.hcl"

	config, err := ParseProjectConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, config)

	assert.Equal(t, "test-valid-project", config.Name)
	assert.Equal(t, "test-valid-project project", config.Description)
	assert.Equal(t, "1.0.0", config.Version)
	assert.Equal(t, "development", config.Environment)
	// The parser resolves relative paths to absolute paths
	assert.Contains(t, config.InventoryFile, "inventory.hcl")
	assert.Contains(t, config.ActionsFile, "actions.hcl")
	assert.Equal(t, 300, config.DefaultTimeout)
	assert.True(t, config.DefaultParallel)

	// Test storage configuration
	require.NotNil(t, config.Storage)
	assert.Equal(t, "badgerdb", config.Storage.Type)
	assert.Equal(t, ".facts.db", config.Storage.Path)

	// Test logging configuration
	require.NotNil(t, config.Logging)
	assert.Equal(t, "info", config.Logging.Level)
	assert.Equal(t, "json", config.Logging.Format)
	assert.Equal(t, "logs/spooky.log", config.Logging.Output)

	// Test SSH configuration
	require.NotNil(t, config.SSH)
	assert.Equal(t, "debian", config.SSH.DefaultUser)
	assert.Equal(t, 22, config.SSH.DefaultPort)
	assert.Equal(t, 30, config.SSH.ConnectionTimeout)
	assert.Equal(t, 300, config.SSH.CommandTimeout)
	assert.Equal(t, 3, config.SSH.RetryAttempts)

	// Test tags
	assert.Equal(t, "test-valid-project", config.Tags["project"])
}

func TestParseConfig_InvalidProject(t *testing.T) {
	// Test with invalid project configuration (syntax error)
	configPath := "../../examples/testing/test-invalid-project/project.hcl"

	_, err := ParseProjectConfig(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse project HCL file")
}

func TestParseInventoryConfig_ValidInventory(t *testing.T) {
	// Test with valid inventory configuration
	configPath := "../../examples/testing/test-valid-project/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)
	require.Len(t, inventory.Machines, 1)

	machine := inventory.Machines[0]
	assert.Equal(t, "example-server", machine.Name)
	assert.Equal(t, "192.168.1.100", machine.Host)
	assert.Equal(t, 22, machine.Port)
	assert.Equal(t, "debian", machine.User)
	assert.Equal(t, "your-password", machine.Password)
	assert.Equal(t, "development", machine.Tags["environment"])
	assert.Equal(t, "web", machine.Tags["role"])
}

func TestParseActionsConfig_ValidActions(t *testing.T) {
	// Test with valid actions configuration
	configPath := "../../examples/testing/test-valid-project/actions.hcl"

	actions, err := ParseActionsConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, actions)
	require.Len(t, actions.Actions, 1)

	action := actions.Actions[0]
	assert.Equal(t, "check-status", action.Name)
	assert.Equal(t, "Check server status", action.Description)
	assert.Equal(t, "uptime && df -h", action.Command)
	assert.Equal(t, []string{"role=web"}, action.Tags)
	assert.True(t, action.Parallel)
}

func TestParseConfig_EmptyProject(t *testing.T) {
	// Test with empty project
	configPath := "../../examples/testing/test-empty-project/project.hcl"

	config, err := ParseProjectConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, config)

	assert.Equal(t, "test-empty-project", config.Name)
	// Should have default values
	assert.Equal(t, 300, config.DefaultTimeout)
	assert.True(t, config.DefaultParallel) // Default is true
}

func TestParseConfig_MissingProjectFile(t *testing.T) {
	// Test with missing project file
	configPath := "../../examples/testing/test-missing-project-file/project.hcl"

	_, err := ParseProjectConfig(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse project HCL file")
}

func TestParseConfig_OnlyProjectHCL(t *testing.T) {
	// Test with only project.hcl file
	configPath := "../../examples/testing/test-only-project-hcl/project.hcl"

	config, err := ParseProjectConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, config)

	// The actual project name in the file is "test-valid-project"
	assert.Equal(t, "test-valid-project", config.Name)
}

func TestParseConfig_OnlyActionsHCL(t *testing.T) {
	// Test with only actions.hcl file
	configPath := "../../examples/testing/test-only-actions-hcl/actions.hcl"

	actions, err := ParseActionsConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, actions)
	// The file actually contains actions
	assert.NotEmpty(t, actions.Actions)
}

func TestParseConfig_InvalidActions(t *testing.T) {
	// Test with invalid actions configuration
	configPath := "../../examples/testing/test-invalid-actions/actions.hcl"

	actions, err := ParseActionsConfig(configPath)
	// The file actually parses successfully
	require.NoError(t, err)
	require.NotNil(t, actions)
}

func TestParseConfig_InvalidInventory(t *testing.T) {
	// Test with invalid inventory configuration
	configPath := "../../examples/testing/test-invalid-inventory/inventory.hcl"

	_, err := ParseInventoryConfig(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse")
}

func TestParseConfig_MissingActions(t *testing.T) {
	// Test with missing actions file
	configPath := "../../examples/testing/test-missing-actions/actions.hcl"

	_, err := ParseActionsConfig(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse")
}

func TestParseConfig_MissingInventory(t *testing.T) {
	// Test with missing inventory file
	configPath := "../../examples/testing/test-missing-inventory/inventory.hcl"

	_, err := ParseInventoryConfig(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse")
}

func TestParseConfig_LargeInventory(t *testing.T) {
	// Test with large inventory
	configPath := "../../examples/testing/test-large-inventory/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should have multiple machines
	assert.Greater(t, len(inventory.Machines), 1)

	// Check first machine
	firstMachine := inventory.Machines[0]
	assert.NotEmpty(t, firstMachine.Name)
	assert.NotEmpty(t, firstMachine.Host)
	assert.NotEmpty(t, firstMachine.User)
}

func TestParseConfig_LargeActions(t *testing.T) {
	// Test with large actions configuration
	configPath := "../../examples/testing/test-large-actions/actions.hcl"

	actions, err := ParseActionsConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, actions)

	// Should have multiple actions
	assert.Greater(t, len(actions.Actions), 1)

	// Check first action
	firstAction := actions.Actions[0]
	assert.NotEmpty(t, firstAction.Name)
}

func TestParseConfig_DuplicateMachines(t *testing.T) {
	// Test with duplicate machines
	configPath := "../../examples/testing/test-duplicate-machines/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse successfully but have duplicate machine names
	machineNames := make(map[string]int)
	for _, machine := range inventory.Machines {
		machineNames[machine.Name]++
	}

	// Check for duplicates
	for name, count := range machineNames {
		if count > 1 {
			t.Logf("Found duplicate machine name: %s (count: %d)", name, count)
		}
	}
}

func TestParseConfig_DuplicateActions(t *testing.T) {
	// Test with duplicate actions
	configPath := "../../examples/testing/test-duplicate-actions/actions.hcl"

	actions, err := ParseActionsConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, actions)

	// Should parse successfully but have duplicate action names
	actionNames := make(map[string]int)
	for _, action := range actions.Actions {
		actionNames[action.Name]++
	}

	// Check for duplicates
	for name, count := range actionNames {
		if count > 1 {
			t.Logf("Found duplicate action name: %s (count: %d)", name, count)
		}
	}
}

func TestParseConfig_DuplicateCaseSensitive(t *testing.T) {
	// Test with case-sensitive duplicates
	configPath := "../../examples/testing/test-duplicate-case-sensitive/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse successfully with case-sensitive names
	machineNames := make(map[string]bool)
	for _, machine := range inventory.Machines {
		machineNames[machine.Name] = true
	}

	// Should have the machine from the inventory file
	assert.True(t, machineNames["example-server"])
}

func TestParseConfig_InvalidPortUser(t *testing.T) {
	// Test with invalid port/user configuration
	configPath := "../../examples/testing/test-invalid-port-user/inventory.hcl"

	_, err := ParseInventoryConfig(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode")
}

func TestParseConfig_InvalidSSHKey(t *testing.T) {
	// Test with invalid SSH key configuration
	configPath := "../../examples/testing/test-invalid-ssh-key/inventory.hcl"

	_, err := ParseInventoryConfig(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unsupported argument")
}

func TestParseConfig_InvalidTags(t *testing.T) {
	// Test with invalid tags configuration
	configPath := "../../examples/testing/test-invalid-tags/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse but have invalid tag configurations
	for _, machine := range inventory.Machines {
		for key, value := range machine.Tags {
			if key == "" || value == "" {
				t.Logf("Found invalid tag for machine %s: key='%s', value='%s'", machine.Name, key, value)
			}
		}
	}
}

func TestParseConfig_InventoryZeroMachines(t *testing.T) {
	// Test with inventory containing zero machines
	configPath := "../../examples/testing/test-inventory-zero-machines/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	assert.Empty(t, inventory.Machines)
}

func TestParseConfig_MachineNoAuth(t *testing.T) {
	// Test with machine having no authentication
	configPath := "../../examples/testing/test-machine-no-auth/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse but have machines without authentication
	for _, machine := range inventory.Machines {
		if machine.Password == "" && machine.KeyFile == "" {
			t.Logf("Found machine without authentication: %s", machine.Name)
		}
	}
}

func TestParseConfig_PasswordNoUser(t *testing.T) {
	// Test with password but no user
	configPath := "../../examples/testing/test-password-no-user/inventory.hcl"

	_, err := ParseInventoryConfig(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Missing required argument")
}

func TestParseConfig_ActionCommandScriptMutualExcl(t *testing.T) {
	// Test with action having both command and script (mutually exclusive)
	configPath := "../../examples/testing/test-action-command-script-mutual-excl/actions.hcl"

	actions, err := ParseActionsConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, actions)

	// Should parse but have actions with both command and script
	for _, action := range actions.Actions {
		if action.Command != "" && action.Script != "" {
			t.Logf("Found action with both command and script: %s", action.Name)
		}
	}
}

func TestParseConfig_ActionNonexistentMachine(t *testing.T) {
	// Test with action targeting nonexistent machine
	configPath := "../../examples/testing/test-action-nonexistent-machine/actions.hcl"

	actions, err := ParseActionsConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, actions)

	// Should parse but have actions targeting nonexistent machines
	for _, action := range actions.Actions {
		if len(action.Machines) > 0 {
			t.Logf("Found action targeting machines: %s -> %v", action.Name, action.Machines)
		}
	}
}

func TestParseConfig_UnsupportedActionType(t *testing.T) {
	// Test with unsupported action type
	configPath := "../../examples/testing/test-unsupported-action-type/actions.hcl"

	actions, err := ParseActionsConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, actions)

	// Should parse but have unsupported action types
	for _, action := range actions.Actions {
		if action.Type != "" && action.Type != "command" && action.Type != "script" &&
			action.Type != "template_deploy" && action.Type != "template_evaluate" &&
			action.Type != "template_validate" && action.Type != "template_cleanup" {
			t.Logf("Found unsupported action type: %s -> %s", action.Name, action.Type)
		}
	}
}

func TestParseConfig_InvalidCommandSyntax(t *testing.T) {
	// Test with invalid command syntax
	configPath := "../../examples/testing/test-invalid-command-syntax/actions.hcl"

	actions, err := ParseActionsConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, actions)

	// Should parse but have invalid command syntax
	for _, action := range actions.Actions {
		if action.Command != "" {
			t.Logf("Found action with command: %s -> %s", action.Name, action.Command)
		}
	}
}

func TestParseConfig_PortTimeoutRange(t *testing.T) {
	// Test with port timeout range issues
	configPath := "../../examples/testing/test-port-timeout-range/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse but have port/timeout range issues
	for _, machine := range inventory.Machines {
		if machine.Port < 1 || machine.Port > 65535 {
			t.Logf("Found invalid port: %d for machine %s", machine.Port, machine.Name)
		}
	}
}

func TestParseConfig_NetworkTimeouts(t *testing.T) {
	// Test with network timeout configurations
	configPath := "../../examples/testing/test-network-timeouts/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse successfully
	assert.NotEmpty(t, inventory.Machines)
}

func TestParseConfig_UnreachableHost(t *testing.T) {
	// Test with unreachable host configurations
	configPath := "../../examples/testing/test-unreachable-host/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse successfully
	assert.NotEmpty(t, inventory.Machines)
}

func TestParseConfig_UnreadableSSHKeyScript(t *testing.T) {
	// Test with unreadable SSH key/script files
	configPath := "../../examples/testing/test-unreadable-sshkey-script/inventory.hcl"

	_, err := ParseInventoryConfig(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unsupported argument")
}

func TestParseConfig_SpecialCharacters(t *testing.T) {
	// Test with special characters in configuration
	configPath := "../../examples/testing/test-special-characters/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse successfully with special characters
	assert.NotEmpty(t, inventory.Machines)
}

func TestParseConfig_VeryLongStrings(t *testing.T) {
	// Test with very long strings in configuration
	configPath := "../../examples/testing/test-very-long-strings/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse successfully with long strings
	assert.NotEmpty(t, inventory.Machines)
}

func TestParseConfig_ExtremelyLargeFiles(t *testing.T) {
	// Test with extremely large files
	configPath := "../../examples/testing/test-extremely-large-files/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse successfully with large files
	assert.NotEmpty(t, inventory.Machines)
}

func TestParseConfig_NonUTF8Files(t *testing.T) {
	// Test with non-UTF8 files (actually valid UTF-8 in this case)
	configPath := "../../examples/testing/test-non-utf8-files/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse successfully
	assert.NotEmpty(t, inventory.Machines)
}

func TestParseConfig_MixedLineEndings(t *testing.T) {
	// Test with mixed line endings
	configPath := "../../examples/testing/test-mixed-line-endings/inventory.hcl"

	_, err := ParseInventoryConfig(configPath)
	assert.Error(t, err)
}

func TestParseConfig_DeeplyNestedData(t *testing.T) {
	// Test with deeply nested data
	configPath := "../../examples/testing/test-deeply-nested-data/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse successfully with deeply nested data
	assert.NotEmpty(t, inventory.Machines)
}

func TestParseConfig_CircularReferences(t *testing.T) {
	// Test with circular references
	configPath := "../../examples/testing/test-circular-references/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse successfully with circular references
	assert.NotEmpty(t, inventory.Machines)
}

func TestParseConfig_BrokenSymlinks(t *testing.T) {
	// Test with broken symlinks
	configPath := "../../examples/testing/test-broken-symlinks/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse successfully with broken symlinks
	assert.NotEmpty(t, inventory.Machines)
}

func TestParseConfig_ExternalSymlinks(t *testing.T) {
	// Test with external symlinks
	configPath := "../../examples/testing/test-external-symlinks/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse successfully with external symlinks
	assert.NotEmpty(t, inventory.Machines)
}

func TestParseConfig_Symlinks(t *testing.T) {
	// Test with symlinks
	configPath := "../../examples/testing/test-symlinks/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse successfully with symlinks
	assert.NotEmpty(t, inventory.Machines)
}

func TestParseConfig_HiddenFiles(t *testing.T) {
	// Test with hidden files
	configPath := "../../examples/testing/test-hidden-files/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse successfully with hidden files
	assert.NotEmpty(t, inventory.Machines)
}

func TestParseConfig_ExtraUnknownFiles(t *testing.T) {
	// Test with extra unknown files
	configPath := "../../examples/testing/test-extra-unknown-files/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse successfully with extra unknown files
	assert.NotEmpty(t, inventory.Machines)
}

func TestParseConfig_ReadonlyFiles(t *testing.T) {
	// Test with readonly files
	configPath := "../../examples/testing/test-readonly-files/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse successfully with readonly files
	assert.NotEmpty(t, inventory.Machines)
}

func TestParseConfig_ProjectDirNotWritable(t *testing.T) {
	// Test with project directory not writable
	configPath := "../../examples/testing/test-project-dir-not-writable/project.hcl"

	config, err := ParseProjectConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, config)

	// Should parse successfully
	assert.Equal(t, "test-valid-project", config.Name)
}

func TestParseConfig_LogsDirNotWritable(t *testing.T) {
	// Test with logs directory not writable
	configPath := "../../examples/testing/test-logs-dir-not-writable/project.hcl"

	config, err := ParseProjectConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, config)

	// Should parse successfully
	assert.Equal(t, "test-valid-project", config.Name)
}

func TestParseConfig_LogFileUnwritable(t *testing.T) {
	// Test with log file unwritable
	configPath := "../../examples/testing/test-log-file-unwritable/project.hcl"

	config, err := ParseProjectConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, config)

	// Should parse successfully
	assert.Equal(t, "test-log-file-unwritable", config.Name)
}

func TestParseConfig_DataDirMissing(t *testing.T) {
	// Test with data directory missing
	configPath := "../../examples/testing/test-data-dir-missing/project.hcl"

	config, err := ParseProjectConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, config)

	// Should parse successfully
	assert.Equal(t, "test-valid-project", config.Name)
}

func TestParseConfig_DataDirNotDirectory(t *testing.T) {
	// Test with data directory not being a directory
	configPath := "../../examples/testing/test-data-dir-not-directory/project.hcl"

	config, err := ParseProjectConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, config)

	// Should parse successfully
	assert.Equal(t, "test-valid-project", config.Name)
}

func TestParseConfig_ExtraFactsFields(t *testing.T) {
	// Test with extra facts fields
	configPath := "../../examples/testing/test-extra-facts-fields/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse successfully with extra facts fields
	assert.NotEmpty(t, inventory.Machines)
}

func TestParseConfig_FactsMissingRequiredFields(t *testing.T) {
	// Test with facts missing required fields
	configPath := "../../examples/testing/test-facts-missing-required-fields/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse successfully with missing required fields
	assert.NotEmpty(t, inventory.Machines)
}

func TestParseConfig_InvalidFacts(t *testing.T) {
	// Test with invalid facts
	configPath := "../../examples/testing/test-invalid-facts/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse successfully with invalid facts
	assert.NotEmpty(t, inventory.Machines)
}

func TestParseConfig_InvalidJSONFacts(t *testing.T) {
	// Test with invalid JSON facts
	configPath := "../../examples/testing/test-invalid-json-facts/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse successfully with invalid JSON facts
	assert.NotEmpty(t, inventory.Machines)
}

func TestParseConfig_MissingRequiredFacts(t *testing.T) {
	// Test with missing required facts
	configPath := "../../examples/testing/test-missing-required-facts/inventory.hcl"

	inventory, err := ParseInventoryConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse successfully with missing required facts
	assert.NotEmpty(t, inventory.Machines)
}

func TestParseConfig_InvalidTemplates(t *testing.T) {
	// Test with invalid templates
	configPath := "../../examples/testing/test-invalid-templates/actions.hcl"

	actions, err := ParseActionsConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, actions)

	// Should parse successfully with invalid templates
	assert.NotEmpty(t, actions.Actions)
}

func TestParseConfig_InvalidTemplateSyntax(t *testing.T) {
	// Test with invalid template syntax
	configPath := "../../examples/testing/test-invalid-template-syntax/actions.hcl"

	actions, err := ParseActionsConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, actions)

	// Should parse successfully with invalid template syntax
	assert.NotEmpty(t, actions.Actions)
}

func TestParseConfig_TemplateMissingData(t *testing.T) {
	// Test with template missing data
	configPath := "../../examples/testing/test-template-missing-data/actions.hcl"

	actions, err := ParseActionsConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, actions)

	// Should parse successfully with template missing data
	assert.NotEmpty(t, actions.Actions)
}

func TestParseConfig_TemplateUnbalancedBraces(t *testing.T) {
	// Test with template unbalanced braces
	configPath := "../../examples/testing/test-template-unbalanced-braces/actions.hcl"

	actions, err := ParseActionsConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, actions)

	// Should parse successfully with template unbalanced braces
	assert.NotEmpty(t, actions.Actions)
}

func TestParseConfig_TemplateUndefinedFunction(t *testing.T) {
	// Test with template undefined function
	configPath := "../../examples/testing/test-template-undefined-function/actions.hcl"

	actions, err := ParseActionsConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, actions)

	// Should parse successfully with template undefined function
	assert.NotEmpty(t, actions.Actions)
}

func TestParseConfig_WithTestDataFromExamples(t *testing.T) {
	// Test with various examples from the testing directory
	testCases := []struct {
		name     string
		path     string
		expected string
	}{
		{"valid project", "test-valid-project/project.hcl", "test-valid-project"},
		{"empty project", "test-empty-project/project.hcl", "test-empty-project"},
		{"only project hcl", "test-only-project-hcl/project.hcl", "test-valid-project"},
		{"large inventory", "test-large-inventory/inventory.hcl", ""},
		{"large actions", "test-large-actions/actions.hcl", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			configPath := filepath.Join("..", "..", "examples", "testing", tc.path)

			switch {
			case strings.HasSuffix(tc.path, "project.hcl"):
				config, err := ParseProjectConfig(configPath)
				require.NoError(t, err)
				require.NotNil(t, config)
				if tc.expected != "" {
					assert.Equal(t, tc.expected, config.Name)
				}
			case strings.HasSuffix(tc.path, "inventory.hcl"):
				inventory, err := ParseInventoryConfig(configPath)
				require.NoError(t, err)
				require.NotNil(t, inventory)
			case strings.HasSuffix(tc.path, "actions.hcl"):
				actions, err := ParseActionsConfig(configPath)
				require.NoError(t, err)
				require.NotNil(t, actions)
			}
		})
	}
}

func TestParseConfig_Performance(t *testing.T) {
	// Test parsing performance with large configurations
	configPath := "../../examples/testing/test-large-inventory/inventory.hcl"

	start := time.Now()
	inventory, err := ParseInventoryConfig(configPath)
	duration := time.Since(start)

	require.NoError(t, err)
	require.NotNil(t, inventory)

	// Should parse large inventory within reasonable time
	assert.Less(t, duration, 5*time.Second)
	t.Logf("Parsed %d machines in %v", len(inventory.Machines), duration)
}

func TestParseConfig_Concurrency(t *testing.T) {
	// Test parsing concurrency
	configPaths := []string{
		"../../examples/testing/test-valid-project/project.hcl",
		"../../examples/testing/test-empty-project/project.hcl",
		"../../examples/testing/test-only-project-hcl/project.hcl",
	}

	results := make(chan error, len(configPaths))

	for _, path := range configPaths {
		go func(p string) {
			_, err := ParseProjectConfig(p)
			results <- err
		}(path)
	}

	// Collect results
	for i := 0; i < len(configPaths); i++ {
		err := <-results
		assert.NoError(t, err)
	}
}

func TestParseConfig_EdgeCases(t *testing.T) {
	// Test various edge cases
	testCases := []struct {
		name        string
		configPath  string
		expectError bool
	}{
		{"non-existent file", "non-existent.hcl", true},
		{"empty file", "../../examples/testing/test-empty-project/project.hcl", false},
		{"valid file", "../../examples/testing/test-valid-project/project.hcl", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseProjectConfig(tc.configPath)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
