# Issue Management Workflow

This directory contains markdown files for GitHub issues that can be created using the `gh` CLI tool.

## Workflow

### 1. Create Issue File

Create a new markdown file in this directory following the template:

```bash
cp docs/issues/00-issue-template.md docs/issues/my-new-issue.md
```

### 2. Edit the Issue

Edit the markdown file with your issue details. The first line should be the issue title:

```markdown
# My Issue Title

## Problem Statement
...
```

### 3. Create the Issue

Use the just command:

```bash
# Create issue
just create-issue docs/issues/my-new-issue.md

# Create as draft
just create-issue-draft docs/issues/my-new-issue.md
```

### 4. Add Labels (Optional)

After creation, add labels to the issue:

```bash
gh issue edit <number> --add-label enhancement,infrastructure
```

## Available Commands

```bash
# List available issue files
just list-issues

# Create issue from file
just create-issue docs/issues/my-issue.md

# Create issue as draft
just create-issue-draft docs/issues/my-issue.md
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
├── 00-issue-template.md          # Template (don't create issues from this)
├── cli-logging-flags.md          # CLI logging enhancement
├── fix-ssh-connection-bug.md     # Bug fix
├── improve-documentation.md      # Documentation improvement
└── README.md                     # This file
```

## Benefits of This Approach

1. **Version Control**: Issues are tracked in git
2. **Better Quality**: More time to think through problems
3. **Collaboration**: Easy to review and discuss before creation
4. **Templates**: Consistent issue structure
5. **Reusability**: Can copy and modify existing issues

## Examples

See `cli-logging-flags.md` for a complete example of a well-structured issue.
