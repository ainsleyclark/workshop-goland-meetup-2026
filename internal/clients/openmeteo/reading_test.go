package openmeteo

import (
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

type test struct {
	client *Client
	gotURL *url.URL
}

func setup(t *testing.T, status int, body []byte) test {
	t.Helper()

	var gotURL url.URL
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotURL = *r.URL
		w.WriteHeader(status)
		_, err := w.Write(body)
		require.NoError(t, err)
	}))
	t.Cleanup(srv.Close)

	return test{
		client: &Client{http: httputil.New(httputil.Config{BaseURL: srv.URL})},
		gotURL: &gotURL,
	}
}

var berlin = Coordinates{Latitude: 52.54833, Longitude: 13.407822}

// validRequest passes validation, so the client always reaches the stub.
var validRequest = Request{
	Coordinates: berlin,
	Time:        time.Date(2023, 9, 7, 10, 0, 0, 0, time.UTC),
}

func TestClient_Reading(t *testing.T) {
	t.Parallel()

	t.Run("Validation Error", func(t *testing.T) {
		t.Parallel()

		_, err := New("").Reading(t.Context(), Request{
			Coordinates: Coordinates{Latitude: 100, Longitude: 200},
			Time:        time.Date(2023, 9, 7, 10, 0, 0, 0, time.UTC),
		})
		require.Error(t, err)
	})

	t.Run("Invalid Request", func(t *testing.T) {
		t.Parallel()

		test := setup(t, http.StatusBadRequest, nil)

		_, err := test.client.Reading(t.Context(), Request{
			Coordinates: berlin,
			Time:        time.Date(2023, 9, 7, 10, 0, 0, 0, time.UTC),
		})
		require.Error(t, err)
	})

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		body, err := os.ReadFile("testdata/hourly-response.json")
		require.NoError(t, err)

		test := setup(t, http.StatusOK, body)

		// A moment part-way through the hour, to prove it is rounded down.
		got, err := test.client.Reading(t.Context(), Request{
			Coordinates: berlin,
			Time:        time.Date(2023, 9, 7, 10, 47, 31, 0, time.UTC),
		})
		require.NoError(t, err)

		t.Log("The request asks the archive for a single hour")
		{
			q := test.gotURL.Query()
			assert.Equal(t, "/archive", test.gotURL.Path)
			assert.Equal(t, "2023-09-07T10:00", q.Get("start_hour"))
			assert.Equal(t, q.Get("start_hour"), q.Get("end_hour"))
			assert.Equal(t, "52.54833", q.Get("latitude"))
			assert.Equal(t, "13.407822", q.Get("longitude"))
			assert.Empty(t, q.Get("start_date"))
			assert.Empty(t, q.Get("daily"))
		}

		t.Log("Every hourly variable the Reading needs is requested")
		{
			for _, param := range hourlyParams {
				assert.Contains(t, test.gotURL.Query().Get("hourly"), param)
			}
		}

		t.Log("The decoded reading matches the response")
		{
			assert.Equal(t, Reading{
				Coordinates:         berlin,
				Time:                time.Date(2023, 9, 7, 10, 0, 0, 0, time.UTC),
				Temperature:         26.3,
				ApparentTemperature: 24.7,
				Humidity:            41,
				Precipitation:       0,
				Rain:                0,
				Snowfall:            0,
				CloudCover:          13,
				WindSpeed:           17.5,
				WindGusts:           34.6,
				WindDirection:       109,
				Code:                ClearSky,
			}, got)
		}
	})

	t.Run("Zero Long", func(t *testing.T) {
		t.Parallel()

		body, err := os.ReadFile("testdata/hourly-response-zero-long.json")
		require.NoError(t, err)

		test := setup(t, http.StatusOK, body)

		got, err := test.client.Reading(t.Context(), Request{
			Coordinates: Coordinates{Latitude: 54.544587, Longitude: 10.227487},
			Time:        time.Date(2023, 9, 7, 10, 0, 0, 0, time.UTC),
		})
		require.NoError(t, err)

		t.Log("The requested longitude stands in for the zero the API returned")
		{
			assert.Equal(t, 52.54833, got.Coordinates.Latitude)
			assert.Equal(t, 10.227487, got.Coordinates.Longitude)
		}
	})
}

func TestRequest_Validate(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input   Request
		wantErr bool
	}{
		"Valid": {
			input:   Request{Coordinates: berlin, Time: time.Date(2023, 9, 7, 10, 0, 0, 0, time.UTC)},
			wantErr: false,
		},
		"Zero time": {
			input:   Request{Coordinates: berlin},
			wantErr: true,
		},
		"Invalid coordinates": {
			input: Request{
				Coordinates: Coordinates{Latitude: 100, Longitude: 200},
				Time:        time.Date(2023, 9, 7, 10, 0, 0, 0, time.UTC),
			},
			wantErr: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.wantErr, test.input.Validate() != nil)
		})
	}
}

func TestRequest_Query(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input time.Time
		want  string
	}{
		"On the hour":     {time.Date(2023, 9, 7, 10, 0, 0, 0, time.UTC), "2023-09-07T10:00"},
		"Part way though": {time.Date(2023, 9, 7, 10, 59, 59, 0, time.UTC), "2023-09-07T10:00"},
		"Non-UTC zone":    {time.Date(2023, 9, 7, 12, 30, 0, 0, time.FixedZone("CEST", 2*60*60)), "2023-09-07T10:00"},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			q := Request{Coordinates: berlin, Time: test.input}.query()
			assert.Equal(t, test.want, q.Get("start_hour"))
			assert.Equal(t, test.want, q.Get("end_hour"))
		})
	}
}

func TestClient_Reading_Response(t *testing.T) {
	t.Parallel()

	t.Run("No data", func(t *testing.T) {
		t.Parallel()

		test := setup(t, http.StatusOK, []byte(`{"hourly":{"time":[]}}`))

		_, err := test.client.Reading(t.Context(), validRequest)
		assert.ErrorContains(t, err, "no hourly data returned")
	})

	t.Run("Unparsable time", func(t *testing.T) {
		t.Parallel()

		test := setup(t, http.StatusOK, []byte(`{"hourly":{"time":["wrong"]}}`))

		_, err := test.client.Reading(t.Context(), validRequest)
		assert.ErrorContains(t, err, "parsing time")
	})

	t.Run("Missing variables", func(t *testing.T) {
		t.Parallel()

		// Open-Meteo can leave a variable out entirely; the reading
		// should fall back to zero rather than panic.
		test := setup(t, http.StatusOK, []byte(`{"hourly":{"time":["2023-09-07T10:00"]}}`))

		got, err := test.client.Reading(t.Context(), validRequest)
		require.NoError(t, err)
		assert.Equal(t, time.Date(2023, 9, 7, 10, 0, 0, 0, time.UTC), got.Time)
		assert.Zero(t, got.Temperature)
		assert.Zero(t, got.WindDirection)
	})
}
