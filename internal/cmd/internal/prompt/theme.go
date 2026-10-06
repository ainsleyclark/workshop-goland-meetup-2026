package prompt

import (
	"workshop/internal/common/printer/styles"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// Theme styles the interactive menus and forms to
// match the printer's colours.
func Theme() *huh.Theme {
	t := huh.ThemeCharm()

	t.Focused.Base = t.Focused.Base.BorderForeground(styles.ColorAccent)
	t.Focused.Card = t.Focused.Base
	t.Focused.Title = t.Focused.Title.Foreground(styles.ColorAccent)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(styles.ColorAccent)
	t.Focused.Description = t.Focused.Description.Foreground(styles.ColorMuted)
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(styles.ColorError)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(styles.ColorError)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(styles.ColorAccent)
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(styles.ColorAccent)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(styles.ColorAccent)
	t.Focused.Option = t.Focused.Option.Foreground(styles.ColorText)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(styles.ColorAccent)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(styles.ColorSuccess)
	t.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(styles.ColorSuccess).SetString(styles.IconSuccess + " ")
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(styles.ColorMuted).SetString(styles.IconDot + " ")
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(styles.ColorText)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Background(styles.ColorAccent)
	t.Focused.Next = t.Focused.FocusedButton
	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(styles.ColorAccent)
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(styles.ColorAccent)

	t.Blurred = t.Focused
	t.Blurred.Base = t.Focused.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Card = t.Blurred.Base
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description

	return t
}
