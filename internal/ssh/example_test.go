package ssh

import (
	"context"
	"fmt"
	"time"

	spookylogging "spooky/internal/logging"
	spookytypes "spooky/internal/types"
)

// Example demonstrates basic SSH client usage
func Example() {
	// Create a log manager
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("ssh-example")

	// Create SSH client configuration
	config := &spookytypes.ClientConfig{
		DefaultPort:       22,
		DefaultTimeout:    30 * time.Second,
		MaxConnections:    5,
		MaxRetryAttempts:  3,
		RetryDelay:        5 * time.Second,
		IdleTimeout:       300 * time.Second,
		DefaultKeyPath:    "~/.ssh/id_rsa",
		DefaultAuthMethod: spookytypes.AuthMethodPublicKey,
	}

	// Create SSH client
	client := NewSimpleClient(config, logger)

	// Create connection request
	request := &spookytypes.ConnectionRequest{
		Host:        "example.com",
		Port:        22,
		User:        "user",
		KeyPath:     "~/.ssh/id_rsa",
		Timeout:     30 * time.Second,
		RequestID:   "example-connection",
		RequestedAt: time.Now(),
	}

	// Connect to remote host
	ctx := context.Background()
	result, err := client.Connect(ctx, request)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}

	if !result.Success {
		fmt.Printf("Connection failed: %s\n", result.Error)
		return
	}

	fmt.Printf("Connected to %s:%d as %s\n",
		result.Connection.Host,
		result.Connection.Port,
		result.Connection.User)

	// Execute a command
	command := &spookytypes.SSHCommand{
		Command:       "uname",
		Args:          []string{"-a"},
		CaptureOutput: true,
		Timeout:       10 * time.Second,
	}

	cmdResult, err := client.RunCommand(ctx, result.Connection, command)
	if err != nil {
		fmt.Printf("Failed to run command: %v\n", err)
		return
	}

	if cmdResult.Success {
		fmt.Printf("Command output: %s\n", cmdResult.Stdout)
	} else {
		fmt.Printf("Command failed: %s\n", cmdResult.Stderr)
	}

	// Close the client
	client.Close(ctx)
}

// ExampleConnectionRequest demonstrates creating a connection request
func ExampleConnectionRequest() {
	request := &spookytypes.ConnectionRequest{
		Host:        "server.example.com",
		Port:        22,
		User:        "admin",
		KeyPath:     "~/.ssh/admin_key",
		Passphrase:  "secret",
		Timeout:     30 * time.Second,
		RequestID:   "admin-connection",
		RequestedAt: time.Now(),
	}

	fmt.Printf("Connection request for %s@%s:%d\n",
		request.User, request.Host, request.Port)
}

// ExampleSSHCommand demonstrates creating an SSH command
func ExampleSSHCommand() {
	command := &spookytypes.SSHCommand{
		Command:       "ls",
		Args:          []string{"-la", "/tmp"},
		WorkingDir:    "/home/user",
		Environment:   map[string]string{"LANG": "en_US.UTF-8"},
		CaptureOutput: true,
		Timeout:       30 * time.Second,
		Priority:      1,
		ScheduledAt:   time.Now(),
		Tags:          []string{"filesystem", "listing"},
	}

	fmt.Printf("Command: %s %v\n", command.Command, command.Args)
}
