# Linting Scripts

This directory contains scripts for running golangci-lint on the spooky codebase in different ways to manage large numbers of linting issues.

## Go Programs

### `lint-todos.go` (AI-Optimized)
A Go program that can handle both single file and batch linting with comprehensive reporting, **optimized for AI consumption**.

**Usage:**
```bash
# Lint all files and generate AI-optimized report
go run scripts/lint-todos.go

# Lint a single file
LINT_FILE=./cmd/root.go go run scripts/lint-todos.go

# Customize output file
OUTPUT_FILE=./docs/plans/lint-todos.md go run scripts/lint-todos.go

# Adjust number of workers (default: CPU cores)
LINT_WORKERS=4 go run scripts/lint-todos.go

# Use fast mode (skips some linters for speed)
LINT_FAST=true go run scripts/lint-todos.go
```

**AI-Optimized Features:**
- **Issue Categorization** - Groups issues by type (typecheck, gocritic, etc.)
- **Priority Classification** - High (build-breaking), Medium (code quality), Low (style)
- **Actionable Fix Suggestions** - Specific code patterns and examples
- **Batch Fix Opportunities** - Identifies issues that appear in multiple files
- **Progress Tracking** - Shows completion status and recommended action plan
- **Context & Dependencies** - File relationships and impact analysis
- **Fix Patterns** - Common solutions with code examples
- **Concurrent processing** - Uses all CPU cores for maximum speed
- **Minimal output** - Only shows summary at the end

**Generated Report Sections:**
1. **📊 Summary Statistics** - Overview with priority breakdown
2. **🏷️ Issue Categories** - Grouped by linter type with descriptions
3. **📁 File-by-File Breakdown** - Detailed issues per file with priority levels
4. **🔧 Batch Fix Opportunities** - Similar issues across multiple files
5. **📈 Progress Tracking** - Current status and recommended action plan
6. **🛠️ Common Fix Patterns** - Code examples for common issues

## Bash Scripts

### `lint-single.sh`
Minimal script for use with `find -exec` to process files one by one.

**Usage:**
```bash
# Process all Go files
find . -name "*.go" -exec ./scripts/lint-single.sh {} \;

# Process specific files
./scripts/lint-single.sh ./cmd/root.go
```

**Output:**
- ✓ for files with no issues
- ✗ for files with issues, including TODO comments and the specific errors

### `lint-file.sh`
Simple script to lint a single file and show issues as TODO comments.

**Usage:**
```bash
./scripts/lint-file.sh ./cmd/root.go
```

### `lint-todos.sh`
Comprehensive bash script that processes all files in batches and generates a markdown report.

**Usage:**
```bash
./scripts/lint-todos.sh
```

**Features:**
- Processes files in configurable batches
- Generates a markdown report (`lint-todos.md`)
- Provides summary statistics
- Includes suggested fixes

## Just Targets

The following just targets are available:

```bash
# Run standard golangci-lint on entire codebase
just lint

# Run AI-optimized linter on each file individually and capture TODOs
just lint-todos

# Run linter on a specific file (Go program)
just lint-file <file.go>

# Run linter on all files using find -exec (captures TODOs)
just lint-all-files
```

## Why File-by-File Linting?

When a codebase has many linting issues, running `golangci-lint run` on the entire codebase can:

1. **Overwhelm the output** - Too many errors to process at once
2. **Timeout** - golangci-lint may give up on large codebases
3. **Make it hard to track progress** - Difficult to see what's been fixed

File-by-file linting helps by:

1. **Managing scope** - Focus on one file at a time
2. **Tracking progress** - Clear indication of which files need work
3. **Generating TODOs** - Structured output for systematic fixing
4. **Avoiding timeouts** - Each file is processed independently

## AI-Optimized Workflow

The enhanced `lint-todos.go` program provides an AI-friendly workflow:

### 1. **Issue Prioritization**
- **🔴 High Priority**: Fix build-breaking issues first (typecheck errors)
- **🟡 Medium Priority**: Address code quality issues (gocritic, govet)
- **🟢 Low Priority**: Handle style and formatting issues last

### 2. **Batch Processing**
- Identify similar issues across multiple files
- Apply systematic fixes using patterns
- Use sed/find commands for bulk changes

### 3. **Progress Tracking**
- Monitor completion percentage
- Track issues by priority level
- Update status as fixes are applied

### 4. **Fix Patterns**
- Common solutions for each issue type
- Code examples for implementation
- Command-line tools for automation

## Performance

**Concurrent Processing:**
- Uses worker pool pattern for maximum efficiency
- Configurable number of workers (default: CPU cores)
- ~2x speedup over sequential processing

**Fast Mode:**
- Use `LINT_FAST=true` to skip slower linters
- Reduces processing time by ~30%
- Still captures most important issues

**Memory Usage:**
- Processes files in memory-efficient batches
- Minimal memory footprint per worker
- Suitable for large codebases

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `LINT_CONFIG` | `.golangci.yml` | Path to golangci-lint config |
| `OUTPUT_FILE` | `lint-todos.md` | Output file path |
| `LINT_FILE` | `""` | Single file to lint (if set) |
| `LINT_WORKERS` | `runtime.NumCPU()` | Number of concurrent workers |
| `LINT_FAST` | `true` | Use fast mode (skip slow linters) |

## Example Output

The AI-optimized report includes:

```markdown
# Linting TODOs - AI-Optimized Report

## 📊 Summary Statistics
- **Total files processed:** 88
- **Files with issues:** 54
- **Total issues:** 223

### Priority Breakdown
- **🔴 High Priority (Build-breaking):** 149 issues (66.8%)
- **🟡 Medium Priority (Code quality):** 73 issues (32.7%)
- **🟢 Low Priority (Style):** 1 issues (0.4%)

## 🏷️ Issue Categories
### typecheck (103 issues)
**Description:** Build-breaking errors that prevent compilation
**Common Fix Patterns:**
- Add missing import: `import "package/path"`
- Define missing variable: `var variableName Type`

## 🔧 Batch Fix Opportunities
### Issues Found in Multiple Files
#### revive (18 files)
**Pattern:** unused-parameter: parameter 'ctx' seems to be unused
**Affected files:**
- internal/cli/commands/integrations.go:25
- internal/integration/manager.go:140

## 📈 Progress Tracking
### Recommended Action Plan
1. **🔴 Fix all High Priority issues first** (149 issues) - prevents compilation
2. **🟡 Address Medium Priority issues** (73 issues) - improves code quality
3. **🟢 Fix Low Priority issues** (1 issues) - style and formatting
```

This format makes it much easier for AI tools to understand the issues, prioritize fixes, and provide actionable solutions.
