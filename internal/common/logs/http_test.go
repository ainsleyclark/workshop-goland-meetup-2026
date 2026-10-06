package logs

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLogger(t *testing.T) {
	t.Parallel()

	t.Run("Successful response", func(t *testing.T) {
		t.Parallel()

		log, output := setup(t)
		next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("X-Test", "preserved")
			w.WriteHeader(http.StatusOK)

			_, err := w.Write([]byte("response"))
			require.NoError(t, err)
		})
		handler := RequestLogger(log)(next)
		response := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/sightings?token=secret", nil).WithContext(t.Context())

		handler.ServeHTTP(response, request)

		var got map[string]any
		require.NoError(t, json.Unmarshal(output.Bytes(), &got))

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "response", response.Body.String())
		assert.Equal(t, "preserved", response.Header().Get("X-Test"))
		assert.Equal(t, float64(http.StatusOK), got["status"])
		assert.Equal(t, float64(len("response")), got["bytes"])
		assert.Equal(t, "INFO", got["level"])
		assert.Equal(t, "GET", got["method"])
		assert.Equal(t, "/sightings", got["path"])
		assert.Contains(t, got, "duration")
		assert.NotContains(t, output.String(), "secret")
	})

	t.Run("Client error", func(t *testing.T) {
		t.Parallel()

		log, output := setup(t)
		next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("X-Test", "preserved")
			w.WriteHeader(http.StatusNotFound)

			_, err := w.Write([]byte("response"))
			require.NoError(t, err)
		})
		handler := RequestLogger(log)(next)
		response := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/sightings?token=secret", nil).WithContext(t.Context())

		handler.ServeHTTP(response, request)

		var got map[string]any
		require.NoError(t, json.Unmarshal(output.Bytes(), &got))

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, "response", response.Body.String())
		assert.Equal(t, "preserved", response.Header().Get("X-Test"))
		assert.Equal(t, float64(http.StatusNotFound), got["status"])
		assert.Equal(t, float64(len("response")), got["bytes"])
		assert.Equal(t, "WARN", got["level"])
		assert.Equal(t, "GET", got["method"])
		assert.Equal(t, "/sightings", got["path"])
		assert.Contains(t, got, "duration")
		assert.NotContains(t, output.String(), "secret")
	})

	t.Run("Server error", func(t *testing.T) {
		t.Parallel()

		log, output := setup(t)
		next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("X-Test", "preserved")
			w.WriteHeader(http.StatusInternalServerError)

			_, err := w.Write([]byte("response"))
			require.NoError(t, err)
		})
		handler := RequestLogger(log)(next)
		response := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/sightings?token=secret", nil).WithContext(t.Context())

		handler.ServeHTTP(response, request)

		var got map[string]any
		require.NoError(t, json.Unmarshal(output.Bytes(), &got))

		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, "response", response.Body.String())
		assert.Equal(t, "preserved", response.Header().Get("X-Test"))
		assert.Equal(t, float64(http.StatusInternalServerError), got["status"])
		assert.Equal(t, float64(len("response")), got["bytes"])
		assert.Equal(t, "ERROR", got["level"])
		assert.Equal(t, "GET", got["method"])
		assert.Equal(t, "/sightings", got["path"])
		assert.Contains(t, got, "duration")
		assert.NotContains(t, output.String(), "secret")
	})

	t.Run("Empty response defaults to success", func(t *testing.T) {
		t.Parallel()

		log, output := setup(t)
		handler := RequestLogger(log)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
		request := httptest.NewRequest("GET", "/", nil).WithContext(t.Context())

		handler.ServeHTTP(httptest.NewRecorder(), request)

		var got map[string]any
		require.NoError(t, json.Unmarshal(output.Bytes(), &got))

		assert.Equal(t, float64(http.StatusOK), got["status"])
		assert.Equal(t, float64(0), got["bytes"])
		assert.Equal(t, "INFO", got["level"])
	})

	t.Run("Writing a body implies success", func(t *testing.T) {
		t.Parallel()

		log, output := setup(t)
		next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, err := w.Write([]byte("hello"))
			require.NoError(t, err)
		})
		handler := RequestLogger(log)(next)
		request := httptest.NewRequest("GET", "/", nil).WithContext(t.Context())

		handler.ServeHTTP(httptest.NewRecorder(), request)

		var got map[string]any
		require.NoError(t, json.Unmarshal(output.Bytes(), &got))

		assert.Equal(t, float64(http.StatusOK), got["status"])
		assert.Equal(t, float64(len("hello")), got["bytes"])
	})

	t.Run("Preserves response flushing", func(t *testing.T) {
		t.Parallel()

		log, _ := setup(t)
		next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			require.NoError(t, http.NewResponseController(w).Flush())

			_, err := w.Write([]byte("stream"))
			require.NoError(t, err)
		})
		handler := RequestLogger(log)(next)
		got := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/", nil).WithContext(t.Context())

		handler.ServeHTTP(got, request)

		assert.True(t, got.Flushed)
		assert.Equal(t, "stream", got.Body.String())
	})

	t.Run("Propagates panics and logs interruption", func(t *testing.T) {
		t.Parallel()

		log, output := setup(t)
		next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			panic("broken")
		})
		handler := RequestLogger(log)(next)
		request := httptest.NewRequest("GET", "/", nil).WithContext(t.Context())

		assert.PanicsWithValue(t, "broken", func() {
			handler.ServeHTTP(httptest.NewRecorder(), request)
		})

		var got map[string]any
		require.NoError(t, json.Unmarshal(output.Bytes(), &got))

		assert.Equal(t, "ERROR", got["level"])
		assert.Equal(t, "HTTP request interrupted", got["msg"])
		assert.Equal(t, float64(0), got["status"])
	})
}
