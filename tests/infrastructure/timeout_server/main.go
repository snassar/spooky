// This file is based on an example from https://github.com/gliderlabs/ssh/tree/master/_examples
//
// Original source: https://github.com/gliderlabs/ssh/tree/master/_examples

package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/gliderlabs/ssh"
)

var (
	DeadlineTimeout = 30 * time.Second
	IdleTimeout     = 10 * time.Second
)

// findAvailablePort finds an available port between 3100-4100
func findAvailablePort() (int, error) {
	rand.Seed(time.Now().UnixNano())

	for attempt := 0; attempt < 5; attempt++ {
		port := rand.Intn(1001) + 3100 // Random port between 3100-4100

		// Try to listen on the port to check if it's available
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			listener.Close()
			return port, nil
		}
	}

	return 0, fmt.Errorf("could not find available port after 5 attempts")
}

func main() {
	port, err := findAvailablePort()
	if err != nil {
		log.Fatalf("Failed to find available port: %v", err)
	}

	ssh.Handle(func(s ssh.Session) {
		log.Println("new connection")
		i := 0
		for {
			i += 1
			log.Println("active seconds:", i)
			select {
			case <-time.After(time.Second):
				continue
			case <-s.Context().Done():
				log.Println("connection closed")
				return
			}
		}
	})

	fmt.Printf("Starting timeout SSH server on :%d\n", port)
	fmt.Printf("Connections will only last %s\n", DeadlineTimeout)
	fmt.Printf("and timeout after %s of no activity\n", IdleTimeout)
	fmt.Println("Press Ctrl+C to stop the server")

	// Allow all connections with any password
	passwordAuth := ssh.PasswordAuth(func(ctx ssh.Context, password string) bool {
		return true // accept any password
	})

	// Create server with timeout settings
	server := &ssh.Server{
		Addr:        fmt.Sprintf(":%d", port),
		MaxTimeout:  DeadlineTimeout,
		IdleTimeout: IdleTimeout,
	}

	// Set up authentication
	server.SetOption(passwordAuth)

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// For testing purposes, exit after a reasonable time
	// This prevents the server from hanging indefinitely
	time.Sleep(DeadlineTimeout + 5*time.Second)
	log.Println("Server shutting down after timeout period")
}
