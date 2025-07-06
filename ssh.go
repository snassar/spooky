package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHClient represents an SSH client connection
type SSHClient struct {
	client *ssh.Client
	server *Server
}

// NewSSHClient creates a new SSH client connection
func NewSSHClient(server *Server, timeout int) (*SSHClient, error) {
	var authMethods []ssh.AuthMethod

	// Add password authentication if provided
	if server.Password != "" {
		authMethods = append(authMethods, ssh.Password(server.Password))
	}

	// Add key-based authentication if provided
	if server.KeyFile != "" {
		key, err := os.ReadFile(server.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read key file %s: %w", server.KeyFile, err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("no authentication method available for server %s", server.Name)
	}

	// SSH client configuration
	config := &ssh.ClientConfig{
		User:            server.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(timeout) * time.Second,
	}

	// Connect to the server
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", server.Host, server.Port), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s@%s:%d: %w", server.User, server.Host, server.Port, err)
	}

	return &SSHClient{
		client: client,
		server: server,
	}, nil
}

// Close closes the SSH connection
func (s *SSHClient) Close() error {
	return s.client.Close()
}

// ExecuteCommand executes a command on the remote server
func (s *SSHClient) ExecuteCommand(command string) (string, error) {
	session, err := s.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Capture stdout and stderr
	var stdout, stderr strings.Builder
	session.Stdout = &stdout
	session.Stderr = &stderr

	// Execute the command
	err = session.Run(command)
	if err != nil {
		return "", fmt.Errorf("command execution failed: %s, stderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// ExecuteScript executes a script file on the remote server
func (s *SSHClient) ExecuteScript(scriptPath string) (string, error) {
	// Read the script file
	scriptContent, err := os.ReadFile(scriptPath)
	if err != nil {
		return "", fmt.Errorf("failed to read script file %s: %w", scriptPath, err)
	}

	// Execute the script content
	return s.ExecuteCommand(string(scriptContent))
}

// executeConfig executes all actions in the configuration
func executeConfig(config *Config) error {
	fmt.Printf("🚀 Starting execution of %d actions...\n", len(config.Actions))

	for _, action := range config.Actions {
		fmt.Printf("\n⚡ Executing action: %s\n", action.Name)
		if action.Description != "" {
			fmt.Printf("📝 Description: %s\n", action.Description)
		}

		// Get target servers for this action
		targetServers, err := getServersForAction(&action, config)
		if err != nil {
			return fmt.Errorf("failed to get servers for action %s: %w", action.Name, err)
		}

		fmt.Printf("🌐 Target servers: %d\n", len(targetServers))

		// Execute on each server
		if action.Parallel || parallel {
			err = executeActionParallel(&action, targetServers)
		} else {
			err = executeActionSequential(&action, targetServers)
		}

		if err != nil {
			return fmt.Errorf("failed to execute action %s: %w", action.Name, err)
		}
	}

	fmt.Println("\n✅ All actions completed successfully!")
	return nil
}

// executeActionSequential executes an action sequentially on all target servers
func executeActionSequential(action *Action, servers []*Server) error {
	for _, server := range servers {
		fmt.Printf("  🔗 Connecting to %s (%s@%s:%d)...\n", server.Name, server.User, server.Host, server.Port)

		// Create SSH client
		client, err := NewSSHClient(server, timeout)
		if err != nil {
			fmt.Printf("  ❌ Failed to connect to %s: %v\n", server.Name, err)
			continue
		}
		defer client.Close()

		// Execute the action
		var output string
		if action.Command != "" {
			output, err = client.ExecuteCommand(action.Command)
		} else if action.Script != "" {
			output, err = client.ExecuteScript(action.Script)
		}

		if err != nil {
			fmt.Printf("  ❌ Failed to execute action on %s: %v\n", server.Name, err)
			continue
		}

		fmt.Printf("  ✅ Success on %s\n", server.Name)
		if output != "" {
			fmt.Printf("  📄 Output:\n%s\n", indentOutput(output))
		}
	}

	return nil
}

// executeActionParallel executes an action in parallel on all target servers
func executeActionParallel(action *Action, servers []*Server) error {
	var wg sync.WaitGroup
	results := make(chan string, len(servers))
	errors := make(chan error, len(servers))

	for _, server := range servers {
		wg.Add(1)
		go func(s *Server) {
			defer wg.Done()

			fmt.Printf("  🔗 Connecting to %s (%s@%s:%d)...\n", s.Name, s.User, s.Host, s.Port)

			// Create SSH client
			client, err := NewSSHClient(s, timeout)
			if err != nil {
				errors <- fmt.Errorf("failed to connect to %s: %w", s.Name, err)
				return
			}
			defer client.Close()

			// Execute the action
			var output string
			if action.Command != "" {
				output, err = client.ExecuteCommand(action.Command)
			} else if action.Script != "" {
				output, err = client.ExecuteScript(action.Script)
			}

			if err != nil {
				errors <- fmt.Errorf("failed to execute action on %s: %w", s.Name, err)
				return
			}

			results <- fmt.Sprintf("✅ Success on %s\n%s", s.Name, indentOutput(output))
		}(server)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(results)
	close(errors)

	// Collect results
	for result := range results {
		fmt.Println(result)
	}

	// Check for errors
	select {
	case err := <-errors:
		return err
	default:
		return nil
	}
}

// indentOutput indents the output for better readability
func indentOutput(output string) string {
	if output == "" {
		return ""
	}

	var indented strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		indented.WriteString("    " + scanner.Text() + "\n")
	}
	return indented.String()
}
