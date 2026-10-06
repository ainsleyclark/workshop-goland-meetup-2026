package web

import (
	"net/http"

	"05-feature/web/handlers"
)

// New returns a ServeMux serving every page. It is the one place that
// knows a URL, so no handler knows another and main.go knows no paths.
func New(sightings handlers.SightingsLister) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /", handlers.Home(sightings))

	return mux
}
