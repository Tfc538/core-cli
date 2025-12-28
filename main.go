package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Tfc538/core-cli/cli"
	"github.com/Tfc538/core-cli/tui"
)

func main() {
	// Args-first: if arguments are provided, run CLI.
	// If invoked without arguments, launch TUI (unless explicitly requesting help/version in some cases).
	if len(os.Args) > 1 {
		// Arguments provided, run CLI
		runCLI()
	} else {
		// No arguments, launch TUI
		runTUI()
	}
}

// runCLI runs the CLI interface.
func runCLI() {
	rootCmd := cli.NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// runTUI launches the interactive TUI.
func runTUI() {
	model := tui.New()
	p := tea.NewProgram(model)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
