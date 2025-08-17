# Scripts Directory

This directory contains shell scripts for build automation, deployment, and CI/CD processes.

## Script Organization

The Go-based developer tools have been moved to the `tools/` directory. This directory now contains only shell scripts for:

- Build automation
- Deployment processes
- CI/CD workflows
- System administration tasks

## Available Scripts

### Build Scripts
- `build.sh` - Main build automation script
- `deploy.sh` - Deployment automation script

### CI/CD Scripts
- `ci/` - Continuous integration and deployment scripts

## Usage

```bash
# Run build script
./scripts/build.sh

# Run deployment script
./scripts/deploy.sh

# Run CI scripts
./scripts/ci/setup.sh
```

## Tool Integration

For Go-based developer tools, see the `tools/` directory:

```bash
# Run developer tools
just build-tools
just run-tools

# Individual tools
just --directory tools/dependency-graph run
just --directory tools/todo-linter run
```

## Migration Notes

The following Go scripts were moved to `tools/` as individual tools:

- `generate-adrs.go` → `tools/adr-generator/`
- `analyze-git-adrs.go` → `tools/git-adr-analyzer/`
- `focused-adr-analysis.go` → `tools/focused-adr-analyzer/`
- `generate-adrs-simple.go` → `tools/simple-adr-generator/`
- `lint-todos.go` → `tools/todo-linter/`

Each tool now has its own justfile, documentation, and can be used independently or through the tools orchestration system.
