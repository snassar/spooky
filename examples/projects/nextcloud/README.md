# Nextcloud Project

This project demonstrates how to organize spooky actions for Nextcloud installation and configuration. It includes both a **multi-file conceptual approach** (for future spooky features) and a **single-file practical approach** (for current spooky limitations).

## Project Structure

```
nextcloud/
├── project.hcl              # Project configuration
├── inventory.hcl            # Machine definitions (wrapper block format)
├── actions.hcl              # Single actions file (wrapper block format)
├── actions/                 # Multi-file structure (conceptual)
│   ├── 01-dependencies.hcl  # System dependencies
│   ├── 02-database.hcl      # Database configuration
│   ├── 03-web-server.hcl    # Apache and SSL configuration
│   ├── 04-installation.hcl  # Nextcloud download and install
│   ├── 05-optimization.hcl  # Performance optimization
│   └── 06-verification.hcl  # Service restart and verification
├── scripts/                 # Scripts referenced by actions
│   ├── configure-mariadb.sh
│   ├── configure-apache-nextcloud.sh
│   ├── install-nextcloud.sh
│   └── ...
└── README.md               # This file
```

## Configuration Format

This project uses the new wrapper block format for all HCL files:

### Inventory File (`inventory.hcl`)
```hcl
inventory {
  machine "example-server" {
    host     = "192.168.1.100"
    port     = 22
    user     = "debian"
    password = "your-password"
    tags = {
      environment = "development"
      role = "web"
    }
  }
}
```

### Actions File (`actions.hcl`)
```hcl
actions {
  action "check-status" {
    description = "Check server status"
    command     = "uptime && df -h"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
  }
}
```

## Current Spooky Limitation

**Important**: Currently, spooky only supports a single `actions.hcl` file per project. The multi-file structure in the `actions/` directory is **conceptual** and demonstrates how spooky could work in the future.

### Current Working Approach
- Use `actions.hcl` for all actions
- Organize with clear section comments
- Use tags for logical grouping

### Future Multi-File Approach
- Separate files for different phases
- Better maintainability and organization
- Requires spooky implementation changes

## Action Organization

### 01-dependencies.hcl
**Purpose**: Install and configure system dependencies
- `install-nextcloud-dependencies` - Install required packages
- `configure-php` - Configure PHP settings

### 02-database.hcl
**Purpose**: Database setup and configuration
- `configure-mariadb` - Configure MariaDB server
- `create-nextcloud-database` - Create database and user

### 03-web-server.hcl
**Purpose**: Web server configuration
- `configure-apache` - Configure Apache virtual host
- `setup-ssl-certificate` - Setup SSL certificates

### 04-installation.hcl
**Purpose**: Nextcloud installation
- `download-nextcloud` - Download and extract Nextcloud
- `install-nextcloud` - Run installation wizard

### 05-optimization.hcl
**Purpose**: Performance and maintenance
- `configure-nextcloud-optimization` - Performance optimization
- `setup-backup-cron` - Automated backup setup

### 06-verification.hcl
**Purpose**: Final steps and verification
- `restart-services` - Restart all services
- `verify-nextcloud-installation` - Verify installation

## Prerequisites

- Target machines must be running Ubuntu/Debian
- Machines must have the appropriate tags:
  - `tag:role=web-server` for web servers
  - `tag:role=database` for database servers
- Internet connectivity for downloading Nextcloud and packages
- SSH access with sudo privileges

## Installation Sequence

The actions should be executed in the following order:

```bash
# Phase 1: Dependencies
spooky action run --project ./nextcloud install-nextcloud-dependencies
spooky action run --project ./nextcloud configure-php

# Phase 2: Database
spooky action run --project ./nextcloud configure-mariadb
spooky action run --project ./nextcloud create-nextcloud-database

# Phase 3: Web Server
spooky action run --project ./nextcloud configure-apache
spooky action run --project ./nextcloud setup-ssl-certificate

# Phase 4: Installation
spooky action run --project ./nextcloud download-nextcloud
spooky action run --project ./nextcloud install-nextcloud

# Phase 5: Optimization
spooky action run --project ./nextcloud configure-nextcloud-optimization
spooky action run --project ./nextcloud setup-backup-cron

# Phase 6: Verification
spooky action run --project ./nextcloud restart-services
spooky action run --project ./nextcloud verify-nextcloud-installation
```

## Usage Examples

### Execute All Actions
```bash
spooky action run --project ./nextcloud --tags nextcloud
```

### Execute Specific Phase
```bash
# Only dependencies
spooky action run --project ./nextcloud --tags dependencies

# Only database setup
spooky action run --project ./nextcloud --tags database

# Only web server configuration
spooky action run --project ./nextcloud --tags apache
```

### Execute Actions by Name
```bash
# Execute all actions from dependencies phase
spooky action run --project ./nextcloud install-nextcloud-dependencies configure-php

# Execute all actions from database phase
spooky action run --project ./nextcloud configure-mariadb create-nextcloud-database
```

### Conditional Execution
```bash
# Skip SSL setup (optional)
spooky action run --project ./nextcloud --exclude-tags ssl

# Only run verification actions
spooky action run --project ./nextcloud --tags verification
```

### Target Specific Machines
```bash
# Execute on web servers only
spooky action run --project ./nextcloud --machines "tag:role=web-server" install-nextcloud-dependencies

# Execute on database servers only
spooky action run --project ./nextcloud --machines "tag:role=database" configure-mariadb
```

## Configuration

### Database Configuration

The default database configuration uses:
- Database name: `nextcloud`
- Database user: `nextcloud`
- Database password: `nextcloud_secure_password`

You can modify these in the `scripts/configure-mariadb.sh` file.

### Admin User

The default admin user is:
- Username: `admin`
- Password: `admin_secure_password`

You can modify these in the `scripts/install-nextcloud.sh`