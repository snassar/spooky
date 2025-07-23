# Custom Fact Import and Merging Implementation Plan

## Overview

This document outlines the implementation plan for [issue #58](https://github.com/snassar/spooky/issues/58) - implementing custom fact import and merging via JSON (local and remote sources).

## Background

Users need to import custom facts from external sources and merge them with automatically collected facts. This enables integration with external systems, manual fact overrides, and hybrid fact collection strategies.

## Current State Analysis

### ✅ Already Implemented

1. **Basic Import Infrastructure**
   - `spooky facts import` command with basic functionality
   - `ImportCustomFacts()` method in Manager
   - HTTP collector for remote JSON imports
   - JSON collector for local file imports
   - Fact merging logic with multiple policies (replace, merge, skip, append)
   - Basic CLI flags: `--merge`, `--validate`, `--format`

2. **Core Data Structures**
   - `FactCollection` and `Fact` types
   - `MergePolicy` enum with all required modes
   - `FactMerger` with comprehensive merging logic
   - Storage interface and implementations

### ❌ Missing Implementation

## Implementation Requirements

### 1. Custom Fact Format Support

**Current**: Only supports flat key-value format
**Required**: Server-specific custom facts and overrides format

```json
{
  "web-001": {
    "custom": {
      "application": {
        "name": "nginx",
        "version": "1.18.0",
        "config_path": "/etc/nginx/nginx.conf"
      },
      "environment": {
        "datacenter": "fra00",
        "rack": "A01",
        "power_zone": "PZ-1"
      },
      "monitoring": {
        "prometheus_port": 9100,
        "alert_manager": "alert.example.com"
      }
    },
    "overrides": {
      "os": {
        "name": "ubuntu",
        "version": "22.04.2"
      }
    }
  },
  "web-002": {
    "custom": {
      "application": {
        "name": "apache",
        "version": "2.4.54"
      }
    }
  }
}
```

### 2. Enhanced CLI Interface

**Current Flags**:
```bash
--merge, --validate, --format, --mapping
```

**Required Additional Flags**:
```bash
--source (local, http)
--file (for local files)
--url (for HTTP endpoints)
--merge-mode (replace, merge, append, select)
--select-facts (comma-separated list)
--override
--dry-run
--server (specific server filtering)
```



### 3. Advanced Merge Logic

**Missing**: Deep merge functionality for nested structures and override handling

### 4. Fact Validation

**Missing**: Comprehensive validation for custom fact format

### 5. Template Integration

**Missing**: Make custom facts available in templates

## Implementation Plan

### Phase 1: Core Custom Fact Support (Priority: High)

#### 1.1 Data Structure Implementation

**File**: `internal/facts/types.go`

```go
// CustomFacts represents the custom fact format for a server
type CustomFacts struct {
    Custom    map[string]interface{} `json:"custom"`
    Overrides map[string]interface{} `json:"overrides"`
    Source    string                 `json:"source,omitempty"`
}

// ImportOptions defines import configuration
type ImportOptions struct {
    Source      string            `json:"source"`
    Path        string            `json:"path"`
    MergeMode   MergeMode         `json:"merge_mode"`
    SelectFacts []string          `json:"select_facts"`
    Override    bool              `json:"override"`
    Validate    bool              `json:"validate"`
    DryRun      bool              `json:"dry_run"`
    Server      string            `json:"server"`
}

// MergeMode defines merge behavior
type MergeMode string

const (
    MergeModeReplace MergeMode = "replace"
    MergeModeMerge   MergeMode = "merge"
    MergeModeAppend  MergeMode = "append"
    MergeModeSelect  MergeMode = "select"
)
```

#### 1.2 Enhanced JSON Collector

**File**: `internal/facts/json_collector.go`

Update to support server-specific format:

```go
func (c *JSONCollector) parseCustomFactsFormat(data []byte, server string) (*FactCollection, error) {
    var customFactsMap map[string]*CustomFacts
    if err := json.Unmarshal(data, &customFactsMap); err != nil {
        return nil, fmt.Errorf("failed to parse custom facts format: %w", err)
    }
    
    collection := &FactCollection{
        Server:    server,
        Timestamp: time.Now(),
        Facts:     make(map[string]*Fact),
    }
    
    // Process custom facts for the specified server
    if customFacts, exists := customFactsMap[server]; exists {
        // Add custom facts
        for category, facts := range customFacts.Custom {
            for key, value := range facts.(map[string]interface{}) {
                factKey := fmt.Sprintf("custom.%s.%s", category, key)
                collection.Facts[factKey] = &Fact{
                    Key:       factKey,
                    Value:     value,
                    Source:    string(SourceCustom),
                    Server:    server,
                    Timestamp: collection.Timestamp,
                    TTL:       DefaultTTL,
                    Metadata:  map[string]interface{}{"category": category},
                }
            }
        }
        
        // Add overrides
        for category, facts := range customFacts.Overrides {
            for key, value := range facts.(map[string]interface{}) {
                factKey := fmt.Sprintf("override.%s.%s", category, key)
                collection.Facts[factKey] = &Fact{
                    Key:       factKey,
                    Value:     value,
                    Source:    string(SourceCustom),
                    Server:    server,
                    Timestamp: collection.Timestamp,
                    TTL:       DefaultTTL,
                    Metadata:  map[string]interface{}{"category": category, "override": true},
                }
            }
        }
    }
    
    return collection, nil
}
```

#### 1.3 Deep Merge Implementation

**File**: `internal/facts/merge.go`

```go
// DeepMerge performs deep merging of nested structures
func DeepMerge(existing, custom interface{}) interface{} {
    if existing == nil {
        return custom
    }
    if custom == nil {
        return existing
    }
    
    switch existingVal := existing.(type) {
    case map[string]interface{}:
        if customMap, ok := custom.(map[string]interface{}); ok {
            merged := make(map[string]interface{})
            
            // Copy existing values
            for k, v := range existingVal {
                merged[k] = v
            }
            
            // Merge custom values
            for k, v := range customMap {
                if existingVal, exists := existingVal[k]; exists {
                    merged[k] = DeepMerge(existingVal, v)
                } else {
                    merged[k] = v
                }
            }
            
            return merged
        }
    case []interface{}:
        if customSlice, ok := custom.([]interface{}); ok {
            // For arrays, append custom values
            return append(existingVal, customSlice...)
        }
    }
    
    // For primitive types, prefer custom value
    return custom
}

// ApplyOverrides applies overrides to existing facts
func ApplyOverrides(facts *FactCollection, overrides map[string]interface{}) *FactCollection {
    if overrides == nil {
        return facts
    }
    
    merged := facts.Clone()
    
    for category, categoryOverrides := range overrides {
        if categoryMap, ok := categoryOverrides.(map[string]interface{}); ok {
            for key, value := range categoryMap {
                factKey := fmt.Sprintf("%s.%s", category, key)
                merged.Facts[factKey] = &Fact{
                    Key:       factKey,
                    Value:     value,
                    Source:    string(SourceCustom),
                    Server:    facts.Server,
                    Timestamp: time.Now(),
                    TTL:       DefaultTTL,
                    Metadata:  map[string]interface{}{"override": true, "category": category},
                }
            }
        }
    }
    
    return merged
}
```

### Phase 2: Enhanced CLI Interface (Priority: High)

#### 2.1 Updated CLI Commands

**File**: `internal/cli/facts.go`

```go
var (
    importSource      string
    importFile        string
    importURL         string
    importRepo        string
    importPath        string
    importMergeMode   string
    importSelectFacts []string
    importOverride    bool
    importValidate    bool
    importDryRun      bool
    importServer      string
)

func initFactsImportFlags() {
    factsImportCmd.Flags().StringVar(&importSource, "source", "local", "Import source: local, http")
    factsImportCmd.Flags().StringVar(&importFile, "file", "", "Path to local JSON file")
    factsImportCmd.Flags().StringVar(&importURL, "url", "", "HTTP/HTTPS URL for remote import")
    factsImportCmd.Flags().StringVar(&importMergeMode, "merge-mode", "replace", "Merge mode: replace, merge, append, select")
    factsImportCmd.Flags().StringSliceVar(&importSelectFacts, "select-facts", nil, "Comma-separated list of facts to import")
    factsImportCmd.Flags().BoolVar(&importOverride, "override", false, "Allow fact overrides")
    factsImportCmd.Flags().BoolVar(&importValidate, "validate", false, "Validate facts before importing")
    factsImportCmd.Flags().BoolVar(&importDryRun, "dry-run", false, "Show what would be imported without importing")
    factsImportCmd.Flags().StringVar(&importServer, "server", "", "Specific server to import facts for")
}
```

#### 2.2 Enhanced Import Function

```go
func runFactsImport(_ *cobra.Command, args []string) error {
    options := &facts.ImportOptions{
        Source:      importSource,
        Path:        importPath,
        MergeMode:   facts.MergeMode(importMergeMode),
        SelectFacts: importSelectFacts,
        Override:    importOverride,
        Validate:    importValidate,
        DryRun:      importDryRun,
        Server:      importServer,
    }
    
    // Determine source path
    var sourcePath string
    switch importSource {
    case "local":
        sourcePath = importFile
    case "http":
        sourcePath = importURL
    default:
        return fmt.Errorf("unsupported source: %s", importSource)
    }
    
    // Create storage and manager
    storage, err := facts.NewFactStorage(facts.StorageOptions{
        Type: facts.StorageTypeBadger,
        Path: getFactsDBPath(),
    })
    if err != nil {
        return fmt.Errorf("failed to create storage: %w", err)
    }
    defer storage.Close()
    
    manager := facts.NewManagerWithStorage(nil, storage)
    
    // Import facts
    if err := manager.ImportCustomFactsWithOptions(sourcePath, options); err != nil {
        return fmt.Errorf("failed to import facts: %w", err)
    }
    
    if importDryRun {
        fmt.Println("DRY RUN: Facts would be imported successfully")
    } else {
        fmt.Println("Facts imported successfully")
    }
    
    return nil
}
```



### Phase 3: Fact Validation (Priority: Medium)

#### 3.1 Validation Implementation

**File**: `internal/facts/validation.go`

```go
package facts

import (
    "encoding/json"
    "fmt"
    "regexp"
    "strings"
)

// ValidationError represents a validation error
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
    Value   interface{} `json:"value,omitempty"`
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error in %s: %s", e.Field, e.Message)
}

// ValidationResult contains validation results
type ValidationResult struct {
    Valid   bool              `json:"valid"`
    Errors  []*ValidationError `json:"errors,omitempty"`
    Warnings []*ValidationError `json:"warnings,omitempty"`
}

// ValidateCustomFacts validates custom facts format
func ValidateCustomFacts(facts map[string]*CustomFacts) *ValidationResult {
    result := &ValidationResult{Valid: true}
    
    for serverID, customFacts := range facts {
        // Validate server ID
        if err := validateServerID(serverID); err != nil {
            result.Valid = false
            result.Errors = append(result.Errors, &ValidationError{
                Field:   "server_id",
                Message: err.Error(),
                Value:   serverID,
            })
        }
        
        // Validate custom facts
        if customFacts.Custom != nil {
            if err := validateCustomFactStructure(customFacts.Custom); err != nil {
                result.Valid = false
                result.Errors = append(result.Errors, &ValidationError{
                    Field:   "custom",
                    Message: err.Error(),
                    Value:   customFacts.Custom,
                })
            }
        }
        
        // Validate overrides
        if customFacts.Overrides != nil {
            if err := validateOverrideStructure(customFacts.Overrides); err != nil {
                result.Valid = false
                result.Errors = append(result.Errors, &ValidationError{
                    Field:   "overrides",
                    Message: err.Error(),
                    Value:   customFacts.Overrides,
                })
            }
        }
    }
    
    return result
}

// validateServerID validates server identifier
func validateServerID(serverID string) error {
    if serverID == "" {
        return fmt.Errorf("server ID cannot be empty")
    }
    
    // Check for valid characters
    validPattern := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
    if !validPattern.MatchString(serverID) {
        return fmt.Errorf("server ID contains invalid characters")
    }
    
    return nil
}

// validateCustomFactStructure validates custom fact structure
func validateCustomFactStructure(custom map[string]interface{}) error {
    for category, facts := range custom {
        if category == "" {
            return fmt.Errorf("category name cannot be empty")
        }
        
        if factsMap, ok := facts.(map[string]interface{}); ok {
            for key, value := range factsMap {
                if key == "" {
                    return fmt.Errorf("fact key cannot be empty in category %s", category)
                }
                
                if value == nil {
                    return fmt.Errorf("fact value cannot be nil for %s.%s", category, key)
                }
            }
        } else {
            return fmt.Errorf("category %s must be an object", category)
        }
    }
    
    return nil
}

// validateOverrideStructure validates override structure
func validateOverrideStructure(overrides map[string]interface{}) error {
    for category, facts := range overrides {
        if category == "" {
            return fmt.Errorf("override category cannot be empty")
        }
        
        if factsMap, ok := facts.(map[string]interface{}); ok {
            for key, value := range factsMap {
                if key == "" {
                    return fmt.Errorf("override key cannot be empty in category %s", category)
                }
                
                if value == nil {
                    return fmt.Errorf("override value cannot be nil for %s.%s", category, key)
                }
            }
        } else {
            return fmt.Errorf("override category %s must be an object", category)
        }
    }
    
    return nil
}
```

### Phase 4: Template Integration (Priority: Low)

#### 4.1 Template Support

**File**: `internal/facts/manager.go`

```go
// GetCustomFacts retrieves custom facts for template usage
func (m *Manager) GetCustomFacts(server string) (map[string]interface{}, error) {
    if m.storage == nil {
        return nil, fmt.Errorf("no storage configured")
    }
    
    // Load persisted facts
    collection, err := m.LoadPersistedFacts(server)
    if err != nil {
        return nil, err
    }
    
    // Extract custom facts
    customFacts := make(map[string]interface{})
    for key, fact := range collection.Facts {
        if strings.HasPrefix(key, "custom.") {
            parts := strings.Split(key, ".")
            if len(parts) >= 3 {
                category := parts[1]
                factKey := parts[2]
                
                if customFacts[category] == nil {
                    customFacts[category] = make(map[string]interface{})
                }
                
                if categoryMap, ok := customFacts[category].(map[string]interface{}); ok {
                    categoryMap[factKey] = fact.Value
                }
            }
        }
    }
    
    return customFacts, nil
}
```

#### 4.2 Template Engine Update

**File**: `internal/cli/templates.go`

```go
// Update template context to include custom facts
func buildTemplateContext(server string, manager *facts.Manager) (map[string]interface{}, error) {
    context := make(map[string]interface{})
    
    // Add system facts
    if collection, err := manager.CollectAllFacts(server); err == nil {
        for key, fact := range collection.Facts {
            if !strings.HasPrefix(key, "custom.") && !strings.HasPrefix(key, "override.") {
                context[key] = fact.Value
            }
        }
    }
    
    // Add custom facts
    if customFacts, err := manager.GetCustomFacts(server); err == nil {
        context["custom"] = customFacts
    }
    
    return context, nil
}
```

### Phase 5: Testing Strategy (Priority: High)

#### 5.1 Unit Tests

**File**: `internal/facts/import_test.go`

```go
package facts

import (
    "encoding/json"
    "testing"
    "time"
)

func TestCustomFactsParsing(t *testing.T) {
    testData := map[string]*CustomFacts{
        "web-001": {
            Custom: map[string]interface{}{
                "application": map[string]interface{}{
                    "name":    "nginx",
                    "version": "1.18.0",
                },
            },
            Overrides: map[string]interface{}{
                "os": map[string]interface{}{
                    "name": "ubuntu",
                },
            },
        },
    }
    
    data, err := json.Marshal(testData)
    if err != nil {
        t.Fatalf("Failed to marshal test data: %v", err)
    }
    
    collector := NewJSONCollector("", MergePolicyReplace)
    collection, err := collector.parseCustomFactsFormat(data, "web-001")
    if err != nil {
        t.Fatalf("Failed to parse custom facts: %v", err)
    }
    
    // Verify custom facts
    if fact, exists := collection.Facts["custom.application.name"]; !exists {
        t.Error("custom.application.name fact not found")
    } else if fact.Value != "nginx" {
        t.Errorf("Expected nginx, got %v", fact.Value)
    }
    
    // Verify overrides
    if fact, exists := collection.Facts["override.os.name"]; !exists {
        t.Error("override.os.name fact not found")
    } else if fact.Value != "ubuntu" {
        t.Errorf("Expected ubuntu, got %v", fact.Value)
    }
}

func TestDeepMerge(t *testing.T) {
    existing := map[string]interface{}{
        "app": map[string]interface{}{
            "name": "old-app",
            "config": map[string]interface{}{
                "port": 8080,
            },
        },
    }
    
    custom := map[string]interface{}{
        "app": map[string]interface{}{
            "name": "new-app",
            "config": map[string]interface{}{
                "host": "localhost",
            },
        },
    }
    
    merged := DeepMerge(existing, custom).(map[string]interface{})
    
    app := merged["app"].(map[string]interface{})
    if app["name"] != "new-app" {
        t.Errorf("Expected new-app, got %v", app["name"])
    }
    
    config := app["config"].(map[string]interface{})
    if config["port"] != 8080 {
        t.Errorf("Expected port 8080, got %v", config["port"])
    }
    if config["host"] != "localhost" {
        t.Errorf("Expected host localhost, got %v", config["host"])
    }
}
```

#### 5.2 Integration Tests

**File**: `tests/integration/custom_facts_test.go`

```go
package integration

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "os"
    "path/filepath"
    "testing"
    
    "spooky/internal/facts"
)

func TestCustomFactsImportIntegration(t *testing.T) {
    // Create test facts file
    testFacts := map[string]*facts.CustomFacts{
        "test-server": {
            Custom: map[string]interface{}{
                "application": map[string]interface{}{
                    "name": "test-app",
                },
            },
        },
    }
    
    tempDir := t.TempDir()
    factsFile := filepath.Join(tempDir, "facts.json")
    
    data, err := json.MarshalIndent(testFacts, "", "  ")
    if err != nil {
        t.Fatalf("Failed to marshal test facts: %v", err)
    }
    
    if err := os.WriteFile(factsFile, data, 0644); err != nil {
        t.Fatalf("Failed to write test file: %v", err)
    }
    
    // Test import
    storage, err := facts.NewFactStorage(facts.StorageOptions{
        Type: facts.StorageTypeBadger,
        Path: filepath.Join(tempDir, "facts.db"),
    })
    if err != nil {
        t.Fatalf("Failed to create storage: %v", err)
    }
    defer storage.Close()
    
    manager := facts.NewManagerWithStorage(nil, storage)
    
    collection, err := manager.ImportCustomFacts(factsFile, "test-server", facts.MergePolicyReplace)
    if err != nil {
        t.Fatalf("Failed to import facts: %v", err)
    }
    
    if len(collection.Facts) == 0 {
        t.Error("No facts imported")
    }
    
    if fact, exists := collection.Facts["custom.application.name"]; !exists {
        t.Error("custom.application.name fact not found")
    } else if fact.Value != "test-app" {
        t.Errorf("Expected test-app, got %v", fact.Value)
    }
}

func TestHTTPCustomFactsIntegration(t *testing.T) {
    // Create test HTTP server
    testFacts := map[string]*facts.CustomFacts{
        "http-server": {
            Custom: map[string]interface{}{
                "monitoring": map[string]interface{}{
                    "port": 9100,
                },
            },
        },
    }
    
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(testFacts)
    }))
    defer server.Close()
    
    // Test HTTP import
    collector := facts.NewHTTPCollector(server.URL, nil, 30*time.Second, facts.MergePolicyReplace)
    collection, err := collector.Collect("http-server")
    if err != nil {
        t.Fatalf("Failed to collect facts from HTTP: %v", err)
    }
    
    if len(collection.Facts) == 0 {
        t.Error("No facts collected from HTTP")
    }
}
```

### Phase 6: Documentation and Examples (Priority: Low)

#### 6.1 CLI Documentation Update

**File**: `docs/cli-specification.md`

Update the facts import section with new flags and examples.

#### 6.2 Usage Examples

**File**: `examples/custom-facts/`

Create example custom fact files and usage documentation.

## Success Criteria

1. **Local imports**: Successfully import from local JSON files with custom format
2. **HTTP imports**: Successfully import from HTTP/HTTPS endpoints
3. **Merging**: Correctly merge custom facts with system facts using all merge modes
4. **Validation**: Validate fact format and content
5. **CLI integration**: Easy import commands with all required flags
6. **Template integration**: Custom facts available in templates
7. **Performance**: Efficient import for large fact sets
8. **Error handling**: Graceful handling of import errors

## Dependencies

- **HTTP client**: Standard library (already implemented)
- **JSON parsing**: Standard library (already implemented)

## File Structure

```
spooky/
├── internal/
│ ├── facts/
│ │ ├── types.go              # Updated with CustomFacts and ImportOptions
│ │ ├── json_collector.go     # Updated for custom format
│ │ ├── merge.go              # Updated with deep merge
│ │ ├── validation.go         # New validation logic
│ │ └── manager.go            # Updated with enhanced import
│ └── cli/
│   └── facts.go              # Updated CLI commands
├── tests/
│ └── integration/
│   └── custom_facts_test.go  # Integration tests
└── examples/
  └── custom-facts/           # Example files
```

## Timeline

- **Phase 1**: 1-2 weeks (Core custom fact support)
- **Phase 2**: 1 week (Enhanced CLI interface)
- **Phase 3**: 1 week (Fact validation)
- **Phase 4**: 1 week (Template integration)
- **Phase 5**: 1-2 weeks (Testing)
- **Phase 6**: 1 week (Documentation)

**Total Estimated Time**: 6-8 weeks

## Risk Assessment

### High Risk
- Deep merge logic for complex nested structures
- Performance with large fact sets

### Medium Risk
- Template integration with existing systems
- Validation logic completeness
- CLI flag compatibility

### Low Risk
- Documentation updates
- Example creation
- Basic JSON parsing

## Out of Scope

- **Real-time sync**: No continuous import from remote sources
- **Complex transformations**: No data transformation during import
- **External APIs**: No integration with external fact APIs
- **Authentication**: Basic auth only for remote sources
- **AWS/S3**: Explicitly excluded from this implementation 