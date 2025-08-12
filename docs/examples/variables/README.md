# Variables Examples

This directory contains examples of variable configurations for the spooky variables system. These examples demonstrate various patterns and best practices for organizing and managing variables.

## Example Files

### [variables-basic-config.hcl](variables-basic-config.hcl)
A basic variable configuration showing fundamental variable types and attributes.

**What it demonstrates:**
- Basic variable types (string, number, bool, list, map)
- Default values and descriptions
- Simple validation rules
- Basic constraints

**Use this when:**
- Getting started with variables
- Learning basic variable syntax
- Creating simple configurations

### [variables-multi-file.hcl](variables-multi-file.hcl)
An example of the main variables file in a multi-file configuration.

**What it demonstrates:**
- Project-level variables
- Environment validation
- Basic project structure

**Use this when:**
- Organizing variables across multiple files
- Setting up project-level configuration
- Managing environment-specific settings

### [variables-with-dependencies.hcl](variables-with-dependencies.hcl)
A complex example showing variable dependencies and advanced features.

**What it demonstrates:**
- Variable dependencies and relationships
- Complex validation rules
- Sensitive variable handling
- Environment-specific configuration
- Advanced constraint validation

**Use this when:**
- Building complex variable relationships
- Managing sensitive data
- Implementing advanced validation
- Creating production-ready configurations

## Using the Examples

### Copy and Customize

1. **Choose an appropriate example** based on your needs
2. **Copy the file** to your project directory
3. **Customize the variables** for your specific use case
4. **Test the configuration** using spooky commands

```bash
# Copy an example to your project
cp docs/examples/variables/variables-basic-config.hcl ./my-project/variables.hcl

# Test the configuration
spooky variables validate ./my-project
spooky variables list ./my-project
spooky variables resolve ./my-project
```

### Multi-File Organization

For larger projects, organize variables into multiple files:

```bash
my-project/
├── variables.hcl                    # Main variables file
└── variables/
    ├── app.hcl                      # Application variables
    ├── database.hcl                 # Database variables
    ├── security.hcl                 # Security variables
    └── monitoring.hcl               # Monitoring variables
```

### Environment-Specific Configuration

Use environment variables to override values:

```bash
# Set environment variables
export SPOOKY_VAR_ENVIRONMENT="production"
export SPOOKY_VAR_DB_HOST="prod-db.example.com"
export SPOOKY_VAR_API_KEY="your-production-api-key"

# Resolve with environment overrides
spooky variables resolve --env-override
```

## Best Practices

### Variable Naming

- Use descriptive, clear names
- Follow consistent naming conventions
- Use underscores for multi-word names
- Avoid reserved words

```hcl
# ✅ Good naming
variable "app_name" { ... }
variable "database_host" { ... }
variable "api_timeout_seconds" { ... }

# ❌ Poor naming
variable "name" { ... }
variable "host" { ... }
variable "timeout" { ... }
```

### Organization

- Group related variables together
- Use separate files for different concerns
- Keep dependencies simple and clear
- Document complex relationships

### Security

- Mark sensitive variables with `sensitive = true`
- Use `encrypted = true` for highly sensitive data
- Use environment variables for secrets
- Never commit sensitive values to version control

### Validation

- Always validate important variables
- Use appropriate constraints
- Provide clear error messages
- Test validation rules thoroughly

## Testing Your Configuration

### Basic Testing

```bash
# Validate syntax and rules
spooky variables validate ./my-project

# List all variables
spooky variables list ./my-project

# Resolve variables
spooky variables resolve ./my-project
```

### Advanced Testing

```bash
# Test with environment overrides
export SPOOKY_VAR_ENVIRONMENT="production"
spooky variables resolve ./my-project --env-override

# Test with specific context
spooky variables resolve ./my-project --context production

# Debug dependency resolution
spooky variables resolve ./my-project --show-dependencies
```

### Troubleshooting

If you encounter issues:

1. **Check syntax** with `spooky variables validate`
2. **Review error messages** carefully
3. **Test with minimal configuration** first
4. **Check for circular dependencies**
5. **Verify environment variable names**

## Integration with Other Systems

### With Templates

Variables can be used in templates for dynamic configuration:

```hcl
# In variables.hcl
variable "app_name" {
  type = "string"
  default = "my-app"
}

# In template.tmpl
# Application: {{.app_name}}
```

### With Actions

Variables can be passed to actions for dynamic execution:

```hcl
# In variables.hcl
variable "deploy_target" {
  type = "string"
  default = "staging"
}

# In actions.hcl
action "deploy" {
  command = "deploy.sh"
  environment = {
    TARGET = "${deploy_target}"
  }
}
```

### With Facts

Variables can be used to configure fact collection:

```hcl
# In variables.hcl
variable "collect_system_info" {
  type = "bool"
  default = true
}

# In facts.hcl
fact "system_info" {
  enabled = "${collect_system_info}"
  command = "uname -a"
}
```

## Next Steps

After reviewing these examples:

1. **Start with a basic configuration** and expand gradually
2. **Organize variables logically** as your project grows
3. **Implement proper validation** for all important variables
4. **Use environment variables** for sensitive data
5. **Test thoroughly** before deploying to production

For more detailed information, refer to:
- [Variables User Guide](../../VARIABLES_USER_GUIDE.md)
- [Variables API Reference](../../VARIABLES_API_REFERENCE.md)
- [Variables Troubleshooting Guide](../../VARIABLES_TROUBLESHOOTING.md)
