//go:build windows

package utilities

import (
	"syscall"
	"unsafe"
)

// isRunningAsAdministrator checks if the current process is running with administrator privileges on Windows
func isRunningAsAdministrator() bool {
	// Windows API constants
	const (
		TOKEN_QUERY    = 0x0008
		TokenElevation = 20
	)

	// Get current process handle
	handle, err := syscall.GetCurrentProcess()
	if err != nil {
		return false
	}

	// Open process token
	var token syscall.Token
	err = syscall.OpenProcessToken(handle, TOKEN_QUERY, &token)
	if err != nil {
		return false
	}
	defer token.Close()

	// Get token elevation information
	var elevation uint32
	var size uint32
	err = syscall.GetTokenInformation(token, TokenElevation, (*byte)(unsafe.Pointer(&elevation)), uint32(unsafe.Sizeof(elevation)), &size)
	if err != nil {
		return false
	}

	// Check if token is elevated (has administrator privileges)
	return elevation != 0
}
