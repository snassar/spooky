package utilities

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// SanitizePathForError removes sensitive path information from error messages
func SanitizePathForError(path string) string {
	// Replace home directory with ~
	if homeDir, err := os.UserHomeDir(); err == nil {
		if strings.HasPrefix(path, homeDir) {
			path = strings.Replace(path, homeDir, "~", 1)
		}
	}

	// Replace common system paths with generic names
	systemPaths := map[string]string{
		"/etc/":       "config/",
		"/root/":      "admin/",
		"/var/log/":   "logs/",
		"/tmp/":       "temp/",
		"/usr/local/": "local/",
		"/opt/":       "opt/",
	}

	for sysPath, replacement := range systemPaths {
		if strings.HasPrefix(path, sysPath) {
			path = strings.Replace(path, sysPath, replacement, 1)
			break
		}
	}

	return path
}

// LogSensitiveInfo logs sensitive information at debug level only
func LogSensitiveInfo(logger interface{}, message string, sensitiveData map[string]interface{}) {
	// This function would be used to log sensitive information at debug level
	// Implementation depends on the logging framework being used
	// For now, this is a placeholder that ensures sensitive data is only logged at debug level
}
