# ScalVer Versioning Implementation

This package implements [ScalVer (Scalable Calendar Versioning)](https://scalver.org/) for the spooky project.

## Overview

ScalVer is a calendar-aware, SemVer-compatible versioning scheme expressed as `<MAJOR>.<DATE>.<PATCH>` where the `DATE` segment may lengthen over time within a `MAJOR` line: `YYYY` → `YYYYMM` → `YYYYMMDD`.

## Current Version

The current version of spooky is: **0.20250905.0**

This means:
- **Major**: 0 (alpha/experimental phase)
- **Date**: 20250905 (September 5, 2025)
- **Patch**: 0 (first release of this date)

## Usage

### Basic Version Information

```go
import "spooky/internal/version"

// Get current version
currentVersion := version.GetCurrentVersion()
fmt.Println(currentVersion.String()) // "0.20250905.0"

// Get full version info with build metadata
versionInfo := version.GetFullVersionInfo()
fmt.Println(versionInfo.String())
```

### Building with Specific Versions

```bash
# Build with default version (0.20250905.0)
just build-version

# Build with custom version
just build-version "1.20250905.0"

# Build with dev version including git commit
just build-version "0.20250905.0-dev-$(git rev-parse --short HEAD)"
```

### Version Comparison

```go
v1, _ := version.ParseScalVer("0.20250905.0")
v2, _ := version.ParseScalVer("0.20250905.1")

if v1.Less(v2) {
    fmt.Println("v1 is older than v2")
}
```

## ScalVer Rules

### Date-Only-Grows (DOG) Rule

Within any single MAJOR line, the DATE can stay the same length or grow, but never shrink:

- ✅ `1.2025.0` → `1.202503.0` (yearly → monthly)
- ✅ `1.202503.0` → `1.20250301.0` (monthly → daily)
- ❌ `1.20250301.0` → `1.2026.0` (daily → yearly, requires MAJOR bump)

### Major Version Bumps

Bump MAJOR for:
- Breaking changes
- When DATE would otherwise need to contract

### Examples

| Version | Meaning |
|---------|---------|
| `0.20250905.0` | First alpha release on September 5, 2025 |
| `0.20250905.1` | Second alpha release on same day |
| `1.20250905.0` | First stable release on September 5, 2025 |
| `1.202509.0` | First September 2025 release (monthly cadence) |
| `1.2025.0` | First 2025 release (yearly cadence) |

## Pre-release Identifiers

ScalVer supports SemVer pre-release identifiers:

- `0.20250905.0-dev-abc123`
- `0.20250905.0-alpha.1`
- `0.20250905.0-beta.2`
- `0.20250905.0-rc.1`

## Build Metadata

Build metadata is automatically included:

- Build time (UTC)
- Git commit hash
- Go version
- Build info (dev/release)

Example output:
```
spooky 0.20250905.0
Automation and configuration management tool
Build: dev
Built: 2025-09-04T22:18:41Z
Commit: 05cf558f3a2988a117e076cde534e288708e041e
Go: go1.24.6
```

## Migration from SemVer

Since every ScalVer tag is syntactically valid SemVer, existing tooling works unchanged:

- Package managers (npm, Cargo, Maven, etc.)
- CI/CD systems
- Release dashboards
- Version comparison tools

## Testing

Run the version package tests:

```bash
go test ./internal/version/... -v
```

## References

- [ScalVer Specification](https://scalver.org/)
- [Semantic Versioning 2.0](https://semver.org/)
- [Calendar Versioning](https://calver.org/)
