package loading

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// XDGManagerImpl implements XDGManager interface
type XDGManagerImpl struct {
	configHome string
	configDirs []string
	dataHome   string
	dataDirs   []string
	cacheHome  string
	runtimeDir string
}

// NewXDGManager creates a new XDG manager
func NewXDGManager() *XDGManagerImpl {
	return &XDGManagerImpl{
		configHome: getXDGConfigHome(),
		configDirs: getXDGConfigDirs(),
		dataHome:   getXDGDataHome(),
		dataDirs:   getXDGDataDirs(),
		cacheHome:  getXDGCacheHome(),
		runtimeDir: getXDGRuntimeDir(),
	}
}

// GetConfigHome returns the XDG config home directory
func (x *XDGManagerImpl) GetConfigHome() string {
	return x.configHome
}

// GetConfigDirs returns the XDG config directories
func (x *XDGManagerImpl) GetConfigDirs() []string {
	return x.configDirs
}

// GetDataHome returns the XDG data home directory
func (x *XDGManagerImpl) GetDataHome() string {
	return x.dataHome
}

// GetDataDirs returns the XDG data directories
func (x *XDGManagerImpl) GetDataDirs() []string {
	return x.dataDirs
}

// GetCacheHome returns the XDG cache home directory
func (x *XDGManagerImpl) GetCacheHome() string {
	return x.cacheHome
}

// GetRuntimeDir returns the XDG runtime directory
func (x *XDGManagerImpl) GetRuntimeDir() string {
	return x.runtimeDir
}

// CreateConfigDirectories creates the necessary XDG directories
func (x *XDGManagerImpl) CreateConfigDirectories() error {
	dirs := []string{
		x.configHome,
		x.dataHome,
		x.cacheHome,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// Helper functions to get XDG directories
func getXDGConfigHome() string {
	if home := os.Getenv("XDG_CONFIG_HOME"); home != "" {
		return home
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		return ".config"
	}

	return filepath.Join(userHome, ".config")
}

func getXDGConfigDirs() []string {
	if dirs := os.Getenv("XDG_CONFIG_DIRS"); dirs != "" {
		return strings.Split(dirs, ":")
	}

	return []string{"/etc/xdg"}
}

func getXDGDataHome() string {
	if home := os.Getenv("XDG_DATA_HOME"); home != "" {
		return home
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		return ".local/share"
	}

	return filepath.Join(userHome, ".local", "share")
}

func getXDGDataDirs() []string {
	if dirs := os.Getenv("XDG_DATA_DIRS"); dirs != "" {
		return strings.Split(dirs, ":")
	}

	return []string{"/usr/local/share", "/usr/share"}
}

func getXDGCacheHome() string {
	if home := os.Getenv("XDG_CACHE_HOME"); home != "" {
		return home
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		return ".cache"
	}

	return filepath.Join(userHome, ".cache")
}

func getXDGRuntimeDir() string {
	if dir := os.Getenv("XDG_RUNTIME_DIR"); dir != "" {
		return dir
	}

	// Fallback to /tmp if XDG_RUNTIME_DIR is not set
	return "/tmp"
}
