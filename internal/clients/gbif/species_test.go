package gbif

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"workshop/internal/common/httputil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Species(t *testing.T) {
	t.Parallel()

	type Test struct {
		client *Client
		gotURL *url.URL
	}

	setup := func(t *testing.T, status int, body []byte) Test {
		t.Helper()

		var gotURL url.URL
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotURL = *r.URL
			w.WriteHeader(status)
			_, err := w.Write(body)
			require.NoError(t, err)
		}))
		t.Cleanup(srv.Close)

		return Test{
			client: &Client{http: httputil.New(httputil.Config{BaseURL: srv.URL})},
			gotURL: &gotURL,
		}
	}

	t.Run("Fetches by key", func(t *testing.T) {
		t.Parallel()

		body, err := os.ReadFile("testdata/species.json")
		require.NoError(t, err)

		test := setup(t, http.StatusOK, body)

		res, err := test.client.Species(t.Context(), 2476674)
		require.NoError(t, err)

		t.Log("Request is sent to the species endpoint with the key in the path")
		{
			assert.Equal(t, "/species/2476674", test.gotURL.Path)
		}

		t.Log("The decoded result matches the requested species")
		{
			assert.Equal(t, 2476674, res.Key)
			assert.Equal(t, "Calypte anna (R.Lesson, 1829)", res.ScientificName)
			assert.Equal(t, "SPECIES", res.Rank)
		}
	})

	t.Run("Key not found returns an error", func(t *testing.T) {
		t.Parallel()

		body, err := os.ReadFile("testdata/species_not_found.json")
		require.NoError(t, err)

		test := setup(t, http.StatusNotFound, body)

		_, err = test.client.Species(t.Context(), 2147483647)
		require.Error(t, err)

		t.Log("The error carries the GBIF 404 status and body verbatim")
		{
			var httpErr *httputil.Error
			require.True(t, errors.As(err, &httpErr))
			assert.Equal(t, http.StatusNotFound, httpErr.StatusCode)
			assert.Contains(t, httpErr.Body, "Entity not found for uri: /species/2147483647")
		}
	})
}
