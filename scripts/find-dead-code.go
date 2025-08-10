package main

import (
	spookytypes "spooky/internal/types"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type DeadCodeAnalyzer struct {
	usedSymbols    map[string]bool
	definedSymbols map[string][]string // symbol -> []file locations
	imports        map[string][]string // package -> []file locations
}

func NewDeadCodeAnalyzer() *DeadCodeAnalyzer {
	return &DeadCodeAnalyzer{
		usedSymbols:    make(map[string]bool),
		definedSymbols: make(map[string][]string),
		imports:        make(map[string][]string),
	}
}

func (dca *DeadCodeAnalyzer) analyzeFile(filePath string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	// Track imports
	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, "\"")
		dca.imports[importPath] = append(dca.imports[importPath], filePath)
	}

	// Track defined symbols
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Recv != nil && len(x.Recv.List) > 0 {
				// Method
				recvType := dca.getTypeName(x.Recv.List[0].Type)
				symbolName := fmt.Sprintf("%s.%s", recvType, x.Name.Name)
				dca.definedSymbols[symbolName] = append(dca.definedSymbols[symbolName], filePath)
			} else {
				// Function
				dca.definedSymbols[x.Name.Name] = append(dca.definedSymbols[x.Name.Name], filePath)
			}
		case *ast.TypeSpec:
			dca.definedSymbols[x.Name.Name] = append(dca.definedSymbols[x.Name.Name], filePath)
		case *ast.ValueSpec:
			for _, name := range x.Names {
				dca.definedSymbols[name.Name] = append(dca.definedSymbols[name.Name], filePath)
			}
		}
		return true
	})

	// Track used symbols
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			if x.Obj == nil {
				// This might be a reference to an imported symbol
				dca.usedSymbols[x.Name] = true
			}
		case *ast.SelectorExpr:
			if ident, ok := x.X.(*ast.Ident); ok {
				symbolName := fmt.Sprintf("%s.%s", ident.Name, x.Sel.Name)
				dca.usedSymbols[symbolName] = true
			}
		}
		return true
	})

	return nil
}

func (dca *DeadCodeAnalyzer) getTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + dca.getTypeName(t.X)
	case *ast.SelectorExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return fmt.Sprintf("%s.%s", ident.Name, t.Sel.Name)
		}
	}
	return "unknown"
}

func (dca *DeadCodeAnalyzer) findDeadCode() []string {
	var deadCode []string

	for symbol, locations := range dca.definedSymbols {
		if !dca.usedSymbols[symbol] {
			for _, location := range locations {
				deadCode = append(deadCode, fmt.Sprintf("UNUSED: %s in %s", symbol, location))
			}
		}
	}

	return deadCode
}

func (dca *DeadCodeAnalyzer) findUnusedImports() []string {
	var unusedImports []string

	for importPath, locations := range dca.imports {
		// Skip standard library imports for now
		if !strings.Contains(importPath, ".") {
			continue
		}

		// Check if any symbols from this import are used
		importUsed := false
		for symbol := range dca.usedSymbols {
			if strings.HasPrefix(symbol, filepath.Base(importPath)+".") {
				importUsed = true
				break
			}
		}

		if !importUsed {
			for _, location := range locations {
				unusedImports = append(unusedImports, fmt.Sprintf("UNUSED IMPORT: %s in %s", importPath, location))
			}
		}
	}

	return unusedImports
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run find-dead-code.go <directory>")
		os.Exit(1)
	}

	rootDir := os.Args[1]
	analyzer := NewDeadCodeAnalyzer()

	// Walk through all Go files
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			// Skip test files for now
			if !strings.HasSuffix(path, "_test.go") {
				if err := analyzer.analyzeFile(path); err != nil {
					fmt.Printf("Warning: %v\n", err)
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}

	// Find dead code
	deadCode := analyzer.findDeadCode()
	unusedImports := analyzer.findUnusedImports()

	// Report results
	fmt.Println("=== DEAD CODE ANALYSIS ===")
	fmt.Printf("Analyzed directory: %s\n\n", rootDir)

	if len(deadCode) == 0 && len(unusedImports) == 0 {
		fmt.Println("✅ No dead code found!")
	} else {
		if len(deadCode) > 0 {
			fmt.Println("🔍 POTENTIAL DEAD CODE:")
			for _, item := range deadCode {
				fmt.Printf("  %s\n", item)
			}
			fmt.Println()
		}

		if len(unusedImports) > 0 {
			fmt.Println("📦 UNUSED IMPORTS:")
			for _, item := range unusedImports {
				fmt.Printf("  %s\n", item)
			}
			fmt.Println()
		}

		fmt.Printf("📊 SUMMARY:\n")
		fmt.Printf("  - Potential dead code items: %d\n", len(deadCode))
		fmt.Printf("  - Unused imports: %d\n", len(unusedImports))
	}
}
