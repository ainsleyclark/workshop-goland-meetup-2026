package prompt

import (
	"fmt"
	"strings"
	"workshop/internal/common/printer/styles"

	"github.com/charmbracelet/huh"
	"github.com/urfave/cli/v3"
)

// Menu entries that aren't commands.
var (
	menuBack = &cli.Command{Name: "back"}
	menuQuit = &cli.Command{Name: "quit"}
)

// Menu walks the command tree one list at a time until a command
// without subcommands is picked, returning it along with its parents.
// It returns nil when the user quits.
func Menu(cmds []*cli.Command) ([]*cli.Command, error) {
	var chain []*cli.Command
	for {
		current, title := cmds, "What do you want to do?"
		if len(chain) > 0 {
			current = chain[len(chain)-1].Commands
			title = Path(chain)
		}

		width := 0
		for _, c := range current {
			width = max(width, len(c.Name))
		}
		opts := make([]huh.Option[*cli.Command], 0, len(current)+1)
		for _, c := range current {
			if c.Hidden {
				continue
			}
			opts = append(opts, huh.NewOption(fmt.Sprintf("%-*s  %s", width, c.Name, c.Usage), c))
		}
		if len(chain) > 0 {
			opts = append(opts, huh.NewOption(styles.IconArrow+" Back", menuBack))
		} else {
			opts = append(opts, huh.NewOption("Quit", menuQuit))
		}

		var picked *cli.Command
		err := huh.NewSelect[*cli.Command]().
			Title(title).
			Options(opts...).
			Value(&picked).
			WithTheme(Theme()).
			Run()
		if err != nil {
			return nil, err
		}

		switch {
		case picked == menuQuit:
			return nil, nil
		case picked == menuBack:
			chain = chain[:len(chain)-1]
		case len(picked.Commands) == 0:
			return append(chain, picked), nil
		default:
			chain = append(chain, picked)
		}
	}
}

// Path joins command names, e.g. "sightings ingest".
func Path(chain []*cli.Command) string {
	names := make([]string, 0, len(chain))
	for _, c := range chain {
		names = append(names, c.Name)
	}
	return strings.Join(names, " ")
}
