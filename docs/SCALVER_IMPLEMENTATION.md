# ScalVer Implementation in Spooky

## Overview

Spooky implements comprehensive support for [ScalVer (Scalable Versioning)](https://scalver.org/) as its primary versioning system. ScalVer combines the benefits of SemVer and CalVer to provide calendar-aware, semantic versioning that can adapt to different release cadences.

**Official Resources:**
- [ScalVer Specification](https://scalver.org/) - The official ScalVer specification
- [ScalVer GitHub Repository](https://github.com/scalver/scalver) - Reference implementation and documentation

## ScalVer Format

The ScalVer format is `MAJOR.DATE.PATCH` where:

- **MAJOR**: Mirrors SemVer's MAJOR component (0 = alpha/experimental, 1+ = stable)
- **DATE**: Calendar date in UTC (YYYY, YYYYMM, or YYYYMMDD)
- **PATCH**: Monotonically increasing counter for backward-compatible updates

### Supported Date Formats

- **YYYY**: Yearly cadence (e.g., `0.2025.0`)
- **YYYYMM**: Monthly cadence (e.g., `0.202508.0`)
- **YYYYMMDD**: Daily cadence (e.g., `0.20250812.0`)

### Development Versions

Development versions use the format `MAJOR.DATE.PATCH-dev-COMMIT` where:
- `COMMIT` is the short git commit hash (e.g., `0.20250812.0-dev-abc123`)

## Implementation

### Core Types

The ScalVer implementation is located in `internal/types/common/scalver.go`:

```go
type ScalVer struct {
    Major int
    Date  string // YYYY, YYYYMM, or YYYYMMDD
    Patch int
}
```

### Key Functions

#### Parsing and Validation

- `ParseScalVer(version string) (*ScalVer, error)`: Parse a ScalVer string
- `IsValidScalVerFormat(version string) bool`: Validate ScalVer format
- `isValidDateComponent(date string) bool`: Validate date component

#### Generation

- `GenerateScalVer(major int, datePrecision string, patch int) (string, error)`: Generate ScalVer version
- `GenerateDevelopmentScalVer(gitCommit string) (string, error)`: Generate development version

#### Comparison and Compatibility

- `Compare(other *ScalVer) int`: Compare two ScalVer versions
- `ValidateScalVerCompatibility(version1, version2 string) (bool, error)`: Check compatibility

#### Information and Analysis

- `GetScalVerInfo(version string) (map[string]interface{}, error)`: Get detailed version info
- `GetDatePrecision() string`: Get date precision (yearly/monthly/daily)
- `IsDevelopment() bool`: Check if version is development
- `IsStable() bool`: Check if version is stable

## Usage Examples

### Basic Parsing

```go
import spookytypescommon "spooky/internal/types/common"

// Parse a ScalVer version
scalver, err := spookytypescommon.ParseScalVer("0.20250812.0")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Major: %d, Date: %s, Patch: %d\n", 
    scalver.Major, scalver.Date, scalver.Patch)
```

### Version Generation

```go
// Generate versions with different cadences
yearly, _ := spookytypescommon.GenerateScalVer(0, "yearly", 0)
monthly, _ := spookytypescommon.GenerateScalVer(0, "monthly", 0)
daily, _ := spookytypescommon.GenerateScalVer(0, "daily", 0)

// Generate development version
devVersion, _ := spookytypescommon.GenerateDevelopmentScalVer("abc123")
```

### Version Comparison

```go
v1, _ := spookytypescommon.ParseScalVer("0.2025.0")
v2, _ := spookytypescommon.ParseScalVer("0.2025.1")

comparison := v1.Compare(v2)
switch comparison {
case -1:
    fmt.Println("v1 is less than v2")
case 0:
    fmt.Println("v1 equals v2")
case 1:
    fmt.Println("v1 is greater than v2")
}
```

### Compatibility Checking

```go
compatible, err := spookytypescommon.ValidateScalVerCompatibility(
    "0.2025.0", "0.2025.1")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Versions are compatible: %t\n", compatible)
```

## Integration Points

### Schema Validation

ScalVer validation is integrated into the schema system:

```go
// In internal/schemas/manager.go
func (m *Manager) isValidScalVerFormat(version string) bool {
    return spookytypescommon.IsValidScalVerFormat(version)
}
```

### Template Validation

Template version validation uses ScalVer:

```go
// In internal/templates/manager.go
func isValidVersion(version string) bool {
    return spookytypescommon.IsValidScalVerFormat(version)
}
```

### Build System

The build system supports ScalVer through the justfile:

```makefile
# Build with specific ScalVer version
build-scalver major date patch:
    go build -ldflags "-X spooky/cmd.Version={{major}}.{{date}}.{{patch}}" -o build/spooky main.go

# Build with yearly cadence
build-yearly:
    go build -ldflags "-X spooky/cmd.Version={{major}}.{{yearly_date}}.0" -o build/spooky main.go

# Build with monthly cadence
build-monthly:
    go build -ldflags "-X spooky/cmd.Version={{major}}.{{monthly_date}}.0" -o build/spooky main.go
```

## Testing

Comprehensive tests are provided in `internal/types/common/scalver_test.go`:

- Parsing tests for all supported formats
- Validation tests for valid and invalid versions
- Generation tests for different cadences
- Comparison and compatibility tests
- Development version tests

Run tests with:

```bash
go test ./internal/types/common -v
```

## Example Usage

See `examples/scalver-usage/main.go` for a complete demonstration of all ScalVer features.

Run the example with:

```bash
cd examples/scalver-usage && go run main.go
```

## Benefits

### Calendar-Aware
- You know when something was released
- Supports yearly, monthly, and daily cadences
- Date component only grows (DOG principle)

### Semantic Versioning Compatible
- All existing tooling works unchanged
- Major version indicates stability (0 = development, 1+ = stable)
- Patch version for backward-compatible updates

### Adjustable Cadence
- Can adapt from yearly → monthly → daily releases
- Supports different cadences for different components
- Maintains version ordering across cadences

### Development Support
- Clear distinction between development and stable versions
- Git commit integration for development builds
- Consistent versioning across all build types

## Migration from Other Versioning Systems

### From SemVer
- Major version maps directly
- Patch version maps directly
- Minor version becomes date component

### From CalVer
- Date component maps directly
- Add major version for stability indication
- Add patch version for ordering within date

### From Custom Systems
- Map stability indicator to major version
- Map release date to date component
- Map build number to patch version

## Best Practices

### Version Selection
- Use yearly cadence for stable releases
- Use monthly cadence for feature releases
- Use daily cadence for development builds

### Development Workflow
- Always use development versions (major = 0) for development
- Use stable versions (major = 1+) for production releases
- Include git commit in development versions

### Compatibility
- Major version changes indicate breaking changes
- Same major version indicates compatibility
- Development versions are compatible with each other

### Documentation
- Document the cadence used for each component
- Explain version meaning to users
- Provide migration guides for version changes

## Future Enhancements

### Planned Features
- Version range support (e.g., "0.2025.*")
- Version constraints (e.g., ">=0.2025.0")
- Version metadata (build info, dependencies)
- Version migration tools

### Integration Opportunities
- CI/CD pipeline integration
- Release automation
- Dependency management
- Version analytics

## Conclusion

ScalVer provides a robust, flexible versioning system that combines the best aspects of semantic and calendar versioning. The implementation in Spooky supports all ScalVer features while maintaining compatibility with existing tools and workflows.

**For more information about ScalVer, visit:**
- [ScalVer Specification](https://scalver.org/)
- [ScalVer GitHub Repository](https://github.com/scalver/scalver)
