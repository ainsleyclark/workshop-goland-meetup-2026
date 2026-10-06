package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"workshop/cmd/mockapi/internal/gbif"

	"github.com/urfave/cli/v3"
)

func scrapeCmd(app *app) *cli.Command {
	return &cli.Command{
		Name:  "scrape",
		Usage: "Refresh the mock's GBIF data from the live API",
		Action: func(ctx context.Context, _ *cli.Command) error {
			ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
			defer stop()

			return gbif.Scrape(ctx, app.Logger)
		},
	}
}
