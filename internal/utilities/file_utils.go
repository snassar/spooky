package utilities

import (
	"fmt"
	"os"
	"path/filepath"
	"spooky/internal/logging"
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
// This function ensures sensitive data is only logged when debug logging is enabled
func LogSensitiveInfo(logger interface{}, message string, sensitiveData map[string]interface{}) {
	// Get the global logger if no specific logger is provided
	var log *logging.Logger
	if logger == nil {
		log = logging.GetGlobalLogger()
	} else if l, ok := logger.(*logging.Logger); ok {
		log = l
	} else {
		// If logger is not the expected type, use global logger
		log = logging.GetGlobalLogger()
	}

	// Only log sensitive information at debug level
	// Convert sensitive data to slog attributes
	attrs := make([]any, 0, len(sensitiveData)*2)
	for key, value := range sensitiveData {
		// Sanitize sensitive keys to avoid accidental exposure
		sanitizedKey := sanitizeSensitiveKey(key)
		attrs = append(attrs, sanitizedKey, value)
	}

	// Log at debug level only
	log.Debug(message, attrs...)
}

// sanitizeSensitiveKey sanitizes sensitive key names to prevent accidental exposure
func sanitizeSensitiveKey(key string) string {
	// Replace common sensitive key patterns with generic names
	sensitivePatterns := map[string]string{
		"password":       "auth_data",
		"passwd":         "auth_data",
		"secret":         "auth_data",
		"token":          "auth_data",
		"key":            "auth_data",
		"private_key":    "auth_data",
		"privatekey":     "auth_data",
		"private-key":    "auth_data",
		"api_key":        "auth_data",
		"apikey":         "auth_data",
		"api-key":        "auth_data",
		"access_key":     "auth_data",
		"accesskey":      "auth_data",
		"access-key":     "auth_data",
		"secret_key":     "auth_data",
		"secretkey":      "auth_data",
		"secret-key":     "auth_data",
		"credential":     "auth_data",
		"credentials":    "auth_data",
		"auth":           "auth_data",
		"authentication": "auth_data",
	}

	keyLower := strings.ToLower(key)
	for pattern, replacement := range sensitivePatterns {
		if strings.Contains(keyLower, pattern) {
			return replacement
		}
	}

	return key
}
