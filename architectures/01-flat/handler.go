package main

import (
	_ "embed"
	"html/template"
	"log"
	"log/slog"
	"net/http"
)

//go:embed sightings.html
var pageHTML string
var page = template.Must(template.New("sightings").Parse(pageHTML))

// Homepage serves GET / and shows all occurrences.
func Homepage(store *MemStore) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		occurrences := store.ListOccurrences()

		type pageRow struct {
			Occurrence Occurrence
			Species    Species
		}

		// Obtain the species for each occurrence.
		rows := make([]pageRow, 0, len(occurrences))
		for _, o := range occurrences {
			species, err := store.FindSpecies(o.SpeciesKey)
			if err != nil {
				slog.Error("failed to find species", "species_key", o.SpeciesKey, "err", err)
				continue
			}
			rows = append(rows, pageRow{Occurrence: o, Species: species})
		}

		if err := page.Execute(w, rows); err != nil {
			log.Println(err)
		}
	}
}
