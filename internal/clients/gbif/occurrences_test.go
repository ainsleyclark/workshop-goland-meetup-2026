package gbif

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
	"workshop/internal/common/httputil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Occurrences(t *testing.T) {
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

	t.Run("Filters by taxa and country", func(t *testing.T) {
		t.Parallel()

		test := setup(t, http.StatusOK, []byte(`{"results":[]}`))
		_, err := test.client.Occurrences(t.Context(), OccurrencesSearchArgs{
			TaxonKeys: []int{5219404, 2435350, 2435349, 5219461},
			Country:   "KE",
		})
		require.NoError(t, err)
		assert.Equal(t, url.Values{
			"taxonKey":     {"5219404", "2435350", "2435349", "5219461"},
			"checklistKey": {"d7dddbf4-2cf0-4f39-9b2a-bb099caae36c"},
			"country":      {"KE"},
		}, test.gotURL.Query())
	})

	t.Run("Paginates filtered occurrences", func(t *testing.T) {
		t.Parallel()

		test := setup(t, http.StatusOK, []byte(`{"offset":40,"limit":20,"endOfRecords":true,"count":41,"results":[{"key":123}]}`))
		res, err := test.client.Occurrences(t.Context(), OccurrencesSearchArgs{
			Limit:   20,
			Offset:  40,
			Country: "KE",
		})

		require.NoError(t, err)
		assert.Equal(t, url.Values{
			"country": {"KE"},
			"limit":   {"20"},
			"offset":  {"40"},
		}, test.gotURL.Query())
		assert.Equal(t, 40, res.Offset)
		assert.Equal(t, 20, res.Limit)
		assert.True(t, res.EndOfRecords)
		assert.Equal(t, int64(41), res.Count)
		require.Len(t, res.Results, 1)
		assert.Equal(t, int64(123), res.Results[0].Key)
	})

	t.Run("Filters by country", func(t *testing.T) {
		t.Parallel()

		body, err := os.ReadFile("testdata/occurrences_country.json")
		require.NoError(t, err)

		test := setup(t, http.StatusOK, body)

		res, err := test.client.Occurrences(t.Context(), OccurrencesSearchArgs{Country: "KE"})
		require.NoError(t, err)

		t.Log("Request is sent to the search endpoint with the query encoded")
		{
			assert.Equal(t, "/occurrence/search", test.gotURL.Path)
			assert.Equal(t, "country=KE", test.gotURL.RawQuery)
		}

		t.Log("Every decoded result actually belongs to the requested country")
		{
			require.Len(t, res.Results, 2)
			for _, result := range res.Results {
				assert.Equal(t, "KE", result.CountryCode)
			}
		}
	})

	t.Run("Decodes date-only creation metadata in multimedia extensions", func(t *testing.T) {
		t.Parallel()

		// Captured occurrence 5933004352 has a raw date-only creation value,
		// while GBIF's normalized media.created contains a full timestamp.
		body, err := os.ReadFile("testdata/occurrences_date_only_created.json")
		require.NoError(t, err)
		test := setup(t, http.StatusOK, body)

		res, err := test.client.Occurrences(t.Context(), OccurrencesSearchArgs{})
		require.NoError(t, err)
		require.Len(t, res.Results, 1)
		result := res.Results[0]
		assert.Equal(t, int64(5933004352), result.Key)
		require.Len(t, result.Extensions.HttpRsGbifOrgTerms10Multimedia, 1)
		assert.Equal(t, "2026-01-06", result.Extensions.HttpRsGbifOrgTerms10Multimedia[0].HttpPurlOrgDcTermsCreated)
		require.Len(t, result.Media, 1)
		assert.True(t, time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC).Equal(result.Media[0].Created))
		assert.True(t, time.Date(2026, 1, 6, 11, 0, 0, 0, time.UTC).Equal(result.EventDate))
	})

	t.Run("Filters by latitude and longitude", func(t *testing.T) {
		t.Parallel()

		body, err := os.ReadFile("testdata/occurrences_latlong.json")
		require.NoError(t, err)

		test := setup(t, http.StatusOK, body)

		res, err := test.client.Occurrences(t.Context(), OccurrencesSearchArgs{Latitude: -3.4306, Longitude: 39.964508})
		require.NoError(t, err)

		t.Log("Request is sent to the search endpoint with the query encoded")
		{
			assert.Equal(t, "/occurrence/search", test.gotURL.Path)
			assert.Equal(t, "decimalLatitude=-3.4306&decimalLongitude=39.964508", test.gotURL.RawQuery)
		}

		t.Log("The decoded result sits at the exact coordinates requested")
		{
			require.Len(t, res.Results, 1)
			assert.Equal(t, -3.4306, res.Results[0].DecimalLatitude)
			assert.Equal(t, 39.964508, res.Results[0].DecimalLongitude)
		}
	})

	t.Run("Filters by date range", func(t *testing.T) {
		t.Parallel()

		body, err := os.ReadFile("testdata/occurrences_daterange.json")
		require.NoError(t, err)

		test := setup(t, http.StatusOK, body)

		res, err := test.client.Occurrences(t.Context(), OccurrencesSearchArgs{
			DateFrom: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			DateTo:   time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC),
		})
		require.NoError(t, err)

		t.Log("Request is sent to the search endpoint with the query encoded")
		{
			assert.Equal(t, "/occurrence/search", test.gotURL.Path)
			assert.Equal(t, "eventDate=2020-01-01%2C2020-01-31", test.gotURL.RawQuery)
		}

		t.Log("Every decoded result falls inside the requested date range")
		{
			require.Len(t, res.Results, 2)
			for _, result := range res.Results {
				assert.Equal(t, 2020, result.Year)
				assert.Equal(t, 1, result.Month)
			}
		}
	})

	t.Run("Invalid latitude returns an error", func(t *testing.T) {
		t.Parallel()

		// Captured from a live call to /occurrence/search?decimalLatitude=999.
		const invalidLatitudeBody = "999 is not valid value, latitude must be between -90 and 90."

		test := setup(t, http.StatusBadRequest, []byte(invalidLatitudeBody))

		_, err := test.client.Occurrences(t.Context(), OccurrencesSearchArgs{Latitude: 999})
		require.Error(t, err)

		t.Log("The error carries the GBIF 400 status and body verbatim")
		{
			var httpErr *httputil.Error
			require.True(t, errors.As(err, &httpErr))
			assert.Equal(t, http.StatusBadRequest, httpErr.StatusCode)
			assert.Contains(t, httpErr.Body, "latitude must be between -90 and 90")
		}
	})

	t.Run("Decodes real taxon search responses", func(t *testing.T) {
		t.Parallel()

		type Input struct {
			taxonKeys []int
			fixture   string
		}

		tt := map[string]struct {
			input       Input
			wantKeys    []string
			wantSpecies []int
		}{
			"Humpback whale including a subspecies": {
				input:       Input{[]int{5220086}, "humpback_whale"},
				wantKeys:    []string{"5220086"},
				wantSpecies: []int{5220086, 5220086},
			},
			"Common bottlenose dolphin": {
				input:       Input{[]int{2440447}, "bottlenose_dolphin"},
				wantKeys:    []string{"2440447"},
				wantSpecies: []int{2440447, 2440447},
			},
			"Alligator group overlapping with a species": {
				input:       Input{[]int{2441370, 2441368, 2441370}, "alligators"},
				wantKeys:    []string{"2441370", "2441368"},
				wantSpecies: []int{2441370, 2441370},
			},
			"Multiple whales and a dolphin": {
				input:       Input{[]int{2440735, 5220086, 2440447}, "marine_animals"},
				wantKeys:    []string{"2440735", "5220086", "2440447"},
				wantSpecies: []int{2440447, 2440447},
			},
		}

		for name, test := range tt {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				body, err := os.ReadFile("testdata/occurrences_" + test.input.fixture + ".json")
				require.NoError(t, err)

				client := setup(t, http.StatusOK, body)
				res, err := client.client.Occurrences(t.Context(), OccurrencesSearchArgs{TaxonKeys: test.input.taxonKeys})
				require.NoError(t, err)

				assert.Equal(t, "/occurrence/search", client.gotURL.Path)
				assert.Equal(t, url.Values{
					"taxonKey":     test.wantKeys,
					"checklistKey": {"d7dddbf4-2cf0-4f39-9b2a-bb099caae36c"},
				}, client.gotURL.Query())
				require.Len(t, res.Results, len(test.wantSpecies))
				assert.Equal(t, 2, res.Limit)
				assert.Positive(t, res.Count)

				for i, result := range res.Results {
					assert.Positive(t, result.Key)
					assert.Equal(t, test.wantSpecies[i], result.SpeciesKey)
					assert.NotEmpty(t, result.ScientificName)
				}
			})
		}
	})
}

func TestOccurrencesSearchArgs_Query(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input OccurrencesSearchArgs
		want  string
	}{
		"Empty": {
			input: OccurrencesSearchArgs{},
			want:  "",
		},
		"First page with custom limit": {
			input: OccurrencesSearchArgs{PaginationOptions: PaginationOptions{Limit: 300}},
			want:  "limit=300",
		},
		"Offset with default limit": {
			input: OccurrencesSearchArgs{PaginationOptions: PaginationOptions{Offset: 40}},
			want:  "offset=40",
		},
		"Country only": {
			input: OccurrencesSearchArgs{Country: "KE"},
			want:  "country=KE",
		},
		"Single taxon": {
			input: OccurrencesSearchArgs{TaxonKeys: []int{5219404}},
			want:  "checklistKey=d7dddbf4-2cf0-4f39-9b2a-bb099caae36c&taxonKey=5219404",
		},
		"Several taxa": {
			input: OccurrencesSearchArgs{TaxonKeys: []int{2440892, 2440888, 2440894}},
			want:  "checklistKey=d7dddbf4-2cf0-4f39-9b2a-bb099caae36c&taxonKey=2440892&taxonKey=2440888&taxonKey=2440894",
		},
		"Repeated taxa are deduplicated": {
			input: OccurrencesSearchArgs{TaxonKeys: []int{5219404, 5219404}},
			want:  "checklistKey=d7dddbf4-2cf0-4f39-9b2a-bb099caae36c&taxonKey=5219404",
		},
		"Radius takes precedence over lat long": {
			input: OccurrencesSearchArgs{
				Latitude:    1.5,
				Longitude:   2.5,
				RadiusMiles: 10,
			},
			want: "geo_distance=10mi%2C1.5%2C2.5",
		},
		"Latitude only": {
			input: OccurrencesSearchArgs{Latitude: 1.5},
			want:  "decimalLatitude=1.5",
		},
		"Longitude only": {
			input: OccurrencesSearchArgs{Longitude: 2.5},
			want:  "decimalLongitude=2.5",
		},
		"Latitude and longitude": {
			input: OccurrencesSearchArgs{Latitude: 1.5, Longitude: 2.5},
			want:  "decimalLatitude=1.5&decimalLongitude=2.5",
		},
		"Date from only": {
			input: OccurrencesSearchArgs{DateFrom: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			want:  "eventDate=2020-01-01",
		},
		"Date to only": {
			input: OccurrencesSearchArgs{DateTo: time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)},
			want:  "eventDate=2020-12-31",
		},
		"Date range": {
			input: OccurrencesSearchArgs{
				DateFrom: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				DateTo:   time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC),
			},
			want: "eventDate=2020-01-01%2C2020-12-31",
		},
		"All fields": {
			input: OccurrencesSearchArgs{
				Country:  "KE",
				Latitude: 1.5,
				DateFrom: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: "country=KE&decimalLatitude=1.5&eventDate=2020-01-01",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.input.query()
			assert.Equal(t, test.want, got.Encode())
		})
	}
}

func TestOccurrence_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input string
		want  string
	}{
		"UTC timestamp":       {`"2026-01-03T19:52:53Z"`, "2026-01-03T19:52:53Z"},
		"Offset timestamp":    {`"2026-01-03T19:52:53-05:00"`, "2026-01-03T19:52:53-05:00"},
		"Fractional seconds":  {`"2026-01-03T19:52:53.123456789Z"`, "2026-01-03T19:52:53.123456789Z"},
		"Local seconds":       {`"2026-01-01T10:18:12"`, "2026-01-01T10:18:12Z"},
		"Local fractions":     {`"2026-01-01T10:18:12.123"`, "2026-01-01T10:18:12.123Z"},
		"Local minutes":       {`"2026-01-01T05:55"`, "2026-01-01T05:55:00Z"},
		"Minutes with offset": {`"2026-01-01T05:55-05:00"`, "2026-01-01T05:55:00-05:00"},
		"Date only":           {`"2026-01-20"`, "2026-01-20T00:00:00Z"},
		"Month only":          {`"2020-01"`, "2020-01-01T00:00:00Z"},
		"Year only":           {`"2020"`, "2020-01-01T00:00:00Z"},
		"Date interval":       {`"2026-01-20/2026-01-23"`, "2026-01-20T00:00:00Z"},
		"Month interval":      {`"2020-01/2020-03"`, "2020-01-01T00:00:00Z"},
		"Timestamp interval":  {`"2026-01-03T19:52:53Z/2026-01-03T20:00:00Z"`, "2026-01-03T19:52:53Z"},
		"Empty":               {`""`, "0001-01-01T00:00:00Z"},
		"Null":                {`null`, "0001-01-01T00:00:00Z"},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var occurrence Occurrence
			err := json.Unmarshal([]byte(`{"key":5277301953,"eventDate":`+test.input+`,"speciesKey":2441370,"countryCode":"US"}`), &occurrence)
			require.NoError(t, err)
			assert.Equal(t, test.want, occurrence.EventDate.Format(time.RFC3339Nano))
			assert.Equal(t, int64(5277301953), occurrence.Key)
			assert.Equal(t, 2441370, occurrence.SpeciesKey)
			assert.Equal(t, "US", occurrence.CountryCode)

			encoded, err := json.Marshal(occurrence)
			require.NoError(t, err)
			var roundTrip Occurrence
			require.NoError(t, json.Unmarshal(encoded, &roundTrip))
			assert.True(t, occurrence.EventDate.Equal(roundTrip.EventDate))
		})
	}

	t.Run("Missing date", func(t *testing.T) {
		t.Parallel()

		var occurrence Occurrence
		require.NoError(t, json.Unmarshal([]byte(`{"key":123}`), &occurrence))
		assert.True(t, occurrence.EventDate.IsZero())
		assert.Equal(t, int64(123), occurrence.Key)
	})

	t.Run("Missing and null fields preserve existing values", func(t *testing.T) {
		t.Parallel()

		date := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)
		occurrence := Occurrence{Key: 123, EventDate: date}
		for _, input := range []string{`{}`, `{"eventDate":null}`} {
			require.NoError(t, json.Unmarshal([]byte(input), &occurrence))
			assert.Equal(t, date, occurrence.EventDate)
			assert.Equal(t, int64(123), occurrence.Key)
		}
	})

	for _, value := range []string{
		`"invalid"`, `"2026-02-30"`, `"2026-13"`, `"2026-01-20/invalid"`,
		`"2026-01-23/2026-01-20"`, `"2026/2027/2028"`, `"/2026"`, `"2026/"`, `123`,
	} {
		t.Run("Invalid date "+value, func(t *testing.T) {
			t.Parallel()

			occurrence := Occurrence{Key: 123}
			err := json.Unmarshal([]byte(`{"key":456,"eventDate":`+value+`}`), &occurrence)
			require.ErrorContains(t, err, "eventDate")
			assert.Equal(t, int64(123), occurrence.Key)
		})
	}
}
