package main

import (
	"fmt"
	"os"
)

// Version information - set at build time via ldflags
var (
	version = "dev"
	commit  = "unknown"
)

func main() {
	fmt.Printf("Spooky version %s-%s\n", version, commit)
	fmt.Println("CLI interface not yet implemented - interface mismatches need to be resolved")
	os.Exit(0)
}
