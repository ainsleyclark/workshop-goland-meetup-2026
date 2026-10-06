package weather

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// hourLayout is the format Open-Meteo expects request hours in. It must
// match the unexported constant of the same name in
// internal/clients/openmeteo/reading.go: this mock deliberately doesn't
// import the production client, so the two can't share it directly.
const hourLayout = "2006-01-02T15:04"

// maxRangeHours bounds how large a range a single request can ask for, so a
// malformed request can't make the mock allocate without limit.
const maxRangeHours = 24 * 366

// archiveRequest is a parsed GET /archive request.
type archiveRequest struct {
	latitude, longitude float64
	startHour, endHour  time.Time
}

// parseRequest validates and parses the query parameters GET /archive
// takes.
func parseRequest(q url.Values) (archiveRequest, error) {
	lat, err := parseCoordinate(q, "latitude", -90, 90)
	if err != nil {
		return archiveRequest{}, err
	}

	lon, err := parseCoordinate(q, "longitude", -180, 180)
	if err != nil {
		return archiveRequest{}, err
	}

	start, err := parseHour(q, "start_hour")
	if err != nil {
		return archiveRequest{}, err
	}

	end, err := parseHour(q, "end_hour")
	if err != nil {
		return archiveRequest{}, err
	}

	if end.Before(start) {
		return archiveRequest{}, fmt.Errorf("end_hour must not be before start_hour")
	}
	if end.Sub(start) > maxRangeHours*time.Hour {
		return archiveRequest{}, fmt.Errorf("requested range is too large")
	}

	// The only caller always asks for the same fixed set of variables, so
	// the mock doesn't filter its response by hourly's contents - it just
	// requires the parameter to be present, as the real API does.
	if q.Get("hourly") == "" {
		return archiveRequest{}, fmt.Errorf("hourly is required")
	}

	return archiveRequest{latitude: lat, longitude: lon, startHour: start, endHour: end}, nil
}

// parseCoordinate reads and range-checks a required latitude or longitude
// parameter.
func parseCoordinate(q url.Values, key string, lo, hi float64) (float64, error) {
	raw := q.Get(key)
	if raw == "" {
		return 0, fmt.Errorf("%s is required", key)
	}

	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number", key)
	}
	if v < lo || v > hi {
		return 0, fmt.Errorf("%s must be between %g and %g", key, lo, hi)
	}

	return v, nil
}

// parseHour reads a required start_hour or end_hour parameter.
func parseHour(q url.Values, key string) (time.Time, error) {
	raw := q.Get(key)
	if raw == "" {
		return time.Time{}, fmt.Errorf("%s is required", key)
	}

	t, err := time.Parse(hourLayout, raw)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s must be in the form %s", key, hourLayout)
	}

	return t, nil
}
