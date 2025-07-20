package config

import (
	"fmt"
	"testing"
	"time"
)

func TestBuildEnterpriseIndex(t *testing.T) {
	machines := []Machine{
		{
			Name: "machine1",
			Tags: map[string]string{"env": "prod", "region": "us-west"},
		},
		{
			Name: "machine2",
			Tags: map[string]string{"env": "prod", "region": "us-east"},
		},
		{
			Name: "machine3",
			Tags: map[string]string{"env": "dev", "region": "us-west"},
		},
	}

	index := buildEnterpriseIndex(machines)

	// Test tag index
	if len(index.TagIndex) != 4 { // 2 env values + 2 region values
		t.Errorf("Expected 4 tag entries, got %d", len(index.TagIndex))
	}

	// Test specific tag lookups
	if machines, exists := index.TagIndex["env=prod"]; !exists || len(machines) != 2 {
		t.Errorf("Expected 2 machines with env=prod, got %d", len(machines))
	}

	if machines, exists := index.TagIndex["region=us-west"]; !exists || len(machines) != 2 {
		t.Errorf("Expected 2 machines with region=us-west, got %d", len(machines))
	}

	// Test tag count
	if index.TagCount["env"] != 1 {
		t.Errorf("Expected env tag count 1, got %d", index.TagCount["env"])
	}

	// Test metrics
	if index.Metrics == nil {
		t.Error("Expected metrics to be set")
	}
	if index.Metrics.MachineCount != 3 {
		t.Errorf("Expected 3 machines in metrics, got %d", index.Metrics.MachineCount)
	}
}

func TestGetMachinesForActionLarge(t *testing.T) {
	config := &Config{
		Machines: []Machine{
			{Name: "machine1", Tags: map[string]string{"env": "prod", "region": "us-west"}},
			{Name: "machine2", Tags: map[string]string{"env": "prod", "region": "us-east"}},
			{Name: "machine3", Tags: map[string]string{"env": "dev", "region": "us-west"}},
		},
	}

	index := buildEnterpriseIndex(config.Machines)

	tests := []struct {
		name     string
		action   Action
		expected int
	}{
		{
			name:     "specific machines",
			action:   Action{Machines: []string{"machine1", "machine2"}},
			expected: 2,
		},
		{
			name:     "tag-based selection",
			action:   Action{Tags: []string{"env=prod"}},
			expected: 2,
		},
		{
			name:     "multiple tags",
			action:   Action{Tags: []string{"env=prod", "region=us-west"}},
			expected: 1,
		},
		{
			name:     "all machines",
			action:   Action{},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			machines, err := GetMachinesForActionLarge(config, &tt.action, index)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if len(machines) != tt.expected {
				t.Errorf("Expected %d machines, got %d", tt.expected, len(machines))
			}
		})
	}
}

func TestIndexCache(t *testing.T) {
	config := &Config{
		Machines: []Machine{
			{Name: "machine1", Tags: map[string]string{"env": "prod"}},
			{Name: "machine2", Tags: map[string]string{"env": "dev"}},
		},
	}

	cache := &IndexCache{}

	// First call should build index
	index1 := cache.GetIndex(config)
	if index1 == nil {
		t.Error("Expected index to be built")
	}

	// Second call should return cached index
	index2 := cache.GetIndex(config)
	if index1 != index2 {
		t.Error("Expected cached index to be returned")
	}

	// Test metrics
	metrics := cache.GetIndexMetrics()
	if metrics.MachineCount != 2 {
		t.Errorf("Expected 2 machines in metrics, got %d", metrics.MachineCount)
	}
}

func TestIndexPerformance(t *testing.T) {
	// Create large test dataset
	machines := make([]Machine, 1000)
	for i := 0; i < 1000; i++ {
		machines[i] = Machine{
			Name: fmt.Sprintf("machine%d", i),
			Tags: map[string]string{
				"env":      fmt.Sprintf("env%d", i%10),
				"region":   fmt.Sprintf("region%d", i%5),
				"instance": fmt.Sprintf("instance%d", i%20),
			},
		}
	}

	config := &Config{Machines: machines}

	// Test index building performance
	start := time.Now()
	index := buildEnterpriseIndex(config.Machines)
	buildTime := time.Since(start)

	if buildTime > 100*time.Millisecond {
		t.Errorf("Index building took too long: %v", buildTime)
	}

	// Test lookup performance
	action := &Action{Tags: []string{"env=env1", "region=region1"}}

	start = time.Now()
	targetMachines, err := GetMachinesForActionLarge(config, action, index)
	lookupTime := time.Since(start)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if lookupTime > 1*time.Millisecond {
		t.Errorf("Lookup took too long: %v", lookupTime)
	}

	// Verify results
	expectedCount := 100
	if len(targetMachines) != expectedCount {
		t.Errorf("Expected %d machines, got %d", expectedCount, len(targetMachines))
	}
}

func TestSortTagsByPopularity(t *testing.T) {
	tagCount := map[string]int{
		"env":      100,
		"region":   50,
		"instance": 25,
	}

	tags := []string{"instance", "env", "region"}
	sorted := sortTagsByPopularity(tags, tagCount)

	expected := []string{"env", "region", "instance"}
	for i, tag := range sorted {
		if tag != expected[i] {
			t.Errorf("Expected %s at position %d, got %s", expected[i], i, tag)
		}
	}
}
