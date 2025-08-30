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
    rm -rf tools/spooky-test-bundle-generator/bundles/
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

build-facts:
    #!/usr/bin/env bash
    echo "🔧 Building spooky-facts gatherer..."
    cd tools/spooky-facts-gatherer
    go build -o spooky-facts .
    echo "✅ Facts gatherer built: tools/spooky-facts-gatherer/spooky-facts"

test-facts:
    #!/usr/bin/env bash
    echo "🧪 Testing facts gatherer..."
    cd tools/spooky-facts-gatherer
    ./spooky-facts preview
    echo "✅ Facts gatherer test completed"

gather-facts:
    #!/usr/bin/env bash
    echo "🔍 Gathering system facts..."
    cd tools/spooky-facts-gatherer
    ./spooky-facts gather
    echo "✅ Facts gathered and saved"

gather-facts-verbose:
	#!/usr/bin/env bash
	echo "🔍 Gathering system facts (verbose mode)..."
	cd tools/spooky-facts-gatherer
	./spooky-facts gather --verbose
	echo "✅ Facts gathered and saved"

# Gather facts from remote machines
gather-remote-facts:
	#!/usr/bin/env bash
	echo "🌐 Gathering facts from remote machines..."
	./build/spooky facts gather

# Gather facts from remote machines with custom output
gather-remote-facts-to output:
	#!/usr/bin/env bash
	echo "🌐 Gathering facts from remote machines..."
	./build/spooky facts gather {{output}}

# Clean test projects
clean-tests:
    #!/usr/bin/env bash
    echo "🧹 Cleaning test projects..."
    rm -rf testing/
    mkdir -p testing
    echo "✅ Test projects cleaned"

# Multi-stage build for spooky-test-bundle-generator
build-spooky-test-bundle-generator:
    #!/usr/bin/env bash
    echo "🔨 Building spooky-test-bundle-generator..."
    
    # Stage 1: Build spooky binary for bundle generator
    echo "📦 Stage 1: Building spooky binary..."
    go build -o tools/spooky-test-bundle-generator/spooky main.go
    if [ $? -ne 0 ]; then
        echo "❌ Failed to build spooky binary"
        exit 1
    fi
    echo "✅ Spooky binary built"
    
    # Stage 2: Build bundle generator
    echo "📦 Stage 2: Building bundle generator..."
    cd tools/spooky-test-bundle-generator
    go build -o spooky-test-bundle-generator main.go
    if [ $? -ne 0 ]; then
        echo "❌ Failed to build bundle generator"
        exit 1
    fi
    echo "✅ Bundle generator built"
    
    echo "🎉 Multi-stage build completed successfully!"
    echo "📁 Output: tools/spooky-test-bundle-generator/spooky-test-bundle-generator"

# Alternative: Separate stages for more control
build-spooky-for-bundles:
    #!/usr/bin/env bash
    echo "🔨 Building spooky binary for bundle generator..."
    go build -o tools/spooky-test-bundle-generator/spooky main.go
    echo "✅ Spooky binary built"

build-bundle-generator: build-spooky-for-bundles
    #!/usr/bin/env bash
    echo "🔨 Building bundle generator..."
    cd tools/spooky-test-bundle-generator
    go build -o spooky-test-bundle-generator main.go
    echo "✅ Bundle generator built"

# Run bundle generator
run-bundle-generator profile output: build-spooky-test-bundle-generator
    #!/usr/bin/env bash
    echo "🔧 Running bundle generator..."
    cd tools/spooky-test-bundle-generator
    ./spooky-test-bundle-generator generate --profile {{profile}} --output {{output}}
    echo "✅ Bundle generated: {{output}}"

# Clean bundle generator artifacts
clean-bundle-generator:
    #!/usr/bin/env bash
    echo "🧹 Cleaning bundle generator artifacts..."
    rm -f tools/spooky-test-bundle-generator/spooky
    rm -f tools/spooky-test-bundle-generator/spooky-test-bundle-generator
    echo "✅ Cleanup completed"

# Clean test bundle generated files
clean-test-bundles:
    #!/usr/bin/env bash
    echo "🧹 Cleaning test bundle generated files..."
    rm -rf tools/spooky-test-bundle-generator/bundles/
    echo "✅ Test bundle files cleaned"

# Show project statistics
stats:
    #!/usr/bin/env bash
    echo "📊 Project Statistics:"
    echo "Go files (excluding test bundles): $(find . -name "*.go" -not -path "./tools/spooky-test-bundle-generator/bundles/*" | wc -l)"
    echo "Total lines (excluding test bundles): $(find . -name "*.go" -not -path "./tools/spooky-test-bundle-generator/bundles/*" | xargs wc -l | tail -1)"
    echo "Test files (excluding test bundles): $(find . -name "*_test.go" -not -path "./tools/spooky-test-bundle-generator/bundles/*" | wc -l)"
    echo "Documentation files: $(find . -name "*.md" | wc -l)"
    echo "Build artifacts: $(find build/ -type f 2>/dev/null | wc -l || echo 0)"
    echo "Test artifacts: $(find testing/ -type f 2>/dev/null | wc -l || echo 0)"
    echo ""
    echo "📦 Test Bundle Files (excluded from main stats):"
    echo "Test bundle Go files: $(find ./tools/spooky-test-bundle-generator/bundles -name "*.go" 2>/dev/null | wc -l || echo 0)"

# Show all available commands
list:
    @just --list
