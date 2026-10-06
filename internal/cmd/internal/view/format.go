package view

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
	"unicode"
	"uuid"
	"workshop/internal/common/printer/styles"
	"workshop/internal/domain/sighting"
	"workshop/internal/domain/species"
	"workshop/internal/domain/weather"
)

// DateFormat is used whenever a date is shown to the user.
const DateFormat = "2006-01-02"

// Ptr formats an optional value, returning an empty string when it
// isn't set so key/value views can leave the row out.
func Ptr[T any](v *T, format string) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf(format, *v)
}

// ShortID returns the first block of a UUID, which is
// enough to identify a record in a local database.
func ShortID(id uuid.UUID) string {
	return id.String()[:8]
}

// SpeciesName returns the friendliest name available,
// e.g. "Humpback Whale (Megaptera novaeangliae)".
func SpeciesName(s species.Species) string {
	scientific := cmp.Or(s.CanonicalName, s.ScientificName)
	vernacular := s.VernacularName
	if vernacular == strings.ToLower(vernacular) {
		vernacular = titleCase(vernacular)
	}
	switch {
	case vernacular != "" && scientific != "":
		return fmt.Sprintf("%s (%s)", vernacular, scientific)
	case vernacular != "":
		return vernacular
	case scientific != "":
		return scientific
	default:
		return "Unknown species"
	}
}

// Location joins the human-readable parts of a sighting's location.
func Location(l sighting.Location) string {
	parts := make([]string, 0, 3)
	for _, p := range []string{l.Locality, l.StateProvince, l.Country} {
		p = tidy(strings.TrimSpace(p))
		// GBIF often repeats the country or state in the locality.
		if p != "" && !slices.ContainsFunc(parts, func(v string) bool { return strings.EqualFold(v, p) }) {
			parts = append(parts, p)
		}
	}
	if len(parts) == 0 {
		return Coordinates(l)
	}
	return strings.Join(parts, ", ")
}

// Coordinates formats a location's latitude and longitude.
func Coordinates(l sighting.Location) string {
	out := fmt.Sprintf("%.4f, %.4f", l.Latitude, l.Longitude)
	if l.CoordinateUncertaintyInMetres != nil {
		out += fmt.Sprintf(" (±%.0fm)", *l.CoordinateUncertaintyInMetres)
	}
	return out
}

// Weather summarises a sighting's recorded conditions, or flags that
// none has been saved yet. weatherPerSighting is a stub until the
// ingest exercise is done, so freshly ingested sightings have a
// zero weather ID rather than a saved record.
func Weather(w weather.Weather) string {
	if w.ID == uuid.Nil() {
		return styles.Muted.Render(styles.IconWarn + " missing")
	}
	return styles.Success.Render(fmt.Sprintf("%s %.0f°C", styles.IconSuccess, w.Temperature.Actual))
}

// Humanise turns GBIF enum values such as HUMAN_OBSERVATION
// into "Human observation".
func Humanise(s string) string {
	if s == "" {
		return ""
	}
	s = strings.ToLower(strings.ReplaceAll(s, "_", " "))
	return strings.ToUpper(s[:1]) + s[1:]
}

// tidy converts SHOUTED place names from GBIF into title case,
// leaving anything that already uses mixed case alone.
func tidy(s string) string {
	if s == strings.ToUpper(s) && s != strings.ToLower(s) {
		return titleCase(s)
	}
	return s
}

// titleCase upper-cases the first letter of every word, treating
// spaces and hyphens as word boundaries.
func titleCase(s string) string {
	r := []rune(strings.ToLower(s))
	for i := range r {
		if i == 0 || r[i-1] == ' ' || r[i-1] == '-' {
			r[i] = unicode.ToUpper(r[i])
		}
	}
	return string(r)
}

// MediaItems renders each media URL with its creator and licence,
// where GBIF gave us one, ready for a List.
func MediaItems(items []sighting.Media) []string {
	out := make([]string, 0, len(items))
	for _, m := range items {
		meta := strings.Join(slices.DeleteFunc([]string{m.Creator, m.License}, func(v string) bool {
			return v == ""
		}), " "+styles.IconDot+" ")

		item := m.URL
		if meta != "" {
			item += " " + styles.Muted.Render(meta)
		}
		out = append(out, item)
	}
	return out
}
