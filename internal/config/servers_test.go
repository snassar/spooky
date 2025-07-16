package config

import (
	"fmt"
	"testing"
	"time"
)

func TestBuildEnterpriseIndex(t *testing.T) {
	servers := []Server{
		{
			Name: "server1",
			Tags: map[string]string{"env": "prod", "region": "us-west"},
		},
		{
			Name: "server2",
			Tags: map[string]string{"env": "prod", "region": "us-east"},
		},
		{
			Name: "server3",
			Tags: map[string]string{"env": "dev", "region": "us-west"},
		},
	}

	index := buildEnterpriseIndex(servers)

	// Test tag index
	if len(index.TagIndex) != 4 { // 2 env values + 2 region values
		t.Errorf("Expected 4 tag entries, got %d", len(index.TagIndex))
	}

	// Test specific tag lookups
	if servers, exists := index.TagIndex["env=prod"]; !exists || len(servers) != 2 {
		t.Errorf("Expected 2 servers with env=prod, got %d", len(servers))
	}

	if servers, exists := index.TagIndex["region=us-west"]; !exists || len(servers) != 2 {
		t.Errorf("Expected 2 servers with region=us-west, got %d", len(servers))
	}

	// Test tag count
	if index.TagCount["env"] != 1 {
		t.Errorf("Expected env tag count 1, got %d", index.TagCount["env"])
	}

	// Test metrics
	if index.Metrics == nil {
		t.Error("Expected metrics to be set")
	}
	if index.Metrics.ServerCount != 3 {
		t.Errorf("Expected 3 servers in metrics, got %d", index.Metrics.ServerCount)
	}
}

func TestGetServersForActionLarge(t *testing.T) {
	config := &Config{
		Servers: []Server{
			{Name: "server1", Tags: map[string]string{"env": "prod", "region": "us-west"}},
			{Name: "server2", Tags: map[string]string{"env": "prod", "region": "us-east"}},
			{Name: "server3", Tags: map[string]string{"env": "dev", "region": "us-west"}},
		},
	}

	index := buildEnterpriseIndex(config.Servers)

	tests := []struct {
		name     string
		action   Action
		expected int
	}{
		{
			name:     "specific servers",
			action:   Action{Servers: []string{"server1", "server2"}},
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
			name:     "all servers",
			action:   Action{},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			servers, err := GetServersForActionLarge(config, &tt.action, index)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if len(servers) != tt.expected {
				t.Errorf("Expected %d servers, got %d", tt.expected, len(servers))
			}
		})
	}
}

func TestIndexCache(t *testing.T) {
	config := &Config{
		Servers: []Server{
			{Name: "server1", Tags: map[string]string{"env": "prod"}},
			{Name: "server2", Tags: map[string]string{"env": "dev"}},
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
	if metrics.ServerCount != 2 {
		t.Errorf("Expected 2 servers in metrics, got %d", metrics.ServerCount)
	}
}

func TestIndexPerformance(t *testing.T) {
	// Create large test dataset
	servers := make([]Server, 1000)
	for i := 0; i < 1000; i++ {
		servers[i] = Server{
			Name: fmt.Sprintf("server%d", i),
			Tags: map[string]string{
				"env":      fmt.Sprintf("env%d", i%10),
				"region":   fmt.Sprintf("region%d", i%5),
				"instance": fmt.Sprintf("instance%d", i%20),
			},
		}
	}

	config := &Config{Servers: servers}

	// Test index building performance
	start := time.Now()
	index := buildEnterpriseIndex(config.Servers)
	buildTime := time.Since(start)

	if buildTime > 100*time.Millisecond {
		t.Errorf("Index building took too long: %v", buildTime)
	}

	// Test lookup performance
	action := &Action{Tags: []string{"env=env1", "region=region1"}}

	start = time.Now()
	targetServers, err := GetServersForActionLarge(config, action, index)
	lookupTime := time.Since(start)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if lookupTime > 1*time.Millisecond {
		t.Errorf("Lookup took too long: %v", lookupTime)
	}

	// Verify results
	expectedCount := 100
	if len(targetServers) != expectedCount {
		t.Errorf("Expected %d servers, got %d", expectedCount, len(targetServers))
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
