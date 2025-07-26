# test-valid-project Project

This is a spooky project with separated inventory and actions configuration.

## Project Structure

test-valid-project/
├── project.hcl          # Project configuration and settings
├── inventory.hcl        # Machine definitions
├── actions.hcl          # Main action definitions (optional)
├── actions/             # Directory for organized action files
│   ├── 01-dependencies.hcl
│   ├── 02-system-update.hcl
│   └── 03-monitoring.hcl
├── templates/           # Template files for dynamic content
├── data/               # Data files for templates
├── files/              # Static files to be deployed
├── logs/               # Log files
├── .gitignore          # Git ignore rules
└── README.md           # This file

## Actions Organization

Spooky supports flexible action organization:

1. **actions.hcl** - Main actions file (optional)
2. **actions/** directory - Organized action files
   - Files are loaded in alphabetical order
   - Use numbered prefixes (01-, 02-, etc.) for ordering
   - Each file can contain multiple actions

## Usage

### List all actions:
spooky list-actions

### Execute an action:
spooky execute actions.hcl --inventory inventory.hcl --action check-status

### Execute with tag targeting:
spooky execute actions.hcl --inventory inventory.hcl --action update-system --tags "role=web"

## Configuration

- **project.hcl**: Project settings, storage, logging, and SSH configuration
- **inventory.hcl**: Machine definitions with tags for targeting
- **actions.hcl**: Main actions (optional)
- **actions/**: Organized action files for better maintainability

## Benefits

1. **Reusability**: Actions can be applied to different inventories
2. **Maintainability**: Clear separation of concerns
3. **Flexibility**: Mix and match actions with different machine groups
4. **Organization**: Group related actions in separate files
5. **Version Control**: Better tracking of changes to machines vs actions
