// This file is based on an example from https://github.com/gliderlabs/ssh/tree/master/_examples
//
// Original source: https://github.com/gliderlabs/ssh/tree/master/_examples

package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/gliderlabs/ssh"
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
		io.WriteString(s, fmt.Sprintf("Hello %s\n", s.User()))
	})

	fmt.Printf("Starting simple SSH server on :%d\n", port)
	fmt.Println("Press Ctrl+C to stop the server")
	log.Fatal(ssh.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
