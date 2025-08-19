# Issue Management Workflow

This directory contains markdown files for Codeberg issues that can be created using the Codeberg API. Issue management commands are provided by the local `justfile`.

## Quick Start

```bash
# Navigate to issues directory
cd docs/issues

# Set up environment variables
export CODEBERG_TOKEN="your-api-token"
export CODEBERG_OWNER="your-username"
export CODEBERG_REPO="your-repo-name"

# List available commands
just help

# Create new issue file from template
just new my-new-issue

# Edit the issue file
just edit issue-my-new-issue.md

# Create Codeberg issue from file (only if it doesn't exist)
just create issue-my-new-issue.md

# List local and Codeberg issues
just list
```

## Workflow

### 1. Create Issue File

Create a new markdown file using the justfile:

```bash
cd docs/issues
just new my-new-issue.md
```

### 2. Edit the Issue

Edit the markdown file with your issue details. The first line should be the issue title:

```markdown
# My Issue Title

## Problem Statement
...
```

### 3. Create the Issue

Use the justfile commands:

```bash
# Create issue
just create my-new-issue.md

# Create as draft
just create-draft my-new-issue.md
```

### 4. Add Labels (Optional)

After creation, add labels to the issue using the Codeberg web interface or API:

```bash
# You can add labels through the Codeberg web interface
# or use the API directly with curl
```

## Available Commands

```bash
# List both local and Codeberg issues
just list

# List local issue files only
just list-local

# List Codeberg issues only
just list-codeberg

# Create new issue file from template (issue-<filename>.md)
just new <filename>

# Create Codeberg issue from file (only if not exists)
just create <file>

# Create Codeberg issue (draft mode not supported)
just create-draft <file>

# Validate issue file format
just validate <file>

# Show issue file content
just show <file>

# Edit issue file
just edit <file>

# Show help
just help
```

## Best Practices

### 1. **Use Descriptive Titles**
- Start with a verb (Add, Fix, Improve, etc.)
- Be specific about what's being changed
- Keep it concise but informative

### 2. **Follow the Template Structure**
- Problem Statement: Clear description of the issue
- Use Cases: Who needs this and why
- Proposed Solution: Detailed implementation plan
- Examples: Practical usage examples
- Acceptance Criteria: Clear checklist for completion

### 3. **Include Code Examples**
- Show the proposed API or interface
- Include usage examples
- Reference existing code when relevant

### 4. **Add Appropriate Labels**
- `enhancement` for new features
- `bug` for bug fixes
- `documentation` for docs improvements
- `infrastructure` for tooling/build changes

### 5. **Set Priority**
- **High**: Blocks development or critical functionality
- **Medium**: Important but not blocking
- **Low**: Nice to have or future enhancement

## File Naming Convention

Use kebab-case for issue files:

```
docs/issues/
├── 00-issue-template.md                    # Template (don't create issues from this)
├── hardcoded-actions-hcl-generation.md     # Schema-driven actions generation
├── hardcoded-spooky-config-generation.md   # Schema-driven config generation
├── hardcoded-logging-config-generation.md  # Schema-driven logging generation
├── hardcoded-variables-hcl-generation.md   # Schema-driven variables generation
├── hardcoded-machines-hcl-generation.md    # Schema-driven machines generation
├── cli-logging-flags.md                    # CLI logging enhancement
└── README.md                               # This file
```

## Benefits of This Approach

1. **Version Control**: Issues are tracked in git
2. **Better Quality**: More time to think through problems
3. **Collaboration**: Easy to review and discuss before creation
4. **Templates**: Consistent issue structure
5. **Reusability**: Can copy and modify existing issues

## Examples

See `cli-logging-flags.md` for a complete example of a well-structured issue.
