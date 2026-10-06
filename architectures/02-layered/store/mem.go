package store

import (
	"fmt"

	"02-layered/models"
)

// MemStore defines an in memory store for occurrence and species data.
type MemStore struct {
	occurrences map[int]models.Occurrence
	species     map[int]models.Species
}

// NewMemStore returns a new in memory datastore.
func NewMemStore() *MemStore {
	return &MemStore{
		occurrences: make(map[int]models.Occurrence),
		species:     make(map[int]models.Species),
	}
}

// SaveOccurrence saves an occurrence to the datastore.
// Returns an error if it exists by key.
func (s *MemStore) SaveOccurrence(o models.Occurrence) error {
	_, ok := s.occurrences[o.Key]
	if ok {
		return fmt.Errorf("occurrence already exists: %d", o.Key)
	}
	s.occurrences[o.Key] = o
	return nil
}

// ListOccurrences returns all occurrences in the datastore.
func (s *MemStore) ListOccurrences() []models.Occurrence {
	var occurrences []models.Occurrence
	for _, o := range s.occurrences {
		occurrences = append(occurrences, o)
	}
	return occurrences
}

// SaveSpecies saves a species to the datastore.
// Returns an error if it exists by key.
func (s *MemStore) SaveSpecies(sp models.Species) error {
	_, ok := s.species[sp.Key]
	if ok {
		return fmt.Errorf("species already exists: %d", sp.Key)
	}
	s.species[sp.Key] = sp
	return nil
}

// FindSpecies returns a species by key.
// Returns an error if it does not exist.
func (s *MemStore) FindSpecies(key int) (models.Species, error) {
	sp, ok := s.species[key]
	if !ok {
		return models.Species{}, fmt.Errorf("species not found: %d", key)
	}
	return sp, nil
}
