package services

import (
	"context"
	"log/slog"

	"03-onion/domain"
	"03-onion/gbif"
)

// Ingest holds everything needed to pull records from GBIF into storage.
type Ingest struct {
	Client         *gbif.Client
	SightingsStore domain.SightingStore
	SpeciesStore   domain.SpeciesStore
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
		_, err = i.SpeciesStore.Find(o.SpeciesKey)
		if err != nil {
			species, err := i.Client.Species(ctx, o.SpeciesKey)
			if err != nil {
				slog.Error("fetching species: " + err.Error())
				continue
			}
			if err = i.SpeciesStore.Save(transformSpecies(species)); err != nil {
				slog.Error("saving species: " + err.Error())
			}
		}

		if err = i.SightingsStore.Save(transformOccurrence(o)); err != nil {
			return err
		}
	}

	return nil
}

func transformOccurrence(o gbif.Occurrence) domain.Sighting {
	return domain.Sighting{
		Key:            o.Key,
		SpeciesKey:     o.SpeciesKey,
		ScientificName: o.ScientificName,
		Country:        o.Country,
		Locality:       o.Locality,
		Latitude:       o.Latitude,
		Longitude:      o.Longitude,
		EventDate:      o.EventDate,
	}
}

func transformSpecies(s gbif.Species) domain.Species {
	return domain.Species{
		Key:            s.Key,
		Kingdom:        s.Kingdom,
		Phylum:         s.Phylum,
		Class:          s.Class,
		Order:          s.Order,
		Family:         s.Family,
		Genus:          s.Genus,
		Species:        s.Species,
		Rank:           s.Rank,
		ScientificName: s.ScientificName,
		CanonicalName:  s.CanonicalName,
		VernacularName: s.VernacularName,
	}
}
