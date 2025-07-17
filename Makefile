.PHONY: build clean test run validate list help check-coverage test-unit test-integration test-coverage install-coverage-tool install-pre-commit-hook build-pre-commit-hook clean-config generate-config

# Build the spooky binary
build:
	go mod tidy
	go build -o build/spooky

# Clean build artifacts
clean:
	rm -rf build/
	go clean -testcache

# Install dependencies
get-dependencies:
	go mod download

# Install development tools
install-development-tools:
	go install github.com/vladopajic/go-test-coverage/v2@latest
	go install github.com/wadey/gocovmerge@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Create release binary
release: clean
	GOOS=linux GOARCH=amd64 go build -o build/spooky-linux-amd64
	GOOS=darwin GOARCH=amd64 go build -o build/spooky-darwin-amd64
	GOOS=windows GOARCH=amd64 go build -o build/spooky-windows-amd64.exe

# Build pre-commit hook
build-pre-commit-hook:
ifeq ($(OS),Windows_NT)
	go build -o build/pre-commit.exe tools/pre-commit/main.go
else
	go build -o build/pre-commit tools/pre-commit/main.go
endif

# Install pre-commit hook (builds and copies to .git/hooks)
install-pre-commit-hook: build-pre-commit-hook
ifeq ($(OS),Windows_NT)
	copy build\pre-commit.exe .git\hooks\pre-commit.exe
else
	cp build/pre-commit .git/hooks/pre-commit
endif
	@echo "Pre-commit hook installed successfully"

# Build test environment tool
build-test-env:
ifeq ($(OS),Windows_NT)
	go build -o build/spooky-test-env.exe tools/spooky-test-env/main.go
else
	go build -o build/spooky-test-env tools/spooky-test-env/main.go
endif
	@echo "Test environment tool built successfully"

# Install test environment tool
install-test-env: build-test-env
ifeq ($(OS),Windows_NT)
	cp build/spooky-test-env.exe ~/spooky-test-env/spooky-test-env.exe
else
	cp build/spooky-test-env ~/spooky-test-env/spooky-test-env
	chmod +x ~/spooky-test-env/spooky-test-env
endif
	@echo "Test environment tool installed successfully"

# Run tests
test: test-unit test-integration check-coverage 

# Run unit tests only
test-unit:
	go test -v ./...

# Run unit tests with coverage and go-test-coverage
test-unit-coverage:
	go test ./... -coverprofile=./tests/coverage.out -covermode=atomic -coverpkg=./...
	go run github.com/vladopajic/go-test-coverage/v2@latest --config=./tests/testcoverage.yml

# Run integration tests only
test-integration:
	go test -v -tags=integration ./tests/integration/...

# Run Podman-based integration tests
test-integration-podman:
	go test -v -podman ./tests/integration/podman_integration_test.go

# Run basic Podman environment tests
test-podman-basic:
	go test -v -podman-basic ./tests/integration/podman_basic_test.go

# Run all integration tests (legacy + Podman)
test-integration-all: test-integration test-integration-podman

# Run all tests with coverage
test-coverage:
	go test -v -coverprofile=./tests/coverage-unit.out ./...
	cd tests && go test -v -tags=integration -coverprofile=./coverage-integration.out ./...
	gocovmerge ./tests/coverage-unit.out ./tests/coverage-integration.out > ./tests/coverage.out
	go tool cover -html=./tests/coverage.out -o ./tests/reports/coverage.html

# Check test coverage using go-test-coverage tool
# Note: Uses 'go run' to avoid requiring local installation
check-coverage:
	go test ./... -coverprofile=./tests/coverage.out -covermode=atomic -coverpkg=./...
	go run github.com/vladopajic/go-test-coverage/v2@latest --config=./tests/testcoverage.yml

# Generate HTML coverage report locally
coverage-html:
	go test ./... -coverprofile=./tests/coverage.out -covermode=atomic -coverpkg=./...
	go tool cover -html=./tests/coverage.out -o tests/reports/coverage.html
	@echo "HTML coverage report generated: tests/reports/coverage.html"
	@echo "Open tests/reports/coverage.html in your browser to view the report"

# Run the tool with example configuration
run: build
	./build/spooky execute examples/configuration/example.hcl

# Validate example configuration
validate: build
	./build/spooky validate examples/configuration/example.hcl

# List servers and actions from example configuration
list: build
	./build/spooky list examples/configuration/example.hcl

# Show help
help: build
	./build/spooky --help

# Clean up generated configuration files
clean-config:
	@echo "Cleaning up generated configuration files..."
	@find examples/configuration/ -name "*-scale-example-*.hcl" -type f -delete
	@echo "Generated configuration files cleaned up"

# Generate configuration files for testing
generate-config:
	@echo "Generating configuration files for testing..."
	@go run tools/generate-config/main.go
	@echo "Configuration files generated successfully"

# Generate and clean config files (clean first, then generate)
config: clean-config generate-config

# Default target
all: deps build test 