# Implementation Plan: Machine Filtering Implementation

## Overview
Implement proper machine filtering for import/export operations, replacing placeholder implementations with real filtering capabilities based on tags, environment, and custom criteria.

## Task Details
- **Task ID**: 7.5
- **Priority**: Medium
- **Files**: 
  - `internal/machines/importexport/manager.go`
- **Functions**: Machine filtering, import/export with filters

## Current State Analysis

### Existing Patterns
1. **Machine Types**: Machine configurations defined in `internal/machines/types/`
2. **Import/Export**: Basic import/export functionality exists
3. **Tag System**: Machine tagging system implemented
4. **Error Handling**: Consistent error wrapping

### Current Placeholder Code
```go
// internal/machines/importexport/manager.go line 45
// This is a placeholder implementation
```

## Implementation Requirements

### Interface Compliance
The machine filtering system must:
1. **Filter by tags** and tag combinations
2. **Filter by environment** and environment variables
3. **Filter by custom criteria** and expressions
4. **Support complex queries** with AND/OR logic
5. **Provide performance optimization** for large inventories
6. **Support export filtering** for selective exports
7. **Handle import filtering** for selective imports

### Required Dependencies
- Machine inventory system
- Tag management system
- Query parsing system
- Performance optimization system

## Detailed Implementation Plan

### Step 1: Implement Filter Criteria System

#### 1.1 Filter Criteria Definition
```go
// internal/machines/importexport/filter.go
package importexport

import (
    "fmt"
    "regexp"
    "strings"
    "time"
)

// FilterCriteria represents machine filtering criteria
type FilterCriteria struct {
    Tags        []string            `json:"tags,omitempty"`
    Environments []string           `json:"environments,omitempty"`
    HostPattern  string             `json:"host_pattern,omitempty"`
    NamePattern  string             `json:"name_pattern,omitempty"`
    UserPattern  string             `json:"user_pattern,omitempty"`
    PortRange    *PortRange         `json:"port_range,omitempty"`
    CustomRules  []CustomRule       `json:"custom_rules,omitempty"`
    Logic        FilterLogic        `json:"logic"`
    Limit        int                `json:"limit,omitempty"`
    Offset       int                `json:"offset,omitempty"`
}

// PortRange represents a port range filter
type PortRange struct {
    Min int `json:"min"`
    Max int `json:"max"`
}

// CustomRule represents a custom filtering rule
type CustomRule struct {
    Field    string      `json:"field"`
    Operator string      `json:"operator"`
    Value    interface{} `json:"value"`
}

// FilterLogic represents the logical operation for combining filters
type FilterLogic string

const (
    FilterLogicAND FilterLogic = "AND"
    FilterLogicOR  FilterLogic = "OR"
)

// NewFilterCriteria creates new filter criteria
func NewFilterCriteria() *FilterCriteria {
    return &FilterCriteria{
        Logic: FilterLogicAND,
    }
}

// AddTag adds a tag filter
func (fc *FilterCriteria) AddTag(tag string) *FilterCriteria {
    fc.Tags = append(fc.Tags, tag)
    return fc
}

// AddEnvironment adds an environment filter
func (fc *FilterCriteria) AddEnvironment(env string) *FilterCriteria {
    fc.Environments = append(fc.Environments, env)
    return fc
}

// SetHostPattern sets host pattern filter
func (fc *FilterCriteria) SetHostPattern(pattern string) *FilterCriteria {
    fc.HostPattern = pattern
    return fc
}

// SetNamePattern sets name pattern filter
func (fc *FilterCriteria) SetNamePattern(pattern string) *FilterCriteria {
    fc.NamePattern = pattern
    return fc
}

// SetPortRange sets port range filter
func (fc *FilterCriteria) SetPortRange(min, max int) *FilterCriteria {
    fc.PortRange = &PortRange{Min: min, Max: max}
    return fc
}

// AddCustomRule adds a custom filtering rule
func (fc *FilterCriteria) AddCustomRule(field, operator string, value interface{}) *FilterCriteria {
    fc.CustomRules = append(fc.CustomRules, CustomRule{
        Field:    field,
        Operator: operator,
        Value:    value,
    })
    return fc
}

// SetLogic sets the logical operation
func (fc *FilterCriteria) SetLogic(logic FilterLogic) *FilterCriteria {
    fc.Logic = logic
    return fc
}

// SetLimit sets the result limit
func (fc *FilterCriteria) SetLimit(limit int) *FilterCriteria {
    fc.Limit = limit
    return fc
}

// SetOffset sets the result offset
func (fc *FilterCriteria) SetOffset(offset int) *FilterCriteria {
    fc.Offset = offset
    return fc
}
```

#### 1.2 Filter Engine Implementation
```go
// internal/machines/importexport/engine.go
package importexport

import (
    "fmt"
    "regexp"
    "strings"
    
    "spooky/internal/machines/types"
)

// FilterEngine applies filtering criteria to machines
type FilterEngine struct{}

// NewFilterEngine creates a new filter engine
func NewFilterEngine() *FilterEngine {
    return &FilterEngine{}
}

// FilterMachines filters machines based on criteria
func (e *FilterEngine) FilterMachines(machines []*types.Machine, criteria *FilterCriteria) ([]*types.Machine, error) {
    if criteria == nil {
        return machines, nil
    }

    var filtered []*types.Machine

    for _, machine := range machines {
        matches, err := e.machineMatches(machine, criteria)
        if err != nil {
            return nil, fmt.Errorf("filter evaluation failed: %w", err)
        }

        if matches {
            filtered = append(filtered, machine)
        }
    }

    // Apply limit and offset
    filtered = e.applyLimitOffset(filtered, criteria)

    return filtered, nil
}

// machineMatches checks if a machine matches the filter criteria
func (e *FilterEngine) machineMatches(machine *types.Machine, criteria *FilterCriteria) (bool, error) {
    var results []bool

    // Check tag filters
    if len(criteria.Tags) > 0 {
        matches := e.matchesTags(machine, criteria.Tags)
        results = append(results, matches)
    }

    // Check environment filters
    if len(criteria.Environments) > 0 {
        matches := e.matchesEnvironments(machine, criteria.Environments)
        results = append(results, matches)
    }

    // Check host pattern
    if criteria.HostPattern != "" {
        matches, err := e.matchesPattern(machine.Host, criteria.HostPattern)
        if err != nil {
            return false, fmt.Errorf("host pattern matching failed: %w", err)
        }
        results = append(results, matches)
    }

    // Check name pattern
    if criteria.NamePattern != "" {
        matches, err := e.matchesPattern(machine.Name, criteria.NamePattern)
        if err != nil {
            return false, fmt.Errorf("name pattern matching failed: %w", err)
        }
        results = append(results, matches)
    }

    // Check port range
    if criteria.PortRange != nil {
        matches := e.matchesPortRange(machine.Port, criteria.PortRange)
        results = append(results, matches)
    }

    // Check custom rules
    for _, rule := range criteria.CustomRules {
        matches, err := e.matchesCustomRule(machine, rule)
        if err != nil {
            return false, fmt.Errorf("custom rule evaluation failed: %w", err)
        }
        results = append(results, matches)
    }

    // Apply logical operation
    return e.applyLogic(results, criteria.Logic), nil
}

// matchesTags checks if machine matches tag criteria
func (e *FilterEngine) matchesTags(machine *types.Machine, tags []string) bool {
    machineTags := make(map[string]bool)
    for _, tag := range machine.Tags {
        machineTags[tag] = true
    }

    for _, requiredTag := range tags {
        if !machineTags[requiredTag] {
            return false
        }
    }
    return true
}

// matchesEnvironments checks if machine matches environment criteria
func (e *FilterEngine) matchesEnvironments(machine *types.Machine, environments []string) bool {
    machineEnv := make(map[string]bool)
    for _, env := range machine.Environment {
        machineEnv[env] = true
    }

    for _, requiredEnv := range environments {
        if !machineEnv[requiredEnv] {
            return false
        }
    }
    return true
}

// matchesPattern checks if value matches pattern
func (e *FilterEngine) matchesPattern(value, pattern string) (bool, error) {
    // Handle wildcard patterns
    if strings.Contains(pattern, "*") {
        regexPattern := strings.ReplaceAll(pattern, "*", ".*")
        regex, err := regexp.Compile(regexPattern)
        if err != nil {
            return false, fmt.Errorf("invalid pattern: %s", pattern)
        }
        return regex.MatchString(value), nil
    }

    // Handle regex patterns
    if strings.HasPrefix(pattern, "/") && strings.HasSuffix(pattern, "/") {
        regexPattern := pattern[1 : len(pattern)-1]
        regex, err := regexp.Compile(regexPattern)
        if err != nil {
            return false, fmt.Errorf("invalid regex pattern: %s", regexPattern)
        }
        return regex.MatchString(value), nil
    }

    // Exact match
    return value == pattern, nil
}

// matchesPortRange checks if port is in range
func (e *FilterEngine) matchesPortRange(port int, portRange *PortRange) bool {
    return port >= portRange.Min && port <= portRange.Max
}

// matchesCustomRule checks if machine matches custom rule
func (e *FilterEngine) matchesCustomRule(machine *types.Machine, rule CustomRule) (bool, error) {
    // Get field value
    var fieldValue interface{}
    switch rule.Field {
    case "name":
        fieldValue = machine.Name
    case "host":
        fieldValue = machine.Host
    case "port":
        fieldValue = machine.Port
    case "user":
        fieldValue = machine.User
    default:
        return false, fmt.Errorf("unknown field: %s", rule.Field)
    }

    // Apply operator
    return e.applyOperator(fieldValue, rule.Operator, rule.Value)
}

// applyOperator applies comparison operator
func (e *FilterEngine) applyOperator(fieldValue interface{}, operator string, ruleValue interface{}) (bool, error) {
    switch operator {
    case "eq", "==":
        return fieldValue == ruleValue, nil
    case "ne", "!=":
        return fieldValue != ruleValue, nil
    case "gt", ">":
        return e.compareValues(fieldValue, ruleValue) > 0, nil
    case "gte", ">=":
        return e.compareValues(fieldValue, ruleValue) >= 0, nil
    case "lt", "<":
        return e.compareValues(fieldValue, ruleValue) < 0, nil
    case "lte", "<=":
        return e.compareValues(fieldValue, ruleValue) <= 0, nil
    case "contains":
        return e.containsValue(fieldValue, ruleValue), nil
    case "in":
        return e.inValue(fieldValue, ruleValue), nil
    default:
        return false, fmt.Errorf("unknown operator: %s", operator)
    }
}

// compareValues compares two values
func (e *FilterEngine) compareValues(a, b interface{}) int {
    // Implementation would handle type conversion and comparison
    // For now, return 0 (equal)
    return 0
}

// containsValue checks if value contains another value
func (e *FilterEngine) containsValue(fieldValue, ruleValue interface{}) bool {
    // Implementation would handle string containment
    return false
}

// inValue checks if value is in a list
func (e *FilterEngine) inValue(fieldValue, ruleValue interface{}) bool {
    // Implementation would handle list membership
    return false
}

// applyLogic applies logical operation to results
func (e *FilterEngine) applyLogic(results []bool, logic FilterLogic) bool {
    if len(results) == 0 {
        return true
    }

    switch logic {
    case FilterLogicAND:
        for _, result := range results {
            if !result {
                return false
            }
        }
        return true
    case FilterLogicOR:
        for _, result := range results {
            if result {
                return true
            }
        }
        return false
    default:
        return false
    }
}

// applyLimitOffset applies limit and offset to results
func (e *FilterEngine) applyLimitOffset(machines []*types.Machine, criteria *FilterCriteria) []*types.Machine {
    if criteria.Offset >= len(machines) {
        return []*types.Machine{}
    }

    start := criteria.Offset
    end := len(machines)

    if criteria.Limit > 0 && start+criteria.Limit < end {
        end = start + criteria.Limit
    }

    return machines[start:end]
}
```

### Step 2: Implement Import/Export Manager

#### 2.1 Enhanced Import/Export Manager
```go
// internal/machines/importexport/manager.go
package importexport

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    
    "spooky/internal/machines/types"
    "spooky/internal/logging"
)

// Manager manages machine import/export operations
type Manager struct {
    filterEngine *FilterEngine
    logger       logging.Logger
}

// NewManager creates a new import/export manager
func NewManager(logger logging.Logger) *Manager {
    return &Manager{
        filterEngine: NewFilterEngine(),
        logger:       logger,
    }
}

// ExportMachines exports machines with filtering
func (m *Manager) ExportMachines(ctx context.Context, machines []*types.Machine, criteria *FilterCriteria, outputPath string) error {
    m.logger.Info("Exporting machines",
        logging.Int("total_machines", len(machines)),
        logging.String("output_path", outputPath))

    // Apply filters
    filteredMachines, err := m.filterEngine.FilterMachines(machines, criteria)
    if err != nil {
        return fmt.Errorf("failed to filter machines: %w", err)
    }

    m.logger.Info("Machines filtered for export",
        logging.Int("filtered_machines", len(filteredMachines)))

    // Create export data
    exportData := &ExportData{
        Machines: filteredMachines,
        Metadata: ExportMetadata{
            TotalMachines:    len(machines),
            FilteredMachines: len(filteredMachines),
            ExportTime:       time.Now(),
            FilterCriteria:   criteria,
        },
    }

    // Ensure output directory exists
    outputDir := filepath.Dir(outputPath)
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        return fmt.Errorf("failed to create output directory: %w", err)
    }

    // Write to file
    file, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")

    if err := encoder.Encode(exportData); err != nil {
        return fmt.Errorf("failed to encode export data: %w", err)
    }

    m.logger.Info("Machines exported successfully",
        logging.String("output_path", outputPath),
        logging.Int("exported_machines", len(filteredMachines)))

    return nil
}

// ImportMachines imports machines with filtering
func (m *Manager) ImportMachines(ctx context.Context, inputPath string, criteria *FilterCriteria) ([]*types.Machine, error) {
    m.logger.Info("Importing machines",
        logging.String("input_path", inputPath))

    // Read input file
    file, err := os.Open(inputPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open input file: %w", err)
    }
    defer file.Close()

    // Decode export data
    var exportData ExportData
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&exportData); err != nil {
        return nil, fmt.Errorf("failed to decode export data: %w", err)
    }

    m.logger.Info("Import data loaded",
        logging.Int("total_machines", exportData.Metadata.TotalMachines),
        logging.Int("filtered_machines", exportData.Metadata.FilteredMachines))

    // Apply filters if specified
    if criteria != nil {
        filteredMachines, err := m.filterEngine.FilterMachines(exportData.Machines, criteria)
        if err != nil {
            return nil, fmt.Errorf("failed to filter imported machines: %w", err)
        }

        m.logger.Info("Imported machines filtered",
            logging.Int("imported_machines", len(filteredMachines)))

        return filteredMachines, nil
    }

    return exportData.Machines, nil
}

// ListMachines lists machines with filtering
func (m *Manager) ListMachines(ctx context.Context, machines []*types.Machine, criteria *FilterCriteria) ([]*types.Machine, error) {
    m.logger.Debug("Listing machines with filters",
        logging.Int("total_machines", len(machines)))

    // Apply filters
    filteredMachines, err := m.filterEngine.FilterMachines(machines, criteria)
    if err != nil {
        return nil, fmt.Errorf("failed to filter machines: %w", err)
    }

    m.logger.Info("Machines listed",
        logging.Int("filtered_machines", len(filteredMachines)))

    return filteredMachines, nil
}

// GetMachineCount gets count of machines matching criteria
func (m *Manager) GetMachineCount(ctx context.Context, machines []*types.Machine, criteria *FilterCriteria) (int, error) {
    m.logger.Debug("Getting machine count with filters",
        logging.Int("total_machines", len(machines)))

    // Apply filters
    filteredMachines, err := m.filterEngine.FilterMachines(machines, criteria)
    if err != nil {
        return 0, fmt.Errorf("failed to filter machines: %w", err)
    }

    count := len(filteredMachines)

    m.logger.Info("Machine count retrieved",
        logging.Int("filtered_count", count))

    return count, nil
}
```

### Step 3: Implement Export Data Structures

#### 3.1 Export Data Definition
```go
// internal/machines/importexport/data.go
package importexport

import (
    "time"
    
    "spooky/internal/machines/types"
)

// ExportData represents exported machine data
type ExportData struct {
    Machines []*types.Machine `json:"machines"`
    Metadata ExportMetadata   `json:"metadata"`
}

// ExportMetadata represents export metadata
type ExportMetadata struct {
    TotalMachines    int            `json:"total_machines"`
    FilteredMachines int            `json:"filtered_machines"`
    ExportTime       time.Time      `json:"export_time"`
    FilterCriteria   *FilterCriteria `json:"filter_criteria,omitempty"`
    Version          string         `json:"version"`
    Source           string         `json:"source"`
}

// NewExportData creates new export data
func NewExportData(machines []*types.Machine, criteria *FilterCriteria) *ExportData {
    return &ExportData{
        Machines: machines,
        Metadata: ExportMetadata{
            TotalMachines:    len(machines),
            FilteredMachines: len(machines),
            ExportTime:       time.Now(),
            FilterCriteria:   criteria,
            Version:          "1.0.0",
            Source:           "spooky",
        },
    }
}
```

## Configuration Options

### Supported Options
- **FilteringEnabled**: Enable/disable filtering
- **PerformanceOptimization**: Enable/disable performance optimization
- **ExportFormat**: JSON, YAML, CSV export formats
- **ImportValidation**: Enable/disable import validation

## Dependencies

### Internal Dependencies
- `spooky/internal/machines/types`
- `spooky/internal/logging`

### External Dependencies
- `context` (standard library)
- `encoding/json` (standard library)
- `fmt` (standard library)
- `os` (standard library)
- `path/filepath` (standard library)
- `regexp` (standard library)
- `strings` (standard library)
- `time` (standard library)

## Implementation Order

1. Implement filter criteria system
2. Add filter engine implementation
3. Create enhanced import/export manager
4. Implement export data structures
5. Add performance optimization
6. Add comprehensive tests
7. Documentation and cleanup
