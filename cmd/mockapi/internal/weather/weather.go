// Package weather mocks the Open-Meteo archive API, generating a
// deterministic reading for any coordinate and hour.
package weather

import (
	"net/http"
	"time"
	"workshop/internal/common/httputil"
)

// Handler serves the mocked Open-Meteo archive endpoint under /weather/v1.
func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /weather/v1/archive", handleArchive)
	mux.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "Not Found")
	})
	return mux
}

// handleArchive mocks GET /archive.
func handleArchive(w http.ResponseWriter, r *http.Request) {
	req, err := parseRequest(r.URL.Query())
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var values hourlyValues
	for hour := req.startHour; !hour.After(req.endHour); hour = hour.Add(time.Hour) {
		v := generate(req.latitude, req.longitude, hour)

		values.Time = append(values.Time, hour.Format(hourLayout))
		values.Temperature2M = append(values.Temperature2M, v.Temperature)
		values.ApparentTemperature = append(values.ApparentTemperature, v.ApparentTemperature)
		values.RelativeHumidity2M = append(values.RelativeHumidity2M, v.Humidity)
		values.Precipitation = append(values.Precipitation, v.Precipitation)
		values.Rain = append(values.Rain, v.Rain)
		values.Snowfall = append(values.Snowfall, v.Snowfall)
		values.CloudCover = append(values.CloudCover, v.CloudCover)
		values.WeatherCode = append(values.WeatherCode, v.Code)
		values.WindSpeed10M = append(values.WindSpeed10M, v.WindSpeed)
		values.WindGusts10M = append(values.WindGusts10M, v.WindGusts)
		values.WindDirection10M = append(values.WindDirection10M, v.WindDirection)
	}

	httputil.WriteJSON(w, http.StatusOK, hourlyResponse{
		Latitude:  req.latitude,
		Longitude: req.longitude,
		Hourly:    values,
	})
}
