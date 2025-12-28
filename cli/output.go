package cli

import (
	"fmt"
	"io"
	"os"
)

// OutputHelper provides formatted output for CLI commands.
type OutputHelper struct {
	out io.Writer
	err io.Writer
}

// NewOutputHelper creates a new output helper.
func NewOutputHelper() *OutputHelper {
	return &OutputHelper{
		out: os.Stdout,
		err: os.Stderr,
	}
}

// Info prints an informational message.
func (h *OutputHelper) Info(msg string) {
	fmt.Fprintln(h.out, msg)
}

// Success prints a success message with checkmark.
func (h *OutputHelper) Success(msg string) {
	fmt.Fprintf(h.out, "✓  %s\n", msg)
}

// Error prints an error message with X mark.
func (h *OutputHelper) Error(msg string) {
	fmt.Fprintf(h.err, "✗  %s\n", msg)
}

// Warning prints a warning message.
func (h *OutputHelper) Warning(msg string) {
	fmt.Fprintf(h.out, "⚠  %s\n", msg)
}

// Progress prints a progress message with down arrow.
func (h *OutputHelper) Progress(msg string) {
	fmt.Fprintf(h.out, "⬇  %s\n", msg)
}

// ProgressPercent prints progress with percentage.
func (h *OutputHelper) ProgressPercent(msg string, percent int) {
	fmt.Fprintf(h.out, "⬇  %s... %d%%\r", msg, percent)
}

// Table prints a simple two-column formatted table.
func (h *OutputHelper) Table(label, value string) {
	fmt.Fprintf(h.out, "%-20s: %s\n", label, value)
}

// Separator prints a blank line.
func (h *OutputHelper) Separator() {
	fmt.Fprintln(h.out)
}

// Heading prints a formatted heading.
func (h *OutputHelper) Heading(msg string) {
	fmt.Fprintf(h.out, "\n%s\n", msg)
}
