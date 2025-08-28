# Quick Reference: Using JSON Prompts in Cursor

## 🚀 Quick Start

1. **Copy a prompt**: `cat .cursor/prompts/bug-fix.json | pbcopy`
2. **Paste in Cursor**: Paste the JSON into Cursor's chat
3. **Add context**: "I'm working on a Go SSH application. Please use this structured approach to find and fix bugs."
4. **Follow the workflow**: The AI will follow the structured steps in the prompt

## 📋 Available Prompts

| Task | File | When to Use |
|------|------|-------------|
| **Bug Fixing** | `bug-fix.json` | When you suspect bugs or want proactive issue detection |
| **TODO Completion** | `todo-finder.json` | For code cleanup and completing placeholder implementations |
| **Code Review** | `code-review.json` | Before merging code or for quality audits |
| **Refactoring** | `refactoring.json` | To improve code quality and reduce complexity |
| **Feature Implementation** | `feature-implementation.json` | When adding new functionality |
| **Feature Audit** | `feature-audit.json` | To identify unnecessary features and complexity |

## 🎯 Common Use Cases

### Finding and Fixing Bugs
```bash
cat .cursor/prompts/bug-fix.json | pbcopy
```
Then in Cursor: "Please analyze this Go SSH application for bugs and fix them using the structured approach."

### Completing TODO Items
```bash
cat .cursor/prompts/todo-finder.json | pbcopy
```
Then in Cursor: "Find all TODO items and incomplete code in this project and complete them properly."

### Code Review
```bash
cat .cursor/prompts/code-review.json | pbcopy
```
Then in Cursor: "Conduct a comprehensive code review of this SSH application focusing on security and quality."

### Refactoring Code
```bash
cat .cursor/prompts/refactoring.json | pbcopy
```
Then in Cursor: "Identify refactoring opportunities in this codebase and implement improvements."

### Adding New Features
```bash
cat .cursor/prompts/feature-implementation.json | pbcopy
```
Then in Cursor: "I want to add a new SSH key management feature. Please implement it following the structured approach."

### Auditing Unnecessary Features
```bash
cat .cursor/prompts/feature-audit.json | pbcopy
```
Then in Cursor: "Please audit this Spooky automation tool to identify unnecessary features, over-engineering, and complexity that doesn't serve the core use case."

## 🔧 Customization Tips

### For Different Languages
- Replace Go-specific patterns with your language's conventions
- Update `target_directories` to match your project structure
- Modify error handling patterns to match your language

### For Different Project Types
- Update `project_type` in the context section
- Adjust security considerations for your domain
- Modify `codebase_patterns` to match your project

### For Specific Tasks
- Add new tasks to the `tasks` section
- Define custom output formats
- Update constraints for your requirements

## 📊 Expected Outputs

Each prompt provides structured output:

- **Bug Fix**: Bug reports with severity, fixes with validation
- **TODO Finder**: Incomplete items report, completion plan, implementation results
- **Code Review**: Review summary, issues by priority, recommendations
- **Refactoring**: Analysis, plan, results with metrics
- **Feature Implementation**: Analysis, plan, results with verification steps
- **Feature Audit**: Core purpose analysis, complexity issues, simplification plan

## ⚡ Pro Tips

1. **Be Specific**: Always provide context about your project
2. **Use Constraints**: Set clear boundaries to prevent unwanted changes
3. **Follow Up**: Ask for clarification if the AI doesn't follow the structure
4. **Iterate**: Refine prompts based on results
5. **Combine**: Use multiple prompts for complex tasks

## 🛠️ Troubleshooting

### AI Ignores the Structure
- Make sure the JSON is properly formatted
- Explicitly ask the AI to follow the structured approach
- Provide more specific context about your project

### Results Are Too Generic
- Add more specific constraints
- Include examples of what you want
- Specify the exact files or areas to focus on

### Missing Context
- Always describe your project type and goals
- Mention any specific requirements or constraints
- Include relevant file paths or directories

## 📚 Advanced Usage

### Combining Prompts
You can combine multiple prompts for complex tasks:
1. Use code review to identify issues
2. Use bug fix to address critical problems
3. Use refactoring to improve overall quality

### Custom Prompts
Create your own prompts by:
1. Copying an existing prompt
2. Modifying the context and tasks
3. Testing on a small codebase first
4. Iterating based on results

### Team Usage
- Share prompts with your team
- Standardize on common prompts
- Create project-specific variations
- Document any customizations
