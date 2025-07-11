.PHONY: build clean test run validate list help check-coverage test-unit test-integration test-coverage install-coverage-tool install-pre-commit-hook build-pre-commit-hook

# Build the spooky binary
build:
	go mod tidy
	go build -o spooky

# Clean build artifacts
clean:
	rm -f spooky
	rm -f coverage.out coverage-integration.out
	rm -f coverage.html coverage-integration.html
	go clean -testcache

# Run tests
test: test-unit test-integration check-coverage 

# Run unit tests only
test-unit:
	go test -v ./...

# Run integration tests only
test-integration:
	go test -v -tags=integration ./tests/integration/...

# Run all tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go test -v -tags=integration -coverprofile=coverage-integration.out ./tests/integration/...
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -html=coverage-integration.out -o coverage-integration.html

# Install go-test-coverage tool
install-coverage-tool:
	go install github.com/vladopajic/go-test-coverage/v2@latest

# Check test coverage using go-test-coverage tool
# Note: Uses 'go run' to avoid requiring local installation
check-coverage:
	go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
	go run github.com/vladopajic/go-test-coverage/v2@latest --config=./tests/testcoverage.yml

# Generate HTML coverage report locally
coverage-html:
	go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
	go tool cover -html=./cover.out -o coverage.html
	@echo "HTML coverage report generated: coverage.html"
	@echo "Open coverage.html in your browser to view the report"

# Run the tool with example configuration
run: build
	./spooky execute example.hcl

# Validate example configuration
validate: build
	./spooky validate example.hcl

# List servers and actions from example configuration
list: build
	./spooky list example.hcl

# Show help
help: build
	./spooky --help

# Install dependencies
deps:
	go mod download

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Create release binary
release: clean
	GOOS=linux GOARCH=amd64 go build -o spooky-linux-amd64
	GOOS=darwin GOARCH=amd64 go build -o spooky-darwin-amd64
	GOOS=windows GOARCH=amd64 go build -o spooky-windows-amd64.exe

# Build pre-commit hook
build-pre-commit-hook:
ifeq ($(OS),Windows_NT)
	go build -o scripts/pre-commit.exe scripts/pre-commit.go
else
	go build -o scripts/pre-commit scripts/pre-commit.go
endif

# Install pre-commit hook (builds and copies to .git/hooks)
install-pre-commit-hook: build-pre-commit-hook
ifeq ($(OS),Windows_NT)
	copy scripts\pre-commit.exe .git\hooks\pre-commit
else
	cp scripts/pre-commit .git/hooks/pre-commit
endif
	@echo "Pre-commit hook installed successfully"

# Default target
all: deps build test 