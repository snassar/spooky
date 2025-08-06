package outputs

import (
	"os"
	"path/filepath"
)

// FileOutput implements file output for log entries
type FileOutput struct {
	file *os.File
	path string
}

// NewFileOutput creates a new file output
func NewFileOutput(path string) (*FileOutput, error) {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	// Open file for writing (create if doesn't exist, append if exists)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}

	return &FileOutput{
		file: file,
		path: path,
	}, nil
}

// Write writes data to the file
func (o *FileOutput) Write(data []byte) error {
	_, err := o.file.Write(data)
	return err
}

// Close closes the file output
func (o *FileOutput) Close() error {
	return o.file.Close()
}

// GetName returns the output name
func (o *FileOutput) GetName() string {
	return "file:" + o.path
}
