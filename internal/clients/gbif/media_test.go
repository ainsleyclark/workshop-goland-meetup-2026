package gbif

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"workshop/internal/common/httputil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Media(t *testing.T) {
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

		body, err := os.ReadFile("testdata/media.json")
		require.NoError(t, err)

		test := setup(t, http.StatusOK, body)

		res, err := test.client.Media(t.Context(), 5231190)
		require.NoError(t, err)

		t.Log("Request is sent to the species media endpoint with the key in the path")
		{
			assert.Equal(t, "/species/5231190/media", test.gotURL.Path)
		}

		t.Log("Every decoded result belongs to the requested taxon")
		{
			require.Len(t, res.Results, 2)
			for _, result := range res.Results {
				assert.Equal(t, 5231190, result.TaxonKey)
				assert.Equal(t, "StillImage", result.Type)
			}
		}
	})

	t.Run("Non 2xx returns an error", func(t *testing.T) {
		t.Parallel()

		test := setup(t, http.StatusInternalServerError, []byte("boom"))

		_, err := test.client.Media(t.Context(), 5231190)
		assert.Equal(t, true, err != nil)
	})
}
