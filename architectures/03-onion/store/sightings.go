package store

import (
	"fmt"

	"03-onion/domain"
)

// Sighting defines an in memory, sightings store.
type Sighting struct {
	data map[int]domain.Sighting
}

var _ domain.SightingStore = (*Sighting)(nil)

// NewSightingsStore returns a new in memory, sightings store.
func NewSightingsStore() *Sighting {
	return &Sighting{
		data: make(map[int]domain.Sighting),
	}
}

// Save saves a sighting to the datastore.
// Returns an error if it exists by key.
func (s *Sighting) Save(si domain.Sighting) error {
	_, ok := s.data[si.Key]
	if ok {
		return fmt.Errorf("sighting already exists: %d", si.Key)
	}
	s.data[si.Key] = si
	return nil
}

// List returns all sightings in the datastore.
func (s *Sighting) List() []domain.Sighting {
	var out []domain.Sighting
	for _, o := range s.data {
		out = append(out, o)
	}
	return out
}
