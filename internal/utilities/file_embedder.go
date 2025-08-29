package utilities

import (
	"embed"
	"fmt"
	"io/fs"

	"github.com/pkg/errors"
)

//go:embed files
var embeddedFiles embed.FS

// FileEmbedder manages embedded files
type FileEmbedder struct {
	files map[string]string
}

// NewFileEmbedder creates a new file embedder
func NewFileEmbedder() (*FileEmbedder, error) {
	embedder := &FileEmbedder{
		files: make(map[string]string),
	}

	if err := embedder.loadEmbeddedFiles(); err != nil {
		return nil, errors.Wrap(err, "failed to load embedded files")
	}

	return embedder, nil
}

// loadEmbeddedFiles loads all embedded files into memory
func (embedder *FileEmbedder) loadEmbeddedFiles() error {
	return fs.WalkDir(embeddedFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.Wrapf(err, "failed to walk embedded files at %s", path)
		}

		if d.IsDir() {
			return nil
		}

		// Read the file content
		content, err := embeddedFiles.ReadFile(path)
		if err != nil {
			return errors.Wrapf(err, "failed to read embedded file %s", path)
		}

		// Store the file content
		embedder.files[path] = string(content)

		return nil
	})
}

// GetFile retrieves a file by name
func (embedder *FileEmbedder) GetFile(name string) (string, bool) {
	content, exists := embedder.files[name]
	return content, exists
}

// ListFiles returns all available file names
func (embedder *FileEmbedder) ListFiles() []string {
	files := make([]string, 0, len(embedder.files))
	for name := range embedder.files {
		files = append(files, name)
	}
	return files
}

// GetDefaultConfig retrieves the default configuration
func (embedder *FileEmbedder) GetDefaultConfig() (string, error) {
	content, exists := embedder.GetFile("files/default-config.hcl")
	if !exists {
		return "", errors.New("default configuration file not found")
	}
	return content, nil
}

// PrintFileSummary prints a summary of embedded files
func (embedder *FileEmbedder) PrintFileSummary() {
	fmt.Println("=== Embedded Files Summary ===")

	files := embedder.ListFiles()
	if len(files) == 0 {
		fmt.Println("No files embedded")
		return
	}

	for _, name := range files {
		content, _ := embedder.GetFile(name)
		fmt.Printf("📄 %s (%d bytes)\n", name, len(content))
	}

	fmt.Printf("\n📊 Total: %d files embedded\n", len(files))
}
