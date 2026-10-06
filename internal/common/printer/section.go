package printer

import (
	"fmt"
	"strings"
	"workshop/internal/common/printer/styles"

	"github.com/charmbracelet/lipgloss"
)

// Title prints a section heading, preceded by a blank line
// so consecutive sections have room to breathe.
func (c *Console) Title(msg string) {
	c.LineBreak()
	c.Println(styles.Title.Render(styles.IconTitle + " " + msg))
}

// Titlef formats and prints a section heading.
func (c *Console) Titlef(format string, args ...any) {
	c.Title(fmt.Sprintf(format, args...))
}

// Muted prints a dim hint, such as a footer or empty state.
func (c *Console) Muted(msg string) {
	c.Println(styles.Muted.Render(msg))
}

// Mutedf formats and prints a dim hint.
func (c *Console) Mutedf(format string, args ...any) {
	c.Muted(fmt.Sprintf(format, args...))
}

// KV is a single row in a key/value view.
type KV struct {
	Key   string
	Value string
}

// KeyValue prints aligned key/value rows, which suits detail
// views. Rows with an empty value are skipped.
func (c *Console) KeyValue(pairs ...KV) {
	c.Println(KeyValueString(pairs...))
}

// KeyValueString renders aligned key/value rows without printing
// them, so they can be embedded inside a Box.
func KeyValueString(pairs ...KV) string {
	width := 0
	for _, p := range pairs {
		if p.Value != "" {
			width = max(width, lipgloss.Width(p.Key))
		}
	}

	label := styles.Label.Width(width + 2)
	lines := make([]string, 0, len(pairs))
	for _, p := range pairs {
		if p.Value == "" {
			continue
		}
		lines = append(lines, label.Render(p.Key)+styles.Value.Render(p.Value))
	}

	return strings.Join(lines, "\n")
}

// Box prints body inside a rounded card with an optional title.
func (c *Console) Box(title, body string) {
	content := body
	if title != "" {
		content = styles.Title.Render(title) + "\n" + body
	}
	c.Println(styles.Box.Render(content))
}

// Badge returns text as an inline pill in the given colour,
// for use inside tables, boxes and messages.
func Badge(text string, color lipgloss.TerminalColor) string {
	return lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Render("[" + strings.ToUpper(text) + "]")
}

// Status prints a single-line progress message that is replaced
// by the next call. It's a no-op when output isn't a terminal.
func (c *Console) Status(msg string) {
	if !c.isTerminal() {
		return
	}
	c.write("\r\x1b[2K" + styles.Muted.Render(msg))
}

// Statusf formats and prints a single-line progress message.
func (c *Console) Statusf(format string, args ...any) {
	c.Status(fmt.Sprintf(format, args...))
}

// ClearStatus removes the line written by Status.
func (c *Console) ClearStatus() {
	if !c.isTerminal() {
		return
	}
	c.write("\r\x1b[2K")
}
