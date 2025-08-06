package completion

// CompletionManager defines the interface for completion generation
type CompletionManager interface {
	GenerateCompletion(shell string) (string, error)
	GenerateCompletionFile(shell, outputPath string) error
	GetSupportedShells() []string
}

// CompletionGenerator defines the interface for completion generation
type CompletionGenerator interface {
	GenerateBashCompletion() (string, error)
	GenerateZshCompletion() (string, error)
	GenerateFishCompletion() (string, error)
	GeneratePowerShellCompletion() (string, error)
}
