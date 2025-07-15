.PHONY: build clean test run validate list help check-coverage test-unit test-integration test-coverage install-coverage-tool install-pre-commit-hook build-pre-commit-hook

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
	go build -o build/pre-commit.exe scripts/pre-commit.go
else
	go build -o build/pre-commit scripts/pre-commit.go
endif

# Install pre-commit hook (builds and copies to .git/hooks)
install-pre-commit-hook: build-pre-commit-hook
ifeq ($(OS),Windows_NT)
	copy build\pre-commit.exe .git\hooks\pre-commit.exe
else
	cp build/pre-commit .git/hooks/pre-commit
endif
	@echo "Pre-commit hook installed successfully"

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

# Default target
all: deps build test 