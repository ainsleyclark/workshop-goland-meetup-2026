package logs

import (
	"io"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

// NewProduction writes structured JSON for log collectors.
func NewProduction(w io.Writer, level slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: level,
	}))
}

// NewLocal writes readable logs, respecting NO_COLOR and CLICOLOR_FORCE.
func NewLocal(w io.Writer, level slog.Level) *slog.Logger {
	file, ok := w.(*os.File)
	colour := ok && isatty.IsTerminal(file.Fd())
	if force := os.Getenv("CLICOLOR_FORCE"); force != "" && force != "0" {
		colour = true
	}
	if os.Getenv("NO_COLOR") != "" {
		colour = false
	}

	return slog.New(tint.NewTextHandler(w, &tint.Options{
		Level:      level,
		TimeFormat: "15:04:05.000",
		NoColor:    !colour,
	}))
}
