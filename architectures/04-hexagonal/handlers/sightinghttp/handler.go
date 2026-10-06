package sightinghttp

import (
	"log/slog"
	"net/http"

	"04-hexagonal/domain/sighting"
	"04-hexagonal/domain/species"
	"04-hexagonal/web"
)

// Handler serves the sighting pages.
type Handler struct {
	sightings *sighting.Service
	species   *species.Service
}

// New returns a Handler backed by the given services.
func New(sightings *sighting.Service, speciesSvc *species.Service) *Handler {
	return &Handler{
		sightings: sightings,
		species:   speciesSvc,
	}
}

// Routes registers the handler's routes, so the wiring in main
// doesn't need to know any URLs.
func (h *Handler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.homepage)
}

// homepage serves GET / and shows all sightings.
func (h *Handler) homepage(w http.ResponseWriter, _ *http.Request) {
	sightings := h.sightings.List()

	// Obtain the species for each sighting.
	rows := make([]web.Row, 0, len(sightings))
	for _, s := range sightings {
		sp, err := h.species.Find(s.SpeciesKey)
		if err != nil {
			slog.Error("failed to find species", "species_key", s.SpeciesKey, "err", err)
			continue
		}
		rows = append(rows, web.Row{Sighting: s, Species: sp})
	}

	if err := web.Sightings(w, rows); err != nil {
		slog.Error("rendering sightings: " + err.Error())
	}
}
