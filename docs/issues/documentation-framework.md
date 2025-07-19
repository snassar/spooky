# Implement comprehensive documentation framework

## Background
As spooky grows with new features like facts management, template processing, and remote configuration sources, we need comprehensive documentation that helps users understand, implement, and troubleshoot the tool effectively.

## Requirements

### User Documentation

#### Getting Started Guide
```markdown
# Getting Started with Spooky

## Installation
```bash
go install github.com/snassar/spooky@latest
```

## Quick Start
1. Create your first configuration file
2. Gather facts from your servers
3. Execute your first action
4. Use templates for dynamic configuration
```

#### Configuration Guide
```markdown
# Configuration Guide

## HCL Configuration Files
Spooky uses HCL (HashiCorp Configuration Language) for configuration files.

### Basic Structure
```hcl
servers {
    web-001 {
        host = "192.168.1.10"
        user = "ubuntu"
        port = 22
    }
}

actions {
    deploy-nginx {
        description = "Deploy nginx configuration"
        command = "sudo systemctl restart nginx"
        parallel = true
        timeout = 30
    }
}
```

### Advanced Configuration
- Template usage
- Fact integration
- Remote sources
- Validation rules
```

#### CLI Reference
```markdown
# CLI Reference

## Core Commands
- `spooky execute` - Execute configuration files
- `spooky validate` - Validate configurations
- `spooky list` - List resources

## Facts Management
- `spooky facts gather` - Gather server facts
- `spooky facts import` - Import facts from external sources
- `spooky facts export` - Export facts to various formats

## Template Management
- `spooky templates list` - List available templates
- `spooky templates validate` - Validate template syntax
- `spooky templates render` - Render templates locally
```

### Developer Documentation

#### API Documentation
```go
// Package spooky provides server configuration and automation capabilities
package spooky

// Config represents a spooky configuration
type Config struct {
    Servers map[string]*Server `hcl:"servers"`
    Actions map[string]*Action `hcl:"actions"`
}

// Server represents a target server
type Server struct {
    Host string `hcl:"host"`
    User string `hcl:"user"`
    Port int    `hcl:"port,optional"`
}

// Action represents an action to execute
type Action struct {
    Description string `hcl:"description"`
    Command     string `hcl:"command"`
    Parallel    bool   `hcl:"parallel,optional"`
    Timeout     int    `hcl:"timeout,optional"`
}
```

#### Architecture Documentation
```markdown
# Architecture Overview

## Components
- **CLI Layer**: Command-line interface using Cobra
- **Configuration Layer**: HCL parsing and validation
- **SSH Layer**: Remote server connectivity
- **Facts Layer**: Server fact collection and storage
- **Template Layer**: Template processing and rendering
- **Storage Layer**: BadgerDB and JSON storage backends

## Data Flow
1. User executes command
2. CLI parses arguments and flags
3. Configuration is loaded and validated
4. Facts are collected from servers
5. Templates are processed with facts
6. Actions are executed on target servers
```

### Examples and Tutorials

#### Basic Examples
```markdown
# Basic Examples

## Simple Server Configuration
```hcl
servers {
    web-server {
        host = "192.168.1.10"
        user = "ubuntu"
    }
}

actions {
    update-system {
        description = "Update system packages"
        command = "sudo apt update && sudo apt upgrade -y"
    }
}
```

## Using Facts in Templates
```hcl
actions {
    configure-nginx {
        description = "Configure nginx based on server facts"
        command = "echo '{{ .facts.os.name }}' > /tmp/os-info"
    }
}
```
```

#### Advanced Examples
```markdown
# Advanced Examples

## Multi-Server Deployment
```hcl
servers {
    web-001 { host = "192.168.1.10", user = "ubuntu" }
    web-002 { host = "192.168.1.11", user = "ubuntu" }
    db-001  { host = "192.168.1.20", user = "postgres" }
}

actions {
    deploy-web {
        description = "Deploy web application"
        command = "sudo systemctl restart nginx"
        tags = ["web", "deploy"]
    }
    
    deploy-db {
        description = "Deploy database"
        command = "sudo systemctl restart postgresql"
        tags = ["database", "deploy"]
    }
}
```

## Remote Configuration Sources
```bash
# Execute from Git repository
spooky execute github.com/org/repo/configs/

# Execute from S3 bucket
spooky execute s3://my-bucket/configs/

# Execute from HTTP endpoint
spooky execute https://example.com/configs/
```
```

### Troubleshooting Guide

#### Common Issues
```markdown
# Troubleshooting Guide

## SSH Connection Issues
### Problem: Connection timeout
**Solution**: Check SSH key permissions and server connectivity
```bash
spooky servers ping web-001
```

### Problem: Authentication failed
**Solution**: Verify SSH key path and server user
```bash
spooky --ssh-key-path ~/.ssh/id_ed25519 execute config.hcl
```

## Configuration Issues
### Problem: HCL parsing errors
**Solution**: Validate configuration syntax
```bash
spooky config validate config.hcl
```

### Problem: Template rendering errors
**Solution**: Check template syntax and fact availability
```bash
spooky templates validate template.hcl
```
```

#### Debugging Tools
```markdown
# Debugging Tools

## Verbose Output
```bash
spooky execute config.hcl --verbose
```

## Dry Run Mode
```bash
spooky execute config.hcl --dry-run
```

## Debug Logging
```bash
spooky --log-level debug execute config.hcl
```
```

## Implementation

### Phase 1: Core Documentation
1. **Create documentation structure**
2. **Write getting started guide**
3. **Create CLI reference**
4. **Add basic examples**

### Phase 2: Advanced Documentation
1. **Write configuration guide**
2. **Create troubleshooting guide**
3. **Add advanced examples**
4. **Create API documentation**

### Phase 3: Developer Documentation
1. **Write architecture documentation**
2. **Create development guide**
3. **Add contribution guidelines**
4. **Create testing documentation**

### Phase 4: Maintenance and Updates
1. **Set up documentation CI/CD**
2. **Create documentation templates**
3. **Implement versioning strategy**
4. **Add search and navigation**

## Documentation Structure

### Website Structure
```
docs/
├── index.md              # Homepage
├── getting-started/      # Getting started guides
│   ├── installation.md
│   ├── quick-start.md
│   └── first-config.md
├── user-guide/           # User documentation
│   ├── configuration.md
│   ├── cli-reference.md
│   ├── templates.md
│   └── facts.md
├── examples/             # Examples and tutorials
│   ├── basic.md
│   ├── advanced.md
│   └── real-world.md
├── troubleshooting/      # Troubleshooting guides
│   ├── common-issues.md
│   ├── debugging.md
│   └── faq.md
├── developer/            # Developer documentation
│   ├── architecture.md
│   ├── api-reference.md
│   ├── contributing.md
│   └── testing.md
└── reference/            # Reference documentation
    ├── cli.md
    ├── configuration.md
    ├── templates.md
    └── facts.md
```

### Documentation Formats

#### Markdown Documentation
- **User guides**: Step-by-step instructions
- **Examples**: Code examples with explanations
- **Troubleshooting**: Problem-solution pairs
- **Reference**: Complete API and CLI reference

#### Code Documentation
- **GoDoc**: API documentation
- **Inline comments**: Code explanations
- **Package documentation**: Overview and usage

#### Interactive Documentation
- **Command-line help**: Built-in help system
- **Interactive examples**: Hands-on tutorials
- **Validation tools**: Configuration validation

## Documentation Tools

### Static Site Generator
- **Hugo**: Fast static site generation
- **Docusaurus**: React-based documentation
- **MkDocs**: Python-based documentation

### Documentation CI/CD
```yaml
# .github/workflows/docs.yml
name: Documentation

on: [push, pull_request]

jobs:
  build-docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Build documentation
        run: hugo --source docs/
      - name: Deploy to GitHub Pages
        run: |
          # Deploy documentation
```

### Documentation Validation
```bash
# Validate markdown links
markdown-link-check docs/**/*.md

# Validate code examples
go run examples/validate-examples.go

# Check documentation coverage
go run tools/doc-coverage.go
```

## Success Criteria

### Content Requirements
- [ ] Complete getting started guide
- [ ] Comprehensive CLI reference
- [ ] Configuration guide with examples
- [ ] Troubleshooting guide
- [ ] API documentation
- [ ] Architecture documentation

### Quality Requirements
- [ ] All documentation is accurate and up-to-date
- [ ] Examples are tested and working
- [ ] Links are valid and functional
- [ ] Search functionality works correctly
- [ ] Documentation is accessible and readable

### Maintenance Requirements
- [ ] Documentation is versioned with releases
- [ ] Automated documentation builds
- [ ] Documentation review process
- [ ] User feedback collection

## Dependencies

### Documentation Tools
- **Hugo**: Static site generation
- **markdown-link-check**: Link validation
- **remark**: Markdown linting
- **prettier**: Code formatting

### CI/CD Tools
- **GitHub Actions**: Documentation deployment
- **GitHub Pages**: Documentation hosting
- **Netlify**: Alternative hosting option

## Implementation Notes

### Documentation Standards
- **Consistent formatting**: Use standard markdown
- **Clear structure**: Logical organization
- **Regular updates**: Keep documentation current
- **User feedback**: Incorporate user suggestions

### Content Guidelines
- **Clear and concise**: Easy to understand
- **Comprehensive**: Cover all features
- **Practical**: Include real-world examples
- **Accessible**: Multiple skill levels

### Maintenance Strategy
- **Regular reviews**: Monthly documentation reviews
- **Version tracking**: Document changes with releases
- **User feedback**: Collect and incorporate feedback
- **Automated testing**: Validate documentation accuracy 