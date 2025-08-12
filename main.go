// Package main provides the main entry point for the spooky application.
// This package initializes and runs the spooky CLI using Cobra.
package main

import (
	"spooky/cmd"
)

func main() {
	// Run the CLI using the cmd package
	cmd.Execute()
}
