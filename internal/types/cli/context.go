package cli

// Context represents CLI execution context
type Context struct {
	ProjectPath string
	CommandName string
	Args        []string
	Flags       map[string]interface{}
}

// ExecutionContext represents command execution context
type ExecutionContext struct {
	Context     *Context
	Coordinator interface{} // Will be typed when coordinator interface is available
	Logger      interface{} // Will be typed when logger interface is available
}

// projectPathKey is used as a context key for project path
type projectPathKey struct{}
