package main

import (
	"fmt"
)

// MemStore defines an in memory store for occurrence and species data.
type MemStore struct {
	occurrences map[int]Occurrence
	species     map[int]Species
}

// NewStore returns a new in memory datastore.
func NewStore() *MemStore {
	return &MemStore{
		occurrences: make(map[int]Occurrence),
		species:     make(map[int]Species),
	}
}

// SaveOccurrence saves an occurrence to the datastore.
// Returns an error if it exists by key.
func (s *MemStore) SaveOccurrence(o Occurrence) error {
	_, ok := s.occurrences[o.Key]
	if ok {
		return fmt.Errorf("occurrence already exists: %d", o.Key)
	}
	s.occurrences[o.Key] = o
	return nil
}

// ListOccurrences returns all occurrences in the datastore.
func (s *MemStore) ListOccurrences() []Occurrence {
	var occurrences []Occurrence
	for _, o := range s.occurrences {
		occurrences = append(occurrences, o)
	}
	return occurrences
}

// SaveSpecies saves a species to the datastore.
// Returns an error if it exists by key.
func (s *MemStore) SaveSpecies(sp Species) error {
	_, ok := s.species[sp.Key]
	if ok {
		return fmt.Errorf("species already exists: %d", sp.Key)
	}
	s.species[sp.Key] = sp
	return nil
}

// FindSpecies returns a species by key.
// Returns an error if it does not exist.
func (s *MemStore) FindSpecies(key int) (Species, error) {
	sp, ok := s.species[key]
	if !ok {
		return Species{}, fmt.Errorf("species not found: %d", key)
	}
	return sp, nil
}
