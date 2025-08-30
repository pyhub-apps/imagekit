package main

import (
	"fmt"
	"os"
	
	"github.com/allieus/pyhub-imagekit/pkg/cli"
)

// Version is set at build time
var Version = "dev"

func main() {
	// Set version for CLI
	cli.SetVersion(Version)
	
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}