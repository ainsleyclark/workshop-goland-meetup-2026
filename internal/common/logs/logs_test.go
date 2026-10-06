package logs

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) (*slog.Logger, *bytes.Buffer) {
	t.Helper()

	output := new(bytes.Buffer)
	log := NewProduction(output, slog.LevelInfo)

	return log, output
}

func TestNewProduction(t *testing.T) {
	t.Run("Filters entries below the configured level", func(t *testing.T) {
		log, output := setup(t)

		log.Debug("hidden")

		got := output.String()
		assert.Empty(t, got)
	})

	t.Run("Writes structured JSON with bound attributes", func(t *testing.T) {
		log, output := setup(t)
		log = log.With("service", "web")

		log.Error("failed", "status", 500)

		var got map[string]any
		require.NoError(t, json.Unmarshal(output.Bytes(), &got))

		assert.Equal(t, "failed", got["msg"])
		assert.Equal(t, "web", got["service"])
		assert.Equal(t, "ERROR", got["level"])
		assert.Equal(t, float64(500), got["status"])
	})
}

func TestNewLocal(t *testing.T) {
	t.Setenv("CLICOLOR_FORCE", "")
	t.Setenv("NO_COLOR", "")

	t.Run("Filters entries below the configured level", func(t *testing.T) {
		var output bytes.Buffer
		log := NewLocal(&output, slog.LevelInfo)

		log.Debug("hidden")

		got := output.String()
		assert.Empty(t, got)
	})

	t.Run("Writes readable attributes without terminal escapes", func(t *testing.T) {
		var output bytes.Buffer
		log := NewLocal(&output, slog.LevelInfo).With("service", "web")

		log.Info("HTTP request", "status", 200)

		got := output.String()
		assert.Contains(t, got, "HTTP request")
		assert.Contains(t, got, "service=web")
		assert.Contains(t, got, "status=200")
		assert.NotContains(t, got, "\x1b[")
	})
}

func TestNewLocalForcedColour(t *testing.T) {
	for _, tc := range []struct {
		name, force, noColour string
		wantColour            bool
	}{
		{"forced for piped output", "1", "", true},
		{"force disabled", "0", "", false},
		{"NO_COLOR takes precedence", "1", "1", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("CLICOLOR_FORCE", tc.force)
			t.Setenv("NO_COLOR", tc.noColour)
			var output bytes.Buffer
			NewLocal(&output, slog.LevelInfo).Info("HTTP request")
			assert.Equal(t, tc.wantColour, bytes.Contains(output.Bytes(), []byte("\x1b[")))
		})
	}
}
