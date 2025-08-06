package outputs

import (
	"os"
)

// ConsoleOutput implements console output for log entries
type ConsoleOutput struct {
	writer *os.File
}

// NewConsoleOutput creates a new console output
func NewConsoleOutput() *ConsoleOutput {
	return &ConsoleOutput{
		writer: os.Stdout,
	}
}

// NewStderrOutput creates a new stderr output
func NewStderrOutput() *ConsoleOutput {
	return &ConsoleOutput{
		writer: os.Stderr,
	}
}

// Write writes data to the console
func (o *ConsoleOutput) Write(data []byte) error {
	_, err := o.writer.Write(data)
	return err
}

// Close closes the console output
func (o *ConsoleOutput) Close() error {
	// Console output doesn't need to be closed
	return nil
}

// GetName returns the output name
func (o *ConsoleOutput) GetName() string {
	if o.writer == os.Stderr {
		return "stderr"
	}
	return "stdout"
}
