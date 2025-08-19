# Spooky project justfile
# Run with: just <command>

# Default target
default:
    @just --list

# Build the spooky binary (development build)
build: generate-version
    @echo "Building spooky binary (development)..."
    go build -ldflags "-X spooky/cmd.Version={{dev_version}} -X spooky/cmd.BuildTime={{build_time}} -X spooky/cmd.GitCommit={{git_commit}}" -o build/spooky main.go
    @echo "Development build complete: build/spooky"

# Build for development (with dev version)
build-dev:
    @echo "Building spooky binary (dev version)..."
    go build -ldflags "-X spooky/cmd.Version=dev -X spooky/cmd.BuildTime={{build_time}} -X spooky/cmd.GitCommit={{git_commit}}" -o build/spooky main.go
    @echo "Dev build complete: build/spooky"

# Generate version information
generate-version:
    @echo "Generating version information..."
    @mkdir -p build
    @echo "{{version}}" > build/version.txt
    @echo "{{build_time}}" > build/build-time.txt
    @echo "{{git_commit}}" > build/git-commit.txt

# Clean build artifacts
clean:
    @echo "Cleaning build artifacts..."
    rm -rf build/
    go clean -cache

# Run all tests
test:
    @echo "Running tests..."
    @mkdir -p build
    go test ./... -v -race -coverprofile=build/coverage.out

# Run tests with coverage report
test-coverage: test
    @echo "Coverage report:"
    go tool cover -func=build/coverage.out | tail -1
    @echo "Generating HTML coverage report..."
    go tool cover -html=build/coverage.out -o build/coverage.html

# Run unit tests only
test-unit:
    @echo "Running unit tests..."
    @mkdir -p build
    go test ./... -v -race -coverprofile=build/coverage.out -run "^Test"

# Run integration tests only
test-integration:
    @echo "Running integration tests..."
    @mkdir -p build
    go test ./... -v -race -coverprofile=build/coverage.out -run "^TestIntegration"

# Run linter
lint:
    @echo "Running linter..."
    golangci-lint run --timeout=5m

# Run terminology checker
lint-terminology:
    @echo "Running terminology checker..."
    @cd tools/terminology-checker && just run

# Run linter on each file individually and capture TODOs (AI-optimized)
lint-todos:
    @echo "Running AI-optimized file-by-file linting analysis..."
    @cd tools/todo-linter && just run

# Run AI-optimized linter on a specific file
lint-file file:
    @echo "Running AI-optimized linter on {{file}}..."
    @cd tools/todo-linter && LINT_FILE={{file}} just run



# Format code
fmt:
    @echo "Formatting code..."
    gofmt -s -w .

# Check code formatting
check-fmt:
    @echo "Checking code formatting..."
    @bash -c 'UNFORMATTED_FILES=$(gofmt -s -l . 2>/dev/null || true); if [ -n "$$UNFORMATTED_FILES" ]; then echo "Code is not formatted. Please run '\''just fmt'\''"; echo "Unformatted files: $$UNFORMATTED_FILES"; exit 1; fi; echo "Code formatting is correct"'

# Check for unused dependencies
check-deps:
    @echo "Checking for unused dependencies..."
    @go mod tidy
    @bash -c 'if [ -n "$$(git status --porcelain)" ]; then echo "go.mod or go.sum has uncommitted changes. Please run '\''go mod tidy'\''"; git status; exit 1; fi; echo "Dependencies are clean"'

# Run all checks (lint, format, deps, tests)
check: check-fmt check-deps lint lint-terminology test

# Install development tools
install-tools:
    @echo "Installing development tools..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install github.com/vladopajic/go-test-coverage/v2@latest

# Show version information
version:
    @echo "Development Version: {{dev_version}}"
    @echo "Official Version: {{version}}"
    @echo "Major: {{major}}"
    @echo "Date: {{date}}"
    @echo "Patch: {{patch}}"
    @echo "Build Time: {{build_time}}"
    @echo "Git Commit: {{git_commit}}"

# Build with specific ScalVer version
build-scalver major date patch:
    @echo "Building with ScalVer {{major}}.{{date}}.{{patch}}..."
    go build -ldflags "-X spooky/cmd.Version={{major}}.{{date}}.{{patch}} -X spooky/cmd.BuildTime={{build_time}} -X spooky/cmd.GitCommit={{git_commit}}" -o build/spooky main.go
    @echo "ScalVer build complete: {{major}}.{{date}}.{{patch}}"

# Build official release from git tag
build-release:
    @echo "Building official release from git tag..."
    @TAG=$$(git describe --tags --exact-match 2>/dev/null || echo "not-tagged"); \
    if [ "$$TAG" = "not-tagged" ]; then \
        echo "Error: Current commit is not tagged. Use 'git tag' to create a release tag."; \
        exit 1; \
    fi; \
    echo "Building release: $$TAG"; \
    go build -ldflags "-X spooky/cmd.Version=$$TAG -X spooky/cmd.BuildTime={{build_time}} -X spooky/cmd.GitCommit={{git_commit}}" -o build/spooky main.go; \
    echo "Release build complete: $$TAG"

# Build with yearly cadence (YYYY)
build-yearly:
    @echo "Building with yearly ScalVer..."
    go build -ldflags "-X spooky/cmd.Version={{major}}.{{yearly_date}}.0 -X spooky/cmd.BuildTime={{build_time}} -X spooky/cmd.GitCommit={{git_commit}}" -o build/spooky main.go
    @echo "Yearly build complete: {{major}}.{{yearly_date}}.0"

# Build with monthly cadence (YYYYMM)
build-monthly:
    @echo "Building with monthly ScalVer..."
    go build -ldflags "-X spooky/cmd.Version={{major}}.{{monthly_date}}.0 -X spooky/cmd.BuildTime={{build_time}} -X spooky/cmd.GitCommit={{git_commit}}" -o build/spooky main.go
    @echo "Monthly build complete: {{major}}.{{monthly_date}}.0"

# Show help
help:
    @just --list

# Variables
# ScalVer format: MAJOR.DATE.PATCH where DATE can be YYYY, YYYYMM, or YYYYMMDD
major := "0"
date := `date +%Y%m%d`
yearly_date := `date +%Y`
monthly_date := `date +%Y%m`
# Development builds always use patch 0 with git commit suffix
dev_version := major + "." + date + ".0-dev-" + `git rev-parse --short HEAD 2>/dev/null || echo "unknown"`
# Official builds use incremented patch (for tagged releases)
patch := "0"
version := major + "." + date + "." + patch
build_time := `date -u +"%Y-%m-%dT%H:%M:%SZ"`
git_commit := `git rev-parse --short HEAD 2>/dev/null || echo "unknown"`

# Show current git status
git-status:
    @echo "Git status:"
    git status --porcelain || echo "No git repository found"

# Show current branch and commit
git-info:
    @echo "Current branch: $(git branch --show-current 2>/dev/null || echo 'unknown')"
    @echo "Current commit: $(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
    @echo "Commit message: $(git log -1 --pretty=format:'%s' 2>/dev/null || echo 'unknown')"

# Generate Architecture Decision Records from git history
generate-adrs:
    @echo "Generating Architecture Decision Records..."
    @cd tools/adr-generator && just run
    @echo "ADRs generated in docs/adr/data/"

# Analyze git history for potential ADRs
analyze-git-adrs:
    @echo "Analyzing git history for potential ADRs..."
    @cd tools/git-adr-analyzer && just run
    @echo "Analysis complete. Check docs/adr/analysis/git-analysis.md"

# Generate focused ADR recommendations
focused-adr-analysis:
    @echo "Generating focused ADR recommendations..."
    @cd tools/focused-adr-analyzer && just run
    @echo "Analysis complete. Check docs/adr/recommendations/focused-adr-recommendations.md"

# Generate comprehensive git analysis
analyze-git-comprehensive:
    @echo "Generating comprehensive git analysis..."
    @cd tools/git-adr-analyzer && just run
    @echo "Analysis complete. Check docs/adr/analysis/git-analysis.md"

# Generate simple ADRs with full paths
generate-adrs-simple:
    @echo "Generating simple ADRs..."
    @cd tools/simple-adr-generator && just run
    @echo "ADRs generated in docs/adr/"

# Generate dependency graphs
dependency-graph:
    @echo "Generating dependency graphs..."
    cd tools/dependency-graph && just analyze
    @echo "Dependency graphs generated in docs/dependencies/"

# Generate specific dependency graph
dependency-graph-package pkg:
    @echo "Generating dependency graph for package: {{pkg}}"
    cd tools/dependency-graph && just package {{pkg}}
    @echo "Package dependency graph generated in docs/dependencies/"

# Generate interface dependency graph
dependency-graph-interface:
    @echo "Generating interface dependency graph..."
    cd tools/dependency-graph && just interface
    @echo "Interface dependency graph generated in docs/dependencies/"

# Generate type dependency graph
dependency-graph-type:
    @echo "Generating type dependency graph..."
    cd tools/dependency-graph && just type
    @echo "Type dependency graph generated in docs/dependencies/"

# Generate function dependency graph
dependency-graph-function:
    @echo "Generating function dependency graph..."
    cd tools/dependency-graph && just function
    @echo "Function dependency graph generated in docs/dependencies/"

# Build all developer tools
build-tools:
    @echo "Building all developer tools..."
    @cd tools && just build

# Run all developer tools
run-tools:
    @echo "Running all developer tools..."
    @cd tools && just run-all

# Clean all developer tools
clean-tools:
    @echo "Cleaning all developer tools..."
    @cd tools && just clean

# Show tools help
tools-help:
    @echo "Developer Tools Help:"
    @cd tools && just help

# Development workflow: build, test, and show version
dev: build test version
    @echo "Development build complete!"

# Release workflow: clean, check, build, test
release: clean check build test version
    @echo "Release build complete!"
    @echo "Binary ready: build/spooky"

# Issue management commands are now in docs/issues/justfile
# Use: cd docs/issues && just <command>
