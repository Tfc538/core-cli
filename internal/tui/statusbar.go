package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// StatusBar provides a status bar component.
type StatusBar struct {
	styles *Styles
	width  int
}

// NewStatusBar creates a new status bar.
func NewStatusBar(width int) *StatusBar {
	return &StatusBar{
		styles: NewStyles(),
		width:  width,
	}
}

// RenderVersionStatus renders the version status section.
func (sb *StatusBar) RenderVersionStatus(currentVersion string, isLatest bool) string {
	if isLatest {
		return sb.styles.Success.Render("✓ Version: " + currentVersion)
	}
	return "Version: " + currentVersion
}

// RenderUpdateStatus renders the update status section.
func (sb *StatusBar) RenderUpdateStatus(latestVersion string, updateAvailable bool) string {
	if updateAvailable {
		return sb.styles.Update.Render("↑ Update: v" + latestVersion + " available")
	}
	return sb.styles.Success.Render("✓ Up to date")
}

// RenderCheckingStatus renders the checking status.
func (sb *StatusBar) RenderCheckingStatus() string {
	return "⏳ Checking for updates..."
}

// RenderError renders an error status.
func (sb *StatusBar) RenderError(err string) string {
	return sb.styles.Error.Render("✗ Error: " + err)
}

// RenderProgress renders progress information.
func (sb *StatusBar) RenderProgress(stage string, percent int) string {
	switch stage {
	case "downloading":
		bar := lipgloss.NewStyle().
			Foreground(lipgloss.Color("4")).
			Render(fmt.Sprintf("⬇ Downloading... %d%%", percent))
		return bar
	case "verifying":
		return "🔍 Verifying..."
	case "replacing":
		return "🔄 Replacing..."
	case "complete":
		return sb.styles.Success.Render("✓ Update complete!")
	case "failed":
		return sb.styles.Error.Render("✗ Update failed")
	}
	return ""
}

// RenderBar renders a complete status bar.
func (sb *StatusBar) RenderBar(content string) string {
	// Pad content to fit width
	padding := sb.width - lipgloss.Width(content)
	if padding < 0 {
		padding = 0
	}
	padded := content + fmt.Sprintf("%*s", padding, "")
	return lipgloss.NewStyle().
		Background(lipgloss.Color("8")).
		Foreground(lipgloss.Color("0")).
		Render(padded)
}
