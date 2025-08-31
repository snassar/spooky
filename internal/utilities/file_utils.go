package utilities

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// FilePermissions represents standard file permissions (rw-r--r--)
	FilePermissions = 0o644
	// DirectoryPermissions represents standard directory permissions (rwxr-xr-x)
	DirectoryPermissions = 0o755
	// ExecutableFile represents executable file permissions (rwxr-xr-x)
	ExecutableFile = 0o755
)

// WriteFileWithPermissions writes content to file with consistent permissions and directory creation
func WriteFileWithPermissions(path, content string, perm os.FileMode) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, DirectoryPermissions); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return os.WriteFile(path, []byte(content), perm)
}

// WriteFile writes content to file with default file permissions
func WriteFile(path, content string) error {
	return WriteFileWithPermissions(path, content, FilePermissions)
}

// WriteExecutableFile writes content to file with executable permissions
func WriteExecutableFile(path, content string) error {
	return WriteFileWithPermissions(path, content, ExecutableFile)
}
