package printer

import (
	"slices"
	"workshop/internal/common/printer/styles"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

// Tree prints a directory-style tree view.
func (c *Console) Tree(root string, children ...any) {
	t := tree.
		Root(root).
		Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(lipgloss.NewStyle().Foreground(styles.ColorAccent)).
		RootStyle(lipgloss.NewStyle().Foreground(styles.ColorAccent).Bold(true)).
		ItemStyle(lipgloss.NewStyle().Foreground(styles.ColorInfo))

	for _, child := range children {
		switch v := child.(type) {
		case string:
			t = t.Child(v)
		case *tree.Tree:
			t = t.Child(v)
		}
	}

	c.Println(t.String())
}

// Lineage prints a staircase of nested levels, such as a taxonomy,
// with each value followed by its muted label. Empty values are skipped.
func (c *Console) Lineage(levels ...KV) {
	var node *tree.Tree
	for _, level := range slices.Backward(levels) {
		if level.Value == "" {
			continue
		}
		label := styles.Value.Bold(true).Render(level.Value) + " " + styles.Muted.Render(level.Key)
		t := tree.Root(label).
			Enumerator(tree.RoundedEnumerator).
			EnumeratorStyle(lipgloss.NewStyle().Foreground(styles.ColorAccent).MarginLeft(1))
		if node != nil {
			t = t.Child(node)
		}
		node = t
	}
	if node == nil {
		return
	}
	c.Println(node.String())
}
