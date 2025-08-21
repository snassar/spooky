package main

import (
	"fmt"
	"log"

	"spooky/internal/schemas"
)

func main() {
	fmt.Println("🔮 Spooky Schema Embedder")
	fmt.Println("==========================")

	// Initialize the schema embedder
	embedder, err := schemas.NewSchemaEmbedder()
	if err != nil {
		log.Fatalf("Failed to initialize schema embedder: %v", err)
	}

	// Print a summary of all embedded schemas
	embedder.PrintSchemaSummary()

	// Demonstrate getting a specific schema
	fmt.Println("\n📋 Example: Getting 'project' schema")
	if schema, exists := embedder.GetSchema("project"); exists {
		fmt.Printf("✅ Found project schema (%d bytes)\n", len(schema))
		fmt.Println("First 200 characters:")
		if len(schema) > 200 {
			fmt.Println(schema[:200] + "...")
		} else {
			fmt.Println(schema)
		}
	} else {
		fmt.Println("❌ Project schema not found")
	}

	// Demonstrate getting validation rules
	fmt.Println("\n✅ Example: Getting 'project' validation rules")
	if rules, exists := embedder.GetValidationRules("project"); exists {
		fmt.Printf("✅ Found project validation rules (%d bytes)\n", len(rules))
		fmt.Println("First 200 characters:")
		if len(rules) > 200 {
			fmt.Println(rules[:200] + "...")
		} else {
			fmt.Println(rules)
		}
	} else {
		fmt.Println("❌ Project validation rules not found")
	}

	// List all available schemas
	fmt.Println("\n📁 All Available Schemas:")
	schemas := embedder.ListSchemas()
	for i, name := range schemas {
		fmt.Printf("  %d. %s\n", i+1, name)
	}

	// List all available validation rules
	fmt.Println("\n✅ All Available Validation Rules:")
	rules := embedder.ListValidationRules()
	for i, name := range rules {
		fmt.Printf("  %d. %s\n", i+1, name)
	}

	// List all available metadata
	fmt.Println("\n📋 All Available Metadata:")
	metadata := embedder.ListMetadata()
	for i, name := range metadata {
		fmt.Printf("  %d. %s\n", i+1, name)
	}

	// List all available test data
	fmt.Println("\n🧪 All Available Test Data:")
	testdata := embedder.ListTestData()
	for i, name := range testdata {
		fmt.Printf("  %d. %s\n", i+1, name)
	}

	// Demonstrate getting test data
	fmt.Println("\n🧪 Example: Getting 'valid-project' test data")
	if testdata, exists := embedder.GetTestData("valid-project"); exists {
		fmt.Printf("✅ Found valid-project test data (%d bytes)\n", len(testdata))
		fmt.Println("First 200 characters:")
		if len(testdata) > 200 {
			fmt.Println(testdata[:200] + "...")
		} else {
			fmt.Println(testdata)
		}
	} else {
		fmt.Println("❌ valid-project test data not found")
	}

	fmt.Println("\n🎉 Schema embedding with test data is working!")
}
