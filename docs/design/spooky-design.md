# Spooky Design Document

## Overview

Spooky is a modern, Go-based infrastructure automation tool designed for simplicity, performance, and scalability. It provides a unified approach to server management, configuration, and deployment through a clean CLI interface and declarative configuration using HCL (HashiCorp Configuration Language).

### **Core Philosophy**

- **Simplicity First**: Clean, intuitive interfaces that reduce cognitive load
- **Performance**: Fast execution and efficient resource usage
- **Scalability**: Handle deployments from single servers to thousands
- **Portability**: Projects work across environments without modification
- **Security**: Secure handling of sensitive data with encryption support
- **Extensibility**: Plugin architecture for custom functionality

## Architecture Overview

### **High-Level Architecture**

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CLI Layer     │    │  Project Layer  │    │  Execution Layer │
│                 │    │                 │    │                 │
│ spooky noun verb│───▶│  Project Config │───▶│  SSH Execution  │
│                 │    │  Variables      │    │  Local Execution │
│                 │    │  Facts          │    │  Template Render │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Global Config  │    │  Storage Layer  │    │  Schema System  │
│                 │    │                 │    │                 │
│ ~/.config/spooky│    │  BadgerDB       │    │  HCL Schemas    │
│ ~/.local/state  │    │  JSON Files     │    │  Validation     │
│                 │    │  HCL Files      │    │  Type Safety    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### **Core Components**

1. **Schema System**: HCL schema validation and type safety (FOUNDATION)
2. **CLI System**: Noun-verb command structure with shell completion
3. **Project System**: Isolated, portable project environments
4. **Facts System**: Dynamic machine data collection and storage
5. **Variables System**: Typed configuration variables with validation
6. **Actions System**: Declarative task execution with dependencies and templating
7. **Configuration System**: Global configuration with XDG integration

## System Dependencies

### **Foundation Layer**
- **Schemas**: Provides schema validation, type safety, and validation infrastructure for all other systems

### **Core Layer**
- **CLI System**: Defines command patterns and user interface (depends on Schema System)
- **Configuration System**: Global configuration management (depends on Schema System)

### **Project Layer**
- **Project System**: Project lifecycle and structure (depends on Schema System, CLI System)
- **Facts System**: Data collection and storage (depends on Schema System, CLI System)
- **Variables System**: Configuration variables (depends on Schema System, CLI System)
- **Actions System**: Task execution, dependencies, and templating (depends on Schema System, CLI System, Variables System, Facts System, Project System)

### **Dependency Graph**

```mermaid
graph TD
    A[schemas --> B[cli]
    A --> C[configuration-system]
    A --> D[project-system]
    A --> E[facts-system]
    A --> F[variables-system]
    D --> G[actions-system]
    E --> G
    F --> G
    B --> G
    
    style A fill:#e1f5fe
    style B fill:#f3e5f5
    style C fill:#f3e5f5
    style D fill:#e8f5e8
    style E fill:#e8f5e8
    style F fill:#e8f5e8
    style G fill:#fff3e0
```

## Design Decisions

### **1. Language and Configuration**

**Decision**: Use Go with HCL for configuration
- **Go**: Performance, static typing, excellent tooling
- **HCL**: Human-readable, supports complex data structures, schema validation
- **Rationale**: Combines performance with developer experience

### **2. CLI Design Pattern**

**Decision**: Noun-verb command structure
```bash
spooky <noun> <verb> <project> [flags]
```

**Examples**:
- `spooky actions run myproject`
- `spooky facts gather myproject`
- `spooky machines ping myproject`
- `spooky variables validate myproject`

**Rationale**: 
- Intuitive and discoverable
- Follows established CLI patterns
- Enables comprehensive shell completion
- Scales well with new features

### **3. Project Isolation**

**Decision**: Isolated projects with global integration
```
project/
├── project.hcl          # Project metadata and config
├── machines.hcl         # Machine inventory
├── actions.hcl          # Action definitions
├── variables.hcl        # Variable definitions
├── templates/           # Template files
├── facts.db/           # Project-specific facts (optional)
└── logs/               # Execution logs
```

**Global Integration**:
- `~/.config/spooky/spooky.hcl` - Global configuration
- `~/.local/state/spooky/global-facts.db` - Global facts
- `/etc/spooky/` - System-wide facts and variables

**Rationale**: 
- Projects are portable and self-contained
- Global resources reduce duplication
- Clear separation of concerns

### **4. Facts System Design**

**Decision**: Hybrid storage with global and project facts

**Global Facts** (shared across projects):
- System information (hardware, OS, network)
- SSH host keys and connectivity data
- Package manager and service manager detection
- Machine identification data

**Project Facts** (project-specific):
- Application versions and deployment states
- Configuration values and custom data
- Environment-specific information

**Storage Options**:
- **BadgerDB**: High-performance key-value store for large deployments
- **JSON**: Human-readable, portable format
- **HCL**: Structured, schema-validated format

**Rationale**:
- Efficient storage for large deployments
- Multiple format support for different use cases
- Clear separation between global and project data

### **5. Variables System Design**

**Decision**: Typed variables with validation and interpolation

**Variable Types**:
- `string`, `number`, `bool`, `list`, `map`, `object`

**Features**:
- Schema validation with HCL schemas
- Variable interpolation in HCL files and templates
- Dependency management with circular reference detection
- File merging for `variables.hcl` and `variables/` directory
- Sensitive data handling with age encryption

**Precedence Order**:
1. Environment variables
2. Project variables
3. Machine facts
4. System defaults

**Rationale**:
- Type safety prevents configuration errors
- Validation catches issues early
- Interpolation enables dynamic configuration
- Security for sensitive data

### **6. Actions System Design**

**Decision**: Declarative actions with dependency management

**Action Properties**:
- `command`: Command to execute
- `description`: Human-readable description
- `depends_on`: Action dependencies
- `timeout`: Execution timeout
- `retries`: Retry attempts
- `parallel`: Parallel execution flag
- `tags`: Action categorization
- `machine`: Target machine specification

**Dependency Management**:
- Directed Acyclic Graph (DAG) for dependency tracking
- Topological sorting for execution order
- Circular reference detection
- Cross-file dependency resolution

**Rationale**:
- Declarative approach is more maintainable
- Dependencies ensure correct execution order
- Parallel execution improves performance
- Tags enable selective execution

### **7. Template System Design**

**Decision**: HCL-based templating with variable and fact integration

**Template Functions**:
- `{{var "name"}}` - Variable access
- `{{fact "name"}}` - Fact access
- `{{env "VAR"}}` - Environment variable access
- `{{machine "name"}}` - Machine data access

**Integration**:
- Variables and facts available in templates
- Template context includes all project data
- Support for complex data structures
- Error handling for missing values

**Rationale**:
- HCL provides familiar syntax
- Integration with variables and facts enables dynamic templates
- Type safety through schema validation
- Consistent with configuration language

### **8. Schema System Design**

**Decision**: HCL schemas for validation and type safety

**Schema Features**:
- Type validation for all configuration
- Required field validation
- Custom validation rules
- Schema documentation generation
- IDE support through schema files

**Schema Location**:
- `internal/schemas/schemas/` - Embedded schemas
- Schemas compiled into binary for portability
- No external schema dependencies

**Rationale**:
- Catches configuration errors early
- Provides IDE support and autocomplete
- Ensures consistency across projects
- Self-contained binary distribution

### **9. Storage System Design**

**Decision**: Hybrid storage with multiple backends

**Storage Backends**:
- **BadgerDB**: High-performance, embedded key-value store
- **JSON**: Human-readable, portable format
- **HCL**: Structured, validated format

**Storage Strategy**:
- Global facts in BadgerDB for performance
- Project facts in BadgerDB or JSON based on size
- Configuration files in HCL for validation
- Export/import support for all formats

**Rationale**:
- BadgerDB provides excellent performance for large datasets
- JSON enables easy data portability
- HCL provides validation and structure
- Multiple formats support different use cases

### **10. Security Design**

**Decision**: Secure by default with encryption support

**Security Features**:
- SSH key-based authentication
- Sensitive variable encryption with age
- Redaction of sensitive data in logs
- Secure handling of credentials
- No plaintext secrets in configuration

**Encryption**:
- **age**: Modern, simple encryption for sensitive variables
- File-level and variable-level encryption
- Key management through age keyrings
- Integration with existing SSH key infrastructure

**Rationale**:
- Security is critical for infrastructure automation
- age provides simple, secure encryption
- SSH integration leverages existing security practices
- Redaction prevents accidental exposure

## Implementation Phases

### **Phase 1: Schema Foundation**
- [x] Schema system implementation and validation
- [x] HCL schema definitions for all components
- [x] Schema validation infrastructure
- [x] Schema-based configuration validation

### **Phase 2: Core Infrastructure**
- [x] Basic CLI structure and commands
- [x] Project initialization and structure
- [x] HCL parsing and validation
- [x] SSH client implementation
- [x] Basic facts collection

### **Phase 3: Storage and Data**
- [x] BadgerDB storage implementation
- [x] JSON and HCL storage backends
- [x] Facts system with global and project storage
- [x] Schema system for validation

### **Phase 4: Variables and Templates**
- [ ] Variables system implementation
- [ ] Template system with variable integration
- [ ] File merging for variables
- [ ] Dependency management

### **Phase 5: Actions and Execution**
- [ ] Actions system implementation
- [ ] Dependency resolution and execution
- [ ] Parallel execution support
- [ ] Action validation and testing

### **Phase 6: Advanced Features**
- [ ] Advanced security features
- [ ] Performance optimizations
- [ ] Advanced templating features
- [ ] Plugin system

## Success Metrics

### **Performance**
- Project initialization: < 1 second
- Facts gathering: < 5 seconds per machine
- Action execution: < 10 seconds for simple actions
- Template rendering: < 100ms per template

### **Scalability**
- Support for 10,000+ machines
- Handle 1GB+ fact databases
- Parallel execution of 100+ actions
- Memory usage < 100MB for typical projects

### **Usability**
- CLI discoverability through help and completion
- Clear error messages with actionable suggestions
- Comprehensive documentation and examples
- Intuitive project structure and configuration

### **Reliability**
- Graceful handling of network failures
- Comprehensive error handling and recovery
- Data integrity through validation

## Error Handling and Resilience

### **Network and Machine Connectivity**

**Decision**: Fail-fast with retry logic and graceful degradation

#### **Machine Connectivity Failures**

**Detection**:
- SSH connection timeouts (configurable per project)
- Network unreachability detection
- Machine-specific authentication failures

**Response Strategy**:
- Continue execution for other machines when one fails
- Maximum retry attempts per machine with exponential backoff
- Parallel execution limits to prevent resource exhaustion
- Configurable failure handling behavior

See [error-handling.go](../snippets/error-handling.go) for implementation examples.

**Error Handling Patterns**:
1. **Individual Machine Failures**: Log error, continue with other machines
2. **Authentication Failures**: Skip machine, report in summary
3. **Network Timeouts**: Retry with exponential backoff
4. **Critical Failures**: Stop execution if `--fail-fast` flag is set

#### **URL and External Resource Failures**

**Detection**:
- HTTP/HTTPS connection timeouts
- DNS resolution failures
- SSL/TLS certificate issues
- Rate limiting responses

**Response Strategy**:
- Retry with exponential backoff
- Configurable timeout and retry limits
- Graceful handling of rate limiting
- Clear error reporting for debugging

See [error-handling.go](../snippets/error-handling.go) for URL fetching implementation.

#### **CLI Error Reporting**

**Error Categories**:
- **Warnings**: Non-critical issues (e.g., machine unreachable)
- **Errors**: Action failures that prevent completion
- **Fatal**: Critical system errors that stop execution

**Error Output Format**:
- Visual indicators: ✓ (success), ✗ (error), ⚠ (warning)
- Summary statistics showing total machines, successful, failed, and warnings
- Clear error messages with actionable information
- Progress indication for long-running operations

### **Data Integrity**

**Decision**: Schema validation and file integrity checks

#### **Schema Validation**
- HCL schema validation for all configuration files
- Type checking for variables and facts
- Required field validation
- Custom validation rules

#### **File Integrity**
- File existence checks before processing
- Read permissions validation
- File size limits to prevent resource exhaustion
- UTF-8 encoding validation for text files

**Note**: Checksums and cryptographic integrity verification are not currently implemented. Data integrity is ensured through schema validation and file system checks.

## Risk Assessment

### **Technical Risks**
- **Performance**: Large deployments may exceed performance targets
  - *Mitigation*: Profiling and optimization, caching strategies
- **Complexity**: Feature creep may impact simplicity
  - *Mitigation*: Strict feature requirements, user feedback
- **Compatibility**: HCL schema changes may break existing projects
  - *Mitigation*: Versioned schemas, migration tools

### **Operational Risks**
- **Security**: Sensitive data exposure through misconfiguration
  - *Mitigation*: Secure defaults, encryption, audit logging
- **Reliability**: Network failures during execution
  - *Mitigation*: Retry logic, circuit breakers, graceful degradation
- **Maintenance**: Complex dependency management
  - *Mitigation*: Clear documentation, automated testing