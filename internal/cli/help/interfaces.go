package help

// HelpManager defines the interface for help rendering
type HelpManager interface {
	ShowHelp(commandName string) (string, error)
	ShowUsage(commandName string) (string, error)
	ShowExamples(commandName string) (string, error)
}

// HelpRenderer defines the interface for help rendering
type HelpRenderer interface {
	RenderHelp(commandName string) (string, error)
	RenderUsage(commandName string) (string, error)
	RenderExamples(commandName string) (string, error)
}
