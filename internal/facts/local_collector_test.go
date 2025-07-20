package facts

import (
	"testing"
)

func TestLocalCollector(t *testing.T) {
	collector := NewLocalCollector()

	// Test collecting all facts
	collection, err := collector.Collect("local")
	if err != nil {
		t.Fatalf("Failed to collect all facts: %v", err)
	}

	if collection == nil {
		t.Fatal("Collection is nil")
	}

	if collection.Server != "local" {
		t.Errorf("Expected server 'local', got '%s'", collection.Server)
	}

	// Check for essential facts
	essentialFacts := []string{
		FactHostname,
		FactOSName,
		FactOSVersion,
		FactOSArch,
		FactCPUCores,
		FactCPUModel,
		FactCPUArch,
		FactCPUFreq,
		FactMemoryTotal,
		FactNetworkIPs,
	}

	for _, factKey := range essentialFacts {
		if fact, exists := collection.Facts[factKey]; !exists {
			t.Errorf("Essential fact '%s' not found", factKey)
		} else if fact == nil {
			t.Errorf("Fact '%s' is nil", factKey)
		} else if fact.Source != string(SourceLocal) {
			t.Errorf("Expected fact source '%s', got '%s'", SourceLocal, fact.Source)
		}
	}
}

func TestLocalCollectorSpecificFacts(t *testing.T) {
	collector := NewLocalCollector()

	// Test collecting specific facts
	specificFacts := []string{FactHostname, FactOSName, FactCPUCores}
	collection, err := collector.CollectSpecific("local", specificFacts)
	if err != nil {
		t.Fatalf("Failed to collect specific facts: %v", err)
	}

	if collection == nil {
		t.Fatal("Collection is nil")
	}

	// Check that we have exactly the requested facts
	for _, expected := range specificFacts {
		if _, exists := collection.Facts[expected]; !exists {
			t.Errorf("Expected fact '%s' not found in collection", expected)
		}
	}

	// Note: The collector might collect additional facts beyond the requested ones
	// This is expected behavior as some facts are collected together
}

func TestLocalCollectorSingleFact(t *testing.T) {
	collector := NewLocalCollector()

	// Test collecting a single fact
	fact, err := collector.GetFact("local", FactHostname)
	if err != nil {
		t.Fatalf("Failed to collect specific fact: %v", err)
	}

	if fact == nil {
		t.Fatal("Fact is nil")
	}

	if fact.Key != FactHostname {
		t.Errorf("Expected fact key '%s', got '%s'", FactHostname, fact.Key)
	}

	if fact.Source != string(SourceLocal) {
		t.Errorf("Expected fact source '%s', got '%s'", SourceLocal, fact.Source)
	}

	// Verify the fact has a value
	if fact.Value == nil {
		t.Error("Fact value is nil")
	}

	// Verify it's a string
	if _, ok := fact.Value.(string); !ok {
		t.Errorf("Expected string value, got %T", fact.Value)
	}
}

func TestLocalCollectorCPUInfo(t *testing.T) {
	collector := NewLocalCollector()

	// Test CPU facts specifically
	cpuFacts := []string{FactCPUCores, FactCPUModel, FactCPUArch, FactCPUFreq}
	collection, err := collector.CollectSpecific("local", cpuFacts)
	if err != nil {
		t.Fatalf("Failed to collect CPU facts: %v", err)
	}

	// Check CPU cores
	if fact, exists := collection.Facts[FactCPUCores]; exists {
		if cores, ok := fact.Value.(int); ok {
			if cores <= 0 {
				t.Errorf("CPU cores should be positive, got %d", cores)
			}
		} else {
			t.Errorf("CPU cores should be int, got %T", fact.Value)
		}
	} else {
		t.Error("CPU cores fact not found")
	}

	// Check CPU model
	if fact, exists := collection.Facts[FactCPUModel]; exists {
		if model, ok := fact.Value.(string); ok {
			if model == "" {
				t.Error("CPU model should not be empty")
			}
		} else {
			t.Errorf("CPU model should be string, got %T", fact.Value)
		}
	} else {
		t.Error("CPU model fact not found")
	}

	// Check CPU frequency
	if fact, exists := collection.Facts[FactCPUFreq]; exists {
		if freq, ok := fact.Value.(string); ok {
			if freq == "" {
				t.Error("CPU frequency should not be empty")
			}
		} else {
			t.Errorf("CPU frequency should be string, got %T", fact.Value)
		}
	} else {
		t.Error("CPU frequency fact not found")
	}
}

func TestLocalCollectorMemoryInfo(t *testing.T) {
	collector := NewLocalCollector()

	// Test memory facts
	memoryFacts := []string{FactMemoryTotal, FactMemoryUsed, FactMemoryAvail}
	collection, err := collector.CollectSpecific("local", memoryFacts)
	if err != nil {
		t.Fatalf("Failed to collect memory facts: %v", err)
	}

	// Check total memory
	if fact, exists := collection.Facts[FactMemoryTotal]; exists {
		if total, ok := fact.Value.(uint64); ok {
			if total == 0 {
				t.Error("Total memory should not be zero")
			}
		} else {
			t.Errorf("Total memory should be uint64, got %T", fact.Value)
		}
	} else {
		t.Error("Total memory fact not found")
	}

	// Check used memory
	if fact, exists := collection.Facts[FactMemoryUsed]; exists {
		if used, ok := fact.Value.(uint64); ok {
			// Used memory can be 0 in some cases, but should be reasonable
			if used > 0 {
				// If we have used memory, check that it's not larger than total
				if totalFact, totalExists := collection.Facts[FactMemoryTotal]; totalExists {
					if total, totalOk := totalFact.Value.(uint64); totalOk {
						if used > total {
							t.Errorf("Used memory (%d) should not exceed total memory (%d)", used, total)
						}
					}
				}
			}
		} else {
			t.Errorf("Used memory should be uint64, got %T", fact.Value)
		}
	} else {
		t.Error("Used memory fact not found")
	}
}

func TestLocalCollectorNetworkInfo(t *testing.T) {
	collector := NewLocalCollector()

	// Test network facts
	fact, err := collector.GetFact("local", FactNetworkIPs)
	if err != nil {
		t.Fatalf("Failed to collect network facts: %v", err)
	}

	if fact == nil {
		t.Fatal("Fact is nil")
	}

	if fact.Value == nil {
		t.Error("Network IPs fact value is nil")
	}

	if ips, ok := fact.Value.([]string); ok {
		if len(ips) == 0 {
			t.Error("Should have at least one IP address")
		}

		// Check that we have at least one non-loopback IP
		hasNonLoopback := false
		for _, ip := range ips {
			if !isLoopbackIP(ip) {
				hasNonLoopback = true
				break
			}
		}

		if !hasNonLoopback {
			t.Log("Warning: Only loopback IPs found - this might be expected in some environments")
		}
	} else {
		t.Errorf("Network IPs should be []string, got %T", fact.Value)
	}
}

func TestLocalCollectorOSInfo(t *testing.T) {
	collector := NewLocalCollector()

	// Test OS facts
	osFacts := []string{FactOSName, FactOSVersion, FactOSArch}
	collection, err := collector.CollectSpecific("local", osFacts)
	if err != nil {
		t.Fatalf("Failed to collect OS facts: %v", err)
	}

	// Check OS name
	if fact, exists := collection.Facts[FactOSName]; exists {
		if osName, ok := fact.Value.(string); ok {
			if osName == "" {
				t.Error("OS name should not be empty")
			}
		} else {
			t.Errorf("OS name should be string, got %T", fact.Value)
		}
	} else {
		t.Error("OS name fact not found")
	}

	// Check OS version (might be empty on some systems)
	if fact, exists := collection.Facts[FactOSVersion]; exists {
		if osVersion, ok := fact.Value.(string); ok {
			// OS version can be empty on some systems, so we just check the type
			_ = osVersion // Use the value to avoid unused variable warning
		} else {
			t.Errorf("OS version should be string, got %T", fact.Value)
		}
	} else {
		t.Error("OS version fact not found")
	}

	// Check OS architecture
	if fact, exists := collection.Facts[FactOSArch]; exists {
		if osArch, ok := fact.Value.(string); ok {
			if osArch == "" {
				t.Error("OS architecture should not be empty")
			}
		} else {
			t.Errorf("OS architecture should be string, got %T", fact.Value)
		}
	} else {
		t.Error("OS architecture fact not found")
	}
}

func TestLocalCollectorInvalidFact(t *testing.T) {
	collector := NewLocalCollector()

	// Test collecting an invalid fact
	_, err := collector.GetFact("local", "invalid-fact")
	if err == nil {
		t.Error("Expected error when collecting invalid fact")
	}
}

func TestLocalCollectorCaching(t *testing.T) {
	collector := NewLocalCollector()

	// Collect facts twice
	collection1, err := collector.Collect("local")
	if err != nil {
		t.Fatalf("Failed to collect facts first time: %v", err)
	}

	collection2, err := collector.Collect("local")
	if err != nil {
		t.Fatalf("Failed to collect facts second time: %v", err)
	}

	// Note: The local collector doesn't implement caching, so these will be different instances
	// This is expected behavior
	if collection1 == collection2 {
		t.Log("Collections are the same instance (caching might be implemented)")
	} else {
		t.Log("Collections are different instances (no caching implemented)")
	}
}

// Helper function to check if an IP is loopback
func isLoopbackIP(ip string) bool {
	return len(ip) >= 3 && (ip[:3] == "127" || ip[:3] == "::1")
}
