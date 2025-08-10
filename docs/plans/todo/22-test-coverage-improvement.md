# Implementation Plan: Test Coverage Improvement

## Overview
Improve test coverage across the spooky codebase to ensure code quality, reliability, and maintainability.

## Task Details
- **Task ID**: 8.3
- **Priority**: Medium
- **Files**: All `internal/` packages
- **Functions**: Test coverage improvement, test quality enhancement, test infrastructure

## Current State Analysis

### Existing Test Patterns
1. **Inconsistent Coverage**: Test coverage varies significantly across packages (~60% average)
2. **Missing Tests**: Some packages lack comprehensive tests
3. **Test Quality**: Some tests are basic and don't cover edge cases
4. **Integration Tests**: Limited integration test coverage

### Test Coverage Issues Found
1. **Low Coverage**: Many packages have coverage below 70%
2. **Missing Edge Cases**: Tests don't cover error conditions and edge cases
3. **Integration Gaps**: Limited integration testing between packages
4. **Performance Tests**: No performance testing infrastructure

## Implementation Requirements

### Test Coverage Compliance
The test coverage improvement must:
1. **Achieve minimum coverage** of 80% across all packages
2. **Cover edge cases** and error conditions
3. **Add integration tests** for package interactions
4. **Implement performance tests** for critical paths
5. **Improve test quality** and maintainability
6. **Create test infrastructure** for future development

### Required Dependencies
- All existing packages
- Testing framework (Go testing)
- Coverage tools
- Mocking framework

## Detailed Implementation Plan

### Step 1: Coverage Analysis and Planning

#### 1.1 Coverage Analysis Tool
```go
// internal/testing/coverage/analyzer.go
package coverage

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

// CoverageAnalyzer analyzes test coverage across packages
type CoverageAnalyzer struct {
    packages map[string]*PackageCoverage
    overall  *OverallCoverage
}

// PackageCoverage represents coverage for a package
type PackageCoverage struct {
    Name           string
    Coverage       float64
    LinesCovered   int
    LinesTotal     int
    FunctionsCovered int
    FunctionsTotal int
    BranchesCovered int
    BranchesTotal  int
    Issues         []CoverageIssue
}

// OverallCoverage represents overall coverage
type OverallCoverage struct {
    Coverage       float64
    Packages       int
    PackagesBelowThreshold int
    Issues         []CoverageIssue
}

// CoverageIssue represents a coverage issue
type CoverageIssue struct {
    Type        IssueType
    Package     string
    File        string
    Line        int
    Description string
    Severity    Severity
}

// IssueType represents the type of coverage issue
type IssueType string

const (
    IssueTypeLowCoverage    IssueType = "low_coverage"
    IssueTypeMissingTests   IssueType = "missing_tests"
    IssueTypeNoTests        IssueType = "no_tests"
    IssueTypeEdgeCase       IssueType = "edge_case"
)

// Severity represents issue severity
type Severity string

const (
    SeverityCritical Severity = "critical"
    SeverityHigh     Severity = "high"
    SeverityMedium   Severity = "medium"
    SeverityLow      Severity = "low"
)

// NewCoverageAnalyzer creates a new coverage analyzer
func NewCoverageAnalyzer() *CoverageAnalyzer {
    return &CoverageAnalyzer{
        packages: make(map[string]*PackageCoverage),
        overall:  &OverallCoverage{},
    }
}

// AnalyzeCoverage analyzes coverage across all packages
func (ca *CoverageAnalyzer) AnalyzeCoverage(rootPath string) error {
    // Walk through all packages
    err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        // Skip vendor and test directories
        if info.IsDir() && (strings.Contains(path, "vendor") || strings.Contains(path, "test")) {
            return filepath.SkipDir
        }
        
        // Check if this is a Go package
        if info.IsDir() && ca.isGoPackage(path) {
            if err := ca.analyzePackage(path); err != nil {
                return fmt.Errorf("failed to analyze package %s: %w", path, err)
            }
        }
        
        return nil
    })
    
    if err != nil {
        return fmt.Errorf("failed to walk directory: %w", err)
    }
    
    // Calculate overall coverage
    ca.calculateOverallCoverage()
    
    return nil
}

// analyzePackage analyzes coverage for a specific package
func (ca *CoverageAnalyzer) analyzePackage(packagePath string) error {
    packageName := filepath.Base(packagePath)
    
    // Run coverage analysis
    coverage, err := ca.runCoverageAnalysis(packagePath)
    if err != nil {
        return err
    }
    
    // Analyze coverage issues
    issues := ca.analyzeCoverageIssues(packagePath, coverage)
    
    ca.packages[packageName] = &PackageCoverage{
        Name:           packageName,
        Coverage:       coverage.Percentage,
        LinesCovered:   coverage.LinesCovered,
        LinesTotal:     coverage.LinesTotal,
        FunctionsCovered: coverage.FunctionsCovered,
        FunctionsTotal: coverage.FunctionsTotal,
        BranchesCovered: coverage.BranchesCovered,
        BranchesTotal:  coverage.BranchesTotal,
        Issues:         issues,
    }
    
    return nil
}

// isGoPackage checks if a directory is a Go package
func (ca *CoverageAnalyzer) isGoPackage(path string) bool {
    // Check for .go files
    entries, err := os.ReadDir(path)
    if err != nil {
        return false
    }
    
    for _, entry := range entries {
        if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
            return true
        }
    }
    
    return false
}
```

#### 1.2 Coverage Report Generator
```go
// internal/testing/coverage/report.go
package coverage

import (
    "fmt"
    "os"
    "path/filepath"
    "text/template"
)

// ReportGenerator generates coverage reports
type ReportGenerator struct {
    analyzer *CoverageAnalyzer
}

// NewReportGenerator creates a new report generator
func NewReportGenerator(analyzer *CoverageAnalyzer) *ReportGenerator {
    return &ReportGenerator{
        analyzer: analyzer,
    }
}

// GenerateReport generates a coverage report
func (rg *ReportGenerator) GenerateReport(outputPath string) error {
    // Create output directory
    if err := os.MkdirAll(outputPath, 0755); err != nil {
        return fmt.Errorf("failed to create output directory: %w", err)
    }
    
    // Generate HTML report
    if err := rg.generateHTMLReport(outputPath); err != nil {
        return err
    }
    
    // Generate markdown report
    if err := rg.generateMarkdownReport(outputPath); err != nil {
        return err
    }
    
    return nil
}

// generateHTMLReport generates an HTML coverage report
func (rg *ReportGenerator) generateHTMLReport(outputPath string) error {
    templateContent := `<!DOCTYPE html>
<html>
<head>
    <title>Test Coverage Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .package { margin: 10px 0; padding: 10px; border: 1px solid #ccc; }
        .coverage { font-weight: bold; }
        .low { color: red; }
        .medium { color: orange; }
        .high { color: green; }
        .issue { margin: 5px 0; padding: 5px; background: #f0f0f0; }
    </style>
</head>
<body>
    <h1>Test Coverage Report</h1>
    <div class="overall">
        <h2>Overall Coverage: {{printf "%.1f" .Overall.Coverage}}%</h2>
        <p>Packages: {{.Overall.Packages}}</p>
        <p>Packages below threshold: {{.Overall.PackagesBelowThreshold}}</p>
    </div>
    
    {{range $name, $pkg := .Packages}}
    <div class="package">
        <h3>{{$pkg.Name}} - {{printf "%.1f" $pkg.Coverage}}%</h3>
        <p>Lines: {{$pkg.LinesCovered}}/{{$pkg.LinesTotal}}</p>
        <p>Functions: {{$pkg.FunctionsCovered}}/{{$pkg.FunctionsTotal}}</p>
        <p>Branches: {{$pkg.BranchesCovered}}/{{$pkg.BranchesTotal}}</p>
        
        {{if $pkg.Issues}}
        <h4>Issues:</h4>
        {{range $pkg.Issues}}
        <div class="issue">
            <strong>{{.Type}}</strong>: {{.Description}}
        </div>
        {{end}}
        {{end}}
    </div>
    {{end}}
</body>
</html>`
    
    tmpl, err := template.New("coverage").Parse(templateContent)
    if err != nil {
        return fmt.Errorf("failed to parse template: %w", err)
    }
    
    file, err := os.Create(filepath.Join(outputPath, "coverage.html"))
    if err != nil {
        return fmt.Errorf("failed to create HTML report: %w", err)
    }
    defer file.Close()
    
    data := struct {
        Overall  *OverallCoverage
        Packages map[string]*PackageCoverage
    }{
        Overall:  rg.analyzer.overall,
        Packages: rg.analyzer.packages,
    }
    
    return tmpl.Execute(file, data)
}
```

### Step 2: Test Infrastructure

#### 2.1 Test Utilities
```go
// internal/testing/utils/utils.go
package utils

import (
    "testing"
    "time"
)

// TestUtils provides test utility functions
type TestUtils struct{}

// NewTestUtils creates new test utilities
func NewTestUtils() *TestUtils {
    return &TestUtils{}
}

// AssertError asserts that an error occurred
func (tu *TestUtils) AssertError(t *testing.T, err error, message string) {
    t.Helper()
    if err == nil {
        t.Errorf("%s: expected error but got nil", message)
    }
}

// AssertNoError asserts that no error occurred
func (tu *TestUtils) AssertNoError(t *testing.T, err error, message string) {
    t.Helper()
    if err != nil {
        t.Errorf("%s: unexpected error: %v", message, err)
    }
}

// AssertEqual asserts that two values are equal
func (tu *TestUtils) AssertEqual(t *testing.T, expected, actual interface{}, message string) {
    t.Helper()
    if expected != actual {
        t.Errorf("%s: expected %v but got %v", message, expected, actual)
    }
}

// AssertNotNil asserts that a value is not nil
func (tu *TestUtils) AssertNotNil(t *testing.T, value interface{}, message string) {
    t.Helper()
    if value == nil {
        t.Errorf("%s: expected non-nil value", message)
    }
}

// WaitForCondition waits for a condition to be true
func (tu *TestUtils) WaitForCondition(t *testing.T, condition func() bool, timeout time.Duration, message string) {
    t.Helper()
    deadline := time.Now().Add(timeout)
    for time.Now().Before(deadline) {
        if condition() {
            return
        }
        time.Sleep(10 * time.Millisecond)
    }
    t.Errorf("%s: condition not met within %v", message, timeout)
}
```

#### 2.2 Mock Framework
```go
// internal/testing/mocks/mocks.go
package mocks

import (
    "spooky/internal/interfaces"
)

// MockFactsIntegration provides a mock facts integration
type MockFactsIntegration struct {
    LoadFactsFunc func(machineNames []string) (*interfaces.FactsContext, error)
    CollectFactsFunc func(machineNames []string) (*interfaces.FactsContext, error)
}

// LoadFacts calls the mock LoadFacts function
func (m *MockFactsIntegration) LoadFacts(machineNames []string) (*interfaces.FactsContext, error) {
    if m.LoadFactsFunc != nil {
        return m.LoadFactsFunc(machineNames)
    }
    return nil, nil
}

// CollectFacts calls the mock CollectFacts function
func (m *MockFactsIntegration) CollectFacts(machineNames []string) (*interfaces.FactsContext, error) {
    if m.CollectFactsFunc != nil {
        return m.CollectFactsFunc(machineNames)
    }
    return nil, nil
}

// MockActionsIntegration provides a mock actions integration
type MockActionsIntegration struct {
    LoadActionsFunc func(projectPath string) (*interfaces.ActionsContext, error)
    ValidateActionFunc func(action interface{}, context *interfaces.ActionExecutionContext) error
}

// LoadActions calls the mock LoadActions function
func (m *MockActionsIntegration) LoadActions(projectPath string) (*interfaces.ActionsContext, error) {
    if m.LoadActionsFunc != nil {
        return m.LoadActionsFunc(projectPath)
    }
    return nil, nil
}

// ValidateAction calls the mock ValidateAction function
func (m *MockActionsIntegration) ValidateAction(action interface{}, context *interfaces.ActionExecutionContext) error {
    if m.ValidateActionFunc != nil {
        return m.ValidateActionFunc(action, context)
    }
    return nil
}
```

### Step 3: Integration Tests

#### 3.1 Integration Test Framework
```go
// internal/testing/integration/framework.go
package integration

import (
    "context"
    "testing"
    "time"
    "spooky/internal/coordinator"
    "spooky/internal/logging"
)

// IntegrationTestFramework provides integration testing framework
type IntegrationTestFramework struct {
    coordinator *coordinator.CoordinatorManager
    logger      logging.Logger
}

// NewIntegrationTestFramework creates a new integration test framework
func NewIntegrationTestFramework() *IntegrationTestFramework {
    logger := logging.NewConsoleLogger(logging.InfoLevel)
    
    return &IntegrationTestFramework{
        logger: logger,
    }
}

// SetupTest sets up a test environment
func (itf *IntegrationTestFramework) SetupTest(t *testing.T) *coordinator.CoordinatorManager {
    t.Helper()
    
    // Create coordinator for testing
    coord, err := coordinator.NewCoordinatorManagerFromProject("testdata/test-project", itf.logger)
    if err != nil {
        t.Fatalf("failed to create coordinator: %v", err)
    }
    
    itf.coordinator = coord
    return coord
}

// TeardownTest tears down a test environment
func (itf *IntegrationTestFramework) TeardownTest(t *testing.T) {
    t.Helper()
    
    if itf.coordinator != nil {
        // Clean up any resources
        if err := itf.coordinator.ClearAllCaches(); err != nil {
            t.Logf("failed to clear caches: %v", err)
        }
    }
}

// RunIntegrationTest runs an integration test
func (itf *IntegrationTestFramework) RunIntegrationTest(t *testing.T, testName string, testFunc func(*testing.T, *coordinator.CoordinatorManager)) {
    t.Run(testName, func(t *testing.T) {
        coord := itf.SetupTest(t)
        defer itf.TeardownTest(t)
        
        // Set timeout for integration tests
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        // Run test with context
        testFunc(t, coord)
    })
}
```

## Implementation Strategy

### Phase 1: Analysis and Planning (Week 1)
1. **Analyze current coverage** - Run coverage analysis across all packages
2. **Identify gaps** - Identify packages with low coverage
3. **Create test plan** - Create comprehensive test plan

### Phase 2: Infrastructure (Week 2)
1. **Implement test utilities** - Create test utility functions
2. **Create mock framework** - Implement mock framework
3. **Set up integration tests** - Create integration test framework

### Phase 3: Test Implementation (Week 3-4)
1. **Implement unit tests** - Add comprehensive unit tests
2. **Add integration tests** - Implement integration tests
3. **Add performance tests** - Implement performance tests

### Phase 4: Quality Assurance (Week 5)
1. **Review test quality** - Review and improve test quality
2. **Update documentation** - Update test documentation
3. **Monitor coverage** - Set up coverage monitoring

## Success Criteria

### Functional Requirements
- [ ] Minimum 80% coverage across all packages
- [ ] All edge cases covered
- [ ] Integration tests implemented
- [ ] Performance tests implemented

### Quality Requirements
- [ ] Test quality improved
- [ ] Test maintainability enhanced
- [ ] Test documentation complete
- [ ] Coverage monitoring in place

## Dependencies

### Required Dependencies
- All existing packages
- Testing framework (Go testing)
- Coverage tools
- Mocking framework

### Optional Dependencies
- Performance testing tools
- Coverage monitoring tools

## Risk Assessment

### High Risk
- **Time Investment**: Improving coverage requires significant time
- **Breaking Changes**: Tests may reveal existing bugs

### Medium Risk
- **Test Maintenance**: Tests need ongoing maintenance
- **Performance Impact**: Tests may impact build time

### Low Risk
- **Tool Integration**: Integration with existing tools
- **Documentation**: Documentation updates

## Next Steps

1. **Start with analysis** - Begin with coverage analysis
2. **Implement infrastructure** - Create test infrastructure
3. **Add tests gradually** - Add tests incrementally
4. **Monitor and maintain** - Monitor coverage and maintain tests
