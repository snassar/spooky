package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type Dependency struct {
	From string
	To   string
	Type string // "import", "interface", "type", "function"
}

type Package struct {
	Name         string
	Path         string
	Dependencies []string
	Interfaces   []string
	Types        []string
	Functions    []string
}

type GraphData struct {
	Packages     map[string]*Package
	Dependencies []Dependency
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: dependency-graph <command> [options]")
		fmt.Println("Commands:")
		fmt.Println("  analyze    - Analyze all dependencies")
		fmt.Println("  package    - Analyze specific package")
		fmt.Println("  interface  - Show interface dependencies")
		fmt.Println("  type       - Show type dependencies")
		fmt.Println("  function   - Show function dependencies")
		os.Exit(1)
	}

	command := os.Args[1]
	outputDir := "../../docs/dependencies"

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	switch command {
	case "analyze":
		analyzeAllDependencies(outputDir)
	case "package":
		if len(os.Args) < 3 {
			fmt.Println("Usage: dependency-graph package <package-path>")
			os.Exit(1)
		}
		analyzePackage(os.Args[2], outputDir)
	case "interface":
		analyzeInterfaceDependencies(outputDir)
	case "type":
		analyzeTypeDependencies(outputDir)
	case "function":
		analyzeFunctionDependencies(outputDir)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func analyzeAllDependencies(outputDir string) {
	fmt.Println("Analyzing all dependencies...")

	graphData := &GraphData{
		Packages: make(map[string]*Package),
	}

	// Analyze internal packages
	internalPackages := []string{
		"internal/actions", "internal/facts", "internal/machines",
		"internal/ssh", "internal/templates", "internal/variables",
		"internal/interfaces", "internal/types", "internal/schemas",
		"internal/config", "internal/logging", "internal/secrets",
		"internal/project", "internal/cli",
	}

	for _, pkg := range internalPackages {
		analyzePackageForGraph(pkg, graphData)
	}

	// Generate different types of graphs
	generatePackageGraph(graphData, outputDir)
	if err := generateInterfaceGraph(graphData, outputDir); err != nil {
		fmt.Printf("Error generating interface graph: %v\n", err)
	}
	if err := generateTypeGraph(graphData, outputDir); err != nil {
		fmt.Printf("Error generating type graph: %v\n", err)
	}
	if err := generateFunctionGraph(graphData, outputDir); err != nil {
		fmt.Printf("Error generating function graph: %v\n", err)
	}
	generateSummary(graphData, outputDir)

	fmt.Printf("Dependency analysis complete. Check %s/\n", outputDir)
}

func analyzePackage(pkgPath string, outputDir string) {
	fmt.Printf("Analyzing package: %s\n", pkgPath)

	graphData := &GraphData{
		Packages: make(map[string]*Package),
	}

	analyzePackageForGraph(pkgPath, graphData)
	generatePackageGraph(graphData, outputDir)

	fmt.Printf("Package analysis complete. Check %s/\n", outputDir)
}

func analyzePackageForGraph(pkgPath string, graphData *GraphData) {
	pkg := &Package{
		Name:         filepath.Base(pkgPath),
		Path:         pkgPath,
		Dependencies: []string{},
		Interfaces:   []string{},
		Types:        []string{},
		Functions:    []string{},
	}

	// Find all Go files in the package
	files, err := filepath.Glob(filepath.Join(pkgPath, "*.go"))
	if err != nil {
		fmt.Printf("Error finding files in %s: %v\n", pkgPath, err)
		return
	}

	for _, file := range files {
		analyzeGoFile(file, pkg, graphData)
	}

	graphData.Packages[pkgPath] = pkg
}

func analyzeGoFile(filePath string, pkg *Package, graphData *GraphData) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return
	}

	contentStr := string(content)

	// Extract imports
	importRegex := regexp.MustCompile(`import\s+\(([\s\S]*?)\)`)
	importMatches := importRegex.FindAllStringSubmatch(contentStr, -1)
	for _, match := range importMatches {
		imports := extractImports(match[1])
		for _, imp := range imports {
			if strings.HasPrefix(imp, "spooky/") {
				pkg.Dependencies = append(pkg.Dependencies, imp)
			}
		}
	}

	// Extract interfaces
	interfaceRegex := regexp.MustCompile(`type\s+(\w+)\s+interface\s*\{`)
	interfaceMatches := interfaceRegex.FindAllStringSubmatch(contentStr, -1)
	for _, match := range interfaceMatches {
		pkg.Interfaces = append(pkg.Interfaces, match[1])
	}

	// Extract types
	typeRegex := regexp.MustCompile(`type\s+(\w+)\s+(?:struct|interface)`)
	typeMatches := typeRegex.FindAllStringSubmatch(contentStr, -1)
	for _, match := range typeMatches {
		pkg.Types = append(pkg.Types, match[1])
	}

	// Extract functions
	funcRegex := regexp.MustCompile(`func\s+(?:\(\w+\s+\*?\w+\)\s+)?(\w+)\s*\(`)
	funcMatches := funcRegex.FindAllStringSubmatch(contentStr, -1)
	for _, match := range funcMatches {
		pkg.Functions = append(pkg.Functions, match[1])
	}
}

func extractImports(importBlock string) []string {
	var imports []string
	lines := strings.Split(importBlock, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		// Remove quotes and extract package path
		line = strings.Trim(line, `"`)
		if line != "" {
			imports = append(imports, line)
		}
	}
	return imports
}

func generatePackageGraph(graphData *GraphData, outputDir string) {
	filename := filepath.Join(outputDir, "package-dependencies.md")
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	fmt.Fprintf(writer, "# Package Dependencies\n\n")
	fmt.Fprintf(writer, "Generated on: %s\n\n", getCurrentTime())

	// Generate Mermaid diagram
	fmt.Fprintf(writer, "```mermaid\ngraph TD\n")

	// Add nodes
	for pkgPath, pkg := range graphData.Packages {
		nodeName := getNodeName(pkgPath)
		fmt.Fprintf(writer, "    %s[\"%s\"]\n", nodeName, pkg.Name)
	}

	// Add edges
	for pkgPath, pkg := range graphData.Packages {
		fromNode := getNodeName(pkgPath)
		for _, dep := range pkg.Dependencies {
			if _, exists := graphData.Packages[dep]; exists {
				toNode := getNodeName(dep)
				fmt.Fprintf(writer, "    %s --> %s\n", fromNode, toNode)
			}
		}
	}

	fmt.Fprintf(writer, "```\n\n")

	// Add package details
	fmt.Fprintf(writer, "## Package Details\n\n")
	for _, pkg := range graphData.Packages {
		fmt.Fprintf(writer, "### %s\n\n", pkg.Name)
		fmt.Fprintf(writer, "- **Path:** `%s`\n", pkg.Path)
		fmt.Fprintf(writer, "- **Dependencies:** %s\n", strings.Join(pkg.Dependencies, ", "))
		fmt.Fprintf(writer, "- **Interfaces:** %s\n", strings.Join(pkg.Interfaces, ", "))
		fmt.Fprintf(writer, "- **Types:** %s\n", strings.Join(pkg.Types, ", "))
		fmt.Fprintf(writer, "- **Functions:** %s\n\n", strings.Join(pkg.Functions, ", "))
	}
}

// Generic graph generator function to eliminate code duplication
func generateGraph(
	outputFile string,
	sectionTitle string,
	data map[string]*Package,
	extractItems func(*Package) []string,
) error {
	// Create output file
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", outputFile, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Write header
	fmt.Fprintf(writer, "# %s\n\n", sectionTitle)
	fmt.Fprintf(writer, "Generated on: %s\n\n", getCurrentTime())

	// Write graph description
	fmt.Fprintf(writer, "```mermaid\ngraph TD\n")

	// Process each package
	for pkgPath, pkg := range data {
		items := extractItems(pkg)
		for _, item := range items {
			nodeName := fmt.Sprintf("%s_%s", getNodeName(pkgPath), item)
			fmt.Fprintf(writer, "    %s[\"%s.%s\"]\n", nodeName, pkg.Name, item)
		}
	}

	// Write footer
	fmt.Fprintf(writer, "```\n\n")

	// Write details section
	fmt.Fprintf(writer, "## %s Details\n\n", sectionTitle)
	for _, pkg := range data {
		items := extractItems(pkg)
		if len(items) > 0 {
			fmt.Fprintf(writer, "### %s\n\n", pkg.Name)
			for _, item := range items {
				fmt.Fprintf(writer, "- **%s**\n", item)
			}
			fmt.Fprintf(writer, "\n")
		}
	}

	return nil
}

// Specific implementations using the generic function
func generateInterfaceGraph(graphData *GraphData, outputDir string) error {
	return generateGraph(
		filepath.Join(outputDir, "interface-dependencies.md"),
		"Interface Dependencies",
		graphData.Packages,
		func(pkg *Package) []string { return pkg.Interfaces },
	)
}

func generateTypeGraph(graphData *GraphData, outputDir string) error {
	return generateGraph(
		filepath.Join(outputDir, "type-dependencies.md"),
		"Type Dependencies",
		graphData.Packages,
		func(pkg *Package) []string { return pkg.Types },
	)
}

func generateFunctionGraph(graphData *GraphData, outputDir string) error {
	return generateGraph(
		filepath.Join(outputDir, "function-dependencies.md"),
		"Function Dependencies",
		graphData.Packages,
		func(pkg *Package) []string { return pkg.Functions },
	)
}

func generateSummary(graphData *GraphData, outputDir string) {
	filename := filepath.Join(outputDir, "README.md")
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	fmt.Fprintf(writer, "# Dependency Analysis\n\n")
	fmt.Fprintf(writer, "Generated on: %s\n\n", getCurrentTime())

	fmt.Fprintf(writer, "## Summary\n\n")
	fmt.Fprintf(writer, "- **Total packages analyzed:** %d\n", len(graphData.Packages))

	totalInterfaces := 0
	totalTypes := 0
	totalFunctions := 0
	for _, pkg := range graphData.Packages {
		totalInterfaces += len(pkg.Interfaces)
		totalTypes += len(pkg.Types)
		totalFunctions += len(pkg.Functions)
	}

	fmt.Fprintf(writer, "- **Total interfaces:** %d\n", totalInterfaces)
	fmt.Fprintf(writer, "- **Total types:** %d\n", totalTypes)
	fmt.Fprintf(writer, "- **Total functions:** %d\n\n", totalFunctions)

	fmt.Fprintf(writer, "## Generated Reports\n\n")
	fmt.Fprintf(writer, "- [Package Dependencies](package-dependencies.md) - Overall package dependency graph\n")
	fmt.Fprintf(writer, "- [Interface Dependencies](interface-dependencies.md) - Interface relationships\n")
	fmt.Fprintf(writer, "- [Type Dependencies](type-dependencies.md) - Type relationships\n")
	fmt.Fprintf(writer, "- [Function Dependencies](function-dependencies.md) - Function relationships\n\n")

	fmt.Fprintf(writer, "## Usage\n\n")
	fmt.Fprintf(writer, "```bash\n")
	fmt.Fprintf(writer, "# Generate all dependency graphs\n")
	fmt.Fprintf(writer, "just dependency-graph\n\n")
	fmt.Fprintf(writer, "# Analyze specific package\n")
	fmt.Fprintf(writer, "just dependency-graph-package internal/ssh\n\n")
	fmt.Fprintf(writer, "# Generate specific graph type\n")
	fmt.Fprintf(writer, "just dependency-graph-interface\n")
	fmt.Fprintf(writer, "just dependency-graph-type\n")
	fmt.Fprintf(writer, "just dependency-graph-function\n")
	fmt.Fprintf(writer, "```\n")
}

func analyzeInterfaceDependencies(outputDir string) {
	fmt.Println("Analyzing interface dependencies...")
	graphData := &GraphData{Packages: make(map[string]*Package)}

	internalPackages := []string{
		"internal/actions", "internal/facts", "internal/machines",
		"internal/ssh", "internal/templates", "internal/variables",
		"internal/interfaces", "internal/types", "internal/schemas",
		"internal/config", "internal/logging", "internal/secrets",
		"internal/project", "internal/cli",
	}

	for _, pkg := range internalPackages {
		analyzePackageForGraph(pkg, graphData)
	}

	if err := generateInterfaceGraph(graphData, outputDir); err != nil {
		fmt.Printf("Error generating interface graph: %v\n", err)
	}
	fmt.Printf("Interface analysis complete. Check %s/\n", outputDir)
}

func analyzeTypeDependencies(outputDir string) {
	fmt.Println("Analyzing type dependencies...")
	graphData := &GraphData{Packages: make(map[string]*Package)}

	internalPackages := []string{
		"internal/actions", "internal/facts", "internal/machines",
		"internal/ssh", "internal/templates", "internal/variables",
		"internal/interfaces", "internal/types", "internal/schemas",
		"internal/config", "internal/logging", "internal/secrets",
		"internal/project", "internal/cli",
	}

	for _, pkg := range internalPackages {
		analyzePackageForGraph(pkg, graphData)
	}

	if err := generateTypeGraph(graphData, outputDir); err != nil {
		fmt.Printf("Error generating type graph: %v\n", err)
	}
	fmt.Printf("Type analysis complete. Check %s/\n", outputDir)
}

func analyzeFunctionDependencies(outputDir string) {
	fmt.Println("Analyzing function dependencies...")
	graphData := &GraphData{Packages: make(map[string]*Package)}

	internalPackages := []string{
		"internal/actions", "internal/facts", "internal/machines",
		"internal/ssh", "internal/templates", "internal/variables",
		"internal/interfaces", "internal/types", "internal/schemas",
		"internal/config", "internal/logging", "internal/secrets",
		"internal/project", "internal/cli",
	}

	for _, pkg := range internalPackages {
		analyzePackageForGraph(pkg, graphData)
	}

	if err := generateFunctionGraph(graphData, outputDir); err != nil {
		fmt.Printf("Error generating function graph: %v\n", err)
	}
	fmt.Printf("Function analysis complete. Check %s/\n", outputDir)
}

func getNodeName(pkgPath string) string {
	return strings.ReplaceAll(pkgPath, "/", "_")
}

func getCurrentTime() string {
	output, err := exec.Command("date").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}
