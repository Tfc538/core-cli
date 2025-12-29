package cli

import (
	"github.com/spf13/cobra"
)

// NewUpdateCmd creates the `core update` parent command.
func NewUpdateCmd() *cobra.Command {
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Manage CORE CLI updates",
		Long:  "Check for and apply updates to CORE CLI.",
	}

	// Add subcommands
	updateCmd.AddCommand(NewUpdateCheckCmd())
	updateCmd.AddCommand(NewUpdateApplyCmd())

	return updateCmd
}
