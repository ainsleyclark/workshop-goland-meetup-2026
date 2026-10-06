package gbif

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteError(t *testing.T) {
	rec := httptest.NewRecorder()
	before := time.Now().UTC().Truncate(time.Millisecond)
	writeError(rec, http.StatusNotFound, "Entity not found for uri: /species/0")
	after := time.Now().UTC()

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, "application/json", rec.Result().Header.Get("Content-Type"))
	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Len(t, body, 4)
	assert.Equal(t, float64(http.StatusNotFound), body["status"])
	assert.Equal(t, "", body["error"])
	assert.Equal(t, "Entity not found for uri: /species/0", body["message"])
	timestamp, ok := body["timestamp"].(string)
	require.True(t, ok)
	assert.Regexp(t, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}\+00:00$`, timestamp)
	parsed, err := time.Parse("2006-01-02T15:04:05.000+00:00", timestamp)
	require.NoError(t, err)
	assert.False(t, parsed.Before(before))
	assert.False(t, parsed.After(after))
}
