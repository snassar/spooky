# Spooky Project Justfile
# A task runner for the Spooky automation framework

# Default recipe - show available commands
default:
    @just --list

# Build the spooky binary
build:
    #!/usr/bin/env bash
    echo "🔨 Building spooky binary..."
    mkdir -p build
    go build -o build/spooky .
    echo "✅ Build complete: ./build/spooky"

# Run all tests
test:
    #!/usr/bin/env bash
    echo "🧪 Running all tests..."
    go test -v ./...
    echo "✅ Tests complete"

# Run tests with coverage
test-coverage:
    #!/usr/bin/env bash
    echo "📊 Running tests with coverage..."
    mkdir -p testing
    go test -v -coverprofile=testing/coverage.out ./...
    go tool cover -html=testing/coverage.out -o testing/coverage.html
    echo "✅ Coverage report generated: testing/coverage.html"

# Run linting with golangci-lint
lint:
    #!/usr/bin/env bash
    echo "🔍 Running golangci-lint..."
    export PATH=$PATH:$(go env GOPATH)/bin
    golangci-lint run
    echo "✅ Linting complete"

# Format code with gofmt
fmt:
    #!/usr/bin/env bash
    echo "🎨 Formatting code..."
    gofmt -s -w .
    echo "✅ Code formatting complete"

# Organize imports
imports:
    #!/usr/bin/env bash
    echo "📦 Organizing imports..."
    export PATH=$PATH:$(go env GOPATH)/bin
    goimports -w .
    echo "✅ Import organization complete"

# Run all code quality checks
check: fmt imports lint
    #!/usr/bin/env bash
    echo "✅ All code quality checks complete"

# Run full CI pipeline
ci: check test
    #!/usr/bin/env bash
    echo "✅ CI pipeline complete"

# Clean build artifacts
clean:
    #!/usr/bin/env bash
    echo "🧹 Cleaning build artifacts..."
    rm -rf build/
    rm -rf testing/
    go clean -cache -testcache
    echo "✅ Clean complete"

# Install development dependencies
install-dev:
    #!/usr/bin/env bash
    echo "📦 Installing development dependencies..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install golang.org/x/tools/cmd/goimports@latest
    echo "✅ Development dependencies installed"

# Create test project
test-project name:
    #!/usr/bin/env bash
    echo "🧪 Creating test project: {{name}}"
    mkdir -p testing
    cd testing
    ../build/spooky project init --name {{name}} --description "Test project for {{name}}" {{name}}
    echo "✅ Test project created: testing/{{name}}"

# Clean test projects
clean-tests:
    #!/usr/bin/env bash
    echo "🧹 Cleaning test projects..."
    rm -rf testing/
    mkdir -p testing
    echo "✅ Test projects cleaned"

# Show project statistics
stats:
    #!/usr/bin/env bash
    echo "📊 Project Statistics:"
    echo "Go files: $(find . -name "*.go" | wc -l)"
    echo "Total lines: $(find . -name "*.go" | xargs wc -l | tail -1)"
    echo "Test files: $(find . -name "*_test.go" | wc -l)"
    echo "Documentation files: $(find . -name "*.md" | wc -l)"
    echo "Build artifacts: $(find build/ -type f 2>/dev/null | wc -l || echo 0)"
    echo "Test artifacts: $(find testing/ -type f 2>/dev/null | wc -l || echo 0)"

# Show all available commands
list:
    @just --list
