package config

import (
	"testing"

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

	// Test tag index
	assert.Len(t, index.TagIndex, 4) // environment=production, environment=staging, role=web, role=database
	assert.Len(t, index.TagIndex["environment=production"], 2)
	assert.Len(t, index.TagIndex["role=web"], 2)

	// Test machine tag index
	assert.Len(t, index.MachineTagIndex, 3)
	for i := range machines {
		machine := &machines[i]
		machineTags, exists := index.MachineTagIndex[machine]
		assert.True(t, exists)
		assert.NotEmpty(t, machineTags)
	}

	// Test tag count
	assert.Equal(t, 2, index.TagCount["environment=production"])
	assert.Equal(t, 2, index.TagCount["role=web"])
	assert.Equal(t, 1, index.TagCount["role=database"])

	// Test metrics
	assert.NotNil(t, index.Metrics)
	assert.Equal(t, 3, index.Metrics.MachineCount)
	assert.Equal(t, 4, index.Metrics.TagCount)
	assert.True(t, index.Metrics.BuildTime > 0)
}

func TestGetMachinesForActionLarge(t *testing.T) {
	// Test getting machines for action with large index
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

	// Test action with machine names
	action := &Action{
		Name:     "test-action",
		Machines: []string{"server1", "server2"},
	}

	resultMachines, err := GetMachinesForActionLarge(machines, action, index)
	require.NoError(t, err)
	assert.Len(t, resultMachines, 2)

	// Test action with tags
	action = &Action{
		Name: "test-action",
		Tags: []string{"role=web"},
	}

	resultMachines, err = GetMachinesForActionLarge(machines, action, index)
	require.NoError(t, err)
	assert.Len(t, resultMachines, 2) // server1 and server3 have role=web

	// Test action with no criteria (should return all machines)
	action = &Action{
		Name: "test-action",
	}

	resultMachines, err = GetMachinesForActionLarge(machines, action, index)
	require.NoError(t, err)
	assert.Len(t, resultMachines, 3)
}

func TestGetMachinesForActionLarge_InvalidMachines(t *testing.T) {
	// Test with invalid machine names
	machines := []Machine{
		{
			Name: "server1",
			Tags: map[string]string{
				"environment": "production",
			},
		},
	}

	index := buildEnterpriseIndex(machines)

	action := &Action{
		Name:     "test-action",
		Machines: []string{"nonexistent-server"},
	}

	_, err := GetMachinesForActionLarge(machines, action, index)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "machine 'nonexistent-server' not found")
}

func TestIndexCache_GetIndex(t *testing.T) {
	// Test index cache functionality
	cache := &IndexCache{}
	require.NotNil(t, cache)

	machines := []Machine{
		{
			Name: "server1",
			Tags: map[string]string{
				"environment": "production",
			},
		},
	}

	// First call should build index
	index1 := cache.GetIndex(machines)
	require.NotNil(t, index1)
	assert.Equal(t, 1, index1.Metrics.MachineCount)

	// Second call should return cached index
	index2 := cache.GetIndex(machines)
	require.NotNil(t, index2)
	assert.Equal(t, index1, index2)

	// Test with different machines
	machines2 := []Machine{
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
	}

	index3 := cache.GetIndex(machines2)
	require.NotNil(t, index3)
	assert.Equal(t, 2, index3.Metrics.MachineCount)
	assert.NotEqual(t, index1, index3)
}

func TestIndexCache_GetIndexMetrics(t *testing.T) {
	// Test getting index metrics
	cache := &IndexCache{}
	machines := []Machine{
		{
			Name: "server1",
			Tags: map[string]string{
				"environment": "production",
			},
		},
	}

	// Build index first
	cache.GetIndex(machines)

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
	// Test getting machines for action
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
	}

	// Test action with machine names
	action := &Action{
		Name:     "test-action",
		Machines: []string{"server1"},
	}

	resultMachines, err := GetMachinesForAction(action, machines)
	require.NoError(t, err)
	assert.Len(t, resultMachines, 1)
	assert.Equal(t, "server1", resultMachines[0].Name)

	// Test action with tags
	action = &Action{
		Name: "test-action",
		Tags: []string{"role=web"},
	}

	resultMachines, err = GetMachinesForAction(action, machines)
	require.NoError(t, err)
	assert.Len(t, resultMachines, 1)
	assert.Equal(t, "server1", resultMachines[0].Name)

	// Test action with no criteria (should return all machines)
	action = &Action{
		Name: "test-action",
	}

	resultMachines, err = GetMachinesForAction(action, machines)
	require.NoError(t, err)
	assert.Len(t, resultMachines, 2)
}

func TestComputeConfigHash(t *testing.T) {
	// Test config hash computation
	machines1 := []Machine{
		{Name: "server1", Host: "192.168.1.1"},
		{Name: "server2", Host: "192.168.1.2"},
	}

	machines2 := []Machine{
		{Name: "server1", Host: "192.168.1.1"},
		{Name: "server2", Host: "192.168.1.2"},
	}

	machines3 := []Machine{
		{Name: "server1", Host: "192.168.1.1"},
		{Name: "server3", Host: "192.168.1.3"},
	}

	hash1 := computeConfigHash(machines1)
	hash2 := computeConfigHash(machines2)
	hash3 := computeConfigHash(machines3)

	assert.Equal(t, hash1, hash2)    // Same machines, same hash
	assert.NotEqual(t, hash1, hash3) // Different machines, different hash
}

func TestIndexCache_IsValid(t *testing.T) {
	// Test index cache validity
	cache := &IndexCache{}
	machines := []Machine{
		{Name: "server1"},
	}

	// Initially not valid
	assert.False(t, cache.isValid(machines))

	// Build index
	cache.GetIndex(machines)

	// Should be valid now
	assert.True(t, cache.isValid(machines))

	// Test with different machines
	machines2 := []Machine{
		{Name: "server2"},
	}
	assert.False(t, cache.isValid(machines2))
}
