package cmd

import (
	"workshop/internal/infra/config"

	"github.com/urfave/cli/v3"
)

// localGBIFBaseURL and localOpenMeteoBaseURL point at the mock API server
// (cmd/mockapi, started with make mock), for --local.
const (
	localGBIFBaseURL      = "http://localhost:8082/gbif/v1"
	localOpenMeteoBaseURL = "http://localhost:8082/weather/v1"
)

// withLocalAPIs points GBIF and Open-Meteo at the local mock server instead
// of the live ones, overriding any GBIF_BASE_URL/OPEN_METEO_BASE_URL from
// .env, so switching over when the internet drops doesn't need an edit.
func withLocalAPIs(cfg config.Config, local bool) config.Config {
	if local {
		cfg.GBIFBaseURL = localGBIFBaseURL
		cfg.OpenMeteoBaseURL = localOpenMeteoBaseURL
	}
	return cfg
}

// countryFlag is the shared country selector. The usage differs per
// command, since one filters saved sightings and the other fetches them.
func countryFlag(usage string) cli.Flag {
	return &cli.StringFlag{
		Name:    "country",
		Aliases: []string{"c"},
		Usage:   usage,
	}
}

// limitFlag caps how many records a command fetches or shows.
func limitFlag(value int, usage string) cli.Flag {
	return &cli.IntFlag{
		Name:    "limit",
		Aliases: []string{"n"},
		Usage:   usage,
		Value:   value,
	}
}
