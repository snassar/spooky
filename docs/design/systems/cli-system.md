# CLI System: Comprehensive Implementation Plan

## Overview

This document is the authoritative source for all CLI system implementation details in spooky. It covers command patterns, user interface design, shell completion, help systems, and integration with all other spooky systems.

**Schema Integration**: This CLI system follows the schema validation patterns and infrastructure defined in [Schema System](../schema-system.md) for all configuration validation and import operations.

**Architecture Integration**: CLI integrates with the overall spooky architecture as described in [Spooky Design](../spooky-design.md), providing the primary user interface for all system operations and management.

## System Integration

This CLI system integrates with other core Spooky systems to provide comprehensive command-line management:

### **Actions System Integration**
- **Action Commands**: `spooky actions` commands for action management and execution (see [Actions System](../actions-system.md))
- **Action Validation**: `spooky actions validate` for action file validation
- **Action Execution**: `spooky actions run` and `spooky actions dry-run` for action execution
- **Action Discovery**: `spooky actions list` for action discovery and filtering
- **Dependency Management**: Action dependency resolution and circular reference detection

### **Facts System Integration**
- **Facts Commands**: `spooky facts` commands for facts collection and management (see [Facts System](../facts-system.md))
- **Facts Collection**: `spooky facts gather` for machine facts collection
- **Facts Validation**: `spooky facts validate` for facts database validation
- **Facts Export**: `spooky facts export` for facts data export
- **Facts Discovery**: `spooky facts list` for facts discovery and filtering

### **Project System Integration**
- **Project Commands**: `spooky project` commands for project management (see [Project System](../project-system.md))
- **Project Initialization**: `spooky project init` for project creation
- **Project Validation**: `spooky project validate` for project structure validation
- **Project Information**: `spooky project show` for project information display
- **Project Discovery**: Project discovery using standard Unix tools

### **Variables System Integration**
- **Variables Commands**: `spooky variables` commands for variable management (see [Variables System](../variables-system.md))
- **Variables Validation**: `spooky variables validate` for variable file validation
- **Variables Discovery**: `spooky variables list` for variable discovery and filtering
- **Variables Display**: `spooky variables show` for detailed variable information
- **Variables Export**: `spooky variables export` for variable data export

### **Configuration System Integration**
- **Configuration Commands**: `spooky config` commands for configuration management (see [Configuration System](../configuration-system.md))
- **Configuration Display**: `spooky config show` for configuration information
- **Configuration Validation**: `spooky config validate` for configuration file validation
- **Configuration Discovery**: `spooky config list` for configuration section discovery
- **Global Configuration**: Integration with global configuration system

### **Templates System Integration**
- **Template Commands**: `spooky templates` commands for template management and rendering (see [Template System](../template-system.md))
- **Template Validation**: `spooky templates validate` for template syntax and structure validation
- **Template Rendering**: `spooky templates render` for template rendering with data context
- **Template Discovery**: `spooky templates list` for template discovery and filtering
- **Template Schema Integration**: Template validation uses schemas from [`template-metadata.hcl`](../../internal/schemas/schemas/template-metadata.hcl), [`template-context.hcl`](../../internal/schemas/schemas/template-context.hcl), and [`template-functions.hcl`](../../internal/schemas/schemas/template-functions.hcl)

### **Machines System Integration**
- **Machine Commands**: Machine management through `spooky machines` commands (see [Machines System](../machines-system.md))
- **Command Patterns**: Machine commands follow the established `spooky noun verb` CLI pattern
- **Validation Commands**: `spooky machines validate` for inventory validation
- **Management Commands**: `spooky machines list`, `spooky machines ping`, `spooky machines connect`
- **Import/Export Commands**: `spooky machines import`, `spooky machines export` for external system integration

### **Schema System Integration**
- **Schema Validation**: CLI commands use schema validation for all configuration and data files (see [Schema System](../schema-system.md))
- **Schema Commands**: Schema validation integrated into all CLI validation commands
- **Schema Discovery**: CLI provides access to embedded schemas for validation
- **Schema Evolution**: CLI supports schema versioning and migration

## Proposed CLI Design

```
# Core execution and management
spooky actions {list, validate, run, dry-run} <project directory>
spooky facts {gather, list, validate, export} <project directory>
spooky machines {list, validate, ping, connect, import, export} <project directory>
spooky templates {list, validate, render} <project directory>
spooky variables {validate, list, show, export} <project directory>
spooky project {show, validate, init} <project directory>

# Configuration and system management
spooky config {list, validate, show} [key] [value]
spooky logs {list, show, tail, clear, export} <project directory>

# System status and utilities
spooky status {show, check, report} <project directory>
spooky help {show, search, examples} [topic]
spooky completion {generate} [--shell SHELL] [--output FILE]

# Common aliases for frequent operations
spooky run <project directory>     # alias for 'spooky actions run'
spooky ls <project directory>      # alias for 'spooky project show'
spooky init <project directory>    # alias for 'spooky project init'

# Global flags
spooky --version                   # show version information
spooky --help                      # show help
```

## Implementation Details

### **Shell Completion Implementation**

#### **Command Structure**
```bash
spooky completion generate                    # auto-detect current shell, output to stdout
spooky completion generate --shell bash       # generate for specific shell
spooky completion generate --shell zsh
spooky completion generate --shell fish
spooky completion generate --shell powershell
spooky completion generate --shell bash --output ~/spooky-bash-completion.sh  # save to file
```

#### **Auto-Detection Logic**
1. **Read `$SHELL` environment variable** to determine current shell
2. **Map shell paths to names**:
   - `/bin/bash` → `bash`
   - `/bin/zsh` → `zsh`
   - `/usr/bin/fish` → `fish`
   - `powershell.exe` → `powershell`
3. **Fallback to `bash`** if detection fails

#### **Supported Shells**
- **bash**: Most common on Linux/macOS
- **zsh**: Default on macOS, popular on Linux
- **fish**: Modern shell with advanced completion
- **powershell**: Windows PowerShell

#### **Output Format**
- **Output to stdout** (user can redirect or copy)
- **No automatic file writing** (requires user action)
- **Include installation instructions** in help text

#### **Example Output for `spooky completion generate --shell bash`**
```bash
$ spooky completion generate --shell bash

# spooky bash completion script
# Generated by spooky completion generate --shell bash
# Install by adding to ~/.bashrc: source <(spooky completion generate --shell bash)

_spooky() {
    local cur prev opts cmds
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    # Main commands (nouns)
    cmds="actions facts machines templates variables project config logs status version help completion"

    # If completing the first word (command)
    if [[ $COMP_CWORD -eq 1 ]]; then
        COMPREPLY=( $(compgen -W "${cmds}" -- "${cur}") )
        return 0
    fi

    # Complete verbs based on noun
    case "${prev}" in
        actions)
            COMPREPLY=( $(compgen -W "list validate run dry-run" -- "${cur}") )
            ;;
        facts)
            COMPREPLY=( $(compgen -W "gather list validate export" -- "${cur}") )
            ;;
        machines)
            COMPREPLY=( $(compgen -W "list validate ping connect import export" -- "${cur}") )
            ;;
        templates)
            COMPREPLY=( $(compgen -W "list validate render" -- "${cur}") )
            ;;
        variables)
            COMPREPLY=( $(compgen -W "validate list show export" -- "${cur}") )
            ;;
        project)
            COMPREPLY=( $(compgen -W "list validate init info" -- "${cur}") )
            ;;
        config)
            COMPREPLY=( $(compgen -W "list validate show edit set" -- "${cur}") )
            ;;
        logs)
            COMPREPLY=( $(compgen -W "list show tail clear export" -- "${cur}") )
            ;;
        status)
            COMPREPLY=( $(compgen -W "show check report" -- "${cur}") )
            ;;
        version)
            COMPREPLY=( $(compgen -W "show check update" -- "${cur}") )
            ;;
        help)
            COMPREPLY=( $(compgen -W "show search examples" -- "${cur}") )
            ;;
        completion)
            COMPREPLY=( $(compgen -W "generate" -- "${cur}") )
            ;;
        generate)
            if [[ "${COMP_WORDS[COMP_CWORD-2]}" == "completion" ]]; then
                COMPREPLY=( $(compgen -W "--shell --output" -- "${cur}") )
            fi
            ;;
        --shell)
            COMPREPLY=( $(compgen -W "bash zsh fish powershell" -- "${cur}") )
            ;;
        --output)
            # Complete file paths for output
            COMPREPLY=( $(compgen -f -- "${cur}") )
            ;;
        *)
            # Try to complete project directories
            if [[ -d "${cur}" ]]; then
                COMPREPLY=( $(compgen -d -- "${cur}") )
            fi
            ;;
    esac
}

complete -F _spooky spooky
```

#### **Installation Instructions**
The completion script includes installation instructions:
```bash
# For bash (direct installation)
spooky completion generate --shell bash >> ~/.bashrc

# For bash (review first)
spooky completion generate --shell bash --output ~/spooky-bash-completion.sh
# Review the file, then: source ~/spooky-bash-completion.sh

# For zsh (direct installation)
spooky completion generate --shell zsh >> ~/.zshrc

# For zsh (review first)
spooky completion generate --shell zsh --output ~/spooky-zsh-completion.sh
# Review the file, then: source ~/spooky-zsh-completion.sh

# For fish (review first)
spooky completion generate --shell fish --output $XDG_CONFIG_HOME/fish/completions/spooky.fish

# For PowerShell (review first)
spooky completion generate --shell powershell --output spooky-completion.ps1
# Review the file, then add to PowerShell profile: . spooky-completion.ps1
```

#### **Help Text**
```bash
$ spooky completion generate --help
Usage: spooky completion generate [--shell SHELL] [--output FILE]

Generates shell completion scripts for spooky commands.

Examples:
  spooky completion generate                    # auto-detect current shell, output to stdout
  spooky completion generate --shell bash       # generate for bash, output to stdout
  spooky completion generate --shell zsh        # generate for zsh, output to stdout
  spooky completion generate --shell bash --output ~/spooky-bash-completion.sh  # save to file

Installation:
  # For bash (direct installation)
  spooky completion generate --shell bash >> ~/.bashrc
  
  # For bash (review first)
  spooky completion generate --shell bash --output ~/spooky-bash-completion.sh
  # Review the file, then: source ~/spooky-bash-completion.sh
  
  # For zsh (direct installation)
  spooky completion generate --shell zsh >> ~/.zshrc
  
  # For zsh (review first)
  spooky completion generate --shell zsh --output ~/spooky-zsh-completion.sh
  # Review the file, then: source ~/spooky-zsh-completion.sh
  
  # For fish (review first)
  spooky completion generate --shell fish --output $XDG_CONFIG_HOME/fish/completions/spooky.fish
  
  # For PowerShell (review first)
  spooky completion generate --shell powershell --output spooky-completion.ps1
  # Review the file, then add to PowerShell profile: . spooky-completion.ps1
```

## Evaluation Against clig.dev Guidelines

### ✅ **Strengths**

#### **1. Human-First Design**
- **Clear hierarchy**: The `noun verb` pattern creates an intuitive mental model
- **Discoverable**: Users can easily guess available commands (`spooky facts --help`)
- **Consistent**: All commands follow the same pattern, reducing cognitive load

#### **2. Simple Parts That Work Together**
- **Modular design**: Each noun represents a distinct domain (facts, actions, machines, etc.)
- **Composable**: Commands can be chained or used in scripts
- **Single responsibility**: Each verb has a clear, focused purpose

#### **3. Consistency Across Programs**
- **Follows established patterns**: Similar to `git <noun> <verb>` and `docker <noun> <verb>`
- **Predictable structure**: Users familiar with modern CLIs will find this intuitive
- **Standard argument ordering**: `<command> <target> <arguments>`

#### **4. Ease of Discovery**
- **Self-documenting**: Command structure reveals available functionality
- **Helpful help**: `spooky facts --help` clearly shows available facts operations
- **Logical grouping**: Related operations are grouped under the same noun

### ⚠️ **Areas for Improvement**

#### **1. Command Length and Typing**
- **Longer commands**: `spooky facts gather` vs current `spooky gather-facts`
- **More typing required**: Users will type `spooky` + noun + verb frequently
- **Potential for aliases**: Consider common shortcuts for frequent operations

#### **2. Verb Consistency**
- **Mixed verb patterns**: Some use `list` (facts, actions, machines), others use `validate` (facts, actions, templates)
- **Inconsistent naming**: `act` vs `gather` vs `list` - verbs should be consistent across nouns

#### **3. Argument Positioning**
- **Project directory placement**: Currently at the end, but might be clearer as a required positional argument
- **Flag consistency**: `--shell` for completion vs no flags for other commands

## Recommended Refinements

### **1. Standardize Verb Patterns**

```
# Consistent verb patterns across all nouns:
spooky facts {gather, list, validate} <project directory>
spooky actions {list, validate, run} <project directory>    # 'run' instead of 'act'
spooky machines {list, validate} <project directory>        # add 'validate'
spooky templates {list, validate, render} <project directory>
spooky project {show, validate, init} <project directory>
spooky completion {generate} --shell
```

### **2. Consider Common Aliases**

```bash
# For frequently used commands, consider aliases:
spooky run <project directory>     # alias for 'spooky actions run'
spooky ls <project directory>      # alias for 'spooky project show'
spooky init <project directory>    # alias for 'spooky project init'
```

### **3. Improve Argument Structure**

```bash
# Make project directory a required positional argument:
spooky facts gather <project directory>
spooky actions list <project directory>
spooky templates render <project> <relative template path> [--variables variables/] [--machine name|ID] [--dry-run] [--preview]

# For commands that don't need project context:
spooky completion generate --shell
spooky version
spooky help
```

## Implementation Considerations

### **1. Clean Implementation**
- **Direct replacement**: Replace current commands with new noun-verb structure
- **No legacy support**: Remove old command patterns entirely
- **Documentation updates**: Update all examples and documentation to reflect new structure

### **2. Help System Design**
```bash
# Main help should show the noun hierarchy:
$ spooky --help
Usage: spooky <noun> <verb> [options]

Core Execution:
  actions     Manage and run actions
  facts       Manage machine facts and data
  machines    Manage machine inventory and operations
  templates   Manage template files
  variables   Manage project variables and configuration
  project     Manage project configuration

Configuration & System:
  config      Manage configuration settings
  logs        Manage log files

System & Utilities:
  status      Check system status
  help        Show help documentation
  completion  Generate shell completions

Global Flags:
  --version   Show version information
  --help      Show this help message

# Noun-specific help should show available verbs:
$ spooky facts --help
Usage: spooky facts <verb> <project directory>

Verbs:
  gather    Collect facts from machines
  list      List available facts
  validate  Validate facts data
  export    Export facts to various formats

$ spooky machines --help
Usage: spooky machines <verb> <project directory>

Verbs:
  list      List all machines in inventory
  validate  Validate machine inventory configuration
  ping      Test connectivity to machines (DNS + ping + SSH)
  connect   Connect to a specific machine via SSH
  import    Import machine inventory from external sources
  export    Export machine inventory to JSON format

Examples:
  spooky machines list <project directory>
  spooky machines ping <project directory>
  spooky machines connect web-server <project directory>
  spooky machines import --from kubernetes k8s-nodes.json <project directory>
  spooky machines export --output machines.json <project directory>

$ spooky project --help
Usage: spooky project <verb> <project directory>

Verbs:
  validate  Validate project structure and configuration
  init      Initialize a new project
  show      Show detailed project information

Examples:
  spooky project init my-app
  spooky project validate <project directory>
  spooky project show <project directory>

# For discovering projects, use standard Unix tools:
find /path/to/projects -name "project.hcl" -exec dirname {} \;
find /path/to/projects -name "project.hcl" -exec sh -c 'echo "Project: $(basename $(dirname {}))"; echo "Path: {}"' \;
```

### **3. Error Handling**
- **Clear error messages**: "Unknown noun 'foo'. Available nouns: facts, actions, machines, templates, project"
- **Suggestions**: "Did you mean 'spooky facts list'?" when user types `spooky fact list`
- **Context-aware help**: Show relevant help when user makes common mistakes

### **4. Configuration Integration**
```bash
# Support for configuration-driven defaults:
spooky facts gather                    # uses current directory as project
spooky actions run --project ./myapp   # explicit project specification
spooky --config /path/to/config.hcl facts gather  # custom config
```

### **Proposed (Noun-Verb Pattern)**
```bash
# Core operations
spooky project init <project directory>
spooky project validate <project directory>
spooky project show <project directory>
spooky facts gather <project directory>
spooky templates validate <template>

# Enhanced functionality
spooky config show
spooky machines export <project directory>
spooky machines connect <machine> <project directory>
spooky logs tail <project directory>
spooky status check <project directory>

# Global flags
spooky --version
spooky --help
```

## Benefits of the New Design

1. **Better discoverability**: Users can explore available functionality by noun
2. **Clearer mental model**: Each noun represents a domain concept
3. **Easier to extend**: New verbs can be added to existing nouns
4. **More intuitive**: Follows modern CLI conventions (git, docker, kubectl)
5. **Better help system**: Hierarchical help structure
6. **Comprehensive coverage**: All major Spooky operations are covered by logical nouns
7. **Consistent patterns**: All commands follow the same noun-verb structure
8. **Logical grouping**: Related operations are grouped under intuitive nouns

## Noun Coverage Analysis

### **Core Execution (6 nouns)**
- `actions` - Action execution and management
- `facts` - Machine facts and data collection
- `machines` - Machine inventory operations (list, validate, ping, connect, import, export)
- `templates` - Template file management
- `variables` - Project variables and configuration management (validate, list, show, export)
- `project` - Project configuration management (validate, init, show)

### **Configuration & System (2 nouns)**
- `config` - Configuration management
- `logs` - Log file management

### **System & Utilities (3 nouns)**
- `status` - System health monitoring
- `help` - Help system
- `completion` - Shell completion generation

### **Aliases (3 commands)**
- `run` - Quick action execution
- `ls` - Quick project information display
- `init` - Quick project initialization

### **Global Flags (2 flags)**
- `--version` - Show version information
- `--help` - Show help

## Project Discovery

Since `spooky project list` is not supported, users should use standard Unix tools to discover projects:

```bash
# Find all Spooky projects in a directory
find /path/to/projects -name "project.hcl" -exec dirname {} \;

# Find projects with more details
find /path/to/projects -name "project.hcl" -exec sh -c 'echo "Project: $(basename $(dirname {}))"; echo "Path: {}"' \;

# Find projects and show their info
find /path/to/projects -name "project.hcl" -exec sh -c 'echo "=== $(basename $(dirname {})) ==="; spooky project show {}' \;

# List projects in current directory
find . -maxdepth 2 -name "project.hcl" -exec dirname {} \;
```

This approach:
- **Follows Unix philosophy** - Use standard tools for discovery
- **Reduces complexity** - No need to implement directory scanning in Spooky
- **More flexible** - Users can use `find`, `ls`, `grep`, etc. as needed
- **Better performance** - Standard tools are optimized for file system operations

## Design Decisions

### **Why No `files`

#### **5. Consistent with Spooky Philosophy**
Spooky is about **managing remote infrastructure through actions**, not about generic file operations. File operations should be part of the action workflow.

### **Why No `ssh` Noun**

The `ssh` noun was intentionally excluded from the CLI design for the following reasons:

#### **1. SSH Operations Are Machine-Specific**
SSH connections are always to specific machines, so they belong under the `machines` noun:
```bash
# Instead of generic SSH operations:
spooky ssh connect [machine] <project directory>
spooky ssh test [machine] <project directory>

# Use machine-specific operations:
spooky machines connect [machine] <project directory>
spooky machines ping <project directory>  # includes SSH testing
```

#### **2. SSH Key Generation Issues**
- **OpenSSL compatibility problems** - Previous issues with SSH key generation
- **Users prefer OpenSSH** - Most users want to use standard OpenSSH tools (`ssh-keygen`)
- **Not core to Spooky** - Key generation isn't Spooky's primary purpose

#### **3. Redundant with `machines ping`**
- **`spooky machines ping`** already tests SSH connectivity as part of its comprehensive check
- **`spooky ssh test`** would be redundant and confusing
- **`spooky machines ping`** does DNS lookup + ping + SSH availability

#### **4. Let Users Handle SSH Keys**
- **No key generation** in Spooky
- **Users manage keys** with standard OpenSSH tools
- **Spooky uses existing keys** from standard locations (`~/.ssh/`)

#### **5. Consistent with Machine Management**
SSH operations are fundamentally about connecting to and managing machines, so they logically belong under the `machines` noun.

### **Why No `version` Noun**

The `version` noun was intentionally excluded from the CLI design for the following reasons:

#### **1. Standard CLI Convention**
Version information is universally provided via the `--version` flag:
```bash
# Standard approach (recommended):
spooky --version

# Instead of noun-verb approach:
spooky version show
```

#### **2. Avoids Confusion with Verbose Flag**
- **`-v` is commonly used for verbose output** in many CLI tools
- **`--version` is the standard** for version information
- **No ambiguity** about what the flag does

#### **3. Simpler Implementation**
- **Single global flag** instead of a noun with verbs
- **Standard cobra behavior** - no custom implementation needed
- **Consistent with other tools** - users expect `--version`

#### **4. Limited Use Cases**
Version information typically only needs to be displayed, not managed:
- **Show version**: `spooky --version`
- **Check for updates**: Not typically done via CLI (use package managers)
- **Update**: Not typically done via CLI (use package managers)

#### **5. Follows clig.dev Guidelines**
The [clig.dev](https://clig.dev/) guidelines recommend using `--version` for version information, which is the established convention across CLI tools.

#### **Schema-Based Import Implementation**

The `spooky machines import` command will use well-defined schemas for different data sources:

##### **Kubernetes Schema**
```bash
spooky machines import --from kubernetes k8s-nodes.json <project directory>
```

**Schema Source**: Kubernetes OpenAPI v3 specification
- **Official Schema**: `https://raw.githubusercontent.com/kubernetes/kubernetes/master/api/openapi-spec/v3/api__v1_openapi.json`
- **Node List Schema**: `io.k8s.api.core.v1.NodeList`
- **Node Schema**: `io.k8s.api.core.v1.Node`

**Data Extraction**:
```json
{
  "apiVersion": "v1",
  "kind": "NodeList",
  "items": [
    {
      "metadata": {
        "name": "worker-node-1",
        "labels": {
          "node-role.kubernetes.io/worker": "true",
          "environment": "production"
        }
      },
      "status": {
        "addresses": [
          {"type": "InternalIP", "address": "10.0.1.10"},
          {"type": "Hostname", "address": "worker-node-1"}
        ],
        "conditions": [
          {"type": "Ready", "status": "True"}
        ]
      }
    }
  ]
}
```

**Transformation to Spooky Format**:
```hcl
inventory {
  machine "worker-node-1" {
    host = "10.0.1.10"
    user = "admin"  # from project defaults
    tags = {
      role = "worker"
      environment = "production"
      source = "kubernetes"
    }
  }
}
```

##### **Podman Schema**
```bash
spooky machines import --from podman podman-containers.json <project directory>
```

**Schema Source**: Podman container inspection format
- **Command**: `podman ps --format json`
- **Schema**: Podman container inspection schema

**Data Extraction**:
```json
[
  {
    "Id": "173363d8a146717d81be76dad38c1ef573981e91724fe438674bf959fe294439",
    "Names": ["web-server"],
    "Image": "localhost/spookylab-ssh:latest",
    "ExposedPorts": {
      "22": ["tcp"],
      "80": ["tcp"]
    },
    "NetworkSettings": {
      "IPAddress": "172.17.0.2",
      "Gateway": "172.17.0.1"
    },
    "State": "running"
  }
]
```

**Transformation to Spooky Format**:
```hcl
inventory {
  machine "web-server" {
    host = "172.17.0.2"
    port = 22
    user = "admin"  # from project defaults
    tags = {
      source = "podman"
      container_id = "173363d8a146"
      image = "localhost/spookylab-ssh:latest"
    }
  }
}
```

##### **CMDB Systems**

Various CMDB systems have different schemas. Here are some common examples:

###### **ServiceNow CMDB**
```bash
spooky machines import --from servicenow servicenow-servers.json <project directory>
```

**Schema**: ServiceNow CMDB API format
```json
{
  "result": [
    {
      "sys_id": "abc123def456",
      "name": "prod-web-01",
      "ip_address": "192.168.1.10",
      "sys_class_name": "cmdb_ci_server",
      "environment": "production",
      "operational_status": "1",
      "install_status": "1",
      "os": "Linux",
      "os_version": "Ubuntu 22.04"
    }
  ]
}
```

###### **Ansible Tower/AWX Inventory**
```bash
spooky machines import --from ansible-tower tower-inventory.json <project directory>
```

**Schema**: Ansible Tower inventory format
```json
{
  "results": [
    {
      "id": 1,
      "name": "web-server-01",
      "description": "Production web server",
      "variables": {
        "ansible_host": "192.168.1.10",
        "ansible_user": "admin",
        "ansible_port": 22
      },
      "groups": ["web-servers", "production"]
    }
  ]
}
```

###### **Terraform State**
```bash
spooky machines import --from terraform terraform.tfstate <project directory>
```

**Schema**: Terraform state file format
```json
{
  "version": 4,
  "terraform_version": "1.5.0",
  "resources": [
    {
      "type": "aws_instance",
      "name": "web_server",
      "instances": [
        {
          "attributes": {
            "id": "i-1234567890abcdef0",
            "public_ip": "54.123.45.67",
            "private_ip": "10.0.1.10",
            "tags": {
              "Name": "web-server-01",
              "Environment": "production"
            }
          }
        }
      ]
    }
  ]
}
```

###### **AWS EC2 Instances**
```bash
spooky machines import --from aws-ec2 aws-instances.json <project directory>
```

**Schema**: AWS EC2 describe-instances format
```json
{
  "Reservations": [
    {
      "Instances": [
        {
          "InstanceId": "i-1234567890abcdef0",
          "PublicIpAddress": "54.123.45.67",
          "PrivateIpAddress": "10.0.1.10",
          "State": {"Name": "running"},
          "Tags": [
            {"Key": "Name", "Value": "web-server-01"},
            {"Key": "Environment", "Value": "production"}
          ],
          "InstanceType": "t3.micro"
        }
      ]
    }
  ]
}
```

###### **Azure Virtual Machines**
```bash
spooky machines import --from azure-vm azure-vms.json <project directory>
```

**Schema**: Azure VM list format
```json
{
  "value": [
    {
      "id": "/subscriptions/.../virtualMachines/web-vm-01",
      "name": "web-vm-01",
      "properties": {
        "networkProfile": {
          "networkInterfaces": [
            {
              "properties": {
                "ipConfigurations": [
                  {
                    "properties": {
                      "privateIPAddress": "10.0.1.10",
                      "publicIPAddress": {
                        "properties": {
                          "ipAddress": "20.123.45.67"
                        }
                      }
                    }
                  }
                ]
              }
            }
          ]
        }
      },
      "tags": {
        "Environment": "production",
        "Role": "web"
      }
    }
  ]
}
```

###### **Google Cloud Compute Instances**
```bash
spooky machines import --from gcp-compute gcp-instances.json <project directory>
```

**Schema**: GCP compute instances format
```json
{
  "items": [
    {
      "id": "1234567890123456789",
      "name": "web-instance-01",
      "networkInterfaces": [
        {
          "networkIP": "10.0.1.10",
          "accessConfigs": [
            {
              "natIP": "34.123.45.67"
            }
          ]
        }
      ],
      "labels": {
        "environment": "production",
        "role": "web"
      },
      "status": "RUNNING"
    }
  ]
}
```

###### **Custom CMDB Format**
```bash
spooky machines import --from custom custom-cmdb.json <project directory>
```

**Schema**: User-defined format
```json
{
  "machines": [
    {
      "name": "web-server-01",
      "host": "192.168.1.10",
      "user": "admin",
      "port": 22,
      "tags": {
        "environment": "production",
        "role": "web",
        "datacenter": "us-west-1"
      }
    }
  ]
}
```

##### **Standard JSON Format**
```bash
spooky machines import machines-standard.json <project directory>
```

#### **3. Template Management**

The `spooky templates` command provides template analysis and rendering capabilities:

```bash
spooky templates list <project directory>                                    # List available templates
spooky templates validate <template> <project directory>                     # Validate template syntax
spooky templates render <project> <relative template path>                   # Render template with real data
spooky templates render <project> <relative template path> --dry-run         # Show what would be rendered
spooky templates render <project> <relative template path> --preview         # Show template analysis + sample
```

**Key Differences:**
- **`--dry-run`**: Shows what the template would render with real project data (no file writing)
- **`--preview`**: Shows template analysis, variable detection, and sample rendering with mock data

**Examples:**
```bash
# Basic rendering (automatically loads all HCL files from data/ directory)
spooky templates render ./my-project templates/nginx.conf.tmpl

# With specific data file
spooky templates render ./my-project templates/config.tmpl --data data/variables.hcl

# With data directory (loads all HCL files from directory)
spooky templates render ./my-project templates/config.tmpl --data data/

# With output file
spooky templates render ./my-project templates/nginx.conf.tmpl --output /etc/nginx/nginx.conf

# Preview mode
spooky templates render ./my-project templates/config.tmpl --preview

# Dry run with real data
spooky templates render ./my-project templates/config.tmpl --data data/variables.hcl --dry-run

# With server facts
spooky templates render ./my-project templates/config.tmpl --machine web-001 --data data/facts.hcl
```

**Data File Support:**
- **Automatic Loading**: All `.hcl` files in `<project>/data/` are automatically loaded
- **File-based Namespaces**: Each HCL file becomes a namespace (e.g., `variables.app_name`, `environment.datacenter`)
- **Single File**: Use `--data data/file.hcl` to load a specific file
- **Directory**: Use `--data data/` to load all HCL files from a directory
- **HCL Format**: Data files use HCL variable syntax with `variable` blocks and `default` values

**Example Output for `--preview`:**
```bash
$ spooky templates render ./my-project templates/nginx.conf.tmpl --preview

=== TEMPLATE PREVIEW ===
Template: nginx.conf.tmpl
Full path: ./my-project/templates/nginx.conf.tmpl
Size: 1024 bytes

--- Template Analysis ---
Variables/Functions detected:
  - hostname
  - port
  - server_name
  - custom.app_name
  - system.os_version

Available template functions:
  - custom(path)     - Access custom facts
  - system(path)     - Access system facts
  - env(key)         - Access environment variables
  - data(path)       - Access additional data

--- Sample Rendering (with mock data) ---
server {
    listen 8080;
    server_name example.com;
    root /var/www/html;
}
=== END PREVIEW ===
```

## CLI Examples by Noun

This section provides comprehensive examples of all CLI compositions, organized by noun and then by noun-verb combinations.

### **1. Actions Noun**

Actions manage and run deployment and configuration tasks.

#### **spooky actions list**
```bash
# List all actions in a project
spooky actions list ./my-project

# List actions with details
spooky actions list ./my-project --verbose

# List actions matching a pattern
spooky actions list ./my-project --filter "deploy*"

# List actions for specific machines
spooky actions list ./my-project --machine "web-001"

# List actions with specific tags
spooky actions list ./my-project --tags "production,web"

# Complex filtering
spooky actions list ./my-project --filter "name:deploy* AND tags:production"
```

#### **spooky actions validate**
```bash
# Validate all actions in a project
spooky actions validate ./my-project

# Validate specific action files
spooky actions validate ./my-project --file actions/deploy.hcl

# Validate with detailed output
spooky actions validate ./my-project --verbose
```

#### **spooky actions run**
```bash
# Run all actions in a project
spooky actions run ./my-project

# Run specific actions
spooky actions run ./my-project --action "deploy-web"

# Run actions on specific machines
spooky actions run ./my-project --machine "web-001"

# Run actions with tags
spooky actions run ./my-project --tags "production,web"

# Run actions in parallel
spooky actions run ./my-project --parallel

# Run actions with timeout
spooky actions run ./my-project --timeout 300

# Complex filtering
spooky actions run ./my-project --filter "name:deploy* AND tags:production"

# Dry run mode
spooky actions run ./my-project --dry-run

# Dry run specific actions
spooky actions run ./my-project --action "deploy-web" --dry-run
```

### **2. Facts Noun**

Facts manage machine data collection and storage.

#### **spooky facts gather**
```bash
# Gather facts from all machines
spooky facts gather ./my-project

# Gather facts from specific machines
spooky facts gather ./my-project --machine "web-001"

# Gather facts with specific tags
spooky facts gather ./my-project --tags "production"

# Gather specific fact types
spooky facts gather ./my-project --facts "os,network,hardware"

# Gather facts with custom collectors
spooky facts gather ./my-project --collector "custom-app"

# Gather facts with timeout
spooky facts gather ./my-project --timeout 60

# Gather facts with limited parallelism
spooky facts gather ./my-project --parallel 10

# Gather facts with conservative parallelism
spooky facts gather ./my-project --parallel 2

# Complex filtering
spooky facts gather ./my-project --filter "tags:production AND hostname:web*"

# Complex filtering with parallel execution
spooky facts gather ./my-project --filter "tags:production" --parallel 5
```

#### **spooky facts list**
```bash
# List all available facts
spooky facts list ./my-project

# List facts for specific machine
spooky facts list ./my-project --machine "web-001"

# List facts with specific tags
spooky facts list ./my-project --tags "production"

# List facts by category
spooky facts list ./my-project --category "system"

# List facts with details
spooky facts list ./my-project --verbose

# Complex filtering
spooky facts list ./my-project --filter "tags:production AND os:linux"
```

#### **spooky facts validate**
```bash
# Validate stored facts data integrity
spooky facts validate ./my-project

# Gather fresh facts and compare with stored
spooky facts validate ./my-project --compare

# Validate with detailed output
spooky facts validate ./my-project --verbose
```

#### **spooky facts export**
```bash
# Export facts to JSON (format must be specified)
spooky facts export ./my-project --format json

# Export facts to HCL
spooky facts export ./my-project --format hcl

# Export specific facts
spooky facts export ./my-project --facts "os,network" --format json

# Export to file
spooky facts export ./my-project --output facts.json --format json

# Export with filtering
spooky facts export ./my-project --machine "web-001" --format json --output web-001-facts.json
```

### **3. Machines Noun**

Machines manage inventory and connectivity operations.

#### **spooky machines list**
```bash
# List all machines in inventory
spooky machines list ./my-project

# List machines with details
spooky machines list ./my-project --verbose

# List machines by tags
spooky machines list ./my-project --tags "production"

# List machines with filtering
spooky machines list ./my-project --filter "tags:production AND hostname:web*"
```

#### **spooky machines validate**
```bash
# Validate machine inventory configuration
spooky machines validate ./my-project

# Validate with detailed output
spooky machines validate ./my-project --verbose
```

#### **spooky machines ping**
```bash
# Test connectivity to all machines
spooky machines ping ./my-project

# Test specific machine
spooky machines ping ./my-project --machine "web-001"

# Test machines with specific tags
spooky machines ping ./my-project --tags "production"

# Test with different protocols
spooky machines ping ./my-project --protocol "dns,ping,ssh"

# Test with timeout
spooky machines ping ./my-project --timeout 30

# Test with limited parallelism
spooky machines ping ./my-project --parallel 10

# Test with conservative parallelism
spooky machines ping ./my-project --parallel 2

# Complex filtering
spooky machines ping ./my-project --filter "tags:production AND hostname:web*" --protocol "dns,ping"

# Complex filtering with parallel execution
spooky machines ping ./my-project --filter "tags:production" --parallel 5
```

#### **spooky machines connect**
```bash
# Connect to specific machine
spooky machines connect web-001 ./my-project

# Connect with specific user
spooky machines connect web-001 ./my-project --user admin

# Connect with port
spooky machines connect web-001 ./my-project --port 2222

# Connect and run command (for complex scenarios)
spooky machines connect web-001 ./my-project --command "register-myself https://cmdb.example.com/"
```

#### **spooky machines import**
```bash
# Import from Kubernetes
spooky machines import --from kubernetes k8s-nodes.json ./my-project

# Import from AWS EC2
spooky machines import --from aws-ec2 aws-instances.json ./my-project

# Import from Terraform state
spooky machines import --from terraform terraform.tfstate ./my-project

# Import from custom format
spooky machines import --from custom custom-inventory.json ./my-project
```

#### **spooky machines export**
```bash
# Export to JSON format (format must be specified)
spooky machines export ./my-project --format json --output machines.json

# Export to HCL format
spooky machines export ./my-project --format hcl --output machines.hcl

# Export specific machines
spooky machines export ./my-project --machine "web-001" --format json --output web-001.json

# Export with filtering
spooky machines export ./my-project --tags "production" --format json --output production-machines.json

# Complex filtering
spooky machines export ./my-project --filter "tags:production AND hostname:web*" --format hcl --output web-servers.hcl
```

### **4. Templates Noun**

Templates manage template files and rendering.

#### **spooky templates list**
```bash
# List all templates in project
spooky templates list ./my-project

# List templates with details
spooky templates list ./my-project --verbose

# List templates by pattern
spooky templates list ./my-project --pattern "*.conf.tmpl"
```

#### **spooky templates validate**
```bash
# Validate all templates in project
spooky templates validate ./my-project

# Validate specific template
spooky templates validate ./my-project --template templates/nginx.conf.tmpl

# Validate multiple templates
spooky templates validate ./my-project --templates templates/*.tmpl

# Validate with detailed output
spooky templates validate ./my-project --verbose
```

#### **spooky templates render**
```bash
# Basic rendering
spooky templates render ./my-project templates/nginx.conf.tmpl

# Render with specific data file
spooky templates render ./my-project templates/config.tmpl --data data/variables.hcl

# Render with data directory
spooky templates render ./my-project templates/config.tmpl --data data/

# Render with machine facts
spooky templates render ./my-project templates/config.tmpl --machine web-001

# Render with output file
spooky templates render ./my-project templates/nginx.conf.tmpl --output /etc/nginx/nginx.conf

# Preview mode
spooky templates render ./my-project templates/config.tmpl --preview

# Dry run mode
spooky templates render ./my-project templates/config.tmpl --dry-run

# Combined example
spooky templates render ./my-project templates/config.tmpl --machine web-001 --data data/facts.hcl --dry-run
```

### **5. Variables Noun**

Variables manage project variables and configuration values. See [Variables System Plan](../variables-system.md) for detailed implementation details.

#### **spooky variables validate**
```bash
# Validate all variable files in project
spooky variables validate ./my-project

# Validate with detailed output
spooky variables validate ./my-project --verbose

# Validate specific variable file
spooky variables validate ./my-project --file variables/environment.hcl

# Validate multiple specific files
spooky variables validate ./my-project --files variables/environment.hcl,variables/secrets.hcl

# Validate with custom schema
spooky variables validate ./my-project --schema custom-schema.hcl

# Validate with strict mode (treat warnings as errors)
spooky variables validate ./my-project --strict

# Validate and show variable dependencies
spooky variables validate ./my-project --show-dependencies

# Validate with format output
spooky variables validate ./my-project --format json --output validation-report.json
```

#### **spooky variables list**
```bash
# List all variables in project
spooky variables list ./my-project

# List variables with details
spooky variables list ./my-project --verbose

# List variables by type
spooky variables list ./my-project --type string

# List variables by file
spooky variables list ./my-project --file variables/environment.hcl

# List required variables only
spooky variables list ./my-project --required

# List sensitive variables only
spooky variables list ./my-project --sensitive

# List variables with filtering
spooky variables list ./my-project --filter "name:app*"

# List with format output
spooky variables list ./my-project --format json --output variables.json
```

#### **spooky variables show**
```bash
# Show specific variable details
spooky variables show app_name ./my-project

# Show variable with full context
spooky variables show app_name ./my-project --verbose

# Show variable dependencies
spooky variables show database_url ./my-project --dependencies

# Show variable usage across project
spooky variables show app_name ./my-project --usage

# Show variable validation rules
spooky variables show server_count ./my-project --validation
```

#### **spooky variables export**
```bash
# Export all variables to JSON
spooky variables export ./my-project --format json --output variables.json

# Export variables to HCL
spooky variables export ./my-project --format hcl --output variables.hcl

# Export specific variables
spooky variables export ./my-project --variables "app_name,environment" --format json

# Export by type
spooky variables export ./my-project --type string --format json

# Export non-sensitive variables only
spooky variables export ./my-project --exclude-sensitive --format json

# Export with resolved values
spooky variables export ./my-project --resolve --format json
```

**Variable File Locations:**
- `variables.hcl` in project root
- `variables/*.hcl` files in variables directory

**Validation Features:**
- File location validation (only project directory files)
- HCL syntax validation (treats invalid HCL as errors)
- Variable type validation (string, number, bool, list, map)
- Custom validation condition evaluation
- Dependency and circular reference detection
- Detailed error reporting with file and line numbers
- Schema validation against embedded schemas
- Cross-variable reference validation

**Variable Types Supported:**
- `string` - Text values
- `number` - Numeric values (integers and floats)
- `bool` - Boolean values (true/false)
- `list` - Ordered lists of values
- `map` - Key-value mappings

**Example Output:**
```bash
$ spooky variables validate ./my-project
Validating variables for project: ./my-project
✓ variables.hcl: 5 variables validated
✓ variables/environment.hcl: 3 variables validated
✓ variables/secrets.hcl: 2 variables validated
✅ Variables validation passed (10 variables in 3 files)

$ spooky variables validate ./my-project --verbose
Validating variables for project: ./my-project
Reading variables.hcl...
  ✓ app_name: string, default="myapp"
  ✓ environment: string, default="development"
  ✓ server_count: number, default=3
  ✓ ssl_enabled: bool, default=false
  ✓ database_config: map, default={host="localhost", port=5432}
Reading variables/environment.hcl...
  ✓ log_level: string, default="info"
  ✓ debug_mode: bool, default=false
  ✓ timeout: number, default=30
Reading variables/secrets.hcl...
  ✓ database_url: string, required=true, sensitive=true
  ✓ api_key: string, required=true, sensitive=true
✅ Variables validation passed (10 variables in 3 files)

$ spooky variables list ./my-project --verbose
Variables in ./my-project:
File: variables.hcl
  app_name: string, default="myapp", description="Application name"
  environment: string, default="development", description="Deployment environment"
  server_count: number, default=3, description="Number of servers"
  ssl_enabled: bool, default=false, description="Enable SSL"
  database_config: map, default={host="localhost", port=5432}, description="Database configuration"

File: variables/environment.hcl
  log_level: string, default="info", description="Logging level"
  debug_mode: bool, default=false, description="Debug mode"
  timeout: number, default=30, description="Operation timeout"

File: variables/secrets.hcl
  database_url: string, required=true, sensitive=true, description="Database URL"
  api_key: string, required=true, sensitive=true, description="API key"

Total: 10 variables in 3 files

$ spooky variables show app_name ./my-project --verbose
Variable: app_name
  Type: string
  Default: "myapp"
  Description: Application name
  Required: false
  Sensitive: false
  Source: variables.hcl:2
  Validation: none
  Dependencies: none
  Usage: 
    - actions/deploy.hcl:5 (command interpolation)
    - templates/nginx.conf.tmpl:3 (template variable)

$ spooky variables export ./my-project --format json --output variables.json
Exported 10 variables to variables.json
  - 8 variables with defaults
  - 2 required variables
  - 2 sensitive variables

$ spooky variables validate ./my-project
Error: invalid HCL in variables/environment.hcl: unexpected '=' at line 3
Error: variable 'server_count' in variables/defaults.hcl: type 'string' does not match value 3
Error: required variable 'database_url' in variables/secrets.hcl has no default value
Error: circular dependency detected: app_name → database_url → app_name
❌ Variables validation failed (4 errors)
```

### **6. Project Noun**

Project manages project configuration and structure.

#### **spooky project validate**
```bash
# Validate project structure
spooky project validate ./my-project

# Validate with detailed output
spooky project validate ./my-project --verbose
```

#### **spooky project init**
```bash
# Initialize new project
spooky project init my-new-project

# Initialize with custom settings
spooky project init my-new-project --user admin --port 22
```

#### **spooky project show**
```bash
# Show project information
spooky project show ./my-project

# Show with details
spooky project show ./my-project --verbose

# Show specific information
spooky project show ./my-project --info "machines,actions"
```

### **7. Config Noun**

Config manages global and project configuration.

#### **spooky config list**
```bash
# List all configuration
spooky config list

# List specific configuration
spooky config list --section "ssh"

# List with details
spooky config list --verbose
```

#### **spooky config validate**
```bash
# Validate configuration file
spooky config validate /path/to/config.hcl

# Validate with detailed output
spooky config validate /path/to/config.hcl --verbose
```

#### **spooky config show**
```bash
# Show configuration file path
spooky config show

# Show with details
spooky config show --verbose
```

### **8. Logs Noun**

Logs manage log files and output.

#### **spooky logs list**
```bash
# List log files
spooky logs list ./my-project

# List with details
spooky logs list ./my-project --verbose

# List by pattern
spooky logs list ./my-project --pattern "*.log"
```

#### **spooky logs show**
```bash
# Show log content
spooky logs show ./my-project

# Show specific log file
spooky logs show ./my-project --file spooky.log

# Show with line numbers
spooky logs show ./my-project --line-numbers

# Show with grep filter
spooky logs show ./my-project --grep "ERROR"
```

#### **spooky logs tail**
```bash
# Tail log files
spooky logs tail ./my-project

# Tail specific file
spooky logs tail ./my-project --file spooky.log

# Tail with follow
spooky logs tail ./my-project --follow

# Tail with lines
spooky logs tail ./my-project --lines 50
```

#### **spooky logs clear**
```bash
# Clear all logs
spooky logs clear ./my-project

# Clear specific log file
spooky logs clear ./my-project --file spooky.log

# Clear with confirmation
spooky logs clear ./my-project --confirm
```

#### **spooky logs export**
```bash
# Export logs to file
spooky logs export ./my-project --output logs.txt

# Export specific log file
spooky logs export ./my-project --file spooky.log --output spooky.txt

# Export with format
spooky logs export ./my-project --format json --output logs.json
```

### **9. Help Noun**

Help provides documentation and assistance.

#### **spooky help show**
```bash
# Show help for command
spooky help show actions

# Show help for specific topic
spooky help show templates

# Show help with examples
spooky help show --examples
```

#### **spooky help search**
```bash
# Search help topics
spooky help search "template"

# Search with pattern
spooky help search --pattern "render*"

# Search with category
spooky help search --category "commands"
```

#### **spooky help examples**
```bash
# Show examples for command
spooky help examples actions run

# Show examples for topic
spooky help examples templates

# Show all examples
spooky help examples
```

### **10. Completion Noun**

Completion generates shell completion scripts.

#### **spooky completion generate**
```bash
# Generate for current shell
spooky completion generate

# Generate for specific shell
spooky completion generate --shell bash

# Generate for zsh
spooky completion generate --shell zsh

# Generate for fish
spooky completion generate --shell fish

# Generate for PowerShell
spooky completion generate --shell powershell

# Generate to file
spooky completion generate --shell bash --output ~/spooky-bash-completion.sh
```

### **11. Global Commands**

Commands that don't follow the noun-verb pattern.

#### **spooky run**
```bash
# Quick action execution (alias for spooky actions run)
spooky run ./my-project

# Run specific action
spooky run ./my-project --action "deploy-web"

# Run with machine
spooky run ./my-project --machine "web-001"

# Dry run mode
spooky run ./my-project --dry-run
```

#### **spooky ls**
```bash
# Quick project listing (alias for spooky project show)
spooky ls ./my-project

# List with details
spooky ls ./my-project --verbose
```

#### **spooky init**
```bash
# Quick project initialization (alias for spooky project init)
spooky init my-new-project

# Init with custom settings
spooky init my-new-project --user admin --port 22
```

#### **spooky ping**
```bash
# Quick connectivity test (alias for spooky machines ping)
spooky ping ./my-project

# Ping specific machine
spooky ping ./my-project --machine "web-001"

# Ping with protocols
spooky ping ./my-project --protocol "dns,ping,ssh"

# Ping with limited parallelism
spooky ping ./my-project --parallel 10

# Ping with conservative parallelism
spooky ping ./my-project --parallel 2
```

### **12. Global Flags**

Flags available on all commands.

#### **spooky --version**
```bash
# Show version information
spooky --version
```

#### **spooky --help**
```bash
# Show help
spooky --help

# Show help for specific command
spooky actions --help
spooky facts --help
spooky machines --help
```

#### **spooky --config**
```bash
# Use alternative config file
spooky --config /path/to/custom-config.hcl actions list ./my-project

# Use config file for all commands
spooky --config /path/to/custom-config.hcl machines ping ./my-project
```

## Common Usage Patterns

### **Project Discovery**
```bash
# Find all Spooky projects
find /path/to/projects -name "project.hcl" -exec dirname {} \;

# Find projects with details
find /path/to/projects -name "project.hcl" -exec sh -c 'echo "Project: $(basename $(dirname {}))"; echo "Path: {}"' \;

# Find and validate projects
find /path/to/projects -name "project.hcl" -exec sh -c 'echo "=== $(basename $(dirname {})) ==="; spooky project validate {}' \;
```

### **Workflow Examples**
```bash
# Complete deployment workflow
spooky project validate ./my-project
spooky machines ping ./my-project
spooky facts gather ./my-project
spooky templates render ./my-project templates/config.tmpl --dry-run
spooky actions run ./my-project

# Monitoring workflow
spooky machines ping ./my-project --parallel 10
spooky facts gather ./my-project --parallel 5
spooky facts export ./my-project --format json --output facts-report.json

# Development workflow
spooky project init my-new-feature
spooky templates render ./my-new-feature templates/config.tmpl --preview
spooky actions run ./my-new-feature --dry-run
spooky project validate ./my-new-feature
```