# Quick Reference: Multi-File Actions

## File Structure
```
actions/
├── 01-dependencies.hcl    # System packages, PHP config
├── 02-database.hcl        # MariaDB setup
├── 03-web-server.hcl      # Apache, SSL
├── 04-installation.hcl    # Download, install Nextcloud
├── 05-optimization.hcl    # Performance, backups
├── 06-verification.hcl    # Services, health checks
└── nextcloud-installation.hcl  # Single file reference
```

## Common Commands

### Full Deployment
```bash
spooky action run --project ./nextcloud --tags nextcloud
```

### Phase-by-Phase
```bash
# Dependencies
spooky action run --project ./nextcloud --tags dependencies

# Database
spooky action run --project ./nextcloud --tags database

# Web Server
spooky action run --project ./nextcloud --tags apache

# Installation
spooky action run --project ./nextcloud --tags installation

# Optimization
spooky action run --project ./nextcloud --tags optimization

# Verification
spooky action run --project ./nextcloud --tags verification
```

### Selective Execution
```bash
# Skip SSL (development)
spooky action run --project ./nextcloud --exclude-tags ssl --tags nextcloud

# Only specific actions
spooky action run --project ./nextcloud install-nextcloud-dependencies configure-php

# Target specific machines
spooky action run --project ./nextcloud --machines "tag:role=web-server" --tags dependencies
```

## Action Tags

| Tag | Actions | Purpose |
|-----|---------|---------|
| `dependencies` | install-nextcloud-dependencies, configure-php | System setup |
| `database` | configure-mariadb, create-nextcloud-database | Database setup |
| `apache` | configure-apache, setup-ssl-certificate | Web server |
| `installation` | download-nextcloud, install-nextcloud | App installation |
| `optimization` | configure-nextcloud-optimization, setup-backup-cron | Performance |
| `verification` | restart-services, verify-nextcloud-installation | Health checks |
| `ssl` | setup-ssl-certificate | SSL configuration |
| `nextcloud` | All actions | Complete deployment |

## Troubleshooting

### Check Available Actions
```bash
spooky action list --project ./nextcloud
```

### Debug Action
```bash
spooky action run --project ./nextcloud ACTION_NAME --verbose
```

### Check Project Status
```bash
spooky project validate ./nextcloud
```

## Best Practices

1. **Execute in order**: 01 → 02 → 03 → 04 → 05 → 06
2. **Use tags**: Target specific phases with tags
3. **Test phases**: Run phases individually for debugging
4. **Monitor logs**: Check logs for detailed error information
5. **Rollback**: Keep backup of working configurations

## Migration Tips

### From Single File
1. Copy actions from single file to numbered files
2. Group related actions together
3. Update any cross-references
4. Test each file independently

### Adding New Actions
1. Create new numbered file (e.g., `07-monitoring.hcl`)
2. Add appropriate tags
3. Update documentation
4. Test thoroughly

## Environment Variations

### Development
```bash
spooky action run --project ./nextcloud --exclude-tags ssl --tags nextcloud
```

### Production
```bash
spooky action run --project ./nextcloud --tags nextcloud
```

### Database Only
```bash
spooky action run --project ./nextcloud --tags database
```

### Web Server Only
```bash
spooky action run --project ./nextcloud --tags apache
``` 