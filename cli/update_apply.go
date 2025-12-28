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
	out := NewOutputHelper()

	// First, check for available updates
	checker := update.NewChecker(update.CheckerConfig{
		GitHubOwner:    "Tfc538",
		GitHubRepo:     "core-cli",
		CurrentVersion: version.Version,
		GitHubToken:    githubToken(),
	})

	info, err := checker.Check()
	if err != nil {
		out.Error(fmt.Sprintf("Check failed: %v", err))
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if !info.UpdateAvailable {
		out.Info("You are already on the latest version.")
		return nil
	}

	// Get current binary path
	binaryPath, err := os.Executable()
	if err != nil {
		out.Error(fmt.Sprintf("Failed to locate binary: %v", err))
		return fmt.Errorf("failed to determine current binary path: %w", err)
	}

	// Show confirmation prompt
	if !skipConfirm {
		out.Heading("Update Available")
		out.Table("Current version", info.CurrentVersion)
		out.Table("Latest version", info.LatestVersion)
		out.Table("Target location", binaryPath)
		out.Separator()

		response := promptUser("Continue with update? [y/N]: ")
		if response != "y" && response != "Y" {
			out.Info("Update cancelled.")
			return nil
		}
	}

	out.Separator()
	out.Progress("Starting update")
	out.Separator()

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
				mb := progress.BytesDone / 1024 / 1024
				totalMB := progress.BytesTotal / 1024 / 1024
				fmt.Printf("⬇  Downloading... %d%% (%d/%d MB)\r",
					progress.Percent, mb, totalMB)
			}
		case "verifying":
			fmt.Println("                                        ")
			out.Progress("Verifying checksum")
		case "replacing":
			out.Success("Checksum verified")
			out.Progress("Replacing binary")
		case "complete":
			fmt.Println("                                        ")
			out.Success("Update complete!")
		case "failed":
			out.Error(fmt.Sprintf("Update failed: %v", progress.Error))
		}
	})

	// Apply the update
	if err := updater.Apply(); err != nil {
		out.Error(fmt.Sprintf("Apply failed: %v", err))
		return fmt.Errorf("update failed: %w", err)
	}

	out.Separator()
	fmt.Printf("✓ CORE CLI updated to v%s\n", info.LatestVersion)
	return nil
}

// promptUser prompts the user for input and returns the response.
func promptUser(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	return strings.TrimSpace(response)
}
