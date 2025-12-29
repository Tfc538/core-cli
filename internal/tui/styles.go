package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles contains the Lip Gloss styles for the TUI.
type Styles struct {
	Title    lipgloss.Style
	Subtitle lipgloss.Style
	Status   lipgloss.Style
	Success  lipgloss.Style
	Error    lipgloss.Style
	Update   lipgloss.Style
	Progress lipgloss.Style
}

// NewStyles creates and returns the TUI styles.
func NewStyles() *Styles {
	return &Styles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("7")).
			MarginBottom(1),

		Subtitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Italic(true),

		Status: lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			MarginTop(1),

		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Bold(true),

		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			Bold(true),

		Update: lipgloss.NewStyle().
			Foreground(lipgloss.Color("3")).
			Bold(true),

		Progress: lipgloss.NewStyle().
			Foreground(lipgloss.Color("4")).
			Bold(true),
	}
}
