package cmd

import (
	"context"
	"strconv"
	"strings"
	"workshop/internal/cmd/internal/view"
	"workshop/internal/domain/species"

	"github.com/urfave/cli/v3"
)

func animalsCmd(app *app) *cli.Command {
	return &cli.Command{
		Name:  "animals",
		Usage: "List the animals the app recognises and the GBIF taxa behind them",
		Action: func(_ context.Context, _ *cli.Command) error {
			view.Render(app.Printer,
				view.Table[species.Animal]{
					Title:   "Animals",
					Headers: []string{"ANIMAL", "GBIF TAXON KEYS"},
					Items:   species.Animals(),
					Row: func(name species.Animal) []string {
						taxa := name.TaxonKeys()
						keys := make([]string, 0, len(taxa))
						for _, k := range taxa {
							keys = append(keys, strconv.Itoa(k))
						}
						return []string{string(name), strings.Join(keys, ", ")}
					},
				},
			)
			return nil
		},
	}
}
