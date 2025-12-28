package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/Tfc538/core-cli/engine/update"
	"github.com/Tfc538/core-cli/version"
)

// NewUpdateCheckCmd creates the `core update check` command.
func NewUpdateCheckCmd() *cobra.Command {
	var jsonOutput bool

	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Check for available CORE CLI updates",
		Long:  "Check the GitHub Releases to see if a newer version of CORE CLI is available.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdateCheck(jsonOutput)
		},
	}

	checkCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	return checkCmd
}

// runUpdateCheck performs the update check.
func runUpdateCheck(jsonOutput bool) error {
	checker := update.NewChecker(update.CheckerConfig{
		GitHubOwner:    "Tfc538",
		GitHubRepo:     "core-cli",
		CurrentVersion: version.Version,
		GitHubToken:    githubToken(),
	})

	info, err := checker.Check()
	if err != nil {
		return fmt.Errorf("update check failed: %w", err)
	}

	if jsonOutput {
		return outputJSON(info)
	}

	// Human-readable output with formatted table
	out := NewOutputHelper()

	out.Table("Current version", info.CurrentVersion)
	out.Table("Latest version", info.LatestVersion)

	if info.UpdateAvailable {
		out.Separator()
		out.Success("Update available!")
		out.Separator()
		fmt.Println("Run 'core update apply' to update.")
	} else {
		out.Separator()
		out.Info("You are already on the latest version.")
	}

	return nil
}
