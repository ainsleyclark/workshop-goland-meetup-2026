package weather

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func get(t *testing.T, target string) (int, []byte) {
	t.Helper()
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	body, err := io.ReadAll(rec.Body)
	require.NoError(t, err)
	return rec.Code, body
}

func TestHandleArchive(t *testing.T) {
	t.Parallel()

	t.Run("Single hour", func(t *testing.T) {
		t.Parallel()

		status, body := get(t, "/weather/v1/archive?latitude=-3.4&longitude=39.95&"+
			"start_hour=2024-01-15T12:00&end_hour=2024-01-15T12:00&hourly=temperature_2m,weather_code")
		require.Equal(t, http.StatusOK, status, string(body))

		var got hourlyResponse
		require.NoError(t, json.Unmarshal(body, &got))
		assert.Equal(t, -3.4, got.Latitude)
		assert.Equal(t, 39.95, got.Longitude)
		assert.Equal(t, []string{"2024-01-15T12:00"}, got.Hourly.Time)
		assert.Len(t, got.Hourly.Temperature2M, 1)
	})

	t.Run("Multi-hour range", func(t *testing.T) {
		t.Parallel()

		status, body := get(t, "/weather/v1/archive?latitude=-3.4&longitude=39.95&"+
			"start_hour=2024-01-15T00:00&end_hour=2024-01-15T03:00&hourly=temperature_2m")
		require.Equal(t, http.StatusOK, status, string(body))

		var got hourlyResponse
		require.NoError(t, json.Unmarshal(body, &got))
		assert.Len(t, got.Hourly.Time, 4)
		assert.Len(t, got.Hourly.Temperature2M, 4)
	})

	t.Run("Same request twice is byte-identical", func(t *testing.T) {
		t.Parallel()

		target := "/weather/v1/archive?latitude=-75&longitude=0&" +
			"start_hour=2024-07-15T00:00&end_hour=2024-07-15T00:00&hourly=temperature_2m"
		_, first := get(t, target)
		_, second := get(t, target)
		assert.Equal(t, first, second)
	})

	t.Run("Rejects bad parameters", func(t *testing.T) {
		t.Parallel()

		for name, query := range map[string]string{
			"Missing latitude":      "longitude=0&start_hour=2024-01-01T00:00&end_hour=2024-01-01T00:00&hourly=temperature_2m",
			"Out of range latitude": "latitude=999&longitude=0&start_hour=2024-01-01T00:00&end_hour=2024-01-01T00:00&hourly=temperature_2m",
			"Missing hourly":        "latitude=0&longitude=0&start_hour=2024-01-01T00:00&end_hour=2024-01-01T00:00",
			"End before start":      "latitude=0&longitude=0&start_hour=2024-01-01T05:00&end_hour=2024-01-01T00:00&hourly=temperature_2m",
			"Unparsable hour":       "latitude=0&longitude=0&start_hour=yesterday&end_hour=2024-01-01T00:00&hourly=temperature_2m",
		} {
			status, body := get(t, "/weather/v1/archive?"+query)
			assert.Equal(t, http.StatusBadRequest, status, name)

			var errBody struct {
				Error  bool   `json:"error"`
				Reason string `json:"reason"`
			}
			require.NoError(t, json.Unmarshal(body, &errBody), name)
			assert.True(t, errBody.Error, name)
			assert.NotEmpty(t, errBody.Reason, name)
		}
	})

	t.Run("Unknown endpoint", func(t *testing.T) {
		t.Parallel()
		status, _ := get(t, "/weather/v1/forecast")
		assert.Equal(t, http.StatusNotFound, status)
	})
}
