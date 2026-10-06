package services

import (
	"context"
	"log/slog"

	"02-layered/gbif"
	"02-layered/store"
)

// Ingest holds everything needed to pull records from GBIF into storage.
type Ingest struct {
	Client *gbif.Client
	Store  store.Storage
}

// Ingest fetches occurrences from the GBIF API and then saves them to
// storage. Species will be obtained per occurrence.
func (i Ingest) Ingest(ctx context.Context) error {
	occurrences, err := i.Client.Occurrences(ctx)
	if err != nil {
		return err
	}

	for _, o := range occurrences {
		// For every occurrence, we need to fetch its species.
		_, err = i.Store.FindSpecies(o.SpeciesKey)
		if err != nil {
			species, err := i.Client.Species(ctx, o.SpeciesKey)
			if err != nil {
				slog.Error("fetching species: " + err.Error())
				continue
			}
			if err = i.Store.SaveSpecies(species); err != nil {
				slog.Error("saving species: " + err.Error())
			}
		}

		if err = i.Store.SaveOccurrence(o); err != nil {
			return err
		}
	}

	return nil
}
