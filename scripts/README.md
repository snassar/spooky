# Scripts Directory

This directory contains utility scripts for the spooky project.

## Pre-commit Hook

The `pre-commit.go` file is a Go script that runs test coverage checks before each commit.

### Building the Pre-commit Hook

To build the pre-commit hook binary:

```bash
# From the project root (cross-platform)
make build-pre-commit-hook

# Manual build (Unix/Linux/macOS)
go build -o scripts/pre-commit scripts/pre-commit.go

# Manual build (Windows)
go build -o scripts/pre-commit.exe scripts/pre-commit.go

# Or from the scripts directory
cd scripts
go build pre-commit.go
```

### Installing the Pre-commit Hook

After building, install the hook:

```bash
# Automated install (cross-platform)
make install-pre-commit-hook

# Manual install (Unix/Linux/macOS)
cp scripts/pre-commit .git/hooks/pre-commit

# Manual install (Windows)
copy scripts\pre-commit.exe .git\hooks\pre-commit
```

### What the Hook Does

The pre-commit hook:
1. Checks if you're in a git repository
2. Identifies staged Go files
3. Runs tests with coverage profiling
4. Verifies coverage meets thresholds defined in `tests/testcoverage.yml`
5. Blocks commits if coverage is insufficient
6. Allows commits if coverage passes

### Manual Usage

You can also run the coverage check manually:

```bash
go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./... && go run github.com/vladopajic/go-test-coverage/v2@latest --config=./tests/testcoverage.yml
```

### Requirements

- Go 1.24 or later
- The `go-test-coverage/v2` tool (installed automatically via `go run`)
- Coverage configuration in `tests/testcoverage.yml` 