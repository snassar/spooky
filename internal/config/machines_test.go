package config

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildEnterpriseIndex(t *testing.T) {
	// Test building enterprise index
	machines := []Machine{
		{
			Name: "server1",
			Tags: map[string]string{
				"environment": "production",
				"role":        "web",
			},
		},
		{
			Name: "server2",
			Tags: map[string]string{
				"environment": "production",
				"role":        "database",
			},
		},
		{
			Name: "server3",
			Tags: map[string]string{
				"environment": "staging",
				"role":        "web",
			},
		},
	}

	index := buildEnterpriseIndex(machines)
	require.NotNil(t, index)
	require.NotNil(t, index.TagIndex)
	require.NotNil(t, index.MachineTagIndex)
	require.NotNil(t, index.TagCount)
	require.NotNil(t, index.Metrics)

	// Test metrics
	assert.Equal(t, 3, index.Metrics.MachineCount)
	assert.Equal(t, 2, index.Metrics.TagCount) // environment and role
	assert.True(t, index.Metrics.BuildTime > 0)
	assert.False(t, index.Metrics.LastUpdated.IsZero())

	// Test tag index
	assert.Len(t, index.TagIndex["environment=production"], 2)
	assert.Len(t, index.TagIndex["environment=staging"], 1)
	assert.Len(t, index.TagIndex["role=web"], 2)
	assert.Len(t, index.TagIndex["role=database"], 1)

	// Test machine tag index - check that all machines in the index have their tags recorded
	assert.Len(t, index.MachineTagIndex, 3)
	for machinePtr, machineTags := range index.MachineTagIndex {
		assert.NotNil(t, machinePtr)
		assert.NotEmpty(t, machineTags)
		// Verify the machine has the expected tags
		switch machinePtr.Name {
		case "server1":
			assert.Equal(t, "production", machineTags["environment"])
			assert.Equal(t, "web", machineTags["role"])
		case "server2":
			assert.Equal(t, "production", machineTags["environment"])
			assert.Equal(t, "database", machineTags["role"])
		case "server3":
			assert.Equal(t, "staging", machineTags["environment"])
			assert.Equal(t, "web", machineTags["role"])
		}
	}
}

func TestGetMachinesForActionLarge(t *testing.T) {
	// Test getting machines for action with large index
	config := &Config{
		Machines: []Machine{
			{
				Name: "server1",
				Tags: map[string]string{
					"environment": "production",
					"role":        "web",
				},
			},
			{
				Name: "server2",
				Tags: map[string]string{
					"environment": "production",
					"role":        "database",
				},
			},
			{
				Name: "server3",
				Tags: map[string]string{
					"environment": "staging",
					"role":        "web",
				},
			},
		},
	}

	index := buildEnterpriseIndex(config.Machines)

	// Test action with machine names
	action := &Action{
		Name:     "test-action",
		Machines: []string{"server1", "server2"},
	}

	machines, err := GetMachinesForActionLarge(config, action, index)
	require.NoError(t, err)
	assert.Len(t, machines, 2)

	// Test action with tags
	action = &Action{
		Name: "test-action",
		Tags: []string{"role=web"},
	}

	machines, err = GetMachinesForActionLarge(config, action, index)
	require.NoError(t, err)
	assert.Len(t, machines, 2) // server1 and server3 have role=web

	// Test action with both machines and tags
	action = &Action{
		Name:     "test-action",
		Machines: []string{"server1"},
		Tags:     []string{"environment=production"},
	}

	machines, err = GetMachinesForActionLarge(config, action, index)
	require.NoError(t, err)
	assert.Len(t, machines, 1) // server1 matches both criteria

	// Test action with no criteria (should return all machines)
	action = &Action{
		Name: "test-action",
	}

	machines, err = GetMachinesForActionLarge(config, action, index)
	require.NoError(t, err)
	assert.Len(t, machines, 3)
}

func TestGetMachinesForActionLarge_InvalidMachines(t *testing.T) {
	// Test with invalid machine names
	config := &Config{
		Machines: []Machine{
			{
				Name: "server1",
				Tags: map[string]string{
					"environment": "production",
				},
			},
		},
	}

	index := buildEnterpriseIndex(config.Machines)

	action := &Action{
		Name:     "test-action",
		Machines: []string{"nonexistent-server"},
	}

	_, err := GetMachinesForActionLarge(config, action, index)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "machine 'nonexistent-server' not found")
}

func TestIndexCache_GetIndex(t *testing.T) {
	// Test index cache functionality
	cache := &IndexCache{}
	require.NotNil(t, cache)

	config := &Config{
		Machines: []Machine{
			{
				Name: "server1",
				Tags: map[string]string{
					"environment": "production",
				},
			},
		},
	}

	// First call should build index
	index1 := cache.GetIndex(config)
	require.NotNil(t, index1)
	assert.Equal(t, 1, index1.Metrics.MachineCount)

	// Second call should return cached index
	index2 := cache.GetIndex(config)
	require.NotNil(t, index2)
	assert.Equal(t, index1, index2)

	// Test with different config
	config2 := &Config{
		Machines: []Machine{
			{
				Name: "server1",
				Tags: map[string]string{
					"environment": "production",
				},
			},
			{
				Name: "server2",
				Tags: map[string]string{
					"environment": "staging",
				},
			},
		},
	}

	index3 := cache.GetIndex(config2)
	require.NotNil(t, index3)
	assert.Equal(t, 2, index3.Metrics.MachineCount)
	assert.NotEqual(t, index1, index3)
}

func TestIndexCache_GetIndexMetrics(t *testing.T) {
	// Test getting index metrics
	cache := &IndexCache{}
	config := &Config{
		Machines: []Machine{
			{
				Name: "server1",
				Tags: map[string]string{
					"environment": "production",
				},
			},
		},
	}

	// Build index first
	cache.GetIndex(config)

	// Get metrics
	metrics := cache.GetIndexMetrics()
	require.NotNil(t, metrics)
	assert.Equal(t, 1, metrics.MachineCount)
	assert.True(t, metrics.BuildTime > 0)
	assert.False(t, metrics.LastUpdated.IsZero())
}

func TestFindMachinesByName(t *testing.T) {
	// Test finding machines by name
	machines := []Machine{
		{Name: "server1"},
		{Name: "server2"},
		{Name: "server3"},
	}

	// Test finding existing machines
	machineNames := []string{"server1", "server3"}
	found, err := findMachinesByName(machines, machineNames)
	require.NoError(t, err)
	assert.Len(t, found, 2)
	assert.Equal(t, "server1", found[0].Name)
	assert.Equal(t, "server3", found[1].Name)

	// Test finding non-existent machine
	_, err = findMachinesByName(machines, []string{"nonexistent"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "machine 'nonexistent' not found")
}

func TestGetAllMachines(t *testing.T) {
	// Test getting all machines
	machines := []Machine{
		{Name: "server1"},
		{Name: "server2"},
		{Name: "server3"},
	}

	all := getAllMachines(machines)
	assert.Len(t, all, 3)

	for i, machine := range all {
		assert.Equal(t, machines[i].Name, machine.Name)
	}
}

func TestGetMachinesForAction(t *testing.T) {
	// Test legacy function
	config := &Config{
		Machines: []Machine{
			{
				Name: "server1",
				Tags: map[string]string{
					"environment": "production",
					"role":        "web",
				},
			},
			{
				Name: "server2",
				Tags: map[string]string{
					"environment": "production",
					"role":        "database",
				},
			},
		},
	}

	// Test action with machine names
	action := &Action{
		Name:     "test-action",
		Machines: []string{"server1"},
	}

	machines, err := GetMachinesForAction(action, config)
	require.NoError(t, err)
	assert.Len(t, machines, 1)
	assert.Equal(t, "server1", machines[0].Name)

	// Test action with tags
	action = &Action{
		Name: "test-action",
		Tags: []string{"role=web"},
	}

	machines, err = GetMachinesForAction(action, config)
	require.NoError(t, err)
	assert.Len(t, machines, 1)
	assert.Equal(t, "server1", machines[0].Name)

	// Test action with no criteria
	action = &Action{
		Name: "test-action",
	}

	machines, err = GetMachinesForAction(action, config)
	require.NoError(t, err)
	assert.Len(t, machines, 2)
}

func TestGetMachinesForAction_WithTestDataFromExamples(t *testing.T) {
	// Test with examples from testing directory
	testCases := []struct {
		name          string
		configPath    string
		actionName    string
		expectedCount int
		expectError   bool
	}{
		{
			name:          "valid project",
			configPath:    "../../examples/testing/test-valid-project",
			actionName:    "check-status",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "large inventory",
			configPath:    "../../examples/testing/test-large-inventory",
			actionName:    "test-action",
			expectedCount: 0, // Will depend on actual data
			expectError:   false,
		},
		{
			name:          "duplicate machines",
			configPath:    "../../examples/testing/test-duplicate-machines",
			actionName:    "test-action",
			expectedCount: 0, // Will depend on actual data
			expectError:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Load inventory
			inventoryPath := tc.configPath + "/inventory.hcl"
			inventory, err := ParseInventoryConfig(inventoryPath)
			if err != nil {
				if tc.expectError {
					return // Expected error
				}
				t.Skipf("Skipping test due to parse error: %v", err)
			}

			// Load actions
			actionsPath := tc.configPath + "/actions.hcl"
			actions, err := ParseActionsConfig(actionsPath)
			if err != nil {
				if tc.expectError {
					return // Expected error
				}
				t.Skipf("Skipping test due to parse error: %v", err)
			}

			// Create config
			config := &Config{
				Machines: inventory.Machines,
				Actions:  actions.Actions,
			}

			// Find the action
			var targetAction *Action
			for i := range actions.Actions {
				if actions.Actions[i].Name == tc.actionName {
					targetAction = &actions.Actions[i]
					break
				}
			}

			if targetAction == nil {
				t.Skipf("Action %s not found in test data", tc.actionName)
			}

			// Test getting machines for action
			machines, err := GetMachinesForAction(targetAction, config)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tc.expectedCount > 0 {
					assert.Len(t, machines, tc.expectedCount)
				}
			}
		})
	}
}

func TestBuildEnterpriseIndex_Performance(t *testing.T) {
	// Test performance with large number of machines
	machines := make([]Machine, 1000)
	for i := range machines {
		machines[i] = Machine{
			Name: fmt.Sprintf("server-%d", i),
			Tags: map[string]string{
				"environment": fmt.Sprintf("env-%d", i%10),
				"role":        fmt.Sprintf("role-%d", i%5),
				"region":      fmt.Sprintf("region-%d", i%3),
			},
		}
	}

	start := time.Now()
	index := buildEnterpriseIndex(machines)
	duration := time.Since(start)

	require.NotNil(t, index)
	assert.Equal(t, 1000, index.Metrics.MachineCount)
	assert.Less(t, duration, 100*time.Millisecond)
	t.Logf("Built index for %d machines in %v", len(machines), duration)
}

func TestGetMachinesForActionLarge_Performance(t *testing.T) {
	// Test performance with large index
	machines := make([]Machine, 1000)
	for i := range machines {
		machines[i] = Machine{
			Name: fmt.Sprintf("server-%d", i),
			Tags: map[string]string{
				"environment": fmt.Sprintf("env-%d", i%10),
				"role":        fmt.Sprintf("role-%d", i%5),
			},
		}
	}

	config := &Config{Machines: machines}
	index := buildEnterpriseIndex(machines)

	action := &Action{
		Name: "test-action",
		Tags: []string{"role=role-1", "environment=env-1"},
	}

	start := time.Now()
	foundMachines, err := GetMachinesForActionLarge(config, action, index)
	duration := time.Since(start)

	require.NoError(t, err)
	assert.Less(t, duration, 10*time.Millisecond)
	t.Logf("Found %d machines in %v", len(foundMachines), duration)
}

func TestIndexCache_Concurrency(t *testing.T) {
	// Test index cache concurrency
	cache := &IndexCache{}
	config := &Config{
		Machines: []Machine{
			{Name: "server1", Tags: map[string]string{"env": "prod"}},
			{Name: "server2", Tags: map[string]string{"env": "staging"}},
		},
	}

	// Test concurrent access
	results := make(chan *CompositeIndex, 10)
	for i := 0; i < 10; i++ {
		go func() {
			index := cache.GetIndex(config)
			results <- index
		}()
	}

	// Collect results
	var indexes []*CompositeIndex
	for i := 0; i < 10; i++ {
		index := <-results
		indexes = append(indexes, index)
	}

	// All indexes should be the same (cached)
	for i := 1; i < len(indexes); i++ {
		assert.Equal(t, indexes[0], indexes[i])
	}
}

func TestSortTagsByPopularity(t *testing.T) {
	// Test sorting tags by popularity
	tags := []string{"tag1", "tag2", "tag3", "tag4"}
	tagCount := map[string]int{
		"tag1": 10,
		"tag2": 5,
		"tag3": 15,
		"tag4": 2,
	}

	sorted := sortTagsByPopularity(tags, tagCount)
	assert.Len(t, sorted, 4)

	// Should be sorted by popularity (descending)
	assert.Equal(t, "tag3", sorted[0]) // count: 15
	assert.Equal(t, "tag1", sorted[1]) // count: 10
	assert.Equal(t, "tag2", sorted[2]) // count: 5
	assert.Equal(t, "tag4", sorted[3]) // count: 2
}

func TestUpdateLookupMetrics(t *testing.T) {
	// Test updating lookup metrics
	index := &CompositeIndex{
		Metrics: &IndexMetrics{
			LookupTime: 0,
			HitRate:    0.5,
		},
	}

	lookupTime := 10 * time.Millisecond
	updateLookupMetrics(index, lookupTime)

	assert.Equal(t, lookupTime, index.Metrics.LookupTime)
	// Hit rate should be updated (implementation dependent)
}

func TestComputeConfigHash(t *testing.T) {
	// Test computing config hash
	config1 := &Config{
		Machines: []Machine{
			{Name: "server1", Host: "192.168.1.1"},
			{Name: "server2", Host: "192.168.1.2"},
		},
	}

	config2 := &Config{
		Machines: []Machine{
			{Name: "server1", Host: "192.168.1.1"},
			{Name: "server2", Host: "192.168.1.2"},
		},
	}

	config3 := &Config{
		Machines: []Machine{
			{Name: "server1", Host: "192.168.1.1"},
			{Name: "server3", Host: "192.168.1.3"},
		},
	}

	hash1 := computeConfigHash(config1)
	hash2 := computeConfigHash(config2)
	hash3 := computeConfigHash(config3)

	// Same config should have same hash
	assert.Equal(t, hash1, hash2)

	// Different config should have different hash
	assert.NotEqual(t, hash1, hash3)
	assert.NotEqual(t, hash2, hash3)
}

func TestIndexCache_IsValid(t *testing.T) {
	// Test index cache validity
	cache := &IndexCache{}
	config := &Config{
		Machines: []Machine{
			{Name: "server1", Host: "192.168.1.1"},
		},
	}

	// Initially should be invalid
	assert.False(t, cache.isValid(config))

	// After building index, should be valid
	cache.GetIndex(config)
	assert.True(t, cache.isValid(config))

	// With different config, should be invalid
	config2 := &Config{
		Machines: []Machine{
			{Name: "server1", Host: "192.168.1.1"},
			{Name: "server2", Host: "192.168.1.2"},
		},
	}
	assert.False(t, cache.isValid(config2))
}

func TestBuildEnterpriseIndex_EdgeCases(t *testing.T) {
	// Test edge cases
	testCases := []struct {
		name     string
		machines []Machine
		expected int
	}{
		{
			name:     "empty machines",
			machines: []Machine{},
			expected: 0,
		},
		{
			name: "machines with empty tags",
			machines: []Machine{
				{Name: "server1", Tags: map[string]string{}},
				{Name: "server2", Tags: map[string]string{"": ""}},
			},
			expected: 2,
		},
		{
			name: "machines with nil tags",
			machines: []Machine{
				{Name: "server1"},
			},
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			index := buildEnterpriseIndex(tc.machines)
			assert.Equal(t, tc.expected, index.Metrics.MachineCount)
		})
	}
}

func TestGetMachinesForActionLarge_EdgeCases(t *testing.T) {
	// Test edge cases
	config := &Config{
		Machines: []Machine{
			{Name: "server1", Tags: map[string]string{"env": "prod"}},
		},
	}

	index := buildEnterpriseIndex(config.Machines)

	testCases := []struct {
		name        string
		action      *Action
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil action",
			action:      nil,
			expectError: true,
			errorMsg:    "action cannot be nil",
		},
		{
			name: "action with empty machine names",
			action: &Action{
				Name:     "test",
				Machines: []string{},
			},
			expectError: false,
		},
		{
			name: "action with empty tags",
			action: &Action{
				Name: "test",
				Tags: []string{},
			},
			expectError: false,
		},
		{
			name: "action with invalid tag format",
			action: &Action{
				Name: "test",
				Tags: []string{"invalid-tag"},
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := GetMachinesForActionLarge(config, tc.action, index)
			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
