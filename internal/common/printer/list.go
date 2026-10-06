package printer

import (
	"workshop/internal/common/printer/styles"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
)

// List prints a simple bullet list with brand colors.
func (c *Console) List(items ...any) {
	baseColor := styles.Base.GetForeground()
	itemStyle := lipgloss.NewStyle().Foreground(baseColor).Bold(true)
	enumeratorStyle := lipgloss.NewStyle().Foreground(baseColor).MarginRight(1).Bold(true)

	l := list.New(items...).
		Enumerator(list.Dash).
		EnumeratorStyle(enumeratorStyle).
		ItemStyle(itemStyle)

	c.Println(l.String())
}
