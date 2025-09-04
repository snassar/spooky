//go:build !windows

package utilities

// isRunningAsAdministrator is a stub for non-Windows systems
// On Unix-like systems, administrator privileges are handled by isRunningAsRoot()
func isRunningAsAdministrator() bool {
	// This function should never be called on non-Windows systems
	// as isRunningAsRoot() handles the logic directly
	return false
}
