package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewUpdateApplyCmd creates the `core update apply` command.
func NewUpdateApplyCmd() *cobra.Command {
	var skipConfirm bool

	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply the latest CORE CLI update",
		Long:  "Download and apply the latest version of CORE CLI, replacing the current binary.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement self-update in Phase 3
			fmt.Println("Update applying not yet implemented")
			return nil
		},
	}

	applyCmd.Flags().BoolVar(&skipConfirm, "yes", false, "Skip confirmation prompt")

	return applyCmd
}
