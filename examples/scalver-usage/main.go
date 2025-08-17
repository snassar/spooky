package main

import (
	"fmt"
	"log"

	spookytypescommon "spooky/internal/types/common"
)

func main() {
	fmt.Println("ScalVer (Scalable Versioning) Examples")
	fmt.Println("======================================")

	// Example 1: Parse existing ScalVer versions
	examples := []string{
		"0.2025.0",                // Yearly cadence
		"0.202508.0",              // Monthly cadence
		"0.20250812.0",            // Daily cadence
		"0.20250812.0-dev-abc123", // Development version
		"1.2025.0",                // Stable version
	}

	fmt.Println("\n1. Parsing ScalVer versions:")
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

	// Example 2: Generate new ScalVer versions
	fmt.Println("\n2. Generating ScalVer versions:")

	// Generate yearly version
	yearly, err := spookytypescommon.GenerateScalVer(0, "yearly", 0)
	if err != nil {
		log.Printf("Error generating yearly version: %v", err)
	} else {
		fmt.Printf("  Yearly: %s\n", yearly)
	}

	// Generate monthly version
	monthly, err := spookytypescommon.GenerateScalVer(0, "monthly", 0)
	if err != nil {
		log.Printf("Error generating monthly version: %v", err)
	} else {
		fmt.Printf("  Monthly: %s\n", monthly)
	}

	// Generate daily version
	daily, err := spookytypescommon.GenerateScalVer(0, "daily", 0)
	if err != nil {
		log.Printf("Error generating daily version: %v", err)
	} else {
		fmt.Printf("  Daily: %s\n", daily)
	}

	// Generate development version
	devVersion, err := spookytypescommon.GenerateDevelopmentScalVer("abc123")
	if err != nil {
		log.Printf("Error generating development version: %v", err)
	} else {
		fmt.Printf("  Development: %s\n", devVersion)
	}

	// Example 3: Version comparison
	fmt.Println("\n3. Version comparison:")
	versions := []string{"0.2025.0", "0.2025.1", "0.2026.0", "1.2025.0"}

	for i := 0; i < len(versions)-1; i++ {
		v1, _ := spookytypescommon.ParseScalVer(versions[i])
		v2, _ := spookytypescommon.ParseScalVer(versions[i+1])

		comparison := v1.Compare(v2)
		var result string
		switch comparison {
		case -1:
			result = "less than"
		case 0:
			result = "equal to"
		case 1:
			result = "greater than"
		}

		fmt.Printf("  %s is %s %s\n", versions[i], result, versions[i+1])
	}

	// Example 4: Version compatibility
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

	// Example 5: Get detailed version information
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

	// Example 6: Validation
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

	fmt.Println("\nScalVer provides:")
	fmt.Println("- Calendar-aware versioning (YYYY, YYYYMM, YYYYMMDD)")
	fmt.Println("- Semantic versioning compatibility")
	fmt.Println("- Adjustable release cadence")
	fmt.Println("- Date-Only-Grows (DOG) principle")
	fmt.Println("- Development and stable version support")
}
