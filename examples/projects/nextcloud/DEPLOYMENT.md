# Nextcloud Deployment Guide

This guide shows how to deploy Nextcloud using the multi-file action approach.

## Prerequisites

1. **Target Machines**: Ensure your inventory has machines tagged with:
   - `tag:role=web-server` for web servers
   - `tag:role=database` for database servers

2. **Network Access**: Machines need internet access to download packages and Nextcloud

3. **Permissions**: SSH access with sudo privileges

## Quick Deployment

### Option 1: Full Installation (All Phases)

```bash
# Execute all Nextcloud actions in sequence
spooky action run --project ./nextcloud --tags nextcloud
```

### Option 2: Phase-by-Phase Deployment

```bash
# Phase 1: Dependencies
echo "Installing dependencies..."
spooky action run --project ./nextcloud install-nextcloud-dependencies
spooky action run --project ./nextcloud configure-php

# Phase 2: Database
echo "Configuring database..."
spooky action run --project ./nextcloud configure-mariadb
spooky action run --project ./nextcloud create-nextcloud-database

# Phase 3: Web Server
echo "Configuring web server..."
spooky action run --project ./nextcloud configure-apache
spooky action run --project ./nextcloud setup-ssl-certificate

# Phase 4: Installation
echo "Installing Nextcloud..."
spooky action run --project ./nextcloud download-nextcloud
spooky action run --project ./nextcloud install-nextcloud

# Phase 5: Optimization
echo "Optimizing configuration..."
spooky action run --project ./nextcloud configure-nextcloud-optimization
spooky action run --project ./nextcloud setup-backup-cron

# Phase 6: Verification
echo "Verifying installation..."
spooky action run --project ./nextcloud restart-services
spooky action run --project ./nextcloud verify-nextcloud-installation
```

## Selective Deployment Scenarios

### Scenario 1: Development Environment (No SSL)

```bash
# Skip SSL setup for development
spooky action run --project ./nextcloud --exclude-tags ssl --tags nextcloud
```

### Scenario 2: Database-Only Setup

```bash
# Only configure database (for separate database server)
spooky action run --project ./nextcloud --tags database
```

### Scenario 3: Web Server Only

```bash
# Only configure web server (assumes database is already configured)
spooky action run --project ./nextcloud --tags apache
```

### Scenario 4: Update Dependencies Only

```bash
# Only update system dependencies
spooky action run --project ./nextcloud --tags dependencies
```

## Troubleshooting

### Check Action Status

```bash
# List all available actions
spooky action list --project ./nextcloud

# Check specific action details
spooky action show --project ./nextcloud install-nextcloud-dependencies
```

### Debug Failed Actions

```bash
# Run with verbose output
spooky action run --project ./nextcloud install-nextcloud-dependencies --verbose

# Check logs
spooky action logs --project ./nextcloud
```

### Common Issues

1. **Package Installation Fails**
   ```bash
   # Check package availability
   spooky action run --project ./nextcloud --machines "tag:role=web-server" "apt update"
   ```

2. **Database Connection Issues**
   ```bash
   # Verify database service
   spooky action run --project ./nextcloud --machines "tag:role=database" "systemctl status mariadb"
   ```

3. **Permission Issues**
   ```bash
   # Check file permissions
   spooky action run --project ./nextcloud --machines "tag:role=web-server" "ls -la /var/www/nextcloud"
   ```

## Rollback Procedures

### Rollback to Previous State

```bash
# If installation fails, you can rollback by removing Nextcloud
spooky action run --project ./nextcloud --machines "tag:role=web-server" "rm -rf /var/www/nextcloud"
```

### Database Rollback

```bash
# Drop and recreate database if needed
spooky action run --project ./nextcloud --machines "tag:role=database" "mysql -e 'DROP DATABASE IF EXISTS nextcloud;'"
```

## Monitoring and Maintenance

### Health Checks

```bash
# Regular health check
spooky action run --project ./nextcloud verify-nextcloud-installation

# Check service status
spooky action run --project ./nextcloud --machines "tag:role=web-server" "systemctl status apache2"
spooky action run --project ./nextcloud --machines "tag:role=database" "systemctl status mariadb"
```

### Backup Verification

```bash
# Check backup cron job
spooky action run --project ./nextcloud --machines "tag:role=web-server" "crontab -l | grep nextcloud"
```

## Customization

### Modify Configuration

1. **Edit action files** in the `actions/` directory
2. **Update scripts** in the `scripts/` directory
3. **Modify inventory** to target different machines

### Add Custom Actions

```bash
# Create new action file
touch actions/07-custom.hcl

# Add custom actions
cat > actions/07-custom.hcl << 'EOF'
action "custom-action" {
  description = "Custom action for your needs"
  command = "echo 'Custom action executed'"
  tags = ["custom"]
  machines = ["tag:role=web-server"]
  parallel = true
  timeout = 60
}
EOF
```

## Performance Optimization

### Large Deployments

For large deployments with many machines:

```bash
# Run in parallel where possible
spooky action run --project ./nextcloud --tags dependencies --parallel

# Use batch processing
spooky action run --project ./nextcloud --tags nextcloud --batch-size 10
```

### Resource Management

```bash
# Monitor resource usage during deployment
spooky action run --project ./nextcloud --machines "tag:role=web-server" "top -n 1"
```

## Security Considerations

### Production Deployment

1. **Change default passwords** in scripts
2. **Enable SSL** (use `setup-ssl-certificate` action)
3. **Configure firewall** rules
4. **Set up monitoring** and alerting

### Access Control

```bash
# Restrict access to specific IPs
spooky action run --project ./nextcloud --machines "tag:role=web-server" "ufw allow from 192.168.1.0/24 to any port 80,443"
```

## Integration with CI/CD

### Automated Deployment

```bash
#!/bin/bash
# deploy-nextcloud.sh

set -e

PROJECT_DIR="./nextcloud"
ENVIRONMENT="${1:-production}"

echo "Deploying Nextcloud to $ENVIRONMENT..."

# Validate project
spooky project validate "$PROJECT_DIR"

# Deploy based on environment
if [ "$ENVIRONMENT" = "production" ]; then
    spooky action run --project "$PROJECT_DIR" --tags nextcloud
else
    spooky action run --project "$PROJECT_DIR" --exclude-tags ssl --tags nextcloud
fi

# Verify deployment
spooky action run --project "$PROJECT_DIR" verify-nextcloud-installation

echo "Nextcloud deployment completed successfully!"
```

This multi-file approach provides the flexibility to deploy Nextcloud in various scenarios while maintaining clear organization and easy troubleshooting. 