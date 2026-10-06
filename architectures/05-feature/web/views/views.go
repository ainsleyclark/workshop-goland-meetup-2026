// Package views holds the templates. They take our domain types
// directly, rather than a second set of view structs to keep in step
// with them.
package views

import (
	_ "embed"
	"html/template"
	"io"

	"05-feature/domain/sighting"
)

//go:embed sightings.html
var pageHTML string
var page = template.Must(template.New("sightings").Parse(pageHTML))

// Sightings writes the sightings page.
func Sightings(w io.Writer, list []sighting.Sighting) error {
	return page.Execute(w, list)
}
