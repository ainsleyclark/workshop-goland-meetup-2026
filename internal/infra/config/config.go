package config

import (
	"errors"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Config is the base configuration struct for the application.
// Which is read from the environment.
type Config struct {
	Env         Environment `env:"APP_ENV" envDefault:"development"`
	DatabaseURL string      `env:"DB_URI" envDefault:"file:workshop.db"`
	Port        int         `env:"PORT" envDefault:"8080"`
	GBIFBaseURL string      `env:"GBIF_BASE_URL"`

	// OpenMeteoBaseURL overrides the Open-Meteo archive endpoint. Empty
	// uses the real API.
	OpenMeteoBaseURL string `env:"OPEN_METEO_BASE_URL"`
}

// Environment names where the app is running.
type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
)

// IsProduction reports whether the app is running on the deployed site.
func (c Config) IsProduction() bool {
	return c.Env == Production
}

// Load reads the .env file into the environment.
//
// Variables already set in the environment win over the file.
// A missing .env file is fine, deployments set real environment variables.
func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, err
	}
	return env.ParseAs[Config]()
}
