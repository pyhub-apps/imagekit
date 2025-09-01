package main

import (
	"fmt"
	"os"
	"strings"
	
	"github.com/allieus/imagekit/pkg/cli"
	"github.com/allieus/imagekit/pkg/update"
)

// Version is set at build time
var Version = "dev"

func main() {
	// Set version for CLI
	cli.SetVersion(Version)
	
	// Execute CLI command
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	
	// Show update notification (except for update command itself)
	if len(os.Args) >= 2 && os.Args[1] != "update" && !strings.Contains(strings.Join(os.Args, " "), "--no-update-check") {
		if updater, err := update.NewUpdater(Version); err == nil {
			updater.ShowUpdateNotification()
		}
	}
}