// Package prompt asks for input interactively using huh: the command
// menu, the forms that fill in a command's flags, and confirmations.
package prompt

import (
	"context"
	"errors"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/x/term"
	"github.com/urfave/cli/v3"
)

// Height caps how many rows a list shows before it scrolls.
const Height = 12

// metadataKey is the command Metadata key holding its Func.
const metadataKey = "prompt"

// Func asks for a command's flags and arguments, returning them
// as they would be typed on the command line.
type Func func(ctx context.Context) ([]string, error)

// With attaches fn to a command's Metadata, so the menu asks for
// its input before running it.
func With(fn Func) map[string]any {
	return map[string]any{metadataKey: fn}
}

// For returns the Func attached to cmd, if it has one.
func For(cmd *cli.Command) (Func, bool) {
	fn, ok := cmd.Metadata[metadataKey].(Func)
	return fn, ok
}

// CanInteract reports whether there's a terminal to ask on.
func CanInteract() bool {
	return term.IsTerminal(os.Stdin.Fd()) && term.IsTerminal(os.Stdout.Fd())
}

// Form shows fields together, filling in the values they point to.
func Form(fields ...huh.Field) error {
	return huh.NewForm(huh.NewGroup(fields...)).WithTheme(Theme()).Run()
}

// Choose asks the user to pick one of options, which can be filtered
// by typing /, returning the picked value.
func Choose(title string, options ...huh.Option[string]) (string, error) {
	var picked string
	err := huh.NewSelect[string]().
		Title(title).
		Description("Type / to filter").
		Options(options...).
		Height(Height).
		Value(&picked).
		WithTheme(Theme()).
		Run()
	return picked, err
}

// Confirm asks a yes/no question, defaulting to no. It refuses
// to guess when there's no terminal to ask on.
func Confirm(question string) (bool, error) {
	if !CanInteract() {
		return false, errors.New("not running in a terminal, pass --yes to confirm")
	}

	var ok bool
	err := huh.NewConfirm().
		Title(question).
		Affirmative("Yes").
		Negative("No").
		Value(&ok).
		WithTheme(Theme()).
		Run()
	if errors.Is(err, huh.ErrUserAborted) {
		return false, nil
	}

	return ok, err
}
