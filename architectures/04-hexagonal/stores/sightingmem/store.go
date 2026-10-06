package sightingmem

import (
	"fmt"

	"04-hexagonal/domain/sighting"
)

// Store defines an in memory, sightings store.
type Store struct {
	data map[int]sighting.Sighting
}

var _ sighting.Store = (*Store)(nil)

// New returns a new in memory, sightings store.
func New() *Store {
	return &Store{
		data: make(map[int]sighting.Sighting),
	}
}

// Save saves a sighting to the datastore.
// Returns an error if it exists by key.
func (s *Store) Save(si sighting.Sighting) error {
	_, ok := s.data[si.Key]
	if ok {
		return fmt.Errorf("%w: %d", sighting.ErrAlreadyExists, si.Key)
	}
	s.data[si.Key] = si
	return nil
}

// List returns all sightings in the datastore.
func (s *Store) List() []sighting.Sighting {
	var out []sighting.Sighting
	for _, si := range s.data {
		out = append(out, si)
	}
	return out
}
