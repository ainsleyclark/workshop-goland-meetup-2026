package handlers

import (
	_ "embed"
	"html/template"
	"log"
	"log/slog"
	"net/http"

	"03-onion/domain"
)

//go:embed sightings.html
var pageHTML string
var page = template.Must(template.New("sightings").Parse(pageHTML))

// Homepage serves GET / and shows all sightings.
func Homepage(sightings domain.SightingStore, species domain.SpeciesStore) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		list := sightings.List()

		type pageRow struct {
			Sighting domain.Sighting
			Species  domain.Species
		}

		// Obtain the species for each sighting.
		rows := make([]pageRow, 0, len(list))
		for _, s := range list {
			sp, err := species.Find(s.SpeciesKey)
			if err != nil {
				slog.Error("failed to find species", "species_key", s.SpeciesKey, "err", err)
				continue
			}
			rows = append(rows, pageRow{Sighting: s, Species: sp})
		}

		if err := page.Execute(w, rows); err != nil {
			log.Println(err)
		}
	}
}
