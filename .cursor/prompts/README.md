# Cursor JSON Prompts for Software Development

This directory contains structured JSON prompts for Cursor AI to handle common software development tasks. These prompts follow the principles outlined in the article about JSON prompts being the future of AI automation.

## How to Use These Prompts

1. **Copy the JSON content** from any prompt file
2. **Paste it into Cursor's chat** when starting a conversation
3. **Modify the context** as needed for your specific project
4. **Follow the structured workflow** outlined in the prompt

## Available Prompts

### 1. Bug Fix Assistant (`bug-fix.json`)
- **Purpose**: Identify and fix bugs in Go code
- **Best for**: When you suspect there are bugs or want to proactively find issues
- **Key features**:
  - Systematic bug identification
  - Structured fix implementation
  - Validation and testing requirements

### 2. TODO Finder (`todo-finder.json`)
- **Purpose**: Find and complete TODO items and incomplete code
- **Best for**: Code cleanup, completing placeholder implementations
- **Key features**:
  - Comprehensive TODO scanning
  - Priority-based completion planning
  - Ensures no placeholder code remains

### 3. Code Review Assistant (`code-review.json`)
- **Purpose**: Conduct thorough code reviews
- **Best for**: Before merging code, during refactoring, or for quality audits
- **Key features**:
  - Multi-dimensional review (security, quality, performance, tests)
  - Structured output with actionable recommendations
  - Positive feedback highlighting good practices

### 4. Code Refactoring Assistant (`refactoring.json`)
- **Purpose**: Refactor code to improve quality and maintainability
- **Best for**: Code cleanup, reducing complexity, improving architecture
- **Key features**:
  - Identifies refactoring opportunities
  - Plans safe refactoring strategies
  - Ensures no functionality is lost
  - Updates tests and documentation

### 5. Feature Implementation Assistant (`feature-implementation.json`)
- **Purpose**: Implement new features with proper architecture
- **Best for**: Adding new functionality, extending existing features
- **Key features**:
  - Follows existing code patterns
  - Includes comprehensive testing
  - Updates documentation
  - Maintains backward compatibility

### 6. Feature Audit Assistant (`feature-audit.json`)
- **Purpose**: Identify unnecessary or over-engineered features
- **Best for**: Code simplification, reducing complexity, removing bloat
- **Key features**:
  - Audits development tools for necessity
  - Identifies over-engineering and complexity
  - Finds unused or rarely-used features
  - Evaluates external dependencies
  - Provides simplification recommendations

## JSON Prompt Structure

Each prompt follows this structure:

```json
{
  "name": "Prompt Name",
  "description": "What this prompt does",
  "version": "1.0.0",
  "context": {
    "role": "AI's role and expertise",
    "project_type": "Type of project",
    "standards": ["List of standards to follow"]
  },
  "tasks": {
    "task_name": {
      "description": "What this task does",
      "steps": [
        {
          "action": "tool_to_use",
          "parameters": "tool parameters"
        }
      ]
    }
  },
  "output_format": {
    "structured_output": "How results should be formatted"
  },
  "constraints": {
    "rule": "constraint description"
  }
}
```

## Customizing Prompts

### For Different Languages
- Replace Go-specific patterns with your language's patterns
- Update tool queries to match your language's conventions
- Modify error handling patterns to match your language

### For Different Project Types
- Update `project_type` and `codebase_patterns`
- Modify `target_directories` in search actions
- Adjust security considerations for your domain

### For Different Tasks
- Add new tasks to the `tasks` section
- Define new output formats as needed
- Update constraints for your specific requirements

## Best Practices

1. **Be Specific**: Include detailed context about your project
2. **Use Structured Output**: Define clear output formats for consistent results
3. **Include Constraints**: Set boundaries to prevent unwanted changes
4. **Test Prompts**: Try prompts on small codebases first
5. **Iterate**: Refine prompts based on results

## Example Usage

```bash
# Copy a prompt and use it in Cursor
cat .cursor/prompts/bug-fix.json | pbcopy  # On macOS
cat .cursor/prompts/bug-fix.json | xclip -selection clipboard  # On Linux
```

Then paste into Cursor's chat and add your specific context:

```
I'm working on a Go SSH application. Please use this structured approach to find and fix any bugs in the codebase.
```

## Contributing

To add new prompts:
1. Create a new JSON file in this directory
2. Follow the established structure
3. Update this README with the new prompt's description
4. Test the prompt on a sample codebase

## Tips for Effective Prompts

- **Start with context**: Always provide clear context about your project
- **Be explicit about constraints**: Specify what should NOT be changed
- **Define clear outputs**: Structure how you want results presented
- **Include examples**: Provide examples of good and bad patterns
- **Set priorities**: Help the AI understand what's most important
