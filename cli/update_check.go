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
	})

	info, err := checker.Check()
	if err != nil {
		return fmt.Errorf("update check failed: %w", err)
	}

	if jsonOutput {
		return outputJSON(info)
	}

	// Human-readable output
	fmt.Printf("Current version: %s\n", info.CurrentVersion)
	fmt.Printf("Latest version:  %s\n", info.LatestVersion)
	fmt.Printf("Update available: ")

	if info.UpdateAvailable {
		fmt.Println("Yes")
		fmt.Println("\nRun 'core update apply' to update.")
	} else {
		fmt.Println("No")
		fmt.Println("\nYou are already on the latest version.")
	}

	return nil
}
