// Package handlers holds one file per page, and a page is a function.
// It sits beside the views it renders rather than inside a domain,
// because a page is usually about more than one noun and belongs to
// none of them.
package handlers

import (
	"log/slog"
	"net/http"

	"05-feature/domain/sighting"
	"05-feature/web/views"
)

// SightingsLister is all the home page needs to draw itself. Declaring
// it here rather than taking *sighting.Service keeps the arrow pointing
// inwards: web knows the domain, the domain knows nothing of web.
type SightingsLister interface {
	List() []sighting.Sighting
}

// Home serves GET / and shows all sightings.
func Home(sightings SightingsLister) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if err := views.Sightings(w, sightings.List()); err != nil {
			slog.Error("rendering sightings: " + err.Error())
		}
	}
}
