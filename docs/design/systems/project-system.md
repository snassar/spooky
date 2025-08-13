# Project System: Comprehensive Implementation Plan

## Overview

This document outlines the implementation of the spooky project system, which provides isolated, portable, and configurable environments for infrastructure automation. Projects integrate with the facts system and configuration while maintaining their own identity and settings.

**Schema Integration**: This project system implements the schema validation patterns and project structure definitions defined in [Schema System](../schema-system.md) for comprehensive project validation, metadata schema enforcement, and schema-based project lifecycle management.

## System Integration

This project system integrates with other core Spooky systems to provide isolated, portable, and configurable project environments:

### **Facts System Integration**
- **Project Facts**: Project-specific facts database and collection settings (see [Facts System](../facts-system.md))
- **Project Isolation**: Each project maintains its own facts database
- **Project Context**: Facts collection within project context and configuration
- **Project Validation**: Facts validation integrated with project validation
- **Facts Storage**: Project-specific facts storage and caching

### **Variables System Integration**
- **Variables Management**: Project variables stored in `variables.hcl` and `variables/` directory (see [Variables System](../variables-system.md))
- **Variable Resolution**: Project variables resolved in project context
- **Template Integration**: Variables available in project templates
- **Configuration Precedence**: Project variables override global defaults

### **Schema System Integration**
- **Project Schema**: Project configuration validated against embedded schemas (see [Schema System](../schema-system.md))
- **Structure Validation**: Project directory structure validated against schema
- **Metadata Schema**: Project metadata enforced through schema validation
- **Schema Evolution**: Project schemas evolve with system changes

### **CLI System Integration**
- **Project Commands**: Project management through `spooky project` commands (see [CLI System](../cli-system.md))
- **Project Initialization**: `spooky project init` for project creation
- **Project Validation**: `spooky project validate` for project structure validation
- **Project Information**: `spooky project show` for project information display
- **Project Discovery**: Project discovery using standard Unix tools

### **Configuration System Integration**
- **Project Configuration**: Project-specific configuration overrides (see [Configuration System](../configuration-system.md))
- **Project Isolation**: Project configuration isolation from global settings
- **Project Defaults**: Default configuration values for new projects
- **Project Validation**: Configuration validation for project settings
- **Project Context**: Configuration resolution within project context

### **Machines System Integration**
- **Project Inventory**: Machine inventory stored in project-specific `machines.hcl` files (see [Machines System](../machines-system.md))
- **Project Isolation**: Each project maintains its own machine inventory
- **Project Context**: Machine operations run within project context
- **Project Validation**: Machine inventory validation integrated with project validation
- **Enterprise Scale**: Project supports large machine inventories with efficient indexing

### **Actions System Integration**
- **Project Actions**: Actions stored in project-specific `actions.hcl` and `actions/` directory (see [Actions System](../actions-system.md))
- **Project Context**: Actions run within project run context
- **Project Configuration**: Action settings configured in project context
- **Project Isolation**: Project actions isolated from global actions
- **Dependency Management**: Project-level action dependency tracking

### **Template System Integration**
- **Project Templates**: Templates stored in project-specific `templates/` directories (see [Template System](../template-system.md))
- **Project Context**: Templates run within project context and configuration
- **Project Data**: Templates have access to project-specific data and configuration
- **Project Validation**: Template validation integrated with project validation
- **Template Deployment**: Project supports template deployment workflows

## Current State Analysis

### **What We Have**
- ✅ **Project initialization** with `spooky project init`
- ✅ **Basic project structure** (project.hcl, machines.hcl, actions.hcl)
- ✅ **Template system** with facts integration
- ✅ **SSH client** for remote running
- ✅ **HCL parsing** for configuration files

### **What We Need**
- 🔄 **Enhanced project schema** with validation
- 🔄 **Project isolation** with global integration
- 🔄 **Portable project names** for cross-environment use
- 🔄 **Project-specific configuration** with global defaults
- 🔄 **Project metadata** and versioning
- 🔄 **Project dependencies** and imports

## Project System Design

### **1. Project Structure and Identity**

**Standard Project Structure**:
```
project-root/
├── project.hcl                   # Main project configuration
├── machines.hcl                  # Machine inventory (optional)
├── actions.hcl                   # Action definitions (optional)
├── variables.hcl                 # Project variables (optional)
├── machines/                     # Machine inventory files (optional)
├── actions/                      # Action definition files (optional)
├── variables/                    # Variable definition files (optional)
├── templates/                    # Template files (optional)
├── files/                        # Static files for deployment (optional)
├── logs/                         # Log files directory (optional)
└── README.md                     # Project documentation (optional)
```

**Project Identity and Metadata**:
```hcl
# project.hcl
project "nextcloud-production" {
  description = "Nextcloud production deployment"
  version = "1.2.0"
  environment = "production"
  
  # Machine inventory configuration
  machine_inventory = "machines.hcl"
  
  # Actions configuration
  actions_file = "actions.hcl"
  
  # Default run settings
  default_timeout = 300
  parallel = 10
  
  # Facts storage configuration
  facts_storage {
    type = "memory"
    path = "memory"
  }
  
  # Logging configuration
  logging {
    level = "info"
    format = "text"
    output = "logs/spooky.log"
  }
  
  # SSH configuration
  ssh {
    default_user = "ubuntu"
    default_port = 22
    connection_timeout = 30
    command_timeout = 300
    retry_attempts = 3
    retry_delay = 5
    key_path = "~/.ssh/id_rsa"
    strict_host_key_checking = true
    allow_password_auth = false
  }
  
  # Templates configuration
  templates {
    data_directory = "data"
    templates_directory = "templates"
    auto_load_data = true
    strict_validation = true
  }
  
  # Security configuration
  security {
    allow_insecure_connections = false
    allow_password_auth = false
    require_https_imports = true
    max_file_size = "100MB"
  }
}
```

### **2. Project Configuration Schema**

```go
type Project struct {
	Path   string
	Config *ProjectConfig
	// Facts are stored in memory during export operations
}
```

### **3. Configuration Precedence and Inheritance**

**Configuration Precedence Hierarchy** (highest to lowest priority):

1. **CLI flags** - `--ssh-timeout 60s`
2. **Project configuration** - `project.hcl`
3. **Global configuration** - `$XDG_CONFIG_HOME/spooky/spooky.hcl`
4. **Environment variables** - `SPOOKY_SSH_TIMEOUT=45s`
5. **Default values** - Built-in defaults

```go
// Configuration inheritance and precedence
func loadProjectConfiguration(projectPath string) (*Project, error) {
    // 1. Load global configuration (base)
    globalConfig, err := loadGlobalConfiguration()
    if err != nil {
        return nil, fmt.Errorf("failed to load global config: %w", err)
    }
    
    // 2. Load project configuration
    projectConfig, err := loadProjectFile(projectPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load project config: %w", err)
    }
    
    // 3. Apply environment variable overrides
    applyEnvironmentOverrides(projectConfig)
    
    // 4. Apply CLI flag overrides
    applyCLIOverrides(projectConfig)
    
    return projectConfig, nil
}

// Apply environment variable overrides to project configuration
func applyEnvironmentOverrides(project *Project) {
    // SSH configuration overrides
    if sshTimeout := os.Getenv("SPOOKY_SSH_TIMEOUT"); sshTimeout != "" {
        if timeout, err := strconv.Atoi(sshTimeout); err == nil {
            if project.SSH == nil {
                project.SSH = &ProjectSSH{}
            }
            project.SSH.ConnectionTimeout = timeout
        }
    }
    
    // Parallel run override
    if parallel := os.Getenv("SPOOKY_PARALLEL"); parallel != "" {
        if p, err := strconv.Atoi(parallel); err == nil {
            project.Parallel = p
        }
    }
    
    // Default timeout override
    if timeout := os.Getenv("SPOOKY_DEFAULT_TIMEOUT"); timeout != "" {
        if t, err := strconv.Atoi(timeout); err == nil {
            project.DefaultTimeout = t
        }
    }
}
```

### **4. Project Isolation and Global Integration**

**Project Isolation**:
- **File system isolation**: Projects can only access files within their directory
- **Configuration isolation**: Project-specific settings override global defaults
- **Facts isolation**: Project facts are stored in project-specific `facts.db` directory
- **Template isolation**: Templates can only access project-specific data

**Global Integration**:
- **Global configuration**: Default settings come from global config
- **Global SSH settings**: Default SSH behavior from global config
- **Environment variables**: Can override project settings

```go
// Project context with isolation and global integration
type ProjectContext struct {
    Project      *Project
    GlobalConfig *GlobalConfig
    Path         string
    Name         string
    
    // Project-specific components
    FactsDB      *BadgerDB
    Templates    *TemplateManager
    Actions      *ActionManager
    
    // Global integration
    SSHClient    *SSHClient
    Logger       *Logger
}

// Project facts manager with project-specific export capabilities
type ProjectFactsManager struct {
    projectPath string
    config      *ProjectFactsConfig
}

// Integrates with the facts system (see [Facts System](../facts-system.md)) for project-specific fact gathering and export
func (p *ProjectFactsManager) ExportFacts(machines []*Machine, format string, outputPath string) error {
    // Gather and export facts directly from machines
    return p.factsManager.ExportFacts(context.Background(), machines, format, outputPath)
}

func (p *ProjectFactsManager) GetStorageStats() (map[string]interface{}, error) {
    // Get storage statistics for debugging
    return p.factsManager.GetStorageStats()
}
```

### **5. Project Portability**

**Portable Project Names**:
- **Unique identification**: Projects are identified by name across environments
- **Version tracking**: Projects have versions for tracking changes
- **Cross-environment**: Same project can run in dev, staging, production
- **Team collaboration**: Projects can be shared between teams

**Project Identity**:
```go
// Project identity and portability
type ProjectIdentity struct {
    Name        string
    Version     string
    Environment string
    Description string
}

func (p *ProjectIdentity) IsPortable() bool {
    // Check if project has required identity fields
    return p.Name != "" && p.Version != ""
}

func (p *ProjectIdentity) GetFullName() string {
    if p.Version != "" {
        return fmt.Sprintf("%s@%s", p.Name, p.Version)
    }
    return p.Name
}
```

### **6. Project Validation**

**Project Structure Validation**:
```go
// Project validator based on actual schema
type ProjectValidator struct {
    project *Project
    path    string
}

func (p *ProjectValidator) Validate() error {
    var errors []string
    
    // 1. Validate project.hcl structure
    if err := p.validateProjectConfig(); err != nil {
        errors = append(errors, fmt.Sprintf("project config: %v", err))
    }
    
    // 2. Validate project directory structure
    if err := p.validateDirectoryStructure(); err != nil {
        errors = append(errors, fmt.Sprintf("directory structure: %v", err))
    }
    
    if len(errors) > 0 {
        return fmt.Errorf("project validation failed:\n  %s", strings.Join(errors, "\n  "))
    }
    
    return nil
}

func (p *ProjectValidator) validateProjectConfig() error {
    // Validate required fields
    if p.project.Name == "" {
        return fmt.Errorf("project name is required")
    }
    
    // Validate optional fields if present
    if p.project.Parallel < 2 {
        return fmt.Errorf("parallel must be 2 or greater")
    }
    
    if p.project.DefaultTimeout < 1 {
        return fmt.Errorf("default_timeout must be 1 or greater")
    }
    
    return nil
}

func (p *ProjectValidator) validateDirectoryStructure() error {
    // Check for optional directories
    optionalDirs := []string{"templates", "files", "logs"}
    for _, dir := range optionalDirs {
        dirPath := filepath.Join(p.path, dir)
        if _, err := os.Stat(dirPath); err == nil {
            // Directory exists, validate it's actually a directory
            if info, err := os.Stat(dirPath); err == nil && !info.IsDir() {
                return fmt.Errorf("%s exists but is not a directory", dir)
            }
        }
    }
    
    return nil
}
```