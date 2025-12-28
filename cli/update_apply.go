package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Tfc538/core-cli/engine/update"
	"github.com/Tfc538/core-cli/version"
)

// NewUpdateApplyCmd creates the `core update apply` command.
func NewUpdateApplyCmd() *cobra.Command {
	var skipConfirm bool

	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply the latest CORE CLI update",
		Long:  "Download and apply the latest version of CORE CLI, replacing the current binary.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdateApply(skipConfirm)
		},
	}

	applyCmd.Flags().BoolVar(&skipConfirm, "yes", false, "Skip confirmation prompt")

	return applyCmd
}

// runUpdateApply performs the update application.
func runUpdateApply(skipConfirm bool) error {
	// First, check for available updates
	checker := update.NewChecker(update.CheckerConfig{
		GitHubOwner:    "Tfc538",
		GitHubRepo:     "core-cli",
		CurrentVersion: version.Version,
	})

	info, err := checker.Check()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if !info.UpdateAvailable {
		fmt.Println("You are already on the latest version.")
		return nil
	}

	// Get current binary path
	binaryPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to determine current binary path: %w", err)
	}

	// Show confirmation prompt
	if !skipConfirm {
		fmt.Printf("This will update CORE CLI from %s to %s.\n", info.CurrentVersion, info.LatestVersion)
		fmt.Printf("The binary will be replaced at: %s\n\n", binaryPath)

		response := promptUser("Continue? [y/N]: ")
		if response != "y" && response != "Y" {
			fmt.Println("Update cancelled.")
			return nil
		}
	}

	fmt.Println("\nStarting update...")

	// Create updater and set up progress reporting
	updater := update.NewUpdater(update.UpdaterConfig{
		DownloadURL: info.DownloadURL,
		ChecksumURL: info.ChecksumURL,
		TargetPath:  binaryPath,
	})

	updater.SetProgressCallback(func(progress update.UpdateProgress) {
		switch progress.Stage {
		case "downloading":
			if progress.BytesTotal > 0 {
				fmt.Printf("⬇  Downloading update... %d%% (%d MB / %d MB)\r",
					progress.Percent,
					progress.BytesDone/1024/1024,
					progress.BytesTotal/1024/1024)
			}
		case "verifying":
			fmt.Println("⬇  Verifying checksum...              ")
			fmt.Println("✓  Checksum verified")
		case "replacing":
			fmt.Println("⬇  Replacing binary...")
		case "complete":
			fmt.Println("✓  Update complete!                    ")
		case "failed":
			fmt.Printf("✗  Update failed: %v\n", progress.Error)
		}
	})

	// Apply the update
	if err := updater.Apply(); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	fmt.Printf("\nCORE CLI updated to v%s\n", info.LatestVersion)
	return nil
}

// promptUser prompts the user for input and returns the response.
func promptUser(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	return strings.TrimSpace(response)
}
