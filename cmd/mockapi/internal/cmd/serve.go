package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"workshop/cmd/mockapi/internal/gbif"
	"workshop/cmd/mockapi/internal/weather"
	"workshop/internal/common/logs"
	"workshop/internal/common/server"

	"github.com/urfave/cli/v3"
)

func serveCmd(app *app) *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "Serve the mock APIs",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "port to listen on",
				Value:   8082,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
			defer stop()

			start := time.Now()
			store, err := gbif.Open(ctx)
			if err != nil {
				return err
			}
			defer store.Close() //nolint:errcheck

			occurrences, species, err := store.Stats(ctx)
			if err != nil {
				return err
			}
			app.Logger.Info("Loaded GBIF data",
				"occurrences", occurrences,
				"species", species,
				"took", time.Since(start).Round(time.Millisecond),
			)

			mux := http.NewServeMux()
			mux.Handle("/gbif/", store.Handler())
			mux.Handle("/weather/", weather.Handler())

			port := cmd.Int("port")
			base := fmt.Sprintf("http://localhost:%d", port)
			app.Logger.Info("Listening",
				"gbif", base+"/gbif/v1",
				"search", base+"/gbif/v1/occurrence/search",
				"weather", base+"/weather/v1",
				"archive", base+"/weather/v1/archive",
			)

			return server.New(port, logs.RequestLogger(app.Logger)(mux)).Start(ctx)
		},
	}
}
