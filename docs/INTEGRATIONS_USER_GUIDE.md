# Integrations System User Guide

## Overview

The spooky integrations system provides comprehensive integration capabilities for connecting with external systems, APIs, and services. This guide covers everything from basic integration setup to advanced features like custom integrations, webhooks, and third-party service connections.

**Status: Partially Implemented** - The integrations system has basic functionality but many advanced features are still in development.

## Related Documentation

- [Actions User Guide](ACTIONS_USER_GUIDE.md) - Using integrations in actions
- [Variables User Guide](VARIABLES_USER_GUIDE.md) - Integration variables and configuration
- [Logging User Guide](LOGGING_USER_GUIDE.md) - Integration logging and monitoring
- [Secrets User Guide](SECRETS_USER_GUIDE.md) - Secure integration credentials

> **See also**: [User Guides Index](USER_GUIDES_INDEX.md) - Complete overview of all user guides

## Getting Started

### Prerequisites

- spooky CLI installed and configured
- Basic understanding of API concepts
- Access to external services and APIs
- Understanding of authentication methods

### Quick Start

1. **Check Available Integration Commands**
   ```bash
   spooky integrations --help
   ```

2. **List Available Integrations**
   ```bash
   spooky integrations list
   ```

3. **Test Integration Connection**
   ```bash
   spooky integrations test --integration github
   ```

## Core Concepts

### Integration Types

spooky supports multiple integration types:

- **API Integrations** - REST and GraphQL API connections
- **Webhook Integrations** - Incoming webhook processing
- **Database Integrations** - Database connections and queries
- **Cloud Integrations** - Cloud provider services
- **Custom Integrations** - User-defined integrations

### Integration Components

The integrations system provides:

- **Authentication Management** - Secure credential storage
- **Connection Pooling** - Efficient connection management
- **Error Handling** - Robust error handling and retry logic
- **Rate Limiting** - API rate limit management

### Integration Features

Key features include:

- **Configurable Authentication** - Multiple authentication methods
- **Request/Response Handling** - Structured request and response processing
- **Data Transformation** - Convert between different data formats
- **Event Processing** - Handle real-time events and webhooks

### Integration with Other Systems

The integrations system works with other spooky systems:

- **Actions**: Use integrations in [action orchestration](ACTIONS_USER_GUIDE.md)
- **Variables**: Store integration credentials as [variables](VARIABLES_USER_GUIDE.md)
- **Logging**: Monitor integration operations through [logging](LOGGING_USER_GUIDE.md)
- **Secrets**: Secure integration credentials with [encryption](SECRETS_USER_GUIDE.md)

## Configuration

### Integration Configuration

Configure integrations in your `spooky.hcl` file:

```hcl
integrations {
  github {
    enabled = true
    api_url = "https://api.github.com"
    token = "ghp_..."
    
    rate_limit {
      requests_per_hour = 5000
      burst_limit = 100
    }
    
    webhooks {
      enabled = true
      secret = "webhook-secret"
      events = ["push", "pull_request"]
    }
  }
  
  slack {
    enabled = true
    webhook_url = "https://hooks.slack.com/services/..."
    channel = "#alerts"
    
    notifications {
      events = ["action_started", "action_completed", "action_failed"]
    }
  }
}
```

### API Integration Configuration

For API integrations:

```hcl
integrations {
  custom_api {
    enabled = true
    base_url = "https://api.example.com"
    
    authentication {
      method = "bearer"
      token = "api-token"
    }
    
    endpoints {
      users = "/users"
      projects = "/projects"
      deployments = "/deployments"
    }
  }
}
```

### Webhook Configuration

For webhook integrations:

```hcl
integrations {
  webhook {
    enabled = true
    url = "https://webhook.example.com/events"
    
    authentication {
      method = "basic"
      username = "webhook-user"
      password = "webhook-password"
    }
    
    events {
      action_started = true
      action_completed = true
      action_failed = true
    }
  }
}
```

## CLI Commands

### Integration Management

Manage integrations:

```bash
# List available integrations
spooky integrations list

# Show integration status
spooky integrations status

# Test integration connection
spooky integrations test --integration github

# Configure integration
spooky integrations configure --integration slack
```

### Integration Operations

Perform integration operations:

```bash
# Send notification
spooky integrations notify --integration slack --message "Deployment completed"

# Trigger webhook
spooky integrations webhook --integration github --event push

# Query API
spooky integrations query --integration github --endpoint /repos/owner/repo
```

### Integration Monitoring

Monitor integration health:

```bash
# Check integration health
spooky integrations health --integration github

# View integration metrics
spooky integrations metrics --integration slack

# Test integration connectivity
spooky integrations ping --integration custom_api
```

## Advanced Features

### Custom Integrations

Create custom integrations:

```hcl
integrations {
  custom_service {
    enabled = true
    type = "custom"
    
    config {
      base_url = "https://api.customservice.com"
      api_version = "v1"
      timeout_seconds = 30
    }
    
    authentication {
      method = "oauth2"
      client_id = "client-id"
      client_secret = "client-secret"
      token_url = "https://auth.customservice.com/token"
    }
    
    endpoints {
      users = "/users"
      projects = "/projects"
    }
  }
}
```

### Event Processing

Configure event processing:

```hcl
integrations {
  event_processor {
    enabled = true
    
    events {
      action_started {
        integrations = ["slack", "webhook"]
        template = "Action {{.action_name}} started on {{.machine}}"
      }
      
      action_completed {
        integrations = ["slack", "webhook"]
        template = "Action {{.action_name}} completed successfully"
      }
      
      action_failed {
        integrations = ["slack", "webhook", "email"]
        template = "Action {{.action_name}} failed: {{.error}}"
      }
    }
  }
}
```

### Data Transformation

Configure data transformation:

```hcl
integrations {
  data_transformer {
    enabled = true
    
    transforms {
      github_to_slack {
        source = "github"
        target = "slack"
        
        mapping {
          "repository.name" = "repo_name"
          "pull_request.title" = "pr_title"
          "user.login" = "author"
        }
      }
    }
  }
}
```

## Security Best Practices

### Authentication Security

- Use secure authentication methods (OAuth2, API keys)
- Store credentials securely with encryption
- Rotate credentials regularly
- Use least privilege access

### API Security

- Validate all API responses
- Implement proper error handling
- Use HTTPS for all API communications
- Monitor API usage and rate limits

### Webhook Security

- Validate webhook signatures
- Use secure webhook endpoints
- Implement webhook authentication
- Monitor webhook delivery

## Troubleshooting

### Common Integration Issues

**Authentication Failed**
```bash
# Test authentication
spooky integrations test --integration github --auth-only

# Check credentials
spooky integrations config show --integration github
```

**API Rate Limiting**
```bash
# Check rate limit status
spooky integrations metrics --integration github --rate-limits

# Adjust rate limiting
spooky integrations config set --integration github --rate-limit 1000
```

**Webhook Delivery Issues**
```bash
# Test webhook endpoint
spooky integrations webhook --integration slack --test

# Check webhook configuration
spooky integrations config show --integration slack --webhooks
```

### Debugging Integrations

Enable debug logging for troubleshooting:

```bash
# Enable debug logging
spooky integrations test --integration github --debug

# View integration logs
spooky logging files tail --filter "integrations"
```

### Performance Optimization

Optimize integration performance:

```bash
# Check connection pooling
spooky integrations metrics --integration github --connections

# Adjust timeout settings
spooky integrations config set --integration github --timeout 60
```

## Integration with Other Systems

### Actions Integration

Integrations work with the actions system:

```hcl
actions {
  action "deploy-with-notifications" {
    description = "Deploy with integration notifications"
    
    machines = ["web-server"]
    parallel = true
    
    integrations {
      slack {
        event = "action_started"
        message = "Starting deployment on {{.machine}}"
      }
    }
    
    command = "deploy.sh"
  }
}
```

### Facts Integration

Integrations can collect facts from external sources:

```hcl
integrations {
  cloud_facts {
    enabled = true
    provider = "aws"
    
    facts {
      instance_metadata = true
      tags = true
      security_groups = true
    }
  }
}
```

### Variables Integration

Integrations can provide variables:

```hcl
integrations {
  config_store {
    enabled = true
    provider = "vault"
    
    variables {
      database_url = "/secret/database/url"
      api_key = "/secret/api/key"
    }
  }
}
```

## Examples

### GitHub Integration

```hcl
integrations {
  github {
    enabled = true
    api_url = "https://api.github.com"
    token = "ghp_..."
    
    repositories {
      "owner/repo" {
        events = ["push", "pull_request"]
        webhook_secret = "webhook-secret"
      }
    }
    
    notifications {
      events = ["deployment_started", "deployment_completed"]
      channel = "#deployments"
    }
  }
}
```

### Slack Integration

```hcl
integrations {
  slack {
    enabled = true
    webhook_url = "https://hooks.slack.com/services/..."
    channel = "#alerts"
    
    notifications {
      events = ["action_started", "action_completed", "action_failed"]
      
      templates {
        action_started = "🚀 Action {{.action_name}} started on {{.machine}}"
        action_completed = "✅ Action {{.action_name}} completed successfully"
        action_failed = "❌ Action {{.action_name}} failed: {{.error}}"
      }
    }
  }
}
```

### Custom API Integration

```hcl
integrations {
  monitoring_api {
    enabled = true
    base_url = "https://monitoring.example.com/api"
    
    authentication {
      method = "bearer"
      token = "monitoring-token"
    }
    
    endpoints {
      metrics = "/metrics"
      alerts = "/alerts"
      dashboards = "/dashboards"
    }
    
    operations {
      send_metric {
        endpoint = "/metrics"
        method = "POST"
        data_template = {
          name = "{{.metric_name}}"
          value = "{{.metric_value}}"
          timestamp = "{{.timestamp}}"
        }
      }
    }
  }
}
```

## Best Practices

### Integration Design

- Use consistent authentication patterns
- Implement proper error handling
- Use rate limiting and backoff strategies
- Monitor integration health

### Security

- Secure all credentials and tokens
- Validate all inputs and outputs
- Use HTTPS for all communications
- Implement proper access controls

### Performance

- Use connection pooling
- Implement caching where appropriate
- Monitor API rate limits
- Optimize request/response handling

### Monitoring

- Monitor integration health
- Track API usage and costs
- Monitor error rates and response times
- Set up alerts for integration failures

## Next Steps

- Explore the [Integrations API Reference](INTEGRATIONS_API_REFERENCE.md) for detailed technical information
- Check the [Integrations Troubleshooting Guide](INTEGRATIONS_TROUBLESHOOTING.md) for common issues
- Review the [Integrations Documentation Summary](INTEGRATIONS_DOCUMENTATION_SUMMARY.md) for implementation details
- Learn about [Custom Integration Development](INTEGRATIONS_USER_GUIDE.md) for advanced usage
