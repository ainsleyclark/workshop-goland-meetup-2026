package cmd

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"slices"
	"strings"
	"workshop/internal/cmd/internal/prompt"
	"workshop/internal/common/printer"
	"workshop/internal/infra/config"

	"github.com/charmbracelet/huh"
	"github.com/urfave/cli/v3"
)

// Run executes the cli command and runs the program.
func Run() {
	ctx := context.Background()

	// Commands are built before the flags are parsed, so they share
	// this pointer and Before fills it in once we know the flags.
	app := &app{Printer: printer.New(os.Stdout)}
	teardown := func() {}

	cfg, err := config.Load()
	if err != nil {
		app.Printer.Error(err.Error())
		os.Exit(1)
	}
	app.Config = cfg

	cmd := &cli.Command{
		Name:                  "workshop",
		Usage:                 "Explore wildlife sightings from GBIF",
		EnableShellCompletion: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "show debug logs on stderr",
			},
			&cli.StringFlag{
				Name:  "db",
				Usage: "SQLite database URL",
				Value: cfg.DatabaseURL,
			},
			&cli.BoolFlag{
				Name:  "local",
				Usage: "use the local mock APIs for GBIF and Open-Meteo, starting them in the background if not already running",
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			// Commands picked from the interactive menu run this hook
			// again, since it belongs to their parent, so only
			// bootstrap the first time round.
			if app.DB == nil {
				level := slog.LevelInfo
				if cmd.Bool("verbose") {
					level = slog.LevelDebug
				}

				if cmd.Bool("local") {
					if err := ensureLocalMock(ctx, app.Printer); err != nil {
						return ctx, err
					}
				}
				cfg = withLocalAPIs(cfg, cmd.Bool("local"))
				cfg.DatabaseURL = cmd.String("db")

				bootstrapped, td, err := bootstrap(ctx, cfg, level)
				teardown = td
				if err != nil {
					return ctx, err
				}

				*app = *bootstrapped
			}

			// The first argument names the subcommand, if there is one.
			var command string
			if sub := cmd.Command(cmd.Args().First()); sub != nil {
				command = sub.Name
			}

			return ctx, checkMigrations(ctx, app, command)
		},
		Action:   interactive(app),
		Commands: commands(app),
	}

	err = cmd.Run(ctx, os.Args)
	teardown()
	if err != nil && !errors.Is(err, context.Canceled) {
		app.Printer.Error(err.Error())
		os.Exit(1)
	}
}

func commands(app *app) []*cli.Command {
	return []*cli.Command{
		sightingsCmd(app),
		speciesCmd(app),
		statsCmd(app),
		animalsCmd(app),
		weatherCmd(app),
		webCmd(app),
		dbCmd(app),
	}
}

// interactive shows a menu of commands, asks for the picked command's
// flags and runs it, then comes back to the menu until the user quits
// or runs one of restartCommands. Without a terminal it prints the
// usual help instead.
func interactive(app *app) cli.ActionFunc {
	return func(ctx context.Context, root *cli.Command) error {
		if !prompt.CanInteract() {
			return cli.ShowRootCommandHelp(root)
		}

		for {
			// Commands hold on to parsed flag values, so the
			// menu gets a freshly built set each time round.
			chain, err := prompt.Menu(commands(app))
			if errors.Is(err, huh.ErrUserAborted) || (err == nil && chain == nil) {
				return nil
			} else if err != nil {
				return err
			}

			ran, err := runPicked(ctx, app, chain)
			if errors.Is(err, context.Canceled) {
				return nil
			} else if err != nil {
				app.Printer.Error(err.Error())
			}
			app.Printer.LineBreak()

			// The migrations and code were built when the menu opened,
			// so a fix saved since wouldn't be seen. Stopping here
			// means the next try is a fresh build.
			if ran && slices.Contains(restartCommands, prompt.Path(chain)) {
				app.Printer.Muted("Run make run again to pick up any changes, migrations included.")
				return nil
			}
		}
	}
}

// restartCommands are the menu commands that change the database, after
// which the menu stops rather than coming back round.
var restartCommands = []string{
	"db migrate up",
	"db migrate down",
	"db reset",
	"db seed",
	"sightings ingest",
}

// runPicked asks for the command's flags, if it has a prompt, and runs
// it, reporting whether it ran. Cancelling the prompt with ctrl+c goes
// back to the menu.
func runPicked(ctx context.Context, app *app, chain []*cli.Command) (bool, error) {
	if err := checkMigrations(ctx, app, chain[0].Name); err != nil {
		return false, err
	}

	cmd := chain[len(chain)-1]
	var args []string
	if ask, ok := prompt.For(cmd); ok {
		var err error
		args, err = ask(ctx)
		if errors.Is(err, huh.ErrUserAborted) {
			return false, nil
		} else if err != nil {
			return false, err
		}
	}

	app.Printer.Muted(strings.TrimSpace("$ workshop " + prompt.Path(chain) + " " + strings.Join(args, " ")))

	return true, cmd.Run(ctx, append([]string{cmd.Name}, args...))
}
