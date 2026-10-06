package sighting

import (
	"context"
	"errors"
	"log/slog"

	"04-hexagonal/clients/gbif"
	"04-hexagonal/domain/species"
)

// Ingest fetches occurrences from the GBIF API and then saves them to
// storage. Species will be obtained per occurrence.
func (s Service) Ingest(ctx context.Context) error {
	occurrences, err := s.gbif.Occurrences(ctx)
	if err != nil {
		return err
	}

	for _, o := range occurrences {
		// For every occurrence, we need to fetch its species.
		_, err = s.species.Find(o.SpeciesKey)
		if errors.Is(err, species.ErrNotFound) {
			sp, err := s.gbif.Species(ctx, o.SpeciesKey)
			if err != nil {
				slog.Error("fetching species: " + err.Error())
				continue
			}
			if err = s.species.Save(transformSpecies(sp)); err != nil {
				slog.Error("saving species: " + err.Error())
			}
		} else if err != nil {
			return err
		}

		// Save rather than store.Save, so a sighting without
		// coordinates cannot reach the store.
		if err = s.Save(transformOccurrence(o)); err != nil {
			if errors.Is(err, ErrInvalid) || errors.Is(err, ErrAlreadyExists) {
				slog.Warn("skipping occurrence", "key", o.Key, "err", err)
				continue
			}
			return err
		}
	}

	return nil
}

func transformOccurrence(o gbif.Occurrence) Sighting {
	return Sighting{
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

func transformSpecies(s gbif.Species) species.Species {
	return species.Species{
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
