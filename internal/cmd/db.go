package cmd

import (
	"context"
	"fmt"
	"workshop/internal/cmd/internal/prompt"
	"workshop/internal/domain/sighting"

	"github.com/urfave/cli/v3"
)

func dbCmd(app *app) *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "Manage the database",
		Commands: []*cli.Command{
			{
				Name:  "ping",
				Usage: "Check database connectivity",
				Action: func(ctx context.Context, _ *cli.Command) error {
					if err := app.DB.PingContext(ctx); err != nil {
						return fmt.Errorf("pinging database: %w", err)
					}
					app.Printer.Success("Database connection OK")
					return nil
				},
			},
			{
				Name:  "reset",
				Usage: "Roll back every migration and re-apply them, deleting all data",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "yes",
						Aliases: []string{"y"},
						Usage:   "skip the confirmation prompt",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if !cmd.Bool("yes") {
						ok, err := prompt.Confirm("This deletes every sighting and species. Continue?")
						if err != nil {
							return err
						} else if !ok {
							app.Printer.Muted("Reset cancelled")
							return nil
						}
					}
					if _, err := app.Migrations.DownTo(ctx, 0); err != nil {
						return fmt.Errorf("rolling back migrations: %w", err)
					}
					if _, err := app.Migrations.Up(ctx); err != nil {
						return fmt.Errorf("applying migrations: %w", err)
					}
					app.Printer.Success("Database reset")
					return nil
				},
			},
			{
				Name:  "seed",
				Usage: "Apply migrations and ingest a starter set of sightings of every recognised animal",
				Action: func(ctx context.Context, _ *cli.Command) error {
					if _, err := app.Migrations.Up(ctx); err != nil {
						return fmt.Errorf("applying migrations: %w", err)
					}
					return runIngest(ctx, app, sighting.IngestRequest{
						Limit: 5000,
					})
				},
			},
			migrateCmd(app),
		},
	}
}
