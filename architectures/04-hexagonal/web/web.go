package web

import (
	_ "embed"
	"html/template"
	"io"

	"04-hexagonal/domain/sighting"
	"04-hexagonal/domain/species"
)

//go:embed sightings.html
var pageHTML string
var page = template.Must(template.New("sightings").Parse(pageHTML))

// Row is one line of the sightings page: what was seen, and what it was.
// The template takes our domain types directly, rather than a second set
// of view structs to keep in step with them.
type Row struct {
	Sighting sighting.Sighting
	Species  species.Species
}

// Sightings writes the sightings page.
func Sightings(w io.Writer, rows []Row) error {
	return page.Execute(w, rows)
}
