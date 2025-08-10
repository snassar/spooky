# Implementation Plan: Merging Functions

## Overview
Implement action merging logic that can handle conflicts between actions, support different merge policies, and provide intelligent conflict resolution strategies.

## Task Details
- **Task ID**: 3.6
- **Priority**: Low
- **File**: `internal/actions/manager.go`
- **Functions**: `MergeActions`, `MergeWithPolicy`

## Current State Analysis

### Existing Patterns
1. **Action Management**: Existing action manager provides basic action operations
2. **Conflict Resolution**: Basic conflict detection exists in validation system
3. **Merge Policies**: Some merge policy definitions exist
4. **Action Types**: Action structures defined with mergeable fields
5. **Validation System**: Action validation framework exists

### Existing Implementation Examples
- **Action Manager**: `internal/actions/manager.go` provides action management
- **Action Types**: `internal/actions/types/action.go` defines action structures
- **Validation**: `internal/actions/validation/manager.go` provides validation
- **Merging**: `internal/actions/merging/manager.go` provides basic merging framework

## Implementation Requirements

### Interface Compliance
The merging functions must:
1. **Merge multiple actions** into a single coherent action
2. **Handle conflicts** between actions gracefully
3. **Support different merge policies** (replace, append, merge, skip)
4. **Provide conflict resolution strategies** for different field types
5. **Validate merged actions** for correctness
6. **Generate merge reports** with conflict details
7. **Support custom merge rules** for specific action types
8. **Handle nested structures** and complex dependencies

### Required Dependencies
- Action manager for action operations
- Validation system for merged action validation
- Merge policy definitions for conflict resolution
- Conflict detection system for identifying issues

## Detailed Implementation Plan

### Step 1: Enhance Action Manager with Merging Functions

**File**: `internal/actions/manager.go`

```go
// MergePolicy represents the policy for merging actions
type MergePolicy string

const (
    MergePolicyReplace MergePolicy = "replace"  // Replace with newer action
    MergePolicyAppend  MergePolicy = "append"   // Append to existing action
    MergePolicyMerge   MergePolicy = "merge"    // Merge fields intelligently
    MergePolicySkip    MergePolicy = "skip"     // Skip conflicting actions
    MergePolicyCustom  MergePolicy = "custom"   // Use custom merge rules
)

// MergeStrategy represents the strategy for handling conflicts
type MergeStrategy string

const (
    MergeStrategyFirstWins  MergeStrategy = "first_wins"   // First action wins
    MergeStrategyLastWins   MergeStrategy = "last_wins"    // Last action wins
    MergeStrategyCombine    MergeStrategy = "combine"      // Combine values
    MergeStrategyPrompt     MergeStrategy = "prompt"       // Prompt user
    MergeStrategyFail       MergeStrategy = "fail"         // Fail on conflict
)

// MergeConflict represents a conflict during merging
type MergeConflict struct {
    Field           string                    `json:"field"`
    Action1         *types.Action             `json:"action1"`
    Action2         *types.Action             `json:"action2"`
    Value1          interface{}               `json:"value1"`
    Value2          interface{}               `json:"value2"`
    ConflictType    ConflictType              `json:"conflict_type"`
    Severity        ConflictSeverity          `json:"severity"`
    Resolution      ConflictResolution        `json:"resolution"`
    Message         string                    `json:"message"`
    Suggestions     []string                  `json:"suggestions"`
}

// ConflictType represents the type of conflict
type ConflictType string

const (
    ConflictTypeFieldValue    ConflictType = "field_value"     // Different field values
    ConflictTypeFieldMissing  ConflictType = "field_missing"   // Field missing in one action
    ConflictTypeTypeMismatch  ConflictType = "type_mismatch"   // Different field types
    ConflictTypeDependency    ConflictType = "dependency"      // Dependency conflict
    ConflictTypeResource      ConflictType = "resource"        // Resource conflict
    ConflictTypeSecurity      ConflictType = "security"        // Security conflict
)

// ConflictSeverity represents the severity of a conflict
type ConflictSeverity string

const (
    ConflictSeverityLow      ConflictSeverity = "low"
    ConflictSeverityMedium   ConflictSeverity = "medium"
    ConflictSeverityHigh     ConflictSeverity = "high"
    ConflictSeverityCritical ConflictSeverity = "critical"
)

// ConflictResolution represents how a conflict was resolved
type ConflictResolution struct {
    Method          string                    `json:"method"`
    Value           interface{}               `json:"value"`
    Applied         bool                      `json:"applied"`
    Reason          string                    `json:"reason"`
    Timestamp       time.Time                 `json:"timestamp"`
}

// MergeResult represents the result of a merge operation
type MergeResult struct {
    MergedAction    *types.Action             `json:"merged_action"`
    Conflicts       []*MergeConflict          `json:"conflicts"`
    Resolved        int                       `json:"resolved"`
    Unresolved      int                       `json:"unresolved"`
    Skipped         int                       `json:"skipped"`
    Warnings        []string                  `json:"warnings"`
    Errors          []string                  `json:"errors"`
    Duration        time.Duration             `json:"duration"`
    Metadata        map[string]interface{}    `json:"metadata"`
}

// MergeOptions represents options for merging actions
type MergeOptions struct {
    Policy          MergePolicy               `json:"policy"`
    Strategy        MergeStrategy             `json:"strategy"`
    ResolveConflicts bool                     `json:"resolve_conflicts"`
    ValidateResult  bool                      `json:"validate_result"`
    GenerateReport  bool                      `json:"generate_report"`
    CustomRules     map[string]MergeRule      `json:"custom_rules"`
    ConflictLimit   int                       `json:"conflict_limit"`
    Timeout         time.Duration             `json:"timeout"`
}

// MergeRule represents a custom merge rule
type MergeRule struct {
    Field           string                    `json:"field"`
    Condition       string                    `json:"condition"`
    Action          string                    `json:"action"`
    Priority        int                       `json:"priority"`
    Description     string                    `json:"description"`
}

// MergeActions merges multiple actions into a single action
func (m *Manager) MergeActions(ctx context.Context, actions []*types.Action, options *MergeOptions) (*MergeResult, error) {
    startTime := time.Now()
    m.logger.Info("Merging actions",
        logging.Int("action_count", len(actions)),
        logging.String("policy", string(options.Policy)),
        logging.String("strategy", string(options.Strategy)))

    if len(actions) == 0 {
        return nil, fmt.Errorf("no actions to merge")
    }

    if len(actions) == 1 {
        return &MergeResult{
            MergedAction: actions[0],
            Conflicts:    make([]*MergeConflict, 0),
            Resolved:     0,
            Unresolved:   0,
            Skipped:      0,
            Warnings:     make([]string, 0),
            Errors:       make([]string, 0),
            Duration:     time.Since(startTime),
            Metadata:     make(map[string]interface{}),
        }, nil
    }

    // Sort actions by priority and timestamp
    sortedActions := m.sortActionsForMerge(actions)

    // Initialize merge result
    result := &MergeResult{
        MergedAction: &types.Action{},
        Conflicts:    make([]*MergeConflict, 0),
        Resolved:     0,
        Unresolved:   0,
        Skipped:      0,
        Warnings:     make([]string, 0),
        Errors:       make([]string, 0),
        Metadata:     make(map[string]interface{}),
    }

    // Start with the first action as base
    result.MergedAction = m.cloneAction(sortedActions[0])

    // Merge remaining actions
    for i := 1; i < len(sortedActions); i++ {
        mergeResult, err := m.mergeTwoActions(result.MergedAction, sortedActions[i], options)
        if err != nil {
            return nil, fmt.Errorf("failed to merge action %d: %w", i, err)
        }

        // Update result
        result.MergedAction = mergeResult.MergedAction
        result.Conflicts = append(result.Conflicts, mergeResult.Conflicts...)
        result.Resolved += mergeResult.Resolved
        result.Unresolved += mergeResult.Unresolved
        result.Skipped += mergeResult.Skipped
        result.Warnings = append(result.Warnings, mergeResult.Warnings...)
        result.Errors = append(result.Errors, mergeResult.Errors...)

        // Check conflict limit
        if len(result.Conflicts) >= options.ConflictLimit {
            result.Warnings = append(result.Warnings, "Conflict limit reached, stopping merge")
            break
        }
    }

    // Validate merged action if requested
    if options.ValidateResult {
        if err := m.validateMergedAction(result.MergedAction); err != nil {
            result.Errors = append(result.Errors, fmt.Sprintf("Validation failed: %v", err))
        }
    }

    result.Duration = time.Since(startTime)
    result.Metadata["action_count"] = len(actions)
    result.Metadata["conflict_count"] = len(result.Conflicts)

    return result, nil
}

// MergeWithPolicy merges actions using a specific merge policy
func (m *Manager) MergeWithPolicy(ctx context.Context, actions []*types.Action, policy MergePolicy, strategy MergeStrategy) (*MergeResult, error) {
    options := &MergeOptions{
        Policy:          policy,
        Strategy:        strategy,
        ResolveConflicts: true,
        ValidateResult:   true,
        GenerateReport:   true,
        ConflictLimit:    100,
        Timeout:          30 * time.Second,
    }

    return m.MergeActions(ctx, actions, options)
}
```

### Step 2: Implement Action Sorting and Preparation

```go
// sortActionsForMerge sorts actions for optimal merging
func (m *Manager) sortActionsForMerge(actions []*types.Action) []*types.Action {
    sorted := make([]*types.Action, len(actions))
    copy(sorted, actions)

    // Sort by priority (highest first), then by timestamp (oldest first)
    sort.Slice(sorted, func(i, j int) bool {
        if sorted[i].Priority != sorted[j].Priority {
            return sorted[i].Priority > sorted[j].Priority
        }
        return sorted[i].Timestamp.Before(sorted[j].Timestamp)
    })

    return sorted
}

// cloneAction creates a deep copy of an action
func (m *Manager) cloneAction(action *types.Action) *types.Action {
    if action == nil {
        return nil
    }

    cloned := &types.Action{
        Name:           action.Name,
        Type:           action.Type,
        Description:    action.Description,
        Machines:       make([]string, len(action.Machines)),
        Dependencies:   make([]string, len(action.Dependencies)),
        Command:        action.Command,
        Script:         action.Script,
        ScriptPath:     action.ScriptPath,
        Template:       action.Template,
        TemplatePath:   action.TemplatePath,
        Parameters:     make(map[string]interface{}),
        Environment:    make(map[string]string),
        Timeout:        action.Timeout,
        RetryCount:     action.RetryCount,
        RetryDelay:     action.RetryDelay,
        Priority:       action.Priority,
        Tags:           make([]string, len(action.Tags)),
        Metadata:       make(map[string]interface{}),
        Timestamp:      action.Timestamp,
    }

    // Copy slices
    copy(cloned.Machines, action.Machines)
    copy(cloned.Dependencies, action.Dependencies)
    copy(cloned.Tags, action.Tags)

    // Copy maps
    for k, v := range action.Parameters {
        cloned.Parameters[k] = v
    }
    for k, v := range action.Environment {
        cloned.Environment[k] = v
    }
    for k, v := range action.Metadata {
        cloned.Metadata[k] = v
    }

    return cloned
}
```

### Step 3: Implement Two-Action Merging

```go
// mergeTwoActions merges two actions together
func (m *Manager) mergeTwoActions(action1, action2 *types.Action, options *MergeOptions) (*MergeResult, error) {
    result := &MergeResult{
        MergedAction: m.cloneAction(action1),
        Conflicts:    make([]*MergeConflict, 0),
        Resolved:     0,
        Unresolved:   0,
        Skipped:      0,
        Warnings:     make([]string, 0),
        Errors:       make([]string, 0),
    }

    // Merge basic fields
    if err := m.mergeBasicFields(result.MergedAction, action2, result, options); err != nil {
        return nil, fmt.Errorf("failed to merge basic fields: %w", err)
    }

    // Merge complex fields
    if err := m.mergeComplexFields(result.MergedAction, action2, result, options); err != nil {
        return nil, fmt.Errorf("failed to merge complex fields: %w", err)
    }

    // Merge type-specific fields
    if err := m.mergeTypeSpecificFields(result.MergedAction, action2, result, options); err != nil {
        return nil, fmt.Errorf("failed to merge type-specific fields: %w", err)
    }

    return result, nil
}

// mergeBasicFields merges basic string and primitive fields
func (m *Manager) mergeBasicFields(merged, action2 *types.Action, result *MergeResult, options *MergeOptions) error {
    // Merge name
    if err := m.mergeField("name", merged.Name, action2.Name, merged, result, options); err != nil {
        return err
    }

    // Merge type
    if err := m.mergeField("type", merged.Type, action2.Type, merged, result, options); err != nil {
        return err
    }

    // Merge description
    if err := m.mergeField("description", merged.Description, action2.Description, merged, result, options); err != nil {
        return err
    }

    // Merge timeout
    if err := m.mergeField("timeout", merged.Timeout, action2.Timeout, merged, result, options); err != nil {
        return err
    }

    // Merge retry count
    if err := m.mergeField("retry_count", merged.RetryCount, action2.RetryCount, merged, result, options); err != nil {
        return err
    }

    // Merge retry delay
    if err := m.mergeField("retry_delay", merged.RetryDelay, action2.RetryDelay, merged, result, options); err != nil {
        return err
    }

    // Merge priority
    if err := m.mergeField("priority", merged.Priority, action2.Priority, merged, result, options); err != nil {
        return err
    }

    return nil
}

// mergeComplexFields merges complex fields like slices and maps
func (m *Manager) mergeComplexFields(merged, action2 *types.Action, result *MergeResult, options *MergeOptions) error {
    // Merge machines
    if err := m.mergeStringSlice("machines", merged.Machines, action2.Machines, merged, result, options); err != nil {
        return err
    }

    // Merge dependencies
    if err := m.mergeStringSlice("dependencies", merged.Dependencies, action2.Dependencies, merged, result, options); err != nil {
        return err
    }

    // Merge tags
    if err := m.mergeStringSlice("tags", merged.Tags, action2.Tags, merged, result, options); err != nil {
        return err
    }

    // Merge parameters
    if err := m.mergeMap("parameters", merged.Parameters, action2.Parameters, merged, result, options); err != nil {
        return err
    }

    // Merge environment
    if err := m.mergeMap("environment", merged.Environment, action2.Environment, merged, result, options); err != nil {
        return err
    }

    // Merge metadata
    if err := m.mergeMap("metadata", merged.Metadata, action2.Metadata, merged, result, options); err != nil {
        return err
    }

    return nil
}

// mergeTypeSpecificFields merges fields specific to action types
func (m *Manager) mergeTypeSpecificFields(merged, action2 *types.Action, result *MergeResult, options *MergeOptions) error {
    // Merge command-specific fields
    if merged.Type == "command" || action2.Type == "command" {
        if err := m.mergeField("command", merged.Command, action2.Command, merged, result, options); err != nil {
            return err
        }
    }

    // Merge script-specific fields
    if merged.Type == "script" || action2.Type == "script" {
        if err := m.mergeField("script", merged.Script, action2.Script, merged, result, options); err != nil {
            return err
        }
        if err := m.mergeField("script_path", merged.ScriptPath, action2.ScriptPath, merged, result, options); err != nil {
            return err
        }
    }

    // Merge template-specific fields
    if merged.Type == "template" || action2.Type == "template" {
        if err := m.mergeField("template", merged.Template, action2.Template, merged, result, options); err != nil {
            return err
        }
        if err := m.mergeField("template_path", merged.TemplatePath, action2.TemplatePath, merged, result, options); err != nil {
            return err
        }
    }

    return nil
}
```

### Step 4: Implement Field Merging Logic

```go
// mergeField merges a single field with conflict resolution
func (m *Manager) mergeField(fieldName string, value1, value2 interface{}, merged *types.Action, result *MergeResult, options *MergeOptions) error {
    // Check if values are equal
    if m.valuesEqual(value1, value2) {
        return nil
    }

    // Check if one value is empty/null
    if m.isEmptyValue(value1) && !m.isEmptyValue(value2) {
        m.setFieldValue(merged, fieldName, value2)
        return nil
    }

    if !m.isEmptyValue(value1) && m.isEmptyValue(value2) {
        return nil
    }

    // Both values are non-empty and different - conflict
    conflict := &MergeConflict{
        Field:        fieldName,
        Action1:      merged,
        Action2:      nil, // Will be set by caller
        Value1:       value1,
        Value2:       value2,
        ConflictType: ConflictTypeFieldValue,
        Severity:     m.determineConflictSeverity(fieldName),
        Resolution:   ConflictResolution{},
    }

    // Apply merge strategy
    resolution, err := m.resolveConflict(conflict, options)
    if err != nil {
        result.Unresolved++
        result.Conflicts = append(result.Conflicts, conflict)
        return err
    }

    conflict.Resolution = resolution
    if resolution.Applied {
        m.setFieldValue(merged, fieldName, resolution.Value)
        result.Resolved++
    } else {
        result.Skipped++
    }

    result.Conflicts = append(result.Conflicts, conflict)
    return nil
}

// mergeStringSlice merges string slices
func (m *Manager) mergeStringSlice(fieldName string, slice1, slice2 []string, merged *types.Action, result *MergeResult, options *MergeOptions) error {
    if len(slice2) == 0 {
        return nil
    }

    if len(slice1) == 0 {
        m.setFieldValue(merged, fieldName, slice2)
        return nil
    }

    // Check for conflicts
    conflicts := m.findSliceConflicts(slice1, slice2)
    if len(conflicts) == 0 {
        // No conflicts, merge based on policy
        mergedSlice := m.mergeStringSlices(slice1, slice2, options.Policy)
        m.setFieldValue(merged, fieldName, mergedSlice)
        return nil
    }

    // Handle conflicts
    for _, conflict := range conflicts {
        resolution, err := m.resolveSliceConflict(conflict, options)
        if err != nil {
            result.Unresolved++
            result.Conflicts = append(result.Conflicts, conflict)
        } else {
            conflict.Resolution = resolution
            if resolution.Applied {
                result.Resolved++
            } else {
                result.Skipped++
            }
            result.Conflicts = append(result.Conflicts, conflict)
        }
    }

    return nil
}

// mergeMap merges maps
func (m *Manager) mergeMap(fieldName string, map1, map2 map[string]interface{}, merged *types.Action, result *MergeResult, options *MergeOptions) error {
    if len(map2) == 0 {
        return nil
    }

    if len(map1) == 0 {
        m.setFieldValue(merged, fieldName, map2)
        return nil
    }

    // Check for conflicts
    conflicts := m.findMapConflicts(map1, map2)
    if len(conflicts) == 0 {
        // No conflicts, merge based on policy
        mergedMap := m.mergeMaps(map1, map2, options.Policy)
        m.setFieldValue(merged, fieldName, mergedMap)
        return nil
    }

    // Handle conflicts
    for _, conflict := range conflicts {
        resolution, err := m.resolveMapConflict(conflict, options)
        if err != nil {
            result.Unresolved++
            result.Conflicts = append(result.Conflicts, conflict)
        } else {
            conflict.Resolution = resolution
            if resolution.Applied {
                result.Resolved++
            } else {
                result.Skipped++
            }
            result.Conflicts = append(result.Conflicts, conflict)
        }
    }

    return nil
}
```

### Step 5: Implement Conflict Resolution

```go
// resolveConflict resolves a merge conflict based on strategy
func (m *Manager) resolveConflict(conflict *MergeConflict, options *MergeOptions) (ConflictResolution, error) {
    resolution := ConflictResolution{
        Timestamp: time.Now(),
    }

    switch options.Strategy {
    case MergeStrategyFirstWins:
        resolution.Method = "first_wins"
        resolution.Value = conflict.Value1
        resolution.Applied = true
        resolution.Reason = "First action value preserved"

    case MergeStrategyLastWins:
        resolution.Method = "last_wins"
        resolution.Value = conflict.Value2
        resolution.Applied = true
        resolution.Reason = "Last action value applied"

    case MergeStrategyCombine:
        combined, err := m.combineValues(conflict.Value1, conflict.Value2, conflict.Field)
        if err != nil {
            return resolution, fmt.Errorf("failed to combine values: %w", err)
        }
        resolution.Method = "combine"
        resolution.Value = combined
        resolution.Applied = true
        resolution.Reason = "Values combined"

    case MergeStrategyFail:
        return resolution, fmt.Errorf("conflict resolution failed for field %s", conflict.Field)

    case MergeStrategyPrompt:
        // In a real implementation, this would prompt the user
        resolution.Method = "prompt"
        resolution.Value = conflict.Value1
        resolution.Applied = false
        resolution.Reason = "User prompt required"

    default:
        return resolution, fmt.Errorf("unknown merge strategy: %s", options.Strategy)
    }

    return resolution, nil
}

// combineValues combines two values based on field type
func (m *Manager) combineValues(value1, value2 interface{}, fieldName string) (interface{}, error) {
    switch fieldName {
    case "machines", "dependencies", "tags":
        return m.combineStringSlices(value1, value2)
    case "parameters", "environment", "metadata":
        return m.combineMaps(value1, value2)
    case "description":
        return m.combineDescriptions(value1, value2)
    default:
        return value2, nil // Default to second value
    }
}

// combineStringSlices combines string slices
func (m *Manager) combineStringSlices(value1, value2 interface{}) (interface{}, error) {
    slice1, ok1 := value1.([]string)
    slice2, ok2 := value2.([]string)

    if !ok1 || !ok2 {
        return nil, fmt.Errorf("values are not string slices")
    }

    // Create a map to track unique values
    unique := make(map[string]bool)
    result := make([]string, 0)

    // Add all values from both slices
    for _, item := range slice1 {
        if !unique[item] {
            unique[item] = true
            result = append(result, item)
        }
    }

    for _, item := range slice2 {
        if !unique[item] {
            unique[item] = true
            result = append(result, item)
        }
    }

    return result, nil
}

// combineMaps combines maps
func (m *Manager) combineMaps(value1, value2 interface{}) (interface{}, error) {
    map1, ok1 := value1.(map[string]interface{})
    map2, ok2 := value2.(map[string]interface{})

    if !ok1 || !ok2 {
        return nil, fmt.Errorf("values are not maps")
    }

    result := make(map[string]interface{})

    // Copy all values from first map
    for k, v := range map1 {
        result[k] = v
    }

    // Add/override with values from second map
    for k, v := range map2 {
        result[k] = v
    }

    return result, nil
}

// combineDescriptions combines descriptions
func (m *Manager) combineDescriptions(value1, value2 interface{}) (interface{}, error) {
    desc1, ok1 := value1.(string)
    desc2, ok2 := value2.(string)

    if !ok1 || !ok2 {
        return nil, fmt.Errorf("values are not strings")
    }

    if desc1 == "" {
        return desc2, nil
    }

    if desc2 == "" {
        return desc1, nil
    }

    // Combine descriptions with separator
    return fmt.Sprintf("%s\n\n%s", desc1, desc2), nil
}
```

### Step 6: Implement Validation and Reporting

```go
// validateMergedAction validates the merged action
func (m *Manager) validateMergedAction(action *types.Action) error {
    // Basic validation
    if action.Name == "" {
        return fmt.Errorf("merged action has no name")
    }

    if action.Type == "" {
        return fmt.Errorf("merged action has no type")
    }

    if len(action.Machines) == 0 {
        return fmt.Errorf("merged action has no target machines")
    }

    // Type-specific validation
    switch action.Type {
    case "command":
        if action.Command == "" {
            return fmt.Errorf("command action has no command")
        }
    case "script":
        if action.Script == "" && action.ScriptPath == "" {
            return fmt.Errorf("script action has no script content or path")
        }
    case "template":
        if action.Template == "" && action.TemplatePath == "" {
            return fmt.Errorf("template action has no template content or path")
        }
    }

    return nil
}

// generateMergeReport generates a detailed merge report
func (m *Manager) generateMergeReport(result *MergeResult) string {
    var report strings.Builder

    report.WriteString("=== Action Merge Report ===\n\n")

    // Summary
    report.WriteString(fmt.Sprintf("Total Conflicts: %d\n", len(result.Conflicts)))
    report.WriteString(fmt.Sprintf("Resolved: %d\n", result.Resolved))
    report.WriteString(fmt.Sprintf("Unresolved: %d\n", result.Unresolved))
    report.WriteString(fmt.Sprintf("Skipped: %d\n", result.Skipped))
    report.WriteString(fmt.Sprintf("Duration: %v\n\n", result.Duration))

    // Warnings
    if len(result.Warnings) > 0 {
        report.WriteString("Warnings:\n")
        for _, warning := range result.Warnings {
            report.WriteString(fmt.Sprintf("  - %s\n", warning))
        }
        report.WriteString("\n")
    }

    // Errors
    if len(result.Errors) > 0 {
        report.WriteString("Errors:\n")
        for _, err := range result.Errors {
            report.WriteString(fmt.Sprintf("  - %s\n", err))
        }
        report.WriteString("\n")
    }

    // Conflicts
    if len(result.Conflicts) > 0 {
        report.WriteString("Conflicts:\n")
        for i, conflict := range result.Conflicts {
            report.WriteString(fmt.Sprintf("  %d. Field: %s\n", i+1, conflict.Field))
            report.WriteString(fmt.Sprintf("     Type: %s\n", conflict.ConflictType))
            report.WriteString(fmt.Sprintf("     Severity: %s\n", conflict.Severity))
            report.WriteString(fmt.Sprintf("     Value 1: %v\n", conflict.Value1))
            report.WriteString(fmt.Sprintf("     Value 2: %v\n", conflict.Value2))
            if conflict.Resolution.Applied {
                report.WriteString(fmt.Sprintf("     Resolution: %s (%s)\n", conflict.Resolution.Method, conflict.Resolution.Reason))
            } else {
                report.WriteString("     Resolution: Not applied\n")
            }
            report.WriteString("\n")
        }
    }

    return report.String()
}
```







## Configuration Options

### Supported Options
- **Merge policy**: Replace, append, merge, skip, custom
- **Conflict strategy**: First wins, last wins, combine, prompt, fail
- **Validation**: Enable/disable result validation
- **Reporting**: Enable/disable detailed reports

## Dependencies

### Internal Dependencies
- `spooky/internal/actions/types`
- `spooky/internal/actions/validation`
- `spooky/internal/actions/merging`
- `spooky/internal/logging`

### External Dependencies
- `sort` (standard library)
- `strings` (standard library)
- `time` (standard library)



## Implementation Order

1. Enhance action manager with merging functions
2. Implement action sorting and preparation
3. Add two-action merging logic
4. Implement field merging logic
5. Add conflict resolution
6. Implement validation and reporting
7. Write comprehensive tests
8. Performance testing and optimization
9. Documentation and cleanup


