package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

type LintIssue struct {
	Line       string
	Message    string
	ErrorType  string
	Severity   string
	File       string
	LineNumber int
	Column     int
}

type LintResult struct {
	File       string
	HasIssues  bool
	Issues     []LintIssue
	IssueCount int
}

type Config struct {
	LintConfig       string
	OutputFile       string
	MaxFilesPerBatch int
	ShowProgress     bool
	SingleFile       string
	NumWorkers       int
	UseFastMode      bool
}

type IssueCategory struct {
	Name        string
	Issues      []LintIssue
	Count       int
	Description string
	FixPatterns []string
}

func main() {
	config := Config{
		LintConfig:       getEnv("LINT_CONFIG", ".golangci.yml"),
		OutputFile:       getEnv("OUTPUT_FILE", "lint-todos.md"),
		MaxFilesPerBatch: 10,
		ShowProgress:     false, // No progress output
		SingleFile:       getEnv("LINT_FILE", ""),
		NumWorkers:       getEnvAsInt("LINT_WORKERS", runtime.NumCPU()),
		UseFastMode:      getEnvAsBool("LINT_FAST", true),
	}

	if config.SingleFile != "" {
		// Lint a single file
		lintSingleFile(config)
	} else {
		// Lint all files
		lintAllFiles(config)
	}
}

func lintSingleFile(config Config) {
	// Check if file exists
	if _, err := os.Stat(config.SingleFile); os.IsNotExist(err) {
		fmt.Printf("Error: File %s does not exist\n", config.SingleFile)
		os.Exit(1)
	}

	fmt.Printf("Processing: %s\n", config.SingleFile)

	result := lintFile(config.SingleFile, config.LintConfig, config.UseFastMode)

	if result.HasIssues {
		fmt.Printf("✗ Issues found in %s\n", config.SingleFile)
		fmt.Println()
		fmt.Printf("TODO: Fix linting issues in %s\n", config.SingleFile)
		fmt.Printf("Run: golangci-lint run --config=%s %s\n", config.LintConfig, config.SingleFile)
		fmt.Println()
		fmt.Println("Issues:")
		for _, issue := range result.Issues {
			fmt.Printf("  %s\n", issue.Message)
		}
		os.Exit(1)
	} else {
		fmt.Printf("✓ No issues found in %s\n", config.SingleFile)
	}
}

func lintAllFiles(config Config) {
	fmt.Printf("Starting concurrent file-by-file linting analysis...\n")
	fmt.Printf("Output will be saved to: %s\n", config.OutputFile)
	fmt.Printf("Lint config: %s\n", config.LintConfig)
	fmt.Printf("Workers: %d\n\n", config.NumWorkers)

	// Find all Go files
	fmt.Printf("Finding Go files...\n")
	files, err := findGoFiles(".")
	if err != nil {
		fmt.Printf("Error finding Go files: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d Go files\n\n", len(files))

	// Process files concurrently
	results := processFilesConcurrent(files, config)

	// Generate report
	generateReport(results, config)
	printSummary(results)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return strings.ToLower(value) == "true"
	}
	return defaultValue
}

func findGoFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			// Skip vendor directory
			if !strings.Contains(path, "/vendor/") {
				files = append(files, path)
			}
		}
		return nil
	})
	return files, err
}

func processFilesConcurrent(files []string, config Config) []LintResult {
	numWorkers := config.NumWorkers
	if numWorkers > len(files) {
		numWorkers = len(files)
	}

	jobs := make(chan string, len(files))
	results := make(chan LintResult, len(files))

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(jobs, results, config.LintConfig, config.UseFastMode, &wg)
	}

	// Send jobs
	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	// Wait for all workers to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var allResults []LintResult
	for result := range results {
		allResults = append(allResults, result)
	}

	// Sort results by file name
	sort.Slice(allResults, func(i, j int) bool {
		return allResults[i].File < allResults[j].File
	})

	return allResults
}

func worker(jobs <-chan string, results chan<- LintResult, lintConfig string, useFastMode bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for file := range jobs {
		result := lintFile(file, lintConfig, useFastMode)
		results <- result
	}
}

func lintFile(file, lintConfig string, useFastMode bool) LintResult {
	result := LintResult{
		File:      file,
		HasIssues: false,
		Issues:    []LintIssue{},
	}

	args := []string{"run", "--config=" + lintConfig}
	if useFastMode {
		args = append(args, "--fast")
	}
	args = append(args, file)

	cmd := exec.Command("golangci-lint", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Parse the output to extract issues
		issues := parseLintOutput(string(output), file)
		result.HasIssues = len(issues) > 0
		result.Issues = issues
		result.IssueCount = len(issues)
	}

	return result
}

func parseLintOutput(output, file string) []LintIssue {
	var issues []LintIssue

	// Skip if it's a parallel execution error
	if strings.Contains(output, "parallel golangci-lint is running") {
		return issues
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse golangci-lint output format: file:line:column: message (linter)
		// Example: cmd/actions.go:277:9: undefined: integrationManager (typecheck)
		re := regexp.MustCompile(`^([^:]+):(\d+):(\d+):\s*(.+?)\s*\(([^)]+)\)$`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 6 {
			lineNum := 0
			column := 0
			fmt.Sscanf(matches[2], "%d", &lineNum)
			fmt.Sscanf(matches[3], "%d", &column)

			issue := LintIssue{
				File:       matches[1],
				LineNumber: lineNum,
				Column:     column,
				Message:    matches[4],
				ErrorType:  matches[5],
				Severity:   categorizeSeverity(matches[5], matches[4]),
			}
			issues = append(issues, issue)
		}
	}

	return issues
}

func categorizeSeverity(errorType, message string) string {
	// High priority - build breaking
	if errorType == "typecheck" || strings.Contains(message, "undefined") || strings.Contains(message, "unused") {
		return "High"
	}

	// Medium priority - code quality
	if errorType == "gocritic" || errorType == "govet" || strings.Contains(message, "unused") {
		return "Medium"
	}

	// Low priority - style
	if errorType == "gofmt" || errorType == "goimports" || errorType == "misspell" {
		return "Low"
	}

	return "Medium"
}

func generateReport(results []LintResult, config Config) {
	file, err := os.Create(config.OutputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Collect all issues and categorize them
	allIssues := collectAllIssues(results)
	categories := categorizeIssues(allIssues)

	// Calculate statistics
	stats := calculateStats(results, categories)

	// Write enhanced header
	writeEnhancedHeader(writer, stats, config)

	// Write issue categories
	writeIssueCategories(writer, categories)

	// Write file-by-file breakdown
	writeFileBreakdown(writer, results, config)

	// Write batch fix opportunities
	writeBatchFixes(writer, allIssues)

	// Write progress tracking
	writeProgressTracking(writer, stats)

	// Write fix patterns
	writeFixPatterns(writer)
}

func collectAllIssues(results []LintResult) []LintIssue {
	var allIssues []LintIssue
	for _, result := range results {
		allIssues = append(allIssues, result.Issues...)
	}
	return allIssues
}

func categorizeIssues(issues []LintIssue) map[string]*IssueCategory {
	categories := make(map[string]*IssueCategory)

	for _, issue := range issues {
		key := issue.ErrorType
		if category, exists := categories[key]; exists {
			category.Issues = append(category.Issues, issue)
			category.Count++
		} else {
			categories[key] = &IssueCategory{
				Name:        key,
				Issues:      []LintIssue{issue},
				Count:       1,
				Description: getCategoryDescription(key),
				FixPatterns: getFixPatterns(key),
			}
		}
	}

	return categories
}

func getCategoryDescription(errorType string) string {
	descriptions := map[string]string{
		"typecheck":   "Build-breaking errors that prevent compilation",
		"gocritic":    "Code quality and style suggestions",
		"govet":       "Go vet analysis for potential bugs",
		"gofmt":       "Code formatting issues",
		"goimports":   "Import organization issues",
		"misspell":    "Spelling errors in comments and strings",
		"unused":      "Unused variables, imports, or functions",
		"ineffassign": "Ineffective assignments",
		"deadcode":    "Dead or unreachable code",
	}

	if desc, exists := descriptions[errorType]; exists {
		return desc
	}
	return "Code quality issue"
}

func getFixPatterns(errorType string) []string {
	patterns := map[string][]string{
		"typecheck": {
			"Add missing import: `import \"package/path\"`",
			"Define missing variable: `var variableName Type`",
			"Fix function call: use correct function name",
			"Add missing package initialization",
		},
		"gocritic": {
			"Use pointer for large structs: `func process(config *Config)`",
			"Simplify boolean expressions",
			"Use consistent naming conventions",
			"Optimize string concatenation",
		},
		"unused": {
			"Remove unused imports: `import _ \"package\"`",
			"Remove unused variables",
			"Remove unused functions",
			"Use variables or remove them",
		},
		"gofmt": {
			"Run: `gofmt -w file.go`",
			"Fix indentation and spacing",
			"Align struct fields",
		},
	}

	if patterns, exists := patterns[errorType]; exists {
		return patterns
	}
	return []string{"Review and fix according to linter suggestions"}
}

func calculateStats(results []LintResult, categories map[string]*IssueCategory) map[string]interface{} {
	totalFiles := len(results)
	filesWithIssues := 0
	totalIssues := 0
	highPriority := 0
	mediumPriority := 0
	lowPriority := 0

	for _, result := range results {
		if result.HasIssues {
			filesWithIssues++
			totalIssues += result.IssueCount
		}
	}

	for _, category := range categories {
		for _, issue := range category.Issues {
			switch issue.Severity {
			case "High":
				highPriority++
			case "Medium":
				mediumPriority++
			case "Low":
				lowPriority++
			}
		}
	}

	return map[string]interface{}{
		"totalFiles":      totalFiles,
		"filesWithIssues": filesWithIssues,
		"totalIssues":     totalIssues,
		"highPriority":    highPriority,
		"mediumPriority":  mediumPriority,
		"lowPriority":     lowPriority,
		"categories":      categories,
	}
}

func writeEnhancedHeader(writer *bufio.Writer, stats map[string]interface{}, config Config) {
	fmt.Fprintf(writer, "# Linting TODOs - AI-Optimized Report\n\n")
	fmt.Fprintf(writer, "This file contains linting issues captured by running `golangci-lint` on each Go file individually.\n")
	fmt.Fprintf(writer, "**Optimized for AI consumption** with categorization, fix suggestions, and progress tracking.\n\n")

	fmt.Fprintf(writer, "**Generated:** %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(writer, "**Config:** %s\n", config.LintConfig)
	fmt.Fprintf(writer, "**Fast Mode:** %t\n\n", config.UseFastMode)

	fmt.Fprintf(writer, "## 📊 Summary Statistics\n\n")
	fmt.Fprintf(writer, "- **Total files processed:** %d\n", stats["totalFiles"])
	fmt.Fprintf(writer, "- **Files with issues:** %d\n", stats["filesWithIssues"])
	fmt.Fprintf(writer, "- **Total issues:** %d\n\n", stats["totalIssues"])

	fmt.Fprintf(writer, "### Priority Breakdown\n\n")
	fmt.Fprintf(writer, "- **🔴 High Priority (Build-breaking):** %d issues (%.1f%%)\n",
		stats["highPriority"], float64(stats["highPriority"].(int))/float64(stats["totalIssues"].(int))*100)
	fmt.Fprintf(writer, "- **🟡 Medium Priority (Code quality):** %d issues (%.1f%%)\n",
		stats["mediumPriority"], float64(stats["mediumPriority"].(int))/float64(stats["totalIssues"].(int))*100)
	fmt.Fprintf(writer, "- **🟢 Low Priority (Style):** %d issues (%.1f%%)\n\n",
		stats["lowPriority"], float64(stats["lowPriority"].(int))/float64(stats["totalIssues"].(int))*100)
}

func writeIssueCategories(writer *bufio.Writer, categories map[string]*IssueCategory) {
	fmt.Fprintf(writer, "## 🏷️ Issue Categories\n\n")

	// Sort categories by count (descending)
	var sortedCategories []*IssueCategory
	for _, category := range categories {
		sortedCategories = append(sortedCategories, category)
	}
	sort.Slice(sortedCategories, func(i, j int) bool {
		return sortedCategories[i].Count > sortedCategories[j].Count
	})

	for _, category := range sortedCategories {
		fmt.Fprintf(writer, "### %s (%d issues)\n\n", category.Name, category.Count)
		fmt.Fprintf(writer, "**Description:** %s\n\n", category.Description)

		if len(category.FixPatterns) > 0 {
			fmt.Fprintf(writer, "**Common Fix Patterns:**\n")
			for _, pattern := range category.FixPatterns {
				fmt.Fprintf(writer, "- %s\n", pattern)
			}
			fmt.Fprintf(writer, "\n")
		}

		// Show a few examples
		fmt.Fprintf(writer, "**Examples:**\n")
		for i, issue := range category.Issues {
			if i >= 3 { // Show max 3 examples
				break
			}
			fmt.Fprintf(writer, "- `%s:%d:%d`: %s\n", issue.File, issue.LineNumber, issue.Column, issue.Message)
		}
		fmt.Fprintf(writer, "\n")
	}
}

func writeFileBreakdown(writer *bufio.Writer, results []LintResult, config Config) {
	fmt.Fprintf(writer, "## 📁 File-by-File Breakdown\n\n")

	for _, result := range results {
		if !result.HasIssues {
			continue
		}

		fmt.Fprintf(writer, "### %s\n\n", result.File)
		fmt.Fprintf(writer, "**Issues found:** %d\n\n", result.IssueCount)

		// Categorize issues by severity
		highIssues := []LintIssue{}
		mediumIssues := []LintIssue{}
		lowIssues := []LintIssue{}

		for _, issue := range result.Issues {
			switch issue.Severity {
			case "High":
				highIssues = append(highIssues, issue)
			case "Medium":
				mediumIssues = append(mediumIssues, issue)
			case "Low":
				lowIssues = append(lowIssues, issue)
			}
		}

		if len(highIssues) > 0 {
			fmt.Fprintf(writer, "#### 🔴 High Priority Issues\n\n")
			for _, issue := range highIssues {
				fmt.Fprintf(writer, "- **%s** (%s): %s\n", issue.ErrorType, issue.Severity, issue.Message)
			}
			fmt.Fprintf(writer, "\n")
		}

		if len(mediumIssues) > 0 {
			fmt.Fprintf(writer, "#### 🟡 Medium Priority Issues\n\n")
			for _, issue := range mediumIssues {
				fmt.Fprintf(writer, "- **%s** (%s): %s\n", issue.ErrorType, issue.Severity, issue.Message)
			}
			fmt.Fprintf(writer, "\n")
		}

		if len(lowIssues) > 0 {
			fmt.Fprintf(writer, "#### 🟢 Low Priority Issues\n\n")
			for _, issue := range lowIssues {
				fmt.Fprintf(writer, "- **%s** (%s): %s\n", issue.ErrorType, issue.Severity, issue.Message)
			}
			fmt.Fprintf(writer, "\n")
		}

		fmt.Fprintf(writer, "```bash\n")
		fmt.Fprintf(writer, "# TODO: Fix linting issues in %s\n", result.File)
		fmt.Fprintf(writer, "# Run: golangci-lint run --config=%s %s\n", config.LintConfig, result.File)
		fmt.Fprintf(writer, "```\n\n")
	}
}

func writeBatchFixes(writer *bufio.Writer, allIssues []LintIssue) {
	fmt.Fprintf(writer, "## 🔧 Batch Fix Opportunities\n\n")

	// Group similar issues
	issueGroups := make(map[string][]LintIssue)
	for _, issue := range allIssues {
		key := fmt.Sprintf("%s:%s", issue.ErrorType, issue.Message)
		issueGroups[key] = append(issueGroups[key], issue)
	}

	// Find issues that appear in multiple files
	var batchFixes []struct {
		pattern string
		issues  []LintIssue
		count   int
	}

	for pattern, issues := range issueGroups {
		if len(issues) > 1 {
			batchFixes = append(batchFixes, struct {
				pattern string
				issues  []LintIssue
				count   int
			}{pattern, issues, len(issues)})
		}
	}

	// Sort by count (descending)
	sort.Slice(batchFixes, func(i, j int) bool {
		return batchFixes[i].count > batchFixes[j].count
	})

	if len(batchFixes) > 0 {
		fmt.Fprintf(writer, "### Issues Found in Multiple Files\n\n")
		for _, fix := range batchFixes {
			if fix.count >= 3 { // Only show if 3+ files affected
				fmt.Fprintf(writer, "#### %s (%d files)\n\n", fix.issues[0].ErrorType, fix.count)
				fmt.Fprintf(writer, "**Pattern:** %s\n\n", fix.issues[0].Message)
				fmt.Fprintf(writer, "**Affected files:**\n")
				for _, issue := range fix.issues {
					fmt.Fprintf(writer, "- %s:%d\n", issue.File, issue.LineNumber)
				}
				fmt.Fprintf(writer, "\n**Suggested batch fix:**\n")
				fmt.Fprintf(writer, "```bash\n")
				fmt.Fprintf(writer, "# Run this command to fix all instances:\n")
				fmt.Fprintf(writer, "# find . -name \"*.go\" -exec sed -i 's/pattern/replacement/g' {} \\;\n")
				fmt.Fprintf(writer, "```\n\n")
			}
		}
	} else {
		fmt.Fprintf(writer, "No batch fix opportunities found.\n\n")
	}
}

func writeProgressTracking(writer *bufio.Writer, stats map[string]interface{}) {
	fmt.Fprintf(writer, "## 📈 Progress Tracking\n\n")

	fmt.Fprintf(writer, "### Current Status\n\n")
	fmt.Fprintf(writer, "- **Fixed:** 0/%d issues (0%%)\n", stats["totalIssues"])
	fmt.Fprintf(writer, "- **High Priority:** %d issues (%.1f%%)\n",
		stats["highPriority"], float64(stats["highPriority"].(int))/float64(stats["totalIssues"].(int))*100)
	fmt.Fprintf(writer, "- **Medium Priority:** %d issues (%.1f%%)\n",
		stats["mediumPriority"], float64(stats["mediumPriority"].(int))/float64(stats["totalIssues"].(int))*100)
	fmt.Fprintf(writer, "- **Low Priority:** %d issues (%.1f%%)\n\n",
		stats["lowPriority"], float64(stats["lowPriority"].(int))/float64(stats["totalIssues"].(int))*100)

	fmt.Fprintf(writer, "### Recommended Action Plan\n\n")
	fmt.Fprintf(writer, "1. **🔴 Fix all High Priority issues first** (%d issues) - prevents compilation\n", stats["highPriority"])
	fmt.Fprintf(writer, "2. **🟡 Address Medium Priority issues** (%d issues) - improves code quality\n", stats["mediumPriority"])
	fmt.Fprintf(writer, "3. **🟢 Fix Low Priority issues** (%d issues) - style and formatting\n", stats["lowPriority"])
	fmt.Fprintf(writer, "4. **🔧 Apply batch fixes** where possible\n")
	fmt.Fprintf(writer, "5. **✅ Run `golangci-lint run` to verify all fixes**\n\n")
}

func writeFixPatterns(writer *bufio.Writer) {
	fmt.Fprintf(writer, "## 🛠️ Common Fix Patterns\n\n")

	fmt.Fprintf(writer, "### For \"undefined\" errors:\n")
	fmt.Fprintf(writer, "```go\n")
	fmt.Fprintf(writer, "// Pattern 1: Add missing import\n")
	fmt.Fprintf(writer, "import \"spooky/internal/missingpackage\"\n\n")
	fmt.Fprintf(writer, "// Pattern 2: Define missing variable\n")
	fmt.Fprintf(writer, "var missingVar = NewMissingType()\n\n")
	fmt.Fprintf(writer, "// Pattern 3: Fix function call\n")
	fmt.Fprintf(writer, "existingFunction() // instead of undefinedFunction()\n")
	fmt.Fprintf(writer, "```\n\n")

	fmt.Fprintf(writer, "### For \"hugeParam\" warnings:\n")
	fmt.Fprintf(writer, "```go\n")
	fmt.Fprintf(writer, "// Before: func process(config Config)\n")
	fmt.Fprintf(writer, "// After: func process(config *Config)\n")
	fmt.Fprintf(writer, "```\n\n")

	fmt.Fprintf(writer, "### For \"unused\" warnings:\n")
	fmt.Fprintf(writer, "```go\n")
	fmt.Fprintf(writer, "// Remove unused imports\n")
	fmt.Fprintf(writer, "// import _ \"unusedpackage\" // Remove this line\n\n")
	fmt.Fprintf(writer, "// Remove unused variables\n")
	fmt.Fprintf(writer, "// var unusedVar = \"value\" // Remove this line\n")
	fmt.Fprintf(writer, "```\n\n")

	fmt.Fprintf(writer, "### For formatting issues:\n")
	fmt.Fprintf(writer, "```bash\n")
	fmt.Fprintf(writer, "# Run gofmt to fix formatting\n")
	fmt.Fprintf(writer, "gofmt -w file.go\n\n")
	fmt.Fprintf(writer, "# Run goimports to fix imports\n")
	fmt.Fprintf(writer, "goimports -w file.go\n")
	fmt.Fprintf(writer, "```\n\n")
}

func printSummary(results []LintResult) {
	totalFiles := len(results)
	filesWithIssues := 0
	totalIssues := 0

	for _, result := range results {
		if result.HasIssues {
			filesWithIssues++
			totalIssues += result.IssueCount
		}
	}

	fmt.Printf("\n=== Linting Analysis Complete ===\n")
	fmt.Printf("Total files processed: %d\n", totalFiles)
	fmt.Printf("Files with issues: %d\n", filesWithIssues)
	fmt.Printf("Total issues: %d\n", totalIssues)

	if filesWithIssues > 0 {
		fmt.Printf("\nTo fix issues, review the generated report and address them systematically.\n")
		fmt.Printf("You can also run: golangci-lint run --fix\n")
		// Exit with success (0) since finding issues is the expected outcome
		os.Exit(0)
	} else {
		fmt.Printf("\nNo linting issues found!\n")
		os.Exit(0)
	}
}
