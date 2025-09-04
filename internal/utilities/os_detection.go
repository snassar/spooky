package utilities

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// OS constants to avoid repeated string literals
const (
	OSLinux   = "linux"
	OSDarwin  = "darwin"
	OSWindows = "windows"
	OSFreeBSD = "freebsd"
	OSOpenBSD = "openbsd"
	OSNetBSD  = "netbsd"
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
	StateDir    string // State directory (for databases, logs, etc.)
	TempDir     string // Temporary directory
	UserHomeDir string // User home directory
	AppName     string // Application name for paths
	ConfigFile  string // Main config file path
	LogFile     string // Main log file path
	CacheFile   string // Main cache file path
	StateFile   string // Main state file path (for SQLite database)
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
	case OSLinux:
		return detectLinux(info)
	case OSDarwin:
		return detectDarwin(info)
	case OSWindows:
		return detectWindows(info)
	case OSFreeBSD, OSOpenBSD, OSNetBSD:
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
	case OSFreeBSD:
		version, err := runCommand("freebsd-version")
		if err == nil {
			info.Version = strings.TrimSpace(version)
		}
	case OSOpenBSD:
		version, err := runCommand("uname", "-r")
		if err == nil {
			info.Version = strings.TrimSpace(version)
		}
	case OSNetBSD:
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
	case OSWindows:
		return isRunningAsAdministrator()
	default:
		// On Unix-like systems, check if UID is 0
		return os.Geteuid() == 0
	}
}

// runCommand runs a command and returns its output
func runCommand(name string, args ...string) (string, error) {
	// Create a context with timeout to prevent hanging commands
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create the command
	cmd := exec.CommandContext(ctx, name, args...)

	// Execute the command and capture output
	output, err := cmd.Output()
	if err != nil {
		// Check if it's a timeout error
		if ctx.Err() == context.DeadlineExceeded {
			return "", errors.Wrap(err, "command timed out after 10 seconds")
		}
		// Check if command was not found
		if execErr, ok := err.(*exec.Error); ok && execErr.Err == exec.ErrNotFound {
			return "", errors.Wrap(err, "command not found")
		}
		// Other execution errors
		return "", errors.Wrap(err, "failed to execute command")
	}

	// Return the output as a string, trimming whitespace
	return strings.TrimSpace(string(output)), nil
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
	case OSLinux, OSDarwin, OSFreeBSD, OSOpenBSD, OSNetBSD:
		setXDGPaths(config, osInfo)
	case OSWindows:
		setWindowsPaths(config, osInfo)
	default:
		return nil, errors.Errorf("unsupported OS for path configuration: %s", osInfo.OS)
	}

	// Set file paths
	config.ConfigFile = filepath.Join(config.ConfigDir, appName+".hcl")
	config.LogFile = filepath.Join(config.LogDir, appName+".log")
	config.CacheFile = filepath.Join(config.CacheDir, appName+".cache")
	config.StateFile = filepath.Join(config.StateDir, appName+".db")

	return config, nil
}

// setXDGPaths sets paths following XDG Base Directory Specification
// Used for Linux, macOS, and BSD systems
func setXDGPaths(config *PathConfig, osInfo *OSInfo) {
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

	// XDG_STATE_HOME for application state (databases, logs, etc.)
	xdgStateHome := os.Getenv("XDG_STATE_HOME")
	if xdgStateHome == "" {
		xdgStateHome = filepath.Join(config.UserHomeDir, ".local", "state")
	}

	config.ConfigDir = filepath.Join(xdgConfigHome, config.AppName)
	config.DataDir = filepath.Join(xdgDataHome, config.AppName)
	config.CacheDir = filepath.Join(xdgCacheHome, config.AppName)
	config.StateDir = filepath.Join(xdgStateHome, config.AppName)
	config.LogDir = filepath.Join(xdgStateHome, config.AppName, "logs")
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
	config.StateDir = filepath.Join(localAppData, config.AppName, "State")
	config.LogDir = filepath.Join(localAppData, config.AppName, "State", "logs")
	config.TempDir = os.TempDir()
}

// EnsureDirectories creates all necessary directories for the application
func EnsureDirectories(config *PathConfig) error {
	dirs := []string{
		config.ConfigDir,
		config.DataDir,
		config.CacheDir,
		config.StateDir,
		config.LogDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
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
	case OSLinux, OSDarwin, OSFreeBSD, OSOpenBSD, OSNetBSD:
		return true
	case OSWindows:
		// Only Windows 11+ is supported
		return strings.Contains(info.Version, "11") || strings.Contains(info.Version, "12")
	default:
		return false
	}
}
