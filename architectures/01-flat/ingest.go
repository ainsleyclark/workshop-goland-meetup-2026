package main

import (
	"context"
	"log/slog"
)

// Ingest fetches occurrences from the GBIF API and then saves them to
// storage. Species will be obtained per occurrence.
func Ingest(ctx context.Context, client *GBIF, store *MemStore) error {
	occurrences, err := client.Occurrences(ctx)
	if err != nil {
		return err
	}

	for _, o := range occurrences {
		// For every occurrence, we need to fetch its species.
		_, err = store.FindSpecies(o.SpeciesKey)
		if err != nil {
			species, err := client.Species(ctx, o.SpeciesKey)
			if err != nil {
				slog.Error("fetching species: " + err.Error())
				continue
			}
			if err = store.SaveSpecies(species); err != nil {
				slog.Error("saving species: " + err.Error())
			}
		}

		if err = store.SaveOccurrence(o); err != nil {
			return err
		}
	}

	return nil
}
