.PHONY: build clean test run validate list help

# Build the spooky binary
build:
	go mod tidy
	go build -o spooky

# Clean build artifacts
clean:
	rm -f spooky
	go clean

# Run tests
test:
	go test ./...

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

# Default target
all: deps build test 