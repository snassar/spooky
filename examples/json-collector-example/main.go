package main

import (
	spookytypes "spooky/internal/types"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	spookyfactscollectorsjson "spooky/internal/facts/collectors/json"
	spookylogging "spooky/internal/logging"
)

func main() {
	// Create a temporary directory for the example
	tempDir, err := os.MkdirTemp("", "json-collector-example")
	if err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create sample JSON files
	if err := createSampleFiles(tempDir); err != nil {
		log.Fatalf("Failed to create sample files: %v", err)
	}

	// Initialize logger
	logger := spookylogging.GetLogger()

	// Example 1: Collect from a single JSON file
	fmt.Println("=== Example 1: Single JSON File ===")
	singleFileCollector := spookyfactscollectorsjson.NewCollector(filepath.Join(tempDir, "server1.json"), logger)
	collection, err := singleFileCollector.Collect("server1")
	if err != nil {
		log.Fatalf("Failed to collect facts: %v", err)
	}

	fmt.Printf("Collected %d facts from server1.json:\n", len(collection.Facts))
	for key, fact := range collection.Facts {
		fmt.Printf("  %s: %v (source: %s)\n", key, fact.Value, fact.Source)
	}

	// Example 2: Collect from a directory with multiple JSON files
	fmt.Println("\n=== Example 2: Directory with Multiple JSON Files ===")
	dirCollector := spookyfactscollectorsjson.NewCollector(tempDir, logger)
	collection, err = dirCollector.Collect("all-servers")
	if err != nil {
		log.Fatalf("Failed to collect facts: %v", err)
	}

	fmt.Printf("Collected %d facts from directory:\n", len(collection.Facts))
	for key, fact := range collection.Facts {
		fmt.Printf("  %s: %v (source: %s)\n", key, fact.Value, fact.Source)
	}

	// Example 3: Collect specific facts
	fmt.Println("\n=== Example 3: Collect Specific Facts ===")
	specificCollection, err := singleFileCollector.CollectSpecific("server1", []string{"name", "config.port"})
	if err != nil {
		log.Fatalf("Failed to collect specific facts: %v", err)
	}

	fmt.Printf("Collected %d specific facts:\n", len(specificCollection.Facts))
	for key, fact := range specificCollection.Facts {
		fmt.Printf("  %s: %v\n", key, fact.Value)
	}

	// Example 4: Get a single fact
	fmt.Println("\n=== Example 4: Get Single Fact ===")
	fact, err := singleFileCollector.GetFact("server1", "name")
	if err != nil {
		log.Fatalf("Failed to get fact: %v", err)
	}

	fmt.Printf("Single fact 'name': %v\n", fact.Value)
}

func createSampleFiles(tempDir string) error {
	// Server 1 configuration
	server1 := map[string]interface{}{
		"name": "web-server-1",
		"config": map[string]interface{}{
			"port":    8080,
			"enabled": true,
			"ssl": map[string]interface{}{
				"enabled": true,
				"cert":    "/etc/ssl/certs/server.crt",
			},
		},
		"tags": []string{"web", "production", "load-balancer"},
		"metadata": map[string]interface{}{
			"environment": "production",
			"region":      "us-west-2",
		},
	}

	// Server 2 configuration
	server2 := map[string]interface{}{
		"name": "db-server-1",
		"config": map[string]interface{}{
			"port":    5432,
			"enabled": true,
			"database": map[string]interface{}{
				"name":    "app_db",
				"version": "13.4",
				"backup":  true,
			},
		},
		"tags": []string{"database", "production", "primary"},
		"metadata": map[string]interface{}{
			"environment": "production",
			"region":      "us-west-2",
		},
	}

	// Write server1.json
	server1Data, err := json.MarshalIndent(server1, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal server1: %w", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "server1.json"), server1Data, 0644); err != nil {
		return fmt.Errorf("failed to write server1.json: %w", err)
	}

	// Write server2.json
	server2Data, err := json.MarshalIndent(server2, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal server2: %w", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "server2.json"), server2Data, 0644); err != nil {
		return fmt.Errorf("failed to write server2.json: %w", err)
	}

	// Write a non-JSON file (should be ignored)
	if err := os.WriteFile(filepath.Join(tempDir, "ignore.txt"), []byte("This should be ignored"), 0644); err != nil {
		return fmt.Errorf("failed to write ignore.txt: %w", err)
	}

	return nil
}
