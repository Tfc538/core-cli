package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewUpdateCheckCmd creates the `core update check` command.
func NewUpdateCheckCmd() *cobra.Command {
	var jsonOutput bool

	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Check for available CORE CLI updates",
		Long:  "Check the GitHub Releases to see if a newer version of CORE CLI is available.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement update checking in Phase 2
			fmt.Println("Update checking not yet implemented")
			return nil
		},
	}

	checkCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	return checkCmd
}
