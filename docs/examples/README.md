# Spooky Examples Directory

This directory contains comprehensive examples for all spooky systems and components, organized by system to maintain clarity and avoid disorganization as the project grows.

## Directory Structure

```
docs/examples/
├── README.md                    # This file - general overview
├── machines/                    # Machine inventory examples
│   ├── README.md               # Machines-specific documentation
│   ├── machines-basic-inventory.hcl
│   ├── machines-multi-environment.hcl
│   └── machines-kubernetes-nodes.hcl
├── facts/                       # Facts collection examples (future)
│   └── README.md
├── variables/                   # Variables management examples
│   ├── README.md               # Variables-specific documentation
│   ├── variables-basic-config.hcl
│   ├── variables-multi-file.hcl
│   └── variables-with-dependencies.hcl
├── logging/                     # Logging configuration examples
│   ├── README.md               # Logging-specific documentation
│   ├── global-logging-config.hcl
│   ├── project-logging-config.hcl
│   └── logging-formats.hcl
├── actions/                     # Actions and orchestration examples (future)
│   └── README.md
├── templates/                   # Template rendering examples (future)
│   └── README.md
├── ssh-basic-connection.hcl     # Basic SSH connection examples
├── ssh-key-types.hcl           # Different key type configurations
└── ssh-certificates.hcl        # SSH certificate authentication examples
└── projects/                    # Complete project examples (future)
    └── README.md
```

## Available Examples

### ✅ Machines System
**Location:** [`machines/`](machines/)
**Status:** Complete with comprehensive examples

### ✅ Variables System
**Location:** [`variables/`](variables/)
**Status:** Complete with comprehensive examples

**Examples Available:**
- **Basic Inventory** - Simple 3-machine setup for getting started
- **Multi-Environment** - Production/staging/development management
- **Kubernetes Nodes** - K8s cluster management with GPU support

**Use Cases:**
- Machine inventory management
- Multi-environment infrastructure
- Kubernetes cluster administration
- Resource planning and capacity management

### ✅ SSH System
**Location:** SSH examples in root directory
**Status:** Complete with comprehensive examples

**Examples Available:**
- **Basic Connection** - Simple SSH connection setup with different authentication methods
- **Key Types** - Examples for ED25519, ED25519-SK, and RSA key configurations
- **Certificates** - SSH certificate authentication with CA setup

**Use Cases:**
- SSH connection configuration and management
- Key type selection and validation
- Certificate-based authentication
- Security best practices implementation

### ✅ Logging System
**Location:** [`logging/`](logging/)
**Status:** Complete with comprehensive examples

**Examples Available:**
- **Global Configuration** - System-wide logging setup
- **Project Configuration** - Project-specific logging overrides
- **Format Examples** - Different output format configurations

**Use Cases:**
- Logging configuration and management
- Output format selection and customization
- Component-specific logging configuration
- Performance and audit logging setup

### 🚧 Future Systems (Planned)

#### Facts System
**Location:** `facts/` (to be created)
**Purpose:** Examples for fact collection, storage, and processing

**Planned Examples:**
- Basic fact collection
- Custom fact collectors
- Fact storage and retrieval
- Fact validation and processing

#### Variables System
**Location:** [`variables/`](variables/)
**Status:** Complete with comprehensive examples

**Examples Available:**
- **Basic Configuration** - Simple variable setup for getting started
- **Multi-File Organization** - Project-level variable management
- **Complex Dependencies** - Advanced variable relationships and validation

**Use Cases:**
- Variable management and organization
- Dependency resolution and validation
- Environment-specific configuration
- Sensitive data handling

#### Actions System
**Location:** `actions/` (to be created)
**Purpose:** Examples for action orchestration and running

**Planned Examples:**
- Basic action definitions
- Complex orchestration workflows
- Action templates and scripts
- Action validation and planning

#### Templates System
**Location:** `templates/` (to be created)
**Purpose:** Examples for template rendering and management

**Planned Examples:**
- Basic template rendering
- Template functions and helpers
- Template validation
- Complex template workflows

#### SSH System
**Location:** `ssh/` (to be created)
**Purpose:** Examples for SSH connection management

**Planned Examples:**
- SSH connection configuration
- Authentication methods
- Connection pooling
- SSH command running

#### Complete Projects
**Location:** `projects/` (to be created)
**Purpose:** End-to-end project examples

**Planned Examples:**
- Web application deployment
- Database migration projects
- Infrastructure provisioning
- Multi-service orchestration

## Using Examples

### Getting Started
1. **Choose Your System** - Navigate to the appropriate system directory
2. **Read the README** - Each system has its own detailed documentation
3. **Copy Examples** - Use examples as starting points for your configurations
4. **Customize** - Adapt examples to your specific environment and needs
5. **Test** - Validate your configurations using spooky commands

### Example Workflow
```bash
# 1. Navigate to the system you're working with
cd docs/examples/machines

# 2. Read the system-specific documentation
cat README.md

# 3. Copy an example to your project
cp machines-basic-inventory.hcl ../../my-project/machines.hcl

# 4. Customize the configuration
# Edit hostnames, IP addresses, users, SSH keys, etc.

# 5. Test your configuration
cd ../../my-project
spooky machines validate .
spooky machines list .
spooky machines ping .
```

## Organization Principles

### Why This Structure?

1. **Scalability** - Each system has its own space to grow
2. **Clarity** - Easy to find examples for specific systems
3. **Maintainability** - Changes to one system don't affect others
4. **Documentation** - Each system can have detailed, focused documentation
5. **Contributions** - Clear structure for adding new examples

### Naming Conventions

- **System directories** use lowercase: `machines/`, `facts/`, `variables/`
- **Example files** use descriptive names: `machines-basic-inventory.hcl`
- **README files** in each directory explain system-specific examples
- **Consistent structure** across all system directories

### Adding New Examples

When adding examples for new systems:

1. **Create system directory** - `mkdir docs/examples/new-system`
2. **Add system README** - Document examples and usage patterns
3. **Include example files** - Provide practical, working examples
4. **Update this README** - Add the new system to the overview
5. **Test examples** - Ensure all examples work with current spooky version

## Best Practices

### For Example Creation

1. **Start Simple** - Provide basic examples for getting started
2. **Show Progression** - Include intermediate and advanced examples
3. **Document Everything** - Clear comments and explanations
4. **Test Thoroughly** - Ensure examples work with current spooky version
5. **Follow Patterns** - Use consistent naming and structure

### For Example Usage

1. **Read Documentation** - Understand the system before using examples
2. **Start with Basics** - Begin with simple examples and progress
3. **Customize Appropriately** - Adapt examples to your environment
4. **Validate Configurations** - Always test your configurations
5. **Follow Best Practices** - Use examples as guides, not templates

## Contributing

When contributing examples:

1. **Follow Structure** - Place examples in appropriate system directories
2. **Add Documentation** - Include README files with usage instructions
3. **Test Examples** - Ensure examples work with current spooky version
4. **Use Clear Names** - Descriptive filenames and comments
5. **Update Overview** - Keep this README current with new systems

## Support

For help with examples:

1. **Check System README** - Each system has detailed documentation
2. **Review Main Documentation** - See the main docs directory for comprehensive guides
3. **Test Examples** - Use spooky commands to validate configurations
4. **Check Troubleshooting** - Review troubleshooting guides for common issues
5. **Ask Questions** - Use the project's support channels for specific help

## Current Status

- ✅ **Machines System** - Complete with comprehensive examples
- ✅ **Variables System** - Complete with comprehensive examples
- ✅ **Logging System** - Complete with comprehensive examples
- 🚧 **Other Systems** - Planned and ready for implementation
- 📋 **Documentation** - Structure established for future growth

As spooky evolves and new systems are implemented, this examples directory will grow systematically, maintaining organization and clarity for all users.
