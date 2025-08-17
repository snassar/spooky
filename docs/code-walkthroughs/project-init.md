# Code Walkthrough: `spooky project init`

## Command Overview

**Command**: `spooky project init <project-name> [flags]`

**Purpose**: Creates a new spooky project with schema-driven directory structure and default configuration files.

**Key Features**:
- Schema-driven project creation using `project-directory.hcl` schema
- Customizable project metadata via flags
- Automatic generation of default configuration files
- Template-based file generation with Go templates

## Command Structure with All Flags

```bash
spooky project init testing-code-walkthrough \
  --name "Testing Code Walkthrough Project" \
  --description "A comprehensive test project for code walkthrough documentation" \
  --version "1.0.0" \
  --author "Spooky Developer" \
  --email "developer@spooky.example.com" \
  --url "https://github.com/spooky/testing-code-walkthrough"
```

## Execution Flow

### 1. Command Entry Point: `main()`

**File**: `main.go`
```go
func main() {
    cmd.Execute()
}
```

**What Happens**:
- Calls the root command's `Execute()` method
- Initializes the CLI framework and command structure

### 2. Root Command Setup: `RootCmd.Execute()`

**File**: `cmd/root.go`
```go
func Execute() {
    err := RootCmd.Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

**What Happens**:
- Executes the root command with all subcommands
- Handles any errors and exits with appropriate status code

### 3. Persistent Pre-Run: `RootCmd.PersistentPreRunE`

**File**: `cmd/root.go`
```go
PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
    // Skip auto-setup for version and help commands
    if cmd.Name() == "version" || cmd.Name() == "help" {
        return nil
    }
    // Auto-setup configuration for all other commands
    if err := spookyconfig.AutoSetupConfig(); err != nil {
        return fmt.Errorf("configuration setup failed: %w", err)
    }
    // Initialize integration manager for commands that need it
    if err := InitializeIntegrationsDependencies(); err != nil {
        return fmt.Errorf("integration initialization failed: %w", err)
    }
    return nil
}
```

**What Happens**:
- Sets up global configuration and logging
- Initializes integration dependencies
- Ensures proper environment setup for all commands

### 4. Project Init Command: `projectInitCmd`

**File**: `cmd/project.go`
```go
var projectInitCmd = &cobra.Command{
    Use:   "init [project-name]",
    Short: "Initialize a new spooky project",
    Long: `Initialize a new spooky project with the specified name.
Flags allow customization of project metadata during initialization.`,
    Args: cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        name, _ := cmd.Flags().GetString("name")
        description, _ := cmd.Flags().GetString("description")
        version, _ := cmd.Flags().GetString("version")
        author, _ := cmd.Flags().GetString("author")
        email, _ := cmd.Flags().GetString("email")
        url, _ := cmd.Flags().GetString("url")
        return handleProjectInit(args[0], name, description, version, author, email, url)
    },
}
```

**What Happens**:
- Defines the command structure and help text
- Validates exactly one argument (project name)
- Extracts all flag values for project metadata
- Calls the main handler function with all parameters

### 5. Project Init Handler: `handleProjectInit()`

**File**: `cmd/project.go`
```go
func handleProjectInit(projectPath, name, description, version, author, email, url string) error {
    ctx := context.Background()
    if projectManager == nil {
        if err := InitializeProjectDependencies(); err != nil {
            return fmt.Errorf("failed to initialize project dependencies: %w", err)
        }
    }
    
    project, err := projectManager.Initialize(ctx, projectPath)
    if err != nil {
        return fmt.Errorf("failed to initialize project: %w", err)
    }
    
    // Update project configuration with provided values if they were specified
    if name != "" || description != "" || version != "" || author != "" || url != "" {
        if name != "" {
            project.Config.Name = name
        }
        if description != "" {
            project.Config.Description = description
        }
        if version != "" {
            project.Config.Metadata.Version = version
        }
        if author != "" {
            project.Config.Metadata.Author = author
        }
        if email != "" {
            project.Config.Metadata.Email = email
        }
        if url != "" {
            project.Config.Metadata.URL = url
        }
        
        if err := projectManager.Save(ctx, project); err != nil {
            return fmt.Errorf("failed to save updated project configuration: %w", err)
        }
    }
    
    fmt.Printf("✅ Project initialized successfully: %s\n", projectPath)
    fmt.Printf("📁 Project directory: %s\n", project.Path)
    fmt.Printf("📋 Project name: %s\n", project.Config.Name)
    fmt.Printf("📝 Description: %s\n", project.Config.Description)
    fmt.Printf("🔧 Next steps:\n")
    fmt.Printf("   - Edit project.hcl to customize settings\n")
    fmt.Printf("   - Add machines to machines.hcl\n")
    fmt.Printf("   - Define actions in actions.hcl\n")
    fmt.Printf("   - Configure variables in variables.hcl\n")
    
    return nil
}
```

**What Happens**:
- Initializes project dependencies if needed
- Calls the project manager to create the project structure
- Updates project metadata with flag values if provided
- Saves the updated configuration
- Displays success message and next steps

### 6. Project Manager Initialization: `projectManager.Initialize()`

**File**: `internal/project/manager.go`
```go
func (m *Manager) Initialize(_ context.Context, projectPath string) (*spookytypes.Project, error) {
    m.logger.Info("Initializing new project", map[string]interface{}{
        "project_path": projectPath,
    })

    // Resolve absolute path
    absPath, err := filepath.Abs(projectPath)
    if err != nil {
        return nil, fmt.Errorf("failed to resolve project path: %w", err)
    }

    // Check if project directory already exists
    if _, err := os.Stat(absPath); err == nil {
        return nil, fmt.Errorf("project directory already exists: %s", absPath)
    }

    // Create project directory
    if err := os.MkdirAll(absPath, 0o755); err != nil {
        return nil, fmt.Errorf("failed to create project directory: %w", err)
    }

    // Load project directory schema for schema-driven creation
    schema, err := spookyschemas.LoadProjectDirectorySchema()
    if err != nil {
        return nil, fmt.Errorf("failed to load project directory schema: %w", err)
    }

    // Create project structure based on schema
    if err := m.createProjectFromSchema(absPath, schema); err != nil {
        return nil, fmt.Errorf("failed to create project structure: %w", err)
    }

    // Create default project configuration
    projectName := filepath.Base(absPath)
    project := &spookytypes.Project{
        Path: absPath,
        Config: &spookytypesproject.ProjectConfig{
            Name:        projectName,
            Description: fmt.Sprintf("A spooky project: %s", projectName),
            Metadata: &spookytypesproject.ProjectMetadata{
                Version: "0.1.0",
                Author:  "spooky-user",
                Created: time.Now(),
                Updated: time.Now(),
            },
            Settings: &spookytypesproject.ProjectSettings{
                ParallelWorkers:    4,
                TimeoutSeconds:     300,
                LogLevel:           "info",
                DefaultDryRun:      false,
                ValidateBeforeRun:  true,
                MaxRetries:         3,
                RetryDelaySeconds:  5,
            },
        },
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    // Create project.hcl file
    if err := m.createProjectHCL(project); err != nil {
        return nil, fmt.Errorf("failed to create project.hcl: %w", err)
    }

    // Create README.md file
    if err := m.createREADME(project); err != nil {
        return nil, fmt.Errorf("failed to create README.md: %w", err)
    }

    m.logger.Info("Project initialized successfully", map[string]interface{}{
        "project_path": absPath,
        "project_name": projectName,
    })

    return project, nil
}
```

**What Happens**:
- Resolves the absolute project path
- Checks if the project directory already exists
- Creates the project directory
- Loads the project directory schema for schema-driven creation
- Creates project structure based on schema definitions
- Creates default project configuration
- Generates project.hcl and README.md files
- Returns the created project

### 7. Schema-Driven Project Creation: `createProjectFromSchema()`

**File**: `internal/project/manager.go`
```go
func (m *Manager) createProjectFromSchema(projectPath string, schema *spookyschemas.ProjectDirectorySchema) error {
    m.logger.Debug("Creating project structure from schema", map[string]interface{}{
        "project_path": projectPath,
    })

    // Create required directories from schema
    for _, dir := range schema.GetRequiredDirectories() {
        dirPath := filepath.Join(projectPath, dir.Name)
        if err := os.MkdirAll(dirPath, 0o755); err != nil {
            return fmt.Errorf("failed to create required directory %s: %w", dir.Name, err)
        }
        m.logger.Debug("Created required directory", map[string]interface{}{
            "directory": dir.Name,
            "path":      dirPath,
        })
    }

    // Create optional directories from schema
    for _, dir := range schema.GetOptionalDirectories() {
        dirPath := filepath.Join(projectPath, dir.Name)
        if err := os.MkdirAll(dirPath, 0o755); err != nil {
            return fmt.Errorf("failed to create optional directory %s: %w", dir.Name, err)
        }
        m.logger.Debug("Created optional directory", map[string]interface{}{
            "directory": dir.Name,
            "path":      dirPath,
        })
    }

    // Create default files for optional files that don't exist
    optionalFiles := []struct {
        name    string
        content string
    }{
        {"machines.hcl", `machines {
  # Add your machines here
  # Example:
  # machine "web-server" {
  #   hostname = "web.example.com"
  #   port = 22
  #   user = "admin"
  #   authentication {
  #     method = "ssh_key"
  #     key_path = "~/.ssh/id_rsa"
  #   }
  # }
}`},
        {"actions.hcl", `actions {
  # Add your actions here
  # Example:
  # action "update-system" {
  #   description = "Update system packages"
  #   machines = ["web-server"]
  #   command = "sudo apt update && sudo apt upgrade -y"
  # }
}`},
        {"variables.hcl", `variables {
  # Add your variables here
  # Example:
  # variable "app_version" {
  #   description = "Application version to deploy"
  #   default = "1.0.0"
  # }
}`},
        {"recipients.txt", `# Add age encryption recipients here
# One recipient per line
# Example:
# age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p`},
    }

    for _, file := range optionalFiles {
        filePath := filepath.Join(projectPath, file.name)
        if _, err := os.Stat(filePath); os.IsNotExist(err) {
            if err := os.WriteFile(filePath, []byte(file.content), 0o644); err != nil {
                return fmt.Errorf("failed to create %s: %w", file.name, err)
            }
            m.logger.Debug("Created optional file", map[string]interface{}{
                "file": file.name,
                "path": filePath,
            })
        }
    }

    return nil
}
```

**What Happens**:
- Creates required directories based on schema definition
- Creates optional directories based on schema definition
- Creates default content for optional files if they don't exist
- Logs the creation of each directory and file
- Returns error if any creation fails

### 8. Schema Loading: `LoadProjectDirectorySchema()`

**File**: `internal/schemas/project_directory_schema.go`
```go
func LoadProjectDirectorySchema() (*ProjectDirectorySchema, error) {
    // Load the embedded project directory schema
    schemaData, err := embeddedSchemas.ReadFile("schemas/structure/project-directory.hcl")
    if err != nil {
        return nil, fmt.Errorf("failed to read embedded project directory schema: %w", err)
    }

    // Parse the schema
    schema, err := ParseProjectDirectorySchema(string(schemaData))
    if err != nil {
        return nil, fmt.Errorf("failed to parse project directory schema: %w", err)
    }

    return schema, nil
}
```

**What Happens**:
- Loads the embedded project directory schema from the binary
- Parses the HCL schema content
- Returns the parsed schema structure
- Handles any loading or parsing errors

### 9. Project HCL Creation: `createProjectHCL()`

**File**: `internal/project/manager.go`
```go
func (m *Manager) createProjectHCL(project *spookytypes.Project) error {
    projectHCLPath := filepath.Join(project.Path, "project.hcl")

    // Parse the template
    tmpl, err := template.New("project").Parse(projectTemplate)
    if err != nil {
        return fmt.Errorf("failed to parse project template: %w", err)
    }

    // Create the file
    file, err := os.Create(projectHCLPath)
    if err != nil {
        return fmt.Errorf("failed to create project.hcl: %w", err)
    }
    defer file.Close()

    // Execute the template
    if err := tmpl.Execute(file, project.Config); err != nil {
        return fmt.Errorf("failed to execute project template: %w", err)
    }

    m.logger.Debug("Created project.hcl", map[string]interface{}{
        "path": projectHCLPath,
    })

    return nil
}
```

**What Happens**:
- Creates the project.hcl file path
- Parses the Go template for project configuration
- Creates the file on disk
- Executes the template with project configuration data
- Logs the successful creation

### 10. README Creation: `createREADME()`

**File**: `internal/project/manager.go`
```go
func (m *Manager) createREADME(project *spookytypes.Project) error {
    readmePath := filepath.Join(project.Path, "README.md")

    // Parse the template
    tmpl, err := template.New("readme").Parse(readmeTemplate)
    if err != nil {
        return fmt.Errorf("failed to parse README template: %w", err)
    }

    // Create the file
    file, err := os.Create(readmePath)
    if err != nil {
        return fmt.Errorf("failed to create README.md: %w", err)
    }
    defer file.Close()

    // Execute the template
    if err := tmpl.Execute(file, project.Config); err != nil {
        return fmt.Errorf("failed to execute README template: %w", err)
    }

    m.logger.Debug("Created README.md", map[string]interface{}{
        "path": readmePath,
    })

    return nil
}
```

**What Happens**:
- Creates the README.md file path
- Parses the Go template for README content
- Creates the file on disk
- Executes the template with project configuration data
- Logs the successful creation

## Generated Project Structure

Based on the schema-driven approach, the following structure is created:

```
testing-code-walkthrough/
├── project.hcl          # Project configuration (required)
├── README.md            # Project documentation (optional)
├── machines.hcl         # Machine inventory (optional)
├── actions.hcl          # Action definitions (optional)
├── variables.hcl        # Variable definitions (optional)
├── recipients.txt       # Age encryption recipients (optional)
├── machines/            # Machine-specific files (optional)
├── actions/             # Action-specific files (optional)
├── variables/           # Variable-specific files (optional)
├── templates/           # Template files (optional)
├── files/               # Project files (optional)
└── logs/                # Log files (optional)
```

## Generated Files Content

### project.hcl
```hcl
project {
  name = "testing-code-walkthrough"
  description = "A comprehensive test project for code walkthrough documentation"
  
  metadata {
    version = "1.0.0"
    author = "Spooky Developer"
    email = "developer@spooky.example.com"
    url = "https://github.com/spooky/testing-code-walkthrough"
    created = "2024-01-15T10:30:00Z"
    updated = "2024-01-15T10:30:00Z"
  }
  
  settings {
    parallel_workers = 4
    timeout_seconds = 300
    log_level = "info"
    default_dry_run = false
    validate_before_run = true
    max_retries = 3
    retry_delay_seconds = 5
  }
}
```

### README.md
```markdown
# Testing Code Walkthrough Project

A comprehensive test project for code walkthrough documentation.

## Project Information

- **Name**: testing-code-walkthrough
- **Description**: A comprehensive test project for code walkthrough documentation
- **Version**: 1.0.0
- **Author**: Spooky Developer
- **Email**: developer@spooky.example.com
- **URL**: https://github.com/spooky/testing-code-walkthrough

## Configuration

This project uses the following configuration files:

- `project.hcl` - Project configuration and metadata
- `machines.hcl` - Machine inventory and connectivity
- `actions.hcl` - Action definitions and scripts
- `variables.hcl` - Variable definitions and values

## Getting Started

1. Review and customize the configuration files
2. Add your machines to `machines.hcl`
3. Define your actions in `actions.hcl`
4. Configure variables in `variables.hcl`
5. Run your first action: `spooky actions run .`

## Project Structure

```
.
├── project.hcl          # Project configuration
├── README.md            # This file
├── machines.hcl         # Machine inventory
├── actions.hcl          # Action definitions
├── variables.hcl        # Variable definitions
├── recipients.txt       # Age encryption recipients
├── machines/            # Machine-specific files
├── actions/             # Action-specific files
├── variables/           # Variable-specific files
├── templates/           # Template files
├── files/               # Project files
└── logs/                # Log files
```

## Next Steps

- Customize the project configuration in `project.hcl`
- Add your target machines to `machines.hcl`
- Define automation actions in `actions.hcl`
- Configure project variables in `variables.hcl`
- Test your setup with `spooky project validate .`
```

## Key Components

### Project Manager
- **Purpose**: Coordinates project creation and management
- **Location**: `internal/project/manager.go`
- **Key Methods**: `Initialize()`, `createProjectFromSchema()`, `createProjectHCL()`, `createREADME()`

### Schema Parser
- **Purpose**: Parses project directory schema for schema-driven creation
- **Location**: `internal/schemas/project_directory_schema.go`
- **Key Methods**: `LoadProjectDirectorySchema()`, `ParseProjectDirectorySchema()`

### Project Directory Schema
- **Purpose**: Defines the expected project structure
- **Location**: `internal/schemas/schemas/structure/project-directory.hcl`
- **Content**: Required/optional files and directories

## Flag Processing

### Available Flags
- `--name`: Project name (defaults to directory name)
- `--description`: Project description
- `--version`: Project version (defaults to "0.1.0")
- `--author`: Project author (defaults to "spooky-user")
- `--email`: Author email address
- `--url`: Project URL

### Flag Processing Flow
1. **Flag Extraction**: All flags are extracted in `projectInitCmd.RunE`
2. **Default Values**: Default values are applied if flags are not provided
3. **Project Creation**: Project is created with default configuration
4. **Flag Override**: Flag values override defaults in `handleProjectInit`
5. **Configuration Save**: Updated configuration is saved to disk

## Error Handling

### Common Error Scenarios
1. **Directory Exists**: Project directory already exists
2. **Permission Denied**: Insufficient permissions to create directory
3. **Invalid Path**: Invalid project path provided
4. **Schema Loading**: Failed to load project directory schema
5. **File Creation**: Failed to create configuration files

### Error Recovery
- **Graceful Degradation**: Continues with defaults if optional components fail
- **Detailed Error Messages**: Provides specific error context
- **Cleanup**: Removes partially created directories on failure
- **Logging**: Comprehensive logging for troubleshooting

## Architecture Patterns

### Schema-Driven Design
- Project structure is defined by embedded schemas
- Schema parsing provides flexibility for structure changes
- Validation ensures schema compliance

### Template-Based Generation
- Go templates for file generation
- Consistent file formatting and structure
- Customizable content based on project metadata

### Dependency Injection
- Project manager is injected into CLI layer
- Schema parser is injected into project manager
- Logger is injected into all components

### Interface-Based Design
- All components use interface contracts
- Loose coupling between components
- Testable and extensible architecture

## Integration Points

### Schema System
- Uses embedded HCL schemas for structure definition
- Schema parsing for dynamic structure creation
- Validation integration for compliance checking

### Configuration System
- Integrates with global configuration setup
- Uses XDG base directory specification
- Supports environment variable overrides

### Logging System
- Structured logging with context
- Debug information for troubleshooting
- Error tracking and reporting

## Example Usage

```bash
# Basic project initialization
spooky project init my-project

# Project with custom metadata
spooky project init my-project \
  --name "My Awesome Project" \
  --description "A comprehensive automation project" \
  --version "2.0.0" \
  --author "John Doe" \
  --email "john@example.com" \
  --url "https://github.com/johndoe/my-project"

# Output
✅ Project initialized successfully: my-project
📁 Project directory: /path/to/my-project
📋 Project name: My Awesome Project
📝 Description: A comprehensive automation project
🔧 Next steps:
   - Edit project.hcl to customize settings
   - Add machines to machines.hcl
   - Define actions in actions.hcl
   - Configure variables in variables.hcl
```

## Exit Codes

- **0**: Project initialization successful
- **1**: Project initialization failed

## Performance Considerations

- Schema loading is cached after first load
- File system operations are optimized
- Template parsing is done once per file
- Parallel file creation where possible
