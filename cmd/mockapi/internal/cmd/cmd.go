package cmd

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"workshop/internal/common/logs"

	"github.com/urfave/cli/v3"
)

type app struct {
	Logger *slog.Logger
}

// Run executes the cli command and runs the program.
func Run() {
	ctx := context.Background()

	logger := logs.NewLocal(os.Stderr, slog.LevelInfo)
	a := &app{
		Logger: logger,
	}

	cmd := &cli.Command{
		Name:           "mockapi",
		Usage:          "Serve mock copies of the APIs the workshop uses",
		DefaultCommand: "serve",
		Commands: []*cli.Command{
			serveCmd(a),
			scrapeCmd(a),
		},
	}

	if err := cmd.Run(ctx, os.Args); err != nil && !errors.Is(err, context.Canceled) {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
