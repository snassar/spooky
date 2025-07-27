# Makefile for spooky-generate tool

.PHONY: build clean test test-coverage help

# Default target
all: build

# Build the spooky-generate binary
build:
	@echo "Building spooky-generate tool..."
	go build -o spooky-generate
	@echo "Build complete! Binary created: ./spooky-generate"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f spooky-generate coverage.out coverage.html
	@echo "Clean complete!"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	@echo "Dependencies installed!"

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the spooky-generate binary (default)"
	@echo "  clean         - Remove build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  deps          - Install dependencies"
	@echo "  help          - Show this help message"
	@echo ""
	@echo "Usage examples:"
	@echo "  make build"
	@echo "  ./spooky-generate --help"
	@echo "  ./spooky-generate --actions 100 --machines 50"
	@echo "  ./spooky-generate --generate inventory-only --machines 1000" 