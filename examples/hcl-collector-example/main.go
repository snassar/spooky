package main

import (
	spookytypes "spooky/internal/types"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	spookyfactscollectorshcl "spooky/internal/facts/collectors/hcl"
)

func main() {
	fmt.Println("=== Spooky HCL Fact Collector Example ===")

	// Create a temporary directory for our example
	tempDir, err := os.MkdirTemp("", "hcl-collector-example")
	if err != nil {
		log.Printf("Failed to create temp directory: %v", err)
		return
	}
	defer os.RemoveAll(tempDir)

	// Create sample HCL files
	if err := createSampleHCLFiles(tempDir); err != nil {
		log.Printf("Failed to create sample HCL files: %v", err)
		return
	}

	// Example 1: Single HCL file collection
	fmt.Println("1. Collecting facts from a single HCL file:")
	collector1 := spookyfactscollectorshcl.NewCollector([]string{filepath.Join(tempDir, "server-config.hcl")})

	collection1, err := collector1.Collect("web-server-01")
	if err != nil {
		log.Printf("Failed to collect facts: %v", err)
		return
	}

	fmt.Printf("   Server: %s\n", collection1.Server)
	fmt.Printf("   Facts collected: %d\n", len(collection1.Facts))

	// Display some key facts
	keyFacts := []string{"name", "version", "port", "tags"}
	for _, key := range keyFacts {
		if fact, exists := collection1.Facts[key]; exists {
			fmt.Printf("   %s: %v\n", key, fact.Value)
		}
	}
	fmt.Println()

	// Example 2: Directory collection (multiple HCL files)
	fmt.Println("2. Collecting facts from a directory (multiple HCL files):")
	collector2 := spookyfactscollectorshcl.NewCollector([]string{tempDir})

	collection2, err := collector2.Collect("all-servers")
	if err != nil {
		log.Printf("Failed to collect facts from directory: %v", err)
		return
	}

	fmt.Printf("   Server: %s\n", collection2.Server)
	fmt.Printf("   Total facts collected: %d\n", len(collection2.Facts))

	// Show facts from different files
	fmt.Println("   Facts by source:")
	for key, fact := range collection2.Facts {
		if metadata, ok := fact.Metadata["source_file"]; ok {
			fmt.Printf("     %s (from %s): %v\n", key, metadata, fact.Value)
		}
	}
	fmt.Println()

	// Example 3: Selective fact collection
	fmt.Println("3. Collecting specific facts only:")
	collector3 := spookyfactscollectorshcl.NewCollector([]string{filepath.Join(tempDir, "server-config.hcl")})

	specificKeys := []string{"name", "version"}
	collection3, err := collector3.CollectSpecific("web-server-01", specificKeys)
	if err != nil {
		log.Printf("Failed to collect specific facts: %v", err)
		return
	}

	fmt.Printf("   Requested facts: %v\n", specificKeys)
	fmt.Printf("   Facts collected: %d\n", len(collection3.Facts))
	for key, fact := range collection3.Facts {
		fmt.Printf("     %s: %v\n", key, fact.Value)
	}
	fmt.Println()

	// Example 4: Individual fact retrieval
	fmt.Println("4. Retrieving individual facts:")
	collector4 := spookyfactscollectorshcl.NewCollector([]string{filepath.Join(tempDir, "server-config.hcl")})

	fact, err := collector4.GetFact("web-server-01", "name")
	if err != nil {
		log.Printf("Failed to get fact: %v", err)
		return
	}
	fmt.Printf("   Retrieved fact 'name': %v\n", fact.Value)
	fmt.Printf("   Fact metadata: %+v\n", fact.Metadata)
	fmt.Println()

	// Example 5: HCL file validation
	fmt.Println("5. Validating HCL files:")
	collector5 := spookyfactscollectorshcl.NewCollector([]string{})

	validFile := filepath.Join(tempDir, "server-config.hcl")
	if err := collector5.ValidateHCLFile(validFile); err != nil {
		fmt.Printf("   ❌ Invalid HCL file: %v\n", err)
	} else {
		fmt.Printf("   ✅ Valid HCL file: %s\n", validFile)
	}

	// Create an invalid HCL file for testing
	invalidFile := filepath.Join(tempDir, "invalid.hcl")
	invalidHCL := `name = "test" = invalid`
	if err := os.WriteFile(invalidFile, []byte(invalidHCL), 0644); err != nil {
		log.Printf("Failed to write invalid HCL file: %v", err)
		return
	}

	if err := collector5.ValidateHCLFile(invalidFile); err != nil {
		fmt.Printf("   ❌ Invalid HCL file (expected): %v\n", err)
	} else {
		fmt.Printf("   ✅ Valid HCL file (unexpected): %s\n", invalidFile)
	}
	fmt.Println()

	// Example 6: Finding HCL files in directory
	fmt.Println("6. Finding HCL files in directory:")
	hclFiles, err := collector5.FindHCLFiles(tempDir)
	if err != nil {
		log.Printf("Failed to find HCL files: %v", err)
		return
	}

	fmt.Printf("   Found %d HCL files:\n", len(hclFiles))
	for _, file := range hclFiles {
		fmt.Printf("     - %s\n", filepath.Base(file))
	}
	fmt.Println()

	// Example 7: Exporting facts to HCL
	fmt.Println("7. Exporting facts to HCL format:")
	facts := map[string]interface{}{
		"name":    "exported-server",
		"version": "2.0.0",
		"enabled": true,
		"port":    9090,
		"tags":    []string{"exported", "example"},
		"config": map[string]interface{}{
			"host": "example.com",
			"ssl":  true,
		},
	}

	outputFile := filepath.Join(tempDir, "exported-facts.hcl")
	file, err := os.Create(outputFile)
	if err != nil {
		log.Printf("Failed to create output file: %v", err)
		return
	}
	defer file.Close()

	if err := collector5.ExportFactsToHCL(facts, file); err != nil {
		log.Printf("Failed to export facts: %v", err)
		return
	}

	fmt.Printf("   Facts exported to: %s\n", outputFile)

	// Read and display the exported content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		log.Printf("Failed to read exported file: %v", err)
		return
	}
	fmt.Printf("   Exported content:\n%s\n", string(content))

	// Example 8: JSON output for comparison
	fmt.Println("8. Converting collection to JSON for comparison:")
	jsonData, err := json.MarshalIndent(collection1, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal to JSON: %v", err)
		return
	}
	fmt.Printf("   JSON representation:\n%s\n", string(jsonData))

	fmt.Println("\n=== Example completed successfully! ===")
}

func createSampleHCLFiles(tempDir string) error {
	// Server configuration file
	serverConfig := `
# Web server configuration
name = "web-server-01"
version = "1.2.3"
enabled = true
port = 8080

# Tags
tags = ["web", "api", "production", "load-balanced"]

# Environment variables
env = {
  NODE_ENV = "production"
  LOG_LEVEL = "info"
  DEBUG = false
}
`

	// Database configuration file
	dbConfig := `
# Database configuration
name = "db-server-01"
version = "5.7.0"
enabled = true
port = 3306

# Backup configuration
backup = {
  enabled = true
  schedule = "0 2 * * *"
  retention_days = 30
  storage_path = "/backups"
}

# Monitoring
monitoring = {
  enabled = true
  metrics_port = 9104
  health_check_interval = 30
}
`

	// Network configuration file
	networkConfig := `
# Network configuration
name = "network-config"
version = "1.0.0"

# Load balancer configuration
load_balancer_algorithm = "round_robin"
health_check_path = "/health"
health_check_interval = 10
timeout = 5

# Backend configuration
backend_web1_host = "10.0.1.10"
backend_web1_port = 8080
backend_web1_weight = 1

backend_web2_host = "10.0.1.11"
backend_web2_port = 8080
backend_web2_weight = 1

# Firewall rules
firewall_web_port = 80
firewall_web_protocol = "tcp"
firewall_web_source = "0.0.0.0/0"
firewall_web_action = "allow"

firewall_api_port = 443
firewall_api_protocol = "tcp"
firewall_api_source = "10.0.0.0/8"
firewall_api_action = "allow"
`

	// Write files
	files := map[string]string{
		"server-config.hcl":  serverConfig,
		"db-config.hcl":      dbConfig,
		"network-config.hcl": networkConfig,
	}

	for filename, content := range files {
		filePath := filepath.Join(tempDir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}
	}

	return nil
}
