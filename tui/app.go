package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Tfc538/core-cli/engine/update"
	"github.com/Tfc538/core-cli/version"
)

// Model is the main Bubble Tea model for the TUI.
type Model struct {
	// Version info
	currentVersion string

	// Update status
	updateInfo     *update.UpdateInfo
	updateError    string
	hasCheckedUpdate bool
	updateInProgress bool

	// Update progress
	updateProgress update.UpdateProgress

	// UI state
	width  int
	height int
}

// New creates a new TUI model.
func New() *Model {
	return &Model{
		currentVersion: version.Version,
	}
}

// Init initializes the model and returns an initial command.
func (m Model) Init() tea.Cmd {
	// Check for updates in background
	return m.checkForUpdatesCmd()
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "u":
			// 'u' key: show/trigger update
			if m.updateInfo != nil && m.updateInfo.UpdateAvailable && !m.updateInProgress {
				m.updateInProgress = true
				// In a full implementation, this would launch the update
				// For now, we just flag it
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case updateCheckCompleteMsg:
		m.updateInfo = msg.info
		m.hasCheckedUpdate = true
		m.updateError = ""
		if msg.err != nil {
			m.updateError = msg.err.Error()
		}

	case updateProgressMsg:
		m.updateProgress = msg.progress
		if msg.progress.Stage == "complete" || msg.progress.Stage == "failed" {
			m.updateInProgress = false
		}
	}

	return m, nil
}

// View renders the TUI.
func (m Model) View() string {
	s := ""

	// Title
	s += renderTitle()

	// Main content (placeholder)
	s += "\n"
	s += "  Core CLI - Intent-driven Developer Control Plane\n"
	s += "\n"
	s += "  Press 'u' to check for updates or 'q' to quit\n"
	s += "\n"

	// Status bar
	s += renderStatusBar(m)

	return s
}

// checkForUpdatesCmd creates a command to check for updates.
func (m Model) checkForUpdatesCmd() tea.Cmd {
	return func() tea.Msg {
		checker := update.NewChecker(update.CheckerConfig{
			GitHubOwner:    "Tfc538",
			GitHubRepo:     "core-cli",
			CurrentVersion: m.currentVersion,
		})

		info, err := checker.Check()
		return updateCheckCompleteMsg{info, err}
	}
}

// updateCheckCompleteMsg is sent when an update check completes.
type updateCheckCompleteMsg struct {
	info *update.UpdateInfo
	err  error
}

// updateProgressMsg is sent to report update progress.
type updateProgressMsg struct {
	progress update.UpdateProgress
}

// renderTitle returns the title section.
func renderTitle() string {
	return "CORE CLI\n"
}

// renderStatusBar returns the status bar section.
func renderStatusBar(m Model) string {
	status := "Status: "

	if !m.hasCheckedUpdate {
		status += "checking for updates..."
	} else if m.updateError != "" {
		status += fmt.Sprintf("update check failed: %s", m.updateError)
	} else if m.updateInfo != nil && m.updateInfo.UpdateAvailable {
		status += fmt.Sprintf("↑ Update available: v%s", m.updateInfo.LatestVersion)
	} else {
		status += "✓ Up to date"
	}

	return "  " + status + "\n"
}
