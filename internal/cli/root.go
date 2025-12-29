package cli

import (
	"github.com/spf13/cobra"
)

// NewRootCmd creates and returns the root command for CORE CLI.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "core",
		Short: "CORE CLI - Intent-driven developer control plane",
		Long: `CORE CLI is a local, intent-driven developer control plane with a CLI-first interface.
All commands can be run headlessly via arguments, or launch an interactive TUI when invoked without arguments.`,
		// Run TUI if no args provided - handled in main.go
		RunE: func(cmd *cobra.Command, args []string) error {
			// If we reach here with args, show help
			return cmd.Help()
		},
	}

	// Add subcommands
	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.AddCommand(NewUpdateCmd())

	return rootCmd
}
