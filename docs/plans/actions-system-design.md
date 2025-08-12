# Actions System: Comprehensive Implementation Plan

## Overview

This document is the authoritative source for all actions system implementation details in spooky. It covers action definition, dependency management, file merging, run planning, and integration with all other spooky systems.

**Schema Integration**: This actions system follows the schema validation patterns and infrastructure defined in [Schema System](../schema-system.md) for action definition validation, dependency schema enforcement, and schema-based run planning.

**Architecture Integration**: Actions integrate with the overall spooky architecture as described in [Spooky Design](../spooky-design.md), providing dynamic running capabilities and workflow automation for all system components.

## System Integration

This actions system integrates with other core Spooky systems to provide comprehensive action management and running:

### **CLI System Integration**
- **Action Commands**: Actions managed through `spooky actions` commands (see [CLI System](../cli-system.md))
- **Command Patterns**: Actions follow the established `spooky noun verb` CLI pattern
- **Validation Commands**: `spooky actions validate` for action file validation
- **Run Commands**: `spooky actions run` and `spooky actions dry-run` for action running
- **List Commands**: `spooky actions list` for action discovery and filtering

### **Facts System Integration**
- **Facts in Actions**: Actions can access machine facts for dynamic running (see [Facts System](../facts-system.md))
- **Facts-Based Commands**: Actions can use facts data in command running
- **Template Context**: Actions have access to facts through template context
- **Facts Resolution**: Actions resolve facts at run time for each machine

### **Project System Integration**
- **Project Actions**: Actions stored in project-specific `actions.hcl` and `actions/` directory (see [Project System](../project-system.md))
- **Project Context**: Actions run within project run context
- **Project Configuration**: Action settings configured in project context
- **Project Isolation**: Project actions isolated from global actions

### **Variables System Integration**
- **Variables in Actions**: Actions can use project variables for dynamic configuration (see [Variables System](../variables-system.md))
- **Variable Resolution**: Actions resolve variables at run time
- **Template Integration**: Actions can use variables in command templates
- **Variable Precedence**: Actions follow variable resolution precedence (Environment → Project → Facts → Defaults)

### **Configuration System Integration**
- **Global Configuration**: Action run settings from `$XDG_CONFIG_HOME/spooky/spooky.hcl` (see [Configuration System](../configuration-system.md))
- **Run Settings**: Timeouts, retries, parallel run limits
- **SSH Configuration**: SSH settings for remote action running
- **Logging Configuration**: Action run logging and output settings

### **Machines System Integration**
- **Target Selection**: Actions use machine inventory for target selection via tags, names, or filters (see [Machines System](../machines-system.md))
- **Machine Resolution**: Actions resolve machine targets through the machines system's indexing and lookup capabilities
- **Execution Context**: Actions run within machine-specific contexts with authentication and connection details
- **Parallel Running**: Machine inventory supports parallel action running across multiple targets

### **Template System Integration**
- **Template Actions**: Actions can run template rendering and deployment (see [Template System](../template-system.md))
- **Template Context**: Actions provide template context with facts, variables, and machine data
- **Template Deployment**: Actions can deploy rendered templates to target machines
- **Template Validation**: Actions validate templates before running

## Current State Analysis

### **What We Have**
- ✅ **Basic file parsing** for HCL configuration files
- ✅ **Individual file validation** for syntax and structure
- ✅ **Template context** with facts, machines, and actions

### **What We Need**
- 🔄 **Actions system** with `actions.hcl` and `actions/` directory
- 🔄 **File merging** for action files
- 🔄 **Dependency tracking** across actions
- 🔄 **Circular reference detection** to prevent infinite loops
- 🔄 **Dependency resolution** with topological sorting
- 🔄 **Cross-file dependency management** for merged files

## Actions System Design

### **1. File Structure**

```
project/
├── actions.hcl            # Main actions file
├── actions/               # Organized action files
│   ├── setup.hcl          # Setup actions
│   ├── deploy.hcl         # Deployment actions
│   ├── cleanup.hcl        # Cleanup actions
│   └── maintenance.hcl    # Maintenance actions
├── variables.hcl          # Main variables file
├── variables/             # Organized variable files
└── ...
```

### **2. Action Definition Schema**

```hcl
# actions.hcl
actions {
  action "setup_database" {
    command = "echo 'Setting up database'"
    description = "Initialize database and create tables"
    depends_on = ["validate_config", "check_permissions"]
    timeout = 300
    retries = 3
  }
  
  action "deploy_app" {
    command = "echo 'Deploying application'"
    description = "Deploy the application to servers"
    depends_on = ["setup_database", "build_app"]
    timeout = 600
    retries = 2
  }
  
  action "health_check" {
    command = "echo 'Checking application health'"
    description = "Verify application is running correctly"
    depends_on = ["deploy_app"]
    timeout = 60
    retries = 1
  }
}
```

### **3. Action Types and Properties**

```hcl
# Supported action properties
action "example" {
  command     = "string"           # Required: command to run
  description = "string"           # Optional: action description
  depends_on  = ["list", "of", "actions"]  # Optional: dependencies
  timeout     = 300                # Optional: timeout in seconds
  retries     = 3                  # Optional: number of retries
  parallel    = false              # Optional: run in parallel with others
  tags        = ["tag1", "tag2"]   # Optional: action tags
  machine     = "machine_name"     # Optional: specific machine to run on
  user        = "username"         # Optional: user to run as
  working_dir = "/path/to/dir"     # Optional: working directory
}
```

## File Merging Implementation

### **1. Merging Strategy**

```go
type ActionsFileMerger struct {
    MaxFileSize    int64  // Default: 100MB
    TempFileSize   int64  // Default: 10MB
    MergeStrategy  string // "memory" or "tempfile"
    projectPath    string
}

type MergedActionsContent struct {
    Actions    map[string]*Action
    SourceFiles []string
    FileOrder   []string  // Order of file processing
    Dependencies map[string][]string
    CircularRefs []string
}
```

### **2. Merging Process**

```go
func (afm *ActionsFileMerger) MergeActions(projectPath string) (*MergedActionsContent, error) {
    // 1. Validate project path and ensure it exists
    // 2. Read actions.hcl (only from project root)
    // 3. Read all .hcl files in actions/ directory (only from project root)
    // 4. Validate all files are valid HCL
    // 5. Merge in memory if total size < TempFileSize
    // 6. Use temporary file if size > TempFileSize
    // 7. Build dependency graph and detect circular references
    // 8. Return merged content with metadata
}
```

### **3. File Size Handling**

```go
func (afm *ActionsFileMerger) shouldUseTempFile(files []string) bool {
    totalSize := int64(0)
    for _, file := range files {
        if info, err := os.Stat(file); err == nil {
            totalSize += info.Size()
        }
    }
    return totalSize > afm.TempFileSize
}
```

### **4. Merging Rules**

```go
type ActionsMerger struct {
    projectPath string
    // Precedence: actions/ files override actions.hcl
    // Last file wins for duplicate action names
    // All files are processed in alphabetical order within directories
    // Only files within project directory are read
    // All files must be valid HCL
    
    func (am *ActionsMerger) MergeActions(projectPath string) (map[string]*Action, error) {
        actions := make(map[string]*Action)
        
        // Validate project path
        if err := am.validateProjectPath(projectPath); err != nil {
            return nil, err
        }
        
        // Process main actions.hcl file (only from project root)
        mainFile := filepath.Join(projectPath, "actions.hcl")
        if fileExists(mainFile) {
            if err := am.processFile(mainFile, actions); err != nil {
                return nil, fmt.Errorf("error in actions.hcl: %w", err)
            }
        }
        
        // Process actions/ directory (only from project root)
        actionsDir := filepath.Join(projectPath, "actions")
        if dirExists(actionsDir) {
            if err := am.processActionsDirectory(actionsDir, actions); err != nil {
                return nil, err
            }
        }
        
        return actions, nil
    }
}
```

## Dependency Management

### **1. Dependency Graph Structure**

```go
type DependencyNode struct {
    Name         string
    Dependencies []string
    Dependents   []string
    File         string
    Line         int
    Type         string  // "action"
    Content      interface{}
}

type DependencyGraph struct {
    Nodes map[string]*DependencyNode
    Edges map[string][]string
    Metadata map[string]interface{}
}
```

### **2. Dependency Declaration Syntax**

```hcl
# Action dependencies
action "setup_database" {
  command = "echo 'Setting up database'"
  depends_on = ["validate_config", "check_permissions"]
}

action "deploy_app" {
  command = "echo 'Deploying application'"
  depends_on = ["setup_database", "build_app"]
}

action "run_tests" {
  command = "echo 'Running tests'"
  depends_on = ["deploy_app", "setup_test_env"]
}

action "cleanup" {
  command = "echo 'Cleaning up'"
  depends_on = ["run_tests", "backup_data"]
}
```

### **3. Cross-File Dependencies**

```hcl
# actions/setup.hcl
actions {
  action "validate_config" {
    command = "echo 'Validating configuration'"
  }
  
  action "check_permissions" {
    command = "echo 'Checking permissions'"
    depends_on = ["validate_config"]
  }
}

# actions/deploy.hcl
actions {
  action "build_app" {
    command = "echo 'Building application'"
  }
  
  action "deploy_app" {
    command = "echo 'Deploying application'"
    depends_on = ["setup_database", "build_app"]  # References action in another file
  }
}
```

## Circular Reference Detection

### **1. Detection Algorithm**

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

### **2. Circular Reference Reporting**

```go
type CircularRefReport struct {
    Cycle        []string
    Files        []string
    Lines        []int
    Types        []string
    Description  string
    Suggestions  []string
}

func (dg *DependencyGraph) GenerateCircularRefReport(circular []string) *CircularRefReport {
    report := &CircularRefReport{
        Cycle: circular,
        Files: make([]string, len(circular)),
        Lines: make([]int, len(circular)),
        Types: make([]string, len(circular)),
    }
    
    for i, nodeName := range circular {
        node := dg.Nodes[nodeName]
        report.Files[i] = node.File
        report.Lines[i] = node.Line
        report.Types[i] = node.Type
    }
    
    report.Description = fmt.Sprintf("Circular dependency detected: %s", strings.Join(circular, " -> "))
    report.Suggestions = dg.generateSuggestions(circular)
    
    return report
}
```

## Dependency Resolution

### **1. Topological Sorting**

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

### **2. Resolution Context**

```go
type ActionsResolutionContext struct {
    Actions map[string]*Action
    Order   []string
    Errors  []error
}

func (arc *ActionsResolutionContext) ResolveDependencies(graph *DependencyGraph) error {
    // 1. Detect circular references
    if circular := graph.DetectCircularRefs(); len(circular) > 0 {
        report := graph.GenerateCircularRefReport(circular)
        return fmt.Errorf("circular dependency: %s", report.Description)
    }
    
    // 2. Get resolution order
    order, err := graph.ResolveOrder()
    if err != nil {
        return err
    }
    
    // 3. Resolve in order
    for _, nodeName := range order {
        if err := arc.resolveAction(nodeName, graph); err != nil {
            arc.Errors = append(arc.Errors, err)
        }
    }
    
    if len(arc.Errors) > 0 {
        return fmt.Errorf("dependency resolution failed: %v", arc.Errors)
    }
    
    return nil
}
```

## File Location and Validation Requirements

### **1. Strict File Location Rules**

```go
// Only these locations are allowed for action files:
// - project/actions.hcl (project root only)
// - project/actions/*.hcl (actions directory only)

func (am *ActionsMerger) validateProjectPath(projectPath string) error {
    absPath, err := filepath.Abs(projectPath)
    if err != nil {
        return fmt.Errorf("invalid project path: %w", err)
    }
    
    if !dirExists(absPath) {
        return fmt.Errorf("project directory does not exist: %s", absPath)
    }
    
    am.projectPath = absPath
    return nil
}

func (am *ActionsMerger) validateFilePath(filePath string) error {
    absPath, err := filepath.Abs(filePath)
    if err != nil {
        return fmt.Errorf("invalid file path: %w", err)
    }
    
    // Must be within project directory
    if !am.isWithinProject(absPath) {
        return fmt.Errorf("file must be within project directory: %s", filePath)
    }
    
    // Must be actions.hcl or in actions/ directory
    if !am.isValidActionFile(absPath) {
        return fmt.Errorf("invalid action file location: %s", filePath)
    }
    
    return nil
}

func (am *ActionsMerger) isWithinProject(filePath string) bool {
    return strings.HasPrefix(filePath, am.projectPath+string(filepath.Separator))
}

func (am *ActionsMerger) isValidActionFile(filePath string) bool {
    fileName := filepath.Base(filePath)
    dirName := filepath.Base(filepath.Dir(filePath))
    
    // Allow actions.hcl in project root
    if fileName == "actions.hcl" && dirName == filepath.Base(am.projectPath) {
        return true
    }
    
    // Allow .hcl files in actions/ directory
    if strings.HasSuffix(fileName, ".hcl") && dirName == "actions" {
        return true
    }
    
    return false
}
```

### **2. Strict HCL Validation**

```go
func (am *ActionsMerger) processFile(filePath string, actions map[string]*Action) error {
    // Validate file location first
    if err := am.validateFilePath(filePath); err != nil {
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
    
    // Parse as actions with strict schema validation
    var actionBlocks []ActionBlock
    if err := hcl.Unmarshal(content, &actionBlocks); err != nil {
        return fmt.Errorf("invalid action definition in %s: %w", filepath.Base(filePath), err)
    }
    
    // Validate each action block
    for i, actionBlock := range actionBlocks {
        if err := am.validateActionBlock(actionBlock, filePath, i+1); err != nil {
            return fmt.Errorf("action %d in %s: %w", i+1, filepath.Base(filePath), err)
        }
        
        // Add to actions map
        actions[actionBlock.Name] = &Action{
            Name:        actionBlock.Name,
            Command:     actionBlock.Command,
            Description: actionBlock.Description,
            DependsOn:   actionBlock.DependsOn,
            Timeout:     actionBlock.Timeout,
            Retries:     actionBlock.Retries,
            Parallel:    actionBlock.Parallel,
            Tags:        actionBlock.Tags,
            Machine:     actionBlock.Machine,
            User:        actionBlock.User,
            WorkingDir:  actionBlock.WorkingDir,
            SourceFile:  filePath,
            Line:        actionBlock.Line,
        }
    }
    
    return nil
}

func (am *ActionsMerger) processActionsDirectory(dirPath string, actions map[string]*Action) error {
    // Ensure we're only reading from the project's actions directory
    if !am.isWithinProject(dirPath) {
        return fmt.Errorf("actions directory must be within project: %s", dirPath)
    }
    
    // Only read .hcl files
    files, err := filepath.Glob(filepath.Join(dirPath, "*.hcl"))
    if err != nil {
        return fmt.Errorf("error reading actions directory: %w", err)
    }
    
    // Process files in alphabetical order
    sort.Strings(files)
    for _, file := range files {
        if err := am.processFile(file, actions); err != nil {
            return err
        }
    }
    
    return nil
}
```

### **3. Error Examples**

```bash
# Invalid HCL syntax
$ spooky actions validate
Error: invalid HCL in actions.hcl: unexpected token '}' at line 5

# File outside project directory
$ spooky actions validate
Error: invalid action file location: /etc/actions.hcl

# Invalid file in actions directory
$ spooky actions validate
Error: invalid action file location: actions/config.txt

# Missing required command
$ spooky actions validate
Error: action 2 in deploy.hcl: required field 'command' is missing

# Circular dependency
$ spooky actions validate
Error: circular dependency detected: setup_database -> deploy_app -> setup_database

# Multiple validation errors
$ spooky actions validate
Error: invalid HCL in actions/setup.hcl: unexpected '=' at line 3
Error: action 'deploy_app' in actions/deploy.hcl: depends_on references undefined action 'build_app'
Error: invalid action file location: actions/config.txt

# Success case
$ spooky actions validate
✓ All action files are valid
✓ Found 8 actions in 4 files
✓ No circular dependencies detected
```

## Implementation Phases

### **Phase 1: Core Actions System (Week 1-2)**

1. **Action Definition Schema**
   - Create Go structs for action definitions
   - Implement HCL parsing for actions
   - Add validation rules and type checking

2. **Strict File Location Validation**
   - Implement project path validation
   - Restrict file reading to project directory only
   - Validate file locations (actions.hcl and actions/*.hcl only)

3. **Strict HCL Validation**
   - Implement HCL syntax validation for all action files
   - Treat invalid HCL as errors, not warnings
   - Provide clear error messages with file and line information

4. **Basic File Merging**
   - Implement `actions.hcl` reading (project root only)
   - Implement `actions/` directory scanning (project directory only)
   - Basic in-memory merging with validation

5. **Action Storage**
   - Create action registry
   - Implement action resolution
   - Add to project context

6. **CLI Command: actions validate**
   - Implement `spooky actions validate` command
   - Validate all action files in project
   - Report file location and HCL validation issues
   - Provide detailed error messages with file and line numbers
   - Exit with non-zero code on validation failures

### **Phase 2: Dependency Management (Week 3)**

1. **Dependency Graph**
   - Build dependency tracking system
   - Implement node and edge management
   - Add file and line tracking

2. **Circular Reference Detection**
   - Implement cycle detection algorithm
   - Add detailed error reporting
   - Create visualization tools

3. **Basic Resolution**
   - Implement topological sorting
   - Add dependency order resolution
   - Basic error handling

### **Phase 3: Cross-File Dependencies (Week 4)**

1. **File Merging Integration**
   - Integrate with file merging system
   - Handle dependencies across merged files
   - Validate cross-file references

2. **Reference Resolution**
   - Implement cross-file reference tracking
   - Add validation for undefined references
   - Handle action references

3. **Error Reporting**
   - Clear error messages for dependency issues
   - Dependency chain visualization
   - File and line number reporting

### **Phase 4: Advanced Features (Week 5-6)**

1. **Performance Optimization**
   - Efficient graph algorithms
   - Caching for large dependency graphs
   - Incremental dependency updates

2. **Advanced Validation**
   - Conditional dependencies
   - Dynamic dependency resolution
   - Cross-project dependencies

3. **Documentation and Testing**
   - Comprehensive test coverage
   - User documentation
   - Migration guide

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
  
  file "machines.hcl" {
    required = true
    validate = "hcl_machines_config"
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
  directory "actions" {
    required = false
    description = "Organized action files"
    validate = "hcl_actions_files"
  }
  
  directory "variables" {
    required = false
    description = "Organized variable files"
    validate = "hcl_variables_files"
  }
  
  # ... other directories
}
```

### **2. Actions Configuration Schema**

```hcl
# internal/schemas/schemas/actions.hcl
schema "actions_config" {
  block "actions" {
    required = true
    min_blocks = 0
    
    block "action" {
      required = true
      min_blocks = 0
      
      attribute "name" {
        type = "string"
        required = true
        description = "Action name"
      }
      
      attribute "command" {
        type = "string"
        required = true
        description = "Command to run"
      }
      
      attribute "description" {
        type = "string"
        required = false
        description = "Action description"
      }
      
      attribute "depends_on" {
        type = "list(string)"
        required = false
        description = "Dependencies on other actions"
      }
      
      attribute "timeout" {
        type = "number"
        required = false
        description = "Timeout in seconds"
        default = 300
      }
      
      attribute "retries" {
        type = "number"
        required = false
        description = "Number of retries"
        default = 1
      }
      
      attribute "parallel" {
        type = "bool"
        required = false
        description = "Run in parallel with other actions"
        default = false
      }
      
      attribute "tags" {
        type = "list(string)"
        required = false
        description = "Action tags"
      }
      
      attribute "machine" {
        type = "string"
        required = false
        description = "Specific machine to run on"
      }
      
      attribute "user" {
        type = "string"
        required = false
        description = "User to run as"
      }
      
      attribute "working_dir" {
        type = "string"
        required = false
        description = "Working directory"
      }
    }
  }
}
```

## Benefits

### **Immediate Benefits**
- **Centralized action management** with organized files
- **Dependency tracking** prevents run order issues
- **Cross-file validation** for merged configurations
- **Clear error messages** for action and dependency issues

### **Long-term Benefits**
- **Scalable action management** for large projects
- **Team collaboration** with organized action files
- **Automated validation** of complex action dependencies
- **Better debugging** with dependency visualization

## Success Metrics

### **Functionality Metrics**
- [ ] Actions can be defined in `actions.hcl` and `actions/` directory
- [ ] File merging works for both small and large files
- [ ] Circular dependencies are detected and reported
- [ ] Dependency resolution works for complex graphs
- [ ] Cross-file dependencies are properly tracked
- [ ] Merged files maintain dependency integrity
- [ ] Only project directory action files are read
- [ ] Invalid HCL in action files is treated as an error
- [ ] Clear error messages for file location and HCL validation issues
- [ ] `spooky actions validate` command validates all action files
- [ ] `spooky actions validate` provides clear, actionable error messages
- [ ] Command exits with appropriate exit codes (0 for success, non-zero for errors)

### **Performance Metrics**
- [ ] File merging handles files up to 2GB efficiently
- [ ] Memory usage stays reasonable for large projects
- [ ] Dependency resolution is fast for complex graphs
- [ ] Circular detection doesn't impact parsing performance

### **User Experience Metrics**
- [ ] Intuitive action definition syntax
- [ ] Good documentation and examples
- [ ] Backward compatibility with existing projects
- [ ] `spooky actions validate` provides clear, actionable error messages
- [ ] Command can be used in CI/CD pipelines for validation
- [ ] Success messages show summary of validated actions and files
- [ ] Clear error messages for dependency issues
- [ ] Dependency visualization tools

## Rationale: No Automatic Circular Dependency Resolution

### **1. We Have Automatic Dependency Resolution**

Spooky **does** provide automatic dependency resolution through topological sorting:

```go
// This works automatically - determines run order
func (dg *DependencyGraph) ResolveOrder() ([]string, error) {
    // Topological sort automatically resolves dependencies
    // Returns actions in correct run order
    // Handles complex dependency chains automatically
}
```

**What's Automatic:**
- **Execution order determination** via topological sorting
- **Dependency graph building** from merged files
- **Parallel running** of independent actions
- **Dependency validation** and reference checking

### **2. We Don't Fix Circular Dependencies**

Spooky **does not** automatically fix circular dependencies:

```hcl
# ❌ Circular dependency - spooky detects but doesn't fix
action "setup_podman" {
  command = "dnf install -y podman"
  depends_on = ["install_dependencies"]
}

action "install_dependencies" {
  command = "dnf update -y"
  depends_on = ["setup_podman"]  # Circular!
}
```

**What's Not Automatic:**
- **Removing circular dependencies** from configuration
- **Guessing user intent** for dependency relationships
- **Modifying configuration files** automatically
- **Breaking existing functionality** without user knowledge

### **3. Why No Automatic Circular Dependency Fixing**

#### **A. Dependency Intent is Business Logic**

```hcl
# Circular dependency could mean different things:
action "deploy_web" {
  depends_on = ["deploy_backend"]
}

action "deploy_backend" {
  depends_on = ["deploy_web"]  # Circular!

# Possible user intents:
# 1. "Web should start first, then backend"
# 2. "Backend should start first, then web" 
# 3. "They should start in parallel"
# 4. "I need a health check between them"
# 5. "I made a mistake, remove one dependency"

# Spooky cannot know which is correct!
```

#### **B. Breaking Changes Risk**

```bash
# Automatic "fix" could break functionality:
$ spooky act . --auto-fix-dependencies
Warning: Removing circular dependency: deploy_web -> deploy_backend

# Result: Web starts before backend is ready
# Result: Backend fails because web isn't available
# Result: User's application doesn't work
# Result: User blames spooky for "breaking" their deployment
```

#### **C. Complex Dependency Scenarios**

```hcl
# Real-world scenario with multiple circular chains:
action "setup_database" {
  depends_on = ["setup_network", "install_packages"]
}

action "setup_network" {
  depends_on = ["setup_database"]  # Needs DB for network config
}

action "install_packages" {
  depends_on = ["setup_network"]  # Needs network to download packages
}

# Automatic resolution would be disastrous:
# - Remove setup_database -> setup_network? Breaks network config
# - Remove setup_network -> install_packages? Breaks package installation
# - Remove setup_network -> setup_database? Breaks database setup
```

#### **D. User Ownership and Understanding**

```bash
# ❌ Automatic fix (user doesn't understand what happened):
$ spooky act . --auto-fix
Warning: Fixed circular dependencies automatically
✓ Actions run successfully

# User thinks: "Great! It worked!"
# Reality: User doesn't know what was changed or why
# Future: User repeats the same mistake

# ✅ Manual fix (user learns and understands):
$ spooky act .
Error: circular dependency detected: setup_podman -> install_dependencies -> setup_podman

# User thinks: "I need to understand my dependencies better"
# User learns: "Dependencies should form a directed acyclic graph"
# Future: User writes better dependencies
```

#### **E. Configuration as Code Philosophy**

```hcl
# Configuration should be explicit and version-controlled:
# actions/setup.hcl
action "setup_podman" {
  command = "dnf install -y podman"
  # User explicitly decides what this depends on
  # Changes are tracked in git
  # Changes are reviewed in pull requests
  # Changes are tested before deployment
}

# Automatic fixes would:
# - Change configuration without user knowledge
# - Make git diffs confusing
# - Hide real problems
# - Make debugging harder
```

#### **F. Different Fixes for Different Contexts**

```hcl
# Context 1: Development environment
action "start_dev_server" {
  depends_on = ["start_dev_db"]
}

action "start_dev_db" {
  depends_on = ["start_dev_server"]  # Circular!
}
# Fix: Remove dependency, run in parallel

# Context 2: Production environment  
action "deploy_prod_web" {
  depends_on = ["deploy_prod_db"]
}

action "deploy_prod_db" {
  depends_on = ["deploy_prod_web"]  # Circular!
}
# Fix: Add health checks, proper startup sequence

# Context 3: CI/CD pipeline
action "build_image" {
  depends_on = ["test_code"]
}

action "test_code" {
  depends_on = ["build_image"]  # Circular!
}
# Fix: Split into separate stages
```

### **4. Benefits of This Approach**

#### **A. User Education**
- **Users learn** about dependency management
- **Users understand** their configuration
- **Users write better** dependencies over time

#### **B. Predictable Behavior**
- **No surprises** from automatic changes
- **Configuration stays** as written
- **Git history** remains clean

#### **C. Business Logic Preservation**
- **User intent** is never overridden
- **Complex scenarios** are handled correctly
- **No breaking changes** from automatic fixes

#### **D. Debugging and Maintenance**
- **Clear error messages** point to real issues
- **Configuration is explicit** and understandable
- **Problems are visible** and fixable

### **5. Error Prevention vs. Error Correction**

```bash
# Spooky philosophy: Prevent errors, don't correct them
# - Clear validation messages
# - Early detection
# - Helpful suggestions
# - User education

# Not: "Let's fix this for you"
# But: "Here's what's wrong and how you can fix it"
```

This approach aligns perfectly with spooky's philosophy of being explicit, predictable, and educational rather than trying to be "smart" and potentially making things worse.

## Lessons from Ansible, Puppet, and Chef

### **1. Ansible Insights: Playbook Organization**

#### **A. Role-Based Organization (Adapted for Spooky)**
```hcl
# actions/roles/webserver.hcl
actions {
  action "install_nginx" {
    command = "dnf install -y nginx"
    tags = ["webserver", "install"]
  }
  
  action "configure_nginx" {
    command = "cp nginx.conf /etc/nginx/"
    depends_on = ["install_nginx"]
    tags = ["webserver", "configure"]
  }
  
  action "start_nginx" {
    command = "systemctl start nginx"
    depends_on = ["configure_nginx"]
    tags = ["webserver", "service"]
  }
}

# actions/roles/database.hcl
actions {
  action "install_postgresql" {
    command = "dnf install -y postgresql-server"
    tags = ["database", "install"]
  }
  
  action "init_database" {
    command = "postgresql-setup initdb"
    depends_on = ["install_postgresql"]
    tags = ["database", "init"]
  }
}
```

**Key Lessons:**
- **Tag-based organization** for selective running
- **Role-based file structure** for team collaboration
- **Clear separation** of install/configure/start phases

#### **B. Handler Pattern (Adapted for Spooky)**
```hcl
# actions/handlers.hcl
actions {
  action "restart_nginx" {
    command = "systemctl restart nginx"
    tags = ["handler", "nginx"]
    # Only runs when explicitly called or when config changes
  }
  
  action "reload_nginx" {
    command = "systemctl reload nginx"
    tags = ["handler", "nginx"]
    # Lighter than restart
  }
}

# actions/webserver.hcl
actions {
  action "update_nginx_config" {
    command = "cp new-nginx.conf /etc/nginx/nginx.conf"
    depends_on = ["reload_nginx"]  # Triggers handler
    tags = ["webserver", "config"]
  }
}
```

### **2. Puppet Insights: Resource Dependencies**

#### **A. Explicit Resource Dependencies**
```hcl
# actions/resources.hcl
actions {
  action "create_user" {
    command = "useradd -m -s /bin/bash appuser"
    tags = ["user", "create"]
  }
  
  action "create_directory" {
    command = "mkdir -p /home/appuser/app"
    depends_on = ["create_user"]  # Explicit dependency
    tags = ["directory", "create"]
  }
  
  action "set_permissions" {
    command = "chown -R appuser:appuser /home/appuser/app"
    depends_on = ["create_directory"]
    tags = ["permissions", "set"]
  }
}
```

**Key Lessons:**
- **Explicit dependencies** prevent race conditions
- **Resource ordering** is critical for system consistency
- **Clear resource lifecycle** (create → configure → start)

#### **B. Idempotency Pattern**
```hcl
# actions/idempotent.hcl
actions {
  action "install_package" {
    command = "dnf install -y nginx || true"  # Idempotent
    tags = ["package", "install"]
  }
  
  action "create_config" {
    command = "test -f /etc/nginx/nginx.conf || cp nginx.conf /etc/nginx/"
    tags = ["config", "create"]
  }
  
  action "ensure_service_running" {
    command = "systemctl is-active nginx || systemctl start nginx"
    depends_on = ["install_package", "create_config"]
    tags = ["service", "ensure"]
  }
}
```

### **3. Chef Insights: Recipe Organization**

#### **A. Recipe-Style Organization**
```hcl
# actions/recipes/default.hcl
actions {
  action "update_system" {
    command = "dnf update -y"
    tags = ["system", "update"]
  }
  
  action "install_dependencies" {
    command = "dnf install -y git curl wget"
    depends_on = ["update_system"]
    tags = ["dependencies", "install"]
  }
}

# actions/recipes/application.hcl
actions {
  action "clone_repository" {
    command = "git clone https://github.com/user/app.git /opt/app"
    depends_on = ["install_dependencies"]
    tags = ["app", "clone"]
  }
  
  action "install_app_dependencies" {
    command = "cd /opt/app && npm install"
    depends_on = ["clone_repository"]
    tags = ["app", "dependencies"]
  }
}
```

#### **B. Environment-Specific Actions**
```hcl
# actions/environments/development.hcl
actions {
  action "setup_dev_environment" {
    command = "echo 'Setting up development environment'"
    tags = ["dev", "setup"]
  }
  
  action "start_dev_services" {
    command = "docker-compose up -d"
    depends_on = ["setup_dev_environment"]
    tags = ["dev", "services"]
  }
}

# actions/environments/production.hcl
actions {
  action "setup_prod_environment" {
    command = "echo 'Setting up production environment'"
    tags = ["prod", "setup"]
  }
  
  action "deploy_prod_services" {
    command = "kubectl apply -f k8s/"
    depends_on = ["setup_prod_environment"]
    tags = ["prod", "deploy"]
  }
}
```

### **4. Cross-Tool Best Practices**

#### **A. State Management Pattern**
```hcl
# actions/state.hcl
actions {
  action "check_current_state" {
    command = "systemctl is-active nginx"
    tags = ["state", "check"]
    # Returns 0 if running, non-zero if not
  }
  
  action "ensure_desired_state" {
    command = "systemctl start nginx"
    depends_on = ["check_current_state"]
    tags = ["state", "ensure"]
    # Only runs if check_current_state fails
  }
}
```

#### **B. Rollback Pattern**
```hcl
# actions/rollback.hcl
actions {
  action "backup_current_config" {
    command = "cp /etc/nginx/nginx.conf /etc/nginx/nginx.conf.backup"
    tags = ["rollback", "backup"]
  }
  
  action "apply_new_config" {
    command = "cp new-nginx.conf /etc/nginx/nginx.conf"
    depends_on = ["backup_current_config"]
    tags = ["rollback", "apply"]
  }
  
  action "rollback_if_failed" {
    command = "cp /etc/nginx/nginx.conf.backup /etc/nginx/nginx.conf"
    tags = ["rollback", "restore"]
    # Only runs if apply_new_config fails
  }
}
```

#### **C. Health Check Pattern**
```hcl
# actions/health.hcl
actions {
  action "health_check_web" {
    command = "curl -f http://localhost:80/health || exit 1"
    tags = ["health", "web"]
    retries = 3
    timeout = 30
  }
  
  action "health_check_db" {
    command = "pg_isready -h localhost || exit 1"
    tags = ["health", "database"]
    retries = 3
    timeout = 30
  }
  
  action "full_health_check" {
    command = "echo 'All systems operational'"
    depends_on = ["health_check_web", "health_check_db"]
    tags = ["health", "full"]
  }
}
```

### **5. Anti-Patterns to Avoid**

#### **A. Don't: Complex Conditional Logic**
```hcl
# ❌ Avoid complex conditionals in actions
action "complex_deploy" {
  command = "if [ $ENV = 'prod' ]; then deploy_prod; else deploy_dev; fi"
  # Hard to test, debug, and maintain
}

# ✅ Use separate actions with tags
action "deploy_prod" {
  command = "deploy_prod.sh"
  tags = ["deploy", "prod"]
}

action "deploy_dev" {
  command = "deploy_dev.sh"
  tags = ["deploy", "dev"]
}
```

#### **B. Don't: Long-Running Actions**
```hcl
# ❌ Avoid very long-running actions
action "long_deploy" {
  command = "sleep 3600 && deploy_app"  # 1 hour timeout
  timeout = 7200
}

# ✅ Break into smaller, manageable actions
action "prepare_deploy" {
  command = "prepare_deployment.sh"
  timeout = 300
}

action "deploy_app" {
  command = "deploy_application.sh"
  depends_on = ["prepare_deploy"]
  timeout = 600
}

action "verify_deploy" {
  command = "verify_deployment.sh"
  depends_on = ["deploy_app"]
  timeout = 300
}
```

#### **C. Don't: Hidden Dependencies**
```hcl
# ❌ Avoid hidden dependencies in commands
action "start_app" {
  command = "cd /opt/app && ./start.sh"  # Depends on app being installed
}

# ✅ Make dependencies explicit
action "install_app" {
  command = "install_application.sh"
  tags = ["app", "install"]
}

action "start_app" {
  command = "cd /opt/app && ./start.sh"
  depends_on = ["install_app"]  # Explicit dependency
  tags = ["app", "start"]
}
```

### **6. Spooky-Specific Enhancements**

#### **A. Template Integration**
```hcl
# actions/templated.hcl
actions {
  action "generate_config" {
    command = "spooky render-template templates/nginx.conf.tmpl /etc/nginx/nginx.conf"
    tags = ["config", "generate"]
  }
  
  action "reload_nginx" {
    command = "systemctl reload nginx"
    depends_on = ["generate_config"]
    tags = ["service", "reload"]
  }
}
```

#### **B. Facts Integration**
```hcl
# actions/facts-aware.hcl
actions {
  action "install_package" {
    command = "dnf install -y ${var.package_name}"
    tags = ["package", "install"]
    # Uses variables that can be set from facts
  }
  
  action "configure_service" {
    command = "cp config-${var.environment}.conf /etc/service/"
    depends_on = ["install_package"]
    tags = ["service", "configure"]
  }
}
```

These insights from Ansible, Puppet, and Chef provide valuable patterns and anti-patterns that can guide spooky's actions system design and help users create more maintainable and reliable automation.

## Risk Assessment

### **Technical Risks**
- **Large file handling** - Mitigation: Implement size thresholds and temporary files
- **Performance impact** - Mitigation: Efficient algorithms and caching
- **Path traversal attacks** - Mitigation: Strict path validation and project directory restrictions
- **Invalid HCL files** - Mitigation: Strict HCL validation with clear error reporting
- **Circular dependency complexity** - Mitigation: Clear error reporting and visualization
- **Cross-file complexity** - Mitigation: Incremental implementation and testing

### **User Experience Risks**
- **Learning curve** - Mitigation: Good documentation and examples
- **Migration complexity** - Mitigation: Backward compatibility and migration tools
- **Error debugging** - Mitigation: Clear error messages and dependency chains

### **Implementation Risks**
- **Scope creep** - Mitigation: Focus on core functionality first
- **Testing complexity** - Mitigation: Comprehensive test scenarios
- **Integration issues** - Mitigation: Incremental implementation and testing 