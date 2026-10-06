package gbif

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testdata holds real GBIF responses recorded for the starter's client.
const testdata = "../../../../internal/clients/gbif/testdata"

// fixtureDataset builds a dataset from the GBIF client's recorded
// responses, the same way a scrape does.
func fixtureDataset(t *testing.T) *dataset {
	t.Helper()

	d := newDataset()
	pages, err := filepath.Glob(filepath.Join(testdata, "occurrences_*.json"))
	require.NoError(t, err)
	for _, page := range pages {
		body, err := os.ReadFile(page)
		require.NoError(t, err)
		_, err = d.addPage(body)
		require.NoError(t, err)
	}

	sp, err := os.ReadFile(filepath.Join(testdata, "species.json"))
	require.NoError(t, err)
	d.Species[2476674] = sp

	media, err := os.ReadFile(filepath.Join(testdata, "media.json"))
	require.NoError(t, err)
	d.Media[2476674] = media

	return d
}

func setup(t *testing.T) http.Handler {
	t.Helper()

	var buf bytes.Buffer
	require.NoError(t, fixtureDataset(t).writeSQL(&buf))

	st, err := open(t.Context(), buf.Bytes())
	require.NoError(t, err)
	t.Cleanup(func() { _ = st.Close() })

	return st.Handler()
}

func get(t *testing.T, h http.Handler, target string) (int, []byte) {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	body, err := io.ReadAll(rec.Body)
	require.NoError(t, err)
	return rec.Code, body
}

type searchPage struct {
	Offset       int               `json:"offset"`
	Limit        int               `json:"limit"`
	EndOfRecords bool              `json:"endOfRecords"`
	Count        int               `json:"count"`
	Results      []json.RawMessage `json:"results"`
	Facets       []any             `json:"facets"`
}

func keys(t *testing.T, p searchPage) []int64 {
	t.Helper()
	out := make([]int64, 0, len(p.Results))
	for _, r := range p.Results {
		var v struct {
			Key int64 `json:"key"`
		}
		require.NoError(t, json.Unmarshal(r, &v))
		out = append(out, v.Key)
	}
	return out
}

func TestSearch(t *testing.T) {
	t.Parallel()

	h := setup(t)

	tt := map[string]struct {
		query string
		want  []int64
	}{
		"Species": {
			query: "taxonKey=5220086",
			want:  []int64{5938073638, 5938076605},
		},
		"Subspecies matches its species": {
			query: "taxonKey=7388406",
			want:  []int64{5938076605},
		},
		"Several taxa, as the starter sends them": {
			query: "taxonKey=2441370&taxonKey=2441368&checklistKey=" + backbone,
			want:  []int64{5938029567, 5938030576},
		},
		"Higher rank": {
			query: "taxonKey=212",
			want:  []int64{5938027647, 5933004352, 1039645472},
		},
		"Country": {
			query: "country=ke",
			want:  []int64{5938027617, 5938027647},
		},
		"Country and taxon": {
			query: "country=MX&taxonKey=5220086",
			want:  []int64{5938073638},
		},
		"Exact coordinates": {
			query: "decimalLatitude=-3.4306&decimalLongitude=39.964508",
			want:  []int64{5938027617},
		},
		"Geo distance, radius first as the starter sends it": {
			query: "geo_distance=10mi,-3.4,39.95",
			want:  []int64{5938027617},
		},
		"Geo distance, GBIF order": {
			query: "geo_distance=-3.4,39.95,10km",
			want:  []int64{5938027617},
		},
		"Date range": {
			query: "eventDate=2026-01-04,2026-01-06",
			want:  []int64{5933004352, 5938073638, 5938076605},
		},
		"Partial date": {
			query: "eventDate=2020-01",
			want:  []int64{1039645472, 1089042729},
		},
		"Open ended date": {
			query: "eventDate=*,2021",
			want:  []int64{1039645472, 1089042729},
		},
		"Nothing matches": {
			query: "taxonKey=1&country=AQ&eventDate=1900",
			want:  []int64{},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			status, body := get(t, h, "/gbif/v1/occurrence/search?"+test.query)
			require.Equal(t, http.StatusOK, status, string(body))

			var p searchPage
			require.NoError(t, json.Unmarshal(body, &p))
			assert.ElementsMatch(t, test.want, keys(t, p))
			assert.Equal(t, len(test.want), p.Count)
			assert.True(t, p.EndOfRecords)
			assert.NotNil(t, p.Facets)
		})
	}
}

func TestSearch_Paging(t *testing.T) {
	t.Parallel()

	h := setup(t)

	t.Run("Pages through every record once", func(t *testing.T) {
		t.Parallel()

		var seen []int64
		for offset := 0; ; offset += 4 {
			_, body := get(t, h, "/gbif/v1/occurrence/search?limit=4&offset="+strconv.Itoa(offset))
			var p searchPage
			require.NoError(t, json.Unmarshal(body, &p))
			assert.Equal(t, 11, p.Count)
			seen = append(seen, keys(t, p)...)
			if p.EndOfRecords {
				break
			}
		}
		assert.Len(t, seen, 11)
	})

	t.Run("Defaults and clamps the limit", func(t *testing.T) {
		t.Parallel()

		for query, want := range map[string]int{"": defaultLimit, "limit=1000": maxLimit} {
			_, body := get(t, h, "/gbif/v1/occurrence/search?"+query)
			var p searchPage
			require.NoError(t, json.Unmarshal(body, &p))
			assert.Equal(t, want, p.Limit, query)
		}
	})

	t.Run("Rejects bad parameters", func(t *testing.T) {
		t.Parallel()

		for _, query := range []string{
			"limit=-1",
			"offset=abc",
			"offset=99900&limit=300",
			"taxonKey=whale",
			"geo_distance=1,2",
			"eventDate=yesterday",
		} {
			status, _ := get(t, h, "/gbif/v1/occurrence/search?"+query)
			assert.Equal(t, http.StatusBadRequest, status, query)
		}
	})
}

func TestSearch_ResultsAreUntouched(t *testing.T) {
	t.Parallel()

	h := setup(t)

	raw, err := os.ReadFile(filepath.Join(testdata, "occurrences_humpback_whale.json"))
	require.NoError(t, err)
	var fixture searchPage
	require.NoError(t, json.Unmarshal(raw, &fixture))

	_, body := get(t, h, "/gbif/v1/occurrence/search?taxonKey=5220086")
	var p searchPage
	require.NoError(t, json.Unmarshal(body, &p))

	require.Len(t, p.Results, len(fixture.Results))
	for i := range fixture.Results {
		var want bytes.Buffer
		require.NoError(t, json.Compact(&want, fixture.Results[i]))
		assert.Equal(t, want.String(), string(p.Results[i]))
	}
}

func TestSpecies(t *testing.T) {
	t.Parallel()

	h := setup(t)

	t.Run("Found", func(t *testing.T) {
		t.Parallel()
		status, body := get(t, h, "/gbif/v1/species/2476674")
		assert.Equal(t, http.StatusOK, status)
		assert.Contains(t, string(body), `"canonicalName":"Calypte anna"`)
	})

	t.Run("Not found", func(t *testing.T) {
		t.Parallel()
		status, body := get(t, h, "/gbif/v1/species/2147483647")
		assert.Equal(t, http.StatusNotFound, status)
		assert.Contains(t, string(body), `"message":"Entity not found for uri: /species/2147483647"`)
	})

	t.Run("Media", func(t *testing.T) {
		t.Parallel()
		status, body := get(t, h, "/gbif/v1/species/2476674/media")
		assert.Equal(t, http.StatusOK, status)
		assert.Contains(t, string(body), `"rightsHolder":"Michiel Oversteegen"`)
	})

	t.Run("Media for an unknown species is an empty page", func(t *testing.T) {
		t.Parallel()
		status, body := get(t, h, "/gbif/v1/species/1/media")
		assert.Equal(t, http.StatusOK, status)
		assert.JSONEq(t, `{"offset":0,"limit":20,"endOfRecords":true,"results":[]}`, string(body))
	})

	t.Run("Unknown endpoint", func(t *testing.T) {
		t.Parallel()
		status, _ := get(t, h, "/gbif/v1/dataset/search")
		assert.Equal(t, http.StatusNotFound, status)
	})
}

func TestEmbeddedDumpLoads(t *testing.T) {
	t.Parallel()

	st, err := Open(t.Context())
	require.NoError(t, err)
	t.Cleanup(func() { _ = st.Close() })

	occurrences, _, err := st.Stats(t.Context())
	require.NoError(t, err)
	assert.Positive(t, occurrences)
}
