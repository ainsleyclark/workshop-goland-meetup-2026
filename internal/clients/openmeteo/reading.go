package openmeteo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"workshop/internal/common/httputil"
)

// hourLayout is the format Open-Meteo expects request hours in.
const hourLayout = "2006-01-02T15:04"

// Request asks for the weather at one place at one moment.
type Request struct {
	Coordinates Coordinates
	Time        time.Time
}

// Validate checks if the request is valid.
func (r Request) Validate() error {
	if r.Time.IsZero() {
		return errors.New("openmeteo: time is required")
	}
	return r.Coordinates.Validate()
}

// query builds the URL query values for the archive endpoint. Asking for
// the same start and end hour returns exactly one reading.
func (r Request) query() url.Values {
	hour := r.Time.UTC().Truncate(time.Hour).Format(hourLayout)

	q := url.Values{}
	q.Set("latitude", strconv.FormatFloat(r.Coordinates.Latitude, 'f', -1, 64))
	q.Set("longitude", strconv.FormatFloat(r.Coordinates.Longitude, 'f', -1, 64))
	q.Set("start_hour", hour)
	q.Set("end_hour", hour)
	q.Set("hourly", strings.Join(hourlyParams, ","))

	return q
}

// Reading is the weather at a single place and hour.
type Reading struct {
	Coordinates         Coordinates
	Time                time.Time
	Temperature         float64
	ApparentTemperature float64
	Humidity            int
	Precipitation       float64
	Rain                float64
	Snowfall            float64
	CloudCover          int
	WindSpeed           float64
	WindGusts           float64
	WindDirection       int
	Code                WeatherCode
}

// Reading retrieves the weather at a given location and time. The
// archive holds hourly readings, so the request is rounded down to the
// hour the moment falls in.
func (c *Client) Reading(ctx context.Context, req Request) (Reading, error) {
	if err := req.Validate(); err != nil {
		return Reading{}, fmt.Errorf("openmeteo: invalid request: %w", err)
	}

	result, err := c.http.Do[hourlyResponse](ctx, httputil.Request{
		Method: http.MethodGet,
		Path:   "/archive",
		Query:  req.query(),
	})
	if err != nil {
		return Reading{}, err
	}

	// Sometimes the Open-Meteo API returns 0 as a longitude for
	// some reason, if we don't assign it a zero value, downstream
	// services may fail.
	if result.Longitude == 0 {
		result.Longitude = req.Coordinates.Longitude
	}

	return result.reading()
}

// reading takes the single hour out of the archive response. Every metric
// arrives as an array with one entry, so each is read at index zero.
func (h hourlyResponse) reading() (Reading, error) {
	if len(h.Hourly.Time) == 0 {
		return Reading{}, errors.New("openmeteo: no hourly data returned")
	}

	at, err := time.Parse(hourLayout, h.Hourly.Time[0])
	if err != nil {
		return Reading{}, fmt.Errorf("openmeteo: parsing time: %w", err)
	}

	return Reading{
		Coordinates:         Coordinates{Latitude: h.Latitude, Longitude: h.Longitude},
		Time:                at,
		Temperature:         first(h.Hourly.Temperature2M),
		ApparentTemperature: first(h.Hourly.ApparentTemperature),
		Humidity:            first(h.Hourly.RelativeHumidity2M),
		Precipitation:       first(h.Hourly.Precipitation),
		Rain:                first(h.Hourly.Rain),
		Snowfall:            first(h.Hourly.Snowfall),
		CloudCover:          first(h.Hourly.CloudCover),
		WindSpeed:           first(h.Hourly.WindSpeed10M),
		WindGusts:           first(h.Hourly.WindGusts10M),
		WindDirection:       first(h.Hourly.WindDirection10M),
		Code:                WeatherCode(first(h.Hourly.WeatherCode)),
	}, nil
}

// first returns the first value of an hourly array, or the zero value if
// Open-Meteo left the variable out of the response.
func first[T any](values []T) T {
	var zero T
	if len(values) == 0 {
		return zero
	}
	return values[0]
}

var hourlyParams = []string{
	"temperature_2m",
	"apparent_temperature",
	"relative_humidity_2m",
	"precipitation",
	"rain",
	"snowfall",
	"cloud_cover",
	"weather_code",
	"wind_speed_10m",
	"wind_gusts_10m",
	"wind_direction_10m",
}

// hourlyResponse is the raw archive API response. Its metrics are arraysl
// one entry per requested hour.
//
// See: https://open-meteo.com/en/docs/historical-weather-api
type (
	hourlyResponse struct {
		Latitude  float64      `json:"latitude"`
		Longitude float64      `json:"longitude"`
		Hourly    hourlyValues `json:"hourly"`
	}
	hourlyValues struct {
		Time                []string  `json:"time"`
		Temperature2M       []float64 `json:"temperature_2m"`
		ApparentTemperature []float64 `json:"apparent_temperature"`
		RelativeHumidity2M  []int     `json:"relative_humidity_2m"`
		Precipitation       []float64 `json:"precipitation"`
		Rain                []float64 `json:"rain"`
		Snowfall            []float64 `json:"snowfall"`
		CloudCover          []int     `json:"cloud_cover"`
		WeatherCode         []int     `json:"weather_code"`
		WindSpeed10M        []float64 `json:"wind_speed_10m"`
		WindGusts10M        []float64 `json:"wind_gusts_10m"`
		WindDirection10M    []int     `json:"wind_direction_10m"`
	}
)
