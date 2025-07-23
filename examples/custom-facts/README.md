# Custom Facts Examples

This directory contains examples of custom facts import and usage in the Spooky CLI.

## Features

- **Local JSON Import**: Import custom facts from local JSON files
- **HTTPS Import**: Import custom facts from HTTPS endpoints (HTTPS required for security)
- **Select Facts Filtering**: Filter which facts to import using patterns
- **Template Integration**: Use custom facts in configuration templates
- **Deep Merge**: Merge custom facts with existing system facts
- **Overrides**: Override system facts with custom values

## File Structure

```
examples/custom-facts/
├── README.md                    # This file
├── example-custom-facts.json    # Example custom facts file
├── simple-example.json          # Simple custom facts example
├── example-template.hcl         # Example template using custom facts
└── http-example.json           # Example for HTTP import testing
```

## Custom Facts Format

Custom facts use a server-specific format with `custom` and `overrides` sections:

```json
{
  "web-001": {
    "custom": {
      "application": {
        "name": "nginx",
        "version": "1.18.0",
        "config_path": "/etc/nginx/nginx.conf"
      },
      "environment": {
        "datacenter": "fra00",
        "rack": "A01",
        "power_zone": "PZ-1"
      },
      "monitoring": {
        "prometheus_port": 9100,
        "alert_manager": "alert.example.com"
      }
    },
    "overrides": {
      "os": {
        "name": "ubuntu",
        "version": "22.04.2"
      }
    }
  },
  "web-002": {
    "custom": {
      "application": {
        "name": "apache",
        "version": "2.4.54"
      }
    }
  }
}
```

## Import Commands

### Local File Import

```bash
# Basic import
spooky facts import --source local --file custom-facts.json

# Import with merge mode
spooky facts import --source local --file custom-facts.json --merge-mode merge

# Import specific facts only
spooky facts import --source local --file custom-facts.json --select-facts application,environment

# Import with validation
spooky facts import --source local --file custom-facts.json --validate

# Dry run import
spooky facts import --source local --file custom-facts.json --dry-run

# Import for specific server
spooky facts import --source local --file custom-facts.json --server web-001
```

### HTTPS Import

```bash
# Import from HTTPS endpoint (HTTPS is required for security)
spooky facts import --source http --url https://api.example.com/facts.json

# Import with merge mode
spooky facts import --source http --url https://api.example.com/facts.json --merge-mode merge

# Import specific facts from HTTPS
spooky facts import --source http --url https://api.example.com/facts.json --select-facts application.name,environment.datacenter
```

## Select Facts Filtering

The `--select-facts` flag allows you to filter which facts to import:

```bash
# Import only application facts
spooky facts import --file custom-facts.json --select-facts application

# Import specific fact keys
spooky facts import --file custom-facts.json --select-facts application.name,application.version

# Import facts matching wildcard patterns
spooky facts import --file custom-facts.json --select-facts *.port,*.name

# Import multiple categories
spooky facts import --file custom-facts.json --select-facts application,environment,monitoring
```

## Template Integration

Custom facts can be used in configuration templates:

### Template Functions

- `{{custom 'path'}}` - Access custom facts
- `{{system 'path'}}` - Access system facts  
- `{{env 'key'}}` - Access environment variables
- `{{data 'path'}}` - Access additional data

### Example Template Usage

```hcl
action "deploy-app" {
  description = "Deploy application with custom facts"
  
  # Use custom facts
  command = "deploy.sh --app {{custom 'application.name'}} --version {{custom 'application.version'}}"
  
  # Use environment variables
  environment = {
    DATACENTER = "{{custom 'environment.datacenter'}}"
    NODE_ENV = "{{env 'NODE_ENV'}}"
  }
  
  # Use system facts
  server_info = {
    HOSTNAME = "{{system 'hostname'}}"
    OS_NAME = "{{system 'os.name'}}"
  }
  
  servers = ["web-001"]
}
```

### Rendering Templates

```bash
# Render template with fact integration
spooky templates render template.hcl --server web-001

# Render with additional data
spooky templates render template.hcl --server web-001 --data-file data.json

# Dry run template rendering
spooky templates render template.hcl --server web-001 --dry-run

# Output to file
spooky templates render template.hcl --server web-001 --output rendered.conf
```

## Merge Modes

- **replace**: Replace existing facts (default)
- **merge**: Merge with existing facts (deep merge)
- **append**: Append to existing facts
- **select**: Select specific facts to merge

## Validation

Custom facts are validated for:
- Server ID format
- Fact structure
- Required fields
- Data types

```bash
# Validate custom facts
spooky facts validate custom-facts.json
```

## HTTPS Endpoint Requirements

For HTTPS import, the endpoint must:
- Use HTTPS protocol (HTTP is not allowed for security reasons)
- Return JSON with Content-Type: application/json
- Use the custom facts format (server-specific with custom/overrides)
- Be accessible via GET request
- Return HTTP 200 status code
- Have a valid SSL/TLS certificate

Example HTTP endpoint response:
```json
{
  "web-001": {
    "custom": {
      "application": {
        "name": "nginx",
        "version": "1.18.0"
      }
    }
  }
}
```

## Testing

Test HTTPS import with a local server:

```bash
# Start a test HTTPS server (you'll need to set up SSL certificates)
# For development, you can use tools like mkcert or self-signed certificates
python3 -m http.server 8080 --bind 127.0.0.1

# Note: HTTP URLs are not allowed for security reasons
# You must use HTTPS with valid certificates for production use
```

## Integration with Existing Features

Custom facts integrate with:
- **Fact Collection**: Custom facts are included in fact queries
- **Templates**: Available in configuration templates
- **Storage**: Persisted in the facts database
- **CLI**: Available through all fact-related commands 