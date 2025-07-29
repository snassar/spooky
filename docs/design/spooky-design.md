# Spooky Design Document

## Overview

Spooky is a modern infrastructure automation tool designed for simplicity, performance, and reliability. It provides a declarative approach to managing infrastructure through HCL (HashiCorp Configuration Language) configurations, with built-in support for facts collection, variable management, and template rendering.

## Architecture

Spooky follows a layered architecture with clear separation of concerns:

### High-Level Architecture

```mermaid
graph TD
    subgraph "Schema Foundation"
        SCHEMAS[Schema System<br/>internal/schemas/schemas/<br/>Defines all configuration formats]
        SCHEMA_COMPOSER[Schema Composer<br/>Runtime composition<br/>Embedded validation]
        SCHEMA_VALIDATOR[Schema Validator<br/>Runtime validation<br/>Error reporting]
    end
    
    subgraph "Global Configuration"
        GLOBAL[Global Spooky Config<br/>~/.config/spooky/spooky.hcl<br/>Must conform to schema]
        GLOBAL_FACTS[Global Facts Database<br/>~/.local/state/spooky/global-facts.db<br/>Shared across projects]
    end
    
    subgraph "Project Structure"
        PROJECT[Project Configuration<br/>project.hcl<br/>Defines project metadata]
        
        subgraph "Project Data"
            MACHINES[machines.hcl<br/>Machine inventory]
            ACTIONS[actions.hcl + actions/<br/>Action definitions]
            VARIABLES[variables.hcl + variables/<br/>Project variables]
            FACTS[facts.db/facts.hcl/facts.json<br/>Project-specific facts]
        end
        
        subgraph "Project Assets"
            TEMPLATES[templates/<br/>Template files]
            FILES[files/<br/>Non-templated files]
        end
    end
    
    subgraph "CLI System"
        CLI_COMMANDS[CLI Commands<br/>spooky noun verb<br/>Schema validation]
        CLI_VALIDATION[Validation Commands<br/>validate, list, show, export<br/>Schema-based validation]
    end
    
    %% Schema validation relationships
    SCHEMAS --> SCHEMA_COMPOSER
    SCHEMA_COMPOSER --> SCHEMA_VALIDATOR
    SCHEMA_VALIDATOR --> GLOBAL
    SCHEMA_VALIDATOR --> PROJECT
    SCHEMA_VALIDATOR --> MACHINES
    SCHEMA_VALIDATOR --> ACTIONS
    SCHEMA_VALIDATOR --> VARIABLES
    SCHEMA_VALIDATOR --> FACTS
    
    %% CLI validation relationships
    CLI_COMMANDS --> SCHEMA_VALIDATOR
    CLI_VALIDATION --> SCHEMA_VALIDATOR
    
    %% Configuration relationships
    GLOBAL --> PROJECT
    GLOBAL_FACTS --> FACTS
    
    %% Project data relationships
    PROJECT --> MACHINES
    PROJECT --> ACTIONS
    PROJECT --> VARIABLES
    PROJECT --> FACTS
    PROJECT --> TEMPLATES
    PROJECT --> FILES
    
    %% Data flow for execution
    MACHINES -.->|provides targets| ACTIONS
    VARIABLES -.->|provides config| ACTIONS
    FACTS -.->|provides data| ACTIONS
    TEMPLATES -.->|provides templates| ACTIONS
    FILES -.->|provides files| ACTIONS
    
    %% Schema validation flow
    SCHEMA_VALIDATOR -.->|validates| CLI_COMMANDS
    SCHEMA_VALIDATOR -.->|validates| CLI_VALIDATION
    
    style SCHEMAS fill:#e1f5fe
    style SCHEMA_COMPOSER fill:#e1f5fe
    style SCHEMA_VALIDATOR fill:#e1f5fe
    style GLOBAL fill:#f3e5f5
    style GLOBAL_FACTS fill:#f3e5f5
    style PROJECT fill:#e8f5e8
    style MACHINES fill:#fff3e0
    style ACTIONS fill:#fff3e0
    style VARIABLES fill:#fff3e0
    style FACTS fill:#fff3e0
    style TEMPLATES fill:#fce4ec
    style FILES fill:#fce4ec
    style CLI_COMMANDS fill:#f0f4c3
    style CLI_VALIDATION fill:#f0f4c3
```

### System Dependencies

```mermaid
graph TD
    A[schemas] --> B[cli]
    A --> C[configuration]
    A --> D[project]
    A --> E[facts]
    A --> F[variables]
    A --> G[actions]
    C --> B
    D --> G
    E --> G
    F --> G
    
    style A fill:#e1f5fe
    style B fill:#f3e5f5
    style C fill:#f3e5f5
    style D fill:#e8f5e8
    style E fill:#e8f5e8
    style F fill:#e8f5e8
    style G fill:#fff3e0
```

## Core Components

### Schema System
- **Location**: `internal/schemas/schemas/`
- **Purpose**: Defines all configuration formats and validation rules
- **Features**: Embedded schema composition, runtime validation, facts system integration
- **Integration**: Used by all validation commands (`spooky validate`, `spooky facts validate`, etc.)
- **Reference**: See [Schema System](../systems/schema-system.md) for detailed implementation

### Configuration System
- **Global Config**: `~/.config/spooky/spooky.hcl`
- **Project Config**: `project.hcl`
- **Validation**: All configs must conform to embedded schemas

### Facts System
- **Storage**: BadgerDB (default), HCL, JSON
- **Scope**: Global (`~/.local/state/spooky/global-facts.db`) and project-specific
- **Collection**: Via SSH or local access using gopsutil
- **Schema**: Comprehensive system information coverage

### Variables System
- **Storage**: HCL files in project directory (`variables.hcl` and `variables/*.hcl`)
- **Purpose**: Project-specific configuration values and runtime parameters
- **Integration**: Used by actions system for template rendering and configuration
- **Validation**: Schema-based validation with embedded variable schemas
- **Reference**: See [Variables System](../systems/variables-system.md) for detailed implementation

### Actions System
- **Definition**: HCL files with command, script, or template actions
- **Execution**: Parallel execution with configurable concurrency
- **Dependencies**: Support for action dependencies and ordering
- **Templates**: Integrated template rendering with facts and variables

## Schema Integration

The schema system provides embedded validation for all configuration components. See [Schema System](../systems/schema-system.md) for detailed implementation, composition approach, and build process.

## Detailed Implementation References

For comprehensive implementation details of each system component, see:

### **Core Systems**
- **[Project System](../systems/project-system.md)** - Complete project system implementation plan with schema validation, configuration precedence, and CLI integration
- **[Schema System](../systems/schema-system.md)** - Schema validation and composition system with embedded runtime validation
- **[Variables System](../systems/variables-system.md)** - Variable management, file merging, dependency resolution, and template integration

### **Integration Points**
- **Project Configuration**: Project metadata, isolation, and portability patterns
- **Schema Validation**: Embedded schema composition and runtime validation
- **Variable Management**: File merging, conflict detection, and precedence resolution
- **CLI Integration**: Command patterns, error handling, and global flag support
- **Template System**: Variable access, facts integration, and context management

### **Implementation Patterns**
- **Schema-First Development**: All configuration formats defined by schemas before implementation
- **Embedded Validation**: Runtime schema composition with no external dependencies
- **File Merging**: Support for single files and directories with strict conflict detection
- **Configuration Precedence**: Environment → Project → Global → Defaults hierarchy
- **Error Handling**: Comprehensive error reporting with file and line numbers

## Design Decisions

### 1. Schema-First Development
- All configuration formats defined by schemas before implementation
- Ensures consistency and validation across the system
- Enables tooling and IDE support

### 2. HCL as Primary Language
- Human-readable configuration format
- Rich type system and validation
- Excellent tooling ecosystem

### 3. BadgerDB for Facts Storage
- High-performance key-value store
- ACID transactions and durability
- Native Go implementation

### 4. SSH-First Connectivity
- Standard protocol for remote access
- No agent installation required
- Secure by default

### 5. Template Integration
- Templates embedded in actions system
- Support for facts and variables
- Multiple template engines

### 6. Parallel Execution
- Configurable concurrency levels
- Dependency-aware execution
- Progress tracking and error handling

### 7. File Merging Strategy
- Support for single files and directories
- Strict validation prevents conflicts
- Clear precedence rules

### 8. Error Handling and Resilience
- Comprehensive error reporting
- Graceful handling of network failures
- Retry mechanisms for transient errors

### 9. Embedded Schema Composition
- Runtime schema composition for flexibility
- Ensures schema consistency across formats
- Maintains embeddability and performance

### 10. No Backward Compatibility
- Freedom to evolve the system
- No legacy code maintenance burden
- Clean, modern architecture

## File Merging Rules

## File Merging and Conflict Resolution

Spooky supports merging configuration from multiple files with strict validation:

- **Machine Inventory**: `machines.hcl` and `machines/*.hcl` files
- **Variables**: `variables.hcl` and `variables/*.hcl` files  
- **Actions**: `actions.hcl` and `actions/*.hcl` files

All file merging follows strict conflict detection and resolution rules. See [Variables System](../systems/variables-system.md) for detailed variable file merging rules and conflict resolution strategies.