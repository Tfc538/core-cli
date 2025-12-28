package main

import (
	"fmt"
	"os"

	"github.com/Tfc538/core-cli/cli"
)

func main() {
	// Args-first: if arguments are provided, run CLI. Otherwise, we'd launch TUI.
	// For now, always run CLI (TUI implemented in Phase 5).
	rootCmd := cli.NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
