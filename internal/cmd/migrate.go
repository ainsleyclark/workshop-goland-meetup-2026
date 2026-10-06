package cmd

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"workshop/internal/cmd/internal/view"
	"workshop/internal/common/printer"

	"github.com/pressly/goose/v3"
	"github.com/urfave/cli/v3"
)

func migrateCmd(app *app) *cli.Command {
	return &cli.Command{
		Name:  "migrate",
		Usage: "Manage database migrations",
		Commands: []*cli.Command{
			{
				Name:  "status",
				Usage: "List applied and pending migrations",
				Action: func(ctx context.Context, _ *cli.Command) error {
					statuses, err := app.Migrations.Status(ctx)
					if err != nil {
						return fmt.Errorf("reading migration status: %w", err)
					}

					view.Render(app.Printer, view.Table[*goose.MigrationStatus]{
						Title:   "Migration status",
						Headers: []string{"VERSION", "STATUS", "FILE"},
						Items:   statuses,
						Empty:   "No migrations found",
						Row: func(s *goose.MigrationStatus) []string {
							return []string{
								fmt.Sprint(s.Source.Version),
								fmt.Sprint(s.State),
								filepath.Base(s.Source.Path),
							}
						},
					})
					return nil
				},
			},
			{
				Name:  "version",
				Usage: "Show the current and latest available migration versions.",
				Action: func(ctx context.Context, _ *cli.Command) error {
					current, target, err := app.Migrations.GetVersions(ctx)
					if err != nil {
						return fmt.Errorf("reading migration versions: %w", err)
					}

					view.Render(app.Printer, view.Details{Title: "Database version", Rows: []printer.KV{
						{Key: "Current", Value: fmt.Sprint(current)},
						{Key: "Latest", Value: fmt.Sprint(target)},
					}})
					return nil
				},
			},
			{
				Name:  "up",
				Usage: "Apply pending database migrations",
				Action: func(ctx context.Context, _ *cli.Command) error {
					results, err := app.Migrations.Up(ctx)
					if err != nil {
						return fmt.Errorf("applying migrations: %w", err)
					}
					if len(results) == 0 {
						app.Printer.Info("Database is already up to date")
						return nil
					}
					app.Printer.Successf("Applied %d migration(s)", len(results))
					return nil
				},
			},
			{
				Name:  "down",
				Usage: "Roll back the most recent migration",
				Action: func(ctx context.Context, _ *cli.Command) error {
					if _, err := app.Migrations.Down(ctx); err != nil {
						return fmt.Errorf("rolling back migration: %w", err)
					}
					app.Printer.Success("Rolled back the most recent migration")
					return nil
				},
			},
		},
	}
}

// checkMigrations returns an error when migrations are pending, unless
// the named top level command manages migrations itself or doesn't
// touch the database.
func checkMigrations(ctx context.Context, app *app, command string) error {
	exemptCommands := []string{"", "db", "help", "animals", "weather", "completion"}
	if slices.Contains(exemptCommands, command) {
		return nil
	}

	pending, err := app.Migrations.HasPending(ctx)
	if err != nil {
		return fmt.Errorf("checking database migrations: %w", err)
	} else if pending {
		return errors.New("database migrations are pending, run 'workshop db migrate up' first")
	}

	return nil
}
