package tui

import (
	"fmt"
	"strings"

	"github.com/Tfc538/core-cli/internal/engine/update"
	"github.com/charmbracelet/lipgloss"
)

// UpdateView displays the update progress.
type UpdateView struct {
	styles *Styles
	width  int
	height int
}

// NewUpdateView creates a new update view.
func NewUpdateView(width, height int) *UpdateView {
	return &UpdateView{
		styles: NewStyles(),
		width:  width,
		height: height,
	}
}

// Render renders the update view.
func (uv *UpdateView) Render(info *update.UpdateInfo, progress update.UpdateProgress) string {
	var s strings.Builder

	// Title
	s.WriteString(uv.styles.Title.Render("🔄 Update CORE CLI"))
	s.WriteString("\n\n")

	// Version info
	s.WriteString(fmt.Sprintf("  Current version: v%s\n", info.CurrentVersion))
	s.WriteString(fmt.Sprintf("  Latest version:  v%s\n", info.LatestVersion))
	s.WriteString("\n")

	// Progress bar
	s.WriteString(uv.renderProgressBar(progress))
	s.WriteString("\n\n")

	// Status message
	s.WriteString(uv.renderStatusMessage(progress))

	return s.String()
}

// renderProgressBar renders a visual progress bar.
func (uv *UpdateView) renderProgressBar(progress update.UpdateProgress) string {
	barWidth := uv.width - 10
	if barWidth < 20 {
		barWidth = 20
	}

	filled := (progress.Percent * barWidth) / 100
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("4")).
		Bold(true)

	label := ""
	switch progress.Stage {
	case "downloading":
		label = "Downloading"
		if progress.BytesTotal > 0 {
			mb := progress.BytesDone / 1024 / 1024
			totalMB := progress.BytesTotal / 1024 / 1024
			label += fmt.Sprintf(" (%d/%d MB)", mb, totalMB)
		}
	case "verifying":
		label = "Verifying"
	case "replacing":
		label = "Replacing"
	}

	return fmt.Sprintf("  [%s] %d%% %s", style.Render(bar), progress.Percent, label)
}

// renderStatusMessage renders the current status message.
func (uv *UpdateView) renderStatusMessage(progress update.UpdateProgress) string {
	var msg string

	switch progress.Stage {
	case "downloading":
		msg = "⬇ Downloading update..."
	case "verifying":
		msg = "🔍 Verifying checksum..."
	case "replacing":
		msg = "🔄 Installing update..."
	case "complete":
		return uv.styles.Success.Render("✓ Update completed successfully!")
	case "failed":
		return uv.styles.Error.Render(fmt.Sprintf("✗ Update failed: %v", progress.Error))
	default:
		msg = "Initializing..."
	}

	return "  " + msg
}

// RenderCheckingModal renders a modal while checking for updates.
func (uv *UpdateView) RenderCheckingModal() string {
	var s strings.Builder

	width := uv.width
	if width > 50 {
		width = 50
	}

	padding := (uv.width - width) / 2
	paddingStr := strings.Repeat(" ", padding)

	s.WriteString("\n")
	s.WriteString(paddingStr)
	s.WriteString(strings.Repeat("─", width))
	s.WriteString("\n")

	s.WriteString(paddingStr)
	s.WriteString(fmt.Sprintf("│ %-*s │\n", width-4, "⏳ Checking for updates..."))

	s.WriteString(paddingStr)
	s.WriteString(strings.Repeat("─", width))
	s.WriteString("\n")

	return s.String()
}

// RenderUpdateAvailableModal renders a modal notifying of available update.
func (uv *UpdateView) RenderUpdateAvailableModal(info *update.UpdateInfo) string {
	var s strings.Builder

	width := uv.width
	if width > 60 {
		width = 60
	}

	padding := (uv.width - width) / 2
	paddingStr := strings.Repeat(" ", padding)

	s.WriteString("\n")
	s.WriteString(paddingStr)
	s.WriteString(strings.Repeat("─", width))
	s.WriteString("\n")

	s.WriteString(paddingStr)
	s.WriteString(fmt.Sprintf("│ %-*s │\n", width-4, "↑ Update Available!"))

	s.WriteString(paddingStr)
	s.WriteString(fmt.Sprintf("│ %-*s │\n", width-4, fmt.Sprintf("v%s → v%s", info.CurrentVersion, info.LatestVersion)))

	s.WriteString(paddingStr)
	s.WriteString(fmt.Sprintf("│ %-*s │\n", width-4, "Press 'u' to update or 'q' to quit"))

	s.WriteString(paddingStr)
	s.WriteString(strings.Repeat("─", width))
	s.WriteString("\n")

	return s.String()
}
