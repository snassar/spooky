package main

import (
	"fmt"
	"log"

	spookytypescommon "spooky/internal/types/common"
)

func main() {
	fmt.Println("ScalVer (Scalable Versioning) Examples")
	fmt.Println("======================================")

	runParsingExamples()
	runGenerationExamples()
	runComparisonExamples()
	runCompatibilityExamples()
	runDetailedInfoExamples()
	runValidationExamples()

	fmt.Println("\nScalVer provides:")
	fmt.Println("- Calendar-aware versioning (YYYY, YYYYMM, YYYYMMDD)")
	fmt.Println("- Semantic versioning compatibility")
	fmt.Println("- Adjustable release cadence")
	fmt.Println("- Date-Only-Grows (DOG) principle")
	fmt.Println("- Development and stable version support")
}

func runParsingExamples() {
	fmt.Println("\n1. Parsing ScalVer versions:")
	examples := []string{
		"0.2025.0",                // Yearly cadence
		"0.202508.0",              // Monthly cadence
		"0.20250812.0",            // Daily cadence
		"0.20250812.0-dev-abc123", // Development version
		"1.2025.0",                // Stable version
	}

	for _, version := range examples {
		scalver, err := spookytypescommon.ParseScalVer(version)
		if err != nil {
			log.Printf("Error parsing %s: %v", version, err)
			continue
		}

		fmt.Printf("  %s -> Major: %d, Date: %s, Patch: %d, Precision: %s, Development: %t\n",
			version, scalver.Major, scalver.Date, scalver.Patch,
			scalver.GetDatePrecision(), scalver.IsDevelopment())
	}
}

func runGenerationExamples() {
	fmt.Println("\n2. Generating ScalVer versions:")

	generateVersion("Yearly", func() (string, error) {
		return spookytypescommon.GenerateScalVer(0, "yearly", 0)
	})

	generateVersion("Monthly", func() (string, error) {
		return spookytypescommon.GenerateScalVer(0, "monthly", 0)
	})

	generateVersion("Daily", func() (string, error) {
		return spookytypescommon.GenerateScalVer(0, "daily", 0)
	})

	generateVersion("Development", func() (string, error) {
		return spookytypescommon.GenerateDevelopmentScalVer("abc123")
	})
}

func generateVersion(name string, generator func() (string, error)) {
	version, err := generator()
	if err != nil {
		log.Printf("Error generating %s version: %v", name, err)
	} else {
		fmt.Printf("  %s: %s\n", name, version)
	}
}

func runComparisonExamples() {
	fmt.Println("\n3. Version comparison:")
	versions := []string{"0.2025.0", "0.2025.1", "0.2026.0", "1.2025.0"}

	for i := 0; i < len(versions)-1; i++ {
		v1, _ := spookytypescommon.ParseScalVer(versions[i])
		v2, _ := spookytypescommon.ParseScalVer(versions[i+1])

		comparison := v1.Compare(v2)
		result := getComparisonResult(comparison)

		fmt.Printf("  %s is %s %s\n", versions[i], result, versions[i+1])
	}
}

func getComparisonResult(comparison int) string {
	switch comparison {
	case -1:
		return "less than"
	case 0:
		return "equal to"
	case 1:
		return "greater than"
	default:
		return "unknown"
	}
}

func runCompatibilityExamples() {
	fmt.Println("\n4. Version compatibility:")
	compatibilityTests := [][]string{
		{"0.2025.0", "0.2025.1"}, // Compatible development versions
		{"1.2025.0", "1.2025.1"}, // Compatible stable versions
		{"0.2025.0", "1.2025.0"}, // Incompatible major versions
	}

	for _, test := range compatibilityTests {
		compatible, err := spookytypescommon.ValidateScalVerCompatibility(test[0], test[1])
		if err != nil {
			fmt.Printf("  %s vs %s: Error - %v\n", test[0], test[1], err)
		} else {
			fmt.Printf("  %s vs %s: Compatible = %t\n", test[0], test[1], compatible)
		}
	}
}

func runDetailedInfoExamples() {
	fmt.Println("\n5. Detailed version information:")
	infoVersion := "0.20250812.0-dev-abc123"
	info, err := spookytypescommon.GetScalVerInfo(infoVersion)
	if err != nil {
		log.Printf("Error getting version info: %v", err)
	} else {
		fmt.Printf("  Version: %s\n", info["version"])
		fmt.Printf("  Major: %v\n", info["major"])
		fmt.Printf("  Date: %s\n", info["date"])
		fmt.Printf("  Patch: %v\n", info["patch"])
		fmt.Printf("  Is Development: %t\n", info["is_development"])
		fmt.Printf("  Is Stable: %t\n", info["is_stable"])
		fmt.Printf("  Date Precision: %s\n", info["date_precision"])
		fmt.Printf("  Format: %s\n", info["format"])
	}
}

func runValidationExamples() {
	fmt.Println("\n6. Version validation:")
	validationTests := []string{
		"0.2025.0",                // Valid
		"0.202508.0",              // Valid
		"0.20250812.0",            // Valid
		"0.20250812.0-dev-abc123", // Valid
		"0.2025",                  // Invalid
		"0.2025.0.1",              // Invalid
		"abc.2025.0",              // Invalid
		"",                        // Invalid
	}

	for _, test := range validationTests {
		valid := spookytypescommon.IsValidScalVerFormat(test)
		fmt.Printf("  %s: %t\n", test, valid)
	}
}
