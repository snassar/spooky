# Variables System: Comprehensive Implementation Plan

## Overview

This document is the authoritative source for all variable system implementation details in spooky. It covers variable types, validation, interpolation, dependency management, file merging, and integration with the facts and template systems.

**Schema Integration**: This variables system leverages the schema validation infrastructure and patterns defined in [Schema System](../schema-system.md) for comprehensive variable validation, type checking, and schema-based export functionality.

**Architecture Integration**: Variables integrate with the overall spooky architecture as described in [Spooky Design](../spooky-design.md), providing configuration values to the actions system for template rendering and running.

## System Integration

This variables system integrates with other core Spooky systems to provide comprehensive variable management and configuration:

### **Facts System Integration**
- **Facts as Variables**: Facts automatically available as variables in templates (see [Facts System](../facts-system.md))
- **Variable Resolution**: Facts can be referenced and resolved in variable definitions
- **Template Context**: Facts and variables merge in unified template context
- **Export Integration**: Variable exports include facts context information

### **Project System Integration**
- **Project Variables**: Variables stored in project-specific `variables.hcl` and `variables/` directory (see [Project System](../project-system.md))
- **Project Context**: Variables resolved within project run context
- **Project Isolation**: Project variables isolated from global variables
- **Project Configuration**: Variable settings configured in project context

### **Schema System Integration**
- **Variable Schema**: Variable definitions validated against embedded schemas (see [Schema System](../schema-system.md))
- **Type Validation**: Variable types enforced through schema validation
- **Schema Evolution**: Variable schemas evolve with system changes
- **Schema Documentation**: Variable schemas documented and versioned

### **CLI System Integration**
- **Variables Commands**: Variables management through `spooky variables` commands (see [CLI System](../cli-system.md))
- **Variables Validation**: `spooky variables validate` for variable file validation
- **Variables Discovery**: `spooky variables list` for variable discovery and filtering
- **Variables Display**: `spooky variables show` for detailed variable information
- **Variables Export**: `spooky variables export` for variable data export

### **Actions System Integration**
- **Variables in Actions**: Actions can use project variables for dynamic configuration (see [Actions System](../actions-system.md))
- **Variable Resolution**: Actions resolve variables at run time
- **Template Integration**: Actions can use variables in command templates
- **Variable Precedence**: Actions follow variable resolution precedence (Environment → Project → Facts → Defaults)

### **Configuration System Integration**
- **Global Configuration**: Variable resolution settings from `$XDG_CONFIG_HOME/spooky/spooky.hcl` (see [Configuration System](../configuration-system.md))
- **Variable Settings**: Variable storage, caching, and performance settings
- **Security Settings**: Variable encryption and sensitive data handling
- **Performance Settings**: Variable resolution performance tuning
- **Configuration Precedence**: Variable resolution follows configuration precedence hierarchy
- **Variable Resolution**: Configuration settings for variable resolution (see [Configuration System](../configuration-system.md))
- **Variable Storage**: Configuration for variable storage and caching
- **Variable Security**: Configuration for sensitive variable handling
- **Variable Performance**: Configuration for variable resolution performance
- **Variable Validation**: Configuration for variable validation settings

### **Machines System Integration**
- **Machine Variables**: Machine inventory can use project variables for dynamic configuration (see [Machines System](../machines-system.md))
- **Variable Interpolation**: Machine hostnames, usernames, and tags can use variable interpolation
- **Environment Variables**: Machine configuration can reference environment variables
- **Dynamic Inventory**: Machine inventory can be generated from variables and external sources
- **Machine Context**: Variables available in machine-specific contexts

### **Template System Integration**
- **Variables in Templates**: Templates can access project variables for dynamic configuration (see [Template System](../template-system.md))
- **Variable Interpolation**: Templates support variable interpolation and resolution
- **Variable Functions**: Template functions for accessing variables (`var()`, `varOrDefault()`)
- **Variable Precedence**: Templates follow variable resolution precedence
- **Template Context**: Variables available in template rendering context

## Implementation Status

### **✅ Completed**
- **Environment variables** in templates via `{{env "VAR_NAME"}}`
- **Template context** (facts, machines, actions, custom facts)
- **Tags** for machines and projects
- **Actions merging** (conceptually similar to what we need for variables)
- **Variables system** with `variables.hcl` and `variables/` directory (schema defined)
- **Variable scoping** and precedence rules (Environment > Project > Facts > Defaults)

### **🔄 In Progress / Planned**
- **File merging** for variables with conflict detection and resolution
- **Dependency management** to prevent circular references
- **Variable interpolation** in HCL files (`${var.name}` syntax)
- **Schema validation** integration with embedded variable schemas
- **CLI commands** for variable management (`spooky variables validate`, `spooky variables list`, etc.)
- **Template integration** for variable usage in actions

## Variables System Design

### **1. File Structure**

```
project/
├── variables.hcl          # Main variables file
├── variables/             # Organized variable files
│   ├── environment.hcl    # Environment-specific variables
│   ├── secrets.hcl        # Sensitive variables
│   ├── defaults.hcl       # Default values
│   └── computed.hcl       # Computed/dependent variables
├── actions.hcl            # Main actions file
├── actions/               # Organized action files
│   ├── setup.hcl          # Setup actions
│   ├── deploy.hcl         # Deployment actions
│   └── cleanup.hcl        # Cleanup actions
└── ...
```

### **2. Variable Definition Schema**

```hcl
# variables.hcl
variable "app_name" {
  type        = "string"
  description = "Application name"
  default     = "myapp"
  required    = false
  sensitive   = false
}

variable "database_url" {
  type        = "string"
  description = "Database connection URL"
  required    = true
  sensitive   = true
  depends_on  = ["database_host", "database_name"]
}

variable "server_count" {
  type        = "number"
  description = "Number of servers to deploy"
  default     = 3
  validation {
    condition     = var.server_count > 0
    error_message = "Server count must be greater than 0"
  }
}

variable "environment" {
  type        = "string"
  description = "Deployment environment"
  default     = "development"
  validation {
    condition     = contains(["development", "staging", "production"], var.environment)
    error_message = "Environment must be development, staging, or production"
  }
}
```

### **3. Variable Types and Validation**

```hcl
# Supported variable types
variable "string_var" {
  type = "string"
}

variable "number_var" {
  type = "number"
}

variable "bool_var" {
  type = "bool"
}

variable "list_var" {
  type = "list(string)"
}

variable "map_var" {
  type = "map(string)"
}

variable "object_var" {
  type = "object({
    name = string
    port = number
    ssl  = bool
  })"
}
```

## File Merging Implementation

### **1. Merging Strategy**

Based on best practices from Ansible, Chef, and Puppet, file merging should be integrated into the existing config parsing system rather than implemented as a separate package.

```go
// Extend existing config parsing functions
// Current system uses individual parsing functions:
// - ParseProjectConfig()
// - ParseInventoryConfig() 
// - ParseActionsConfig()
// - LoadActionsConfig()

// Add new variable parsing functions following the same pattern
func ParseVariablesConfig(filename string) (*VariablesConfig, error) {
    // Use existing HCL parsing infrastructure
    parser := hclparse.NewParser()
    file, diags := parser.ParseHCLFile(filename)
    if diags.HasErrors() {
        return nil, fmt.Errorf("failed to parse variables HCL file: %s", diags.Error())
    }
    
    // Decode using wrapper pattern (consistent with existing code)
    var wrapper VariablesWrapper
    diags = gohcl.DecodeBody(file.Body, nil, &wrapper)
    if diags.HasErrors() {
        return nil, fmt.Errorf("failed to decode variables configuration: %s", diags.Error())
    }
    
    if wrapper.Variables == nil {
        return nil, errors.New("no variables block found in configuration")
    }
    
    return wrapper.Variables, nil
}

// LoadVariablesConfig loads variables from multiple sources and merges them
// Following the same pattern as LoadActionsConfig
func LoadVariablesConfig(projectPath string) (*VariablesConfig, error) {
    logger := logging.GetLogger()
    
    // Initialize merged config
    mergedConfig := &VariablesConfig{
        Variables: []Variable{},
    }
    
    // 1. Try to load variables.hcl from project root
    rootVariablesFile := filepath.Join(projectPath, "variables.hcl")
    if _, err := os.Stat(rootVariablesFile); err == nil {
        logger.Info("Loading variables from root file", logging.String("file", rootVariablesFile))
        rootConfig, err := ParseVariablesConfig(rootVariablesFile)
        if err != nil {
            logger.Error("Failed to parse root variables file", err, logging.String("file", rootVariablesFile))
            return nil, fmt.Errorf("failed to parse root variables file: %w", err)
        }
        mergedConfig.Variables = append(mergedConfig.Variables, rootConfig.Variables...)
        logger.Info("Loaded variables from root file", logging.Int("variables", len(rootConfig.Variables)))
    }
    
    // 2. Try to load all .hcl files from variables/ directory
    variablesDir := filepath.Join(projectPath, "variables")
    if _, err := os.Stat(variablesDir); err == nil {
        logger.Info("Loading variables from directory", logging.String("dir", variablesDir))
        
        entries, err := os.ReadDir(variablesDir)
        if err != nil {
            logger.Error("Failed to read variables directory", err, logging.String("dir", variablesDir))
            return nil, fmt.Errorf("failed to read variables directory: %w", err)
        }
        
        // Sort entries to ensure consistent loading order
        var variableFiles []string
        for _, entry := range entries {
            if !entry.IsDir() && filepath.Ext(entry.Name()) == ".hcl" {
                variableFiles = append(variableFiles, entry.Name())
            }
        }
        
        // Sort files to ensure consistent loading order
        sort.Strings(variableFiles)
        
        for _, fileName := range variableFiles {
            filePath := filepath.Join(variablesDir, fileName)
            logger.Info("Loading variable file", logging.String("file", filePath))
            
            fileConfig, err := ParseVariablesConfig(filePath)
            if err != nil {
                logger.Error("Failed to parse variable file", err, logging.String("file", filePath))
                return nil, fmt.Errorf("failed to parse variable file %s: %w", fileName, err)
            }
            
            mergedConfig.Variables = append(mergedConfig.Variables, fileConfig.Variables...)
            logger.Info("Loaded variables from file",
                logging.String("file", fileName),
                logging.Int("variables", len(fileConfig.Variables)))
        }
    }
    
    // Check if we loaded any variables
    if len(mergedConfig.Variables) == 0 {
        logger.Warn("No variables found in project", logging.String("project_path", projectPath))
        return mergedConfig, nil
    }
    
    logger.Info("Successfully loaded all variables",
        logging.String("project_path", projectPath),
        logging.Int("total_variables", len(mergedConfig.Variables)))
    
    return mergedConfig, nil
}

// Integrate with existing project loading in TemplateContext
func (ctx *TemplateContext) loadVariables(logger logging.Logger, projectPath string) error {
    logger.Info("Loading variables from project", logging.String("project_path", projectPath))
    
    // Use the new LoadVariablesConfig function that can load from multiple sources
    variablesConfig, err := config.LoadVariablesConfig(projectPath)
    if err != nil {
        logger.Error("Failed to load variables", err, logging.String("project_path", projectPath))
        return fmt.Errorf("failed to load variables: %w", err)
    }
    
    // Convert to VariableContext
    ctx.Variables = &VariableContext{
        ProjectVars: make(map[string]interface{}),
        SourceMap:   make(map[string]VariableSource),
    }
    
    // Process variables and build source map
    for _, variable := range variablesConfig.Variables {
        ctx.Variables.ProjectVars[variable.Name] = variable.Default
        ctx.Variables.SourceMap[variable.Name] = VariableSource{
            Type: "project",
            File: variable.SourceFile,
            Line: variable.Line,
        }
    }
    
    logger.Info("Variables loaded successfully",
        logging.String("project_path", projectPath),
        logging.Int("variables", len(ctx.Variables.ProjectVars)))
    
    return nil
}

type FileMerger struct {
    MaxFileSize    int64  // Default: 100MB
    TempFileSize   int64  // Default: 10MB
    MergeStrategy  string // "memory" or "tempfile"
    StrictMode     bool   // Default: true for validation
}

type MergedContent struct {
    Content     interface{}
    SourceFiles []string
    Dependencies map[string][]string
    CircularRefs []string
}
```

### **2. Conflict Handling**

Variable conflicts should be treated as validation errors. If the same variable is defined in both `variables.hcl` and `variables/environment.hcl`, `spooky variables validate` should error on this.

```go
type VariableMerger struct {
    // ... existing fields
    strictMode bool // Default: true for validation
}

func (vm *VariableMerger) MergeFiles(files []string, project *Project) error {
    variables := make(map[string]*VariableDefinition)
    conflicts := []VariableConflict{}
    
    for _, file := range files {
        fileVars, err := vm.parseVariableFile(file)
        if err != nil {
            return fmt.Errorf("error parsing %s: %w", filepath.Base(file), err)
        }
        
        for name, varDef := range fileVars {
            if existing, exists := variables[name]; exists {
                // Conflict detected
                conflict := VariableConflict{
                    VariableName: name,
                    FirstFile:    existing.SourceFile,
                    FirstLine:    existing.Line,
                    SecondFile:   file,
                    SecondLine:   varDef.Line,
                }
                conflicts = append(conflicts, conflict)
                
                // In strict mode (validation), continue collecting all conflicts
                if vm.strictMode {
                    continue
                }
                
                // In non-strict mode, last file wins (for backward compatibility)
                // But log a warning
                log.Warnf("Variable '%s' redefined in %s:%d (previously defined in %s:%d)", 
                    name, filepath.Base(file), varDef.Line, 
                    filepath.Base(existing.SourceFile), existing.Line)
            }
            
            variables[name] = varDef
        }
    }
    
    // Report all conflicts in validation mode
    if len(conflicts) > 0 {
        return &VariableConflictError{Conflicts: conflicts}
    }
    
    project.Variables = variables
    return nil
}

type VariableConflict struct {
    VariableName string
    FirstFile    string
    FirstLine    int
    SecondFile   string
    SecondLine   int
}

type VariableConflictError struct {
    Conflicts []VariableConflict
}

func (e *VariableConflictError) Error() string {
    var buf strings.Builder
    buf.WriteString("Variable conflicts detected:\n")
    for _, conflict := range e.Conflicts {
        buf.WriteString(fmt.Sprintf("  Variable '%s' defined in:\n", conflict.VariableName))
        buf.WriteString(fmt.Sprintf("    %s:%d\n", filepath.Base(conflict.FirstFile), conflict.FirstLine))
        buf.WriteString(fmt.Sprintf("    %s:%d\n", filepath.Base(conflict.SecondFile), conflict.SecondLine))
    }
    return buf.String()
}
```

**CLI Integration Example:**
```bash
$ spooky variables validate ./my-project
Error: Variable conflicts detected:
  Variable 'app_name' defined in:
    variables.hcl:5
    variables/environment.hcl:3
  Variable 'database_url' defined in:
    variables.hcl:12
    variables/secrets.hcl:8
❌ Variables validation failed (2 conflicts)
```

### **3. Parse-Time vs Runtime Merging**

Based on how Ansible, Chef, and Puppet handle this, merging should be done at **parse time** with runtime validation.

#### **How Other Tools Handle This:**

**Ansible:**
```yaml
# Ansible loads all variable files at parse time
# Variables are resolved at runtime, but file loading is done upfront
- name: Load variables
  include_vars: "{{ item }}"
  loop:
    - variables.yml
    - variables/{{ environment }}.yml
    - variables/secrets.yml
```

**Puppet:**
```puppet
# Puppet loads all manifests at parse time
# Hiera data is loaded at parse time, but resolved at runtime
class myapp {
  $app_name = hiera('app_name')
  $database_url = hiera('database_url')
}
```

**Chef:**
```ruby
# Chef loads all cookbook files at parse time
# Attributes are merged at parse time, but resolved at runtime
default['myapp']['name'] = 'myapp'
default['myapp']['database_url'] = 'postgresql://localhost/myapp'
```

#### **Recommended Approach for Spooky:**

```go
// Parse time: Load and merge all variable files
func (vm *VariableMerger) LoadVariables(projectPath string) (*VariableContext, error) {
    files := vm.discoverVariableFiles(projectPath)
    variables, conflicts, err := vm.mergeFiles(files)
    if err != nil {
        return nil, err
    }
    
    if len(conflicts) > 0 {
        return nil, &VariableConflictError{Conflicts: conflicts}
    }
    
    return &VariableContext{
        ProjectVars: variables,
        SourceMap:   vm.buildSourceMap(files),
        Dependencies: vm.buildDependencyGraph(variables),
    }, nil
}

// Runtime: Resolve variables with facts and environment
func (vr *VariableResolver) Resolve(name string) (interface{}, error) {
    // Variables are already loaded and merged at parse time
    // Runtime resolution only handles facts, environment, and interpolation
    if value, exists := vr.context.Variables.ProjectVars[name]; exists {
        return vr.interpolateProjectVar(value)
    }
    
    // Check other sources (facts, env, defaults)
    return vr.resolveFromOtherSources(name)
}
```

**Benefits of Parse-Time Merging:**

1. **Early Error Detection**: Conflicts and syntax errors caught before running
2. **Performance**: No file I/O during variable resolution
3. **Consistency**: All variable files loaded once, ensuring consistent state
4. **Validation**: Can validate dependencies and circular references upfront
5. **Caching**: Merged variables can be cached for multiple operations

**Runtime Resolution:**
- **Facts**: Resolved at runtime for each machine
- **Environment Variables**: Checked at runtime for security
- **Interpolation**: Variable references resolved at runtime
- **Defaults**: Applied at runtime if no other source provides value

**Implementation Example:**
```go
// Parse time (during project loading)
func (p *Project) Load() error {
    // Load and merge all variable files
    variables, err := p.variableMerger.LoadVariables(p.Path)
    if err != nil {
        return fmt.Errorf("failed to load variables: %w", err)
    }
    
    p.VariableContext = variables
    return nil
}

// Runtime (during template rendering or action running)
func (p *Project) ResolveVariable(name string) (interface{}, error) {
    resolver := &VariableResolver{context: p.TemplateContext}
    return resolver.Resolve(name)
}
```

### **4. File Size Handling**

```go
func (fm *FileMerger) shouldUseTempFile(files []string) bool {
    totalSize := int64(0)
    for _, file := range files {
        if info, err := os.Stat(file); err == nil {
            totalSize += info.Size()
        }
    }
    return totalSize > fm.TempFileSize
}
```

### **5. Merging Process**

```go
func (fm *FileMerger) MergeVariables(projectPath string) (*MergedContent, error) {
    // 1. Discover variable files using existing patterns
    files := fm.discoverVariableFiles(projectPath)
    
    // 2. Validate file locations (only project directory files)
    for _, file := range files {
        if err := fm.validateFilePath(file, projectPath); err != nil {
            return nil, err
        }
    }
    
    // 3. Choose merging strategy based on file size
    if fm.shouldUseTempFile(files) {
        return fm.mergeWithTempFile(files)
    }
    
    // 4. Merge in memory
    return fm.mergeInMemory(files)
    
    // 5. Validate dependencies and detect circular references
    // 6. Return merged content with metadata
}
```

## Dependency Management

### **1. Dependency Graph**

```go
type DependencyNode struct {
    Name         string
    Dependencies []string
    Dependents   []string
    File         string
    Line         int
}

type DependencyGraph struct {
    Nodes map[string]*DependencyNode
    Edges map[string][]string
}
```

### **2. Circular Reference Detection**

```go
func (dg *DependencyGraph) DetectCircularRefs() []string {
    visited := make(map[string]bool)
    recStack := make(map[string]bool)
    circular := []string{}
    
    for node := range dg.Nodes {
        if !visited[node] {
            if dg.hasCycle(node, visited, recStack, &circular) {
                return circular
            }
        }
    }
    return nil
}

func (dg *DependencyGraph) hasCycle(node string, visited, recStack map[string]bool, circular *[]string) bool {
    visited[node] = true
    recStack[node] = true
    
    for _, dep := range dg.Nodes[node].Dependencies {
        if !visited[dep] {
            if dg.hasCycle(dep, visited, recStack, circular) {
                *circular = append(*circular, node)
                return true
            }
        } else if recStack[dep] {
            *circular = append(*circular, node)
            return true
        }
    }
    
    recStack[node] = false
    return false
}
```

### **3. Dependency Resolution**

```go
func (dg *DependencyGraph) ResolveOrder() ([]string, error) {
    // Topological sort for dependency resolution
    inDegree := make(map[string]int)
    queue := []string{}
    result := []string{}
    
    // Calculate in-degrees
    for _, node := range dg.Nodes {
        inDegree[node.Name] = len(node.Dependencies)
        if inDegree[node.Name] == 0 {
            queue = append(queue, node.Name)
        }
    }
    
    // Process queue
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        result = append(result, current)
        
        for _, dependent := range dg.Nodes[current].Dependents {
            inDegree[dependent]--
            if inDegree[dependent] == 0 {
                queue = append(queue, dependent)
            }
        }
    }
    
    if len(result) != len(dg.Nodes) {
        return nil, fmt.Errorf("circular dependency detected")
    }
    
    return result, nil
}
```

## Variable Storage and Resolution

### **1. Variable Storage Architecture**

Based on best practices from Ansible, Puppet, and Chef, variables should be integrated into the existing template context for seamless access across all systems.

```go
// Extend existing TemplateContext for variable integration
type TemplateContext struct {
    // Existing context fields
    Facts     map[string]interface{}
    Machines  map[string]interface{}
    Actions   map[string]interface{}
    Tags      map[string]string
    Env       map[string]string
    
    // Enhanced variable context
    Variables *VariableContext
}

type VariableContext struct {
    // Hierarchical storage for different variable sources
    ProjectVars    map[string]interface{} // variables.hcl and variables/*.hcl
    FactsVars      map[string]interface{} // facts-derived variables (same as Facts)
    EnvironmentVars map[string]string     // environment variables
    DefaultVars    map[string]interface{} // system defaults
    
    // Metadata for debugging and validation
    SourceMap      map[string]VariableSource // tracks where each variable came from
    Dependencies   *DependencyGraph
    ResolutionOrder []string // order variables were resolved
    Cache          map[string]interface{} // cache resolved values for performance
}

type VariableSource struct {
    Type     string // "project", "facts", "env", "default"
    File     string // source file (if applicable)
    Line     int    // line number (if applicable)
    Priority int    // resolution priority
}
```

### **2. Variable Resolution with Facts System**

The variables system integrates with the facts system (see [Facts System](../facts-system.md)) to provide seamless access to machine facts as variables in templates and actions:

#### **How Other Tools Handle This:**

**Ansible:**
```yaml
# Facts are automatically available as variables
- name: Use facts in variables
  vars:
    server_name: "{{ ansible_hostname }}"
    os_version: "{{ ansible_distribution_version }}"
    log_path: "/var/log/{{ ansible_hostname }}"
```

**Puppet:**
```puppet
# Facts are automatically available via Facter
$hostname = $facts['hostname']
$log_directory = "/var/log/${facts['hostname']}"
$config_path = "/etc/${facts['os']['name']}"
```

**Chef:**
```ruby
# Facts are available as node attributes via Ohai
hostname = node['hostname']
log_directory = "/var/log/#{node['hostname']}"
config_path = "/etc/#{node['platform']}"
```

#### **Recommended Approach for Spooky:**

```go
// Variable resolution precedence (highest to lowest)
const (
    PriorityEnvironment = 1000  // Environment variables
    PriorityProject    = 900   // Project variables (variables.hcl, variables/*.hcl)
    PriorityFacts      = 800   // Machine facts
    PriorityDefault    = 100   // System defaults
)

type VariableResolver struct {
    context *TemplateContext
    cache   map[string]interface{} // Cache resolved values
}

func (vr *VariableResolver) Resolve(name string) (interface{}, error) {
    // Check cache first for performance
    if value, exists := vr.cache[name]; exists {
        return value, nil
    }
    
    // 1. Environment Variables (highest priority)
    if value, exists := os.LookupEnv("SPOOKY_" + strings.ToUpper(name)); exists {
        vr.cache[name] = value
        return value, nil
    }
    
    // 2. Project Variables (can reference facts and other variables)
    if value, exists := vr.context.Variables.ProjectVars[name]; exists {
        resolved, err := vr.interpolateProjectVar(value)
        if err != nil {
            return nil, err
        }
        vr.cache[name] = resolved
        return resolved, nil
    }
    
    // 3. Facts (automatically available as variables)
    if value, exists := vr.context.Facts[name]; exists {
        vr.cache[name] = value
        return value, nil
    }
    
    // 4. System Defaults (lowest priority)
    if value, exists := vr.context.Variables.DefaultVars[name]; exists {
        vr.cache[name] = value
        return value, nil
    }
    
    return nil, fmt.Errorf("variable '%s' not found in any source", name)
}

// Handle nested variable interpolation within project variables
func (vr *VariableResolver) interpolateProjectVar(value interface{}) (interface{}, error) {
    switch v := value.(type) {
    case string:
        // Replace ${var.fact_name} with actual fact values
        return vr.interpolateString(v)
    case map[string]interface{}:
        // Handle nested maps
        result := make(map[string]interface{})
        for key, val := range v {
            resolved, err := vr.interpolateProjectVar(val)
            if err != nil {
                return nil, err
            }
            result[key] = resolved
        }
        return result, nil
    default:
        return value, nil
    }
}

// Example usage in variables.hcl:
// variable "log_path" {
//   type = "string"
//   default = "/var/log/${var.hostname}"  // References fact
// }
// 
// variable "config_dir" {
//   type = "string"
//   default = "/etc/${var.os_family}"     // References fact
// }
```

### **3. Integration with Existing Facts System**

The variables system extends the existing facts system (see [Facts System](../facts-system.md)) to support variable resolution and template integration:

```go
// Extend the existing facts system to support variable resolution
type FactsCollector struct {
    // ... existing fields
    variableResolver *VariableResolver
}

func (fc *FactsCollector) CollectFacts(machine string) (map[string]interface{}, error) {
    facts := make(map[string]interface{})
    
    // Collect system facts as before
    systemFacts, err := fc.collectSystemFacts(machine)
    if err != nil {
        return nil, err
    }
    
    // Add facts to variable context for seamless access
    fc.variableResolver.context.Variables.FactsVars = systemFacts
    fc.variableResolver.context.Facts = systemFacts // Maintain backward compatibility
    
    // Now facts can be referenced in variable interpolation
    // e.g., ${var.hostname} will resolve to the machine's hostname from facts
    
    return facts, nil
}
```

### **4. Benefits of This Approach**

1. **Automatic Integration**: Facts are automatically available as variables, no manual loading required
2. **Consistent Access**: Facts can be accessed using the same `${var.fact_name}` syntax as other variables
3. **Dynamic Resolution**: Facts are resolved at runtime for each machine
4. **Template Compatibility**: Facts work seamlessly in templates alongside other variables
5. **Performance**: Caching prevents repeated fact lookups
6. **Debugging Support**: Source tracking helps with troubleshooting
7. **Clear Precedence**: Environment variables can override project settings for security

## Variable Interpolation

### **1. Interpolation Syntax**

```hcl
# In any HCL file
machine "web-server" {
  hostname = "${var.app_name}-web-${var.environment}"
  port     = var.server_port
  tags = {
    environment = var.environment
    app         = var.app_name
  }
}

action "deploy" {
  command = "echo 'Deploying ${var.app_name} to ${var.environment}'"
  depends_on = ["${var.database_action}"]
}
```

### **2. Interpolation Context**

```go
type InterpolationContext struct {
    Variables map[string]interface{}
    Facts     map[string]interface{}
    Machines  map[string]interface{}
    Actions   map[string]interface{}
    Tags      map[string]string
    Env       map[string]string
}

func (ic *InterpolationContext) Interpolate(input string) (string, error) {
    // Replace ${var.name} with actual values
    // Handle nested interpolation
    // Validate variable existence
    // Type checking for interpolation
}
```

### **3. Template Context Enhancement**

The existing `TemplateContext` is enhanced to include the new `VariableContext` for seamless variable access:

```go
// Enhanced TemplateContext with VariableContext integration
type TemplateContext struct {
    // Existing fields...
    Facts     map[string]interface{}
    Machines  map[string]interface{}
    Actions   map[string]interface{}
    Tags      map[string]string
    Env       map[string]string
    
    // Enhanced variable context
    Variables *VariableContext
    VarFuncs  map[string]interface{}
}

// Template functions for variable access with precedence resolution
func (tc *TemplateContext) Var(name string) interface{} {
    if tc.Variables == nil {
        return nil
    }
    
    // Use the variable resolver for proper precedence handling
    resolver := &VariableResolver{context: tc}
    if value, err := resolver.Resolve(name); err == nil {
        return value
    }
    return nil
}

func (tc *TemplateContext) VarOrDefault(name string, defaultValue interface{}) interface{} {
    if value := tc.Var(name); value != nil {
        return value
    }
    return defaultValue
}

func (tc *TemplateContext) VarRequired(name string) (interface{}, error) {
    if value := tc.Var(name); value != nil {
        return value, nil
    }
    return nil, fmt.Errorf("required variable '%s' not found", name)
}

// Enhanced template functions for facts access (maintains backward compatibility)
func (tc *TemplateContext) Fact(name string) interface{} {
    if value, exists := tc.Facts[name]; exists {
        return value
    }
    return nil
}

// New function for accessing variable source information
func (tc *TemplateContext) VarSource(name string) *VariableSource {
    if tc.Variables != nil {
        if source, exists := tc.Variables.SourceMap[name]; exists {
            return &source
        }
    }
    return nil
}
```

### **4. Template Syntax for Variables**

With the enhanced VariableContext, templates can seamlessly access variables and facts:

```hcl
# Template file: templates/nginx.conf.tmpl
server {
    # Variables with precedence resolution (env → project → facts → defaults)
    listen {{var "server_port" | default 80}};
    server_name {{var "domain_name"}};
    
    # Using variable functions with fallbacks
    root {{varOrDefault "web_root" "/var/www/html"}};
    
    # Facts are automatically available as variables
    # {{var "hostname"}} resolves to the machine's hostname from facts
    access_log /var/log/nginx/{{var "hostname"}}.access.log;
    
    # Required variables with error handling
    {{if var "ssl_enabled"}}
    ssl_certificate {{varRequired "ssl_cert_path"}};
    ssl_key {{varRequired "ssl_key_path"}};
    {{end}}
    
    # Variable interpolation in templates
    location / {
        proxy_pass http://{{var "app_host"}}:{{var "app_port"}};
        proxy_set_header Host {{var "domain_name"}};
    }
    
    # Facts can be accessed directly (backward compatibility)
    # {{fact "os_family"}} - same as {{var "os_family"}}
    # {{fact "ipaddress"}} - same as {{var "ipaddress"}}
    
    # Variable source information for debugging
    {{if varSource "database_url"}}
    # Database URL from: {{varSource "database_url".Type}} ({{varSource "database_url".File}}:{{varSource "database_url".Line}})
    {{end}}
}
```

**Variable Resolution Examples:**

```hcl
# variables.hcl - Project variables can reference facts
variable "log_path" {
  type = "string"
  default = "/var/log/${var.hostname}"  # References fact automatically
}

variable "config_dir" {
  type = "string"
  default = "/etc/${var.os_family}"     # References fact automatically
}

variable "server_name" {
  type = "string"
  default = "${var.app_name}-${var.hostname}"  # References project var + fact
}

# Environment variables can override project variables
# SPOOKY_DATABASE_URL will override any database_url in variables.hcl

# Facts are automatically available
# hostname, os_family, ipaddress, etc. are available as variables
```

## Schema Integration

### **1. Leveraging Existing Schema Infrastructure**

The variables system integrates with the existing schema validation infrastructure to provide comprehensive validation and export capabilities.

```go
// Use existing SchemaValidator infrastructure
func runVariablesValidate(cmd *cobra.Command, args []string) error {
    projectDir := args[0]
    logger := logging.GetLogger()
    
    // Create schema validator (reuse existing infrastructure)
    validator := schemas.NewSchemaValidator()
    
    // Load variables schema specifically
    if err := validator.LoadSchema(schemas.SchemaTypeVariables); err != nil {
        return fmt.Errorf("failed to load variables schema: %w", err)
    }
    
    // Validate variables using schema
    result := validator.ValidateVariables(projectDir)
    if !result.Valid {
        fmt.Println("❌ Variables validation failed:")
        for _, err := range result.Errors {
            fmt.Printf("  - %s: %s\n", err.Field, err.Message)
        }
        return fmt.Errorf("variables validation failed")
    }
    
    fmt.Println("✅ Variables validation passed")
    return nil
}

// Extend SchemaValidator with variables validation
func (sv *SchemaValidator) ValidateVariables(projectPath string) *ValidationResult {
    result := &ValidationResult{Valid: true}
    
    // Load variables schema
    variablesSchema, exists := sv.schemas[SchemaTypeVariables]
    if !exists {
        result.Valid = false
        result.Errors = append(result.Errors, ValidationError{
            Field:   "schema",
            Message: "variables schema not loaded",
        })
        return result
    }
    
    // Validate variables.hcl if exists
    variablesFile := filepath.Join(projectPath, "variables.hcl")
    if _, err := os.Stat(variablesFile); err == nil {
        if err := sv.validateVariablesFile(variablesFile, variablesSchema); err != nil {
            result.Valid = false
            result.Errors = append(result.Errors, ValidationError{
                Field:   "variables.hcl",
                Message: err.Error(),
            })
        }
    }
    
    // Validate variables/*.hcl files
    variablesDir := filepath.Join(projectPath, "variables")
    if _, err := os.Stat(variablesDir); err == nil {
        entries, err := os.ReadDir(variablesDir)
        if err != nil {
            result.Valid = false
            result.Errors = append(result.Errors, ValidationError{
                Field:   "variables/",
                Message: fmt.Sprintf("failed to read variables directory: %v", err),
            })
            return result
        }
        
        for _, entry := range entries {
            if !entry.IsDir() && filepath.Ext(entry.Name()) == ".hcl" {
                filePath := filepath.Join(variablesDir, entry.Name())
                if err := sv.validateVariablesFile(filePath, variablesSchema); err != nil {
                    result.Valid = false
                    result.Errors = append(result.Errors, ValidationError{
                        Field:   entry.Name(),
                        Message: err.Error(),
                    })
                }
            }
        }
    }
    
    return result
}
```

### **2. Schema-Based Validation Rules**

The variables schema (`internal/schemas/schemas/variables.hcl`) provides comprehensive validation rules:

```hcl
# Schema validation rules from variables.hcl
validation = {
  # File location validation
  file_location = {
    rule = "path"
    pattern = "^(variables\\.hcl|variables/[^/]+\\.hcl)$"
    message = "Variable files must be variables.hcl in project root or .hcl files in variables/ directory"
  }
  
  # Variable name validation
  variable_name = {
    rule = "regex"
    pattern = "^[a-z][a-z0-9_]*$"
    message = "Variable names must be lowercase with underscores, starting with a letter"
  }
  
  # Type validation
  variable_type = {
    rule = "enum"
    allowed_values = ["string", "number", "bool", "list", "map"]
    message = "Variable type must be one of: string, number, bool, list, map"
  }
  
  # Required variable validation
  required_variable = {
    rule = "required"
    condition = "variable_is_required"
    message = "Required variables must have a default value or be provided via environment variable"
  }
  
  # Validation condition syntax
  validation_condition = {
    rule = "hcl"
    message = "Validation condition must be valid HCL syntax"
  }
  
  # No circular dependencies
  no_circular_deps = {
    rule = "acyclic"
    message = "Variables cannot have circular dependencies"
  }
}
```

### **3. Schema-Enhanced Export**

Variables export includes schema validation status and metadata:

```go
// Export with schema validation status
func runVariablesExport(cmd *cobra.Command, args []string) error {
    projectDir := args[0]
    format, _ := cmd.Flags().GetString("format")
    output, _ := cmd.Flags().GetString("output")
    resolve, _ := cmd.Flags().GetBool("resolve")
    
    // Load variables with schema validation
    validator := schemas.NewSchemaValidator()
    if err := validator.LoadSchema(schemas.SchemaTypeVariables); err != nil {
        return fmt.Errorf("failed to load variables schema: %w", err)
    }
    
    // Load and validate variables
    variablesConfig, err := config.LoadVariablesConfig(projectDir)
    if err != nil {
        return fmt.Errorf("failed to load variables: %w", err)
    }
    
    // Create export structure with schema metadata
    export := &VariablesExport{
        Variables: variablesConfig.Variables,
        Schema: &SchemaMetadata{
            Version: "1.0",
            Type:    "variables",
            Rules:   validator.GetSchemaRules(schemas.SchemaTypeVariables),
        },
        Validation: validator.ValidateVariables(projectDir),
        Metadata: &ExportMetadata{
            ExportedAt: time.Now(),
            Project:    projectDir,
            Format:     format,
            Resolved:   resolve,
        },
    }
    
    // Export based on format
    switch format {
    case "json":
        return exportToJSON(export, output)
    case "hcl":
        return exportToHCL(export, output)
    default:
        return fmt.Errorf("unsupported export format: %s", format)
    }
}

type VariablesExport struct {
    Variables  []Variable       `json:"variables"`
    Schema     *SchemaMetadata  `json:"schema"`
    Validation *ValidationResult `json:"validation"`
    Metadata   *ExportMetadata  `json:"metadata"`
}

type SchemaMetadata struct {
    Version string                 `json:"version"`
    Type    string                 `json:"type"`
    Rules   map[string]interface{} `json:"rules"`
}

type ExportMetadata struct {
    ExportedAt time.Time `json:"exported_at"`
    Project    string    `json:"project"`
    Format     string    `json:"format"`
    Resolved   bool      `json:"resolved"`
}
```

### **4. Schema Integration Benefits**

1. **Consistent Validation**: Uses the same validation infrastructure as other spooky components
2. **Comprehensive Rules**: Leverages all schema validation rules (file location, names, types, dependencies)
3. **Embedded Schemas**: No external schema files needed - schemas are embedded in the binary
4. **Extensible**: Easy to add new validation rules to the schema
5. **Export Metadata**: Export includes schema information for tooling integration
6. **Error Consistency**: Uses the same error format as other validation commands
7. **Performance**: Reuses existing schema loading and validation infrastructure

### **5. Schema Usage Examples**

```bash
# Validate using embedded schema
$ spooky variables validate ./my-project
Validating variables using embedded schema...
✓ variables.hcl: 5 variables validated against schema
✓ variables/environment.hcl: 3 variables validated against schema
✅ Variables validation passed (8 variables in 2 files)

# Export with schema metadata
$ spooky variables export ./my-project --format json --output variables.json
Exported 8 variables with schema validation metadata to variables.json

# Show schema validation results
$ spooky variables show app_name ./my-project --verbose
Variable: app_name
  Type: string
  Default: "myapp"
  Schema Validation: ✓ Valid
  Schema Rules Applied:
    - variable_name: ✓ Matches pattern ^[a-z][a-z0-9_]*$
    - variable_type: ✓ Valid type "string"
    - file_location: ✓ Valid location variables.hcl
```

## File Location and Validation Requirements

### **1. Strict File Location Rules**

```go
// Only these locations are allowed for variable files:
// - project/variables.hcl (project root only)
// - project/variables/*.hcl (variables directory only)

func (vm *VariableMerger) validateProjectPath(projectPath string) error {
    absPath, err := filepath.Abs(projectPath)
    if err != nil {
        return fmt.Errorf("invalid project path: %w", err)
    }
    
    if !dirExists(absPath) {
        return fmt.Errorf("project directory does not exist: %s", absPath)
    }
    
    vm.projectPath = absPath
    return nil
}

func (vm *VariableMerger) validateFilePath(filePath string) error {
    absPath, err := filepath.Abs(filePath)
    if err != nil {
        return fmt.Errorf("invalid file path: %w", err)
    }
    
    // Must be within project directory
    if !vm.isWithinProject(absPath) {
        return fmt.Errorf("file must be within project directory: %s", filePath)
    }
    
    // Must be variables.hcl or in variables/ directory
    if !vm.isValidVariableFile(absPath) {
        return fmt.Errorf("invalid variable file location: %s", filePath)
    }
    
    return nil
}

func (vm *VariableMerger) isValidVariableFile(filePath string) bool {
    fileName := filepath.Base(filePath)
    dirName := filepath.Base(filepath.Dir(filePath))
    
    // Allow variables.hcl in project root
    if fileName == "variables.hcl" && dirName == filepath.Base(vm.projectPath) {
        return true
    }
    
    // Allow .hcl files in variables/ directory
    if strings.HasSuffix(fileName, ".hcl") && dirName == "variables" {
        return true
    }
    
    return false
}
```

### **2. Strict HCL Validation**

```go
func (vm *VariableMerger) processFile(filePath string, variables map[string]*Variable) error {
    // Validate file location first
    if err := vm.validateFilePath(filePath); err != nil {
        return err
    }
    
    // Read file content
    content, err := os.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("error reading file %s: %w", filePath, err)
    }
    
    // Parse HCL with strict validation
    var hclContent interface{}
    if err := hcl.Unmarshal(content, &hclContent); err != nil {
        return fmt.Errorf("invalid HCL in %s: %w", filepath.Base(filePath), err)
    }
    
    // Parse as variables with strict schema validation
    var varBlocks []VariableBlock
    if err := hcl.Unmarshal(content, &varBlocks); err != nil {
        return fmt.Errorf("invalid variable definition in %s: %w", filepath.Base(filePath), err)
    }
    
    // Validate each variable block
    for i, varBlock := range varBlocks {
        if err := vm.validateVariableBlock(varBlock, filePath, i+1); err != nil {
            return fmt.Errorf("variable %d in %s: %w", i+1, filepath.Base(filePath), err)
        }
        
        // Add to variables map
        variables[varBlock.Name] = &Variable{
            Name:       varBlock.Name,
            Type:       varBlock.Type,
            Value:      varBlock.Default,
            Required:   varBlock.Required,
            Sensitive:  varBlock.Sensitive,
            SourceFile: filePath,
            Line:       varBlock.Line,
        }
    }
    
    return nil
}
```

## Implementation Phases

### **Phase 1: Core Variables System (Week 1-2)**

1. **Variable Definition Schema**
   - Create Go structs for variable definitions
   - Implement HCL parsing for variables
   - Add validation rules and type checking

2. **Integration with Existing Config Parsing**
   - Extend existing parsing functions (`ParseProjectConfig`, `ParseInventoryConfig`, `ParseActionsConfig`) to include variable parsing
   - Add `ParseVariablesConfig()` and `LoadVariablesConfig()` functions following the same pattern as `LoadActionsConfig`
   - Integrate variable loading into `TemplateContext.loadVariables()` method
   - Use existing file discovery patterns for `variables.hcl` and `variables/*.hcl`
   - Leverage existing HCL parsing infrastructure (`hclparse.NewParser()`, `gohcl.DecodeBody()`)

3. **Parse-Time File Merging**
   - Implement `VariableMerger` with conflict detection
   - Load and merge all variable files at parse time
   - Validate file locations (only project directory files)
   - Implement strict mode for validation (treat conflicts as errors)
   - Support both memory and temporary file merging strategies

4. **Strict File Location Validation**
   - Implement project path validation
   - Restrict file reading to project directory only
   - Validate file locations (variables.hcl and variables/*.hcl only)

5. **Strict HCL Validation**
   - Implement HCL syntax validation for all variable files
   - Treat invalid HCL as errors, not warnings
   - Provide clear error messages with file and line information
   - **Use embedded variables schema for comprehensive validation**
   - **Validate against schema rules: file location, variable names, types, validation conditions**
   - **Leverage existing SchemaValidator infrastructure**

6. **Variable Storage**
   - Create `VariableContext` with hierarchical storage
   - Implement variable resolution with precedence
   - Add to project context with source tracking

7. **CLI Command: validate-variables**
   - Implement `spooky variables validate` command (see [CLI System Design](../cli-system.md#5-variables-noun))
   - **Use embedded variables schema for validation**
   - **Leverage existing SchemaValidator infrastructure**
   - Validate all variable files in project
   - Report file location and HCL validation issues
   - Detect and report variable conflicts with detailed file/line information
   - **Validate against schema rules: variable names, types, validation conditions, circular dependencies**
   - Provide detailed error messages with file and line numbers
   - Exit with non-zero code on validation failures

8. **CLI Command: list-variables**
   - Implement `spooky variables list` command (see [CLI System Design](../cli-system.md#5-variables-noun))
   - List all variables in project with filtering options
   - Support filtering by type, file, required status, and sensitivity
   - Provide verbose output with detailed variable information
   - Support JSON/HCL output formats
   - **Include schema validation status for each variable**

9. **CLI Command: show-variable**
   - Implement `spooky variables show <variable_name>` command (see [CLI System Design](../cli-system.md#5-variables-noun))
   - Display detailed information about specific variables
   - Show variable dependencies and usage across project
   - Display validation rules and metadata
   - Support verbose output with full context
   - **Show schema validation results for the variable**

10. **CLI Command: export-variables**
    - Implement `spooky variables export` command (see [CLI System Design](../cli-system.md#5-variables-noun))
    - Export variables to JSON or HCL format
    - Support filtering by type, sensitivity, and specific variables
    - Include resolved values and metadata
    - Support output to files with proper formatting
    - **Export with schema validation status**
    - **Include schema metadata in export output**

### **Phase 2: Dependency Management (Week 3)**

1. **Dependency Graph**
   - Build dependency tracking system
   - Implement circular reference detection
   - Add topological sorting for resolution

2. **Validation System**
   - Validate variable dependencies
   - Check for undefined variables
   - Type validation for dependencies

3. **Error Handling**
   - Clear error messages for circular refs
   - Dependency chain visualization
   - File and line number reporting

### **Phase 3: Interpolation and Integration (Week 4)**

1. **Variable Interpolation**
   - Implement `${var.name}` syntax
   - Add to HCL parsing pipeline
   - Handle nested interpolation

2. **Actions Integration**
   - Apply same merging to actions
   - Add dependency management to actions
   - Cross-file action references

3. **Template Integration**
   - Add variables to template context
   - Update template functions
   - Variable access in templates

### **Phase 4: Advanced Features (Week 5-6)**

1. **File Size Optimization**
   - Implement temporary file merging
   - Add size thresholds and limits
   - Memory-efficient processing

2. **Advanced Validation**
   - Custom validation functions
   - Conditional validation
   - Cross-variable validation

3. **Documentation and Testing**
   - Comprehensive test coverage
   - User documentation
   - Migration guide

## CLI Command Implementation

### **Integration with CLI System**

The variables system integrates with the existing CLI system as described in [CLI System Design](../cli-system.md#5-variables-noun). All variables commands follow the established `spooky noun verb` pattern and integrate with the global configuration and logging systems.

### **Command Implementation Requirements**

#### **1. spooky variables validate**
```go
// Command structure
var variablesValidateCmd = &cobra.Command{
    Use:   "validate <project directory>",
    Short: "Validate project variables",
    Long:  `Validate all variable files in a project for syntax, types, and dependencies.`,
    Args:  cobra.ExactArgs(1),
    RunE:  runVariablesValidate,
}

// Implementation requirements
func runVariablesValidate(cmd *cobra.Command, args []string) error {
    // 1. Validate project directory exists
    // 2. Load and parse variables.hcl (if exists)
    // 3. Load and parse all variables/*.hcl files
    // 4. Validate HCL syntax for all files
    // 5. Validate variable definitions against schema
    // 6. Check for circular dependencies
    // 7. Validate custom validation conditions
    // 8. Report results with appropriate exit codes
}
```

#### **2. spooky variables list**
```go
// Command structure
var variablesListCmd = &cobra.Command{
    Use:   "list <project directory>",
    Short: "List project variables",
    Long:  `List all variables in a project with optional filtering.`,
    Args:  cobra.ExactArgs(1),
    RunE:  runVariablesList,
}

// Implementation requirements
func runVariablesList(cmd *cobra.Command, args []string) error {
    // 1. Load all variable files
    // 2. Apply filters (type, file, required, sensitive)
    // 3. Format output based on --format flag
    // 4. Support --verbose for detailed output
    // 5. Support --output for file output
}
```

#### **3. spooky variables show**
```go
// Command structure
var variablesShowCmd = &cobra.Command{
    Use:   "show <variable_name> <project directory>",
    Short: "Show variable details",
    Long:  `Show detailed information about a specific variable.`,
    Args:  cobra.ExactArgs(2),
    RunE:  runVariablesShow,
}

// Implementation requirements
func runVariablesShow(cmd *cobra.Command, args []string) error {
    // 1. Load all variable files
    // 2. Find specified variable
    // 3. Show variable metadata and definition
    // 4. Show dependencies (if --dependencies flag)
    // 5. Show usage across project (if --usage flag)
    // 6. Show validation rules (if --validation flag)
}
```

#### **4. spooky variables export**
```go
// Command structure
var variablesExportCmd = &cobra.Command{
    Use:   "export <project directory>",
    Short: "Export variables",
    Long:  `Export project variables to JSON or HCL format.`,
    Args:  cobra.ExactArgs(1),
    RunE:  runVariablesExport,
}

// Implementation requirements
func runVariablesExport(cmd *cobra.Command, args []string) error {
    // 1. Load all variable files
    // 2. Apply filters (type, sensitivity, specific variables)
    // 3. Resolve values (if --resolve flag)
    // 4. Format output (JSON or HCL)
    // 5. Write to file or stdout
    // 6. Handle sensitive variable export appropriately
}
```

### **Global Flag Integration**

All variables commands integrate with the global CLI system:

```go
// Global flags supported by all variables commands
variablesValidateCmd.Flags().Bool("verbose", false, "Show detailed output")
variablesValidateCmd.Flags().String("format", "text", "Output format (text, json)")
variablesValidateCmd.Flags().String("output", "", "Output file path")

variablesListCmd.Flags().Bool("verbose", false, "Show detailed output")
variablesListCmd.Flags().String("type", "", "Filter by variable type")
variablesListCmd.Flags().String("file", "", "Filter by source file")
variablesListCmd.Flags().Bool("required", false, "Show only required variables")
variablesListCmd.Flags().Bool("sensitive", false, "Show only sensitive variables")
variablesListCmd.Flags().String("filter", "", "Filter using complex expressions")
variablesListCmd.Flags().String("format", "text", "Output format (text, json, hcl)")
variablesListCmd.Flags().String("output", "", "Output file path")

variablesShowCmd.Flags().Bool("verbose", false, "Show detailed output")
variablesShowCmd.Flags().Bool("dependencies", false, "Show variable dependencies")
variablesShowCmd.Flags().Bool("usage", false, "Show variable usage across project")
variablesShowCmd.Flags().Bool("validation", false, "Show validation rules")

variablesExportCmd.Flags().String("format", "", "Export format (json, hcl) - required")
variablesExportCmd.Flags().String("output", "", "Output file path - required")
variablesExportCmd.Flags().String("variables", "", "Export specific variables")
variablesExportCmd.Flags().String("type", "", "Export by variable type")
variablesExportCmd.Flags().Bool("exclude-sensitive", false, "Exclude sensitive variables")
variablesExportCmd.Flags().Bool("resolve", false, "Include resolved values")
```

### **Error Handling and Exit Codes**

All variables commands follow consistent error handling:

```go
// Exit codes
const (
    ExitSuccess = 0
    ExitValidationError = 1
    ExitFileError = 2
    ExitSchemaError = 3
    ExitDependencyError = 4
)

// Error types
type ValidationError struct {
    File    string
    Line    int
    Message string
}

type DependencyError struct {
    Variable string
    Cycle    []string
    Message  string
}
```

### **Integration with Template System**

Variables commands integrate with the enhanced template system:

```go
// Enhanced TemplateContext with VariableContext integration
type TemplateContext struct {
    // Existing fields...
    Facts     map[string]interface{}
    Machines  map[string]interface{}
    Actions   map[string]interface{}
    Tags      map[string]string
    Env       map[string]string
    
    // Enhanced variable context
    Variables *VariableContext
    VarFuncs  map[string]interface{}
}

// Template functions for variable access with precedence resolution
func (tc *TemplateContext) Var(name string) interface{} {
    if tc.Variables == nil {
        return nil
    }
    
    // Use the variable resolver for proper precedence handling
    resolver := &VariableResolver{context: tc}
    if value, err := resolver.Resolve(name); err == nil {
        return value
    }
    return nil
}

func (tc *TemplateContext) VarOrDefault(name string, defaultValue interface{}) interface{} {
    if value := tc.Var(name); value != nil {
        return value
    }
    return defaultValue
}

func (tc *TemplateContext) VarRequired(name string) (interface{}, error) {
    if value := tc.Var(name); value != nil {
        return value, nil
    }
    return nil, fmt.Errorf("required variable '%s' not found", name)
}

// Enhanced template functions for facts access (maintains backward compatibility)
func (tc *TemplateContext) Fact(name string) interface{} {
    if value, exists := tc.Facts[name]; exists {
        return value
    }
    return nil
}

// New function for accessing variable source information
func (tc *TemplateContext) VarSource(name string) *VariableSource {
    if tc.Variables != nil {
        if source, exists := tc.Variables.SourceMap[name]; exists {
            return &source
        }
    }
    return nil
}
```

## Benefits

### **Immediate Benefits**
- **Centralized configuration** with variables
- **Reusable values** across project files
- **Environment-specific** configurations
- **Sensitive data handling** with proper validation
- **Dependency management** prevents circular references

### **Long-term Benefits**
- **Scalable configuration** for large projects
- **Team collaboration** with organized variable files
- **CI/CD integration** with environment variables
- **Configuration reuse** across multiple projects
- **Better error handling** with dependency tracking

## Success Metrics

### **Functionality Metrics**
- [ ] Variables can be defined in `variables.hcl` and `variables/` directory
- [ ] File merging works for both small and large files
- [ ] Variable conflicts are detected and reported as validation errors
- [ ] File merging is integrated into existing config parsing system
- [ ] Parse-time merging provides early error detection
- [ ] Circular dependencies are detected and reported
- [ ] Variable interpolation works in all HCL files
- [ ] Actions support the same merging and dependency system
- [ ] Variable types are properly validated
- [ ] Only project directory variable files are read
- [ ] Invalid HCL in variable files is treated as an error
- [ ] Clear error messages for file location and HCL validation issues
- [ ] **Variables schema is used for comprehensive validation**
- [ ] **Schema validation covers file location, variable names, types, and validation conditions**
- [ ] **Schema validation detects circular dependencies**
- [ ] **Export includes schema validation status and metadata**
- [ ] `spooky variables validate` command validates all variable files (see [CLI System Design](../cli-system.md#5-variables-noun))
- [ ] `spooky variables validate` detects and reports variable conflicts with file/line details
- [ ] `spooky variables validate` uses embedded variables schema for validation
- [ ] `spooky variables list` command lists variables with filtering options (see [CLI System Design](../cli-system.md#5-variables-noun))
- [ ] `spooky variables list` includes schema validation status for each variable
- [ ] `spooky variables show` command displays detailed variable information (see [CLI System Design](../cli-system.md#5-variables-noun))
- [ ] `spooky variables show` displays schema validation results for variables
- [ ] `spooky variables export` command exports variables to JSON/HCL format (see [CLI System Design](../cli-system.md#5-variables-noun))
- [ ] `spooky variables export` includes schema metadata in export output
- [ ] Command provides detailed error reporting with file and line numbers
- [ ] Command exits with appropriate exit codes (0 for success, non-zero for errors)

### **Performance Metrics**
- [ ] File merging handles files up to 2GB efficiently
- [ ] Memory usage stays reasonable for large projects
- [ ] Dependency resolution is fast for complex graphs
- [ ] Interpolation performance doesn't impact parsing
- [ ] Variable loading handles files efficiently

### **User Experience Metrics**
- [ ] Clear error messages for dependency issues
- [ ] Intuitive variable definition syntax
- [ ] Good documentation and examples
- [ ] Backward compatibility with existing projects
- [ ] `spooky variables validate` provides clear, actionable error messages (see [CLI System Design](../cli-system.md#5-variables-noun))
- [ ] `spooky variables list` provides clear, organized variable listings (see [CLI System Design](../cli-system.md#5-variables-noun))
- [ ] `spooky variables show` provides comprehensive variable details (see [CLI System Design](../cli-system.md#5-variables-noun))
- [ ] `spooky variables export` provides flexible export options (see [CLI System Design](../cli-system.md#5-variables-noun))
- [ ] Command can be used in CI/CD pipelines for validation
- [ ] Success messages show summary of validated variables and files 

## Risk Assessment

### **Technical Risks**
- **Large file handling** - Mitigation: Implement size thresholds and temporary files
- **Circular dependency complexity** - Mitigation: Clear error reporting and visualization
- **Performance impact** - Mitigation: Efficient algorithms and caching

### **User Experience Risks**
- **Learning curve** - Mitigation: Good documentation and examples
- **Migration complexity** - Mitigation: Backward compatibility and migration tools
- **Error debugging** - Mitigation: Clear error messages and dependency chains

### **Implementation Risks**
- **Scope creep** - Mitigation: Focus on core functionality first
- **Testing complexity** - Mitigation: Comprehensive test scenarios
- **Integration issues** - Mitigation: Incremental implementation and testing 

## CLI Integration

### **1. Variables Commands Implementation**

Variables commands should be implemented in the existing `internal/cli/` structure following the same pattern as actions/facts/machines commands.

```go
// internal/cli/variables.go
package cli

import (
    "fmt"
    "path/filepath"
    
    "github.com/spf13/cobra"
    "spooky/internal/config"
    "spooky/internal/logging"
    "spooky/internal/schemas"
)

// Variables commands following the same pattern as other CLI commands
var variablesCmd = &cobra.Command{
    Use:   "variables",
    Short: "Manage project variables",
    Long:  `Manage project variables including validation, listing, and export.`,
}

var variablesValidateCmd = &cobra.Command{
    Use:   "validate <project directory>",
    Short: "Validate project variables",
    Long:  `Validate all variable files in a project for syntax, types, and dependencies.`,
    Args:  cobra.ExactArgs(1),
    RunE:  runVariablesValidate,
}

var variablesListCmd = &cobra.Command{
    Use:   "list <project directory>",
    Short: "List project variables",
    Long:  `List all variables in a project with optional filtering.`,
    Args:  cobra.ExactArgs(1),
    RunE:  runVariablesList,
}

var variablesShowCmd = &cobra.Command{
    Use:   "show <variable_name> <project directory>",
    Short: "Show variable details",
    Long:  `Show detailed information about a specific variable.`,
    Args:  cobra.ExactArgs(2),
    RunE:  runVariablesShow,
}

var variablesExportCmd = &cobra.Command{
    Use:   "export <project directory>",
    Short: "Export variables",
    Long:  `Export project variables to JSON or HCL format.`,
    Args:  cobra.ExactArgs(1),
    RunE:  runVariablesExport,
}

func init() {
    // Add variables commands to root command (same pattern as other nouns)
    rootCmd.AddCommand(variablesCmd)
    
    // Add subcommands to variables command
    variablesCmd.AddCommand(variablesValidateCmd)
    variablesCmd.AddCommand(variablesListCmd)
    variablesCmd.AddCommand(variablesShowCmd)
    variablesCmd.AddCommand(variablesExportCmd)
    
    // Add flags following the same pattern as other commands
    addVariablesFlags()
}

func addVariablesFlags() {
    // Validate flags
    variablesValidateCmd.Flags().Bool("verbose", false, "Show detailed output")
    variablesValidateCmd.Flags().String("format", "text", "Output format (text, json)")
    variablesValidateCmd.Flags().String("output", "", "Output file path")
    
    // List flags
    variablesListCmd.Flags().Bool("verbose", false, "Show detailed output")
    variablesListCmd.Flags().String("type", "", "Filter by variable type")
    variablesListCmd.Flags().String("file", "", "Filter by source file")
    variablesListCmd.Flags().Bool("required", false, "Show only required variables")
    variablesListCmd.Flags().Bool("sensitive", false, "Show only sensitive variables")
    variablesListCmd.Flags().String("filter", "", "Filter using complex expressions")
    variablesListCmd.Flags().String("format", "text", "Output format (text, json, hcl)")
    variablesListCmd.Flags().String("output", "", "Output file path")
    
    // Show flags
    variablesShowCmd.Flags().Bool("verbose", false, "Show detailed output")
    variablesShowCmd.Flags().Bool("dependencies", false, "Show variable dependencies")
    variablesShowCmd.Flags().Bool("usage", false, "Show variable usage across project")
    variablesShowCmd.Flags().Bool("validation", false, "Show validation rules")
    
    // Export flags
    variablesExportCmd.Flags().String("format", "", "Export format (json, hcl) - required")
    variablesExportCmd.Flags().String("output", "", "Output file path - required")
    variablesExportCmd.Flags().String("variables", "", "Export specific variables")
    variablesExportCmd.Flags().String("type", "", "Export by variable type")
    variablesExportCmd.Flags().Bool("exclude-sensitive", false, "Exclude sensitive variables")
    variablesExportCmd.Flags().Bool("resolve", false, "Include resolved values")
}
```

### **2. Integration with Existing Commands**

Variables declared in `variables.hcl` and `variables/*.hcl` are automatically available to all commands that need them, following the same pattern as facts and actions.

#### **Automatic Variables Loading**

```go
// internal/cli/template_context.go - Enhanced to load variables automatically
func (ctx *TemplateContext) Load(projectPath string) error {
    logger := logging.GetLogger()
    
    // Load existing context (following current pattern)
    if err := ctx.loadProjectConfig(logger, projectPath); err != nil {
        return err
    }
    
    if err := ctx.loadFacts(logger, projectPath); err != nil {
        return err
    }
    
    if err := ctx.loadActions(logger, projectPath); err != nil {
        return err
    }
    
    // NEW: Load variables automatically (same pattern as facts/actions)
    if err := ctx.loadVariables(logger, projectPath); err != nil {
        return err
    }
    
    if err := ctx.loadCustomData(logger, projectPath); err != nil {
        return err
    }
    
    ctx.loadEnvironment()
    return nil
}

// Variables are loaded automatically, no explicit flags needed
func (ctx *TemplateContext) loadVariables(logger logging.Logger, projectPath string) error {
    logger.Info("Loading variables from project", logging.String("project_path", projectPath))
    
    // Use the new LoadVariablesConfig function that can load from multiple sources
    variablesConfig, err := config.LoadVariablesConfig(projectPath)
    if err != nil {
        logger.Error("Failed to load variables", err, logging.String("project_path", projectPath))
        return fmt.Errorf("failed to load variables: %w", err)
    }
    
    // Convert to VariableContext
    ctx.Variables = &VariableContext{
        ProjectVars: make(map[string]interface{}),
        SourceMap:   make(map[string]VariableSource),
    }
    
    // Process variables and build source map
    for _, variable := range variablesConfig.Variables {
        ctx.Variables.ProjectVars[variable.Name] = variable.Default
        ctx.Variables.SourceMap[variable.Name] = VariableSource{
            Type: "project",
            File: variable.SourceFile,
            Line: variable.Line,
        }
    }
    
    logger.Info("Variables loaded successfully",
        logging.String("project_path", projectPath),
        logging.Int("variables", len(ctx.Variables.ProjectVars)))
    
    return nil
}
```

### **3. Commands That Use Variables**

Variables are available in all commands that need them, following these patterns:

#### **Commands That Always Use Variables:**

1. **Template Rendering Commands**
```bash
# Variables are automatically available in templates
$ spooky templates render ./my-project templates/nginx.conf.tmpl --output nginx.conf
# Template can use: {{var "app_name"}}, {{var "server_port"}}, etc.

# Variables are available in template context
$ spooky templates render ./my-project templates/config.tmpl --data data.json --output config.conf
# Both variables and data are available in template context
```

2. **Action Execution Commands**
```bash
# Variables are automatically available in action commands and scripts
$ spooky actions run ./my-project --action deploy
# Action can use: ${var.app_name}, ${var.environment}, etc. in commands

$ spooky actions run ./my-project --action setup --machine web-server
# Variables are available in action scripts and commands
```

3. **Facts Gathering Commands**
```bash
# Variables can be used in custom facts collection
$ spooky facts gather ./my-project
# Custom facts can reference variables for dynamic fact collection
```

#### **Commands That Optionally Use Variables:**

1. **Project Validation Commands**
```bash
# Variables are validated as part of project validation
$ spooky project validate ./my-project
# Includes variables validation along with other project components

# Explicit variables validation
$ spooky variables validate ./my-project
# Dedicated variables validation with detailed output
```

2. **Export Commands**
```bash
# Variables can be exported alongside other data
$ spooky machines export ./my-project --format json --output machines.json
# Export includes variable interpolation in machine configurations

$ spooky actions export ./my-project --format json --output actions.json
# Export includes variable interpolation in action configurations
```

#### **Commands That Don't Use Variables:**

1. **System Commands**
```bash
$ spooky --version
$ spooky --help
$ spooky completion generate
# These don't need variables
```

2. **Configuration Commands**
```bash
$ spooky config show
$ spooky config validate $XDG_CONFIG_HOME/spooky/spooky.hcl
# These work with global config, not project variables
```

### **4. Variable Availability Examples**

#### **Template Context Example:**
```go
// Variables are automatically available in all template contexts
func (ctx *TemplateContext) RenderTemplate(templatePath, outputPath string) error {
    // Variables are already loaded and available
    // Template can access:
    // - {{var "app_name"}} - from variables.hcl
    // - {{var "environment"}} - from variables/environment.hcl
    // - {{var "hostname"}} - from facts (automatic)
    // - {{env "SPOOKY_DEBUG"}} - from environment variables
    
    return ctx.templateEngine.Render(templatePath, outputPath, ctx)
}
```

#### **Action Context Example:**
```go
// Variables are automatically available in action running
func (ctx *TemplateContext) RunAction(action *config.Action, machine string) error {
    // Variables are already loaded and available
    // Action command can use:
    // - ${var.app_name} - from variables.hcl
    // - ${var.server_port} - from variables/environment.hcl
    // - ${var.hostname} - from facts (automatic)
    
    // Interpolate variables in action command
    interpolatedCommand := ctx.interpolateVariables(action.Command)
    
    return ctx.sshClient.RunCommand(machine, interpolatedCommand)
}
```

#### **Facts Context Example:**
```go
// Variables are available in facts collection
func (ctx *TemplateContext) CollectFacts(machine string) error {
    // Variables can be used in custom facts collection
    // Custom facts can reference variables for dynamic collection
    
    customFacts := map[string]interface{}{
        "app_name": ctx.Var("app_name"),
        "environment": ctx.Var("environment"),
        "log_path": ctx.Var("log_path"),
    }
    
    return ctx.factsCollector.CollectCustomFacts(machine, customFacts)
}
```

### **5. Variable Precedence in Commands**

Variables follow the same precedence order in all commands:

1. **Environment Variables** (highest priority)
   ```bash
   $ SPOOKY_APP_NAME=production spooky actions run ./my-project --action deploy
   # SPOOKY_APP_NAME overrides any app_name in variables.hcl
   ```

2. **Project Variables** (variables.hcl and variables/*.hcl)
   ```bash
   $ spooky templates render ./my-project templates/config.tmpl
   # Uses app_name from variables.hcl or variables/environment.hcl
   ```

3. **Facts** (automatic)
   ```bash
   $ spooky actions run ./my-project --action setup
   # hostname, os_family, etc. are automatically available as variables
   ```

4. **System Defaults** (lowest priority)
   ```bash
   $ spooky project init my-project
   # Uses system defaults for project creation
   ```

### **6. Error Handling for Missing Variables**

```go
// Commands handle missing variables gracefully
func (ctx *TemplateContext) Var(name string) interface{} {
    if ctx.Variables == nil {
        return nil
    }
    
    // Use the variable resolver for proper precedence handling
    resolver := &VariableResolver{context: ctx}
    if value, err := resolver.Resolve(name); err == nil {
        return value
    }
    
    // Return nil for missing variables (commands handle this)
    return nil
}

// Commands can check for required variables
func (ctx *TemplateContext) VarRequired(name string) (interface{}, error) {
    if value := ctx.Var(name); value != nil {
        return value, nil
    }
    return nil, fmt.Errorf("required variable '%s' not found", name)
}
```

### **7. CLI Integration Benefits**

1. **Automatic Loading**: Variables are loaded automatically with no explicit flags needed
2. **Consistent Access**: All commands use the same variable access patterns
3. **Precedence Handling**: Environment variables can override project settings
4. **Error Handling**: Graceful handling of missing variables
5. **Performance**: Variables are loaded once and cached for all commands
6. **Debugging**: Clear error messages for missing required variables
7. **Flexibility**: Commands can choose when to use variables

## File Location and Validation Requirements

### **1. Strict File Location Rules**

```go
// Only these locations are allowed for variable files:
// - project/variables.hcl (project root only)
// - project/variables/*.hcl (variables directory only)

func (vm *VariableMerger) validateProjectPath(projectPath string) error {
    absPath, err := filepath.Abs(projectPath)
    if err != nil {
        return fmt.Errorf("invalid project path: %w", err)
    }
    
    if !dirExists(absPath) {
        return fmt.Errorf("project directory does not exist: %s", absPath)
    }
    
    vm.projectPath = absPath
    return nil
}

func (vm *VariableMerger) validateFilePath(filePath string) error {
    absPath, err := filepath.Abs(filePath)
    if err != nil {
        return fmt.Errorf("invalid file path: %w", err)
    }
    
    // Must be within project directory
    if !vm.isWithinProject(absPath) {
        return fmt.Errorf("file must be within project directory: %s", filePath)
    }
    
    // Must be variables.hcl or in variables/ directory
    if !vm.isValidVariableFile(absPath) {
        return fmt.Errorf("invalid variable file location: %s", filePath)
    }
    
    return nil
}

func (vm *VariableMerger) isValidVariableFile(filePath string) bool {
    fileName := filepath.Base(filePath)
    dirName := filepath.Base(filepath.Dir(filePath))
    
    // Allow variables.hcl in project root
    if fileName == "variables.hcl" && dirName == filepath.Base(vm.projectPath) {
        return true
    }
    
    // Allow .hcl files in variables/ directory
    if strings.HasSuffix(fileName, ".hcl") && dirName == "variables" {
        return true
    }
    
    return false
}
```

### **2. Strict HCL Validation**

```go
func (vm *VariableMerger) processFile(filePath string, variables map[string]*Variable) error {
    // Validate file location first
    if err := vm.validateFilePath(filePath); err != nil {
        return err
    }
    
    // Read file content
    content, err := os.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("error reading file %s: %w", filePath, err)
    }
    
    // Parse HCL with strict validation
    var hclContent interface{}
    if err := hcl.Unmarshal(content, &hclContent); err != nil {
        return fmt.Errorf("invalid HCL in %s: %w", filepath.Base(filePath), err)
    }
    
    // Parse as variables with strict schema validation
    var varBlocks []VariableBlock
    if err := hcl.Unmarshal(content, &varBlocks); err != nil {
        return fmt.Errorf("invalid variable definition in %s: %w", filepath.Base(filePath), err)
    }
    
    // Validate each variable block
    for i, varBlock := range varBlocks {
        if err := vm.validateVariableBlock(varBlock, filePath, i+1); err != nil {
            return fmt.Errorf("variable %d in %s: %w", i+1, filepath.Base(filePath), err)
        }
        
        // Add to variables map
        variables[varBlock.Name] = &Variable{
            Name:       varBlock.Name,
            Type:       varBlock.Type,
            Value:      varBlock.Default,
            Required:   varBlock.Required,
            Sensitive:  varBlock.Sensitive,
            SourceFile: filePath,
            Line:       varBlock.Line,
        }
    }
    
    return nil
}
```

## Implementation Phases

### **Phase 1: Core Variables System (Week 1-2)**

1. **Variable Definition Schema**
   - Create Go structs for variable definitions
   - Implement HCL parsing for variables
   - Add validation rules and type checking

2. **Integration with Existing Config Parsing**
   - Extend existing parsing functions (`ParseProjectConfig`, `ParseInventoryConfig`, `ParseActionsConfig`) to include variable parsing
   - Add `ParseVariablesConfig()` and `LoadVariablesConfig()` functions following the same pattern as `LoadActionsConfig`
   - Integrate variable loading into `TemplateContext.loadVariables()` method
   - Use existing file discovery patterns for `variables.hcl` and `variables/*.hcl`
   - Leverage existing HCL parsing infrastructure (`hclparse.NewParser()`, `gohcl.DecodeBody()`)

3. **Parse-Time File Merging**
   - Implement `VariableMerger` with conflict detection
   - Load and merge all variable files at parse time
   - Validate file locations (only project directory files)
   - Implement strict mode for validation (treat conflicts as errors)
   - Support both memory and temporary file merging strategies

4. **Strict File Location Validation**
   - Implement project path validation
   - Restrict file reading to project directory only
   - Validate file locations (variables.hcl and variables/*.hcl only)

5. **Strict HCL Validation**
   - Implement HCL syntax validation for all variable files
   - Treat invalid HCL as errors, not warnings
   - Provide clear error messages with file and line information
   - **Use embedded variables schema for comprehensive validation**
   - **Validate against schema rules: file location, variable names, types, validation conditions**
   - **Leverage existing SchemaValidator infrastructure**

6. **Variable Storage**
   - Create `VariableContext` with hierarchical storage
   - Implement variable resolution with precedence
   - Add to project context with source tracking

7. **CLI Command: validate-variables**
   - Implement `spooky variables validate` command (see [CLI System Design](../cli-system.md#5-variables-noun))
   - **Use embedded variables schema for validation**
   - **Leverage existing SchemaValidator infrastructure**
   - Validate all variable files in project
   - Report file location and HCL validation issues
   - Detect and report variable conflicts with detailed file/line information
   - **Validate against schema rules: variable names, types, validation conditions, circular dependencies**
   - Provide detailed error messages with file and line numbers
   - Exit with non-zero code on validation failures

8. **CLI Command: list-variables**
   - Implement `spooky variables list` command (see [CLI System Design](../cli-system.md#5-variables-noun))
   - List all variables in project with filtering options
   - Support filtering by type, file, required status, and sensitivity
   - Provide verbose output with detailed variable information
   - Support JSON/HCL output formats
   - **Include schema validation status for each variable**

9. **CLI Command: show-variable**
   - Implement `spooky variables show <variable_name>` command (see [CLI System Design](../cli-system.md#5-variables-noun))
   - Display detailed information about specific variables
   - Show variable dependencies and usage across project
   - Display validation rules and metadata
   - Support verbose output with full context
   - **Show schema validation results for the variable**

10. **CLI Command: export-variables**
    - Implement `spooky variables export` command (see [CLI System Design](../cli-system.md#5-variables-noun))
    - Export variables to JSON or HCL format
    - Support filtering by type, sensitivity, and specific variables
    - Include resolved values and metadata
    - Support output to files with proper formatting
    - **Export with schema validation status**
    - **Include schema metadata in export output**

### **Phase 2: Dependency Management (Week 3)**

1. **Dependency Graph**
   - Build dependency tracking system
   - Implement circular reference detection
   - Add topological sorting for resolution

2. **Validation System**
   - Validate variable dependencies
   - Check for undefined variables
   - Type validation for dependencies

3. **Error Handling**
   - Clear error messages for circular refs
   - Dependency chain visualization
   - File and line number reporting

### **Phase 3: Interpolation and Integration (Week 4)**

1. **Variable Interpolation**
   - Implement `${var.name}` syntax
   - Add to HCL parsing pipeline
   - Handle nested interpolation

2. **Actions Integration**
   - Apply same merging to actions
   - Add dependency management to actions
   - Cross-file action references

3. **Template Integration**
   - Add variables to template context
   - Update template functions
   - Variable access in templates

### **Phase 4: Advanced Features (Week 5-6)**

1. **File Size Optimization**
   - Implement temporary file merging
   - Add size thresholds and limits
   - Memory-efficient processing

2. **Advanced Validation**
   - Custom validation functions
   - Conditional validation
   - Cross-variable validation

3. **Documentation and Testing**
   - Comprehensive test coverage
   - User documentation
   - Migration guide

## CLI Command Implementation

### **Integration with CLI System**

The variables system integrates with the existing CLI system as described in [CLI System Design](../cli-system.md#5-variables-noun). All variables commands follow the established `spooky noun verb` pattern and integrate with the global configuration and logging systems.

### **Command Implementation Requirements**

#### **1. spooky variables validate**
```go
// Command structure
var variablesValidateCmd = &cobra.Command{
    Use:   "validate <project directory>",
    Short: "Validate project variables",
    Long:  `Validate all variable files in a project for syntax, types, and dependencies.`,
    Args:  cobra.ExactArgs(1),
    RunE:  runVariablesValidate,
}

// Implementation requirements
func runVariablesValidate(cmd *cobra.Command, args []string) error {
    // 1. Validate project directory exists
    // 2. Load and parse variables.hcl (if exists)
    // 3. Load and parse all variables/*.hcl files
    // 4. Validate HCL syntax for all files
    // 5. Validate variable definitions against schema
    // 6. Check for circular dependencies
    // 7. Validate custom validation conditions
    // 8. Report results with appropriate exit codes
}
```

#### **2. spooky variables list**
```go
// Command structure
var variablesListCmd = &cobra.Command{
    Use:   "list <project directory>",
    Short: "List project variables",
    Long:  `List all variables in a project with optional filtering.`,
    Args:  cobra.ExactArgs(1),
    RunE:  runVariablesList,
}

// Implementation requirements
func runVariablesList(cmd *cobra.Command, args []string) error {
    // 1. Load all variable files
    // 2. Apply filters (type, file, required, sensitive)
    // 3. Format output based on --format flag
    // 4. Support --verbose for detailed output
    // 5. Support --output for file output
}
```

#### **3. spooky variables show**
```go
// Command structure
var variablesShowCmd = &cobra.Command{
    Use:   "show <variable_name> <project directory>",
    Short: "Show variable details",
    Long:  `Show detailed information about a specific variable.`,
    Args:  cobra.ExactArgs(2),
    RunE:  runVariablesShow,
}

// Implementation requirements
func runVariablesShow(cmd *cobra.Command, args []string) error {
    // 1. Load all variable files
    // 2. Find specified variable
    // 3. Show variable metadata and definition
    // 4. Show dependencies (if --dependencies flag)
    // 5. Show usage across project (if --usage flag)
    // 6. Show validation rules (if --validation flag)
}
```

#### **4. spooky variables export**
```go
// Command structure
var variablesExportCmd = &cobra.Command{
    Use:   "export <project directory>",
    Short: "Export variables",
    Long:  `Export project variables to JSON or HCL format.`,
    Args:  cobra.ExactArgs(1),
    RunE:  runVariablesExport,
}

// Implementation requirements
func runVariablesExport(cmd *cobra.Command, args []string) error {
    // 1. Load all variable files
    // 2. Apply filters (type, sensitivity, specific variables)
    // 3. Resolve values (if --resolve flag)
    // 4. Format output (JSON or HCL)
    // 5. Write to file or stdout
    // 6. Handle sensitive variable export appropriately
}
```

### **Global Flag Integration**

All variables commands integrate with the global CLI system:

```go
// Global flags supported by all variables commands
variablesValidateCmd.Flags().Bool("verbose", false, "Show detailed output")
variablesValidateCmd.Flags().String("format", "text", "Output format (text, json)")
variablesValidateCmd.Flags().String("output", "", "Output file path")

variablesListCmd.Flags().Bool("verbose", false, "Show detailed output")
variablesListCmd.Flags().String("type", "", "Filter by variable type")
variablesListCmd.Flags().String("file", "", "Filter by source file")
variablesListCmd.Flags().Bool("required", false, "Show only required variables")
variablesListCmd.Flags().Bool("sensitive", false, "Show only sensitive variables")
variablesListCmd.Flags().String("filter", "", "Filter using complex expressions")
variablesListCmd.Flags().String("format", "text", "Output format (text, json, hcl)")
variablesListCmd.Flags().String("output", "", "Output file path")

variablesShowCmd.Flags().Bool("verbose", false, "Show detailed output")
variablesShowCmd.Flags().Bool("dependencies", false, "Show variable dependencies")
variablesShowCmd.Flags().Bool("usage", false, "Show variable usage across project")
variablesShowCmd.Flags().Bool("validation", false, "Show validation rules")

variablesExportCmd.Flags().String("format", "", "Export format (json, hcl) - required")
variablesExportCmd.Flags().String("output", "", "Output file path - required")
variablesExportCmd.Flags().String("variables", "", "Export specific variables")
variablesExportCmd.Flags().String("type", "", "Export by variable type")
variablesExportCmd.Flags().Bool("exclude-sensitive", false, "Exclude sensitive variables")
variablesExportCmd.Flags().Bool("resolve", false, "Include resolved values")
```

### **Error Handling and Exit Codes**

All variables commands follow consistent error handling:

```go
// Exit codes
const (
    ExitSuccess = 0
    ExitValidationError = 1
    ExitFileError = 2
    ExitSchemaError = 3
    ExitDependencyError = 4
)

// Error types
type ValidationError struct {
    File    string
    Line    int
    Message string
}

type DependencyError struct {
    Variable string
    Cycle    []string
    Message  string
}
```

### **Integration with Template System**

Variables commands integrate with the enhanced template system:

```go
// Enhanced TemplateContext with VariableContext integration
type TemplateContext struct {
    // Existing fields...
    Facts     map[string]interface{}
    Machines  map[string]interface{}
    Actions   map[string]interface{}
    Tags      map[string]string
    Env       map[string]string
    
    // Enhanced variable context
    Variables *VariableContext
    VarFuncs  map[string]interface{}
}

// Template functions for variable access with precedence resolution
func (tc *TemplateContext) Var(name string) interface{} {
    if tc.Variables == nil {
        return nil
    }
    
    // Use the variable resolver for proper precedence handling
    resolver := &VariableResolver{context: tc}
    if value, err := resolver.Resolve(name); err == nil {
        return value
    }
    return nil
}

func (tc *TemplateContext) VarOrDefault(name string, defaultValue interface{}) interface{} {
    if value := tc.Var(name); value != nil {
        return value
    }
    return defaultValue
}

func (tc *TemplateContext) VarRequired(name string) (interface{}, error) {
    if value := tc.Var(name); value != nil {
        return value, nil
    }
    return nil, fmt.Errorf("required variable '%s' not found", name)
}

// Enhanced template functions for facts access (maintains backward compatibility)
func (tc *TemplateContext) Fact(name string) interface{} {
    if value, exists := tc.Facts[name]; exists {
        return value
    }
    return nil
}

// New function for accessing variable source information
func (tc *TemplateContext) VarSource(name string) *VariableSource {
    if tc.Variables != nil {
        if source, exists := tc.Variables.SourceMap[name]; exists {
            return &source
        }
    }
    return nil
}
```