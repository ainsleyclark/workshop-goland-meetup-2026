package printer

import (
	"strings"
	"workshop/internal/common/printer/styles"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// maxCellWidth caps how wide a single cell can grow before
// it's truncated, so long remarks don't wreck the layout.
const maxCellWidth = 48

// Table prints a styled table with headers and rows.
func (c *Console) Table(headers []string, rows [][]string) {
	cellStyle := lipgloss.NewStyle().Padding(0, 1).MaxWidth(maxCellWidth + 2)

	truncated := make([][]string, len(rows))
	for i, row := range rows {
		truncated[i] = make([]string, len(row))
		for j, cell := range row {
			truncated[i][j] = truncate(cell, maxCellWidth)
		}
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(styles.ColorMuted)).
		BorderRow(false).
		BorderColumn(false).
		StyleFunc(func(row, _ int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return styles.Header
			case row%2 == 0:
				return cellStyle.Foreground(styles.ColorText)
			default:
				return cellStyle.Foreground(styles.ColorSubtle)
			}
		}).
		Headers(headers...).
		Rows(truncated...)

	c.Println(t.String())
}

// TableOrEmpty prints the table, or the empty message
// when there are no rows to display.
func (c *Console) TableOrEmpty(headers []string, rows [][]string, empty string) {
	if len(rows) == 0 {
		c.Muted(empty)
		return
	}
	c.Table(headers, rows)
}

// truncate shortens s to at most n visible characters, adding an
// ellipsis when anything is cut. Styled strings are left alone.
func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n || strings.Contains(s, "\x1b") {
		return s
	}
	return string(r[:n-1]) + "…"
}
