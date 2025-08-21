# Test Automation Project

This is a test automation project for validating spooky schema functionality.

## Project Structure

```
.
├── project.hcl          # Project configuration
├── machines.hcl         # Machine inventory
├── actions.hcl          # Action definitions
├── variables.hcl        # Variable definitions
├── templates/           # Template files
│   ├── nginx.conf.tmpl  # Nginx configuration
│   ├── app.conf.tmpl    # Application configuration
│   ├── app.service.tmpl # Systemd service
│   ├── docker-compose.yml.tmpl # Docker Compose
│   └── env.tmpl         # Environment variables
├── files/               # Static files
│   ├── deploy.sh        # Deployment script
│   └── backup.sh        # Backup script
└── logs/                # Log files directory
```

## Usage

This project demonstrates a complete spooky automation setup with:

- **Project Configuration**: Basic project settings and metadata
- **Machine Inventory**: Web and database servers
- **Actions**: Deployment, backup, and configuration management
- **Variables**: Environment-specific configuration
- **Templates**: Dynamic configuration generation
- **Static Files**: Deployment and maintenance scripts

## Validation

This project should pass all schema validations:

- ✅ Project structure validation
- ✅ Machine inventory validation
- ✅ Action definitions validation
- ✅ Variable definitions validation
- ✅ Template syntax validation
- ✅ Cross-file reference validation

## Testing

Use this project to test:

1. Schema validation functionality
2. Project directory structure validation
3. Cross-file relationship validation
4. Template rendering and validation
5. Action execution workflows
