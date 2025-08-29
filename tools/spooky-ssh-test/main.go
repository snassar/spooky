package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"spooky/internal/encryption"
	"spooky/internal/logging"
	"spooky/internal/schemas"
	"spooky/internal/ssh"
)

func main() {
	fmt.Println("🔧 Spooky SSH Client Test")
	fmt.Println("==========================")

	// Create a test SSH configuration
	sshConfig := &schemas.SpookySSHV1{
		Timeout:            30,
		KeepaliveInterval:  60,
		KeepaliveCount:     3,
		KeyScanTimeout:     10,
		KnownHostsStrict:   false, // For testing
		ConnectionPoolSize: 10,
	}

	// Create age encryption (optional for testing)
	var ageEncryption *encryption.AgeEncryption
	if os.Getenv("SSH_AUTH_SOCK") != "" {
		// Try to create age encryption if SSH agent is available
		var err error
		ageEncryption, err = encryption.NewAgeEncryption("", "")
		if err != nil {
			logger := logging.GetGlobalLogger()
			logger.Warn("failed to initialize age encryption for SSH test, continuing without encryption support",
				slog.String("error", err.Error()))
			// Continue with nil encryption - SSH manager will handle this gracefully
		}
	}

	// Create SSH manager
	manager := ssh.NewSimpleSSHManager(ageEncryption, sshConfig)

	// Create a test machine configuration
	testMachine := &schemas.MachinesMachineV1{
		Hostname: "localhost",
		Port:     22,
		User:     os.Getenv("USER"),
		Authentication: schemas.MachinesMachineAuthenticationV1{
			// Use password authentication for testing
			Password: schemas.MachinesMachineAuthenticationPasswordV1{
				Value:     "test_password", // This would be encrypted in real usage
				Encrypted: false,
			},
		},
	}

	fmt.Printf("Testing SSH connection to %s@%s:%d\n", testMachine.User, testMachine.Hostname, testMachine.Port)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := manager.TestConnection(ctx, testMachine)
	if err != nil {
		fmt.Printf("❌ SSH connection test failed: %v\n", err)
		fmt.Println("\nNote: This is expected if:")
		fmt.Println("- SSH server is not running on localhost")
		fmt.Println("- Authentication credentials are not configured")
		fmt.Println("- SSH agent is not available")
		return
	}

	fmt.Println("✅ SSH connection test successful!")

	// Test command execution
	fmt.Println("\nTesting command execution...")
	result, err := manager.ExecuteCommandOnMachine(ctx, testMachine, "echo 'Hello from Spooky SSH!' && uname -a")
	if err != nil {
		fmt.Printf("❌ Command execution failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Command executed successfully!\n")
	fmt.Printf("Exit Code: %d\n", result.ExitCode)
	fmt.Printf("Stdout: %s\n", result.Stdout)
	if result.Stderr != "" {
		fmt.Printf("Stderr: %s\n", result.Stderr)
	}
}
