package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type ADR struct {
	ID           string
	Title        string
	Date         time.Time
	Author       string
	Context      string
	Decision     string
	Consequences []string
	Evidence     []string
	Status       string
	Files        []string
	CommitHash   string
}

type ADRCategory struct {
	Name string
	ADRs []ADR
}

func main() {
	// Create output directory
	outputDir := "../docs/adr/data"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Generate different types of ADR reports
	categories := []ADRCategory{
		generateArchitectureADRs(),
		generateFeatureADRs(),
		generateRefactoringADRs(),
		generateInterfaceADRs(),
		generateSecurityADRs(),
	}

	// Write individual ADR files
	for _, category := range categories {
		for _, adr := range category.ADRs {
			writeADRFile(outputDir, adr)
		}
	}

	// Write summary report
	writeSummaryReport(outputDir, categories)

	fmt.Printf("Generated %d ADR reports in %s/\n", countADRs(categories), outputDir)
}

func generateArchitectureADRs() ADRCategory {
	adrs := []ADR{}

	// ADR-001: Interface-First Architecture
	adrs = append(adrs, ADR{
		ID:       "ADR-001",
		Title:    "Interface-First Architecture",
		Date:     time.Date(2025, 8, 16, 0, 0, 0, 0, time.UTC),
		Author:   "Samir Nassar",
		Context:  "The project needed a consistent way to coordinate between different system components (facts, actions, machines, templates, variables, secrets).",
		Decision: "Implemented the IntegrationManager pattern as the central coordinator, with all major components implementing specific interfaces.",
		Consequences: []string{
			"✅ Loose coupling between components",
			"✅ Testable through interface mocking",
			"✅ Clear separation of concerns",
			"❌ Additional complexity in coordination layer",
		},
		Evidence: []string{
			"IntegrationManager appears in recent commits",
			"Interface files in internal/interfaces/",
			"Manager pattern used throughout (SSHManager, FactManager, etc.)",
		},
		Status: "Accepted",
		Files:  []string{"internal/interfaces/", "internal/ssh/", "internal/facts/", "internal/actions/"},
	})

	// ADR-002: Manager Pattern
	adrs = append(adrs, ADR{
		ID:       "ADR-002",
		Title:    "Manager Pattern for System Components",
		Date:     time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
		Author:   "Samir Nassar",
		Context:  "Needed consistent patterns for managing different system domains (facts, actions, machines, templates, variables, secrets).",
		Decision: "Implemented the Manager pattern where each domain has a manager that coordinates specialized components and implements domain-specific interfaces.",
		Consequences: []string{
			"✅ Consistent architectural patterns",
			"✅ Clear domain boundaries",
			"✅ Interface-based coordination",
			"❌ Potential over-abstraction",
		},
		Evidence: []string{
			"Multiple manager implementations across packages",
			"Consistent naming pattern (*Manager)",
			"Interface-based coordination",
		},
		Status: "Accepted",
		Files:  []string{"internal/ssh/manager.go", "internal/facts/manager.go", "internal/actions/manager.go"},
	})

	return ADRCategory{Name: "Architecture", ADRs: adrs}
}

func generateFeatureADRs() ADRCategory {
	adrs := []ADR{}

	// ADR-003: Age-Based Encryption
	adrs = append(adrs, ADR{
		ID:       "ADR-003",
		Title:    "Age-Based Encryption for Secrets",
		Date:     time.Date(2025, 8, 16, 20, 10, 31, 0, time.UTC),
		Author:   "Samir Nassar",
		Context:  "Needed to replace AES-GCM encryption with a more modern, secure encryption system for secrets management.",
		Decision: "Implemented age-based encryption using filippo.io/age library with support for multiple recipients and identity-based encryption.",
		Consequences: []string{
			"✅ Modern, secure encryption",
			"✅ Multiple recipient support",
			"✅ Identity-based encryption",
			"❌ Migration complexity from existing AES-GCM",
		},
		Evidence: []string{
			"Commit 126a911: implement comprehensive age-based secrets management system",
			"Replaced AES-GCM with age encryption",
			"Added CLI commands for encryption/decryption",
		},
		Status:     "Accepted",
		CommitHash: "126a911",
		Files:      []string{"internal/secrets/", "cmd/secrets.go"},
	})

	// ADR-004: ScalVer Versioning
	adrs = append(adrs, ADR{
		ID:       "ADR-004",
		Title:    "ScalVer Versioning System",
		Date:     time.Date(2025, 8, 17, 7, 16, 21, 0, time.UTC),
		Author:   "Samir Nassar",
		Context:  "Needed a versioning system that could handle both development builds and official releases with clear semantic meaning.",
		Decision: "Implemented ScalVer (Scalable Versioning) with MAJOR.DATE.PATCH format, supporting flexible date components and git commit suffixes for development builds.",
		Consequences: []string{
			"✅ Clear version semantics",
			"✅ Development build support",
			"✅ Flexible date components",
			"❌ Non-standard versioning format",
		},
		Evidence: []string{
			"Commit 8d37df9: implement ScalVer versioning system",
			"Added internal/types/common/scalver.go",
			"642 lines of test coverage",
		},
		Status:     "Accepted",
		CommitHash: "8d37df9",
		Files:      []string{"internal/types/common/scalver.go"},
	})

	return ADRCategory{Name: "Features", ADRs: adrs}
}

func generateRefactoringADRs() ADRCategory {
	adrs := []ADR{}

	// ADR-005: HCL Schema-First Development
	adrs = append(adrs, ADR{
		ID:       "ADR-005",
		Title:    "HCL Schema-First Development",
		Date:     time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
		Author:   "Samir Nassar",
		Context:  "Needed consistent configuration validation and project structure validation across the system.",
		Decision: "Use HCL schemas for all configuration validation, with schemas defined in external .hcl files rather than embedded in code.",
		Consequences: []string{
			"✅ Consistent validation across system",
			"✅ Clear configuration contracts",
			"✅ Schema evolution support",
			"❌ Additional schema maintenance overhead",
		},
		Evidence: []string{
			"Multiple schema files in internal/schemas/schemas/",
			"Schema-driven validation mentioned in commits",
			"Project structure validation using schemas",
		},
		Status: "Accepted",
		Files:  []string{"internal/schemas/", "internal/schemas/schemas/"},
	})

	return ADRCategory{Name: "Refactoring", ADRs: adrs}
}

func generateInterfaceADRs() ADRCategory {
	adrs := []ADR{}

	// ADR-006: Interface-Based Coordination
	adrs = append(adrs, ADR{
		ID:       "ADR-006",
		Title:    "Interface-Based Coordination",
		Date:     time.Date(2025, 8, 16, 0, 0, 0, 0, time.UTC),
		Author:   "Samir Nassar",
		Context:  "Needed to ensure loose coupling between system components and enable testability through mocking.",
		Decision: "All major system components implement interfaces, with coordination happening through interface contracts rather than direct dependencies.",
		Consequences: []string{
			"✅ Loose coupling between components",
			"✅ Testable through interface mocking",
			"✅ Clear separation of concerns",
			"❌ Additional interface complexity",
		},
		Evidence: []string{
			"Interface definitions in internal/interfaces/",
			"Manager implementations satisfy interfaces",
			"IntegrationManager coordinates through interfaces",
		},
		Status: "Accepted",
		Files:  []string{"internal/interfaces/interfaces.go"},
	})

	return ADRCategory{Name: "Interfaces", ADRs: adrs}
}

func generateSecurityADRs() ADRCategory {
	adrs := []ADR{}

	// ADR-007: SSH Security Standards
	adrs = append(adrs, ADR{
		ID:       "ADR-007",
		Title:    "SSH Security Standards",
		Date:     time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
		Author:   "Samir Nassar",
		Context:  "Needed to establish security standards for SSH connections and authentication.",
		Decision: "Support only Ed25519 and RSA 4096-bit keys, validate key permissions (600), and implement proper host key validation.",
		Consequences: []string{
			"✅ Modern, secure key types",
			"✅ Proper key validation",
			"✅ Host key validation",
			"❌ Limited key type support",
		},
		Evidence: []string{
			"SSH key validation in internal/ssh/",
			"Host key validation implementation",
			"Key permission validation",
		},
		Status: "Accepted",
		Files:  []string{"internal/ssh/", "internal/types/ssh/"},
	})

	return ADRCategory{Name: "Security", ADRs: adrs}
}

func writeADRFile(outputDir string, adr ADR) {
	filename := filepath.Join(outputDir, fmt.Sprintf("%s.md", adr.ID))
	fmt.Printf("Creating file: %s\n", filename)
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Write ADR content
	fmt.Fprintf(writer, "# %s: %s\n\n", adr.ID, adr.Title)
	fmt.Fprintf(writer, "**Date:** %s\n", adr.Date.Format("2006-01-02"))
	fmt.Fprintf(writer, "**Author:** %s\n", adr.Author)
	fmt.Fprintf(writer, "**Status:** %s\n\n", adr.Status)

	if adr.CommitHash != "" {
		fmt.Fprintf(writer, "**Commit:** %s\n\n", adr.CommitHash)
	}

	fmt.Fprintf(writer, "## Context\n\n%s\n\n", adr.Context)
	fmt.Fprintf(writer, "## Decision\n\n%s\n\n", adr.Decision)

	if len(adr.Consequences) > 0 {
		fmt.Fprintf(writer, "## Consequences\n\n")
		for _, consequence := range adr.Consequences {
			fmt.Fprintf(writer, "- %s\n", consequence)
		}
		fmt.Fprintf(writer, "\n")
	}

	if len(adr.Evidence) > 0 {
		fmt.Fprintf(writer, "## Evidence\n\n")
		for _, evidence := range adr.Evidence {
			fmt.Fprintf(writer, "- %s\n", evidence)
		}
		fmt.Fprintf(writer, "\n")
	}

	if len(adr.Files) > 0 {
		fmt.Fprintf(writer, "## Related Files\n\n")
		for _, file := range adr.Files {
			fmt.Fprintf(writer, "- `%s`\n", file)
		}
		fmt.Fprintf(writer, "\n")
	}
}

func writeSummaryReport(outputDir string, categories []ADRCategory) {
	filename := filepath.Join(outputDir, "SUMMARY.md")
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating summary file: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	fmt.Fprintf(writer, "# Architecture Decision Records Summary\n\n")
	fmt.Fprintf(writer, "Generated on: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// Sort ADRs by date
	var allADRs []ADR
	for _, category := range categories {
		allADRs = append(allADRs, category.ADRs...)
	}
	sort.Slice(allADRs, func(i, j int) bool {
		return allADRs[i].Date.Before(allADRs[j].Date)
	})

	fmt.Fprintf(writer, "## Timeline\n\n")
	for _, adr := range allADRs {
		fmt.Fprintf(writer, "- **%s** (%s): [%s](%s.md) - %s\n",
			adr.Date.Format("2006-01-02"), adr.Status, adr.Title, adr.ID, adr.Author)
	}

	fmt.Fprintf(writer, "\n## By Category\n\n")
	for _, category := range categories {
		fmt.Fprintf(writer, "### %s (%d ADRs)\n\n", category.Name, len(category.ADRs))
		for _, adr := range category.ADRs {
			fmt.Fprintf(writer, "- [%s](%s.md): %s\n", adr.ID, adr.ID, adr.Title)
		}
		fmt.Fprintf(writer, "\n")
	}

	fmt.Fprintf(writer, "## Statistics\n\n")
	fmt.Fprintf(writer, "- **Total ADRs:** %d\n", countADRs(categories))
	fmt.Fprintf(writer, "- **Categories:** %d\n", len(categories))
	fmt.Fprintf(writer, "- **Accepted:** %d\n", countByStatus(allADRs, "Accepted"))
	fmt.Fprintf(writer, "- **Deprecated:** %d\n", countByStatus(allADRs, "Deprecated"))
	fmt.Fprintf(writer, "- **Superseded:** %d\n", countByStatus(allADRs, "Superseded"))
}

func countADRs(categories []ADRCategory) int {
	count := 0
	for _, category := range categories {
		count += len(category.ADRs)
	}
	return count
}

func countByStatus(adrs []ADR, status string) int {
	count := 0
	for _, adr := range adrs {
		if adr.Status == status {
			count++
		}
	}
	return count
}
