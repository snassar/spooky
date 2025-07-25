# hcl_collector Project

This is a spooky project with separated inventory and actions configuration.

## Project Structure

hcl_collector/
├── project.hcl          # Project configuration and settings
├── inventory.hcl        # Machine definitions
├── actions.hcl          # Action definitions for automation
├── .gitignore          # Git ignore rules
├── templates/           # Template files for dynamic content
├── files/              # Static files to be deployed
└── README.md           # This file

## Usage

### Execute an action:
spooky execute actions.hcl --inventory inventory.hcl --action check-status

### Execute with tag targeting:
spooky execute actions.hcl --inventory inventory.hcl --action update-system --tags "role=web"

## Configuration

- **project.hcl**: Project settings, storage, logging, and SSH configuration
- **inventory.hcl**: Machine definitions with tags for targeting
- **actions.hcl**: Automation actions with tag-based targeting

## Benefits

1. **Reusability**: Actions can be applied to different inventories
2. **Maintainability**: Clear separation of concerns
3. **Flexibility**: Mix and match actions with different machine groups
4. **Version Control**: Better tracking of changes to machines vs actions
