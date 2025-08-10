# Script Actions vs Files Directory: Why Temporary Files?

Based on analysis of the codebase, here's the key distinction and why temporary files are necessary:

## Two Different Concepts

### 1. Script Actions (`type = "script"`)
- **Purpose**: Execute script **content** defined inline in the action
- **Content Source**: The `script` field contains the actual script content as a string
- **Example**:
```hcl
action "setup_nginx" {
  type = "script"
  script = """
    #!/bin/bash
    apt-get update
    apt-get install -y nginx
    systemctl enable nginx
    systemctl start nginx
  """
}
```

### 2. Files Directory (`<project>/files/`)
- **Purpose**: Store static files for deployment (not execution)
- **Content Source**: Pre-existing files in the project's `files/` directory
- **Use Case**: Configuration files, static assets, documentation, etc.

## Why Temporary Files Are Needed

### The Core Issue: Script Content vs File Path

**Script Actions contain script content, not file paths:**

```go
// From internal/actions/types/action.go
type Action struct {
    Script string `hcl:"script,optional"`  // Contains script content, not file path
}
```

**The validation confirms this:**
```go
// From internal/actions/acting/actor.go
if a.action.Script == "" {
    return fmt.Errorf("script content is required")  // Not "script file path"
}
```

### Execution Requirements

**To execute a script, you need:**
1. **A file on disk** (executable programs require file paths)
2. **Executable permissions** (chmod 0755)
3. **Proper file extension** (e.g., `.sh` for shell scripts)

**Script content in a string cannot be executed directly:**
- `exec.Command("script content")` - ❌ Invalid
- `exec.Command("/path/to/script.sh")` - ✅ Valid

## The Solution: Temporary File Creation

### Process Flow:
1. **Extract script content** from the action's `script` field
2. **Create temporary file** with unique name (`spooky-script-*.sh`)
3. **Write script content** to the temporary file
4. **Set executable permissions** (`chmod 0755`)
5. **Execute the file** using `exec.CommandContext(ctx, scriptFile)`
6. **Clean up** temporary file with `defer os.Remove(scriptFile)`

### Code Implementation:
```go
func (a *actorImpl) createTemporaryScript(scriptContent string) (string, error) {
    // Create temporary file
    tmpFile, err := os.CreateTemp("", "spooky-script-*.sh")
    if err != nil {
        return "", fmt.Errorf("failed to create temporary file: %w", err)
    }

    // Write script content
    if _, err := tmpFile.WriteString(scriptContent); err != nil {
        tmpFile.Close()
        os.Remove(tmpFile.Name())
        return "", fmt.Errorf("failed to write script content: %w", err)
    }

    // Make script executable
    if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
        os.Remove(tmpFile.Name())
        return "", fmt.Errorf("failed to make script executable: %w", err)
    }

    return tmpFile.Name(), nil
}
```

## Alternative Approaches (Not Currently Implemented)

### Option 1: Script File References
```hcl
action "setup_nginx" {
  type = "script"
  script_file = "files/setup_nginx.sh"  # Reference to file in files/ directory
}
```

### Option 2: Script Directory
```hcl
action "setup_nginx" {
  type = "script"
  script_file = "scripts/setup_nginx.sh"  # Reference to file in scripts/ directory
}
```

### Option 3: Hybrid Approach
```hcl
action "setup_nginx" {
  type = "script"
  script = "files/setup_nginx.sh"  # Could be either content or file path
}
```

## Current Design Rationale

### Advantages of Current Approach:
- **Self-contained actions**: Script content is embedded in the action definition
- **Version control friendly**: Script changes are tracked with action changes
- **No file dependencies**: Actions don't depend on external files
- **Atomic operations**: Everything needed is in the action definition

### Disadvantages:
- **Temporary file overhead**: Creates and cleans up files for each execution
- **No script reuse**: Can't reference the same script from multiple actions
- **Limited script size**: Large scripts make action definitions unwieldy

## Future Considerations

The current implementation could be enhanced to support:
- **Script file references** alongside inline content
- **Script templates** with variable substitution
- **Script libraries** for reusable components
- **Caching** of frequently used scripts

But for now, temporary files are the necessary bridge between inline script content and executable file requirements.
