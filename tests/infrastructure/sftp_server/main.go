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
	"github.com/pkg/sftp"
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

// SftpHandler handler for SFTP subsystem
func SftpHandler(sess ssh.Session) {
	debugStream := io.Discard
	serverOptions := []sftp.ServerOption{
		sftp.WithDebug(debugStream),
	}
	server, err := sftp.NewServer(
		sess,
		serverOptions...,
	)
	if err != nil {
		log.Printf("sftp server init error: %s\n", err)
		return
	}
	if err := server.Serve(); err == io.EOF {
		server.Close()
		fmt.Println("sftp client exited session.")
	} else if err != nil {
		fmt.Println("sftp server completed with error:", err)
	}
}

func main() {
	port, err := findAvailablePort()
	if err != nil {
		log.Fatalf("Failed to find available port: %v", err)
	}

	ssh_server := ssh.Server{
		Addr: fmt.Sprintf("127.0.0.1:%d", port),
		SubsystemHandlers: map[string]ssh.SubsystemHandler{
			"sftp": SftpHandler,
		},
	}

	fmt.Printf("Starting SFTP server on 127.0.0.1:%d\n", port)
	fmt.Println("Press Ctrl+C to stop the server")
	log.Fatal(ssh_server.ListenAndServe())
}
