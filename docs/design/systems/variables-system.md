# Variables System: Comprehensive Implementation Plan

## Overview

This document is the authoritative source for all variable system implementation details in spooky. It covers variable types, validation, interpolation, dependency management, file merging, and integration with the facts and template systems.

**Schema Integration**: This variables system leverages the schema validation infrastructure and patterns defined in [Schema System](../schema-system.md) for comprehensive variable validation, type checking, and schema-based export functionality.

**Architecture Integration**: Variables integrate with the overall spooky architecture as described in [Spooky Design](../spooky-design.md), providing configuration values to the actions system for template rendering and execution.

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

1. **Early Error Detection**: Conflicts and syntax errors caught before execution
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

// Runtime (during template rendering or action execution)
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

Following the patterns established by Ansible, Puppet, and Chef, facts should be automatically available as variables without explicit loading.

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

## Schema Updates

### **1. Project Directory Schema**

```hcl
# internal/schemas/project-directory.hcl
schema "project_directory" {
  # Required files
  file "project.hcl" {
    required = true
    validate = "hcl_project_config"
  }
  
  file "inventory.hcl" {
    required = true
    validate = "hcl_inventory_config"
  }
  
  file "actions.hcl" {
    required = false  # Can be in actions/ directory
    validate = "hcl_actions_config"
  }
  
  file "variables.hcl" {
    required = false  # Can be in variables/ directory
    validate = "hcl_variables_config"
  }
  
  # Required directories
  directory "facts.db" {
    required = true
    description = "Facts database directory"
    validate = "badgerdb_initialized"
  }
  
  # Optional directories
  directory "variables" {
    required = false
    description = "Organized variable files"
    validate = "hcl_variables_files"
  }
  
  directory "actions" {
    required = false
    description = "Organized action files"
    validate = "hcl_actions_files"
  }
  
  # ... other directories
}
```

### **2. Variables Configuration Schema**

```hcl
# internal/schemas/variables-config.hcl
schema "variables_config" {
  block "variable" {
    required = true
    min_blocks = 0
    
    attribute "name" {
      type = "string"
      required = true
      description = "Variable name"
      validation = "valid_variable_name"
    }
    
    attribute "type" {
      type = "string"
      required = true
      description = "Variable type"
      validation = "valid_variable_type"
      allowed_values = ["string", "number", "bool", "list", "map"]
    }
    
    attribute "description" {
      type = "string"
      required = false
      description = "Variable description"
    }
    
    attribute "default" {
      type = "any"
      required = false
      description = "Default value"
      validation = "type_matches_variable_type"
    }
    
    attribute "required" {
      type = "bool"
      required = false
      default = false
      description = "Whether variable is required"
    }
    
    attribute "sensitive" {
      type = "bool"
      required = false
      default = false
      description = "Whether variable contains sensitive data"
    }
    
    attribute "depends_on" {
      type = "list(string)"
      required = false
      description = "Dependencies on other variables"
    }
    
    block "validation" {
      required = false
      max_blocks = 1
      
      attribute "condition" {
        type = "string"
        required = true
        description = "Validation condition using HCL syntax"
        validation = "valid_hcl_validation_condition"
      }
      
      attribute "error_message" {
        type = "string"
        required = true
        description = "Error message for validation failure"
      }
    }
  }
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
$ spooky config validate ~/.config/spooky/spooky.hcl
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
// Variables are automatically available in action execution
func (ctx *TemplateContext) ExecuteAction(action *config.Action, machine string) error {
    // Variables are already loaded and available
    // Action command can use:
    // - ${var.app_name} - from variables.hcl
    // - ${var.server_port} - from variables/environment.hcl
    // - ${var.hostname} - from facts (automatic)
    
    // Interpolate variables in action command
    interpolatedCommand := ctx.interpolateVariables(action.Command)
    
    return ctx.sshClient.ExecuteCommand(machine, interpolatedCommand)
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

## Schema Updates

### **1. Project Directory Schema**

```hcl
# internal/schemas/project-directory.hcl
schema "project_directory" {
  # Required files
  file "project.hcl" {
    required = true
    validate = "hcl_project_config"
  }
  
  file "inventory.hcl" {
    required = true
    validate = "hcl_inventory_config"
  }
  
  file "actions.hcl" {
    required = false  # Can be in actions/ directory
    validate = "hcl_actions_config"
  }
  
  file "variables.hcl" {
    required = false  # Can be in variables/ directory
    validate = "hcl_variables_config"
  }
  
  # Required directories
  directory "facts.db" {
    required = true
    description = "Facts database directory"
    validate = "badgerdb_initialized"
  }
  
  # Optional directories
  directory "variables" {
    required = false
    description = "Organized variable files"
    validate = "hcl_variables_files"
  }
  
  directory "actions" {
    required = false
    description = "Organized action files"
    validate = "hcl_actions_files"
  }
  
  # ... other directories
}
```

### **2. Variables Configuration Schema**

```hcl
# internal/schemas/variables-config.hcl
schema "variables_config" {
  block "variable" {
    required = true
    min_blocks = 0
    
    attribute "name" {
      type = "string"
      required = true
      description = "Variable name"
      validation = "valid_variable_name"
    }
    
    attribute "type" {
      type = "string"
      required = true
      description = "Variable type"
      validation = "valid_variable_type"
      allowed_values = ["string", "number", "bool", "list", "map"]
    }
    
    attribute "description" {
      type = "string"
      required = false
      description = "Variable description"
    }
    
    attribute "default" {
      type = "any"
      required = false
      description = "Default value"
      validation = "type_matches_variable_type"
    }
    
    attribute "required" {
      type = "bool"
      required = false
      default = false
      description = "Whether variable is required"
    }
    
    attribute "sensitive" {
      type = "bool"
      required = false
      default = false
      description = "Whether variable contains sensitive data"
    }
    
    attribute "depends_on" {
      type = "list(string)"
      required = false
      description = "Dependencies on other variables"
    }
    
    block "validation" {
      required = false
      max_blocks = 1
      
      attribute "condition" {
        type = "string"
        required = true
        description = "Validation condition using HCL syntax"
        validation = "valid_hcl_validation_condition"
      }
      
      attribute "error_message" {
        type = "string"
        required = true
        description = "Error message for validation failure"
      }
    }
  }
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

## Template Integration

### **1. Variable Context Integration Strategy**

**Recommendation: Add variables to existing TemplateContext**

```go
// Enhanced TemplateContext with integrated variables
type TemplateContext struct {
    // Existing fields
    Project   *Project
    Machine   *Machine
    Facts     map[string]interface{}
    
    // NEW: Integrated variable context
    Variables *VariableContext
    
    // NEW: Variable resolver for precedence
    varResolver *VariableResolver
}

// VariableContext provides hierarchical storage
type VariableContext struct {
    ProjectVars    map[string]interface{} // variables.hcl and variables/*.hcl
    FactsVars      map[string]interface{} // facts-derived variables (same as Facts)
    EnvironmentVars map[string]string     // environment variables
    DefaultVars    map[string]interface{} // system defaults
    
    // Metadata for debugging and validation
    SourceMap      map[string]VariableSource
    Dependencies   *DependencyGraph
    ResolutionOrder []string
    Cache          map[string]interface{}
}
```

**Pros of Integration:**
- ✅ **Single Source of Truth**: All template data in one context
- ✅ **Consistent Access**: Same pattern for facts, variables, and project data
- ✅ **Performance**: No need to merge contexts at runtime
- ✅ **Debugging**: Easy to inspect all available data
- ✅ **Template Functions**: Can access variables via existing function patterns

**Cons of Integration:**
- ❌ **Increased Complexity**: TemplateContext becomes larger
- ❌ **Tight Coupling**: Variables become part of core template system
- ❌ **Migration**: Existing template code may need updates

**Constraints:**
- Must maintain backward compatibility with existing template functions
- Must handle variable precedence correctly
- Must provide clear error messages for undefined variables

### **2. Variable Interpolation Implementation**

**Recommendation: Support both HCL and template interpolation**

#### **HCL File Interpolation: `${var.name}`**

```hcl
# variables.hcl
variables {
  variable "app_name" {
    type = "string"
    default = "myapp"
  }
  
  variable "port" {
    type = "number"
    default = 8080
  }
}

# machines.hcl - variables automatically available
machines {
  machine "web" {
    hostname = "${var.app_name}-web.example.com"
    port = var.port
    tags = ["web", "${var.app_name}"]
  }
}

# actions.hcl - variables in command interpolation
actions {
  action "deploy" {
    command = "docker run -p ${var.port}:${var.port} ${var.app_name}"
    script = <<-EOF
      #!/bin/bash
      echo "Deploying ${var.app_name} on port ${var.port}"
      docker build -t ${var.app_name} .
    EOF
  }
}
```

**Implementation Details:**
```go
// HCL interpolation happens at parse time
func interpolateHCLVariables(content []byte, variables map[string]interface{}) ([]byte, error) {
    // Use Go template engine for HCL interpolation
    tmpl, err := template.New("hcl").Parse(string(content))
    if err != nil {
        return nil, fmt.Errorf("failed to parse HCL template: %w", err)
    }
    
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, variables); err != nil {
        return nil, fmt.Errorf("failed to interpolate HCL variables: %w", err)
    }
    
    return buf.Bytes(), nil
}

// Integration with existing HCL parsing
func ParseMachinesInventory(filename string) (*MachinesInventory, error) {
    content, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    
    // Load variables first
    variables, err := LoadVariablesConfig(filepath.Dir(filename))
    if err != nil {
        return nil, err
    }
    
    // Interpolate variables in HCL content
    interpolated, err := interpolateHCLVariables(content, variables.ProjectVars)
    if err != nil {
        return nil, err
    }
    
    // Parse interpolated HCL
    parser := hclparse.NewParser()
    file, diags := parser.ParseHCL(interpolated, filename)
    // ... rest of parsing logic
}
```

#### **Template Interpolation: `{{var "name"}}`**

```go
// Enhanced template functions
func (tc *TemplateContext) Var(name string) (interface{}, error) {
    return tc.varResolver.Resolve(name)
}

func (tc *TemplateContext) VarOrDefault(name string, defaultValue interface{}) interface{} {
    if value, err := tc.varResolver.Resolve(name); err == nil {
        return value
    }
    return defaultValue
}

func (tc *TemplateContext) VarRequired(name string) (interface{}, error) {
    value, err := tc.varResolver.Resolve(name)
    if err != nil {
        return nil, fmt.Errorf("required variable '%s' not found", name)
    }
    return value, nil
}

// Backward compatibility
func (tc *TemplateContext) Fact(name string) (interface{}, error) {
    return tc.Var(name) // Facts are now variables too
}

// Debug function
func (tc *TemplateContext) VarSource(name string) string {
    if source, exists := tc.Variables.SourceMap[name]; exists {
        return string(source)
    }
    return "unknown"
}
```

### **3. Practical Impact of HCL Interpolation**

**What HCL interpolation gives us in spooky:**

#### **1. Dynamic Configuration**
```hcl
# variables.hcl
variables {
  variable "environment" {
    type = "string"
    default = "development"
  }
  
  variable "base_domain" {
    type = "string"
    default = "example.com"
  }
}

# machines.hcl - environment-specific configuration
machines {
  machine "web" {
    hostname = "web-${var.environment}.${var.base_domain}"
    port = var.environment == "production" ? 443 : 8080
  }
}
```

#### **2. Reusable Action Templates**
```hcl
# variables.hcl
variables {
  variable "app_name" {
    type = "string"
    default = "myapp"
  }
  
  variable "docker_registry" {
    type = "string"
    default = "docker.io"
  }
}

# actions.hcl - parameterized actions
actions {
  action "deploy" {
    command = "docker pull ${var.docker_registry}/${var.app_name}:latest"
    script = <<-EOF
      docker run -d --name ${var.app_name} ${var.docker_registry}/${var.app_name}:latest
    EOF
  }
}
```

#### **3. Environment-Specific Facts**
```hcl
# variables.hcl
variables {
  variable "monitoring_enabled" {
    type = "bool"
    default = true
  }
}

# facts.hcl - conditional fact collection
facts {
  fact "monitoring_status" {
    condition = var.monitoring_enabled
    command = "systemctl is-active prometheus"
  }
}
```

### **4. Integration with Existing Template Functions**

**Recommendation: Extend existing functions, don't replace them**

```go
// Existing template functions remain unchanged
func (tc *TemplateContext) Env(name string) string {
    return os.Getenv(name)
}

// NEW: Variable functions follow same pattern
func (tc *TemplateContext) Var(name string) (interface{}, error) {
    return tc.varResolver.Resolve(name)
}

// Enhanced template syntax examples
```

**Template Usage Examples:**
```bash
# Existing functions still work
{{env "PATH"}}
{{fact "hostname"}}

# NEW: Variable functions
{{var "app_name"}}
{{varOrDefault "port" 8080}}
{{varRequired "api_key"}}

# Debug functions
{{varSource "app_name"}}  # Returns: "project", "facts", "environment", "default"
```

### **5. Implementation Strategy for spooky**

**Phase 1: Template Context Enhancement**
```go
// 1. Extend TemplateContext
type TemplateContext struct {
    // ... existing fields
    Variables *VariableContext
    varResolver *VariableResolver
}

// 2. Add variable loading to template creation
func NewTemplateContext(project *Project, machine *Machine) *TemplateContext {
    tc := &TemplateContext{
        Project: project,
        Machine: machine,
        Facts:   make(map[string]interface{}),
    }
    
    // Load variables
    variables, err := LoadVariablesConfig(project.Path)
    if err != nil {
        // Log warning, continue without variables
        logging.GetLogger().Warn("failed to load variables", "error", err)
    } else {
        tc.Variables = variables
        tc.varResolver = NewVariableResolver(tc)
    }
    
    return tc
}
```

**Phase 2: HCL Interpolation**
```go
// Add interpolation to all HCL parsing functions
func ParseProjectConfig(filename string) (*ProjectConfig, error) {
    // 1. Load variables first
    variables, err := LoadVariablesConfig(filepath.Dir(filename))
    if err != nil {
        return nil, err
    }
    
    // 2. Interpolate variables in HCL content
    content, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    
    interpolated, err := interpolateHCLVariables(content, variables.ProjectVars)
    if err != nil {
        return nil, err
    }
    
    // 3. Parse interpolated HCL
    // ... existing parsing logic
}
```

**Phase 3: Template Function Integration**
```go
// Add variable functions to template engine
func (tc *TemplateContext) registerVariableFunctions(tmpl *template.Template) {
    tmpl.Funcs(template.FuncMap{
        // Existing functions
        "env": tc.Env,
        "fact": tc.Fact,
        
        // NEW: Variable functions
        "var": tc.Var,
        "varOrDefault": tc.VarOrDefault,
        "varRequired": tc.VarRequired,
        "varSource": tc.VarSource,
    })
}
```

### **6. Benefits for spooky**

1. **Configuration Reusability**: Same variables across machines, actions, and templates
2. **Environment Management**: Easy switching between dev/staging/prod
3. **Reduced Duplication**: Define once, use everywhere
4. **Dynamic Configuration**: Runtime variable resolution with facts
5. **Better Debugging**: Clear variable sources and precedence
6. **Consistent Access**: Same syntax in HCL and templates
7. **Backward Compatibility**: Existing templates continue to work

### **7. Migration Strategy**

1. **Phase 1**: Add variables to TemplateContext (non-breaking)
2. **Phase 2**: Add variable functions to templates (non-breaking)
3. **Phase 3**: Add HCL interpolation (breaking for invalid variable references)
4. **Phase 4**: Deprecate old fact-only access (with warnings)
5. **Phase 5**: Remove deprecated functions (major version bump)

## Schema and Validation

### **1. Leveraging Existing HCL Schema System**

**Recommendation: Use and extend the existing schema system in `internal/schemas/`**

The existing variables schema at `internal/schemas/schemas/variables.hcl` provides a solid foundation, but needs extensions for the complete variables system.

#### **Current Schema Analysis**
```hcl
# internal/schemas/schemas/variables.hcl (existing)
variables {
  variable_blocks = {
    type = "list(object)"
    required = true
    min_items = 0
    properties = {
      name = { type = "string", required = true, pattern = "^[a-z][a-z0-9_]*$" }
      type = { type = "string", required = true, enum = ["string", "number", "bool", "list", "map"] }
      description = { type = "string", required = false, max_length = 500 }
      default = { type = "any", required = false }
      required = { type = "bool", required = false, default = false }
      sensitive = { type = "bool", required = false, default = false }
      validation = {
        type = "object", required = false, max_occurrences = 1
        properties = {
          condition = { type = "string", required = true, max_length = 1000 }
          error_message = { type = "string", required = true, max_length = 500 }
        }
      }
    }
  }
  validation = {
    file_location = { rule = "path", pattern = "^(variables\\.hcl|variables/[^/]+\\.hcl)$" }
    variable_name = { rule = "regex", pattern = "^[a-z][a-z0-9_]*$" }
    variable_type = { rule = "enum", allowed_values = ["string", "number", "bool", "list", "map"] }
    required_variable = { rule = "required", condition = "variable_is_required" }
    validation_condition = { rule = "hcl" }
    no_circular_deps = { rule = "acyclic" }
  }
}
```

#### **Required Schema Extensions**

**1. Add `object_var` type support**
```hcl
# Extended schema for object_var type
variables {
  variable_blocks = {
    type = "list(object)"
    required = true
    min_items = 0
    properties = {
      name = { type = "string", required = true, pattern = "^[a-z][a-z0-9_]*$" }
      type = { type = "string", required = true, enum = ["string", "number", "bool", "list", "map", "object"] }
      description = { type = "string", required = false, max_length = 500 }
      default = { type = "any", required = false }
      required = { type = "bool", required = false, default = false }
      sensitive = { type = "bool", required = false, default = false }
      
      # NEW: Object type schema definition
      object_schema = {
        type = "object", required = false, max_occurrences = 1
        condition = "type == 'object'"
        properties = {
          properties = {
            type = "map(object)"
            required = true
            properties = {
              type = { type = "string", required = true, enum = ["string", "number", "bool", "list", "map"] }
              required = { type = "bool", required = false, default = false }
              description = { type = "string", required = false, max_length = 500 }
            }
          }
          required_properties = { type = "list(string)", required = false }
        }
      }
      
      # NEW: Dependency tracking
      depends_on = { type = "list(string)", required = false }
      
      validation = {
        type = "object", required = false, max_occurrences = 1
        properties = {
          condition = { type = "string", required = true, max_length = 1000 }
          error_message = { type = "string", required = true, max_length = 500 }
        }
      }
    }
  }
  
  # Enhanced validation rules
  validation = {
    file_location = { rule = "path", pattern = "^(variables\\.hcl|variables/[^/]+\\.hcl)$" }
    variable_name = { rule = "regex", pattern = "^[a-z][a-z0-9_]*$" }
    variable_type = { rule = "enum", allowed_values = ["string", "number", "bool", "list", "map", "object"] }
    required_variable = { rule = "required", condition = "variable_is_required" }
    validation_condition = { rule = "hcl" }
    no_circular_deps = { rule = "acyclic" }
    
    # NEW: Object schema validation
    object_schema_required = { rule = "required", condition = "type == 'object' && object_schema == null" }
    object_properties_valid = { rule = "object_properties", condition = "type == 'object'" }
    
    # NEW: Dependency validation
    dependency_exists = { rule = "dependency_exists", condition = "depends_on != null" }
    dependency_acyclic = { rule = "dependency_acyclic", condition = "depends_on != null" }
  }
}
```

**2. Add file merging validation schema**
```hcl
# internal/schemas/schemas/variable-merging.hcl (new)
variable_merging {
  file_patterns = {
    type = "list(string)"
    required = true
    default = ["variables.hcl", "variables/*.hcl"]
  }
  
  conflict_resolution = {
    type = "string"
    required = true
    enum = ["error", "first_wins", "last_wins"]
    default = "error"
  }
  
  validation = {
    no_duplicate_variables = { rule = "no_duplicates", scope = "project" }
    file_readable = { rule = "file_readable" }
    file_parseable = { rule = "hcl_parseable" }
  }
}
```

### **2. Validation Error Reporting**

**Recommendation: Use existing error handling patterns with enhancements**

#### **Existing Error Patterns Analysis**
```go
// Current error patterns in spooky
type ValidationError struct {
    File    string
    Line    int
    Column  int
    Message string
    Code    string
}

type ValidationResult struct {
    Errors   []ValidationError
    Warnings []ValidationError
    Valid    bool
}
```

#### **Enhanced Error Reporting for Variables**
```go
// Enhanced validation error types
type VariableValidationError struct {
    File        string
    Line        int
    Column      int
    Variable    string
    Message     string
    Code        string
    Severity    string // "error" or "warning"
    Source      string // "schema", "dependency", "interpolation", "type"
    Context     map[string]interface{}
}

type VariableValidationResult struct {
    Errors      []VariableValidationError
    Warnings    []VariableValidationError
    Valid       bool
    FileCount   int
    VariableCount int
    Duration    time.Duration
}

// Error codes for variables
const (
    ErrVarDuplicateName     = "VAR_DUPLICATE_NAME"
    ErrVarInvalidName       = "VAR_INVALID_NAME"
    ErrVarInvalidType       = "VAR_INVALID_TYPE"
    ErrVarCircularDependency = "VAR_CIRCULAR_DEPENDENCY"
    ErrVarMissingDependency = "VAR_MISSING_DEPENDENCY"
    ErrVarInvalidValidation = "VAR_INVALID_VALIDATION"
    ErrVarInterpolationFail = "VAR_INTERPOLATION_FAIL"
    ErrVarFileConflict       = "VAR_FILE_CONFLICT"
    ErrVarFileUnreadable     = "VAR_FILE_UNREADABLE"
    ErrVarFileInvalidHCL     = "VAR_FILE_INVALID_HCL"
)
```

#### **Error Reporting Implementation**
```go
// Enhanced validation with development mode
func ValidateVariables(projectPath string, development bool) (*VariableValidationResult, error) {
    result := &VariableValidationResult{
        Errors:   make([]VariableValidationError, 0),
        Warnings: make([]VariableValidationError, 0),
    }
    
    start := time.Now()
    defer func() {
        result.Duration = time.Since(start)
    }()
    
    // Load and validate variables
    variables, err := LoadVariablesConfig(projectPath)
    if err != nil {
        return nil, err
    }
    
    // Schema validation
    if err := validateVariableSchema(variables, result, development); err != nil {
        return nil, err
    }
    
    // Dependency validation
    if err := validateVariableDependencies(variables, result, development); err != nil {
        return nil, err
    }
    
    // Interpolation validation
    if err := validateVariableInterpolation(variables, result, development); err != nil {
        return nil, err
    }
    
    // File conflict validation
    if err := validateFileConflicts(projectPath, result, development); err != nil {
        return nil, err
    }
    
    result.Valid = len(result.Errors) == 0
    result.VariableCount = len(variables.ProjectVars)
    
    return result, nil
}

// Development mode error handling
func addValidationError(result *VariableValidationResult, err VariableValidationError, development bool) {
    if development && err.Severity == "error" {
        // Convert errors to warnings in development mode
        err.Severity = "warning"
        result.Warnings = append(result.Warnings, err)
    } else if err.Severity == "error" {
        result.Errors = append(result.Errors, err)
    } else {
        result.Warnings = append(result.Warnings, err)
    }
}
```

### **3. Strict Validation with Development Mode**

**Recommendation: Strict validation by default with `--development` flag**

#### **CLI Implementation**
```go
// internal/cli/variables.go
var variablesValidateCmd = &cobra.Command{
    Use:   "validate [project directory]",
    Short: "Validate project variables",
    Long:  `Validate variables in the specified project directory.`,
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        projectDir := args[0]
        development, _ := cmd.Flags().GetBool("development")
        
        logger := logging.GetLogger()
        logger.Info("validating variables", "project", projectDir, "development", development)
        
        result, err := ValidateVariables(projectDir, development)
        if err != nil {
            return fmt.Errorf("validation failed: %w", err)
        }
        
        // Report results
        if len(result.Errors) > 0 {
            logger.Error("validation errors found", "count", len(result.Errors))
            for _, err := range result.Errors {
                logger.Error("validation error", 
                    "file", err.File, 
                    "line", err.Line, 
                    "variable", err.Variable,
                    "code", err.Code,
                    "message", err.Message)
            }
            return fmt.Errorf("validation failed with %d errors", len(result.Errors))
        }
        
        if len(result.Warnings) > 0 {
            logger.Warn("validation warnings found", "count", len(result.Warnings))
            for _, warn := range result.Warnings {
                logger.Warn("validation warning",
                    "file", warn.File,
                    "line", warn.Line,
                    "variable", warn.Variable,
                    "code", warn.Code,
                    "message", warn.Message)
            }
        }
        
        logger.Info("validation completed successfully",
            "variables", result.VariableCount,
            "files", result.FileCount,
            "duration", result.Duration)
        
        return nil
    },
}

func init() {
    variablesValidateCmd.Flags().Bool("development", false, "Convert validation errors to warnings")
    variablesCmd.AddCommand(variablesValidateCmd)
}
```

#### **Validation Strictness Levels**

**1. Strict Mode (Default)**
```go
// All validation errors are treated as errors
func validateStrict(variables *VariableContext) []VariableValidationError {
    var errors []VariableValidationError
    
    // Schema validation - all errors
    if err := validateSchema(variables); err != nil {
        errors = append(errors, err)
    }
    
    // Dependency validation - all errors
    if err := validateDependencies(variables); err != nil {
        errors = append(errors, err)
    }
    
    // Type validation - all errors
    if err := validateTypes(variables); err != nil {
        errors = append(errors, err)
    }
    
    return errors
}
```

**2. Development Mode**
```go
// Convert certain errors to warnings
func validateDevelopment(variables *VariableContext) (*VariableValidationResult, error) {
    result := &VariableValidationResult{}
    
    // Schema validation - errors become warnings
    if err := validateSchema(variables); err != nil {
        result.Warnings = append(result.Warnings, convertToWarning(err))
    }
    
    // Dependency validation - still errors (critical)
    if err := validateDependencies(variables); err != nil {
        result.Errors = append(result.Errors, err)
    }
    
    // Type validation - errors become warnings
    if err := validateTypes(variables); err != nil {
        result.Warnings = append(result.Warnings, convertToWarning(err))
    }
    
    return result, nil
}
```

### **4. Schema Integration Implementation**

```go
// Use existing schema infrastructure
func validateVariableSchema(variables *VariableContext, result *VariableValidationResult, development bool) error {
    // Load variables schema
    schema, err := schemas.LoadSchema(schemas.SchemaTypeVariables)
    if err != nil {
        return fmt.Errorf("failed to load variables schema: %w", err)
    }
    
    // Create schema validator
    validator := schemas.NewSchemaValidator()
    if err := validator.AddSchema(schema); err != nil {
        return fmt.Errorf("failed to add variables schema: %w", err)
    }
    
    // Validate each variable file
    for filePath, vars := range variables.ProjectVars {
        if err := validator.ValidateFile(filePath, vars); err != nil {
            validationErr := VariableValidationError{
                File:     filePath,
                Message:  err.Error(),
                Code:     ErrVarInvalidSchema,
                Severity: "error",
                Source:   "schema",
            }
            addValidationError(result, validationErr, development)
        }
    }
    
    return nil
}
```

### **5. Benefits of This Approach**

1. **Consistent with Existing System**: Uses same schema infrastructure as other spooky components
2. **Comprehensive Validation**: Covers schema, dependencies, types, and interpolation
3. **Flexible Error Handling**: Strict by default, lenient in development
4. **Clear Error Messages**: Detailed error reporting with file, line, and context
5. **Performance**: Efficient validation with caching and early termination
6. **Extensible**: Easy to add new validation rules and error types
7. **Integration**: Seamless integration with existing CLI and logging systems

### **6. Migration Strategy**

1. **Phase 1**: Extend existing variables schema (non-breaking)
2. **Phase 2**: Add new validation rules (non-breaking)
3. **Phase 3**: Implement strict validation (breaking for invalid configurations)
4. **Phase 4**: Add development mode (non-breaking)
5. **Phase 5**: Deprecate old validation patterns (with warnings)

## Performance and Scalability

### **1. Large Variable File Handling**

**Recommendation: Optimize for typical use cases, not extreme edge cases**

#### **Real-World Variable File Sizes**

**Ansible Deployments:**
- **Small projects**: 1-10 variables, <1KB files
- **Medium projects**: 50-200 variables, 5-50KB files
- **Large projects**: 500-2000 variables, 100-500KB files
- **Enterprise deployments**: 2000-10000 variables, 500KB-2MB files

**Puppet Deployments:**
- **Small environments**: 10-50 variables, 2-10KB files
- **Medium environments**: 100-500 variables, 20-100KB files
- **Large environments**: 1000-5000 variables, 200KB-1MB files
- **Enterprise environments**: 5000-20000 variables, 1-5MB files

**Spooky Target Sizes:**
- **Typical projects**: 10-100 variables, 2-20KB files
- **Large projects**: 100-1000 variables, 20-200KB files
- **Enterprise projects**: 1000-5000 variables, 200KB-1MB files
- **Maximum supported**: 10MB per file (reasonable limit)

#### **Implementation Strategy**

```go
// Optimized variable loading for typical sizes
func LoadVariablesConfig(projectPath string) (*VariableContext, error) {
    // Use standard file I/O for typical sizes
    content, err := os.ReadFile(filepath.Join(projectPath, "variables.hcl"))
    if err != nil {
        if os.IsNotExist(err) {
            return &VariableContext{}, nil // No variables file
        }
        return nil, err
    }
    
    // Check file size before processing
    if len(content) > 10*1024*1024 { // 10MB limit
        return nil, fmt.Errorf("variables file too large: %d bytes (max 10MB)", len(content))
    }
    
    // Parse with standard HCL parser
    return parseVariablesContent(content, projectPath)
}

// Efficient directory scanning for variables/*.hcl
func loadVariableFiles(projectPath string) (map[string][]byte, error) {
    variablesDir := filepath.Join(projectPath, "variables")
    files := make(map[string][]byte)
    
    entries, err := os.ReadDir(variablesDir)
    if err != nil {
        if os.IsNotExist(err) {
            return files, nil // No variables directory
        }
        return nil, err
    }
    
    totalSize := 0
    for _, entry := range entries {
        if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".hcl") {
            continue
        }
        
        filePath := filepath.Join(variablesDir, entry.Name())
        content, err := os.ReadFile(filePath)
        if err != nil {
            return nil, fmt.Errorf("failed to read %s: %w", filePath, err)
        }
        
        // Track cumulative size
        totalSize += len(content)
        if totalSize > 50*1024*1024 { // 50MB total limit
            return nil, fmt.Errorf("total variables files too large: %d bytes (max 50MB)", totalSize)
        }
        
        files[entry.Name()] = content
    }
    
    return files, nil
}
```

### **2. Variable Resolution Caching**

**Recommendation: Implement intelligent caching with invalidation**

#### **Caching Strategy**

```go
// Variable resolver with caching
type VariableResolver struct {
    context *VariableContext
    cache   map[string]interface{}
    cacheHits   int64
    cacheMisses int64
    mu      sync.RWMutex
}

func (vr *VariableResolver) Resolve(name string) (interface{}, error) {
    // Check cache first (read lock)
    vr.mu.RLock()
    if value, exists := vr.cache[name]; exists {
        atomic.AddInt64(&vr.cacheHits, 1)
        vr.mu.RUnlock()
        return value, nil
    }
    vr.mu.RUnlock()
    
    atomic.AddInt64(&vr.cacheMisses, 1)
    
    // Resolve variable (write lock for cache update)
    vr.mu.Lock()
    defer vr.mu.Unlock()
    
    // Double-check cache after acquiring write lock
    if value, exists := vr.cache[name]; exists {
        return value, nil
    }
    
    // Resolve variable with precedence
    value, err := vr.resolveWithPrecedence(name)
    if err != nil {
        return nil, err
    }
    
    // Cache resolved value
    if vr.cache == nil {
        vr.cache = make(map[string]interface{})
    }
    vr.cache[name] = value
    
    return value, nil
}

// Cache statistics for monitoring
func (vr *VariableResolver) CacheStats() (hits, misses int64) {
    return atomic.LoadInt64(&vr.cacheHits), atomic.LoadInt64(&vr.cacheMisses)
}

// Cache invalidation for development
func (vr *VariableResolver) InvalidateCache() {
    vr.mu.Lock()
    defer vr.mu.Unlock()
    vr.cache = make(map[string]interface{})
}
```

#### **Caching Benefits**

**Performance Impact:**
- **First resolution**: Full precedence lookup (environment → project → facts → defaults)
- **Subsequent resolutions**: O(1) cache lookup
- **Typical cache hit rate**: 80-95% in most deployments
- **Memory overhead**: ~100 bytes per cached variable

**Real-World Examples:**
```go
// Template rendering with high variable reuse
func renderTemplate(tc *TemplateContext, templateName string) (string, error) {
    // Variables like 'app_name', 'environment', 'version' used 50+ times
    // Caching provides significant performance benefit
    
    tmpl, err := template.New(templateName).Parse(templateContent)
    if err != nil {
        return "", err
    }
    
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, tc); err != nil {
        return "", err
    }
    
    return buf.String(), nil
}
```

### **3. Variable Resolution Cost Analysis**

**Is variable resolution costly? It depends on the implementation.**

#### **Cost Breakdown**

**1. Simple Variable Lookup (Cached)**
```go
// O(1) - Direct map lookup
func (vr *VariableResolver) getCached(name string) (interface{}, bool) {
    vr.mu.RLock()
    defer vr.mu.RUnlock()
    value, exists := vr.cache[name]
    return value, exists
}
// Cost: ~10-50 nanoseconds
```

**2. Full Precedence Resolution (Uncached)**
```go
// O(n) where n = number of variable sources
func (vr *VariableResolver) resolveWithPrecedence(name string) (interface{}, error) {
    // 1. Environment variables (fast)
    if value, exists := os.LookupEnv("SPOOKY_" + strings.ToUpper(name)); exists {
        return value, nil
    }
    
    // 2. Project variables (map lookup)
    if value, exists := vr.context.ProjectVars[name]; exists {
        return vr.interpolateProjectVar(value)
    }
    
    // 3. Facts (map lookup)
    if value, exists := vr.context.Facts[name]; exists {
        return value, nil
    }
    
    // 4. Defaults (map lookup)
    if value, exists := vr.context.DefaultVars[name]; exists {
        return value, nil
    }
    
    return nil, fmt.Errorf("variable '%s' not found", name)
}
// Cost: ~100-500 nanoseconds (first time)
```

**3. Variable Interpolation (Expensive)**
```go
// O(m) where m = complexity of interpolation
func (vr *VariableResolver) interpolateProjectVar(value interface{}) (interface{}, error) {
    switch v := value.(type) {
    case string:
        // Check if string contains variable references
        if strings.Contains(v, "${var.") {
            return vr.interpolateString(v)
        }
        return v, nil
    default:
        return v, nil
    }
}
// Cost: ~1-10 microseconds (depends on interpolation complexity)
```

#### **Performance Optimization Strategies**

**1. Lazy Interpolation**
```go
// Only interpolate when actually needed
type LazyVariable struct {
    rawValue interface{}
    resolved interface{}
    resolver *VariableResolver
}

func (lv *LazyVariable) Value() (interface{}, error) {
    if lv.resolved != nil {
        return lv.resolved, nil
    }
    
    resolved, err := lv.resolver.interpolateProjectVar(lv.rawValue)
    if err != nil {
        return nil, err
    }
    
    lv.resolved = resolved
    return resolved, nil
}
```

**2. Pre-computed Interpolation**
```go
// Interpolate all variables at load time
func (vc *VariableContext) PrecomputeInterpolation() error {
    for name, value := range vc.ProjectVars {
        if interpolated, err := vc.interpolateProjectVar(value); err != nil {
            return fmt.Errorf("failed to interpolate variable '%s': %w", name, err)
        } else {
            vc.ProjectVars[name] = interpolated
        }
    }
    return nil
}
```

### **4. Variable Immutability**

**Recommendation: Variables should be immutable once loaded**

#### **How Other Projects Handle This**

**Ansible:**
- Variables are immutable during playbook execution
- Changes require playbook restart
- **Pros**: Predictable behavior, no race conditions
- **Cons**: Cannot adapt to runtime changes

**Puppet:**
- Variables are immutable during catalog compilation
- Changes require catalog recompilation
- **Pros**: Consistent state, no side effects
- **Cons**: Cannot respond to dynamic changes

**Chef:**
- Variables can be modified during chef-client run
- **Pros**: Dynamic adaptation
- **Cons**: Complex debugging, potential race conditions

**Terraform:**
- Variables are immutable during plan/apply
- Changes require new plan
- **Pros**: Predictable, reproducible
- **Cons**: Cannot adapt to runtime conditions

#### **Spooky Implementation: Immutable Variables**

```go
// Immutable variable context
type VariableContext struct {
    ProjectVars    map[string]interface{} // Read-only after initialization
    FactsVars      map[string]interface{} // Read-only after initialization
    EnvironmentVars map[string]string     // Read-only after initialization
    DefaultVars    map[string]interface{} // Read-only after initialization
    
    // Metadata
    SourceMap      map[string]VariableSource
    Dependencies   *DependencyGraph
    ResolutionOrder []string
    Cache          map[string]interface{}
    
    // Immutability flag
    initialized    bool
    mu             sync.RWMutex
}

// Thread-safe read-only access
func (vc *VariableContext) GetProjectVar(name string) (interface{}, bool) {
    vc.mu.RLock()
    defer vc.mu.RUnlock()
    value, exists := vc.ProjectVars[name]
    return value, exists
}

// Prevent modification after initialization
func (vc *VariableContext) SetProjectVar(name string, value interface{}) error {
    vc.mu.Lock()
    defer vc.mu.Unlock()
    
    if vc.initialized {
        return fmt.Errorf("cannot modify variables after initialization")
    }
    
    if vc.ProjectVars == nil {
        vc.ProjectVars = make(map[string]interface{})
    }
    vc.ProjectVars[name] = value
    return nil
}

// Mark as initialized (immutable)
func (vc *VariableContext) Initialize() {
    vc.mu.Lock()
    defer vc.mu.Unlock()
    vc.initialized = true
}
```

#### **Benefits of Immutability**

**1. Predictable Behavior**
```go
// Variables remain consistent throughout execution
func runAction(action *Action, tc *TemplateContext) error {
    // Variables won't change during action execution
    appName, _ := tc.Var("app_name")
    port, _ := tc.Var("port")
    
    // Safe to use these values throughout the action
    command := fmt.Sprintf("docker run -p %v:%v %v", port, port, appName)
    return executeCommand(command)
}
```

**2. Thread Safety**
```go
// Multiple goroutines can safely access variables
func runParallelActions(actions []*Action, tc *TemplateContext) error {
    var wg sync.WaitGroup
    errChan := make(chan error, len(actions))
    
    for _, action := range actions {
        wg.Add(1)
        go func(a *Action) {
            defer wg.Done()
            // Safe concurrent access to immutable variables
            if err := runAction(a, tc); err != nil {
                errChan <- err
            }
        }(action)
    }
    
    wg.Wait()
    close(errChan)
    
    // Check for errors
    for err := range errChan {
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

**3. Debugging and Testing**
```go
// Variables state is predictable for debugging
func debugVariableState(tc *TemplateContext) {
    logger := logging.GetLogger()
    
    logger.Info("variable state", 
        "project_vars", len(tc.Variables.ProjectVars),
        "facts_vars", len(tc.Variables.FactsVars),
        "cache_hits", tc.varResolver.CacheStats())
    
    // State won't change during debugging
    for name, source := range tc.Variables.SourceMap {
        logger.Debug("variable source", "name", name, "source", source)
    }
}
```

### **5. Performance Monitoring**

```go
// Performance metrics for variables
type VariableMetrics struct {
    LoadTime       time.Duration
    ParseTime      time.Duration
    CacheHits      int64
    CacheMisses    int64
    Interpolations int64
    FileCount      int
    VariableCount  int
    TotalSize      int64
}

func (vm *VariableMetrics) CacheHitRate() float64 {
    total := vm.CacheHits + vm.CacheMisses
    if total == 0 {
        return 0
    }
    return float64(vm.CacheHits) / float64(total) * 100
}

// Integration with existing metrics
func (tc *TemplateContext) GetVariableMetrics() *VariableMetrics {
    hits, misses := tc.varResolver.CacheStats()
    return &VariableMetrics{
        CacheHits:   hits,
        CacheMisses: misses,
        // ... other metrics
    }
}
```

### **6. Best Practices for Performance**

**1. File Organization**
```hcl
# Good: Separate files by concern
# variables/environment.hcl - 50 variables
# variables/application.hcl - 30 variables
# variables/infrastructure.hcl - 20 variables

# Avoid: Single massive file
# variables.hcl - 1000+ variables
```

**2. Caching Strategy**
```go
// Enable caching for production, disable for development
func NewVariableResolver(context *VariableContext, enableCache bool) *VariableResolver {
    vr := &VariableResolver{context: context}
    if enableCache {
        vr.cache = make(map[string]interface{})
    }
    return vr
}
```

**3. Lazy Loading**
```go
// Only load variables when needed
func (tc *TemplateContext) ensureVariablesLoaded() error {
    if tc.Variables != nil {
        return nil
    }
    
    variables, err := LoadVariablesConfig(tc.Project.Path)
    if err != nil {
        return err
    }
    
    tc.Variables = variables
    tc.varResolver = NewVariableResolver(variables, true)
    return nil
}
```

### **7. Migration Strategy**

1. **Phase 1**: Implement basic variable loading with size limits
2. **Phase 2**: Add caching with monitoring
3. **Phase 3**: Implement immutability with thread safety
4. **Phase 4**: Add performance metrics and optimization
5. **Phase 5**: Document best practices and performance guidelines

## Error Handling and User Experience

### **1. Variable Interpolation Error Handling**

**Recommendation: Error immediately unless `--development` mode is used**

#### **Error Handling Strategy**

```go
// Variable interpolation with development mode support
type InterpolationError struct {
    Variable    string
    File        string
    Line        int
    Column      int
    Message     string
    Code        string
    Context     string
}

// Interpolation error codes
const (
    ErrVarUndefined     = "VAR_UNDEFINED"
    ErrVarCircularRef   = "VAR_CIRCULAR_REF"
    ErrVarInvalidSyntax = "VAR_INVALID_SYNTAX"
    ErrVarTypeMismatch  = "VAR_TYPE_MISMATCH"
    ErrVarInterpolation = "VAR_INTERPOLATION_FAIL"
)

// Enhanced variable resolver with development mode
func (vr *VariableResolver) interpolateString(content string, development bool) (string, error) {
    // Use Go template engine for interpolation
    tmpl, err := template.New("interpolation").Parse(content)
    if err != nil {
        return "", &InterpolationError{
            Variable: "syntax",
            Message:  fmt.Sprintf("invalid interpolation syntax: %v", err),
            Code:     ErrVarInvalidSyntax,
        }
    }
    
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, vr); err != nil {
        // Check if it's an undefined variable error
        if strings.Contains(err.Error(), "undefined variable") {
            varName := extractVariableName(err.Error())
            if development {
                // In development mode, use empty string as default
                return "", nil
            } else {
                // In production mode, return error
                return "", &InterpolationError{
                    Variable: varName,
                    Message:  fmt.Sprintf("undefined variable '%s'", varName),
                    Code:     ErrVarUndefined,
                }
            }
        }
        return "", &InterpolationError{
            Variable: "interpolation",
            Message:  fmt.Sprintf("interpolation failed: %v", err),
            Code:     ErrVarInterpolation,
        }
    }
    
    return buf.String(), nil
}

// CLI integration with development mode
func runVariablesValidate(cmd *cobra.Command, args []string) error {
    projectDir := args[0]
    development, _ := cmd.Flags().GetBool("development")
    
    logger := logging.GetLogger()
    logger.Info("validating variables", "project", projectDir, "development", development)
    
    // Load and validate variables with development mode
    result, err := ValidateVariables(projectDir, development)
    if err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    // Handle interpolation errors based on mode
    for _, err := range result.Errors {
        if err.Code == ErrVarUndefined && development {
            // Convert undefined variable errors to warnings in development mode
            logger.Warn("undefined variable (development mode)", 
                "variable", err.Variable,
                "file", err.File,
                "line", err.Line,
                "message", "using empty string as default")
        } else {
            logger.Error("validation error", 
                "variable", err.Variable,
                "file", err.File,
                "line", err.Line,
                "code", err.Code,
                "message", err.Message)
        }
    }
    
    return nil
}
```

#### **Development Mode Behavior**

```go
// Development mode examples
func (vr *VariableResolver) Resolve(name string) (interface{}, error) {
    value, err := vr.resolveWithPrecedence(name)
    if err != nil {
        if vr.developmentMode {
            // In development mode, return empty string for undefined variables
            return "", nil
        }
        return nil, err
    }
    return value, nil
}

// Template function with development mode
func (tc *TemplateContext) Var(name string) (interface{}, error) {
    value, err := tc.varResolver.Resolve(name)
    if err != nil {
        if tc.developmentMode {
            // Log warning and return empty string
            logging.GetLogger().Warn("undefined variable in template", 
                "variable", name,
                "template", tc.currentTemplate,
                "using_default", "")
            return "", nil
        }
        return nil, err
    }
    return value, nil
}
```

### **2. Helpful Error Messages**

**Recommendation: Provide intelligent error suggestions with configurable verbosity**

#### **Error Suggestion Implementation**

```go
// Intelligent error suggestions
type ErrorSuggestion struct {
    Variable    string
    Suggestions []string
    Confidence  float64
    Reason      string
}

func suggestVariableNames(undefinedVar string, availableVars []string) *ErrorSuggestion {
    suggestions := make([]string, 0)
    
    // Calculate similarity scores
    scores := make(map[string]float64)
    for _, available := range availableVars {
        score := calculateSimilarity(undefinedVar, available)
        if score > 0.6 { // 60% similarity threshold
            scores[available] = score
        }
    }
    
    // Sort by similarity score
    type scorePair struct {
        name  string
        score float64
    }
    var pairs []scorePair
    for name, score := range scores {
        pairs = append(pairs, scorePair{name, score})
    }
    sort.Slice(pairs, func(i, j int) bool {
        return pairs[i].score > pairs[j].score
    })
    
    // Take top 3 suggestions
    for i := 0; i < 3 && i < len(pairs); i++ {
        suggestions = append(suggestions, pairs[i].name)
    }
    
    confidence := 0.0
    if len(suggestions) > 0 {
        confidence = scores[suggestions[0]]
    }
    
    return &ErrorSuggestion{
        Variable:    undefinedVar,
        Suggestions: suggestions,
        Confidence:  confidence,
        Reason:      determineSuggestionReason(undefinedVar, suggestions),
    }
}

// Similarity calculation using Levenshtein distance
func calculateSimilarity(a, b string) float64 {
    distance := levenshteinDistance(a, b)
    maxLen := max(len(a), len(b))
    if maxLen == 0 {
        return 1.0
    }
    return 1.0 - float64(distance)/float64(maxLen)
}

// Enhanced error reporting with suggestions
func (vr *VariableResolver) ResolveWithSuggestions(name string) (interface{}, error) {
    value, err := vr.resolveWithPrecedence(name)
    if err != nil {
        if strings.Contains(err.Error(), "not found") {
            // Get all available variable names
            availableVars := make([]string, 0)
            for varName := range vr.context.ProjectVars {
                availableVars = append(availableVars, varName)
            }
            for varName := range vr.context.FactsVars {
                availableVars = append(availableVars, varName)
            }
            
            // Generate suggestions
            suggestion := suggestVariableNames(name, availableVars)
            
            // Create enhanced error message
            enhancedErr := fmt.Errorf("variable '%s' not found", name)
            if len(suggestion.Suggestions) > 0 {
                enhancedErr = fmt.Errorf("variable '%s' not found. Did you mean: %s?", 
                    name, strings.Join(suggestion.Suggestions, ", "))
            }
            
            return nil, enhancedErr
        }
        return nil, err
    }
    return value, nil
}
```

#### **Configurable Verbosity**

```go
// CLI flag for error verbosity
var variablesValidateCmd = &cobra.Command{
    Use:   "validate [project directory]",
    Short: "Validate project variables",
    Long:  `Validate variables in the specified project directory.`,
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        projectDir := args[0]
        development, _ := cmd.Flags().GetBool("development")
        verbose, _ := cmd.Flags().GetBool("verbose")
        
        logger := logging.GetLogger()
        
        result, err := ValidateVariables(projectDir, development)
        if err != nil {
            return fmt.Errorf("validation failed: %w", err)
        }
        
        // Report errors with optional suggestions
        for _, err := range result.Errors {
            if verbose && err.Code == ErrVarUndefined {
                // Show suggestions in verbose mode
                suggestion := suggestVariableNames(err.Variable, getAvailableVariables(result))
                if len(suggestion.Suggestions) > 0 {
                    logger.Error("validation error with suggestions", 
                        "variable", err.Variable,
                        "file", err.File,
                        "line", err.Line,
                        "suggestions", suggestion.Suggestions,
                        "confidence", suggestion.Confidence)
                } else {
                    logger.Error("validation error", 
                        "variable", err.Variable,
                        "file", err.File,
                        "line", err.Line,
                        "message", err.Message)
                }
            } else {
                logger.Error("validation error", 
                    "variable", err.Variable,
                    "file", err.File,
                    "line", err.Line,
                    "message", err.Message)
            }
        }
        
        return nil
    },
}

func init() {
    variablesValidateCmd.Flags().Bool("development", false, "Convert validation errors to warnings")
    variablesValidateCmd.Flags().Bool("verbose", false, "Show detailed error suggestions")
    variablesCmd.AddCommand(variablesValidateCmd)
}
```

### **3. Sensitive Variable Handling**

**Recommendation: Single sensitivity level with age encryption integration**

#### **Sensitive Variable Definition**

```hcl
# variables.hcl
variables {
  variable "database_password" {
    type = "string"
    description = "Database connection password"
    sensitive = true  # Single sensitivity level
    default = "changeme"
  }
  
  variable "api_key" {
    type = "string"
    description = "API authentication key"
    sensitive = true
    required = true
  }
  
  variable "app_name" {
    type = "string"
    description = "Application name"
    sensitive = false  # Default
    default = "myapp"
  }
}
```

#### **Age Encryption Integration**

```go
// Age encryption for sensitive variables
type AgeEncryption struct {
    PublicKey  string
    PrivateKey string
    Enabled    bool
}

type SensitiveVariable struct {
    Name      string
    Value     string
    Encrypted bool
    AgeKey    string
}

// Encrypt sensitive variable using age
func encryptSensitiveVariable(value string, publicKey string) (string, error) {
    // Use age library for encryption
    recipient, err := age.ParseX25519Recipient(publicKey)
    if err != nil {
        return "", fmt.Errorf("invalid age public key: %w", err)
    }
    
    var buf bytes.Buffer
    w, err := age.Encrypt(&buf, recipient)
    if err != nil {
        return "", fmt.Errorf("failed to create age encryptor: %w", err)
    }
    
    if _, err := w.Write([]byte(value)); err != nil {
        return "", fmt.Errorf("failed to encrypt value: %w", err)
    }
    
    if err := w.Close(); err != nil {
        return "", fmt.Errorf("failed to finalize encryption: %w", err)
    }
    
    return buf.String(), nil
}

// Decrypt sensitive variable using age
func decryptSensitiveVariable(encryptedValue string, privateKey string) (string, error) {
    // Use age library for decryption
    identity, err := age.ParseX25519Identity(privateKey)
    if err != nil {
        return "", fmt.Errorf("invalid age private key: %w", err)
    }
    
    r, err := age.Decrypt(strings.NewReader(encryptedValue), identity)
    if err != nil {
        return "", fmt.Errorf("failed to create age decryptor: %w", err)
    }
    
    decrypted, err := io.ReadAll(r)
    if err != nil {
        return "", fmt.Errorf("failed to decrypt value: %w", err)
    }
    
    return string(decrypted), nil
}
```

#### **Variable Storage with Encryption**

```go
// Enhanced variable context with encryption support
type VariableContext struct {
    ProjectVars    map[string]interface{}
    FactsVars      map[string]interface{}
    EnvironmentVars map[string]string
    DefaultVars    map[string]interface{}
    
    // Encryption support
    AgeEncryption  *AgeEncryption
    SensitiveVars  map[string]*SensitiveVariable
    
    // Metadata
    SourceMap      map[string]VariableSource
    Dependencies   *DependencyGraph
    ResolutionOrder []string
    Cache          map[string]interface{}
    
    // Immutability
    initialized    bool
    mu             sync.RWMutex
}

// Load variables with encryption support
func LoadVariablesConfig(projectPath string, agePrivateKey string) (*VariableContext, error) {
    vc := &VariableContext{
        ProjectVars: make(map[string]interface{}),
        SensitiveVars: make(map[string]*SensitiveVariable),
    }
    
    // Load age encryption if private key provided
    if agePrivateKey != "" {
        vc.AgeEncryption = &AgeEncryption{
            PrivateKey: agePrivateKey,
            Enabled:    true,
        }
    }
    
    // Load and parse variables
    if err := loadVariableFiles(projectPath, vc); err != nil {
        return nil, err
    }
    
    // Decrypt sensitive variables
    if err := decryptSensitiveVariables(vc); err != nil {
        return nil, err
    }
    
    return vc, nil
}

// Decrypt sensitive variables
func decryptSensitiveVariables(vc *VariableContext) error {
    if !vc.AgeEncryption.Enabled {
        return nil
    }
    
    for name, sensitiveVar := range vc.SensitiveVars {
        if sensitiveVar.Encrypted {
            decrypted, err := decryptSensitiveVariable(sensitiveVar.Value, vc.AgeEncryption.PrivateKey)
            if err != nil {
                return fmt.Errorf("failed to decrypt variable '%s': %w", name, err)
            }
            vc.ProjectVars[name] = decrypted
        }
    }
    
    return nil
}
```

#### **Logging and Output Redaction**

```go
// Redact sensitive variables in logs and output
func (vc *VariableContext) RedactSensitiveData(data interface{}) interface{} {
    switch v := data.(type) {
    case map[string]interface{}:
        redacted := make(map[string]interface{})
        for key, value := range v {
            if vc.isSensitiveVariable(key) {
                redacted[key] = "[REDACTED]"
            } else {
                redacted[key] = vc.RedactSensitiveData(value)
            }
        }
        return redacted
    case []interface{}:
        redacted := make([]interface{}, len(v))
        for i, value := range v {
            redacted[i] = vc.RedactSensitiveData(value)
        }
        return redacted
    default:
        return v
    }
}

// Check if variable is sensitive
func (vc *VariableContext) isSensitiveVariable(name string) bool {
    if sensitiveVar, exists := vc.SensitiveVars[name]; exists {
        return sensitiveVar.Encrypted
    }
    return false
}

// Enhanced logging with redaction
func (tc *TemplateContext) LogVariableAccess(name string, value interface{}) {
    logger := logging.GetLogger()
    
    if tc.Variables.isSensitiveVariable(name) {
        logger.Debug("accessed sensitive variable", 
            "name", name,
            "value", "[REDACTED]",
            "source", tc.Variables.SourceMap[name])
    } else {
        logger.Debug("accessed variable", 
            "name", name,
            "value", value,
            "source", tc.Variables.SourceMap[name])
    }
}
```

#### **Export with Encryption**

```go
// Export variables with encryption support
func (vc *VariableContext) ExportVariables(format string, encrypt bool, agePublicKey string) ([]byte, error) {
    var data interface{}
    
    if encrypt {
        // Create encrypted export
        encryptedVars := make(map[string]interface{})
        for name, value := range vc.ProjectVars {
            if vc.isSensitiveVariable(name) {
                // Encrypt sensitive variables
                encrypted, err := encryptSensitiveVariable(fmt.Sprintf("%v", value), agePublicKey)
                if err != nil {
                    return nil, fmt.Errorf("failed to encrypt variable '%s': %w", name, err)
                }
                encryptedVars[name] = map[string]interface{}{
                    "encrypted": true,
                    "value":     encrypted,
                    "age_key":   agePublicKey,
                }
            } else {
                encryptedVars[name] = value
            }
        }
        data = encryptedVars
    } else {
        // Create redacted export
        data = vc.RedactSensitiveData(vc.ProjectVars)
    }
    
    // Export in requested format
    switch format {
    case "json":
        return json.MarshalIndent(data, "", "  ")
    case "hcl":
        return vc.exportToHCL(data)
    default:
        return nil, fmt.Errorf("unsupported export format: %s", format)
    }
}
```

### **4. CLI Integration for Sensitive Variables**

```go
// CLI commands with encryption support
var variablesExportCmd = &cobra.Command{
    Use:   "export [project directory]",
    Short: "Export project variables",
    Long:  `Export variables from the specified project directory.`,
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        projectDir := args[0]
        format, _ := cmd.Flags().GetString("format")
        output, _ := cmd.Flags().GetString("output")
        encrypt, _ := cmd.Flags().GetBool("encrypt")
        ageKey, _ := cmd.Flags().GetString("age-key")
        
        logger := logging.GetLogger()
        
        // Load variables
        variables, err := LoadVariablesConfig(projectDir, "")
        if err != nil {
            return fmt.Errorf("failed to load variables: %w", err)
        }
        
        // Export with encryption if requested
        var exportData []byte
        if encrypt {
            if ageKey == "" {
                return fmt.Errorf("--age-key required when --encrypt is specified")
            }
            exportData, err = variables.ExportVariables(format, true, ageKey)
        } else {
            exportData, err = variables.ExportVariables(format, false, "")
        }
        
        if err != nil {
            return fmt.Errorf("failed to export variables: %w", err)
        }
        
        // Write to output
        if output == "" {
            fmt.Println(string(exportData))
        } else {
            if err := os.WriteFile(output, exportData, 0644); err != nil {
                return fmt.Errorf("failed to write output file: %w", err)
            }
            logger.Info("variables exported", "file", output, "format", format, "encrypted", encrypt)
        }
        
        return nil
    },
}

func init() {
    variablesExportCmd.Flags().String("format", "json", "Export format (json, hcl)")
    variablesExportCmd.Flags().String("output", "", "Output file (default: stdout)")
    variablesExportCmd.Flags().Bool("encrypt", false, "Encrypt sensitive variables")
    variablesExportCmd.Flags().String("age-key", "", "Age public key for encryption")
    variablesCmd.AddCommand(variablesExportCmd)
}
```

### **5. Benefits of This Approach**

1. **Clear Error Handling**: Immediate errors in production, warnings in development
2. **Intelligent Suggestions**: Helpful error messages with variable name suggestions
3. **Single Sensitivity Level**: Simple `sensitive = true` flag
4. **Age Encryption Integration**: Secure encryption for sensitive variables
5. **Comprehensive Redaction**: Automatic redaction in logs and exports
6. **Flexible Export**: Support for encrypted and redacted exports
7. **Development-Friendly**: Verbose mode for debugging and suggestions

### **6. Migration Strategy**

1. **Phase 1**: Implement basic error handling with development mode
2. **Phase 2**: Add intelligent error suggestions (optional verbose mode)
3. **Phase 3**: Implement sensitive variable handling with redaction
4. **Phase 4**: Add age encryption integration
5. **Phase 5**: Add encrypted export functionality

## Testing Strategy

### **1. Integration Tests Structure**

**Recommendation: Follow existing test structure in `examples/testing/`**

#### **Test Project Organization**

```bash
# Follow existing test structure
examples/testing/
├── test-variables-basic/
│   ├── project.hcl
│   ├── variables.hcl
│   ├── machines.hcl
│   ├── actions.hcl
│   └── README.md
├── test-variables-dependencies/
│   ├── project.hcl
│   ├── variables/
│   │   ├── environment.hcl
│   │   ├── application.hcl
│   │   └── infrastructure.hcl
│   ├── machines.hcl
│   ├── actions.hcl
│   └── README.md
├── test-variables-circular-deps/
│   ├── project.hcl
│   ├── variables.hcl
│   ├── machines.hcl
│   ├── actions.hcl
│   └── README.md
├── test-variables-encryption/
│   ├── project.hcl
│   ├── variables.hcl
│   ├── machines.hcl
│   ├── actions.hcl
│   ├── keys/
│   │   ├── age-key.txt
│   │   └── age-key.pub
│   └── README.md
└── test-variables-performance/
    ├── project.hcl
    ├── variables.hcl
    ├── machines.hcl
    ├── actions.hcl
    └── README.md
```

#### **Basic Variables Test**

```hcl
# examples/testing/test-variables-basic/variables.hcl
variables {
  variable "app_name" {
    type = "string"
    description = "Application name"
    default = "myapp"
  }
  
  variable "environment" {
    type = "string"
    description = "Deployment environment"
    default = "development"
  }
  
  variable "port" {
    type = "number"
    description = "Application port"
    default = 8080
  }
  
  variable "debug_enabled" {
    type = "bool"
    description = "Enable debug mode"
    default = false
  }
  
  variable "database_config" {
    type = "map"
    description = "Database configuration"
    default = {
      host = "localhost"
      port = 5432
      name = "myapp_db"
    }
  }
}

# examples/testing/test-variables-basic/machines.hcl
machines {
  machine "web" {
    hostname = "${var.app_name}-web.${var.environment}.example.com"
    port = var.port
    tags = ["web", var.environment]
  }
}

# examples/testing/test-variables-basic/actions.hcl
actions {
  action "deploy" {
    command = "docker run -p ${var.port}:${var.port} ${var.app_name}"
    script = <<-EOF
      #!/bin/bash
      echo "Deploying ${var.app_name} on port ${var.port}"
      echo "Environment: ${var.environment}"
      echo "Debug: ${var.debug_enabled}"
    EOF
  }
}
```

#### **Test Implementation**

```go
// tests/integration/variables_test.go
package integration

import (
    "testing"
    "path/filepath"
    
    "spooky/internal/config"
    "spooky/internal/facts"
    "spooky/internal/logging"
)

func TestVariablesBasic(t *testing.T) {
    testDir := "examples/testing/test-variables-basic"
    
    // Test variable loading
    t.Run("LoadVariables", func(t *testing.T) {
        variables, err := config.LoadVariablesConfig(testDir)
        if err != nil {
            t.Fatalf("failed to load variables: %v", err)
        }
        
        // Verify variables loaded correctly
        expectedVars := map[string]interface{}{
            "app_name": "myapp",
            "environment": "development",
            "port": 8080,
            "debug_enabled": false,
        }
        
        for name, expectedValue := range expectedVars {
            if value, exists := variables.ProjectVars[name]; !exists {
                t.Errorf("variable '%s' not found", name)
            } else if value != expectedValue {
                t.Errorf("variable '%s' = %v, expected %v", name, value, expectedValue)
            }
        }
    })
    
    // Test variable interpolation in machines
    t.Run("MachineInterpolation", func(t *testing.T) {
        machines, err := config.ParseMachinesInventory(filepath.Join(testDir, "machines.hcl"))
        if err != nil {
            t.Fatalf("failed to parse machines: %v", err)
        }
        
        webMachine := machines.Machines["web"]
        expectedHostname := "myapp-web.development.example.com"
        if webMachine.Hostname != expectedHostname {
            t.Errorf("hostname = %s, expected %s", webMachine.Hostname, expectedHostname)
        }
    })
    
    // Test variable interpolation in actions
    t.Run("ActionInterpolation", func(t *testing.T) {
        actions, err := config.LoadActionsConfig(testDir)
        if err != nil {
            t.Fatalf("failed to load actions: %v", err)
        }
        
        deployAction := actions.Actions["deploy"]
        expectedCommand := "docker run -p 8080:8080 myapp"
        if deployAction.Command != expectedCommand {
            t.Errorf("command = %s, expected %s", deployAction.Command, expectedCommand)
        }
    })
}
```

### **2. Dependency Resolution Testing**

**Recommendation: Create comprehensive test cases for complex dependency graphs**

#### **Directed Acyclic Graph (DAG) Implementation**

```go
// internal/variables/dependency.go
package variables

import (
    "fmt"
    "sort"
)

// DependencyGraph represents a directed acyclic graph
type DependencyGraph struct {
    nodes map[string]*DependencyNode
    edges map[string][]string // node -> dependencies
}

type DependencyNode struct {
    Name         string
    Dependencies []string
    Dependents   []string
    Visited      bool
    InProgress   bool
}

// NewDependencyGraph creates a new dependency graph
func NewDependencyGraph() *DependencyGraph {
    return &DependencyGraph{
        nodes: make(map[string]*DependencyNode),
        edges: make(map[string][]string),
    }
}

// AddNode adds a node to the dependency graph
func (dg *DependencyGraph) AddNode(name string) {
    if _, exists := dg.nodes[name]; !exists {
        dg.nodes[name] = &DependencyNode{
            Name: name,
        }
    }
}

// AddDependency adds a dependency relationship
func (dg *DependencyGraph) AddDependency(from, to string) error {
    dg.AddNode(from)
    dg.AddNode(to)
    
    // Check for circular dependency
    if dg.wouldCreateCycle(from, to) {
        return fmt.Errorf("circular dependency detected: %s -> %s", from, to)
    }
    
    // Add edge
    dg.edges[from] = append(dg.edges[from], to)
    dg.nodes[from].Dependencies = append(dg.nodes[from].Dependencies, to)
    dg.nodes[to].Dependents = append(dg.nodes[to].Dependents, from)
    
    return nil
}

// Detect cycles using DFS
func (dg *DependencyGraph) wouldCreateCycle(from, to string) bool {
    // Reset visited flags
    for _, node := range dg.nodes {
        node.Visited = false
        node.InProgress = false
    }
    
    return dg.hasCycleDFS(to, from)
}

func (dg *DependencyGraph) hasCycleDFS(current, target string) bool {
    if current == target {
        return true
    }
    
    node := dg.nodes[current]
    if node.Visited {
        return false
    }
    
    node.Visited = true
    node.InProgress = true
    
    for _, dep := range dg.edges[current] {
        if dg.hasCycleDFS(dep, target) {
            return true
        }
    }
    
    node.InProgress = false
    return false
}

// TopologicalSort returns variables in dependency order
func (dg *DependencyGraph) TopologicalSort() ([]string, error) {
    // Reset visited flags
    for _, node := range dg.nodes {
        node.Visited = false
        node.InProgress = false
    }
    
    var result []string
    var stack []string
    
    // Get all nodes
    for name := range dg.nodes {
        stack = append(stack, name)
    }
    
    // Sort for deterministic order
    sort.Strings(stack)
    
    for _, node := range stack {
        if !dg.nodes[node].Visited {
            if err := dg.topologicalSortDFS(node, &result); err != nil {
                return nil, err
            }
        }
    }
    
    // Reverse to get dependency order
    for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
        result[i], result[j] = result[j], result[i]
    }
    
    return result, nil
}

func (dg *DependencyGraph) topologicalSortDFS(node string, result *[]string) error {
    dgNode := dg.nodes[node]
    
    if dgNode.InProgress {
        return fmt.Errorf("circular dependency detected involving %s", node)
    }
    
    if dgNode.Visited {
        return nil
    }
    
    dgNode.Visited = true
    dgNode.InProgress = true
    
    // Process dependencies first
    for _, dep := range dg.edges[node] {
        if err := dg.topologicalSortDFS(dep, result); err != nil {
            return err
        }
    }
    
    dgNode.InProgress = false
    *result = append(*result, node)
    
    return nil
}
```

#### **Complex Dependency Test Cases**

```hcl
# examples/testing/test-variables-dependencies/variables/environment.hcl
variables {
  variable "environment" {
    type = "string"
    description = "Deployment environment"
    default = "development"
  }
  
  variable "base_domain" {
    type = "string"
    description = "Base domain for services"
    default = "example.com"
  }
  
  variable "region" {
    type = "string"
    description = "AWS region"
    default = "us-west-2"
  }
}

# examples/testing/test-variables-dependencies/variables/application.hcl
variables {
  variable "app_name" {
    type = "string"
    description = "Application name"
    default = "myapp"
    depends_on = ["environment"]
  }
  
  variable "app_version" {
    type = "string"
    description = "Application version"
    default = "1.0.0"
  }
  
  variable "app_domain" {
    type = "string"
    description = "Application domain"
    default = "${var.app_name}.${var.environment}.${var.base_domain}"
    depends_on = ["app_name", "environment", "base_domain"]
  }
}

# examples/testing/test-variables-dependencies/variables/infrastructure.hcl
variables {
  variable "database_host" {
    type = "string"
    description = "Database host"
    default = "db-${var.environment}.${var.region}.${var.base_domain}"
    depends_on = ["environment", "region", "base_domain"]
  }
  
  variable "redis_host" {
    type = "string"
    description = "Redis host"
    default = "redis-${var.environment}.${var.region}.${var.base_domain}"
    depends_on = ["environment", "region", "base_domain"]
  }
  
  variable "load_balancer_domain" {
    type = "string"
    description = "Load balancer domain"
    default = "lb-${var.app_name}.${var.environment}.${var.base_domain}"
    depends_on = ["app_name", "environment", "base_domain"]
  }
}
```

#### **Dependency Resolution Tests**

```go
// tests/integration/dependency_test.go
package integration

import (
    "testing"
    "reflect"
    
    "spooky/internal/variables"
)

func TestDependencyResolution(t *testing.T) {
    testDir := "examples/testing/test-variables-dependencies"
    
    t.Run("LoadDependencies", func(t *testing.T) {
        variables, err := config.LoadVariablesConfig(testDir)
        if err != nil {
            t.Fatalf("failed to load variables: %v", err)
        }
        
        // Verify dependency graph
        graph := variables.Dependencies
        
        // Check specific dependencies
        expectedDeps := map[string][]string{
            "app_name": {"environment"},
            "app_domain": {"app_name", "environment", "base_domain"},
            "database_host": {"environment", "region", "base_domain"},
            "redis_host": {"environment", "region", "base_domain"},
            "load_balancer_domain": {"app_name", "environment", "base_domain"},
        }
        
        for varName, expectedDependencies := range expectedDeps {
            node := graph.nodes[varName]
            if node == nil {
                t.Errorf("dependency node for '%s' not found", varName)
                continue
            }
            
            if !reflect.DeepEqual(node.Dependencies, expectedDependencies) {
                t.Errorf("dependencies for '%s' = %v, expected %v", 
                    varName, node.Dependencies, expectedDependencies)
            }
        }
    })
    
    t.Run("TopologicalSort", func(t *testing.T) {
        variables, err := config.LoadVariablesConfig(testDir)
        if err != nil {
            t.Fatalf("failed to load variables: %v", err)
        }
        
        sorted, err := variables.Dependencies.TopologicalSort()
        if err != nil {
            t.Fatalf("topological sort failed: %v", err)
        }
        
        // Verify that dependencies come before dependents
        positions := make(map[string]int)
        for i, name := range sorted {
            positions[name] = i
        }
        
        // Check that app_domain comes after its dependencies
        if pos := positions["app_domain"]; pos <= positions["app_name"] || 
           pos <= positions["environment"] || pos <= positions["base_domain"] {
            t.Errorf("app_domain should come after its dependencies in topological sort")
        }
    })
    
    t.Run("VariableResolution", func(t *testing.T) {
        variables, err := config.LoadVariablesConfig(testDir)
        if err != nil {
            t.Fatalf("failed to load variables: %v", err)
        }
        
        // Test resolved values
        expectedValues := map[string]string{
            "app_domain": "myapp.development.example.com",
            "database_host": "db-development.us-west-2.example.com",
            "redis_host": "redis-development.us-west-2.example.com",
            "load_balancer_domain": "lb-myapp.development.example.com",
        }
        
        for varName, expectedValue := range expectedValues {
            if value, exists := variables.ProjectVars[varName]; !exists {
                t.Errorf("variable '%s' not found", varName)
            } else if value != expectedValue {
                t.Errorf("variable '%s' = %v, expected %v", varName, value, expectedValue)
            }
        }
    })
}
```

### **3. Circular Dependency Testing**

```hcl
# examples/testing/test-variables-circular-deps/variables.hcl
variables {
  variable "app_name" {
    type = "string"
    description = "Application name"
    default = "myapp"
    depends_on = ["app_domain"]  # Circular dependency!
  }
  
  variable "app_domain" {
    type = "string"
    description = "Application domain"
    default = "${var.app_name}.example.com"
    depends_on = ["app_name"]  # Circular dependency!
  }
  
  variable "database_url" {
    type = "string"
    description = "Database URL"
    default = "postgresql://${var.db_user}:${var.db_password}@${var.db_host}"
    depends_on = ["db_user", "db_password", "db_host"]
  }
  
  variable "db_user" {
    type = "string"
    description = "Database user"
    default = "app_user"
  }
  
  variable "db_password" {
    type = "string"
    description = "Database password"
    default = "secret"
    sensitive = true
  }
  
  variable "db_host" {
    type = "string"
    description = "Database host"
    default = "localhost"
  }
}
```

```go
// tests/integration/circular_dependency_test.go
package integration

import (
    "testing"
    "strings"
)

func TestCircularDependencyDetection(t *testing.T) {
    testDir := "examples/testing/test-variables-circular-deps"
    
    t.Run("DetectCircularDependency", func(t *testing.T) {
        _, err := config.LoadVariablesConfig(testDir)
        if err == nil {
            t.Fatal("expected error for circular dependency, got none")
        }
        
        if !strings.Contains(err.Error(), "circular dependency") {
            t.Errorf("expected circular dependency error, got: %v", err)
        }
        
        if !strings.Contains(err.Error(), "app_name") || !strings.Contains(err.Error(), "app_domain") {
            t.Errorf("error should mention both variables in circular dependency: %v", err)
        }
    })
    
    t.Run("ValidDependencies", func(t *testing.T) {
        // Test that valid dependencies work
        variables, err := config.LoadVariablesConfig("examples/testing/test-variables-basic")
        if err != nil {
            t.Fatalf("failed to load valid variables: %v", err)
        }
        
        // Verify no circular dependencies
        _, err = variables.Dependencies.TopologicalSort()
        if err != nil {
            t.Errorf("topological sort failed for valid dependencies: %v", err)
        }
    })
}
```

### **4. How DAGs Help with Conflict Resolution and Dependency Management**

#### **Directed Acyclic Graph Benefits**

**1. Conflict Resolution**
```go
// DAGs prevent circular dependencies at build time
func (dg *DependencyGraph) AddDependency(from, to string) error {
    // Before adding edge, check if it would create a cycle
    if dg.wouldCreateCycle(from, to) {
        return fmt.Errorf("circular dependency detected: %s -> %s", from, to)
    }
    
    // Safe to add edge
    dg.edges[from] = append(dg.edges[from], to)
    return nil
}
```

**2. Dependency Resolution Order**
```go
// Topological sort ensures correct resolution order
func (dg *DependencyGraph) TopologicalSort() ([]string, error) {
    var result []string
    
    // Process nodes in dependency order
    for _, node := range dg.nodes {
        if !node.Visited {
            if err := dg.topologicalSortDFS(node.Name, &result); err != nil {
                return nil, err
            }
        }
    }
    
    return result, nil
}
```

**3. Parallel Processing**
```go
// DAGs enable parallel processing of independent variables
func (dg *DependencyGraph) GetIndependentNodes() []string {
    var independent []string
    
    for name, node := range dg.nodes {
        if len(node.Dependencies) == 0 {
            independent = append(independent, name)
        }
    }
    
    return independent
}

func (dg *DependencyGraph) ProcessInParallel() error {
    // Process independent nodes in parallel
    independent := dg.GetIndependentNodes()
    
    var wg sync.WaitGroup
    errChan := make(chan error, len(independent))
    
    for _, nodeName := range independent {
        wg.Add(1)
        go func(name string) {
            defer wg.Done()
            if err := processVariable(name); err != nil {
                errChan <- err
            }
        }(nodeName)
    }
    
    wg.Wait()
    close(errChan)
    
    // Check for errors
    for err := range errChan {
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

**4. Change Impact Analysis**
```go
// DAGs enable efficient change impact analysis
func (dg *DependencyGraph) GetImpactedNodes(changedNode string) []string {
    var impacted []string
    visited := make(map[string]bool)
    
    dg.getImpactedNodesDFS(changedNode, &impacted, visited)
    
    return impacted
}

func (dg *DependencyGraph) getImpactedNodesDFS(node string, impacted *[]string, visited map[string]bool) {
    if visited[node] {
        return
    }
    
    visited[node] = true
    *impacted = append(*impacted, node)
    
    // All dependents are impacted
    for _, dependent := range dg.nodes[node].Dependents {
        dg.getImpactedNodesDFS(dependent, impacted, visited)
    }
}
```

### **5. Performance Testing**

```go
// tests/integration/performance_test.go
package integration

import (
    "testing"
    "time"
    
    "spooky/internal/variables"
)

func TestVariablePerformance(t *testing.T) {
    testDir := "examples/testing/test-variables-performance"
    
    t.Run("LoadPerformance", func(t *testing.T) {
        start := time.Now()
        
        variables, err := config.LoadVariablesConfig(testDir)
        if err != nil {
            t.Fatalf("failed to load variables: %v", err)
        }
        
        duration := time.Since(start)
        
        // Performance assertions
        if duration > 100*time.Millisecond {
            t.Errorf("variable loading took too long: %v", duration)
        }
        
        if len(variables.ProjectVars) < 100 {
            t.Errorf("expected at least 100 variables, got %d", len(variables.ProjectVars))
        }
    })
    
    t.Run("ResolutionPerformance", func(t *testing.T) {
        variables, err := config.LoadVariablesConfig(testDir)
        if err != nil {
            t.Fatalf("failed to load variables: %v", err)
        }
        
        resolver := variables.NewVariableResolver()
        
        start := time.Now()
        
        // Resolve variables multiple times to test caching
        for i := 0; i < 1000; i++ {
            for varName := range variables.ProjectVars {
                _, err := resolver.Resolve(varName)
                if err != nil {
                    t.Errorf("failed to resolve variable '%s': %v", varName, err)
                }
            }
        }
        
        duration := time.Since(start)
        
        // Performance assertions
        if duration > 1*time.Second {
            t.Errorf("variable resolution took too long: %v", duration)
        }
        
        // Check cache hit rate
        hits, misses := resolver.CacheStats()
        hitRate := float64(hits) / float64(hits+misses) * 100
        
        if hitRate < 80 {
            t.Errorf("cache hit rate too low: %.2f%%", hitRate)
        }
    })
}
```

### **6. Benefits of DAG-Based Dependency Management**

1. **Conflict Prevention**: DAGs prevent circular dependencies at build time
2. **Correct Resolution Order**: Topological sort ensures variables are resolved in correct order
3. **Parallel Processing**: Independent variables can be processed in parallel
4. **Change Impact Analysis**: Efficiently determine which variables are impacted by changes
5. **Validation**: Easy to validate dependency relationships
6. **Debugging**: Clear visualization of dependency relationships
7. **Performance**: Efficient algorithms for dependency resolution

### **7. Migration Strategy**

1. **Phase 1**: Implement basic DAG structure and circular dependency detection
2. **Phase 2**: Add topological sorting for resolution order
3. **Phase 3**: Implement parallel processing for independent variables
4. **Phase 4**: Add change impact analysis
5. **Phase 5**: Create comprehensive test suite with performance benchmarks

## Backward Compatibility

### **1. Migration Strategy**

**Recommendation: No backward compatibility - create new projects**

Since spooky is still in early development and we don't care about backward compatibility yet, we'll focus on creating new projects rather than providing migration tools.

#### **New Project Approach**

```bash
# New projects will use the variables system from the start
spooky project init my-new-project
cd my-new-project

# Project structure with variables
my-new-project/
├── project.hcl
├── variables.hcl          # NEW: Variables file
├── variables/             # NEW: Variables directory
│   ├── environment.hcl
│   └── application.hcl
├── machines.hcl
├── actions.hcl
├── templates/
└── facts.db/             # Existing facts database
```

#### **No Migration Tools**

```go
// No migration tools needed - focus on new projects
func initProject(projectPath string) error {
    // Create new project structure
    if err := createProjectStructure(projectPath); err != nil {
        return fmt.Errorf("failed to create project structure: %w", err)
    }
    
    // Initialize variables system
    if err := initializeVariablesSystem(projectPath); err != nil {
        return fmt.Errorf("failed to initialize variables system: %w", err)
    }
    
    // Initialize existing systems (facts, actions, etc.)
    if err := initializeExistingSystems(projectPath); err != nil {
        return fmt.Errorf("failed to initialize existing systems: %w", err)
    }
    
    return nil
}

func initializeVariablesSystem(projectPath string) error {
    // Create variables.hcl with basic structure
    variablesContent := `variables {
  # Define your project variables here
  variable "app_name" {
    type = "string"
    description = "Application name"
    default = "myapp"
  }
  
  variable "environment" {
    type = "string"
    description = "Deployment environment"
    default = "development"
  }
}`
    
    if err := os.WriteFile(filepath.Join(projectPath, "variables.hcl"), []byte(variablesContent), 0644); err != nil {
        return fmt.Errorf("failed to create variables.hcl: %w", err)
    }
    
    // Create variables directory
    variablesDir := filepath.Join(projectPath, "variables")
    if err := os.MkdirAll(variablesDir, 0755); err != nil {
        return fmt.Errorf("failed to create variables directory: %w", err)
    }
    
    return nil
}
```

### **3. Benefits of This Approach**

1. **No Migration Complexity**: Focus on new projects without backward compatibility concerns
2. **Data File Integration**: Variables can reference and enhance existing data files
3. **Flexible Data Access**: Support for JSON, YAML, and HCL data files
4. **Enhanced Templates**: Templates can access both variables and data files
5. **Validation**: Data file validation and management tools
6. **Environment Overrides**: Variables can override data file values based on environment
7. **Clear Separation**: Variables for configuration, data files for static data

### **4. Migration Strategy**

1. **Phase 1**: Implement variables system for new projects only
2. **Phase 2**: Add data file integration and validation
3. **Phase 3**: Add CLI tools for data file management
4. **Phase 4**: Document best practices for data file organization
5. **Phase 5**: Add advanced data file features (encryption, validation rules)