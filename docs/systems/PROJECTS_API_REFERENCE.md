# Projects API Reference

## Overview

The Projects API provides comprehensive project management capabilities for the spooky codebase. This reference documents all interfaces, types, and methods available for project initialization, configuration management, validation, and lifecycle management.

## Core Interfaces

### ProjectIntegration Interface

The main interface for project integration operations.

```go
type ProjectIntegration interface {
    LoadProject(ctx context.Context, projectPath string) (*spookytypes.Project, error)
    ValidateProject(ctx context.Context, project *spookytypes.Project) (*spookytypes.ValidationResult, error)
    InitializeProject(ctx context.Context, config *spookytypes.ProjectConfig) (*spookytypes.Project, error)
    GetProjectInfo(ctx context.Context, projectPath string) (*spookytypes.ProjectInfo, error)
}
```

#### Methods

##### LoadProject
Loads a project configuration from the specified path.

```go
func (i *ProjectIntegration) LoadProject(ctx context.Context, projectPath string) (*spookytypes.Project, error)
```

**Parameters:**
- `ctx context.Context` - Context for the operation
- `projectPath string` - Path to the project directory

**Returns:**
- `*spookytypes.Project` - Loaded project configuration
- `error` - Error if loading fails

**Example:**
```go
project, err := projectIntegration.LoadProject(ctx, "./my-project")
if err != nil {
    return fmt.Errorf("failed to load project: %w", err)
}
```

##### ValidateProject
Validates a project configuration and structure.

```go
func (i *ProjectIntegration) ValidateProject(ctx context.Context, project *spookytypes.Project) (*spookytypes.ValidationResult, error)
```

**Parameters:**
- `ctx context.Context` - Context for the operation
- `project *spookytypes.Project` - Project to validate

**Returns:**
- `*spookytypes.ValidationResult` - Validation results
- `error` - Error if validation fails

**Example:**
```go
result, err := projectIntegration.ValidateProject(ctx, project)
if err != nil {
    return fmt.Errorf("failed to validate project: %w", err)
}

if !result.IsValid() {
    for _, error := range result.Errors {
        log.Printf("Validation error: %s", error.Message)
    }
    return fmt.Errorf("project validation failed")
}
```

##### InitializeProject
Initializes a new project with the specified configuration.

```go
func (i *ProjectIntegration) InitializeProject(ctx context.Context, config *spookytypes.ProjectConfig) (*spookytypes.Project, error)
```

**Parameters:**
- `ctx context.Context` - Context for the operation
- `config *spookytypes.ProjectConfig` - Project configuration

**Returns:**
- `*spookytypes.Project` - Initialized project
- `error` - Error if initialization fails

**Example:**
```go
config := &spookytypes.ProjectConfig{
    Project: &spookytypes.ProjectMetadata{
        Name: "my-project",
        Description: "A sample project",
        Version: "1.0.0",
        Author: "spooky-user",
    },
    Settings: &spookytypes.ProjectSettings{
        ParallelWorkers: 4,
        TimeoutSeconds: 300,
        LogLevel: "info",
    },
}

project, err := projectIntegration.InitializeProject(ctx, config)
if err != nil {
    return fmt.Errorf("failed to initialize project: %w", err)
}
```

##### GetProjectInfo
Retrieves information about a project.

```go
func (i *ProjectIntegration) GetProjectInfo(ctx context.Context, projectPath string) (*spookytypes.ProjectInfo, error)
```

**Parameters:**
- `ctx context.Context` - Context for the operation
- `projectPath string` - Path to the project directory

**Returns:**
- `*spookytypes.ProjectInfo` - Project information
- `error` - Error if retrieval fails

**Example:**
```go
info, err := projectIntegration.GetProjectInfo(ctx, "./my-project")
if err != nil {
    return fmt.Errorf("failed to get project info: %w", err)
}

log.Printf("Project: %s", info.Name)
log.Printf("Version: %s", info.Version)
log.Printf("Status: %s", info.Status)
```

## Project Manager Interface

The project manager provides high-level project management operations.

```go
type ProjectManager interface {
    LoadProject(ctx context.Context, projectPath string) (*spookytypes.Project, error)
    ValidateProject(ctx context.Context, project *spookytypes.Project) (*spookytypes.ValidationResult, error)
    InitializeProject(ctx context.Context, config *spookytypes.ProjectConfig) (*spookytypes.Project, error)
    GetProjectInfo(ctx context.Context, projectPath string) (*spookytypes.ProjectInfo, error)
    GetProjectStatus(ctx context.Context, projectPath string) (*spookytypes.ProjectStatus, error)
    GetProjectConfig(ctx context.Context, projectPath string) (*spookytypes.ProjectConfig, error)
    StartProject(ctx context.Context, projectPath string) error
    StopProject(ctx context.Context, projectPath string) error
    RestartProject(ctx context.Context, projectPath string) error
    CleanProject(ctx context.Context, projectPath string) error
}
```

### Additional Methods

#### GetProjectStatus
Retrieves the current status of a project.

```go
func (m *ProjectManager) GetProjectStatus(ctx context.Context, projectPath string) (*spookytypes.ProjectStatus, error)
```

**Parameters:**
- `ctx context.Context` - Context for the operation
- `projectPath string` - Path to the project directory

**Returns:**
- `*spookytypes.ProjectStatus` - Project status information
- `error` - Error if retrieval fails

**Example:**
```go
status, err := projectManager.GetProjectStatus(ctx, "./my-project")
if err != nil {
    return fmt.Errorf("failed to get project status: %w", err)
}

log.Printf("Project status: %s", status.Status)
log.Printf("Last updated: %s", status.LastUpdated)
log.Printf("Health: %s", status.Health)
```

#### GetProjectConfig
Retrieves the configuration of a project.

```go
func (m *ProjectManager) GetProjectConfig(ctx context.Context, projectPath string) (*spookytypes.ProjectConfig, error)
```

**Parameters:**
- `ctx context.Context` - Context for the operation
- `projectPath string` - Path to the project directory

**Returns:**
- `*spookytypes.ProjectConfig` - Project configuration
- `error` - Error if retrieval fails

**Example:**
```go
config, err := projectManager.GetProjectConfig(ctx, "./my-project")
if err != nil {
    return fmt.Errorf("failed to get project config: %w", err)
}

log.Printf("Project name: %s", config.Project.Name)
log.Printf("Parallel workers: %d", config.Settings.ParallelWorkers)
log.Printf("Log level: %s", config.Settings.LogLevel)
```

#### StartProject
Starts a project and its associated components.

```go
func (m *ProjectManager) StartProject(ctx context.Context, projectPath string) error
```

**Parameters:**
- `ctx context.Context` - Context for the operation
- `projectPath string` - Path to the project directory

**Returns:**
- `error` - Error if starting fails

**Example:**
```go
err := projectManager.StartProject(ctx, "./my-project")
if err != nil {
    return fmt.Errorf("failed to start project: %w", err)
}

log.Printf("Project started successfully")
```

#### StopProject
Stops a project and its associated components.

```go
func (m *ProjectManager) StopProject(ctx context.Context, projectPath string) error
```

**Parameters:**
- `ctx context.Context` - Context for the operation
- `projectPath string` - Path to the project directory

**Returns:**
- `error` - Error if stopping fails

**Example:**
```go
err := projectManager.StopProject(ctx, "./my-project")
if err != nil {
    return fmt.Errorf("failed to stop project: %w", err)
}

log.Printf("Project stopped successfully")
```

#### RestartProject
Restarts a project and its associated components.

```go
func (m *ProjectManager) RestartProject(ctx context.Context, projectPath string) error
```

**Parameters:**
- `ctx context.Context` - Context for the operation
- `projectPath string` - Path to the project directory

**Returns:**
- `error` - Error if restarting fails

**Example:**
```go
err := projectManager.RestartProject(ctx, "./my-project")
if err != nil {
    return fmt.Errorf("failed to restart project: %w", err)
}

log.Printf("Project restarted successfully")
```

#### CleanProject
Cleans up project resources and temporary files.

```go
func (m *ProjectManager) CleanProject(ctx context.Context, projectPath string) error
```

**Parameters:**
- `ctx context.Context` - Context for the operation
- `projectPath string` - Path to the project directory

**Returns:**
- `error` - Error if cleaning fails

**Example:**
```go
err := projectManager.CleanProject(ctx, "./my-project")
if err != nil {
    return fmt.Errorf("failed to clean project: %w", err)
}

log.Printf("Project cleaned successfully")
```

## Core Types

### Project

Represents a complete project configuration.

```go
type Project struct {
    ID              string                 // Project identifier
    Name            string                 // Project name
    Description     string                 // Project description
    Path            string                 // Project path
    Type            string                 // Project type
    Version         string                 // Project version
    Author          string                 // Project author
    Tags            []string               // Project tags
    Metadata        map[string]interface{} // Project metadata
    Configuration   *ProjectConfig         // Project configuration
    Status          string                 // Project status
    CreatedAt       time.Time              // Creation timestamp
    UpdatedAt       time.Time              // Update timestamp
}
```

### ProjectConfig

Represents the configuration of a project.

```go
type ProjectConfig struct {
    Project         *ProjectMetadata       // Project metadata
    Settings        *ProjectSettings       // Project settings
    Components      *ProjectComponents     // Project components
    Dependencies    *ProjectDependencies   // Project dependencies
    Environment     *ProjectEnvironment    // Project environment
    Security        *ProjectSecurity       // Project security
}
```

### ProjectMetadata

Represents metadata about a project.

```go
type ProjectMetadata struct {
    Name            string                 // Project name
    Description     string                 // Project description
    Version         string                 // Project version
    Author          string                 // Project author
    License         string                 // Project license
    Tags            []string               // Project tags
    Repository      string                 // Repository URL
    Documentation   string                 // Documentation URL
}
```

### ProjectSettings

Represents project settings and configuration.

```go
type ProjectSettings struct {
    ParallelWorkers int                    // Number of parallel workers
    TimeoutSeconds  int                    // Operation timeout in seconds
    LogLevel        string                 // Logging level
    DebugMode       bool                   // Debug mode flag
    CacheEnabled    bool                   // Cache enabled flag
    CacheTTL        string                 // Cache TTL
}
```

### ProjectComponents

Represents project component configuration.

```go
type ProjectComponents struct {
    Facts           *ComponentConfig       // Facts component
    Actions         *ComponentConfig       // Actions component
    Machines        *ComponentConfig       // Machines component
    Templates       *ComponentConfig       // Templates component
    Variables       *ComponentConfig       // Variables component
    Secrets         *ComponentConfig       // Secrets component
}
```

### ComponentConfig

Represents configuration for a project component.

```go
type ComponentConfig struct {
    Enabled         bool                   // Component enabled flag
    Config          map[string]interface{} // Component configuration
    Metadata        map[string]interface{} // Component metadata
}
```

### ProjectEnvironment

Represents project environment configuration.

```go
type ProjectEnvironment struct {
    Development     *EnvironmentConfig     // Development environment
    Staging         *EnvironmentConfig     // Staging environment
    Production      *EnvironmentConfig     // Production environment
    Testing         *EnvironmentConfig     // Testing environment
}
```

### EnvironmentConfig

Represents configuration for a project environment.

```go
type EnvironmentConfig struct {
    Enabled         bool                   // Environment enabled flag
    Machines        []string               // Environment machines
    Variables       []string               // Environment variables
    SecurityLevel   string                 // Security level
    BackupEnabled   bool                   // Backup enabled flag
    Metadata        map[string]interface{} // Environment metadata
}
```

### ProjectSecurity

Represents project security configuration.

```go
type ProjectSecurity struct {
    EncryptSensitiveData    bool                   // Encrypt sensitive data flag
    ValidateFilePermissions bool                   // Validate file permissions flag
    AuditAllOperations      bool                   // Audit all operations flag
    AccessControlEnabled    bool                   // Access control enabled flag
    AllowedUsers           []string               // Allowed users
    AllowedGroups          []string               // Allowed groups
    Encryption             *EncryptionConfig      // Encryption configuration
    Compliance             *ComplianceConfig      // Compliance configuration
}
```

### EncryptionConfig

Represents encryption configuration.

```go
type EncryptionConfig struct {
    Method              string                 // Encryption method
    KeyPath             string                 // Key path
    BackupKeys          bool                   // Backup keys flag
    RotationInterval    string                 // Key rotation interval
    Metadata            map[string]interface{} // Encryption metadata
}
```

### ComplianceConfig

Represents compliance configuration.

```go
type ComplianceConfig struct {
    SOXEnabled       bool                   // SOX compliance enabled
    PCIEnabled       bool                   // PCI compliance enabled
    GDPREnabled      bool                   // GDPR compliance enabled
    AuditRetention   string                 // Audit retention period
    Metadata         map[string]interface{} // Compliance metadata
}
```

### ProjectInfo

Represents information about a project.

```go
type ProjectInfo struct {
    Name            string                 // Project name
    Version         string                 // Project version
    Status          string                 // Project status
    Path            string                 // Project path
    Type            string                 // Project type
    Author          string                 // Project author
    CreatedAt       time.Time              // Creation timestamp
    UpdatedAt       time.Time              // Update timestamp
    Metadata        map[string]interface{} // Project metadata
}
```

### ProjectStatus

Represents the status of a project.

```go
type ProjectStatus struct {
    Status          string                 // Project status
    Health          string                 // Project health
    LastUpdated     time.Time              // Last update timestamp
    Components      map[string]string      // Component status
    Errors          []string               // Error messages
    Warnings        []string               // Warning messages
    Metadata        map[string]interface{} // Status metadata
}
```

## Error Types

### ProjectError

Represents errors specific to project operations.

```go
type ProjectError struct {
    Code            string                 // Error code
    Message         string                 // Error message
    Operation       string                 // Operation that failed
    ProjectPath     string                 // Project path
    Details         map[string]interface{} // Error details
    Timestamp       time.Time              // Error timestamp
}
```

### ValidationError

Represents validation errors for project configuration.

```go
type ValidationError struct {
    Field           string                 // Field that failed validation
    Message         string                 // Validation message
    Value           interface{}            // Invalid value
    Constraint      string                 // Constraint that failed
    Severity        string                 // Error severity
    Metadata        map[string]interface{} // Validation metadata
}
```

## Constants

### Project Status Constants

```go
const (
    ProjectStatusInitializing = "initializing"
    ProjectStatusActive       = "active"
    ProjectStatusStopped      = "stopped"
    ProjectStatusError        = "error"
    ProjectStatusMaintenance  = "maintenance"
)
```

### Project Health Constants

```go
const (
    ProjectHealthHealthy   = "healthy"
    ProjectHealthWarning   = "warning"
    ProjectHealthCritical  = "critical"
    ProjectHealthUnknown   = "unknown"
)
```

### Project Type Constants

```go
const (
    ProjectTypeBasic        = "basic"
    ProjectTypeWebApp       = "web-app"
    ProjectTypeDatabase     = "database"
    ProjectTypeInfrastructure = "infrastructure"
    ProjectTypeMonitoring   = "monitoring"
    ProjectTypeSecurity     = "security"
)
```

### Security Level Constants

```go
const (
    SecurityLevelLow      = "low"
    SecurityLevelMedium   = "medium"
    SecurityLevelHigh     = "high"
    SecurityLevelMaximum  = "maximum"
)
```

## Usage Examples

### Basic Project Management

```go
// Load and validate a project
project, err := projectIntegration.LoadProject(ctx, "./my-project")
if err != nil {
    return fmt.Errorf("failed to load project: %w", err)
}

result, err := projectIntegration.ValidateProject(ctx, project)
if err != nil {
    return fmt.Errorf("failed to validate project: %w", err)
}

if !result.IsValid() {
    return fmt.Errorf("project validation failed")
}

// Get project information
info, err := projectIntegration.GetProjectInfo(ctx, "./my-project")
if err != nil {
    return fmt.Errorf("failed to get project info: %w", err)
}

log.Printf("Project: %s (v%s)", info.Name, info.Version)
log.Printf("Status: %s", info.Status)
```

### Project Initialization

```go
// Create project configuration
config := &spookytypes.ProjectConfig{
    Project: &spookytypes.ProjectMetadata{
        Name: "my-web-app",
        Description: "A web application project",
        Version: "1.0.0",
        Author: "spooky-user",
        License: "MIT",
        Tags: []string{"web", "application", "production"},
    },
    Settings: &spookytypes.ProjectSettings{
        ParallelWorkers: 8,
        TimeoutSeconds: 600,
        LogLevel: "info",
        DebugMode: false,
        CacheEnabled: true,
        CacheTTL: "1h",
    },
    Components: &spookytypes.ProjectComponents{
        Facts: &spookytypes.ComponentConfig{
            Enabled: true,
            Config: map[string]interface{}{
                "collection_interval": "30m",
                "storage_backend": "badgerdb",
            },
        },
        Actions: &spookytypes.ComponentConfig{
            Enabled: true,
            Config: map[string]interface{}{
                "execution_mode": "parallel",
                "max_retries": 3,
            },
        },
        Machines: &spookytypes.ComponentConfig{
            Enabled: true,
            Config: map[string]interface{}{
                "connection_pool_size": 10,
                "health_check_interval": "5m",
            },
        },
    },
    Environment: &spookytypes.ProjectEnvironment{
        Development: &spookytypes.EnvironmentConfig{
            Enabled: true,
            Machines: []string{"dev-server-1", "dev-server-2"},
            Variables: []string{"dev-vars.hcl"},
            SecurityLevel: "medium",
        },
        Production: &spookytypes.EnvironmentConfig{
            Enabled: true,
            Machines: []string{"prod-server-1", "prod-server-2", "prod-server-3"},
            Variables: []string{"prod-vars.hcl"},
            SecurityLevel: "high",
            BackupEnabled: true,
        },
    },
    Security: &spookytypes.ProjectSecurity{
        EncryptSensitiveData: true,
        ValidateFilePermissions: true,
        AuditAllOperations: true,
        AccessControlEnabled: true,
        AllowedUsers: []string{"admin", "operator"},
        AllowedGroups: []string{"devops", "sre"},
        Encryption: &spookytypes.EncryptionConfig{
            Method: "age",
            KeyPath: "~/.config/spooky/keys",
            BackupKeys: true,
            RotationInterval: "90d",
        },
    },
}

// Initialize the project
project, err := projectIntegration.InitializeProject(ctx, config)
if err != nil {
    return fmt.Errorf("failed to initialize project: %w", err)
}

log.Printf("Project initialized: %s", project.Name)
```

### Project Lifecycle Management

```go
// Start a project
err := projectManager.StartProject(ctx, "./my-project")
if err != nil {
    return fmt.Errorf("failed to start project: %w", err)
}

// Get project status
status, err := projectManager.GetProjectStatus(ctx, "./my-project")
if err != nil {
    return fmt.Errorf("failed to get project status: %w", err)
}

log.Printf("Project status: %s", status.Status)
log.Printf("Project health: %s", status.Health)

// Check component status
for component, componentStatus := range status.Components {
    log.Printf("Component %s: %s", component, componentStatus)
}

// Stop the project
err = projectManager.StopProject(ctx, "./my-project")
if err != nil {
    return fmt.Errorf("failed to stop project: %w", err)
}

// Clean up project resources
err = projectManager.CleanProject(ctx, "./my-project")
if err != nil {
    return fmt.Errorf("failed to clean project: %w", err)
}
```

### Error Handling

```go
// Handle project errors
project, err := projectIntegration.LoadProject(ctx, "./my-project")
if err != nil {
    var projectErr *spookytypes.ProjectError
    if errors.As(err, &projectErr) {
        log.Printf("Project error: %s (code: %s)", projectErr.Message, projectErr.Code)
        log.Printf("Operation: %s", projectErr.Operation)
        log.Printf("Project path: %s", projectErr.ProjectPath)
        
        for key, value := range projectErr.Details {
            log.Printf("Detail %s: %v", key, value)
        }
    }
    return fmt.Errorf("failed to load project: %w", err)
}

// Handle validation errors
result, err := projectIntegration.ValidateProject(ctx, project)
if err != nil {
    return fmt.Errorf("failed to validate project: %w", err)
}

if !result.IsValid() {
    for _, validationErr := range result.Errors {
        log.Printf("Validation error: %s", validationErr.Message)
        log.Printf("Field: %s", validationErr.Field)
        log.Printf("Value: %v", validationErr.Value)
        log.Printf("Constraint: %s", validationErr.Constraint)
        log.Printf("Severity: %s", validationErr.Severity)
    }
    return fmt.Errorf("project validation failed")
}
```

## Related Documentation

- [Projects System](PROJECTS_SYSTEM.md) - Complete system overview
- [Projects User Guide](PROJECTS_USER_GUIDE.md) - User guide and examples
- [Projects Troubleshooting](PROJECTS_TROUBLESHOOTING.md) - Troubleshooting guide
- [Configuration Management](CONFIGURATION_SYSTEM.md) - Configuration integration
- [Schema System](SCHEMA_SYSTEM.md) - Schema integration
- [CLI Reference](CLI_REFERENCE.md) - CLI integration
