package secrets

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// PasswordPrompter interface for prompting for passwords
type PasswordPrompter interface {
	PromptForPassword(prompt string) (string, error)
}

// TerminalPasswordPrompter implements password prompting using terminal
type TerminalPasswordPrompter struct{}

// NewTerminalPasswordPrompter creates a new terminal password prompter
func NewTerminalPasswordPrompter() *TerminalPasswordPrompter {
	return &TerminalPasswordPrompter{}
}

// PromptForPassword prompts for a password securely using the terminal
func (t *TerminalPasswordPrompter) PromptForPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	defer fmt.Println() // Add newline after password input

	// Get the file descriptor for stdin
	fd := int(os.Stdin.Fd())

	// Check if stdin is a terminal
	if !term.IsTerminal(fd) {
		return "", fmt.Errorf("stdin is not a terminal")
	}

	// Read password without echo
	password, err := term.ReadPassword(fd)
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}

	return string(password), nil
}

// PromptForAgePassword prompts for the age identity password
func PromptForAgePassword(identityPath string) (string, error) {
	prompter := NewTerminalPasswordPrompter()
	prompt := fmt.Sprintf("Enter password for age identity %s: ", identityPath)
	return prompter.PromptForPassword(prompt)
}
