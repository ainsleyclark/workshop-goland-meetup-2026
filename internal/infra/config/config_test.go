package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"workshop/internal/infra/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func envDir(t *testing.T, contents string) {
	t.Helper()

	for _, key := range []string{"DB_URI", "PORT", "APP_ENV", "GBIF_BASE_URL"} {
		t.Setenv(key, "")
		require.NoError(t, os.Unsetenv(key))
	}

	dir := t.TempDir()
	if contents != "" {
		require.NoError(t, os.WriteFile(filepath.Join(dir, ".env"), []byte(contents), 0o600))
	}
	t.Chdir(dir)
}

func TestLoad(t *testing.T) {
	t.Run("Defaults", func(t *testing.T) {
		envDir(t, "")

		got, err := config.Load()
		require.NoError(t, err)
		assert.Equal(t, config.Config{DatabaseURL: "file:workshop.db", Port: 8080, Env: config.Development}, got)
	})

	t.Run("From file", func(t *testing.T) {
		envDir(t, "DB_URI=file:test.db\nPORT=9090\nAPP_ENV=production\nGBIF_BASE_URL=http://localhost:8081/gbif/v1\n")

		got, err := config.Load()
		require.NoError(t, err)
		assert.Equal(t, config.Config{
			DatabaseURL: "file:test.db",
			Port:        9090,
			Env:         config.Production,
			GBIFBaseURL: "http://localhost:8081/gbif/v1",
		}, got)
	})

	t.Run("Environment wins", func(t *testing.T) {
		envDir(t, "DB_URI=file:test.db\n")
		t.Setenv("DB_URI", "file:env.db")

		got, err := config.Load()
		require.NoError(t, err)
		assert.Equal(t, "file:env.db", got.DatabaseURL)
	})

	t.Run("Invalid port", func(t *testing.T) {
		envDir(t, "PORT=nope\n")

		_, err := config.Load()
		assert.Error(t, err)
	})

	t.Run("Malformed file", func(t *testing.T) {
		envDir(t, "PORT 8080\n")

		_, err := config.Load()
		assert.Error(t, err)
	})
}

func TestConfig_IsProduction(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input config.Environment
		want  bool
	}{
		"Production":  {input: config.Production, want: true},
		"Development": {input: config.Development, want: false},
		"Unset":       {input: "", want: false},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := config.Config{Env: test.input}.IsProduction()
			assert.Equal(t, test.want, got)
		})
	}
}
