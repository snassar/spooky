package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// bannedTerms defines the banned terminology patterns
var bannedTerms = []struct {
	pattern string
	message string
}{
	{
		pattern: `(?i)\bexec\b`, // exec (will be filtered by allowed patterns)
		message: "Use 'run' instead of 'exec'",
	},
	{
		pattern: `(?i)\bexecute\b`,
		message: "Use 'run' or 'orchestrate' instead of 'execute'",
	},
	{
		pattern: `(?i)\bexecution\b`,
		message: "Use 'orchestration' or 'running' instead of 'execution'",
	},
	{
		pattern: `(?i)\bexecutor\b`,
		message: "Use 'orchestrator' instead of 'executor'",
	},
	{
		pattern: `(?i)\bexecuting\b`,
		message: "Use 'running' instead of 'executing'",
	},
}

// allowedPatterns defines patterns that are allowed despite containing banned terms
var allowedPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)os/exec`),                   // Standard library
	regexp.MustCompile(`(?i)exec\.Command`),             // Standard library usage
	regexp.MustCompile(`(?i)\.Execute\(\)`),             // Framework methods (cobra, templates)
	regexp.MustCompile(`(?i)"exec"`),                    // String literals for pattern matching
	regexp.MustCompile(`(?i)exec.*system.*eval.*shell`), // Security pattern lists
	regexp.MustCompile(`(?i)"exec".*"system"`),          // Security pattern lists in arrays
	regexp.MustCompile(`(?i)forbiddenPatterns.*exec`),   // Security pattern configuration
	regexp.MustCompile(`(?i)forbiddenPatterns.*"exec"`), // Security pattern configuration with quotes
	regexp.MustCompile(`(?i)strings\.Contains.*"exec"`), // Security pattern checking
	regexp.MustCompile(`(?i)template\.Content.*"exec"`), // Template content checking
	regexp.MustCompile(`(?i)^Execute$`),                 // Cobra Execute method
	regexp.MustCompile(`(?i)cmd\.Execute`),              // Cobra command execution
	regexp.MustCompile(`(?i)tmpl\.Execute`),             // Template execution
	regexp.MustCompile(`(?i)compiledTemplate\.Execute`), // Template execution
}

func main() {
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage: terminology-checker [-v] <directory>")
		fmt.Println("Example: terminology-checker -v ./internal")
		os.Exit(1)
	}

	var violations []string
	for _, arg := range args {
		if err := filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if !strings.HasSuffix(path, ".go") {
				return nil
			}

			if verbose {
				fmt.Printf("Checking: %s\n", path)
			}

			fileViolations := checkFile(path)
			violations = append(violations, fileViolations...)

			return nil
		}); err != nil {
			fmt.Printf("Error walking directory %s: %v\n", arg, err)
			os.Exit(1)
		}
	}

	if len(violations) > 0 {
		fmt.Printf("\nFound %d terminology violations:\n\n", len(violations))
		for _, violation := range violations {
			fmt.Println(violation)
		}
		os.Exit(1)
	}

	fmt.Println("✅ No terminology violations found!")
}

func checkFile(filePath string) []string {
	var violations []string

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing %s: %v\n", filePath, err)
		return violations
	}

	// Check comments
	for _, comment := range node.Comments {
		checkText(&violations, filePath, fset.Position(comment.Pos()), comment.Text(), "comment")
	}

	// Check identifiers and string literals
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			checkText(&violations, filePath, fset.Position(x.Pos()), x.Name, "identifier")
		case *ast.BasicLit:
			if x.Kind == token.STRING {
				checkText(&violations, filePath, fset.Position(x.Pos()), x.Value, "string literal")
			}
		}
		return true
	})

	return violations
}

func checkText(violations *[]string, filePath string, pos token.Position, text, context string) {
	// Remove quotes from string literals
	text = strings.Trim(text, `"`)

	for _, banned := range bannedTerms {
		pattern := regexp.MustCompile(banned.pattern)
		if pattern.MatchString(text) {
			// Check if this is an allowed pattern
			if isAllowedPattern(text) {
				continue
			}

			violation := fmt.Sprintf("%s:%d:%d: %s in %s: %s",
				filePath, pos.Line, pos.Column, banned.message, context, text)
			*violations = append(*violations, violation)
		}
	}
}

func isAllowedPattern(text string) bool {
	// Check exact string matches first
	exactAllowed := []string{
		"exec",   // In forbiddenPatterns arrays
		"system", // In forbiddenPatterns arrays
		"eval",   // In forbiddenPatterns arrays
		"shell",  // In forbiddenPatterns arrays
	}

	for _, allowed := range exactAllowed {
		if text == allowed {
			return true
		}
	}

	// Check regex patterns
	for _, allowed := range allowedPatterns {
		if allowed.MatchString(text) {
			return true
		}
	}
	return false
}
