package utilities

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

// OSInfo represents detailed operating system information
type OSInfo struct {
	OS          string // "linux", "darwin", "windows", "freebsd", "openbsd", "netbsd"
	Version     string // OS version information
	Arch        string // Architecture (amd64, arm64, etc.)
	Distro      string // Linux distribution (ubuntu, centos, etc.)
	Kernel      string // Kernel version
	IsContainer bool   // Whether running in a container
	IsWSL       bool   // Whether running in Windows Subsystem for Linux
	IsRoot      bool   // Whether running as root/administrator
}

// PathConfig represents OS-specific path configuration
type PathConfig struct {
	ConfigDir   string // Configuration directory
	LogDir      string // Log directory
	CacheDir    string // Cache directory
	DataDir     string // Data directory
	TempDir     string // Temporary directory
	UserHomeDir string // User home directory
	AppName     string // Application name for paths
	ConfigFile  string // Main config file path
	LogFile     string // Main log file path
	CacheFile   string // Main cache file path
}

// DetectOS detects the current operating system and provides detailed information
func DetectOS() (*OSInfo, error) {
	info := &OSInfo{
		OS:      runtime.GOOS,
		Arch:    runtime.GOARCH,
		Version: runtime.Version(),
	}

	// Detect if running as root/administrator
	info.IsRoot = isRunningAsRoot()

	// OS-specific detection
	switch info.OS {
	case "linux":
		return detectLinux(info)
	case "darwin":
		return detectDarwin(info)
	case "windows":
		return detectWindows(info)
	case "freebsd", "openbsd", "netbsd":
		return detectBSD(info)
	default:
		return nil, errors.Errorf("unsupported operating system: %s", info.OS)
	}
}

// detectLinux detects Linux-specific information
func detectLinux(info *OSInfo) (*OSInfo, error) {
	// Detect distribution
	info.Distro = detectLinuxDistro()

	// Detect kernel version
	info.Kernel = detectKernelVersion()

	// Detect if running in container
	info.IsContainer = detectContainer()

	// Detect if running in WSL
	info.IsWSL = detectWSL()

	return info, nil
}

// detectDarwin detects macOS-specific information
func detectDarwin(info *OSInfo) (*OSInfo, error) {
	// macOS version detection
	version, err := runCommand("sw_vers", "-productVersion")
	if err == nil {
		info.Version = strings.TrimSpace(version)
	}

	// Detect if running in container (Docker Desktop, etc.)
	info.IsContainer = detectContainer()

	return info, nil
}

// detectWindows detects Windows-specific information
func detectWindows(info *OSInfo) (*OSInfo, error) {
	// Windows version detection
	version, err := runCommand("cmd", "/c", "ver")
	if err == nil {
		info.Version = strings.TrimSpace(version)
	}

	// Check if Windows 11+ (simplified check)
	// In a real implementation, you'd parse the version string more carefully
	if strings.Contains(strings.ToLower(info.Version), "11") ||
		strings.Contains(strings.ToLower(info.Version), "12") {
		info.Version = "11+"
	} else {
		return nil, errors.New("Windows 11+ is required, older versions are not supported")
	}

	return info, nil
}

// detectBSD detects BSD-specific information
func detectBSD(info *OSInfo) (*OSInfo, error) {
	// BSD version detection
	switch info.OS {
	case "freebsd":
		version, err := runCommand("freebsd-version")
		if err == nil {
			info.Version = strings.TrimSpace(version)
		}
	case "openbsd":
		version, err := runCommand("uname", "-r")
		if err == nil {
			info.Version = strings.TrimSpace(version)
		}
	case "netbsd":
		version, err := runCommand("uname", "-r")
		if err == nil {
			info.Version = strings.TrimSpace(version)
		}
	}

	return info, nil
}

// detectLinuxDistro detects the Linux distribution
func detectLinuxDistro() string {
	// Try common distribution files
	distroFiles := []string{
		"/etc/os-release",
		"/etc/lsb-release",
		"/etc/redhat-release",
		"/etc/debian_version",
		"/etc/gentoo-release",
		"/etc/SuSE-release",
	}

	for _, file := range distroFiles {
		if content, err := os.ReadFile(file); err == nil {
			contentStr := strings.ToLower(string(content))

			// Parse common distribution identifiers
			switch {
			case strings.Contains(contentStr, "ubuntu"):
				return "ubuntu"
			case strings.Contains(contentStr, "debian"):
				return "debian"
			case strings.Contains(contentStr, "centos"):
				return "centos"
			case strings.Contains(contentStr, "redhat"):
				return "redhat"
			case strings.Contains(contentStr, "fedora"):
				return "fedora"
			case strings.Contains(contentStr, "arch"):
				return "arch"
			case strings.Contains(contentStr, "gentoo"):
				return "gentoo"
			case strings.Contains(contentStr, "opensuse"):
				return "opensuse"
			case strings.Contains(contentStr, "alpine"):
				return "alpine"
			}
		}
	}

	return "unknown"
}

// detectKernelVersion detects the Linux kernel version
func detectKernelVersion() string {
	if content, err := os.ReadFile("/proc/version"); err == nil {
		parts := strings.Fields(string(content))
		if len(parts) >= 3 {
			return parts[2] // Kernel version is typically the 3rd field
		}
	}
	return "unknown"
}

// detectContainer detects if running in a container
func detectContainer() bool {
	// Check for /.dockerenv
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	// Check cgroup for container indicators
	if content, err := os.ReadFile("/proc/1/cgroup"); err == nil {
		contentStr := string(content)
		if strings.Contains(contentStr, "docker") ||
			strings.Contains(contentStr, "kubepods") ||
			strings.Contains(contentStr, "containerd") {
			return true
		}
	}

	return false
}

// detectWSL detects if running in Windows Subsystem for Linux
func detectWSL() bool {
	// Check for WSL-specific files
	wslFiles := []string{
		"/proc/version",
		"/proc/sys/kernel/osrelease",
	}

	for _, file := range wslFiles {
		if content, err := os.ReadFile(file); err == nil {
			contentStr := strings.ToLower(string(content))
			if strings.Contains(contentStr, "microsoft") ||
				strings.Contains(contentStr, "wsl") {
				return true
			}
		}
	}

	return false
}

// isRunningAsRoot detects if the process is running as root/administrator
func isRunningAsRoot() bool {
	switch runtime.GOOS {
	case "windows":
		// On Windows, check if running as administrator
		// This is a simplified check - in production you'd use Windows API
		return false // Placeholder
	default:
		// On Unix-like systems, check if UID is 0
		return os.Geteuid() == 0
	}
}

// runCommand runs a command and returns its output
func runCommand(name string, args ...string) (string, error) {
	// This is a simplified implementation
	// In production, you'd use os/exec to run commands
	return "", errors.New("command execution not implemented")
}

// GetPathConfig returns OS-specific path configuration
func GetPathConfig(appName string) (*PathConfig, error) {
	osInfo, err := DetectOS()
	if err != nil {
		return nil, errors.Wrap(err, "failed to detect OS")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user home directory")
	}

	config := &PathConfig{
		AppName:     appName,
		UserHomeDir: homeDir,
	}

	// Set OS-specific paths
	switch osInfo.OS {
	case "linux":
		setLinuxPaths(config, osInfo)
	case "darwin":
		setDarwinPaths(config, osInfo)
	case "windows":
		setWindowsPaths(config, osInfo)
	case "freebsd", "openbsd", "netbsd":
		setBSDPaths(config, osInfo)
	default:
		return nil, errors.Errorf("unsupported OS for path configuration: %s", osInfo.OS)
	}

	// Set file paths
	config.ConfigFile = filepath.Join(config.ConfigDir, appName+".hcl")
	config.LogFile = filepath.Join(config.LogDir, appName+".log")
	config.CacheFile = filepath.Join(config.CacheDir, appName+".cache")

	return config, nil
}

// setLinuxPaths sets Linux-specific paths
func setLinuxPaths(config *PathConfig, osInfo *OSInfo) {
	// Follow XDG Base Directory Specification
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		xdgConfigHome = filepath.Join(config.UserHomeDir, ".config")
	}

	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome == "" {
		xdgDataHome = filepath.Join(config.UserHomeDir, ".local", "share")
	}

	xdgCacheHome := os.Getenv("XDG_CACHE_HOME")
	if xdgCacheHome == "" {
		xdgCacheHome = filepath.Join(config.UserHomeDir, ".cache")
	}

	config.ConfigDir = filepath.Join(xdgConfigHome, config.AppName)
	config.DataDir = filepath.Join(xdgDataHome, config.AppName)
	config.CacheDir = filepath.Join(xdgCacheHome, config.AppName)
	config.LogDir = filepath.Join(xdgDataHome, config.AppName, "logs")
	config.TempDir = os.TempDir()
}

// setDarwinPaths sets macOS-specific paths
func setDarwinPaths(config *PathConfig, osInfo *OSInfo) {
	config.ConfigDir = filepath.Join(config.UserHomeDir, "Library", "Application Support", config.AppName)
	config.DataDir = filepath.Join(config.UserHomeDir, "Library", "Application Support", config.AppName)
	config.CacheDir = filepath.Join(config.UserHomeDir, "Library", "Caches", config.AppName)
	config.LogDir = filepath.Join(config.UserHomeDir, "Library", "Logs", config.AppName)
	config.TempDir = os.TempDir()
}

// setWindowsPaths sets Windows-specific paths
func setWindowsPaths(config *PathConfig, osInfo *OSInfo) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		appData = filepath.Join(config.UserHomeDir, "AppData", "Roaming")
	}

	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		localAppData = filepath.Join(config.UserHomeDir, "AppData", "Local")
	}

	config.ConfigDir = filepath.Join(appData, config.AppName)
	config.DataDir = filepath.Join(localAppData, config.AppName)
	config.CacheDir = filepath.Join(localAppData, config.AppName, "Cache")
	config.LogDir = filepath.Join(localAppData, config.AppName, "Logs")
	config.TempDir = os.TempDir()
}

// setBSDPaths sets BSD-specific paths
func setBSDPaths(config *PathConfig, osInfo *OSInfo) {
	// BSD systems typically follow Unix conventions
	config.ConfigDir = filepath.Join(config.UserHomeDir, ".config", config.AppName)
	config.DataDir = filepath.Join(config.UserHomeDir, ".local", "share", config.AppName)
	config.CacheDir = filepath.Join(config.UserHomeDir, ".cache", config.AppName)
	config.LogDir = filepath.Join(config.UserHomeDir, ".local", "share", config.AppName, "logs")
	config.TempDir = os.TempDir()
}

// EnsureDirectories creates all necessary directories for the application
func EnsureDirectories(config *PathConfig) error {
	dirs := []string{
		config.ConfigDir,
		config.DataDir,
		config.CacheDir,
		config.LogDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return errors.Wrapf(err, "failed to create directory: %s", dir)
		}
	}

	return nil
}

// String returns a string representation of OS information
func (info *OSInfo) String() string {
	return fmt.Sprintf("OS: %s, Version: %s, Arch: %s, Distro: %s, Container: %v, WSL: %v, Root: %v",
		info.OS, info.Version, info.Arch, info.Distro, info.IsContainer, info.IsWSL, info.IsRoot)
}

// IsSupported returns true if the OS is supported
func (info *OSInfo) IsSupported() bool {
	switch info.OS {
	case "linux", "darwin", "freebsd", "openbsd", "netbsd":
		return true
	case "windows":
		// Only Windows 11+ is supported
		return strings.Contains(info.Version, "11") || strings.Contains(info.Version, "12")
	default:
		return false
	}
}
