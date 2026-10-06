package sighting

import (
	"context"
	"errors"
	"log/slog"

	"05-feature/clients/gbif"
	"05-feature/domain/species"
)

// Ingest fetches occurrences from the GBIF API and then saves them to
// storage. Species will be obtained per occurrence.
func (s Service) Ingest(ctx context.Context) error {
	occurrences, err := s.gbif.Occurrences(ctx)
	if err != nil {
		return err
	}

	for _, o := range occurrences {
		sp, err := s.resolveSpecies(ctx, o.SpeciesKey)
		if err != nil {
			slog.Error("resolving species", "key", o.SpeciesKey, "err", err)
			continue
		}

		// Save rather than store.Save, so a sighting without
		// coordinates cannot reach the store.
		if err = s.Save(transformOccurrence(o, sp)); err != nil {
			if errors.Is(err, ErrInvalid) || errors.Is(err, ErrAlreadyExists) {
				slog.Warn("skipping occurrence", "key", o.Key, "err", err)
				continue
			}
			return err
		}
	}

	return nil
}

// resolveSpecies returns the species for an occurrence, fetching and
// saving it the first time that species is seen.
func (s Service) resolveSpecies(ctx context.Context, key int) (species.Species, error) {
	sp, err := s.species.Find(key)
	switch {
	case err == nil:
		return sp, nil
	case !errors.Is(err, species.ErrNotFound):
		return species.Species{}, err
	}

	fetched, err := s.gbif.Species(ctx, key)
	if err != nil {
		return species.Species{}, err
	}

	sp = transformSpecies(fetched)
	if err = s.species.Save(sp); err != nil {
		return species.Species{}, err
	}

	return sp, nil
}

func transformOccurrence(o gbif.Occurrence, sp species.Species) Sighting {
	return Sighting{
		Key:            o.Key,
		Species:        sp,
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
