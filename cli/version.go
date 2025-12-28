package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/Tfc538/core-cli/version"
)

// NewVersionCmd creates the `core version` command.
func NewVersionCmd() *cobra.Command {
	var jsonOutput bool

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display the current CORE CLI version",
		Long:  "Display the current version of CORE CLI along with build metadata.",
		RunE: func(cmd *cobra.Command, args []string) error {
			info := version.Get()

			if jsonOutput {
				return outputJSON(info)
			}

			fmt.Println(info.String())
			return nil
		},
	}

	versionCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	return versionCmd
}

// outputJSON marshals and prints data as JSON.
func outputJSON(data interface{}) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(b))
	return nil
}
