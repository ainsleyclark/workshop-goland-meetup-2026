package weather

import (
	"net/http"
	"workshop/internal/common/httputil"
)

// hourlyResponse mirrors the archive endpoint's real response shape, field
// for field, matching the unexported struct of the same name in
// internal/clients/openmeteo/reading.go - that one can't be imported here,
// since it's unexported, so keep the two in sync by hand.
type hourlyResponse struct {
	Latitude  float64      `json:"latitude"`
	Longitude float64      `json:"longitude"`
	Hourly    hourlyValues `json:"hourly"`
}

// hourlyValues holds one entry per requested hour.
type hourlyValues struct {
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

// writeError writes an error in the shape Open-Meteo uses: a JSON object
// with an HTTP 400 status for a bad request.
// See: https://open-meteo.com/en/docs/historical-weather-api
func writeError(w http.ResponseWriter, status int, reason string) {
	httputil.WriteJSON(w, status, struct {
		Error  bool   `json:"error"`
		Reason string `json:"reason"`
	}{
		Error:  true,
		Reason: reason,
	})
}
