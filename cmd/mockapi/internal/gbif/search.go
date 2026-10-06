package gbif

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// GBIF's paging rules for occurrence search.
// See: https://techdocs.gbif.org/en/openapi/#paging
const (
	defaultLimit = 20
	maxLimit     = 300
	maxOffset    = 100_000
)

// searchWhere filters occurrences. Every filter is optional: a NULL
// parameter switches it off, so the query text never changes.
const searchWhere = ` WHERE
	(:taxa IS NULL OR key IN (
		SELECT occurrence_key FROM occurrence_taxa WHERE taxon_key IN (SELECT value FROM json_each(:taxa))
	))
	AND (:countries IS NULL OR country IN (SELECT value FROM json_each(:countries)))
	AND (:lat_min IS NULL OR lat BETWEEN :lat_min AND :lat_max)
	AND (:lng_min IS NULL OR lng BETWEEN :lng_min AND :lng_max)
	AND (:geo_km IS NULL OR distance_km(lat, lng, :geo_lat, :geo_lng) <= :geo_km)
	AND (:date_from IS NULL OR event_date BETWEEN :date_from AND :date_to)`

const (
	searchCountQuery = "SELECT count(*) FROM occurrences" + searchWhere
	searchQuery      = "SELECT raw FROM occurrences" + searchWhere + " ORDER BY rowid LIMIT :limit OFFSET :offset"
)

// search holds the parameters for searchWhere. Nil means no filter.
type search struct {
	taxa, countries  *string // JSON arrays.
	latMin, latMax   *float64
	lngMin, lngMax   *float64
	geoLat, geoLng   *float64
	geoKm            *float64
	dateFrom, dateTo *string
	limit, offset    int
}

func (s search) args() []any {
	return []any{
		sql.Named("taxa", s.taxa),
		sql.Named("countries", s.countries),
		sql.Named("lat_min", s.latMin),
		sql.Named("lat_max", s.latMax),
		sql.Named("lng_min", s.lngMin),
		sql.Named("lng_max", s.lngMax),
		sql.Named("geo_lat", s.geoLat),
		sql.Named("geo_lng", s.geoLng),
		sql.Named("geo_km", s.geoKm),
		sql.Named("date_from", s.dateFrom),
		sql.Named("date_to", s.dateTo),
		sql.Named("limit", s.limit),
		sql.Named("offset", s.offset),
	}
}

// parseSearch supports the occurrence search parameters the workshop uses.
// Anything else, such as checklistKey, is accepted and ignored.
func parseSearch(q url.Values) (search, error) {
	s := search{limit: defaultLimit}

	var err error
	if v := q.Get("limit"); v != "" {
		if s.limit, err = strconv.Atoi(v); err != nil || s.limit < 0 {
			return s, fmt.Errorf("invalid limit: %s", v)
		}
		s.limit = min(s.limit, maxLimit) // GBIF clamps rather than rejects.
	}
	if v := q.Get("offset"); v != "" {
		if s.offset, err = strconv.Atoi(v); err != nil || s.offset < 0 {
			return s, fmt.Errorf("invalid offset: %s", v)
		}
	}
	if s.offset+s.limit > maxOffset {
		return s, fmt.Errorf("max offset of %d exceeded: %d + %d", maxOffset, s.offset, s.limit)
	}

	if keys := q["taxonKey"]; len(keys) > 0 {
		taxa := make([]int, len(keys))
		for i, k := range keys {
			if taxa[i], err = strconv.Atoi(k); err != nil {
				return s, fmt.Errorf("invalid taxonKey: %s", k)
			}
		}
		s.taxa = jsonArray(taxa)
	}

	if countries := q["country"]; len(countries) > 0 {
		upper := make([]string, len(countries))
		for i, c := range countries {
			upper[i] = strings.ToUpper(c)
		}
		s.countries = jsonArray(upper)
	}

	if v := q.Get("decimalLatitude"); v != "" {
		lo, hi, err := parseRange(v)
		if err != nil {
			return s, fmt.Errorf("invalid decimalLatitude: %s", v)
		}
		s.latMin, s.latMax = &lo, &hi
	}
	if v := q.Get("decimalLongitude"); v != "" {
		lo, hi, err := parseRange(v)
		if err != nil {
			return s, fmt.Errorf("invalid decimalLongitude: %s", v)
		}
		s.lngMin, s.lngMax = &lo, &hi
	}

	if v := q.Get("geo_distance"); v != "" {
		lat, lng, km, err := parseGeoDistance(v)
		if err != nil {
			return s, fmt.Errorf("invalid geo_distance: %s", v)
		}
		s.geoLat, s.geoLng, s.geoKm = &lat, &lng, &km
	}

	if v := q.Get("eventDate"); v != "" {
		from, to, err := parseDateRange(v)
		if err != nil {
			return s, fmt.Errorf("invalid eventDate: %s", v)
		}
		s.dateFrom, s.dateTo = &from, &to
	}

	return s, nil
}

func jsonArray[T any](values []T) *string {
	b, _ := json.Marshal(values) // Slices of ints and strings always marshal.
	s := string(b)
	return &s
}

// parseRange reads "5" as exactly 5 and "1,5" as 1 to 5 inclusive.
func parseRange(v string) (lo, hi float64, err error) {
	from, to, isRange := strings.Cut(v, ",")
	if lo, err = strconv.ParseFloat(from, 64); err != nil {
		return 0, 0, err
	}
	if !isRange {
		return lo, lo, nil
	}
	hi, err = strconv.ParseFloat(to, 64)
	return lo, hi, err
}

// parseGeoDistance reads GBIF's "lat,lng,distance", where distance ends in
// km, m or mi. The distance may also come first, which is how the starter's
// client sends it.
func parseGeoDistance(v string) (lat, lng, km float64, err error) {
	parts := strings.Split(v, ",")
	if len(parts) != 3 {
		return 0, 0, 0, errors.New("want three parts")
	}
	var coords []float64
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if d, ok := parseDistance(p); ok {
			km = d
			continue
		}
		f, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return 0, 0, 0, err
		}
		coords = append(coords, f)
	}
	if len(coords) != 2 || km == 0 {
		return 0, 0, 0, errors.New("want a lat, lng and distance")
	}
	return coords[0], coords[1], km, nil
}

func parseDistance(v string) (km float64, ok bool) {
	for _, unit := range []struct {
		suffix string
		toKm   float64
	}{{"km", 1}, {"mi", 1.609344}, {"m", 0.001}} {
		if n, found := strings.CutSuffix(v, unit.suffix); found {
			f, err := strconv.ParseFloat(n, 64)
			return f * unit.toKm, err == nil && f > 0
		}
	}
	return 0, false
}

var partialDate = regexp.MustCompile(`^\d{4}(-\d{2}(-\d{2})?)?$`)

// parseDateRange reads "2020", "2020-01", "2020-01-01" or a comma
// separated range of them, where "*" leaves that end open.
func parseDateRange(v string) (from, to string, err error) {
	fromStr, toStr, isRange := strings.Cut(v, ",")
	if !isRange {
		toStr = fromStr
	}
	from, to = "0000-01-01", "9999-12-31"
	if fromStr != "*" {
		if !partialDate.MatchString(fromStr) {
			return "", "", errors.New("invalid date")
		}
		from = firstDay(fromStr)
	}
	if toStr != "*" {
		if !partialDate.MatchString(toStr) {
			return "", "", errors.New("invalid date")
		}
		to = lastDay(toStr)
	}
	return from, to, nil
}

// lastDay pads a partial date to its last day. Day 31 is fine for any
// month, as dates are only ever compared as strings.
func lastDay(date string) string {
	switch len(date) {
	case 4:
		return date + "-12-31"
	case 7:
		return date + "-31"
	}
	return date
}
